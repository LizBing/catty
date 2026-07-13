package interpreter

import (
	"fmt"
	"math"

	"catty/lowering"
	"catty/opcode"
	"catty/rtda"
)

// LoopIR is the IR executor: it runs methods from their lowered (stack-eliminated)
// form instead of raw bytecode. It is the validation harness for the lowering
// pass — if it produces output byte-identical to the tree-walking interpreter
// (Loop), the lowering is semantics-preserving.
//
// Per-instruction cost: the IR is predecoded (operands read once at lowering
// time), and the common arithmetic/branch ops address their inputs/outputs by
// the lowering's precomputed slot indices, so the hot loop does no operand
// parsing. The lowering for a method is computed lazily and cached.
func LoopIR(thread *rtda.Thread) {
	cache := map[*rtda.Method]*lowering.IR{}
	// The IR for the current frame is reused across every instruction in a loop
	// body, so we remember it and only re-resolve on a frame change (call/return)
	// — avoiding a map lookup per instruction, which otherwise dominates.
	var ir *lowering.IR
	var lastFrame *rtda.Frame
	for !thread.IsStackEmpty() {
		frame := thread.CurrentFrame()
		if frame != lastFrame {
			m := frame.Method()
			ir = cache[m]
			if ir == nil {
				var err error
				ir, err = lowering.Lower(m)
				if err != nil {
					panic(err)
				}
				cache[m] = ir
			}
			lastFrame = frame
		}
		pc := frame.PC()
		execIR(thread, frame, ir)
		if thread.HasException() {
			handleException(thread, pc)
			lastFrame = nil // force IR re-lookup after frame changes
		}
	}
}

