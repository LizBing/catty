# Roadmap

## Vision: catty as an experimental JRE

catty is not just "a JVM written in Go" — it explores a platform that compiles
Java programs into Go programs while reusing appropriate Go runtime services.
The exact Thread, memory-model, class-library, I/O, and AOT production boundaries
remain subject to Accepted ADRs and measured workstreams.

Withdrawn ADRs 0008–0013 preserve early hypotheses but do not define the
current architecture.

## Completed

### Phase 1 — Interpreter MVP ✅
Switch-dispatched bytecode interpreter, ~140 opcodes, 5 native core classes,
8 fixtures byte-identical to `java`.

### A0–A4 — AOT transpiler ✅
Bytecode IR + stack elimination → type tracking → fresh-per-def type-aware
emitter → invoke bridge → loops → diamonds (phi) → OOP → long/float/double →
edge items → `catty build` (whole-program offline AOT). fib(35) at native
speed (~44 ms, on par with HotSpot JIT).

## Strategic decisions and open questions

| ADR | Status | Decision / question | Impact |
|---|---|---|---|
| 0016 | Accepted | AOT is the primary product path; the interpreter remains a permanent semantic fallback | Multi-engine execution boundary |
| 0017–0024 | Accepted | Java 25 semantics, Go-runtime boundary, dissolution, representation, bootstrap, String, and interpreter policy | R2 governing constraints |
| 0008 | Withdrawn | Earlier AOT-first proposal | Replaced by ADR-0016's multi-engine model |
| 0009 | Withdrawn | Evaluate a hybrid class library | Bootstrap control and compatibility |
| 0010 | Withdrawn | Evaluate Java Thread mapping onto Go runtime mechanisms | Thread identity and lifecycle |
| 0011 | Withdrawn | Determine the required Java memory semantics | Compatibility and optimization boundary |
| 0012 | Withdrawn | Evaluate Go escape analysis for Java objects | Allocation optimization |
| 0013 | Withdrawn | Evaluate direct Go runtime integration | Runtime service boundary |

Withdrawn proposals are retained as historical context only and do not
authorize implementation. Their unresolved questions require research and new
Proposed ADRs.

## Implementation phases

### Phase R1 — Run real Java programs (exceptions + opcodes + bootstrap)
**Status:** ✅ Complete

- **Exceptions** (`try/catch/athrow`): full mechanism — athrow, runtime
  exceptions (NPE/Arithmetic/CCE/AIOOBE), try/catch/finally, frame unwinding.
  Native Throwable/Exception hierarchy (~13 classes).
- **Remaining opcodes**: `invokeinterface`, `multianewarray`, `wide` (~145/201
  opcodes). `invokedynamic` deferred to R3.
- **Bootstrap classpath**: `catty` auto-detects `$CATTY_BOOT` / `$JAVA_HOME` /
  `java_home`, prepends java.base to the user classpath. Uses the JDK's own
  `jimage extract` tool (no runtime jimage parser — keeps catty lean). The R1
  implementation currently serves six bootstrap classes synthetically; ADR-0022
  makes that a revisable capability boundary rather than a permanent class list.
- **Native layer expansion**: ~18 synthetic classes + ~40 native
  method registrations. `NoSuchMethodError` (not a crash) on gaps.

**Milestone** ✅: `catty -cp . HelloWorld` with real java.base — one command,
auto-detected. `RealBaseSmoke` (18 assertions) byte-identical to `java`.

### Phase R2 — Runtime semantics and concurrency planning
**Status:** ✅ Complete (bounded concurrency surface)

Initialization and bounded UTF-16 String slices complete. Bounded Java 25
concurrency (Thread/monitor foundation Slices A–E) complete: race-free SC heap
(ADR-0030), concurrency-safe class loading with canonical Class mirrors, stable
Java Thread identity/lifecycle with one goroutine carrier per started platform
Thread, VM daemon liveness, interruptible wait/join/sleep, Java object monitors
and wait sets with notify/interrupt ordering (ADR-0029), `holdsLock`/`wait`
argument validation, and per-Class `initMu`/`initCond` with JVMS §5.5
cross-context initialization protocol. The 19-fixture concurrency matrix
matches Temurin 25 in Interpreter and IR (1× and race-built 100× stress); AOT
reports all concurrency fixtures as `Not implemented`/build rejection. The
multi-threaded producer-consumer milestone is achieved. Timed `wait`/`join`,
`Unsafe`, virtual threads, `ThreadGroup`/`ThreadLocal`,
`java.util.concurrent`, and AOT concurrency remain out of scope and are
governed by later phases.

The completed
[`r2-runtime-semantics-research`](./workstreams/r2-runtime-semantics-research.md)
produced the accepted initialization and UTF-16 String decisions and their two
completed implementation slices. The completed
[`r2-concurrency-semantics-research`](./workstreams/r2-concurrency-semantics-research.md)
established Java-visible Thread, monitor, wait-set, interruption, cross-thread
initialization, and minimum memory-ordering decisions. The completed
[`r2-thread-monitor-foundation-slice`](./workstreams/r2-thread-monitor-foundation-slice.md)
delivered all five implementation slices (A–E). Unsafe-backed class-library
compatibility remains a separate future workstream.

**Milestone** ✅: multi-threaded producer-consumer program.

### Phase R3 — Reflection & dynamic features
**Status:** After R2

- `java.lang.reflect` on retained runtime metadata, subject to a dedicated
  reflection design under ADR-0016's multi-engine boundary.
- `invokedynamic` full support (LambdaMetafactory, dynamic proxies).
- Annotation parsing (`RuntimeVisibleAnnotations`). ~3–4 weeks.

### Phase R4 — I/O & network (future API-family decisions required)
**Status:** After R3

- Define supported file, socket, selector, and native-integration semantics.
- Evaluate direct Go runtime adapters against compatibility and maintenance cost.

### Phase R5 — AOT coverage expansion (governed by ADR-0016)
**Status:** After R3/R4

- Instance method AOT. Exception handling in emitted code.
- `invokedynamic` AOT (CallSite resolution at build time).
- Expand AOT coverage with explicit interpreter fallback; engine selection and
  any runtime promotion policy remain evidence-driven. ~4–6 weeks.

### Phase R6 — Performance & polish (evidence-driven)
**Status:** After R5

- Escape-analysis and allocation-strategy validation.
- Object layout optimization (field reordering).
- Value-type identification. Core class optimization.
- JCK testing. ~8+ weeks.

## Performance targets

| Milestone | fib(35) target | Notes |
|---|---|---|
| Today (A4) | ~44 ms | achieved; on par with HotSpot JIT |
| R1 complete | ~44 ms | (same; interpreter for new classes) |
| R5 (real-program AOT) | near hand-written Go | whole-program native |
| R6 (optimized) | ≤ hand-written Go | escape analysis + layout opt |

## What can't catty run yet?

Minecraft? No. Real Java applications need: exceptions (R1), concurrency (R2),
reflection (R3), I/O (R4), and much broader class-library compatibility. Each
is a documented phase. The AOT path proves the performance thesis; runtime
semantics and class-library scope remain the primary research work.
