package emitter

import (
	"catty/internal/classfile"
)

// qItem is a worklist entry for depth simulation.
type qItem struct{ pc, depth int }

// computeCanonicalDepths simulates operand-stack depth for every reachable
// pc using SM frame canonical values at merge points and per-opcode stack
// effects elsewhere.
func (e *methodEmitter) computeCanonicalDepths(reach map[int]bool) map[int]int {
	depths := make(map[int]int)
	code := e.code

	entryDepth := 0
	argDescs, _, err := splitMethodDesc(e.m.Desc)
	if err == nil {
		for _, a := range argDescs {
			entryDepth += descSlots(a)
		}
		if e.m.AccessFlags&classfile.AccStatic == 0 {
			entryDepth++
		}
	}

	// Build SM frame lookup.
	smDepth := make(map[int]int)
	for i := range e.m.Code.StackMaps {
		fr := &e.m.Code.StackMaps[i]
		n := 0
		for _, it := range fr.Stack {
			n++
			if it.Tag == classfile.VItemLong || it.Tag == classfile.VItemDouble {
				n++
			}
		}
		smDepth[int(fr.Offset)] = n
	}

	handlerSet := make(map[int]bool)
	for pc := range e.handlerAt {
		handlerSet[pc] = true
	}

	d := entryDepth
	for pc := 0; pc < len(code); pc++ {
		if !reach[pc] {
			continue
		}
		op := code[pc]

		// Exception handler entries: stack cleared + exception pushed = 1.
		if handlerSet[pc] {
			d = 1
		} else if sd, ok := smDepth[pc]; ok {
			d = sd
		}
		depths[pc] = d

		effect := e.instrEffect(pc, op)
		d += effect
		if d < 0 {
			d = 0
		}

		// Skip non-linear successors.
		switch op {
		case 0xa7, 0xc8: // goto/goto_w: jump target will re-sync from SM
			pc += instrSizeSafe(code, pc) - 1
		case 0xaa, 0xab: // switches
			pc += instrSizeSafe(code, pc) - 1
		default:
			// fallthrough handled by loop increment
		}
	}
	return depths
}

// netStackEffect returns how many slots an opcode adds to the operand stack.
func netStackEffect(op byte) int {
	switch op {
	case 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x10, 0x11, 0x12, 0x13, 0x15, 0x17, 0x19,
		0x1a, 0x1b, 0x1c, 0x1d, 0x2a, 0x2b, 0x2c, 0x2d:
		return 1
	case 0x09, 0x0a, 0x0e, 0x0f, 0x14,
		0x16, 0x18, 0x1e, 0x1f, 0x20, 0x21, 0x26, 0x27, 0x28, 0x29:
		return 2
	case 0x36, 0x38, 0x3a,
		0x3b, 0x3c, 0x3d, 0x3e,
		0x43, 0x44, 0x45, 0x46,
		0x4b, 0x4c, 0x4d, 0x4e,
		0x57, // pop
		0x2e, 0x30, 0x32, 0x33, 0x34, 0x35: // cat1 *aload
		return -1
	case 0x37, 0x3f, 0x40, 0x41, 0x42, // lstore family (cat2)
		0x39, 0x47, 0x48, 0x49, 0x4a: // dstore family (cat2)
		return -2
	case 0x2f, 0x31: // laload/daload (cat2 load)
		return 2
	case 0x58:
		return -2
	case 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e:
		return 1
	case 0x5f:
		return 0
	case 0x60, 0x62, 0x64, 0x65, 0x66, 0x67,
		0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
		0x70, 0x71, 0x72, 0x73,
		0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
		0x80, 0x81, 0x82, 0x83, 0x94, 0x95, 0x96, 0x97, 0x98:
		return -1
	case 0x84: // iinc
		return 0
	case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
		0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
		0xc6, 0xc7:
		return -1 // pop condition
	case 0xa7, 0xc8: // goto
		return 0
	case 0xb2: // getstatic
		return 1
	case 0xb3: // putstatic
		return -1
	case 0xb4: // getfield
		return 0 // pop ref + push value = net 0
	case 0xb5: // putfield
		return -2
	case 0xb6, 0xb7, 0xb9: // invoke virtual/special/interface
		// Can't know exact without descriptor; conservative estimate 0.
		return 0
	case 0xb8: // invokestatic
		return 0 // conservative estimate
	case 0xbb: // new
		return 1
	case 0xbc: // newarray
		return 0 // pop count + push arrref
	case 0xbe: // arraylength
		return 0
	case 0xbf: // athrow
		return -1
	case 0xc0: // checkcast
		return 0
	case 0xc1: // instanceof
		return 0
	case 0xc2, 0xc3: // monitorenter/exit
		return -1
	default:
		return 0
	}
}

// instrSizeSafe returns the byte length of the instruction at pc.
func instrSizeSafe(code []byte, pc int) int {
	if pc >= len(code) {
		return 1
	}
	op := code[pc]
	switch op {
	case 0xaa, 0xab: // tableswitch / lookupswitch
		p := pc + 1
		for p%4 != 0 {
			p++
		}
		rd32 := func(i int) int32 {
			return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
				uint32(code[i+2])<<8 | uint32(code[i+3]))
		}
		p += 4 // default
		if op == 0xab {
			n := int(rd32(p))
			return p - pc + 4 + 8*n
		}
		low := int(rd32(p))
		high := int(rd32(p + 4))
		return p - pc + 8 + 4*(high-low+1)
	case 0xc4: // wide
		if code[pc+1] == 0x84 {
			return 6
		}
		return 4
	case 0x10, 0x12, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x36, 0x37, 0x38, 0x39, 0x3a, 0xa9, 0xbc:
		return 2
	case 0x11, 0x13, 0x14, 0x84,
		0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
		0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
		0xa7, 0xa8,
		0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
		0xbb, 0xbd, 0xc0, 0xc1, 0xc6, 0xc7:
		return 3
	case 0xb9, 0xc5, 0xc8, 0xc9:
		return 5
	}
	return 1
}

// instrEffect computes the exact net stack effect of the instruction at pc,
// resolving method/field descriptors from the constant pool.
func (e *methodEmitter) instrEffect(pc int, op byte) int {
	code := e.code
	switch op {
	case 0xb6, 0xb7, 0xb8, 0xb9: // invokes
		_, _, desc, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return 0
		}
		argDescs, retDesc, derr := splitMethodDesc(desc)
		if derr != nil {
			return 0
		}
		total := 0
		for _, a := range argDescs {
			total += descSlots(a)
		}
		if op != 0xb8 {
			total++ // receiver
		}
		ret := descSlots(retDesc)
		return ret - total
	case 0xb4: // getfield
		_, _, desc, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return 0
		}
		return descSlots(desc) // pop recv + push value
	case 0xb5: // putfield
		_, _, desc, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return 0
		}
		return -(1 + descSlots(desc)) // pop recv + value
	case 0xb2: // getstatic
		_, _, desc, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return 1
		}
		return descSlots(desc)
	case 0xb3: // putstatic
		_, _, desc, _, err := refAt(e.cf, code, pc)
		if err != nil {
			return -1
		}
		return -descSlots(desc)
	}
	return netStackEffect(op)
}
