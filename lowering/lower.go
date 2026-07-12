package lowering

import (
	"fmt"

	"catty/classfile"
	"catty/opcode"
	"catty/rtda"
)

// Lower converts a method's bytecode into the register-form IR in three passes:
//  1. decode — walk instructions linearly, predecoding operands and resolving
//     branch/switch targets to absolute pcs.
//  2. depth dataflow — a forward worklist over control-flow edges computes the
//     operand-stack depth (in slots) on entry to each instruction.
//  3. vreg assignment — turn each instruction's (entry depth, slot effect) into
//     concrete Uses/Defs slot indices: the operand stack, eliminated.
//
// Lowering is a pure function of the method's bytecode and its owner's constant
// pool — it needs no Loader, because field/invoke slot effects come from the
// descriptor in the constant-pool ref.
func Lower(method *rtda.Method) (*IR, error) {
	code := method.Code()
	ir := &IR{Insts: make([]IRInst, len(code))}
	if len(code) == 0 {
		return ir, nil
	}
	cp := method.Owner().ConstantPool()

	starts := decode(code, ir)
	depth, err := depthDataflow(code, cp, ir, starts, method.ExceptionTable())
	if err != nil {
		return nil, err
	}
	for _, pc := range starts {
		d := depth[pc]
		if d < 0 {
			continue // unreachable code — never executed, so its vregs are irrelevant
		}
		inst := &ir.Insts[pc]
		inst.Depth = d
		assignVregs(inst, cp)
	}
	typeDataflow(method, ir, starts, cp)

	return ir, nil
}

// decode walks instructions from pc 0, filling ir.Insts[pc] and returning the
// ordered list of instruction-start pcs. Bytecode is a linear, non-overlapping
// instruction sequence, so pc += length disassembles it correctly (switch
// padding/tables are skipped by their full computed length).
func decode(code []byte, ir *IR) []int {
	var starts []int
	for pc := 0; pc < len(code); {
		inst, length := decodeInst(code, pc)
		ir.Insts[pc] = inst
		starts = append(starts, pc)
		pc += length
	}
	return starts
}

func decodeInst(code []byte, pc int) (IRInst, int) {
	op := opcode.Opcode(code[pc])
	inst := IRInst{Op: op, Present: true}
	switch op {
	// --- u1 operands ---
	case opcode.Ldc:
		inst.Index = uint16(code[pc+1])
		inst.Length = 2
	case opcode.Bipush:
		inst.Const8 = int8(code[pc+1])
		inst.Length = 2
	case opcode.Newarray:
		inst.Atype = code[pc+1]
		inst.Length = 2
	case opcode.Iload, opcode.Lload, opcode.Fload, opcode.Dload, opcode.Aload,
		opcode.Istore, opcode.Lstore, opcode.Fstore, opcode.Dstore, opcode.Astore:
		inst.Index = uint16(code[pc+1])
		inst.Length = 2

	// --- u2 operands ---
	case opcode.Sipush:
		inst.Const16 = int16(be16(code, pc+1))
		inst.Length = 3
	case opcode.LdcW, opcode.Ldc2W,
		opcode.Getstatic, opcode.Putstatic, opcode.Getfield, opcode.Putfield,
		opcode.New, opcode.Anewarray, opcode.Checkcast, opcode.Instanceof,
		opcode.Invokevirtual, opcode.Invokespecial, opcode.Invokestatic:
		inst.Index = be16(code, pc+1)
		inst.Length = 3

	// --- iinc: u1 index + s1 delta ---
	case opcode.Iinc:
		inst.IncIndex = code[pc+1]
		inst.Const8 = int8(code[pc+2])
		inst.Length = 3

	// --- conditional branches: s2 offset → absolute target ---
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne,
		opcode.Ifnull, opcode.Ifnonnull:
		inst.Branch = pc + int(int16(be16(code, pc+1)))
		inst.Length = 3
	case opcode.Goto:
		inst.Branch = pc + int(int16(be16(code, pc+1)))
		inst.Length = 3
	case opcode.GotoW:
		inst.Branch = pc + int(int32(be32(code, pc+1)))
		inst.Length = 5

	// --- switches ---
	case opcode.Tableswitch:
		inst.Switch = decodeTableSwitch(code, pc)
		inst.Length = switchLength(code, pc, true)
	case opcode.Lookupswitch:
		inst.Switch = decodeLookupSwitch(code, pc)
		inst.Length = switchLength(code, pc, false)

	// --- invokeinterface (length 5: u2 index + u1 count + u1 zero) ---
	case opcode.Invokeinterface:
		inst.Index = be16(code, pc+1)
		inst.Length = 5

	// --- multianewarray (length 4: u2 class-index + u1 dimensions) ---
	case opcode.Multianewarray:
		inst.Index = be16(code, pc+1)
		inst.Count = code[pc+3] // dimensions
		inst.Length = 4

	// --- wide (length 4 for load/store, 6 for iinc) ---
	case opcode.Wide:
		modifiedOp := opcode.Opcode(code[pc+1])
		inst.Index = be16(code, pc+2)
		if modifiedOp == opcode.Iinc {
			inst.IncIndex = uint8(inst.Index)         // reuse IncIndex for the u2 index
			inst.Const8 = int8(be16(code, pc+4) >> 8) // approximate; wide iinc needs special handling
			inst.Length = 6
		} else {
			inst.Length = 4
		}

	// --- still unsupported ---
	case opcode.Invokedynamic:
		panic("catty/lowering: invokedynamic not yet supported")

	default:
		inst.Length = 1
	}
	return inst, inst.Length
}

