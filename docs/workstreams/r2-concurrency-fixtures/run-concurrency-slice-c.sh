#!/usr/bin/env bash
# Slice C fail-closed conformance runner (ADR-0029).
#
# Runs the 11 Slice C concurrency fixtures against Temurin 25 (reference),
# catty Interpreter, and catty IR. Every fixture must match the reference
# in combined stdout+stderr and exit code. Any mismatch, timeout, missing
# fixture, or build failure is a hard failure.
#
# The candidate commit is resolved to a full SHA and checked out in a
# temporary detached worktree so the build and fixtures are sourced from
# exactly that commit — the caller's working state is irrelevant.
#
# Usage:
#   bash docs/workstreams/r2-concurrency-fixtures/run-concurrency-slice-c.sh <candidate>
#
#   R2_CONCURRENCY_STRESS=20 bash <script> <candidate>
#     Runs each fixture that many times per engine.  Default is 1.
#
# Output:
#   docs/workstreams/r2-concurrency-candidate-evidence/<candidate>/slice-c/
#     results.txt              — 1× run
#     results-stress-<N>x.txt  — stress run (when STRESS > 1)
#
# Guard: refuses to overwrite an existing evidence file.
set -eu

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

fail_closed() { echo "slice-c-runner: $*" >&2; exit 1; }

