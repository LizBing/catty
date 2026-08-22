package vm

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"catty/internal/kernel"
)

// runFixture loads testdata/<name> and runs its static main with an empty
// String[] argument, returning captured stdout.
func runFixture(t *testing.T, name string) (string, error) {
	t.Helper()
	data, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	cls, err := k.LoadClassBytes(data)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	th := New(k)
	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, err := k.NewArray("Ljava/lang/String;", 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = th.Call(main, nil, []kernel.Value{argsArr})
	return out.String(), err
}

func TestAcceptanceHelloWorld(t *testing.T) {
	out, err := runFixture(t, "HelloWorld.class")
	if err != nil {
		var th *kernel.Thrown
		if errors.As(err, &th) {
			t.Fatalf("uncaught %v", err)
		}
		t.Fatal(err)
	}
	if out != "Hello, Catty!\n" {
		t.Errorf("stdout = %q, want %q", out, "Hello, Catty!\n")
	}
}

func TestAcceptanceCollectionsDemo(t *testing.T) {
	out, err := runFixture(t, "CollectionsDemo.class")
	if err != nil {
		t.Fatalf("execution failed: %v", err)
	}
	want := strings.Join([]string{
		"150",
		"caught: Index: 99, Size: 5",
		"div-by-zero caught",
		"sum=150",
		"has30",
		"",
	}, "\n")
	if out != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out, want)
	}
}
