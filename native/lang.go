package native

import (
	"fmt"
	"os"

	"catty/rtda"
)

// buildObjectClass makes java.lang.Object. It has no superclass. Its constructor
// is a nop; getClass returns a stub Class object (catty does not model
// java.lang.Class richly for MVP, so most reflection paths are out of scope).
func buildObjectClass(_ rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/Object", nil)
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", nop))
	c.AddMethod(rtda.NativeMethod(c, "hashCode", "()I", objectHashCode))
	c.AddMethod(rtda.NativeMethod(c, "getClass", "()Ljava/lang/Class;", objectGetClass))
	c.AddMethod(rtda.NativeMethod(c, "clone", "()Ljava/lang/Object;", objectClone))
	c.AddMethod(rtda.NativeMethod(c, "equals", "(Ljava/lang/Object;)Z", objectEquals))
	c.AddMethod(rtda.NativeMethod(c, "toString", "()Ljava/lang/String;", objectToString))
	return c
}

func objectEquals(f *rtda.Frame) {
	this := f.GetRef(0)
	other := f.GetRef(1)
	if this == other {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func objectToString(f *rtda.Frame) {
	this := f.GetRef(0)
	name := javaToDot(this.Class().Name())
	hash := int32(uintptr(0)) // simplified
	result := name + "@" + itoaHex(hash)
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(result)
	f.PushRef(out)
}

func itoaHex(n int32) string {
	if n == 0 {
		return "0"
	}
	const hex = "0123456789abcdef"
	var b [8]byte
	i := len(b)
	u := uint32(n)
	for u > 0 {
		i--
		b[i] = hex[u&0xf]
		u >>= 4
	}
	return string(b[i:])
}

// buildStringClass makes java.lang.String. Instances created by the `ldc`
// instruction carry their Go string value in extra; these constructors cover the
// few cases where user code calls `new String(...)`.
func buildStringClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/String", loader.LoadClass("java/lang/Object"))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", stringInit))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "(Ljava/lang/String;)V", stringInitString))
	c.AddMethod(rtda.NativeMethod(c, "length", "()I", stringLength))
	return c
}

func stringInit(f *rtda.Frame) {
	f.GetRef(0).SetExtra("")
}

func stringInitString(f *rtda.Frame) {
	this := f.GetRef(0)
	arg := f.GetRef(1)
	this.SetExtra(stringValue(arg))
}

func stringLength(f *rtda.Frame) {
	this := f.GetRef(0)
	f.PushInt(int32(len(stringValue(this))))
}

// stringValue returns the Go string held by a java.lang.String object (or "" if
// the object carries no payload, e.g. it was allocated but never assigned).
func stringValue(obj *rtda.Object) string {
	if obj == nil {
		return ""
	}
	if s, ok := obj.Extra().(string); ok {
		return s
	}
	return ""
}

// buildStringBuilderClass makes java.lang.StringBuilder with append/toString,
// backed by a Go strings.Builder stored in extra.
func buildStringBuilderClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/StringBuilder", loader.LoadClass("java/lang/Object"))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", sbInit))
	c.AddMethod(rtda.NativeMethod(c, "append", "(Ljava/lang/String;)Ljava/lang/StringBuilder;", sbAppendString))
	c.AddMethod(rtda.NativeMethod(c, "append", "(I)Ljava/lang/StringBuilder;", sbAppendInt))
	c.AddMethod(rtda.NativeMethod(c, "append", "(J)Ljava/lang/StringBuilder;", sbAppendLong))
	c.AddMethod(rtda.NativeMethod(c, "toString", "()Ljava/lang/String;", sbToString))
	return c
}

func sbInit(f *rtda.Frame) {
	f.GetRef(0).SetExtra(&stringsBuilder{})
}

func sbAppendString(f *rtda.Frame) {
	this := f.GetRef(0)
	this.Extra().(*stringsBuilder).WriteString(stringValue(f.GetRef(1)))
	f.PushRef(this)
}

func sbAppendInt(f *rtda.Frame) {
	this := f.GetRef(0)
	this.Extra().(*stringsBuilder).WriteString(fmt.Sprintf("%d", f.GetInt(1)))
	f.PushRef(this)
}

// sbAppendLong: long arg occupies locals[1] (high) and locals[2] (low).
func sbAppendLong(f *rtda.Frame) {
	this := f.GetRef(0)
	this.Extra().(*stringsBuilder).WriteString(fmt.Sprintf("%d", f.GetLong(1)))
	f.PushRef(this)
}

func sbToString(f *rtda.Frame) {
	this := f.GetRef(0)
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(this.Extra().(*stringsBuilder).String())
	f.PushRef(out)
}

// buildSystemClass makes java.lang.System with the static `out`/`err` fields,
// each holding a PrintStream whose extra is the underlying io.Writer (os.Stdout
// / os.Stderr). Initialized eagerly here since System has no <clinit> to run.
func buildSystemClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/System", loader.LoadClass("java/lang/Object"))
	ps := loader.LoadClass("java/io/PrintStream")

	out := c.AddStaticField("out", "Ljava/io/PrintStream;")
	errf := c.AddStaticField("err", "Ljava/io/PrintStream;")

	outObj := rtda.NewObject(ps)
	outObj.SetExtra(os.Stdout)
	c.SetStaticRef(out.SlotID(), outObj)

	errObj := rtda.NewObject(ps)
	errObj.SetExtra(os.Stderr)
	c.SetStaticRef(errf.SlotID(), errObj)
	return c
}

// stringsBuilder is a thin wrapper around a growable byte buffer, avoiding an
// import of strings.Builder solely to keep the native package dependency-free
// of stdlib quirks (and letting us swap implementations trivially).
type stringsBuilder struct{ buf []byte }

func (b *stringsBuilder) WriteString(s string) { b.buf = append(b.buf, s...) }
func (b *stringsBuilder) String() string       { return string(b.buf) }
