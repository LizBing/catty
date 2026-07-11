// Package opcode holds the JVM opcode set (JVMS §6.5) as a typed constant
// block. It is a leaf package imported by every stage that switches on bytecode
// (the interpreter and the lowering pass), so neither has to depend on the
// other just to share opcode values.
package opcode

// Opcode is a single JVM bytecode opcode.
type Opcode byte

// JVM opcode mnemonics; values are the opcode bytes. Constants catty does not
// yet implement (invokedynamic, jsr, ...) are present too so a dispatcher can
// name them precisely when it panics.
const (
	Nop             Opcode = 0x00
	AconstNull      Opcode = 0x01
	IconstM1        Opcode = 0x02
	Iconst0         Opcode = 0x03
	Iconst1         Opcode = 0x04
	Iconst2         Opcode = 0x05
	Iconst3         Opcode = 0x06
	Iconst4         Opcode = 0x07
	Iconst5         Opcode = 0x08
	Lconst0         Opcode = 0x09
	Lconst1         Opcode = 0x0a
	Fconst0         Opcode = 0x0b
	Fconst1         Opcode = 0x0c
	Fconst2         Opcode = 0x0d
	Dconst0         Opcode = 0x0e
	Dconst1         Opcode = 0x0f
	Bipush          Opcode = 0x10
	Sipush          Opcode = 0x11
	Ldc             Opcode = 0x12
	LdcW            Opcode = 0x13
	Ldc2W           Opcode = 0x14
	Iload           Opcode = 0x15
	Lload           Opcode = 0x16
	Fload           Opcode = 0x17
	Dload           Opcode = 0x18
	Aload           Opcode = 0x19
	Iload0          Opcode = 0x1a
	Iload1          Opcode = 0x1b
	Iload2          Opcode = 0x1c
	Iload3          Opcode = 0x1d
	Lload0          Opcode = 0x1e
	Lload1          Opcode = 0x1f
	Lload2          Opcode = 0x20
	Lload3          Opcode = 0x21
	Fload0          Opcode = 0x22
	Fload1          Opcode = 0x23
	Fload2          Opcode = 0x24
	Fload3          Opcode = 0x25
	Dload0          Opcode = 0x26
	Dload1          Opcode = 0x27
	Dload2          Opcode = 0x28
	Dload3          Opcode = 0x29
	Aload0          Opcode = 0x2a
	Aload1          Opcode = 0x2b
	Aload2          Opcode = 0x2c
	Aload3          Opcode = 0x2d
	Iaload          Opcode = 0x2e
	Laload          Opcode = 0x2f
	Faload          Opcode = 0x30
	Daload          Opcode = 0x31
	Aaload          Opcode = 0x32
	Baload          Opcode = 0x33
	Caload          Opcode = 0x34
	Saload          Opcode = 0x35
	Istore          Opcode = 0x36
	Lstore          Opcode = 0x37
	Fstore          Opcode = 0x38
	Dstore          Opcode = 0x39
	Astore          Opcode = 0x3a
	Istore0         Opcode = 0x3b
	Istore1         Opcode = 0x3c
	Istore2         Opcode = 0x3d
	Istore3         Opcode = 0x3e
	Lstore0         Opcode = 0x3f
	Lstore1         Opcode = 0x40
	Lstore2         Opcode = 0x41
	Lstore3         Opcode = 0x42
	Fstore0         Opcode = 0x43
	Fstore1         Opcode = 0x44
	Fstore2         Opcode = 0x45
	Fstore3         Opcode = 0x46
	Dstore0         Opcode = 0x47
	Dstore1         Opcode = 0x48
	Dstore2         Opcode = 0x49
	Dstore3         Opcode = 0x4a
	Astore0         Opcode = 0x4b
	Astore1         Opcode = 0x4c
	Astore2         Opcode = 0x4d
	Astore3         Opcode = 0x4e
	Iastore         Opcode = 0x4f
	Lastore         Opcode = 0x50
	Fastore         Opcode = 0x51
	Dastore         Opcode = 0x52
	Aastore         Opcode = 0x53
	Bastore         Opcode = 0x54
	Castore         Opcode = 0x55
	Sastore         Opcode = 0x56
	Pop             Opcode = 0x57
	Pop2            Opcode = 0x58
	Dup             Opcode = 0x59
	DupX1           Opcode = 0x5a
	DupX2           Opcode = 0x5b
	Dup2            Opcode = 0x5c
	Dup2X1          Opcode = 0x5d
	Dup2X2          Opcode = 0x5e
	Swap            Opcode = 0x5f
	Iadd            Opcode = 0x60
	Ladd            Opcode = 0x61
	Fadd            Opcode = 0x62
	Dadd            Opcode = 0x63
	Isub            Opcode = 0x64
	Lsub            Opcode = 0x65
	Fsub            Opcode = 0x66
	Dsub            Opcode = 0x67
	Imul            Opcode = 0x68
	Lmul            Opcode = 0x69
	Fmul            Opcode = 0x6a
	Dmul            Opcode = 0x6b
	Idiv            Opcode = 0x6c
	Ldiv            Opcode = 0x6d
	Fdiv            Opcode = 0x6e
	Ddiv            Opcode = 0x6f
	Irem            Opcode = 0x70
	Lrem            Opcode = 0x71
	Frem            Opcode = 0x72
	Drem            Opcode = 0x73
	Ineg            Opcode = 0x74
	Lneg            Opcode = 0x75
	Fneg            Opcode = 0x76
	Dneg            Opcode = 0x77
	Ishl            Opcode = 0x78
	Lshl            Opcode = 0x79
	Ishr            Opcode = 0x7a
	Lshr            Opcode = 0x7b
	Iushr           Opcode = 0x7c
	Lushr           Opcode = 0x7d
	Iand            Opcode = 0x7e
	Land            Opcode = 0x7f
	Ior             Opcode = 0x80
	Lor             Opcode = 0x81
	Ixor            Opcode = 0x82
	Lxor            Opcode = 0x83
	Iinc            Opcode = 0x84
	I2l             Opcode = 0x85
	I2f             Opcode = 0x86
	I2d             Opcode = 0x87
	L2i             Opcode = 0x88
	L2f             Opcode = 0x89
	L2d             Opcode = 0x8a
	F2i             Opcode = 0x8b
	F2l             Opcode = 0x8c
	F2d             Opcode = 0x8d
	D2i             Opcode = 0x8e
	D2l             Opcode = 0x8f
	D2f             Opcode = 0x90
	I2b             Opcode = 0x91
	I2c             Opcode = 0x92
	I2s             Opcode = 0x93
	Lcmp            Opcode = 0x94
	Fcmpl           Opcode = 0x95
	Fcmpg           Opcode = 0x96
	Dcmpl           Opcode = 0x97
	Dcmpg           Opcode = 0x98
	Ifeq            Opcode = 0x99
	Ifne            Opcode = 0x9a
	Iflt            Opcode = 0x9b
	Ifge            Opcode = 0x9c
	Ifgt            Opcode = 0x9d
	Ifle            Opcode = 0x9e
	IfIcmpeq        Opcode = 0x9f
	IfIcmpne        Opcode = 0xa0
	IfIcmplt        Opcode = 0xa1
	IfIcmpge        Opcode = 0xa2
	IfIcmpgt        Opcode = 0xa3
	IfIcmple        Opcode = 0xa4
	IfAcmpeq        Opcode = 0xa5
	IfAcmpne        Opcode = 0xa6
	Goto            Opcode = 0xa7
	Tableswitch     Opcode = 0xaa
	Lookupswitch    Opcode = 0xab
	Ireturn         Opcode = 0xac
	Lreturn         Opcode = 0xad
	Freturn         Opcode = 0xae
	Dreturn         Opcode = 0xaf
	Areturn         Opcode = 0xb0
	Return          Opcode = 0xb1
	Getstatic       Opcode = 0xb2
	Putstatic       Opcode = 0xb3
	Getfield        Opcode = 0xb4
	Putfield        Opcode = 0xb5
	Invokevirtual   Opcode = 0xb6
	Invokespecial   Opcode = 0xb7
	Invokestatic    Opcode = 0xb8
	Invokeinterface Opcode = 0xb9
	Invokedynamic   Opcode = 0xba
	New             Opcode = 0xbb
	Newarray        Opcode = 0xbc
	Anewarray       Opcode = 0xbd
	Arraylength     Opcode = 0xbe
	Athrow          Opcode = 0xbf
	Checkcast       Opcode = 0xc0
	Instanceof      Opcode = 0xc1
	Monitorenter    Opcode = 0xc2
	Monitorexit     Opcode = 0xc3
	Wide            Opcode = 0xc4
	Multianewarray  Opcode = 0xc5
	Ifnull          Opcode = 0xc6
	Ifnonnull       Opcode = 0xc7
	GotoW           Opcode = 0xc8
)

// Name returns a human-readable mnemonic for opcodes catty declines to run, for
// use in panic messages. Implemented opcodes need not be named here.
func Name(o Opcode) string {
	switch o {
	case Invokedynamic:
		return "invokedynamic"
	case Invokeinterface:
		return "invokeinterface"
	case Multianewarray:
		return "multianewarray"
	case Wide:
		return "wide"
	default:
		return "unknown"
	}
}