// execIR runs one IR instruction. Every instruction first seeds the operand
// stack pointer from the lowering's known entry depth (so shared Push/Pop-based
// helpers and the shuffle ops behave exactly as in the tree-walker). Pure ops
// then read/write the precomputed Uses/Defs slots directly — exercising the
// stack elimination. pc is advanced by the instruction's decoded length unless a
// branch/return/switch overrides it.
func execIR(thread *rtda.Thread, frame *rtda.Frame, ir *lowering.IR) {
	pc := frame.PC()
	inst := &ir.Insts[pc]
	frame.SetPC(pc + inst.Length) // default: fall through to the next instruction
	frame.SetStackTop(inst.Depth)
	cp := frame.Method().Owner().ConstantPool()

	switch inst.Op {

	// ---------- constants (slot-index) ----------
	case opcode.Nop:
	case opcode.AconstNull:
		frame.SetStackSlotRef(int(inst.Defs[0]), nil)
	case opcode.IconstM1:
		frame.SetStackSlotNum(int(inst.Defs[0]), -1)
	case opcode.Iconst0, opcode.Iconst1, opcode.Iconst2,
		opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(inst.Op-opcode.Iconst0))
	case opcode.Lconst0, opcode.Lconst1:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(inst.Op-opcode.Lconst0))
	case opcode.Fconst0, opcode.Fconst1, opcode.Fconst2:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(float32(inst.Op-opcode.Fconst0))))
	case opcode.Dconst0, opcode.Dconst1:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(float64(inst.Op-opcode.Dconst0))))
	case opcode.Bipush:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(inst.Const8))
	case opcode.Sipush:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(inst.Const16))
	case opcode.Ldc, opcode.LdcW:
		pushConstant(thread, frame, cp, inst.Index)
	case opcode.Ldc2W:
		switch cp.Tag(inst.Index) {
		case 5: // CONSTANT_Long
			frame.PushLong(cp.Long(inst.Index))
		case 6: // CONSTANT_Double
			frame.PushDouble(cp.Double(inst.Index))
		}

	// ---------- loads (slot-index) ----------
	case opcode.Iload:
		frame.SetStackSlotNum(int(inst.Defs[0]), frame.GetInt(int(inst.Index)))
	case opcode.Aload:
		frame.SetStackSlotRef(int(inst.Defs[0]), frame.GetRef(int(inst.Index)))
	case opcode.Fload:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(frame.GetFloat(int(inst.Index)))))
	case opcode.Lload, opcode.Dload:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), frame.GetLong(int(inst.Index)))
	case opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3:
		frame.SetStackSlotNum(int(inst.Defs[0]), frame.GetInt(int(inst.Op-opcode.Iload0)))
	case opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		frame.SetStackSlotRef(int(inst.Defs[0]), frame.GetRef(int(inst.Op-opcode.Aload0)))
	case opcode.Lload0, opcode.Lload1, opcode.Lload2, opcode.Lload3:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), frame.GetLong(int(inst.Op-opcode.Lload0)))
	case opcode.Fload0, opcode.Fload1, opcode.Fload2, opcode.Fload3:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(frame.GetFloat(int(inst.Op-opcode.Fload0)))))
	case opcode.Dload0, opcode.Dload1, opcode.Dload2, opcode.Dload3:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(frame.GetDouble(int(inst.Op-opcode.Dload0)))))

	// ---------- stores (slot-index) ----------
	case opcode.Istore:
		frame.SetInt(int(inst.Index), frame.StackSlotNum(int(inst.Uses[0])))
	case opcode.Astore:
		frame.SetRef(int(inst.Index), frame.StackSlotRef(int(inst.Uses[0])))
	case opcode.Fstore:
		frame.SetFloat(int(inst.Index), math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0])))))
	case opcode.Lstore, opcode.Dstore:
		frame.SetLong(int(inst.Index), slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1])))
	case opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		frame.SetInt(int(inst.Op-opcode.Istore0), frame.StackSlotNum(int(inst.Uses[0])))
	case opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		frame.SetRef(int(inst.Op-opcode.Astore0), frame.StackSlotRef(int(inst.Uses[0])))
	case opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3:
		frame.SetLong(int(inst.Op-opcode.Lstore0), slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1])))
	case opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3:
		frame.SetFloat(int(inst.Op-opcode.Fstore0), math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0])))))
	case opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		frame.SetDouble(int(inst.Op-opcode.Dstore0), math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1])))))

	// ---------- array load (slot-index) ----------
	case opcode.Iaload, opcode.Baload, opcode.Caload, opcode.Saload:
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		frame.SetStackSlotNum(int(inst.Defs[0]), arr.Cells()[int(i)].GetInt())
	case opcode.Laload, opcode.Daload:
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		v := readTwoSlots(arr, int(i))
		if inst.Op == opcode.Daload {
			frame.SetStackSlotNum(int(inst.Defs[1]), int32(math.Float64bits(math.Float64frombits(uint64(v)))))
			frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float64bits(math.Float64frombits(uint64(v)))>>32))
		} else {
			setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), v)
		}
	case opcode.Faload:
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		frame.SetStackSlotNum(int(inst.Defs[0]), arr.Cells()[int(i)].GetInt())
	case opcode.Aaload:
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		frame.SetStackSlotRef(int(inst.Defs[0]), arr.Cells()[int(i)].GetRef())

	// ---------- array store (slot-index) ----------
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore:
		v := frame.StackSlotNum(int(inst.Uses[2]))
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		arr.Cells()[int(i)].SetInt(v)
	case opcode.Lastore, opcode.Dastore:
		v := slotLong(frame, int(inst.Uses[2]), int(inst.Uses[3]))
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		if inst.Op == opcode.Dastore {
			writeTwoSlots(arr, int(i), int64(math.Float64bits(math.Float64frombits(uint64(v)))))
		} else {
			writeTwoSlots(arr, int(i), v)
		}
	case opcode.Fastore:
		v := frame.StackSlotNum(int(inst.Uses[2]))
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		arr.Cells()[int(i)].SetInt(v)
	case opcode.Aastore:
		v := frame.StackSlotRef(int(inst.Uses[2]))
		i := frame.StackSlotNum(int(inst.Uses[1]))
		arr := frame.StackSlotRef(int(inst.Uses[0]))
		arr.Cells()[int(i)].SetRef(v)

	// ---------- stack shuffles (Push/Pop, seeded) ----------
	case opcode.Pop:
		frame.PopSlot()
	case opcode.Pop2:
		frame.PopSlot()
		frame.PopSlot()
	case opcode.Dup:
		s := frame.PopSlot()
		frame.PushSlot(s)
		frame.PushSlot(s)
	case opcode.DupX1:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opcode.DupX2:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		s3 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s3)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opcode.Dup2:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s2)
		frame.PushSlot(s1)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opcode.Swap:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s2)

	// ---------- int arithmetic (slot-index) ----------
	case opcode.Iadd, opcode.Isub, opcode.Imul, opcode.Idiv, opcode.Irem,
		opcode.Iand, opcode.Ior, opcode.Ixor:
		a := frame.StackSlotNum(int(inst.Uses[0]))
		b := frame.StackSlotNum(int(inst.Uses[1]))
		if (inst.Op == opcode.Idiv || inst.Op == opcode.Irem) && b == 0 {
			throwRuntime(thread, pc, "java/lang/ArithmeticException", "/ by zero")
			return
		}
		frame.SetStackSlotNum(int(inst.Defs[0]), iapply(inst.Op, a, b))
	case opcode.Ishl, opcode.Ishr, opcode.Iushr:
		s := uint32(frame.StackSlotNum(int(inst.Uses[1])))
		v := frame.StackSlotNum(int(inst.Uses[0]))
		frame.SetStackSlotNum(int(inst.Defs[0]), ishift(inst.Op, v, s))
	case opcode.Ineg:
		frame.SetStackSlotNum(int(inst.Defs[0]), -frame.StackSlotNum(int(inst.Uses[0])))
	case opcode.Iinc:
		idx := int(inst.IncIndex)
		frame.SetInt(idx, frame.GetInt(idx)+int32(inst.Const8))

	// ---------- long arithmetic (slot-index) ----------
	case opcode.Ladd, opcode.Lsub, opcode.Lmul, opcode.Ldiv, opcode.Lrem,
		opcode.Land, opcode.Lor, opcode.Lxor:
		a := slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))
		b := slotLong(frame, int(inst.Uses[2]), int(inst.Uses[3]))
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), lapply(inst.Op, a, b))
	case opcode.Lshl, opcode.Lshr, opcode.Lushr:
		s := uint32(frame.StackSlotNum(int(inst.Uses[2]))) & 0x3f
		v := slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), lshift(inst.Op, v, s))
	case opcode.Lneg:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), -slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1])))

	// ---------- float / double arithmetic (slot-index) ----------
	case opcode.Fadd, opcode.Fsub, opcode.Fmul, opcode.Fdiv, opcode.Frem:
		a := math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))
		b := math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[1]))))
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(fapply(inst.Op, a, b))))
	case opcode.Fneg:
		a := math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(-a)))
	case opcode.Dadd, opcode.Dsub, opcode.Dmul, opcode.Ddiv, opcode.Drem:
		a := math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))
		b := math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[2]), int(inst.Uses[3]))))
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(dapply(inst.Op, a, b))))
	case opcode.Dneg:
		a := math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(-a)))

	// ---------- conversions (slot-index) ----------
	case opcode.I2l:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(frame.StackSlotNum(int(inst.Uses[0]))))
	case opcode.I2f:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(float32(frame.StackSlotNum(int(inst.Uses[0]))))))
	case opcode.I2d:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(float64(frame.StackSlotNum(int(inst.Uses[0]))))))
	case opcode.L2i:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))
	case opcode.L2f:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(float32(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))))
	case opcode.L2d:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(float64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))))
	case opcode.F2i:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))))
	case opcode.F2l:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))))
	case opcode.F2d:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64bits(float64(math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))))))
	case opcode.D2i:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))))
	case opcode.D2l:
		setSlotLong(frame, int(inst.Defs[0]), int(inst.Defs[1]), int64(math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))))
	case opcode.D2f:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(math.Float32bits(float32(math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))))))
	case opcode.I2b:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(int8(frame.StackSlotNum(int(inst.Uses[0])))))
	case opcode.I2c:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(uint16(frame.StackSlotNum(int(inst.Uses[0])))))
	case opcode.I2s:
		frame.SetStackSlotNum(int(inst.Defs[0]), int32(int16(frame.StackSlotNum(int(inst.Uses[0])))))

	// ---------- comparisons (slot-index) ----------
	case opcode.Lcmp:
		b := slotLong(frame, int(inst.Uses[2]), int(inst.Uses[3]))
		a := slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))
		switch {
		case a > b:
			frame.SetStackSlotNum(int(inst.Defs[0]), 1)
		case a < b:
			frame.SetStackSlotNum(int(inst.Defs[0]), -1)
		default:
			frame.SetStackSlotNum(int(inst.Defs[0]), 0)
		}
	case opcode.Fcmpl, opcode.Fcmpg:
		b := math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[1]))))
		a := math.Float32frombits(uint32(frame.StackSlotNum(int(inst.Uses[0]))))
		frame.SetStackSlotNum(int(inst.Defs[0]), cmpFloat(a, b, inst.Op==opcode.Fcmpg))
	case opcode.Dcmpl, opcode.Dcmpg:
		b := math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[2]), int(inst.Uses[3]))))
		a := math.Float64frombits(uint64(slotLong(frame, int(inst.Uses[0]), int(inst.Uses[1]))))
		frame.SetStackSlotNum(int(inst.Defs[0]), cmpDouble(a, b, inst.Op==opcode.Dcmpg))

	// ---------- branches (slot-index for operands, predecoded target) ----------
	case opcode.Ifeq:
		if frame.StackSlotNum(int(inst.Uses[0])) == 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifne:
		if frame.StackSlotNum(int(inst.Uses[0])) != 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.Iflt:
		if frame.StackSlotNum(int(inst.Uses[0])) < 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifge:
		if frame.StackSlotNum(int(inst.Uses[0])) >= 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifgt:
		if frame.StackSlotNum(int(inst.Uses[0])) > 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifle:
		if frame.StackSlotNum(int(inst.Uses[0])) <= 0 {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmpeq:
		if frame.StackSlotNum(int(inst.Uses[0])) == frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmpne:
		if frame.StackSlotNum(int(inst.Uses[0])) != frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmplt:
		if frame.StackSlotNum(int(inst.Uses[0])) < frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmpge:
		if frame.StackSlotNum(int(inst.Uses[0])) >= frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmpgt:
		if frame.StackSlotNum(int(inst.Uses[0])) > frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfIcmple:
		if frame.StackSlotNum(int(inst.Uses[0])) <= frame.StackSlotNum(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfAcmpeq:
		if frame.StackSlotRef(int(inst.Uses[0])) == frame.StackSlotRef(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.IfAcmpne:
		if frame.StackSlotRef(int(inst.Uses[0])) != frame.StackSlotRef(int(inst.Uses[1])) {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifnull:
		if frame.StackSlotRef(int(inst.Uses[0])) == nil {
			frame.SetPC(inst.Branch)
		}
	case opcode.Ifnonnull:
		if frame.StackSlotRef(int(inst.Uses[0])) != nil {
			frame.SetPC(inst.Branch)
		}

	case opcode.Goto, opcode.GotoW:
		frame.SetPC(inst.Branch)
	case opcode.Tableswitch:
		key := frame.StackSlotNum(int(inst.Uses[0]))
		st := inst.Switch
		if key >= st.Low && key <= st.High {
			frame.SetPC(st.Targets[key-st.Low])
		} else {
			frame.SetPC(st.Default)
		}
	case opcode.Lookupswitch:
		key := frame.StackSlotNum(int(inst.Uses[0]))
		st := inst.Switch
		t := st.Default
		for i, k := range st.Keys {
			if k == key {
				t = st.Targets[i]
				break
			}
		}
		frame.SetPC(t)

	// ---------- returns (shared helpers) ----------
	case opcode.Return:
		thread.PopFrame()
	case opcode.Ireturn:
		returnInt(frame, thread)
	case opcode.Areturn:
		returnRef(frame, thread)
	case opcode.Lreturn:
		returnLong(frame, thread)
	case opcode.Freturn:
		returnFloat(frame, thread)
	case opcode.Dreturn:
		returnDouble(frame, thread)

	// ---------- fields (shared helpers, seeded stackTop) ----------
	case opcode.Getstatic:
		cls, name, desc := cp.MemberRef(inst.Index)
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		ensureInitialized(thread, field.Owner())
		loadFieldValue(frame, field.Owner().StaticCells(), field.SlotID(), desc)
	case opcode.Putstatic:
		cls, name, desc := cp.MemberRef(inst.Index)
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		ensureInitialized(thread, field.Owner())
		storeFieldValue(frame, field.Owner().StaticCells(), field.SlotID(), desc)
	case opcode.Getfield:
		cls, name, desc := cp.MemberRef(inst.Index)
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		obj := frame.PopRef()
		loadFieldValue(frame, obj.Cells(), field.SlotID(), desc)
	case opcode.Putfield:
		// Stack: [..., objref, value], value on top — pop value first, then objref.
		cls, name, desc := cp.MemberRef(inst.Index)
		field := thread.Loader().LoadClass(cls).LookupField(name, desc)
		slotID := field.SlotID()
		switch desc[0] {
		case 'J':
			v := frame.PopLong()
			obj := frame.PopRef()
			obj.Cells()[slotID].SetLong(v)
		case 'D':
			v := frame.PopDouble()
			obj := frame.PopRef()
			obj.Cells()[slotID].SetDouble(v)
		case 'F':
			v := frame.PopFloat()
			obj := frame.PopRef()
			obj.Cells()[slotID].SetFloat(v)
		case 'L', '[':
			v := frame.PopRef()
			obj := frame.PopRef()
			obj.Cells()[slotID].SetRef(v)
		default:
			v := frame.PopInt()
			obj := frame.PopRef()
			obj.Cells()[slotID].SetInt(v)
		}

	// ---------- invocations (shared helpers) ----------
	case opcode.Invokevirtual:
		cls, name, desc := cp.MemberRef(inst.Index)
		class := thread.Loader().LoadClass(cls)
		spec := class.LookupMethod(name, desc)
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			throwRuntime(thread, pc, "java/lang/NullPointerException", "")
			return
		}
		invokeMethod(thread, receiver.Class().LookupMethod(name, desc))
	case opcode.Invokespecial:
		cls, name, desc := cp.MemberRef(inst.Index)
		class := thread.Loader().LoadClass(cls)
		invokeMethod(thread, class.LookupMethod(name, desc))
	case opcode.Invokestatic:
		cls, name, desc := cp.MemberRef(inst.Index)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		invokeMethod(thread, class.LookupMethod(name, desc))
	case opcode.Invokeinterface:
		cls, name, desc := cp.MemberRef(inst.Index)
		spec := thread.Loader().LoadClass(cls).LookupMethod(name, desc)
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			throwRuntime(thread, pc, "java/lang/NullPointerException", "")
			return
		}
		invokeMethod(thread, receiver.Class().LookupMethod(name, desc))

	// ---------- object / array / misc (shared helpers) ----------
	case opcode.New:
		class := thread.Loader().LoadClass(cp.ClassName(inst.Index))
		ensureInitialized(thread, class)
		frame.PushRef(rtda.NewObject(class))
	case opcode.Newarray:
		frame.PushRef(newPrimitiveArray(thread, inst.Atype, int(frame.PopInt())))
	case opcode.Anewarray:
		elemName := cp.ClassName(inst.Index)
		frame.PushRef(newRefArray(thread, elemName, int(frame.PopInt())))
	case opcode.Multianewarray:
		className := cp.ClassName(inst.Index)
		class := thread.Loader().LoadClass(className)
		dims := int(inst.Count)
		sizes := make([]int, dims)
		for i := dims - 1; i >= 0; i-- {
			sizes[i] = int(frame.PopInt())
		}
		frame.PushRef(rtda.NewMultiArray(class, sizes, thread.Loader()))
	case opcode.Arraylength:
		arr := frame.PopRef()
		frame.PushInt(int32(arr.ArrayLength()))
	case opcode.Checkcast:
		target := thread.Loader().LoadClass(cp.ClassName(inst.Index))
		obj := frame.PeekRef(0)
		if obj != nil && !obj.IsInstanceOf(target) {
			throwRuntime(thread, pc, "java/lang/ClassCastException", "")
			return
		}
	case opcode.Instanceof:
		target := thread.Loader().LoadClass(cp.ClassName(inst.Index))
		obj := frame.PopRef()
		if obj != nil && obj.IsInstanceOf(target) {
			frame.PushInt(1)
		} else {
			frame.PushInt(0)
		}
	case opcode.Monitorenter, opcode.Monitorexit:
		frame.PopRef()

	case opcode.Athrow:
		excObj := frame.PopRef()
		if excObj == nil {
			throwRuntime(thread, pc, "java/lang/NullPointerException", "")
			return
		}
		thread.Throw(excObj, pc)

	default:
		panic(fmt.Sprintf("catty: IR opcode 0x%02x (%s) not implemented", byte(inst.Op), opcode.Name(inst.Op)))
	}
}

