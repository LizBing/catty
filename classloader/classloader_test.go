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
	cls := &rtda.Class{}
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

func (p *singleProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if name == p.name {
		return ProviderClass(p.class)
	}
	return ProviderMiss()
}
