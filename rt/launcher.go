package rt

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Launch runs image startups and the application entry point on the launcher's
// canonical main logical thread. args are the Java method arguments; for a Java
// main method that is one reference value denoting String[].
func (r *Runtime) Launch(args []Value) Outcome {
	if !r.installed {
		return Outcome{Value: VoidValue()}
	}
	t := r.MainThread()
	if tr := r.trace; tr != nil {
		tr.threadStart(t.ID, 0, false, "main")
	}
	out := r.launchMain(t, args)
	// The canonical main platform thread terminates when main returns or
	// throws; joiners and isAlive observers of the main peer must see that.
	r.thrMu.Lock()
	r.releaseHeldMonitorsLocked(t)
	if p, ok := peerOfObject(t.peer); ok {
		p.state = threadTerminated
		p.thread = nil
		if t.peer.monitor != nil {
			for _, w := range append([]*waiter(nil), t.peer.monitor.waiters...) {
				r.wakeWaiterLocked(t.peer.monitor, w, waiterNotified)
			}
		}
	}
	r.thrMu.Unlock()
	result := "terminated"
	if out.Thrown != nil {
		result = "uncaught"
	}
	if tr := r.trace; tr != nil {
		tr.threadEnd(t.ID, result, out.Thrown)
	}
	return out
}

func (r *Runtime) launchMain(t *Thread, args []Value) Outcome {
	for _, startup := range r.Image.Startups {
		if startup == nil {
			continue
		}
		if out := startup(t); out.Thrown != nil {
			return out
		}
	}
	if r.Image.MainFunc != nil {
		return r.Image.MainFunc(t, args)
	}
	main, err := r.ResolveMethod(r.Image.Main.Owner, r.Image.Main.Name, r.Image.Main.Descriptor)
	if err != nil {
		return ThrowOutcome(r.NewThrowable(t, "java/lang/NoSuchMethodError"))
	}
	return t.Call(main, nil, args)
}

// RunApp installs an embedded image and launches it as a standalone process.
// It returns the process exit status: zero on a normal Java main return, one
// for an uncaught Java throwable, and two for a host-side installation error.
func RunApp(img *Image, args []string) int {
	r := NewRuntime(img)
	if path := os.Getenv("CATTY_TRACE"); path != "" {
		if err := r.SetTracePath(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
	}
	defer r.CloseTrace()
	if err := r.Install(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	argv := r.NewStringArray(args)
	out := r.Launch([]Value{RefValue(argv)})
	if out.Thrown != nil {
		r.printUncaughtFor("main", out.Thrown)
		r.WaitForNonDaemon()
		return 1
	}
	// After main returns or throws, wait for started non-daemon executions.
	// Daemons may still be running at process exit.
	r.WaitForNonDaemon()
	return 0
}

// printUncaughtFor reports one uncaught Java throwable with the terminating
// thread's Java name and frames. Each diagnostic is written as one unit so
// concurrent worker failures cannot interleave.
func (r *Runtime) printUncaughtFor(name string, thrown *Object) {
	var b strings.Builder
	className := "java.lang.Throwable"
	if thrown != nil && thrown.Class != nil {
		className = strings.ReplaceAll(thrown.Class.Name, "/", ".")
	}
	fmt.Fprintf(&b, "Exception in thread \"%s\" %s\n", name, className)
	if thrown != nil {
		if data := ThrowableOf(thrown); data != nil {
			for i := len(data.Frames) - 1; i >= 0; i-- {
				m := data.Frames[i].Method
				if m == nil || m.Owner == nil {
					continue
				}
				owner := strings.ReplaceAll(m.Owner.Name, "/", ".")
				fmt.Fprintf(&b, "\tat %s.%s%s\n", owner, m.Name, m.Descriptor)
			}
		}
	}
	r.diagMu.Lock()
	_, _ = io.WriteString(os.Stderr, b.String())
	r.diagMu.Unlock()
}

// printUncaught keeps the original single-thread diagnostic helper.
func printUncaught(w io.Writer, thrown *Object) {
	className := "java.lang.Throwable"
	if thrown != nil && thrown.Class != nil {
		className = strings.ReplaceAll(thrown.Class.Name, "/", ".")
	}
	fmt.Fprintf(w, "Exception in thread \"main\" %s\n", className)
	if thrown == nil {
		return
	}
	data := ThrowableOf(thrown)
	if data == nil {
		return
	}
	for i := len(data.Frames) - 1; i >= 0; i-- {
		m := data.Frames[i].Method
		if m == nil || m.Owner == nil {
			continue
		}
		owner := strings.ReplaceAll(m.Owner.Name, "/", ".")
		fmt.Fprintf(w, "\tat %s.%s%s\n", owner, m.Name, m.Descriptor)
	}
}
