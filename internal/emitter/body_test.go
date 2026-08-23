package emitter

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"catty/internal/classfile"
)

// emitMainForTest parses a fixture class and returns the emitted Go body of
// its static main([Ljava/lang/String;)V.
func emitMainForTest(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")
	if m == nil {
		t.Fatalf("%s: no main method", path)
	}
	src, err := emitMethodBody(cf, m)
	if err != nil {
		t.Fatalf("emit main: %v", err)
	}
	return src
}

// TestHandlerEntryPushAfterLabel pins the handler-entry depth fix: at every
// exception handler entry the exception receiver push (`s0 = exc.Obj`) must
// be the first statement AFTER the handler label, because all arrivals come
// through excDispatch's goto — a push emitted before the label is dead code
// and the handler observes a stale operand slot. Handler-entry canonical
// depth is always 1 (JVMS §2.6): stack cleared, exception pushed.
func TestHandlerEntryPushAfterLabel(t *testing.T) {
	src := emitMainForTest(t, "../../testdata/cp/CollectionsDemo.class")
	lines := strings.Split(src, "\n")

	// Collect handler entry pcs from the class file itself.
	data, err := os.ReadFile("../../testdata/cp/CollectionsDemo.class")
	if err != nil {
		t.Fatal(err)
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	m := cf.FindMethod("main", "([Ljava/lang/String;)V")

	for _, h := range m.Code.Handlers {
		want := []string{"L" + strconv.Itoa(int(h.HandlerPc)) + ":", "s0 = exc.Obj"}
		if !containsConsecutive(lines, want) {
			t.Errorf("handler pc=%d: expected consecutive lines %q in emitted body:\n%s",
				h.HandlerPc, want, src)
		}
	}

	// The exception receiver must be loaded exactly once per distinct
	// handler pc and nowhere else (no dead pushes before labels).
	distinct := map[uint16]bool{}
	for _, h := range m.Code.Handlers {
		distinct[h.HandlerPc] = true
	}
	got := strings.Count(src, "= exc.Obj")
	if got != len(distinct) {
		t.Errorf("exc.Obj pushes = %d, want %d (one per distinct handler pc):\n%s",
			got, len(distinct), src)
	}
}

func containsConsecutive(lines []string, want []string) bool {
	for i := 0; i+len(want) <= len(lines); i++ {
		ok := true
		for j, w := range want {
			if strings.TrimSpace(lines[i+j]) != w {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
