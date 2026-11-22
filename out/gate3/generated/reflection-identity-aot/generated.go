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
	var locals [22]rt.Value
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
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 2, 9, stack[0], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 2, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeStatic(m, 9, 15, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 9, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeStatic(m, 20, 15, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 32, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 32, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 87, 111, 114, 107, 101, 114})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 41, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 41, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 84, 111, 107, 101, 110})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 50, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 84, 111, 107, 101, 110})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 59, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 59, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = locals[6]
			stack[2] = locals[7]
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
			out := t.InvokeVirtual(m, 79, 39, stack[0], "java/io/PrintStream", "println", "(Z)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 79, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[6]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 87, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 90, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[8] = stack[0]
			stack[0] = locals[8]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 99, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 102, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 102, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = locals[7]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 112, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 115, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 115, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 124, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 127, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 127, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = locals[4]
			stack[1] = t.StringConstant([]uint16{97, 99, 99, 101, 112, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 137, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[4] = rt.IntValue(0)
			stack[5] = locals[6]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 144, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 145, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 145, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[12] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 150, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = locals[12]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 157, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[6] = locals[9]
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 164, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 165, 51, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 165, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.CheckCast(m, 59, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 168, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 171, 61, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 171, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 174, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 174, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			stack[1] = t.StringConstant([]uint16{97, 99, 99, 101, 112, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 182, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[4] = rt.IntValue(0)
			stack[5] = locals[7]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 189, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 190, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 190, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 194, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 198, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 198, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
		case 4:
			locals[13] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 206, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 210, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 210, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = locals[12]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 217, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[4] = rt.IntValue(0)
			stack[5] = locals[11]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 224, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 225, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 225, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 229, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 233, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 233, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			locals[13] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 241, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 245, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 245, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[0] = t.StringConstant([]uint16{91, 76, 101, 120, 116, 46, 84, 111, 107, 101, 110, 59})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 252, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 252, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[13] = stack[0]
			stack[0] = t.StringConstant([]uint16{91, 76, 101, 120, 116, 46, 84, 111, 107, 101, 110, 59})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[3]
			{
			out := t.InvokeStatic(m, 261, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 261, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[14] = stack[0]
			stack[0] = locals[4]
			stack[1] = t.StringConstant([]uint16{97, 99, 99, 101, 112, 116})
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 271, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[4] = rt.IntValue(0)
			stack[5] = locals[13]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 278, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 279, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 279, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[15] = stack[0]
			stack[0] = locals[4]
			stack[1] = t.StringConstant([]uint16{116, 111, 107, 101, 110, 115})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 289, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 292, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 292, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			if hpc, ok := t.HandleThrow(m, 297, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 300, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 300, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[16] = stack[0]
			stack[0] = locals[5]
			stack[1] = t.StringConstant([]uint16{116, 111, 107, 101, 110, 115})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 310, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 313, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 313, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			if hpc, ok := t.HandleThrow(m, 318, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 321, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 321, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[17] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 326, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = locals[15]
			stack[2] = rt.RefValue(nil)
			stack[3] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[3]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 333, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			stack[6] = locals[16]
			{
			thrown := t.ArrayStore(stack[4], rt.AsInt(stack[5]), stack[6], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 340, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 341, 51, stack[1], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[2:4])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 341, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.CheckCast(m, 59, "java/lang/Integer", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 344, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 347, 61, stack[1], "java/lang/Integer", "intValue", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 347, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 350, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 350, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[15]
			stack[1] = rt.RefValue(nil)
			stack[2] = rt.IntValue(1)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 357, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[4] = rt.IntValue(0)
			stack[5] = locals[17]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5], 'L')
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 364, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			out := t.InvokeVirtual(m, 365, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 365, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 369, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeVirtual(m, 373, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 373, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			locals[18] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 381, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = rt.IntValue(3)
			{
			out := t.InvokeVirtual(m, 385, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 385, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
		case 9:
			stack[0] = t.StringConstant([]uint16{101, 120, 116, 46, 67, 111, 118})
			stack[1] = rt.IntValue(0)
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 392, 27, "java/lang/Class", "forName", "(Ljava/lang/String;ZLjava/lang/ClassLoader;)Ljava/lang/Class;", stack[0:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 392, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[18] = stack[0]
			stack[0] = locals[18]
			stack[1] = t.StringConstant([]uint16{99, 114, 101, 97, 116, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 402, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 405, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 405, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			if hpc, ok := t.HandleThrow(m, 410, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 413, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 413, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[19] = stack[0]
			stack[0] = locals[18]
			stack[1] = t.StringConstant([]uint16{109, 97, 107, 101})
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 10, "java/lang/Class", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 423, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 426, 47, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 426, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[20] = stack[0]
			stack[0] = locals[20]
			stack[1] = locals[19]
			stack[2] = rt.IntValue(0)
			{
			obj, thrown := t.NewReferenceArray(m, 2, "java/lang/Object", rt.AsInt(stack[2]))
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 436, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			{
			out := t.InvokeVirtual(m, 439, 51, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 439, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[21] = stack[0]
			{
			v, thrown := t.GetStatic(m, 33, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 444, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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
			stack[1] = locals[21]
			{
			v, thrown := t.CheckCast(m, 80, "java/lang/String", stack[1])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 449, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 452, 82, stack[1], "java/lang/String", "length", "()I", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 452, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
			block = 8
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
			out := t.InvokeVirtual(m, 455, 65, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 455, thrown); ok {
			stack = [7]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 204:
			block = 4
			case 239:
			block = 6
			case 379:
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

