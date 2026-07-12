package rtda

// Slot is the atomic storage cell of both the operand stack and local variable
// arrays (JVMS §2.6). It is a tagged union: every slot holds exactly one of
//   - a category-1 numeric/returnAddress value in num (covers byte, char,
//     short, int, boolean, float-as-bits, and the returnAddress used by jsr/ret);
//   - a reference in ref (object or array; nil is Java null).
//
// Category-2 values (long, double) occupy TWO consecutive slots: the high
// 32 bits in the first, the low 32 bits in the second — see Frame.PushLong.
//
// This is the HotSpot "stack word" model. 16 bytes per slot (8-byte-aligned int
// + pointer) is larger than a native word; the documented performance-arc
// upgrade is to split into parallel []int32 / []*Object arrays to halve memory
// traffic on numeric code. That is intentionally deferred past MVP.
type Slot struct {
	num int32
	ref *Object
}

// Num/Ref read the slot's value; the interpreter uses these for field and array
// access. Same-package code (Frame) touches num/ref directly.
func (s Slot) Num() int32        { return s.num }
func (s Slot) Ref() *Object      { return s.ref }
func (s *Slot) SetNum(v int32)   { s.num = v }
func (s *Slot) SetRef(r *Object) { s.ref = r }

// RefSlot / IntSlot construct a Slot from a typed value. Cross-package code (the
// AOT runtime bridge, emitted programs) can't build Slot{...} directly since the
// fields are unexported; these are the constructors it uses to box call args.
func RefSlot(r *Object) Slot { return Slot{ref: r} }
func IntSlot(n int32) Slot   { return Slot{num: n} }
