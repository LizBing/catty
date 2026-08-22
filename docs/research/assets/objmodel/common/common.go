// Package common holds scaffolding shared by the two object-model benchmark
// candidates (hybrid and uniform) so their workloads stay bit-identical where
// it matters for a fair comparison.
package common

// Seed is an exported, non-constant value used to populate object fields.
// Because it is exported, the compiler cannot assume it is never mutated and
// therefore cannot constant-fold benchmark workloads that read it.
var Seed int64 = 0x5f3759df

// Species multipliers shared by both candidates so that the arithmetic each
// benchmark performs is identical regardless of representation. Eight species
// give the megamorphic dispatch a realistic number of concrete targets.
const (
	MA = 3
	MB = 5
	MC = 7
	MD = 11
	ME = 13
	MF = 17
	MG = 19
	MH = 23
)
