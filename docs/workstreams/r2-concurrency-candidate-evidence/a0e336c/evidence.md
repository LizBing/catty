# Slice B Rework Evidence — `a0e336c`

**Candidate:** `a0e336c` (branch: `worktree-r2-thread-monitor-foundation`)
**Base candidate:** `505d3ee` (original Slice B)
**Rework base:** `07d5d16` (Slice B governance)
**Date:** 2026-07-14
**Status:** Complete, awaiting Owner review

## Blockers resolved

### 1. Thread final action exactly-once
- `Terminate()` uses `atomic.CAS(stateRunnable, stateTerminated)` before `close(done)`
- Repeated or concurrent `Terminate()` returns immediately — no panic, no double-close
- VM liveness unaffected: the goroutine carrier calls `Terminate()` once after `DefaultRunLoop`, and the CAS ensures only the first call closes `done`

### 2. Interrupted() drains stale waker signal
- `Interrupted()` drains the waker channel after the atomic Swap
- `Sleep()` re-checks `Interrupted()` on waker instead of unconditionally clearing the flag
- Stale waker signals (left by a previous `Interrupt()` whose flag was already cleared by `Interrupted()`) return normally from `Sleep`
- `Join` pattern already uses `Interrupted()` after waker — no regression
- All interrupt/clear/sleep sequences are race-free

### 3. setDaemon lifecycle rules and synchronization
- `SetDaemon(v bool) bool`: checks `state == stateNew` under `configMu`; returns false if started or terminated
- `ConsumeDaemonForStart()`: reads daemon under `configMu`, establishing happens-before with any `SetDaemon` that completed before the `SetStarted` CAS
- `threadSetDaemon` native: throws `IllegalThreadStateException` when `SetDaemon` returns false
- `threadStart`: uses `ConsumeDaemonForStart()` instead of raw `IsDaemon()`
- No Go data race between `SetDaemon` and the daemon read at start time

## Rework diff (relative to `07d5d16`)

| File | Δ | What |
|---|---|---|
| `rtda/thread.go` | +91/−18 | CAS Terminate, Interrupted waker drain, Sleep re-check, SetDaemon lifecycle + configMu, ConsumeDaemonForStart |
| `native/thread.go` | +10/−6 | ConsumeDaemonForStart in start, throw ITSE in setDaemon |
| `rtda/thread_test.go` | +243/−49 | Double-terminate safe, Interrupted drain + Sleep interaction, SetDaemon lifecycle, concurrent SetDaemon/start race |
| `native/thread_test.go` | +107 | SetDaemon after start → ITSE, after terminate → ITSE, before start succeeds |
| `evidence.md`, `handoff.md` | whitespace fix | Removed trailing whitespace from `505d3ee` evidence |

**Total rework:** 6 files, +421/−40

## Verification gates

| Gate | Command | Result |
|---|---|---|
| Build | `go build ./...` | **Pass** |
| Vet | `go vet ./...` | **Pass** |
| Unit tests | `go test -race -count=1 ./...` | **Pass** — all 8 packages |
| Regression | `bash tests/run.sh` | **Pass** — 10/10 fixtures |
| Whitespace | `git diff --check 3034e05..a0e336c` | **Pass** |

## New tests added (all under `-race`)

### rtda/thread_test.go (+8 test functions, 15 subtests)
- `TestInterruptedDrainsWaker` (2 subtests): sleep-after-clear normal, sleep-interrupted-during throws
- `TestInterruptBoundaryCases` (4 subtests): multi-interrupt-clear, check-clear-check-sleep, pre-interrupted-sleep, concurrent Interrupt+Interrupted drain
- `TestSetDaemonLifecycle` (3 subtests): NEW ok, after start fails, after terminate fails
- `TestConcurrentSetDaemonAndStart` (100 iterations): race-free, outcome reflects valid ordering

### native/thread_test.go (+3 tests)
- `TestThreadSetDaemonAfterStartThrowsITSE`: start then setDaemon → ITSE
- `TestThreadSetDaemonAfterTerminateThrowsITSE`: start, wait terminate, setDaemon → ITSE
- `TestThreadSetDaemonBeforeStartSucceeds`: setDaemon before start succeeds

## Contract gates not yet run

These are full-workstream acceptance gates covering Slices A–E:
- 19-fixture differential matrix
- AOT rejection matrix
- Race stress (`R2_CONCURRENCY_STRESS=100`)
- Evidence isolation check

All remain for Slice D–E completion. Not passed here.
