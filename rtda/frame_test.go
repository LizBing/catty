package rtda

import (
	"testing"
)

// testLoader is a minimal Loader used in rtda-internal tests.
type testLoader struct {
	classes map[string]*Class
	id      *LoaderIdentity
}

func (l *testLoader) LoadClass(name string) *Class {
	return l.classes[name]
}

func (l *testLoader) LoadClassResult(name string) ClassLoadResult {
	c := l.LoadClass(name)
	if c != nil {
		return NewClassResult(c)
	}
	return NewFailureResult(&ClassLoadFailure{Kind: FailureNotFound, Name: name})
}

func (l *testLoader) LoaderIdentity() *LoaderIdentity {
	if l.id == nil {
		l.id = NewLoaderIdentity()
	}
	return l.id
}

// TestFrameLongRoundTrip guards the two-slot long encoding: the interpreter's
// long/double support (loads/stores, fields, returns) all depend on PushLong /
// PopLong splitting and rejoining a 64-bit value across two adjacent slots.
func TestFrameLongRoundTrip(t *testing.T) {
	m := &Method{maxStack: 4, maxLocals: 4}
	f := NewFrame(nil, m)

	for _, v := range []int64{0, 1, -1, 1<<40, -(1 << 40), 0x123456789ABCDEF0} {
		f.PushLong(v)
		if got := f.PopLong(); got != v {
			t.Errorf("PopLong = %#x, want %#x", got, v)
		}
	}
}

// TestFrameDoubleRoundTrip exercises the same path via the double view.
func TestFrameDoubleRoundTrip(t *testing.T) {
	m := &Method{maxStack: 4, maxLocals: 4}
	f := NewFrame(nil, m)

	for _, v := range []float64{0, 1.5, -2.25, 3.14159, 1e300, -1e-300} {
		f.PushDouble(v)
		if got := f.PopDouble(); got != v {
			t.Errorf("PopDouble = %v, want %v", got, v)
		}
	}
}

// TestFrameLocalsLong covers SetLong/GetLong (local-variable two-slot access),
// used by lload/lstore and long parameter passing.
func TestFrameLocalsLong(t *testing.T) {
	m := &Method{maxStack: 2, maxLocals: 4}
	f := NewFrame(nil, m)

	v := int64(-0x7777AAAA55550000)
	f.SetLong(1, v)
	if got := f.GetLong(1); got != v {
		t.Errorf("GetLong = %#x, want %#x", got, v)
	}
}

// TestFrameRefGCNil checks that popping a reference clears the slot so the GC
// can reclaim the object — a subtle correctness/perf property of PopRef.
func TestFrameRefGCNil(t *testing.T) {
	m := &Method{maxStack: 2, maxLocals: 1}
	f := NewFrame(nil, m)

	obj := &Object{}
	f.PushRef(obj)
	_ = f.PopRef()
	// The just-popped slot must no longer hold a reference.
	if f.stack[0].ref != nil {
		t.Error("PopRef left a stale reference in the freed slot")
	}
}

// --- ACC_SYNCHRONIZED frame-entry contract (Slice C, ADR-0029) ---

// newSyncMethod creates a virtual method with ACC_SYNCHRONIZED set. The method
// is non-native and has no bytecode; it exists only to test frame-level sync
// monitor entry and exit. maxStack/maxLocals are minimal to avoid OOB.
func newSyncMethod(owner *Class, name, desc string, isStatic bool) *Method {
	flags := accPublic | accSynchronized
	if isStatic {
		flags |= accStatic
	}
	md := ParseMethodDescriptor(desc)
	return &Method{
		owner:        owner,
		name:         name,
		descriptor:   desc,
		accessFlags:  flags,
		argSlotCount: uint(md.ArgSlots()),
		maxLocals:    uint(md.ArgSlots() + 1),
		maxStack:     2,
	}
}

