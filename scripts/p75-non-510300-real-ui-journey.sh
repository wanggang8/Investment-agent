#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p75-non-510300-real-ui-journey"
CONFIG_PATH="$TMP_DIR/config.p75.yaml"
LOCAL_CONFIG="${P75_LOCAL_CONFIG:-$ROOT_DIR/configs/config.local.yaml}"
SQLITE_PATH="$TMP_DIR/investment-agent-p75.db"
VECLITE_PATH="$TMP_DIR/veclite.json"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SOURCE_LOG="$TMP_DIR/source.log"
SOURCE_PORT_FILE="$TMP_DIR/source.port"
REQUEST_LOG="$TMP_DIR/source-requests.json"
SERVER_PORT="${P75_SERVER_PORT:-18095}"
WEB_PORT="${P75_WEB_PORT:-14195}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P75_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-20-p75-non-510300}"
PRECHECK_LOG="$ARTIFACT_DIR/precheck.log"
DB_CHECK_LOG="$ARTIFACT_DIR/db-impact-check.log"
SERVER_PID=""
WEB_PID=""
SOURCE_PID=""

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
  stop_process "$WEB_PID"
  stop_process "$SERVER_PID"
  stop_process "$SOURCE_PID"
}
trap cleanup EXIT

require_log_contains() {
  local file="$1"
  local needle="$2"
  if ! grep -Fq "$needle" "$file"; then
    echo "P75 required output missing: $needle" >&2
    echo "Log: $file" >&2
    exit 1
  fi
}

assert_real_acceptance_config() {
  local file="$1"
  if grep -Eq '^[[:space:]]*use_stub:[[:space:]]*true[[:space:]]*$' "$file"; then
    echo "P75 rejects mock-only acceptance config: data_sources.use_stub=true" >&2
    exit 1
  fi
  if ! grep -Fq 'accepted_local' "$file"; then
    echo "P75 requires accepted-local HTTP source names for request-construction evidence" >&2
    exit 1
  fi
  if ! grep -Fq 'market_endpoint:' "$file"; then
    echo "P75 requires configured market_endpoint" >&2
    exit 1
  fi
  if ! grep -Fq 'public_evidence:' "$file"; then
    echo "P75 requires public evidence collectors" >&2
    exit 1
  fi
}

if [[ ! -f "$LOCAL_CONFIG" ]]; then
  echo "P75 requires a local config with a real LLM key: $LOCAL_CONFIG" >&2
  exit 1
fi

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR" "$ARTIFACT_DIR"

node "$ROOT_DIR/scripts/p75_accepted_local_source_server.mjs" "$SOURCE_PORT_FILE" "$REQUEST_LOG" >"$SOURCE_LOG" 2>&1 &
SOURCE_PID="$!"

for _ in {1..90}; do
  if [[ -s "$SOURCE_PORT_FILE" ]]; then
    break
  fi
  sleep 0.25
done
if [[ ! -s "$SOURCE_PORT_FILE" ]]; then
  echo "P75 accepted-local source server did not start" >&2
  cat "$SOURCE_LOG" >&2 || true
  exit 1
fi
SOURCE_PORT="$(tr -d '[:space:]' < "$SOURCE_PORT_FILE")"
SOURCE_BASE_URL="http://127.0.0.1:$SOURCE_PORT"

python3 - "$LOCAL_CONFIG" "$CONFIG_PATH" "$SQLITE_PATH" "$VECLITE_PATH" "$SERVER_PORT" "$SOURCE_BASE_URL" <<'PY'
import re
import sys
from pathlib import Path

source, target, sqlite_path, veclite_path, server_port, source_base_url = sys.argv[1:7]
text = Path(source).read_text(encoding="utf-8")

def value(name, default=""):
    match = re.search(rf"(?m)^\s*{re.escape(name)}:\s*(.*)$", text)
    if not match:
        return default
    return match.group(1).strip().strip('"').strip("'")

api_key = value("api_key")
base_url = value("base_url", "https://api.deepseek.com")
model = value("model", "deepseek-chat")
timeout = value("timeout_seconds", "60")
if not api_key:
    raise SystemExit("deepseek.api_key is empty in local config")

Path(target).write_text(f'''server:
  host: "127.0.0.1"
  port: {server_port}

sqlite:
  path: "{sqlite_path}"

veclite:
  path: "{veclite_path}"

deepseek:
  api_key: "{api_key}"
  base_url: "{base_url}"
  model: "{model}"
  timeout_seconds: {timeout}

data_sources:
  enabled:
    - "accepted_local"
  use_stub: false
  market_endpoint: "{source_base_url}/market"
  intelligence_endpoint: ""
  public_evidence:
    enabled: true
    sources:
      - "cninfo"
      - "szse"
    cninfo_base_url: "{source_base_url}"
    cninfo_org_ids:
      "159915": "accepted-local-159915"
      "510300": "9900000091"
    szse_base_url: "{source_base_url}"
    csrc_base_url: "{source_base_url}"
  market_collectors:
    enabled: false
    sources:
      - "csindex"
      - "eastmoney_fund"
    csindex_base_url: "https://www.csindex.com.cn"
    eastmoney_fund_base_url: "https://fund.eastmoney.com"

daily_auto_run:
  enabled: false
  run_time: "08:30"
  timezone: "Asia/Shanghai"
  scope: "holdings"
  retry: 1
  timeout_seconds: 900
  max_symbols: 20

log:
  level: "error"
''', encoding="utf-8")
PY

assert_real_acceptance_config "$CONFIG_PATH"

export INVESTMENT_AGENT_CONFIG="$CONFIG_PATH"

{
  echo "P75 non-510300 precheck started at $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  go run ./cmd/smoke-seed
  go run ./cmd/agent --task market-refresh --symbol 159915
  go run ./cmd/agent --task public-evidence-refresh --symbol 159915 --start-date 2026-06-01 --end-date 2026-06-30
  go run ./cmd/agent --task evidence-index
} 2>&1 | tee "$PRECHECK_LOG"

require_log_contains "$PRECHECK_LOG" "task market-refresh completed"
require_log_contains "$PRECHECK_LOG" "task public-evidence-refresh completed"
require_log_contains "$PRECHECK_LOG" "task evidence-index completed"

go build -o "$SERVER_BIN" ./cmd/server

"$SERVER_BIN" >"$SERVER_LOG" 2>&1 &
SERVER_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null

VITE_API_PROXY_TARGET="http://127.0.0.1:$SERVER_PORT" bash -c 'cd "$1" && exec env VITE_API_PROXY_TARGET="$2" ./node_modules/.bin/vite --host 127.0.0.1 --port "$3" --strictPort' _ "$ROOT_DIR/web" "http://127.0.0.1:$SERVER_PORT" "$WEB_PORT" >"$WEB_LOG" 2>&1 &
WEB_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null

P75_CAPTURE_SCREENSHOTS=1 \
P75_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p75-non-510300-real-ui-journey.spec.ts

python3 "$ROOT_DIR/scripts/p75_non_510300_sqlite_check.py" "$SQLITE_PATH" "$REQUEST_LOG" "$ARTIFACT_DIR" 2>&1 | tee "$DB_CHECK_LOG"

cleanup
SERVER_PID=""
WEB_PID=""
SOURCE_PID=""
