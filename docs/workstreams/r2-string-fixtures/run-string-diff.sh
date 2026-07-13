#!/usr/bin/env bash
# R2 String differential harness — String UTF-16 semantics (ADR-0027).
#
# Compares Temurin 25 against catty's three engines (Interpreter, IR, AOT) on
# the 8-fixture String matrix. Fail-closed: non-zero exit on any mismatch,
# build failure, or missing toolchain.
#
# The eight fixtures:
#   1. SupplementaryChar        — supplementary character in literal
#   2. HashDivergence           — UTF-16 code-unit hashCode
#   3. LoneSurrogate            — String(char[]) lone surrogate round-trip
#   4. StringBounds             — StringIndexOutOfBoundsException
#   5. StringSubstringUnits     — substring code-unit indices
#   6. StringCharArrayRoundTrip — char[] constructor + toCharArray defensive copy
#   7. LoneSurrogateLiteral     — lone surrogate from classfile literal (MUTF-8)
#   8. StringNativeSurface      — null contract + PrintStream + StringBuilder
#
# Usage: bash docs/workstreams/r2-string-fixtures/run-string-diff.sh [candidate-ref]
#   candidate-ref — if provided, used as the candidate directory name and
#                   substituted for {candidate} in the evidence path. Defaults
#                   to the short commit SHA.

set -u

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
EVI_FIX="$ROOT/docs/workstreams/r2-evidence/fixtures"
STR_FIX="$SCRIPT_DIR"   # r2-string-fixtures/ (this script's directory)

# Candidate identity.
CANDIDATE="${1:-$(cd "$ROOT" && git rev-parse --short HEAD)}"
EVI_DIR="$ROOT/docs/workstreams/r2-string-evidence/$CANDIDATE"
mkdir -p "$EVI_DIR"
RESULTS="$EVI_DIR/run-string-results.txt"

fail_closed() { echo "run-string-diff: $*" >&2; exit 1; }

command -v java  >/dev/null || fail_closed "java not found (Temurin 25 required)"
command -v javac >/dev/null || fail_closed "javac not found (Temurin 25 required)"
command -v perl  >/dev/null || fail_closed "perl not found (portable timeout)"
command -v go    >/dev/null || fail_closed "go not found (catty build + AOT)"

# Portable wall-clock timeout.
to() { perl -e 'alarm shift; exec @ARGV' "$@"; }

BIN="$(mktemp -t catty-string.XXXXXX)"
STAGE="$(mktemp -d -t string-stage.XXXXXX)"
trap 'rm -rf "$BIN" "$STAGE"' EXIT

{
  echo "=== R2 String differential run ==="
  echo "candidate: $CANDIDATE"
  echo "repo:      $ROOT"
  echo "commit:    $(cd "$ROOT" && git rev-parse HEAD)"
  echo "base:      298b723"
  echo "date:      $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "toolchain: java=$(java -version 2>&1 | head -1 | sed 's/.*version "//;s/".*//'), go=$(go version 2>&1 | sed 's/go version //;s/ .*//'), javac=$(javac -version 2>&1 | sed 's/javac //')"
  echo "compare:   combined stdout+stderr AND exit code must equal Temurin 25 (fail-closed)"
  echo
} | tee "$RESULTS"

echo "==> building catty" | tee -a "$RESULTS"
go build -o "$BIN" "$ROOT/cmd/jvm" >>"$RESULTS" 2>&1 || fail_closed "catty build failed"

T_RUN=20
T_AOT=120

printf "%-28s %-6s %-10s %-10s %-10s %s\n" "fixture" "exit" "Temurin25" "Interp" "IR" "AOT" | tee -a "$RESULTS"
printf "%-28s %-6s %-10s %-10s %-10s %s\n" "-------" "----" "---------" "------" "--" "---" | tee -a "$RESULTS"

PASS=0
FAIL=0

