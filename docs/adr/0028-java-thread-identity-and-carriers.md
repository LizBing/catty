# ADR-0028: Java Thread identity, lifecycle, and execution carriers

- **Status:** Accepted
- **Date:** 2026-07-13
- **Governing workstream:** `r2-concurrency-semantics-research`
- **Evidence:** `docs/workstreams/r2-concurrency-evidence/reports/{runtime-boundary-map,java25-concurrency-contract,go-mechanism-experiments}.md`

## Context

ADR-0018 permits goroutines as runtime infrastructure but rejects the withdrawn
`Thread == goroutine` decision. Catty currently has one `rtda.Thread` with a frame stack
and exception signal, a synthetic `java.lang.Thread` with only a constructor, and one
package-global Thread in the AOT runtime. The Java Thread object, lifecycle, daemon state,
interrupt status, completion, and execution carrier are not represented.

The 19-fixture research baseline matches Temurin on no engine. Go experiments show that a
goroutine is a viable initial carrier, but Go exposes no supported goroutine identity and
goroutine completion alone is not a join publication edge.

## Proposed decision

Catty SHALL distinguish three concepts:

1. the stable Java-visible `java.lang.Thread` object;
2. the runtime execution context containing frames, pending exception, interrupt status,
   lifecycle, and its stable Java Thread mirror; and
3. an internal execution carrier.

For the first supported platform-Thread slice, each successful `Thread.start()` SHALL
create one goroutine carrier that closes over exactly one execution context. This is an
initial mechanism, not Java identity and not a commitment for virtual threads.

### Identity and lifecycle

- The initial main execution context has one canonical Java Thread object.
- `currentThread()` returns the canonical object attached to the calling execution
  context. Repeated calls in that Java Thread return the same reference.
- An unstarted Thread has a stable runtime record. `start()` atomically transitions it to
  started and may succeed only once; later calls throw `IllegalThreadStateException`.
- The carrier invokes the Java Thread's selected `run()` method with ordinary virtual
  dispatch. Normal return or uncaught abrupt completion performs one final lifecycle
  action and marks the Thread terminated.
- `isAlive()` observes started-but-not-terminated state. `join()` waits on an explicit
  completion synchronization object so the target's actions happen-before successful
  return.
- VM liveness is supervised independently of the main carrier. Started non-daemon
  Threads are counted before carrier launch and removed at their final action; daemon
  Threads do not keep the VM alive.
- Interrupt status belongs to the Java Thread runtime record, not the goroutine. Its
  interaction with monitor waits, join, and sleep is governed by ADR-0029 if accepted.

### Engine boundary

Interpreter, IR, native calls, and class initialization SHALL receive the execution
context explicitly. No code may discover Java identity from a goroutine ID.

Concurrent AOT execution requires generated methods and every runtime/fallback bridge to
accept an explicit execution context. Until a bounded AOT workstream implements that ABI,
concurrency-triggering AOT programs SHALL be rejected as `Not implemented` at build time;
they may not build and then use the package-global Thread.

### Facade boundary

The first implementation may expand the synthetic `java.lang.Thread` facade with the
bounded methods required by its accepted workstream. The Thread state is a Go-backed
kernel capability. This does not make the whole facade permanently synthetic or decide
broad java.base Thread participation under ADR-0022.

## Consequences

- `rtda.Thread` becomes a Java execution context rather than a carrier synonym.
- Thread IDs must be allocated atomically, but IDs do not replace object identity.
- The launcher needs a VM supervisor and the native Thread facade needs stable mirror
  attachment, lifecycle, completion, daemon, and interrupt state.
- One-goroutine-per-supported-platform-Thread is simple and measurable, but scale and
  virtual-thread claims remain unsupported.
- AOT remains explicit `Not implemented` for the first concurrency slice rather than
  receiving an unsafe goroutine-local or global-context approximation.

## Non-scope

Virtual threads, ThreadGroup, priorities, ThreadLocal, context-class-loader inheritance,
uncaught-exception handler policy, stop/suspend/resume, scheduling fairness, and a
permanent Thread facade/provider split.
