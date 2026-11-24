package rt

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// threadState is the lifecycle state of a Java Thread peer. It is guarded by
// Runtime.thrMu.
type threadState uint8

const (
	threadNew threadState = iota
	threadStarted
	threadTerminated
)

// threadPeer is the runtime-owned state of one java.lang.Thread object. The
// peer object's Payload points at it, so a terminated thread's name, daemon
// status, identity and pending interrupt remain observable after termination.
type threadPeer struct {
	id       int64
	name     string
	daemon   bool
	state    threadState
	runnable *Object
	thread   *Thread
	// pendingInterrupt is the Java interrupt status. It survives termination
	// and can be set on a peer that has not started yet.
	pendingInterrupt bool
	// hookDone closes exactly once when this peer terminates while it is a
	// member of the frozen shutdown snapshot. It is nil for ordinary peers and
	// for snapshot members that had already terminated at freeze time.
	hookDone chan struct{}
}

func peerData(o *Object) *threadPeer {
	if o == nil {
		return nil
	}
	p, _ := o.Payload.(*threadPeer)
	return p
}

// peerOfObject returns the peer state when o is a java.lang.Thread object.
func peerOfObject(o *Object) (*threadPeer, bool) {
	if o == nil {
		return nil, false
	}
	p, ok := o.Payload.(*threadPeer)
	return p, ok
}

// waiterState records the single winning transition of one wait record.
type waiterState uint8

const (
	waiterPending waiterState = iota
	waiterNotified
	waiterTimedOut
	waiterInterrupted
)

// waiter is one wait-set record. Every field is guarded by Runtime.thrMu; the
// channel is closed exactly once by the winning transition.
type waiter struct {
	thread     *Thread
	monitor    *monitor
	savedDepth int
	state      waiterState
	ch         chan struct{}
	// object keeps the monitored object reachable while the wait record is
	// pending, so an unreachable-but-waited object cannot be collected out
	// from under its wait set.
	object *Object
	// timer and timerDone are the current timeout segment and its cancellation
	// signal. Setting a timer replaces any previous segment; the winning
	// transition stops the timer and closes timerDone so the timer goroutine
	// exits instead of retaining the waiter until the Java deadline.
	timer     runtimeTimer
	timerDone chan struct{}
}

// waitTimer is one absolute Java wait deadline. start is sampled from the
// runtime's monotonic clock and millis is the rounded Java timeout; remaining
// returns a bounded timer segment and whether the deadline has passed. A
// duration that exceeds Go's representable range is handled in segments, so
// the deadline keeps its full Java length while each host timer stays bounded.
type waitTimer struct {
	start  time.Time
	millis int64
}

// wallElapsedMillis returns the non-negative elapsed whole milliseconds between
// start and now without using time.Time.Sub, which saturates beyond Go's
// duration range. It is used only after the monotonic reading saturates.
func wallElapsedMillis(start, now time.Time) int64 {
	seconds := now.Unix() - start.Unix()
	nanos := int64(now.Nanosecond()) - int64(start.Nanosecond())
	if nanos < 0 {
		seconds--
		nanos += int64(time.Second)
	}
	if seconds < 0 {
		return 0
	}
	maxMillis := int64(math.MaxInt64)
	if seconds > (maxMillis-nanos/int64(time.Millisecond))/1000 {
		return maxMillis
	}
	return seconds*1000 + nanos/int64(time.Millisecond)
}

func (r *Runtime) newWaitTimer(millis int64) waitTimer {
	return waitTimer{start: r.clock.Now(), millis: millis}
}

