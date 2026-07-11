package interpreter

import (
	"math"

	"catty/rtda"
)

// loadFieldValue reads a field/array cell from a slot slice and pushes its value
// onto the operand stack, interpreting the bits by descriptor.
func loadFieldValue(frame *rtda.Frame, storage []rtda.Slot, slotID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		frame.PushInt(storage[slotID].Num())
	case 'F':
		frame.PushFloat(math.Float32frombits(uint32(storage[slotID].Num())))
	case 'J':
		high := uint32(storage[slotID].Num())
		low := uint32(storage[slotID+1].Num())
		frame.PushLong(int64(high)<<32 | int64(low))
	case 'D':
		high := uint32(storage[slotID].Num())
		low := uint32(storage[slotID+1].Num())
		frame.PushDouble(math.Float64frombits(uint64(high)<<32 | uint64(low)))
	default: // 'L' object or '[' array
		frame.PushRef(storage[slotID].Ref())
	}
}

// storeFieldValue pops a value off the operand stack and writes it into a slot
// slice, encoding by descriptor (long/double span two slots).
func storeFieldValue(frame *rtda.Frame, storage []rtda.Slot, slotID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		storage[slotID].SetNum(frame.PopInt())
	case 'F':
		storage[slotID].SetNum(int32(math.Float32bits(frame.PopFloat())))
	case 'J':
		v := frame.PopLong()
		storage[slotID].SetNum(int32(uint64(v) >> 32)) // high
		storage[slotID+1].SetNum(int32(v))             // low
	case 'D':
		bits := math.Float64bits(frame.PopDouble())
		storage[slotID].SetNum(int32(bits >> 32)) // high
		storage[slotID+1].SetNum(int32(bits))     // low
	default:
		storage[slotID].SetRef(frame.PopRef())
	}
}
