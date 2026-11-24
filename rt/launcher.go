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
	if r.lifecycleEnabled() {
		r.emitAppStart(t, args)
	}
	if tr := r.trace; tr != nil {
		tr.threadStart(t.ID, 0, false, "main")
	}
	out := r.launchMain(t, args)
	if out.Thrown != nil && r.lifecycleEnabled() {
		// The natural-shutdown status must be published, and the diagnostic
		// written, before main retires and the first-zero transition commits.
		r.printUncaughtFor("main", out.Thrown)
		r.thrMu.Lock()
		r.mainStatus = 1
		r.thrMu.Unlock()
	}
	// The canonical main platform thread terminates when main returns or
	// throws; joiners and isAlive observers of the main peer must see that.
	// Publish the thread's last trace event before termination publication can
	// commit a natural shutdown and the terminal process_exit record.
	result := "terminated"
	if out.Thrown != nil {
		result = "uncaught"
	}
	if tr := r.trace; tr != nil {
		tr.threadEnd(t.ID, result, out.Thrown)
	}
	r.publishTermination(t, t.peer, false)
	return out
}

// emitAppStart records the one application start record before the first Java
// entry. Arguments are counted, never recorded.
func (r *Runtime) emitAppStart(t *Thread, args []Value) {
	r.appStartOnce.Do(func() {
		argc := 0
		if len(args) > 0 && args[0].Kind == KindRef && args[0].Ref != nil {
			if array := args[0].Ref.Array; array != nil {
				argc = array.Len()
			}
		}
		if tr := r.trace; tr != nil {
			tr.appStart(r.Image.Profile, r.Image.ABI, argc, t.ID)
		}
	})
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
	// JVMS 5.5: the main class is initialized before main is invoked, so class
	// initialization output and any exit from <clinit> precede the main body.
	if thrown := t.EnsureInitialized(main.Owner); thrown != nil {
		return ThrowOutcome(thrown)
	}
	return t.Call(main, nil, args)
}

// RunApp installs an embedded image and launches it as a standalone process.
// It returns the process exit status: the selected lifecycle status on a Gate 5
// image, zero on a normal Java main return, one for an uncaught Java throwable,
// and two for a host-side installation or diagnostics setup error.
func RunApp(img *Image, args []string) int {
	return runApp(img, args, nil)
}

// RunAppConfigured is RunApp with a pre-install runtime configuration seam. A
// Go subprocess fixture uses it to install a host-fault probe.
func RunAppConfigured(img *Image, args []string, configure func(*Runtime)) int {
	return runApp(img, args, configure)
}

func runApp(img *Image, args []string, configure func(*Runtime)) int {
	r := NewRuntime(img)
	if configure != nil {
		configure(r)
	}
	if path := os.Getenv("CATTY_TRACE"); path != "" {
		if err := r.SetTracePath(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
	}
	defer r.CloseTrace()
	// Metrics are a Gate 5 feature. Open the sink before Install so class and
	// method publication counters include bootstrap and application
	// definitions; the profile/ABI pair is checked again by Install.
	lifecycleImage := false
	if profile, ok := ProfileByID(img.Profile); ok && profile.Lifecycle && profile.ABI == img.ABI {
		lifecycleImage = true
	}
	if lifecycleImage {
		if path := os.Getenv("CATTY_METRICS"); path != "" {
			if err := r.SetMetricsPath(path); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 2
			}
			defer r.CloseMetrics()
		}
	}
	if err := r.Install(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	argv := r.NewStringArray(args)
	if !r.profile.Lifecycle {
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

	// Gate 5: the canonical main logical thread executes on its own goroutine
	// so the launcher can act on a terminal result even while main is running,
	// blocked, or parked in a successful exit/halt call.
	mainThread := r.MainThread()
	go func() {
		defer func() {
			if value := recover(); value != nil {
				r.reportHostFault(mainThread, value)
				panic(value)
			}
		}()
		r.Launch([]Value{RefValue(argv)})
	}()
	result := <-r.lifecycleTerminal
	// The terminal record is the last trace event: process_exit seals the sink
	// even when ordinary Java threads are still live. Metrics are a separate
	// artifact and never go to application stdout.
	if tr := r.trace; tr != nil {
		tr.processExit(result.reason, result.status, result.requester)
	}
	if err := r.finishMetrics(result); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return int(result.status)
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
	fmt.Fprintf(&b, "Exception in thread \"%s\" %s", name, className)
	if data := ThrowableOf(thrown); data != nil && data.Message != "" {
		fmt.Fprintf(&b, ": %s", data.Message)
	}
	b.WriteByte('\n')
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
