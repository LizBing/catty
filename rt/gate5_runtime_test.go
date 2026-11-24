package rt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// ---------------------------------------------------------------- helpers

func gate5Runtime(t *testing.T, classes ...*ClassDef) *Runtime {
	t.Helper()
	img := &Image{
		ABI:       Gate5ABI,
		Profile:   Gate5ProfileID,
		Plugins:   true,
		Threads:   true,
		Lifecycle: true,
		Limits:    Gate3PluginLimits,
		Classes:   classes,
	}
	r := NewRuntime(img)
	if err := r.Install(); err != nil {
		t.Fatalf("install: %v", err)
	}
	return r
}

func gate5RunnableClass(name string, run func(t *Thread)) *ClassDef {
	return &ClassDef{
		Name: name, Super: "java/lang/Object", Flags: AccPublic,
		Interfaces: []string{"java/lang/Runnable"},
		Methods: []*MethodDef{
			aotMethod(name, "run", "()V", AccPublic, func(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
				run(t)
				return Normal()
			}),
		},
	}
}

func gate5Peer(r *Runtime, name string, daemon bool, runnable *Object) *Object {
	peer := r.newPeerObject(name, daemon, threadNew)
	if runnable != nil {
		peerData(peer).runnable = runnable
	}
	return peer
}

func gate5Runnable(r *Runtime, class string) *Object {
	return r.NewInstance(r.mustClass(class))
}

func gate5Register(t *testing.T, r *Runtime, peer *Object) {
	t.Helper()
	caller := testThread(r, "registrar")
	method := r.mustClass("java/lang/Runtime").MethodIndex["addShutdownHook(Ljava/lang/Thread;)V"]
	out := intrinsicRuntimeAddShutdownHook(caller, method, r.runtimeObject, []Value{RefValue(peer)})
	if out.Thrown != nil {
		t.Fatalf("addShutdownHook threw %v", out.Thrown.Class.Name)
	}
}

func gate5Start(t *testing.T, r *Runtime, peer *Object) {
	t.Helper()
	caller := testThread(r, "starter")
	method := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	out := intrinsicThreadStart(caller, method, peer, nil)
	if out.Thrown != nil {
		t.Fatalf("Thread.start threw %v", out.Thrown.Class.Name)
	}
}

func gate5WaitTerminal(t *testing.T, r *Runtime) terminalResult {
	t.Helper()
	select {
	case result := <-r.lifecycleTerminal:
		return result
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for a terminal lifecycle result")
		return terminalResult{}
	}
}

func gate5WaitState(t *testing.T, r *Runtime, state lifecycleState) {
	t.Helper()
	waitFor(t, "lifecycle state", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return r.lifeState == state
	})
}

// gate5WaitShutdownClaimed waits until the running state has been left; a
// snapshot with no hooks may already be terminal.
func gate5WaitShutdownClaimed(t *testing.T, r *Runtime) {
	t.Helper()
	waitFor(t, "shutdown claimed", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return r.lifeState != lifecycleRunning
	})
}

// -------------------------------------------------------- profile boundary

