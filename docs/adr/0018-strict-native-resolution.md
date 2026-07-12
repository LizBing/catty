# ADR-0018: Strict unresolved-native failure

- **Status:** Accepted
- **Date:** 2026-07-12

## Context

`rtda.InterpretedMethod` currently assigns every unresolved `ACC_NATIVE` method
a return-type-based function that pushes zero or null. This was useful for
dependency exploration, but it turns missing behavior into plausible output.

The R2 research probes demonstrate the failure mode: Integer/Long decimal
conversion exits successfully while printing NUL bytes because
`Unsafe.putByte` silently does nothing. A successful exit is therefore not
evidence of compatible execution.

Some explicit native functions are also named `nop`, `nopBool0`, `nopRef0`, or
return zero. Those names do not establish that the Java contract permits the
behavior; every registration needs classification.

## Decision

Remove the generic native zero/null stub. A declared but unresolved native has
no executable implementation. Invocation throws a catchable Java
`UnsatisfiedLinkError` containing class, method, and descriptor context.

The same rule applies to interpreter, IR, AOT/runtime bridge, and real
java.base classes. Dependency discovery may inventory unresolved signatures but
still throws on invocation.

Every registered native is classified as:

- Implemented;
- Semantic no-op;
- Compatibility adapter;
- Unsupported.

Only a reviewed signature can be a no-op or compatibility adapter. There is no
global “continue with zero” mode. Temporary shims require a signature, test,
owner, and removal criterion.

## Consequences

**Positive**

- Missing semantics fail at their real boundary.
- Capability tests cannot pass because a native returned a convenient zero.
- Unsafe and java.base expansion becomes caller-driven and auditable.

**Negative**

- Programs that previously limped forward may now fail earlier.
- The native invocation path must create Java exceptions consistently across
  all execution engines.
- Existing no-op/zero registrations require an explicit audit and tests.
