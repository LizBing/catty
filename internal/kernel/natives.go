package kernel

import (
	"fmt"
	"io"
	"strconv"
	"sync/atomic"
	"time"
)

// ---- small arg helpers -----------------------------------------------------

func argI(args []Value, i int) int32   { return args[i].(int32) }
func argL(args []Value, i int) int64   { return args[i].(int64) }
func argF(args []Value, i int) float32 { return args[i].(float32) }
func argD(args []Value, i int) float64 { return args[i].(float64) }

func boolV(b bool) Value {
	if b {
		return int32(1)
	}
	return int32(0)
}

func writeOut(ctx *CallContext, s string) {
	_, _ = io.WriteString(ctx.K.Stdout(), s)
}

// ---- java/lang/Object -------------------------------------------------------

func natObjectInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return nil, nil
}

func natObjectHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return int32(recv.(*Instance).IdentityHash()), nil
}

// RefIdentical is Java `==` on references.
func RefIdentical(a, b Value) bool {
	switch a.(type) {
	case *Instance, *ArrayObj, *JString:
		return a == b
	default:
		if a == nil || b == nil {
			return a == nil && b == nil
		}
		return a == b
	}
}

func natObjectEquals(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(RefIdentical(recv, args[0])), nil
}

func natObjectToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in := recv.(*Instance)
	s := dotted(in.Class.Name) + "@" + strconv.FormatUint(uint64(in.IdentityHash()), 16)
	return ctx.K.MakeJStringFromGo(s), nil
}

// ---- java/lang/String (backed by *JString) ----------------------------------

func natStringLength(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return int32(len(recv.(*JString).Chars)), nil
}

func natStringCharAt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	chars := recv.(*JString).Chars
	i := argI(args, 0)
	if i < 0 || int(i) >= len(chars) {
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException",
			"String index out of range: "+FormatInt(i))
	}
	return chars[i], nil
}

func natStringEquals(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, okA := recv.(*JString)
	b, okB := args[0].(*JString)
	if !okA || !okB {
		return boolV(false), nil
	}
	if len(a.Chars) != len(b.Chars) {
		return boolV(false), nil
	}
	for i := range a.Chars {
		if a.Chars[i] != b.Chars[i] {
			return boolV(false), nil
		}
	}
	return boolV(true), nil
}

func natStringHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	h := int32(0)
	for _, c := range recv.(*JString).Chars {
		h = 31*h + int32(c)
	}
	return h, nil
}

func natStringToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return recv, nil
}

// ---- java/lang/StringBuilder (payload *sbBuf) --------------------------------

type sbBuf struct{ buf []uint16 }

func sbOf(recv Value) *sbBuf { return recv.(*Instance).Payload.(*sbBuf) }

func natStringBuilderInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &sbBuf{}
	return nil, nil
}

func natSBAppendString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	if s, ok := args[0].(*JString); ok && s != nil {
		b.buf = append(b.buf, s.Chars...)
	} else {
		b.buf = append(b.buf, utf16Encode([]rune("null"))...)
	}
	return recv, nil
}

func natSBAppendInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	b.buf = append(b.buf, utf16Encode([]rune(FormatInt(argI(args, 0))))...)
	return recv, nil
}

func natSBAppendLong(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	b.buf = append(b.buf, utf16Encode([]rune(FormatLong(argL(args, 0))))...)
	return recv, nil
}

func natSBAppendChar(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	b.buf = append(b.buf, uint16(argI(args, 0)))
	return recv, nil
}

func natSBAppendBool(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	t := "false"
	if argI(args, 0) != 0 {
		t = "true"
	}
	b.buf = append(b.buf, utf16Encode([]rune(t))...)
	return recv, nil
}

func natSBToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	chars := make([]uint16, len(b.buf))
	copy(chars, b.buf)
	return ctx.K.MakeJString(chars), nil
}

// ---- java/io/PrintStream (payload io.Writer) ----------------------------------

func natPrintlnString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if s, ok := args[0].(*JString); ok && s != nil {
		writeOut(ctx, s.String()+"\n")
	} else {
		writeOut(ctx, "null\n")
	}
	return nil, nil
}

func natPrintlnInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, FormatInt(argI(args, 0))+"\n")
	return nil, nil
}

func natPrintlnLong(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, FormatLong(argL(args, 0))+"\n")
	return nil, nil
}

