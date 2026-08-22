package kernel

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestRegistryInterruptSleep(t *testing.T) {
	r := NewThreadRegistry()
	obj := &Instance{} // stand-in object; registry only stores the pointer
	j := r.Register(1, obj, "sleeper")

	done := make(chan bool, 1) // completed=false ⇒ interrupted
	go func() {
		done <- r.Sleep(1, 5*time.Second)
	}()

	deadline := time.Now().Add(time.Second)
	for j.SleepersNow() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("sleeper never registered")
		}
		time.Sleep(time.Millisecond)
	}
	if !r.Interrupt(1) {
		t.Fatal("Interrupt must report it woke a sleeper")
	}
	if completed := <-done; completed {
		t.Fatal("sleep must report interrupted")
	}
	if !j.interrupted.Load() {
		t.Fatal("flag set by Interrupt was consumed — Sleep contract clears it")
	}
}

func TestRegistryInterruptWaiter(t *testing.T) {
	r := NewThreadRegistry()
	m := &Monitor{}
	j := r.Register(2, &Instance{}, "waiter")
	_ = j

	outCh := make(chan waitOutcome, 1)
	go func() {
		m.Enter(2)
		outCh <- m.Wait(r, 2, -1)
	}()

	waiters := func() int {
		r.mu.Lock()
		defer r.mu.Unlock()
		return len(r.waiting[2])
	}
	dl := time.Now().Add(time.Second)
	for waiters() == 0 {
		if time.Now().After(dl) {
			t.Fatal("Object.wait never published to registry")
		}
		time.Sleep(time.Millisecond)
	}
	if !r.Interrupt(2) {
		t.Fatal("Interrupt must wake monitor waiters")
	}
	if out := <-outCh; out != waitGotInterrupt {
		t.Fatalf("outcome=%v", out)
	}
	m.Exit(2)
}

func TestJoinLatchAndStates(t *testing.T) {
	r := NewThreadRegistry()
	obj := &Instance{}
	j := r.Register(3, obj, "t3")
	if j.IsAlive() {
		t.Fatal("fresh record not alive yet")
	}
	j.alive.Store(true)
	j.state.Store(int32(threadRunning))

	reached := make(chan bool, 1)
	go func() {
		ok, _ := r.Join(0, j, time.Second)
		reached <- ok
	}()
	time.Sleep(20 * time.Millisecond)
	select {
	case <-reached:
		t.Fatal("join finished before terminate")
	default:
	}
	j.Terminate()
	if ok := <-reached; !ok {
		t.Fatal("join must complete after Terminate")
	}
	if j.IsAlive() || j.StateName() != "TERMINATED" {
		t.Fatalf("state after termination: alive=%v %s", j.IsAlive(), j.StateName())
	}
	// double terminate is a no-op
	j.Terminate()
}

func TestInterruptFlagPeekClear(t *testing.T) {
	r := NewThreadRegistry()
	r.Register(4, &Instance{}, "t4")
	if r.PeekInterrupted(4) {
		t.Fatal("fresh flag set")
	}
	r.Interrupt(4)
	if !r.PeekInterrupted(4) {
		t.Fatal("interrupt must set flag")
	}
	if !r.ClearInterrupted(4) {
		t.Fatal("clear returns previous")
	}
	if r.PeekInterrupted(4) {
		t.Fatal("cleared flag persists")
	}
}

// TestConcurrentCounters exercises registry under contention (race detector).
func TestConcurrentCounters(t *testing.T) {
	r := NewThreadRegistry()
	var live atomic.Int32
	for key := uint64(100); key < 108; key++ {
		r.Register(key, &Instance{}, "c")
		live.Add(1)
		go func(k uint64) {
			defer live.Add(-1)
			_ = r.Sleep(k, time.Millisecond)
		}(key)
	}
	dl := time.Now().Add(time.Second)
	for live.Load() > 0 && time.Now().Before(dl) {
		time.Sleep(time.Millisecond)
	}
	if live.Load() != 0 {
		t.Fatal("sleepers did not finish")
	}
}
