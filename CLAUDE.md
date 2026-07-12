# Claude project instructions

Claude is a primary implementation agent for catty. The repository, not a chat
session, is the source of truth.

Before changing the project, read in this order:

1. [`docs/PROJECT_STATUS.md`](docs/PROJECT_STATUS.md)
2. The active contract linked from `PROJECT_STATUS.md`, if one exists
3. [`docs/COLLABORATION.md`](docs/COLLABORATION.md)
4. Relevant ADRs and [`docs/DEVELOPMENT.md`](docs/DEVELOPMENT.md)
5. [`docs/handoffs/codex-latest.md`](docs/handoffs/codex-latest.md) when
   implementing or continuing work planned or reviewed by Codex

For non-trivial work, do not implement from chat memory alone. Confirm that a
workstream contract states the base commit, owner, goals, non-goals, semantic
constraints, and acceptance gates. If it does not, stop and create or request
the missing contract before making broad changes.

Work on a Claude-specific branch or worktree when Codex may work concurrently.
Keep changes inside the contract. Before handoff, run the required checks and
update the contract's evidence plus `docs/handoffs/claude-latest.md`. Update
`docs/PROJECT_STATUS.md` only when project-level state has genuinely changed.

Do not edit an Accepted ADR to reverse its decision. Add a superseding ADR.
Do not substitute Go behavior for required JVM/JRE semantics without an
explicit decision. Missing semantics must fail loudly rather than return a
plausible zero value unless the contract explicitly authorizes a discovery stub.