func natPrintlnChar(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, string(rune(argI(args, 0)))+"\n")
	return nil, nil
}

func natPrintlnBool(ctx *CallContext, recv Value, args []Value) (Value, error) {
	t := "false"
	if argI(args, 0) != 0 {
		t = "true"
	}
	writeOut(ctx, t+"\n")
	return nil, nil
}

func natPrintlnObject(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, ctx.Stringify(args[0])+"\n")
	return nil, nil
}

func natPrintlnVoid(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, "\n")
	return nil, nil
}

func natPrintString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if s, ok := args[0].(*JString); ok && s != nil {
		writeOut(ctx, s.String())
	} else {
		writeOut(ctx, "null")
	}
	return nil, nil
}

func natPrintObject(ctx *CallContext, recv Value, args []Value) (Value, error) {
	writeOut(ctx, ctx.Stringify(args[0]))
	return nil, nil
}

// ---- java/lang/System ---------------------------------------------------------

func systemOutInit(k *Kernel, c *Class) error {
	psCls := k.lookupClass("java/io/PrintStream")
	if psCls == nil {
		return fmt.Errorf("PrintStream missing during System init")
	}
	ps, err := k.NewInstance(psCls)
	if err != nil {
		return err
	}
	ps.Payload = io.Writer(k.Stdout())
	f := c.FindField("out", "Ljava/io/PrintStream;")
	if f == nil || !f.Static {
		return fmt.Errorf("System.out field missing")
	}
	c.Statics[f.StaticSlot] = ps
	return nil
}

func natSystemArraycopy(ctx *CallContext, recv Value, args []Value) (Value, error) {
	src := args[0].(*ArrayObj)
	srcPos := argI(args, 1)
	dst := args[2].(*ArrayObj)
	dstPos := argI(args, 3)
	length := argI(args, 4)
	check := func(arr *ArrayObj, pos, n int32, what string) error {
		if pos < 0 || n < 0 || int64(pos)+int64(n) > int64(len(arr.Elems)) {
			return ctx.Throw("java/lang/IndexOutOfBoundsException",
				fmt.Sprintf("arraycopy: %s index out of range", what))
		}
		return nil
	}
	if err := check(src, srcPos, length, "source"); err != nil {
		return nil, err
	}
	if err := check(dst, dstPos, length, "destination"); err != nil {
		return nil, err
	}
	tmp := make([]Value, length)
	copy(tmp, src.Elems[srcPos:int32(srcPos)+length])
	copy(dst.Elems[dstPos:], tmp)
	return nil, nil
}

func natCurrentTimeMillis(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return time.Now().UnixMilli(), nil
}

func natNanoTime(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return time.Now().UnixNano(), nil
}

func natIdentityHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if args[0] == nil {
		return int32(0), nil
	}
	switch o := args[0].(type) {
	case *Instance:
		return int32(o.IdentityHash()), nil
	case *ArrayObj:
		return int32(o.IdentityHash()), nil
	case *JString:
		return int32(o.IdentityHash()), nil
	}
	return int32(0), nil
}

// ---- java/lang/Integer ----------------------------------------------------------

const intValueSlotFallback = 0

func integerValue(recv Value) int32 {
	in := recv.(*Instance)
	if v, ok := in.Fields[0].(int32); ok {
		return v
	}
	return 0
}

func natIntegerInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Fields[0] = argI(args, 0)
	return nil, nil
}

func natIntegerIntValue(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return integerValue(recv), nil
}

func natIntegerHashCode(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return integerValue(recv), nil
}

func natIntegerEquals(ctx *CallContext, recv Value, args []Value) (Value, error) {
	o, ok := args[0].(*Instance)
	if !ok || o.Class.Name != "java/lang/Integer" {
		return boolV(false), nil
	}
	return boolV(integerValue(recv) == integerValue(o)), nil
}

func natIntegerToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.MakeJStringFromGo(FormatInt(integerValue(recv))), nil
}

func natStaticIntegerValueOf(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.IntegerOf(argI(args, 0)), nil
}

func natStaticIntegerToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.MakeJStringFromGo(FormatInt(argI(args, 0))), nil
}

// ---- java.lang.Throwable family ---------------------------------------------

