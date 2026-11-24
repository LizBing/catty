package classfile

import (
	"fmt"

	"catty/rt"
)

func s8(v byte) int32    { return int32(int8(v)) }
func s16(v uint16) int32 { return int32(int16(v)) }

// decodeInstructions validates operand lengths and produces retained,
// offset-preserving instructions plus a per-method symbolic reference table.
func decodeInstructions(cp *constantPool, code []byte, m *Method, profile string) ([]rt.Insn, []rt.Ref, error) {
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
		case rt.OpNop, rt.OpAconstNull,
			rt.OpIconstM1, rt.OpIconst0, rt.OpIconst1, rt.OpIconst2,
			rt.OpIconst3, rt.OpIconst4, rt.OpIconst5,
			rt.OpLconst0, rt.OpLconst1, rt.OpFconst0, rt.OpFconst1, rt.OpFconst2,
			rt.OpDconst0, rt.OpDconst1,
			rt.OpIaload, rt.OpLaload, rt.OpFaload, rt.OpDaload, rt.OpAaload,
			rt.OpBaload, rt.OpCaload, rt.OpSaload,
			rt.OpIastore, rt.OpLastore, rt.OpFastore, rt.OpDastore, rt.OpAastore,
			rt.OpBastore, rt.OpCastore, rt.OpSastore,
			rt.OpPop, rt.OpPop2, rt.OpDup, rt.OpDupX1, rt.OpDupX2,
			rt.OpDup2, rt.OpDup2X1, rt.OpDup2X2, rt.OpSwap,
			rt.OpIadd, rt.OpLadd, rt.OpFadd, rt.OpDadd,
			rt.OpIsub, rt.OpLsub, rt.OpFsub, rt.OpDsub,
			rt.OpImul, rt.OpLmul, rt.OpFmul, rt.OpDmul,
			rt.OpIdiv, rt.OpLdiv, rt.OpFdiv, rt.OpDdiv,
			rt.OpIrem, rt.OpLrem, rt.OpFrem, rt.OpDrem,
			rt.OpIneg, rt.OpLneg, rt.OpFneg, rt.OpDneg,
			rt.OpIshl, rt.OpLshl, rt.OpIshr, rt.OpLshr, rt.OpIushr, rt.OpLushr,
			rt.OpIand, rt.OpLand, rt.OpIor, rt.OpLor, rt.OpIxor, rt.OpLxor,
			rt.OpI2l, rt.OpI2f, rt.OpI2d, rt.OpL2i, rt.OpL2f, rt.OpL2d,
			rt.OpF2i, rt.OpF2l, rt.OpF2d, rt.OpD2i, rt.OpD2l, rt.OpD2f,
			rt.OpI2b, rt.OpI2c, rt.OpI2s,
			rt.OpLcmp, rt.OpFcmpl, rt.OpFcmpg, rt.OpDcmpl, rt.OpDcmpg,
			rt.OpIreturn, rt.OpLreturn, rt.OpFreturn, rt.OpDreturn,
			rt.OpAreturn, rt.OpReturn,
			rt.OpArraylength, rt.OpAthrow:
		case 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21,
			0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29,
			0x2a, 0x2b, 0x2c, 0x2d,
			0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40, 0x41, 0x42,
			0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a,
			0x4b, 0x4c, 0x4d, 0x4e:
			// Normalize the one-byte local-variable forms to their base
			// opcode with an explicit local index. The original offset and
			// length are preserved.
			base, idx := canonicalLocalShort(op)
			insn.Op = base
			insn.A = idx
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
		case rt.OpLdc, rt.OpLdcW, rt.OpLdc2W:
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
					return fail(pc, "truncated ldc_w/ldc2_w")
				}
				idx = v
			}
			ref, err := ldcRef(cp, idx, op == rt.OpLdc2W)
			if err != nil {
				return fail(pc, "%v", err)
			}
			insn.Ref = addRef(ref)
		case rt.OpIload, rt.OpLload, rt.OpFload, rt.OpDload, rt.OpAload,
			rt.OpIstore, rt.OpLstore, rt.OpFstore, rt.OpDstore, rt.OpAstore:
			v, err := r.u1()
			if err != nil {
				return fail(pc, "truncated local instruction")
			}
			insn.A = int32(v)
		case rt.OpIinc:
			idx, err := r.u1()
			if err != nil {
				return fail(pc, "truncated iinc")
			}
			amount, err := r.u1()
			if err != nil {
				return fail(pc, "truncated iinc")
			}
			insn.A = int32(idx)
			insn.B = s8(amount)
		case rt.OpIfeq, rt.OpIfne, rt.OpIflt, rt.OpIfge, rt.OpIfgt, rt.OpIfle,
			rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfIcmplt, rt.OpIfIcmpge,
			rt.OpIfIcmpgt, rt.OpIfIcmple, rt.OpIfAcmpeq, rt.OpIfAcmpne,
			rt.OpGoto, rt.OpIfnull, rt.OpIfnonnull:
			v, err := r.u2()
			if err != nil {
				return fail(pc, "truncated branch")
			}
			insn.A = int32(pc) + s16(v)
		case rt.OpGotoW:
			v, err := r.u4()
			if err != nil {
				return fail(pc, "truncated goto_w")
			}
			insn.A = int32(int64(pc) + int64(int32(v)))
		case rt.OpTableswitch, rt.OpLookupswitch:
			if err := decodeSwitch(cp, r, pc, op, &insn); err != nil {
				return fail(pc, "%v", err)
			}
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
		case rt.OpNew, rt.OpAnewarray, rt.OpCheckcast, rt.OpInstanceof:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated class reference")
			}
			name, err := cp.className(idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			insn.Ref = addRef(rt.Ref{Kind: rt.RefClass, Class: name, CP: int(idx)})
		case rt.OpNewarray:
			atype, err := r.u1()
			if err != nil {
				return fail(pc, "truncated newarray")
			}
			if _, ok := rt.NewarrayDescriptor(int32(atype)); !ok {
				return fail(pc, "invalid newarray atype %d", atype)
			}
			insn.A = int32(atype)
		case rt.OpMultianewarray:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated multianewarray")
			}
			dims, err := r.u1()
			if err != nil {
				return fail(pc, "truncated multianewarray dimensions")
			}
			if dims == 0 {
				return fail(pc, "multianewarray has zero dimensions")
			}
			name, err := cp.className(idx)
			if err != nil {
				return fail(pc, "%v", err)
			}
			insn.Ref = addRef(rt.Ref{Kind: rt.RefClass, Class: name, CP: int(idx)})
			insn.B = int32(dims)
		case rt.OpInvokedynamic:
			idx, err := r.u2()
			if err != nil {
				return fail(pc, "truncated invokedynamic")
			}
			if _, err := r.u2(); err != nil {
				return fail(pc, "truncated invokedynamic")
			}
			insn.A = int32(idx)
		case rt.OpWide:
			if err := decodeWide(r, pc, &insn); err != nil {
				return fail(pc, "%v", err)
			}
		case rt.OpJsr, rt.OpRet, rt.OpJsrW:
			return fail(pc, "legacy subroutine opcode 0x%02x is unsupported", op)
		case rt.OpMonitorenter, rt.OpMonitorexit:
			if !rt.ProfileHasThreads(profile) {
				return fail(pc, "monitor opcode 0x%02x is unsupported", op)
			}
			// Both monitor instructions are one byte with no operands. The
			// verifier consumes exactly one reference operand. An explicitly
			// retained length keeps the decoder independent of a typed nil
			// value so a hand-authored body still round-trips exactly.
		case rt.OpBreakpoint, rt.OpImpdep1, rt.OpImpdep2:
			return fail(pc, "reserved opcode 0x%02x is unsupported", op)
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

