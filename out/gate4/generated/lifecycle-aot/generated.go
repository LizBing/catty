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
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = locals[1]
			{
			thrown := t.PutField(m, 1, stack[0], "Lifecycle$1", "val$worker", "LLifecycle$Worker;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 6, 7, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
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
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 1, stack[0], "Lifecycle$1", "val$worker", "LLifecycle$Worker;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 10:
			block = 1
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
			out := t.InvokeVirtual(m, 4, 13, stack[0], "Lifecycle$Worker", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 10:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 2
			continue
		case 1:
			locals[1] = stack[0]
			return rt.Normal()
		case 2:
			{
			obj, thrown := t.NewObject(m, 20, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 12, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 10:
			block = 1
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
			stack[2] = t.StringConstant([]uint16{97, 32, 108, 105, 118, 101, 32, 116, 104, 114, 101, 97, 100, 32, 115, 116, 97, 114, 116, 101, 100, 32, 116, 119, 105, 99, 101})
			{
			out := t.InvokeSpecial(m, 18, 24, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 10:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 10:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
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
			stack[1] = t.StringConstant([]uint16{99, 111, 117, 110, 116, 101, 100})
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
	var locals [1]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = stack[0]
			{
			v, thrown := t.GetField(m, 9, stack[1], "Lifecycle$Counted", "runs", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			thrown := t.PutField(m, 9, stack[0], "Lifecycle$Counted", "runs", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
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
			stack[0] = locals[0]
			stack[1] = locals[1]
			{
			thrown := t.PutField(m, 7, stack[0], "Lifecycle$Joiner", "target", "Ljava/lang/Thread;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [4]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 13, "Lifecycle", "JOIN_GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			thrown := t.MonitorEnter(m, 5, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 5, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 19, "Lifecycle", "joining", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 7, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "Lifecycle$Joiner", "target", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			out := t.InvokeVirtual(m, 14, 23, stack[0], "java/lang/Thread", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 6
			continue
		case 1:
			locals[2] = stack[0]
			{
			out := t.InvokeStatic(m, 21, 30, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			out := t.InvokeVirtual(m, 24, 34, stack[0], "java/lang/Thread", "isInterrupted", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 4
			continue
			}
			block = 2
			continue
		case 2:
			{
			v, thrown := t.GetStatic(m, 13, "Lifecycle", "JOIN_GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			out := t.InvokeStatic(m, 33, 38, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 33, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = rt.IntValue(1)
			block = 5
			continue
		case 4:
			stack[0] = rt.IntValue(0)
			block = 5
			continue
		case 5:
			{
			thrown := t.PutStatic(m, 42, "Lifecycle", "joinInterrupted", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 6
			continue
		case 6:
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 48, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 48, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 54, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 54, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
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
			if hpc, ok := t.HandleThrow(m, 56, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 20:
			block = 1
			case 52:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 8:
			return rt.Normal()
		}
	}
}

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Thread", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[1]
			{
			thrown := t.PutField(m, 7, stack[0], "Lifecycle$Starter", "target", "LLifecycle$Counted;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 13, "Lifecycle", "startRace", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 27:
			block = 3
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			{
			out := t.InvokeStatic(m, 6, 19, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 27:
			block = 3
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
		case 2:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "Lifecycle$Starter", "target", "LLifecycle$Counted;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 13, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 27:
			block = 3
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
			out := t.InvokeVirtual(m, 16, 22, stack[0], "Lifecycle$Counted", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 16, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 27:
			block = 3
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			{
			thrown := t.PutField(m, 27, stack[0], "Lifecycle$Starter", "won", "Z", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 27:
			block = 3
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
		case 3:
			locals[1] = stack[0]
			block = 4
			continue
		case 4:
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
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Thread", "<init>", "()V", stack[1:1])
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
			out := t.InvokeStatic(m, 0, 7, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			thrown := t.PutStatic(m, 11, "Lifecycle", "observed", "Ljava/lang/Thread;", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 11, "Lifecycle", "observed", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			v, thrown := t.GetStatic(m, 17, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			if stack[0].Ref != stack[1].Ref {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 20, "Lifecycle", "directCalls", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			thrown := t.PutStatic(m, 20, "Lifecycle", "directCalls", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		case 2:
			{
			v, thrown := t.GetStatic(m, 24, "Lifecycle", "GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			thrown := t.MonitorEnter(m, 29, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 29, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 28, "Lifecycle", "started", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 31, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 3
			continue
		case 3:
			{
			v, thrown := t.GetStatic(m, 32, "Lifecycle", "release", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 6
			continue
			}
			block = 4
			continue
		case 4:
			{
			v, thrown := t.GetStatic(m, 24, "Lifecycle", "GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 40, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			out := t.InvokeVirtual(m, 43, 35, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 43, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 3
			continue
		case 5:
			locals[2] = stack[0]
			{
			obj, thrown := t.NewObject(m, 42, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			stack[2] = t.StringConstant([]uint16{117, 110, 101, 120, 112, 101, 99, 116, 101, 100, 32, 119, 111, 114, 107, 101, 114, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeSpecial(m, 56, 46, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 56, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 59, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 6:
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 61, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 67, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			case 49:
			block = 5
			case 65:
			block = 7
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 8:
			{
			v, thrown := t.GetStatic(m, 49, "Lifecycle", "input", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 70, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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
			thrown := t.PutStatic(m, 52, "Lifecycle", "result", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			case 65:
			block = 7
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

func cattyMethod13(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [5]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, 7, "Lifecycle$Counted")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, 9, stack[1], "Lifecycle$Counted", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[0] = stack[0]
			stack[0] = locals[0]
			{
			out := t.InvokeVirtual(m, 9, 10, stack[0], "Lifecycle$Counted", "threadId", "()J", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[1] = stack[0]
			{
			obj, thrown := t.NewObject(m, 14, "Lifecycle$Starter")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = locals[0]
			{
			out := t.InvokeSpecial(m, 18, 16, stack[1], "Lifecycle$Starter", "<init>", "(LLifecycle$Counted;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[3] = stack[0]
			{
			obj, thrown := t.NewObject(m, 14, "Lifecycle$Starter")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = locals[0]
			{
			out := t.InvokeSpecial(m, 27, 16, stack[1], "Lifecycle$Starter", "<init>", "(LLifecycle$Counted;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[4] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 33, 19, stack[0], "Lifecycle$Starter", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			{
			out := t.InvokeVirtual(m, 38, 19, stack[0], "Lifecycle$Starter", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 22, "Lifecycle", "startRace", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 46, 28, stack[0], "Lifecycle$Starter", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			{
			out := t.InvokeVirtual(m, 51, 28, stack[0], "Lifecycle$Starter", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			out := t.InvokeVirtual(m, 55, 31, stack[0], "Lifecycle$Counted", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			{
			v, thrown := t.GetField(m, 32, stack[0], "Lifecycle$Starter", "won", "Z")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			v, thrown := t.GetField(m, 32, stack[1], "Lifecycle$Starter", "won", "Z")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 3
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 35, stack[0], "Lifecycle$Counted", "runs", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = rt.IntValue(1)
			block = 4
			continue
		case 3:
			stack[0] = rt.IntValue(0)
			block = 4
			continue
		case 4:
			stack[1] = t.StringConstant([]uint16{116, 119, 111, 32, 102, 105, 114, 115, 116, 32, 115, 116, 97, 114, 116, 115, 32, 115, 117, 99, 99, 101, 101, 100, 101, 100})
			{
			out := t.InvokeStatic(m, 85, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			stack[1] = rt.LongValue(0)
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) <= 0 {
			block = 8
			continue
			}
			block = 5
			continue
		case 5:
			stack[0] = locals[1]
			stack[1] = locals[0]
			{
			out := t.InvokeVirtual(m, 96, 10, stack[1], "Lifecycle$Counted", "threadId", "()J", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) != 0 {
			block = 8
			continue
			}
			block = 6
			continue
		case 6:
			stack[0] = locals[1]
			{
			v, thrown := t.GetStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 107, 51, stack[1], "java/lang/Thread", "threadId", "()J", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) == 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			stack[0] = rt.IntValue(1)
			block = 9
			continue
		case 8:
			stack[0] = rt.IntValue(0)
			block = 9
			continue
		case 9:
			stack[1] = t.StringConstant([]uint16{105, 110, 99, 111, 114, 114, 101, 99, 116, 32, 115, 116, 97, 98, 108, 101, 32, 116, 104, 114, 101, 97, 100, 32, 73, 68})
			{
			out := t.InvokeStatic(m, 121, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [6]rt.Value
	var stack [5]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			out := t.InvokeStatic(m, 0, 56, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			thrown := t.PutStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 60, "Lifecycle$Worker")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeSpecial(m, 10, 62, stack[1], "Lifecycle$Worker", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 10, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			{
			v, thrown := t.GetStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeStatic(m, 17, 56, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 17, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			if stack[0].Ref != stack[1].Ref {
			block = 3
			continue
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 23, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 26, 63, stack[0], "java/lang/Thread", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 26, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = rt.IntValue(1)
			block = 4
			continue
		case 3:
			stack[0] = rt.IntValue(0)
			block = 4
			continue
		case 4:
			stack[1] = t.StringConstant([]uint16{109, 97, 105, 110, 32, 116, 104, 114, 101, 97, 100, 32, 105, 100, 101, 110, 116, 105, 116, 121})
			{
			out := t.InvokeStatic(m, 39, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 39, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 43, 69, stack[0], "Lifecycle$Worker", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 43, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 47, 70, stack[0], "Lifecycle$Worker", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			stack[0] = rt.IntValue(1)
			block = 7
			continue
		case 6:
			stack[0] = rt.IntValue(0)
			block = 7
			continue
		case 7:
			stack[1] = t.StringConstant([]uint16{106, 111, 105, 110, 105, 110, 103, 32, 97, 32, 110, 101, 119, 32, 116, 104, 114, 101, 97, 100})
			{
			out := t.InvokeStatic(m, 60, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 64, 73, stack[0], "Lifecycle$Worker", "run", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 76, "Lifecycle", "directCalls", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 10
			continue
			}
			block = 8
			continue
		case 8:
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 75, 70, stack[0], "Lifecycle$Worker", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			stack[0] = rt.IntValue(1)
			block = 11
			continue
		case 10:
			stack[0] = rt.IntValue(0)
			block = 11
			continue
		case 11:
			stack[1] = t.StringConstant([]uint16{100, 105, 114, 101, 99, 116, 32, 114, 117, 110, 32, 109, 117, 115, 116, 32, 114, 101, 109, 97, 105, 110, 32, 115, 121, 110, 99, 104, 114, 111, 110, 111, 117, 115})
			{
			out := t.InvokeStatic(m, 88, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 88, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(41)
			{
			thrown := t.PutStatic(m, 81, "Lifecycle", "input", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 93, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 97, 84, stack[0], "Lifecycle$Worker", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 97, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 12
			continue
		case 12:
			{
			v, thrown := t.GetStatic(m, 85, "Lifecycle", "started", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 100, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			{
			out := t.InvokeStatic(m, 106, 88, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 106, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 12
			continue
		case 14:
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 113, 70, stack[0], "Lifecycle$Worker", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 113, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{119, 111, 114, 107, 101, 114, 32, 115, 104, 111, 117, 108, 100, 32, 98, 101, 32, 97, 108, 105, 118, 101, 32, 97, 116, 32, 103, 97, 116, 101})
			{
			out := t.InvokeStatic(m, 118, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 118, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			stack[1] = rt.LongValue(1)
			stack[2] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 124, 93, stack[0], "Lifecycle$Worker", "join", "(JI)V", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 124, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 128, 70, stack[0], "Lifecycle$Worker", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 128, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{116, 105, 109, 101, 100, 32, 106, 111, 105, 110, 32, 115, 104, 111, 117, 108, 100, 32, 114, 101, 116, 117, 114, 110, 32, 119, 105, 116, 104, 32, 119, 111, 114, 107, 101, 114, 32, 115, 116, 105, 108, 108, 32, 97, 108, 105, 118, 101})
			{
			out := t.InvokeStatic(m, 133, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 133, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 52, "java/lang/Thread")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 136, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			obj, thrown := t.NewObject(m, 98, "Lifecycle$Joiner")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 140, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			stack[4] = locals[1]
			{
			out := t.InvokeSpecial(m, 145, 100, stack[3], "Lifecycle$Joiner", "<init>", "(Ljava/lang/Thread;)V", stack[4:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 145, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeSpecial(m, 148, 103, stack[1], "java/lang/Thread", "<init>", "(Ljava/lang/Runnable;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 148, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 153, 106, stack[0], "java/lang/Thread", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 153, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 15
			continue
		case 15:
			{
			v, thrown := t.GetStatic(m, 107, "Lifecycle", "joining", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 156, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 17
			continue
			}
			block = 16
			continue
		case 16:
			{
			out := t.InvokeStatic(m, 162, 88, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 162, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 15
			continue
		case 17:
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 169, 110, stack[0], "java/lang/Thread", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 169, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 173, 113, stack[0], "java/lang/Thread", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 173, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 114, "Lifecycle", "joinInterrupted", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 176, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{106, 111, 105, 110, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116, 32, 115, 104, 111, 117, 108, 100, 32, 99, 108, 101, 97, 114, 32, 115, 116, 97, 116, 117, 115, 32, 97, 110, 100, 32, 114, 101, 116, 97, 105, 110, 32, 111, 116, 104, 101, 114, 32, 109, 111, 110, 105, 116, 111, 114, 115})
			{
			out := t.InvokeStatic(m, 181, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 181, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 52, "java/lang/Thread")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 184, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			obj, thrown := t.NewObject(m, 119, "Lifecycle$1")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 188, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			stack[4] = locals[1]
			{
			out := t.InvokeSpecial(m, 193, 121, stack[3], "Lifecycle$1", "<init>", "(LLifecycle$Worker;)V", stack[4:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 193, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeSpecial(m, 196, 103, stack[1], "java/lang/Thread", "<init>", "(Ljava/lang/Runnable;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 196, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[3] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 201, 106, stack[0], "java/lang/Thread", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 201, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 205, 113, stack[0], "java/lang/Thread", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 205, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 124, "Lifecycle", "GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 208, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			locals[4] = stack[1]
			{
			thrown := t.MonitorEnter(m, 214, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 214, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 128, "Lifecycle", "release", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 216, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 124, "Lifecycle", "GATE", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 219, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 222, 131, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 222, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			thrown := t.MonitorExit(m, 227, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 227, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 19
			continue
		case 18:
			locals[5] = stack[0]
			stack[0] = locals[4]
			{
			thrown := t.MonitorExit(m, 235, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 235, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 238, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 19:
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 240, 69, stack[0], "Lifecycle$Worker", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 240, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 134, "Lifecycle", "observed", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 243, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			if stack[0].Ref != stack[1].Ref {
			block = 23
			continue
			}
			block = 20
			continue
		case 20:
			{
			v, thrown := t.GetStatic(m, 137, "Lifecycle", "result", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 250, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(42)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 23
			continue
			}
			block = 21
			continue
		case 21:
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 259, 70, stack[0], "Lifecycle$Worker", "isAlive", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 259, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 23
			continue
			}
			block = 22
			continue
		case 22:
			stack[0] = rt.IntValue(1)
			block = 24
			continue
		case 23:
			stack[0] = rt.IntValue(0)
			block = 24
			continue
		case 24:
			stack[1] = t.StringConstant([]uint16{119, 111, 114, 107, 101, 114, 32, 105, 100, 101, 110, 116, 105, 116, 121, 32, 97, 110, 100, 32, 106, 111, 105, 110, 32, 112, 117, 98, 108, 105, 99, 97, 116, 105, 111, 110})
			{
			out := t.InvokeStatic(m, 272, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 272, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 276, 84, stack[0], "Lifecycle$Worker", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 276, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 142, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 279, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			stack[2] = t.StringConstant([]uint16{116, 101, 114, 109, 105, 110, 97, 116, 101, 100, 32, 116, 104, 114, 101, 97, 100, 32, 114, 101, 115, 116, 97, 114, 116, 101, 100})
			{
			out := t.InvokeSpecial(m, 285, 146, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 285, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 288, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 25:
			locals[4] = stack[0]
			{
			out := t.InvokeStatic(m, 291, 56, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 291, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 294, 110, stack[0], "java/lang/Thread", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 294, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 297, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 300, 151, stack[0], "java/lang/Thread", "isInterrupted", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 300, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 28
			continue
			}
			block = 26
			continue
		case 26:
			{
			v, thrown := t.GetStatic(m, 47, "Lifecycle", "main", "Ljava/lang/Thread;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 306, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 309, 151, stack[0], "java/lang/Thread", "isInterrupted", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 309, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 28
			continue
			}
			block = 27
			continue
		case 27:
			stack[0] = rt.IntValue(1)
			block = 29
			continue
		case 28:
			stack[0] = rt.IntValue(0)
			block = 29
			continue
		case 29:
			stack[1] = t.StringConstant([]uint16{115, 116, 97, 116, 117, 115, 32, 113, 117, 101, 114, 121, 32, 99, 108, 101, 97, 114, 101, 100, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeStatic(m, 322, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 322, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 325, 156, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 325, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 32
			continue
			}
			block = 30
			continue
		case 30:
			{
			out := t.InvokeStatic(m, 331, 156, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 331, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 32
			continue
			}
			block = 31
			continue
		case 31:
			stack[0] = rt.IntValue(1)
			block = 33
			continue
		case 32:
			stack[0] = rt.IntValue(0)
			block = 33
			continue
		case 33:
			stack[1] = t.StringConstant([]uint16{116, 101, 115, 116, 45, 97, 110, 100, 45, 99, 108, 101, 97, 114, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeStatic(m, 344, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 344, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 52, "java/lang/Thread")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 347, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeSpecial(m, 351, 161, stack[1], "java/lang/Thread", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 351, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[4] = stack[0]
			stack[0] = locals[4]
			{
			out := t.InvokeVirtual(m, 358, 162, stack[0], "java/lang/Thread", "isDaemon", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 358, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 35
			continue
			}
			block = 34
			continue
		case 34:
			stack[0] = rt.IntValue(1)
			block = 36
			continue
		case 35:
			stack[0] = rt.IntValue(0)
			block = 36
			continue
		case 36:
			stack[1] = t.StringConstant([]uint16{109, 97, 105, 110, 32, 105, 115, 32, 110, 111, 110, 45, 100, 97, 101, 109, 111, 110})
			{
			out := t.InvokeStatic(m, 371, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 371, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 377, 167, stack[0], "java/lang/Thread", "setDaemon", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 377, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			out := t.InvokeVirtual(m, 382, 162, stack[0], "java/lang/Thread", "isDaemon", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 382, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{110, 101, 119, 32, 116, 104, 114, 101, 97, 100, 32, 100, 97, 101, 109, 111, 110, 32, 115, 101, 116, 116, 105, 110, 103})
			{
			out := t.InvokeStatic(m, 387, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 387, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 390, 173, "Lifecycle", "simultaneousStart", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 390, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			stack[1] = rt.LongBitsValue(0xffffffffffffffff)
			{
			out := t.InvokeVirtual(m, 397, 178, stack[0], "Lifecycle$Worker", "join", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 397, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 142, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 400, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			stack[2] = t.StringConstant([]uint16{110, 101, 103, 97, 116, 105, 118, 101, 32, 106, 111, 105, 110, 32, 97, 99, 99, 101, 112, 116, 101, 100})
			{
			out := t.InvokeSpecial(m, 406, 146, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 406, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 409, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 37:
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 185, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 412, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{108, 105, 102, 101, 99, 121, 99, 108, 101})
			{
			out := t.InvokeVirtual(m, 417, 193, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 417, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 231:
			block = 18
			case 289:
			block = 25
			case 410:
			block = 37
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
			thrown := t.PutStatic(m, 124, "Lifecycle", "GATE", "Ljava/lang/Object;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 2, "java/lang/Object")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 14, 1, stack[1], "java/lang/Object", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutStatic(m, 198, "Lifecycle", "JOIN_GATE", "Ljava/lang/Object;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

