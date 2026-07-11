#!/usr/bin/env bash
# E2E verification: compile each .java fixture with javac, run it through both
# the real `java` and through catty, and diff stdout. A fixture passes only when
# catty's output is byte-identical to java's.
#
# Usage: ./tests/run.sh
set -u

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
FIX="$ROOT/tests/fixtures"

BINARY="$(mktemp -t catty.XXXXXX)"
trap 'rm -f "$BINARY"; rm -f "$FIX"/*.class' EXIT

echo "==> building catty"
go build -o "$BINARY" "$ROOT/cmd/jvm" || exit 1

echo "==> compiling fixtures"
javac -source 8 -target 8 -nowarn -d "$FIX" "$FIX"/*.java 2>/dev/null || {
    echo "javac failed"; exit 1; }

# Each main class to exercise. (Point has no main; it's a helper for OOPDemo.)
MAIN_CLASSES=(HelloWorld Fibonacci Factorial ArraySum OOPDemo StaticFields SwitchDemo)

pass=0; fail=0
for cls in "${MAIN_CLASSES[@]}"; do
    java_out=$(cd "$FIX" && java "$cls" 2>&1)
    catty_out=$(cd "$FIX" && "$BINARY" -cp . "$cls" 2>&1)
    if [ "$java_out" = "$catty_out" ]; then
        echo "PASS  $cls"
        pass=$((pass+1))
    else
        echo "FAIL  $cls"
        echo "  --- java ---"; echo "$java_out" | sed 's/^/  /'
        echo "  --- catty ---"; echo "$catty_out" | sed 's/^/  /'
        fail=$((fail+1))
    fi
done

echo
echo "==> $pass passed, $fail failed"
[ "$fail" -eq 0 ]
