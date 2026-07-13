package rtda

import (
	"sync/atomic"

	"catty/classfile"
)

// Class is the runtime representation of a loaded class (JVMS §2.5.1 method area
// metadata). It is built by the classloader from a classfile.ClassFile. Core
// classes with no on-disk class file (java.lang.Object, System, ...) are
// synthesized directly as Class values by the native package.
type Class struct {
	name           string
	superName      string
	superClass     *Class
	interfaceNames []string
	interfaces     []*Class
	accessFlags    uint16
	cp             *classfile.ConstantPool

	instanceFields []*Field
	instCellCount  uint // total instance cell count (own + inherited), 1 per field per ADR-0030
	staticFields   []*Field
	staticCells    []HeapCell // heap-cell storage for static fields per ADR-0030

	methods     []*Method
	methodTable map[string]*Method // key = name + descriptor

	// Array support. A class with isArray is an array type; componentClass is the
	// element class (for object arrays) and componentKind tags primitive arrays.
	isArray        bool
	componentClass *Class
	componentKind  int // kindByte, kindChar, ...; 0 for object arrays
	arrayClass     *Class // cached "[Lthis;" / "[Ithis" array class

	// Class initialization bookkeeping (ADR-0025: Java 25 single-execution-context
	// state machine).
	initState int32  // one of the four init* constants
	initOwner uint64 // identity of the execution context currently initializing this class (0 = none)

	// classObject is the canonical java.lang.Class mirror for this class, created
	// lazily with a compare-and-swap so all goroutines see the same identity.
	// The field stores nil until the first request triggers lazy materialization.
	classObject atomic.Pointer[Object]
}

// Class initialization states (JVMS §5.5 via ADR-0025).
const (
	initNotStarted int32 = iota // not-initialized
	initInProgress              // initializing — initOwner names the owning execution context
	initInitialized             // successfully initialized
	initErroneous               // initialization failed — class is erroneous
)

const (
	kindNone int = iota
	kindBoolean
	kindByte
	kindChar
	kindShort
	kindInt
	kindLong
	kindFloat
	kindDouble
)

// --- Accessors used by the loader and interpreter ---

func (c *Class) Name() string         { return c.name }
func (c *Class) SuperClass() *Class   { return c.superClass }
func (c *Class) AccessFlags() uint16  { return c.accessFlags }
func (c *Class) ConstantPool() *classfile.ConstantPool { return c.cp }
func (c *Class) InstCellCount() uint    { return c.instCellCount }
func (c *Class) IsArray() bool          { return c.isArray }

// --- Typed static cell accessors (ADR-0030) ---
// staticCells is unexported; all external access goes through these methods.

// GetStaticInt returns the static cell at slotID as int32.
func (c *Class) GetStaticInt(slotID uint) int32 { return c.staticCells[slotID].GetInt() }

// SetStaticInt stores v into the static cell at slotID.
func (c *Class) SetStaticInt(slotID uint, v int32) { c.staticCells[slotID].SetInt(v) }

// GetStaticLong returns the static cell at slotID as int64.
func (c *Class) GetStaticLong(slotID uint) int64 { return c.staticCells[slotID].GetLong() }

// SetStaticLong stores v into the static cell at slotID.
func (c *Class) SetStaticLong(slotID uint, v int64) { c.staticCells[slotID].SetLong(v) }

// GetStaticFloat returns the static cell at slotID as float32.
func (c *Class) GetStaticFloat(slotID uint) float32 { return c.staticCells[slotID].GetFloat() }

// SetStaticFloat stores v into the static cell at slotID.
func (c *Class) SetStaticFloat(slotID uint, v float32) { c.staticCells[slotID].SetFloat(v) }

// GetStaticDouble returns the static cell at slotID as float64.
func (c *Class) GetStaticDouble(slotID uint) float64 { return c.staticCells[slotID].GetDouble() }

// SetStaticDouble stores v into the static cell at slotID.
func (c *Class) SetStaticDouble(slotID uint, v float64) { c.staticCells[slotID].SetDouble(v) }

// GetStaticRef returns the static cell at slotID as an object reference.
func (c *Class) GetStaticRef(slotID uint) *Object { return c.staticCells[slotID].GetRef() }

// StaticCellToSlot returns the static cell at slotID as a frame Slot for AOT
// bridge interop. For long/double the result is truncated (Slot.num is int32);
// callers must use GetStaticLong/GetStaticDouble for full 64-bit values.
func (c *Class) StaticCellToSlot(slotID uint, desc string) Slot {
	return c.staticCells[slotID].ToSlot(desc)
}

// IsInterface / IsAbstract etc.
func (c *Class) IsInterface() bool { return c.accessFlags&accInterface != 0 }
func (c *Class) IsAbstract() bool  { return c.accessFlags&accAbstract != 0 }

func (c *Class) ComponentClass() *Class { return c.componentClass }

// ClassObject returns the canonical java.lang.Class mirror for this class, creating
// it lazily via CAS on first access. All callers see the same Object identity, so
// obj.getClass() == obj.getClass() holds even across goroutines (ADR-0029).
// The caller must provide a factory that allocates a java.lang.Class Object and
// sets its Extra to this *Class.
func (c *Class) ClassObject(factory func() *Object) *Object {
	if obj := c.classObject.Load(); obj != nil {
		return obj
	}
	obj := factory()
	if obj == nil {
		return nil
	}
	obj.SetExtra(c)
	if c.classObject.CompareAndSwap(nil, obj) {
		return obj
	}
	// Lost the race — return the winner.
	return c.classObject.Load()
}

