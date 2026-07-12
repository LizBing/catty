# Project status

**As of:** 2026-07-12  
**Baseline verified through:** `89037c4`  
**Branch:** `main`  
**Milestone:** R1 complete and hardened  
**Active workstream:** None — Codex is preparing R2 architecture gates for LizBing's review

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
| `Integer/Long.toString`, `Double.parseDouble`, representative `HashMap` | R2 prerequisite | `DecimalDigits` reaches `jdk.internal.misc.Unsafe` |
| Java concurrency and monitors | R2 | Thread lifecycle, monitor semantics, JMM guarantees |
| `invokedynamic`, reflection, annotations | R3 | Dynamic metadata and call-site semantics |
| Broad Java I/O and networking | R4 | Native/runtime integration design |
| Broad real-program AOT coverage | R5 | Instance methods, exceptions, dynamic calls |

“18/18 smoke tests pass” means the selected R1 corpus passes; it is not a claim
that arbitrary `java.base` applications are supported.

## Decisions required before R2 implementation

1. Supersede or revise the direction proposed by ADR-0011. Catty may implement
   JMM-observable semantics with Go synchronization primitives; adopting the Go
   memory model as externally visible behavior conflicts with the stated goal
   of preserving JRE semantics.
2. Define strict handling for unresolved native methods. Zero/null discovery
   stubs must not silently become the default compatibility behavior.
3. Specify the minimum Unsafe semantic surface by behavior—offset identity,
   CAS, volatile access, fences, array layout, park/unpark—not by method count.
4. Define representative concurrency and `java.base` acceptance programs before
   implementation begins.

These items should become the first R2 workstream contract and, where needed,
superseding ADRs. No R2 implementation owner is currently assigned.

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

Codex drafts the R2 architecture-gates workstream for LizBing's review. No R2
implementation begins before that contract resolves the four open decisions,
reaches **Ready**, and names one implementation owner.
