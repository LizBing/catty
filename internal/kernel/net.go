package kernel

import (
	"fmt"
	"io"
	"net"
)

// Minimal java.net / stream surface over Go's net package (M2).
//
// Blocking reads run on the Java thread's goroutine — the payoff of the
// goroutine-backed thread model. Interrupt-wiring for socket reads
// (SetDeadline-on-interrupt) is deliberately deferred and registered as a
// deviation; see ledger DEV-0008.

func netListenerOf(v Value) (net.Listener, error) {
	in, ok := v.(*Instance)
	if !ok {
		return nil, fmt.Errorf("receiver is not a heap object")
	}
	ln, ok := in.Payload.(net.Listener)
	if !ok {
		return nil, fmt.Errorf("ServerSocket payload missing/closed")
	}
	return ln, nil
}

func netConnOf(v Value) (net.Conn, error) {
	in, ok := v.(*Instance)
	if !ok {
		return nil, fmt.Errorf("receiver is not a heap object")
	}
	c, ok := in.Payload.(net.Conn)
	if !ok {
		return nil, fmt.Errorf("Socket payload missing/closed")
	}
	return c, nil
}

// newStreamInstance wraps an io reader/writer pair into the synthesized
// stream classes. Both halves of one Socket share the same net.Conn.
func (k *Kernel) newStreamInstance(className string, r io.Reader, w io.Writer) (*Instance, error) {
	cls, ok := k.ClassByName(className)
	if !ok {
		return nil, fmt.Errorf("bootstrap %s missing", className)
	}
	in, err := k.NewInstance(cls)
	if err != nil {
		return nil, err
	}
	in.Payload = &streamHandle{r: r, w: w}
	return in, nil
}

type streamHandle struct {
	r io.Reader
	w io.Writer
}

func streamOf(recv Value) (*streamHandle, error) {
	in, ok := recv.(*Instance)
	if !ok {
		return nil, fmt.Errorf("receiver is not a heap object")
	}
	h, ok := in.Payload.(*streamHandle)
	if !ok || h == nil {
		return nil, fmt.Errorf("stream closed")
	}
	return h, nil
}

// ---- ServerSocket ----------------------------------------------------------------

func natServerSocketInit(ctx *CallContext, recv Value, args []Value) (Value, error) {
	port := argI(args, 0)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, ctx.Throw("java/net/BindException", err.Error())
	}
	recv.(*Instance).Payload = ln
	return nil, nil
}

func natServerSocketAccept(ctx *CallContext, recv Value, args []Value) (Value, error) {
	ln, err := netListenerOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	conn, aerr := ln.Accept()
	if aerr != nil {
		return nil, ctx.Throw("java/net/SocketException", aerr.Error())
	}
	sockCls, ok := ctx.K.ClassByName("java/net/Socket")
	if !ok {
		return nil, fmt.Errorf("bootstrap java/net/Socket missing")
	}
	sock, cerr := ctx.K.NewInstance(sockCls)
	if cerr != nil {
		return nil, cerr
	}
	sock.Payload = conn
	return sock, nil
}

func natServerSocketLocalPort(ctx *CallContext, recv Value, args []Value) (Value, error) {
	ln, err := netListenerOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	if tcp, ok := ln.Addr().(*net.TCPAddr); ok {
		return int32(tcp.Port), nil
	}
	return int32(0), nil
}

func natServerSocketClose(ctx *CallContext, recv Value, args []Value) (Value, error) {
	ln, err := netListenerOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	cerr := ln.Close()
	if cerr != nil {
		return nil, ctx.Throw("java/net/SocketException", cerr.Error())
	}
	recv.(*Instance).Payload = nil
	return nil, nil
}

// ---- Socket ------------------------------------------------------------------------

func natSocketGetInputStream(ctx *CallContext, recv Value, args []Value) (Value, error) {
	conn, err := netConnOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	return ctx.K.newStreamInstance("java/io/InputStream", conn, conn)
}

func natSocketGetOutputStream(ctx *CallContext, recv Value, args []Value) (Value, error) {
	conn, err := netConnOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	return ctx.K.newStreamInstance("java/io/OutputStream", conn, conn)
}

func natSocketClose(ctx *CallContext, recv Value, args []Value) (Value, error) {
	conn, err := netConnOf(recv)
	if err == nil {
		if cerr := conn.Close(); cerr != nil {
			return nil, ctx.Throw("java/net/SocketException", cerr.Error())
		}
		recv.(*Instance).Payload = nil
	}
	return nil, nil
}

// ---- streams -----------------------------------------------------------------------

func natStreamReadB(ctx *CallContext, recv Value, args []Value) (Value, error) {
	h, err := streamOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	arr := args[0].(*ArrayObj)
	tmp := make([]byte, len(arr.Elems))
	n, rerr := h.r.Read(tmp)
	for i := 0; i < n; i++ {
		arr.Elems[i] = int32(tmp[i])
	}
	if rerr != nil {
		if n > 0 {
			return int32(n), nil // partial read then error next call
		}
		return int32(-1), nil // EOF convention (IOException mapping deferred)
	}
	return int32(n), nil
}

func natStreamWriteBII(ctx *CallContext, recv Value, args []Value) (Value, error) {
	h, err := streamOf(recv)
	if err != nil {
		return nil, ctx.Throw("java/net/SocketException", err.Error())
	}
	arr := args[0].(*ArrayObj)
	off := argI(args, 1)
	length := argI(args, 2)
	b, berr := elemsToBytes(arr, off, length)
	if berr != nil {
		return nil, ctx.Throw("java/lang/IndexOutOfBoundsException", berr.Error())
	}
	if _, werr := h.w.Write(b); werr != nil {
		return nil, ctx.Throw("java/net/SocketException", werr.Error())
	}
	return nil, nil
}

func natStreamClose(ctx *CallContext, recv Value, args []Value) (Value, error) {
	h, err := streamOf(recv)
	if err != nil {
		return nil, nil // already closed: idempotent per Close contract
	}
	if c, ok := h.r.(io.Closer); ok {
		_ = c.Close()
	}
	recv.(*Instance).Payload = nil
	return nil, nil
}