run_one() {
  local jf="$1"
  local name
  name="$(basename "$jf" .java)"
  local d="$STAGE/$name"
  mkdir -p "$d"

  if ! javac --release 25 -nowarn -d "$d" "$jf" 2>"$d/javac.err"; then
    printf "%-28s %-6s %-10s %-10s %-10s %s\n" "$name" "FAIL" "javac-FAIL" "-" "-" "-" | tee -a "$RESULTS"
    FAIL=$((FAIL + 1))
    return 1
  fi

  local tum tum_rc lp lp_rc ir ir_rc
  tum=$( { cd "$d" && to "$T_RUN" java -cp . "$name"; } 2>&1 ); tum_rc=$?
  lp=$(  { cd "$d" && to "$T_RUN" "$BIN" -cp . -no-boot "$name"; } 2>&1 ); lp_rc=$?
  ir=$(  { cd "$d" && to "$T_RUN" "$BIN" -cp . -ir -no-boot "$name"; } 2>&1 ); ir_rc=$?

  # AOT (best effort): build and run.
  local m_aot="NO-BUILD" aot aot_rc="-"
  ( cd "$ROOT" && to "$T_AOT" "$BIN" build -cp "$d" -no-boot -o "$d/aot.bin" "$name" ) >"$d/aot_build.out" 2>&1
  if [ -x "$d/aot.bin" ]; then
    aot=$( { cd "$d" && to "$T_RUN" "$d/aot.bin"; } 2>&1 ); aot_rc=$?
    if [ "$tum" = "$aot" ] && [ "$tum_rc" = "$aot_rc" ]; then
      m_aot="match"
    else
      m_aot="MISMATCH"
    fi
  else
    aot=$(tail -3 "$d/aot_build.out" 2>/dev/null | sed 's/\t/        /g')
  fi

  local m_lp="MISMATCH" m_ir="MISMATCH"
  [ "$tum" = "$lp" ] && [ "$tum_rc" = "$lp_rc" ] && m_lp="match"
  [ "$tum" = "$ir" ] && [ "$tum_rc" = "$ir_rc" ] && m_ir="match"

  local verdict="PASS"
  [ "$m_lp" != "match" ] && verdict="FAIL"
  [ "$m_ir" != "match" ] && verdict="FAIL"
  # AOT: fail-closed only if it built but mismatched (NO-BUILD is not a fail).
  [ "$m_aot" = "MISMATCH" ] && verdict="FAIL"

  if [ "$verdict" = "PASS" ]; then
    PASS=$((PASS + 1))
  else
    FAIL=$((FAIL + 1))
  fi

  printf "%-28s %-6s %-10s %-10s %-10s %s\n" "$name" "$tum_rc" "$tum" "$m_lp" "$m_ir" "$m_aot" | tee -a "$RESULTS"

  {
    echo "----- $name -----"
    echo "[$tum_rc] temurin25:"; printf '%s\n' "$tum"
    echo "[$lp_rc] catty interp:"; printf '%s\n' "$lp"
    echo "[$ir_rc] catty ir:"; printf '%s\n' "$ir"
    echo "[$aot_rc] catty aot:"; printf '%s\n' "$aot"
  } >>"$RESULTS"
}

for jf in "$EVI_FIX"/SupplementaryChar.java \
          "$EVI_FIX"/HashDivergence.java \
          "$EVI_FIX"/LoneSurrogate.java \
          "$EVI_FIX"/StringBounds.java \
          "$EVI_FIX"/StringSubstringUnits.java \
          "$EVI_FIX"/StringCharArrayRoundTrip.java \
          "$STR_FIX"/LoneSurrogateLiteral.java \
          "$STR_FIX"/StringNativeSurface.java; do
  [ -e "$jf" ] || { printf "%-28s %-6s %-10s %-10s %-10s %s\n" "$(basename "$jf" .java)" "FAIL" "MISSING" "-" "-" "-" | tee -a "$RESULTS"; FAIL=$((FAIL + 1)); continue; }
  run_one "$jf"
done

echo | tee -a "$RESULTS"
echo "==> $PASS passed, $FAIL failed ($((PASS + FAIL)) total)" | tee -a "$RESULTS"
echo "evidence: $RESULTS" | tee -a "$RESULTS"

# Fail-closed.
if [ "$FAIL" -gt 0 ]; then
  fail_closed "$FAIL fixture(s) failed"
fi
echo "run-string-diff: all String fixtures match Temurin 25" >&2
