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

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			thrown := t.PutStatic(m, 13, "WaitSet", "sleeping", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 7, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.LongBitsValue(0x7fffffffffffffff)
			stack[1] = rt.IntValue(999999)
			{
			out := t.InvokeStatic(m, 15, 20, "java/lang/Thread", "sleep", "(JI)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 24, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			stack[2] = t.StringConstant([]uint16{108, 111, 110, 103, 32, 115, 108, 101, 101, 112, 32, 114, 101, 116, 117, 114, 110, 101, 100, 32, 119, 105, 116, 104, 111, 117, 116, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeSpecial(m, 24, 28, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			if hpc, ok := t.HandleThrow(m, 27, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 1:
			locals[2] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 29, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			out := t.InvokeStatic(m, 32, 33, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 32, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			block = 2
			continue
		case 2:
			{
			out := t.InvokeStatic(m, 38, 37, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 38, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			thrown := t.PutStatic(m, 41, "WaitSet", "sleepInterrupted", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 49, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			thrown := t.MonitorExit(m, 53, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 53, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 7
			continue
		case 6:
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 59, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 59, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
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
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 1
			case 57:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 7:
			return rt.Normal()
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			thrown := t.PutField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			out := t.InvokeSpecial(m, 6, 7, stack[0], "java/lang/Thread", "<init>", "()V", stack[1:1])
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
	var locals [6]rt.Value
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
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := t.MonitorEnter(m, 9, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			locals[2] = stack[1]
			{
			thrown := t.MonitorEnter(m, 19, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			{
			v, thrown := t.GetField(m, 19, stack[1], "WaitSet$Trial", "ready", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 25, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			thrown := t.PutField(m, 19, stack[0], "WaitSet$Trial", "ready", "I", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 23, stack[0], "WaitSet$Trial", "signalled", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 37, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeVirtual(m, 50, 27, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 3:
			stack[0] = locals[0]
			{
			out := t.InvokeStatic(m, 57, 32, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeVirtual(m, 60, 36, stack[1], "java/lang/Thread", "isInterrupted", "()Z", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := t.PutField(m, 40, stack[0], "WaitSet$Trial$Waiter", "pending", "Z", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 63, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := t.PutField(m, 43, stack[0], "WaitSet$Trial", "normalWake", "Z", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 71, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 78, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeVirtual(m, 81, 46, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			{
			thrown := t.PutField(m, 51, stack[0], "WaitSet$Trial$Waiter", "interrupted", "Z", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 93, 32, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 93, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeVirtual(m, 96, 36, stack[0], "java/lang/Thread", "isInterrupted", "()Z", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 96, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			stack[1] = t.StringConstant([]uint16{119, 97, 105, 116, 32, 102, 97, 105, 108, 101, 100, 32, 116, 111, 32, 99, 108, 101, 97, 114, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeStatic(m, 109, 56, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 109, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 113, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 116, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeStatic(m, 119, 62, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 119, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{119, 97, 105, 116, 32, 114, 101, 116, 117, 114, 110, 101, 100, 32, 119, 105, 116, 104, 111, 117, 116, 32, 105, 116, 115, 32, 109, 111, 110, 105, 116, 111, 114})
			{
			out := t.InvokeStatic(m, 124, 56, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 124, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := t.MonitorExit(m, 128, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 128, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			locals[4] = stack[0]
			stack[0] = locals[2]
			{
			thrown := t.MonitorExit(m, 135, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 135, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 138, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 10:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 1, stack[0], "WaitSet$Trial$Waiter", "this$0", "LWaitSet$Trial;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 140, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			v, thrown := t.GetField(m, 13, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 143, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			out := t.InvokeStatic(m, 146, 62, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 146, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{119, 97, 105, 116, 32, 102, 97, 105, 108, 101, 100, 32, 116, 111, 32, 114, 101, 115, 116, 111, 114, 101, 32, 114, 101, 99, 117, 114, 115, 105, 111, 110, 32, 100, 101, 112, 116, 104})
			{
			out := t.InvokeStatic(m, 151, 56, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 151, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			thrown := t.MonitorExit(m, 155, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 155, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
		case 11:
			locals[5] = stack[0]
			stack[0] = locals[1]
			{
			thrown := t.MonitorExit(m, 162, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 162, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
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
			if hpc, ok := t.HandleThrow(m, 165, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 87:
			block = 4
			case 132:
			block = 9
			case 159:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 12:
			return rt.Normal()
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			{
			out := t.InvokeSpecial(m, 1, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			obj, thrown := t.NewObject(m, 2, "java/lang/Object")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 9, 1, stack[2], "java/lang/Object", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [6]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject(m, 13, "WaitSet$Trial$Waiter")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			stack[2] = locals[0]
			{
			out := t.InvokeSpecial(m, 5, 15, stack[1], "WaitSet$Trial$Waiter", "<init>", "(LWaitSet$Trial;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 5, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			{
			obj, thrown := t.NewObject(m, 13, "WaitSet$Trial$Waiter")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			stack[2] = locals[0]
			{
			out := t.InvokeSpecial(m, 14, 15, stack[1], "WaitSet$Trial$Waiter", "<init>", "(LWaitSet$Trial;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[3] = stack[0]
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 19, 18, stack[0], "WaitSet$Trial$Waiter", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			out := t.InvokeVirtual(m, 23, 18, stack[0], "WaitSet$Trial$Waiter", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 23, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 21, stack[0], "WaitSet$Trial", "ready", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 27, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			stack[1] = rt.IntValue(2)
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			{
			out := t.InvokeStatic(m, 34, 25, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 3:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 41, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			stack[1] = stack[0]
			locals[4] = stack[1]
			{
			thrown := t.MonitorEnter(m, 47, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			thrown := t.PutField(m, 30, stack[0], "WaitSet$Trial", "signalled", "Z", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			if rt.AsInt(stack[0]) == 0 {
			block = 5
			continue
			}
			block = 4
			continue
		case 4:
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 58, 34, stack[0], "WaitSet$Trial$Waiter", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 58, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			v, thrown := t.GetField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 62, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			out := t.InvokeVirtual(m, 65, 37, stack[0], "java/lang/Object", "notify", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 65, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
		case 5:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 72, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			out := t.InvokeVirtual(m, 75, 37, stack[0], "java/lang/Object", "notify", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			out := t.InvokeVirtual(m, 79, 34, stack[0], "WaitSet$Trial$Waiter", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 79, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 40, stack[0], "WaitSet$Trial", "normalWake", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 83, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			if rt.AsInt(stack[0]) != 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			out := t.InvokeVirtual(m, 93, 43, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 93, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
		case 8:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "WaitSet$Trial", "monitor", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 100, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			out := t.InvokeVirtual(m, 103, 46, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 103, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			thrown := t.MonitorExit(m, 108, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 108, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			locals[5] = stack[0]
			stack[0] = locals[4]
			{
			thrown := t.MonitorExit(m, 116, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 116, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			if hpc, ok := t.HandleThrow(m, 119, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 10:
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 121, 49, stack[0], "WaitSet$Trial$Waiter", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 121, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			out := t.InvokeVirtual(m, 125, 49, stack[0], "WaitSet$Trial$Waiter", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 125, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
			block = 9
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
			v, thrown := t.GetField(m, 52, stack[0], "WaitSet$Trial$Waiter", "interrupted", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 129, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			stack[1] = locals[2]
			{
			v, thrown := t.GetField(m, 55, stack[1], "WaitSet$Trial$Waiter", "pending", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 133, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 14
			continue
			}
			block = 11
			continue
		case 11:
			stack[0] = locals[3]
			{
			v, thrown := t.GetField(m, 52, stack[0], "WaitSet$Trial$Waiter", "interrupted", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 140, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			if rt.AsInt(stack[0]) != 0 {
			block = 14
			continue
			}
			block = 12
			continue
		case 12:
			stack[0] = locals[3]
			{
			v, thrown := t.GetField(m, 55, stack[0], "WaitSet$Trial$Waiter", "pending", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 147, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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
			if rt.AsInt(stack[0]) != 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			stack[0] = rt.IntValue(1)
			block = 15
			continue
		case 14:
			stack[0] = rt.IntValue(0)
			block = 15
			continue
		case 15:
			stack[1] = t.StringConstant([]uint16{110, 111, 116, 105, 102, 121, 47, 105, 110, 116, 101, 114, 114, 117, 112, 116, 32, 114, 101, 116, 117, 114, 110, 101, 100, 32, 97, 32, 102, 111, 114, 98, 105, 100, 100, 101, 110, 32, 111, 117, 116, 99, 111, 109, 101})
			{
			out := t.InvokeStatic(m, 160, 60, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 160, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 112:
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

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeVirtual(m, 3, 13, stack[0], "java/lang/Object", "notify", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{110, 111, 116, 105, 102, 121, 32, 119, 105, 116, 104, 111, 117, 116, 32, 111, 119, 110, 101, 114, 115, 104, 105, 112})
			{
			out := t.InvokeSpecial(m, 12, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 12, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 1:
			locals[0] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 17, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeVirtual(m, 20, 25, stack[0], "java/lang/Object", "notifyAll", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 23, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{110, 111, 116, 105, 102, 121, 65, 108, 108, 32, 119, 105, 116, 104, 111, 117, 116, 32, 111, 119, 110, 101, 114, 115, 104, 105, 112})
			{
			out := t.InvokeSpecial(m, 29, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 29, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 32, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 2:
			locals[0] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.LongValue(1)
			{
			out := t.InvokeVirtual(m, 38, 30, stack[0], "java/lang/Object", "wait", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 38, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 41, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{119, 97, 105, 116, 32, 119, 105, 116, 104, 111, 117, 116, 32, 111, 119, 110, 101, 114, 115, 104, 105, 112})
			{
			out := t.InvokeSpecial(m, 47, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 3:
			locals[0] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 52, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			locals[0] = stack[1]
			{
			thrown := t.MonitorEnter(m, 57, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 58, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.LongValue(0)
			stack[2] = rt.IntValue(1000000)
			{
			out := t.InvokeVirtual(m, 64, 37, stack[0], "java/lang/Object", "wait", "(JI)V", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{105, 110, 118, 97, 108, 105, 100, 32, 110, 97, 110, 111, 115, 32, 97, 99, 99, 101, 112, 116, 101, 100})
			{
			out := t.InvokeSpecial(m, 73, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 73, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 76, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 4:
			locals[1] = stack[0]
			{
			out := t.InvokeStatic(m, 78, 44, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 78, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeVirtual(m, 81, 50, stack[0], "java/lang/Thread", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 84, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeVirtual(m, 87, 53, stack[0], "java/lang/Object", "wait", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 87, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{112, 114, 101, 45, 105, 110, 116, 101, 114, 114, 117, 112, 116, 101, 100, 32, 119, 97, 105, 116, 32, 114, 101, 116, 117, 114, 110, 101, 100, 32, 110, 111, 114, 109, 97, 108, 108, 121})
			{
			out := t.InvokeSpecial(m, 96, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 96, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 99, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 5:
			locals[1] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 101, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeStatic(m, 104, 59, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 104, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			block = 8
			continue
			}
			block = 6
			continue
		case 6:
			{
			out := t.InvokeStatic(m, 110, 63, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 110, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[1] = t.StringConstant([]uint16{112, 114, 101, 45, 105, 110, 116, 101, 114, 114, 117, 112, 116, 101, 100, 32, 119, 97, 105, 116, 32, 115, 101, 109, 97, 110, 116, 105, 99, 115})
			{
			out := t.InvokeStatic(m, 123, 69, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 123, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 126, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.LongValue(0)
			stack[2] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 131, 37, stack[0], "java/lang/Object", "wait", "(JI)V", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 131, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 134, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeStatic(m, 137, 59, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 137, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{116, 105, 109, 101, 100, 32, 119, 97, 105, 116, 32, 108, 111, 115, 116, 32, 109, 111, 110, 105, 116, 111, 114})
			{
			out := t.InvokeStatic(m, 142, 69, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 142, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 145, 44, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 145, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeVirtual(m, 148, 50, stack[0], "java/lang/Thread", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 148, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.LongBitsValue(0x0000000000002710)
			{
			out := t.InvokeStatic(m, 154, 79, "java/lang/Thread", "sleep", "(J)V", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 154, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 16, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 157, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			stack[2] = t.StringConstant([]uint16{112, 114, 101, 45, 105, 110, 116, 101, 114, 114, 117, 112, 116, 101, 100, 32, 115, 108, 101, 101, 112, 32, 114, 101, 116, 117, 114, 110, 101, 100, 32, 110, 111, 114, 109, 97, 108, 108, 121})
			{
			out := t.InvokeSpecial(m, 163, 20, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 163, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			if hpc, ok := t.HandleThrow(m, 166, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 10:
			locals[1] = stack[0]
			{
			v, thrown := t.GetStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 168, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			out := t.InvokeStatic(m, 171, 59, "java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 171, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			block = 13
			continue
			}
			block = 11
			continue
		case 11:
			{
			out := t.InvokeStatic(m, 177, 63, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 177, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			block = 13
			continue
			}
			block = 12
			continue
		case 12:
			stack[0] = rt.IntValue(1)
			block = 14
			continue
		case 13:
			stack[0] = rt.IntValue(0)
			block = 14
			continue
		case 14:
			stack[1] = t.StringConstant([]uint16{115, 108, 101, 101, 112, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116, 32, 115, 101, 109, 97, 110, 116, 105, 99, 115})
			{
			out := t.InvokeStatic(m, 190, 69, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 190, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			thrown := t.MonitorExit(m, 194, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 194, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			locals[2] = stack[0]
			stack[0] = locals[0]
			{
			thrown := t.MonitorExit(m, 200, stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 200, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
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
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 202, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 16:
			block = 1
			case 33:
			block = 2
			case 51:
			block = 3
			case 77:
			block = 4
			case 100:
			block = 5
			case 167:
			block = 10
			case 198:
			block = 15
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 16:
			return rt.Normal()
		}
	}
}

func cattyMethod10(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			out := t.InvokeStatic(m, 0, 86, "WaitSet", "ownershipAndTimeouts", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 89, "WaitSet$Sleeper")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 7, 91, stack[1], "WaitSet$Sleeper", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 12, 92, stack[0], "WaitSet$Sleeper", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 95, "WaitSet", "sleeping", "Z")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			{
			out := t.InvokeStatic(m, 21, 99, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 3:
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 28, 102, stack[0], "WaitSet$Sleeper", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 32, 103, stack[0], "WaitSet$Sleeper", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 106, "WaitSet", "sleepInterrupted", "Z")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{115, 108, 101, 101, 112, 32, 100, 105, 100, 32, 110, 111, 116, 32, 114, 101, 116, 97, 105, 110, 32, 111, 119, 110, 101, 114, 115, 104, 105, 112, 32, 97, 110, 100, 32, 99, 108, 101, 97, 114, 32, 116, 104, 101, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeStatic(m, 40, 69, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			block = 4
			continue
		case 4:
			stack[0] = locals[2]
			stack[1] = rt.IntValue(16)
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 9
			continue
			}
			block = 5
			continue
		case 5:
			{
			obj, thrown := t.NewObject(m, 111, "WaitSet$Trial")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 55, 113, stack[1], "WaitSet$Trial", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[1] = locals[2]
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) & rt.AsInt(stack[2]))
			if rt.AsInt(stack[1]) != 0 {
			block = 7
			continue
			}
			block = 6
			continue
		case 6:
			stack[1] = rt.IntValue(1)
			block = 8
			continue
		case 7:
			stack[1] = rt.IntValue(0)
			block = 8
			continue
		case 8:
			{
			out := t.InvokeVirtual(m, 69, 114, stack[0], "WaitSet$Trial", "run", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1)
			block = 4
			continue
		case 9:
			{
			v, thrown := t.GetStatic(m, 118, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{119, 97, 105, 116, 45, 115, 101, 116})
			{
			out := t.InvokeVirtual(m, 83, 126, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
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
			thrown := t.PutStatic(m, 7, "WaitSet", "LOCK", "Ljava/lang/Object;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

