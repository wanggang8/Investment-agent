#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p89-real-provider-dynamic-probability"
CONFIG_PATH="$TMP_DIR/config.p89.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p89.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P89_SERVER_PORT:-18189}"
WEB_PORT="${P89_WEB_PORT:-14289}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P89_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability}"
SUMMARY_PATH="$ARTIFACT_DIR/p89-acceptance-summary.json"
SOURCE_JSON="$ARTIFACT_DIR/p89-source-preverification.json"
SOURCE_PREVERIFY_LOG="$ARTIFACT_DIR/source-preverification.log"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
GO_WORKFLOW_LOG="$ARTIFACT_DIR/go-workflow-tests.log"
GO_HANDLER_LOG="$ARTIFACT_DIR/go-handler-tests.log"
WEB_TEST_LOG="$ARTIFACT_DIR/web-component-tests.log"
FRONTEND_BUILD_LOG="$ARTIFACT_DIR/frontend-build.log"
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

run_logged "$SOURCE_PREVERIFY_LOG" python3 "$ROOT_DIR/scripts/p89_source_preverification.py"
go run ./cmd/smoke-seed >"$SEED_LOG" 2>&1

NOW="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
sqlite3 "$SQLITE_PATH" <<SQL
DELETE FROM position_snapshots;
DELETE FROM portfolio_snapshots;
DELETE FROM positions WHERE symbol IN ('510300','159915','600000');
DELETE FROM market_snapshots WHERE symbol IN ('510300','159915','600000');
DELETE FROM rag_chunks WHERE summary_id IN (SELECT summary_id FROM intelligence_summary WHERE symbol IN ('510300','159915','600000'));
DELETE FROM intelligence_summary WHERE symbol IN ('510300','159915','600000');
DELETE FROM intelligence_items WHERE intelligence_id LIKE 'intel_p89_%';
DELETE FROM source_verifications WHERE symbol IN ('510300','159915','600000');
DELETE FROM capability_configs;

INSERT OR IGNORE INTO rule_versions (rule_version,status,rules_json,effective_at,created_at)
VALUES ('v_p89','active','{}','$NOW','$NOW');

INSERT OR REPLACE INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at)
VALUES ('cap_p89','["510300","159915","600000"]','[]','["etf","fund","stock"]','["hold_review","sell_review"]','$NOW');

INSERT OR REPLACE INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at)
VALUES ('snap_p89','$NOW',100,1300,0.0769230769,0,3,'manual','$NOW');
INSERT OR REPLACE INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,updated_at)
VALUES
('pos_p89_510300','510300','沪深300ETF',100,3,3.3,330,0.1,'normal','2026-01-05','P89 baseline','core','$NOW'),
('pos_p89_159915','159915','创业板ETF',100,2,2,200,0,'normal','2026-01-06','P89 dynamic downshift','satellite','$NOW'),
('pos_p89_600000','600000','浦发银行',100,10,10.5,1050,0.05,'normal','2026-01-07','P89 extreme fear','equity','$NOW');
INSERT OR REPLACE INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,created_at)
VALUES
('ps_p89_510300','snap_p89','pos_p89_510300','510300','沪深300ETF',100,3,3.3,330,0.1,'normal','2026-01-05','P89 baseline','core','$NOW'),
('ps_p89_159915','snap_p89','pos_p89_159915','159915','创业板ETF',100,2,2,200,0,'normal','2026-01-06','P89 dynamic downshift','satellite','$NOW'),
('ps_p89_600000','snap_p89','pos_p89_600000','600000','浦发银行',100,10,10.5,1050,0.05,'normal','2026-01-07','P89 extreme fear','equity','$NOW');

