package vm

import (
	"bytes"
	"strings"
	"testing"

	"catty/internal/gen"
	"catty/internal/kernel"
)

// TestAOTCollectionsDemoMatchesInterpreter is the P-0006 T7 dual-path check
// for the exception-handling fixture: emitter-generated bodies (gen.Install)
// must produce byte-identical stdout to the interpreted path. It pins the
// handler-entry depth fix — before it, the AOT body crashed at the first
// handler with NoSuchMethodError because the exception receiver was never
// loaded into the operand slot.
func TestAOTCollectionsDemoMatchesInterpreter(t *testing.T) {
	interOut, err := runFixture(t, "CollectionsDemo.class")
	if err != nil {
		t.Fatalf("interpreted run failed: %v", err)
	}

	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)
	cls, err := loader.Load("CollectionsDemo")
	if err != nil {
		t.Fatalf("load CollectionsDemo: %v", err)
	}
	gen.Install(k) // emitted bodies take precedence over interpreted ones

	th := New(k)
	mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, err := k.NewArray("Ljava/lang/String;", 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := th.Call(mainM, nil, []kernel.Value{argsArr}); err != nil {
		t.Fatalf("AOT run failed: %v", err)
	}

	want := strings.Join([]string{
		"150",
		"caught: Index: 99, Size: 5",
		"div-by-zero caught",
		"sum=150",
		"has30",
		"",
	}, "\n")
	if out.String() != want {
		t.Errorf("AOT stdout =\n%q\nwant\n%q", out.String(), want)
	}
	if interOut != want {
		t.Errorf("interpreter baseline drifted: got %q, want %q", interOut, want)
	}
}
