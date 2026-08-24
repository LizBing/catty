// Command genemit translates javac-produced class files into Go source
// committed under internal/gen (emitter-abi.md). Output carries a DO NOT
// EDIT header; regenerate instead of editing.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"catty/internal/classfile"
	"catty/internal/emitter"
)

func main() {
	cp := flag.String("cp", "", "directory classpath")
	out := flag.String("out", "", "output .go file")
	flag.Parse()
	if *cp == "" || *out == "" || flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: genemit -cp <dir> -out <file.go> <Class[ Class...]>")
		os.Exit(2)
	}

	var files []*classfile.ClassFile
	dirs := filepath.SplitList(*cp)
	for _, name := range flag.Args() {
		var data []byte
		var err error
		rel := filepath.FromSlash(strings.ReplaceAll(name, ".", "/") + ".class")
		for _, dir := range dirs {
			var d []byte
			d, err = os.ReadFile(filepath.Join(dir, rel))
			if err == nil {
				data = d
				break
			}
		}
		if data == nil {
			fmt.Fprintln(os.Stderr, "genemit:", err)
			os.Exit(1)
		}
		cf, err := classfile.Parse(data)
		if err != nil {
			fmt.Fprintln(os.Stderr, "genemit:", err)
			os.Exit(1)
		}
		files = append(files, cf)
	}

	src, err := emitter.EmitProgram(files...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "genemit:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, []byte(src), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "genemit:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "genemit: wrote %s (%d classes)\n", *out, len(files))
}
