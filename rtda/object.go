package rtda

// Object is the runtime representation of a Java object OR array. Both are
// references on the operand stack and in fields/locals.
//
//   - Class instances store their per-instance state in heapCells (one HeapCell per
//     declared instance field per ADR-0030). The cell count per field is always 1;
//     long and double occupy a single 64-bit atomic cell — there is no 2-cell layout.
//   - Arrays store their elements in heapCells too (one HeapCell per element for
//     all component types). Element count == len(heapCells).
//
// Extra carries a native payload: java.lang.String stores its Go string here so
// System.out.println needs no interpreter-visible char array.
//
// heapCells is unexported per ADR-0030: no code outside rtda may access the
// mutable backing slice. All access goes through the typed per-cell accessors.
type Object struct {
	class     *Class
	heapCells []HeapCell
	extra     any
}

func NewObject(class *Class) *Object {
	return &Object{class: class, heapCells: NewHeapCells(int(class.instCellCount))}
}

// Class returns the object's runtime class.
func (o *Object) Class() *Class { return o.class }

// IsInstanceOf reports whether o can be treated as an instance of target,
// implementing the JVM instanceof and checkcast rules (JVMS §6.5.instanceof).
func (o *Object) IsInstanceOf(target *Class) bool {
	if o == nil {
		return false
	}
	return target.isAssignableFrom(o.class)
}

// SetExtra / Extra access the native payload slot.
func (o *Object) SetExtra(v any) { o.extra = v }
func (o *Object) Extra() any     { return o.extra }

// --- Array support (the class is flagged isArray) ---

func NewArray(class *Class, length int) *Object {
	return &Object{class: class, heapCells: NewHeapCells(length)}
}

// NewMultiArray recursively creates a multi-dimensional array. dims[0] is the
// outermost dimension; the class is the array type of the outermost level (e.g.
// "[[I" for new int[3][4]). For each level, sub-arrays are created whose
// component class is the next-inner array type (or a primitive/base type).
func NewMultiArray(class *Class, dims []int, loader Loader) *Object {
	arr := NewArray(class, dims[0])
	if len(dims) == 1 {
		return arr
	}
	componentClass := class.ComponentClass()
	if componentClass == nil {
		return arr // primitive innermost — already zero-initialized
	}
	subDims := dims[1:]
	for i := 0; i < dims[0]; i++ {
		arr.heapCells[i].SetRef(NewMultiArray(componentClass, subDims, loader))
	}
	return arr
}

func (o *Object) ArrayLength() int {
	return len(o.heapCells)
}

// cellCount returns the number of heap cells. Unexported — callers inside rtda
// may use it; external code uses the typed accessors or ArrayLength/InstCellCount.
func (o *Object) cellCount() int { return len(o.heapCells) }

// --- Typed per-cell accessors (ADR-0030) ---
// One cell per field/element regardless of Java type. Long and double are single
// 64-bit atomic cells. These are the ONLY way to access heap storage from outside
// the rtda package.

// GetIntCell returns the cell at index i as int32 (for Z, B, C, S, I and
// float-as-bits fields/elements).
func (o *Object) GetIntCell(i int) int32 { return o.heapCells[i].GetInt() }

// SetIntCell stores v into the cell at index i.
func (o *Object) SetIntCell(i int, v int32) { o.heapCells[i].SetInt(v) }

// GetLongCell returns the cell at index i as int64 (for J fields/elements).
func (o *Object) GetLongCell(i int) int64 { return o.heapCells[i].GetLong() }

// SetLongCell stores v into the cell at index i.
func (o *Object) SetLongCell(i int, v int64) { o.heapCells[i].SetLong(v) }

// GetFloatCell returns the cell at index i as float32 (for F fields/elements).
func (o *Object) GetFloatCell(i int) float32 { return o.heapCells[i].GetFloat() }

// SetFloatCell stores v into the cell at index i.
func (o *Object) SetFloatCell(i int, v float32) { o.heapCells[i].SetFloat(v) }

// GetDoubleCell returns the cell at index i as float64 (for D fields/elements).
func (o *Object) GetDoubleCell(i int) float64 { return o.heapCells[i].GetDouble() }

