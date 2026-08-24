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
	"strings"

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

// EnsureInit drives <clinit> for emitted-code static access, using a
// registry-backed recursion guard.
func EnsureInit(th kernel.OwnerKey, clsName string) {
	c, err := K.ResolveClass(clsName)
	if err != nil {
		panic("genrt EnsureInit: " + err.Error())
	}
	tr := kernel.KeyTracker{R: K.Threads, Key: thrKey(th)}
	_ = K.EnsureInitialized(tr, c)
}

func engineErr(err error) *kernel.Thrown {
	panic("genrt engine error: " + err.Error())
}

// --- resolution -----------------------------------------------------------------

type methodKey struct {
	cls, name, desc string
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
	if recv == nil {
		return nil, Throw(th, "java/lang/NullPointerException",
			fmt.Sprintf("invokevirtual %s.%s on null", cls, name))
	}
	dyn := ClassOf(recv)
	m := methodForDyn(dyn, name, desc)
	return invokeChecked(th, m, recv, args)
}

// CallSpecial resolves in the referenced class walking supers.
func CallSpecial(th kernel.OwnerKey, recv kernel.Value, cls, name, desc string, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	m := methodFor(cls, name, desc)
	return invokeChecked(th, m, recv, args)
}

func invokeChecked(th kernel.OwnerKey, m *kernel.Method, recv kernel.Value, args []kernel.Value) (kernel.Value, *kernel.Thrown) {
	key := thrKey(th)
	if err := K.Threads.FrameEnter(key); err != nil {
		return nil, Throw(th, "java/lang/StackOverflowError", "")
	}
	defer K.Threads.FrameExit(key)

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

// New allocates an instance of an internal class name.
func New(clsName string) *kernel.Instance {
	c, err := K.ResolveClass(clsName)
	if err != nil {
		panic("genrt new: " + err.Error())
	}
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
