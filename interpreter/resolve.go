// Package interpreter — typed class-resolution adapter (K2).
//
// Java-reachable symbolic-resolution paths use resolveClass, which calls
// LoadClassResult and maps typed failures to Java throwables on the thread.
// Internal bootstrap must-load helpers may still use LoadClass (which panics
// on failure) — those paths are reachable only when a catty invariant is
// broken, and a Go panic is the correct response.
package interpreter

import (
	"catty/rtda"
)

// resolveClass resolves a class by name via the Loader's typed LoadClassResult.
//
// On success it returns the fully linked *rtda.Class. On failure it maps the
// typed failure to a Java throwable (NoClassDefFoundError, ClassFormatError,
// ClassCircularityError, or LinkageError), signals it on the thread, and
// returns nil. The caller must check for nil and return immediately so the
// interpreter loop picks up the pending exception.
func resolveClass(thread *rtda.Thread, pc int, name string) *rtda.Class {
	result := thread.Loader().LoadClassResult(name)
	if result.IsSuccess() {
		return result.Class()
	}
	f := result.Failure()
	exClass := mapFailureToExceptionClass(f.Kind)
	message := f.Error()
	throwRuntime(thread, pc, exClass, message)
	return nil
}

// mapFailureToExceptionClass maps a typed class-load failure kind to the
// corresponding Java exception class name (JLS 12.2.1, JVMS §5.3).
func mapFailureToExceptionClass(kind rtda.FailureKind) string {
	switch kind {
	case rtda.FailureNotFound:
		return "java/lang/NoClassDefFoundError"
	case rtda.FailureFormat:
		return "java/lang/ClassFormatError"
	case rtda.FailureCircularity:
		return "java/lang/ClassCircularityError"
	case rtda.FailureDuplicateDefinition:
		return "java/lang/LinkageError"
	case rtda.FailureLinkage:
		return "java/lang/LinkageError"
	default:
		return "java/lang/LinkageError"
	}
}
