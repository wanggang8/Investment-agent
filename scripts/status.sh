#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

ensure_env_file
compose_cmd -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" ps
if wait_for_health; then
  echo "status=healthy"
else
  echo "status=unhealthy"
  exit 1
fi
