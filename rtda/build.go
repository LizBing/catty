package rtda

import (
	"fmt"
	"sync"

	"catty/classfile"
)

// BuildResult is the typed return from BuildClass.
type BuildResult struct {
	Class   *Class
	Failure *ClassLoadFailure
}

// BuildClass builds a runtime Class from a parsed class file using the typed
// Loader interface. Superclass and interface resolution failures are propagated
// as typed results rather than compressed into a generic nil.
//
// NewClass is retained as a legacy must-load wrapper for bootstrap callers
// that cannot tolerate typed failures (synthetic class builders).
func BuildClass(cf *classfile.ClassFile, loader Loader) BuildResult {
	c := &Class{
		name:             cf.ClassName(),
		superName:        cf.SuperClassName(),
		accessFlags:      cf.AccessFlags(),
		cp:               cf.ConstantPool(),
		bootstrapMethods: cf.BootstrapMethods(),
		interfaceNames:   cf.InterfaceNames(),
		methodTable:      make(map[string]*Method),
	}
	c.initCond = sync.NewCond(&c.initMu)

	if c.superName != "" {
		result := loader.LoadClassResult(c.superName)
		if !result.IsSuccess() {
			f := result.Failure()
			return BuildResult{Failure: &ClassLoadFailure{
				Kind:  f.Kind,
				Name:  cf.ClassName(),
				Cause: fmt.Errorf("superclass %s: %w", c.superName, f),
			}}
		}
		c.superClass = result.Class()
		c.instCellCount = c.superClass.instCellCount
	}
	c.interfaces = make([]*Class, len(c.interfaceNames))
	for i, n := range c.interfaceNames {
		result := loader.LoadClassResult(n)
		if !result.IsSuccess() {
			f := result.Failure()
			return BuildResult{Failure: &ClassLoadFailure{
				Kind:  f.Kind,
				Name:  cf.ClassName(),
				Cause: fmt.Errorf("interface %s: %w", n, f),
			}}
		}
		c.interfaces[i] = result.Class()
	}

	for _, f := range cf.Fields() {
		if f.AccessFlags()&accStatic != 0 {
			cellID := uint(len(c.staticCells))
			c.staticCells = append(c.staticCells, HeapCell{})
			c.staticFields = append(c.staticFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), true, cellID))
		} else {
			cellID := c.instCellCount
			c.instCellCount++
			c.instanceFields = append(c.instanceFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), false, cellID))
		}
	}

	for _, m := range cf.Methods() {
		code := m.Code()
		var exTable []exceptionEntry
		var maxStack, maxLocals uint
		var bytecode []byte
		if code != nil {
			maxStack = uint(code.MaxStack())
			maxLocals = uint(code.MaxLocals())
			bytecode = code.Code()
			exTable = convertExceptionTable(code.ExceptionTable(), c.cp)
		}
		method := InterpretedMethod(c, m.Name(), m.Descriptor(),
			m.AccessFlags(), maxStack, maxLocals, bytecode, exTable)
		if code != nil {
			method.SetStackMap(code.StackMapTable())
		}
		c.AddMethod(method)
	}
	return BuildResult{Class: c}
}

// NewClass builds a runtime Class from a parsed class file, resolving the
// superclass and interfaces through loader (which avoids an import cycle with
// the classloader package). Field slot layout is computed here:
//   - instance fields are laid out contiguously after the superclass's instance
//     slots (so subclasses reuse inherited offsets);
//   - static fields get their own slots in staticVars;
//   - long/double occupy two slots everywhere.
//
// Returns nil when superclass or interface resolution fails.
// Prefer BuildClass for typed failure propagation; NewClass is kept for
// legacy bootstrap callers that cannot tolerate typed failures.
func NewClass(cf *classfile.ClassFile, loader Loader) *Class {
	c := &Class{
		name:             cf.ClassName(),
		superName:        cf.SuperClassName(),
		accessFlags:      cf.AccessFlags(),
		cp:               cf.ConstantPool(),
		bootstrapMethods: cf.BootstrapMethods(),
		interfaceNames:   cf.InterfaceNames(),
		methodTable:      make(map[string]*Method),
	}
	c.initCond = sync.NewCond(&c.initMu)

	if c.superName != "" {
		c.superClass = loader.LoadClass(c.superName)
		if c.superClass == nil {
			return nil // superclass resolution failed
		}
		c.instCellCount = c.superClass.instCellCount
	}
	c.interfaces = make([]*Class, len(c.interfaceNames))
	for i, n := range c.interfaceNames {
		c.interfaces[i] = loader.LoadClass(n)
		if c.interfaces[i] == nil {
			return nil // interface resolution failed
		}
	}

	for _, f := range cf.Fields() {
		// ADR-0030: every field occupies exactly one heap cell regardless of type.
		if f.AccessFlags()&accStatic != 0 {
			cellID := uint(len(c.staticCells))
			c.staticCells = append(c.staticCells, HeapCell{})
			c.staticFields = append(c.staticFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), true, cellID))
		} else {
			cellID := c.instCellCount
			c.instCellCount++
			c.instanceFields = append(c.instanceFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), false, cellID))
		}
	}

	for _, m := range cf.Methods() {
		code := m.Code()
		var exTable []exceptionEntry
		var maxStack, maxLocals uint
		var bytecode []byte
		if code != nil {
			maxStack = uint(code.MaxStack())
			maxLocals = uint(code.MaxLocals())
			bytecode = code.Code()
			exTable = convertExceptionTable(code.ExceptionTable(), c.cp)
		}
		method := InterpretedMethod(c, m.Name(), m.Descriptor(),
			m.AccessFlags(), maxStack, maxLocals, bytecode, exTable)
		if code != nil {
			method.SetStackMap(code.StackMapTable())
		}
		c.AddMethod(method)
	}
	return c
}

