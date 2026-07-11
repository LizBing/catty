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
	}
	return nil
}

// nop is the body of native methods that exist only for spec compliance (e.g.
// Object.<init>, which does nothing).
func nop(*rtda.Frame) {}
