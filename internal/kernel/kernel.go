package kernel

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"catty/internal/classfile"
	"catty/internal/verify"
)

// NativeFunc is a synthesized method implementation. recv is nil for
// static methods; args excludes the receiver and follows ParseMethodDesc
// argument order.
type NativeFunc func(ctx *CallContext, recv Value, args []Value) (Value, error)

// OwnerKey is the opaque thread identity used for monitor ownership.
// Implemented by the execution layer's thread abstraction.
type OwnerKey interface {
	OwnerKey() uint64
}

// ErrStackOverflow signals frame-budget exhaustion to FrameMeter users.
var ErrStackOverflow = errors.New("stack overflow")

// FrameMeter is an optional OwnerKey extension: threads that meter frame
// depth on themselves (plain field, no registry locking). The emitted
// path (genrt.invokeChecked) prefers it over ThreadRegistry.FrameEnter so
// the hot dispatch chain stays lock-free; interpreted frames share the
// same budget, matching the JVM single-stack model.
type FrameMeter interface {
	EnterFrame() error // ErrStackOverflow when budget exhausted
	ExitFrame()
}

// CallContext gives natives access to runtime services.
type CallContext struct {
	K     *Kernel
	Owner OwnerKey // executing thread identity (nil in thread-less paths)
}

func (ctx *CallContext) ownerKey() uint64 {
	if ctx.Owner != nil {
		return ctx.Owner.OwnerKey()
	}
	return 0
}

// OwnerKeyValue exposes the calling thread's monitor key to natives/tests.
func (ctx *CallContext) OwnerKeyValue() uint64 { return ctx.ownerKey() }

