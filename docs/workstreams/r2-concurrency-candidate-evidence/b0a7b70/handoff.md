# Slice B Final Rework Handoff

**Candidate:** `b0a7b70`
**Branch:** `worktree-r2-thread-monitor-foundation`
**Previous rework:** `a0e336c`
**Original Slice B:** `505d3ee`
**Status:** Awaiting Owner review

## What changed

One fix: `IsDaemon()` now holds `configMu` so concurrent `SetDaemon()` on an
unstarted thread is Go-race-free. This is the last known data race in the
Slice B daemon subsystem.

All three configMu users now uniformly synchronize on the same mutex:

| Method | configMu role |
|---|---|
| `SetDaemon(bool) bool` | Lock → check state==NEW → write daemon → Unlock |
| `IsDaemon() bool` | Lock → read daemon → Unlock |
| `ConsumeDaemonForStart() bool` | Lock → read daemon → Unlock |

## Integrity of prior fixes

All three blocker fixes from `a0e336c` remain intact:

1. **Terminate exactly-once** — CAS from RUNNABLE to TERMINATED; only the
   winner closes `done`. Double Terminate is safe, no panic.
2. **Interrupted drains waker** — `Interrupted()` drains after Swap;
   `Sleep()` re-checks `Interrupted()` on waker signal.
3. **setDaemon lifecycle** — `SetDaemon` checks `state==NEW` under
   `configMu`; `ConsumeDaemonForStart` reads under same mutex.

## Next step (unchanged)

Slice C: monitors, synchronized methods, wait sets, and interruption per ADR-0029.
