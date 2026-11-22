package rt

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Launch runs image startups and the application entry point on the launcher's
// single logical thread. args are the Java method arguments; for a Java main
// method that is one reference value denoting String[].
func (r *Runtime) Launch(args []Value) Outcome {
	if !r.installed {
		return Outcome{Value: VoidValue()}
	}
	t := r.MainThread()
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
		printUncaught(os.Stderr, out.Thrown)
		return 1
	}
	return 0
}

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
