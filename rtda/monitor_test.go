package rtda

import (
	"sync"
	"testing"
)

// --- Monitor identity and basic ownership ---

func TestMonitorLazyIdentity(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m1 := obj.Monitor()
	m2 := obj.Monitor()
	if m1 != m2 {
		t.Fatal("Monitor() returned different pointers")
	}
}

func TestMonitorDefaultUnlocked(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	if m.HoldsLock(1) {
		t.Fatal("fresh monitor should not be held by anyone")
	}
	if d := m.RecursionDepth(); d != 0 {
		t.Fatalf("fresh monitor depth = %d, want 0", d)
	}
}

// --- Enter / Exit ---

func TestMonitorEnterExit(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(42)

	m.Enter(ec)
	if !m.HoldsLock(ec) {
		t.Fatal("HoldsLock false after Enter")
	}
	if d := m.RecursionDepth(); d != 1 {
		t.Fatalf("depth = %d, want 1", d)
	}
	if !m.Exit(ec) {
		t.Fatal("Exit returned false for owner")
	}
	if m.HoldsLock(ec) {
		t.Fatal("HoldsLock true after Exit")
	}
	if d := m.RecursionDepth(); d != 0 {
		t.Fatalf("depth = %d after exit, want 0", d)
	}
}

func TestMonitorReentrant(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(99)

	m.Enter(ec)
	m.Enter(ec)
	m.Enter(ec)
	if d := m.RecursionDepth(); d != 3 {
		t.Fatalf("depth = %d, want 3", d)
	}
	if !m.HoldsLock(ec) {
		t.Fatal("HoldsLock false at depth 3")
	}
	m.Exit(ec)
	if d := m.RecursionDepth(); d != 2 {
		t.Fatalf("depth = %d after one exit, want 2", d)
	}
	m.Exit(ec)
	m.Exit(ec)
	if m.HoldsLock(ec) {
		t.Fatal("HoldsLock true after all exits")
	}
}

func TestMonitorNonOwnerExit(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()

	m.Enter(1)
	if m.Exit(2) {
		t.Fatal("Exit as non-owner should return false")
	}
	// Owner should still hold the lock.
	if !m.HoldsLock(1) {
		t.Fatal("original owner lost lock after non-owner Exit")
	}
	m.Exit(1)
}

// --- Exclusion ---

func TestMonitorExclusion(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()

	var entered1, entered2 sync.WaitGroup
	entered1.Add(1)
	entered2.Add(1)
	var done sync.WaitGroup
	done.Add(2)

	// EC 1 enters and holds the lock.
	m.Enter(1)

	// EC 2 tries to enter from a goroutine — must block.
	go func() {
		entered2.Done() // signal: we're about to block
		m.Enter(2)
		if !m.HoldsLock(2) {
			t.Error("EC 2 entered but doesn't hold lock")
		}
		m.Exit(2)
		done.Done()
	}()

	// EC 1 reentrant entry from another goroutine (same EC, different goroutine).
	go func() {
		m.Enter(1) // reentrant — should not block (main still holds lock)
		entered1.Done() // signal AFTER entering so depth ≥ 2 is stable
		if d := m.RecursionDepth(); d < 2 {
			t.Errorf("reentrant depth = %d, want >= 2", d)
		}
		m.Exit(1)
		done.Done()
	}()

	// Wait for both goroutines to signal they've started.
	entered1.Wait()
	entered2.Wait()

	// Both should still be blocked (or about to block) — the lock is held by EC 1 at depth 1.
	// Release.
	m.Exit(1)

	done.Wait()

	if m.HoldsLock(1) {
		t.Fatal("EC 1 still holds after full release")
	}
	if m.HoldsLock(2) {
		t.Fatal("EC 2 still holds after all exits")
	}
}

// --- Wait set ---

