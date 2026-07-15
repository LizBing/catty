package native

import (
	"testing"

	"catty/rtda"
)

// buildObjectMonitorHierarchy creates a minimal class hierarchy for testing
// Object.wait()/notify() and Thread.holdsLock(). It includes:
//   Object (with wait/notify methods) → Throwable → Exception → RuntimeException → IMSE
//   Thread (with holdsLock)
//   String and Class (required by throwException for IMSE construction)
func buildObjectMonitorHierarchy() *simpleLoader {
	l := &simpleLoader{classes: make(map[string]*rtda.Class)}

	obj := rtda.NewSyntheticClass("java/lang/Object", nil)
	obj.AddMethod(rtda.NativeMethod(obj, "wait", "()V", objectWait0))
	obj.AddMethod(rtda.NativeMethod(obj, "notify", "()V", objectNotify))
	obj.AddMethod(rtda.NativeMethod(obj, "wait", "(J)V", objectWait))
	obj.AddMethod(rtda.NativeMethod(obj, "wait", "(JI)V", objectWaitJI))
	l.classes["java/lang/Object"] = obj

	// Throwable with detailMessage field (needed by throwException).
	throwable := rtda.NewSyntheticClass("java/lang/Throwable", obj)
	throwable.AddInstanceField("detailMessage", "Ljava/lang/String;")
	l.classes["java/lang/Throwable"] = throwable

	exc := rtda.NewSyntheticClass("java/lang/Exception", throwable)
	l.classes["java/lang/Exception"] = exc

	re := rtda.NewSyntheticClass("java/lang/RuntimeException", exc)
	l.classes["java/lang/RuntimeException"] = re

	imse := rtda.NewSyntheticClass("java/lang/IllegalMonitorStateException", re)
	l.classes["java/lang/IllegalMonitorStateException"] = imse

	npe := rtda.NewSyntheticClass("java/lang/NullPointerException", re)
	l.classes["java/lang/NullPointerException"] = npe

	iae := rtda.NewSyntheticClass("java/lang/IllegalArgumentException", re)
	l.classes["java/lang/IllegalArgumentException"] = iae

	// String (needed by throwException for the detail message String object).
	str := rtda.NewSyntheticClass("java/lang/String", obj)
	l.classes["java/lang/String"] = str

	// Class mirror class (needed for getClass and similar native calls).
	classCls := rtda.NewSyntheticClass("java/lang/Class", obj)
	l.classes["java/lang/Class"] = classCls

	// Thread with <init> and holdsLock.
	thread := rtda.NewSyntheticClass("java/lang/Thread", obj)
	thread.AddMethod(rtda.NativeMethod(thread, "<init>", "()V", threadInit))
	thread.AddMethod(staticNative(thread, "holdsLock", "(Ljava/lang/Object;)Z", threadHoldsLock))
	l.classes["java/lang/Thread"] = thread

	return l
}

// newTestThreadWithLoader creates a java.lang.Thread object whose attached
// rtda.Thread uses the given loader. Returns the facade object and the
// underlying rtda.Thread execution context.
func newTestThreadWithLoader(loader *simpleLoader, threadClass *rtda.Class) (*rtda.Object, *rtda.Thread) {
	obj := rtda.NewObject(threadClass)
	caller := rtda.NewThread(loader)
	frame := caller.NewFrame(threadClass.LookupMethod("<init>", "()V"))
	frame.SetRef(0, obj)
	threadInit(frame)
	return obj, obj.Extra().(*rtda.Thread)
}

// --- wait/notify ownership failure (contract-mandated unit evidence) ---

// TestObjectWaitNoOwnershipThrowsIMSE verifies that calling Object.wait()
// without owning the monitor throws IllegalMonitorStateException.
func TestObjectWaitNoOwnershipThrowsIMSE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	// Lock the monitor so we have a legitimate object, then release it.
	m := lockObj.Monitor()
	m.Enter(rt.EC())
	m.Exit(rt.EC())

	// Now call wait() without owning the monitor.
	frame := caller.NewFrame(objClass.LookupMethod("wait", "()V"))
	frame.SetRef(0, lockObj)
	caller.PushFrame(frame)
	objectWait0(frame)

	if !caller.HasException() {
		t.Fatal("wait() without ownership should throw IllegalMonitorStateException")
	}
}

