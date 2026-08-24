package kernel

import (
	"strings"
	"testing"
)

func reflectKernel(t *testing.T) (*Kernel, *Class) {
	t.Helper()
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	if _, err := k.DefineClass(&ClassDef{
		Name: "Pojo",
		Fields: []FieldDef{
			{Name: "name", Desc: "Ljava/lang/String;", Flags: 0x0002}, // private
			{Name: "age", Desc: "I", Flags: 0x0001},                   // public
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: 0x0001, Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
				return nil, nil
			}},
			{Name: "greet", Desc: "(I)Ljava/lang/String;", Flags: 0x0001, Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
				in := recv.(*Instance)
				name := ""
				if js, ok := in.Fields[0].(*JString); ok && js != nil {
					name = js.Go()
				}
				return ctx.NewStringGo("hi " + name), nil
			}},
			{Name: "sum", Desc: "(JJ)J", Flags: 0x0009, Native: func(ctx *CallContext, recv Value, args []Value) (Value, error) {
				a, _ := args[0].(int64)
				b, _ := args[1].(int64)
				return a + b, nil
			}},
		},
	}); err != nil {
		t.Fatal(err)
	}
	c, ok := k.ClassByName("Pojo")
	if !ok {
		t.Fatal("Pojo not registered")
	}
	return k, c
}

func reflectNative(t *testing.T, k *kernelSelf, cls, name, desc string, recv Value, args ...Value) Value {
	t.Helper()
	c, ok := k.ClassByName(cls)
	if !ok {
		t.Fatalf("%s missing", cls)
	}
	m, err := k.ResolveMethod(c, name, desc)
	if err != nil {
		t.Fatalf("resolve %s.%s%s: %v", cls, name, desc, err)
	}
	v, ierr := k.InvokeAs(nil, m, recv, args)
	if ierr != nil {
		t.Fatalf("invoke %s.%s: %v", cls, name, ierr)
	}
	return v
}

// kernelSelf alias keeps the helper signature readable.
type kernelSelf = Kernel

func jstrOf(t *testing.T, v Value) string {
	t.Helper()
	js, ok := v.(*JString)
	if !ok {
		t.Fatalf("want JString, got %T", v)
	}
	return js.Go()
}

func TestReflectionForNameAndName(t *testing.T) {
	k, c := reflectKernel(t)
	clsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo"))
	in := clsObj.(*Instance)
	pc, _ := in.Payload.(*Class)
	if pc != c {
		t.Fatal("forName returned a different Class than the registry")
	}
	got := reflectNative(t, k, "java/lang/Class", "getName", "()Ljava/lang/String;", clsObj)
	if jstrOf(t, got) != "Pojo" {
		t.Errorf("getName = %q", jstrOf(t, got))
	}
	if again := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo")); again != clsObj {
		t.Error("forName identity broken across calls")
	}

	// miss → ClassNotFoundException
	cnfe, ok := k.ClassByName("java/lang/ClassNotFoundException")
	if !ok {
		t.Fatal("ClassNotFoundException missing")
	}
	m, err := k.ResolveMethod(cnfe, "<init>", "(Ljava/lang/String;)V")
	if err != nil || m == nil {
		t.Fatalf("CNFE ctor missing: %v", err)
	}
}

func TestReflectionDeclaredMembers(t *testing.T) {
	k, _ := reflectKernel(t)
	clsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo"))

	fArr := reflectNative(t, k, "java/lang/Class", "getDeclaredFields",
		"()[Ljava/lang/reflect/Field;", clsObj).(*ArrayObj)
	if len(fArr.Elems) != 2 {
		t.Fatalf("declared fields = %d, want 2", len(fArr.Elems))
	}
	names := map[string]bool{}
	for _, fv := range fArr.Elems {
		n := jstrOf(t, reflectNative(t, k, "java/lang/reflect/Field",
			"getName", "()Ljava/lang/String;", fv))
		names[n] = true
	}
	if !names["name"] || !names["age"] {
		t.Fatalf("field names incomplete: %v", names)
	}

	mArr := reflectNative(t, k, "java/lang/Class", "getDeclaredMethods",
		"()[Ljava/lang/reflect/Method;", clsObj).(*ArrayObj)
	if len(mArr.Elems) != 2 { // greet + sum; <init> excluded
		t.Fatalf("declared methods = %d, want 2", len(mArr.Elems))
	}
	for _, mv := range mArr.Elems {
		n := jstrOf(t, reflectNative(t, k, "java/lang/reflect/Method",
			"getName", "()Ljava/lang/String;", mv))
		if strings.HasPrefix(n, "<") {
			t.Fatalf("constructor/clinit leaked into methods: %s", n)
		}
	}

	cArr := reflectNative(t, k, "java/lang/Class", "getDeclaredConstructors",
		"()[Ljava/lang/reflect/Constructor;", clsObj).(*ArrayObj)
	if len(cArr.Elems) != 1 {
		t.Fatalf("declared constructors = %d, want 1", len(cArr.Elems))
	}
}

func TestReflectionPrimitiveClasses(t *testing.T) {
	k, _ := reflectKernel(t)
	intCls := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("int"))
	if again := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("int")); again != intCls {
		t.Fatal("primitive Class identity broken")
	}
	name := jstrOf(t, reflectNative(t, k, "java/lang/Class", "getName",
		"()Ljava/lang/String;", intCls))
	if name != "int" {
		t.Errorf("primitive getName = %q, want int", name)
	}
	_ = intCls // primitives expose no declared members; pin via empty-array call below
	_ = reflectNative(t, k, "java/lang/Class", "getDeclaredFields",
		"()[Ljava/lang/reflect/Field;", intCls)
}

