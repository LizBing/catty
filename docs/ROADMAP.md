# Roadmap

## Vision: catty as an experimental JRE

catty is not just "a JVM written in Go" — it is **a platform that compiles Java
programs into Go programs**, where the final product is a native Go binary running
on Go's GC, scheduler, and network stack. Java's Thread/synchronized/volatile/
GC/IO are "dissolved" into Go's goroutine/mutex/atomic/GC/netpoll at compile time.

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
| 0008 | AOT-first (interpreter is dev tier) | No JIT warmup, no safepoints |
| 0009 | Hybrid class library (~50 native + ~7000 interpreted) | Bootstrap control + semantic compat |
| 0010 | Thread = goroutine | Virtual threads from day one |
| 0011 | Go memory model (not JMM) | Superseded by accepted ADR-0016 |
| 0012 | Escape analysis replaces generational GC | Stack-allocate Java objects |
| 0013 | Direct Go runtime integration | I/O = native Go performance |

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
- **Native layer expansion** (ADR-0009): ~18 synthetic classes + ~40 native
  method registrations. `NoSuchMethodError` (not a crash) on gaps.

**Milestone** ✅: `catty -cp . HelloWorld` with real java.base — one command,
auto-detected. `RealBaseSmoke` (18 assertions) byte-identical to `java`.

### Phase R2 — Concurrency (ADR-0010)
**Status:** Next

Prerequisite sizing is caller-driven, not a native-method count. On Temurin
25.0.3, Integer/Long decimal conversion reaches a narrow DecimalDigits Unsafe
array-write path; Double parsing and ordinary HashMap operations currently fail
for separate initialization/class-library reasons. R2 begins with strict native
failure and the checked caller graph, then implements logical Unsafe profiles
alongside Thread and memory semantics (ADRs 0016–0019).

- `java.lang.Thread` → goroutine. Thread.start = `go func()`.
- Lazy reentrant Monitor per object; explicit owner/depth/wait-set state wraps
  Go synchronization primitives (ADR-0017).
- `Thread.sleep` uses Go timers; interrupt is explicit Java status plus waiter
  wakeup, not merely `context.Cancel`.
- Protected Java memory semantics plus measured Strict/Go-native/Hybrid study
  (ADR-0016). Schedule is re-estimated after R2-GATE closes.

**Milestone**: multi-threaded producer-consumer program.

### Phase R3 — Reflection & dynamic features
**Status:** After R2

- `java.lang.reflect` on retained `rtda` metadata (ADR-0007 tiered model).
- `invokedynamic` full support (LambdaMetafactory, dynamic proxies).
- Annotation parsing (`RuntimeVisibleAnnotations`). ~3–4 weeks.

### Phase R4 — I/O & network (ADR-0013)
**Status:** After R3

- java.io → Go os.File. java.net → Go net.Conn/Listener.
- NIO Selector → Go netpoll. JNI bridge (long-term, for LWJGL).
- ~4–8 weeks (scope-dependent).

### Phase R5 — AOT coverage expansion (ADR-0008)
**Status:** After R3/R4

- Instance method AOT. Exception handling in emitted code.
- `invokedynamic` AOT (CallSite resolution at build time).
- Tiered: interpret cold, AOT hot. ~4–6 weeks.

### Phase R6 — Performance & polish (ADR-0012)
**Status:** After R5

- Escape analysis validation (stack-allocated Java objects).
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
