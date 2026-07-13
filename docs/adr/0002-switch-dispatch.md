# ADR-0002: Switch dispatch for the interpreter

- **Status:** Superseded by [ADR-0024](./0024-interpreter-policy-and-evidence-driven-optimization.md)
- **Date:** 2026-07-11

## Context

The interpreter's core cost is **dispatch**: reading the next opcode and jumping
to its handler, once per instruction. The classic taxonomy (Eli Bendersky;
Peter Liniker's benchmarks):

- **switch dispatch** — one central `switch (op) { … }`, one indirect branch per
  instruction, subject to branch-prediction misses on the unpredictable opcode
  stream.
- **direct threading** (computed goto) — each handler ends with `goto *table[next]`,
  giving every opcode its own branch site; the predictor sees each opcode's own
  history → far fewer mispredicts. Typically ~1.5–2× faster than switch.
- **call threading** — each handler is a function that tail-calls the next.

Go has **no computed goto** (the GCC extension). Direct threading is therefore
unavailable. The remaining candidates are switch and call-threading.

## Decision

Use a single dense `switch` on the opcode byte (`interpreter/interpreter.go`),
not a function-pointer dispatch table.

## Consequences

**Positive**
- Idiomatic Go, readable, and trivially debuggable (a single function to step
  through). The 136 handlers stay co-located and grouped by category.
- Go's compiler turns a dense value-range `switch` over opcode bytes into an
  efficient jump table; the one indirect branch per instruction is well-predicted
  enough for MVP throughput.
- A call-threaded design would pay a function-call (with Go's calling convention
  and lack of guaranteed tail-call optimization) per instruction — measured to
  cost *more* per instruction than the branch mispredicts it saves here.

**Negative**
- catty runs `fib(35)` ~7× slower than HotSpot's interpreter, which uses an
  assembly template interpreter with computed-goto-style dispatch. A meaningful
  chunk of that gap is dispatch.
- The dispatch loop is one large function ( readability was traded for keeping
  it inlinable; `Loop`/`exec` may be re-merged if profiling favors it — see
  `ROADMAP.md` Theme B).

## Alternatives considered
- **Assembly computed-goto**: would regain the ~2× but breaks portability and
  the "pure Go" property. Rejected for Phase 1; revisit only if interpreter
  throughput becomes the binding constraint *before* AOT lands.
- **Predecoded instruction + switch**: decode bytecode into a flat
  `[]Instruction` once, removing operand parsing from the hot loop. Orthogonal
  to this ADR and the highest-payoff interpreter tuning on the roadmap.
