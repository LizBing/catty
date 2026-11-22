package rt

// The bootstrap profile is an explicit, closed inventory. Names, descriptors
// and hierarchies here are metadata for the builder and for runtime
// installation. No application class may define anything in these namespaces.

type bsField struct {
	name       string
	descriptor string
	flags      uint16
	initial    func(r *Runtime) Value
}

type bsMethod struct {
	name       string
	descriptor string
	flags      uint16
	intrinsic  func(t *Thread, m *Method, receiver *Object, args []Value) Outcome
}

type bsClass struct {
	name       string
	super      string
	interfaces []string
	flags      uint16
	fields     []bsField
	methods    []bsMethod
}

func typeInitial(desc byte) func(r *Runtime) Value {
	return func(r *Runtime) Value {
		return RefValue(r.PrimitiveClass(desc).Mirror)
	}
}

func throwableInit() []bsMethod {
	return []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}
}

var bootstrapInventory = []*bsClass{
	{name: "java/lang/Object", super: "", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getClass", "()Ljava/lang/Class;", AccPublic, intrinsicObjectGetClass},
	}},
	{name: "java/lang/Throwable", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"<init>", "(Ljava/lang/Throwable;)V", AccPublic, intrinsicThrowableInitCause},
		{"<init>", "(Ljava/lang/String;)V", AccPublic, intrinsicThrowableInitMessage},
		{"getCause", "()Ljava/lang/Throwable;", AccPublic, intrinsicThrowableGetCause},
	}},
	{name: "java/lang/Exception", super: "java/lang/Throwable", methods: throwableInit()},
	{name: "java/lang/Error", super: "java/lang/Throwable", methods: throwableInit()},
	{name: "java/lang/RuntimeException", super: "java/lang/Exception", methods: throwableInit()},
	{name: "java/lang/ReflectiveOperationException", super: "java/lang/Exception", methods: throwableInit()},
	{name: "java/lang/IllegalAccessException", super: "java/lang/ReflectiveOperationException", methods: throwableInit()},
	{name: "java/lang/NoSuchMethodException", super: "java/lang/ReflectiveOperationException", methods: throwableInit()},
	{name: "java/lang/IllegalArgumentException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/NullPointerException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/ArithmeticException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/IndexOutOfBoundsException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/ArrayIndexOutOfBoundsException", super: "java/lang/IndexOutOfBoundsException", methods: throwableInit()},
	{name: "java/lang/ArrayStoreException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/NegativeArraySizeException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/ClassCastException", super: "java/lang/RuntimeException", methods: throwableInit()},
	{name: "java/lang/LinkageError", super: "java/lang/Error", methods: throwableInit()},
	{name: "java/lang/IncompatibleClassChangeError", super: "java/lang/LinkageError", methods: throwableInit()},
	{name: "java/lang/IllegalAccessError", super: "java/lang/IncompatibleClassChangeError", methods: throwableInit()},
	{name: "java/lang/NoSuchFieldError", super: "java/lang/IncompatibleClassChangeError", methods: throwableInit()},
	{name: "java/lang/NoSuchMethodError", super: "java/lang/IncompatibleClassChangeError", methods: throwableInit()},
	{name: "java/lang/InstantiationError", super: "java/lang/IncompatibleClassChangeError", methods: throwableInit()},
	{name: "java/lang/ExceptionInInitializerError", super: "java/lang/LinkageError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getCause", "()Ljava/lang/Throwable;", AccPublic, intrinsicThrowableGetCause},
	}},
	{name: "java/lang/NoClassDefFoundError", super: "java/lang/LinkageError", methods: throwableInit()},
	{name: "java/lang/AbstractMethodError", super: "java/lang/IncompatibleClassChangeError", methods: throwableInit()},
	{name: "java/lang/VerifyError", super: "java/lang/LinkageError", methods: throwableInit()},
	{name: "java/lang/StackOverflowError", super: "java/lang/Error", methods: throwableInit()},
	{name: "java/lang/Class", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", AccPublic, intrinsicClassGetMethod},
	}},
	{name: "java/lang/String", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"length", "()I", AccPublic, intrinsicStringLength},
	}},
	{name: "java/lang/Number", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Boolean", super: "java/lang/Object",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('Z')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(Z)Ljava/lang/Boolean;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"booleanValue", "()Z", AccPublic, intrinsicUnbox},
		}},
	{name: "java/lang/Byte", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('B')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(B)Ljava/lang/Byte;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"byteValue", "()B", AccPublic, intrinsicUnbox},
		}},
	{name: "java/lang/Character", super: "java/lang/Object",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('C')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(C)Ljava/lang/Character;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"charValue", "()C", AccPublic, intrinsicUnbox},
		}},
	{name: "java/lang/Short", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('S')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(S)Ljava/lang/Short;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"shortValue", "()S", AccPublic, intrinsicUnbox},
		}},
	{name: "java/lang/Integer", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('I')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(I)Ljava/lang/Integer;", AccPublic | AccStatic, intrinsicIntegerValueOf},
			{"intValue", "()I", AccPublic, intrinsicIntegerIntValue},
		}},
	{name: "java/lang/Long", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('J')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(J)Ljava/lang/Long;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"longValue", "()J", AccPublic, intrinsicUnbox},
		}},
	{name: "java/lang/Float", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('F')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(F)Ljava/lang/Float;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"floatValue", "()F", AccPublic, intrinsicUnbox},
			{"intBitsToFloat", "(I)F", AccPublic | AccStatic, intrinsicFloatIntBitsToFloat},
			{"floatToRawIntBits", "(F)I", AccPublic | AccStatic, intrinsicFloatFloatToRawIntBits},
			{"floatToIntBits", "(F)I", AccPublic | AccStatic, intrinsicFloatFloatToIntBits},
		}},
	{name: "java/lang/Double", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: typeInitial('D')},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(D)Ljava/lang/Double;", AccPublic | AccStatic, intrinsicBoxValueOf},
			{"doubleValue", "()D", AccPublic, intrinsicUnbox},
			{"longBitsToDouble", "(J)D", AccPublic | AccStatic, intrinsicDoubleLongBitsToDouble},
			{"doubleToRawLongBits", "(D)J", AccPublic | AccStatic, intrinsicDoubleDoubleToRawLongBits},
			{"doubleToLongBits", "(D)J", AccPublic | AccStatic, intrinsicDoubleDoubleToLongBits},
		}},
	{name: "java/lang/System", super: "java/lang/Object",
		fields: []bsField{
			{name: "out", descriptor: "Ljava/io/PrintStream;", flags: AccPublic | AccStatic | AccFinal, initial: func(r *Runtime) Value {
				return RefValue(r.NewInstance(r.mustClass("java/io/PrintStream")))
			}},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
		}},
	{name: "java/io/PrintStream", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"println", "(I)V", AccPublic, intrinsicPrintlnInt},
		{"println", "(J)V", AccPublic, intrinsicPrintlnLong},
		{"println", "(Z)V", AccPublic, intrinsicPrintlnBoolean},
		{"println", "(C)V", AccPublic, intrinsicPrintlnChar},
		{"println", "(B)V", AccPublic, intrinsicPrintlnByte},
		{"println", "(S)V", AccPublic, intrinsicPrintlnShort},
		{"println", "(Ljava/lang/String;)V", AccPublic, intrinsicPrintlnString},
	}},
	{name: "java/lang/reflect/Method", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", AccPublic, intrinsicMethodInvoke},
	}},
	{name: "java/lang/reflect/InvocationTargetException", super: "java/lang/ReflectiveOperationException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getCause", "()Ljava/lang/Throwable;", AccPublic, intrinsicThrowableGetCause},
	}},
	{name: "java/lang/Cloneable", super: "java/lang/Object", flags: AccInterface | AccAbstract},
	{name: "java/io/Serializable", super: "java/lang/Object", flags: AccInterface | AccAbstract},
}