// depthDataflow computes, for each instruction-start pc, the operand-stack depth
// on entry. A pc reached at two different depths ⇒ malformed bytecode (javac
// output is verifier-valid, so this won't fire on the fixtures). Returns the
// depth slice (−1 at unreachable pcs).
func depthDataflow(code []byte, cp *classfile.ConstantPool, ir *IR, starts []int, exTable []rtda.ExceptionEntry) ([]int, error) {
	const unset = -1
	depth := make([]int, len(code))
	for i := range depth {
		depth[i] = unset
	}
	depth[starts[0]] = 0
	worklist := []int{starts[0]}
	for len(worklist) > 0 {
		n := len(worklist) - 1
		pc := worklist[n]
		worklist = worklist[:n]
		inst := &ir.Insts[pc]
		var pop, push int
		if inst.Op == opcode.Multianewarray {
			pop = int(inst.Count) // pop `dimensions` counts
			push = 1
		} else {
			pop, push = slotEffect(inst.Op, cp, inst.Index)
		}
		next := depth[pc] + push - pop
		fall, hasFall, targets := successors(inst, pc)
		for _, s := range targets {
			if err := propagate(depth, s, next, &worklist); err != nil {
				return nil, err
			}
		}
		if hasFall {
			if err := propagate(depth, fall, next, &worklist); err != nil {
				return nil, err
			}
		}
		// Exception edges: if this PC is inside any protected range, the
		// handler is a successor with depth 1 (the exception on the stack).
		for _, entry := range exTable {
			if pc >= entry.StartPc() && pc < entry.EndPc() {
				if err := propagate(depth, entry.HandlerPc(), 1, &worklist); err != nil {
					return nil, err
				}
			}
		}
	}
	return depth, nil
}

func propagate(depth []int, pc, next int, worklist *[]int) error {
	switch {
	case depth[pc] == -1:
		depth[pc] = next
		*worklist = append(*worklist, pc)
	case depth[pc] != next:
		return fmt.Errorf("catty/lowering: stack-depth mismatch at pc %d (%d ≠ %d)", pc, depth[pc], next)
	}
	return nil
}

// successors returns the explicit branch targets, whether control falls through
// to pc+Length, and (if so) the fall-through pc.
func successors(inst *IRInst, pc int) (fall int, hasFall bool, targets []int) {
	switch inst.Op {
	case opcode.Goto, opcode.GotoW:
		return 0, false, []int{inst.Branch}
	case opcode.Return, opcode.Ireturn, opcode.Lreturn, opcode.Freturn,
		opcode.Dreturn, opcode.Areturn, opcode.Athrow:
		return 0, false, nil
	case opcode.Tableswitch, opcode.Lookupswitch:
		ts := inst.Switch
		tgts := make([]int, 0, 1+len(ts.Targets))
		tgts = append(tgts, ts.Default)
		tgts = append(tgts, ts.Targets...)
		return 0, false, tgts
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne,
		opcode.Ifnull, opcode.Ifnonnull:
		return pc + inst.Length, true, []int{inst.Branch}
	default:
		return pc + inst.Length, true, nil
	}
}

// assignVregs turns an instruction's entry depth + slot effect into concrete
// Uses/Defs slot indices: the input slots are the top `pop`, the output slots
// overwrite starting at depth-pop. Stack-shuffle ops (dup_x1/swap/…) get net
// effects here for record-keeping; the executor special-cases their movement.
func assignVregs(inst *IRInst, cp *classfile.ConstantPool) {
	pop, push := slotEffect(inst.Op, cp, inst.Index)
	d := inst.Depth
	for i := d - pop; i < d; i++ {
		inst.Uses = append(inst.Uses, uint8(i))
	}
	for i := d - pop; i < d-pop+push; i++ {
		inst.Defs = append(inst.Defs, uint8(i))
	}
}

// --- switch decoding ---

func decodeTableSwitch(code []byte, pc int) *SwitchTable {
	base := pc + 1
	p := base + padTo4(base)
	defOff := int(int32(be32(code, p)))
	low := int32(be32(code, p+4))
	high := int32(be32(code, p+8))
	n := high - low + 1
	st := &SwitchTable{Default: pc + defOff, Low: low, High: high, Targets: make([]int, n)}
	q := p + 12
	for i := range st.Targets {
		st.Targets[i] = pc + int(int32(be32(code, q)))
		q += 4
	}
	return st
}

func decodeLookupSwitch(code []byte, pc int) *SwitchTable {
	base := pc + 1
	p := base + padTo4(base)
	defOff := int(int32(be32(code, p)))
	npairs := int32(be32(code, p+4))
	st := &SwitchTable{Default: pc + defOff, Keys: make([]int32, npairs), Targets: make([]int, npairs)}
	q := p + 8
	for i := range st.Targets {
		st.Keys[i] = int32(be32(code, q))
		st.Targets[i] = pc + int(int32(be32(code, q+4)))
		q += 8
	}
	return st
}

// switchLength returns the byte length of a tableswitch/lookupswitch, including
// padding and the offset table.
func switchLength(code []byte, pc int, table bool) int {
	base := pc + 1
	p := base + padTo4(base)
	if table {
		low := int32(be32(code, p+4))
		high := int32(be32(code, p+8))
		n := high - low + 1
		return (p - pc) + 12 + int(n)*4
	}
	npairs := int32(be32(code, p+4))
	return (p - pc) + 8 + int(npairs)*8
}

func padTo4(base int) int { return (4 - base%4) % 4 }

// --- big-endian readers ---

func be16(code []byte, off int) uint16 { return uint16(code[off])<<8 | uint16(code[off+1]) }

func be32(code []byte, off int) uint32 {
	return uint32(code[off])<<24 | uint32(code[off+1])<<16 | uint32(code[off+2])<<8 | uint32(code[off+3])
}
