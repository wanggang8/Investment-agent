#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$ROOT_DIR/tmp/p83-governance-traceability"
CONFIG_PATH="$TMP_DIR/config.p83.yaml"
SQLITE_PATH="$TMP_DIR/investment-agent-p83.db"
VECLITE_PATH="$TMP_DIR/veclite"
SERVER_LOG="$TMP_DIR/server.log"
WEB_LOG="$TMP_DIR/web.log"
SEED_LOG="$TMP_DIR/seed.log"
SERVER_PORT="${P83_SERVER_PORT:-18183}"
WEB_PORT="${P83_WEB_PORT:-14283}"
SERVER_BIN="$TMP_DIR/server"
ARTIFACT_DIR="${P83_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability}"
SUMMARY_PATH="$ARTIFACT_DIR/governance-traceability-summary.json"
HANDLER_LOG="$ARTIFACT_DIR/go-handler-tests.log"
WORKFLOW_LOG="$ARTIFACT_DIR/go-workflow-tests.log"
AGENT_LOG="$ARTIFACT_DIR/go-agent-tests.log"
DB_CHECK_LOG="$ARTIFACT_DIR/db-readback-check.log"
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
VALUES
('decision_p83_monthly_exec','req_p83_monthly_exec','consultation','510300','P83 monthly attribution executed manually','degraded','formal_trade_advice','normal','in_scope','P83 governance traceability scenario','failed','missing_evidence','{"heat":"neutral"}','["calm"]','["valuation_review","discipline_check"]','["source_gap"]','hold','P83 月度复盘样本：线下已执行并需归因','["自动交易","外部推送"]','["记录复盘"]','executed_manually',NULL,NULL,'v_p83','[{"agent_name":"P83ReviewAnalyst","conclusion":"用于月度归因和纪律审计验收","key_reasons":["纪律遵守审计","证据缺失统计"],"risk_warnings":["不自动交易"],"confidence":"medium","evidence_ids":[]}]','{"precision_status":"insufficient","reason":"review sample","sample_count":5,"sample_window":"P83 monthly","screening_condition":"governance traceability","scenarios":[]}','[{"step":"review","result":"degraded"}]','{"review_period":"monthly","discipline_audit":"manual_execution"}','$NOW'),
('decision_p83_monthly_plan','req_p83_monthly_plan','consultation','159915','P83 monthly attribution planned action','completed','formal_trade_advice','normal','in_scope','P83 governance traceability scenario','satisfied','','{"heat":"low"}','["neutral"]','["risk_review"]','[]','hold','P83 月度复盘样本：记录计划','["自动交易"]','["记录计划"]','planned',NULL,NULL,'v_p83','[{"agent_name":"P83ReviewAnalyst","conclusion":"用于错误案例统计和计划动作验收","key_reasons":["记录计划","审计追踪"],"risk_warnings":["人工确认"],"confidence":"medium","evidence_ids":[]}]','{"precision_status":"insufficient","reason":"review sample","sample_count":5,"sample_window":"P83 monthly","screening_condition":"governance traceability","scenarios":[]}','[{"step":"review","result":"watch"}]','{"review_period":"monthly","discipline_audit":"planned"}','$NOW');

INSERT OR REPLACE INTO operation_confirmations (confirmation_id,decision_id,confirmation_type,operation_type,symbol,quantity,price,fees,executed_at,payload_json,note,created_at)
VALUES
('confirm_p83_exec','decision_p83_monthly_exec','executed_manually','sell','510300',100,4.2,1,'$NOW','{"p83":"monthly_attribution"}','P83 月度复盘线下执行样本','$NOW'),
('confirm_p83_plan','decision_p83_monthly_plan','planned',NULL,NULL,0,0,0,NULL,'{"p83":"monthly_attribution"}','P83 月度复盘记录计划样本','$NOW');

INSERT OR REPLACE INTO error_cases (error_case_id,decision_id,confirmation_id,actual_outcome,root_cause_tag,lesson_learned,created_at)
VALUES ('err_p83_monthly','decision_p83_monthly_exec','confirm_p83_exec','P83 月度归因：证据缺失导致降级','rule_threshold_issue','后续季度规则提案必须保留守门人和最终确认。','$NOW');

