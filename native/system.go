package native

import (
	"math"
	"time"

	"catty/rtda"
)

// init registers native methods that real JDK classes declare but catty
// implements in Go. These are looked up by (className, methodName, descriptor)
// when the classloader builds a real JDK class with ACC_NATIVE methods.
func init() {
	// --- System ---
	RegisterNative("java/lang/System", "arraycopy", "(Ljava/lang/Object;ILjava/lang/Object;II)V", systemArrayCopy)
	RegisterNative("java/lang/System", "currentTimeMillis", "()J", systemCurrentTimeMillis)
	RegisterNative("java/lang/System", "nanoTime", "()J", systemNanoTime)
	RegisterNative("java/lang/System", "identityHashCode", "(Ljava/lang/Object;)I", systemIdentityHashCode)

	// --- Object ---
	RegisterNative("java/lang/Object", "hashCode", "()I", objectHashCode)
	RegisterNative("java/lang/Object", "getClass", "()Ljava/lang/Class;", objectGetClass)
	RegisterNative("java/lang/Object", "clone", "()Ljava/lang/Object;", objectClone)

	// --- Class ---
	RegisterNative("java/lang/Class", "getName", "()Ljava/lang/String;", classGetName)
	RegisterNative("java/lang/Class", "getSimpleName", "()Ljava/lang/String;", classGetSimpleName)
	RegisterNative("java/lang/Class", "desiredAssertionStatus", "()Z", classDesiredAssertionStatus)
	RegisterNative("java/lang/Class", "isInterface", "()Z", classIsInterface)
	RegisterNative("java/lang/Class", "isArray", "()Z", classIsArray)
	RegisterNative("java/lang/Class", "getModifiers", "()I", classGetModifiers)

	// --- Thread ---
	RegisterNative("java/lang/Thread", "currentThread", "()Ljava/lang/Thread;", threadCurrentThread)

	// --- Float/Double bit conversion (native in JDK) ---
	RegisterNative("java/lang/Float", "floatToRawIntBits", "(F)I", floatToRawIntBits)
	RegisterNative("java/lang/Float", "intBitsToFloat", "(I)F", intBitsToFloat)
	RegisterNative("java/lang/Double", "doubleToRawLongBits", "(D)J", doubleToRawLongBits)
	RegisterNative("java/lang/Double", "longBitsToDouble", "(J)D", longBitsToDouble)
}

// --- System native methods ---

func systemArrayCopy(f *rtda.Frame) {
	dst := f.GetRef(2)
	dstPos := int(f.GetInt(3))
	src := f.GetRef(0)
	srcPos := int(f.GetInt(1))
	length := int(f.GetInt(4))

	if src == nil || dst == nil {
		// NPE would be thrown by JVM; for now, just return
		return
	}

	srcFields := src.Fields()
	dstFields := dst.Fields()

	// Determine element width (1 for cat-1, 2 for long/double arrays)
	width := 1
	if src.Class().ComponentKind() == 5 || src.Class().ComponentKind() == 6 { // long or double
		width = 2
	}

	for i := length - 1; i >= 0; i-- {
		sIdx := (srcPos + i) * width
		dIdx := (dstPos + i) * width
		for j := 0; j < width; j++ {
			dstFields[dIdx+j] = srcFields[sIdx+j]
		}
	}
}

func systemCurrentTimeMillis(f *rtda.Frame) {
	f.PushLong(time.Now().UnixMilli())
}

func systemNanoTime(f *rtda.Frame) {
	f.PushLong(time.Now().UnixNano())
}

func systemIdentityHashCode(f *rtda.Frame) {
	obj := f.GetRef(0)
	f.PushInt(int32(uintptr(unsafePointer(obj)) & 0x7FFFFFFF))
}

// --- Object native methods ---

func objectHashCode(f *rtda.Frame) {
	obj := f.GetRef(0)
	f.PushInt(int32(uintptr(unsafePointer(obj)) & 0x7FFFFFFF))
}

func objectGetClass(f *rtda.Frame) {
	obj := f.GetRef(0)
	if obj == nil {
		f.PushRef(nil)
		return
	}
	clsObj := getClassObject(f.Thread(), obj.Class())
	f.PushRef(clsObj)
}

func objectClone(f *rtda.Frame) {
	this := f.GetRef(0)
	if this == nil {
		f.PushRef(nil)
		return
	}
	// Shallow clone: copy all fields.
	clone := rtda.NewObject(this.Class())
	copy(clone.Fields(), this.Fields())
	f.PushRef(clone)
}

// --- Class native methods ---

func classGetName(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	name := javaToDot(cls.Name())
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(name)
	f.PushRef(out)
}

func classGetSimpleName(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	name := cls.Name()
	// Extract simple name after last '/'
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(name)
	f.PushRef(out)
}

func classDesiredAssertionStatus(f *rtda.Frame) {
	f.PushInt(0) // assertions disabled
}

func classIsInterface(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	if cls.IsInterface() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func classIsArray(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	if cls.IsArray() {
		f.PushInt(1)
	} else {
		f.PushInt(0)
	}
}

func classGetModifiers(f *rtda.Frame) {
	this := f.GetRef(0)
	cls := getClassFromExtra(this)
	f.PushInt(int32(cls.AccessFlags()))
}

func threadCurrentThread(f *rtda.Frame) {
	threadClass := f.Thread().Loader().LoadClass("java/lang/Thread")
	obj := rtda.NewObject(threadClass)
	f.PushRef(obj)
}

// --- Float/Double bit conversion ---

func floatToRawIntBits(f *rtda.Frame) {
	v := math.Float32frombits(uint32(f.GetInt(0)))
	f.PushInt(int32(math.Float32bits(v)))
}

func intBitsToFloat(f *rtda.Frame) {
	f.PushFloat(math.Float32frombits(uint32(f.GetInt(0))))
}

func doubleToRawLongBits(f *rtda.Frame) {
	bits := uint64(f.GetLong(0))
	f.PushLong(int64(bits))
}

func longBitsToDouble(f *rtda.Frame) {
	bits := uint64(f.GetLong(0))
	f.PushDouble(math.Float64frombits(bits))
}
