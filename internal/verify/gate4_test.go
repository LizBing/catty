package verify

import (
	"strings"
	"testing"

	"catty/internal/classfile"
	"catty/rt"
)

func gate4CodeMethod(name, descriptor string, flags uint16, insns []rt.Insn) *classfile.Method {
	codeLength := 0
	if len(insns) > 0 {
		last := insns[len(insns)-1]
		codeLength = int(last.PC + last.Len)
	}
	return &classfile.Method{
		Owner: "Fixture", Name: name, Descriptor: descriptor, Flags: flags,
		MaxStack: 2, MaxLocals: 1, CodeLength: codeLength, Insns: insns,
	}
}

func gate4Class(methods ...*classfile.Method) *classfile.Class {
	return &classfile.Class{
		Name: "Fixture", Super: "java/lang/Object", Flags: rt.AccPublic,
		Methods: methods, Major: 65,
	}
}

func verifyGate4(c *classfile.Class) error {
	return VerifyClass(c, &Context{Classes: map[string]*classfile.Class{}, Profile: rt.Gate4ProfileID})
}

func insn(op uint8, pc, length int32) rt.Insn {
	return rt.Insn{Op: op, PC: pc, Len: length, Ref: -1}
}

func TestGate4MonitorInstructionsVerify(t *testing.T) {
	// null is a valid reference operand for both monitor instructions.
	body := []rt.Insn{
		insn(rt.OpAconstNull, 0, 1),
		insn(rt.OpMonitorenter, 1, 1),
		insn(rt.OpAconstNull, 2, 1),
		insn(rt.OpMonitorexit, 3, 1),
		insn(rt.OpReturn, 4, 1),
	}
	if err := verifyGate4(gate4Class(gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccStatic, body))); err != nil {
		t.Fatalf("valid monitor body rejected: %v", err)
	}
	// A monitorenter with an int operand matches OpenJDK's VerifyError shape.
	badEnter := []rt.Insn{
		insn(rt.OpIconst0, 0, 1),
		insn(rt.OpMonitorenter, 1, 1),
		insn(rt.OpReturn, 2, 1),
	}
	err := verifyGate4(gate4Class(gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccStatic, badEnter)))
	if err == nil || !strings.Contains(err.Error(), "pc 1") ||
		!strings.Contains(err.Error(), "expected a reference") {
		t.Fatalf("invalid monitorenter operand error = %v", err)
	}
	// monitorexit has the same operand rule.
	badExit := []rt.Insn{
		insn(rt.OpIconst0, 0, 1),
		insn(rt.OpMonitorexit, 1, 1),
		insn(rt.OpReturn, 2, 1),
	}
	if err := verifyGate4(gate4Class(gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccStatic, badExit))); err == nil ||
		!strings.Contains(err.Error(), "expected a reference") {
		t.Fatalf("invalid monitorexit operand error = %v", err)
	}
	// Earlier profiles keep rejecting both opcodes at the same offset.
	for _, profile := range []string{rt.Gate2ProfileID, rt.Gate3ProfileID} {
		ctx := &Context{Classes: map[string]*classfile.Class{}, Profile: profile}
		err := VerifyClass(gate4Class(gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccStatic, body)), ctx)
		if err == nil || !strings.Contains(err.Error(), "monitor instruction") {
			t.Fatalf("%s monitor body error = %v", profile, err)
		}
	}
}

func TestGate4SynchronizedMethodFlags(t *testing.T) {
	returnOnly := []rt.Insn{insn(rt.OpReturn, 0, 1)}
	valid := gate4Class(gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccStatic|rt.AccSynchronized, returnOnly))
	if err := verifyGate4(valid); err != nil {
		t.Fatalf("valid synchronized static method rejected: %v", err)
	}
	// The synchronization flag on <clinit> is ignored for execution.
	clinit := gate4Class(gate4CodeMethod("<clinit>", "()V", rt.AccStatic|rt.AccSynchronized, returnOnly))
	if err := verifyGate4(clinit); err != nil {
		t.Fatalf("<clinit> synchronization flag rejected: %v", err)
	}
	// A class initializer must still be static ()V.
	if err := verifyGate4(gate4Class(gate4CodeMethod("<clinit>", "()V", rt.AccSynchronized, returnOnly))); err == nil ||
		!strings.Contains(err.Error(), "static") {
		t.Fatalf("non-static class initializer error = %v", err)
	}

	constructor := gate4Class(gate4CodeMethod("<init>", "()V", rt.AccPublic|rt.AccSynchronized, returnOnly))
	if err := verifyGate4(constructor); err == nil || !strings.Contains(err.Error(), "constructor") {
		t.Fatalf("synchronized constructor error = %v", err)
	}
	abstract := gate4Class(gate4CodeMethod("a", "()V", rt.AccPublic|rt.AccAbstract|rt.AccSynchronized, nil))
	if err := verifyGate4(abstract); err == nil || !strings.Contains(err.Error(), "abstract") {
		t.Fatalf("synchronized abstract method error = %v", err)
	}
	iface := &classfile.Class{
		Name: "Fixture", Super: "java/lang/Object", Flags: rt.AccPublic | rt.AccInterface | rt.AccAbstract,
		Methods: []*classfile.Method{gate4CodeMethod("m", "()V", rt.AccPublic|rt.AccSynchronized, returnOnly)},
		Major:   65,
	}
	if err := verifyGate4(iface); err == nil || !strings.Contains(err.Error(), "interface method") {
		t.Fatalf("synchronized interface method error = %v", err)
	}
	native := gate4Class(gate4CodeMethod("n", "()V", rt.AccPublic|rt.AccNative|rt.AccSynchronized, nil))
	if err := verifyGate4(native); err == nil || !strings.Contains(err.Error(), "native") {
		t.Fatalf("synchronized native method error = %v", err)
	}
}

func TestGate4FinalVolatileField(t *testing.T) {
	fieldClass := func(flags uint16) *classfile.Class {
		return &classfile.Class{
			Name: "Fixture", Super: "java/lang/Object", Flags: rt.AccPublic, Major: 65,
			Fields: []*classfile.Field{{Owner: "Fixture", Name: "f", Descriptor: "I", Flags: flags}},
		}
	}
	if err := verifyGate4(fieldClass(rt.AccFinal | rt.AccVolatile)); err == nil ||
		!strings.Contains(err.Error(), "final") {
		t.Fatalf("final+volatile field error = %v", err)
	}
	if err := verifyGate4(fieldClass(rt.AccFinal)); err != nil {
		t.Fatalf("final field rejected: %v", err)
	}
	ctx := &Context{Classes: map[string]*classfile.Class{}, Profile: rt.Gate3ProfileID}
	if err := VerifyClass(fieldClass(rt.AccFinal|rt.AccVolatile), ctx); err != nil {
		t.Fatalf("Gate 3 final+volatile field rejected: %v", err)
	}
}
