package kernel

import "testing"

// BenchmarkBoxedIntFast measures the out-of-cache boxing fast path
// (P-0009 U5): atomic class cache + specialized single-field layout,
// no locking. The pre-U5 path paid k.mu + registry map lookup + generic
// field default-fill per box.
func BenchmarkBoxedIntFast(b *testing.B) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	_ = k.IntegerOf(0) // warm class cache
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = k.newBoxedIntFast(int32(i))
	}
}

// BenchmarkBoxedIntLegacy is the pre-U5 shape kept side-by-side for the
// R-0007/R-0009 evidence trail: global lock + registry lookup + generic
// NewInstance-style default fill.
func BenchmarkBoxedIntLegacy(b *testing.B) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})
	cls := k.integerCls()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k.mu.Lock()
		box := &Instance{}
		box.Class = cls
		box.Fields = make([]Value, cls.layoutSize)
		f := cls.fieldsByKey[memberKey("value", "I")]
		if f != nil {
			box.Fields[f.Slot] = int32(i)
		} else {
			box.Fields[0] = int32(i)
		}
		k.mu.Unlock()
	}
}
