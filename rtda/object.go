package rtda

// Object is the runtime representation of a Java object OR array. Both are
// references on the operand stack and in fields/locals.
//
//   - Class instances store their per-instance state in fields (one Slot per
//     declared field, laid out by Field.slotID). Go's GC traces these pointers
//     natively — catty writes no write barriers.
//   - Arrays store their elements in fields too (one Slot per element for
//     category-1 components; two for long[]/double[]), with the class marked
//     isArray. elementWidth() gives the per-element slot stride.
//
// Extra carries a native payload: java.lang.String stores its Go string here so
// System.out.println needs no interpreter-visible char array. Storing the value
// directly avoids a round-trip through a char[] for the common path.
type Object struct {
	class  *Class
	fields []Slot
	extra  any
}

func NewObject(class *Class) *Object {
	return &Object{class: class, fields: make([]Slot, class.instSlotCount)}
}

func (o *Object) Class() *Class { return o.class }

// Fields returns the slot storage backing instance fields / array elements, so
// the interpreter can read and write field slots by id and array elements by
// index without this package growing a typed accessor per type.
func (o *Object) Fields() []Slot { return o.fields }

// IsInstanceOf reports whether o can be treated as an instance of target,
// implementing the JVM instanceof and checkcast rules (JVMS §6.5.instanceof).
func (o *Object) IsInstanceOf(target *Class) bool {
	if o == nil {
		return false
	}
	return o.class.isAssignableFrom(target)
}

// SetExtra / Extra access the native payload slot.
func (o *Object) SetExtra(v any) { o.extra = v }
func (o *Object) Extra() any     { return o.extra }

// --- Array support (the class is flagged isArray) ---

func NewArray(class *Class, length int) *Object {
	width := 1
	if class.componentLongOrDouble() {
		width = 2
	}
	return &Object{class: class, fields: make([]Slot, length*width)}
}

func (o *Object) ArrayLength() int {
	width := 1
	if o.class.componentLongOrDouble() {
		width = 2
	}
	return len(o.fields) / width
}

// ArrayElementSlot returns the Slot at array index i (caller knows the type).
func (o *Object) ArrayElementSlot(i int) *Slot {
	width := 1
	if o.class.componentLongOrDouble() {
		width = 2
	}
	return &o.fields[i*width]
}
