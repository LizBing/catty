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
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
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
	var locals [2]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			{
			obj, thrown := t.NewObject(m, 7, "java/lang/RuntimeException")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = locals[1]
			{
			out := t.InvokeSpecial(m, 9, 9, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			return rt.ThrowOutcome(thrown)
			}
		case 2:
			return rt.Normal()
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
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0x00000000000003e8)
			{
			out := t.InvokeStatic(m, 3, 9, "java/lang/Thread", "sleep", "(J)V", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 9:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 0
			continue
		case 1:
			locals[1] = stack[0]
			block = 0
			continue
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 97})
			{
			out := t.InvokeSpecial(m, 3, 3, stack[0], "java/lang/Thread", "<init>", "(Ljava/lang/String;)V", stack[1:2])
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
	var locals [4]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 9, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			out := t.InvokeStatic(m, 4, 15, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			if stack[1].Ref != stack[2].Ref {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[1] = t.StringConstant([]uint16{105, 100, 101, 110, 116, 105, 116, 121, 45, 97, 61, 49})
			block = 3
			continue
		case 2:
			stack[1] = t.StringConstant([]uint16{105, 100, 101, 110, 116, 105, 116, 121, 45, 97, 61, 48})
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 17, 23, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 17, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = stack[0]
			locals[1] = stack[1]
			{
			thrown := t.MonitorEnter(m, 25, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 25, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 34, "Hooks", "arrived", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 26, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, 34, "Hooks", "arrived", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 31, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			out := t.InvokeVirtual(m, 37, 38, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 37, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 4:
			{
			v, thrown := t.GetStatic(m, 34, "Hooks", "arrived", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 40, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 7
			continue
			}
			block = 5
			continue
		case 5:
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			out := t.InvokeVirtual(m, 50, 44, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 6:
			locals[2] = stack[0]
			block = 4
			continue
		case 7:
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 61, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 9
			continue
		case 8:
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 67, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 69, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 9:
			{
			obj, thrown := t.NewObject(m, 4, "java/lang/Thread")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 70, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			obj, thrown := t.NewObject(m, 49, "Hooks$Blocked")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 74, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			out := t.InvokeSpecial(m, 78, 51, stack[3], "Hooks$Blocked", "<init>", "()V", stack[4:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 78, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 119, 111, 114, 107, 101, 114})
			{
			out := t.InvokeSpecial(m, 83, 55, stack[1], "java/lang/Thread", "<init>", "(Ljava/lang/Runnable;Ljava/lang/String;)V", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 83, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 86, 58, stack[0], "java/lang/Thread", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 86, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 9, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 89, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 97})
			{
			out := t.InvokeVirtual(m, 94, 23, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 94, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 98})
			{
			out := t.InvokeSpecial(m, 3, 3, stack[0], "java/lang/Thread", "<init>", "(Ljava/lang/String;)V", stack[1:2])
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
	var locals [4]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 9, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			out := t.InvokeStatic(m, 4, 15, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			if stack[1].Ref != stack[2].Ref {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[1] = t.StringConstant([]uint16{105, 100, 101, 110, 116, 105, 116, 121, 45, 98, 61, 49})
			block = 3
			continue
		case 2:
			stack[1] = t.StringConstant([]uint16{105, 100, 101, 110, 116, 105, 116, 121, 45, 98, 61, 48})
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 17, 23, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 17, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = stack[0]
			locals[1] = stack[1]
			{
			thrown := t.MonitorEnter(m, 25, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 25, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 34, "Hooks", "arrived", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 26, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, 34, "Hooks", "arrived", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 31, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			out := t.InvokeVirtual(m, 37, 38, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 37, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 4:
			{
			v, thrown := t.GetStatic(m, 34, "Hooks", "arrived", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 40, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 7
			continue
			}
			block = 5
			continue
		case 5:
			{
			v, thrown := t.GetStatic(m, 28, "Hooks", "BARRIER", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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
			out := t.InvokeVirtual(m, 50, 44, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 6:
			locals[2] = stack[0]
			block = 4
			continue
		case 7:
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 61, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 9
			continue
		case 8:
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 67, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 69, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 9:
			{
			v, thrown := t.GetStatic(m, 9, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 70, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 98})
			{
			out := t.InvokeVirtual(m, 75, 23, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 56:
			block = 6
			case 65:
			block = 8
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

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = t.StringConstant([]uint16{98, 97, 100, 45, 104, 111, 111, 107})
			{
			out := t.InvokeSpecial(m, 3, 3, stack[0], "java/lang/Thread", "<init>", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 9, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{116, 104, 114, 111, 119, 105, 110, 103})
			{
			out := t.InvokeVirtual(m, 5, 17, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 22, "java/lang/RuntimeException")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = t.StringConstant([]uint16{104, 111, 111, 107, 45, 98, 111, 111, 109})
			{
			out := t.InvokeSpecial(m, 14, 26, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
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
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
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
	var locals [3]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			out := t.InvokeStatic(m, 0, 7, "java/lang/Runtime", "getRuntime", "()Ljava/lang/Runtime;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[1] = stack[0]
			stack[0] = locals[1]
			{
			obj, thrown := t.NewObject(m, 13, "Hooks$HookA")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 9, 15, stack[2], "Hooks$HookA", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 12, 16, stack[0], "java/lang/Runtime", "addShutdownHook", "(Ljava/lang/Thread;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 20, "Hooks$HookB")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 19, 22, stack[1], "Hooks$HookB", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			stack[0] = locals[2]
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 25, 23, stack[0], "Hooks$HookB", "setDaemon", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 30, 16, stack[0], "java/lang/Runtime", "addShutdownHook", "(Ljava/lang/Thread;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			obj, thrown := t.NewObject(m, 27, "Hooks$Throwing")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 38, 29, stack[2], "Hooks$Throwing", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 41, 16, stack[0], "java/lang/Runtime", "addShutdownHook", "(Ljava/lang/Thread;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 30, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{109, 97, 105, 110})
			{
			out := t.InvokeVirtual(m, 49, 38, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(7)
			{
			out := t.InvokeStatic(m, 54, 44, "java/lang/System", "exit", "(I)V", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod12(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, 2, "java/lang/Object")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, 1, stack[1], "java/lang/Object", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutStatic(m, 48, "Hooks", "BARRIER", "Ljava/lang/Object;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