INSERT OR REPLACE INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at)
VALUES
('market_p89_510300','510300','2026-06-22',3.3,1.10,20,18,45,30,'normal','neutral','{"metadata":{"nav_history":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28],"expected_return_historical_samples":[{"scenario":"upside","count":6,"return_range":"8.00%~15.00%","lower_bound":0.08,"upper_bound":0.15,"return_rate":0.12,"trigger":"估值修复"},{"scenario":"base","count":18,"return_range":"0.00%~8.00%","lower_bound":0,"upper_bound":0.08,"return_rate":0.04,"trigger":"维持当前"},{"scenario":"downside","count":6,"return_range":"-12.00%~0.00%","lower_bound":-0.12,"upper_bound":0,"return_rate":-0.06,"trigger":"估值收缩"}],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"}}}}','$NOW'),
('market_p89_159915','159915','2026-06-22',2,1.10,65,58,45,45,'normal','neutral','{"metadata":{"nav_history":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28],"expected_return_historical_samples":[{"scenario":"upside","count":6,"return_range":"8.00%~15.00%","lower_bound":0.08,"upper_bound":0.15,"return_rate":0.12,"trigger":"估值修复"},{"scenario":"base","count":18,"return_range":"0.00%~8.00%","lower_bound":0,"upper_bound":0.08,"return_rate":0.04,"trigger":"维持当前"},{"scenario":"downside","count":6,"return_range":"-12.00%~0.00%","lower_bound":-0.12,"upper_bound":0,"return_rate":-0.06,"trigger":"估值收缩"}],"expected_return_market_state":"stress","expected_return_fundamental_state":"below_expectation","expected_return_pessimistic_path_months":1,"expected_return_assumption_checks":[{"name":"盈利增速","expected":0.08,"actual":0.01,"months_below":2}],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"}}}}','$NOW'),
('market_p89_600000','600000','2026-06-22',10.5,1.10,40,38,45,60,'normal','extreme','{"metadata":{"nav_history":[1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25],"expected_return_historical_samples":[{"scenario":"upside","count":4,"return_range":"6.00%~12.00%","lower_bound":0.06,"upper_bound":0.12,"return_rate":0.09,"trigger":"恐慌后估值修复"},{"scenario":"base","count":10,"return_range":"-2.00%~6.00%","lower_bound":-0.02,"upper_bound":0.06,"return_rate":0.02,"trigger":"情绪低位震荡"},{"scenario":"downside","count":6,"return_range":"-18.00%~-4.00%","lower_bound":-0.18,"upper_bound":-0.04,"return_rate":-0.10,"trigger":"恐慌继续扩散"}],"expected_return_historical_contexts":[{"label":"极端恐惧样本","window":"2018Q4, 2020Q1, 2022Q4","sample_count":20,"outcome":"暂停主动交易建议","max_drawdown":-0.18,"recovery":"3-9 个月","source":"local_public_history"}],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"}}}}','$NOW');

INSERT OR REPLACE INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at)
VALUES
('intel_p89_510300_a','P89OfficialA','A','https://example.invalid/p89/510300-a','2026-06-22T00:00:00Z','$NOW','hash_p89_510300_a','P89 510300 正式证据 A','local-p89','$NOW'),
('intel_p89_510300_b','P89OfficialB','A','https://example.invalid/p89/510300-b','2026-06-22T00:05:00Z','$NOW','hash_p89_510300_b','P89 510300 正式证据 B','local-p89','$NOW'),
('intel_p89_159915_a','P89OfficialA','A','https://example.invalid/p89/159915-a','2026-06-22T00:00:00Z','$NOW','hash_p89_159915_a','P89 159915 正式证据 A','local-p89','$NOW'),
('intel_p89_159915_b','P89OfficialB','A','https://example.invalid/p89/159915-b','2026-06-22T00:05:00Z','$NOW','hash_p89_159915_b','P89 159915 正式证据 B','local-p89','$NOW'),
('intel_p89_600000_a','P89OfficialA','A','https://example.invalid/p89/600000-a','2026-06-22T00:00:00Z','$NOW','hash_p89_600000_a','P89 600000 正式证据 A','local-p89','$NOW'),
('intel_p89_600000_b','P89OfficialB','A','https://example.invalid/p89/600000-b','2026-06-22T00:05:00Z','$NOW','hash_p89_600000_b','P89 600000 正式证据 B','local-p89','$NOW');
INSERT OR REPLACE INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at)
VALUES
('sum_p89_510300_a','intel_p89_510300_a','510300','510300','normal','neutral','P89 baseline 正式证据 A。','A','formal',1,1,'vg_p89_510300','$NOW'),
('sum_p89_510300_b','intel_p89_510300_b','510300','510300','normal','neutral','P89 baseline 正式证据 B。','A','formal',1,1,'vg_p89_510300','$NOW'),
('sum_p89_159915_a','intel_p89_159915_a','159915','159915','normal','negative','P89 dynamic 正式证据 A。','A','formal',1,1,'vg_p89_159915','$NOW'),
('sum_p89_159915_b','intel_p89_159915_b','159915','159915','normal','negative','P89 dynamic 正式证据 B。','A','formal',1,1,'vg_p89_159915','$NOW'),
('sum_p89_600000_a','intel_p89_600000_a','600000','600000','normal','negative','P89 extreme fear 正式证据 A。','A','formal',1,1,'vg_p89_600000','$NOW'),
('sum_p89_600000_b','intel_p89_600000_b','600000','600000','normal','negative','P89 extreme fear 正式证据 B。','A','formal',1,1,'vg_p89_600000','$NOW');
INSERT OR REPLACE INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at)
VALUES
('ver_p89_510300','vg_p89_510300','event_p89_510300','510300','normal','formal','satisfied',2,2,'A','2026-06-22T00:05:00Z','["sum_p89_510300_a","sum_p89_510300_b"]','$NOW'),
('ver_p89_159915','vg_p89_159915','event_p89_159915','159915','normal','formal','satisfied',2,2,'A','2026-06-22T00:05:00Z','["sum_p89_159915_a","sum_p89_159915_b"]','$NOW'),
('ver_p89_600000','vg_p89_600000','event_p89_600000','600000','normal','formal','satisfied',2,2,'A','2026-06-22T00:05:00Z','["sum_p89_600000_a","sum_p89_600000_b"]','$NOW');
SQL

