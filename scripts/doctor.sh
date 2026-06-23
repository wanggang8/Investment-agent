#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

require_docker
ensure_env_file
create_runtime_dirs

if [[ -z "$(env_value DEEPSEEK_API_KEY)" ]]; then
  echo "warning=DEEPSEEK_API_KEY is empty; LLM-backed analysis will degrade safely until configured"
else
  echo "deepseek_api_key=configured"
fi

echo "docker=ok"
echo "compose=ok"
echo "env_file=$ENV_FILE"
echo "data_dir=$(data_dir)"
echo "sqlite_dir=$(data_dir)/data/sqlite"
echo "veclite_dir=$(data_dir)/data/veclite"
