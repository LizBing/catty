package lowering

import (
	"catty/classfile"
	"catty/opcode"
	"catty/rtda"
)

// SlotType is the verified type of one operand-stack slot, in the granularity
// the emitter needs to pick a Go type. Category-2 values (long/double) span two
// slots: [TypeLong, TypeTop] / [TypeDouble, TypeTop].
type SlotType uint8

const (
	TypeTop    SlotType = iota // unused slot or the 2nd slot of a category-2 value
	TypeInt                    // byte/char/short/boolean/int → int32
	TypeLong                   // → int64 (followed by TypeTop)
	TypeFloat                  // → float32
	TypeDouble                 // → float64 (followed by TypeTop)
	TypeRef                    // object/array/null → *rtda.Object
)

// typeDataflow fills IRInst.InTypes: the operand-stack slot types at each
// instruction's entry. It is a single linear pass over `starts` (decoded order),
// tracking only the operand stack — loads are opcode-derived (iload reads an
// int, aload a ref, … the verifier guarantees this), so locals need not be
// tracked. At each StackMapTable frame pc the stack is reset to the frame's
// stack types, pinning every merge (no type-lattice merge logic).
//
// If the method has no StackMapTable, InTypes is left nil — the A2 emitter will
// treat types as unknown and skip AOT for that method.
func typeDataflow(method *rtda.Method, ir *IR, starts []int, cp *classfile.ConstantPool) {
	framesByPc := map[int][]SlotType{}
	if smt := method.StackMap(); smt != nil {
		for _, f := range smt.Reconstruct(initialVerifLocals(method)) {
			framesByPc[f.Offset] = verifToSlots(f.Stack)
		}
	}
	var stack []SlotType
	for _, pc := range starts {
		if fs, ok := framesByPc[pc]; ok {
			stack = fs // reset at a merge — the frame's stack is authoritative
		}
		inst := &ir.Insts[pc]
		inst.InTypes = append([]SlotType(nil), stack...)
		applyTypeEffect(inst, &stack, cp)
	}
}

// initialVerifLocals builds the implicit-entry-frame locals (VerifType form, one
// entry per logical local — long/double are a single ItemLong/ItemDouble) from
// the method descriptor. Needed only to seed StackMapTable.Reconstruct's delta
// walk; the resulting frames' Stack is what typeDataflow uses.
func initialVerifLocals(method *rtda.Method) []classfile.VerifType {
	md := rtda.ParseMethodDescriptor(method.Descriptor())
	var locals []classfile.VerifType
	if !method.IsStatic() {
		locals = append(locals, classfile.VerifType{Tag: classfile.ItemObject}) // `this`
	}
	for _, p := range md.ParameterTypes {
		locals = append(locals, paramVerif(p))
	}
	return locals
}

func paramVerif(p string) classfile.VerifType {
	switch p {
	case "J":
		return classfile.VerifType{Tag: classfile.ItemLong}
	case "D":
		return classfile.VerifType{Tag: classfile.ItemDouble}
	case "F":
		return classfile.VerifType{Tag: classfile.ItemFloat}
	default: // I, B, C, S, Z, L...;, [...
		return classfile.VerifType{Tag: classfile.ItemObject} // refs; ints use Integer below
	}
}

// verifToSlots expands a frame's verification types (one per logical slot, with
// long/double as single entries) into per-slot SlotTypes (long/double → 2 slots).
func verifToSlots(vts []classfile.VerifType) []SlotType {
	out := make([]SlotType, 0, len(vts))
	for _, vt := range vts {
		switch vt.Tag {
		case classfile.ItemLong:
			out = append(out, TypeLong, TypeTop)
		case classfile.ItemDouble:
			out = append(out, TypeDouble, TypeTop)
		case classfile.ItemInteger:
			out = append(out, TypeInt)
		case classfile.ItemFloat:
			out = append(out, TypeFloat)
		case classfile.ItemTop:
			out = append(out, TypeTop)
		default: // Null, Object, UninitializedThis, Uninitialized → ref
			out = append(out, TypeRef)
		}
	}
	return out
}