// Stringify renders a value the way PrintStream.println(Object) does,
// dispatching toString() with the calling thread's identity.
func (ctx *CallContext) Stringify(v Value) string {
	if v == nil {
		return "null"
	}
	switch o := v.(type) {
	case *JString:
		return o.String()
	case *Instance:
		if m, err := ctx.K.ResolveMethod(o.Class, "toString", "()Ljava/lang/String;"); err == nil {
			if r, err2 := ctx.Invoke(m, o, nil); err2 == nil {
				if s, ok := r.(*JString); ok {
					return s.String()
				}
			}
		}
		return dotted(o.Class.Name) + "@" + fmt.Sprintf("%x", o.IdentityHash())
	case *ArrayObj:
		return dotted(o.Class.Name) + "@" + fmt.Sprintf("%x", o.IdentityHash())
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Invoke calls a resolved method through the kernel dispatcher, carrying
// the calling thread's identity so nested monitor/sleep ops attribute
// correctly.
func (ctx *CallContext) Invoke(m *Method, recv Value, args []Value) (Value, error) {
	return ctx.K.InvokeAs(ctx.Owner, m, recv, args)
}

// NewStringGo builds a non-interned String from Go text.
func (ctx *CallContext) NewStringGo(s string) *JString { return ctx.K.MakeJStringFromGo(s) }

// Throw constructs a Java throwable wrapped for VM unwinding.
func (ctx *CallContext) Throw(className, msg string) error {
	return ctx.K.Throw(className, msg)
}

// Thrown wraps a Java throwable propagating through Go stacks. The VM
// distinguishes it from engine errors via errors.As.
type Thrown struct {
	Obj *Instance
}

func (t *Thrown) Error() string {
	msg := ""
	if s, ok := t.Obj.fieldByName("detailMessage").(*JString); ok && s != nil {
		msg = ": " + s.String()
	}
	return "uncaught " + dotted(t.Obj.Class.Name) + msg
}

func dotted(internal string) string {
	return strings.ReplaceAll(internal, "/", ".")
}

// Invoker is implemented by the VM so the kernel can execute interpreted
// code (natives calling back into Java, <clinit>, …).
type Invoker interface {
	// InvokeInterpreted executes an interpreted method. recv is nil for
	// static methods; args excludes the receiver.
	InvokeInterpreted(m *Method, recv Value, args []Value) (Value, error)
}

// Options configures a Kernel.
type Options struct {
	Stdout io.Writer
	Stderr io.Writer

	// SkipVerify disables bytecode verification (JVMS §4.10). Default is
	// verified-on-load once the verifier is wired; trusted-input tooling
	// may turn it off.
	SkipVerify bool

	// MaxFrames bounds interpreted-call depth for StackOverflowError
	// detection (0 selects the default 4096).
	MaxFrames int

	// Resolver optionally backs runtime class resolution (VM
	// resolveClassIdx / ResolveClass misses). Set to a ClassPathLoader.Load
	// bound method by embedders that support dynamic loading.
	Resolver func(name string) (*Class, error)
}

// Kernel is the runtime registry and object-model entry point.
type Kernel struct {
	opts    Options
	mu      sync.Mutex
	classes map[string]*Class
	strPool map[string]*JString
	intBox  [256]*Instance // Integer.valueOf cache [-128,127]
	intCls  atomic.Pointer[Class] // lazily cached bootstrap Integer class

	// Threads tracks java.lang.Thread identities and blocking operations.
	Threads *ThreadRegistry

	nextKey atomic.Uint64

	// SpawnJavaThread is installed by the VM: start() delegates goroutine
	// creation upward (kernel must not import the VM). The hook runs the
	// thread's run() on a fresh execution context and terminates the record.
	SpawnJavaThread func(j *JThread)

	// UncaughtHandler formats an uncaught throwable from a Java thread.
	// Installed by the VM alongside SpawnJavaThread.
	UncaughtHandler func(j *JThread, thrown *Thrown)

	// classLoadHook runs after a classfile-backed class is registered
	// (lazy loads included). It is how the emitter layer attaches
	// generated bodies to classes that resolve AFTER gen.Install ran —
	// nested classes and library dependencies load lazily at first
	// `new`, which previously left them interpreted forever. Guarded by
	// mu; called outside it.
	classLoadHook func(*Class)

	// invoker is the VM execution bridge. Accessed via setInvoker/
	// invokerFallback because spawned threads race to install it.
	invoker Invoker
}

// SetClassLoadHook installs the post-registration callback. One hook per
// kernel; the last installation wins (single emitter layer per process).
func (k *Kernel) SetClassLoadHook(h func(*Class)) {
	k.mu.Lock()
	k.classLoadHook = h
	k.mu.Unlock()
}

// ctxPool backs InvokeAs's native-dispatch branch (CallContext reuse).
var ctxPool sync.Pool

// InstallInvoker registers the VM execution bridge. Idempotent: the first
// installation wins (all threads of one kernel share one engine bridge).
func (k *Kernel) InstallInvoker(i Invoker) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.invoker == nil {
		k.invoker = i
	}
}

// invokerFallback returns the installed bridge (may be nil).
func (k *Kernel) invokerFallback() Invoker {
	k.mu.Lock()
	defer k.mu.Unlock()
	return k.invoker
}

// New builds a Kernel with the full bootstrap (synthesized java.*) surface.
func New(opts Options) *Kernel {
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}
	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}
	k := &Kernel{
		opts:    opts,
		classes: make(map[string]*Class),
		strPool: make(map[string]*JString),
	}
	k.Threads = NewThreadRegistry(k.MaxFrames())
	bootstrap(k)
	return k
}

// Stderr returns the diagnostic writer (uncaught exceptions etc).
func (k *Kernel) Stderr() io.Writer { return k.opts.Stderr }

// MintKey allocates a unique monitor-owner/thread identity.
func (k *Kernel) MintKey() uint64 { return k.nextKey.Add(1) }

// MaxFrames returns the interpreted-frame budget for StackOverflowError
// detection (default 4096).
func (k *Kernel) MaxFrames() int {
	if k.opts.MaxFrames > 0 {
		return k.opts.MaxFrames
	}
	return 4096
}

// Stdout returns the writer backing System.out.
func (k *Kernel) Stdout() io.Writer { return k.opts.Stdout }

// SetResolver installs a runtime class-resolution fallback (classpath
// loader). Call before executing Java code.
func (k *Kernel) SetResolver(r func(name string) (*Class, error)) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.opts.Resolver = r
}

// Invoke dispatches any resolved method without a thread identity.
// Prefer InvokeAs from code running on behalf of a Java thread.
func (k *Kernel) Invoke(m *Method, recv Value, args []Value) (Value, error) {
	return k.InvokeAs(nil, m, recv, args)
}

