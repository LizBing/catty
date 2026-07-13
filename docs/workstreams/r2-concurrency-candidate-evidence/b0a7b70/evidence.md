# Slice B Final Rework Evidence — `b0a7b70`

**Candidate:** `b0a7b70` (branch: `worktree-r2-thread-monitor-foundation`)
**Previous rework:** `a0e336c`
**Original Slice B:** `505d3ee`
**Date:** 2026-07-14
**Status:** Complete, Owner accepted 2026-07-14

## Fix: IsDaemon race-free with SetDaemon

`IsDaemon()` now holds `configMu` for the read, matching the write path in
`SetDaemon()`. Previously `IsDaemon()` read `t.daemon` without any
synchronization, creating a Go data race when called concurrently with
`SetDaemon()` on an unstarted thread.

### Change

| Before | After |
|---|---|
| `func (t *Thread) IsDaemon() bool { return t.daemon }` | Acquires `configMu`, reads under lock, releases |

### Why configMu and not atomic

`configMu` already serializes `SetDaemon` (write) with `ConsumeDaemonForStart`
(read at start time). Adding `IsDaemon` to the same mutex keeps one
synchronization boundary for all daemon state access. An `atomic.Bool` would
also work but adds a second mechanism for the same concern.

### No regression on accepted fixes

- `Terminate()`: still CAS-guarded, exactly-once (unchanged)
- `Interrupted()`: still drains waker after Swap (unchanged)
- `Sleep()`: still re-checks `Interrupted()` on waker (unchanged)
- `SetDaemon` lifecycle: still checks state==NEW under configMu (unchanged)
- `ConsumeDaemonForStart`: still holds configMu for the start-time read (unchanged)

## Diff (relative to `a0e336c`)

| File | Δ | What |
|---|---|---|
| `rtda/thread.go` | +8/−3 | IsDaemon holds configMu |
| `rtda/thread_test.go` | +68 | TestConcurrentSetDaemonAndIsDaemon (3 subtests) |

**Total:** 2 files, +76/−3

## Verification gates

| Gate | Command | Result |
|---|---|---|
| Build | `go build ./...` | **Pass** |
| Vet | `go vet ./...` | **Pass** |
| Unit tests | `go test ./...` | **Pass** |
| Race tests | `go test -race ./...` | **Pass** — all 8 packages |
| Regression | `bash tests/run.sh` | **Pass** — 10/10 fixtures |
| Whitespace | `git diff --check 3034e05..b0a7b70` | **Pass** |

## New test: TestConcurrentSetDaemonAndIsDaemon (3 subtests, all under `-race`)

1. **concurrent read write on NEW thread is race-free** — 200 goroutines,
   100 writers calling `SetDaemon(bool)`, 100 readers calling `IsDaemon()`.
   Assertion: zero Go race detector warnings.

2. **daemon frozen after start** — `SetDaemon(true)`, `SetStarted()`,
   verify `IsDaemon()` still returns true, verify `SetDaemon(false)` fails,
   verify value unchanged.

3. **daemon frozen after terminate** — same pattern with `Terminate()`.

## Contract gates not yet run

19-fixture matrix, AOT rejection matrix, race stress — all belong to Slices D–E.
