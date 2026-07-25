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
