# <workstream>: <title>

**Status:** Proposed
**Type:** research / implementation
**Review:** self / owner
**Base commit:** `<commit>`
**Roadmap item:** `<milestone / item / owner-directed exception>`
**Governing ADRs:** `<Accepted ADRs or None>`
**Prerequisites:** `<workstreams, capabilities, decisions, or None>`
**Acceptance anchor:** `N/A while Proposed; <accepted-contract commit> before implementation`

## Outcome

一句话描述完成后可观察到的结果。

## Scope

- 本 workstream 实现什么。

## Non-scope

- 明确不实现什么。

## Semantic constraints

- 必须保持的 JVM/JRE 行为、失败边界和引擎一致性。

## Acceptance

| Gate | Command / artifact | Result |
|---|---|---|
| 示例行为 | `<exact command>` | Not run |
| 回归 | `<exact command>` | Not run |
| Evidence isolation | `<exact check proving historical evidence unchanged and candidate output path>` | Not run |

结果只使用 `Pass`、`Fail`、`Not run`、`Not implemented`。

## Amendments

Accepted 后只在此追加由 Owner 接受的需求变化，不回写降低原合同。

---

## Implementation preflight

开始 production implementation 前记录：

- **Acceptance anchor / actual base:** `<commits; worktree must descend from anchor>`
- **Historical evidence check:** `<exact command and exit status>`
- **Candidate evidence destination:** `<workstream-specific path containing candidate ID>`
- **Harness output policy:** `<must not default to historical/shared path>`

任何一项缺失时，保持 `Accepted`，不得转为 `In Progress`。

---

## Plan

| Slice | Status | Evidence |
|---|---|---|
| A | Pending | — |

状态使用 `Pending`、`In progress`、`Complete`。

---

## Handoff

- **Branch / candidate:**
- **Acceptance anchor / base:**
- **Dirty files:**
- **Historical evidence check:**
- **Candidate evidence path:**
- **Last location:**
- **Checks run / not run:**
- **Blocker:**
- **Next action:**
- **Non-derivable context:**
