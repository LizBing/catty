package q1_devirt

import (
	"testing"

	"gcopt/common"
)

// sink accumulates results so the measured work is never dead-code eliminated.
var sink int64

// Direct concrete call, expected to fully inline to a couple of ALU ops.
func BenchmarkDirect(b *testing.B) {
	objs := [4]*C1{{n: 1}, {n: 1}, {n: 1}, {n: 1}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += objs[i&3].Add(common.Seed, int64(i))
	}
}

// Interface local assigned a concrete value in-loop: devirtualization case.
func BenchmarkInterfaceLocal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var v Adder = &C1{n: 1}
		sink += v.Add(common.Seed, int64(i))
	}
}

// Assertion-then-call over an opaque interface (monomorphic target).
func BenchmarkAssertThenCall(b *testing.B) {
	as := [4]Adder{makeC1(), makeC1(), makeC1(), makeC1()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if c, ok := as[i&3].(*C1); ok {
			sink += c.Add(common.Seed, int64(i))
		}
	}
}

// Genuinely virtual call: opaque interface, monomorphic but invisible type.
func BenchmarkOpaqueMono(b *testing.B) {
	as := [4]Adder{makeC1(), makeC1(), makeC1(), makeC1()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += as[i&3].Add(common.Seed, int64(i))
	}
}

// Bimorphic virtual call: two concrete types behind the same interface.
func BenchmarkOpaqueBi(b *testing.B) {
	as := [4]Adder{makeC1(), makeC2(), makeC1(), makeC2()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += as[i&3].Add(common.Seed, int64(i))
	}
}

// Opaque interface behind a struct field (realistic receiver slot shape).
func BenchmarkField(b *testing.B) {
	hs := [4]*Holder{{v: makeC1()}, {v: makeC1()}, {v: makeC1()}, {v: makeC1()}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += hs[i&3].v.Add(common.Seed, int64(i))
	}
}

// Budget cliff: a method just under the inline budget (M05, five statements)
// vs one just over it. The generator's -m output pins the exact boundary; this
// pair shows the runtime consequence of crossing it.
func BenchmarkBudgetUnder(b *testing.B) {
	c := &C1{n: 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += c.M05(common.Seed + int64(i))
	}
}

func BenchmarkBudgetOver(b *testing.B) {
	c := &C1{n: 1}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += c.M30(common.Seed + int64(i))
	}
}
