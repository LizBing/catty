// Package vm implements the Catty bytecode interpreter (M0 execution
// engine). It executes classfile bytecode against the kernel object
// model. Performance is explicitly not a goal in M0 (P-0001); the AOT
// path (ADR-0003) is the performance vehicle.
package vm

import (
	"errors"
	"fmt"
	"math"
	"os"

	"catty/internal/classfile"
	"catty/internal/kernel"
)

// Thread is an execution context: one goroutine of Java execution.
type Thread struct {
	K         *kernel.Kernel
	id        uint64
	jobj      *kernel.Instance // bound java.lang.Thread (nil for primordial v0)
	depth     int              // interpreted-frame depth (SOE metering)
	maxDepth  int
	initStack []string // class names with <clinit> in progress (JVMS §5.5)

	jstack []kernel.JavaFrame // Java call stack for backfill; touched only by the owning goroutine
}

// --- kernel.FrameTracker ------------------------------------------------------

// PushJavaFrame appends one frame (stack-backfill bookkeeping). Only the
// owning goroutine executes Java on this Thread, so no lock is required;
// cross-thread readers go through JavaFrames' snapshot.
func (t *Thread) PushJavaFrame(f kernel.JavaFrame) { t.jstack = append(t.jstack, f) }

// PopJavaFrame drops the top frame. A mismatched pop is an engine bug and
// panics loudly instead of silently corrupting traces.
func (t *Thread) PopJavaFrame() {
	n := len(t.jstack)
	if n == 0 {
		panic("vm: PopJavaFrame on empty stack")
	}
	t.jstack[n-1] = kernel.JavaFrame{} // drop reference for GC hygiene
	t.jstack = t.jstack[:n-1]
}

// JavaFrames snapshots the current stack, outermost first.
func (t *Thread) JavaFrames() []kernel.JavaFrame {
	cp := make([]kernel.JavaFrame, len(t.jstack))
	copy(cp, t.jstack)
	return cp
}

// SetTopJavaLine implements kernel.FrameTracker: records the source line
// the top frame is executing (call sites and line-segment entries). This
// is what makes non-leaf trace frames print their call-site line and the
// leaf print its throw-site line (U3).
func (t *Thread) SetTopJavaLine(line int32) {
	if n := len(t.jstack); n > 0 {
		t.jstack[n-1].Line = line
	}
}

// lineAt resolves the source line for pc in m's LineNumberTable.
// Returns 0 when absent (Unknown Source).
func lineAt(m *kernel.Method, pc int) int32 {
	if m == nil || m.Code == nil || len(m.Code.LineNumbers) == 0 {
		return 0
	}
	lns := m.Code.LineNumbers
	line := int32(0)
	for i := range lns {
		if int(lns[i].StartPc) <= pc {
			line = int32(lns[i].Line)
		} else {
			break
		}
	}
	return line
}

// New creates the primordial thread for a kernel and installs the VM's
// bridges (Invoker, spawn hook, main-thread record).
func New(k *kernel.Kernel) *Thread {
	id := k.MintKey()
	t := &Thread{K: k, id: id}
	k.InstallInvoker(t)

	// Bind java.lang.Thread for main and install goroutine spawning.
	if cls, ok := k.ClassByName("java/lang/Thread"); ok {
		if obj, err := k.NewInstance(cls); err == nil {
			name := "main"
			j := k.Threads.Register(id, obj, name)
			obj.Payload = j
			if f := cls.FindField("name", "Ljava/lang/String;"); f != nil {
				obj.Fields[f.Slot] = k.InternGo(name)
			}
			t.jobj = obj
			k.Threads.SetMain(j)
		}
	}
	k.SpawnJavaThread = func(j *kernel.JThread) { go runJavaThread(k, j) }
	k.UncaughtHandler = defaultUncaughtHandler(k)
	return t
}

// OwnerKey implements kernel.OwnerKey.
func (t *Thread) OwnerKey() uint64 { return t.id }

// InvokeInterpreted implements kernel.Invoker with SOE depth metering.
// The meter is shared with the emitted path (genrt.invokeChecked) so one
// budget spans both engines like a real JVM stack.
func (t *Thread) InvokeInterpreted(m *kernel.Method, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
	if err := t.EnterFrame(); err != nil {
		return nil, t.throwNamed("java/lang/StackOverflowError", "")
	}
	defer t.ExitFrame()
	return t.exec(m, recv, args)
}

