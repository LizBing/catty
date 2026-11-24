package rt

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Gate 5 application lifecycle. One process owns one lifecycle: a running
// state, a shutdown state with a frozen hook snapshot, and a terminal state
// with the selected Java or host status. Lifecycle state, hook registry
// mutations and live non-daemon publication are serialized by Runtime.thrMu.
// Java execution, hook startup, diagnostics and process finalization happen
// outside that lock.

type lifecycleState uint8

const (
	lifecycleRunning lifecycleState = iota
	lifecycleShutdown
	lifecycleTerminal
)

// terminalResult is the one committed process result. reason is "natural",
// "exit" or "halt"; requester is the logical thread that requested the
// transition, or zero for a natural shutdown.
type terminalResult struct {
	reason    string
	status    int32
	requester int64
}

// lifecyclePlan is the frozen shutdown snapshot and the executions it owns.
type lifecyclePlan struct {
	reason    string
	status    int32
	requester int64
	hooks     []*Object
	starts    []*hookStart
	waits     []<-chan struct{}
}

// hookStart is one snapshot hook that the VM must start.
type hookStart struct {
	peer   *Object
	worker *Thread
	name   string
	daemon bool
}

// processEpoch is the arbitrary process-local monotonic origin for
// System.nanoTime.
var processEpoch = time.Now()

// lifecycleEnabled reports whether the installed profile publishes the
// application lifecycle.
func (r *Runtime) lifecycleEnabled() bool {
	return r.profile != nil && r.profile.Lifecycle
}

// threadDisplayName returns the Java thread name for diagnostics.
func (r *Runtime) threadDisplayName(t *Thread) string {
	if t != nil {
		if p, ok := peerOfObject(t.peer); ok {
			return p.name
		}
	}
	return "main"
}

// ------------------------------------------------------------ transitions

// requestExit is the shared System.exit and Runtime.exit request. The first
// request that observes the running state freezes the registry and starts
// graceful shutdown; every other request parks until process termination. A
// successful exit call never returns into Java, so its frames and monitors
// stay live until the process exits.
func (r *Runtime) requestExit(t *Thread, status int32) {
	if hook := r.ExitRequestHook; hook != nil {
		hook(status)
	}
	requester := int64(0)
	if t != nil {
		requester = t.ID
	}
	r.thrMu.Lock()
	if r.lifeState != lifecycleRunning {
		r.thrMu.Unlock()
		parkJavaThread()
	}
	r.lifeState = lifecycleShutdown
	r.lifeReason = "exit"
	r.lifeStatus = status
	r.lifeRequester = requester
	plan := r.freezeHooksLocked("exit", status, requester)
	r.thrMu.Unlock()
	go r.activateShutdown(plan)
	parkJavaThread()
}

// requestHalt commits a terminal halt result immediately. A halt in the
// running state skips the shutdown event and hook snapshot; a halt during
// graceful shutdown cuts it short by publishing the halt terminal result
// before the remaining hooks complete. The requesting thread parks without
// unwinding.
func (r *Runtime) requestHalt(t *Thread, status int32) {
	requester := int64(0)
	if t != nil {
		requester = t.ID
	}
	if r.publishTerminal("halt", status, requester) {
		// The committed result signals the launcher through requestTerminal.
	}
	parkJavaThread()
}

// publishTerminal commits the one terminal result. Later requests cannot
// change it; callers that lose wait for process termination.
func (r *Runtime) publishTerminal(reason string, status int32, requester int64) bool {
	r.thrMu.Lock()
	if r.lifeState == lifecycleTerminal {
		r.thrMu.Unlock()
		return false
	}
	r.lifeState = lifecycleTerminal
	r.lifeReason = reason
	r.lifeStatus = status
	r.lifeRequester = requester
	// Release registry and snapshot references: process termination follows.
	r.hooks = nil
	r.hookSet = nil
	r.thrMu.Unlock()
	// A halt may win the race after the registry froze but before the
	// coordinator recorded shutdown_begin. Emit the frozen begin record first,
	// so process_exit can never precede it. A halt in the running state has no
	// pending begin and still skips the event.
	r.ensureShutdownBegin()
	select {
	case r.lifecycleTerminal <- terminalResult{reason: reason, status: status, requester: requester}:
	default:
		// The channel is buffered and only the first terminal publication
		// sends; a full channel means a committed result already exists.
	}
	return true
}

// beginNaturalShutdownLocked performs the first-zero transition. The caller
// holds thrMu and has already updated mainRetired/liveNonDaemon; state,
// status and the hook snapshot therefore commit atomically with the last
// participating thread's termination publication.
func (r *Runtime) beginNaturalShutdownLocked() *lifecyclePlan {
	if !r.lifecycleEnabled() || r.lifeState != lifecycleRunning || !r.mainRetired || r.liveNonDaemon != 0 {
		return nil
	}
	status := r.mainStatus
	r.lifeState = lifecycleShutdown
	r.lifeReason = "natural"
	r.lifeStatus = status
	return r.freezeHooksLocked("natural", status, 0)
}