// SetDoubleCell stores v into the cell at index i.
func (o *Object) SetDoubleCell(i int, v float64) { o.heapCells[i].SetDouble(v) }

// GetRefCell returns the cell at index i as an object reference (for L/[
// fields/elements). Returns nil for null.
func (o *Object) GetRefCell(i int) *Object { return o.heapCells[i].GetRef() }

// SetRefCell stores r into the cell at index i.
func (o *Object) SetRefCell(i int, r *Object) { o.heapCells[i].SetRef(r) }

// --- Bulk copy (clone / System.arraycopy) ---

// copyCellsTyped copies length cells from src[srcOff:] to dst[dstOff:], using the
// given component kind to decide which atomic operations to use. It is the one
// implementation behind Object.clone and System.arraycopy.
func copyCellsTyped(dst, src []HeapCell, dstOff, srcOff, length, kind int) {
	switch kind {
	case kindLong:
		for i := 0; i < length; i++ {
			dst[dstOff+i].SetLong(src[srcOff+i].GetLong())
		}
	case kindDouble:
		for i := 0; i < length; i++ {
			dst[dstOff+i].SetDouble(src[srcOff+i].GetDouble())
		}
	case kindFloat:
		for i := 0; i < length; i++ {
			dst[dstOff+i].SetFloat(src[srcOff+i].GetFloat())
		}
	case kindBoolean, kindByte, kindChar, kindShort, kindInt:
		for i := 0; i < length; i++ {
			dst[dstOff+i].SetInt(src[srcOff+i].GetInt())
		}
	default:
		// reference arrays (including multi-dimensional)
		for i := 0; i < length; i++ {
			dst[dstOff+i].SetRef(src[srcOff+i].GetRef())
		}
	}
}

// CopyObjectCells copies length cells from src[srcOff:] to dst[dstOff:], using the
// component kind for typed dispatch. Preserves full 64-bit values for long/double.
// Handles overlapping regions within the same array with Java memmove semantics.
func CopyObjectCells(dst, src *Object, dstOff, srcOff, length, kind int) {
	if length <= 0 {
		return
	}
	// Java memmove semantics: if src and dst are the same array and the
	// regions overlap, copy backward to preserve the source data.
	if src == dst && srcOff < dstOff {
		// Overlap with source before destination: copy backward.
		copyCellsTypedReverse(dst.heapCells, src.heapCells, dstOff, srcOff, length, kind)
		return
	}
	copyCellsTyped(dst.heapCells, src.heapCells, dstOff, srcOff, length, kind)
}

// copyCellsTypedReverse copies length cells backward (from end to start).
func copyCellsTypedReverse(dst, src []HeapCell, dstOff, srcOff, length, kind int) {
	switch kind {
	case kindLong:
		for i := length - 1; i >= 0; i-- {
			dst[dstOff+i].SetLong(src[srcOff+i].GetLong())
		}
	case kindDouble:
		for i := length - 1; i >= 0; i-- {
			dst[dstOff+i].SetDouble(src[srcOff+i].GetDouble())
		}
	case kindFloat:
		for i := length - 1; i >= 0; i-- {
			dst[dstOff+i].SetFloat(src[srcOff+i].GetFloat())
		}
	case kindBoolean, kindByte, kindChar, kindShort, kindInt:
		for i := length - 1; i >= 0; i-- {
			dst[dstOff+i].SetInt(src[srcOff+i].GetInt())
		}
	default:
		for i := length - 1; i >= 0; i-- {
			dst[dstOff+i].SetRef(src[srcOff+i].GetRef())
		}
	}
}

// CloneCells returns a new independent slice with atomic cells copied from src.
// Used by Object.clone.
func CloneCells(src *Object) []HeapCell {
	dstHeap := make([]HeapCell, len(src.heapCells))
	CopyHeapCells(dstHeap, src.heapCells)
	return dstHeap
}

// CloneObject returns a shallow clone of obj: same class, independent copy of
// every cell. The monitor (if any) is NOT copied per ADR-0029. Extra is nil.
func CloneObject(obj *Object) *Object {
	clone := &Object{
		class:     obj.class,
		heapCells: CloneCells(obj),
	}
	return clone
}
