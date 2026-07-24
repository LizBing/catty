# R3 execution-context and Java Thread facade boundary slice

**Status:** Ready
**Type:** implementation
**Review:** owner
**Profile:** Catty JVMS Core shared kernel; Java Thread facade remains bounded
**Roadmap item:** Phase R3 — Reflection & dynamic features, shared-kernel track
**Governing ADRs:** ADR-0016, ADR-0018, ADR-0020, ADR-0025,
ADR-0028 through ADR-0031, ADR-0034
**Prerequisites:** `r3-runtime-identity-definition-slice` Done at integration
commit `ee088a0`
**Acceptance anchor:** `ee088a0ce7b1d0ba6b2f97d59f0c07274464373e`; fixed by
the implementation preflight on 2026-07-24

## Outcome

Runtime execution context, Java-visible `java.lang.Thread` facade state, and
internal execution carriers have explicit boundaries in names, APIs, and tests.
Interpreter, IR, native calls, class initialization, monitor ownership, and
future typed invocation consumers receive an explicit execution-context handle
without treating the Java Thread facade object or a goroutine as the context.

This slice preserves the existing bounded Java Thread/monitor behavior and adds
no new Java-visible Thread, reflection, Host ABI, InvokeDynamic, or AOT
concurrency capability.

## Scope

- Introduce a runtime execution-context abstraction in `rtda` that owns frame
  stack, loader, execution-context ID, pending Java throwable, bridge return,
  class-initialization ownership, and monitor ownership identity.
- Keep `rtda.Thread` as a compatibility alias or adapter only where needed
  during migration; production APIs should move toward `ExecutionContext`
  naming at the shared kernel boundary.
- Move Java Thread facade attachment, lifecycle, daemon, interrupt, completion,
  sleep/join wait channels, and monitor-wait enrollment behind a dedicated
  Java Thread runtime/facade sidecar API rather than letting generic execution
  consumers reach those fields directly.
- Preserve `java.lang.Thread.currentThread`, start-once, `isAlive`, untimed
  `join`, interrupt, sleep, daemon, monitor wait/notify, and main-thread
  liveness behavior already verified by R2.
- Update Interpreter, IR, native, launch, runtime bridge, class initialization,
  monitor, and tests to use execution-context terminology where that boundary
  is semantically relevant.
- Document remaining temporary aliases and the exact removal point, if any, in
  the handoff.

## Non-scope

Timed `wait`/`join`; virtual threads; `ThreadGroup`; `ThreadLocal`;
context-class-loader inheritance; priorities; uncaught-exception handler
policy; Java SE broad `java.lang.Thread` compatibility; public Catty Runtime
Profile APIs; Host ABI provider implementation; typed dynamic invocation
value/result implementation; generated classes; InvokeDynamic; AOT
concurrency; removing all AOT runtime globals; multi-VM instance support; or
changing the launcher default.

## Semantic constraints

- Java Thread object identity remains the facade object returned by
  `Thread.currentThread()`; goroutine identity is never observable and never
  used to discover Java identity.
- Execution-context ID remains the owner token for monitors and class
  initialization. Monitor recursion, wait/notify/interrupt ordering,
  synchronized-method release, and class-initialization same-owner/other-owner
  behavior must not change.
- Interpreter and IR must share the same execution-context and facade services.
  Tests must prove the rename/boundary split does not create a second object,
  Class, Thread, monitor, or throwable world.
- AOT concurrency remains precise build rejection. This slice may rename or
  isolate the AOT bridge's single context, but it must not make concurrent AOT
  programs build or fall back through a global context.
- Native/profile code may use Java Thread facade APIs only when it is
  implementing Java Thread behavior; generic runtime, initialization, monitor,
  and future typed-invocation consumers must use execution-context APIs.
- Unsupported facade/native bindings must remain explicit failures under
  ADR-0034; no zero/null success-looking stub may be introduced.

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| Boundary map | Documented map of execution-context state vs Java Thread facade state vs carrier-only state in this workstream | Pass for Slice A inventory; full closure pending |
| API migration | `rg` audit shows generic runtime/engine APIs use execution-context naming; Java Thread facade APIs remain confined to facade/lifecycle/native Thread code or documented temporary aliases | Pass for Slice D: engine/native/launcher/runtime bridge use `ExecutionContext` terminology; remaining `Thread` names are documented temporary compatibility |
| Identity/context | focused tests prove `currentThread`, facade attachment, EC owner ID, monitor ownership, synchronized static mirror, and throwable identity survive the boundary split | Partial: ExecutionContext constructor, ID/EC compatibility, frame identity, bridge return, exception isolation, sidecar facade/lifecycle ownership, currentThread, start/join/interrupt/daemon, existing monitor tests, and focused engine/native compile tests pass; full regression pending |
| R2 concurrency regression | 19-fixture concurrency matrix 1x and race-built 100x: 19/19 Interpreter + IR Match, 19/19 AOT NO-BUILD | Pass at candidate `6dc325c`: 1x 19/19 Interpreter + IR Match, 19/19 AOT NO-BUILD; 100x race-built stress 19/19 Interpreter + IR Match, 19/19 AOT NO-BUILD |
| Core regression | `GOCACHE=/private/tmp/catty-go-cache go test ./...`; `GOCACHE=/private/tmp/catty-go-cache go test -race ./...`; `GOCACHE=/private/tmp/catty-go-cache bash tests/run.sh` | Dirty-tree Pass: `go test ./...`, `go test -race ./...`, `go vet ./...`, and `bash tests/run.sh` 10/10 |
| Capability honesty | R3 24-row baseline remains Interpreter 0/24 Match, IR 0/24 Match, AOT 24/24 NO-BUILD; no new Java-visible R3 row claimed Supported | Dirty-tree Pass: 24/24 rows; Interpreter 0/24 MATCH; IR 0/24 MATCH; AOT 24/24 NO-BUILD |
| Governance | isolated candidate evidence path, historical evidence unchanged, `git diff --check` Pass | Pass: candidate fixed at `6dc325c`; R2 candidate evidence isolated under `docs/workstreams/r2-concurrency-candidate-evidence/6dc325c/`; historical evidence check and `git diff --check` pass |

