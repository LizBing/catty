package native

import (
	"sync/atomic"
	"testing"
	"time"

	"catty/rtda"
)

// simpleLoader is a minimal rtda.Loader used in tests to avoid an import cycle
// with catty/classloader (which imports catty/native).
type simpleLoader struct {
	classes map[string]*rtda.Class
	id      *rtda.LoaderIdentity
}

func (l *simpleLoader) LoadClass(name string) *rtda.Class {
	return l.classes[name]
}

func (l *simpleLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	c := l.LoadClass(name)
	if c != nil {
		return rtda.NewClassResult(c)
	}
	return rtda.NewFailureResult(&rtda.ClassLoadFailure{Kind: rtda.FailureNotFound, Name: name})
}

func (l *simpleLoader) LoaderIdentity() *rtda.LoaderIdentity {
	if l.id == nil {
		l.id = rtda.NewLoaderIdentity()
	}
	return l.id
}

// buildTestHierarchy constructs a minimal class hierarchy for Thread tests:
// Object → Throwable → Exception → RuntimeException → IllegalArgumentException → ITSE
// Object → Throwable → Exception → InterruptedException
// Object → Thread (with native methods wired)
func buildTestHierarchy() *simpleLoader {
	l := &simpleLoader{classes: make(map[string]*rtda.Class)}

	obj := rtda.NewSyntheticClass("java/lang/Object", nil)
	l.classes["java/lang/Object"] = obj

	throwable := rtda.NewSyntheticClass("java/lang/Throwable", obj)
	l.classes["java/lang/Throwable"] = throwable

	exc := rtda.NewSyntheticClass("java/lang/Exception", throwable)
	l.classes["java/lang/Exception"] = exc

	re := rtda.NewSyntheticClass("java/lang/RuntimeException", exc)
	l.classes["java/lang/RuntimeException"] = re

	iae := rtda.NewSyntheticClass("java/lang/IllegalArgumentException", re)
	l.classes["java/lang/IllegalArgumentException"] = iae

	itse := rtda.NewSyntheticClass("java/lang/IllegalThreadStateException", iae)
	l.classes["java/lang/IllegalThreadStateException"] = itse

	ie := rtda.NewSyntheticClass("java/lang/InterruptedException", exc)
	l.classes["java/lang/InterruptedException"] = ie

	thread := rtda.NewSyntheticClass("java/lang/Thread", obj)
	thread.AddMethod(rtda.NativeMethod(thread, "<init>", "()V", threadInit))
	thread.AddMethod(rtda.NativeMethod(thread, "run", "()V", nop))
	thread.AddMethod(rtda.NativeMethod(thread, "start", "()V", threadStart))
	thread.AddMethod(rtda.NativeMethod(thread, "isAlive", "()Z", threadIsAlive))
	thread.AddMethod(rtda.NativeMethod(thread, "join", "()V", threadJoin))
	thread.AddMethod(rtda.NativeMethod(thread, "interrupt", "()V", threadInterrupt))
	thread.AddMethod(rtda.NativeMethod(thread, "isInterrupted", "()Z", threadIsInterrupted))
	thread.AddMethod(rtda.NativeMethod(thread, "setDaemon", "(Z)V", threadSetDaemon))
	thread.AddMethod(rtda.NativeMethod(thread, "isDaemon", "()Z", threadIsDaemon))
	thread.AddMethod(staticNative(thread, "currentThread", "()Ljava/lang/Thread;", threadCurrentThread))
	thread.AddMethod(staticNative(thread, "interrupted", "()Z", threadInterrupted))
	thread.AddMethod(staticNative(thread, "sleep", "(J)V", threadSleep))
	l.classes["java/lang/Thread"] = thread

	return l
}

// testRunLoop is a minimal interpreter loop that executes native methods directly.
// The goroutine carrier in threadStart calls this via rtda.DefaultRunLoop.
func testRunLoop(t *rtda.Thread) {
	for !t.IsStackEmpty() {
		f := t.CurrentFrame()
		m := f.Method()
		if m.IsNative() {
			m.NativeFunc()(f)
		}
		t.PopFrame()
	}
}

