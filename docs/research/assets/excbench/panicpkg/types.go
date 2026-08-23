// Package panicpkg implements candidate P: Java throw = Go panic of a
// *common.JException sentinel, caught at method boundaries by defer/recover.
//
// Every method that owns a handler entry installs a deferred recover. The
// deferred function type-asserts the recovered value to the sentinel, compares
// its class tag against the tags the method handles, runs the handler on a
// match, and re-panics on a mismatch (Java exception-table "not handled here,
// propagate"). A non-sentinel recovered value (a genuine Go panic such as a
// nil dereference) is also re-panicked, which is what keeps the two kinds of
// panic strictly distinct.
package panicpkg

import "excbench/common"

// Handler tags recognized at each layer. SentinelA.Tag == tagA and
// SentinelC.Tag == tagC, so deep propagation is caught at layer a and shallow
// throw-catch at layer c; the other layers re-panic.
const (
	tagA = uint32(1)
	tagB = uint32(2)
	tagC = uint32(3)
	tagD = uint32(4)
)

// Per-layer handler results, chosen distinct so a correctness test can detect
// an exception swallowed by the wrong layer (each layer only matches its own
// tag; everything else must re-panic).
const (
	resultA = -200
	resultB = -201
	resultC = -100
	resultD = -101
)

// --- Normal path, every layer owns a (non-matching) handler ----------------

//go:noinline
func aNorm(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagA {
				r = resultA
			} else {
				panic(e)
			}
		}
	}()
	return bNorm(x)
}

//go:noinline
func bNorm(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagB {
				r = resultB
			} else {
				panic(e)
			}
		}
	}()
	return cNorm(x)
}

//go:noinline
func cNorm(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagC {
				r = resultC
			} else {
				panic(e)
			}
		}
	}()
	return dNorm(x)
}

//go:noinline
func dNorm(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagD {
				r = resultD
			} else {
				panic(e)
			}
		}
	}()
	return common.Work(x)
}

// --- Normal path, no layer owns a handler (the realistic majority) ---------

// Methods with an empty exception table do not intercept exceptions at all; a
// panic simply unwinds their frames for free. This is the "zero normal-path
// cost" end of candidate P.

//go:noinline
func aNoHandlers(x int) int {
	return bNoHandlers(x)
}

//go:noinline
func bNoHandlers(x int) int {
	return cNoHandlers(x)
}

//go:noinline
func cNoHandlers(x int) int {
	return dNoHandlers(x)
}

//go:noinline
func dNoHandlers(x int) int {
	return common.Work(x)
}

// --- Normal path, every layer owns a plain (non-recovering) defer ----------
//
// A finally block or a synchronized block maps to a plain defer in P: it runs
// on both the normal and the unwinding path, but it does not recover, so it
// does not pay the recover + type-assert + tag-match + re-panic machinery. It
// sits between the no-handler and the catch-handler extremes.

//go:noinline
func aFinally(x int) int {
	defer common.Cleanup()
	return bFinally(x)
}

//go:noinline
func bFinally(x int) int {
	defer common.Cleanup()
	return cFinally(x)
}

//go:noinline
func cFinally(x int) int {
	defer common.Cleanup()
	return dFinally(x)
}

//go:noinline
func dFinally(x int) int {
	defer common.Cleanup()
	return common.Work(x)
}

// --- Shallow throw-catch: d throws, c catches (one boundary) ---------------

//go:noinline
func dShallow(x int) int {
	panic(common.SentinelC)
}

//go:noinline
func cShallow(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagC {
				r = resultC
			} else {
				panic(e)
			}
		}
	}()
	return dShallow(x)
}

// --- Deep propagation: d throws, caught at a -------------------------------

// Variant 1 (realistic): intermediate layers own no handler, so the panic
// unwinds their frames for free.

//go:noinline
func aDeepFree(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagA {
				r = resultA
			} else {
				panic(e)
			}
		}
	}()
	return bDeepFree(x)
}

//go:noinline
func bDeepFree(x int) int {
	return cDeepFree(x)
}

//go:noinline
func cDeepFree(x int) int {
	return dDeepFree(x)
}

//go:noinline
func dDeepFree(x int) int {
	panic(common.SentinelA)
}

// Variant 2 (worst case, per the task's spec): every intermediate layer owns
// a non-matching handler, so each recovers and re-panics before the exception
// reaches the matching handler at a.

//go:noinline
func aDeepRethrow(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagA {
				r = resultA
			} else {
				panic(e)
			}
		}
	}()
	return bDeepRethrow(x)
}

//go:noinline
func bDeepRethrow(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagB {
				r = resultB
			} else {
				panic(e)
			}
		}
	}()
	return cDeepRethrow(x)
}

//go:noinline
func cDeepRethrow(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagC {
				r = resultC
			} else {
				panic(e)
			}
		}
	}()
	return dDeepRethrow(x)
}

//go:noinline
func dDeepRethrow(x int) (r int) {
	defer func() {
		if e := recover(); e != nil {
			if je, ok := e.(*common.JException); ok && je.Tag == tagD {
				r = resultD
			} else {
				panic(e)
			}
		}
	}()
	panic(common.SentinelA)
}
