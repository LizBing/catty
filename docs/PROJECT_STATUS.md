# Project status

**As of:** 2026-07-12  
**Baseline verified through:** `89037c4`  
**Branch:** `main`  
**Milestone:** R1 complete and hardened  
**Active workstream:** [`R2-A`](./workstreams/R2-A-STRICT-NATIVE.md) — Ready; Claude/DeepSeek owner, Codex reviewer

This is catty's single current-state summary. The strategic sequence lives in
[`ROADMAP.md`](./ROADMAP.md); session detail belongs in
[`worklog/`](./worklog/) and [`handoffs/`](./handoffs/).

## Verified baseline

- Interpreter: approximately 145 opcodes, exceptions, interface dispatch,
  multidimensional arrays, and class initialization.
- Class loading: provider chain plus real `java.base` auto-detection through a
  JDK-extracted image.
- Native/bootstrap layer: six irreducible synthetic bootstrap classes,
  additional synthetic fallbacks, and approximately 40 native registrations.
- AOT: standalone Go binary path; `fib(35)` remains approximately 40–60 ms on
  the recorded development machine.
- Regression evidence: unit tests, three-engine fixture comparison, and the
  real `java.base` smoke path run in CI.
- Historical reconstruction and independent review are archived in
  [`CLAUDE_SESSION_HISTORY.md`](./CLAUDE_SESSION_HISTORY.md) and
  [`CODEX_SESSION_REVIEW.md`](./CODEX_SESSION_REVIEW.md).

## Explicit capability boundary

The following are not R1 capabilities:

| Capability | Planned phase | Current blocker |
|---|---|---|
| `Integer/Long.toString` | R2 prerequisite | DecimalDigits uses unresolved Unsafe array writes; generic stub produces NUL output |
| `Double.parseDouble` | R2 prerequisite investigation | FloatingDecimal probe times out; no direct Unsafe edge found |
| Representative basic `HashMap` | R2 prerequisite investigation | `jdk.internal.misc.VM` reports “Not yet initialized”; basic path has no direct Unsafe edge |
| Java concurrency and monitors | R2 | Thread lifecycle, monitor semantics, JMM guarantees |
| `invokedynamic`, reflection, annotations | R3 | Dynamic metadata and call-site semantics |
| Broad Java I/O and networking | R4 | Native/runtime integration design |
| Broad real-program AOT coverage | R5 | Instance methods, exceptions, dynamic calls |

“18/18 smoke tests pass” means the selected R1 corpus passes; it is not a claim
that arbitrary `java.base` applications are supported.

## Accepted R2 directions and remaining gates

LizBing accepted G1–G4 on 2026-07-12:

1. Protect DRF/final/volatile/monitor/Thread/class-init semantics and measure
   Strict, Go-native, and Hybrid storage before considering a named
   racy-program deviation.
2. Unresolved natives throw `UnsatisfiedLinkError`; no generic zero/null stub.
3. Unsafe uses logical offsets and caller-backed U0–U4 profiles.
4. R2 requires deterministic differential tests, stress, timeouts, and race
   evidence.

ADR-0016 through ADR-0019 are Accepted and ADR-0011 is superseded. The exact
JDK 25 caller graph and R2 test plan are attached. No R2 implementation owner
is currently assigned.

## Known architecture risks

- Synthetic String's Go payload and OpenJDK's field layout form a dual
  representation that must be reconciled before deep reflection support.
- Unsafe, monitors, atomics, and Thread lifecycle are one semantic cluster;
  isolated zero-value stubs can create plausible but incorrect execution.
- The current smoke corpus is curated and should grow into capability matrices
  plus representative application tests.
- Interpreter/AOT/native paths must keep one observable semantic contract even
  when their implementation mechanisms differ.

## Next coordination action

Claude/DeepSeek implements R2-A in `claude/r2-a-strict-native` and leaves a
tested handoff. Codex independently reviews the diff, inventory, Java exception
semantics, and all engine evidence before any integration.
