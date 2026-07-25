package interpreter

import (
	"math"
	"testing"

	"catty/opcode"
	"catty/rtda"
)

func TestInvokeDirectNativeTypedArgumentsAndResults(t *testing.T) {
	method := rtda.NativeMethod(nil, "sum", "(JD)J", func(frame *rtda.Frame) {
		frame.PushLong(frame.GetLong(0) + int64(frame.GetDouble(2)))
	})
	method.SetStatic()
	result := InvokeDirect(rtda.InvocationRequest{
		Context:   rtda.NewExecutionContext(nil),
		Method:    method,
		Arguments: []rtda.JavaValue{rtda.LongValue(math.MinInt64 + 10), rtda.DoubleValue(3)},
	})
	if !result.IsNormal() {
		t.Fatalf("InvokeDirect returned %v: %v", result.Kind(), result.Failure())
	}
	value, _ := result.Value()
	got, ok := value.Long()
	if !ok || got != math.MinInt64+13 {
		t.Fatalf("typed direct long result = %d, %v; want %d", got, ok, math.MinInt64+13)
	}
}

func TestInvokeDirectNativeConstructorAndThrowable(t *testing.T) {
	rtda.InitVMTypes()
	receiver := rtda.NewObject(rtda.VMPrimitiveInt)
	constructor := rtda.NativeMethod(nil, "<init>", "(I)V", func(frame *rtda.Frame) {
		frame.GetRef(0).SetExtra(frame.GetInt(1))
	})
	result := InvokeDirect(rtda.InvocationRequest{
		Context: rtda.NewExecutionContext(nil), Method: constructor,
		Receiver: receiver, Arguments: []rtda.JavaValue{rtda.IntValue(42)},
	})
	if !result.IsNormal() {
		t.Fatalf("constructor direct invocation failed: %v", result.Failure())
	}
	if value, _ := result.Value(); !value.IsVoid() || receiver.Extra() != int32(42) {
		t.Fatal("constructor did not preserve receiver/void result")
	}

	throwable := rtda.NewObject(rtda.VMPrimitiveInt)
	throwing := rtda.NativeMethod(nil, "throws", "()V", func(frame *rtda.Frame) {
		frame.Context().Throw(throwable, 0)
	})
	throwing.SetStatic()
	result = InvokeDirect(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: throwing})
	if got, ok := result.Throwable(); !ok || got != throwable {
		t.Fatal("native Java throwable identity was not preserved")
	}
}

func TestInvokeDirectInterpretedNormalAndAbrupt(t *testing.T) {
	rtda.InitVMTypes()
	receiver := rtda.NewObject(rtda.VMPrimitiveInt)
	normal := rtda.InterpretedMethod(nil, "constant", "()I", 0, 1, 1,
		[]byte{byte(opcode.Iconst5), byte(opcode.Ireturn)}, nil)
	result := InvokeDirect(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: normal, Receiver: receiver})
	value, ok := result.Value()
	got, intOK := value.Int()
	if !ok || !intOK || got != 5 {
		t.Fatalf("interpreted direct normal result = %#v, %v, %v; want int 5", value, ok, intOK)
	}

	throwable := rtda.NewObject(rtda.VMPrimitiveInt)
	abrupt := rtda.InterpretedMethod(nil, "abrupt", "()V", 0, 1, 1,
		[]byte{byte(opcode.Aload0), byte(opcode.Athrow)}, nil)
	result = InvokeDirect(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: abrupt, Receiver: throwable})
	if got, ok := result.Throwable(); !ok || got != throwable {
		t.Fatal("interpreted Java throwable identity was not preserved")
	}
}

func TestInvokeDirectRejectsBoundaryMisuse(t *testing.T) {
	method := rtda.NativeMethod(nil, "one", "(I)I", func(frame *rtda.Frame) { frame.PushInt(frame.GetInt(0)) })
	method.SetStatic()
	context := rtda.NewExecutionContext(nil)
	if result := InvokeDirect(rtda.InvocationRequest{Context: context, Method: method}); !result.IsInternalFailure() {
		t.Fatal("arity mismatch must be an internal failure")
	}
	if result := InvokeDirect(rtda.InvocationRequest{Context: context, Method: method, Arguments: []rtda.JavaValue{rtda.LongValue(1)}}); !result.IsInternalFailure() {
		t.Fatal("primitive widening must not be implicit")
	}
	context.PushFrame(context.NewFrame(method))
	if result := InvokeDirect(rtda.InvocationRequest{Context: context, Method: method, Arguments: []rtda.JavaValue{rtda.IntValue(1)}}); !result.IsInternalFailure() {
		t.Fatal("non-empty execution context stack must be rejected")
	}
	context.PopFrame()
}
