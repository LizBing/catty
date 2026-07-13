# ADR-0008: AOT-first architecture (interpreter is the dev tier, not production)

- **Status:** Withdrawn
- **Date:** 2026-07-12
- **Withdrawn:** 2026-07-13 — a single `fib(35)` result cannot justify an
  AOT-only production contract. ADR-0016 subsequently establishes AOT as the
  primary product path while retaining a permanent interpreter fallback.

## Context

catty's AOT path (`catty build` → `go build`) runs `fib(35)` in 44ms, on par
with HotSpot JIT. The interpreter runs it in 4.5s — 100× slower. The interpreter
exists for development/debugging and as a fallback for dynamically-loaded classes
that can't be AOT'd.

Traditional JVMs use interpret-first → JIT-on-hot (tiered compilation with
warmup). catty flips this: **AOT is the production path, the interpreter is the
development tool.** There is no JIT warmup, no safepoint polling (Go's GC is
concurrent), no on-stack replacement, no deoptimization.

## Decision

catty's production execution path is **AOT-only**.

- `catty build` is the production entry point.
- `catty -cp . Main` (interpreter) is the development/debugging entry point.
- Runtime dynamic class loading falls back to the interpreter (the fallback tier).
- **No JIT is implemented** — no warmup, no safepoints, no tiered compilation.

## Consequences

**Positive**
- No JIT warmup (first request is full speed).
- No safepoint checks (Go GC is concurrent).
- Simpler codebase (no profile collection, no deoptimization).
- Better optimization opportunities (Go compiler does whole-program analysis; JIT
  only sees hotspots).

**Negative**
- Dynamically-loaded classes (`Class.forName`) run interpreted (slow).
- No runtime profile-guided optimization (JIT inlines based on actual call
  frequency).
- AOT requires closed-world analysis (which classes might be loaded?), limiting
  dynamic scenarios.
