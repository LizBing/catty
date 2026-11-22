package main

import "catty/rt"

func cattyMethod0(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			obj, thrown := t.NewObject("Crossing$Signal")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			stack[2] = stack[1]
			{
			out := t.InvokeSpecial(m, 9, stack[2], "Crossing$Signal", "<init>", "()V", stack[3:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(stack[0], "Crossing$Box", "signal", "LCrossing$Signal;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod1(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			v, thrown := t.GetField(stack[1], "Crossing$Box", "count", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			thrown := t.PutField(stack[0], "Crossing$Box", "count", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/RuntimeException", "<init>", "()V", stack[1:1])
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
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
	var stack [6]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			{
			obj, thrown := t.NewObject("Crossing$Box")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			{
			out := t.InvokeSpecial(m, 4, stack[1], "Crossing$Box", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[0] = stack[0]
			{
			obj, thrown := t.NewObject("Crossing$Worker")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 8, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			{
			out := t.InvokeSpecial(m, 12, stack[1], "Crossing$Worker", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 12, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[0]
			{
			out := t.InvokeInterface(m, 18, stack[0], "Crossing$Step", "apply", "(LCrossing$Box;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 18, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(-1)
			return rt.OkValue(stack[0])
		case 1:
			locals[2] = stack[0]
			stack[0] = locals[2]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(stack[1], "Crossing$Box", "signal", "LCrossing$Signal;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 28, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			if stack[0].Ref != stack[1].Ref {
			block = 3
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(stack[0], "Crossing$Box", "count", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 35, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = rt.IntValue(-2)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = t.ClassConstant("Crossing$Step")
			stack[1] = t.StringConstant([]uint16{97, 112, 112, 108, 121})
			stack[2] = rt.IntValue(1)
			{
			arr, thrown := t.NewArray(rt.AsInt(stack[2]), "java/lang/Class")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 50, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(arr)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = t.ClassConstant("Crossing$Box")
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			out := t.InvokeVirtual(m, 58, stack[0], "java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 58, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			locals[2] = stack[0]
			stack[0] = locals[2]
			stack[1] = locals[1]
			stack[2] = rt.IntValue(1)
			{
			arr, thrown := t.NewArray(rt.AsInt(stack[2]), "java/lang/Object")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 65, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
			block = 5
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = rt.RefValue(arr)
			}
			stack[3] = stack[2]
			stack[4] = rt.IntValue(0)
			stack[5] = locals[0]
			{
			thrown := t.ArrayStore(stack[3], rt.AsInt(stack[4]), stack[5])
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 71, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			out := t.InvokeVirtual(m, 72, stack[0], "java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 72, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			stack[0] = rt.IntValue(-3)
			return rt.OkValue(stack[0])
		case 5:
			locals[3] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 81, stack[0], "java/lang/reflect/InvocationTargetException", "getCause", "()Ljava/lang/Throwable;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 81, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(stack[1], "Crossing$Box", "signal", "LCrossing$Signal;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 85, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			if stack[0].Ref == stack[1].Ref {
			block = 7
			continue
			}
			block = 6
			continue
		case 6:
			stack[0] = rt.IntValue(-4)
			return rt.OkValue(stack[0])
		case 7:
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(stack[0], "Crossing$Box", "count", "I")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 95, thrown); ok {
			stack = [6]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 25:
			block = 1
			case 79:
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
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			out := t.InvokeSpecial(m, 1, stack[0], "java/lang/Object", "<init>", "()V", stack[1:1])
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
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic("java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			{
			out := t.InvokeStatic(m, 3, "Crossing", "probe", "()I", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 6, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

