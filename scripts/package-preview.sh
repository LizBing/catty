#!/bin/sh
set -eu

if [ "$#" -ne 2 ]; then
  printf 'usage: %s VERSION OUTPUT_DIRECTORY\n' "$0" >&2
  exit 2
fi

version=$1
case "$version" in
  ''|[!0-9]*|*[!0-9A-Za-z.-]*)
    printf 'Invalid preview version: %s\n' "$version" >&2
    exit 2
    ;;
esac

mkdir -p "$2"
output=$(CDPATH= cd -P "$2" && pwd -P)
repo=$(CDPATH= cd -P "$(dirname "$0")/.." && pwd -P)
cd "$repo"

release="catty-$version-darwin-arm64"
work=$(mktemp -d "${TMPDIR:-/tmp}/catty-preview.XXXXXX")
trap 'rm -rf "$work"' EXIT
root="$work/$release"
mkdir -p "$root/bin" "$root/libexec"

GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 GOWORK=off GOTOOLCHAIN=local \
  go build -trimpath -buildvcs=false -o "$root/libexec/catty" ./cmd/catty

cp go.mod README.md "$root/"
if [ -f LICENSE ]; then
  cp LICENSE "$root/"
fi
for pkg in rt plugin internal/admission internal/classfile internal/hierarchy internal/verify; do
  mkdir -p "$root/$pkg"
  for source in "$pkg"/*.go; do
    case "$source" in
      *_test.go) ;;
      *) cp "$source" "$root/$pkg/" ;;
    esac
  done
done

cat > "$root/bin/catty" <<'EOF'
#!/bin/sh
root=$(CDPATH= cd -P "$(dirname "$0")/.." && pwd -P) || exit 1
CATTY_MODULE_ROOT="$root" exec "$root/libexec/catty" "$@"
EOF
chmod 755 "$root/bin/catty"

(
  cd "$work"
  tar -czf "$release.tar.gz" "$release"
  shasum -a 256 "$release.tar.gz" > "$release.tar.gz.sha256"
)
mv "$work/$release.tar.gz" "$work/$release.tar.gz.sha256" "$output"
printf '%s\n' "$output/$release.tar.gz"
