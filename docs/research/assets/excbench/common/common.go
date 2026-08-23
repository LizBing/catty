// Package common holds the sentinel type and the identical work payload shared
// by both exception-mechanism candidates, so any measured difference is
// attributable to the mechanism and not to divergent arithmetic.
package common

// JException is the sentinel value a Java exception is carried in. Only its
// class tag is modeled, because that is all exception-table matching
// (instanceof) reads. The real runtime payload (message, cause, Java stack
// trace) is orthogonal to the propagation mechanism and omitted.
type JException struct {
	Tag uint32
}

// Sink defeats dead-code elimination; both candidates accumulate into it.
var Sink int64

// SentinelA is the pre-allocated exception caught at the outermost layer
// (deep propagation). SentinelC is the exception caught one level up (shallow
// throw-catch).
//
// They are pre-allocated and reused so the measured cost is the propagation
// mechanism, not per-throw allocation. A real emitter allocates a fresh
// exception per throw for the Java stack trace; that allocation is identical
// in both candidates and cancels out of a P-vs-F comparison. -benchmem still
// reports whatever the mechanism itself allocates (e.g. panic bookkeeping).
var (
	SentinelA = &JException{Tag: 1}
	SentinelC = &JException{Tag: 3}
)

// Work is the identical per-iteration payload for both candidates. It is a
// small multiplicative mix that depends on the input, so the compiler cannot
// constant-fold it across the chain; it is noinline so it survives as a real
// leaf call in every variant.
//
//go:noinline
func Work(x int) int {
	return (x*2654435761 + 17) >> 5
}

// Guard is a read-only global that keeps Cleanup non-eliminable while keeping
// the cleanup itself nearly free, so the finally-only benchmark measures the
// defer machinery rather than the cleanup work.
var Guard int64

// Cleanup stands in for a finally body or monitor exit: present but nearly
// free. The read of Guard (exported, mutable in principle) prevents the
// compiler from deleting the deferred call.
//
//go:noinline
func Cleanup() {
	_ = Guard
}
