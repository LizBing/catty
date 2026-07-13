package native

import (
	"fmt"
	"os"

	"catty/rtda"
)

func init() {
	registerSynthetic("java/lang/Object", buildObjectClass)
	registerSynthetic("java/lang/String", buildStringClass)
	registerSynthetic("java/lang/StringBuilder", buildStringBuilderClass)
	registerSynthetic("java/lang/System", buildSystemClass)
}

// buildObjectClass makes java.lang.Object.
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
	hash := int32(uintptr(0))
	result := name + "@" + itoaHex(hash)
	f.PushRef(newStringFromGo(f.Thread(), result))
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

// buildStringClass makes java.lang.String. The canonical backing is an
// immutable *rtda.StringValue ([]uint16) stored in Object.Extra().
func buildStringClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/String", loader.LoadClass("java/lang/Object"))
	// Static fields referenced by real JDK methods.
	compact := c.AddStaticField("COMPACT_STRINGS", "Z")
	c.StaticCells()[compact.SlotID()].SetInt(1) // true
	latin1 := c.AddStaticField("LATIN1", "B")
	c.StaticCells()[latin1.SlotID()].SetInt(0)
	utf16 := c.AddStaticField("UTF16", "B")
	c.StaticCells()[utf16.SlotID()].SetInt(1)
	// Instance field: byte coder (0 = LATIN1).
	c.AddInstanceField("coder", "B")
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", stringInit))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "(Ljava/lang/String;)V", stringInitString))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "([C)V", stringInitChars))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "([BB)V", stringInitBytes))
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

// --- String native methods ---

func stringInit(f *rtda.Frame) {
	f.GetRef(0).SetExtra(rtda.NewStringValue([]uint16{}))
}

func stringInitString(f *rtda.Frame) {
	this := f.GetRef(0)
	arg := f.GetRef(1)
	if arg == nil {
		throwNPE(f, "Cannot invoke \"String.<init>(String)\" because \"original\" is null")
		return
	}
	this.SetExtra(stringValueSV(arg))
}

// stringInitChars builds a String from a char[] with a defensive copy, per
// JLS String(char[]). Each char array element is treated as a UTF-16 code unit.
func stringInitChars(f *rtda.Frame) {
	this := f.GetRef(0)
	chars := f.GetRef(1) // char[]
	if chars == nil {
		throwNPE(f, "Cannot invoke \"String.<init>(char[])\" because \"value\" is null")
		return
	}
	n := chars.ArrayLength()
	units := make([]uint16, n)
	for i := 0; i < n; i++ {
		units[i] = uint16(chars.Cells()[i].GetInt())
	}
	this.SetExtra(rtda.NewStringValue(units))
}

// stringInitBytes decodes a byte[] per its coder into a UTF-16 StringValue.
// coder=0 is LATIN-1 (each byte → code point → uint16), coder=1 is big-endian UTF-16.
func stringInitBytes(f *rtda.Frame) {
	this := f.GetRef(0)
	buf := f.GetRef(1) // byte[]
	coder := f.GetInt(2)
	if buf == nil {
		this.SetExtra(rtda.NewStringValue([]uint16{}))
		return
	}
	n := buf.ArrayLength()
	if n == 0 {
		this.SetExtra(rtda.NewStringValue([]uint16{}))
		return
	}
	if coder == 0 {
		// LATIN-1: each byte is a Unicode code point.
		units := make([]uint16, n)
		for i := 0; i < n; i++ {
			units[i] = uint16(buf.Cells()[i].GetInt() & 0xFF)
		}
		this.SetExtra(rtda.NewStringValue(units))
	} else {
		// UTF-16 big-endian.
		if n%2 != 0 {
			this.SetExtra(rtda.NewStringValue([]uint16{}))
			return
		}
		units := make([]uint16, n/2)
		for i := 0; i < n; i += 2 {
			hi := uint16(buf.Cells()[i].GetInt())
			lo := uint16(buf.Cells()[i+1].GetInt())
			units[i/2] = hi<<8 | lo
		}
		this.SetExtra(rtda.NewStringValue(units))
	}
}

// stringValueSV returns the StringValue held by a String object, or an empty
// StringValue if the object is nil or carries no payload.
func stringValueSV(obj *rtda.Object) *rtda.StringValue {
	if obj == nil {
		return rtda.NewStringValue([]uint16{})
	}
	if sv, ok := obj.Extra().(*rtda.StringValue); ok {
		return sv
	}
	return rtda.NewStringValue([]uint16{})
}

// throwStringBounds throws a StringIndexOutOfBoundsException with the given
// message and signals it on the thread.
func throwStringBounds(f *rtda.Frame, message string) {
	throwException(f, "java/lang/StringIndexOutOfBoundsException", message)
}

// throwNPE throws a NullPointerException with the given message.
func throwNPE(f *rtda.Frame, message string) {
	throwException(f, "java/lang/NullPointerException", message)
}

// throwException throws a named exception with the given detail message on the
// thread's current frame.
func throwException(f *rtda.Frame, className, message string) {
	thread := f.Thread()
	pc := 0
	if cf := thread.CurrentFrame(); cf != nil {
		pc = cf.PC()
	}
	cls := thread.Loader().LoadClass(className)
	obj := rtda.NewObject(cls)
	if message != "" {
		for c := cls; c != nil; c = c.SuperClass() {
			if mf := c.LookupField("detailMessage", "Ljava/lang/String;"); mf != nil {
				msgObj := newStringFromGo(thread, message)
				obj.Cells()[mf.SlotID()].SetRef(msgObj)
				break
			}
		}
	}
	thread.Throw(obj, pc)
}

func stringLength(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	f.PushInt(int32(sv.Len()))
}

func stringCharAt(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	idx := int(f.GetInt(1))
	if idx < 0 || idx >= sv.Len() {
		throwStringBounds(f, "String index out of range: "+itoaInt(idx))
		return
	}
	f.PushInt(int32(sv.CharAt(idx)))
}

func stringEquals(f *rtda.Frame) {
	a := stringValueSV(f.GetRef(0))
	other := f.GetRef(1)
	if other == nil || other.Class().Name() != "java/lang/String" {
		f.PushInt(0)
		return
	}
	b := stringValueSV(other)
	if a.Equals(b) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringHashCode(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	f.PushInt(sv.HashCode())
}

func stringIsEmpty(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	if sv.IsEmpty() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringSubstring(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	begin := int(f.GetInt(1))
	if begin < 0 || begin > sv.Len() {
		throwStringBounds(f, "String index out of range: "+itoaInt(begin))
		return
	}
	f.PushRef(newStringFromSV(f.Thread(), sv.Substring(begin, sv.Len())))
}

func stringSubstringII(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	begin := int(f.GetInt(1))
	end := int(f.GetInt(2))
	if begin < 0 {
		throwStringBounds(f, "String index out of range: "+itoaInt(begin))
		return
	}
	if end > sv.Len() {
		throwStringBounds(f, "String index out of range: "+itoaInt(end))
		return
	}
	if begin > end {
		throwStringBounds(f, "String index out of range: "+itoaInt(begin-end))
		return
	}
	f.PushRef(newStringFromSV(f.Thread(), sv.Substring(begin, end)))
}

func stringConcat(f *rtda.Frame) {
	a := stringValueSV(f.GetRef(0))
	other := f.GetRef(1)
	if other == nil {
		throwNPE(f, "Cannot invoke \"String.concat(String)\" because \"str\" is null")
		return
	}
	b := stringValueSV(other)
	f.PushRef(newStringFromSV(f.Thread(), a.Concat(b)))
}

func stringIndexOf(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	ch := int(f.GetInt(1))
	f.PushInt(int32(sv.IndexOf(ch)))
}

func stringStartsWith(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	prefixObj := f.GetRef(1)
	if prefixObj == nil {
		throwNPE(f, "Cannot invoke \"String.startsWith(String)\" because \"prefix\" is null")
		return
	}
	prefix := stringValueSV(prefixObj)
	if sv.StartsWith(prefix) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringEndsWith(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	suffixObj := f.GetRef(1)
	if suffixObj == nil {
		throwNPE(f, "Cannot invoke \"String.endsWith(String)\" because \"suffix\" is null")
		return
	}
	suffix := stringValueSV(suffixObj)
	if sv.EndsWith(suffix) {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func stringCompareTo(f *rtda.Frame) {
	a := stringValueSV(f.GetRef(0))
	other := f.GetRef(1)
	if other == nil {
		throwNPE(f, "Cannot invoke \"String.compareTo(String)\" because \"anotherString\" is null")
		return
	}
	b := stringValueSV(other)
	f.PushInt(int32(a.CompareTo(b)))
}

func stringToCharArray(f *rtda.Frame) {
	sv := stringValueSV(f.GetRef(0))
	units := sv.ToCharArray()
	charClass := f.Thread().Loader().LoadClass("[C")
	arr := rtda.NewArray(charClass, len(units))
	for i, u := range units {
		arr.Cells()[i].SetInt(int32(u))
	}
	f.PushRef(arr)
}

// --- String construction helpers ---

// newStringFromSV creates a new java.lang.String backed by sv.
func newStringFromSV(thread *rtda.Thread, sv *rtda.StringValue) *rtda.Object {
	class := thread.Loader().LoadClass("java/lang/String")
	obj := rtda.NewObject(class)
	obj.SetExtra(sv)
	return obj
}

// newStringFromGo creates a new java.lang.String from a Go string by
// converting each rune to UTF-16 code units (validating surrogates as-is).
// Used for exception messages, property values, etc. where the source is
// known-safe ASCII or BMP text.
func newStringFromGo(thread *rtda.Thread, s string) *rtda.Object {
	class := thread.Loader().LoadClass("java/lang/String")
	obj := rtda.NewObject(class)
	units := goStringToUTF16(s)
	obj.SetExtra(rtda.NewStringValue(units))
	return obj
}

// goStringToUTF16 converts a Go string to UTF-16 code units. Each rune < 0x10000
// becomes one unit; supplementary runes become two units (surrogate pair).
func goStringToUTF16(s string) []uint16 {
	if s == "" {
		return []uint16{}
	}
	// Fast path: ASCII.
	ascii := true
	for _, r := range s {
		if r >= 0x80 {
			ascii = false
			break
		}
	}
	if ascii {
		units := make([]uint16, len(s))
		for i, b := range []byte(s) {
			units[i] = uint16(b)
		}
		return units
	}
	var units []uint16
	for _, r := range s {
		if r < 0x10000 {
			units = append(units, uint16(r))
		} else {
			r -= 0x10000
			units = append(units, uint16((r>>10)&0x3FF)+0xD800)
			units = append(units, uint16(r&0x3FF)+0xDC00)
		}
	}
	return units
}

// --- StringBuilder ---

// buildStringBuilderClass makes java.lang.StringBuilder with append/toString,
// backed by a []uint16 buffer stored in extra.
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

// utf16Builder is a growable UTF-16 code-unit buffer.
type utf16Builder struct{ units []uint16 }

func (b *utf16Builder) appendUnits(u []uint16) { b.units = append(b.units, u...) }
func (b *utf16Builder) appendRune(r rune) {
	if r < 0x10000 {
		b.units = append(b.units, uint16(r))
	} else {
		r -= 0x10000
		b.units = append(b.units, uint16((r>>10)&0x3FF)+0xD800, uint16(r&0x3FF)+0xDC00)
	}
}
func (b *utf16Builder) toSV() *rtda.StringValue { return rtda.NewStringValue(b.units) }

func sbInit(f *rtda.Frame) {
	f.GetRef(0).SetExtra(&utf16Builder{})
}

func sbAppendString(f *rtda.Frame) {
	this := f.GetRef(0)
	arg := f.GetRef(1)
	if arg == nil {
		this.Extra().(*utf16Builder).appendUnits(goStringToUTF16("null"))
	} else {
		sv := stringValueSV(arg)
		this.Extra().(*utf16Builder).appendUnits(sv.Units())
	}
	f.PushRef(this)
}

func sbAppendInt(f *rtda.Frame) {
	this := f.GetRef(0)
	s := fmt.Sprintf("%d", f.GetInt(1))
	this.Extra().(*utf16Builder).appendUnits(goStringToUTF16(s))
	f.PushRef(this)
}

func sbAppendLong(f *rtda.Frame) {
	this := f.GetRef(0)
	s := fmt.Sprintf("%d", f.GetLong(1))
	this.Extra().(*utf16Builder).appendUnits(goStringToUTF16(s))
	f.PushRef(this)
}

func sbAppendBool(f *rtda.Frame) {
	this := f.GetRef(0)
	v := f.GetInt(1)
	if v != 0 {
		this.Extra().(*utf16Builder).appendUnits(goStringToUTF16("true"))
	} else {
		this.Extra().(*utf16Builder).appendUnits(goStringToUTF16("false"))
	}
	f.PushRef(this)
}

func sbAppendChar(f *rtda.Frame) {
	this := f.GetRef(0)
	this.Extra().(*utf16Builder).appendRune(rune(f.GetInt(1)))
	f.PushRef(this)
}

func sbToString(f *rtda.Frame) {
	this := f.GetRef(0)
	f.PushRef(newStringFromSV(f.Thread(), this.Extra().(*utf16Builder).toSV()))
}

// --- System ---

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

	c.AddMethod(staticNative(c, "identityHashCode", "(Ljava/lang/Object;)I", systemIdentityHashCode))
	c.AddMethod(staticNative(c, "getProperty", "(Ljava/lang/String;)Ljava/lang/String;", systemGetProperty))
	return c
}

var knownProperties = map[string]string{
	"line.separator":  "\n",
	"file.separator":  "/",
	"path.separator":  ":",
	"file.encoding":   "UTF-8",
	"java.version":    "25",
	"java.vm.version": "25",
	"java.vm.name":    "catty",
	"java.home":       ".",
	"user.dir":        ".",
	"user.home":       ".",
	"java.class.path": ".",
	"java.io.tmpdir":  "/tmp",
	"os.name":         "unknown",
	"os.arch":         "unknown",
	"os.version":      "unknown",
	"java.class.version": "65.0",
}

func systemGetProperty(f *rtda.Frame) {
	keyObj := f.GetRef(0)
	if keyObj == nil {
		throwNPE(f, "Cannot invoke \"System.getProperty(String)\" because \"key\" is null")
		return
	}
	key := stringValueSV(keyObj)
	keyStr := key.GoString()
	if val, ok := knownProperties[keyStr]; ok {
		f.PushRef(newStringFromGo(f.Thread(), val))
		return
	}
	f.PushRef(nil)
}

func itoaInt(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
