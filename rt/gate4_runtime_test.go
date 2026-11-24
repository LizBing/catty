package rt

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------- helpers

// testClock is a manual monotonic clock. Every Java timer wait consumes it, so
// tests can force a timeout transition at a chosen point.
type testClock struct {
	mu     sync.Mutex
	now    time.Time
	timers []*testTimer
}

type testTimer struct {
	at time.Time
	ch chan time.Time
}

func newTestClock() *testClock {
	return &testClock{now: time.Unix(1_700_000_000, 0)}
}

func (c *testClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *testClock) After(d time.Duration) <-chan time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	ch := make(chan time.Time, 1)
	c.timers = append(c.timers, &testTimer{at: c.now.Add(d), ch: ch})
	return ch
}

// Advance moves the clock and fires every due timer exactly once.
func (c *testClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	var due []*testTimer
	keep := c.timers[:0]
	for _, tm := range c.timers {
		if !tm.at.After(c.now) {
			due = append(due, tm)
			continue
		}
		keep = append(keep, tm)
	}
	c.timers = keep
	c.mu.Unlock()
	for _, tm := range due {
		tm.ch <- tm.at
	}
}

// pendingTimers reports armed timers (test observability).
func (c *testClock) pendingTimers() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.timers)
}

// gate4Runtime installs a minimal Gate 4 image. It uses no application
// classes unless the test supplies them.
func gate4Runtime(t *testing.T, classes ...*ClassDef) *Runtime {
	t.Helper()
	img := &Image{
		ABI:     Gate4ABI,
		Profile: Gate4ProfileID,
		Plugins: true,
		Threads: true,
		Limits:  Gate3PluginLimits,
		Classes: classes,
	}
	r := NewRuntime(img)
	if err := r.Install(); err != nil {
		t.Fatalf("install: %v", err)
	}
	return r
}

// testThread creates a logical thread bound to a started peer.
func testThread(r *Runtime, name string) *Thread {
	t := r.NewThread()
	peer := r.newPeerObject(name, false, threadStarted)
	r.thrMu.Lock()
	if p, ok := peerOfObject(peer); ok {
		p.thread = t
	}
	t.peer = peer
	r.thrMu.Unlock()
	return t
}

// newTestObject allocates an ordinary java.lang.Object.
func newTestObject(r *Runtime) *Object {
	return r.NewInstance(r.mustClass("java/lang/Object"))
}

func aotMethod(owner, name, descriptor string, flags uint16, entry Entry) *MethodDef {
	return &MethodDef{
		Owner: owner, Name: name, Descriptor: descriptor, Flags: flags,
		Engine: "aot", Entry: entry, EntryIndex: -1,
	}
}

func staticField(c *Class, name string) *Field {
	for _, f := range c.StaticField {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func instanceField(c *Class, name string) *Field {
	for _, f := range c.InstanceField {
		if f.Name == name {
			return f
		}
	}
	return nil
}

// waitFor polls cond until it is true or the watchdog expires.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		runtime.Gosched()
		time.Sleep(time.Millisecond)
	}
}

// waitForGoroutine is waitFor with an explicit done channel.
func waitForGoroutine(t *testing.T, what string, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
}

func monitorOf(o *Object) *monitor { return o.monitor }

func monitorDepth(r *Runtime, o *Object) int {
	if o == nil || o.monitor == nil {
		return 0
	}
	return o.monitor.depth
}

func monitorOwner(r *Runtime, o *Object) *Thread {
	if o == nil || o.monitor == nil {
		return nil
	}
	return o.monitor.owner
}

func readEvents(t *testing.T, path string) []map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read trace: %v", err)
	}
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		var record map[string]any
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			t.Fatalf("decode trace line %q: %v", line, err)
		}
		out = append(out, record)
	}
	return out
}

func containsName(names []string, want string) bool {
	for _, name := range names {
		if name == want {
			return true
		}
	}
	return false
}

func eventNames(events []map[string]any, name string) []map[string]any {
	var out []map[string]any
	for _, event := range events {
		if event["event"] == name {
			out = append(out, event)
		}
	}
	return out
}

// ---------------------------------------------------------------- monitors

func TestGate4MonitorReentrantOwnership(t *testing.T) {
	r := gate4Runtime(t)
	owner := testThread(r, "owner")
	contender := testThread(r, "contender")
	obj := newTestObject(r)

	for i := 0; i < 3; i++ {
		if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
			t.Fatalf("reenter %d: %v", i, thrown)
		}
	}
	if got := monitorDepth(r, obj); got != 3 {
		t.Fatalf("recursion depth = %d, want 3", got)
	}

	done := make(chan Outcome, 1)
	go func() {
		done <- Outcome{Thrown: contender.MonitorEnter(nil, 1, RefValue(obj)), Value: RefValue(obj)}
	}()
	waitFor(t, "contender parked", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return obj.monitor.contenders == 1
	})
	if got := monitorOwner(r, obj); got != owner {
		t.Fatalf("contender stole the monitor at depth %d", monitorDepth(r, obj))
	}
	if thrown := owner.MonitorExit(nil, 2, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case <-done:
		t.Fatal("contender acquired the monitor after an intermediate exit")
	default:
	}
	if thrown := owner.MonitorExit(nil, 3, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case <-done:
		t.Fatal("contender acquired the monitor before the final exit")
	default:
	}
	if thrown := owner.MonitorExit(nil, 4, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case out := <-done:
		if out.Thrown != nil {
			t.Fatalf("contender failed: %v", out.Thrown)
		}
	case <-time.After(time.Second):
		t.Fatal("contender stayed blocked after the final exit")
	}
	if got := monitorOwner(r, obj); got != contender {
		t.Fatal("contender did not acquire the monitor after the final exit")
	}
	if thrown := contender.MonitorExit(nil, 4, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if got := monitorDepth(r, obj); got != 0 {
		t.Fatalf("monitor leaked depth %d", got)
	}
}

func TestGate4MonitorIllegalExitAndNull(t *testing.T) {
	r := gate4Runtime(t)
	owner := testThread(r, "owner")
	other := testThread(r, "other")
	obj := newTestObject(r)

	if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if thrown := other.MonitorExit(nil, 1, RefValue(obj)); thrown == nil ||
		thrown.Class == nil || thrown.Class.Name != "java/lang/IllegalMonitorStateException" {
		t.Fatalf("non-owner exit outcome = %v", thrown)
	}
	if got := monitorOwner(r, obj); got != owner || monitorDepth(r, obj) != 1 {
		t.Fatal("illegal exit changed owner or depth")
	}
	if thrown := owner.MonitorEnter(nil, 2, RefValue(nil)); thrown == nil ||
		thrown.Class.Name != "java/lang/NullPointerException" {
		t.Fatalf("null monitorenter outcome = %v", thrown)
	}
	if thrown := owner.MonitorExit(nil, 3, RefValue(nil)); thrown == nil ||
		thrown.Class.Name != "java/lang/NullPointerException" {
		t.Fatalf("null monitorexit outcome = %v", thrown)
	}
}

func TestGate4MonitoredObjectCollectible(t *testing.T) {
	r := gate4Runtime(t)
	owner := testThread(r, "owner")
	finalized := make(chan struct{}, 1)
	func() {
		obj := newTestObject(r)
		runtime.SetFinalizer(obj, func(*Object) { finalized <- struct{}{} })
		if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
			t.Fatal(thrown)
		}
		if thrown := owner.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
			t.Fatal(thrown)
		}
		obj = nil
		_ = obj
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		runtime.GC()
		select {
		case <-finalized:
			return
		case <-time.After(5 * time.Millisecond):
			if time.Now().After(deadline) {
				t.Fatal("a previously monitored object was retained by monitor state")
			}
		}
	}
}

// ------------------------------------------------------------ method entry

