package rtda

import "catty/classfile"

// NewClass builds a runtime Class from a parsed class file, resolving the
// superclass and interfaces through loader (which avoids an import cycle with
// the classloader package). Field slot layout is computed here:
//   - instance fields are laid out contiguously after the superclass's instance
//     slots (so subclasses reuse inherited offsets);
//   - static fields get their own slots in staticVars;
//   - long/double occupy two slots everywhere.
func NewClass(cf *classfile.ClassFile, loader Loader) *Class {
	c := &Class{
		name:           cf.ClassName(),
		superName:      cf.SuperClassName(),
		accessFlags:    cf.AccessFlags(),
		cp:             cf.ConstantPool(),
		interfaceNames: cf.InterfaceNames(),
		methodTable:    make(map[string]*Method),
	}

	if c.superName != "" {
		c.superClass = loader.LoadClass(c.superName)
		c.instCellCount = c.superClass.instCellCount
	}
	c.interfaces = make([]*Class, len(c.interfaceNames))
	for i, n := range c.interfaceNames {
		c.interfaces[i] = loader.LoadClass(n)
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

// InitState returns the class's four-state initialization value (initNotStarted,
// initInProgress, initInitialized, or initErroneous).
func (c *Class) InitState() int32 { return c.initState }

// IsInitialized reports whether the class has completed initialization.
func (c *Class) IsInitialized() bool { return c.initState == initInitialized }

// InitOwner returns the identity of the execution context that is currently
// initializing this class, or 0 if no context owns the initializing state.
func (c *Class) InitOwner() uint64 { return c.initOwner }

// SetInitOwner records the execution context that owns the current initializing
// state. It is only meaningful when initState == initInProgress.
func (c *Class) SetInitOwner(owner uint64) { c.initOwner = owner }

// MarkInitInProgress transitions the state from not-initialized to initializing.
// Returns true if the transition succeeded (caller now owns initialization).
func (c *Class) MarkInitInProgress(owner uint64) bool {
	if c.initState != initNotStarted {
		return false
	}
	c.initState = initInProgress
	c.initOwner = owner
	return true
}

// MarkInitialized transitions the state from initializing to initialized.
func (c *Class) MarkInitialized() {
	c.initState = initInitialized
	c.initOwner = 0
}

// MarkErroneous transitions the state from initializing to erroneous.
func (c *Class) MarkErroneous() {
	c.initState = initErroneous
	c.initOwner = 0
}

// InitStarted is the legacy accessor; it returns true for any state past
// not-initialized. Kept for compatibility with existing callers that only need
// to know whether init has been attempted.
func (c *Class) InitStarted() bool { return c.initState != initNotStarted }

// MarkInitStarted is the legacy mutator — its only remaining safe use is to
// set the state to initializing before the shared service takes over.
// Prefer MarkInitInProgress in new code.
func (c *Class) MarkInitStarted() {
	if c.initState == initNotStarted {
		c.initState = initInProgress
	}
}

// NewArrayClass builds the runtime class for an array type name ("[I",
// "[Ljava/lang/Object;", "[[C", ...). Component resolution for object and array
// components goes through loader so each element type is cached exactly once.
func NewArrayClass(name string, loader Loader) *Class {
	c := &Class{name: name, isArray: true, methodTable: make(map[string]*Method)}
	comp := name[1:] // drop leading '['
	switch comp[0] {
	case 'L':
		c.componentClass = loader.LoadClass(comp[1 : len(comp)-1])
	case '[':
		c.componentClass = loader.LoadClass(comp)
	case 'I':
		c.componentKind = kindInt
	case 'J':
		c.componentKind = kindLong
	case 'F':
		c.componentKind = kindFloat
	case 'D':
		c.componentKind = kindDouble
	case 'B':
		c.componentKind = kindByte
	case 'C':
		c.componentKind = kindChar
	case 'S':
		c.componentKind = kindShort
	case 'Z':
		c.componentKind = kindBoolean
	}
	return c
}

// ComponentKind returns the primitive element kind (0 for object arrays).
func (c *Class) ComponentKind() int { return c.componentKind }
