package rt

import (
	"fmt"
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
}

// startWaitTimer arms the single timeout transition for one wait record using
// the runtime clock.

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
	fired  bool
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
	if tr := t.trace(); tr != nil {
		if contended {
			tr.monitor(t.ID, o.Identity, "contended", depth, m, pc, "")
		}
		tr.monitor(t.ID, o.Identity, "enter", depth, m, pc, "")
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
	if tr := t.trace(); tr != nil {
		tr.monitor(t.ID, o.Identity, "exit", depth, m, pc, "")
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
			mo.waiters = append(mo.waiters[:i], mo.waiters[i+1:]...)
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
	close(w.ch)
	return true
}

// startWaitTimer arms the single timeout transition for one wait record using
// the runtime clock.
func (r *Runtime) startWaitTimer(w *waiter, mo *monitor, d time.Duration) {
	if d <= 0 {
		return
	}
	fired := r.clock.After(d)
	go func() {
		<-fired
		r.thrMu.Lock()
		r.wakeWaiterLocked(mo, w, waiterTimedOut)
		r.thrMu.Unlock()
	}()
}

// objectWait implements Object.wait()V, wait(J)V and wait(JI)V. Argument
// validation has already happened; ownership invalidates with
// IllegalMonitorStateException before any interrupt status is consumed.
func (t *Thread) objectWait(obj *Object, timeout time.Duration, timed bool) Outcome {
	r := t.RT
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
	mo.waiters = append(mo.waiters, w)
	t.blocker = &threadBlocker{kind: blockerWait, waiter: w}
	r.releaseMonitorForWaitLocked(t, obj, mo)
	r.thrMu.Unlock()
	if tr := t.trace(); tr != nil {
		tr.monitor(t.ID, obj.Identity, "wait", saved, nil, -1, "")
	}
	if timed {
		r.startWaitTimer(w, mo, timeout)
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
// deadline across unrelated wakeups. Interruption throws from the joiner with
// its status cleared and releases only the target's monitor.
func (t *Thread) joinPeer(target *Object, timeout time.Duration, timed bool) Outcome {
	r := t.RT
	if thrown := t.monitorEnterObject(target, nil, -1); thrown != nil {
		return ThrowOutcome(thrown)
	}
	var deadline time.Time
	if timed {
		deadline = r.clock.Now().Add(timeout)
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
		var remaining time.Duration
		if timed {
			remaining = deadline.Sub(r.clock.Now())
			if remaining <= 0 {
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
		mo.waiters = append(mo.waiters, w)
		t.blocker = &threadBlocker{kind: blockerJoin, waiter: w}
		r.releaseMonitorForWaitLocked(t, target, mo)
		r.thrMu.Unlock()
		if timed {
			r.startWaitTimer(w, mo, remaining)
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
// owns. A pending interrupt makes even a zero-duration sleep throw.
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
	fired := r.clock.After(d)
	done := make(chan struct{})
	go func() {
		select {
		case <-fired:
			r.thrMu.Lock()
			if b.sleep == waiterPending {
				b.sleep = waiterTimedOut
				if !b.fired {
					b.fired = true
					close(b.ch)
				}
			}
			r.thrMu.Unlock()
		case <-done:
		}
	}()
	<-b.ch
	close(done)
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
			if !b.fired {
				b.fired = true
				close(b.ch)
			}
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
	}
	result := t.peer
	r.thrMu.Unlock()
	return result
}
