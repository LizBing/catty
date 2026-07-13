# R2 concurrency current-runtime boundary map

**Code baseline:** `63d5658`
**Differential baseline:** `../baseline-63d5658/run-concurrency-results-v5.txt`
**Scope:** read-only map of production code; research prototypes are described separately.

## Baseline result

All 19 fixed fixtures completed on Temurin 25.0.3. Catty matched none in Interpreter,
IR, or AOT. Interpreter and IR mismatched 19/19; AOT produced 15 `NO-BUILD` results and
four built binaries that failed at run time. The most important single-context defect is
`MonitorNull`: both Interpreter and IR print the body and exit zero instead of throwing
`NullPointerException`.

The failures form five dependency clusters rather than 19 independent missing methods:

1. no stable Java Thread object or lifecycle;
2. no monitor, wait-set, synchronized-method, or interrupt service;
3. no race-free shared heap/classloader/class-initialization boundary;
4. no canonical Class mirror for static synchronized locking; and
5. no explicit execution context across the AOT/runtime/fallback boundary.

## Runtime ownership and mutable-state map

| Boundary | Current code evidence | Concurrency delta |
|---|---|---|
| Java execution context | [`rtda/thread.go:15`](../../../../rtda/thread.go) owns the frame stack, bridge return, pending exception, and `ecID`; the global sequence at line 36 is a plain increment | one `rtda.Thread` per Java Thread; atomic ID allocation; no concurrent use of one frame/exception stack |
| Java Thread facade | [`native/registry.go:93`](../../../../native/registry.go) declares only a no-op constructor; registered methods are absent from the synthetic facade | stable Thread mirror, lifecycle state, target/run dispatch, start/join/isAlive/interrupt/currentThread surface |
| Current Thread identity | [`native/system.go:253`](../../../../native/system.go) allocates a fresh object on every handler call | attach exactly one Java Thread object to each execution context; never infer identity from a goroutine |
| Launch/liveness | [`launch/launch.go:22`](../../../../launch/launch.go) constructs one Thread and returns when its loop drains | supervisor must count started non-daemon Java Threads and delay VM termination until the count reaches zero |
| Object and array storage | [`rtda/object.go:18`](../../../../rtda/object.go) stores raw `[]Slot` and returns it directly from `Fields`; array access returns a mutable slot pointer | introduce non-escaping race-free heap-cell access; update Interpreter, IR, AOT, runtime, and native callers |
| Static storage | [`rtda/class.go:20`](../../../../rtda/class.go) stores raw `[]Slot` and returns it from `StaticVars` | same race-free heap-cell boundary as instance fields and arrays |
| Field metadata | [`rtda/field.go:7`](../../../../rtda/field.go) retains access flags but exposes no volatile/final classification | expose semantic flags and route access through the heap boundary |
| Method metadata | [`rtda/method.go:152`](../../../../rtda/method.go) has no `ACC_SYNCHRONIZED` constant/helper | preserve and classify synchronized methods; invocation and every normal/abrupt exit must pair monitor entry/exit |
| Explicit monitor opcodes | [`interpreter/interpreter.go:864`](../../../../interpreter/interpreter.go) and [`interpreter/ir.go:546`](../../../../interpreter/ir.go) only pop the reference | shared monitor service with null, owner, recursion, blocking, and exit failures |
| Wait/notify facade | synthetic Object at [`native/lang.go:18`](../../../../native/lang.go) declares none of the methods; registry fallbacks are no-ops | facade methods must delegate to the same per-object monitor/wait set as bytecodes and synchronized methods |
| Class mirrors | both [`interpreter/invoke.go:184`](../../../../interpreter/invoke.go) and [`native/native_registry.go:42`](../../../../native/native_registry.go) allocate a new Class object for a runtime class | canonical mirror per runtime class/loader identity; required for static synchronized monitor identity |
| Class loading | [`classloader/classloader.go:93`](../../../../classloader/classloader.go) uses an unprotected cache and performs recursive load/build outside an in-flight protocol | synchronize cache and same-name load; preserve one runtime Class identity without holding a global lock through recursive loading |
| Class initialization | [`rtda/init.go:33`](../../../../rtda/init.go) reads/writes plain state; [`rtda/build.go:163`](../../../../rtda/build.go) explicitly assumes the racing case will not occur | unique init lock per Class, other-owner wait/retry, notify-all on success/failure, uninterrupted wait, and release/acquire publication |
| Native registry | [`native/native_registry.go:13`](../../../../native/native_registry.go) already protects its map with `sync.RWMutex` | no structural change; registrations remain startup-only |
| Native payloads | `Object.Extra()` is an unprotected `any`; String values are immutable, while StringBuilder and some facade payloads are mutable | each mutable supported payload needs its own synchronization or must be rejected from the concurrent surface |
| AOT runtime | [`runtime/runtime.go:20`](../../../../runtime/runtime.go) stores one package-global loader and Thread; every helper and fallback call uses it | loader may be VM-global, but execution context must be an explicit call/bridge parameter |
| AOT signatures | [`transpile/emit.go:669`](../../../../transpile/emit.go) emits only Java parameters | future concurrent AOT needs an internal execution-context parameter; the first bounded slice may explicitly reject concurrency in AOT |
| AOT heap access | [`transpile/emit.go:384`](../../../../transpile/emit.go) emits direct `Fields()[slot]` access | emitter must use the same heap-cell API before any concurrent AOT claim |

## Cross-engine implications

- Interpreter and IR can share one `rtda.Thread`, monitor service, heap service, loader,
  and initialization service. Their dispatch loops are per Thread and must never execute
  the same frame stack concurrently.
- A Java Thread starting in interpreted code may remain interpreted in the first bounded
  slice. This avoids pretending that the current AOT global Thread is carrier-local.
- An AOT method cannot call `runtime.Thread()` to discover its Java identity. Go exposes
  no supported goroutine-ID API, carriers may eventually change, and Java virtual/thread
  identity is explicitly independent of a carrier.
- AOT concurrency is therefore `Not implemented` until generated methods and runtime
  helpers pass an explicit execution context and all emitted heap accesses use the shared
  race-free boundary.

## Required dependency order

1. Canonical execution context and Java Thread mirror.
2. Concurrency-safe loader identity and canonical Class mirrors.
3. Race-free heap cells and audited mutable native payloads.
4. Per-object monitor plus synchronized invocation/exit.
5. Wait sets, interrupt ordering, lifecycle, join, and VM liveness.
6. Cross-thread extension of ADR-0025.
7. Explicit AOT execution-context ABI and emitted heap/monitor access.

This order is an implementation dependency graph, not authorization. The proposed first
slice may combine 1–6 for Interpreter/IR while leaving 7 explicitly `Not implemented`.