func TestGate4ImplicitSynchronizedMethodLocks(t *testing.T) {
	var entered, reentered int32
	var instLock, staticLock *Object
	var outerDepth, callbackDepth int32
	instDecl := aotMethod("app/Monitor", "inst", "()V", AccPublic|AccSynchronized,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			atomic.AddInt32(&entered, 1)
			if !th.holdsLockObject(receiver) {
				t.Error("instance synchronized body does not own the receiver")
			}
			depth := monitorDepth(th.RT, receiver)
			instLock = receiver
			if depth == 1 {
				outerDepth = int32(depth)
				out := th.Call(m, receiver, nil)
				if out.Thrown != nil {
					return out
				}
				atomic.AddInt32(&reentered, 1)
				return Normal()
			}
			callbackDepth = int32(depth)
			return Normal()
		})
	staticDecl := aotMethod("app/Monitor", "stat", "()V", AccPublic|AccStatic|AccSynchronized,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			if !th.holdsLockObject(m.Owner.Mirror) {
				t.Error("static synchronized body does not own the declaring class mirror")
			}
			staticLock = m.Owner.Mirror
			return Normal()
		})
	parent := &ClassDef{Name: "app/Monitor", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{instDecl, staticDecl}}
	child := &ClassDef{Name: "app/MonitorChild", Super: "app/Monitor", Flags: AccPublic}
	r := gate4Runtime(t, parent, child)

	c, err := r.Class("app/Monitor")
	if err != nil {
		t.Fatal(err)
	}
	inst := c.MethodIndex["inst()V"]
	receiver := r.NewInstance(c)
	th := testThread(r, "main")
	out := th.Call(inst, receiver, nil)
	if out.Thrown != nil {
		t.Fatalf("synchronized instance call: %v", out.Thrown)
	}
	if atomic.LoadInt32(&entered) != 2 || atomic.LoadInt32(&reentered) != 1 {
		t.Fatalf("callback accounting = %d/%d, want 2/1", entered, reentered)
	}
	if instLock != receiver {
		t.Fatal("instance method locked a different object")
	}
	if outerDepth != 1 || callbackDepth != 2 {
		t.Fatalf("implicit recursion depths = %d/%d, want 1/2", outerDepth, callbackDepth)
	}
	if monitorDepth(r, receiver) != 0 {
		t.Fatal("implicit instance monitor leaked after return")
	}

	// Inherited static selection locks the declaring class, not the receiver
	// subclass.
	stat, err := r.ResolveMethod("app/MonitorChild", "stat", "()V")
	if err != nil {
		t.Fatal(err)
	}
	out = th.Call(stat, nil, nil)
	if out.Thrown != nil {
		t.Fatalf("inherited static synchronized call: %v", out.Thrown)
	}
	if staticLock != c.Mirror || staticLock == nil {
		t.Fatal("inherited static method locked a different class mirror")
	}
	if monitorDepth(r, c.Mirror) != 0 {
		t.Fatal("implicit static monitor leaked after return")
	}

	// Two receivers use independent instance monitors: holding one receiver's
	// monitor must not block a call on the other.
	other := r.NewInstance(c)
	if thrown := th.MonitorEnter(nil, 9, RefValue(other)); thrown != nil {
		t.Fatal(thrown)
	}
	done := make(chan Outcome, 1)
	go func() { done <- th.Call(inst, receiver, nil) }()
	select {
	case o := <-done:
		if o.Thrown != nil {
			t.Fatalf("second receiver call failed: %v", o.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("per-receiver instance locks were shared")
	}
	if thrown := th.MonitorExit(nil, 10, RefValue(other)); thrown != nil {
		t.Fatal(thrown)
	}
}

func TestGate4SynchronizedMethodUnwindAndCatch(t *testing.T) {
	// A throwing synchronized method preserves throwable identity and releases
	// exactly one implicit acquisition on unwind.
	var thrownObject *Object
	unwindEntry := func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		if !th.holdsLockObject(receiver) {
			t.Error("throwing synchronized body lost its monitor")
		}
		return ThrowOutcome(thrownObject)
	}
	unwindClass := &ClassDef{Name: "app/Guard", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{
			aotMethod("app/Guard", "thrower", "()V", AccPublic|AccSynchronized, unwindEntry),
		}}
	r := gate4Runtime(t, unwindClass)
	thrownObject = r.NewThrowable(nil, "java/lang/RuntimeException")
	c := r.mustClass("app/Guard")
	receiver := r.NewInstance(c)
	th := testThread(r, "main")
	out := th.Call(c.MethodIndex["thrower()V"], receiver, nil)
	if out.Thrown != thrownObject {
		t.Fatal("unwind replaced the throwable identity")
	}
	if monitorDepth(r, receiver) != 0 {
		t.Fatal("synchronized method unwound without releasing its monitor")
	}

	// A same-frame catch executes with the implicit monitor still held. The
	// handler calls Thread.holdsLock(receiver) and stores the result.
	class := &ClassDef{
		Name: "app/Catcher", Super: "java/lang/Object", Flags: AccPublic,
		StaticFields: []FieldDef{{Name: "HELD", Descriptor: "I", Flags: AccStatic}},
	}
	method := &MethodDef{
		Owner: "app/Catcher", Name: "guard", Descriptor: "()V",
		Flags: AccPublic | AccSynchronized, Engine: "interp",
		MaxStack: 2, MaxLocals: 2, CodeSize: 17, EntryIndex: -1,
		Insns: []Insn{
			{Op: OpNew, PC: 0, Ref: 0, Len: 3},
			{Op: OpDup, PC: 3, Ref: -1, Len: 1},
			{Op: OpInvokespecial, PC: 4, Ref: 1, Len: 3},
			{Op: OpAthrow, PC: 7, Ref: -1, Len: 1},
			{Op: OpAstore, PC: 8, A: 1, Ref: -1, Len: 1},
			{Op: OpAload, PC: 9, A: 0, Ref: -1, Len: 1},
			{Op: OpInvokestatic, PC: 10, Ref: 2, Len: 3},
			{Op: OpPutstatic, PC: 13, Ref: 3, Len: 3},
			{Op: OpReturn, PC: 16, Ref: -1, Len: 1},
		},
		Refs: []Ref{
			{Kind: RefClass, Class: "java/lang/RuntimeException", CP: 1},
			{Kind: RefMethod, Owner: "java/lang/RuntimeException", Name: "<init>", Descriptor: "()V", CP: 2},
			{Kind: RefMethod, Owner: "java/lang/Thread", Name: "holdsLock", Descriptor: "(Ljava/lang/Object;)Z", CP: 3},
			{Kind: RefField, Owner: "app/Catcher", Name: "HELD", Descriptor: "I", CP: 4},
		},
		Handlers: []Handler{{StartPC: 0, EndPC: 8, HandlerPC: 8, CatchType: "java/lang/RuntimeException"}},
	}
	class.Methods = []*MethodDef{method}
	r3 := gate4Runtime(t, class)
	c3 := r3.mustClass("app/Catcher")
	receiver3 := r3.NewInstance(c3)
	th3 := testThread(r3, "main")
	out3 := th3.Call(c3.MethodIndex["guard()V"], receiver3, nil)
	if out3.Thrown != nil {
		t.Fatalf("same-frame catch escaped: %v", out3.Thrown)
	}
	if v := r3.staticRead(c3, staticField(c3, "HELD").Index); AsInt(v) != 1 {
		t.Fatalf("handler did not observe the implicit monitor: HELD=%d", AsInt(v))
	}
	if monitorDepth(r3, receiver3) != 0 {
		t.Fatal("implicit monitor leaked after the handler returned")
	}
}

// ------------------------------------------------------------- wait sets

type waitResult struct {
	out   Outcome
	depth int
	ok    bool
}

// startWaiter enters obj depth times, records the restored depth, exits the
// same number of times and reports the wait outcome.
func startWaiter(link *Thread, obj *Object, depth int, timeout time.Duration, timed bool, ready chan<- struct{}, done chan<- waitResult) {
	if thrown := link.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		done <- waitResult{out: ThrowOutcome(thrown)}
		close(ready)
		return
	}
	for i := 1; i < depth; i++ {
		if thrown := link.MonitorEnter(nil, int32(i), RefValue(obj)); thrown != nil {
			done <- waitResult{out: ThrowOutcome(thrown)}
			close(ready)
			return
		}
	}
	close(ready)
	out := link.objectWait(obj, timeout, timed)
	restored := monitorDepth(link.RT, obj)
	for i := 0; i < depth; i++ {
		if thrown := link.MonitorExit(nil, int32(i), RefValue(obj)); thrown != nil {
			done <- waitResult{out: ThrowOutcome(thrown)}
			return
		}
	}
	done <- waitResult{out: out, depth: restored, ok: true}
}

func TestGate4WaitSetNotifyAndNotifyAll(t *testing.T) {
	r := gate4Runtime(t)
	obj := newTestObject(r)
	const depth = 2
	waiters := make([]*Thread, 3)
	ready := make([]chan struct{}, 3)
	done := make([]chan waitResult, 3)
	for i := range waiters {
		waiters[i] = testThread(r, "waiter")
		ready[i] = make(chan struct{})
		done[i] = make(chan waitResult, 1)
	}
	for i := range waiters {
		go startWaiter(waiters[i], obj, depth, 0, false, ready[i], done[i])
	}
	for i := range ready {
		waitForGoroutine(t, "waiter enrollment", ready[i])
	}
	waitFor(t, "three enrolled waiters", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 3
	})
	owner := testThread(r, "owner")
	if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if out := owner.notifyObject(obj, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	waitFor(t, "one notification", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 2
	})
	if out := owner.notifyObject(obj, true); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if thrown := owner.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	for i := range done {
		select {
		case result := <-done[i]:
			if result.out.Thrown != nil {
				t.Fatalf("waiter %d outcome: %v", i, result.out.Thrown)
			}
			if !result.ok || result.depth != depth {
				t.Fatalf("waiter %d restored depth = %d, want %d", i, result.depth, depth)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("waiter %d stranded", i)
		}
	}
	if monitorDepth(r, obj) != 0 || monitorOwner(r, obj) != nil {
		t.Fatal("notify path leaked the monitor")
	}
}

