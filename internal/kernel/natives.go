package kernel

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
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
		switch v := m.(type) {
		case *JString:
			msg = ": " + v.String()
		case *Instance:
			if js, jsErr := AsJString(v); jsErr == nil {
				msg = ": " + js.String()
			}
		default:
			msg = fmt.Sprintf(": %v", v)
		}
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

// ---- java/lang/Class + Object.getClass ----------------------------------------

func natObjectGetClass(ctx *CallContext, recv Value, args []Value) (Value, error) {
	cls := heapHeader(recv).Class
	return ctx.K.ClassObjectOf(cls)
}

func natClassGetName(ctx *CallContext, recv Value, args []Value) (Value, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("getClass receiver payload missing")
	}
	c, ok := in.Payload.(*Class)
	if !ok {
		return nil, fmt.Errorf("Class instance lacks class payload")
	}
	if strings.HasPrefix(c.Name, "[") {
		return ctx.K.MakeJStringFromGo(c.Name), nil // arrays keep descriptor form
	}
	return ctx.K.MakeJStringFromGo(dotted(c.Name)), nil
}

// ---- java/lang/String byte conversions + search --------------------------------

// elemsToBytes extracts [off,off+len) of a byte array as Go bytes.
func elemsToBytes(arr *ArrayObj, off, length int32) ([]byte, error) {
	if off < 0 || length < 0 || int64(off)+int64(length) > int64(len(arr.Elems)) {
		return nil, fmt.Errorf("bad range off=%d len=%d size=%d", off, length, len(arr.Elems))
	}
	out := make([]byte, length)
	for i := 0; i < int(length); i++ {
		out[i] = byte(arr.Elems[int(off)+i].(int32))
	}
	return out, nil
}

func bytesToElems(b []byte) []Value {
	out := make([]Value, len(b))
	for i, c := range b {
		out[i] = int32(c)
	}
	return out
}

func natStringGetBytes(ctx *CallContext, recv Value, args []Value) (Value, error) {
	chars := recv.(*JString).Chars
	rs := utf16Decode(chars)
	b := make([]byte, 0, len(rs))
	for _, r := range rs {
		b = appendRune(b, r)
	}
	arr, err := ctx.K.NewArray("[B", len(b))
	if err != nil {
		return nil, err
	}
	copy(arr.Elems, bytesToElems(b))
	return arr, nil
}

func appendRune(b []byte, r rune) []byte {
	switch {
	case r < 0x80:
		return append(b, byte(r))
	case r < 0x800:
		return append(b, byte(0xC0|r>>6), byte(0x80|r&0x3F))
	case r < 0x10000:
		return append(b, byte(0xE0|r>>12), byte(0x80|(r>>6)&0x3F), byte(0x80|r&0x3F))
	default:
		r -= 0x10000
		hi, lo := rune(0xD800+(r>>10)), rune(0xDC00+(r&0x3FF))
		return append(b,
			byte(0xE0|hi>>12), byte(0x80|(hi>>6)&0x3F), byte(0x80|hi&0x3F),
			byte(0xE0|lo>>12), byte(0x80|(lo>>6)&0x3F), byte(0x80|lo&0x3F))
	}
}

func natStringInitBytesRange(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	off := argI(args, 1)
	length := argI(args, 2)
	b, err := elemsToBytes(arr, off, length)
	if err != nil {
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException", err.Error())
	}
	js := ctx.K.MakeJString(utf16Encode(decodeUTF8Runes(b)))
	recv.(*JString).Chars = js.Chars
	return nil, nil
}

func decodeUTF8Runes(b []byte) []rune {
	var rs []rune
	for i := 0; i < len(b); {
		c := b[i]
		switch {
		case c < 0x80:
			rs = append(rs, rune(c))
			i++
		case c>>5 == 0b110 && i+1 < len(b):
			rs = append(rs, rune(c&0x1F)<<6|rune(b[i+1]&0x3F))
			i += 2
		case c>>4 == 0b1110 && i+2 < len(b):
			rs = append(rs, rune(c&0x0F)<<12|rune(b[i+1]&0x3F)<<6|rune(b[i+2]&0x3F))
			i += 3
		default:
			rs = append(rs, 0xFFFD)
			i++
		}
	}
	return rs
}

func natStringIndexOf(ctx *CallContext, recv Value, args []Value) (Value, error) {
	hay := recv.(*JString).Chars
	pat, ok := args[0].(*JString)
	if !ok || pat == nil {
		return int32(-1), nil
	}
	if len(pat.Chars) == 0 || len(pat.Chars) > len(hay) {
		if len(pat.Chars) == 0 {
			return 0, nil
		}
		return int32(-1), nil
	}
	for i := 0; i+len(pat.Chars) <= len(hay); i++ {
		match := true
		for j := range pat.Chars {
			if hay[i+j] != pat.Chars[j] {
				match = false
				break
			}
		}
		if match {
			return int32(i), nil
		}
	}
	return int32(-1), nil
}

func natStaticIntegerParseInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	s, ok := args[0].(*JString)
	if !ok || s == nil {
		return nil, ctx.Throw("java/lang/NumberFormatException", "null")
	}
	text := s.String()
	neg := false
	body := text
	if strings.HasPrefix(body, "-") {
		neg = true
		body = body[1:]
	} else if strings.HasPrefix(body, "+") {
		body = body[1:]
	}
	if body == "" {
		return nil, ctx.Throw("java/lang/NumberFormatException", "for input string: \""+text+"\"")
	}
	var v int64
	for _, ch := range []byte(body) {
		if ch < '0' || ch > '9' {
			return nil, ctx.Throw("java/lang/NumberFormatException",
				"for input string: \""+text+"\"")
		}
		v = v*10 + int64(ch-'0')
		if v > 1<<31 {
			v = 1 << 31 // clamp; overflow wrap handled below
		}
	}
	if neg {
		v = -v
	}
	if v < math.MinInt32 || v > math.MaxInt32 {
		return nil, ctx.Throw("java/lang/NumberFormatException",
			"for input string: \""+text+"\"")
	}
	return int32(v), nil
}