func convertExceptionTable(entries []*classfile.ExceptionTableEntry, cp *classfile.ConstantPool) []exceptionEntry {
	if len(entries) == 0 {
		return nil
	}
	out := make([]exceptionEntry, len(entries))
	for i, e := range entries {
		var catchName string
		if e.CatchType() != 0 {
			catchName = cp.ClassName(e.CatchType())
		}
		out[i] = exceptionEntry{
			startPc:   int(e.StartPc()),
			endPc:     int(e.EndPc()),
			handlerPc: int(e.HandlerPc()),
			catchType: catchName,
		}
	}
	return out
}

// --- Synthetic class construction (for core classes built natively, not from
// a class file: java.lang.Object/System/String/PrintStream/...) ---

// NewSyntheticClass creates an empty class with the given name and superclass.
// The native package adds fields and native methods to it.
func NewSyntheticClass(name string, super *Class) *Class {
	c := &Class{name: name, superClass: super, methodTable: make(map[string]*Method)}
	c.initCond = sync.NewCond(&c.initMu)
	if super != nil {
		c.instCellCount = super.instCellCount
	}
	return c
}

// AddStaticField declares a static field, allocating one heap cell, and
// returns the field's cell ID so the builder can initialize it.
func (c *Class) AddStaticField(name, descriptor string) *Field {
	cellID := uint(len(c.staticCells))
	c.staticCells = append(c.staticCells, HeapCell{})
	f := NewField(c, name, descriptor, accStatic, true, cellID)
	c.staticFields = append(c.staticFields, f)
	return f
}

func (c *Class) SetStaticRef(cellID uint, ref *Object) {
	c.staticCells[cellID].SetRef(ref)
}

// AddInstanceField declares an instance field on a synthetic class, allocating
// the next cell in the layout. Used by native exception classes (Throwable's
// detailMessage).
func (c *Class) AddInstanceField(name, descriptor string) *Field {
	cellID := c.instCellCount
	c.instCellCount++
	f := NewField(c, name, descriptor, 0, false, cellID)
	c.instanceFields = append(c.instanceFields, f)
	return f
}

func (c *Class) SetSuper(super *Class) { c.superClass = super }

// MarkInterface sets the ACC_INTERFACE flag so IsInterface() returns true.
// Used when building synthetic interface classes in tests.
func (c *Class) MarkInterface() { c.accessFlags |= accInterface }

// AddInterface adds a direct superinterface to a synthetic class.
// Used when building synthetic class hierarchies in tests.
func (c *Class) AddInterface(iface *Class) {
	c.interfaces = append(c.interfaces, iface)
}

// --- Class initialization (<clinit>) bookkeeping (ADR-0025) ---
//
// All accessors are guarded by initMu (ADR-0029). The lock is distinct from the
// Class mirror's Java monitor and protects initState, initOwner, and initCond.
// InitializeClass itself holds initMu across the full protocol, so the
// Mark*/Init* methods below are safe for external callers but NOT called
// internally by InitializeClass (which manipulates fields directly).

// InitState returns the class's four-state initialization value (initNotStarted,
// initInProgress, initInitialized, or initErroneous). Lock-guarded.
func (c *Class) InitState() int32 {
	c.initMu.Lock()
	s := c.initState
	c.initMu.Unlock()
	return s
}

