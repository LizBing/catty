package native

import "catty/rtda"

func init() {
	// Exception classes needed by Thread fixtures.
	registerSynthetic("java/lang/IllegalThreadStateException", buildITSE)
	registerSynthetic("java/lang/InterruptedException", buildIE)
}

func buildITSE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/IllegalThreadStateException", "java/lang/IllegalArgumentException", loader)
}

func buildIE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/InterruptedException", "java/lang/Exception", loader)
}

// --- Thread native methods (ADR-0028) ---

// threadInit implements Thread.<init>(). It creates a runtime execution context
// (rtda.Thread) for the new Java Thread object and attaches it bidirectionally.
func threadInit(f *rtda.Frame) {
	this := f.GetRef(0)
	loader := f.Thread().Loader()
	t := rtda.NewThread(loader)
	t.SetJavaThread(this)
	this.SetExtra(t)
}

// threadCurrentThread implements the static Thread.currentThread() native.
// Returns the canonical java.lang.Thread facade attached to the calling
// execution context (ADR-0028: stable identity per Java Thread).
func threadCurrentThread(f *rtda.Frame) {
	f.PushRef(f.Thread().JavaThread())
}

// threadStart implements Thread.start(). It atomically transitions the thread
// from NEW to RUNNABLE (fails if already started), launches a goroutine carrier,
// and runs the Thread's run() method via virtual dispatch.
func threadStart(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		throwNPE(f, "null Thread object")
		return
	}

	extra := this.Extra()
	if extra == nil {
		// Thread was not properly initialized.
		throwIllegalThreadState(f)
		return
	}
	t := extra.(*rtda.Thread)

	// Atomic start-once: CAS NEW → RUNNABLE.
	if !t.SetStarted() {
		throwIllegalThreadState(f)
		return
	}

	vm := rtda.GetVM()
	vm.ThreadStarted(t.IsDaemon())

	// Find run() via virtual dispatch on the actual class of the Thread object
	// (so subclass overrides like Worker.run() are found).
	runMethod := this.Class().LookupMethod("run", "()V")

	go func() {
		defer func() {
			t.Terminate()
			vm.ThreadTerminated(t.IsDaemon())
		}()

		if runMethod != nil {
			frame := t.NewFrame(runMethod)
			frame.SetRef(0, this) // 'this' = the Thread object
			t.PushFrame(frame)
			rtda.DefaultRunLoop(t)
		}
		// If runMethod is nil (shouldn't happen — buildThread adds a native
		// run() and subclasses override it), the goroutine just terminates.
	}()
}

// threadIsAlive implements Thread.isAlive(). Returns true if the thread has
// been started and has not yet terminated.
func threadIsAlive(f *rtda.Frame) {
	this := f.GetRef(0)
	if extra := this.Extra(); extra != nil {
		t := extra.(*rtda.Thread)
		if t.IsAlive() {
			f.PushInt(1)
			return
		}
	}
	f.PushInt(0)
}

// threadJoin implements Thread.join() (untimed). Blocks the calling thread
// until the target thread terminates, or until the calling thread is interrupted.
func threadJoin(f *rtda.Frame) {
	this := f.GetRef(0)
	extra := this.Extra()
	if extra == nil {
		return // not started — return immediately per Java spec
	}
	target := extra.(*rtda.Thread)
	if !target.IsAlive() {
		return // already terminated
	}

	joining := f.Thread()
	for {
		select {
		case <-target.Done():
			return
		case <-joining.Waker():
			if joining.Interrupted() {
				throwInterruptedException(f)
				return
			}
			// Spurious wake — retry.
		}
	}
}

// threadInterrupt implements Thread.interrupt(). Sets the interrupt flag on
// the target thread and signals its waker channel.
func threadInterrupt(f *rtda.Frame) {
	this := f.GetRef(0)
	if extra := this.Extra(); extra != nil {
		t := extra.(*rtda.Thread)
		t.Interrupt()
	}
}

// threadIsInterrupted implements Thread.isInterrupted() (instance method).
// Reads the interrupt flag without clearing it.
func threadIsInterrupted(f *rtda.Frame) {
	this := f.GetRef(0)
	if extra := this.Extra(); extra != nil {
		t := extra.(*rtda.Thread)
		if t.IsInterrupted() {
			f.PushInt(1)
			return
		}
	}
	f.PushInt(0)
}

// threadInterrupted implements the static Thread.interrupted().
// Reads and clears the interrupt flag of the calling thread.
func threadInterrupted(f *rtda.Frame) {
	if f.Thread().Interrupted() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

// threadSleep implements the static Thread.sleep(long millis).
// The calling thread sleeps for the given duration, or until interrupted.
// If interrupted, InterruptedException is thrown and the interrupt flag is cleared.
func threadSleep(f *rtda.Frame) {
	millis := f.GetLong(0) // static method: local 0 = first arg
	if millis < 0 {
		// Negative sleep is an error in Java
		f.Thread().Throw(newException(f.Thread(), "java/lang/IllegalArgumentException", "timeout value is negative"), 0)
		return
	}
	t := f.Thread()
	if !t.Sleep(millis) {
		// Was interrupted — throw InterruptedException.
		throwInterruptedException(f)
	}
}

// threadOnSpinWait implements the static Thread.onSpinWait().
// Java's version is a CPU hint for spin loops; we make it a nop.
func threadOnSpinWait(f *rtda.Frame) {}

// threadSetDaemon implements Thread.setDaemon(boolean).
// May only be called before the thread is started.
func threadSetDaemon(f *rtda.Frame) {
	this := f.GetRef(0)
	v := f.GetInt(1) != 0
	if extra := this.Extra(); extra != nil {
		t := extra.(*rtda.Thread)
		t.SetDaemon(v)
	}
}

// threadIsDaemon implements Thread.isDaemon().
func threadIsDaemon(f *rtda.Frame) {
	this := f.GetRef(0)
	if extra := this.Extra(); extra != nil {
		t := extra.(*rtda.Thread)
		if t.IsDaemon() {
			f.PushInt(1)
			return
		}
	}
	f.PushInt(0)
}

// --- helpers ---

// throwIllegalThreadState creates an IllegalThreadStateException and signals
// it on the calling thread. The interpreter's exception dispatch handles
// propagation after the native method returns.
func throwIllegalThreadState(f *rtda.Frame) {
	obj := newException(f.Thread(), "java/lang/IllegalThreadStateException", "")
	f.Thread().Throw(obj, 0)
}

// throwInterruptedException creates an InterruptedException and signals it on
// the calling thread. The interrupt flag must be cleared before calling this
// (which Sleep/Interrupted already do).
func throwInterruptedException(f *rtda.Frame) {
	obj := newException(f.Thread(), "java/lang/InterruptedException", "")
	f.Thread().Throw(obj, 0)
}

// newException creates a new exception object of the given class with an
// optional detail message. Returns the rtda.Object.
func newException(thread *rtda.Thread, className, message string) *rtda.Object {
	cls := thread.Loader().LoadClass(className)
	obj := rtda.NewObject(cls)
	if message != "" {
		slot := detailMessageSlot(obj)
		msgObj := newStringFromGo(thread, message)
		obj.SetRefCell(int(slot), msgObj)
	}
	return obj
}