// canonicalLocalShort maps a one-byte local-variable opcode to its base opcode
// and local index.
func canonicalLocalShort(op uint8) (uint8, int32) {
	switch {
	case op >= 0x1a && op <= 0x1d:
		return rt.OpIload, int32(op - 0x1a)
	case op >= 0x1e && op <= 0x21:
		return rt.OpLload, int32(op - 0x1e)
	case op >= 0x22 && op <= 0x25:
		return rt.OpFload, int32(op - 0x22)
	case op >= 0x26 && op <= 0x29:
		return rt.OpDload, int32(op - 0x26)
	case op >= 0x2a && op <= 0x2d:
		return rt.OpAload, int32(op - 0x2a)
	case op >= 0x3b && op <= 0x3e:
		return rt.OpIstore, int32(op - 0x3b)
	case op >= 0x3f && op <= 0x42:
		return rt.OpLstore, int32(op - 0x3f)
	case op >= 0x43 && op <= 0x46:
		return rt.OpFstore, int32(op - 0x43)
	case op >= 0x47 && op <= 0x4a:
		return rt.OpDstore, int32(op - 0x47)
	case op >= 0x4b && op <= 0x4e:
		return rt.OpAstore, int32(op - 0x4b)
	}
	return op, 0
}

func decodeSwitch(cp *constantPool, r *reader, pc int, op uint8, insn *rt.Insn) error {
	_ = cp
	// The switch payload is 4-byte aligned relative to the start of the code
	// array. r.off is pc+1 here.
	pad := (4 - (r.off % 4)) % 4
	padding, err := r.bytes(pad)
	if err != nil {
		return fmt.Errorf("truncated switch padding")
	}
	for _, b := range padding {
		if b != 0 {
			return fmt.Errorf("nonzero switch padding byte 0x%02x", b)
		}
	}
	def, err := r.u4()
	if err != nil {
		return fmt.Errorf("truncated switch default")
	}
	defaultTarget := int32(int64(pc) + int64(int32(def)))
	insn.A = defaultTarget
	if op == rt.OpTableswitch {
		lowRaw, err := r.u4()
		if err != nil {
			return fmt.Errorf("truncated tableswitch low")
		}
		highRaw, err := r.u4()
		if err != nil {
			return fmt.Errorf("truncated tableswitch high")
		}
		low := int64(int32(lowRaw))
		high := int64(int32(highRaw))
		if high < low {
			return fmt.Errorf("tableswitch high %d is below low %d", int32(highRaw), int32(lowRaw))
		}
		count := high - low + 1
		if count <= 0 || count > 65536 {
			return fmt.Errorf("tableswitch range %d..%d is out of bounds", low, high)
		}
		for i := int64(0); i < count; i++ {
			off, err := r.u4()
			if err != nil {
				return fmt.Errorf("truncated tableswitch target %d", i)
			}
			insn.Keys = append(insn.Keys, int32(low+i))
			insn.Targets = append(insn.Targets, int32(int64(pc)+int64(int32(off))))
		}
		return nil
	}
	npairsRaw, err := r.u4()
	if err != nil {
		return fmt.Errorf("truncated lookupswitch pair count")
	}
	npairs := int64(int32(npairsRaw))
	if npairs < 0 || npairs > 65536 {
		return fmt.Errorf("lookupswitch pair count %d is out of bounds", int32(npairsRaw))
	}
	var previous int32
	for i := int64(0); i < npairs; i++ {
		keyRaw, err := r.u4()
		if err != nil {
			return fmt.Errorf("truncated lookupswitch key %d", i)
		}
		off, err := r.u4()
		if err != nil {
			return fmt.Errorf("truncated lookupswitch target %d", i)
		}
		key := int32(keyRaw)
		if i > 0 && key <= previous {
			return fmt.Errorf("lookupswitch keys are not strictly increasing (%d after %d)", key, previous)
		}
		previous = key
		insn.Keys = append(insn.Keys, key)
		insn.Targets = append(insn.Targets, int32(int64(pc)+int64(int32(off))))
	}
	return nil
}

