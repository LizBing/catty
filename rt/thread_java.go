package rt

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Gate 4 extends the Gate 3 bootstrap inventory with java.lang.Thread,
// java.lang.Runnable and the monitor/interrupt throwable types. The inventory
// is assembled once on first use so the Gate 3 inventory is fully initialized
// first.

var (
	gate4Once      sync.Once
	gate4Inventory []*bsClass
	gate4Index     map[string]*bsClass
)

func gate4BootstrapInventory() []*bsClass {
	gate4Once.Do(func() {
		gate4Inventory = buildGate4Inventory()
		gate4Index = map[string]*bsClass{}
		for _, c := range gate4Inventory {
			gate4Index[c.name] = c
		}
	})
	return gate4Inventory
}

func gate4BootstrapByName() map[string]*bsClass {
	gate4BootstrapInventory()
	return gate4Index
}

func buildGate4Inventory() []*bsClass {
	out := make([]*bsClass, 0, len(gate3Inventory)+5)
	merged := map[string]*bsClass{}
	for _, c := range gate3Inventory {
		nc := cloneBootstrapClass(c)
		out = append(out, nc)
		merged[nc.name] = nc
	}
	if object := merged["java/lang/Object"]; object != nil {
		object.methods = append(object.methods,
			bsMethod{"wait", "()V", AccPublic | AccFinal, intrinsicObjectWait},
			bsMethod{"wait", "(J)V", AccPublic | AccFinal, intrinsicObjectWait},
			bsMethod{"wait", "(JI)V", AccPublic | AccFinal, intrinsicObjectWait},
			bsMethod{"notify", "()V", AccPublic | AccFinal, intrinsicObjectNotify},
			bsMethod{"notifyAll", "()V", AccPublic | AccFinal, intrinsicObjectNotify},
		)
	}
	merged["java/lang/Thread"] = &bsClass{
		name:       "java/lang/Thread",
		super:      "java/lang/Object",
		interfaces: []string{"java/lang/Runnable"},
		methods: []bsMethod{
			{"<init>", "()V", AccPublic, intrinsicThreadInit},
			{"<init>", "(Ljava/lang/Runnable;)V", AccPublic, intrinsicThreadInit},
			{"<init>", "(Ljava/lang/String;)V", AccPublic, intrinsicThreadInit},
			{"<init>", "(Ljava/lang/Runnable;Ljava/lang/String;)V", AccPublic, intrinsicThreadInit},
			{"start", "()V", AccPublic, intrinsicThreadStart},
			{"run", "()V", AccPublic, intrinsicThreadRun},
			{"interrupt", "()V", AccPublic, intrinsicThreadInterrupt},
			{"isInterrupted", "()Z", AccPublic, intrinsicThreadIsInterrupted},
			{"join", "()V", AccPublic | AccFinal, intrinsicThreadJoin},
			{"join", "(J)V", AccPublic | AccFinal, intrinsicThreadJoin},
			{"join", "(JI)V", AccPublic | AccFinal, intrinsicThreadJoin},
			{"isAlive", "()Z", AccPublic | AccFinal, intrinsicThreadIsAlive},
			{"setDaemon", "(Z)V", AccPublic | AccFinal, intrinsicThreadSetDaemon},
			{"isDaemon", "()Z", AccPublic | AccFinal, intrinsicThreadIsDaemon},
			{"getName", "()Ljava/lang/String;", AccPublic | AccFinal, intrinsicThreadGetName},
			{"threadId", "()J", AccPublic | AccFinal, intrinsicThreadID},
			{"currentThread", "()Ljava/lang/Thread;", AccPublic | AccStatic, intrinsicThreadCurrent},
			{"interrupted", "()Z", AccPublic | AccStatic, intrinsicThreadInterrupted},
			{"holdsLock", "(Ljava/lang/Object;)Z", AccPublic | AccStatic, intrinsicThreadHoldsLock},
			{"sleep", "(J)V", AccPublic | AccStatic, intrinsicThreadSleep},
			{"sleep", "(JI)V", AccPublic | AccStatic, intrinsicThreadSleep},
			{"yield", "()V", AccPublic | AccStatic, intrinsicThreadYield},
		},
	}
	merged["java/lang/Runnable"] = &bsClass{
		name:  "java/lang/Runnable",
		super: "java/lang/Object",
		flags: AccPublic | AccInterface | AccAbstract,
		methods: []bsMethod{
			{"run", "()V", AccPublic | AccAbstract | AccInterface, nil},
		},
	}
	merged["java/lang/InterruptedException"] = &bsClass{
		name: "java/lang/InterruptedException", super: "java/lang/Exception", methods: throwableInit(),
	}
	merged["java/lang/IllegalThreadStateException"] = &bsClass{
		name: "java/lang/IllegalThreadStateException", super: "java/lang/IllegalArgumentException", methods: throwableInit(),
	}
	merged["java/lang/IllegalMonitorStateException"] = &bsClass{
		name: "java/lang/IllegalMonitorStateException", super: "java/lang/RuntimeException", methods: throwableInit(),
	}
	order := []string{
		"java/lang/Thread",
		"java/lang/Runnable",
		"java/lang/InterruptedException",
		"java/lang/IllegalThreadStateException",
		"java/lang/IllegalMonitorStateException",
	}
	for _, name := range order {
		out = append(out, merged[name])
	}
	return out
}