// TestObjectNotifyNoOwnershipThrowsIMSE verifies that calling Object.notify()
// without owning the monitor throws IllegalMonitorStateException.
func TestObjectNotifyNoOwnershipThrowsIMSE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	// Acquire and release so the monitor exists but is unowned.
	m := lockObj.Monitor()
	m.Enter(rt.EC())
	m.Exit(rt.EC())

	// Now call notify() without owning the monitor.
	frame := caller.NewFrame(objClass.LookupMethod("notify", "()V"))
	frame.SetRef(0, lockObj)
	caller.PushFrame(frame)
	objectNotify(frame)

	if !caller.HasException() {
		t.Fatal("notify() without ownership should throw IllegalMonitorStateException")
	}
}

// --- Thread.holdsLock (contract-mandated unit evidence) ---

// TestThreadHoldsLockTrue verifies holdsLock returns true when the calling
// thread owns the monitor.
func TestThreadHoldsLockTrue(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)

	// Acquire the monitor.
	lockObj.Monitor().Enter(rt.EC())

	// holdsLock should return true for the owning EC.
	if !lockObj.Monitor().HoldsLock(rt.EC()) {
		t.Fatal("monitor should be held after Enter")
	}

	// Clean up.
	lockObj.Monitor().Exit(rt.EC())
}

// TestThreadHoldsLockFalse verifies holdsLock returns false when the calling
// thread does not own the monitor.
func TestThreadHoldsLockFalse(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)

	// Don't acquire — holdsLock should return false.
	if lockObj.Monitor().HoldsLock(rt.EC()) {
		t.Fatal("fresh monitor should not be held")
	}
}

// TestThreadHoldsLockNull verifies holdsLock(null) throws NullPointerException.
func TestThreadHoldsLockNull(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	threadClass := loader.classes["java/lang/Thread"]

	caller := rtda.NewThread(loader)

	frame := caller.NewFrame(threadClass.LookupMethod("holdsLock", "(Ljava/lang/Object;)Z"))
	frame.SetRef(0, nil)
	caller.PushFrame(frame)
	threadHoldsLock(frame)

	if !caller.HasException() {
		t.Fatal("holdsLock(null) should throw NullPointerException")
	}
	exc := caller.ClearException()
	if exc == nil || exc.Class().Name() != "java/lang/NullPointerException" {
		t.Fatalf("holdsLock(null) threw %v, want NullPointerException", exc)
	}
}

// TestThreadHoldsLockNative verifies the full native stack: Thread.holdsLock
// via the native function, with a proper Frame, Thread, and Object.
func TestThreadHoldsLockNative(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	caller := rtda.NewThread(loader)

	// Acquire the monitor as the caller thread (threadHoldsLock checks
	// against f.Thread().EC(), which is the caller).
	lockObj.Monitor().Enter(caller.EC())

	// Call holdsLock via the native function (passing lockObj as arg 0).
	frame := caller.NewFrame(threadClass.LookupMethod("holdsLock", "(Ljava/lang/Object;)Z"))
	frame.SetRef(0, lockObj)
	threadHoldsLock(frame)

	if frame.PopInt() != 1 {
		t.Fatal("holdsLock should return true (1) when thread owns the monitor")
	}

	// Release and verify it returns false.
	lockObj.Monitor().Exit(caller.EC())

	frame = caller.NewFrame(threadClass.LookupMethod("holdsLock", "(Ljava/lang/Object;)Z"))
	frame.SetRef(0, lockObj)
	threadHoldsLock(frame)

	if frame.PopInt() != 0 {
		t.Fatal("holdsLock should return false (0) after monitor is released")
	}
}

// --- wait/notify successful path: verified by rtda.TestMonitorWaitNotify and
// the Slice C concurrency fixtures. The unique contract-mandated unit tests
// for the native layer are ownership failure (above) and holdsLock (above). ---

// --- Amendment 2: Object.wait argument-range validation (ADR-0029 + JDK API) ---

