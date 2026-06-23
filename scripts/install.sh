#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

ensure_env_file
create_runtime_dirs
MODE="$(detect_install_mode)"

if [[ "$MODE" == "upgrade" ]]; then
  echo "existing_install=detected"
  bash "$ROOT_DIR/scripts/upgrade.sh"
  exit 0
fi

echo "first_install=starting"
bash "$ROOT_DIR/scripts/doctor.sh"
compose_cmd -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d --build
wait_for_health
write_release_state "first_install"
print_next_steps
