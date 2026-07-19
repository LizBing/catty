#!/usr/bin/env bash
# R3 reflection/dynamic research baseline harness.
#
# Runs the fixed 24-fixture matrix against Temurin 25, catty Interpreter, IR,
# and AOT. This is descriptive research evidence: catty mismatches, parse
# failures, panics, and NO-BUILD are recorded rather than treated as supported.
# Missing tools/fixtures/rows, reference compilation/execution failure, or an
# unbounded process fail closed.
#
# Usage:
#   R3_RESULTS_DIR=<empty-output-dir> bash run-r3-baseline.sh
#
# Optional:
#   R3_JAVA_HOME=<jdk-home>  JDK whose java.base image is extracted for catty.
#   CATTY_BOOT=<dir>         Reuse an already extracted java.base directory.
set -u

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
FIXTURE_DIR="$SCRIPT_DIR"

die() { echo "r3-baseline: $*" >&2; exit 1; }
to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

[ -n "${R3_RESULTS_DIR:-}" ] || die "R3_RESULTS_DIR is required"
[ ! -e "$R3_RESULTS_DIR" ] || die "refusing existing R3_RESULTS_DIR: $R3_RESULTS_DIR"
mkdir -p "$R3_RESULTS_DIR" || die "cannot create results directory"
RESULTS="$R3_RESULTS_DIR/results.txt"

command -v java >/dev/null || die "java not found"
command -v javac >/dev/null || die "javac not found"
command -v jimage >/dev/null || die "jimage not found"
command -v go >/dev/null || die "go not found"
command -v perl >/dev/null || die "perl not found"
command -v shasum >/dev/null || die "shasum not found"
java -version 2>&1 | head -1 | grep '"25\.' >/dev/null || die "Java 25 required"
javac -version 2>&1 | grep '^javac 25\.' >/dev/null || die "javac 25 required"

FIXTURES=(
  ClassIdentity
  ClassForNameInit
  ClassQueries
  DeclaredMembers
  PrimitiveAndArrayClass
  MissingClass
  MethodInvoke
  ConstructorInvoke
  FieldGetSet
  StaticReflectiveInit
  ReflectiveConversions
  ReflectiveFailures
  ClassAnnotation
  MemberAnnotation
  AnnotationDefaults
  InheritedRepeatableAnnotation
  StringConcatIndy
  StatelessLambda
  CapturingLambda
  MethodReference
  BootstrapFailureOnce
  ProxyDispatch
  ProxyObjectMethods
  ProxyFailureAndDefault
)
EXPECTED=24
[ "${#FIXTURES[@]}" -eq "$EXPECTED" ] || die "internal fixture count is not 24"

for name in "${FIXTURES[@]}"; do
  [ -f "$FIXTURE_DIR/$name.java" ] || die "missing fixture: $name.java"
done
actual_java=$(find "$FIXTURE_DIR" -maxdepth 1 -name '*.java' -type f | wc -l | tr -d ' ')
[ "$actual_java" -eq "$EXPECTED" ] || die "fixture directory has $actual_java Java files, expected 24"
(cd "$FIXTURE_DIR" && shasum -a 256 -c manifest.sha256) >/dev/null 2>&1 ||
  die "fixture source hash mismatch"

BIN="$(mktemp -t catty-r3.XXXXXX)"
STAGE="$(mktemp -d -t r3-stage.XXXXXX)"
BOOT_STAGE=""
cleanup() {
  rm -f "$BIN"
  rm -rf "$STAGE"
  if [ -n "$BOOT_STAGE" ]; then rm -rf "$BOOT_STAGE"; fi
}
trap cleanup EXIT

BOOT_DIR="${CATTY_BOOT:-}"
if [ -z "$BOOT_DIR" ] || [ ! -d "$BOOT_DIR/java/lang" ]; then
  JAVA_HOME_VALUE="${R3_JAVA_HOME:-}"
  if [ -z "$JAVA_HOME_VALUE" ] && [ -x /usr/libexec/java_home ]; then
    JAVA_HOME_VALUE="$(/usr/libexec/java_home 2>/dev/null || true)"
  fi
  [ -n "$JAVA_HOME_VALUE" ] || die "set CATTY_BOOT or R3_JAVA_HOME"
  [ -f "$JAVA_HOME_VALUE/lib/modules" ] || die "JDK modules image not found under $JAVA_HOME_VALUE"
  BOOT_STAGE="$(mktemp -d -t r3-java-base.XXXXXX)"
  to 120 jimage extract --dir "$BOOT_STAGE" "$JAVA_HOME_VALUE/lib/modules" >/dev/null 2>&1 ||
    die "jimage extraction failed/timed out"
  BOOT_DIR="$BOOT_STAGE/java.base"
fi
[ -d "$BOOT_DIR/java/lang" ] || die "invalid extracted java.base: $BOOT_DIR"

T_RUN=20
T_AOT=120

