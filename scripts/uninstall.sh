#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

PURGE="0"
CONFIRMATION_PHRASE="DELETE INVESTMENT AGENT DATA"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --purge)
      PURGE="1"
      shift
      ;;
    -h|--help)
      echo "Usage: bash scripts/uninstall.sh [--purge]"
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

ensure_env_file
compose_cmd -f "$ROOT_DIR/docker-compose.yml" --env-file "$ENV_FILE" down

if [[ "$PURGE" == "1" ]]; then
  echo "This will delete local SQLite, VecLite, backups, logs, release state, and .env."
  echo "Type: $CONFIRMATION_PHRASE"
  read -r confirmation
  if [[ "$confirmation" != "$CONFIRMATION_PHRASE" ]]; then
    echo "purge=cancelled"
    exit 1
  fi
  rm -rf "$(data_dir)"
  rm -f "$ENV_FILE"
  echo "purge=completed"
else
  echo "uninstall=completed"
  echo "data_preserved=$(data_dir)"
fi