// newTestThreadObj creates an initialized java.lang.Thread object: builds the
// hierarchy, creates the Object, and runs <init> to attach the rtda.Thread.
func newTestThreadObj() (*simpleLoader, *rtda.Class, *rtda.Object) {
	loader := buildTestHierarchy()
	threadClass := loader.classes["java/lang/Thread"]
	obj := rtda.NewObject(threadClass)

	caller := rtda.NewThread(loader)
	frame := caller.NewFrame(threadClass.LookupMethod("<init>", "()V"))
	frame.SetRef(0, obj)
	threadInit(frame)

	return loader, threadClass, obj
}

// TestThreadInitAttachesRuntimeRecord verifies that <init> creates and
// bidirectionally attaches the rtda.Thread execution context.
func TestThreadInitAttachesRuntimeRecord(t *testing.T) {
	_, _, obj := newTestThreadObj()

	extra := obj.Extra()
	if extra == nil {
		t.Fatal("threadInit did not set Extra on Thread object")
	}
	rt := extra.(*rtda.Thread)
	if rt.JavaThread() != obj {
		t.Error("rtda.Thread.JavaThread should point back to the facade object")
	}
	if rt.IsAlive() {
		t.Error("newly initialized thread should not be alive")
	}
}

// TestThreadCurrentThreadIdentity verifies currentThread returns the canonical
// facade object attached to the calling execution context.
func TestThreadCurrentThreadIdentity(t *testing.T) {
	_, threadClass, facade := newTestThreadObj()

	caller := rtda.NewThread(nil)
	caller.SetJavaThread(facade)

	frame := caller.NewFrame(threadClass.LookupMethod("currentThread", "()Ljava/lang/Thread;"))
	threadCurrentThread(frame)

	result := frame.PopRef()
	if result != facade {
		t.Error("currentThread must return the calling thread's canonical facade")
	}
}

// TestThreadIsAlive verifies isAlive returns correct values across lifecycle states.
func TestThreadIsAlive(t *testing.T) {
	_, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(nil)

	// Not started → not alive.
	frame := caller.NewFrame(threadClass.LookupMethod("isAlive", "()Z"))
	frame.SetRef(0, obj)
	threadIsAlive(frame)
	if frame.PopInt() != 0 {
		t.Error("isAlive should return false before start")
	}

	// After SetStarted → alive.
	rt.SetStarted()
	frame = caller.NewFrame(threadClass.LookupMethod("isAlive", "()Z"))
	frame.SetRef(0, obj)
	threadIsAlive(frame)
	if frame.PopInt() != 1 {
		t.Error("isAlive should return true after start")
	}

	// After Terminate → not alive.
	rt.Terminate()
	frame = caller.NewFrame(threadClass.LookupMethod("isAlive", "()Z"))
	frame.SetRef(0, obj)
	threadIsAlive(frame)
	if frame.PopInt() != 0 {
		t.Error("isAlive should return false after termination")
	}
}

// TestThreadStartAndJoinHappyPath verifies the basic start→join lifecycle:
// launching a goroutine carrier, executing run(), and unblocking join.
func TestThreadStartAndJoinHappyPath(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)
	rtda.DefaultRunLoop = testRunLoop

	loader, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(loader)

	// Start the thread.
	startFrame := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	startFrame.SetRef(0, obj)
	threadStart(startFrame)

	if caller.HasException() {
		t.Fatal("start should not throw on first call")
	}

	// The goroutine may or may not have completed already — join should handle both.
	joinFrame := caller.NewFrame(threadClass.LookupMethod("join", "()V"))
	joinFrame.SetRef(0, obj)

	joinDone := make(chan struct{})
	go func() {
		threadJoin(joinFrame)
		close(joinDone)
	}()

	select {
	case <-joinDone:
		// expected
	case <-time.After(3 * time.Second):
		t.Fatal("join timed out — goroutine carrier may not have terminated")
	}

	// After join, the thread should not be alive.
	if rt.IsAlive() {
		t.Error("thread should not be alive after join completes")
	}

	// WaitForNonDaemonThreads should unblock — the thread is non-daemon by default.
	vmDone := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(vmDone)
	}()
	select {
	case <-vmDone:
		// expected
	case <-time.After(time.Second):
		t.Error("VM should unblock after all non-daemon threads terminate")
	}
}

