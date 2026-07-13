#!/usr/bin/env bash
# R2 research differential harness.
#
# Compares Temurin 25 against catty's three engines (Interpreter, IR, AOT) on a
# Java 25 fixture matrix that targets class/interface initialization, init
# failure semantics, the bootstrap resolution boundary, and String UTF-16
# semantics.
#
# Research-only (ADR-0017 pins Temurin 25 as the differential reference;
# ADR-0016 requires per-engine evidence). catty is run in pure-synthetic mode
# (-no-boot) for a CONTROLLED SEMANTIC differential — the java.base availability
# boundary is documented separately by the ReachUnsafe fixture. This harness
# does not modify production code; it only produces evidence.
#
# Usage: bash docs/workstreams/r2-evidence/run-r2-diff.sh
#   R2_RESULTS_DIR  — if set, write results to $R2_RESULTS_DIR/run-r2-results.txt
#                     instead of the default baseline location.
set -u

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"      # docs/workstreams/r2-evidence -> repo root
FIX="$SCRIPT_DIR/fixtures"
RESULTS="${R2_RESULTS_DIR:-$SCRIPT_DIR}/run-r2-results.txt"

fail_closed() { echo "r2-diff: $*" >&2; exit 1; }

command -v java  >/dev/null || fail_closed "java not found (Temurin 25 required)"
command -v javac >/dev/null || fail_closed "javac not found (Temurin 25 required)"
command -v perl  >/dev/null || fail_closed "perl not found (used for portable timeout)"
command -v go    >/dev/null || fail_closed "go not found (catty build + AOT)"

# Portable wall-clock timeout: perl alarm then exec. $1=seconds, rest=cmd.
to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

BIN="$(mktemp -t catty-r2.XXXXXX)"
STAGE="$(mktemp -d -t r2-stage.XXXXXX)"
trap 'rm -rf "$BIN" "$STAGE"' EXIT

{
  echo "=== R2 differential run ==="
  echo "repo:    $ROOT"
  echo "commit:  $(cd "$ROOT" && git rev-parse HEAD 2>/dev/null || echo unknown)"
  echo "branch:  $(cd "$ROOT" && git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)"
  echo "base:    ecb086e"
  echo "java:    $(java -version 2>&1 | head -1)"
  echo "javac:   $(javac -version 2>&1)"
  echo "boot:    CATTY_BOOT=${CATTY_BOOT:-<unset>}; JAVA_HOME=${JAVA_HOME:-<unset>}"
  echo "mode:    catty -no-boot (pure-synthetic, controlled semantic differential)"
  echo "compare: combined stdout+stderr AND exit code must equal Temurin 25"
  echo "results: $RESULTS"
  echo
} | tee "$RESULTS"

echo "==> building catty" | tee -a "$RESULTS"
go build -o "$BIN" "$ROOT/cmd/jvm" >>"$RESULTS" 2>&1 || fail_closed "catty build failed"

T_RUN=20   # seconds: java / interpreter / IR / aot-binary run
T_AOT=120  # seconds: catty build (shells out to go build)

printf "%-26s %-10s %-10s %-10s %-10s\n" "fixture" "Temurin25" "Interp" "IR" "AOT" | tee -a "$RESULTS"
printf "%-26s %-10s %-10s %-10s %-10s\n" "-------" "---------" "------" "--" "---" | tee -a "$RESULTS"

for jf in "$FIX"/*.java; do
  [ -e "$jf" ] || continue
  name="$(basename "$jf" .java)"
  d="$STAGE/$name"; mkdir -p "$d"

  if ! javac --release 25 -nowarn -d "$d" "$jf" 2>"$d/javac.err"; then
    printf "%-26s %-10s\n" "$name" "javac-FAIL" | tee -a "$RESULTS"
    { echo "----- $name : javac FAILED"; cat "$d/javac.err"; } >>"$RESULTS"
    continue
  fi

  tum=$( { cd "$d" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); tum_rc=$?
  lp=$(  { cd "$d" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); lp_rc=$?
  ir=$(  { cd "$d" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?

  # AOT (best effort): build, then run the produced binary.
  m_aot="NO-BUILD"; aot="(no binary)"; aot_rc="-"
  # Run catty build from REPO ROOT so the emitted program's `go build` has
  # module context (go.mod → module mode → resolves catty/rtda & catty/runtime).
  ( cd "$ROOT" && to "$T_AOT" "$BIN" build -cp "$d" -no-boot -o "$d/aot.bin" "$name" ) >"$d/aot_build.out" 2>&1
  if [ -x "$d/aot.bin" ]; then
    aot=$( { cd "$d" && to "$T_RUN" "$d/aot.bin"; } 2>&1 ); aot_rc=$?
    m_aot="MISMATCH"; [ "$tum" = "$aot" ] && [ "$tum_rc" = "$aot_rc" ] && m_aot="match"
  else
    aot=$(tail -3 "$d/aot_build.out" 2>/dev/null)
  fi

  ref="ref"; [ "$tum_rc" -ne 0 ] 2>/dev/null && ref="ref($tum_rc)"
  m_lp="MISMATCH"; [ "$tum" = "$lp" ] && [ "$tum_rc" = "$lp_rc" ] && m_lp="match"
  m_ir="MISMATCH"; [ "$tum" = "$ir" ] && [ "$tum_rc" = "$ir_rc" ] && m_ir="match"

  printf "%-26s %-10s %-10s %-10s %-10s\n" "$name" "$ref" "$m_lp" "$m_ir" "$m_aot" | tee -a "$RESULTS"

  {
    echo "----- $name -----"
    echo "[$tum_rc] temurin25:"; echo "$tum" | sed 's/^/      /'
    echo "[$lp_rc] catty interp:"; echo "$lp" | sed 's/^/      /'
    echo "[$ir_rc] catty ir:"; echo "$ir" | sed 's/^/      /'
    echo "[$aot_rc] catty aot:"; echo "$aot" | sed 's/^/      /'
  } >>"$RESULTS"
done

echo | tee -a "$RESULTS"
echo "==> full evidence written to: $RESULTS" | tee -a "$RESULTS"
