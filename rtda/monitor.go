package rtda

import (
	"sync"
	"sync/atomic"
)

// Monitor is the runtime kernel state for a Java object monitor (ADR-0029).
// Every Object, array, and canonical Class mirror may lazily attach one.
// The monitor is NOT copied by CloneObject.
//
// A Monitor is a reentrant mutex keyed by execution-context identity.
// It manages an ordered collection of wait-set entries, each carrying a
// private wake signal so notify, notifyAll, and waiter interruption can
// be totally ordered under the monitor state lock.
type Monitor struct {
	mu    sync.Mutex
	cond  sync.Cond // signaled on ownership release (depth → 0)

	owner uint64 // execution-context ID, 0 = no owner
	depth int    // recursion count

	waitSet []*waiterEntry
}

// waiter states (atomic int32).
const (
	waitWaiting    int32 = iota // enrolled, not yet notified or interrupted
	waitNotified                // notify won the race
	waitInterrupted             // interrupt won the race
)

// waiterEntry represents a Java thread blocked in Object.wait() on this monitor.
// State transitions are single-shot: waiting → notified or waiting → interrupted.
// Wake signals are exactly-once via a buffered (cap 1) channel.
type waiterEntry struct {
	state int32        // atomic; one of waitWaiting, waitNotified, waitInterrupted
	ecID  uint64       // the waiting execution context
	wake  chan struct{} // buffered 1 — private wake signal
}

// --- Monitor creation (lazy, CAS-published via Object.Monitor()) ---

func newMonitor() *Monitor {
	m := &Monitor{}
	m.cond.L = &m.mu
	return m
}

// --- Monitor() accessor on Object ---

// Monitor returns the lazy CAS-published monitor for this object per ADR-0029.
// The first call for a given object allocates the Monitor; all subsequent
// callers (including concurrent ones) see the same *Monitor.
func (obj *Object) Monitor() *Monitor {
	if m := obj.monitor.Load(); m != nil {
		return m
	}
	m := newMonitor()
	if obj.monitor.CompareAndSwap(nil, m) {
		return m
	}
	// Lost the race — return the winner.
	return obj.monitor.Load()
}

// --- Monitor operations ---

// Enter acquires the monitor for ecID. If ecID already owns the monitor
// (reentrant entry), depth is incremented. If another execution context owns
// it, the calling goroutine blocks until the monitor is released. Entry is
// non-interruptible per ADR-0029.
func (m *Monitor) Enter(ecID uint64) {
	m.mu.Lock()
	if m.owner == ecID {
		m.depth++
		m.mu.Unlock()
		return
	}
	for m.owner != 0 {
		m.cond.Wait()
	}
	m.owner = ecID
	m.depth = 1
	m.mu.Unlock()
}

// Exit releases one recursion level of the monitor. The caller ecID MUST be
// the current owner or Exit returns false (the caller should throw
// IllegalMonitorStateException). When depth reaches zero, ownership is
// released and one blocked enterer is woken.
func (m *Monitor) Exit(ecID uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != ecID {
		return false
	}
	m.depth--
	if m.depth == 0 {
		m.owner = 0
		m.cond.Signal()
	}
	return true
}

// HoldsLock reports whether ecID currently owns this monitor.
func (m *Monitor) HoldsLock(ecID uint64) bool {
	m.mu.Lock()
	held := m.owner == ecID
	m.mu.Unlock()
	return held
}

// RecursionDepth returns the current recursion depth (0 if unowned).
func (m *Monitor) RecursionDepth() int {
	m.mu.Lock()
	d := m.depth
	m.mu.Unlock()
	return d
}

// OwnerEC returns the execution-context ID of the current owner, or 0 if unowned.
func (m *Monitor) OwnerEC() uint64 {
	m.mu.Lock()
	o := m.owner
	m.mu.Unlock()
	return o
}

// --- Wait / Notify / NotifyAll ---

