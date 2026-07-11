package rtda

// Loader is the subset of the classloader that rtda needs at run time for class
// resolution (new, anewarray, checkcast, ldc class, invokeinterface, ...).
// Declaring it here as an interface keeps rtda free of any import cycle with the
// classloader package, which implements Loader concretely.
type Loader interface {
	LoadClass(name string) *Class
}

// Thread models a single JVM execution thread's stack of frames (JVMS §2.5.1).
// The MVP is single-threaded: the program runs on one Thread, the interpreter
// pushing a frame per method call and popping on return. The concurrency arc
// later promotes Thread to a goroutine backed by the Go scheduler.
type Thread struct {
	stack  []*Frame
	loader Loader
	// bridgeReturn captures a method's return value when run from the AOT bridge
	// (interpreter.RunMethod): there is no caller frame, so the return helpers
	// write here instead of pushing. nil outside bridge mode.
	bridgeReturn *Slot
}

func NewThread(loader Loader) *Thread {
	return &Thread{loader: loader}
}

func (t *Thread) Loader() Loader { return t.loader }

func (t *Thread) PushFrame(frame *Frame) {
	t.stack = append(t.stack, frame)
}

func (t *Thread) PopFrame() *Frame {
	n := len(t.stack)
	f := t.stack[n-1]
	t.stack[n-1] = nil // let the frame (and its slots' refs) be GC'd
	t.stack = t.stack[:n-1]
	return f
}

func (t *Thread) CurrentFrame() *Frame {
	return t.stack[len(t.stack)-1]
}

func (t *Thread) IsStackEmpty() bool {
	return len(t.stack) == 0
}

// Bridge-mode accessors: used by the AOT bridge (interpreter.RunMethod) to capture
// a method's return when there is no caller frame to push it onto.
func (t *Thread) SetBridgeReturn(s *Slot) { t.bridgeReturn = s }
func (t *Thread) HasBridgeReturn() bool    { return t.bridgeReturn != nil }
func (t *Thread) BridgeReturn(s Slot)      { *t.bridgeReturn = s }

// NewFrame allocates a frame for a method on this thread.
func (t *Thread) NewFrame(method *Method) *Frame {
	return NewFrame(t, method)
}
