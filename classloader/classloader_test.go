package classloader

import (
	"sync"
	"testing"
	"time"

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
	classClass := rtda.NewSyntheticClass("java/lang/Class", nil)
	cls := rtda.NewSyntheticClass("test/ConcurrentMirror", nil)
	cl := NewCustom()
	if result := cl.DefineClassResult("java/lang/Class", classClass); !result.IsSuccess() {
		t.Fatalf("define java/lang/Class: %v", result.FailureOrNil())
	}
	if result := cl.DefineClassResult("test/ConcurrentMirror", cls); !result.IsSuccess() {
		t.Fatalf("define test class: %v", result.FailureOrNil())
	}

	const N = 32
	var wg sync.WaitGroup
	results := make([]*rtda.Object, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cls.ClassObject()
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
	if second := cls.ClassObject(); second != first {
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

// --- Cross-session: per-defRecord serialisation ---

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

func TestLoadClassConcurrentDefSerialisation(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/Slow", nil)
	sp := &slowProvider{
		class:   cls,
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	cl := NewCustom(sp)

	// Goroutine 1: start loading — per-defRecord lock is acquired, provider blocks.
	var wg sync.WaitGroup
	wg.Add(1)
	var result1 rtda.ClassLoadResult
	go func() {
		defer wg.Done()
		result1 = cl.LoadClassResult("test/Slow")
	}()

	// Wait until the provider has started.
	<-sp.started

	// Goroutine 2: try to load the same class — must wait on defRecord.cond.
	// Since the defRecord state machine serialises same-name definitions, this
	// goroutine blocks until goroutine 1 publishes the terminal state.
	wg.Add(1)
	var result2 rtda.ClassLoadResult
	go func() {
		defer wg.Done()
		result2 = cl.LoadClassResult("test/Slow")
	}()

	// Release goroutine 1 to complete the definition.
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

// --- K2 Block 1: Typed definition (defineClassDirect) ---

// TestDefineClassDuplicateSameLoader verifies that a second defineClassDirect
// call for the same name on the same loader returns FailureDuplicateDefinition
// with the original class preserved.
func TestDefineClassDuplicateSameLoader(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/DupDirect", nil)
	cl := NewCustom()

	// First definition succeeds.
	result1 := cl.defineClassDirect("test/DupDirect", cls)
	if !result1.IsSuccess() {
		t.Fatalf("first defineClassDirect: expected success, got %v", result1.FailureOrNil())
	}
	if result1.Class() != cls {
		t.Errorf("first definition: got %p, want %p", result1.Class(), cls)
	}

	// Second definition must return FailureDuplicateDefinition.
	cls2 := rtda.NewSyntheticClass("test/DupDirect", nil)
	result2 := cl.defineClassDirect("test/DupDirect", cls2)
	if result2.IsSuccess() {
		t.Fatal("second defineClassDirect: expected failure, got success")
	}
	f := result2.Failure()
	if f.Kind != rtda.FailureDuplicateDefinition {
		t.Errorf("failure kind = %v, want FailureDuplicateDefinition", f.Kind)
	}
	if f.DefinedClass != cls {
		t.Errorf("DefinedClass = %p, want %p (original class)", f.DefinedClass, cls)
	}

	// Cache must still contain the original class.
	cl.mu.RLock()
	cached := cl.initiatingCache["test/DupDirect"]
	cl.mu.RUnlock()
	if cached != cls {
		t.Errorf("cache entry = %p, want %p (original preserved)", cached, cls)
	}
}

// TestTwoLoadersDefineSameName verifies that two loaders can independently
// define a class with the same name — different Class pointer and different
// DefiningLoader.
func TestTwoLoadersDefineSameName(t *testing.T) {
	cls1 := rtda.NewSyntheticClass("test/SameName", nil)
	cls2 := rtda.NewSyntheticClass("test/SameName", nil)

	cl1 := NewCustom()
	cl2 := NewCustom()

	r1 := cl1.defineClassDirect("test/SameName", cls1)
	if !r1.IsSuccess() {
		t.Fatalf("loader 1: expected success, got %v", r1.FailureOrNil())
	}

	r2 := cl2.defineClassDirect("test/SameName", cls2)
	if !r2.IsSuccess() {
		t.Fatalf("loader 2: expected success, got %v", r2.FailureOrNil())
	}

	c1 := r1.Class()
	c2 := r2.Class()
	if c1 == c2 {
		t.Error("different loaders: same Class pointer — must be distinct")
	}
	if c1.DefiningLoader() == c2.DefiningLoader() {
		t.Error("different loaders: same DefiningLoader identity — must be distinct")
	}
	if c1.Name() != "test/SameName" || c2.Name() != "test/SameName" {
		t.Errorf("name mismatch: %q / %q", c1.Name(), c2.Name())
	}
}

// TestDefineClassDirectConcurrentDefining verifies that when two goroutines
// race to defineClassDirect the same name, exactly one succeeds and the other
// receives FailureDuplicateDefinition with the winning Class preserved.
func TestDefineClassDirectConcurrentDefining(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/ConcurrentDefine", nil)
	cl := NewCustom()

	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	var r1, r2 rtda.ClassLoadResult

	go func() {
		defer wg.Done()
		<-start
		r1 = cl.defineClassDirect("test/ConcurrentDefine", cls)
	}()
	go func() {
		defer wg.Done()
		<-start
		r2 = cl.defineClassDirect("test/ConcurrentDefine", cls)
	}()

	// Release both goroutines simultaneously.
	close(start)

	// Wait with timeout.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: concurrent defineClassDirect may have deadlocked")
	}

	// One must succeed, one must return duplicate.
	successCount := 0
	dupCount := 0
	for _, r := range []rtda.ClassLoadResult{r1, r2} {
		if r.IsSuccess() {
			successCount++
			if r.Class() != cls {
				t.Errorf("successful definition returned %p, want %p", r.Class(), cls)
			}
		} else if r.Failure().Kind == rtda.FailureDuplicateDefinition {
			dupCount++
			if r.Failure().DefinedClass != cls {
				t.Errorf("duplicate DefinedClass = %p, want %p", r.Failure().DefinedClass, cls)
			}
		} else {
			t.Errorf("unexpected result: success=%v failure=%v", r.IsSuccess(), r.FailureOrNil())
		}
	}
	if successCount != 1 {
		t.Errorf("expected 1 success, got %d", successCount)
	}
	if dupCount != 1 {
		t.Errorf("expected 1 duplicate, got %d", dupCount)
	}
}

// --- K2 Item 1: Definition concurrency protocol ---

// TestSameSessionCircularity verifies that same-loader A→B→A circularity is
// detected through the session's seen map. The provider for A loads B, and the
// provider for B loads A — the session.seen map catches the re-entrant load.
func TestSameSessionCircularity(t *testing.T) {
	classA := rtda.NewSyntheticClass("test/A", nil)
	classB := rtda.NewSyntheticClass("test/B", nil)

	cp := &circularProvider{classA: classA, classB: classB}
	cl := NewCustom(cp)

	result := cl.LoadClassResult("test/A")
	if result.IsSuccess() {
		t.Fatal("expected circularity failure, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureCircularity {
		t.Errorf("failure kind = %v, want FailureCircularity", f.Kind)
	}

	// The second call (after circularity resolved) must still produce the
	// cached failure.
	result2 := cl.LoadClassResult("test/A")
	if result2.IsSuccess() {
		t.Fatal("second call: expected failure, got success")
	}
	if result2.Failure().Kind != rtda.FailureCircularity {
		t.Errorf("second call: failure kind = %v, want FailureCircularity", result2.Failure().Kind)
	}
}

// crossLoaderDelegatingProvider delegates resolution of a specific name to
// another ClassLoader, propagating the session for circularity detection.
type crossLoaderDelegatingProvider struct {
	delegateName string
	targetName   string
	target       *ClassLoader
}

func (p *crossLoaderDelegatingProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name != p.delegateName {
		return ProviderMiss()
	}
	targetName := p.targetName
	if targetName == "" {
		targetName = name
	}
	result := DelegateLoad(loader, p.target, targetName)
	if result.IsSuccess() {
		return ProviderClass(result.Class())
	}
	return ProviderFailure(result.Failure())
}

// TestCrossLoaderCircularity verifies that A→B→A circularity is detected
// across loader boundaries. Loader L1 resolves "test/XCA" whose provider
// delegates to L2 for "test/XCB". L2's provider calls back to L1 for
// "test/XCA" with session propagation — seenInChain detects the cycle.
func TestCrossLoaderCircularity(t *testing.T) {
	classB := rtda.NewSyntheticClass("test/XCB", nil)

	l1 := NewCustom()
	l2 := NewCustom()

	// L2: "test/XCB" depends on "test/XCA" — call back to L1.
	l2prov := &crossLoaderCallingProvider{
		triggerName:  "test/XCB",
		dependName:   "test/XCA",
		targetLoader: l1,
		class:        classB,
	}
	l2.providers = []ClassProvider{l2prov}

	// L1: "test/XCA" delegates to L2.
	l1prov := &crossLoaderDelegatingProvider{
		delegateName: "test/XCA",
		targetName:   "test/XCB",
		target:       l2,
	}
	l1.providers = []ClassProvider{l1prov}

	result := l1.LoadClassResult("test/XCA")
	if result.IsSuccess() {
		t.Fatal("expected cross-loader circularity failure, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureCircularity {
		t.Errorf("failure kind = %v, want FailureCircularity", f.Kind)
	}
}

// crossLoaderCallingProvider provides a class that triggers a dependency
// lookup on another loader via session propagation.
type crossLoaderCallingProvider struct {
	triggerName  string
	dependName   string
	targetLoader *ClassLoader
	class        *rtda.Class
}

func (p *crossLoaderCallingProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name != p.triggerName {
		return ProviderMiss()
	}
	// Resolve the dependency on the target loader, propagating the session.
	result := DelegateLoad(loader, p.targetLoader, p.dependName)
	if result.IsSuccess() {
		return ProviderClass(p.class)
	}
	return ProviderFailure(&rtda.ClassLoadFailure{
		Kind:  rtda.FailureCircularity,
		Name:  p.triggerName,
		Cause: result.Failure(),
	})
}

// TestConcurrentDefinitionDeadlockFree verifies that concurrent definitions of
// different names on different loaders proceed in parallel without deadlock.
// Each defRecord lock is independent — holding dr.mu for name A does not block
// definition of name B, even when both are triggered concurrently.
//
// This also exercises cross-loader session propagation with a hierarchical
// delegation pattern: L1's provider delegates "test/DeadA" to L2. L2's provider
// for "test/DeadA" calls back to L1 for "test/DeadDep" (the circularity path
// through a shared session chain). A second goroutine concurrently loads
// "test/DeadB" from L2 without cross-delegation, proving independent defRecord
// operations don't interfere.
func TestConcurrentDefinitionDeadlockFree(t *testing.T) {
	classA := rtda.NewSyntheticClass("test/DeadA", nil)
	classB := rtda.NewSyntheticClass("test/DeadB", nil)
	classDep := rtda.NewSyntheticClass("test/DeadDep", nil)

	l1 := NewCustom()
	l2 := NewCustom()

	// L1 delegates "test/DeadA" to L2.
	l1.providers = []ClassProvider{&crossLoaderDelegatingProvider{
		delegateName: "test/DeadA",
		target:       l2,
	}}

	// L2 provides "test/DeadA" — its provider calls back to L1 for
	// "test/DeadDep" via session propagation (exercises session chain).
	l2.providers = []ClassProvider{&crossLoaderCallingProvider{
		triggerName:  "test/DeadA",
		dependName:   "test/DeadDep",
		targetLoader: l1,
		class:        classA,
	}}

	// L1 has a provider for "test/DeadDep" (prevents not-found).
	l1.providers = append(l1.providers, &singleProvider{
		name:  "test/DeadDep",
		class: classDep,
	})

	// L2 directly provides "test/DeadB" (no cross-delegation).
	l2.providers = append(l2.providers, &singleProvider{
		name:  "test/DeadB",
		class: classB,
	})

	done := make(chan struct{}, 2)

	go func() {
		_ = l1.LoadClassResult("test/DeadA")
		done <- struct{}{}
	}()
	go func() {
		_ = l2.LoadClassResult("test/DeadB")
		done <- struct{}{}
	}()

	// Both must complete within the timeout. A deadlock would cause
	// the test to hang and timeout.
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout: concurrent definitions may have deadlocked")
		}
	}
}

type mutualDelegationBarrier struct {
	mu    sync.Mutex
	count int
	ready chan struct{}
}

func (b *mutualDelegationBarrier) arrive() {
	b.mu.Lock()
	b.count++
	if b.count == 2 {
		close(b.ready)
	}
	b.mu.Unlock()
	<-b.ready
}

type mutualDelegatingProvider struct {
	name       string
	dependName string
	target     *ClassLoader
	barrier    *mutualDelegationBarrier
}

func (p *mutualDelegatingProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name != p.name {
		return ProviderMiss()
	}
	p.barrier.arrive()
	result := DelegateLoad(loader, p.target, p.dependName)
	if result.IsSuccess() {
		return ProviderFailure(&rtda.ClassLoadFailure{Kind: rtda.FailureLinkage, Name: name})
	}
	return ProviderFailure(&rtda.ClassLoadFailure{
		Kind:  result.Failure().Kind,
		Name:  name,
		Cause: result.Failure(),
	})
}

// TestMutualCrossSessionDelegationDoesNotDeadlock pins the exact K2 lock-ring
// case: two independent top-level contexts own different definitions and then
// synchronously delegate to each other. The wait graph must break the cycle and
// publish terminal failures to both callers.
func TestMutualCrossSessionDelegationDoesNotDeadlock(t *testing.T) {
	l1 := NewCustom()
	l2 := NewCustom()
	barrier := &mutualDelegationBarrier{ready: make(chan struct{})}
	l1.providers = []ClassProvider{&mutualDelegatingProvider{
		name: "test/MutualA", dependName: "test/MutualB", target: l2, barrier: barrier,
	}}
	l2.providers = []ClassProvider{&mutualDelegatingProvider{
		name: "test/MutualB", dependName: "test/MutualA", target: l1, barrier: barrier,
	}}

	results := make(chan rtda.ClassLoadResult, 2)
	go func() { results <- l1.LoadClassResult("test/MutualA") }()
	go func() { results <- l2.LoadClassResult("test/MutualB") }()

	for i := 0; i < 2; i++ {
		select {
		case result := <-results:
			if result.IsSuccess() {
				t.Fatal("mutual delegation unexpectedly succeeded")
			}
			if result.Failure().Kind != rtda.FailureCircularity {
				t.Fatalf("failure kind = %v, want FailureCircularity", result.Failure().Kind)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timeout: mutual cross-session delegation deadlocked")
		}
	}
}

// TestMultiWaiterSameTerminalSuccess verifies that N concurrent goroutines
// waiting on the same defRecord all observe the identical Class pointer
// after the definition completes.
func TestMultiWaiterSameTerminalSuccess(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/MultiOk", nil)
	sp := &slowProvider{
		class:   cls,
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	cl := NewCustom(sp)

	const N = 16
	var wg sync.WaitGroup
	results := make([]rtda.ClassLoadResult, N)

	// Launch all goroutines — they will contend for the same defRecord.
	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cl.LoadClassResult("test/MultiOk")
		}(i)
	}

	// Wait for the provider to start, then release it.
	<-sp.started
	close(sp.release)

	wg.Wait()

	// All must see the same success.
	for i, r := range results {
		if !r.IsSuccess() {
			t.Errorf("goroutine %d: expected success, got failure: %v", i, r.FailureOrNil())
			continue
		}
		if r.Class() != cls {
			t.Errorf("goroutine %d: got %p, want %p", i, r.Class(), cls)
		}
	}
}

// TestMultiWaiterSameTerminalFailure verifies that N concurrent goroutines
// waiting on the same defRecord all observe the identical failure after
// the definition attempt fails.
func TestMultiWaiterSameTerminalFailure(t *testing.T) {
	sp := &slowFailureProvider{
		name:    "test/MultiFail",
		kind:    rtda.FailureLinkage,
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	cl := NewCustom(sp)

	const N = 16
	var wg sync.WaitGroup
	results := make([]rtda.ClassLoadResult, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cl.LoadClassResult("test/MultiFail")
		}(i)
	}

	// Wait for the provider to start, then release it.
	<-sp.started
	close(sp.release)

	wg.Wait()

	// All must see the same failure.
	for i, r := range results {
		if r.IsSuccess() {
			t.Errorf("goroutine %d: expected failure, got success", i)
			continue
		}
		f := r.Failure()
		if f.Kind != rtda.FailureLinkage {
			t.Errorf("goroutine %d: failure kind = %v, want FailureLinkage", i, f.Kind)
		}
		if f.Name != "test/MultiFail" {
			t.Errorf("goroutine %d: failure name = %q, want test/MultiFail", i, f.Name)
		}
	}
}

// slowFailureProvider is like slowProvider but returns a failure instead of
// a class. Used by TestMultiWaiterSameTerminalFailure.
type slowFailureProvider struct {
	name    string
	kind    rtda.FailureKind
	started chan struct{}
	release chan struct{}
}

func (p *slowFailureProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name != p.name {
		return ProviderMiss()
	}
	close(p.started)
	<-p.release
	return ProviderFailure(&rtda.ClassLoadFailure{Kind: p.kind, Name: name})
}

// --- K2 Item 2: Separate lookup/initiation from definition ---

// TestLookupMissThenDefinitionSuccess verifies that when LoadClassResult misses
// (defRecord transitions to defFailed due to provider miss), a subsequent
// defineClassDirect can still succeed — a failed lookup is NOT a duplicate.
func TestLookupMissThenDefinitionSuccess(t *testing.T) {
	cl := NewCustom() // empty provider chain — all lookups miss

	// Lookup misses — defRecord transitions to defFailed.
	result1 := cl.LoadClassResult("test/LookupThenDefine")
	if result1.IsSuccess() {
		t.Fatal("expected lookup miss, got success")
	}
	if result1.Failure().Kind != rtda.FailureNotFound {
		t.Fatalf("lookup failure kind = %v, want FailureNotFound", result1.Failure().Kind)
	}

	// defineClassDirect must succeed — failed lookup ≠ duplicate definition.
	cls := rtda.NewSyntheticClass("test/LookupThenDefine", nil)
	result2 := cl.defineClassDirect("test/LookupThenDefine", cls)
	if !result2.IsSuccess() {
		t.Fatalf("defineClassDirect after lookup miss: expected success, got %v", result2.FailureOrNil())
	}
	if result2.Class() != cls {
		t.Errorf("defineClassDirect returned %p, want %p", result2.Class(), cls)
	}
}

// TestLookupFailureThenCorrectedDefinition verifies that when defineClassDirect
// fails (e.g., name mismatch), a corrected second call succeeds rather than
// receiving FailureDuplicateDefinition.
func TestLookupFailureThenCorrectedDefinition(t *testing.T) {
	cl := NewCustom()

	// First attempt: name mismatch → failure.
	mismatched := rtda.NewSyntheticClass("wrong/Name", nil)
	result1 := cl.defineClassDirect("requested/Name", mismatched)
	if result1.IsSuccess() {
		t.Fatal("first defineClassDirect (mismatched): expected failure, got success")
	}
	if result1.Failure().Kind != rtda.FailureFormat {
		t.Fatalf("first failure kind = %v, want FailureFormat", result1.Failure().Kind)
	}

	// Second attempt: corrected name → must succeed.
	correct := rtda.NewSyntheticClass("requested/Name", nil)
	result2 := cl.defineClassDirect("requested/Name", correct)
	if !result2.IsSuccess() {
		t.Fatalf("second defineClassDirect (corrected): expected success, got %v", result2.FailureOrNil())
	}
	if result2.Class() != correct {
		t.Errorf("defineClassDirect returned %p, want %p", result2.Class(), correct)
	}

	// Verify the class is cached with the correct loader identity.
	cl.mu.RLock()
	cached := cl.initiatingCache["requested/Name"]
	cl.mu.RUnlock()
	if cached != correct {
		t.Errorf("cache entry = %p, want %p", cached, correct)
	}
}

// TestSuccessfulDefinitionThenDuplicate verifies that after a successful
// defineClassDirect, a second call for the same name returns
// FailureDuplicateDefinition with the original class preserved.
func TestSuccessfulDefinitionThenDuplicate(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/FirstWins", nil)
	cl := NewCustom()

	// First definition succeeds.
	result1 := cl.defineClassDirect("test/FirstWins", cls)
	if !result1.IsSuccess() {
		t.Fatalf("first defineClassDirect: expected success, got %v", result1.FailureOrNil())
	}

	// Second definition must return duplicate with original preserved.
	cls2 := rtda.NewSyntheticClass("test/FirstWins", nil)
	result2 := cl.defineClassDirect("test/FirstWins", cls2)
	if result2.IsSuccess() {
		t.Fatal("second defineClassDirect: expected failure, got success")
	}
	f := result2.Failure()
	if f.Kind != rtda.FailureDuplicateDefinition {
		t.Errorf("failure kind = %v, want FailureDuplicateDefinition", f.Kind)
	}
	if f.DefinedClass != cls {
		t.Errorf("DefinedClass = %p, want %p (original)", f.DefinedClass, cls)
	}

	// Cache must still contain the original class.
	cl.mu.RLock()
	cached := cl.initiatingCache["test/FirstWins"]
	cl.mu.RUnlock()
	if cached != cls {
		t.Errorf("cache entry = %p, want %p", cached, cls)
	}
}

// TestDelegatedLookupThenDefinitionRejected verifies that if initiatingCache
// already holds a class defined by a different loader, defineClassDirect on
// this loader returns FailureDuplicateDefinition with the delegated class
// preserved.
func TestDelegatedLookupThenDefinitionRejected(t *testing.T) {
	cl := NewCustom()

	// Simulate a delegated class: manually insert a class with a different
	// loader identity into the cache.
	delegated := rtda.NewSyntheticClass("test/Delegated", nil)
	otherID := rtda.NewLoaderIdentity()
	delegated.BindLoader(otherID)

	cl.mu.Lock()
	cl.initiatingCache["test/Delegated"] = delegated
	cl.mu.Unlock()

	// defineClassDirect must reject — delegated class already in cache.
	cls := rtda.NewSyntheticClass("test/Delegated", nil)
	result := cl.defineClassDirect("test/Delegated", cls)
	if result.IsSuccess() {
		t.Fatal("defineClassDirect for delegated name: expected failure, got success")
	}
	f := result.Failure()
	if f.Kind != rtda.FailureDuplicateDefinition {
		t.Errorf("failure kind = %v, want FailureDuplicateDefinition", f.Kind)
	}
	if f.DefinedClass != delegated {
		t.Errorf("DefinedClass = %p, want %p (delegated class)", f.DefinedClass, delegated)
	}

	// Cache must still contain the delegated class (not overwritten).
	cl.mu.RLock()
	cached := cl.initiatingCache["test/Delegated"]
	cl.mu.RUnlock()
	if cached != delegated {
		t.Errorf("cache entry = %p, want %p (delegated preserved)", cached, delegated)
	}
}

func TestDefineClassResultRejectsPreboundCandidate(t *testing.T) {
	owner := NewCustom()
	target := NewCustom()
	class := rtda.NewSyntheticClass("test/Prebound", nil)
	class.BindLoader(owner.LoaderIdentity())

	result := target.DefineClassResult("test/Prebound", class)
	if result.IsSuccess() {
		t.Fatal("pre-bound definition candidate unexpectedly succeeded")
	}
	if result.Failure().Kind != rtda.FailureLinkage {
		t.Fatalf("failure kind = %v, want FailureLinkage", result.Failure().Kind)
	}
	if got := target.LoadClassResult("test/Prebound"); got.IsSuccess() {
		t.Fatal("rejected pre-bound candidate polluted the initiating cache")
	}
}

// TestConcurrentDefinitionRace verifies that when N goroutines race to define
// the same name via defineClassDirect, exactly one wins and the rest receive
// FailureDuplicateDefinition.
func TestConcurrentDefinitionRace(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/RaceDefine", nil)
	cl := NewCustom()

	const N = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]rtda.ClassLoadResult, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			<-start
			results[idx] = cl.defineClassDirect("test/RaceDefine", cls)
		}(i)
	}

	// Release all goroutines simultaneously.
	close(start)

	// Wait with timeout.
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: concurrent definition race may have deadlocked")
	}

	successCount := 0
	dupCount := 0
	for _, r := range results {
		if r.IsSuccess() {
			successCount++
			if r.Class() != cls {
				t.Errorf("successful definition returned %p, want %p", r.Class(), cls)
			}
		} else if r.Failure().Kind == rtda.FailureDuplicateDefinition {
			dupCount++
			if r.Failure().DefinedClass != cls {
				t.Errorf("duplicate DefinedClass = %p, want %p", r.Failure().DefinedClass, cls)
			}
		} else {
			t.Errorf("unexpected result: success=%v failure=%v", r.IsSuccess(), r.FailureOrNil())
		}
	}
	if successCount != 1 {
		t.Errorf("expected 1 success, got %d", successCount)
	}
	if dupCount != N-1 {
		t.Errorf("expected %d duplicates, got %d", N-1, dupCount)
	}
}

