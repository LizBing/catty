package interpreter

import (
	"fmt"

	"catty/rtda"
)

// InvokeDispatch resolves and executes a typed direct call with the tree
// interpreter. Static calls initialize the resolved declarer's Class before
// execution; virtual/interface/special/constructor calls preserve the selected
// receiver and do not imply a separate initialization trigger.
func InvokeDispatch(request rtda.InvocationLookupRequest) rtda.DynamicResult {
	return invokeDispatch(request, InvokeDirect)
}

// InvokeDispatchIR is the IR counterpart to InvokeDispatch.
func InvokeDispatchIR(request rtda.InvocationLookupRequest) rtda.DynamicResult {
	return invokeDispatch(request, InvokeDirectIR)
}

func invokeDispatch(request rtda.InvocationLookupRequest, invoke func(rtda.InvocationRequest) rtda.DynamicResult) (result rtda.DynamicResult) {
	if request.Context == nil {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct dispatch: nil execution context"))
	}
	if request.Context.HasException() {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct dispatch: execution context has a pending exception"))
	}
	if request.Target == nil {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct dispatch: nil target class"))
	}
	if request.Kind != rtda.InvokeStatic && request.Receiver == nil {
		return directJavaThrowable(request.Context, "java/lang/NullPointerException", "direct invocation receiver is null")
	}
	if request.Kind != rtda.InvokeStatic && !request.Receiver.IsInstanceOf(request.Target) {
		return directJavaThrowable(request.Context, "java/lang/IncompatibleClassChangeError", "direct invocation receiver is not assignable to target")
	}
	method := resolveDirectMethod(request)
	if method == nil {
		return directJavaThrowable(request.Context, "java/lang/NoSuchMethodError", request.Target.Name()+"."+request.Name+request.Descriptor)
	}
	if !hasDirectAccess(request.Caller, method.Owner(), method.AccessFlags(), request.Receiver) {
		return directJavaThrowable(request.Context, "java/lang/IllegalAccessError", "direct invocation access denied")
	}
	if request.Kind == rtda.InvokeStatic {
		if !method.IsStatic() {
			return directJavaThrowable(request.Context, "java/lang/IncompatibleClassChangeError", "direct static invocation selected an instance method")
		}
		ensureInitialized(request.Context, method.Owner())
		if request.Context.HasException() {
			return rtda.ThrowableResult(request.Context.ClearException())
		}
	} else if method.IsStatic() {
		return directJavaThrowable(request.Context, "java/lang/IncompatibleClassChangeError", "direct instance invocation selected a static method")
	}
	return invoke(rtda.InvocationRequest{
		Context: request.Context, Caller: request.Caller, Method: method,
		Receiver: request.Receiver, Arguments: request.Arguments,
	})
}

const (
	directAccPublic    uint16 = 0x0001
	directAccPrivate   uint16 = 0x0002
	directAccProtected uint16 = 0x0004
)

func hasDirectAccess(caller, owner *rtda.Class, flags uint16, receiver *rtda.Object) bool {
	if flags&directAccPublic != 0 {
		return true
	}
	if caller == nil || owner == nil {
		return false
	}
	if flags&directAccPrivate != 0 {
		return caller == owner
	}
	samePackage := classPackage(caller) == classPackage(owner)
	if flags&directAccProtected == 0 {
		return samePackage
	}
	if samePackage {
		return true
	}
	if !isSubclassOf(caller, owner) {
		return false
	}
	// JVMS protected cross-package instance access requires the receiver to be
	// an instance of the accessing subclass (or one of its subclasses).
	return receiver == nil || receiver.IsInstanceOf(caller)
}

func isSubclassOf(class, ancestor *rtda.Class) bool {
	for current := class; current != nil; current = current.SuperClass() {
		if current == ancestor {
			return true
		}
	}
	return false
}

func classPackage(class *rtda.Class) string {
	name := class.Name()
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			return name[:i]
		}
	}
	return ""
}

func resolveDirectMethod(request rtda.InvocationLookupRequest) *rtda.Method {
	switch request.Kind {
	case rtda.InvokeConstructor:
		if request.Name != "<init>" {
			return nil
		}
		return request.Target.GetMethod(request.Name, request.Descriptor)
	case rtda.InvokeStatic, rtda.InvokeSpecial:
		return request.Target.LookupMethod(request.Name, request.Descriptor)
	case rtda.InvokeVirtual, rtda.InvokeInterface:
		return request.Receiver.Class().LookupMethod(request.Name, request.Descriptor)
	default:
		return nil
	}
}

