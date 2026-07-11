package classfile

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// compileFixture compiles tests/fixtures/<name>.java into t.TempDir() and
// returns the bytes of the resulting <name>.class. It fails the test (rather
// than skipping) if javac is missing — the whole project depends on a JDK.
func compileFixture(t *testing.T, name string) []byte {
	t.Helper()
	src := filepath.Join("..", "tests", "fixtures", name+".java")
	out := t.TempDir()
	cmd := exec.Command("javac", "-source", "8", "-target", "8", "-d", out, src)
	if err := cmd.Run(); err != nil {
		t.Fatalf("javac %s failed: %v", name, err)
	}
	data, err := os.ReadFile(filepath.Join(out, name+".class"))
	if err != nil {
		t.Fatalf("read class: %v", err)
	}
	return data
}

func TestParseHelloWorld(t *testing.T) {
	cf, err := Parse(compileFixture(t, "HelloWorld"))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got := cf.ClassName(); got != "HelloWorld" {
		t.Errorf("ClassName = %q, want HelloWorld", got)
	}
	if got := cf.SuperClassName(); got != "java/lang/Object" {
		t.Errorf("SuperClassName = %q, want java/lang/Object", got)
	}

	// HelloWorld declares no fields and exactly 2 methods: <init>() and main.
	if len(cf.Fields()) != 0 {
		t.Errorf("Fields = %d, want 0", len(cf.Fields()))
	}
	if len(cf.Methods()) != 2 {
		t.Fatalf("Methods = %d, want 2", len(cf.Methods()))
	}

	var main *MemberInfo
	for _, m := range cf.Methods() {
		if m.Name() == "main" && m.Descriptor() == "([Ljava/lang/String;)V" {
			main = m
		}
	}
	if main == nil {
		t.Fatalf("main method not found")
	}
	code := main.Code()
	if code == nil {
		t.Fatalf("main has no Code attribute")
	}
	if len(code.Code()) == 0 {
		t.Errorf("main Code is empty")
	}
	if code.MaxLocals() < 3 { // args + a + b
		t.Errorf("MaxLocals = %d, want >= 3", code.MaxLocals())
	}
}
