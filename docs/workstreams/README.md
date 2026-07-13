# Workstreams

Workstreams use the model-neutral protocol in
[`../COLLABORATION.md`](../COLLABORATION.md) and the
[`TEMPLATE.md`](./TEMPLATE.md).

**Active workstream:** [`r2-initialization-slice`](./r2-initialization-slice.md)
— Accepted; implementation has not started.
Ready for owner review.

Proposed ADRs from this workstream:
- [ADR-0025](../adr/0025-class-initialization-state-machine.md) — Java 25 class/interface
  initialization state machine (Accepted)
- [ADR-0027](../adr/0027-string-utf16-representation.md) — UTF-16 String kernel backing
  (Accepted)

Proposed implementation contract:
- [`r2-initialization-slice`](./r2-initialization-slice.md) — first R2 implementation slice
  (class-init state machine, owner review). Accepted; not yet started.
- [`r2-string-utf16-slice`](./r2-string-utf16-slice.md) — native String UTF-16 slice
  (owner review). Proposed; not yet accepted.

R2 research is authorized after the governance baseline is integrated. It is a
research-only workstream and does not authorize production runtime changes.
Earlier R2 experiments remain reachable on archived/history branches and are
not current project state.
