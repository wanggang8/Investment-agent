#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p75-sop-failure-real-ui"
CONFIG_PATH="$TMP_DIR/config.p75.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p75-sop.db"
VECLITE_PATH="$TMP_DIR/veclite.json"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P75_SERVER_PORT:-18096}"
WEB_PORT="${P75_WEB_PORT:-14196}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P75_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-20-p75-sop-failure}"
DB_CHECK_LOG="$ARTIFACT_DIR/db-impact-check.log"
FINAL_RULE_APPLY="${P75_FINAL_RULE_APPLY:-0}"
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

require_equal() {
  local actual="$1"
  local expected="$2"
  local message="$3"
  if [[ "$actual" != "$expected" ]]; then
    echo "P75 check failed: $message actual=$actual expected=$expected" >&2
    exit 1
  fi
}

require_at_least() {
  local actual="$1"
  local minimum="$2"
  local message="$3"
  if (( actual < minimum )); then
    echo "P75 check failed: $message actual=$actual minimum=$minimum" >&2
    exit 1
  fi
}

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR" "$ARTIFACT_DIR"

cat >"$CONFIG_PATH" <<YAML
server:
  host: "127.0.0.1"
  port: $SERVER_PORT

sqlite:
  path: "$SQLITE_PATH"

veclite:
  path: "$VECLITE_PATH"

deepseek:
  api_key: "p75-sop-failure-ui-no-llm-call"
  base_url: "https://api.deepseek.com"
  model: "deepseek-chat"
  timeout_seconds: 15

data_sources:
  enabled:
    - "accepted_local"
  use_stub: false
  market_endpoint: "http://127.0.0.1:$SERVER_PORT/p75-unused-market"
  intelligence_endpoint: ""
  public_evidence:
    enabled: false
    sources: []
  market_collectors:
    enabled: false
    sources: []

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

sqlite3 "$SQLITE_PATH" <<'SQL'
DELETE FROM capability_configs WHERE capability_id <> 'cap_p75_sop';
DELETE FROM market_snapshots WHERE market_snapshot_id IN ('market_smoke_p39');

INSERT OR REPLACE INTO capability_configs (capability_id,asset_types_json,symbols_json,excluded_symbols_json,strategy_scope_json,updated_at)
VALUES ('cap_p75_sop','["ETF","fund"]','["510300","159915","512000","588000","513100","515000"]','["999999"]','["long_term_etf","discipline_review"]','2026-06-20T08:00:00Z');

INSERT OR REPLACE INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at)
VALUES (
  'market_p75_source_health',
  '510300',
  '2026-06-20',
  4.20,
  1.20,
  85,
  82,
  70,
  65,
  'warning',
  'neutral',
  '{"source_name":"csindex_p75","source_level":"A","source_type":"public_market","metadata":{"p34_source_health":{"index_valuation_files":{"source_name":"csindex_p75","source_level":"A","source_type":"index_valuation","freshness":"stale","data_date":"2026-06-10","last_success_at":"2026-06-10T08:00:00Z","last_failure_at":"2026-06-20T08:00:00Z","failure_category":"stale","affected_symbols":["510300"]},"capital_flow":{"source_name":"eastmoney_p75","source_level":"B","source_type":"fund_flow","freshness":"parse_error","data_date":"2026-06-19","last_success_at":"2026-06-19T08:00:00Z","last_failure_at":"2026-06-20T08:00:00Z","failure_category":"parse_error","affected_symbols":["159915"]}},"p34_data_categories":["index_valuation_files","capital_flow"]}}',
  '2026-06-20T08:00:00Z'
);

