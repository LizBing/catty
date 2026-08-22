package vm

import (
	"fmt"
	"math"
	"strings"

	"catty/internal/classfile"
	"catty/internal/kernel"
)

// --- value/object helpers -----------------------------------------------------

func dotted(internal string) string { return strings.ReplaceAll(internal, "/", ".") }

func bool32v(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

// objClass returns the dynamic class of any heap value.
func objClass(v kernel.Value) *kernel.Class {
	switch o := v.(type) {
	case *kernel.Instance:
		return o.Class
	case *kernel.ArrayObj:
		return o.Class
	case *kernel.JString:
		return o.Class
	default:
		panic(fmt.Sprintf("vm: %T is not a heap object", v))
	}
}

func objHeader(v kernel.Value) *kernel.Header {
	switch o := v.(type) {
	case *kernel.Instance:
		return &o.Header
	case *kernel.ArrayObj:
		return &o.Header
	case *kernel.JString:
		return &o.Header
	default:
		panic(fmt.Sprintf("vm: %T is not a heap object", v))
	}
}

func objClassName(v kernel.Value) string { return objClass(v).Name }

// objFields returns instance field storage (getfield/putfield targets).
func objFields(v kernel.Value) []kernel.Value {
	in, ok := v.(*kernel.Instance)
	if !ok {
		panic(fmt.Sprintf("vm: field access on non-instance %T", v))
	}
	return in.Fields
}

// --- comparisons / conversions --------------------------------------------------

func cmpI64(a, b int64) int32 {
	switch {
	case a > b:
		return 1
	case a < b:
		return -1
	default:
		return 0
	}
}

// cmpFP implements fcmpl/fcmpg: NaN yields -1 for the _l form, +1 for _g.
func cmpFP(a, b float64, gForm bool) int32 {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	case a == b:
		return 0
	default: // NaN
		if gForm {
			return 1
		}
		return -1
	}
}

func condCmp1(op byte, v int32) bool {
	switch op {
	case opIfeq:
		return v == 0
	case opIfne:
		return v != 0
	case opIflt:
		return v < 0
	case opIfge:
		return v >= 0
	case opIfgt:
		return v > 0
	case opIfle:
		return v <= 0
	}
	return false
}

func condCmp2(op byte, a, b int32) bool {
	switch op {
	case opIfIcmpeq:
		return a == b
	case opIfIcmpne:
		return a != b
	case opIfIcmplt:
		return a < b
	case opIfIcmpge:
		return a >= b
	case opIfIcmpgt:
		return a > b
	case opIfIcmple:
		return a <= b
	}
	return false
}

// floatToI32 implements Java d2i/f2i: NaN→0, saturating clamp, truncate.
func floatToI32(f float64) int32 {
	if math.IsNaN(f) {
		return 0
	}
	if f >= 2147483647.0 {
		return 2147483647
	}
	if f <= -2147483648.0 {
		return -2147483648
	}
	return int32(f)
}

func floatToI64(f float64) int64 {
	if math.IsNaN(f) {
		return 0
	}
	if f >= 9223372036854775807.0 {
		return 9223372036854775807
	}
	if f <= -9223372036854775808.0 {
		return -9223372036854775808
	}
	return int64(f)
}

func atypeDesc(at byte) (string, bool) {
	switch at {
	case 4:
		return "Z", true
	case 5:
		return "C", true
	case 6:
		return "F", true
	case 7:
		return "D", true
	case 8:
		return "B", true
	case 9:
		return "S", true
	case 10:
		return "I", true
	case 11:
		return "J", true
	}
	return "", false
}

func returnsValue(desc string) bool {
	i := strings.IndexByte(desc, ')')
	if i < 0 || i+1 >= len(desc) {
		return false
	}
	return desc[i+1] != 'V'
}

// --- throwing -------------------------------------------------------------------

// throwNamed builds a Java throwable of the given bootstrap class.
// Panics only on engine bugs (missing bootstrap class).
func (t *Thread) throwNamed(className, msg string) *kernel.Thrown {
	obj, err := t.K.NewThrowable(className, msg)
	if err != nil {
		panic("vm: " + err.Error())
	}
	return &kernel.Thrown{Obj: obj}
}

