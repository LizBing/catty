# Slice B Rework Handoff

**Candidate:** `a0e336c`
**Branch:** `worktree-r2-thread-monitor-foundation`
**Original Slice B:** `505d3ee`
**Status:** Awaiting Owner review

## Rework summary

Three acceptance blockers on `505d3ee` resolved:

### 1. Terminate exactly-once
`Terminate()` was unconditionally closing `done`, causing panic on double-call.
Now uses `CAS(stateRunnable, stateTerminated)` — only the first transition
closes `done`. Concurrent/repeated calls return immediately.

### 2. Interrupted waker drain
`Interrupted()` cleared the flag but left a stale waker signal, causing
subsequent `Sleep()` to incorrectly throw `InterruptedException`.
Now `Interrupted()` drains the waker after the atomic Swap, and `Sleep()`
re-checks `Interrupted()` on waker instead of unconditionally clearing.

### 3. setDaemon lifecycle synchronization
`SetDaemon` had no lifecycle check (could be called after start) and no
synchronization with the daemon read at start time (potential data race).
Now uses `configMu` to serialize `SetDaemon` with `ConsumeDaemonForStart()`,
checks `state == stateNew`, and returns false (caller throws ITSE) otherwise.

## Key design changes

- **Terminate CAS:** `atomic.CAS(stateRunnable, stateTerminated)` before `close(done)`.
  If CAS fails (already terminated or never started), returns immediately.
- **Interrupted drain:** `select { case <-t.waker: default: }` after Swap.
  A concurrent `Interrupt()` that fires after the Swap re-sets the flag, so
  no real interrupt is lost.
- **Sleep re-check:** On waker signal, calls `Interrupted()` instead of
  unconditionally clearing. If `Interrupted()` returns false (stale wake),
  `Sleep` returns true (completed normally).
- **configMu:** `sync.Mutex` serializing SetDaemon write and
  ConsumeDaemonForStart read. The SetStarted CAS provides the state transition;
  configMu provides the happens-before for the daemon value.

## Next step (unchanged)

Slice C: monitors, synchronized methods, wait sets, and interruption per ADR-0029.

## Dependencies verified

- Thread identity, lifecycle, interrupt, join, sleep (Slice B — reworked)
- SC heap cells, canonical Class mirrors (Slice A — unchanged)
