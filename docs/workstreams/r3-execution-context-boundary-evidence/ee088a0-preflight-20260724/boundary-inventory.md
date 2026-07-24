# Execution-context boundary inventory

Preflight candidate: `ee088a0ce7b1d0ba6b2f97d59f0c07274464373e`
Date: 2026-07-24

This is an inventory only. It does not claim that the boundary has been
implemented.

## Current `rtda.Thread` ownership

| Current state | Intended owner | Evidence |
|---|---|---|
| frame stack, current frame, frame count | `ExecutionContext` | `rtda/thread.go:40-49`, `PushFrame`/`PopFrame` |
| loader and execution-context ID | `ExecutionContext` | `rtda/thread.go:41-47`, `Loader`, `EC` |
| bridge return slot | `ExecutionContext` bridge state | `rtda/thread.go:49-53`, `SetBridgeReturn` |
| pending throwable and throw PC | `ExecutionContext` | `rtda/thread.go:54-57`, `Throw`/`ClearException` |
| Java Thread facade attachment | Java Thread facade sidecar | `rtda/thread.go:60-62`, `native/thread.go:20-24` |
| lifecycle, interrupt, daemon, completion/waker channels | Java Thread facade/lifecycle service | `rtda/thread.go:63-83`, `native/thread.go:31-64` |
| monitor waiter enrollment | execution-context wait service, coordinated with facade interrupt | `rtda/thread.go:85-94`, `MonitorWait` |

## Consumer map

- Interpreter and IR pass `*rtda.Thread` through execution, invocation,
  exception, class-initialization, monitor, and native paths.
- `native/thread.go` is the Java Thread facade/lifecycle boundary: it creates
  the runtime object, returns `currentThread`, starts the goroutine carrier, and
  manages lifecycle-visible operations.
- `launch/launch.go` creates the primordial context and attaches its canonical
  Java Thread object.
- `runtime/runtime.go` owns one package-global loader and thread for the AOT
  bridge; this remains a separate containment concern and is not resolved by
  this inventory.
- Tests in `rtda/thread_test.go`, `native/thread_test.go`, and
  `native/object_test.go` directly construct and inspect `rtda.Thread`; these
  are migration-sensitive fixtures.

## First migration seam

The lowest-risk first seam is to introduce execution-context naming and APIs
around frame/loader/exception/EC ownership while keeping a compatibility
adapter for existing interpreter signatures. Java Thread lifecycle and facade
operations should migrate separately, with `Thread.currentThread()` continuing
to return the attached Java object.