func setDetailMessage(ctx *CallContext, recv Value, msg string) {
	in := recv.(*Instance)
	f := in.Class.FindField("detailMessage", "Ljava/lang/String;")
	if f != nil && msg != "" {
		in.Fields[f.Slot] = ctx.K.MakeJStringFromGo(msg)
	}
}

func natThrowableInitVoid(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return nil, nil
}

func natThrowableInitString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if s, ok := args[0].(*JString); ok && s != nil {
		setDetailMessage(ctx, recv, s.String())
	}
	return nil, nil
}

func natThrowableGetMessage(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return recv.(*Instance).fieldByName("detailMessage"), nil
}

func natThrowableToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in := recv.(*Instance)
	name := dotted(in.Class.Name)
	var msg string
	if m := in.fieldByName("detailMessage"); m != nil {
		msg = ": " + m.(*JString).String()
	}
	return ctx.K.MakeJStringFromGo(name + msg), nil
}

// ---- java/util/ArrayList (payload *alBuf) --------------------------------------

type alBuf struct{ data []Value }

func alOf(recv Value) *alBuf { return recv.(*Instance).Payload.(*alBuf) }

// elementMatches implements M0 contains(): identity, plus Integer value
// equality (real equals() dispatch lands with the M1 reflection work).
func elementMatches(elem, target Value) bool {
	if RefIdentical(elem, target) {
		return true
	}
	ei, okE := IntValueOf(elem)
	ti, okT := IntValueOf(target)
	return okE && okT && ei == ti
}

func natArrayListInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &alBuf{}
	return nil, nil
}

func natArrayListAdd(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := alOf(recv)
	b.data = append(b.data, args[0])
	return boolV(true), nil
}

func natArrayListGet(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := alOf(recv)
	i := argI(args, 0)
	if i < 0 || int(i) >= len(b.data) {
		return nil, ctx.Throw("java/lang/IndexOutOfBoundsException",
			fmt.Sprintf("Index: %d, Size: %d", i, len(b.data)))
	}
	return b.data[i], nil
}

func natArrayListSet(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := alOf(recv)
	i := argI(args, 0)
	if i < 0 || int(i) >= len(b.data) {
		return nil, ctx.Throw("java/lang/IndexOutOfBoundsException",
			fmt.Sprintf("Index: %d, Size: %d", i, len(b.data)))
	}
	old := b.data[i]
	b.data[i] = args[1]
	return old, nil
}

func natArrayListSize(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return int32(len(alOf(recv).data)), nil
}

func natArrayListIsEmpty(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(len(alOf(recv).data) == 0), nil
}

func natArrayListContains(ctx *CallContext, recv Value, args []Value) (Value, error) {
	for _, e := range alOf(recv).data {
		if elementMatches(e, args[0]) {
			return boolV(true), nil
		}
	}
	return boolV(false), nil
}


// ---- java/lang/Object monitor operations (wait/notify) ----------------------

// heapHeader extracts the common header of any heap value.
func heapHeader(v Value) *Header {
	switch o := v.(type) {
	case *Instance:
		return &o.Header
	case *ArrayObj:
		return &o.Header
	case *JString:
		return &o.Header
	default:
		panic(fmt.Sprintf("kernel: %T is not a heap object", v))
	}
}

func natObjectWaitMillis(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if ctx.Owner == nil {
		return nil, fmt.Errorf("wait called without thread context")
	}
	if ctx.K.Threads.ClearInterrupted(ctx.ownerKey()) {
		return nil, ctx.Throw("java/lang/InterruptedException", "flag set before wait")
	}
	out := heapHeader(recv).Monitor().Wait(ctx.K.Threads, ctx.ownerKey(), argL(args, 0))
	if out == waitGotInterrupt {
		return nil, ctx.Throw("java/lang/InterruptedException", "wait interrupted")
	}
	return nil, nil
}

func natObjectNotify(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if _, err := heapHeader(recv).Monitor().Notify(ctx.ownerKey()); err != nil {
		return nil, ctx.Throw("java/lang/IllegalMonitorStateException", "notify without ownership")
	}
	return nil, nil
}

func natObjectNotifyAll(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if _, err := heapHeader(recv).Monitor().NotifyAll(ctx.ownerKey()); err != nil {
		return nil, ctx.Throw("java/lang/IllegalMonitorStateException", "notifyAll without ownership")
	}
	return nil, nil
}

// ---- java/lang/Thread --------------------------------------------------------

