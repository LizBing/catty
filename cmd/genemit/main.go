// Command genemit translates javac-produced class files into Go source
// committed under internal/gen (emitter-abi.md). Output carries a DO NOT
// EDIT header; regenerate instead of editing.
package main

import (
	"flag"
	"fmt"
	"os"
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
	for _, name := range flag.Args() {
		path := filepath.Join(*cp, filepath.FromSlash(name)+".class")
		data, err := os.ReadFile(path)
		if err != nil {
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
