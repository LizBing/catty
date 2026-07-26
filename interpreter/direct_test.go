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
	normalOwner := rtda.NewSyntheticClass("test/DirectNormal", nil)
	normal := rtda.InterpretedMethod(normalOwner, "constant", "()I", 0x0009, 1, 0,
		[]byte{byte(opcode.Iconst5), byte(opcode.Ireturn)}, nil)
	normalOwner.AddMethod(normal)
	result := InvokeDirect(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: normal})
	value, ok := result.Value()
	got, intOK := value.Int()
	if !ok || !intOK || got != 5 {
		t.Fatalf("interpreted direct normal result = %#v, %v, %v; want int 5", value, ok, intOK)
	}
	result = InvokeDirectIR(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: normal})
	value, ok = result.Value()
	got, intOK = value.Int()
	if !ok || !intOK || got != 5 {
		t.Fatalf("IR direct normal result = %#v, %v, %v; want int 5", value, ok, intOK)
	}

	abruptOwner := rtda.NewSyntheticClass("test/DirectAbrupt", nil)
	throwable := rtda.NewObject(abruptOwner)
	abrupt := rtda.InterpretedMethod(abruptOwner, "abrupt", "()V", 0, 1, 1,
		[]byte{byte(opcode.Aload0), byte(opcode.Athrow)}, nil)
	abruptOwner.AddMethod(abrupt)
	result = InvokeDirect(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: abrupt, Receiver: throwable})
	if got, ok := result.Throwable(); !ok || got != throwable {
		t.Fatal("interpreted Java throwable identity was not preserved")
	}
	result = InvokeDirectIR(rtda.InvocationRequest{Context: rtda.NewExecutionContext(nil), Method: abrupt, Receiver: throwable})
	if got, ok := result.Throwable(); !ok || got != throwable {
		t.Fatal("IR Java throwable identity was not preserved")
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
	assertDirectInt(t, InvokeDirect(rtda.InvocationRequest{Context: context, Method: method, Arguments: []rtda.JavaValue{rtda.IntValue(1)}}), 1)
	if context.FrameCount() != 1 || context.CurrentFrame().StackSize() != 0 {
		t.Fatal("direct native invocation must preserve the existing Java caller frame")
	}
	context.PopFrame()
}

func TestInvokeDirectValidatesReferenceArgumentAssignment(t *testing.T) {
	loader := newTestLoader()
	target := rtda.NewSyntheticClass("test/ArgumentTarget", nil)
	other := rtda.NewSyntheticClass("test/ArgumentOther", nil)
	owner := rtda.NewSyntheticClass("test/ArgumentOwner", nil)
	for _, class := range []*rtda.Class{target, other, owner} {
		loader.addClass(class)
	}
	method := rtda.InterpretedMethod(owner, "accept", "(Ltest/ArgumentTarget;)V", 0x0009, 0, 1,
		[]byte{byte(opcode.Return)}, nil)
	owner.AddMethod(method)

	for _, invoke := range []struct {
		name string
		fn   func(rtda.InvocationRequest) rtda.DynamicResult
	}{
		{"tree", InvokeDirect},
		{"ir", InvokeDirectIR},
	} {
		t.Run(invoke.name, func(t *testing.T) {
			context := rtda.NewExecutionContext(loader)
			result := invoke.fn(rtda.InvocationRequest{
				Context: context, Method: method,
				Arguments: []rtda.JavaValue{rtda.ReferenceValue(rtda.NewObject(other))},
			})
			if !result.IsInternalFailure() {
				t.Fatalf("incompatible reference argument result = %#v; want internal failure", result)
			}
			if context.FrameCount() != 0 {
				t.Fatal("rejected reference argument modified the execution context")
			}

			result = invoke.fn(rtda.InvocationRequest{
				Context: context, Method: method,
				Arguments: []rtda.JavaValue{rtda.ReferenceValue(nil)},
			})
			if !result.IsNormal() {
				t.Fatalf("null reference argument result = %#v; want normal", result)
			}
		})
	}
}

