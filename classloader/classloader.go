package classloader

import (
	"strings"
	"sync"

	"catty/classfile"
	"catty/classpath"
	"catty/native"
	"catty/rtda"
)

// --- Provider result ---

// ProviderResult is the typed return from a ClassProvider.
// Exactly one of Miss, Class, or Failure is meaningful.
type ProviderResult struct {
	Miss    bool
	Class   *rtda.Class
	Failure *rtda.ClassLoadFailure
}

// ProviderMiss means this provider does not handle the requested name.
func ProviderMiss() ProviderResult { return ProviderResult{Miss: true} }

// ProviderClass wraps a fully linked Class.
func ProviderClass(c *rtda.Class) ProviderResult { return ProviderResult{Class: c} }

// ProviderFailure wraps a terminal failure.
func ProviderFailure(f *rtda.ClassLoadFailure) ProviderResult { return ProviderResult{Failure: f} }

// --- Provider chain ---

// ClassProvider is one source of class definitions. Each provider is asked in
// order. Returns ProviderMiss to delegate to the next provider; ProviderClass
// for a fully linked class; or ProviderFailure for a terminal error that stops
// the chain. Providers receive the Loader so they can recursively resolve
// superclasses / component types.
type ClassProvider interface {
	Provide(name string, loader rtda.Loader) ProviderResult
}

// ArrayProvider handles array type names ("[I", "[Ljava/lang/Object;", …).
type ArrayProvider struct{}

func (ArrayProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if !strings.HasPrefix(name, "[") {
		return ProviderMiss()
	}
	c := rtda.NewArrayClass(name, loader)
	if c == nil {
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind: rtda.FailureLinkage,
			Name: name,
		})
	}
	return ProviderClass(c)
}

// BootstrapProvider serves the irreducible synthetic bootstrap classes
// (Object, String, Class, System, Thread, Throwable). These are ALWAYS served
// from the native registry, never from a class file — they form the
// Go↔Java bridge and carry catty-specific native payloads (extra fields).
type BootstrapProvider struct{}

func (BootstrapProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if !native.IsBootstrap(name) {
		return ProviderMiss()
	}
	classes := native.SyntheticClasses()
	if fn := classes[name]; fn != nil {
		return ProviderClass(fn(loader))
	}
	return ProviderMiss()
}

// SyntheticProvider serves all synthetic classes EXCEPT the bootstrap set.
// These are classes like StringBuilder, PrintStream, exception subclasses, and
// Comparable — they work as synthetic implementations when no real JDK is on
// the classpath, but can be replaced by real bytecode when java.base is available.
type SyntheticProvider struct{}

func (SyntheticProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if native.IsBootstrap(name) {
		return ProviderMiss() // BootstrapProvider handles these
	}
	classes := native.SyntheticClasses()
	if fn := classes[name]; fn != nil {
		return ProviderClass(fn(loader))
	}
	return ProviderMiss()
}

// ClasspathProvider reads .class files from a classpath.
type ClasspathProvider struct {
	CP *classpath.Classpath
}

func (p ClasspathProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if p.CP == nil {
		return ProviderMiss()
	}
	data, _, err := p.CP.ReadClass(strings.ReplaceAll(name, ".", "/"))
	if err != nil {
		// Distinguish typed miss from real I/O errors.
		if classpath.IsNotFound(err) {
			return ProviderMiss()
		}
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind:  rtda.FailureNotFound,
			Name:  name,
			Cause: err,
		})
	}
	cf, err := classfile.Parse(data)
	if err != nil {
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind:  rtda.FailureFormat,
			Name:  name,
			Cause: err,
		})
	}
	c := rtda.NewClass(cf, loader)
	if c == nil {
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind: rtda.FailureLinkage,
			Name: name,
		})
	}
	return ProviderClass(c)
}

// --- Definition state machine ---

// defState tracks the lifecycle of one (loader, name) definition.
type defState int

