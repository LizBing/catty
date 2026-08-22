// Package hybrid implements object-model candidate A from ADR-0004 (the
// "mixed representation"): every Java class becomes a Go struct whose
// embedded fields mirror the inheritance chain, and a reference whose static
// type is too wide is promoted to a Go interface so dispatch goes through the
// itab. This is the representation the emitter would use under the current
// provisional decision.
package hybrid

import "objbench/common"

// sink accumulates benchmark results so the compiler cannot eliminate work.
var sink int64

// Animal is the Java-style supertype used for virtual (itab) dispatch.
type Animal interface {
	Compute() int64
}

// Depth-3 inheritance chain: D3Base -> D3Mid -> D3Leaf.
type D3Base struct {
	f0 int64
}

type D3Mid struct {
	D3Base
	f1 int64
}

type D3Leaf struct {
	D3Mid
	f2 int64
}

// Compute reads one field per level of the chain.
func (l *D3Leaf) Compute() int64 { return l.f0 + l.f1 + l.f2 }

// Depth-6 inheritance chain: D6A -> D6B -> D6C -> D6D -> D6E -> D6Leaf.
type D6A struct {
	f0 int64
}

type D6B struct {
	D6A
	f1 int64
}

type D6C struct {
	D6B
	f2 int64
}

type D6D struct {
	D6C
	f3 int64
}

type D6E struct {
	D6D
	f4 int64
}

type D6Leaf struct {
	D6E
	f5 int64
}

func (l *D6Leaf) Compute() int64 {
	return l.f0 + l.f1 + l.f2 + l.f3 + l.f4 + l.f5
}

// Megamorphic leaf types at depth 3: each embeds D3Mid and adds its own
// field, yielding a distinct concrete type that satisfies Animal.
type Mega3A struct {
	D3Mid
	z int64
}

func (m *Mega3A) Compute() int64 { return m.f0*common.MA + m.f1*(common.MA+1) + m.z*(common.MA+2) }

type Mega3B struct {
	D3Mid
	z int64
}

func (m *Mega3B) Compute() int64 { return m.f0*common.MB + m.f1*(common.MB+1) + m.z*(common.MB+2) }

type Mega3C struct {
	D3Mid
	z int64
}

func (m *Mega3C) Compute() int64 { return m.f0*common.MC + m.f1*(common.MC+1) + m.z*(common.MC+2) }

type Mega3D struct {
	D3Mid
	z int64
}

func (m *Mega3D) Compute() int64 { return m.f0*common.MD + m.f1*(common.MD+1) + m.z*(common.MD+2) }

type Mega3E struct {
	D3Mid
	z int64
}

func (m *Mega3E) Compute() int64 { return m.f0*common.ME + m.f1*(common.ME+1) + m.z*(common.ME+2) }

type Mega3F struct {
	D3Mid
	z int64
}

func (m *Mega3F) Compute() int64 { return m.f0*common.MF + m.f1*(common.MF+1) + m.z*(common.MF+2) }

type Mega3G struct {
	D3Mid
	z int64
}

func (m *Mega3G) Compute() int64 { return m.f0*common.MG + m.f1*(common.MG+1) + m.z*(common.MG+2) }

type Mega3H struct {
	D3Mid
	z int64
}

func (m *Mega3H) Compute() int64 { return m.f0*common.MH + m.f1*(common.MH+1) + m.z*(common.MH+2) }

// Megamorphic leaf types at depth 6: each embeds D6E and adds its own field.
type Mega6A struct {
	D6E
	z int64
}

func (m *Mega6A) Compute() int64 {
	return m.f0*common.MA + m.f1*(common.MA+1) + m.f2*(common.MA+2) + m.f3*(common.MA+3) + m.f4*(common.MA+4) + m.z*(common.MA+5)
}

type Mega6B struct {
	D6E
	z int64
}

func (m *Mega6B) Compute() int64 {
	return m.f0*common.MB + m.f1*(common.MB+1) + m.f2*(common.MB+2) + m.f3*(common.MB+3) + m.f4*(common.MB+4) + m.z*(common.MB+5)
}

type Mega6C struct {
	D6E
	z int64
}

