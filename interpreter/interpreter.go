package interpreter

import (
	"fmt"
	"os"

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
		if thread.HasException() {
			handleException(thread, opcodePc)
		}
	}
}

// handleException searches exception tables frame-by-frame for a handler
// matching the throw PC and exception type. If found, it clears the operand
// stack, pushes the exception object, and jumps to the handler. If no frame
// catches it, the exception is uncaught — print a message and exit.
func handleException(thread *rtda.Thread, throwPC int) {
	excObj := thread.ClearException()

	for !thread.IsStackEmpty() {
		frame := thread.CurrentFrame()
		method := frame.Method()

		for _, entry := range method.ExceptionTable() {
			if throwPC >= entry.StartPc() && throwPC < entry.EndPc() {
				// catchType "" = catch-all (finally); otherwise check type.
				if entry.CatchType() == "" || excObj.IsInstanceOf(thread.Loader().LoadClass(entry.CatchType())) {
					// Found a handler: clear operand stack, push exception, jump.
					frame.ClearStack()
					frame.PushRef(excObj)
					frame.SetPC(entry.HandlerPc())
					return
				}
			}
		}

		// Not caught in this frame — pop it and try the caller.
		thread.PopFrame()
		if !thread.IsStackEmpty() {
			throwPC = thread.CurrentFrame().PC() - 1 // PC was advanced past the invoke opcode
		}
	}

	// Uncaught exception — print and exit.
	fmt.Fprintf(os.Stderr, "Exception in thread \"main\" %s", javaClassName(excObj.Class().Name()))
	msgSlot := findDetailMessage(excObj)
	if msgSlot >= 0 && excObj.GetRefCell(msgSlot) != nil {
		s := stringValueFromObj(excObj.GetRefCell(msgSlot))
		if s != "" {
			fmt.Fprintf(os.Stderr, ": %s", s)
		}
	}
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

// findDetailMessage returns the slot offset of Throwable's detailMessage, or -1.
func findDetailMessage(obj *rtda.Object) int {
	for cls := obj.Class(); cls != nil; cls = cls.SuperClass() {
		if f := cls.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
			return int(f.SlotID())
		}
	}
	return -1
}

// stringValueFromObj extracts a human-readable Go string from a String object
// for diagnostic output (uncaught exception handler). Returns "" if nil or
// not a StringValue.
func stringValueFromObj(obj *rtda.Object) string {
	if obj == nil {
		return ""
	}
	if sv, ok := obj.Extra().(*rtda.StringValue); ok {
		return sv.GoString()
	}
	return ""
}

func javaClassName(internal string) string {
	out := make([]byte, len(internal))
	for i := 0; i < len(internal); i++ {
		if internal[i] == '/' {
			out[i] = '.'
		} else {
			out[i] = internal[i]
		}
	}
	return string(out)
}

// throwRuntime creates an exception object of the given class with an optional
// message and signals it on the thread. Used for runtime exceptions (NPE,
// ArithmeticException, ClassCastException, etc.).
func throwRuntime(thread *rtda.Thread, pc int, className, message string) {
	cls := thread.Loader().LoadClass(className)
	obj := rtda.NewObject(cls)
	if message != "" {
		// Set detailMessage on the Throwable ancestor.
		for c := cls; c != nil; c = c.SuperClass() {
			if f := c.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
				// Build a StringValue from the Go string message.
				strClass := thread.Loader().LoadClass("java/lang/String")
				msgObj := rtda.NewObject(strClass)
				msgObj.SetExtra(rtda.NewStringValue(goStringToUTF16Units(message)))
				obj.SetRefCell(int(f.SlotID()), msgObj)
				break
			}
		}
	}
	thread.Throw(obj, pc)
}

