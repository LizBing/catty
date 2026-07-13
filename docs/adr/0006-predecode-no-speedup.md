# ADR-0006: Predecode in the interpreter does not pay off — AOT is the perf path

- **Status:** Superseded by [ADR-0024](./0024-interpreter-policy-and-evidence-driven-optimization.md)
- **Date:** 2026-07-11

## Context

A common interpreter optimization is **predecode**: decode each instruction's
operands once into a flat instruction array, so the hot dispatch loop does no
operand parsing. ROADMAP Theme B listed it as "often the single biggest
interpreter win."

A0 built exactly this: `lowering.Lower` predecodes a method into `[]IRInst`
(resolved operands, branch targets, computed Uses/Defs), and `interpreter.LoopIR`
dispatches over that array instead of raw bytecode. The expectation (stated in
the A0 plan) was that `-ir` would be measurably faster than the tree-walker.

## Decision

Keep the IR executor as a **validation tool and the AOT emitter's input**, but
do *not* pursue it (or further interpreter-internal predecode) for speed. Route
all performance work through the AOT transpiler (emit Go source → `go build`).

## Consequences

**Measured** on `BenchFib` (fib(35), ~29M recursive calls):

| Engine | Time |
|---|---|
| tree-walker (`Loop`) | ~4.5 s |
| IR executor (`LoopIR`) | ~4.8 s (~6% slower) |

The IR executor is **slower**, not faster. Per-instruction overhead — seeding
the operand-stack pointer from the IR depth, the `IRInst` struct dereference,
and slot-accessor method calls — outweighs the saved operand parsing. Even after
eliminating the per-instruction `map[*Method]*IR` lookup (by reusing the IR
across a frame's instructions), the tree-walker still wins. This is consistent
with Go's lack of computed-goto dispatch (ADR-0002): the dispatch itself, not
operand parsing, dominates.

**Positive**
- A genuine, measured data point that redirects effort: stop tuning the
  interpreter, build the AOT transpiler. The north star ("最优性能") is only
  reachable by compiling bytecode to native code, not by a faster interpreter.
- The IR executor still earns its place — it is the equivalence gate proving the
  lowering is semantics-preserving, and it is the exact input the AOT emitter
  consumes.

**Negative**
- The "Theme B — interpreter tuning" roadmap item is mostly closed: predecode is
  discounted, and the remaining items (frame pooling, slot splitting) are
  marginal next to the JIT gap. The honest path to speed is AOT, not interpreter
  micro-optimization.

## Alternatives considered
- **Optimize the IR executor** (cache IR on the frame, skip stack-pointer
  seeding for pure ops, flatten the dispatch). Tried the IR-reuse cache;
  marginal. The rest is diminishing returns against the ~87× JIT gap and not
  worth the complexity.
- **Hand-written assembly dispatch** (computed goto) would help the interpreter
  but breaks the "pure Go" property and is dwarfed by what AOT achieves.
