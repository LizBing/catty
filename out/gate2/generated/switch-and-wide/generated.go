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

func cattyMethod1(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case 10:
			block = 1
			continue
			case 11:
			block = 2
			continue
			case 12:
			block = 3
			continue
			case 13:
			block = 4
			continue
			case 14:
			block = 5
			continue
			case 15:
			block = 6
			continue
			case 16:
			block = 7
			continue
			case 17:
			block = 8
			continue
			case 18:
			block = 9
			continue
			case 19:
			block = 10
			continue
			default:
			block = 11
			continue
			}
		case 1:
			stack[0] = rt.IntValue(201)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(202)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(203)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = rt.IntValue(204)
			return rt.OkValue(stack[0])
		case 5:
			stack[0] = rt.IntValue(205)
			return rt.OkValue(stack[0])
		case 6:
			stack[0] = rt.IntValue(206)
			return rt.OkValue(stack[0])
		case 7:
			stack[0] = rt.IntValue(207)
			return rt.OkValue(stack[0])
		case 8:
			stack[0] = rt.IntValue(208)
			return rt.OkValue(stack[0])
		case 9:
			stack[0] = rt.IntValue(209)
			return rt.OkValue(stack[0])
		case 10:
			stack[0] = rt.IntValue(210)
			return rt.OkValue(stack[0])
		case 11:
			stack[0] = rt.IntValue(-10)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod2(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case -1000:
			block = 1
			continue
			case -7:
			block = 2
			continue
			case 0:
			block = 3
			continue
			case 3:
			block = 4
			continue
			case 100:
			block = 5
			continue
			case 10000:
			block = 6
			continue
			case 1000000:
			block = 7
			continue
			default:
			block = 8
			continue
			}
		case 1:
			stack[0] = rt.IntValue(5)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(6)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(7)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = rt.IntValue(8)
			return rt.OkValue(stack[0])
		case 5:
			stack[0] = rt.IntValue(9)
			return rt.OkValue(stack[0])
		case 6:
			stack[0] = rt.IntValue(10)
			return rt.OkValue(stack[0])
		case 7:
			stack[0] = rt.IntValue(11)
			return rt.OkValue(stack[0])
		case 8:
			stack[0] = rt.IntValue(-5)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod3(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [1]rt.Value
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
			switch rt.AsInt(stack[0]) {
			case 0:
			block = 1
			continue
			case 1:
			block = 6
			continue
			case 2:
			block = 10
			continue
			default:
			block = 11
			continue
			}
		case 1:
			stack[0] = locals[1]
			switch rt.AsInt(stack[0]) {
			case 0:
			block = 2
			continue
			case 1:
			block = 3
			continue
			case 2:
			block = 4
			continue
			default:
			block = 5
			continue
			}
		case 2:
			stack[0] = rt.IntValue(13)
			locals[2] = stack[0]
			block = 12
			continue
		case 3:
			stack[0] = rt.IntValue(14)
			locals[2] = stack[0]
			block = 12
			continue
		case 4:
			stack[0] = rt.IntValue(15)
			locals[2] = stack[0]
			block = 12
			continue
		case 5:
			stack[0] = rt.IntValue(16)
			locals[2] = stack[0]
			block = 12
			continue
		case 6:
			stack[0] = locals[1]
			switch rt.AsInt(stack[0]) {
			case 5:
			block = 7
			continue
			case 6:
			block = 8
			continue
			default:
			block = 9
			continue
			}
		case 7:
			stack[0] = rt.IntValue(17)
			locals[2] = stack[0]
			block = 12
			continue
		case 8:
			stack[0] = rt.IntValue(18)
			locals[2] = stack[0]
			block = 12
			continue
		case 9:
			stack[0] = rt.IntValue(19)
			locals[2] = stack[0]
			block = 12
			continue
		case 10:
			stack[0] = rt.IntValue(20)
			locals[2] = stack[0]
			block = 12
			continue
		case 11:
			stack[0] = rt.IntValue(21)
			locals[2] = stack[0]
			block = 12
			continue
		case 12:
			stack[0] = locals[2]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod4(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [3]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[1] = stack[0]
			stack[0] = rt.IntValue(0)
			locals[2] = stack[0]
			block = 1
			continue
		case 1:
			stack[0] = locals[2]
			stack[1] = locals[0]
			if rt.AsInt(stack[0]) >= rt.AsInt(stack[1]) {
			block = 13
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = locals[2]
			switch rt.AsInt(stack[0]) {
			case 0:
			block = 3
			continue
			case 1:
			block = 4
			continue
			case 2:
			block = 5
			continue
			case 3:
			block = 6
			continue
			case 4:
			block = 7
			continue
			case 5:
			block = 8
			continue
			default:
			block = 9
			continue
			}
		case 3:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 1)
			block = 10
			continue
		case 4:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 2)
			block = 10
			continue
		case 5:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 3)
			block = 10
			continue
		case 6:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 4)
			block = 10
			continue
		case 7:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 5)
			block = 10
			continue
		case 8:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 6)
			block = 10
			continue
		case 9:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 7)
			block = 10
			continue
		case 10:
			stack[0] = locals[1]
			stack[1] = rt.IntValue(15)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 12
			continue
			}
			block = 11
			continue
		case 11:
			block = 13
			continue
		case 12:
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + 1)
			block = 1
			continue
		case 13:
			stack[0] = locals[0]
			locals[2] = stack[0]
			block = 14
			continue
		case 14:
			stack[0] = locals[2]
			if rt.AsInt(stack[0]) <= 0 {
			block = 16
			continue
			}
			block = 15
			continue
		case 15:
			stack[0] = locals[1]
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[1] = stack[0]
			locals[2] = rt.IntValue(rt.AsInt(locals[2]) + -1)
			block = 14
			continue
		case 16:
			stack[0] = locals[1]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod5(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case 97:
			block = 1
			continue
			case 98:
			block = 2
			continue
			case 122:
			block = 3
			continue
			default:
			block = 4
			continue
			}
		case 1:
			stack[0] = rt.IntValue(31)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(32)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(33)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = rt.IntValue(34)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod6(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case -2:
			block = 1
			continue
			case 0:
			block = 2
			continue
			case 5:
			block = 3
			continue
			default:
			block = 4
			continue
			}
		case 1:
			stack[0] = rt.IntValue(41)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(42)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(43)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = rt.IntValue(44)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod7(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case -300:
			block = 1
			continue
			case 0:
			block = 2
			continue
			case 300:
			block = 3
			continue
			default:
			block = 4
			continue
			}
		case 1:
			stack[0] = rt.IntValue(51)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(52)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(53)
			return rt.OkValue(stack[0])
		case 4:
			stack[0] = rt.IntValue(54)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod8(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[1] = stack[0]
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case 5:
			block = 1
			continue
			case 6:
			block = 2
			continue
			case 7:
			block = 3
			continue
			case 8:
			block = 2
			continue
			case 9:
			block = 4
			continue
			default:
			block = 2
			continue
			}
		case 1:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 1)
			block = 2
			continue
		case 2:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 2)
			block = 3
			continue
		case 3:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 4)
			block = 5
			continue
		case 4:
			locals[1] = rt.IntValue(rt.AsInt(locals[1]) + 8)
			block = 5
			continue
		case 5:
			stack[0] = locals[1]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod9(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(100)
			locals[1] = stack[0]
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case 0:
			block = 1
			continue
			case 1:
			block = 2
			continue
			default:
			block = 3
			continue
			}
		case 1:
			stack[0] = rt.IntValue(101)
			locals[1] = stack[0]
			block = 3
			continue
		case 2:
			stack[0] = rt.IntValue(102)
			locals[1] = stack[0]
			block = 3
			continue
		case 3:
			stack[0] = locals[1]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod10(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [2]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[1] = stack[0]
			stack[0] = locals[0]
			stack[1] = rt.IntValue(100)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 8
			continue
			}
			block = 1
			continue
		case 1:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(1000)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 5
			continue
			}
			block = 2
			continue
		case 2:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(10000)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 4
			continue
			}
			block = 3
			continue
		case 3:
			stack[0] = rt.IntValue(4)
			locals[1] = stack[0]
			block = 17
			continue
		case 4:
			stack[0] = rt.IntValue(3)
			locals[1] = stack[0]
			block = 17
			continue
		case 5:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(500)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 7
			continue
			}
			block = 6
			continue
		case 6:
			stack[0] = rt.IntValue(2)
			locals[1] = stack[0]
			block = 17
			continue
		case 7:
			stack[0] = rt.IntValue(1)
			locals[1] = stack[0]
			block = 17
			continue
		case 8:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(10)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 12
			continue
			}
			block = 9
			continue
		case 9:
			stack[0] = locals[0]
			stack[1] = rt.IntValue(50)
			if rt.AsInt(stack[0]) <= rt.AsInt(stack[1]) {
			block = 11
			continue
			}
			block = 10
			continue
		case 10:
			stack[0] = rt.IntValue(8)
			locals[1] = stack[0]
			block = 17
			continue
		case 11:
			stack[0] = rt.IntValue(7)
			locals[1] = stack[0]
			block = 17
			continue
		case 12:
			stack[0] = locals[0]
			if rt.AsInt(stack[0]) <= 0 {
			block = 14
			continue
			}
			block = 13
			continue
		case 13:
			stack[0] = rt.IntValue(6)
			locals[1] = stack[0]
			block = 17
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
			stack[0] = rt.IntValue(5)
			locals[1] = stack[0]
			block = 17
			continue
		case 16:
			stack[0] = rt.IntValue(9)
			locals[1] = stack[0]
			block = 17
			continue
		case 17:
			stack[0] = locals[1]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod11(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [302]rt.Value
	var stack [2]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.IntValue(0)
			locals[1] = stack[0]
			stack[0] = rt.IntValue(1)
			locals[2] = stack[0]
			stack[0] = rt.IntValue(2)
			locals[3] = stack[0]
			stack[0] = rt.IntValue(3)
			locals[4] = stack[0]
			stack[0] = rt.IntValue(4)
			locals[5] = stack[0]
			stack[0] = rt.IntValue(5)
			locals[6] = stack[0]
			stack[0] = rt.IntValue(6)
			locals[7] = stack[0]
			stack[0] = rt.IntValue(7)
			locals[8] = stack[0]
			stack[0] = rt.IntValue(8)
			locals[9] = stack[0]
			stack[0] = rt.IntValue(9)
			locals[10] = stack[0]
			stack[0] = rt.IntValue(10)
			locals[11] = stack[0]
			stack[0] = rt.IntValue(11)
			locals[12] = stack[0]
			stack[0] = rt.IntValue(12)
			locals[13] = stack[0]
			stack[0] = rt.IntValue(13)
			locals[14] = stack[0]
			stack[0] = rt.IntValue(14)
			locals[15] = stack[0]
			stack[0] = rt.IntValue(15)
			locals[16] = stack[0]
			stack[0] = rt.IntValue(16)
			locals[17] = stack[0]
			stack[0] = rt.IntValue(17)
			locals[18] = stack[0]
			stack[0] = rt.IntValue(18)
			locals[19] = stack[0]
			stack[0] = rt.IntValue(19)
			locals[20] = stack[0]
			stack[0] = rt.IntValue(20)
			locals[21] = stack[0]
			stack[0] = rt.IntValue(21)
			locals[22] = stack[0]
			stack[0] = rt.IntValue(22)
			locals[23] = stack[0]
			stack[0] = rt.IntValue(23)
			locals[24] = stack[0]
			stack[0] = rt.IntValue(24)
			locals[25] = stack[0]
			stack[0] = rt.IntValue(25)
			locals[26] = stack[0]
			stack[0] = rt.IntValue(26)
			locals[27] = stack[0]
			stack[0] = rt.IntValue(27)
			locals[28] = stack[0]
			stack[0] = rt.IntValue(28)
			locals[29] = stack[0]
			stack[0] = rt.IntValue(29)
			locals[30] = stack[0]
			stack[0] = rt.IntValue(30)
			locals[31] = stack[0]
			stack[0] = rt.IntValue(31)
			locals[32] = stack[0]
			stack[0] = rt.IntValue(32)
			locals[33] = stack[0]
			stack[0] = rt.IntValue(33)
			locals[34] = stack[0]
			stack[0] = rt.IntValue(34)
			locals[35] = stack[0]
			stack[0] = rt.IntValue(35)
			locals[36] = stack[0]
			stack[0] = rt.IntValue(36)
			locals[37] = stack[0]
			stack[0] = rt.IntValue(37)
			locals[38] = stack[0]
			stack[0] = rt.IntValue(38)
			locals[39] = stack[0]
			stack[0] = rt.IntValue(39)
			locals[40] = stack[0]
			stack[0] = rt.IntValue(40)
			locals[41] = stack[0]
			stack[0] = rt.IntValue(41)
			locals[42] = stack[0]
			stack[0] = rt.IntValue(42)
			locals[43] = stack[0]
			stack[0] = rt.IntValue(43)
			locals[44] = stack[0]
			stack[0] = rt.IntValue(44)
			locals[45] = stack[0]
			stack[0] = rt.IntValue(45)
			locals[46] = stack[0]
			stack[0] = rt.IntValue(46)
			locals[47] = stack[0]
			stack[0] = rt.IntValue(47)
			locals[48] = stack[0]
			stack[0] = rt.IntValue(48)
			locals[49] = stack[0]
			stack[0] = rt.IntValue(49)
			locals[50] = stack[0]
			stack[0] = rt.IntValue(50)
			locals[51] = stack[0]
			stack[0] = rt.IntValue(51)
			locals[52] = stack[0]
			stack[0] = rt.IntValue(52)
			locals[53] = stack[0]
			stack[0] = rt.IntValue(53)
			locals[54] = stack[0]
			stack[0] = rt.IntValue(54)
			locals[55] = stack[0]
			stack[0] = rt.IntValue(55)
			locals[56] = stack[0]
			stack[0] = rt.IntValue(56)
			locals[57] = stack[0]
			stack[0] = rt.IntValue(57)
			locals[58] = stack[0]
			stack[0] = rt.IntValue(58)
			locals[59] = stack[0]
			stack[0] = rt.IntValue(59)
			locals[60] = stack[0]
			stack[0] = rt.IntValue(60)
			locals[61] = stack[0]
			stack[0] = rt.IntValue(61)
			locals[62] = stack[0]
			stack[0] = rt.IntValue(62)
			locals[63] = stack[0]
			stack[0] = rt.IntValue(63)
			locals[64] = stack[0]
			stack[0] = rt.IntValue(64)
			locals[65] = stack[0]
			stack[0] = rt.IntValue(65)
			locals[66] = stack[0]
			stack[0] = rt.IntValue(66)
			locals[67] = stack[0]
			stack[0] = rt.IntValue(67)
			locals[68] = stack[0]
			stack[0] = rt.IntValue(68)
			locals[69] = stack[0]
			stack[0] = rt.IntValue(69)
			locals[70] = stack[0]
			stack[0] = rt.IntValue(70)
			locals[71] = stack[0]
			stack[0] = rt.IntValue(71)
			locals[72] = stack[0]
			stack[0] = rt.IntValue(72)
			locals[73] = stack[0]
			stack[0] = rt.IntValue(73)
			locals[74] = stack[0]
			stack[0] = rt.IntValue(74)
			locals[75] = stack[0]
			stack[0] = rt.IntValue(75)
			locals[76] = stack[0]
			stack[0] = rt.IntValue(76)
			locals[77] = stack[0]
			stack[0] = rt.IntValue(77)
			locals[78] = stack[0]
			stack[0] = rt.IntValue(78)
			locals[79] = stack[0]
			stack[0] = rt.IntValue(79)
			locals[80] = stack[0]
			stack[0] = rt.IntValue(80)
			locals[81] = stack[0]
			stack[0] = rt.IntValue(81)
			locals[82] = stack[0]
			stack[0] = rt.IntValue(82)
			locals[83] = stack[0]
			stack[0] = rt.IntValue(83)
			locals[84] = stack[0]
			stack[0] = rt.IntValue(84)
			locals[85] = stack[0]
			stack[0] = rt.IntValue(85)
			locals[86] = stack[0]
			stack[0] = rt.IntValue(86)
			locals[87] = stack[0]
			stack[0] = rt.IntValue(87)
			locals[88] = stack[0]
			stack[0] = rt.IntValue(88)
			locals[89] = stack[0]
			stack[0] = rt.IntValue(89)
			locals[90] = stack[0]
			stack[0] = rt.IntValue(90)
			locals[91] = stack[0]
			stack[0] = rt.IntValue(91)
			locals[92] = stack[0]
			stack[0] = rt.IntValue(92)
			locals[93] = stack[0]
			stack[0] = rt.IntValue(93)
			locals[94] = stack[0]
			stack[0] = rt.IntValue(94)
			locals[95] = stack[0]
			stack[0] = rt.IntValue(95)
			locals[96] = stack[0]
			stack[0] = rt.IntValue(96)
			locals[97] = stack[0]
			stack[0] = rt.IntValue(97)
			locals[98] = stack[0]
			stack[0] = rt.IntValue(98)
			locals[99] = stack[0]
			stack[0] = rt.IntValue(99)
			locals[100] = stack[0]
			stack[0] = rt.IntValue(100)
			locals[101] = stack[0]
			stack[0] = rt.IntValue(101)
			locals[102] = stack[0]
			stack[0] = rt.IntValue(102)
			locals[103] = stack[0]
			stack[0] = rt.IntValue(103)
			locals[104] = stack[0]
			stack[0] = rt.IntValue(104)
			locals[105] = stack[0]
			stack[0] = rt.IntValue(105)
			locals[106] = stack[0]
			stack[0] = rt.IntValue(106)
			locals[107] = stack[0]
			stack[0] = rt.IntValue(107)
			locals[108] = stack[0]
			stack[0] = rt.IntValue(108)
			locals[109] = stack[0]
			stack[0] = rt.IntValue(109)
			locals[110] = stack[0]
			stack[0] = rt.IntValue(110)
			locals[111] = stack[0]
			stack[0] = rt.IntValue(111)
			locals[112] = stack[0]
			stack[0] = rt.IntValue(112)
			locals[113] = stack[0]
			stack[0] = rt.IntValue(113)
			locals[114] = stack[0]
			stack[0] = rt.IntValue(114)
			locals[115] = stack[0]
			stack[0] = rt.IntValue(115)
			locals[116] = stack[0]
			stack[0] = rt.IntValue(116)
			locals[117] = stack[0]
			stack[0] = rt.IntValue(117)
			locals[118] = stack[0]
			stack[0] = rt.IntValue(118)
			locals[119] = stack[0]
			stack[0] = rt.IntValue(119)
			locals[120] = stack[0]
			stack[0] = rt.IntValue(120)
			locals[121] = stack[0]
			stack[0] = rt.IntValue(121)
			locals[122] = stack[0]
			stack[0] = rt.IntValue(122)
			locals[123] = stack[0]
			stack[0] = rt.IntValue(123)
			locals[124] = stack[0]
			stack[0] = rt.IntValue(124)
			locals[125] = stack[0]
			stack[0] = rt.IntValue(125)
			locals[126] = stack[0]
			stack[0] = rt.IntValue(126)
			locals[127] = stack[0]
			stack[0] = rt.IntValue(127)
			locals[128] = stack[0]
			stack[0] = rt.IntValue(128)
			locals[129] = stack[0]
			stack[0] = rt.IntValue(129)
			locals[130] = stack[0]
			stack[0] = rt.IntValue(130)
			locals[131] = stack[0]
			stack[0] = rt.IntValue(131)
			locals[132] = stack[0]
			stack[0] = rt.IntValue(132)
			locals[133] = stack[0]
			stack[0] = rt.IntValue(133)
			locals[134] = stack[0]
			stack[0] = rt.IntValue(134)
			locals[135] = stack[0]
			stack[0] = rt.IntValue(135)
			locals[136] = stack[0]
			stack[0] = rt.IntValue(136)
			locals[137] = stack[0]
			stack[0] = rt.IntValue(137)
			locals[138] = stack[0]
			stack[0] = rt.IntValue(138)
			locals[139] = stack[0]
			stack[0] = rt.IntValue(139)
			locals[140] = stack[0]
			stack[0] = rt.IntValue(140)
			locals[141] = stack[0]
			stack[0] = rt.IntValue(141)
			locals[142] = stack[0]
			stack[0] = rt.IntValue(142)
			locals[143] = stack[0]
			stack[0] = rt.IntValue(143)
			locals[144] = stack[0]
			stack[0] = rt.IntValue(144)
			locals[145] = stack[0]
			stack[0] = rt.IntValue(145)
			locals[146] = stack[0]
			stack[0] = rt.IntValue(146)
			locals[147] = stack[0]
			stack[0] = rt.IntValue(147)
			locals[148] = stack[0]
			stack[0] = rt.IntValue(148)
			locals[149] = stack[0]
			stack[0] = rt.IntValue(149)
			locals[150] = stack[0]
			stack[0] = rt.IntValue(150)
			locals[151] = stack[0]
			stack[0] = rt.IntValue(151)
			locals[152] = stack[0]
			stack[0] = rt.IntValue(152)
			locals[153] = stack[0]
			stack[0] = rt.IntValue(153)
			locals[154] = stack[0]
			stack[0] = rt.IntValue(154)
			locals[155] = stack[0]
			stack[0] = rt.IntValue(155)
			locals[156] = stack[0]
			stack[0] = rt.IntValue(156)
			locals[157] = stack[0]
			stack[0] = rt.IntValue(157)
			locals[158] = stack[0]
			stack[0] = rt.IntValue(158)
			locals[159] = stack[0]
			stack[0] = rt.IntValue(159)
			locals[160] = stack[0]
			stack[0] = rt.IntValue(160)
			locals[161] = stack[0]
			stack[0] = rt.IntValue(161)
			locals[162] = stack[0]
			stack[0] = rt.IntValue(162)
			locals[163] = stack[0]
			stack[0] = rt.IntValue(163)
			locals[164] = stack[0]
			stack[0] = rt.IntValue(164)
			locals[165] = stack[0]
			stack[0] = rt.IntValue(165)
			locals[166] = stack[0]
			stack[0] = rt.IntValue(166)
			locals[167] = stack[0]
			stack[0] = rt.IntValue(167)
			locals[168] = stack[0]
			stack[0] = rt.IntValue(168)
			locals[169] = stack[0]
			stack[0] = rt.IntValue(169)
			locals[170] = stack[0]
			stack[0] = rt.IntValue(170)
			locals[171] = stack[0]
			stack[0] = rt.IntValue(171)
			locals[172] = stack[0]
			stack[0] = rt.IntValue(172)
			locals[173] = stack[0]
			stack[0] = rt.IntValue(173)
			locals[174] = stack[0]
			stack[0] = rt.IntValue(174)
			locals[175] = stack[0]
			stack[0] = rt.IntValue(175)
			locals[176] = stack[0]
			stack[0] = rt.IntValue(176)
			locals[177] = stack[0]
			stack[0] = rt.IntValue(177)
			locals[178] = stack[0]
			stack[0] = rt.IntValue(178)
			locals[179] = stack[0]
			stack[0] = rt.IntValue(179)
			locals[180] = stack[0]
			stack[0] = rt.IntValue(180)
			locals[181] = stack[0]
			stack[0] = rt.IntValue(181)
			locals[182] = stack[0]
			stack[0] = rt.IntValue(182)
			locals[183] = stack[0]
			stack[0] = rt.IntValue(183)
			locals[184] = stack[0]
			stack[0] = rt.IntValue(184)
			locals[185] = stack[0]
			stack[0] = rt.IntValue(185)
			locals[186] = stack[0]
			stack[0] = rt.IntValue(186)
			locals[187] = stack[0]
			stack[0] = rt.IntValue(187)
			locals[188] = stack[0]
			stack[0] = rt.IntValue(188)
			locals[189] = stack[0]
			stack[0] = rt.IntValue(189)
			locals[190] = stack[0]
			stack[0] = rt.IntValue(190)
			locals[191] = stack[0]
			stack[0] = rt.IntValue(191)
			locals[192] = stack[0]
			stack[0] = rt.IntValue(192)
			locals[193] = stack[0]
			stack[0] = rt.IntValue(193)
			locals[194] = stack[0]
			stack[0] = rt.IntValue(194)
			locals[195] = stack[0]
			stack[0] = rt.IntValue(195)
			locals[196] = stack[0]
			stack[0] = rt.IntValue(196)
			locals[197] = stack[0]
			stack[0] = rt.IntValue(197)
			locals[198] = stack[0]
			stack[0] = rt.IntValue(198)
			locals[199] = stack[0]
			stack[0] = rt.IntValue(199)
			locals[200] = stack[0]
			stack[0] = rt.IntValue(200)
			locals[201] = stack[0]
			stack[0] = rt.IntValue(201)
			locals[202] = stack[0]
			stack[0] = rt.IntValue(202)
			locals[203] = stack[0]
			stack[0] = rt.IntValue(203)
			locals[204] = stack[0]
			stack[0] = rt.IntValue(204)
			locals[205] = stack[0]
			stack[0] = rt.IntValue(205)
			locals[206] = stack[0]
			stack[0] = rt.IntValue(206)
			locals[207] = stack[0]
			stack[0] = rt.IntValue(207)
			locals[208] = stack[0]
			stack[0] = rt.IntValue(208)
			locals[209] = stack[0]
			stack[0] = rt.IntValue(209)
			locals[210] = stack[0]
			stack[0] = rt.IntValue(210)
			locals[211] = stack[0]
			stack[0] = rt.IntValue(211)
			locals[212] = stack[0]
			stack[0] = rt.IntValue(212)
			locals[213] = stack[0]
			stack[0] = rt.IntValue(213)
			locals[214] = stack[0]
			stack[0] = rt.IntValue(214)
			locals[215] = stack[0]
			stack[0] = rt.IntValue(215)
			locals[216] = stack[0]
			stack[0] = rt.IntValue(216)
			locals[217] = stack[0]
			stack[0] = rt.IntValue(217)
			locals[218] = stack[0]
			stack[0] = rt.IntValue(218)
			locals[219] = stack[0]
			stack[0] = rt.IntValue(219)
			locals[220] = stack[0]
			stack[0] = rt.IntValue(220)
			locals[221] = stack[0]
			stack[0] = rt.IntValue(221)
			locals[222] = stack[0]
			stack[0] = rt.IntValue(222)
			locals[223] = stack[0]
			stack[0] = rt.IntValue(223)
			locals[224] = stack[0]
			stack[0] = rt.IntValue(224)
			locals[225] = stack[0]
			stack[0] = rt.IntValue(225)
			locals[226] = stack[0]
			stack[0] = rt.IntValue(226)
			locals[227] = stack[0]
			stack[0] = rt.IntValue(227)
			locals[228] = stack[0]
			stack[0] = rt.IntValue(228)
			locals[229] = stack[0]
			stack[0] = rt.IntValue(229)
			locals[230] = stack[0]
			stack[0] = rt.IntValue(230)
			locals[231] = stack[0]
			stack[0] = rt.IntValue(231)
			locals[232] = stack[0]
			stack[0] = rt.IntValue(232)
			locals[233] = stack[0]
			stack[0] = rt.IntValue(233)
			locals[234] = stack[0]
			stack[0] = rt.IntValue(234)
			locals[235] = stack[0]
			stack[0] = rt.IntValue(235)
			locals[236] = stack[0]
			stack[0] = rt.IntValue(236)
			locals[237] = stack[0]
			stack[0] = rt.IntValue(237)
			locals[238] = stack[0]
			stack[0] = rt.IntValue(238)
			locals[239] = stack[0]
			stack[0] = rt.IntValue(239)
			locals[240] = stack[0]
			stack[0] = rt.IntValue(240)
			locals[241] = stack[0]
			stack[0] = rt.IntValue(241)
			locals[242] = stack[0]
			stack[0] = rt.IntValue(242)
			locals[243] = stack[0]
			stack[0] = rt.IntValue(243)
			locals[244] = stack[0]
			stack[0] = rt.IntValue(244)
			locals[245] = stack[0]
			stack[0] = rt.IntValue(245)
			locals[246] = stack[0]
			stack[0] = rt.IntValue(246)
			locals[247] = stack[0]
			stack[0] = rt.IntValue(247)
			locals[248] = stack[0]
			stack[0] = rt.IntValue(248)
			locals[249] = stack[0]
			stack[0] = rt.IntValue(249)
			locals[250] = stack[0]
			stack[0] = rt.IntValue(250)
			locals[251] = stack[0]
			stack[0] = rt.IntValue(251)
			locals[252] = stack[0]
			stack[0] = rt.IntValue(252)
			locals[253] = stack[0]
			stack[0] = rt.IntValue(253)
			locals[254] = stack[0]
			stack[0] = rt.IntValue(254)
			locals[255] = stack[0]
			stack[0] = rt.IntValue(255)
			locals[256] = stack[0]
			stack[0] = rt.IntValue(256)
			locals[257] = stack[0]
			stack[0] = rt.IntValue(257)
			locals[258] = stack[0]
			stack[0] = rt.IntValue(258)
			locals[259] = stack[0]
			stack[0] = rt.IntValue(259)
			locals[260] = stack[0]
			stack[0] = rt.IntValue(260)
			locals[261] = stack[0]
			stack[0] = rt.IntValue(261)
			locals[262] = stack[0]
			stack[0] = rt.IntValue(262)
			locals[263] = stack[0]
			stack[0] = rt.IntValue(263)
			locals[264] = stack[0]
			stack[0] = rt.IntValue(264)
			locals[265] = stack[0]
			stack[0] = rt.IntValue(265)
			locals[266] = stack[0]
			stack[0] = rt.IntValue(266)
			locals[267] = stack[0]
			stack[0] = rt.IntValue(267)
			locals[268] = stack[0]
			stack[0] = rt.IntValue(268)
			locals[269] = stack[0]
			stack[0] = rt.IntValue(269)
			locals[270] = stack[0]
			stack[0] = rt.IntValue(270)
			locals[271] = stack[0]
			stack[0] = rt.IntValue(271)
			locals[272] = stack[0]
			stack[0] = rt.IntValue(272)
			locals[273] = stack[0]
			stack[0] = rt.IntValue(273)
			locals[274] = stack[0]
			stack[0] = rt.IntValue(274)
			locals[275] = stack[0]
			stack[0] = rt.IntValue(275)
			locals[276] = stack[0]
			stack[0] = rt.IntValue(276)
			locals[277] = stack[0]
			stack[0] = rt.IntValue(277)
			locals[278] = stack[0]
			stack[0] = rt.IntValue(278)
			locals[279] = stack[0]
			stack[0] = rt.IntValue(279)
			locals[280] = stack[0]
			stack[0] = rt.IntValue(280)
			locals[281] = stack[0]
			stack[0] = rt.IntValue(281)
			locals[282] = stack[0]
			stack[0] = rt.IntValue(282)
			locals[283] = stack[0]
			stack[0] = rt.IntValue(283)
			locals[284] = stack[0]
			stack[0] = rt.IntValue(284)
			locals[285] = stack[0]
			stack[0] = rt.IntValue(285)
			locals[286] = stack[0]
			stack[0] = rt.IntValue(286)
			locals[287] = stack[0]
			stack[0] = rt.IntValue(287)
			locals[288] = stack[0]
			stack[0] = rt.IntValue(288)
			locals[289] = stack[0]
			stack[0] = rt.IntValue(289)
			locals[290] = stack[0]
			stack[0] = rt.IntValue(290)
			locals[291] = stack[0]
			stack[0] = rt.IntValue(291)
			locals[292] = stack[0]
			stack[0] = rt.IntValue(292)
			locals[293] = stack[0]
			stack[0] = rt.IntValue(293)
			locals[294] = stack[0]
			stack[0] = rt.IntValue(294)
			locals[295] = stack[0]
			stack[0] = rt.IntValue(295)
			locals[296] = stack[0]
			stack[0] = rt.IntValue(296)
			locals[297] = stack[0]
			stack[0] = rt.IntValue(297)
			locals[298] = stack[0]
			stack[0] = rt.IntValue(298)
			locals[299] = stack[0]
			stack[0] = rt.IntValue(299)
			locals[300] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[1]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			stack[1] = locals[2]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[301] = stack[0]
			locals[269] = rt.IntValue(rt.AsInt(locals[269]) + 5)
			stack[0] = locals[301]
			stack[1] = locals[269]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[301] = stack[0]
			stack[0] = locals[301]
			stack[1] = locals[261]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[301] = stack[0]
			stack[0] = locals[299]
			stack[1] = rt.IntValue(1000)
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[299] = stack[0]
			stack[0] = locals[301]
			stack[1] = locals[299]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[301] = stack[0]
			stack[0] = locals[301]
			stack[1] = locals[300]
			stack[0] = rt.IntValue(rt.AsInt(stack[0]) + rt.AsInt(stack[1]))
			locals[301] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[301]
			{
			out := t.InvokeVirtual(m, 1539, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[301]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod12(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [264]rt.Value
	var stack [4]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = rt.LongBitsValue(0x0000000000001388)
			locals[2] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001389)
			locals[4] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138a)
			locals[6] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138b)
			locals[8] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138c)
			locals[10] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138d)
			locals[12] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138e)
			locals[14] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000138f)
			locals[16] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001390)
			locals[18] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001391)
			locals[20] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001392)
			locals[22] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001393)
			locals[24] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001394)
			locals[26] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001395)
			locals[28] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001396)
			locals[30] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001397)
			locals[32] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001398)
			locals[34] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001399)
			locals[36] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139a)
			locals[38] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139b)
			locals[40] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139c)
			locals[42] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139d)
			locals[44] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139e)
			locals[46] = stack[0]
			stack[0] = rt.LongBitsValue(0x000000000000139f)
			locals[48] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a0)
			locals[50] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a1)
			locals[52] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a2)
			locals[54] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a3)
			locals[56] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a4)
			locals[58] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a5)
			locals[60] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a6)
			locals[62] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a7)
			locals[64] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a8)
			locals[66] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013a9)
			locals[68] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013aa)
			locals[70] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ab)
			locals[72] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ac)
			locals[74] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ad)
			locals[76] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ae)
			locals[78] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013af)
			locals[80] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b0)
			locals[82] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b1)
			locals[84] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b2)
			locals[86] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b3)
			locals[88] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b4)
			locals[90] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b5)
			locals[92] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b6)
			locals[94] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b7)
			locals[96] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b8)
			locals[98] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013b9)
			locals[100] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ba)
			locals[102] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013bb)
			locals[104] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013bc)
			locals[106] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013bd)
			locals[108] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013be)
			locals[110] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013bf)
			locals[112] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c0)
			locals[114] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c1)
			locals[116] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c2)
			locals[118] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c3)
			locals[120] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c4)
			locals[122] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c5)
			locals[124] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c6)
			locals[126] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c7)
			locals[128] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c8)
			locals[130] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013c9)
			locals[132] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ca)
			locals[134] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013cb)
			locals[136] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013cc)
			locals[138] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013cd)
			locals[140] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ce)
			locals[142] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013cf)
			locals[144] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d0)
			locals[146] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d1)
			locals[148] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d2)
			locals[150] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d3)
			locals[152] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d4)
			locals[154] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d5)
			locals[156] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d6)
			locals[158] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d7)
			locals[160] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d8)
			locals[162] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013d9)
			locals[164] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013da)
			locals[166] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013db)
			locals[168] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013dc)
			locals[170] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013dd)
			locals[172] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013de)
			locals[174] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013df)
			locals[176] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e0)
			locals[178] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e1)
			locals[180] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e2)
			locals[182] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e3)
			locals[184] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e4)
			locals[186] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e5)
			locals[188] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e6)
			locals[190] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e7)
			locals[192] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e8)
			locals[194] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013e9)
			locals[196] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ea)
			locals[198] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013eb)
			locals[200] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ec)
			locals[202] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ed)
			locals[204] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ee)
			locals[206] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ef)
			locals[208] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f0)
			locals[210] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f1)
			locals[212] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f2)
			locals[214] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f3)
			locals[216] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f4)
			locals[218] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f5)
			locals[220] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f6)
			locals[222] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f7)
			locals[224] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f8)
			locals[226] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013f9)
			locals[228] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013fa)
			locals[230] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013fb)
			locals[232] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013fc)
			locals[234] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013fd)
			locals[236] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013fe)
			locals[238] = stack[0]
			stack[0] = rt.LongBitsValue(0x00000000000013ff)
			locals[240] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001400)
			locals[242] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001401)
			locals[244] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001402)
			locals[246] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001403)
			locals[248] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001404)
			locals[250] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001405)
			locals[252] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001406)
			locals[254] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001407)
			locals[256] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001408)
			locals[258] = stack[0]
			stack[0] = rt.LongBitsValue(0x0000000000001409)
			locals[260] = stack[0]
			stack[0] = locals[0]
			stack[1] = locals[258]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			stack[1] = locals[260]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[262] = stack[0]
			stack[0] = locals[256]
			stack[1] = rt.LongBitsValue(0x0000000000000007)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[256] = stack[0]
			stack[0] = locals[262]
			stack[1] = locals[256]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[262] = stack[0]
			stack[0] = locals[258]
			stack[1] = rt.LongBitsValue(0x000000000000000b)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[258] = stack[0]
			stack[0] = locals[262]
			stack[1] = locals[258]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[262] = stack[0]
			stack[0] = locals[260]
			stack[1] = rt.LongBitsValue(0x000000000000000d)
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[260] = stack[0]
			stack[0] = locals[262]
			stack[1] = locals[260]
			stack[0] = rt.LongValue(rt.AsLong(stack[0]) + rt.AsLong(stack[1]))
			locals[262] = stack[0]
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = locals[262]
			{
			out := t.InvokeVirtual(m, 752, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			stack[0] = locals[262]
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod13(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [1]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			stack[0] = locals[0]
			switch rt.AsInt(stack[0]) {
			case 1:
			block = 1
			continue
			case 2:
			block = 2
			continue
			default:
			block = 3
			continue
			}
		case 1:
			stack[0] = rt.IntValue(1)
			return rt.OkValue(stack[0])
		case 2:
			stack[0] = rt.IntValue(2)
			return rt.OkValue(stack[0])
		case 3:
			stack[0] = rt.IntValue(0)
			return rt.OkValue(stack[0])
		}
	}
}

