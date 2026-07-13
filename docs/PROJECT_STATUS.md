# Project status

**As of:** 2026-07-13
**Stable baseline:** R1 complete and hardened
**Baseline commit:** `5720147`
**Active workstream:** [`r2-runtime-semantics-research`](./workstreams/r2-runtime-semantics-research.md) (Accepted; not started)
**Current phase:** R2 runtime-semantics research

This is the single model-neutral current-state entry. Strategy lives in
[`ROADMAP.md`](./ROADMAP.md); decisions live in [`adr/`](./adr/); scoped work
lives in [`workstreams/`](./workstreams/).

## Verified R1 capability

- Interpreter: approximately 145 opcodes, exceptions, interface dispatch,
  multidimensional arrays, and class initialization.
- Class loading: provider chain plus real `java.base` auto-detection through a
  JDK-extracted image.
- Native/bootstrap layer: the current R1 implementation has six synthetic
  bootstrap classes, additional synthetic fallbacks, and approximately 40
  native registrations; ADR-0022 does not treat that class list as permanent.
- AOT: standalone Go binary path; `fib(35)` recorded at approximately 40–60 ms
  on the development machine.
- Regression baseline: unit tests, three-engine fixture comparison, and the
  real `java.base` smoke path in CI.

## Governance-reset validation

Revalidated locally on 2026-07-13:

- `go vet ./...` — Pass
- `go test ./...` — Pass
- `go test -race ./...` — Pass
- `bash tests/run.sh` — Pass, 10/10 fixtures

## Explicit boundary

R1 does not claim Java concurrency, monitors, Unsafe, broad reflection,
`invokedynamic`, broad I/O/networking, arbitrary `java.base` application
compatibility, full class/interface initialization semantics, or full Java String
UTF-16 behavior. `Integer/Long.toString`, `Double.parseDouble`, and representative
`HashMap` behavior remain blocked by unresolved runtime/library dependencies.

## Decision state

ADRs 0016–0024 are Accepted. ADRs 0001–0007 and 0014–0015 are superseded;
ADRs 0008–0013 are withdrawn. ADR-0017 fixes Java 25 as the supported-capability
semantic baseline; ADR-0016 fixes AOT as the primary product path with a
permanent interpreter fallback. Detailed class/interface initialization,
bootstrap capability mapping, String representation, Thread/monitor/JMM,
Unsafe, and allocation decisions remain R2 research outputs rather than
implementation authority.

## Next action

Integrate the governance baseline, then begin Slice A of the Accepted R2
research workstream. No R2 production implementation is authorized yet.