func TestGate5ProfileBoundary(t *testing.T) {
	profile, ok := ProfileByID(Gate5ProfileID)
	if !ok || profile.ABI != Gate5ABI || !profile.Lifecycle || !profile.Threads || !profile.Plugins {
		t.Fatalf("Gate 5 profile = %+v", profile)
	}
	caps := ProfileCapabilities(Gate5ProfileID)
	for _, want := range []string{"application-lifecycle", "shutdown-hooks", "application-diagnostics",
		"threads", "monitors", "plugin-admission"} {
		found := false
		for _, cap := range caps {
			if cap == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Gate 5 capabilities missing %q: %v", want, caps)
		}
	}
	if !ProfileHasThreads(Gate5ProfileID) || !ProfileAdmitsPlugins(Gate5ProfileID) {
		t.Fatal("Gate 5 must inherit threading and plugin admission")
	}
	for _, earlier := range []string{Gate2ProfileID, Gate3ProfileID, Gate4ProfileID} {
		if BootstrapHasMethodFor(earlier, "java/lang/System", "exit", "(I)V") {
			t.Fatalf("%s publishes System.exit", earlier)
		}
		if BootstrapHasClassFor(earlier, "java/lang/Runtime") {
			t.Fatalf("%s publishes java.lang.Runtime", earlier)
		}
	}
	for _, member := range []struct{ owner, name, desc string }{
		{"java/lang/System", "exit", "(I)V"},
		{"java/lang/System", "nanoTime", "()J"},
		{"java/lang/Runtime", "getRuntime", "()Ljava/lang/Runtime;"},
		{"java/lang/Runtime", "exit", "(I)V"},
		{"java/lang/Runtime", "halt", "(I)V"},
		{"java/lang/Runtime", "addShutdownHook", "(Ljava/lang/Thread;)V"},
		{"java/lang/Runtime", "removeShutdownHook", "(Ljava/lang/Thread;)Z"},
		{"java/lang/Runtime", "gc", "()V"},
	} {
		if !BootstrapHasMethodFor(Gate5ProfileID, member.owner, member.name, member.desc) {
			t.Fatalf("Gate 5 bootstrap missing %s.%s%s", member.owner, member.name, member.desc)
		}
	}
	info, ok := BootstrapClassDefFor(Gate5ProfileID, "java/lang/IllegalStateException")
	if !ok || info.Super != "java/lang/RuntimeException" {
		t.Fatalf("java.lang.IllegalStateException = %+v", info)
	}
}

func TestGate5ImageFeatureFlagConsistency(t *testing.T) {
	img := &Image{ABI: Gate5ABI, Profile: Gate5ProfileID, Plugins: true, Threads: true, Lifecycle: true,
		Limits: Gate3PluginLimits}
	if err := NewRuntime(img).Install(); err != nil {
		t.Fatalf("consistent Gate 5 image: %v", err)
	}
	// Omission means false; an older image keeps decoding unchanged.
	for _, p := range []*Profile{Gate2Profile, Gate3Profile, Gate4Profile} {
		old := &Image{ABI: p.ABI, Profile: p.ID, Plugins: p.Plugins, Threads: p.Threads}
		if err := NewRuntime(old).Install(); err != nil {
			t.Fatalf("%s image with omitted lifecycle flag: %v", p.ID, err)
		}
	}
	bad := &Image{ABI: Gate5ABI, Profile: Gate5ProfileID, Plugins: true, Threads: true, Lifecycle: false,
		Limits: Gate3PluginLimits}
	if err := NewRuntime(bad).Install(); err == nil || !strings.Contains(err.Error(), "lifecycle flag") {
		t.Fatalf("mismatched lifecycle flag error = %v", err)
	}
	bad = &Image{ABI: Gate5ABI - 1, Profile: Gate5ProfileID, Plugins: true, Threads: true, Lifecycle: true}
	if err := NewRuntime(bad).Install(); err == nil {
		t.Fatal("mismatched Gate 5 ABI was accepted")
	}
	bad = &Image{ABI: Gate5ABI, Profile: Gate5ProfileID, Plugins: false, Threads: true, Lifecycle: true,
		Limits: Gate3PluginLimits}
	if err := NewRuntime(bad).Install(); err == nil || !strings.Contains(err.Error(), "plugin flag") {
		t.Fatalf("mismatched plugin flag error = %v", err)
	}
}

func TestGate5BootstrapIntrinsics(t *testing.T) {
	r := gate5Runtime(t)
	if r.runtimeObject == nil {
		t.Fatal("Runtime singleton was not installed")
	}
	main := testThread(r, "main")
	getRuntime := r.mustClass("java/lang/Runtime").MethodIndex["getRuntime()Ljava/lang/Runtime;"]
	first := intrinsicRuntimeGetRuntime(main, getRuntime, nil, nil)
	second := intrinsicRuntimeGetRuntime(main, getRuntime, nil, nil)
	if first.Thrown != nil || second.Thrown != nil ||
		first.Value.Ref == nil || first.Value.Ref != second.Value.Ref || first.Value.Ref != r.runtimeObject {
		t.Fatal("Runtime.getRuntime did not return one non-null singleton")
	}
	nano := r.mustClass("java/lang/System").MethodIndex["nanoTime()J"]
	a := intrinsicSystemNanoTime(main, nano, nil, nil)
	time.Sleep(time.Millisecond)
	b := intrinsicSystemNanoTime(main, nano, nil, nil)
	if AsLong(b.Value) <= AsLong(a.Value) {
		t.Fatalf("nanoTime did not advance: %d then %d", AsLong(a.Value), AsLong(b.Value))
	}
}

