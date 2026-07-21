package interpreter

import (
	"sync"
	"testing"
	"time"

	"catty/rtda"
)

// testLoader is a minimal rtda.Loader for interpreter exception tests.
// It returns failures for specific class names and resolved classes for others.
type testLoader struct {
	mu           sync.Mutex
	failingNames map[string]rtda.FailureKind
	classes      map[string]*rtda.Class
	id           *rtda.LoaderIdentity
}

func newTestLoader() *testLoader {
	return &testLoader{
		failingNames: make(map[string]rtda.FailureKind),
		classes:      make(map[string]*rtda.Class),
		id:           rtda.NewLoaderIdentity(),
	}
}

func (l *testLoader) addFailing(name string, kind rtda.FailureKind) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.failingNames[name] = kind
}

func (l *testLoader) addClass(c *rtda.Class) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.classes[c.Name()] = c
}

func (l *testLoader) LoadClass(name string) *rtda.Class {
	l.mu.Lock()
	defer l.mu.Unlock()
	if c, ok := l.classes[name]; ok {
		return c
	}
	return nil
}

func (l *testLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	l.mu.Lock()
	defer l.mu.Unlock()
	if c, ok := l.classes[name]; ok {
		return rtda.NewClassResult(c)
	}
	if kind, ok := l.failingNames[name]; ok {
		return rtda.NewFailureResult(&rtda.ClassLoadFailure{Kind: kind, Name: name})
	}
	return rtda.NewFailureResult(&rtda.ClassLoadFailure{Kind: rtda.FailureNotFound, Name: name})
}

func (l *testLoader) LoaderIdentity() *rtda.LoaderIdentity {
	return l.id
}

// setupBaseClasses creates the minimal class hierarchy needed for
// resolveClass to construct error throwables (NoClassDefFoundError, etc.).
// All classes are registered with the loader.
func setupBaseClasses(loader *testLoader) {
	obj := rtda.NewSyntheticClass("java/lang/Object", nil)
	loader.addClass(obj)

	throwable := rtda.NewSyntheticClass("java/lang/Throwable", obj)
	loader.addClass(throwable)

	error_ := rtda.NewSyntheticClass("java/lang/Error", throwable)
	loader.addClass(error_)

	linkageError := rtda.NewSyntheticClass("java/lang/LinkageError", error_)
	loader.addClass(linkageError)

	ncdfe := rtda.NewSyntheticClass("java/lang/NoClassDefFoundError", linkageError)
	loader.addClass(ncdfe)

	classFormatError := rtda.NewSyntheticClass("java/lang/ClassFormatError", linkageError)
	loader.addClass(classFormatError)

	classCircularityError := rtda.NewSyntheticClass("java/lang/ClassCircularityError", linkageError)
	loader.addClass(classCircularityError)

	runtimeException := rtda.NewSyntheticClass("java/lang/RuntimeException", throwable)
	loader.addClass(runtimeException)
	npe := rtda.NewSyntheticClass("java/lang/NullPointerException", runtimeException)
	loader.addClass(npe)
}

// makeTestClass creates a synthetic class with a single interpreted method
// that has the given exception table entries. The method has minimal bytecode
// (just a 'return' instruction at PC 0).
func makeTestClass(name string, methodName string, exTable []rtda.ExceptionEntry, loader *testLoader) *rtda.Class {
	return makeTestClassWithCode(name, methodName, []byte{0xb1}, exTable, loader)
}

func makeTestClassWithCode(name, methodName string, code []byte, exTable []rtda.ExceptionEntry, loader *testLoader) *rtda.Class {
	cls := rtda.NewSyntheticClass(name, nil)
	// Method descriptor: static void ().
	method := rtda.InterpretedMethod(cls, methodName, "()V",
		0x0009, // ACC_PUBLIC | ACC_STATIC
		1, 0,   // maxStack=1, maxLocals=0
		code,
		exTable,
	)
	cls.AddMethod(method)
	loader.addClass(cls)
	return cls
}

// makeExceptionObject creates a minimal exception object of the given class.
func makeExceptionObject(cls *rtda.Class) *rtda.Object {
	return rtda.NewObject(cls)
}

