# R3 reflection/dynamic research baseline summary

**Acceptance anchor:** `6cf3636`
**Branch at run:** `codex/r3-reflection-dynamic-research`
**Date:** 2026-07-17
**Reference:** Temurin 25.0.3
**Toolchain:** javac/jimage 25.0.3; Go 1.26.5 darwin/arm64

## Command

```bash
R3_RESULTS_DIR=docs/workstreams/r3-reflection-dynamic-evidence/baseline-6cf3636 \
  bash docs/workstreams/r3-reflection-dynamic-fixtures/run-r3-baseline.sh
```

The harness extracted the local JDK 25 modules image to a temporary directory
and placed its `java.base` directory in `CATTY_BOOT` for every catty build/run.
Every process had a 20-second run timeout; AOT builds had 120 seconds.

## Fixture freeze

- Exactly 24 Java sources were present.
- `manifest.sha256` verified all source hashes before execution.
- All 24 compiled with `javac --release 25` and ran successfully on Temurin.
- For `BootstrapFailureOnce`, the harness removed the compiled
  `BootstrapFailureTarget.class` to force failure while resolving the method
  reference's bootstrap MethodHandle argument. Both attempts produced
  `NoClassDefFoundError`; Temurin reported distinct Throwable references.

## Engine result

| Engine | Result |
|---|---|
| Temurin 25.0.3 | Pass — 24/24 references exited 0 |
| Interpreter | Baseline complete — 0 Match, 24 Exit(1) |
| IR | Baseline complete — 0 Match, 24 Exit(1) |
| AOT | Baseline complete — 24/24 NO-BUILD |
| Rows/timeouts | Pass — 24/24 rows, no timeout or omitted engine |

Baseline failure is expected research data and is not an acceptance failure.
No catty row is classified Supported or Fallback.

## Dominant failures

- Interpreter reflection rows reach missing synthetic Class/member methods and
  normally throw `NoSuchMethodError`; IR frequently resolves nil and panics.
- `ClassQueries` prints `false` for superclass mirror identity before failing,
  exposing a canonical-mirror bypass in `classGetSuperclass`.
- Interpreter rejects opcode `0xba` explicitly; IR reports InvokeDynamic not
  supported by lowering.
- AOT no-build diagnostics are generic; they are not yet the precise R3
  reachability rejection required by a future candidate.

Full combined stdout/stderr and exit status for every engine/fixture are in
`results.txt`.
