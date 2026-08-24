// Package genrt is the runtime bridge between emitter-generated Go code and
// the Catty kernel (emitter-abi.md §5). Generated functions carry the
// calling convention:
//
//	func Catty_<mangled>(thr kernel.OwnerKey, recv kernel.Value,
//	    args []kernel.Value) (kernel.Value, *kernel.Thrown)
//
// All semantics-risky operations (resolution, NPE checks, monitor enter/exit,
// bounds checks, division) are centralized here so generated code stays a
// thin translation of bytecode. Engine-level failures (missing classes,
// broken registries) panic loudly — they are bugs, not Java exceptions.
package genrt

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"catty/internal/kernel"
)

// K is the kernel installed by InstallKernel before any Java execution.
var K *kernel.Kernel

// resolved caches pool-driven method resolutions against the INSTALLED
// kernel. Entries are invalidated on every InstallKernel: they point into
// a specific kernel's class graph and must never leak across installs
// (in-process multi-kernel runs — tests, embedding — would otherwise
// execute stale bodies against mismatched class identities).
var resolved = map[methodKey]*kernel.Method{}

// InstallKernel wires the bridge to a kernel (done once by the entrypoint).
func InstallKernel(k *kernel.Kernel) {
	K = k
	resolved = make(map[methodKey]*kernel.Method)
	classCache = sync.Map{} // entries reference the previous kernel's graph
	for i := range icTable {
		icTable[i].Store(nil) // entries reference the previous kernel's graph
	}
}

func thrKey(th kernel.OwnerKey) uint64 {
	if th == nil {
		return 0
	}
	return th.OwnerKey()
}

func ownerOf(th kernel.OwnerKey) kernel.OwnerKey { return th }

// ClassOf returns the runtime class of any heap value.
func ClassOf(v kernel.Value) *kernel.Class {
	switch o := v.(type) {
	case *kernel.Instance:
		return o.Class
	case *kernel.ArrayObj:
		return o.Class
	case *kernel.JString:
		return o.Class
	default:
		panic(fmt.Sprintf("genrt: %T is not a heap object", v))
	}
}

func headerOf(v kernel.Value) *kernel.Header {
	switch o := v.(type) {
	case *kernel.Instance:
		return &o.Header
	case *kernel.ArrayObj:
		return &o.Header
	case *kernel.JString:
		return &o.Header
	default:
		panic(fmt.Sprintf("genrt: %T is not a heap object", v))
	}
}

func mustInstance(v kernel.Value, why string) *kernel.Instance {
	in, ok := v.(*kernel.Instance)
	if !ok {
		panic(fmt.Sprintf("genrt: %s on non-instance (%T)", why, v))
	}
	return in
}

// trackerFor returns the init tracker for a thread. Threads implementing
// kernel.InitTracker (the VM's Thread) must be used DIRECTLY: mixing this
// with a KeyTracker splits the in-progress bookkeeping across two stores,
// and a <clinit> that drops from emitted into interpreted code then
// deadlocks on its own class mutex (self-referential init, DEV-0002
// family).
func trackerFor(th kernel.OwnerKey) kernel.InitTracker {
	if it, ok := th.(kernel.InitTracker); ok {
		return it
	}
	return kernel.KeyTracker{R: K.Threads, Key: thrKey(th)}
}

// EnsureInit drives <clinit> for emitted-code static access, using a
// registry-backed recursion guard.
func EnsureInit(th kernel.OwnerKey, clsName string) {
	c := classFor(clsName)
	_ = K.EnsureInitialized(trackerFor(th), c)
}

func engineErr(err error) *kernel.Thrown {
	panic("genrt engine error: " + err.Error())
}

// --- resolution -----------------------------------------------------------------

type methodKey struct {
	cls, name, desc string
}

// icEntry is one monomorphic inline-cache slot payload. Immutable after
// construction; swapped in via atomic pointer so concurrent callers never
// observe a torn {dyn, m} pair.
type icEntry struct {
	cls, name, desc string // call-site identity (compile-time constants)
	dyn             *kernel.Class
	m               *kernel.Method
}

// icTable is the direct-mapped virtual-dispatch cache (P-0009 T1).
// Generated code passes the same constant strings per call site, so hit
// verification is mostly pointer equality. Cleared on InstallKernel:
// entries point into a specific kernel's class graph.
var icTable [1024]atomic.Pointer[icEntry]

