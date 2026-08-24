// Reflection minimal surface (P-0011): java/lang/Class metadata natives
// plus java/lang/reflect.{Field,Method,Constructor} synthesized classes.
//
// Design notes:
//   - Metadata is read from the runtime's own *Class/*Method/*Field
//     structures — both synthesized and classfile-loaded classes reflect
//     uniformly.
//   - Field/Method/Constructor instances carry the runtime object in
//     Payload; identity is stable per member.
//   - Primitive pseudo-classes (int.class …) are cached per descriptor and
//     render their keyword name via getName.
//   - v1 deviations (DEV-0010): access control ignored (private settable),
//     getDeclared* only (no inherited-member walk), no generics/signature.
package kernel

import (
	"fmt"
	"sort"
	"strings"

	"catty/internal/classfile"
)

const classfileAccPublicFinal = classfile.AccPublic | classfile.AccFinal

const reflectionException = "java/lang/ClassNotFoundException"

// primitiveDescs maps Java keyword names to descriptors for the primitive
// Class constants (plus void).
var primitiveDescs = []struct{ name, desc string }{
	{"int", "I"}, {"long", "J"}, {"float", "F"}, {"double", "D"},
	{"boolean", "Z"}, {"byte", "B"}, {"char", "C"}, {"short", "S"},
	{"void", "V"},
}

// primitiveClass builds (once) the Class mirror for a primitive keyword.
// It deliberately avoids k.mu and the class registry: StaticInit hooks run
// while DefineClass holds k.mu, so touching either would self-deadlock.
// The pseudo-class stays out of the registry — forName resolves primitive
// keywords before consulting it, and nothing else looks them up by name.
func (k *Kernel) primitiveClass(desc, name string) (*Instance, error) {
	if v, ok := k.primOnce.Load(desc); ok {
		return v.(*Instance), nil
	}
	cls := &Class{Name: name, Flags: classfileAccPublicFinal}
	// Primitive mirrors share java/lang/Class's method table so
	// getName/getSimpleName/… resolve on them too (when defined already).
	if cc := k.lookupClass("java/lang/Class"); cc != nil {
		cls.Methods = cc.Methods
		cls.methodsByKey = make(map[string]*Method, len(cc.Methods))
		for _, m := range cc.Methods {
			cls.methodsByKey[memberKey(m.Name, m.Desc)] = m
		}
	}
	cls.setState(StateInitialized)
	in := &Instance{}
	in.Class = cls
	in.Payload = &primitiveInfo{name: name, desc: desc}
	actual, _ := k.primOnce.LoadOrStore(desc, in)
	return actual.(*Instance), nil
}

type primitiveInfo struct{ name, desc string }

func classPayload(in *Instance) (*Class, bool) {
	c, ok := in.Payload.(*Class)
	return c, ok
}

// descToClass resolves a field/method descriptor component into its Class
// mirror (primitives, reference, and array forms).
func (k *Kernel) descToClass(desc string) (*Instance, error) {
	if len(desc) == 1 {
		for _, p := range primitiveDescs {
			if p.desc == desc {
				return k.primitiveClass(desc, p.name)
			}
		}
		return nil, fmt.Errorf("reflect: unknown primitive descriptor %q", desc)
	}
	if strings.HasPrefix(desc, "[") {
		arrCls := k.ArrayClassOf(desc)
		return k.ClassObjectOf(arrCls)
	}
	if strings.HasPrefix(desc, "L") && strings.HasSuffix(desc, ";") {
		internal := desc[1 : len(desc)-1]
		c, err := k.ResolveClass(internal)
		if err != nil {
			return nil, err
		}
		return k.ClassObjectOf(c)
	}
	return nil, fmt.Errorf("reflect: bad descriptor %q", desc)
}

// boxValue wraps a raw field/return value for reflection results according
// to the declared descriptor character (V → null).
func (k *Kernel) boxValue(ctx *CallContext, descByte byte, v Value) Value {
	switch descByte {
	case 'I':
		return k.IntegerOf(v.(int32))
	case 'J':
		return k.boxLong(v.(int64))
	case 'Z':
		if v.(int32) != 0 {
			return k.IntegerOf(1)
		}
		return k.IntegerOf(0)
	case 'B':
		return k.boxViaWrapper("java/lang/Byte", int64(v.(int32)))
	case 'C':
		return k.boxViaWrapper("java/lang/Character", int64(v.(int32)))
	case 'S':
		return k.boxViaWrapper("java/lang/Short", int64(v.(int32)))
	case 'F':
		return k.boxFloat(v.(float32))
	case 'D':
		return k.boxDouble(v.(float64))
	case 'V':
		return nil
	}
	return v
}

