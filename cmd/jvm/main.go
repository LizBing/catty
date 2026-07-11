package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"catty/classloader"
	"catty/classpath"
	"catty/interpreter"
	"catty/rtda"
)

// main is the catty launcher: `jvm [-cp path] [-ir] <main class>`.
// It loads the named class, finds its main([Ljava/lang/String;)V, and hands a
// fresh thread + main frame to the interpreter. Unhandled VM panics (e.g.
// unsupported opcodes, NullPointerException) print a diagnostic and exit 1.
//
// -ir selects the stack-eliminated IR executor (lowering.Lower → LoopIR) instead
// of the default tree-walking interpreter. Both must produce identical output;
// tests/run.sh runs both and diffs against java.
func main() {
	cpOpt := flag.String("cp", ".", "classpath (colon-separated directories/jars)")
	useIR := flag.Bool("ir", false, "use the lowered IR executor instead of the tree-walking interpreter")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: jvm [-cp path] [-ir] <main class>")
		os.Exit(2)
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "catty:", r)
			if os.Getenv("CATTY_DEBUG") != "" {
				debug.PrintStack()
			}
			os.Exit(1)
		}
	}()

	loader := classloader.New(classpath.Parse(*cpOpt))
	class := loader.LoadClass(args[0])
	mainMethod := class.GetMethod("main", "([Ljava/lang/String;)V")
	if mainMethod == nil {
		fmt.Fprintf(os.Stderr, "catty: main method not found in %s\n", args[0])
		os.Exit(1)
	}

	thread := rtda.NewThread(loader)
	// Push the main frame first so that <clinit> (if the main class has one) is
	// pushed on top of it and therefore runs before main — matching the JVMS
	// initialization order, while keeping a caller frame on the stack for the
	// clinit invocation's argument copy.
	frame := thread.NewFrame(mainMethod)
	frame.SetRef(0, nil) // args = null (programs that read args are out of MVP scope)
	thread.PushFrame(frame)
	interpreter.InitClass(thread, class)
	if *useIR {
		interpreter.LoopIR(thread)
	} else {
		interpreter.Loop(thread)
	}
}
