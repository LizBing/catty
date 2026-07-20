package classloader

import (
	"sync"
	"testing"

	"catty/rtda"
)

// countProvider is a ClassProvider that counts invocations.
type countProvider struct {
	class *rtda.Class
	mu    sync.Mutex
	count int
}

func (p *countProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name != "test/Dup" {
		return ProviderMiss()
	}
	p.mu.Lock()
	p.count++
	p.mu.Unlock()
	return ProviderClass(p.class)
}

// TestConcurrentLoadSingleIdentity verifies that loading the same class from N
// concurrent goroutines returns the SAME *Class pointer (exactly one cache entry,
// no duplicates). The K2 definition protocol ensures only one goroutine executes
// the provider chain; the rest wait on the defRecord and receive the published
// result.
func TestConcurrentLoadSingleIdentity(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/Dup", nil)
	cp := &countProvider{class: cls}
	cl := NewCustom(cp)

	const N = 32
	var wg sync.WaitGroup
	results := make([]*rtda.Class, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cl.LoadClass("test/Dup")
		}(i)
	}

	wg.Wait()

	// Every goroutine must see the same pointer.
	for i, c := range results {
		if c != cls {
			t.Errorf("goroutine %d got %p, want %p", i, c, cls)
		}
	}

	// Provider must have been called exactly once.
	if cp.count != 1 {
		t.Errorf("provider called %d times, want 1", cp.count)
	}

	// Cache must contain exactly one entry for the name.
	cl.mu.RLock()
	cached := cl.initiatingCache["test/Dup"]
	cl.mu.RUnlock()
	if cached != cls {
		t.Errorf("cache entry = %p, want %p", cached, cls)
	}
}

// TestConcurrentClassMirrorIdentity verifies that ClassObject returns the
// same java.lang.Class mirror (same *Object) when called from concurrent
// goroutines (ADR-0029 canonical mirror requirement).
func TestConcurrentClassMirrorIdentity(t *testing.T) {
	cls := &rtda.Class{}

	const N = 32
	var wg sync.WaitGroup
	results := make([]*rtda.Object, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cls.ClassObject(func() *rtda.Object {
				return &rtda.Object{}
			})
		}(i)
	}

	wg.Wait()

	// Every goroutine must see the same Object pointer.
	first := results[0]
	if first == nil {
		t.Fatal("first mirror is nil")
	}
	for i, obj := range results {
		if obj != first {
			t.Errorf("goroutine %d got mirror %p, want %p", i, obj, first)
		}
	}

	// A second call from the same goroutine should return the same.
	if second := cls.ClassObject(func() *rtda.Object {
		t.Fatal("factory should not be called again")
		return nil
	}); second != first {
		t.Errorf("second call from same goroutine: got %p, want %p", second, first)
	}
}

// TestLoadClassCachesSingleEntry verifies the basic single-goroutine case:
// loading the same class twice returns the same pointer.
func TestLoadClassCachesSingleEntry(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/Cached", nil)
	cl := NewCustom(&singleProvider{
		name:  "test/Cached",
		class: cls,
	})

	c1 := cl.LoadClass("test/Cached")
	c2 := cl.LoadClass("test/Cached")

	if c1 != cls {
		t.Errorf("first load: got %p, want %p", c1, cls)
	}
	if c2 != cls {
		t.Errorf("second load: got %p, want %p", c2, cls)
	}
	if c1 != c2 {
		t.Error("same-name loads returned different pointers")
	}
}

// singleProvider provides a single named class without any coordination.
type singleProvider struct {
	name  string
	class *rtda.Class
}

func (p *singleProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name == p.name {
		return ProviderClass(p.class)
	}
	return ProviderMiss()
}

// --- K2 typed failure propagation tests ---

// failureProvider returns a typed failure for a specific class name.
type failureProvider struct {
	className string
	kind      rtda.FailureKind
}

func (p *failureProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name == p.className {
		return ProviderFailure(&rtda.ClassLoadFailure{Kind: p.kind, Name: name})
	}
	return ProviderMiss()
}

func TestLoadClassResultTypedFailure(t *testing.T) {
	cl := NewCustom(&failureProvider{className: "test/Doomed", kind: rtda.FailureLinkage})

	result := cl.LoadClassResult("test/Doomed")
	if result.IsSuccess() {
		t.Fatal("expected failure for test/Doomed, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureLinkage {
		t.Errorf("failure kind = %v, want FailureLinkage", f.Kind)
	}
	if f.Name != "test/Doomed" {
		t.Errorf("failure name = %q, want test/Doomed", f.Name)
	}

	// Second call must return the same failure (cached on defRecord).
	result2 := cl.LoadClassResult("test/Doomed")
	if result2.IsSuccess() {
		t.Fatal("second call: expected failure, got success")
	}
	f2 := result2.Failure()
	if f2.Kind != rtda.FailureLinkage {
		t.Errorf("second call: failure kind = %v, want FailureLinkage", f2.Kind)
	}
}

// nameMismatchProvider returns a Class whose Name() differs from the requested name.
type nameMismatchProvider struct {
	class *rtda.Class
}

func (p *nameMismatchProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	// p.class.Name() != name — the loader must detect this.
	return ProviderClass(p.class)
}

func TestLoadClassNameMismatch(t *testing.T) {
	mismatched := rtda.NewSyntheticClass("wrong/Name", nil)
	cl := NewCustom(&nameMismatchProvider{class: mismatched})

	result := cl.LoadClassResult("requested/Name")
	if result.IsSuccess() {
		t.Fatal("expected failure for name mismatch, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureFormat {
		t.Errorf("failure kind = %v, want FailureFormat", f.Kind)
	}
}

// --- Circular dependency detection ---

// circularProvider simulates A depends on B, B depends on A.
type circularProvider struct {
	mu             sync.Mutex
	resolvingA     bool
	callCount      int
	classA, classB *rtda.Class
}

func (p *circularProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	p.mu.Lock()
	p.callCount++
	p.mu.Unlock()

	switch name {
	case "test/A":
		// Load B recursively, which triggers circularity back to A.
		result := loader.LoadClassResult("test/B")
		if result.IsSuccess() {
			return ProviderClass(p.classA)
		}
		// Circularity detected — propagate the failure.
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind:  rtda.FailureCircularity,
			Name:  "test/A",
			Cause: result.Failure(),
		})
	case "test/B":
		// Load A recursively — this triggers the circularity.
		result := loader.LoadClassResult("test/A")
		if result.IsSuccess() {
			return ProviderClass(p.classB)
		}
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind:  rtda.FailureCircularity,
			Name:  "test/B",
			Cause: result.Failure(),
		})
	}
	return ProviderMiss()
}

