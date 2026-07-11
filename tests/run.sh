#!/usr/bin/env bash
# E2E verification: compile each .java fixture with javac, run it through the
# real `java` and through BOTH catty execution engines (the tree-walking
# interpreter and the -ir lowered executor), and diff stdout. A fixture passes
# only when all three produce byte-identical output.
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
    loop_out=$(cd "$FIX" && "$BINARY" -cp . "$cls" 2>&1)
    ir_out=$(cd "$FIX" && "$BINARY" -cp . -ir "$cls" 2>&1)
    if [ "$java_out" = "$loop_out" ] && [ "$java_out" = "$ir_out" ]; then
        echo "PASS  $cls"
        pass=$((pass+1))
    else
        echo "FAIL  $cls"
        echo "  --- java ---"; echo "$java_out" | sed 's/^/  /'
        echo "  --- loop ---"; echo "$loop_out" | sed 's/^/  /'
        echo "  --- ir   ---"; echo "$ir_out" | sed 's/^/  /'
        fail=$((fail+1))
    fi
done

echo
echo "==> $pass passed, $fail failed"
[ "$fail" -eq 0 ]
