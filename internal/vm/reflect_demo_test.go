package vm

import (
	"bytes"
	"strings"
	"testing"

	"catty/internal/gen"
	"catty/internal/kernel"
)

// TestReflectionDemoPinned pins the reflection minimal surface (P-0011)
// through the demo.ReflectDemo fixture on BOTH engines, byte-identical to
// the reference JVM output captured at fixture time. Covers forName,
// getDeclaredFields/Constructors, Field.getType/set, Method.invoke,
// isInstance and primitive Class identity — the framework-shaped surface.
func TestReflectionDemoPinned(t *testing.T) {
	want := strings.Join([]string{
		"class=demo.ReflectDemo",
		"field:tag:String",
		"field:count:int",
		"tag=hello,count=7",
		"isInstance=true",
		"intClass=true",
		`{"tag":"hello","count":7}`,
		"",
	}, "\n")
	for _, aot := range []bool{false, true} {
		var out bytes.Buffer
		k := kernel.New(kernel.Options{Stdout: &out})
		loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
		k.SetResolver(loader.Load)
		cls, err := loader.Load("demo/ReflectDemo")
		if err != nil {
			t.Fatalf("load: %v", err)
		}
		if aot {
			gen.Install(k)
		}
		th := New(k)
		mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
		if err != nil {
			t.Fatal(err)
		}
		argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
		if _, err := th.Call(mainM, nil, []kernel.Value{argsArr}); err != nil {
			t.Fatalf("[aot=%v] ReflectDemo failed: %v", aot, err)
		}
		if out.String() != want {
			t.Errorf("[aot=%v] stdout =\n%q\nwant\n%q", aot, out.String(), want)
		}
	}
}
