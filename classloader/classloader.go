package classloader

import (
	"fmt"
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
	result := rtda.NewArrayClassResult(name, loader)
	if result.Class != nil {
		return ProviderClass(result.Class)
	}
	// Propagate typed failure with component cause preserved.
	return ProviderFailure(result.Failure)
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
	owner   *loadContext
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
//
// context is shared by every loader participating in one synchronous
// delegation chain. It records (loader,name), rather than name alone, because
// equal binary names in different defining loaders are distinct identities.
type loadSession struct {
	loader  *ClassLoader
	context *loadContext
}

func newLoadSession(cl *ClassLoader) *loadSession {
	return &loadSession{
		loader:  cl,
		context: newLoadContext(),
	}
}

type loadKey struct {
	loader *ClassLoader
	name   string
}

type loadContext struct {
	mu   sync.Mutex
	seen map[loadKey]bool
}

func newLoadContext() *loadContext {
	return &loadContext{seen: make(map[loadKey]bool)}
}

// seenInChain reports whether this defining identity is already active in the
// synchronous delegation chain.
func (s *loadSession) seenInChain(name string) bool {
	key := loadKey{loader: s.loader, name: name}
	s.context.mu.Lock()
	defer s.context.mu.Unlock()
	return s.context.seen[key]
}

func (s *loadSession) enter(name string) {
	key := loadKey{loader: s.loader, name: name}
	s.context.mu.Lock()
	s.context.seen[key] = true
	s.context.mu.Unlock()
}

func (s *loadSession) leave(name string) {
	key := loadKey{loader: s.loader, name: name}
	s.context.mu.Lock()
	delete(s.context.seen, key)
	s.context.mu.Unlock()
}

// definitionWaitGraph records cross-context waits. Per-record locks prevent
// duplicate publication; this graph prevents two independent top-level load
// contexts from waiting on each other forever (A owns X and waits for Y while
// B owns Y and waits for X).
var definitionWaitGraph = struct {
	sync.Mutex
	waiting map[*loadContext]*loadContext
}{waiting: make(map[*loadContext]*loadContext)}

// beginDefinitionWait registers waiter -> owner. It returns false if adding
// the edge would close a wait cycle.
func beginDefinitionWait(waiter, owner *loadContext) bool {
	if waiter == nil || owner == nil || waiter == owner {
		return false
	}
	definitionWaitGraph.Lock()
	defer definitionWaitGraph.Unlock()
	for current := owner; current != nil; current = definitionWaitGraph.waiting[current] {
		if current == waiter {
			return false
		}
	}
	definitionWaitGraph.waiting[waiter] = owner
	return true
}

func endDefinitionWait(waiter *loadContext) {
	definitionWaitGraph.Lock()
	delete(definitionWaitGraph.waiting, waiter)
	definitionWaitGraph.Unlock()
}

// sessionLoader adapts a loadSession into a rtda.Loader. Recursive calls
// within the same session use this view so circularity detection shares the
// same session.seen map.
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
//   - Per-defRecord locking (dr.mu) serialises same-name definition.
//     The defRecord state machine (defUnresolved → defDefining → terminal)
//     ensures exactly-once definition without a global loader lock.
//     Concurrent definitions of different names proceed in parallel.
type ClassLoader struct {
	identity  *rtda.LoaderIdentity
	providers []ClassProvider

	mu              sync.RWMutex
	initiatingCache map[string]*rtda.Class // fast lookup, may contain delegated Classes
	defRecords      map[string]*defRecord  // definition state for THIS loader only
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
//  1. Create a loadSession and dispatch to loadClassResultInternal.
//  2. loadClassResultInternal uses per-defRecord locking to serialise
//     same-name definitions. Concurrent calls for different names
//     proceed in parallel.
func (cl *ClassLoader) LoadClassResult(name string) rtda.ClassLoadResult {
	return cl.loadClassResultInternal(name, newLoadSession(cl))
}

// loadClassResultWithSession creates a child loader view sharing the parent's
// load context. Kept as the internal implementation of DelegateLoad.
func (cl *ClassLoader) loadClassResultWithSession(name string, parent *loadSession) rtda.ClassLoadResult {
	child := &loadSession{
		loader:  cl,
		context: parent.context,
	}
	return cl.loadClassResultInternal(name, child)
}

// DelegateLoad performs a provider-to-loader delegation while preserving the
// active load context. Providers must use this operation instead of calling a
// target loader's top-level LoadClassResult directly.
func DelegateLoad(providerLoader rtda.Loader, target *ClassLoader, name string) rtda.ClassLoadResult {
	if sl, ok := providerLoader.(*sessionLoader); ok {
		return target.loadClassResultWithSession(name, sl.session)
	}
	return target.LoadClassResult(name)
}

// loadClassResultInternal drives the typed definition protocol.
//
// Protocol:
//  1. Primitive/void types bypass the provider chain (canonical VM identities).
//  2. Fast path: check initiatingCache (RLock). May return delegated Classes.
//  3. If miss, get or create a defRecord for this name within THIS loader.
//  4. If defRecord is Defined or Failed, return its published result.
//  5. If defRecord is Defining (another goroutine is defining this name):
//     a. Check session.seenInChain for circularity.
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
		if session.seenInChain(name) {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		owner := dr.owner
		if !beginDefinitionWait(session.context, owner) {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		for dr.state == defDefining {
			dr.cond.Wait()
		}
		endDefinitionWait(session.context)
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
// Per-defRecord locking serialises same-name definition. Multiple goroutines
// may be in defineClass concurrently for different names, each holding a
// different defRecord.mu. The double-check under dr.mu prevents TOCTOU races
// between loadClassResultInternal's state check and defineClass entry.
func (cl *ClassLoader) defineClass(name string, dr *defRecord, session *loadSession) rtda.ClassLoadResult {
	// Double-check: another goroutine may have resolved this between
	// loadClassResultInternal's state check and defineClass entry.
	dr.mu.Lock()
	switch dr.state {
	case defDefined:
		dr.mu.Unlock()
		return rtda.NewClassResult(dr.class)
	case defFailed:
		dr.mu.Unlock()
		return rtda.NewFailureResult(dr.failure)
	case defDefining:
		// Another goroutine transitioned to defDefining between
		// loadClassResultInternal's state check and defineClass entry.
		// Check for intra-session circularity.
		if session.seenInChain(name) {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		owner := dr.owner
		if !beginDefinitionWait(session.context, owner) {
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind: rtda.FailureCircularity,
				Name: name,
			})
		}
		for dr.state == defDefining {
			dr.cond.Wait()
		}
		endDefinitionWait(session.context)
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
	if session.seenInChain(name) {
		dr.mu.Unlock()
		return rtda.NewFailureResult(&rtda.ClassLoadFailure{
			Kind: rtda.FailureCircularity,
			Name: name,
		})
	}

	// Mark as defining and record in session.
	dr.state = defDefining
	dr.owner = session.context
	session.enter(name)
	dr.mu.Unlock()
	defer session.leave(name)

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
				dr.owner = nil
				dr.cond.Broadcast()
				dr.mu.Unlock()
				return rtda.NewFailureResult(failure)
			}

			return cl.registerDefinition(result.Class, name, dr)
		}
		if result.Failure != nil {
			// Terminal failure: publish only in defRecord (no partial cache).
			dr.mu.Lock()
			dr.state = defFailed
			dr.failure = result.Failure
			dr.owner = nil
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
	dr.owner = nil
	dr.cond.Broadcast()
	dr.mu.Unlock()
	return rtda.NewFailureResult(failure)
}

// registerDefinition installs a fully linked Class as this loader's definition
// for the given name. It binds the loader identity, resolves native methods,
// publishes to the initiatingCache and defRecord, and cleans up defRecords for
// delegated classes.
//
// The caller must have already transitioned dr.state to defDefining. This
// function acquires dr.mu when publishing the terminal defDefined state.
func (cl *ClassLoader) registerDefinition(class *rtda.Class, name string, dr *defRecord) rtda.ClassLoadResult {
	// Bind loader identity (if not already bound). Array classes created via
	// GetArrayClass already have their defining loader set from the component.
	isDelegated := class.DefiningLoader() != nil
	if !isDelegated {
		class.BindLoader(cl.identity)
		class.BindLoaderRef(cl)
	}

	resolveNativeMethods(class)

	// K2: add to initiatingCache FIRST so the FastPath catches
	// subsequent calls within the same session.
	cl.mu.Lock()
	cl.initiatingCache[name] = class
	cl.mu.Unlock()

	// Publish in defRecord to wake any waiters.
	dr.mu.Lock()
	dr.state = defDefined
	dr.class = class
	dr.owner = nil
	dr.cond.Broadcast()
	dr.mu.Unlock()

	// K2: delegated classes (from parent/synthetic providers that already have
	// a defining loader) must NOT leave a defRecord — defRecords tracks only
	// classes THIS loader defines.
	if isDelegated {
		cl.mu.Lock()
		delete(cl.defRecords, name)
		cl.mu.Unlock()
	}

	return rtda.NewClassResult(class)
}

// DefineClassResult installs a pre-built, fully linked Class as this loader's
// definition. It is the production typed definition service, distinct from
// LoadClassResult lookup/initiation.
// for the given name. This is the typed definition operation separate from
// lookup/initiation (LoadClassResult).
//
// Protocol:
//   - Delegated class check: if initiatingCache already holds a Class defined
//     by a different loader, return FailureDuplicateDefinition with that Class
//     preserved. This prevents a child loader from overwriting a parent-defined
//     class.
//   - defDefined: this loader has successfully defined the name → duplicate.
//   - defDefining: a concurrent definition is in progress → wait for terminal
//     state; if success → duplicate, if failure → retry.
//   - defFailed: a prior lookup/initiation failed. This is NOT a duplicate —
//     the caller may have corrected the error. Reset to defUnresolved and retry.
//   - defUnresolved: proceed with definition.
func (cl *ClassLoader) DefineClassResult(name string, class *rtda.Class) rtda.ClassLoadResult {
	// 0. Delegated identity protection: if initiatingCache already holds a
	//    delegated Class (defined by a different loader), reject the definition.
	cl.mu.RLock()
	if cached := cl.initiatingCache[name]; cached != nil {
		if cached.DefiningLoader() != cl.identity {
			cl.mu.RUnlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind:         rtda.FailureDuplicateDefinition,
				Name:         name,
				DefinedClass: cached,
			})
		}
	}
	cl.mu.RUnlock()

	// 1. Get or create defRecord for this name.
	cl.mu.Lock()
	dr := cl.defRecords[name]
	if dr == nil {
		dr = newDefRecord()
		cl.defRecords[name] = dr
	}
	cl.mu.Unlock()

	// 2. Check defRecord state under lock.
	dr.mu.Lock()
	switch dr.state {
	case defDefined:
		// This loader has already successfully defined the name.
		c := dr.class
		dr.mu.Unlock()
		return rtda.NewFailureResult(&rtda.ClassLoadFailure{
			Kind:         rtda.FailureDuplicateDefinition,
			Name:         name,
			DefinedClass: c,
		})
	case defDefining:
		// A concurrent definition is in progress. Wait for terminal state.
		for dr.state == defDefining {
			dr.cond.Wait()
		}
		if dr.state == defDefined {
			c := dr.class
			dr.mu.Unlock()
			return rtda.NewFailureResult(&rtda.ClassLoadFailure{
				Kind:         rtda.FailureDuplicateDefinition,
				Name:         name,
				DefinedClass: c,
			})
		}
		// Concurrent definition failed — reset and retry.
		dr.state = defUnresolved
		dr.failure = nil
	case defFailed:
		// A prior lookup/initiation failed. This is NOT a duplicate —
		// the caller may have corrected the error (e.g., format failure
		// on first attempt, corrected class on second).
		dr.state = defUnresolved
		dr.failure = nil
	}
	// dr.state is now defUnresolved (either was already unresolved, or was
	// defFailed/defDefining→defFailed and has been reset).
	if class == nil {
		failure := &rtda.ClassLoadFailure{
			Kind:  rtda.FailureFormat,
			Name:  name,
			Cause: &classNameMismatchError{requested: name, declared: "<nil>"},
		}
		dr.state = defFailed
		dr.failure = failure
		dr.cond.Broadcast()
		dr.mu.Unlock()
		return rtda.NewFailureResult(failure)
	}
	// A definition candidate must be unbound. Pre-bound Classes are valid only
	// as delegated lookup results and must never be installed by this service.
	if class.DefiningLoader() != nil {
		failure := &rtda.ClassLoadFailure{
			Kind:  rtda.FailureLinkage,
			Name:  name,
			Cause: fmt.Errorf("catty: definition candidate %s is already bound to a loader", class.Name()),
		}
		dr.state = defFailed
		dr.failure = failure
		dr.cond.Broadcast()
		dr.mu.Unlock()
		return rtda.NewFailureResult(failure)
	}

	// 3. Validate name match.
	if class.Name() != name {
		failure := &rtda.ClassLoadFailure{
			Kind:  rtda.FailureFormat,
			Name:  name,
			Cause: &classNameMismatchError{requested: name, declared: class.Name()},
		}
		dr.state = defFailed
		dr.failure = failure
		dr.cond.Broadcast()
		dr.mu.Unlock()
		return rtda.NewFailureResult(failure)
	}

	// 4. Transition to defDefining and install.
	dr.state = defDefining
	dr.mu.Unlock()

	return cl.registerDefinition(class, name, dr)
}

// defineClassDirect is retained as a compatibility shim for existing
// same-package tests. New production and test code should use DefineClassResult.
func (cl *ClassLoader) defineClassDirect(name string, class *rtda.Class) rtda.ClassLoadResult {
	return cl.DefineClassResult(name, class)
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
