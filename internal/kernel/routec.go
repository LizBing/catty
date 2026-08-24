package kernel

import (
	"catty/internal/classfile"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// bootstrapRouteC closes the JDK-surface gaps surfaced by the first
// third-party library (minimal-json): Serializable marker, String reader/
// writer, Double wrapper, List/Iterator protocol, Arrays/Collections utils.

func bootstrapRouteC(k *Kernel) {
	// Serializable: pure marker.
	mustDefine(k, &ClassDef{
		Name:  "java/io/Serializable",
		Flags: classfile.AccPublic | classfile.AccInterface | classfile.AccAbstract,
	})

	// Writer abstract + StringWriter (payload-backed *strings.Builder).
	mustDefine(k, &ClassDef{
		Name:  "java/io/Writer",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccAbstract,
		Methods: []MethodDef{
			{Name: "write", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic | classfile.AccAbstract},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic | classfile.AccAbstract},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/StringWriter",
		Super: "java/io/Writer",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "()V", Flags: classfile.AccPublic, Native: natStringWriterInit},
			{Name: "write", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natStringWriterWrite},
			{Name: "write", Desc: "([CII)V", Flags: classfile.AccPublic, Native: natStringWriterWriteCII},
			{Name: "write", Desc: "(I)V", Flags: classfile.AccPublic, Native: natStringWriterWriteChar},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natStringWriterToString},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natObjectInit},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/io/StringReader",
		Super: "java/io/Reader",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(Ljava/lang/String;)V", Flags: classfile.AccPublic, Native: natStringReaderInit},
			{Name: "read", Desc: "([CII)I", Flags: classfile.AccPublic, Native: natStringReaderRead},
			{Name: "read", Desc: "()I", Flags: classfile.AccPublic, Native: natStringReaderRead1},
			{Name: "close", Desc: "()V", Flags: classfile.AccPublic, Native: natFileClose},
		},
	})

	// Double wrapper (cat2 field).
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Double",
		Super: "java/lang/Number",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Fields: []FieldDef{
			{Name: "value", Desc: "D", Flags: classfile.AccPrivate | classfile.AccFinal},
		},
		Methods: []MethodDef{
			{Name: "<init>", Desc: "(D)V", Flags: classfile.AccPublic, Native: natDoubleInit},
			{Name: "doubleValue", Desc: "()D", Flags: classfile.AccPublic, Native: natWrapperValue("D")},
			{Name: "isNaN", Desc: "(D)Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleIsNaN},
			{Name: "isInfinite", Desc: "(D)Z", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleIsInfinite},
			{Name: "parseDouble", Desc: "(Ljava/lang/String;)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleParse},
			{Name: "toString", Desc: "(D)Ljava/lang/String;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natDoubleToString},
			{Name: "valueOf", Desc: "(D)Ljava/lang/Double;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natWrapperValueOf("java/lang/Double")},
			{Name: "toString", Desc: "()Ljava/lang/String;", Flags: classfile.AccPublic, Native: natWrapperToString},
		},
	})

	// List interface + iterator protocol on ArrayList.
	mustDefine(k, &ClassDef{
		Name:   "java/util/ArrayList$Itr",
		Super:  "java/lang/Object",
		Ifaces: []string{"java/util/Iterator"},
		Flags:  classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "hasNext", Desc: "()Z", Flags: classfile.AccPublic, Native: natItrHasNext},
			{Name: "next", Desc: "()Ljava/lang/Object;", Flags: classfile.AccPublic, Native: natItrNext},
		},
	})

	// Arrays.asList / Collections.unmodifiableList → pragmatic aliases.
	mustDefine(k, &ClassDef{
		Name:  "java/util/Arrays",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "asList", Desc: "([Ljava/lang/Object;)Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysAsList},
			{Name: "fill", Desc: "([CC)V", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysFillChar},
			{Name: "fill", Desc: "([II)V", Flags: classfile.AccPublic | classfile.AccStatic, Native: natArraysFillInt},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/lang/Math",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic | classfile.AccFinal,
		Methods: []MethodDef{
			{Name: "max", Desc: "(II)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxInt},
			{Name: "min", Desc: "(II)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinInt},
			{Name: "max", Desc: "(JJ)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxLong},
			{Name: "min", Desc: "(JJ)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinLong},
			{Name: "max", Desc: "(DD)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMaxDouble},
			{Name: "min", Desc: "(DD)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathMinDouble},
			{Name: "abs", Desc: "(I)I", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsInt},
			{Name: "abs", Desc: "(J)J", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsLong},
			{Name: "abs", Desc: "(D)D", Flags: classfile.AccPublic | classfile.AccStatic, Native: natMathAbsDouble},
		},
	})
	mustDefine(k, &ClassDef{
		Name:  "java/util/Collections",
		Super: "java/lang/Object",
		Flags: classfile.AccPublic,
		Methods: []MethodDef{
			{Name: "unmodifiableList", Desc: "(Ljava/util/List;)Ljava/util/List;", Flags: classfile.AccPublic | classfile.AccStatic, Native: natCollectionsUnmodList},
		},
	})
}

func natStringWriterInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Payload = &strings.Builder{}
	return nil, nil
}

func natStringWriterWrite(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	b.WriteString(javaStr(args[0]))
	return nil, nil
}

func natStringWriterWriteCII(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	arr := args[0].(*ArrayObj)
	off, length := int(argI(args, 1)), int(argI(args, 2))
	for i := off; i < off+length && i < len(arr.Elems); i++ {
		if c, ok := arr.Elems[i].(int32); ok {
			b.WriteRune(rune(c))
		}
	}
	return nil, nil
}

func natStringWriterWriteChar(ctx *CallContext, recv Value, args []Value) (Value, error) {
	b := recv.(*Instance).Payload.(*strings.Builder)
	if c, ok := args[0].(int32); ok {
		b.WriteRune(rune(c))
	}
	return nil, nil
}

func natStringWriterToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(recv.(*Instance).Payload.(*strings.Builder).String()), nil
}

type strRdr struct {
	cur  int
	s    string
	done bool
}

func natStringReaderInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	recv.(*Instance).Payload = &strRdr{s: js.String()}
	return nil, nil
}

func (sr *strRdr) readOne() int {
	if sr.done {
		return -1
	}
	r := []rune(sr.s)
	if sr.cur >= len(r) {
		sr.done = true
		return -1
	}
	ch := r[sr.cur]
	sr.cur++
	return int(ch)
}

func natStringReaderRead1(ctx *CallContext, recv Value, args []Value) (Value, error) {
	sr := recv.(*Instance).Payload.(*strRdr)
	return int32(sr.readOne()), nil
}

func natStringReaderRead(ctx *CallContext, recv Value, args []Value) (Value, error) {
	sr := recv.(*Instance).Payload.(*strRdr)
	arr := args[0].(*ArrayObj)
	off := argI(args, 1)
	length := argI(args, 2)
	n := 0
	if os.Getenv("CATTY_SRDBG") != "" {
		js, _ := AsJString(args[0])
		_ = js
	}
	for i := 0; i < int(length); i++ {
		c := sr.readOne()
		if c == -1 {
			break
		}
		arr.Elems[off+int32(i)] = int32(c)
		n++
	}
	if n == 0 {
		return int32(-1), nil // JDK Reader contract: EOF is -1, never 0
	}
	if os.Getenv("CATTY_SRDBG") != "" {
		first := ""
		for i := 0; i < n && i < 12; i++ {
			first += string(rune(arr.Elems[int(off)+i].(int32)))
		}
		println("[SR] read off=", int(off), "len=", int(length), "->n=", n, "first=12:", first)
	}
	return int32(n), nil
}

func natDoubleInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	recv.(*Instance).Fields[0] = args[0]
	return nil, nil
}

func natDoubleIsNaN(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(float64)
	return boolV(v != v), nil
}

func natDoubleIsInfinite(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(float64)
	return boolV(v > 1.7976931348623157e308 || v < -1.7976931348623157e308), nil
}

func natDoubleParse(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	v, err := strconv.ParseFloat(strings.TrimSpace(js.String()), 64)
	if err != nil {
		return nil, ctx.Throw("java/lang/NumberFormatException", js.String())
	}
	return v, nil
}

func natDoubleToString(ctx *CallContext, recv Value, args []Value) (Value, error) {
	return ctx.K.InternGo(fmt.Sprintf("%v", args[0].(float64))), nil
}

// ---- ArrayList iterator ----

type itrBuf struct {
	list *alBuf
	idx  int
}

func natArrayListIterator(ctx *CallContext, recv Value, args []Value) (Value, error) {
	itrCls, ok := ctx.K.ClassByName("java/util/ArrayList$Itr")
	if !ok {
		return nil, fmt.Errorf("ArrayList$Itr missing")
	}
	obj, err := ctx.K.NewInstance(itrCls)
	if err != nil {
		return nil, err
	}
	obj.Payload = &itrBuf{list: alOf(recv)}
	return obj, nil
}

func natItrHasNext(ctx *CallContext, recv Value, args []Value) (Value, error) {
	it := recv.(*Instance).Payload.(*itrBuf)
	return boolV(it.idx < len(it.list.data)), nil
}

func natItrNext(ctx *CallContext, recv Value, args []Value) (Value, error) {
	it := recv.(*Instance).Payload.(*itrBuf)
	if it.idx >= len(it.list.data) {
		return nil, ctx.Throw("java/util/NoSuchElementException", "iterator exhausted")
	}
	v := it.list.data[it.idx]
	it.idx++
	return v, nil
}

// ---- arrays-as-list / unmodifiable passthrough ----

func natArraysAsList(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arrCls, ok := ctx.K.ClassByName("java/util/ArrayList")
	if !ok {
		return nil, fmt.Errorf("ArrayList missing")
	}
	obj, err := ctx.K.NewInstance(arrCls)
	if err != nil {
		return nil, err
	}
	buf := &alBuf{}
	if len(args) > 0 {
		if a, ok := args[0].(*ArrayObj); ok {
			buf.data = append(buf.data, a.Elems...)
		}
	}
	obj.Payload = buf
	return obj, nil
}

func natCollectionsUnmodList(ctx *CallContext, recv Value, args []Value) (Value, error) {
	// Deviation (registered): returns the list itself — mutability not
	// enforced; sufficient for read-only third-party consumers.
	return args[0], nil
}