// TestFrameEnterSyncMonitorInstance verifies that when a frame for an
// ACC_SYNCHRONIZED instance method calls EnterSyncMonitor, the monitor on the
// receiver object (local 0) is entered and syncObject is set.
func TestFrameEnterSyncMonitorInstance(t *testing.T) {
	cls := newMinimalClass("TestSync")
	obj := NewObject(cls)
	thr := NewThread(nil) // nil loader is fine for instance methods
	method := newSyncMethod(cls, "syncMethod", "()V", false)
	frame := NewFrame(thr, method)
	frame.SetRef(0, obj) // 'this' in local 0

	frame.EnterSyncMonitor()

	mon := obj.Monitor()
	if !mon.HoldsLock(thr.EC()) {
		t.Fatal("EnterSyncMonitor did not acquire instance monitor")
	}
	if frame.SyncObject() != obj {
		t.Fatal("frame.SyncObject() != receiver after EnterSyncMonitor")
	}

	// Clean up.
	mon.Exit(thr.EC())
}

// TestFrameEnterSyncMonitorStatic verifies that EnterSyncMonitor on a static
// synchronized method enters the monitor on the declaring class's canonical
// Class mirror.
func TestFrameEnterSyncMonitorStatic(t *testing.T) {
	// Build a minimal loader that can resolve java/lang/Class.
	classClass := NewSyntheticClass("java/lang/Class", nil)
	loader := &testLoader{classes: map[string]*Class{
		"java/lang/Class": classClass,
	}}

	cls := newMinimalClass("TestStaticSync")
	method := newSyncMethod(cls, "staticSync", "()V", true)
	thr := NewThread(loader)
	frame := NewFrame(thr, method)

	frame.EnterSyncMonitor()

	// The Class mirror should have been created lazily and its monitor entered.
	mirror := cls.ClassObject(nil)
	if mirror == nil {
		t.Fatal("ClassObject returned nil — Class mirror was not created")
	}
	if !mirror.Monitor().HoldsLock(thr.EC()) {
		t.Fatal("EnterSyncMonitor did not acquire Class mirror monitor for static method")
	}
	if frame.SyncObject() != mirror {
		t.Fatal("frame.SyncObject() != Class mirror after EnterSyncMonitor (static)")
	}

	// Clean up.
	mirror.Monitor().Exit(thr.EC())
}

// TestFrameEnterSyncMonitorNonSync verifies that EnterSyncMonitor is a no-op
// for methods without ACC_SYNCHRONIZED.
func TestFrameEnterSyncMonitorNonSync(t *testing.T) {
	cls := newMinimalClass("TestNonSync")
	obj := NewObject(cls)
	thr := NewThread(nil)
	method := NativeMethod(cls, "plain", "()V", func(f *Frame) {}) // no ACC_SYNCHRONIZED
	frame := NewFrame(thr, method)
	frame.SetRef(0, obj)

	frame.EnterSyncMonitor()

	if frame.SyncObject() != nil {
		t.Fatal("non-synchronized method should not set syncObject")
	}
	if obj.Monitor().HoldsLock(thr.EC()) {
		t.Fatal("non-synchronized method should not acquire monitor")
	}
}

// TestFramePopFrameReleasesSyncMonitor verifies that when a frame with a
// non-nil syncObject is popped, the implicit synchronized monitor is released.
func TestFramePopFrameReleasesSyncMonitor(t *testing.T) {
	cls := newMinimalClass("TestSyncPop")
	obj := NewObject(cls)
	thr := NewThread(nil)
	method := newSyncMethod(cls, "syncMethod", "()V", false)
	frame := NewFrame(thr, method)
	frame.SetRef(0, obj)

	frame.EnterSyncMonitor()
	thr.PushFrame(frame)

	if !obj.Monitor().HoldsLock(thr.EC()) {
		t.Fatal("monitor not held after EnterSyncMonitor + PushFrame")
	}

	thr.PopFrame()

	if obj.Monitor().HoldsLock(thr.EC()) {
		t.Fatal("monitor still held after PopFrame — implicit sync exit missing")
	}
}

// TestFramePopFrameNonSyncDoesNotPanic verifies that PopFrame on a frame
// without a syncObject is a no-op (no nil pointer dereference).
func TestFramePopFrameNonSyncDoesNotPanic(t *testing.T) {
	cls := newMinimalClass("TestPlain")
	thr := NewThread(nil)
	method := NativeMethod(cls, "plain", "()V", func(f *Frame) {})
	frame := NewFrame(thr, method)
	thr.PushFrame(frame)

	// PopFrame must not panic on a non-synchronized frame.
	popped := thr.PopFrame()
	if popped != frame {
		t.Fatal("PopFrame returned wrong frame")
	}
}