// maxTimerDuration is the largest representable monotonic wait segment. A
// positive Java timer that exceeds it waits in bounded segments.
const maxTimerDuration = time.Duration(math.MaxInt64)

// millisToDuration converts a nonnegative Java millisecond count, saturating at
// the representable timer ceiling.
func millisToDuration(ms int64) time.Duration {
	if ms <= 0 {
		return 0
	}
	if ms > int64(maxTimerDuration/time.Millisecond) {
		return maxTimerDuration
	}
	return time.Duration(ms) * time.Millisecond
}

// waitTimeout converts the two-argument wait/join forms into one duration. A
// zero remainder stays indefinite (timed=false); a positive nanosecond remainder
// rounds up to one millisecond with saturation.
func waitTimeout(ms int64, nanos int32) (time.Duration, bool) {
	if ms == 0 && nanos == 0 {
		return 0, false
	}
	d := millisToDuration(ms)
	if nanos > 0 {
		if d > maxTimerDuration-time.Millisecond {
			d = maxTimerDuration
		} else {
			d += time.Millisecond
		}
	}
	if d <= 0 {
		d = time.Millisecond
	}
	return d, true
}

// sleepNanos sums a sleep duration in nanoseconds, saturating at the signed
// 64-bit nanosecond ceiling.
func sleepNanos(ms int64, nanos int32) time.Duration {
	const maxMillis = int64(math.MaxInt64) / int64(time.Millisecond)
	if ms > maxMillis {
		return maxTimerDuration
	}
	total := ms * int64(time.Millisecond)
	if total > int64(math.MaxInt64)-int64(nanos) {
		return maxTimerDuration
	}
	return time.Duration(total + int64(nanos))
}

// newPeerObject creates one java.lang.Thread object with peer state. The caller
// publishes it under thrMu.
func (r *Runtime) newPeerObject(name string, daemon bool, state threadState) *Object {
	class := r.mustClass("java/lang/Thread")
	return &Object{
		Class:    class,
		Payload:  &threadPeer{id: atomic.AddInt64(&r.nextThread, 1), name: name, daemon: daemon, state: state},
		Identity: r.newIdentity(),
	}
}

// ------------------------------------------------------------ Object methods

func intrinsicObjectWait(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	switch m.Descriptor {
	case "()V":
		return t.objectWait(receiver, 0, false)
	case "(J)V":
		ms := AsLong(args[0])
		if ms < 0 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		return t.objectWait(receiver, millisToDuration(ms), ms > 0)
	case "(JI)V":
		ms, nanos := AsLong(args[0]), AsInt(args[1])
		if ms < 0 || nanos < 0 || nanos > 999999 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		d, timed := waitTimeout(ms, nanos)
		return t.objectWait(receiver, d, timed)
	}
	return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
}

func intrinsicObjectNotify(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	return t.notifyObject(receiver, m.Name == "notifyAll")
}

// ------------------------------------------------------------ Thread methods

