# K2 owner-fix provisional handoff

**Date:** 2026-07-21  
**Branch:** `codex/r3-runtime-identity-definition-slice`  
**Parent HEAD:** `cc1cf3a801090e746cb7181a71fa38e88464bcce`  
**State:** dirty working tree; no commit, merge, push, or publish authorization

This record supersedes the conclusions in `handoff-r3.md` and
`gates-r3.txt`. Those files are retained as prior-round evidence, but their
statements that independent top-level mutual delegation is outside K2 and that
no R3 harness exists are incorrect.

## Blocking defects closed

- Cross-loader recursion now shares an explicit load context, and a wait graph
  detects cycles between independent top-level load contexts. The deterministic
  two-owner mutual-delegation test terminates with typed circularity failures.
- `ClassLoader.DefineClassResult` is the production typed definition entry,
  separate from lookup/initiation. It preserves existing definitions, rejects
  delegated-cache replacement and pre-bound definition candidates, waits for
  in-flight definitions, and permits retry after failed lookup/definition.
- Class mirrors resolve `java/lang/Class` through the defining loader. Reference
  arrays inherit the component's defining loader; primitive and void mirrors use
  the set-once VM bootstrap loader. Build-time transpilation no longer mutates
  that process-global runtime lifecycle.
- Interpreter, IR, and `<clinit>` tests execute the real exception loops and
  prove that catch-type resolution failure replaces the throwable and escapes
  later handlers in the same frame.
- Nested malformed array descriptors, including `[[V`, `[[X`, `[[II`, and
  `[[Lfoo`, return typed format failures before component lookup.

## Current working-tree gates

| Gate | Result |
|---|---|
| `gofmt -l` on changed Go files | Pass |
| `git diff --check` | Pass |
| `go vet ./...` | Pass |
| `go test -count=1 ./...` | Pass |
| `go test -race -count=1 ./...` | Pass |
| `bash tests/run.sh` | Pass, 10/10 |
| focused `go test -race -count=10 ./classloader ./interpreter ./rtda ./transpile` | Pass |
| stress `go test -race -count=100 ./classloader ./interpreter` | Pass |
| R3 24-row baseline harness | Pass; Interpreter 24/24 `EXIT(1)`, IR 24/24 `EXIT(1)`, AOT 24/24 `NO-BUILD` |

The commands above use `GOCACHE=/private/tmp/catty-go-cache` where needed so
the repository remains the only persistent write target. The R3 run used
`R3_RESULTS_DIR=/private/tmp/catty-k2-owner-fix-r3`; it did not mutate frozen
historical evidence.

## Gate that is intentionally not claimed

The formal R2 concurrency candidate harness has **not** been run for this dirty
tree. `docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh`
accepts a commit and creates a detached worktree, so it cannot authenticate
uncommitted contents. Its required 19-row 1x/100x result remains pending a
candidate commit explicitly authorized by the Owner.

Accordingly, this is implementation evidence, not immutable candidate evidence,
and K2 must not be marked Done from this file alone.
