# excbench — exception propagation mechanism microbenchmark (P-0006 / R-0002)

Independent Go module (`module excbench`) comparing the two exception-propagation
candidates behind ADR-0005 / ADR-0009. It lives under
`docs/research/assets/excbench/` with its own `go.mod`, so the root build
ignores it.

## Candidates

| Package | Candidate | Throw | Propagate | Catch |
|---|---|---|---|---|
| `panicpkg/` | P (panic/recover) | `panic(sentinel)` | deferred `recover` + re-panic on mismatch | deferred `recover` + tag match |
| `flagpkg/` | F (flag return) | `return 0, sentinel` | `if err != nil { return err }` | `if err != nil { handle }` |

`common/` holds the sentinel (`*JException`, only its class tag is modeled)
and the identical `Work` payload so both candidates do the same arithmetic.

## The modeled chain

Four layers `a -> b -> c -> d`. `d` does the work and, on the throwing paths,
raises at every iteration (no RNG, so the RNG cost cannot mask the mechanism).
The three scenarios from the plan:

1. **Normal path** — zero exceptions, steady-state loop (the main P-vs-F
   battlefield).
2. **Throw-catch (shallow)** — high-frequency raise, caught one boundary up.
3. **Deep propagation** — exception crosses all four frames.

Candidate P is measured across its cost spectrum: with every layer owning a
handler (worst case), with no layer owning a handler (the realistic majority,
true zero cost), and for deep propagation both the free-unwind and the
re-panic-per-layer variants.

## Reproduce

```sh
cd docs/research/assets/excbench

go build ./... && go vet ./...

# correctness guards (re-panic / swallow logic)
go test ./...

# primary timing + allocation data cited in R-0002
go test ./panicpkg ./flagpkg -run '^$' -bench . -benchmem -count=10

# benchstat comparison (installed locally), if available
go test ./panicpkg -run '^$' -bench . -benchmem -count=10 > /tmp/p.txt
go test ./flagpkg  -run '^$' -bench . -benchmem -count=10 > /tmp/f.txt
benchstat /tmp/p.txt /tmp/f.txt
```

Raw output is archived under `results/`.

## Methodological notes

- Every chain method is `//go:noinline` to model the realistic case where AOT
  methods are too large to inline; the mechanism cost is per-boundary and the
  emitter cannot rely on inlining to hide it.
- The exception object is pre-allocated and reused, so the measured cost is the
  mechanism, not per-throw allocation (allocation is identical across
  candidates and orthogonal). `-benchmem` still reports what the mechanism
  itself allocates.
- Hardware: Apple M1, darwin/arm64, go1.26.x. Apple Silicon has heterogeneous
  cores plus frequency scaling, so individual ns/op samples vary roughly
  ±15–25 %; medians over many interleaved samples are stable to about ±5–7 %.
  Conclusions treat differences under ~5 % as within noise.
