package panicpkg

import (
	"testing"

	"excbench/common"
)

// Normal path, zero exceptions. Every layer owns a handler, so each call pays
// a defer + recover(nil) check but nothing is ever thrown.
func BenchmarkNormalPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(aNorm(i))
	}
}

// Normal path, zero exceptions, no layer owns a handler: the pure call chain
// with no exception machinery at all (the realistic majority of methods).
func BenchmarkNormalPathNoHandlers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(aNoHandlers(i))
	}
}

// Normal path, zero exceptions, every layer owns a plain defer (finally /
// synchronized), but no recover.
func BenchmarkNormalPathFinally(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(aFinally(i))
	}
}

// High-frequency throw-catch: d panics every iteration, c catches one
// boundary up. Measures the throw + catch primitive.
func BenchmarkThrowShallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(cShallow(i))
	}
}

// Deep propagation (realistic): d panics, intermediate frames have no handler,
// a catches. Exception crosses four frames.
func BenchmarkDeepPropFree(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(aDeepFree(i))
	}
}

// Deep propagation (worst case): every intermediate frame recovers and
// re-panics before the matching handler at a runs.
func BenchmarkDeepPropRethrow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.Sink += int64(aDeepRethrow(i))
	}
}
