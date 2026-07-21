package rtda

import "sync"

// Canonical VM primitive type and void singletons.
// These Classes have VMIdentity as their defining loader and are shared
// across all initiating loaders (ADR-0033).
var (
	VMPrimitiveBool   *Class
	VMPrimitiveByte   *Class
	VMPrimitiveChar   *Class
	VMPrimitiveShort  *Class
	VMPrimitiveInt    *Class
	VMPrimitiveLong   *Class
	VMPrimitiveFloat  *Class
	VMPrimitiveDouble *Class
	VMVoid            *Class

	vmTypesOnce sync.Once
)

// InitVMTypes initializes the canonical VM primitive and void Classes.
// Safe to call multiple times (once.Do). Must be called before any class
// loading that might trigger primitive/void identity lookup.
func InitVMTypes() {
	vmTypesOnce.Do(func() {
		VMPrimitiveBool = makeVMPrimitive("boolean")
		VMPrimitiveByte = makeVMPrimitive("byte")
		VMPrimitiveChar = makeVMPrimitive("char")
		VMPrimitiveShort = makeVMPrimitive("short")
		VMPrimitiveInt = makeVMPrimitive("int")
		VMPrimitiveLong = makeVMPrimitive("long")
		VMPrimitiveFloat = makeVMPrimitive("float")
		VMPrimitiveDouble = makeVMPrimitive("double")
		VMVoid = makeVMPrimitive("void")
	})
}

func makeVMPrimitive(name string) *Class {
	info := primitiveInfo[name]
	c := &Class{
		name:           name,
		definingLoader: VMIdentity,
		componentKind:  info.kind,
		methodTable:    make(map[string]*Method),
	}
	c.initCond = sync.NewCond(&c.initMu)
	// Primitive and void types are always already initialized.
	c.initState = initInitialized
	return c
}

// VMPrimitiveForKind returns the canonical VM Class for a primitive kind,
// or nil for kindNone. Safe to call before InitVMTypes (self-initialising).
func VMPrimitiveForKind(kind int) *Class {
	InitVMTypes()
	switch kind {
	case kindBoolean:
		return VMPrimitiveBool
	case kindByte:
		return VMPrimitiveByte
	case kindChar:
		return VMPrimitiveChar
	case kindShort:
		return VMPrimitiveShort
	case kindInt:
		return VMPrimitiveInt
	case kindLong:
		return VMPrimitiveLong
	case kindFloat:
		return VMPrimitiveFloat
	case kindDouble:
		return VMPrimitiveDouble
	default:
		return nil
	}
}

// VMPrimitiveForName returns the canonical VM Class for a primitive type
// name ("int", "boolean", "void", etc.), or nil if name is not a primitive.
// Safe to call before InitVMTypes (self-initialising).
func VMPrimitiveForName(name string) *Class {
	InitVMTypes()
	switch name {
	case "boolean":
		return VMPrimitiveBool
	case "byte":
		return VMPrimitiveByte
	case "char":
		return VMPrimitiveChar
	case "short":
		return VMPrimitiveShort
	case "int":
		return VMPrimitiveInt
	case "long":
		return VMPrimitiveLong
	case "float":
		return VMPrimitiveFloat
	case "double":
		return VMPrimitiveDouble
	case "void":
		return VMVoid
	default:
		return nil
	}
}

// IsVMPrimitive reports whether name is a primitive type or void.
func IsVMPrimitive(name string) bool {
	switch name {
	case "boolean", "byte", "char", "short", "int", "long", "float", "double", "void":
		return true
	default:
		return false
	}
}

// GetArrayClass returns the canonical array class for this component class.
// For reference arrays, the component Class owns the array class identity.
// For primitive arrays, the VM primitive Class owns it.
//
// The array class name is a valid JVM descriptor:
//   - primitive component: "[" + descriptor byte (e.g. "[I" for int[])
//   - reference component: "[L" + name + ";" (e.g. "[Ljava/lang/String;")
//   - array component:     "[" + name        (e.g. "[[I" for int[][])
//
// Thread-safe via CAS on arrayClass.
func (c *Class) GetArrayClass() *Class {
	// Fast path: already cached.
	if arr := c.arrayClass.Load(); arr != nil {
		return arr
	}

	// Build the correct JVM descriptor for this array type.
	var name string
	if desc := PrimitiveDescriptor(c.name); desc != 0 {
		// Primitive component: "[" + descriptor byte.
		name = "[" + string(desc)
	} else if c.isArray {
		// Array component: "[" + existing descriptor.
		name = "[" + c.name
	} else {
		// Reference (object) component: "[L" + binary name + ";".
		name = "[L" + c.name + ";"
	}

	// Primitive arrays inherit the component's kind; reference and array
	// arrays get kindNone (elements are references, not primitives).
	kind := kindNone
	if !c.isArray {
		kind = c.componentKind
	}

	arr := &Class{
		name:              name,
		isArray:           true,
		componentClass:    c,
		componentKind:     kind,
		definingLoader:    c.definingLoader,
		definingLoaderRef: c.definingLoaderRef,
		methodTable:       make(map[string]*Method),
	}
	arr.initCond = sync.NewCond(&arr.initMu)

	// CAS publication — exactly one array Class identity per component.
	if c.arrayClass.CompareAndSwap(nil, arr) {
		return arr
	}
	return c.arrayClass.Load()
}
