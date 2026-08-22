// Package uniform implements object-model candidate B from ADR-0004 (the
// "unified *Object" representation): every reference is a header pointer
// whose fields live in a slot array, and every access goes through a type
// assertion against the header's class tag.
package uniform

import "objbench/common"

// sink accumulates benchmark results so the compiler cannot eliminate work.
var sink int64

// Class identifies a concrete Java class; it plays the role of the type tag
// in the unified *Object representation.
type Class struct {
	name string
}

// Object is the unified representation: a header pointer plus a fixed slot
// array. A fixed array (rather than a slice) keeps both candidates at one
// allocation per object so the comparison isolates access cost rather than
// allocation count.
type Object struct {
	cls *Class
	f   [16]int64
}

var (
	classD3Leaf = &Class{"D3Leaf"}
	classD6Leaf = &Class{"D6Leaf"}

	classMega3A = &Class{"Mega3A"}
	classMega3B = &Class{"Mega3B"}
	classMega3C = &Class{"Mega3C"}
	classMega3D = &Class{"Mega3D"}
	classMega3E = &Class{"Mega3E"}
	classMega3F = &Class{"Mega3F"}
	classMega3G = &Class{"Mega3G"}
	classMega3H = &Class{"Mega3H"}

	classMega6A = &Class{"Mega6A"}
	classMega6B = &Class{"Mega6B"}
	classMega6C = &Class{"Mega6C"}
	classMega6D = &Class{"Mega6D"}
	classMega6E = &Class{"Mega6E"}
	classMega6F = &Class{"Mega6F"}
	classMega6G = &Class{"Mega6G"}
	classMega6H = &Class{"Mega6H"}
)

// d3LeafCall is the monomorphic "assert then call": check the class tag, then
// read the leaf's slots.
func d3LeafCall(o *Object) int64 {
	if o.cls != classD3Leaf {
		panic("unexpected class")
	}
	return o.f[0] + o.f[1] + o.f[2]
}

// d6LeafCall is the depth-6 counterpart of d3LeafCall.
func d6LeafCall(o *Object) int64 {
	if o.cls != classD6Leaf {
		panic("unexpected class")
	}
	return o.f[0] + o.f[1] + o.f[2] + o.f[3] + o.f[4] + o.f[5]
}

// bi3Dispatch is bimorphic dispatch via a two-way type switch.
func bi3Dispatch(o *Object) int64 {
	switch o.cls {
	case classMega3A:
		return o.f[0]*common.MA + o.f[1]*(common.MA+1) + o.f[2]*(common.MA+2)
	case classMega3B:
		return o.f[0]*common.MB + o.f[1]*(common.MB+1) + o.f[2]*(common.MB+2)
	}
	return 0
}

// bi6Dispatch is the depth-6 counterpart of bi3Dispatch.
func bi6Dispatch(o *Object) int64 {
	switch o.cls {
	case classMega6A:
		return o.f[0]*common.MA + o.f[1]*(common.MA+1) + o.f[2]*(common.MA+2) + o.f[3]*(common.MA+3) + o.f[4]*(common.MA+4) + o.f[5]*(common.MA+5)
	case classMega6B:
		return o.f[0]*common.MB + o.f[1]*(common.MB+1) + o.f[2]*(common.MB+2) + o.f[3]*(common.MB+3) + o.f[4]*(common.MB+4) + o.f[5]*(common.MB+5)
	}
	return 0
}