{
  echo "=== R3 reflection/dynamic research baseline (24 fixtures) ==="
  echo "acceptance-anchor: 6cf3636"
  echo "worktree-head:     $(cd "$ROOT" && git rev-parse HEAD)"
  echo "branch:            $(cd "$ROOT" && git rev-parse --abbrev-ref HEAD)"
  echo "worktree-dirty:    $(cd "$ROOT" && test -n "$(git status --porcelain)" && echo yes || echo no)"
  echo "java:              $(java -version 2>&1 | head -1)"
  echo "javac:             $(javac -version 2>&1)"
  echo "jimage:            $(jimage --version 2>&1)"
  echo "go:                $(go version 2>&1)"
  echo "boot:              extracted java.base ($BOOT_DIR)"
  echo "timeouts:          run=${T_RUN}s aot-build=${T_AOT}s"
  echo "compare:           combined stdout+stderr and exit status"
  echo "policy:            descriptive baseline; no catty state is counted as support"
  echo
} >"$RESULTS"

(cd "$ROOT" && go build -o "$BIN" ./cmd/jvm) >>"$RESULTS" 2>&1 || die "catty build failed"

printf '%-32s %-11s %-16s %-16s %s\n' "fixture" "Temurin25" "Interpreter" "IR" "AOT" >>"$RESULTS"
printf '%-32s %-11s %-16s %-16s %s\n' "-------" "---------" "-----------" "--" "---" >>"$RESULTS"

classify_run() {
  local ref_out="$1" ref_rc="$2" actual_out="$3" actual_rc="$4"
  if [ "$actual_rc" -eq 142 ]; then
    printf 'TIMEOUT'
  elif [ "$actual_rc" -ne 0 ]; then
    printf 'EXIT(%s)' "$actual_rc"
  elif [ "$ref_rc" -eq "$actual_rc" ] && [ "$ref_out" = "$actual_out" ]; then
    printf 'MATCH'
  else
    printf 'MISMATCH'
  fi
}

rows=0
for name in "${FIXTURES[@]}"; do
  dir="$STAGE/$name"
  mkdir -p "$dir" || die "cannot create stage for $name"
  javac --release 25 -nowarn -d "$dir" "$FIXTURE_DIR/$name.java" >"$dir/javac.out" 2>&1 ||
    die "javac failed for $name"
  if [ "$name" = "BootstrapFailureOnce" ]; then
    rm -f "$dir/BootstrapFailureTarget.class"
  fi

  ref=$( { cd "$dir" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); ref_rc=$?
  [ "$ref_rc" -eq 0 ] || die "Temurin failed/timed out for $name (rc=$ref_rc)"

  interp=$( { cd "$dir" && CATTY_BOOT="$BOOT_DIR" to "$T_RUN" "$BIN" -cp . "$name"; } 2>&1 ); interp_rc=$?
  ir=$( { cd "$dir" && CATTY_BOOT="$BOOT_DIR" to "$T_RUN" "$BIN" -cp . -ir "$name"; } 2>&1 ); ir_rc=$?
  interp_state=$(classify_run "$ref" "$ref_rc" "$interp" "$interp_rc")
  ir_state=$(classify_run "$ref" "$ref_rc" "$ir" "$ir_rc")

  aot_build_out="$dir/aot-build.out"
  (cd "$ROOT" && CATTY_BOOT="$BOOT_DIR" to "$T_AOT" "$BIN" build -cp "$dir" -o "$dir/aot.bin" "$name") >"$aot_build_out" 2>&1
  aot_build_rc=$?
  aot=""
  aot_rc="-"
  if [ -x "$dir/aot.bin" ]; then
    aot=$( { cd "$dir" && to "$T_RUN" "$dir/aot.bin"; } 2>&1 ); aot_rc=$?
    aot_state=$(classify_run "$ref" "$ref_rc" "$aot" "$aot_rc")
  elif [ "$aot_build_rc" -eq 142 ]; then
    aot_state="BUILD-TIMEOUT"
  else
    aot_state="NO-BUILD"
    aot=$(tail -5 "$aot_build_out" 2>/dev/null)
  fi

  printf '%-32s %-11s %-16s %-16s %s\n' "$name" "REF" "$interp_state" "$ir_state" "$aot_state" >>"$RESULTS"
  {
    echo "----- $name -----"
    echo "[$ref_rc] Temurin25"
    printf '%s\n' "$ref"
    echo "[$interp_rc] Interpreter ($interp_state)"
    printf '%s\n' "$interp"
    echo "[$ir_rc] IR ($ir_state)"
    printf '%s\n' "$ir"
    echo "[$aot_rc; build=$aot_build_rc] AOT ($aot_state)"
    printf '%s\n' "$aot"
  } >>"$RESULTS"
  rows=$((rows + 1))
done

[ "$rows" -eq "$EXPECTED" ] || die "produced $rows rows, expected 24"
echo >>"$RESULTS"
echo "rows: $rows/24 complete" >>"$RESULTS"
echo "evidence: $RESULTS"
