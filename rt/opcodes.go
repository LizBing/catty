package rt

import "sort"

// JVM opcode values admitted by the Gate 2 profile.
const (
	OpNop             = 0x00
	OpAconstNull      = 0x01
	OpIconstM1        = 0x02
	OpIconst0         = 0x03
	OpIconst1         = 0x04
	OpIconst2         = 0x05
	OpIconst3         = 0x06
	OpIconst4         = 0x07
	OpIconst5         = 0x08
	OpLconst0         = 0x09
	OpLconst1         = 0x0a
	OpFconst0         = 0x0b
	OpFconst1         = 0x0c
	OpFconst2         = 0x0d
	OpDconst0         = 0x0e
	OpDconst1         = 0x0f
	OpBipush          = 0x10
	OpSipush          = 0x11
	OpLdc             = 0x12
	OpLdcW            = 0x13
	OpLdc2W           = 0x14
	OpIload           = 0x15
	OpLload           = 0x16
	OpFload           = 0x17
	OpDload           = 0x18
	OpAload           = 0x19
	OpIaload          = 0x2e
	OpLaload          = 0x2f
	OpFaload          = 0x30
	OpDaload          = 0x31
	OpAaload          = 0x32
	OpBaload          = 0x33
	OpCaload          = 0x34
	OpSaload          = 0x35
	OpIstore          = 0x36
	OpLstore          = 0x37
	OpFstore          = 0x38
	OpDstore          = 0x39
	OpAstore          = 0x3a
	OpIastore         = 0x4f
	OpLastore         = 0x50
	OpFastore         = 0x51
	OpDastore         = 0x52
	OpAastore         = 0x53
	OpBastore         = 0x54
	OpCastore         = 0x55
	OpSastore         = 0x56
	OpPop             = 0x57
	OpPop2            = 0x58
	OpDup             = 0x59
	OpDupX1           = 0x5a
	OpDupX2           = 0x5b
	OpDup2            = 0x5c
	OpDup2X1          = 0x5d
	OpDup2X2          = 0x5e
	OpSwap            = 0x5f
	OpIadd            = 0x60
	OpLadd            = 0x61
	OpFadd            = 0x62
	OpDadd            = 0x63
	OpIsub            = 0x64
	OpLsub            = 0x65
	OpFsub            = 0x66
	OpDsub            = 0x67
	OpImul            = 0x68
	OpLmul            = 0x69
	OpFmul            = 0x6a
	OpDmul            = 0x6b
	OpIdiv            = 0x6c
	OpLdiv            = 0x6d
	OpFdiv            = 0x6e
	OpDdiv            = 0x6f
	OpIrem            = 0x70
	OpLrem            = 0x71
	OpFrem            = 0x72
	OpDrem            = 0x73
	OpIneg            = 0x74
	OpLneg            = 0x75
	OpFneg            = 0x76
	OpDneg            = 0x77
	OpIshl            = 0x78
	OpLshl            = 0x79
	OpIshr            = 0x7a
	OpLshr            = 0x7b
	OpIushr           = 0x7c
	OpLushr           = 0x7d
	OpIand            = 0x7e
	OpLand            = 0x7f
	OpIor             = 0x80
	OpLor             = 0x81
	OpIxor            = 0x82
	OpLxor            = 0x83
	OpIinc            = 0x84
	OpI2l             = 0x85
	OpI2f             = 0x86
	OpI2d             = 0x87
	OpL2i             = 0x88
	OpL2f             = 0x89
	OpL2d             = 0x8a
	OpF2i             = 0x8b
	OpF2l             = 0x8c
	OpF2d             = 0x8d
	OpD2i             = 0x8e
	OpD2l             = 0x8f
	OpD2f             = 0x90
	OpI2b             = 0x91
	OpI2c             = 0x92
	OpI2s             = 0x93
	OpLcmp            = 0x94
	OpFcmpl           = 0x95
	OpFcmpg           = 0x96
	OpDcmpl           = 0x97
	OpDcmpg           = 0x98
	OpIfeq            = 0x99
	OpIfne            = 0x9a
	OpIflt            = 0x9b
	OpIfge            = 0x9c
	OpIfgt            = 0x9d
	OpIfle            = 0x9e
	OpIfIcmpeq        = 0x9f
	OpIfIcmpne        = 0xa0
	OpIfIcmplt        = 0xa1
	OpIfIcmpge        = 0xa2
	OpIfIcmpgt        = 0xa3
	OpIfIcmple        = 0xa4
	OpIfAcmpeq        = 0xa5
	OpIfAcmpne        = 0xa6
	OpGoto            = 0xa7
	OpJsr             = 0xa8
	OpRet             = 0xa9
	OpTableswitch     = 0xaa
	OpLookupswitch    = 0xab
	OpIreturn         = 0xac
	OpLreturn         = 0xad
	OpFreturn         = 0xae
	OpDreturn         = 0xaf
	OpAreturn         = 0xb0
	OpReturn          = 0xb1
	OpGetstatic       = 0xb2
	OpPutstatic       = 0xb3
	OpGetfield        = 0xb4
	OpPutfield        = 0xb5
	OpInvokevirtual   = 0xb6
	OpInvokespecial   = 0xb7
	OpInvokestatic    = 0xb8
	OpInvokeinterface = 0xb9
	OpInvokedynamic   = 0xba
	OpNew             = 0xbb
	OpNewarray        = 0xbc
	OpAnewarray       = 0xbd
	OpArraylength     = 0xbe
	OpAthrow          = 0xbf
	OpCheckcast       = 0xc0
	OpInstanceof      = 0xc1
	OpMonitorenter    = 0xc2
	OpMonitorexit     = 0xc3
	OpWide            = 0xc4
	OpMultianewarray  = 0xc5
	OpIfnull          = 0xc6
	OpIfnonnull       = 0xc7
	OpGotoW           = 0xc8
	OpJsrW            = 0xc9
	OpBreakpoint      = 0xca
	OpImpdep1         = 0xfe
	OpImpdep2         = 0xff
)