func cattyMethod14(t *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
	var locals [1]rt.Value
	var stack [3]rt.Value
	_ = locals
	_ = stack
	locals[0] = args[0]
	block := 0
	for {
		switch block {
		case 0:
			{
			v, thrown := t.GetStatic(m, "java/lang/System", "out", "Ljava/io/PrintStream;")
			if thrown != nil {
			return rt.ThrowOutcome(thrown)
			}
			stack[0] = v
			}
			stack[1] = rt.IntValue(3)
			{
			out := t.InvokeStatic(m, 4, "SwitchWide", "dense", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 7, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(7)
			{
			out := t.InvokeStatic(m, 15, "SwitchWide", "dense", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 18, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(11)
			{
			out := t.InvokeStatic(m, 26, "SwitchWide", "dense", "(I)I", stack[1:2])
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
			stack[1] = rt.IntValue(15)
			{
			out := t.InvokeStatic(m, 37, "SwitchWide", "offsetDense", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 40, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			out := t.InvokeStatic(m, 47, "SwitchWide", "offsetDense", "(I)I", stack[1:2])
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
			stack[1] = rt.IntValue(1000000)
			{
			out := t.InvokeStatic(m, 59, "SwitchWide", "sparse", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 62, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[2] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 70, "SwitchWide", "nested", "(II)I", stack[1:3])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 73, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(9)
			{
			out := t.InvokeStatic(m, 81, "SwitchWide", "loopSwitch", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 84, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(97)
			{
			out := t.InvokeStatic(m, 92, "SwitchWide", "charSwitch", "(C)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 95, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(113)
			{
			out := t.InvokeStatic(m, 103, "SwitchWide", "charSwitch", "(C)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 106, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(-2)
			{
			out := t.InvokeStatic(m, 114, "SwitchWide", "byteSwitch", "(B)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 117, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(9)
			{
			out := t.InvokeStatic(m, 125, "SwitchWide", "byteSwitch", "(B)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 128, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(-300)
			{
			out := t.InvokeStatic(m, 137, "SwitchWide", "shortSwitch", "(S)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 140, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(301)
			{
			out := t.InvokeStatic(m, 149, "SwitchWide", "shortSwitch", "(S)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 152, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(5)
			{
			out := t.InvokeStatic(m, 159, "SwitchWide", "fallthrough", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 162, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(4)
			{
			out := t.InvokeStatic(m, 169, "SwitchWide", "fallthrough", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 172, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(7)
			{
			out := t.InvokeStatic(m, 180, "SwitchWide", "fallthrough", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 183, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(9)
			{
			out := t.InvokeStatic(m, 191, "SwitchWide", "fallthrough", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 194, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			out := t.InvokeStatic(m, 201, "SwitchWide", "noDefault", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 204, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			out := t.InvokeStatic(m, 211, "SwitchWide", "noDefault", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 214, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(2)
			{
			out := t.InvokeStatic(m, 221, "SwitchWide", "noDefault", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 224, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(20000)
			{
			out := t.InvokeStatic(m, 233, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 236, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(700)
			{
			out := t.InvokeStatic(m, 245, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 248, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(20)
			{
			out := t.InvokeStatic(m, 256, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 259, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(5)
			{
			out := t.InvokeStatic(m, 266, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 269, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			out := t.InvokeStatic(m, 276, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 279, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(-3)
			{
			out := t.InvokeStatic(m, 287, "SwitchWide", "deepIf", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 290, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.IntValue(7)
			{
			out := t.InvokeStatic(m, 298, "SwitchWide", "wideInts", "(I)I", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 301, stack[0], "java/io/PrintStream", "println", "(I)V", stack[1:2])
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
			stack[1] = rt.LongBitsValue(0x0000000000000009)
			{
			out := t.InvokeStatic(m, 310, "SwitchWide", "wideLongs", "(J)J", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			stack[1] = out.Value
			}
			{
			out := t.InvokeVirtual(m, 313, stack[0], "java/io/PrintStream", "println", "(J)V", stack[1:2])
			if out.Thrown != nil {
			thrown := out.Thrown
			return rt.ThrowOutcome(thrown)
			}
			}
			return rt.Normal()
		}
	}
}

