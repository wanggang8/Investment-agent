#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
# shellcheck source=scripts/deploy-lib.sh
source "$ROOT_DIR/scripts/deploy-lib.sh"

REASON="manual"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --reason)
      if [[ $# -lt 2 ]]; then
        echo "--reason requires a value" >&2
        exit 1
      fi
      REASON="$2"
      shift 2
      ;;
    -h|--help)
      echo "Usage: bash scripts/backup.sh [--reason VALUE]"
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

ensure_env_file
create_runtime_dirs
STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BACKUP_DIR="$(data_dir)/backups/${STAMP}-${REASON}"
mkdir -p "$BACKUP_DIR"

if [[ -f "$(data_dir)/data/sqlite/investment-agent.db" ]]; then
  cp "$(data_dir)/data/sqlite/investment-agent.db" "$BACKUP_DIR/investment-agent.db"
fi
if [[ -d "$(data_dir)/data/veclite" ]]; then
  tar -czf "$BACKUP_DIR/veclite.tar.gz" -C "$(data_dir)/data" veclite
fi
cp "$ENV_FILE" "$BACKUP_DIR/env.redacted"
python3 - "$BACKUP_DIR/env.redacted" <<'PY'
import re
import sys
from pathlib import Path

path = Path(sys.argv[1])
text = path.read_text(encoding="utf-8")
text = re.sub(r"(?m)^(DEEPSEEK_API_KEY=).*$", r"\1<redacted>", text)
path.write_text(text, encoding="utf-8")
PY
if [[ -f "$(state_file)" ]]; then
  cp "$(state_file)" "$BACKUP_DIR/release-state.json"
fi

echo "backup=passed path=$BACKUP_DIR"
