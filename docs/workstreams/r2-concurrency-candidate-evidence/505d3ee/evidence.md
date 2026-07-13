# Slice B Evidence — `505d3ee`

**Candidate:** `505d3ee` (branch: `worktree-r2-thread-monitor-foundation`)  
**Base:** `3034e05` (Slice A acceptance + governance)  
**Date:** 2026-07-14  
**Status:** Complete, pending Owner review

## Scope delivered

ADR-0028 three-way split: Java Thread facade → runtime execution context → goroutine carrier.

| Capability | Status |
|---|---|
| Thread lifecycle (NEW → RUNNABLE → TERMINATED) | Implemented |
| Start-once atomic CAS | Implemented |
| currentThread() canonical identity | Implemented |
| Goroutine carrier per started Thread | Implemented |
| Join (interruptible) | Implemented |
| Interrupt flag (set/read/clear) | Implemented |
| Sleep (interruptible) | Implemented |
| VM daemon liveness supervisor | Implemented |
| Daemon/non-daemon flag | Implemented |
| Main thread special handling | Implemented |
| Uncaught exception: non-main returns, main os.Exit | Implemented |
| Import cycle resolution (DefaultRunLoop callback) | Implemented |

## Files changed

| File | Δ | Role |
|---|---|---|
| `rtda/vm.go` | +66 NEW | VM supervisor (sync.Cond, non-daemon count) |
| `rtda/thread.go` | +154/−8 | Lifecycle, interrupt, daemon, sleep, done, waker, DefaultRunLoop |
| `native/thread.go` | +236 NEW | 15 native Thread methods, exception classes |
| `native/registry.go` | +25/−1 | buildThread expanded from stub to 15 methods |
| `native/system.go` | −9 | Removed dead threadCurrentThread |
| `launch/launch.go` | +32 | VM init, main thread facade, liveness wait |
| `interpreter/interpreter.go` | +8/−3 | Non-main uncaught exception returns |
| `rtda/thread_test.go` | +292 NEW | 32 test cases (lifecycle, interrupt, sleep, daemon, done, waker, concurrent ecID) |
| `rtda/vm_test.go` | +146 NEW | 5 test cases (count tracking, daemon exclusion, concurrent) |
| `native/thread_test.go` | +519 NEW | 14 test cases (init, currentThread, isAlive, start+join, start-twice, interrupt, sleep) |

**Total:** 10 files, +1464/−23

## Verification gates

| Gate | Result |
|---|---|
| `go build ./...` | **Pass** |
| `go vet ./...` | **Pass** |
| `go test -race ./...` | **Pass** — all packages |
| `rtda` tests | **Pass** — 107 cases (32 new for Slice B) |
| `native` tests | **Pass** — 14 cases (all new for Slice B) |
| `git diff --check` | **Pass** — no whitespace issues |

## Test coverage

### rtda/thread_test.go (32 cases)
- Lifecycle: NEW not alive, SetStarted transitions, SetStarted once/second fails, terminated not alive, Done channel close, double-terminate safety
- Interrupt: initial state, Interrupt sets flag, IsInterrupted no-clear, Interrupted read-and-clear
- Waker: Interrupt signals, second Interrupt signals after drain
- Sleep: normal completion, interrupted-before, interrupted-during, zero sleep
- Daemon: default/SetDaemon/IsDaemon
- Done: blocks until Terminate
- JavaThread: nil initially, set/get round-trip
- Concurrent: 100 goroutine ecID uniqueness
- Main flag: default/SetMain/IsMain

### rtda/vm_test.go (5 cases)
- Non-daemon count: increment/decrement, WaitForNonDaemon blocks/unblocks
- Daemon exclusion: daemon threads don't affect count
- Zero-count immediate unblock
- Concurrent start/terminate (100 goroutines)
- SetVM/GetVM round-trip

### native/thread_test.go (14 cases)
- init attachment: Extra ↔ JavaThread bidirectional
- currentThread identity: returns calling thread's canonical facade
- isAlive: NEW/RUNNABLE/TERMINATED states
- start+join happy path: goroutine carrier, join unblocks, VM unblocks
- start twice → ITSE
- interrupt flag: isInterrupted no-clear, interrupted read+clear, static vs instance
- sleep: normal, interrupted-before, interrupted-during
- setDaemon/isDaemon
- VM liveness: non-daemon keeps alive, daemon doesn't
- Daemon thread start liveness

## Non-scope (not touched)

- Monitors, synchronized methods, wait sets (Slice C)
- Concurrent class initialization (Slice D)
- AOT execution context (Slice E)
- Timed wait/join, ThreadGroup, ThreadLocal, priorities, virtual threads
