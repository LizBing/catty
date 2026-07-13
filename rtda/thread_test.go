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

	t.Run("Terminate double-call is safe", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		tr.Terminate()
		// Second Terminate must be safe — CAS prevents double-close.
		// No panic allowed; test fails if Terminate panics.
		tr.Terminate()
		// State must still be TERMINATED.
		if tr.IsAlive() {
			t.Error("thread should not be alive after double Terminate")
		}
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

// --- Blocker 2: Interrupted() waker drain + Sleep interaction ---

// TestInterruptedDrainsWaker verifies that Interrupted() drains the stale
// waker signal so subsequent Sleep returns normally.
func TestInterruptedDrainsWaker(t *testing.T) {
	t.Run("sleep after interrupted clears normally", func(t *testing.T) {
		tr := NewThread(nil)
		tr.Interrupt() // sets flag + sends to waker

		// interrupted() should return true and drain the waker.
		if !tr.Interrupted() {
			t.Error("Interrupted should return true after Interrupt")
		}
		// Flag should be cleared.
		if tr.IsInterrupted() {
			t.Error("flag should be cleared after Interrupted")
		}

		// Subsequent sleep must complete normally — no stale waker signal.
		if !tr.Sleep(50) {
			t.Error("sleep after interrupted() should complete normally, not report interrupted")
		}
	})

	t.Run("sleep interrupted during sleep throws", func(t *testing.T) {
		tr := NewThread(nil)
		done := make(chan bool)
		go func() {
			done <- tr.Sleep(60000) // long sleep
		}()
		time.Sleep(10 * time.Millisecond) // let sleep start
		tr.Interrupt()
		result := <-done
		if result {
			t.Error("sleep should return false (interrupted) when interrupt arrives during sleep")
		}
		if tr.IsInterrupted() {
			t.Error("sleep should clear interrupt flag via Interrupted re-check")
		}
	})
}

// TestInterruptBoundaryCases verifies multiple interrupt/clear/sleep sequences.
func TestInterruptBoundaryCases(t *testing.T) {
	t.Run("multiple interrupt then clear all", func(t *testing.T) {
		tr := NewThread(nil)

		// First interrupt.
		tr.Interrupt()
		if !tr.Interrupted() {
			t.Error("first Interrupted should return true")
		}

		// Second interrupt.
		tr.Interrupt()
		if !tr.IsInterrupted() {
			t.Error("second Interrupt should set flag")
		}

		// Sleep should see it.
		if tr.Sleep(1000) {
			t.Error("sleep should return false when interrupted")
		}

		// State is now clear.
		if tr.IsInterrupted() {
			t.Error("flag should be clear after sleep interrupted check")
		}

		// Another sleep should complete normally.
		if !tr.Sleep(10) {
			t.Error("sleep after clear should complete normally")
		}
	})

	t.Run("interrupt check clear check sleep", func(t *testing.T) {
		tr := NewThread(nil)

		tr.Interrupt()
		if !tr.IsInterrupted() {
			t.Error("IsInterrupted should return true after Interrupt")
		}
		if !tr.Interrupted() {
			t.Error("Interrupted should return true")
		}
		if tr.IsInterrupted() {
			t.Error("IsInterrupted should return false after Interrupted")
		}
		// Sleep must complete normally.
		if !tr.Sleep(20) {
			t.Error("sleep must complete normally after interrupt is cleared")
		}
	})

	t.Run("interrupt during sleep preserves state before clear", func(t *testing.T) {
		tr := NewThread(nil)

		// Interrupt, then sleep — should return false immediately.
		tr.Interrupt()
		if tr.Sleep(1000) {
			t.Error("sleep should return false when pre-interrupted")
		}
		// Flag cleared by sleep's Interrupted check.
		if tr.IsInterrupted() {
			t.Error("flag should be clear after sleep detects interrupt")
		}
	})

	t.Run("concurrent interrupt and interrupted drain", func(t *testing.T) {
		tr := NewThread(nil)
		const N = 50

		var wg sync.WaitGroup
		for i := 0; i < N; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				tr.Interrupt()
			}()
			wg.Add(1)
			go func() {
				defer wg.Done()
				tr.Interrupted() // must be race-free
			}()
		}
		wg.Wait()

		// After draining, flag may be set (if last op was Interrupt) or clear
		// (if last op was Interrupted). Both are valid. Just verify no race.
	})
}

// --- Blocker 3: SetDaemon lifecycle rules ---