// InternalWait is the monitor-side half of Object.wait(). It MUST be called
// after the Thread's active-waiter pre-check has succeeded (the Thread has
// already published its waitingOn pointer).
//
// InternalWait saves the current recursion depth, fully releases the monitor,
// enqueues a waiter entry, and blocks. On wake, it reacquires the monitor and
// restores the saved depth. It returns true if the wait completed normally
// (notify won), or false if the waiter was interrupted.
func (m *Monitor) InternalWait(ecID uint64, savedDepth int) bool {
	m.mu.Lock()
	// Verify ownership — should never fail if the caller already holds the lock.
	if m.owner != ecID {
		m.mu.Unlock()
		panic("catty: InternalWait called by non-owner")
	}
	if savedDepth < 1 {
		savedDepth = 1 // defensive: at least one level
	}

	// Enqueue waiter.
	we := &waiterEntry{
		state: waitWaiting,
		ecID:  ecID,
		wake:  make(chan struct{}, 1),
	}
	m.waitSet = append(m.waitSet, we)

	// Fully release: depth reset and signal one blocked enterer.
	m.depth = 0
	m.owner = 0
	m.cond.Signal()
	m.mu.Unlock()

	// Block until woken by notify, notifyAll, or interrupt.
	<-we.wake

	// Reacquire the monitor.
	m.Enter(ecID)

	// Restore saved recursion depth.
	m.mu.Lock()
	if m.owner == ecID {
		m.depth = savedDepth
	}
	// Remove this waiter entry from the wait set so a subsequent wait()
	// by the same execution context produces a fresh entry that
	// InterruptWaiter and Notify/NotifyAll can correctly locate.
	m.removeWaiterEntry(we)
	m.mu.Unlock()

	// Determine wake cause.
	finalState := atomic.LoadInt32(&we.state)
	return finalState != waitInterrupted
}

// removeWaiterEntry removes a specific entry from the wait set. The caller
// must hold m.mu. Called from InternalWait after the waiter has reacquired
// the monitor so stale entries do not accumulate.
func (m *Monitor) removeWaiterEntry(target *waiterEntry) {
	for i, we := range m.waitSet {
		if we == target {
			m.waitSet = append(m.waitSet[:i], m.waitSet[i+1:]...)
			return
		}
	}
}

// Notify transitions one currently-waiting waiter to notified and wakes it.
// The caller ecID must own the monitor; returns false for non-owner.
func (m *Monitor) Notify(ecID uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != ecID {
		return false
	}
	for _, we := range m.waitSet {
		if atomic.CompareAndSwapInt32(&we.state, waitWaiting, waitNotified) {
			select {
			case we.wake <- struct{}{}:
			default:
			}
			return true
		}
	}
	return true // no waiting waiter to notify — not an error
}

// NotifyAll transitions all currently-waiting waiters to notified and wakes
// each of them. The caller ecID must own the monitor; returns false for
// non-owner.
func (m *Monitor) NotifyAll(ecID uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != ecID {
		return false
	}
	for _, we := range m.waitSet {
		if atomic.CompareAndSwapInt32(&we.state, waitWaiting, waitNotified) {
			select {
			case we.wake <- struct{}{}:
			default:
			}
		}
	}
	return true
}

// InterruptWaiter is called from Thread.Interrupt when the target thread has
// an active wait-set enrollment. It scans for the first waiter entry whose EC
// ID matches and whose state is still waitWaiting, then attempts to atomically
// transition it to waitInterrupted. Stale entries (already notified or
// interrupted) for the same EC are skipped, so a second wait() by the same
// execution context is correctly located.
//
// If notify already won the race on every matching entry, the interrupt is
// silently recorded (interruptState is already set) — the waiter will wake
// normally and the interrupt remains pending post-wait.
func (m *Monitor) InterruptWaiter(ecID uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, we := range m.waitSet {
		if we.ecID != ecID {
			continue
		}
		if atomic.LoadInt32(&we.state) != waitWaiting {
			continue // stale entry — a previous wait() already completed
		}
		if atomic.CompareAndSwapInt32(&we.state, waitWaiting, waitInterrupted) {
			select {
			case we.wake <- struct{}{}:
			default:
			}
		}
		// Found a waiting entry for this EC — whether the CAS succeeded or
		// a concurrent notify won, there is at most one active waiter per
		// execution context. Stop searching.
		return
	}
}

// WaitSetLen returns the current number of entries in the wait set.
// Exported for testing.
func (m *Monitor) WaitSetLen() int {
	m.mu.Lock()
	n := len(m.waitSet)
	m.mu.Unlock()
	return n
}
