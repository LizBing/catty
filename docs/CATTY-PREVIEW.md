# Catty 0.1 preview and installation plan

**Draft.** The first preview is for developers packaging small Java 21 command-line or batch applications as standalone executables. Start with macOS arm64 and one real pilot application. Advertise another target once the release checks pass on that target.

The developer installs Catty and Go on the build machine, supplies compiled
Java classes or JARs, and gets one application executable and a JSON build
report. The deployment machine needs the matching OS and CPU architecture and
any inputs the application reads. The generated executable carries the Catty
runtime and runs without Java or Go installed.

## What the preview supports

- **Inputs.** Ordinary Java 21 class files (65.0) from class directories or
  single-release classpath JARs. Pass the fully qualified main class explicitly.
  For multiple entries, use the application's classpath order. Catty admits
  the class files in each supplied entry, so the compatibility check covers
  the supplied dependencies as well as the main class.
- **Profile.** Use `--enable-lifecycle --mode auto` for the preview workflow.
  It selects `catty-jvm21-gate5`, image ABI 5, and automatic selection of
  AOT compilation or interpretation per method. The build report lists the
  supported opcodes and bootstrap members, per-method engine choices, and
  application references retained for lazy linkage.
- **Compatibility.** The supported surface is the published Catty profile.
  Check dependencies for `invokedynamic`, including common lambda and string
  concatenation forms, and for requirements around modules, multi-release
  JARs, JAR resources, JNI, or JDK APIs outside the bootstrap inventory. The
  first pilot determines which missing features are release blockers.
- **Plugins.** Use trusted, explicitly supplied JARs. Plugin classes run
  interpreted inside the application process and remain subject to Gate 3's
  admission rules.
- **Application shutdown.** Gate 5 supports natural shutdown, explicit exit
  and halt, and shutdown hooks. OS signal integration is outside the Gate 5
  profile. Include signal-driven shutdown in the acceptance checks if the
  pilot runs under a supervisor that requires it.