// TestThreadStartTwiceThrowsITSE verifies that calling start() on an already-started
// thread throws IllegalThreadStateException.
func TestThreadStartTwiceThrowsITSE(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)
	rtda.DefaultRunLoop = testRunLoop

	loader, threadClass, obj := newTestThreadObj()
	caller := rtda.NewThread(loader)

	// First start — should succeed.
	f1 := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	f1.SetRef(0, obj)
	threadStart(f1)
	if caller.HasException() {
		t.Fatal("first start should not throw, got exception")
	}

	// Wait for the goroutine carrier to finish so Terminate is called.
	rt := obj.Extra().(*rtda.Thread)
	select {
	case <-rt.Done():
	case <-time.After(time.Second):
		t.Fatal("thread did not terminate in time")
	}

	// Clear any stale state.
	caller.ClearException()

	// Second start — should throw ITSE.
	f2 := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	f2.SetRef(0, obj)
	threadStart(f2)
	if !caller.HasException() {
		t.Fatal("second start should throw IllegalThreadStateException")
	}
}

// TestThreadInterruptFlag verifies interrupt/isInterrupted/interrupted semantics.
func TestThreadInterruptFlag(t *testing.T) {
	_, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(nil)

	// Initially not interrupted.
	frame := caller.NewFrame(threadClass.LookupMethod("isInterrupted", "()Z"))
	frame.SetRef(0, obj)
	threadIsInterrupted(frame)
	if frame.PopInt() != 0 {
		t.Error("isInterrupted should initially return false")
	}

	// Interrupt the target thread (via native instance method).
	intFrame := caller.NewFrame(threadClass.LookupMethod("interrupt", "()V"))
	intFrame.SetRef(0, obj)
	threadInterrupt(intFrame)

	if !rt.IsInterrupted() {
		t.Error("thread should be interrupted after interrupt()")
	}

	// isInterrupted reads without clearing.
	frame = caller.NewFrame(threadClass.LookupMethod("isInterrupted", "()Z"))
	frame.SetRef(0, obj)
	threadIsInterrupted(frame)
	if frame.PopInt() != 1 {
		t.Error("isInterrupted should return true after interrupt")
	}
	// Second read should still return true (no clear).
	frame = caller.NewFrame(threadClass.LookupMethod("isInterrupted", "()Z"))
	frame.SetRef(0, obj)
	threadIsInterrupted(frame)
	if frame.PopInt() != 1 {
		t.Error("isInterrupted should still return true (does not clear)")
	}

	// Static interrupted() checks the CALLING thread, not the target.
	// Interrupt the caller to test the static method.
	caller.Interrupt()
	staticFrame := caller.NewFrame(threadClass.LookupMethod("interrupted", "()Z"))
	threadInterrupted(staticFrame)
	if staticFrame.PopInt() != 1 {
		t.Error("interrupted() should return true when calling thread is interrupted")
	}
	// Now it should be cleared on the caller.
	if caller.IsInterrupted() {
		t.Error("interrupted() should clear the calling thread's interrupt flag")
	}
	// Target thread's flag should be unaffected.
	if !rt.IsInterrupted() {
		t.Error("interrupted() on caller should not clear target's interrupt flag")
	}
}

// TestThreadInterruptOnCallingThread verifies static interrupted() works on the
// calling thread (not a target thread).
func TestThreadInterruptOnCallingThread(t *testing.T) {
	_, threadClass, _ := newTestThreadObj()
	caller := rtda.NewThread(nil)
	caller.Interrupt()

	frame := caller.NewFrame(threadClass.LookupMethod("interrupted", "()Z"))
	threadInterrupted(frame)
	if frame.PopInt() != 1 {
		t.Error("interrupted() should return true after interrupt on calling thread")
	}
	if caller.IsInterrupted() {
		t.Error("interrupted() should clear calling thread's flag")
	}
}

// TestThreadSleepNormal verifies uninterrupted sleep returns normally.
func TestThreadSleepNormal(t *testing.T) {
	loader, threadClass, _ := newTestThreadObj()
	caller := rtda.NewThread(loader)

	frame := caller.NewFrame(threadClass.LookupMethod("sleep", "(J)V"))
	frame.SetLong(0, 10) // 10ms sleep
	threadSleep(frame)

	if caller.HasException() {
		t.Error("normal sleep should not throw InterruptedException")
	}
}

