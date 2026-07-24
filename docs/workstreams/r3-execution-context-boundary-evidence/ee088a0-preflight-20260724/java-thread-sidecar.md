# Java Thread sidecar implementation evidence

Date: 2026-07-25

Slice C moved Java Thread facade/runtime state behind `rtda.JavaThreadState`
while preserving the temporary `rtda.Thread` compatibility alias for execution
contexts.

## Runtime ownership

- `rtda.ExecutionContext` owns execution state: frame stack, loader, execution
  context ID, bridge return, pending throwable, and throw PC.
- `rtda.JavaThreadState` owns Java Thread facade/runtime state: canonical
  `java.lang.Thread` facade object, lifecycle, interrupt flag, daemon flag,
  completion channel, interrupt waker, main marker, and monitor-wait enrollment.
- `rtda.Thread` remains a temporary alias for `ExecutionContext`.
- `ExecutionContext` compatibility methods such as `SetStarted`, `Interrupt`,
  `Done`, `MonitorWait`, and `Sleep` delegate to the sidecar.

## Facade object boundary

- `Thread.<init>` creates an `ExecutionContext`, obtains its
  `JavaThreadState`, attaches the facade object to the sidecar, and stores the
  sidecar in `Object.Extra()`.
- The launcher creates the main execution context, attaches the canonical main
  Thread facade to the sidecar, and stores `thread.JavaThreadState()` in the
  facade object's `Extra()`.
- Thread native methods extract `*rtda.JavaThreadState` from Thread facade
  objects. Carrier creation explicitly retrieves the owning execution context
  via `JavaThreadState.Context()`.

## Preserved behavior

- `Thread.currentThread()` still returns the canonical facade object attached
  to the calling execution context.
- `Thread.start`, `isAlive`, untimed `join`, interrupt, static
  `interrupted`, `sleep`, `setDaemon`, `isDaemon`, VM liveness, and
  monitor-wait interruption continue to use the same synchronization semantics.
- Goroutine identity remains unobservable and is not used for Java identity.

## Checks

- `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda ./native ./interpreter ./launch ./runtime` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...` — Pass
- `git diff --check` — Pass
