package native

import "catty/rtda"

// NativeClass builds a synthetic core class for the given internal name, or
// returns nil if name is not one of catty's natively-implemented core classes.
// The classloader consults this before falling back to the classpath: classes
// like java.lang.Object/System/String have no on-disk class file in MVP (we do
// not ship a JRE), so they are synthesized here with native Go methods.
func NativeClass(loader rtda.Loader, name string) *rtda.Class {
	switch name {
	case "java/lang/Object":
		return buildObjectClass(loader)
	case "java/lang/String":
		return buildStringClass(loader)
	case "java/lang/StringBuilder":
		return buildStringBuilderClass(loader)
	case "java/io/PrintStream":
		return buildPrintStreamClass(loader)
	case "java/lang/System":
		return buildSystemClass(loader)
	case "java/lang/Throwable":
		return buildThrowable(loader)
	case "java/lang/Exception":
		return buildException(loader)
	case "java/lang/RuntimeException":
		return buildRuntimeException(loader)
	case "java/lang/NullPointerException":
		return buildNPE(loader)
	case "java/lang/ArithmeticException":
		return buildArithmeticException(loader)
	case "java/lang/IndexOutOfBoundsException":
		return buildIndexOutOfBounds(loader)
	case "java/lang/ArrayIndexOutOfBoundsException":
		return buildAIOOBE(loader)
	case "java/lang/ClassCastException":
		return buildCCE(loader)
	case "java/lang/IllegalArgumentException":
		return buildIAE(loader)
	case "java/lang/Error":
		return buildExceptionSubclass("java/lang/Error", "java/lang/Throwable", loader)
	case "java/lang/LinkageError":
		return buildExceptionSubclass("java/lang/LinkageError", "java/lang/Error", loader)
	case "java/lang/IncompatibleClassChangeError":
		return buildExceptionSubclass("java/lang/IncompatibleClassChangeError", "java/lang/LinkageError", loader)
	case "java/lang/NoSuchMethodError":
		return buildExceptionSubclass("java/lang/NoSuchMethodError", "java/lang/IncompatibleClassChangeError", loader)
	case "java/lang/Comparable":
		return buildInterface("java/lang/Comparable", loader)
	case "java/lang/Class":
		return buildClass(loader)
	case "java/lang/Thread":
		return buildThread(loader)
	}
	return nil
}

// nop is the body of native methods that exist only for spec compliance (e.g.
// Object.<init>, which does nothing).
func nop(*rtda.Frame) {}

// buildInterface creates a synthetic interface class. Interfaces have no fields
// and no method bodies — they just declare method signatures. catty treats them
// as empty synthetic classes so the classloader can resolve them.
func buildInterface(name string, loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass(name, loader.LoadClass("java/lang/Object"))
	return c
}

// buildClass creates java.lang.Class as a native class. Class objects store
// the rtda.Class they represent in extra. Native methods (getName, isInterface,
// etc.) are registered in native_registry.go's init().
func buildClass(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/Class", loader.LoadClass("java/lang/Object"))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", nop))
	c.AddMethod(rtda.NativeMethod(c, "desiredAssertionStatus", "()Z", classDesiredAssertionStatus))
	c.AddMethod(rtda.NativeMethod(c, "getName", "()Ljava/lang/String;", classGetName))
	c.AddMethod(rtda.NativeMethod(c, "getSimpleName", "()Ljava/lang/String;", classGetSimpleName))
	c.AddMethod(rtda.NativeMethod(c, "isInterface", "()Z", classIsInterface))
	c.AddMethod(rtda.NativeMethod(c, "isArray", "()Z", classIsArray))
	c.AddMethod(rtda.NativeMethod(c, "getModifiers", "()I", classGetModifiers))
	c.AddMethod(rtda.NativeMethod(c, "isInstance", "(Ljava/lang/Object;)Z", classIsInstance))
	c.AddMethod(rtda.NativeMethod(c, "isAssignableFrom", "(Ljava/lang/Class;)Z", classIsAssignableFrom))
	c.AddMethod(rtda.NativeMethod(c, "getSuperclass", "()Ljava/lang/Class;", classGetSuperclass))
	c.AddMethod(rtda.NativeMethod(c, "isHidden", "()Z", classIsHidden))
	return c
}

// buildThread creates a minimal java.lang.Thread for Thread.currentThread().
func buildThread(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/Thread", loader.LoadClass("java/lang/Object"))
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", nop))
	return c
}
