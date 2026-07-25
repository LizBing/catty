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

// ExecutionContext models one JVM execution context's stack of frames
// (JVMS §2.5.1) plus the Java Thread state that has not yet moved behind the
// facade boundary. Slice B makes ExecutionContext the canonical runtime type;
// Thread remains a temporary compatibility alias until native Thread facade
// state is split out.
type ExecutionContext struct {
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
	// bridgeDynamicReturn captures one logical Java result for the typed direct
	// invocation adapters. It is separate from bridgeReturn so the legacy AOT
	// Slot bridge remains an adapter rather than the stable result boundary.
	bridgeDynamicReturn *JavaValue
	// pendingException is non-nil when an exception is in flight (athrow or a
	// runtime error like NPE). The interpreter Loop checks HasException after
	// each instruction and dispatches to handleException.
	pendingException *Object
	throwPC          int // PC of the instruction that threw (for exception-table search)

	threadState *JavaThreadState
}

// JavaThreadState is the runtime sidecar for the bounded java.lang.Thread
// facade. It owns Java-visible Thread lifecycle, interrupt, daemon, completion,
// sleep/join wakeup, main-thread marker, and monitor-wait enrollment state.
// The owning ExecutionContext remains the execution stack, loader, throwable,
// bridge-return, class-init, and monitor-owner identity.
type JavaThreadState struct {
	context *ExecutionContext

	// facade is the canonical java.lang.Thread object attached to this runtime
	// record. currentThread() returns this object.
	facade *Object
	// state is the lifecycle state (stateNew / stateRunnable / stateTerminated).
	// Managed with atomic ops: SetStarted uses CAS; Terminate and IsAlive use
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
	// done is closed when the thread terminates (state -> stateTerminated).
	// join() reads from this channel to detect completion.
	done chan struct{}
	// waker is a buffered (cap 1) channel signaled by Interrupt() to wake a
	// thread blocked in sleep() or join().
	waker chan struct{}
	// isMain is true only for the primordial main thread. The interpreter uses
	// this to decide whether an uncaught exception should call os.Exit.
	isMain bool

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

// Thread is a compatibility alias retained while Java Thread facade callers are
// migrated. New code should name ExecutionContext when it means execution
// state, and reserve Java Thread terminology for the java.lang.Thread facade.
type Thread = ExecutionContext

// threadECSeq is a monotonically increasing counter for execution-context
// identity assignment. It starts at 1 so that 0 is reserved for "no owner".
// Protected by atomic for concurrent thread creation.
var threadECSeq uint64 = 1

// DefaultRunLoop is the interpreter loop function used by spawned threads.
// Set by the launcher before the main thread starts. native/thread.go calls
// this instead of importing interpreter directly, avoiding an import cycle.
var DefaultRunLoop func(*ExecutionContext)

func NewExecutionContext(loader Loader) *ExecutionContext {
	ecID := atomic.AddUint64(&threadECSeq, 1) - 1
	t := &ExecutionContext{
		loader: loader,
		ecID:   ecID,
	}
	t.threadState = NewJavaThreadState(t)
	return t
}

func NewThread(loader Loader) *Thread {
	return NewExecutionContext(loader)
}

func (t *ExecutionContext) Loader() Loader { return t.loader }
func (t *ExecutionContext) ID() uint64     { return t.ecID }
func (t *ExecutionContext) EC() uint64     { return t.ID() }

func NewJavaThreadState(context *ExecutionContext) *JavaThreadState {
	return &JavaThreadState{
		context: context,
		done:    make(chan struct{}),
		waker:   make(chan struct{}, 1),
	}
}

func (t *ExecutionContext) JavaThreadState() *JavaThreadState { return t.threadState }

func (s *JavaThreadState) Context() *ExecutionContext { return s.context }

func (t *ExecutionContext) PushFrame(frame *Frame) {
	t.stack = append(t.stack, frame)
}

func (t *ExecutionContext) PopFrame() *Frame {
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

func (t *ExecutionContext) CurrentFrame() *Frame {
	return t.stack[len(t.stack)-1]
}

func (t *ExecutionContext) IsStackEmpty() bool {
	return len(t.stack) == 0
}

// FrameCount returns the current number of frames on the thread's stack.
func (t *ExecutionContext) FrameCount() int { return len(t.stack) }

// Bridge-mode accessors: used by the AOT bridge (interpreter.RunMethod) to capture
// a method's return when there is no caller frame to push it onto.
func (t *ExecutionContext) SetBridgeReturn(s *Slot) { t.bridgeReturn = s }
func (t *ExecutionContext) HasBridgeReturn() bool   { return t.bridgeReturn != nil }
func (t *ExecutionContext) BridgeReturn(s Slot)     { *t.bridgeReturn = s }

func (t *ExecutionContext) SetBridgeDynamicReturn(v *JavaValue) { t.bridgeDynamicReturn = v }
func (t *ExecutionContext) HasBridgeDynamicReturn() bool        { return t.bridgeDynamicReturn != nil }
func (t *ExecutionContext) BridgeDynamicReturn(v JavaValue)     { *t.bridgeDynamicReturn = v }

// --- Exception handling ---

func (t *ExecutionContext) Throw(obj *Object, pc int) { t.pendingException = obj; t.throwPC = pc }
func (t *ExecutionContext) HasException() bool        { return t.pendingException != nil }
func (t *ExecutionContext) ClearException() *Object {
	obj := t.pendingException
	t.pendingException = nil
	return obj
}
func (t *ExecutionContext) ThrowPC() int { return t.throwPC }

// NewFrame allocates a frame for a method on this thread.
func (t *ExecutionContext) NewFrame(method *Method) *Frame {
	return NewFrame(t, method)
}

// --- Lifecycle (ADR-0028) ---

// SetJavaThread attaches the canonical java.lang.Thread facade object.
func (t *ExecutionContext) SetJavaThread(obj *Object) { t.threadState.SetJavaThread(obj) }

// JavaThread returns the canonical java.lang.Thread facade object.
func (t *ExecutionContext) JavaThread() *Object { return t.threadState.JavaThread() }

func (s *JavaThreadState) SetJavaThread(obj *Object) { s.facade = obj }

func (s *JavaThreadState) JavaThread() *Object { return s.facade }

// SetStarted atomically transitions state from NEW to RUNNABLE. Returns true on
// success; false means the thread was already started or terminated.
func (t *ExecutionContext) SetStarted() bool { return t.threadState.SetStarted() }

func (s *JavaThreadState) SetStarted() bool {
	return atomic.CompareAndSwapInt32(&s.state, stateNew, stateRunnable)
}

// IsAlive reports whether this thread has been started but not yet terminated.
func (t *ExecutionContext) IsAlive() bool { return t.threadState.IsAlive() }

func (s *JavaThreadState) IsAlive() bool { return atomic.LoadInt32(&s.state) == stateRunnable }

// Terminate marks the thread as terminated and closes the done channel,
// unblocking any join() callers. The CAS ensures exactly-once semantics:
// repeated or concurrent calls are harmless — only the first transition
// from RUNNABLE to TERMINATED closes done.
func (t *ExecutionContext) Terminate() { t.threadState.Terminate() }

func (s *JavaThreadState) Terminate() {
	if !atomic.CompareAndSwapInt32(&s.state, stateRunnable, stateTerminated) {
		return // already terminated, or never started
	}
	close(s.done)
}

// --- Interrupt (ADR-0028) ---

// Interrupt sets the interrupt flag and signals the waker channel to unblock
// any sleep/join operation on this thread. If the target thread is waiting on
// a monitor (Object.wait), Interrupt also tries to atomically claim the waiter
// entry under the monitor's state lock (ADR-0029). The ordering between notify
// and interrupt is determined under that lock — one wins, the other does not.
func (t *ExecutionContext) Interrupt() { t.threadState.Interrupt() }

func (s *JavaThreadState) Interrupt() {
	atomic.StoreInt32(&s.interruptState, 1)

	// Wake any sleep/join waiter.
	select {
	case s.waker <- struct{}{}:
	default:
	}

	// Wake any monitor waiter (Object.wait).
	s.waiterMu.Lock()
	m := s.waitingOn
	s.waiterMu.Unlock()
	if m != nil {
		m.InterruptWaiter(s.context.ID())
	}
}

// IsInterrupted returns the interrupt state without clearing it.
func (t *ExecutionContext) IsInterrupted() bool { return t.threadState.IsInterrupted() }

func (s *JavaThreadState) IsInterrupted() bool { return atomic.LoadInt32(&s.interruptState) == 1 }

// Interrupted atomically reads and clears the interrupt state and drains any
// stale waker signal. Returns the old value (whether the thread was interrupted).
// This is the static Thread.interrupted() semantic.
//
// Draining the waker prevents a stale signal (left behind by a previous
// Interrupt() whose flag has now been cleared) from being consumed by a
// subsequent Sleep or Join as a spurious interrupt.
func (t *ExecutionContext) Interrupted() bool { return t.threadState.Interrupted() }

func (s *JavaThreadState) Interrupted() bool {
	wasInterrupted := atomic.SwapInt32(&s.interruptState, 0) == 1
	if wasInterrupted {
		// Drain the waker. A concurrent Interrupt() that fires after the
		// Swap above sets interruptState back to 1 before sending to waker,
		// so a real interrupt cannot be lost — the flag was re-set.
		select {
		case <-s.waker:
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
func (t *ExecutionContext) SetDaemon(v bool) bool { return t.threadState.SetDaemon(v) }

func (s *JavaThreadState) SetDaemon(v bool) bool {
	s.configMu.Lock()
	defer s.configMu.Unlock()
	if atomic.LoadInt32(&s.state) != stateNew {
		return false
	}
	s.daemon = v
	return true
}

// IsDaemon reports whether this thread is a daemon thread. Holds configMu
// so concurrent SetDaemon on a not-yet-started thread is race-free.
func (t *ExecutionContext) IsDaemon() bool { return t.threadState.IsDaemon() }

func (s *JavaThreadState) IsDaemon() bool {
	s.configMu.Lock()
	d := s.daemon
	s.configMu.Unlock()
	return d
}

// ConsumeDaemonForStart reads the daemon flag under configMu, establishing a
// happens-before edge with any SetDaemon call that completed before start.
// Must be called once, immediately after SetStarted succeeds, to determine
// whether the thread counts toward VM liveness.
func (t *ExecutionContext) ConsumeDaemonForStart() bool { return t.threadState.ConsumeDaemonForStart() }

func (s *JavaThreadState) ConsumeDaemonForStart() bool {
	s.configMu.Lock()
	d := s.daemon
	s.configMu.Unlock()
	return d
}

// --- Completion (for join) ---

// Done returns a channel that is closed when the thread terminates.
func (t *ExecutionContext) Done() <-chan struct{} { return t.threadState.Done() }

func (s *JavaThreadState) Done() <-chan struct{} { return s.done }

// Waker returns a channel signaled on Interrupt.
func (t *ExecutionContext) Waker() <-chan struct{} { return t.threadState.Waker() }

func (s *JavaThreadState) Waker() <-chan struct{} { return s.waker }

// --- Main thread ---

func (t *ExecutionContext) SetMain(v bool) { t.threadState.SetMain(v) }
func (t *ExecutionContext) IsMain() bool   { return t.threadState.IsMain() }

func (s *JavaThreadState) SetMain(v bool) { s.isMain = v }
func (s *JavaThreadState) IsMain() bool   { return s.isMain }

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
func (t *ExecutionContext) MonitorWait(m *Monitor, savedDepth int) (normal bool, interrupted bool) {
	return t.threadState.MonitorWait(m, savedDepth)
}

func (s *JavaThreadState) MonitorWait(m *Monitor, savedDepth int) (normal bool, interrupted bool) {
	// Phase 1: pre-check + publish under waiterMu (closes the race).
	s.waiterMu.Lock()
	if atomic.LoadInt32(&s.interruptState) == 1 {
		// Pre-interrupted: clear, do NOT release monitor.
		s.Interrupted()
		s.waiterMu.Unlock()
		return false, true
	}
	s.waitingOn = m
	s.waiterMu.Unlock()

	// Phase 2: monitor handles release/enqueue/block/reacquire/depth restore.
	normal = m.InternalWait(s.context.ID(), savedDepth)

	// Phase 3: cleanup.
	s.waiterMu.Lock()
	s.waitingOn = nil
	s.waiterMu.Unlock()

	// If interrupt won the race against notify, clear the interrupt flag
	// per JLS so InterruptedException is thrown with status cleared.
	if !normal {
		s.Interrupted()
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
func (t *ExecutionContext) Sleep(millis int64) bool { return t.threadState.Sleep(millis) }

func (s *JavaThreadState) Sleep(millis int64) bool {
	if s.Interrupted() { // check, clear, and drain before sleeping
		return false
	}
	if millis <= 0 {
		return true
	}
	select {
	case <-time.After(time.Duration(millis) * time.Millisecond):
		return true
	case <-s.waker:
		// Re-check: if the interrupt flag was cleared concurrently
		// (stale waker), return normally. If still interrupted, the
		// flag is now cleared by Interrupted() and we return false.
		if s.Interrupted() {
			return false
		}
		// Stale wake — interrupt was already consumed.
		return true
	}
}
