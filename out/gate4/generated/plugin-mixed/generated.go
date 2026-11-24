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

func cattyMethod1(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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

func cattyMethod2(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [3]rt.Value
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
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			{
			out := t.InvokeStatic(m, 3, 7, "catty/runtime/Plugins", "readJar", "(Ljava/lang/String;)[B", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, 13, "app/Main", "jar", "[B", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.ClassConstant(m, 14, "app/Main")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			{
			out := t.InvokeVirtual(m, 11, 19, stack[0], "java/lang/Class", "getClassLoader", "()Ljava/lang/ClassLoader;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			thrown := t.PutStatic(m, 25, "app/Main", "parent", "Ljava/lang/ClassLoader;", stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, 29, "app/Main$Registrar")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 21, 31, stack[1], "app/Main$Registrar", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			{
			obj, thrown := t.NewObject(m, 29, "app/Main$Registrar")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 29, 31, stack[1], "app/Main$Registrar", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[2] = stack[0]
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 34, 32, stack[0], "app/Main$Registrar", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 38, 32, stack[0], "app/Main$Registrar", "start", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			out := t.InvokeVirtual(m, 42, 35, stack[0], "app/Main$Registrar", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			{
			out := t.InvokeVirtual(m, 46, 35, stack[0], "app/Main$Registrar", "join", "()V", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			{
			v, thrown := t.GetField(m, 38, stack[0], "app/Main$Registrar", "workerClass", "Ljava/lang/Class;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if stack[0].Ref == nil {
			block = 5
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[2]
			{
			v, thrown := t.GetField(m, 38, stack[0], "app/Main$Registrar", "workerClass", "Ljava/lang/Class;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			if stack[0].Ref == nil {
			block = 5
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = locals[1]
			{
			v, thrown := t.GetField(m, 38, stack[0], "app/Main$Registrar", "workerClass", "Ljava/lang/Class;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			v, thrown := t.GetField(m, 38, stack[1], "app/Main$Registrar", "workerClass", "Ljava/lang/Class;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			if stack[0].Ref == stack[1].Ref {
			block = 5
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = locals[1]
			{
			v, thrown := t.GetField(m, 42, stack[0], "app/Main$Registrar", "visits", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 5
			continue
			}
			block = 4
			continue
		case 4:
			stack[0] = locals[2]
			{
			v, thrown := t.GetField(m, 42, stack[0], "app/Main$Registrar", "visits", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			{
			obj, thrown := t.NewObject(m, 46, "java/lang/RuntimeException")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			stack[2] = t.StringConstant([]uint16{108, 111, 97, 100, 101, 114, 45, 108, 111, 99, 97, 108, 32, 116, 104, 114, 101, 97, 100, 32, 111, 114, 32, 115, 116, 97, 116, 105, 99, 32, 105, 100, 101, 110, 116, 105, 116, 121, 32, 119, 97, 115, 32, 108, 111, 115, 116})
			{
			out := t.InvokeSpecial(m, 96, 50, stack[1], "java/lang/RuntimeException", "<init>", "(Ljava/lang/String;)V", stack[2:3])
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
		case 6:
			{
			v, thrown := t.GetStatic(m, 53, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = t.StringConstant([]uint16{112, 108, 117, 103, 105, 110})
			{
			out := t.InvokeVirtual(m, 105, 61, stack[0], "java/io/PrintStream", "println", "(Ljava/lang/String;)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

