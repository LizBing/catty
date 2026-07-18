# Roadmap

## Vision: catty as an experimental JVMS runtime platform

catty is not a replacement Java SE JRE. It explores a platform that executes
the JVMS-defined core of Java classfiles and compiles supported programs into
Go programs while reusing appropriate Go runtime services. Its default product
is a small Catty Runtime Profile, with an explicit host-interoperation boundary;
Java SE APIs and an OpenJDK `java.base` image are optional compatibility
profiles, not the definition of the core platform.

The exact JVMS coverage, Thread and memory semantics, Catty runtime APIs, host
services, Java SE compatibility profiles, I/O, and AOT production boundaries
remain subject to Accepted ADRs and measured workstreams. ADR-0034 governs the
profile and host-interoperation boundary; it does not by itself authorize
implementation or change verified capability claims.

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
| 0017 | Superseded | Earlier universal Java 25 supported-capability semantic baseline | Superseded by ADR-0034's profile-scoped contract |
| 0018–0024 | Accepted | Go-runtime boundary, dissolution, representation, bootstrap, String, and interpreter policy | R2 governing constraints |
| 0031–0033 | Accepted | Metadata/typed invocation, InvokeDynamic linkage, and defining-loader/generated-class kernels | R3 shared-kernel constraints; Java SE facades remain optional profile work |
| 0034 | Accepted | Define a JVMS core, parallel Catty/Java SE runtime profiles, and a typed Host ABI | Product and native-boundary direction |
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

### Phase R1 — Classfile execution and optional `java.base` compatibility (exceptions + opcodes + bootstrap)
**Status:** ✅ Complete

- **Exceptions** (`try/catch/athrow`): full mechanism — athrow, runtime
  exceptions (NPE/Arithmetic/CCE/AIOOBE), try/catch/finally, frame unwinding.
  Native Throwable/Exception hierarchy (~13 classes).
- **Remaining opcodes**: `invokeinterface`, `multianewarray`, `wide` (~145/201
  opcodes). `invokedynamic` deferred to R3.
- **Bootstrap classpath**: the current launcher auto-detects `$CATTY_BOOT` /
  `$JAVA_HOME` / `java_home` and prepends `java.base` to the user classpath.
  This is an R1 compatibility mechanism and differential-test reference, not a
  permanent core-runtime requirement. It uses the JDK's own `jimage extract`
  tool (no runtime jimage parser — keeps catty lean). The R1 implementation
  currently serves six bootstrap classes synthetically; ADR-0022 makes that a
  revisable capability boundary rather than a permanent class list.
- **Native layer expansion**: ~18 synthetic classes + ~40 native
  method registrations. `NoSuchMethodError` (not a crash) on gaps.

**Milestone** ✅: `catty -cp . HelloWorld` with an auto-detected real `java.base`
compatibility image. `RealBaseSmoke` (18 assertions) is byte-identical to
`java`; this does not claim complete Java SE or JRE compatibility.

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

- JVMS dynamic-linkage core: retained `BootstrapMethods` metadata,
  MethodHandle/MethodType/CallSite state, and bounded `invokedynamic` linkage.
- Catty Runtime Profile dynamic services: runtime metadata, generated-class
  identity, and host-independent invocation facilities.
- Optional Java SE compatibility facades: `java.lang.reflect`, annotations,
  LambdaMetafactory, and dynamic proxies, each only where an Accepted
  workstream explicitly claims their Java SE behavior.

R3's current research remains a Java 25 compatibility investigation.
ADR-0031 through ADR-0033 have been reconciled with ADR-0034 and Accepted; the
ten Accepted implementation contracts now separate five shared-kernel slices
from five optional Java SE compatibility slices. No production contract is
active until its prerequisites and acceptance anchor are fixed, and no complete
Java reflection or dynamic-feature surface is implied. The previous 3–4 week
estimate is superseded by the profile-separated research estimate; capability
branches are selected independently.

### Phase R4 — Host services, I/O & network (future API-family decisions required)
**Status:** After R3

- Define Catty Host ABI providers for file, socket, selector, graphics, and
  other host services, with typed arguments/results and explicit failures.
- Decide whether individual Java SE APIs are compatibility facades over those
  providers; JNI is not a required platform boundary.

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
- JVMS/core differential and conformance testing; JCK is relevant only to a
  separately declared Java SE compatibility profile. ~8+ weeks.

## Performance targets

| Milestone | fib(35) target | Notes |
|---|---|---|
| Today (A4) | ~44 ms | achieved; on par with HotSpot JIT |
| R1 complete | ~44 ms | (same; interpreter for new classes) |
| R5 (real-program AOT) | near hand-written Go | whole-program native |
| R6 (optimized) | ≤ hand-written Go | escape analysis + layout opt |

## Current boundary and future compatibility

catty currently has a bounded, evidence-backed Java 25-compatible execution
surface, including an optional real-`java.base` smoke path. It does not claim
to run arbitrary Java applications, Minecraft, arbitrary Java SE libraries, or
JNI-dependent software. Such software may require optional Java SE
compatibility profiles plus host-service providers; those are not automatic
milestones of the Catty Runtime Profile.

The AOT path proves the performance thesis. JVMS semantics, the small runtime
profile, host interoperation, and any opted-in library compatibility remain
separate research and implementation tracks.
