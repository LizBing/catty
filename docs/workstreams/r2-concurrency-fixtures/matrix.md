# R2 concurrency research fixture matrix

This is the fixed 19-fixture baseline set for the accepted
`r2-concurrency-semantics-research` contract. Every fixture compiles with
`javac --release 25`, produces deterministic Temurin 25 output, and is run with bounded
wall-clock time in Interpreter, IR, and AOT modes. A catty mismatch or `NO-BUILD` is a
research observation; a missing row, reference failure, or timeout of the reference is a
harness failure.

| Fixture | Primary obligation | Determinism mechanism |
|---|---|---|
| `CurrentThreadIdentity` | stable Java Thread identity and liveness | same-thread identity comparison |
| `ThreadStartJoin` | lifecycle, start, termination, join visibility | join before reads |
| `ThreadStartTwice` | start-at-most-once failure | first run joined before second start |
| `NonDaemonLiveness` | VM waits for started non-daemon threads | volatile handoff fixes output order |
| `DaemonLiveness` | daemon does not prevent VM shutdown | long sleeper must not print before main exits |
| `SynchronizedBlocks` | exclusion and reentrancy | joins plus exact final counter |
| `SynchronizedMethods` | instance/static implicit monitors | `holdsLock` within method bodies |
| `MonitorNull` | null `monitorenter` failure | single execution context |
| `MonitorOwnership` | notify ownership failure | single execution context |
| `WaitNotify` | wait releases/reacquires; notification | condition loops and join |
| `NotifyAll` | wait-set broadcast and resumption | two waiters, exact resumed count |
| `InterruptStatus` | interrupt observe/clear behavior | self-interruption |
| `InterruptWait` | interrupt removes waiter and clears status on throw | readiness handshake and join |
| `InterruptSleep` | sleep interruption and clear-on-throw | volatile readiness and join |
| `InterruptJoin` | join interruption and clear-on-throw | volatile readiness; both threads terminated |
| `VolatilePublication` | volatile release/acquire visibility | volatile handoff and join |
| `FinalFieldPublication` | final-field read through safe publication | volatile handoff; does not alone prove freeze semantics |
| `CrossThreadClassInitialization` | single initializer, contention visibility | two readers and joins |
| `ProducerConsumer` | R2 milestone: monitor-based one-slot handoff | condition loops, sentinel, and joins |

Final-field freeze semantics and improperly synchronized executions require specification
analysis and supplemental bounded experiments; the `FinalFieldPublication` differential
fixture intentionally proves only the safely-published path.
