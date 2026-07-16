package transpile

import "catty/rtda"

// monitorenter / monitorexit opcode values and invoke family opcodes.
const (
	opcMonitorenter  = 0xc2
	opcMonitorexit   = 0xc3
	opcInvokevirtual = 0xb6
	opcInvokespecial = 0xb7
	opcInvokestatic  = 0xb8
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
//  3. invokevirtual/invokespecial/invokestatic targeting Thread or
//     Object.wait/notify/notifyAll
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

		// Single-byte monitor opcodes.
		if b == opcMonitorenter {
			return "monitorenter bytecode, concurrency not supported in AOT"
		}
		if b == opcMonitorexit {
			return "monitorexit bytecode, concurrency not supported in AOT"
		}

		// Three-byte invoke instructions — check constant pool target.
		if (b == opcInvokevirtual || b == opcInvokespecial || b == opcInvokestatic) && cp != nil && i+2 < len(code) {
			index := uint16(code[i+1])<<8 | uint16(code[i+2])
			className, name, _ := cp.MemberRef(index)
			key := className + "." + name
			if concurrencyClasses[className] {
				return "invoke on " + key + " (Thread), concurrency not supported in AOT"
			}
			if concurrencyMethods[key] {
				return "invoke on " + key + ", concurrency not supported in AOT"
			}
			i += 3
			continue
		}

		i++
	}

	return ""
}
