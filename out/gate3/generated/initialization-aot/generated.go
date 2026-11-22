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
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 2, 9, stack[0], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			{
			v, thrown := t.ArrayLoad(stack[0], rt.AsInt(stack[1]), 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 8, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeStatic(m, 9, 15, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 15, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			{
			out := t.InvokeStatic(m, 19, 25, "app/Probe", "reset", "()V", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 73, 110, 105, 116})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 26, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 26, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			locals[4] = stack[0]
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 70, 97, 99, 116, 111, 114, 121})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 35, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 35, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[1] = t.StringConstant([]uint16{109, 97, 107, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 41, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 44, 40, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 44, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			if hpc, ok := t.HandleThrow(m, 49, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 52, 44, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 52, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeStatic(m, 60, 56, "app/Probe", "get", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 63, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 63, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 73, 110, 105, 116})
			stack[1] = rt.IntValue(1)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 70, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 70, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 74, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeStatic(m, 77, 56, "app/Probe", "get", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 77, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 80, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 80, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 83, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[1] = locals[4]
			stack[2] = t.StringConstant([]uint16{118, 97, 108, 117, 101})
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 91, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
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
			out := t.InvokeVirtual(m, 94, 40, stack[1], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 94, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 99, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
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
			out := t.InvokeVirtual(m, 102, 44, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 102, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.CheckCast(m, 68, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 105, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 108, 70, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 108, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 111, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 111, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 116, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 116, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			locals[6] = stack[0]
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 73, 110, 105, 116})
			stack[1] = rt.IntValue(1)
			stack[2] = locals[6]
			{
			out := t.InvokeStatic(m, 126, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 126, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 130, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeStatic(m, 133, 56, "app/Probe", "get", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 133, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 136, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 136, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			obj, thrown := t.NewObject(m, 73, "java/lang/RuntimeException")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 139, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[2] = t.StringConstant([]uint16{115, 117, 112, 112, 108, 105, 101, 100})
			{
			out := t.InvokeSpecial(m, 145, 77, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 145, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			thrown := t.PutStatic(m, 80, "app/Probe", "supplied", "Ljava/lang/Throwable;", stack[0])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 148, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 153, 21, "catty/runtime/Plugins", "register", "([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 153, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 66, 97, 100})
			stack[1] = rt.IntValue(1)
			stack[2] = locals[7]
			{
			out := t.InvokeStatic(m, 163, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 163, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 167, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 171, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 171, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 179, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 184, 87, stack[1], "java/lang/ExceptionInInitializerError", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 184, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 80, "app/Probe", "supplied", "Ljava/lang/Throwable;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 187, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
			block = 6
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
			out := t.InvokeVirtual(m, 198, 91, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 198, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 66, 97, 100})
			stack[1] = rt.IntValue(1)
			stack[2] = locals[7]
			{
			out := t.InvokeStatic(m, 206, 32, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 206, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 210, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			out := t.InvokeVirtual(m, 214, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 214, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, 50, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 222, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeVirtual(m, 226, 60, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 226, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 177:
			block = 1
			case 220:
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
		case 7:
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

