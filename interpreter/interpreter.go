package interpreter

import (
	"fmt"

	"catty/opcode"
	"catty/rtda"
)

// Loop is the bytecode interpreter. It runs the current thread's top frame's
// bytecode until the JVM stack drains (the bottom frame — main — returned).
//
// Dispatch is a single dense switch on the opcode byte. Go has no computed
// goto, so this switch — not a function-pointer table — is the fastest portable
// dispatch (a call-threaded design costs more per instruction than it saves here,
// given Go's calling convention). The pc is managed explicitly so branch
// targets are computed against opcodePc (the opcode's own address).
func Loop(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.CurrentFrame()
		opcodePc := frame.PC()
		op := opcode.Opcode(frame.Code()[opcodePc])
		frame.SetPC(opcodePc + 1)
		exec(thread, frame, op, opcodePc)
	}
}

// exec runs one instruction. Split out of Loop only for readability; it is the
// whole dispatch switch. Cases that complete normally let the loop re-read the
// (possibly new) current frame.
func exec(thread *rtda.Thread, frame *rtda.Frame, op opcode.Opcode, opcodePc int) {
	switch op {

	// ---------- constants ----------
	case opcode.Nop:
	case opcode.AconstNull:
		frame.PushRef(nil)
	case opcode.IconstM1:
		frame.PushInt(-1)
	case opcode.Iconst0, opcode.Iconst1, opcode.Iconst2, opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		frame.PushInt(int32(op - opcode.Iconst0))
	case opcode.Lconst0, opcode.Lconst1:
		frame.PushLong(int64(op - opcode.Lconst0))
	case opcode.Fconst0, opcode.Fconst1, opcode.Fconst2:
		frame.PushFloat(float32(op - opcode.Fconst0))
	case opcode.Dconst0, opcode.Dconst1:
		frame.PushDouble(float64(op - opcode.Dconst0))
	case opcode.Bipush:
		frame.PushInt(int32(int8(frame.ReadUint8())))
	case opcode.Sipush:
		frame.PushInt(int32(frame.ReadInt16()))
	case opcode.Ldc:
		pushConstant(thread, frame, frame.Method().Owner().ConstantPool(), uint16(frame.ReadUint8()))
	case opcode.LdcW:
		pushConstant(thread, frame, frame.Method().Owner().ConstantPool(), frame.ReadUint16())
	case opcode.Ldc2W:
		cp := frame.Method().Owner().ConstantPool()
		idx := frame.ReadUint16()
		switch cp.Tag(idx) {
		case 5: // CONSTANT_Long
			frame.PushLong(cp.Long(idx))
		case 6: // CONSTANT_Double
			frame.PushDouble(cp.Double(idx))
		}

	// ---------- loads ----------
	case opcode.Iload, opcode.Lload, opcode.Fload, opcode.Dload, opcode.Aload:
		loadLocal(frame, op, int(frame.ReadUint8()))
	case opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3:
		frame.PushInt(frame.GetInt(int(op - opcode.Iload0)))
	case opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		frame.PushRef(frame.GetRef(int(op - opcode.Aload0)))
	case opcode.Lload0, opcode.Lload1, opcode.Lload2, opcode.Lload3:
		frame.PushLong(frame.GetLong(int(op - opcode.Lload0)))
	case opcode.Fload0, opcode.Fload1, opcode.Fload2, opcode.Fload3:
		frame.PushFloat(frame.GetFloat(int(op - opcode.Fload0)))
	case opcode.Dload0, opcode.Dload1, opcode.Dload2, opcode.Dload3:
		frame.PushDouble(frame.GetDouble(int(op - opcode.Dload0)))

	// ---------- stores ----------
	case opcode.Istore, opcode.Lstore, opcode.Fstore, opcode.Dstore, opcode.Astore:
		storeLocal(frame, op, int(frame.ReadUint8()))
	case opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		frame.SetInt(int(op-opcode.Istore0), frame.PopInt())
	case opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		frame.SetRef(int(op-opcode.Astore0), frame.PopRef())
	case opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3:
		frame.SetLong(int(op-opcode.Lstore0), frame.PopLong())
	case opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3:
		frame.SetFloat(int(op-opcode.Fstore0), frame.PopFloat())
	case opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		frame.SetDouble(int(op-opcode.Dstore0), frame.PopDouble())

	// ---------- array load ----------
	case opcode.Iaload, opcode.Baload, opcode.Caload, opcode.Saload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushInt(arr.ArrayElementSlot(int(i)).Num())
	case opcode.Laload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushLong(readTwoSlots(arr, int(i)))
	case opcode.Faload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushFloat(float32frombits(uint32(arr.ArrayElementSlot(int(i)).Num())))
	case opcode.Daload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushDouble(float64frombits(uint64(readTwoSlots(arr, int(i)))))
	case opcode.Aaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushRef(arr.ArrayElementSlot(int(i)).Ref())

	// ---------- array store ----------
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore:
		v := frame.PopInt()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetNum(v)
	case opcode.Lastore:
		v := frame.PopLong()
		i := frame.PopInt()
		arr := frame.PopRef()
		writeTwoSlots(arr, int(i), v)
	case opcode.Fastore:
		v := frame.PopFloat()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetNum(int32(float32bits(v)))
	case opcode.Dastore:
		v := frame.PopDouble()
		i := frame.PopInt()
		arr := frame.PopRef()
		writeTwoSlots(arr, int(i), int64(float64bits(v)))
	case opcode.Aastore:
		v := frame.PopRef()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetRef(v)

	// ---------- stack manipulation ----------
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

	// ---------- arithmetic: int ----------
	case opcode.Iadd:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a + b)
	case opcode.Isub:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a - b)
	case opcode.Imul:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a * b)
	case opcode.Idiv:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a / b) // Go integer division truncates toward zero, matching Java
	case opcode.Irem:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a % b)
	case opcode.Ineg:
		frame.PushInt(-frame.PopInt())
	case opcode.Ishl:
		s := uint32(frame.PopInt()) & 0x1f
		v := frame.PopInt()
		frame.PushInt(v << s)
	case opcode.Ishr:
		s := uint32(frame.PopInt()) & 0x1f
		v := frame.PopInt()
		frame.PushInt(v >> s) // arithmetic: signed shift
	case opcode.Iushr:
		s := uint32(frame.PopInt()) & 0x1f
		v := uint32(frame.PopInt())
		frame.PushInt(int32(v >> s)) // logical shift
	case opcode.Iand:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a & b)
	case opcode.Ior:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a | b)
	case opcode.Ixor:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a ^ b)
	case opcode.Iinc:
		idx := int(frame.ReadUint8())
		delta := int8(frame.ReadUint8())
		frame.SetInt(idx, frame.GetInt(idx)+int32(delta))

	// ---------- arithmetic: long ----------
	case opcode.Ladd:
		frame.PushLong(frame.PopLong() + frame.PopLong())
	case opcode.Lsub:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() - b)
	case opcode.Lmul:
		frame.PushLong(frame.PopLong() * frame.PopLong())
	case opcode.Ldiv:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() / b)
	case opcode.Lrem:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() % b)
	case opcode.Lneg:
		frame.PushLong(-frame.PopLong())
	case opcode.Lshl:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(frame.PopLong() << s)
	case opcode.Lshr:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(frame.PopLong() >> s)
	case opcode.Lushr:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(int64(uint64(frame.PopLong()) >> s))
	case opcode.Land:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a & b)
	case opcode.Lor:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a | b)
	case opcode.Lxor:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a ^ b)

	// ---------- arithmetic: float / double ----------
	case opcode.Fadd:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() + b)
	case opcode.Fsub:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() - b)
	case opcode.Fmul:
		frame.PushFloat(frame.PopFloat() * frame.PopFloat())
	case opcode.Fdiv:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() / b)
	case opcode.Frem:
		b := frame.PopFloat()
		frame.PushFloat(float32(remF(frame.PopFloat(), b)))
	case opcode.Fneg:
		frame.PushFloat(-frame.PopFloat())
	case opcode.Dadd:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() + b)
	case opcode.Dsub:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() - b)
	case opcode.Dmul:
		frame.PushDouble(frame.PopDouble() * frame.PopDouble())
	case opcode.Ddiv:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() / b)
	case opcode.Drem:
		b := frame.PopDouble()
		frame.PushDouble(remF64(frame.PopDouble(), b))
	case opcode.Dneg:
		frame.PushDouble(-frame.PopDouble())

	// ---------- conversions ----------
	case opcode.I2l:
		frame.PushLong(int64(frame.PopInt()))
	case opcode.I2f:
		frame.PushFloat(float32(frame.PopInt()))
	case opcode.I2d:
		frame.PushDouble(float64(frame.PopInt()))
	case opcode.L2i:
		frame.PushInt(int32(frame.PopLong()))
	case opcode.L2f:
		frame.PushFloat(float32(frame.PopLong()))
	case opcode.L2d:
		frame.PushDouble(float64(frame.PopLong()))
	case opcode.F2i:
		frame.PushInt(int32(frame.PopFloat()))
	case opcode.F2l:
		frame.PushLong(int64(frame.PopFloat()))
	case opcode.F2d:
		frame.PushDouble(float64(frame.PopFloat()))
	case opcode.D2i:
		frame.PushInt(int32(frame.PopDouble()))
	case opcode.D2l:
		frame.PushLong(int64(frame.PopDouble()))
	case opcode.D2f:
		frame.PushFloat(float32(frame.PopDouble()))
	case opcode.I2b:
		frame.PushInt(int32(int8(frame.PopInt())))
	case opcode.I2c:
		frame.PushInt(int32(uint16(frame.PopInt())))
	case opcode.I2s:
		frame.PushInt(int32(int16(frame.PopInt())))

	// ---------- comparisons ----------
	case opcode.Lcmp:
		b := frame.PopLong()
		a := frame.PopLong()
		switch {
		case a > b:
			frame.PushInt(1)
		case a < b:
			frame.PushInt(-1)
		default:
			frame.PushInt(0)
		}
	case opcode.Fcmpl, opcode.Fcmpg:
		b := frame.PopFloat()
		a := frame.PopFloat()
		frame.PushInt(cmpFloat(a, b, op == opcode.Fcmpg))
	case opcode.Dcmpl, opcode.Dcmpg:
		b := frame.PopDouble()
		a := frame.PopDouble()
		frame.PushInt(cmpDouble(a, b, op == opcode.Dcmpg))

	// ---------- conditional branches ----------
	case opcode.Ifeq:
		if frame.PopInt() == 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Ifne:
		if frame.PopInt() != 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Iflt:
		if frame.PopInt() < 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Ifge:
		if frame.PopInt() >= 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Ifgt:
		if frame.PopInt() > 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Ifle:
		if frame.PopInt() <= 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.IfIcmpeq:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a == b)
	case opcode.IfIcmpne:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a != b)
	case opcode.IfIcmplt:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a < b)
	case opcode.IfIcmpge:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a >= b)
	case opcode.IfIcmpgt:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a > b)
	case opcode.IfIcmple:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a <= b)
	case opcode.IfAcmpeq:
		b := frame.PopRef()
		a := frame.PopRef()
		condBranch(frame, opcodePc, a == b)
	case opcode.IfAcmpne:
		b := frame.PopRef()
		a := frame.PopRef()
		condBranch(frame, opcodePc, a != b)
	case opcode.Ifnull:
		if frame.PopRef() == nil {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opcode.Ifnonnull:
		if frame.PopRef() != nil {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}

	case opcode.Goto:
		branch(frame, opcodePc, int(frame.ReadInt16()))
	case opcode.GotoW:
		branch(frame, opcodePc, int(frame.ReadInt32()))
	case opcode.Tableswitch:
		tableSwitch(frame, opcodePc)
	case opcode.Lookupswitch:
		lookupSwitch(frame, opcodePc)

	// ---------- returns ----------
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

	// ---------- fields ----------
	case opcode.Getstatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		field := class.LookupField(name, desc)
		loadFieldValue(frame, class.StaticVars(), field.SlotID(), desc)
	case opcode.Putstatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		field := class.LookupField(name, desc)
		storeFieldValue(frame, class.StaticVars(), field.SlotID(), desc)
	case opcode.Getfield:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		field := class.LookupField(name, desc)
		obj := frame.PopRef()
		loadFieldValue(frame, obj.Fields(), field.SlotID(), desc)
	case opcode.Putfield:
		// Stack layout is [objref, value] with value on top, so the value must
		// be popped before the object reference.
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		field := thread.Loader().LoadClass(cls).LookupField(name, desc)
		slotID := field.SlotID()
		switch desc[0] {
		case 'J':
			v := frame.PopLong()
			obj := frame.PopRef()
			s := obj.Fields()
			s[slotID].SetNum(int32(uint64(v) >> 32))
			s[slotID+1].SetNum(int32(v))
		case 'D':
			v := frame.PopDouble()
			obj := frame.PopRef()
			bits := float64bits(v)
			s := obj.Fields()
			s[slotID].SetNum(int32(bits >> 32))
			s[slotID+1].SetNum(int32(bits))
		case 'F':
			v := frame.PopFloat()
			obj := frame.PopRef()
			obj.Fields()[slotID].SetNum(int32(float32bits(v)))
		case 'L', '[':
			v := frame.PopRef()
			obj := frame.PopRef()
			obj.Fields()[slotID].SetRef(v)
		default: // 'Z','B','C','S','I'
			v := frame.PopInt()
			obj := frame.PopRef()
			obj.Fields()[slotID].SetNum(v)
		}

	// ---------- invocations ----------
	case opcode.Invokevirtual:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		spec := class.LookupMethod(name, desc)
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			panic("catty: NullPointerException")
		}
		invokeMethod(thread, receiver.Class().LookupMethod(name, desc))
	case opcode.Invokespecial:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		invokeMethod(thread, class.LookupMethod(name, desc))
	case opcode.Invokestatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		invokeMethod(thread, class.LookupMethod(name, desc))

	// ---------- object / array creation ----------
	case opcode.New:
		idx := frame.ReadUint16()
		class := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		ensureInitialized(thread, class)
		frame.PushRef(rtda.NewObject(class))
	case opcode.Newarray:
		atype := frame.ReadUint8()
		length := int(frame.PopInt())
		frame.PushRef(newPrimitiveArray(thread, atype, length))
	case opcode.Anewarray:
		idx := frame.ReadUint16()
		length := int(frame.PopInt())
		elemName := frame.Method().Owner().ConstantPool().ClassName(idx)
		frame.PushRef(newRefArray(thread, elemName, length))
	case opcode.Arraylength:
		arr := frame.PopRef()
		frame.PushInt(int32(arr.ArrayLength()))
	case opcode.Checkcast:
		idx := frame.ReadUint16()
		target := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		obj := frame.PeekRef(0)
		if obj != nil && !obj.IsInstanceOf(target) {
			panic("catty: ClassCastException")
		}
	case opcode.Instanceof:
		idx := frame.ReadUint16()
		target := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		obj := frame.PopRef()
		if obj != nil && obj.IsInstanceOf(target) {
			frame.PushInt(1)
		} else {
			frame.PushInt(0)
		}
	case opcode.Monitorenter, opcode.Monitorexit:
		frame.PopRef() // concurrency deferred: monitors are nops in the single-threaded MVP

	default:
		panic(fmt.Sprintf("catty: opcode 0x%02x (%s) not implemented", op, opcode.Name(opcode.Opcode(op))))
	}
}
