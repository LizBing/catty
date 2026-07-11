package rtda

import "math"

// Frame is one JVM stack frame: the local variable array, the operand stack,
// the method being executed, and the current program counter (JVMS §2.5.2).
//
// Fields are unexported and accessed through methods so the interpreter stays
// decoupled from the slot layout; the methods are small enough for the compiler
// to inline in the dispatch loop.
type Frame struct {
	lower    *Frame // caller frame; threads keep frames as a linked stack
	thread   *Thread
	method   *Method
	code     []byte // cached method.code for the dispatch loop (nil if native)
	locals   []Slot
	stack    []Slot
	stackTop int  // index of the next free slot = current operand-stack size
	pc       int  // index into code of the next instruction
}

func NewFrame(thread *Thread, method *Method) *Frame {
	return &Frame{
		thread: thread,
		method: method,
		code:   method.code,
		locals: make([]Slot, method.maxLocals),
		stack:  make([]Slot, method.maxStack),
	}
}

func (f *Frame) Thread() *Thread { return f.thread }
func (f *Frame) Method() *Method { return f.method }
func (f *Frame) Code() []byte    { return f.method.code }

func (f *Frame) PC() int       { return f.pc }
func (f *Frame) SetPC(pc int)  { f.pc = pc }
func (f *Frame) NextPC() int   { return f.pc } // pc already advanced past the opcode

// Operand decoders. The dispatch loop reads opcodes and their inline operands
// straight from code via these; they advance pc past the bytes consumed.
func (f *Frame) ReadUint8() uint8 {
	b := f.code[f.pc]
	f.pc++
	return b
}

func (f *Frame) ReadUint16() uint16 {
	v := uint16(f.code[f.pc])<<8 | uint16(f.code[f.pc+1])
	f.pc += 2
	return v
}

// ReadInt16 reads a signed branch offset, used by goto and the conditional jumps.
func (f *Frame) ReadInt16() int16 {
	return int16(f.ReadUint16())
}

// ReadInt32 reads a big-endian signed 4-byte operand (switch offsets/keys).
func (f *Frame) ReadInt32() int32 {
	v := int32(f.code[f.pc])<<24 | int32(f.code[f.pc+1])<<16 |
		int32(f.code[f.pc+2])<<8 | int32(f.code[f.pc+3])
	f.pc += 4
	return v
}

// --- Operand stack: category-1 (int / float / ref / returnAddress) ---

func (f *Frame) PushInt(val int32) {
	f.stack[f.stackTop].num = val
	f.stackTop++
}

func (f *Frame) PopInt() int32 {
	f.stackTop--
	return f.stack[f.stackTop].num
}

func (f *Frame) PushFloat(val float32) {
	f.stack[f.stackTop].num = int32(math.Float32bits(val))
	f.stackTop++
}

func (f *Frame) PopFloat() float32 {
	f.stackTop--
	return math.Float32frombits(uint32(f.stack[f.stackTop].num))
}

func (f *Frame) PushRef(ref *Object) {
	f.stack[f.stackTop].ref = ref
	f.stackTop++
}

// PopRef nils the underlying slot so the GC can reclaim the object: leaving a
// stale reference in a freed slot would keep it alive until the slot is reused.
func (f *Frame) PopRef() *Object {
	f.stackTop--
	r := f.stack[f.stackTop].ref
	f.stack[f.stackTop].ref = nil
	return r
}

// PeekRef returns the reference `below` slots beneath the stack top (0 = top)
// without popping. Used by invokevirtual to read `this` for dynamic dispatch
// before the args are copied into the callee frame.
func (f *Frame) PeekRef(below int) *Object {
	return f.stack[f.stackTop-1-below].ref
}

// PushRefSlot / PopRefSlot move a slot by value without type interpretation —
// used by dup, pop (on refs), and any instruction that copies a slot opaquely.
func (f *Frame) PushSlot(s Slot) {
	f.stack[f.stackTop] = s
	f.stackTop++
}

func (f *Frame) PopSlot() Slot {
	f.stackTop--
	return f.stack[f.stackTop]
}

// --- Operand stack: category-2 (long / double take two slots) ---

func (f *Frame) PushLong(val int64) {
	f.stack[f.stackTop].num = int32(uint64(val) >> 32) // high word
	f.stackTop++
	f.stack[f.stackTop].num = int32(val) // low word
	f.stackTop++
}

func (f *Frame) PopLong() int64 {
	f.stackTop--
	low := uint32(f.stack[f.stackTop].num)
	f.stackTop--
	high := uint32(f.stack[f.stackTop].num)
	return int64(high)<<32 | int64(low)
}

func (f *Frame) PushDouble(val float64) {
	f.PushLong(int64(math.Float64bits(val)))
}

func (f *Frame) PopDouble() float64 {
	return math.Float64frombits(uint64(f.PopLong()))
}

// --- Local variables ---

func (f *Frame) SetInt(index int, val int32)    { f.locals[index].num = val }
func (f *Frame) GetInt(index int) int32         { return f.locals[index].num }
func (f *Frame) SetFloat(index int, val float32) { f.locals[index].num = int32(math.Float32bits(val)) }
func (f *Frame) GetFloat(index int) float32     { return math.Float32frombits(uint32(f.locals[index].num)) }
func (f *Frame) SetRef(index int, ref *Object)  { f.locals[index].ref = ref }
func (f *Frame) GetRef(index int) *Object       { return f.locals[index].ref }

func (f *Frame) SetLong(index int, val int64) {
	f.locals[index].num = int32(uint64(val) >> 32)
	f.locals[index+1].num = int32(val)
}

func (f *Frame) GetLong(index int) int64 {
	high := uint32(f.locals[index].num)
	low := uint32(f.locals[index+1].num)
	return int64(high)<<32 | int64(low)
}

func (f *Frame) SetDouble(index int, val float64) { f.SetLong(index, int64(math.Float64bits(val))) }
func (f *Frame) GetDouble(index int) float64      { return math.Float64frombits(uint64(f.GetLong(index))) }

// SetSlot / Slot copy a raw local slot by value, used when handing arguments to
// a callee frame (the slot's category is the callee's concern).
func (f *Frame) SetSlot(index int, s Slot) { f.locals[index] = s }

// --- indexed operand-stack access (for the IR executor) ---
//
// The lowering pass eliminates the operand stack into slot-indexed virtual
// registers. These accessors let the IR executor read/write those slots by the
// index the lowering computed, instead of via Push/Pop. SetStackTop seeds the
// stack pointer from the lowering's known depth so shared Push/Pop-based helpers
// (invokes, fields, …) still work when the executor calls them.

func (f *Frame) SetStackTop(n int) { f.stackTop = n }

func (f *Frame) StackSlotNum(i int) int32    { return f.stack[i].num }
func (f *Frame) StackSlotRef(i int) *Object  { return f.stack[i].ref }
func (f *Frame) SetStackSlotNum(i int, v int32)   { f.stack[i].num = v }
func (f *Frame) SetStackSlotRef(i int, r *Object) { f.stack[i].ref = r }

// CopyStackSlot copies a whole slot (num + ref) by value, for the IR executor's
// stack-shuffle instructions (dup, swap, …) which move category-unknown values.
func (f *Frame) CopyStackSlot(dst, src int) { f.stack[dst] = f.stack[src] }
