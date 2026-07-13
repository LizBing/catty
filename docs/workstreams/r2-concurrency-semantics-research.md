# R2 concurrency semantics research

**Status:** Done
**Type:** research
**Review:** owner
**Base commit:** `63d5658`
**Roadmap item:** Phase R2 — Runtime semantics and concurrency planning
**Governing ADRs:** ADR-0016 through ADR-0025, excluding unused ADR-0026; ADR-0027 where String output is used by fixtures
**Prerequisites:** R2 initialization and bounded UTF-16 String workstreams Done
**Acceptance anchor:** N/A — research-only; any follow-up implementation contract requires its own accepted-contract anchor

## Outcome

Produce Java 25 differential evidence, mechanism experiments, and Proposed decisions for
catty's Thread identity/lifecycle, monitors and wait sets, cross-thread class initialization,
and the minimum memory-ordering boundary needed by the first concurrency implementation
slice. Deliver a separate bounded implementation contract; do not add production
concurrency in this workstream.

## Scope

- Establish a Temurin 25 differential fixture set covering the minimum observable
  concurrency surface: stable `currentThread` identity, start-once/lifecycle, join,
  daemon/liveness, synchronized block and method behavior, monitor reentrancy and
  ownership failures, `wait`/`notify`/`notifyAll`, interruption at blocking points,
  cross-thread class initialization, and representative start/join, monitor, volatile,
  and final-field visibility obligations.
- Record Interpreter, IR, and AOT baseline behavior independently. The harness must use
  deterministic coordination where possible, bound every subprocess, preserve stdout,
  stderr, and exit status, and classify unsupported paths explicitly.
- Map the current execution-context, synthetic `java.lang.Thread`, object/monitor,
  method-flag, field-storage, class-initialization, exception, launcher, and AOT bridge
  boundaries. Identify every singleton or mutable runtime structure that becomes shared.
- Compare candidate Go mechanisms only after stating the Java-visible contract. At
  minimum evaluate carrier ownership, Java Thread identity, monitor ownership and
  reentrancy, wait-set signaling, interruption races, VM liveness, and race-free
  publication. Prototypes remain research artifacts and must be exercised with the Go
  race detector where applicable.
- Propose the durable decisions required before production work. The proposal(s) must
  distinguish Java Thread objects from execution carriers, define the monitor/wait-set
  and interruption contracts, state the supported Java memory-ordering boundary, extend
  ADR-0025 for cross-thread initialization, and define cross-engine behavior under
  ADR-0016.
- Define the smallest production implementation slice on the path to the R2
  producer-consumer milestone, including engine matrix, explicit non-support behavior,
  semantic constraints, fixture set, immutable evidence destination, and acceptance
  gates.

## Non-scope

- Production goroutines, locks, atomics, monitor state, Thread lifecycle methods,
  `wait`/`notify`, interrupt delivery, volatile accessors, or cross-thread class-init
  changes.
- Selecting `Thread == goroutine`, adopting the Go memory model as Java semantics, or
  treating `sync.Mutex`, `sync.Cond`, channels, or Go atomics as correct without
  Java-visible evidence.
- `jdk.internal.misc.Unsafe`, `sun.misc.Unsafe`, broad `java.base` enablement,
  `java.util.concurrent`, thread pools, virtual threads, `ThreadLocal`, reflection,
  `invokedynamic`, I/O, JNI, or a complete Java Memory Model implementation.
- Promoting a research prototype into production because it passes a fixture or race
  run. Follow-up production work requires Owner-accepted ADRs and an Accepted
  implementation workstream.
- Performance or scale claims such as one goroutine per Java Thread. Measurements may
  inform a decision but do not expand the supported contract.

## Semantic constraints

- Java SE/JLS/JVMS 25 define supported behavior; repository-pinned Temurin 25 is the
  differential reference under ADR-0017.
- A Java Thread is a Java semantic object. A goroutine is only a candidate carrier under
  ADR-0018; carrier reuse or replacement must not change Java identity, lifecycle,
  interruption, monitor ownership, or class-initialization ownership.
- Monitor operations are per Java object, reentrant for their owning Java Thread, and
  include Java-specified failure behavior. `wait` participates in the same monitor's wait
  set and must release then reacquire ownership according to Java 25.
- The research must identify required happens-before edges rather than substituting the
  Go memory model. Catty-internal Go data races are defects even for racy Java programs.
- Cross-thread initialization preserves ADR-0025's four semantic states and adds waiting,
  wakeup, failure, and visibility behavior behind the same shared service; engines may not
  create separate initialization state machines.
- Interpreter, IR, and AOT capability is reported separately as `Supported`, `Fallback`,
  or `Not implemented`. Research baselines may mismatch Temurin, but no mismatch may be
  relabeled as support.
