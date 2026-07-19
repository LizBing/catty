# R3 fixed research fixtures

This directory contains exactly the 24 Java 25 sources frozen by the Accepted
`r3-reflection-dynamic-research` contract. `manifest.sha256` records their
source hashes. Verify it from this directory with:

```bash
shasum -a 256 -c manifest.sha256
```

`BootstrapFailureOnce.java` intentionally declares a package-private
`BootstrapFailureTarget`. The harness compiles both classes and removes
`BootstrapFailureTarget.class` before execution. The method reference then
forces resolution of an unavailable MethodHandle bootstrap argument at one
InvokeDynamic call site; two executions record the repeated failure behavior.

The baseline harness requires an explicit, previously nonexistent output
directory. It uses `CATTY_BOOT` when it points at an extracted `java.base`, or
extracts the JDK 25 modules image into a temporary directory:

```bash
R3_RESULTS_DIR=docs/workstreams/r3-reflection-dynamic-evidence/baseline-6cf3636 \
  bash docs/workstreams/r3-reflection-dynamic-fixtures/run-r3-baseline.sh
```

The harness is descriptive. A catty mismatch, panic, timeout, or AOT NO-BUILD
is baseline evidence, not a capability claim. Missing fixtures or rows and a
failed Temurin reference fail closed.
