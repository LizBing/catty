# ONBOARD-1: Claude/DeepSeek collaboration onboarding

**Status:** Ready  
**Owner:** Claude (DeepSeek backend selected manually by LizBing)  
**Reviewer:** Codex  
**Integrator:** Codex when explicitly accepted by LizBing  
**Base commit:** `cc1379775e8a8407882d89615b62d3c5546db5e9`  
**Branch:** `claude/onboarding-deepseek`  
**Target milestone:** Collaboration validation before R2

## Outcome

Demonstrate that a fresh Claude Code session with no private Codex chat context
can recover catty's verified state, authority boundaries, and next decisions
from the repository alone, then leave a precise durable handoff.

## Context

LizBing, Codex, and Claude adopted the repository-first protocol in COLLAB-1.
Before assigning runtime implementation, this low-risk exercise validates that
Claude can follow the protocol while backed by DeepSeek and that Codex can
supervise and review the resulting session.

## In scope

- Read `CLAUDE.md` and every document it requires for an unassigned workstream.
- Inspect the repository and Git baseline using read-only commands.
- Explain, in Claude's own words:
  - source-of-truth precedence and participant roles;
  - verified R1 capabilities versus explicitly unsupported behavior;
  - the four decisions required before R2 implementation;
  - Claude's branch, commit, test, and handoff permissions.
- Replace `docs/handoffs/claude-latest.md` with a live onboarding handoff.
- Record commands actually run and distinguish observations from inference.

## Out of scope

- Any Go, Java, CI, ADR, roadmap, or architecture modification.
- Designing or implementing R2.
- Editing `docs/PROJECT_STATUS.md`.
- Committing, pushing, merging, or changing branches during the Claude session.
- Reading recovered private session JSONL; repository documents are sufficient.

## Semantic constraints

- No JVM/JRE behavior changes.
- Claude must recognize that Go implementation mechanisms do not automatically
  replace required Java observable semantics.
- Green curated tests must not be presented as arbitrary `java.base` support.
- Unresolved native behavior must not be normalized as compatible zero values.

## Design and decision questions

| Question | Owner | Resolution |
|---|---|---|
| Can DeepSeek-backed Claude recover state from the repository alone? | Codex | To be evaluated |
| May Claude begin R2 implementation in this session? | LizBing/Codex | No |
| Who may accept the onboarding evidence? | Codex | Codex reviews; LizBing owns final process direction |

## Implementation slices

| Slice | Result | Dependencies | Status |
|---|---|---|---|
| A | Read and verify shared project state | Clean worktree at base | Pending |
| B | Produce structured understanding report | A | Pending |
| C | Write live `claude-latest.md` handoff only | B | Pending |
| D | Codex independently reviews diff and claims | C | Pending |

## Acceptance gates

- [ ] Claude identifies `docs/PROJECT_STATUS.md` as the current-state source.
- [ ] Claude distinguishes Accepted ADRs from Proposed ADRs.
- [ ] Claude states the R1 corpus boundary without overclaiming compatibility.
- [ ] Claude identifies JMM, strict unresolved natives, minimum Unsafe semantics,
      and representative concurrency/java.base programs as pre-R2 decisions.
- [ ] Claude makes no runtime, ADR, roadmap, project-status, or CI changes.
- [ ] The only working-tree diff is `docs/handoffs/claude-latest.md`.
- [ ] The handoff records branch, base/head, commands, findings, risks, dirty
      state, and one next action.
- [ ] Codex can reproduce all repository claims without the Claude session.

## Evidence

| Gate | Command/artifact | Result |
|---|---|---|
| Repository baseline | `git status`, `git rev-parse HEAD` | Pending |
| Claude understanding | Supervised Claude Code output | Pending |
| Durable handoff | `docs/handoffs/claude-latest.md` | Pending |
| Scope control | `git diff --name-only` | Pending |

## Risks and rollback

- Risk: Claude begins implementation despite the narrow contract.
- Detection: any diff outside `docs/handoffs/claude-latest.md`.
- Risk: provider switching did not take effect.
- Detection: Claude reports unexpected provider/model metadata when available;
  otherwise LizBing's CC Switch selection is the authority.
- Rollback: discard the isolated worktree diff and remove the onboarding branch.
  The integrated repository remains unchanged.

## Handoff history

| Date | From | To | Commit | Summary |
|---|---|---|---|---|
| 2026-07-12 | Codex | Claude | `cc13797` | Fresh-session onboarding contract ready |