## Amendments

Accepted 后只在此追加由 Owner 接受的需求变化，不回写降低原合同。

---

## Implementation preflight

开始 production implementation 前记录：

- **Acceptance anchor / actual base:** `ee088a0ce7b1d0ba6b2f97d59f0c07274464373e`
  (current HEAD; worktree descends from this commit)
- **Historical evidence check:** `git diff --quiet ee088a0 -- docs/workstreams/r2-evidence
  docs/workstreams/r2-initialization-evidence docs/workstreams/r2-string-evidence
  docs/workstreams/r2-concurrency-evidence docs/workstreams/r2-concurrency-candidate-evidence
  docs/workstreams/r3-metadata-slice/evidence` — Pass (exit 0)
- **Candidate evidence destination:** `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/`
- **Harness output policy:** all candidate commands must set an explicit output
  directory beneath the candidate evidence destination; no command may write to
  historical/shared evidence paths or overwrite prior candidate evidence.

任何一项缺失时，保持 `Accepted`，不得转为 `In Progress`。

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A. Boundary inventory | Complete | `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/boundary-inventory.md` |
| B. ExecutionContext kernel type/API | Complete | `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/execution-context-design.md`; `rtda.NewExecutionContext`, temporary `rtda.Thread` alias, `ExecutionContext.ID`, `Frame.Context`; focused tests in `rtda/thread_test.go` and `rtda/frame_test.go`; `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go test ./interpreter ./native ./launch ./runtime` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` Pass |
| C. Java Thread facade sidecar/API | Complete | `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/java-thread-sidecar.md`; `rtda.JavaThreadState`; Thread facade `Object.Extra()` stores sidecar; launcher main Thread facade stores sidecar; native Thread facade extracts sidecar; focused tests in `rtda/thread_test.go`, `native/thread_test.go`, and `native/object_test.go`; `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda ./native ./interpreter ./launch ./runtime` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...` Pass |
| D. Engine/native migration | Complete | `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/engine-native-migration.md`; interpreter bytecode/IR/invoke/init/bridge helpers use `*rtda.ExecutionContext`; generic native helpers use `Frame.Context`; launcher and runtime bridge use execution-context naming; compatibility audit leaves only `rtda.Thread` alias, `NewThread`, `Frame.Thread`, `runtime.Thread`, and focused compatibility test; `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda ./interpreter ./native ./launch ./runtime` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` Pass; `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...` Pass |
| E. Regression/evidence closure | Complete | Dirty-tree evidence: `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/slice-e-dirty-20260725/summary.md`; fixed candidate: `6dc325c`; R2 candidate evidence: `docs/workstreams/r2-concurrency-candidate-evidence/6dc325c/results.txt` and `results-stress-100x.txt`; passed `go test ./...`, `go test -race ./...`, `go vet ./...`, `bash tests/run.sh`, R3 baseline, historical evidence check, `git diff --check`, R2 concurrency 1x, and R2 concurrency 100x race-built stress |

---

## Handoff

- **Branch / candidate:** `codex/r3-execution-context-thread-facade-boundary`;
  implementation candidate `6dc325cae30dd0ecfc544c7fc9e3b536f04efefb`
- **Acceptance anchor / base:** `ee088a0ce7b1d0ba6b2f97d59f0c07274464373e`
- **Dirty files:** candidate evidence/status updates after fixing implementation
  candidate; implementation commit is fixed
- **Historical evidence check:** Pass on 2026-07-24; exact command recorded above
- **Candidate evidence path:** `docs/workstreams/r3-execution-context-boundary-evidence/ee088a0-preflight-20260724/`
- **Last location:** Slice E fixed-candidate regression evidence completed; owner
  review is pending.
- **Checks run / not run:** historical evidence check, `git diff --check`,
  `GOCACHE=/private/tmp/catty-review-go-cache go test ./...`,
  `GOCACHE=/private/tmp/catty-review-go-cache go test -race ./...`,
  `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...`,
  `GOCACHE=/private/tmp/catty-review-go-cache bash tests/run.sh`, and R3
  baseline Pass on the dirty worktree. R2 concurrency candidate 1x and 100x
  race-built stress Pass at fixed candidate `6dc325c`.
- **Blocker:** No known technical blocker remains. Because Review is `owner`,
  Done/integration still requires Owner acceptance.
- **Next action:** Owner reviews candidate `6dc325c` and its evidence; if
  accepted, mark Done/integrate per project protocol.
- **Non-derivable context:** This is a boundary-hardening prerequisite for the accepted typed dynamic-invocation kernel, not an implementation of typed values or Java Thread feature expansion.
