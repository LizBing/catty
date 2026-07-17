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
# When R2_CONCURRENCY_STRESS>1, fixtures are processed concurrently:
# each fixture runs as a background subshell executing its full stress
# loop independently. Concurrency defaults to 4 (configurable via
# R2_STRESS_CONCURRENCY). Per-fixture outputs are written to independent
# temp files and merged in fixture-list order after completion.
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
#   R2_CONCURRENCY_STRESS=100 R2_STRESS_CONCURRENCY=8 bash .../run-concurrency-candidate.sh <candidate>
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

# --- Concurrency control (stress mode only) ---
CONCURRENCY="${R2_STRESS_CONCURRENCY:-4}"
[ "$CONCURRENCY" -ge 1 ] || die "R2_STRESS_CONCURRENCY must be >= 1, got $CONCURRENCY"

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

# FIFO for concurrent job-server token passing (stress mode only)
JOBSERVER_FIFO=""
CONC_OUT=""

cleanup() {
  # Close job server fd if open
  exec 3>&- 2>/dev/null || true
  exec 4>&- 2>/dev/null || true
  if [ -n "$JOBSERVER_FIFO" ] && [ -p "$JOBSERVER_FIFO" ]; then
    rm -f "$JOBSERVER_FIFO"
  fi
  if [ -n "$CONC_OUT" ] && [ -d "$CONC_OUT" ]; then
    rm -rf "$CONC_OUT"
  fi
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
printf_log() { local fmt="$1"; shift; printf "$fmt\n" "$@" >> "$RESULTS"; printf "$fmt\n" "$*"; }

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
if [ "$STRESS" -gt 1 ]; then
  log_line "concurrency:         ${CONCURRENCY} (max concurrent fixtures)"
fi
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

# --- Table header and row format ---
if [ "$STRESS" -eq 1 ]; then
  printf_log "%-35s %-14s %-14s %-14s %s" "fixture" "Temurin25" "Interpreter" "IR" "AOT"
  printf_log "%-35s %-14s %-14s %-14s %s" "-------" "---------" "-----------" "--" "---"
  ROW_FMT="%-35s %-14s %-14s %-14s %s"
else
  printf_log "%-35s %-14s %-22s %-22s %s" "fixture" "Temurin25" "Interpreter(${STRESS}x)" "IR(${STRESS}x)" "AOT"
  printf_log "%-35s %-14s %-22s %-22s %s" "-------" "---------" "--------------------" "--------------------" "---"
  ROW_FMT="%-35s %-14s %-22s %-22s %s"
fi

# --- Per-fixture processing function ---
# Shared by both sequential (1x) and concurrent (stress) modes.
# Arguments: $1=name $2=output-file (empty string = write to $RESULTS directly)
# Uses outer variables: STAGE, FIXTURE_DIR, BIN, ROOT, T_RUN, T_AOT, STRESS
# Returns 0 on success, 1 on failure.
process_fixture() {
  local name="$1"
  local outfile="$2"
  local dir="$STAGE/$name"
  mkdir -p "$dir"

  # Compile fixture from the candidate worktree.
  javac --release 25 -nowarn -d "$dir" "$FIXTURE_DIR/$name.java" 2>"$dir/javac.err"
  if [ $? -ne 0 ]; then
    echo "javac failed for $name" >> "$outfile"
    return 1
  fi

  # Reference (Temurin 25) — run once.
  local ref ref_rc
  ref=$( { cd "$dir" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); ref_rc=$?
  if [ "$ref_rc" -ne 0 ]; then
    echo "Temurin reference failed/timed out for $name (rc=$ref_rc)" >> "$outfile"
    return 1
  fi

  # Interpreter — STRESS times; every run must match.
  local interp_status="Match" interp interp_rc
  local i
  for i in $(seq 1 "$STRESS"); do
    interp=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); interp_rc=$?
    if [ "$ref" != "$interp" ] || [ "$ref_rc" != "$interp_rc" ]; then
      interp_status="FAIL"
      break
    fi
  done

  # IR — STRESS times; every run must match.
  local ir_status="Match" ir ir_rc
  for i in $(seq 1 "$STRESS"); do
    ir=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?
    if [ "$ref" != "$ir" ] || [ "$ref_rc" != "$ir_rc" ]; then
      ir_status="FAIL"
      break
    fi
  done

  # AOT — assert NO-BUILD.
  local aot_status="NO-BUILD" aot_out="(build rejected)" aot_build_rc
  ( cd "$ROOT" && to "$T_AOT" "$BIN" build -cp "$dir" -no-boot -o "$dir/aot.bin" "$name" ) >"$dir/aot_build.out" 2>&1; aot_build_rc=$?
  if [ -x "$dir/aot.bin" ]; then
    aot_status="FAIL(built)"
    local aot_run aot_run_rc
    aot_out="UNEXPECTED BUILD SUCCESS"
    aot_run=$( { cd "$dir" && to "$T_RUN" "$dir/aot.bin"; } 2>&1 ); aot_run_rc=$?
    aot_out="$aot_out
Run output: $aot_run"
  else
    aot_out=$(tail -3 "$dir/aot_build.out" 2>/dev/null | sed 's/\t/        /g')
  fi

  # Write table row and detail block to output file.
  {
    printf "$ROW_FMT\n" "$name" "ref" "$interp_status" "$ir_status" "$aot_status"
    echo "----- $name -----"
    echo "[$ref_rc] temurin25:"; printf '%s\n' "$ref"
    echo "[$interp_rc] catty interpreter:"; printf '%s\n' "$interp"
    echo "[$ir_rc] catty IR:"; printf '%s\n' "$ir"
    echo "[AOT] catty aot:"; printf '%s\n' "$aot_out"
  } >> "$outfile"

  # Write status for summary counting.
  echo "$interp_status" > "$STAGE/$name.interp_status"
  echo "$ir_status"     > "$STAGE/$name.ir_status"
  echo "$aot_status"    > "$STAGE/$name.aot_status"

  # Fail closed on any mismatch.
  [ "$interp_status" = "Match" ] || { echo "$name Interpreter mismatch (rc=$interp_rc)" >&2; return 1; }
  [ "$ir_status" = "Match" ]     || { echo "$name IR mismatch (rc=$ir_rc)" >&2; return 1; }
  [ "$aot_status" = "NO-BUILD" ] || { echo "$name AOT: expected NO-BUILD, got $aot_status" >&2; return 1; }
  return 0
}

# --- Main fixture processing ---
passed_i=0
passed_ir=0
passed_aot=0

if [ "$STRESS" -eq 1 ]; then
  # === Sequential mode (1x): unchanged behavior ===
  for name in "${FIXTURES[@]}"; do
    process_fixture "$name" "$RESULTS" || die "fixture $name failed"
  done
else
  # === Concurrent mode (stress>1): FIFO-based job server ===
  CONC_OUT="$(mktemp -d -t candidate-conc.XXXXXX)"
  JOBSERVER_FIFO="$CONC_OUT/.jobserver"
  mkfifo "$JOBSERVER_FIFO" || die "mkfifo failed for $JOBSERVER_FIFO"

  # Seed the FIFO with CONCURRENCY tokens.
  exec 3<>"$JOBSERVER_FIFO"
  i=0
  while [ "$i" -lt "$CONCURRENCY" ]; do
    printf '\n' >&3
    i=$((i + 1))
  done

  # Also open write fd for token return from subshells.
  exec 4<>"$JOBSERVER_FIFO"

  # Launch all fixtures as background subshells.
  # Each acquires a FIFO token (blocks when all slots are taken),
  # runs process_fixture, then returns the token.
  pids=""
  for name in "${FIXTURES[@]}"; do
    (
      # Acquire token (blocks if no slot available).
      read -r _ <&3
      process_fixture "$name" "$CONC_OUT/$name.out"
      rc=$?
      # Return token so another waiting fixture can proceed.
      printf '\n' >&4
      exit $rc
    ) &
    pids="$pids $!"
  done

  # Wait for all background jobs; track failures.
  failed=0
  for pid in $pids; do
    wait "$pid" || failed=$((failed + 1))
  done

  # Close job server fds.
  exec 3>&- 2>/dev/null || true
  exec 4>&- 2>/dev/null || true

  # Merge per-fixture outputs in fixture-list order (deterministic table layout).
  for name in "${FIXTURES[@]}"; do
    if [ -f "$CONC_OUT/$name.out" ]; then
      cat "$CONC_OUT/$name.out" >> "$RESULTS"
    else
      echo "MISSING  $name  (no output — fixture may have crashed before writing)" >> "$RESULTS"
      failed=$((failed + 1))
    fi
  done

  # Count passes from per-fixture status files.
  for name in "${FIXTURES[@]}"; do
    is="$(cat "$STAGE/$name.interp_status" 2>/dev/null)"
    irs="$(cat "$STAGE/$name.ir_status" 2>/dev/null)"
    as="$(cat "$STAGE/$name.aot_status" 2>/dev/null)"
    [ "$is" = "Match" ] && passed_i=$((passed_i + 1))
    [ "$irs" = "Match" ] && passed_ir=$((passed_ir + 1))
    [ "$as" = "NO-BUILD" ] && passed_aot=$((passed_aot + 1))
  done

  [ "$failed" -eq 0 ] || die "$failed fixture(s) failed (see detail blocks above)"
fi

# --- Summary ---
{
  echo
  echo "==> Candidate summary"
  echo "fixtures:           $EXPECTED"
  echo "interpreter match:  $passed_i/$EXPECTED"
  echo "IR match:           $passed_ir/$EXPECTED"
  echo "AOT NO-BUILD:       $passed_aot/$EXPECTED"
  echo "stress:             ${STRESS}x"
  if [ "$STRESS" -gt 1 ]; then
    echo "concurrency:        ${CONCURRENCY}"
  fi
  echo "result:             Pass"
} >> "$RESULTS"

echo "candidate-runner: all $EXPECTED fixtures passed (Interpreter + IR + AOT NO-BUILD)" >&2
