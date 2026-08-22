#!/usr/bin/env bash
# Mechanical validation gate (protocol §16).
# Fast-fail structural checks for the repository. Code-level gates
# (go vet/build/test) activate automatically once go.mod exists.
set -euo pipefail
cd "$(dirname "$0")/.."

fail=0
err() { echo "FAIL: $*"; fail=1; }

required=(
  AGENTS.md
  README.md
  Makefile
  .gitignore
  docs/protocol/autonomous-development-protocol.md
  docs/vision.md
  docs/status.md
  docs/architecture.md
  docs/quality/testing.md
  docs/quality/performance.md
  docs/quality/security.md
  docs/quality/architecture-rules.md
  docs/quality/escalation-policy.md
  docs/debt/register.md
  docs/research/TEMPLATE.md
  docs/specifications/README.md
  docs/specifications/deviation-ledger.md
  docs/specifications/task-packet-template.md
  docs/briefs/README.md
)

for f in "${required[@]}"; do
  [[ -f "$f" ]] || err "missing required file: $f"
done

# Active plans must exist (at least one while project is running).
ls docs/plans/active/P-*.md >/dev/null 2>&1 || err "no active plan in docs/plans/active/"

# ADR filename convention.
for f in docs/decisions/ADR-*.md; do
  [[ -e "$f" ]] || { err "no ADR found"; break; }
  [[ "$f" =~ ^docs/decisions/ADR-[0-9]{4}\.md$ ]] || err "bad ADR filename: $f"
done
# Duplicate ADR numbers.
dups=$(ls docs/decisions | sed -E 's/^ADR-([0-9]{4})\.md$/\1/' | sort | uniq -d || true)
[[ -z "$dups" ]] || err "duplicate ADR numbers: $dups"

# Protocol file integrity marker.
grep -q "AI 自治研发协议" docs/protocol/autonomous-development-protocol.md \
  || err "protocol file missing version marker"

# Unexplained TODO/FIXME in code files (markdown exempt; tools/ self-exempt).
while IFS= read -r -d '' f; do
  n=$(grep -E '(TODO|FIXME)' "$f" | grep -cvE '(TODO|FIXME)\(' || true)
  (( n > 0 )) && err "unexplained TODO/FIXME ($n) in $f — tag as TODO(P-xxxx|ADR-xxxx|DEBT-xxxx)"
done < <(find . -path ./.git -prune -o -path ./tools -prune -o -type f \
        \( -name '*.go' -o -name '*.sh' -o -name '*.js' -o -name '*.ts' \) -print0)

# Code gates, active once Go sources exist.
if [[ -n "$(find . -path ./.git -prune -o -name '*.go' -print -quit 2>/dev/null)" ]]; then
  go vet ./... || err "go vet failed"
  go build ./... || err "go build failed"
  go test ./... > /tmp/catty_go_test.out 2>&1 \
    || { err "go test failed (tail below)"; tail -40 /tmp/catty_go_test.out; }
fi

if (( fail )); then
  echo "check: FAILED"
  exit 1
fi
echo "check: OK"
