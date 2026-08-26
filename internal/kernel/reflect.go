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
	// The mirror INSTANCE is an instance of java/lang/Class itself (so
	// gson's `type instanceof Class` works), carrying a primitiveInfo
	// payload that getName/isPrimitive read. The pseudo-Class named after
	// the keyword exists only so methodsByKey resolution has a home.
	cls := &Class{Name: name, Flags: classfileAccPublicFinal}
	cls.Methods = append(cls.Methods, k.primObjMethods...)
	cls.Methods = append(cls.Methods, k.primClsMethods...)
	cls.methodsByKey = make(map[string]*Method, len(cls.Methods))
	for _, m := range cls.Methods {
		cls.methodsByKey[memberKey(m.Name, m.Desc)] = m
	}
	cls.setState(StateInitialized)
	in := &Instance{}
	in.Class = k.classSelfCls
	in.Payload = &primitiveInfo{name: name, desc: desc, metaCls: cls}
	actual, _ := k.primOnce.LoadOrStore(desc, in)
	return actual.(*Instance), nil
}

type primitiveInfo struct {
	name    string
	desc    string
	metaCls *Class // pseudo-Class named "int" etc. for method resolution
}

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
		return k.boxViaWrapper("java/lang/Boolean", int64(v.(int32)))
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
	// Map class name to its primitive descriptor for valueOf lookup.
	primDescs := map[string]string{
		"java/lang/Byte":      "B",
		"java/lang/Short":     "S",
		"java/lang/Character": "C",
		"java/lang/Boolean":   "Z",
		"java/lang/Long":      "J",
	}
	pd, ok := primDescs[cls]
	if !ok {
		pd = "J"
	}
	m, err := k.ResolveMethod(mustLookup(k, cls), "valueOf", "("+pd+")"+wrapperDesc(cls))
	if err != nil {
		return k.IntegerOf(int32(v)) // fallback
	}
	var arg Value
	switch pd {
	case "Z":
		arg = int32(v)
	default:
		arg = v
	}
	out, ierr := k.InvokeAs(nil, m, nil, []Value{arg})
	if ierr != nil {
		return k.IntegerOf(int32(v))
	}
	return out
}

func wrapperDesc(cls string) string { return "L" + cls + ";" }

func mustLookup(k *Kernel, name string) *Class {
	c, _ := k.ClassByName(name)
	return c
}

// ---- registration ----------------------------------------------------------

