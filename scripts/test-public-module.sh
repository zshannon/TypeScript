#!/usr/bin/env bash

set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
public_source="$repo_root/public"
consumer_source="$repo_root/scripts/testdata/public-consumer"
test_root=$(mktemp -d "${TMPDIR:-/tmp}/typescript-public-module.XXXXXX")
trap 'rm -rf "$test_root"' EXIT

if [[ ! -f "$public_source/go.mod" ]]; then
  echo "FAIL: generated public module is missing; run go run ./scripts/export-go.go" >&2
  exit 1
fi
public_module=$(awk '$1 == "module" { print $2; exit }' "$public_source/go.mod")
if [[ ! "$public_module" =~ ^github\.com/zshannon/TypeScript/public/v([1-9][0-9]*)$ ]]; then
  echo "FAIL: generated module path is not a major-versioned public module: $public_module" >&2
  exit 1
fi
public_major=${BASH_REMATCH[1]}
source_version=$(go -C "$repo_root" run ./scripts/export-go.go --version)
source_major=${source_version#v}
source_major=${source_major%%.*}
if [[ "$source_major" != "$public_major" ]]; then
  echo "FAIL: generated module major v$public_major does not match source version $source_version" >&2
  exit 1
fi

mkdir "$test_root/consumer"
cp "$consumer_source/go.mod" "$test_root/consumer/go.mod"
sed "s|PUBLIC_MODULE_PATH|$public_module|g" \
  "$consumer_source/main.go" >"$test_root/consumer/main.go"

if [[ -n ${PUBLIC_MODULE_VERSION:-} ]]; then
  (
    cd "$test_root/consumer"
    GOWORK=off go get "$public_module@$PUBLIC_MODULE_VERSION"
  )
else
  GOWORK=off go -C "$repo_root/tsc" run "$repo_root/scripts/check-public-module-zip.go" \
    "$public_source" "$public_module" "$source_version" "$test_root/public"
  (
    cd "$test_root/public"
    GOWORK=off go build ./...
  )
  (
    cd "$test_root/consumer"
    GOWORK=off go mod edit -require="$public_module@$source_version"
    GOWORK=off go mod edit -replace="$public_module=../public"
  )
fi

(
  cd "$test_root/consumer"
  GOWORK=off go mod tidy
  EXPECTED_TYPESCRIPT_VERSION=${source_version#v} GOWORK=off go run .
)