func TestGate4WaitSetTimeoutAndInterrupt(t *testing.T) {
	r := gate4Runtime(t)
	clock := newTestClock()
	r.clock = clock
	obj := newTestObject(r)
	const depth = 2

	// Timeout: the owner keeps the monitor while the timer fires, so the waiter
	// stays blocked in reacquisition until it is released.
	timed := testThread(r, "timed")
	timedReady := make(chan struct{})
	timedDone := make(chan waitResult, 1)
	go startWaiter(timed, obj, depth, time.Second, true, timedReady, timedDone)
	waitForGoroutine(t, "timed waiter enrollment", timedReady)
	owner := testThread(r, "owner")
	if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	clock.Advance(time.Second)
	waitFor(t, "timeout transition", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 0
	})
	select {
	case <-timedDone:
		t.Fatal("timed waiter resumed before the owner released the monitor")
	default:
	}
	if thrown := owner.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-timedDone:
		if result.out.Thrown != nil {
			t.Fatalf("timed wait outcome: %v", result.out.Thrown)
		}
		if result.depth != depth {
			t.Fatalf("timed wait restored depth %d, want %d", result.depth, depth)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed waiter stranded")
	}

	// Interrupt: the winning transition throws after reacquiring the full depth
	// and clears the interrupt status.
	interrupted := testThread(r, "interrupted")
	interruptReady := make(chan struct{})
	interruptDone := make(chan waitResult, 1)
	go startWaiter(interrupted, obj, depth, 0, false, interruptReady, interruptDone)
	waitForGoroutine(t, "interrupt waiter enrollment", interruptReady)
	owner = testThread(r, "owner2")
	if thrown := owner.MonitorEnter(nil, 2, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	owner.interruptPeer(interrupted.peer)
	if thrown := owner.MonitorExit(nil, 3, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-interruptDone:
		if result.out.Thrown == nil || result.out.Thrown.Class.Name != "java/lang/InterruptedException" {
			t.Fatalf("interrupted wait outcome: %v", result.out.Thrown)
		}
		if result.depth != depth {
			t.Fatalf("interrupted wait restored depth %d, want %d", result.depth, depth)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("interrupted waiter stranded")
	}
	r.thrMu.Lock()
	if p, _ := peerOfObject(interrupted.peer); p != nil && p.pendingInterrupt {
		t.Error("interrupted wait did not clear the interrupt status")
	}
	r.thrMu.Unlock()
}

func TestGate4WaitSetWinnerTransitions(t *testing.T) {
	r := gate4Runtime(t)
	obj := newTestObject(r)

	// Notification wins; a later interrupt stays pending on the normal return.
	first := testThread(r, "first")
	firstReady := make(chan struct{})
	firstDone := make(chan waitResult, 1)
	go startWaiter(first, obj, 1, 0, false, firstReady, firstDone)
	waitForGoroutine(t, "first enrollment", firstReady)
	owner := testThread(r, "owner")
	if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if out := owner.notifyObject(obj, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	owner.interruptPeer(first.peer)
	if thrown := owner.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-firstDone:
		if result.out.Thrown != nil {
			t.Fatalf("notified waiter threw: %v", result.out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("notified waiter stranded")
	}
	r.thrMu.Lock()
	p, _ := peerOfObject(first.peer)
	latePending := p != nil && p.pendingInterrupt
	if p != nil {
		p.pendingInterrupt = false
	}
	r.thrMu.Unlock()
	if !latePending {
		t.Fatal("late interrupt was not left pending after a normal return")
	}

	// Interrupt wins while a second waiter stays eligible; the later notify
	// claims the second waiter and never the interrupted one.
	loser := testThread(r, "loser")
	winner := testThread(r, "winner")
	loserReady := make(chan struct{})
	winnerReady := make(chan struct{})
	loserDone := make(chan waitResult, 1)
	winnerDone := make(chan waitResult, 1)
	go startWaiter(loser, obj, 1, 0, false, loserReady, loserDone)
	go startWaiter(winner, obj, 1, 0, false, winnerReady, winnerDone)
	waitForGoroutine(t, "loser enrollment", loserReady)
	waitForGoroutine(t, "winner enrollment", winnerReady)
	waitFor(t, "two enrolled", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 2
	})
	owner = testThread(r, "owner2")
	if thrown := owner.MonitorEnter(nil, 2, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	owner.interruptPeer(loser.peer)
	if out := owner.notifyObject(obj, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if thrown := owner.MonitorExit(nil, 3, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-loserDone:
		if result.out.Thrown == nil || result.out.Thrown.Class.Name != "java/lang/InterruptedException" {
			t.Fatalf("interrupted waiter outcome: %v", result.out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("interrupted waiter stranded")
	}
	select {
	case result := <-winnerDone:
		if result.out.Thrown != nil {
			t.Fatalf("second waiter outcome: %v", result.out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("second waiter stranded after notify")
	}

	// Timeout wins and a later interrupt stays pending; the stale interrupt
	// cannot reclassify the completed record.
	clock := newTestClock()
	r.clock = clock
	timed := testThread(r, "timed")
	timedReady := make(chan struct{})
	timedDone := make(chan waitResult, 1)
	go startWaiter(timed, obj, 1, time.Second, true, timedReady, timedDone)
	waitForGoroutine(t, "timed enrollment", timedReady)
	owner = testThread(r, "owner3")
	if thrown := owner.MonitorEnter(nil, 4, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	clock.Advance(time.Second)
	waitFor(t, "timeout transition", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 0
	})
	owner.interruptPeer(timed.peer)
	if thrown := owner.MonitorExit(nil, 5, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-timedDone:
		if result.out.Thrown != nil {
			t.Fatalf("timed-out waiter threw: %v", result.out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed-out waiter stranded")
	}
	r.thrMu.Lock()
	p, _ = peerOfObject(timed.peer)
	latePending = p != nil && p.pendingInterrupt
	if p != nil {
		p.pendingInterrupt = false
	}
	r.thrMu.Unlock()
	if !latePending {
		t.Fatal("late interrupt after timeout was not left pending")
	}
}

func TestGate4WaitRestoresDepthAndRetainsUnrelatedMonitor(t *testing.T) {
	r := gate4Runtime(t)
	obj := newTestObject(r)
	other := newTestObject(r)
	waiter := testThread(r, "waiter")
	ready := make(chan struct{})
	done := make(chan waitResult, 1)
	go func() {
		if thrown := waiter.MonitorEnter(nil, 0, RefValue(other)); thrown != nil {
			done <- waitResult{out: ThrowOutcome(thrown)}
			return
		}
		if thrown := waiter.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
			done <- waitResult{out: ThrowOutcome(thrown)}
			return
		}
		for i := 1; i < 3; i++ {
			if thrown := waiter.MonitorEnter(nil, int32(i), RefValue(obj)); thrown != nil {
				done <- waitResult{out: ThrowOutcome(thrown)}
				return
			}
		}
		close(ready)
		out := waiter.objectWait(obj, 0, false)
		restored := monitorDepth(waiter.RT, obj)
		unrelated := waiter.holdsLockObject(other)
		for i := 0; i < 3; i++ {
			if thrown := waiter.MonitorExit(nil, int32(i), RefValue(obj)); thrown != nil {
				done <- waitResult{out: ThrowOutcome(thrown)}
				return
			}
		}
		if thrown := waiter.MonitorExit(nil, 0, RefValue(other)); thrown != nil {
			done <- waitResult{out: ThrowOutcome(thrown)}
			return
		}
		if !unrelated {
			out = ThrowOutcome(waiter.RT.NewThrowable(nil, "java/lang/IllegalMonitorStateException"))
		}
		done <- waitResult{out: out, depth: restored, ok: true}
	}()
	waitForGoroutine(t, "deep waiter enrollment", ready)
	waitFor(t, "deep waiter recorded", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(obj.monitor.waiters) == 1
	})
	owner := testThread(r, "owner")
	if thrown := owner.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if out := owner.notifyObject(obj, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if thrown := owner.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	select {
	case result := <-done:
		if result.out.Thrown != nil {
			t.Fatalf("deep wait outcome: %v", result.out.Thrown)
		}
		if !result.ok || result.depth != 3 {
			t.Fatalf("restored recursion depth = %d, want 3", result.depth)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("deep waiter stranded")
	}
	if monitorDepth(r, obj) != 0 || monitorDepth(r, other) != 0 {
		t.Fatal("waiter leaked a monitor")
	}
}

// ----------------------------------------------------- timers and overloads

func TestGate4TimerOverloadsAndSaturation(t *testing.T) {
	if d, timed := waitTimeout(0, 0); timed || d != 0 {
		t.Fatalf("waitTimeout(0,0) = %v,%v; want indefinite", d, timed)
	}
	if d, timed := waitTimeout(0, 1); !timed || d != time.Millisecond {
		t.Fatalf("waitTimeout(0,1) = %v,%v; want one millisecond", d, timed)
	}
	if d, timed := waitTimeout(5, 999999); !timed || d != 6*time.Millisecond {
		t.Fatalf("waitTimeout(5,999999) = %v,%v; want rounded six milliseconds", d, timed)
	}
	if d, timed := waitTimeout(math.MaxInt64, 999999); !timed || d != maxTimerDuration {
		t.Fatalf("overflowing wait did not saturate: %v", d)
	}
	if d := sleepNanos(math.MaxInt64, 999999); d != maxTimerDuration {
		t.Fatalf("overflowing sleep did not saturate: %v", d)
	}
	if d := sleepNanos(1, 999999); d != time.Millisecond+999999*time.Nanosecond {
		t.Fatalf("sleepNanos(1,999999) = %v", d)
	}

	r := gate4Runtime(t)
	thread := testThread(r, "main")
	sleep := r.mustClass("java/lang/Thread").MethodIndex["sleep(J)V"]
	if out := intrinsicThreadSleep(thread, sleep, nil, []Value{LongValue(-1)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("negative sleep outcome = %v", out.Thrown)
	}
	sleep2 := r.mustClass("java/lang/Thread").MethodIndex["sleep(JI)V"]
	if out := intrinsicThreadSleep(thread, sleep2, nil, []Value{LongValue(0), IntValue(1000000)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("nanos overflow outcome = %v", out.Thrown)
	}
	if out := intrinsicThreadSleep(thread, sleep2, nil, []Value{LongValue(0), IntValue(0)}); out.Thrown != nil {
		t.Fatalf("zero sleep threw: %v", out.Thrown)
	}

	obj := newTestObject(r)
	wait := r.mustClass("java/lang/Object").MethodIndex["wait(JI)V"]
	target := testThread(r, "target")
	if thrown := target.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if out := intrinsicObjectWait(target, wait, obj, []Value{LongValue(0), IntValue(1000000)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("invalid wait nanos outcome = %v", out.Thrown)
	}
	// A pending interrupt wins over an otherwise zero-duration timed wait.
	target.setInterruptLocked()
	if out := intrinsicObjectWait(target, wait, obj, []Value{LongValue(0), IntValue(1)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/InterruptedException" {
		t.Fatalf("pre-interrupted wait outcome = %v", out.Thrown)
	}
	if thrown := target.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	// Zero-duration sleep throws for a pending interrupt even without waiting.
	target.setInterruptLocked()
	if out := intrinsicThreadSleep(target, sleep, nil, []Value{LongValue(0)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/InterruptedException" {
		t.Fatalf("zero sleep with pending interrupt = %v", out.Thrown)
	}
	join := r.mustClass("java/lang/Thread").MethodIndex["join(J)V"]
	fresh := r.newPeerObject("fresh", false, threadNew)
	if out := intrinsicThreadJoin(thread, join, fresh, []Value{LongValue(-1)}); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("negative join outcome = %v", out.Thrown)
	}
	if out := intrinsicThreadJoin(thread, join, fresh, []Value{LongValue(0)}); out.Thrown != nil {
		t.Fatalf("join on a new peer threw: %v", out.Thrown)
	}
	if out := intrinsicThreadYield(thread, nil, nil, nil); out.Thrown != nil {
		t.Fatalf("yield threw: %v", out.Thrown)
	}
}

// ----------------------------------------------------------- join and sleep

func TestGate4JoinLifecycleAndInterrupt(t *testing.T) {
	release := make(chan struct{})
	var ran int32
	task := aotMethod("app/Task", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		atomic.AddInt32(&ran, 1)
		<-release
		return Normal()
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Task", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{task}},
	)
	main := testThread(r, "main")
	peer := r.newPeerObject("worker", false, threadNew)
	payload := peerData(peer)
	payload.runnable = r.NewInstance(r.mustClass("app/Task"))
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(main, start, peer, nil); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	waitFor(t, "task entered", func() bool { return atomic.LoadInt32(&ran) == 1 })

	joined := make(chan Outcome, 2)
	go func() { joined <- main.joinPeer(peer, 0, false) }()
	joiner2 := testThread(r, "joiner2")
	go func() { joined <- joiner2.joinPeer(peer, 0, false) }()
	waitFor(t, "both joiners enrolled", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return peer.monitor != nil && len(peer.monitor.waiters) == 2
	})
	// Interrupting one joiner affects only that joiner; the target keeps
	// running and the other joiner keeps waiting.
	main.interruptPeer(joiner2.peer)
	select {
	case out := <-joined:
		if out.Thrown == nil || out.Thrown.Class.Name != "java/lang/InterruptedException" {
			t.Fatalf("interrupted join outcome = %v", out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("interrupted joiner did not return")
	}
	r.thrMu.Lock()
	if p, _ := peerOfObject(joiner2.peer); p != nil && p.pendingInterrupt {
		t.Error("interrupted join did not clear the interrupt status")
	}
	r.thrMu.Unlock()
	select {
	case <-joined:
		t.Fatal("second joiner returned before the target terminated")
	default:
	}
	close(release)
	select {
	case out := <-joined:
		if out.Thrown != nil {
			t.Fatalf("join after termination threw: %v", out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("joiner stranded after termination")
	}
	waitFor(t, "worker terminated", func() bool {
		p := peerData(peer)
		return p != nil && p.state == threadTerminated
	})
	// Joining a terminated peer returns immediately, and a restarted peer is
	// rejected forever.
	if out := main.joinPeer(peer, 0, false); out.Thrown != nil {
		t.Fatalf("join on a terminated peer threw: %v", out.Thrown)
	}
	if out := intrinsicThreadStart(main, start, peer, nil); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalThreadStateException" {
		t.Fatalf("restart outcome = %v", out.Thrown)
	}
	if got := r.LiveNonDaemon(); got != 0 {
		t.Fatalf("live non-daemon count = %d, want 0", got)
	}
}

func TestGate4TimedJoinOneDeadline(t *testing.T) {
	r := gate4Runtime(t)
	clock := newTestClock()
	r.clock = clock
	main := testThread(r, "main")
	peer := r.newPeerObject("target", false, threadStarted)
	targetThread := testThread(r, "target-thread")
	r.thrMu.Lock()
	peerData(peer).thread = targetThread
	r.thrMu.Unlock()

	done := make(chan Outcome, 1)
	go func() { done <- main.joinPeer(peer, 10*time.Second, true) }()
	waitFor(t, "join enrolled", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return peer.monitor != nil && len(peer.monitor.waiters) == 1
	})
	// An unrelated notify wakes the joiner early; the single deadline means it
	// recomputes the remaining time and waits again instead of restarting.
	notifier := testThread(r, "notifier")
	if thrown := notifier.MonitorEnter(nil, 0, RefValue(peer)); thrown != nil {
		t.Fatal(thrown)
	}
	if out := notifier.notifyObject(peer, true); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if thrown := notifier.MonitorExit(nil, 1, RefValue(peer)); thrown != nil {
		t.Fatal(thrown)
	}
	clock.Advance(4 * time.Second)
	waitFor(t, "join re-enrolled", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return len(peer.monitor.waiters) == 1
	})
	select {
	case <-done:
		t.Fatal("join returned before its single deadline")
	default:
	}
	clock.Advance(6 * time.Second)
	select {
	case out := <-done:
		if out.Thrown != nil {
			t.Fatalf("timed join threw: %v", out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed join did not return at its deadline")
	}
	if pending := clock.pendingTimers(); pending != 0 {
		t.Fatalf("%d timer segments were left armed", pending)
	}
}

func TestGate4SleepInterruptRetainsMonitors(t *testing.T) {
	r := gate4Runtime(t)
	clock := newTestClock()
	r.clock = clock
	sleeper := testThread(r, "sleeper")
	other := newTestObject(r)
	if thrown := sleeper.MonitorEnter(nil, 0, RefValue(other)); thrown != nil {
		t.Fatal(thrown)
	}
	done := make(chan Outcome, 1)
	go func() { done <- sleeper.sleepFor(time.Duration(math.MaxInt64)) }()
	waitFor(t, "sleep blocker", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return sleeper.blocker != nil && sleeper.blocker.kind == blockerSleep
	})
	if !sleeper.holdsLockObject(other) {
		t.Fatal("sleep released an owned monitor")
	}
	sleeper.interruptPeer(sleeper.peer)
	select {
	case out := <-done:
		if out.Thrown == nil || out.Thrown.Class.Name != "java/lang/InterruptedException" {
			t.Fatalf("interrupted sleep outcome: %v", out.Thrown)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("interrupted sleep did not return")
	}
	if !sleeper.holdsLockObject(other) {
		t.Fatal("interrupted sleep lost an owned monitor")
	}
	r.thrMu.Lock()
	if p, _ := peerOfObject(sleeper.peer); p != nil && p.pendingInterrupt {
		t.Error("interrupted sleep did not clear the interrupt status")
	}
	r.thrMu.Unlock()
	if thrown := sleeper.MonitorExit(nil, 1, RefValue(other)); thrown != nil {
		t.Fatal(thrown)
	}
}

// ------------------------------------------------------------- accounting

func TestGate4ThreadStartAndDaemonRaces(t *testing.T) {
	var runs int32
	task := aotMethod("app/Task", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		atomic.AddInt32(&runs, 1)
		return Normal()
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Task", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{task}},
	)
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	makePeer := func() *Object {
		peer := r.newPeerObject("race", false, threadNew)
		peerData(peer).runnable = r.NewInstance(r.mustClass("app/Task"))
		return peer
	}

	// Two simultaneous first starts: exactly one wins, even after termination.
	peer := makePeer()
	results := make(chan Outcome, 2)
	gate := make(chan struct{})
	for i := 0; i < 2; i++ {
		go func() {
			<-gate
			results <- intrinsicThreadStart(testThread(r, "starter"), start, peer, nil)
		}()
	}
	close(gate)
	var winners int
	for i := 0; i < 2; i++ {
		out := <-results
		if out.Thrown == nil {
			winners++
		} else if out.Thrown.Class.Name != "java/lang/IllegalThreadStateException" {
			t.Fatalf("start race loser outcome: %v", out.Thrown)
		}
	}
	if winners != 1 {
		t.Fatalf("%d first starts won, want 1", winners)
	}
	waitFor(t, "single run", func() bool { return r.LiveNonDaemon() == 0 })
	if got := atomic.LoadInt32(&runs); got != 1 {
		t.Fatalf("task ran %d times, want 1", got)
	}

	// setDaemon and start serialize on the peer: a daemon change is either
	// observed before the transition or rejected by IllegalThreadStateException.
	for i := 0; i < 32; i++ {
		p := makePeer()
		setDaemon := r.mustClass("java/lang/Thread").MethodIndex["setDaemon(Z)V"]
		started := make(chan Outcome, 1)
		daemoned := make(chan Outcome, 1)
		gate := make(chan struct{})
		go func() {
			<-gate
			started <- intrinsicThreadStart(testThread(r, "starter"), start, p, nil)
		}()
		go func() {
			<-gate
			daemoned <- intrinsicThreadSetDaemon(testThread(r, "daemoner"), setDaemon, p, []Value{IntValue(1)})
		}()
		close(gate)
		startOut := <-started
		daemonOut := <-daemoned
		payload := peerData(p)
		if startOut.Thrown != nil {
			t.Fatalf("start failed under setDaemon race: %v", startOut.Thrown)
		}
		if daemonOut.Thrown == nil {
			if !payload.daemon {
				t.Fatal("successful setDaemon was not published")
			}
		} else if daemonOut.Thrown.Class.Name != "java/lang/IllegalThreadStateException" {
			t.Fatalf("setDaemon race outcome: %v", daemonOut.Thrown)
		}
		r.WaitForNonDaemon()
	}
}

func TestGate4LauncherWaitsForDescendant(t *testing.T) {
	var children, completed int32
	var childPeer *Object
	childRun := aotMethod("app/Child", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		atomic.AddInt32(&completed, 1)
		return Normal()
	})
	parentRun := aotMethod("app/Parent", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		atomic.AddInt32(&children, 1)
		out := intrinsicThreadStart(th, th.RT.mustClass("java/lang/Thread").MethodIndex["start()V"], childPeer, nil)
		return out
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Child", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{childRun}},
		&ClassDef{Name: "app/Parent", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{parentRun}},
	)
	main := testThread(r, "main")
	childPeer = r.newPeerObject("child", false, threadNew)
	peerData(childPeer).runnable = r.NewInstance(r.mustClass("app/Child"))
	parentPeer := r.newPeerObject("parent", false, threadNew)
	peerData(parentPeer).runnable = r.NewInstance(r.mustClass("app/Parent"))
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(main, start, parentPeer, nil); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	// The launcher's accounting must include the child started by the parent
	// after the parent itself was counted.
	r.WaitForNonDaemon()
	if atomic.LoadInt32(&children) != 1 || atomic.LoadInt32(&completed) != 1 {
		t.Fatalf("descendant accounting = %d/%d, want 1/1", children, completed)
	}
	if got := r.LiveNonDaemon(); got != 0 {
		t.Fatalf("live non-daemon count = %d after both workers", got)
	}
}

// ------------------------------------------------------ class initialization

func TestGate4ClassInitializationContention(t *testing.T) {
	var value *Field
	started := make(chan struct{})
	release := make(chan struct{})
	clinit := aotMethod("app/Ready", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			c := th.RT.mustClass("app/Ready")
			if thrown := th.EnsureInitialized(c); thrown != nil {
				return ThrowOutcome(thrown)
			}
			close(started)
			<-release
			th.RT.staticWrite(c, value.Index, IntValue(73))
			return Normal()
		})
	r := gate4Runtime(t, &ClassDef{
		Name: "app/Ready", Super: "java/lang/Object", Flags: AccPublic,
		StaticFields: []FieldDef{{Name: "VALUE", Descriptor: "I", Flags: AccStatic}},
		Methods:      []*MethodDef{clinit},
	})
	c := r.mustClass("app/Ready")
	value = staticField(c, "VALUE")
	owner := testThread(r, "owner")
	contender := testThread(r, "contender")
	contender.setInterruptLocked()

	ownerDone := make(chan *Object, 1)
	go func() { ownerDone <- owner.EnsureInitialized(c) }()
	waitForGoroutine(t, "initializer start", started)

	waited := make(chan struct{})
	r.ClassInitWaitHook = func(hookClass *Class, waiter *Thread) {
		if hookClass == c && waiter == contender {
			select {
			case <-waited:
			default:
				close(waited)
			}
		}
	}
	contenderDone := make(chan *Object, 1)
	go func() { contenderDone <- contender.EnsureInitialized(c) }()
	waitForGoroutine(t, "competing initialization wait", waited)
	if got := r.staticRead(c, value.Index); AsInt(got) != 0 {
		t.Fatal("competing thread observed a partially published initializer write")
	}
	r.ClassInitWaitHook = nil
	close(release)
	if thrown := <-ownerDone; thrown != nil {
		t.Fatalf("owner initialization failed: %v", thrown)
	}
	if thrown := <-contenderDone; thrown != nil {
		t.Fatalf("contender initialization failed: %v", thrown)
	}
	if got := r.staticRead(c, value.Index); AsInt(got) != 73 {
		t.Fatalf("initializer write was not published: VALUE=%d", AsInt(got))
	}
	r.thrMu.Lock()
	p, _ := peerOfObject(contender.peer)
	pending := p != nil && p.pendingInterrupt
	r.thrMu.Unlock()
	if !pending {
		t.Fatal("class initialization wait consumed the competing thread's interrupt")
	}
}

func TestGate4ClassInitializationFailureOutcomes(t *testing.T) {
	// Owner receives an ExceptionInInitializerError wrapping the original
	// non-Error throwable.
	var original *Object
	nonError := aotMethod("app/Raises", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			return ThrowOutcome(original)
		})
	r := gate4Runtime(t, &ClassDef{Name: "app/Raises", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{nonError}})
	original = r.NewThrowable(nil, "java/lang/RuntimeException")
	c := r.mustClass("app/Raises")
	th := testThread(r, "owner")
	thrown := th.EnsureInitialized(c)
	if thrown == nil || thrown.Class.Name != "java/lang/ExceptionInInitializerError" {
		t.Fatalf("non-Error initializer outcome = %v", thrown)
	}
	if CauseOf(thrown) != original {
		t.Fatal("initializer wrapper lost the original cause identity")
	}
	first := th.EnsureInitialized(c)
	second := th.EnsureInitialized(c)
	if first == nil || second == nil || first.Class.Name != "java/lang/NoClassDefFoundError" ||
		second.Class.Name != "java/lang/NoClassDefFoundError" || first == second {
		t.Fatal("erroneous class did not yield a fresh NoClassDefFoundError per use")
	}

	// An Error from the initializer propagates unchanged.
	var errorValue *Object
	errorMethod := aotMethod("app/Errors", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			return ThrowOutcome(errorValue)
		})
	r2 := gate4Runtime(t, &ClassDef{Name: "app/Errors", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{errorMethod}})
	errorValue = r2.NewThrowable(nil, "java/lang/Error")
	c2 := r2.mustClass("app/Errors")
	if got := testThread(r2, "owner").EnsureInitialized(c2); got != errorValue {
		t.Fatal("initializer Error identity was replaced")
	}

	// A failed required parent marks the child erroneous and propagates the
	// parent's failure; the child initializer never runs.
	var parentError *Object
	var childRan int32
	parentMethod := aotMethod("app/Parent", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			return ThrowOutcome(parentError)
		})
	childMethod := aotMethod("app/Child", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			atomic.AddInt32(&childRan, 1)
			return Normal()
		})
	r3 := gate4Runtime(t,
		&ClassDef{Name: "app/Parent", Super: "java/lang/Object", Flags: AccPublic,
			Methods: []*MethodDef{parentMethod}},
		&ClassDef{Name: "app/Child", Super: "app/Parent", Flags: AccPublic,
			Methods: []*MethodDef{childMethod}})
	parent := r3.mustClass("app/Parent")
	child := r3.mustClass("app/Child")
	parentError = r3.NewThrowable(nil, "java/lang/RuntimeException")
	// Initialize the child directly: the required parent fails first and the
	// child inherits that failure, then every later use gets a fresh
	// NoClassDefFoundError.
	childThrown := testThread(r3, "child-owner").EnsureInitialized(child)
	if childThrown == nil || childThrown.Class.Name != "java/lang/ExceptionInInitializerError" ||
		CauseOf(childThrown) != parentError {
		t.Fatalf("parent failure was not propagated to the child: %v", childThrown)
	}
	if later := testThread(r3, "later").EnsureInitialized(child); later == nil ||
		later.Class.Name != "java/lang/NoClassDefFoundError" {
		t.Fatalf("later child use = %v, want NoClassDefFoundError", later)
	}
	if thrown := testThread(r3, "parent-later").EnsureInitialized(parent); thrown == nil ||
		thrown.Class.Name != "java/lang/NoClassDefFoundError" {
		t.Fatalf("later parent use = %v, want NoClassDefFoundError", thrown)
	}
	if atomic.LoadInt32(&childRan) != 0 {
		t.Fatal("child initializer ran after its parent failed")
	}
	child.initMu.Lock()
	childState := child.State
	child.initMu.Unlock()
	if childState != classErroneous {
		t.Fatal("child class was not left erroneous")
	}
}

func TestGate4ClassInitializationOrder(t *testing.T) {
	var order []string
	var orderMu sync.Mutex
	appendOrder := func(name string) {
		orderMu.Lock()
		defer orderMu.Unlock()
		order = append(order, name)
	}
	parentInit := aotMethod("app/OrderParent", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			appendOrder("parent")
			return Normal()
		})
	ifaceInit := aotMethod("app/OrderIface", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			appendOrder("iface")
			return Normal()
		})
	ifaceDefault := aotMethod("app/OrderIface", "d", "()V", AccPublic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome { return Normal() })
	childInit := aotMethod("app/OrderChild", "<clinit>", "()V", AccStatic,
		func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
			appendOrder("child")
			return Normal()
		})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/OrderParent", Super: "java/lang/Object", Flags: AccPublic,
			Methods: []*MethodDef{parentInit}},
		&ClassDef{Name: "app/OrderIface", Super: "java/lang/Object",
			Flags:   AccPublic | AccInterface | AccAbstract,
			Methods: []*MethodDef{ifaceInit, ifaceDefault}},
		&ClassDef{Name: "app/OrderChild", Super: "app/OrderParent",
			Interfaces: []string{"app/OrderIface"}, Flags: AccPublic,
			Methods: []*MethodDef{childInit}},
	)
	if thrown := testThread(r, "owner").EnsureInitialized(r.mustClass("app/OrderChild")); thrown != nil {
		t.Fatalf("initialization failed: %v", thrown)
	}
	orderMu.Lock()
	got := strings.Join(order, ",")
	orderMu.Unlock()
	if got != "parent,iface,child" {
		t.Fatalf("initialization order = %q, want parent,iface,child", got)
	}
}

func TestGate4ClassInitializationPerPluginLoader(t *testing.T) {
	clinit := &MethodDef{
		Owner: "ext/Boom", Name: "<clinit>", Descriptor: "()V", Flags: AccStatic,
		Engine: "interp", MaxStack: 1, MaxLocals: 0, CodeSize: 5, EntryIndex: -1,
		Insns: []Insn{
			{Op: OpIconst5, PC: 0, Ref: -1, Len: 1},
			{Op: OpPutstatic, PC: 1, Ref: 0, Len: 3},
			{Op: OpReturn, PC: 4, Ref: -1, Len: 1},
		},
		Refs: []Ref{{Kind: RefField, Owner: "ext/Boom", Name: "COUNT", Descriptor: "I", CP: 1}},
	}
	def := func() *ClassDef {
		return &ClassDef{Name: "ext/Boom", Super: "java/lang/Object", Flags: AccPublic,
			StaticFields: []FieldDef{{Name: "COUNT", Descriptor: "I", Flags: AccStatic}},
			Methods:      []*MethodDef{clinit}}
	}
	r := gate4Runtime(t)
	loaderA, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{def()}, Defined: []string{"ext/Boom"}})
	if thrown != nil {
		t.Fatal(thrown)
	}
	loaderB, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{def()}, Defined: []string{"ext/Boom"}})
	if thrown != nil {
		t.Fatal(thrown)
	}
	classA, _ := loaderA.Local("ext/Boom")
	classB, _ := loaderB.Local("ext/Boom")
	if classA == classB || classA.Mirror == classB.Mirror {
		t.Fatal("identical plugin definitions shared class identity")
	}
	if thrown := testThread(r, "a").EnsureInitialized(classA); thrown != nil {
		t.Fatal(thrown)
	}
	if thrown := testThread(r, "b").EnsureInitialized(classB); thrown != nil {
		t.Fatal(thrown)
	}
	fieldA := staticField(classA, "COUNT")
	fieldB := staticField(classB, "COUNT")
	if AsInt(r.staticRead(classA, fieldA.Index)) != 5 || AsInt(r.staticRead(classB, fieldB.Index)) != 5 {
		t.Fatal("plugin-local static initialization was not independent")
	}
}

// ------------------------------------------------------- runtime metadata

func TestGate4ResolutionCacheRace(t *testing.T) {
	target := aotMethod("app/Target", "m", "()I", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		return OkValue(IntValue(1))
	})
	r := gate4Runtime(t, &ClassDef{Name: "app/Target", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{target}})
	const goroutines = 32
	results := make([]*Method, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			th := testThread(r, "resolver")
			m, thrown := th.resolveMethodSymbol(nil, 11, "app/Target", "m", "()I")
			if thrown != nil {
				t.Errorf("resolution failed: %v", thrown)
				return
			}
			results[i] = m
		}(i)
	}
	wg.Wait()
	for i := 1; i < goroutines; i++ {
		if results[i] == nil || results[i] != results[0] {
			t.Fatal("concurrent first resolutions produced distinct methods")
		}
	}
	r.metaMu.Lock()
	count := len(r.resolutions)
	r.metaMu.Unlock()
	if count != 1 {
		t.Fatalf("resolution cache holds %d entries, want 1", count)
	}

	// Cached failures are sticky: every caller receives the same throwable.
	first, firstThrown := testThread(r, "r1").resolveMethodSymbol(nil, 12, "missing/Type", "m", "()V")
	second, secondThrown := testThread(r, "r2").resolveMethodSymbol(nil, 12, "missing/Type", "m", "()V")
	if first != nil || second != nil || firstThrown == nil || firstThrown != secondThrown {
		t.Fatal("cached resolution failure was not sticky")
	}
}

