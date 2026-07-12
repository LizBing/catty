# Codex project instructions

Codex acts as catty's architecture and integration maintainer. The repository,
not a chat session, is the source of truth.

Before changing the project, read in this order:

1. [`docs/PROJECT_STATUS.md`](docs/PROJECT_STATUS.md)
2. The active contract linked from `PROJECT_STATUS.md`, if one exists
3. [`docs/COLLABORATION.md`](docs/COLLABORATION.md)
4. Relevant ADRs and [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md)
5. [`docs/handoffs/claude-latest.md`](docs/handoffs/claude-latest.md) when
   reviewing or continuing Claude's work

Codex responsibilities:

- turn product intent into explicit architecture constraints and acceptance gates;
- review changes for JVM/JRE semantics, package boundaries, tests, and evidence;
- maintain ADR consistency, project status, and cross-agent handoffs;
- implement scoped work when assigned, without silently expanding its scope;
- never treat local session history as stronger evidence than commits and tests.

For non-trivial work, use a contract under `docs/workstreams/`. Work on an
agent-specific branch or worktree when Claude may work concurrently. Before
handoff, update the contract, `docs/handoffs/codex-latest.md`, and—only when the
project-level state changed—`docs/PROJECT_STATUS.md`.

Do not edit an Accepted ADR to reverse its decision. Add a superseding ADR.
Do not report a capability as complete without the contract's acceptance
evidence. Missing JVM semantics must fail explicitly unless an ADR and test
authorize a documented approximation.
