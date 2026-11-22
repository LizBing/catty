package rt

// JVM opcode values for the Gate 1 profile.
const (
	OpAconstNull      = 0x01
	OpIconstM1        = 0x02
	OpIconst0         = 0x03
	OpIconst1         = 0x04
	OpIconst2         = 0x05
	OpIconst3         = 0x06
	OpIconst4         = 0x07
	OpIconst5         = 0x08
	OpBipush          = 0x10
	OpSipush          = 0x11
	OpLdc             = 0x12
	OpLdcW            = 0x13
	OpIload           = 0x15
	OpAload           = 0x19
	OpIstore          = 0x36
	OpAstore          = 0x3a
	OpPop             = 0x57
	OpDup             = 0x59
	OpIadd            = 0x60
	OpIrem            = 0x70
	OpIreturn         = 0xac
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
	OpAnewarray       = 0xbd
	OpArraylength     = 0xbe
	OpAthrow          = 0xbf
	OpIfeq            = 0x99
	OpIfne            = 0x9a
	OpIfIcmpeq        = 0x9f
	OpIfIcmpne        = 0xa0
	OpIfAcmpeq        = 0xa5
	OpIfAcmpne        = 0xa6
	OpGoto            = 0xa7
	OpAaload          = 0x32
	OpAastore         = 0x53
)

func localIndex(op uint8, operand int32) (int32, bool) {
	switch op {
	case OpIload, OpAload, OpIstore, OpAstore:
		return operand, true
	case 0x1a, 0x1b, 0x1c, 0x1d: // iload_0..3
		return int32(op - 0x1a), true
	case 0x2a, 0x2b, 0x2c, 0x2d: // aload_0..3
		return int32(op - 0x2a), true
	case 0x3b, 0x3c, 0x3d, 0x3e: // istore_0..3
		return int32(op - 0x3b), true
	case 0x4b, 0x4c, 0x4d, 0x4e: // astore_0..3
		return int32(op - 0x4b), true
	}
	return 0, false
}