func TestGate4DispatchCacheRace(t *testing.T) {
	baseDef := aotMethod("app/Base", "n", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		return Normal()
	})
	subDef := aotMethod("app/Sub", "n", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		return Normal()
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Base", Super: "java/lang/Object", Flags: AccPublic, Methods: []*MethodDef{baseDef}},
		&ClassDef{Name: "app/Sub", Super: "app/Base", Flags: AccPublic, Methods: []*MethodDef{subDef}},
	)
	base := r.mustClass("app/Base")
	sub := r.mustClass("app/Sub")
	decl := base.MethodIndex["n()V"]
	override := sub.MethodIndex["n()V"]
	caller := base.MethodIndex["n()V"]
	const goroutines = 32
	var wg sync.WaitGroup
	mismatch := make(chan *Method, 1)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			th := testThread(r, "dispatcher")
			target, hit, thrown := th.selectGuarded(caller, 0, OpInvokevirtual, sub, decl)
			if thrown != nil {
				t.Errorf("dispatch failed: %v", thrown)
				return
			}
			if target != override {
				select {
				case mismatch <- target:
				default:
				}
			}
			_ = hit
		}()
	}
	wg.Wait()
	select {
	case target := <-mismatch:
		t.Fatalf("concurrent dispatch selected %s.%s", target.Owner.Name, target.Name)
	default:
	}
	r.metaMu.Lock()
	count := len(r.dispatch)
	r.metaMu.Unlock()
	if count != 1 {
		t.Fatalf("dispatch cache holds %d entries, want 1", count)
	}
}

