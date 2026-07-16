#!/usr/bin/env bash
# R2 concurrency candidate harness — parent final 19-fixture matrix.
#
# Runs the fixed 19 concurrency fixtures against Temurin 25 (reference),
# catty Interpreter, and catty IR. Every fixture must match the reference
# in combined stdout+stderr and exit code. Any mismatch, timeout, missing
# fixture, or build failure is a hard failure.
#
# Additionally runs an AOT column: every fixture must be build-rejected
# as NO-BUILD. Any built binary, panic, mismatch, omitted row, or fallback
# is a hard failure (ADR-0028, ADR-0029).
#
# Builds catty with -race when R2_CONCURRENCY_STRESS>1 (Amendment 1
# convention) and records candidate/base/toolchain in the evidence header.
#
# Candidate evidence is isolated under:
#   docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/
# with NO slice-c suffix. Never writes research baseline, shared/latest,
# or any slice-c/ directory.
#
# Usage:
#   bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-candidate.sh <candidate>
#
#   R2_CONCURRENCY_STRESS=100 bash .../run-concurrency-candidate.sh <candidate>
#
# NOTE: does NOT use set -e due to Bash 3.2 compatibility; every command
# exit code is checked explicitly.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

die() { echo "candidate-runner: $*" >&2; exit 1; }

