package rtda

import (
	"math"
	"sync/atomic"
)

// HeapCell is a race-free, sequentially consistent unit of Java heap storage
// per ADR-0030. Each instance field, static field, and array element is backed
// by exactly one HeapCell regardless of its Java type. Long and double values
// are stored atomically in a single 64-bit cell; there is no 2-slot category-2
// layout in the heap.
//
// Frame operand stacks and local variables continue to use the existing Slot
// type. HeapCell is only for shared heap storage.
type HeapCell struct {
	bits atomic.Int64           // primitive value (32-bit sign/zero-extended to 64; 64-bit direct)
	ref  atomic.Pointer[Object] // object reference (nil is Java null)
}

// --- Primitive typed accessors ---

// GetInt returns the cell's value as int32. Caller must know the field is
// a 32-bit or smaller primitive (int, short, char, byte, boolean, or float-as-bits).
func (c *HeapCell) GetInt() int32 { return int32(c.bits.Load()) }

// SetInt stores v as a primitive cell value.
func (c *HeapCell) SetInt(v int32) { c.bits.Store(int64(v)) }

// GetLong returns the cell's value as int64 (also used for double bits).
func (c *HeapCell) GetLong() int64 { return c.bits.Load() }

// SetLong stores v as a 64-bit primitive cell value.
func (c *HeapCell) SetLong(v int64) { c.bits.Store(v) }

// GetFloat returns the cell's int32 bits as float32.
func (c *HeapCell) GetFloat() float32 { return math.Float32frombits(uint32(c.bits.Load())) }

// SetFloat stores v's bits as a primitive cell value.
func (c *HeapCell) SetFloat(v float32) { c.bits.Store(int64(math.Float32bits(v))) }

// GetDouble returns the cell's int64 bits as float64.
func (c *HeapCell) GetDouble() float64 { return math.Float64frombits(uint64(c.bits.Load())) }

// SetDouble stores v's bits as a 64-bit primitive cell value.
func (c *HeapCell) SetDouble(v float64) { c.bits.Store(int64(math.Float64bits(v))) }

// --- Reference typed accessors ---

// GetRef returns the stored object reference, or nil (Java null).
func (c *HeapCell) GetRef() *Object { return c.ref.Load() }

// SetRef stores an object reference. Pass nil for Java null.
func (c *HeapCell) SetRef(r *Object) { c.ref.Store(r) }

// --- Bulk helpers ---

// NewHeapCells allocates n zero-initialized heap cells.
func NewHeapCells(n int) []HeapCell { return make([]HeapCell, n) }

// CopyHeapCells copies src into dst by element, using atomic stores.
// Used for clone and arraycopy.
func CopyHeapCells(dst, src []HeapCell) {
	for i := range src {
		if i >= len(dst) {
			break
		}
		dst[i].bits.Store(src[i].bits.Load())
		dst[i].ref.Store(src[i].ref.Load())
	}
}

// cloneHeapCells returns a new slice with independent atomic cells copied from src.
func cloneHeapCells(src []HeapCell) []HeapCell {
	dst := make([]HeapCell, len(src))
	CopyHeapCells(dst, src)
	return dst
}

// ToSlot converts a heap cell to a frame Slot for AOT bridge interop.
// desc is the field/element descriptor. Caller must use a pointer to avoid
// copying the atomic value inside HeapCell.
//
// Long and double cannot be represented in a single Slot (Slot.num is int32).
// Callers must use the long/double-specific bridge path or direct typed access.
func (c *HeapCell) ToSlot(desc string) Slot {
	switch desc[0] {
	case 'L', '[':
		return Slot{ref: c.GetRef()}
	case 'F':
		return Slot{num: int32(math.Float32bits(c.GetFloat()))}
	case 'J':
		panic("HeapCell.ToSlot: long values cannot be represented as a single Slot; use direct cell access")
	case 'D':
		panic("HeapCell.ToSlot: double values cannot be represented as a single Slot; use direct cell access")
	default: // Z, B, C, S, I
		return Slot{num: c.GetInt()}
	}
}