func (t *Thread) npe(f *frame, why string) *kernel.Thrown {
	return t.throwNamed("java/lang/NullPointerException", why)
}

// arrayLoad: stack is …, arrayref, index.
func (t *Thread) arrayLoad(f *frame) (kernel.Value, *kernel.Thrown, error) {
	idx := f.popI()
	arr := f.popRef()
	if arr == nil {
		return nil, t.npe(f, "array load from null"), nil
	}
	a := arr.(*kernel.ArrayObj)
	if idx < 0 || int(idx) >= len(a.Elems) {
		return nil, t.throwNamed("java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems))), nil
	}
	return a.Elems[idx], nil, nil
}

// arrayStore: stack is …, arrayref, index, value.
func (t *Thread) arrayStore(f *frame, v kernel.Value) *kernel.Thrown {
	idx := f.popI()
	arr := f.popRef()
	if arr == nil {
		return t.npe(f, "array store into null")
	}
	a := arr.(*kernel.ArrayObj)
	if idx < 0 || int(idx) >= len(a.Elems) {
		return t.throwNamed("java/lang/ArrayIndexOutOfBoundsException",
			fmt.Sprintf("Index %d out of bounds for length %d", idx, len(a.Elems)))
	}
	a.Elems[idx] = v
	return nil
}

// --- constant pool resolution (cached per holder class) ---------------------------

type resKind uint8

const (
	resClass resKind = iota + 1
	resString
)

func (t *Thread) cached(holder *kernel.Class, kind resKind, idx uint16, compute func() (any, error)) (any, error) {
	key := uint32(kind)<<16 | uint32(idx)
	if v, ok := holder.CacheGet(key); ok {
		return v, nil
	}
	v, err := compute()
	if err != nil {
		return nil, err
	}
	holder.CacheSet(key, v)
	return v, nil
}

// resolveClassIdx resolves a CONSTANT_Class entry to a runtime class.
func (t *Thread) resolveClassIdx(holder *kernel.Class, idx uint16) (*kernel.Class, error) {
	v, err := t.cached(holder, resClass, idx, func() (any, error) {
		name, err := holder.CF.ClassName(idx)
		if err != nil {
			return nil, err
		}
		return t.K.ResolveClass(name)
	})
	if err != nil {
		return nil, err
	}
	return v.(*kernel.Class), nil
}

func (t *Thread) resolveClassNamed(name string) (*kernel.Class, error) {
	return t.K.ResolveClass(name)
}

// ldc resolves ldc/ldc_w operands (Integer/Float/String; Class unsupported).
func (t *Thread) ldc(holder *kernel.Class, idx uint16) (kernel.Value, error) {
	v, err := t.cached(holder, resString, idx, func() (any, error) {
		e, err := holder.CF.Entry(idx)
		if err != nil {
			return nil, err
		}
		switch e.Tag {
		case classfile.CInteger:
			return e.IntVal, nil
		case classfile.CFloat:
			return e.FloatVal, nil
		case classfile.CString:
			s, err := holder.CF.UTF8(e.Idx1)
			if err != nil {
				return nil, err
			}
			return t.K.InternGo(s), nil
		case classfile.CClass:
			return nil, fmt.Errorf("ldc of Class constant unsupported in M0")
		default:
			return nil, fmt.Errorf("ldc on %s", e.Tag)
		}
	})
	if err != nil {
		return nil, err
	}
	return v.(kernel.Value), nil
}

// buildMultiArray recursively creates multi-dimensional arrays.
func (t *Thread) buildMultiArray(desc string, dims []int32) (kernel.Value, error) {
	n := int(dims[0])
	arr, err := t.K.NewArray(desc[1:], n)
	if err != nil {
		return nil, err
	}
	if len(dims) > 1 && n > 0 {
		for i := range arr.Elems {
			sub, err := t.buildMultiArray(desc[1:], dims[1:])
			if err != nil {
				return nil, err
			}
			arr.Elems[i] = sub
		}
	}
	return arr, nil
}
