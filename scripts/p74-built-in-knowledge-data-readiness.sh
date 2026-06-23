#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p74-built-in-knowledge-data-readiness"
CONFIG_PATH="$TMP_DIR/config.p74.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p74.db"
VECLITE_PATH="$TMP_DIR/veclite.json"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SERVER_PORT="${P74_SERVER_PORT:-18094}"
WEB_PORT="${P74_WEB_PORT:-14194}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P74_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-19-p74}"
API_CHECK_LOG="$ARTIFACT_DIR/api-readiness-check.log"
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

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR" "$ARTIFACT_DIR"

cat > "$CONFIG_PATH" <<YAML
server:
  host: "127.0.0.1"
  port: $SERVER_PORT

sqlite:
  path: "$SQLITE_PATH"

veclite:
  path: "$VECLITE_PATH"

deepseek:
  api_key: ""
  base_url: "https://api.deepseek.com"
  model: "deepseek-chat"
  timeout_seconds: 60

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
  timezone: "Asia/Shanghai"
  scope: "holdings"
  retry: 1
  timeout_seconds: 900
  max_symbols: 20

log:
  level: "error"
YAML

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

python3 "$ROOT_DIR/scripts/p74_readiness_api_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR" "http://127.0.0.1:$SERVER_PORT" 2>&1 | tee "$API_CHECK_LOG"

VITE_API_PROXY_TARGET="http://127.0.0.1:$SERVER_PORT" bash -c 'cd "$1" && exec env VITE_API_PROXY_TARGET="$2" ./node_modules/.bin/vite --host 127.0.0.1 --port "$3" --strictPort' _ "$ROOT_DIR/web" "http://127.0.0.1:$SERVER_PORT" "$WEB_PORT" >"$WEB_LOG" 2>&1 &
WEB_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null

P74_CAPTURE_SCREENSHOTS=1 \
P74_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p74-built-in-knowledge-data-readiness.spec.ts

cleanup
SERVER_PID=""
WEB_PID=""
