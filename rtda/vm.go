package rtda

import "sync"

// VM supervises thread liveness per ADR-0028: started non-daemon threads keep
// the launcher alive; daemon threads do not. The launcher creates a VM before
// running main and calls WaitForNonDaemonThreads after the main thread returns.
type VM struct {
	mu        sync.Mutex
	cond      *sync.Cond
	nonDaemon int
}

// NewVM creates a new VM supervisor. Exactly one should be created per process
// and registered via SetVM before the main thread starts.
func NewVM() *VM {
	v := &VM{}
	v.cond = sync.NewCond(&v.mu)
	return v
}

// ThreadStarted records that a thread has started. Daemon threads are ignored;
// non-daemon threads increment the liveness counter.
func (v *VM) ThreadStarted(daemon bool) {
	if daemon {
		return
	}
	v.mu.Lock()
	v.nonDaemon++
	v.mu.Unlock()
}

// ThreadTerminated records that a thread has finished. Daemon threads are
// ignored; non-daemon threads decrement the liveness counter and broadcast when
// it reaches zero so WaitForNonDaemonThreads unblocks.
func (v *VM) ThreadTerminated(daemon bool) {
	if daemon {
		return
	}
	v.mu.Lock()
	v.nonDaemon--
	if v.nonDaemon == 0 {
		v.cond.Broadcast()
	}
	v.mu.Unlock()
}

// WaitForNonDaemonThreads blocks until the non-daemon thread count reaches zero.
func (v *VM) WaitForNonDaemonThreads() {
	v.mu.Lock()
	for v.nonDaemon != 0 {
		v.cond.Wait()
	}
	v.mu.Unlock()
}

// --- package-level VM reference (set once by launch, read by native methods) ---

var theVM *VM

// SetVM registers the process-global VM supervisor. Must be called once before
// any thread-starting native methods execute.
func SetVM(v *VM) { theVM = v }

// GetVM returns the process-global VM supervisor, or nil if not set.
func GetVM() *VM { return theVM }
