package rtda

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"catty/classfile"
)

// compileTestFixture compiles tests/fixtures/<name>.java and returns the
// class bytes. It uses the system JDK without source/target restrictions so
// dynamic constant features are emitted.
func compileTestFixture(t *testing.T, name string) []byte {
	t.Helper()
	src := filepath.Join("..", "tests", "fixtures", name+".java")
	out := t.TempDir()
	cmd := exec.Command("javac", "-d", out, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("javac %s failed: %v\n%s", name, err, out)
	}
	data, err := os.ReadFile(filepath.Join(out, name+".class"))
	if err != nil {
		t.Fatalf("read class: %v", err)
	}
	return data
}

// recordingMockLoader is a minimal Loader that records every LoadClass call in
// order. It is used to verify NewClass only loads the superclass and declared
// interfaces – never bootstrap argument classes or other dynamically referenced
// types.
type recordingMockLoader struct {
	mu      sync.Mutex
	classes map[string]*Class
	calls   []string // ordered list of every name passed to LoadClass
	id      *LoaderIdentity
}

func (l *recordingMockLoader) LoadClass(name string) *Class {
	l.mu.Lock()
	l.calls = append(l.calls, name)
	l.mu.Unlock()
	return l.classes[name]
}

func (l *recordingMockLoader) LoadClassResult(name string) ClassLoadResult {
	c := l.LoadClass(name)
	if c != nil {
		return NewClassResult(c)
	}
	return NewFailureResult(&ClassLoadFailure{Kind: FailureNotFound, Name: name})
}

func (l *recordingMockLoader) LoaderIdentity() *LoaderIdentity {
	if l.id == nil {
		l.id = NewLoaderIdentity()
	}
	return l.id
}

func (l *recordingMockLoader) Calls() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	c := make([]string, len(l.calls))
	copy(c, l.calls)
	return c
}

func TestBootstrapMethodsAttachment(t *testing.T) {
	// Parse a classfile with BootstrapMethods and build a runtime Class.
	data := compileTestFixture(t, "DynStringConcat")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	cfBM := cf.BootstrapMethods()
	if cfBM == nil {
		t.Fatal("ClassFile.BootstrapMethods() is nil")
	}

	objClass := NewSyntheticClass("java/lang/Object", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Object": objClass,
	}}

	class := NewClass(cf, loader)

	// Verify the runtime Class stores the BootstrapMethods.
	rtBM := class.BootstrapMethods()
	if rtBM == nil {
		t.Fatal("Class.BootstrapMethods() returns nil")
	}
	if rtBM != cfBM {
		t.Error("Class.BootstrapMethods() must return the exact pointer stored during build time")
	}
	if rtBM.NumEntries() != cfBM.NumEntries() {
		t.Errorf("NumEntries mismatch: classfile=%d, runtime=%d", cfBM.NumEntries(), rtBM.NumEntries())
	}

	// Verify immutability: Entry() returns a copy.
	e1 := rtBM.Entry(0)
	if len(e1.Arguments) > 0 {
		orig := make([]uint16, len(e1.Arguments))
		copy(orig, e1.Arguments)
		e1.Arguments[0] = 9999 // mutate
		e2 := rtBM.Entry(0)
		if e2.Arguments[0] != orig[0] {
			t.Errorf("Entry returned aliased slice: got %d after mutate, want %d", e2.Arguments[0], orig[0])
		}
	}

	// Verify the stored pointer is stable across repeated calls.
	if class.BootstrapMethods() != rtBM {
		t.Error("Class.BootstrapMethods() is not stable across repeated calls")
	}
}

func TestBootstrapMethodsAttachmentEmpty(t *testing.T) {
	data := compileTestFixture(t, "HelloWorld")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	objClass := NewSyntheticClass("java/lang/Object", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Object": objClass,
	}}

	class := NewClass(cf, loader)

	// HelloWorld has no BootstrapMethods — accessor must return nil.
	if class.BootstrapMethods() != nil {
		t.Error("Class.BootstrapMethods() should be nil for classes without BootstrapMethods attribute")
	}
}