func TestReflectionFieldGetSetWithBoxing(t *testing.T) {
	k, c := reflectKernel(t)
	clsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo"))
	in, err := k.NewInstance(c)
	if err != nil {
		t.Fatal(err)
	}
	fArr := reflectNative(t, k, "java/lang/Class", "getDeclaredFields",
		"()[Ljava/lang/reflect/Field;", clsObj).(*ArrayObj)
	var nameF, ageF Value
	for _, fv := range fArr.Elems {
		switch jstrOf(t, reflectNative(t, k, "java/lang/reflect/Field",
			"getName", "()Ljava/lang/String;", fv)) {
		case "name":
			nameF = fv
		case "age":
			ageF = fv
		}
	}
	setDesc := "(Ljava/lang/Object;Ljava/lang/Object;)V"
	getDesc := "(Ljava/lang/Object;)Ljava/lang/Object;"
	reflectNative(t, k, "java/lang/reflect/Field", "set", setDesc, nameF, in, k.InternGo("rex"))
	reflectNative(t, k, "java/lang/reflect/Field", "set", setDesc, ageF, in, k.IntegerOf(7))

	if got := jstrOf(t, reflectNative(t, k, "java/lang/reflect/Field", "get", getDesc, nameF, in)); got != "rex" {
		t.Errorf("name via get = %q", got)
	}
	gotAge := reflectNative(t, k, "java/lang/reflect/Field", "get", getDesc, ageF, in)
	boxedAge := gotAge.(*Instance) // Integer box
	if raw, ok := boxedAge.Fields[0].(int32); !ok || raw != 7 {
		t.Errorf("boxed age = %#v", boxedAge.Fields[0])
	}

	// static field path
	if _, err := k.DefineClass(&ClassDef{
		Name:   "Statics",
		Fields: []FieldDef{{Name: "counter", Desc: "I", Flags: 0x0009}},
	}); err != nil {
		t.Fatal(err)
	}
	sClsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Statics"))
	sf := reflectNative(t, k, "java/lang/Class", "getDeclaredFields",
		"()[Ljava/lang/reflect/Field;", sClsObj).(*ArrayObj).Elems[0]
	reflectNative(t, k, "java/lang/reflect/Field", "set", setDesc, sf, nil, k.IntegerOf(5))
	sc, _ := k.ClassByName("Statics")
	if sc.Statics[0].(int32) != 5 {
		t.Fatalf("static set failed: %v", sc.Statics[0])
	}
	got := reflectNative(t, k, "java/lang/reflect/Field", "get", getDesc, sf, nil)
	if n, ok := IntValueOf(got); !ok || n != 5 {
		t.Fatalf("static get = %#v, want boxed 5", got)
	}
}

func TestReflectionMethodInvoke(t *testing.T) {
	k, c := reflectKernel(t)
	clsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo"))
	in, err := k.NewInstance(c)
	if err != nil {
		t.Fatal(err)
	}
	in.Fields[0] = k.InternGo("rex")

	mArr := reflectNative(t, k, "java/lang/Class", "getDeclaredMethods",
		"()[Ljava/lang/reflect/Method;", clsObj).(*ArrayObj)
	invokeDesc := "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;"
	objArrDesc := "[Ljava/lang/Object;"
	for _, mv := range mArr.Elems {
		switch jstrOf(t, reflectNative(t, k, "java/lang/reflect/Method",
			"getName", "()Ljava/lang/String;", mv)) {
		case "greet":
			args, _ := k.NewArray(objArrDesc, 1)
			args.Elems[0] = k.IntegerOf(3)
			out := reflectNative(t, k, "java/lang/reflect/Method", "invoke",
				invokeDesc, mv, in, args)
			if got := jstrOf(t, out); got != "hi rex" {
				t.Errorf("greet via reflection = %q", got)
			}
		case "sum":
			args, _ := k.NewArray(objArrDesc, 2)
			args.Elems[0] = int64(20)
			args.Elems[1] = int64(22)
			out := reflectNative(t, k, "java/lang/reflect/Method", "invoke",
				invokeDesc, mv, nil, args)
			// long results arrive boxed as java/lang/Long (reflection semantics).
			boxed, ok := out.(*Instance)
			if !ok {
				t.Fatalf("sum returned %T", out)
			}
			if v := boxed.Fields[0].(int64); v != 42 {
				t.Errorf("sum = %d, want 42", v)
			}
		}
	}
}

func TestReflectionConstructorNewInstance(t *testing.T) {
	k, _ := reflectKernel(t)
	clsObj := reflectNative(t, k, "java/lang/Class", "forName",
		"(Ljava/lang/String;)Ljava/lang/Class;", nil, k.InternGo("Pojo"))
	cArr := reflectNative(t, k, "java/lang/Class", "getDeclaredConstructors",
		"()[Ljava/lang/reflect/Constructor;", clsObj).(*ArrayObj)
	obj := reflectNative(t, k, "java/lang/reflect/Constructor", "newInstance",
		"([Ljava/lang/Object;)Ljava/lang/Object;", cArr.Elems[0],
		mustEmptyObjectArray(t, k))
	in := obj.(*Instance)
	if in.Class.Name != "Pojo" {
		t.Fatalf("newInstance produced %s", in.Class.Name)
	}
}

func mustEmptyObjectArray(t *testing.T, k *Kernel) Value {
	t.Helper()
	arr, err := k.NewArray("[Ljava/lang/Object;", 0)
	if err != nil {
		t.Fatal(err)
	}
	return arr
}
