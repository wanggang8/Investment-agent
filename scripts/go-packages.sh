#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT"

packages="$(go list ./cmd/... ./internal/... ./pkg/...)"

if printf '%s\n' "$packages" | grep -E '/node_modules(/|$)|^investment-agent/web(/|$)' >/dev/null; then
  echo "go package selection included frontend dependency packages" >&2
  printf '%s\n' "$packages" >&2
  exit 1
fi

printf '%s\n' "$packages"
