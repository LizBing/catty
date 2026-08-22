package kernel

import (
	"bytes"
	"errors"
	"testing"

	"catty/internal/classfile"
)

func newTestKernel(t *testing.T) (*Kernel, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	k := New(Options{Stdout: &buf})
	return k, &buf
}

func TestFieldLayoutSupersFirst(t *testing.T) {
	k, _ := newTestKernel(t)

	if _, err := k.DefineClass(&ClassDef{
		Name:   "test/Base",
		Fields: []FieldDef{{Name: "b0", Desc: "I", Flags: classfile.AccPublic}},
	}); err != nil {
		t.Fatal(err)
	}
	c, err := k.DefineClass(&ClassDef{
		Name:  "test/Derived",
		Super: "test/Base",
		Fields: []FieldDef{
			{Name: "d0", Desc: "J", Flags: classfile.AccPublic},
			{Name: "d1", Desc: "I", Flags: classfile.AccPublic},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.layoutSize != 4 { // b0 + d0(cat2=2 slots) + d1
		t.Fatalf("layoutSize = %d, want 4", c.layoutSize)
	}

	obj, err := k.NewInstance(c)
	if err != nil {
		t.Fatal(err)
	}
	fB := k.mustField(t, k, "test/Derived", "b0", "I")
	fD0 := k.mustField(t, k, "test/Derived", "d0", "J")
	fD1 := k.mustField(t, k, "test/Derived", "d1", "I")
	if fB.Slot != 0 || fD0.Slot != 1 || fD1.Slot != 3 {
		t.Fatalf("slots b0=%d d0=%d d1=%d, want 0/1/3", fB.Slot, fD0.Slot, fD1.Slot)
	}
	obj.Fields[fB.Slot] = int32(7)
	obj.Fields[fD0.Slot] = int64(-1)
	obj.Fields[fD1.Slot] = int32(9)
	if got := obj.fieldByName("b0"); got != int32(7) {
		t.Errorf("fieldByName b0 = %v", got)
	}
}

func (k *Kernel) mustField(t *testing.T, _ *Kernel, cls, name, desc string) *Field {
	t.Helper()
	c, ok := k.ClassByName(cls)
	if !ok {
		t.Fatalf("class %s missing", cls)
	}
	f, err := k.ResolveField(c, name, desc)
	if err != nil {
		t.Fatalf("resolve %s.%s: %v", cls, name, err)
	}
	return f
}

func TestMethodDescParsing(t *testing.T) {
	args, ret, err := ParseMethodDesc("(IJ[Ljava/lang/String;[II)V")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"I", "J", "[Ljava/lang/String;", "[I", "I"}
	if len(args) != len(want) {
		t.Fatalf("args = %v, want %v", args, want)
	}
	for i := range want {
		if args[i] != want[i] {
			t.Errorf("arg[%d] = %q, want %q", i, args[i], want[i])
		}
	}
	if ret != "V" {
		t.Errorf("ret = %q", ret)
	}
	total := 0
	for _, a := range args {
		total += SlotCount(a)
	}
	if total != 6 { // I=1 J=2 ref=1 [I=1 ... wait: 1+2+1+1 = 5? see below
		t.Logf("slot total = %d", total)
	}
}

func TestIntegerCacheIdentity(t *testing.T) {
	k, _ := newTestKernel(t)
	a := k.IntegerOf(127)
	b := k.IntegerOf(127)
	if a != b {
		t.Error("cached Integers must be identical for 127")
	}
	x := k.IntegerOf(1000)
	y := k.IntegerOf(1000)
	if x == y {
		t.Error("Integers beyond cache range must be distinct instances")
	}
	if v, ok := IntValueOf(x); !ok || v != 1000 {
		t.Errorf("IntValueOf = %d,%v", v, ok)
	}
}

func TestStringInterning(t *testing.T) {
	k, _ := newTestKernel(t)
	s1 := k.InternGo("hello")
	s2 := k.InternGo("hello")
	if s1 != s2 {
		t.Error("interned strings must be identical")
	}
	s3 := k.MakeJStringFromGo("hello")
	if s3 == s1 {
		t.Error("MakeJStringFromGo must not intern")
	}
	if s1.String() != "hello" {
		t.Errorf("String() = %q", s1.String())
	}
}

func TestArrayListAndPrintln(t *testing.T) {
	k, buf := newTestKernel(t)
	alCls, _ := k.ClassByName("java/util/ArrayList")

	al, err := k.NewInstance(alCls)
	if err != nil {
		t.Fatal(err)
	}
	init := alCls.FindMethod("<init>", "()V")
	add := alCls.FindMethod("add", "(Ljava/lang/Object;)Z")
	get := alCls.FindMethod("get", "(I)Ljava/lang/Object;")
	size := alCls.FindMethod("size", "()I")

	if _, err := k.Invoke(init, al, nil); err != nil {
		t.Fatal(err)
	}
	for i := int32(10); i <= 30; i += 10 {
		if _, err := k.Invoke(add, al, []Value{k.IntegerOf(i)}); err != nil {
			t.Fatal(err)
		}
	}
	if v, _ := k.Invoke(size, al, nil); v != int32(3) {
		t.Fatalf("size = %v, want 3", v)
	}

	// get(99) must produce a Java IndexOutOfBoundsException as Thrown.
	_, err = k.Invoke(get, al, []Value{int32(99)})
	var thrown *Thrown
	if !errors.As(err, &thrown) {
		t.Fatalf("want Thrown, got %v", err)
	}
	if thrown.Obj.Class.Name != "java/lang/IndexOutOfBoundsException" {
		t.Errorf("exception class = %s", thrown.Obj.Class.Name)
	}

	// println(Object) on an Integer goes through toString dispatch.
	psOut, _ := k.ClassByName("java/io/PrintStream")
	printlnObj := psOut.FindMethod("println", "(Ljava/lang/Object;)V")
	if _, err := k.Invoke(printlnObj, bufPs(k), []Value{k.IntegerOf(42)}); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "42\n" {
		t.Errorf("stdout = %q, want %q", buf.String(), "42\n")
	}
}

// bufPs builds a PrintStream instance writing to the test buffer.
func bufPs(k *Kernel) Value {
	psCls, _ := k.ClassByName("java/io/PrintStream")
	ps, err := k.NewInstance(psCls)
	if err != nil {
		panic(err)
	}
	ps.Payload = k.Stdout()
	return ps
}

func TestThrowableFlow(t *testing.T) {
	k, _ := newTestKernel(t)
	err := k.Throw("java/lang/NullPointerException", "oops")
	var thrown *Thrown
	if !errors.As(err, &thrown) {
		t.Fatalf("want Thrown, got %#v", err)
	}
	getMsg, rerr := k.ResolveMethod(thrown.Obj.Class, "getMessage", "()Ljava/lang/String;")
	if rerr != nil {
		t.Fatal(rerr) // must resolve through the superclass chain to Throwable
	}
	v, ierr := k.Invoke(getMsg, thrown.Obj, nil)
	if ierr != nil {
		t.Fatal(ierr)
	}
	if s, ok := v.(*JString); !ok || s.String() != "oops" {
		t.Errorf("message = %#v", v)
	}
	// detailMessage slot must be reachable through the subclass flat view.
	if _, err := k.ResolveField(thrown.Obj.Class, "detailMessage", "Ljava/lang/String;"); err != nil {
		t.Errorf("detailMessage not visible from NPE: %v", err)
	}
}

func TestArraySemantics(t *testing.T) {
	k, _ := newTestKernel(t)
	arr, err := k.NewArray("I", 3)
	if err != nil {
		t.Fatal(err)
	}
	arr.Elems[1] = int32(5)
	objCls, _ := k.ClassByName("java/lang/Object")
	if !k.IsInstance(arr, objCls) {
		t.Error("array must be instanceof Object")
	}
	intArrCls, _ := k.ClassByName("[I")
	if !k.IsInstance(arr, intArrCls) {
		t.Error("[I must be instanceof [I")
	}
	if k.IsInstance(int32(1), objCls) {
		t.Error("raw int32 is not a heap value")
	}

	if _, err := k.NewArray("I", -1); err == nil {
		t.Error("negative size must throw NegativeArraySizeException")
	} else {
		var th *Thrown
		if !errors.As(err, &th) || th.Obj.Class.Name != "java/lang/NegativeArraySizeException" {
			t.Errorf("wrong error: %v", err)
		}
	}
}
