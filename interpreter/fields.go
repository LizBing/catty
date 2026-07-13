package interpreter

import (
	"math"

	"catty/rtda"
)

// loadFieldValue reads a field/array heap cell and pushes its value onto the
// operand stack, interpreting the bits by descriptor. storage is a slice of
// rtda.HeapCell per ADR-0030; long/double occupy exactly one cell.
func loadFieldValue(frame *rtda.Frame, storage []rtda.HeapCell, cellID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		frame.PushInt(storage[cellID].GetInt())
	case 'F':
		frame.PushFloat(storage[cellID].GetFloat())
	case 'J':
		frame.PushLong(storage[cellID].GetLong())
	case 'D':
		frame.PushDouble(storage[cellID].GetDouble())
	default: // 'L' object or '[' array
		frame.PushRef(storage[cellID].GetRef())
	}
}

// storeFieldValue pops a value off the operand stack and writes it into a heap
// cell, encoding by descriptor. storage is a slice of rtda.HeapCell per ADR-0030.
func storeFieldValue(frame *rtda.Frame, storage []rtda.HeapCell, cellID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		storage[cellID].SetInt(frame.PopInt())
	case 'F':
		storage[cellID].SetFloat(frame.PopFloat())
	case 'J':
		storage[cellID].SetLong(frame.PopLong())
	case 'D':
		storage[cellID].SetDouble(frame.PopDouble())
	default:
		storage[cellID].SetRef(frame.PopRef())
	}
}

// bits32 extracts the float bits from a uint32 for use with HeapCell.SetFloat.
func floatbits(v float32) uint32 { return math.Float32bits(v) }
