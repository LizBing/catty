package vm

import (
	"fmt"
	"io"

	"catty/internal/kernel"
)

// maxDefaultFrames mirrors kernel's default; kept local for clarity.
const maxDefaultFrames = 4096

// initMaxDepth resolves the frame budget once per thread.
func (t *Thread) initMaxDepth() {
	if t.maxDepth == 0 {
		t.maxDepth = t.K.MaxFrames()
	}
}

// runJavaThread is the goroutine body behind Thread.start(): it builds a
// fresh execution context bound to the JThread record, dispatches run(),
// terminates the record and reports uncaught throwables.
func runJavaThread(k *kernel.Kernel, j *kernel.JThread) {
	nt := NewWithID(k, j.Key)
	nt.jobj = j.Obj
	nt.initMaxDepth()

	defer func() {
		j.Terminate()
		if r := recover(); r != nil {
			fmt.Fprintf(k.Stderr(), "fatal: engine panic in thread %s: %v\n", j.Name, r)
		}
	}()

	rm, err := k.ResolveMethod(j.Obj.Class, "run", "()V")
	if err != nil {
		fmt.Fprintf(k.Stderr(), "Exception in thread \"%s\": %v\n", j.Name, err)
		return
	}
	if _, err := nt.Call(rm, j.Obj, nil); err != nil {
		if thrown, ok := err.(*kernel.Thrown); ok {
			if h := k.UncaughtHandler; h != nil {
				h(j, thrown)
			}
			return
		}
		// Engine errors are bugs, but a worker must not take down the
		// process: report loudly and terminate the record.
		fmt.Fprintf(k.Stderr(), "fatal: engine error in thread %s: %v\n", j.Name, err)
	}
}

func defaultUncaughtHandler(k *kernel.Kernel) func(j *kernel.JThread, th *kernel.Thrown) {
	var w io.Writer = k.Stderr()
	return func(j *kernel.JThread, th *kernel.Thrown) {
		// Stack backfill: frames captured at construction ride along on
		// the throwable (DEBT-0019 diagnostic infrastructure).
		fmt.Fprint(w, kernel.FormatUncaught(j.Name, th))
	}
}

// NewWithID creates a thread bound to an existing identity (start() path).
// The kernel's invoker bridge is installed once by the primordial New();
// interpreted dispatch is owner-routed, so per-thread installation is
// neither needed nor race-free.
func NewWithID(k *kernel.Kernel, id uint64) *Thread {
	return &Thread{K: k, id: id, maxDepth: k.MaxFrames()}
}
