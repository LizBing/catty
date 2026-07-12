# R2 concurrency and Unsafe test plan

## Purpose

R2 adds schedule-sensitive behavior that cannot be validated by the existing
stdout-only fixture loop alone. This plan defines deterministic semantic tests,
bounded CI stress, longer local stress, HotSpot differential evidence, and
performance comparisons for Strict, Go-native research, and Hybrid backends.

The plan does not claim that every engine supports concurrency on day one.
Each block declares its engine matrix; an unsupported engine fails explicitly
or is marked not implemented in capability metadata, never silently counted as
passing.

## Harness layout

Proposed repository shape:

```text
tests/
├── fixtures/                 existing R1 deterministic stdout tests
├── research/                 expected-failure and dependency probes
├── concurrency/
│   ├── deterministic/        protocol-controlled semantic cases
│   ├── litmus/               outcome-set tests derived from JMM/jcstress
│   └── stress/               repeated schedule-sensitive cases
└── run-r2.sh                 bounded differential runner with timeouts
```

Every R2 case records:

- case name and semantic clause;
- engine and backend;
- pinned JDK reference version;
- timeout and iteration count;
- allowed, forbidden, and interesting outcomes;
- seed when randomized;
- stdout, stderr, and thread-state snapshot on failure.

## Deterministic programs

| Program | Required assertions |
|---|---|
| StrictNativeProbe | unresolved declaration throws catchable `UnsatisfiedLinkError`; message identifies signature |
| StartJoinPublication | parent-before-start visible to child; child-before-termination visible after join; second start throws |
| MonitorCounter | mutual exclusion, reentrancy, instance/static synchronized, abrupt-exit release, non-owner errors |
| WaitNotifyProtocol | wait releases/restores recursion; notify one; notifyAll all; timeout; ownership error |
| InterruptProtocol | flag set/observe/clear; wait/sleep/join throw and clear; notify/interrupt race legal outcomes |
| VolatilePublication | payload visible after volatile publication; volatile total-order litmus |
| AtomicCAS | reference/int/long CAS; compare-exchange; linearizable counter |
| ParkPermit | unpark-before-park, one permit, no permit accumulation, timeout, interrupt return |
| ClassInitRace | exactly once, recursive request, publication, erroneous state and subsequent failure |
| NonDaemonLifetime | non-daemon keeps VM alive; daemon-only does not; main termination ordering |
| UnsafeJavaBaseSmoke | Integer/Long decimal conversion plus only caller-verified Unsafe paths |
| HashMapBasic | put/get/remove/collision/null without incorrectly requiring serialization Unsafe path |
| DoubleParse | finite, signed zero, NaN/infinity, exponent, invalid input; timeout treated as failure |

Tests use latches/barriers/handshakes implemented by the capability under test.
`Thread.sleep` may create a timeout scenario but is never the sole correctness
coordination mechanism.

## Litmus matrix

Port or adapt 40–60 small tests from published JMM/jcstress categories:

- store buffering and message passing;
- volatile Dekker and IRIW;
- final-field safe/unsafe publication;
- monitor publication and reentrancy;
- start/join/interrupt ordering;
- long/double/reference atomicity and tearing;
- CAS, exchange, acquire/release/opaque/volatile VarHandle modes;
- class initialization;
- Unsafe publication and ordered writes.

For each test:

```text
HotSpot outcomes
Strict Catty outcomes
Go-native research outcomes
Hybrid Catty outcomes
```

A production backend fails if it produces an outcome forbidden by its declared
Java profile. Absence of an interesting weak outcome is not failure when Catty
provides stronger ordering.

## Execution tiers

### Per-commit local/CI

- deterministic programs once per claimed engine;
- bounded litmus set, at least 10,000 actor iterations or a fixed short budget;
- per-case timeout;
- `go test -race ./...` for production runtime code;
- existing R1 regression suite with real java.base.

### Nightly/manual stress

- full litmus matrix;
- at least 1,000,000 actor iterations or 10–30 minutes per category;
- multiple `GOMAXPROCS` values including 1, 2, and host CPU count;
- repeated seeds and forced scheduling (`runtime.Gosched` only as perturbation,
  not semantics);
- deadlock watchdog and state dump.

## Engine matrix

| Capability | Interpreter | IR | AOT/runtime bridge |
|---|---:|---:|---:|
| Strict native failure | Required | Required | Required |
| Thread lifecycle | Required | Required | Required before R2 closes |
| monitor bytecodes | Required | Required | Required when AOT emits them |
| synchronized method | Required | Required | Required before AOT claims method support |
| volatile/shared heap | Required | Required | Required before concurrent AOT |
| Unsafe profiles | Required | Required | Required for bridged/compiled callers |

An engine may remain explicitly unsupported during an intermediate block, but
R2 closure requires either semantic parity or a documented later milestone that
prevents that engine from claiming the capability.

## Commands

Baseline gates:

```sh
gofmt -w <changed-go-files>
go vet ./...
go test ./...
go test -race ./...
bash tests/run.sh
bash tests/run-r2.sh --ci
```

Stress and benchmark examples, finalized with the harness:

```sh
bash tests/run-r2.sh --stress --duration 20m --seed <seed>
go test -run '^$' -bench 'Monitor|Volatile|CAS|Thread|Park' -benchmem ./...
```

## Memory-model study decision report

The Strict/Go-native/Hybrid study records:

- CPU, OS, Go, JDK, Catty commit, engine, and `GOMAXPROCS`;
- raw benchmark samples, median, p95/p99, and variance;
- forbidden/interesting litmus outcomes;
- race detector findings separated into intentional research races and runtime
  implementation defects;
- static/dynamic synchronization profile of representative workloads.

A proposed semantic waiver names the exact lost behavior. It is considered only
at ≥2× affected microbenchmark, ≥15% application throughput, or ≥20% p99
latency improvement, and still requires its own ADR.

## Failure policy

- timeout is failure, not skip;
- unexpected Java exception is failure;
- Go panic/race/deadlock is failure;
- missing capability is an explicit expected-failure entry only before its
  assigned block, and cannot count toward PASS totals;
- expected failures remain in the repository until resolved or reassigned by an
  accepted architecture decision.
