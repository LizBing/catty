// Package classfile parses JVM class files into immutable definitions with
// original bytecode offsets. It performs binary-level validation only; the
// verifier and builder decide profile admission.
package classfile

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

type cpEntry struct {
	tag  uint8
	text string
	unit []uint16

	classIndex      uint16
	nameTypeIndex   uint16
	stringIndex     uint16
	first           uint16
	second          uint16
	intVal          int32
	nameIndex       uint16
	descriptorIndex uint16
	bootstrapIndex  uint16
	refIndex        uint16
	refKind         uint8
}

type constantPool struct {
	entries []cpEntry
}

func (cp *constantPool) entry(i uint16) (*cpEntry, error) {
	if int(i) >= len(cp.entries) || i == 0 {
		return nil, fmt.Errorf("constant pool index %d out of range", i)
	}
	e := &cp.entries[i]
	if e.tag == 0 {
		return nil, fmt.Errorf("constant pool index %d names an unusable slot", i)
	}
	return e, nil
}

func (cp *constantPool) utf8(i uint16) (string, []uint16, error) {
	e, err := cp.entry(i)
	if err != nil {
		return "", nil, err
	}
	if e.tag != 1 {
		return "", nil, fmt.Errorf("constant pool index %d is not a Utf8 constant", i)
	}
	return e.text, e.unit, nil
}

func (cp *constantPool) className(i uint16) (string, error) {
	e, err := cp.entry(i)
	if err != nil {
		return "", err
	}
	if e.tag != 7 {
		return "", fmt.Errorf("constant pool index %d is not a Class constant", i)
	}
	name, _, err := cp.utf8(e.classIndex)
	return name, err
}

func (cp *constantPool) nameAndType(i uint16) (string, string, error) {
	e, err := cp.entry(i)
	if err != nil {
		return "", "", err
	}
	if e.tag != 12 {
		return "", "", fmt.Errorf("constant pool index %d is not a NameAndType constant", i)
	}
	name, _, err := cp.utf8(e.nameIndex)
	if err != nil {
		return "", "", err
	}
	desc, _, err := cp.utf8(e.descriptorIndex)
	return name, desc, err
}

func (cp *constantPool) memberRef(i uint16) (owner, name, desc, kind string, err error) {
	e, err := cp.entry(i)
	if err != nil {
		return "", "", "", "", err
	}
	switch e.tag {
	case 9:
		kind = "field"
	case 10:
		kind = "method"
	case 11:
		kind = "interfacemethod"
	default:
		return "", "", "", "", fmt.Errorf("constant pool index %d is not a member reference", i)
	}
	owner, err = cp.className(e.first)
	if err != nil {
		return "", "", "", "", err
	}
	name, desc, err = cp.nameAndType(e.second)
	return owner, name, desc, kind, err
}

func (cp *constantPool) ldcString(i uint16) ([]uint16, error) {
	e, err := cp.entry(i)
	if err != nil {
		return nil, err
	}
	if e.tag != 8 {
		return nil, fmt.Errorf("constant pool index %d is not a String constant", i)
	}
	_, units, err := cp.utf8(e.stringIndex)
	return units, err
}

func (cp *constantPool) ldcInt(i uint16) (int32, error) {
	e, err := cp.entry(i)
	if err != nil {
		return 0, err
	}
	if e.tag != 3 {
		return 0, fmt.Errorf("constant pool index %d is not an Integer constant", i)
	}
	return e.intVal, nil
}

type reader struct {
	data []byte
	off  int
}

func (r *reader) need(n int) error {
	if n < 0 || r.off+n > len(r.data) {
		return fmt.Errorf("truncated class file at offset %d (need %d bytes)", r.off, n)
	}
	return nil
}

func (r *reader) u1() (uint8, error) {
	if err := r.need(1); err != nil {
		return 0, err
	}
	v := r.data[r.off]
	r.off++
	return v, nil
}

func (r *reader) u2() (uint16, error) {
	if err := r.need(2); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint16(r.data[r.off:])
	r.off += 2
	return v, nil
}

func (r *reader) u4() (uint32, error) {
	if err := r.need(4); err != nil {
		return 0, err
	}
	v := binary.BigEndian.Uint32(r.data[r.off:])
	r.off += 4
	return v, nil
}

func (r *reader) bytes(n int) ([]byte, error) {
	if err := r.need(n); err != nil {
		return nil, err
	}
	v := r.data[r.off : r.off+n]
	r.off += n
	return v, nil
}

// decodeMUTF8 decodes the JVM's modified UTF-8 into UTF-16 code units without
// collapsing isolated surrogates.
func decodeMUTF8(b []byte) ([]uint16, error) {
	units := make([]uint16, 0, len(b))
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c == 0:
			return nil, fmt.Errorf("modified UTF-8 contains a zero byte")
		case c < 0x80:
			units = append(units, uint16(c))
			i++
		case c&0xE0 == 0xC0:
			if i+1 >= len(b) {
				return nil, fmt.Errorf("truncated modified UTF-8")
			}
			v := uint16(c&0x1F)<<6 | uint16(b[i+1]&0x3F)
			units = append(units, v)
			i += 2
		case c&0xF0 == 0xE0:
			if i+2 >= len(b) {
				return nil, fmt.Errorf("truncated modified UTF-8")
			}
			v := uint16(c&0x0F)<<12 | uint16(b[i+1]&0x3F)<<6 | uint16(b[i+2]&0x3F)
			units = append(units, v)
			i += 3
		default:
			return nil, fmt.Errorf("invalid modified UTF-8 lead byte 0x%02x", c)
		}
	}
	return units, nil
}

func unitsText(units []uint16) string {
	return string(utf16.Decode(units))
}
