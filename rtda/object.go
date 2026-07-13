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
type Object struct {
	class     *Class
	heapCells []HeapCell
	extra     any
}

func NewObject(class *Class) *Object {
	return &Object{class: class, heapCells: NewHeapCells(int(class.instCellCount))}
}

func (o *Object) Class() *Class { return o.class }

// Cells returns the heap cell storage backing instance fields / array elements.
// Callers use typed accessors (GetInt/SetInt/GetRef/SetRef/…) on individual cells.
func (o *Object) Cells() []HeapCell { return o.heapCells }

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
	// Each element of the outer array is itself an array of the component type.
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

// --- Typed array-element accessors (for AOT-emitted code and interpreter) ---
// One cell per element regardless of type. Long/double get a single 64-bit atomic cell.

func (o *Object) GetLongElement(i int) int64 {
	return o.heapCells[i].GetLong()
}

func (o *Object) SetLongElement(i int, v int64) {
	o.heapCells[i].SetLong(v)
}

func (o *Object) GetFloatElement(i int) float32 {
	return o.heapCells[i].GetFloat()
}

func (o *Object) SetFloatElement(i int, v float32) {
	o.heapCells[i].SetFloat(v)
}

func (o *Object) GetDoubleElement(i int) float64 {
	return o.heapCells[i].GetDouble()
}

func (o *Object) SetDoubleElement(i int, v float64) {
	o.heapCells[i].SetDouble(v)
}
