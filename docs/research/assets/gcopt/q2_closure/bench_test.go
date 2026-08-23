package q2_closure

import (
	"testing"

	"gcopt/common"
)

var sink int64

// Baseline: direct call.
func BenchmarkDirectCall(b *testing.B) {
	c := &C{n: 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += c.Add(common.Seed, int64(i))
	}
}

// Closure capturing a local by value, called in-loop (should inline, 0 allocs).
func BenchmarkClosureByValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := common.Seed + int64(i)
		f := func(y int64) int64 { return x + y }
		sink += f(common.Seed)
	}
}

// Closure capturing a local by reference, called in-loop (0 allocs if inlined).
func BenchmarkClosureByRef(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := common.Seed + int64(i)
		f := func() int64 { x++; return x }
		sink += f()
	}
}

// Closure that must heap-allocate because it crosses a noinline boundary.
func BenchmarkClosureHeap(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := common.Seed + int64(i)
		f := func() int64 { return x * 2 }
		sink += callFunc(f)
	}
}

// Method value on a concrete receiver (2 words, no heap alloc; call inlines).
func BenchmarkMethodValueConcrete(b *testing.B) {
	c := &C{n: 1}
	f := c.Add
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += f(common.Seed, int64(i))
	}
}

// Method value on an interface receiver: the captured itab keeps the call
// virtual even though the method value itself allocates nothing.
func BenchmarkMethodValueInterface(b *testing.B) {
	a := makeAdder()
	f := a.Add
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += f(common.Seed, int64(i))
	}
}

// Escaping closure capturing x by value: one heap allocation per iteration
// (the closure box embeds the copy of x).
func BenchmarkClosureEscapeByValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := makeAdderClosure(common.Seed + int64(i))
		sink += f(common.Seed)
	}
}

// Escaping closure capturing x by reference: two heap allocations per
// iteration (the closure plus the boxed x).
func BenchmarkClosureEscapeByRef(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := makeCounter(common.Seed + int64(i))
		sink += f()
	}
}