// ------------------------------------------------------ natural shutdown

func TestGate5NaturalShutdownFirstZero(t *testing.T) {
	var hookRuns, workerRuns int32
	release := make(chan struct{})
	r := gate5Runtime(t,
		gate5RunnableClass("app/Hook", func(*Thread) { atomic.AddInt32(&hookRuns, 1) }),
		gate5RunnableClass("app/Worker", func(*Thread) {
			<-release
			atomic.AddInt32(&workerRuns, 1)
		}),
	)
	hook := gate5Peer(r, "hook", false, gate5Runnable(r, "app/Hook"))
	gate5Register(t, r, hook)
	worker := gate5Peer(r, "worker", false, gate5Runnable(r, "app/Worker"))
	gate5Start(t, r, worker)

	main := r.MainThread()
	r.publishTermination(main, main.peer, false)
	select {
	case result := <-r.lifecycleTerminal:
		t.Fatalf("natural shutdown committed while a non-daemon worker was live: %+v", result)
	case <-time.After(50 * time.Millisecond):
	}
	if atomic.LoadInt32(&hookRuns) != 0 {
		t.Fatal("shutdown hook ran before the first zero count")
	}
	close(release)
	result := gate5WaitTerminal(t, r)
	if result.reason != "natural" || result.status != 0 {
		t.Fatalf("natural result = %+v", result)
	}
	if atomic.LoadInt32(&workerRuns) != 1 || atomic.LoadInt32(&hookRuns) != 1 {
		t.Fatalf("worker/hook runs = %d/%d", workerRuns, hookRuns)
	}
	r.thrMu.Lock()
	defer r.thrMu.Unlock()
	if r.hooks != nil || r.hookSet != nil {
		t.Fatal("hook registry was not released at terminal publication")
	}
	select {
	case extra := <-r.lifecycleTerminal:
		t.Fatalf("second terminal publication: %+v", extra)
	default:
	}
}

func TestGate5FirstZeroRacingNewNonDaemonStart(t *testing.T) {
	var runs int32
	r := gate5Runtime(t, gate5RunnableClass("app/Worker", func(*Thread) { atomic.AddInt32(&runs, 1) }))
	main := r.MainThread()
	worker := gate5Peer(r, "late-worker", false, gate5Runnable(r, "app/Worker"))
	reached := make(chan struct{})
	release := make(chan struct{})
	r.ThreadStartHook = func(*Object) {
		close(reached)
		<-release
	}
	started := make(chan Outcome, 1)
	go func() {
		started <- intrinsicThreadStart(testThread(r, "starter"), r.mustClass("java/lang/Thread").MethodIndex["start()V"], worker, nil)
	}()
	<-reached
	// No non-daemon peer was published yet, so main's retirement commits the
	// natural transition. The delayed start then happens during shutdown and
	// does not extend the already committed drain.
	r.publishTermination(main, main.peer, false)
	gate5WaitShutdownClaimed(t, r)
	close(release)
	if out := <-started; out.Thrown != nil {
		t.Fatalf("late start failed: %v", out.Thrown.Class.Name)
	}
	result := gate5WaitTerminal(t, r)
	if result.reason != "natural" || result.status != 0 {
		t.Fatalf("natural result = %+v", result)
	}
	waitFor(t, "late worker run", func() bool { return atomic.LoadInt32(&runs) == 1 })
}

// ------------------------------------------------------ hook registration

