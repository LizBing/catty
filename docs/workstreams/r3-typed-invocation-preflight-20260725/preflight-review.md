# R3 typed invocation implementation preflight

**Date:** 2026-07-25
**Workstream:** [`r3-typed-invocation-kernel-slice`](../r3-typed-invocation-kernel-slice.md)
**Acceptance anchor:** `65630453cdb7579c85f602b21c2d9f459a042f88`
**Actual review base:** `cc1da5dbe3cc343ef11db5cf9b394f7cf943262b`
**Mode:** record and review only; no production implementation started

## Authorization state

The workstream is Accepted with owner review, but implementation remains paused
until the Owner explicitly accepts this preflight and authorizes the start. This
record does not change the frozen Outcome, Scope, Non-scope, Semantic
constraints, Acceptance gates, profile classification, or review type.

The actual review base is a descendant of the acceptance anchor and contains
only the already-integrated execution-context boundary/status governance work
after the fixed anchor. No implementation branch has been created for this
slice.

## Governing constraints checked

- ADR-0016: typed invocation must preserve shared object world, explicit
  normal/exception transitions, and per-engine capability claims.
- ADR-0018: Go panic, goroutine, mutex, and channel behavior cannot substitute
  for Java-visible semantics.
- ADR-0020: interpreter `Slot`, IR registers, heap-cell bits, and AOT bridge
  shapes are adapters, not a stable cross-engine value contract.
- ADR-0021/0025/0029: initialization remains a shared runtime service and
  invocation must preserve synchronized-method and class-initialization
  behavior on normal and abrupt paths.
- ADR-0028: invocation APIs must receive explicit `ExecutionContext`; no code
  may discover Java Thread identity from a carrier.
- ADR-0030: heap access must use typed heap-cell APIs; no direct heap-layout
  exposure can become part of the invocation result model.
- ADR-0031: the slice owns only the shared typed dynamic invocation kernel.
  Java reflection conversions, wrappers, boxing, varargs, access checks,
  `InvocationTargetException`, MethodHandle, InvokeDynamic, generated classes,
  Host provider binding, and AOT dynamic execution remain outside this slice.
- ADR-0033/0034: generated-class and Host ABI vocabulary may reuse the kernel,
  but their public policies and authority checks are separate future contracts.

## Code boundary findings

1. `rtda.ExecutionContext` is now the canonical execution state and preserves
   the temporary `rtda.Thread` compatibility alias. Typed invocation should use
   `*rtda.ExecutionContext` in new APIs and avoid adding new `Thread`-named
   entry points.
2. Interpreter and native invocation still use `*rtda.Frame` plus `rtda.Slot`
   internally. That is acceptable as an adapter layer, but the new stable value
   model must not expose those frame or slot layouts.
3. The AOT runtime bridge still exposes `runtime.InvokeVirtual`,
   `InvokeSpecial`, `InvokeStatic`, and `interpreter.RunMethod` with
   `[]rtda.Slot` and `rtda.Slot` results. This is the most important boundary
   to quarantine: the first implementation should add typed kernel/adapters
   before altering generated AOT bridge behavior.
4. Category-2 values already have typed heap accessors, but bridge returns still
   panic for long/double through the old runtime bridge. The typed result matrix
   must make long/double one logical value.
5. Java throwable transport is currently represented as pending exception on
   `ExecutionContext`; the new result model must distinguish normal zero/null,
   Java throwable, and internal provider failure without relying on Go panic.
6. `rtda.Method` still assigns zero/null default stubs for missing classfile
   natives. ADR-0034 rejects that long-term Host ABI behavior, but this
   workstream does not authorize a native-binding policy rewrite. Treat it as a
   recorded adjacent debt, not a typed-invocation gate.

## Recommended first implementation shape

Start with a small shared kernel package or file owned by runtime data services
rather than the AOT bridge:

- define a logical Java type/value/result vocabulary for void, boolean, byte,
  char, short, int, long, float, double, and reference;
- define result states for normal, Java throwable, and internal failure;
- add descriptor-aware frame and heap adapters with exhaustive unit tests;
- add direct method/field/constructor kernel tests for Interpreter and IR
  consumers only after the value/result matrix is stable;
- keep AOT dynamic invocation precisely rejected and do not mark any R3
  Java-visible row newly Supported.

## Review conclusion

The slice is ready for Owner review of this preflight. Implementation should
not start until the Owner explicitly accepts this record and authorizes moving
the workstream from Accepted to In Progress.
