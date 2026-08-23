# gcopt — Go compiler optimization-behavior probes (P-0006 / R-0003)

Independent Go module (`module gcopt`) that answers: which handwritten Go
idioms keep working (or stop working) inside the code the M3 AOT emitter will
generate? It lives under `docs/research/assets/gcopt/` with its own `go.mod`,
so the root build ignores it.

Toolchain assumed: **go 1.26.x, GOARCH=arm64** (Apple Silicon). Evidence is the
deterministic `-gcflags=-m/-m -m` inliner+escape report, the
`-d=ssa/check_bce/debug=1` bounds-check report, disassembly, and `-benchmem`
benchmarks.

## Layout

| Path | Question |
|---|---|
| `q1_devirt/` | Interface devirtualization + inliner cost budget |
| `q2_closure/` | Closure & method-value escape cost |
| `q3_bigfunc/` | Large straight-line function scaling |
| `q4_concat/` | String concatenation chain vs []byte append |
| `q5_dispatch/` | Type switch vs assertion chain vs int-tag switch |
| `q6_bce/` | Bounds-check elimination loop shapes |
| `gen/` | Emits `q1_devirt/budget_gen.go` and `q3_bigfunc/big_gen.go` |
| `results/` | Raw compiler + benchmark output cited by R-0003 |

## Reproduce

```sh
cd docs/research/assets/gcopt

# regenerate the repetitive probe files
go run ./gen

# build + vet
go build ./... && go vet ./...

# deterministic compiler reports (the primary evidence)
go build -gcflags='-m -m' ./q1_devirt          # inliner budget + devirtualization
go build -gcflags='-m'    ./q2_closure          # closure escape analysis
go build -gcflags='-m -m' ./q3_bigfunc          # big-function inline refusal
go build -gcflags='-d=ssa/check_bce/debug=1' ./q6_bce   # bounds-check report
go build -gcflags='-S'    ./q1_devirt           # devirtualization disassembly
go build -gcflags='-S'    ./q3_bigfunc          # big-function code-size scaling
go test -c -gcflags='-S' -o /dev/null ./q4_concat  # concatstringN vs concatstrings
go test -c -gcflags='-S' -o /dev/null ./q5_dispatch # type-switch jump table

# benchmark timing + allocation counters
go test -run '^$' -bench . -benchmem -count=10 ./...

# compare two runs with benchstat (golang.org/x/perf/cmd/benchstat)
```

## Noise notes

Same Apple Silicon caveats as `assets/objmodel`: heterogeneous cores and active
frequency scaling give ~±15–25 % variance on 3–5 ns samples. Conclusions lean
on the deterministic compiler reports plus allocation counters (which are
exact), and treat timing differences under ~5 % as noise. The inliner budget,
escape, and BCE results are exact and reproducible regardless of thermal state.
