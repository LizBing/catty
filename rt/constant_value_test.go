package rt_test

import (
	"testing"

	"catty/rt"
)

// TestCattyConstantValueInitialization verifies the JVMS 5.5 ConstantValue
// step: a class applies its constant fields before its superclass and before
// its own <clinit>, so a superclass initializer can observe them.
func TestCattyConstantValueInitialization(t *testing.T) {
	var observed int32 = -1
	parentInit := func(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
		value, err := t.RT.StaticValue("Child", "VALUE", "I")
		if err != nil {
			return rt.ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalAccessError"))
		}
		observed = rt.AsInt(value)
		return rt.Normal()
	}
	childInit := func(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
		return rt.Normal()
	}
	image := &rt.Image{
		ABI:     rt.ABIRevision,
		Profile: rt.ProfileID,
		Classes: []*rt.ClassDef{
			{
				Name: "Parent", Super: "java/lang/Object", Flags: rt.AccPublic,
				Methods: []*rt.MethodDef{{
					Owner: "Parent", Name: "<clinit>", Descriptor: "()V",
					Engine: "aot", Entry: parentInit,
				}},
			},
			{
				Name: "Child", Super: "Parent", Flags: rt.AccPublic,
				StaticFields: []rt.FieldDef{{
					Name: "VALUE", Descriptor: "I", Flags: rt.AccStatic | rt.AccFinal,
					Constant: &rt.Constant{Kind: "int", Bits: 42},
				}},
				Methods: []*rt.MethodDef{{
					Owner: "Child", Name: "<clinit>", Descriptor: "()V",
					Engine: "aot", Entry: childInit,
				}},
			},
		},
	}
	r := rt.NewRuntime(image)
	if err := r.Install(); err != nil {
		t.Fatalf("install: %v", err)
	}
	child, err := r.Class("Child")
	if err != nil {
		t.Fatal(err)
	}
	thread := r.MainThread()
	if thrown := thread.EnsureInitialized(child); thrown != nil {
		t.Fatalf("Child initialization threw %v", thrown)
	}
	if observed != 42 {
		t.Fatalf("superclass initializer observed ConstantValue %d, want 42", observed)
	}
	value, err := r.StaticValue("Child", "VALUE", "I")
	if err != nil || rt.AsInt(value) != 42 {
		t.Fatalf("Child.VALUE = %v, %v; want 42", value, err)
	}
}
