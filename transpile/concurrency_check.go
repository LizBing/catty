package transpile

import "catty/rtda"

// monitorenter / monitorexit opcode values and invoke family opcodes.
const (
	opcMonitorenter    = 0xc2
	opcMonitorexit     = 0xc3
	opcInvokevirtual   = 0xb6
	opcInvokespecial   = 0xb7
	opcInvokestatic    = 0xb8
	opcInvokeinterface = 0xb9
)

// concurrencyMethods is the set of classes whose methods (any) trigger the
// AOT concurrency boundary, plus specific Object methods. Any invoke targeting
// one of these class+method keys is a build-time rejection.
var concurrencyClasses = map[string]bool{
	"java/lang/Thread": true,
}

var concurrencyMethods = map[string]bool{
	"java/lang/Object.wait":      true,
	"java/lang/Object.notify":    true,
	"java/lang/Object.notifyAll": true,
}

// concurrencyReason returns a non-empty reason if the method triggers the
// AOT concurrency boundary. The check covers:
//
//  1. monitor bytecodes (monitorenter/monitorexit) in the method's code
//  2. ACC_SYNCHRONIZED flag on the method
//  3. invokevirtual/invokespecial/invokestatic/invokeinterface targeting
//     Thread or Object.wait/notify/notifyAll
//
// This is a conservative build-time rejection; AOT concurrency execution
// is Not implemented (ADR-0028, ADR-0029) and any program that touches
// monitors, synchronized methods, or Thread/Object concurrency methods
// must not produce a native binary.
func concurrencyReason(m *rtda.Method) string {
	if m.IsNative() || len(m.Code()) == 0 {
		return ""
	}

	// ACC_SYNCHRONIZED methods — implicit monitor entry/exit.
	if m.IsSynchronized() {
		return "synchronized method (ACC_SYNCHRONIZED), concurrency not supported in AOT"
	}

	cp := m.Owner().ConstantPool()
	code := m.Code()

	for i := 0; i < len(code); {
		b := code[i]
		length := instLength(code, i)

		// Single-byte monitor opcodes.
		if b == opcMonitorenter {
			return "monitorenter bytecode, concurrency not supported in AOT"
		}
		if b == opcMonitorexit {
			return "monitorexit bytecode, concurrency not supported in AOT"
		}

		// Invoke family — check constant-pool target for Thread/Object
		// concurrency methods. Includes invokevirtual (0xb6), invokespecial
		// (0xb7), invokestatic (0xb8), and invokeinterface (0xb9).
		// invokedynamic (0xba) is correctly stepped by instLength but not
		// target-checked — its bootstrap-method indirection cannot be
		// resolved from the bytecode alone and it cannot reach AOT.
		if (b == opcInvokevirtual || b == opcInvokespecial || b == opcInvokestatic || b == opcInvokeinterface) &&
			cp != nil && i+2 < len(code) {
			index := uint16(code[i+1])<<8 | uint16(code[i+2])
			className, name, _ := cp.MemberRef(index)
			key := className + "." + name
			if concurrencyClasses[className] {
				return "invoke on " + key + " (Thread), concurrency not supported in AOT"
			}
			if concurrencyMethods[key] {
				return "invoke on " + key + ", concurrency not supported in AOT"
			}
		}

		i += length
	}

	return ""
}

// --- instruction-length helpers ---

// instLength returns the byte length of the JVM instruction at code[pc].
// It covers all opcodes including variable-length instructions
// (tableswitch, lookupswitch, wide). Unknown/invalid opcodes default to 1.
//
// This is a self-contained duplicate of the length logic in
// lowering/lower.go:decodeInst — the lowering versions are unexported and
// coupled to full instruction decoding, so we keep a minimal copy here
// for the byte scanner.
func instLength(code []byte, pc int) int {
	if pc >= len(code) {
		return 1
	}
	switch code[pc] {
	// ---- 2-byte instructions ----
	case 0x10: // bipush
		return 2
	case 0x12: // ldc
		return 2
	case 0x15, 0x16, 0x17, 0x18, 0x19: // iload, lload, fload, dload, aload
		return 2
	case 0x36, 0x37, 0x38, 0x39, 0x3a: // istore, lstore, fstore, dstore, astore
		return 2
	case 0xbc: // newarray
		return 2

	// ---- 3-byte instructions ----
	case 0x11: // sipush
		return 3
	case 0x13, 0x14: // ldc_w, ldc2_w
		return 3
	case 0x84: // iinc
		return 3
	case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e: // ifeq..ifle
		return 3
	case 0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4: // if_icmpXX
		return 3
	case 0xa5, 0xa6: // if_acmpeq, if_acmpne
		return 3
	case 0xa7: // goto
		return 3
	case 0xb2, 0xb3, 0xb4, 0xb5: // getstatic, putstatic, getfield, putfield
		return 3
	case 0xb6, 0xb7, 0xb8: // invokevirtual, invokespecial, invokestatic
		return 3
	case 0xbb: // new
		return 3
	case 0xbd: // anewarray
		return 3
	case 0xc0, 0xc1: // checkcast, instanceof
		return 3
	case 0xc6, 0xc7: // ifnull, ifnonnull
		return 3

	// ---- 4-byte instructions ----
	case 0xc5: // multianewarray
		return 4

	// ---- 5-byte instructions ----
	case 0xb9: // invokeinterface
		return 5
	case 0xba: // invokedynamic
		return 5
	case 0xc8: // goto_w
		return 5

	// ---- variable-length ----
	case 0xaa: // tableswitch
		return switchLen(code, pc, true)
	case 0xab: // lookupswitch
		return switchLen(code, pc, false)
	case 0xc4: // wide
		// wide <opcode> — if the modified opcode is iinc, length is 6;
		// otherwise 4 (wide load/store).
		if pc+1 < len(code) && code[pc+1] == 0x84 {
			return 6
		}
		return 4

	// ---- 1-byte (default) ----
	default:
		return 1
	}
}

// switchLen returns the full byte length of a tableswitch or lookupswitch
// instruction, including 0–3 bytes of alignment padding.
func switchLen(code []byte, pc int, table bool) int {
	base := pc + 1
	p := base + padTo4(base)
	if table {
		low := int32(be32(code, p+4))
		high := int32(be32(code, p+8))
		n := high - low + 1
		return (p - pc) + 12 + int(n)*4
	}
	npairs := int32(be32(code, p+4))
	return (p - pc) + 8 + int(npairs)*8
}

// padTo4 returns the number of bytes needed to align n to a 4-byte boundary.
func padTo4(n int) int { return (4 - n%4) % 4 }

// be32 reads a big-endian uint32 from code at offset off.
func be32(code []byte, off int) uint32 {
	return uint32(code[off])<<24 | uint32(code[off+1])<<16 |
		uint32(code[off+2])<<8 | uint32(code[off+3])
}