INSERT OR REPLACE INTO decision_records (decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,capability_status,capability_reason,source_verification_status,risk_reason_code,media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at)
VALUES
('decision_p75_insufficient','req_p75_insufficient','consultation','999999','P75 unsupported symbol insufficient data check','degraded','rejection_record','insufficient_data','out_of_scope','标的不在能力圈或数据依赖缺失','failed','insufficient_evidence','{}','[]','["source_verification_failed"]','["formal evidence missing"]','insufficient_data','证据不足，暂停交易类建议','["新增买入","交易确认"]','["补齐正式证据后再复核"]','not_required',NULL,'market_p75_source_health','v1','[{"agent_name":"P75SafetyAnalyst","conclusion":"证据不足，暂停交易类建议","key_reasons":["标的不在能力圈","正式证据不足"],"risk_warnings":["不允许交易确认"],"confidence":"low","evidence_ids":[]}]','{"precision_status":"unavailable","reason":"insufficient formal evidence","sample_count":0,"sample_window":"none","screening_condition":"none","scenarios":[]}','[{"step":"source_verification","result":"failed"}]','{"market_snapshot":{"symbol":"999999","trade_date":"2026-06-20","close_price":0,"pe_percentile":0,"pb_percentile":0}}','2026-06-20T08:01:00Z'),
('decision_p75_model_unavailable','req_p75_model_unavailable','consultation','510300','P75 model unavailable safe degradation check','degraded','formal_trade_advice','insufficient_data','in_scope','能力圈内但分析服务不可用','satisfied','data_degraded','{}','[]','["llm_unavailable"]','["ANALYST_UNAVAILABLE"]','insufficient_data','LLM 降级，暂停交易类建议','["新增买入","自动确认"]','["仅查看已有规则与数据"]','not_required',NULL,'market_p75_source_health','v1','[{"agent_name":"P75LLMDegraded","conclusion":"分析服务暂不可用","key_reasons":["模型不可用"],"risk_warnings":["LLM 降级，暂停交易类建议"],"confidence":"low","evidence_ids":[]}]','{"precision_status":"unavailable","reason":"model unavailable","sample_count":0,"sample_window":"none","screening_condition":"none","scenarios":[]}','[{"step":"analyst","result":"degraded"}]','{"market_snapshot":{"symbol":"510300","trade_date":"2026-06-20","close_price":4.2,"pe_percentile":85,"pb_percentile":82}}','2026-06-20T08:02:00Z');
INSERT OR REPLACE INTO decision_records (decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,capability_status,capability_reason,source_verification_status,risk_reason_code,media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at)
VALUES
('decision_p75_mark_error','req_p75_mark_error','consultation','510300','P75 mark error critical mutation check','completed','formal_trade_advice','normal','in_scope','能力圈内','satisfied','','{}','[]','["hold_review"]','[]','hold','持有并继续人工复核','["自动交易"]','["记录计划","标记错误"]','pending',NULL,'market_p75_source_health','v1','[{"agent_name":"P75ReviewAnalyst","conclusion":"持有并继续人工复核","key_reasons":["用于错误标记验收"],"risk_warnings":["系统只记录复盘样本"],"confidence":"medium","evidence_ids":[]}]','{"precision_status":"insufficient","reason":"review sample","sample_count":4,"sample_window":"2026-Q2","screening_condition":"P75 UI mark error","scenarios":[{"scenario":"base","return_range":"0%~2%","probability":null,"trigger":"review"}]}','[{"step":"rule","result":"hold"}]','{"market_snapshot":{"symbol":"510300","trade_date":"2026-06-20","close_price":4.2,"pe_percentile":85,"pb_percentile":82}}','2026-06-20T08:03:00Z');

