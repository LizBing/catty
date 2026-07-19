package classfile

// StackMapTableAttribute (JVMS §4.7.4) carries the verification types of locals
// and operand stack at selected bytecode offsets (every branch target for Java
// 6+). catty uses it to pin merge points during type tracking — no type lattice
// or merge logic is needed, because the table gives the exact merged frame.
//
// Frames are delta-encoded off the previous frame, and the first frame's
// "previous" is the method's implicit initial frame (its argument locals). That
// initial locals come from the method descriptor, which this attribute doesn't
// have — so the parser stores raw frames and Reconstruct(initialLocals) produces
// the full frames once the caller supplies the seed.
type StackMapTableAttribute struct {
	raw []rawFrame
}

// StackMapFrame is one fully-reconstructed frame at a bytecode offset.
type StackMapFrame struct {
	Offset int         // bytecode pc this frame applies at
	Locals []VerifType // locals types (in slot order)
	Stack  []VerifType // operand-stack types (bottom-to-top)
}

// VerifType is one verification_type_info (JVMS §4.7.4). Tag is an Item*
// constant; ClassIndex is set for ItemObject (a CONSTANT_Class index), Offset
// for ItemUninitialized.
type VerifType struct {
	Tag        uint8
	ClassIndex uint16
	Offset     uint16
}

// Verification type item tags (JVMS §4.7.4, Table 4.7.4-B).
const (
	ItemTop               = 0
	ItemInteger           = 1
	ItemFloat             = 2
	ItemDouble            = 3
	ItemLong              = 4
	ItemNull              = 5
	ItemUninitializedThis = 6
	ItemObject            = 7
	ItemUninitialized     = 8
)

// rawFrame is one parsed entry before delta-reconstruction. `locals` holds the
// appended types (APPEND) or the full locals (FULL); `stack` holds the single
// stack type (SAME_LOCALS_1) or the full stack (FULL).
type rawFrame struct {
	offsetDelta int
	kind        frameKind
	chop        int // CHOP: number of locals removed
	locals      []VerifType
	stack       []VerifType
}

type frameKind uint8

const (
	kindSame frameKind = iota
	kindSameLocals1
	kindChop
	kindAppend
	kindFull
)

func (a *StackMapTableAttribute) readInfo(_ *ClassReader) {} // parsed by readStackMapTable

// Reconstruct applies the delta-encoding to produce the full frame at each
// offset, given the method's initial locals (the seed = argument types from the
// method descriptor). Long/double remain single VerifType entries here (two
// slots); the lowering expands them.
func (a *StackMapTableAttribute) Reconstruct(initialLocals []VerifType) []StackMapFrame {
	frames := make([]StackMapFrame, len(a.raw))
	prevLocals := initialLocals
	prevOffset := -1
	for i, rf := range a.raw {
		var locals, stack []VerifType
		switch rf.kind {
		case kindSame:
			locals = prevLocals
		case kindSameLocals1:
			locals = prevLocals
			stack = rf.stack
		case kindChop:
			locals = prevLocals[:len(prevLocals)-rf.chop]
		case kindAppend:
			locals = make([]VerifType, len(prevLocals)+len(rf.locals))
			copy(locals, prevLocals)
			copy(locals[len(prevLocals):], rf.locals)
		case kindFull:
			locals = rf.locals
			stack = rf.stack
		}
		offset := prevOffset + rf.offsetDelta + 1
		prevOffset = offset
		prevLocals = append([]VerifType(nil), locals...) // copy so later frames can't mutate this one
		frames[i] = StackMapFrame{Offset: offset, Locals: prevLocals, Stack: stack}
	}
	return frames
}

// readStackMapTable parses number_of_entries + the delta-encoded frames into raw
// form (no reconstruction — the initial-locals seed is applied later).
func readStackMapTable(info []byte) *StackMapTableAttribute {
	r := NewClassReader(info)
	n := int(r.ReadUint16())
	a := &StackMapTableAttribute{raw: make([]rawFrame, 0, n)}
	for range n {
		ft := r.ReadUint8()
		rf := rawFrame{}
		switch {
		case ft <= 63: // same_frame
			rf.kind = kindSame
			rf.offsetDelta = int(ft)
		case ft <= 127: // same_locals_1_stack_item_frame
			rf.kind = kindSameLocals1
			rf.offsetDelta = int(ft - 64)
			rf.stack = []VerifType{readVerifType(r)}
		case ft == 247: // same_locals_1_stack_item_frame_extended
			rf.kind = kindSameLocals1
			rf.offsetDelta = int(r.ReadUint16())
			rf.stack = []VerifType{readVerifType(r)}
		case ft >= 248 && ft <= 250: // chop_frame
			rf.kind = kindChop
			rf.offsetDelta = int(r.ReadUint16())
			rf.chop = 251 - int(ft)
		case ft >= 251 && ft <= 254: // append_frame
			rf.kind = kindAppend
			rf.offsetDelta = int(r.ReadUint16())
			rf.locals = readVerifTypes(r, int(ft)-251)
		case ft == 255: // full_frame
			rf.kind = kindFull
			rf.offsetDelta = int(r.ReadUint16())
			rf.locals = readVerifTypes(r, int(r.ReadUint16()))
			rf.stack = readVerifTypes(r, int(r.ReadUint16()))
		default:
			panicf("StackMapTable", "reserved stack_map_frame type %d", ft)
		}
		a.raw = append(a.raw, rf)
	}
	if r.Len() > 0 {
		panicf("StackMapTable", "trailing bytes after StackMapTable attribute body (%d remaining)", r.Len())
	}
	return a
}

func readVerifTypes(r *ClassReader, n int) []VerifType {
	ts := make([]VerifType, n)
	for i := range ts {
		ts[i] = readVerifType(r)
	}
	return ts
}

func readVerifType(r *ClassReader) VerifType {
	t := VerifType{Tag: r.ReadUint8()}
	switch t.Tag {
	case ItemObject:
		t.ClassIndex = r.ReadUint16()
	case ItemUninitialized:
		t.Offset = r.ReadUint16()
	}
	return t
}
