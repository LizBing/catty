package native

import "catty/rtda"

// builderFunc constructs a synthetic class.
type builderFunc func(loader rtda.Loader) *rtda.Class

// syntheticClasses maps internal class names to their builder functions.
// Populated by init() across native/*.go; read by SyntheticProvider (classloader)
// and by the existing NativeClass() fallback for backward compatibility.
var syntheticClasses = map[string]builderFunc{}

// BootstrapClasses is the set of class names that MUST be served from the
// synthetic registry, even when a real JDK is on the classpath. These classes
// form the irreducible Go↔Java bridge and cannot be replaced by bytecode.
var BootstrapClasses = map[string]bool{
	"java/lang/Object":    true,
	"java/lang/String":    true,
	"java/lang/Class":     true,
	"java/lang/System":    true,
	"java/lang/Thread":    true,
	"java/lang/Throwable": true,
}

// IsBootstrap reports whether name is one of the bootstrap classes.
func IsBootstrap(name string) bool { return BootstrapClasses[name] }

// registerSynthetic registers a builder for a synthetic class. Called from
// init() blocks so a builder that references any unexported Go function can
// live in the same file as that function.
func registerSynthetic(name string, fn builderFunc) {
	syntheticClasses[name] = fn
}

// NativeClass builds a synthetic core class, or returns nil. Kept for
// backward compatibility; the primary lookup path is now SyntheticClasses().
func NativeClass(loader rtda.Loader, name string) *rtda.Class {
	if fn := syntheticClasses[name]; fn != nil {
		return fn(loader)
	}
	return nil
}

// SyntheticClasses returns the full synthetic-class registry for the new
// classloader provider chain.
func SyntheticClasses() map[string]builderFunc { return syntheticClasses }

// --- shared helpers for synthetic class builders ---

// nop is the body of native methods that exist only for spec compliance (e.g.
// Object.<init>, which does nothing).
func nop(*rtda.Frame) {}

// staticNative creates a NativeMethod and marks it as static.
func staticNative(owner *rtda.Class, name, descriptor string, fn func(*rtda.Frame)) *rtda.Method {
	m := rtda.NativeMethod(owner, name, descriptor, fn)
	m.SetStatic()
	return m
}

// buildInterface creates a synthetic interface class. Interfaces have no fields
// and no method bodies — they just declare method signatures. catty treats them
// as empty synthetic classes so the classloader can resolve them.
func buildInterface(name string, loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass(name, loader.LoadClass("java/lang/Object"))
	return c
}

func init() {
	registerSynthetic("java/lang/Class", buildClass)
	registerSynthetic("java/lang/Thread", buildThread)
}

// buildClass creates java.lang.Class as a native class.
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
	c.AddMethod(staticNative(c, "getPrimitiveClass", "(Ljava/lang/String;)Ljava/lang/Class;", classGetPrimitiveClass))
	c.AddMethod(staticNative(c, "registerNatives", "()V", nop))
	return c
}

// buildThread creates java.lang.Thread as a full synthetic class with lifecycle,
// interrupt, daemon, and join support (ADR-0028, Slice B).
func buildThread(loader rtda.Loader) *rtda.Class {
	c := rtda.NewSyntheticClass("java/lang/Thread", loader.LoadClass("java/lang/Object"))
	// Constructor — creates the runtime execution context.
	c.AddMethod(rtda.NativeMethod(c, "<init>", "()V", threadInit))
	// run() is a native nop — subclasses override it with bytecode.
	c.AddMethod(rtda.NativeMethod(c, "run", "()V", nop))
	// Lifecycle
	c.AddMethod(rtda.NativeMethod(c, "start", "()V", threadStart))
	c.AddMethod(rtda.NativeMethod(c, "isAlive", "()Z", threadIsAlive))
	c.AddMethod(rtda.NativeMethod(c, "join", "()V", threadJoin))
	// Interrupt
	c.AddMethod(rtda.NativeMethod(c, "interrupt", "()V", threadInterrupt))
	c.AddMethod(rtda.NativeMethod(c, "isInterrupted", "()Z", threadIsInterrupted))
	c.AddMethod(staticNative(c, "interrupted", "()Z", threadInterrupted))
	// Daemon
	c.AddMethod(rtda.NativeMethod(c, "setDaemon", "(Z)V", threadSetDaemon))
	c.AddMethod(rtda.NativeMethod(c, "isDaemon", "()Z", threadIsDaemon))
	// Static utilities
	c.AddMethod(staticNative(c, "sleep", "(J)V", threadSleep))
	c.AddMethod(staticNative(c, "onSpinWait", "()V", threadOnSpinWait))
	// Native registration stubs (also registered in system.go for real-JDK builds)
	c.AddMethod(staticNative(c, "currentThread", "()Ljava/lang/Thread;", threadCurrentThread))
	c.AddMethod(staticNative(c, "registerNatives", "()V", nop))
	return c
}
