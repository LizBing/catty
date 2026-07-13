package native

import (
	"fmt"
	"io"
	"unicode/utf8"

	"catty/rtda"
)

func init() {
	registerSynthetic("java/io/PrintStream", buildPrintStreamClass)
}

// buildPrintStreamClass makes java.io.PrintStream with the println/print
// overloads catty's test programs use. Each instance's extra is the io.Writer
// it writes to (System.out -> os.Stdout).
func buildPrintStreamClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/io/PrintStream", loader.LoadClass("java/lang/Object"))
	c.AddMethod(rtda.NativeMethod(c, "println", "(Ljava/lang/String;)V", printlnString))
	c.AddMethod(rtda.NativeMethod(c, "println", "(I)V", printlnInt))
	c.AddMethod(rtda.NativeMethod(c, "println", "(J)V", printlnLong))
	c.AddMethod(rtda.NativeMethod(c, "println", "(Z)V", printlnBool))
	c.AddMethod(rtda.NativeMethod(c, "println", "(C)V", printlnChar))
	c.AddMethod(rtda.NativeMethod(c, "println", "()V", printlnEmpty))
	c.AddMethod(rtda.NativeMethod(c, "print", "(Ljava/lang/String;)V", printString))
	c.AddMethod(rtda.NativeMethod(c, "print", "(I)V", printInt))
	return c
}

func writer(f *rtda.Frame) io.Writer {
	if w, ok := f.GetRef(0).Extra().(io.Writer); ok {
		return w
	}
	return io.Discard
}

func printlnString(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(1))
	fmt.Fprintln(writer(f), sv.GoString())
}

func printlnInt(f *rtda.Frame) {
	fmt.Fprintln(writer(f), f.GetInt(1))
}

func printlnLong(f *rtda.Frame) {
	// long takes two local slots: locals[1] (high) and locals[2] (low).
	fmt.Fprintln(writer(f), f.GetLong(1))
}

func printlnBool(f *rtda.Frame) {
	if f.GetInt(1) != 0 {
		fmt.Fprintln(writer(f), "true")
	} else {
		fmt.Fprintln(writer(f), "false")
	}
}

// printlnChar prints the int argument as a Unicode code point, matching Java's
// println(char) which prints the character itself (not its numeric value).
func printlnChar(f *rtda.Frame) {
	ch := rune(f.GetInt(1))
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], ch)
	fmt.Fprintln(writer(f), string(buf[:n]))
}

func printlnEmpty(f *rtda.Frame) {
	fmt.Fprintln(writer(f))
}

func printString(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(1))
	fmt.Fprint(writer(f), sv.GoString())
}

func printInt(f *rtda.Frame) {
	fmt.Fprint(writer(f), f.GetInt(1))
}