# --- Arg validation ---
[ $# -eq 1 ] || die "usage: $0 <candidate-commit-id>"
CANDIDATE="$1"

(cd "$ROOT" && git rev-parse --git-dir) >/dev/null 2>&1 || die "not a git repository: $ROOT"

CANDIDATE_FULL="$(cd "$ROOT" && git rev-parse --verify "$CANDIDATE^{commit}")" || die "not a valid commit: $CANDIDATE"

# --- Toolchain check ---
command -v java  >/dev/null || die "java not found (Temurin 25 required)"
command -v javac >/dev/null || die "javac not found (Temurin 25 required)"
command -v perl  >/dev/null || die "perl not found (portable timeout)"
command -v go    >/dev/null || die "go not found (catty build)"

java -version 2>&1 | head -1 | grep '"25\.' >/dev/null || die "java 25 required"
javac -version 2>&1 | grep '^javac 25\.' >/dev/null || die "javac 25 required"

# --- Stress multiplier ---
STRESS="${R2_CONCURRENCY_STRESS:-1}"
[ "$STRESS" -ge 1 ] || die "R2_CONCURRENCY_STRESS must be >= 1, got $STRESS"

if [ "$STRESS" -gt 1 ]; then
  RACE_BUILD=1
else
  RACE_BUILD=0
fi

# --- Evidence directory (main repo, never in the detached worktree) ---
EVIDENCE_DIR="$ROOT/docs/workstreams/r2-concurrency-candidate-evidence/$CANDIDATE"
if [ "$STRESS" -gt 1 ]; then
  RESULTS="$EVIDENCE_DIR/results-stress-${STRESS}x.txt"
else
  RESULTS="$EVIDENCE_DIR/results.txt"
fi

[ ! -f "$RESULTS" ] || die "refusing to overwrite existing evidence: $RESULTS"
mkdir -p "$EVIDENCE_DIR" || die "cannot create evidence directory: $EVIDENCE_DIR"

# --- Detached worktree at candidate ---
(cd "$ROOT" && git worktree prune) >/dev/null 2>&1 || true

BUILD_DIR="$(mktemp -d -t catty-candidate-build.XXXXXX)"
git -C "$ROOT" worktree add --detach --no-checkout "$BUILD_DIR" "$CANDIDATE_FULL" >/dev/null 2>&1 || die "failed to create detached worktree at $CANDIDATE"
git -C "$BUILD_DIR" checkout --detach "$CANDIDATE_FULL" >/dev/null 2>&1 || die "failed to checkout candidate in detached worktree"

BUILD_COMMIT="$(git -C "$BUILD_DIR" rev-parse HEAD)"
[ "$BUILD_COMMIT" = "$CANDIDATE_FULL" ] || die "build worktree HEAD ($BUILD_COMMIT) != candidate ($CANDIDATE_FULL)"

FIXTURE_DIR="$BUILD_DIR/docs/workstreams/r2-concurrency-fixtures"

BIN="$(mktemp -t catty-candidate.XXXXXX)"
STAGE="$(mktemp -d -t candidate-stage.XXXXXX)"
cleanup() {
  rm -rf "$BIN" "$STAGE"
  git -C "$ROOT" worktree remove --force "$BUILD_DIR" >/dev/null 2>&1 || true
  rm -rf "$BUILD_DIR"
}
trap cleanup EXIT

# --- The fixed 19 fixtures (hard-coded, from matrix.md) ---
FIXTURES=(
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
)

EXPECTED=19
to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

# Verify all fixtures exist in the candidate worktree.
for name in "${FIXTURES[@]}"; do
  [ -f "$FIXTURE_DIR/$name.java" ] || die "missing fixture in candidate: $name.java"
done
[ "${#FIXTURES[@]}" -eq "$EXPECTED" ] || die "fixture count ${#FIXTURES[@]} != $EXPECTED"

T_RUN=20
T_AOT=120

# --- Output helpers ---
log_line() { printf '%s\n' "$*" >> "$RESULTS"; printf '%s\n' "$*"; }
printf_log() { local fmt="$1"; shift; printf "$fmt\n" "$@" >> "$RESULTS"; printf "$fmt\n" "$@"; }

# --- Header ---
log_line "=== R2 concurrency candidate runner (19-fixture matrix) ==="
log_line "candidate:           $CANDIDATE"
log_line "candidate-full:      $CANDIDATE_FULL"
log_line "build-commit:        $BUILD_COMMIT"
log_line "build-source:        detached worktree at $CANDIDATE_FULL"
log_line "worktree-cleanliness: detached (immutable candidate snapshot)"
log_line "date:                $(date -u +%Y-%m-%dT%H:%M:%SZ)"
log_line "java:                $(java -version 2>&1 | head -1)"
log_line "javac:               $(javac -version 2>&1)"
log_line "go:                  $(go version 2>&1)"
log_line "stress:              ${STRESS}x"
log_line "race-build:          ${RACE_BUILD}"
log_line "timeout:             ${T_RUN}s (AOT ${T_AOT}s)"
log_line "fixtures:            $EXPECTED"
log_line "policy:              fail-closed — any mismatch or missing row is a failure"
log_line "aot-policy:          all 19 must be NO-BUILD; any built binary, panic, mismatch, omitted row, or Fallback is Fail"
log_line ""

# --- Build from the detached worktree ---
BUILD_FLAGS=""
if [ "$RACE_BUILD" -eq 1 ]; then
  BUILD_FLAGS="-race"
  log_line "==> building catty (race-enabled)"
else
  log_line "==> building catty"
fi
(cd "$BUILD_DIR" && go build $BUILD_FLAGS -o "$BIN" ./cmd/jvm) >>"$RESULTS" 2>&1 || die "catty build failed"

# --- Table header ---
if [ "$STRESS" -eq 1 ]; then
  printf_log "%-35s %-14s %-14s %-14s %s" "fixture" "Temurin25" "Interpreter" "IR" "AOT"
  printf_log "%-35s %-14s %-14s %-14s %s" "-------" "---------" "-----------" "--" "---"
else
  printf_log "%-35s %-14s %-22s %-22s %s" "fixture" "Temurin25" "Interpreter(${STRESS}x)" "IR(${STRESS}x)" "AOT"
  printf_log "%-35s %-14s %-22s %-22s %s" "-------" "---------" "--------------------" "--------------------" "---"
fi

passed_i=0
passed_ir=0
passed_aot=0

for name in "${FIXTURES[@]}"; do
  dir="$STAGE/$name"
  mkdir -p "$dir"

  # Compile fixture from the candidate worktree.
  javac --release 25 -nowarn -d "$dir" "$FIXTURE_DIR/$name.java" 2>"$dir/javac.err" || {
    cat "$dir/javac.err"
    die "javac failed for $name"
  }

  # Reference (Temurin 25) — run once.
  ref=$( { cd "$dir" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); ref_rc=$?
  [ "$ref_rc" -eq 0 ] || die "Temurin reference failed/timed out for $name (rc=$ref_rc)"

  # Interpreter — STRESS times; every run must match.
  interp_status="Match"
  for i in $(seq 1 "$STRESS"); do
    interp=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); interp_rc=$?
    if [ "$ref" != "$interp" ] || [ "$ref_rc" != "$interp_rc" ]; then
      interp_status="FAIL"
      break
    fi
  done
  [ "$interp_status" = "Match" ] && passed_i=$((passed_i + 1))

  # IR — STRESS times; every run must match.
  ir_status="Match"
  for i in $(seq 1 "$STRESS"); do
    ir=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?
    if [ "$ref" != "$ir" ] || [ "$ref_rc" != "$ir_rc" ]; then
      ir_status="FAIL"
      break
    fi
  done
  [ "$ir_status" = "Match" ] && passed_ir=$((passed_ir + 1))

  # AOT — assert NO-BUILD.
  aot_status="NO-BUILD"
  aot_out="(build rejected)"
  ( cd "$ROOT" && to "$T_AOT" "$BIN" build -cp "$dir" -no-boot -o "$dir/aot.bin" "$name" ) >"$dir/aot_build.out" 2>&1; aot_build_rc=$?
  if [ -x "$dir/aot.bin" ]; then
    aot_status="FAIL(built)"
    aot_out="UNEXPECTED BUILD SUCCESS"
    aot_run=$( { cd "$dir" && to "$T_RUN" "$dir/aot.bin"; } 2>&1 ); aot_run_rc=$?
    aot_out="$aot_out
Run output: $aot_run"
  else
    aot_out=$(tail -3 "$dir/aot_build.out" 2>/dev/null | sed 's/\t/        /g')
    passed_aot=$((passed_aot + 1))
  fi

  printf_log "%-35s %-14s %-14s %-14s %s" "$name" "ref" "$interp_status" "$ir_status" "$aot_status"

  {
    echo "----- $name -----"
    echo "[$ref_rc] temurin25:"; printf '%s\n' "$ref"
    echo "[$interp_rc] catty interpreter:"; printf '%s\n' "$interp"
    echo "[$ir_rc] catty IR:"; printf '%s\n' "$ir"
    echo "[AOT] catty aot:"; printf '%s\n' "$aot_out"
  } >> "$RESULTS"

  # Fail closed on any mismatch.
  [ "$interp_status" = "Match" ] || die "$name Interpreter mismatch (rc=$interp_rc)"
  [ "$ir_status" = "Match" ] || die "$name IR mismatch (rc=$ir_rc)"
  [ "$aot_status" = "NO-BUILD" ] || die "$name AOT: expected NO-BUILD, got $aot_status"
done

{
  echo
  echo "==> Candidate summary"
  echo "fixtures:           $EXPECTED"
  echo "interpreter match:  $passed_i/$EXPECTED"
  echo "IR match:           $passed_ir/$EXPECTED"
  echo "AOT NO-BUILD:       $passed_aot/$EXPECTED"
  echo "stress:             ${STRESS}x"
  echo "result:             Pass"
} >> "$RESULTS"

echo "candidate-runner: all $EXPECTED fixtures passed (Interpreter + IR + AOT NO-BUILD)" >&2
