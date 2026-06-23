#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p85-expected-return-analysis"
CONFIG_PATH="$TMP_DIR/config.p85.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p85.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P85_SERVER_PORT:-18185}"
WEB_PORT="${P85_WEB_PORT:-14285}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P85_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p85-expected-return-analysis}"
SUMMARY_PATH="$ARTIFACT_DIR/expected-return-summary.json"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
WORKFLOW_TEST_LOG="$ARTIFACT_DIR/go-workflow-tests.log"
HANDLER_TEST_LOG="$ARTIFACT_DIR/go-handler-tests.log"
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
DELETE FROM position_snapshots;
DELETE FROM portfolio_snapshots;
DELETE FROM positions WHERE symbol IN ('510300','159915','512000');
DELETE FROM market_snapshots WHERE symbol IN ('510300','159915','512000');
DELETE FROM rag_chunks WHERE summary_id IN (SELECT summary_id FROM intelligence_summary WHERE symbol IN ('510300','159915','512000'));
DELETE FROM intelligence_summary WHERE symbol IN ('510300','159915','512000');
DELETE FROM source_verifications WHERE symbol IN ('510300','159915','512000');
DELETE FROM capability_configs;

INSERT OR IGNORE INTO rule_versions (rule_version,status,rules_json,effective_at,created_at)
VALUES ('v_p85','active','{}','$NOW','$NOW');

INSERT OR REPLACE INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at)
VALUES ('cap_p85','["510300","159915","512000"]','[]','["etf","fund"]','["hold_review","sell_review"]','$NOW');

INSERT OR REPLACE INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at)
VALUES ('snap_p85','$NOW',10000,57000,0.1754385965,0.15,2,'manual','$NOW');
INSERT OR REPLACE INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,updated_at)
VALUES
('pos_p85_510300','510300','沪深300ETF',10000,2.00,3.00,30000,0.50,'normal','2025-01-15','P85 核心持仓预期收益验收','core','$NOW'),
('pos_p85_159915','159915','创业板ETF',10000,2.00,1.70,17000,-0.15,'normal','2025-03-10','P85 下行情景验收','satellite','$NOW');
INSERT OR REPLACE INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,created_at)
VALUES
('ps_p85_510300','snap_p85','pos_p85_510300','510300','沪深300ETF',10000,2.00,3.00,30000,0.50,'normal','2025-01-15','P85 核心持仓预期收益验收','core','$NOW'),
('ps_p85_159915','snap_p85','pos_p85_159915','159915','创业板ETF',10000,2.00,1.70,17000,-0.15,'normal','2025-03-10','P85 下行情景验收','satellite','$NOW');

INSERT OR REPLACE INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at)
VALUES
('market_p85_510300','510300','2026-06-22',3.00,1.10,20,18,45,30,'normal','neutral','{"metadata":{"nav_history":[2.20,2.25,2.30,2.35,2.40,2.45,2.50,2.55,2.60,2.65,2.70,2.75,2.80,2.85,2.90,2.95,2.97,3.00],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"},"source_verification":{"freshness":"fresh"}}}}','$NOW'),
('market_p85_159915','159915','2026-06-22',1.70,1.50,60,58,55,42,'normal','neutral','{"metadata":{"nav_history":[2.10,2.08,2.05,2.03,2.00,1.98,1.95,1.92,1.90,1.88,1.85,1.82,1.80,1.78,1.75,1.73,1.71,1.70],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"},"source_verification":{"freshness":"fresh"}}}}','$NOW'),
('market_p85_512000','512000','2026-06-22',1.20,0.70,0,0,0,0,'normal','neutral','{"metadata":{"p34_source_health":{"market_history":{"freshness":"missing"},"valuation_percentiles":{"freshness":"missing"}}}}','$NOW');

