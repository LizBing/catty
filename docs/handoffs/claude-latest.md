# Claude latest handoff

**Date:** 2026-07-12
**Role:** Implementation owner (onboarding validation)
**Workstream:** [`ONBOARD-1`](../workstreams/ONBOARD-1-CLAUDE-DEEPSEEK.md)
**Branch:** `claude/onboarding-deepseek`
**Base commit:** `cc13797`
**Head commit:** `1e1fd1a`
**Model/backend:** DeepSeek, selected and confirmed by LizBing through CC Switch.
CC Switch routes Claude model aliases to the selected backend, so CLI result
metadata identifies the requested alias (`claude-opus-4-8[1m]`) rather than the
underlying provider.

## What this session did

Completed ONBOARD-1 slices A, B, and C: read `CLAUDE.md` and all documents it
requires, inspected the repository baseline with read-only commands, produced a
structured understanding report, and wrote this live handoff.

## Commands actually run

```sh
git rev-parse HEAD                          # 1e1fd1a
git branch --show-current                   # claude/onboarding-deepseek
git status                                  # clean working tree
git log --oneline -10                       # recent history
git merge-base cc13797 HEAD                 # cc13797 (base is at HEAD)
git log --oneline cc13797..HEAD             # 1 commit ahead: workstream doc
git diff --name-only HEAD~1..HEAD           # files in the initial workstream commit
git diff --check                            # no trailing whitespace
```

All commands were read-only. No `go build`, `go test`, `go vet`, `tests/run.sh`,
or any other build/test invocation was performed — this session is documentation
inspection only.

## Findings

### Observed facts

- `docs/PROJECT_STATUS.md` is the single current-state summary, as required by
  `COLLABORATION.md` §2 and §8.
- Of 15 ADRs, 9 are Accepted (0001–0007, 0014, 0015) and 6 are Proposed
  (0008–0013). ADR-0011 (Go memory model), ADR-0010 (Thread = goroutine),
  ADR-0008 (AOT-first), and ADR-0009 (hybrid class library) are all Proposed,
  not Accepted — they are design proposals, not implementation authority
  (`COLLABORATION.md` §6).
- The 14 Java fixture files in `tests/fixtures/` (including `RealBaseSmoke.java`,
  conditional on `CATTY_BOOT`) define the R1 e2e corpus. `RealBaseSmoke.java`
  explicitly documents 18 assertion groups and calls out what is *deliberately
  excluded* (HashMap, Double.parseDouble, Integer/Long.toString — all blocked
  by the DecimalDigits → Unsafe cascade). The contract (`ONBOARD-1` §Semantic
  constraints) and `PROJECT_STATUS.md` both warn against presenting this as
  arbitrary `java.base` support.
- The four pre-R2 decisions are listed verbatim in `PROJECT_STATUS.md` §
  "Decisions required before R2 implementation" and match the contract's
  acceptance gates.
- HEAD (`1e1fd1a`) is one commit ahead of the contract base (`cc13797`); that
  commit adds the ONBOARD-1 workstream document, `PROJECT_STATUS.md` update,
  `codex-latest.md`, and a worklog entry — all documentation. The working tree
  is otherwise clean.
- `claude-latest.md` previously contained a Codex-authored recovery record
  stating "This is a synthesized recovery record, not a new live handoff
  authored by Claude. Claude should replace this file after its next assigned
  workstream."

### Inferences (not directly observable)

- The session cannot independently inspect CC Switch's routing internals. The
  DeepSeek attribution is configuration information confirmed by LizBing; the
  CLI-visible `claude-opus-4-8[1m]` value is the routed request alias.
- The repository appears well-structured for two-agent collaboration: separate
  branch naming conventions (`claude/` and `codex/`), handoff files, and
  workstream contracts give each agent independent read access to the other's
  state without shared chat context.

## Unresolved risks

- The four pre-R2 decisions (ADR-0011 direction, strict native resolution,
  minimum Unsafe surface, representative acceptance programs) remain open.
  There is no R2 workstream contract yet.
- The `claude-latest.md` being replaced was a Codex-authored synthesis, so
  this is the first Claude-authored handoff on this branch.
- Backend identity cannot be independently derived from Claude CLI metadata;
  future audits should treat CC Switch's user-confirmed selection as the source
  for the routed provider and the CLI value as the requested model alias.

## Dirty state

Before this edit, the working tree was clean. After this edit, exactly one file
is modified and uncommitted:
`docs/handoffs/claude-latest.md`.

No other files are staged, modified, or untracked.

## Next action

**Codex reviews this handoff against the ONBOARD-1 acceptance gates**, then
LizBing confirms whether the collaboration validation is Accepted. No R2
implementation begins before that gate passes and an R2 architecture-gates
workstream contract is drafted and reaches **Ready** (`PROJECT_STATUS.md` §
"Next coordination action").
