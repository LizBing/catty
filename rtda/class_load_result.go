package rtda

import "fmt"

// FailureKind categorises a class loading or linking failure.
type FailureKind int

const (
	// FailureNotFound means no provider could locate the class.
	FailureNotFound FailureKind = iota

	// FailureFormat means the classfile bytes are structurally invalid.
	FailureFormat

	// FailureCircularity means a class hierarchy cycle was detected.
	FailureCircularity

	// FailureDuplicateDefinition means a second definition was attempted
	// for a name already defined by this loader.
	FailureDuplicateDefinition

	// FailureLinkage means a superclass, interface, or component type
	// failed to load or link, or the class fails structural constraints.
	FailureLinkage
)

func (k FailureKind) String() string {
	switch k {
	case FailureNotFound:
		return "not found"
	case FailureFormat:
		return "format error"
	case FailureCircularity:
		return "circular hierarchy"
	case FailureDuplicateDefinition:
		return "duplicate definition"
	case FailureLinkage:
		return "linkage failure"
	default:
		return fmt.Sprintf("unknown failure(%d)", k)
	}
}

// ClassLoadFailure is an immutable description of a class loading or linking
// failure. The Name field always holds the symbolic class name that was being
// resolved when the failure occurred. Cause is the underlying Go error (may
// be nil when the failure is purely structural).
type ClassLoadFailure struct {
	Kind  FailureKind
	Name  string
	Cause error
}

func (f *ClassLoadFailure) Error() string {
	if f.Cause != nil {
		return fmt.Sprintf("catty: %s: %s: %v", f.Kind, f.Name, f.Cause)
	}
	return fmt.Sprintf("catty: %s: %s", f.Kind, f.Name)
}

// ClassLoadResult carries exactly one fully linked Class or one terminal
// ClassLoadFailure. A zero value is invalid — construct via NewClassResult
// or NewFailureResult.
type ClassLoadResult struct {
	class   *Class
	failure *ClassLoadFailure
}

// NewClassResult wraps a successfully loaded and linked Class.
func NewClassResult(class *Class) ClassLoadResult {
	if class == nil {
		panic("catty: NewClassResult with nil class")
	}
	return ClassLoadResult{class: class}
}

// NewFailureResult wraps a terminal ClassLoadFailure.
func NewFailureResult(failure *ClassLoadFailure) ClassLoadResult {
	if failure == nil {
		panic("catty: NewFailureResult with nil failure")
	}
	return ClassLoadResult{failure: failure}
}

// IsSuccess reports whether the result is a fully linked Class.
func (r ClassLoadResult) IsSuccess() bool { return r.class != nil }

// IsFailure reports whether the result is a terminal failure.
func (r ClassLoadResult) IsFailure() bool { return r.failure != nil }

// Class returns the loaded Class. Panics if the result is a failure.
func (r ClassLoadResult) Class() *Class {
	if r.class == nil {
		panic("catty: Class() called on failure result")
	}
	return r.class
}

// ClassOrNil returns the loaded Class, or nil on failure.
func (r ClassLoadResult) ClassOrNil() *Class { return r.class }

// Failure returns the terminal failure. Panics if the result is a success.
func (r ClassLoadResult) Failure() *ClassLoadFailure {
	if r.failure == nil {
		panic("catty: Failure() called on success result")
	}
	return r.failure
}

// FailureOrNil returns the failure, or nil on success.
func (r ClassLoadResult) FailureOrNil() *ClassLoadFailure { return r.failure }
