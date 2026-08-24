package vm

import (
	"bytes"
	"strings"
	"testing"

	"catty/internal/gen"
	"catty/internal/kernel"
)

// runTraceProbe executes a fixture that dies from an uncaught throwable on
// the interpreter path and returns the rendered uncaught report.
func runTraceProbe(t *testing.T, name string, aot bool) string {
	t.Helper()
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/cp"})
	k.SetResolver(loader.Load)
	cls, err := loader.Load(name)
	if err != nil {
		t.Fatalf("load %s: %v", name, err)
	}
	if aot {
		gen.Install(k) // emitted bodies take precedence
	}
	th := New(k)
	mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatalf("resolve main: %v", err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
	_, callErr := th.Call(mainM, nil, []kernel.Value{argsArr})
	if callErr == nil {
		t.Fatalf("%s completed without throwing", name)
	}
	thrown, ok := callErr.(*kernel.Thrown)
	if !ok {
		t.Fatalf("engine error, not a Java throwable: %v", callErr)
	}
	return kernel.FormatUncaught("main", thrown)
}

func assertFrames(t *testing.T, report string, wantFrames []string) {
	t.Helper()
	lines := strings.Split(strings.TrimSuffix(report, "\n"), "\n")
	if len(lines) != 1+len(wantFrames) {
		t.Fatalf("report has %d lines, want %d:\n%s", len(lines), 1+len(wantFrames), report)
	}
	for i, wf := range wantFrames {
		if lines[i+1] != "\tat "+wf+"(Unknown Source)" {
			t.Errorf("frame %d = %q, want %q\nreport:\n%s", i, lines[i+1], "\tat "+wf, report)
		}
	}
}

// TestStackBackfillBootstrapNPE pins DEBT-0019 infrastructure: a bootstrap
// NPE raised inside emitted/interpreted code carries the Java call stack,
// leaf first. Both engines route through kernel.InvokeAs, so both traces
// must agree.
func TestStackBackfillBootstrapNPE(t *testing.T) {
	for _, aot := range []bool{false, true} {
		report := runTraceProbe(t, "TraceProbe", aot)
		engine := map[bool]string{false: "interp", true: "aot"}[aot]
		if !strings.HasPrefix(report, "Exception in thread \"main\" java.lang.NullPointerException") {
			t.Errorf("[%s] header: %q", engine, report)
		}
		assertFrames(t, report, []string{
			"TraceProbe.stepA",
			"TraceProbe.main",
		})
	}
}

// TestStackBackfillUserException pins fillInStackTrace semantics for user
// exceptions: the <init> construction chain is trimmed off the top.
func TestStackBackfillUserException(t *testing.T) {
	for _, aot := range []bool{false, true} {
		report := runTraceProbe(t, "TraceProbe2", aot)
		engine := map[bool]string{false: "interp", true: "aot"}[aot]
		if !strings.HasPrefix(report, "Exception in thread \"main\" java.lang.IllegalStateException: boom:3") {
			t.Errorf("[%s] header: %q", engine, report)
		}
		if strings.Contains(report, "<init>") {
			t.Errorf("[%s] construction chain leaked into trace:\n%s", engine, report)
		}
		assertFrames(t, report, []string{
			"TraceProbe2.fire",
			"TraceProbe2.main",
		})
	}
}

// TestTwoParseDoubleParseBothPaths pins DEBT-0018: parsing the SAME
// document twice in one process must succeed twice (the historical bug
// returned null on parse#2). Interpreter and AOT paths both covered;
// JsonDriver rides along as the minimal-json end-to-end assault target.
func TestTwoParseDoubleParseBothPaths(t *testing.T) {
	want := strings.Join([]string{
		"parse#1 len=94",
		"parse#2 len=94",
		"",
	}, "\n")
	for _, aot := range []bool{false, true} {
		var out bytes.Buffer
		k := kernel.New(kernel.Options{Stdout: &out})
		loader := kernel.NewClassPathLoader(k, []string{"../../testdata/minjson"})
		k.SetResolver(loader.Load)
		cls, err := loader.Load("TwoParse")
		if err != nil {
			t.Fatalf("load TwoParse: %v", err)
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
			t.Fatalf("[aot=%v] TwoParse failed: %v", aot, err)
		}
		if out.String() != want {
			t.Errorf("[aot=%v] stdout =\n%q\nwant\n%q", aot, out.String(), want)
		}
	}
}

// TestAOTJsonDriverEndToEnd pins the minimal-json assault (DEBT-0019
// closure): all six probes plus nested output, byte-identical to the
// reference JVM oracle captured at fixture time.
func TestAOTJsonDriverEndToEnd(t *testing.T) {
	var out bytes.Buffer
	k := kernel.New(kernel.Options{Stdout: &out})
	loader := kernel.NewClassPathLoader(k, []string{"../../testdata/minjson"})
	k.SetResolver(loader.Load)
	cls, err := loader.Load("JsonDriver")
	if err != nil {
		t.Fatalf("load JsonDriver: %v", err)
	}
	gen.Install(k)

	th := New(k)
	mainM, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		t.Fatal(err)
	}
	argsArr, _ := k.NewArray("Ljava/lang/String;", 0)
	if _, err := th.Call(mainM, nil, []kernel.Value{argsArr}); err != nil {
		t.Fatalf("JsonDriver AOT failed: %v", err)
	}

	want := strings.Join([]string{
		"PROBE-OK len=7",
		"PROBE-OK len=18",
		"PROBE-OK len=7",
		"PROBE-OK len=14",
		"PROBE-OK len=26",
		"PROBE-OK len=94",
		"name=catty",
		"year=2026",
		"tags=3",
		"pi=3.5",
		"parse-err caught ok",
		"done",
		"",
	}, "\n")
	if out.String() != want {
		t.Errorf("stdout =\n%q\nwant\n%q", out.String(), want)
	}
}
