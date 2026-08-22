package kernel

import (
	"fmt"
	"sync"

	"catty/internal/classfile"
)

// ClassState implements the JVMS §5.5 initialization state machine (M0 v0:
// single-threaded semantics; cross-thread circular-init deadlock breaking
// is deferred — see deviation ledger).
type ClassState int8

const (
	StateDefined ClassState = iota // loaded, not yet initialized
	StateInitializing
	StateInitialized
	StateErroneous
)

// Method is a resolved method of a class: either interpreted (CF != nil,
// Code != nil) or native (Native != nil).
type Method struct {
	Holder *Class
	Name   string
	Desc   string
	Flags  uint16

	CF     *classfile.ClassFile // interpreted: owner's constant pool
	Code   *classfile.Code      // interpreted: bytecode
	Native NativeFunc           // synthesized implementation
}

// Static reports whether the method is static.
func (m *Method) Static() bool { return m.Flags&classfile.AccStatic != 0 }

// Key is the map lookup key for members.
func memberKey(name, desc string) string { return name + "|" + desc }

// Field is a resolved field. Instance fields carry an absolute slot in the
// object layout (supers first); static fields carry an index into the
// declaring class's Statics.
type Field struct {
	Holder     *Class
	Name       string
	Desc       string
	Flags      uint16
	Slot       int  // instance slot (valid if !Static)
	Static     bool
	StaticSlot int  // index into Holder.Statics (valid if Static)
}

// Class is runtime class metadata.
type Class struct {
	Name  string
	Super *Class
	Ifaces []*Class
	Flags uint16

	Methods    []*Method
	OwnFields  []*Field // fields declared by this class itself (non-static)
	Statics    []Value  // static storage declared here

	CF   *classfile.ClassFile // nil for synthesized classes
	def  *ClassDef            // synthesized origin (nil for CF classes)

	State  ClassState
	initMu sync.Mutex

	// Lookup caches cover this class plus all ancestors (flat view).
	methodsByKey map[string]*Method
	fieldsByKey  map[string]*Field

	// Resolved constant-pool entries for interpreted classes. Key space
	// is owned by callers (VM: kind<<16|index). Guarded by resolveMu.
	resolveMu sync.Mutex
	resolved  map[uint32]any

	layoutSize int // total instance slots including supers

	IsArray  bool
	CompDesc string // array component descriptor (arrays only)
}

// Synthetic reports whether the class was defined by the runtime rather
// than loaded from a class file.
func (c *Class) Synthetic() bool { return c.CF == nil }

// FindMethod looks up name+desc on this class only.
func (c *Class) FindMethod(name, desc string) *Method {
	return c.methodsByKey[memberKey(name, desc)]
}

// FindField looks up name+desc on this class only (instance or static).
func (c *Class) FindField(name, desc string) *Field {
	return c.fieldsByKey[memberKey(name, desc)]
}

// IsSubclassOf reports whether c is anc or inherits from it (class chain
// only, interfaces excluded — use IsInstance for full semantics).
func (c *Class) IsSubclassOf(anc *Class) bool {
	for x := c; x != nil; x = x.Super {
		if x == anc {
			return true
		}
	}
	return false
}

// ResolvedCache returns the resolution cache map (lazy init).
func (c *Class) resolvedCache() map[uint32]any {
	if c.resolved == nil {
		c.resolved = make(map[uint32]any)
	}
	return c.resolved
}

// CacheGet reads a constant-pool resolution cache entry (thread-safe).
// The key space is owned by callers (e.g. the VM's kind<<16|index scheme).
func (c *Class) CacheGet(key uint32) (any, bool) {
	c.resolveMu.Lock()
	defer c.resolveMu.Unlock()
	v, ok := c.resolvedCache()[key]
	return v, ok
}

// CacheSet writes a resolution cache entry (thread-safe).
func (c *Class) CacheSet(key uint32, v any) {
	c.resolveMu.Lock()
	defer c.resolveMu.Unlock()
	c.resolvedCache()[key] = v
}

// InitTracker lets the VM report in-progress class initialization so that
// recursive initialization from a class's own <clinit> proceeds per JVMS §5.5.
type InitTracker interface {
	IsInitializing(name string) bool
	BeginInit(name string)
	EndInit(name string)
}

// EnsureInitialized drives <clinit> execution: supers first, then self.
// Synthesized classes are pre-initialized at define time.
func (k *Kernel) EnsureInitialized(tracker InitTracker, c *Class) error {
	if c.State == StateInitialized {
		return nil
	}
	if tracker != nil && tracker.IsInitializing(c.Name) {
		return nil // recursive init from own <clinit>: proceed (JVMS §5.5)
	}

	c.initMu.Lock()
	defer c.initMu.Unlock()
	switch c.State {
	case StateInitialized:
		return nil
	case StateErroneous:
		return fmt.Errorf("class %s is in erroneous state from a failed <clinit>", c.Name)
	}

	if c.Super != nil {
		if err := k.EnsureInitialized(tracker, c.Super); err != nil {
			c.State = StateErroneous
			return err
		}
	}
	if tracker != nil {
		tracker.BeginInit(c.Name)
		defer tracker.EndInit(c.Name)
	}
	c.State = StateInitializing
	if m := c.FindMethod("<clinit>", "()V"); m != nil {
		if _, err := k.Invoke(m, nil, nil); err != nil {
			c.State = StateErroneous
			return err
		}
	}
	c.State = StateInitialized
	return nil
}
