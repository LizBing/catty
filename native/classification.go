// Package native holds native classification metadata for R2-A strict native
// resolution. Every registered native method has an entry here — each must be
// classified exactly once. Classification invariants are enforced by
// native/classification_test.go.
//
// Classification semantics:
//
//	Category                 | Meaning
//	==========================|==================================================
//	CategoryImplemented       | Full observable Java behaviour implemented in Go
//	CategorySemanticNoOp      | Java spec permits no required observable effect
//	CategoryCompatibilityAdapter | Different Go mechanism, equivalent Java behaviour
//	CategoryUnsupported       | Declared known; no callable implementation
package native

// Classification categorises a registered native method.
type Classification uint8

const (
	CategoryImplemented          Classification = iota // full behaviour
	CategorySemanticNoOp                               // spec permits emptiness
	CategoryCompatibilityAdapter                       // equivalent via different mechanism
	CategoryUnsupported                                // known, no implementation
)

func (c Classification) String() string {
	switch c {
	case CategoryImplemented:
		return "Implemented"
	case CategorySemanticNoOp:
		return "SemanticNoOp"
	case CategoryCompatibilityAdapter:
		return "CompatibilityAdapter"
	case CategoryUnsupported:
		return "Unsupported"
	default:
		return "Unknown"
	}
}

// ClassifiedEntry is the metadata for one registered native method.
type ClassifiedEntry struct {
	ClassName      string
	MethodName     string
	Descriptor     string
	Classification Classification
	Rationale      string // why this classification
}

// classifiedRegistry is the canonical list of every registered native method
// with its classification. Populated in the init() block below.
var classifiedRegistry []ClassifiedEntry

func init() {
	classifiedRegistry = []ClassifiedEntry{
		// --- Implemented: full observable behaviour ---
		{cs, "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", CategoryImplemented, "Element-wise array copy with bounds"},
		{cs, "currentTimeMillis", "()J", CategoryImplemented, "Go time.Now().UnixMilli() — equivalent to System.currentTimeMillis"},
		{cs, "nanoTime", "()J", CategoryImplemented, "Go time.Now().UnixNano() — equivalent to System.nanoTime"},
		{cs, "identityHashCode", "(Ljava/lang/Object;)I", CategoryImplemented, "Identity hash from Go object pointer"},
		{co, "hashCode", "()I", CategoryImplemented, "Identity hash from Go object pointer"},
		{co, "getClass", "()Ljava/lang/Class;", CategoryImplemented, "Returns rtda.Class-wrapped Class object"},
		{co, "clone", "()Ljava/lang/Object;", CategoryImplemented, "Shallow field copy via Go"},
		{cc, "getName", "()Ljava/lang/String;", CategoryImplemented, "Returns internal name as Java dotted string"},
		{cc, "getSimpleName", "()Ljava/lang/String;", CategoryImplemented, "Extracts simple name after last '/'"},
		{cc, "desiredAssertionStatus", "()Z", CategoryImplemented, "Returns false (assertions disabled)"},
		{cc, "isInterface", "()Z", CategoryImplemented, "Delegates to rtda.Class.IsInterface"},
		{cc, "isArray", "()Z", CategoryImplemented, "Delegates to rtda.Class.IsArray"},
		{cc, "getModifiers", "()I", CategoryImplemented, "Returns class access flags"},
		{cc, "isInstance", "(Ljava/lang/Object;)Z", CategoryImplemented, "Delegates to rtda.Object.IsInstanceOf"},
		{cc, "isAssignableFrom", "(Ljava/lang/Class;)Z", CategoryImplemented, "Walks super-class chain from arg"},
		{cc, "getSuperclass", "()Ljava/lang/Class;", CategoryImplemented, "Returns super-class as Class object"},
		{cc, "isHidden", "()Z", CategoryImplemented, "Returns false (no hidden classes in catty)"},
		{cc, "getPrimitiveClass", "(Ljava/lang/String;)Ljava/lang/Class;", CategoryImplemented, "Maps primitive names to array types"},
		{ct, "currentThread", "()Ljava/lang/Thread;", CategoryImplemented, "Returns a new Thread object"},
		{cstr, "intern", "()Ljava/lang/String;", CategoryImplemented, "Identity (returns this)"},
		{cr, "availableProcessors", "()I", CategoryImplemented, "Returns 1 (minimal)"},
		{cfl, "floatToRawIntBits", "(F)I", CategoryImplemented, "Go math.Float32bits"},
		{cfl, "intBitsToFloat", "(I)F", CategoryImplemented, "Go math.Float32frombits"},
		{cdb, "doubleToRawLongBits", "(D)J", CategoryImplemented, "Go math.Float64bits"},
		{cdb, "longBitsToDouble", "(J)D", CategoryImplemented, "Go math.Float64frombits"},

		// --- SemanticNoOp: Java spec permits no observable effect ---
		{co, "registerNatives", "()V", CategorySemanticNoOp, "JNI bootstrap hook — no effect per spec"},
		{cc, "registerNatives", "()V", CategorySemanticNoOp, "JNI bootstrap hook — no effect per spec"},
		{ct, "registerNatives", "()V", CategorySemanticNoOp, "JNI bootstrap hook — no effect per spec"},
		{cr, "gc", "()V", CategorySemanticNoOp, "GC hint — JVM may ignore; Go GC is concurrent"},

		// --- CompatibilityAdapter: equivalent behaviour via different mechanism ---
		{cs, "mapLibraryName", "(Ljava/lang/String;)Ljava/lang/String;", CategoryCompatibilityAdapter, "Identity mapping (no native libs)"},

		// --- Unsupported: declared, no callable implementation ---
		// Object threading methods — R2-C/D implement these properly
		{co, "notify", "()V", CategoryUnsupported, "Requires monitor implementation (R2-D)"},
		{co, "notifyAll", "()V", CategoryUnsupported, "Requires monitor implementation (R2-D)"},
		{co, "wait", "(J)V", CategoryUnsupported, "Requires monitor implementation (R2-D)"},

		// Thread.holdsLock — requires monitor implementation (R2-D)
		{ct, "holdsLock", "(Ljava/lang/Object;)Z", CategoryUnsupported, "Requires monitor implementation (R2-D)"},

		// Runtime memory queries — require Unsafe/memory model (R2-E+)
		{cr, "freeMemory", "()J", CategoryUnsupported, "Requires Unsafe memory model (R2-E+)"},
		{cr, "totalMemory", "()J", CategoryUnsupported, "Requires Unsafe memory model (R2-E+)"},
		{cr, "maxMemory", "()J", CategoryUnsupported, "Requires Unsafe memory model (R2-E+)"},

		// AccessController — security manager not implemented
		{cac, "doPrivileged", "(Ljava/security/PrivilegedAction;)Ljava/lang/Object;", CategoryUnsupported, "Security manager not implemented"},
		{cac, "doPrivileged", "(Ljava/security/PrivilegedExceptionAction;)Ljava/lang/Object;", CategoryUnsupported, "Security manager not implemented"},
	}
}

// Short class-name constants for the registry.
const (
	co   = "java/lang/Object"
	cs   = "java/lang/System"
	cc   = "java/lang/Class"
	ct   = "java/lang/Thread"
	cstr = "java/lang/String"
	cr   = "java/lang/Runtime"
	cfl  = "java/lang/Float"
	cdb  = "java/lang/Double"
	cac  = "java/security/AccessController"
)
