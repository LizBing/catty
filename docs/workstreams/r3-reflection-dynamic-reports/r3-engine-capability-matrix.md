# R3 engine capability matrix

**Acceptance anchor:** `6cf3636`
**Evidence:** `../r3-reflection-dynamic-evidence/baseline-6cf3636/results.txt`
**Reference:** Temurin 25.0.3

## Baseline result

The fail-closed harness produced all 24 rows. Every Temurin reference compiled,
ran within 20 seconds, and exited 0. Catty did not match any fixture:

| Engine | Match | Mismatch | Exit/panic | No-build | Missing/timeout |
|---|---:|---:|---:|---:|---:|
| Temurin 25.0.3 | 24 reference | 0 | 0 | 0 | 0 |
| Interpreter | 0 | 0 | 24 | 0 | 0 |
| IR | 0 | 0 | 24 | 0 | 0 |
| AOT | 0 | 0 | 0 | 24 | 0 |

This is a descriptive baseline. AOT `NO-BUILD` is not yet an accepted
Not-implemented gate because the diagnostic is generic rather than an R3
reachability classification.

## Fixed rows

| Family | Fixture | Interpreter | IR | AOT |
|---|---|---|---|---|
| Class/type | ClassIdentity | Exit 1 | Exit 1 | No-build |
| Class/type | ClassForNameInit | Exit 1 | Exit 1 | No-build |
| Class/type | ClassQueries | Exit 1 | Exit 1 | No-build |
| Class/type | DeclaredMembers | Exit 1 | Exit 1 | No-build |
| Class/type | PrimitiveAndArrayClass | Exit 1 | Exit 1 | No-build |
| Class/type | MissingClass | Exit 1 | Exit 1 | No-build |
| Member | MethodInvoke | Exit 1 | Exit 1 | No-build |
| Member | ConstructorInvoke | Exit 1 | Exit 1 | No-build |
| Member | FieldGetSet | Exit 1 | Exit 1 | No-build |
| Member | StaticReflectiveInit | Exit 1 | Exit 1 | No-build |
| Member | ReflectiveConversions | Exit 1 | Exit 1 | No-build |
| Member | ReflectiveFailures | Exit 1 | Exit 1 | No-build |
| Annotation | ClassAnnotation | Exit 1 | Exit 1 | No-build |
| Annotation | MemberAnnotation | Exit 1 | Exit 1 | No-build |
| Annotation | AnnotationDefaults | Exit 1 | Exit 1 | No-build |
| Annotation | InheritedRepeatableAnnotation | Exit 1 | Exit 1 | No-build |
| InvokeDynamic | StringConcatIndy | Exit 1 | Exit 1 | No-build |
| InvokeDynamic | StatelessLambda | Exit 1 | Exit 1 | No-build |
| InvokeDynamic | CapturingLambda | Exit 1 | Exit 1 | No-build |
| InvokeDynamic | MethodReference | Exit 1 | Exit 1 | No-build |
| InvokeDynamic | BootstrapFailureOnce | Exit 1 | Exit 1 | No-build |
| Proxy | ProxyDispatch | Exit 1 | Exit 1 | No-build |
| Proxy | ProxyObjectMethods | Exit 1 | Exit 1 | No-build |
| Proxy | ProxyFailureAndDefault | Exit 1 | Exit 1 | No-build |

## Failure clusters

| Cluster | Interpreter evidence | IR evidence | Architectural meaning |
|---|---|---|---|
| missing Class/member facade | Java `NoSuchMethodError` on methods such as `Class.forName` and `getClassLoader` | nil resolved target or nil-pointer panic | shared resolution must return Java errors; IR cannot assume every resolved method exists |
| Class identity | `ClassQueries` prints `false` for canonical superclass identity before later failure | fails later in the same path | all Class-producing natives must use the canonical mirror service |
| annotation/reflection bytecode | stops at missing facade/runtime method | nil-pointer panic | metadata parsing alone is insufficient; facade and exception transport are prerequisites |
| invokedynamic | explicit `opcode 0xba ... not implemented` | explicit `lowering: invokedynamic not yet supported` | shared CallSite service and a typed IR operation are required |
| proxy | missing ClassLoader/Proxy/reflection surface | nil-pointer panic | defining-loader/generated-class and handler dispatch are required together |
| AOT | generic unsupported-opcode/transpile failure or build failure | n/a | needs explicit R3 reachability rejection before it can be classified Not implemented |

## Proposed R3 completion state

Phase R5 already owns AOT InvokeDynamic coverage, and current AOT exception and
Thread-context bridges cannot preserve the R3 contract. The bounded R3
implementation should therefore target:

| Family | Interpreter | IR | AOT in R3 |
|---|---|---|---|
| Class/type reflection | Required Supported | Required Supported | Not implemented / explicit build rejection |
| Member access/invocation | Required Supported | Required Supported | Not implemented / explicit build rejection |
| Runtime annotations | Required Supported | Required Supported | Not implemented / explicit build rejection |
| InvokeDynamic + concat/lambda | Required Supported | Required Supported | Not implemented / explicit build rejection (R5 owns expansion) |
| Dynamic proxies | Required Supported | Required Supported | Not implemented / explicit build rejection |

An implementation slice may support only its frozen subset, but every candidate
harness must still classify all relevant rows. Rows outside that slice remain
explicit Not implemented; they may not disappear from the matrix.

## Candidate acceptance policy

- Supported means combined stdout, stderr, and exit status match Temurin 25.0.3.
- Fallback means an accepted, intentionally selected engine transition executes
  the same semantics and is reported as such. R3 does not assume the current
  AOT bridge is a valid fallback.
- Not implemented means build-time or pre-execution rejection with a precise
  R3 capability diagnostic.
- A Go panic, nil dereference, generic unsupported-opcode crash, timeout,
  built-then-panic binary, or omitted row is Fail.
- Interpreter and IR share runtime services but require separate executions and
  separate result columns.