// unboxArg adapts a boxed/reflection argument to the raw storage form a
// descriptor demands (Integer → int32 etc.). Reference descriptors pass
// through untouched.
func unboxArg(v Value, desc string) (Value, error) {
	if v == nil {
		if len(desc) == 1 {
			return nil, fmt.Errorf("interop: null argument")
		}
		return nil, nil
	}
	switch desc[0] {
	case 'I', 'Z', 'B', 'C', 'S':
		if n, ok := v.(int32); ok {
			return n, nil
		}
		if in, ok := v.(*Instance); ok {
			if f := in.fieldByName("value"); f != nil {
				if n, ok := f.(int32); ok {
					return n, nil
				}
			}
		}
		return nil, fmt.Errorf("not an int-boxed value (%T)", v)
	case 'J':
		if n, ok := v.(int64); ok {
			return n, nil
		}
		if in, ok := v.(*Instance); ok {
			if f := in.fieldByName("value"); f != nil {
				if n, ok := f.(int64); ok {
					return n, nil
				}
			}
		}
		return nil, fmt.Errorf("not a long-boxed value (%T)", v)
	case 'F':
		if n, ok := v.(float32); ok {
			return n, nil
		}
	case 'D':
		if n, ok := v.(float64); ok {
			return n, nil
		}
	}
	return v, nil
}

func (k *Kernel) boxLong(v int64) Value {
	if m, err := k.ResolveMethod(mustLookup(k, "java/lang/Long"), "valueOf", "(J)Ljava/lang/Long;"); err == nil {
		if out, ierr := k.InvokeAs(nil, m, nil, []Value{v}); ierr == nil {
			return out
		}
	}
	return v
}

func (k *Kernel) boxFloat(v float32) Value {
	if m, err := k.ResolveMethod(mustLookup(k, "java/lang/Float"), "valueOf", "(F)Ljava/lang/Float;"); err == nil {
		if out, ierr := k.InvokeAs(nil, m, nil, []Value{v}); ierr == nil {
			return out
		}
	}
	return v
}

func (k *Kernel) boxDouble(v float64) Value {
	if m, err := k.ResolveMethod(mustLookup(k, "java/lang/Double"), "valueOf", "(D)Ljava/lang/Double;"); err == nil {
		if out, ierr := k.InvokeAs(nil, m, nil, []Value{v}); ierr == nil {
			return out
		}
	}
	return v
}

func (k *Kernel) boxViaWrapper(cls string, v int64) Value {
	if m, err := k.ResolveMethod(mustLookup(k, cls), "valueOf", "(J)"+wrapperDesc(cls)); err == nil {
		if out, ierr := k.InvokeAs(nil, m, nil, []Value{v}); ierr == nil {
			return out
		}
	}
	return v
}

func wrapperDesc(cls string) string { return "L" + cls + ";" }

func mustLookup(k *Kernel, name string) *Class {
	c, _ := k.ClassByName(name)
	return c
}

// ---- registration ----------------------------------------------------------

func registerReflection(k *Kernel) error {
	// ClassNotFoundException first: forName throws it.
	if _, err := k.DefineClass(&ClassDef{
		Name:  reflectionException,
		Super: "java/lang/Exception",
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: 0x0001,
				Native: natThrowableInitString},
		},
	}); err != nil {
		return err
	}

	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/Class",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natClassGetName},
			{Name: "getSimpleName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natClassGetSimpleName},
			{Name: "forName", Desc: "(Ljava/lang/String;)Ljava/lang/Class;", Flags: 0x0009, Native: natClassForName},
			{Name: "getDeclaredFields", Desc: "()[Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassDeclaredFields},
			{Name: "getField", Desc: "(Ljava/lang/String;)Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassGetField},
			{Name: "getDeclaredMethods", Desc: "()[Ljava/lang/reflect/Method;", Flags: 0x0001, Native: natClassDeclaredMethods},
			{Name: "getDeclaredConstructors", Desc: "()[Ljava/lang/reflect/Constructor;", Flags: 0x0001, Native: natClassDeclaredConstructors},
			{Name: "isInstance", Desc: "(Ljava/lang/Object;)Z", Flags: 0x0001, Native: natClassIsInstance},
			{Name: "getSuperclass", Desc: "()Ljava/lang/Class;", Flags: 0x0001, Native: natClassGetSuperclass},
			{Name: "newInstance", Desc: "()Ljava/lang/Object;", Flags: 0x0001, Native: natClassNewInstance},
		},
	}); err != nil {
		return err
	}

	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/reflect/Field",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natFieldGetName},
			{Name: "getType", Desc: "()Ljava/lang/Class;", Flags: 0x0001, Native: natFieldGetType},
			{Name: "get", Desc: "(Ljava/lang/Object;)Ljava/lang/Object;", Flags: 0x0001, Native: natFieldGet},
			{Name: "set", Desc: "(Ljava/lang/Object;Ljava/lang/Object;)V", Flags: 0x0001, Native: natFieldSet},
		},
	}); err != nil {
		return err
	}

	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/reflect/Method",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natMethodGetName},
			{Name: "invoke", Desc: "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", Flags: 0x0001, Native: natMethodInvoke},
		},
	}); err != nil {
		return err
	}

	_, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/reflect/Constructor",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natConstructorGetName},
			{Name: "newInstance", Desc: "([Ljava/lang/Object;)Ljava/lang/Object;", Flags: 0x0001, Native: natConstructorNewInstance},
		},
	})
	return err
}

