package transpile

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"catty/classloader"
	"catty/classpath"
	"catty/lowering"
)

// compileFixtures compiles tests/fixtures/*.java into a temp dir for the
// classloader to read.
func compileFixtures(t *testing.T) string {
	t.Helper()
	src, err := filepath.Abs(filepath.Join("..", "tests", "fixtures"))
	if err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatal(err)
	}
	args := []string{"-source", "8", "-target", "8", "-nowarn", "-d", out}
	for _, e := range entries {
		if filepath.Ext(e.Name()) == ".java" {
			args = append(args, filepath.Join(src, e.Name()))
		}
	}
	if err := exec.Command("javac", args...).Run(); err != nil {
		t.Fatalf("javac failed: %v", err)
	}
	return out
}

// TestEmitFib is the A1 milestone: lower Fibonacci.fib, emit it to Go, compile +
// run it natively, and check it matches Java's result. The go build is itself the
// gate for the Go goto/label/unused-var rules.
func TestEmitFib(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	fib := cl.LoadClass("Fibonacci").GetMethod("fib", "(I)I")
	if fib == nil {
		t.Fatal("Fibonacci.fib not found")
	}
	ir, err := lowering.Lower(fib)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	src, err := Emit(fib, ir)
	if err != nil {
		t.Fatalf("emit: %v", err)
	}
	t.Logf("emitted Go:\n%s", src)

	// Wrap the emitted function in a runnable main that prints fib(35).
	out := t.TempDir()
	program := "package main\n\nimport \"fmt\"\n\n" + src +
		"\nfunc main() { fmt.Println(Fibonacci_fib(35)) }\n"
	mainPath := filepath.Join(out, "main.go")
	if err := os.WriteFile(mainPath, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}

	// Compile the emitted program (this is the gate for Go's source rules). Build
	// the file directly: the temp dir is outside any module, so `go build <dir>`
	// would fail, but `go build <file>` (command-line-arguments mode) works.
	bin := filepath.Join(out, "fibbin")
	if buildOut, buildErr := exec.Command("go", "build", "-o", bin, mainPath).CombinedOutput(); buildErr != nil {
		t.Fatalf("go build emitted program failed: %v\n%s\n--- source ---\n%s",
			buildErr, buildOut, program)
	}

	// Correctness: fib(35) must equal Java's result.
	got, err := exec.Command(bin).Output()
	if err != nil {
		t.Fatalf("run emitted program: %v", err)
	}
	if want := "9227465"; strings.TrimSpace(string(got)) != want {
		t.Errorf("emitted fib(35) = %q, want %s", strings.TrimSpace(string(got)), want)
	}

	// Speed signal: native execution should be a tiny fraction of the
	// interpreter's 4.5 s (and of `java -Xint`'s 0.6 s).
	start := time.Now()
	if err := exec.Command(bin).Run(); err != nil {
		t.Fatalf("rerun: %v", err)
	}
	t.Logf("emitted fib(35) ran in %v (interpreter ≈ 4.5s, java -Xint ≈ 0.6s)", time.Since(start))
}
