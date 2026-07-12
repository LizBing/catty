package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strings"

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
	noBoot := flag.Bool("no-boot", false, "do not prepend java.base to classpath")
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

	fullCP := *cpOpt
	if !*noBoot {
		if boot := detectBootClasspath(); boot != "" {
			fullCP = boot + string(os.PathListSeparator) + fullCP
		}
	}
	launch.Interpret(fullCP, args[0], *useIR)
}

// detectBootClasspath returns the path to a java.base class directory, or ""
// if no JDK is available. Detection order:
//
//  1. $CATTY_BOOT — explicit path to an extracted java.base directory
//  2. $JAVA_HOME/lib/java.base — pre-extracted
//  3. /usr/libexec/java_home (macOS) — follow symlink to JDK home
//  4. Well-known JDK install paths
func detectBootClasspath() string {
	// 1. Explicit override.
	if dir := os.Getenv("CATTY_BOOT"); dir != "" {
		if isDir(dir) {
			return dir
		}
	}

	// 2. $JAVA_HOME.
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		if p := bootUnder(jh); p != "" {
			return p
		}
	}

	// 3. macOS java_home tool.
	if out, err := exec.Command("/usr/libexec/java_home").Output(); err == nil {
		jh := strings.TrimSpace(string(out))
		if p := bootUnder(jh); p != "" {
			return p
		}
	}

	// 4. Well-known paths.
	for _, jh := range []string{
		"/Library/Java/JavaVirtualMachines/temurin-25.jdk/Contents/Home",
		"/opt/homebrew/opt/openjdk/libexec/openjdk.jdk/Contents/Home",
		"/usr/local/opt/openjdk/libexec/openjdk.jdk/Contents/Home",
		"/usr/lib/jvm/java-25-openjdk",
		"/usr/lib/jvm/java-21-openjdk",
	} {
		if p := bootUnder(jh); p != "" {
			return p
		}
	}

	return ""
}

// bootUnder looks for a java.base class directory under a JDK home.
func bootUnder(javaHome string) string {
	// Pre-extracted directory (from jimage extract --dir=.../java.base).
	if p := filepath.Join(javaHome, "lib", "java.base"); isDir(p) {
		return p
	}
	return ""
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

// buildCmd implements `catty build`.
func buildCmd(args []string) {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	cp := fs.String("cp", ".", "classpath")
	out := fs.String("o", "catty_aot", "output binary path")
	run := fs.Bool("run", false, "build and immediately run the binary")
	noBoot := fs.Bool("no-boot", false, "do not prepend java.base to classpath")
	fs.Parse(args)
	rest := fs.Args()
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "usage: catty build [-cp path] [-o output] [-run] <main class>")
		os.Exit(2)
	}

	fullCP := *cp
	if !*noBoot {
		if boot := detectBootClasspath(); boot != "" {
			fullCP = boot + string(os.PathListSeparator) + fullCP
		}
	}
	launch.Build(fullCP, rest[0], *out, *run)
}
