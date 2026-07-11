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
		c.instSlotCount = c.superClass.instSlotCount
	}
	c.interfaces = make([]*Class, len(c.interfaceNames))
	for i, n := range c.interfaceNames {
		c.interfaces[i] = loader.LoadClass(n)
	}

	for _, f := range cf.Fields() {
		if f.AccessFlags()&accStatic != 0 {
			slot := uint(len(c.staticVars))
			c.staticVars = append(c.staticVars, make([]Slot, fieldSlotSize(f.Descriptor()))...)
			c.staticFields = append(c.staticFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), true, slot))
		} else {
			slot := c.instSlotCount
			c.instSlotCount += uint(fieldSlotSize(f.Descriptor()))
			c.instanceFields = append(c.instanceFields,
				NewField(c, f.Name(), f.Descriptor(), f.AccessFlags(), false, slot))
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

// fieldSlotSize returns 2 for long/double, 1 for everything else.
func fieldSlotSize(descriptor string) int {
	if descriptor == "J" || descriptor == "D" {
		return 2
	}
	return 1
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
		c.instSlotCount = super.instSlotCount
	}
	return c
}

// AddStaticField declares a static field, allocating slots in staticVars, and
// returns the field's slotID so the builder can initialize it.
func (c *Class) AddStaticField(name, descriptor string) *Field {
	slot := uint(len(c.staticVars))
	c.staticVars = append(c.staticVars, make([]Slot, fieldSlotSize(descriptor))...)
	f := NewField(c, name, descriptor, accStatic, true, slot)
	c.staticFields = append(c.staticFields, f)
	return f
}

func (c *Class) SetStaticRef(slotID uint, ref *Object) {
	c.staticVars[slotID].ref = ref
}

// AddInstanceField declares an instance field on a synthetic class, allocating
// the next slot in the layout. Used by native exception classes (Throwable's
// detailMessage).
func (c *Class) AddInstanceField(name, descriptor string) *Field {
	slot := c.instSlotCount
	c.instSlotCount += uint(fieldSlotSize(descriptor))
	f := NewField(c, name, descriptor, 0, false, slot)
	c.instanceFields = append(c.instanceFields, f)
	return f
}

func (c *Class) SetSuper(super *Class) { c.superClass = super }

// --- Class initialization (<clinit>) bookkeeping ---

// InitStarted/InitDone track whether <clinit> has run. The interpreter triggers
// initialization lazily at the JVMS §5.5 events (new/getstatic/putstatic/
// invokestatic); until then a loaded class is linked but not initialized.
func (c *Class) InitStarted() bool { return c.initStarted }
func (c *Class) MarkInitStarted()  { c.initStarted = true }

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
