package kernel

import (
	"bufio"
	"os"
)

// ---- P-0008 Route A: file I/O minimal surface ----
// FileReader is a pragmatic line-oriented reader implemented natively
// over bufio.Scanner; FileInputStream exposes the classic byte API.

type fileBuf struct{ f *os.File }

func natFileExists(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	_, err := os.Stat(js.String())
	return boolV(err == nil), nil
}

func natFileLength(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	fi, err := os.Stat(js.String())
	if err != nil {
		return int64(0), nil
	}
	return fi.Size(), nil
}

func natFileReaderInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	js, _ := AsJString(args[0])
	f, oerr := os.Open(js.String())
	if oerr != nil {
		return nil, ctx.Throw("java/io/FileNotFoundException", oerr.Error())
	}
	recv.(*Instance).Payload = &fileBuf{f: f}
	return nil, nil
}

func natFileClose(ctx *CallContext, recv Value, args []Value) (Value, error) {
	if fb, ok := recv.(*Instance).Payload.(*fileBuf); ok {
		_ = fb.f.Close()
		recv.(*Instance).Payload = nil
	}
	return nil, nil
}

// ---- BufferedReader over FileReader ----

type bufRdr struct{ sc *bufio.Scanner }

func natBufferedReaderInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	inner, ok := args[0].(*Instance)
	if !ok {
		return nil, ctx.Throw("java/lang/IllegalArgumentException", "reader")
	}
	fb, ok := inner.Payload.(*fileBuf)
	if !ok {
		return nil, ctx.Throw("java/io/IOException", "unsupported reader")
	}
	recv.(*Instance).Payload = &bufRdr{sc: bufio.NewScanner(fb.f)}
	return nil, nil
}

func natBufferedReadLine(ctx *CallContext, recv Value, args []Value) (Value, error) {
	br := recv.(*Instance).Payload.(*bufRdr)
	if !br.sc.Scan() {
		if serr := br.sc.Err(); serr != nil {
			return nil, ctx.Throw("java/io/IOException", serr.Error())
		}
		return nil, nil // EOF: null
	}
	return ctx.K.InternGo(br.sc.Text()), nil
}

// ---- java.lang.Math (Route C surface) ----

func natMathMaxInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := argI(args, 0), argI(args, 1)
	if a > b {
		return a, nil
	}
	return b, nil
}

func natMathMinInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := argI(args, 0), argI(args, 1)
	if a < b {
		return a, nil
	}
	return b, nil
}

func natMathMaxLong(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := args[0].(int64), args[1].(int64)
	if a > b {
		return a, nil
	}
	return b, nil
}

func natMathMinLong(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := args[0].(int64), args[1].(int64)
	if a < b {
		return a, nil
	}
	return b, nil
}

func natMathMaxDouble(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := args[0].(float64), args[1].(float64)
	if a > b || b != b {
		return a, nil
	}
	return b, nil
}

func natMathMinDouble(ctx *CallContext, recv Value, args []Value) (Value, error) {
	a, b := args[0].(float64), args[1].(float64)
	if a < b || b != b {
		return a, nil
	}
	return b, nil
}

func natMathAbsInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := argI(args, 0)
	if v < 0 {
		return -v, nil
	}
	return v, nil
}

func natMathAbsLong(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(int64)
	if v < 0 {
		return -v, nil
	}
	return v, nil
}

func natMathAbsDouble(ctx *CallContext, recv Value, args []Value) (Value, error) {
	v := args[0].(float64)
	if v < 0 || v != v {
		return -v, nil
	}
	return v, nil
}

// ---- java.util.Arrays fills ----

func natArraysFillChar(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	fill := rune(args[1].(int32))
	for i := range arr.Elems {
		arr.Elems[i] = int32(fill)
	}
	return nil, nil
}

func natArraysFillInt(ctx *CallContext, recv Value, args []Value) (Value, error) {
	arr := args[0].(*ArrayObj)
	fill := argI(args, 1)
	for i := range arr.Elems {
		arr.Elems[i] = fill
	}
	return nil, nil
}