// IsInitialized reports whether the class has completed initialization.
// Lock-guarded — establishes acquire visibility for the published state,
// including any clinit-written static fields.
func (c *Class) IsInitialized() bool {
	c.initMu.Lock()
	ok := c.initState == initInitialized
	c.initMu.Unlock()
	return ok
}

// InitOwner returns the identity of the execution context that is currently
// initializing this class, or 0 if no context owns the initializing state.
// Lock-guarded.
func (c *Class) InitOwner() uint64 {
	c.initMu.Lock()
	o := c.initOwner
	c.initMu.Unlock()
	return o
}

// SetInitOwner records the execution context that owns the current initializing
// state. It is only meaningful when initState == initInProgress.
// Lock-guarded.
func (c *Class) SetInitOwner(owner uint64) {
	c.initMu.Lock()
	c.initOwner = owner
	c.initMu.Unlock()
}

// MarkInitInProgress transitions the state from not-initialized to initializing.
// Returns true if the transition succeeded (caller now owns initialization).
// Lock-guarded.
func (c *Class) MarkInitInProgress(owner uint64) bool {
	c.initMu.Lock()
	defer c.initMu.Unlock()
	if c.initState != initNotStarted {
		return false
	}
	c.initState = initInProgress
	c.initOwner = owner
	return true
}

// MarkInitialized transitions the state from initializing to initialized and
// broadcasts on initCond to wake all waiters (ADR-0029 terminal publication).
// Lock-guarded.
func (c *Class) MarkInitialized() {
	c.initMu.Lock()
	c.initState = initInitialized
	c.initOwner = 0
	c.initCond.Broadcast()
	c.initMu.Unlock()
}

// MarkErroneous transitions the state from initializing to erroneous and
// broadcasts on initCond to wake all waiters (ADR-0029 terminal publication).
// Lock-guarded.
func (c *Class) MarkErroneous() {
	c.initMu.Lock()
	c.initState = initErroneous
	c.initOwner = 0
	c.initCond.Broadcast()
	c.initMu.Unlock()
}

// InitStarted is the legacy accessor; it returns true for any state past
// not-initialized. Lock-guarded. Kept for compatibility with existing callers
// that only need to know whether init has been attempted.
func (c *Class) InitStarted() bool {
	c.initMu.Lock()
	ok := c.initState != initNotStarted
	c.initMu.Unlock()
	return ok
}

// MarkInitStarted is the legacy mutator — its only remaining safe use is to
// set the state to initializing before the shared service takes over.
// Lock-guarded. Prefer MarkInitInProgress in new code.
func (c *Class) MarkInitStarted() {
	c.initMu.Lock()
	if c.initState == initNotStarted {
		c.initState = initInProgress
	}
	c.initMu.Unlock()
}

// NewArrayClass builds the runtime class for an array type name ("[I",
// "[Ljava/lang/Object;", "[[C", ...). Component resolution for object and array
// components goes through loader so each element type is cached exactly once.
//
// The resulting array Class derives its defining-loader identity from its
// component type (ADR-0033): reference arrays inherit the component's loader;
// primitive arrays use VMIdentity.
//
// Returns nil when component class resolution fails.
func NewArrayClass(name string, loader Loader) *Class {
	comp := name[1:] // drop leading '['

	// Resolve component class.
	var component *Class
	switch comp[0] {
	case 'L':
		component = loader.LoadClass(comp[1 : len(comp)-1])
	case '[':
		component = loader.LoadClass(comp)
	default:
		// Primitive component — use VM canonical identity.
		component = vmPrimitiveForDescriptor(comp[0])
	}
	if component == nil {
		return nil
	}

	// Delegate to component-owned array creation so the array Class
	// identity is canonical and derived from the component.
	return component.GetArrayClass()
}

// vmPrimitiveForDescriptor returns the canonical VM Class for a single
// primitive array component descriptor character.
func vmPrimitiveForDescriptor(ch byte) *Class {
	switch ch {
	case 'I':
		return VMPrimitiveInt
	case 'J':
		return VMPrimitiveLong
	case 'F':
		return VMPrimitiveFloat
	case 'D':
		return VMPrimitiveDouble
	case 'B':
		return VMPrimitiveByte
	case 'C':
		return VMPrimitiveChar
	case 'S':
		return VMPrimitiveShort
	case 'Z':
		return VMPrimitiveBool
	}
	return nil
}

// ComponentKind returns the primitive element kind (0 for object arrays).
func (c *Class) ComponentKind() int { return c.componentKind }