// TestSetDaemonLifecycle verifies that setDaemon only succeeds before start.
func TestSetDaemonLifecycle(t *testing.T) {
	t.Run("setDaemon ok when NEW", func(t *testing.T) {
		tr := NewThread(nil)
		if !tr.SetDaemon(true) {
			t.Error("SetDaemon should succeed when thread is NEW")
		}
		if !tr.IsDaemon() {
			t.Error("IsDaemon should return true after SetDaemon(true)")
		}
		// Change back.
		if !tr.SetDaemon(false) {
			t.Error("SetDaemon should succeed again when thread is still NEW")
		}
		if tr.IsDaemon() {
			t.Error("IsDaemon should return false after SetDaemon(false)")
		}
	})

	t.Run("setDaemon fails after start", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		if tr.SetDaemon(true) {
			t.Error("SetDaemon should fail after start")
		}
		// Default is non-daemon — must not have changed.
		if tr.IsDaemon() {
			t.Error("daemon should remain false after failed SetDaemon")
		}
	})

	t.Run("setDaemon fails after terminate", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetStarted()
		tr.Terminate()
		if tr.SetDaemon(true) {
			t.Error("SetDaemon should fail after termination")
		}
	})
}

// TestConcurrentSetDaemonAndStart verifies no race between setDaemon and start,
// and that the outcome reflects exactly one allowed order.
func TestConcurrentSetDaemonAndStart(t *testing.T) {
	// Run many iterations to exercise the race.
	for i := 0; i < 100; i++ {
		tr := NewThread(nil)

		var setResult bool
		var startResult bool
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			setResult = tr.SetDaemon(true)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			startResult = tr.SetStarted()
		}()

		wg.Wait()

		// If start won, SetDaemon must have failed.
		// If SetDaemon won, the daemon value must be true.
		if startResult && setResult {
			// Both succeeded → SetDaemon happened before the CAS in SetStarted.
			// Daemon must be true.
			if tr.ConsumeDaemonForStart() != true {
				t.Errorf("iter %d: daemon must be true when SetDaemon won race with start", i)
			}
		} else if startResult && !setResult {
			// Start won, SetDaemon lost → SetDaemon saw state != NEW.
			// Daemon should remain false (default).
			if tr.ConsumeDaemonForStart() != false {
				t.Errorf("iter %d: daemon must be false when SetDaemon lost race with start", i)
			}
		}
		// If !startResult, SetDaemon should succeed (no concurrent start).
		if !startResult && !setResult {
			t.Errorf("iter %d: SetDaemon should succeed when start hasn't happened", i)
		}
	}
}

// TestConcurrentSetDaemonAndIsDaemon verifies that concurrent SetDaemon and
// IsDaemon on an unstarted thread are race-free. After start, daemon is
// frozen — subsequent SetDaemon calls fail.
func TestConcurrentSetDaemonAndIsDaemon(t *testing.T) {
	t.Run("concurrent read write on NEW thread is race-free", func(t *testing.T) {
		tr := NewThread(nil)
		const N = 200

		var wg sync.WaitGroup
		// Writers: toggle daemon flag.
		for i := 0; i < N/2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					tr.SetDaemon(j%2 == 0)
				}
			}()
		}
		// Readers: call IsDaemon.
		for i := 0; i < N/2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					_ = tr.IsDaemon()
				}
			}()
		}
		wg.Wait()
		// No Go race detector warning is the assertion.
	})

	t.Run("daemon frozen after start", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetDaemon(true)
		if !tr.IsDaemon() {
			t.Fatal("SetDaemon(true) should work on NEW thread")
		}
		tr.SetStarted()
		// After start, daemon is frozen.
		if !tr.IsDaemon() {
			t.Error("daemon should still be true after start")
		}
		// SetDaemon must fail.
		if tr.SetDaemon(false) {
			t.Error("SetDaemon should fail after start")
		}
		// Value must be unchanged.
		if !tr.IsDaemon() {
			t.Error("daemon must remain true after failed SetDaemon")
		}
	})

	t.Run("daemon frozen after terminate", func(t *testing.T) {
		tr := NewThread(nil)
		tr.SetDaemon(true)
		tr.SetStarted()
		tr.Terminate()
		if tr.SetDaemon(false) {
			t.Error("SetDaemon should fail after termination")
		}
		if !tr.IsDaemon() {
			t.Error("daemon must remain true after termination")
		}
	})
}
