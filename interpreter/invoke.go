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
func invokeMethod(thread *rtda.Thread, method *rtda.Method) {
	if method == nil {
		panic(fmt.Sprintf("catty: invokeMethod received nil method from frame %s.%s%s",
			thread.CurrentFrame().Method().Owner().Name(),
			thread.CurrentFrame().Method().Name(),
			thread.CurrentFrame().Method().Descriptor()))
	}
	if method.IsNative() {
		invokeNative(thread, method)
		return
	}
	if method.IsNative() {
		invokeNative(thread, method)
		return
	}
	caller := thread.CurrentFrame()
	frame := thread.NewFrame(method)
	copyArgs(caller, frame, method)
	thread.PushFrame(frame)
}

func invokeNative(thread *rtda.Thread, method *rtda.Method) {
	caller := thread.CurrentFrame()
	frame := thread.NewFrame(method)
	copyArgs(caller, frame, method)
	method.NativeFunc()(frame)
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

// ensureInitialized runs a class's <clinit> the first time the class is used at
// a JVMS §5.5 initialization point (new/getstatic/putstatic/invokestatic).
// The <clinit> runs synchronously on a throwaway frame so the caller's frame
// and operand stack are undisturbed when ensureInitialized returns.
func ensureInitialized(thread *rtda.Thread, class *rtda.Class) {
	if class.InitStarted() || class.IsInterface() {
		return
	}
	class.MarkInitStarted()
	if class.SuperClass() != nil {
		ensureInitialized(thread, class.SuperClass())
	}
	if clinit := class.GetMethod("<clinit>", "()V"); clinit != nil {
		runClinit(thread, clinit)
	}
}

// runClinit runs a <clinit> method synchronously: push a frame, run only the
// clinit frame (not the whole thread), and return when the clinit frame is
// popped. This prevents Loop from continuing into the caller's frame while
// the caller's opcode handler (e.g. 'new') is still mid-execution.
func runClinit(thread *rtda.Thread, method *rtda.Method) {
	clinitDepth := thread.FrameCount() + 1 // target depth after push
	thread.PushFrame(thread.NewFrame(method))
	for thread.FrameCount() >= clinitDepth && !thread.IsStackEmpty() {
		frame := thread.CurrentFrame()
		opcodePc := frame.PC()
		op := opcode.Opcode(frame.Code()[opcodePc])
		frame.SetPC(opcodePc + 1)
		exec(thread, frame, op, opcodePc)
		if thread.HasException() {
			handleException(thread, opcodePc)
		}
	}
}

// InitClass is the exported form of ensureInitialized, for the launcher to
// initialize the main class before entering main().
func InitClass(thread *rtda.Thread, class *rtda.Class) {
	ensureInitialized(thread, class)
}

// pushConstant handles ldc / ldc_w: pushes an int, float, String, or Class
// constant onto the operand stack per the constant pool tag at index.
func pushConstant(thread *rtda.Thread, frame *rtda.Frame, cp *classfile.ConstantPool, index uint16) {
	switch cp.Tag(index) {
	case classfile.ConstantInteger:
		frame.PushInt(cp.Integer(index))
	case classfile.ConstantFloat:
		frame.PushFloat(cp.Float(index))
	case classfile.ConstantString:
		frame.PushRef(newString(thread, cp.String(index)))
	case classfile.ConstantClass:
		className := cp.ClassName(index)
		cls := thread.Loader().LoadClass(className)
		frame.PushRef(getClassObject(thread, cls))
	}
}

// newString creates a java.lang.String object backed by a Go string value.
func newString(thread *rtda.Thread, value string) *rtda.Object {
	class := thread.Loader().LoadClass("java/lang/String")
	obj := rtda.NewObject(class)
	obj.SetExtra(value)
	return obj
}

// getClassObject creates a java.lang.Class object wrapping cls.
// The Class object stores the rtda.Class in its extra field.
func getClassObject(thread *rtda.Thread, cls *rtda.Class) *rtda.Object {
	classClass := thread.Loader().LoadClass("java/lang/Class")
	obj := rtda.NewObject(classClass)
	obj.SetExtra(cls)
	return obj
}
