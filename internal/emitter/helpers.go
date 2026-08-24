package emitter

import (
	"fmt"

	"catty/internal/classfile"
)

// classAt reads a CONSTANT_Class operand at code[pc+1..pc+2].
func classAt(cf *classfile.ClassFile, code []byte, pc int) (string, error) {
	idx := uint16(code[pc+1])<<8 | uint16(code[pc+2])
	return cf.ClassName(idx)
}

// retOf reports whether a descriptor returns a value.
func retOf(desc string) (struct{}, bool, error) {
	i := len(desc) - 1
	if i < 0 || desc[i] != 'V' {
		return struct{}{}, true, nil
	}
	if len(desc) < 2 || desc[0] != '(' {
		return struct{}{}, false, fmt.Errorf("bad descriptor %q", desc)
	}
	return struct{}{}, false, nil
}

// lineAtPC resolves the source line for pc via the method's
// LineNumberTable (0 when absent — caller skips SetLine emission).
func (e *methodEmitter) lineAtPC(pc int) int32 {
	if e.m.Code == nil || len(e.m.Code.LineNumbers) == 0 {
		return 0
	}
	lns := e.m.Code.LineNumbers
	line := int32(0)
	for i := range lns {
		if int(lns[i].StartPc) <= pc {
			line = int32(lns[i].Line)
		} else {
			break
		}
	}
	return line
}