// --- Method lookup ---

func methodKey(name, descriptor string) string { return name + descriptor }

// AddMethod registers a method in both the list and the lookup table.
func (c *Class) AddMethod(m *Method) {
	c.methods = append(c.methods, m)
	if c.methodTable == nil {
		c.methodTable = make(map[string]*Method)
	}
	c.methodTable[methodKey(m.name, m.descriptor)] = m
}

// GetMethod looks up a method by name+descriptor declared in this class only.
func (c *Class) GetMethod(name, descriptor string) *Method {
	return c.methodTable[methodKey(name, descriptor)]
}

// Methods returns every method declared on this class (used by the lowering
// pass's test harness and future tooling).
func (c *Class) Methods() []*Method { return c.methods }

// LookupMethod walks the class hierarchy (then interfaces) for a method, used by
// invokevirtual/invokespecial/invokestatic resolution.
func (c *Class) LookupMethod(name, descriptor string) *Method {
	for cls := c; cls != nil; cls = cls.superClass {
		if m := cls.methodTable[methodKey(name, descriptor)]; m != nil {
			return m
		}
	}
	// interfaces (superinterfaces of c and its supers)
	for cls := c; cls != nil; cls = cls.superClass {
		for _, iface := range cls.interfaces {
			if m := lookupInterfaceMethod(iface, name, descriptor); m != nil {
				return m
			}
		}
	}
	return nil
}

func lookupInterfaceMethod(iface *Class, name, descriptor string) *Method {
	if m := iface.methodTable[methodKey(name, descriptor)]; m != nil {
		return m
	}
	for _, sub := range iface.interfaces {
		if m := lookupInterfaceMethod(sub, name, descriptor); m != nil {
			return m
		}
	}
	return nil
}

// --- Field lookup ---

// LookupField walks the hierarchy for an instance/static field by name+descriptor.
func (c *Class) LookupField(name, descriptor string) *Field {
	for cls := c; cls != nil; cls = cls.superClass {
		for _, f := range cls.instanceFields {
			if f.name == name && f.descriptor == descriptor {
				return f
			}
		}
		for _, f := range cls.staticFields {
			if f.name == name && f.descriptor == descriptor {
				return f
			}
		}
	}
	return nil
}

// StaticField looks up only this class's static fields by name+descriptor.
func (c *Class) StaticField(name, descriptor string) *Field {
	for _, f := range c.staticFields {
		if f.name == name && f.descriptor == descriptor {
			return f
		}
	}
	return nil
}

// --- Type compatibility (instanceof / checkcast / array covariance) ---

// DefaultBearingSuperInterfaces returns the superinterfaces of c that declare at
// least one non-abstract, non-static method ("default methods"), collected
// recursively and in JVMS §5.5 order (depth-first, left-to-right).
// Each interface appears at most once.
func (c *Class) DefaultBearingSuperInterfaces(seen map[string]bool) []*Class {
	var result []*Class
	for _, iface := range c.interfaces {
		collectDefaultBearing(iface, seen, &result)
	}
	return result
}

func collectDefaultBearing(iface *Class, seen map[string]bool, result *[]*Class) {
	if seen[iface.name] {
		return
	}
	// Recurse into superinterfaces first (JVMS §5.5 order).
	for _, super := range iface.interfaces {
		collectDefaultBearing(super, seen, result)
	}
	// Then this interface if it declares any default method.
	if !seen[iface.name] && ifaceHasDefaultMethod(iface) {
		seen[iface.name] = true
		*result = append(*result, iface)
	}
}

// ifaceHasDefaultMethod reports whether an interface declares at least one
// non-abstract, non-static method (a "default method" as per JVMS §5.5).
func ifaceHasDefaultMethod(iface *Class) bool {
	for _, m := range iface.methods {
		if m.accessFlags&(accAbstract|accStatic) == 0 {
			return true
		}
	}
	return false
}

// isAssignableFrom reports whether a value of class c can be assigned to a
// variable of type target (the inverse direction of instanceof).
func (c *Class) isAssignableFrom(target *Class) bool {
	if target == nil {
		return false
	}
	if c == target {
		return true
	}
	// Everything is assignable to java.lang.Object.
	if c.name == "java/lang/Object" {
		return true
	}
	// Array covariance: an array is assignable to Object, to Cloneable, to
	// java.io.Serializable, and to a compatible array type.
	if target.isArray {
		return c.isArray && c.componentClass != nil &&
			c.componentClass.isAssignableFrom(target.componentClass)
	}
	// Superclass chain.
	for sc := target.superClass; sc != nil; sc = sc.superClass {
		if sc == c {
			return true
		}
	}
	// Interface implementation.
	return c.IsInterface() && target.implementsInterface(c)
}

func (c *Class) implementsInterface(target *Class) bool {
	for cls := c; cls != nil; cls = cls.superClass {
		for _, iface := range cls.interfaces {
			if iface == target || iface.implementsInterface(target) {
				return true
			}
		}
	}
	return false
}
