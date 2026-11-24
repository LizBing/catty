package rt

import (
	"fmt"
	"sync/atomic"
)

// Thread is one logical Java thread. Each live execution owns its frames and
// locals exclusively; the Java peer, lifecycle state, interrupt status and
// monitor/blocker state are shared and guarded by Runtime.thrMu.
type Thread struct {
	RT     *Runtime
	ID     int64
	depth  int
	frames []FrameInfo

	// peer is the java.lang.Thread object for this logical thread. It is
	// created at Thread construction or by Thread.currentThread for a host
	// thread that entered Java without a peer.
	peer *Object

	// engineStack is the thread's stack of non-intrinsic execution engines.
	// Intrinsic frames are adapters for cross-engine accounting, so they are
	// not pushed; a new logical thread starts with an empty stack.
	engineStack []string

	// The following fields are guarded by RT.thrMu.
	blocker       *threadBlocker
	hostInterrupt bool // pending interrupt for a host thread without a peer
	// held maps each object whose monitor this logical thread owns at depth
	// one or more to that monitor. Java releases every monitor when a thread
	// terminates, even if hand-authored code failed to exit it. The map roots
	// only objects whose monitor is currently held; object-lifetime monitor
	// state itself is stored on the object.
	held map[*Object]*monitor
}

// interruptPendingLocked reports the logical thread's Java interrupt status.
// The caller holds thrMu.
func (t *Thread) interruptPendingLocked() bool {
	if p, ok := peerOfObject(t.peer); ok {
		return p.pendingInterrupt
	}
	return t.hostInterrupt
}

// setInterruptLocked records the Java interrupt status. The caller holds thrMu.
func (t *Thread) setInterruptLocked() {
	if p, ok := peerOfObject(t.peer); ok {
		p.pendingInterrupt = true
		return
	}
	t.hostInterrupt = true
}

// clearInterruptLocked consumes the Java interrupt status. The caller holds
// thrMu.
func (t *Thread) clearInterruptLocked() {
	if p, ok := peerOfObject(t.peer); ok {
		p.pendingInterrupt = false
		return
	}
	t.hostInterrupt = false
}

const maxJavaDepth = 20000

// NewThread creates a logical Java thread with a fresh runtime-local identity.
// The returned thread has no Java peer until it is bound by Thread.start, the
// launcher, or Thread.currentThread.
func (r *Runtime) NewThread() *Thread {
	id := atomic.AddInt64(&r.nextThread, 1)
	return &Thread{RT: r, ID: id}
}

// MainThread returns the launcher thread, creating it and its canonical "main"
// peer on first use.
func (r *Runtime) MainThread() *Thread {
	r.thrMu.Lock()
	if r.mainThread != nil {
		t := r.mainThread
		r.thrMu.Unlock()
		return t
	}
	r.thrMu.Unlock()
	var peer *Object
	if r.profile != nil && r.profile.Threads {
		peer = r.newPeerObject("main", false, threadStarted)
	}
	t := r.NewThread()
	r.thrMu.Lock()
	if r.mainThread == nil {
		r.mainThread = t
		if p, ok := peerOfObject(peer); ok {
			p.thread = t
		}
		t.peer = peer
		r.noteLiveThreadStartLocked()
	}
	result := r.mainThread
	r.thrMu.Unlock()
	return result
}

// SetTracePath opens the opt-in JSON-lines trace. It is used by RunApp.
func (r *Runtime) SetTracePath(path string) error {
	tr, err := openTrace(path)
	if err != nil {
		return err
	}
	r.trace = tr
	return nil
}

// CloseTrace flushes and closes the trace channel. The trace object stays
// reachable so a daemon that is still running can safely emit into a closed
// channel: emit observes the closed file under the trace mutex and drops the
// record.
func (r *Runtime) CloseTrace() {
	if tr := r.trace; tr != nil {
		tr.close()
	}
}

// Call is the one owner of Java frame push/pop and trace enter/exit for every
// engine. The selected method's entry runs synchronously on the same thread.
// A synchronized method acquires its implicit monitor after resolution,
// receiver/access checks and active-use initialization have already happened in
// the caller, and releases exactly one acquisition on Java return or unwind
// before the outcome reaches the caller. The synchronization flag on <clinit>
// is ignored for execution.
func (t *Thread) Call(m *Method, receiver *Object, args []Value) Outcome {
	if m == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if t.depth >= maxJavaDepth {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/StackOverflowError"))
	}
	var lock *Object
	if m.Flags&AccSynchronized != 0 && m.Name != "<clinit>" {
		if m.Flags&AccStatic != 0 {
			lock = m.Owner.Mirror
		} else {
			lock = receiver
		}
		if lock == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
		}
		if thrown := t.monitorEnterObject(lock, m, -1); thrown != nil {
			return ThrowOutcome(thrown)
		}
	}
	t.depth++
	t.frames = append(t.frames, FrameInfo{Method: m})
	depth := t.depth - 1
	if tr := t.trace(); tr != nil {
		tr.enter(t.ID, m, depth)
	}
	if hook := t.RT.HostFaultHook; hook != nil {
		hook(t, m)
	}
	if metrics := t.RT.metrics; metrics != nil && (m.Engine == "aot" || m.Engine == "interp") {
		if n := len(t.engineStack); n > 0 && t.engineStack[n-1] != m.Engine {
			metrics.crossEngineCall()
		}
		t.engineStack = append(t.engineStack, m.Engine)
		defer func() { t.engineStack = t.engineStack[:len(t.engineStack)-1] }()
	}
	var out Outcome
	switch m.Engine {
	case "aot":
		if m.Entry == nil {
			out = ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NoSuchMethodError"))
		} else {
			out = m.Entry(t, m, receiver, args)
		}
	case "interp":
		out = t.RT.execute(t, m, receiver, args)
	case "intrinsic":
		if m.Intrinsic == nil {
			out = ThrowOutcome(t.RT.NewThrowable(t, "java/lang/AbstractMethodError"))
		} else {
			out = m.Intrinsic(t, m, receiver, args)
		}
	case "abstract":
		out = ThrowOutcome(t.RT.NewThrowable(t, "java/lang/AbstractMethodError"))
	default:
		out = ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NoSuchMethodError"))
	}
	if t.RT.metrics != nil && out.Thrown != nil {
		t.RT.metrics.javaException()
	}
	if lock != nil {
		// Release exactly one implicit acquisition. A body that already
		// relinquished ownership is reported at the method boundary.
		if thrown := t.monitorExitObject(lock, m, -1); thrown != nil {
			out = ThrowOutcome(thrown)
		}
	}
	if tr := t.trace(); tr != nil {
		if out.Thrown != nil {
			tr.unwind(t.ID, m, depth, out.Thrown)
		} else {
			tr.exit(t.ID, m, depth)
		}
	}
	t.frames = t.frames[:len(t.frames)-1]
	t.depth--
	return out
}

