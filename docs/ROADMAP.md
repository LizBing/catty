# Roadmap

## Vision: catty as an experimental JRE

catty is not just "a JVM written in Go" — it explores a platform that compiles
Java programs into Go programs while reusing appropriate Go runtime services.
The exact Thread, memory-model, class-library, I/O, and AOT production boundaries
remain subject to Accepted ADRs and measured workstreams.

See the strategic plan and ADRs 0008–0013 for the architectural vision.

## Completed

### Phase 1 — Interpreter MVP ✅
Switch-dispatched bytecode interpreter, ~140 opcodes, 5 native core classes,
8 fixtures byte-identical to `java`.

### A0–A4 — AOT transpiler ✅
Bytecode IR + stack elimination → type tracking → fresh-per-def type-aware
emitter → invoke bridge → loops → diamonds (phi) → OOP → long/float/double →
edge items → `catty build` (whole-program offline AOT). fib(35) at native
speed (~44 ms, on par with HotSpot JIT).

## Strategic ADRs (proposed)

| ADR | Decision | Impact |
|---|---|---|
| 0008 | Evaluate an AOT-first production model | Production-tier boundary |
| 0009 | Evaluate a hybrid class library | Bootstrap control and compatibility |
| 0010 | Evaluate Java Thread mapping onto Go runtime mechanisms | Thread identity and lifecycle |
| 0011 | Determine the required Java memory semantics | Compatibility and optimization boundary |
| 0012 | Evaluate Go escape analysis for Java objects | Allocation optimization |
| 0013 | Evaluate direct Go runtime integration | Runtime service boundary |

These ADRs are discussion inputs only. They do not authorize implementation.

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
  `jimage extract` tool (no runtime jimage parser — keeps catty lean). The 6
  bootstrap classes (Object/String/Class/System/Thread/Throwable) stay synthetic;
  everything else loads from real java.base.
- **Native layer expansion**: ~18 synthetic classes + ~40 native
  method registrations. `NoSuchMethodError` (not a crash) on gaps.

**Milestone** ✅: `catty -cp . HelloWorld` with real java.base — one command,
auto-detected. `RealBaseSmoke` (18 assertions) byte-identical to `java`.

### Phase R2 — Runtime semantics and concurrency planning
**Status:** Requires research and Accepted decisions

JDK 25's `Integer.toString`/`Double.parseDouble`/`HashMap` paths reach
`jdk.internal.misc.Unsafe`; concurrency additionally requires explicit Thread,
monitor, class-initialization, volatile/final, interrupt, liveness, and memory
ordering contracts. The first post-R1 work should establish evidence and
Accepted decisions before selecting implementation mechanisms.

Candidate planning outputs include a caller-backed Unsafe profile, protected
Java memory semantics, explicit Thread identity/lifecycle, monitor behavior,
and deterministic differential/race/timeout gates. None is authorized until an
Accepted ADR and workstream require it.

**Milestone**: multi-threaded producer-consumer program.

### Phase R3 — Reflection & dynamic features
**Status:** After R2

- `java.lang.reflect` on retained `rtda` metadata (ADR-0007 tiered model).
- `invokedynamic` full support (LambdaMetafactory, dynamic proxies).
- Annotation parsing (`RuntimeVisibleAnnotations`). ~3–4 weeks.

### Phase R4 — I/O & network (subject to ADR-0013)
**Status:** After R3

- Define supported file, socket, selector, and native-integration semantics.
- Evaluate direct Go runtime adapters against compatibility and maintenance cost.

### Phase R5 — AOT coverage expansion (subject to ADR-0008)
**Status:** After R3/R4

- Instance method AOT. Exception handling in emitted code.
- `invokedynamic` AOT (CallSite resolution at build time).
- Tiered: interpret cold, AOT hot. ~4–6 weeks.

### Phase R6 — Performance & polish (subject to ADR-0012)
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
reflection (R3), I/O (R4), and the full JDK class library (ADR-0009 hybrid).
Each is a documented phase. The AOT path proves the performance thesis; the
class-library gap is the remaining work.
