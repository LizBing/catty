package kernel

import "testing"

func TestRouteCArraysFill(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	c, ok := k.ClassByName("java/util/Arrays")
	if !ok {
		t.Fatal("Arrays not defined")
	}
	if m := c.methodsByKey[memberKey("fill", "([CC)V")]; m == nil {
		keys := []string{}
		for kk := range c.methodsByKey {
			keys = append(keys, kk)
		}
		t.Fatalf("fill missing; have %v", keys)
	}
}

func TestRouteCArraysFillInvoke(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	c, ok := k.ClassByName("java/util/Arrays")
	if !ok {
		t.Fatal("no Arrays")
	}
	m, err := k.ResolveMethod(c, "fill", "([CC)V")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	owner := &intTestOwner{key: k.MintKey()}
	arr := &ArrayObj{Elems: make([]Value, 3)}
	if _, err := k.InvokeAs(owner, m, nil, []Value{arr, int32('x')}); err != nil {
		t.Fatalf("invoke: %v", err)
	}
	for i, e := range arr.Elems {
		if e.(int32) != 'x' {
			t.Fatalf("elem %d = %v", i, e)
		}
	}
}
