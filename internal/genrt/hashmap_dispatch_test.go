package genrt

import (
	"testing"

	"catty/internal/kernel"
)

type discardWriter struct{}

func (discardWriter) Write(p []byte) (int, error) { return len(p), nil }

type testOwner struct{ key uint64 }

func (o *testOwner) OwnerKey() uint64 { return o.key }

// TestHashMapThroughGenRTPath exercises the exact same call chain as
// generated AOT code: New → CallSpecial(<init>) → CallVirtual(put) → get.
func TestHashMapThroughGenRTPath(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: discardWriter{}, Stderr: discardWriter{}})
	k.SetResolver(func(name string) (*kernel.Class, error) {
		return nil, nil
	})
	K = k // set genrt global

	testKey := &testOwner{key: k.MintKey()}
	th := testKey

	// Step 1: new HashMap
	obj := New("java/util/HashMap")
	if obj == nil {
		t.Fatal("New returned nil")
	}

	// Step 2: <init>()
	_, thrown := CallSpecial(th, obj, "java/util/HashMap", "<init>", "()V", []kernel.Value{})
	if thrown != nil {
		t.Fatalf("<init> threw: %v", thrown)
	}

	// Verify Payload was set by <init>
	in := obj
	if in.Payload == nil {
		t.Fatal("Payload is nil after <init> — natHashMapInit was not called!")
	}

	// Step 3: put("the", Integer(3))
	key := k.InternGo("the")
	intCls, _ := k.ClassByName("java/lang/Integer")
	intObj, err := k.NewInstance(intCls)
	if err != nil {
		t.Fatal(err)
	}
	intObj.Fields[0] = int32(3)

	_, thrown = CallVirtual(th, obj, "java/util/HashMap", "put",
		"(Ljava/lang/Object;Ljava/lang/Object;)Ljava/lang/Object;",
		[]kernel.Value{key, intObj})
	if thrown != nil {
		t.Fatalf("put threw: %v", thrown)
	}

	// Step 4: get("the")
	got, thrown := CallVirtual(th, obj, "java/util/HashMap", "get",
		"(Ljava/lang/Object;)Ljava/lang/Object;",
		[]kernel.Value{key})
	if thrown != nil {
		t.Fatalf("get threw: %v", thrown)
	}
	if got == nil {
		t.Fatal("get returned nil after successful put — THIS IS THE BUG")
	}

	// Verify value
	gotInst, ok := got.(*kernel.Instance)
	if !ok {
		t.Fatalf("get returned %T, want *Instance", got)
	}
	if gotInst.Fields[0].(int32) != 3 {
		t.Fatalf("expected 3, got %v", gotInst.Fields[0])
	}

	t.Log("HashMap works correctly through genrt dispatch path ✓")
}
