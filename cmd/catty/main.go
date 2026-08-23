// Command catty runs Java programs on the Catty interpreter.
//
// Usage:
//
//	catty run <path/to/Class.class>
//	catty [-cp <dir>] run <dotted.MainClass>   (directory classpath)
package main

import (
	"fmt"
	"os"
	"strings"

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

	k := kernel.New(kernel.Options{Stdout: os.Stdout})
	var cls *kernel.Class
	var err error
	if strings.HasSuffix(target, ".class") {
		var data []byte
		data, err = os.ReadFile(target)
		if err == nil {
			cls, err = k.LoadClassBytes(data)
		}
	} else if cp != "" {
		loader := kernel.NewClassPathLoader(k, strings.Split(cp, string(os.PathListSeparator)))
		cls, err = loader.Load(dottedToInternal(target))
		if err == nil {
			k.SetResolver(loader.Load)
		}
	} else {
		err = fmt.Errorf("named class requires -cp")
	}
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
	argsArr, err := javaArgs(k, rest[2:])
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

func usage() {
	fmt.Fprintln(os.Stderr, "usage: catty [-cp <dir>] run <File.class | dotted.Main>")
	os.Exit(2)
}

func dottedToInternal(s string) string {
	return strings.ReplaceAll(strings.TrimSuffix(s, ".class"), ".", "/")
}

func asThrown(err error, target **kernel.Thrown) bool {
	t, ok := err.(*kernel.Thrown)
	if ok {
		*target = t
	}
	return ok
}


// javaArgs converts CLI arguments after the main class into a Java
// String[] for main().
func javaArgs(k *kernel.Kernel, cli []string) (kernel.Value, error) {
	arr := &kernel.ArrayObj{Elems: make([]kernel.Value, len(cli))}
	for i, a := range cli {
		arr.Elems[i] = k.InternGo(a)
	}
	return arr, nil
}
