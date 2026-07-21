package rtda

import (
	"sync"
	"sync/atomic"

	"catty/classfile"
)

// Class is the runtime representation of a loaded class (JVMS §2.5.1 method area
// metadata). It is built by the classloader from a classfile.ClassFile. Core
// classes with no on-disk class file (java.lang.Object, System, ...) are
// synthesized directly as Class values by the native package.
//
// Class identity is (definingLoader, name). The definingLoader is set exactly
// once during construction and is immutable thereafter. Primitive and void
// Classes use VMIdentity.
type Class struct {
	name              string
	definingLoader    *LoaderIdentity // immutable after construction; nil until bound
	definingLoaderRef Loader          // the defining Loader interface; nil for VM primitives
	superName         string
	superClass        *Class
	interfaceNames    []string
	interfaces        []*Class
	accessFlags       uint16
	cp                *classfile.ConstantPool
	bootstrapMethods  *classfile.BootstrapMethodsAttr

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
	componentKind  int                   // kindByte, kindChar, ...; 0 for object arrays
	arrayClass     atomic.Pointer[Class] // cached "[Lthis;" / "[Ithis" array class

	// Class initialization bookkeeping (ADR-0025: Java 25 state machine).
	// Protected by initMu. The lock is distinct from the Class mirror's Java
	// monitor (ADR-0029) and guards initState, initOwner, and initCond.
	initMu    sync.Mutex
	initCond  *sync.Cond // lazy-init via initCondOnce; signals terminal state transitions
	initState int32      // one of the four init* constants
	initOwner uint64     // identity of the execution context currently initializing this class (0 = none)

	// classObject is the canonical java.lang.Class mirror for this class, created
	// lazily with double-checked locking so all goroutines see the same identity.
	// The field stores nil until the first request triggers lazy materialization.
	// classObjectMu serialises factory invocation so at most one goroutine calls
	// the factory (prevents wasted allocation and ensures loader-protected creation).
	classObject   atomic.Pointer[Object]
	classObjectMu sync.Mutex
}

// Class initialization states (JVMS §5.5 via ADR-0025).
const (
	initNotStarted  int32 = iota // not-initialized
	initInProgress               // initializing — initOwner names the owning execution context
	initInitialized              // successfully initialized
	initErroneous                // initialization failed — class is erroneous
)

// BootstrapLoader provides the java/lang/Class class for VM primitive and
// void type mirrors which have no defining loader of their own. It is set
// exactly once by the launcher before any Class mirror is created.
// Use SetBootstrapLoader to set it; direct access is not allowed.
var (
	bootstrapLoader    Loader
	bootstrapLoaderMu  sync.Mutex
	bootstrapLoaderSet bool
)

// SetBootstrapLoader sets the bootstrap loader exactly once. Subsequent calls
// with the same loader are idempotent. A call with a different loader panics
// (invariant violation — the bootstrap loader is a VM-wide singleton).
func SetBootstrapLoader(l Loader) {
	bootstrapLoaderMu.Lock()
	defer bootstrapLoaderMu.Unlock()
	if bootstrapLoaderSet {
		if bootstrapLoader != l {
			panic("catty: SetBootstrapLoader: already set to a different loader")
		}
		return
	}
	bootstrapLoader = l
	bootstrapLoaderSet = true
}

// getBootstrapLoader returns the current bootstrap loader.
func getBootstrapLoader() Loader {
	bootstrapLoaderMu.Lock()
	defer bootstrapLoaderMu.Unlock()
	return bootstrapLoader
}

// resetBootstrapLoaderForTesting resets the bootstrap loader so tests may
// install a fresh loader for each sub-test. NOT for production use.
func resetBootstrapLoaderForTesting() {
	bootstrapLoaderMu.Lock()
	defer bootstrapLoaderMu.Unlock()
	bootstrapLoader = nil
	bootstrapLoaderSet = false
	// VM types are process singletons, so reset their already-materialized
	// mirrors as part of test isolation. Production never calls this helper.
	for _, class := range []*Class{
		VMPrimitiveBool, VMPrimitiveByte, VMPrimitiveChar, VMPrimitiveShort,
		VMPrimitiveInt, VMPrimitiveLong, VMPrimitiveFloat, VMPrimitiveDouble, VMVoid,
	} {
		if class == nil {
			continue
		}
		class.classObject.Store(nil)
		if arrayClass := class.arrayClass.Load(); arrayClass != nil {
			arrayClass.classObject.Store(nil)
		}
	}
}

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