func TestGate5ShutdownHookRegistrationRules(t *testing.T) {
	releaseLive := make(chan struct{})
	r := gate5Runtime(t,
		gate5RunnableClass("app/Hook", func(*Thread) {}),
		gate5RunnableClass("app/Blocking", func(*Thread) { <-releaseLive }),
	)
	caller := testThread(r, "caller")
	add := r.mustClass("java/lang/Runtime").MethodIndex["addShutdownHook(Ljava/lang/Thread;)V"]
	remove := r.mustClass("java/lang/Runtime").MethodIndex["removeShutdownHook(Ljava/lang/Thread;)Z"]
	addValue := func(peer *Object) *Object {
		out := intrinsicRuntimeAddShutdownHook(caller, add, r.runtimeObject, []Value{RefValue(peer)})
		return out.Thrown
	}
	removeValue := func(peer *Object) Outcome {
		return intrinsicRuntimeRemoveShutdownHook(caller, remove, r.runtimeObject, []Value{RefValue(peer)})
	}
	if thrown := addValue(nil); thrown == nil || thrown.Class.Name != "java/lang/NullPointerException" {
		t.Fatalf("add(null) = %v", thrown)
	}
	if thrown := removeValue(nil).Thrown; thrown == nil || thrown.Class.Name != "java/lang/NullPointerException" {
		t.Fatalf("remove(null) = %v", thrown)
	}
	hook := gate5Peer(r, "hook", false, gate5Runnable(r, "app/Hook"))
	if thrown := addValue(hook); thrown != nil {
		t.Fatalf("first add: %v", thrown.Class.Name)
	}
	if thrown := addValue(hook); thrown == nil || thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("duplicate add = %v", thrown)
	}
	if out := removeValue(hook); out.Thrown != nil || AsInt(out.Value) != 1 {
		t.Fatalf("remove registered = %+v", out)
	}
	if out := removeValue(hook); out.Thrown != nil || AsInt(out.Value) != 0 {
		t.Fatalf("remove absent = %+v", out)
	}
	// A live peer cannot be registered, but a retired peer can. Keep one
	// live until after the rejection, then let it retire and register it.
	live := gate5Peer(r, "live", true, gate5Runnable(r, "app/Blocking"))
	gate5Start(t, r, live)
	waitFor(t, "live peer started", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		p, _ := peerOfObject(live)
		return p.state == threadStarted
	})
	if thrown := addValue(live); thrown == nil || thrown.Class.Name != "java/lang/IllegalArgumentException" {
		t.Fatalf("live add = %v", thrown)
	}
	close(releaseLive)
	waitFor(t, "live peer termination", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		p, _ := peerOfObject(live)
		return p.state == threadTerminated
	})
	if thrown := addValue(live); thrown != nil {
		t.Fatalf("retired add = %v", thrown.Class.Name)
	}
	if out := removeValue(live); out.Thrown != nil || AsInt(out.Value) != 1 {
		t.Fatalf("remove retired registered peer = %+v", out)
	}
	// Freeze the registry and require IllegalStateException before argument
	// validation for both operations.
	r.ExitRequestHook = func(int32) {}
	go intrinsicSystemExit(caller, r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(0)})
	gate5WaitShutdownClaimed(t, r)
	if thrown := addValue(nil); thrown == nil || thrown.Class.Name != "java/lang/IllegalStateException" {
		t.Fatalf("add(null) after freeze = %v", thrown)
	}
	if thrown := removeValue(nil).Thrown; thrown == nil || thrown.Class.Name != "java/lang/IllegalStateException" {
		t.Fatalf("remove(null) after freeze = %v", thrown)
	}
	gate5WaitTerminal(t, r)
	if thrown := addValue(live); thrown == nil || thrown.Class.Name != "java/lang/IllegalStateException" {
		t.Fatalf("add after terminal = %v", thrown)
	}
}