// EnterFrame implements kernel.FrameMeter (lock-free per-thread depth).
func (t *Thread) EnterFrame() error {
	if t.maxDepth == 0 { // primordial thread lazy-init
		t.maxDepth = t.K.MaxFrames()
	}
	t.depth++
	if t.depth > t.maxDepth {
		return kernel.ErrStackOverflow
	}
	return nil
}

// ExitFrame implements kernel.FrameMeter.
func (t *Thread) ExitFrame() { t.depth-- }

// --- kernel.InitTracker -----------------------------------------------------

func (t *Thread) IsInitializing(name string) bool {
	for _, n := range t.initStack {
		if n == name {
			return true
		}
	}
	return false
}

func (t *Thread) BeginInit(name string) { t.initStack = append(t.initStack, name) }
func (t *Thread) EndInit(name string) {
	if n := len(t.initStack); n > 0 {
		t.initStack = t.initStack[:n-1]
	}
}

// EnsureInitialized drives <clinit> for a class through the kernel.
func (t *Thread) EnsureInitialized(c *kernel.Class) error {
	return t.K.EnsureInitialized(t, c)
}

// InvokeInterpreted implements kernel.Invoker.

// Call invokes any method (native or interpreted) with full dispatch,
// attributing monitor operations to this thread.
func (t *Thread) Call(m *kernel.Method, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
	return t.K.InvokeAs(t, m, recv, args)
}

// --- frame -------------------------------------------------------------------

type frame struct {
	m      *kernel.Method
	locals []kernel.Value
	stack  []kernel.Value
	sp     int
	code   []byte
	pc     int
}

// operand readers (advance f.pc)
func (f *frame) u1() uint8 { v := f.code[f.pc]; f.pc++; return v }
func (f *frame) u2() uint16 {
	v := uint16(f.code[f.pc])<<8 | uint16(f.code[f.pc+1])
	f.pc += 2
	return v
}
func (f *frame) s2() int32 { return int32(int16(f.u2())) }
func (f *frame) s4() int32 {
	v := uint32(f.code[f.pc])<<24 | uint32(f.code[f.pc+1])<<16 |
		uint32(f.code[f.pc+2])<<8 | uint32(f.code[f.pc+3])
	f.pc += 4
	return int32(v)
}
func (f *frame) s4pad() int32 { // switch padding to 4-byte alignment (relative to method start)
	for f.pc%4 != 0 {
		f.pc++
	}
	return f.s4()
}

// branch sets pc relative to the branch opcode address.
func (f *frame) branch(from int, offset int32) { f.pc = from + int(offset) }

// push is category-aware (cat2 consumes two raw slots). It also
// normalizes stray numeric widths (e.g. uint16 leaked from a native
// helper) into the canonical int32 so type switches downstream are
// total.
func (f *frame) push(v kernel.Value) {
	if u, ok := v.(uint16); ok {
		v = int32(u)
	}
	f.stack[f.sp] = v
	f.sp++
	if kernel.IsCat2(v) {
		f.stack[f.sp] = kernel.TopSentinel
		f.sp++
	}
}

func (f *frame) push1(v kernel.Value) { f.stack[f.sp] = v; f.sp++ }

func (f *frame) popRaw() kernel.Value {
	f.sp--
	v := f.stack[f.sp]
	f.stack[f.sp] = nil
	return v
}

func (f *frame) pop1() kernel.Value { return f.popRaw() }

// pop2s pops a category-2 value laid out as [value][sentinel].
func (f *frame) pop2s() kernel.Value {
	v := f.stack[f.sp-2]
	f.sp -= 2
	f.stack[f.sp] = nil
	f.stack[f.sp+1] = nil
	return v
}

func (f *frame) popI() int32          { return f.pop1().(int32) }
func (f *frame) popL() int64          { return f.pop2s().(int64) }
func (f *frame) popF() float32        { return f.pop1().(float32) }
func (f *frame) popD() float64        { return f.pop2s().(float64) }
func (f *frame) popRef() kernel.Value { return f.pop1() }

// asThrown converts an error into *kernel.Thrown when it is one.
func asThrown(err error) (*kernel.Thrown, bool) {
	var th *kernel.Thrown
	if errors.As(err, &th) {
		return th, true
	}
	return nil, false
}

