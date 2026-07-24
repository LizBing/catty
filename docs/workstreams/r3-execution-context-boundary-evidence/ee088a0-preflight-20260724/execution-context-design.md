# ExecutionContext kernel design

Candidate base: `ee088a0ce7b1d0ba6b2f97d59f0c07274464373e`
Date: 2026-07-24
Status: design complete; implementation pending

## Decision

Introduce `rtda.ExecutionContext` as the concrete owner of JVM execution
state. Preserve `rtda.Thread` temporarily as a Go type alias and keep
`NewThread` as a compatibility constructor. This permits generic engine APIs
to migrate by name without duplicating runtime identity or forcing an atomic
repository-wide rewrite.

Do not introduce an execution-context interface. The context is a runtime
identity-bearing object used on interpreter hot paths; an interface would add
indirection without providing a second valid implementation.

## Target ownership

`ExecutionContext` owns:

- frame stack and frame creation;
- initiating loader reference;
- immutable context/monitor/class-initialization owner ID;
- pending Java throwable and throw PC;
- invocation bridge return state;
- an optional binding to Java Thread facade state.

The Java Thread facade sidecar introduced in Slice C owns:

- canonical `java.lang.Thread` object attachment;
- NEW/RUNNABLE/TERMINATED lifecycle;
- interrupt state;
- daemon configuration;
- completion and wake channels;
- primordial-main marker;
- active monitor-wait enrollment needed to coordinate interrupt with a
  context-owned monitor owner ID.

A goroutine is only a carrier chosen by `Thread.start`; it is not stored as an
identity and cannot be used to discover an execution context or Java Thread.

## Kernel API

The first implementation step adds the following concrete API shape:

```go
type ExecutionContext struct {
    // execution-owned fields only after Slice C
}

// Deprecated compatibility name during this workstream.
type Thread = ExecutionContext

func NewExecutionContext(loader Loader) *ExecutionContext

// Deprecated compatibility constructor. During migration it creates the same
// single context object and the lifecycle state required by existing callers.
func NewThread(loader Loader) *ExecutionContext

func (c *ExecutionContext) ID() uint64
func (c *ExecutionContext) Loader() Loader
func (c *ExecutionContext) NewFrame(method *Method) *Frame
```

Existing stack, bridge-return, and throwable methods move unchanged to
`*ExecutionContext`. `EC()` remains as a deprecated forwarding method to
`ID()` until all production consumers have migrated.

`Frame` stores `*ExecutionContext` and gains `Context() *ExecutionContext`.
`Frame.Thread()` remains as a compatibility accessor while native and engine
callers migrate. Both accessors return the same pointer; no wrapper, copy, or
second identity is created.

## Facade binding contract

Slice B does not move lifecycle state. It establishes the naming and
compatibility seam only. Slice C adds a one-to-one optional binding:

```text
ExecutionContext ── optional binding ── JavaThreadState ── facade object
       │                                      │
       └── immutable owner ID                 └── lifecycle/interrupt/wait state
```

Binding is explicit and one-time. `Thread.currentThread()` reads the bound
facade object. A context without a Java facade is valid for internal or AOT
bridge execution, but Java Thread native operations must fail explicitly if
they require a missing binding.

The Java object's `Extra()` representation is not changed in Slice B. Any
change from `*ExecutionContext` to a sidecar handle belongs to Slice C and must
be covered by native Thread identity/lifecycle tests.

## Monitor and class-initialization invariants

- `ExecutionContext.ID()` remains the sole owner token for monitors and class
  initialization; Java object identity and goroutine identity are forbidden as
  substitutes.
- Frame synchronized-monitor enter/exit uses `frame.Context().ID()`.
- Monitor wait remains callable from the execution context. After Slice C its
  interrupt/wait enrollment mechanics delegate to the bound facade sidecar,
  while monitor ownership continues to use the context ID.
- The numeric ID allocation algorithm and zero-as-no-owner reservation do not
  change in this workstream.

## Migration sequence

1. Add `ExecutionContext`, `Thread` alias, `NewExecutionContext`, compatibility
   `NewThread`, `ID`, and compatibility `EC` without changing behavior.
2. Change `Frame` storage and generic `rtda` APIs to execution-context naming;
   keep compatibility accessors.
3. Migrate Interpreter and IR signatures to `*rtda.ExecutionContext`.
4. Introduce the Java Thread facade sidecar and move lifecycle state in Slice C.
5. Migrate native Thread, launcher, and AOT bridge boundaries in Slice D.
6. Remove compatibility names only when an `rg` audit proves there are no
   production consumers and the accepted workstream explicitly authorizes it.

Each step must compile and pass focused tests before the next step begins.

## Rejected alternatives

- A wrapper `ExecutionContext` containing `*Thread`: rejected because it
  creates two apparent execution identities and leaves ownership inverted.
- Immediate repository-wide rename with no alias: rejected because it couples
  semantic separation to mechanical churn and enlarges the regression surface.
- Goroutine-local current-context lookup: rejected because carrier identity is
  not Java-visible identity and goroutines may not define future invocation
  boundaries.
- Moving lifecycle fields in Slice B: rejected because Slice C has separate
  identity, interrupt, join, sleep, and monitor-wait acceptance obligations.

## Slice B acceptance

- `ExecutionContext` is the defined type; `Thread` is only a compatibility
  alias.
- `NewExecutionContext` and `NewThread` return the same runtime shape without a
  wrapper object.
- New focused tests prove context IDs remain unique under concurrent creation,
  frame/context pointer identity is stable, monitor ownership is unchanged,
  and throwable/bridge state remains context-local.
- Existing rtda, interpreter, IR, and native tests pass.
- No Java-visible behavior or AOT concurrency capability is added.