func icHash(cls, name, desc string) uint32 {
	h := uint32(2166136261)
	for _, s := range [3]string{cls, name, desc} {
		for i := 0; i < len(s); i++ {
			h ^= uint32(s[i])
			h *= 16777619
		}
	}
	return h % uint32(len(icTable))
}

// ICSlot is the public form of icHash for the emitter: genemit bakes its
// result as a literal into each virtual/interface call site
// (TestEmitterICSlotAgreement pins both implementations together).
func ICSlot(cls, name, desc string) uint32 { return icHash(cls, name, desc) }

// methodForDynCached resolves a virtual dispatch through the inline
// cache. Monomorphic sites — the overwhelming majority in practice —
// pay one atomic load plus pointer compares instead of a superclass-chain
// walk with per-call memberKey allocation.
func methodForDynCached(dyn *kernel.Class, cls, name, desc string) *kernel.Method {
	return methodForDynCachedSlot(icHash(cls, name, desc), dyn, cls, name, desc)
}

func methodForDynCachedSlot(slot uint32, dyn *kernel.Class, cls, name, desc string) *kernel.Method {
	entry := &icTable[slot]
	if e := entry.Load(); e != nil && e.dyn == dyn &&
		e.cls == cls && e.name == name && e.desc == desc {
		return e.m
	}
	m := methodForDyn(dyn, name, desc)
	entry.Store(&icEntry{cls: cls, name: name, desc: desc, dyn: dyn, m: m})
	return m
}

func methodForDyn(dyn *kernel.Class, name, desc string) *kernel.Method {
	m, err := K.ResolveMethod(dyn, name, desc)
	if err != nil {
		panic("genrt resolve virtual: " + err.Error())
	}
	return m
}

func methodFor(cls, name, desc string) *kernel.Method {
	key := methodKey{cls, name, desc}
	if m, ok := resolved[key]; ok {
		return m
	}
	c, err := K.ResolveClass(cls)
	if err != nil {
		panic("genrt resolve: " + err.Error())
	}
	m, err := K.ResolveMethod(c, name, desc)
	if err != nil {
		panic("genrt resolve: " + err.Error())
	}
	resolved[key] = m
	return m
}

// --- calls ------------------------------------------------------------------------

// CallStatic invokes a static method (initializes the declaring class via
// the kernel's init path inside InvokeAs).
func CallStatic(th kernel.OwnerKey, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	m := methodFor(cls, name, desc)
	return invokeChecked(th, m, nil, args)
}

// CallVirtual resolves against recv's dynamic class. Null receivers raise NPE.
func CallVirtual(th kernel.OwnerKey, recv kernel.Value, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	return CallVirtualIC(icHash(cls, name, desc), th, recv, cls, name, desc, args)
}

// CallVirtualIC is CallVirtual with the call-site's inline-cache slot
// precomputed at emission time (no per-call string hashing).
func CallVirtualIC(slot uint32, th kernel.OwnerKey, recv kernel.Value, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	if recv == nil {
		return nil, Throw(th, "java/lang/NullPointerException",
			fmt.Sprintf("invokevirtual %s.%s on null", cls, name))
	}
	dyn := ClassOf(recv)
	m := methodForDynCachedSlot(slot, dyn, cls, name, desc)
	return invokeChecked(th, m, recv, args)
}

// CallSpecial resolves in the referenced class walking supers.
func CallSpecial(th kernel.OwnerKey, recv kernel.Value, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	m := methodFor(cls, name, desc)
	return invokeChecked(th, m, recv, args)
}

func invokeChecked(th kernel.OwnerKey, m *kernel.Method, recv kernel.Value, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	// Frame metering: prefer the thread's own lock-free counter (shared
	// budget with interpreted frames); fall back to the registry for
	// owner keys without a meter.
	if fm, ok := th.(kernel.FrameMeter); ok {
		if err := fm.EnterFrame(); err != nil {
			return nil, Throw(th, "java/lang/StackOverflowError", "")
		}
		defer fm.ExitFrame()
	} else {
		key := thrKey(th)
		if err := K.Threads.FrameEnter(key); err != nil {
			return nil, Throw(th, "java/lang/StackOverflowError", "")
		}
		defer K.Threads.FrameExit(key)
	}

	v, err := K.InvokeAs(ownerOf(th), m, recv, args)
	if err == nil {
		return v, nil
	}
	if t, ok := err.(*kernel.Thrown); ok {
		return nil, t
	}
	panic(fmt.Sprintf("genrt engine error in %s.%s%s: %v",
		m.Holder.Name, m.Name, m.Desc, err))
}

// --- fields --------------------------------------------------------------------------

