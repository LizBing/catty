# R2-GATE: Concurrency and Unsafe architecture gates

**Status:** Closed
**Owner:** Codex (architecture only)
**Reviewer:** LizBing
**Integrator:** Unassigned
**Base commit:** `0b82986444d79b957468fda68595a2c279807497`
**Branch:** `main` (draft only; no runtime implementation authorized)
**Target milestone:** R2 readiness

## Outcome

Turn the four open R2 risks into accepted semantic contracts before any Thread,
monitor, volatile, Unsafe, or broad native implementation begins. Closing this
workstream authorizes later implementation contracts; it does not itself claim
R2 functionality.

Catty's public promise remains Java/JRE observable semantics. Go's goroutines,
atomics, locks, GC, and compiler are implementation mechanisms, not permission
to expose Go semantics where Java specifies different behavior.

## Decisions requested from LizBing

The recommendations below are the proposed architecture. Each must be accepted,
changed, or rejected before this workstream can become **Ready**.

| Gate | Recommended decision | Why it is the safe default |
|---|---|---|
| G1 — Memory model | Supersede ADR-0011. Preserve DRF-Java plus final/volatile/monitor/Thread/class-init semantics. Compare Strict, Go-native, and Hybrid storage before deciding whether to waive full JMM outcomes for racy programs. | The product goal is JRE semantics, but a named, measured deviation may be justified by material application-level performance. |
| G2 — Missing natives | Strict failure by default: a declared but unresolved native throws `UnsatisfiedLinkError`. No generic zero/null stub. Explicit, reviewed semantic no-ops remain registered by name. | Plausible zero values silently corrupt program behavior and invalidate compatibility evidence. |
| G3 — Unsafe | Use opaque logical offsets translated by Catty; never expose Go addresses. Implement only behavior-backed profiles and throw strictly outside them. | Go GC may move implementation details, object layout is not OpenJDK layout, and raw pointer arithmetic would bypass GC and type safety. |
| G4 — Acceptance | Require deterministic differential programs, repeated stress, timeout/deadlock detection, and `go test -race`; a curated PASS count alone is insufficient. | Concurrency failures are schedule-sensitive and can survive ordinary unit tests. |

## First-principles invariants

These invariants constrain every later R2 implementation contract.

1. **No accidental Go data races.** Catty must not accidentally implement Java shared fields, static
   fields, arrays, thread state, class initialization, or monitor state using
   unsynchronized concurrent Go reads and writes. A research-only Go-native
   backend may intentionally model a Java-level race, but it must be isolated,
   classified, and excluded from production until its outcomes and safety are
   understood. `go test -race` remains a production release gate.
2. **Allowed-behavior subset.** An implementation may initially be stronger
   than the JMM (for example, sequentially consistent accesses) if every
   behavior it produces is legal Java behavior. It may not weaken a specified
   happens-before edge.
3. **One semantic access layer.** Interpreter, IR executor, AOT code, native
   methods, and Unsafe operations must reach shared Java state through the same
   semantic access abstraction. AOT may bypass it only after a proven
   thread-local/non-escaping optimization.
4. **Explicit Java thread identity.** A goroutine executing Java code carries a
   specific `*rtda.Thread`; monitor ownership and `currentThread` use that
   object. Catty must not depend on an unavailable or unstable Go goroutine ID.
5. **Java exceptions, not Go accidents.** Illegal monitor use, double start,
   interruption, unresolved natives, invalid Unsafe offsets, and unsupported
   operations surface through Java exception semantics on every engine.
6. **No fabricated compatibility.** Unsupported Unsafe/off-heap/JNI behavior
   fails loudly. It must never return a value chosen only to let bootstrap code
   continue.

## G1 — JMM and concurrency semantic contract

### Proposed decision

ADR-0011 (“Adopt Go memory model”) is rejected in its current form and must be
superseded before implementation. LizBing accepted the following replacement
direction on 2026-07-12:

> Catty preserves Java semantics for data-race-free programs and for final,
> volatile, monitor, Thread, class-initialization, and selected Unsafe/VarHandle
> operations. Before promising full outcomes for incorrectly synchronized Java
> programs, Catty measures Strict, Go-native, and Hybrid implementations. Any
> waiver must name the lost behavior and demonstrate material macro-level gain.

The JLS permits any execution strategy whose observable executions are legal
under the JMM. R2 therefore starts with a conservative implementation and
optimizes only after proof.

