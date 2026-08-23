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
