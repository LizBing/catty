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

// TestThreadHoldsLockNull verifies holdsLock returns false for a null object.
func TestThreadHoldsLockNull(t *testing.T) {
	loader := buildObjectMonitorHierarchy()
	threadClass := loader.classes["java/lang/Thread"]

	caller := rtda.NewThread(loader)

	frame := caller.NewFrame(threadClass.LookupMethod("holdsLock", "(Ljava/lang/Object;)Z"))
	frame.SetRef(0, nil)
	threadHoldsLock(frame)

	if frame.PopInt() != 0 {
		t.Fatal("holdsLock(null) should return false (0)")
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
// for the native layer are ownership failure (above) and holdsLock (below). ---
