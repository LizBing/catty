package rtda

// InvocationRequest is the engine-neutral input to a direct typed invocation.
// Caller is retained even where the first Interpreter adapter does not yet
// enforce access policy, so later reflection/Host ABI layers cannot invent a
// different caller-identity boundary. Method is already resolved by the
// consumer; receiver is required exactly for non-static methods.
type InvocationRequest struct {
	Context   *ExecutionContext
	Caller    *Class
	Method    *Method
	Receiver  *Object
	Arguments []JavaValue
}

// InvocationKind describes Java dispatch selection for an unresolved direct
// call. It is separate from InvocationRequest because the latter always carries
// an already-resolved Method.
type InvocationKind uint8

const (
	InvokeStatic InvocationKind = iota
	InvokeVirtual
	InvokeInterface
	InvokeSpecial
	InvokeConstructor
)

// InvocationLookupRequest is the typed input for direct dispatch. Target is
// the symbolic target class already resolved by the caller; Method selection is
// performed by the engine adapter according to Kind.
type InvocationLookupRequest struct {
	Context    *ExecutionContext
	Caller     *Class
	Target     *Class
	Kind       InvocationKind
	Name       string
	Descriptor string
	Receiver   *Object
	Arguments  []JavaValue
}