func registerReflection(k *Kernel) error {
	// java.lang.reflect.Type — marker interface for gson's TypeToken world.
	// Defined BEFORE java/lang/Class so the latter can declare it as an
	// interface (Class implements Type per JVM semantics).
	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/reflect/Type",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	}); err != nil {
		return err
	}
	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/annotation/Annotation",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	}); err != nil {
		return err
	}
	// GenericArrayType / ParameterizedType / WildcardType / TypeVariable —
	// marker interfaces for the generic-type family (gson TypeToken).
	for _, nm := range []string{
		"java/lang/reflect/GenericArrayType",
		"java/lang/reflect/WildcardType",
		"java/lang/reflect/TypeVariable",
	} {
		if _, err := k.DefineClass(&ClassDef{
			Name:   nm,
			Super:  "java/lang/Object",
			Ifaces: []string{"java/lang/reflect/Type"},
			Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
		}); err != nil {
			return err
		}
	}

	// ClassNotFoundException next: forName throws it.
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
		Name:   "java/lang/Class",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/lang/reflect/Type"}, // Class implements Type (JVM)
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natClassGetName},
			{Name: "getSimpleName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natClassGetSimpleName},
			{Name: "isArray", Desc: "()Z", Flags: 0x0001, Native: natClassIsArray},
			{Name: "getModifiers", Desc: "()I", Flags: 0x0001, Native: natClassGetModifiers},
			{Name: "isAnonymousClass", Desc: "()Z", Flags: 0x0001, Native: natClassIsAnonymousClass},
			{Name: "isLocalClass", Desc: "()Z", Flags: 0x0001, Native: natClassIsLocalClass},
			{Name: "isMemberClass", Desc: "()Z", Flags: 0x0001, Native: natClassIsMemberClass},
			{Name: "getAnnotation", Desc: "(Ljava/lang/Class;)Ljava/lang/annotation/Annotation;", Flags: 0x0001,
				Native: natReflectGetAnnotation},
			{Name: "getAnnotations", Desc: "()[Ljava/lang/annotation/Annotation;", Flags: 0x0001,
				Native: natReflectGetAnnotations},
			{Name: "isEnum", Desc: "()Z", Flags: 0x0001, Native: natClassIsEnum},
			{Name: "isInterface", Desc: "()Z", Flags: 0x0001, Native: natClassIsInterface},
			{Name: "isPrimitive", Desc: "()Z", Flags: 0x0001, Native: natClassIsPrimitive},
			{Name: "getInterfaces", Desc: "()[Ljava/lang/Class;", Flags: 0x0001, Native: natClassGetInterfaces},
			{Name: "getGenericInterfaces", Desc: "()[Ljava/lang/reflect/Type;", Flags: 0x0001, Native: natClassGetInterfaces},
			{Name: "forName", Desc: "(Ljava/lang/String;)Ljava/lang/Class;", Flags: 0x0009, Native: natClassForName},
			{Name: "getDeclaredFields", Desc: "()[Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassDeclaredFields},
			{Name: "getFields", Desc: "()[Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassGetFields},
			{Name: "getMethods", Desc: "()[Ljava/lang/reflect/Method;", Flags: 0x0001, Native: natClassGetMethods},
			{Name: "getMethod", Desc: "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", Flags: 0x0001, Native: natClassGetMethod},
			{Name: "getConstructors", Desc: "()[Ljava/lang/reflect/Constructor;", Flags: 0x0001, Native: natClassGetConstructors},
			{Name: "getField", Desc: "(Ljava/lang/String;)Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassGetField},
			{Name: "getDeclaredField", Desc: "(Ljava/lang/String;)Ljava/lang/reflect/Field;", Flags: 0x0001, Native: natClassGetDeclaredField},
			{Name: "getDeclaredMethods", Desc: "()[Ljava/lang/reflect/Method;", Flags: 0x0001, Native: natClassDeclaredMethods},
			{Name: "getDeclaredConstructors", Desc: "()[Ljava/lang/reflect/Constructor;", Flags: 0x0001, Native: natClassDeclaredConstructors},
			{Name: "getDeclaredConstructor", Desc: "([Ljava/lang/Class;)Ljava/lang/reflect/Constructor;", Flags: 0x0001, Native: natClassGetDeclaredConstructor},
			{Name: "getConstructor", Desc: "([Ljava/lang/Class;)Ljava/lang/reflect/Constructor;", Flags: 0x0001, Native: natClassGetDeclaredConstructor},
			{Name: "isInstance", Desc: "(Ljava/lang/Object;)Z", Flags: 0x0001, Native: natClassIsInstance},
			{Name: "isAssignableFrom", Desc: "(Ljava/lang/Class;)Z", Flags: 0x0001, Native: natClassIsAssignableFrom},
			{Name: "getSuperclass", Desc: "()Ljava/lang/Class;", Flags: 0x0001, Native: natClassGetSuperclass},
			{Name: "getGenericSuperclass", Desc: "()Ljava/lang/reflect/Type;", Flags: 0x0001, Native: natClassGetSuperclass},
			{Name: "newInstance", Desc: "()Ljava/lang/Object;", Flags: 0x0001, Native: natClassNewInstance},
			{Name: "cast", Desc: "(Ljava/lang/Object;)Ljava/lang/Object;", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					in, ok := recv.(*Instance)
					if !ok {
						return nil, ctx.Throw("java/lang/RuntimeException", "bad Class recv")
					}
					c, ok2 := in.Payload.(*Class)
					if !ok2 {
						return args[0], nil
					} // primitives: pass through
					if len(args) > 0 && args[0] != nil && !ctx.K.IsInstance(args[0], c) {
						return nil, ctx.Throw("java/lang/ClassCastException",
							in.Class.Name+" cannot be cast to class "+c.Name)
					}
					if len(args) > 0 {
						return args[0], nil
					}
					return nil, nil
				}},
		},
	}); err != nil {
		return err
	}

	// Snapshot method tables for primitive mirrors — AFTER java/lang/Class
	// exists (its Methods list is what int.class etc. share).
	k.primObjMethods = mustLookup(k, "java/lang/Object").Methods
	k.primClsMethods = mustLookup(k, "java/lang/Class").Methods
	k.classSelfCls = mustLookup(k, "java/lang/Class")

	if _, err := k.DefineClass(&ClassDef{
		Name:  "java/lang/reflect/Field",
		Super: "java/lang/Object",
		Methods: []MethodDef{
			{Name: "getName", Desc: "()Ljava/lang/String;", Flags: 0x0001, Native: natFieldGetName},
			{Name: "setAccessible", Desc: "(Z)V", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return nil, nil
				}},
			{Name: "isAccessible", Desc: "()Z", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return int32(1), nil
				}},
			{Name: "getType", Desc: "()Ljava/lang/Class;", Flags: 0x0001, Native: natFieldGetType},
			{Name: "getGenericType", Desc: "()Ljava/lang/reflect/Type;", Flags: 0x0001, Native: natFieldGetType}, // v1: generic == declared
			{Name: "getModifiers", Desc: "()I", Flags: 0x0001, Native: natFieldGetModifiers},
			{Name: "isSynthetic", Desc: "()Z", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					f, ferr := fieldPayload(recv)
					if ferr != nil {
						return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
					}
					return boolV(f.Flags&0x1000 != 0), nil // ACC_SYNTHETIC
				}},
			{Name: "getAnnotation", Desc: "(Ljava/lang/Class;)Ljava/lang/annotation/Annotation;", Flags: 0x0001,
				Native: natReflectGetAnnotation},
			{Name: "getAnnotations", Desc: "()[Ljava/lang/annotation/Annotation;", Flags: 0x0001,
				Native: natReflectGetAnnotations},
			{Name: "isEnumConstant", Desc: "()Z", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					f, ferr := fieldPayload(recv)
					if ferr != nil {
						return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
					}
					return boolV(f.Flags&0x4000 != 0), nil // ACC_ENUM
				}},
			{Name: "getDeclaringClass", Desc: "()Ljava/lang/Class;", Flags: 0x0001, Native: natFieldGetDeclaringClass},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: 0x0001, Native: natReflectEquals},
			{Name: "hashCode", Desc: "()I", Flags: 0x0001, Native: natReflectHashCode},
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
			{Name: "isAccessible", Desc: "()Z", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return int32(1), nil // access control not enforced (DEV-0010)
				}},
			{Name: "setAccessible", Desc: "(Z)V", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return nil, nil // no-op: everything accessible
				}},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: 0x0001, Native: natReflectEquals},
			{Name: "hashCode", Desc: "()I", Flags: 0x0001, Native: natReflectHashCode},
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
			{Name: "isAccessible", Desc: "()Z", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return int32(1), nil
				}},
			{Name: "setAccessible", Desc: "(Z)V", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return nil, nil
				}},
			{Name: "equals", Desc: "(Ljava/lang/Object;)Z", Flags: 0x0001, Native: natReflectEquals},
			{Name: "hashCode", Desc: "()I", Flags: 0x0001, Native: natReflectHashCode},
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