func (m *Mega6C) Compute() int64 {
	return m.f0*common.MC + m.f1*(common.MC+1) + m.f2*(common.MC+2) + m.f3*(common.MC+3) + m.f4*(common.MC+4) + m.z*(common.MC+5)
}

type Mega6D struct {
	D6E
	z int64
}

func (m *Mega6D) Compute() int64 {
	return m.f0*common.MD + m.f1*(common.MD+1) + m.f2*(common.MD+2) + m.f3*(common.MD+3) + m.f4*(common.MD+4) + m.z*(common.MD+5)
}

type Mega6E struct {
	D6E
	z int64
}

func (m *Mega6E) Compute() int64 {
	return m.f0*common.ME + m.f1*(common.ME+1) + m.f2*(common.ME+2) + m.f3*(common.ME+3) + m.f4*(common.ME+4) + m.z*(common.ME+5)
}

type Mega6F struct {
	D6E
	z int64
}

func (m *Mega6F) Compute() int64 {
	return m.f0*common.MF + m.f1*(common.MF+1) + m.f2*(common.MF+2) + m.f3*(common.MF+3) + m.f4*(common.MF+4) + m.z*(common.MF+5)
}

type Mega6G struct {
	D6E
	z int64
}

func (m *Mega6G) Compute() int64 {
	return m.f0*common.MG + m.f1*(common.MG+1) + m.f2*(common.MG+2) + m.f3*(common.MG+3) + m.f4*(common.MG+4) + m.z*(common.MG+5)
}

type Mega6H struct {
	D6E
	z int64
}

func (m *Mega6H) Compute() int64 {
	return m.f0*common.MH + m.f1*(common.MH+1) + m.f2*(common.MH+2) + m.f3*(common.MH+3) + m.f4*(common.MH+4) + m.z*(common.MH+5)
}

// Constructors are marked noinline so the compiler cannot see through them
// and constant-fold an object's field values or its type identity inside a
// benchmark loop. This keeps the measured access cost honest.

//go:noinline
func newD3Leaf() *D3Leaf {
	l := &D3Leaf{}
	l.f0 = common.Seed
	l.f1 = common.Seed + 1
	l.f2 = common.Seed + 2
	return l
}

//go:noinline
func newD6Leaf() *D6Leaf {
	l := &D6Leaf{}
	l.f0 = common.Seed
	l.f1 = common.Seed + 1
	l.f2 = common.Seed + 2
	l.f3 = common.Seed + 3
	l.f4 = common.Seed + 4
	l.f5 = common.Seed + 5
	return l
}

//go:noinline
func newMega3A() *Mega3A {
	m := &Mega3A{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3B() *Mega3B {
	m := &Mega3B{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3C() *Mega3C {
	m := &Mega3C{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3D() *Mega3D {
	m := &Mega3D{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3E() *Mega3E {
	m := &Mega3E{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3F() *Mega3F {
	m := &Mega3F{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3G() *Mega3G {
	m := &Mega3G{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega3H() *Mega3H {
	m := &Mega3H{}
	m.f0, m.f1, m.z = common.Seed, common.Seed+1, common.Seed+2
	return m
}

//go:noinline
func newMega6A() *Mega6A {
	m := &Mega6A{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6B() *Mega6B {
	m := &Mega6B{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6C() *Mega6C {
	m := &Mega6C{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6D() *Mega6D {
	m := &Mega6D{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6E() *Mega6E {
	m := &Mega6E{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6F() *Mega6F {
	m := &Mega6F{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6G() *Mega6G {
	m := &Mega6G{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func newMega6H() *Mega6H {
	m := &Mega6H{}
	m.f0, m.f1, m.f2, m.f3, m.f4, m.z = common.Seed, common.Seed+1, common.Seed+2, common.Seed+3, common.Seed+4, common.Seed+5
	return m
}

//go:noinline
func setD3Leaf(l *D3Leaf, v int64) {
	l.f0 = v
	l.f1 = v
	l.f2 = v
}

//go:noinline
func setD6Leaf(l *D6Leaf, v int64) {
	l.f0 = v
	l.f1 = v
	l.f2 = v
	l.f3 = v
	l.f4 = v
	l.f5 = v
}
