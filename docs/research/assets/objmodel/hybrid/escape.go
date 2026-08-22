package hybrid

import "objbench/common"

// LocalUseD3 creates a leaf locally, consumes it, and returns only the
// result. Whether the leaf escapes to the heap is the escape-analysis
// question this function exists to answer under `go build -gcflags=-m`.
func LocalUseD3() int64 {
	l := &D3Leaf{}
	l.f0 = common.Seed
	l.f1 = common.Seed + 1
	l.f2 = common.Seed + 2
	return l.Compute()
}

// LocalUseD6 is the depth-6 counterpart of LocalUseD3.
func LocalUseD6() int64 {
	l := &D6Leaf{}
	l.f0 = common.Seed
	l.f1 = common.Seed + 1
	l.f2 = common.Seed + 2
	l.f3 = common.Seed + 3
	l.f4 = common.Seed + 4
	l.f5 = common.Seed + 5
	return l.Compute()
}

// BoxedUseD3 boxes a leaf in an interface and calls through it, which forces
// the concrete value into the interface data word.
func BoxedUseD3() int64 {
	var a Animal = newD3Leaf()
	return a.Compute()
}

// SinkNoinline is a noinline consumer that forces its argument to escape.
//
//go:noinline
func SinkNoinline(l *D3Leaf) int64 {
	return l.f0 + l.f1 + l.f2
}

// EscapeUseD3 passes a leaf to a noinline consumer, forcing a heap
// allocation for the object.
func EscapeUseD3() int64 {
	return SinkNoinline(newD3Leaf())
}
