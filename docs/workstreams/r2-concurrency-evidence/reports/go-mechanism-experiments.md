# Go mechanism experiments for R2 concurrency

**Status:** research evidence only; no production code is changed.
**Prototype package:** `../prototypes/`
**Toolchain:** Go 1.26.5 darwin/arm64

## Commands and results

| Command | Exit | Result |
|---|---:|---|
| `go test -race -count=50 ./docs/workstreams/r2-concurrency-evidence/prototypes` | 0 | Pass, `1.670s` |
| `go test -race -count=100 ./docs/workstreams/r2-concurrency-evidence/prototypes` | 0 | Pass, `1.956s` |
| `go test -count=1000 ./docs/workstreams/r2-concurrency-evidence/prototypes` | 0 | Pass, `4.412s` |

The first sandboxed race invocation failed before test execution because the managed
sandbox denied Go's default build-cache path. The same exact test was rerun with the
approved external cache access and passed; this environmental failure is not counted as a
semantic test result.

## Go guarantees used as mechanisms

The [Go memory model](https://go.dev/ref/mem) supplies usable implementation edges:

- a `go` statement is synchronized-before the new goroutine begins;
- Mutex unlock is synchronized-before a later successful lock;
- channel send/receive and close/receive provide documented edges;
- Go atomic operations are sequentially consistent; and
- goroutine exit alone has no publication guarantee, so Java join needs an explicit
  completion synchronization operation.

These facts justify mechanisms only. Java-visible behavior remains the contract in the
Java 25 report.

## Candidate comparison

| Concern | Candidate | Finding |
|---|---|---|
| Java Thread carrier | one goroutine per started platform Thread | viable first mechanism; Java Thread object/state remains distinct and is captured explicitly by the goroutine closure |
| Current-thread lookup | goroutine ID / global Thread | reject: Go has no supported goroutine ID and the AOT global aliases all Java Threads |
| Current-thread lookup | explicit `*rtda.Thread`/execution-context argument | select long-term; works across interpreter, native, and AOT/fallback boundaries |
| Java monitor | direct `sync.Mutex` | reject: not reentrant and has no Java owner identity or recursion depth |
| Java monitor | `sync.Cond` alone | reject: no ownership, wait-set membership, interrupt state, or full recursion-depth release/restore |
| Java monitor | per-object sidecar with internal mutex/condition and explicit Java owner ID | select; prototype passes exclusion, reentrancy, ownership, wait restore, and interrupt/notify ordering tests |
| Waiter wakeup | shared broadcast channel only | reject: cannot atomically order notify versus interrupt or reassign a notification after interrupt wins |
| Waiter wakeup | per-waiter state + signal, state transitions under monitor lock | select; notification and interrupt have a total order and notification skips non-waiting entries |
| Class initialization | `sync.Once` | reject: cannot model same-owner recursion, erroneous state, distinct first/later exceptions, or notify/retry procedure |
| Class initialization | four-state cell + owner + condition | select; prototype runs initializer once, blocks other owner, and publishes completion |
| VM liveness | wait for main goroutine only | reject: violates non-daemon lifetime |
| VM liveness | supervisor count updated before carrier start and at final action | select; prototype distinguishes daemon/non-daemon and joins through an explicit completion channel |
| Shared heap | raw `Slot` fields/slices | reject: Java sharing becomes catty Go data races and category-2/reference copies are not a stable atomic boundary |
| Shared heap | one coarse lock per object/class | viable but forces all access through a lock, complicates array/clone/native bulk operations, and couples monitor contention to storage unless locks are separate |
| Shared heap | dedicated SC heap cells (`uint64` bits or atomic object reference) | select for first concurrency slice; race-free, gives volatile/final at least required ordering, and respects ADR-0020's separate heap domain |

## Prototype findings

The prototype deliberately keys ownership by a supplied Java execution-context ID. It
shows that a monitor can:

- block a different Java owner until the final recursive exit;
- release all recursion depth for wait and restore the same depth before return;
- order waiter notification and interruption under one state lock; and
- avoid losing a notification when interruption removes one of two waiters first.

The lifecycle prototype registers a non-daemon before launching its carrier, closes a
completion channel as the Thread's final action, rejects a second start, and allows join
and the VM supervisor to acquire the worker's writes. A daemon carrier is not counted for
VM liveness.

The class-init prototype validates the mechanism shape required by ADR-0025/JVMS §5.5.
It is not a replacement implementation: production still needs Java exception objects,
predecessor ordering, recursive execution through the interpreter, and the exact
first-failure/later-`NoClassDefFoundError` behavior.

## Rejected shortcuts and limitations

- Passing `go test -race` proves the prototypes contain no detected Go race under the
  exercised schedules; it does not prove the Java Memory Model.
- `FinalFieldPublication` uses a volatile handoff and therefore does not isolate final
  freeze semantics. The JLS contract, not a probabilistic litmus, is the acceptance source.
- The prototype does not implement timed wait, sleep, scheduling priority, ThreadGroup,
  virtual threads, ThreadLocal, uncaught-exception handlers, or broad java.base Thread.
- Per-waiter channels are notification mechanisms, not Java wait-set semantics by
  themselves; all selection/state changes must remain under the monitor state lock.
- SC heap cells intentionally over-order ordinary Java field/array accesses. This is a
  bounded correctness-first choice whose performance cost must be measured only after
  the R2 semantic milestone.

## Recommendation

Use one goroutine as the initial carrier for each started supported platform Thread, but
keep stable Java Thread state and identity in the runtime. Use a per-object lazy monitor
sidecar, a dedicated per-Class initialization condition, explicit execution-context
passing, and a distinct SC heap-cell representation. Keep concurrent AOT execution
`Not implemented` in the first production slice; adding the AOT execution-context ABI and
emitted heap accesses is a subsequent bounded slice.