[Compatibility contract](ARCHITECTURE.md#1-compatibility-contract) ·
[CLI and profile details](../README.md#cli-reference) ·
[Plugin admission contract](GATE3.md#1-profile-and-build-contract)

## Install Catty

Download the archive and its SHA-256 checksum file from the same preview release.
The example name below identifies the proposed first preview artifact. Check
`uname -s` and `uname -m` for `Darwin` and `arm64`, and use the release's
validated Go version on the build machine. The dry run uses Go 1.27.1.
Use a JDK 21 when compiling Java sources or comparing results with OpenJDK.
Building from an existing JAR requires Go on `PATH`.

From the download directory, run the following commands.

```sh
release=catty-0.1.0-preview.1-darwin-arm64
shasum -a 256 -c "$release.tar.gz.sha256" &&
  mkdir -p "$HOME/.local/opt" &&
  tar -xzf "$release.tar.gz" -C "$HOME/.local/opt" &&
  export PATH="$HOME/.local/opt/$release/bin:$PATH" &&
  catty --help
```

Keep the extracted directory together. To use Catty in future shells, put the
versioned `bin` directory on `PATH` in your shell configuration. To upgrade,
extract the new release beside the old one and point `PATH` at the new
directory. To roll back, point it at the previous release. To uninstall, remove
that release's directory and its `PATH` entry.

## Build and run an application

From the Java application's directory, replace `app.jar` and
`com.example.Main` with your JAR and entry class, and pass your application's
arguments in place of `hello`. On macOS, use a colon to separate classpath
entries, as in `app.jar:lib/dependency.jar`.

```sh
catty build --enable-lifecycle --mode auto \
  --classpath app.jar --main com.example.Main \
  --report dist/app.json --output dist/app &&
  ./dist/app hello
```

Keep `dist/app.json` with the build records. Deploy `dist/app` and the inputs
your application reads. After upgrading Catty, rebuild the executable from
the original classpath. The generated Go and image ABI are implementation
interfaces during the preview.

| Symptom | What to check |
| --- | --- |
| `catty: command not found` | Check that the selected release's `bin` directory is on `PATH`. |
| Catty cannot find its module root | Extract the full archive and invoke its `bin/catty` launcher. The bundled module should be beside `bin` and `libexec`. |
| The generated Go build fails | Check `go version`, the target OS and architecture, and the Go error printed by Catty. |
| Catty rejects a class or bootstrap member | Keep the first diagnostic and record the release, main class, and ordered classpath. Check the named dependency against the preview profile. |
| The executable builds but fails on an input | Compare its output and exit status with OpenJDK 21. Rerun with `CATTY_TRACE=trace.jsonl` and keep `dist/app.json`. |

Share a minimal reproducer with the release's issue tracker or pilot contact.
Remove private paths and application data from reports, traces, and logs before
sharing them.

## Build the preview archive

Use one archive per OS and architecture. Each archive contains a prebuilt CLI,
the matching Go runtime sources, and a small launcher that selects the bundled
module even when the working directory is another Catty checkout. The app
builder invokes the user's Go toolchain to compile the generated executable.

```text
catty-0.1.0-preview.1-darwin-arm64/
  bin/catty
  libexec/catty
  go.mod
  rt/
  plugin/
  internal/{admission,classfile,hierarchy,verify}/
  README.md
  LICENSE
```

Run the packaging script from the source revision being tested. It writes
the archive and its checksum to the output directory and copies `LICENSE`
when present. Add the approved license before distributing the archive.

```sh
sh scripts/package-preview.sh 0.1.0-preview.1 dist
```

The runtime package list is the dependency closure of `./rt` and `./plugin`
in the current module. Review it when imports or runtime assets change.
Building a Gate 5 JAR from the extracted archive exercises that closure.
Record the Git revision, validated toolchains, and archive checksum alongside
the tagged release.

## Preview release checks

1. **Prove the document works.** On a clean machine for the advertised target,
   download the archive through the intended release channel, verify its
   checksum, extract it, and build an existing Java 21 JAR. Repeat with
   `GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off`. Verify the build report selects
   Gate 5 and the copied executable runs on a matching machine without the
   builder, the JAR, or either toolchain.
2. **Prove the pilot works.** Have an independent developer follow the quickstart
   with their own app. Compare stdout, stderr, and exit status with OpenJDK 21
   on representative inputs. Include one unsupported-input case and check that
   the diagnostic identifies the blocking dependency or feature.
3. **Verify the release revision.** Require a passing CI run for the tagged
   revision. The workflow runs Go vet and race tests, Gates 1–5 acceptance tests,
   and an archive smoke check on macOS arm64. [CI workflow](../.github/workflows/ci.yml)
4. **Measure what matters to the pilot.** Compare startup to the first useful
   result, package size, peak RSS, and sustained throughput on the advertised
   target. Set acceptable limits with the pilot before tagging. Keep the
   measured revision, workload inputs, and environment with the results.
5. **Publish the support contract.** Include the validated Go version and target,
   supported profile, known compatibility limits, trusted-plugin policy,
   release notes, license, checksum, and issue-reporting location.

Ship the preview when an independent developer can build and deploy a real
application from the release archive and the published checks pass. Use that
experience to decide which Java APIs and deployment targets to support next.

[Gate 5 evidence and benchmarks](../out/gate5/EVIDENCE.md) ·
[Recorded Gate 4 test flake](../out/gate5/EVIDENCE.md#pre-existing-gate-4-timer-arming-flake-recorded-separately)

## Local installation dry run

On macOS arm64 with Go 1.27.1, the recipe produced an archive that passed its
checksum check. The extracted CLI built a Java 21 JAR from an unrelated
directory with Go downloads disabled and with Java absent from `PATH`. The
build report selected Gate 5, and the copied executable matched OpenJDK 21
after the input JAR and class directory were removed. The clean-machine and
independent-pilot checks remain release criteria.
