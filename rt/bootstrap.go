package rt

// The bootstrap profile is an explicit, closed inventory. Names, descriptors
// and hierarchies here are metadata for the builder and for runtime
// installation. No application class may define anything in these
// namespaces.

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

var bootstrapInventory = []*bsClass{
	{name: "java/lang/Object", super: "", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Throwable", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getCause", "()Ljava/lang/Throwable;", AccPublic, intrinsicThrowableGetCause},
	}},
	{name: "java/lang/Exception", super: "java/lang/Throwable", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Error", super: "java/lang/Throwable", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/RuntimeException", super: "java/lang/Exception", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ReflectiveOperationException", super: "java/lang/Exception", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/NoSuchMethodException", super: "java/lang/ReflectiveOperationException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/IllegalArgumentException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/NullPointerException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ArithmeticException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/IndexOutOfBoundsException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ArrayIndexOutOfBoundsException", super: "java/lang/IndexOutOfBoundsException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ArrayStoreException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/NegativeArraySizeException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ClassCastException", super: "java/lang/RuntimeException", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/LinkageError", super: "java/lang/Error", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/IncompatibleClassChangeError", super: "java/lang/LinkageError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/NoSuchFieldError", super: "java/lang/IncompatibleClassChangeError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/NoSuchMethodError", super: "java/lang/IncompatibleClassChangeError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/InstantiationError", super: "java/lang/IncompatibleClassChangeError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/ExceptionInInitializerError", super: "java/lang/LinkageError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getCause", "()Ljava/lang/Throwable;", AccPublic, intrinsicThrowableGetCause},
	}},
	{name: "java/lang/NoClassDefFoundError", super: "java/lang/LinkageError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/AbstractMethodError", super: "java/lang/IncompatibleClassChangeError", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/StackOverflowError", super: "java/lang/Error", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Class", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
		{"getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", AccPublic, intrinsicClassGetMethod},
	}},
	{name: "java/lang/String", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Number", super: "java/lang/Object", methods: []bsMethod{
		{"<init>", "()V", AccPublic, intrinsicObjectInit},
	}},
	{name: "java/lang/Integer", super: "java/lang/Number",
		fields: []bsField{
			{name: "TYPE", descriptor: "Ljava/lang/Class;", flags: AccPublic | AccStatic | AccFinal, initial: func(r *Runtime) Value {
				return RefValue(r.IntClass().Mirror)
			}},
		},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicObjectInit},
			{"valueOf", "(I)Ljava/lang/Integer;", AccPublic | AccStatic, intrinsicIntegerValueOf},
			{"intValue", "()I", AccPublic, intrinsicIntegerIntValue},
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
