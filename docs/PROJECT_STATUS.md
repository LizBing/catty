# Project status

**As of:** 2026-07-13
**Stable baseline:** R1 complete and hardened
**Baseline commit:** `5720147`
**Active workstream:** None
**Current phase:** Governance reset before post-R1 planning

This is the single model-neutral current-state entry. Strategy lives in
[`ROADMAP.md`](./ROADMAP.md); decisions live in [`adr/`](./adr/); scoped work
lives in [`workstreams/`](./workstreams/).

## Verified R1 capability

- Interpreter: approximately 145 opcodes, exceptions, interface dispatch,
  multidimensional arrays, and class initialization.
- Class loading: provider chain plus real `java.base` auto-detection through a
  JDK-extracted image.
- Native/bootstrap layer: six irreducible synthetic bootstrap classes,
  additional synthetic fallbacks, and approximately 40 native registrations.
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
`invokedynamic`, broad I/O/networking, or arbitrary `java.base` application
compatibility. `Integer/Long.toString`, `Double.parseDouble`, and representative
`HashMap` behavior remain blocked by unresolved runtime/library dependencies.

## Decision state

ADRs 0001–0007 and 0014–0015 are Accepted. ADRs 0008–0013 are Proposed and do
not authorize implementation. In particular, Thread mapping, Java memory
semantics, Unsafe boundaries, AOT production scope, and direct Go runtime
integration require renewed planning and, where appropriate, Accepted ADRs.

## Next action

Propose the next research or implementation workstream under the Accepted
collaboration protocol. Do not begin non-trivial post-R1 implementation before
the Project Owner accepts its contract.
