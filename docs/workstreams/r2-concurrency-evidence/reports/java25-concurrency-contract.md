# Java 25 concurrency semantic contract for R2 planning

This report extracts the minimum Java-visible obligations needed to plan catty's first
concurrency slice. Java SE/JLS/JVMS 25 govern; Temurin 25.0.3 is only the differential
reference. Links below point to Oracle's Java 25 specifications and API documentation.

## Thread identity and lifecycle

- Every Java thread is associated with a `java.lang.Thread` object; `currentThread()`
  returns the object for the current Java thread. A carrier is not part of this identity
  contract. [JLS 17](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html),
  [Thread API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Thread.html)
- `start()` schedules the Thread's `run` behavior and succeeds at most once; a second
  start throws `IllegalThreadStateException`. A thread is alive after start and before
  termination.
- `join()` detects termination. A call to start happens-before every action in the
  started thread, and all actions in a thread happen-before a successful join return or
  another specified termination detection.
- The VM shutdown sequence begins when all started non-daemon threads have terminated;
  an unstarted non-daemon Thread does not keep the VM alive.

## Monitors and synchronized methods

- Every object has one monitor. At most one Java thread owns it; the owner may enter it
  repeatedly, with one exit required per entry. Other Java threads block until ownership
  becomes available. [JLS §17.1](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html#jls-17.1)
- `monitorenter(null)` and `monitorexit(null)` throw `NullPointerException`.
  `monitorexit` by a non-owner throws `IllegalMonitorStateException`.
  [JVMS monitorenter/monitorexit](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-6.html#jvms-6.5.monitorenter)
- `ACC_SYNCHRONIZED` is implicit VM behavior, not normally a bytecode pair. Instance
  methods lock the receiver; static methods lock the canonical `Class` object for the
  declaring class. Normal return and abrupt escape both release the monitor.
  [JVMS §2.11.10](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-2.html#jvms-2.11.10)
- An unlock on a monitor synchronizes-with every subsequent lock on that monitor in the
  synchronization order.

## Wait sets and notification

- Every object has a wait set associated with the same monitor. `wait`, `notify`, and
  `notifyAll` require ownership or throw `IllegalMonitorStateException`.
  [JLS §17.2](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html#jls-17.2),
  [Object API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Object.html)
- If the current recursion depth is `n`, wait atomically joins the wait set and performs
  `n` unlock actions. Before normal return or `InterruptedException`, it performs `n`
  lock actions and thus restores ownership and recursion depth.
- Wait may resume because of notify, notify-all, interrupt, timeout, or a permitted
  spurious wakeup. Programs must use a condition loop; catty must not promise absence of
  spurious wakeups as Java-visible semantics even if its first mechanism never creates one.
- Notify chooses any one eligible waiter; notify-all makes all eligible waiters compete
  to reacquire. A notified waiter does not run until it later reacquires the monitor.
- Notification and interruption of the same waiter must be ordered. If interrupt wins,
  a later notify cannot be silently consumed by that already-interrupted waiter; Java's
  specified alternatives prevent lost notification.

## Interruption

- Interrupt sets the target Thread's interrupt status and synchronizes-with a later
  specified observation of that interruption.
- `isInterrupted()` observes without clearing; static `Thread.interrupted()` observes and
  clears the current Thread's status.
- A Thread interrupted in Object.wait, Thread.join, or Thread.sleep is enabled to resume,
  clears its interrupt status, and receives `InterruptedException` after any required
  monitor reacquisition. [Thread API](https://docs.oracle.com/en/java/javase/25/docs/api/java.base/java/lang/Thread.html)
- Ordinary blocking while acquiring a Java monitor is not an interruptible operation.
  Interrupt may set status while the Thread remains blocked on entry.

## Cross-thread class initialization

ADR-0025's four states remain correct but need the synchronization procedure in
[JVMS §5.5](https://docs.oracle.com/javase/specs/jvms/se25/html/jvms-5.html#jvms-5.5):

- each runtime class identity has one initialization lock;
- a different owner waits, then retries after success or failure notification;
- the same owner recognizes recursive initialization and returns normally;
- success and every erroneous transition notify all waiters;
- waiting for initialization does not alter interrupt status; and
- eliding the lock after initialization is legal only if the same happens-before edges
  remain.

## Minimum memory-ordering boundary

Java's synchronization order includes monitor lock/unlock, volatile access, Thread start,
termination detection, and interrupt detection. The directly required edges are:

| Release | Acquire / observation |
|---|---|
| monitor unlock | subsequent lock of the same monitor |
| volatile write | subsequent read of the same volatile variable |
| Thread start | first and all later actions in the started Thread |
| Thread final action | successful join/termination detection |
| interrupt action | specified interrupt observation |
| class-init completion/failure lock release | waiter retry/fast path with preserved initialization visibility |

For correctly synchronized programs, Java provides the DRF-SC guarantee: executions
appear sequentially consistent. [JLS §17.4.4–17.4.5](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html#jls-17.4.4)

Final fields add a separate rule. A constructor exit freezes each written final field;
if `this` does not escape before construction completes, another thread that obtains the
object reference must see the initialized final values even when the publication itself
has a data race. Referenced objects/arrays receive the associated minimum visibility.
[JLS §17.5](https://docs.oracle.com/javase/specs/jls/se25/html/jls-17.html#jls-17.5)

## Bounded policy conclusion

The first catty concurrency implementation should select a conservative subset of
Java-allowed executions:

- all Java heap cells used by the supported concurrent surface are accessed through a
  race-free, sequentially consistent Go mechanism;
- volatile fields therefore receive at least their required ordering;
- constructor writes and reference publication also pass through that mechanism, giving
  final fields at least their required visibility; and
- stronger ordering of ordinary fields is an allowed restriction of nondeterministic
  outcomes, not a substitution of Go semantics for Java semantics.

This policy must be an Accepted ADR because it creates a durable heap/AOT/native boundary
and deliberately chooses stronger ordering than Java minimally requires. It can later be
optimized only with evidence that the Java 25 guarantees and catty's Go race freedom are
preserved.