func TestGate4MetadataInterningConcurrency(t *testing.T) {
	r := gate4Runtime(t, &ClassDef{Name: "app/Comp", Super: "java/lang/Object", Flags: AccPublic})
	comp := r.mustClass("app/Comp")
	units := StringToUnits("catty")
	const goroutines = 32
	arrays := make([]*Class, goroutines)
	stringsSeen := make([]*Object, goroutines)
	boxes := make([]*Object, goroutines)
	mirrors := make([]*Object, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			arrays[i] = r.arrayClassForComponent(comp)
			stringsSeen[i] = r.NewStringUnits(units, true)
			boxes[i] = r.BoxedInteger(42)
			mirrors[i] = comp.Mirror
		}(i)
	}
	wg.Wait()
	for i := 1; i < goroutines; i++ {
		if arrays[i] != arrays[0] {
			t.Fatal("array classes were interned inconsistently under concurrency")
		}
		if stringsSeen[i] != stringsSeen[0] {
			t.Fatal("interned strings lost canonical identity under concurrency")
		}
		if boxes[i] != boxes[0] {
			t.Fatal("boxed cache lost canonical identity under concurrency")
		}
		if mirrors[i] != mirrors[0] {
			t.Fatal("class mirrors lost canonical identity")
		}
	}
}

func TestGate4ConcurrentAdmissionQuotaAndRollback(t *testing.T) {
	r := gate4Runtime(t)
	limits := *Gate3PluginLimits
	limits.Loaders = 1
	r.Image.Limits = &limits
	gate := make(chan struct{})
	results := make(chan *Object, 2)
	for i := 0; i < 2; i++ {
		name := []string{"ext/One", "ext/Two"}[i]
		go func(name string) {
			<-gate
			loader, thrown := r.installPlugin(r.Loader, &AdmissionResult{
				Classes: []*ClassDef{{Name: name, Super: "java/lang/Object", Flags: AccPublic}},
				Defined: []string{name}, ClassBytes: 8,
			})
			if loader != nil {
				results <- nil
				return
			}
			results <- thrown
		}(name)
	}
	close(gate)
	var successes, limitFailures int
	for i := 0; i < 2; i++ {
		thrown := <-results
		if thrown == nil {
			successes++
			continue
		}
		if data := ThrowableOf(thrown); data != nil && data.Code == "LIMIT" {
			limitFailures++
		} else {
			t.Fatalf("quota race rejection = %v", thrown)
		}
	}
	if successes != 1 || limitFailures != 1 {
		t.Fatalf("quota race = %d successes, %d limit failures", successes, limitFailures)
	}
	if got := len(r.pluginLoaders); got != 1 {
		t.Fatalf("published loaders = %d, want 1", got)
	}

	// A rejected candidate must not undo a concurrent successful admission.
	limits.Loaders = 2
	r.Image.Limits = &limits
	valid := make(chan error, 1)
	invalid := make(chan error, 1)
	go func() {
		_, thrown := r.installPlugin(r.Loader, &AdmissionResult{
			Classes: []*ClassDef{{Name: "ext/Good", Super: "java/lang/Object", Flags: AccPublic}},
			Defined: []string{"ext/Good"}, ClassBytes: 4,
		})
		if thrown != nil {
			valid <- &PluginError{Code: "VERIFY", Message: "valid candidate rejected"}
			return
		}
		valid <- nil
	}()
	go func() {
		_, thrown := r.installPlugin(r.Loader, &AdmissionResult{
			Classes: []*ClassDef{{Name: "ext/Bad", Super: "missing/Super", Flags: AccPublic}},
			Defined: []string{"ext/Bad"}, ClassBytes: 9,
		})
		if thrown == nil {
			invalid <- nil
			return
		}
		data := ThrowableOf(thrown)
		invalid <- errorFromCode(data)
	}()
	if err := <-valid; err != nil {
		t.Fatalf("valid admission failed: %v", err)
	}
	if err := <-invalid; err == nil || (err.Error() != "DEPENDENCY" && err.Error() != "LIMIT") {
		t.Fatalf("invalid admission code = %v", err)
	}
	r.metaMu.Lock()
	loaders := len(r.pluginLoaders)
	bytes := r.pluginClassBytes
	id := r.nextLoaderID
	r.metaMu.Unlock()
	if loaders != 2 || bytes != 12 {
		t.Fatalf("rollback accounting = %d loaders, %d bytes", loaders, bytes)
	}
	if id != 3 {
		t.Fatalf("next loader ID = %d, want 3", id)
	}
}

