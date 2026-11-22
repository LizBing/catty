package rt_test

import (
	"testing"

	"catty/rt"
)

// TestCattySwapAndGotoW exercises the interpreter paths for stack manipulation
// and wide branch opcodes that javac does not normally emit.
func TestCattySwapAndGotoW(t *testing.T) {
	newImage := func(method *rt.MethodDef) *rt.Image {
		return &rt.Image{
			ABI:     rt.ABIRevision,
			Profile: rt.ProfileID,
			Classes: []*rt.ClassDef{{
				Name: "OpcodeProbe", Super: "java/lang/Object", Flags: rt.AccPublic,
				Methods: []*rt.MethodDef{method},
			}},
		}
	}
	run := func(image *rt.Image) int32 {
		r := rt.NewRuntime(image)
		if err := r.Install(); err != nil {
			t.Fatalf("install: %v", err)
		}
		method, err := r.ResolveMethod("OpcodeProbe", "run", "()I")
		if err != nil {
			t.Fatal(err)
		}
		out := r.MainThread().Call(method, nil, nil)
		if out.Thrown != nil {
			t.Fatalf("run threw %v", out.Thrown)
		}
		return rt.AsInt(out.Value)
	}

	swapMethod := &rt.MethodDef{
		Owner: "OpcodeProbe", Name: "run", Descriptor: "()I", Flags: rt.AccStatic,
		Engine: "interp", MaxStack: 4, MaxLocals: 0, CodeSize: 5,
		Insns: []rt.Insn{
			{Op: rt.OpIconst1, PC: 0, Ref: -1, Len: 1},
			{Op: rt.OpIconst2, PC: 1, Ref: -1, Len: 1},
			{Op: rt.OpSwap, PC: 2, Ref: -1, Len: 1},
			{Op: rt.OpIsub, PC: 3, Ref: -1, Len: 1},
			{Op: rt.OpIreturn, PC: 4, Ref: -1, Len: 1},
		},
	}
	if got := run(newImage(swapMethod)); got != 1 {
		t.Fatalf("swap result = %d, want 1", got)
	}

	gotoWMethod := &rt.MethodDef{
		Owner: "OpcodeProbe", Name: "run", Descriptor: "()I", Flags: rt.AccStatic,
		Engine: "interp", MaxStack: 1, MaxLocals: 0, CodeSize: 9,
		Insns: []rt.Insn{
			{Op: rt.OpGotoW, PC: 0, A: 7, Ref: -1, Len: 5},
			{Op: rt.OpIconst5, PC: 5, Ref: -1, Len: 1},
			{Op: rt.OpIreturn, PC: 6, Ref: -1, Len: 1},
			{Op: rt.OpIconst0, PC: 7, Ref: -1, Len: 1},
			{Op: rt.OpBipush, PC: 8, A: 9, Ref: -1, Len: 2},
		},
	}
	// Add the final return at pc 10 so the taken path yields 9.
	gotoWMethod.Insns = append(gotoWMethod.Insns, rt.Insn{Op: rt.OpIreturn, PC: 10, Ref: -1, Len: 1})
	gotoWMethod.CodeSize = 11
	if got := run(newImage(gotoWMethod)); got != 9 {
		t.Fatalf("goto_w result = %d, want 9", got)
	}
}