func TestGate5RegistrationVersusSnapshot(t *testing.T) {
	var runs int32
	r := gate5Runtime(t, gate5RunnableClass("app/Hook", func(*Thread) { atomic.AddInt32(&runs, 1) }))
	hook := gate5Peer(r, "hook", false, gate5Runnable(r, "app/Hook"))
	gate5Register(t, r, hook)
	caller := testThread(r, "caller")
	remove := r.mustClass("java/lang/Runtime").MethodIndex["removeShutdownHook(Ljava/lang/Thread;)Z"]
	reached := make(chan struct{})
	release := make(chan struct{})
	r.ExitRequestHook = func(int32) {
		close(reached)
		<-release
	}
	go intrinsicSystemExit(caller, r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(4)})
	<-reached
	// Removal wins the race against the snapshot: the hook is absent from the
	// frozen registry and must never run.
	out := intrinsicRuntimeRemoveShutdownHook(caller, remove, r.runtimeObject, []Value{RefValue(hook)})
	if out.Thrown != nil || AsInt(out.Value) != 1 {
		t.Fatalf("remove before snapshot = %+v", out)
	}
	close(release)
	result := gate5WaitTerminal(t, r)
	if result.reason != "exit" || result.status != 4 {
		t.Fatalf("result = %+v", result)
	}
	if atomic.LoadInt32(&runs) != 0 {
		t.Fatal("removed hook executed")
	}

	// Snapshot wins: the frozen registry rejects removal and runs the hook.
	var runs2 int32
	release2 := make(chan struct{})
	r2 := gate5Runtime(t, gate5RunnableClass("app/Hook", func(*Thread) {
		<-release2
		atomic.AddInt32(&runs2, 1)
	}))
	hook2 := gate5Peer(r2, "hook", false, gate5Runnable(r2, "app/Hook"))
	gate5Register(t, r2, hook2)
	go intrinsicSystemExit(testThread(r2, "exit"), r2.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(0)})
	gate5WaitShutdownClaimed(t, r2)
	out2 := intrinsicRuntimeRemoveShutdownHook(testThread(r2, "remover"), remove, r2.runtimeObject, []Value{RefValue(hook2)})
	if out2.Thrown == nil || out2.Thrown.Class.Name != "java/lang/IllegalStateException" {
		t.Fatalf("remove after snapshot = %+v", out2)
	}
	close(release2)
	gate5WaitTerminal(t, r2)
	if atomic.LoadInt32(&runs2) != 1 {
		t.Fatalf("snapshot hook runs = %d", runs2)
	}
}

func TestGate5ManualStartVersusSnapshotStart(t *testing.T) {
	var runs int32
	started := make(chan struct{})
	release := make(chan struct{})
	r := gate5Runtime(t, gate5RunnableClass("app/Hook", func(*Thread) {
		atomic.AddInt32(&runs, 1)
		close(started)
		<-release
	}))
	hook := gate5Peer(r, "hook", false, gate5Runnable(r, "app/Hook"))
	gate5Register(t, r, hook)
	caller := testThread(r, "caller")
	reached := make(chan struct{})
	exitRelease := make(chan struct{})
	r.ExitRequestHook = func(int32) {
		close(reached)
		<-exitRelease
	}
	go intrinsicSystemExit(caller, r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(0)})
	<-reached
	gate5Start(t, r, hook)
	<-started
	// The manual start wins before the registry freezes. The snapshot counts
	// the live execution and must not start the peer a second time.
	close(exitRelease)
	gate5WaitState(t, r, lifecycleShutdown)
	close(release)
	result := gate5WaitTerminal(t, r)
	if result.status != 0 || atomic.LoadInt32(&runs) != 1 {
		t.Fatalf("manual-start result = %+v, runs = %d", result, runs)
	}
	start := r.mustClass("java/lang/Thread").MethodIndex["start()V"]
	if out := intrinsicThreadStart(caller, start, hook, nil); out.Thrown == nil ||
		out.Thrown.Class.Name != "java/lang/IllegalThreadStateException" {
		t.Fatalf("second start after termination = %+v", out.Thrown)
	}
}

// ------------------------------------------------------ exit and halt

func TestGate5SimultaneousExitRequests(t *testing.T) {
	r := gate5Runtime(t)
	caller1 := testThread(r, "exit-3")
	caller2 := testThread(r, "exit-7")
	gate := make(chan struct{})
	systemExit := r.mustClass("java/lang/System").MethodIndex["exit(I)V"]
	go func() {
		<-gate
		intrinsicSystemExit(caller1, systemExit, nil, []Value{IntValue(3)})
	}()
	go func() {
		<-gate
		intrinsicSystemExit(caller2, systemExit, nil, []Value{IntValue(7)})
	}()
	close(gate)
	result := gate5WaitTerminal(t, r)
	if result.status != 3 && result.status != 7 {
		t.Fatalf("selected status = %d, want one of the two requests", result.status)
	}
	if result.requester != caller1.ID && result.requester != caller2.ID {
		t.Fatalf("requester = %d", result.requester)
	}
	select {
	case extra := <-r.lifecycleTerminal:
		t.Fatalf("second terminal result: %+v", extra)
	default:
	}
}