func TestInvokeDirectIRRestoresEntryDepthAfterLoweringFailure(t *testing.T) {
	parent := rtda.NativeMethod(nil, "parent", "()V", func(*rtda.Frame) {})
	parent.SetStatic()
	owner := rtda.NewSyntheticClass("test/InvalidIR", nil)
	invalid := rtda.InterpretedMethod(owner, "invalid", "()V", 0x0009, 1, 0,
		[]byte{
			byte(opcode.Iconst0),
			byte(opcode.Ifeq), 0, 4,
			byte(opcode.Iconst1),
			byte(opcode.Return),
		}, nil)
	owner.AddMethod(invalid)
	context := rtda.NewExecutionContext(nil)
	caller := context.NewFrame(parent)
	context.PushFrame(caller)

	result := InvokeDirectIR(rtda.InvocationRequest{Context: context, Method: invalid})
	if !result.IsInternalFailure() {
		t.Fatalf("invalid IR result = %#v; want internal failure", result)
	}
	if context.FrameCount() != 1 || context.CurrentFrame() != caller || caller.StackSize() != 0 {
		t.Fatal("IR internal failure did not restore the typed invocation entry depth")
	}
	context.PopFrame()
}

func TestInvokeDirectNestedInterpretedBoundary(t *testing.T) {
	parent := rtda.NativeMethod(nil, "parent", "()V", func(*rtda.Frame) {})
	parent.SetStatic()
	normalOwner := rtda.NewSyntheticClass("test/NestedNormal", nil)
	normal := rtda.InterpretedMethod(normalOwner, "constant", "()J", 0x0009, 2, 0,
		[]byte{byte(opcode.Lconst1), byte(opcode.Lreturn)}, nil)
	normalOwner.AddMethod(normal)
	abruptOwner := rtda.NewSyntheticClass("test/NestedAbrupt", nil)
	throwable := rtda.NewObject(abruptOwner)
	abrupt := rtda.InterpretedMethod(abruptOwner, "abrupt", "()V", 0, 1, 1,
		[]byte{byte(opcode.Aload0), byte(opcode.Athrow)}, nil)
	abruptOwner.AddMethod(abrupt)

	for _, invoke := range []struct {
		name string
		fn   func(rtda.InvocationRequest) rtda.DynamicResult
	}{
		{"tree", InvokeDirect},
		{"ir", InvokeDirectIR},
	} {
		t.Run(invoke.name, func(t *testing.T) {
			context := rtda.NewExecutionContext(nil)
			context.PushFrame(context.NewFrame(parent))
			result := invoke.fn(rtda.InvocationRequest{Context: context, Method: normal})
			value, ok := result.Value()
			got, longOK := value.Long()
			if !ok || !longOK || got != 1 {
				t.Fatalf("nested normal result = %#v, normal=%v long=%v", result, ok, longOK)
			}
			if context.FrameCount() != 1 || context.CurrentFrame().StackSize() != 0 {
				t.Fatal("nested normal invocation modified the Java caller frame")
			}
			result = invoke.fn(rtda.InvocationRequest{Context: context, Method: abrupt, Receiver: throwable})
			if gotThrowable, ok := result.Throwable(); !ok || gotThrowable != throwable {
				t.Fatal("nested abrupt invocation lost throwable identity")
			}
			if context.FrameCount() != 1 || context.CurrentFrame().StackSize() != 0 {
				t.Fatal("nested abrupt invocation unwound the Java caller frame")
			}
			context.PopFrame()
		})
	}
}

