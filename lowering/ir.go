// Package lowering converts a method's stack-based bytecode into a register-form
// IR: the operand stack is eliminated into virtual registers (each vreg is the
// slot index it occupies). This is the keystone analysis for the AOT transpiler
// (ROADMAP Theme A): it proves the stack can be statically resolved, and as a
// byproduct gives the interpreter a predecoded form to run.
//
// A0 is depth-only: each opcode's slot effect is statically known (resolved from
// the constant-pool descriptor for fields/invokes), so a depth dataflow suffices
// to assign vregs. No type inference, no SSA, no phis — position-stable slot
// indices are sound for execution because JVMS guarantees equal stack depth at
// every merge point.
package lowering

import (
	"catty/classfile"
	"catty/opcode"
	"catty/rtda"
)

// IR is the lowered form of one method: a slice indexed by bytecode pc. Entries
// at pcs that are not instruction starts are zero-valued (Present=false).
type IR struct {
	Insts []IRInst
}

// IRInst is one decoded, stack-eliminated instruction.
type IRInst struct {
	Op     opcode.Opcode
	Present bool
	Length int    // bytes consumed in the original code (for fallthrough / decode walk)
	Depth  int    // operand-stack depth (in slots) on entry, from the depth dataflow
	Uses   []uint8 // slot indices this instruction reads
	Defs   []uint8 // slot indices this instruction writes

	// Predecoded immediate operands (only the relevant field is set per op).
	Index    uint16       // local index / constant-pool index
	Const8   int8         // bipush value / iinc delta
	Const16  int16        // sipush value
	IncIndex uint8        // iinc local index
	Atype    uint8        // newarray element type
	Count    uint8        // invokeinterface count / multianewarray dims
	Branch   int          // resolved absolute target pc (goto / if*)
	Switch   *SwitchTable // resolved switch table (tableswitch / lookupswitch)
}

// SwitchTable is the decoded, pc-resolved form of tableswitch / lookupswitch.
type SwitchTable struct {
	Default  int      // absolute pc of the default handler
	Low      int32    // tableswitch only: lowest key
	High     int32    // tableswitch only: highest key
	Keys     []int32  // lookupswitch only: match keys (parallel to Targets)
	Targets  []int    // absolute target pcs; tableswitch indexes by (key-Low)
}