// natStreamWriteB backs write([B)V.
func natStreamWriteB(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	return natStreamWriteBII(ctx, recv, []Value{arr, int32(0), int32(len(arr.Elems))})
}

func natSBAppendObject(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	if len(args) > 0 && args[0] != nil {
		s, err := valueToDisplayString(ctx, args[0])
		if err != nil {
			return nil, err
		}
		b.buf = append(b.buf, utf16Encode([]rune(s))...)
	} else {
		b.buf = append(b.buf, utf16Encode([]rune("null"))...)
	}
	return recv, nil
}

// valueToDisplayString calls toString() on an object value.
func valueToDisplayString(ctx *CallContext, v Value) (string, error) {
	if js, ok := v.(*JString); ok {
		return js.String(), nil
	}
	in, ok := v.(*Instance)
	if !ok {
		return fmt.Sprintf("%v", v), nil
	}
	m, err := ctx.K.ResolveMethod(in.Class, "toString", "()Ljava/lang/String;")
	if err != nil {
		return fmt.Sprintf("%s@%x", in.Class.Name, 0), nil
	}
	res, err := ctx.K.InvokeAs(ctx.Owner, m, v, nil)
	if err != nil {
		return "", err
	}
	if js, ok := res.(*JString); ok {
		return js.String(), nil
	}
	return "", nil
}

func natStaticNanoTime(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return time.Now().UnixNano(), nil
}

func natStaticCurrentTimeMillis(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return time.Now().UnixMilli(), nil
}

func natSBInitString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := &sbBuf{}
	if len(args) > 0 && args[0] != nil {
		js, err := AsJString(args[0])
		if err == nil {
			b.buf = append(b.buf, js.Chars...)
		}
	}
	recv.(*Instance).Payload = b
	return nil, nil
}

var procStart = time.Now()

// tickMillis reports milliseconds since process start as int32 — a
// cat1-free clock for benchmark loops in emitted code (long-return
// canonical-depth unification pending).
func natStaticTickMillis(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return int32(time.Since(procStart).Milliseconds()), nil
}

func natSBLength(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return int32(len(sbOf(recv).buf)), nil
}

func natStringInitChars(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	rs := make([]rune, 0, len(arr.Elems))
	for _, e := range arr.Elems {
		if c, ok := e.(int32); ok {
			rs = append(rs, rune(c))
		}
	}
	js := ctx.K.MakeJString(utf16Encode(rs))
	recv.(*JString).Chars = js.Chars
	return nil, nil
}

func natStringInitCharsRange(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	off, length := int(argI(args, 1)), int(argI(args, 2))
	rs := make([]rune, 0, length)
	for i := off; i < off+length && i < len(arr.Elems); i++ {
		if c, ok := arr.Elems[i].(int32); ok {
			rs = append(rs, rune(c))
		}
	}
	js := ctx.K.MakeJString(utf16Encode(rs))
	recv.(*JString).Chars = js.Chars
	return nil, nil
}

func natStringGetChars(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js := recv.(*JString)
	srcBegin := int(argI(args, 0))
	srcEnd := int(argI(args, 1))
	dst := args[2].(*ArrayObj)
	dstOff := int(argI(args, 3))
	if srcBegin < 0 || srcEnd > len(js.Chars) || srcBegin > srcEnd {
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException",
			fmt.Sprintf("begin %d, end %d, length %d", srcBegin, srcEnd, len(js.Chars)))
	}
	for i := srcBegin; i < srcEnd; i++ {
		dst.Elems[dstOff+i-srcBegin] = int32(js.Chars[i])
	}
	return nil, nil
}

func natSBAppendChars(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	arr := args[0].(*ArrayObj)
	off, length := int(argI(args, 1)), int(argI(args, 2))
	for i := off; i < off+length && i < len(arr.Elems); i++ {
		if c, ok := arr.Elems[i].(int32); ok {
			b.buf = append(b.buf, uint16(c))
		}
	}
	return recv, nil
}

func natSBSetLength(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	n := argI(args, 0)
	switch {
	case n < 0:
		return nil, ctx.Throw("java/lang/StringIndexOutOfBoundsException", fmt.Sprintf("length %d", n))
	case int(n) < len(b.buf):
		b.buf = b.buf[:n]
	default:
		for i := len(b.buf); i < int(n); i++ {
			b.buf = append(b.buf, 0)
		}
	}
	return nil, nil
}

func natStaticIntegerParseIntRadix(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	radix := argI(args, 1)
	v, err := strconv.ParseInt(js.String(), int(radix), 32)
	if err != nil {
		return nil, ctx.Throw("java/lang/NumberFormatException", js.String())
	}
	return int32(v), nil
}

func natStaticIntegerToStringRadix(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := argI(args, 0)
	radix := int(argI(args, 1))
	return ctx.K.InternGo(strconv.FormatInt(int64(v), radix)), nil
}

func natSBAppendDouble(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := sbOf(recv)
	d := args[0].(float64)
	s := strconv.FormatFloat(d, 'g', -1, 64)
	for _, r := range s {
		b.buf = append(b.buf, uint16(r))
	}
	return recv, nil
}
