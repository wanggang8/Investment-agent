#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p84-portfolio-confirmation"
CONFIG_PATH="$TMP_DIR/config.p84.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p84.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P84_SERVER_PORT:-18184}"
WEB_PORT="${P84_WEB_PORT:-14284}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P84_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation}"
SUMMARY_PATH="$ARTIFACT_DIR/portfolio-confirmation-summary.json"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
HANDLER_LOG="$ARTIFACT_DIR/go-handler-tests.log"
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
export GOCACHE="${GOCACHE:-$ROOT_DIR/.cache/go-build}"
mkdir -p "$GOCACHE"

go run ./cmd/smoke-seed >"$SEED_LOG" 2>&1

NOW="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
sqlite3 "$SQLITE_PATH" <<SQL
INSERT OR REPLACE INTO decision_records (decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,capability_status,capability_reason,source_verification_status,risk_reason_code,media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at)
VALUES ('decision_p84_pending','req_p84_pending','consultation','510300','P84 pending manual confirmation','completed','formal_trade_advice','normal','in_scope','P84 portfolio confirmation scenario','satisfied','','{"heat":"neutral"}','["calm"]','["portfolio_rebalance_review","manual_confirmation_required"]','[]','hold','P84 本地组合确认场景：等待用户线下处理结果','["自动交易","外部推送"]','["记录线下处理","继续观察"]','pending',NULL,NULL,'v_p84','[{"agent_name":"P84LocalAnalyst","conclusion":"仅用于本地确认链路验收","key_reasons":["持仓比例需人工复核","确认只能记录线下动作"],"risk_warnings":["不自动交易"],"confidence":"medium","evidence_ids":[]}]','{"precision_status":"insufficient","reason":"P84 portfolio data-impact scenario","sample_count":0,"scenarios":[]}','[{"step":"rule","result":"hold"}]','{"p84":"portfolio_confirmation"}','$NOW');
SQL

go build -o "$SERVER_BIN" ./cmd/server
"$SERVER_BIN" >"$SERVER_LOG" 2>&1 &
SERVER_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null

VITE_API_PROXY_TARGET="http://127.0.0.1:$SERVER_PORT" bash -c 'cd "$1" && exec env VITE_API_PROXY_TARGET="$2" ./node_modules/.bin/vite --host 127.0.0.1 --port "$3" --strictPort' _ "$ROOT_DIR/web" "http://127.0.0.1:$SERVER_PORT" "$WEB_PORT" >"$WEB_LOG" 2>&1 &
WEB_PID="$!"

for _ in {1..90}; do
  if curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null

P84_CAPTURE_SCREENSHOTS=1 \
P84_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p84-portfolio-confirmation-data-impact.spec.ts

run_logged "$HANDLER_LOG" go test ./internal/application/handler -run 'Portfolio|Confirmation|DecisionLoop|Review|Audit' -count=1

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

python3 "$ROOT_DIR/scripts/p84_portfolio_confirmation_sqlite_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR_DISPLAY" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR_DISPLAY" "$DB_CHECK_LOG" "$HANDLER_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, handler_log, browser_results = sys.argv[1:6]

def read_kv(path):
    out = {}
    for line in Path(path).read_text(encoding="utf-8").splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key] = value
    return out

handler_text = Path(handler_log).read_text(encoding="utf-8")
go_status = "passed" if "FAIL" not in handler_text and re.search(r"(?m)^ok\s+", handler_text) else "failed"
browser = json.loads(Path(browser_results).read_text(encoding="utf-8"))
db = read_kv(db_log)
payload = {
    "status": "passed" if browser.get("status") == "passed" and db.get("status") == "passed" and go_status == "passed" else "failed",
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "go_tests": {"handler": {"status": go_status, "log": f"{artifact_dir}/go-handler-tests.log"}},
    "safety": {
        "forbidden_broker_order_push_tables": int(db.get("forbidden_broker_order_push_tables", "-1")),
        "auto_confirmation_rows": int(db.get("auto_confirmation_rows", "-1")),
        "claim_boundary": "P84 verifies local portfolio/manual confirmation data impact; it does not claim broker sync, automatic trading, or complete allocation/rebalance policy coverage.",
    },
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
