package native

import (
	"catty/rtda"
)

func init() {
	registerSynthetic("java/lang/Throwable", buildThrowable)
	registerSynthetic("java/lang/Exception", buildException)
	registerSynthetic("java/lang/RuntimeException", buildRuntimeException)
	registerSynthetic("java/lang/NullPointerException", buildNPE)
	registerSynthetic("java/lang/ArithmeticException", buildArithmeticException)
	registerSynthetic("java/lang/ArrayIndexOutOfBoundsException", buildAIOOBE)
	registerSynthetic("java/lang/IndexOutOfBoundsException", buildIndexOutOfBounds)
	registerSynthetic("java/lang/ClassCastException", buildCCE)
	registerSynthetic("java/lang/IllegalArgumentException", buildIAE)
	registerSynthetic("java/lang/Error", func(loader rtda.Loader) *rtda.Class {
		return buildExceptionSubclass("java/lang/Error", "java/lang/Throwable", loader)
	})
	registerSynthetic("java/lang/LinkageError", func(loader rtda.Loader) *rtda.Class {
		return buildExceptionSubclass("java/lang/LinkageError", "java/lang/Error", loader)
	})
	registerSynthetic("java/lang/IncompatibleClassChangeError", func(loader rtda.Loader) *rtda.Class {
		return buildExceptionSubclass("java/lang/IncompatibleClassChangeError", "java/lang/LinkageError", loader)
	})
	registerSynthetic("java/lang/NoSuchMethodError", func(loader rtda.Loader) *rtda.Class {
		return buildExceptionSubclass("java/lang/NoSuchMethodError", "java/lang/IncompatibleClassChangeError", loader)
	})
	registerSynthetic("java/lang/Comparable", buildComparable)
}

func buildComparable(loader rtda.Loader) *rtda.Class {
	return buildInterface("java/lang/Comparable", loader)
}

// buildThrowable creates java.lang.Throwable with a detailMessage field and
// the getMessage/toString methods. All exception classes chain to this.
func buildThrowable(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/Throwable", loader.LoadClass("java/lang/Object"))
	// detailMessage: String — stored as an instance field on Throwable.
	c.AddInstanceField("detailMessage", "Ljava/lang/String;")
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", nop))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "(Ljava/lang/String;)V", throwableInitMsg))
	c.AddMethod(rtda.NativeMethod(c, "getMessage", "()Ljava/lang/String;", throwableGetMessage))
	c.AddMethod(rtda.NativeMethod(c, "toString", "()Ljava/lang/String;", throwableToString))
	return c
}

// detailMessageSlot walks the class hierarchy to find the slot offset of
// Throwable's detailMessage field. This avoids hardcoding the offset.
func detailMessageSlot(obj *rtda.Object) uint {
	for cls := obj.Class(); cls != nil; cls = cls.SuperClass() {
		if f := cls.LookupField("detailMessage", "Ljava/lang/String;"); f != nil {
			return f.SlotID()
		}
	}
	panic("catty: detailMessage field not found")
}

func throwableInitMsg(f *rtda.Frame) {
	this := f.GetRef(0)
	msg := f.GetRef(1)
	slot := detailMessageSlot(this)
	this.Fields()[slot].SetRef(msg)
}

func throwableGetMessage(f *rtda.Frame) {
	this := f.GetRef(0)
	slot := detailMessageSlot(this)
	f.PushRef(this.Fields()[slot].Ref())
}

func throwableToString(f *rtda.Frame) {
	this := f.GetRef(0)
	className := javaClassName(this.Class().Name())
	msgSlot := detailMessageSlot(this)
	msgObj := this.Fields()[msgSlot].Ref()
	msg := ""
	if msgObj != nil {
		if s, ok := msgObj.Extra().(string); ok {
			msg = s
		}
	}
	result := className
	if msg != "" {
		result += ": " + msg
	}
	strClass := f.Thread().Loader().LoadClass("java/lang/String")
	out := rtda.NewObject(strClass)
	out.SetExtra(result)
	f.PushRef(out)
}

// javaClassName converts an internal class name ("java/lang/NullPointerException")
// to the Java dotted form ("java.lang.NullPointerException").
func javaClassName(internal string) string {
	out := make([]byte, len(internal))
	for i := 0; i < len(internal); i++ {
		if internal[i] == '/' {
			out[i] = '.'
		} else {
			out[i] = internal[i]
		}
	}
	return string(out)
}

// buildExceptionSubclass is a shared builder for the exception hierarchy.
// Each subclass inherits Throwable's detailMessage and methods via the class
// hierarchy (LookupMethod walks up). It only needs its own <init> constructors
// that chain to super.
func buildExceptionSubclass(name string, superName string, loader rtda.Loader) *rtda.Class {
	super := loader.LoadClass(superName)
	c := rtda.NewSyntheticClass(name, super)
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", nop))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "(Ljava/lang/String;)V", throwableInitMsg))
	return c
}

func buildException(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/Exception", "java/lang/Throwable", loader)
}

func buildRuntimeException(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/RuntimeException", "java/lang/Exception", loader)
}

func buildNPE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/NullPointerException", "java/lang/RuntimeException", loader)
}

func buildArithmeticException(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/ArithmeticException", "java/lang/RuntimeException", loader)
}

func buildAIOOBE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/ArrayIndexOutOfBoundsException", "java/lang/IndexOutOfBoundsException", loader)
}

func buildIndexOutOfBounds(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/IndexOutOfBoundsException", "java/lang/RuntimeException", loader)
}

func buildCCE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/ClassCastException", "java/lang/RuntimeException", loader)
}

func buildIAE(loader rtda.Loader) *rtda.Class {
	return buildExceptionSubclass("java/lang/IllegalArgumentException", "java/lang/RuntimeException", loader)
}
