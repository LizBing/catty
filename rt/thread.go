package rt

import "fmt"

// Thread is one logical Java thread. Gate 1 runs one thread, but the explicit
// parameter is retained for the threaded profiles.
type Thread struct {
	RT     *Runtime
	ID     int64
	depth  int
	frames []FrameInfo
}

const maxJavaDepth = 20000

// NewThread creates a logical Java thread with a fresh identity.
func (r *Runtime) NewThread() *Thread {
	r.nextThread++
	return &Thread{RT: r, ID: r.nextThread}
}

// MainThread returns the launcher thread, creating it on first use.
func (r *Runtime) MainThread() *Thread {
	if r.mainThread == nil {
		r.mainThread = r.NewThread()
	}
	return r.mainThread
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

// CloseTrace flushes and closes the trace channel.
func (r *Runtime) CloseTrace() {
	if r.trace != nil {
		r.trace.close()
		r.trace = nil
	}
}

// Call is the one owner of Java frame push/pop and trace enter/exit for every
// engine. The selected method's entry runs synchronously on the same thread.
func (t *Thread) Call(m *Method, receiver *Object, args []Value) Outcome {
	if m == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if t.depth >= maxJavaDepth {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/StackOverflowError"))
	}
	t.depth++
	t.frames = append(t.frames, FrameInfo{Method: m})
	depth := t.depth - 1
	if tr := t.trace(); tr != nil {
		tr.enter(t.ID, m, depth)
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

func (r *Runtime) bootstrapThrowable(class string) *Object {
	return r.NewThrowable(nil, class)
}

// EnsureInitialized performs the single-threaded JVM initialization procedure
// for c. It returns a Java throwable on failure and nil on success. The owning
// thread may reenter while state is initializing.
func (t *Thread) EnsureInitialized(c *Class) *Object {
	if c == nil || c.IsPrimitive || c.IsArray {
		return nil
	}
	switch c.State {
	case classInitialized:
		return nil
	case classInitializing:
		return nil
	case classErroneous:
		return t.RT.NewThrowable(t, "java/lang/NoClassDefFoundError")
	}
	c.State = classInitializing
	if c.Super != nil {
		if thrown := t.EnsureInitialized(c.Super); thrown != nil {
			c.State = classErroneous
			c.InitError = thrown
			return thrown
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
			c.State = classErroneous
			c.InitError = err
			return err
		}
	}
	c.State = classInitialized
	return nil
}

// ResolveMethod resolves owner/name/descriptor through the class hierarchy.
func (r *Runtime) ResolveMethod(owner, name, descriptor string) (*Method, error) {
	m, thrown := r.resolveMethod(owner, name, descriptor)
	if thrown != nil {
		return nil, fmt.Errorf("catty: cannot resolve method %s.%s%s", owner, name, descriptor)
	}
	return m, nil
}

func (r *Runtime) resolveMethod(owner, name, descriptor string) (*Method, *Object) {
	c, err := r.classForName(owner)
	if err != nil {
		return nil, r.bootstrapThrowable("java/lang/NoSuchMethodError")
	}
	key := name + descriptor
	for cur := c; cur != nil; cur = cur.Super {
		if m := cur.MethodIndex[key]; m != nil {
			return m, nil
		}
	}
	if m := findInterfaceMethod(c, key); m != nil {
		return m, nil
	}
	return nil, r.bootstrapThrowable("java/lang/NoSuchMethodError")
}

func findInterfaceMethod(c *Class, key string) *Method {
	seen := map[*Class]bool{}
	var visit func(c *Class) *Method
	visit = func(c *Class) *Method {
		if c == nil || seen[c] {
			return nil
		}
		seen[c] = true
		for _, in := range c.Interfaces {
			if m := in.MethodIndex[key]; m != nil {
				return m
			}
			if m := visit(in); m != nil {
				return m
			}
		}
		return visit(c.Super)
	}
	return visit(c)
}

// selectVirtual performs receiver-based method selection for invokevirtual and
// invokeinterface after resolution has produced the declaration.
func (r *Runtime) selectVirtual(receiver *Class, declaration *Method) (*Method, *Object) {
	key := declaration.Name + declaration.Descriptor
	for c := receiver; c != nil; c = c.Super {
		if m := c.MethodIndex[key]; m != nil && m.Flags&AccStatic == 0 {
			if m.Flags&AccAbstract != 0 {
				return nil, r.bootstrapThrowable("java/lang/AbstractMethodError")
			}
			return m, nil
		}
	}
	if m := findInterfaceMethod(receiver, key); m != nil && m.Flags&AccStatic == 0 {
		if m.Flags&AccAbstract != 0 {
			return nil, r.bootstrapThrowable("java/lang/AbstractMethodError")
		}
		return m, nil
	}
	return nil, r.bootstrapThrowable("java/lang/AbstractMethodError")
}

func (r *Runtime) resolveField(owner, name, descriptor string) (*Field, *Object) {
	c, err := r.classForName(owner)
	if err != nil {
		return nil, r.bootstrapThrowable("java/lang/NoSuchFieldError")
	}
	for cur := c; cur != nil; cur = cur.Super {
		for _, f := range cur.InstanceField {
			if f.Name == name && f.Descriptor == descriptor {
				return f, nil
			}
		}
		for _, f := range cur.StaticField {
			if f.Name == name && f.Descriptor == descriptor {
				return f, nil
			}
		}
	}
	if f := findInterfaceField(c, name, descriptor); f != nil {
		return f, nil
	}
	return nil, r.bootstrapThrowable("java/lang/NoSuchFieldError")
}

func findInterfaceField(c *Class, name, descriptor string) *Field {
	seen := map[*Class]bool{}
	var visit func(c *Class) *Field
	visit = func(c *Class) *Field {
		if c == nil || seen[c] {
			return nil
		}
		seen[c] = true
		for _, in := range c.Interfaces {
			for _, f := range in.StaticField {
				if f.Name == name && f.Descriptor == descriptor {
					return f
				}
			}
			if f := visit(in); f != nil {
				return f
			}
		}
		return visit(c.Super)
	}
	return visit(c)
}

func refOrNull(t *Thread, v Value) (*Object, *Object) {
	if v.Kind != KindRef {
		return nil, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	if v.Ref == nil {
		return nil, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return v.Ref, nil
}

// GetField implements getfield.
func (t *Thread) GetField(ref Value, owner, name, descriptor string) (Value, *Object) {
	obj, thrown := refOrNull(t, ref)
	if thrown != nil {
		return Value{}, thrown
	}
	f, thrown := t.RT.resolveField(owner, name, descriptor)
	if thrown != nil {
		return Value{}, thrown
	}
	if f.Static {
		return Value{}, t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if f.Offset < 0 || f.Offset >= len(obj.Fields) {
		return Value{}, t.RT.NewThrowable(t, "java/lang/NoSuchFieldError")
	}
	return obj.Fields[f.Offset], nil
}

// PutField implements putfield.
func (t *Thread) PutField(ref Value, owner, name, descriptor string, v Value) *Object {
	obj, thrown := refOrNull(t, ref)
	if thrown != nil {
		return thrown
	}
	f, thrown := t.RT.resolveField(owner, name, descriptor)
	if thrown != nil {
		return thrown
	}
	if f.Static {
		return t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if f.Offset < 0 || f.Offset >= len(obj.Fields) {
		return t.RT.NewThrowable(t, "java/lang/NoSuchFieldError")
	}
	obj.Fields[f.Offset] = v
	return nil
}

// GetStatic implements getstatic, including class initialization.
func (t *Thread) GetStatic(owner, name, descriptor string) (Value, *Object) {
	f, thrown := t.RT.resolveField(owner, name, descriptor)
	if thrown != nil {
		return Value{}, thrown
	}
	if !f.Static {
		return Value{}, t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if thrown := t.EnsureInitialized(f.Owner); thrown != nil {
		return Value{}, thrown
	}
	return f.Owner.Statics[f.Index], nil
}

// PutStatic implements putstatic, including class initialization.
func (t *Thread) PutStatic(owner, name, descriptor string, v Value) *Object {
	f, thrown := t.RT.resolveField(owner, name, descriptor)
	if thrown != nil {
		return thrown
	}
	if !f.Static {
		return t.RT.NewThrowable(t, "java/lang/IncompatibleClassChangeError")
	}
	if thrown := t.EnsureInitialized(f.Owner); thrown != nil {
		return thrown
	}
	f.Owner.Statics[f.Index] = v
	return nil
}

// NewObject implements new.
func (t *Thread) NewObject(className string) (*Object, *Object) {
	c, err := t.RT.classForName(className)
	if err != nil {
		return nil, t.RT.NewThrowable(t, "java/lang/NoClassDefFoundError")
	}
	if c.IsInterface() || c.Flags&AccAbstract != 0 || c.IsArray || c.IsPrimitive {
		return nil, t.RT.NewThrowable(t, "java/lang/InstantiationError")
	}
	if thrown := t.EnsureInitialized(c); thrown != nil {
		return nil, thrown
	}
	return t.RT.NewInstance(c), nil
}

// NewArray implements anewarray.
func (t *Thread) NewArray(count int32, componentName string) (*Object, *Object) {
	if count < 0 {
		return nil, t.RT.NewThrowable(t, "java/lang/NegativeArraySizeException")
	}
	component, err := t.RT.classForName(componentName)
	if err != nil {
		return nil, t.RT.NewThrowable(t, "java/lang/NoClassDefFoundError")
	}
	return t.RT.NewArray(component, int(count)), nil
}

// ArrayLoad implements aaload.
func (t *Thread) ArrayLoad(array Value, index int32) (Value, *Object) {
	obj, thrown := refOrNull(t, array)
	if thrown != nil {
		return Value{}, thrown
	}
	if obj.Array == nil {
		return Value{}, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	if index < 0 || int(index) >= len(obj.Array.Elems) {
		return Value{}, t.RT.NewThrowable(t, "java/lang/ArrayIndexOutOfBoundsException")
	}
	return obj.Array.Elems[index], nil
}

// ArrayStore implements aastore, including the reified component check.
func (t *Thread) ArrayStore(array Value, index int32, v Value) *Object {
	obj, thrown := refOrNull(t, array)
	if thrown != nil {
		return thrown
	}
	if obj.Array == nil {
		return t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	if index < 0 || int(index) >= len(obj.Array.Elems) {
		return t.RT.NewThrowable(t, "java/lang/ArrayIndexOutOfBoundsException")
	}
	if v.Kind != KindRef {
		return t.RT.NewThrowable(t, "java/lang/ArrayStoreException")
	}
	if v.Ref != nil && !t.RT.IsAssignable(v.Ref.Class, obj.Array.Component) {
		return t.RT.NewThrowable(t, "java/lang/ArrayStoreException")
	}
	obj.Array.Elems[index] = v
	return nil
}

// ArrayLength implements arraylength.
func (t *Thread) ArrayLength(array Value) (int32, *Object) {
	obj, thrown := refOrNull(t, array)
	if thrown != nil {
		return 0, thrown
	}
	if obj.Array == nil {
		return 0, t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return int32(len(obj.Array.Elems)), nil
}

// ClassConstant implements ldc of a Class constant.
func (t *Thread) ClassConstant(name string) Value {
	c, err := t.RT.classForName(name)
	if err != nil {
		panic("catty: class constant refers to unknown class " + name)
	}
	return RefValue(c.Mirror)
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
// escapes this method. The caller is responsible for replacing its operand
// stack with the throwable at the handler entry.
func (t *Thread) HandleThrow(m *Method, throwPC int32, thrown *Object) (int32, bool) {
	if thrown == nil {
		return 0, false
	}
	for _, h := range m.Handlers {
		if throwPC < h.StartPC || throwPC >= h.EndPC {
			continue
		}
		if h.CatchType != "" {
			catchClass, err := t.RT.classForName(h.CatchType)
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

func (t *Thread) invocationEnv(caller *Method, pc int32) {
	t.SetPC(pc)
}

// InvokeStatic implements invokestatic.
func (t *Thread) InvokeStatic(caller *Method, pc int32, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	m, thrown := t.RT.resolveMethod(owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if thrown := t.EnsureInitialized(m.Owner); thrown != nil {
		return ThrowOutcome(thrown)
	}
	return t.Call(m, nil, args)
}

// InvokeSpecial implements invokespecial.
func (t *Thread) InvokeSpecial(caller *Method, pc int32, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	m, thrown := t.RT.resolveMethod(owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	return t.Call(m, recv, args)
}

// InvokeVirtual implements invokevirtual.
func (t *Thread) InvokeVirtual(caller *Method, pc int32, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	decl, thrown := t.RT.resolveMethod(owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	target, thrown := t.RT.selectVirtual(recv.Class, decl)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	return t.Call(target, recv, args)
}

// InvokeInterface implements invokeinterface and emits the direct dispatch
// evidence before entering the selected body.
func (t *Thread) InvokeInterface(caller *Method, pc int32, receiver Value, owner, name, descriptor string, args []Value) Outcome {
	t.invocationEnv(caller, pc)
	recv, thrown := refOrNull(t, receiver)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	decl, thrown := t.RT.resolveMethod(owner, name, descriptor)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	target, thrown := t.RT.selectVirtual(recv.Class, decl)
	if thrown != nil {
		return ThrowOutcome(thrown)
	}
	if tr := t.trace(); tr != nil {
		tr.dispatch(t.ID, "invokeinterface", decl, target)
	}
	return t.Call(target, recv, args)
}

// StaticValue resolves and returns one static field value (used by tests and
// host bindings).
func (r *Runtime) StaticValue(owner, name, descriptor string) (Value, error) {
	f, thrown := r.resolveField(owner, name, descriptor)
	if thrown != nil {
		return Value{}, fmt.Errorf("catty: cannot resolve static field %s.%s%s", owner, name, descriptor)
	}
	if !f.Static {
		return Value{}, fmt.Errorf("catty: field %s.%s%s is not static", owner, name, descriptor)
	}
	return f.Owner.Statics[f.Index], nil
}
