// Package runtime is the AOT bridge: the entry points that AOT-transpiled Go
// code calls to reach catty's runtime (the "world transition" of ADR-0007).
//
// Transpiled methods can't, on their own, resolve classes/fields/methods or run
// native/interpreted code — they call into this package, which holds the
// classloader + thread and resolves targets by (class, name, descriptor) at run
// time. A2.2 supports native targets (e.g. System.out.println); interpreted
// targets need a catcher frame and come later.
package runtime

import (
	"math"

	"catty/classloader"
	"catty/classpath"
	"catty/interpreter"
	"catty/rtda"
)

// loader and thread are set by Bootstrap and shared across bridge calls.
var (
	loader rtda.Loader
	thread *rtda.Thread
)

// Bootstrap loads the main class (and its dependencies, including the native
// core classes) and runs its <clinit>, so the bridge can resolve targets. Called
// by the emitted program's main before the transpiled main method runs.
func Bootstrap(classpathStr, mainClass string) {
	cl := classloader.New(classpath.Parse(classpathStr))
	loader = cl
	thread = rtda.NewThread(cl)
	interpreter.InitClass(thread, cl.LoadClass(mainClass))
}

// GetStatic reads a static field, resolving the declaring class at run time.
func GetStatic(class, name, desc string) rtda.Slot {
	c := loader.LoadClass(class)
	interpreter.InitClass(thread, c)
	field := c.LookupField(name, desc)
	return c.StaticVars()[field.SlotID()]
}

// InvokeVirtual dispatches a virtual call: args[0] is `this`, and the target is
// resolved on the receiver's runtime class (dynamic dispatch). Runs native or
// interpreted targets.
func InvokeVirtual(class, name, desc string, args []rtda.Slot) rtda.Slot {
	_ = loader.LoadClass(class) // ensures the class (and its methods) are loaded
	recv := args[0].Ref()
	if recv == nil {
		panic("catty: NullPointerException")
	}
	method := recv.Class().LookupMethod(name, desc)
	return runMethod(method, args)
}

// InvokeSpecial resolves the target on the declared class (used for <init> and
// super/private calls). Runs native or interpreted targets.
func InvokeSpecial(class, name, desc string, args []rtda.Slot) rtda.Slot {
	method := loader.LoadClass(class).LookupMethod(name, desc)
	return runMethod(method, args)
}

// NewObject allocates an instance of class (without running <init> — the caller
// follows with InvokeSpecial("<init>") to initialize it, matching `new`).
func NewObject(class string) *rtda.Object {
	return rtda.NewObject(loader.LoadClass(class))
}

// runMethod runs a native target (synchronously) or an interpreted target (via
// the bridge's RunMethod).
func runMethod(method *rtda.Method, args []rtda.Slot) rtda.Slot {
	if method.IsNative() {
		return runNative(method, args)
	}
	return interpreter.RunMethod(thread, method, args)
}

// NewString creates a java.lang.String carrying the Go string in its extra
// payload, matching how the interpreter represents ldc strings.
func NewString(s string) *rtda.Object {
	obj := rtda.NewObject(loader.LoadClass("java/lang/String"))
	obj.SetExtra(s)
	return obj
}

// NewIntArray builds a Java int[] from Go values — a convenience for transpiled
// test programs that need an array input without constructing one via bytecode.
func NewIntArray(values ...int32) *rtda.Object {
	arr := rtda.NewArray(loader.LoadClass("[I"), len(values))
	for i, v := range values {
		arr.ArrayElementSlot(i).SetNum(v)
	}
	return arr
}

// FloatMod / DoubleMod implement Java's floating-point % (fmod semantics:
// result sign = dividend). Go has no `%` operator on floats, so the emitted
// code calls these instead.
func FloatMod(a, b float32) float32  { return float32(math.Mod(float64(a), float64(b))) }
func DoubleMod(a, b float64) float64 { return math.Mod(a, b) }

// runNative sets up a frame with the given argument slots, runs the native
// method, and returns its result slot (zero for void).
func runNative(method *rtda.Method, args []rtda.Slot) rtda.Slot {
	if !method.HasNativeImplementation() {
		thread.ThrowUnsatisfiedLinkError(
			method.Owner().Name(),
			method.Name(),
			method.ReturnType(),
		)
		return rtda.Slot{}
	}
	frame := thread.NewFrame(method)
	for i, a := range args {
		frame.SetSlot(i, a)
	}
	method.NativeFunc()(frame)
	return popReturn(frame, method.ReturnType())
}

// popReturn extracts the native method's return value from its throwaway frame.
func popReturn(frame *rtda.Frame, ret string) rtda.Slot {
	switch ret {
	case "V", "":
		return rtda.Slot{}
	case "J", "D":
		// long/double span two slots; representing them through a single Slot
		// return needs typed bridge variants — deferred (not needed for A2.2).
		panic("catty/runtime: long/double invoke return not supported in A2.2")
	}
	return frame.PopSlot()
}

// InvokeStatic resolves and runs a static method (native or interpreted) by
// (class, name, desc). Used by the AOT bridge when an invokestatic target isn't
// AOT'd in the emitted binary.
func InvokeStatic(class, name, desc string, args []rtda.Slot) rtda.Slot {
	return runMethod(loader.LoadClass(class).LookupMethod(name, desc), args)
}

// Thread returns the runtime's thread (for the fallback interpreter path).
func Thread() *rtda.Thread { return thread }

// Loader returns the runtime's class loader.
func Loader() rtda.Loader { return loader }
