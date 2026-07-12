# <workstream ID>: <title>

**Status:** Proposed  
**Owner:** Unassigned  
**Reviewer:** Unassigned  
**Integrator:** LizBing  
**Base commit:** `<full commit>`  
**Branch:** `<agent/workstream>`  
**Target milestone:** `<milestone>`

## Outcome

State one externally observable result. Avoid an implementation checklist as
the definition of success.

## Context

Explain why this work is needed and link relevant ADRs, failures, benchmarks,
or earlier contracts.

## In scope

- ...

## Out of scope

- ...

## Semantic constraints

- Name the relevant JVMS/JLS/JRE behavior and reference implementation.
- State which execution paths must agree: interpreter, IR, AOT, native.
- State how unsupported behavior fails.

## Design and decision questions

| Question | Owner | Resolution |
|---|---|---|
| ... | LizBing/Codex/Claude | Open |

## Implementation slices

| Slice | Result | Dependencies | Status |
|---|---|---|---|
| A | ... | None | Pending |

## Acceptance gates

- [ ] Required Go unit tests pass.
- [ ] Required Java differential fixtures pass against the named JDK.
- [ ] Real `java.base` path passes when runtime/class-library behavior changes.
- [ ] Unsupported paths fail with the specified error.
- [ ] Documentation and ADRs match the implementation.
- [ ] Performance evidence is attached when a performance claim changes.

Replace these defaults with concrete test names and commands before changing
the status to **Ready**.

## Evidence

| Gate | Command/artifact | Result |
|---|---|---|
| ... | ... | Pending |

## Risks and rollback

- Risk: ...
- Detection: ...
- Rollback or containment: ...

## Handoff history

| Date | From | To | Commit | Summary |
|---|---|---|---|---|
