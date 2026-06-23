#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENV_FILE="${ROOT_DIR}/.env"
ENV_EXAMPLE="${ROOT_DIR}/.env.example"
STATE_DIR_DEFAULT="${ROOT_DIR}/.investment-agent"
STATE_FILE_NAME="release-state.json"

compose_cmd() {
  if docker compose version >/dev/null 2>&1; then
    docker compose "$@"
    return
  fi
  if command -v docker-compose >/dev/null 2>&1; then
    docker-compose "$@"
    return
  fi
  echo "Docker Compose is required. Install Docker Desktop or the Docker Compose plugin." >&2
  return 1
}

require_docker() {
  if ! command -v docker >/dev/null 2>&1; then
    echo "Docker is required. Install Docker Desktop or Docker Engine first." >&2
    return 1
  fi
  if ! docker info >/dev/null 2>&1; then
    echo "Docker is installed but not running or not accessible." >&2
    return 1
  fi
  compose_cmd version >/dev/null
}

ensure_env_file() {
  if [[ ! -f "$ENV_FILE" ]]; then
    cp "$ENV_EXAMPLE" "$ENV_FILE"
  fi
  ensure_env_value "INVESTMENT_AGENT_DATA_DIR" "$(absolute_data_dir)"
}

env_value() {
  local key="$1"
  if [[ ! -f "$ENV_FILE" ]]; then
    return 0
  fi
  grep -E "^${key}=" "$ENV_FILE" | tail -n 1 | cut -d= -f2- || true
}

ensure_env_value() {
  local key="$1"
  local value="$2"
  if grep -qE "^${key}=" "$ENV_FILE"; then
    python3 - "$ENV_FILE" "$key" "$value" <<'PY'
import sys
from pathlib import Path

path = Path(sys.argv[1])
key = sys.argv[2]
value = sys.argv[3]
lines = path.read_text(encoding="utf-8").splitlines()
out = []
changed = False
for line in lines:
    if line.startswith(key + "="):
        out.append(f"{key}={value}")
        changed = True
    else:
        out.append(line)
if not changed:
    out.append(f"{key}={value}")
path.write_text("\n".join(out) + "\n", encoding="utf-8")
PY
  else
    printf '%s=%s\n' "$key" "$value" >>"$ENV_FILE"
  fi
}

absolute_data_dir() {
  local raw
  raw="$(env_value INVESTMENT_AGENT_DATA_DIR)"
  if [[ -z "$raw" || "$raw" == ".investment-agent" ]]; then
    printf '%s\n' "$STATE_DIR_DEFAULT"
    return
  fi
  python3 - "$ROOT_DIR" "$raw" <<'PY'
import sys
from pathlib import Path

root = Path(sys.argv[1]).resolve()
raw = Path(sys.argv[2]).expanduser()
print(raw.resolve(strict=False) if raw.is_absolute() else (root / raw).resolve(strict=False))
PY
}

data_dir() {
  absolute_data_dir
}

state_file() {
  printf '%s/%s\n' "$(data_dir)" "$STATE_FILE_NAME"
}

create_runtime_dirs() {
  local dir
  dir="$(data_dir)"
  mkdir -p "$dir/data/sqlite" "$dir/data/veclite" "$dir/backups" "$dir/logs"
}

detect_install_mode() {
  local dir
  dir="$(data_dir)"
  if [[ -f "$(state_file)" || -f "$dir/data/sqlite/investment-agent.db" ]]; then
    echo "upgrade"
  else
    echo "first_install"
  fi
}

write_release_state() {
  local mode="$1"
  local dir
  dir="$(data_dir)"
  mkdir -p "$dir"
  python3 - "$dir/$STATE_FILE_NAME" "$mode" "$ROOT_DIR" <<'PY'
import json
import subprocess
import sys
from datetime import datetime, timezone
from pathlib import Path

path = Path(sys.argv[1])
mode = sys.argv[2]
root = sys.argv[3]
try:
    commit = subprocess.check_output(["git", "-C", root, "rev-parse", "HEAD"], text=True).strip()
except Exception:
    commit = "unknown"
data = {
    "status": "installed",
    "last_mode": mode,
    "updated_at": datetime.now(timezone.utc).isoformat(),
    "source_commit": commit,
    "safety": "install/upgrade preserves SQLite, VecLite, backups, logs, and .env",
}
path.write_text(json.dumps(data, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY
}

wait_for_health() {
  local web_port
  web_port="$(env_value INVESTMENT_AGENT_WEB_PORT)"
  web_port="${web_port:-4173}"
  for _ in {1..120}; do
    if curl -fsS "http://127.0.0.1:${web_port}/api/v1/health" >/dev/null 2>&1; then
      echo "health=passed url=http://127.0.0.1:${web_port}"
      return 0
    fi
    sleep 1
  done
  echo "health=failed url=http://127.0.0.1:${web_port}" >&2
  return 1
}

print_next_steps() {
  local web_port
  web_port="$(env_value INVESTMENT_AGENT_WEB_PORT)"
  web_port="${web_port:-4173}"
  echo "Investment Agent is running."
  echo "Open: http://127.0.0.1:${web_port}"
  echo "Data directory: $(data_dir)"
}
