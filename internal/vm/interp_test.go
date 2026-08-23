package vm

import (
	"errors"
	"math"
	"testing"

	"catty/internal/classfile"
	"catty/internal/kernel"
)

// opcodes used by hand-written tests
const (
	bIconst0  = 0x03
	bIconst1  = 0x04
	bBipush   = 0x10
	bLdcW     = 0x13
	bLdc2W    = 0x14
	bIload0   = 0x1A
	bIstore0  = 0x3B
	bIinc     = 0x84
	bIfIcmpgt = 0xA3
	bGoto     = 0xA7
	bLookupsw = 0xAB
	bIreturn  = 0xAC
	bLreturn  = 0xAD
	bAreturn  = 0xB0
	bIdiv     = 0x6C
	bLdiv     = 0x6D
	bD2i      = 0x8E
	bNew      = 0xBB
	bInvokeV  = 0xB6
	bInvokeS  = 0xB7
	bDup      = 0x59
)

func TestLongDivisionWraps(t *testing.T) {
	b := newClassBuilder("LongOps")
	lmin := b.long(math.MinInt64)
	lone := b.long(-1)
	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()J",
		maxStack: 4, maxLocals: 0,
		code: []byte{
			bLdc2W, byte(lmin >> 8), byte(lmin),
			bLdc2W, byte(lone >> 8), byte(lone),
			bLdiv,
			bLreturn,
		},
	})
	k, cls := b.build(t)
	v, err := callStatic(t, k, cls, "m", "()J", nil)
	if err != nil {
		t.Fatal(err)
	}
	lv, ok := v.(int64)
	if !ok {
		t.Fatalf("type %T", v)
	}
	if lv != math.MinInt64 {
		t.Errorf("MinInt64 / -1 = %v", lv)
	}
}

func TestDivByZeroThrows(t *testing.T) {
	b := newClassBuilder("DivZero")
	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()I",
		maxStack: 2, maxLocals: 0,
		code: []byte{
			bIconst1,
			bIconst0,
			bIdiv,
			bIreturn,
		},
	})
	k, cls := b.build(t)
	_, err := callStatic(t, k, cls, "m", "()I", nil)
	var th *kernel.Thrown
	if !errors.As(err, &th) || th.Obj.Class.Name != "java/lang/ArithmeticException" {
		t.Fatalf("want ArithmeticException, got %#v", err)
	}
}

func TestExceptionHandlerCatches(t *testing.T) {
	// try { return 1/0; } catch (ArithmeticException e) { return 42; }
	b := newClassBuilder("CatchIt")
	fortyTwo := b.integer(42)
	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()I",
		maxStack: 2, maxLocals: 1,
		code: []byte{
			bIconst1, // 0
			bIconst0, // 1
			bIdiv,    // 2 — fault
			bIreturn, // 3
			bLdcW,    // 4? handler at 4: ldc_w 42
			byte(fortyTwo >> 8), byte(fortyTwo),
			bIreturn, // 7
		},
		handlers: []excHandler{{start: 0, end: 4, handler: 4, catchClassName: "java/lang/ArithmeticException"}},
	})
	k, cls := b.build(t)
	v, err := callStatic(t, k, cls, "m", "()I", nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(42) {
		t.Errorf("handler result = %v, want 42", v)
	}
}

func TestSumLoopAligned(t *testing.T) {
	// Exact-pc layout (verified by hand, pcs annotated):
	// 0  iconst_0
	// 1  istore_0
	// 2  iconst_0
	// 3  istore_1
	// 4  iload_1            (loop head)
	// 5  bipush 10
	// 7  if_icmple → 20     (7 + 13)
	// 10 iload_0
	// 11 iload_1
	// 12 iadd
	// 13 istore_0
	// 14 iinc 1,+1
	// 17 goto -13           (17-13 = 4)
	// 20 iload_0
	// 21 ireturn
	code := []byte{
		bIconst0,     // 0
		bIstore0,     // 1
		bIconst0,     // 2
		bIstore0 + 1, // 3
		bIload0 + 1,  // 4
		bBipush, 10,  // 5..6
		bIfIcmpgt, 0, 13, // 7..9: exit when i > 10 → 20 (iload_0)
		bIload0,     // 10
		bIload0 + 1, // 11
		0x60,        // 12 iadd
		bIstore0,    // 13
		bIinc, 1, 1, // 14..16 (idx=1, const=+1 as signed byte)
		bGoto, 0xFF, 0xF3, // 17..19 → 4
		bIload0,  // 20
		bIreturn, // 21
	}
	b := newClassBuilder("SumLoop")
	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()I",
		maxStack: 2, maxLocals: 2,
		code: code,
	})
	k, cls := b.build(t)
	v, err := callStatic(t, k, cls, "m", "()I", nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(55) {
		t.Errorf("sum 0..10 = %v, want 55", v)
	}
}

