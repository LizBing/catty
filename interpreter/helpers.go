package interpreter

import (
	"math"
	"strings"

	"catty/opcode"
	"catty/rtda"
)

// ---------- typed local-variable load/store for the indexed forms ----------

func loadLocal(frame *rtda.Frame, op opcode.Opcode, idx int) {
	switch op {
	case opcode.Iload:
		frame.PushInt(frame.GetInt(idx))
	case opcode.Lload:
		frame.PushLong(frame.GetLong(idx))
	case opcode.Fload:
		frame.PushFloat(frame.GetFloat(idx))
	case opcode.Dload:
		frame.PushDouble(frame.GetDouble(idx))
	case opcode.Aload:
		frame.PushRef(frame.GetRef(idx))
	}
}

func storeLocal(frame *rtda.Frame, op opcode.Opcode, idx int) {
	switch op {
	case opcode.Istore:
		frame.SetInt(idx, frame.PopInt())
	case opcode.Lstore:
		frame.SetLong(idx, frame.PopLong())
	case opcode.Fstore:
		frame.SetFloat(idx, frame.PopFloat())
	case opcode.Dstore:
		frame.SetDouble(idx, frame.PopDouble())
	case opcode.Astore:
		frame.SetRef(idx, frame.PopRef())
	}
}

// ---------- two-slot array elements (long[] / double[]) ----------

// readTwoSlots reads a category-2 element at array index i (each element spans
// two slots in long[]/double[]).
func readTwoSlots(arr *rtda.Object, i int) int64 {
	f := arr.Fields()
	base := i * 2
	high := uint32(f[base].Num())
	low := uint32(f[base+1].Num())
	return int64(high)<<32 | int64(low)
}

func writeTwoSlots(arr *rtda.Object, i int, v int64) {
	f := arr.Fields()
	base := i * 2
	f[base].SetNum(int32(uint64(v) >> 32)) // high
	f[base+1].SetNum(int32(v))             // low
}

// ---------- float/double bit helpers (named wrappers keep the switch tidy) ----------

func float32bits(f float32) uint32        { return math.Float32bits(f) }
func float32frombits(u uint32) float32    { return math.Float32frombits(u) }
func float64bits(f float64) uint64        { return math.Float64bits(f) }
func float64frombits(u uint64) float64    { return math.Float64frombits(u) }

// remF / remF64 implement Java's floating-point %, which is fmod semantics
// (result sign follows the dividend) — exactly Go's math.Mod.
func remF(a, b float32) float32 { return float32(math.Mod(float64(a), float64(b))) }
func remF64(a, b float64) float64 { return math.Mod(a, b) }

// cmpFloat / cmpDouble implement fcmpl/fcmpg and dcmpl/dcmpg: NaN yields -1 for
// the "l" variant and +1 for the "g" variant; otherwise the usual ordering.
func cmpFloat(a, b float32, isCmpg bool) int32 {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		if isCmpg {
			return 1
		}
		return -1
	}
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

func cmpDouble(a, b float64, isCmpg bool) int32 {
	if math.IsNaN(a) || math.IsNaN(b) {
		if isCmpg {
			return 1
		}
		return -1
	}
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

// ---------- branches ----------

func branch(frame *rtda.Frame, opcodePc, offset int) {
	frame.SetPC(opcodePc + offset)
}

// condBranch reads the (always-present) branch offset and applies it if cond.
func condBranch(frame *rtda.Frame, opcodePc int, cond bool) {
	offset := frame.ReadInt16()
	if cond {
		frame.SetPC(opcodePc + int(offset))
	}
}

// tableSwitch implements the tableswitch instruction: a dense jump table over
// contiguous integer keys [low, high], with a default offset.
func tableSwitch(frame *rtda.Frame, opcodePc int) {
	base := opcodePc + 1
	frame.SetPC(base + padTo4(base))
	defaultOff := frame.ReadInt32()
	low := frame.ReadInt32()
	high := frame.ReadInt32()
	key := frame.PopInt()

	off := defaultOff
	if key >= low && key <= high {
		frame.SetPC(frame.PC() + int(key-low)*4) // seek to the matching offset
		off = frame.ReadInt32()
	}
	frame.SetPC(opcodePc + int(off))
}

// lookupSwitch implements lookupswitch: a sparse, sorted (match, offset) pair
// table with a default. Linear scan is correct; pairs are sorted so a real JVM
// binary-searches — fine to defer for MVP workloads.
func lookupSwitch(frame *rtda.Frame, opcodePc int) {
	base := opcodePc + 1
	frame.SetPC(base + padTo4(base))
	defaultOff := frame.ReadInt32()
	npairs := frame.ReadInt32()
	key := frame.PopInt()

	off := defaultOff
	for range npairs {
		match := frame.ReadInt32()
		offset := frame.ReadInt32()
		if match == key {
			off = offset
			break
		}
	}
	frame.SetPC(opcodePc + int(off))
}

// padTo4 returns the 0-3 bytes of padding needed so an address after a switch
// opcode aligns to a 4-byte boundary relative to the method's code start.
func padTo4(base int) int { return (4 - base%4) % 4 }

// ---------- returns ----------

func returnInt(frame *rtda.Frame, thread *rtda.Thread) {
	v := frame.PopInt()
	thread.PopFrame()
	if !thread.IsStackEmpty() {
		thread.CurrentFrame().PushInt(v)
	}
}

func returnRef(frame *rtda.Frame, thread *rtda.Thread) {
	v := frame.PopRef()
	thread.PopFrame()
	if !thread.IsStackEmpty() {
		thread.CurrentFrame().PushRef(v)
	}
}

func returnLong(frame *rtda.Frame, thread *rtda.Thread) {
	v := frame.PopLong()
	thread.PopFrame()
	if !thread.IsStackEmpty() {
		thread.CurrentFrame().PushLong(v)
	}
}

func returnFloat(frame *rtda.Frame, thread *rtda.Thread) {
	v := frame.PopFloat()
	thread.PopFrame()
	if !thread.IsStackEmpty() {
		thread.CurrentFrame().PushFloat(v)
	}
}

func returnDouble(frame *rtda.Frame, thread *rtda.Thread) {
	v := frame.PopDouble()
	thread.PopFrame()
	if !thread.IsStackEmpty() {
		thread.CurrentFrame().PushDouble(v)
	}
}

// ---------- array creation ----------

// newPrimitiveArray builds an array of a primitive type (newarray). The atype
// operand encodes the element type per JVMS §6.5.newarray.
func newPrimitiveArray(thread *rtda.Thread, atype byte, length int) *rtda.Object {
	class := thread.Loader().LoadClass(primitiveArrayName(atype))
	return rtda.NewArray(class, length)
}

func primitiveArrayName(atype byte) string {
	switch atype {
	case 4:
		return "[Z"
	case 5:
		return "[C"
	case 6:
		return "[F"
	case 7:
		return "[D"
	case 8:
		return "[B"
	case 9:
		return "[S"
	case 10:
		return "[I"
	case 11:
		return "[J"
	}
	panic("catty: invalid newarray atype")
}

// newRefArray builds an array of references (anewarray), including arrays of
// arrays. elemName is the element's internal class name ("java/lang/String" or
// an array descriptor like "[I").
func newRefArray(thread *rtda.Thread, elemName string, length int) *rtda.Object {
	var arrName string
	if strings.HasPrefix(elemName, "[") {
		arrName = "[" + elemName
	} else {
		arrName = "[L" + elemName + ";"
	}
	class := thread.Loader().LoadClass(arrName)
	return rtda.NewArray(class, length)
}
