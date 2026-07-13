# Workstreams

Workstreams use the model-neutral protocol in
[`../COLLABORATION.md`](../COLLABORATION.md) and the
[`TEMPLATE.md`](./TEMPLATE.md).

**Active workstream:** Accepted
[`r2-thread-monitor-foundation-slice`](./r2-thread-monitor-foundation-slice.md).

The R2 initialization slice is Done.

Proposed ADRs from this workstream:
- [ADR-0025](../adr/0025-class-initialization-state-machine.md) — Java 25 class/interface
  initialization state machine (Accepted)
- [ADR-0027](../adr/0027-string-utf16-representation.md) — UTF-16 String kernel backing
  (Accepted)

Implementation contracts:
- [`r2-initialization-slice`](./r2-initialization-slice.md) — first R2 implementation slice
  (class-init state machine, owner review). Done.
- [`r2-string-utf16-slice`](./r2-string-utf16-slice.md) — native String UTF-16 slice
  (owner review). Done; eight-fixture engine matrix accepted and integrated.
- [`r2-concurrency-semantics-research`](./r2-concurrency-semantics-research.md) — Done
  research-only contract for Thread identity/lifecycle, monitors/wait sets,
  cross-thread initialization, and the minimum memory-ordering boundary. It does not
  authorize production concurrency.
- [`r2-thread-monitor-foundation-slice`](./r2-thread-monitor-foundation-slice.md) — Accepted
  implementation contract produced by the concurrency research. ADR-0028 through ADR-0030
  and the contract were accepted by the Owner on 2026-07-14; its acceptance anchor is
  required before production work.

The original R2 runtime-semantics research workstream is Done. The concurrency
research above is Done. Its successor contract was accepted by the Owner on
2026-07-14 and authorizes only its frozen production scope after its acceptance
anchor and implementation preflight. Earlier R2 experiments remain
reachable on archived/history branches and are not current project state.
