package interpreter

import (
	"fmt"

	"catty/lowering"
	"catty/opcode"
	"catty/rtda"
)

// InvokeDirect executes one already-resolved method through the typed dynamic
// invocation boundary. It may run above an ordinary Java caller frame; the
// boundary captures only the return/throwable that crosses its entry depth.
// Caller identity is carried in InvocationRequest for later access policy.
func InvokeDirect(request rtda.InvocationRequest) (result rtda.DynamicResult) {
	return invokeDirect(request, invokeDirectInterpreted)
}

// InvokeDirectIR is the IR counterpart to InvokeDirect. It shares the typed
// request/result and exception boundary while executing interpreted bytecode
// through the lowering IR adapter.
func InvokeDirectIR(request rtda.InvocationRequest) (result rtda.DynamicResult) {
	return invokeDirect(request, invokeDirectIR)
}

func invokeDirect(request rtda.InvocationRequest, interpreted func(*rtda.ExecutionContext, *rtda.Frame, int) rtda.DynamicResult) (result rtda.DynamicResult) {
	if failure := validateDirectRequest(request); failure != nil {
		return rtda.InternalFailureResult(failure)
	}
	context := request.Context
	baseDepth := context.FrameCount()
	defer func() {
		if recovered := recover(); recovered != nil {
			result = rtda.InternalFailureResult(fmt.Errorf("typed direct invocation panic: %v", recovered))
		}
		if result.IsInternalFailure() {
			for context.FrameCount() > baseDepth {
				context.PopFrame()
			}
		}
	}()

	frame, failure := directFrame(request)
	if failure != nil {
		return rtda.InternalFailureResult(failure)
	}
	if request.Method.IsNative() {
		return invokeDirectNative(context, request.Method, frame)
	}
	return interpreted(context, frame, baseDepth)
}

func validateDirectRequest(request rtda.InvocationRequest) error {
	if request.Context == nil {
		return fmt.Errorf("typed direct invocation: nil execution context")
	}
	if request.Context.HasException() {
		return fmt.Errorf("typed direct invocation: execution context has a pending exception")
	}
	if request.Method == nil {
		return fmt.Errorf("typed direct invocation: nil method")
	}
	if request.Method.IsStatic() {
		if request.Receiver != nil {
			return fmt.Errorf("typed direct invocation %s%s: static method has receiver", request.Method.Name(), request.Method.Descriptor())
		}
	} else {
		if request.Receiver == nil {
			return fmt.Errorf("typed direct invocation %s%s: nil receiver", request.Method.Name(), request.Method.Descriptor())
		}
		if owner := request.Method.Owner(); owner != nil && !request.Receiver.IsInstanceOf(owner) {
			return fmt.Errorf("typed direct invocation %s%s: receiver is not assignable to owner", request.Method.Name(), request.Method.Descriptor())
		}
	}
	parameters := rtda.ParseMethodDescriptor(request.Method.Descriptor()).ParameterTypes
	if len(parameters) != len(request.Arguments) {
		return fmt.Errorf("typed direct invocation %s%s: got %d arguments, want %d",
			request.Method.Name(), request.Method.Descriptor(), len(request.Arguments), len(parameters))
	}
	for i, descriptor := range parameters {
		kind, ok := rtda.JavaValueKindForDescriptor(descriptor)
		if !ok || kind == rtda.JavaValueVoid || request.Arguments[i].Kind() != kind {
			return fmt.Errorf("typed direct invocation %s%s: argument %d kind %s, want %s",
				request.Method.Name(), request.Method.Descriptor(), i, request.Arguments[i].Kind(), kind)
		}
		if kind == rtda.JavaValueReference {
			if failure := validateDirectReferenceArgument(request.Context, descriptor, request.Arguments[i]); failure != nil {
				return fmt.Errorf("typed direct invocation %s%s: argument %d: %w",
					request.Method.Name(), request.Method.Descriptor(), i, failure)
			}
		}
	}
	return nil
}

func validateDirectReferenceArgument(context *rtda.ExecutionContext, descriptor string, value rtda.JavaValue) error {
	reference, _ := value.Reference()
	if reference == nil {
		return nil
	}
	if context.Loader() == nil {
		return fmt.Errorf("cannot validate reference assignment without a class loader")
	}
	className := descriptor
	if descriptor[0] == 'L' {
		className = descriptor[1 : len(descriptor)-1]
	}
	load := context.Loader().LoadClassResult(className)
	if !load.IsSuccess() {
		return fmt.Errorf("cannot resolve parameter type %s: %w", descriptor, load.Failure())
	}
	if !reference.IsInstanceOf(load.Class()) {
		return fmt.Errorf("reference is not assignable to parameter type %s", descriptor)
	}
	return nil
}

func directFrame(request rtda.InvocationRequest) (*rtda.Frame, error) {
	frame := request.Context.NewFrame(request.Method)
	index := 0
	if !request.Method.IsStatic() {
		if result := rtda.WriteFrameLocal(frame, index, rtda.ReferenceValue(request.Receiver)); !result.IsNormal() {
			return nil, result.Failure()
		}
		index++
	}
	for _, argument := range request.Arguments {
		if result := rtda.WriteFrameLocal(frame, index, argument); !result.IsNormal() {
			return nil, result.Failure()
		}
		if argument.Kind() == rtda.JavaValueLong || argument.Kind() == rtda.JavaValueDouble {
			index += 2
		} else {
			index++
		}
	}
	return frame, nil
}

