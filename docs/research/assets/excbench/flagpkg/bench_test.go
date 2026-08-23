package flagpkg

import (
	"testing"

	"excbench/common"
)

// Normal path, zero exceptions: every layer checks err, nothing is ever
// thrown.
func BenchmarkNormalPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, _ := aNorm(i)
		common.Sink += int64(v)
	}
}

// High-frequency throw-catch: d returns err every iteration, c catches one
// boundary up.
func BenchmarkThrowShallow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, _ := cShallow(i)
		common.Sink += int64(v)
	}
}

// Deep propagation: d errs, b/c re-return, a catches. Exception crosses four
// frames.
func BenchmarkDeepProp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, _ := aDeep(i)
		common.Sink += int64(v)
	}
}
