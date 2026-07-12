package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"

	"catty/launch"
)

// main is the catty CLI: `catty [-cp path] [-ir] <main class>` for
// interpretation, or `catty build [-cp path] [-o output] [-run] <main class>`
// for offline AOT compilation. All runtime logic lives in catty/launch; this
// file is only CLI argument parsing + environment setup.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "build" {
		buildCmd(os.Args[2:])
		return
	}

	cpOpt := flag.String("cp", ".", "classpath (colon-separated directories/jars)")
	useIR := flag.Bool("ir", false, "use the lowered IR executor")
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: catty [-cp path] [-ir] <main class>")
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

	launch.Interpret(*cpOpt, args[0], *useIR)
}

// buildCmd implements `catty build`.
func buildCmd(args []string) {
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
	launch.Build(*cp, rest[0], *out, *run)
}
