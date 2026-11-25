# Gate 1 acceptance tests

Run the full suite from the repository root.

```sh
python3 -B tests/acceptance.py
```

The suite requires Python 3.9 or newer, Go, and the `java`, `javac`, `javap`, and `jar` tools from a JDK that accepts `javac --release 21`. A missing or incompatible tool fails the run. The runner builds `./cmd/catty` once and uses a fresh temporary directory. Its first output line gives the location of the class files, command logs, bytecode disassemblies, binaries, build reports, generated Go source, and runtime traces. Command logs include arguments, exit status, stdout, and stderr.

```sh
python3 -B tests/acceptance.py --prepare-only
```

Use this command to verify the fixtures and OpenJDK expectations before Catty is implemented. It compiles the Java sources, checks their instructions against §3 of [IMPLEMENTATION.md](IMPLEMENTATION.md), and runs each applicable oracle with `-Xverify:all`. The two malformed verification fixtures must fail in OpenJDK. The exception boundary fixture must pass OpenJDK verification after the test changes the handler's `end_pc` to the offset of the throwing call.

## Cases

| Case | Evidence the suite checks |
| --- | --- |
| Crossing in three modes | Exact `2\n` stdout. The same class files select the expected engines in all three reports. The traces show direct and reflective interface selection, nested callbacks, original throw and wrapper identity, exact catch offsets, and balanced frames in one process and one logical Java thread. |
| Standalone launch and generated Go | The mixed binary runs alone in an empty directory with an empty search path and the input class directory temporarily moved. The suite locates the generated functions for `Box.bump` and `Crossing.probe` and checks their Go syntax trees for an update and control flow. The source review below completes the AOT evidence. |
| Automatic fallback | Positive and negative signed remainder, the minimum int with divisor −1, and the zero-divisor exception match OpenJDK. Only `rem` keeps JVM code and its reason names `irem`. Forced AOT fails, preserves the existing executable byte-for-byte, and the preserved binary runs again. |
| Dispatch and reflection | Two receivers reach the same interface call site in opposite orders. A name from `args[0]` finds a method invoked only reflectively. The fixture checks inherited lookup, override selection, returned reference identity, a missing method, a wrong argument type, int boxing and unboxing, a null zero-argument array, and original target throwable identity. |
| Class state and fields | The class mirror precedes initialization. One AOT caller and an interpreted method share the static field. Inherited and hidden fields remain independent. Two uses of a failing class produce the expected first-use and subsequent-use errors. |
| Reader, policy, and packaging | Directory and JAR launch agree. Reversing the classpath changes the winning definition and origin. Version 66, an invalid branch into an operand, a malformed stack map, synchronized methods, a dynamic bootstrap, a reserved class name, and a missing main are rejected. Policy negatives cover unknown, duplicate, conflicting, abstract, and unlowerable selections. |
| Exceptions | In both engines, null, bounds and array-store errors reach the Java handlers. A call inside the range catches the original throwable before a broader handler. A call at the exclusive end unwinds to the caller with the same throwable. |
| Runtime invariants | A fresh `go test -count=1 -json ./...` run must include the two passing runtime cases below. |

## Report and trace fields

These field names give the suite and builder a shared format for the evidence in §§3 and 6 of the implementation spec. Additional fields are allowed.

| Record | Required fields |
| --- | --- |
| Build report | `profile` as a string, `imageABI` as a string or integer, `main` as the supplied Java binary name, and `toolchains.go` as a version string. The fixture logs record the JDK version independently. |
| `inputs[]` | `origin` as the absolute classpath entry and `digest` as a nonempty digest string. Preserve classpath order. |
| `classes[]` | `name` as the JVM internal class name and `origin` as the winning absolute class file path or `/absolute/archive.jar!/name.class`. |
| `methods[]` | `owner`, `name`, `descriptor`, `engine`, `reason`, and `retainedCodeBytes`. Use JVM internal owner names. For interpreted methods count the original `Code` byte array. Use zero for all other entry kinds. |
| Generated method metadata | When `--generated-source DIR` is requested, copy the generated Go source used by the compiler to `DIR` and add `generated.file` and `generated.symbol` to each AOT method record. The file path is relative to `DIR` and the symbol is a top-level Go function name. |

Each JSON line produced by `CATTY_TRACE` carries `event`, the actual OS `pid`, and a non-null scalar logical `thread` ID. Entries and exits use zero-based depth, with the outermost Java frame at depth zero. Method IDs are objects with `owner`, `name`, and `descriptor`. Throwable IDs are opaque scalars that remain stable within the trace.

| `event` | Additional fields |
| --- | --- |
| `enter` | Method ID fields, `engine`, `depth` |
| `return` | Method ID fields, `depth` |
| `unwind` | Method ID fields, `depth`, `throwable` ID |
| `dispatch` | `opcode`, `declaration` method ID, `target` method ID. Emit this for the direct `invokeinterface` before entering the selected body. |
| `reflect` | `declaration` method ID, `target` method ID. Emit this before entering the selected target. |
| `catch` | Handler method ID fields, `throwPC`, `handlerPC`, `throwable` ID |
| `wrap` | `wrapper` and `cause` IDs |

For the AOT review, open the functions printed in `logs/source-bump.json` and `logs/source-probe.json`. Follow the field read, int increment, and write on the original receiver in `bump`. Follow the throw checks in `probe` to Go handler blocks for the offsets independently checked by the trace. Verify that the generated bodies themselves execute their lowered control flow. Inspect their call paths for any handoff to a bytecode interpreter or opcode-dispatch loop.

## Runtime cases to wire into the implementation

Add these Go tests using the runtime's concrete image and interpreter APIs. The acceptance runner requires both to execute and pass.

### `TestCattyAcceptanceImageABI`

- Construct a valid minimal image with a startup initializer and main that each set an observable marker. Install and launch the valid image as the positive control.
- Use a fresh image and runtime instance for the negative case. Change only the ABI revision to an unsupported value.
- Assert that installation fails with an ABI-specific error and leaves both startup markers unset. Run the same check with an unsupported profile ID and a valid ABI revision.

### `TestCattyAcceptanceInterpreterOperandGC`

- Execute a retained method that pushes a Java object with a known nonzero int field onto the operand stack. Use an object that is absent from the Java locals, static storage, and caches while it is on that stack. The nested call takes no reference to this object.
- Make the nested call go through the image's compiled entry table. Trigger Go garbage collection in that compiled entry. Use a non-rooting identity witness or finalization probe so the test itself does not keep the object alive. Check that the compiled entry actually ran and that the test has no other strong root to the object during collection.
- After the compiled entry returns, read the object from the operand stack through interpreter instructions. Assert the original field value and object identity.

## Fixture verification

The prepare-only run succeeded with Go 1.27.1 and Temurin OpenJDK 25.0.4.1 on darwin/arm64. It verified 15 OpenJDK outcomes, including the expected arithmetic exception, and two additional verifier rejections. The Catty executable checks and runtime cases require the implementation.
