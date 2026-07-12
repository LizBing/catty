# Three-party collaboration protocol

This protocol coordinates LizBing, Codex, and Claude without requiring either
agent to see the other's private chat context.

## 1. Roles

| Participant | Primary authority |
|---|---|
| LizBing | Product direction, scope, compatibility target, and final trade-offs |
| Codex | Architecture stewardship, task contracts, integration review, risk audit |
| Claude | Primary implementation, focused debugging, tests, and implementation handoff |

These are defaults, not capability limits. Codex may implement and Claude may
help design, but every active workstream has exactly one implementation owner
and one reviewer. LizBing resolves scope or product conflicts.

## 2. Sources of truth

When information conflicts, use this precedence:

1. Executable behavior and reproducible test evidence
2. Accepted ADRs
3. The active workstream contract
4. [`PROJECT_STATUS.md`](./PROJECT_STATUS.md)
5. Architecture, roadmap, and development documentation
6. Handoffs and work logs
7. Chat sessions and reconstructed session archives

Chat is a design workspace, not durable project state. Session JSONL can help
recover intent, but it cannot override the checked-in project record.

`PROJECT_STATUS.md` is the single current-state summary. `ROADMAP.md` describes
the longer strategic sequence. They must not become competing status reports.

## 3. Workstream lifecycle

Any change spanning packages, changing semantics, adding a runtime capability,
or expected to take more than one focused session gets a contract copied from
[`workstreams/TEMPLATE.md`](./workstreams/TEMPLATE.md).

Lifecycle:

1. **Proposed** — scope and decision questions are being shaped.
2. **Ready** — base commit, owner, reviewer, constraints, and gates are complete.
3. **In progress** — one implementation owner is writing code.
4. **In review** — implementation and evidence are ready for the reviewer.
5. **Accepted** — gates pass and the designated integrator accepts the result.
6. **Closed** — merged to `main`, status/docs updated, no required work remains.

Only one active implementation owner may change a workstream at a time. A
reviewer may run read-only checks and leave findings, but should not create a
competing implementation on the same branch.

## 4. Branches and worktrees

- `main` is the integrated, reviewable line.
- Concurrent work uses separate branches and preferably separate worktrees.
- Suggested names: `claude/<workstream>` and `codex/<workstream>`.
- A contract records the exact base commit. If `main` moves materially, the
  owner rebases or records why the old base remains valid before handoff.
- Do not let two agents edit the same working tree concurrently.
- Do not push directly to `main` unless LizBing explicitly requests it or the
  active contract designates the actor as integrator.

Commits use the repository-configured human identity unless LizBing requests a
different attribution. Agent participation belongs in the workstream and
handoff; agents must not invent or spoof author identity.

## 5. Handoff contract

Before yielding non-trivial work, the owner updates its `*-latest.md` file with:

- workstream and branch;
- base and head commits;
- what changed and why;
- commands run and exact outcomes;
- unresolved failures, semantic uncertainties, and dirty files;
- the next concrete action.

A handoff must distinguish facts from hypotheses. “Tests pass” names the tests;
“feature complete” points to each acceptance gate. If no code was changed, say
so explicitly.

The receiving agent first verifies branch, commit, diff, and evidence. Private
session history is optional context, not a prerequisite for continuation.

## 6. Architecture and semantic control

- Accepted ADRs are immutable records. Reverse one with a new superseding ADR.
- Proposed ADRs are not implementation authority until accepted by LizBing.
- Implementation strategy (for example, goroutines) does not weaken required
  observable Java semantics (for example, JMM happens-before guarantees).
- Unsupported behavior fails explicitly. Discovery-only stubs require an
  explicit allowlist, test, and removal criterion in the workstream.
- Every compatibility claim includes a reference behavior, usually OpenJDK,
  and a representative—not merely curated—test corpus.

## 7. Integration gates

The active contract selects proportional gates. The default set is:

```sh
gofmt -w <changed-go-files>
go vet ./...
go test ./...
bash tests/run.sh
```

Runtime and class-library changes also require the real `java.base` path used by
CI. Performance claims include the benchmark command, environment, baseline,
and repeated measurements. Documentation-only changes require link checks,
`git diff --check`, and a review for stale duplicated state.

The reviewer records one outcome: accepted, changes requested, or blocked by a
named product/architecture decision. Green CI is necessary but does not by
itself prove JVM semantic compatibility.

## 8. Status update rule

Update `PROJECT_STATUS.md` only when a merge changes a milestone, active
workstream, capability boundary, known blocker, or verified evidence. Put
session-level detail in `docs/worklog/` and the latest handoff. This keeps the
shared state short enough that every agent can safely reread it before work.