INSERT OR REPLACE INTO risk_alerts (alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,prohibited_actions_json,suggested_actions_json,related_decision_id,related_report_id,related_notification_id,related_audit_event_id,last_triggered_at,resolved_at,resolution_reason,created_at,updated_at)
VALUES
('risk_p75_sop_a','buy_thesis_broken','critical','active','510300','P75 SOP-A 持仓下跌超过5%，先检查买入逻辑破坏和正式证据','{"sop":"SOP-A","data_prerequisites":["position_snapshot","buy_thesis","formal_evidence"],"llm_role":"explain_only"}','["新增买入","摊低成本"]','["复核买入逻辑","记录继续观察"]','decision_p75_insufficient',NULL,NULL,'audit_p75_sop_a_trigger','2026-06-20T08:10:00Z',NULL,NULL,'2026-06-20T08:10:00Z','2026-06-20T08:10:00Z'),
('risk_p75_sop_b','valuation_high','warning','active','159915','P75 SOP-B 持仓上涨超过20%，复核止盈分段和回归核心资产','{"sop":"SOP-B","data_prerequisites":["profit_ratio","valuation_percentile","allocation"],"llm_role":"explain_only"}','["追涨加仓"]','["复核分段止盈","止盈资金优先回归核心资产"]',NULL,NULL,NULL,'audit_p75_sop_b_trigger','2026-06-20T08:11:00Z',NULL,NULL,'2026-06-20T08:11:00Z','2026-06-20T08:11:00Z'),
('risk_p75_sop_c','sentiment_extreme','warning','active','512000','P75 SOP-C 热点追涨冲动，情绪极端时先冷静复核','{"sop":"SOP-C","data_prerequisites":["sentiment_percentile","media_heat","user_emotion_tag"],"llm_role":"cooldown_explanation"}','["追热点新增买入"]','["至少冷静1个交易日","检查估值和证据"]',NULL,NULL,NULL,'audit_p75_sop_c_trigger','2026-06-20T08:12:00Z',NULL,NULL,'2026-06-20T08:12:00Z','2026-06-20T08:12:00Z'),
('risk_p75_sop_d','sentiment_extreme','critical','active','588000','P75 SOP-D 恐慌清仓语言，先执行冷静期和客观数据复核','{"sop":"SOP-D","data_prerequisites":["panic_language","sentiment_percentile","historical_analog"],"llm_role":"calm_summary"}','["立即清仓","情绪化卖出"]','["冷静期后再复核","查看历史类比样本"]',NULL,NULL,NULL,'audit_p75_sop_d_trigger','2026-06-20T08:13:00Z',NULL,NULL,'2026-06-20T08:13:00Z','2026-06-20T08:13:00Z'),
('risk_p75_sop_e','insufficient_evidence','critical','active','513100','P75 SOP-E 宏观灰犀牛，需要两类正式证据和波动分支复核','{"sop":"SOP-E","data_prerequisites":["formal_evidence","volatility","expected_return"],"llm_role":"scenario_reassessment"}','["无证据减仓","新增风险暴露"]','["补齐两类正式证据","复核情景概率来源"]',NULL,NULL,NULL,'audit_p75_sop_e_trigger','2026-06-20T08:14:00Z',NULL,NULL,'2026-06-20T08:14:00Z','2026-06-20T08:14:00Z'),
('risk_p75_sop_f','data_degraded','critical','active','515000','P75 SOP-F 黑天鹅事件，24小时冻结主动动作并等待A级信源复核','{"sop":"SOP-F","data_prerequisites":["event_time","two_a_level_sources","freeze_expiry"],"llm_role":"freeze_boundary"}','["冻结期交易确认","自动触发交易"]','["冻结观察24小时","等待两类A级信源影响评估"]',NULL,NULL,NULL,'audit_p75_sop_f_trigger','2026-06-20T08:15:00Z',NULL,NULL,'2026-06-20T08:15:00Z','2026-06-20T08:15:00Z');

