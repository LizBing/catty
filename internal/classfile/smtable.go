package classfile

import (
	"fmt"
	"os"
)

// Verification-type tags (JVMS §4.7.4).
const (
	VItemTop      = 0
	VItemInteger  = 1
	VItemFloat    = 2
	VItemDouble   = 3
	VItemLong     = 4
	VItemNull     = 5
	VItemUninitTh = 6
	VItemObject   = 7
	VItemUninit   = 8
)

// VItem is one verification_type_info entry. For VItemObject, CPoolIdx is
// the CONSTANT_Class index; for VItemUninit, Offset is the code offset of
// the `new` that created the value.
type VItem struct {
	Tag      byte
	CPoolIdx uint16
	Offset   uint16
}

// TagName renders the tag symbolically (diagnostics/tests).
func (v VItem) TagName() string {
	switch v.Tag {
	case VItemTop:
		return "top"
	case VItemInteger:
		return "int"
	case VItemFloat:
		return "float"
	case VItemDouble:
		return "double"
	case VItemLong:
		return "long"
	case VItemNull:
		return "null"
	case VItemUninitTh:
		return "uninitThis"
	case VItemObject:
		return "object"
	case VItemUninit:
		return "uninit"
	default:
		return fmt.Sprintf("vitem(%d)", v.Tag)
	}
}

// StackMapFrame is one full frame with absolute bytecode offset.
type StackMapFrame struct {
	Offset int32
	Locals []VItem
	Stack  []VItem

	// Kind is the raw frame_type byte; consumers derive semantics
	// (0-63 same, 64-127 same+1 stack, 247 ext-same, 248-250 chop,
	// 251 ext-same, 252-254 append, 255 full).
	Kind     byte
}

// ParseStackMapTable decodes the StackMapTable attribute payload
// (JVMS §4.7.4) into frames with absolute offsets. Long/double local items
// are expanded with an implicit top slot so consumers can index flat.
func ParseStackMapTable(cf *ClassFile, data []byte) ([]StackMapFrame, error) {
	return decodeFrames(cf, data)
}

// decodeFrames is the authoritative StackMapTable decoder.
var debugSM = os.Getenv("CATTY_SM_DEBUG") != ""

func decodeFrames(cf *ClassFile, data []byte) ([]StackMapFrame, error) {
	r := &reader{data: data}
	n := int(r.u2())
	if debugSM {
		fmt.Fprintf(os.Stderr, "[sm] len=%d count=%d\n", len(data), n)
	}
	frames := make([]StackMapFrame, 0, n)
	prevAbs := int32(-1)
	var prevLocals []VItem

	readItem := func() (VItem, error) {
		tag := r.u1()
		it := VItem{Tag: tag}
		switch tag {
		case VItemObject:
			it.CPoolIdx = r.u2()
		case VItemUninit:
			it.Offset = r.u2()
		case VItemTop, VItemInteger, VItemFloat, VItemDouble,
			VItemLong, VItemNull, VItemUninitTh:
		default:
			return VItem{}, fmt.Errorf("bad verification_type_info tag %d", tag)
		}
		return it, nil
	}

	for i := 0; i < n; i++ {
		ft := r.u1()
		if debugSM {
			fmt.Fprintf(os.Stderr, "[sm] ft=%d pos=%d\n", ft, r.pos)
		}
		var abs int32
		var locs []VItem
		var stk []VItem

		switch {
		case ft <= 63: // same_frame: offset_delta implied
			abs = prevAbs + int32(ft) + 1
		case ft <= 127: // same_locals_1_stack_item_frame
			abs = prevAbs + int32(ft-64) + 1
			it, err := readItem()
			if err != nil {
				return nil, err
			}
			stk = []VItem{it}
		case ft <= 246:
			return nil, fmt.Errorf("reserved frame_type %d", ft)
		case ft == 247: // same_frame_extended
			abs = prevAbs + int32(r.u2()) + 1
			locs = prevLocals
		case ft <= 250: // chop_frame
			abs = prevAbs + int32(r.u2()) + 1
			// chop: locals are the previous frame's minus k; the verifier
			// applies this over its running absolute state.
		case ft == 251: // same_frame_extended (u2 form)
			abs = prevAbs + int32(r.u2()) + 1
			locs = prevLocals
		case ft <= 254: // append_frame_extended: (ft-251) new locals
			abs = prevAbs + int32(r.u2()) + 1
			add := make([]VItem, 0, ft-251)
			for j := byte(0); j < ft-251; j++ {
				it, err := readItem()
				if err != nil {
					return nil, err
				}
				add = append(add, it)
			}
			locs = add
		default: // full_frame
			abs = prevAbs + int32(r.u2()) + 1
			nl := int(r.u2())
			for j := 0; j < nl; j++ {
				it, err := readItem()
				if err != nil {
					return nil, err
				}
				locs = append(locs, it)
			}
			ns := int(r.u2())
			for j := 0; j < ns; j++ {
				it, err := readItem()
				if err != nil {
					return nil, err
				}
				stk = append(stk, it)
			}
		}

		frames = append(frames, StackMapFrame{
			Offset: abs, Locals: locs, Stack: stk, Kind: ft,
		})
		prevAbs = abs
		prevLocals = locs
	}
	return frames, nil
}