// ---- java/lang/Class natives ------------------------------------------------

func natClassGetSimpleName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, ctx.Throw("java/lang/RuntimeException", "bad Class receiver")
	}
	if pi, ok := in.Payload.(*primitiveInfo); ok {
		return ctx.NewStringGo(pi.name), nil
	}
	c, ok := in.Payload.(*Class)
	if !ok {
		return nil, ctx.Throw("java/lang/RuntimeException", "Class payload missing")
	}
	n := c.Name
	if strings.HasPrefix(n, "[") {
		// arrays: strip descriptor prefixes, keep component simple name
		for strings.HasPrefix(n, "[") {
			n = n[1:]
		}
		n = strings.TrimSuffix(strings.TrimPrefix(n, "L"), ";")
	}
	if i := strings.LastIndexAny(n, "$/."); i >= 0 {
		n = n[i+1:]
	}
	return ctx.NewStringGo(n), nil
}

func classPayloadOrThrow(recv Value) (*Class, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("reflect: Class receiver payload missing")
	}
	c, ok := in.Payload.(*Class)
	if !ok {
		if pi, ok2 := in.Payload.(*primitiveInfo); ok2 {
			return nil, &primitiveSignal{pi}
		}
		return nil, fmt.Errorf("reflect: unexpected Class payload %T", in.Payload)
	}
	return c, nil
}

type primitiveSignal struct{ pi *primitiveInfo }

func (e *primitiveSignal) Error() string { return "primitive class " + e.pi.name }

func natClassForName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := args[0].(*JString)
	if js == nil {
		return nil, ctx.Throw(reflectionException, "null name")
	}
	dottedName := js.Go()
	for _, p := range primitiveDescs {
		if p.name == dottedName {
			return ctx.K.primitiveClass(p.desc, p.name)
		}
	}
	internal := strings.ReplaceAll(dottedName, ".", "/")
	c, err := ctx.K.ResolveClass(internal)
	if err != nil {
		return nil, ctx.Throw(reflectionException, dottedName)
	}
	return ctx.K.ClassObjectOf(c)
}

func natClassDeclaredFields(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		if _, ok := err.(*primitiveSignal); ok {
			return emptyArray(ctx.K, "Ljava/lang/reflect/Field;")
		}
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	fields := c.DeclaredFields
	arr, aerr := ctx.K.NewArray("Ljava/lang/reflect/Field;", len(fields))
	if aerr != nil {
		return nil, aerr
	}
	for i, f := range fields {
		in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Field", f)
		if ierr != nil {
			return nil, ierr
		}
		arr.Elems[i] = in
	}
	return arr, nil
}

func natClassGetField(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	js, _ := args[0].(*JString)
	if js == nil {
		return nil, ctx.Throw("java/lang/NullPointerException", "getField(null)")
	}
	// Walk the hierarchy (public fields only, per Java semantics; access
	// control is not enforced in this runtime — see DEV-0010).
	for x := c; x != nil; x = x.Super {
		for _, f := range x.DeclaredFields {
			if f.Name == js.Go() {
				in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Field", f)
				if ierr != nil {
					return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
				}
				return in, nil
			}
		}
	}
	return nil, ctx.Throw("java/lang.NoSuchFieldException", js.Go())
}

