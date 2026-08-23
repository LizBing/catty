// Package q1_devirt probes Go's interface devirtualization: when the concrete
// type behind an interface value is statically provable, does gc turn the
// indirect (itab) call into a direct call and inline it? It also locates the
// inliner's cost-budget boundary using generated methods M00..M40.
package q1_devirt

// Adder is the minimal interface under test. Its method body is tiny so that
// any failure to devirtualize/inline shows up as a clear indirect-call cost.
type Adder interface {
	Add(a, b int64) int64
}

// C1 is a concrete implementation with a single field so field access through
// the interface forces an itab+data word load unless devirtualized.
type C1 struct {
	n int64
}

// Add is small enough to inline (cost well under 80).
func (c *C1) Add(a, b int64) int64 {
	return a + b + c.n
}

// C2 is a second concrete type so bimorphic call sites can be tested too.
type C2 struct {
	n int64
}

func (c *C2) Add(a, b int64) int64 {
	return a + b + c.n*2
}

// ProbeDirect is the baseline: a concrete (pointer) receiver call.
func ProbeDirect(c *C1, a, b int64) int64 {
	return c.Add(a, b)
}

// ProbeInterfaceLocal assigns a concrete value to an interface local and calls
// through it in the same function, the strongest devirtualization case.
func ProbeInterfaceLocal(a, b int64) int64 {
	var v Adder = &C1{n: 1}
	return v.Add(a, b)
}

// ProbeAssertThenCall narrows an opaque interface with a type assertion and
// then calls the concrete method.
func ProbeAssertThenCall(v Adder, a, b int64) int64 {
	if c, ok := v.(*C1); ok {
		return c.Add(a, b)
	}
	return 0
}

// ProbeOpaqueParam calls through an interface that crosses a function
// boundary, so the concrete type is genuinely unknown: no devirtualization is
// possible.
func ProbeOpaqueParam(v Adder, a, b int64) int64 {
	return v.Add(a, b)
}

// Holder stores the interface in a struct field; the field write is visible
// only if the constructor is inlined.
type Holder struct {
	v Adder
}

// ProbeField calls through an interface stored in a struct field.
func ProbeField(h *Holder, a, b int64) int64 {
	return h.v.Add(a, b)
}

// makeC1 returns the interface from behind a noinline boundary so call sites
// cannot see the concrete type — used to build a genuinely virtual baseline.
//
//go:noinline
func makeC1() Adder {
	return &C1{n: 1}
}

// makeC2 is the second opaque factory.
//
//go:noinline
func makeC2() Adder {
	return &C2{n: 1}
}
