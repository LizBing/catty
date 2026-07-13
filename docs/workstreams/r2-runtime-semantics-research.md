# R2: runtime semantics research

**Status:** Done
**Type:** research
**Review:** owner
**Base commit:** `ecb086e`
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
| Initialization evidence | `docs/workstreams/r2-evidence/matrix.md` + `docs/workstreams/r2-evidence/run-r2-results.txt` | Pass — 16 fixtures, 3-engine matrix run against Temurin 25.0.3; 9 initialization fixtures with per-engine results recorded |
| Bootstrap evidence | `docs/workstreams/r2-evidence/reports/r2-bootstrap-graph.md` | Pass — 7 candidate capabilities mapped to R1 providers and minimum observed responsibilities; facade/provider decisions explicitly deferred |
| String evidence | `docs/workstreams/r2-evidence/reports/r2-string-matrix.md` + `docs/workstreams/r2-evidence/reports/r2-string-representation.md` | Pass — UTF-16 fixture matrix and representation analysis: `[]uint16` proposed as kernel backing; facade, bridge ABI, and host-text policy explicitly deferred |
| R2 proposal | `docs/adr/0025-*.md` (Accepted) + `docs/adr/0027-*.md` (Accepted) + two Proposed implementation contracts | Pass — both architecture decisions accepted; both bounded R2 implementation contracts remain Proposed for owner review |
| Governance consistency | `git diff --check` | Pass — document diff is clean; production validation was not re-run because this workstream changed no production code |

## Amendments

None.

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A — baseline and fixture design | Complete | `docs/workstreams/r2-evidence/`: 16 fixtures, `run-r2-diff.sh` harness, `matrix.md`, `run-r2-results.txt` — 16/16 fixtures compiled, 3-engine differential run against Temurin 25.0.3 |
| B — initialization and bootstrap analysis | Complete | `reports/r2-init-deltas.md` (7-area JVMS §5.5 gap analysis, GetstaticOwner crash confirmed as bug), `reports/r2-bootstrap-graph.md` (7-capability → R1 mapping, java.base gating analysis) |
| C — String and representation analysis | Complete | `reports/r2-string-matrix.md` (six-fixture UTF-16 matrix), `reports/r2-string-representation.md` (four-candidate analysis; `[]uint16` kernel-backing recommendation) |
| D — proposed R2 implementation contracts | Complete | `docs/adr/0025-*.md` (Accepted — class-init), `docs/adr/0027-*.md` (Accepted — UTF-16 String kernel backing), `r2-initialization-slice.md` and `r2-string-utf16-slice.md` (Proposed) |

---

## Handoff

- **Branch / candidate:** `main` (dirty — research artifacts uncommitted)
- **Dirty files:** Research evidence and fixtures under `docs/workstreams/r2-evidence/`, 2 Proposed ADRs under `docs/adr/`, 1 Proposed implementation contract under `docs/workstreams/`, plus their indexes/status documents
- **Last location:** All 4 slices Complete; 5/5 gates scored Pass. Owner accepted ADR-0025, ADR-0027, and `r2-initialization-slice`; research workstream is Done.
- **Checks run / not run:** Differential harness run (16/16 fixtures compiled, full matrix populated) and `git diff --check` passed; `go vet ./...`, `go test ./...`, `go test -race ./...`, `bash tests/run.sh` were not re-run because production code is unchanged
- **Blocker:** None — `r2-initialization-slice` is Accepted. The separate String contract remains Proposed.
- **Next action:** Assign one Active Agent to `r2-initialization-slice` on an implementation branch/worktree.
- **Non-derivable context:** The research harness requires `catty build` to be invoked from the repo root for AOT (module-context issue documented in `matrix.md`). The 16 differential fixtures are concat-free (Java 25 `javac` emits `invokedynamic` for `+`; catty defers `invokedynamic` to R3). The `ReachUnsafe` fixture exercises the bootstrap boundary in pure-synthetic mode (no java.base extracted) — extracting java.base with `jimage` is an Owner decision for broader coverage.