func decodeWide(r *reader, pc int, insn *rt.Insn) error {
	sub, err := r.u1()
	if err != nil {
		return fmt.Errorf("truncated wide opcode")
	}
	switch sub {
	case rt.OpIload, rt.OpLload, rt.OpFload, rt.OpDload, rt.OpAload,
		rt.OpIstore, rt.OpLstore, rt.OpFstore, rt.OpDstore, rt.OpAstore:
		idx, err := r.u2()
		if err != nil {
			return fmt.Errorf("truncated wide local index")
		}
		insn.Op = sub
		insn.A = int32(idx)
		return nil
	case rt.OpIinc:
		idx, err := r.u2()
		if err != nil {
			return fmt.Errorf("truncated wide iinc index")
		}
		amount, err := r.u2()
		if err != nil {
			return fmt.Errorf("truncated wide iinc amount")
		}
		insn.Op = rt.OpIinc
		insn.A = int32(idx)
		insn.B = s16(amount)
		return nil
	case rt.OpRet:
		return fmt.Errorf("wide ret is unsupported at pc %d", pc)
	default:
		return fmt.Errorf("invalid wide opcode 0x%02x", sub)
	}
}

func ldcRef(cp *constantPool, idx uint16, wide bool) (rt.Ref, error) {
	e, err := cp.entry(idx)
	if err != nil {
		return rt.Ref{}, err
	}
	if wide {
		switch e.tag {
		case 5:
			return rt.Ref{Kind: rt.RefLong, Bits: e.bits, CP: int(idx)}, nil
		case 6:
			return rt.Ref{Kind: rt.RefDouble, Bits: e.bits, CP: int(idx)}, nil
		default:
			return rt.Ref{}, fmt.Errorf("ldc2_w constant tag %d is not long or double", e.tag)
		}
	}
	switch e.tag {
	case 3:
		return rt.Ref{Kind: rt.RefInt, Int: e.intVal, CP: int(idx)}, nil
	case 4:
		return rt.Ref{Kind: rt.RefFloat, Bits: e.bits, CP: int(idx)}, nil
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
	case 5, 6:
		return rt.Ref{}, fmt.Errorf("ldc cannot load a category-two constant")
	case 17:
		return rt.Ref{}, fmt.Errorf("dynamically computed constant at index %d is unsupported", idx)
	case 15, 16:
		return rt.Ref{}, fmt.Errorf("method handle/type constant at index %d is unsupported", idx)
	default:
		return rt.Ref{}, fmt.Errorf("unsupported ldc constant tag %d", e.tag)
	}
}