“99.9% compatible” is not an acceptance metric: it has no defined workload
denominator. The investigation reports semantic classes, forbidden outcomes,
real-workload prevalence, and measured performance. A waiver is eligible only
if it does not weaken the protected semantics above and yields at least one of:

- 2× improvement in the affected synchronization/access microbenchmark;
- 15% application throughput improvement; or
- 20% application p99 latency reduction.

Smaller gains keep the Strict/Hybrid semantics. These thresholds are decision
triggers, not automatic approval; LizBing still accepts each named deviation.

Primary references:

- [JLS 25 §17.1–17.5 — monitors, wait sets, memory model, final fields](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html)
- [JLS 25 §12.4.2 — detailed class initialization procedure](https://docs.oracle.com/javase/specs/jls/se25/html/jls-12.html#jls-12.4.2)
- [JVMS 25 §2.11.10 and §6.5 — implicit/explicit monitor operations](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-6.html)

### Required happens-before edges

The first R2 implementation must preserve at least these specified edges:

- program order within one Java thread;
- monitor unlock → every subsequent lock on the same monitor;
- volatile write → every subsequent volatile read of that variable;
- successful `Thread.start()` → first action of the started thread;
- final action of a thread → successful termination observation through
  `join()`/`isAlive()`;
- interrupt action → observation through `InterruptedException`,
  `interrupted()`, or `isInterrupted()`;
- class initialization completion → subsequent active use of the class;
- constructor/final-field publication rules from JLS §17.5.

Passing these examples is not a claim of full JMM/JCK conformance. It is the R2
compatibility floor; causality and broader litmus coverage remain tracked until
formal conformance evidence exists.

### Shared-memory representation

Current `rtda.Slot` storage cannot be concurrently accessed directly. R2 must
introduce a semantic heap-access layer covering:

- instance fields;
- static fields;
- primitive and reference array elements;
- long/double values as one logical access;
- synthetic payloads that are shared across threads.

The initial implementation may serialize accesses conservatively (per object,
per class, per array, or via atomic cells). The exact representation is an
implementation benchmark decision, provided that:

- Go's race detector sees no race;
- reference writes remain visible to Go GC;
- volatile operations are sequentially consistent or otherwise proven to
  satisfy Java volatile ordering;
- plain Java accesses still produce only JMM-legal executions;
- interpreter and emitted AOT code cannot disagree.

Performance recovery belongs to a later measured optimization slice. Escape and
thread-local analysis may remove synchronization only when it proves that no
other Java thread can observe the cell.

### Thread runtime state

Each Java `Thread` maps one-to-one to a goroutine while it is running, but the
Java object owns an explicit runtime payload:

- lifecycle: `NEW → RUNNABLE → TERMINATED` plus observable blocked/wait states;
- associated `*rtda.Thread` execution context;
- completion signal used by `join` and VM shutdown;
- interrupt status with atomic observation/clear operations;
- park permit (at most one), separate from interrupt status;
- daemon flag and uncaught termination outcome;
- reference to the Java `Thread` object returned by `currentThread`.

`start` is exactly-once and throws `IllegalThreadStateException` on a second
start. VM process lifetime waits for live non-daemon Java threads after `main`
returns. A goroutine is not itself used as thread identity.

### Monitor contract

`sync.Mutex` alone is not a Java monitor because Java monitors are reentrant and
have ownership, recursion depth, and a wait set. Each Java object therefore has
a lazily allocated monitor abstraction containing at least:

- owner `*rtda.Thread`;
- recursion depth;
- entry contention mechanism;
- wait set with individually wakeable waiters;
- timed wait and interrupt integration.

Required behavior:

- reentrant acquire/release;
- `IllegalMonitorStateException` for non-owner exit/wait/notify;
- `wait` atomically records the recursion count, releases all acquisitions,
  blocks, then reacquires and restores the count before returning/throwing;
- `notify` selects one waiter without promising fairness; `notifyAll` releases
  all waiters to contend for reacquisition;
- interrupt during wait removes/wakes the waiter and clears status when
  `InterruptedException` is delivered;
- timed waits may wake spuriously but not before the requested timeout solely
  because the timer expired early;
- `sleep` does not release monitors and creates no synchronization edge.

Both bytecode `monitorenter`/`monitorexit` and `ACC_SYNCHRONIZED` methods use
this abstraction. Synchronized methods release the monitor on normal return and
on every abrupt exit, including Java exceptions.

### Class initialization under concurrency

The current boolean `initStarted` is insufficient once multiple threads exist.
R2 must implement the JLS §12.4.2 protocol:

- one initialization lock per class;
- at most one initializing thread;
- recursive request by that same thread is permitted;
- other threads wait for completion;
- success publishes initialized static state;
- failure records erroneous state and propagates the specified initialization
  errors to current and subsequent users.

## G2 — Strict native resolution contract

### Proposed decision

Remove the generic return-type-based `rtda.nativeStub`. A real class-file method
with `ACC_NATIVE` remains unresolved until the native registry attaches a
specific implementation. Invoking an unresolved declaration throws Java
`UnsatisfiedLinkError` with class, method, and descriptor context.

This rule applies identically to:

- interpreter invocation;
- IR executor invocation;
- AOT/runtime bridge invocation.

The loader may discover and report unresolved natives, but discovery must not
change execution semantics.

### Native classifications

Every registered native is classified in code or generated inventory as one of:

| Classification | Meaning | Example |
|---|---|---|
| Implemented | Catty implements specified observable behavior | `System.arraycopy` |
| Semantic no-op | The Java contract permits no externally required effect | `registerNatives`; possibly `System.gc` request semantics |
| Compatibility adapter | Different mechanism, equivalent Java behavior | time, identity hash, thread operations |
| Unsupported | No implementation; invocation throws `UnsatisfiedLinkError` | raw off-heap access until explicitly implemented |

`nop`, `nopBool0`, `nopRef0`, and zero-return helpers require an audit. Their
names are not evidence that no-op/zero behavior is legal. In particular,
`Object.wait`, `notify`, `notifyAll`, `Thread.holdsLock`, and Runtime memory
queries cannot retain placeholder behavior under strict mode.

### Discovery mode

If dependency exploration needs aggregation, add a diagnostics-only mode that
records unresolved signatures and still throws at the first invocation. Do not
provide a generic “continue with zero” flag. Any temporary behavior shim must be
an explicit signature registration, test, owner, and removal criterion.

## G3 — Minimum Unsafe semantic profiles

### Proposed decision

Unsafe offsets are opaque Catty logical tokens. They identify a Java storage
location; they are not byte addresses into Go objects and must never be
convertible to `unsafe.Pointer` by Java code.

The token/translation layer must distinguish:

- instance field + declaring class + logical slot;
- static field + class storage + logical slot;
- array component kind + logical index;
- invalid/stale/unsupported tokens.

For array APIs whose Java bytecode performs `base + index * scale`, Catty may
return stable synthetic base/scale values and translate the resulting logical
offset back to an element. These values are process-internal compatibility
tokens, not promises about Go layout.

### Profile U0 — bootstrap and layout discovery

Required to investigate and unblock current JDK 25 bootstrap paths:

- singleton/getUnsafe path and legal `registerNatives` no-op;
- `objectFieldOffset(Class, String)` where no reflection object is required;
- `arrayBaseOffset`, `arrayIndexScale`;
- address size/page size/endianness constants only where the caller's behavior
  can be implemented honestly.

Before implementation, produce a checked-in call graph from the exact pinned
Temurin 25 classes used by Integer/Long conversion, floating parsing, and the
selected HashMap operations. “Approximately 50 methods” is not a scope.

### Profile U1 — heap access

- plain and volatile get/put for reference, int, long, boolean, byte, short,
  char, float, and double over logical instance/static/array locations;
- null base or absolute-address forms remain unsupported unless separately
  contracted;
- all operations reuse the G1 semantic heap layer.

### Profile U2 — atomics

- compare-and-set for reference, int, and long first;
- compare-and-exchange and weak forms only when demanded by the checked call
  graph, with their documented plain/acquire/release/volatile semantics;
- atomic read-modify-write helpers may be real JDK bytecode built on the above;
- comparison and update operate on one logical cell and are linearizable.

### Profile U3 — fences

- full, acquire/load, and release/store fence behavior sufficient for the JDK
  methods reached by the acceptance corpus;
- fence operations participate in the same ordering design as volatile and
  CAS; they are not empty natives merely because Go atomics are strong.

### Profile U4 — parking

- one-permit `park`/`unpark` semantics tied to Java Thread runtime state;
- `unpark` before `park` makes one future park return immediately;
- repeated `unpark` does not accumulate more than one permit;
- interrupt wakes/permits return according to Java LockSupport/Unsafe behavior;
- timeout uses monotonic duration for relative waits and documented wall-clock
  interpretation for absolute waits.

### Explicitly outside the initial Unsafe profile

Until separate contracts exist, these throw strictly rather than emulate raw
memory with fabricated values:

- arbitrary absolute-address get/put;
- `allocateMemory`/`reallocateMemory`/`freeMemory`;
- bulk raw memory copy/set and swap;
- `defineClass`, anonymous/hidden class injection;
- raw monitor APIs, load-average probes, writeback/cache-line operations;
- arbitrary object-header or compressed-reference assumptions.

If a representative R2 program reaches one of these, the architecture gate is
reopened or that program is explicitly assigned to a later milestone; it is not
silently removed from the evidence set.

## G4 — Acceptance and evidence contract

### Test topology

R2 implementation contracts must add a dedicated concurrency harness rather
than putting nondeterministic checks into the existing stdout-only fixture loop.
Every test has a wall-clock timeout, deterministic synchronization where
possible, repeat count for schedule-sensitive paths, and HotSpot reference
behavior.

The matrix must cover both pure-synthetic startup where meaningful and the real
Temurin 25 `java.base` path. Interpreter, IR, AOT/runtime bridge are included
only when that engine claims the feature; unsupported engines must fail
explicitly rather than silently skip semantics.

### Required representative programs

| Program | Observable contract |
|---|---|
| StrictNativeProbe | A declared unresolved native throws catchable `UnsatisfiedLinkError`; no zero/null continuation |
| StartJoinPublication | Writes before `start` visible in child; child writes visible after `join`; second start fails |
| MonitorCounter | Mutual exclusion, reentrancy, static and instance synchronized methods, release on exception |
| WaitNotifyProtocol | ownership errors, wait release/reacquire, notify and notifyAll, timeout, guarded-loop behavior |
| InterruptProtocol | interrupt flag observe/clear; wait/sleep/join interruption and `InterruptedException` clearing |
| VolatilePublication | volatile message passing and no stale payload after publication |
| AtomicCAS | reference/int/long CAS and linearizable multi-thread increment |
| ParkPermit | unpark-before-park, one-permit behavior, timeout, interrupt wakeup |
| ClassInitRace | one `<clinit>`, publication to all users, recursive init, failure propagation |
| NonDaemonLifetime | VM remains alive after main while a non-daemon Java thread runs; daemon alone does not keep it alive |
| UnsafeJavaBaseSmoke | exact Integer/Long conversion, floating parsing cases, and representative HashMap put/get/remove/collision/null behavior |

Concurrency tests must not use `sleep` as their sole correctness coordination.
Stress loops supplement deterministic barriers; they do not replace semantic
assertions.

### Required commands/evidence

At minimum, each implementation merge records:

```sh
gofmt -w <changed-go-files>
go vet ./...
go test ./...
go test -race ./...
bash tests/run.sh
```

The R2 harness runs each relevant program against pinned Temurin 25 and Catty,
with per-test timeout and repeated stress mode. CI must run a bounded version;
a longer stress job may run separately. Failure artifacts include seed,
iteration, engine, thread dump/state snapshot where available, and stdout/stderr.

### Performance guardrails

Correctness gates precede optimization, but R2 records costs so conservative
serialization does not become invisible debt:

- R1 startup and `fib(35)` remain regression baselines;
- uncontended monitor enter/exit, volatile access, CAS, thread start/join, and
  park/unpark get microbenchmarks against the previous Catty commit and Go
  primitives;
- no performance target may justify removing a specified Java ordering edge;
- optimization requires before/after profiles and unchanged semantic/race tests.

## Proposed architecture outputs

Closing R2-GATE should create or approve these durable decisions:

1. **ADR-0016:** Supersede ADR-0011 — implement JMM-observable semantics with a
   conservative race-free shared-memory layer.
2. **ADR-0017:** Explicit Java Thread context and reentrant monitor state
   machine on goroutines.
3. **ADR-0018:** Strict unresolved-native failure and audited native
   classifications.
4. **ADR-0019:** Unsafe logical offsets and phased semantic profiles.
5. **R2 test plan:** named programs, harness shape, CI/stress split, and exact
   pinned JDK evidence.

These ADR numbers are reserved by this draft but are not created or Accepted
until LizBing approves the corresponding decisions.

## Implementation sequence after gate approval

Later workstreams should be small and independently reviewable:

| Block | Result | Prerequisites |
|---|---|---|
| R2-A | Strict native resolution + inventory + StrictNativeProbe | ADR-0018 |
| R2-B | Race-free heap access foundation + volatile/final/class-init rules | ADR-0016, test harness skeleton |
| R2-C | Explicit Thread lifecycle, start/join/current/sleep/interrupt, VM liveness | R2-B, ADR-0017 |
| R2-D | Reentrant monitors, bytecodes, synchronized methods, wait/notify | R2-C |
| R2-E | Unsafe U0/U1 logical offsets and heap access; java.base numeric smoke | R2-A/B, ADR-0019, call graph |
| R2-F | Unsafe U2/U3 atomics and fences | R2-E |
| R2-G | Unsafe U4 park/unpark and broader concurrency smoke | R2-C/F |
| R2-H | Integrated stress, race, performance, docs, and R2 closure audit | R2-A–G |

An implementation block may be reassigned between Claude and Codex, but each
has one owner and one reviewer. DeepSeek is the default for bounded slices;
GLM is reserved for high-semantic-risk debugging when manually selected by
LizBing through CC Switch.

## Workstream acceptance gates

- [x] LizBing decides G1: protected Java semantics plus measured Strict/Go-native/Hybrid study.
- [x] LizBing decides G2: strict `UnsatisfiedLinkError`, no generic zero stub.
- [x] LizBing decides G3: logical Unsafe offsets and U0–U4 profile boundaries.
- [x] LizBing decides G4: representative differential/race/stress evidence.
- [x] ADR-0016 through ADR-0019 are drafted.
- [x] ADR-0016 through ADR-0019 are reviewed and Accepted.
- [x] Exact JDK 25 call graphs for selected Unsafe-dependent paths are attached.
- [x] R2 test-plan contract names commands, timeouts, repetitions, and engine coverage.
- [x] No runtime implementation started before the above decisions were Accepted.

## Evidence

| Gate | Command/artifact | Result |
|---|---|---|
| Existing architecture | ADR-0010/0011, current `rtda`, native, interpreter paths | Reviewed |
| Current native behavior | `rtda.nativeStub`, registry and loader inspection | Generic zero/null fallback confirmed |
| Current concurrency behavior | Thread/Object/Class structures and opcode docs | Single-threaded/nop monitor baseline confirmed |
| JDK baseline | `java -version` | Temurin 25.0.3+9 LTS |
| Unsafe surface | `javap -p jdk.internal.misc.Unsafe` | Large mixed surface; profile scoping required |
| Java semantics | JLS/JVMS/Java SE 25 references linked above | Reviewed |
| Product decisions | LizBing review | G1–G4 accepted on 2026-07-12 |
| Caller graph | `docs/research/R2_UNSAFE_CALL_GRAPH.md` | Attached; corrects grouped Unsafe assumption |
| Test plan | `docs/R2_TEST_PLAN.md` | Drafted with deterministic/CI/stress tiers |
| ADR acceptance | LizBing review, 2026-07-12 | ADR-0016 through ADR-0019 Accepted |

## Risks and rollback

- Risk: conservative shared-memory access damages the AOT performance thesis.
  Containment: measure it, then remove synchronization only with escape or
  thread-local proof; never optimize by reintroducing Go data races.
- Risk: logical Unsafe offsets fail code that assumes real HotSpot addresses.
  Containment: such code is outside the initial profile and fails explicitly;
  compatibility expansion requires a separate contract.
- Risk: scope expands from concurrency into full reflection/off-heap/JNI.
  Containment: U0–U4 are caller-driven; unsupported profiles remain strict.
- Risk: schedule-sensitive tests become flaky.
  Containment: deterministic protocols first, bounded stress second, explicit
  seeds/timeouts, no sleep-only assertions.
- Rollback: this workstream is documentation-only. Rejecting a gate changes or
  supersedes the draft before any runtime code is authorized.

## Handoff history

| Date | From | To | Commit | Summary |
|---|---|---|---|---|
| 2026-07-12 | Codex | LizBing | Uncommitted draft | Four R2 architecture gates proposed for review |
| 2026-07-12 | LizBing | Codex | Decision recorded in draft | G1–G4 directions accepted; architecture artifacts authorized |
| 2026-07-12 | LizBing | Codex | Closure commit containing this row | ADR-0016–0019 accepted; R2-GATE closed |
