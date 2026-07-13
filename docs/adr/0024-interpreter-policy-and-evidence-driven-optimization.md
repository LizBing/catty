# ADR-0024: Interpreter simplicity and evidence-driven optimization

- **Status:** Accepted
- **Date:** 2026-07-13
- **Supersedes:** [ADR-0002](./0002-switch-dispatch.md), [ADR-0006](./0006-predecode-no-speedup.md)

## Context

R1 measured a switch-dispatched tree-walker faster than its first IR executor
on one recursive benchmark. Those measurements are useful history, but neither
a Go dispatch technique nor one predecode result should freeze catty's
long-term interpreter architecture.

## Decision

The interpreter prioritizes semantic coverage, safe fallback, diagnosis, and
maintainability. The current switch dispatcher and current IR executor are
allowed implementations, not long-term architecture constraints. Interpreter
optimization is justified only by measured benefit on representative workloads
and must not undermine ADR-0016's AOT-primary path or cross-engine parity.

No specific dispatch technique, frame layout, predecode scheme, or IR executor
performance conclusion is permanent. A workstream may replace any of them
after defining semantic gates and reporting relevant performance evidence.

## Consequences

- The R1 `BenchFib` predecode result remains historical evidence rather than a
  prohibition on future interpreter optimization.
- AOT remains the principal performance path under ADR-0016; interpreter work
  is scoped by fallback and mixed-engine needs.