func invokeDirectNative(context *rtda.ExecutionContext, method *rtda.Method, frame *rtda.Frame) rtda.DynamicResult {
	frame.EnterSyncMonitor()
	defer func() {
		if lock := frame.SyncObject(); lock != nil {
			lock.Monitor().Exit(context.ID())
		}
	}()
	method.NativeFunc()(frame)
	if context.HasException() {
		return rtda.ThrowableResult(context.ClearException())
	}
	return directFrameResult(frame, method.ReturnType())
}

func invokeDirectInterpreted(context *rtda.ExecutionContext, frame *rtda.Frame, baseDepth int) rtda.DynamicResult {
	var value rtda.JavaValue
	context.PushBridgeDynamicReturn(baseDepth, &value)
	defer context.PopBridgeDynamicReturn()
	context.PushFrame(frame)
	for context.FrameCount() > baseDepth {
		current := context.CurrentFrame()
		pc := current.PC()
		op := opcode.Opcode(current.Code()[pc])
		current.SetPC(pc + 1)
		exec(context, current, op, pc)
		for context.HasException() {
			thrown, uncaught := unwindDirectException(context, pc, baseDepth)
			if uncaught {
				return rtda.ThrowableResult(thrown)
			}
		}
	}
	return rtda.NormalResult(value)
}

func invokeDirectIR(context *rtda.ExecutionContext, frame *rtda.Frame, baseDepth int) rtda.DynamicResult {
	var value rtda.JavaValue
	context.PushBridgeDynamicReturn(baseDepth, &value)
	defer context.PopBridgeDynamicReturn()
	context.PushFrame(frame)
	cache := map[*rtda.Method]*lowering.IR{}
	for context.FrameCount() > baseDepth {
		current := context.CurrentFrame()
		ir := cache[current.Method()]
		if ir == nil {
			var err error
			ir, err = lowering.Lower(current.Method())
			if err != nil {
				return rtda.InternalFailureResult(fmt.Errorf("typed direct IR lowering: %w", err))
			}
			cache[current.Method()] = ir
		}
		pc := current.PC()
		execIR(context, current, ir)
		for context.HasException() {
			thrown, uncaught := unwindDirectException(context, pc, baseDepth)
			if uncaught {
				return rtda.ThrowableResult(thrown)
			}
		}
	}
	return rtda.NormalResult(value)
}

func directFrameResult(frame *rtda.Frame, descriptor string) rtda.DynamicResult {
	kind, ok := rtda.JavaValueKindForDescriptor(descriptor)
	if !ok {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct invocation: invalid return descriptor %q", descriptor))
	}
	switch kind {
	case rtda.JavaValueVoid:
		return rtda.NormalResult(rtda.VoidValue())
	case rtda.JavaValueBoolean:
		return rtda.NormalResult(rtda.BooleanValue(frame.PopInt() != 0))
	case rtda.JavaValueByte:
		return rtda.NormalResult(rtda.ByteValue(int8(frame.PopInt())))
	case rtda.JavaValueChar:
		return rtda.NormalResult(rtda.CharValue(uint16(frame.PopInt())))
	case rtda.JavaValueShort:
		return rtda.NormalResult(rtda.ShortValue(int16(frame.PopInt())))
	case rtda.JavaValueInt:
		return rtda.NormalResult(rtda.IntValue(frame.PopInt()))
	case rtda.JavaValueLong:
		return rtda.NormalResult(rtda.LongValue(frame.PopLong()))
	case rtda.JavaValueFloat:
		return rtda.NormalResult(rtda.FloatValue(frame.PopFloat()))
	case rtda.JavaValueDouble:
		return rtda.NormalResult(rtda.DoubleValue(frame.PopDouble()))
	case rtda.JavaValueReference:
		return rtda.NormalResult(rtda.ReferenceValue(frame.PopRef()))
	default:
		return rtda.InternalFailureResult(fmt.Errorf("typed direct invocation: invalid return kind %s", kind))
	}
}

// unwindDirectException mirrors the interpreter's exception-table walk but
// returns an exception that escapes the direct invocation boundary instead of
// printing it or terminating a Java Thread.
func unwindDirectException(context *rtda.ExecutionContext, throwPC, baseDepth int) (*rtda.Object, bool) {
	thrown := context.ClearException()
forFrames:
	for context.FrameCount() > baseDepth {
		frame := context.CurrentFrame()
		for _, entry := range frame.Method().ExceptionTable() {
			if throwPC < entry.StartPc() || throwPC >= entry.EndPc() {
				continue
			}
			if entry.CatchType() == "" {
				frame.ClearStack()
				frame.PushRef(thrown)
				frame.SetPC(entry.HandlerPc())
				return nil, false
			}
			catchClass := resolveClass(context, throwPC, entry.CatchType())
			if catchClass == nil {
				thrown = context.ClearException()
				context.PopFrame()
				if context.FrameCount() == baseDepth {
					return thrown, true
				}
				throwPC = context.CurrentFrame().PC() - 1
				continue forFrames
			}
			if thrown.IsInstanceOf(catchClass) {
				frame.ClearStack()
				frame.PushRef(thrown)
				frame.SetPC(entry.HandlerPc())
				return nil, false
			}
		}
		context.PopFrame()
		if context.FrameCount() == baseDepth {
			return thrown, true
		}
		throwPC = context.CurrentFrame().PC() - 1
	}
	return thrown, true
}
