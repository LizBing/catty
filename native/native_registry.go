package native

import (
	"sync"

	"catty/rtda"
)

// nativeMethodRegistry maps (className, methodName, descriptor) → Go function.
// Populated by init() in system.go and other native files. The classloader
// checks this when building a real JDK class with ACC_NATIVE methods — if a
// Go implementation is registered, it replaces the default stub.
var (
	nativeMu       sync.RWMutex
	nativeRegistry = map[string]func(*rtda.Frame){}
)

func nativeKey(className, methodName, descriptor string) string {
	return className + "\x00" + methodName + "\x00" + descriptor
}

// RegisterNative associates a Go implementation with a native method on a
// specific class. Called from init() blocks in native/*.go files.
func RegisterNative(className, methodName, descriptor string, fn func(*rtda.Frame)) {
	nativeMu.Lock()
	defer nativeMu.Unlock()
	nativeRegistry[nativeKey(className, methodName, descriptor)] = fn
}

// GetNative returns the registered Go implementation for a native method, or
// nil if none is registered (the caller uses the default stub instead).
func GetNative(className, methodName, descriptor string) func(*rtda.Frame) {
	nativeMu.RLock()
	defer nativeMu.RUnlock()
	return nativeRegistry[nativeKey(className, methodName, descriptor)]
}

// --- helpers for Class objects ---

// getClassObject returns (or creates) a java.lang.Class object wrapping the
// given rtda.Class. The Class object stores the rtda.Class in its extra field.
func getClassObject(thread *rtda.Thread, cls *rtda.Class) *rtda.Object {
	classClass := thread.Loader().LoadClass("java/lang/Class")
	obj := rtda.NewObject(classClass)
	obj.SetExtra(cls)
	return obj
}

func getClassFromExtra(obj *rtda.Object) *rtda.Class {
	if cls, ok := obj.Extra().(*rtda.Class); ok {
		return cls
	}
	return nil
}

func javaToDot(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			out[i] = '.'
		} else {
			out[i] = s[i]
		}
	}
	return string(out)
}

// unsafePointer returns the pointer of an object as uintptr for hash codes.
// We avoid importing unsafe by using a helper that the rtda package can provide.
func unsafePointer(obj *rtda.Object) uintptr {
	// Use the address of the fields slice header as a proxy — good enough
	// for identity hashing (which only needs to be stable, not cryptographic).
	if obj == nil {
		return 0
	}
	return uintptr(0) // simplified — real hash can be added later
}
