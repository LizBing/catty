package rt

import (
	"fmt"
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
		base := 1
		for i, a := range args {
			if base+i < len(locals) {
				locals[base+i] = a
			}
		}
	} else {
		for i, a := range args {
			if i < len(locals) {
				locals[i] = a
			}
		}
	}
	stack := make([]Value, m.MaxStack)
	sp := 0
	ip := 0

	push := func(v Value) {
		stack[sp] = v
		sp++
	}
	pop := func() Value {
		sp--
		v := stack[sp]
		stack[sp] = Value{}
		return v
	}
	abort := func(pc int32, thrown *Object) Outcome {
		if hpc, ok := t.HandleThrow(m, pc, thrown); ok {
			for i := range stack {
				stack[i] = Value{}
			}
			stack[0] = RefValue(thrown)
			sp = 1
			idx, ok := m.pcIndex[hpc]
			if !ok {
				panic(fmt.Sprintf("catty: handler offset %d in %s.%s%s is not an instruction", hpc, m.Owner.Name, m.Name, m.Descriptor))
			}
			ip = idx
			return Outcome{}
		}
		return ThrowOutcome(thrown)
	}
	branch := func(target int32) {
		idx, ok := m.pcIndex[target]
		if !ok {
			panic(fmt.Sprintf("catty: branch target %d in %s.%s%s is not an instruction", target, m.Owner.Name, m.Name, m.Descriptor))
		}
		ip = idx
	}

	for ip < len(m.Insns) {
		insn := m.Insns[ip]
		pc := insn.PC
		t.SetPC(pc)
		switch insn.Op {
		case OpAconstNull:
			push(RefValue(nil))
		case OpIconstM1, OpIconst0, OpIconst1, OpIconst2, OpIconst3, OpIconst4, OpIconst5:
			push(IntValue(int32(insn.Op) - 3))
		case OpBipush, OpSipush:
			push(IntValue(insn.A))
		case OpLdc, OpLdcW:
			if insn.Ref < 0 || int(insn.Ref) >= len(m.Refs) {
				panic("catty: ldc without a constant reference")
			}
			ref := m.Refs[insn.Ref]
			switch ref.Kind {
			case RefInt:
				push(IntValue(ref.Int))
			case RefString:
				push(t.StringConstant(ref.Units))
			case RefClass:
				push(t.ClassConstant(ref.Class))
			default:
				panic("catty: unsupported ldc constant kind")
			}
		case OpIload, OpAload, 0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d:
			idx, _ := localIndex(insn.Op, insn.A)
			push(locals[idx])
		case OpIstore, OpAstore, 0x3b, 0x3c, 0x3d, 0x3e, 0x4b, 0x4c, 0x4d, 0x4e:
			idx, _ := localIndex(insn.Op, insn.A)
			locals[idx] = pop()
		case OpDup:
			push(stack[sp-1])
		case OpPop:
			_ = pop()
		case OpNew:
			ref := m.Refs[insn.Ref]
			obj, thrown := t.NewObject(ref.Class)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpGetfield:
			ref := m.Refs[insn.Ref]
			obj := pop()
			v, thrown := t.GetField(obj, ref.Owner, ref.Name, ref.Descriptor)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(v)
		case OpPutfield:
			ref := m.Refs[insn.Ref]
			v := pop()
			obj := pop()
			if thrown := t.PutField(obj, ref.Owner, ref.Name, ref.Descriptor, v); thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
		case OpGetstatic:
			ref := m.Refs[insn.Ref]
			v, thrown := t.GetStatic(ref.Owner, ref.Name, ref.Descriptor)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(v)
		case OpPutstatic:
			ref := m.Refs[insn.Ref]
			v := pop()
			thrown := t.PutStatic(ref.Owner, ref.Name, ref.Descriptor, v)
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
		case OpAnewarray:
			ref := m.Refs[insn.Ref]
			count := AsInt(pop())
			obj, thrown := t.NewArray(count, ref.Class)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(RefValue(obj))
		case OpAaload:
			index := AsInt(pop())
			array := pop()
			v, thrown := t.ArrayLoad(array, index)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(v)
		case OpAastore:
			v := pop()
			index := AsInt(pop())
			array := pop()
			if thrown := t.ArrayStore(array, index, v); thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
		case OpArraylength:
			array := pop()
			n, thrown := t.ArrayLength(array)
			if thrown != nil {
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(IntValue(n))
		case OpIadd:
			b := AsInt(pop())
			a := AsInt(pop())
			push(IntValue(a + b))
		case OpIrem:
			b := AsInt(pop())
			a := AsInt(pop())
			if b == 0 {
				thrown := r.NewThrowable(t, "java/lang/ArithmeticException")
				if out := abort(pc, thrown); out.Thrown != nil {
					return out
				}
				continue
			}
			push(IntValue(a % b))
		case OpIfeq, OpIfne:
			a := AsInt(pop())
			take := (insn.Op == OpIfeq && a == 0) || (insn.Op == OpIfne && a != 0)
			if take {
				branch(insn.A)
				continue
			}
		case OpIfIcmpeq, OpIfIcmpne:
			b := AsInt(pop())
			a := AsInt(pop())
			take := (insn.Op == OpIfIcmpeq && a == b) || (insn.Op == OpIfIcmpne && a != b)
			if take {
				branch(insn.A)
				continue
			}
		case OpIfAcmpeq, OpIfAcmpne:
			b := pop()
			a := pop()
			same := a.Kind == KindRef && b.Kind == KindRef && a.Ref == b.Ref
			take := (insn.Op == OpIfAcmpeq && same) || (insn.Op == OpIfAcmpne && !same)
			if take {
				branch(insn.A)
				continue
			}
		case OpGoto:
			branch(insn.A)
			continue
		case OpInvokestatic, OpInvokevirtual, OpInvokespecial, OpInvokeinterface:
			ref := m.Refs[insn.Ref]
			mt, err := ParseMethodType(ref.Descriptor)
			if err != nil {
				panic("catty: bad method descriptor in verified body: " + err.Error())
			}
			n := mt.ArgCount
			callArgs := make([]Value, n)
			copy(callArgs, stack[sp-n:sp])
			for i := sp - n; i < sp; i++ {
				stack[i] = Value{}
			}
			sp -= n
			var out Outcome
			switch insn.Op {
			case OpInvokestatic:
				out = t.InvokeStatic(m, pc, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokespecial:
				recv := pop()
				out = t.InvokeSpecial(m, pc, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokevirtual:
				recv := pop()
				out = t.InvokeVirtual(m, pc, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			case OpInvokeinterface:
				recv := pop()
				out = t.InvokeInterface(m, pc, recv, ref.Owner, ref.Name, ref.Descriptor, callArgs)
			}
			runtime.KeepAlive(stack)
			runtime.KeepAlive(locals)
			if out.Thrown != nil {
				if res := abort(pc, out.Thrown); res.Thrown != nil {
					return res
				}
				continue
			}
			if mt.Return.Kind != 'V' {
				push(out.Value)
			}
		case OpReturn:
			return Normal()
		case OpIreturn, OpAreturn:
			return OkValue(pop())
		case OpAthrow:
			v := pop()
			var thrown *Object
			if v.Kind != KindRef || v.Ref == nil {
				thrown = r.NewThrowable(t, "java/lang/NullPointerException")
			} else {
				thrown = v.Ref
			}
			if out := abort(pc, thrown); out.Thrown != nil {
				return out
			}
			continue
		default:
			panic(fmt.Sprintf("catty: interpreter reached unsupported opcode 0x%02x at %s.%s%s pc %d",
				insn.Op, m.Owner.Name, m.Name, m.Descriptor, pc))
		}
		ip++
	}
	panic(fmt.Sprintf("catty: fell off the end of %s.%s%s", m.Owner.Name, m.Name, m.Descriptor))
}