type codeError string

func (e codeError) Error() string { return string(e) }

func errorFromCode(data *ThrowableData) error {
	if data == nil {
		return codeError("VERIFY")
	}
	return codeError(data.Code)
}

func TestGate4AdmissionSnapshotDetached(t *testing.T) {
	r := gate4Runtime(t)
	original := []byte{1, 2, 3, 4, 5}
	captured := make(chan []byte, 1)
	release := make(chan struct{})
	r.SetAdmission(func(req *AdmissionRequest) (*AdmissionResult, error) {
		captured <- append([]byte(nil), req.Bytes...)
		<-release
		return nil, &PluginError{Code: "PROFILE", Message: "test stop"}
	})
	th := testThread(r, "registrar")
	byteClass := r.arrayClassForComponent(r.PrimitiveClass('B'))
	arr := r.NewArrayOfClass(byteClass, len(original))
	for i, b := range original {
		arr.Array.Bytes[i] = int8(b)
	}
	plugins := r.mustClass("catty/runtime/Plugins")
	register := plugins.MethodIndex["register([BLjava/lang/ClassLoader;)Ljava/lang/ClassLoader;"]
	done := make(chan Outcome, 1)
	go func() {
		done <- intrinsicPluginsRegister(th, register, nil, []Value{RefValue(arr), RefValue(nil)})
	}()
	snapshot := <-captured
	for i := range arr.Array.Bytes {
		arr.Array.Bytes[i] = 99
	}
	close(release)
	out := <-done
	if out.Thrown == nil || out.Thrown.Class.Name != "catty/runtime/PluginAdmissionException" {
		t.Fatalf("admission stop outcome = %v", out.Thrown)
	}
	for i := range original {
		if i >= len(snapshot) || snapshot[i] != original[i] {
			t.Fatalf("admission snapshot aliased the caller's byte array: %v", snapshot)
		}
	}
}

// --------------------------------------------------------- shared memory

func TestGate4SharedMemorySlotOperations(t *testing.T) {
	class := &ClassDef{
		Name: "app/Mem", Super: "java/lang/Object", Flags: AccPublic,
		StaticFields: []FieldDef{
			{Name: "SI", Descriptor: "I", Flags: AccStatic},
			{Name: "SL", Descriptor: "J", Flags: AccStatic},
			{Name: "SR", Descriptor: "Ljava/lang/Object;", Flags: AccStatic | AccVolatile},
		},
		InstanceFields: []FieldDef{
			{Name: "II", Descriptor: "I", Flags: 0},
			{Name: "IL", Descriptor: "J", Flags: 0},
			{Name: "IR", Descriptor: "Ljava/lang/Object;", Flags: 0},
		},
	}
	r := gate4Runtime(t, class)
	c := r.mustClass("app/Mem")
	receiver := r.NewInstance(c)
	a := r.NewInstance(r.mustClass("java/lang/Object"))
	b := r.NewInstance(r.mustClass("java/lang/Object"))
	si, sl, sr := staticField(c, "SI"), staticField(c, "SL"), staticField(c, "SR")
	ii, il, ir := instanceField(c, "II"), instanceField(c, "IL"), instanceField(c, "IR")

	const iterations = 400
	var writers, readers sync.WaitGroup
	stop := make(chan struct{})
	for g := 0; g < 4; g++ {
		writers.Add(1)
		go func(g int) {
			defer writers.Done()
			for i := 0; i < iterations; i++ {
				if (i+g)&1 == 0 {
					r.staticWrite(c, si.Index, IntValue(11))
					r.staticWrite(c, sl.Index, LongValue(0x1122334455667788))
					r.staticWrite(c, sr.Index, RefValue(a))
					r.fieldWrite(receiver, ii.Offset, IntValue(22))
					r.fieldWrite(receiver, il.Offset, LongValue(-0x1122334455667788))
					r.fieldWrite(receiver, ir.Offset, RefValue(b))
				} else {
					r.staticWrite(c, si.Index, IntValue(33))
					r.staticWrite(c, sl.Index, LongValue(0x0102030405060708))
					r.staticWrite(c, sr.Index, RefValue(b))
					r.fieldWrite(receiver, ii.Offset, IntValue(44))
					r.fieldWrite(receiver, il.Offset, LongValue(0x0102030405060708))
					r.fieldWrite(receiver, ir.Offset, RefValue(a))
				}
			}
		}(g)
	}
	for g := 0; g < 4; g++ {
		readers.Add(1)
		go func() {
			defer readers.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				if v := AsInt(r.staticRead(c, si.Index)); v != 0 && v != 11 && v != 33 {
					t.Errorf("torn static int: %d", v)
					return
				}
				switch AsLong(r.staticRead(c, sl.Index)) {
				case 0, 0x1122334455667788, 0x0102030405060708:
				default:
					t.Errorf("torn static long: %x", AsLong(r.staticRead(c, sl.Index)))
					return
				}
				if ref := r.staticRead(c, sr.Index).Ref; ref != nil && ref != a && ref != b {
					t.Errorf("invalid static reference value")
					return
				}
				if v := AsInt(r.fieldRead(receiver, ii.Offset)); v != 0 && v != 22 && v != 44 {
					t.Errorf("torn instance int: %d", v)
					return
				}
				switch AsLong(r.fieldRead(receiver, il.Offset)) {
				case 0, int64(-0x1122334455667788), 0x0102030405060708:
				default:
					t.Errorf("torn instance long: %x", AsLong(r.fieldRead(receiver, il.Offset)))
					return
				}
				if ref := r.fieldRead(receiver, ir.Offset).Ref; ref != nil && ref != a && ref != b {
					t.Errorf("invalid instance reference value")
					return
				}
			}
		}()
	}
	writers.Wait()
	close(stop)
	readers.Wait()

	// Array backing kinds through the same memory path.
	ints := r.NewArrayOfClass(r.arrayClassForComponent(r.PrimitiveClass('I')), 8)
	longs := r.NewArrayOfClass(r.arrayClassForComponent(r.PrimitiveClass('J')), 8)
	refs := r.NewArrayOfClass(r.arrayClassForComponent(r.mustClass("java/lang/Object")), 8)
	var awg sync.WaitGroup
	for g := 0; g < 4; g++ {
		awg.Add(1)
		go func(g int) {
			defer awg.Done()
			for i := 0; i < iterations; i++ {
				idx := i % 8
				if (i+g)&1 == 0 {
					r.arrayWrite(ints.Array, idx, IntValue(7))
					r.arrayWrite(longs.Array, idx, LongValue(0x0102030405060708))
					r.arrayWrite(refs.Array, idx, RefValue(a))
				} else {
					r.arrayWrite(ints.Array, idx, IntValue(9))
					r.arrayWrite(longs.Array, idx, LongValue(0x1122334455667788))
					r.arrayWrite(refs.Array, idx, RefValue(b))
				}
				for j := 0; j < 8; j++ {
					if v := AsInt(r.arrayRead(ints.Array, j)); v != 0 && v != 7 && v != 9 {
						t.Errorf("torn int array element: %d", v)
						return
					}
					switch AsLong(r.arrayRead(longs.Array, j)) {
					case 0, 0x0102030405060708, 0x1122334455667788:
					default:
						t.Errorf("torn long array element")
						return
					}
					if ref := r.arrayRead(refs.Array, j).Ref; ref != nil && ref != a && ref != b {
						t.Errorf("invalid reference array element")
						return
					}
				}
			}
		}(g)
	}
	awg.Wait()
}

func TestGate4StartJoinPublication(t *testing.T) {
	pub := aotMethod("app/Pub", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		c := th.RT.mustClass("app/Pub")
		th.RT.staticWrite(c, staticField(c, "VALUE").Index, IntValue(99))
		return Normal()
	})
	def := &ClassDef{Name: "app/Pub", Super: "java/lang/Object", Flags: AccPublic,
		StaticFields: []FieldDef{{Name: "VALUE", Descriptor: "I", Flags: AccStatic}},
		Interfaces:   []string{"java/lang/Runnable"},
		Methods:      []*MethodDef{pub}}
	r := gate4Runtime(t, def)
	main := testThread(r, "main")
	c := r.mustClass("app/Pub")
	peer := r.newPeerObject("publisher", false, threadNew)
	peerData(peer).runnable = r.NewInstance(c)
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(main, start, peer, nil); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if out := main.joinPeer(peer, 0, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if got := AsInt(r.staticRead(c, staticField(c, "VALUE").Index)); got != 99 {
		t.Fatalf("join did not publish the worker write: %d", got)
	}
}

func TestGate4TraceCloseDuringDaemonActivity(t *testing.T) {
	stop := make(chan struct{})
	task := aotMethod("app/Noisy", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		obj := receiver
		for {
			select {
			case <-stop:
				return Normal()
			default:
			}
			if thrown := th.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
				return ThrowOutcome(thrown)
			}
			if thrown := th.MonitorExit(nil, 1, RefValue(obj)); thrown != nil {
				return ThrowOutcome(thrown)
			}
		}
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Noisy", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{task}},
	)
	tracePath := t.TempDir() + "/trace.jsonl"
	if err := r.SetTracePath(tracePath); err != nil {
		t.Fatal(err)
	}
	peer := r.newPeerObject("daemon", true, threadNew)
	peerData(peer).runnable = r.NewInstance(r.mustClass("app/Noisy"))
	main := testThread(r, "main")
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(main, start, peer, nil); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	waitFor(t, "trace monitor activity", func() bool {
		data, err := os.ReadFile(tracePath)
		return err == nil && strings.Contains(string(data), `"event":"monitor"`)
	})
	// Close while the daemon is still emitting; the trace drops records
	// instead of racing the file.
	r.CloseTrace()
	time.Sleep(5 * time.Millisecond)
	close(stop)
	r.WaitForNonDaemon()
	events := readEvents(t, tracePath)
	if len(events) == 0 {
		t.Fatal("trace contains no records")
	}
	for _, event := range events {
		if event["event"] == "monitor" && event["phase"] == "exit" {
			return
		}
	}
	t.Fatal("trace close lost every monitor record")
}

func TestGate4UncaughtWorkerDiagnosticAndTrace(t *testing.T) {
	task := aotMethod("app/Failing", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		return ThrowOutcome(th.RT.NewThrowable(th, "java/lang/RuntimeException"))
	})
	r := gate4Runtime(t,
		&ClassDef{Name: "app/Failing", Super: "java/lang/Object", Flags: AccPublic,
			Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{task}},
	)
	tracePath := t.TempDir() + "/trace.jsonl"
	if err := r.SetTracePath(tracePath); err != nil {
		t.Fatal(err)
	}
	peer := r.newPeerObject("worker", false, threadNew)
	peerData(peer).runnable = r.NewInstance(r.mustClass("app/Failing"))
	main := testThread(r, "main")
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	oldStderr := os.Stderr
	os.Stderr = writer
	out := intrinsicThreadStart(main, start, peer, nil)
	if out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	r.WaitForNonDaemon()
	r.CloseTrace()
	os.Stderr = oldStderr
	_ = writer.Close()
	diagnostic, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	text := string(diagnostic)
	if !strings.Contains(text, `Exception in thread "worker" java.lang.RuntimeException`) ||
		!strings.Contains(text, "app.Failing.run") {
		t.Fatalf("uncaught worker diagnostic = %q", text)
	}
	events := readEvents(t, tracePath)
	starts := eventNames(events, "thread_start")
	ends := eventNames(events, "thread_end")
	var workerStart, workerEnd map[string]any
	for _, event := range starts {
		if event["name"] == "worker" {
			workerStart = event
		}
	}
	for _, event := range ends {
		if event["result"] == "uncaught" {
			workerEnd = event
		}
	}
	if workerStart == nil || workerEnd == nil {
		t.Fatalf("missing worker lifecycle evidence: %v / %v", workerStart, workerEnd)
	}
	if workerStart["daemon"] != false || workerEnd["throwable"] == nil {
		t.Fatalf("worker trace = %v / %v", workerStart, workerEnd)
	}
}