func (d waitTimer) remaining(now time.Time) (time.Duration, bool) {
	elapsed := now.Sub(d.start)
	var elapsedMillis int64
	if elapsed < 0 {
		elapsedMillis = 0
	} else if elapsed < maxTimerDuration {
		// The monotonic reading covers every duration Go can represent.
		elapsedMillis = int64(elapsed / time.Millisecond)
	} else {
		// Sub saturates at MaxInt64 nanoseconds, which is less than the largest
		// Java timeouts. Fall back to whole seconds and nanoseconds for the
		// overflow range so the exact Java deadline is preserved.
		elapsedMillis = wallElapsedMillis(d.start, now)
	}
	remainingMillis := d.millis - elapsedMillis
	if remainingMillis <= 0 {
		return 0, true
	}
	maxMillis := int64(maxTimerDuration / time.Millisecond)
	if remainingMillis > maxMillis {
		return maxTimerDuration, false
	}
	return time.Duration(remainingMillis) * time.Millisecond, false
}

// timerSegment returns the next bounded host timer for a remaining Java wait,
// capped by the runtime's segment limit.
func (r *Runtime) timerSegment(remaining time.Duration) time.Duration {
	maxMillis := int64(maxTimerDuration / time.Millisecond)
	limit := r.maxTimerSegmentMillis
	if limit <= 0 || limit > maxMillis {
		limit = maxMillis
	}
	maxSegment := time.Duration(limit) * time.Millisecond
	if maxSegment <= 0 {
		maxSegment = time.Millisecond
	}
	segment := remaining
	if segment > maxSegment {
		segment = maxSegment
	}
	if segment < time.Millisecond {
		segment = time.Millisecond
	}
	return segment
}

// monitor is the optional monitor state of one Java object. owner and depth
// describe ownership, cond wakes monitor entry contenders, and waiters is the
// object's wait set. Every field is guarded by Runtime.thrMu.
type monitor struct {
	owner      *Thread
	depth      int
	cond       *sync.Cond
	waiters    []*waiter
	contenders int
}

// blockerKind identifies the interruptible operation one thread is parked in.
type blockerKind uint8

const (
	blockerWait blockerKind = iota
	blockerJoin
	blockerSleep
)

// threadBlocker is the active interruptible blocker of one logical thread. The
// wake channel is closed exactly once by whichever transition wins.
type threadBlocker struct {
	kind   blockerKind
	waiter *waiter
	ch     chan struct{}
	sleep  waiterState // pending, timedOut or interrupted for a sleep blocker
}

// ---------------------------------------------------------------- monitors

// monitorForLocked returns o's monitor, creating it on first use. The caller
// holds thrMu. Monitor state lives on the object so its lifetime is the
// object's; the runtime never keeps a map that roots every monitored object.
func (r *Runtime) monitorForLocked(o *Object) *monitor {
	if o.monitor == nil {
		o.monitor = &monitor{}
	}
	return o.monitor
}

func (r *Runtime) monitorFor(o *Object) *monitor {
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	return r.monitorForLocked(o)
}

