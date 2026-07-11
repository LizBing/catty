package transpile

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestBuildProgram is the A4 milestone: BuildProgram transpiles a whole program
// (all emittable methods), go builds it, runs the native binary, and asserts
// output is byte-identical to java. Validated on HelloWorld (native-invoke
// only) and Fibonacci (main + fib both AOT'd, invokestatic → direct call).
func TestBuildProgram(t *testing.T) {
	dir := compileFixtures(t)

	for _, tc := range []struct {
		name string
		want string
	}{
		{"HelloWorld", "Hello, World!\n42\n"},
		{"Fibonacci", "0\n1\n1\n2\n3\n5\n8\n13\n21\n34\n55\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src, err := BuildProgram(tc.name, dir)
			if err != nil {
				t.Fatalf("BuildProgram(%s): %v", tc.name, err)
			}
			t.Logf("emitted program:\n%s", src)

			out := t.TempDir()
			mainPath := filepath.Join(out, "program.go")
			if err := os.WriteFile(mainPath, []byte(src), 0o644); err != nil {
				t.Fatal(err)
			}
			bin := filepath.Join(out, "program")
			if buildOut, buildErr := exec.Command("go", "build", "-o", bin, mainPath).CombinedOutput(); buildErr != nil {
				t.Fatalf("go build failed: %v\n%s\n--- source ---\n%s", buildErr, buildOut, src)
			}
			got, err := exec.Command(bin).Output()
			if err != nil {
				t.Fatalf("run binary: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("%s: got %q, want %q", tc.name, string(got), tc.want)
			}
		})
	}
}