// mega3Dispatch is megamorphic dispatch via an eight-way type switch.
func mega3Dispatch(o *Object) int64 {
	switch o.cls {
	case classMega3A:
		return o.f[0]*common.MA + o.f[1]*(common.MA+1) + o.f[2]*(common.MA+2)
	case classMega3B:
		return o.f[0]*common.MB + o.f[1]*(common.MB+1) + o.f[2]*(common.MB+2)
	case classMega3C:
		return o.f[0]*common.MC + o.f[1]*(common.MC+1) + o.f[2]*(common.MC+2)
	case classMega3D:
		return o.f[0]*common.MD + o.f[1]*(common.MD+1) + o.f[2]*(common.MD+2)
	case classMega3E:
		return o.f[0]*common.ME + o.f[1]*(common.ME+1) + o.f[2]*(common.ME+2)
	case classMega3F:
		return o.f[0]*common.MF + o.f[1]*(common.MF+1) + o.f[2]*(common.MF+2)
	case classMega3G:
		return o.f[0]*common.MG + o.f[1]*(common.MG+1) + o.f[2]*(common.MG+2)
	case classMega3H:
		return o.f[0]*common.MH + o.f[1]*(common.MH+1) + o.f[2]*(common.MH+2)
	}
	return 0
}

// mega6Dispatch is the depth-6 counterpart of mega3Dispatch.
func mega6Dispatch(o *Object) int64 {
	switch o.cls {
	case classMega6A:
		return o.f[0]*common.MA + o.f[1]*(common.MA+1) + o.f[2]*(common.MA+2) + o.f[3]*(common.MA+3) + o.f[4]*(common.MA+4) + o.f[5]*(common.MA+5)
	case classMega6B:
		return o.f[0]*common.MB + o.f[1]*(common.MB+1) + o.f[2]*(common.MB+2) + o.f[3]*(common.MB+3) + o.f[4]*(common.MB+4) + o.f[5]*(common.MB+5)
	case classMega6C:
		return o.f[0]*common.MC + o.f[1]*(common.MC+1) + o.f[2]*(common.MC+2) + o.f[3]*(common.MC+3) + o.f[4]*(common.MC+4) + o.f[5]*(common.MC+5)
	case classMega6D:
		return o.f[0]*common.MD + o.f[1]*(common.MD+1) + o.f[2]*(common.MD+2) + o.f[3]*(common.MD+3) + o.f[4]*(common.MD+4) + o.f[5]*(common.MD+5)
	case classMega6E:
		return o.f[0]*common.ME + o.f[1]*(common.ME+1) + o.f[2]*(common.ME+2) + o.f[3]*(common.ME+3) + o.f[4]*(common.ME+4) + o.f[5]*(common.ME+5)
	case classMega6F:
		return o.f[0]*common.MF + o.f[1]*(common.MF+1) + o.f[2]*(common.MF+2) + o.f[3]*(common.MF+3) + o.f[4]*(common.MF+4) + o.f[5]*(common.MF+5)
	case classMega6G:
		return o.f[0]*common.MG + o.f[1]*(common.MG+1) + o.f[2]*(common.MG+2) + o.f[3]*(common.MG+3) + o.f[4]*(common.MG+4) + o.f[5]*(common.MG+5)
	case classMega6H:
		return o.f[0]*common.MH + o.f[1]*(common.MH+1) + o.f[2]*(common.MH+2) + o.f[3]*(common.MH+3) + o.f[4]*(common.MH+4) + o.f[5]*(common.MH+5)
	}
	return 0
}

// Constructors are noinline so the compiler cannot see through them and
// constant-fold an object's class tag or field values inside a benchmark
// loop.

//go:noinline
func newObject3(cls *Class) *Object {
	o := &Object{cls: cls}
	o.f[0] = common.Seed
	o.f[1] = common.Seed + 1
	o.f[2] = common.Seed + 2
	return o
}

//go:noinline
func newObject6(cls *Class) *Object {
	o := &Object{cls: cls}
	o.f[0] = common.Seed
	o.f[1] = common.Seed + 1
	o.f[2] = common.Seed + 2
	o.f[3] = common.Seed + 3
	o.f[4] = common.Seed + 4
	o.f[5] = common.Seed + 5
	return o
}

//go:noinline
func setObject3(o *Object, v int64) {
	o.f[0] = v
	o.f[1] = v
	o.f[2] = v
}

//go:noinline
func setObject6(o *Object, v int64) {
	o.f[0] = v
	o.f[1] = v
	o.f[2] = v
	o.f[3] = v
	o.f[4] = v
	o.f[5] = v
}