func GetStatic(th kernel.OwnerKey, cls, name, desc string) kernel.Value {
	EnsureInit(th, cls)
	c, _ := K.ResolveClass(cls)
	f := fieldFor(c, name, desc)
	return c.Statics[f.StaticSlot]
}

func SetStatic(th kernel.OwnerKey, cls, name, desc string, v kernel.Value) {
	EnsureInit(th, cls)
	c, _ := K.ResolveClass(cls)
	f := fieldFor(c, name, desc)
	c.Statics[f.StaticSlot] = v
}

func GetField(recv kernel.Value, cls, name, desc string) kernel.Value {
	in := mustInstance(recv, "getfield")
	f := fieldFor(in.Class, name, desc)
	return in.Fields[f.Slot]
}

func SetField(recv kernel.Value, cls, name, desc string, v kernel.Value) {
	in := mustInstance(recv, "putfield")
	f := fieldFor(in.Class, name, desc)
	in.Fields[f.Slot] = v
}

func fieldFor(cls *kernel.Class, name, desc string) *kernel.Field {
	f, err := K.ResolveField(cls, name, desc)
	if err != nil {
		panic("genrt field: " + err.Error())
	}
	return f
}

// --- allocation / typing ---------------------------------------------------------------

// classCache memoizes name→*Class for the installed kernel (cleared by
// InstallKernel). Allocation paths resolve classes by name on every `new`.
var classCache sync.Map // string -> *kernel.Class

func classFor(clsName string) *kernel.Class {
	if v, ok := classCache.Load(clsName); ok {
		return v.(*kernel.Class)
	}
	c, err := K.ResolveClass(clsName)
	if err != nil {
		panic("genrt resolve: " + err.Error())
	}
	classCache.Store(clsName, c)
	return c
}

// New allocates an instance of an internal class, mirroring the
// interpreter's `new` semantics exactly: java/lang/String gets its magic
// JString representation, and the class is driven to initialized state
// before the constructor runs (JVMS §5.5).
func New(th kernel.OwnerKey, clsName string) kernel.Value {
	c := classFor(clsName)
	if c.Name == "java/lang/String" {
		return K.MakeJString(nil)
	}
	_ = K.EnsureInitialized(trackerFor(th), c)
	in, ierr := K.NewInstance(c)
	if ierr != nil {
		panic("genrt new: " + ierr.Error())
	}
	return in
}

// Str interns a string literal.
func Str(s string) *kernel.JString { return K.InternGo(s) }

func InstanceOf(v kernel.Value, clsName string) bool {
	c, ok := K.ClassByName(clsName)
	if !ok {
		return false
	}
	return K.IsInstance(v, c)
}

// CheckCast returns v, or the exception object for a failed cast.
func CheckCast(th kernel.OwnerKey, v kernel.Value, clsName string) kernel.Value {
	c, ok := K.ClassByName(clsName)
	if !ok || c == nil {
		panic("genrt checkcast target missing: " + clsName)
	}
	if v != nil && !K.IsInstance(v, c) {
		t := Throw(th, "java/lang/ClassCastException",
			fmt.Sprintf("class %s cannot be cast to class %s", TypeName(v), dotted(clsName)))
		return t.Obj
	}
	return v
}

func TypeName(v kernel.Value) string {
	switch o := v.(type) {
	case *kernel.Instance:
		return dotted(o.Class.Name)
	case *kernel.ArrayObj:
		return o.Class.Name
	case *kernel.JString:
		return "java.lang.String"
	}
	return "?"
}

func dotted(internal string) string { return strings.ReplaceAll(internal, "/", ".") }

// --- arithmetic with Java semantics -------------------------------------------------------

func IDiv(th kernel.OwnerKey, a, b int32) (int32, *kernel.Thrown) {
	if b == 0 {
		return 0, Throw(th, "java/lang/ArithmeticException", "/ by zero")
	}
	return a / b, nil
}

func IRem(th kernel.OwnerKey, a, b int32) (int32, *kernel.Thrown) {
	if b == 0 {
		return 0, Throw(th, "java/lang/ArithmeticException", "/ by zero")
	}
	return a % b, nil
}

// --- arrays ----------------------------------------------------------------------------------

