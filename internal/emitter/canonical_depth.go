package emitter

import (
	"catty/internal/classfile"
)

// qItem is a worklist entry for depth simulation.
type qItem struct{ pc, depth int }

// computeCanonicalDepths simulates operand-stack depth for every reachable
// pc via worklist propagation along CFG edges. Ground truth: StackMap
// frames override at their offsets; exception-handler entries are always
// depth 1 (JVMS §2.6 stack clear + exception ref).
//
// History note: the original implementation walked pcs LINEARLY, letting
// a taken-branch path's depth leak into unreachable straight-line
// instructions (e.g. the fall-through arm after goto) — masked for years
// by dense SM resyncs, exposed by minimal-json's getDouble merge.
func (e *methodEmitter) computeCanonicalDepths(reach map[int]bool) map[int]int {
	depths := make(map[int]int)
	code := e.code

	// Build SM frame lookup (raw slot count, cat2 = 2).
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

	applyGroundTruth := func(pc, d int) int {
		if handlerSet[pc] {
			return 1
		}
		if sd, ok := smDepth[pc]; ok {
			return sd
		}
		return d
	}

	enqueue := func(list []qItem, pc, d int) []qItem {
		if !reach[pc] {
			return list
		}
		d = applyGroundTruth(pc, d)
		if prev, seen := depths[pc]; seen {
			if prev == d {
				return list // converged
			}
			// Verifier-checked input guarantees consistent frames; on
			// any residual disagreement keep the established value and
			// do not re-propagate (deterministic, no oscillation).
			return list
		}
		depths[pc] = d
		return append(list, qItem{pc, d})
	}

	work := make([]qItem, 0, len(reach))
	work = enqueue(work, 0, 0) // method entry: empty operand stack

	for len(work) > 0 {
		it := work[len(work)-1]
		work = work[:len(work)-1]
		op := code[it.pc]

		next := -1 // linear successor when the opcode falls through
		switch op {
		case 0xa7, 0xc8: // goto, goto_w
			next = branchTarget(it.pc, code)
			if op == 0xc8 {
				p := it.pc + 1
				next = it.pc + int(int32(uint32(code[p])<<24|uint32(code[p+1])<<16|
					uint32(code[p+2])<<8|uint32(code[p+3])))
			}
		case 0xaa, 0xab: // tableswitch, lookupswitch
			p := it.pc + 1
			for p%4 != 0 {
				p++
			}
			rd32 := func(i int) int32 {
				return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
					uint32(code[i+2])<<8 | uint32(code[i+3]))
			}
			emit := func(tgt int) { work = enqueue(work, tgt, it.depth) }
			if op == 0xaa {
				def := it.pc + int(rd32(p))
				low, high := int(rd32(p+4)), int(rd32(p+8))
				p += 12
				emit(def)
				for v := low; v <= high; v++ {
					emit(it.pc + int(rd32(p)))
					p += 4
				}
			} else {
				def := it.pc + int(rd32(p))
				pairs := int(rd32(p+4))
				p += 8
				emit(def)
				for j := 0; j < pairs; j++ {
					emit(it.pc + int(rd32(p+4)))
					p += 8
				}
			}
		case 0xa9, 0xbf, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb1:
			// ret/athrow/returns: no linear successor. Exception edges
			// land on handler entries, seeded globally below.
		case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
			0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
			0xc6, 0xc7: // conditional branches: both edges
			work = enqueue(work, branchTarget(it.pc, code), it.depth)
			next = it.pc + instrSizeSafe(code, it.pc)
		default:
			next = it.pc + instrSizeSafe(code, it.pc)
		}
		if next >= 0 && next <= len(code) {
			d := it.depth + e.instrEffect(it.pc, op)
			if d < 0 {
				d = 0
			}
			work = enqueue(work, next, d)
		}
	}

	// Handler entries not reached by explicit edges still canonicalize.
	for pc := range handlerSet {
		if _, ok := depths[pc]; !ok && reach[pc] {
			depths[pc] = 1
		}
	}
	return depths
}

// netStackEffect returns how many raw slots an opcode adds to the operand
// stack. Fixed-effect opcodes delegate to classfile.StackEffect — the
// single source of truth shared across engines (DEBT-0015/0017/0019
// drift class); only descriptor-shaped opcodes fall through to
// conservative defaults here, and instrEffect refines those from the
// constant pool.
func netStackEffect(op byte) int {
	if n, class := classfile.StackEffect(op); class == classfile.EffectFixed {
		return n
	}
	switch op {
	case 0xb2: // getstatic: conservative before pool resolution
		return 1
	case 0xbf, 0xc2, 0xc3, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
		0xc6, 0xc7: // unreachable via EffectFixed above; kept for clarity
		return -1
	}
	return 0
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
		// Net effect: pop receiver, push value ⇒ slots-1. The missing
		// -1 here inflated straight-line depth by one per getfield —
		// masked for years by dense StackMap resyncs, exposed when the
		// lazy-install hook first ran minimal-json's generated bodies.
		return descSlots(desc) - 1
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
