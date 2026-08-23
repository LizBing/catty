package q3_bigfunc

import (
	"testing"

	"gcopt/common"
)

var sink int64

// The size sweep: one function, N-step dependent chain. Measures whether
// per-step cost or compile time degrades as the function body grows.
func BenchmarkBigChain50(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain50(common.Seed + int64(i))
	}
}

func BenchmarkBigChain100(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain100(common.Seed + int64(i))
	}
}

func BenchmarkBigChain200(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain200(common.Seed + int64(i))
	}
}

func BenchmarkBigChain500(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain500(common.Seed + int64(i))
	}
}

func BenchmarkBigChain1000(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain1000(common.Seed + int64(i))
	}
}

func BenchmarkBigChain2000(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain2000(common.Seed + int64(i))
	}
}

func BenchmarkBigChain4000(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += BigChain4000(common.Seed + int64(i))
	}
}

// One 200-step function vs eight 25-step functions composed inline: does
// splitting into smaller bodies change the emitted code at all?
func BenchmarkSplitChain200(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += SplitChain200(common.Seed + int64(i))
	}
}