// -------------------------------------------------------- profile and image

func TestGate4ProfileAndImageValidation(t *testing.T) {
	base := func(abi int, profile string, plugins, threads bool) *Image {
		return &Image{ABI: abi, Profile: profile, Plugins: plugins, Threads: threads, Limits: Gate3PluginLimits}
	}
	valid := []*Image{
		base(Gate2ABI, Gate2ProfileID, false, false),
		base(Gate3ABI, Gate3ProfileID, true, false),
		base(Gate4ABI, Gate4ProfileID, true, true),
	}
	for i, img := range valid {
		r := NewRuntime(img)
		if err := r.Install(); err != nil {
			t.Fatalf("valid pairing %d rejected: %v", i, err)
		}
	}
	invalid := []struct {
		img  *Image
		want string
	}{
		{base(Gate3ABI, Gate4ProfileID, true, true), "ABI"},
		{base(Gate4ABI, Gate3ProfileID, true, false), "ABI"},
		{base(Gate3ABI, Gate3ProfileID, true, true), "threading"},
		{base(Gate4ABI, Gate4ProfileID, false, true), "plugin"},
		{base(Gate4ABI, Gate4ProfileID, true, false), "threading"},
	}
	for i, item := range invalid {
		r := NewRuntime(item.img)
		err := r.Install()
		if err == nil || !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(item.want)) {
			t.Fatalf("invalid pairing %d error = %v, want %q", i, err, item.want)
		}
		if r.Installed() {
			t.Fatalf("invalid pairing %d reported installed", i)
		}
	}

	// Canonical opcode and capability inventories are profile scoped.
	gate3Opcodes := SupportedOpcodeNamesFor(Gate3ProfileID)
	gate4Opcodes := SupportedOpcodeNamesFor(Gate4ProfileID)
	if containsName(gate3Opcodes, "monitorenter") || containsName(gate3Opcodes, "monitorexit") {
		t.Fatal("Gate 3 inventory admits monitor opcodes")
	}
	if !containsName(gate4Opcodes, "monitorenter") || !containsName(gate4Opcodes, "monitorexit") {
		t.Fatal("Gate 4 inventory omits monitor opcodes")
	}
	if !containsName(ProfileCapabilities(Gate4ProfileID), "synchronized-memory") {
		t.Fatal("Gate 4 capabilities omit synchronized-memory")
	}

	// A declared inventory may not name an opcode the profile rejects, even
	// when the image carries no method that uses it.
	badInventory := base(Gate2ABI, Gate2ProfileID, false, false)
	badInventory.Opcodes = []string{"monitorenter"}
	if err := NewRuntime(badInventory).Install(); err == nil ||
		!strings.Contains(err.Error(), "unsupported opcode") {
		t.Fatalf("Gate 2 image admitted a monitor inventory: %v", err)
	}
	goodInventory := base(Gate4ABI, Gate4ProfileID, true, true)
	goodInventory.Opcodes = gate4Opcodes
	goodInventory.Capabilities = ProfileCapabilities(Gate4ProfileID)
	if err := NewRuntime(goodInventory).Install(); err != nil {
		t.Fatalf("canonical Gate 4 inventory rejected: %v", err)
	}
	badCapability := base(Gate2ABI, Gate2ProfileID, false, false)
	badCapability.Capabilities = []string{"threads"}
	if err := NewRuntime(badCapability).Install(); err == nil ||
		!strings.Contains(err.Error(), "capability") {
		t.Fatalf("Gate 2 image admitted a Gate 4 capability: %v", err)
	}
}

func TestGate4RetainedThreadingValidation(t *testing.T) {
	monitorMethod := func(flags uint16) *MethodDef {
		return &MethodDef{
			Owner: "app/M", Name: "m", Descriptor: "()V", Flags: flags, Engine: "interp",
			MaxStack: 1, MaxLocals: 0, CodeSize: 2, EntryIndex: -1,
			Insns: []Insn{
				{Op: OpMonitorenter, PC: 0, Ref: -1, Len: 1},
				{Op: OpReturn, PC: 1, Ref: -1, Len: 1},
			},
		}
	}
	gate2 := &Image{ABI: Gate2ABI, Profile: Gate2ProfileID, Classes: []*ClassDef{{
		Name: "app/M", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{monitorMethod(AccPublic)},
	}}}
	if err := NewRuntime(gate2).Install(); err == nil || !strings.Contains(err.Error(), "monitorenter") {
		t.Fatalf("Gate 2 retained monitor instruction error = %v", err)
	}

	// Earlier profiles reject synchronized method flags even in a crafted
	// image whose inventory was omitted.
	longBody := &MethodDef{Owner: "app/S", Name: "m", Descriptor: "()V",
		Flags: AccPublic | AccStatic | AccSynchronized, Engine: "interp",
		MaxStack: 1, MaxLocals: 0, CodeSize: 1, EntryIndex: -1,
		Insns: []Insn{{Op: OpReturn, PC: 0, Ref: -1, Len: 1}}}
	gate3 := &Image{ABI: Gate3ABI, Profile: Gate3ProfileID, Plugins: true, Limits: Gate3PluginLimits,
		Classes: []*ClassDef{{Name: "app/S", Super: "java/lang/Object", Flags: AccPublic,
			Methods: []*MethodDef{longBody}}}}
	if err := NewRuntime(gate3).Install(); err == nil || !strings.Contains(err.Error(), "synchronized") {
		t.Fatalf("Gate 3 synchronized flag error = %v", err)
	}

	// Gate 4 accepts a valid synchronized static method and still rejects the
	// illegal flag combinations.
	gate4 := func(md *MethodDef) *Image {
		return &Image{ABI: Gate4ABI, Profile: Gate4ProfileID, Plugins: true, Threads: true, Limits: Gate3PluginLimits,
			Classes: []*ClassDef{{Name: "app/S", Super: "java/lang/Object", Flags: AccPublic,
				Methods: []*MethodDef{md}}}}
	}
	if err := NewRuntime(gate4(longBody)).Install(); err != nil {
		t.Fatalf("valid synchronized static method rejected: %v", err)
	}
	cases := []struct {
		name string
		md   *MethodDef
	}{
		{"constructor", &MethodDef{Owner: "app/S", Name: "<init>", Descriptor: "()V",
			Flags: AccPublic | AccSynchronized, Engine: "abstract"}},
		{"abstract", &MethodDef{Owner: "app/S", Name: "a", Descriptor: "()V",
			Flags: AccPublic | AccAbstract | AccSynchronized, Engine: "abstract"}},
		{"native", &MethodDef{Owner: "app/S", Name: "n", Descriptor: "()V",
			Flags: AccPublic | AccNative | AccSynchronized, Engine: "abstract"}},
	}
	for _, item := range cases {
		if err := NewRuntime(gate4(item.md)).Install(); err == nil ||
			!strings.Contains(err.Error(), "synchronized") {
			t.Fatalf("%s synchronized flag error = %v", item.name, err)
		}
	}
	// A synchronized interface method is illegal; the interface class flag
	// makes the check apply.
	iface := &Image{ABI: Gate4ABI, Profile: Gate4ProfileID, Plugins: true, Threads: true, Limits: Gate3PluginLimits,
		Classes: []*ClassDef{{Name: "app/I", Super: "java/lang/Object",
			Flags: AccPublic | AccInterface | AccAbstract,
			Methods: []*MethodDef{{Owner: "app/I", Name: "m", Descriptor: "()V",
				Flags: AccPublic | AccSynchronized, Engine: "interp",
				MaxStack: 1, MaxLocals: 1, CodeSize: 1, EntryIndex: -1,
				Insns: []Insn{{Op: OpReturn, PC: 0, Ref: -1, Len: 1}}}}}}}
	if err := NewRuntime(iface).Install(); err == nil || !strings.Contains(err.Error(), "interface method") {
		t.Fatalf("interface synchronized method error = %v", err)
	}
}

// TestGate4StaticLockIdentityPerLoader proves that a synchronized static method
// locks the canonical mirror of its defining class, so the same binary name
// under two plugin loaders uses two distinct locks.
func TestGate4StaticLockIdentityPerLoader(t *testing.T) {
	method := func() *MethodDef {
		return &MethodDef{
			Owner: "ext/Sync", Name: "probe", Descriptor: "()Z", Flags: AccPublic | AccStatic | AccSynchronized,
			Engine: "interp", MaxStack: 1, MaxLocals: 0, CodeSize: 4, EntryIndex: -1,
			Insns: []Insn{
				{Op: OpLdc, PC: 0, Ref: 0, Len: 2},
				{Op: OpInvokestatic, PC: 2, Ref: 1, Len: 3},
				{Op: OpIreturn, PC: 5, Ref: -1, Len: 1},
			},
			Refs: []Ref{
				{Kind: RefClass, Class: "ext/Sync", CP: 1},
				{Kind: RefMethod, Owner: "java/lang/Thread", Name: "holdsLock",
					Descriptor: "(Ljava/lang/Object;)Z", CP: 2},
			},
		}
	}
	// Correct the code size to match the emitted instruction offsets.
	m := method()
	m.CodeSize = 6
	def := func() *ClassDef {
		return &ClassDef{Name: "ext/Sync", Super: "java/lang/Object", Flags: AccPublic,
			Methods: []*MethodDef{method()}}
	}
	r := gate4Runtime(t)
	var classes []*Class
	for i := 0; i < 2; i++ {
		loader, thrown := r.installPlugin(r.Loader, &AdmissionResult{Classes: []*ClassDef{def()}, Defined: []string{"ext/Sync"}})
		if thrown != nil {
			t.Fatal(thrown)
		}
		class, ok := loader.Local("ext/Sync")
		if !ok {
			t.Fatal("plugin class missing")
		}
		classes = append(classes, class)
	}
	if classes[0] == classes[1] || classes[0].Mirror == classes[1].Mirror {
		t.Fatal("identical plugin definitions shared class identity")
	}
	th := testThread(r, "main")
	for i, class := range classes {
		out := th.Call(class.MethodIndex["probe()Z"], nil, nil)
		if out.Thrown != nil {
			t.Fatalf("loader %d probe: %v", i, out.Thrown)
		}
		if AsInt(out.Value) != 1 {
			t.Fatalf("loader %d synchronized static method locked the wrong mirror", i)
		}
		if class.Mirror.monitor != nil && class.Mirror.monitor.depth != 0 {
			t.Fatalf("loader %d leaked its static monitor", i)
		}
	}
}