INSERT OR REPLACE INTO rule_proposals (proposal_id,proposal_type,status,source_error_case_id,title,proposal_version,before_rule_json,after_rule_json,reason,impact_scope_json,risk_notes_json,sample_count,final_confirmed_at,final_confirmed_note,applied_rule_version,related_error_cases_json,created_at)
VALUES
('prop_p83_quarterly','threshold','pending_final_confirm','err_p83_monthly','P83 季度规则效果复盘提案','v_p83_draft','{"content":"旧阈值"}','{"content":"新阈值：证据缺失进入人工复核"}','quarterly benchmark comparison and rule-effect review; review_period=quarterly:2026-Q2','{"scope":"review_summary","auto_apply":false}','["不自动应用规则","需用户最终确认"]',6,NULL,NULL,NULL,'["err_p83_monthly"]','$NOW'),
('prop_p83_master_weight','capability','pending_final_confirm','err_p83_monthly','大师权重调整提案','v_p83_master','{"target_rule":"master.graham.margin_of_safety","weight":0.25}','{"target_rule":"master.graham.margin_of_safety","weight":0.30,"proposal_subtype":"master_weight_adjustment"}','master_weight_adjustment target_rule=master.graham.margin_of_safety review_period=monthly:2026-06','{"scope":"master_wisdom_weight","auto_apply":false}','["大师经验权重只能作为提案","需守门人和用户最终确认"]',5,NULL,NULL,NULL,'["err_p83_monthly"]','$NOW');

INSERT OR REPLACE INTO gatekeeper_audits (gatekeeper_audit_id,proposal_id,audit_result,audit_reason,required_changes,violates_fundamental_rule,has_rule_conflict,backtest_metrics_json,allow_apply,audited_rule_version,created_at)
VALUES
('gate_p83_quarterly','prop_p83_quarterly','approved','P83 季度复盘提案可送最终确认，但不自动应用。','保留人工最终确认',0,0,'{"sample_count":6,"benchmark_compared":true,"passed":true}',1,'v_p83_draft','$NOW'),
('gate_p83_master','prop_p83_master_weight','approved','P83 大师权重调整提案可送最终确认，但不自动应用。','保留用户风格确认',0,0,'{"sample_count":5,"master_weight_adjustment":true,"passed":true}',1,'v_p83_master','$NOW');

INSERT OR REPLACE INTO rule_effect_validations (validation_id,proposal_id,candidate_rule_version,validation_status,sample_count,sample_window,representativeness_status,overfit_risk,replay_result,guardrail_decision,source_explanation_json,metrics_json,risk_notes_json,related_error_cases_json,related_decision_ids_json,related_risk_alert_ids_json,related_audit_event_ids_json,safety_note,created_at,updated_at)
VALUES
('validation_p83_quarterly','prop_p83_quarterly','v_p83_draft','passed',6,'2026-Q2','passed','low','passed','passed','{"review_period":"quarterly:2026-Q2","benchmark":"沪深300","source_case_count":6}','{"benchmark_return":0.02,"strategy_return":0.018,"hit_count":6,"misjudgment_count":1,"missing_evidence_count":1}','["规则效果只读验证，不自动应用"]','["err_p83_monthly"]','["decision_p83_monthly_exec","decision_p83_monthly_plan"]','["risk_smoke_p39"]','["audit_p83_review"]','规则效果验证只用于本地规则治理；规则生效仍需用户手动最终确认。','$NOW','$NOW'),
('validation_p83_master','prop_p83_master_weight','v_p83_master','passed',5,'2026-Q2','passed','low','passed','passed','{"proposal_subtype":"master_weight_adjustment","target_rule":"master.graham.margin_of_safety"}','{"source_case_count":5,"style_fit":0.8}','["大师经验权重调整只生成提案"]','["err_p83_monthly"]','["decision_p83_monthly_exec"]','[]','["audit_p83_review"]','大师权重调整提案不会自动应用。','$NOW','$NOW');

INSERT OR REPLACE INTO rule_effect_tracking (tracking_id,applied_rule_version,proposal_id,period,hit_count,misjudgment_count,missing_evidence_count,degraded_count,risk_alert_count,trend_direction,metrics_json,related_proposal_ids_json,related_audit_event_ids_json,related_risk_alert_ids_json,safety_note,created_at,updated_at)
VALUES
('track_p83_quarterly','v_p83_observed','prop_p83_quarterly','quarterly',6,1,1,1,1,'flat','{"benchmark_return":0.02,"strategy_return":0.018,"discipline_audit":"passed"}','["prop_p83_quarterly","prop_p83_master_weight"]','["audit_p83_review"]','["risk_smoke_p39"]','只读追踪，不自动应用规则。','$NOW','$NOW');

INSERT OR REPLACE INTO audit_events (audit_event_id,request_id,decision_id,workflow_type,node_name,actor,action,node_action,proposal_id,confirmation_id,error_case_id,status,error_code,before_state,after_state,rule_version,snapshot_id,input_ref_type,input_ref,output_ref_type,output_ref,created_at)
VALUES
('audit_p83_review','req_p83_review','decision_p83_monthly_exec','review_summary','P83ReviewSeed','system','run_local_task','seed_governance_traceability','prop_p83_quarterly','confirm_p83_exec','err_p83_monthly','success',NULL,NULL,'review_seeded','v_p83',NULL,'review_summary','monthly:2026-06','rule_proposal','prop_p83_quarterly','$NOW'),
('audit_p83_master','req_p83_master','decision_p83_monthly_exec','review_summary','P83MasterSeed','system','create_proposal','seed_master_weight_adjustment','prop_p83_master_weight',NULL,'err_p83_monthly','success',NULL,NULL,'proposal_seeded','v_p83',NULL,'error_case','err_p83_monthly','rule_proposal','prop_p83_master_weight','$NOW');

