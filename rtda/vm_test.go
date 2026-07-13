package rtda

import (
	"sync"
	"testing"
	"time"
)

// TestVMNonDaemonCount verifies increment/decrement of non-daemon count.
func TestVMNonDaemonCount(t *testing.T) {
	vm := NewVM()

	// Start 3 non-daemon threads.
	vm.ThreadStarted(false)
	vm.ThreadStarted(false)
	vm.ThreadStarted(false)

	// Terminate 2.
	vm.ThreadTerminated(false)
	vm.ThreadTerminated(false)

	// 1 remains — WaitForNonDaemon should NOT unblock.
	done := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(done)
	}()

	select {
	case <-done:
		t.Error("WaitForNonDaemonThreads unblocked while non-daemon count > 0")
	case <-time.After(30 * time.Millisecond):
	}

	// Terminate the last.
	vm.ThreadTerminated(false)

	select {
	case <-done:
		// expected
	case <-time.After(time.Second):
		t.Error("WaitForNonDaemonThreads did not unblock after count reached 0")
	}
}

// TestVMDaemonDoesNotAffectCount verifies daemon threads don't change non-daemon count.
func TestVMDaemonDoesNotAffectCount(t *testing.T) {
	vm := NewVM()

	// Start a non-daemon thread.
	vm.ThreadStarted(false)

	// Daemon threads should not affect the count.
	vm.ThreadStarted(true)
	vm.ThreadStarted(true)
	vm.ThreadTerminated(true)

	// Non-daemon count is still 1 — should not unblock.
	done := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(done)
	}()

	select {
	case <-done:
		t.Error("WaitForNonDaemonThreads unblocked with daemon threads only")
	case <-time.After(30 * time.Millisecond):
	}

	// Terminate the non-daemon.
	vm.ThreadTerminated(false)

	select {
	case <-done:
		// expected
	case <-time.After(time.Second):
		t.Error("WaitForNonDaemonThreads did not unblock")
	}
}

// TestVMZeroNonDaemonImmediateUnblock verifies WaitForNonDaemonThreads returns
// immediately when non-daemon count is already zero.
func TestVMZeroNonDaemonImmediateUnblock(t *testing.T) {
	vm := NewVM()

	done := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(done)
	}()

	select {
	case <-done:
		// expected
	case <-time.After(time.Second):
		t.Error("WaitForNonDaemonThreads blocked when count is zero")
	}
}

// TestVMConcurrentStartTerminate stress-tests concurrent ThreadStarted/Terminated.
func TestVMConcurrentStartTerminate(t *testing.T) {
	vm := NewVM()
	const N = 100

	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			vm.ThreadStarted(false)
			// Simulate some work.
			vm.ThreadTerminated(false)
		}()
	}
	wg.Wait()

	// All non-daemon threads have terminated — should not block.
	done := make(chan struct{})
	go func() {
		vm.WaitForNonDaemonThreads()
		close(done)
	}()

	select {
	case <-done:
		// expected
	case <-time.After(time.Second):
		t.Error("WaitForNonDaemonThreads blocked after all threads terminated")
	}
}

// TestVMGetSetVM verifies SetVM/GetVM round-trip.
func TestVMGetSetVM(t *testing.T) {
	vm1 := NewVM()
	SetVM(vm1)
	if GetVM() != vm1 {
		t.Error("GetVM returned different VM than SetVM")
	}

	vm2 := NewVM()
	SetVM(vm2)
	if GetVM() != vm2 {
		t.Error("GetVM after second SetVM returned wrong VM")
	}
}
