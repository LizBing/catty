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
	var locals [6]rt.Value
	var stack [4]rt.Value
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
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeStatic(m, 13, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 13, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeStatic(m, 20, 15, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeStatic(m, 24, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 24, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 28, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[2]
			{
			out := t.InvokeStatic(m, 35, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 35, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{103, 111, 111, 100})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 41, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 44, 39, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 49, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 52, 43, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 52, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.CheckCast(m, 49, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 55, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 58, 51, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 58, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 61, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 61, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[3]
			{
			out := t.InvokeStatic(m, 71, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 71, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{103, 111, 111, 100})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 77, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 80, 39, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 80, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 85, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 88, 43, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 88, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.CheckCast(m, 49, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 91, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 94, 51, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 94, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 97, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 97, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			locals[4] = stack[0]
			block = 1
			continue
		case 1:
			stack[0] = locals[4]
			stack[1] = rt.IntValue(2)
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 8
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 113, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 113, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{102, 105, 101, 108, 100})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 119, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 122, 39, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 122, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 127, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 130, 43, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 130, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 134, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 138, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 138, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
		case 3:
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 146, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			{
			out := t.InvokeVirtual(m, 151, 65, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 151, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			n, thrown := t.InstanceOf(m, 69, "java/lang/NoSuchFieldError", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 154, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			stack[1] = rt.IntValue(9)
			block = 6
			continue
		case 6:
			{
			out := t.InvokeVirtual(m, 166, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 166, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
		case 7:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 1)
			block = 1
			continue
		case 8:
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 175, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[2]
			{
			out := t.InvokeStatic(m, 182, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 182, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{102, 105, 101, 108, 100})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 188, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 191, 39, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 191, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 196, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 199, 43, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 199, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.CheckCast(m, 49, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 202, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 205, 51, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 205, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 208, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 208, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			locals[4] = stack[0]
			block = 9
			continue
		case 9:
			stack[0] = locals[4]
			stack[1] = rt.IntValue(2)
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 16
			continue
			}
			block = 10
			continue
		case 10:
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 224, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 224, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{109, 101, 116, 104, 111, 100})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 230, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 233, 39, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 233, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 238, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 241, 43, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 241, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 245, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 249, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 249, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
		case 11:
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 257, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			{
			out := t.InvokeVirtual(m, 262, 65, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 262, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			n, thrown := t.InstanceOf(m, 73, "java/lang/NoSuchMethodError", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 265, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			block = 13
			continue
			}
			block = 12
			continue
		case 12:
			stack[1] = rt.IntValue(2)
			block = 14
			continue
		case 13:
			stack[1] = rt.IntValue(9)
			block = 14
			continue
		case 14:
			{
			out := t.InvokeVirtual(m, 277, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 277, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 1)
			block = 9
			continue
		case 16:
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 286, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[2]
			{
			out := t.InvokeStatic(m, 293, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 293, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{109, 101, 116, 104, 111, 100})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 299, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 302, 39, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 302, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 307, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 310, 43, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 310, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.CheckCast(m, 49, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 313, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 316, 51, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 316, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 319, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 319, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			locals[4] = stack[0]
			block = 17
			continue
		case 17:
			stack[0] = locals[4]
			stack[1] = rt.IntValue(2)
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 24
			continue
			}
			block = 18
			continue
		case 18:
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 335, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 335, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = t.StringConstant([]uint16{107, 108, 97, 115, 115})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 341, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 344, 39, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 344, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 349, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 352, 43, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 352, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			if hpc, ok := t.HandleThrow(m, 356, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 360, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 360, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 23
			continue
		case 19:
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 368, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			{
			out := t.InvokeVirtual(m, 373, 65, stack[1], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 373, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			n, thrown := t.InstanceOf(m, 77, "java/lang/NoClassDefFoundError", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 376, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			block = 21
			continue
			}
			block = 20
			continue
		case 20:
			stack[1] = rt.IntValue(3)
			block = 22
			continue
		case 21:
			stack[1] = rt.IntValue(9)
			block = 22
			continue
		case 22:
			{
			out := t.InvokeVirtual(m, 388, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 388, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 23
			continue
		case 23:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 1)
			block = 17
			continue
		case 24:
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 397, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{101, 120, 116, 46, 76, 105, 110, 107})
			stack[2] = rt.IntValue(0)
			stack[3] = locals[2]
			{
			out := t.InvokeStatic(m, 404, 33, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[1:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 404, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = t.StringConstant([]uint16{103, 111, 111, 100})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 410, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 413, 39, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 413, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 418, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = rt.RefValue(obj)
			}
			{
			out := t.InvokeVirtual(m, 421, 43, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 421, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			v, thrown := t.CheckCast(m, 49, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 424, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 427, 51, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 427, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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
			out := t.InvokeVirtual(m, 430, 55, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 430, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 144:
			block = 3
			case 255:
			block = 11
			case 366:
			block = 19
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

