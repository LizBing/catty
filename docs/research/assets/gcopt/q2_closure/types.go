// Package q2_closure measures what the emitter pays if it models bytecode
// locals with Go closures or method values: does the captured state stay on
// the stack (zero escape) or heap-allocate, and how expensive is the call?
package q2_closure

import "gcopt/common"

// C is a concrete type used to build method values.
type C struct {
	n int64
}

// Add is a tiny inlinable method.
func (c *C) Add(a, b int64) int64 {
	return a + b + c.n
}

// Closures that capture locals, called synchronously in the same function.
// These should be stack-allocated and inlined (see -m evidence).

func ProbeClosureByValue() int64 {
	x := common.Seed
	f := func(y int64) int64 { return x + y }
	return f(1) + f(2)
}

func ProbeClosureByRef() int64 {
	x := common.Seed
	f := func() int64 { x++; return x }
	return f() + f()
}

// ProbeClosureEscape returns a closure, forcing it (and any by-reference
// captured locals) to the heap.

func ProbeClosureEscapeByValue() func(int64) int64 {
	x := common.Seed
	return func(y int64) int64 { return x + y }
}

func ProbeClosureEscapeByRef() func() int64 {
	x := common.Seed
	return func() int64 { x++; return x }
}

// ProbeMethodValueConcrete binds a concrete receiver.
func ProbeMethodValueConcrete() func(int64, int64) int64 {
	c := &C{n: 1}
	return c.Add
}

// Adder is an interface so we can also measure interface method values.
type Adder interface {
	Add(a, b int64) int64
}

// makeAdder hides the concrete type behind a noinline boundary.
//
//go:noinline
func makeAdder() Adder {
	return &C{n: 1}
}

// ProbeMethodValueInterface binds an interface method value (captures the
// itab, so a later call stays virtual).
func ProbeMethodValueInterface() func(int64, int64) int64 {
	a := makeAdder()
	return a.Add
}

// callFunc is a noinline consumer that forces any closure passed to it onto
// the heap, modelling "store the generated callable in a slot then call later".
//
//go:noinline
func callFunc(f func() int64) int64 {
	return f()
}

// makeAdderClosure returns a closure that captures x by value. The return
// forces the closure (and its embedded copy of x) to the heap.
//
//go:noinline
func makeAdderClosure(x int64) func(int64) int64 {
	return func(y int64) int64 { return x + y }
}

// makeCounter returns a closure that captures x by reference (it mutates x),
// so returning it heap-allocates BOTH the closure and the boxed x.
//
//go:noinline
func makeCounter(x int64) func() int64 {
	return func() int64 { x++; return x }
}