// slotLong reads a category-2 value from two stack slots (high at hi, low at lo).
func slotLong(f *rtda.Frame, hi, lo int) int64 {
	return int64(f.StackSlotNum(hi))<<32 | int64(uint32(f.StackSlotNum(lo)))
}

func setSlotLong(f *rtda.Frame, hi, lo int, v int64) {
	f.SetStackSlotNum(hi, int32(uint64(v)>>32))
	f.SetStackSlotNum(lo, int32(v))
}

// iapply/lapply/fapply/dapply/ishift/lshift fold the binary-op switch for the
// slot-index arithmetic cases, keeping execIR compact.
func iapply(op opcode.Opcode, a, b int32) int32 {
	switch op {
	case opcode.Iadd:
		return a + b
	case opcode.Isub:
		return a - b
	case opcode.Imul:
		return a * b
	case opcode.Idiv:
		return a / b
	case opcode.Irem:
		return a % b
	case opcode.Iand:
		return a & b
	case opcode.Ior:
		return a | b
	case opcode.Ixor:
		return a ^ b
	}
	panic("iapply " + opcode.Name(op))
}

func ishift(op opcode.Opcode, v int32, s uint32) int32 {
	switch op {
	case opcode.Ishl:
		return v << (s & 0x1f)
	case opcode.Ishr:
		return v >> (s & 0x1f)
	case opcode.Iushr:
		return int32(uint32(v) >> (s & 0x1f))
	}
	panic("ishift " + opcode.Name(op))
}

