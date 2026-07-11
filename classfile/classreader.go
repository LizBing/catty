package classfile

import "encoding/binary"

// ClassReader wraps a class file's bytes and reads JVMS u1/u2/u4 values in
// big-endian order. It is the single low-level input primitive every parser in
// this package is built on, so read methods are kept allocation-free and inlined
// by the compiler. Panic on underflow: a malformed class file is a fatal error
// during loading, and recovering only to report "bad class file" adds no value.
type ClassReader struct {
	data []byte
}

func NewClassReader(data []byte) *ClassReader {
	return &ClassReader{data: data}
}

func (r *ClassReader) ReadUint8() uint8 {
	v := r.data[0]
	r.data = r.data[1:]
	return v
}

func (r *ClassReader) ReadUint16() uint16 {
	v := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]
	return v
}

func (r *ClassReader) ReadUint32() uint32 {
	v := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]
	return v
}

func (r *ClassReader) ReadUint64() uint64 {
	v := binary.BigEndian.Uint64(r.data)
	r.data = r.data[8:]
	return v
}

// ReadUint16s reads a u2 count followed by that many u2 values, returning the
// values without the leading count. This is the layout of interfaces[] and
// most "length-prefixed" arrays in the class file format.
func (r *ClassReader) ReadUint16s() []uint16 {
	n := r.ReadUint16()
	s := make([]uint16, n)
	for i := range s {
		s[i] = r.ReadUint16()
	}
	return s
}

// ReadBytes reads a u4 length followed by that many raw bytes.
func (r *ClassReader) ReadBytes() []byte {
	n := r.ReadUint32()
	b := make([]byte, n)
	copy(b, r.data[:n])
	r.data = r.data[n:]
	return b
}
