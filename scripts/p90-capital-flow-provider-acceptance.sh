#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p90-capital-flow-provider"
CONFIG_PATH="$TMP_DIR/config.p90.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p90.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P90_SERVER_PORT:-18190}"
WEB_PORT="${P90_WEB_PORT:-14290}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P90_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider}"
SUMMARY_PATH="$ARTIFACT_DIR/p90-acceptance-summary.json"
SOURCE_JSON="$ARTIFACT_DIR/p90-source-preverification.json"
SOURCE_PREVERIFY_LOG="$ARTIFACT_DIR/source-preverification.log"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
GO_WORKFLOW_LOG="$ARTIFACT_DIR/go-workflow-tests.log"
WEB_TEST_LOG="$ARTIFACT_DIR/web-component-tests.log"
FRONTEND_BUILD_LOG="$ARTIFACT_DIR/frontend-build.log"
E2E_LOG="$ARTIFACT_DIR/p90-browser-e2e.log"
SERVER_PID=""
WEB_PID=""

stop_process() {
  local pid="$1"
  if [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null; then
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
  fi
}

cleanup() {
  stop_process "$WEB_PID"
  stop_process "$SERVER_PID"
}
trap cleanup EXIT

run_logged() {
  local log_file="$1"
  shift
  "$@" >"$log_file" 2>&1
}

rm -rf "$TMP_DIR"
mkdir -p "$VECLITE_PATH" "$ARTIFACT_DIR"

cat >"$CONFIG_PATH" <<YAML
server:
  host: "127.0.0.1"
  port: $SERVER_PORT

sqlite:
  path: "$SQLITE_PATH"

veclite:
  path: "$VECLITE_PATH"

deepseek:
  api_key: "${DEEPSEEK_API_KEY:-}"
  base_url: "${DEEPSEEK_BASE_URL:-https://api.deepseek.com}"
  model: "${DEEPSEEK_MODEL:-deepseek-chat}"
  timeout_seconds: 30

data_sources:
  enabled: []
  use_stub: false
  market_endpoint: ""
  intelligence_endpoint: ""
  public_evidence:
    enabled: false
    sources:
      - "cninfo"
    cninfo_base_url: "https://www.cninfo.com.cn"
    szse_base_url: "https://www.szse.cn"
    csrc_base_url: "https://www.csrc.gov.cn"
  market_collectors:
    enabled: true
    sources:
      - "p89_structured_public"

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
export GOCACHE="${GOCACHE:-$ROOT_DIR/.cache/go-build}"
mkdir -p "$GOCACHE"

run_logged "$SOURCE_PREVERIFY_LOG" python3 "$ROOT_DIR/scripts/p90_source_preverification.py" --check
go run ./cmd/smoke-seed >"$SEED_LOG" 2>&1

NOW="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
sqlite3 "$SQLITE_PATH" <<SQL
DELETE FROM position_snapshots;
DELETE FROM portfolio_snapshots;
DELETE FROM positions WHERE symbol = '600000';
DELETE FROM market_snapshots WHERE symbol = '600000';
DELETE FROM capability_configs;

INSERT OR IGNORE INTO rule_versions (rule_version,status,rules_json,effective_at,created_at)
VALUES ('v_p90','active','{}','$NOW','$NOW');

INSERT OR REPLACE INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at)
VALUES ('cap_p90','["600000"]','[]','["stock"]','["hold_review","sell_review"]','$NOW');

INSERT OR REPLACE INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at)
VALUES ('snap_p90','$NOW',10000,10000,1,0,0,'manual','$NOW');

INSERT OR REPLACE INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at)
VALUES ('market_p90_seed','600000','2026-06-20',9.00,1.00,50,50,50,50,'normal','neutral','{"metadata":{"p90_seed":"no capital_flow field before real UI refresh","p34_source_health":{"market_history":{"freshness":"fresh","data_date":"2026-06-20","source_name":"p90_seed_without_capital_flow"}}}}','$NOW');
SQL

go build -o "$SERVER_BIN" ./cmd/server
"$SERVER_BIN" >"$SERVER_LOG" 2>&1 &
SERVER_PID="$!"

for _ in {1..120}; do
  if curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null

bash -c 'cd "$1" && exec env VITE_API_PROXY_TARGET="$2" ./node_modules/.bin/vite --host 127.0.0.1 --port "$3" --strictPort' _ "$ROOT_DIR/web" "http://127.0.0.1:$SERVER_PORT" "$WEB_PORT" >"$WEB_LOG" 2>&1 &
WEB_PID="$!"

for _ in {1..120}; do
  if curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null

P90_CAPTURE_SCREENSHOTS=1 \
P90_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p90-capital-flow-provider.spec.ts >"$E2E_LOG" 2>&1

run_logged "$GO_WORKFLOW_LOG" go test ./internal/application/workflow -run 'TestP90Structured|TestP89StructuredPublicCollector|TestP88StructuredData|TestP88MarketRefreshPersistsStructuredDataReadback' -count=1
run_logged "$WEB_TEST_LOG" npm --prefix "$ROOT_DIR/web" test -- SettingsPage.test.tsx
run_logged "$FRONTEND_BUILD_LOG" npm --prefix "$ROOT_DIR/web" run build
python3 "$ROOT_DIR/scripts/p90_sqlite_readback_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR/browser-results.json" "$SOURCE_JSON" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR" "$DB_CHECK_LOG" "$SOURCE_JSON" "$GO_WORKFLOW_LOG" "$WEB_TEST_LOG" "$FRONTEND_BUILD_LOG" "$E2E_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, source_json, workflow_log, web_log, build_log, e2e_log, browser_results = sys.argv[1:10]

def read_kv(path: str) -> dict[str, str]:
    out = {}
    for line in Path(path).read_text(encoding="utf-8").splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key] = value
    return out

browser = json.loads(Path(browser_results).read_text(encoding="utf-8"))
source = json.loads(Path(source_json).read_text(encoding="utf-8"))
db = read_kv(db_log)
payload = {
    "status": "passed",
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "source_preverification": source,
    "go_tests": {"workflow": {"status": "passed", "log": f"{artifact_dir}/go-workflow-tests.log"}},
    "web_tests": {"settings_page": {"status": "passed", "log": f"{artifact_dir}/web-component-tests.log"}},
    "frontend_build": {"status": "passed", "log": f"{artifact_dir}/frontend-build.log"},
    "browser_e2e": {"status": "passed", "log": f"{artifact_dir}/p90-browser-e2e.log"},
    "claim_boundary": "P90 upgrades only REQ-04-016 and REQ-05-003 when Eastmoney H5 public capital-flow source, Settings UI refresh, market snapshot API readback, and SQLite readback all pass without fixture/stub/accepted-local/manual-seed capital-flow values.",
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