// MonitorEnter implements the monitorenter instruction. A null reference
// throws NullPointerException; otherwise the current logical thread owns the
// object's monitor on return. Monitor entry preserves the interrupt status.
func (t *Thread) MonitorEnter(m *Method, pc int32, ref Value) *Object {
	if ref.Kind != KindRef || ref.Ref == nil {
		return t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return t.monitorEnterObject(ref.Ref, m, pc)
}

func (t *Thread) monitorEnterObject(o *Object, m *Method, pc int32) *Object {
	return t.monitorEnter(o, m, pc, true)
}

// monitorEnterSilent acquires a monitor for a runtime-internal transition that
// has no Java-visible monitor operation. Termination uses it to publish state
// and notify joiners while owning the target's peer monitor without emitting
// trace events after the thread_end record.
func (t *Thread) monitorEnterSilent(o *Object) *Object {
	return t.monitorEnter(o, nil, -1, false)
}

func (t *Thread) monitorEnter(o *Object, m *Method, pc int32, traced bool) *Object {
	if o == nil {
		return t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	r := t.RT
	r.thrMu.Lock()
	mo := r.monitorForLocked(o)
	contended := false
	for mo.owner != nil && mo.owner != t {
		contended = true
		if mo.cond == nil {
			mo.cond = sync.NewCond(&r.thrMu)
		}
		mo.contenders++
		mo.cond.Wait()
		mo.contenders--
	}
	mo.owner = t
	mo.depth++
	if t.held == nil {
		t.held = map[*Object]*monitor{}
	}
	t.held[o] = mo
	depth := mo.depth
	r.thrMu.Unlock()
	if contended && r.metrics != nil {
		r.metrics.monitorContention()
	}
	if traced {
		if tr := t.trace(); tr != nil {
			if contended {
				tr.monitor(t.ID, o.Identity, "contended", depth, m, pc, "")
			}
			tr.monitor(t.ID, o.Identity, "enter", depth, m, pc, "")
		}
	}
	return nil
}

// MonitorExit implements the monitorexit instruction. Only the owner may exit;
// an exit by any other thread throws IllegalMonitorStateException and leaves
// the monitor unchanged.
func (t *Thread) MonitorExit(m *Method, pc int32, ref Value) *Object {
	if ref.Kind != KindRef || ref.Ref == nil {
		return t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	return t.monitorExitObject(ref.Ref, m, pc)
}

func (t *Thread) monitorExitObject(o *Object, m *Method, pc int32) *Object {
	return t.monitorExit(o, m, pc, true)
}

func (t *Thread) monitorExitSilent(o *Object) *Object {
	return t.monitorExit(o, nil, -1, false)
}

func (t *Thread) monitorExit(o *Object, m *Method, pc int32, traced bool) *Object {
	if o == nil {
		return t.RT.NewThrowable(t, "java/lang/NullPointerException")
	}
	r := t.RT
	r.thrMu.Lock()
	if o.monitor == nil || o.monitor.owner != t || o.monitor.depth <= 0 {
		r.thrMu.Unlock()
		return r.NewThrowable(t, "java/lang/IllegalMonitorStateException")
	}
	mo := o.monitor
	mo.depth--
	depth := mo.depth
	if depth == 0 {
		mo.owner = nil
		delete(t.held, o)
		if mo.cond != nil {
			mo.cond.Broadcast()
		}
	}
	r.thrMu.Unlock()
	if traced {
		if tr := t.trace(); tr != nil {
			tr.monitor(t.ID, o.Identity, "exit", depth, m, pc, "")
		}
	}
	return nil
}

// holdsLockObject reports whether the current logical thread owns o's monitor.
func (t *Thread) holdsLockObject(o *Object) bool {
	if o == nil {
		return false
	}
	r := t.RT
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	mo := o.monitor
	return mo != nil && mo.owner == t && mo.depth > 0
}

// releaseMonitorForWaitLocked fully releases mo for a wait and returns the
// saved recursion depth. The caller holds thrMu.
func (r *Runtime) releaseMonitorForWaitLocked(t *Thread, o *Object, mo *monitor) int {
	saved := mo.depth
	mo.owner = nil
	mo.depth = 0
	delete(t.held, o)
	if mo.cond != nil {
		mo.cond.Broadcast()
	}
	return saved
}

// reacquireMonitor reacquires mo to its full saved recursion depth before a
// wait returns or throws. Monitor reacquisition is not interruptible; a later
// interrupt stays pending.
func (r *Runtime) reacquireMonitor(t *Thread, o *Object, mo *monitor, savedDepth int) {
	r.thrMu.Lock()
	for mo.owner != nil && mo.owner != t {
		if mo.cond == nil {
			mo.cond = sync.NewCond(&r.thrMu)
		}
		mo.cond.Wait()
	}
	mo.owner = t
	mo.depth = savedDepth
	if t.held == nil {
		t.held = map[*Object]*monitor{}
	}
	t.held[o] = mo
	r.thrMu.Unlock()
}

// releaseHeldMonitorsLocked releases every monitor a terminating thread still
// owns, waking entry contenders. The caller holds thrMu. Java releases all
// monitors when a thread terminates, including after an uncaught throwable.
func (r *Runtime) releaseHeldMonitorsLocked(t *Thread) {
	for o, mo := range t.held {
		if mo.owner == t {
			mo.owner = nil
			mo.depth = 0
			if mo.cond != nil {
				mo.cond.Broadcast()
			}
		}
		delete(t.held, o)
	}
}

// removeWaiterLocked removes w from its monitor's wait set. The caller holds
// thrMu.
func removeWaiterLocked(mo *monitor, w *waiter) {
	for i, existing := range mo.waiters {
		if existing == w {
			copy(mo.waiters[i:], mo.waiters[i+1:])
			last := len(mo.waiters) - 1
			// Clear the vacated slot: shortening the slice alone would keep the
			// completed waiter and its Java thread reachable through the
			// backing array for the lifetime of a long-lived monitor.
			mo.waiters[last] = nil
			mo.waiters = mo.waiters[:last]
			return
		}
	}
}

// wakeWaiterLocked performs the single winning transition of one wait record.
// Only a pending record can transition, so a stale timer or a later interrupt
// cannot overwrite the winner. The caller holds thrMu.
func (r *Runtime) wakeWaiterLocked(mo *monitor, w *waiter, state waiterState) bool {
	if w == nil || w.state != waiterPending {
		return false
	}
	w.state = state
	removeWaiterLocked(mo, w)
	if w.thread != nil && w.thread.blocker != nil && w.thread.blocker.waiter == w {
		w.thread.blocker = nil
	}
	// Release the timeout segment and its owner goroutine immediately. A stale
	// timer callback can no longer observe a pending record, and stopping the
	// host timer keeps its resources from being retained until the Java
	// deadline after an early notification or interruption.
	if w.timer != nil {
		w.timer.Stop()
		w.timer = nil
	}
	if w.timerDone != nil {
		close(w.timerDone)
		w.timerDone = nil
	}
	close(w.ch)
	return true
}

// startWaitTimer owns the timeout transition for one wait record. It walks the
// absolute deadline in bounded segments so Java timeouts longer than Go's
// duration range remain exact, stops its current segment when the record wins
// early, and exits as soon as timerDone closes. At most one timer is armed per
// waiter, and a stale callback never sees a pending record.
func (r *Runtime) startWaitTimer(w *waiter, mo *monitor, deadline waitTimer, done chan struct{}) {
	go func() {
		r.timerWorkers.Add(1)
		defer r.timerWorkers.Add(-1)
		for {
			remaining, expired := deadline.remaining(r.clock.Now())
			if expired {
				r.thrMu.Lock()
				r.wakeWaiterLocked(mo, w, waiterTimedOut)
				r.thrMu.Unlock()
				return
			}
			timer := r.clock.NewTimer(r.timerSegment(remaining))
			r.thrMu.Lock()
			if w.state != waiterPending {
				r.thrMu.Unlock()
				timer.Stop()
				return
			}
			w.timer = timer
			r.thrMu.Unlock()
			select {
			case <-timer.C():
			case <-done:
				timer.Stop()
				return
			}
			r.thrMu.Lock()
			pending := w.state == waiterPending
			r.thrMu.Unlock()
			if !pending {
				return
			}
		}
	}()
}

// objectWait implements Object.wait()V, wait(J)V and wait(JI)V. Argument
// validation has already happened; ownership invalidates with
// IllegalMonitorStateException before any interrupt status is consumed. A
// timed wait parks against one absolute deadline; its timer goroutine releases
// itself immediately when notification or interruption wins.
func (t *Thread) objectWait(obj *Object, timeoutMillis int64, timed bool) Outcome {
	r := t.RT
	var deadline waitTimer
	if timed {
		deadline = r.newWaitTimer(timeoutMillis)
	}
	r.thrMu.Lock()
	mo := r.monitorForLocked(obj)
	if mo.owner != t {
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalMonitorStateException"))
	}
	if t.interruptPendingLocked() {
		t.clearInterruptLocked()
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/InterruptedException"))
	}
	saved := mo.depth
	w := &waiter{thread: t, monitor: mo, savedDepth: saved, state: waiterPending,
		ch: make(chan struct{}), object: obj}
	var timerDone chan struct{}
	if timed {
		timerDone = make(chan struct{})
		w.timerDone = timerDone
	}
	mo.waiters = append(mo.waiters, w)
	t.blocker = &threadBlocker{kind: blockerWait, waiter: w}
	r.releaseMonitorForWaitLocked(t, obj, mo)
	r.thrMu.Unlock()
	if tr := t.trace(); tr != nil {
		tr.monitor(t.ID, obj.Identity, "wait", saved, nil, -1, "")
	}
	if timed {
		r.startWaitTimer(w, mo, deadline, timerDone)
	}
	<-w.ch
	r.reacquireMonitor(t, obj, mo, saved)
	r.thrMu.Lock()
	state := w.state
	if state == waiterInterrupted {
		t.clearInterruptLocked()
	}
	if t.blocker != nil && t.blocker.waiter == w {
		t.blocker = nil
	}
	r.thrMu.Unlock()
	reason := "notified"
	switch state {
	case waiterTimedOut:
		reason = "timed_out"
	case waiterInterrupted:
		reason = "interrupted"
	}
	if tr := t.trace(); tr != nil {
		tr.monitor(t.ID, obj.Identity, "resume", saved, nil, -1, reason)
	}
	if state == waiterInterrupted {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/InterruptedException"))
	}
	return Normal()
}

// notifyObject implements Object.notify()V. The notifier retains ownership and
// notification is never buffered for future waiters.
func (t *Thread) notifyObject(obj *Object, all bool) Outcome {
	r := t.RT
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	mo := r.monitorForLocked(obj)
	if mo.owner != t {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/IllegalMonitorStateException"))
	}
	for i := 0; i < len(mo.waiters); {
		w := mo.waiters[i]
		if w.state != waiterPending {
			i++
			continue
		}
		r.wakeWaiterLocked(mo, w, waiterNotified)
		if !all {
			break
		}
	}
	return Normal()
}

// joinPeer waits for target to terminate. A timed join preserves one monotonic
// deadline across unrelated wakeups: every wakeup recomputes the remaining time
// from the same deadline, and each wait record owns at most one bounded timer
// segment that is stopped when the join is woken early. Interruption throws
// from the joiner with its status cleared and releases only the target's
// monitor.
func (t *Thread) joinPeer(target *Object, timeoutMillis int64, timed bool) Outcome {
	r := t.RT
	if thrown := t.monitorEnterObject(target, nil, -1); thrown != nil {
		return ThrowOutcome(thrown)
	}
	var deadline waitTimer
	if timed {
		deadline = r.newWaitTimer(timeoutMillis)
	}
	interrupted := false
	for {
		r.thrMu.Lock()
		p, ok := peerOfObject(target)
		if !ok || p.state != threadStarted {
			r.thrMu.Unlock()
			break
		}
		if t.interruptPendingLocked() {
			t.clearInterruptLocked()
			r.thrMu.Unlock()
			interrupted = true
			break
		}
		if timed {
			if _, expired := deadline.remaining(r.clock.Now()); expired {
				r.thrMu.Unlock()
				break
			}
		}
		mo := r.monitorForLocked(target)
		if mo.owner != t {
			// The target monitor was relinquished by hand-authored code; the
			// wait protocol cannot proceed.
			r.thrMu.Unlock()
			interrupted = false
			break
		}
		saved := mo.depth
		w := &waiter{thread: t, monitor: mo, savedDepth: saved, state: waiterPending,
			ch: make(chan struct{}), object: target}
		var timerDone chan struct{}
		if timed {
			timerDone = make(chan struct{})
			w.timerDone = timerDone
		}
		mo.waiters = append(mo.waiters, w)
		t.blocker = &threadBlocker{kind: blockerJoin, waiter: w}
		r.releaseMonitorForWaitLocked(t, target, mo)
		r.thrMu.Unlock()
		if timed {
			r.startWaitTimer(w, mo, deadline, timerDone)
		}
		<-w.ch
		r.reacquireMonitor(t, target, mo, saved)
		r.thrMu.Lock()
		if w.state == waiterInterrupted {
			t.clearInterruptLocked()
			interrupted = true
			r.thrMu.Unlock()
			break
		}
		r.thrMu.Unlock()
	}
	if thrown := t.monitorExitObject(target, nil, -1); thrown != nil {
		return ThrowOutcome(thrown)
	}
	if interrupted {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/InterruptedException"))
	}
	return Normal()
}

// sleepFor implements Thread.sleep. The logical thread keeps every monitor it
// owns. A pending interrupt makes even a zero-duration sleep throw. The timer
// is owned directly by the sleeping goroutine, so an interrupted sleep stops it
// instead of leaving armed timer work behind.
func (t *Thread) sleepFor(d time.Duration) Outcome {
	r := t.RT
	r.thrMu.Lock()
	if t.interruptPendingLocked() {
		t.clearInterruptLocked()
		r.thrMu.Unlock()
		return ThrowOutcome(r.NewThrowable(t, "java/lang/InterruptedException"))
	}
	if d <= 0 {
		r.thrMu.Unlock()
		return Normal()
	}
	b := &threadBlocker{kind: blockerSleep, ch: make(chan struct{}), sleep: waiterPending}
	t.blocker = b
	r.thrMu.Unlock()
	timer := r.clock.NewTimer(d)
	select {
	case <-timer.C():
		// The timer won; a later interrupt only leaves its status pending.
		r.thrMu.Lock()
		if b.sleep == waiterPending {
			b.sleep = waiterTimedOut
		}
		r.thrMu.Unlock()
	case <-b.ch:
		timer.Stop()
	}
	r.thrMu.Lock()
	state := b.sleep
	if t.blocker == b {
		t.blocker = nil
	}
	interrupted := state == waiterInterrupted
	if interrupted {
		t.clearInterruptLocked()
	}
	r.thrMu.Unlock()
	if interrupted {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/InterruptedException"))
	}
	return Normal()
}

// interruptPeer records the Java interrupt status and wakes an interruptible
// blocker. A peer that has not started, or has terminated, keeps the bit.
func (t *Thread) interruptPeer(peer *Object) {
	r := t.RT
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	p, ok := peerOfObject(peer)
	if !ok {
		return
	}
	p.pendingInterrupt = true
	owner := p.thread
	if owner == nil || owner.blocker == nil {
		return
	}
	b := owner.blocker
	switch b.kind {
	case blockerWait, blockerJoin:
		if b.waiter != nil {
			r.wakeWaiterLocked(b.waiter.monitor, b.waiter, waiterInterrupted)
		}
	case blockerSleep:
		if b.sleep == waiterPending {
			b.sleep = waiterInterrupted
			close(b.ch)
		}
	}
}

// currentPeer returns the calling logical thread's Java peer, creating one for
// a host thread that entered Java without a peer.
func (t *Thread) currentPeer() *Object {
	r := t.RT
	r.thrMu.Lock()
	if t.peer != nil {
		peer := t.peer
		r.thrMu.Unlock()
		return peer
	}
	r.thrMu.Unlock()
	peer := r.newPeerObject(fmt.Sprintf("Thread-%d", atomic.AddInt64(&r.unnamedPeers, 1)), false, threadStarted)
	r.thrMu.Lock()
	if t.peer == nil {
		t.peer = peer
		if p, ok := peerOfObject(peer); ok {
			p.thread = t
		}
		r.noteLiveThreadStartLocked()
	}
	result := t.peer
	r.thrMu.Unlock()
	return result
}
