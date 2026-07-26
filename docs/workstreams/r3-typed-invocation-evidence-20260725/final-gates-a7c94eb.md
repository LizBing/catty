# Typed-invocation final-gate record — candidate `a7c94eb`

**Date:** 2026-07-26
**Candidate:** `a7c94eb8e81fcd613341d2fab5e29c1464576eeb`
**Status:** Pass evidence — Owner requested a corrected successor candidate

This record is the final-gate evidence for the fixed typed-invocation candidate.
It is deliberately fail-closed: a gate is not recorded as passing unless its
complete required output exists.

| Gate | Result | Evidence / note |
|---|---|---|
| Static analysis | Pass | `go vet ./...` |
| Unit/regression | Pass | `go test ./...` |
| Race regression | Pass | `go test -race ./...`; the `transpile` package was also rerun separately after the aggregate runner stopped displaying output, and passed |
| Differential fixtures | Pass | `bash tests/run.sh`: 10/10 fixtures |
| R2 concurrency matrix, 1× | Pass | 19/19 Interpreter match; 19/19 IR match; 19/19 AOT NO-BUILD in [`results.txt`](../r2-concurrency-candidate-evidence/a7c94eb8e81fcd613341d2fab5e29c1464576eeb/results.txt) |
| R2 concurrency matrix, 100× race-built | **Pass** | 19/19 Interpreter match; 19/19 IR match; 19/19 AOT NO-BUILD in [`results-stress-100x.txt`](../r2-concurrency-candidate-evidence/a7c94eb8e81fcd613341d2fab5e29c1464576eeb/results-stress-100x.txt) |
| Evidence isolation | Pass | This candidate wrote only under [`r2-concurrency-candidate-evidence/a7c94eb8e81fcd613341d2fab5e29c1464576eeb/`](../r2-concurrency-candidate-evidence/a7c94eb8e81fcd613341d2fab5e29c1464576eeb/) |
| Governance diff check | Pass | `git diff --check` at the fixed candidate before evidence recording |

## 100× gate execution record

The required command completed successfully on 2026-07-26:

```sh
R2_CONCURRENCY_STRESS=100 R2_STRESS_CONCURRENCY=4 \
  bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh \
  a7c94eb8e81fcd613341d2fab5e29c1464576eeb
```

**Execution environment:** macOS darwin/arm64, Go 1.26.5, Temurin 25.0.3

**Header verification:**
- candidate-full: `a7c94eb8e81fcd613341d2fab5e29c1464576eeb`
- build-commit: `a7c94eb8e81fcd613341d2fab5e29c1464576eeb`
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

Historical incomplete evidence files preserved in the same directory:
- `results-stress-100x.incomplete-20260726*.txt` (8 diagnostic artifacts)
- `results-stress-100x.txt.incomplete-20260726T072827Z.16878` (killed artifact)

No detached worktree remained after execution.

## Workstream disposition

All gates required of candidate `a7c94eb` passed. During Owner review on
2026-07-26, Owner accepted this evidence but requested three bounded code
corrections before Done. The workstream therefore returned to In Progress and
this record remains historical evidence for `a7c94eb`, not acceptance evidence
for the corrected successor candidate.