func TestGate4AbstractAndMissingMethodErrors(t *testing.T) {
	abstract := &MethodDef{Owner: "app/Abs", Name: "a", Descriptor: "()V",
		Flags: AccPublic | AccAbstract, Engine: "abstract"}
	r := gate4Runtime(t, &ClassDef{Name: "app/Abs", Super: "java/lang/Object", Flags: AccPublic,
		Methods: []*MethodDef{abstract}})
	th := testThread(r, "main")
	c := r.mustClass("app/Abs")
	receiver := r.NewInstance(c)
	if out := th.Call(c.MethodIndex["a()V"], receiver, nil); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/AbstractMethodError" {
		t.Fatalf("abstract method call outcome = %v", out.Thrown)
	}
	if th.holdsLockObject(receiver) {
		t.Fatal("abstract method call acquired the receiver monitor")
	}
	if out := th.Call(nil, receiver, nil); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/NullPointerException" {
		t.Fatalf("nil method call outcome = %v", out.Thrown)
	}
	if th.holdsLockObject(receiver) {
		t.Fatal("nil method call acquired the receiver monitor")
	}
}

// TestGate4ArrayBackingKinds exercises every array backing kind through the
// shared helpers used by both engines.
func TestGate4ArrayBackingKinds(t *testing.T) {
	r := gate4Runtime(t)
	object := r.mustClass("java/lang/Object")
	refA := r.NewInstance(object)
	refB := r.NewInstance(object)
	components := []byte{'Z', 'B', 'C', 'S', 'I', 'J', 'F', 'D'}
	for _, kind := range components {
		arr := r.NewArrayOfClass(r.arrayClassForComponent(r.PrimitiveClass(kind)), 2)
		var first, second Value
		switch kind {
		case 'Z':
			first, second = IntValue(1), IntValue(0)
		case 'B':
			first, second = IntValue(-7), IntValue(7)
		case 'C':
			first, second = IntValue('a'), IntValue('z')
		case 'S':
			first, second = IntValue(-300), IntValue(300)
		case 'I':
			first, second = IntValue(-70000), IntValue(70000)
		case 'J':
			first, second = LongValue(-0x1122334455667788), LongValue(0x1122334455667788)
		case 'F':
			first, second = FloatValue(1.5), FloatValue(-2.5)
		case 'D':
			first, second = DoubleValue(1.5e300), DoubleValue(-2.5e-300)
		}
		r.arrayWrite(arr.Array, 0, first)
		r.arrayWrite(arr.Array, 1, second)
		if got := r.arrayRead(arr.Array, 0); got.Kind != first.Kind || got.Bits != first.Bits {
			t.Fatalf("array kind %c first element = %+v, want %+v", kind, got, first)
		}
		if got := r.arrayRead(arr.Array, 1); got.Kind != second.Kind || got.Bits != second.Bits {
			t.Fatalf("array kind %c second element = %+v, want %+v", kind, got, second)
		}
	}
	refs := r.NewArrayOfClass(r.arrayClassForComponent(object), 2)
	r.arrayWrite(refs.Array, 0, RefValue(refA))
	r.arrayWrite(refs.Array, 1, RefValue(refB))
	if r.arrayRead(refs.Array, 0).Ref != refA || r.arrayRead(refs.Array, 1).Ref != refB {
		t.Fatal("reference array elements lost identity")
	}
}

// TestGate4ReflectionSnapshotDetached verifies that the reflection argument
// snapshot does not alias a concurrently mutated Java array.
func TestGate4ReflectionSnapshotDetached(t *testing.T) {
	r := gate4Runtime(t)
	arr := r.NewArrayOfClass(r.arrayClassForComponent(r.mustClass("java/lang/Object")), 4)
	a := r.NewInstance(r.mustClass("java/lang/Object"))
	b := r.NewInstance(r.mustClass("java/lang/Object"))
	arr.Array.Elems[0], arr.Array.Elems[1] = RefValue(a), RefValue(b)
	var wg sync.WaitGroup
	stop := make(chan struct{})
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			r.arrayWrite(arr.Array, 0, RefValue(b))
			r.arrayWrite(arr.Array, 0, RefValue(a))
		}
	}()
	for i := 0; i < 200; i++ {
		snapshot := r.arraySnapshotElems(arr.Array)
		for j, value := range snapshot {
			if value.Ref != nil && value.Ref != a && value.Ref != b {
				t.Fatalf("snapshot slot %d holds an invalid reference", j)
			}
		}
	}
	close(stop)
	wg.Wait()
}

// TestGate4TerminationReleasesMonitors covers JVM thread termination: a dying
// thread releases every monitor it still owns, even when hand-authored code
// never exited it.
func TestGate4TerminationReleasesMonitors(t *testing.T) {
	r := gate4Runtime(t)
	obj := newTestObject(r)
	task := aotMethod("app/Leaker", "run", "()V", AccPublic, func(th *Thread, m *Method, receiver *Object, args []Value) Outcome {
		if thrown := th.MonitorEnter(nil, 0, RefValue(obj)); thrown != nil {
			return ThrowOutcome(thrown)
		}
		return Normal()
	})
	def := &ClassDef{Name: "app/Leaker", Super: "java/lang/Object", Flags: AccPublic,
		Interfaces: []string{"java/lang/Runnable"}, Methods: []*MethodDef{task}}
	// Rebuild with the class because the runtime image is immutable.
	r = gate4Runtime(t, def)
	peer := r.newPeerObject("leaker", false, threadNew)
	peerData(peer).runnable = r.NewInstance(r.mustClass("app/Leaker"))
	main := testThread(r, "main")
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(main, start, peer, nil); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if out := main.joinPeer(peer, 0, false); out.Thrown != nil {
		t.Fatal(out.Thrown)
	}
	if owner := monitorOwner(r, obj); owner != nil || monitorDepth(r, obj) != 0 {
		t.Fatalf("terminated thread retained a monitor: owner=%v depth=%d", owner, monitorDepth(r, obj))
	}
	contender := testThread(r, "contender")
	if thrown := contender.MonitorEnter(nil, 1, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
	if thrown := contender.MonitorExit(nil, 2, RefValue(obj)); thrown != nil {
		t.Fatal(thrown)
	}
}

// TestGate4MonitorDoesNotRootUnrelatedObjects proves that a held monitor keeps
// only its own object reachable, not other objects.
func TestGate4MonitorDoesNotRootUnrelatedObjects(t *testing.T) {
	r := gate4Runtime(t)
	owner := testThread(r, "owner")
	monitored := newTestObject(r)
	if thrown := owner.MonitorEnter(nil, 0, RefValue(monitored)); thrown != nil {
		t.Fatal(thrown)
	}
	finalized := make(chan struct{}, 1)
	func() {
		other := newTestObject(r)
		runtime.SetFinalizer(other, func(*Object) { finalized <- struct{}{} })
	}()
	deadline := time.Now().Add(5 * time.Second)
	for {
		runtime.GC()
		select {
		case <-finalized:
			if monitorOwner(r, monitored) != owner {
				t.Fatal("collecting an unrelated object released the held monitor")
			}
			if thrown := owner.MonitorExit(nil, 1, RefValue(monitored)); thrown != nil {
				t.Fatal(thrown)
			}
			return
		case <-time.After(5 * time.Millisecond):
			if time.Now().After(deadline) {
				t.Fatal("monitor state rooted an unrelated object")
			}
		}
	}
}

// gate4PrepFixture mirrors the Gate 3 preparation fixture under the Gate 4
// profile: two published plugin loaders have already resolved the same binary
// name to different classes.
type gate4PrepFixture struct {
	r     *Runtime
	l1    *Loader
	l2    *Loader
	base  *Class
	iface *Class
	x1    *Class
	x2    *Class
}

func newGate4PrepFixture(t *testing.T) *gate4PrepFixture {
	t.Helper()
	r := gate4Runtime(t)
	object := r.mustClass("java/lang/Object")
	l2 := &Loader{ID: 2, Name: "prep-root", Runtime: r, Parent: r.Loader, Classes: map[string]*Class{}}
	l1 := &Loader{ID: 3, Name: "prep-base", Runtime: r, Parent: l2, Classes: map[string]*Class{}}
	r.nextLoaderID = 3
	r.pluginLoaders = append(r.pluginLoaders, l2, l1)
	base := testHiddenClass(l1, "ext/Base", AccPublic, object)
	iface := testHiddenClass(l2, "ext/IFace", AccPublic|AccInterface|AccAbstract, object)
	testHiddenMethod(base, "m", "(Lext/N;)V", AccPublic)
	testHiddenMethod(iface, "m", "(Lext/N;)V", AccPublic|AccAbstract)
	x1 := testHiddenClass(l1, "ext/N", AccPublic, object)
	x2 := testHiddenClass(l2, "ext/N", AccPublic, object)
	if r.recordInitiation(l1, "ext/N", x1) || r.recordInitiation(l2, "ext/N", x2) {
		t.Fatal("initial ext/N initiations conflicted")
	}
	return &gate4PrepFixture{r: r, l1: l1, l2: l2, base: base, iface: iface, x1: x1, x2: x2}
}

// TestGate4AdmissionRollbackAfterMetadataStaging rejects a candidate only after
// its classes, methods and array metadata have been staged, then verifies the
// successful transaction state is intact and the candidate did not publish a
// loader, consume a loader ID, or mutate constraints and quota accounting.
func TestGate4AdmissionRollbackAfterMetadataStaging(t *testing.T) {
	f := newGate4PrepFixture(t)
	beforeInitiations := cloneInitiations(f.r.initiations)
	beforeConstraints := cloneConstraints(f.r.constraints)
	loadersBefore := len(f.r.pluginLoaders)
	bytesBefore := f.r.pluginClassBytes
	idBefore := f.r.nextLoaderID
	identityBefore := f.r.nextIdentity.Load()

	def := &ClassDef{
		Name: "ext/Child", Super: "ext/Base", Interfaces: []string{"ext/IFace"},
		InstanceFields: []FieldDef{{Name: "value", Descriptor: "[Lext/N;", Flags: 0}},
	}
	loader, thrown := f.r.installPlugin(f.l1, &AdmissionResult{Classes: []*ClassDef{def}, Defined: []string{"ext/Child"}, ClassBytes: 64})
	if loader != nil || thrown == nil {
		t.Fatal("conflicting candidate was published")
	}
	if data := ThrowableOf(thrown); data == nil || data.Code != "VERIFY" {
		t.Fatalf("wrong diagnostic for inherited implementation conflict: %#v", data)
	}
	if len(f.r.pluginLoaders) != loadersBefore || f.r.nextLoaderID != idBefore {
		t.Fatal("rejected candidate changed the loader registry or ID counter")
	}
	if f.r.pluginClassBytes != bytesBefore {
		t.Fatalf("rejected candidate changed quota accounting: %d -> %d", bytesBefore, f.r.pluginClassBytes)
	}
	if !sameInitiations(f.r.initiations, beforeInitiations) {
		t.Fatal("rejected candidate changed initiation records")
	}
	if !sameConstraints(f.r.constraints, beforeConstraints) {
		t.Fatal("rejected candidate published staged loading constraints")
	}
	if f.r.nextIdentity.Load() < identityBefore {
		t.Fatal("rejected candidate rewound the object identity counter")
	}
	// A later valid candidate still admits normally.
	valid := &ClassDef{Name: "ext/Fine", Super: "java/lang/Object", Flags: AccPublic}
	loader, thrown = f.r.installPlugin(f.l1, &AdmissionResult{Classes: []*ClassDef{valid}, Defined: []string{"ext/Fine"}, ClassBytes: 8})
	if thrown != nil || loader == nil {
		t.Fatalf("valid admission after rejection failed: %v", thrown)
	}
	if len(f.r.pluginLoaders) != loadersBefore+1 {
		t.Fatal("valid admission after rejection did not publish exactly one loader")
	}
}
