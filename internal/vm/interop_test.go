package vm

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"catty/internal/gen"
	"catty/internal/kernel"
)

// registerTestBridge wires an embedder-side Go bridge using the raw
// DefineClass surface — exactly what an embedding application would do
// today (R-0008 spike A). At runtime the synthesized class shadows the
// compiled stub in the classpath.
func registerTestBridge(k *kernel.Kernel) {
	if _, err := k.DefineClass(&kernel.ClassDef{
		Name: "go/Bridge",
		Methods: []kernel.MethodDef{
			{Name: "add", Desc: "(II)I", Flags: 0x0009, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				a, _ := args[0].(int32)
				b, _ := args[1].(int32)
				return a + b, nil
			}},
			{Name: "greet", Desc: "(Ljava/lang/String;)Ljava/lang/String;", Flags: 0x0009, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				js, _ := args[0].(*kernel.JString)
				return ctx.NewStringGo("hello, " + js.Go() + "!"), nil
			}},
			{Name: "fail", Desc: "()V", Flags: 0x0009, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				return nil, ctx.Throw("java/lang/RuntimeException", "go-side failure")
			}},
		},
	}); err != nil {
		panic(err)
	}
}

// TestInteropFromAOT pins R-0008 spike A: emitted Java code calls
// embedder-registered Go functions — value conversion both ways and error
// propagation into Java's catch clause.
func TestInteropFromAOT(t *testing.T) {
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	registerTestBridge(k)
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)
	cls, err := loader.Load("GoCall")
	if err != nil {
		t.Fatalf("load GoCall: %v", err)
	}
	gen.Install(k)

	th := New(k)
	mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatal(err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
	if _, err := th.Call(mainM, nil, []kernel.Value{argsArr}); err != nil {
		t.Fatalf("GoCall failed: %v", err)
	}
	want := strings.Join([]string{
		"42",
		"hello, catty!",
		"caught:go-side failure",
		"",
	}, "\n")
	if out.String() != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out.String(), want)
	}
	fmt.Fprintln(&bytes.Buffer{}) // keep fmt import stable for debug tweaks
}
