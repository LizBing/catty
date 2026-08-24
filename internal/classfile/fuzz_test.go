package classfile

import (
	"os"
	"path/filepath"
	"testing"
)

// FuzzParse hammers the classfile parser — the largest untrusted input
// surface (DEBT-0001). Any panic is a bug; well-formedness errors are not.
// Run via `make fuzz` (bounded) or plain go test -fuzz for long sessions.
func FuzzParse(f *testing.F) {
	seeds, _ := filepath.Glob("../../testdata/**/*.class")
	for _, s := range seeds {
		if b, err := os.ReadFile(s); err == nil {
			f.Add(b)
		}
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parser panic: %v", r)
			}
		}()
		_, _ = Parse(data)
	})
}