// TestThreadSleepInterruptedBefore verifies sleep throws IE if already interrupted.
func TestThreadSleepInterruptedBefore(t *testing.T) {
	loader, threadClass, _ := newTestThreadObj()
	caller := rtda.NewThread(loader)
	caller.Interrupt()

	frame := caller.NewFrame(threadClass.LookupMethod("sleep", "(J)V"))
	frame.SetLong(0, 1000)
	threadSleep(frame)

	if !caller.HasException() {
		t.Error("sleep should throw InterruptedException when already interrupted")
	}
}

// TestThreadSleepInterruptedDuring verifies sleep is interrupted mid-sleep.
func TestThreadSleepInterruptedDuring(t *testing.T) {
	loader, threadClass, _ := newTestThreadObj()
	caller := rtda.NewThread(loader)

	var sleepDone int32
	go func() {
		frame := caller.NewFrame(threadClass.LookupMethod("sleep", "(J)V"))
		frame.SetLong(0, 60000) // long sleep
		threadSleep(frame)
		atomic.StoreInt32(&sleepDone, 1)
	}()

	time.Sleep(10 * time.Millisecond) // let sleep start
	caller.Interrupt()

	// Wait for sleep to return.
	deadline := time.After(time.Second)
	for atomic.LoadInt32(&sleepDone) == 0 {
		select {
		case <-deadline:
			t.Fatal("sleep did not return after interrupt")
		case <-time.After(time.Millisecond):
		}
	}

	if !caller.HasException() {
		t.Error("interrupted sleep should throw InterruptedException")
	}
}

// TestThreadSetDaemonIsDaemon verifies setDaemon/isDaemon native methods.
func TestThreadSetDaemonIsDaemon(t *testing.T) {
	_, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(nil)

	// Default: not daemon.
	frame := caller.NewFrame(threadClass.LookupMethod("isDaemon", "()Z"))
	frame.SetRef(0, obj)
	threadIsDaemon(frame)
	if frame.PopInt() != 0 {
		t.Error("new thread should not be daemon by default")
	}

	// Set daemon.
	setFrame := caller.NewFrame(threadClass.LookupMethod("setDaemon", "(Z)V"))
	setFrame.SetRef(0, obj)
	setFrame.SetInt(1, 1)
	threadSetDaemon(setFrame)

	if !rt.IsDaemon() {
		t.Error("setDaemon(true) should mark thread as daemon")
	}

	// Verify isDaemon returns true.
	frame = caller.NewFrame(threadClass.LookupMethod("isDaemon", "()Z"))
	frame.SetRef(0, obj)
	threadIsDaemon(frame)
	if frame.PopInt() != 1 {
		t.Error("isDaemon should return true after setDaemon(true)")
	}
}

// TestThreadSetDaemonAfterStartThrowsITSE verifies that calling setDaemon
// after a successful start throws IllegalThreadStateException.
func TestThreadSetDaemonAfterStartThrowsITSE(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)
	rtda.DefaultRunLoop = testRunLoop

	loader, threadClass, obj := newTestThreadObj()
	caller := rtda.NewThread(loader)

	// Start the thread first.
	startFrame := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	startFrame.SetRef(0, obj)
	threadStart(startFrame)
	if caller.HasException() {
		t.Fatal("first start should not throw")
	}

	// Now try to setDaemon — must throw ITSE.
	setFrame := caller.NewFrame(threadClass.LookupMethod("setDaemon", "(Z)V"))
	setFrame.SetRef(0, obj)
	setFrame.SetInt(1, 1) // true
	threadSetDaemon(setFrame)

	if !caller.HasException() {
		t.Error("setDaemon after start should throw IllegalThreadStateException")
	}

	// Cleanup: wait for the goroutine carrier.
	rt := obj.Extra().(*rtda.Thread)
	select {
	case <-rt.Done():
	case <-time.After(time.Second):
		t.Fatal("thread did not terminate")
	}
}

