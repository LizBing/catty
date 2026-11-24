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
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 7, "Initialization", "attempts", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, 7, "Initialization", "attempts", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 8, 13, "Initialization", "blockInitializer", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 16, "Initialization", "TOKEN", "Ljava/lang/RuntimeException;")
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

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 20, "Initialization$Broken", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
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
			out := t.InvokeStatic(m, 0, 25, "Initialization$Broken", "initialize", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, 20, "Initialization$Broken", "value", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	locals[2] = args[1]
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
			thrown := t.PutField(m, 7, stack[0], "Initialization$Reader", "owner", "Z", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[2]
			{
			thrown := t.PutField(m, 13, stack[0], "Initialization$Reader", "broken", "Z", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "Initialization$Reader", "owner", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 1, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			out := t.InvokeStatic(m, 7, 16, "java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 7, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			out := t.InvokeVirtual(m, 10, 20, stack[0], "java/lang/Thread", "interrupt", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 10, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			thrown := t.PutStatic(m, 23, "Initialization", "attempting", "Z", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
		case 2:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 13, stack[0], "Initialization$Reader", "broken", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			{
			out := t.InvokeStatic(m, 24, 28, "Initialization$Broken", "read", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			block = 5
			continue
		case 4:
			{
			out := t.InvokeStatic(m, 30, 34, "Initialization$Ready", "read", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			block = 5
			continue
		case 5:
			locals[1] = stack[0]
			block = 11
			continue
		case 6:
			locals[2] = stack[0]
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 39, 39, stack[0], "java/lang/ExceptionInInitializerError", "getCause", "()Ljava/lang/Throwable;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 39, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			v, thrown := t.GetStatic(m, 43, "Initialization", "TOKEN", "Ljava/lang/RuntimeException;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 42, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			stack[1] = t.StringConstant([]uint16{105, 110, 105, 116, 105, 97, 108, 105, 122, 101, 114, 32, 108, 111, 115, 116, 32, 111, 114, 105, 103, 105, 110, 97, 108, 32, 99, 97, 117, 115, 101})
			{
			out := t.InvokeStatic(m, 55, 49, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 55, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(-1)
			locals[1] = stack[0]
			block = 11
			continue
		case 10:
			locals[2] = stack[0]
			stack[0] = rt.IntValue(-2)
			locals[1] = stack[0]
			block = 11
			continue
		case 11:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "Initialization$Reader", "owner", "Z")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 68, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) == 0 {
			block = 13
			continue
			}
			block = 12
			continue
		case 12:
			stack[0] = locals[1]
			{
			thrown := t.PutStatic(m, 57, "Initialization", "ownerResult", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 14
			continue
		case 13:
			{
			out := t.InvokeStatic(m, 81, 61, "java/lang/Thread", "interrupted", "()Z", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{105, 110, 105, 116, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110, 32, 99, 111, 110, 115, 117, 109, 101, 100, 32, 105, 110, 116, 101, 114, 114, 117, 112, 116})
			{
			out := t.InvokeStatic(m, 86, 49, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 86, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
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
			thrown := t.PutStatic(m, 67, "Initialization", "contenderResult", "I", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 37:
			block = 6
			case 63:
			block = 10
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 14
			continue
		case 14:
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
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 7, "Initialization", "attempts", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, 7, "Initialization", "attempts", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 8, 13, "Initialization$Ready", "read", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = rt.IntValue(1)
			block = 3
			continue
		case 2:
			stack[0] = rt.IntValue(0)
			block = 3
			continue
		case 3:
			stack[1] = t.StringConstant([]uint16{114, 101, 99, 117, 114, 115, 105, 118, 101, 32, 105, 110, 105, 116, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110, 32, 115, 104, 111, 117, 108, 100, 32, 115, 101, 101, 32, 116, 104, 101, 32, 100, 101, 102, 97, 117, 108, 116})
			{
			out := t.InvokeStatic(m, 21, 21, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 24, 27, "Initialization", "blockInitializer", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(73)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod10(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 30, "Initialization$Ready", "value", "I")
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

func cattyMethod12(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 7, "Initialization", "initializing", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 13, "Initialization", "release", "Z")
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
			out := t.InvokeStatic(m, 10, 16, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 3:
			return rt.Normal()
		}
	}
}

func cattyMethod13(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			{
			thrown := t.PutStatic(m, 7, "Initialization", "initializing", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			thrown := t.PutStatic(m, 21, "Initialization", "attempting", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			thrown := t.PutStatic(m, 13, "Initialization", "release", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 24, "Initialization$Reader")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = rt.IntValue(1)
			stack[3] = locals[0]
			{
			out := t.InvokeSpecial(m, 18, 26, stack[1], "Initialization$Reader", "<init>", "(ZZ)V", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			{
			obj, thrown := t.NewObject(m, 24, "Initialization$Reader")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = rt.IntValue(0)
			stack[3] = locals[0]
			{
			out := t.InvokeSpecial(m, 28, 26, stack[1], "Initialization$Reader", "<init>", "(ZZ)V", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 33, 29, stack[0], "Initialization$Reader", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 7, "Initialization", "initializing", "Z")
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
			out := t.InvokeStatic(m, 42, 16, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 1
			continue
		case 3:
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 49, 29, stack[0], "Initialization$Reader", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 4:
			{
			v, thrown := t.GetStatic(m, 21, "Initialization", "attempting", "Z")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if rt.AsInt(stack[0]) != 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			{
			out := t.InvokeStatic(m, 58, 16, "java/lang/Thread", "yield", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 4
			continue
		case 6:
			stack[0] = rt.IntValue(1)
			{
			thrown := t.PutStatic(m, 13, "Initialization", "release", "Z", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 69, 32, stack[0], "Initialization$Reader", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 73, 32, stack[0], "Initialization$Reader", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 35, "Initialization", "ownerResult", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			if rt.AsInt(stack[1]) == 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			stack[1] = rt.IntValue(-1)
			block = 9
			continue
		case 8:
			stack[1] = rt.IntValue(73)
			block = 9
			continue
		case 9:
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 11
			continue
			}
			block = 10
			continue
		case 10:
			stack[0] = rt.IntValue(1)
			block = 12
			continue
		case 11:
			stack[0] = rt.IntValue(0)
			block = 12
			continue
		case 12:
			stack[1] = t.StringConstant([]uint16{98, 97, 100, 32, 105, 110, 105, 116, 105, 97, 108, 105, 122, 101, 114, 32, 111, 119, 110, 105, 110, 103, 32, 111, 117, 116, 99, 111, 109, 101})
			{
			out := t.InvokeStatic(m, 99, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 47, "Initialization", "contenderResult", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			if rt.AsInt(stack[1]) == 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			stack[1] = rt.IntValue(-2)
			block = 15
			continue
		case 14:
			stack[1] = rt.IntValue(73)
			block = 15
			continue
		case 15:
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 17
			continue
			}
			block = 16
			continue
		case 16:
			stack[0] = rt.IntValue(1)
			block = 18
			continue
		case 17:
			stack[0] = rt.IntValue(0)
			block = 18
			continue
		case 18:
			stack[1] = t.StringConstant([]uint16{98, 97, 100, 32, 105, 110, 105, 116, 105, 97, 108, 105, 122, 101, 114, 32, 99, 111, 109, 112, 101, 116, 105, 110, 103, 32, 111, 117, 116, 99, 111, 109, 101})
			{
			out := t.InvokeStatic(m, 126, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
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
	var locals [4]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.ClassConstant(m, 52, "Initialization$Ready")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			locals[1] = stack[0]
			{
			v, thrown := t.ClassConstant(m, 54, "Initialization$Broken")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			locals[2] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[2]
			if stack[0].Ref == stack[1].Ref {
			block = 3
			continue
			}
			block = 1
			continue
		case 1:
			{
			v, thrown := t.GetStatic(m, 56, "Initialization", "attempts", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			stack[0] = rt.IntValue(1)
			block = 4
			continue
		case 3:
			stack[0] = rt.IntValue(0)
			block = 4
			continue
		case 4:
			stack[1] = t.StringConstant([]uint16{99, 108, 97, 115, 115, 32, 108, 105, 116, 101, 114, 97, 108, 32, 116, 114, 105, 103, 103, 101, 114, 101, 100, 32, 105, 110, 105, 116, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110})
			{
			out := t.InvokeStatic(m, 24, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 28, 61, "Initialization", "trial", "(Z)V", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 28, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			out := t.InvokeStatic(m, 32, 61, "Initialization", "trial", "(Z)V", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 32, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 35, 65, "Initialization$Broken", "read", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 35, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			obj, thrown := t.NewObject(m, 69, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 39, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			stack[2] = t.StringConstant([]uint16{101, 114, 114, 111, 110, 101, 111, 117, 115, 32, 99, 108, 97, 115, 115, 32, 114, 97, 110, 32, 97, 103, 97, 105, 110})
			{
			out := t.InvokeSpecial(m, 45, 73, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 45, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			if hpc, ok := t.HandleThrow(m, 48, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 5:
			locals[3] = stack[0]
			{
			v, thrown := t.GetStatic(m, 56, "Initialization", "attempts", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 7
			continue
			}
			block = 6
			continue
		case 6:
			stack[0] = rt.IntValue(1)
			block = 8
			continue
		case 7:
			stack[0] = rt.IntValue(0)
			block = 8
			continue
		case 8:
			stack[1] = t.StringConstant([]uint16{105, 110, 105, 116, 105, 97, 108, 105, 122, 101, 114, 32, 100, 105, 100, 32, 110, 111, 116, 32, 114, 117, 110, 32, 101, 120, 97, 99, 116, 108, 121, 32, 111, 110, 99, 101, 32, 112, 101, 114, 32, 99, 108, 97, 115, 115})
			{
			out := t.InvokeStatic(m, 64, 41, "Check", "that", "(ZLjava/lang/String;)V", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 80, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 67, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{105, 110, 105, 116, 105, 97, 108, 105, 122, 97, 116, 105, 111, 110})
			{
			out := t.InvokeVirtual(m, 72, 88, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 72, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 49:
			block = 5
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
			obj, thrown := t.NewObject(m, 69, "java/lang/RuntimeException")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 4, 93, stack[1], "java/lang/RuntimeException", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutStatic(m, 94, "Initialization", "TOKEN", "Ljava/lang/RuntimeException;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

