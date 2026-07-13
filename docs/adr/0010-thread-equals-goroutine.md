# ADR-0010: Thread = goroutine (virtual threads from day one)

- **Status:** Withdrawn
- **Date:** 2026-07-12
- **Withdrawn:** 2026-07-13 — it conflates a possible goroutine execution
  carrier with Java Thread identity, lifecycle, monitor, wait-set, interrupt,
  daemon, join, and VM-liveness semantics. R2 research must define the Java
  observable contract before selecting Go mechanisms.

## Context

Java 21's Virtual Threads (Project Loom) took years to implement M:N threading on
top of OS threads. Go has had goroutines since day one — M:N scheduling,
work-stealing, ~2KB initial stack, ~100ns context switch. catty doesn't need to
"implement" virtual threads: **every Java `Thread` IS a goroutine.**

## Decision

`java.lang.Thread.start()` directly calls `go func(){}`. Each Java `Thread` = one
goroutine.

- `Thread`'s JVM stack = goroutine's growable stack.
- `synchronized` → `sync.Mutex` (lazily allocated per object).
- `Object.wait`/`notify` → `sync.Cond`.
- `Thread.sleep` → `time.Sleep`.
- `Thread.interrupt` → `context.Cancel` or channel.

No OS thread pool, no Loom, no `ForkJoinPool` mapping.

## Consequences

**Positive**
- Virtual threads advantages (million-scale concurrency, low memory) from day
  one.
- No Loom implementation cost.
- Go's GMP scheduler provides work-stealing automatically.
- Stack grows on demand (Go runtime handles it, unlike fixed 1MB OS thread
  stacks).

**Negative**
- `synchronized` holding blocks the goroutine (Go has no coroutine-level
  preemptive mutex). Long critical sections block the scheduler — mitigation:
  recommend `java.util.concurrent` lock-free structures.
- Thread ID / `ThreadLocal` need goroutine ID mapping (Go doesn't expose goroutine
  ID; use `sync.Map` side-table).
