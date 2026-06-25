#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="${P116_TMP_DIR:-$(mktemp -d "$ROOT_DIR/tmp/p116-multi-fund-transaction-ledger-acceptance.XXXXXX")}"
CONFIG_PATH="$TMP_DIR/config.p116.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p116.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
FRONTEND_LOG="$TMP_DIR/frontend.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P116_SERVER_PORT:-18116}"
FRONTEND_PORT="${P116_FRONTEND_PORT:-19116}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P116_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance}"
BROWSER_SUMMARY="$ARTIFACT_DIR/browser/p116-browser-results.json"
FINAL_SUMMARY="$ARTIFACT_DIR/p116-scenario-summary.json"
SERVER_PID=""
FRONTEND_PID=""

stop_process() {
  local pid="$1"
  if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
  fi
}

cleanup() {
  stop_process "$FRONTEND_PID"
  stop_process "$SERVER_PID"
}
trap cleanup EXIT

mkdir -p "$TMP_DIR" "$VECLITE_PATH" "$ARTIFACT_DIR/api_sqlite" "$ARTIFACT_DIR/browser" "$ARTIFACT_DIR/degradation"

cat >"$CONFIG_PATH" <<YAML
runtime:
  mode: "development"

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
  timeout_seconds: 15

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

go run ./cmd/smoke-seed >"$SEED_LOG" 2>&1
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

python3 "$ROOT_DIR/scripts/p116_multi_fund_transaction_ledger_acceptance.py" \
  --base-url "http://127.0.0.1:$SERVER_PORT" \
  --sqlite "$SQLITE_PATH" \
  --artifact-dir "$ARTIFACT_DIR"

(
  cd "$ROOT_DIR/web"
  VITE_API_PROXY_TARGET="http://127.0.0.1:$SERVER_PORT" npm run dev -- --host 127.0.0.1 --port "$FRONTEND_PORT" >"$FRONTEND_LOG" 2>&1
) &
FRONTEND_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$FRONTEND_PORT" >/dev/null 2>&1; then
    break
  fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$FRONTEND_PORT" >/dev/null

(
  cd "$ROOT_DIR/web"
  P116_ARTIFACT_DIR="$ARTIFACT_DIR/browser" E2E_BASE_URL="http://127.0.0.1:$FRONTEND_PORT" npx playwright test e2e/p116-multi-fund-transaction-ledger-acceptance.spec.ts --project=chromium
)

python3 "$ROOT_DIR/scripts/p116_multi_fund_transaction_ledger_acceptance.py" \
  --base-url "http://127.0.0.1:$SERVER_PORT" \
  --sqlite "$SQLITE_PATH" \
  --artifact-dir "$ARTIFACT_DIR" \
  --browser-summary "$BROWSER_SUMMARY" \
  --merge-only

cat "$FINAL_SUMMARY"