func ALoad(arr kernel.Value, idx int32) kernel.Value {
	a := arr.(*kernel.ArrayObj)
	if idx < 0 || int(idx) >= len(a.Elems) {
		t := Throw(nil, "java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems)))
		return t.Obj
	}
	return a.Elems[idx]
}

// AStore stores and returns a non-nil exception object when out of bounds.
func AStore(arr kernel.Value, idx int32, v kernel.Value) kernel.Value {
	a := arr.(*kernel.ArrayObj)
	if idx < 0 || int(idx) >= len(a.Elems) {
		t := Throw(nil, "java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems)))
		return t.Obj
	}
	a.Elems[idx] = v
	return nil
}

func ArrayLength(arr kernel.Value) int32 { return int32(len(arr.(*kernel.ArrayObj).Elems)) }

// ArrayLengthChecked implements arraylength with NPE semantics.
func ArrayLengthChecked(th kernel.OwnerKey, arr kernel.Value) (int32, *kernel.Thrown) {
	a, ok := arr.(*kernel.ArrayObj)
	if !ok {
		return 0, Throw(th, "java/lang/NullPointerException", "array length")
	}
	return int32(len(a.Elems)), nil
}

// --- monitors -----------------------------------------------------------------------------------

func MonitorEnter(th kernel.OwnerKey, obj kernel.Value) *kernel.Thrown {
	h := headerOf(obj)
	h.Monitor().Enter(thrKey(th))
	return nil
}

func MonitorExit(th kernel.OwnerKey, obj kernel.Value) *kernel.Thrown {
	h := headerOf(obj)
	if err := h.Monitor().Exit(thrKey(th)); err != nil {
		return Throw(th, "java/lang/IllegalMonitorStateException",
			"monitorexit without ownership")
	}
	return nil
}

// SetLine records the source line the current (top) Java frame is
// executing. Emitted code calls this at line-segment entries and labels;
// it is what gives stack traces exact leaf/call-site lines (U3).
func SetLine(th kernel.OwnerKey, line int32) {
	if ft, ok := th.(kernel.FrameTracker); ok {
		ft.SetTopJavaLine(line)
	}
}

// --- string helpers for folded StringBuilder chains (P-0009 U1) -----------

// MakeStr allocates a fresh (non-interned) JString — toString semantics.
func MakeStr(s string) kernel.Value { return K.MakeJStringFromGo(s) }

// StrOf extracts Go text from a value known to be a String.
func StrOf(v kernel.Value) string { return v.(*kernel.JString).Go() }

// ItoA/JtoA render integral appends; CtoA renders a UTF-16 unit; ZtoA a bool.
func ItoA(v int32) string { return strconv.Itoa(int(v)) }
func JtoA(v int64) string { return strconv.FormatInt(v, 10) }
func CtoA(u int32) string { return string(rune(u)) }
func ZtoA(v int32) string {
	if v != 0 {
		return "true"
	}
	return "false"
}

// --- throw helper ----------------------------------------------------------------------------------

// Throw builds a bootstrap exception instance wrapped as Thrown. The
// current Java call stack (maintained by kernel.InvokeAs) is backfilled
// onto the throwable, mirroring fillInStackTrace (DEBT-0019).
func Throw(th kernel.OwnerKey, cls, msg string) *kernel.Thrown {
	obj, err := K.NewThrowable(cls, msg)
	if err != nil {
		panic("genrt: " + err.Error())
	}
	kernel.AttachTraceTo(th, obj)
	return &kernel.Thrown{Obj: obj}
}

// RefEq implements if_acmp_eq/ne reference identity.
func RefEq(a, b kernel.Value) bool { return a == b }

// BoolValue boxes an instanceof-style boolean result.
func BoolValue(b bool) kernel.Value {
	if b {
		return int32(1)
	}
	return int32(0)
}

// GetFieldChecked performs getfield with null-receiver NPE semantics.
func GetFieldChecked(th kernel.OwnerKey, recv kernel.Value, name, desc string) (kernel.Value, *kernel.Thrown) {
	if recv == nil {
		return nil, Throw(th, "java/lang/NullPointerException",
			fmt.Sprintf("getfield %s on null", name))
	}
	in := recv.(*kernel.Instance)
	f, err := K.ResolveField(in.Class, name, desc)
	if err != nil {
		panic("genrt field: " + err.Error())
	}
	return in.Fields[f.Slot], nil
}

// SetFieldChecked performs putfield with null-receiver NPE semantics.
func SetFieldChecked(th kernel.OwnerKey, recv kernel.Value, name, desc string, v kernel.Value) *kernel.Thrown {
	if recv == nil {
		return Throw(th, "java/lang/NullPointerException",
			fmt.Sprintf("putfield %s on null", name))
	}
	in := recv.(*kernel.Instance)
	f, err := K.ResolveField(in.Class, name, desc)
	if err != nil {
		panic("genrt field: " + err.Error())
	}
	in.Fields[f.Slot] = v
	return nil
}

// NewPrimitiveArray allocates a primitive array by component descriptor.
func NewPrimitiveArray(th kernel.OwnerKey, compDesc string, size int32) (kernel.Value, *kernel.Thrown) {
	if size < 0 {
		return nil, Throw(th, "java/lang/NegativeArraySizeException", fmt.Sprintf("%d", size))
	}
	a, aerr := K.NewArray(compDesc, int(size))
	if aerr != nil {
		panic("genrt NewPrimitiveArray: " + aerr.Error())
	}
	return a, nil
}

// ALoadChecked reads an array element with NPE/bounds semantics.
func ALoadChecked(th kernel.OwnerKey, arr kernel.Value, idx int32) (kernel.Value, *kernel.Thrown) {
	a, ok := arr.(*kernel.ArrayObj)
	if !ok {
		return nil, Throw(th, "java/lang/NullPointerException", "array load")
	}
	if idx < 0 || int(idx) >= len(a.Elems) {
		return nil, Throw(th, "java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems)))
	}
	return a.Elems[idx], nil
}

// AStoreChecked writes an array element with NPE/bounds semantics.
func AStoreChecked(th kernel.OwnerKey, arr kernel.Value, idx int32, v kernel.Value) *kernel.Thrown {
	if os.Getenv("CATTY_ASD") != "" {
		fmt.Fprintf(os.Stderr, "[asd] arr_nil=%v idx=%d\n", arr == nil, idx)
	}
	a, ok := arr.(*kernel.ArrayObj)
	if !ok {
		return Throw(th, "java/lang/NullPointerException", "array store into null")
	}
	if idx < 0 || int(idx) >= len(a.Elems) {
		return Throw(th, "java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems)))
	}
	a.Elems[idx] = v
	return nil
}

// CallInterface resolves like invokevirtual (dynamic dispatch).
func CallInterface(th kernel.OwnerKey, recv kernel.Value, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	return CallVirtual(th, recv, cls, name, desc, args)
}

// LDiv implements JLS 15.17.1 long division incl. div-by-zero throw.
func LDiv(th kernel.OwnerKey, a, b int64) (int64, *kernel.Thrown) {
	if b == 0 {
		return 0, Throw(th, "java/lang/ArithmeticException", "/ by zero")
	}
	return a / b, nil
}

// LRem implements JLS 15.17.3 long remainder.
func LRem(th kernel.OwnerKey, a, b int64) (int64, *kernel.Thrown) {
	if b == 0 {
		return 0, Throw(th, "java/lang/ArithmeticException", "/ by zero")
	}
	return a % b, nil
}

// --- comparisons (JVMS §6.5 lcmp/fcmp/dcmp NaN semantics) -----------------

func LCmp(a, b int64) int32 {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// fCmp32 shares the fcmpl/fcmpg shape: v<c → -1, v>c → +1, equal → 0,
// NaN → nanResult (JLS 15.20.1: direction depends on the opcode).
func fCmp32(v, c float32, nanResult int32) int32 {
	if v < c {
		return -1
	}
	if v > c {
		return 1
	}
	if v == c {
		return 0
	}
	return nanResult // NaN involved
}

func fCmp64(v, c float64, nanResult int32) int32 {
	if v < c {
		return -1
	}
	if v > c {
		return 1
	}
	if v == c {
		return 0
	}
	return nanResult
}

func FCmpl(a, b float32) int32 { return fCmp32(a, b, -1) }
func FCmpg(a, b float32) int32 { return fCmp32(a, b, 1) }
func DCmpl(a, b float64) int32 { return fCmp64(a, b, -1) }
func DCmpg(a, b float64) int32 { return fCmp64(a, b, 1) }

// NewRefArray allocates a reference-component array with nil elements.
func NewRefArray(th kernel.OwnerKey, compClass string, n int32) (kernel.Value, *kernel.Thrown) {
	if n < 0 {
		return nil, Throw(th, "java/lang/NegativeArraySizeException", fmt.Sprintf("%d", n))
	}
	arrCls := K.ArrayClassOf("L" + compClass + ";")
	if arrCls == nil {
		panic("[nra-dbg] no array class for " + compClass)
	}
	a := &kernel.ArrayObj{CompDesc: "L" + compClass + ";", Elems: make([]kernel.Value, n)}
	a.Class = arrCls
	return a, nil
}