// slotEffect returns the (pop, push) operand-stack effect in *slots* for op.
// For fields and invokes the effect depends on the referenced descriptor, read
// from the constant pool — no class loading is needed, only the cp ref, which is
// why lowering is a pure function of the method bytecode + its constant pool.
func slotEffect(op opcode.Opcode, cp *classfile.ConstantPool, cpIndex uint16) (pop, push int) {
	switch op {
	// --- constants ---
	case opcode.AconstNull:
		return 0, 1
	case opcode.IconstM1, opcode.Iconst0, opcode.Iconst1, opcode.Iconst2,
		opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		return 0, 1
	case opcode.Lconst0, opcode.Lconst1:
		return 0, 2
	case opcode.Fconst0, opcode.Fconst1, opcode.Fconst2:
		return 0, 1
	case opcode.Dconst0, opcode.Dconst1:
		return 0, 2
	case opcode.Bipush, opcode.Sipush, opcode.Ldc, opcode.LdcW:
		return 0, 1
	case opcode.Ldc2W:
		return 0, 2

	// --- loads ---
	case opcode.Iload, opcode.Fload, opcode.Aload,
		opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3,
		opcode.Fload0, opcode.Fload1, opcode.Fload2, opcode.Fload3,
		opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		return 0, 1
	case opcode.Lload, opcode.Dload,
		opcode.Lload0, opcode.Lload1, opcode.Lload2, opcode.Lload3,
		opcode.Dload0, opcode.Dload1, opcode.Dload2, opcode.Dload3:
		return 0, 2

	// --- stores ---
	case opcode.Istore, opcode.Fstore, opcode.Astore,
		opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3,
		opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3,
		opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		return 1, 0
	case opcode.Lstore, opcode.Dstore,
		opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3,
		opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		return 2, 0

	// --- array load ---
	case opcode.Iaload, opcode.Baload, opcode.Caload, opcode.Saload, opcode.Faload, opcode.Aaload:
		return 2, 1
	case opcode.Laload, opcode.Daload:
		return 2, 2

	// --- array store ---
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore, opcode.Fastore, opcode.Aastore:
		return 3, 0
	case opcode.Lastore, opcode.Dastore:
		return 4, 0

	// --- stack manipulation (net effect for depth; executor special-cases the copy) ---
	case opcode.Pop:
		return 1, 0
	case opcode.Pop2:
		return 2, 0
	case opcode.Dup:
		return 1, 2
	case opcode.DupX1:
		return 2, 3
	case opcode.DupX2:
		return 3, 4
	case opcode.Dup2:
		return 2, 4
	case opcode.Dup2X1:
		return 3, 5
	case opcode.Dup2X2:
		return 4, 6
	case opcode.Swap:
		return 2, 2

	// --- arithmetic: int ---
	case opcode.Iadd, opcode.Isub, opcode.Imul, opcode.Idiv, opcode.Irem,
		opcode.Ishl, opcode.Ishr, opcode.Iushr, opcode.Iand, opcode.Ior, opcode.Ixor:
		return 2, 1
	case opcode.Ineg:
		return 1, 1
	case opcode.Iinc:
		return 0, 0

	// --- arithmetic: long ---
	case opcode.Ladd, opcode.Lsub, opcode.Lmul, opcode.Ldiv, opcode.Lrem,
		opcode.Land, opcode.Lor, opcode.Lxor:
		return 4, 2
	case opcode.Lshl, opcode.Lshr, opcode.Lushr:
		return 3, 2 // shift count is int (1) + long (2)
	case opcode.Lneg:
		return 2, 2

	// --- arithmetic: float / double ---
	case opcode.Fadd, opcode.Fsub, opcode.Fmul, opcode.Fdiv, opcode.Frem:
		return 2, 1
	case opcode.Fneg:
		return 1, 1
	case opcode.Dadd, opcode.Dsub, opcode.Dmul, opcode.Ddiv, opcode.Drem:
		return 4, 2
	case opcode.Dneg:
		return 2, 2

	// --- conversions ---
	case opcode.I2l, opcode.I2d, opcode.F2l, opcode.F2d:
		return 1, 2
	case opcode.I2f, opcode.F2i, opcode.I2b, opcode.I2c, opcode.I2s:
		return 1, 1
	case opcode.L2i, opcode.L2f, opcode.D2i, opcode.D2f:
		return 2, 1
	case opcode.L2d, opcode.D2l:
		return 2, 2

	// --- comparisons ---
	case opcode.Lcmp, opcode.Dcmpl, opcode.Dcmpg:
		return 4, 1
	case opcode.Fcmpl, opcode.Fcmpg:
		return 2, 1

	// --- conditional branches (pop operands; push nothing) ---
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.Ifnull, opcode.Ifnonnull:
		return 1, 0
	case opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne:
		return 2, 0

	// --- unconditional / switches / returns ---
	case opcode.Goto, opcode.GotoW:
		return 0, 0
	case opcode.Tableswitch, opcode.Lookupswitch:
		return 1, 0 // pop the key, jump
	case opcode.Ireturn, opcode.Freturn, opcode.Areturn:
		return 1, 0
	case opcode.Lreturn, opcode.Dreturn:
		return 2, 0
	case opcode.Return:
		return 0, 0

	// --- fields (effect depends on descriptor read from the cp ref) ---
	case opcode.Getstatic, opcode.Putstatic, opcode.Getfield, opcode.Putfield:
		return fieldEffect(op, cp, cpIndex)

	// --- invokes (effect depends on descriptor) ---
	case opcode.Invokevirtual, opcode.Invokespecial, opcode.Invokestatic, opcode.Invokeinterface:
		return invokeEffect(op, cp, cpIndex)

	// --- object / array / misc ---
	case opcode.New:
		return 0, 1
	case opcode.Newarray, opcode.Anewarray:
		return 1, 1 // consume the length operand, push the array reference
	case opcode.Arraylength:
		return 1, 1
	case opcode.Athrow:
		return 1, 0
	case opcode.Checkcast:
		return 1, 1 // consumes objectref, pushes the same reference
	case opcode.Instanceof:
		return 1, 1
	case opcode.Monitorenter, opcode.Monitorexit:
		return 1, 0
	case opcode.Multianewarray:
		return 0, 1 // pop dims; for depth purposes the caller passes the exact pop via Count
	}
	panic("catty/lowering: unknown opcode for slot effect: " + opcode.Name(op))
}

// fieldEffect computes the slot effect of get/putstatic/get/putfield from the
// referenced field's descriptor (category 2 = long/double = 2 slots).
func fieldEffect(op opcode.Opcode, cp *classfile.ConstantPool, cpIndex uint16) (pop, push int) {
	_, _, desc := cp.MemberRef(cpIndex)
	cat := 1
	if desc == "J" || desc == "D" {
		cat = 2
	}
	switch op {
	case opcode.Getstatic:
		return 0, cat
	case opcode.Putstatic:
		return cat, 0
	case opcode.Getfield:
		return 1, cat // objectref + result
	case opcode.Putfield:
		return 1 + cat, 0 // objectref + value
	}
	return 0, 0
}

// invokeEffect computes the slot effect of an invoke from the referenced
// method descriptor (args + return type), without loading any class.
func invokeEffect(op opcode.Opcode, cp *classfile.ConstantPool, cpIndex uint16) (pop, push int) {
	_, _, desc := cp.MemberRef(cpIndex)
	md := rtda.ParseMethodDescriptor(desc)
	popSlots := md.ArgSlots()
	if op != opcode.Invokestatic {
		popSlots++ // `this`
	}
	return popSlots, returnSlots(md.ReturnType)
}

// returnSlots maps a return-type descriptor to its slot width (V=0, J/D=2, else 1).
func returnSlots(ret string) int {
	if ret == "" || ret == "V" {
		return 0
	}
	if ret == "J" || ret == "D" {
		return 2
	}
	return 1
}
