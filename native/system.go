package native

import (
	"math"
	"time"

	"catty/rtda"
)

// init registers native methods that real JDK classes declare but catty
// implements in Go. These are looked up by (className, methodName, descriptor)
// when the classloader builds a real JDK class with ACC_NATIVE methods.
func init() {
	// --- System ---
	RegisterNative("java/lang/System", "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", systemArrayCopy)
	RegisterNative("java/lang/System", "currentTimeMillis", "()J", systemCurrentTimeMillis)
	RegisterNative("java/lang/System", "nanoTime", "()J", systemNanoTime)
	RegisterNative("java/lang/System", "identityHashCode", "(Ljava/lang/Object;)I", systemIdentityHashCode)
	RegisterNative("java/lang/System", "mapLibraryName", "(Ljava/lang/String;)Ljava/lang/String;", systemMapLibraryName)

	// --- Object ---
	RegisterNative("java/lang/Object", "hashCode", "()I", objectHashCode)
	RegisterNative("java/lang/Object", "getClass", "()Ljava/lang/Class;", objectGetClass)
	RegisterNative("java/lang/Object", "clone", "()Ljava/lang/Object;", objectClone)
	RegisterNative("java/lang/Object", "notify", "()V", objectNotify)
	RegisterNative("java/lang/Object", "notifyAll", "()V", objectNotifyAll)
	RegisterNative("java/lang/Object", "wait", "()V", objectWait0)
	RegisterNative("java/lang/Object", "wait", "(J)V", objectWait)
	RegisterNative("java/lang/Object", "wait", "(JI)V", objectWaitJI)
	RegisterNative("java/lang/Object", "registerNatives", "()V", nop)

	// --- Class ---
	RegisterNative("java/lang/Class", "getName", "()Ljava/lang/String;", classGetName)
	RegisterNative("java/lang/Class", "getSimpleName", "()Ljava/lang/String;", classGetSimpleName)
	RegisterNative("java/lang/Class", "desiredAssertionStatus", "()Z", classDesiredAssertionStatus)
	RegisterNative("java/lang/Class", "isInterface", "()Z", classIsInterface)
	RegisterNative("java/lang/Class", "isArray", "()Z", classIsArray)
	RegisterNative("java/lang/Class", "getModifiers", "()I", classGetModifiers)
	RegisterNative("java/lang/Class", "isInstance", "(Ljava/lang/Object;)Z", classIsInstance)
	RegisterNative("java/lang/Class", "isAssignableFrom", "(Ljava/lang/Class;)Z", classIsAssignableFrom)
	RegisterNative("java/lang/Class", "getSuperclass", "()Ljava/lang/Class;", classGetSuperclass)
	RegisterNative("java/lang/Class", "isHidden", "()Z", classIsHidden)
	RegisterNative("java/lang/Class", "getPrimitiveClass", "(Ljava/lang/String;)Ljava/lang/Class;", classGetPrimitiveClass)
	RegisterNative("java/lang/Class", "registerNatives", "()V", nop)

	// --- Thread (implementations in thread.go; registered via synthetic class) ---
	RegisterNative("java/lang/Thread", "holdsLock", "(Ljava/lang/Object;)Z", threadHoldsLock)
	RegisterNative("java/lang/Thread", "registerNatives", "()V", nop)

	// --- String ---
	RegisterNative("java/lang/String", "intern", "()Ljava/lang/String;", stringIntern)

	// --- Runtime ---
	RegisterNative("java/lang/Runtime", "availableProcessors", "()I", runtimeAvailableProcessors)
	RegisterNative("java/lang/Runtime", "freeMemory", "()J", runtimeZeroLong)
	RegisterNative("java/lang/Runtime", "totalMemory", "()J", runtimeZeroLong)
	RegisterNative("java/lang/Runtime", "maxMemory", "()J", runtimeZeroLong)
	RegisterNative("java/lang/Runtime", "gc", "()V", nop)

	// --- Float/Double bit conversion (native in JDK) ---
	RegisterNative("java/lang/Float", "floatToRawIntBits", "(F)I", floatToRawIntBits)
	RegisterNative("java/lang/Float", "intBitsToFloat", "(I)F", intBitsToFloat)
	RegisterNative("java/lang/Double", "doubleToRawLongBits", "(D)J", doubleToRawLongBits)
	RegisterNative("java/lang/Double", "longBitsToDouble", "(J)D", longBitsToDouble)

	// --- AccessController / SecurityManager (stubs) ---
	RegisterNative("java/security/AccessController", "doPrivileged", "(Ljava/security/PrivilegedAction;)Ljava/lang/Object;", nopRef0)
	RegisterNative("java/security/AccessController", "doPrivileged", "(Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;", nopRef0)
}

// --- System native methods ---

