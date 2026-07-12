# R2 Unsafe caller graph — Temurin 25.0.3+9

## Scope and method

This report replaces the earlier estimate that Integer/Long conversion,
Double parsing, and ordinary HashMap operations all enter a generic “~50 Unsafe
native” cascade.

Environment:

```text
openjdk version 25.0.3 2026-04-21 LTS
Temurin-25.0.3+9
```

Evidence commands use the installed JDK's class files:

```sh
javap -p -c java.lang.Integer
javap -p -c java.lang.Long
javap -p -c java.lang.Double
javap -p -c jdk.internal.util.DecimalDigits
javap -p -c jdk.internal.math.FloatingDecimal
javap -p -c jdk.internal.misc.Unsafe
javap -p -c 'java.util.HashMap$UnsafeHolder'
javap -p -c java.util.HashMap
```

Catty probes live in
[`../../tests/research/JavaBasePathProbes.java`](../../tests/research/JavaBasePathProbes.java).
They were compiled separately and run with a JDK image extracted exactly as CI
does, using a 12-second per-process timeout.

This is a method-level graph for the named entry points, not a promise about
every method in the containing classes.

## Integer.toString(int)

```text
Integer.toString(int)
├─ DecimalDigits.stringSize(int)                         bytecode only
├─ new byte[]
├─ DecimalDigits.uncheckedGetCharsLatin1/UTF16
│  ├─ DecimalDigits.<clinit>
│  │  └─ Unsafe.getUnsafe                               bytecode getter
│  │     └─ Unsafe.<clinit>
│  │        ├─ Unsafe.registerNatives                   native, semantic setup
│  │        ├─ Unsafe.arrayBaseOffset0 × 9 array kinds  native
│  │        ├─ Unsafe.arrayIndexScale0 × 9 array kinds  native
│  │        └─ UnsafeConstants VM-provided values       runtime injection needed
│  ├─ Latin1: Unsafe.putByte                            native
│  └─ UTF16: Unsafe.putCharUnaligned                    bytecode
│     └─ Unsafe.putShort or Unsafe.putByte parts        native
└─ synthetic String(byte[], coder) bridge
```

The checked probe exits zero but prints NUL bytes. That is not a hang and not a
large Unsafe dependency: the current generic native stub makes `putByte` a
silent no-op.

Minimum honest profile for the selected path:

- strict `registerNatives` classification/setup;
- all array base/scale initialization required by `Unsafe.<clinit>`;
- meaningful UnsafeConstants values or a reviewed Catty replacement;
- `putByte`; and `putShort`/byte-parts for UTF-16.

## Long.toString(long)

```text
Long.toString(long)
├─ DecimalDigits.stringSize(long)
├─ new byte[]
├─ DecimalDigits.uncheckedGetCharsLatin1/UTF16
│  └─ same Unsafe path as Integer.toString
└─ synthetic String(byte[], coder) bridge
```

The probe likewise exits zero and prints NUL bytes due to the silent put stub.
No long-specific CAS, field offset, fence, or raw-memory method is needed by
this conversion path.

## Double.parseDouble(String)

```text
Double.parseDouble
└─ FloatingDecimal.parseDouble
   ├─ FloatingDecimal.readJavaFormatString
   │  ├─ String.length/charAt
   │  ├─ byte[] digit buffer
   │  └─ FloatingDecimal ASCIIToBinary converter objects
   └─ ASCIIToBinaryConverter.doubleValue
```

No direct Unsafe or DecimalDigits call appears in the inspected parse path.
The current Catty probe exceeds 12 seconds. This is a separate runtime/class
library investigation, not evidence that `Double.parseDouble` requires the U0
or U1 profile. FloatingDecimal also initializes a ThreadLocal used by formatting;
the parse-specific timeout must be minimized before assigning a root cause.

## HashMap basic operations

Selected operations:

```text
new HashMap
put / get / remove / size
```

Their inspected bytecode has no Unsafe edge. The Catty probe fails immediately
with:

```text
java.lang.IllegalStateException: Not yet initialized
```

The message exists in `jdk.internal.misc.VM`, indicating a VM initialization
dependency rather than an Unsafe cascade. Root-cause minimization remains open.

HashMap does have a separate Unsafe path for Java serialization:

```text
HashMap.readObject
└─ HashMap$UnsafeHolder.putLoadFactor
   ├─ Unsafe.getUnsafe
   ├─ Unsafe.objectFieldOffset(Class, "loadFactor")
   └─ Unsafe.putFloat
```

That path belongs to serialization/I/O coverage and is not required for the
ordinary HashMap R2 smoke program.

## Corrected conclusions

| Entry point | Current result | Unsafe requirement |
|---|---|---|
| Integer.toString | exits 0, incorrect NUL output | narrow U0/U1 array-write path |
| Long.toString | exits 0, incorrect NUL output | same narrow U0/U1 path |
| Double.parseDouble | timeout | no direct Unsafe edge found; separate investigation |
| HashMap basic operations | `VM: Not yet initialized` | no Unsafe edge; separate VM-init investigation |
| HashMap deserialization | not in current smoke scope | field token + `putFloat` |

The full Unsafe profiles remain necessary for R2 concurrency and broader JDK
code, but they must not be justified by incorrectly grouped entry points.
