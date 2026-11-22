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
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "InitSequence", "failures", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, "InitSequence", "failures", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "InitSequence", "signal", "LInitSequence$Signal;")
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

func cattyMethod2(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "InitSequence$Broken", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			out := t.InvokeStatic(m, 0, "InitSequence$Broken", "fail", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, "InitSequence$Broken", "value", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
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
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, stack[0], "InitSequence$Parent", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod5(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "InitSequence$Child", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(30)
			{
			out := t.InvokeStatic(m, 2, "InitSequence", "mark", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, "InitSequence$Child", "value", "I", stack[0])
			if thrown != nil {
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

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(10)
			{
			out := t.InvokeStatic(m, 2, "InitSequence", "mark", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, "InitSequence$Parent", "value", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/RuntimeException", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
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
			{
			v, thrown := t.GetStatic(m, "InitSequence$WithDefault", "MARKER", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod11(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(90)
			{
			out := t.InvokeStatic(m, 2, "InitSequence", "mark", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, "InitSequence$WithoutDefault", "MARKER", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
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
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
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
			stack[1] = locals[0]
			{
			out := t.InvokeVirtual(m, 4, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.ClassConstant("InitSequence$Child")
			stack[2] = t.ClassConstant("InitSequence$Child")
			if stack[1].Ref != stack[2].Ref {
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
			stack[1] = rt.IntValue(-1)
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 15, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeStatic(m, 21, "InitSequence$Child", "read", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeVirtual(m, 24, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeStatic(m, 30, "InitSequence$Broken", "read", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeVirtual(m, 33, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 33, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
		case 4:
			locals[1] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 40, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			{
			out := t.InvokeVirtual(m, 44, stack[1], "java/lang/ExceptionInInitializerError", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			v, thrown := t.GetStatic(m, "InitSequence", "signal", "LInitSequence$Signal;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			stack[1] = rt.IntValue(40)
			block = 7
			continue
		case 6:
			stack[1] = rt.IntValue(-40)
			block = 7
			continue
		case 7:
			{
			out := t.InvokeVirtual(m, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 63, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeStatic(m, 66, "InitSequence$Broken", "read", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 66, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeVirtual(m, 69, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 69, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 10
			continue
		case 9:
			locals[1] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 76, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(50)
			{
			out := t.InvokeVirtual(m, 81, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 10
			continue
		case 10:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 84, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			v, thrown := t.GetStatic(m, "InitSequence", "failures", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 87, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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
			out := t.InvokeVirtual(m, 90, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 39:
			block = 4
			case 75:
			block = 9
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

func cattyMethod15(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, "InitSequence$Signal")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, stack[1], "InitSequence$Signal", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutStatic(m, "InitSequence", "signal", "LInitSequence$Signal;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

