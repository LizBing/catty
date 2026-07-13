// Package prototypes contains research-only concurrency mechanism experiments.
// It is not linked into catty production code.
package prototypes

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var (
	errIllegalMonitorState = errors.New("illegal monitor state")
	errInterrupted         = errors.New("interrupted")
)

type waiterState uint8

const (
	waiting waiterState = iota
	notified
	interrupted
)

type waiter struct {
	thread uint64
	state  waiterState
	wake   chan struct{}
}

// researchMonitor models the minimum mechanism needed behind a Java object.
// Ownership is a Java execution-context ID, never an inferred goroutine ID.
type researchMonitor struct {
	mu      sync.Mutex
	entry   *sync.Cond
	owner   uint64
	depth   int
	waiters []*waiter
}

func newResearchMonitor() *researchMonitor {
	m := &researchMonitor{}
	m.entry = sync.NewCond(&m.mu)
	return m
}

func (m *researchMonitor) enter(thread uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for m.owner != 0 && m.owner != thread {
		m.entry.Wait()
	}
	if m.owner == thread {
		m.depth++
		return
	}
	m.owner = thread
	m.depth = 1
}

func (m *researchMonitor) exit(thread uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != thread || m.depth == 0 {
		return errIllegalMonitorState
	}
	m.depth--
	if m.depth == 0 {
		m.owner = 0
		m.entry.Broadcast()
	}
	return nil
}

func (m *researchMonitor) wait(thread uint64, installed chan<- struct{}) error {
	m.mu.Lock()
	if m.owner != thread || m.depth == 0 {
		m.mu.Unlock()
		return errIllegalMonitorState
	}

	savedDepth := m.depth
	w := &waiter{thread: thread, state: waiting, wake: make(chan struct{})}
	m.waiters = append(m.waiters, w)
	m.owner = 0
	m.depth = 0
	m.entry.Broadcast()
	m.mu.Unlock()
	if installed != nil {
		installed <- struct{}{}
	}

	<-w.wake
	m.enter(thread)

	m.mu.Lock()
	m.depth = savedDepth
	state := w.state
	m.removeWaiterLocked(w)
	m.mu.Unlock()
	if state == interrupted {
		return errInterrupted
	}
	return nil
}

func (m *researchMonitor) notifyOne(thread uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != thread || m.depth == 0 {
		return errIllegalMonitorState
	}
	for _, w := range m.waiters {
		if w.state == waiting {
			w.state = notified
			close(w.wake)
			return nil
		}
	}
	return nil
}

func (m *researchMonitor) notifyAll(thread uint64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.owner != thread || m.depth == 0 {
		return errIllegalMonitorState
	}
	for _, w := range m.waiters {
		if w.state == waiting {
			w.state = notified
			close(w.wake)
		}
	}
	return nil
}

// interrupt orders interruption with notification under the same state lock.
// If interruption wins, a later notify skips this waiter and selects another.
func (m *researchMonitor) interrupt(thread uint64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, w := range m.waiters {
		if w.thread == thread && w.state == waiting {
			w.state = interrupted
			close(w.wake)
			return true
		}
	}
	return false
}

func (m *researchMonitor) removeWaiterLocked(target *waiter) {
	for i, w := range m.waiters {
		if w == target {
			copy(m.waiters[i:], m.waiters[i+1:])
			m.waiters[len(m.waiters)-1] = nil
			m.waiters = m.waiters[:len(m.waiters)-1]
			return
		}
	}
}

func TestGoMutexIsNotJavaReentrant(t *testing.T) {
	var mu sync.Mutex
	mu.Lock()
	if mu.TryLock() {
		t.Fatal("sync.Mutex unexpectedly allowed recursive acquisition")
	}
	mu.Unlock()
}

func TestMonitorExclusionAndReentrancy(t *testing.T) {
	m := newResearchMonitor()
	m.enter(1)
	m.enter(1)

	entered := make(chan struct{})
	go func() {
		m.enter(2)
		close(entered)
		_ = m.exit(2)
	}()

	select {
	case <-entered:
		t.Fatal("second Java thread entered while recursion depth was nonzero")
	case <-time.After(time.Millisecond):
	}
	if err := m.exit(1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
		t.Fatal("second Java thread entered before final recursive exit")
	case <-time.After(time.Millisecond):
	}
	if err := m.exit(1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("second Java thread never entered")
	}
}

