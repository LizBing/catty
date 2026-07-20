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
type SyntheticProvider struct{}

func (SyntheticProvider) Provide(name string, loader rtda.Loader) ProviderResult {
	if native.IsBootstrap(name) {
		return ProviderMiss()
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

	// Validate: classfile's declared name must match the requested binary name.
	if cf.ClassName() != name {
		return ProviderFailure(&rtda.ClassLoadFailure{
			Kind: rtda.FailureFormat,
			Name: name,
			Cause: &classNameMismatchError{
				requested: name,
				declared:  cf.ClassName(),
			},
		})
	}

	result := rtda.BuildClass(cf, loader)
	if result.Class != nil {
		return ProviderClass(result.Class)
	}
	// Propagate the typed dependency failure.
	return ProviderFailure(result.Failure)
}

// classNameMismatchError records a mismatch between the requested binary name
// and the classfile's declared this_class.
type classNameMismatchError struct {
	requested string
	declared  string
}

func (e *classNameMismatchError) Error() string {
	return "catty: classfile name mismatch: requested " + e.requested + ", declared " + e.declared
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
//
// defRecord is owned by the defining loader. Delegated classes (from
// a parent or synthetic provider) do NOT create defRecords — they
// enter only the initiating cache.
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
// within the same session use this view. It bypasses the top-level defMu
// acquisition since defMu is already held by the top-level caller.
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
// It implements rtda.Loader so the interpreter can resolve classes at run time.
//
// K2 identity semantics:
//   - initiatingCache maps (name) → Class for fast lookup; may return a Class
//     defined by a different (delegated) loader without rebinding identity.
//   - defRecords maps (name) → *defRecord for classes DEFINED by THIS loader.
//     Delegated classes enter only the initiating cache, never defRecords.
//   - defMu serialises ALL top-level definition. It is acquired once per
//     LoadClassResult call and held across the entire provider chain including
//     recursive dependency resolution. This prevents cross-session deadlock
//     without depending on goroutine identity.
type ClassLoader struct {
	identity  *rtda.LoaderIdentity
	providers []ClassProvider

	mu              sync.RWMutex
	initiatingCache map[string]*rtda.Class // fast lookup, may contain delegated Classes
	defRecords      map[string]*defRecord  // definition state for THIS loader only
	defMu           sync.Mutex             // serialises all top-level definition
}

// New builds a ClassLoader with the standard provider chain:
//
//	ArrayProvider → BootstrapProvider → SyntheticProvider → ClasspathProvider
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

// NewCustom builds a ClassLoader with explicit providers.
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

// LoadClass is the must-load convenience method. Panics on failure.
// Only valid for bootstrap invariants; Java-reachable paths use LoadClassResult.
func (cl *ClassLoader) LoadClass(name string) *rtda.Class {
	result := cl.LoadClassResult(name)
	if result.IsSuccess() {
		return result.Class()
	}
	panic("catty: " + result.Failure().Error())
}

// LoadClassResult is the typed lookup method — the top-level entry point.
//
// Protocol:
//  1. Acquire defMu (global definition serialisation).
//  2. Create a loadSession and dispatch to loadClassResultInternal.
//  3. Release defMu on return.
//
// defMu serialises ALL definitions so that no two sessions can be defining
// simultaneously. Recursive resolution within the provider chain bypasses
// defMu (the sessionLoader calls loadClassResultInternal directly).
func (cl *ClassLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	cl.defMu.Lock()
	defer cl.defMu.Unlock()
	return cl.loadClassResultInternal(name, newLoadSession(cl))
}

// loadClassResultInternal drives the typed definition protocol.
//
// defMu MUST be held by the caller (acquired at the top-level LoadClassResult
// or inherited from the enclosing top-level call via sessionLoader recursion).
//
// Protocol:
//  1. Primitive/void types bypass the provider chain (canonical VM identities).
//  2. Fast path: check initiatingCache (RLock). May return delegated Classes.
//  3. If miss, get or create a defRecord for this name within THIS loader.
//  4. If defRecord is Defined or Failed, return its published result.
//  5. If defRecord is Defining (another goroutine within the same top-level
//     session or a different goroutine that is the defMu holder):
//     a. Check session.seen for THIS session's circularity.
//     b. Wait on defRecord.cond for terminal publication.
//  6. Otherwise (defUnresolved): define via the provider chain.
func (cl *ClassLoader) loadClassResultInternal(name string, session *loadSession) rtda.ClassLoadResult {
	// 0. Primitive and void types are canonical VM identities.
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

	// 2. Get or create definition record for THIS loader.
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
		// Another call is defining this name. Check session circularity.
		if session.seen[name] {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		// Wait for terminal publication.
		dr.cond.Wait()
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
	dr.mu.Unlock()

	// 4. defUnresolved — define this class.
	return cl.defineClass(name, dr, session)
}

// defineClass runs the provider chain to define a class.
//
// defMu is already held by the top-level caller. defineClass is called
// from loadClassResultInternal under that lock, or recursively from a
// provider via sessionLoader → loadClassResultInternal → defineClass.
//
// Only ONE goroutine can be inside defineClass at any time (serialised
// by defMu), so cross-session deadlock is impossible.
func (cl *ClassLoader) defineClass(name string, dr *defRecord, session *loadSession) rtda.ClassLoadResult {
	// Double-check: another call (within the same defMu scope or a prior
	// concurrent call) may have resolved this while we were setting up.
	dr.mu.Lock()
	switch dr.state {
	case defDefined:
		dr.mu.Unlock()
		return rtda.NewClassResult(dr.class)
	case defFailed:
		dr.mu.Unlock()
		return rtda.NewFailureResult(dr.failure)
	case defDefining:
		// This should not happen under defMu serialisation unless
		// there's intra-session reentrancy (circular dependency).
		if session.seen[name] {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		// Wait for terminal publication (same-session reentrancy).
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
		return rtda.NewFailureResult(&rtda.ClassLoadFailure{
			Kind: rtda.FailureCircularity,
			Name: name,
		})
	}

	// Mark as defining and record in session.
	dr.state = defDefining
	session.seen[name] = true
	dr.mu.Unlock()

	// Build a session-aware loader for recursive calls.
	sl := &sessionLoader{session: session}

	// Run the provider chain.
	for _, p := range cl.providers {
		result := p.Provide(name, sl)
		if result.Miss {
			continue
		}
		if result.Class != nil {
			// Validate: provider must return a Class whose Name() matches
			// the requested binary name.
			if result.Class.Name() != name {
				failure := &rtda.ClassLoadFailure{
					Kind:  rtda.FailureFormat,
					Name:  name,
					Cause: &classNameMismatchError{requested: name, declared: result.Class.Name()},
				}
				dr.mu.Lock()
				dr.state = defFailed
				dr.failure = failure
				dr.cond.Broadcast()
				dr.mu.Unlock()
				return rtda.NewFailureResult(failure)
			}

			// Bind loader identity (if not already bound). Array classes
			// created via GetArrayClass already have their defining loader
			// set from the component type.
			isDelegated := result.Class.DefiningLoader() != nil
			if !isDelegated {
				result.Class.BindLoader(cl.identity)
			}

			resolveNativeMethods(result.Class)

			// K2: add to initiatingCache FIRST so the FastPath catches
			// subsequent calls within the same session.
			cl.mu.Lock()
			cl.initiatingCache[name] = result.Class
			cl.mu.Unlock()

			// Publish in defRecord to wake any waiters.
			dr.mu.Lock()
			dr.state = defDefined
			dr.class = result.Class
			dr.cond.Broadcast()
			dr.mu.Unlock()

			// K2: delegated classes (from parent/synthetic providers that
			// already have a defining loader) must NOT leave a defRecord
			// — defRecords tracks only classes THIS loader defines.
			if isDelegated {
				cl.mu.Lock()
				delete(cl.defRecords, name)
				cl.mu.Unlock()
			}

			return rtda.NewClassResult(result.Class)
		}
		if result.Failure != nil {
			// Terminal failure: publish only in defRecord (no partial cache).
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
// Returns nil on failure (typed failure is discarded — the caller within the
// provider chain will observe the nil and propagate its own typed failure).
func (cl *ClassLoader) loadClassInternal(name string, session *loadSession) *rtda.Class {
	result := cl.loadClassResultInternal(name, session)
	if result.IsSuccess() {
		return result.Class()
	}
	return nil
}

// Classes returns every class the loader has cached so far.
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