func systemArrayCopy(f *rtda.Frame) {
	src := f.GetRef(0)
	srcPos := int(f.GetInt(1))
	dst := f.GetRef(2)
	dstPos := int(f.GetInt(3))
	length := int(f.GetInt(4))

	if src == nil || dst == nil {
		return
	}
	if !src.Class().IsArray() || !dst.Class().IsArray() {
		return
	}
	// Use typed copy with overlap-safe memmove semantics and per-component-kind
	// dispatch so long/double preserve the full 64-bit value.
	rtda.CopyObjectCells(dst, src, dstPos, srcPos, length, src.Class().ComponentKind())
}

func systemCurrentTimeMillis(f *rtda.Frame) {
	f.PushLong(time.Now().UnixMilli())
}

func systemNanoTime(f *rtda.Frame) {
	f.PushLong(time.Now().UnixNano())
}

func systemIdentityHashCode(f *rtda.Frame) {
	obj := f.GetRef(0)
	f.PushInt(int32(uintptr(unsafePointer(obj)) & 0x7FFFFFFF))
}

// --- Object native methods ---

func objectHashCode(f *rtda.Frame) {
	obj := f.GetRef(0)
	f.PushInt(int32(uintptr(unsafePointer(obj)) & 0x7FFFFFFF))
}

func objectGetClass(f *rtda.Frame) {
	obj := f.GetRef(0)
	if obj == nil {
		f.PushRef(nil)
		return
	}
	clsObj := getClassObject(f.Thread(), obj.Class())
	f.PushRef(clsObj)
}

func objectClone(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		f.PushRef(nil)
		return
	}
	// Shallow clone: copy all fields.
	clone := rtda.CloneObject(this)
	f.PushRef(clone)
}

// --- Class native methods ---

func classGetName(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	name := javaToDot(cls.Name())
	f.PushRef(newStringFromGo(f.Thread(), name))
}

func classGetSimpleName(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	name := cls.Name()
	// Extract simple name after last '/'
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	f.PushRef(newStringFromGo(f.Thread(), name))
}

func classDesiredAssertionStatus(f *rtda.Frame) {
	f.PushInt(0) // assertions disabled
}

