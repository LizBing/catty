package kernel

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMonitorReentrant(t *testing.T) {
	m := &Monitor{}
	const key = 7
	m.Enter(key)
	m.Enter(key)
	m.Enter(key)
	if err := m.Exit(key); err != nil {
		t.Fatal(err)
	}
	if !m.IsOwnedBy(key) {
		t.Fatal("still owned after one exit at depth 3")
	}
	if err := m.Exit(key); err != nil {
		t.Fatal(err)
	}
	if err := m.Exit(key); err != nil {
		t.Fatal(err)
	}
	if m.IsOwnedBy(key) {
		t.Fatal("must be free after full depth release")
	}
}

func TestMonitorExitNotOwner(t *testing.T) {
	m := &Monitor{}
	if err := m.Exit(42); err != ErrNotMonitorOwner {
		t.Fatalf("want ErrNotMonitorOwner, got %v", err)
	}
	m.Enter(1)
	if err := m.Exit(2); err != ErrNotMonitorOwner {
		t.Fatalf("want ErrNotMonitorOwner for foreign key, got %v", err)
	}
}

func TestMonitorHandoff(t *testing.T) {
	m := &Monitor{}
	m.Enter(1)

	acquired := make(chan struct{})
	go func() {
		m.Enter(2)
		close(acquired)
	}()

	// give B time to park
	time.Sleep(20 * time.Millisecond)
	select {
	case <-acquired:
		t.Fatal("B acquired while A still owns")
	default:
	}

	if err := m.Exit(1); err != nil {
		t.Fatal(err)
	}
	<-acquired // must acquire after handoff

	// cleanup: B exits so the monitor is free for later tests
	if err := m.Exit(2); err != nil {
		t.Fatal(err)
	}
}

func TestWaitNotifyHandshake(t *testing.T) {
	m := &Monitor{}
	done := make(chan waitOutcome, 1)

	go func() {
		m.Enter(10)
		out := m.Wait(nil, 10, -1) // indefinitely
		done <- out           // reacquired on wake; ownership verified below
	}()

	// A holds-and-waits → lock free; B acquires and notifies once the
	// waiter record is actually queued.
	waiters := func() int {
		m.mu.Lock()
		defer m.mu.Unlock()
		return len(m.waits)
	}
	deadline := time.Now().Add(time.Second)
	for waiters() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("waiter never parked")
		}
		time.Sleep(time.Millisecond)
	}

	m.Enter(20)
	notified, err := m.Notify(20)
	if err != nil || !notified {
		t.Fatalf("notify = %v,%v", notified, err)
	}
	m.Exit(20)

	out := <-done
	if out != waitGotNotify {
		t.Fatalf("outcome = %v, want waitGotNotify", out)
	}
	if !m.IsOwnedBy(10) {
		t.Fatal("waiter must own the monitor again before Wait returns")
	}
	if err := m.Exit(10); err != nil {
		t.Fatal(err)
	}
}

func TestWaitTimeoutRestoresDepth(t *testing.T) {
	m := &Monitor{}
	m.Enter(5)
	m.Enter(5) // depth 2
	start := time.Now()
	out := m.Wait(nil, 5, 40)
	if out != waitHitDeadline {
		t.Fatalf("outcome = %v, want waitHitDeadline", out)
	}
	if elapsed := time.Since(start); elapsed < 35*time.Millisecond {
		t.Fatalf("returned too early: %v", elapsed)
	}
	if !m.IsOwnedBy(5) {
		t.Fatal("ownership must be restored after timeout")
	}
	m.Exit(5)
	m.Exit(5) // depth restored: second exit must succeed
	if m.IsOwnedBy(5) {
		t.Fatal("fully released")
	}
}

func TestInterruptedWait(t *testing.T) {
	m := &Monitor{}
	outCh := make(chan waitOutcome, 1)
	go func() {
		m.Enter(9)
		outCh <- m.Wait(nil, 9, -1)
	}()
	waiters := func() int {
		m.mu.Lock()
		defer m.mu.Unlock()
		return len(m.waits)
	}
	deadline := time.Now().Add(time.Second)
	for waiters() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("waiter never parked")
		}
		time.Sleep(time.Millisecond)
	}
	if n := m.InterruptWaiters(9); n != 1 {
		t.Fatalf("interrupted=%d, want 1", n)
	}
	if out := <-outCh; out != waitGotInterrupt {
		t.Fatalf("outcome = %v, want waitGotInterrupt", out)
	}
	m.Exit(9)
}

func TestNotifyAllWakesAll(t *testing.T) {
	m := &Monitor{}
	var woke atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Enter(100)
			if o := m.Wait(nil, 100, 2000); o == waitGotNotify {
				woke.Add(1)
			}
			m.Exit(100)
		}()
	}
	// wait until all three are parked in m.waits
	dl := time.Now().Add(time.Second)
	for {
		m.mu.Lock()
		n := len(m.waits)
		m.mu.Unlock()
		if n == 3 {
			break
		}
		if time.Now().After(dl) {
			t.Fatalf("only %d waiters parked", n)
		}
		time.Sleep(time.Millisecond)
	}
	m.Enter(200)
	if _, err := m.NotifyAll(200); err != nil {
		t.Fatal(err)
	}
	m.Exit(200)
	wg.Wait()
	if got := woke.Load(); got != 3 {
		t.Fatalf("woke=%d, want 3", got)
	}
}

func TestNotifyWithoutWaitersIsNoop(t *testing.T) {
	m := &Monitor{}
	m.Enter(3)
	notified, err := m.Notify(3)
	if err != nil || notified {
		t.Fatalf("notify with no waiters = %v,%v", notified, err)
	}
	m.Exit(3)
}
