package classfile

// OpcodeStackEffect is the SINGLE SOURCE OF TRUTH for the numeric operand
// stack effect of every JVM opcode (raw slots; a category-2 value counts
// as two). It exists because three engines previously grew private tables
// that drifted apart (DEBT-0015 interpreter pop/cat2 rows, DEBT-0017
// invoke-cat2 return depth, DEBT-0019 emitter array stores) — each drift
// surfaced as runtime corruption instead of a test failure.
//
// Contract: classification is total over all 256 byte values.
//
//	EffectFixed      → slots is the exact net effect.
//	EffectDescriptor → effect depends on constant-pool operands
//	                   (invokes, field access, multianewarray, wide);
//	                   callers resolve them from the pool.
//	EffectIllegal    → not reachable in a v52 class file accepted by this
//	                   project (jsr/ret outlawed by project rule;
//	                   invokedynamic desugared at load; reserved bytes).
//
// Exhaustiveness is enforced by TestOpcodeStackEffectTotal.
type EffectClass uint8

const (
	EffectFixed EffectClass = iota
	EffectDescriptor
	EffectIllegal
)

// StackEffect classifies op. For EffectFixed it also returns the net raw
// slot delta; other classes return 0.
func StackEffect(op byte) (slots int, class EffectClass) {
	switch op {
	case 0x00: // nop
		return 0, EffectFixed

	case 0x01: // aconst_null
		return 1, EffectFixed
	case 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08: // iconst_m1..iconst_5
		return 1, EffectFixed
	case 0x09, 0x0a: // lconst_0, lconst_1
		return 2, EffectFixed
	case 0x0b, 0x0c, 0x0d: // fconst_0..2
		return 1, EffectFixed
	case 0x0e, 0x0f: // dconst_0, dconst_1
		return 2, EffectFixed
	case 0x10, 0x11: // bipush, sipush
		return 1, EffectFixed
	case 0x12, 0x13: // ldc, ldc_w
		return 1, EffectFixed
	case 0x14: // ldc2_w
		return 2, EffectFixed

	case 0x15, 0x17, 0x19, // iload, fload, aload
		0x1a, 0x1b, 0x1c, 0x1d, // iload_0..3
		0x22, 0x23, 0x24, 0x25, // fload_0..3
		0x2a, 0x2b, 0x2c, 0x2d: // aload_0..3
		return 1, EffectFixed
	case 0x16, 0x18, // lload, dload
		0x1e, 0x1f, 0x20, 0x21, // lload_0..3
		0x26, 0x27, 0x28, 0x29: // dload_0..3
		return 2, EffectFixed

	case 0x2e, 0x30, 0x32, 0x33, 0x34, 0x35: // iaload faload aaload baload caload saload
		return -1, EffectFixed
	case 0x2f, 0x31: // laload, daload
		return 0, EffectFixed

	case 0x36, 0x38, 0x3a, // istore, fstore, astore
		0x3b, 0x3c, 0x3d, 0x3e, // istore_0..3
		0x43, 0x44, 0x45, 0x46, // fstore_0..3
		0x4b, 0x4c, 0x4d, 0x4e, // astore_0..3
		0x57: // pop
		return -1, EffectFixed
	case 0x37, 0x39, // lstore, dstore
		0x3f, 0x40, 0x41, 0x42, // lstore_0..3
		0x47, 0x48, 0x49, 0x4a, // dstore_0..3
		0x58: // pop2
		return -2, EffectFixed

	case 0x4f, 0x51, 0x53, 0x54, 0x55, 0x56: // iastore fastore aastore bastore castore sastore
		return -3, EffectFixed
	case 0x50, 0x52: // lastore, dastore
		return -4, EffectFixed

	case 0x59, 0x5a, 0x5b: // dup, dup_x1, dup_x2
		return 1, EffectFixed
	case 0x5c, 0x5d, 0x5e: // dup2, dup2_x1, dup2_x2
		return 2, EffectFixed
	case 0x5f: // swap
		return 0, EffectFixed

	case 0x60, 0x62, 0x64, 0x66, 0x68, 0x6a, 0x6c, 0x6e, 0x70, 0x72, // iadd fadd isub fsub imul fmul idiv fdiv irem frem
		0x74, 0x76, // ineg, fneg
		0x78, 0x7a, 0x7c, // ishl, ishr, iushr
		0x7e, 0x80, 0x82: // iand, ior, ixor
		return -1, EffectFixed

	case 0x61, 0x63, 0x65, 0x67, 0x69, 0x6b, 0x6d, 0x6f, 0x71, 0x73: // ladd dadd lsub dsub lmul dmul ldiv ddiv lrem drem
		return -2, EffectFixed
	case 0x75, 0x77: // lneg, dneg
		return 0, EffectFixed

	case 0x79, 0x7b, 0x7d: // lshl, lshr, lushr
		return -1, EffectFixed

	case 0x85, 0x87, 0x8c, 0x8d: // i2l, i2d, f2l, f2d
		return 1, EffectFixed
	case 0x86, 0x8b, 0x91, 0x92, 0x93, // i2f, f2i, i2b, i2c, i2s
		0x8a, 0x8f: // l2d, d2l
		return 0, EffectFixed
	case 0x88, 0x89, 0x8e, 0x90: // l2i, l2f, d2i, d2f
		return -1, EffectFixed

	case 0x94, 0x97, 0x98: // lcmp, dcmpl, dcmpg
		return -3, EffectFixed
	case 0x95, 0x96: // fcmpl, fcmpg
		return -1, EffectFixed

	case 0x84: // iinc
		return 0, EffectFixed

	case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, // ifeq..ifle
		0xc6, 0xc7: // ifnull, ifnonnull
		return -1, EffectFixed
	case 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6: // if_icmpXX, if_acmpeq/ne
		return -2, EffectFixed
	case 0xa7, 0xc8: // goto, goto_w
		return 0, EffectFixed

	case 0xaa, 0xab: // tableswitch, lookupswitch
		return -1, EffectFixed

	case 0xac, 0xae, 0xb0: // ireturn, freturn, areturn
		return -1, EffectFixed
	case 0xad, 0xaf: // lreturn, dreturn
		return -2, EffectFixed
	case 0xb1: // return
		return 0, EffectFixed

	case 0xb2, 0xb3, 0xb4, 0xb5: // getstatic, putstatic, getfield, putfield
		return 0, EffectDescriptor
	case 0xb6, 0xb7, 0xb8, 0xb9: // invokevirtual/special/static/interface
		return 0, EffectDescriptor

	case 0xbb: // new
		return 1, EffectFixed
	case 0xbc, 0xbd, 0xbe, 0xc0, 0xc1: // newarray anewarray arraylength checkcast instanceof
		return 0, EffectFixed
	case 0xbf, 0xc2, 0xc3: // athrow, monitorenter, monitorexit
		return -1, EffectFixed

	case 0xc4, 0xc5: // wide, multianewarray
		return 0, EffectDescriptor

	case 0xa8, 0xa9, 0xba: // jsr, ret, invokedynamic — outside accepted v52 input
		return 0, EffectIllegal
	}
	return 0, EffectIllegal // reserved / unassigned opcode bytes
}