// Anonymous/local/member detection via the binary-name conventions javac
// emits: Outer$1 (anonymous), Outer$1Local (local), Outer$Inner (member).
func anonLocalMember(name string) (anon, local, member bool) {
	i := strings.LastIndex(name, "$")
	if i < 0 {
		return false, false, strings.Contains(name, ".") || name != ""
	}
	tail := name[i+1:]
	if tail == "" {
		return false, false, true
	}
	c := tail[0]
	if c >= '0' && c <= '9' {
		return true, false, true
	}
	// $1Local style: digit followed by identifier → local
	for j := 1; j < len(tail); j++ {
		if tail[j] >= '0' && tail[j] <= '9' {
			continue
		}
		if j > 1 || !(tail[j] >= 'A' && tail[j] <= 'Z') {
			break
		}
	}
	if len(tail) > 1 && tail[0] >= '0' && tail[0] <= '9' {
		rest := tail[1:]
		isDigits := true
		for k := range rest {
			if rest[k] < '0' || rest[k] > '9' {
				isDigits = false
				break
			}
		}
		if isDigits {
			return true, false, true // pure numeric: anonymous
		}
		return false, true, true // local class: digits + name
	}
	return false, false, true
}

func natClassIsAnonymousClass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 {
		return int32(0), nil
	}
	a, _, _ := anonLocalMember(c.Name)
	return boolV(a), nil
}