# --- Arg validation ---
[ $# -eq 1 ] || fail_closed "usage: $0 <candidate-commit-id>"
CANDIDATE="$1"

# Verify we are in a real Git checkout.
(cd "$ROOT" && git rev-parse --git-dir) >/dev/null 2>&1 \
  || fail_closed "not a git repository: $ROOT"

# Resolve candidate to a full, immutable SHA.
CANDIDATE_FULL="$(cd "$ROOT" && git rev-parse --verify "$CANDIDATE^{commit}")" \
  || fail_closed "not a valid commit: $CANDIDATE"

# --- Toolchain check ---
command -v java  >/dev/null || fail_closed "java not found (Temurin 25 required)"
command -v javac >/dev/null || fail_closed "javac not found (Temurin 25 required)"
command -v perl  >/dev/null || fail_closed "perl not found (portable timeout)"
command -v go    >/dev/null || fail_closed "go not found (catty build)"

java -version 2>&1 | head -1 | grep '"25\.' >/dev/null || fail_closed "java 25 required"
javac -version 2>&1 | grep '^javac 25\.' >/dev/null || fail_closed "javac 25 required"

# --- Stress multiplier ---
STRESS="${R2_CONCURRENCY_STRESS:-1}"
[ "$STRESS" -ge 1 ] || fail_closed "R2_CONCURRENCY_STRESS must be >= 1, got $STRESS"

# --- Evidence directory (main repo, never in the detached worktree) ---
EVIDENCE_DIR="$ROOT/docs/workstreams/r2-concurrency-candidate-evidence/$CANDIDATE/slice-c"
if [ "$STRESS" -gt 1 ]; then
  RESULTS="$EVIDENCE_DIR/results-stress-${STRESS}x.txt"
else
  RESULTS="$EVIDENCE_DIR/results.txt"
fi

if [ -f "$RESULTS" ]; then
  fail_closed "refusing to overwrite existing evidence: $RESULTS"
fi
mkdir -p "$EVIDENCE_DIR" || fail_closed "cannot create evidence directory: $EVIDENCE_DIR"

# --- Detached worktree at candidate ---
# Prune any stale worktree metadata from previous aborted runs.
(cd "$ROOT" && git worktree prune) >/dev/null 2>&1 || true

BUILD_DIR="$(mktemp -d -t catty-slice-c-build.XXXXXX)"
git -C "$ROOT" worktree add --detach --no-checkout "$BUILD_DIR" "$CANDIDATE_FULL" >/dev/null 2>&1 \
  || fail_closed "failed to create detached worktree at $CANDIDATE"
# Checkout the full tree now that the worktree exists.
git -C "$BUILD_DIR" checkout --detach "$CANDIDATE_FULL" >/dev/null 2>&1 \
  || fail_closed "failed to checkout candidate in detached worktree"

BUILD_COMMIT="$(git -C "$BUILD_DIR" rev-parse HEAD)"
[ "$BUILD_COMMIT" = "$CANDIDATE_FULL" ] \
  || fail_closed "build worktree HEAD ($BUILD_COMMIT) != candidate ($CANDIDATE_FULL)"

FIXTURE_DIR="$BUILD_DIR/docs/workstreams/r2-concurrency-fixtures"

# Cleanup: remove detached worktree + temp files.
BIN="$(mktemp -t catty-slice-c.XXXXXX)"
STAGE="$(mktemp -d -t slice-c-stage.XXXXXX)"
cleanup() {
  rm -rf "$BIN" "$STAGE"
  git -C "$ROOT" worktree remove --force "$BUILD_DIR" >/dev/null 2>&1 || true
  rm -rf "$BUILD_DIR"
}
trap cleanup EXIT

# --- Slice C fixture list (11) ---
FIXTURES="
SynchronizedBlocks
SynchronizedMethods
MonitorNull
MonitorOwnership
WaitNotify
NotifyAll
InterruptWait
InterruptStatus
InterruptSleep
InterruptJoin
ProducerConsumer
"

EXPECTED=11
to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

# Verify all fixtures exist in the candidate worktree.
count=0
for name in $FIXTURES; do
  [ -f "$FIXTURE_DIR/$name.java" ] || fail_closed "missing fixture in candidate: $name.java"
  count=$((count + 1))
done
[ "$count" -eq "$EXPECTED" ] || fail_closed "fixture count $count != $EXPECTED"

T_RUN=20

# --- Header ---
{
  echo "=== Slice C runner ==="
  echo "candidate:           $CANDIDATE"
  echo "candidate-full:      $CANDIDATE_FULL"
  echo "build-commit:        $BUILD_COMMIT"
  echo "build-source:        detached worktree at $CANDIDATE_FULL"
  echo "worktree-cleanliness: detached (immutable candidate snapshot)"
  echo "date:                $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "java:                $(java -version 2>&1 | head -1)"
  echo "javac:               $(javac -version 2>&1)"
  echo "go:                  $(go version 2>&1)"
  echo "stress:              ${STRESS}x"
  echo "timeout:             ${T_RUN}s"
  echo "fixtures:            $EXPECTED"
  echo "policy:              fail-closed — any mismatch or missing row is a failure"
  echo
} | tee "$RESULTS"

# --- Build from the detached worktree ---
echo "==> building catty" | tee -a "$RESULTS"
go build -o "$BIN" "$BUILD_DIR/cmd/jvm" >>"$RESULTS" 2>&1 \
  || fail_closed "catty build failed"

if [ "$STRESS" -eq 1 ]; then
  printf "%-30s %-14s %-14s %-14s\n" "fixture" "Temurin25" "Interpreter" "IR" | tee -a "$RESULTS"
  printf "%-30s %-14s %-14s %-14s\n" "-------" "---------" "-----------" "--" | tee -a "$RESULTS"
else
  printf "%-30s %-14s %-14s %-14s\n" "fixture" "Temurin25" "Interpreter(${STRESS}x)" "IR(${STRESS}x)" | tee -a "$RESULTS"
  printf "%-30s %-14s %-14s %-14s\n" "-------" "---------" "-----------------" "-----------" | tee -a "$RESULTS"
fi

passed_i=0
passed_ir=0

for name in $FIXTURES; do
  dir="$STAGE/$name"
  mkdir -p "$dir"

  # Compile fixture from the candidate worktree.
  if ! javac --release 25 -nowarn -d "$dir" "$FIXTURE_DIR/$name.java" 2>"$dir/javac.err"; then
    cat "$dir/javac.err"
    fail_closed "javac failed for $name"
  fi

  # Reference (Temurin 25) — run once.
  ref=$( { cd "$dir" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); ref_rc=$?
  [ "$ref_rc" -eq 0 ] || fail_closed "Temurin reference failed/timed out for $name (rc=$ref_rc)"

  # Interpreter — run STRESS times; every run must match.
  interp_status="Match"
  for i in $(seq 1 "$STRESS"); do
    interp=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); interp_rc=$?
    if [ "$ref" != "$interp" ] || [ "$ref_rc" != "$interp_rc" ]; then
      interp_status="FAIL"
      break
    fi
  done
  [ "$interp_status" = "Match" ] && passed_i=$((passed_i + 1))

  # IR — run STRESS times; every run must match.
  ir_status="Match"
  for i in $(seq 1 "$STRESS"); do
    ir=$( { cd "$dir" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?
    if [ "$ref" != "$ir" ] || [ "$ref_rc" != "$ir_rc" ]; then
      ir_status="FAIL"
      break
    fi
  done
  [ "$ir_status" = "Match" ] && passed_ir=$((passed_ir + 1))

  printf "%-30s %-14s %-14s %-14s\n" "$name" "ref" "$interp_status" "$ir_status" | tee -a "$RESULTS"

  {
    echo "----- $name -----"
    echo "[0] temurin25:"; printf '%s\n' "$ref"
    echo "[$interp_rc] catty interpreter:"; printf '%s\n' "$interp"
    echo "[$ir_rc] catty IR:"; printf '%s\n' "$ir"
  } >>"$RESULTS"

  # Fail closed on any mismatch.
  [ "$interp_status" = "Match" ] || fail_closed "$name Interpreter mismatch (rc=$interp_rc)"
  [ "$ir_status" = "Match" ] || fail_closed "$name IR mismatch (rc=$ir_rc)"
done

{
  echo
  echo "==> Slice C summary"
  echo "fixtures:          $EXPECTED"
  echo "interpreter match: $passed_i/$EXPECTED"
  echo "IR match:          $passed_ir/$EXPECTED"
  echo "stress:            ${STRESS}x"
  echo "result:            Pass"
} | tee -a "$RESULTS"

echo "slice-c-runner: all $EXPECTED fixtures passed (Interpreter + IR)" >&2
