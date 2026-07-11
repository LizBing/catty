package interpreter

import (
	"fmt"

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
		op := frame.Code()[opcodePc]
		frame.SetPC(opcodePc + 1)
		exec(thread, frame, op, opcodePc)
	}
}

// exec runs one instruction. Split out of Loop only for readability; it is the
// whole dispatch switch. Cases that complete normally let the loop re-read the
// (possibly new) current frame.
func exec(thread *rtda.Thread, frame *rtda.Frame, op byte, opcodePc int) {
	switch op {

	// ---------- constants ----------
	case opNop:
	case opAconstNull:
		frame.PushRef(nil)
	case opIconstM1:
		frame.PushInt(-1)
	case opIconst0, opIconst1, opIconst2, opIconst3, opIconst4, opIconst5:
		frame.PushInt(int32(op - opIconst0))
	case opLconst0, opLconst1:
		frame.PushLong(int64(op - opLconst0))
	case opFconst0, opFconst1, opFconst2:
		frame.PushFloat(float32(op - opFconst0))
	case opDconst0, opDconst1:
		frame.PushDouble(float64(op - opDconst0))
	case opBipush:
		frame.PushInt(int32(int8(frame.ReadUint8())))
	case opSipush:
		frame.PushInt(int32(frame.ReadInt16()))
	case opLdc:
		pushConstant(thread, frame, frame.Method().Owner().ConstantPool(), uint16(frame.ReadUint8()))
	case opLdcW:
		pushConstant(thread, frame, frame.Method().Owner().ConstantPool(), frame.ReadUint16())
	case opLdc2W:
		cp := frame.Method().Owner().ConstantPool()
		idx := frame.ReadUint16()
		switch cp.Tag(idx) {
		case 5: // CONSTANT_Long
			frame.PushLong(cp.Long(idx))
		case 6: // CONSTANT_Double
			frame.PushDouble(cp.Double(idx))
		}

	// ---------- loads ----------
	case opIload, opLload, opFload, opDload, opAload:
		loadLocal(frame, op, int(frame.ReadUint8()))
	case opIload0, opIload1, opIload2, opIload3:
		frame.PushInt(frame.GetInt(int(op - opIload0)))
	case opAload0, opAload1, opAload2, opAload3:
		frame.PushRef(frame.GetRef(int(op - opAload0)))
	case opLload0, opLload1, opLload2, opLload3:
		frame.PushLong(frame.GetLong(int(op - opLload0)))
	case opFload0, opFload1, opFload2, opFload3:
		frame.PushFloat(frame.GetFloat(int(op - opFload0)))
	case opDload0, opDload1, opDload2, opDload3:
		frame.PushDouble(frame.GetDouble(int(op - opDload0)))

	// ---------- stores ----------
	case opIstore, opLstore, opFstore, opDstore, opAstore:
		storeLocal(frame, op, int(frame.ReadUint8()))
	case opIstore0, opIstore1, opIstore2, opIstore3:
		frame.SetInt(int(op-opIstore0), frame.PopInt())
	case opAstore0, opAstore1, opAstore2, opAstore3:
		frame.SetRef(int(op-opAstore0), frame.PopRef())
	case opLstore0, opLstore1, opLstore2, opLstore3:
		frame.SetLong(int(op-opLstore0), frame.PopLong())
	case opFstore0, opFstore1, opFstore2, opFstore3:
		frame.SetFloat(int(op-opFstore0), frame.PopFloat())
	case opDstore0, opDstore1, opDstore2, opDstore3:
		frame.SetDouble(int(op-opDstore0), frame.PopDouble())

	// ---------- array load ----------
	case opIaload, opBaload, opCaload, opSaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushInt(arr.ArrayElementSlot(int(i)).Num())
	case opLaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushLong(readTwoSlots(arr, int(i)))
	case opFaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushFloat(float32frombits(uint32(arr.ArrayElementSlot(int(i)).Num())))
	case opDaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushDouble(float64frombits(uint64(readTwoSlots(arr, int(i)))))
	case opAaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		frame.PushRef(arr.ArrayElementSlot(int(i)).Ref())

	// ---------- array store ----------
	case opIastore, opBastore, opCastore, opSastore:
		v := frame.PopInt()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetNum(v)
	case opLastore:
		v := frame.PopLong()
		i := frame.PopInt()
		arr := frame.PopRef()
		writeTwoSlots(arr, int(i), v)
	case opFastore:
		v := frame.PopFloat()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetNum(int32(float32bits(v)))
	case opDastore:
		v := frame.PopDouble()
		i := frame.PopInt()
		arr := frame.PopRef()
		writeTwoSlots(arr, int(i), int64(float64bits(v)))
	case opAastore:
		v := frame.PopRef()
		i := frame.PopInt()
		arr := frame.PopRef()
		arr.ArrayElementSlot(int(i)).SetRef(v)

	// ---------- stack manipulation ----------
	case opPop:
		frame.PopSlot()
	case opPop2:
		frame.PopSlot()
		frame.PopSlot()
	case opDup:
		s := frame.PopSlot()
		frame.PushSlot(s)
		frame.PushSlot(s)
	case opDupX1:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opDupX2:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		s3 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s3)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opDup2:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s2)
		frame.PushSlot(s1)
		frame.PushSlot(s2)
		frame.PushSlot(s1)
	case opSwap:
		s1 := frame.PopSlot()
		s2 := frame.PopSlot()
		frame.PushSlot(s1)
		frame.PushSlot(s2)

	// ---------- arithmetic: int ----------
	case opIadd:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a + b)
	case opIsub:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a - b)
	case opImul:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a * b)
	case opIdiv:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a / b) // Go integer division truncates toward zero, matching Java
	case opIrem:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a % b)
	case opIneg:
		frame.PushInt(-frame.PopInt())
	case opIshl:
		s := uint32(frame.PopInt()) & 0x1f
		v := frame.PopInt()
		frame.PushInt(v << s)
	case opIshr:
		s := uint32(frame.PopInt()) & 0x1f
		v := frame.PopInt()
		frame.PushInt(v >> s) // arithmetic: signed shift
	case opIushr:
		s := uint32(frame.PopInt()) & 0x1f
		v := uint32(frame.PopInt())
		frame.PushInt(int32(v >> s)) // logical shift
	case opIand:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a & b)
	case opIor:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a | b)
	case opIxor:
		b := frame.PopInt()
		a := frame.PopInt()
		frame.PushInt(a ^ b)
	case opIinc:
		idx := int(frame.ReadUint8())
		delta := int8(frame.ReadUint8())
		frame.SetInt(idx, frame.GetInt(idx)+int32(delta))

	// ---------- arithmetic: long ----------
	case opLadd:
		frame.PushLong(frame.PopLong() + frame.PopLong())
	case opLsub:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() - b)
	case opLmul:
		frame.PushLong(frame.PopLong() * frame.PopLong())
	case opLdiv:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() / b)
	case opLrem:
		b := frame.PopLong()
		frame.PushLong(frame.PopLong() % b)
	case opLneg:
		frame.PushLong(-frame.PopLong())
	case opLshl:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(frame.PopLong() << s)
	case opLshr:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(frame.PopLong() >> s)
	case opLushr:
		s := uint32(frame.PopInt()) & 0x3f
		frame.PushLong(int64(uint64(frame.PopLong()) >> s))
	case opLand:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a & b)
	case opLor:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a | b)
	case opLxor:
		b := frame.PopLong()
		a := frame.PopLong()
		frame.PushLong(a ^ b)

	// ---------- arithmetic: float / double ----------
	case opFadd:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() + b)
	case opFsub:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() - b)
	case opFmul:
		frame.PushFloat(frame.PopFloat() * frame.PopFloat())
	case opFdiv:
		b := frame.PopFloat()
		frame.PushFloat(frame.PopFloat() / b)
	case opFrem:
		b := frame.PopFloat()
		frame.PushFloat(float32(remF(frame.PopFloat(), b)))
	case opFneg:
		frame.PushFloat(-frame.PopFloat())
	case opDadd:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() + b)
	case opDsub:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() - b)
	case opDmul:
		frame.PushDouble(frame.PopDouble() * frame.PopDouble())
	case opDdiv:
		b := frame.PopDouble()
		frame.PushDouble(frame.PopDouble() / b)
	case opDrem:
		b := frame.PopDouble()
		frame.PushDouble(remF64(frame.PopDouble(), b))
	case opDneg:
		frame.PushDouble(-frame.PopDouble())

	// ---------- conversions ----------
	case opI2l:
		frame.PushLong(int64(frame.PopInt()))
	case opI2f:
		frame.PushFloat(float32(frame.PopInt()))
	case opI2d:
		frame.PushDouble(float64(frame.PopInt()))
	case opL2i:
		frame.PushInt(int32(frame.PopLong()))
	case opL2f:
		frame.PushFloat(float32(frame.PopLong()))
	case opL2d:
		frame.PushDouble(float64(frame.PopLong()))
	case opF2i:
		frame.PushInt(int32(frame.PopFloat()))
	case opF2l:
		frame.PushLong(int64(frame.PopFloat()))
	case opF2d:
		frame.PushDouble(float64(frame.PopFloat()))
	case opD2i:
		frame.PushInt(int32(frame.PopDouble()))
	case opD2l:
		frame.PushLong(int64(frame.PopDouble()))
	case opD2f:
		frame.PushFloat(float32(frame.PopDouble()))
	case opI2b:
		frame.PushInt(int32(int8(frame.PopInt())))
	case opI2c:
		frame.PushInt(int32(uint16(frame.PopInt())))
	case opI2s:
		frame.PushInt(int32(int16(frame.PopInt())))

	// ---------- comparisons ----------
	case opLcmp:
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
	case opFcmpl, opFcmpg:
		b := frame.PopFloat()
		a := frame.PopFloat()
		frame.PushInt(cmpFloat(a, b, op == opFcmpg))
	case opDcmpl, opDcmpg:
		b := frame.PopDouble()
		a := frame.PopDouble()
		frame.PushInt(cmpDouble(a, b, op == opDcmpg))

	// ---------- conditional branches ----------
	case opIfeq:
		if frame.PopInt() == 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfne:
		if frame.PopInt() != 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIflt:
		if frame.PopInt() < 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfge:
		if frame.PopInt() >= 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfgt:
		if frame.PopInt() > 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfle:
		if frame.PopInt() <= 0 {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfIcmpeq:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a == b)
	case opIfIcmpne:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a != b)
	case opIfIcmplt:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a < b)
	case opIfIcmpge:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a >= b)
	case opIfIcmpgt:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a > b)
	case opIfIcmple:
		b := frame.PopInt()
		a := frame.PopInt()
		condBranch(frame, opcodePc, a <= b)
	case opIfAcmpeq:
		b := frame.PopRef()
		a := frame.PopRef()
		condBranch(frame, opcodePc, a == b)
	case opIfAcmpne:
		b := frame.PopRef()
		a := frame.PopRef()
		condBranch(frame, opcodePc, a != b)
	case opIfnull:
		if frame.PopRef() == nil {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}
	case opIfnonnull:
		if frame.PopRef() != nil {
			branch(frame, opcodePc, int(frame.ReadInt16()))
		} else {
			frame.ReadInt16()
		}

	case opGoto:
		branch(frame, opcodePc, int(frame.ReadInt16()))
	case opGotoW:
		branch(frame, opcodePc, int(frame.ReadInt32()))
	case opTableswitch:
		tableSwitch(frame, opcodePc)
	case opLookupswitch:
		lookupSwitch(frame, opcodePc)

	// ---------- returns ----------
	case opReturn:
		thread.PopFrame()
	case opIreturn:
		returnInt(frame, thread)
	case opAreturn:
		returnRef(frame, thread)
	case opLreturn:
		returnLong(frame, thread)
	case opFreturn:
		returnFloat(frame, thread)
	case opDreturn:
		returnDouble(frame, thread)

	// ---------- fields ----------
	case opGetstatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		field := class.LookupField(name, desc)
		loadFieldValue(frame, class.StaticVars(), field.SlotID(), desc)
	case opPutstatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		field := class.LookupField(name, desc)
		storeFieldValue(frame, class.StaticVars(), field.SlotID(), desc)
	case opGetfield:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		field := class.LookupField(name, desc)
		obj := frame.PopRef()
		loadFieldValue(frame, obj.Fields(), field.SlotID(), desc)
	case opPutfield:
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
	case opInvokevirtual:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		spec := class.LookupMethod(name, desc)
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			panic("catty: NullPointerException")
		}
		invokeMethod(thread, receiver.Class().LookupMethod(name, desc))
	case opInvokespecial:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		invokeMethod(thread, class.LookupMethod(name, desc))
	case opInvokestatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		ensureInitialized(thread, class)
		invokeMethod(thread, class.LookupMethod(name, desc))

	// ---------- object / array creation ----------
	case opNew:
		idx := frame.ReadUint16()
		class := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		ensureInitialized(thread, class)
		frame.PushRef(rtda.NewObject(class))
	case opNewarray:
		atype := frame.ReadUint8()
		length := int(frame.PopInt())
		frame.PushRef(newPrimitiveArray(thread, atype, length))
	case opAnewarray:
		idx := frame.ReadUint16()
		length := int(frame.PopInt())
		elemName := frame.Method().Owner().ConstantPool().ClassName(idx)
		frame.PushRef(newRefArray(thread, elemName, length))
	case opArraylength:
		arr := frame.PopRef()
		frame.PushInt(int32(arr.ArrayLength()))
	case opCheckcast:
		idx := frame.ReadUint16()
		target := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		obj := frame.PeekRef(0)
		if obj != nil && !obj.IsInstanceOf(target) {
			panic("catty: ClassCastException")
		}
	case opInstanceof:
		idx := frame.ReadUint16()
		target := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		obj := frame.PopRef()
		if obj != nil && obj.IsInstanceOf(target) {
			frame.PushInt(1)
		} else {
			frame.PushInt(0)
		}
	case opMonitorenter, opMonitorexit:
		frame.PopRef() // concurrency deferred: monitors are nops in the single-threaded MVP

	default:
		panic(fmt.Sprintf("catty: opcode 0x%02x (%s) not implemented", op, opcodeName[op]))
	}
}
