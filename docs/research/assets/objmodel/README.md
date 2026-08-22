# objbench — object-model microbenchmark prototypes (P-0002)

Independent Go module (`module objbench`) that compares the two object
representation candidates from ADR-0004. It lives under
`docs/research/assets/objmodel/` with its own `go.mod`, so the root build
ignores it.

## Layout

| Path | Contents |
|---|---|
| `common/` | Shared seed value and species multipliers so both candidates do identical arithmetic. |
| `hybrid/` | Candidate A (mixed): concrete structs + embedding inheritance + interface promotion. |
| `uniform/` | Candidate B (unified): `*Object` header + fixed slot array + type-assertion dispatch. |
| `results/` | Raw `go test -bench` output for the runs cited in R-0001. |

## Inheritance model

- Depth 3: `Base -> Mid -> Leaf` (`hybrid` uses `D3Base -> D3Mid -> D3Leaf`).
- Depth 6: six levels (`D6A -> ... -> D6Leaf`).
- Eight species (`A`..`H`, multipliers 3,5,7,11,13,17,19,23) give the
  megamorphic dispatch a realistic number of concrete targets.
- Candidate B stores fields in a fixed `[16]int64` slot array so both
  candidates are a single allocation per object; a slice would unfairly add a
  second allocation to B. Field slot `i` maps to the field at level `i`.

## Benchmark groups

1. Monomorphic call (concrete direct call vs assert-then-call)
2. Bimorphic + megamorphic call (itab dispatch vs per-call type switch)
3. Field read / write (struct selector vs slot indexing)
4. Downcast (interface type assertion vs class-tag comparison)

Each group is measured at inheritance depth 3 and depth 6.

## Reproduce

```sh
cd docs/research/assets/objmodel

# build + vet
go build ./... && go vet ./...

# single clean run with allocation counters (used for the allocs observation)
go test ./hybrid ./uniform -run '^$' -bench . -benchmem -count=10

# interleaved run (primary timing data in R-0001); alternates packages so both
# candidates share thermal/frequency state on Apple Silicon
for round in $(seq 1 15); do
  for pkg in hybrid uniform; do
    go test -run '^$' -bench . -benchtime=400ms ./"$pkg"
  done
done

# escape-analysis summary
go build -gcflags='-m' ./hybrid ./uniform
```

## Noise notes

Apple M1 has heterogeneous (performance/efficiency) cores plus active
frequency scaling. Individual 3–5 ns samples vary by roughly ±15–25 %;
medians over many interleaved samples are stable to about ±5–7 %. Conclusions
in R-0001 therefore lean on medians plus deterministic `-gcflags=-m` and
disassembly evidence, and treat differences under ~5 % as within noise.
