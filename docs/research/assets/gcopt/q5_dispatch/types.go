// Package q5_dispatch compares three ways to emit a Java "instanceof /
// downcast / virtual-dispatch" ladder: a Go type switch, an if-else type
// assertion chain, and a precomputed int tag switch (the unified
// representation's shape). Eight concrete types make the megamorphic case.
package q5_dispatch

// I is the interface value being dispatched on.
type I interface {
	Val() int64
}

// Concrete types T1..T8. Val returns the type index so the compiler cannot
// fold all branches to one constant.
type T1 struct{}

func (t *T1) Val() int64 { return 1 }

type T2 struct{}

func (t *T2) Val() int64 { return 2 }

type T3 struct{}

func (t *T3) Val() int64 { return 3 }

type T4 struct{}

func (t *T4) Val() int64 { return 4 }

type T5 struct{}

func (t *T5) Val() int64 { return 5 }

type T6 struct{}

func (t *T6) Val() int64 { return 6 }

type T7 struct{}

func (t *T7) Val() int64 { return 7 }

type T8 struct{}

func (t *T8) Val() int64 { return 8 }

// Factories hide the concrete type so the compiler must keep a real itab
// dispatch at every call site.

//go:noinline
func newI1() I { return &T1{} }

//go:noinline
func newI2() I { return &T2{} }

//go:noinline
func newI3() I { return &T3{} }

//go:noinline
func newI4() I { return &T4{} }

//go:noinline
func newI5() I { return &T5{} }

//go:noinline
func newI6() I { return &T6{} }

//go:noinline
func newI7() I { return &T7{} }

//go:noinline
func newI8() I { return &T8{} }

// Obj is the unified-representation shape: one object with a tag field, so
// dispatch is a plain int comparison.
type Obj struct {
	tag int64
	n   int64
}

// tagOf reads the tag through a noinline boundary so the tag switch stays a
// real data-dependent switch.
//
//go:noinline
func tagOf(o *Obj) int64 {
	return o.tag
}
