package vm

import (
	"bytes"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"catty/internal/kernel"

	"catty/internal/classfile"
)

// TestSpawnJoinSynthetic exercises Thread.start/join/isAlive against a
// synthesized subclass whose run() is a Go native bumping a counter.
func TestSpawnJoinSynthetic(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: &bytes.Buffer{}})
	var bumped atomic.Int32

	_, err := k.DefineClass(&kernel.ClassDef{
		Name:  "test/Worker",
		Super: "java/lang/Thread",
		Methods: []kernel.MethodDef{
			{Name: "run", Desc: "()V", Flags: classfile.AccPublic, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				for i := 0; i < 50; i++ {
					bumped.Add(1)
					ctx.K.Threads.Sleep(ctx.OwnerKeyValue(), time.Millisecond)
				}
				return nil, nil
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	th := New(k)
	wcls, _ := k.ClassByName("test/Worker")
	obj, err := k.NewInstance(wcls)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := k.InvokeAs(th, mustResolve(t, k, wcls, "<init>", "()V"), obj, nil); err != nil {
		t.Fatal(err)
	}
	j := obj.Payload.(*kernel.JThread)
	if j.StateName() != "NEW" {
		t.Fatalf("fresh thread state = %s", j.StateName())
	}

	start := mustResolve(t, k, wcls, "start", "()V")
	if _, err := k.InvokeAs(th, start, obj, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := k.InvokeAs(th, start, obj, nil); err == nil {
		t.Fatal("double start must throw IllegalThreadStateException")
	}

	// join via the registry latch (native path covered by fixture test)
	ok, interrupted := k.Threads.Join(th.OwnerKey(), j, 5*time.Second)
	if interrupted || !ok {
		t.Fatalf("join = %v,%v", ok, interrupted)
	}
	if n := bumped.Load(); n != 50 {
		t.Errorf("bumped = %d, want 50", n)
	}
	alive := mustResolve(t, k, wcls, "isAlive", "()Z")
	v, err := k.InvokeAs(th, alive, obj, nil)
	if err != nil || v != int32(0) {
		t.Errorf("isAlive after termination = %v,%v", v, err)
	}
}

// TestInterruptSleepingThread verifies the full Java-level path:
// Thread.sleep throws InterruptedException when interrupt() fires.
func TestInterruptSleepingThread(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: &bytes.Buffer{}})
	_, err := k.DefineClass(&kernel.ClassDef{
		Name:  "test/Sleeper",
		Super: "java/lang/Thread",
		Methods: []kernel.MethodDef{
			{Name: "run", Desc: "()V", Flags: classfile.AccPublic, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				tcls := ctx.K.Threads.Main().Obj.Class
				sl, serr := ctx.K.ResolveMethod(tcls, "sleep", "(J)V")
				if serr != nil {
					return nil, serr
				}
				_, verr := ctx.Invoke(sl, nil, []kernel.Value{int64(30_000)})
				return nil, verr // propagate InterruptedException upward
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	th := New(k)
	scls, _ := k.ClassByName("test/Sleeper")
	obj, _ := k.NewInstance(scls)
	k.InvokeAs(th, mustResolve(t, k, scls, "<init>", "(Ljava/lang/String;)V"), obj, []kernel.Value{k.InternGo("s")})

	if _, err := k.InvokeAs(th, mustResolve(t, k, scls, "start", "()V"), obj, nil); err != nil {
		t.Fatal(err)
	}

	// Wait until the sleeper is parked inside sleep, then interrupt.
	j := obj.Payload.(*kernel.JThread)
	dl := time.Now().Add(time.Second)
	for j.SleepersNow() == 0 {
		if time.Now().After(dl) {
			t.Fatalf("sleeper never parked; state=%s alive=%v key=%d",
				j.StateName(), j.IsAlive(), j.Key)
		}
		time.Sleep(time.Millisecond)
	}
	interrupt := mustResolve(t, k, scls, "interrupt", "()V")
	if _, err := k.InvokeAs(th, interrupt, obj, nil); err != nil {
		t.Fatal(err)
	}

	ok, interrupted := k.Threads.Join(th.OwnerKey(), j, 5*time.Second)
	if !ok || interrupted {
		t.Fatalf("join=%v interrupted=%v", ok, interrupted)
	}
	// Uncaught handler ran? We didn't install one — default writes stderr.
	// The assertion that matters: the thread TERMINATED quickly rather than
	// sleeping 30s (join above proves it).
}

func TestStackOverflowDetected(t *testing.T) {
	// static int m() { return m(); } → infinite recursion → SOE
	b := newClassBuilder("Rec")
	self := b.memberRef(classfile.CMethodref, "Rec", "m", "()I")
	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()I",
		maxStack: 1, maxLocals: 0,
		code: []byte{
			0xb8, byte(self >> 8), byte(self), // invokestatic m
			0xac, // ireturn
		},
	})
	k := kernel.New(kernel.Options{SkipVerify: true, MaxFrames: 200})
	cls, err := k.LoadClassBytes(b.serialize(t))
	if err != nil {
		t.Fatal(err)
	}
	th := New(k)
	m, err := k.ResolveMethod(cls, "m", "()I")
	if err != nil {
		t.Fatal(err)
	}
	_, err = th.Call(m, nil, nil)
	var thrownErr *kernel.Thrown
	if !errors.As(err, &thrownErr) ||
		thrownErr.Obj.Class.Name != "java/lang/StackOverflowError" {
		t.Fatalf("want StackOverflowError, got %#v", err)
	}
}

func TestUncaughtHandlerReceivesThrowable(t *testing.T) {
	k := kernel.New(kernel.Options{Stdout: &bytes.Buffer{}, Stderr: &bytes.Buffer{}})
	type captured struct {
		name string
		msg  string
	}
	got := make(chan captured, 1)

	_, err := k.DefineClass(&kernel.ClassDef{
		Name:  "test/Bomb",
		Super: "java/lang/Thread",
		Methods: []kernel.MethodDef{
			{Name: "run", Desc: "()V", Flags: classfile.AccPublic, Native: func(ctx *kernel.CallContext, recv kernel.Value, args []kernel.Value) (kernel.Value, error) {
				return nil, ctx.Throw("java/lang/IllegalStateException", "boom")
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	vt := New(k)
	k.UncaughtHandler = func(j *kernel.JThread, thrown *kernel.Thrown) {
		got <- captured{name: j.Name, msg: strings.TrimSpace(thrown.Error())}
	}
	bcls, _ := k.ClassByName("test/Bomb")
	obj, _ := k.NewInstance(bcls)
	initM := mustResolve(t, k, bcls, "<init>", "()V")
	startM := mustResolve(t, k, bcls, "start", "()V")
	if _, err := k.InvokeAs(vt, initM, obj, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := k.InvokeAs(vt, startM, obj, nil); err != nil {
		t.Fatal(err)
	}
	select {
	case c := <-got:
		if c.msg == "" || !strings.Contains(c.msg, "boom") {
			t.Errorf("uncaught payload = %+v", c)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("handler never fired")
	}
}

func mustClass(t *testing.T, k *kernel.Kernel, name string) *kernel.Class {
	t.Helper()
	c, ok := k.ClassByName(name)
	if !ok {
		t.Fatalf("class %s missing", name)
	}
	return c
}

func mustResolve(t *testing.T, k *kernel.Kernel, cls *kernel.Class, name, desc string) *kernel.Method {
	t.Helper()
	m, err := k.ResolveMethod(cls, name, desc)
	if err != nil {
		t.Fatalf("resolve %s.%s%s: %v", cls.Name, name, desc, err)
	}
	return m
}