// freezeHooksLocked detaches the registry once and prepares the snapshot's
// start and completion bookkeeping. Terminated peers are skipped (they have
// already run); a peer already started by Java is only awaited. The caller
// holds thrMu.
func (r *Runtime) freezeHooksLocked(reason string, status int32, requester int64) *lifecyclePlan {
	plan := &lifecyclePlan{
		reason:    reason,
		status:    status,
		requester: requester,
		hooks:     append([]*Object(nil), r.hooks...),
	}
	r.hooks = nil
	r.hookSet = nil
	r.beginMu.Lock()
	r.shutdownBeginPending = true
	r.shutdownBeginWritten = false
	r.frozenReason = reason
	r.frozenStatus = status
	r.frozenRequester = requester
	r.frozenHookCount = len(plan.hooks)
	r.beginMu.Unlock()
	for _, hook := range plan.hooks {
		p, ok := peerOfObject(hook)
		if !ok || p == nil {
			continue
		}
		switch p.state {
		case threadNew:
			if p.hookDone == nil {
				p.hookDone = make(chan struct{})
			}
			plan.waits = append(plan.waits, p.hookDone)
			plan.starts = append(plan.starts, &hookStart{peer: hook, name: p.name, daemon: p.daemon})
		case threadStarted:
			// A manual start won before the snapshot. Count that execution only
			// while it is still live, and never start the peer again.
			if p.hookDone == nil {
				p.hookDone = make(chan struct{})
			}
			plan.waits = append(plan.waits, p.hookDone)
		case threadTerminated:
			// Registration of a retired peer is allowed; it is skipped because
			// it has already run.
		}
	}
	return plan
}

// activateShutdown emits the one shutdown event and starts every eligible hook
// before waiting for any. Hook frames, Java calls and diagnostics all run
// outside the lifecycle lock.
func (r *Runtime) activateShutdown(plan *lifecyclePlan) {
	if plan == nil {
		return
	}
	r.keepShutdownAlive()
	if hook := r.ActivateShutdownHook; hook != nil {
		hook()
	}
	// Record the shutdown event before checking for a terminal race. A halt
	// that already won still owes the trace the begin record for the frozen
	// snapshot; terminal publication emits it if this coordinator is late.
	r.ensureShutdownBegin()
	r.thrMu.Lock()
	terminal := r.lifeState == lifecycleTerminal
	r.thrMu.Unlock()
	if terminal {
		return
	}
	if hook := r.BeforeHookStartHook; hook != nil {
		hook()
	}

	// Publish the new-to-started transition for every eligible snapshot member
	// under the lifecycle lock, then start the goroutines outside it.
	r.thrMu.Lock()
	var started []*hookStart
	for _, start := range plan.starts {
		p, ok := peerOfObject(start.peer)
		if !ok || p == nil {
			continue
		}
		if p.state != threadNew {
			// A concurrent manual start won; its hookDone is already waited on.
			continue
		}
		// Capture the daemon flag at the one new-to-started transition, after
		// any successful setDaemon that raced the frozen snapshot. The trace
		// and the live non-daemon accounting must use the same value.
		daemon := p.daemon
		p.state = threadStarted
		if !daemon {
			r.liveNonDaemon++
		}
		worker := r.NewThread()
		p.thread = worker
		worker.peer = start.peer
		start.worker = worker
		start.daemon = daemon
		r.noteLiveThreadStartLocked()
		started = append(started, start)
	}
	r.thrMu.Unlock()

	for _, start := range started {
		if tr := r.trace; tr != nil {
			tr.threadStartShutdownHook(start.worker.ID, start.daemon, start.name)
		}
		go r.runPeer(peerOfObjectUnsafe(start.peer), start.peer, start.worker)
	}

	go func() {
		for _, done := range plan.waits {
			<-done
		}
		r.publishTerminal(plan.reason, plan.status, plan.requester)
	}()
}

// ensureShutdownBegin records the one frozen-registry shutdown event. It is
// safe to call from both the shutdown coordinator and terminal publication:
// beginMu serializes the check and the write, so exactly one event precedes
// the final process_exit record.
func (r *Runtime) ensureShutdownBegin() {
	r.beginMu.Lock()
	defer r.beginMu.Unlock()
	if !r.shutdownBeginPending || r.shutdownBeginWritten {
		return
	}
	r.shutdownBeginWritten = true
	if tr := r.trace; tr != nil {
		tr.shutdownBegin(r.frozenReason, r.frozenStatus, r.frozenHookCount, r.frozenRequester)
	}
}

