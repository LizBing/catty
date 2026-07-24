# Engine/native execution-context migration evidence

Date: 2026-07-25

Slice D migrated engine, native helper, launcher, and AOT bridge APIs toward
execution-context terminology while preserving the temporary compatibility
surface.

## Production migration

- Interpreter bytecode loop, IR loop, invoke helpers, class-initialization
  helpers, bridge `RunMethod`, array helpers, and resolution helpers now accept
  `*rtda.ExecutionContext` and use `context` naming.
- Generic native helpers use `Frame.Context()` instead of the compatibility
  `Frame.Thread()` accessor.
- The launcher creates the main execution context with
  `rtda.NewExecutionContext`.
- The AOT runtime bridge package-global is now `context *rtda.ExecutionContext`.
  It exposes `ExecutionContext()` as the clear accessor.
- Java Thread native facade code keeps Java Thread terminology for facade
  operations, but uses `Frame.Context()` for the calling execution context and
  `JavaThreadState` for target facade state.

## Remaining compatibility surface

`rg -n "\*rtda\.Thread|\*Thread|type Thread = ExecutionContext|func NewThread|func \(f \*Frame\) Thread|func Thread\(\)" . -g '*.go'`
returns only:

- `rtda/thread.go`: temporary `type Thread = ExecutionContext`
- `rtda/thread.go`: compatibility constructor `NewThread`
- `rtda/frame.go`: compatibility accessor `Frame.Thread`
- `runtime/runtime.go`: compatibility accessor `Thread`
- `rtda/frame_test.go`: focused compatibility assertion

The compatibility surface remains until final regression closure authorizes
removal.

## Checks

- `GOCACHE=/private/tmp/catty-review-go-cache go test ./rtda ./interpreter ./native ./launch ./runtime` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go test ./...` — Pass
- `GOCACHE=/private/tmp/catty-review-go-cache go vet ./...` — Pass
