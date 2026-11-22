package verify

import (
	"errors"
	"strings"
	"testing"

	"catty/internal/classfile"
	"catty/rt"
)

// handlerFixture is a static method that always throws RuntimeException and
// catches it at pc 8. The handler's stack-map frame controls the locals the
// verifier must require on every exception edge.
func handlerFixture(handlerLocals []classfile.VType) *classfile.Class {
	m := &classfile.Method{
		Owner: "HandlerFixture", Name: "f", Descriptor: "(I)I", Flags: rt.AccPublic | rt.AccStatic,
		MaxStack: 2, MaxLocals: 2, CodeLength: 11,
		Insns: []rt.Insn{
			{Op: rt.OpNew, PC: 0, A: 0, Ref: 0, Len: 3},
			{Op: rt.OpDup, PC: 3, Ref: -1, Len: 1},
			{Op: rt.OpInvokespecial, PC: 4, Ref: 1, Len: 3},
			{Op: rt.OpAthrow, PC: 7, Ref: -1, Len: 1},
			{Op: rt.OpAstore, PC: 8, A: 1, Ref: -1, Len: 1},
			{Op: rt.OpIconst2, PC: 9, Ref: -1, Len: 1},
			{Op: rt.OpIreturn, PC: 10, Ref: -1, Len: 1},
		},
		Refs: []rt.Ref{
			{Kind: rt.RefClass, Class: "java/lang/RuntimeException", CP: 1},
			{Kind: rt.RefMethod, Owner: "java/lang/RuntimeException", Name: "<init>", Descriptor: "()V", CP: 2},
		},
		Handlers:    []rt.Handler{{StartPC: 0, EndPC: 8, HandlerPC: 8, CatchType: "java/lang/RuntimeException"}},
		HasStackMap: true,
		Frames: []classfile.Frame{{
			Offset: 8, Locals: handlerLocals,
			Stack: []classfile.VType{{Kind: classfile.VTObject, Class: "java/lang/RuntimeException"}},
		}},
	}
	return &classfile.Class{
		Name: "HandlerFixture", Super: "java/lang/Object", Flags: rt.AccPublic,
		Methods: []*classfile.Method{m},
		Major:   65,
	}
}

func TestHandlerFrameInputsAccepted(t *testing.T) {
	c := handlerFixture([]classfile.VType{{Kind: classfile.VTInteger}})
	if err := VerifyClass(c, &Context{Classes: map[string]*classfile.Class{}}); err != nil {
		t.Fatalf("valid handler frame rejected: %v", err)
	}
	// A shared handler target may be reached from several catch types; the
	// frame must still match the local state at every edge. A catch-all can
	// deliver any Throwable, so its frame declares the Throwable supertype.
	c.Methods[0].Frames[0].Stack = []classfile.VType{{Kind: classfile.VTObject, Class: "java/lang/Throwable"}}
	c.Methods[0].Handlers = append(c.Methods[0].Handlers,
		rt.Handler{StartPC: 0, EndPC: 8, HandlerPC: 8, CatchType: ""})
	if err := VerifyClass(c, &Context{Classes: map[string]*classfile.Class{}}); err != nil {
		t.Fatalf("shared catch-all handler rejected: %v", err)
	}
}

func TestHandlerFrameLocalMismatchRejected(t *testing.T) {
	c := handlerFixture([]classfile.VType{{Kind: classfile.VTFloat}})
	err := VerifyClass(c, &Context{Classes: map[string]*classfile.Class{}})
	if err == nil {
		t.Fatal("inconsistent handler frame accepted")
	}
	msg := err.Error()
	for _, want := range []string{"exception edge", "pc 0", "handler pc 8", "local 0"} {
		if !strings.Contains(msg, want) {
			t.Fatalf("handler diagnostic %q does not contain %q", msg, want)
		}
	}
}

func TestHandlerFrameUninitializedMismatchRejected(t *testing.T) {
	// new + dup leaves an uninitialized value on the stack; the handler frame
	// must not claim an initialized reference for it. This fixture uses a local
	// slot to carry the uninitialized value into the exception edge.
	m := handlerFixture([]classfile.VType{{Kind: classfile.VTInteger}}).Methods[0]
	m.MaxLocals = 3
	m.Insns = []rt.Insn{
		{Op: rt.OpNew, PC: 0, Ref: 0, Len: 3},
		{Op: rt.OpDup, PC: 3, Ref: -1, Len: 1},
		{Op: rt.OpAstore, PC: 4, A: 1, Ref: -1, Len: 1},
		{Op: rt.OpAthrow, PC: 5, Ref: -1, Len: 1},
		{Op: rt.OpAconstNull, PC: 6, Ref: -1, Len: 1},
		{Op: rt.OpAstore, PC: 7, A: 1, Ref: -1, Len: 1},
		{Op: rt.OpIconst2, PC: 8, Ref: -1, Len: 1},
		{Op: rt.OpIreturn, PC: 9, Ref: -1, Len: 1},
	}
	m.MaxStack = 1
	m.CodeLength = 10
	m.Handlers = []rt.Handler{{StartPC: 0, EndPC: 5, HandlerPC: 6, CatchType: ""}}
	m.Frames = []classfile.Frame{{
		Offset: 6,
		Locals: []classfile.VType{{Kind: classfile.VTInteger}, {Kind: classfile.VTObject, Class: "java/lang/Object"}},
		Stack:  []classfile.VType{{Kind: classfile.VTObject, Class: "java/lang/Throwable"}},
	}}
	c := &classfile.Class{Name: "HandlerFixture", Super: "java/lang/Object", Flags: rt.AccPublic, Methods: []*classfile.Method{m}}
	err := VerifyClass(c, &Context{Classes: map[string]*classfile.Class{}})
	if err == nil {
		t.Fatal("initialized local frame matched an uninitialized exception edge")
	}
	if !strings.Contains(err.Error(), "exception edge") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

func TestVerifierWorkExactBoundary(t *testing.T) {
	c := handlerFixture([]classfile.VType{{Kind: classfile.VTInteger}})
	classes := map[string]*classfile.Class{}
	work := &Work{Limit: 1 << 20}
	if err := VerifyClass(c, &Context{Classes: classes, Work: work}); err != nil {
		t.Fatal(err)
	}
	if work.Used == 0 {
		t.Fatal("verifier did not charge work")
	}
	exact := &Work{Limit: work.Used}
	if err := VerifyClass(c, &Context{Classes: classes, Work: exact}); err != nil {
		t.Fatalf("exact verifier work boundary rejected: %v", err)
	}
	short := &Work{Limit: work.Used - 1}
	err := VerifyClass(c, &Context{Classes: classes, Work: short})
	if err == nil {
		t.Fatal("one unit below the verifier work boundary accepted")
	}
	var workErr *WorkError
	if !errors.As(err, &workErr) {
		t.Fatalf("expected WorkError, got %T: %v", err, err)
	}
}