// intrinsicThreadInit creates the peer state exactly once. A null explicit name
// throws NullPointerException. The daemon flag is inherited from the
// constructing logical thread.
func intrinsicThreadInit(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	var runnable, nameArg *Object
	switch m.Descriptor {
	case "()V":
	case "(Ljava/lang/Runnable;)V":
		runnable = args[0].Ref
	case "(Ljava/lang/String;)V":
		nameArg = args[0].Ref
		if nameArg == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
		}
	case "(Ljava/lang/Runnable;Ljava/lang/String;)V":
		runnable = args[0].Ref
		nameArg = args[1].Ref
		if nameArg == nil {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
		}
	default:
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
	}
	r := t.RT
	r.thrMu.Lock()
	if _, ok := peerOfObject(receiver); ok {
		r.thrMu.Unlock()
		return Normal()
	}
	name := ""
	if nameArg != nil {
		if text, ok := StringText(nameArg); ok {
			name = text
		} else {
			r.thrMu.Unlock()
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/ClassCastException"))
		}
	}
	daemon := false
	if constructing, ok := peerOfObject(t.peer); ok {
		daemon = constructing.daemon
	}
	if nameArg == nil {
		name = fmt.Sprintf("Thread-%d", atomic.AddInt64(&r.unnamedPeers, 1))
	}
	p := &threadPeer{
		id:       atomic.AddInt64(&r.nextThread, 1),
		name:     name,
		daemon:   daemon,
		state:    threadNew,
		runnable: runnable,
	}
	receiver.Payload = p
	r.thrMu.Unlock()
	return Normal()
}

// intrinsicThreadStart performs the single new-to-started transition, publishes
// identity, daemon status and launcher accounting, and only then schedules the
// worker goroutine. A second start throws IllegalThreadStateException even
// after termination.
func intrinsicThreadStart(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	if !ok || p.state != threadNew {
		r.thrMu.Unlock()
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalThreadStateException"))
	}
	p.state = threadStarted
	if !p.daemon {
		r.liveNonDaemon++
	}
	// Create the logical thread and publish thread_start before start returns,
	// so the parent's own termination cannot precede its child's first trace
	// event.
	worker := r.NewThread()
	p.thread = worker
	worker.peer = receiver
	name, daemon := p.name, p.daemon
	parent := t.ID
	r.thrMu.Unlock()
	if tr := r.trace; tr != nil {
		tr.threadStart(worker.ID, parent, daemon, name)
	}
	go r.runPeer(p, receiver, worker)
	return Normal()
}

// runPeer runs the virtual run()V selection with the peer as receiver on the
// already-published worker logical thread, then publishes termination.
func (r *Runtime) runPeer(p *threadPeer, peer *Object, t *Thread) {
	var out Outcome
	decl := r.mustClass("java/lang/Thread").MethodIndex["run()V"]
	target := r.selectVirtual(peer.Class, decl)
	if target == nil {
		out = ThrowOutcome(r.NewThrowable(t, "java/lang/AbstractMethodError"))
	} else {
		out = t.Call(target, peer, nil)
	}
	r.finishThread(t, peer, out)
}

// finishThread publishes termination and notifies joiners before the launcher
// counts the thread as completed. A worker's uncaught Java throwable is reported
// with that thread's Java name and frames.
func (r *Runtime) finishThread(t *Thread, peer *Object, out Outcome) {
	result := "terminated"
	if out.Thrown != nil {
		result = "uncaught"
		r.printUncaughtFor(peerName(peer), out.Thrown)
	}
	if tr := r.trace; tr != nil {
		tr.threadEnd(t.ID, result, out.Thrown)
	}
	r.thrMu.Lock()
	r.releaseHeldMonitorsLocked(t)
	p, _ := peerOfObject(peer)
	if p != nil {
		p.state = threadTerminated
		p.thread = nil
		if peer.monitor != nil {
			for _, w := range append([]*waiter(nil), peer.monitor.waiters...) {
				r.wakeWaiterLocked(peer.monitor, w, waiterNotified)
			}
		}
		if !p.daemon {
			r.liveNonDaemon--
			r.liveCond.Broadcast()
		}
	}
	r.thrMu.Unlock()
}

func peerName(peer *Object) string {
	r := peerData(peer)
	if r == nil {
		return ""
	}
	return r.name
}