func TestGate5HaltCutsShortGracefulShutdown(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	var hookRuns int32
	r := gate5Runtime(t, gate5RunnableClass("app/Hook", func(*Thread) {
		atomic.AddInt32(&hookRuns, 1)
		close(started)
		<-release
	}))
	hook := gate5Peer(r, "hook", false, gate5Runnable(r, "app/Hook"))
	gate5Register(t, r, hook)
	caller := testThread(r, "exit")
	go intrinsicSystemExit(caller, r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(5)})
	<-started
	halt := r.mustClass("java/lang/Runtime").MethodIndex["halt(I)V"]
	go intrinsicRuntimeHalt(testThread(r, "halter"), halt, r.runtimeObject, []Value{IntValue(19)})
	result := gate5WaitTerminal(t, r)
	if result.reason != "halt" || result.status != 19 {
		t.Fatalf("halt result = %+v", result)
	}
	// The hook that was already running finishes independently; its late
	// terminal request cannot replace the committed halt result.
	close(release)
	waitFor(t, "hook completion", func() bool { return atomic.LoadInt32(&hookRuns) == 1 })
	select {
	case extra := <-r.lifecycleTerminal:
		t.Fatalf("halt was replaced by %+v", extra)
	case <-time.After(50 * time.Millisecond):
	}
}

func TestGate5ExitFromHookStallsUntilHalt(t *testing.T) {
	before := make(chan struct{})
	after := make(chan struct{})
	r := gate5Runtime(t, gate5RunnableClass("app/ExitHook", func(t *Thread) {
		close(before)
		intrinsicSystemExit(t, t.RT.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(9)})
		close(after)
	}))
	hook := gate5Peer(r, "exiting-hook", false, gate5Runnable(r, "app/ExitHook"))
	gate5Register(t, r, hook)
	go intrinsicSystemExit(testThread(r, "main-exit"), r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(0)})
	<-before
	select {
	case result := <-r.lifecycleTerminal:
		t.Fatalf("shutdown completed despite an exit from a hook: %+v", result)
	case <-time.After(100 * time.Millisecond):
	}
	select {
	case <-after:
		t.Fatal("code after exit in a hook executed")
	default:
	}
	halt := r.mustClass("java/lang/Runtime").MethodIndex["halt(I)V"]
	go intrinsicRuntimeHalt(testThread(r, "halter"), halt, r.runtimeObject, []Value{IntValue(21)})
	result := gate5WaitTerminal(t, r)
	if result.reason != "halt" || result.status != 21 {
		t.Fatalf("coordinated halt result = %+v", result)
	}
}

func TestGate5ExitPreservesCallerFramesAndMonitors(t *testing.T) {
	r := gate5Runtime(t)
	caller := testThread(r, "exit-caller")
	guard := newTestObject(r)
	if thrown := caller.MonitorEnter(nil, 0, RefValue(guard)); thrown != nil {
		t.Fatal(thrown)
	}
	var peer *Object
	done := make(chan struct{})
	go func() {
		defer close(done)
		intrinsicSystemExit(caller, r.mustClass("java/lang/System").MethodIndex["exit(I)V"], nil, []Value{IntValue(0)})
	}()
	waitFor(t, "exit caller parked", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		p, _ := peerOfObject(caller.peer)
		return p != nil && p.state == threadStarted && r.lifeState != lifecycleRunning
	})
	// The parked caller remains live, keeps its monitor and has not unwound.
	if !caller.holdsLockObject(guard) {
		t.Fatal("successful exit released the caller's monitor")
	}
	r.thrMu.Lock()
	peer = caller.peer
	r.thrMu.Unlock()
	if frames := caller.FramesSnapshot(); len(frames) != 0 {
		t.Fatalf("host caller unexpectedly had Java frames: %d", len(frames))
	}
	result := gate5WaitTerminal(t, r)
	if result.status != 0 {
		t.Fatalf("result = %+v", result)
	}
	select {
	case <-done:
		t.Fatal("successful exit returned into Go")
	default:
	}
	r.thrMu.Lock()
	if p, _ := peerOfObject(peer); p == nil || p.state != threadStarted {
		r.thrMu.Unlock()
		t.Fatal("parked exit caller was republished as terminated")
	}
	r.thrMu.Unlock()
}

// ------------------------------------------------------------- metrics

