# ADR-0025: Java 25 class and interface initialization state machine

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** None (codifies the research decision left open by [ADR-0021](./0021-class-initialization-boundary.md))
- **Governing workstream:** `r2-runtime-semantics-research` (Slice B)
- **Differential evidence:** `docs/workstreams/r2-evidence/run-r2-results.txt`,
  `docs/workstreams/r2-evidence/reports/r2-init-deltas.md`

## Context

ADR-0021 separated loading, linking, and initialization and deferred the concrete Java 25
state machine. R1 has only `initStarted bool`, skips interface initialization, has no
erroneous state, and mishandles static fields referenced through a non-declaring class.

This ADR fixes the semantic contract for one Java execution context. Cross-thread waiting,
Java Memory Model visibility, reflection, method handles, and VM startup are not silently
claimed by this slice; the state representation must leave a synchronization boundary for
later work under ADR-0018.

## Decision

catty SHALL use one shared class/interface initialization service across Interpreter, IR,
and AOT. Each runtime class or interface identity (including its defining loader identity)
has exactly one of four states: **not-initialized**, **initializing**, **initialized**, or
**erroneous**.

```text
not-initialized --(initialization request)--> initializing --(success)--> initialized
                                                |
                                                +--(abrupt completion)--> erroneous
```

### Initialization requests in the R2 slice

The bytecode use points in scope are exactly:

1. `new` for the class to be instantiated;
2. `getstatic` for the class or interface that actually declares the resolved field,
   unless the field is a constant variable;
3. `putstatic` for the class or interface that actually declares the resolved field; and
4. `invokestatic` for the class or interface that actually declares the resolved method.

`invokevirtual`, ordinary `invokespecial`, and `invokeinterface` are not initialization
triggers. In particular, invoking a default method does not itself initialize its declaring
interface. Instead, when a **class** is initialized, its superclass and the recursively
enumerated superinterfaces that declare at least one non-abstract, non-static method are
initialized first, in the JVMS §5.5 order. Initializing an **interface** does not initialize
its superinterfaces.

An `assert` statement is not a separate trigger. The VM determines a class's desired
assertion status while performing that class's initialization.

JVMS §5.5 also defines requests originating from specific method handles, reflection, and
VM startup. Those mechanisms are outside this bounded slice. They SHALL eventually call
the same shared service; no separate initialization semantics may be introduced.

### State transitions and ordering

- A request for an **initialized** class/interface completes normally without executing
  `<clinit>` again.
- A recursive request from the execution context already initializing that same identity
  completes normally without re-entering or re-running `<clinit>`.
- To initialize a class, initialize its superclass first, then the required default-bearing
  superinterfaces in JVMS order, then invoke the class's `<clinit>` if present.
- If superclass or required-superinterface initialization completes abruptly, mark the
  class erroneous and complete its request abruptly with the same reason.
- To initialize an interface, invoke only that interface's `<clinit>` if present; do not
  initialize its superinterfaces merely because they are extended.
- A class or interface with no `<clinit>` still reaches **initialized** after its required
  predecessors complete.
- Preparation supplies `ConstantValue` fields before this procedure. Reading a constant
  variable does not request initialization.

### Failure semantics

- If `<clinit>` completes by throwing an instance of `java.lang.Error`, propagate that
  object and mark the class/interface erroneous.
- Otherwise attempt to wrap the thrown object in `ExceptionInInitializerError`, propagate
  the wrapper, and mark the class/interface erroneous, following JVMS §5.5.
- A later request for an erroneous identity throws `NoClassDefFoundError`. This ADR does
  not require a particular observable cause chain beyond Java 25 behavior established by
  fixtures; it does not prescribe reusing the first `ExceptionInInitializerError` object.

### Execution-context boundary

R2 records which Java execution context owns an **initializing** state so recursive requests
can be recognized. This ADR does not equate a Java thread with a goroutine and does not add
cross-context waiting or visibility guarantees. A future concurrency ADR may place locking,
waiting, and happens-before behavior behind the shared initialization service without
changing the Java-visible state transitions above.

### Engine obligations

- Interpreter and IR SHALL request initialization at the four in-scope bytecode use points.
- AOT SHALL preserve the same requests in emitted code or shared runtime bridges. In
  particular, its `invokestatic` path needs the guard; virtual, special, and interface
  invocation paths do not acquire guards merely because they invoke methods.
- Field and method resolution SHALL pass the actual declarer to the shared service.
- Unsupported AOT patterns remain explicit `Fallback` or `Not implemented` under ADR-0016;
  they may not approximate initialization silently.

## Non-scope

- Cross-Java-thread initialization contention, blocking, deadlock behavior, and JMM
  happens-before guarantees.
- Reflection, method-handle initialization requests, VM-startup initialization, broad
  `java.base` compatibility, `invokedynamic`, and JIT compilation.
- A prescribed Go lock type, goroutine mapping, or engine-specific state machine.

## Consequences

- R1's boolean is insufficient; implementation needs four states plus initializing-owner
  identity at the shared runtime boundary.
- `<clinit>` executes at most once for each runtime class/interface identity, not once per
  source name, VM globally, goroutine, or Java thread.
- The declarer-owner, constant-field, recursive-request, predecessor-ordering, and
  erroneous-state rules become acceptance evidence rather than inferred implementation
  behavior.
- `InterfaceDefaultInit` tests initialization of a class's default-bearing superinterface
  during class initialization; it is not evidence for an `invokeinterface` trigger.
