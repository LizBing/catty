package rt

import (
	"fmt"
	"math"
	"runtime"
)

// execute interprets one verified method body. Call has already pushed the
// logical frame; the interpreter owns only locals and operand slots.
func (r *Runtime) execute(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	locals := make([]Value, m.MaxLocals)
	if m.Flags&AccStatic == 0 {
		if receiver == nil {
			return ThrowOutcome(r.NewThrowable(t, "java/lang/NullPointerException"))
		}
		if m.MaxLocals > 0 {
			locals[0] = RefValue(receiver)
		}
	}
	base := 0
	if m.Flags&AccStatic == 0 {
		base = 1
	}
	for i, a := range args {
		slot := base + m.Type.ParamSlot[i]
		if slot < len(locals) {
			locals[slot] = a
		}
	}
	stack := make([]Value, 0, m.MaxStack)
	ip := 0

	push := func(v Value) { stack = append(stack, v) }
	pop := func() Value {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		return v
	}
	popInt := func() int32 { return AsInt(pop()) }
	popLong := func() int64 { return AsLong(pop()) }
	popFloat := func() float32 { return AsFloat(pop()) }
	popDouble := func() float64 { return AsDouble(pop()) }
	popRef := func() *Object {
		v := pop()
		if v.Kind != KindRef {
			return nil
		}
		return v.Ref
	}
	branch := func(target int32) bool {
		idx, ok := m.pcIndex[target]
		if !ok {
			panic(fmt.Sprintf("catty: branch target %d in %s.%s%s is not an instruction", target, m.Owner.Name, m.Name, m.Descriptor))
		}
		ip = idx
		return true
	}
	abort := func(pc int32, thrown *Object) (Outcome, bool) {
		if hpc, ok := t.HandleThrow(m, pc, thrown); ok {
			stack = stack[:0]
			stack = append(stack, RefValue(thrown))
			idx, ok := m.pcIndex[hpc]
			if !ok {
				panic(fmt.Sprintf("catty: handler offset %d in %s.%s%s is not an instruction", hpc, m.Owner.Name, m.Name, m.Descriptor))
			}
			ip = idx
			return Outcome{}, true
		}
		return ThrowOutcome(thrown), false
	}
	storeLocal := func(idx int32, v Value) {
		i := int(idx)
		if i > 0 && (locals[i-1].Kind == KindLong || locals[i-1].Kind == KindDouble) {
			locals[i-1] = Value{}
		}
		locals[i] = v
		if v.Width() == 2 && i+1 < len(locals) {
			locals[i+1] = Value{}
		}
	}

	for ip < len(m.Insns) {
		insn := m.Insns[ip]
		pc := insn.PC
		t.SetPC(pc)
		switch insn.Op {
		case OpNop:
		case OpAconstNull:
			push(RefValue(nil))
		case OpIconstM1, OpIconst0, OpIconst1, OpIconst2, OpIconst3, OpIconst4, OpIconst5:
			push(IntValue(int32(insn.Op) - 3))
		case OpLconst0, OpLconst1:
			push(LongValue(int64(insn.Op - OpLconst0)))
		case OpFconst0, OpFconst1, OpFconst2:
			push(FloatValue(float32(insn.Op - OpFconst0)))
		case OpDconst0, OpDconst1:
			push(DoubleValue(float64(insn.Op - OpDconst0)))
		case OpBipush, OpSipush:
			push(IntValue(insn.A))
		case OpLdc, OpLdcW, OpLdc2W:
			if insn.Ref < 0 || int(insn.Ref) >= len(m.Refs) {
				panic("catty: ldc without a constant reference")
			}
			ref := m.Refs[insn.Ref]
			switch ref.Kind {
			case RefInt:
				push(IntValue(ref.Int))
			case RefLong:
				push(LongValue(int64(ref.Bits)))
			case RefFloat:
				push(FloatBitsValue(uint32(ref.Bits)))
			case RefDouble:
				push(DoubleBitsValue(ref.Bits))
			case RefString:
				push(t.StringConstant(ref.Units))
			case RefClass:
				v, thrown := t.ClassConstant(m, ref.CP, ref.Class)
				if thrown != nil {
					if out, ok := abort(pc, thrown); !ok {
						return out
					}
					continue
				}
				push(v)
			default:
				panic("catty: unsupported ldc constant kind")
			}
		case OpIload, OpLload, OpFload, OpDload, OpAload:
			idx, _ := LocalIndex(insn.Op, insn.A)
			push(locals[idx])
		case OpIstore, OpLstore, OpFstore, OpDstore, OpAstore:
			idx, _ := LocalIndex(insn.Op, insn.A)
			storeLocal(idx, pop())
		case OpIinc:
			idx, _ := LocalIndex(OpIinc, insn.A)
			storeLocal(idx, IntValue(AsInt(locals[idx])+insn.B))
		case OpIaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'I')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpLaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'J')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpFaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'F')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpDaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'D')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpAaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'L')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpBaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'B')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpCaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'C')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpSaload:
			index := popInt()
			array := pop()
			v, thrown := t.ArrayLoad(array, index, 'S')
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpIastore, OpLastore, OpFastore, OpDastore, OpAastore, OpBastore, OpCastore, OpSastore:
			value := pop()
			index := popInt()
			array := pop()
			var kind byte
			switch insn.Op {
			case OpIastore:
				kind = 'I'
			case OpLastore:
				kind = 'J'
			case OpFastore:
				kind = 'F'
			case OpDastore:
				kind = 'D'
			case OpAastore:
				kind = 'L'
			case OpBastore:
				kind = 'B'
			case OpCastore:
				kind = 'C'
			case OpSastore:
				kind = 'S'
			}
			if thrown := t.ArrayStore(array, index, value, kind); thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
		case OpPop:
			_ = pop()
		case OpPop2:
			v := pop()
			if v.Width() == 1 {
				_ = pop()
			}
		case OpDup:
			push(stack[len(stack)-1])
		case OpDupX1:
			n := len(stack)
			v1, v2 := stack[n-1], stack[n-2]
			stack = append(stack[:n-2], v1, v2, v1)
		case OpDupX2:
			n := len(stack)
			v1 := stack[n-1]
			if stack[n-2].Width() == 2 {
				v2 := stack[n-2]
				stack = append(stack[:n-2], v1, v2, v1)
			} else {
				v2, v3 := stack[n-2], stack[n-3]
				stack = append(stack[:n-3], v1, v3, v2, v1)
			}
		case OpDup2:
			n := len(stack)
			if stack[n-1].Width() == 2 {
				push(stack[n-1])
				break
			}
			v1, v2 := stack[n-1], stack[n-2]
			stack = append(stack[:n-2], v2, v1, v2, v1)
		case OpDup2X1:
			n := len(stack)
			v1 := stack[n-1]
			if v1.Width() == 2 {
				v2 := stack[n-2]
				stack = append(stack[:n-2], v1, v2, v1)
			} else {
				v2, v3 := stack[n-2], stack[n-3]
				stack = append(stack[:n-3], v2, v1, v3, v2, v1)
			}
		case OpDup2X2:
			n := len(stack)
			v1 := stack[n-1]
			if v1.Width() == 2 {
				if stack[n-2].Width() == 2 {
					v2 := stack[n-2]
					stack = append(stack[:n-2], v1, v2, v1)
				} else {
					v2, v3 := stack[n-2], stack[n-3]
					stack = append(stack[:n-3], v1, v3, v2, v1)
				}
			} else {
				v2 := stack[n-2]
				if stack[n-3].Width() == 2 {
					v3 := stack[n-3]
					stack = append(stack[:n-3], v2, v1, v3, v2, v1)
				} else {
					v3, v4 := stack[n-3], stack[n-4]
					stack = append(stack[:n-4], v2, v1, v4, v3, v2, v1)
				}
			}
		case OpSwap:
			n := len(stack)
			stack[n-1], stack[n-2] = stack[n-2], stack[n-1]
		case OpIadd:
			b := popInt()
			a := popInt()
			push(IntValue(a + b))
		case OpLadd:
			b := popLong()
			a := popLong()
			push(LongValue(a + b))
		case OpFadd:
			b := popFloat()
			a := popFloat()
			push(FloatValue(a + b))
		case OpDadd:
			b := popDouble()
			a := popDouble()
			push(DoubleValue(a + b))
		case OpIsub:
			b := popInt()
			a := popInt()
			push(IntValue(a - b))
		case OpLsub:
			b := popLong()
			a := popLong()
			push(LongValue(a - b))
		case OpFsub:
			b := popFloat()
			a := popFloat()
			push(FloatValue(a - b))
		case OpDsub:
			b := popDouble()
			a := popDouble()
			push(DoubleValue(a - b))
		case OpImul:
			b := popInt()
			a := popInt()
			push(IntValue(a * b))
		case OpLmul:
			b := popLong()
			a := popLong()
			push(LongValue(a * b))
		case OpFmul:
			b := popFloat()
			a := popFloat()
			push(FloatValue(a * b))
		case OpDmul:
			b := popDouble()
			a := popDouble()
			push(DoubleValue(a * b))
		case OpIdiv, OpIrem:
			b := popInt()
			a := popInt()
			if b == 0 {
				thrown := r.NewThrowable(t, "java/lang/ArithmeticException")
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			if insn.Op == OpIdiv {
				push(IntValue(a / b))
			} else {
				push(IntValue(a % b))
			}
		case OpLdiv, OpLrem:
			b := popLong()
			a := popLong()
			if b == 0 {
				thrown := r.NewThrowable(t, "java/lang/ArithmeticException")
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			if insn.Op == OpLdiv {
				push(LongValue(a / b))
			} else {
				push(LongValue(a % b))
			}
		case OpFdiv:
			b := popFloat()
			a := popFloat()
			push(FloatValue(a / b))
		case OpDdiv:
			b := popDouble()
			a := popDouble()
			push(DoubleValue(a / b))
		case OpFrem:
			b := popFloat()
			a := popFloat()
			push(FloatValue(float32(math.Mod(float64(a), float64(b)))))
		case OpDrem:
			b := popDouble()
			a := popDouble()
			push(DoubleValue(math.Mod(a, b)))
		case OpIneg:
			push(IntValue(-popInt()))
		case OpLneg:
			push(LongValue(-popLong()))
		case OpFneg:
			push(FloatValue(-popFloat()))
		case OpDneg:
			push(DoubleValue(-popDouble()))
		case OpIshl:
			d := popInt()
			a := popInt()
			push(IntValue(a << IntShiftDistance(d)))
		case OpIshr:
			d := popInt()
			a := popInt()
			push(IntValue(a >> IntShiftDistance(d)))
		case OpIushr:
			d := popInt()
			a := popInt()
			push(IntValue(int32(uint32(a) >> IntShiftDistance(d))))
		case OpLshl:
			d := popInt()
			a := popLong()
			push(LongValue(a << LongShiftDistance(d)))
		case OpLshr:
			d := popInt()
			a := popLong()
			push(LongValue(a >> LongShiftDistance(d)))
		case OpLushr:
			d := popInt()
			a := popLong()
			push(LongValue(int64(uint64(a) >> LongShiftDistance(d))))
		case OpIand:
			b := popInt()
			a := popInt()
			push(IntValue(a & b))
		case OpIor:
			b := popInt()
			a := popInt()
			push(IntValue(a | b))
		case OpIxor:
			b := popInt()
			a := popInt()
			push(IntValue(a ^ b))
		case OpLand:
			b := popLong()
			a := popLong()
			push(LongValue(a & b))
		case OpLor:
			b := popLong()
			a := popLong()
			push(LongValue(a | b))
		case OpLxor:
			b := popLong()
			a := popLong()
			push(LongValue(a ^ b))
		case OpI2l:
			push(LongValue(int64(popInt())))
		case OpI2f:
			push(FloatValue(float32(popInt())))
		case OpI2d:
			push(DoubleValue(float64(popInt())))
		case OpL2i:
			push(IntValue(int32(popLong())))
		case OpL2f:
			push(FloatValue(float32(popLong())))
		case OpL2d:
			push(DoubleValue(float64(popLong())))
		case OpF2i:
			push(IntValue(F2I(popFloat())))
		case OpF2l:
			push(LongValue(F2L(popFloat())))
		case OpF2d:
			push(DoubleValue(float64(popFloat())))
		case OpD2i:
			push(IntValue(D2I(popDouble())))
		case OpD2l:
			push(LongValue(D2L(popDouble())))
		case OpD2f:
			push(FloatValue(float32(popDouble())))
		case OpI2b:
			push(IntValue(I2B(popInt())))
		case OpI2c:
			push(IntValue(I2C(popInt())))
		case OpI2s:
			push(IntValue(I2S(popInt())))
		case OpLcmp:
			b := popLong()
			a := popLong()
			push(IntValue(Lcmp(a, b)))
		case OpFcmpl:
			b := popFloat()
			a := popFloat()
			push(IntValue(Fcmp(a, b, -1)))
		case OpFcmpg:
			b := popFloat()
			a := popFloat()
			push(IntValue(Fcmp(a, b, 1)))
		case OpDcmpl:
			b := popDouble()
			a := popDouble()
			push(IntValue(Dcmp(a, b, -1)))
		case OpDcmpg:
			b := popDouble()
			a := popDouble()
			push(IntValue(Dcmp(a, b, 1)))
		case OpIfeq, OpIfne, OpIflt, OpIfge, OpIfgt, OpIfle:
			a := popInt()
			var take bool
			switch insn.Op {
			case OpIfeq:
				take = a == 0
			case OpIfne:
				take = a != 0
			case OpIflt:
				take = a < 0
			case OpIfge:
				take = a >= 0
			case OpIfgt:
				take = a > 0
			case OpIfle:
				take = a <= 0
			}
			if take {
				branch(insn.A)
				continue
			}
		case OpIfIcmpeq, OpIfIcmpne, OpIfIcmplt, OpIfIcmpge, OpIfIcmpgt, OpIfIcmple:
			b := popInt()
			a := popInt()
			var take bool
			switch insn.Op {
			case OpIfIcmpeq:
				take = a == b
			case OpIfIcmpne:
				take = a != b
			case OpIfIcmplt:
				take = a < b
			case OpIfIcmpge:
				take = a >= b
			case OpIfIcmpgt:
				take = a > b
			case OpIfIcmple:
				take = a <= b
			}
			if take {
				branch(insn.A)
				continue
			}
		case OpIfAcmpeq, OpIfAcmpne:
			b := popRef()
			a := popRef()
			take := (insn.Op == OpIfAcmpeq) == (a == b)
			if take {
				branch(insn.A)
				continue
			}
		case OpIfnull, OpIfnonnull:
			a := popRef()
			take := (insn.Op == OpIfnull) == (a == nil)
			if take {
				branch(insn.A)
				continue
			}
		case OpGoto, OpGotoW:
			branch(insn.A)
			continue
		case OpTableswitch, OpLookupswitch:
			key := popInt()
			target := insn.A
			for i, k := range insn.Keys {
				if k == key {
					target = insn.Targets[i]
					break
				}
			}
			branch(target)
			continue
		case OpNew:
			ref := m.Refs[insn.Ref]
			obj, thrown := t.NewObject(m, ref.CP, ref.Class)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpNewarray:
			count := popInt()
			obj, thrown := t.NewPrimitiveArray(m, 0, insn.A, count)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpAnewarray:
			ref := m.Refs[insn.Ref]
			count := popInt()
			obj, thrown := t.NewReferenceArray(m, ref.CP, ref.Class, count)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpMultianewarray:
			ref := m.Refs[insn.Ref]
			dims := int(insn.B)
			counts := make([]int32, dims)
			for i := dims - 1; i >= 0; i-- {
				counts[i] = popInt()
			}
			obj, thrown := t.MultiANewArray(m, ref.CP, ref.Class, dims, counts)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpArraylength:
			array := pop()
			n, thrown := t.ArrayLength(array)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(IntValue(n))
		case OpMonitorenter:
			ref := pop()
			if thrown := t.MonitorEnter(m, pc, ref); thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
		case OpMonitorexit:
			ref := pop()
			if thrown := t.MonitorExit(m, pc, ref); thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
		case OpAthrow:
			thrown := popRef()
			if thrown == nil {
				thrown = r.NewThrowable(t, "java/lang/NullPointerException")
			}
			if out, ok := abort(pc, thrown); !ok {
				return out
			}
			continue
		case OpCheckcast:
			ref := m.Refs[insn.Ref]
			v, thrown := t.CheckCast(m, ref.CP, ref.Class, pop())
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpInstanceof:
			ref := m.Refs[insn.Ref]
			n, thrown := t.InstanceOf(m, ref.CP, ref.Class, pop())
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(IntValue(n))
		case OpGetfield:
			ref := m.Refs[insn.Ref]
			obj := pop()
			v, thrown := t.GetField(m, ref.CP, obj, ref.Owner, ref.Name, ref.Descriptor)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpPutfield:
			ref := m.Refs[insn.Ref]
			v := pop()
			obj := pop()
			if thrown := t.PutField(m, ref.CP, obj, ref.Owner, ref.Name, ref.Descriptor, v); thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
		case OpGetstatic:
			ref := m.Refs[insn.Ref]
			v, thrown := t.GetStatic(m, ref.CP, ref.Owner, ref.Name, ref.Descriptor)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
			push(v)
		case OpPutstatic:
			ref := m.Refs[insn.Ref]
			v := pop()
			thrown := t.PutStatic(m, ref.CP, ref.Owner, ref.Name, ref.Descriptor, v)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out, ok := abort(pc, thrown); !ok {
					return out
				}
				continue
			}
		case OpInvokestatic, OpInvokevirtual, OpInvokespecial, OpInvokeinterface:
			ref := m.Refs[insn.Ref]
			mt, err := ParseMethodType(ref.Descriptor)
			if err != nil {
				panic("catty: bad method descriptor in verified body: " + err.Error())
			}
			n := mt.ArgCount
			callArgs := make([]Value, n)
			copy(callArgs, stack[len(stack)-n:])
			stack = stack[:len(stack)-n]
			var out Outcome
			switch insn.Op {
			case OpInvokestatic:
				out = t.InvokeStatic(m, pc, ref.CP, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokespecial:
				recv := pop()
				out = t.InvokeSpecial(m, pc, ref.CP, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokevirtual:
				recv := pop()
				out = t.InvokeVirtual(m, pc, ref.CP, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokeinterface:
				recv := pop()
				out = t.InvokeInterface(m, pc, ref.CP, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			}
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if out.Thrown != nil {
				if res, ok := abort(pc, out.Thrown); !ok {
					return res
				}
				continue
			}
			if mt.Return.Kind != 'V' {
				push(out.Value)
			}
		case OpReturn:
			return Normal()
		case OpIreturn, OpLreturn, OpFreturn, OpDreturn, OpAreturn:
			return OkValue(pop())
		default:
			panic(fmt.Sprintf("catty: interpreter reached unsupported opcode 0x%02x at %s.%s%s pc %d",
				insn.Op, m.Owner.Name, m.Name, m.Descriptor, pc))
		}
		ip++
	}
	panic(fmt.Sprintf("catty: fell off the end of %s.%s%s", m.Owner.Name, m.Name, m.Descriptor))
}
