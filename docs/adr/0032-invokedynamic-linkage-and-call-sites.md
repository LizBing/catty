# ADR-0032: InvokeDynamic linkage kernel and typed call-site state

- **Status:** Accepted
- **Date:** 2026-07-17
- **Revised:** 2026-07-18 for ADR-0034 profile boundaries
- **Governing workstream:** `r3-reflection-dynamic-research`
- **Evidence:** `docs/workstreams/r3-reflection-dynamic-reports/{r3-java25-semantic-contract,r3-classfile-metadata-gap,r3-runtime-boundary-map,r3-engine-capability-matrix}.md`

## Context

Catty parses constant-pool tags for MethodType, MethodHandle, and InvokeDynamic
but hides their operands from consumers and discards BootstrapMethods.
Interpreter explicitly rejects opcode `0xba`, IR lowering rejects it, and AOT
no-builds all fixed rows with generic diagnostics. There is no runtime method
type, method handle, call site, per-instruction resolution state, bootstrap
invocation, or failure publication model.

JVMS 25 resolution binds a dynamically computed call site to one
`invokedynamic` instruction, not merely to a shared constant-pool entry.
Successful and failed resolution are reused by later executions of that same
instruction; a failed instruction does not execute its bootstrap again.

ADR-0034 requires the JVMS linkage kernel to be distinguished from public Java
SE `java.lang.invoke` APIs and OpenJDK bootstrap factories. Catty must report
the exact bounded linkage capability it supports rather than treating
StringConcatFactory or LambdaMetafactory as the definition of opcode support.

## Decision

### Capability levels and profile ownership

Catty SHALL publish separate capabilities for:

1. **InvokeDynamic linkage kernel** — symbolic metadata, typed method/handle
   identities, per-instruction linkage state, bootstrap invocation, target
   validation, and invocation semantics for the declared bootstrap surface;
2. **Catty Runtime bootstrap providers** — bounded bootstraps supplied by the
   Catty Runtime Profile; and
3. **Java SE dynamic compatibility** — declared `java.lang.invoke` APIs,
   arbitrary supported Java bootstraps, StringConcatFactory,
   LambdaMetafactory, and related Java SE behavior.

These levels share the same kernel and object world but have independent
profile and engine states. Catty SHALL NOT claim general `invokedynamic`
support when only named bootstrap providers or fixture shapes are accepted.
Kernel readiness by itself is not opcode support: a site is `Supported` only
when its declared bootstrap protocol, linkage result, invocation, and
JVMS-required failure behavior are all implemented.

### Symbolic and resolved forms

Classfile SHALL expose validated BootstrapMethods, MethodType, MethodHandle,
and InvokeDynamic operands as immutable symbolic metadata. It SHALL
structurally parse ConstantDynamic so metadata is not lost, but execution of
ConstantDynamic remains a separate capability.

The runtime SHALL represent canonical logical method types and bounded typed
method handles for the invocation and field kinds required by accepted
consumers. A method handle is a typed executable reference, not a Go function
pointer or pre-filled interpreter frame. Its invocation consumes and produces
ADR-0031 typed dynamic values if ADR-0031 is accepted.

The call-site kernel SHALL represent a typed target and its update semantics.
The first bounded capability may support only immutable/constant target state.
Mutable and volatile target-update semantics require an explicit later
capability decision. Profile-visible MethodType, MethodHandle, Lookup, and
CallSite objects are adapters governed by their profile contracts.

### Per-instruction linkage state

Each runtime Method and `invokedynamic` bytecode PC SHALL own one lazy linkage
state with these logical states:

```text
unresolved -> resolving(owner) -> linked(call site)
                              \-> failed(Java linkage error)
```

The state is protected by internal VM synchronization separate from Java
object monitors. It SHALL publish exactly one terminal result to concurrent
callers, wait/retry across execution contexts, define bounded same-owner
recursive behavior, and preserve interrupt status while waiting.

The state key includes runtime Method identity and instruction PC even when
another instruction references the same InvokeDynamic constant-pool entry.

### Bootstrap and invocation boundary

For every declared bootstrap capability, resolution SHALL:

1. obtain the call-site name and canonical logical method type;
2. resolve the bootstrap handle and ordered symbolic/static arguments in the
   caller Class lookup/access context required by that capability;
3. adapt those values through the profile's bootstrap facade or registered
   Catty Runtime bootstrap provider;
4. invoke through the typed dynamic invocation kernel;
5. validate a non-null call site and a target matching the symbolic type; and
6. publish linked state or the terminal Java linkage failure.

Invocation SHALL consume arguments and produce or throw results according to
the symbolic method type and current target. Java exceptions and object/Class
identity remain in the shared runtime world.

JVMS-required bootstrap failure translation, including the declared
`BootstrapMethodError` behavior for supported sites, belongs to the linkage
capability rather than an optional facade. Broad public Lookup behavior,
MethodHandle adaptations, CallSite APIs, and arbitrary Java bootstrap
compatibility are required only by a declared Java SE dynamic-compatibility
capability. A Catty Runtime bootstrap provider may use an internal typed
adapter, but it cannot change JVMS-visible descriptor, linkage-once,
target-type, invocation, or terminal-failure guarantees claimed by the
linkage-kernel capability.

### Engine boundary

Interpreter and IR SHALL call the same linkage service with explicit Thread,
runtime Method, and PC. IR may use a typed dynamic-call instruction or shared
service operation; its schema is not a permanent ABI. AOT support is reported
per bootstrap capability. Unsupported AOT sites are precisely build-rejected
unless a later Accepted workstream defines a typed Interpreter fallback or
direct-link policy preserving values, exceptions, initialization, and Thread
context.

## Consequences

- Runtime Method gains per-PC linkage metadata and synchronization.
- Bootstrap failures become reproducible Java linkage results rather than Go
  panics or repeated bootstrap execution.
- Catty Runtime bootstraps and Java SE-compatible bootstraps can share one
  kernel without making their public APIs or compatibility claims identical.
- String concatenation and Java lambdas become optional named compatibility
  capabilities rather than proof that the entire InvokeDynamic surface exists.
- MethodHandle adaptation remains a significant, capability-scoped semantic
  surface; broad combinators are not implied by the kernel.
- Per-site state and resolved handles increase retained runtime metadata.

## Non-scope

ConstantDynamic execution; MutableCallSite or VolatileCallSite semantics;
accepting a Catty Runtime bootstrap API; complete `java.lang.invoke` or Lookup;
arbitrary MethodHandle combinators; `invokeWithArguments`; general user-defined
bootstrap compatibility; StringConcatFactory; LambdaMetafactory; Java lambda
compatibility; AOT execution/fallback; JIT; hidden classes; or VarHandle.

## Acceptance record

Accepted by Owner on 2026-07-18 after reconciliation with ADR-0034. Acceptance
fixes the profile-scoped InvokeDynamic linkage kernel and typed call-site state.
It does not accept complete `java.lang.invoke`, StringConcatFactory,
LambdaMetafactory, arbitrary bootstrap compatibility, or production
implementation.
