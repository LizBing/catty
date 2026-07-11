package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"

	"catty/classloader"
	"catty/classpath"
	"catty/interpreter"
	"catty/rtda"
	"catty/transpile"
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
	// `catty build` — offline AOT: transpile → go build → standalone binary.
	if len(os.Args) > 1 && os.Args[1] == "build" {
		build(os.Args[2:])
		return
	}

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

// build implements `catty build [-cp path] [-o output] [-run] <MainClass>`:
// transpiles the program via transpile.BuildProgram, compiles with go build,
// and produces (or runs) a standalone native binary.
func build(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	cp := fs.String("cp", ".", "classpath")
	out := fs.String("o", "catty_aot", "output binary path")
	run := fs.Bool("run", false, "build and immediately run the binary")
	fs.Parse(args)
	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "usage: catty build [-cp path] [-o output] [-run] <main class>")
		os.Exit(2)
	}
	mainClass := rest[0]

	src, err := transpile.BuildProgram(mainClass, *cp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "catty build: %v\n", err)
		os.Exit(1)
	}

	tmp, err := os.MkdirTemp("", "catty-build-*.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "catty build: %v\n", err)
		os.Exit(1)
	}
	srcPath := filepath.Join(tmp, "program.go")
	if err := os.WriteFile(srcPath, []byte(src), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "catty build: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("go", "build", "-o", *out, srcPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "catty build: go build failed: %v\n", err)
		os.Exit(1)
	}

	if *run {
		runCmd := exec.Command(*out)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "catty build: run failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("built: %s\n", *out)
	}
}