func TestInvokeDispatchInterpreterAndIR(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	owner := rtda.NewSyntheticClass("test/Dispatch", nil)
	loader.addClass(owner)
	static := rtda.NativeMethod(owner, "staticValue", "()I", func(frame *rtda.Frame) { frame.PushInt(7) })
	static.SetStatic()
	owner.AddMethod(static)
	virtual := rtda.NativeMethod(owner, "virtualValue", "()I", func(frame *rtda.Frame) { frame.PushInt(9) })
	owner.AddMethod(virtual)
	constructor := rtda.NativeMethod(owner, "<init>", "(I)V", func(frame *rtda.Frame) {
		frame.GetRef(0).SetExtra(frame.GetInt(1))
	})
	owner.AddMethod(constructor)

	staticRequest := rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner,
		Kind: rtda.InvokeStatic, Name: "staticValue", Descriptor: "()I",
	}
	assertDirectInt(t, InvokeDispatch(staticRequest), 7)
	if !owner.IsInitialized() {
		t.Fatal("direct static dispatch did not initialize declaring class")
	}
	staticRequest.Context = rtda.NewExecutionContext(loader)
	assertDirectInt(t, InvokeDispatchIR(staticRequest), 7)
	nestedContext := rtda.NewExecutionContext(loader)
	nestedContext.PushFrame(nestedContext.NewFrame(static))
	staticRequest.Context = nestedContext
	assertDirectInt(t, InvokeDispatch(staticRequest), 7)
	if nestedContext.FrameCount() != 1 || nestedContext.CurrentFrame().StackSize() != 0 {
		t.Fatal("nested direct dispatch modified the Java caller frame")
	}
	nestedContext.PopFrame()

	receiver := rtda.NewObject(owner)
	virtualRequest := rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner, Receiver: receiver,
		Kind: rtda.InvokeVirtual, Name: "virtualValue", Descriptor: "()I",
	}
	assertDirectInt(t, InvokeDispatch(virtualRequest), 9)
	virtualRequest.Context = rtda.NewExecutionContext(loader)
	virtualRequest.Kind = rtda.InvokeInterface
	assertDirectInt(t, InvokeDispatchIR(virtualRequest), 9)
	virtualRequest.Context = rtda.NewExecutionContext(loader)
	virtualRequest.Kind = rtda.InvokeSpecial
	assertDirectInt(t, InvokeDispatch(virtualRequest), 9)

	constructorRequest := rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner, Receiver: receiver,
		Kind: rtda.InvokeConstructor, Name: "<init>", Descriptor: "(I)V",
		Arguments: []rtda.JavaValue{rtda.IntValue(11)},
	}
	if result := InvokeDispatch(constructorRequest); !result.IsNormal() {
		t.Fatalf("constructor dispatch failed: %v", result.Failure())
	}
	if receiver.Extra() != int32(11) {
		t.Fatal("constructor dispatch did not preserve receiver state")
	}
}

func TestInvokeConstructorDoesNotInheritSuperclassConstructor(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addClass(rtda.NewSyntheticClass("java/lang/NoSuchMethodError", nil))
	super := rtda.NewSyntheticClass("test/ConstructorSuper", nil)
	constructor := rtda.NativeMethod(super, "<init>", "()V", func(frame *rtda.Frame) {
		frame.GetRef(0).SetExtra("super constructor invoked")
	})
	super.AddMethod(constructor)
	target := rtda.NewSyntheticClass("test/ConstructorTarget", super)
	loader.addClass(super)
	loader.addClass(target)
	receiver := rtda.NewObject(target)

	result := InvokeDispatch(rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: target, Receiver: receiver,
		Kind: rtda.InvokeConstructor, Name: "<init>", Descriptor: "()V",
	})
	throwable, ok := result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/NoSuchMethodError" {
		t.Fatalf("inherited constructor result = %#v, throwable=%v; want NoSuchMethodError", result, ok)
	}
	if receiver.Extra() != nil {
		t.Fatal("constructor dispatch invoked an inherited superclass constructor")
	}
}

func TestInvokeDispatchNullReceiverIsJavaThrowable(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addClass(rtda.NewSyntheticClass("java/lang/NoSuchMethodError", nil))
	loader.addClass(rtda.NewSyntheticClass("java/lang/IncompatibleClassChangeError", nil))
	loader.addClass(rtda.NewSyntheticClass("java/lang/ClassCastException", nil))
	owner := rtda.NewSyntheticClass("test/NullReceiver", nil)
	loader.addClass(owner)
	result := InvokeDispatch(rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner,
		Kind: rtda.InvokeVirtual, Name: "missing", Descriptor: "()V",
	})
	throwable, ok := result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/NullPointerException" {
		t.Fatalf("null receiver result = %#v, throwable=%v; want NullPointerException", result, ok)
	}
	field := rtda.NewField(owner, "i", "I", 0x0001, false, 0)
	result = ReadDirectField(rtda.NewExecutionContext(loader), field, nil)
	throwable, ok = result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/NullPointerException" {
		t.Fatalf("null field receiver result = %#v, throwable=%v; want NullPointerException", result, ok)
	}
	result = InvokeDispatch(rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner,
		Kind: rtda.InvokeStatic, Name: "missing", Descriptor: "()V",
	})
	throwable, ok = result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/NoSuchMethodError" {
		t.Fatalf("missing method result = %#v, throwable=%v; want NoSuchMethodError", result, ok)
	}
	other := rtda.NewSyntheticClass("test/Other", nil)
	result = InvokeDispatch(rtda.InvocationLookupRequest{
		Context: rtda.NewExecutionContext(loader), Target: owner, Receiver: rtda.NewObject(other),
		Kind: rtda.InvokeVirtual, Name: "missing", Descriptor: "()V",
	})
	throwable, ok = result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/IncompatibleClassChangeError" {
		t.Fatalf("receiver mismatch result = %#v, throwable=%v; want IncompatibleClassChangeError", result, ok)
	}
	targetClass := rtda.NewSyntheticClass("test/FieldTarget", nil)
	otherClass := rtda.NewSyntheticClass("test/FieldOther", nil)
	loader.addClass(targetClass)
	loader.addClass(otherClass)
	referenceField := rtda.NewField(owner, "ref", "Ltest/FieldTarget;", 0x0001, false, 0)
	result = WriteDirectField(rtda.NewExecutionContext(loader), referenceField, rtda.NewObject(owner), rtda.ReferenceValue(rtda.NewObject(otherClass)))
	throwable, ok = result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/ClassCastException" {
		t.Fatalf("reference mismatch result = %#v, throwable=%v; want ClassCastException", result, ok)
	}
	missingField := rtda.NewField(owner, "missing", "Ltest/Missing;", 0x0001, false, 0)
	result = WriteDirectField(rtda.NewExecutionContext(loader), missingField, rtda.NewObject(owner), rtda.ReferenceValue(rtda.NewObject(targetClass)))
	throwable, ok = result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/NoClassDefFoundError" {
		t.Fatalf("reference resolution result = %#v, throwable=%v; want NoClassDefFoundError", result, ok)
	}
}

