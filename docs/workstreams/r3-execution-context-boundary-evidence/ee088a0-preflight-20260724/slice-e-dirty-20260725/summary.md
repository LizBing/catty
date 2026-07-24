# Slice E dirty-tree regression summary

Date: 2026-07-25

This evidence was produced from the current dirty worktree on branch
`codex/r3-execution-context-thread-facade-boundary`. It is useful regression
evidence, but it is not fixed-candidate acceptance evidence because no
candidate commit has been authorized or created.

## Passed dirty-tree checks

- `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go test -race ./...` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache bash tests/run.sh` — Pass, 10/10 fixtures
- `GOCACHE=/private/tmp/catty-review-go-cache R3_RESULTS_DIR=.../r3-baseline bash docs/workstreams/r3-reflection-dynamic-fixtures/run-r3-baseline.sh` — Pass
- `git diff --check` — Pass
- Historical evidence check against `ee088a0` — Pass

## R3 baseline statistics

Machine-counted from `r3-baseline/results.txt`:

- rows: 24/24
- Interpreter MATCH: 0/24
- IR MATCH: 0/24
- AOT NO-BUILD: 24/24

No Java-visible R3 row is newly claimed supported.

## Not completed

The R2 concurrency candidate runner was not run because
`docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh` requires
a fixed candidate commit id and executes a detached worktree at that commit.
Running it against `ee088a0` would exclude the current dirty implementation
changes; running it against the dirty worktree would violate the candidate
evidence protocol.

Required remaining gates after a candidate commit is authorized:

- `bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>`
- `R2_CONCURRENCY_STRESS=100 bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>`

Until those pass, the workstream must remain `In Progress`, not `Ready` or
`Done`.