const (
	defUnresolved defState = iota
	defDefining
	defDefined
	defFailed
)

// defRecord is the atomic definition state for a single class name.
// It publishes either one fully linked Class or one terminal failure
// to all concurrent waiters.
type defRecord struct {
	mu      sync.Mutex
	cond    *sync.Cond
	state   defState
	class   *rtda.Class
	failure *rtda.ClassLoadFailure
}

func newDefRecord() *defRecord {
	d := &defRecord{}
	d.cond = sync.NewCond(&d.mu)
	return d
}

// loadSession tracks in-flight definitions to detect same-chain circularity.
// It is passed through provider-facing Loader views so recursive superclass/
// interface/component lookups within one top-level LoadClassResult share the
// session without depending on goroutine identity.
type loadSession struct {
	loader *ClassLoader
	seen   map[string]bool // names currently being defined in this session
}

func newLoadSession(cl *ClassLoader) *loadSession {
	return &loadSession{
		loader: cl,
		seen:   make(map[string]bool),
	}
}

// sessionLoader adapts a loadSession into a rtda.Loader. Recursive calls
// within the same session use this view.
type sessionLoader struct {
	session *loadSession
}

func (sl *sessionLoader) LoadClass(name string) *rtda.Class {
	return sl.session.loader.loadClassInternal(name, sl.session)
}

func (sl *sessionLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	return sl.session.loader.loadClassResultInternal(name, sl.session)
}

func (sl *sessionLoader) LoaderIdentity() *rtda.LoaderIdentity {
	return sl.session.loader.identity
}

// --- ClassLoader ---

// ClassLoader loads, links, and caches classes through a chain of ClassProviders.
// It implements rtda.Loader so the interpreter can resolve classes at run time
// without an import cycle.
//
// Init (<clinit>) is NOT done here: it is triggered lazily by the interpreter
// at JVMS §5.5 events (new / getstatic / putstatic / invokestatic).
//
// K2 identity semantics:
//   - initiatingCache maps (name) → Class for fast lookup; may return a Class
//     defined by a different (delegated) loader without rebinding identity.
//   - defRecords maps (name) → *defRecord, the defining-loader-owned definition
//     state. Only classes defined by THIS loader appear here.
//   - A delegated class (defined by parent/synthetic) enters only the initiating
//     cache, not the definition records.
type ClassLoader struct {
	identity  *rtda.LoaderIdentity
	providers []ClassProvider

	mu              sync.RWMutex
	initiatingCache map[string]*rtda.Class // fast lookup, may contain delegated Classes
	defRecords      map[string]*defRecord  // definition state for THIS loader
	defMu           sync.Mutex             // serializes first-time definition
}

// New builds a ClassLoader with the standard provider chain:
//
//	ArrayProvider → BootstrapProvider → SyntheticProvider → ClasspathProvider
//
// Callers that need a different order can construct the loader manually.
func New(cp *classpath.Classpath) *ClassLoader {
	return &ClassLoader{
		identity: rtda.NewLoaderIdentity(),
		providers: []ClassProvider{
			ArrayProvider{},
			BootstrapProvider{},
			SyntheticProvider{},
			ClasspathProvider{CP: cp},
		},
		initiatingCache: make(map[string]*rtda.Class),
		defRecords:      make(map[string]*defRecord),
	}
}

// NewCustom builds a ClassLoader with explicit providers (for testing or
// custom class-resolution strategies).
func NewCustom(providers ...ClassProvider) *ClassLoader {
	return &ClassLoader{
		identity:        rtda.NewLoaderIdentity(),
		providers:       providers,
		initiatingCache: make(map[string]*rtda.Class),
		defRecords:      make(map[string]*defRecord),
	}
}

// LoaderIdentity returns the opaque identity of this class loader.
func (cl *ClassLoader) LoaderIdentity() *rtda.LoaderIdentity {
	return cl.identity
}

