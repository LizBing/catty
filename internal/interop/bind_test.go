package interop

import (
	"errors"
	"math"
	"strings"
	"testing"

	"catty/internal/kernel"
)

type discard struct{}

func (discard) Write(p []byte) (int, error) { return len(p), nil }

func newK(t *testing.T) *kernel.Kernel {
	t.Helper()
	return kernel.New(kernel.Options{Stdout: discard{}, Stderr: discard{}})
}

func mustBind(t *testing.T, k *kernel.Kernel, spec Spec) {
	t.Helper()
	if err := Bind(k, spec); err != nil {
		t.Fatalf("Bind %s: %v", spec.Class, err)
	}
}

// call runs a bound static through the standard dispatch path and returns
// the raw result; a thrown Java exception comes back as error.
func call(t *testing.T, k *kernel.Kernel, cls, name string, args ...kernel.Value) (kernel.Value, error) {
	t.Helper()
	c, ok := k.ClassByName(cls)
	if !ok {
		t.Fatalf("class %s not registered", cls)
	}
	var m *kernel.Method
	for _, cand := range c.Methods {
		if cand.Name == name {
			m = cand
			break
		}
	}
	if m == nil {
		t.Fatalf("method %s.%s not found", cls, name)
	}
	return k.InvokeAs(nil, m, nil, args)
}

func wantThrow(t *testing.T, err error, msgPart string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected throw containing %q, got none", msgPart)
	}
	if !strings.Contains(err.Error(), msgPart) {
		t.Fatalf("throw %q missing %q", err.Error(), msgPart)
	}
}

func TestBindStringRoundTrip(t *testing.T) {
	k := newK(t)
	mustBind(t, k, Spec{Class: "go/T", Funcs: map[string]any{
		"echo":  func(s string) string { return "<<" + s + ">>" },
		"empty": func(s string) int32 { return int32(len(s)) },
	}})
	v, err := call(t, k, "go/T", "echo", k.InternGo("catty"))
	if err != nil {
		t.Fatal(err)
	}
	if got := v.(*kernel.JString).Go(); got != "<<catty>>" {
		t.Errorf("got %q", got)
	}
	// empty string survives the boundary
	v, err = call(t, k, "go/T", "empty", k.InternGo(""))
	if err != nil {
		t.Fatal(err)
	}
	if v.(int32) != 0 {
		t.Errorf("empty string length = %v", v)
	}
}

func TestBindNumericAndBool(t *testing.T) {
	k := newK(t)
	mustBind(t, k, Spec{Class: "go/N", Funcs: map[string]any{
		"i32":   func(v int32) int32 { return v * 2 },
		"i64":   func(v int64) int64 { return v * 2 },
		"gi":    func(v int) int64 { return int64(v) * 2 },
		"f32":   func(v float32) float32 { return v / 2 },
		"f64":   func(v float64) float64 { return v / 2 },
		"neg":   func(v bool) bool { return !v },
		"nanIn": func(v float64) bool { return math.IsNaN(v) },
	}})
	step := func(name string, in kernel.Value) kernel.Value {
		v, err := call(t, k, "go/N", name, in)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		return v
	}
	if v := step("i32", int32(21)); v.(int32) != 42 {
		t.Errorf("i32 = %v", v)
	}
	if v := step("i64", int64(1<<40)); v.(int64) != 1<<41 {
		t.Errorf("i64 = %v", v)
	}
	if v := step("gi", int64(21)); v.(int64) != 42 {
		t.Errorf("gi = %v", v)
	}
	if v := step("f32", float32(3.5)); v.(float32) != 1.75 {
		t.Errorf("f32 = %v", v)
	}
	if v := step("f64", 3.5); v.(float64) != 1.75 {
		t.Errorf("f64 = %v", v)
	}
	if v := step("neg", int32(0)); v.(int32) != 1 {
		t.Errorf("neg(false) = %v", v)
	}
	// NaN and ±Inf pass through the D boundary intact.
	if v := step("nanIn", math.NaN()); v.(int32) != 1 {
		t.Errorf("NaN not observed: %v", v)
	}
	if _, err := call(t, k, "go/N", "f64", math.Inf(-1)); err != nil {
		t.Errorf("-Inf rejected: %v", err)
	}
}

