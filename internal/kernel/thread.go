package kernel

import (
	"sync"
	"sync/atomic"
	"time"
)

// JThread is the kernel-side record behind one java.lang.Thread instance.
// The execution layer (VM) owns goroutine mechanics and installs itself via
// Kernel.SpawnJavaThread; the kernel owns identity, lifecycle state,
// interrupt flags and join latches so natives can implement Thread
// semantics without importing the VM.
type JThread struct {
	Obj  *Instance // the java.lang.Thread instance
	Key  uint64    // monitor-owner key (== VM thread id)
	Name string

	state atomic.Int32 // threadState constants
	alive atomic.Bool

	interrupted atomic.Bool

	joinLatch chan struct{} // closed on termination
	once      sync.Once

	mu       sync.Mutex
	sleepers map[*sleepRec]struct{} // active sleeps of THIS thread
}

type threadState int32

const (
	threadNew threadState = iota
	threadRunning
	threadTerminated
)

// StateName renders the lifecycle state (diagnostics/tests).
func (j *JThread) StateName() string {
	switch threadState(j.state.Load()) {
	case threadNew:
		return "NEW"
	case threadRunning:
		return "RUNNABLE"
	default:
		return "TERMINATED"
	}
}

func (j *JThread) IsAlive() bool { return j.alive.Load() }

// sleepRec is one interruptible sleep in progress.
type sleepRec struct {
	ch chan struct{} // closed by interrupt
}

func (j *JThread) addSleeper(rec *sleepRec) {
	j.mu.Lock()
	if j.sleepers == nil {
		j.sleepers = make(map[*sleepRec]struct{})
	}
	j.sleepers[rec] = struct{}{}
	j.mu.Unlock()
}

func (j *JThread) removeSleeper(rec *sleepRec) {
	j.mu.Lock()
	delete(j.sleepers, rec)
	j.mu.Unlock()
}

// interruptSleeps closes every active sleep; reports whether any existed.
func (j *JThread) interruptSleeps() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	if len(j.sleepers) == 0 {
		return false
	}
	for rec := range j.sleepers {
		close(rec.ch)
	}
	j.sleepers = make(map[*sleepRec]struct{})
	return true
}

// Terminate marks the thread dead and releases joiners. Idempotent.
func (j *JThread) Terminate() {
	j.once.Do(func() {
		j.alive.Store(false)
		j.state.Store(int32(threadTerminated))
		close(j.joinLatch)
	})
}

// ThreadRegistry tracks live threads and their blocking operations, giving
// Interrupt a path to wake Object.wait waiters and sleepers regardless of
// which monitor they were parked on.
type ThreadRegistry struct {
	mu   sync.Mutex
	byID map[uint64]*JThread
	main *JThread

	waiting map[uint64][]*waitRec // key -> active Object.wait records
}

func NewThreadRegistry() *ThreadRegistry {
	return &ThreadRegistry{
		byID:    make(map[uint64]*JThread),
		waiting: make(map[uint64][]*waitRec),
	}
}

// NewRecord creates an unattached record for a not-yet-started thread.
// attach() inserts it once start() mints its identity.
func (r *ThreadRegistry) NewRecord(obj *Instance, name string) *JThread {
	return &JThread{
		Obj:       obj,
		Name:      name,
		joinLatch: make(chan struct{}),
	}
}

// attach inserts a started record under its minted key.
func (r *ThreadRegistry) attach(j *JThread) {
	r.mu.Lock()
	r.byID[j.Key] = j
	r.mu.Unlock()
}

// InterruptByKey is the target-object entry point for Thread.interrupt:
// unstarted threads only get the flag; started ones additionally wake
// blocking operations.
func (r *ThreadRegistry) InterruptByKey(j *JThread) bool {
	if j.Key == 0 {
		j.interrupted.Store(true)
		return false
	}
	return r.Interrupt(j.Key)
}

// Register creates the record for a new thread identity.
func (r *ThreadRegistry) Register(key uint64, obj *Instance, name string) *JThread {
	j := &JThread{
		Obj:       obj,
		Key:       key,
		Name:      name,
		joinLatch: make(chan struct{}),
	}
	j.state.Store(int32(threadNew))
	r.mu.Lock()
	r.byID[key] = j
	r.mu.Unlock()
	return j
}

// ByKey returns the record for a monitor-owner key (nil when unknown).
func (r *ThreadRegistry) ByKey(key uint64) *JThread {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.byID[key]
}