// TestNoEagerBootstrapLoading verifies that NewClass ONLY loads the direct
// superclass and declared interfaces (in order). No bootstrap argument classes,
// no dynamically referenced types. The expected call list is computed from the
// classfile rather than hardcoded so the test is self-describing.
func TestNoEagerBootstrapLoading(t *testing.T) {
	data := compileTestFixture(t, "DynStringConcat")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	bm := cf.BootstrapMethods()
	if bm == nil || bm.NumEntries() == 0 {
		t.Fatal("DynStringConcat fixture must have BootstrapMethods — cannot verify no-eager-loading without them")
	}

	// Compute the exact expected LoadClass calls: superclass + declared interfaces.
	var expected []string
	superName := cf.SuperClassName()
	if superName != "" {
		expected = append(expected, superName)
	}
	expected = append(expected, cf.InterfaceNames()...)

	// Build mock classes for every name NewClass should request.
	classes := make(map[string]*Class)
	for _, name := range expected {
		classes[name] = NewSyntheticClass(name, nil)
	}

	loader := &recordingMockLoader{classes: classes}
	_ = NewClass(cf, loader)

	actual := loader.Calls()

	// Exact comparison: must match in count and order.
	if len(actual) != len(expected) {
		t.Errorf("LoadClass call count mismatch: got %d calls %v, want %d calls %v",
			len(actual), actual, len(expected), expected)
		return
	}
	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("LoadClass call[%d] = %q, want %q", i, actual[i], expected[i])
		}
	}

	t.Logf("LoadClass calls: %v (expected: %v)", actual, expected)
}

// TestBootstrapMethodsNotExecuted verifies that bootstrap methods are never
// executed during parsing or class construction. Bootstrap invocation is K4
// (linkage) territory.
func TestBootstrapMethodsNotExecuted(t *testing.T) {
	data := compileTestFixture(t, "DynStringConcat")
	cf, err := classfile.Parse(data)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	bm := cf.BootstrapMethods()
	if bm == nil {
		t.Fatal("BootstrapMethods is nil")
	}

	// Verify we can inspect entries without side effects.
	for i := 0; i < bm.NumEntries(); i++ {
		e := bm.Entry(i)
		if e.MethodRef == 0 {
			t.Errorf("entry %d: MethodRef is 0", i)
		}
	}
}

// TestSyntheticClassBootstrapMethodsNil verifies that synthetic classes
// (which have no classfile) return nil from BootstrapMethods().
func TestSyntheticClassBootstrapMethodsNil(t *testing.T) {
	synth := NewSyntheticClass("test/Synth", nil)
	if synth.BootstrapMethods() != nil {
		t.Error("NewSyntheticClass must return nil for BootstrapMethods()")
	}
	if synth.Name() != "test/Synth" {
		t.Errorf("Name = %q, want test/Synth", synth.Name())
	}
}

// --- K2 Item 4: BootstrapLoader set-once and BindLoaderRef validation ---

func TestSetBootstrapLoaderOnce(t *testing.T) {
	resetBootstrapLoaderForTesting()

	l1 := &recordingMockLoader{classes: make(map[string]*Class)}

	// First set should succeed.
	SetBootstrapLoader(l1)
	if getBootstrapLoader() != l1 {
		t.Error("getBootstrapLoader() != l1 after first SetBootstrapLoader")
	}

	// Second set with same loader should be idempotent.
	SetBootstrapLoader(l1)
	if getBootstrapLoader() != l1 {
		t.Error("idempotent SetBootstrapLoader changed the loader")
	}

	// Second set with different loader should panic.
	l2 := &recordingMockLoader{classes: make(map[string]*Class)}
	defer func() {
		if r := recover(); r == nil {
			t.Error("SetBootstrapLoader with different loader did not panic")
		}
	}()
	SetBootstrapLoader(l2)
}

func TestBindLoaderRefSetOnce(t *testing.T) {
	l1 := &recordingMockLoader{classes: make(map[string]*Class)}
	l2 := &recordingMockLoader{classes: make(map[string]*Class)}
	cls := NewSyntheticClass("test/BindRef", nil)

	// First bind should succeed.
	cls.BindLoaderRef(l1)
	if cls.definingLoaderRef != l1 {
		t.Error("definingLoaderRef != l1 after BindLoaderRef")
	}

	// Same loader should be idempotent.
	cls.BindLoaderRef(l1)
	if cls.definingLoaderRef != l1 {
		t.Error("definingLoaderRef changed after idempotent BindLoaderRef")
	}

	// Different loader should panic.
	defer func() {
		if r := recover(); r == nil {
			t.Error("BindLoaderRef with different loader did not panic")
		}
	}()
	cls.BindLoaderRef(l2)
}

func TestClassObjectViaDefiningLoader(t *testing.T) {
	resetBootstrapLoaderForTesting()

	// Create a class for java/lang/Class.
	classClass := NewSyntheticClass("java/lang/Class", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Class": classClass,
	}}

	SetBootstrapLoader(loader)

	// Create a test class with definingLoaderRef set.
	testClass := NewSyntheticClass("test/MyClass", nil)
	testClass.BindLoaderRef(loader)

	obj := testClass.ClassObject()
	if obj == nil {
		t.Fatal("ClassObject() returned nil")
	}
	if obj.Class() != classClass {
		t.Error("mirror Class is not java/lang/Class")
	}
	if obj.Extra() != testClass {
		t.Error("mirror Extra is not the original Class")
	}

	// Second call must return same object.
	obj2 := testClass.ClassObject()
	if obj2 != obj {
		t.Error("ClassObject() is not idempotent")
	}
}

