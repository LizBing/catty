package kernel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ClassPathLoader resolves class names against a list of directories
// containing `<name>.class` files (name in internal form, e.g.
// `com/example/Foo`). JARs arrive later (DEBT-0008 follow-up).
//
// Loading is lazy along the super chain: loading a class pulls in its
// superclass and interfaces on demand, with a cycle guard. Verification
// runs at link time when enabled (Options.Verify).
//
// Concurrency: worker threads racing the first `new X` all land here.
// Loads are serialized by mu; waiters re-check the registry after acquiring
// and receive the winner's class. Genuine same-thread cycles (super chain
// loop) still raise ClassCircularityError.
type ClassPathLoader struct {
	k     *Kernel
	dirs  []string
	mu    sync.Mutex      // serializes concurrent duplicate loads
	cycle map[string]bool // same-goroutine recursion guard, guarded by mu
}

// NewClassPathLoader attaches a directory-based loader to the kernel.
func NewClassPathLoader(k *Kernel, dirs []string) *ClassPathLoader {
	return &ClassPathLoader{k: k, dirs: dirs, cycle: make(map[string]bool)}
}

// Load resolves and loads name (internal form). Returns the cached class
// when already present. Errors are LinkageError-style messages; the VM maps
// NoClassDefFoundError/ClassNotFoundException upstream as needed.
func (l *ClassPathLoader) Load(name string) (*Class, error) {
	if strings.HasPrefix(name, "[") {
		return l.k.ArrayClassOf(name[1:]), nil
	}
	if c, ok := l.k.ClassByName(name); ok {
		return c, nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.loadLocked(name)
}

func (l *ClassPathLoader) loadLocked(name string) (*Class, error) {
	// Winner may have registered while we waited on mu.
	if c, ok := l.k.ClassByName(name); ok {
		return c, nil
	}
	if l.cycle[name] {
		return nil, fmt.Errorf("ClassCircularityError: %s", name)
	}
	l.cycle[name] = true
	defer func() { delete(l.cycle, name) }()

	data, err := l.readClassFile(name)
	if err != nil {
		return nil, err
	}
	return l.k.LoadClassBytesWith(data, l.loadLocked) // recursion stays under mu
}

// dependency resolves superclasses/interfaces lazily; called with mu held.
func (l *ClassPathLoader) dependency(name string) (*Class, error) {
	return l.loadLocked(name)
}

func (l *ClassPathLoader) readClassFile(name string) ([]byte, error) {
	for _, dir := range l.dirs {
		path := filepath.Join(dir, filepath.FromSlash(name)+".class")
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("java.lang.ClassNotFoundException: %s", dotted(name))
}
