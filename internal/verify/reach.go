package verify

import (
	"fmt"
	"os"
)

// branchTarget resolves a two-byte signed branch offset at pc.
func branchTarget(pc int, code []byte) int {
	return pc + int(int16(uint16(code[pc+1])<<8 | uint16(code[pc+2])))
}

// Reachable computes reachable instruction start offsets from pc 0,
// following fallthrough and all control transfers. handlerPCs are treated
// as additional roots (exception entries are always reachable).
func Reachable(code []byte, handlerPCs []int) (map[int]bool, error) {
	reach := map[int]bool{}
	queue := []int{0}
	for len(queue) > 0 {
		pc := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		if pc < 0 || pc >= len(code) || reach[pc] {
			continue
		}
		reach[pc] = true
		if os.Getenv("CATTY_RTRACE") != "" {
			fmt.Fprintf(os.Stderr, "[reach] marking pc=%d op=%#x\n", pc, code[pc])
		}
		op := code[pc]
		switch op {
		case 0xa9, 0xba, 0xc9: // ret / invokedynamic / jsr_w
			return nil, fmt.Errorf("unsupported opcode %#x", op)
		case 0xbf, 0xac, 0xad, 0xae, 0xaf, 0xb0, 0xb1:
			// athrow & returns: no linear successor
		case 0xa7, 0xc8: // goto / goto_w
			tgt := branchTarget(pc, code)
			if tgt < 0 || tgt >= len(code) {
				return nil, fmt.Errorf("goto target %d out of range", tgt)
			}
			queue = append(queue, tgt)
			next := instrNext(pc)
			if next < len(code) && !reach[next] {
				queue = append(queue, next) // dead-code after goto still scanned? No—unreachable; skip.
			}
		case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
			0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
			0xc6, 0xc7:
			tgt := branchTarget(pc, code)
			queue = append(queue, tgt)
			queue = append(queue, pc+3)
		case 0xaa, 0xab:
			for _, t := range switchTargets(code, pc) {
				queue = append(queue, t)
			}
		default:
			nxt, ok := instrNextSafe(code, pc)
			if !ok {
				return nil, fmt.Errorf("unknown opcode %#x at %d", op, pc)
			}
			queue = append(queue, nxt)
		}
	}
	for _, hp := range handlerPCs {
		reach[hp] = true
	}
	return reach, nil
}

// instrNext returns the next linear pc accounting for operand widths.
func instrNext(pc int) int { return pc + 3 } // used for goto-family only

func instrNextSafe(code []byte, pc int) (int, bool) {
	op := code[pc]
	switch op {
	case 0xaa, 0xab, 0xc4: // variable-length handled by callers where needed
		if op == 0xc4 {
			m := code[pc+1]
			if m == 0x84 || m == 0xa9 {
				return pc + 6, true
			}
			return pc + 4, true
		}
		return switchNext(code, pc)
	}
	n, ok := fixedLenCache[op]
	if !ok {
		return 0, false
	}
	return pc + n, true
}

func switchLen(code []byte, pc int) int {
	p := pc + 1
	for p%4 != 0 {
		p++
	}
	rd32 := func(i int) int32 {
		return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
			uint32(code[i+2])<<8 | uint32(code[i+3]))
	}
	def := rd32(p)
	_ = def
	p += 4
	if code[pc] == 0xab {
		n := int(rd32(p))
		p += 4 + 8*n
	} else {
		low := int(rd32(p))
		high := int(rd32(p+4))
		p += 8 + 4*(high-low+1)
	}
	return p
}

// switchTargets enumerates jump targets of a table/lookup switch at pc.
func switchTargets(code []byte, pc int) []int {
	var out []int
	p := pc + 1
	for p%4 != 0 {
		p++
	}
	rd32 := func(i int) int32 {
		return int32(uint32(code[i])<<24 | uint32(code[i+1])<<16 |
			uint32(code[i+2])<<8 | uint32(code[i+3]))
	}
	out = append(out, pc+int(rd32(p)))
	p += 4
	if code[pc] == 0xab {
		n := int(rd32(p))
		p += 4
		for j := 0; j < n; j++ {
			out = append(out, pc+int(rd32(p+4)))
			p += 8
		}
		return out
	}
	low := int(rd32(p))
	high := int(rd32(p + 4))
	p += 12
	for v := low; v <= high; v++ {
		_ = v
		out = append(out, pc+int(rd32(p)))
		p += 4
	}
	return out
}


// switchNext returns the next linear pc after a variable-length switch.
func switchNext(code []byte, pc int) (int, bool) {
	return switchLen(code, pc), true
}

// fixedLenCache holds total instruction length for fixed-size opcodes.
var fixedLenCache = buildFixedLenTable()

func buildFixedLenTable() map[byte]int {
	m := map[byte]int{}
	for o := 0; o < 256; o++ {
		m[byte(o)] = 1
	}
	for _, o := range []int{0x10, 0x12, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x36, 0x37, 0x38, 0x39, 0x3a, 0xa9, 0xbc} {
		m[byte(o)] = 2
	}
	for _, o := range []int{0x11, 0x13, 0x14, 0x84,
		0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
		0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
		0xa7, 0xa8,
		0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
		0xbb, 0xbd, 0xc0, 0xc1, 0xc6, 0xc7} {
		m[byte(o)] = 3
	}
	for _, o := range []int{0xb9, 0xc5, 0xc8, 0xc9} {
		m[byte(o)] = 5
	}
	delete(m, 0xaa)
	delete(m, 0xab)
	delete(m, 0xc4)
	return m
}


// SwitchNext returns the next linear pc after a variable-length switch.
func SwitchNext(code []byte, pc int) (int, bool) {
	return switchNext(code, pc)
}

// FixedLen returns the fixed instruction length for an opcode, or false
// for variable-length opcodes.
func FixedLen(op byte) (int, bool) {
	n, ok := fixedLenCache[op]
	return n, ok
}