func TestMonitorWaitNotify(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(5)

	// Main goroutine holds the monitor at depth 2.
	m.Enter(ec)
	m.Enter(ec)

	var normal bool
	var waiterReady sync.WaitGroup
	waiterReady.Add(1)
	var waiterDone sync.WaitGroup
	waiterDone.Add(1)

	go func() {
		defer waiterDone.Done()
		// Reentrant enter — same EC.
		m.Enter(ec) // depth → 3
		waiterReady.Done()
		normal = m.InternalWait(ec, 3)
		// After wake, monitor is reacquired and depth restored.
		m.Exit(ec) // release the extra Enter
	}()

	// Wait for the waiter to be enqueued and monitor released.
	waiterReady.Wait()
	for m.WaitSetLen() == 0 {
	}

	// Monitor was fully released by InternalWait. Re-enter.
	m.Enter(ec) // blocks until free — should be immediate

	if !m.Notify(ec) {
		t.Fatal("Notify returned false for owner")
	}

	// Release so the waiter can reacquire.
	m.Exit(ec)
	m.Exit(ec) // depth back to 1
	m.Exit(ec) // fully release for waiter to reacquire

	waiterDone.Wait()

	if !normal {
		t.Fatal("InternalWait returned false (interrupted) after notify")
	}
	if !m.HoldsLock(ec) {
		t.Fatal("monitor not reacquired after wait")
	}
	m.Exit(ec)
}

func TestMonitorWaitNotifyAll(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()

	// Distinct execution contexts simulate distinct Java threads.
	mainEC := uint64(101)
	ecs := []uint64{102, 103, 104}

	// Main holds the monitor.
	m.Enter(mainEC)

	results := make([]bool, 3)
	var allReady sync.WaitGroup
	var allDone sync.WaitGroup

	for i := 0; i < 3; i++ {
		allReady.Add(1)
		allDone.Add(1)
		go func(idx int) {
			defer allDone.Done()
			m.Enter(ecs[idx])
			allReady.Done()
			results[idx] = m.InternalWait(ecs[idx], 1)
			m.Exit(ecs[idx]) // release after InternalWait reacquires
		}(i)
	}

	// Release so waiters can enter one at a time (each releases after InternalWait).
	m.Exit(mainEC)

	// Wait for all three to be enqueued.
	allReady.Wait()
	for m.WaitSetLen() < 3 {
	}

	// Re-enter as mainEC and notify all.
	m.Enter(mainEC)
	if !m.NotifyAll(mainEC) {
		t.Fatal("NotifyAll returned false for owner")
	}
	m.Exit(mainEC)

	allDone.Wait()

	for i, r := range results {
		if !r {
			t.Errorf("waiter %d got interrupted result after NotifyAll", i)
		}
	}

	// Each waiter reacquired — clean up any remaining holders.
	// Waiters exit after InternalWait returns — they exit the monitor before finishing.
	for i := 0; i < 3; i++ {
		if m.HoldsLock(ecs[i]) {
			m.Exit(ecs[i])
		}
	}
	if m.HoldsLock(mainEC) {
		m.Exit(mainEC)
	}
}

func TestMonitorNotifyNoOwner(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()

	// notify without ownership should return false.
	if m.Notify(99) {
		t.Fatal("Notify as non-owner should return false")
	}

	m.Enter(1)
	if m.Notify(2) {
		t.Fatal("Notify as wrong owner should return false")
	}
	m.Exit(1)
}

// --- Interrupt ordering (ADR-0029) ---

func TestMonitorInterruptBeforeNotify(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(10)

	// Main holds monitor.
	m.Enter(ec)

	var result bool
	var ready sync.WaitGroup
	ready.Add(1)
	var done sync.WaitGroup
	done.Add(1)

	go func() {
		defer done.Done()
		m.Enter(ec) // reentrant
		ready.Done()
		result = m.InternalWait(ec, 2)
		m.Exit(ec) // release the extra Enter
	}()

	ready.Wait()
	for m.WaitSetLen() == 0 {
	}

	// Waiter released. InterruptWaiter does NOT require ownership — it only
	// locks the monitor's state mutex. Interrupt wins the race vs notify.
	m.InterruptWaiter(ec)

	done.Wait()

	if result {
		t.Fatal("InternalWait returned true after interrupt won (should be false)")
	}

	// Clean up.
	for m.HoldsLock(ec) {
		m.Exit(ec)
	}
}

func TestMonitorNotifyBeforeInterrupt(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(11)

	m.Enter(ec)

	var normal bool
	var ready sync.WaitGroup
	ready.Add(1)
	var done sync.WaitGroup
	done.Add(1)

	go func() {
		defer done.Done()
		m.Enter(ec) // reentrant
		ready.Done()
		normal = m.InternalWait(ec, 2)
		m.Exit(ec) // release the extra Enter
	}()

	ready.Wait()
	for m.WaitSetLen() == 0 {
	}

	// Waiter released. Re-enter, then Notify (must own monitor).
	m.Enter(ec)
	if !m.Notify(ec) {
		t.Fatal("Notify returned false for owner")
	}
	m.Exit(ec) // release for waiter to reacquire

	done.Wait()

	if !normal {
		t.Fatal("InternalWait returned false after notify won (should be true)")
	}

	// Interrupt after notify — waiter already notified, InterruptWaiter is a no-op.
	m.InterruptWaiter(ec)

	// Clean up.
	for m.HoldsLock(ec) {
		m.Exit(ec)
	}
}

