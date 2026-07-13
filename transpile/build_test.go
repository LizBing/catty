package transpile

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"catty/rtda"
)

// JVM access flag literals (stable values from JVMS §4.6).
const (
	flagStatic = 0x0008
	flagPublic = 0x0001
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

// TestInitClosureHasClinit verifies that the ADR-0025 initialization
// predecessor closure correctly identifies <clinit>-bearing classes and
// interfaces across the superclass chain and default-bearing superinterface
// graph. These are direct evidence for the AOT build-time rejection check.
func TestInitClosureHasClinit(t *testing.T) {
	// Case 1: target has no <clinit>, superclass has <clinit>.
	// new Child — must identify Parent.<clinit> in the closure.
	t.Run("superclass-has-clinit", func(t *testing.T) {
		parent := rtda.NewSyntheticClass("Parent", nil)
		parent.AddMethod(rtda.InterpretedMethod(parent, "<clinit>", "()V", flagStatic, 0, 0, nil, nil))
		child := rtda.NewSyntheticClass("Child", parent)
		got := initClosureHasClinit(child)
		if got != "Parent" {
			t.Errorf("initClosureHasClinit(Child) = %q, want %q", got, "Parent")
		}
	})

	// Case 2: target has no <clinit>, default-bearing interface has <clinit>.
	// new Impl — must identify Iface.<clinit> in the closure.
	t.Run("interface-has-clinit", func(t *testing.T) {
		iface := rtda.NewSyntheticClass("Iface", nil)
		iface.MarkInterface()
		iface.AddMethod(rtda.InterpretedMethod(iface, "<clinit>", "()V", flagStatic, 0, 0, nil, nil))
		iface.AddMethod(rtda.InterpretedMethod(iface, "m", "()V", flagPublic, 0, 0, nil, nil))
		impl := rtda.NewSyntheticClass("Impl", nil)
		impl.AddInterface(iface)
		got := initClosureHasClinit(impl)
		if got != "Iface" {
			t.Errorf("initClosureHasClinit(Impl) = %q, want %q", got, "Iface")
		}
	})

	// Case 3: no <clinit> anywhere in the closure.
	t.Run("no-clinit", func(t *testing.T) {
		parent := rtda.NewSyntheticClass("Parent", nil)
		child := rtda.NewSyntheticClass("Child", parent)
		got := initClosureHasClinit(child)
		if got != "" {
			t.Errorf("initClosureHasClinit(Child) = %q, want %q", got, "")
		}
	})

	// Case 4: interface target only checks itself, not superinterfaces.
	// Per ADR-0025, initializing an interface does NOT initialize
	// superinterfaces. Only class initialization cascades to interfaces.
	t.Run("interface-ignores-superinterface", func(t *testing.T) {
		superIface := rtda.NewSyntheticClass("SuperIface", nil)
		superIface.MarkInterface()
		superIface.AddMethod(rtda.InterpretedMethod(superIface, "<clinit>", "()V", flagStatic, 0, 0, nil, nil))
		superIface.AddMethod(rtda.InterpretedMethod(superIface, "dm", "()V", flagPublic, 0, 0, nil, nil))
		subIface := rtda.NewSyntheticClass("SubIface", nil)
		subIface.MarkInterface()
		subIface.AddInterface(superIface)
		got := initClosureHasClinit(subIface)
		if got != "" {
			t.Errorf("initClosureHasClinit(SubIface) = %q, want %q (interface init does not recurse into superinterfaces)", got, "")
		}
	})
}
