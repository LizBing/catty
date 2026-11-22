package classfile

import (
	"fmt"

	"catty/rt"
)

// Verification types used by StackMapTable frames.
const (
	VTTop uint8 = iota
	VTInteger
	VTFloat
	VTDouble
	VTLong
	VTNull
	VTUninitializedThis
	VTObject
	VTUninitialized
)

// VType is one decoded verification_type_info.
type VType struct {
	Kind   uint8
	Class  string
	Offset int32
}

// Frame is one decoded StackMapTable frame with its absolute bytecode offset.
type Frame struct {
	Offset int32
	Locals []VType
	Stack  []VType
}

func parseStackMapTable(cp *constantPool, body []byte, m *Method) ([]Frame, error) {
	r := &reader{data: body}
	count, err := r.u2()
	if err != nil {
		return nil, fmt.Errorf("malformed stack map: %w", err)
	}
	initial, err := initialLocals(m)
	if err != nil {
		return nil, err
	}
	locals := append([]VType(nil), initial...)
	var frames []Frame
	prevOffset := int32(-1)
	for i := 0; i < int(count); i++ {
		frameType, err := r.u1()
		if err != nil {
			return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
		}
		offset := int32(0)
		stack := []VType(nil)
		nextLocals := locals
		switch {
		case frameType <= 63:
			offset = prevOffset + int32(frameType) + 1
			nextLocals = locals
		case frameType <= 127:
			offset = prevOffset + int32(frameType-64) + 1
			nextLocals = locals
			vt, err := readVType(r, cp)
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			stack = []VType{vt}
		case frameType == 247:
			delta, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			offset = prevOffset + int32(delta) + 1
			nextLocals = locals
			vt, err := readVType(r, cp)
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			stack = []VType{vt}
		case frameType >= 248 && frameType <= 250:
			delta, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			offset = prevOffset + int32(delta) + 1
			chop := int(251 - frameType)
			if chop > len(locals) {
				return nil, fmt.Errorf("malformed stack map frame %d: chop %d exceeds %d locals", i, chop, len(locals))
			}
			nextLocals = locals[:len(locals)-chop]
		case frameType == 251:
			delta, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			offset = prevOffset + int32(delta) + 1
			nextLocals = locals
		case frameType >= 252 && frameType <= 254:
			delta, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			offset = prevOffset + int32(delta) + 1
			appendCount := int(frameType - 251)
			nextLocals = append([]VType(nil), locals...)
			for j := 0; j < appendCount; j++ {
				vt, err := readVType(r, cp)
				if err != nil {
					return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
				}
				nextLocals = append(nextLocals, vt)
			}
		case frameType == 255:
			delta, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			offset = prevOffset + int32(delta) + 1
			nlocals, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			nextLocals = make([]VType, 0, int(nlocals))
			for j := 0; j < int(nlocals); j++ {
				vt, err := readVType(r, cp)
				if err != nil {
					return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
				}
				nextLocals = append(nextLocals, vt)
			}
			nstack, err := r.u2()
			if err != nil {
				return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
			}
			stack = make([]VType, 0, int(nstack))
			for j := 0; j < int(nstack); j++ {
				vt, err := readVType(r, cp)
				if err != nil {
					return nil, fmt.Errorf("malformed stack map frame %d: %w", i, err)
				}
				stack = append(stack, vt)
			}
		default:
			return nil, fmt.Errorf("malformed stack map: reserved frame type %d at entry %d", frameType, i)
		}
		if offset < 0 || int(offset) >= m.CodeLength {
			return nil, fmt.Errorf("malformed stack map frame %d: offset %d outside code length %d", i, offset, m.CodeLength)
		}
		if offset <= prevOffset {
			return nil, fmt.Errorf("malformed stack map frame %d: offset %d does not increase", i, offset)
		}
		if len(nextLocals) > m.MaxLocals {
			return nil, fmt.Errorf("malformed stack map frame %d: %d locals exceed max_locals %d", i, len(nextLocals), m.MaxLocals)
		}
		if len(stack) > m.MaxStack {
			return nil, fmt.Errorf("malformed stack map frame %d: stack %d exceeds max_stack %d", i, len(stack), m.MaxStack)
		}
		frames = append(frames, Frame{Offset: offset, Locals: append([]VType(nil), nextLocals...), Stack: append([]VType(nil), stack...)})
		locals = nextLocals
		prevOffset = offset
	}
	if r.off != len(body) {
		return nil, fmt.Errorf("malformed stack map: %d trailing bytes", len(body)-r.off)
	}
	return frames, nil
}

func initialLocals(m *Method) ([]VType, error) {
	mt, err := rt.ParseMethodType(m.Descriptor)
	if err != nil {
		return nil, err
	}
	var locals []VType
	if m.Flags&rt.AccStatic == 0 {
		if m.Name == "<init>" {
			locals = append(locals, VType{Kind: VTUninitializedThis})
		} else {
			locals = append(locals, VType{Kind: VTObject, Class: m.Owner})
		}
	}
	for _, p := range mt.Params {
		vt, err := vtypeFromType(p)
		if err != nil {
			return nil, err
		}
		locals = append(locals, vt)
	}
	return locals, nil
}

func vtypeFromType(t rt.Type) (VType, error) {
	switch t.Kind {
	case 'I':
		return VType{Kind: VTInteger}, nil
	case 'L':
		return VType{Kind: VTObject, Class: t.ClassName}, nil
	case 'A':
		return VType{Kind: VTObject, Class: t.Descriptor}, nil
	case 'V':
		return VType{}, fmt.Errorf("void is not a verification type")
	default:
		return VType{}, fmt.Errorf("unsupported type in stack map")
	}
}

func readVType(r *reader, cp *constantPool) (VType, error) {
	tag, err := r.u1()
	if err != nil {
		return VType{}, err
	}
	switch tag {
	case 0:
		return VType{Kind: VTTop}, nil
	case 1:
		return VType{Kind: VTInteger}, nil
	case 5:
		return VType{Kind: VTNull}, nil
	case 6:
		return VType{Kind: VTUninitializedThis}, nil
	case 7:
		idx, err := r.u2()
		if err != nil {
			return VType{}, err
		}
		name, err := cp.className(idx)
		if err != nil {
			return VType{}, err
		}
		return VType{Kind: VTObject, Class: name}, nil
	case 8:
		off, err := r.u2()
		if err != nil {
			return VType{}, err
		}
		return VType{Kind: VTUninitialized, Offset: int32(off)}, nil
	case 2, 3, 4:
		return VType{}, fmt.Errorf("category-two/floating verification type %d is outside the Gate 1 profile", tag)
	default:
		return VType{}, fmt.Errorf("unknown verification type tag %d", tag)
	}
}