// LoadClass is the must-load convenience method. It returns a fully linked
// Class or panics. Only bootstrap invariants and legacy callers proven
// unreachable from supported classfiles may use this method.
func (cl *ClassLoader) LoadClass(name string) *rtda.Class {
	result := cl.LoadClassResult(name)
	if result.IsSuccess() {
		return result.Class()
	}
	panic("catty: " + result.Failure().Error())
}

// LoadClassResult is the typed lookup method. It returns either a fully linked
// Class or a terminal ClassLoadFailure. Java-reachable resolution paths use
// this method so failures propagate as Java throwables rather than Go panics.
func (cl *ClassLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	return cl.loadClassResultInternal(name, newLoadSession(cl))
}

// loadClassResultInternal drives the typed definition protocol.
//
// Protocol:
//  1. Fast path: check initiatingCache (RLock).
//  2. If miss, get or create a defRecord for this name.
//  3. If defRecord is already Defined or Failed, return its published result.
//  4. Otherwise, serialize under defMu and run the provider chain.
//  5. Providers that return Miss delegate to the next provider.
//  6. The first ProviderClass or ProviderFailure terminates the chain.
//  7. Published Class enters both initiatingCache and defRecord.
//  8. Published Failure enters only defRecord (no partial cache entry).
func (cl *ClassLoader) loadClassResultInternal(name string, session *loadSession) rtda.ClassLoadResult {
	// 0. Primitive and void types are canonical VM identities (ADR-0033).
	// They bypass the provider chain and definition state entirely.
	if rtda.IsVMPrimitive(name) {
		c := rtda.VMPrimitiveForName(name)
		if c != nil {
			return rtda.NewClassResult(c)
		}
	}

	// 1. Fast path: initiating cache hit (may be delegated).
	cl.mu.RLock()
	if c := cl.initiatingCache[name]; c != nil {
		cl.mu.RUnlock()
		return rtda.NewClassResult(c)
	}
	cl.mu.RUnlock()

	// 2. Get or create definition record.
	cl.mu.Lock()
	dr := cl.defRecords[name]
	if dr == nil {
		dr = newDefRecord()
		cl.defRecords[name] = dr
	}
	cl.mu.Unlock()

	// 3. Check if already resolved.
	dr.mu.Lock()
	switch dr.state {
	case defDefined:
		c := dr.class
		dr.mu.Unlock()
		return rtda.NewClassResult(c)
	case defFailed:
		f := dr.failure
		dr.mu.Unlock()
		return rtda.NewFailureResult(f)
	case defDefining:
		// Another goroutine is defining — wait for terminal publication.
		// Check for circularity within this session.
		if session.seen[name] {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		dr.cond.Wait()
		// After wakeup, re-read the published state.
		switch dr.state {
		case defDefined:
			c := dr.class
			dr.mu.Unlock()
			return rtda.NewClassResult(c)
		case defFailed:
			f := dr.failure
			dr.mu.Unlock()
			return rtda.NewFailureResult(f)
		default:
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureLinkage,
				Name: name,
			})
		}
	default:
		// defUnresolved — we need to define.
	}
	dr.mu.Unlock()

	// 4. Serialize first-time definition and run the provider chain.
	return cl.defineClass(name, dr, session)
}

