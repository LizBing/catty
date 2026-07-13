# Slice B Handoff

**Candidate:** `505d3ee`**Branch:** `worktree-r2-thread-monitor-foundation`**Previous slice (A) candidate:** `ec1b398`**Status:** Awaiting Owner review

## What Slice B delivers

A stable Thread execution model per ADR-0028:

1. **Thread identity:** `currentThread()` returns the canonical Java Thread facade object attached to the calling goroutine's execution context. Every rtda.Thread has one java.lang.Thread Object; every started Thread Object has one goroutine carrier.

2. **Lifecycle:** NEW (0) → CAS to RUNNABLE (1) via `SetStarted()` → TERMINATED (2) via `Terminate()`. Start-once enforced by atomic CAS. `isAlive()` returns true only in RUNNABLE.

3. **Goroutine carrier:** `threadStart()` launches a goroutine that pushes a `run()V` frame (virtual dispatch) and calls `rtda.DefaultRunLoop`. On return, `Terminate()` closes the `done` channel and `ThreadTerminated()` decrements the VM count.

4. **Join:** `join()` selects on `target.Done()` and `joiningThread.Waker()`. On waker signal, checks `Interrupted()` and throws `InterruptedException` if set.

5. **Interrupt:** Atomic int32 flag + buffered (cap 1) waker channel. `Interrupt()` stores 1 and does non-blocking send to waker. `isInterrupted()` reads without clear. `interrupted()` (static) swaps to 0.

6. **Sleep:** `Sleep(millis)` checks interrupted first (returns false if set), then `time.After` + waker select. On waker, clears flag and returns false. Caller throws `InterruptedException`.

7. **VM liveness:** `sync.Cond`-based supervisor tracking non-daemon thread count. `WaitForNonDaemonThreads()` blocks until zero. Daemon threads don't count. Main thread is non-daemon by default.

8. **Import cycle:** `rtda.DefaultRunLoop` callback pattern — set by `launch.go` before main, called by `native/thread.go` in the goroutine carrier. Breaks the `native → interpreter → lowering → classloader → native` cycle.

## Key design decisions

- **Thread state as int32 with atomic ops** (not mutex): CAS for start, Load for isAlive, Store for terminate. Simple and race-free.
- **Waker as buffered channel (cap 1):** Non-blocking send from Interrupt, select-based receive from Sleep/Join. Capacity 1 prevents goroutine leaks.
- **Extras pattern:** `obj.SetExtra(t)` / `obj.Extra().(*rtda.Thread)` — same pattern as String → StringValue. No extra allocation.
- **Done channel (unbuffered):** Closed once on Terminate. All joiners see the close simultaneously.
- **throwIllegalThreadState / throwInterruptedException in thread.go:** Direct `newException` + `Throw` — avoids the `throwException` helper that accesses `thread.CurrentFrame()` (which panics on empty stack).

## Files to review

| File | Review focus |
|---|---|
| `rtda/vm.go` | sync.Cond correctness, broadcast at zero |
| `rtda/thread.go:44-226` | Atomic ops, channel lifecycle, Sleep select |
| `native/thread.go` | Goroutine carrier defer, join waker select, exception creation |
| `launch/launch.go:34-73` | VM/DefaultRunLoop init order, main thread setup, liveness wait |
| `interpreter/interpreter.go` | Non-main uncaught exception (was os.Exit, now returns) |

## Next step (Slice C)

Monitors, synchronized methods, wait sets, and interruption per ADR-0029:
- Lazy monitor sidecar per Object
- `monitorenter`/`monitorexit` with recursion and ownership
- `ACC_SYNCHRONIZED` instance/static method entry/exit
- `Object.wait()`, `notify()`, `notifyAll()`
- Interrupt of wait
- `IllegalMonitorStateException`

## Dependencies for Slice C

- Thread identity (done — Slice B)
- Interrupt mechanism (done — Slice B)
- SC heap cells (done — Slice A)
- Canonical Class mirrors (done — Slice A)