INSERT OR REPLACE INTO audit_events (audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at)
VALUES
('audit_p75_sop_a_trigger','req_p75_sop_a','decision_p75_insufficient','risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-A 持仓下跌超过5%，先检查买入逻辑破坏和正式证据','risk_alert','risk_p75_sop_a','2026-06-20T08:10:00Z'),
('audit_p75_sop_b_trigger','req_p75_sop_b',NULL,'risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-B 持仓上涨超过20%，复核止盈分段和回归核心资产','risk_alert','risk_p75_sop_b','2026-06-20T08:11:00Z'),
('audit_p75_sop_c_trigger','req_p75_sop_c',NULL,'risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-C 热点追涨冲动，情绪极端时先冷静复核','risk_alert','risk_p75_sop_c','2026-06-20T08:12:00Z'),
('audit_p75_sop_d_trigger','req_p75_sop_d',NULL,'risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-D 恐慌清仓语言，先执行冷静期和客观数据复核','risk_alert','risk_p75_sop_d','2026-06-20T08:13:00Z'),
('audit_p75_sop_e_trigger','req_p75_sop_e',NULL,'risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-E 宏观灰犀牛，需要两类正式证据和波动分支复核','risk_alert','risk_p75_sop_e','2026-06-20T08:14:00Z'),
('audit_p75_sop_f_trigger','req_p75_sop_f',NULL,'risk_alert_sop','P75SOPSeed','system','risk_alert','trigger_risk_alert',NULL,NULL,NULL,'success',NULL,NULL,'active','v1',NULL,'sop_summary','P75 SOP-F 黑天鹅事件，24小时冻结主动动作并等待A级信源复核','risk_alert','risk_p75_sop_f','2026-06-20T08:15:00Z');

INSERT OR REPLACE INTO rule_proposals (proposal_id,proposal_type,status,source_error_case_id,title,proposal_version,before_rule_json,after_rule_json,reason,impact_scope_json,risk_notes_json,sample_count,final_confirmed_at,final_confirmed_note,applied_rule_version,related_error_cases_json,created_at)
VALUES
('prop_p75_ui_send_gatekeeper','threshold','pending_user_confirm',NULL,'P75 UI 送审提案','draft','{"threshold":80}','{"threshold":82}','P75 真实 UI 送入守门人审计','{"scope":"p75_ui_gatekeeper_mutation"}','["不自动应用正式规则"]',5,NULL,NULL,NULL,'[]','2026-06-20T08:19:00Z'),
('prop_p75_gatekeeper_denied','risk_rule','pending_final_confirm',NULL,'P75 守门人否决样例','draft','{"rule":"old"}','{"rule":"would_violate_root_rule"}','根本规则冲突，不允许应用','{"scope":"p75_ui_failure_state"}','["违反根本规则"]',5,NULL,NULL,NULL,'[]','2026-06-20T08:20:00Z'),
('prop_p75_gatekeeper_review','sop','pending_final_confirm',NULL,'P75 守门人用户复核样例','draft','{"rule":"old"}','{"rule":"needs_user_context"}','样本不足且涉及用户行为模式，需要用户复核','{"scope":"p75_ui_failure_state"}','["样本不足","情绪偏差"]',2,NULL,NULL,NULL,'[]','2026-06-20T08:21:00Z');

INSERT OR REPLACE INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,required_changes,violates_fundamental_rule,has_rule_conflict,backtest_metrics_json,allow_apply,audited_rule_version,created_at)
VALUES
('gate_p75_denied','prop_p75_gatekeeper_denied','rejected','违反根本规则：不允许把证据不足转成交易确认','保留人工复核边界',1,1,'{"sample_count":5,"degradation":true}',0,'draft','2026-06-20T08:22:00Z'),
('gate_p75_review','prop_p75_gatekeeper_review','needs_user_review','样本不足且用户情绪偏差，需要人工复核','补充样本和用户确认',0,0,'{"sample_count":2,"emotion_bias":true}',0,'draft','2026-06-20T08:23:00Z');

