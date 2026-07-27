# Typed-invocation final-gate record — candidate `047ffca`

**Date:** 2026-07-26
**Candidate:** `047ffca5b0126262865e13884dc48528b74a2985`
**Status:** Pass — Owner accepted; workstream Done

This record is the final-gate evidence for the corrected typed-invocation candidate.
It is deliberately fail-closed: a gate is not recorded as passing unless its
complete required output exists.

This candidate contains three bounded corrections requested by Owner on
2026-07-26 within the frozen contract: reference-argument assignment checks,
IR internal-failure stack restoration, and declared-only constructor resolution.
No ADR or new workstream was required.

| Gate | Result | Evidence / note |
|---|---|---|
| Static analysis | Pass | `go vet ./...` |
| Unit/regression | Pass | `go test ./...` |
| Race regression | Pass | `go test -race ./...` |
| Differential fixtures | Pass | `bash tests/run.sh`: 10/10 fixtures |
| R2 concurrency matrix, 1× | Pass | 19/19 Interpreter match; 19/19 IR match; 19/19 AOT NO-BUILD in [`results.txt`](../r2-concurrency-candidate-evidence/047ffca5b0126262865e13884dc48528b74a2985/results.txt) |
| R2 concurrency matrix, 100× race-built | **Pass** | 19/19 Interpreter match; 19/19 IR match; 19/19 AOT NO-BUILD in [`results-stress-100x.txt`](../r2-concurrency-candidate-evidence/047ffca5b0126262865e13884dc48528b74a2985/results-stress-100x.txt) |
| Evidence isolation | Pass | This candidate wrote only under [`r2-concurrency-candidate-evidence/047ffca5b0126262865e13884dc48528b74a2985/`](../r2-concurrency-candidate-evidence/047ffca5b0126262865e13884dc48528b74a2985/) |
| Governance diff check | Pass | `git diff --check` at the fixed candidate before evidence recording |

## 100× gate execution record

The required command completed successfully on 2026-07-26:

```sh
R2_CONCURRENCY_STRESS=100 R2_STRESS_CONCURRENCY=4 \
  bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh \
  047ffca5b0126262865e13884dc48528b74a2985
```

**Execution environment:** macOS darwin/arm64, Go 1.26.5, Temurin 25.0.3

**Header verification:**
- candidate-full: `047ffca5b0126262865e13884dc48528b74a2985`
- build-commit: `047ffca5b0126262865e13884dc48528b74a2985`
- build-source: detached immutable candidate snapshot
- stress: 100×, race-build: 1, concurrency: 4
- fixtures: 19

**Summary:**
- interpreter match: 19/19
- IR match: 19/19
- AOT NO-BUILD: 19/19
- result: Pass

**Exit code:** 0

All 19 fixtures match Temurin 25 in both Interpreter and IR across 100 stress
iterations each. AOT correctly rejects all 19 fixtures as NO-BUILD. No race
conditions, timeouts, deadlocks, missing rows, or semantic mismatches were
detected.

No detached worktree remained after execution.

### Relation to prior candidate

Candidate `a7c94eb8e81fcd613341d2fab5e29c1464576eeb` (recorded in
[`final-gates-a7c94eb.md`](./final-gates-a7c94eb.md)) passed all gates but Owner
requested three bounded code corrections. This candidate `047ffca` is the
corrected successor; it supersedes `a7c94eb` as the Ready candidate for Owner
review.

## Workstream disposition

All gates required of candidate `047ffca` passed. Owner accepted the review and
marked the workstream Done on 2026-07-26. The integration commit is recorded in
the workstream acceptance record and project status source.