// --- Clone isolation ---

func TestMonitorCloneIsolation(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m1 := obj.Monitor()

	m1.Enter(1)

	clone := CloneObject(obj)
	m2 := clone.Monitor()

	if m1 == m2 {
		t.Fatal("clone shares the same Monitor pointer")
	}
	if m2.HoldsLock(1) {
		t.Fatal("clone monitor should not be held by original owner")
	}
	if d := m2.RecursionDepth(); d != 0 {
		t.Fatalf("clone monitor depth = %d, want 0", d)
	}

	m1.Exit(1)
}

// --- Monitor on arrays and Class mirrors ---

func TestMonitorOnArray(t *testing.T) {
	cls := newMinimalArrayClass("[I")
	obj := NewArray(cls, 10)
	m := obj.Monitor()

	m.Enter(42)
	if !m.HoldsLock(42) {
		t.Fatal("array monitor not held after Enter")
	}
	m.Exit(42)
}

func TestMonitorOnClassMirror(t *testing.T) {
	cls := newMinimalClass("Test")
	mirror := cls.ClassObject(func() *Object { return NewObject(newMinimalClass("java/lang/Class")) })
	m := mirror.Monitor()

	m.Enter(77)
	if !m.HoldsLock(77) {
		t.Fatal("Class mirror monitor not held after Enter")
	}
	m.Exit(77)
}

// --- Concrete depth test (exact restore) ---

func TestMonitorWaitExactDepthRestore(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(55)

	// Enter to depth 5.
	m.Enter(ec) // 1
	m.Enter(ec) // 2
	m.Enter(ec) // 3
	m.Enter(ec) // 4
	m.Enter(ec) // 5

	var finalDepth int
	var ready sync.WaitGroup
	ready.Add(1)
	var done sync.WaitGroup
	done.Add(1)

	go func() {
		defer done.Done()
		m.Enter(ec) // reentrant → 6
		ready.Done()
		m.InternalWait(ec, 6)
		m.mu.Lock()
		finalDepth = m.depth
		m.mu.Unlock()
		m.Exit(ec) // release the extra Enter
	}()

	ready.Wait()
	for m.WaitSetLen() == 0 {
	}

	// The monitor should be released now (waiter fully released it).
	if m.HoldsLock(ec) {
		t.Fatal("monitor not released during wait")
	}

	// Re-enter and notify.
	m.Enter(ec)
	m.Notify(ec)
	m.Exit(ec) // release for waiter to reacquire

	done.Wait()

	if finalDepth != 6 {
		t.Fatalf("restored depth = %d, want 6", finalDepth)
	}
	if !m.HoldsLock(ec) {
		t.Fatal("monitor not reacquired after wait")
	}

	// The waiter goroutine did m.Exit(ec) once (for the extra Enter).
	// The restored depth is 6, minus the exit = 5.
	if d := m.RecursionDepth(); d != 5 {
		t.Fatalf("depth after wait = %d, want 5", d)
	}

	// Exit remaining levels.
	for i := 0; i < 5; i++ {
		m.Exit(ec)
	}
	if m.HoldsLock(ec) {
		t.Fatal("monitor still held after all exits")
	}
}

// --- Concurrent CAS: many goroutines racing on the same object ---

func TestMonitorConcurrentCASSameIdentity(t *testing.T) {
	obj := NewObject(newMinimalClass("Race"))
	var monitors [50]*Monitor
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			monitors[idx] = obj.Monitor()
		}(i)
	}
	wg.Wait()
	for i := 1; i < 50; i++ {
		if monitors[i] != monitors[0] {
			t.Fatalf("concurrent Monitor() returned different pointers at index %d", i)
		}
	}
}

// --- helpers ---

// newMinimalClass builds a bare synthetic class with a name for testing.
// It won't have a real loader or superclass chain; only the name matters.
func newMinimalClass(name string) *Class {
	c := &Class{
		name:        name,
		methodTable: make(map[string]*Method),
	}
	c.initCond = sync.NewCond(&c.initMu)
	return c
}