INSERT OR REPLACE INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,source_explanation_json,metrics_json,risk_notes_json,related_error_cases_json,related_decision_ids_json,related_risk_alert_ids_json,related_audit_event_ids_json,safety_note,created_at,updated_at)
VALUES
('validation_p75_ui_send_gatekeeper','prop_p75_ui_send_gatekeeper','draft','passed',5,'2026-Q2','passed','low','passed','passed','{"source_case_count":5}','{"hit_count":5,"missing_evidence_count":0}','["P75 UI 送审提案已通过本地回放门禁"]','[]','["decision_p75_mark_error"]','["risk_p75_sop_b"]','["audit_p75_sop_b_trigger"]','规则效果验证只用于本地规则治理；规则生效仍需用户手动最终确认。','2026-06-20T08:24:30Z','2026-06-20T08:24:30Z'),
('validation_p75_review','prop_p75_gatekeeper_review','draft','insufficient',2,'2026-Q2','insufficient','high','mixed','needs_user_review','{"source_case_count":2}','{"hit_count":1,"missing_evidence_count":1}','["样本不足"]','[]','["decision_p75_insufficient"]','["risk_p75_sop_d"]','["audit_p75_sop_d_trigger"]','规则效果验证只用于本地规则治理；规则生效仍需用户手动确认。','2026-06-20T08:24:00Z','2026-06-20T08:24:00Z');
SQL

BASE_CONFIRMATIONS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM operation_confirmations;")"
BASE_POSITION_TX="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM position_transactions;")"
BASE_RULE_VERSIONS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_versions;")"
BASE_ERROR_CASES="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM error_cases;")"
BASE_GATEKEEPER_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM gatekeeper_audits;")"

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
P75_FINAL_RULE_APPLY="$FINAL_RULE_APPLY" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p75-sop-failure-real-ui.spec.ts

AFTER_CONFIRMATIONS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM operation_confirmations;")"
AFTER_POSITION_TX="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM position_transactions;")"
AFTER_RULE_VERSIONS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_versions;")"
AFTER_ERROR_CASES="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM error_cases;")"
AFTER_GATEKEEPER_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM gatekeeper_audits;")"
LIFECYCLE_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM audit_events WHERE workflow_type='risk_alert_sop' AND node_action='update_risk_alert_lifecycle' AND output_ref LIKE 'risk_p75_sop_%';")"
SOP_TERMINAL_COUNT="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM risk_alerts WHERE alert_id IN ('risk_p75_sop_a','risk_p75_sop_b','risk_p75_sop_c','risk_p75_sop_d','risk_p75_sop_e','risk_p75_sop_f') AND sop_status IN ('observing','escalated','resolved');")"
FORBIDDEN_TABLES="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%');")"
MARK_ERROR_CASES="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM error_cases WHERE decision_id='decision_p75_mark_error' AND root_cause_tag='rule_threshold_issue';")"
MARK_ERROR_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM audit_events WHERE decision_id='decision_p75_mark_error' AND action='mark_error' AND after_state='marked_error';")"
GATEKEEPER_NODE_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM audit_events WHERE proposal_id='prop_p75_ui_send_gatekeeper' AND action='audit_rule_change' AND node_name IN ('ProposalLoadNode','FundamentalRuleCheckNode','ConflictCheckNode','BacktestNode','AuditDecisionNode','AuditRecordNode');")"
GATEKEEPER_STATUS="$(sqlite3 "$SQLITE_PATH" "SELECT status FROM rule_proposals WHERE proposal_id='prop_p75_ui_send_gatekeeper';")"
FINAL_RULE_AUDITS="$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM audit_events WHERE proposal_id='prop_p75_ui_send_gatekeeper' AND action='update_rule' AND after_state='applied' AND output_ref_type='rule_version';")"
FINAL_RULE_VERSION="$(sqlite3 "$SQLITE_PATH" "SELECT COALESCE(applied_rule_version,'') FROM rule_proposals WHERE proposal_id='prop_p75_ui_send_gatekeeper';")"

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