// Call executes m on t through the common ABI.
func (r *Runtime) Call(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return t.Call(m, receiver, args)
}

// SetPC records the caller's current bytecode offset immediately before a
// potentially throwing instruction.
func (t *Thread) SetPC(pc int32) {
	if n := len(t.frames); n > 0 {
		t.frames[n-1].PC = pc
	}
}

// PC returns the current bytecode offset of the active frame.
func (t *Thread) PC() int32 {
	if n := len(t.frames); n > 0 {
		return t.frames[n-1].PC
	}
	return 0
}

// CurrentMethod returns the active method, if any.
func (t *Thread) CurrentMethod() *Method {
	if n := len(t.frames); n > 0 {
		return t.frames[n-1].Method
	}
	return nil
}

// CallerMethod returns the logical Java caller of the active frame. Runtime
// intrinsics such as Method.invoke use it for caller-sensitive access checks.
func (t *Thread) CallerMethod() *Method {
	if n := len(t.frames); n >= 2 {
		return t.frames[n-2].Method
	}
	return nil
}

func (r *Runtime) bootstrapThrowable(class string) *Object {
	return r.NewThrowable(nil, class)
}

// ---------------------------------------------------------------- class init

// EnsureInitialized performs the JVMS 5.5 initialization procedure for c. It
// returns a Java throwable on failure and nil on success. Each defining class
// owns its own initialization guard, owning logical thread and completion
// signal. A competing thread copies the completion signal and waits outside the
// guard, then rechecks state; its interrupt status is untouched. Reentry by the
// owning thread succeeds with the current static values.
func (t *Thread) EnsureInitialized(c *Class) *Object {
	if c == nil || c.IsPrimitive || c.IsArray {
		return nil
	}
	for {
		c.initMu.Lock()
		switch c.State {
		case classInitialized:
			c.initMu.Unlock()
			return nil
		case classErroneous:
			c.initMu.Unlock()
			// Every later active use receives a fresh error, never a rerun.
			return t.RT.NewThrowable(t, "java/lang/NoClassDefFoundError")
		case classInitializing:
			if c.initOwner == t {
				c.initMu.Unlock()
				return nil
			}
			done, owner := c.initDone, c.initOwner
			// Emit the competing wait while still holding the initialization
			// guard, so the owner cannot publish completion first and the
			// trace demonstrates start < wait < complete.
			if tr := t.trace(); tr != nil {
				tr.classInitWait(t.ID, c, owner)
			}
			c.initMu.Unlock()
			if hook := t.RT.ClassInitWaitHook; hook != nil {
				hook(c, t)
			}
			<-done
			continue
		}
		// First active use: publish the initializing state and owner, and
		// record the start while the guard is held so no competing waiter can
		// precede it in the trace.
		c.State = classInitializing
		c.initOwner = t
		c.initDone = make(chan struct{})
		if tr := t.trace(); tr != nil {
			tr.classInit(t.ID, c, true, nil)
		}
		c.initMu.Unlock()
		err := t.runClassInitialization(c)
		c.initMu.Lock()
		if err != nil {
			c.State = classErroneous
			c.InitError = err
		} else {
			c.State = classInitialized
		}
		done := c.initDone
		c.initOwner = nil
		// Publish the completion observation before any waiter can proceed.
		if tr := t.trace(); tr != nil {
			tr.classInit(t.ID, c, false, err)
		}
		c.initMu.Unlock()
		close(done)
		return err
	}
}

// runClassInitialization executes JVMS 5.5 steps 6-9 outside the guard: apply
// ConstantValue attributes in declaration order through the shared memory path,
// initialize the required superclass and default-method interfaces in order,
// then invoke <clinit> through Thread.Call. The synchronization flag on
// <clinit> is ignored for execution.
func (t *Thread) runClassInitialization(c *Class) *Object {
	t.RT.applyConstants(c)
	if !c.IsInterface() {
		if c.Super != nil {
			if err := t.EnsureInitialized(c.Super); err != nil {
				return err
			}
		}
		for _, in := range t.RT.superInterfacesWithDefaults(c) {
			if err := t.EnsureInitialized(in); err != nil {
				return err
			}
		}
	}
	if clinit, ok := c.MethodIndex["<clinit>()V"]; ok {
		out := t.Call(clinit, nil, nil)
		if out.Thrown != nil {
			err := out.Thrown
			if !t.RT.IsAssignable(err.Class, t.RT.mustClass("java/lang/Error")) {
				wrapper := t.RT.NewThrowable(t, "java/lang/ExceptionInInitializerError")
				SetCause(wrapper, err)
				err = wrapper
			}
			return err
		}
	}
	return nil
}

// applyConstants installs ConstantValue attributes during initialization. The
// builder has already validated every constant against its field descriptor.
// The whole sequence is one ordered memory operation so a concurrently
// publishing initializer observes either none or all of the preceding writes,
// never a torn intermediate state.
func (r *Runtime) applyConstants(c *Class) {
	r.memMu.Lock()
	defer r.memMu.Unlock()
	for i, f := range c.StaticField {
		if f.Constant == nil || i >= len(c.Statics) {
			continue
		}
		if v, err := r.constantValue(f.Constant); err == nil {
			c.Statics[i] = v
		}
	}
}

