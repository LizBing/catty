package rtda

import (
	"sync/atomic"
	"time"
)

// Loader is the subset of the classloader that rtda needs at run time for class
// resolution (new, anewarray, checkcast, ldc class, invokeinterface, ...).
// Declaring it here as an interface keeps rtda free of any import cycle with the
// classloader package, which implements Loader concretely.
type Loader interface {
	LoadClass(name string) *Class
}

// thread states (atomic int32).
const (
	stateNew int32 = iota
	stateRunnable
	stateTerminated
)

// Thread models a single JVM execution thread's stack of frames (JVMS §2.5.1)
// plus its Java-level identity, lifecycle, interrupt status, and daemon state
// (ADR-0028).
type Thread struct {
	stack  []*Frame
	loader Loader
	// ecID identifies the execution context (ADR-0025). In the single-context
	// interpreter it is a fixed sentinel; AOT and future multi-threaded runtimes
	// assign a distinct value per Java thread so recursive same-owner <clinit>
	// requests complete normally without re-running <clinit>.
	ecID uint64
	// bridgeReturn captures a method's return value when run from the AOT bridge
	// (interpreter.RunMethod): there is no caller frame, so the return helpers
	// write here instead of pushing. nil outside bridge mode.
	bridgeReturn *Slot
	// pendingException is non-nil when an exception is in flight (athrow or a
	// runtime error like NPE). The interpreter Loop checks HasException after
	// each instruction and dispatches to handleException.
	pendingException *Object
	throwPC          int // PC of the instruction that threw (for exception-table search)

	// --- Slice B: Thread identity, lifecycle, interrupt, and daemon ---

	// javaThread is the canonical java.lang.Thread facade object attached to
	// this execution context. currentThread() returns this object.
	javaThread *Object
	// state is the lifecycle state (stateNew / stateRunnable / stateTerminated).
	// Managed with atomic ops — SetStarted uses CAS; Terminate and IsAlive use
	// Load/Store.
	state int32
	// interruptState is 0 (clear) or 1 (interrupt pending). Managed with atomic
	// ops: Interrupt stores 1 + signals waker; IsInterrupted Loads; Interrupted
	// Swaps to 0.
	interruptState int32
	// daemon is set before start and read immutably after start.
	daemon bool
	// done is closed when the thread terminates (state → stateTerminated).
	// join() reads from this channel to detect completion.
	done chan struct{}
	// waker is a buffered (cap 1) channel signaled by Interrupt() to wake a
	// thread blocked in sleep() or join().
	waker chan struct{}
	// isMain is true only for the primordial main thread. The interpreter uses
	// this to decide whether an uncaught exception should call os.Exit.
	isMain bool
}

// threadECSeq is a monotonically increasing counter for execution-context
// identity assignment. It starts at 1 so that 0 is reserved for "no owner".
// Protected by atomic for concurrent thread creation.
var threadECSeq uint64 = 1

// DefaultRunLoop is the interpreter loop function used by spawned threads.
// Set by the launcher before the main thread starts. native/thread.go calls
// this instead of importing interpreter directly, avoiding an import cycle.
var DefaultRunLoop func(*Thread)

func NewThread(loader Loader) *Thread {
	ecID := atomic.AddUint64(&threadECSeq, 1) - 1
	return &Thread{
		loader: loader,
		ecID:   ecID,
		done:   make(chan struct{}),
		waker:  make(chan struct{}, 1),
	}
}

func (t *Thread) Loader() Loader { return t.loader }
func (t *Thread) EC() uint64     { return t.ecID }

func (t *Thread) PushFrame(frame *Frame) {
	t.stack = append(t.stack, frame)
}

func (t *Thread) PopFrame() *Frame {
	n := len(t.stack)
	f := t.stack[n-1]
	t.stack[n-1] = nil // let the frame (and its slots' refs) be GC'd
	t.stack = t.stack[:n-1]
	return f
}

func (t *Thread) CurrentFrame() *Frame {
	return t.stack[len(t.stack)-1]
}

func (t *Thread) IsStackEmpty() bool {
	return len(t.stack) == 0
}