// primitiveInfo maps a primitive type name to its kind and JVM descriptor byte.
var primitiveInfo = map[string]struct {
	kind int
	desc byte
}{
	"boolean": {kindBoolean, 'Z'},
	"byte":    {kindByte, 'B'},
	"char":    {kindChar, 'C'},
	"short":   {kindShort, 'S'},
	"int":     {kindInt, 'I'},
	"long":    {kindLong, 'J'},
	"float":   {kindFloat, 'F'},
	"double":  {kindDouble, 'D'},
	"void":    {0, 'V'},
}

// PrimitiveDescriptor returns the JVM descriptor byte for a primitive type
// name, or 0 if name is not a primitive.
func PrimitiveDescriptor(name string) byte {
	if info, ok := primitiveInfo[name]; ok {
		return info.desc
	}
	return 0
}

// --- Accessors used by the loader and interpreter ---

func (c *Class) Name() string                                      { return c.name }
func (c *Class) DefiningLoader() *LoaderIdentity                   { return c.definingLoader }
func (c *Class) SuperClass() *Class                                { return c.superClass }
func (c *Class) AccessFlags() uint16                               { return c.accessFlags }
func (c *Class) ConstantPool() *classfile.ConstantPool             { return c.cp }
func (c *Class) BootstrapMethods() *classfile.BootstrapMethodsAttr { return c.bootstrapMethods }
func (c *Class) InstCellCount() uint                               { return c.instCellCount }
func (c *Class) IsArray() bool                                     { return c.isArray }

// BindLoader sets the defining loader identity exactly once.
// Subsequent calls with the same identity are idempotent; a call with a
// different identity panics (invariant violation).
func (c *Class) BindLoader(id *LoaderIdentity) {
	if c.definingLoader == nil {
		c.definingLoader = id
		return
	}
	if c.definingLoader != id {
		panic("catty: Class.BindLoader: loader identity already bound to a different value")
	}
}

// BindLoaderRef sets the defining Loader interface reference exactly once.
// For VM primitives the ref stays nil; mirror creation falls back to BootstrapLoader.
// Subsequent calls with a different loader panic (invariant violation).
func (c *Class) BindLoaderRef(loader Loader) {
	if c.definingLoaderRef == nil {
		c.definingLoaderRef = loader
		return
	}
	if c.definingLoaderRef != loader {
		panic("catty: Class.BindLoaderRef: loader ref already bound to a different value")
	}
}

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

// ClassObject returns the canonical java.lang.Class mirror for this class,
// creating it lazily with double-checked locking on first access. All callers
// see the same Object identity (ADR-0029).
//
// The mirror's java/lang/Class is resolved through the class's defining loader.
// For VM primitives and void (which have no defining loader), BootstrapLoader
// is used as a fallback.
func (c *Class) ClassObject() *Object {
	if obj := c.classObject.Load(); obj != nil {
		return obj
	}
	c.classObjectMu.Lock()
	defer c.classObjectMu.Unlock()
	if obj := c.classObject.Load(); obj != nil {
		return obj
	}

	// Resolve java/lang/Class through the defining loader.
	loader := c.definingLoaderRef
	if loader == nil {
		loader = getBootstrapLoader()
	}
	if loader == nil {
		return nil // BootstrapLoader not yet set
	}
	result := loader.LoadClassResult("java/lang/Class")
	if !result.IsSuccess() {
		return nil
	}
	classClass := result.Class()

	obj := NewObject(classClass)
	obj.SetExtra(c)
	c.classObject.Store(obj)
	return obj
}

// classObjectWithFactory is test-only support for low-level rtda tests that do
// not construct a Loader. Production mirror creation must use ClassObject.
func (c *Class) classObjectWithFactory(factory func() *Object) *Object {
	if obj := c.classObject.Load(); obj != nil {
		return obj
	}
	c.classObjectMu.Lock()
	defer c.classObjectMu.Unlock()
	if obj := c.classObject.Load(); obj != nil {
		return obj
	}
	obj := factory()
	if obj == nil {
		return nil
	}
	obj.SetExtra(c)
	c.classObject.Store(obj)
	return obj
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
