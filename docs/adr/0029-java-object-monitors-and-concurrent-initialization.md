# ADR-0029: Java object monitors, wait sets, and concurrent class initialization

- **Status:** Accepted
- **Date:** 2026-07-13
- **Governing workstream:** `r2-concurrency-semantics-research`
- **Evidence:** `docs/workstreams/r2-concurrency-evidence/reports/{java25-concurrency-contract,go-mechanism-experiments,runtime-boundary-map}.md`

## Context

Interpreter and IR currently treat `monitorenter` and `monitorexit` as operand pops.
Synthetic Object omits wait/notify, synchronized method flags are ignored, Class mirrors
are not canonical, and ADR-0025's initialization state is intentionally unsynchronized.

A direct Go mutex is not a Java monitor: it is not reentrant and has no Java owner,
recursion depth, wait set, or interrupt ordering. A condition variable or channel alone
also cannot prevent notification loss when notify and interrupt race.

## Proposed decision

Every Java object, including arrays and canonical Class mirrors, SHALL have access to one
lazy monitor sidecar. A monitor is runtime kernel state and is not an ordinary Java field.
It is not copied by `clone()`.

### Monitor state

The sidecar SHALL maintain under one internal state lock:

- owner execution-context identity, or no owner;
- recursion depth;
- entry coordination for blocked contenders; and
- an ordered collection of waiters whose state is `waiting`, `notified`, or
  `interrupted`, each with a private wake signal.

Entry by the owner increments depth; entry by another Java Thread blocks. Exit decrements
depth and releases ownership at zero. Null entry/exit throws `NullPointerException`; exit,
wait, or notify without ownership throws `IllegalMonitorStateException`.

### Synchronized methods and Class mirrors

`ACC_SYNCHRONIZED` SHALL be implemented at the shared invocation/return/unwind boundary.
Instance methods use the receiver monitor; static methods use one canonical Java Class
mirror for the method's declaring runtime class identity. Normal and abrupt completion
both release exactly the implicit entry.

Catty SHALL create at most one Class mirror per runtime class/defining-loader identity.
This canonicality is required for static synchronized semantics and does not authorize
broad reflection.

### Wait, notification, and interruption

- Wait atomically enqueues the caller, saves recursion depth, fully releases the monitor,
  and blocks. Before normal return or `InterruptedException`, it reacquires the monitor
  and restores the saved depth.
- Notify changes one currently waiting entry to notified; notify-all changes all currently
  waiting entries. Selection need not be fair.
- Notify and interrupt transitions are totally ordered under the monitor state lock. If
  interrupt wins, later notify skips that waiter and may select another, preventing lost
  notification.
- Interrupt state belongs to ADR-0028's Java Thread record. Waiting, joining, and sleeping
  observe and clear it when throwing `InterruptedException`; ordinary monitor entry is not
  interruptible.

### Concurrent class initialization

Each runtime Class identity SHALL have one dedicated initialization lock/condition behind
the existing shared ADR-0025 service. It need not reuse the Java Class mirror's monitor.
The service SHALL implement JVMS §5.5 other-owner wait/retry, same-owner recursion,
notify-all on success and every erroneous transition, unchanged interrupt status while
waiting, and release/acquire visibility for initialized state.

### Engine boundary

Interpreter and IR use the same monitor, Thread, mirror, and initialization services.
AOT monitor/concurrency support remains `Not implemented` until ADR-0028's explicit
execution-context ABI is implemented; build-time rejection is required.

## Consequences

- Object representation gains a lazy monitor reference or equivalent side table.
- Native Object and Thread methods become thin facades over shared runtime services.
- Invocation and exception unwinding require paired implicit-monitor cleanup.
- Class loading must become concurrency-safe so a name/loader pair cannot produce two
  runtime Class identities or mirrors.
- Deadlock detection and fairness are not required by Java 25 and remain out of scope.

## Non-scope

Timed-wait precision beyond the accepted implementation contract, monitor deflation,
biased/thin locks, lock elision, deadlock detection, fairness guarantees, virtual-thread
pinning policy, and concurrent AOT execution.