// ReadDirectField applies the direct-call initialization trigger to static
// fields before delegating storage access to the typed rtda field boundary.
func ReadDirectField(context *rtda.ExecutionContext, field *rtda.Field, receiver *rtda.Object) rtda.DynamicResult {
	return ReadDirectFieldFrom(context, nil, field, receiver)
}

// ReadDirectFieldFrom reads a field with an explicit Java caller identity.
func ReadDirectFieldFrom(context *rtda.ExecutionContext, caller *rtda.Class, field *rtda.Field, receiver *rtda.Object) rtda.DynamicResult {
	if failure := validateDirectFieldContext(context, field); failure != nil {
		return rtda.InternalFailureResult(failure)
	}
	if !hasDirectAccess(caller, field.Owner(), field.AccessFlags(), receiver) {
		return directJavaThrowable(context, "java/lang/IllegalAccessError", "direct field access denied")
	}
	if field.IsStatic() {
		ensureInitialized(context, field.Owner())
		if context.HasException() {
			return rtda.ThrowableResult(context.ClearException())
		}
	}
	if !field.IsStatic() && receiver == nil {
		return directJavaThrowable(context, "java/lang/NullPointerException", "direct field receiver is null")
	}
	return rtda.ReadFieldValue(field, receiver)
}

// WriteDirectField is the write counterpart to ReadDirectField.
func WriteDirectField(context *rtda.ExecutionContext, field *rtda.Field, receiver *rtda.Object, value rtda.JavaValue) rtda.DynamicResult {
	return WriteDirectFieldFrom(context, nil, field, receiver, value)
}

// WriteDirectFieldFrom writes a field with an explicit Java caller identity.
func WriteDirectFieldFrom(context *rtda.ExecutionContext, caller *rtda.Class, field *rtda.Field, receiver *rtda.Object, value rtda.JavaValue) rtda.DynamicResult {
	if failure := validateDirectFieldContext(context, field); failure != nil {
		return rtda.InternalFailureResult(failure)
	}
	if !hasDirectAccess(caller, field.Owner(), field.AccessFlags(), receiver) {
		return directJavaThrowable(context, "java/lang/IllegalAccessError", "direct field access denied")
	}
	if field.IsStatic() {
		ensureInitialized(context, field.Owner())
		if context.HasException() {
			return rtda.ThrowableResult(context.ClearException())
		}
	}
	if !field.IsStatic() && receiver == nil {
		return directJavaThrowable(context, "java/lang/NullPointerException", "direct field receiver is null")
	}
	if failure, handled := validateDirectReferenceAssignment(context, field, value); handled {
		return failure
	}
	return rtda.WriteFieldValue(field, receiver, value)
}

func validateDirectFieldContext(context *rtda.ExecutionContext, field *rtda.Field) error {
	if context == nil {
		return fmt.Errorf("typed direct field access: nil execution context")
	}
	if context.HasException() {
		return fmt.Errorf("typed direct field access: execution context has a pending exception")
	}
	if field == nil || field.Owner() == nil {
		return fmt.Errorf("typed direct field access: invalid field")
	}
	return nil
}

func directJavaThrowable(context *rtda.ExecutionContext, className, message string) (result rtda.DynamicResult) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = rtda.InternalFailureResult(fmt.Errorf("typed direct Java failure construction: %v", recovered))
		}
	}()
	if context.Loader() == nil {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct Java failure %s: no class loader", className))
	}
	throwRuntime(context, 0, className, message)
	return rtda.ThrowableResult(context.ClearException())
}

func validateDirectReferenceAssignment(context *rtda.ExecutionContext, field *rtda.Field, value rtda.JavaValue) (rtda.DynamicResult, bool) {
	kind, ok := rtda.JavaValueKindForDescriptor(field.Descriptor())
	if !ok || kind != rtda.JavaValueReference || value.Kind() != rtda.JavaValueReference {
		return rtda.DynamicResult{}, false
	}
	reference, _ := value.Reference()
	if reference == nil {
		return rtda.DynamicResult{}, false
	}
	if context.Loader() == nil {
		return rtda.InternalFailureResult(fmt.Errorf("typed direct field reference assignment: no class loader")), true
	}
	className := field.Descriptor()
	if className[0] == 'L' {
		className = className[1 : len(className)-1]
	}
	load := context.Loader().LoadClassResult(className)
	if !load.IsSuccess() {
		failure := load.Failure()
		return directJavaThrowable(context, mapFailureToExceptionClass(failure.Kind), failure.Error()), true
	}
	if !reference.IsInstanceOf(load.Class()) {
		return directJavaThrowable(context, "java/lang/ClassCastException", "direct field reference is not assignable"), true
	}
	return rtda.DynamicResult{}, false
}
