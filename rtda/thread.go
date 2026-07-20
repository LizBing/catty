package rtda

import (
	"sync"
	"sync/atomic"
	"time"
)

// Loader is the subset of the classloader that rtda needs at run time for class
// resolution (new, anewarray, checkcast, ldc class, invokeinterface, ...).
// Declaring it here as an interface keeps rtda free of any import cycle with the
// classloader package, which implements Loader concretely.
type Loader interface {
	// LoadClass is the must-load convenience method. It returns a fully linked
	// Class or panics. Only bootstrap invariants and legacy callers proven
	// unreachable from supported classfiles may use this method.
	LoadClass(name string) *Class

	// LoadClassResult returns a typed result: either a fully linked Class or a
	// terminal ClassLoadFailure. Java-reachable resolution paths MUST use this
	// method so failures propagate as Java throwables rather than Go panics.
	LoadClassResult(name string) ClassLoadResult

	// LoaderIdentity returns the opaque identity of this loader.
	// Primitive and void types use VMIdentity.
	LoaderIdentity() *LoaderIdentity
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
	// daemon flag (ADR-0028). Written by SetDaemon under configMu while
	// state==NEW; read immutably after start via ConsumeDaemonForStart.
	daemon   bool
	configMu sync.Mutex // serializes SetDaemon with the daemon read in start
	// done is closed when the thread terminates (state → stateTerminated).
	// join() reads from this channel to detect completion.
	done chan struct{}
	// waker is a buffered (cap 1) channel signaled by Interrupt() to wake a
	// thread blocked in sleep() or join().
	waker chan struct{}
	// isMain is true only for the primordial main thread. The interpreter uses
	// this to decide whether an uncaught exception should call os.Exit.
	isMain bool

	// --- Slice C: active-waiter protocol (ADR-0029) ---

	// waiterMu serializes the pre-wait interrupt check and active-waiter
	// publication. It closes the race between wait() observing the interrupt
	// flag and Interrupt() finding the waiter.
	waiterMu sync.Mutex
	// waitingOn is non-nil while this execution context is enrolled in a
	// monitor's wait set. Set under waiterMu before releasing the monitor;
	// cleared under waiterMu after reacquisition. Read by Interrupt under
	// waiterMu to locate and interrupt the waiter without holding the
	// monitor's state lock.
	waitingOn *Monitor
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
	// Release the implicit synchronized-method monitor (ADR-0029).
	// Only the implicit ACC_SYNCHRONIZED entry is attached to frame cleanup;
	// explicit block entries are governed by bytecode monitorexit.
	if f.syncObject != nil {
		f.syncObject.Monitor().Exit(t.ecID)
	}
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
// unblocking any join() callers. The CAS ensures exactly-once semantics:
// repeated or concurrent calls are harmless — only the first transition
// from RUNNABLE to TERMINATED closes done.
func (t *Thread) Terminate() {
	if !atomic.CompareAndSwapInt32(&t.state, stateRunnable, stateTerminated) {
		return // already terminated, or never started
	}
	close(t.done)
}

// --- Interrupt (ADR-0028) ---

// Interrupt sets the interrupt flag and signals the waker channel to unblock
// any sleep/join operation on this thread. If the target thread is waiting on
// a monitor (Object.wait), Interrupt also tries to atomically claim the waiter
// entry under the monitor's state lock (ADR-0029). The ordering between notify
// and interrupt is determined under that lock — one wins, the other does not.
func (t *Thread) Interrupt() {
	atomic.StoreInt32(&t.interruptState, 1)

	// Wake any sleep/join waiter.
	select {
	case t.waker <- struct{}{}:
	default:
	}

	// Wake any monitor waiter (Object.wait).
	t.waiterMu.Lock()
	m := t.waitingOn
	t.waiterMu.Unlock()
	if m != nil {
		m.InterruptWaiter(t.ecID)
	}
}

// IsInterrupted returns the interrupt state without clearing it.
func (t *Thread) IsInterrupted() bool {
	return atomic.LoadInt32(&t.interruptState) == 1
}

// Interrupted atomically reads and clears the interrupt state and drains any
// stale waker signal. Returns the old value (whether the thread was interrupted).
// This is the static Thread.interrupted() semantic.
//
// Draining the waker prevents a stale signal (left behind by a previous
// Interrupt() whose flag has now been cleared) from being consumed by a
// subsequent Sleep or Join as a spurious interrupt.
func (t *Thread) Interrupted() bool {
	wasInterrupted := atomic.SwapInt32(&t.interruptState, 0) == 1
	if wasInterrupted {
		// Drain the waker. A concurrent Interrupt() that fires after the
		// Swap above sets interruptState back to 1 before sending to waker,
		// so a real interrupt cannot be lost — the flag was re-set.
		select {
		case <-t.waker:
		default:
		}
	}
	return wasInterrupted
}

// --- Daemon ---

// SetDaemon sets the daemon flag. May only be called before the thread is
// started (state == NEW). Returns true on success; false means the thread has
// already been started or terminated, and the caller should throw
// IllegalThreadStateException.
//
// configMu serializes SetDaemon with ConsumeDaemonForStart so the daemon value
// read at start time is stable and the write is race-free.
func (t *Thread) SetDaemon(v bool) bool {
	t.configMu.Lock()
	defer t.configMu.Unlock()
	if atomic.LoadInt32(&t.state) != stateNew {
		return false
	}
	t.daemon = v
	return true
}

// IsDaemon reports whether this thread is a daemon thread. Holds configMu
// so concurrent SetDaemon on a not-yet-started thread is race-free.
func (t *Thread) IsDaemon() bool {
	t.configMu.Lock()
	d := t.daemon
	t.configMu.Unlock()
	return d
}

// ConsumeDaemonForStart reads the daemon flag under configMu, establishing a
// happens-before edge with any SetDaemon call that completed before start.
// Must be called once, immediately after SetStarted succeeds, to determine
// whether the thread counts toward VM liveness.
func (t *Thread) ConsumeDaemonForStart() bool {
	t.configMu.Lock()
	d := t.daemon
	t.configMu.Unlock()
	return d
}

// --- Completion (for join) ---

// Done returns a channel that is closed when the thread terminates.
func (t *Thread) Done() <-chan struct{} { return t.done }

// Waker returns a channel signaled on Interrupt.
func (t *Thread) Waker() <-chan struct{} { return t.waker }

// --- Main thread ---

func (t *Thread) SetMain(v bool) { t.isMain = v }
func (t *Thread) IsMain() bool   { return t.isMain }

// --- Monitor wait (Slice C, ADR-0029) ---

// MonitorWait implements the execution-context side of Object.wait().
//
// Phase 1 (under waiterMu): atomically check interrupt status AND publish the
// active waiter. If already interrupted, clear interrupt status and return
// false — the caller MUST throw InterruptedException WITHOUT releasing the
// monitor or altering recursion depth.
//
// Phase 2 (monitor): if the pre-check passed, delegate to Monitor.InternalWait
// which fully releases the monitor, blocks, and reacquires. The monitor
// restores the exact saved recursion depth before returning.
//
// Phase 3 (cleanup): clear waitingOn under waiterMu. If the waiter was
// interrupted (notify lost the race), the interrupt flag is cleared per JLS
// and the caller MUST throw InterruptedException — the monitor has already
// been reacquired and depth restored. If the waiter was notified normally,
// the caller returns normally; any pending interrupt flag remains set.
//
// Returns (normal, interrupted). The caller is responsible for throwing
// InterruptedException when interrupted is true.
func (t *Thread) MonitorWait(m *Monitor, savedDepth int) (normal bool, interrupted bool) {
	// Phase 1: pre-check + publish under waiterMu (closes the race).
	t.waiterMu.Lock()
	if atomic.LoadInt32(&t.interruptState) == 1 {
		// Pre-interrupted: clear, do NOT release monitor.
		t.Interrupted()
		t.waiterMu.Unlock()
		return false, true
	}
	t.waitingOn = m
	t.waiterMu.Unlock()

	// Phase 2: monitor handles release/enqueue/block/reacquire/depth restore.
	normal = m.InternalWait(t.ecID, savedDepth)

	// Phase 3: cleanup.
	t.waiterMu.Lock()
	t.waitingOn = nil
	t.waiterMu.Unlock()

	// If interrupt won the race against notify, clear the interrupt flag
	// per JLS so InterruptedException is thrown with status cleared.
	if !normal {
		t.Interrupted()
	}

	return normal, !normal
}

// --- Sleep ---

// Sleep blocks the calling goroutine for millis milliseconds, or until the
// thread is interrupted. If interrupted before sleep, returns false (caller
// should throw InterruptedException). Returns true if sleep completed normally.
//
// On waker signal, re-checks Interrupted() rather than unconditionally clearing
// the flag. This avoids treating a stale waker signal (drained by Interrupted()
// but re-delivered due to channel buffering) as a real interrupt.
func (t *Thread) Sleep(millis int64) bool {
	if t.Interrupted() { // check, clear, and drain before sleeping
		return false
	}
	if millis <= 0 {
		return true
	}
	select {
	case <-time.After(time.Duration(millis) * time.Millisecond):
		return true
	case <-t.waker:
		// Re-check: if the interrupt flag was cleared concurrently
		// (stale waker), return normally. If still interrupted, the
		// flag is now cleared by Interrupted() and we return false.
		if t.Interrupted() {
			return false
		}
		// Stale wake — interrupt was already consumed.
		return true
	}
}
