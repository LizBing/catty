package hybrid

import (
	"testing"

	"objbench/common"
)

// Group 1: monomorphic call site (concrete type direct call).
func BenchmarkCallMonoD3(b *testing.B) {
	objs := [4]*D3Leaf{newD3Leaf(), newD3Leaf(), newD3Leaf(), newD3Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += objs[i&3].Compute()
	}
}

func BenchmarkCallMonoD6(b *testing.B) {
	objs := [4]*D6Leaf{newD6Leaf(), newD6Leaf(), newD6Leaf(), newD6Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += objs[i&3].Compute()
	}
}

// Group 2a: bimorphic call site (interface dispatch over two concrete types).
func BenchmarkCallBiD3(b *testing.B) {
	animals := [2]Animal{newMega3A(), newMega3B()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += animals[i&1].Compute()
	}
}

func BenchmarkCallBiD6(b *testing.B) {
	animals := [2]Animal{newMega6A(), newMega6B()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += animals[i&1].Compute()
	}
}

// Group 2b: megamorphic call site (interface dispatch over eight types).
func BenchmarkCallMegaD3(b *testing.B) {
	animals := [8]Animal{
		newMega3A(), newMega3B(), newMega3C(), newMega3D(),
		newMega3E(), newMega3F(), newMega3G(), newMega3H(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += animals[i&7].Compute()
	}
}

func BenchmarkCallMegaD6(b *testing.B) {
	animals := [8]Animal{
		newMega6A(), newMega6B(), newMega6C(), newMega6D(),
		newMega6E(), newMega6F(), newMega6G(), newMega6H(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += animals[i&7].Compute()
	}
}

// Group 3a: field read (struct selector).
func BenchmarkFieldReadD3(b *testing.B) {
	objs := [4]*D3Leaf{newD3Leaf(), newD3Leaf(), newD3Leaf(), newD3Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := objs[i&3]
		sink += l.f0 + l.f1 + l.f2
	}
}

func BenchmarkFieldReadD6(b *testing.B) {
	objs := [4]*D6Leaf{newD6Leaf(), newD6Leaf(), newD6Leaf(), newD6Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := objs[i&3]
		sink += l.f0 + l.f1 + l.f2 + l.f3 + l.f4 + l.f5
	}
}

// Group 3b: field write (struct selector via noinline setter).
func BenchmarkFieldWriteD3(b *testing.B) {
	objs := [4]*D3Leaf{newD3Leaf(), newD3Leaf(), newD3Leaf(), newD3Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setD3Leaf(objs[i&3], common.Seed+int64(i))
		sink += objs[i&3].f2
	}
}

func BenchmarkFieldWriteD6(b *testing.B) {
	objs := [4]*D6Leaf{newD6Leaf(), newD6Leaf(), newD6Leaf(), newD6Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setD6Leaf(objs[i&3], common.Seed+int64(i))
		sink += objs[i&3].f5
	}
}

// Group 4: interface assertion downcast (always succeeds, monomorphic).
func BenchmarkDowncastD3(b *testing.B) {
	as := [4]Animal{newD3Leaf(), newD3Leaf(), newD3Leaf(), newD3Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if l, ok := as[i&3].(*D3Leaf); ok {
			sink += l.f2
		}
	}
}

func BenchmarkDowncastD6(b *testing.B) {
	as := [4]Animal{newD6Leaf(), newD6Leaf(), newD6Leaf(), newD6Leaf()}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if l, ok := as[i&3].(*D6Leaf); ok {
			sink += l.f5
		}
	}
}
