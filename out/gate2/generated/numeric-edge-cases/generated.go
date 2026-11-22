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

func cattyMethod1(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			out := t.InvokeSpecial(m, 2, stack[0], "java/lang/RuntimeException", "<init>", "(Ljava/lang/Throwable;)V", stack[1:2])
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

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [21]rt.Value
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
			stack[1] = rt.IntValue(0)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(31)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(32)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(33)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[6] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(-1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) << rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[7] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[8] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[9] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(31)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[10] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(32)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[11] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(33)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[12] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(-1)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) >> rt.IntShiftDistance(rt.AsInt(stack[1])))
			locals[13] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[14] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[15] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(31)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[16] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(32)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[17] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(33)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[18] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(-1)
			stack[0] = rt.IntValue(int32(uint32(rt.AsInt(stack[0])) >> rt.IntShiftDistance(rt.AsInt(stack[1]))))
			locals[19] = stack[0]
			stack[0] = locals[2]
			stack[1] = locals[3]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[4]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[5]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[6]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[8]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[9]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[10]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[11]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[12]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[13]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[14]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[15]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[16]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[17]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[18]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[19]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[20] = stack[0]
			stack[0] = locals[20]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[20] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[20]
			{
			out := t.InvokeVirtual(m, 161, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[20]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [36]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(63)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[7] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(64)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[9] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(65)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) << rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[11] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) >> rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[13] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) >> rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[15] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(63)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) >> rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[17] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(64)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) >> rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[19] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(65)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) >> rt.LongShiftDistance(rt.AsInt(stack[1])))
			locals[21] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(0)
			stack[0] = rt.LongValue(int64(uint64(rt.AsLong(stack[0])) >> rt.LongShiftDistance(rt.AsInt(stack[1]))))
			locals[23] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			stack[0] = rt.LongValue(int64(uint64(rt.AsLong(stack[0])) >> rt.LongShiftDistance(rt.AsInt(stack[1]))))
			locals[25] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(63)
			stack[0] = rt.LongValue(int64(uint64(rt.AsLong(stack[0])) >> rt.LongShiftDistance(rt.AsInt(stack[1]))))
			locals[27] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(64)
			stack[0] = rt.LongValue(int64(uint64(rt.AsLong(stack[0])) >> rt.LongShiftDistance(rt.AsInt(stack[1]))))
			locals[29] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(65)
			stack[0] = rt.LongValue(int64(uint64(rt.AsLong(stack[0])) >> rt.LongShiftDistance(rt.AsInt(stack[1]))))
			locals[31] = stack[0]
			stack[0] = locals[3]
			stack[1] = locals[5]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[7]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[9]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[11]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[13]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[15]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[17]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[19]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[21]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[23]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[25]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[27]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[29]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			stack[1] = locals[31]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			locals[33] = stack[0]
			stack[0] = locals[33]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			stack[1] = locals[33]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[35] = stack[0]
			stack[0] = locals[35]
			stack[1] = locals[0]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 147, "NumericEdgeCases", "intShiftCheck", "(II)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[35] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[35]
			{
			out := t.InvokeVirtual(m, 158, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[35]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod5(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [42]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0xffffffffffffffff)
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[3]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			locals[5] = stack[0]
			stack[0] = rt.LongBitsValue(0x0102030405060708)
			locals[7] = stack[0]
			stack[0] = locals[5]
			stack[1] = locals[7]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[9] = stack[0]
			stack[0] = locals[5]
			stack[1] = locals[7]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) - rt.AsLong(stack[1]))
			locals[11] = stack[0]
			stack[0] = locals[5]
			stack[1] = locals[7]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) * rt.AsLong(stack[1]))
			locals[13] = stack[0]
			stack[0] = locals[7]
			stack[1] = rt.LongValue(1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) | rt.AsLong(stack[1]))
			stack[1] = locals[2]
			stack[2] = rt.IntValue(1)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) & rt.AsInt(stack[2]))
			stack[1] = rt.LongValue(int64(rt.AsInt(stack[1])))
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[15] = stack[0]
			stack[0] = locals[9]
			stack[1] = locals[15]
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) / d)
			}
			locals[17] = stack[0]
			stack[0] = locals[11]
			stack[0] = rt.LongValue(-rt.AsLong(stack[0]))
			locals[19] = stack[0]
			stack[0] = locals[9]
			stack[1] = rt.LongBitsValue(0xff00ff00ff00ff00)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) & rt.AsLong(stack[1]))
			locals[21] = stack[0]
			stack[0] = locals[11]
			stack[1] = rt.LongBitsValue(0x00ff00ff00ff00ff)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) | rt.AsLong(stack[1]))
			locals[23] = stack[0]
			stack[0] = locals[13]
			stack[1] = locals[19]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			locals[25] = stack[0]
			stack[0] = rt.LongBitsValue(0x8000000000000000)
			locals[27] = stack[0]
			stack[0] = rt.LongBitsValue(0xffffffffffffffff)
			locals[29] = stack[0]
			stack[0] = locals[27]
			stack[1] = locals[29]
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) / d)
			}
			locals[31] = stack[0]
			stack[0] = rt.LongBitsValue(0x7fffffffffffffff)
			locals[33] = stack[0]
			stack[0] = rt.LongValue(1)
			locals[35] = stack[0]
			stack[0] = locals[33]
			stack[1] = locals[35]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[37] = stack[0]
			stack[0] = locals[27]
			stack[1] = locals[35]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) - rt.AsLong(stack[1]))
			locals[39] = stack[0]
			stack[0] = locals[9]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[9]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[11]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[11]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[13]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[13]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[17]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[17]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[21]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[23]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[25]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[31]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[31]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[37]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			stack[0] = locals[41]
			stack[1] = locals[39]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[41] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[9]
			{
			out := t.InvokeVirtual(m, 257, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[13]
			{
			out := t.InvokeVirtual(m, 265, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[31]
			{
			out := t.InvokeVirtual(m, 273, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[37]
			{
			out := t.InvokeVirtual(m, 281, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[41]
			stack[1] = locals[5]
			stack[2] = locals[2]
			{
			out := t.InvokeStatic(m, 289, "NumericEdgeCases", "longShiftCheck", "(JI)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[41] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[41]
			{
			out := t.InvokeVirtual(m, 300, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[41]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [15]rt.Value
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
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) - rt.AsInt(stack[1]))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) * rt.AsInt(stack[1]))
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) & rt.AsInt(stack[1]))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) | rt.AsInt(stack[1]))
			locals[6] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[7] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.IntValue(-rt.AsInt(stack[0]))
			locals[8] = stack[0]
			stack[0] = rt.IntValue(-2147483648)
			locals[9] = stack[0]
			stack[0] = rt.IntValue(-1)
			locals[10] = stack[0]
			stack[0] = locals[9]
			stack[1] = locals[10]
			{
			d := rt.AsInt(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 43, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) / d)
			}
			locals[11] = stack[0]
			stack[0] = rt.IntValue(0)
			locals[12] = stack[0]
			stack[0] = rt.IntValue(0)
			locals[13] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[12]
			{
			d := rt.AsInt(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 55, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) / d)
			}
			locals[13] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 58, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			stack[1] = rt.IntValue(-7001)
			{
			out := t.InvokeVirtual(m, 64, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 64, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			locals[14] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 72, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			stack[1] = rt.IntValue(7001)
			{
			out := t.InvokeVirtual(m, 78, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 78, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			stack[0] = locals[2]
			stack[1] = locals[3]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[4]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[5]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[6]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[8]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[11]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[13]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[14] = stack[0]
			stack[0] = locals[14]
			stack[1] = locals[0]
			{
			out := t.InvokeStatic(m, 110, "NumericEdgeCases", "leaf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 110, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			stack[2] = locals[1]
			{
			out := t.InvokeStatic(m, 114, "NumericEdgeCases", "leaf", "(I)I", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 114, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[14] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 121, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
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
			stack[1] = locals[14]
			{
			out := t.InvokeVirtual(m, 126, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 126, thrown); ok {
			stack = [3]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 70:
			block = 1
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[14]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			locals[1] = stack[0]
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 7)
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + -3)
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 100)
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + -250)
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 1)
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + -1)
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			{
			out := t.InvokeVirtual(m, 27, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [20]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[0] = stack[0]
			stack[0] = rt.IntValue(1)
			locals[1] = stack[0]
			stack[0] = rt.IntValue(2)
			locals[2] = stack[0]
			stack[0] = rt.IntValue(3)
			locals[3] = stack[0]
			stack[0] = rt.IntValue(4)
			locals[4] = stack[0]
			stack[0] = rt.IntValue(5)
			locals[5] = stack[0]
			stack[0] = rt.IntValue(-1)
			locals[6] = stack[0]
			stack[0] = rt.IntValue(100)
			locals[7] = stack[0]
			stack[0] = rt.IntValue(30000)
			locals[8] = stack[0]
			stack[0] = rt.IntValue(100000)
			locals[9] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[3]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[4]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[5]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[6]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[8]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[9]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[10] = stack[0]
			stack[0] = rt.LongBitsValue(0x1122334455667788)
			locals[11] = stack[0]
			stack[0] = rt.LongBitsValue(0x8877665544332211)
			locals[13] = stack[0]
			stack[0] = locals[11]
			stack[1] = locals[13]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) ^ rt.AsLong(stack[1]))
			locals[15] = stack[0]
			stack[0] = rt.LongValue(0)
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186a9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186aa)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ab)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ac)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ad)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ae)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186af)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186b9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ba)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186bb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186bc)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186bd)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186be)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186bf)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186c9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ca)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186cb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186cc)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186cd)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ce)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186cf)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186d9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186da)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186db)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186dc)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186dd)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186de)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186df)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186e9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ea)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186eb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ec)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ed)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ee)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ef)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186f9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186fa)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186fb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186fc)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186fd)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186fe)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000186ff)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018700)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018701)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018702)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018703)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018704)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018705)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018706)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018707)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018708)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018709)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001870f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018710)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018711)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018712)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018713)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018714)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018715)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018716)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018717)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018718)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018719)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001871f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018720)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018721)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018722)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018723)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018724)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018725)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018726)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018727)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018728)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018729)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001872f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018730)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018731)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018732)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018733)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018734)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018735)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018736)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018737)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018738)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018739)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001873f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018740)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018741)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018742)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018743)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018744)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018745)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018746)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018747)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018748)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018749)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001874f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018750)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018751)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018752)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018753)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018754)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018755)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018756)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018757)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018758)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018759)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001875f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018760)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018761)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018762)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018763)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018764)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018765)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018766)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018767)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018768)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018769)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001876f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018770)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018771)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018772)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018773)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018774)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018775)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018776)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018777)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018778)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018779)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001877f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018780)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018781)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018782)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018783)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018784)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018785)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018786)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018787)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018788)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018789)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001878f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018790)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018791)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018792)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018793)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018794)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018795)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018796)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018797)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018798)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x0000000000018799)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879a)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879c)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879e)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x000000000001879f)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187a9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187aa)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ab)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ac)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ad)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ae)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187af)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187b9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ba)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187bb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187bc)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187bd)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187be)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187bf)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c0)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c1)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c2)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c3)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c4)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c5)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c6)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c7)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c8)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187c9)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187ca)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[17]
			stack[1] = rt.LongBitsValue(0x00000000000187cb)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[17] = stack[0]
			stack[0] = locals[10]
			locals[19] = stack[0]
			stack[0] = locals[19]
			stack[1] = locals[15]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[19] = stack[0]
			stack[0] = locals[19]
			stack[1] = locals[15]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[19] = stack[0]
			stack[0] = locals[19]
			stack[1] = locals[17]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[19] = stack[0]
			stack[0] = locals[19]
			stack[1] = locals[17]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[19] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[19]
			{
			out := t.InvokeVirtual(m, 2524, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[19]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [17]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[0] = rt.LongValue(int64(rt.AsInt(stack[0])))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.FloatValue(float32(rt.AsInt(stack[0])))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.DoubleValue(float64(rt.AsInt(stack[0])))
			locals[6] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			locals[8] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.FloatValue(float32(rt.AsLong(stack[0])))
			locals[9] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.DoubleValue(float64(rt.AsLong(stack[0])))
			locals[10] = stack[0]
			stack[0] = locals[3]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[3]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[8]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[5]
			{
			out := t.InvokeStatic(m, 48, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[6]
			{
			out := t.InvokeStatic(m, 56, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[13] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[13]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[13]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[9]
			{
			out := t.InvokeStatic(m, 84, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[10]
			{
			out := t.InvokeStatic(m, 92, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[15] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[15]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			stack[0] = locals[12]
			stack[1] = locals[15]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[12] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[12]
			{
			out := t.InvokeVirtual(m, 121, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[12]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod10(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [9]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[0] = rt.DoubleValue(float64(rt.AsFloat(stack[0])))
			locals[3] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.FloatValue(float32(rt.AsDouble(stack[0])))
			locals[5] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeStatic(m, 8, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[6] = stack[0]
			stack[0] = locals[5]
			{
			out := t.InvokeStatic(m, 15, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 25, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[8]
			{
			out := t.InvokeVirtual(m, 33, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[6]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			stack[1] = locals[6]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[8]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod11(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [9]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[0] = rt.IntValue(rt.F2I(rt.AsFloat(stack[0])))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.LongValue(rt.F2L(rt.AsFloat(stack[0])))
			locals[4] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.IntValue(rt.D2I(rt.AsDouble(stack[0])))
			locals[6] = stack[0]
			stack[0] = locals[1]
			stack[0] = rt.LongValue(rt.D2L(rt.AsDouble(stack[0])))
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			{
			out := t.InvokeVirtual(m, 19, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeVirtual(m, 27, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 35, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			out := t.InvokeVirtual(m, 43, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			stack[1] = locals[4]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[6]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod12(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [5]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[0] = rt.IntValue(rt.I2B(rt.AsInt(stack[0])))
			locals[1] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.IntValue(rt.I2C(rt.AsInt(stack[0])))
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.IntValue(rt.I2S(rt.AsInt(stack[0])))
			locals[3] = stack[0]
			stack[0] = locals[1]
			stack[1] = rt.IntValue(31)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) * rt.AsInt(stack[1]))
			stack[1] = locals[2]
			stack[2] = rt.IntValue(17)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) * rt.AsInt(stack[2]))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[3]
			stack[2] = rt.IntValue(7)
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) * rt.AsInt(stack[2]))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[4] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeVirtual(m, 30, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod13(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.FloatValue(0)
			locals[1] = stack[0]
			stack[0] = rt.FloatValue(1)
			locals[2] = stack[0]
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) / rt.AsFloat(stack[1]))
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = locals[2]
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) / rt.AsFloat(stack[1]))
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = locals[2]
			stack[0] = rt.FloatValue(-rt.AsFloat(stack[0]))
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) / rt.AsFloat(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [5]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.DoubleValue(0)
			locals[1] = stack[0]
			stack[0] = rt.DoubleValue(1)
			locals[3] = stack[0]
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) / rt.AsDouble(stack[1]))
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = locals[3]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) / rt.AsDouble(stack[1]))
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = locals[3]
			stack[0] = rt.DoubleValue(-rt.AsDouble(stack[0]))
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) / rt.AsDouble(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod15(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [8]rt.Value
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
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) + rt.AsFloat(stack[1]))
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) - rt.AsFloat(stack[1]))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) * rt.AsFloat(stack[1]))
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) / rt.AsFloat(stack[1]))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.FloatValue(-rt.AsFloat(stack[0]))
			locals[6] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeStatic(m, 26, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 29, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			{
			out := t.InvokeStatic(m, 36, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 39, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeStatic(m, 47, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 50, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[5]
			{
			out := t.InvokeStatic(m, 58, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 61, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeStatic(m, 69, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 72, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 77, "NumericEdgeCases", "compareFloat", "(FF)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			out := t.InvokeVirtual(m, 87, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			{
			out := t.InvokeStatic(m, 91, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = locals[7]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod16(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), 1))
			if rt.AsInt(stack[0]) >= 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), 1))
			if rt.AsInt(stack[0]) > 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) <= 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 4)
			block = 6
			continue
		case 6:
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) < 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 8)
			block = 8
			continue
		case 8:
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 16)
			block = 10
			continue
		case 10:
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 32)
			block = 12
			continue
		case 12:
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 64)
			block = 14
			continue
		case 14:
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 128)
			block = 16
			continue
		case 16:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 81, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod17(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [15]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) + rt.AsDouble(stack[1]))
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) - rt.AsDouble(stack[1]))
			locals[6] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) * rt.AsDouble(stack[1]))
			locals[8] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) / rt.AsDouble(stack[1]))
			locals[10] = stack[0]
			stack[0] = locals[0]
			stack[0] = rt.DoubleValue(-rt.AsDouble(stack[0]))
			locals[12] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeStatic(m, 29, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 32, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeStatic(m, 40, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 43, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[8]
			{
			out := t.InvokeStatic(m, 51, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 54, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[10]
			{
			out := t.InvokeStatic(m, 62, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 65, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[12]
			{
			out := t.InvokeStatic(m, 73, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 76, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[2]
			{
			out := t.InvokeStatic(m, 81, "NumericEdgeCases", "compareDouble", "(DD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[14] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[14]
			{
			out := t.InvokeVirtual(m, 91, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			{
			out := t.InvokeStatic(m, 96, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			stack[1] = locals[14]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod18(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [5]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), 1))
			if rt.AsInt(stack[0]) >= 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), 1))
			if rt.AsInt(stack[0]) > 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) <= 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 4)
			block = 6
			continue
		case 6:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) < 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 8)
			block = 8
			continue
		case 8:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 16)
			block = 10
			continue
		case 10:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 32)
			block = 12
			continue
		case 12:
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 64)
			block = 14
			continue
		case 14:
			stack[0] = locals[2]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) != 0 {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 128)
			block = 16
			continue
		case 16:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeVirtual(m, 83, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod19(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [16]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) + rt.AsFloat(stack[1]))
			locals[3] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) - rt.AsFloat(stack[1]))
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) * rt.AsFloat(stack[1]))
			locals[5] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[0]
			stack[0] = rt.FloatValue(rt.AsFloat(stack[0]) / rt.AsFloat(stack[1]))
			locals[6] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) + rt.AsDouble(stack[1]))
			locals[7] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) - rt.AsDouble(stack[1]))
			locals[9] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) * rt.AsDouble(stack[1]))
			locals[11] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[1]
			stack[0] = rt.DoubleValue(rt.AsDouble(stack[0]) / rt.AsDouble(stack[1]))
			locals[13] = stack[0]
			stack[0] = rt.IntValue(0)
			locals[15] = stack[0]
			stack[0] = locals[3]
			stack[1] = locals[3]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[4]
			stack[1] = locals[4]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[5]
			stack[1] = locals[5]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 4)
			block = 6
			continue
		case 6:
			stack[0] = locals[6]
			stack[1] = locals[6]
			stack[0] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[0]), rt.AsFloat(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 8)
			block = 8
			continue
		case 8:
			stack[0] = locals[7]
			stack[1] = locals[7]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 16)
			block = 10
			continue
		case 10:
			stack[0] = locals[9]
			stack[1] = locals[9]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 32)
			block = 12
			continue
		case 12:
			stack[0] = locals[11]
			stack[1] = locals[11]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 64)
			block = 14
			continue
		case 14:
			stack[0] = locals[13]
			stack[1] = locals[13]
			stack[0] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[0]), rt.AsDouble(stack[1]), -1))
			if rt.AsInt(stack[0]) == 0 {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			locals[15] = rt.IntValue(rt.AsInt(locals[15]) + 128)
			block = 16
			continue
		case 16:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[15]
			{
			out := t.InvokeVirtual(m, 136, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[15]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod20(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [9]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeStatic(m, 1, "java/lang/Float", "intBitsToFloat", "(I)F", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[3] = stack[0]
			stack[0] = locals[1]
			{
			out := t.InvokeStatic(m, 6, "java/lang/Double", "longBitsToDouble", "(J)D", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[4] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeStatic(m, 12, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[6] = stack[0]
			stack[0] = locals[4]
			{
			out := t.InvokeStatic(m, 19, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 29, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			out := t.InvokeVirtual(m, 37, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[6]
			stack[1] = locals[7]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = locals[7]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod21(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [5]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[4] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) >= 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) != 0 {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[0]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.Lcmp(rt.AsLong(stack[0]), rt.AsLong(stack[1])))
			if rt.AsInt(stack[0]) <= 0 {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[4] = rt.IntValue(rt.AsInt(locals[4]) + 4)
			block = 6
			continue
		case 6:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeVirtual(m, 35, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[4]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod22(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) < rt.AsInt(stack[1]) {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 4)
			block = 6
			continue
		case 6:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) > rt.AsInt(stack[1]) {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 8)
			block = 8
			continue
		case 8:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 16)
			block = 10
			continue
		case 10:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if rt.AsInt(stack[0]) == rt.AsInt(stack[1]) {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 32)
			block = 12
			continue
		case 12:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) == 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 64)
			block = 14
			continue
		case 14:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) != 0 {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 128)
			block = 16
			continue
		case 16:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) >= 0 {
			block = 18
			continue
			}
			block = 17
			continue
		case 17:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 256)
			block = 18
			continue
		case 18:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) < 0 {
			block = 20
			continue
			}
			block = 19
			continue
		case 19:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 512)
			block = 20
			continue
		case 20:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) <= 0 {
			block = 22
			continue
			}
			block = 21
			continue
		case 21:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1024)
			block = 22
			continue
		case 22:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) > 0 {
			block = 24
			continue
			}
			block = 23
			continue
		case 23:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 2048)
			block = 24
			continue
		case 24:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 111, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod23(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[1] = args[1]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			if stack[0].Ref != stack[1].Ref {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1)
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = locals[1]
			if stack[0].Ref == stack[1].Ref {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 2)
			block = 4
			continue
		case 4:
			stack[0] = locals[0]
			if stack[0].Ref != nil {
			block = 6
			continue
			}
			block = 5
			continue
		case 5:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 4)
			block = 6
			continue
		case 6:
			stack[0] = locals[0]
			if stack[0].Ref == nil {
			block = 8
			continue
			}
			block = 7
			continue
		case 7:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 8)
			block = 8
			continue
		case 8:
			stack[0] = locals[1]
			if stack[0].Ref != nil {
			block = 10
			continue
			}
			block = 9
			continue
		case 9:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 16)
			block = 10
			continue
		case 10:
			stack[0] = locals[1]
			if stack[0].Ref == nil {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 32)
			block = 12
			continue
		case 12:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 50, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod24(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [4]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			out := t.InvokeStatic(m, 1, "NumericEdgeCases", "leaf", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[1] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(17)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			out := t.InvokeStatic(m, 9, "NumericEdgeCases", "leaf", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[2] = stack[0]
			stack[0] = locals[1]
			stack[1] = rt.IntValue(31)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) * rt.AsInt(stack[1]))
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[3] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			{
			out := t.InvokeVirtual(m, 24, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[3]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod25(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
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
			obj, thrown := t.NewObject(m, "NumericEdgeCases$Marker")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 0, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			out := t.InvokeSpecial(m, 4, stack[1], "NumericEdgeCases$Marker", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[1] = stack[0]
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			{
			obj, thrown := t.NewObject(m, "NumericEdgeCases$Marker")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 14, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			stack[2] = locals[1]
			{
			out := t.InvokeSpecial(m, 19, stack[1], "NumericEdgeCases$Marker", "<init>", "(Ljava/lang/Throwable;)V", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 19, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			if hpc, ok := t.HandleThrow(m, 22, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 2:
			stack[0] = rt.IntValue(5)
			locals[2] = stack[0]
			block = 10
			continue
		case 3:
			locals[3] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 30, stack[0], "NumericEdgeCases$Marker", "getCause", "()Ljava/lang/Throwable;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 30, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			if stack[0].Ref != stack[1].Ref {
			block = 5
			continue
			}
			block = 4
			continue
		case 4:
			stack[0] = rt.IntValue(1)
			block = 6
			continue
		case 5:
			stack[0] = rt.IntValue(0)
			block = 6
			continue
		case 6:
			locals[4] = stack[0]
			stack[0] = locals[3]
			{
			out := t.InvokeVirtual(m, 45, stack[0], "NumericEdgeCases$Marker", "getCause", "()Ljava/lang/Throwable;", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 45, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = locals[3]
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
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 59, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1000)
			stack[2] = locals[4]
			stack[3] = rt.IntValue(10)
			stack[2] = rt.IntValue(rt.AsInt(stack[2]) * rt.AsInt(stack[3]))
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			stack[2] = locals[5]
			stack[1] = rt.IntValue(rt.AsInt(stack[1]) + rt.AsInt(stack[2]))
			{
			out := t.InvokeVirtual(m, 74, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 74, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(3)
			locals[2] = stack[0]
			block = 10
			continue
		case 10:
			{
			obj, thrown := t.NewObject(m, "NumericEdgeCases$Marker")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 79, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			out := t.InvokeSpecial(m, 83, stack[1], "NumericEdgeCases$Marker", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 83, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[3] = stack[0]
			stack[0] = locals[3]
			{
			thrown := stack[0].Ref
			if thrown == nil {
			thrown = t.Throwable("java/lang/NullPointerException")
			}
			if hpc, ok := t.HandleThrow(m, 88, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
		case 11:
			locals[4] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 91, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			stack[2] = locals[3]
			if stack[1].Ref != stack[2].Ref {
			block = 13
			continue
			}
			block = 12
			continue
		case 12:
			stack[1] = rt.IntValue(2001)
			block = 14
			continue
		case 13:
			stack[1] = rt.IntValue(-2001)
			block = 14
			continue
		case 14:
			{
			out := t.InvokeVirtual(m, 109, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 109, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 112, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
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
			{
			out := t.InvokeVirtual(m, 117, stack[1], "NumericEdgeCases$Marker", "getCause", "()Ljava/lang/Throwable;", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 117, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			if stack[1].Ref != nil {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			stack[1] = rt.IntValue(2002)
			block = 17
			continue
		case 16:
			stack[1] = rt.IntValue(-2002)
			block = 17
			continue
		case 17:
			{
			out := t.InvokeVirtual(m, 132, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 132, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 135, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 139, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 139, thrown); ok {
			stack = [4]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 28:
			block = 3
			case 89:
			block = 11
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[2]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod26(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [6]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0x0000000100000001)
			locals[1] = stack[0]
			stack[0] = locals[1]
			stack[1] = locals[0]
			stack[1] = rt.LongValue(int64(rt.AsInt(stack[1])))
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) * rt.AsLong(stack[1]))
			locals[3] = stack[0]
			stack[0] = locals[3]
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			stack[1] = locals[3]
			stack[2] = rt.IntValue(32)
			stack[1] = rt.LongValue(int64(uint64(rt.AsLong(stack[1])) >> rt.LongShiftDistance(rt.AsInt(stack[2]))))
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			stack[1] = rt.IntValue(1437226410)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) ^ rt.AsInt(stack[1]))
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[5]
			{
			out := t.InvokeVirtual(m, 28, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[5]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod27(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [19]rt.Value
	var stack [5]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	locals[2] = args[1]
	locals[3] = args[2]
	locals[4] = args[3]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.LongBitsValue(0x0000000000000003)
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 4, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) % d)
			}
			locals[6] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.LongBitsValue(0xfffffffffffffff9)
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 11, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) % d)
			}
			locals[8] = stack[0]
			stack[0] = locals[3]
			stack[1] = rt.FloatBitsValue(0x40200000)
			stack[0] = rt.FloatValue(rt.FRem(rt.AsFloat(stack[0]), rt.AsFloat(stack[1])))
			locals[10] = stack[0]
			stack[0] = locals[3]
			stack[1] = rt.FloatBitsValue(0x80000000)
			stack[0] = rt.FloatValue(rt.FRem(rt.AsFloat(stack[0]), rt.AsFloat(stack[1])))
			locals[11] = stack[0]
			stack[0] = locals[4]
			stack[1] = rt.DoubleBitsValue(0x400a000000000000)
			stack[0] = rt.DoubleValue(rt.DRem(rt.AsDouble(stack[0]), rt.AsDouble(stack[1])))
			locals[12] = stack[0]
			stack[0] = locals[4]
			stack[1] = rt.DoubleValue(0)
			stack[0] = rt.DoubleValue(rt.DRem(rt.AsDouble(stack[0]), rt.AsDouble(stack[1])))
			locals[14] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 42, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 48, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 48, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 51, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 57, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 57, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 60, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[10]
			{
			out := t.InvokeStatic(m, 65, "java/lang/Float", "floatToRawIntBits", "(F)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 65, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			out := t.InvokeVirtual(m, 68, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 68, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 71, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			stack[2] = locals[11]
			stack[1] = rt.IntValue(rt.Fcmp(rt.AsFloat(stack[1]), rt.AsFloat(stack[2]), -1))
			if rt.AsInt(stack[1]) != 0 {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[1] = rt.IntValue(0)
			block = 3
			continue
		case 2:
			stack[1] = rt.IntValue(1)
			block = 3
			continue
		case 3:
			{
			out := t.InvokeVirtual(m, 87, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 87, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 90, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			out := t.InvokeStatic(m, 95, "java/lang/Double", "doubleToRawLongBits", "(D)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 95, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			out := t.InvokeVirtual(m, 98, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 98, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 101, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[14]
			stack[2] = locals[14]
			stack[1] = rt.IntValue(rt.Dcmp(rt.AsDouble(stack[1]), rt.AsDouble(stack[2]), -1))
			if rt.AsInt(stack[1]) != 0 {
			block = 5
			continue
			}
			block = 4
			continue
		case 4:
			stack[1] = rt.IntValue(0)
			block = 6
			continue
		case 5:
			stack[1] = rt.IntValue(1)
			block = 6
			continue
		case 6:
			{
			out := t.InvokeVirtual(m, 117, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 117, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			locals[16] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[16]
			stack[1] = rt.LongValue(int64(rt.AsInt(stack[1])))
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 127, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) / d)
			}
			locals[17] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 130, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[17]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 136, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 136, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 8
			continue
		case 7:
			locals[17] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 144, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(31337)
			{
			out := t.InvokeVirtual(m, 150, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 150, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			}
			block = 8
			continue
		case 8:
			stack[0] = locals[0]
			stack[1] = locals[16]
			stack[1] = rt.LongValue(int64(rt.AsInt(stack[1])))
			{
			d := rt.AsLong(stack[1])
			if d == 0 {
			thrown := t.Throwable("java/lang/ArithmeticException")
			if hpc, ok := t.HandleThrow(m, 157, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) % d)
			}
			locals[17] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 160, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[17]
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 166, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 166, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
		case 9:
			locals[17] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			if hpc, ok := t.HandleThrow(m, 174, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
			default:
			panic("catty: generated handler transfer has no target block")
			}
			continue
			}
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(31338)
			{
			out := t.InvokeVirtual(m, 180, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			if hpc, ok := t.HandleThrow(m, 180, thrown); ok {
			stack = [5]rt.Value{}
			stack[0] = rt.RefValue(thrown)
			switch hpc {
			case 142:
			block = 7
			case 172:
			block = 9
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
			stack[0] = locals[6]
			stack[1] = locals[8]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			stack[0] = rt.IntValue(int32(rt.AsLong(stack[0])))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod28(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0x1122334455667788)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod29(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [6]rt.Value
	var stack [9]rt.Value
	_ = locals
	_ = stack
	locals[0] = rt.RefValue(receiver)
	locals[1] = args[0]
	locals[2] = args[1]
	locals[3] = args[2]
	locals[5] = args[3]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(2)
			{
			obj, thrown := t.NewReferenceArray(m, "java/lang/Object", rt.AsInt(stack[1]))
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "objectArray", "[Ljava/lang/Object;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = rt.IntValue(2)
			{
			obj, thrown := t.NewPrimitiveArray(m, 10, rt.AsInt(stack[1]))
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "intArray", "[I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = rt.IntValue(2)
			{
			obj, thrown := t.NewPrimitiveArray(m, 11, rt.AsInt(stack[1]))
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = rt.RefValue(obj)
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "longArray", "[J", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "objectArray", "[Ljava/lang/Object;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			stack[3] = locals[5]
			{
			v1 := stack[3]
			v3 := stack[1]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v3
			stack[1] = v1
			}
			{
			thrown := t.ArrayStore(stack[2], rt.AsInt(stack[3]), stack[4], 'L')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "objectField", "Ljava/lang/Object;", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "intArray", "[I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			stack[3] = locals[2]
			{
			v1 := stack[3]
			v3 := stack[1]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v3
			stack[1] = v1
			}
			{
			thrown := t.ArrayStore(stack[2], rt.AsInt(stack[3]), stack[4], 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "intField", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			stack[3] = locals[3]
			{
			v3 := stack[1]
			v2 := stack[2]
			v1 := stack[3]
			stack[1] = v1
			stack[2] = v3
			stack[3] = v2
			stack[4] = v1
			}
			{
			thrown := t.ArrayStore(stack[2], rt.AsInt(stack[3]), stack[4], 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "longField", "J", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "objectArray", "[Ljava/lang/Object;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			stack[2] = locals[0]
			stack[3] = locals[5]
			{
			v1 := stack[3]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v1
			}
			{
			thrown := t.PutField(m, stack[3], "NumericEdgeCases", "objectField", "Ljava/lang/Object;", stack[4])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'L')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "intArray", "[I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			stack[2] = locals[0]
			stack[3] = locals[2]
			stack[4] = rt.IntValue(1)
			stack[3] = rt.IntValue(rt.AsInt(stack[3]) + rt.AsInt(stack[4]))
			{
			v1 := stack[3]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v1
			}
			{
			thrown := t.PutField(m, stack[3], "NumericEdgeCases", "intField", "I", stack[4])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			stack[2] = locals[0]
			stack[3] = locals[3]
			stack[4] = rt.LongBitsValue(0x0000000000000002)
			stack[3] = rt.LongValue(rt.AsLong(stack[3]) + rt.AsLong(stack[4]))
			{
			v1 := stack[3]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v1
			}
			{
			thrown := t.PutField(m, stack[3], "NumericEdgeCases", "longField", "J", stack[4])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "intArray", "[I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			{
			v2 := stack[0]
			v1 := stack[1]
			stack[2] = v2
			stack[3] = v1
			}
			{
			v, thrown := t.ArrayLoad(stack[2], rt.AsInt(stack[3]), 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			stack[3] = rt.IntValue(3)
			stack[2] = rt.IntValue(rt.AsInt(stack[2]) << rt.IntShiftDistance(rt.AsInt(stack[3])))
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			stack[2] = locals[0]
			{
			v, thrown := t.GetField(m, stack[2], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			stack[3] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[2], rt.AsInt(stack[3]), 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			stack[3] = locals[3]
			stack[2] = rt.LongValue(rt.AsLong(stack[2]) + rt.AsLong(stack[3]))
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "intArray", "[I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v2 := stack[1]
			v1 := stack[2]
			stack[3] = v2
			stack[4] = v1
			}
			{
			v, thrown := t.ArrayLoad(stack[3], rt.AsInt(stack[4]), 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			v1 := stack[3]
			v3 := stack[1]
			v2 := stack[2]
			stack[4] = v1
			stack[3] = v2
			stack[2] = v3
			stack[1] = v1
			}
			stack[5] = rt.IntValue(1)
			stack[4] = rt.IntValue(rt.AsInt(stack[4]) + rt.AsInt(stack[5]))
			{
			thrown := t.ArrayStore(stack[2], rt.AsInt(stack[3]), stack[4], 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "intField", "I", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v2 := stack[1]
			v1 := stack[2]
			stack[3] = v2
			stack[4] = v1
			}
			{
			v, thrown := t.ArrayLoad(stack[3], rt.AsInt(stack[4]), 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[3] = v
			}
			{
			v3 := stack[1]
			v2 := stack[2]
			v1 := stack[3]
			stack[1] = v1
			stack[2] = v3
			stack[3] = v2
			stack[4] = v1
			}
			stack[5] = rt.LongValue(1)
			stack[4] = rt.LongValue(rt.AsLong(stack[4]) + rt.AsLong(stack[5]))
			{
			thrown := t.ArrayStore(stack[2], rt.AsInt(stack[3]), stack[4], 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			thrown := t.PutField(m, stack[0], "NumericEdgeCases", "longField", "J", stack[1])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "intField", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 163, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "longField", "J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 174, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "intArray", "[I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[1], rt.AsInt(stack[2]), 'I')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			{
			out := t.InvokeVirtual(m, 186, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[1], rt.AsInt(stack[2]), 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 199, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "objectArray", "[Ljava/lang/Object;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[1], rt.AsInt(stack[2]), 'L')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[5]
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
			out := t.InvokeVirtual(m, 221, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[1]
			stack[2] = locals[0]
			{
			v, thrown := t.GetField(m, stack[2], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			stack[3] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[2], rt.AsInt(stack[3]), 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = v
			}
			stack[3] = locals[3]
			stack[2] = rt.LongValue(rt.AsLong(stack[2]) + rt.AsLong(stack[3]))
			{
			thrown := t.ArrayStore(stack[0], rt.AsInt(stack[1]), stack[2], 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[0]
			{
			v, thrown := t.GetField(m, stack[1], "NumericEdgeCases", "longArray", "[J")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[2] = locals[1]
			{
			v, thrown := t.ArrayLoad(stack[1], rt.AsInt(stack[2]), 'J')
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = v
			}
			stack[1] = rt.IntValue(int32(rt.AsLong(stack[1])))
			{
			out := t.InvokeVirtual(m, 248, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[0]
			{
			out := t.InvokeVirtual(m, 252, stack[0], "NumericEdgeCases", "longMark", "()J", stack[1:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[0]
			{
			v, thrown := t.GetField(m, stack[0], "NumericEdgeCases", "intField", "I")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod30(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(13)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) * rt.AsInt(stack[1]))
			stack[1] = rt.IntValue(99)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod31(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [9]rt.Value
	var stack [7]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			{
			n, thrown := t.ArrayLength(stack[0])
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.IntValue(n)
			}
			locals[1] = stack[0]
			stack[0] = rt.LongBitsValue(0x1122334455667788)
			stack[1] = locals[1]
			{
			out := t.InvokeStatic(m, 7, "NumericEdgeCases", "longOps", "(JI)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[2] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[2]
			{
			out := t.InvokeVirtual(m, 15, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(1898136935)
			stack[1] = rt.IntValue(-37)
			{
			out := t.InvokeStatic(m, 23, "NumericEdgeCases", "intArith", "(II)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[3] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[3]
			{
			out := t.InvokeVirtual(m, 31, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[1]
			stack[1] = rt.IntValue(10)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			out := t.InvokeStatic(m, 38, "NumericEdgeCases", "iincDemo", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[4] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[4]
			{
			out := t.InvokeVirtual(m, 48, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			out := t.InvokeStatic(m, 51, "NumericEdgeCases", "constantKinds", "()I", stack[0:0])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[5] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[5]
			{
			out := t.InvokeVirtual(m, 61, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(305419896)
			stack[1] = rt.LongBitsValue(0x1122334455667788)
			{
			out := t.InvokeStatic(m, 70, "NumericEdgeCases", "convWiden", "(IJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[6] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[6]
			{
			out := t.InvokeVirtual(m, 80, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.FloatBitsValue(0x3fc00000)
			stack[1] = rt.DoubleBitsValue(0xc002000000000000)
			{
			out := t.InvokeStatic(m, 89, "NumericEdgeCases", "convFloat", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			locals[7] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[7]
			{
			out := t.InvokeVirtual(m, 99, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 103, "NumericEdgeCases", "specialFloat", "(I)F", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 107, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeStatic(m, 110, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 115, "NumericEdgeCases", "specialFloat", "(I)F", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 119, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeStatic(m, 122, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 127, "NumericEdgeCases", "specialFloat", "(I)F", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 131, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeStatic(m, 134, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0x4f32d05e)
			stack[1] = rt.DoubleBitsValue(0xc1e0000000200000)
			{
			out := t.InvokeStatic(m, 144, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0x4f000000)
			stack[1] = rt.DoubleBitsValue(0x43e0000000000000)
			{
			out := t.InvokeStatic(m, 154, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0xcf000000)
			stack[1] = rt.DoubleBitsValue(0xc3e0000000000000)
			{
			out := t.InvokeStatic(m, 164, "NumericEdgeCases", "saturation", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0x3fc00000)
			stack[1] = rt.FloatBitsValue(0xc0200000)
			{
			out := t.InvokeStatic(m, 174, "NumericEdgeCases", "floatOps", "(FF)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0x80000000)
			stack[1] = rt.FloatValue(1)
			{
			out := t.InvokeStatic(m, 182, "NumericEdgeCases", "floatOps", "(FF)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatValue(1)
			stack[1] = rt.FloatBitsValue(0x80000000)
			{
			out := t.InvokeStatic(m, 190, "NumericEdgeCases", "floatOps", "(FF)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.FloatBitsValue(0xbf800000)
			stack[1] = rt.FloatValue(0)
			{
			out := t.InvokeStatic(m, 198, "NumericEdgeCases", "floatOps", "(FF)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 206, "NumericEdgeCases", "specialFloat", "(I)F", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 210, "NumericEdgeCases", "specialFloat", "(I)F", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 213, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 216, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 223, "NumericEdgeCases", "specialFloat", "(I)F", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.FloatValue(1)
			{
			out := t.InvokeStatic(m, 227, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 230, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.FloatValue(1)
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 238, "NumericEdgeCases", "specialFloat", "(I)F", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 241, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 244, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.FloatValue(0)
			stack[2] = rt.FloatBitsValue(0x80000000)
			{
			out := t.InvokeStatic(m, 254, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 257, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 264, "NumericEdgeCases", "specialFloat", "(I)F", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 268, "NumericEdgeCases", "specialFloat", "(I)F", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 271, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 274, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 281, "NumericEdgeCases", "specialFloat", "(I)F", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 285, "NumericEdgeCases", "specialFloat", "(I)F", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 288, "NumericEdgeCases", "compareFloat", "(FF)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 291, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.DoubleBitsValue(0x3ff8000000000000)
			stack[1] = rt.DoubleBitsValue(0xc004000000000000)
			{
			out := t.InvokeStatic(m, 300, "NumericEdgeCases", "doubleOps", "(DD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.DoubleBitsValue(0x8000000000000000)
			stack[1] = rt.DoubleValue(1)
			{
			out := t.InvokeStatic(m, 308, "NumericEdgeCases", "doubleOps", "(DD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.DoubleValue(1)
			stack[1] = rt.DoubleBitsValue(0x8000000000000000)
			{
			out := t.InvokeStatic(m, 316, "NumericEdgeCases", "doubleOps", "(DD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.DoubleBitsValue(0x7fe1ccf385ebc8a0)
			stack[1] = rt.DoubleBitsValue(0x000730d67819e8d2)
			{
			out := t.InvokeStatic(m, 326, "NumericEdgeCases", "doubleOps", "(DD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 334, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 338, "NumericEdgeCases", "specialDouble", "(I)D", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 341, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 344, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 351, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.DoubleValue(1)
			{
			out := t.InvokeStatic(m, 355, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 358, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.DoubleValue(1)
			stack[2] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 366, "NumericEdgeCases", "specialDouble", "(I)D", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 369, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 372, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.DoubleValue(0)
			stack[2] = rt.DoubleBitsValue(0x8000000000000000)
			{
			out := t.InvokeStatic(m, 382, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 385, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 392, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 396, "NumericEdgeCases", "specialDouble", "(I)D", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 399, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 402, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 409, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			stack[2] = rt.IntValue(1)
			{
			out := t.InvokeStatic(m, 413, "NumericEdgeCases", "specialDouble", "(I)D", stack[2:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[2] = out.Value
			}
			{
			out := t.InvokeStatic(m, 416, "NumericEdgeCases", "compareDouble", "(DD)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 419, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 423, "NumericEdgeCases", "specialFloat", "(I)F", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 427, "NumericEdgeCases", "specialDouble", "(I)D", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeStatic(m, 430, "NumericEdgeCases", "nanPaths", "(FD)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(2143289345)
			stack[1] = rt.LongBitsValue(0x7ff8000000000001)
			{
			out := t.InvokeStatic(m, 440, "NumericEdgeCases", "bitRoundTrip", "(IJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(-2147483648)
			stack[1] = rt.LongBitsValue(0x8000000000000000)
			{
			out := t.InvokeStatic(m, 449, "NumericEdgeCases", "bitRoundTrip", "(IJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.LongBitsValue(0x1122334455667788)
			stack[1] = rt.LongBitsValue(0x1122334455667788)
			{
			out := t.InvokeStatic(m, 459, "NumericEdgeCases", "longCompare", "(JJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.LongBitsValue(0xfffffffffffffffb)
			stack[1] = rt.LongBitsValue(0x0000000000000009)
			{
			out := t.InvokeStatic(m, 469, "NumericEdgeCases", "longCompare", "(JJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.LongBitsValue(0x0000000000000009)
			stack[1] = rt.LongBitsValue(0xfffffffffffffffb)
			{
			out := t.InvokeStatic(m, 479, "NumericEdgeCases", "longCompare", "(JJ)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(7)
			stack[1] = rt.IntValue(3)
			{
			out := t.InvokeStatic(m, 486, "NumericEdgeCases", "branchKinds", "(II)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(-2)
			stack[1] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 493, "NumericEdgeCases", "branchKinds", "(II)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[0]
			stack[1] = rt.RefValue(nil)
			{
			out := t.InvokeStatic(m, 499, "NumericEdgeCases", "refBranches", "(Ljava/lang/Object;Ljava/lang/Object;)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[0]
			stack[1] = locals[0]
			{
			out := t.InvokeStatic(m, 505, "NumericEdgeCases", "refBranches", "(Ljava/lang/Object;Ljava/lang/Object;)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.RefValue(nil)
			stack[1] = rt.RefValue(nil)
			{
			out := t.InvokeStatic(m, 511, "NumericEdgeCases", "refBranches", "(Ljava/lang/Object;Ljava/lang/Object;)I", stack[0:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = rt.IntValue(305419896)
			{
			out := t.InvokeStatic(m, 518, "NumericEdgeCases", "narrowing", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[1]
			stack[1] = rt.IntValue(2)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			{
			out := t.InvokeStatic(m, 525, "NumericEdgeCases", "leafUse", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.LongBitsValue(0x1122334455667788)
			stack[2] = rt.IntValue(-37)
			stack[3] = rt.FloatBitsValue(0x3fc00000)
			stack[4] = rt.DoubleBitsValue(0xc002000000000000)
			{
			out := t.InvokeStatic(m, 543, "NumericEdgeCases", "remainders", "(JIFD)I", stack[1:5])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 546, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			{
			obj, thrown := t.NewObject(m, "NumericEdgeCases")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = rt.RefValue(obj)
			}
			stack[1] = stack[0]
			{
			out := t.InvokeSpecial(m, 553, stack[1], "NumericEdgeCases", "<init>", "()V", stack[2:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			locals[8] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[8]
			stack[2] = rt.IntValue(1)
			stack[3] = rt.IntValue(7)
			stack[4] = rt.LongBitsValue(0x1122334455667788)
			stack[5] = locals[0]
			{
			out := t.InvokeVirtual(m, 570, stack[1], "NumericEdgeCases", "stackForms", "(IIJLjava/lang/Object;)I", stack[2:6])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 573, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = rt.IntValue(0)
			{
			out := t.InvokeStatic(m, 577, "NumericEdgeCases", "exceptionPath", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			stack[0] = locals[1]
			stack[1] = rt.IntValue(3)
			if rt.AsInt(stack[0]) != rt.AsInt(stack[1]) {
			block = 2
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[1]
			{
			out := t.InvokeStatic(m, 587, "NumericEdgeCases", "extraPath", "(I)I", stack[0:1])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = out.Value
			}
			block = 2
			continue
		case 2:
			return rt.Normal()
		}
	}
}

