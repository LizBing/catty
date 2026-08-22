package uniform

import (
	"testing"

	"objbench/common"
)

// Group 1: monomorphic call site (type assertion, then call).
func BenchmarkCallMonoD3(b *testing.B) {
	objs := [4]*Object{
		newObject3(classD3Leaf), newObject3(classD3Leaf),
		newObject3(classD3Leaf), newObject3(classD3Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += d3LeafCall(objs[i&3])
	}
}

func BenchmarkCallMonoD6(b *testing.B) {
	objs := [4]*Object{
		newObject6(classD6Leaf), newObject6(classD6Leaf),
		newObject6(classD6Leaf), newObject6(classD6Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += d6LeafCall(objs[i&3])
	}
}

// Group 2a: bimorphic call site (two-way type-switch dispatch).
func BenchmarkCallBiD3(b *testing.B) {
	objs := [2]*Object{
		newObject3(classMega3A), newObject3(classMega3B),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += bi3Dispatch(objs[i&1])
	}
}

func BenchmarkCallBiD6(b *testing.B) {
	objs := [2]*Object{
		newObject6(classMega6A), newObject6(classMega6B),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += bi6Dispatch(objs[i&1])
	}
}

// Group 2b: megamorphic call site (eight-way type-switch dispatch).
func BenchmarkCallMegaD3(b *testing.B) {
	objs := [8]*Object{
		newObject3(classMega3A), newObject3(classMega3B),
		newObject3(classMega3C), newObject3(classMega3D),
		newObject3(classMega3E), newObject3(classMega3F),
		newObject3(classMega3G), newObject3(classMega3H),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += mega3Dispatch(objs[i&7])
	}
}

func BenchmarkCallMegaD6(b *testing.B) {
	objs := [8]*Object{
		newObject6(classMega6A), newObject6(classMega6B),
		newObject6(classMega6C), newObject6(classMega6D),
		newObject6(classMega6E), newObject6(classMega6F),
		newObject6(classMega6G), newObject6(classMega6H),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += mega6Dispatch(objs[i&7])
	}
}

// Group 3a: field read (slot array indexing).
func BenchmarkFieldReadD3(b *testing.B) {
	objs := [4]*Object{
		newObject3(classD3Leaf), newObject3(classD3Leaf),
		newObject3(classD3Leaf), newObject3(classD3Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := objs[i&3]
		sink += o.f[0] + o.f[1] + o.f[2]
	}
}

func BenchmarkFieldReadD6(b *testing.B) {
	objs := [4]*Object{
		newObject6(classD6Leaf), newObject6(classD6Leaf),
		newObject6(classD6Leaf), newObject6(classD6Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		o := objs[i&3]
		sink += o.f[0] + o.f[1] + o.f[2] + o.f[3] + o.f[4] + o.f[5]
	}
}

// Group 3b: field write (slot array indexing via noinline setter).
func BenchmarkFieldWriteD3(b *testing.B) {
	objs := [4]*Object{
		newObject3(classD3Leaf), newObject3(classD3Leaf),
		newObject3(classD3Leaf), newObject3(classD3Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setObject3(objs[i&3], common.Seed+int64(i))
		sink += objs[i&3].f[2]
	}
}

func BenchmarkFieldWriteD6(b *testing.B) {
	objs := [4]*Object{
		newObject6(classD6Leaf), newObject6(classD6Leaf),
		newObject6(classD6Leaf), newObject6(classD6Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setObject6(objs[i&3], common.Seed+int64(i))
		sink += objs[i&3].f[5]
	}
}

// Group 4: type-tag downcast (always succeeds, monomorphic).
func BenchmarkDowncastD3(b *testing.B) {
	objs := [4]*Object{
		newObject3(classD3Leaf), newObject3(classD3Leaf),
		newObject3(classD3Leaf), newObject3(classD3Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if objs[i&3].cls == classD3Leaf {
			sink += objs[i&3].f[2]
		}
	}
}

func BenchmarkDowncastD6(b *testing.B) {
	objs := [4]*Object{
		newObject6(classD6Leaf), newObject6(classD6Leaf),
		newObject6(classD6Leaf), newObject6(classD6Leaf),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if objs[i&3].cls == classD6Leaf {
			sink += objs[i&3].f[5]
		}
	}
}
