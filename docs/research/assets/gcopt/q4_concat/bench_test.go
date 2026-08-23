// Package q4_concat compares string-building strategies the emitter might
// generate for javac's StringBuilder pattern: Go `+` chains, strings.Builder,
// and preallocated vs naive []byte append.
package q4_concat

import (
	"strings"
	"testing"
)

// pieces is a package-level var (not a const) so the compiler cannot fold the
// concatenations; indexing by i&7 gives realistic varying operands.
var pieces = []string{"a", "bb", "ccc", "dddd", "eeeee", "ffffff", "ggggggg", "hhhhhhhh"}

var sink int

const loopN = 1000

// Loop cases: dynamic length, the shape that separates linear from quadratic.

func BenchmarkLoopPlus(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := ""
		for j := 0; j < loopN; j++ {
			s = s + pieces[j&7]
		}
		sink += len(s)
	}
}

func BenchmarkLoopBuilder(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.Grow(loopN * 4)
		for j := 0; j < loopN; j++ {
			sb.WriteString(pieces[j&7])
		}
		sink += sb.Len()
	}
}

func BenchmarkLoopAppendPre(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := make([]byte, 0, loopN*4)
		for j := 0; j < loopN; j++ {
			buf = append(buf, pieces[j&7]...)
		}
		sink += len(buf)
	}
}

func BenchmarkLoopAppendNaive(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf []byte
		for j := 0; j < loopN; j++ {
			buf = append(buf, pieces[j&7]...)
		}
		sink += len(buf)
	}
}

// Fixed cases: a statically-known small number of parts, the shape javac
// generates for `a + b + c`. Compare a single + chain against Builder.

func BenchmarkFixedPlus4(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := pieces[i&7] + pieces[(i+1)&7] + pieces[(i+2)&7] + pieces[(i+3)&7]
		sink += len(s)
	}
}

func BenchmarkFixedBuilder4(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.Grow(16)
		sb.WriteString(pieces[i&7])
		sb.WriteString(pieces[(i+1)&7])
		sb.WriteString(pieces[(i+2)&7])
		sb.WriteString(pieces[(i+3)&7])
		sink += sb.Len()
	}
}

func BenchmarkFixedPlus8(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := pieces[i&7] + pieces[(i+1)&7] + pieces[(i+2)&7] + pieces[(i+3)&7] +
			pieces[(i+4)&7] + pieces[(i+5)&7] + pieces[(i+6)&7] + pieces[(i+7)&7]
		sink += len(s)
	}
}

func BenchmarkFixedBuilder8(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var sb strings.Builder
		sb.Grow(32)
		for j := 0; j < 8; j++ {
			sb.WriteString(pieces[(i+j)&7])
		}
		sink += sb.Len()
	}
}