func natClassDeclaredMethods(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		if _, ok := err.(*primitiveSignal); ok {
			return emptyArray(ctx.K, "Ljava/lang/reflect/Method;")
		}
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	var ms []*Method
	for _, m := range c.Methods {
		if m.Name == "<init>" || m.Name == "<clinit>" {
			continue
		}
		ms = append(ms, m)
	}
	sort.Slice(ms, func(i, j int) bool { return ms[i].Name < ms[j].Name })
	arr, aerr := ctx.K.NewArray("Ljava/lang/reflect/Method;", len(ms))
	if aerr != nil {
		return nil, aerr
	}
	for i, m := range ms {
		in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Method", m)
		if ierr != nil {
			return nil, ierr
		}
		arr.Elems[i] = in
	}
	return arr, nil
}

func natClassDeclaredConstructors(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		if _, ok := err.(*primitiveSignal); ok {
			return emptyArray(ctx.K, "Ljava/lang/reflect/Constructor;")
		}
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	var cs []*Method
	for _, m := range c.Methods {
		if m.Name == "<init>" {
			cs = append(cs, m)
		}
	}
	sort.Slice(cs, func(i, j int) bool { return cs[i].Desc < cs[j].Desc })
	arr, aerr := ctx.K.NewArray("Ljava/lang/reflect/Constructor;", len(cs))
	if aerr != nil {
		return nil, aerr
	}
	for i, m := range cs {
		in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Constructor",
			&ctorRef{cls: c, m: m})
		if ierr != nil {
			return nil, ierr
		}
		arr.Elems[i] = in
	}
	return arr, nil
}

func natClassIsInstance(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return int32(0), nil
	}
	if len(args) == 0 || args[0] == nil {
		return int32(0), nil
	}
	if ctx.K.IsInstance(args[0], c) {
		return int32(1), nil
	}
	return int32(0), nil
}

func natClassGetSuperclass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return nil, nil
	}
	if c.Super == nil {
		return nil, nil
	}
	return ctx.K.ClassObjectOf(c.Super)
}

func natClassNewInstance(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	in, ierr := ctx.K.NewInstance(c)
	if ierr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
	}
	init, rerr := ctx.K.ResolveMethod(c, "<init>", "()V")
	if rerr != nil {
		return in, nil // no no-arg constructor: return zero-initialized
	}
	if _, terr := ctx.Invoke(init, in, nil); terr != nil {
		return nil, terr
	}
	return in, nil
}

// ---- Field natives ------------------------------------------------------------

func fieldPayload(recv Value) (*Field, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("reflect: bad Field receiver")
	}
	f, ok := in.Payload.(*Field)
	if !ok {
		return nil, fmt.Errorf("reflect: Field payload missing")
	}
	return f, nil
}

func natFieldGetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, err := fieldPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	return ctx.NewStringGo(f.Name), nil
}

func natFieldGetType(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, err := fieldPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	c, cerr := ctx.K.descToClass(f.Desc)
	if cerr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", cerr.Error())
	}
	return c, nil
}

func natFieldGet(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, ferr := fieldPayload(recv)
	if ferr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
	}
	var raw Value
	if f.Static {
		c := declaringClassOf(recv)
		raw = c.Statics[f.StaticSlot]
	} else {
		if len(args) == 0 || args[0] == nil {
			return nil, ctx.Throw("java/lang/NullPointerException", "Field.get on null instance")
		}
		in := args[0].(*Instance)
		raw = in.Fields[f.Slot]
	}
	return ctx.K.boxValue(ctx, f.Desc[0], raw), nil
}

func natFieldSet(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, ferr := fieldPayload(recv)
	if ferr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
	}
	if len(args) < 2 {
		return nil, ctx.Throw("java/lang/IllegalArgumentException", "set(obj,value)")
	}
	raw, uerr := unboxArg(args[1], f.Desc)
	if uerr != nil {
		return nil, ctx.Throw("java/lang/IllegalArgumentException", uerr.Error())
	}
	if f.Static {
		c := declaringClassOf(recv)
		c.Statics[f.StaticSlot] = raw
		return nil, nil
	}
	if args[0] == nil {
		return nil, ctx.Throw("java/lang/NullPointerException", "Field.set on null instance")
	}
	in := args[0].(*Instance)
	in.Fields[f.Slot] = raw
	return nil, nil
}

// ---- Method / Constructor natives -------------------------------------------------

func methodPayload(recv Value) (*Method, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("reflect: bad Method receiver")
	}
	m, ok := in.Payload.(*Method)
	if !ok {
		return nil, fmt.Errorf("reflect: Method payload missing")
	}
	return m, nil
}

func natMethodGetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m, err := methodPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	return ctx.NewStringGo(m.Name), nil
}

func natMethodInvoke(ctx *CallContext, recv Value, args []Value) (Value, error) {
	m, err := methodPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	var callRecv Value
	if !m.Static() {
		if len(args) == 0 || args[0] == nil {
			return nil, ctx.Throw("java/lang/NullPointerException", "invoke on null receiver")
		}
		callRecv = args[0]
	}
	params, _, perr := ParseMethodDesc(m.Desc)
	if perr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", perr.Error())
	}
	var callArgs []Value
	if len(args) > 1 {
		if arr, ok := args[1].(*ArrayObj); ok {
			if len(arr.Elems) != len(params) {
				return nil, ctx.Throw("java/lang/IllegalArgumentException",
					fmt.Sprintf("wrong number of arguments: %d for %d", len(arr.Elems), len(params)))
			}
			callArgs = make([]Value, len(arr.Elems))
			for i, av := range arr.Elems {
				cv, uerr := unboxArg(av, params[i])
				if uerr != nil {
					return nil, ctx.Throw("java/lang/IllegalArgumentException",
						fmt.Sprintf("argument %d: %v", i, uerr))
				}
				callArgs[i] = cv
			}
		}
	} else if len(params) > 0 {
		return nil, ctx.Throw("java/lang/NullPointerException", "invoke args array null")
	}
	out, ierr := ctx.Invoke(m, callRecv, callArgs)
	if ierr != nil {
		return nil, ierr
	}
	retDesc := m.Desc[strings.Index(m.Desc, ")")+1:]
	if retDesc == "V" {
		return nil, nil
	}
	if len(retDesc) == 1 {
		if out == nil {
			return nil, nil
		}
		return ctx.K.boxValue(ctx, retDesc[0], out), nil
	}
	return out, nil
}

type ctorRef struct {
	cls *Class
	m   *Method
}

func ctorPayload(recv Value) (*ctorRef, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("reflect: bad Constructor receiver")
	}
	cr, ok := in.Payload.(*ctorRef)
	if !ok {
		return nil, fmt.Errorf("reflect: Constructor payload missing")
	}
	return cr, nil
}

func natConstructorGetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	cr, err := ctorPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	return ctx.NewStringGo(cr.cls.Name), nil
}

func natConstructorNewInstance(ctx *CallContext, recv Value, args []Value) (Value, error) {
	cr, err := ctorPayload(recv)
	if err != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	in, ierr := ctx.K.NewInstance(cr.cls)
	if ierr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
	}
	var callArgs []Value
	if len(args) > 0 {
		if arr, ok := args[0].(*ArrayObj); ok {
			pd, _, perr := ParseMethodDesc(cr.m.Desc)
			if perr != nil {
				return nil, ctx.Throw("java/lang/RuntimeException", perr.Error())
			}
			if len(arr.Elems) != len(pd) {
				return nil, ctx.Throw("java/lang/IllegalArgumentException", "wrong number of arguments")
			}
			callArgs = make([]Value, len(arr.Elems))
			for i, av := range arr.Elems {
				cv, uerr := unboxArg(av, pd[i])
				if uerr != nil {
					return nil, ctx.Throw("java/lang/IllegalArgumentException",
						fmt.Sprintf("argument %d: %v", i, uerr))
				}
				callArgs[i] = cv
			}
		}
	}
	if _, terr := ctx.Invoke(cr.m, in, callArgs); terr != nil {
		return nil, terr
	}
	return in, nil
}

// ---- shared helpers --------------------------------------------------------------

func newReflectInstance(k *Kernel, cls string, payload any) (*Instance, error) {
	c, ok := k.ClassByName(cls)
	if !ok {
		return nil, fmt.Errorf("reflect: %s missing", cls)
	}
	in, err := k.NewInstance(c)
	if err != nil {
		return nil, err
	}
	in.Payload = payload
	return in, nil
}

func emptyArray(k *Kernel, comp string) (Value, error) {
	return k.NewArray("L"+comp+";", 0)
}

func declaringClassOf(fieldInst Value) *Class {
	in := fieldInst.(*Instance)
	// Field instances carry the owner implicitly through the payload's Holder.
	if f, ok := in.Payload.(*Field); ok && f.Holder != nil {
		return f.Holder
	}
	panic("reflect: Field without holder")
}
