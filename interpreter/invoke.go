package interpreter

import (
	"fmt"

	"catty/classfile"
	"catty/opcode"
	"catty/rtda"
)

// invokeMethod sets up a method call: it copies arguments from the caller's
// operand stack into a fresh callee frame. Interpreted methods get the frame
// pushed for the dispatch loop to run; native methods run synchronously and
// their return value is handed back to the caller immediately.
func invokeMethod(context *rtda.ExecutionContext, method *rtda.Method) {
	if method == nil {
		panic(fmt.Sprintf("catty: invokeMethod received nil method from frame %s.%s%s",
			context.CurrentFrame().Method().Owner().Name(),
			context.CurrentFrame().Method().Name(),
			context.CurrentFrame().Method().Descriptor()))
	}
	if method.IsNative() {
		invokeNative(context, method)
		return
	}
	caller := context.CurrentFrame()
	frame := context.NewFrame(method)
	copyArgs(caller, frame, method)
	frame.EnterSyncMonitor()
	context.PushFrame(frame)
}

func invokeNative(context *rtda.ExecutionContext, method *rtda.Method) {
	caller := context.CurrentFrame()
	frame := context.NewFrame(method)
	copyArgs(caller, frame, method)

	// Native throwaway frames are never pushed to the stack, so PopFrame won't
	// release the implicit monitor. Enter before the native call and exit after,
	// even on exception (ADR-0029).
	frame.EnterSyncMonitor()

	method.NativeFunc()(frame)

	if so := frame.SyncObject(); so != nil {
		so.Monitor().Exit(context.EC())
	}

	// A native method that throws an exception (via Thread.Throw) must
	// NOT transfer a return value — the callee frame's stack may be in
	// an inconsistent state, and the interpreter will dispatch the
	// pending exception instead.
	if context.HasException() {
		return
	}
	transferReturn(caller, frame, method.ReturnType())
}

// copyArgs pops totalSlots (= parameters + `this` for instance methods) slots
// off the caller's stack into the callee's locals, in reverse so that the
// top-of-stack slot (last argument) lands in the highest local index.
func copyArgs(caller, callee *rtda.Frame, method *rtda.Method) {
	total := int(method.ArgSlotCount())
	if !method.IsStatic() {
		total++ // `this`
	}
	if total > callee.LocalsLen() {
		panic(fmt.Sprintf("catty: copyArgs locals overflow in %s.%s%s: need %d slots, locals len=%d maxLocals=%d",
			method.Owner().Name(), method.Name(), method.Descriptor(), total, callee.LocalsLen(), method.MaxLocals()))
	}
	for i := total - 1; i >= 0; i-- {
		callee.SetSlot(i, caller.PopSlot())
	}
}

// transferReturn moves a native method's return value from its throwaway frame
// onto the caller's stack, by return-type descriptor.
func transferReturn(caller, callee *rtda.Frame, ret string) {
	switch ret {
	case "V":
	case "J", "D": // category-2: two slots, popped high-then-low (pushed low-then-high)
		caller.PushSlot(callee.PopSlot())
		caller.PushSlot(callee.PopSlot())
	default:
		caller.PushSlot(callee.PopSlot())
	}
}

// ensureInitialized runs class/interface initialization at a JVMS §5.5 point.
// This is the interpreter's ClinitRunner callback passed to the shared
// initialization service. The <clinit> runs synchronously via the existing
// runClinit loop so the caller's frame and operand stack are undisturbed.
func ensureInitialized(context *rtda.ExecutionContext, class *rtda.Class) {
	if class.IsInitialized() {
		return
	}
	loader := context.Loader()
	result := rtda.InitializeClass(loader, class, context.EC(), func(c *rtda.Class, m *rtda.Method) rtda.InitResult {
		return runClinit(context, m)
	})
	if result.ErrObj != nil {
		pc := 0
		if !context.IsStackEmpty() {
			pc = context.CurrentFrame().PC()
		}
		context.Throw(rtda.WrapInitFailure(loader, result.ErrObj), pc)
	}
}

