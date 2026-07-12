# ADR-0017: Explicit Java Thread context and reentrant monitors

- **Status:** Accepted
- **Date:** 2026-07-12
- **Refines:** ADR-0010

## Context

ADR-0010 correctly selects goroutines as Catty's scheduling substrate, but its
“Thread = goroutine” shorthand is incomplete. Go exposes no supported goroutine
ID. `sync.Mutex` is not reentrant and does not carry Java monitor ownership,
recursion depth, wait sets, interruption, or timed wait behavior.

Java also specifies lifecycle and publication behavior that a bare goroutine
does not provide: start exactly once, join publication, interrupt status,
daemon process lifetime, currentThread identity, and termination observation.

## Decision

One running Java Thread uses one goroutine, but identity and semantics live in
an explicit Catty runtime context. A Java Thread payload contains at least:

- associated `*rtda.Thread` and Java Thread object;
- lifecycle state and exactly-once start transition;
- completion signal for join and VM liveness;
- atomic interrupt status;
- one-bit park permit;
- daemon flag and termination outcome.

Code executing Java always receives its `*rtda.Thread`; Catty never discovers
identity from a Go goroutine ID.

Each Java object has a lazily allocated reentrant Monitor abstraction with:

- owner `*rtda.Thread` and recursion depth;
- entry contention;
- an explicit wait set;
- individually wakeable timed/interruptible waiters.

The same Monitor implements bytecode `monitorenter`/`monitorexit`, implicit
`ACC_SYNCHRONIZED` method locking, `Object.wait`, `notify`, `notifyAll`, and
`Thread.holdsLock`. Normal and abrupt method exits release implicit monitors.

Class initialization uses a related per-class state machine and lock following
JLS §12.4.2; a boolean `initStarted` is not sufficient under concurrency.

## Consequences

**Positive**

- Retains Go's scheduler and lightweight goroutines while implementing Java
  ownership, lifecycle, and happens-before semantics.
- Avoids unsupported goroutine-ID hacks.
- Provides one monitor contract for interpreter, AOT bridge, and native paths.

**Negative**

- Every live Java thread and contended object carries additional runtime state.
- Correct wait/notify/interrupt races require a state machine more complex than
  `sync.Cond` alone.
- ThreadLocal, daemon shutdown, uncaught termination, and park/unpark become
  explicit runtime responsibilities rather than “free” goroutine behavior.
