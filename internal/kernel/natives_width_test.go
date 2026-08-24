package kernel

import (
	"testing"
)

// TestNativeReturnWidths pins the native/Java boundary contract: every
// native returning a JVM value must hand back the canonical Go type for
// its descriptor (int32 for B/C/S/I/Z, int64, float32, float64). Raw
// internal widths (uint16 from JString.Chars) crashed the AOT path's
// type assertions while the interpreter's push-normalization masked them
// (DEBT-0016 residual root cause).
func TestNativeReturnWidths(t *testing.T) {
	k := New(Options{Stdout: widthDiscard{}, Stderr: widthDiscard{}})

	t.Run("String.charAt(I)C", func(t *testing.T) {
		js := k.InternGo("abc")
		m, err := k.ResolveMethod(js.Class, "charAt", "(I)C")
		if err != nil {
			t.Fatal(err)
		}
		v, err := k.Invoke(m, js, []Value{int32(1)})
		if err != nil {
			t.Fatal(err)
		}
		i, ok := v.(int32)
		if !ok {
			t.Fatalf("charAt returned %T (%v), want int32", v, v)
		}
		if i != 98 { // 'b'
			t.Fatalf("charAt(1) = %d, want 98", i)
		}
	})

	t.Run("String.length()I", func(t *testing.T) {
		js := k.InternGo("abcd")
		m, err := k.ResolveMethod(js.Class, "length", "()I")
		if err != nil {
			t.Fatal(err)
		}
		v, err := k.Invoke(m, js, nil)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := v.(int32); !ok {
			t.Fatalf("length returned %T (%v), want int32", v, v)
		}
	})
}

// discardWriter swallows diagnostics (mirrors kernel_test helper).
type widthDiscard struct{}

func (widthDiscard) Write(p []byte) (int, error) { return len(p), nil }
