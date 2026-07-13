#!/usr/bin/env bash
# R2 concurrency research baseline harness.
#
# Records Temurin 25 and catty Interpreter/IR/AOT behavior for the fixed
# 19-fixture research matrix. Catty mismatches are expected baseline evidence;
# the harness fails only when the reference/toolchain/fixture set is incomplete,
# a required command is not bounded, or an engine row cannot be recorded.
set -u

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
BASELINE_ID="baseline-63d5658"
EVIDENCE_DIR="${R2_CONCURRENCY_RESULTS_DIR:-$ROOT/docs/workstreams/r2-concurrency-evidence/$BASELINE_ID}"
RESULT_NAME="${R2_CONCURRENCY_RESULT_NAME:-run-concurrency-results.txt}"
RESULTS="$EVIDENCE_DIR/$RESULT_NAME"

fail_closed() { echo "r2-concurrency-research: $*" >&2; exit 1; }

command -v java  >/dev/null || fail_closed "java not found (Temurin 25 required)"
command -v javac >/dev/null || fail_closed "javac not found (Temurin 25 required)"
command -v perl  >/dev/null || fail_closed "perl not found (portable timeout)"
command -v go    >/dev/null || fail_closed "go not found (catty build + AOT)"

java -version 2>&1 | head -1 | grep '"25\.' >/dev/null || fail_closed "java 25 required"
javac -version 2>&1 | grep '^javac 25\.' >/dev/null || fail_closed "javac 25 required"

mkdir -p "$EVIDENCE_DIR" || fail_closed "cannot create evidence directory: $EVIDENCE_DIR"
[ ! -e "$RESULTS" ] || fail_closed "refusing to overwrite immutable evidence: $RESULTS"

to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

BIN="$(mktemp -t catty-concurrency.XXXXXX)"
STAGE="$(mktemp -d -t concurrency-stage.XXXXXX)"
trap 'rm -rf "$BIN" "$STAGE"' EXIT

FIXTURES="
CurrentThreadIdentity
ThreadStartJoin
ThreadStartTwice
NonDaemonLiveness
DaemonLiveness
SynchronizedBlocks
SynchronizedMethods
MonitorNull
MonitorOwnership
WaitNotify
NotifyAll
InterruptStatus
InterruptWait
InterruptSleep
InterruptJoin
VolatilePublication
FinalFieldPublication
CrossThreadClassInitialization
ProducerConsumer
"

EXPECTED_COUNT=19
actual_count=0
for name in $FIXTURES; do
  [ -f "$SCRIPT_DIR/$name.java" ] || fail_closed "missing fixture: $name.java"
  actual_count=$((actual_count + 1))
done
[ "$actual_count" -eq "$EXPECTED_COUNT" ] || fail_closed "fixture count $actual_count != $EXPECTED_COUNT"

T_RUN=20
T_AOT=120

{
  echo "=== R2 concurrency research baseline ==="
  echo "baseline-id: $BASELINE_ID"
  echo "catty-base:  63d5658"
  echo "repo-head:   $(cd "$ROOT" && git rev-parse HEAD)"
  echo "branch:      $(cd "$ROOT" && git rev-parse --abbrev-ref HEAD)"
  echo "date:        $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "java:        $(java -version 2>&1 | head -1)"
  echo "javac:       $(javac -version 2>&1)"
  echo "go:          $(go version 2>&1)"
  echo "mode:        catty -no-boot; combined stdout+stderr and exit code"
  echo "timeouts:    run=${T_RUN}s aot-build=${T_AOT}s"
  echo "fixtures:    $EXPECTED_COUNT fixed research fixtures"
  echo "policy:      catty mismatch/NO-BUILD is baseline evidence, not harness failure"
  echo "evidence:    $RESULTS"
  echo
} | tee "$RESULTS"

echo "==> building catty" | tee -a "$RESULTS"
go build -o "$BIN" "$ROOT/cmd/jvm" >>"$RESULTS" 2>&1 || fail_closed "catty build failed"

