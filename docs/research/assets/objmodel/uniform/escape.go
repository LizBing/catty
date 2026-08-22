package uniform

import "objbench/common"

// LocalUseD3 creates an object locally, consumes it, and returns only the
// result. Whether the object escapes to the heap is the escape-analysis
// question this function exists to answer under `go build -gcflags=-m`.
func LocalUseD3() int64 {
	o := &Object{cls: classD3Leaf}
	o.f[0] = common.Seed
	o.f[1] = common.Seed + 1
	o.f[2] = common.Seed + 2
	return d3LeafCall(o)
}

// LocalUseD6 is the depth-6 counterpart of LocalUseD3.
func LocalUseD6() int64 {
	o := &Object{cls: classD6Leaf}
	o.f[0] = common.Seed
	o.f[1] = common.Seed + 1
	o.f[2] = common.Seed + 2
	o.f[3] = common.Seed + 3
	o.f[4] = common.Seed + 4
	o.f[5] = common.Seed + 5
	return d6LeafCall(o)
}

// DispatchUseD3 routes an object through the megamorphic dispatch switch.
func DispatchUseD3() int64 {
	return mega3Dispatch(newObject3(classMega3A))
}

// SinkNoinline is a noinline consumer that forces its argument to escape.
//
//go:noinline
func SinkNoinline(o *Object) int64 {
	return o.f[0] + o.f[1] + o.f[2]
}

// EscapeUseD3 passes an object to a noinline consumer, forcing a heap
// allocation for the object.
func EscapeUseD3() int64 {
	return SinkNoinline(newObject3(classD3Leaf))
}
