# Project status

**As of:** 2026-07-13
**Stable baseline:** R2 initialization and bounded UTF-16 String slices complete
**Baseline commit:** `8171361` (integration; String candidate `00327d6`, evidence `9008b00`)
**Active workstream:** None
**Current phase:** R2 runtime-semantics planning

This is the single model-neutral current-state entry. Strategy lives in
[`ROADMAP.md`](./ROADMAP.md); decisions live in [`adr/`](./adr/); scoped work
lives in [`workstreams/`](./workstreams/).

## Verified capability

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
- R2 initialization: bounded Java 25 single-execution-context class/interface
  initialization at `new`, resolved `getstatic`/`putstatic`, and resolved
  `invokestatic`; 13/13 differential fixtures match in Interpreter and IR.
  AOT supports the constant-field path and explicitly rejects the remaining
  tested initialization paths pending cross-engine exception propagation.
- R2 String: immutable UTF-16 code-unit backing for the bounded synthetic/native
  String surface. All eight differential fixtures match Temurin 25 in Interpreter
  and IR; AOT supports five fixtures and explicitly reports three as Not implemented.

## Governance-reset validation

Revalidated locally on 2026-07-13:

- `go vet ./...` — Pass
- `go test ./...` — Pass
- `go test -race ./...` — Pass
- `bash tests/run.sh` — Pass, 10/10 fixtures

## Explicit boundary

catty does not claim Java concurrency, monitors, Unsafe, broad reflection,
`invokedynamic`, broad I/O/networking, arbitrary `java.base` application
compatibility, cross-Java-thread initialization behavior, cross-engine AOT
exception propagation, or a complete Java String API. `Integer/Long.toString`,
`Double.parseDouble`, and representative `HashMap` behavior remain blocked by
unresolved runtime/library dependencies.

## Decision state

ADRs 0016–0027 (excluding unused 0026) are Accepted. ADRs 0001–0007 and 0014–0015 are superseded;
ADRs 0008–0013 are withdrawn. ADR-0017 fixes Java 25 as the supported-capability
semantic baseline; ADR-0016 fixes AOT as the primary product path with a
permanent interpreter fallback. ADR-0025 is implemented by the completed,
bounded class/interface-initialization workstream; ADR-0027 is implemented by the
completed bounded UTF-16 String workstream. Bootstrap capability
mapping, Thread/monitor/JMM, Unsafe, and allocation remain deferred.

## Next action

Select or accept the next bounded R2 workstream.