func lapply(op opcode.Opcode, a, b int64) int64 {
	switch op {
	case opcode.Ladd:
		return a + b
	case opcode.Lsub:
		return a - b
	case opcode.Lmul:
		return a * b
	case opcode.Ldiv:
		return a / b
	case opcode.Lrem:
		return a % b
	case opcode.Land:
		return a & b
	case opcode.Lor:
		return a | b
	case opcode.Lxor:
		return a ^ b
	}
	panic("lapply " + opcode.Name(op))
}

func lshift(op opcode.Opcode, v int64, s uint32) int64 {
	switch op {
	case opcode.Lshl:
		return v << s
	case opcode.Lshr:
		return v >> s
	case opcode.Lushr:
		return int64(uint64(v) >> s)
	}
	panic("lshift " + opcode.Name(op))
}

func fapply(op opcode.Opcode, a, b float32) float32 {
	switch op {
	case opcode.Fadd:
		return a + b
	case opcode.Fsub:
		return a - b
	case opcode.Fmul:
		return a * b
	case opcode.Fdiv:
		return a / b
	case opcode.Frem:
		return remF(a, b)
	}
	panic("fapply " + opcode.Name(op))
}

func dapply(op opcode.Opcode, a, b float64) float64 {
	switch op {
	case opcode.Dadd:
		return a + b
	case opcode.Dsub:
		return a - b
	case opcode.Dmul:
		return a * b
	case opcode.Ddiv:
		return a / b
	case opcode.Drem:
		return remF64(a, b)
	}
	panic("dapply " + opcode.Name(op))
}