printf "%-34s %-10s %-18s %-18s %s\n" "fixture" "Temurin25" "Interpreter" "IR" "AOT" | tee -a "$RESULTS"
printf "%-34s %-10s %-18s %-18s %s\n" "-------" "---------" "-----------" "--" "---" | tee -a "$RESULTS"

recorded=0
interp_match=0
ir_match=0
aot_match=0
aot_no_build=0

indent() {
  local line
  while IFS= read -r line || [ -n "$line" ]; do
    printf '      %s\n' "$line"
  done
}

comparison() {
  local ref_out="$1" ref_rc="$2" actual_out="$3" actual_rc="$4"
  if [ "$ref_out" = "$actual_out" ] && [ "$ref_rc" = "$actual_rc" ]; then
    printf 'match'
  else
    printf 'MISMATCH(rc=%s)' "$actual_rc"
  fi
}

for name in $FIXTURES; do
  dir="$STAGE/$name"
  mkdir -p "$dir"

  if ! javac --release 25 -nowarn -d "$dir" "$SCRIPT_DIR/$name.java" 2>"$dir/javac.err"; then
    cat "$dir/javac.err" >>"$RESULTS"
    fail_closed "javac failed for $name"
  fi

  ref=$( { cd "$dir" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); ref_rc=$?
  [ "$ref_rc" -eq 0 ] || fail_closed "Temurin reference failed/timed out for $name (rc=$ref_rc)"

  interp=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); interp_rc=$?
  ir=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?

  interp_status=$(comparison "$ref" "$ref_rc" "$interp" "$interp_rc")
  ir_status=$(comparison "$ref" "$ref_rc" "$ir" "$ir_rc")
  [ "$interp_status" = "match" ] && interp_match=$((interp_match + 1))
  [ "$ir_status" = "match" ] && ir_match=$((ir_match + 1))

  aot_status="NO-BUILD"
  aot=""
  aot_rc="-"
  ( cd "$ROOT" && to "$T_AOT" "$BIN" build -cp "$dir" -no-boot -o "$dir/aot.bin" "$name" ) >"$dir/aot-build.out" 2>&1
  aot_build_rc=$?
  if [ -x "$dir/aot.bin" ]; then
    aot=$( { cd "$dir" && to "$T_RUN" "$dir/aot.bin"; } 2>&1 ); aot_rc=$?
    aot_status=$(comparison "$ref" "$ref_rc" "$aot" "$aot_rc")
    [ "$aot_status" = "match" ] && aot_match=$((aot_match + 1))
  else
    aot_no_build=$((aot_no_build + 1))
    aot=$(cat "$dir/aot-build.out" 2>/dev/null)
  fi

  printf "%-34s %-10s %-18s %-18s %s\n" "$name" "ref" "$interp_status" "$ir_status" "$aot_status" | tee -a "$RESULTS"

  {
    echo "----- $name -----"
    echo "[0] temurin25:"; printf '%s' "$ref" | indent
    echo "[$interp_rc] catty interpreter:"; printf '%s' "$interp" | indent
    echo "[$ir_rc] catty IR:"; printf '%s' "$ir" | indent
    echo "[$aot_rc] catty AOT (build-rc=$aot_build_rc):"; printf '%s' "$aot" | indent
  } >>"$RESULTS"

  recorded=$((recorded + 1))
done

[ "$recorded" -eq "$EXPECTED_COUNT" ] || fail_closed "recorded $recorded rows, expected $EXPECTED_COUNT"

{
  echo
  echo "==> baseline summary"
  echo "rows:              $recorded/$EXPECTED_COUNT"
  echo "interpreter match: $interp_match/$EXPECTED_COUNT"
  echo "IR match:          $ir_match/$EXPECTED_COUNT"
  echo "AOT match:         $aot_match/$EXPECTED_COUNT"
  echo "AOT NO-BUILD:      $aot_no_build/$EXPECTED_COUNT"
  echo "result:            Pass (complete research baseline; mismatch is expected evidence)"
} | tee -a "$RESULTS"

echo "r2-concurrency-research: baseline recorded at $RESULTS" >&2
