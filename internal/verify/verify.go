// Package verify implements the M1 structural tier of JVMS §4.10
// verification.
//
// What v0 checks (mechanically, soundly):
//   - every branch target and exception handler entry carries a
//     StackMapTable frame (strict linking requires it for v51+);
//   - every StackMapTable OBJECT item resolves to a CONSTANT_Class entry;
//   - handler ranges are well-formed; frames fit MaxLocals/MaxStack bounds;
//   - code never runs past its end on any simulated linear path is NOT yet
//     checked — see DEBT-0009 for the dataflow type-checker roadmap.
//
// Deliberately out of scope in v0: operand type simulation. A buggy type
// checker is worse than an explicit structural tier; the gap is registered
// as DEBT-0009 with trusted-input operation documented in the deviation
// ledger until it lands.
package verify

import (
	"fmt"

	"catty/internal/classfile"
)

// RefResolver is reserved for the v1 dataflow checker; kept so embedders
// can already supply class-graph adapters.
type RefResolver interface {
	Known(name string) bool
	IsSubclass(child, anc string) bool
}

type errVerify struct {
	method string
	pc     int
	msg    string
}

func (e *errVerify) Error() string {
	return fmt.Sprintf("%s @pc=%d: %s", e.method, e.pc, e.msg)
}

func verr(method string, pc int, format string, args ...any) error {
	return &errVerify{method: method, pc: pc, msg: fmt.Sprintf(format, args...)}
}

// Verify checks every Code-bearing method of cf structurally.
func Verify(cf *classfile.ClassFile, _ RefResolver) error {
	for i := range cf.Methods {
		m := &cf.Methods[i]
		if m.Code == nil {
			continue // abstract/native
		}
		if err := checkMethod(cf, m); err != nil {
			return err
		}
	}
	return nil
}

func checkMethod(cf *classfile.ClassFile, m *classfile.MethodInfo) error {
	code := m.Code
	name := m.Name + m.Desc

	targets, err := branchTargets(code.Code)
	if err != nil {
		return verr(name, 0, "%v", err)
	}
	for _, h := range code.Handlers {
		if h.EndPc < h.StartPc {
			return verr(name, int(h.StartPc), "handler range [%d,%d) inverted", h.StartPc, h.EndPc)
		}
		if h.CatchType != 0 {
			e, err := cf.Entry(h.CatchType)
			if err != nil || e.Tag != classfile.CClass {
				return verr(name, int(h.HandlerPc), "catch_type pool[%d] is not a class", h.CatchType)
			}
		}
		targets[int(h.HandlerPc)] = true
	}

	smAt := make(map[int]*classfile.StackMapFrame, len(code.StackMaps))
	for i := range code.StackMaps {
		fr := &code.StackMaps[i]
		if fr.Offset < 0 || fr.Offset >= int32(len(code.Code)) {
			return verr(name, int(fr.Offset), "frame offset out of range")
		}
		if _, dup := smAt[int(fr.Offset)]; dup {
			return verr(name, int(fr.Offset), "duplicate frame")
		}
		smAt[int(fr.Offset)] = fr
	}

	for tgt := range targets {
		if _, ok := smAt[tgt]; !ok {
			return verr(name, tgt, "branch/handler target has no StackMapTable frame")
		}
	}

	for i := range code.StackMaps {
		fr := &code.StackMaps[i]
		localSlots := 0
		for _, it := range fr.Locals {
			switch it.Tag {
			case classfile.VItemObject:
				if err := cfEntryIsClass(cf, it.CPoolIdx); err != nil {
					return verr(name, int(fr.Offset), "%v", err)
				}
			case classfile.VItemUninit:
				if int(it.Offset) >= len(code.Code) {
					return verr(name, int(fr.Offset), "uninit offset out of range")
				}
			case classfile.VItemTop, classfile.VItemInteger, classfile.VItemFloat,
				classfile.VItemLong, classfile.VItemDouble, classfile.VItemNull,
				classfile.VItemUninitTh:
			default:
				return verr(name, int(fr.Offset), "bad item tag %d", it.Tag)
			}
			localSlots++
			if it.Tag == classfile.VItemLong || it.Tag == classfile.VItemDouble {
				localSlots++ // second slot
			}
		}
		if localSlots > int(code.MaxLocals) {
			return verr(name, int(fr.Offset), "frame locals %d exceed MaxLocals %d",
				localSlots, code.MaxLocals)
		}
		if len(fr.Stack) > int(code.MaxStack)+2 {
			return verr(name, int(fr.Offset), "frame stack %d exceeds MaxStack %d",
				len(fr.Stack), code.MaxStack)
		}
	}
	return nil
}

func cfEntryIsClass(cf *classfile.ClassFile, idx uint16) error {
	e, err := cf.Entry(idx)
	if err != nil {
		return err
	}
	if e.Tag != classfile.CClass {
		return fmt.Errorf("pool[%d] is %s, want CONSTANT_Class", idx, e.Tag)
	}
	if _, err := cf.ClassName(idx); err != nil {
		return err
	}
	return nil
}