// applyTypeEffect advances the operand-stack types past one instruction. Loads
// push opcode-derived types; stores/branches/returns pop; arithmetic/conversions
// pop inputs and push the result type; fields/invokes resolve types from the
// constant-pool descriptor.
func applyTypeEffect(inst *IRInst, stack *[]SlotType, cp *classfile.ConstantPool) {
	s := stack
	push := func(t SlotType) { *s = append(*s, t) }
	push2 := func(t SlotType) { *s = append(*s, t, TypeTop) }
	pop := func() { *s = (*s)[:len(*s)-1] }
	pop2 := func() { *s = (*s)[:len(*s)-2] }

	switch inst.Op {
	// --- constants ---
	case opcode.AconstNull:
		push(TypeRef)
	case opcode.IconstM1, opcode.Iconst0, opcode.Iconst1, opcode.Iconst2,
		opcode.Iconst3, opcode.Iconst4, opcode.Iconst5:
		push(TypeInt)
	case opcode.Lconst0, opcode.Lconst1:
		push2(TypeLong)
	case opcode.Fconst0, opcode.Fconst1, opcode.Fconst2:
		push(TypeFloat)
	case opcode.Dconst0, opcode.Dconst1:
		push2(TypeDouble)
	case opcode.Bipush, opcode.Sipush:
		push(TypeInt)
	case opcode.Ldc, opcode.LdcW:
		switch cp.Tag(inst.Index) {
		case classfile.ConstantInteger:
			push(TypeInt)
		case classfile.ConstantFloat:
			push(TypeFloat)
		case classfile.ConstantString, classfile.ConstantClass:
			push(TypeRef)
		}
	case opcode.Ldc2W:
		switch cp.Tag(inst.Index) {
		case classfile.ConstantLong:
			push2(TypeLong)
		case classfile.ConstantDouble:
			push2(TypeDouble)
		}

	// --- loads (opcode-derived; the verifier guarantees the local's type) ---
	case opcode.Iload, opcode.Iload0, opcode.Iload1, opcode.Iload2, opcode.Iload3:
		push(TypeInt)
	case opcode.Fload, opcode.Fload0, opcode.Fload1, opcode.Fload2, opcode.Fload3:
		push(TypeFloat)
	case opcode.Aload, opcode.Aload0, opcode.Aload1, opcode.Aload2, opcode.Aload3:
		push(TypeRef)
	case opcode.Lload, opcode.Lload0, opcode.Lload1, opcode.Lload2, opcode.Lload3:
		push2(TypeLong)
	case opcode.Dload, opcode.Dload0, opcode.Dload1, opcode.Dload2, opcode.Dload3:
		push2(TypeDouble)

	// --- stores ---
	case opcode.Istore, opcode.Istore0, opcode.Istore1, opcode.Istore2, opcode.Istore3:
		pop()
	case opcode.Fstore, opcode.Fstore0, opcode.Fstore1, opcode.Fstore2, opcode.Fstore3:
		pop()
	case opcode.Astore, opcode.Astore0, opcode.Astore1, opcode.Astore2, opcode.Astore3:
		pop()
	case opcode.Lstore, opcode.Lstore0, opcode.Lstore1, opcode.Lstore2, opcode.Lstore3:
		pop2()
	case opcode.Dstore, opcode.Dstore0, opcode.Dstore1, opcode.Dstore2, opcode.Dstore3:
		pop2()

	// --- array load ---
	case opcode.Iaload, opcode.Baload, opcode.Caload, opcode.Saload:
		pop2()
		push(TypeInt)
	case opcode.Faload:
		pop2()
		push(TypeFloat)
	case opcode.Aaload:
		pop2()
		push(TypeRef)
	case opcode.Laload:
		pop2()
		push2(TypeLong)
	case opcode.Daload:
		pop2()
		push2(TypeDouble)

	// --- array store ---
	case opcode.Iastore, opcode.Bastore, opcode.Castore, opcode.Sastore, opcode.Fastore, opcode.Aastore:
		*s = (*s)[:len(*s)-3]
	case opcode.Lastore, opcode.Dastore:
		*s = (*s)[:len(*s)-4]

	// --- stack shuffles (mirror the value movement on types) ---
	case opcode.Pop:
		pop()
	case opcode.Pop2:
		pop2()
	case opcode.Dup:
		t := (*s)[len(*s)-1]
		push(t)
		push(t)
	case opcode.DupX1:
		t1 := (*s)[len(*s)-1]
		t2 := (*s)[len(*s)-2]
		pop2()
		push(t1)
		push(t2)
		push(t1)
	case opcode.DupX2:
		t1 := (*s)[len(*s)-1]
		t2 := (*s)[len(*s)-2]
		t3 := (*s)[len(*s)-3]
		*s = (*s)[:len(*s)-3]
		push(t1)
		push(t3)
		push(t2)
		push(t1)
	case opcode.Dup2:
		t1 := (*s)[len(*s)-1]
		t2 := (*s)[len(*s)-2]
		pop2()
		push(t2)
		push(t1)
		push(t2)
		push(t1)
	case opcode.Swap:
		t1 := (*s)[len(*s)-1]
		t2 := (*s)[len(*s)-2]
		pop2()
		push(t1)
		push(t2)

	// --- int arithmetic ---
	case opcode.Iadd, opcode.Isub, opcode.Imul, opcode.Idiv, opcode.Irem,
		opcode.Iand, opcode.Ior, opcode.Ixor, opcode.Ishl, opcode.Ishr, opcode.Iushr:
		pop2()
		push(TypeInt)
	case opcode.Ineg:
		pop()
		push(TypeInt)
	case opcode.Iinc:
		// no stack effect

	// --- long arithmetic ---
	case opcode.Ladd, opcode.Lsub, opcode.Lmul, opcode.Ldiv, opcode.Lrem,
		opcode.Land, opcode.Lor, opcode.Lxor:
		*s = (*s)[:len(*s)-4]
		push2(TypeLong)
	case opcode.Lshl, opcode.Lshr, opcode.Lushr:
		*s = (*s)[:len(*s)-3] // int shift count + long
		push2(TypeLong)
	case opcode.Lneg:
		pop2()
		push2(TypeLong)

	// --- float / double arithmetic ---
	case opcode.Fadd, opcode.Fsub, opcode.Fmul, opcode.Fdiv, opcode.Frem:
		pop2()
		push(TypeFloat)
	case opcode.Fneg:
		pop()
		push(TypeFloat)
	case opcode.Dadd, opcode.Dsub, opcode.Dmul, opcode.Ddiv, opcode.Drem:
		*s = (*s)[:len(*s)-4]
		push2(TypeDouble)
	case opcode.Dneg:
		pop2()
		push2(TypeDouble)

	// --- conversions ---
	case opcode.I2l:
		pop()
		push2(TypeLong)
	case opcode.I2f:
		pop()
		push(TypeFloat)
	case opcode.I2d:
		pop()
		push2(TypeDouble)
	case opcode.L2i:
		pop2()
		push(TypeInt)
	case opcode.L2f:
		pop2()
		push(TypeFloat)
	case opcode.L2d:
		pop2()
		push2(TypeDouble)
	case opcode.F2i:
		pop()
		push(TypeInt)
	case opcode.F2l:
		pop()
		push2(TypeLong)
	case opcode.F2d:
		pop()
		push2(TypeDouble)
	case opcode.D2i:
		pop2()
		push(TypeInt)
	case opcode.D2l:
		pop2()
		push2(TypeLong)
	case opcode.D2f:
		pop2()
		push(TypeFloat)
	case opcode.I2b, opcode.I2c, opcode.I2s:
		pop()
		push(TypeInt)

	// --- comparisons ---
	case opcode.Lcmp:
		*s = (*s)[:len(*s)-4]
		push(TypeInt)
	case opcode.Fcmpl, opcode.Fcmpg:
		pop2()
		push(TypeInt)
	case opcode.Dcmpl, opcode.Dcmpg:
		*s = (*s)[:len(*s)-4]
		push(TypeInt)

	// --- branches (pop compared operands) ---
	case opcode.Ifeq, opcode.Ifne, opcode.Iflt, opcode.Ifge, opcode.Ifgt, opcode.Ifle,
		opcode.Ifnull, opcode.Ifnonnull:
		pop()
	case opcode.IfIcmpeq, opcode.IfIcmpne, opcode.IfIcmplt, opcode.IfIcmpge,
		opcode.IfIcmpgt, opcode.IfIcmple, opcode.IfAcmpeq, opcode.IfAcmpne:
		pop2()
	case opcode.Goto, opcode.GotoW:
		// no stack effect
	case opcode.Tableswitch, opcode.Lookupswitch:
		pop() // the key

	// --- returns ---
	case opcode.Ireturn, opcode.Freturn, opcode.Areturn:
		pop()
	case opcode.Lreturn, opcode.Dreturn:
		pop2()
	case opcode.Return:
		// no stack effect

	// --- fields ---
	case opcode.Getstatic:
		pushFieldType(inst, s, cp)
	case opcode.Putstatic:
		popFieldType(inst, s, cp)
	case opcode.Getfield:
		pop() // objectref
		pushFieldType(inst, s, cp)
	case opcode.Putfield:
		popFieldType(inst, s, cp) // value
		pop()                     // objectref

	// --- invokes ---
	case opcode.Invokevirtual, opcode.Invokespecial, opcode.Invokestatic, opcode.Invokeinterface:
		_, _, desc := cp.MemberRef(inst.Index)
		md := rtda.ParseMethodDescriptor(desc)
		// pop args (+ this for instance methods)
		n := md.ArgSlots()
		if inst.Op != opcode.Invokestatic {
			n++
		}
		*s = (*s)[:len(*s)-n]
		pushReturnType(md.ReturnType, s)

	// --- object / array / misc ---
	case opcode.New:
		push(TypeRef)
	case opcode.Newarray, opcode.Anewarray:
		pop() // length
		push(TypeRef)
	case opcode.Arraylength:
		pop()
		push(TypeInt)
	case opcode.Athrow:
		pop()
	case opcode.Checkcast:
		// consumes objectref, pushes the same reference — net zero
	case opcode.Instanceof:
		pop()
		push(TypeInt)
	case opcode.Monitorenter, opcode.Monitorexit:
		pop()
	}
}