go build -o "$SERVER_BIN" ./cmd/server
"$SERVER_BIN" >"$SERVER_LOG" 2>&1 &
SERVER_PID="$!"

for _ in {1..120}; do
  if curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$SERVER_PORT/api/v1/health" >/dev/null

VITE_API_PROXY_TARGET="http://127.0.0.1:$SERVER_PORT" bash -c 'cd "$1" && exec env VITE_API_PROXY_TARGET="$2" ./node_modules/.bin/vite --host 127.0.0.1 --port "$3" --strictPort' _ "$ROOT_DIR/web" "http://127.0.0.1:$SERVER_PORT" "$WEB_PORT" >"$WEB_LOG" 2>&1 &
WEB_PID="$!"

for _ in {1..120}; do
  if curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null 2>&1; then break; fi
  sleep 0.5
done
curl -fsS "http://127.0.0.1:$WEB_PORT" >/dev/null

P89_CAPTURE_SCREENSHOTS=1 \
P89_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p89-real-provider-dynamic-probability.spec.ts

run_logged "$GO_WORKFLOW_LOG" go test ./internal/application/workflow ./internal/domain/rule -run 'TestP89|TestP88ExpectedReturnDynamicMonitoring|TestEvaluateSentiment' -count=1
run_logged "$GO_HANDLER_LOG" go test ./internal/application/handler -run 'TestP89|ConsultDecisionAcceptsExpectedReturnDynamicInputs' -count=1
run_logged "$WEB_TEST_LOG" npm --prefix "$ROOT_DIR/web" test -- DecisionTrace.test.tsx SettingsPage.test.tsx
run_logged "$FRONTEND_BUILD_LOG" npm --prefix "$ROOT_DIR/web" run build
python3 "$ROOT_DIR/scripts/p89_sqlite_readback_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR/browser-results.json" "$SOURCE_JSON" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR" "$DB_CHECK_LOG" "$SOURCE_JSON" "$GO_WORKFLOW_LOG" "$GO_HANDLER_LOG" "$WEB_TEST_LOG" "$FRONTEND_BUILD_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, source_json, workflow_log, handler_log, web_log, build_log, browser_results = sys.argv[1:10]

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

def text_status(path):
    text = Path(path).read_text(encoding="utf-8")
    return "passed" if "failed" not in text.lower() and "error" not in text.lower() else "failed"

browser = json.loads(Path(browser_results).read_text(encoding="utf-8"))
source = json.loads(Path(source_json).read_text(encoding="utf-8"))
db = read_kv(db_log)
status = "passed" if browser.get("status") == "passed" and db.get("status") == "passed" and go_status(workflow_log) == "passed" and go_status(handler_log) == "passed" and text_status(web_log) == "passed" and text_status(build_log) == "passed" else "failed"
payload = {
    "status": status,
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "source_preverification": source,
    "go_tests": {
        "workflow_rule": {"status": go_status(workflow_log), "log": f"{artifact_dir}/go-workflow-tests.log"},
        "handler": {"status": go_status(handler_log), "log": f"{artifact_dir}/go-handler-tests.log"},
    },
    "web_tests": {"status": text_status(web_log), "log": f"{artifact_dir}/web-component-tests.log"},
    "frontend_build": {"status": text_status(build_log), "log": f"{artifact_dir}/frontend-build.log"},
    "claim_boundary": "P89 upgrades only rows with direct provider or UI/API/SQLite evidence. Capital-flow remains partial when provider_status=blocked.",
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
