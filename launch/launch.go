// Package launch is the runtime startup layer: given a classpath and a main
// class name, it builds the classloader, loads the main class, prepares the
// thread, and runs the interpreter or AOT build. This is the shared entry
// point used by both the CLI (cmd/jvm) and the AOT bridge (catty/runtime).
package launch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"catty/classloader"
	"catty/classpath"
	"catty/interpreter"
	"catty/rtda"
	"catty/transpile"
)

// Interpret loads the main class and runs it through the bytecode interpreter
// (or the IR executor if useIR is set). This is the `catty -cp . Main` path.
func Interpret(cpOpt, mainClass string, useIR bool) {
	loader := classloader.New(classpath.Parse(cpOpt))
	class := loader.LoadClass(mainClass)
	mainMethod := class.GetMethod("main", "([Ljava/lang/String;)V")
	if mainMethod == nil {
		fmt.Fprintf(os.Stderr, "catty: main method not found in %s\n", mainClass)
		os.Exit(1)
	}

	// Set the run loop for spawned threads BEFORE any thread starts.
	// native/thread.go calls this instead of importing interpreter directly,
	// avoiding a native → interpreter → lowering → classloader → native cycle.
	if useIR {
		rtda.DefaultRunLoop = interpreter.LoopIR
	} else {
		rtda.DefaultRunLoop = interpreter.Loop
	}

	// Create the VM supervisor (ADR-0028). Set it before any thread starts
	// so native Thread.start() can access it via rtda.GetVM().
	vm := rtda.NewVM()
	rtda.SetVM(vm)

	// Create the main thread execution context and its canonical Java Thread
	// facade. The main thread starts RUNNABLE (it's already executing).
	thread := rtda.NewThread(loader)
	thread.SetStarted() // CAS NEW → RUNNABLE
	thread.SetMain(true)
	threadClass := loader.LoadClass("java/lang/Thread")
	mainThreadObj := rtda.NewObject(threadClass)
	mainThreadObj.SetExtra(thread)
	thread.SetJavaThread(mainThreadObj)

	// Main thread is non-daemon and keeps the VM alive while it runs.
	vm.ThreadStarted(false)

	frame := thread.NewFrame(mainMethod)
	frame.SetRef(0, nil) // args = null
	thread.PushFrame(frame)
	interpreter.InitClass(thread, class)

	if useIR {
		interpreter.LoopIR(thread)
	} else {
		interpreter.Loop(thread)
	}

	// Main thread has completed. Wait for all remaining non-daemon threads
	// before returning (ADR-0028: VM liveness).
	thread.Terminate()
	vm.ThreadTerminated(false)
	vm.WaitForNonDaemonThreads()
}

// Build transpiles a whole program and produces a standalone native binary.
// This is the `catty build` path.
func Build(cpOpt, mainClass, output string, runAfter bool) {
	src, err := transpile.BuildProgram(mainClass, cpOpt)
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

	cmd := exec.Command("go", "build", "-o", output, srcPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "catty build: go build failed: %v\n", err)
		os.Exit(1)
	}

	if runAfter {
		runPath, err := filepath.Abs(output)
		if err != nil {
			runPath = output
		}
		runCmd := exec.Command(runPath)
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr
		if err := runCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "catty build: run failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("built: %s\n", output)
	}
}