// newThreadWithLoader creates a thread with the given loader.
func newThreadWithLoader(loader rtda.Loader) *rtda.Thread {
	return rtda.NewThread(loader)
}

// pushFrame pushes a frame for the given method onto the thread.
func pushFrame(thread *rtda.Thread, method *rtda.Method) {
	frame := thread.NewFrame(method)
	thread.PushFrame(frame)
}

// TestHandleExceptionMissingCatchType verifies that when catch-type resolution
// fails, the current frame is popped and the replacement throwable propagates
// to the caller.
func TestHandleExceptionMissingCatchType(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)

	// The catch type "test/NonExistent" will fail resolution.
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	// Create the throwing class with an exception table entry for the failing type.
	exTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "test/NonExistent"), // catch type that fails
	}
	throwingClass := makeTestClass("test/ThrowingClass", "throwingMethod", exTable, loader)
	throwingMethod := throwingClass.LookupMethod("throwingMethod", "()V")
	if throwingMethod == nil {
		t.Fatal("throwing method not found")
	}

	// Create caller class with no exception table (will receive propagated exception).
	callerClass := makeTestClass("test/CallerClass", "callerMethod", nil, loader)
	callerMethod := callerClass.LookupMethod("callerMethod", "()V")
	if callerMethod == nil {
		t.Fatal("caller method not found")
	}

	// Create an exception object of a class that IS loadable.
	excClass := rtda.NewSyntheticClass("test/TestException", nil)
	loader.addClass(excClass)
	excObj := rtda.NewObject(excClass)

	// Set up thread stack: caller frame -> throwing frame.
	thread := newThreadWithLoader(loader)
	pushFrame(thread, callerMethod)
	pushFrame(thread, throwingMethod)

	// Throw the exception at PC 0 (in range of the failing catch type entry).
	thread.Throw(excObj, 0)

	// handleException should: resolve catch type (fails), pop throwing frame,
	// continue frame-walk with replacement throwable in caller frame.
	handleException(thread, 0)

	// Throwing frame should be popped.
	if thread.IsStackEmpty() {
		// The replacement throwable propagated past all frames.
		// This means the caller didn't catch it either, which is fine —
		// the test verifies the throwing frame was popped.
		return
	}

	// If we're here, there are frames left. The current frame should be
	// the caller frame, NOT the throwing frame.
	currentFrame := thread.CurrentFrame()
	if currentFrame.Method() == throwingMethod {
		t.Error("throwing frame was not popped after catch-type resolution failure")
	}
}

// TestCatchTypeFailureNotCaughtByLaterCatchAll verifies that when a failing
// catch-type entry is followed by a catch-all in the same frame, the catch-all
// does NOT capture the replacement throwable at the original throwPC.
func TestCatchTypeFailureNotCaughtByLaterCatchAll(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)

	// The typed catch type fails resolution.
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	// Exception table: failing typed entry FIRST, then catch-all SECOND.
	// Both cover PC 0. The catch-all must NOT capture the replacement.
	exTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "test/NonExistent"), // fails
		rtda.NewExceptionEntry(0, 1, 20, ""),                 // catch-all
	}
	throwingClass := makeTestClass("test/ThrowingClass2", "throwingMethod", exTable, loader)
	throwingMethod := throwingClass.LookupMethod("throwingMethod", "()V")
	if throwingMethod == nil {
		t.Fatal("throwing method not found")
	}

	// Caller class with no exception table.
	callerClass := makeTestClass("test/CallerClass2", "callerMethod", nil, loader)
	callerMethod := callerClass.LookupMethod("callerMethod", "()V")
	if callerMethod == nil {
		t.Fatal("caller method not found")
	}

	excClass := rtda.NewSyntheticClass("test/TestException2", nil)
	loader.addClass(excClass)
	excObj := rtda.NewObject(excClass)

	// Thread: caller -> throwing.
	thread := newThreadWithLoader(loader)
	pushFrame(thread, callerMethod)
	pushFrame(thread, throwingMethod)

	thread.Throw(excObj, 0)
	handleException(thread, 0)

	// After handleException, the throwing frame should be popped.
	// The catch-all must NOT have captured the replacement throwable.
	if thread.IsStackEmpty() {
		return // propagated past all frames — correct
	}

	currentFrame := thread.CurrentFrame()
	if currentFrame.Method() == throwingMethod {
		t.Error("catch-all incorrectly captured replacement throwable; frame was not popped")
	}
	// Verify the exception is NOT pending on the thread (it was cleared and
	// either caught by caller or propagated).
	if thread.HasException() {
		t.Error("thread still has pending exception after handleException returned")
	}
}