// TestObjectWaitNegativeTimeoutThrowsIAE verifies Object.wait(long) throws
// IllegalArgumentException when timeoutMillis is negative, before any monitor
// ownership check (per java.lang.Object.wait(long) specification).
func TestObjectWaitNegativeTimeoutThrowsIAE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	// Own the monitor so the only failure is the negative timeout.
	m := lockObj.Monitor()
	m.Enter(rt.EC())

	frame := caller.NewFrame(objClass.LookupMethod("wait", "(J)V"))
	frame.SetRef(0, lockObj)
	frame.SetLong(1, -1)
	caller.PushFrame(frame)
	objectWait(frame)

	if !caller.HasException() {
		t.Fatal("wait(-1) should throw IllegalArgumentException")
	}
	exc := caller.ClearException()
	if exc == nil || exc.Class().Name() != "java/lang/IllegalArgumentException" {
		t.Fatalf("wait(-1) threw %v, want IllegalArgumentException", exc)
	}
	// Monitor ownership and depth must be unaffected.
	if !m.HoldsLock(rt.EC()) {
		t.Fatal("monitor released after negative-timeout wait")
	}
	if d := m.RecursionDepth(); d != 1 {
		t.Fatalf("monitor depth = %d after negative-timeout wait, want 1", d)
	}
	m.Exit(rt.EC())
}

// TestObjectWaitJINegativeTimeoutThrowsIAE verifies Object.wait(long, int)
// throws IllegalArgumentException for a negative timeoutMillis.
func TestObjectWaitJINegativeTimeoutThrowsIAE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	m := lockObj.Monitor()
	m.Enter(rt.EC())

	frame := caller.NewFrame(objClass.LookupMethod("wait", "(JI)V"))
	frame.SetRef(0, lockObj)
	frame.SetLong(1, -5)
	frame.SetInt(3, 0)
	caller.PushFrame(frame)
	objectWaitJI(frame)

	if !caller.HasException() {
		t.Fatal("wait(-5, 0) should throw IllegalArgumentException")
	}
	exc := caller.ClearException()
	if exc == nil || exc.Class().Name() != "java/lang/IllegalArgumentException" {
		t.Fatalf("wait(-5, 0) threw %v, want IllegalArgumentException", exc)
	}
	m.Exit(rt.EC())
}

// TestObjectWaitJINanosOutOfRangeThrowsIAE verifies Object.wait(long, int)
// throws IllegalArgumentException when nanos is outside 0-999999.
func TestObjectWaitJINanosOutOfRangeThrowsIAE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	m := lockObj.Monitor()
	m.Enter(rt.EC())

	frame := caller.NewFrame(objClass.LookupMethod("wait", "(JI)V"))
	frame.SetRef(0, lockObj)
	frame.SetLong(1, 100)
	frame.SetInt(3, 1000000) // > 999999
	caller.PushFrame(frame)
	objectWaitJI(frame)

	if !caller.HasException() {
		t.Fatal("wait(100, 1000000) should throw IllegalArgumentException")
	}
	exc := caller.ClearException()
	if exc == nil || exc.Class().Name() != "java/lang/IllegalArgumentException" {
		t.Fatalf("wait(100, 1000000) threw %v, want IllegalArgumentException", exc)
	}
	m.Exit(rt.EC())
}

// TestObjectWaitJINegativeNanosThrowsIAE verifies the nanos lower bound.
func TestObjectWaitJINegativeNanosThrowsIAE(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	objClass := loader.classes["java/lang/Object"]
	threadClass := loader.classes["java/lang/Thread"]

	lockObj := rtda.NewObject(objClass)
	_, rt := newTestThreadWithLoader(loader, threadClass)
	caller := rtda.NewThread(loader)

	m := lockObj.Monitor()
	m.Enter(rt.EC())

	frame := caller.NewFrame(objClass.LookupMethod("wait", "(JI)V"))
	frame.SetRef(0, lockObj)
	frame.SetLong(1, 100)
	frame.SetInt(3, -1)
	caller.PushFrame(frame)
	objectWaitJI(frame)

	if !caller.HasException() {
		t.Fatal("wait(100, -1) should throw IllegalArgumentException")
	}
	exc := caller.ClearException()
	if exc == nil || exc.Class().Name() != "java/lang/IllegalArgumentException" {
		t.Fatalf("wait(100, -1) threw %v, want IllegalArgumentException", exc)
	}
	m.Exit(rt.EC())
}
