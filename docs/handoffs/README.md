# Agent handoffs

[`claude-latest.md`](./claude-latest.md) and
[`codex-latest.md`](./codex-latest.md) are overwrite-in-place pointers to each
agent's latest durable handoff. Workstream files retain milestone history; the
monthly work log retains a concise chronological record.

Each handoff states branch and commit identities, tests actually run, dirty
state, unresolved risks, and exactly one recommended next action. The receiving
agent verifies those claims from the repository before continuing.
