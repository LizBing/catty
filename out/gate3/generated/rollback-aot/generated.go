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
	var locals [8]rt.Value
	var stack [6]rt.Value
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
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			{
			out := t.InvokeVirtual(m, 2, 9, stack[0], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[1] = stack[0]
			{
			obj, thrown := t.NewObject(m, 15, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 6, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeSpecial(m, 11, 17, stack[1], "app/Box", "<init>", "(I)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 17, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			{
			out := t.InvokeStatic(m, 18, 20, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 22, 26, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 22, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[3] = stack[0]
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 87, 111, 114, 107, 101, 114})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 30, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 36, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 39, 38, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 39, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 47, 42, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 47, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			v, thrown := t.CheckCast(m, 48, "app/Base", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[4] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 55, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[4]
			stack[2] = locals[2]
			{
			out := t.InvokeVirtual(m, 61, 56, stack[1], "app/Base", "value", "(Lapp/Box;)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 64, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 69, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			{
			out := t.InvokeStatic(m, 70, 20, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 70, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 74, 26, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 74, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 78, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = t.StringConstant([]uint16{117, 110, 101, 120, 112, 101, 99, 116, 101, 100})
			{
			out := t.InvokeVirtual(m, 83, 67, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 83, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 91, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = t.StringConstant([]uint16{67, 76, 65, 83, 83, 95, 70, 79, 82, 77, 65, 84})
			stack[2] = locals[5]
			{
			out := t.InvokeVirtual(m, 98, 74, stack[2], "catty/runtime/PluginAdmissionException", "getCode", "()Ljava/lang/String;", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 98, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 101, 78, stack[1], "java/lang/String", "equals", "(Ljava/lang/Object;)Z", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 101, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 104, 84, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 104, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 107, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[4]
			stack[2] = locals[2]
			{
			out := t.InvokeVirtual(m, 113, 56, stack[1], "app/Base", "value", "(Lapp/Box;)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 113, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 116, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 116, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 65, 108, 112, 104, 97})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[1]
			{
			out := t.InvokeStatic(m, 123, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 123, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 127, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 131, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 131, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 139, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 143, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 143, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 65, 108, 112, 104, 97})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 150, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 150, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 154, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 158, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 158, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 166, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 170, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 170, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			stack[1] = rt.IntValue(2)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 175, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			{
			out := t.InvokeStatic(m, 176, 20, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 176, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 180, 26, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 180, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[5] = stack[0]
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 65, 108, 112, 104, 97})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[5]
			{
			out := t.InvokeStatic(m, 190, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 190, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 196, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 199, 38, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 199, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 204, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 207, 42, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 207, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[6] = stack[0]
			stack[0] = locals[6]
			{
			out := t.InvokeVirtual(m, 214, 91, stack[0], "java/lang/Object", "getClass", "()Ljava/lang/Class;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 214, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 220, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			v, thrown := t.ClassConstant(m, 15, "app/Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 225, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			if hpc, ok := t.HandleThrow(m, 227, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			out := t.InvokeVirtual(m, 228, 38, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 228, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[6]
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 234, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			stack[5] = locals[2]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 240, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			out := t.InvokeVirtual(m, 241, 42, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 241, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 246, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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
			stack[1] = locals[7]
			{
			v, thrown := t.CheckCast(m, 96, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 251, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 254, 98, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 254, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
			block = 5
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
			out := t.InvokeVirtual(m, 257, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 257, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 89:
			block = 1
			case 137:
			block = 3
			case 164:
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

