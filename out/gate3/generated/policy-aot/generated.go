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
	var stack [2]rt.Value
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
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
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
			out := t.InvokeStatic(m, 3, 7, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 3, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
			block = 1
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
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
			block = 1
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
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
			block = 1
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
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
			block = 1
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
			if hpc, ok := t.HandleThrow(m, 15, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
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
			stack[1] = t.StringConstant([]uint16{97, 99, 99, 101, 112, 116, 101, 100})
			{
			out := t.InvokeVirtual(m, 20, 33, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 20, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
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
			{
			v, thrown := t.GetStatic(m, 25, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 27, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
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
			stack[1] = locals[1]
			{
			out := t.InvokeVirtual(m, 31, 41, stack[1], "catty/runtime/PluginAdmissionException", "getCode", "()Ljava/lang/String;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 31, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
			block = 1
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
			out := t.InvokeVirtual(m, 34, 33, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 34, thrown); ok {
			stack = [2]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 26:
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
		case 2:
			return rt.Normal()
		}
	}
}