// TestCatchTypeFailureCallerCatches verifies that when catch-type resolution
// fails and the frame is popped, the caller can catch the replacement throwable.
func TestCatchTypeFailureCallerCatches(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)

	// The typed catch type fails.
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	// Throwing frame: has failing catch type entry.
	exTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "test/NonExistent"),
	}
	throwingClass := makeTestClass("test/ThrowingClass3", "throwingMethod", exTable, loader)
	throwingMethod := throwingClass.LookupMethod("throwingMethod", "()V")
	if throwingMethod == nil {
		t.Fatal("throwing method not found")
	}

	// Caller frame: catches NoClassDefFoundError (the replacement throwable type).
	callerExTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "java/lang/NoClassDefFoundError"),
	}
	callerClass := makeTestClass("test/CallerClass3", "callerMethod", callerExTable, loader)
	callerMethod := callerClass.LookupMethod("callerMethod", "()V")
	if callerMethod == nil {
		t.Fatal("caller method not found")
	}

	excClass := rtda.NewSyntheticClass("test/TestException3", nil)
	loader.addClass(excClass)
	excObj := rtda.NewObject(excClass)

	// Thread: caller -> throwing.
	thread := newThreadWithLoader(loader)
	pushFrame(thread, callerMethod)
	// Set caller PC to 1 so that after popping the throwing frame,
	// throwPC = 1 - 1 = 0, which is in range of the exception table.
	thread.CurrentFrame().SetPC(1)
	pushFrame(thread, throwingMethod)

	thread.Throw(excObj, 0)
	handleException(thread, 0)

	// After handleException: throwing frame popped. The replacement throwable
	// (NoClassDefFoundError) should be caught by caller's handler.
	// The current frame should still be the caller frame (it caught the exception).
	if thread.IsStackEmpty() {
		t.Error("stack is empty — caller should have caught the replacement throwable")
		return
	}
	currentFrame := thread.CurrentFrame()
	if currentFrame.Method() != callerMethod {
		t.Errorf("current frame method = %s, want callerMethod", currentFrame.Method().Name())
	}
	// PC should be set to handler PC (10).
	if currentFrame.PC() != 10 {
		t.Errorf("caller PC = %d, want 10 (handler PC)", currentFrame.PC())
	}
}

// TestClinitMissingCatchType verifies that when catch-type resolution fails
// inside a <clinit> method, the frame is popped and the replacement propagates.
// When it propagates past the clinit boundary, runClinit returns InitResult with
// the replacement throwable.
func TestClinitMissingCatchType(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)

	// Catch type fails resolution.
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	// Clinit method with failing catch type entry.
	exTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "test/NonExistent"),
	}
	clinitClass := makeTestClass("test/ClinitClass", "<clinit>", exTable, loader)
	clinitMethod := clinitClass.LookupMethod("<clinit>", "()V")
	if clinitMethod == nil {
		t.Fatal("<clinit> method not found")
	}

	excClass := rtda.NewSyntheticClass("test/TestException4", nil)
	loader.addClass(excClass)
	excObj := rtda.NewObject(excClass)

	// Set up thread with just the clinit frame.
	thread := newThreadWithLoader(loader)

	// runClinit pushes its own frame. But we need to set up the exception
	// inside the clinit. Since runClinit exec's bytecode, we need the bytecode
	// to throw. Instead, we test directly by simulating: push the clinit frame,
	// throw an exception, and check that the catch-type failure causes the
	// frame to be popped.

	// Simulate the clinit execution state:
	pushFrame(thread, clinitMethod)
	thread.Throw(excObj, 0)

	// The clinit exception handler should pop the clinit frame on catch-type failure.
	handleException(thread, 0)

	// The clinit frame should be popped (no other frames exist).
	if !thread.IsStackEmpty() {
		t.Error("clinit frame should have been popped after catch-type failure")
	}
}