var (
	bootstrapByName = map[string]*bsClass{}
	bootstrapKeys   []string
)

func init() {
	for _, c := range bootstrapInventory {
		bootstrapByName[c.name] = c
		bootstrapKeys = append(bootstrapKeys, c.name)
	}
}

// BootstrapMethodInfo is the builder-facing view of one intrinsic method.
type BootstrapMethodInfo struct {
	Owner      string
	Name       string
	Descriptor string
}

// BootstrapFieldInfo is the builder-facing view of one bootstrap field.
type BootstrapFieldInfo struct {
	Owner      string
	Name       string
	Descriptor string
}

// BootstrapClassInfo is the builder-facing view of one bootstrap class.
type BootstrapClassInfo struct {
	Name       string
	Super      string
	Interfaces []string
	Flags      uint16
}

// BootstrapClassDef returns bootstrap class metadata.
func BootstrapClassDef(name string) (BootstrapClassInfo, bool) {
	c, ok := bootstrapByName[name]
	if !ok {
		return BootstrapClassInfo{}, false
	}
	return BootstrapClassInfo{Name: c.name, Super: c.super, Interfaces: c.interfaces, Flags: c.flags}, true
}

// BootstrapHasClass reports whether name belongs to the supported bootstrap
// inventory.
func BootstrapHasClass(name string) bool {
	_, ok := bootstrapByName[name]
	return ok
}

// BootstrapHasMethod reports whether owner (or a bootstrap superclass) declares
// a method with the exact name and descriptor.
func BootstrapHasMethod(owner, name, descriptor string) bool {
	c, ok := bootstrapByName[owner]
	if !ok {
		return false
	}
	for c != nil {
		for _, m := range c.methods {
			if m.name == name && m.descriptor == descriptor {
				return true
			}
		}
		if c.super == "" {
			break
		}
		c = bootstrapByName[c.super]
	}
	return false
}

// BootstrapHasField reports whether owner or a bootstrap superclass declares a
// field with the exact name and descriptor.
func BootstrapHasField(owner, name, descriptor string) bool {
	c, ok := bootstrapByName[owner]
	if !ok {
		return false
	}
	for c != nil {
		for _, f := range c.fields {
			if f.name == name && f.descriptor == descriptor {
				return true
			}
		}
		if c.super == "" {
			break
		}
		c = bootstrapByName[c.super]
	}
	return false
}

// BootstrapMethods lists every declared intrinsic method in deterministic
// inventory order. The builder includes these records in its report.
func BootstrapMethods() []BootstrapMethodInfo {
	var out []BootstrapMethodInfo
	for _, c := range bootstrapInventory {
		for _, m := range c.methods {
			out = append(out, BootstrapMethodInfo{Owner: c.name, Name: m.name, Descriptor: m.descriptor})
		}
	}
	return out
}

// BootstrapFields lists every declared bootstrap field in deterministic
// inventory order.
func BootstrapFields() []BootstrapFieldInfo {
	var out []BootstrapFieldInfo
	for _, c := range bootstrapInventory {
		for _, f := range c.fields {
			out = append(out, BootstrapFieldInfo{Owner: c.name, Name: f.name, Descriptor: f.descriptor})
		}
	}
	return out
}

// BootstrapClasses lists the bootstrap classes in inventory order.
func BootstrapClasses() []BootstrapClassInfo {
	out := make([]BootstrapClassInfo, 0, len(bootstrapInventory))
	for _, c := range bootstrapInventory {
		out = append(out, BootstrapClassInfo{Name: c.name, Super: c.super, Interfaces: c.interfaces, Flags: c.flags})
	}
	return out
}