// InvokeAs dispatches a method carrying the calling thread's identity so
// monitors, class-init re-entrancy and interrupt state attribute correctly.
//
// Interpreted methods run on the OWNER's execution context when it provides
// one (the VM's Thread does); k.Invoker is only the fallback for
// owner-less paths (tests, kernel-internal probes).
func (k *Kernel) InvokeAs(owner OwnerKey, m *Method, recv Value, args []Value) (Value, error) {
	// Java call-stack maintenance for stack backfill (DEBT-0019). Every
	// engine funnels through here, so one push/pop pair covers interpreted
	// and emitted frames alike.
	var ft FrameTracker
	if f, ok := owner.(FrameTracker); ok {
		ft = f
		ft.PushJavaFrame(JavaFrame{Class: m.Holder.Name, Method: m.Name})
		defer ft.PopJavaFrame()
	}
	// Emitter-generated bodies take the allocation-free fast path
	// (P-0009 T1): the ABI returns (Value, *Thrown) directly and needs
	// no CallContext.
	if m.EmitBody != nil {
		v, exc := m.EmitBody(owner, recv, args)
		if exc != nil {
			return nil, exc
		}
		return v, nil
	}
	if m.Native != nil {
		// CallContext pooling: natives allocate one per call otherwise
		// (mapops alloc profile showed this as the single largest
		// engine-side source). Pool handles reentrant ctx.Invoke nesting.
		c, _ := ctxPool.Get().(*CallContext)
		if c == nil {
			c = &CallContext{}
		}
		c.K, c.Owner = k, owner
		v, err := m.Native(c, recv, args)
		c.K, c.Owner = nil, nil
		ctxPool.Put(c)
		return v, err
	}
	if io, ok := owner.(Invoker); ok {
		return io.InvokeInterpreted(m, recv, args)
	}
	fb := k.invokerFallback()
	if fb == nil {
		return nil, fmt.Errorf("kernel: no Invoker installed for interpreted method %s.%s", m.Holder.Name, m.Name)
	}
	return fb.InvokeInterpreted(m, recv, args)
}