func TestGate5MetricsCounters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "metrics.json")
	interp := &MethodDef{
		Owner: "app/Main", Name: "interpreted", Descriptor: "()V", Flags: AccPublic | AccStatic,
		Engine: "interp", MaxStack: 0, MaxLocals: 0, CodeSize: 1,
		Insns: []Insn{{Op: OpReturn, PC: 0, Ref: -1}}, EntryIndex: -1,
	}
	aot := &MethodDef{
		Owner: "app/Main", Name: "compiled", Descriptor: "()V", Flags: AccPublic | AccStatic,
		Engine: "aot", EntryIndex: -1,
		Entry: func(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
			return t.Call(t.RT.mustClass("app/Main").MethodIndex["interpreted()V"], nil, nil)
		},
	}
	thrower := &MethodDef{
		Owner: "app/Main", Name: "boom", Descriptor: "()V", Flags: AccPublic | AccStatic,
		Engine: "aot", EntryIndex: -1,
		Entry: func(t *Thread, m *Method, receiver *Object, args []Value) Outcome {
			return ThrowOutcome(t.RT.NewThrowable(t, "java/lang/RuntimeException"))
		},
	}
	img := &Image{
		ABI: Gate5ABI, Profile: Gate5ProfileID, Plugins: true, Threads: true, Lifecycle: true,
		Limits: Gate3PluginLimits,
		Classes: []*ClassDef{{
			Name: "app/Main", Super: "java/lang/Object", Flags: AccPublic,
			Methods: []*MethodDef{interp, aot, thrower},
		}},
	}
	r := NewRuntime(img)
	if err := r.SetMetricsPath(path); err != nil {
		t.Fatal(err)
	}
	defer r.CloseMetrics()
	if err := r.Install(); err != nil {
		t.Fatal(err)
	}
	main := r.MainThread()
	compiled := r.mustClass("app/Main").MethodIndex["compiled()V"]
	if out := r.Call(main, compiled, nil, nil); out.Thrown != nil {
		t.Fatalf("compiled call: %v", out.Thrown)
	}
	if r.metrics.crossEngine.Load() != 1 {
		t.Fatalf("crossEngineCalls = %d", r.metrics.crossEngine.Load())
	}
	boom := r.mustClass("app/Main").MethodIndex["boom()V"]
	if out := r.Call(main, boom, nil, nil); out.Thrown == nil {
		t.Fatal("boom did not throw")
	}
	if r.metrics.javaExceptions.Load() != 1 {
		t.Fatalf("javaExceptions = %d", r.metrics.javaExceptions.Load())
	}
	// Reflection: one successful invoke and one lookup miss.
	decl := r.mustClass("app/Main").MethodIndex["interpreted()V"]
	methodObj := &Object{Class: r.mustClass("java/lang/reflect/Method"), Payload: decl, Identity: r.newIdentity()}
	invoke := r.mustClass("java/lang/reflect/Method").MethodIndex["invoke(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;"]
	if out := intrinsicMethodInvoke(main, invoke, methodObj, []Value{RefValue(nil), RefValue(nil)}); out.Thrown != nil {
		t.Fatalf("reflective invoke: %v", out.Thrown)
	}
	if r.metrics.reflective.Load() != 1 {
		t.Fatalf("reflectiveInvocations = %d", r.metrics.reflective.Load())
	}
	if _, thrown := r.lookupPublicMethod(r.mustClass("app/Main"), "absent", nil); thrown == nil {
		t.Fatal("missing lookup did not throw")
	}
	if r.metrics.reflectionMisses.Load() != 1 {
		t.Fatalf("reflectionLookupMisses = %d", r.metrics.reflectionMisses.Load())
	}
	// Monitor contention: hold the object, let a second thread wait, release.
	guard := newTestObject(r)
	if thrown := main.MonitorEnter(nil, 0, RefValue(guard)); thrown != nil {
		t.Fatal(thrown)
	}
	waiter := testThread(r, "waiter")
	entered := make(chan struct{}, 1)
	go func() {
		if thrown := waiter.MonitorEnter(nil, 1, RefValue(guard)); thrown == nil {
			entered <- struct{}{}
			waiter.MonitorExit(nil, 2, RefValue(guard))
		}
	}()
	waitFor(t, "waiter parked", func() bool {
		r.thrMu.Lock()
		defer r.thrMu.Unlock()
		return guard.monitor != nil && guard.monitor.contenders > 0
	})
	main.MonitorExit(nil, 3, RefValue(guard))
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		t.Fatal("waiter did not enter the monitor")
	}
	waitFor(t, "monitor contention", func() bool { return r.metrics.contentions.Load() == 1 })
	// Thread publication and explicit GC observation.
	release := make(chan struct{})
	runClass := &ClassDef{Name: "app/Task", Super: "java/lang/Object", Flags: AccPublic,
		Interfaces: []string{"java/lang/Runnable"},
		Methods: []*MethodDef{aotMethod("app/Task", "run", "()V", AccPublic, func(*Thread, *Method, *Object, []Value) Outcome {
			<-release
			return Normal()
		})}}
	// A second runtime is required because classes are immutable after install.
	_ = runClass
	r.metrics.liveThreadStarted()
	r.metrics.liveThreadEnded()
	r.collectGarbage()
	if !r.metrics.postGCHeapSeen.Load() {
		t.Fatal("Runtime.gc did not record the post-GC heap observation")
	}
	terminal := terminalResult{reason: "exit", status: 0, requester: main.ID}
	if err := r.finishMetrics(terminal); err != nil {
		t.Fatalf("finishMetrics: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc metricsDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("metrics JSON: %v\n%s", err, data)
	}
	if doc.SchemaVersion != 1 || doc.Profile != Gate5ProfileID || doc.TerminalReason != "exit" || doc.Status != 0 {
		t.Fatalf("metrics identity = %+v", doc)
	}
	if doc.Counters.ClassesLoaded < int64(len(Gate5Profile.bootstrapInventory())+len(PrimitiveClasses)+1) {
		t.Fatalf("classesLoaded = %d", doc.Counters.ClassesLoaded)
	}
	if doc.Counters.MethodsByEngine["aot"] < 2 || doc.Counters.MethodsByEngine["interp"] != 1 ||
		doc.Counters.MethodsByEngine["intrinsic"] == 0 {
		t.Fatalf("methodsByEngine = %+v", doc.Counters.MethodsByEngine)
	}
	if doc.Counters.CrossEngineCalls != 1 || doc.Counters.ReflectiveInvocations != 1 ||
		doc.Counters.ReflectionLookupMisses != 1 || doc.Counters.JavaExceptions != 1 ||
		doc.Counters.MonitorContentions != 1 {
		t.Fatalf("counters = %+v", doc.Counters)
	}
	if doc.Counters.PeakJavaThreads < 2 || doc.Counters.LiveJavaThreads != 1 {
		t.Fatalf("thread counters = live %d peak %d", doc.Counters.LiveJavaThreads, doc.Counters.PeakJavaThreads)
	}
	if doc.GoRuntime.HeapAllocBytes == 0 || doc.GoRuntime.PostGCHeapBytes == nil ||
		!strings.Contains(doc.GoRuntime.HeapAllocScope, "live Go-managed heap") {
		t.Fatalf("goRuntime = %+v", doc.GoRuntime)
	}
}