// opcodeNames maps every admitted opcode to its JVM mnemonic.
var opcodeNames = map[uint8]string{
	OpNop: "nop", OpAconstNull: "aconst_null",
	OpIconstM1: "iconst_m1", OpIconst0: "iconst_0", OpIconst1: "iconst_1",
	OpIconst2: "iconst_2", OpIconst3: "iconst_3", OpIconst4: "iconst_4", OpIconst5: "iconst_5",
	OpLconst0: "lconst_0", OpLconst1: "lconst_1",
	OpFconst0: "fconst_0", OpFconst1: "fconst_1", OpFconst2: "fconst_2",
	OpDconst0: "dconst_0", OpDconst1: "dconst_1",
	OpBipush: "bipush", OpSipush: "sipush",
	OpLdc: "ldc", OpLdcW: "ldc_w", OpLdc2W: "ldc2_w",
	OpIload: "iload", OpLload: "lload", OpFload: "fload", OpDload: "dload", OpAload: "aload",
	0x1a: "iload_0", 0x1b: "iload_1", 0x1c: "iload_2", 0x1d: "iload_3",
	0x1e: "lload_0", 0x1f: "lload_1", 0x20: "lload_2", 0x21: "lload_3",
	0x22: "fload_0", 0x23: "fload_1", 0x24: "fload_2", 0x25: "fload_3",
	0x26: "dload_0", 0x27: "dload_1", 0x28: "dload_2", 0x29: "dload_3",
	0x2a: "aload_0", 0x2b: "aload_1", 0x2c: "aload_2", 0x2d: "aload_3",
	OpIaload: "iaload", OpLaload: "laload", OpFaload: "faload", OpDaload: "daload",
	OpAaload: "aaload", OpBaload: "baload", OpCaload: "caload", OpSaload: "saload",
	OpIstore: "istore", OpLstore: "lstore", OpFstore: "fstore", OpDstore: "dstore", OpAstore: "astore",
	0x3b: "istore_0", 0x3c: "istore_1", 0x3d: "istore_2", 0x3e: "istore_3",
	0x3f: "lstore_0", 0x40: "lstore_1", 0x41: "lstore_2", 0x42: "lstore_3",
	0x43: "fstore_0", 0x44: "fstore_1", 0x45: "fstore_2", 0x46: "fstore_3",
	0x47: "dstore_0", 0x48: "dstore_1", 0x49: "dstore_2", 0x4a: "dstore_3",
	0x4b: "astore_0", 0x4c: "astore_1", 0x4d: "astore_2", 0x4e: "astore_3",
	OpIastore: "iastore", OpLastore: "lastore", OpFastore: "fastore", OpDastore: "dastore",
	OpAastore: "aastore", OpBastore: "bastore", OpCastore: "castore", OpSastore: "sastore",
	OpPop: "pop", OpPop2: "pop2", OpDup: "dup", OpDupX1: "dup_x1", OpDupX2: "dup_x2",
	OpDup2: "dup2", OpDup2X1: "dup2_x1", OpDup2X2: "dup2_x2", OpSwap: "swap",
	OpIadd: "iadd", OpLadd: "ladd", OpFadd: "fadd", OpDadd: "dadd",
	OpIsub: "isub", OpLsub: "lsub", OpFsub: "fsub", OpDsub: "dsub",
	OpImul: "imul", OpLmul: "lmul", OpFmul: "fmul", OpDmul: "dmul",
	OpIdiv: "idiv", OpLdiv: "ldiv", OpFdiv: "fdiv", OpDdiv: "ddiv",
	OpIrem: "irem", OpLrem: "lrem", OpFrem: "frem", OpDrem: "drem",
	OpIneg: "ineg", OpLneg: "lneg", OpFneg: "fneg", OpDneg: "dneg",
	OpIshl: "ishl", OpLshl: "lshl", OpIshr: "ishr", OpLshr: "lshr",
	OpIushr: "iushr", OpLushr: "lushr",
	OpIand: "iand", OpLand: "land", OpIor: "ior", OpLor: "lor", OpIxor: "ixor", OpLxor: "lxor",
	OpIinc: "iinc",
	OpI2l: "i2l", OpI2f: "i2f", OpI2d: "i2d", OpL2i: "l2i", OpL2f: "l2f", OpL2d: "l2d",
	OpF2i: "f2i", OpF2l: "f2l", OpF2d: "f2d", OpD2i: "d2i", OpD2l: "d2l", OpD2f: "d2f",
	OpI2b: "i2b", OpI2c: "i2c", OpI2s: "i2s",
	OpLcmp: "lcmp", OpFcmpl: "fcmpl", OpFcmpg: "fcmpg", OpDcmpl: "dcmpl", OpDcmpg: "dcmpg",
	OpIfeq: "ifeq", OpIfne: "ifne", OpIflt: "iflt", OpIfge: "ifge", OpIfgt: "ifgt", OpIfle: "ifle",
	OpIfIcmpeq: "if_icmpeq", OpIfIcmpne: "if_icmpne", OpIfIcmplt: "if_icmplt",
	OpIfIcmpge: "if_icmpge", OpIfIcmpgt: "if_icmpgt", OpIfIcmple: "if_icmple",
	OpIfAcmpeq: "if_acmpeq", OpIfAcmpne: "if_acmpne",
	OpGoto: "goto", OpJsr: "jsr", OpRet: "ret",
	OpTableswitch: "tableswitch", OpLookupswitch: "lookupswitch",
	OpIreturn: "ireturn", OpLreturn: "lreturn", OpFreturn: "freturn", OpDreturn: "dreturn",
	OpAreturn: "areturn", OpReturn: "return",
	OpGetstatic: "getstatic", OpPutstatic: "putstatic", OpGetfield: "getfield", OpPutfield: "putfield",
	OpInvokevirtual: "invokevirtual", OpInvokespecial: "invokespecial",
	OpInvokestatic: "invokestatic", OpInvokeinterface: "invokeinterface", OpInvokedynamic: "invokedynamic",
	OpNew: "new", OpNewarray: "newarray", OpAnewarray: "anewarray",
	OpArraylength: "arraylength", OpAthrow: "athrow",
	OpCheckcast: "checkcast", OpInstanceof: "instanceof",
	OpMonitorenter: "monitorenter", OpMonitorexit: "monitorexit",
	OpWide: "wide", OpMultianewarray: "multianewarray",
	OpIfnull: "ifnull", OpIfnonnull: "ifnonnull", OpGotoW: "goto_w", OpJsrW: "jsr_w",
	OpBreakpoint: "breakpoint", OpImpdep1: "impdep1", OpImpdep2: "impdep2",
}