func TestWaitReleasesAndRestoresRecursionDepth(t *testing.T) {
	m := newResearchMonitor()
	installed := make(chan struct{}, 1)
	done := make(chan error, 1)

	go func() {
		m.enter(1)
		m.enter(1)
		err := m.wait(1, installed)
		if err == nil {
			err = m.exit(1)
		}
		if err == nil {
			err = m.exit(1)
		}
		done <- err
	}()
	<-installed

	m.enter(2)
	if err := m.notifyOne(2); err != nil {
		t.Fatal(err)
	}
	if err := m.exit(2); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatalf("wait/restore failed: %v", err)
	}
}

func TestInterruptBeforeNotifyDoesNotConsumeNotification(t *testing.T) {
	m := newResearchMonitor()
	installed := make(chan struct{}, 2)
	results := make(chan error, 2)

	for _, id := range []uint64{1, 2} {
		go func(thread uint64) {
			m.enter(thread)
			results <- m.wait(thread, installed)
			_ = m.exit(thread)
		}(id)
	}
	<-installed
	<-installed

	if !m.interrupt(1) {
		t.Fatal("thread 1 was not in the wait set")
	}
	m.enter(3)
	if err := m.notifyOne(3); err != nil {
		t.Fatal(err)
	}
	if err := m.exit(3); err != nil {
		t.Fatal(err)
	}

	first := <-results
	second := <-results
	interruptedCount := 0
	for _, err := range []error{first, second} {
		if errors.Is(err, errInterrupted) {
			interruptedCount++
		} else if err != nil {
			t.Fatalf("unexpected waiter result: %v", err)
		}
	}
	if interruptedCount != 1 {
		t.Fatalf("interrupted waiters=%d, want 1", interruptedCount)
	}
}

type initState uint8

const (
	initNotStarted initState = iota
	initInProgress
	initDone
	initFailed
)

type initCell struct {
	mu    sync.Mutex
	cond  *sync.Cond
	state initState
	owner uint64
}

func newInitCell() *initCell {
	c := &initCell{}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *initCell) initialize(thread uint64, run func() error) error {
	c.mu.Lock()
	for c.state == initInProgress && c.owner != thread {
		c.cond.Wait()
	}
	switch c.state {
	case initDone:
		c.mu.Unlock()
		return nil
	case initFailed:
		c.mu.Unlock()
		return errors.New("erroneous")
	case initInProgress:
		c.mu.Unlock()
		return nil
	}
	c.state = initInProgress
	c.owner = thread
	c.mu.Unlock()

	err := run()
	c.mu.Lock()
	if err == nil {
		c.state = initDone
	} else {
		c.state = initFailed
	}
	c.owner = 0
	c.cond.Broadcast()
	c.mu.Unlock()
	return err
}

func TestCrossThreadInitializationRunsOnceAndPublishes(t *testing.T) {
	cell := newInitCell()
	var runs atomic.Int32
	var published int32
	start := make(chan struct{})
	results := make(chan error, 2)

	for _, id := range []uint64{1, 2} {
		go func(thread uint64) {
			<-start
			results <- cell.initialize(thread, func() error {
				runs.Add(1)
				published = 42
				return nil
			})
		}(id)
	}
	close(start)
	if err := <-results; err != nil {
		t.Fatal(err)
	}
	if err := <-results; err != nil {
		t.Fatal(err)
	}
	if got := runs.Load(); got != 1 {
		t.Fatalf("initializer runs=%d, want 1", got)
	}
	if published != 42 {
		t.Fatalf("published=%d, want 42", published)
	}
}

func TestSequentiallyConsistentAtomicCanBackBoundedVolatile(t *testing.T) {
	var data int32
	var ready atomic.Bool
	done := make(chan int32, 1)
	go func() {
		for !ready.Load() {
		}
		done <- data
	}()
	data = 42
	ready.Store(true)
	if got := <-done; got != 42 {
		t.Fatalf("published=%d, want 42", got)
	}
}
