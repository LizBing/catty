package classfile

import (
	"encoding/binary"
	"fmt"
)

// FormatError is a typed classfile format error returned by Parse. Each error
// names the operation that failed and a human-readable detail.
type FormatError struct {
	Op  string
	Msg string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("catty: classfile format error: %s: %s", e.Op, e.Msg)
}

// parsePanic is the private sentinel the parser uses for control flow. Reader
// primitives and decoders raise it on malformed input so that Parse can recover
// it into a typed *FormatError without threading error returns through every
// nested call. Programming faults (nil dereference, slice bounds in untracked
// code) use a different panic type and are re-raised.
type parsePanic struct {
	err *FormatError
}

// panicf raises a parsePanic with a fmt.Sprintf-formatted message.
func panicf(op, format string, args ...any) {
	panic(parsePanic{err: &FormatError{Op: op, Msg: fmt.Sprintf(format, args...)}})
}

// ClassReader wraps a class file's bytes and reads JVMS u1/u2/u4 values in
// big-endian order. Every read method checks bounds and raises parsePanic on
// truncation — callers must not index into data directly.
type ClassReader struct {
	data []byte
}

// NewClassReader returns a ClassReader over the given byte slice.
func NewClassReader(data []byte) *ClassReader {
	return &ClassReader{data: data}
}

// Len returns the number of bytes remaining.
func (r *ClassReader) Len() int { return len(r.data) }

// ReadUint8 reads a u1. Panics with parsePanic on truncation.
func (r *ClassReader) ReadUint8() uint8 {
	if len(r.data) < 1 {
		panicf("read", "truncated: expected u1, have %d bytes", len(r.data))
	}
	v := r.data[0]
	r.data = r.data[1:]
	return v
}

// ReadUint16 reads a big-endian u2.
func (r *ClassReader) ReadUint16() uint16 {
	if len(r.data) < 2 {
		panicf("read", "truncated: expected u2, have %d bytes", len(r.data))
	}
	v := binary.BigEndian.Uint16(r.data)
	r.data = r.data[2:]
	return v
}

// ReadUint32 reads a big-endian u4.
func (r *ClassReader) ReadUint32() uint32 {
	if len(r.data) < 4 {
		panicf("read", "truncated: expected u4, have %d bytes", len(r.data))
	}
	v := binary.BigEndian.Uint32(r.data)
	r.data = r.data[4:]
	return v
}

// ReadUint64 reads a big-endian u8.
func (r *ClassReader) ReadUint64() uint64 {
	if len(r.data) < 8 {
		panicf("read", "truncated: expected u8, have %d bytes", len(r.data))
	}
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

// ReadBytes reads a u4 length followed by that many raw bytes. The returned
// slice is a copy; the reader's cursor advances past the bytes.
func (r *ClassReader) ReadBytes() []byte {
	n := r.ReadUint32()
	if uint32(len(r.data)) < n {
		panicf("read", "truncated: expected %d bytes, have %d", n, len(r.data))
	}
	b := make([]byte, n)
	copy(b, r.data[:n])
	r.data = r.data[n:]
	return b
}

// ReadSlice reads n raw bytes directly from the reader, returning a sub-slice
// and advancing the cursor. The returned slice aliases the reader's backing
// array — callers must not mutate it or retain it after the reader is gone.
func (r *ClassReader) ReadSlice(n uint32) []byte {
	if uint32(len(r.data)) < n {
		panicf("read", "truncated: expected %d bytes, have %d", n, len(r.data))
	}
	s := r.data[:n]
	r.data = r.data[n:]
	return s
}

// Skip advances the cursor by n bytes without allocating.
func (r *ClassReader) Skip(n uint32) {
	if uint32(len(r.data)) < n {
		panicf("read", "truncated: expected %d bytes to skip, have %d", n, len(r.data))
	}
	r.data = r.data[n:]
}
