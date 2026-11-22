package classfile

import (
	"fmt"

	"catty/rt"
)

func s8(v byte) int32    { return int32(int8(v)) }
func s16(v uint16) int32 { return int32(int16(v)) }

// decodeInstructions validates operand lengths and produces retained,
// offset-preserving instructions plus a per-method symbolic reference table.
func decodeInstructions(cp *constantPool, code []byte, m *Method) ([]rt.Insn, []rt.Ref, error) {
	var insns []rt.Insn
	var refs []rt.Ref
	r := &reader{data: code}
	addRef := func(ref rt.Ref) int32 {
		refs = append(refs, ref)
		return int32(len(refs) - 1)
	}
	fail := func(pc int, format string, args ...any) ([]rt.Insn, []rt.Ref, error) {
		return nil, nil, fmt.Errorf("pc %d: %s", pc, fmt.Sprintf(format, args...))
	}
	for r.off < len(code) {
		pc := r.off
		op, _ := r.u1()
		insn := rt.Insn{Op: op, PC: int32(pc), Ref: -1}
		switch op {
		case rt.OpAconstNull, rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2,
			rt.OpIconst3, rt.OpIconst4, rt.OpIconst5, rt.OpPop, rt.OpDup,
			rt.OpIadd, rt.OpIrem, rt.OpIreturn, rt.OpAreturn, rt.OpReturn,
			rt.OpAaload, rt.OpAastore, rt.OpArraylength, rt.OpAthrow:
		case 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d,
			0x3b, 0x3c, 0x3d, 0x3e, 0x4b, 0x4c, 0x4d, 0x4e:
			if idx, ok := localIndexClassfile(op); ok {
				insn.A = idx
			}
		case rt.OpBipush:
			v, err := r.u1()
			if err != nil {
				return fail(pc, "truncated bipush")
			}
			insn.A = s8(v)
		case rt.OpSipush:
			v, err := r.u2()
			if err != nil {
				return fail(pc, "truncated sipush")
			}
			insn.A = s16(v)
		case rt.OpLdc, rt.OpLdcW:
			var idx uint16
			if op == rt.OpLdc {
				v, err := r.u1()
				if err != nil {
					return fail(pc, "truncated ldc")
				}
				idx = uint16(v)
			} else {
				v, err := r.u2()
				if err != nil {
					return fail(pc, "truncated ldc_w")
				}
				idx = v
			}
			ref, err := ldcRef(cp, idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			insn.Ref = addRef(ref)
		case rt.OpIload, rt.OpAload, rt.OpIstore, rt.OpAstore:
			v, err := r.u1()
			if err != nil {
				return fail(pc, "truncated local instruction")
			}
			insn.A = int32(v)
		case rt.OpIfeq, rt.OpIfne, rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfAcmpeq, rt.OpIfAcmpne, rt.OpGoto:
			v, err := r.u2()
			if err != nil {
				return fail(pc, "truncated branch")
			}
			insn.A = int32(pc) + s16(v)
		case rt.OpGetstatic, rt.OpPutstatic, rt.OpGetfield, rt.OpPutfield:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated field reference")
			}
			owner, name, desc, kind, err := cp.memberRef(idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			if kind != "field" {
				return fail(pc, "opcode 0x%02x references a non-field constant", op)
			}
			insn.Ref = addRef(rt.Ref{Kind: rt.RefField, Owner: owner, Name: name, Descriptor: desc, CP: int(idx)})
		case rt.OpInvokevirtual, rt.OpInvokespecial, rt.OpInvokestatic, rt.OpInvokeinterface:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated method reference")
			}
			owner, name, desc, kind, err := cp.memberRef(idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			if kind == "field" {
				return fail(pc, "invoke references a field constant")
			}
			if op == rt.OpInvokeinterface {
				count, err := r.u1()
				if err != nil {
					return fail(pc, "truncated invokeinterface count")
				}
				zero, err := r.u1()
				if err != nil {
					return fail(pc, "truncated invokeinterface zero byte")
				}
				if zero != 0 {
					return fail(pc, "invokeinterface has nonzero reserved byte")
				}
				insn.B = int32(count)
			}
			insn.Ref = addRef(rt.Ref{Kind: rt.RefMethod, Owner: owner, Name: name, Descriptor: desc, CP: int(idx)})
		case rt.OpNew, rt.OpAnewarray:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated class reference")
			}
			name, err := cp.className(idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			insn.Ref = addRef(rt.Ref{Kind: rt.RefClass, Class: name, CP: int(idx)})
		case rt.OpInvokedynamic:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated invokedynamic")
			}
			if _, err := r.u2(); err != nil {
				return fail(pc, "truncated invokedynamic")
			}
			insn.A = int32(idx)
		default:
			return fail(pc, "unsupported opcode 0x%02x", op)
		}
		insn.Len = int32(r.off - pc)
		insns = append(insns, insn)
	}
	if len(insns) == 0 {
		return nil, nil, fmt.Errorf("empty Code array")
	}
	return insns, refs, nil
}

func localIndexClassfile(op byte) (int32, bool) {
	switch {
	case op >= 0x1a && op <= 0x1d:
		return int32(op - 0x1a), true
	case op >= 0x2a && op <= 0x2d:
		return int32(op - 0x2a), true
	case op >= 0x3b && op <= 0x3e:
		return int32(op - 0x3b), true
	case op >= 0x4b && op <= 0x4e:
		return int32(op - 0x4b), true
	}
	return 0, false
}

func ldcRef(cp *constantPool, idx uint16) (rt.Ref, error) {
	e, err := cp.entry(idx)
	if err != nil {
		return rt.Ref{}, err
	}
	switch e.tag {
	case 3:
		return rt.Ref{Kind: rt.RefInt, Int: e.intVal, CP: int(idx)}, nil
	case 7:
		name, err := cp.className(idx)
		if err != nil {
			return rt.Ref{}, err
		}
		return rt.Ref{Kind: rt.RefClass, Class: name, CP: int(idx)}, nil
	case 8:
		units, err := cp.ldcString(idx)
		if err != nil {
			return rt.Ref{}, err
		}
		return rt.Ref{Kind: rt.RefString, Units: units, CP: int(idx)}, nil
	default:
		return rt.Ref{}, fmt.Errorf("unsupported ldc constant tag %d", e.tag)
	}
}

// BranchTargets returns the relative branch target for one instruction.
func BranchTarget(insn rt.Insn) (int32, bool) {
	switch insn.Op {
	case rt.OpIfeq, rt.OpIfne, rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfAcmpeq, rt.OpIfAcmpne, rt.OpGoto:
		return insn.A, true
	}
	return 0, false
}

// CodeLengthFromPCs returns the total byte length implied by decoded offsets.
func CodeLengthFromPCs(insns []rt.Insn) int {
	if len(insns) == 0 {
		return 0
	}
	last := insns[len(insns)-1]
	return int(last.PC + last.Len)
}
