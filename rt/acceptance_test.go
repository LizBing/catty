package rt_test

import (
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"catty/rt"
)

// TestCattyAcceptanceImageABI constructs a minimal image with a startup
// initializer and a main function, installs and launches it as a positive
// control, then verifies that an unsupported ABI revision and an unsupported
// profile ID are rejected before any startup marker runs.
func TestCattyAcceptanceImageABI(t *testing.T) {
	newImage := func(abi int, profile string, startup, main *int32) *rt.Image {
		return &rt.Image{
			ABI:     abi,
			Profile: profile,
			Startups: []rt.Startup{
				func(t *rt.Thread) rt.Outcome {
					atomic.StoreInt32(startup, 1)
					return rt.Normal()
				},
			},
			MainFunc: func(t *rt.Thread, args []rt.Value) rt.Outcome {
				atomic.StoreInt32(main, 1)
				return rt.Normal()
			},
		}
	}

	var startup, main int32
	valid := newImage(rt.ABIRevision, rt.ProfileID, &startup, &main)
	r := rt.NewRuntime(valid)
	if err := r.Install(); err != nil {
		t.Fatalf("valid image install: %v", err)
	}
	out := r.Launch(nil)
	if out.Thrown != nil {
		t.Fatalf("valid image launch threw %v", out.Thrown)
	}
	if startup != 1 || main != 1 {
		t.Fatalf("startup/main markers = %d/%d, want 1/1", startup, main)
	}

	var badStartup, badMain int32
	badABI := newImage(rt.ABIRevision+1, rt.ProfileID, &badStartup, &badMain)
	rBadABI := rt.NewRuntime(badABI)
	err := rBadABI.Install()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "abi") {
		t.Fatalf("unsupported ABI install error = %v, want an ABI-specific error", err)
	}
	if rBadABI.Installed() {
		t.Fatal("runtime reported installed after an unsupported ABI")
	}
	if badStartup != 0 || badMain != 0 {
		t.Fatalf("unsupported ABI ran markers %d/%d, want 0/0", badStartup, badMain)
	}

	var profileStartup, profileMain int32
	badProfile := newImage(rt.ABIRevision, "catty-unsupported-profile", &profileStartup, &profileMain)
	rBadProfile := rt.NewRuntime(badProfile)
	err = rBadProfile.Install()
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "profile") {
		t.Fatalf("unsupported profile install error = %v, want a profile-specific error", err)
	}
	if profileStartup != 0 || profileMain != 0 {
		t.Fatalf("unsupported profile ran markers %d/%d, want 0/0", profileStartup, profileMain)
	}
}

type gcCanary struct {
	value int64
	// A pointer-containing, non-tiny object so Go's small pointer-free
	// allocation batching cannot legally suppress finalization.
	refs [4]*byte
}

// keepCanaryUnreachable installs a finalizer on an otherwise unreachable Go
// object. Its finalization proves the garbage collector and finalizer machinery
// actually ran during the nested compiled call.
//
//go:noinline
func keepCanaryUnreachable(done *int32) {
	canary := &gcCanary{value: 1}
	runtime.SetFinalizer(canary, func(*gcCanary) {
		atomic.StoreInt32(done, 1)
	})
}

