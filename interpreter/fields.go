package interpreter

import (
	"math"

	"catty/rtda"
)

// loadInstanceField reads an instance field from obj and pushes its typed value
// onto the operand stack. cellID is the field's cell index (0-based, per ADR-0030).
func loadInstanceField(frame *rtda.Frame, obj *rtda.Object, cellID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		frame.PushInt(obj.GetIntCell(int(cellID)))
	case 'F':
		frame.PushFloat(obj.GetFloatCell(int(cellID)))
	case 'J':
		frame.PushLong(obj.GetLongCell(int(cellID)))
	case 'D':
		frame.PushDouble(obj.GetDoubleCell(int(cellID)))
	default: // 'L' object or '[' array
		frame.PushRef(obj.GetRefCell(int(cellID)))
	}
}

// storeInstanceField pops a typed value from the operand stack and writes it into
// an instance field. cellID is the field's cell index.
func storeInstanceField(frame *rtda.Frame, obj *rtda.Object, cellID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		obj.SetIntCell(int(cellID), frame.PopInt())
	case 'F':
		obj.SetFloatCell(int(cellID), frame.PopFloat())
	case 'J':
		obj.SetLongCell(int(cellID), frame.PopLong())
	case 'D':
		obj.SetDoubleCell(int(cellID), frame.PopDouble())
	default:
		obj.SetRefCell(int(cellID), frame.PopRef())
	}
}

// loadStaticField reads a static field from class and pushes its typed value
// onto the operand stack. slotID is the field's static cell index.
func loadStaticField(frame *rtda.Frame, class *rtda.Class, slotID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		frame.PushInt(class.GetStaticInt(slotID))
	case 'F':
		frame.PushFloat(class.GetStaticFloat(slotID))
	case 'J':
		frame.PushLong(class.GetStaticLong(slotID))
	case 'D':
		frame.PushDouble(class.GetStaticDouble(slotID))
	default:
		frame.PushRef(class.GetStaticRef(slotID))
	}
}

// storeStaticField pops a typed value from the operand stack and writes it into
// a static field. slotID is the field's static cell index.
func storeStaticField(frame *rtda.Frame, class *rtda.Class, slotID uint, desc string) {
	switch desc[0] {
	case 'Z', 'B', 'C', 'S', 'I':
		class.SetStaticInt(slotID, frame.PopInt())
	case 'F':
		class.SetStaticFloat(slotID, frame.PopFloat())
	case 'J':
		class.SetStaticLong(slotID, frame.PopLong())
	case 'D':
		class.SetStaticDouble(slotID, frame.PopDouble())
	default:
		class.SetStaticRef(slotID, frame.PopRef())
	}
}

// floatbits extracts the float bits from a uint32 for use with HeapCell.SetFloat.
func floatbits(v float32) uint32 { return math.Float32bits(v) }
