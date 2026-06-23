#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

ensure_env_file
create_runtime_dirs
bash "$ROOT_DIR/scripts/doctor.sh"
bash "$ROOT_DIR/scripts/backup.sh" --reason upgrade
compose_cmd -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" up -d --build
wait_for_health
write_release_state "upgrade"
print_next_steps