func classIsInterface(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	if cls.IsInterface() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func classIsArray(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	if cls.IsArray() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func classGetModifiers(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	f.PushInt(int32(cls.AccessFlags()))
}

func classIsInstance(f *rtda.Frame) {
	this := f.GetRef(0)
	other := f.GetRef(1)
	cls := getClassFromExtra(this)
	if other == nil {
		f.PushInt(0)
		return
	}
	if other.IsInstanceOf(cls) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func classIsAssignableFrom(f *rtda.Frame) {
	this := f.GetRef(0)
	argClassObj := f.GetRef(1)
	cls := getClassFromExtra(this)
	argCls := getClassFromExtra(argClassObj)
	if argCls == nil {
		f.PushInt(0)
		return
	}
	// cls.isAssignableFrom(argCls) — argCls is assignable to cls?
	// rtda.Class.isAssignableFrom(c) checks: can a value of `c` be assigned to `cls`?
	// We invert: argCls.IsInstanceOf(cls)?
	// Actually: cls.isAssignableFrom(argCls) = argCls is a subclass of cls (or same).
	// Check if cls is an ancestor (or same) of argCls.
	for c := argCls; c != nil; c = c.SuperClass() {
		if c == cls {
			f.PushInt(1)
			return
		}
	}
	f.PushInt(0)
}

func classGetSuperclass(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	super := cls.SuperClass()
	if super == nil {
		f.PushRef(nil)
		return
	}
	// Use canonical Class mirror (ADR-0033, K2 mirror continuity).
	superObj := super.ClassObject()
	f.PushRef(superObj)
}

func classIsHidden(f *rtda.Frame) {
	f.PushInt(0) // no hidden classes in catty
}

// --- Float/Double bit conversion ---

func floatToRawIntBits(f *rtda.Frame) {
	v := math.Float32frombits(uint32(f.GetInt(0)))
	f.PushInt(int32(math.Float32bits(v)))
}

func intBitsToFloat(f *rtda.Frame) {
	f.PushFloat(math.Float32frombits(uint32(f.GetInt(0))))
}

func doubleToRawLongBits(f *rtda.Frame) {
	bits := uint64(f.GetLong(0))
	f.PushLong(int64(bits))
}

func longBitsToDouble(f *rtda.Frame) {
	bits := uint64(f.GetLong(0))
	f.PushDouble(math.Float64frombits(bits))
}

// --- Missing native stubs ---

func systemMapLibraryName(f *rtda.Frame) {
	// Identity: just return the input name.
	this := f.GetRef(0)
	f.PushRef(this)
}

func stringIntern(f *rtda.Frame) {
	// No interning: just return this.
	this := f.GetRef(0)
	f.PushRef(this)
}

func runtimeAvailableProcessors(f *rtda.Frame) {
	f.PushInt(1) // minimal
}

func runtimeZeroLong(f *rtda.Frame) {
	f.PushLong(0)
}

func classGetPrimitiveClass(f *rtda.Frame) {
	name := f.GetRef(0)
	if name == nil {
		f.PushRef(nil)
		return
	}
	// Extract the primitive class name from the StringValue.
	sv := stringValueSV(name)
	primName := sv.GoString()

	// Look up the canonical VM primitive Class (ADR-0033, K2).
	cls := rtda.VMPrimitiveForName(primName)
	if cls == nil {
		f.PushRef(nil)
		return
	}
	// Return the canonical Class mirror for this VM type.
	mirror := cls.ClassObject()
	f.PushRef(mirror)
}

func nopBool0(f *rtda.Frame) { f.PushInt(0) }

func nopRef0(f *rtda.Frame) { f.PushRef(nil) }

// --- Object.wait / notify / notifyAll (Slice C, ADR-0029) ---

// objectWait0 implements Object.wait() — the no-argument indefinite wait.
func objectWait0(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Object")
		return
	}
	m := this.Monitor()
	thread := f.Thread()
	ec := thread.EC()

	if !m.HoldsLock(ec) {
		throwIMSE(f)
		return
	}

	savedDepth := m.RecursionDepth()
	_, interrupted := thread.MonitorWait(m, savedDepth)
	if interrupted {
		throwInterruptedException(f)
	}
}

// objectWait implements Object.wait(long timeoutMillis).
// It releases the monitor, blocks until notified or interrupted, then reacquires
// the monitor with the original recursion depth restored. If the thread was
// interrupted before or during the wait, InterruptedException is thrown.
// The timeout is ignored for now (indefinite wait); timed wait is Slice D.
func objectWait(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Object")
		return
	}
	// java.lang.Object.wait(long timeoutMillis) specifies:
	//   IllegalArgumentException - if timeoutMillis is negative
	timeoutMillis := f.GetLong(1)
	if timeoutMillis < 0 {
		throwIllegalArgumentException(f, "timeout value is negative")
		return
	}
	m := this.Monitor()
	thread := f.Thread()
	ec := thread.EC()

	if !m.HoldsLock(ec) {
		throwIMSE(f)
		return
	}

	savedDepth := m.RecursionDepth()
	_, interrupted := thread.MonitorWait(m, savedDepth)
	if interrupted {
		throwInterruptedException(f)
	}
	// Normal return: notify won, thread reacquired monitor with depth restored.
	// No need to clear interrupt — if interrupted, the flag was already handled.
}

// objectWaitJI implements Object.wait(long timeoutMillis, int nanos).
// The nanos are ignored for now (indefinite wait); timed wait is Slice D,
// but the argument range is validated per java.lang.Object.wait(long,int):
//
//	IllegalArgumentException - if timeoutMillis is negative,
//	                           or if the value of nanos is out of range
func objectWaitJI(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Object")
		return
	}
	timeoutMillis := f.GetLong(1)
	nanos := f.GetInt(3)
	if timeoutMillis < 0 {
		throwIllegalArgumentException(f, "timeout value is negative")
		return
	}
	if nanos < 0 || nanos > 999999 {
		throwIllegalArgumentException(f, "nanosecond timeout value out of range")
		return
	}
	m := this.Monitor()
	thread := f.Thread()
	ec := thread.EC()

	if !m.HoldsLock(ec) {
		throwIMSE(f)
		return
	}

	savedDepth := m.RecursionDepth()
	_, interrupted := thread.MonitorWait(m, savedDepth)
	if interrupted {
		throwInterruptedException(f)
	}
}

// throwIllegalArgumentException throws an IllegalArgumentException with the
// given detail message on the calling thread.
func throwIllegalArgumentException(f *rtda.Frame, message string) {
	throwException(f, "java/lang/IllegalArgumentException", message)
}

// objectNotify implements Object.notify(). It wakes a single thread waiting on
// this object's monitor. The caller must own the monitor, or
// IllegalMonitorStateException is thrown.
func objectNotify(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Object")
		return
	}
	m := this.Monitor()
	if !m.Notify(f.Thread().EC()) {
		throwIMSE(f)
	}
}

// objectNotifyAll implements Object.notifyAll(). It wakes all threads waiting on
// this object's monitor. The caller must own the monitor, or
// IllegalMonitorStateException is thrown.
func objectNotifyAll(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Object")
		return
	}
	m := this.Monitor()
	if !m.NotifyAll(f.Thread().EC()) {
		throwIMSE(f)
	}
}

// threadHoldsLock implements the static Thread.holdsLock(Object).
// Returns true if the calling thread owns the monitor on the given object.
func threadHoldsLock(f *rtda.Frame) {
	obj := f.GetRef(0) // static method: local 0 = first arg
	if obj == nil {
		// Per java.lang.Thread: passing null to a method in this class
		// throws NullPointerException unless otherwise specified.
		throwNPE(f, "null")
		return
	}
	if obj.Monitor().HoldsLock(f.Thread().EC()) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

// throwIMSE throws an IllegalMonitorStateException on the calling thread.
func throwIMSE(f *rtda.Frame) {
	throwException(f, "java/lang/IllegalMonitorStateException", "")
}
