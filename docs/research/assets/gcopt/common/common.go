// Package common holds scaffolding shared by the gcopt microbenchmarks so the
// workloads stay identical across candidates. The seed is an exported var, so
// the compiler cannot constant-fold benchmark loops that read it.
package common

// Seed is an exported, non-constant value used to populate operands.
var Seed int64 = 0x5f3759df

// Sink is an exported accumulator that benchmarks write into, preventing the
// compiler from dead-code-eliminating the measured work.
var Sink int64
