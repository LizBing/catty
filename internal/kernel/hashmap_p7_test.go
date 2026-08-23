package kernel

import (
	"testing"
)

// TestHashMapNativesDirect exercises HashMap natives directly through
// the kernel API (bypassing both interpreter and AOT emission) to verify
// that the native implementations are correct in isolation.
func TestHashMapNativesDirect(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	ctx := &CallContext{K: k}

	cls, ok := k.ClassByName("java/util/HashMap")
	if !ok {
		t.Fatal("HashMap not bootstrapped")
	}
	obj, err := k.NewInstance(cls)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := natHashMapInit(ctx, obj, nil); err != nil {
		t.Fatal(err)
	}

	key := k.InternGo("the")
	val, _ := natWrapperValueOf("java/lang/Integer")(ctx, nil, []Value{int32(3)})

	// put("the", 3)
	if _, err := natHashMapPut(ctx, obj, []Value{key, val}); err != nil {
		t.Fatal(err)
	}

	// get("the") should return Integer(3)
	got, err := natHashMapGet(ctx, obj, []Value{key})
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("get returned nil after put")
	}
	in := got.(*Instance)
	if in.Fields[0].(int32) != 3 {
		t.Fatalf("expected 3, got %v", in.Fields[0])
	}

	// containsKey should be true
	c, _ := natHashMapContainsKey(ctx, obj, []Value{key})
	if c.(int32) != 1 {
		t.Fatal("containsKey returned false")
	}

	// size should be 1
	s, _ := natHashMapSize(ctx, obj, nil)
	if s.(int32) != 1 {
		t.Fatalf("size expected 1, got %d", s)
	}
}
