# R2: runtime semantics research

**Status:** Accepted
**Type:** research
**Review:** owner
**Base commit:** `5720147`
**Roadmap item:** Phase R2 — Runtime semantics and concurrency planning
**Governing ADRs:** ADR-0016 through ADR-0024
**Prerequisites:** Governance baseline integrated on `main`

## Outcome

Produce evidence-backed Proposed ADRs and a bounded implementation contract for
the first R2 runtime-semantics slice; do not add production concurrency or
class-library capabilities in this workstream.

## Scope

- Establish a Java 25/Temurin 25 differential fixture set for class and
  interface initialization, failure, and relevant bootstrap behavior.
- Map current R1 initialization, bootstrap, String, Slot, and AOT/interpreter
  bridge behavior against ADR-0017 through ADR-0024.
- Produce proposed designs for the detailed class-initialization state machine,
  bootstrap-kernel capability graph, and String representation experiment.
- Define the smallest candidate R2 implementation slice, its engine matrix,
  explicit non-support behavior, semantic constraints, and acceptance gates.

## Non-scope

- Production Java concurrency, monitors, Java Memory Model, Unsafe,
  `invokedynamic`, broad reflection, broad java.base compatibility, or a JIT.
- Refactoring Slot, object layout, String representation, bootstrap classes, or
  AOT bridge code into production.
- Treating research prototypes or benchmark wins as accepted architecture.

## Semantic constraints

- Java 25 is the semantic baseline; Temurin 25 is the pinned differential
  reference under ADR-0017.
- Unsupported paths must be identified explicitly; no silent approximation is
  accepted as a compatibility result.
- Research reports distinguish Interpreter, IR, and AOT evidence.
- No accepted ADR is reversed by this workstream; any revision is proposed as a
  new ADR or an amendment for Owner acceptance.

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| Initialization evidence | Versioned fixture matrix and Temurin 25 comparison report | Not run |
| Bootstrap evidence | Capability/dependency graph with R1 implementation mapping | Not run |
| String evidence | UTF-16 edge-case matrix and representation trade-off report | Not run |
| R2 proposal | Proposed ADRs plus a bounded implementation workstream contract | Not run |
| Governance consistency | ADR/status/roadmap links and `git diff --check` | Not run |

## Amendments

None.

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — baseline and fixture design | Pending | — |
| B — initialization and bootstrap analysis | Pending | — |
| C — String and representation analysis | Pending | — |
| D — proposed R2 implementation contract | Pending | — |

---

## Handoff

- **Branch / candidate:** Governance changes awaiting integration
- **Dirty files:** Governance ADR and planning documents
- **Last location:** R2 research is accepted but has not started
- **Checks run / not run:** Document checks pending
- **Blocker:** Governance baseline must be committed before research work begins
- **Next action:** Integrate governance baseline, then begin Slice A
- **Non-derivable context:** This is a research-only workstream; no production implementation is authorized.
