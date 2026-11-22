package main

import "catty/rt"

func cattyMethod0(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod1(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(11)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod2(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(12)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(-1)
			{
			out := t.InvokeStatic(m, 1, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(98)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod5(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(99)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "Base", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(22)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(123)
			{
			out := t.InvokeStatic(m, 2, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeVirtual(m, 1, stack[0], "Child", "value", "()Ljava/lang/Integer;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod10(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod11(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(41)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod12(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod13(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "Init", "state", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(30)
			{
			out := t.InvokeVirtual(m, 5, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(14)
			{
			thrown := t.PutStatic(m, "Init", "state", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod15(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod16(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(100)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod17(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.LongBitsValue(0x00000000000000c8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod18(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.DoubleBitsValue(0x4072c00000000000)
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) + rt.AsDouble(stack[1]))
			stack[0] = rt.IntValue(rt.D2I(rt.AsDouble(stack[0])))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod19(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.FloatBitsValue(0x43c80000)
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) + rt.AsFloat(stack[1]))
			stack[0] = rt.IntValue(rt.F2I(rt.AsFloat(stack[0])))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod20(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(500)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod21(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod22(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(31)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod23(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0x0000000000000020)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod24(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.FloatBitsValue(0x42060000)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod25(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.DoubleBitsValue(0x4041400000000000)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod26(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(1)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod27(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(67)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod28(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(35)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod29(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(36)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod30(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [0]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			return rt.Normal()
		}
	}
}

func cattyMethod31(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [30]rt.Value
	var stack [8]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, "Child")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, stack[1], "Child", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			stack[0] = t.ClassConstant("Child")
			stack[1] = t.StringConstant([]uint16{105, 110, 104, 101, 114, 105, 116, 101, 100})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 13, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 16, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 16, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[2] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 24, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 27, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			stack[2] = locals[1]
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 33, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 36, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 36, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 39, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 42, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 42, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 45, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 45, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.ClassConstant("Base")
			stack[1] = t.StringConstant([]uint16{111, 118, 101, 114, 114, 105, 100, 100, 101, 110})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 53, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 56, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 56, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[3] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(2)
			{
			out := t.InvokeVirtual(m, 64, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			stack[2] = locals[1]
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 73, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 76, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 76, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 79, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 82, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 82, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 85, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 85, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 88, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			{
			obj, thrown := t.NewObject(m, "Base")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 92, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			{
			out := t.InvokeSpecial(m, 96, stack[3], "Base", "<init>", "()V", stack[4:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 96, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 100, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 103, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 103, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 106, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 109, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 109, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 112, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 112, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.ClassConstant("Child")
			stack[1] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 120, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 123, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 123, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[4] = stack[0]
			stack[0] = t.ClassConstant("Base")
			stack[1] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 133, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 136, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 136, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 141, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(3)
			{
			out := t.InvokeVirtual(m, 145, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 145, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			stack[1] = locals[1]
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 152, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 155, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 155, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[6] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 160, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			n, thrown := t.InstanceOf("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 165, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[1] = rt.IntValue(1)
			block = 3
			continue
		case 2:
			stack[1] = rt.IntValue(0)
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 176, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 176, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 179, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 184, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 187, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 187, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 190, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 190, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[5]
			{
			obj, thrown := t.NewObject(m, "Base")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 195, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 199, stack[2], "Base", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 199, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 203, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 206, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 206, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 211, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			n, thrown := t.InstanceOf("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 216, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 5
			continue
			}
			block = 4
			continue
		case 4:
			stack[1] = rt.IntValue(1)
			block = 6
			continue
		case 5:
			stack[1] = rt.IntValue(0)
			block = 6
			continue
		case 6:
			{
			out := t.InvokeVirtual(m, 227, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 227, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 230, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 235, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 238, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 238, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 241, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 241, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			{
			obj, thrown := t.NewObject(m, "Base")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 246, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 250, stack[2], "Base", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 250, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 254, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 257, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 257, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 261, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-4)
			{
			out := t.InvokeVirtual(m, 266, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 266, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 8
			continue
		case 7:
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 274, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(4)
			{
			out := t.InvokeVirtual(m, 278, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 278, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 8
			continue
		case 8:
			stack[0] = t.ClassConstant("Init")
			stack[1] = t.StringConstant([]uint16{115, 116, 97, 116, 115})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 286, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 289, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 289, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 294, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(5)
			{
			out := t.InvokeVirtual(m, 298, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 298, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 301, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[8]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 308, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 311, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 311, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 314, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 317, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 317, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 320, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 320, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 323, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[8]
			{
			obj, thrown := t.NewObject(m, "Init")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 328, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			{
			out := t.InvokeSpecial(m, 332, stack[3], "Init", "<init>", "()V", stack[4:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 332, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 336, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 339, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 339, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 342, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 345, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 345, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 348, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 348, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 73, 110, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 356, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Integer", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 361, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 364, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 365, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 365, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[9] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 76, 111, 110, 103})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 375, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Long", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 380, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 383, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 384, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 384, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[10] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 68, 111, 117, 98, 108, 101})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 394, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Double", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 399, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 402, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 403, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 403, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[11] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 70, 108, 111, 97, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 413, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Float", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 418, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 421, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 422, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 422, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[12] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 67, 104, 97, 114})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 432, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Character", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 437, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 440, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 441, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 441, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[13] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{116, 97, 107, 101, 66, 111, 111, 108, 101, 97, 110})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 451, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Boolean", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 456, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 459, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 460, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 460, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[14] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 465, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(6)
			{
			out := t.InvokeVirtual(m, 470, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 470, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 473, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/Integer", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 476, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/Integer", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 479, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			if stack[1].Ref != stack[2].Ref {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			stack[1] = rt.IntValue(1)
			block = 11
			continue
		case 10:
			stack[1] = rt.IntValue(0)
			block = 11
			continue
		case 11:
			{
			out := t.InvokeVirtual(m, 490, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 490, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 493, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[9]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 500, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(5)
			{
			out := t.InvokeStatic(m, 506, "java/lang/Byte", "valueOf", "(B)Ljava/lang/Byte;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 506, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 509, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 510, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 510, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 513, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 516, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 516, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 519, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 519, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 522, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[9]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 529, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(6)
			{
			out := t.InvokeStatic(m, 536, "java/lang/Short", "valueOf", "(S)Ljava/lang/Short;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 536, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 539, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 540, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 540, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 543, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 546, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 546, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 549, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 549, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 552, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[9]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 559, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(65)
			{
			out := t.InvokeStatic(m, 566, "java/lang/Character", "valueOf", "(C)Ljava/lang/Character;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 566, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 569, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 570, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 570, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 573, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 576, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 576, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 579, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 579, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 582, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[10]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 589, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(7)
			{
			out := t.InvokeStatic(m, 596, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 596, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 599, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 600, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 600, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Long", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 603, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 606, stack[1], "java/lang/Long", "longValue", "()J", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 606, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 610, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 610, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 613, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[11]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 620, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.LongBitsValue(0x0000000000000008)
			{
			out := t.InvokeStatic(m, 628, "java/lang/Long", "valueOf", "(J)Ljava/lang/Long;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 628, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 631, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 632, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 632, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 635, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 638, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 638, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 641, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 641, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 644, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[11]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 651, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.FloatBitsValue(0x3fc00000)
			{
			out := t.InvokeStatic(m, 658, "java/lang/Float", "valueOf", "(F)Ljava/lang/Float;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 658, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 661, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 662, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 662, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 665, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 668, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 668, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 671, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 671, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 674, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[12]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 681, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(9)
			{
			out := t.InvokeStatic(m, 688, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 688, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 691, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 692, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 692, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 695, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 698, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 698, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 701, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 701, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 704, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[11]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 711, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.DoubleBitsValue(0x4024000000000000)
			{
			out := t.InvokeStatic(m, 719, "java/lang/Double", "valueOf", "(D)Ljava/lang/Double;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 719, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 722, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 723, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 723, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 726, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 729, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 729, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 732, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 732, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 735, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[12]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 742, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.FloatValue(1)
			{
			out := t.InvokeStatic(m, 748, "java/lang/Float", "valueOf", "(F)Ljava/lang/Float;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 748, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 751, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 752, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 752, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 755, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 758, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 758, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 761, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 761, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 764, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[13]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 771, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(66)
			{
			out := t.InvokeStatic(m, 778, "java/lang/Character", "valueOf", "(C)Ljava/lang/Character;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 778, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 781, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 782, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 782, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 785, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 788, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 788, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 791, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 791, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 794, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[14]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 801, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 807, "java/lang/Boolean", "valueOf", "(Z)Ljava/lang/Boolean;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 807, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 810, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 811, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 811, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Boolean", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 814, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 817, stack[1], "java/lang/Boolean", "booleanValue", "()Z", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 817, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 13
			continue
			}
			block = 12
			continue
		case 12:
			stack[1] = rt.IntValue(1)
			block = 14
			continue
		case 13:
			stack[1] = rt.IntValue(0)
			block = 14
			continue
		case 14:
			{
			out := t.InvokeVirtual(m, 828, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 828, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 831, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(7)
			{
			out := t.InvokeVirtual(m, 836, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 836, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[9]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 843, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = rt.LongValue(1)
			{
			out := t.InvokeStatic(m, 849, "java/lang/Long", "valueOf", "(J)Ljava/lang/Long;", stack[5:6])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 849, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 852, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 853, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 853, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 857, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-8)
			{
			out := t.InvokeVirtual(m, 862, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 862, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 16
			continue
		case 15:
			locals[15] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 870, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(8)
			{
			out := t.InvokeVirtual(m, 875, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 875, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 16
			continue
		case 16:
			stack[0] = locals[10]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 882, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = rt.DoubleValue(1)
			{
			out := t.InvokeStatic(m, 888, "java/lang/Double", "valueOf", "(D)Ljava/lang/Double;", stack[5:6])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 888, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 891, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 892, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 892, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 896, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-9)
			{
			out := t.InvokeVirtual(m, 901, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 901, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 18
			continue
		case 17:
			locals[15] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 909, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(9)
			{
			out := t.InvokeVirtual(m, 914, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 914, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 18
			continue
		case 18:
			stack[0] = locals[9]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 921, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = rt.RefValue(nil)
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 927, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 928, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 928, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 932, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-10)
			{
			out := t.InvokeVirtual(m, 937, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 937, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 20
			continue
		case 19:
			locals[15] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 945, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(10)
			{
			out := t.InvokeVirtual(m, 950, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 950, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 20
			continue
		case 20:
			stack[0] = locals[14]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 957, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 963, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[5:6])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 963, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 966, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 967, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 967, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 971, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-11)
			{
			out := t.InvokeVirtual(m, 976, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 976, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 22
			continue
		case 21:
			locals[15] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 984, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(11)
			{
			out := t.InvokeVirtual(m, 989, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 989, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 22
			continue
		case 22:
			stack[0] = locals[13]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 996, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = rt.IntValue(65)
			{
			out := t.InvokeStatic(m, 1003, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[5:6])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1003, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1006, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 1007, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1007, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1011, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-12)
			{
			out := t.InvokeVirtual(m, 1016, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1016, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 24
			continue
		case 23:
			locals[15] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1024, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(12)
			{
			out := t.InvokeVirtual(m, 1029, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1029, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 24
			continue
		case 24:
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 73, 110, 116})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1037, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1040, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1040, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[15] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 76, 111, 110, 103})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1050, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1053, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1053, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[16] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 70, 108, 111, 97, 116})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1063, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1066, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1066, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[17] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 68, 111, 117, 98, 108, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1076, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1079, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1079, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[18] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 66, 111, 111, 108, 101, 97, 110})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1089, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1092, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1092, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[19] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 67, 104, 97, 114})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1102, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1105, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1105, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[20] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 83, 104, 111, 114, 116})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1115, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1118, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1118, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[21] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 66, 121, 116, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1128, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1131, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1131, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[22] = stack[0]
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{98, 111, 120, 86, 111, 105, 100})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1141, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1144, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1144, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[23] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1149, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(13)
			{
			out := t.InvokeVirtual(m, 1154, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1154, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[15]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1161, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1164, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1164, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1169, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1174, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 26
			continue
			}
			block = 25
			continue
		case 25:
			stack[1] = rt.IntValue(1)
			block = 27
			continue
		case 26:
			stack[1] = rt.IntValue(0)
			block = 27
			continue
		case 27:
			{
			out := t.InvokeVirtual(m, 1185, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1185, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1188, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1193, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1196, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1196, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 1199, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1199, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[16]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1206, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1209, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1209, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1214, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Long", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1219, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 29
			continue
			}
			block = 28
			continue
		case 28:
			stack[1] = rt.IntValue(1)
			block = 30
			continue
		case 29:
			stack[1] = rt.IntValue(0)
			block = 30
			continue
		case 30:
			{
			out := t.InvokeVirtual(m, 1230, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1230, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1233, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Long", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1238, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1241, stack[1], "java/lang/Long", "longValue", "()J", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1241, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 1245, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1245, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[17]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1252, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1255, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1255, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1260, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Float", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1265, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 32
			continue
			}
			block = 31
			continue
		case 31:
			stack[1] = rt.IntValue(1)
			block = 33
			continue
		case 32:
			stack[1] = rt.IntValue(0)
			block = 33
			continue
		case 33:
			{
			out := t.InvokeVirtual(m, 1276, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1276, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1279, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Float", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1284, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1287, stack[1], "java/lang/Float", "floatValue", "()F", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1287, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[1] = rt.IntValue(rt.F2I(rt.AsFloat(stack[1])))
			{
			out := t.InvokeVirtual(m, 1291, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1291, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[18]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1298, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1301, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1301, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1306, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Double", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1311, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 35
			continue
			}
			block = 34
			continue
		case 34:
			stack[1] = rt.IntValue(1)
			block = 36
			continue
		case 35:
			stack[1] = rt.IntValue(0)
			block = 36
			continue
		case 36:
			{
			out := t.InvokeVirtual(m, 1322, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1322, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1325, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Double", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1330, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1333, stack[1], "java/lang/Double", "doubleValue", "()D", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1333, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[1] = rt.IntValue(rt.D2I(rt.AsDouble(stack[1])))
			{
			out := t.InvokeVirtual(m, 1337, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1337, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[19]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1344, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1347, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1347, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1352, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Boolean", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1357, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 38
			continue
			}
			block = 37
			continue
		case 37:
			stack[1] = rt.IntValue(1)
			block = 39
			continue
		case 38:
			stack[1] = rt.IntValue(0)
			block = 39
			continue
		case 39:
			{
			out := t.InvokeVirtual(m, 1368, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1368, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1371, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Boolean", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1376, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1379, stack[1], "java/lang/Boolean", "booleanValue", "()Z", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1379, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 41
			continue
			}
			block = 40
			continue
		case 40:
			stack[1] = rt.IntValue(1)
			block = 42
			continue
		case 41:
			stack[1] = rt.IntValue(0)
			block = 42
			continue
		case 42:
			{
			out := t.InvokeVirtual(m, 1390, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1390, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[20]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1397, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1400, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1400, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1405, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Character", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1410, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 44
			continue
			}
			block = 43
			continue
		case 43:
			stack[1] = rt.IntValue(1)
			block = 45
			continue
		case 44:
			stack[1] = rt.IntValue(0)
			block = 45
			continue
		case 45:
			{
			out := t.InvokeVirtual(m, 1421, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1421, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1424, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Character", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1429, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1432, stack[1], "java/lang/Character", "charValue", "()C", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1432, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 1435, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1435, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[21]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1442, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1445, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1445, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1450, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Short", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1455, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 47
			continue
			}
			block = 46
			continue
		case 46:
			stack[1] = rt.IntValue(1)
			block = 48
			continue
		case 47:
			stack[1] = rt.IntValue(0)
			block = 48
			continue
		case 48:
			{
			out := t.InvokeVirtual(m, 1466, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1466, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1469, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Short", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1474, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1477, stack[1], "java/lang/Short", "shortValue", "()S", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1477, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 1480, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1480, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[22]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1487, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1490, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1490, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1495, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			n, thrown := t.InstanceOf("java/lang/Byte", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1500, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			if rt.AsInt(stack[1]) == 0 {
			block = 50
			continue
			}
			block = 49
			continue
		case 49:
			stack[1] = rt.IntValue(1)
			block = 51
			continue
		case 50:
			stack[1] = rt.IntValue(0)
			block = 51
			continue
		case 51:
			{
			out := t.InvokeVirtual(m, 1511, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1511, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1514, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			{
			v, thrown := t.CheckCast("java/lang/Byte", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1519, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1522, stack[1], "java/lang/Byte", "byteValue", "()B", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1522, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 1525, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1525, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[23]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1532, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1535, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1535, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[24] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1540, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[24]
			if stack[1].Ref != nil {
			block = 53
			continue
			}
			block = 52
			continue
		case 52:
			stack[1] = rt.IntValue(1)
			block = 54
			continue
		case 53:
			stack[1] = rt.IntValue(0)
			block = 54
			continue
		case 54:
			{
			out := t.InvokeVirtual(m, 1553, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1553, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1556, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(14)
			{
			out := t.InvokeVirtual(m, 1561, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1561, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1567, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1570, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1570, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1574, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-15)
			{
			out := t.InvokeVirtual(m, 1579, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1579, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 56
			continue
		case 55:
			locals[25] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1587, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(15)
			{
			out := t.InvokeVirtual(m, 1592, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1592, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 56
			continue
		case 56:
			stack[0] = locals[2]
			{
			obj, thrown := t.NewObject(m, "java/lang/Object")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1596, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 1600, stack[2], "java/lang/Object", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1600, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1604, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1607, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1607, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1611, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-16)
			{
			out := t.InvokeVirtual(m, 1616, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1616, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 58
			continue
		case 57:
			locals[25] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1624, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(16)
			{
			out := t.InvokeVirtual(m, 1629, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1629, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 58
			continue
		case 58:
			stack[0] = t.ClassConstant("Face")
			stack[1] = t.StringConstant([]uint16{99, 97, 108, 108})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1637, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1640, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1640, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[25] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1645, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[25]
			{
			obj, thrown := t.NewObject(m, "FaceImpl")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1650, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			{
			out := t.InvokeSpecial(m, 1654, stack[3], "FaceImpl", "<init>", "()V", stack[4:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1654, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[3])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1658, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			out := t.InvokeVirtual(m, 1661, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1661, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1664, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 1667, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1667, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 1670, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1670, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1673, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(17)
			{
			out := t.InvokeVirtual(m, 1678, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1678, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.ClassConstant("Child")
			stack[1] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1686, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Integer", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1691, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1694, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 1695, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1695, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1699, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-18)
			{
			out := t.InvokeVirtual(m, 1704, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1704, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 60
			continue
		case 59:
			locals[26] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1712, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(18)
			{
			out := t.InvokeVirtual(m, 1717, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1717, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 60
			continue
		case 60:
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{109, 105, 115, 115, 105, 110, 103})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1725, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1728, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1728, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1732, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-19)
			{
			out := t.InvokeVirtual(m, 1737, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1737, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 62
			continue
		case 61:
			locals[26] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1745, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(19)
			{
			out := t.InvokeVirtual(m, 1750, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1750, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 62
			continue
		case 62:
			stack[0] = t.ClassConstant("Base")
			stack[1] = t.StringConstant([]uint16{112, 97, 99, 107, 97, 103, 101, 80, 114, 105, 118, 97, 116, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1758, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1761, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1761, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1765, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-20)
			{
			out := t.InvokeVirtual(m, 1770, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1770, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 64
			continue
		case 63:
			locals[26] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1778, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(20)
			{
			out := t.InvokeVirtual(m, 1783, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1783, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 64
			continue
		case 64:
			stack[0] = t.ClassConstant("Thrower")
			stack[1] = t.StringConstant([]uint16{101, 120, 112, 108, 111, 100, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1791, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1794, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1794, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[26] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1799, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(21)
			{
			out := t.InvokeVirtual(m, 1804, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1804, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[26]
			{
			obj, thrown := t.NewObject(m, "Thrower")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1809, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 1813, stack[2], "Thrower", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1813, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1817, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1820, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1820, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1824, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-22)
			{
			out := t.InvokeVirtual(m, 1829, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1829, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 75
			continue
		case 65:
			locals[27] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1837, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[27]
			{
			v, thrown := t.GetStatic(m, "Thrower", "SAVED", "Ljava/lang/RuntimeException;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1842, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			if stack[1].Ref == stack[2].Ref {
			block = 67
			continue
			}
			block = 66
			continue
		case 66:
			stack[1] = rt.IntValue(1)
			block = 68
			continue
		case 67:
			stack[1] = rt.IntValue(0)
			block = 68
			continue
		case 68:
			{
			out := t.InvokeVirtual(m, 1853, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1853, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1856, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[27]
			{
			out := t.InvokeVirtual(m, 1861, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1861, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "Thrower", "SAVED", "Ljava/lang/RuntimeException;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1864, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			if stack[1].Ref != stack[2].Ref {
			block = 70
			continue
			}
			block = 69
			continue
		case 69:
			stack[1] = rt.IntValue(1)
			block = 71
			continue
		case 70:
			stack[1] = rt.IntValue(0)
			block = 71
			continue
		case 71:
			{
			out := t.InvokeVirtual(m, 1875, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1875, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1878, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[27]
			{
			out := t.InvokeVirtual(m, 1883, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1883, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = locals[27]
			if stack[1].Ref == stack[2].Ref {
			block = 73
			continue
			}
			block = 72
			continue
		case 72:
			stack[1] = rt.IntValue(1)
			block = 74
			continue
		case 73:
			stack[1] = rt.IntValue(0)
			block = 74
			continue
		case 74:
			{
			out := t.InvokeVirtual(m, 1896, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1896, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 75
			continue
		case 75:
			{
			out := t.InvokeStatic(m, 1899, "hidden/Factory", "make", "()Ljava/lang/Object;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1899, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[27] = stack[0]
			stack[0] = locals[27]
			{
			out := t.InvokeVirtual(m, 1906, stack[0], "java/lang/Object", "getClass", "()Ljava/lang/Class;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1906, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Class;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1912, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1915, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1915, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[28] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1920, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(22)
			{
			out := t.InvokeVirtual(m, 1925, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1925, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[28]
			stack[1] = locals[27]
			stack[2] = rt.RefValue(nil)
			{
			v, thrown := t.CheckCast("[Ljava/lang/Object;", stack[2])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1933, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			{
			out := t.InvokeVirtual(m, 1936, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1936, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1940, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(-23)
			{
			out := t.InvokeVirtual(m, 1945, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1945, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 77
			continue
		case 76:
			locals[29] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1953, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(23)
			{
			out := t.InvokeVirtual(m, 1958, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1958, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 77
			continue
		case 77:
			stack[0] = t.ClassConstant("Reflect")
			stack[1] = t.StringConstant([]uint16{108, 101, 97, 102})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1966, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			{
			v, thrown := t.GetStatic(m, "java/lang/Integer", "TYPE", "Ljava/lang/Class;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1971, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[5] = v
			}
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1974, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 1975, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1975, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[29] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1980, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(24)
			{
			out := t.InvokeVirtual(m, 1985, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 1985, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1988, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[29]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1995, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 2001, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2001, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2004, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 2005, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2005, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2008, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 2011, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2011, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 2014, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2014, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2017, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[29]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2024, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			stack[4] = stack[3]
			stack[5] = rt.IntValue(0)
			stack[6] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 2030, "java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;", stack[6:7])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2030, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[6] = out.Value
			}
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2033, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 2034, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2034, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			v, thrown := t.CheckCast("java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2037, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 2040, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2040, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 2043, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2043, thrown); ok {
			stack = [8]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 272:
			block = 7
			case 868:
			block = 15
			case 907:
			block = 17
			case 943:
			block = 19
			case 982:
			block = 21
			case 1022:
			block = 23
			case 1585:
			block = 55
			case 1622:
			block = 57
			case 1710:
			block = 59
			case 1743:
			block = 61
			case 1776:
			block = 63
			case 1835:
			block = 65
			case 1951:
			block = 76
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod32(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod33(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "Thrower", "SAVED", "Ljava/lang/RuntimeException;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			return rt.ThrowOutcome(thrown)
			}
		}
	}
}

func cattyMethod34(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, "java/lang/RuntimeException")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, stack[1], "java/lang/RuntimeException", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutStatic(m, "Thrower", "SAVED", "Ljava/lang/RuntimeException;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod35(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod36(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, "hidden/Hidden")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, stack[1], "hidden/Hidden", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod37(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod38(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(42)
			return rt.OkValue(stack[0])
		}
	}
}

