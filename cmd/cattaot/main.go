// Command cattaot runs Java programs with emitter-generated method bodies
// taking precedence over the interpreter for classes present in internal/gen
// (emitter-abi.md §5 mixed execution).
//
// Usage:
//
//	cattaot [-cp <dir>] run <dotted.MainClass>
package main

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/gen"
	"catty/internal/kernel"
	"catty/internal/vm"
)

func main() {
	cp := ""
	rest := os.Args[1:]
	if len(rest) > 0 && rest[0] == "-cp" {
		if len(rest) < 2 {
			usage()
		}
		cp = rest[1]
		rest = rest[2:]
	}
	if len(rest) < 2 || rest[0] != "run" {
		usage()
	}
	target := rest[1]
	internalName := strings.ReplaceAll(strings.TrimSuffix(target, ".class"), ".", "/")

	k := kernel.New(kernel.Options{Stdout: os.Stdout})
	th := vm.New(k) // invoker bridge + primordial thread record

	loader := kernel.NewClassPathLoader(k, strings.Split(cp, string(os.PathListSeparator)))
	k.SetResolver(loader.Load)

	cls, err := loader.Load(internalName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cattaot:", err)
		os.Exit(1)
	}

	gen.Install(k) // emitted bodies take precedence over interpreted ones

	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		fmt.Fprintln(os.Stderr, "cattaot:", err)
		os.Exit(1)
	}
	argsArr := javaArgs(k, rest[2:])

	if _, err := th.Call(main, nil, []kernel.Value{argsArr}); err != nil {
		var t *kernel.Thrown
		if tt, ok := err.(*kernel.Thrown); ok {
			t = tt
			// Stack backfill renders the Java frames captured at
			// construction (DEBT-0019 diagnostic infrastructure).
			fmt.Fprint(os.Stderr, kernel.FormatUncaught("main", t))
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "cattaot engine error:", err)
		os.Exit(70)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: cattaot [-cp <dir>] run <dotted.Main>")
	os.Exit(2)
}

// javaArgs converts CLI arguments after the main class into a Java
// String[] for main().
func javaArgs(k *kernel.Kernel, cli []string) kernel.Value {
	arr := &kernel.ArrayObj{Elems: make([]kernel.Value, len(cli))}
	for i, a := range cli {
		arr.Elems[i] = k.InternGo(a)
	}
	return arr
}