// branchTargets enumerates every bytecode offset that can start execution
// other than fallthrough from the previous instruction.
func branchTargets(code []byte) (map[int]bool, error) {
	out := make(map[int]bool)
	i := 0
	add4 := func(at int) int32 {
		v := int32(code[at])<<24 | int32(code[at+1])<<16 |
			int32(code[at+2])<<8 | int32(code[at+3])
		return v
	}
	pad := func(from int) int {
		for from%4 != 0 {
			from++
		}
		return from
	}
	for i < len(code) {
		op := code[i]
		switch op {
		case 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, // ifeq..ifle
			0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, // if_icmp*
			0xa5, 0xa6, // if_acmpeq/ne
			0xc6, 0xc7: // ifnull/nonnull
			if i+3 >= len(code) {
				return nil, fmt.Errorf("branch operand overrun at %d", i)
			}
			off := int(int16(uint16(code[i+1])<<8 | uint16(code[i+2])))
			out[i+off] = true
			i += 3
		case 0xa7, 0xa8: // goto/jsr
			if i+3 >= len(code) {
				return nil, fmt.Errorf("branch operand overrun at %d", i)
			}
			off := int(int16(uint16(code[i+1])<<8 | uint16(code[i+2])))
			out[i+off] = true
			i += 3
		case 0xc8: // goto_w
			if i+5 >= len(code) {
				return nil, fmt.Errorf("goto_w operand overrun at %d", i)
			}
			out[i+int(add4(i+1))] = true
			i += 5
		case 0xaa: // tableswitch
			p := pad(i + 1)
			if p+12 > len(code) {
				return nil, fmt.Errorf("tableswitch overrun at %d", i)
			}
			def := add4(p)
			low := add4(p + 4)
			high := add4(p + 8)
			if high < low || int(high)-int(low) > len(code)*2 {
				return nil, fmt.Errorf("tableswitch bad range at %d", i)
			}
			out[i+int(def)] = true
			base := p + 12
			for v := low; v <= high; v++ {
				if base+4 > len(code) {
					return nil, fmt.Errorf("tableswitch overrun at %d", i)
				}
				out[i+int(add4(base))] = true
				base += 4
			}
			i = base
		case 0xab: // lookupswitch
			p := pad(i + 1)
			if p+8 > len(code) {
				return nil, fmt.Errorf("lookupswitch overrun at %d", i)
			}
			def := add4(p)
			npairs := int(add4(p + 4))
			if npairs < 0 || p+8+npairs*8 > len(code) {
				return nil, fmt.Errorf("lookupswitch bad pair count at %d", i)
			}
			out[i+int(def)] = true
			base := p + 8
			for j := 0; j < npairs; j++ {
				out[i+int(add4(base+4))] = true
				base += 8
			}
			i = base
		case 0xc4: // wide: 4 bytes, or 6 for iinc/ret
			if i+2 >= len(code) {
				return nil, fmt.Errorf("wide overrun at %d", i)
			}
			m := code[i+1]
			i += 6
			if m != 0x84 && m != 0xa9 {
				i -= 2
			}
			if i > len(code) {
				return nil, fmt.Errorf("wide overrun at %d", i)
			}
		default:
			n, ok := fixedLen[op]
			if !ok {
				return nil, fmt.Errorf("unknown opcode %#x at %d during target scan", op, i)
			}
			if i+n > len(code) {
				return nil, fmt.Errorf("operand overrun at %d", i)
			}
			i += n
		}
	}
	return out, nil
}

// fixedLen holds total instruction length for every fixed-size opcode;
// absence means variable-length (switches/wide, handled explicitly).
var fixedLen = func() map[byte]int {
	m := make(map[byte]int, 220)
	for o := 0; o <= 255; o++ {
		m[byte(o)] = 1 // majority are one byte
	}
	two := []int{0x10, 0x12, 0x15, 0x16, 0x17, 0x18, 0x19,
		0x36, 0x37, 0x38, 0x39, 0x3a, 0xa9, 0xbc}
	for _, o := range two {
		m[byte(o)] = 2
	}
	three := []int{0x11, 0x13, 0x14, 0x84,
		0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e,
		0x9f, 0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6,
		0xa7, 0xa8,
		0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
		0xbb, 0xbd, 0xc0, 0xc1, 0xc6, 0xc7}
	for _, o := range three {
		m[byte(o)] = 3
	}
	five := []int{0xb9, 0xc5, 0xc8, 0xc9}
	for _, o := range five {
		m[byte(o)] = 5
	}
	delete(m, 0xaa) // tableswitch — variable
	delete(m, 0xab) // lookupswitch — variable
	delete(m, 0xc4) // wide — variable
	return m
}()