// newMinimalArrayClass builds a bare array class for testing.
func newMinimalArrayClass(name string) *Class {
	c := &Class{
		name:    name,
		isArray: true,
	}
	c.initCond = sync.NewCond(&c.initMu)
	return c
}

// TestMonitorSecondWaitInterruptible verifies that after a first wait/notify
// cycle completes, a second wait() on the same monitor by the same execution
// context produces a fresh, correctly interruptible entry. This is the
// regression test for the stale-entries bug: completed waiters were not
// removed from waitSet, and InterruptWaiter returned on the first EC match
// (the stale, now-notified entry) instead of finding the live waiter.
func TestMonitorSecondWaitInterruptible(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(7)

	// === First wait/notify cycle ===
	m.Enter(ec) // depth 1

	var firstResult bool
	var waiterReady sync.WaitGroup
	waiterReady.Add(1)
	var waiterDone sync.WaitGroup
	waiterDone.Add(1)

	go func() {
		defer waiterDone.Done()
		m.Enter(ec) // reentrant → depth 2
		waiterReady.Done()
		firstResult = m.InternalWait(ec, 2)
		m.Exit(ec) // release the extra Enter
	}()

	waiterReady.Wait()
	for m.WaitSetLen() < 1 {
	}

	// Notify the first waiter.
	m.Enter(ec) // depth 2
	m.Notify(ec)
	m.Exit(ec) // depth 1 — release the extra Enter

	// Release so waiter can reacquire.
	m.Exit(ec) // depth 0 — fully release
	waiterDone.Wait()

	if !firstResult {
		t.Fatal("first InternalWait returned false (interrupted), want true (notified)")
	}

	// After the first cycle, the wait set must be empty (stale entry removed).
	if n := m.WaitSetLen(); n != 0 {
		t.Fatalf("wait set len after first cycle = %d, want 0", n)
	}

	// Reacquire the monitor.
	m.Enter(ec) // depth 1

	// === Second wait — this one will be interrupted ===
	var secondResult bool
	waiterReady.Add(1)
	waiterDone.Add(1)

	go func() {
		defer waiterDone.Done()
		m.Enter(ec) // reentrant → depth 2
		waiterReady.Done()
		secondResult = m.InternalWait(ec, 2)
		m.Exit(ec) // release the extra Enter
	}()

	waiterReady.Wait()
	for m.WaitSetLen() < 1 {
	}

	// Interrupt (instead of notify) — must find the LIVE entry, not a stale one.
	m.InterruptWaiter(ec)

	// Release so waiter can reacquire.
	m.Exit(ec) // depth 0
	waiterDone.Wait()

	if secondResult {
		t.Fatal("second InternalWait returned true (notified), want false (interrupted)")
	}

	// Verify cleanup after the second cycle.
	if n := m.WaitSetLen(); n != 0 {
		t.Fatalf("wait set len after second cycle = %d, want 0", n)
	}

	// Clean up any remaining depth.
	for m.HoldsLock(ec) {
		m.Exit(ec)
	}
}

// TestMonitorWaitSetCleanup verifies that completed waiters do not accumulate
// in the wait set, preventing unbounded growth across repeated wait/notify
// cycles.
func TestMonitorWaitSetCleanup(t *testing.T) {
	obj := NewObject(newMinimalClass("Test"))
	m := obj.Monitor()
	ec := uint64(11)

	const rounds = 5
	for r := 0; r < rounds; r++ {
		if n := m.WaitSetLen(); n != 0 {
			t.Fatalf("round %d: wait set len before wait = %d, want 0", r, n)
		}

		m.Enter(ec) // depth 1

		var ready sync.WaitGroup
		ready.Add(1)
		var done sync.WaitGroup
		done.Add(1)

		go func() {
			defer done.Done()
			m.Enter(ec) // depth 2
			ready.Done()
			m.InternalWait(ec, 2)
			m.Exit(ec)
		}()

		ready.Wait()
		for m.WaitSetLen() < 1 {
		}

		// Notify.
		m.Enter(ec) // depth 2
		m.Notify(ec)
		m.Exit(ec) // depth 1

		m.Exit(ec) // depth 0 — release for waiter
		done.Wait()

		// After each round, the wait set must be empty.
		if n := m.WaitSetLen(); n != 0 {
			t.Fatalf("round %d: wait set len after cycle = %d, want 0", r, n)
		}

		// Clean up any remaining ownership.
		for m.HoldsLock(ec) {
			m.Exit(ec)
		}
	}
}
