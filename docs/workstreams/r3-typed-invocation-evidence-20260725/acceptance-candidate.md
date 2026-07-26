# R3 typed invocation acceptance candidate

**Date:** 2026-07-25
**Workstream:** [`../r3-typed-invocation-kernel-slice.md`](../r3-typed-invocation-kernel-slice.md)
**Review:** Owner
**Status:** candidate evidence; not Ready or Done

## Static-field composite evidence

Typed direct static-field behavior is covered at three distinct boundaries:

| Boundary | Evidence | Required property |
|---|---|---|
| Real classfile execution | `bash tests/run.sh` / `StaticFields` | Existing classfile `getstatic`/`putstatic` agrees with Java in tree Interpreter and IR. |
| Shared typed storage | `TestTypedFieldAccessRoundTripsAllKinds` | Static values cover every primitive, category-2 values, references, and null without exposing HeapCell layout. |
| Context-aware direct adapter | `TestInvokeDispatchInterpreterAndIR` | Static dispatch initializes the declarer's Class before execution, in both tree and IR adapters. |

The three layers deliberately remain separate: the first proves bytecode
semantics against Java, the second proves the stable typed value/storage
boundary, and the third proves the direct invocation initialization transition.
No single test is presented as evidence for all three concerns.

## Current direct-invocation matrix

| Concern | Tree Interpreter | IR | AOT |
|---|---|---|---|
| Typed normal/category-2 result | unit tested | unit tested | explicitly rejected |
| Abrupt Java throwable identity | unit tested | unit tested | explicitly rejected |
| Static/virtual/interface/special/constructor dispatch | unit tested | unit tested where execution is interpreted | explicitly rejected |
| Static initialization trigger | unit tested | unit tested | no typed dynamic fallback |
| Caller access | unit tested | shared dispatch rule | explicitly rejected |
| Nested typed invocation boundary | unit tested | unit tested | explicitly rejected |

## Remaining finalization work

- Run and record the complete workstream gate set against a fixed candidate
  commit, including the project-prescribed R2 evidence-isolation checks.
- Review whether the bounded direct access model needs any additional
  Java-visible compatibility policy before Owner acceptance.
- Owner review and explicit transition to `Ready`/`Done`; this record does not
  make either transition.