func TestLoadClassCircularDependency(t *testing.T) {
	classA := rtda.NewSyntheticClass("test/A", nil)
	classB := rtda.NewSyntheticClass("test/B", nil)
	cp := &circularProvider{classA: classA, classB: classB}
	cl := NewCustom(cp)

	// Loading A triggers loading B, which triggers loading A — circular.
	result := cl.LoadClassResult("test/A")
	if result.IsSuccess() {
		t.Fatal("expected circularity failure, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureCircularity {
		t.Errorf("failure kind = %v, want FailureCircularity", f.Kind)
	}
}

// --- Cross-session: defMu serialisation ---

// slowProvider embeds a delay in Provide so we can test that a concurrent
// LoadClassResult waits for the definition to complete.
type slowProvider struct {
	class   *rtda.Class
	started chan struct{}
	release chan struct{}
}

func (p *slowProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	close(p.started) // signal that we've started the provider
	<-p.release      // wait for the test to allow completion
	return ProviderClass(p.class)
}

func TestLoadClassConcurrentDefMuSerialisation(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/Slow", nil)
	sp := &slowProvider{
		class:   cls,
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	cl := NewCustom(sp)

	// Goroutine 1: start loading — defMu is acquired, provider blocks.
	var wg sync.WaitGroup
	wg.Add(1)
	var result1 rtda.ClassLoadResult
	go func() {
		defer wg.Done()
		result1 = cl.LoadClassResult("test/Slow")
	}()

	// Wait until the provider has started.
	<-sp.started

	// Goroutine 2: try to load the same class — must wait on defMu.
	// Since defMu serialises all definitions, this goroutine blocks until
	// goroutine 1 releases.
	wg.Add(1)
	var result2 rtda.ClassLoadResult
	go func() {
		defer wg.Done()
		result2 = cl.LoadClassResult("test/Slow")
	}()

	// Give goroutine 2 a moment to reach defMu.
	// Spin briefly (no proper synchronisation for this — the test just checks
	// that both eventually see the same published result).
	close(sp.release)

	wg.Wait()

	// Both must see the same Class.
	if !result1.IsSuccess() || result1.Class() != cls {
		t.Errorf("goroutine 1: success=%v, class=%p, want %p", result1.IsSuccess(), result1.Class(), cls)
	}
	if !result2.IsSuccess() || result2.Class() != cls {
		t.Errorf("goroutine 2: success=%v, class=%p, want %p", result2.IsSuccess(), result2.Class(), cls)
	}
}

// --- Provider chain miss ---

func TestLoadClassProviderAllMiss(t *testing.T) {
	// No providers handle the requested name.
	cl := NewCustom() // empty provider chain

	result := cl.LoadClassResult("nonexistent/Class")
	if result.IsSuccess() {
		t.Fatal("expected failure for unknown class, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureNotFound {
		t.Errorf("failure kind = %v, want FailureNotFound", f.Kind)
	}
}

// --- Delegated class (from parent-like provider) does not bind to this loader ---

func TestLoadClassDelegatedIdentity(t *testing.T) {
	// A "parent" provider returns a class with its own loader identity.
	parentClass := rtda.NewSyntheticClass("parent/Class", nil)
	parentID := rtda.NewLoaderIdentity()
	parentClass.BindLoader(parentID)

	parentProvider := &singleProvider{name: "parent/Class", class: parentClass}
	cl := NewCustom(parentProvider)

	// Load through cl — the class enters initiatingCache but retains its
	// parent loader identity.
	c := cl.LoadClass("parent/Class")
	if c != parentClass {
		t.Fatalf("load returned %p, want %p", c, parentClass)
	}
	if c.DefiningLoader() != parentID {
		t.Errorf("defining loader = %v, want %v (parent identity preserved)", c.DefiningLoader(), parentID)
	}

	// Cache must contain the class (initiating cache may hold delegated classes).
	cl.mu.RLock()
	cached := cl.initiatingCache["parent/Class"]
	cl.mu.RUnlock()
	if cached != parentClass {
		t.Error("delegated class not in initiating cache")
	}

	// defRecords must NOT contain the delegated class.
	cl.mu.RLock()
	dr := cl.defRecords["parent/Class"]
	cl.mu.RUnlock()
	if dr != nil {
		t.Error("delegated class appeared in defRecords (should only contain own definitions)")
	}
}