// defineClass runs the provider chain to define a class. It must be called
// with the caller holding the right to define (i.e., dr is unresolved and
// no one else is defining this name).
func (cl *ClassLoader) defineClass(name string, dr *defRecord, session *loadSession) rtda.ClassLoadResult {
	// Serialize all first-time definitions under defMu. This avoids
	// cross-name lock cycles without relying on goroutine identity.
	cl.defMu.Lock()

	// Double-check: another goroutine may have defined this while we waited
	// for defMu.
	dr.mu.Lock()
	switch dr.state {
	case defDefined:
		dr.mu.Unlock()
		cl.defMu.Unlock()
		return rtda.NewClassResult(dr.class)
	case defFailed:
		dr.mu.Unlock()
		cl.defMu.Unlock()
		return rtda.NewFailureResult(dr.failure)
	case defDefining:
		dr.mu.Unlock()
		cl.defMu.Unlock()
		// Someone else is defining — wait.
		dr.mu.Lock()
		for dr.state == defDefining {
			dr.cond.Wait()
		}
		switch dr.state {
		case defDefined:
			c := dr.class
			dr.mu.Unlock()
			return rtda.NewClassResult(c)
		case defFailed:
			f := dr.failure
			dr.mu.Unlock()
			return rtda.NewFailureResult(f)
		default:
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureLinkage,
				Name: name,
			})
		}
	}

	// Check for circularity within this session.
	if session.seen[name] {
		dr.mu.Unlock()
		cl.defMu.Unlock()
		return rtda.NewFailureResult(&rtda.ClassLoadFailure{
			Kind: rtda.FailureCircularity,
			Name: name,
		})
	}

	// Mark as defining and record in session.
	dr.state = defDefining
	session.seen[name] = true
	dr.mu.Unlock()

	// Release defMu so recursive LoadClassResult calls within the same
	// top-level session can proceed for *different* names. Same-name
	// circularity is caught by the session.seen check above.
	cl.defMu.Unlock()

	// Build a session-aware loader for recursive calls.
	sl := &sessionLoader{session: session}

	// Run the provider chain.
	for _, p := range cl.providers {
		result := p.Provide(name, sl)
		if result.Miss {
			continue
		}
		if result.Class != nil {
			// Success: bind loader identity (if not already bound) and publish.
			// Array classes created via GetArrayClass already have their
			// defining loader set from the component type (VM identity for
			// primitive arrays, component's loader for reference arrays).
			if result.Class.DefiningLoader() == nil {
				result.Class.BindLoader(cl.identity)
			}
			resolveNativeMethods(result.Class)

			dr.mu.Lock()
			dr.state = defDefined
			dr.class = result.Class
			dr.cond.Broadcast()
			dr.mu.Unlock()

			// Enter into initiating cache.
			cl.mu.Lock()
			if existing := cl.initiatingCache[name]; existing != nil {
				// Another defined-elsewhere class was cached while we
				// were defining. Keep the first winner (ours is newer
				// but we respect the cache invariant).
				cl.mu.Unlock()
				dr.mu.Lock()
				dr.state = defDefined
				dr.class = existing
				dr.mu.Unlock()
				return rtda.NewClassResult(existing)
			}
			cl.initiatingCache[name] = result.Class
			cl.mu.Unlock()

			return rtda.NewClassResult(result.Class)
		}
		if result.Failure != nil {
			// Terminal failure: publish failure only (no partial cache).
			dr.mu.Lock()
			dr.state = defFailed
			dr.failure = result.Failure
			dr.cond.Broadcast()
			dr.mu.Unlock()
			return rtda.NewFailureResult(result.Failure)
		}
	}

	// All providers missed.
	failure := &rtda.ClassLoadFailure{
		Kind: rtda.FailureNotFound,
		Name: name,
	}
	dr.mu.Lock()
	dr.state = defFailed
	dr.failure = failure
	dr.cond.Broadcast()
	dr.mu.Unlock()
	return rtda.NewFailureResult(failure)
}

// loadClassInternal is the session-aware must-load helper used by sessionLoader.
func (cl *ClassLoader) loadClassInternal(name string, session *loadSession) *rtda.Class {
	result := cl.loadClassResultInternal(name, session)
	if result.IsSuccess() {
		return result.Class()
	}
	return nil
}

// Classes returns every class the loader has cached so far (used by the AOT
// build to iterate emittable methods across all loaded classes). Thread-safe.
func (cl *ClassLoader) Classes() []*rtda.Class {
	cl.mu.RLock()
	defer cl.mu.RUnlock()
	out := make([]*rtda.Class, 0, len(cl.initiatingCache))
	for _, c := range cl.initiatingCache {
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
