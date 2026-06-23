#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/recovery-smoke"
SOURCE_CONFIG="$TMP_DIR/source-config.yaml"
RESTORE_CONFIG="$TMP_DIR/restore-config.yaml"
SOURCE_DB="$TMP_DIR/source/investment-agent.db"
RESTORE_DB="$TMP_DIR/restore/investment-agent.db"
SOURCE_VECLITE="$TMP_DIR/source/veclite"
RESTORE_VECLITE="$TMP_DIR/restore/veclite"
BACKUP_DIR="$TMP_DIR/backups"
SERVER_BIN="$TMP_DIR/server"
SERVER_LOG="$TMP_DIR/server.log"
SERVER_PORT="${RECOVERY_SMOKE_SERVER_PORT:-18180}"
SERVER_PID=""

stop_process() {
  local pid="$1"
  if [[ -z "$pid" ]]; then
    return
  fi
  if kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
  fi
}

cleanup() {
  stop_process "$SERVER_PID"
}
trap cleanup EXIT

write_config() {
  local path="$1"
  local sqlite_path="$2"
  local veclite_path="$3"
  cat > "$path" <<YAML
server:
  host: "127.0.0.1"
  port: $SERVER_PORT

sqlite:
  path: "$sqlite_path"

veclite:
  path: "$veclite_path"

deepseek:
  api_key: ""
  base_url: "https://api.deepseek.com"

data_sources:
  enabled:
    - "stub"
  use_stub: true
  market_endpoint: ""
  intelligence_endpoint: ""
  public_evidence:
    enabled: false
    sources:
      - "cninfo"
    cninfo_base_url: "https://www.cninfo.com.cn"
    cninfo_org_ids:
      "510300": "9900000091"
    szse_base_url: "https://www.szse.cn"
    csrc_base_url: "https://www.csrc.gov.cn"
  market_collectors:
    enabled: false
    sources:
      - "csindex"
    csindex_base_url: "https://www.csindex.com.cn"
    eastmoney_fund_base_url: "https://fund.eastmoney.com"

daily_auto_run:
  enabled: false
  run_time: "08:30"
  timezone: "UTC"
  scope: "holdings"
  retry: 1
  timeout_seconds: 900
  max_symbols: 20

log:
  level: "error"
YAML
}

rm -rf "$TMP_DIR"
mkdir -p "$SOURCE_VECLITE" "$RESTORE_VECLITE" "$BACKUP_DIR"
write_config "$SOURCE_CONFIG" "$SOURCE_DB" "$SOURCE_VECLITE"
write_config "$RESTORE_CONFIG" "$RESTORE_DB" "$RESTORE_VECLITE"

INVESTMENT_AGENT_CONFIG="$SOURCE_CONFIG" go run ./cmd/smoke-seed
BACKUP_OUTPUT="$(go run ./cmd/agent --config "$SOURCE_CONFIG" --backup "$BACKUP_DIR")"
BACKUP_FILE="${BACKUP_OUTPUT#backup created:}"

go run ./cmd/agent --config "$RESTORE_CONFIG" --recovery-smoke "$BACKUP_FILE" >/dev/null
go build -o "$SERVER_BIN" ./cmd/server

INVESTMENT_AGENT_CONFIG="$RESTORE_CONFIG" "$SERVER_BIN" >"$SERVER_LOG" 2>&1 &
SERVER_PID="$!"

for _ in {1..60}; do
  if curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/decisions/decision_smoke_p30" | grep -q "P30 本地 E2E smoke 决策"

cleanup
SERVER_PID=""
rm -rf "$TMP_DIR"
