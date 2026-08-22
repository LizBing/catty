package classfile

import (
	"os"
	"testing"
)

const testdataDir = "../../testdata"

func mustParseFixture(t *testing.T, name string) *ClassFile {
	t.Helper()
	data, err := os.ReadFile(testdataDir + "/" + name)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	cf, err := Parse(data)
	if err != nil {
		t.Fatalf("parse %s: %v", name, err)
	}
	return cf
}

func TestParseHelloWorld(t *testing.T) {
	cf := mustParseFixture(t, "HelloWorld.class")

	if cf.MajorVersion != 52 {
		t.Errorf("major = %d, want 52 (javac --release 8)", cf.MajorVersion)
	}
	if cf.ThisClass != "HelloWorld" {
		t.Errorf("ThisClass = %q, want HelloWorld", cf.ThisClass)
	}
	if cf.SuperClass != "java/lang/Object" {
		t.Errorf("SuperClass = %q, want java/lang/Object", cf.SuperClass)
	}

	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	if m == nil {
		t.Fatal("main method not found")
	}
	if m.Code == nil || len(m.Code.Code) == 0 {
		t.Fatal("main has no Code attribute")
	}
	// javac --release 8 compiles println(String) via getstatic(0xB2).
	if m.Code.Code[0] != 0xB2 {
		t.Errorf("first opcode of main = %#x, want 0xB2 (getstatic)", m.Code.Code[0])
	}
	if m.Code.MaxLocals < 1 {
		t.Errorf("main MaxLocals = %d, want >= 1 (args param)", m.Code.MaxLocals)
	}

	ctor := cf.FindMethod("<init>", "()V")
	if ctor == nil || ctor.Code == nil {
		t.Error("default constructor missing or has no Code")
	}
}

func TestParseCollectionsDemo(t *testing.T) {
	cf := mustParseFixture(t, "CollectionsDemo.class")

	if cf.ThisClass != "CollectionsDemo" {
		t.Fatalf("ThisClass = %q", cf.ThisClass)
	}

	main := cf.FindMethod("main", "([Ljava/lang/String;)V")
	if main == nil {
		t.Fatal("main not found")
	}

	// The demo exercises: invokevirtual on ArrayList/PrintStream/StringBuilder,
	// checkcast to Integer, exception handlers for IndexOutOfBounds and Arithmetic.
	sawCheckcast, sawHandlers := false, false
	for _, b := range main.Code.Code {
		if b == 0xC0 { // checkcast
			sawCheckcast = true
		}
	}
	if len(main.Code.Handlers) < 2 {
		sawHandlers = true // inverted below
	}
	if !sawCheckcast {
		t.Error("expected at least one checkcast in main (unboxing xs.get)")
	}
	if len(main.Code.Handlers) != 2 {
		t.Errorf("handlers = %d, want 2 (two try/catch blocks)", len(main.Code.Handlers))
	}
	if sawHandlers {
		t.Error("handler count assertion logic error") // keeps the inverted var honest
	}
	for _, h := range main.Code.Handlers {
		if h.EndPc <= h.StartPc {
			t.Errorf("bad handler range [%d,%d)", h.StartPc, h.EndPc)
		}
	}
}

func TestParseRejects(t *testing.T) {
	cases := []struct {
		name string
		data []byte
	}{
		{"empty", []byte{}},
		{"bad-magic", []byte{0xCA, 0xFE, 0xBA, 0xBD, 0, 0, 0, 52}},
		{"version-53", func() []byte {
			b := []byte{0xCA, 0xFE, 0xBA, 0xBE}
			b = append(b, 0, 0, 0, 53) // minor=0 major=53 > 52
			return append(b, make([]byte, 16)...)
		}()},
		{"truncated-hello", func() []byte {
			full, err := os.ReadFile(testdataDir + "/HelloWorld.class")
			if err != nil {
				t.Fatal(err)
			}
			return full[:len(full)/2]
		}()},
		{"trailing-garbage", func() []byte {
			full, err := os.ReadFile(testdataDir + "/HelloWorld.class")
			if err != nil {
				t.Fatal(err)
			}
			return append(append([]byte{}, full...), 0xDE, 0xAD)
		}()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.data); err == nil {
				t.Errorf("expected error for %s", tc.name)
			}
		})
	}
}

func TestMUTF8RoundTrip(t *testing.T) {
	cases := []string{
		"", "ascii only", "nul\x00inside", "中文测试",
		"emoji \U0001F600 supplementary", "mixéß",
	}
	for _, s := range cases {
		got := decodeMUTF8(EncodeMUTF8(s))
		if got != s {
			t.Errorf("roundtrip %q → %q", s, got)
		}
	}
	// Explicit C0 80 encoding decodes to NUL.
	if got := decodeMUTF8([]byte{0xC0, 0x80}); got != "\x00" {
		t.Errorf("C0 80 → %q, want NUL", got)
	}
}
