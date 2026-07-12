# COLLAB-1: Three-party collaboration foundation

**Status:** Closed  
**Owner:** Codex  
**Reviewer:** LizBing  
**Integrator:** LizBing or Codex when explicitly requested  
**Base commit:** `89037c41a01836c777d85c4d480c4ea4be727af1`  
**Branch:** `main` (single-agent documentation change; uncommitted during review)  
**Target milestone:** Project governance before R2

## Outcome

LizBing, Codex, and Claude can continue catty from durable repository state
without access to one another's private sessions, while preserving one owner,
one reviewer, and explicit evidence for every non-trivial workstream.

## Context

The accidental Claude `/clear` demonstrated that `.claude/progress.md` and
private chat context were insufficient integration boundaries. The recovered
sessions also showed that implementation checklists could be mistaken for
semantic completion. R2 increases that risk because Thread, Unsafe, monitors,
and memory visibility cross multiple packages and execution engines.

## In scope

- Shared roles, source precedence, branches, handoffs, and integration gates.
- A single model-neutral current-state document.
- Tool-specific entry files that point to shared state.
- Reusable workstream, handoff, and work-log structures.
- Initial durable handoffs based on the recovered R1 evidence.

## Out of scope

- Deciding the R2 JMM or Unsafe design.
- Changing runtime code or CI.
- Automating session export or GitHub project management.
- Assigning Claude an R2 implementation task.

## Semantic constraints

- This work changes no JVM behavior.
- Governance must preserve the stated goal of JRE semantics and must not treat
  a Go implementation mechanism as an implicit semantic waiver.
- Missing capabilities are reported explicitly rather than inferred from green
  curated tests.

## Design and decision questions

| Question | Owner | Resolution |
|---|---|---|
| Who owns architecture and integration continuity? | LizBing | Codex |
| Who is the default primary implementation agent? | LizBing | Claude |
| What synchronizes private sessions? | Codex | Repository contracts and handoffs, not bidirectional chat access |
| What is the current-state source? | Codex | `docs/PROJECT_STATUS.md` |

## Implementation slices

| Slice | Result | Dependencies | Status |
|---|---|---|---|
| A | Shared collaboration protocol | Role decision | Complete |
| B | Project status and private-status migration | A | Complete |
| C | Workstream/handoff/worklog templates | A | Complete |
| D | Claude and Codex entry instructions | B, C | Complete |
| E | README/development links and validation | A–D | Complete |

## Acceptance gates

- [x] Both agents have a root entry file with the same shared-state order.
- [x] `.claude/progress.md` no longer duplicates project status.
- [x] Workstream contracts identify owner, reviewer, base, scope, and evidence.
- [x] Handoffs distinguish recovered history from live agent-authored claims.
- [x] Relative Markdown links resolve.
- [x] `git diff --check` passes.
- [x] LizBing accepts the collaboration model.

## Evidence

| Gate | Command/artifact | Result |
|---|---|---|
| Entry symmetry | `AGENTS.md`, `CLAUDE.md` | Pass |
| Single status source | `.claude/progress.md`, `docs/PROJECT_STATUS.md` | Pass |
| Contract schema | `docs/workstreams/TEMPLATE.md` | Pass |
| Relative links | Local relative-link existence scan | Pass; no broken links |
| Whitespace | `git diff --check` | Pass |
| Governance acceptance | LizBing review, 2026-07-12 | Accepted |
| Runtime regression | Not run | Documentation-only change |

## Risks and rollback

- Risk: process overhead slows small fixes.
- Detection: trivial changes repeatedly create unnecessary contracts.
- Containment: contracts are required only for cross-package, semantic, or
  multi-session work; focused local fixes use normal commit messages and tests.
- Rollback: remove the new governance files and restore `.claude/progress.md`
  from the base commit. No runtime state is affected.

## Handoff history

| Date | From | To | Commit | Summary |
|---|---|---|---|---|
| 2026-07-12 | Codex | LizBing | Integration commit containing this contract | Reviewed, accepted, and closed |