func TestRunClinitMissingCatchType(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	// aconst_null; athrow; return. The athrow produces NPE at PC 1, then
	// resolution of the protected catch type produces NoClassDefFoundError.
	exTable := []rtda.ExceptionEntry{rtda.NewExceptionEntry(0, 2, 2, "test/NonExistent")}
	class := makeTestClassWithCode("test/RealClinit", "<clinit>", []byte{0x01, 0xbf, 0xb1}, exTable, loader)
	method := class.LookupMethod("<clinit>", "()V")
	thread := rtda.NewThread(loader)

	done := make(chan rtda.InitResult, 1)
	go func() { done <- runClinit(thread, method) }()
	select {
	case result := <-done:
		if result.ErrObj == nil {
			t.Fatal("runClinit unexpectedly completed normally")
		}
		if got := result.ErrObj.Class().Name(); got != "java/lang/NoClassDefFoundError" {
			t.Fatalf("replacement throwable = %s, want java/lang/NoClassDefFoundError", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runClinit hung resolving a missing catch type")
	}
}

func testCatchResolutionLoop(t *testing.T, loop func(*rtda.Thread)) {
	t.Helper()
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	throwTable := []rtda.ExceptionEntry{rtda.NewExceptionEntry(0, 2, 2, "test/NonExistent")}
	throwClass := makeTestClassWithCode("test/LoopThrow", "thrower", []byte{0x01, 0xbf, 0xb1}, throwTable, loader)
	callerTable := []rtda.ExceptionEntry{rtda.NewExceptionEntry(0, 1, 1, "java/lang/NoClassDefFoundError")}
	callerClass := makeTestClassWithCode("test/LoopCaller", "caller", []byte{0xb1, 0xb1}, callerTable, loader)

	thread := rtda.NewThread(loader)
	callerFrame := thread.NewFrame(callerClass.LookupMethod("caller", "()V"))
	callerFrame.SetPC(1)
	thread.PushFrame(callerFrame)
	thread.PushFrame(thread.NewFrame(throwClass.LookupMethod("thrower", "()V")))

	done := make(chan any, 1)
	go func() {
		defer func() { done <- recover() }()
		loop(thread)
	}()
	select {
	case panicValue := <-done:
		if panicValue != nil {
			t.Fatalf("execution panicked: %v", panicValue)
		}
		if !thread.IsStackEmpty() {
			t.Fatal("execution did not drain the caller handler frame")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("execution hung resolving a missing catch type")
	}
}

func TestLoopCatchTypeResolutionFailure(t *testing.T) {
	testCatchResolutionLoop(t, Loop)
}

func TestLoopIRCatchTypeResolutionFailure(t *testing.T) {
	testCatchResolutionLoop(t, LoopIR)
}

// TestHandleExceptionDeadlockFree verifies that catch-type resolution failure
// does not deadlock when the loader itself delegates back to the same thread's
// loader context (a cross-dependency scenario).
func TestHandleExceptionDeadlockFree(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addFailing("test/NonExistent", rtda.FailureNotFound)

	exTable := []rtda.ExceptionEntry{
		rtda.NewExceptionEntry(0, 1, 10, "test/NonExistent"),
	}
	throwingClass := makeTestClass("test/DeadlockFree", "m", exTable, loader)
	throwingMethod := throwingClass.LookupMethod("m", "()V")
	if throwingMethod == nil {
		t.Fatal("method not found")
	}

	callerClass := makeTestClass("test/DeadlockCaller", "caller", nil, loader)
	callerMethod := callerClass.LookupMethod("caller", "()V")
	if callerMethod == nil {
		t.Fatal("caller method not found")
	}

	excClass := rtda.NewSyntheticClass("test/TestExc", nil)
	loader.addClass(excClass)
	excObj := rtda.NewObject(excClass)

	done := make(chan struct{})
	go func() {
		defer close(done)
		thread := newThreadWithLoader(loader)
		pushFrame(thread, callerMethod)
		pushFrame(thread, throwingMethod)
		thread.Throw(excObj, 0)
		handleException(thread, 0)
	}()

	select {
	case <-done:
		// Pass — no deadlock.
	case <-time.After(2 * time.Second):
		t.Fatal("handleException deadlocked after catch-type resolution failure")
	}
}
