# Catty

Catty builds a supported Java 21 application into one standalone Go executable.
It compiles Java bytecode to Go where possible and runs the remaining methods in
a built-in interpreter. The output runs on the target OS and architecture
without Java or Go installed.

## Install

The builder requires Go 1.21 or newer on `PATH`. Use JDK 21 when compiling
Java source. An existing Java 21 JAR can be built with Go alone.

### From source

```sh
git clone https://github.com/LizBing/catty.git
cd catty
mkdir -p bin
go build -o bin/catty ./cmd/catty
export PATH="$PWD/bin:$PATH"
catty --help
```

Keep `bin/catty` inside the checkout so it can find the Go runtime sources.
To use Catty in future shells, add the checkout's absolute `bin` path to
`PATH` in your shell configuration.

### From a release archive

If you have a release archive for your OS and architecture, place it and its
checksum file in the same directory. Set `release` to your archive's name
without `.tar.gz`.

```sh
release=catty-0.1.0-preview.1-darwin-arm64
shasum -a 256 -c "$release.tar.gz.sha256" &&
  mkdir -p "$HOME/.local/opt" &&
  tar -xzf "$release.tar.gz" -C "$HOME/.local/opt" &&
  export PATH="$HOME/.local/opt/$release/bin:$PATH" &&
  catty --help
```

Keep the extracted directory together. To switch versions, extract releases
side by side and point `PATH` at the version you want.

## Build and run

In your application's directory, pass your existing JAR and main class to
`catty build`. For a quick test, save this as `Hello.java` and compile it with
JDK 21.

```java
package demo;

public class Hello {
    public static void main(String[] args) {
        System.out.println("Hello from Catty");
        if (args.length > 0) {
            System.out.println(args[0]);
        }
    }
}
```

```sh
mkdir -p classes
javac --release 21 -d classes Hello.java
jar --create --file hello.jar -C classes .

catty build --enable-lifecycle --classpath hello.jar --main demo.Hello \
  --report dist/hello.json --output dist/hello
./dist/hello world
```

The executable prints `Hello from Catty` and `world` on separate lines.
Everything after the executable name goes to Java's `main(String[] args)`.
Deploy the executable and any files your application reads. Keep the JSON
report for build diagnostics.

Catty accepts Java 21 class directories and single-release classpath JARs.
Its Java standard library and bytecode support are defined by the selected
profile. The report lists supported opcodes and platform members. Applications
that use lambdas, other `invokedynamic` features, or library APIs outside that
inventory need additional runtime support.

## CLI reference

`catty build` builds an application. `catty help`, `catty -h`, and
`catty --help` print the top-level usage. Run `catty build --help` for the build
flag list.

| Build argument | Meaning |
| --- | --- |
| `--classpath PATH` | Required. Ordered class directories and single-release JARs. Separate entries with `:` on macOS and Linux, or `;` on Windows. If multiple entries define a class, the first entry wins. Example for macOS and Linux is `app.jar:lib/dependency.jar`. |
| `--main CLASS` | Required. Binary name of the entry class with a `public static void main(String[] args)` method, such as `demo.Hello`. Catty takes the entry point from this flag. |
| `--report FILE` | Required. Destination for the JSON build report, including input digests, profile, and engine choice for each method. |
| `--output FILE` | Required. Destination for the standalone executable. Use a different path from `--report`. |
| `--mode MODE` | Defaults to `auto`. See the engine modes below. |
| `--policy FILE` | Optional JSON file that overrides individual method engines in `auto` mode. |
| `--generated-source DIR` | Optional directory for the generated `main.go`, `generated.go`, and `go.mod`. Use it to inspect the Go source Catty generates. |
| `--enable-plugins` | Selects Gate 3. Enables runtime loading of supported plugin JARs. Plugin methods run in the interpreter. |
| `--enable-threads` | Selects Gate 4. Adds Java platform threads, monitors, and wait/notify to the Gate 3 features. |
| `--enable-lifecycle` | Selects Gate 5. Adds application exit, halt, shutdown hooks, and process diagnostics to the Gate 4 features. |
| `-h`, `--help` | Prints the build flags when passed to `catty build`. |

Without a profile flag, Catty selects Gate 2. When several flags are supplied,
`--enable-lifecycle` takes precedence, followed by `--enable-threads` and
`--enable-plugins`.

### Engine modes and policy

- `auto` compiles eligible application methods ahead of time and interprets
  admitted bytecode that lacks an AOT lowering.
- `aot` requires AOT compilation for every concrete application method in the
  image. The build fails if one lacks a lowering.
- `interp` retains all concrete application methods for the interpreter.

A policy identifies a concrete method declared in the application by its
slash-separated owner, method name, and JVM descriptor. Set the engine to
`aot` or `interp`. For example, save the following as `policy.json` and add
`--policy policy.json` to the quickstart build to interpret `demo.Hello.main`.

```json
{
  "methods": [
    {
      "owner": "demo/Hello",
      "name": "main",
      "descriptor": "([Ljava/lang/String;)V",
      "engine": "interp"
    }
  ]
}
```

### Diagnostics

| Environment variable | Set on | Meaning |
| --- | --- | --- |
| `CATTY_MODULE_ROOT` | Builder | Points to the Catty source checkout if you run a standalone copy of the builder CLI. |
| `CATTY_KEEP_STAGE` | Builder | Set to a nonempty value to preserve the temporary Go build directory. The builder prints its path to stderr. |
| `CATTY_TRACE` | Generated app | Path to a JSON Lines execution trace. |
| `CATTY_METRICS` | Gate 5 app | Path to the JSON metrics written at process exit, including execution, threads, and Go heap and GC counters. |

For architecture, specs, test plans, and release planning, see the
[developer documentation](https://github.com/LizBing/catty/tree/main/docs).