// --------------------------------------------------------------- trace

func TestGate5TraceSealsAtTerminal(t *testing.T) {
	path := filepath.Join(t.TempDir(), "trace.jsonl")
	r := gate5Runtime(t)
	if err := r.SetTracePath(path); err != nil {
		t.Fatal(err)
	}
	tr := r.trace
	tr.appStart(Gate5ProfileID, Gate5ABI, 2, 1)
	tr.threadStart(1, 0, false, "main")
	done := make(chan struct{})
	go func() {
		for i := 0; i < 50; i++ {
			tr.monitor(1, uint64(i), "enter", 1, nil, 0, "")
			tr.enter(1, r.mustClass("java/lang/Object").MethodIndex["<init>()V"], 0)
		}
		close(done)
	}()
	tr.processExit("exit", 263, 1)
	<-done
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 3 {
		t.Fatalf("trace lines = %d", len(lines))
	}
	var last map[string]any
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &last); err != nil {
		t.Fatal(err)
	}
	if last["event"] != "process_exit" || last["reason"] != "exit" || last["status"] != float64(263) {
		t.Fatalf("terminal trace record = %v", last)
	}
	for i, line := range lines[:len(lines)-1] {
		var event map[string]any
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			t.Fatalf("line %d: %v", i, err)
		}
		if event["event"] == "process_exit" {
			t.Fatalf("process_exit appeared before the end at line %d", i)
		}
	}
	if tr.enabled() {
		t.Fatal("trace remained enabled after the terminal record")
	}
}