// SetMain designates the primordial thread record.
func (r *ThreadRegistry) SetMain(j *JThread) {
	r.mu.Lock()
	r.main = j
	r.mu.Unlock()
}

// Main returns the primordial thread record.
func (r *ThreadRegistry) Main() *JThread {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.main
}

// publishWait records an Object.wait parking under owner key.
func (r *ThreadRegistry) publishWait(key uint64, rec *waitRec) {
	r.mu.Lock()
	r.waiting[key] = append(r.waiting[key], rec)
	r.mu.Unlock()
}

// retractWait removes a consumed record.
func (r *ThreadRegistry) retractWait(key uint64, rec *waitRec) {
	r.mu.Lock()
	list := r.waiting[key]
	for i, w := range list {
		if w == rec {
			r.waiting[key] = append(list[:i], list[i+1:]...)
			break
		}
	}
	if len(r.waiting[key]) == 0 {
		delete(r.waiting, key)
	}
	r.mu.Unlock()
}

// Interrupt implements Thread.interrupt semantics for the target key:
//
//   - sets the interrupt flag;
//   - wakes Object.wait waiters with outcome=waitGotInterrupt;
//   - wakes in-progress sleeps.
//
// Returns true when the thread was blocked in an interruptible operation.
func (r *ThreadRegistry) Interrupt(key uint64) bool {
	j := r.ByKey(key)
	if j != nil {
		j.interrupted.Store(true)
	} else {
		// Unknown key (e.g. racing termination): flag lives on the record,
		// nothing else to do.
		return false
	}

	woken := false
	r.mu.Lock()
	for _, rec := range r.waiting[key] {
		if rec.outcome.Load() == int32(waitStillWaiting) {
			rec.outcome.Store(int32(waitGotInterrupt))
			close(rec.ch)
			woken = true
		}
	}
	delete(r.waiting, key)
	r.mu.Unlock()

	if j.interruptSleeps() {
		woken = true
	}
	return woken
}

// SleepersNow reports active sleep count (diagnostics/tests).
func (j *JThread) SleepersNow() int {
	j.mu.Lock()
	defer j.mu.Unlock()
	return len(j.sleepers)
}

// ClearInterrupted peeks-and-clears the flag; returns previous value.
func (r *ThreadRegistry) ClearInterrupted(key uint64) bool {
	j := r.ByKey(key)
	if j == nil {
		return false
	}
	return j.interrupted.Swap(false)
}

// PeekInterrupted reads the flag without clearing.
func (r *ThreadRegistry) PeekInterrupted(key uint64) bool {
	j := r.ByKey(key)
	if j == nil {
		return false
	}
	return j.interrupted.Load()
}

// Sleep parks the thread interruptibly for d. Returns false when
// interrupted (flag already cleared by Interrupt's contract upstream).
func (r *ThreadRegistry) Sleep(key uint64, d time.Duration) (completed bool) {
	j := r.ByKey(key)
	if j == nil {
		time.Sleep(d)
		return true
	}
	// double-check pattern: flag may have been set before we parked
	if j.interrupted.Load() {
		j.interrupted.Store(false)
		return false
	}
	rec := &sleepRec{ch: make(chan struct{})}
	j.addSleeper(rec)
	defer j.removeSleeper(rec)

	// re-check after registration closes the race window
	if j.interrupted.Swap(false) {
		return false
	}

	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-rec.ch:
		return false
	case <-timer.C:
		return true
	}
}

// Join waits for the target record to terminate, interruptibly via the
// self thread's sleeper channel (same mechanism as Sleep).
func (r *ThreadRegistry) Join(selfKey uint64, target *JThread, d time.Duration) (reachedEnd bool, interrupted bool) {
	self := r.ByKey(selfKey)
	if self != nil && self.interrupted.Swap(false) {
		return false, true
	}
	var sleepRecCh chan struct{}
	var rec *sleepRec
	if self != nil {
		rec = &sleepRec{ch: make(chan struct{})}
		self.addSleeper(rec)
		sleepRecCh = rec.ch
		defer self.removeSleeper(rec)
		// close the set-flag-before-park race window
		if self.interrupted.Swap(false) {
			return false, true
		}
	}
	var timeout <-chan time.Time
	if d > 0 {
		t := time.NewTimer(d)
		defer t.Stop()
		timeout = t.C
	}
	select {
	case <-target.joinLatch:
		return true, false
	case <-timeout:
		return false, false
	case <-sleepRecCh:
		return false, true
	}
}
