// Command catty runs a Java class file on the Catty interpreter (M0).
//
// Usage: catty run <path/to/Class.class> [main-args…]  (args not yet wired)
package main

import (
	"fmt"
	"os"
	"strings"

	"catty/internal/kernel"
	"catty/internal/vm"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "run" {
		fmt.Fprintln(os.Stderr, "usage: catty run <file.class>")
		os.Exit(2)
	}
	path := os.Args[2]

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "catty:", err)
		os.Exit(2)
	}

	k := kernel.New(kernel.Options{Stdout: os.Stdout})
	cls, err := k.LoadClassBytes(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "catty:", err)
		os.Exit(1)
	}

	t := vm.New(k)
	main, err := k.ResolveMethod(cls, "main", "([Ljava/lang/String;)V")
	if err != nil {
		fmt.Fprintln(os.Stderr, "catty:", err)
		os.Exit(1)
	}
	argsArr, err := k.NewArray("Ljava/lang/String;", 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "catty:", err)
		os.Exit(1)
	}

	if _, err := t.Call(main, nil, []kernel.Value{argsArr}); err != nil {
		var th *kernel.Thrown
		if ok := asThrown(err, &th); ok {
			fmt.Fprintf(os.Stderr, "Exception in thread \"main\" %s\n",
				strings.TrimPrefix(th.Error(), "uncaught "))
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "catty engine error:", err)
		os.Exit(70) // EX_SOFTWARE: engine bug, not a Java-level failure
	}
}

func asThrown(err error, target **kernel.Thrown) bool {
	t, ok := err.(*kernel.Thrown)
	if ok {
		*target = t
	}
	return ok
}