// TestCattyAcceptanceInterpreterOperandGC verifies that a Java object living
// only in interpreter operand slots survives a Go garbage collection triggered
// by a nested compiled entry. The test holds no strong reference to the object
// while it is on the operand stack: it records only the object's stable
// identity and installs a finalization probe.
func TestCattyAcceptanceInterpreterOperandGC(t *testing.T) {
	var gcRuns int32
	var finalized int32
	var canaryDone int32
	var observedIdentity uint64

	initEntry := func(th *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
		return rt.Normal()
	}
	gcEntry := func(th *rt.Thread, m *rt.Method, receiver *rt.Object, args []rt.Value) rt.Outcome {
		atomic.AddInt32(&gcRuns, 1)
		keepCanaryUnreachable(&canaryDone)
		if receiver != nil || len(args) != 0 {
			t.Errorf("compiled entry received a reference: receiver=%v args=%d", receiver, len(args))
		}
		for i := 0; i < 50; i++ {
			runtime.GC()
			time.Sleep(time.Millisecond)
			if atomic.LoadInt32(&canaryDone) == 1 {
				break
			}
		}
		return rt.Normal()
	}

	witness := &rt.ClassDef{
		Name:  "Witness",
		Super: "java/lang/Object",
		InstanceFields: []rt.FieldDef{
			{Name: "value", Descriptor: "I"},
		},
		StaticFields: []rt.FieldDef{
			{Name: "seen", Descriptor: "I"},
		},
		Methods: []*rt.MethodDef{
			{
				Owner: "Witness", Name: "<init>", Descriptor: "()V", Flags: rt.AccPublic,
				Engine: "aot", EntryIndex: 0,
			},
			{
				Owner: "Witness", Name: "run", Descriptor: "()LWitness;",
				Flags: rt.AccPublic | rt.AccStatic, Engine: "interp", MaxStack: 3, MaxLocals: 0, CodeSize: 24,
				EntryIndex: -1,
				Refs: []rt.Ref{
					{Kind: rt.RefClass, Class: "Witness"},
					{Kind: rt.RefMethod, Owner: "Witness", Name: "<init>", Descriptor: "()V"},
					{Kind: rt.RefField, Owner: "Witness", Name: "value", Descriptor: "I"},
					{Kind: rt.RefMethod, Owner: "Driver", Name: "gc", Descriptor: "()V"},
					{Kind: rt.RefField, Owner: "Witness", Name: "seen", Descriptor: "I"},
				},
				Insns: []rt.Insn{
					{Op: rt.OpNew, PC: 0, Ref: 0, Len: 3},
					{Op: rt.OpDup, PC: 3, Ref: -1, Len: 1},
					{Op: rt.OpInvokespecial, PC: 4, Ref: 1, Len: 3},
					{Op: rt.OpDup, PC: 7, Ref: -1, Len: 1},
					{Op: rt.OpBipush, PC: 8, A: 42, Ref: -1, Len: 2},
					{Op: rt.OpPutfield, PC: 10, Ref: 2, Len: 3},
					{Op: rt.OpInvokestatic, PC: 13, Ref: 3, Len: 3},
					{Op: rt.OpDup, PC: 16, Ref: -1, Len: 1},
					{Op: rt.OpGetfield, PC: 17, Ref: 2, Len: 3},
					{Op: rt.OpPutstatic, PC: 20, Ref: 4, Len: 3},
					{Op: rt.OpAreturn, PC: 23, Ref: -1, Len: 1},
				},
			},
		},
	}
	driver := &rt.ClassDef{
		Name:  "Driver",
		Super: "java/lang/Object",
		Methods: []*rt.MethodDef{
			{
				Owner: "Driver", Name: "gc", Descriptor: "()V",
				Flags: rt.AccPublic | rt.AccStatic, Engine: "aot", EntryIndex: 1,
			},
		},
	}
	image := &rt.Image{
		ABI:     rt.ABIRevision,
		Profile: rt.ProfileID,
		Classes: []*rt.ClassDef{witness, driver},
	}
	if err := image.AttachEntries([]rt.Entry{initEntry, gcEntry}); err != nil {
		t.Fatalf("attach compiled entries: %v", err)
	}
	r := rt.NewRuntime(image)
	r.AllocHook = func(o *rt.Object) {
		if o.Class == nil || o.Class.Name != "Witness" {
			return
		}
		observedIdentity = o.Identity
		runtime.SetFinalizer(o, func(*rt.Object) {
			atomic.StoreInt32(&finalized, 1)
		})
	}
	if err := r.Install(); err != nil {
		t.Fatalf("install: %v", err)
	}
	method, err := r.ResolveMethod("Witness", "run", "()LWitness;")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	thread := r.NewThread()
	out := r.Call(thread, method, nil, nil)
	if out.Thrown != nil {
		t.Fatalf("retained method threw %v", out.Thrown)
	}
	if atomic.LoadInt32(&gcRuns) != 1 {
		t.Fatalf("compiled entry ran %d times, want 1", gcRuns)
	}
	obj := out.Value.Ref
	if obj == nil {
		t.Fatal("retained method returned Java null")
	}
	if obj.Class == nil || obj.Class.Name != "Witness" {
		t.Fatalf("returned object class = %v, want Witness", obj.Class)
	}
	if got := rt.AsInt(obj.Fields[0]); got != 42 {
		t.Fatalf("field after GC = %d, want 42", got)
	}
	if obj.Identity != observedIdentity {
		t.Fatalf("object identity = %d, want witness %d", obj.Identity, observedIdentity)
	}
	witnessClass, err := r.Class("Witness")
	if err != nil {
		t.Fatalf("class lookup: %v", err)
	}
	if got := rt.AsInt(witnessClass.Statics[0]); got != 42 {
		t.Fatalf("interpreter-read static seen = %d, want 42", got)
	}

	// Confirm the finalizer machinery actually processed unreachable objects
	// during the nested call; otherwise the survival assertion would be weak.
	deadline := time.Now().Add(2 * time.Second)
	for atomic.LoadInt32(&canaryDone) == 0 && time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(time.Millisecond)
	}
	if atomic.LoadInt32(&canaryDone) != 1 {
		t.Fatal("finalization probe did not run; GC/finalizer verification is inconclusive")
	}
	for i := 0; i < 100; i++ {
		if atomic.LoadInt32(&finalized) != 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if atomic.LoadInt32(&finalized) != 0 {
		t.Fatal("operand-stack object was finalized during the nested compiled call")
	}
	// Keep the returned object strongly reachable through the finalization
	// probe; the probe above detects a collection that happened during the
	// nested compiled call.
	runtime.KeepAlive(obj)
	runtime.KeepAlive(out)
}
