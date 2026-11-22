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
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(41)
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
			stack[0] = rt.IntValue(7)
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
			stack[0] = rt.IntValue(1)
			return rt.OkValue(stack[0])
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
			stack[0] = locals[1]
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 2, 7, stack[0], "app/Box", "add", "(I)I", stack[1:2])
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

func cattyMethod5(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			thrown := t.PutField(m, 7, stack[0], "app/Box", "value", "I", stack[1])
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
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = stack[0]
			{
			v, thrown := t.GetField(m, 7, stack[1], "app/Box", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			thrown := t.PutField(m, 7, stack[0], "app/Box", "value", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = stack[0]
			{
			v, thrown := t.GetField(m, 13, stack[1], "app/Box", "callbacks", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			thrown := t.PutField(m, 13, stack[0], "app/Box", "callbacks", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 7, stack[0], "app/Box", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
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
			v, thrown := t.GetField(m, 7, stack[0], "app/Box", "value", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
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
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, 13, stack[0], "app/Box", "callbacks", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			thrown := t.PutField(m, 7, stack[0], "app/Box", "value", "I", stack[1])
			if thrown != nil {
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
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(100)
			{
			out := t.InvokeVirtual(m, 3, 7, stack[0], "app/Box", "add", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[1]
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
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
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
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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

func cattyMethod15(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [9]rt.Value
	var stack [6]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 2, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			out := t.InvokeStatic(m, 3, 7, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			{
			v, thrown := t.ClassConstant(m, 13, "app/Main")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			out := t.InvokeVirtual(m, 8, 15, stack[1], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 8, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			out := t.InvokeStatic(m, 11, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			locals[1] = stack[0]
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 84, 104, 114, 111, 119, 101, 114})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[1]
			{
			out := t.InvokeStatic(m, 19, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			locals[2] = stack[0]
			stack[0] = locals[2]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 16, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 27, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 30, 33, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 35, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 38, 37, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 38, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			locals[3] = stack[0]
			stack[0] = locals[3]
			{
			v, thrown := t.CheckCast(m, 43, "app/ThrowStep", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 43, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			locals[4] = stack[0]
			{
			obj, thrown := t.NewObject(m, 45, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 48, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[2] = t.StringConstant([]uint16{111, 114, 105, 103, 105, 110, 97, 108})
			{
			out := t.InvokeSpecial(m, 54, 49, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 54, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[5] = stack[0]
			{
			obj, thrown := t.NewObject(m, 52, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 59, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeSpecial(m, 64, 54, stack[1], "app/Box", "<init>", "(I)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[6] = stack[0]
			stack[0] = locals[4]
			stack[1] = locals[6]
			stack[2] = locals[5]
			{
			out := t.InvokeInterface(m, 75, 57, stack[0], "app/ThrowStep", "boom", "(Lapp/Box;Ljava/lang/Throwable;)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 75, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			{
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 85, 67, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 85, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 5
			continue
		case 1:
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 93, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = locals[7]
			stack[2] = locals[5]
			if stack[1].Ref != stack[2].Ref {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			stack[1] = rt.IntValue(1)
			block = 4
			continue
		case 3:
			stack[1] = rt.IntValue(0)
			block = 4
			continue
		case 4:
			{
			out := t.InvokeVirtual(m, 108, 74, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 108, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 111, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 116, 77, stack[1], "app/Box", "get", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 116, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			out := t.InvokeVirtual(m, 119, 67, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 119, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 5
			continue
		case 5:
			stack[0] = locals[2]
			stack[1] = t.StringConstant([]uint16{98, 111, 111, 109})
			stack[2] = rt.IntValue(2)
			{
			obj, thrown := t.NewReferenceArray(m, 16, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 126, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			v, thrown := t.ClassConstant(m, 52, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 131, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			if hpc, ok := t.HandleThrow(m, 133, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(1)
			{
			v, thrown := t.ClassConstant(m, 72, "java/lang/Throwable")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 136, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			if hpc, ok := t.HandleThrow(m, 138, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			out := t.InvokeVirtual(m, 139, 33, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 139, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			locals[7] = stack[0]
			stack[0] = locals[7]
			stack[1] = locals[3]
			stack[2] = rt.IntValue(2)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 148, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			stack[5] = locals[6]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 155, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(1)
			stack[5] = locals[5]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 160, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			out := t.InvokeVirtual(m, 161, 37, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 161, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			{
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 165, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 169, 67, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 169, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
		case 6:
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 177, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = locals[8]
			{
			out := t.InvokeVirtual(m, 182, 84, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 182, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = locals[5]
			if stack[1].Ref != stack[2].Ref {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			stack[1] = rt.IntValue(1)
			block = 9
			continue
		case 8:
			stack[1] = rt.IntValue(0)
			block = 9
			continue
		case 9:
			{
			out := t.InvokeVirtual(m, 195, 74, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 195, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 198, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 203, 77, stack[1], "app/Box", "get", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 203, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			out := t.InvokeVirtual(m, 206, 67, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 206, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			v, thrown := t.GetStatic(m, 61, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 209, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
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
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 214, 88, stack[1], "app/Box", "callbacks", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 214, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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
			out := t.InvokeVirtual(m, 217, 67, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 217, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 91:
			block = 1
			case 175:
			block = 6
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

func cattyMethod16(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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

func cattyMethod17(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 7, "app/Probe", "count", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			thrown := t.PutStatic(m, 7, "app/Probe", "count", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod18(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, 7, "app/Probe", "count", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod19(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [0]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			{
			thrown := t.PutStatic(m, 7, "app/Probe", "count", "I", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