func TestInvokeDispatchCallerAccess(t *testing.T) {
	loader := newTestLoader()
	setupBaseClasses(loader)
	loader.addClass(rtda.NewSyntheticClass("java/lang/IllegalAccessError", nil))
	owner := rtda.NewSyntheticClass("access/Owner", nil)
	samePackage := rtda.NewSyntheticClass("access/Peer", nil)
	outsider := rtda.NewSyntheticClass("other/Outsider", nil)
	subclass := rtda.NewSyntheticClass("other/Subclass", owner)
	for _, class := range []*rtda.Class{owner, samePackage, outsider, subclass} {
		loader.addClass(class)
	}
	private := rtda.InterpretedMethod(owner, "privateValue", "()I", 0x000a, 1, 0,
		[]byte{byte(opcode.Iconst1), byte(opcode.Ireturn)}, nil)
	packageValue := rtda.InterpretedMethod(owner, "packageValue", "()I", 0x0008, 1, 0,
		[]byte{byte(opcode.Iconst2), byte(opcode.Ireturn)}, nil)
	protected := rtda.InterpretedMethod(owner, "protectedValue", "()I", 0x000c, 1, 0,
		[]byte{byte(opcode.Iconst3), byte(opcode.Ireturn)}, nil)
	owner.AddMethod(private)
	owner.AddMethod(packageValue)
	owner.AddMethod(protected)

	invoke := func(caller *rtda.Class, name string) rtda.DynamicResult {
		return InvokeDispatch(rtda.InvocationLookupRequest{
			Context: rtda.NewExecutionContext(loader), Caller: caller, Target: owner,
			Kind: rtda.InvokeStatic, Name: name, Descriptor: "()I",
		})
	}
	assertDirectInt(t, invoke(owner, "privateValue"), 1)
	assertIllegalAccess(t, invoke(outsider, "privateValue"))
	assertDirectInt(t, invoke(samePackage, "packageValue"), 2)
	assertIllegalAccess(t, invoke(outsider, "packageValue"))
	assertDirectInt(t, invoke(subclass, "protectedValue"), 3)
	assertIllegalAccess(t, invoke(outsider, "protectedValue"))
	privateField := rtda.NewField(owner, "privateField", "I", 0x0002, false, 0)
	assertIllegalAccess(t, ReadDirectFieldFrom(rtda.NewExecutionContext(loader), outsider, privateField, rtda.NewObject(owner)))
}

func assertIllegalAccess(t *testing.T, result rtda.DynamicResult) {
	t.Helper()
	throwable, ok := result.Throwable()
	if !ok || throwable.Class().Name() != "java/lang/IllegalAccessError" {
		t.Fatalf("access result = %#v, throwable=%v; want IllegalAccessError", result, ok)
	}
}

func assertDirectInt(t *testing.T, result rtda.DynamicResult, want int32) {
	t.Helper()
	value, ok := result.Value()
	got, intOK := value.Int()
	if !ok || !intOK || got != want {
		t.Fatalf("direct result = %#v, normal=%v int=%v; want %d", result, ok, intOK, want)
	}
}
