# Roadmap

## Completed

### Phase 1 — Interpreter MVP ✅
Switch-dispatched bytecode interpreter running a minimal Java subset
(~140 opcodes), core classes natively in Go, 8 fixtures byte-identical to
`java`. ~3300 LOC.

### A0 — Bytecode IR + stack elimination ✅
`lowering/` package: decode → depth dataflow → vreg assignment. Stack
eliminated into slot-indexed virtual registers. Depth-only; no SSA/phis
(position-stable slots are sound for execution). IR executor (`LoopIR`)
validates the lowering is semantics-preserving. ADR-0006: predecode in the
interpreter does not pay off — AOT is the perf path.

### A1 — AOT emitter (fib natively) ✅
`transpile.Emit`: method IR → Go source. `fib(35)` runs natively in ~44 ms
(~100× interpreter, on par with HotSpot JIT).

### A1.5 — Type tracking ✅
`StackMapTable` parsing + type-aware dataflow. `IRInst.InTypes` carries
operand-stack slot types. No SSA/merge logic (frames pin merges).

### A2.1 — Fresh-per-def, type-aware emitter ✅
Rewrote emitter to fresh-per-def temps (resolves slot-type reuse). Refs,
arrays, typed params/returns. Merge-free gate.

### A2.2 — Invoke bridge ✅
`catty/runtime`: Bootstrap, GetStatic, InvokeVirtual, NewString.
HelloWorld transpiled + run natively. Native targets only.

### A2.3 — Loops ✅
Empty-stack merges (loop heads) need no phi — loop state is in mutable
locals. Gate relaxed. `ArrayOps.sum` (for-loop) AOT-executed.

### A2.4 — Diamonds / phi ✅
Non-empty-stack merges (diamonds) handled via copy-insertion: per-slot
merge temps assigned at predecessor edges, read after the join. `max`
(ternary) AOT-executed. Control-flow complete.

### A2.2b — OOP ✅
`new`/`getfield`/`putfield`/`invokespecial` + interpreted-target bridge
(`interpreter.RunMethod` + Thread.bridgeReturn). `OOPAot` (new + fields +
user invokevirtual) AOT-executed.

### A2.5 — long/float/double ✅
`defTempCat2` (one Go temp for 2 JVM slots). ~50 new opcode cases. All
primitive types covered. `Factorial.fact(10)` (long), `fadd`/`dmul`
(float/double) AOT-executed.

### A2.6 — Edge items ✅
`frem`/`drem` (runtime.FloatMod/DoubleMod), cat-2 merge phi (long/double
at diamond joins), switch (tableswitch/lookupswitch → Go switch + goto).
Emitter covers all opcodes the interpreter supports.

### A4 — `catty build` ✅
`transpile.BuildProgram`: reachability traversal + pass1/2 (emittability +
invokestatic dispatch) + assemble. `catty build` command. HelloWorld +
Fibonacci build to standalone native binaries, output byte-identical to
`java`. The AOT arc's payoff.

## Next

### Theme C — Spec coverage (broaden what runs)

Ordered by value (each unlocks a class of real Java programs):

- [ ] **Exceptions** (`try`/`catch`/`athrow`) — exception tables already
  parsed in `rtda.Method`; wire into the dispatch loop's error path. Bridge
  to Go `panic`/`recover` for native-thrown exceptions. **Prerequisite for
  reflection** (APIs throw `NoSuchMethodException` etc.).
- [ ] **`invokedynamic` / lambdas / string concat (Java 9+)** — model
  `CallSite` + bootstrap methods. `LambdaMetafactory` synthesizes adapter
  classes. Or: compile fixtures with `-source 8` to get StringBuilder-based
  concat (current workaround).
- [ ] **`invokeinterface`** — `rtda.Class` supports itable partially; wire
  the opcode into the interpreter + emitter.
- [ ] **Reflection** (`java.lang.reflect`) — native methods on retained
  `rtda` metadata (ADR-0007). Needs exceptions first.
- [ ] **More core classes** — `Math`, `Integer`/`Long` parse/format,
  `Object` helpers (`equals`/`hashCode`/`toString`), `String` proper
  (char-array-backed, `charAt`/`substring`).

### Theme D — Concurrency (the "Go runtime" payoff)

- [ ] **`java.lang.Thread` → goroutine** — multi-`Thread` runtime.
- [ ] **Per-object monitors** for `synchronized` (lazily-allocated
  `sync.Mutex` per object).
- [ ] **`wait`/`notify`** via `sync.Cond`.
- [ ] **JMM approximation** — `volatile`/happens-before (Java and Go memory
  models differ; document deviations).

### Theme E — Project hygiene

- [ ] More unit tests (every constant-pool tag, object/field layout).
- [ ] Fuzzing the class-file parser.
- [ ] `catty build` from anywhere (go.mod template / SDK path).
- [ ] Graceful AOT fallback (un-AOT-able main → interpreter wrapper).
- [ ] Instance method AOT support.

## Performance summary

| Engine | fib(35) | Notes |
|---|---|---|
| catty AOT (`catty build`) | ~44 ms | native speed; on par with HotSpot JIT |
| `java` (HotSpot JIT) | ~50 ms | baseline |
| `java -Xint` | ~600 ms | ~14× AOT |
| catty interpreter | ~4.5 s | ~100× AOT |
| catty IR executor (`-ir`) | ~4.8 s | validation tier; ADR-0006 |

## What can't catty run yet?

Minecraft? No. catty runs simple static Java programs (HelloWorld, Fibonacci,
Factorial, OOP demos). Real applications need: exceptions, concurrency,
reflection, `invokedynamic`, JNI, and the full JDK class library (~17000
classes) — each is a documented Theme C/D item. The AOT path proves the thesis
(bytecode → Go → native speed); the class-library gap is the remaining work.