// pushFieldType pushes the type of a getstatic/getfield result.
func pushFieldType(inst *IRInst, s *[]SlotType, cp *classfile.ConstantPool) {
	t, cat2 := descToType(fieldDescriptor(cp, inst.Index))
	if cat2 {
		*s = append(*s, t, TypeTop)
	} else {
		*s = append(*s, t)
	}
}

// popFieldType pops a putstatic/putfield value (1 or 2 slots by field category).
func popFieldType(inst *IRInst, s *[]SlotType, cp *classfile.ConstantPool) {
	_, cat2 := descToType(fieldDescriptor(cp, inst.Index))
	if cat2 {
		*s = (*s)[:len(*s)-2]
	} else {
		*s = (*s)[:len(*s)-1]
	}
}

func pushReturnType(ret string, s *[]SlotType) {
	t, cat2 := descToType(ret)
	if ret == "" || ret == "V" {
		return
	}
	if cat2 {
		*s = append(*s, t, TypeTop)
	} else {
		*s = append(*s, t)
	}
}

func fieldDescriptor(cp *classfile.ConstantPool, index uint16) string {
	_, _, desc := cp.MemberRef(index)
	return desc
}

// descToType maps a field/return descriptor to its slot type + category.
func descToType(desc string) (SlotType, bool) {
	switch desc {
	case "J":
		return TypeLong, true
	case "D":
		return TypeDouble, true
	case "F":
		return TypeFloat, false
	case "I", "B", "C", "S", "Z":
		return TypeInt, false
	default: // L...; or [
		return TypeRef, false
	}
}
