package native

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"catty/rtda"
)

func init() {
	registerSynthetic("java/lang/Object", buildObjectClass)
	registerSynthetic("java/lang/String", buildStringClass)
	registerSynthetic("java/lang/StringBuilder", buildStringBuilderClass)
	registerSynthetic("java/lang/System", buildSystemClass)
}

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
	// Static fields referenced by real JDK methods (Integer.toHexString etc).
	compact := c.AddStaticField("COMPACT_STRINGS", "Z")
	c.SetStaticRef(compact.SlotID(), nil)      // Z fields use num, not ref — use SetNum
	c.StaticVars()[compact.SlotID()].SetNum(1) // true
	latin1 := c.AddStaticField("LATIN1", "B")
	c.StaticVars()[latin1.SlotID()].SetNum(0)
	utf16 := c.AddStaticField("UTF16", "B")
	c.StaticVars()[utf16.SlotID()].SetNum(1)
	// Instance field: byte coder (0 = LATIN1).
	c.AddInstanceField("coder", "B")
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", stringInit))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "(Ljava/lang/String;)V", stringInitString))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "([BB)V", stringInitBytes)) // called by JDK Integer/Long toString
	c.AddMethod(rtda.NativeMethod(c, "length", "()I", stringLength))
	c.AddMethod(rtda.NativeMethod(c, "charAt", "(I)C", stringCharAt))
	c.AddMethod(rtda.NativeMethod(c, "equals", "(Ljava/lang/Object;)Z", stringEquals))
	c.AddMethod(rtda.NativeMethod(c, "hashCode", "()I", stringHashCode))
	c.AddMethod(rtda.NativeMethod(c, "isEmpty", "()Z", stringIsEmpty))
	c.AddMethod(rtda.NativeMethod(c, "substring", "(I)Ljava/lang/String;", stringSubstring))
	c.AddMethod(rtda.NativeMethod(c, "substring", "(II)Ljava/lang/String;", stringSubstringII))
	c.AddMethod(rtda.NativeMethod(c, "concat", "(Ljava/lang/String;)Ljava/lang/String;", stringConcat))
	c.AddMethod(rtda.NativeMethod(c, "indexOf", "(I)I", stringIndexOf))
	c.AddMethod(rtda.NativeMethod(c, "startsWith", "(Ljava/lang/String;)Z", stringStartsWith))
	c.AddMethod(rtda.NativeMethod(c, "endsWith", "(Ljava/lang/String;)Z", stringEndsWith))
	c.AddMethod(rtda.NativeMethod(c, "compareTo", "(Ljava/lang/String;)I", stringCompareTo))
	c.AddMethod(rtda.NativeMethod(c, "toCharArray", "()[C", stringToCharArray))
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

// stringInitBytes decodes a byte[] per its coder into a Go string.
// coder=0 is LATIN-1 (each byte → code point), coder=1 is big-endian UTF-16.
func stringInitBytes(f *rtda.Frame) {
	this := f.GetRef(0)
	buf := f.GetRef(1) // byte[]
	coder := f.GetInt(2)
	if buf == nil {
		this.SetExtra("")
		return
	}
	n := buf.ArrayLength()
	if n == 0 {
		this.SetExtra("")
		return
	}
	if coder == 0 {
		// LATIN-1: each byte is a Unicode code point.
		raw := make([]byte, n)
		for i := 0; i < n; i++ {
			raw[i] = byte(buf.ArrayElementSlot(i).Num())
		}
		// Convert Latin-1 bytes to runes for correct Go string.
		runes := make([]rune, n)
		for i, b := range raw {
			runes[i] = rune(b)
		}
		this.SetExtra(string(runes))
	} else {
		// UTF-16 big-endian.
		if n%2 != 0 {
			this.SetExtra("")
			return
		}
		u16 := make([]uint16, n/2)
		for i := 0; i < n; i += 2 {
			hi := uint16(buf.ArrayElementSlot(i).Num())
			lo := uint16(buf.ArrayElementSlot(i + 1).Num())
			u16[i/2] = hi<<8 | lo
		}
		this.SetExtra(string(utf16.Decode(u16)))
	}
}

// runesOf returns the rune slice for a String's Go value.
func runesOf(f *rtda.Frame) string { return stringValue(f.GetRef(0)) }

func stringLength(f *rtda.Frame) {
	f.PushInt(int32(utf8.RuneCountInString(stringValue(f.GetRef(0)))))
}

func stringCharAt(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	idx := f.GetInt(1)
	// Slow path: convert to runes for indexed access.
	r := []rune(s)
	if int(idx) < 0 || int(idx) >= len(r) {
		f.PushInt(0)
		return
	}
	f.PushInt(int32(r[idx]))
}

func stringEquals(f *rtda.Frame) {
	a := stringValue(f.GetRef(0))
	other := f.GetRef(1)
	if other == nil || other.Class().Name() != "java/lang/String" {
		f.PushInt(0)
		return
	}
	b := stringValue(other)
	if a == b {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringHashCode(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	runes := []rune(s)
	h := int32(0)
	for _, c := range runes {
		h = h*31 + c
	}
	f.PushInt(h)
}

func stringIsEmpty(f *rtda.Frame) {
	if len(stringValue(f.GetRef(0))) == 0 {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringSubstring(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	begin := f.GetInt(1)
	runes := []rune(s)
	if int(begin) < 0 || int(begin) > len(runes) {
		begin = int32(len(runes))
	}
	result := string(runes[begin:])
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(result)
	f.PushRef(out)
}

func stringSubstringII(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	begin := f.GetInt(1)
	end := f.GetInt(2)
	runes := []rune(s)
	if int(begin) < 0 {
		begin = 0
	}
	if int(end) > len(runes) {
		end = int32(len(runes))
	}
	if int(begin) > int(end) {
		result := ""
		strClass := f.Thread().Loader().LoadClass("java/lang/String")
		out := rtda.NewObject(strClass)
		out.SetExtra(result)
		f.PushRef(out)
		return
	}
	result := string(runes[begin:end])
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(result)
	f.PushRef(out)
}

func stringConcat(f *rtda.Frame) {
	a := stringValue(f.GetRef(0))
	b := stringValue(f.GetRef(1))
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(a + b)
	f.PushRef(out)
}

func stringIndexOf(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	ch := rune(f.GetInt(1))
	idx := strings.IndexRune(s, ch)
	f.PushInt(int32(idx)) // -1 if not found
}

func stringStartsWith(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	prefix := stringValue(f.GetRef(1))
	if strings.HasPrefix(s, prefix) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringEndsWith(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	suffix := stringValue(f.GetRef(1))
	if strings.HasSuffix(s, suffix) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringCompareTo(f *rtda.Frame) {
	a := stringValue(f.GetRef(0))
	b := stringValue(f.GetRef(1))
	if a < b {
		f.PushInt(-1)
	} else if a > b {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringToCharArray(f *rtda.Frame) {
	s := stringValue(f.GetRef(0))
	runes := []rune(s)
	charClass := f.Thread().Loader().LoadClass("[C")
	arr := rtda.NewArray(charClass, len(runes))
	for i, c := range runes {
		arr.ArrayElementSlot(i).SetNum(c)
	}
	f.PushRef(arr)
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
	c.AddMethod(rtda.NativeMethod(c, "append", "(Z)Ljava/lang/StringBuilder;", sbAppendBool))
	c.AddMethod(rtda.NativeMethod(c, "append", "(C)Ljava/lang/StringBuilder;", sbAppendChar))
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

func sbAppendBool(f *rtda.Frame) {
	this := f.GetRef(0)
	v := f.GetInt(1)
	if v != 0 {
		this.Extra().(*stringsBuilder).WriteString("true")
	} else {
		this.Extra().(*stringsBuilder).WriteString("false")
	}
	f.PushRef(this)
}

func sbAppendChar(f *rtda.Frame) {
	this := f.GetRef(0)
	this.Extra().(*stringsBuilder).WriteString(string(rune(f.GetInt(1))))
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

	// Native methods on System (referenced directly by user code).
	c.AddMethod(staticNative(c, "identityHashCode", "(Ljava/lang/Object;)I", systemIdentityHashCode))
	c.AddMethod(staticNative(c, "getProperty", "(Ljava/lang/String;)Ljava/lang/String;", systemGetProperty))
	return c
}

// knownProperties maps commonly-referenced system properties to their values.
// Unknown keys return null (nil). Expanded as needed.
var knownProperties = map[string]string{
	"line.separator":     "\n",
	"file.separator":     "/",
	"path.separator":     ":",
	"file.encoding":      "UTF-8",
	"java.version":       "25",
	"java.vm.version":    "25",
	"java.vm.name":       "catty",
	"java.home":          ".",
	"user.dir":           ".",
	"user.home":          ".",
	"java.class.path":    ".",
	"java.io.tmpdir":     "/tmp",
	"os.name":            "unknown",
	"os.arch":            "unknown",
	"os.version":         "unknown",
	"java.class.version": "65.0",
}

func systemGetProperty(f *rtda.Frame) {
	key := stringValue(f.GetRef(0))
	if val, ok := knownProperties[key]; ok {
		strClass := f.Thread().Loader().LoadClass("java/lang/String")
		out := rtda.NewObject(strClass)
		out.SetExtra(val)
		f.PushRef(out)
		return
	}
	f.PushRef(nil)
}

// stringsBuilder is a thin wrapper around a growable byte buffer, avoiding an
// import of strings.Builder solely to keep the native package dependency-free
// of stdlib quirks (and letting us swap implementations trivially).
type stringsBuilder struct{ buf []byte }

func (b *stringsBuilder) WriteString(s string) { b.buf = append(b.buf, s...) }
func (b *stringsBuilder) String() string       { return string(b.buf) }