INSERT OR REPLACE INTO notifications (notification_id,type,severity,title,message,source_type,source_id,created_at)
VALUES ('notif_p83_review','review_degraded','warning','P83 复盘验收通知','复盘窗口存在降级、错误案例和规则提案，需要人工复核。','review_summary','monthly','$NOW');
SQL

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

P83_CAPTURE_SCREENSHOTS=1 \
P83_ARTIFACT_DIR="$ARTIFACT_DIR" \
E2E_BASE_URL="http://127.0.0.1:$WEB_PORT" \
npm --prefix "$ROOT_DIR/web" run test:e2e -- --config playwright.config.ts --workers=1 p83-governance-traceability.spec.ts

run_logged "$HANDLER_LOG" go test ./internal/application/handler -run 'Review|Notification|DailyDiscipline|Rule|Local|DataSourceQuality' -count=1
run_logged "$WORKFLOW_LOG" go test ./internal/application/workflow -run 'EvolutionProposal|Gatekeeper|DailyAutoRun|Consultation|EvidenceVerification' -count=1
run_logged "$AGENT_LOG" go test ./cmd/agent -run 'Review|Release|Diagnostics|Daily|Preflight' -count=1

if [[ "$ARTIFACT_DIR" == "$ROOT_DIR/"* ]]; then
  ARTIFACT_DIR_DISPLAY="${ARTIFACT_DIR#$ROOT_DIR/}"
else
  ARTIFACT_DIR_DISPLAY="$ARTIFACT_DIR"
fi

{
  echo "status=passed"
  echo "review_decisions=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM decision_records WHERE decision_id LIKE 'decision_p83_%';")"
  echo "review_confirmations=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_id LIKE 'confirm_p83_%';")"
  echo "review_error_cases=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM error_cases WHERE error_case_id='err_p83_monthly';")"
  echo "review_rule_proposals=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_proposals WHERE proposal_id IN ('prop_p83_quarterly','prop_p83_master_weight');")"
  echo "master_weight_proposals=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_proposals WHERE proposal_id='prop_p83_master_weight' AND proposal_type='capability' AND after_rule_json LIKE '%master_weight_adjustment%';")"
  echo "quarterly_effect_tracking=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_effect_tracking WHERE tracking_id='track_p83_quarterly' AND period='quarterly';")"
  echo "review_notifications=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM notifications WHERE type='review_degraded' AND source_type='review_summary' AND read_at IS NOT NULL;")"
  echo "review_audit_events=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM audit_events WHERE audit_event_id IN ('audit_p83_review','audit_p83_master');")"
  echo "forbidden_broker_order_push_tables=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%');")"
  echo "active_rule_versions_from_p83=$(sqlite3 "$SQLITE_PATH" "SELECT COUNT(*) FROM rule_versions WHERE COALESCE(created_from_proposal_id,'') LIKE 'prop_p83_%';")"
  echo "artifact_dir=$ARTIFACT_DIR_DISPLAY"
} >"$DB_CHECK_LOG"

python3 - "$SUMMARY_PATH" "$ARTIFACT_DIR_DISPLAY" "$DB_CHECK_LOG" "$HANDLER_LOG" "$WORKFLOW_LOG" "$AGENT_LOG" "$ARTIFACT_DIR/browser-results.json" <<'PY'
import json
import re
import sys
from pathlib import Path

summary_path, artifact_dir, db_log, handler_log, workflow_log, agent_log, browser_results = sys.argv[1:8]

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
payload = {
    "status": "passed",
    "artifact_dir": artifact_dir,
    "browser": browser,
    "db_readback": db,
    "go_tests": {
        "handler": {"status": go_status(handler_log), "log": f"{artifact_dir}/go-handler-tests.log"},
        "workflow": {"status": go_status(workflow_log), "log": f"{artifact_dir}/go-workflow-tests.log"},
        "agent": {"status": go_status(agent_log), "log": f"{artifact_dir}/go-agent-tests.log"},
    },
    "safety": {
        "forbidden_broker_order_push_tables": int(db.get("forbidden_broker_order_push_tables", "-1")),
        "active_rule_versions_from_p83": int(db.get("active_rule_versions_from_p83", "-1")),
        "claim_boundary": "P83 verifies governance traceability and related UI/API/readback; it does not refresh packages or claim broker/external push/automatic trading/automatic rule application.",
    },
}
if any(item["status"] != "passed" for item in payload["go_tests"].values()):
    payload["status"] = "failed"
if browser.get("status") != "passed" or db.get("status") != "passed":
    payload["status"] = "failed"
Path(summary_path).write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
PY

cleanup
SERVER_PID=""
WEB_PID=""