// TestThreadSetDaemonAfterTerminateThrowsITSE verifies that calling setDaemon
// after the thread has terminated throws IllegalThreadStateException.
func TestThreadSetDaemonAfterTerminateThrowsITSE(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)
	rtda.DefaultRunLoop = testRunLoop

	loader, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(loader)

	// Start and wait for termination.
	startFrame := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	startFrame.SetRef(0, obj)
	threadStart(startFrame)
	if caller.HasException() {
		t.Fatal("start should not throw")
	}

	select {
	case <-rt.Done():
	case <-time.After(time.Second):
		t.Fatal("thread did not terminate")
	}

	// Now try to setDaemon — must throw ITSE.
	setFrame := caller.NewFrame(threadClass.LookupMethod("setDaemon", "(Z)V"))
	setFrame.SetRef(0, obj)
	setFrame.SetInt(1, 1) // true
	threadSetDaemon(setFrame)

	if !caller.HasException() {
		t.Error("setDaemon after termination should throw IllegalThreadStateException")
	}
}

// TestThreadSetDaemonBeforeStartSucceeds verifies setDaemon works when called
// before start (state == NEW).
func TestThreadSetDaemonBeforeStartSucceeds(t *testing.T) {
	_, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	caller := rtda.NewThread(nil)

	// Set daemon before start — should succeed.
	setFrame := caller.NewFrame(threadClass.LookupMethod("setDaemon", "(Z)V"))
	setFrame.SetRef(0, obj)
	setFrame.SetInt(1, 1) // true
	threadSetDaemon(setFrame)

	if caller.HasException() {
		t.Error("setDaemon before start should not throw")
	}
	if !rt.IsDaemon() {
		t.Error("daemon should be true after setDaemon(true)")
	}

	// Also verify we can change it back before start.
	setFrame2 := caller.NewFrame(threadClass.LookupMethod("setDaemon", "(Z)V"))
	setFrame2.SetRef(0, obj)
	setFrame2.SetInt(1, 0) // false
	threadSetDaemon(setFrame2)

	if caller.HasException() {
		t.Error("setDaemon(false) before start should not throw")
	}
	if rt.IsDaemon() {
		t.Error("daemon should be false after setDaemon(false)")
	}
}

// TestVMNonDaemonKeepsAlive verifies a started non-daemon thread keeps the VM
// from exiting until it terminates.
func TestVMNonDaemonKeepsAlive(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)

	vm.ThreadStarted(false)

	vmDone := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(vmDone)
	}()

	// Should NOT unblock — count is 1.
	select {
	case <-vmDone:
		t.Error("VM should not unblock while non-daemon thread is running")
	case <-time.After(50 * time.Millisecond):
	}

	// Terminate the non-daemon thread.
	vm.ThreadTerminated(false)

	select {
	case <-vmDone:
		// expected
	case <-time.After(time.Second):
		t.Error("VM should unblock after last non-daemon terminates")
	}
}

// TestVMDaemonDoesNotKeepAlive verifies daemon threads don't prevent VM exit.
func TestVMDaemonDoesNotKeepAlive(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)

	vm.ThreadStarted(true)

	vmDone := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(vmDone)
	}()

	// Should unblock immediately — daemon threads don't count.
	select {
	case <-vmDone:
		// expected
	case <-time.After(time.Second):
		t.Error("daemon thread should not keep VM alive")
	}
}

// TestThreadStartDaemonLiveness verifies a daemon thread does not keep VM alive
// when started via threadStart.
func TestThreadStartDaemonLiveness(t *testing.T) {
	vm := rtda.NewVM()
	rtda.SetVM(vm)
	rtda.DefaultRunLoop = testRunLoop

	loader, threadClass, obj := newTestThreadObj()
	rt := obj.Extra().(*rtda.Thread)
	rt.SetDaemon(true)

	caller := rtda.NewThread(loader)

	// Start the daemon thread.
	startFrame := caller.NewFrame(threadClass.LookupMethod("start", "()V"))
	startFrame.SetRef(0, obj)
	threadStart(startFrame)

	if caller.HasException() {
		t.Fatal("start should not throw")
	}

	// VM should be free to exit — daemon thread doesn't count.
	vmDone := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(vmDone)
	}()

	select {
	case <-vmDone:
		// expected — daemon threads don't keep VM alive
	case <-time.After(time.Second):
		t.Error("VM should be able to exit with only daemon threads")
	}

	// Clean up: wait for the goroutine carrier to finish.
	select {
	case <-rt.Done():
	case <-time.After(time.Second):
		t.Fatal("daemon thread did not terminate")
	}
}
