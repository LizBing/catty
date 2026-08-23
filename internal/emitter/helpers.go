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
