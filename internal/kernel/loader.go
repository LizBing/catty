package kernel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ClassPathLoader resolves class names against a list of directories
// containing `<name>.class` files (name in internal form, e.g.
// `com/example/Foo`). JARs arrive later (DEBT-0008 follow-up).
//
// Loading is lazy along the super chain: loading a class pulls in its
// superclass and interfaces on demand, with a cycle guard. Verification
// runs at link time when enabled (Options.Verify).
type ClassPathLoader struct {
	k       *Kernel
	dirs    []string
	cycle   map[string]bool // in-progress loads, guarded by k.mu
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

	l.k.mu.Lock()
	if l.cycle[name] {
		l.k.mu.Unlock()
		return nil, fmt.Errorf("ClassCircularityError: %s", name)
	}
	l.cycle[name] = true
	l.k.mu.Unlock()
	defer func() {
		l.k.mu.Lock()
		delete(l.cycle, name)
		l.k.mu.Unlock()
	}()

	data, err := l.readClassFile(name)
	if err != nil {
		return nil, err
	}
	c, err := l.k.LoadClassBytesWith(data, l.dependency)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// dependency resolves superclasses/interfaces lazily through this loader.
func (l *ClassPathLoader) dependency(name string) (*Class, error) {
	return l.Load(name)
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