// keepShutdownAlive guarantees one runnable-later wakeup while graceful
// shutdown waits on a hook. Without it, an untimed Object.wait or a parked
// caller with no timer can leave every Go goroutine asleep and make the Go
// runtime abort the process instead of matching the JVM's live shutdown.
func (r *Runtime) keepShutdownAlive() {
	r.shutdownKeeperOnce.Do(func() {
		go func() {
			for {
				time.Sleep(time.Hour)
			}
		}()
	})
}

func peerOfObjectUnsafe(o *Object) *threadPeer {
	p, _ := peerOfObject(o)
	return p
}

// parkJavaThread parks the requesting logical thread without unwinding its
// frames or releasing its Java monitors. The sleep loop keeps the Go runtime
// from declaring a deadlock when a deliberately stuck Java shutdown has no
// other runnable work; an external watchdog or a coordinated halt ends the
// process.
func parkJavaThread() {
	for {
		time.Sleep(time.Hour)
	}
}

// ------------------------------------------------------------ Java surface

func intrinsicSystemExit(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	t.RT.requestExit(t, AsInt(args[0]))
	return Normal()
}

func intrinsicSystemNanoTime(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return OkValue(LongValue(time.Since(processEpoch).Nanoseconds()))
}

func intrinsicRuntimeGetRuntime(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if t.RT.runtimeObject == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalStateException"))
	}
	return OkValue(RefValue(t.RT.runtimeObject))
}

func intrinsicRuntimeExit(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	t.RT.requestExit(t, AsInt(args[0]))
	return Normal()
}

func intrinsicRuntimeHalt(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	t.RT.requestHalt(t, AsInt(args[0]))
	return Normal()
}

// intrinsicRuntimeAddShutdownHook implements Runtime.addShutdownHook. The
// registry is receiver-independent and keyed by Java object identity. Once
// shutdown claims the registry both mutation operations throw
// IllegalStateException before validating the hook argument.
func intrinsicRuntimeAddShutdownHook(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 || args[0].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	hook := args[0].Ref
	r := t.RT
	r.thrMu.Lock()
	if r.lifeState != lifecycleRunning {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalStateException"))
	}
	if hook == nil {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if r.threadClass == nil || !r.IsAssignable(hook.Class, r.threadClass) {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/ClassCastException"))
	}
	p, ok := peerOfObject(hook)
	if !ok || p == nil {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if _, exists := r.hookSet[hook]; exists {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if p.state == threadStarted {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	if r.hookSet == nil {
		r.hookSet = map[*Object]struct{}{}
	}
	r.hookSet[hook] = struct{}{}
	r.hooks = append(r.hooks, hook)
	r.thrMu.Unlock()
	return Normal()
}

// intrinsicRuntimeRemoveShutdownHook implements Runtime.removeShutdownHook.
func intrinsicRuntimeRemoveShutdownHook(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if len(args) != 1 || args[0].Kind != KindRef {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	hook := args[0].Ref
	r := t.RT
	r.thrMu.Lock()
	if r.lifeState != lifecycleRunning {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalStateException"))
	}
	if hook == nil {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/NullPointerException"))
	}
	removed := false
	if _, exists := r.hookSet[hook]; exists {
		delete(r.hookSet, hook)
		for i, candidate := range r.hooks {
			if candidate == hook {
				copy(r.hooks[i:], r.hooks[i+1:])
				r.hooks[len(r.hooks)-1] = nil
				r.hooks = r.hooks[:len(r.hooks)-1]
				break
			}
		}
		removed = true
	}
	r.thrMu.Unlock()
	if removed {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicRuntimeGC(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	t.RT.collectGarbage()
	return Normal()
}

// reportHostFault captures the logical thread's Java frame snapshot while its
// frames are still live and writes it together with the original panic value
// and Go stack. It never converts a host fault into a Java throwable and never
// participates in ordinary shutdown completion.
func (r *Runtime) reportHostFault(t *Thread, value any) {
	var b strings.Builder
	if t != nil {
		fmt.Fprintf(&b, "Catty host fault on Java thread %q\n", r.threadDisplayName(t))
		for _, frame := range t.FramesSnapshot() {
			m := frame.Method
			if m == nil || m.Owner == nil {
				continue
			}
			owner := strings.ReplaceAll(m.Owner.Name, "/", ".")
			fmt.Fprintf(&b, "\tat %s.%s%s (pc %d)\n", owner, m.Name, m.Descriptor, frame.PC)
		}
	} else {
		b.WriteString("Catty host fault\n")
	}
	fmt.Fprintf(&b, "panic: %v\n", value)
	b.Write(debug.Stack())
	r.diagMu.Lock()
	_, _ = os.Stderr.WriteString(b.String())
	r.diagMu.Unlock()
}