func TestBindValuePassthroughIdentity(t *testing.T) {
	k := newK(t)
	mustBind(t, k, Spec{Class: "go/P", Funcs: map[string]any{
		"same": func(v kernel.Value) kernel.Value { return v },
	}})
	cls, _ := k.ClassByName("java/lang/Object")
	obj, err := k.NewInstance(cls)
	if err != nil {
		t.Fatal(err)
	}
	got, cerr := call(t, k, "go/P", "same", obj)
	if cerr != nil {
		t.Fatal(cerr)
	}
	if got != kernel.Value(obj) {
		t.Fatal("passthrough broke identity")
	}
	// null passthrough is legal for Value-typed parameters
	if _, cerr := call(t, k, "go/P", "same", nil); cerr != nil {
		t.Fatalf("null Value rejected: %v", cerr)
	}
}

func TestBindNullScalarThrows(t *testing.T) {
	k := newK(t)
	mustBind(t, k, Spec{Class: "go/X", Funcs: map[string]any{
		"slen": func(s string) int32 { return int32(len(s)) },
	}})
	_, err := call(t, k, "go/X", "slen", nil)
	wantThrow(t, err, "interop: null argument 0")
}

func TestBindErrorPropagation(t *testing.T) {
	k := newK(t)
	sentinel := errors.New("disk on fire")
	mustBind(t, k, Spec{Class: "go/E", Funcs: map[string]any{
		"boom": func() error { return sentinel },
		"ok":   func() error { return nil },
		"both": func(s string) (string, error) {
			if s == "bad" {
				return "", errors.New("bad input")
			}
			return strings.ToUpper(s), nil
		},
	}})
	_, err := call(t, k, "go/E", "boom")
	wantThrow(t, err, "disk on fire")

	if _, err := call(t, k, "go/E", "ok"); err != nil {
		t.Fatalf("ok() threw: %v", err)
	}

	v, err := call(t, k, "go/E", "both", k.InternGo("quiet"))
	if err != nil {
		t.Fatal(err)
	}
	if got := v.(*kernel.JString).Go(); got != "QUIET" {
		t.Errorf("got %q", got)
	}
	_, err = call(t, k, "go/E", "both", k.InternGo("bad"))
	wantThrow(t, err, "bad input")
}

func TestBindFailFastAtRegistration(t *testing.T) {
	k := newK(t)
	cases := []struct {
		name string
		fn   any
		want string
	}{
		{"notafunc", 42, "want func"},
		{"unsupported", func(v chan int) string { return "" }, "unsupported type"},
		{"unsupresult", func() chan int { return nil }, "unsupported type"},
		{"twovalues", func() (string, int32) { return "", 0 }, "at most one"},
	}
	for _, c := range cases {
		err := Bind(k, Spec{Class: "go/F", Funcs: map[string]any{c.name: c.fn}})
		if err == nil {
			t.Errorf("%s: Bind succeeded, want fail-fast", c.name)
		} else if !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %q missing %q", c.name, err.Error(), c.want)
		}
	}
	// The failed Bind must not have left a half-registered class behind.
	if _, ok := k.ClassByName("go/F"); ok {
		t.Error("failed Bind registered its class anyway")
	}
}

func TestBindVariadicRejected(t *testing.T) {
	k := newK(t)
	err := Bind(k, Spec{Class: "go/V", Funcs: map[string]any{
		"variadic": func(xs ...int32) int32 { return 0 },
	}})
	wantThrow(t, err, "variadic") // Bind returns plain errors, check text via wrapper below
}