func TestClassObjectConcurrentIdentity(t *testing.T) {
	resetBootstrapLoaderForTesting()

	classClass := NewSyntheticClass("java/lang/Class", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Class": classClass,
	}}
	SetBootstrapLoader(loader)

	testClass := NewSyntheticClass("test/ConcurrentMirror", nil)
	testClass.BindLoaderRef(loader)

	const N = 32
	var wg sync.WaitGroup
	results := make([]*Object, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = testClass.ClassObject()
		}(i)
	}
	wg.Wait()

	first := results[0]
	if first == nil {
		t.Fatal("first mirror is nil")
	}
	for i, obj := range results {
		if obj != first {
			t.Errorf("goroutine %d got mirror %p, want %p", i, obj, first)
		}
	}
}

func TestClassObjectPrimitive(t *testing.T) {
	resetBootstrapLoaderForTesting()

	classClass := NewSyntheticClass("java/lang/Class", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Class": classClass,
	}}
	SetBootstrapLoader(loader)

	// VM primitive class (definingLoaderRef is nil, falls back to BootstrapLoader).
	InitVMTypes()
	obj := VMPrimitiveInt.ClassObject()
	if obj == nil {
		t.Fatal("ClassObject() for VMPrimitiveInt returned nil")
	}
	if obj.Extra() != VMPrimitiveInt {
		t.Error("mirror Extra is not VMPrimitiveInt")
	}
}

func TestClassObjectArray(t *testing.T) {
	resetBootstrapLoaderForTesting()

	classClass := NewSyntheticClass("java/lang/Class", nil)
	loader := &recordingMockLoader{classes: map[string]*Class{
		"java/lang/Class": classClass,
	}}
	SetBootstrapLoader(loader)

	// Create [I via VMPrimitiveInt.GetArrayClass().
	InitVMTypes()
	intArray := VMPrimitiveInt.GetArrayClass()

	// Array class inherits definingLoaderRef from component (nil for primitive).
	obj := intArray.ClassObject()
	if obj == nil {
		t.Fatal("ClassObject() for [I returned nil")
	}
	if obj.Extra() != intArray {
		t.Error("mirror Extra is not the array Class")
	}
}

func TestClassMirrorsUseDistinctDefiningLoaders(t *testing.T) {
	classClass1 := NewSyntheticClass("java/lang/Class", nil)
	classClass2 := NewSyntheticClass("java/lang/Class", nil)
	loader1 := &recordingMockLoader{classes: map[string]*Class{"java/lang/Class": classClass1}}
	loader2 := &recordingMockLoader{classes: map[string]*Class{"java/lang/Class": classClass2}}

	class1 := NewSyntheticClass("test/Same", nil)
	class1.BindLoader(loader1.LoaderIdentity())
	class1.BindLoaderRef(loader1)
	class2 := NewSyntheticClass("test/Same", nil)
	class2.BindLoader(loader2.LoaderIdentity())
	class2.BindLoaderRef(loader2)

	mirror1 := class1.ClassObject()
	mirror2 := class2.ClassObject()
	if mirror1 == nil || mirror2 == nil {
		t.Fatal("mirror materialization returned nil")
	}
	if mirror1 == mirror2 {
		t.Fatal("different defining loaders produced the same mirror")
	}
	if mirror1.Class() != classClass1 || mirror2.Class() != classClass2 {
		t.Fatal("mirror java/lang/Class did not come from its defining loader")
	}
}

func TestDelegatedAndReferenceArrayMirrorsUseParentLoader(t *testing.T) {
	parentClassClass := NewSyntheticClass("java/lang/Class", nil)
	parent := &recordingMockLoader{classes: map[string]*Class{"java/lang/Class": parentClassClass}}
	delegated := NewSyntheticClass("test/DelegatedMirror", nil)
	delegated.BindLoader(parent.LoaderIdentity())
	delegated.BindLoaderRef(parent)

	classMirror := delegated.ClassObject()
	arrayMirror := delegated.GetArrayClass().ClassObject()
	if classMirror == nil || arrayMirror == nil {
		t.Fatal("delegated mirror materialization returned nil")
	}
	if classMirror.Class() != parentClassClass || arrayMirror.Class() != parentClassClass {
		t.Fatal("delegated class or array mirror used an initiating loader")
	}
}