// superInterfacesWithDefaults returns all superinterfaces of a class that
// declare at least one non-abstract, non-static method, in the recursive
// enumeration order required by JVMS 5.5.
func (r *Runtime) superInterfacesWithDefaults(c *Class) []*Class {
	var out []*Class
	seen := map[*Class]bool{}
	var visit func(in *Class)
	visit = func(in *Class) {
		for _, sup := range in.Interfaces {
			visit(sup)
		}
		if !seen[in] && declaresDefault(in) {
			seen[in] = true
			out = append(out, in)
		}
	}
	for _, in := range c.Interfaces {
		visit(in)
	}
	return out
}

func declaresDefault(c *Class) bool {
	if c == nil || !c.IsInterface() {
		return false
	}
	for _, m := range c.Methods {
		if m.Name == "<clinit>" || m.Name == "<init>" {
			continue
		}
		if m.Flags&(AccStatic|AccAbstract|AccPrivate) == 0 {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------- resolution

func methodKey(name, descriptor string) string { return name + descriptor }

func (r *Runtime) findDeclaredMethod(owner *Class, name, descriptor string) *Method {
	key := methodKey(name, descriptor)
	for c := owner; c != nil; c = c.Super {
		if m := c.MethodIndex[key]; m != nil {
			return m
		}
	}
	return nil
}

// findInterfaceMethodClosure returns the first matching method declared by an
// interface in owner's superinterface closure.
func (r *Runtime) findInterfaceMethodClosure(owner *Class, key string) *Method {
	seen := map[*Class]bool{}
	var visit func(c *Class) *Method
	visit = func(c *Class) *Method {
		if c == nil || seen[c] {
			return nil
		}
		seen[c] = true
		if m := c.MethodIndex[key]; m != nil && c.IsInterface() {
			return m
		}
		for _, in := range c.Interfaces {
			if m := visit(in); m != nil {
				return m
			}
		}
		if !c.IsInterface() {
			return visit(c.Super)
		}
		return nil
	}
	return visit(owner)
}

// findMethod is the method-resolution lookup (JVMS 5.4.3.3 and 5.4.3.4).
func (r *Runtime) findMethod(owner *Class, name, descriptor string) *Method {
	key := methodKey(name, descriptor)
	if m := r.findDeclaredMethod(owner, name, descriptor); m != nil {
		return m
	}
	if m := r.findInterfaceMethodClosure(owner, key); m != nil {
		return m
	}
	return nil
}

// findField is field resolution (JVMS 5.4.3.2).
func (r *Runtime) findField(owner *Class, name, descriptor string) *Field {
	for c := owner; c != nil; c = c.Super {
		for _, f := range c.StaticField {
			if f.Name == name && f.Descriptor == descriptor {
				return f
			}
		}
		for _, f := range c.InstanceField {
			if f.Name == name && f.Descriptor == descriptor {
				return f
			}
		}
	}
	// Interface fields are static and inherited.
	seen := map[*Class]bool{}
	var visit func(c *Class) *Field
	visit = func(c *Class) *Field {
		if c == nil || seen[c] {
			return nil
		}
		seen[c] = true
		for _, f := range c.StaticField {
			if f.Name == name && f.Descriptor == descriptor {
				return f
			}
		}
		for _, in := range c.Interfaces {
			if f := visit(in); f != nil {
				return f
			}
		}
		return nil
	}
	return visit(owner)
}

func (r *Runtime) samePackage(a, b *Class) bool {
	if a == nil || b == nil {
		return false
	}
	// JVMS runtime package identity is the package string plus defining loader.
	return a.Loader == b.Loader && a.PackageName() == b.PackageName()
}

func nestHostOf(c *Class) string {
	if c == nil {
		return ""
	}
	if c.NestHost != "" {
		return c.NestHost
	}
	return c.Name
}

func (r *Runtime) sameNest(a, b *Class) bool {
	if a == nil || b == nil || a.Loader != b.Loader {
		return false
	}
	if a == b {
		return true
	}
	if nestHostOf(a) != nestHostOf(b) {
		return false
	}
	// Mutual membership must be present in the defining class metadata.
	hostName := nestHostOf(a)
	if a.Name == hostName {
		return hasNestMember(a, b.Name)
	}
	if b.Name == hostName {
		return hasNestMember(b, a.Name)
	}
	host, ok := a.Loader.Lookup(hostName)
	if !ok {
		return false
	}
	return hasNestMember(host, a.Name) && hasNestMember(host, b.Name)
}

func hasNestMember(host *Class, name string) bool {
	for _, n := range host.NestMembers {
		if n == name {
			return true
		}
	}
	return false
}

func (r *Runtime) isSubclass(c, of *Class) bool {
	for cur := c; cur != nil; cur = cur.Super {
		if cur == of {
			return true
		}
	}
	return false
}

func (r *Runtime) classAccessible(caller, target *Class) bool {
	if target == nil || target.Flags&AccPublic != 0 {
		return true
	}
	return r.samePackage(caller, target) || r.sameNest(caller, target)
}

func (r *Runtime) memberAccessible(caller *Method, owner *Class, flags uint16) bool {
	var callerClass *Class
	if caller != nil {
		callerClass = caller.Owner
	}
	if callerClass == nil {
		// Host-side calls (launcher, tests) act as privileged callers.
		return true
	}
	if !r.classAccessible(callerClass, owner) {
		return false
	}
	switch {
	case flags&AccPublic != 0:
		return true
	case flags&AccPrivate != 0:
		return r.sameNest(callerClass, owner)
	case flags&AccProtected != 0:
		return r.samePackage(callerClass, owner) || r.isSubclass(callerClass, owner) || r.sameNest(callerClass, owner)
	default:
		return r.samePackage(callerClass, owner) || r.sameNest(callerClass, owner)
	}
}

func (r *Runtime) protectedReceiverOK(caller *Method, owner *Class, receiver *Class) bool {
	if caller == nil || caller.Owner == nil || owner == nil {
		return true
	}
	if owner.Flags&AccProtected == 0 {
		return true
	}
	if r.samePackage(caller.Owner, owner) || r.sameNest(caller.Owner, owner) {
		return true
	}
	if !r.isSubclass(caller.Owner, owner) {
		return false
	}
	return receiver == nil || r.isSubclass(receiver, caller.Owner)
}

// resolveClassSymbol resolves a class operand and caches the outcome by the
// referring runtime class and constant-pool entry. A failed lookup is sticky
// just like a failed field or method resolution. Resolution, initiation
// recording and the cache transition are one metadata transaction.
func (t *Thread) resolveClassSymbol(caller *Method, cp int, name string) (*Class, *Object) {
	key := resolutionKey{Kind: 'C', CP: cp, Owner: name}
	if caller != nil && caller.Owner != nil {
		key.Class = caller.Owner
	}
	r := t.RT
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	if res, ok := r.resolutions[key]; ok {
		return res.Class, res.Thrown
	}
	loader := loaderOf(caller)
	c, err := r.classForLoaderNameLocked(loader, name)
	if err != nil {
		return nil, t.cacheResolutionFailure(key, "java/lang/NoClassDefFoundError")
	}
	if r.recordClassInitiationLocked(loader, name, c) {
		return nil, t.cacheResolutionFailure(key, "java/lang/LinkageError")
	}
	r.resolutions[key] = &Resolution{Class: c}
	return c, nil
}

// cacheResolutionFailure allocates a fresh throwable while holding metaMu,
// records it as the sticky outcome and returns it. newThrowableLocked performs
// no Java execution and acquires no other runtime lock.
func (t *Thread) cacheResolutionFailure(key resolutionKey, className string) *Object {
	thrown := t.RT.newThrowableLocked(t, className)
	t.RT.resolutions[key] = &Resolution{Thrown: thrown}
	return thrown
}

// resolveMethodSymbol resolves a method reference and caches the outcome for the
// referring runtime class and constant-pool entry. Access control is applied at
// resolution time.
func (t *Thread) resolveMethodSymbol(caller *Method, cp int, ownerName, name, descriptor string) (*Method, *Object) {
	key := resolutionKey{Kind: 'M', CP: cp, Owner: ownerName, Name: name, Descriptor: descriptor}
	if caller != nil && caller.Owner != nil {
		key.Class = caller.Owner
	}
	r := t.RT
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	if res, ok := r.resolutions[key]; ok {
		return res.Method, res.Thrown
	}
	fail := func(className string) (*Method, *Object) {
		thrown := r.newThrowableLocked(t, className)
		r.resolutions[key] = &Resolution{Thrown: thrown}
		return nil, thrown
	}
	loader := r.Loader
	if caller != nil && caller.Owner != nil && caller.Owner.Loader != nil {
		loader = caller.Owner.Loader
	}
	owner, err := r.classForLoaderNameLocked(loader, ownerName)
	if err != nil {
		return fail("java/lang/NoClassDefFoundError")
	}
	if r.recordInitiationLocked(loader, ownerName, owner) {
		return fail("java/lang/LinkageError")
	}
	m := r.findMethod(owner, name, descriptor)
	if m == nil {
		return fail("java/lang/NoSuchMethodError")
	}
	if !r.memberAccessible(caller, m.Owner, m.Flags) {
		return fail("java/lang/IllegalAccessError")
	}
	if r.imposeMethodConstraintsLocked(loader, m.Owner.Loader, m.Type) {
		return fail("java/lang/LinkageError")
	}
	r.resolutions[key] = &Resolution{Method: m}
	return m, nil
}

// resolveFieldSymbol resolves a field reference with the same caching rule.
func (t *Thread) resolveFieldSymbol(caller *Method, cp int, ownerName, name, descriptor string) (*Field, *Object) {
	key := resolutionKey{Kind: 'F', CP: cp, Owner: ownerName, Name: name, Descriptor: descriptor}
	if caller != nil && caller.Owner != nil {
		key.Class = caller.Owner
	}
	r := t.RT
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	if res, ok := r.resolutions[key]; ok {
		return res.Field, res.Thrown
	}
	fail := func(className string) (*Field, *Object) {
		thrown := r.newThrowableLocked(t, className)
		r.resolutions[key] = &Resolution{Thrown: thrown}
		return nil, thrown
	}
	loader := r.Loader
	if caller != nil && caller.Owner != nil && caller.Owner.Loader != nil {
		loader = caller.Owner.Loader
	}
	owner, err := r.classForLoaderNameLocked(loader, ownerName)
	if err != nil {
		return fail("java/lang/NoClassDefFoundError")
	}
	if r.recordInitiationLocked(loader, ownerName, owner) {
		return fail("java/lang/LinkageError")
	}
	f := r.findField(owner, name, descriptor)
	if f == nil {
		return fail("java/lang/NoSuchFieldError")
	}
	if !r.memberAccessible(caller, f.Owner, f.Flags) {
		return fail("java/lang/IllegalAccessError")
	}
	if r.imposeFieldConstraintsLocked(loader, f.Owner.Loader, f.Type) {
		return fail("java/lang/LinkageError")
	}
	r.resolutions[key] = &Resolution{Field: f}
	return f, nil
}

// ResolveMethod resolves owner/name/descriptor through the class hierarchy
// without an access-checking caller. It is used by host launchers and tests.
func (r *Runtime) ResolveMethod(owner, name, descriptor string) (*Method, error) {
	c, err := r.classForName(owner)
	if err != nil {
		return nil, fmt.Errorf("catty: cannot resolve class %s", owner)
	}
	m := r.findMethod(c, name, descriptor)
	if m == nil {
		return nil, fmt.Errorf("catty: cannot resolve method %s.%s%s", owner, name, descriptor)
	}
	return m, nil
}

// selectVirtual performs receiver-based method selection for invokevirtual and
// invokeinterface after resolution has produced a declaration.
//
// JVMS 5.4.6 searches the receiver class and then its superclass chain for an
// instance method that can override the resolved declaration. The search is
// scoped by receiver type, so the hierarchy relation itself is implicit: an
// inherited superclass implementation can satisfy an interface declaration
// even when that superclass does not implement the interface. If no class
// declaration matches, the maximally-specific superinterface defaults are
// considered.
func (r *Runtime) selectVirtual(receiver *Class, declaration *Method) *Method {
	return r.selectedMethod(receiver, declaration)
}

// selectedMethod performs the receiver-scoped JVMS 5.4.6 search described for
// selectVirtual without materializing the AbstractMethodError diagnostic. It
// returns nil when no concrete implementation is available. Preparation uses
// this form so it can inspect the selected implementation without running
// bootstrap code, creating heap objects, resolving descriptor leaves or
// initializing classes.
func (r *Runtime) selectedMethod(receiver *Class, declaration *Method) *Method {
	key := methodKey(declaration.Name, declaration.Descriptor)
	for c := receiver; c != nil; c = c.Super {
		if m := c.MethodIndex[key]; m != nil {
			if m == declaration && m.Flags&AccStatic == 0 {
				// A private nestmate or same-class invocation resolves to the
				// declaration itself; selection cannot move it to a subclass.
				return m
			}
			if m.Flags&(AccStatic|AccPrivate) != 0 {
				continue
			}
			// Selection uses the JVMS 5.4.5 can-override relation between the
			// candidate and the resolved declaration. In particular a
			// package-private class declaration is overridden only by
			// same-runtime-package candidates or through a transitive chain
			// of overriding declarations; a public interface declaration is
			// also satisfied by an inherited implementation from a superclass
			// that does not itself implement the interface.
			if !r.selectionOverrides(m, declaration, map[*Method]bool{}) {
				continue
			}
			if m.Flags&AccAbstract != 0 {
				return nil
			}
			return m
		}
	}
	candidates := r.collectDefaultMethods(receiver, key)
	if len(candidates) == 0 {
		return nil
	}
	// Keep the maximally specific declarations: remove any candidate whose
	// owner is a superinterface of another candidate's owner.
	maximal := candidates[:0]
	for _, cand := range candidates {
		overridden := false
		for _, other := range candidates {
			if other == cand {
				continue
			}
			if other.Owner != cand.Owner && r.interfaceExtends(other.Owner, cand.Owner) {
				overridden = true
				break
			}
		}
		if !overridden {
			maximal = append(maximal, cand)
		}
	}
	if len(maximal) == 0 {
		return nil
	}
	// With multiple unrelated maximally-specific defaults Java permits any;
	// the profile selects deterministically so both engines agree.
	best := maximal[0]
	for _, cand := range maximal[1:] {
		if cand.Owner.Name < best.Owner.Name {
			best = cand
		}
	}
	return best
}

// selectionOverrides implements the JVMS 5.4.5 can-override relation used by
// method selection, where the hierarchy relation is implicit in the
// receiver-scoped search. methodOverrides below additionally requires the
// declaring types to be related; preparation uses that stronger form to decide
// which inherited declarations a candidate actually overrides.
func (r *Runtime) selectionOverrides(candidate, declaration *Method, seen map[*Method]bool) bool {
	if candidate == declaration {
		return true
	}
	if candidate == nil || declaration == nil || candidate.Owner == nil || declaration.Owner == nil {
		return false
	}
	if candidate.Flags&(AccStatic|AccPrivate) != 0 {
		return false
	}
	if candidate.Name != declaration.Name || candidate.Descriptor != declaration.Descriptor {
		return false
	}
	if declaration.Flags&AccPrivate != 0 {
		return false
	}
	// A public or protected declaration (including every non-private
	// interface method) is satisfied by any same-signature instance method
	// reachable through the receiver's superclass chain.
	if declaration.Flags&(AccPublic|AccProtected) != 0 {
		return true
	}
	if r.samePackage(candidate.Owner, declaration.Owner) {
		return true
	}
	if seen[candidate] {
		return false
	}
	seen[candidate] = true
	key := methodKey(candidate.Name, candidate.Descriptor)
	for c := candidate.Owner.Super; c != nil && c != declaration.Owner; c = c.Super {
		m := c.MethodIndex[key]
		if m == nil || m.Flags&(AccStatic|AccPrivate) != 0 {
			continue
		}
		if r.selectionOverrides(m, declaration, seen) && r.selectionOverrides(candidate, m, seen) {
			return true
		}
	}
	return false
}

// methodOverrides implements the JVMS 5.4.5 instance-method override relation,
// including the package-private runtime-package rule and transitive overriding
// through an intermediate declaration.
func (r *Runtime) methodOverrides(candidate, declaration *Method, seen map[*Method]bool) bool {
	if candidate == declaration {
		return true
	}
	if candidate == nil || declaration == nil || candidate.Owner == nil || declaration.Owner == nil {
		return false
	}
	if candidate.Flags&(AccStatic|AccPrivate) != 0 {
		return false
	}
	if candidate.Name != declaration.Name || candidate.Descriptor != declaration.Descriptor {
		return false
	}
	if declaration.Owner.IsInterface() {
		if candidate.Owner.IsInterface() {
			return r.interfaceExtends(candidate.Owner, declaration.Owner)
		}
		return r.implementsInterface(candidate.Owner, declaration.Owner)
	}
	if !r.isSubclass(candidate.Owner, declaration.Owner) {
		return false
	}
	if declaration.Flags&(AccPublic|AccProtected) != 0 {
		return true
	}
	if r.samePackage(candidate.Owner, declaration.Owner) {
		return true
	}
	if seen[candidate] {
		return false
	}
	seen[candidate] = true
	key := methodKey(candidate.Name, candidate.Descriptor)
	for c := candidate.Owner.Super; c != nil && c != declaration.Owner; c = c.Super {
		m := c.MethodIndex[key]
		if m == nil || m.Flags&(AccStatic|AccPrivate) != 0 {
			continue
		}
		if r.methodOverrides(m, declaration, seen) && r.methodOverrides(candidate, m, seen) {
			return true
		}
	}
	return false
}

func (r *Runtime) collectDefaultMethods(receiver *Class, key string) []*Method {
	var out []*Method
	seen := map[*Class]bool{}
	var visit func(c *Class)
	visit = func(c *Class) {
		if c == nil || seen[c] {
			return
		}
		seen[c] = true
		if c.IsInterface() {
			if m := c.MethodIndex[key]; m != nil && m.Flags&(AccStatic|AccAbstract|AccPrivate) == 0 {
				out = append(out, m)
			}
		}
		for _, in := range c.Interfaces {
			visit(in)
		}
		if !c.IsInterface() {
			visit(c.Super)
		}
	}
	visit(receiver)
	return out
}

func (r *Runtime) interfaceExtends(from, to *Class) bool {
	seen := map[*Class]bool{}
	var visit func(c *Class) bool
	visit = func(c *Class) bool {
		if c == nil || seen[c] {
			return false
		}
		seen[c] = true
		if c == to {
			return true
		}
		for _, in := range c.Interfaces {
			if visit(in) {
				return true
			}
		}
		return false
	}
	for _, in := range from.Interfaces {
		if visit(in) {
			return true
		}
	}
	return false
}

// --------------------------------------------------------------- invocations

func refOrNull(t *Thread, v Value) (*Object, *Object) {
	if v.Kind != KindRef {
		return nil, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	if v.Ref == nil {
		return nil, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return v.Ref, nil
}

func (t *Thread) invocationEnv(caller *Method, pc int32) {
	_ = caller
	t.SetPC(pc)
}

// InvokeStatic implements invokestatic.
func (t *Thread) InvokeStatic(caller *Method, pc int32, cp int, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	m, thrown := t.resolveMethodSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if m.Flags&AccStatic == 0 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError"))
	}
	if thrown := t.EnsureInitialized(m.Owner); thrown != nil {
		return ThrowOutcome(thrown)
	}
	return t.Call(m, nil, args)
}

// InvokeSpecial implements invokespecial.
func (t *Thread) InvokeSpecial(caller *Method, pc int32, cp int, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	m, thrown := t.resolveMethodSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if m.Flags&AccStatic != 0 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError"))
	}
	if !t.RT.protectedReceiverOK(caller, m.Owner, recv.Class) {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessError"))
	}
	return t.Call(m, recv, args)
}

// InvokeVirtual implements invokevirtual.
func (t *Thread) InvokeVirtual(caller *Method, pc int32, cp int, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	decl, thrown := t.resolveMethodSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if decl.Flags&AccStatic != 0 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError"))
	}
	if !t.RT.protectedReceiverOK(caller, decl.Owner, recv.Class) {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessError"))
	}
	target, hit, thrown := t.selectGuarded(caller, pc, OpInvokevirtual, recv.Class, decl)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if tr := t.trace(); tr != nil {
		tr.dispatch(t.ID, "invokevirtual", caller, pc, decl, target, recv.Class, hit)
	}
	return t.Call(target, recv, args)
}

// InvokeInterface implements invokeinterface and emits the direct dispatch
// evidence before entering the selected body.
func (t *Thread) InvokeInterface(caller *Method, pc int32, cp int, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	decl, thrown := t.resolveMethodSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if !decl.Owner.IsInterface() {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError"))
	}
	if !t.RT.protectedReceiverOK(caller, decl.Owner, recv.Class) {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessError"))
	}
	target, hit, thrown := t.selectGuarded(caller, pc, OpInvokeinterface, recv.Class, decl)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	// JVMS 6.5 invokeinterface: a selected method that is neither public nor
	// private makes the instruction throw IllegalAccessError. A package-private
	// superclass implementation therefore does not satisfy a public interface
	// method even though selection can find it.
	if target.Flags&(AccPublic|AccPrivate) == 0 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessError"))
	}
	if tr := t.trace(); tr != nil {
		tr.dispatch(t.ID, "invokeinterface", caller, pc, decl, target, recv.Class, hit)
	}
	return t.Call(target, recv, args)
}

// selectGuarded applies the Gate 3 compiled-dispatch cache. Symbolic
// resolution, receiver and access checks have already run; a hit requires the
// identical resolved declaration and exact receiver class. Admission of a new
// subclass therefore cannot change an existing receiver's successful
// selection, and a new receiver takes the full selection path.
func (t *Thread) selectGuarded(caller *Method, pc int32, op uint8, receiver *Class, decl *Method) (*Method, bool, *Object) {
	key := dispatchKey{Caller: caller, PC: pc, Op: op}
	r := t.RT
	r.metaMu.Lock()
	defer r.metaMu.Unlock()
	if entry, ok := r.dispatch[key]; ok && entry.Declaration == decl && entry.Receiver == receiver {
		return entry.Target, true, nil
	}
	target := r.selectVirtual(receiver, decl)
	if target == nil {
		return nil, false, r.newThrowableLocked(t, "java/lang/AbstractMethodError")
	}
	if target != decl && target.Owner.Loader != decl.Owner.Loader &&
		r.imposeMethodConstraintsLocked(target.Owner.Loader, decl.Owner.Loader, target.Type) {
		return nil, false, r.newThrowableLocked(t, "java/lang/LinkageError")
	}
	r.dispatch[key] = &dispatchEntry{Declaration: decl, Receiver: receiver, Target: target}
	return target, false, nil
}

// ------------------------------------------------------------------- fields

func coerceFieldValue(ft Type, v Value) Value {
	if ft.Kind != 'I' {
		return v
	}
	switch ft.Descriptor {
	case "Z":
		return IntValue(AsInt(v) & 1)
	case "B":
		return IntValue(int32(int8(AsInt(v))))
	case "C":
		return IntValue(int32(uint16(AsInt(v))))
	case "S":
		return IntValue(int32(int16(AsInt(v))))
	}
	return v
}

// GetField implements getfield.
func (t *Thread) GetField(caller *Method, cp int, ref Value, owner, name, descriptor string) (Value, *Object) {
	obj, thrown := refOrNull(t, ref)
	if thrown != nil {
		return Value{}, thrown
	}
	f, thrown := t.resolveFieldSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return Value{}, thrown
	}
	if f.Static {
		return Value{}, t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if !t.RT.protectedReceiverOK(caller, f.Owner, obj.Class) {
		return Value{}, t.RT.NewThrowable(t, "java/lang/IllegalAccessError")
	}
	if f.Offset < 0 || f.Offset >= len(obj.Fields) {
		return Value{}, t.RT.NewThrowable(t, "java/lang/NoSuchFieldError")
	}
	return t.RT.fieldRead(obj, f.Offset), nil
}

// PutField implements putfield.
func (t *Thread) PutField(caller *Method, cp int, ref Value, owner, name, descriptor string, v Value) *Object {
	obj, thrown := refOrNull(t, ref)
	if thrown != nil {
		return thrown
	}
	f, thrown := t.resolveFieldSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return thrown
	}
	if f.Static {
		return t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if !t.RT.protectedReceiverOK(caller, f.Owner, obj.Class) {
		return t.RT.NewThrowable(t, "java/lang/IllegalAccessError")
	}
	if f.Offset < 0 || f.Offset >= len(obj.Fields) {
		return t.RT.NewThrowable(t, "java/lang/NoSuchFieldError")
	}
	t.RT.fieldWrite(obj, f.Offset, coerceFieldValue(f.Type, v))
	return nil
}

// GetStatic implements getstatic, including class initialization.
func (t *Thread) GetStatic(caller *Method, cp int, owner, name, descriptor string) (Value, *Object) {
	f, thrown := t.resolveFieldSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return Value{}, thrown
	}
	if !f.Static {
		return Value{}, t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if thrown := t.EnsureInitialized(f.Owner); thrown != nil {
		return Value{}, thrown
	}
	return t.RT.staticRead(f.Owner, f.Index), nil
}

// PutStatic implements putstatic, including class initialization.
func (t *Thread) PutStatic(caller *Method, cp int, owner, name, descriptor string, v Value) *Object {
	f, thrown := t.resolveFieldSymbol(caller, cp, owner, name, descriptor)
	if thrown != nil {
		return thrown
	}
	if !f.Static {
		return t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if thrown := t.EnsureInitialized(f.Owner); thrown != nil {
		return thrown
	}
	t.RT.staticWrite(f.Owner, f.Index, coerceFieldValue(f.Type, v))
	return nil
}

// StaticValue resolves and returns one static field value (used by tests and
// host bindings).
func (r *Runtime) StaticValue(owner, name, descriptor string) (Value, error) {
	c, err := r.classForName(owner)
	if err != nil {
		return Value{}, fmt.Errorf("catty: cannot resolve class %s", owner)
	}
	f := r.findField(c, name, descriptor)
	if f == nil {
		return Value{}, fmt.Errorf("catty: cannot resolve static field %s.%s%s", owner, name, descriptor)
	}
	if !f.Static {
		return Value{}, fmt.Errorf("catty: field %s.%s%s is not static", owner, name, descriptor)
	}
	return r.staticRead(f.Owner, f.Index), nil
}

// ---------------------------------------------------------------- allocation

// NewObject implements new.
func (t *Thread) NewObject(caller *Method, cp int, className string) (*Object, *Object) {
	c, err := t.resolveClassSymbol(caller, cp, className)
	if err != nil {
		return nil, err
	}
	if !t.RT.classAccessible(callerOwner(caller), c) {
		return nil, t.RT.NewThrowable(t, "java/lang/IllegalAccessError")
	}
	if c.IsInterface() || c.IsArray || c.IsPrimitive || c.Flags&AccAbstract != 0 {
		return nil, t.RT.NewThrowable(t, "java/lang/InstantiationError")
	}
	if thrown := t.EnsureInitialized(c); thrown != nil {
		return nil, thrown
	}
	obj := t.RT.NewInstance(c)
	if t.RT.IsAssignable(c, t.RT.mustClass("java/lang/Throwable")) {
		obj.Payload = &ThrowableData{Frames: t.FramesSnapshot()}
	}
	return obj, nil
}

func callerOwner(caller *Method) *Class {
	if caller == nil {
		return nil
	}
	return caller.Owner
}

// loaderOf returns the loader that owns the referring method, falling back to
// the application loader for synthetic host calls.
func loaderOf(caller *Method) *Loader {
	if caller != nil && caller.Owner != nil && caller.Owner.Loader != nil {
		return caller.Owner.Loader
	}
	return nil
}

// NewPrimitiveArray implements newarray.
func (t *Thread) NewPrimitiveArray(caller *Method, cp int, atype int32, count int32) (*Object, *Object) {
	_ = caller
	if count < 0 {
		return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
	}
	desc, ok := NewarrayDescriptor(atype)
	if !ok {
		return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
	}
	ac, thrown := t.resolveClassSymbol(caller, cp, desc)
	if thrown != nil {
		return nil, thrown
	}
	return t.RT.NewArrayOfClass(ac, int(count)), nil
}

// NewReferenceArray implements anewarray.
func (t *Thread) NewReferenceArray(caller *Method, cp int, componentName string, count int32) (*Object, *Object) {
	_ = caller
	if count < 0 {
		return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
	}
	component, thrown := t.resolveClassSymbol(caller, cp, componentName)
	if thrown != nil {
		return nil, thrown
	}
	return t.RT.NewArrayOfClass(t.RT.arrayClassForComponent(component), int(count)), nil
}

// MultiANewArray implements multianewarray. Dimensions are allocated in order;
// unspecified trailing dimensions stay null.
func (t *Thread) MultiANewArray(caller *Method, cp int, descriptor string, dims int, counts []int32) (*Object, *Object) {
	_ = caller
	ac, thrown := t.resolveClassSymbol(caller, cp, descriptor)
	if thrown != nil {
		return nil, thrown
	}
	if dims < 1 || dims > len(counts) {
		return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
	}
	var allocate func(class *Class, depth int) (*Object, *Object)
	allocate = func(class *Class, depth int) (*Object, *Object) {
		count := counts[depth]
		if count < 0 {
			return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
		}
		obj := t.RT.NewArrayOfClass(class, int(count))
		if depth+1 < dims && class.Component != nil && class.Component.IsArray {
			for i := 0; i < int(count); i++ {
				child, thrown := allocate(class.Component, depth+1)
				if thrown != nil {
					return nil, thrown
				}
				obj.Array.Elems[i] = RefValue(child)
			}
		}
		return obj, nil
	}
	return allocate(ac, 0)
}

// --------------------------------------------------------------- conversions

// CheckCast implements checkcast.
func (t *Thread) CheckCast(caller *Method, cp int, className string, v Value) (Value, *Object) {
	if v.Kind != KindRef {
		return Value{}, t.RT.NewThrowable(t, "java/lang/ClassCastException")
	}
	if v.Ref == nil {
		return v, nil
	}
	c, thrown := t.resolveClassSymbol(caller, cp, className)
	if thrown != nil {
		return Value{}, thrown
	}
	if !t.RT.IsAssignable(v.Ref.Class, c) {
		return Value{}, t.RT.NewThrowable(t, "java/lang/ClassCastException")
	}
	return v, nil
}

// InstanceOf implements instanceof.
func (t *Thread) InstanceOf(caller *Method, cp int, className string, v Value) (int32, *Object) {
	if v.Kind != KindRef {
		return 0, t.RT.NewThrowable(t, "java/lang/ClassCastException")
	}
	if v.Ref == nil {
		return 0, nil
	}
	c, thrown := t.resolveClassSymbol(caller, cp, className)
	if thrown != nil {
		return 0, thrown
	}
	if t.RT.IsAssignable(v.Ref.Class, c) {
		return 1, nil
	}
	return 0, nil
}

// -------------------------------------------------------------------- arrays

func arrayObject(t *Thread, array Value) (*Object, *Object) {
	obj, thrown := refOrNull(t, array)
	if thrown != nil {
		return nil, thrown
	}
	if obj.Array == nil {
		return nil, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return obj, nil
}

func arrayIndex(t *Thread, a *ArrayData, index int32) *Object {
	if index < 0 || int(index) >= a.Len() {
		return t.RT.NewThrowable(t, "java/lang/ArrayIndexOutOfBoundsException")
	}
	return nil
}

// ArrayLoad implements typed array loads. kind is the element descriptor the
// instruction expects ('L' for reference arrays).
func (t *Thread) ArrayLoad(array Value, index int32, kind byte) (Value, *Object) {
	obj, thrown := arrayObject(t, array)
	if thrown != nil {
		return Value{}, thrown
	}
	if thrown := arrayIndex(t, obj.Array, index); thrown != nil {
		return Value{}, thrown
	}
	switch kind {
	case 'L', 'B', 'C', 'S', 'I', 'J', 'F', 'D':
		return t.RT.arrayRead(obj.Array, int(index)), nil
	}
	panic(fmt.Sprintf("catty: unknown array load kind %q", kind))
}

// ArrayStore implements typed array stores. Java checks null, then bounds, then
// reference assignability.
func (t *Thread) ArrayStore(array Value, index int32, value Value, kind byte) *Object {
	obj, thrown := arrayObject(t, array)
	if thrown != nil {
		return thrown
	}
	if thrown := arrayIndex(t, obj.Array, index); thrown != nil {
		return thrown
	}
	switch kind {
	case 'L':
		if value.Kind != KindRef {
			return t.RT.NewThrowable(t, "java/lang/ArrayStoreException")
		}
		if value.Ref != nil && !t.RT.IsAssignable(value.Ref.Class, obj.Array.Component) {
			return t.RT.NewThrowable(t, "java/lang/ArrayStoreException")
		}
		t.RT.arrayWrite(obj.Array, int(index), value)
	case 'B', 'C', 'S', 'I', 'J', 'F', 'D':
		t.RT.arrayWrite(obj.Array, int(index), value)
	default:
		panic(fmt.Sprintf("catty: unknown array store kind %q", kind))
	}
	return nil
}

// ArrayLength implements arraylength.
func (t *Thread) ArrayLength(array Value) (int32, *Object) {
	obj, thrown := arrayObject(t, array)
	if thrown != nil {
		return 0, thrown
	}
	return int32(arrayLength(obj.Array)), nil
}

// ---------------------------------------------------------------- constants

// ClassConstant implements ldc of a Class constant.
func (t *Thread) ClassConstant(caller *Method, cp int, name string) (Value, *Object) {
	c, thrown := t.resolveClassSymbol(caller, cp, name)
	if thrown != nil {
		return Value{}, thrown
	}
	return RefValue(c.Mirror), nil
}

// StringConstant implements ldc of a string constant with JVM interning.
func (t *Thread) StringConstant(units []uint16) Value {
	return RefValue(t.RT.NewStringUnits(units, true))
}

// Throwable allocates a fresh throwable of the given bootstrap class and
// captures the current Java frames.
func (t *Thread) Throwable(className string) *Object {
	return t.RT.NewThrowable(t, className)
}

// HandleThrow implements the shared Java exception transfer. It returns the
// matching handler's original bytecode offset, or false when the throwable
// escapes this method.
func (t *Thread) HandleThrow(m *Method, throwPC int32, thrown *Object) (int32, bool) {
	if thrown == nil {
		return 0, false
	}
	for _, h := range m.Handlers {
		if throwPC < h.StartPC || throwPC >= h.EndPC {
			continue
		}
		if h.CatchType != "" {
			catchClass, err := t.RT.classForLoaderName(loaderOf(m), h.CatchType)
			if err != nil {
				continue
			}
			if !t.RT.IsAssignable(thrown.Class, catchClass) {
				continue
			}
		}
		if tr := t.trace(); tr != nil {
			tr.catchEvent(t.ID, m, throwPC, h.HandlerPC, thrown)
		}
		return h.HandlerPC, true
	}
	return 0, false
}