// ClassByName looks a class up by internal name (array pseudo-classes
// included, created on demand).
func (k *Kernel) ClassByName(name string) (*Class, bool) {
	if strings.HasPrefix(name, "[") {
		return k.ArrayClassOf(name[1:]), true
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	c, ok := k.classes[name]
	return c, ok
}

// lookupClass is the unlocked registry read. Callers must either hold
// k.mu or be inside single-threaded bootstrap/native contexts (M0 scope;
// the M1 loader work replaces this with a proper class loader hierarchy).
func (k *Kernel) lookupClass(name string) *Class {
	return k.classes[name]
}

// ResolveClass is ClassByName with an error for unknown classes. When a
// Resolver is configured, misses fall back to it (and thereby register).
func (k *Kernel) ResolveClass(name string) (*Class, error) {
	c, ok := k.ClassByName(name)
	if ok {
		return c, nil
	}
	if k.opts.Resolver != nil && !strings.HasPrefix(name, "[") {
		return k.opts.Resolver(name)
	}
	return nil, fmt.Errorf("java.lang.ClassNotFoundException: %s", dotted(name))
}

// ArrayClassOf returns the array pseudo-class for a component descriptor,
// creating and registering it on first use.
func (k *Kernel) ArrayClassOf(compDesc string) *Class {
	name := "[" + compDesc
	k.mu.Lock()
	defer k.mu.Unlock()
	if c, ok := k.classes[name]; ok {
		return c
	}
	objCls := k.classes["java/lang/Object"]
	c := &Class{
		Name:     name,
		Super:    objCls,
		Flags:    classfile.AccFinal | classfile.AccPublic,
		IsArray:  true,
		CompDesc: compDesc,
	}
	c.setState(StateInitialized)
	k.classes[name] = c
	return c
}

// ---- Class definition / loading -----------------------------------------

// ClassDef declares a synthesized class.
type ClassDef struct {
	Name       string
	Super      string // "" → java/lang/Object (ignored for Object itself)
	Ifaces     []string
	Flags      uint16
	Fields     []FieldDef
	Methods    []MethodDef
	StaticInit func(k *Kernel, c *Class) error
}

// FieldDef declares a field of a synthesized class.
type FieldDef struct {
	Name  string
	Desc  string
	Flags uint16
}

// MethodDef declares a method of a synthesized class.
type MethodDef struct {
	Name   string
	Desc   string
	Flags  uint16
	Native NativeFunc
}

// DefineClass registers a synthesized class. Synthesized classes are born
// initialized: their static hook runs immediately and there is no <clinit>.
func (k *Kernel) DefineClass(def *ClassDef) (*Class, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	if _, exists := k.classes[def.Name]; exists {
		return nil, fmt.Errorf("class %s already defined", def.Name)
	}

	c := &Class{
		Name: def.Name, Flags: def.Flags, def: def,
		methodsByKey: make(map[string]*Method),
		fieldsByKey:  make(map[string]*Field),
	}
	c.setState(StateInitializing)
	k.classes[def.Name] = c

	superName := ""
	if def.Name != "java/lang/Object" {
		superName = def.Super
		if superName == "" {
			superName = "java/lang/Object"
		}
	}
	if superName != "" {
		sc, ok := k.classes[superName]
		if !ok {
			delete(k.classes, def.Name)
			return nil, fmt.Errorf("define %s: superclass %s not defined yet", def.Name, superName)
		}
		c.Super = sc
	}
	for _, in := range def.Ifaces {
		ic, ok := k.classes[in]
		if !ok {
			delete(k.classes, def.Name)
			return nil, fmt.Errorf("define %s: interface %s not defined yet", def.Name, in)
		}
		c.Ifaces = append(c.Ifaces, ic)
	}

	// Layout: root-first so supers occupy low slots; flat member view gets
	// every ancestor field.
	var chain []*Class
	for x := c; x != nil; x = x.Super {
		chain = append([]*Class{x}, chain...)
	}
	slot := 0
	for _, x := range chain {
		for _, fd := range x.def.Fields {
			f := &Field{Holder: x, Name: fd.Name, Desc: fd.Desc, Flags: fd.Flags}
			if fd.Flags&classfile.AccStatic != 0 {
				f.Static = true
				f.StaticSlot = len(x.Statics)
				x.Statics = append(x.Statics, nil)
			} else {
				f.Slot = slot
				slot += SlotCount(fd.Desc) // category-2 (J/D) occupies two slots
				if x == c {
					c.OwnFields = append(c.OwnFields, f)
				}
			}
			c.flatField(f)
		}
	}
	c.layoutSize = slot

	for _, md := range def.Methods {
		m := &Method{Holder: c, Name: md.Name, Desc: md.Desc, Flags: md.Flags, Native: md.Native}
		c.Methods = append(c.Methods, m)
		c.methodsByKey[memberKey(m.Name, m.Desc)] = m
	}

	if def.StaticInit != nil {
		if err := def.StaticInit(k, c); err != nil {
			c.setState(StateErroneous)
			return nil, fmt.Errorf("static init of %s: %w", def.Name, err)
		}
	}
	c.setState(StateInitialized)
	return c, nil
}

// LoadClassBytes parses and registers a class file image. Superclass and
// interfaces must already be known (M0-style explicit loading).
func (k *Kernel) LoadClassBytes(data []byte) (*Class, error) {
	return k.LoadClassBytesWith(data, nil)
}

// LoadClassBytesWith parses and registers a class file, resolving
// superclasses/interfaces through dep when absent from the registry.
//
// M1 note: registration happens up-front (placeholder pattern) so recursive
// dependencies observe a stable namespace; concurrent duplicate loads of the
// same name are outside M1 scope and serialized by embedders.
func (k *Kernel) LoadClassBytesWith(data []byte, dep func(name string) (*Class, error)) (*Class, error) {
	cf, err := classfile.Parse(data)
	if err != nil {
		return nil, err
	}

	k.mu.Lock()
	if _, exists := k.classes[cf.ThisClass]; exists {
		k.mu.Unlock()
		return nil, fmt.Errorf("class %s already loaded", cf.ThisClass)
	}
	c := &Class{Name: cf.ThisClass, Flags: cf.AccessFlags, CF: cf,
		methodsByKey: make(map[string]*Method), fieldsByKey: make(map[string]*Field)}
	c.setState(StateDefined)
	k.classes[c.Name] = c
	k.mu.Unlock()

	fail := func(err error) (*Class, error) {
		k.mu.Lock()
		delete(k.classes, c.Name)
		k.mu.Unlock()
		return nil, err
	}

	resolveDep := func(name string) (*Class, error) {
		k.mu.Lock()
		existing, ok := k.classes[name]
		k.mu.Unlock()
		if ok {
			return existing, nil
		}
		if dep == nil {
			return nil, fmt.Errorf("load %s: %s not loaded", c.Name, name)
		}
		return dep(name)
	}

	if cf.SuperClass != "" {
		sc, err := resolveDep(cf.SuperClass)
		if err != nil {
			return fail(err)
		}
		c.Super = sc
	}
	for _, in := range cf.Interfaces {
		ic, err := resolveDep(in)
		if err != nil {
			return fail(err)
		}
		c.Ifaces = append(c.Ifaces, ic)
	}

	// Layout across the hierarchy.
	var chain []*Class
	for x := c; x != nil; x = x.Super {
		chain = append([]*Class{x}, chain...)
	}
	slot := 0
	for _, x := range chain {
		if x.CF == nil {
			continue // synthesized ancestors contribute nothing here
		}
		for i := range x.CF.Fields {
			fd := &x.CF.Fields[i]
			f := &Field{Holder: x, Name: fd.Name, Desc: fd.Desc, Flags: fd.AccessFlags}
			if fd.AccessFlags&classfile.AccStatic != 0 {
				f.Static = true
				f.StaticSlot = len(x.Statics)
				x.Statics = append(x.Statics, nil)
				if cvi := fd.ConstantValue(); cvi != 0 {
					v, err := k.constPoolPrimitive(x.CF, cvi, fd.Desc)
					if err != nil {
						return fail(err)
					}
					x.Statics[f.StaticSlot] = v
				}
			} else {
				f.Slot = slot
				slot += SlotCount(fd.Desc) // category-2 fields take two slots
				if x == c {
					c.OwnFields = append(c.OwnFields, f)
				}
			}
			c.flatField(f)
		}
	}
	c.layoutSize = slot

	for i := range cf.Methods {
		md := &cf.Methods[i]
		m := &Method{Holder: c, Name: md.Name, Desc: md.Desc, Flags: md.AccessFlags, CF: cf, Code: md.Code}
		c.Methods = append(c.Methods, m)
		c.methodsByKey[memberKey(m.Name, m.Desc)] = m
	}

	// Java defaults for statics without ConstantValue.
	for _, f := range c.fieldsByKey {
		if f.Static && f.Holder == c && c.Statics[f.StaticSlot] == nil {
			c.Statics[f.StaticSlot] = zeroValue(f.Desc)
		}
	}

	if !k.opts.SkipVerify {
		if err := k.verifyLoaded(c, cf); err != nil {
			return fail(fmt.Errorf("VerifyError in %s: %w", c.Name, err))
		}
	}
	// Post-registration hook (emitter layer attaches generated bodies to
	// lazily-loaded classes). Fired outside k.mu: the hook resolves
	// members on the fresh class only.
	k.mu.Lock()
	hook := k.classLoadHook
	k.mu.Unlock()
	if hook != nil {
		hook(c)
	}
	return c, nil
}

// verifyLoaded runs the structural verifier against a freshly loaded class.
// The resolver adapter answers from the currently-registered class graph
// only (registry reads; no recursive loading during verification).
func (k *Kernel) verifyLoaded(c *Class, cf *classfile.ClassFile) error {
	adapter := &registryResolver{k: k}
	return verify.Verify(cf, adapter)
}

type registryResolver struct{ k *Kernel }

func (r *registryResolver) Known(name string) bool {
	_, ok := r.k.ClassByName(name)
	return ok
}

func (r *registryResolver) IsInterface(name string) bool {
	c, ok := r.k.ClassByName(name)
	if !ok {
		return false
	}
	return c.Flags&classfile.AccInterface != 0
}

func (r *registryResolver) IsSubclass(child, anc string) bool {
	cc, ok1 := r.k.ClassByName(child)
	ac, ok2 := r.k.ClassByName(anc)
	if !ok1 || !ok2 {
		return false
	}
	for x := cc; x != nil; x = x.Super {
		if x == ac {
			return true
		}
	}
	return r.k.implementsIface(cc, ac)
}

func (c *Class) flatField(f *Field) {
	if c.fieldsByKey == nil {
		c.fieldsByKey = make(map[string]*Field)
	}
	c.fieldsByKey[memberKey(f.Name, f.Desc)] = f
}

// ---- Member resolution ---------------------------------------------------

// ResolveField resolves name+desc starting at start, walking supers then
// interfaces (JVMS §5.4.3.2, M0 simplification: no linkage errors).
func (k *Kernel) ResolveField(start *Class, name, desc string) (*Field, error) {
	key := memberKey(name, desc)
	for x := start; x != nil; x = x.Super {
		if f, ok := x.fieldsByKey[key]; ok {
			return f, nil
		}
	}
	seen := make(map[*Class]bool)
	var walkIfaces func(cs []*Class) *Field
	walkIfaces = func(cs []*Class) *Field {
		for _, ic := range cs {
			if seen[ic] {
				continue
			}
			seen[ic] = true
			if f, ok := ic.fieldsByKey[key]; ok {
				return f
			}
			if f := walkIfaces(ic.Ifaces); f != nil {
				return f
			}
		}
		return nil
	}
	if f := walkIfaces(start.Ifaces); f != nil {
		return f, nil
	}
	return nil, fmt.Errorf("NoSuchFieldError: %s.%s:%s", start.Name, name, desc)
}

// ResolveMethod resolves name+desc starting at start (supers, then
// interfaces) — JVMS §5.4.3.3/§5.4.3.4.
func (k *Kernel) ResolveMethod(start *Class, name, desc string) (*Method, error) {
	key := memberKey(name, desc)
	for x := start; x != nil; x = x.Super {
		if m, ok := x.methodsByKey[key]; ok {
			return m, nil
		}
	}
	seen := make(map[*Class]bool)
	var walkIfaces func(cs []*Class) *Method
	walkIfaces = func(cs []*Class) *Method {
		for _, ic := range cs {
			if seen[ic] {
				continue
			}
			seen[ic] = true
			if m, ok := ic.methodsByKey[key]; ok {
				return m
			}
			if m := walkIfaces(ic.Ifaces); m != nil {
				return m
			}
		}
		return nil
	}
	if m := walkIfaces(start.Ifaces); m != nil {
		return m, nil
	}
	return nil, fmt.Errorf("NoSuchMethodError: %s.%s:%s", start.Name, name, desc)
}

// ---- Object construction -------------------------------------------------

// ClassObjectOf returns the unique java.lang.Class instance for c,
// creating it on first use (identity stable for == comparisons).
func (k *Kernel) ClassObjectOf(c *Class) (*Instance, error) {
	if o := c.classObj.Load(); o != nil {
		return o, nil
	}
	cls, ok := k.ClassByName("java/lang/Class")
	if !ok {
		return nil, fmt.Errorf("kernel: bootstrap Class missing")
	}
	in, err := k.NewInstance(cls)
	if err != nil {
		return nil, err
	}
	in.Payload = c
	if !c.classObj.CompareAndSwap(nil, in) {
		return c.classObj.Load(), nil
	}
	return in, nil
}

// NewInstance allocates an object of class c. Fields receive their Java
// default values per descriptor (0/0L/false… for primitives, null for refs).
func (k *Kernel) NewInstance(c *Class) (*Instance, error) {
	if c.IsArray || c.Name == "java/lang/Class" ||
		c.Flags&(classfile.AccAbstract|classfile.AccInterface) != 0 {
		return nil, fmt.Errorf("cannot instantiate %s", c.Name)
	}
	in := &Instance{}
	in.Class = c
	in.Fields = make([]Value, c.layoutSize)
	for _, f := range c.fieldsByKey {
		if !f.Static && f.Slot < c.layoutSize {
			in.Fields[f.Slot] = zeroValue(f.Desc)
		}
	}
	return in, nil
}

// NewArray allocates an array by component descriptor.
func (k *Kernel) NewArray(compDesc string, n int) (*ArrayObj, error) {
	if n < 0 {
		obj, err := k.NewThrowable("java/lang/NegativeArraySizeException", itoa(n))
		if err != nil {
			return nil, err
		}
		return nil, &Thrown{Obj: obj}
	}
	a := &ArrayObj{CompDesc: compDesc, Elems: make([]Value, n)}
	a.Class = k.ArrayClassOf(compDesc)
	z := zeroValue(compDesc)
	for i := range a.Elems {
		a.Elems[i] = z
	}
	return a, nil
}

func itoa(n int) string { return fmt.Sprintf("%d", n) }

// MakeJString allocates a String (not interned).
func (k *Kernel) MakeJString(chars []uint16) *JString {
	s := &JString{Chars: chars}
	strCls := k.lookupClass("java/lang/String")
	if strCls == nil {
		panic("kernel: bootstrap String class missing")
	}
	s.Class = strCls
	return s
}

// MakeJStringFromGo allocates a String from Go text.
func (k *Kernel) MakeJStringFromGo(text string) *JString {
	return k.MakeJString(utf16Encode([]rune(text)))
}

// Intern returns the interned String for chars (ldc semantics).
func (k *Kernel) Intern(chars []uint16) *JString {
	key := utf16Key(chars)
	k.mu.Lock()
	defer k.mu.Unlock()
	if s, ok := k.strPool[key]; ok {
		return s
	}
	s := k.MakeJString(append([]uint16(nil), chars...))
	k.strPool[key] = s
	return s
}

// InternGo interns from Go text.
func (k *Kernel) InternGo(text string) *JString { return k.Intern(utf16Encode([]rune(text))) }

func utf16Key(u []uint16) string {
	b := make([]byte, len(u)*2)
	for i, v := range u {
		b[i*2] = byte(v)
		b[i*2+1] = byte(v >> 8)
	}
	return string(b)
}

// IntegerOf implements Integer.valueOf boxing with the [-128,127] cache.
func (k *Kernel) IntegerOf(v int32) *Instance {
	if v >= -128 && v <= 127 {
		if b := k.intBox[v+128]; b != nil {
			return b
		}
		k.mu.Lock()
		if b := k.intBox[v+128]; b != nil { // lost race re-check
			k.mu.Unlock()
			return b
		}
		b := k.newBoxedIntLocked(v)
		k.intBox[v+128] = b
		k.mu.Unlock()
		return b
	}
	// Out-of-cache boxes touch no shared state: allocation only. The
	// mutex here was pure serialization overhead on hot map workloads
	// (P-0009 U5).
	return k.newBoxedIntFast(v)
}

func (k *Kernel) newBoxedIntLocked(v int32) *Instance {
	cls := k.lookupClass("java/lang/Integer")
	if cls == nil {
		panic("kernel: bootstrap Integer class missing")
	}
	b := &Instance{}
	b.Class = cls
	b.Fields = make([]Value, cls.layoutSize)
	f := cls.fieldsByKey[memberKey("value", "I")]
	if f != nil {
		b.Fields[f.Slot] = v
	} else {
		b.Fields[0] = v
	}
	return b
}

// integerCls caches the bootstrap Integer class (atomic; the registry
// entry is immutable once defined).
func (k *Kernel) integerCls() *Class {
	if c := k.intCls.Load(); c != nil {
		return c
	}
	c := k.lookupClass("java/lang/Integer")
	if c == nil {
		panic("kernel: bootstrap Integer class missing")
	}
	k.intCls.Store(c)
	return c
}

// newBoxedIntFast is the out-of-cache boxing fast path: no locking, and a
// specialized single-field layout instead of the generic default-fill walk.
// Integer's synthesized layout is fixed at one slot holding the value, so
// this is semantically identical to NewInstance + putfield.
func (k *Kernel) newBoxedIntFast(v int32) *Instance {
	b := &Instance{}
	b.Class = k.integerCls()
	b.Fields = []Value{v}
	return b
}

// IntValueOf extracts the int from a boxed Integer.
func IntValueOf(v Value) (int32, bool) {
	in, ok := v.(*Instance)
	if !ok || in.Class.Name != "java/lang/Integer" {
		return 0, false
	}
	i, ok := in.Fields[0].(int32)
	return i, ok
}

// NewThrowable allocates a throwable without running constructors.
func (k *Kernel) NewThrowable(className, msg string) (*Instance, error) {
	c, ok := k.ClassByName(className)
	if !ok {
		return nil, fmt.Errorf("exception class %s not found", className)
	}
	in, err := k.NewInstance(c)
	if err != nil {
		return nil, err
	}
	if msg != "" {
		if f, ferr := k.ResolveField(c, "detailMessage", "Ljava/lang/String;"); ferr == nil && f != nil {
			in.Fields[f.Slot] = k.MakeJStringFromGo(msg)
		}
	}
	return in, nil
}

// Throw builds a Thrown error for the given exception class.
func (k *Kernel) Throw(className, msg string) error {
	obj, err := k.NewThrowable(className, msg)
	if err != nil {
		return err
	}
	return &Thrown{Obj: obj}
}

// ---- Type predicates ------------------------------------------------------

// IsInstance implements instanceof/checkcast semantics including arrays
// and interfaces.
func (k *Kernel) IsInstance(v Value, cls *Class) bool {
	switch o := v.(type) {
	case nil:
		return false
	case *JString:
		return k.classIsa(o.Class, cls)
	case *Instance:
		return k.classIsa(o.Class, cls)
	case *ArrayObj:
		return k.arrayIsa(o, cls)
	default:
		return false
	}
}

func (k *Kernel) classIsa(c, target *Class) bool {
	if c == target || c.IsSubclassOf(target) {
		return true
	}
	return k.implementsIface(c, target)
}

func (k *Kernel) implementsIface(c, target *Class) bool {
	seen := make(map[*Class]bool)
	var visit func(cs []*Class) bool
	visit = func(cs []*Class) bool {
		for _, ic := range cs {
			if seen[ic] {
				continue
			}
			seen[ic] = true
			if ic == target || k.ifaceExtends(ic, target) {
				return true
			}
			if visit(ic.Ifaces) {
				return true
			}
		}
		return false
	}
	for x := c; x != nil; x = x.Super {
		if visit(x.Ifaces) {
			return true
		}
	}
	return false
}

func (k *Kernel) ifaceExtends(i, target *Class) bool {
	for _, sup := range i.Ifaces {
		if sup == target || k.ifaceExtends(sup, target) {
			return true
		}
	}
	return false
}

func (k *Kernel) arrayIsa(a *ArrayObj, target *Class) bool {
	if !target.IsArray {
		return target.Name == "java/lang/Object"
	}
	if a.CompDesc == target.CompDesc {
		return true
	}
	// Covariance: X[] <: Y[] iff X <: Y for reference components.
	if !(strings.HasPrefix(a.CompDesc, "L") || strings.HasPrefix(a.CompDesc, "[")) {
		return false // primitive component arrays are invariant
	}
	compCls, err := k.ResolveClass(a.CompDesc)
	if err != nil {
		return false
	}
	tgtCls, ok := k.ClassByName(target.CompDesc)
	if !ok {
		return false
	}
	return k.classIsa(compCls, tgtCls)
}

// ---- Misc helpers ----------------------------------------------------------

// constPoolPrimitive converts a constant-pool entry into a Value matching
// the field descriptor desc (for ConstantValue attributes).
func (k *Kernel) constPoolPrimitive(cf interface {
	Entry(i uint16) (*classfile.ConstEntry, error)
}, idx uint16, desc string) (Value, error) {
	e, err := cf.Entry(idx)
	if err != nil {
		return nil, err
	}
	switch desc[0] {
	case 'B', 'S', 'C', 'I', 'Z':
		if e.Tag != classfile.CInteger {
			return nil, fmt.Errorf("constant type mismatch for %s", desc)
		}
		v := e.IntVal
		switch desc[0] {
		case 'Z':
			v = bool32(v != 0)
		case 'C':
			v &= 0xFFFF
		case 'B':
			v = int32(int8(v))
		case 'S':
			v = int32(int16(v))
		}
		return v, nil
	case 'J':
		if e.Tag != classfile.CLong {
			return nil, fmt.Errorf("constant type mismatch for J")
		}
		return e.LongVal, nil
	case 'F':
		if e.Tag != classfile.CFloat {
			return nil, fmt.Errorf("constant type mismatch for F")
		}
		return e.FloatVal, nil
	case 'D':
		if e.Tag != classfile.CDouble {
			return nil, fmt.Errorf("constant type mismatch for D")
		}
		return e.DoubleVal, nil
	case 'L':
		if e.Tag != classfile.CString {
			return nil, fmt.Errorf("constant type mismatch for %s", desc)
		}
		return k.InternGo(e.Str), nil
	}
	return nil, fmt.Errorf("unsupported ConstantValue descriptor %s", desc)
}

func bool32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}