func threadOf(recv Value) (*JThread, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("receiver is not a heap object")
	}
	j, ok := in.Payload.(*JThread)
	if !ok || j == nil {
		return nil, fmt.Errorf("receiver is not a java.lang.Thread (no record)")
	}
	return j, nil
}

func natThreadInitVoid(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return initJavaThread(ctx, recv, "")
}

func natThreadInitName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	name := "Thread"
	if s, ok := args[0].(*JString); ok && s != nil {
		name = s.String()
	}
	return initJavaThread(ctx, recv, name)
}

var threadNameCtr atomic.Uint64

func initJavaThread(ctx *CallContext, recv Value, name string) (Value, error) {
	in := recv.(*Instance)
	if name == "" {
		name = fmt.Sprintf("Thread-%d", threadNameCtr.Add(1))
	}
	j := ctx.K.Threads.NewRecord(in, name)
	in.Payload = j
	f := in.Class.FindField("name", "Ljava/lang/String;")
	if f != nil {
		in.Fields[f.Slot] = ctx.K.MakeJStringFromGo(name)
	}
	return nil, nil
}

func natThreadStart(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	if !j.state.CompareAndSwap(int32(threadNew), int32(threadRunning)) {
		return nil, ctx.Throw("java/lang/IllegalThreadStateException", "start on running/finished thread")
	}
	key := ctx.K.MintKey()
	j.Key = key
	ctx.K.Threads.attach(j)
	j.alive.Store(true)
	if ctx.K.SpawnJavaThread == nil {
		return nil, fmt.Errorf("no SpawnJavaThread hook installed")
	}
	ctx.K.SpawnJavaThread(j)
	return nil, nil
}

func natThreadRunDefault(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return nil, nil // default run(); subclasses override via dispatch
}

func natThreadJoinMillis(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	reached, interrupted := ctx.K.Threads.Join(ctx.ownerKey(), j, time.Duration(argL(args, 0))*time.Millisecond)
	if interrupted {
		return nil, ctx.Throw("java/lang/InterruptedException", "join interrupted")
	}
	_ = reached
	return nil, nil
}

func natThreadIsAlive(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	return boolV(j.IsAlive()), nil
}

func natThreadSetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, _ := threadOf(recv)
	if s, ok := args[0].(*JString); ok && s != nil {
		j.Name = s.String()
	}
	in := recv.(*Instance)
	if f := in.Class.FindField("name", "Ljava/lang/String;"); f != nil {
		in.Fields[f.Slot] = args[0]
	}
	return nil, nil
}

func natThreadGetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in := recv.(*Instance)
	if f := in.Class.FindField("name", "Ljava/lang/String;"); f != nil {
		if s, ok := in.Fields[f.Slot].(*JString); ok {
			return s, nil
		}
	}
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	return ctx.K.MakeJStringFromGo(j.Name), nil
}

func natThreadInterrupt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	ctx.K.Threads.InterruptByKey(j)
	return nil, nil
}

func natThreadIsInterrupted(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j, err := threadOf(recv)
	if err != nil {
		return nil, err
	}
	return boolV(j.interrupted.Load()), nil
}

func natStaticCurrentThread(ctx *CallContext, recv Value, args []Value) (Value, error) {
	j := ctx.K.Threads.ByKey(ctx.ownerKey())
	if j == nil {
		return nil, fmt.Errorf("current thread has no record")
	}
	return j.Obj, nil
}

func natStaticSleep(ctx *CallContext, recv Value, args []Value) (Value, error) {
	key := ctx.ownerKey()
	if key == 0 {
		return nil, fmt.Errorf("Thread.sleep without thread context")
	}
	if ctx.K.Threads.ClearInterrupted(key) {
		return nil, ctx.Throw("java/lang/InterruptedException", "sleep interrupted before park")
	}
	if !ctx.K.Threads.Sleep(key, time.Duration(argL(args, 0))*time.Millisecond) {
		return nil, ctx.Throw("java/lang/InterruptedException", "sleep interrupted")
	}
	return nil, nil
}

func natStaticInterruptedFlag(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return boolV(ctx.K.Threads.ClearInterrupted(ctx.ownerKey())), nil
}

// natThreadJoinForever backs join() with no timeout.
func natThreadJoinForever(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return natThreadJoinMillis(ctx, recv, []Value{int64(0)})
}
