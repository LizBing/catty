package kernel

import (
	"archive/zip"
	"io"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ClassPathLoader resolves class names against a list of directories
// containing `<name>.class` files and/or `.jar` archives (names in
// internal form, e.g. `com/example/Foo`). Classpath entries ending in
// `.jar` are indexed once at loader construction; lookup order follows
// classpath order (directories and jars interleaved as given).
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
	jars  []*jarIndex // one per .jar classpath entry, in order
	mu    sync.Mutex  // serializes concurrent duplicate loads
	cycle map[string]bool
}

// jarIndex is the name→entry index of one .jar classpath element.
// The archive's file handle stays open for the loader's lifetime —
// zip entries stream from it on every class load.
type jarIndex struct {
	path   string
	file   *os.File
	reader *zip.Reader
	byName map[string]*zip.File
}

// NewClassPathLoader attaches a directory/jar-based loader to the kernel.
func NewClassPathLoader(k *Kernel, entries []string) *ClassPathLoader {
	l := &ClassPathLoader{k: k, dirs: nil, cycle: make(map[string]bool)}
	for _, e := range entries {
		if strings.HasSuffix(strings.ToLower(e), ".jar") {
			if ji := openJarIndex(e); ji != nil {
				l.jars = append(l.jars, ji)
			}
			continue // missing/unreadable jar: skip like a missing dir
		}
		l.dirs = append(l.dirs, e)
	}
	return l
}

func openJarIndex(path string) *jarIndex {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	st, err := f.Stat()
	if err != nil {
		return nil
	}
	r, err := zip.NewReader(f, st.Size())
	if err != nil {
		return nil
	}
	ji := &jarIndex{path: path, file: f, reader: r, byName: make(map[string]*zip.File)}
	for _, file := range r.File {
		if strings.HasSuffix(file.Name, ".class") {
			ji.byName[file.Name] = file
		}
	}
	return ji
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

// dependency resolves superclasses/interfaces lazily; kept for embedders
// calling it directly.
func (l *ClassPathLoader) dependency(name string) (*Class, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.loadLocked(name)
}

func (l *ClassPathLoader) readClassFile(name string) ([]byte, error) {
	rel := filepath.FromSlash(name) + ".class"
	for _, dir := range l.dirs {
		if data, err := os.ReadFile(filepath.Join(dir, rel)); err == nil {
			return data, nil
		}
	}
	for _, ji := range l.jars {
		if zf, ok := ji.byName[name+".class"]; ok {
			rc, err := zf.Open()
			if err != nil {
				continue
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err == nil {
				return data, nil
			}
		}
	}
	return nil, fmt.Errorf("java.lang.ClassNotFoundException: %s", dotted(name))
}

