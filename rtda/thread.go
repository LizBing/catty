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
	// ecID identifies the execution context (ADR-0025). In the single-context
	// interpreter it is a fixed sentinel; AOT and future multi-threaded runtimes
	// assign a distinct value per Java thread so recursive same-owner <clinit>
	// requests complete normally without re-running <clinit>.
	ecID uint64
	// bridgeReturn captures a method's return value when run from the AOT bridge
	// (interpreter.RunMethod): there is no caller frame, so the return helpers
	// write here instead of pushing. nil outside bridge mode.
	bridgeReturn *Slot
	// pendingException is non-nil when an exception is in flight (athrow or a
	// runtime error like NPE). The interpreter Loop checks HasException after
	// each instruction and dispatches to handleException.
	pendingException *Object
	throwPC          int // PC of the instruction that threw (for exception-table search)
}

// threadECSeq is a monotonically increasing counter for execution-context
// identity assignment. It starts at 1 so that 0 is reserved for "no owner".
var threadECSeq uint64 = 1

func NewThread(loader Loader) *Thread {
	ecID := threadECSeq
	threadECSeq++
	return &Thread{loader: loader, ecID: ecID}
}

func (t *Thread) Loader() Loader { return t.loader }
func (t *Thread) EC() uint64     { return t.ecID }

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

// FrameCount returns the current number of frames on the thread's stack.
func (t *Thread) FrameCount() int { return len(t.stack) }

// Bridge-mode accessors: used by the AOT bridge (interpreter.RunMethod) to capture
// a method's return when there is no caller frame to push it onto.
func (t *Thread) SetBridgeReturn(s *Slot) { t.bridgeReturn = s }
func (t *Thread) HasBridgeReturn() bool    { return t.bridgeReturn != nil }
func (t *Thread) BridgeReturn(s Slot)      { *t.bridgeReturn = s }

// --- Exception handling ---
//
// Exceptions use a signal on the Thread (not Go panic/recover): the opcode
// handler sets the pending exception + throwPC, returns from exec, and the
// Loop's handleException searches exception tables frame-by-frame.

func (t *Thread) Throw(obj *Object, pc int) { t.pendingException = obj; t.throwPC = pc }
func (t *Thread) HasException() bool        { return t.pendingException != nil }
func (t *Thread) ClearException() *Object {
	obj := t.pendingException
	t.pendingException = nil
	return obj
}
func (t *Thread) ThrowPC() int { return t.throwPC }

// NewFrame allocates a frame for a method on this thread.
func (t *Thread) NewFrame(method *Method) *Frame {
	return NewFrame(t, method)
}