INSERT OR REPLACE INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at)
VALUES
('intel_p85_510300','P85OfficialA','A','https://example.invalid/p85/510300','2026-06-22T00:00:00Z','$NOW','hash_p85_510300','P85 510300 正式证据','local-p85','$NOW'),
('intel_p85_159915','P85OfficialA','A','https://example.invalid/p85/159915','2026-06-22T00:00:00Z','$NOW','hash_p85_159915','P85 159915 正式证据','local-p85','$NOW'),
('intel_p85_512000','P85OfficialA','A','https://example.invalid/p85/512000','2026-06-22T00:00:00Z','$NOW','hash_p85_512000','P85 512000 正式证据','local-p85','$NOW');
INSERT OR REPLACE INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at)
VALUES
('sum_p85_510300','intel_p85_510300','510300','510300','normal','neutral','P85 510300 正式证据摘要：用于真实工作流预期收益验收。','A','formal',1,1,'vg_p85_510300','$NOW'),
('sum_p85_159915','intel_p85_159915','159915','159915','normal','neutral','P85 159915 正式证据摘要：用于真实工作流下行情景验收。','A','formal',1,1,'vg_p85_159915','$NOW'),
('sum_p85_512000','intel_p85_512000','512000','512000','normal','neutral','P85 512000 正式证据摘要：用于样本不足安全降级验收。','A','formal',1,1,'vg_p85_512000','$NOW');
INSERT OR REPLACE INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at)
VALUES
('ver_p85_510300','vg_p85_510300','event_p85_510300','510300','normal','formal','satisfied',2,1,'A','2026-06-22T00:00:00Z','["sum_p85_510300"]','$NOW'),
('ver_p85_159915','vg_p85_159915','event_p85_159915','159915','normal','formal','satisfied',2,1,'A','2026-06-22T00:00:00Z','["sum_p85_159915"]','$NOW'),
('ver_p85_512000','vg_p85_512000','event_p85_512000','512000','normal','formal','satisfied',2,1,'A','2026-06-22T00:00:00Z','["sum_p85_512000"]','$NOW');
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

P85_CAPTURE_SCREENSHOTS=1 \
P85_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p85-expected-return-analysis-accuracy.spec.ts

run_logged "$WORKFLOW_TEST_LOG" go test ./internal/application/workflow -run 'ExpectedReturn|expectedReturn' -count=1
run_logged "$HANDLER_TEST_LOG" go test ./internal/application/handler -run 'ExpectedReturn|ConsultDecisionAcceptsExpectedReturnDynamicInputs|DecisionDetail' -count=1

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

python3 "$ROOT_DIR/scripts/p85_expected_return_sqlite_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR/browser-results.json" "$ARTIFACT_DIR_DISPLAY" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR_DISPLAY" "$DB_CHECK_LOG" "$WORKFLOW_TEST_LOG" "$HANDLER_TEST_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import os
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, workflow_log, handler_log, browser_results = sys.argv[1:7]

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
workflow_status = go_status(workflow_log)
handler_status = go_status(handler_log)
payload = {
    "status": "passed" if browser.get("status") == "passed" and db.get("status") == "passed" and workflow_status == "passed" and handler_status == "passed" else "failed",
    "artifact_dir": artifact_dir,
    "llm_mode": "deepseek" if bool(os.environ.get("DEEPSEEK_API_KEY")) else "static_fallback_no_real_llm_claim",
    "browser": browser,
    "db_readback": db,
    "go_tests": {
        "workflow": {"status": workflow_status, "log": f"{artifact_dir}/go-workflow-tests.log"},
        "handler": {"status": handler_status, "log": f"{artifact_dir}/go-handler-tests.log"},
    },
    "safety": {
        "operation_confirmations_p85": int(db.get("operation_confirmations_p85", "-1")),
        "forbidden_broker_order_push_tables": int(db.get("forbidden_broker_order_push_tables", "-1")),
        "auto_confirmation_rows": int(db.get("auto_confirmation_rows", "-1")),
        "claim_boundary": "P85 verifies deterministic expected-return UI/API/SQLite behavior and safety boundaries; it does not claim future return accuracy, broker connectivity, automatic trading, or longitudinal market prediction accuracy.",
    },
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
