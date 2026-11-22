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
	var locals [13]rt.Value
	var stack [7]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.ClassConstant(m, 7, "app/Main")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 2, 9, stack[0], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 8, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeStatic(m, 9, 15, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 15, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[0] = locals[2]
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 21, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 21, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 26, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[2] = locals[4]
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
			stack[1] = rt.IntValue(0)
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 40, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 40, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 87, 111, 114, 107, 101, 114})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 47, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 87, 111, 114, 107, 101, 114})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[4]
			{
			out := t.InvokeStatic(m, 57, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 62, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[5]
			stack[2] = locals[6]
			if stack[1].Ref != stack[2].Ref {
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
			out := t.InvokeVirtual(m, 77, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 77, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 80, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{97, 112, 112, 46, 66, 97, 115, 101})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[3]
			{
			out := t.InvokeStatic(m, 87, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 87, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.ClassConstant(m, 45, "app/Base")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 100, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 100, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 103, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{91, 73})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[3]
			{
			out := t.InvokeStatic(m, 110, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 110, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{91, 73})
			stack[3] = rt.IntValue(0)
			stack[4] = rt.RefValue(nil)
			{
			out := t.InvokeStatic(m, 117, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[2:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 117, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			block = 11
			continue
			}
			block = 10
			continue
		case 10:
			stack[1] = rt.IntValue(1)
			block = 12
			continue
		case 11:
			stack[1] = rt.IntValue(0)
			block = 12
			continue
		case 12:
			{
			out := t.InvokeVirtual(m, 128, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 128, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 131, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{91, 76, 101, 120, 116, 46, 87, 111, 114, 107, 101, 114, 59})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[3]
			{
			out := t.InvokeStatic(m, 138, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 138, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{91, 76, 101, 120, 116, 46, 87, 111, 114, 107, 101, 114, 59})
			stack[3] = rt.IntValue(0)
			stack[4] = locals[4]
			{
			out := t.InvokeStatic(m, 146, 39, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[2:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 146, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			stack[1] = rt.IntValue(1)
			block = 15
			continue
		case 14:
			stack[1] = rt.IntValue(0)
			block = 15
			continue
		case 15:
			{
			out := t.InvokeVirtual(m, 157, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 157, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[5]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 165, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 168, 53, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 168, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			if hpc, ok := t.HandleThrow(m, 173, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 176, 57, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 176, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.CheckCast(m, 45, "app/Base", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 179, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			locals[7] = stack[0]
			stack[0] = locals[6]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 189, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 192, 53, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 192, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			if hpc, ok := t.HandleThrow(m, 197, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 200, 57, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 200, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.CheckCast(m, 45, "app/Base", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 203, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 208, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			obj, thrown := t.NewObject(m, 63, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 213, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeSpecial(m, 218, 65, stack[3], "app/Box", "<init>", "(I)V", stack[4:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 218, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 221, 68, stack[1], "app/Base", "value", "(Lapp/Box;)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 221, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 224, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 224, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 227, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			obj, thrown := t.NewObject(m, 63, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 232, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeSpecial(m, 237, 65, stack[3], "app/Box", "<init>", "(I)V", stack[4:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 237, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 240, 68, stack[1], "app/Base", "value", "(Lapp/Box;)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 240, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 243, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 243, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 246, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			obj, thrown := t.NewObject(m, 63, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 251, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeSpecial(m, 256, 65, stack[3], "app/Box", "<init>", "(I)V", stack[4:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 256, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 259, 68, stack[1], "app/Base", "value", "(Lapp/Box;)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 259, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 262, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 262, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 265, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			n, thrown := t.InstanceOf(m, 74, "app/Step", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 270, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.IntValue(n)
			}
			{
			out := t.InvokeVirtual(m, 273, 31, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 273, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[5]
			stack[1] = t.StringConstant([]uint16{115, 101, 108, 102})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 281, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.ClassConstant(m, 2, "java/lang/Object")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 286, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			if hpc, ok := t.HandleThrow(m, 288, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 289, 53, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 289, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 294, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 301, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[6] = locals[7]
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 308, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 309, 57, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 309, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.CheckCast(m, 78, "java/lang/Boolean", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 312, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 315, 80, stack[1], "java/lang/Boolean", "booleanValue", "()Z", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 315, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			block = 17
			continue
			}
			block = 16
			continue
		case 16:
			stack[1] = rt.IntValue(1)
			block = 18
			continue
		case 17:
			stack[1] = rt.IntValue(0)
			block = 18
			continue
		case 18:
			{
			out := t.InvokeVirtual(m, 326, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 326, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 329, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 336, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[6] = locals[8]
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 343, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 344, 57, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 344, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.CheckCast(m, 78, "java/lang/Boolean", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 347, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 350, 80, stack[1], "java/lang/Boolean", "booleanValue", "()Z", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 350, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			block = 20
			continue
			}
			block = 19
			continue
		case 19:
			stack[1] = rt.IntValue(1)
			block = 21
			continue
		case 20:
			stack[1] = rt.IntValue(0)
			block = 21
			continue
		case 21:
			{
			out := t.InvokeVirtual(m, 361, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 361, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[5]
			stack[1] = t.StringConstant([]uint16{99, 97, 115, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 369, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.ClassConstant(m, 2, "java/lang/Object")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 374, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			if hpc, ok := t.HandleThrow(m, 376, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 377, 53, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 377, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[0] = locals[10]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 386, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[5] = locals[8]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 393, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 394, 57, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 394, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 398, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 402, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 402, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 26
			continue
		case 22:
			locals[11] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 410, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			{
			out := t.InvokeVirtual(m, 415, 88, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 415, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			n, thrown := t.InstanceOf(m, 92, "java/lang/ClassCastException", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 418, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			block = 24
			continue
			}
			block = 23
			continue
		case 23:
			stack[1] = rt.IntValue(3)
			block = 25
			continue
		case 24:
			stack[1] = rt.IntValue(9)
			block = 25
			continue
		case 25:
			{
			out := t.InvokeVirtual(m, 430, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 430, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 26
			continue
		case 26:
			stack[0] = locals[5]
			stack[1] = t.StringConstant([]uint16{115, 116, 111, 114, 101})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 438, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.ClassConstant(m, 2, "java/lang/Object")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 443, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			if hpc, ok := t.HandleThrow(m, 445, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 446, 53, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 446, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[0] = locals[11]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 455, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[5] = locals[8]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 462, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeVirtual(m, 463, 57, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 463, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 467, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			out := t.InvokeVirtual(m, 471, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 471, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 31
			continue
		case 27:
			locals[12] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 479, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			{
			out := t.InvokeVirtual(m, 484, 88, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 484, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			n, thrown := t.InstanceOf(m, 96, "java/lang/ArrayStoreException", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 487, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
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
			stack[1] = rt.IntValue(4)
			block = 30
			continue
		case 29:
			stack[1] = rt.IntValue(9)
			block = 30
			continue
		case 30:
			{
			out := t.InvokeVirtual(m, 499, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 499, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 408:
			block = 22
			case 477:
			block = 27
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 31
			continue
		case 31:
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

