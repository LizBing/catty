package rtda

import "testing"

// TestFrameLongRoundTrip guards the two-slot long encoding: the interpreter's
// long/double support (loads/stores, fields, returns) all depend on PushLong /
// PopLong splitting and rejoining a 64-bit value across two adjacent slots.
func TestFrameLongRoundTrip(t *testing.T) {
	m := &Method{maxStack: 4, maxLocals: 4}
	f := NewFrame(nil, m)

	for _, v := range []int64{0, 1, -1, 1<<40, -(1 << 40), 0x123456789ABCDEF0} {
		f.PushLong(v)
		if got := f.PopLong(); got != v {
			t.Errorf("PopLong = %#x, want %#x", got, v)
		}
	}
}

// TestFrameDoubleRoundTrip exercises the same path via the double view.
func TestFrameDoubleRoundTrip(t *testing.T) {
	m := &Method{maxStack: 4, maxLocals: 4}
	f := NewFrame(nil, m)

	for _, v := range []float64{0, 1.5, -2.25, 3.14159, 1e300, -1e-300} {
		f.PushDouble(v)
		if got := f.PopDouble(); got != v {
			t.Errorf("PopDouble = %v, want %v", got, v)
		}
	}
}

// TestFrameLocalsLong covers SetLong/GetLong (local-variable two-slot access),
// used by lload/lstore and long parameter passing.
func TestFrameLocalsLong(t *testing.T) {
	m := &Method{maxStack: 2, maxLocals: 4}
	f := NewFrame(nil, m)

	v := int64(-0x7777AAAA55550000)
	f.SetLong(1, v)
	if got := f.GetLong(1); got != v {
		t.Errorf("GetLong = %#x, want %#x", got, v)
	}
}

// TestFrameRefGCNil checks that popping a reference clears the slot so the GC
// can reclaim the object — a subtle correctness/perf property of PopRef.
func TestFrameRefGCNil(t *testing.T) {
	m := &Method{maxStack: 2, maxLocals: 1}
	f := NewFrame(nil, m)

	obj := &Object{}
	f.PushRef(obj)
	_ = f.PopRef()
	// The just-popped slot must no longer hold a reference.
	if f.stack[0].ref != nil {
		t.Error("PopRef left a stale reference in the freed slot")
	}
}
