package q6_bce

import "testing"

var s = make([]int64, 256)

// These confirm the safe loop shapes all run at the same speed; the actual
// BCE evidence is the -d=ssa/check_bce/debug=1 output, not timing.

func BenchmarkSumRangeIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumRangeIndex(s)
	}
}

func BenchmarkSumRangeValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumRangeValue(s)
	}
}

func BenchmarkSumForLen(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumForLen(s)
	}
}

func BenchmarkSumForCache(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumForCache(s)
	}
}

func BenchmarkSumReverse(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumReverse(s)
	}
}

// SumWithN keeps the bounds check (n is opaque) without a helper call, so the
// delta vs SumForLen isolates the cost of the check itself.
func BenchmarkSumWithN(b *testing.B) {
	n := len(s)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumWithN(s, n)
	}
}

func BenchmarkSumOpaqueIndex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SumOpaqueIndex(s)
	}
}
