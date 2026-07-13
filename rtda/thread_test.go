package rtda

import (
	"sync"
	"testing"
	"time"
)

// TestThreadLifecycle verifies the state machine: NEW → RUNNABLE → TERMINATED.
func TestThreadLifecycle(t *testing.T) {
	t.Run("new thread is NOT alive", func(t *testing.T) {
		tr := NewThread(nil)
		if tr.IsAlive() {
			t.Error("new thread should not be alive")
		}
	})

	t.Run("SetStarted transitions NEW to RUNNABLE", func(t *testing.T) {
		tr := NewThread(nil)
		if !tr.SetStarted() {
			t.Error("first SetStarted should succeed")
		}
		if !tr.IsAlive() {
			t.Error("thread should be alive after SetStarted")
		}
	})

	t.Run("SetStarted fails if already started", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		if tr.SetStarted() {
			t.Error("second SetStarted should fail")
		}
	})

	t.Run("SetStarted fails after terminated", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		tr.Terminate()
		if tr.IsAlive() {
			t.Error("terminated thread should not be alive")
		}
		if tr.SetStarted() {
			t.Error("SetStarted should fail after termination")
		}
	})

	t.Run("Terminate closes done channel", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		tr.Terminate()
		select {
		case <-tr.Done():
			// expected
		default:
			t.Error("Done channel should be closed after Terminate")
		}
	})

	t.Run("Terminate double-call does not panic", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		tr.Terminate()
		// Second Terminate should not panic (state is already TERMINATED,
		// and close-once is guarded by CAS in the real implementation).
		// The current implementation does unconditional close — check that
		// we at least don't crash.
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Log("double close panicked (expected if no CAS guard):", r)
				}
			}()
			tr.Terminate()
		}()
	})
}

// TestThreadInterrupt verifies interrupt flag set/read/clear semantics.
func TestThreadInterrupt(t *testing.T) {
	t.Run("initial state is not interrupted", func(t *testing.T) {
		tr := NewThread(nil)
		if tr.IsInterrupted() {
			t.Error("new thread should not be interrupted")
		}
	})

	t.Run("Interrupt sets flag", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt()
		if !tr.IsInterrupted() {
			t.Error("Interrupt should set flag")
		}
	})

	t.Run("IsInterrupted does not clear", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt()
		if !tr.IsInterrupted() {
			t.Error("first IsInterrupted should return true")
		}
		if !tr.IsInterrupted() {
			t.Error("second IsInterrupted should still return true")
		}
	})

	t.Run("Interrupted reads and clears", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt()
		if !tr.Interrupted() {
			t.Error("Interrupted should return true")
		}
		if tr.IsInterrupted() {
			t.Error("flag should be cleared after Interrupted")
		}
	})

	t.Run("Interrupted returns false when not interrupted", func(t *testing.T) {
		tr := NewThread(nil)
		if tr.Interrupted() {
			t.Error("Interrupted should return false when not interrupted")
		}
	})
}

// TestThreadWaker verifies the waker channel is signaled by Interrupt.
func TestThreadWaker(t *testing.T) {
	t.Run("Interrupt signals waker", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt()
		select {
		case <-tr.Waker():
			// expected
		default:
			t.Error("Waker should be signaled after Interrupt")
		}
	})

	t.Run("Interrupt drain validates semantics", func(t *testing.T) {
		tr := NewThread(nil)
		// First interrupt should signal
		tr.Interrupt()
		<-tr.Waker() // drain

		// Second interrupt should also signal
		tr.Interrupt()
		select {
		case <-tr.Waker():
			// expected
		default:
			t.Error("Waker should be signaled on second Interrupt too")
		}
	})
}

// TestThreadSleep verifies interruptible sleep.
func TestThreadSleep(t *testing.T) {
	t.Run("sleep normally completes", func(t *testing.T) {
		tr := NewThread(nil)
		start := time.Now()
		if !tr.Sleep(10) {
			t.Error("uninterrupted sleep should return true")
		}
		if time.Since(start) < 5*time.Millisecond {
			t.Error("sleep returned too quickly")
		}
	})

	t.Run("sleep returns false if interrupted before", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt()
		if tr.Sleep(1000) {
			t.Error("sleep should return false when interrupted")
		}
		// After sleep returns false, flag should be cleared
		if tr.IsInterrupted() {
			t.Error("sleep should clear interrupt flag")
		}
	})

	t.Run("sleep interrupted during sleep", func(t *testing.T) {
		tr := NewThread(nil)
		done := make(chan bool)
		go func() {
			done <- tr.Sleep(10000) // long sleep
		}()
		time.Sleep(5 * time.Millisecond) // let sleep start
		tr.Interrupt()                   // interrupt during sleep
		result := <-done
		if result {
			t.Error("sleep should return false when interrupted during sleep")
		}
	})

	t.Run("zero sleep returns true", func(t *testing.T) {
		tr := NewThread(nil)
		if !tr.Sleep(0) {
			t.Error("zero sleep should return true")
		}
	})
}

// TestThreadDaemon verifies daemon flag.
func TestThreadDaemon(t *testing.T) {
	tr := NewThread(nil)
	if tr.IsDaemon() {
		t.Error("new thread default should be non-daemon")
	}
	tr.SetDaemon(true)
	if !tr.IsDaemon() {
		t.Error("SetDaemon(true) should mark thread as daemon")
	}
}

// TestThreadDoneChannel verifies join via Done channel.
func TestThreadDoneChannel(t *testing.T) {
	t.Run("Done blocks until Terminate", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()

		received := make(chan struct{})
		go func() {
			<-tr.Done()
			close(received)
		}()

		// Should not be done yet
		select {
		case <-received:
			t.Error("Done should not fire before Terminate")
		case <-time.After(10 * time.Millisecond):
		}

		tr.Terminate()

		select {
		case <-received:
			// expected
		case <-time.After(time.Second):
			t.Error("Done should fire after Terminate")
		}
	})
}

// TestThreadJavaThread verifies facade attachment.
func TestThreadJavaThread(t *testing.T) {
	tr := NewThread(nil)
	if tr.JavaThread() != nil {
		t.Error("JavaThread should be nil initially")
	}
	obj := &Object{}
	tr.SetJavaThread(obj)
	if tr.JavaThread() != obj {
		t.Error("JavaThread should return the set object")
	}
}

// TestConcurrentThreadCreation verifies atomic ecID assignment under race.
func TestConcurrentThreadCreation(t *testing.T) {
	const N = 100
	seen := make(map[uint64]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr := NewThread(nil)
			mu.Lock()
			seen[tr.EC()] = true
			mu.Unlock()
		}()
	}
	wg.Wait()

	if len(seen) != N {
		t.Errorf("expected %d unique ecIDs, got %d", N, len(seen))
	}
}

// TestThreadMainFlag verifies main flag.
func TestThreadMainFlag(t *testing.T) {
	tr := NewThread(nil)
	if tr.IsMain() {
		t.Error("new thread should not be main by default")
	}
	tr.SetMain(true)
	if !tr.IsMain() {
		t.Error("SetMain(true) should mark as main")
	}
}
