package kernel

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Monitor implements Java intrinsic-monitor semantics (JVMS §2.2.8):
//
//   - reentrant ownership with depth counting;
//   - parked-contender queue for lock handoff (barging permitted — matches
//     thin-lock behavior; the spec does not mandate fairness);
//   - Object.wait: full-depth release, park, reacquire to prior depth;
//   - notify/notifyAll with arbitrary-waiter selection as permitted;
//   - interrupt-aware waiting: waiters carry an outcome so M2's
//     Thread.interrupt can mark+wake them without polling.
//
// Thread identity is an explicit opaque key (Go deliberately exposes no
// goroutine IDs); key 0 is reserved for "no thread".
type Monitor struct {
	mu    sync.Mutex
	owner uint64 // 0 = unowned
	count int    // reentrancy depth

	parked []chan struct{} // contenders for the lock, FIFO

	waits []*waitRec // Object.wait records
}

type waitRec struct {
	key uint64
	ch  chan struct{}
	// outcome is written by notifiers/interrupters under m.mu, but read
	// lock-free by the owning waiter — hence atomic.
	outcome atomic.Int32
}


// Enter acquires the monitor, reentrant for the given owner key.
func (m *Monitor) Enter(key uint64) {
	m.enterDepth(key, 1)
}

func (m *Monitor) enterDepth(key uint64, depth int) {
	if key == 0 {
		panic("kernel: monitor key 0 is reserved")
	}
	m.mu.Lock()
	for {
		if m.owner == 0 {
			m.owner = key
			m.count = depth
			m.mu.Unlock()
			return
		}
		if m.owner == key {
			m.count += depth
			m.mu.Unlock()
			return
		}
		ch := make(chan struct{})
		m.parked = append(m.parked, ch)
		m.mu.Unlock()
		<-ch
		m.mu.Lock()
		// woken: loop and retry acquisition (handoff or barging)
	}
}

// Exit releases one level of ownership. Returns ErrNotMonitorOwner when
// key does not own the monitor (IllegalMonitorStateException upstream).
func (m *Monitor) Exit(key uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != key {
		return ErrNotMonitorOwner
	}
	m.count--
	if m.count > 0 {
		return nil
	}
	m.releaseLocked()
	return nil
}

// releaseLocked clears ownership and hands off to one parked contender.
// Caller must hold m.mu.
func (m *Monitor) releaseLocked() {
	m.owner = 0
	m.count = 0
	if len(m.parked) > 0 {
		ch := m.parked[0]
		m.parked = m.parked[1:]
		close(ch)
	}
}

// IsOwnedBy reports whether key currently owns the monitor.
func (m *Monitor) IsOwnedBy(key uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.owner == key
}

var ErrNotMonitorOwner = errors.New("current thread does not own the monitor")

// waitOutcome is the terminal state of one Wait call.
type waitOutcome int32

const (
	waitStillWaiting waitOutcome = iota
	waitGotNotify
	waitHitDeadline
	waitGotInterrupt
)

// Wait implements Object.wait(timeoutMillis): full-depth release, park
// until notify/interrupt/deadline, then reacquire to the prior depth.
// timeoutMillis <= 0 means indefinitely.
func (m *Monitor) Wait(key uint64, timeoutMillis int64) waitOutcome {
	m.mu.Lock()
	if m.owner != key {
		m.mu.Unlock()
		panic("kernel: Wait without ownership (native layer must guard)")
	}
	depth := m.count
	rec := &waitRec{key: key, ch: make(chan struct{})}
	m.waits = append(m.waits, rec)

	// Release fully and let a contender take the lock.
	m.count = 0
	m.owner = 0
	m.releaseLocked()
	m.mu.Unlock()

	var deadline <-chan time.Time
	if timeoutMillis > 0 {
		t := time.NewTimer(time.Duration(timeoutMillis) * time.Millisecond)
		defer t.Stop()
		deadline = t.C
	}

	for rec.outcome.Load() == int32(waitStillWaiting) {
		if deadline == nil {
			<-rec.ch
			continue
		}
		select {
		case <-rec.ch:
		case <-deadline:
			m.mu.Lock()
			if rec.outcome.Load() == int32(waitStillWaiting) {
				rec.outcome.Store(int32(waitHitDeadline))
				m.removeWaitLocked(rec)
			}
			m.mu.Unlock()
		}
	}

	m.removeWait(rec)

	// Reacquire at prior depth.
	m.reenter(key, depth)
	return waitOutcome(rec.outcome.Load())
}

// removeWait drops a record from the waiter list if still present.
func (m *Monitor) removeWait(rec *waitRec) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, w := range m.waits {
		if w == rec {
			m.waits = append(m.waits[:i], m.waits[i+1:]...)
			return
		}
	}
}

// removeWaitLocked is the locked fast path used inside Wait's deadline arm.
func (m *Monitor) removeWaitLocked(rec *waitRec) {
	for i, w := range m.waits {
		if w == rec {
			m.waits = append(m.waits[:i], m.waits[i+1:]...)
			return
		}
	}
}

// reenter acquires for key at an exact depth (post-wait restore).
func (m *Monitor) reenter(key uint64, depth int) {
	m.mu.Lock()
	for {
		if m.owner == 0 {
			m.owner = key
			m.count = depth
			m.mu.Unlock()
			return
		}
		if m.owner == key { // possible if a same-key notifier re-entered first
			m.count += depth
			m.mu.Unlock()
			return
		}
		ch := make(chan struct{})
		m.parked = append(m.parked, ch)
		m.mu.Unlock()
		<-ch
		m.mu.Lock()
	}
}

// Notify wakes one arbitrary waiter. Returns false when none exist.
func (m *Monitor) Notify(key uint64) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != key {
		return false, ErrNotMonitorOwner
	}
	for i, w := range m.waits {
		if w.outcome.Load() == int32(waitStillWaiting) {
			w.outcome.Store(int32(waitGotNotify))
			m.waits = append(m.waits[:i], m.waits[i+1:]...)
			close(w.ch)
			return true, nil
		}
	}
	return false, nil
}

// NotifyAll wakes every waiter.
func (m *Monitor) NotifyAll(key uint64) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != key {
		return 0, ErrNotMonitorOwner
	}
	n := 0
	for _, w := range m.waits {
		if w.outcome.Load() == int32(waitStillWaiting) {
			w.outcome.Store(int32(waitGotNotify))
			close(w.ch)
			n++
		}
	}
	m.waits = nil
	return n, nil
}

// InterruptWaiters marks and wakes all waiters belonging to key.
// M2's Thread.interrupt calls this; the outcome lets Wait distinguish
// interrupt from notify.
func (m *Monitor) InterruptWaiters(key uint64) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	n := 0
	remaining := m.waits[:0]
	for _, w := range m.waits {
		if w.key == key && w.outcome.Load() == int32(waitStillWaiting) {
			w.outcome.Store(int32(waitGotInterrupt))
			close(w.ch)
			n++
			continue
		}
		remaining = append(remaining, w)
	}
	m.waits = remaining
	return n
}