func natClassIsLocalClass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 {
		return int32(0), nil
	}
	_, l, _ := anonLocalMember(c.Name)
	return boolV(l), nil
}

func natClassIsMemberClass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 {
		return int32(0), nil
	}
	m, _, _ := anonLocalMember(c.Name)
	return boolV(m && !c.IsArray), nil
}

func natClassGetModifiers(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(1), nil // public fallback
	}
	if c, ok2 := in.Payload.(*Class); ok2 {
		return int32(c.Flags), nil
	}
	return int32(0x0011), nil // public+final for primitives
}

func natClassIsPrimitive(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	_, ok2 := in.Payload.(*primitiveInfo)
	return boolV(ok2), nil
}

func natClassIsInterface(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 {
		return int32(0), nil
	}
	return boolV(c.Flags&classfile.AccInterface != 0), nil
}

func natClassGetInterfaces(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return emptyArray(ctx.K, "Ljava/lang/Class;")
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 {
		return emptyArray(ctx.K, "Ljava/lang/Class;")
	}
	arr, aerr := ctx.K.NewArray("Ljava/lang/Class;", len(c.Ifaces))
	if aerr != nil {
		return nil, aerr
	}
	for i, ic := range c.Ifaces {
		mirror, err := ctx.K.ClassObjectOf(ic)
		if err != nil {
			return nil, err
		}
		arr.Elems[i] = mirror
	}
	return arr, nil
}

func natClassIsArray(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok3 := in.Payload.(*Class)
	if !ok3 {
		return int32(0), nil
	}
	if c.IsArray {
		return int32(1), nil
	}
	return int32(0), nil
}

func natClassIsEnum(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	c, ok2 := in.Payload.(*Class)
	if !ok2 || c.Super == nil {
		return int32(0), nil
	}
	if c.Super.Name == "java/lang/Enum" {
		return int32(1), nil
	}
	return int32(0), nil
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

// walkPublicFields collects PUBLIC fields along the superclass chain,
// nearest declaration winning on name shadowing (P-0012 inherited
// traversal).
func walkPublicFields(c *Class) []*Field {
	var out []*Field
	// NOTE: unlike method overriding, a shadowed field is a DISTINCT
	// field in Java — getFields reports BOTH child.base and Base.base.
	for x := c; x != nil; x = x.Super {
		for _, f := range x.DeclaredFields {
			if f.Static || f.Flags&0x0001 == 0 {
				continue
			}
			out = append(out, f)
		}
	}
	return out
}

// walkPublicMethods collects PUBLIC non-constructor methods along the
// superclass chain, nearest declaration winning on name+descriptor.
// Interface default/abstract methods are not included (v1 limitation,
// DEV-0010).
func walkPublicMethods(c *Class) []*Method {
	seen := map[string]bool{}
	var out []*Method
	for x := c; x != nil; x = x.Super {
		for _, m := range x.Methods {
			if m.Name == "<init>" || m.Name == "<clinit>" || m.Flags&0x0001 == 0 {
				continue
			}
			key := memberKey(m.Name, m.Desc)
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, m)
		}
	}
	return out
}

func natClassGetFields(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return emptyArray(ctx.K, "Ljava/lang/reflect/Field;")
	}
	fs := walkPublicFields(c)
	arr, aerr := ctx.K.NewArray("Ljava/lang/reflect/Field;", len(fs))
	if aerr != nil {
		return nil, aerr
	}
	for i, f := range fs {
		in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Field", f)
		if ierr != nil {
			return nil, ierr
		}
		arr.Elems[i] = in
	}
	return arr, nil
}

func natClassGetMethods(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		if _, ok := err.(*primitiveSignal); ok {
			return emptyArray(ctx.K, "Ljava/lang/reflect/Method;")
		}
		return nil, ctx.Throw("java/lang/RuntimeException", err.Error())
	}
	ms := walkPublicMethods(c)
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

