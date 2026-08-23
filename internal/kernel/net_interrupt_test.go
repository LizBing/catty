package kernel

import (
	"net"
	"testing"
	"time"
)

type intTestOwner struct{ key uint64 }

func (o *intTestOwner) OwnerKey() uint64 { return o.key }

// TestSocketReadInterruptible pins the DEBT-0011 contract: a thread parked
// in Socket InputStream.read MUST be woken by Thread.interrupt (via
// SetDeadline-on-interrupt) and surface InterruptedIOException while its
// interrupt flag stays set.
func TestSocketReadInterruptible(t *testing.T) {
	k := New(Options{Stdout: discard{}, Stderr: discard{}})

	// Server: accept one conn, never write.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	type accepted struct {
		conn net.Conn
		err  error
	}
	acc := make(chan accepted, 1)
	go func() {
		c, e := ln.Accept()
		acc <- accepted{c, e}
	}()

	// Client dials directly (bypasses accept-native; same net.Conn shape).
	client, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	srv := <-acc
	if srv.err != nil {
		t.Fatal(srv.err)
	}
	defer srv.conn.Close()

	// Kernel-side stream instance wrapping the client conn.
	stream, err := k.newStreamInstance("java/io/InputStream", client, client)
	if err != nil {
		t.Fatal(err)
	}

	// Thread record for the reader goroutine.
	key := k.MintKey()
	owner := &intTestOwner{key: key}
	threadCls, _ := k.ClassByName("java/lang/Thread")
	var dummy *Instance
	if threadCls != nil {
		dummy, _ = k.NewInstance(threadCls)
	}
	jt := k.Threads.Register(key, dummy, "reader")
	_ = jt

	ctx := &CallContext{K: k, Owner: owner}

	// Park a goroutine in read.
	done := make(chan struct {
		v Value
		e error
	}, 1)
	buf := &ArrayObj{Elems: make([]Value, 32)}
	go func() {
		v, e := natStreamReadB(ctx, stream, []Value{buf})
		done <- struct {
			v Value
			e error
		}{v, e}
	}()

	// Wait until the read has registered its conn (parked).
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if jt.netConn.Load() != nil {
			break
		}
		time.Sleep(2 * time.Millisecond)
	}
	if jt.netConn.Load() == nil {
		t.Fatal("read never registered its conn — cannot test interrupt wakeup")
	}

	// Interrupt: must SetDeadline(now) and unblock the read.
	start := time.Now()
	k.Threads.InterruptByKey(jt)

	res := <-done
	elapsed := time.Since(start)

	if res.e == nil {
		t.Fatalf("read returned %v without error after interrupt", res.v)
	}
	thrown, ok := res.e.(*Thrown)
	if !ok {
		t.Fatalf("expected *Thrown, got %v", res.e)
	}
	// Class check via message marker is fragile; verify through the class
	// name carried by the thrown instance.
	if thrown.Obj == nil || thrown.Obj.Class.Name != "java/io/InterruptedIOException" {
		name := "<nil>"
		if thrown.Obj != nil {
			name = thrown.Obj.Class.Name
		}
		t.Fatalf("want InterruptedIOException, got %s: %v", name, thrown.Error())
	}

	// Wakeup must be prompt.
	if elapsed > time.Second {
		t.Fatalf("interrupt wakeup took %v, want <1s", elapsed)
	}

	// Interrupt flag must remain set (InterruptedIOException does not clear).
	if !jt.interrupted.Load() {
		t.Fatal("interrupt flag was cleared by read wakeup")
	}
}