- Deterministic tests and bounded schedules are primary evidence. Stress/litmus runs are
  supplemental and cannot prove the absence of a memory-ordering bug by themselves.

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| Java 25 differential baseline | `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-research.sh` plus immutable baseline under `docs/workstreams/r2-concurrency-evidence/baseline-63d5658/`; covers every Scope category and fails closed on missing tools, timeouts, or omitted engine rows | Pass — authoritative additive v5 has 19/19 Temurin rows; Interpreter 0/19, IR 0/19, AOT 0/19 (15 `NO-BUILD`, four built then failed) |
| Current-runtime map | `docs/workstreams/r2-concurrency-evidence/reports/runtime-boundary-map.md`; identifies mutable shared state and the Interpreter/IR/AOT transition points with file/line evidence | Pass — five dependency clusters and required implementation order recorded |
| Java semantic contract | `docs/workstreams/r2-concurrency-evidence/reports/java25-concurrency-contract.md`; records lifecycle, monitor, wait-set, interruption, liveness, class-init, and minimum memory-ordering obligations with authoritative references | Pass — Java 25 Thread/monitor/wait/interrupt/JMM/final/class-init obligations mapped |
| Mechanism experiments | `docs/workstreams/r2-concurrency-evidence/reports/go-mechanism-experiments.md` plus bounded prototypes/tests; records exact commands, exit codes, race results, limitations, and rejected alternatives | Pass — race count 100 and stress count 1000 pass; rejected shortcuts and limitations recorded |
| Decisions and production proposal | One or more linked Proposed ADRs plus one Proposed implementation workstream; all open questions that can change the first slice are resolved or made explicit prerequisites | Pass — Proposed ADR-0028/0029/0030 and `r2-thread-monitor-foundation-slice` |
| Governance consistency | `git diff --check` and proof that completed R2 historical evidence directories are unchanged | Pass — `git diff --check`, harness syntax/count checks, and historical-directory diff check exit 0 |

Results use only `Pass`, `Fail`, `Not run`, or `Not implemented`.

## Amendments

Accepted changes are appended here after Owner approval; the frozen contract is not
rewritten to reduce gates.

## Acceptance record

Accepted by Owner on 2026-07-13. The accepted research contract authorizes evidence,
fixtures, bounded research prototypes, Proposed ADRs, and a Proposed implementation
contract within this document's Scope, Non-scope, Semantic constraints, and Acceptance
gates. It does not authorize production concurrency implementation or integration actions.

---

## Investigation findings

- Interpreter and IR currently pop `monitorenter`/`monitorexit` operands without null,
  ownership, or locking behavior. The synthetic Object facade does not declare
  `wait`/`notify`; registry handlers for real declarations are no-ops, and the registered
  `Thread.holdsLock` handler always returns false.
- The synthetic `java.lang.Thread` exposes only a no-op constructor. Although the native
  registry contains a `currentThread` handler, the facade does not declare it, and that
  handler allocates a fresh Java object per invocation. Stable Thread identity is not
  represented yet.
- `rtda.Thread` owns an unsynchronized frame stack and pending-exception signal. Its
  execution-context sequence is incremented without synchronization; this was sufficient
  only for the completed single-context initialization slice.
- Java objects and static/array storage expose mutable `[]Slot` backing with no monitor or
  volatile/final access boundary. Runtime field metadata retains access flags, but the
  current public API does not classify volatile fields, and runtime method flags do not
  classify `ACC_SYNCHRONIZED`.
- The AOT runtime bridge stores one package-global loader and one package-global
  `*rtda.Thread`; fallback calls, native calls, exceptions, and class initialization all
  use that singleton. A concurrency design must replace or scope this boundary rather
  than infer a current goroutine automatically.
- ADR-0025 already records an initialization owner and deliberately leaves a
  synchronization boundary. Cross-thread waiting can extend that service, but the
  current `Class` state and owner fields are not synchronized.
- The classloader cache is an unprotected map, so concurrent first load can race or create
  duplicate runtime Class identities. Class-mirror helpers also allocate a new mirror per
  request, which cannot support static synchronized monitor identity.
- Unsafe-backed class-library compatibility is a separate dependency problem. Combining
  it with Thread/monitor/JMM research would make both the semantic decision and the first
  production slice unbounded.

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — freeze fixture categories, harness policy, and current three-engine baseline | Complete | 19-fixture v5 baseline, 19/19 rows |
| B — map Java 25 semantics and current runtime/engine gaps | Complete | runtime-boundary and Java 25 contract reports |
| C — run bounded Go-mechanism prototypes and race/liveness experiments | Complete | race count 100; stress count 1000 |
| D — propose ADRs and the first bounded concurrency implementation contract | Complete | ADR-0028/0029/0030 + proposed Thread/monitor slice |

Status uses `Pending`, `In progress`, or `Complete`.

---

## Handoff

- **Branch / candidate:** `main`; research artifacts are fixed in the 2026-07-14 governance/acceptance-anchor commit
- **Acceptance anchor / base:** N/A for research; baseline `63d5658`
- **Dirty files:** research fixtures/evidence/reports/prototypes and Proposed decisions/contracts; pre-existing untracked `.claude/worktrees/` is Owner state and untouched
- **Historical evidence check:** Pass — completed initialization/String/original R2 evidence directories have no diff
- **Candidate evidence path:** N/A; research baseline is fixed to `docs/workstreams/r2-concurrency-evidence/baseline-63d5658/`
- **Last location:** all four research slices and all six research acceptance gates complete
- **Checks run / not run:** 19-fixture v5 baseline Pass; prototype race/stress Pass; harness syntax/count and `git diff --check` Pass; production regression gates not run because production code is unchanged
- **Blocker:** None — Owner accepted the research conclusions, ADR-0028/0029/0030, and the successor implementation contract on 2026-07-14
- **Next action:** successor `r2-thread-monitor-foundation-slice` begins from its acceptance anchor
- **Non-derivable context:** keep Unsafe/java.base reachability separate unless Owner explicitly changes the R2 ordering
