#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p88-remaining-full-release-blockers"
CONFIG_PATH="$TMP_DIR/config.p88.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p88.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P88_SERVER_PORT:-18188}"
WEB_PORT="${P88_WEB_PORT:-14288}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P88_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers}"
SUMMARY_PATH="$ARTIFACT_DIR/p88-acceptance-summary.json"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
GO_WORKFLOW_LOG="$ARTIFACT_DIR/go-workflow-tests.log"
GO_HANDLER_LOG="$ARTIFACT_DIR/go-handler-tests.log"
WEB_TEST_LOG="$ARTIFACT_DIR/web-component-tests.log"
SOURCE_PREVERIFY_LOG="$ARTIFACT_DIR/source-preverification.log"
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
DELETE FROM positions WHERE symbol IN ('510300','159915','600000');
DELETE FROM market_snapshots WHERE symbol IN ('510300','159915','600000');
DELETE FROM rag_chunks WHERE summary_id IN (SELECT summary_id FROM intelligence_summary WHERE symbol IN ('510300','159915','600000'));
DELETE FROM intelligence_summary WHERE symbol IN ('510300','159915','600000');
DELETE FROM source_verifications WHERE symbol IN ('510300','159915','600000');
DELETE FROM capability_configs;
DELETE FROM rule_proposals WHERE proposal_version LIKE 'p88-sop-addendum-%';
DELETE FROM notifications WHERE source_type='rule_proposal';

INSERT OR IGNORE INTO rule_versions (rule_version,status,rules_json,effective_at,created_at)
VALUES ('v_p88','active','{}','$NOW','$NOW');

INSERT OR REPLACE INTO capability_configs (capability_id,symbols_json,excluded_symbols_json,asset_types_json,strategy_scope_json,updated_at)
VALUES ('cap_p88','["510300","159915","600000"]','[]','["etf","fund","stock"]','["hold_review","sell_review"]','$NOW');

INSERT OR REPLACE INTO portfolio_snapshots (snapshot_id,snapshot_time,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at)
VALUES ('snap_p88','$NOW',100,1100,0.0909090909,0,3,'manual','$NOW');
INSERT OR REPLACE INTO positions (position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,updated_at)
VALUES
('pos_p88_510300','510300','沪深300ETF',100,3,6,600,1,'normal','2026-01-05','P88 宽基指数 ETF 验收','core','$NOW'),
('pos_p88_159915','159915','创业板ETF',100,2,3,300,0.5,'normal','2026-01-06','P88 行业成长基金验收','satellite','$NOW'),
('pos_p88_600000','600000','浦发银行',100,1,1,100,0,'normal','2026-01-07','P88 金融成分股路径验收','equity','$NOW');
INSERT OR REPLACE INTO position_snapshots (position_snapshot_id,snapshot_id,position_id,symbol,name,quantity,cost_price,current_price,market_value,unrealized_profit_ratio,position_state,buy_date,buy_reason,asset_tag,created_at)
VALUES
('ps_p88_510300','snap_p88','pos_p88_510300','510300','沪深300ETF',100,3,6,600,1,'normal','2026-01-05','P88 宽基指数 ETF 验收','core','$NOW'),
('ps_p88_159915','snap_p88','pos_p88_159915','159915','创业板ETF',100,2,3,300,0.5,'normal','2026-01-06','P88 行业成长基金验收','satellite','$NOW'),
('ps_p88_600000','snap_p88','pos_p88_600000','600000','浦发银行',100,1,1,100,0,'normal','2026-01-07','P88 金融成分股路径验收','equity','$NOW');

INSERT OR REPLACE INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at)
VALUES
('market_p88_510300','510300','2026-06-22',6,1.10,20,18,45,30,'normal','neutral','{"metadata":{"nav_history":[2.20,2.25,2.30,2.35,2.40,2.45,2.50,2.55,2.60,2.65,2.70,2.75,2.80,2.85,2.90,2.95,2.97,3.00,3.05,3.10,3.15,3.20,3.25,3.30,3.35,3.40,3.45,3.50],"expected_return_historical_samples":[{"scenario":"upside","count":6,"return_range":"12.00%~18.00%","lower_bound":0.12,"upper_bound":0.18,"return_rate":0.15,"trigger":"业绩超预期，估值提升至历史高位"},{"scenario":"base","count":18,"return_range":"4.00%~9.00%","lower_bound":0.04,"upper_bound":0.09,"return_rate":0.065,"trigger":"业绩符合预期，估值维持当前水平"},{"scenario":"downside","count":6,"return_range":"-10.00%~-2.00%","lower_bound":-0.10,"upper_bound":-0.02,"return_rate":-0.06,"trigger":"业绩低于预期，估值收缩"}],"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"},"market_history":{"freshness":"fresh"},"source_verification":{"freshness":"fresh"}}}}','$NOW'),
('market_p88_159915','159915','2026-06-22',3,1.10,55,48,45,30,'normal','neutral','{"metadata":{"p34_source_health":{"valuation_percentiles":{"freshness":"fresh"}}}}','$NOW'),
('market_p88_600000','600000','2026-06-22',1,1.10,40,38,45,30,'normal','neutral','{"metadata":{"p88_structured_fields":{"constituent_financial":{"revenue":4523000000,"net_profit":812000000,"growth":0.137,"disclosure_date":"2026-04-30"}}}}','$NOW');