// WaitForNonDaemon waits until every started non-daemon platform thread has
// terminated. Threads started by a terminating worker still count because the
// start increment happens before that worker's decrement.
func (r *Runtime) WaitForNonDaemon() {
	r.thrMu.Lock()
	for r.liveNonDaemon > 0 {
		r.liveCond.Wait()
	}
	r.thrMu.Unlock()
}

// LiveNonDaemon reports the current non-daemon worker count (test hook).
func (r *Runtime) LiveNonDaemon() int {
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	return r.liveNonDaemon
}

func intrinsicThreadRun(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	var runnable *Object
	if ok {
		runnable = p.runnable
	}
	r.thrMu.Unlock()
	if !ok || runnable == nil {
		return Normal()
	}
	decl := r.mustClass("java/lang/Runnable").MethodIndex["run()V"]
	target := r.selectVirtual(runnable.Class, decl)
	if target == nil {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/AbstractMethodError"))
	}
	return t.Call(target, runnable, nil)
}

func intrinsicThreadInterrupt(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	t.interruptPeer(receiver)
	return Normal()
}

func intrinsicThreadIsInterrupted(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	interrupted := ok && p.pendingInterrupt
	r.thrMu.Unlock()
	if interrupted {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicThreadJoin(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	switch m.Descriptor {
	case "()V":
		return t.joinPeer(receiver, 0, false)
	case "(J)V":
		ms := AsLong(args[0])
		if ms < 0 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		return t.joinPeer(receiver, millisToDuration(ms), ms > 0)
	case "(JI)V":
		ms, nanos := AsLong(args[0]), AsInt(args[1])
		if ms < 0 || nanos < 0 || nanos > 999999 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		d, timed := waitTimeout(ms, nanos)
		return t.joinPeer(receiver, d, timed)
	}
	return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
}

func intrinsicThreadIsAlive(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	alive := ok && p.state == threadStarted
	r.thrMu.Unlock()
	if alive {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicThreadSetDaemon(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	daemon := AsInt(args[0]) != 0
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	if !ok {
		r.thrMu.Unlock()
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalThreadStateException"))
	}
	if p.state == threadStarted {
		r.thrMu.Unlock()
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalThreadStateException"))
	}
	p.daemon = daemon
	r.thrMu.Unlock()
	return Normal()
}

func intrinsicThreadIsDaemon(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	p, ok := peerOfObject(receiver)
	daemon := ok && p.daemon
	r.thrMu.Unlock()
	if daemon {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicThreadGetName(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	name := ""
	if p, ok := peerOfObject(receiver); ok {
		name = p.name
	}
	r.thrMu.Unlock()
	return OkValue(RefValue(r.NewStringGo(name)))
}

func intrinsicThreadID(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if receiver == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	r := t.RT
	r.thrMu.Lock()
	var id int64
	if p, ok := peerOfObject(receiver); ok {
		id = p.id
	}
	r.thrMu.Unlock()
	return OkValue(LongValue(id))
}

func intrinsicThreadCurrent(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	return OkValue(RefValue(t.currentPeer()))
}

func intrinsicThreadInterrupted(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	peer := t.currentPeer()
	r := t.RT
	r.thrMu.Lock()
	interrupted := false
	if p, ok := peerOfObject(peer); ok {
		interrupted = p.pendingInterrupt
		p.pendingInterrupt = false
	}
	r.thrMu.Unlock()
	if interrupted {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicThreadHoldsLock(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	if args[0].Kind != KindRef || args[0].Ref == nil {
		return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/NullPointerException"))
	}
	if t.holdsLockObject(args[0].Ref) {
		return OkValue(IntValue(1))
	}
	return OkValue(IntValue(0))
}

func intrinsicThreadSleep(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	switch m.Descriptor {
	case "(J)V":
		ms := AsLong(args[0])
		if ms < 0 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		return t.sleepFor(millisToDuration(ms))
	case "(JI)V":
		ms, nanos := AsLong(args[0]), AsInt(args[1])
		if ms < 0 || nanos < 0 || nanos > 999999 {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
		}
		return t.sleepFor(sleepNanos(ms, nanos))
	}
	return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/IllegalArgumentException"))
}

func intrinsicThreadYield(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
	runtime.Gosched()
	return Normal()
}
