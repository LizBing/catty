# Roadmap

catty is phased: Phase 1 (a correct interpreter) is done; the items below are
the path to both broader spec coverage and the performance ceiling the
premise ("reuse the Go runtime") is meant to reach. Items are grouped by theme,
not by a fixed release order — pick what matters.

## Phase 1 — interpreter MVP ✅

Delivered. A switch-dispatched bytecode interpreter runs a minimal-but-real
Java subset, byte-identical to `java` on the test corpus. See
[`CHANGELOG.md`](./CHANGELOG.md) for the contents and
[`ARCHITECTURE.md`](./ARCHITECTURE.md) §7 for the baseline.

## Theme A — the AOT transpiler (the performance arc)

This is the headline Phase 2 work and the only route to native-class speed. The
idea: statically lower JVM bytecode to **Go source**, then compile it with the
Go toolchain. The Go compiler becomes catty's optimizing "JIT backend", and the
Go runtime remains the GC/scheduler — exactly the premise, extended.

Sketch of the design:

1. **Stack elimination**. Bytecode is stack-based, but every instruction's
   stack effect is statically known. A depth dataflow computes, at each program
   point, the stack depth, turning the operand stack into slot-indexed virtual
   registers. ✅ **Done in A0** (depth-only; type tracking / true SSA is the
   emitter's concern).
2. **Lowering to Go**. Each Java method becomes a Go `func`. Java objects → Go
   structs (already the case); Java calls → direct Go calls after resolution;
   Java control flow → Go control flow. `long`/`double` map to `int64`/`float64`.
3. **Code generation**. Emit one Go package per class to a generated tree, then
   `go build` it (whole-program ahead-of-time) or `plugin` (lazy/hot). The
   interpreter remains as the fallback for code that can't be lowered (e.g.
   reflective or not-yet-supported cases).

Why this works *because* of the Go choice: a C/Rust host language would need a
hand-written SSA backend + register allocator + a GC. catty gets all three for
free from Go's toolchain and runtime.

Milestones:
- [x] **A0** — stack-elimination pass + an IR executor that runs it, verified
  end-to-end (3 engines byte-identical on all fixtures). *De-risks the bet: the
  stack can be statically eliminated.* See `lowering/` and the A0 changelog entry.
- [x] **A1** — emit + compile a trivial method (`fib`) to Go source, run it
  natively. Reuses A0's `lowering.IR`; the new work is the Go-source emitter
  (`transpile.Emit`). *Emitted `fib(35)` runs in ~44 ms — native Go speed, ~100×
  the interpreter and on par with HotSpot JIT.* See `transpile/` and the A1
  changelog entry.
- [ ] **A2** — object model, field access, virtual dispatch in emitted code.
- [ ] **A3** — whole-program transpile of the test corpus, diff vs interpreter.
- [ ] **A4** — hot-method selection: interpret cold code, transpile hot.

Risks to call out up front: `invokedynamic`/dynamic call sites, exceptions
(which need Go `panic`/`recover` bridging or result codes), and the JMM (the
concurrency arc must land first or in lockstep).

## Theme B — interpreter tuning (the ~7× gap to `java -Xint`)

Before AOT lands, the interpreter itself has clear, bounded headroom. Ordered
by expected payoff:

- [~] **Predecode bytecode** into a flat instruction array once per method, so
  the hot loop does no operand parsing. **Tried in A0 (`LoopIR`): it measured
  ~6% *slower* than the tree-walker** — the per-instruction overhead (stack-
  pointer seeding, instruction-struct dereference, slot-accessor calls) outweighs
  the saved operand parsing. Predecode inside a Go interpreter does not pay off;
  the IR's value is validation + as the input to the AOT emitter, not speed.
  Recorded in ADR-0006. The remaining Theme-B items are still open.
- [ ] **`sync.Pool` for `Frame`** — method entry currently allocates a frame
  (locals + operand stack). Pooling eliminates per-call allocation in tight
  recursion.
- [ ] **Split `Slot`** into parallel `[]int32` / `[]*Object` arrays — halves
  memory traffic on numeric-heavy code (16 B/word → 8 B). Touches every handler,
  so do it once and benchmark.
- [ ] **Inline the dispatch loop** — collapse `Loop`/`exec` into one function
  (currently split for readability) and ensure hot cases stay in the branch
  predictor's good graces.
- [ ] **`go tool pprof`** a CPU profile of `BenchFib` to confirm the real
  hotspots before optimizing by guess.

## Theme C — spec coverage

Each closes a class of real Java programs. Roughly in value order:

- [ ] **Exceptions** (`try`/`catch`/`finally`, `athrow`). The exception tables
  are already parsed (`rtda.Method.ExceptionTable()`); wire them into the
  dispatch loop's error path. Bridge to Go errors/panics so native methods can
  throw.
- [ ] **`invokedynamic` / lambdas / string concat (Java 9+)**. Needed to run
  programs compiled without `-source 8`. Requires modeling `CallSite` /
  `bootstrap methods`, or at minimum `StringConcatFactory`.
- [ ] **Interfaces fully** — `invokeinterface` currently panics; resolution
  walks the itable, which `rtda.Class` already supports partially.
- [ ] **More core classes** — `Math`, `Integer`/`Long` parse/format, the
  primitive wrappers, `Object` helpers (`equals`/`hashCode`/`toString`).
- [ ] **`java.lang.String` properly** — backed by a `char`/`byte` array rather
  than a Go string in `extra`, so methods like `charAt`/`substring` work.

## Theme D — concurrency (the "Go runtime" payoff, finally visible)

The premise promises cheap threads; deliver them:

- [ ] **`java.lang.Thread` → goroutine**. The runtime becomes multi-`Thread`.
- [ ] **Per-object monitors** for `synchronized` (lazily-allocated `sync.Mutex`
  per object, or an embedded lock word), and `wait`/`notify` via `sync.Cond`.
- [ ] **JMM approximation** — `volatile` and `happens-before`. Java's and Go's
  memory models differ; document the deviations rather than silently mismatch.

## Theme E — project hygiene

- [ ] **CI** — a GitHub Actions workflow running `go vet`, `go test`, and
  `./tests/run.sh` on push.
- [ ] **More unit tests** for `classfile` (every constant-pool tag) and
  `rtda` (object/field layout, array indexing).
- [ ] **Fuzzing** the class-file parser (`go test -fuzz`) — it's the trust
  boundary for untrusted `.class` input.
- [ ] **`golangci-lint`** with a checked-in config once the codebase stabilizes.

## Performance targets (illustrative)

| Milestone | Target on `fib(35)` |
|---|---|
| Today (Phase 1) | 4.3 s (~7× slower than `java -Xint`) |
| Theme B complete | within ~2–3× of `java -Xint` |
| Theme A1 (single method AOT) | `fib` near hand-written Go speed |
| Theme A3 (whole-program AOT) | within a small constant of `java` JIT |

These are design targets, not commitments — measure against `BenchFib` as each
theme lands.
