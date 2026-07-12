package classloader

import (
	"strings"

	"catty/classfile"
	"catty/classpath"
	"catty/native"
	"catty/rtda"
)

// --- Provider chain ---

// ClassProvider is one source of class definitions. Each provider is asked in
// order; the first one that returns a non-nil result wins. Providers receive
// the Loader so they can recursively resolve superclasses / component types.
type ClassProvider interface {
	Provide(name string, loader rtda.Loader) *rtda.Class
}

// ArrayProvider handles array type names ("[I", "[Ljava/lang/Object;", …).
type ArrayProvider struct{}

func (ArrayProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if !strings.HasPrefix(name, "[") {
		return nil
	}
	return rtda.NewArrayClass(name, loader)
}

// BootstrapProvider serves the irreducible synthetic bootstrap classes
// (Object, String, Class, System, Thread, Throwable). These are ALWAYS served
// from the native registry, never from a class file — they form the
// Go↔Java bridge and carry catty-specific native payloads (extra fields).
type BootstrapProvider struct{}

func (BootstrapProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if !native.IsBootstrap(name) {
		return nil
	}
	classes := native.SyntheticClasses()
	if fn := classes[name]; fn != nil {
		return fn(loader)
	}
	return nil
}

// SyntheticProvider serves all synthetic classes EXCEPT the bootstrap set.
// These are classes like StringBuilder, PrintStream, exception subclasses, and
// Comparable — they work as synthetic implementations when no real JDK is on
// the classpath, but can be replaced by real bytecode when java.base is available.
type SyntheticProvider struct{}

func (SyntheticProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if native.IsBootstrap(name) {
		return nil // BootstrapProvider handles these
	}
	classes := native.SyntheticClasses()
	if fn := classes[name]; fn != nil {
		return fn(loader)
	}
	return nil
}

// ClasspathProvider reads .class files from a classpath.
type ClasspathProvider struct {
	CP *classpath.Classpath
}

func (p ClasspathProvider) Provide(name string, loader rtda.Loader) *rtda.Class {
	if p.CP == nil {
		return nil
	}
	data, _, err := p.CP.ReadClass(strings.ReplaceAll(name, ".", "/"))
	if err != nil {
		return nil
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		return nil
	}
	return rtda.NewClass(cf, loader)
}

// --- ClassLoader ---

// ClassLoader loads, links, and caches classes through a chain of ClassProviders.
// It implements rtda.Loader so the interpreter can resolve classes at run time
// without an import cycle.
//
// Init (<clinit>) is NOT done here: it is triggered lazily by the interpreter
// at JVMS §5.5 events (new / getstatic / putstatic / invokestatic).
type ClassLoader struct {
	providers []ClassProvider
	cache     map[string]*rtda.Class
}

// New builds a ClassLoader with the standard provider chain:
//
//	ArrayProvider → BootstrapProvider → SyntheticProvider → ClasspathProvider
//
// Callers that need a different order can construct the loader manually.
func New(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{
		providers: []ClassProvider{
			ArrayProvider{},
			BootstrapProvider{},
			SyntheticProvider{},
			ClasspathProvider{CP: cp},
		},
		cache: make(map[string]*rtda.Class),
	}
}

// NewCustom builds a ClassLoader with explicit providers (for testing or
// custom class-resolution strategies).
func NewCustom(providers ...ClassProvider) *ClassLoader {
	return &ClassLoader{
		providers: providers,
		cache:     make(map[string]*rtda.Class),
	}
}

// LoadClass returns the loaded class for the given internal name, loading it on
// first access. Names use internal slashes ("java/lang/Object").
func (cl *ClassLoader) LoadClass(name string) *rtda.Class {
	if c := cl.cache[name]; c != nil {
		return c
	}
	for _, p := range cl.providers {
		if c := p.Provide(name, cl); c != nil {
			cl.cache[name] = c
			resolveNativeMethods(c)
			return c
		}
	}
	panic("catty: class not found: " + name)
}

// Classes returns every class the loader has cached so far (used by the AOT
// build to iterate emittable methods across all loaded classes).
func (cl *ClassLoader) Classes() []*rtda.Class {
	out := make([]*rtda.Class, 0, len(cl.cache))
	for _, c := range cl.cache {
		out = append(out, c)
	}
	return out
}

// --- Native method resolution ---

// resolveNativeMethods checks each method on a freshly loaded class: if it's
// native and a Go implementation is registered in the native registry, attach
// it. Methods without a registered implementation keep the default stub.
func resolveNativeMethods(class *rtda.Class) {
	for _, m := range class.Methods() {
		if !m.IsNative() {
			continue
		}
		if fn := native.GetNative(class.Name(), m.Name(), m.Descriptor()); fn != nil {
			m.SetNativeFunc(fn)
		}
	}
}