// natClassGetMethod resolves a public method by name + parameter Classes
// along the superclass chain (reflection v2).
func natClassGetMethod(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, cerr := classPayloadOrThrow(recv)
	if cerr != nil {
		return nil, ctx.Throw("java/lang/NoSuchMethodException", "primitive class")
	}
	js, _ := args[0].(*JString)
	if js == nil {
		return nil, ctx.Throw("java/lang/NullPointerException", "getMethod(null)")
	}
	var want []*Instance
	if len(args) > 1 {
		if arr, ok := args[1].(*ArrayObj); ok {
			for _, e := range arr.Elems {
				want = append(want, e.(*Instance))
			}
		}
	}
	for x := c; x != nil; x = x.Super {
		for _, m := range x.Methods {
			if m.Name != js.Go() || m.Flags&0x0001 == 0 ||
				strings.HasPrefix(m.Name, "<") {
				continue
			}
			pd, _, perr := ParseMethodDesc(m.Desc)
			if perr != nil || len(pd) != len(want) {
				continue
			}
			match := true
			for i, w := range want {
				dc, derr := ctx.K.descToClass(pd[i])
				if derr != nil || dc != w {
					match = false
					break
				}
			}
			if !match {
				continue
			}
			in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Method", m)
			if ierr != nil {
				return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
			}
			return in, nil
		}
	}
	sig := js.Go() + "/"
	for _, w := range want {
		if pi, ok := w.Payload.(*primitiveInfo); ok {
			sig += pi.desc
		} else if cc, ok2 := w.Payload.(*Class); ok2 {
			sig += "L" + cc.Name + ";"
		} else {
			sig += "?"
		}
	}
	return nil, ctx.Throw("java/lang/NoSuchMethodException", sig)
}

