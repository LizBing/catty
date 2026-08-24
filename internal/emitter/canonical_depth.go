package emitter

import (
	"fmt"
	"os"

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
	if os.Getenv("CATTY_SIM") != "" {
		fmt.Fprintf(os.Stderr, "[cc] %s.%s reach=%d\n", e.cf.ThisClass, e.m.Name, len(reach))
	}

	// JVMS: the operand stack is EMPTY at method entry; parameters live
	// in locals. Counting them here inflated every depth by the arg
	// footprint (the root cause of the Json.value(F) s-1 drift).
	entryDepth := 0
	if os.Getenv("CATTY_SIM") != "" {
		fmt.Fprintf(os.Stderr, "[cc-v2] %s.%s entry=0\n", e.cf.ThisClass, e.m.Name)
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
		if os.Getenv("CATTY_SIM") != "" &&
			((pc >= 35 && pc <= 45) || (pc >= 155 && pc <= 170)) {
			fmt.Fprintf(os.Stderr, "[sim] pc=%d op=%#x rec=%d\n", pc, op, d)
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
	if os.Getenv("CATTY_SIM") != "" && e.cf.ThisClass == "com/eclipsesource/json/Json" {
		fmt.Fprintf(os.Stderr, "[cc-done] %s.%s depths=%v\n", e.cf.ThisClass, e.m.Name, depths)
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