// OpcodeName returns the JVM mnemonic for an admitted opcode.
func OpcodeName(op uint8) string {
	return opcodeNames[op]
}

// SupportedOpcodeNames returns the sorted opcode inventory published by the
// profile. Rejected opcodes (jsr, ret, invokedynamic, monitorenter/exit,
// reserved opcodes) are deliberately absent.
func SupportedOpcodeNames() []string {
	out := make([]string, 0, len(opcodeNames))
	for op, name := range opcodeNames {
		if !AdmitsOpcode(op) {
			continue
		}
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// AdmitsOpcode reports whether the profile admits op.
func AdmitsOpcode(op uint8) bool {
	switch op {
	case OpJsr, OpRet, OpInvokedynamic, OpMonitorenter, OpMonitorexit,
		OpJsrW, OpBreakpoint, OpImpdep1, OpImpdep2:
		return false
	}
	_, ok := opcodeNames[op]
	return ok
}

// LocalIndex decodes the effective local index of a local variable
// instruction. It returns false for opcodes that are not local variable
// instructions.
func LocalIndex(op uint8, operand int32) (int32, bool) {
	switch op {
	case OpIload, OpLload, OpFload, OpDload, OpAload,
		OpIstore, OpLstore, OpFstore, OpDstore, OpAstore:
		return operand, true
	case OpIinc:
		return operand, true
	case 0x1a, 0x1b, 0x1c, 0x1d:
		return int32(op - 0x1a), true
	case 0x1e, 0x1f, 0x20, 0x21:
		return int32(op - 0x1e), true
	case 0x22, 0x23, 0x24, 0x25:
		return int32(op - 0x22), true
	case 0x26, 0x27, 0x28, 0x29:
		return int32(op - 0x26), true
	case 0x2a, 0x2b, 0x2c, 0x2d:
		return int32(op - 0x2a), true
	case 0x3b, 0x3c, 0x3d, 0x3e:
		return int32(op - 0x3b), true
	case 0x3f, 0x40, 0x41, 0x42:
		return int32(op - 0x3f), true
	case 0x43, 0x44, 0x45, 0x46:
		return int32(op - 0x43), true
	case 0x47, 0x48, 0x49, 0x4a:
		return int32(op - 0x47), true
	case 0x4b, 0x4c, 0x4d, 0x4e:
		return int32(op - 0x4b), true
	}
	return 0, false
}

// LoadWidth reports the slot width pushed by a local load or constant
// instruction, and false if op is not one of them.
func LoadWidth(op uint8) (int, bool) {
	switch op {
	case OpLload, OpDload, OpLconst0, OpLconst1, OpDconst0, OpDconst1, OpLdc2W:
		return 2, true
	case OpIload, OpFload, OpAload, OpIconstM1, OpIconst0, OpIconst1, OpIconst2,
		OpIconst3, OpIconst4, OpIconst5, OpFconst0, OpFconst1, OpFconst2,
		OpBipush, OpSipush, OpLdc, OpLdcW, OpAconstNull:
		return 1, true
	case 0x1a, 0x1b, 0x1c, 0x1d, 0x22, 0x23, 0x24, 0x25, 0x2a, 0x2b, 0x2c, 0x2d:
		return 1, true
	case 0x1e, 0x1f, 0x20, 0x21, 0x26, 0x27, 0x28, 0x29:
		return 2, true
	}
	return 0, false
}

// NewarrayDescriptor maps a `newarray` atype to its array descriptor.
func NewarrayDescriptor(atype int32) (string, bool) {
	switch atype {
	case 4:
		return "[Z", true
	case 5:
		return "[C", true
	case 6:
		return "[F", true
	case 7:
		return "[D", true
	case 8:
		return "[B", true
	case 9:
		return "[S", true
	case 10:
		return "[I", true
	case 11:
		return "[J", true
	}
	return "", false
}
