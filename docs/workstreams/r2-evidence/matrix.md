# R2 differential fixture matrix

**Baseline commit:** `ecb086e`
**Reference:** Temurin 25.0.3
**catty mode:** pure-synthetic (`-no-boot`)
**Harness:** `docs/workstreams/r2-evidence/run-r2-diff.sh`

`match` means byte-identical combined output and equal exit code. `MISMATCH` differs;
`NO-BUILD` means AOT produced no executable.

| # | Fixture | Category | Interpreter | IR | AOT | Finding |
|---|---|---|---|---|---|---|
| 1 | `ClinitOrder` | init order | match | match | NO-BUILD | superclass-before-subclass baseline |
| 2 | `ClinitThrows` | init failure | MISMATCH | MISMATCH | MISMATCH | raw exception; no EIIE/erroneous behavior |
| 3 | `ConstantFieldNoInit` | constant field | match | match | match | `ConstantValue` read does not initialize |
| 4 | `GetstaticOwner` | static declarer | MISMATCH | MISMATCH | MISMATCH | initializes `Sub`, then wrong-storage crash |
| 5 | `InterfaceDefaultInit` | class predecessors | MISMATCH | MISMATCH | MISMATCH | default-bearing interface not initialized before implementing class |
| 6 | `InterfaceNoInitOnImpl` | interface exclusion | match | match | NO-BUILD | interface without default method remains uninitialized |
| 7 | `InvokeStaticInit` | init trigger | match | match | match | `invokestatic` baseline |
| 8 | `RecursiveInitialization` | recursive request | match | match | match | in-progress read returns default value; no second `<clinit>` |
| 9 | `SuperclassInitializationFailure` | predecessor failure | MISMATCH | MISMATCH | NO-BUILD | raw superclass exception; no erroneous subclass behavior |
| 10 | `HashDivergence` | String UTF-16 | MISMATCH | MISMATCH | MISMATCH | hash divergence |
| 11 | `LoneSurrogate` | String UTF-16 | MISMATCH | MISMATCH | NO-BUILD | `String(char[])` is absent |
| 12 | `StringBounds` | String bounds | MISMATCH | MISMATCH | NO-BUILD | invalid accesses silently return instead of throwing |
| 13 | `StringCharArrayRoundTrip` | String char[] copy | MISMATCH | MISMATCH | NO-BUILD | `String(char[])` is absent |
| 14 | `StringSubstringUnits` | String UTF-16 substring | MISMATCH | MISMATCH | MISMATCH | surrogate units become replacement characters |
| 15 | `SupplementaryChar` | String UTF-16 | MISMATCH | MISMATCH | MISMATCH | code-unit length/charAt divergence |
| 16 | `ReachUnsafe` | bootstrap boundary | MISMATCH | MISMATCH | NO-BUILD | synthetic-mode `java/lang/Integer` boundary |

The 9 initialization fixtures are rows 1–9; the 6 String fixtures are rows 10–15. Exact
reference and engine output is retained in `run-r2-results.txt`.

## AOT classifications

- `NO-BUILD` is currently a measured harness result, not yet an ADR-0016 capability label.
- Semantic mismatches shared by all engines point to shared runtime behavior.
- A future Accepted implementation workstream must convert every unsupported acceptance
  path to an explicit `Fallback` or `Not implemented` classification.