func natClassGetConstructors(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, err := classPayloadOrThrow(recv)
	if err != nil {
		return emptyArray(ctx.K, "Ljava/lang/reflect/Constructor;")
	}
	var cs []*Method
	for _, m := range c.Methods {
		if m.Name == "<init>" && m.Flags&0x0001 != 0 {
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

func natClassGetDeclaredField(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, cerr := classPayloadOrThrow(recv)
	if cerr != nil {
		return nil, ctx.Throw("java/lang/NoSuchFieldException", "primitive")
	}
	js, _ := args[0].(*JString)
	if js == nil {
		return nil, ctx.Throw("java/lang/NoSuchFieldException", "null")
	}
	for _, f := range c.DeclaredFields {
		if f.Name == js.Go() {
			in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Field", f)
			if ierr != nil {
				return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
			}
			return in, nil
		}
	}
	return nil, ctx.Throw("java/lang/NoSuchFieldException", js.Go())
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

func matchParams(m *Method, want []*Instance, k *Kernel) bool {
	pd, _, err := ParseMethodDesc(m.Desc)
	if err != nil || len(pd) != len(want) {
		return false
	}
	for i, w := range want {
		dc, derr := k.descToClass(pd[i])
		if derr != nil || dc != w {
			return false
		}
	}
	return true
}

func findCtor(k *Kernel, c *Class, want []*Instance) *Method {
	for _, m := range c.Methods {
		if m.Name != "<init>" {
			continue
		}
		if len(want) == 0 && len(m.Desc) > 3 && m.Desc[1] == ')' {
			return m // ()V-style no-arg
		}
		if matchParams(m, want, k) {
			return m
		}
	}
	return nil
}

func ctorArrayArg(args []Value) []*Instance {
	var want []*Instance
	if len(args) > 0 {
		if arr, ok := args[0].(*ArrayObj); ok {
			for _, e := range arr.Elems {
				if in, ok2 := e.(*Instance); ok2 {
					want = append(want, in)
				}
			}
		}
	}
	return want
}

func natClassGetDeclaredConstructor(ctx *CallContext, recv Value, args []Value) (Value, error) {
	c, cerr := classPayloadOrThrow(recv)
	if cerr != nil {
		return nil, ctx.Throw("java/lang/NoSuchMethodException", "primitive")
	}
	want := ctorArrayArg(args)
	m := findCtor(kFor(ctx), c, want)
	if m == nil {
		return nil, ctx.Throw("java/lang/NoSuchMethodException", c.Name)
	}
	in, ierr := newReflectInstance(ctx.K, "java/lang/reflect/Constructor",
		&ctorRef{cls: c, m: m})
	if ierr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ierr.Error())
	}
	return in, nil
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

func natClassIsAssignableFrom(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok || len(args) == 0 || args[0] == nil {
		return int32(0), nil
	}
	c, ok := in.Payload.(*Class)
	if !ok {
		return int32(0), nil
	}
	argCls, ok2 := args[0].(*Instance)
	if !ok2 {
		return int32(0), nil
	}
	ac, ok3 := argCls.Payload.(*Class)
	if !ok3 {
		return int32(0), nil
	}
	if c == ac || ctx.K.classIsa(ac, c) {
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
	if pt := ctx.K.buildParameterizedSuper(c); pt != nil {
		return pt, nil
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

func natFieldGetModifiers(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, ferr := fieldPayload(recv)
	if ferr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
	}
	return int32(f.Flags), nil
}

func natFieldGetDeclaringClass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	f, ferr := fieldPayload(recv)
	if ferr != nil {
		return nil, ctx.Throw("java/lang/RuntimeException", ferr.Error())
	}
	if f.Holder == nil {
		return nil, ctx.Throw("java/lang/RuntimeException", "field without holder")
	}
	return ctx.K.ClassObjectOf(f.Holder)
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

// natReflectEquals/natReflectHashCode back equals()/hashCode() on all
// reflective mirror objects (Field/Method/Constructor): two mirrors are
// equal iff they wrap the SAME runtime member — matching java.lang
// reflection semantics where getMethod twice yields equal objects.
func natReflectEquals(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	if len(args) == 0 || args[0] == nil {
		return int32(0), nil
	}
	other, ok := args[0].(*Instance)
	if !ok || other.Class != in.Class {
		return int32(0), nil
	}
	if in.Payload == nil || other.Payload == nil {
		return int32(0), nil
	}
	if in.Payload == other.Payload {
		return int32(1), nil
	}
	return int32(0), nil
}

func natReflectHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return int32(0), nil
	}
	return int32(in.IdentityHash()), nil
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

// natEnumValueOf implements Enum.valueOf(Class,String) by reading the
// enum class's synthetic static $VALUES array directly from the registry —
// no Java-level reflection API needed (the kernel owns the storage).
func natEnumValueOf(ctx *CallContext, recv Value, args []Value) (Value, error) {
	cInst, _ := args[0].(*Instance)
	js, _ := args[1].(*JString)
	if cInst == nil || js == nil {
		return nil, ctx.Throw("java/lang/NullPointerException", "Enum.valueOf")
	}
	c, ok := cInst.Payload.(*Class)
	if !ok {
		return nil, ctx.Throw("java/lang/NullPointerException", "Enum.valueOf: not an enum class")
	}
	f := c.fieldsByKey[memberKey("$VALUES", "[L"+c.Name+";")]
	if f == nil {
		return nil, ctx.Throw("java/lang/IllegalArgumentException",
			"No enum constants in "+c.Name)
	}
	arr, ok := c.Statics[f.StaticSlot].(*ArrayObj)
	if !ok {
		return nil, ctx.Throw("java/lang/IllegalArgumentException",
			"No enum constants in "+c.Name)
	}
	for _, e := range arr.Elems {
		in, ok := e.(*Instance)
		if !ok || in.Class != c {
			continue
		}
		name := ""
		if jsn, ok := in.Fields[0].(*JString); ok && jsn != nil {
			name = jsn.Go()
		}
		if name == js.Go() {
			return e, nil
		}
	}
	return nil, ctx.Throw("java/lang/IllegalArgumentException",
		"No enum constant "+c.Name+"."+js.Go())
}

func kFor(ctx *CallContext) *Kernel { return ctx.K }

// ---- annotation metadata access (P-0013 v3) -------------------------------

type ElementValue = classfile.ElementValue

func annDataFrom(recv Value) []ParsedAnnotation {
	in, ok := recv.(*Instance)
	if !ok || in.Payload == nil {
		return nil
	}
	switch p := in.Payload.(type) {
	case *Field:
		return p.Annotations
	case *Method:
		return p.Annotations
	}
	return nil
}

// natReflectGetAnnotation returns the first annotation whose type matches
// the wanted Class mirror. Works for Class, Field, and Method receivers.
func natReflectGetAnnotation(ctx *CallContext, recv Value, args []Value) (Value, error) {
	anns := annDataFrom(recv)
	if len(anns) == 0 || len(args) == 0 {
		return nil, nil
	}
	wantIn, ok := args[0].(*Instance)
	if !ok {
		return nil, nil
	}
	wantC, hasPayload := classPayloadOf(wantIn)
	if !hasPayload {
		return nil, nil
	}
	target := "L" + wantC.Name + ";"
	for i := range anns {
		if anns[i].TypeDesc != target {
			continue
		}
		inst, berr := buildAnnotationMirror(ctx.K, &anns[i])
		if berr != nil {
			return nil, ctx.Throw("java/lang/RuntimeException", berr.Error())
		}
		return inst, nil
	}
	return nil, nil
}

func natReflectGetAnnotations(ctx *CallContext, recv Value, args []Value) (Value, error) {
	anns := annDataFrom(recv)
	arr, aerr := ctx.K.NewArray("Ljava/lang/annotation/Annotation;", len(anns))
	if aerr != nil {
		return nil, aerr
	}
	for i := range anns {
		inst, berr := buildAnnotationMirror(ctx.K, &anns[i])
		if berr != nil {
			continue
		}
		arr.Elems[i] = inst
	}
	return arr, nil
}

// buildAnnotationMirror creates a Java object implementing the annotation
// interface by registering a synthetic concrete companion class and creating
// an instance of it. Avoids DefineClass entirely for maximum simplicity.
func buildAnnotationMirror(k *Kernel, pa *ParsedAnnotation) (*Instance, error) {
	internalName := strings.TrimSuffix(strings.TrimPrefix(pa.TypeDesc, "L"), ";")
	implName := internalName + "$CattyImpl"

	// Check if already registered.
	if c, ok := k.ClassByName(implName); ok {
		in, err := k.NewInstance(c)
		if err != nil {
			return nil, err
		}
		pl := map[string]ElementValue{}
		for k2, v2 := range pa.Elements {
			pl[k2] = v2
		}
		in.Payload = pl
		return in, nil
	}

	iface, ok := k.ClassByName(internalName)
	if !ok {
		return nil, fmt.Errorf("annotation class %s not loaded", internalName)
	}

	// Build the synthetic class directly.
	c := &Class{
		Name:  implName,
		Super: iface,
		Flags: 0x0001, // public
	}
	c.setState(StateInitialized)

	methodsByKey := map[string]*Method{}
	for _, m := range iface.Methods {
		elemName := m.Name
		cm := &Method{
			Holder: c, Name: elemName, Desc: m.Desc,
			Flags: 0x0001,
			Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
				in, ok3 := recv.(*Instance)
				if !ok3 {
					return nil, nil
				}
				pl, ok4 := in.Payload.(map[string]ElementValue)
				if !ok4 {
					return nil, nil
				}
				ev, exists := pl[elemName]
				if !exists {
					// Element has a default value or is absent — return
					// type-appropriate zero. Check the method's descriptor.
					if strings.HasSuffix(m.Desc, "[Ljava/lang/String;") {
						arr, _ := ctx.K.NewArray("Ljava/lang/String;", 0)
						return arr, nil
					}
					if strings.HasSuffix(m.Desc, "Z") {
						return int32(0), nil
					}
					if strings.HasSuffix(m.Desc, "I") {
						return int32(0), nil
					}
					return nil, nil
				}
				switch ev.Tag {
				case 's':
					return ctx.NewStringGo(ev.EnumName), nil
				case 'I', 'B', 'S':
					return int32(int32(ev.ConstIdx)), nil
				case 'Z':
					if ev.ConstIdx != 0 {
						return int32(1), nil
					}
					return int32(0), nil
				default:
					return nil, nil
				}
			},
		}
		c.Methods = append(c.Methods, cm)
		methodsByKey[memberKey(elemName, m.Desc)] = cm
	}
	c.methodsByKey = methodsByKey

	// Register in the kernel so ResolveMethod/findMethod can find it.
	k.RegisterSynthetic(c)

	in, ierr := k.NewInstance(c)
	if ierr != nil {
		return nil, ierr
	}
	pl := map[string]ElementValue{}
	for k2, v2 := range pa.Elements {
		pl[k2] = v2
	}
	in.Payload = pl
	return in, nil
}

func parseConstIntVal(idx uint16) int32 { return int32(idx) }

func classPayloadOf(in *Instance) (*Class, bool) {
	c, ok := in.Payload.(*Class)
	return c, ok
}

var _ = fmt.Sprintf // suppress unused import if needed

func idx32(v uint16) int32 { return int32(v) }

// registerParameterizedType registers ParameterizedType interface and
// ParameterizedTypeImpl concrete class for generic type support (P-0012).
func registerParameterizedType(k *Kernel) {
	mustDefine(k, &ClassDef{
		Name:   "java/lang/reflect/ParameterizedType",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/lang/reflect/Type"},
		Flags:  classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})
	mustDefine(k, &ClassDef{
		Name:   "java/lang/reflect/ParameterizedTypeImpl",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/lang/reflect/ParameterizedType"},
		Fields: []FieldDef{
			{Name: "rawType", Desc: "Ljava/lang/Class;", Flags: 0x0002},
			{Name: "actualArgs", Desc: "[Ljava/lang/reflect/Type;", Flags: 0x0002},
			{Name: "ownerType", Desc: "Ljava/lang/reflect/Type;", Flags: 0x0002},
		},
		Methods: []MethodDef{
			{Name: "getRawType", Desc: "()Ljava/lang/Class;", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return recv.(*Instance).Fields[0], nil
				}},
			{Name: "getActualTypeArguments", Desc: "()[Ljava/lang/reflect/Type;", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					f := recv.(*Instance).Class.FindField("actualArgs", "[Ljava/lang/reflect/Type;")
					if f != nil && f.Slot < len(recv.(*Instance).Fields) && recv.(*Instance).Fields[f.Slot] != nil {
						return recv.(*Instance).Fields[f.Slot], nil
					}
					return ctx.K.NewArray("Ljava/lang/reflect/Type;", 0)
				}},
			{Name: "getOwnerType", Desc: "()Ljava/lang/reflect/Type;", Flags: 0x0001,
				Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
					return nil, nil
				}},
		},
	})
}