// TestDefinitionLookupInterleaving verifies that defineClassDirect and
// LoadClassResult interleave correctly when the definition completes while
// a lookup is waiting on the defRecord.
func TestDefinitionLookupInterleaving(t *testing.T) {
	cls := rtda.NewSyntheticClass("test/Interleave", nil)
	var lookupResult rtda.ClassLoadResult
	lookupDone := make(chan struct{})
	sp := &slowProvider{
		class:   cls,
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	cl := NewCustom(sp)

	// Goroutine 1: LoadClassResult — starts definition, blocks in provider.
	go func() {
		lookupResult = cl.LoadClassResult("test/Interleave")
		close(lookupDone)
	}()

	// Wait for the provider to start.
	<-sp.started

	// Goroutine 2: defineClassDirect — sees defDefining, waits for terminal state.
	var defineResult rtda.ClassLoadResult
	defineDone := make(chan struct{})
	cls2 := rtda.NewSyntheticClass("test/Interleave", nil)
	go func() {
		defineResult = cl.defineClassDirect("test/Interleave", cls2)
		close(defineDone)
	}()

	// Release the slow provider — lookup completes, publishes defDefined.
	close(sp.release)

	// Both must complete.
	select {
	case <-lookupDone:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: LoadClassResult did not complete")
	}
	select {
	case <-defineDone:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: defineClassDirect did not complete")
	}

	// LoadClassResult must succeed.
	if !lookupResult.IsSuccess() {
		t.Fatalf("LoadClassResult: expected success, got %v", lookupResult.FailureOrNil())
	}
	if lookupResult.Class() != cls {
		t.Errorf("LoadClassResult returned %p, want %p", lookupResult.Class(), cls)
	}

	// defineClassDirect must return duplicate with the winning class.
	if defineResult.IsSuccess() {
		t.Fatal("defineClassDirect: expected failure (duplicate), got success")
	}
	if defineResult.Failure().Kind != rtda.FailureDuplicateDefinition {
		t.Errorf("defineClassDirect failure kind = %v, want FailureDuplicateDefinition", defineResult.Failure().Kind)
	}
	if defineResult.Failure().DefinedClass != cls {
		t.Errorf("defineClassDirect DefinedClass = %p, want %p", defineResult.Failure().DefinedClass, cls)
	}
}