INSERT OR REPLACE INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at)
VALUES
('intel_p88_510300','P88OfficialA','A','https://example.invalid/p88/510300','2026-06-22T00:00:00Z','$NOW','hash_p88_510300','P88 510300 正式证据','local-p88','$NOW'),
('intel_p88_159915_a','P88OfficialA','A','https://example.invalid/p88/159915-a','2026-06-22T00:00:00Z','$NOW','hash_p88_159915_a','P88 159915 买入逻辑破坏正式证据 A','local-p88','$NOW'),
('intel_p88_159915_s','P88OfficialS','S','https://example.invalid/p88/159915-s','2026-06-22T00:05:00Z','$NOW','hash_p88_159915_s','P88 159915 买入逻辑破坏正式证据 S','local-p88','$NOW'),
('intel_p88_600000','P88OfficialA','A','https://example.invalid/p88/600000','2026-06-22T00:00:00Z','$NOW','hash_p88_600000','P88 600000 重大负面事件单源证据','local-p88','$NOW');
INSERT OR REPLACE INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at)
VALUES
('sum_p88_510300','intel_p88_510300','510300','510300','normal','neutral','P88 510300 正式证据摘要：用于历史预期收益验收。','A','formal',1,1,'vg_p88_510300','$NOW'),
('sum_p88_159915_a','intel_p88_159915_a','159915','159915','buy_logic_break','negative','P88 159915 正式证据摘要 A：A 级独立信源确认买入逻辑破坏。','A','formal',1,1,'vg_p88_159915','$NOW'),
('sum_p88_159915_s','intel_p88_159915_s','159915','159915','buy_logic_break','negative','P88 159915 正式证据摘要 S：S 级独立信源确认买入逻辑破坏。','S','formal',1,1,'vg_p88_159915','$NOW'),
('sum_p88_600000','intel_p88_600000','600000','600000','major_negative','negative','P88 600000 正式证据摘要：单一 A 级信源提示重大负面事件，需冻结观察等待第二信源。','A','formal',1,1,'vg_p88_600000','$NOW');
INSERT OR REPLACE INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at)
VALUES
('ver_p88_510300','vg_p88_510300','event_p88_510300','510300','normal','formal','satisfied',2,1,'A','2026-06-22T00:00:00Z','["sum_p88_510300"]','$NOW'),
('ver_p88_159915','vg_p88_159915','event_p88_159915','159915','buy_logic_break','formal','satisfied',2,2,'S','2026-06-22T00:05:00Z','["sum_p88_159915_a","sum_p88_159915_s"]','$NOW'),
('ver_p88_600000','vg_p88_600000','event_p88_600000','600000','major_negative','formal','satisfied',1,1,'A','2026-06-22T00:00:00Z','["sum_p88_600000"]','$NOW');
SQL

run_logged "$SOURCE_PREVERIFY_LOG" python3 "$ROOT_DIR/scripts/p88_source_preverification.py"

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

P88_CAPTURE_SCREENSHOTS=1 \
P88_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p88-remaining-full-release-blockers.spec.ts

run_logged "$GO_WORKFLOW_LOG" go test ./internal/application/workflow ./internal/domain/rule -run 'TestP88|TestBuildExpectedReturn|TestExpectedReturnNode|TestExpectedReturnSampleCount|TestP88StructuredData' -count=1
run_logged "$GO_HANDLER_LOG" go test ./internal/application/handler ./internal/application/service -run 'P88|ExpectedReturn|Rebalance|RuleProposal|SOPAddendum' -count=1
run_logged "$WEB_TEST_LOG" npm --prefix "$ROOT_DIR/web" test -- --run DecisionTrace.test.tsx PortfolioPage.test.tsx RulesPage.test.tsx

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

python3 "$ROOT_DIR/scripts/p88_sqlite_readback_check.py" "$SQLITE_PATH" "$ARTIFACT_DIR/browser-results.json" "$ARTIFACT_DIR_DISPLAY" >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR_DISPLAY" "$DB_CHECK_LOG" "$GO_WORKFLOW_LOG" "$GO_HANDLER_LOG" "$WEB_TEST_LOG" "$SOURCE_PREVERIFY_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, workflow_log, handler_log, web_log, source_log, browser_results = sys.argv[1:9]

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
db = read_kv(db_log)
payload = {
    "status": "passed" if browser.get("status") == "passed" and db.get("status") == "passed" and go_status(workflow_log) == "passed" and go_status(handler_log) == "passed" and text_status(web_log) == "passed" and text_status(source_log) == "passed" else "failed",
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "go_tests": {
        "workflow_rule": {"status": go_status(workflow_log), "log": f"{artifact_dir}/go-workflow-tests.log"},
        "handler_service": {"status": go_status(handler_log), "log": f"{artifact_dir}/go-handler-tests.log"},
    },
    "web_tests": {"status": text_status(web_log), "log": f"{artifact_dir}/web-component-tests.log"},
    "source_preverification": {"status": text_status(source_log), "log": f"{artifact_dir}/source-preverification.log"},
    "claim_boundary": "P88 upgrades only directly evidenced rows. Structured real-provider rows remain blocked unless a non-mock runtime provider and SQLite readback are proven. Forbidden-capability evidence is limited to the exercised P88 UI/API/SQLite paths and does not replace broader product G9 scans.",
}
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