require_equal "$AFTER_CONFIRMATIONS" "$((BASE_CONFIRMATIONS + 1))" "mark-error UI must create exactly one local confirmation"
require_equal "$AFTER_POSITION_TX" "$BASE_POSITION_TX" "SOP/failure UI must not create position transactions"
if [[ "$FINAL_RULE_APPLY" == "1" ]]; then
  require_equal "$AFTER_RULE_VERSIONS" "$((BASE_RULE_VERSIONS + 1))" "P82 final-confirm UI must create exactly one new local rule version"
else
  require_equal "$AFTER_RULE_VERSIONS" "$BASE_RULE_VERSIONS" "gatekeeper deny/user-review UI must not create rule versions"
fi
require_equal "$AFTER_ERROR_CASES" "$((BASE_ERROR_CASES + 1))" "mark-error UI must create exactly one error case"
require_equal "$AFTER_GATEKEEPER_AUDITS" "$((BASE_GATEKEEPER_AUDITS + 1))" "gatekeeper UI send must create exactly one gatekeeper audit"
require_at_least "$LIFECYCLE_AUDITS" 6 "every SOP UI lifecycle action must write audit_events"
require_equal "$SOP_TERMINAL_COUNT" "6" "every SOP risk alert must be read back in updated status"
require_equal "$FORBIDDEN_TABLES" "0" "forbidden broker/order/push tables must not exist"
require_equal "$MARK_ERROR_CASES" "1" "mark-error UI must persist root-cause error case"
require_equal "$MARK_ERROR_AUDITS" "1" "mark-error UI must write audit event"
require_equal "$GATEKEEPER_NODE_AUDITS" "6" "gatekeeper UI send must write all node audit events"
if [[ "$FINAL_RULE_APPLY" == "1" ]]; then
  require_equal "$GATEKEEPER_STATUS" "applied" "P82 final-confirm UI must apply proposal only after explicit user final confirmation"
  require_equal "$FINAL_RULE_AUDITS" "1" "P82 final-confirm UI must write exactly one update_rule audit event"
  if [[ -z "$FINAL_RULE_VERSION" ]]; then
    echo "P75 check failed: P82 final-confirm UI must persist applied_rule_version" >&2
    exit 1
  fi
else
  require_equal "$GATEKEEPER_STATUS" "pending_final_confirm" "gatekeeper UI send must stop before final rule application"
  require_equal "$FINAL_RULE_AUDITS" "0" "default gatekeeper UI journey must not apply a rule"
fi

{
  echo "status=passed"
  echo "final_rule_apply=$FINAL_RULE_APPLY"
  echo "base_confirmations=$BASE_CONFIRMATIONS"
  echo "after_confirmations=$AFTER_CONFIRMATIONS"
  echo "base_position_transactions=$BASE_POSITION_TX"
  echo "after_position_transactions=$AFTER_POSITION_TX"
  echo "base_rule_versions=$BASE_RULE_VERSIONS"
  echo "after_rule_versions=$AFTER_RULE_VERSIONS"
  echo "base_error_cases=$BASE_ERROR_CASES"
  echo "after_error_cases=$AFTER_ERROR_CASES"
  echo "base_gatekeeper_audits=$BASE_GATEKEEPER_AUDITS"
  echo "after_gatekeeper_audits=$AFTER_GATEKEEPER_AUDITS"
  echo "lifecycle_audits=$LIFECYCLE_AUDITS"
  echo "sop_updated_status_count=$SOP_TERMINAL_COUNT"
  echo "mark_error_cases=$MARK_ERROR_CASES"
  echo "mark_error_audits=$MARK_ERROR_AUDITS"
  echo "gatekeeper_node_audits=$GATEKEEPER_NODE_AUDITS"
  echo "gatekeeper_status=$GATEKEEPER_STATUS"
  echo "final_rule_audits=$FINAL_RULE_AUDITS"
  echo "final_rule_version=$FINAL_RULE_VERSION"
  echo "forbidden_broker_order_push_tables=$FORBIDDEN_TABLES"
  echo "artifact_dir=$ARTIFACT_DIR_DISPLAY"
} | tee "$DB_CHECK_LOG"