// BranchTarget returns the single branch target of one instruction.
func BranchTarget(insn rt.Insn) (int32, bool) {
	switch insn.Op {
	case rt.OpIfeq, rt.OpIfne, rt.OpIflt, rt.OpIfge, rt.OpIfgt, rt.OpIfle,
		rt.OpIfIcmpeq, rt.OpIfIcmpne, rt.OpIfIcmplt, rt.OpIfIcmpge,
		rt.OpIfIcmpgt, rt.OpIfIcmple, rt.OpIfAcmpeq, rt.OpIfAcmpne,
		rt.OpGoto, rt.OpIfnull, rt.OpIfnonnull, rt.OpGotoW:
		return insn.A, true
	}
	return 0, false
}

// SwitchTargets returns the default and case targets of a switch.
func SwitchTargets(insn rt.Insn) (defaultTarget int32, targets []int32, ok bool) {
	if insn.Op != rt.OpTableswitch && insn.Op != rt.OpLookupswitch {
		return 0, nil, false
	}
	return insn.A, insn.Targets, true
}

// CodeLengthFromPCs returns the total byte length implied by decoded offsets.
func CodeLengthFromPCs(insns []rt.Insn) int {
	if len(insns) == 0 {
		return 0
	}
	last := insns[len(insns)-1]
	return int(last.PC + last.Len)
}
