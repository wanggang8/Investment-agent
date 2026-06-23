#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p87-portfolio-state-allocation"
CONFIG_PATH="$TMP_DIR/config.p87.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p87.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P87_SERVER_PORT:-18187}"
WEB_PORT="${P87_WEB_PORT:-14287}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P87_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety}"
SUMMARY_PATH="$ARTIFACT_DIR/portfolio-state-allocation-summary.json"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
HANDLER_LOG="$ARTIFACT_DIR/go-handler-tests.log"
RULE_LOG="$ARTIFACT_DIR/go-rule-tests.log"
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
INSERT OR IGNORE INTO rule_versions (rule_version,status,rules_json,effective_at,created_at)
VALUES ('v_p87','active','{}','$NOW','$NOW');

INSERT OR REPLACE INTO decision_records (decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,capability_status,capability_reason,source_verification_status,risk_reason_code,media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at)
VALUES
('decision_p87_sell_only','req_p87_sell_only','consultation','159915','P87 买入逻辑破坏后如何处理？','completed','formal_trade_advice','high_risk','in_scope','P87 sell-only state safety scenario','satisfied','buy_logic_broken','{"heat":"neutral"}','["calm"]','["BUY_LOGIC","manual_confirmation_required"]','[]','sell_only','买入逻辑破坏，持仓进入只卖不买；禁止新增买入和加仓，仅允许用户线下卖出后回填。','["新增买入","加仓","自动交易","外部推送"]','["用户线下卖出后记录确认","继续观察"]','pending',NULL,NULL,'v_p87','[{"agent_name":"P87StateAnalyst","conclusion":"买入逻辑确认破坏，进入只卖不买。","key_reasons":["至少 2 个 A/S 级独立信源确认不利变化","禁止新增买入和加仓"],"risk_warnings":["系统只允许人工线下处理后回填，不自动交易"],"confidence":"high","evidence_ids":[]}]','{"precision_status":"unavailable","reason":"P87 state safety scenario","sample_count":0,"scenarios":[],"disclaimer":"不构成收益承诺。"}','[{"priority":1,"rule_id":"BUY_LOGIC","result":"sell_only"}]','{"p87":"sell_only"}','$NOW'),
('decision_p87_frozen_watch','req_p87_frozen_watch','consultation','511880','P87 多源验证不足时是否能交易？','completed','non_trade_record','frozen_watch','in_scope','P87 frozen-watch state safety scenario','failed','source_insufficient','{"heat":"neutral"}','["calm"]','["SOURCE_VERIFICATION"]','["formal_source_count_below_2"]','frozen_watch','多源验证不足，进入冻结观察；等待更多 A/S 级独立信源，不生成交易确认。','["新增买入","卖出","减仓","自动交易","外部推送"]','["等待更多 A/S 级独立信源","补充正式证据"]','not_required',NULL,NULL,'v_p87','[{"agent_name":"P87SourceAnalyst","conclusion":"正式证据不足，冻结观察。","key_reasons":["少于 2 个 A/S 级独立信源"],"risk_warnings":["当前不生成交易类建议"],"confidence":"low","evidence_ids":[]}]','{"precision_status":"unavailable","reason":"formal evidence insufficient","sample_count":0,"scenarios":[],"disclaimer":"正式证据不足时不声明交易确认或收益预期可靠。"}','[{"priority":1,"rule_id":"SOURCE","result":"frozen_watch"}]','{"p87":"frozen_watch"}','$NOW'),
('decision_p87_insufficient','req_p87_insufficient','consultation','000000','P87 数据不足时是否会给交易建议？','completed','non_trade_record','insufficient_data','unknown','P87 insufficient-data safety scenario','failed','missing_prerequisites','{"heat":"unknown"}','[]','["DATA_GATE"]','["missing_account_or_market_data"]','insufficient_data','账户或行情等必需数据不足；前端进入信息不足状态，不生成交易类建议。','["交易确认","自动交易","外部推送","收益预期可靠通过"]','["补充账户、行情、估值和正式证据"]','not_required',NULL,NULL,'v_p87','[{"agent_name":"P87DataGate","conclusion":"数据不足，无法生成交易类建议。","key_reasons":["账户、行情或正式证据缺失"],"risk_warnings":["不声明安全边际、估值高低或交易确认已通过"],"confidence":"low","evidence_ids":[]}]','{"precision_status":"unavailable","reason":"missing prerequisites","sample_count":0,"scenarios":[],"disclaimer":"必需依赖缺失时不声明安全边际、交易确认或收益预期可靠。"}','[{"priority":1,"rule_id":"DATA_GATE","result":"insufficient_data"}]','{"p87":"insufficient_data"}','$NOW');
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

P87_CAPTURE_SCREENSHOTS=1 \
P87_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p87-portfolio-state-allocation-safety.spec.ts

run_logged "$HANDLER_LOG" go test ./internal/application/handler -run 'Portfolio|Confirmation|DecisionDetail' -count=1
run_logged "$RULE_LOG" go test ./internal/domain/rule -run 'P75PositionStateAndCooldownBoundaries|P75PortfolioAllocationAndTakeProfitReadback|P75RulePriorityAndRootRules' -count=1

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

python3 "$ROOT_DIR/scripts/p87_portfolio_state_allocation_sqlite_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR/browser-results.json" "$ARTIFACT_DIR_DISPLAY" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR_DISPLAY" "$DB_CHECK_LOG" "$HANDLER_LOG" "$RULE_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, handler_log, rule_log, browser_results = sys.argv[1:7]

def read_kv(path):
    out = {}
    for line in Path(path).read_text(encoding="utf-8").splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key] = value
    return out

def go_status(path):
    text = Path(path).read_text(encoding="utf-8")
    return "passed" if "FAIL" not in text and re.search(r"(?m)^ok\s+", text) else "failed"

browser = json.loads(Path(browser_results).read_text(encoding="utf-8"))
db = read_kv(db_log)
handler_status = go_status(handler_log)
rule_status = go_status(rule_log)
payload = {
    "status": "passed" if browser.get("status") == "passed" and db.get("status") == "passed" and handler_status == "passed" and rule_status == "passed" else "failed",
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "go_tests": {
        "handler": {"status": handler_status, "log": f"{artifact_dir}/go-handler-tests.log"},
        "rule": {"status": rule_status, "log": f"{artifact_dir}/go-rule-tests.log"},
    },
    "safety": {
        "frozen_or_insufficient_confirmations": int(db.get("frozen_or_insufficient_confirmations", "-1")),
        "forbidden_broker_order_push_tables": int(db.get("forbidden_broker_order_push_tables", "-1")),
        "auto_confirmation_rows": int(db.get("auto_confirmation_rows", "-1")),
        "claim_boundary": "P87 verifies local UI/API/SQLite portfolio state and allocation facts plus safe-degradation readback; it does not claim broker connectivity, automatic trading, automatic confirmation, complete monthly attribution, or full original-requirement closure.",
    },
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