func TestLookupswitchAligned(t *testing.T) {
	// pcs: 0 iload_0; 1 lookupswitch (pad→4); 4 default@16; 8 npairs=1;
	// 12 match=7; 16 offset=8 → 1+8=9; …
	// 9 bipush 100; 11 ireturn; 17 ldc_w -1; 20 ireturn
	b := newClassBuilder("SwitchAligned")
	neg1 := b.integer(-1)

	// Correct layout:
	// 0: iload_0
	// 1: lookupswitch; 2-3 pad; 4-7 default; 8-11 npairs; 12-15 match; 16-19 off
	// 20: bipush 100; 22: ireturn
	// 23: ldc_w -1; 26: ireturn
	code2 := []byte{}
	code2 = append(code2, 0x1A)                             // 0
	code2 = append(code2, 0xAB)                             // 1
	code2 = append(code2, 0, 0)                             // 2..3 pad → default at pc 4
	code2 = append(code2, 0, 0, 0, 22)                      // 4..7 default → 1+22=23
	code2 = append(code2, 0, 0, 0, 1)                       // 8..11
	code2 = append(code2, 0, 0, 0, 7)                       // 12..15
	code2 = append(code2, 0, 0, 0, 19)                      // 16..19 → 1+19=20
	code2 = append(code2, bBipush, 100)                     // 20..21
	code2 = append(code2, bIreturn)                         // 22
	code2 = append(code2, bLdcW, byte(neg1>>8), byte(neg1)) // 23..25
	code2 = append(code2, bIreturn)                         // 26

	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "(I)I",
		maxStack: 2, maxLocals: 1,
		code: code2,
	})
	k, cls := b.build(t)

	v, err := callStatic(t, k, cls, "m", "(I)I", []kernel.Value{int32(7)})
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(100) {
		t.Errorf("case hit = %v, want 100", v)
	}
	v, err = callStatic(t, k, cls, "m", "(I)I", []kernel.Value{int32(3)})
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(-1) {
		t.Errorf("default = %v, want -1", v)
	}
}

func TestFloatConversionsSaturate(t *testing.T) {
	b := newClassBuilder("Conv")
	dBig := b.double(1e18)
	dNan := b.double(math.NaN())
	b.addMethod(methodBlob{
		flags: 0x0009, name: "big", desc: "()I",
		maxStack: 2, maxLocals: 0,
		code: []byte{
			bLdc2W, byte(dBig >> 8), byte(dBig),
			bD2i,
			bIreturn,
		},
	})
	b.addMethod(methodBlob{
		flags: 0x0009, name: "nan", desc: "()I",
		maxStack: 2, maxLocals: 0,
		code: []byte{
			bLdc2W, byte(dNan >> 8), byte(dNan),
			bD2i,
			bIreturn,
		},
	})
	k, cls := b.build(t)
	v, err := callStatic(t, k, cls, "big", "()I", nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(math.MaxInt32) {
		t.Errorf("d2i(1e18) = %v", v)
	}
	v, err = callStatic(t, k, cls, "nan", "()I", nil)
	if err != nil {
		t.Fatal(err)
	}
	if v != int32(0) {
		t.Errorf("d2i(NaN) = %v", v)
	}
}

func TestStringConcatViaBuilder(t *testing.T) {
	// new StringBuilder; dup; <init>; ldc "x="; append; ldc 7; append;
	// toString; areturn — crosses kernel natives through the interpreter.
	b := newClassBuilder("Concat")
	sbCls := b.class("java/lang/StringBuilder")
	sbInit := b.memberRef(classfile.CMethodref, "java/lang/StringBuilder", "<init>", "()V")
	sbAppendStr := b.memberRef(classfile.CMethodref, "java/lang/StringBuilder", "append", "(Ljava/lang/String;)Ljava/lang/StringBuilder;")
	sbAppendInt := b.memberRef(classfile.CMethodref, "java/lang/StringBuilder", "append", "(I)Ljava/lang/StringBuilder;")
	sbToString := b.memberRef(classfile.CMethodref, "java/lang/StringBuilder", "toString", "()Ljava/lang/String;")
	strX := b.stringConst("x=")
	int7 := b.integer(7)

	b.addMethod(methodBlob{
		flags: 0x0009, name: "m", desc: "()Ljava/lang/String;",
		maxStack: 4, maxLocals: 0,
		code: []byte{
			bNew, byte(sbCls >> 8), byte(sbCls),
			bDup,
			bInvokeS, byte(sbInit >> 8), byte(sbInit),
			bLdcW, byte(strX >> 8), byte(strX),
			bInvokeV, byte(sbAppendStr >> 8), byte(sbAppendStr),
			bLdcW, byte(int7 >> 8), byte(int7),
			bInvokeV, byte(sbAppendInt >> 8), byte(sbAppendInt),
			bInvokeV, byte(sbToString >> 8), byte(sbToString),
			bAreturn,
		},
	})
	k, cls := b.build(t)
	th := New(k)
	m, err := k.ResolveMethod(cls, "m", "()Ljava/lang/String;")
	if err != nil {
		t.Fatal(err)
	}
	v, err := th.Call(m, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	s, ok := v.(*kernel.JString)
	if !ok {
		t.Fatalf("want JString, got %T", v)
	}
	if s.String() != "x=7" {
		t.Errorf("concat result = %q, want x=7", s.String())
	}
}