// buildParameterizedSuper creates a ParameterizedType instance when the
// class's Signature encodes a parameterized superclass (P-0012).
func (k *Kernel) buildParameterizedSuper(c *Class) *Instance {
	sig := c.Signature
	if sig == "" || !strings.Contains(sig, "<") {
		return nil
	}
	if strings.HasPrefix(sig, "<") {
		end := strings.Index(sig, ">")
		if end < 0 {
			return nil
		}
		sig = sig[end+1:]
	}
	bracket := strings.Index(sig, "<")
	if bracket <= 1 {
		return nil
	}
	superInternal := strings.TrimPrefix(sig[:bracket], "L")
	superCls := k.lookupClass(superInternal)
	if superCls == nil {
		return nil
	}
	argsStart := bracket + 1
	argsEnd := strings.LastIndex(sig, ">")
	if argsEnd <= argsStart {
		return nil
	}
	argStr := sig[argsStart:argsEnd]

	var typeArgs []Value
	for len(argStr) > 0 {
		if argStr[0] != 'L' {
			break
		}
		end := 1
		depth := 0
		for end < len(argStr) {
			ch := argStr[end]
			if ch == '<' {
				depth++
			} else if ch == '>' {
				depth--
			} else if ch == ';' && depth == 0 {
				break
			}
			end++
		}
		internal := argStr[1:end]
		if i := strings.Index(internal, "<"); i >= 0 {
			internal = internal[:i]
		}
		tc, tcErr := k.ResolveClass(internal)
		if tcErr != nil {
			return nil
		}
		mirror, mErr := k.ClassObjectOf(tc)
		if mErr != nil {
			return nil
		}
		typeArgs = append(typeArgs, mirror)
		argStr = argStr[end+1:]
	}

	ptCls, ok := k.ClassByName("java/lang/reflect/ParameterizedTypeImpl")
	if !ok {
		return nil
	}
	in, ierr := k.NewInstance(ptCls)
	if ierr != nil {
		return nil
	}
	if f := ptCls.FindField("rawType", "Ljava/lang/Class;"); f != nil && f.Slot < len(in.Fields) {
		obj, _ := k.ClassObjectOf(superCls)
		in.Fields[f.Slot] = obj
	}
	if f := ptCls.FindField("actualArgs", "[Ljava/lang/reflect/Type;"); f != nil && f.Slot < len(in.Fields) {
		arr, _ := k.NewArray("Ljava/lang/reflect/Type;", len(typeArgs))
		copy(arr.Elems, typeArgs)
		in.Fields[f.Slot] = arr
	}
	return in
}
