#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p63-full-ui-regression"
CONFIG_PATH="$TMP_DIR/config.p63.yaml"
LOCAL_CONFIG="${P63_LOCAL_CONFIG:-$ROOT_DIR/configs/config.local.yaml}"
SQLITE_PATH="$TMP_DIR/investment-agent-p63.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SERVER_PORT="${P63_SERVER_PORT:-18081}"
WEB_PORT="${P63_WEB_PORT:-14176}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P63_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-18-p63}"
SERVER_PID=""
WEB_PID=""

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
}
trap cleanup EXIT

if [[ ! -f "$LOCAL_CONFIG" ]]; then
  echo "P63 requires a local config with a temporary LLM key: $LOCAL_CONFIG" >&2
  exit 1
fi

rm -rf "$TMP_DIR"
mkdir -p "$VECLITE_PATH" "$ARTIFACT_DIR"

python3 - "$LOCAL_CONFIG" "$CONFIG_PATH" "$SQLITE_PATH" "$VECLITE_PATH" "$SERVER_PORT" <<'PY'
import re
import sys
from pathlib import Path

source, target, sqlite_path, veclite_path, server_port = sys.argv[1:6]
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
      "000001": "gssz0000001"
    szse_base_url: "https://www.szse.cn"
    csrc_base_url: "https://www.csrc.gov.cn"
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

export INVESTMENT_AGENT_CONFIG="$CONFIG_PATH"
go run ./cmd/smoke-seed
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

P63_CAPTURE_SCREENSHOTS=1 \
P63_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts p63-full-ui-regression.spec.ts

cleanup
SERVER_PID=""
WEB_PID=""