// runClinit runs a <clinit> method synchronously: push a frame, run only the
// clinit frame (not the whole context), and return when the clinit frame is
// popped. This prevents Loop from continuing into the caller's frame while
// the caller's opcode handler (e.g. 'new') is still mid-execution.
//
// If <clinit> completes abruptly (uncaught exception anywhere in the clinit
// call chain), runClinit pops every frame back to the caller and returns the
// thrown object as InitResult.ErrObj.
func runClinit(context *rtda.ExecutionContext, method *rtda.Method) rtda.InitResult {
	clinitDepth := context.FrameCount() + 1 // target depth after push
	frame := context.NewFrame(method)
	frame.EnterSyncMonitor()
	context.PushFrame(frame)
	for context.FrameCount() >= clinitDepth && !context.IsStackEmpty() {
		frame := context.CurrentFrame()
		opcodePc := frame.PC()
		op := opcode.Opcode(frame.Code()[opcodePc])
		frame.SetPC(opcodePc + 1)
		exec(context, frame, op, opcodePc)
		if context.HasException() {
			thrown := context.ClearException()
			// Walk frames from the throwing frame down to (and including)
			// the clinit frame, searching for a handler.
			for context.FrameCount() >= clinitDepth {
				f := context.CurrentFrame()
				caught := false
				framePopped := false
				for _, entry := range f.Method().ExceptionTable() {
					if opcodePc >= entry.StartPc() && opcodePc < entry.EndPc() {
						// catchType "" = catch-all (finally).
						if entry.CatchType() == "" {
							f.ClearStack()
							f.PushRef(thrown)
							f.SetPC(entry.HandlerPc())
							caught = true
							break
						}
						// Typed catch-type resolution (K2).
						// On failure, the replacement throwable is a
						// NEW exception and must propagate from the
						// caller boundary — it must not be caught by
						// later handlers in the same frame.
						catchClass := resolveClass(context, opcodePc, entry.CatchType())
						if catchClass == nil {
							thrown = context.ClearException()
							context.PopFrame()
							framePopped = true
							if context.FrameCount() < clinitDepth {
								return rtda.InitResult{ErrObj: thrown}
							}
							opcodePc = context.CurrentFrame().PC() - 1
							break
						}
						if thrown.IsInstanceOf(catchClass) {
							f.ClearStack()
							f.PushRef(thrown)
							f.SetPC(entry.HandlerPc())
							caught = true
							break
						}
					}
				}
				if caught {
					break // handler found, resume execution
				}
				if framePopped {
					continue // already popped by catch-type failure
				}
				// Not caught in this frame — pop it.
				context.PopFrame()
				if context.FrameCount() < clinitDepth {
					// Exception propagated past clinit boundary.
					return rtda.InitResult{ErrObj: thrown}
				}
				// Set throwPC for the caller's exception-table search.
				opcodePc = context.CurrentFrame().PC() - 1
			}
		}
	}
	return rtda.SuccessInit()
}

// InitClass is the exported form of ensureInitialized, for the launcher to
// initialize the main class before entering main().
func InitClass(context *rtda.ExecutionContext, class *rtda.Class) {
	ensureInitialized(context, class)
}

// pushConstant handles ldc / ldc_w: pushes an int, float, String, or Class
// constant onto the operand stack per the constant pool tag at index.
// pc is the bytecode/IR offset for exception backtraces.
// Returns false if class resolution failed (exception already set on context);
// the caller must return immediately.
func pushConstant(context *rtda.ExecutionContext, frame *rtda.Frame, cp *classfile.ConstantPool, index uint16, pc int) bool {
	switch cp.Tag(index) {
	case classfile.ConstantInteger:
		frame.PushInt(cp.Integer(index))
	case classfile.ConstantFloat:
		frame.PushFloat(cp.Float(index))
	case classfile.ConstantString:
		frame.PushRef(newString(context, cp.UTF16(index)))
	case classfile.ConstantClass:
		className := cp.ClassName(index)
		cls := resolveClass(context, pc, className)
		if cls == nil {
			return false
		}
		frame.PushRef(getClassObject(context, cls))
	}
	return true
}

// newString creates a java.lang.String object from lossless UTF-16 code units
// (obtained from the classfile constant pool via decodeMUTF8ToUTF16).
func newString(context *rtda.ExecutionContext, units []uint16) *rtda.Object {
	class := context.Loader().LoadClass("java/lang/String")
	obj := rtda.NewObject(class)
	obj.SetExtra(rtda.NewStringValue(units))
	return obj
}

// getClassObject returns the canonical java.lang.Class object wrapping cls
// (ADR-0029). The Class object stores the rtda.Class in its extra field.
// All callers see the same Object identity for the same Class.
// K2: mirror creation uses the class's defining loader, not the context's
// initiating loader.
func getClassObject(context *rtda.ExecutionContext, cls *rtda.Class) *rtda.Object {
	return cls.ClassObject()
}