// FrameCount returns the current number of frames on the thread's stack.
func (t *Thread) FrameCount() int { return len(t.stack) }

// Bridge-mode accessors: used by the AOT bridge (interpreter.RunMethod) to capture
// a method's return when there is no caller frame to push it onto.
func (t *Thread) SetBridgeReturn(s *Slot) { t.bridgeReturn = s }
func (t *Thread) HasBridgeReturn() bool    { return t.bridgeReturn != nil }
func (t *Thread) BridgeReturn(s Slot)      { *t.bridgeReturn = s }

// --- Exception handling ---

func (t *Thread) Throw(obj *Object, pc int) { t.pendingException = obj; t.throwPC = pc }
func (t *Thread) HasException() bool        { return t.pendingException != nil }
func (t *Thread) ClearException() *Object {
	obj := t.pendingException
	t.pendingException = nil
	return obj
}
func (t *Thread) ThrowPC() int { return t.throwPC }

// NewFrame allocates a frame for a method on this thread.
func (t *Thread) NewFrame(method *Method) *Frame {
	return NewFrame(t, method)
}

// --- Lifecycle (ADR-0028) ---

// SetJavaThread attaches the canonical java.lang.Thread facade object.
func (t *Thread) SetJavaThread(obj *Object) { t.javaThread = obj }

// JavaThread returns the canonical java.lang.Thread facade object.
func (t *Thread) JavaThread() *Object { return t.javaThread }

// SetStarted atomically transitions state from NEW to RUNNABLE. Returns true on
// success; false means the thread was already started or terminated.
func (t *Thread) SetStarted() bool {
	return atomic.CompareAndSwapInt32(&t.state, stateNew, stateRunnable)
}

// IsAlive reports whether this thread has been started but not yet terminated.
func (t *Thread) IsAlive() bool {
	return atomic.LoadInt32(&t.state) == stateRunnable
}

// Terminate marks the thread as terminated and closes the done channel,
// unblocking any join() callers.
func (t *Thread) Terminate() {
	atomic.StoreInt32(&t.state, stateTerminated)
	close(t.done)
}

// --- Interrupt (ADR-0028) ---

// Interrupt sets the interrupt flag and signals the waker channel to unblock
// any sleep/join operation on this thread.
func (t *Thread) Interrupt() {
	atomic.StoreInt32(&t.interruptState, 1)
	select {
	case t.waker <- struct{}{}:
	default:
	}
}

// IsInterrupted returns the interrupt state without clearing it.
func (t *Thread) IsInterrupted() bool {
	return atomic.LoadInt32(&t.interruptState) == 1
}

// Interrupted atomically reads and clears the interrupt state. Returns the old
// value (whether the thread was interrupted). This is the static
// Thread.interrupted() semantic.
func (t *Thread) Interrupted() bool {
	return atomic.SwapInt32(&t.interruptState, 0) == 1
}

// --- Daemon ---

func (t *Thread) SetDaemon(v bool) { t.daemon = v }
func (t *Thread) IsDaemon() bool   { return t.daemon }

// --- Completion (for join) ---

// Done returns a channel that is closed when the thread terminates.
func (t *Thread) Done() <-chan struct{} { return t.done }

// Waker returns a channel signaled on Interrupt.
func (t *Thread) Waker() <-chan struct{} { return t.waker }

// --- Main thread ---

func (t *Thread) SetMain(v bool) { t.isMain = v }
func (t *Thread) IsMain() bool   { return t.isMain }

// --- Sleep ---

// Sleep blocks the calling goroutine for millis milliseconds, or until the
// thread is interrupted. If interrupted before or during sleep, it clears the
// interrupt flag and returns false (caller should throw InterruptedException).
// Returns true if sleep completed normally.
func (t *Thread) Sleep(millis int64) bool {
	if t.Interrupted() { // check and clear before sleeping
		return false
	}
	if millis <= 0 {
		return true
	}
	select {
	case <-time.After(time.Duration(millis) * time.Millisecond):
		return true
	case <-t.waker:
		atomic.StoreInt32(&t.interruptState, 0) // clear on interrupt
		return false
	}
}