// goStringToUTF16Units converts a Go string to UTF-16 code units.
func goStringToUTF16Units(s string) []uint16 {
	if s == "" {
		return []uint16{}
	}
	ascii := true
	for _, r := range s {
		if r >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		units := make([]uint16, len(s))
		for i, b := range []byte(s) {
			units[i] = uint16(b)
		}
		return units
	}
	var units []uint16
	for _, r := range s {
		if r < 0x10000 {
			units = append(units, uint16(r))
		} else {
			r -= 0x10000
			units = append(units, uint16((r>>10)&0x3FF)+0xD800)
			units = append(units, uint16(r&0x3FF)+0xDC00)
		}
	}
	return units
}

// checkArrayBounds verifies array access bounds, throwing the appropriate
// exception if out of range. Returns true if an exception was thrown (caller
// should return immediately).
func checkArrayBounds(thread *rtda.Thread, pc int, arr *rtda.Object, index int) bool {
	length := arr.ArrayLength()
	if index < 0 || index >= length {
		msg := "Index " + itoaInt(index) + " out of bounds for length " + itoaInt(length)
		throwRuntime(thread, pc, "java/lang/ArrayIndexOutOfBoundsException", msg)
		return true
	}
	return false
}

func itoaInt(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
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
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		frame.PushInt(arr.GetIntCell(int(i)))
	case opcode.Laload:
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		frame.PushLong(readTwoSlots(arr, int(i)))
	case opcode.Faload:
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		frame.PushFloat(float32frombits(uint32(arr.GetIntCell(int(i)))))
	case opcode.Daload:
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		frame.PushDouble(float64frombits(uint64(readTwoSlots(arr, int(i)))))
	case opcode.Aaload:
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		frame.PushRef(arr.GetRefCell(int(i)))

	// ---------- array store ----------
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore:
		v := frame.PopInt()
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		arr.SetIntCell(int(i), v)
	case opcode.Lastore:
		v := frame.PopLong()
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		writeTwoSlots(arr, int(i), v)
	case opcode.Fastore:
		v := frame.PopFloat()
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		arr.SetIntCell(int(i), int32(float32bits(v)))
	case opcode.Dastore:
		v := frame.PopDouble()
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		writeTwoSlots(arr, int(i), int64(float64bits(v)))
	case opcode.Aastore:
		v := frame.PopRef()
		i := frame.PopInt()
		arr := frame.PopRef()
		if arr == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		if err := checkArrayBounds(thread, opcodePc, arr, int(i)); err {
			return
		}
		arr.SetRefCell(int(i), v)

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
		if b == 0 {
			throwRuntime(thread, opcodePc, "java/lang/ArithmeticException", "/ by zero")
			return
		}
		frame.PushInt(a / b) // Go integer division truncates toward zero, matching Java
	case opcode.Irem:
		b := frame.PopInt()
		a := frame.PopInt()
		if b == 0 {
			throwRuntime(thread, opcodePc, "java/lang/ArithmeticException", "/ by zero")
			return
		}
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
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		ensureInitialized(thread, field.Owner())
		loadStaticField(frame, field.Owner(), field.SlotID(), desc)
	case opcode.Putstatic:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		ensureInitialized(thread, field.Owner())
		storeStaticField(frame, field.Owner(), field.SlotID(), desc)
	case opcode.Getfield:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		referencedClass := thread.Loader().LoadClass(cls)
		field := referencedClass.LookupField(name, desc)
		obj := frame.PopRef()
		loadInstanceField(frame, obj, field.SlotID(), desc)
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
			obj.SetLongCell(int(slotID), v)
		case 'D':
			v := frame.PopDouble()
			obj := frame.PopRef()
			obj.SetDoubleCell(int(slotID), v)
		case 'F':
			v := frame.PopFloat()
			obj := frame.PopRef()
			obj.SetFloatCell(int(slotID), v)
		case 'L', '[':
			v := frame.PopRef()
			obj := frame.PopRef()
			obj.SetRefCell(int(slotID), v)
		default: // 'Z','B','C','S','I'
			v := frame.PopInt()
			obj := frame.PopRef()
			obj.SetIntCell(int(slotID), v)
		}

	// ---------- invocations ----------
	case opcode.Invokevirtual:
		idx := frame.ReadUint16()
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		class := thread.Loader().LoadClass(cls)
		spec := class.LookupMethod(name, desc)
		if spec == nil {
			throwRuntime(thread, opcodePc, "java/lang/NoSuchMethodError", cls+"."+name+desc)
			return
		}
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
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
		m := class.LookupMethod(name, desc)
		if m == nil {
			throwRuntime(thread, opcodePc, "java/lang/NoSuchMethodError", cls+"."+name+desc)
			return
		}
		invokeMethod(thread, m)
	case opcode.Invokeinterface:
		idx := frame.ReadUint16()
		frame.ReadUint8() // count — historical, ignored
		frame.ReadUint8() // 0   — historical, ignored
		cls, name, desc := frame.Method().Owner().ConstantPool().MemberRef(idx)
		spec := thread.Loader().LoadClass(cls).LookupMethod(name, desc)
		if spec == nil {
			throwRuntime(thread, opcodePc, "java/lang/NoSuchMethodError", cls+"."+name+desc)
			return
		}
		receiver := frame.PeekRef(int(spec.ArgSlotCount()))
		if receiver == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		invokeMethod(thread, receiver.Class().LookupMethod(name, desc))

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
	case opcode.Multianewarray:
		idx := frame.ReadUint16()
		dims := int(frame.ReadUint8())
		className := frame.Method().Owner().ConstantPool().ClassName(idx)
		class := thread.Loader().LoadClass(className)
		// Pop `dims` dimension sizes from the stack (top = last dimension).
		sizes := make([]int, dims)
		for i := dims - 1; i >= 0; i-- {
			sizes[i] = int(frame.PopInt())
		}
		frame.PushRef(rtda.NewMultiArray(class, sizes, thread.Loader()))
	case opcode.Arraylength:
		arr := frame.PopRef()
		frame.PushInt(int32(arr.ArrayLength()))
	case opcode.Checkcast:
		idx := frame.ReadUint16()
		target := thread.Loader().LoadClass(frame.Method().Owner().ConstantPool().ClassName(idx))
		obj := frame.PeekRef(0)
		if obj != nil && !obj.IsInstanceOf(target) {
			throwRuntime(thread, opcodePc, "java/lang/ClassCastException", "")
			return
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

	case opcode.Athrow:
		excObj := frame.PopRef()
		if excObj == nil {
			throwRuntime(thread, opcodePc, "java/lang/NullPointerException", "")
			return
		}
		thread.Throw(excObj, opcodePc)

	case opcode.Wide:
		// wide extends local-variable indices to u2. Read the modified opcode
		// and u2 index, then execute the non-wide form manually.
		modOp := opcode.Opcode(frame.ReadUint8())
		wideIdx := int(frame.ReadUint16())
		switch modOp {
		case opcode.Iload:
			frame.PushInt(frame.GetInt(wideIdx))
		case opcode.Fload:
			frame.PushFloat(frame.GetFloat(wideIdx))
		case opcode.Aload:
			frame.PushRef(frame.GetRef(wideIdx))
		case opcode.Lload:
			frame.PushLong(frame.GetLong(wideIdx))
		case opcode.Dload:
			frame.PushDouble(frame.GetDouble(wideIdx))
		case opcode.Istore:
			frame.SetInt(wideIdx, frame.PopInt())
		case opcode.Fstore:
			frame.SetFloat(wideIdx, frame.PopFloat())
		case opcode.Astore:
			frame.SetRef(wideIdx, frame.PopRef())
		case opcode.Lstore:
			frame.SetLong(wideIdx, frame.PopLong())
		case opcode.Dstore:
			frame.SetDouble(wideIdx, frame.PopDouble())
		case opcode.Iinc:
			wideConst := int32(int16(frame.ReadUint16()))
			frame.SetInt(wideIdx, frame.GetInt(wideIdx)+wideConst)
		}

	default:
		panic(fmt.Sprintf("catty: opcode 0x%02x (%s) not implemented", op, opcode.Name(opcode.Opcode(op))))
	}
}
