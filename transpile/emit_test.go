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
	src, err := Emit(fib, ir, cl)
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

// TestEmitFirst is A2.1's ref+array validation: emit ArrayOps.first(int[])I
// (aload the array ref, iaload, ireturn) and confirm the fresh-per-def emitter
// produces a ref-typed signature and array-element access that compiles. (We
// compile-check rather than execute: running it needs an int[] *rtda.Object,
// which needs the runtime bridge — A2.2.)
func TestEmitFirst(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	first := cl.LoadClass("ArrayOps").GetMethod("first", "([I)I")
	if first == nil {
		t.Fatal("ArrayOps.first not found")
	}
	ir, err := lowering.Lower(first)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	src, err := Emit(first, ir, cl)
	if err != nil {
		t.Fatalf("emit: %v", err)
	}
	t.Logf("emitted first:\n%s", src)

	if !strings.Contains(src, "func ArrayOps_first(l0 *rtda.Object) int32") {
		t.Errorf("ref-typed signature missing:\n%s", src)
	}
	if !strings.Contains(src, "ArrayElementSlot") {
		t.Errorf("array-element access missing:\n%s", src)
	}

	// Compile-check: the emitted func must be valid Go (gate for the Go-source
	// rules on ref code). main calls first(nil) — compiles, not run.
	out := t.TempDir()
	program := "package main\n\nimport \"catty/rtda\"\n\n" + src +
		"\nfunc main() { _ = ArrayOps_first((*rtda.Object)(nil)) }\n"
	mainPath := filepath.Join(out, "main.go")
	if err := os.WriteFile(mainPath, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	if buildOut, buildErr := exec.Command("go", "build", "-o", filepath.Join(out, "bin"), mainPath).CombinedOutput(); buildErr != nil {
		t.Fatalf("emitted first does not compile: %v\n%s\n--- source ---\n%s", buildErr, buildOut, program)
	}
}

// TestEmitHelloWorld is the A2.2 milestone: transpile HelloWorld.main (getstatic
// System.out, ldc String, invokevirtual println, int math) and run it natively
// via the runtime bridge, asserting output is byte-identical to java.
func TestEmitHelloWorld(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	main := cl.LoadClass("HelloWorld").GetMethod("main", "([Ljava/lang/String;)V")
	if main == nil {
		t.Fatal("HelloWorld.main not found")
	}
	ir, err := lowering.Lower(main)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	src, err := Emit(main, ir, cl)
	if err != nil {
		t.Fatalf("emit: %v", err)
	}
	t.Logf("emitted HelloWorld.main:\n%s", src)

	out := t.TempDir()
	program := "package main\n\nimport (\n\t\"catty/runtime\"\n\t\"catty/rtda\"\n)\n\n" + src +
		"\nfunc main() {\n\truntime.Bootstrap(\".\", \"HelloWorld\")\n\t" +
		"HelloWorld_main((*rtda.Object)(nil))\n}\n"
	mainPath := filepath.Join(out, "hw.go")
	if err := os.WriteFile(mainPath, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}

	bin := filepath.Join(out, "hwbin")
	if buildOut, buildErr := exec.Command("go", "build", "-o", bin, mainPath).CombinedOutput(); buildErr != nil {
		t.Fatalf("go build emitted HelloWorld failed: %v\n%s\n--- source ---\n%s", buildErr, buildOut, program)
	}

	// Run where HelloWorld.class lives, so runtime.Bootstrap(".") finds it.
	cmd := exec.Command(bin)
	cmd.Dir = dir
	got, err := cmd.Output()
	if err != nil {
		t.Fatalf("run emitted HelloWorld: %v", err)
	}
	if want := "Hello, World!\n42\n"; string(got) != want {
		t.Errorf("emitted HelloWorld output = %q, want %q", string(got), want)
	}
}

// TestEmitSum is the A2.3 milestone: AOT-execute a loop. ArrayOps.sum(int[]) has
// a for-loop (a merge with an empty operand stack at the head) — no phi needed,
// since loop state lives in mutable locals. Asserts the correct sum.
func TestEmitSum(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	sum := cl.LoadClass("ArrayOps").GetMethod("sum", "([I)I")
	if sum == nil {
		t.Fatal("ArrayOps.sum not found")
	}
	ir, err := lowering.Lower(sum)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	src, err := Emit(sum, ir, cl)
	if err != nil {
		t.Fatalf("emit: %v", err)
	}
	t.Logf("emitted sum:\n%s", src)

	out := t.TempDir()
	program := "package main\n\nimport (\n\t\"catty/runtime\"\n\t\"catty/rtda\"\n\t\"fmt\"\n)\n\n" + src +
		"\nfunc main() {\n\truntime.Bootstrap(\".\", \"ArrayOps\")\n\t" +
		"fmt.Println(ArrayOps_sum(runtime.NewIntArray(1, 2, 3, 4, 5)))\n}\n"
	mainPath := filepath.Join(out, "sum.go")
	if err := os.WriteFile(mainPath, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(out, "sumbin")
	if buildOut, buildErr := exec.Command("go", "build", "-o", bin, mainPath).CombinedOutput(); buildErr != nil {
		t.Fatalf("go build emitted sum failed: %v\n%s\n--- source ---\n%s", buildErr, buildOut, program)
	}
	cmd := exec.Command(bin)
	cmd.Dir = dir
	got, err := cmd.Output()
	if err != nil {
		t.Fatalf("run emitted sum: %v", err)
	}
	if want := "15\n"; string(got) != want {
		t.Errorf("emitted sum([1,2,3,4,5]) = %q, want %q", string(got), want)
	}
}

// TestEmitDiamondGate confirms the gate refuses merges that leave a value on the
// operand stack (diamonds like `a > b ? a : b`), which need phi insertion.
func TestEmitDiamondGate(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	max := cl.LoadClass("ArrayOps").GetMethod("max", "(II)I")
	if max == nil {
		t.Fatal("ArrayOps.max not found")
	}
	ir, err := lowering.Lower(max)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	if _, err := Emit(max, ir, cl); err == nil {
		t.Error("expected a non-empty-stack-merge error for ArrayOps.max (ternary), got nil")
	}
}

// TestEmitOOP is the A2.2b milestone: AOT-execute OOP — new (+ interpreted
// <init>), putfield, getfield, a user invokevirtual (interpreted via the bridge),
// and a native println. Prints b.v + b.doubled() = 21 + 42 = 63.
func TestEmitOOP(t *testing.T) {
	dir := compileFixtures(t)
	cl := classloader.New(classpath.Parse(dir))
	main := cl.LoadClass("OOPAot").GetMethod("main", "([Ljava/lang/String;)V")
	if main == nil {
		t.Fatal("OOPAot.main not found")
	}
	ir, err := lowering.Lower(main)
	if err != nil {
		t.Fatalf("lower: %v", err)
	}
	src, err := Emit(main, ir, cl)
	if err != nil {
		t.Fatalf("emit: %v", err)
	}
	t.Logf("emitted OOPAot.main:\n%s", src)

	out := t.TempDir()
	program := "package main\n\nimport (\n\t\"catty/runtime\"\n\t\"catty/rtda\"\n)\n\n" + src +
		"\nfunc main() {\n\truntime.Bootstrap(\".\", \"OOPAot\")\n\t" +
		"OOPAot_main((*rtda.Object)(nil))\n}\n"
	mainPath := filepath.Join(out, "oop.go")
	if err := os.WriteFile(mainPath, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(out, "oopbin")
	if buildOut, buildErr := exec.Command("go", "build", "-o", bin, mainPath).CombinedOutput(); buildErr != nil {
		t.Fatalf("go build emitted OOPAot failed: %v\n%s\n--- source ---\n%s", buildErr, buildOut, program)
	}
	cmd := exec.Command(bin)
	cmd.Dir = dir
	got, err := cmd.Output()
	if err != nil {
		t.Fatalf("run emitted OOPAot: %v", err)
	}
	if want := "63\n"; string(got) != want {
		t.Errorf("emitted OOPAot output = %q, want %q", string(got), want)
	}
}
