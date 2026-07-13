package classloader

import (
	"sync"
	"testing"

	"catty/rtda"
)

// slowProvider returns a class only after an explicit sync signal, so
// concurrent LoadClass calls are forced to race on the slow path.
type slowProvider struct {
	name   string
	class  *rtda.Class
	ready  chan struct{}
	give   chan struct{}
}

func (p *slowProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if name != p.name {
		return nil
	}
	// Signal that we are inside Provide (past the cache check).
	p.ready <- struct{}{}
	// Wait for the test to release us.
	<-p.give
	return p.class
}

// TestConcurrentLoadSingleIdentity verifies that loading the same class from N
// concurrent goroutines returns the SAME *Class pointer (exactly one cache entry,
// no duplicates), and that the provider is only called once.
func TestConcurrentLoadSingleIdentity(t *testing.T) {
	cls := &rtda.Class{}
	sp := &slowProvider{
		name:  "test/Dup",
		class: cls,
		ready: make(chan struct{}, 16),
		give:  make(chan struct{}),
	}
	cl := NewCustom(sp)

	const N = 16
	var wg sync.WaitGroup
	results := make([]*rtda.Class, N)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = cl.LoadClass("test/Dup")
		}(i)
	}

	// Wait for ALL goroutines to reach the slow path inside Provide.
	for i := 0; i < N; i++ {
		<-sp.ready
	}
	// Release the provider — one wins the cache CAS, the rest hit the
	// double-check inside LoadClass.
	close(sp.give)

	wg.Wait()

	// Every goroutine must see the same pointer.
	for i, c := range results {
		if c != cls {
			t.Errorf("goroutine %d got %p, want %p", i, c, cls)
		}
	}

	// Cache must contain exactly one entry for the name.
	cl.mu.RLock()
	cached := cl.cache["test/Dup"]
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
	cls := &rtda.Class{}
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

func (p *singleProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if name == p.name {
		return p.class
	}
	return nil
}