// exec runs one interpreted method to completion.
func (t *Thread) exec(m *kernel.Method, recv kernel.Value, args []kernel.Value) (ret kernel.Value, err error) {
	code := m.Code
	f := &frame{
		m:      m,
		locals: make([]kernel.Value, code.MaxLocals),
		stack:  make([]kernel.Value, int(code.MaxStack)+2),
		code:   code.Code,
	}

	// Locals: receiver first, then arguments per descriptor slots.
	idx := 0
	if !m.Static() {
		f.locals[0] = recv
		idx = 1
	}
	argDescs, _, err := kernel.ParseMethodDesc(m.Desc)
	if err != nil {
		return nil, fmt.Errorf("bad descriptor %s: %w", m.Desc, err)
	}
	if len(args) != len(argDescs) {
		return nil, fmt.Errorf("call %s%s%s: %d args, want %d",
			m.Holder.Name, m.Name, m.Desc, len(args), len(argDescs))
	}
	for i, ad := range argDescs {
		f.locals[idx] = args[i]
		idx += kernel.SlotCount(ad)
	}

	// throwOrHandle routes a thrown error into a matching handler in this
	// frame; returns handled=true when control continues here.
	throwOrHandle := func(thrown error, at int) bool {
		th, ok := asThrown(thrown)
		if !ok {
			return false
		}
		for _, h := range code.Handlers {
			if at < int(h.StartPc) || at >= int(h.EndPc) {
				continue
			}
			if h.CatchType != 0 {
				cc, err := t.resolveClassIdx(m.Holder, h.CatchType)
				if err != nil {
					return false
				}
				if !t.K.IsInstance(th.Obj, cc) {
					continue
				}
			}
			for i := 0; i < f.sp; i++ {
				f.stack[i] = nil
			}
			f.sp = 1
			f.stack[0] = th.Obj
			f.pc = int(h.HandlerPc)
			return true
		}
		return false
	}

	trace := os.Getenv("CATTY_TRACE") != ""
	for {
		if f.pc >= len(f.code) {
			return nil, fmt.Errorf("pc overrun in %s.%s%s", m.Holder.Name, m.Name, m.Desc)
		}
		faultPc := f.pc
		op := f.code[f.pc]
		if trace {
			fmt.Fprintf(os.Stderr, "[trace] %s.%s pc=%d op=%#x sp=%d\n",
				m.Holder.Name, m.Name, faultPc, op, f.sp)
		}
		f.pc++

		switch op {

		// ---- constants ----
		case opNop:
		case opAconstNull:
			f.push(nil)
		case opIconstM1:
			f.push(int32(-1))
		case opIconst0:
			f.push(int32(0))
		case opIconst1:
			f.push(int32(1))
		case opIconst2:
			f.push(int32(2))
		case opIconst3:
			f.push(int32(3))
		case opIconst4:
			f.push(int32(4))
		case opIconst5:
			f.push(int32(5))
		case opLconst0:
			f.push(int64(0))
		case opLconst1:
			f.push(int64(1))
		case opFconst0:
			f.push(float32(0))
		case opFconst1:
			f.push(float32(1))
		case opFconst2:
			f.push(float32(2))
		case opDconst0:
			f.push(float64(0))
		case opDconst1:
			f.push(float64(1))
		case opBipush:
			f.push(int32(int8(f.u1())))
		case opSipush:
			f.push(int32(int16(f.u2())))
		case opLdc, opLdcW:
			var ci uint16
			if op == opLdc {
				ci = uint16(f.u1())
			} else {
				ci = f.u2()
			}
			v, err := t.ldc(m.Holder, ci)
			if err != nil {
				return nil, err
			}
			f.push(v)
		case opLdc2W:
			e, err := m.CF.Entry(f.u2())
			if err != nil {
				return nil, err
			}
			switch e.Tag {
			case classfile.CLong:
				f.push(e.LongVal)
			case classfile.CDouble:
				f.push(e.DoubleVal)
			default:
				return nil, fmt.Errorf("ldc2_w on %s", e.Tag)
			}

		// ---- loads ----
		case opIload:
			f.push(f.locals[f.u1()])
		case opLload, opDload:
			f.push(f.locals[f.u1()])
		case opFload:
			f.push(f.locals[f.u1()])
		case opAload:
			f.push(f.locals[f.u1()])
		case opIload0, opIload1, opIload2, opIload3:
			f.push(f.locals[int(op-opIload0)])
		case opFload0, opFload1, opFload2, opFload3:
			f.push(f.locals[int(op-opFload0)])
		case opLload0, opLload1, opLload2, opLload3:
			f.push(f.locals[int(op-opLload0)])
		case opDload0, opDload1, opDload2, opDload3:
			f.push(f.locals[int(op-opDload0)])
		case opAload0, opAload1, opAload2, opAload3:
			f.push(f.locals[int(op-opAload0)])

		// ---- array loads ----
		case opIaload, opBaload, opCaload, opSaload, opLaload, opFaload, opDaload, opAaload:
			v, th, err := t.arrayLoad(f, faultPc)
			if err != nil {
				return nil, err
			}
			if th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			switch op {
			case opCaload:
				v = v.(int32) & 0xFFFF
			case opBaload:
				v = int32(int8(v.(int32)))
			}
			f.push(v)

		// ---- stores ----
		case opIstore, opFstore, opAstore:
			f.locals[f.u1()] = f.pop1()
		case opLstore, opDstore:
			f.locals[f.u1()] = f.pop2s()
		case opIstore0, opIstore1, opIstore2, opIstore3:
			f.locals[int(op-opIstore0)] = f.pop1()
		case opFstore0, opFstore1, opFstore2, opFstore3:
			f.locals[int(op-opFstore0)] = f.pop1()
		case opAstore0, opAstore1, opAstore2, opAstore3:
			f.locals[int(op-opAstore0)] = f.pop1()
		case opLstore0, opLstore1, opLstore2, opLstore3:
			f.locals[int(op-opLstore0)] = f.pop2s()
		case opDstore0, opDstore1, opDstore2, opDstore3:
			f.locals[int(op-opDstore0)] = f.pop2s()

		// ---- array stores ----
		case opIastore, opBastore, opCastore, opSastore:
			v := f.popI()
			switch op {
			case opBastore:
				v = int32(int8(v))
			case opCastore:
				v &= 0xFFFF
			case opSastore:
				v = int32(int16(v))
			}
			if th := t.arrayStore(f, faultPc, v); th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opLastore:
			v := f.popL()
			if th := t.arrayStore(f, faultPc, v); th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opFastore:
			v := f.popF()
			if th := t.arrayStore(f, faultPc, v); th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opDastore:
			v := f.popD()
			if th := t.arrayStore(f, faultPc, v); th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opAastore:
			v := f.popRef()
			if th := t.arrayStore(f, faultPc, v); th != nil {
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}

		// ---- stack ops (raw slots, per JVMS) ----
		case opPop:
			f.pop1()
		case opPop2:
			f.sp -= 2
		case opDup:
			f.stack[f.sp] = f.stack[f.sp-1]
			f.sp++
		case opDupX1:
			t1 := f.stack[f.sp-1]
			f.stack[f.sp] = t1
			f.stack[f.sp-1] = f.stack[f.sp-2]
			f.stack[f.sp-2] = t1
			f.sp++
		case opDupX2:
			t1 := f.stack[f.sp-1]
			f.stack[f.sp] = t1
			f.stack[f.sp-1] = f.stack[f.sp-2]
			f.stack[f.sp-2] = f.stack[f.sp-3]
			f.stack[f.sp-3] = t1
			f.sp++
		case opDup2:
			f.stack[f.sp] = f.stack[f.sp-2]
			f.stack[f.sp+1] = f.stack[f.sp-1]
			f.sp += 2
		case opDup2X1:
			n1 := f.stack[f.sp-1]
			n2 := f.stack[f.sp-2]
			f.stack[f.sp+1] = n1
			f.stack[f.sp] = n2
			f.stack[f.sp-1] = f.stack[f.sp-3]
			f.stack[f.sp-2] = n2
			f.stack[f.sp-3] = n1
			f.sp += 2
		case opDup2X2:
			a := f.stack[f.sp-4]
			b := f.stack[f.sp-3]
			c := f.stack[f.sp-2]
			d := f.stack[f.sp-1]
			f.stack[f.sp+1] = d
			f.stack[f.sp] = c
			f.stack[f.sp-1] = b
			f.stack[f.sp-2] = a
			f.stack[f.sp-3] = c
			f.stack[f.sp-4] = d
			f.sp += 2
		case opSwap:
			f.stack[f.sp-1], f.stack[f.sp-2] = f.stack[f.sp-2], f.stack[f.sp-1]

		// ---- int arithmetic ----
		case opIadd:
			b, a := f.popI(), f.popI()
			f.push(a + b)
		case opIsub:
			b, a := f.popI(), f.popI()
			f.push(a - b)
		case opImul:
			b, a := f.popI(), f.popI()
			f.push(a * b)
		case opIdiv:
			b, a := f.popI(), f.popI()
			if b == 0 {
				t.SetTopJavaLine(lineAt(m, faultPc))
				if th := t.throwNamed("java/lang/ArithmeticException", "/ by zero"); th != nil {
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
			}
			f.push(a / b)
		case opIrem:
			b, a := f.popI(), f.popI()
			if b == 0 {
				t.SetTopJavaLine(lineAt(m, faultPc))
				if th := t.throwNamed("java/lang/ArithmeticException", "/ by zero"); th != nil {
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
			}
			f.push(a % b)
		case opIneg:
			f.push(-f.popI())
		case opIshl:
			s, a := f.popI(), f.popI()
			f.push(a << (uint32(s) & 31))
		case opIshr:
			s, a := f.popI(), f.popI()
			f.push(a >> (uint32(s) & 31))
		case opIushr:
			s, a := f.popI(), f.popI()
			f.push(int32(uint32(a) >> (uint32(s) & 31)))
		case opIand:
			b, a := f.popI(), f.popI()
			f.push(a & b)
		case opIor:
			b, a := f.popI(), f.popI()
			f.push(a | b)
		case opIxor:
			b, a := f.popI(), f.popI()
			f.push(a ^ b)

		// ---- long arithmetic ----
		case opLadd:
			b, a := f.popL(), f.popL()
			f.push(a + b)
		case opLsub:
			b, a := f.popL(), f.popL()
			f.push(a - b)
		case opLmul:
			b, a := f.popL(), f.popL()
			f.push(a * b)
		case opLdiv:
			b, a := f.popL(), f.popL()
			if b == 0 {
				t.SetTopJavaLine(lineAt(m, faultPc))
				if th := t.throwNamed("java/lang/ArithmeticException", "/ by zero"); th != nil {
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
			}
			f.push(a / b) // MIN/-1 wraps to MIN in Go, matching JVMS
		case opLrem:
			b, a := f.popL(), f.popL()
			if b == 0 {
				t.SetTopJavaLine(lineAt(m, faultPc))
				if th := t.throwNamed("java/lang/ArithmeticException", "/ by zero"); th != nil {
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
			}
			f.push(a % b)
		case opLneg:
			f.push(-f.popL())
		case opLshl:
			s, a := f.popI(), f.popL()
			f.push(a << (uint32(s) & 63))
		case opLshr:
			s, a := f.popI(), f.popL()
			f.push(a >> (uint32(s) & 63))
		case opLushr:
			s, a := f.popI(), f.popL()
			f.push(int64(uint64(a) >> (uint32(s) & 63)))
		case opLand:
			b, a := f.popL(), f.popL()
			f.push(a & b)
		case opLor:
			b, a := f.popL(), f.popL()
			f.push(a | b)
		case opLxor:
			b, a := f.popL(), f.popL()
			f.push(a ^ b)

		// ---- float/double arithmetic ----
		case opFadd:
			b, a := f.popF(), f.popF()
			f.push(a + b)
		case opFsub:
			b, a := f.popF(), f.popF()
			f.push(a - b)
		case opFmul:
			b, a := f.popF(), f.popF()
			f.push(a * b)
		case opFdiv:
			b, a := f.popF(), f.popF()
			f.push(a / b)
		case opFrem:
			b, a := f.popF(), f.popF()
			f.push(float32(math.Mod(float64(a), float64(b))))
		case opFneg:
			f.push(-f.popF())
		case opDadd:
			b, a := f.popD(), f.popD()
			f.push(a + b)
		case opDsub:
			b, a := f.popD(), f.popD()
			f.push(a - b)
		case opDmul:
			b, a := f.popD(), f.popD()
			f.push(a * b)
		case opDdiv:
			b, a := f.popD(), f.popD()
			f.push(a / b)
		case opDrem:
			b, a := f.popD(), f.popD()
			f.push(math.Mod(a, b))
		case opDneg:
			f.push(-f.popD())

		case opIinc:
			// Non-wide iinc: index u1 + const one-byte signed (JVMS §6.5).
			n := int(f.u1())
			c := int32(int8(f.u1()))
			f.locals[n] = f.locals[n].(int32) + c

		// ---- conversions ----
		case opI2l:
			f.push(int64(f.popI()))
		case opI2f:
			f.push(float32(f.popI()))
		case opI2d:
			f.push(float64(f.popI()))
		case opL2i:
			f.push(int32(f.popL()))
		case opL2f:
			f.push(float32(f.popL()))
		case opL2d:
			f.push(float64(f.popL()))
		case opF2i:
			f.push(floatToI32(float64(f.popF())))
		case opF2l:
			f.push(floatToI64(float64(f.popF())))
		case opF2d:
			f.push(float64(f.popF()))
		case opD2i:
			f.push(floatToI32(f.popD()))
		case opD2l:
			f.push(floatToI64(f.popD()))
		case opD2f:
			f.push(float32(f.popD()))
		case opI2b:
			f.push(int32(int8(f.popI())))
		case opI2c:
			f.push(f.popI() & 0xFFFF)
		case opI2s:
			f.push(int32(int16(f.popI())))

		// ---- comparisons ----
		case opLcmp:
			b, a := f.popL(), f.popL()
			f.push(cmpI64(a, b))
		case opFcmpl, opFcmpg:
			b, a := f.popF(), f.popF()
			f.push(cmpFP(float64(a), float64(b), op == opFcmpg))
		case opDcmpl, opDcmpg:
			b, a := f.popD(), f.popD()
			f.push(cmpFP(a, b, op == opDcmpg))

		// ---- branches ----
		case opIfeq, opIfne, opIflt, opIfge, opIfgt, opIfle:
			v := f.popI()
			if condCmp1(op, v) {
				f.branch(faultPc, f.s2())
			} else {
				f.pc += 2
			}
		case opIfIcmpeq, opIfIcmpne, opIfIcmplt, opIfIcmpge, opIfIcmpgt, opIfIcmple:
			b, a := f.popI(), f.popI()
			if condCmp2(op, a, b) {
				f.branch(faultPc, f.s2())
			} else {
				f.pc += 2
			}
		case opIfAcmpeq, opIfAcmpne:
			b, a := f.popRef(), f.popRef()
			eq := kernel.RefIdentical(a, b)
			if (op == opIfAcmpeq) == eq {
				f.branch(faultPc, f.s2())
			} else {
				f.pc += 2
			}
		case opGoto:
			f.branch(faultPc, f.s2())
		case opGotow:
			off := int32(f.code[f.pc])<<24 | int32(f.code[f.pc+1])<<16 |
				int32(f.code[f.pc+2])<<8 | int32(f.code[f.pc+3])
			f.pc += 4
			f.branch(faultPc, off)
		case opIfnull, opIfnonnull:
			v := f.popRef()
			isNull := v == nil
			if (op == opIfnull) == isNull {
				f.branch(faultPc, f.s2())
			} else {
				f.pc += 2
			}
		case opTableswitch:
			def := f.s4pad()
			low := f.s4()
			high := f.s4()
			v := f.popI()
			off := def
			if v >= low && v <= high {
				// offsets start at current pc: skip (v-low) 4-byte entries
				f.pc += int(v-low) * 4
				off = f.s4()
			}
			f.branch(faultPc, off)
		case opLookupswitch:
			def := f.s4pad()
			n := int(f.s4())
			v := f.popI()
			off := def
			matched := false
			for i := 0; i < n; i++ {
				k := f.s4()
				o := f.s4()
				if k == v && !matched {
					off = o
					matched = true
				}
			}
			f.branch(faultPc, off)

		// ---- returns ----
		case opIreturn, opFreturn, opAreturn:
			return f.pop1(), nil
		case opLreturn, opDreturn:
			return f.pop2s(), nil
		case opReturn:
			return nil, nil

		// ---- fields ----
		case opGetstatic:
			cls, name, desc, _, err := m.CF.Ref(f.u2())
			if err != nil {
				return nil, err
			}
			owner, err := t.resolveClassNamed(cls)
			if err != nil {
				return nil, err
			}
			if err := t.EnsureInitialized(owner); err != nil {
				return nil, err
			}
			fd, err := t.K.ResolveField(owner, name, desc)
			if err != nil {
				return nil, err
			}
			f.push(owner.Statics[fd.StaticSlot])
		case opPutstatic:
			cls, name, desc, _, err := m.CF.Ref(f.u2())
			if err != nil {
				return nil, err
			}
			owner, err := t.resolveClassNamed(cls)
			if err != nil {
				return nil, err
			}
			if err := t.EnsureInitialized(owner); err != nil {
				return nil, err
			}
			fd, err := t.K.ResolveField(owner, name, desc)
			if err != nil {
				return nil, err
			}
			if kernel.SlotCount(desc) == 2 {
				owner.Statics[fd.StaticSlot] = f.pop2s()
			} else {
				owner.Statics[fd.StaticSlot] = f.pop1()
			}
		case opGetfield:
			_, name, desc, _, err := m.CF.Ref(f.u2())
			if err != nil {
				return nil, err
			}
			obj := f.popRef()
			if obj == nil {
				th := t.npe(f, "getfield on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			owner := objClass(obj)
			fd, err := t.K.ResolveField(owner, name, desc)
			if err != nil {
				return nil, err
			}
			f.push(objFields(obj)[fd.Slot])
		case opPutfield:
			_, name, desc, _, err := m.CF.Ref(f.u2())
			if err != nil {
				return nil, err
			}
			var v kernel.Value
			if kernel.SlotCount(desc) == 2 {
				v = f.pop2s()
			} else {
				v = f.pop1()
			}
			obj := f.popRef()
			if obj == nil {
				th := t.npe(f, "putfield on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			owner := objClass(obj)
			fd, err := t.K.ResolveField(owner, name, desc)
			if err != nil {
				return nil, err
			}
			objFields(obj)[fd.Slot] = v

		// ---- invocations ----
		case opInvokevirtual, opInvokespecial, opInvokestatic, opInvokeinterf:
			cls, name, desc, _, err := m.CF.Ref(f.u2())
			if err != nil {
				return nil, err
			}
			if op == opInvokeinterf {
				f.u1() // count
				f.u1() // zero byte
			}
			argDescs, _, err := kernel.ParseMethodDesc(desc)
			if err != nil {
				return nil, err
			}
			slots := 0
			for _, ad := range argDescs {
				slots += kernel.SlotCount(ad)
			}
			rawArgs := make([]kernel.Value, slots)
			for i := slots - 1; i >= 0; i-- {
				rawArgs[i] = f.pop1()
			}
			var recv kernel.Value
			if op != opInvokestatic {
				recv = f.pop1()
			}
			vals := make([]kernel.Value, 0, len(argDescs))
			for _, rv := range rawArgs {
				if !kernel.IsTop(rv) {
					vals = append(vals, rv)
				}
			}

			var target *kernel.Method
			switch op {
			case opInvokestatic:
				owner, err := t.resolveClassNamed(cls)
				if err != nil {
					return nil, err
				}
				if err := t.EnsureInitialized(owner); err != nil {
					return nil, err
				}
				if target, err = t.K.ResolveMethod(owner, name, desc); err != nil {
					return nil, err
				}
			case opInvokevirtual, opInvokeinterf:
				if recv == nil {
					th := t.npe(f, "invokevirtual on null")
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
				dyn := objClass(recv)
				if target, err = t.K.ResolveMethod(dyn, name, desc); err != nil {
					return nil, err
				}
			case opInvokespecial:
				// M0 simplification: resolve in the named class walking
				// supers (covers <init>, private, super calls; exact
				// receiver-class rules per JVMS §5.4.3.3 deferred).
				owner, err := t.resolveClassNamed(cls)
				if err != nil {
					return nil, err
				}
				if target, err = t.K.ResolveMethod(owner, name, desc); err != nil {
					return nil, err
				}
				if recv == nil {
					th := t.npe(f, "invokespecial on null")
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
			}
			if target.Native == nil && target.Code == nil {
				return nil, fmt.Errorf("method %s.%s%s has neither code nor native", cls, name, desc)
			}
			// Method-level synchronized (JVMS §2.11.10): instance methods
			// lock the receiver, static methods lock the Class object.
			var syncHdr *kernel.Header
			if target.Flags&classfile.AccSynchronized != 0 {
				if target.Static() {
					cobj, cerr := t.K.ClassObjectOf(target.Holder)
					if cerr != nil {
						return nil, cerr
					}
					syncHdr = &cobj.Header
				} else if recv != nil {
					syncHdr = objHeader(recv)
				}
				if syncHdr == nil {
					th := t.npe(f, "synchronized method on null")
					if throwOrHandle(th, faultPc) {
						continue
					}
					return nil, th
				}
				syncHdr.Monitor().Enter(t.OwnerKey())
			}
			t.SetTopJavaLine(lineAt(m, faultPc)) // caller frame records call-site line (U3)
			res, callErr := t.K.InvokeAs(t, target, recv, vals) // attribute to this thread
			if syncHdr != nil {
				if serr := syncHdr.Monitor().Exit(t.OwnerKey()); serr != nil {
					return nil, fmt.Errorf("sync method exit: %w", serr)
				}
			}
			if callErr != nil {
				if throwOrHandle(callErr, faultPc) {
					continue
				}
				return nil, callErr
			}
			if res != nil || returnsValue(desc) {
				f.push(res)
			}

		// ---- object/array creation & misc ----
		case opNew:
			cls, err := t.resolveClassIdx(m.Holder, f.u2())
			if err != nil {
				return nil, err
			}
			if err := t.EnsureInitialized(cls); err != nil {
				return nil, err
			}
			var obj kernel.Value
			if cls.Name == "java/lang/String" {
				obj = t.K.MakeJString(nil)
			} else {
				in, err := t.K.NewInstance(cls)
				if err != nil {
					return nil, err
				}
				obj = in
			}
			f.push(obj)
		case opNewarray:
			at := f.u1()
			n := f.popI()
			desc, ok := atypeDesc(at)
			if !ok {
				return nil, fmt.Errorf("bad newarray atype %d", at)
			}
			arr, err := t.K.NewArray(desc, int(n))
			if err != nil {
				if throwOrHandle(err, faultPc) {
					continue
				}
				return nil, err
			}
			f.push(arr)
		case opAnewarray:
			cls, err := t.resolveClassIdx(m.Holder, f.u2())
			if err != nil {
				return nil, err
			}
			n := f.popI()
			arr, err := t.K.NewArray(cls.Name, int(n))
			if err != nil {
				if throwOrHandle(err, faultPc) {
					continue
				}
				return nil, err
			}
			f.push(arr)
		case opMultianewarr:
			cls, err := t.resolveClassIdx(m.Holder, f.u2())
			if err != nil {
				return nil, err
			}
			dims := int(f.u1())
			dl := make([]int32, dims)
			for i := dims - 1; i >= 0; i-- {
				dl[i] = f.popI()
			}
			arr, err := t.buildMultiArray(cls.Name, dl)
			if err != nil {
				if throwOrHandle(err, faultPc) {
					continue
				}
				return nil, err
			}
			f.push(arr)
		case opArraylength:
			arr := f.popRef()
			if arr == nil {
				th := t.npe(f, "arraylength on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			f.push(int32(len(arr.(*kernel.ArrayObj).Elems)))
		case opAthrow:
			obj := f.popRef()
			if obj == nil {
				th := t.npe(f, "athrow on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			return nil, &kernel.Thrown{Obj: obj.(*kernel.Instance)}
		case opCheckcast:
			cls, err := t.resolveClassIdx(m.Holder, f.u2())
			if err != nil {
				return nil, err
			}
			obj := f.stack[f.sp-1]
			if obj != nil && !t.K.IsInstance(obj, cls) {
				t.SetTopJavaLine(lineAt(m, faultPc))
				th := t.throwNamed("java/lang/ClassCastException",
					fmt.Sprintf("class %s cannot be cast to class %s", objClassName(obj), dotted(cls.Name)))
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opInstanceof:
			cls, err := t.resolveClassIdx(m.Holder, f.u2())
			if err != nil {
				return nil, err
			}
			obj := f.popRef()
			if obj == nil {
				f.push(int32(0))
			} else {
				f.push(bool32v(t.K.IsInstance(obj, cls)))
			}
		case opMonitorenter:
			obj := f.popRef()
			if obj == nil {
				th := t.npe(f, "monitorenter on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			objHeader(obj).Monitor().Enter(t.OwnerKey())
		case opMonitorexit:
			obj := f.popRef()
			if obj == nil {
				th := t.npe(f, "monitorexit on null")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
			if err := objHeader(obj).Monitor().Exit(t.OwnerKey()); err != nil {
				t.SetTopJavaLine(lineAt(m, faultPc))
					th := t.throwNamed("java/lang/IllegalMonitorStateException",
					"monitorexit on unowned monitor")
				if throwOrHandle(th, faultPc) {
					continue
				}
				return nil, th
			}
		case opWide:
			wop := f.u1()
			ni := int(f.u2())
			switch wop {
			case opIload, opFload, opAload:
				f.push(f.locals[ni])
			case opLload, opDload:
				f.push(f.locals[ni])
			case opIstore, opFstore, opAstore:
				f.locals[ni] = f.pop1()
			case opLstore, opDstore:
				f.locals[ni] = f.pop2s()
			case opIinc:
				inc := f.s2()
				f.locals[ni] = f.locals[ni].(int32) + inc
			case opRet:
				return nil, fmt.Errorf("jsr/ret unsupported (v51+)")
			default:
				return nil, fmt.Errorf("bad wide opcode %#x", wop)
			}

		case opJsr, opRet:
			return nil, fmt.Errorf("jsr/ret unsupported (prohibited in class file v51+)")
		case opInvokedynamic:
			return nil, fmt.Errorf("invokedynamic unsupported in M0 (build-time desugaring per ADR-0002)")

		default:
			return nil, fmt.Errorf("unsupported opcode %#x at %s.%s%s pc=%d",
				op, m.Holder.Name, m.Name, m.Desc, faultPc)
		}
	}
}
