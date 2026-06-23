#!/usr/bin/env python3
"""Generate P82 SOP/action UI-to-SQLite closure artifacts."""

from __future__ import annotations

import argparse
import json
import sqlite3
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p81-dynamic-source-field-coverage-matrix.md"
P82_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p82-sop-action-ui-sqlite-matrix.md"
P82_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p82-sop-action-ui-sqlite-closure.md"
P82_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p82-sop-action-ui-sqlite"
P82_SUMMARY = P82_ASSET_DIR / "sop-action-ui-sqlite-summary.json"
P82_BROWSER = P82_ASSET_DIR / "browser-results.json"
P82_DB_CHECK = P82_ASSET_DIR / "db-impact-check.log"
P82_SQLITE = ROOT / "tmp" / "p75-sop-failure-real-ui" / "investment-agent-p75-sop.db"
KNOWLEDGE_REGISTRY = ROOT / "internal" / "application" / "knowledge" / "registry.go"

P82_UI_COMMAND = (
    "P75_FINAL_RULE_APPLY=1 "
    "P75_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p82-sop-action-ui-sqlite "
    "bash scripts/p75-sop-failure-real-ui-acceptance.sh"
)

P82_COLUMNS = [
    "p82_status",
    "p82_closure_basis",
    "p82_fresh_evidence_command",
    "p82_fresh_evidence_artifact",
    "p82_remaining_gap",
    "p82_next_action",
]

P82_PLAN_IDS = {
    "REQ-02-004",
    "REQ-02-018",
    "REQ-02-019",
    "REQ-04-005",
    "REQ-07-011",
    "REQ-08-001",
    "REQ-08-002",
    "REQ-08-003",
    "REQ-08-005",
    "REQ-08-006",
    "REQ-08-007",
    "REQ-08-008",
    "REQ-08-009",
    "REQ-08-010",
    "REQ-08-011",
    "REQ-08-012",
    "REQ-08-013",
    "REQ-08-014",
    "REQ-08-015",
    "REQ-08-016",
    "REQ-08-017",
    "REQ-08-019",
    "REQ-08-021",
    "REQ-08-022",
    "REQ-08-024",
    "REQ-08-025",
    "REQ-08-026",
    "REQ-10-001",
    "REQ-10-002",
    "REQ-10-003",
    "REQ-10-004",
    "REQ-10-005",
    "REQ-12-001",
    "REQ-12-002",
    "REQ-12-003",
    "REQ-13-001",
    "REQ-13-002",
    "REQ-13-003",
    "REQ-13-004",
    "REQ-13-005",
    "REQ-13-007",
    "REQ-13-008",
    "REQ-13-009",
    "REQ-13-011",
    "REQ-13-012",
    "REQ-13-015",
    "REQ-13-016",
    "REQ-13-017",
    "REQ-13-019",
    "REQ-16-016",
    "REQ-16-029",
    "REQ-17-004",
    "REQ-17-010",
}

P82_DEFER_IDS = {
    "REQ-10-001": "Needs direct portfolio/allocation UI evidence for core asset target ratios, not only SOP-B action context.",
    "REQ-10-002": "Needs direct portfolio/allocation UI evidence for satellite asset target ratios, not only SOP-B action context.",
    "REQ-10-003": "Needs direct portfolio/allocation UI evidence for cash target ratios, not only SOP context.",
    "REQ-10-004": "Needs quarterly rebalance UI/API/readback evidence for preset-ratio drift handling.",
    "REQ-12-002": "Needs monthly attribution UI evidence covering P/L attribution, discipline audit, emotion log, and error-case statistics.",
    "REQ-12-003": "Needs quarterly benchmark comparison, rule-effect review, and evolution proposal summary evidence.",
    "REQ-13-011": "Needs a direct master-wisdom weight-adjustment proposal scenario, not a generic rule proposal.",
    "REQ-16-029": "Needs full main dashboard/cockpit evidence, not only P82 SOP/rules/settings routes.",
    "REQ-17-004": "Needs dashboard evidence for account state, data update time, discipline state, and triggered rules.",
}

P82_UPGRADE_IDS = P82_PLAN_IDS - set(P82_DEFER_IDS)

SOP_MAP: dict[str, dict[str, str]] = {
    "REQ-08-001": {"scenario": "SOP-A", "alert": "risk_p75_sop_a", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-002": {"scenario": "SOP-A", "alert": "risk_p75_sop_a", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-003": {"scenario": "SOP-A", "alert": "risk_p75_sop_a", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-005": {"scenario": "SOP-A", "alert": "risk_p75_sop_a", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-006": {"scenario": "SOP-B", "alert": "risk_p75_sop_b", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-007": {"scenario": "SOP-B", "alert": "risk_p75_sop_b", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-008": {"scenario": "SOP-B", "alert": "risk_p75_sop_b", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-009": {"scenario": "SOP-B", "alert": "risk_p75_sop_b", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-010": {"scenario": "SOP-B", "alert": "risk_p75_sop_b", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-011": {"scenario": "SOP-C", "alert": "risk_p75_sop_c", "ui": "/risk-alerts,/decisions/decision_p75_insufficient", "sqlite": "risk_alerts,decision_records,audit_events"},
    "REQ-08-012": {"scenario": "SOP-C", "alert": "risk_p75_sop_c", "ui": "/data-quality?symbol=999999,/decisions/decision_p75_insufficient", "sqlite": "decision_records,risk_alerts"},
    "REQ-08-013": {"scenario": "SOP-C", "alert": "risk_p75_sop_c", "ui": "/risk-alerts,/settings", "sqlite": "risk_alerts,market_snapshots"},
    "REQ-08-014": {"scenario": "SOP-C", "alert": "risk_p75_sop_c", "ui": "/risk-alerts", "sqlite": "risk_alerts"},
    "REQ-08-015": {"scenario": "SOP-C", "alert": "risk_p75_sop_c", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-016": {"scenario": "SOP-D", "alert": "risk_p75_sop_d", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-017": {"scenario": "SOP-D", "alert": "risk_p75_sop_d", "ui": "/risk-alerts,/decisions/decision_p75_model_unavailable", "sqlite": "risk_alerts,decision_records"},
    "REQ-08-019": {"scenario": "SOP-D", "alert": "risk_p75_sop_d", "ui": "/risk-alerts", "sqlite": "risk_alerts"},
    "REQ-08-021": {"scenario": "SOP-E", "alert": "risk_p75_sop_e", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-022": {"scenario": "SOP-E", "alert": "risk_p75_sop_e", "ui": "/risk-alerts", "sqlite": "risk_alerts"},
    "REQ-08-024": {"scenario": "SOP-F", "alert": "risk_p75_sop_f", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-08-025": {"scenario": "SOP-F", "alert": "risk_p75_sop_f", "ui": "/risk-alerts,/settings", "sqlite": "risk_alerts,market_snapshots"},
    "REQ-08-026": {"scenario": "SOP-F", "alert": "risk_p75_sop_f", "ui": "/risk-alerts,/audit", "sqlite": "risk_alerts,audit_events"},
}

GENERAL_MAP: dict[str, dict[str, str]] = {
    "REQ-02-004": {"scenario": "white-box-rules", "ui": "/risk-alerts,/rules,/settings", "sqlite": "risk_alerts,rule_proposals,gatekeeper_audits,market_snapshots"},
    "REQ-02-018": {"scenario": "gatekeeper-approval-required", "ui": "/rules", "sqlite": "rule_proposals,gatekeeper_audits,audit_events"},
    "REQ-02-019": {"scenario": "auditable-standards", "ui": "/risk-alerts,/audit,/rules", "sqlite": "risk_alerts,audit_events,gatekeeper_audits"},
    "REQ-04-005": {"scenario": "risk-alert-center", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-07-011": {"scenario": "stable-knowledge-ids", "ui": "/risk-alerts,/settings", "sqlite": "knowledge_registry_source,risk_alerts"},
    "REQ-10-001": {"scenario": "allocation-sop-context", "ui": "/risk-alerts", "sqlite": "risk_alerts.trigger_context_json"},
    "REQ-10-002": {"scenario": "satellite-sop-context", "ui": "/risk-alerts", "sqlite": "risk_alerts.trigger_context_json"},
    "REQ-10-003": {"scenario": "cash-buffer-sop-context", "ui": "/risk-alerts", "sqlite": "risk_alerts.trigger_context_json"},
    "REQ-10-004": {"scenario": "quarterly-review-boundary", "ui": "/rules,/review", "sqlite": "rule_effect_validations,rule_proposals"},
    "REQ-10-005": {"scenario": "satellite-limit-action", "ui": "/risk-alerts", "sqlite": "risk_alerts.suggested_actions_json"},
    "REQ-12-001": {"scenario": "daily-monitoring", "ui": "/risk-alerts,/audit", "sqlite": "risk_alerts,audit_events"},
    "REQ-12-002": {"scenario": "review-error-stats", "ui": "/review,/audit", "sqlite": "error_cases,operation_confirmations,audit_events"},
    "REQ-12-003": {"scenario": "rule-evolution-summary", "ui": "/rules,/review", "sqlite": "rule_effect_validations,rule_proposals,gatekeeper_audits"},
    "REQ-13-001": {"scenario": "error-case-field-readback", "ui": "/decisions/decision_p75_mark_error,/review", "sqlite": "error_cases.error_case_id"},
    "REQ-13-002": {"scenario": "decision-type-readback", "ui": "/decisions/decision_p75_mark_error", "sqlite": "decision_records.workflow_type,final_verdict_status"},
    "REQ-13-003": {"scenario": "context-snapshot-readback", "ui": "/decisions/decision_p75_mark_error", "sqlite": "decision_records.context_snapshot_json"},
    "REQ-13-004": {"scenario": "agent-reasoning-readback", "ui": "/decisions/decision_p75_mark_error", "sqlite": "decision_records.analyst_reports_json"},
    "REQ-13-005": {"scenario": "actual-outcome-readback", "ui": "/decisions/decision_p75_mark_error,/review", "sqlite": "error_cases.actual_outcome"},
    "REQ-13-007": {"scenario": "lesson-learned-readback", "ui": "/decisions/decision_p75_mark_error,/review", "sqlite": "error_cases.lesson_learned"},
    "REQ-13-008": {"scenario": "timestamp-readback", "ui": "/audit,/review", "sqlite": "error_cases.created_at,audit_events.created_at"},
    "REQ-13-009": {"scenario": "threshold-proposal-gatekeeper", "ui": "/rules", "sqlite": "rule_proposals,gatekeeper_audits"},
    "REQ-13-011": {"scenario": "style-weight-proposal-boundary", "ui": "/rules", "sqlite": "rule_proposals,rule_effect_validations"},
    "REQ-13-012": {"scenario": "behavior-risk-warning", "ui": "/risk-alerts,/review", "sqlite": "risk_alerts,rule_effect_validations"},
    "REQ-13-015": {"scenario": "sample-sufficiency-warning", "ui": "/rules", "sqlite": "rule_effect_validations.sample_count"},
    "REQ-13-016": {"scenario": "emotion-bias-warning", "ui": "/rules", "sqlite": "rule_effect_validations.risk_notes_json"},
    "REQ-13-017": {"scenario": "backtest-regression-gate", "ui": "/rules", "sqlite": "gatekeeper_audits.backtest_metrics_json"},
    "REQ-13-019": {"scenario": "approved-but-not-auto-applied", "ui": "/rules", "sqlite": "rule_proposals,rule_versions"},
    "REQ-16-016": {"scenario": "sop-a-f-real-ui", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"},
    "REQ-16-029": {"scenario": "web-console-real-ui", "ui": "/risk-alerts,/rules,/settings,/decisions", "sqlite": "readback-through-api"},
    "REQ-17-004": {"scenario": "dashboard-discipline-state", "ui": "/risk-alerts,/settings", "sqlite": "risk_alerts,market_snapshots"},
    "REQ-17-010": {"scenario": "user-action-recording", "ui": "/risk-alerts,/decisions/decision_p75_mark_error", "sqlite": "risk_alerts,operation_confirmations,error_cases"},
}


def split_markdown_row(line: str) -> list[str]:
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in line.rstrip("\n")[1:-1]:
        if escaped:
            current.append(char)
            escaped = False
        elif char == "\\":
            escaped = True
        elif char == "|":
            cells.append("".join(current).strip())
            current = []
        else:
            current.append(char)
    cells.append("".join(current).strip())
    return cells


def escape_cell(value: object) -> str:
    return str(value).replace("\n", " ").replace("\r", " ").strip().replace("\\", "\\\\").replace("|", "\\|")


def rel(path: Path) -> str:
    return str(path.relative_to(ROOT))


def read_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        return {"status": "missing", "path": rel(path)}
    with path.open(encoding="utf-8") as fh:
        data = json.load(fh)
    return data if isinstance(data, dict) else {"status": "invalid"}


def read_log_kv(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}
    out: dict[str, str] = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key.strip()] = value.strip().replace(str(ROOT) + "/", "")
    return out


def read_source_rows() -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in SOURCE_MATRIX.read_text(encoding="utf-8").splitlines():
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if header is None:
            if cells and cells[0] == "requirement_id":
                header = cells
            continue
        if set("".join(cells)) <= {"-", ":"}:
            continue
        if len(cells) != len(header):
            raise SystemExit(f"Invalid source matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("Source matrix header not found")
    return header, rows


def db_rows(conn: sqlite3.Connection, sql: str) -> list[dict[str, Any]]:
    conn.row_factory = sqlite3.Row
    return [dict(row) for row in conn.execute(sql).fetchall()]


def db_one(conn: sqlite3.Connection, sql: str) -> dict[str, Any] | None:
    rows = db_rows(conn, sql)
    return rows[0] if rows else None


def open_sqlite_readonly() -> sqlite3.Connection:
    return sqlite3.connect(f"file:{P82_SQLITE}?mode=ro", uri=True)


def extract_sqlite_evidence() -> dict[str, Any]:
    if not P82_SQLITE.exists():
        return {"status": "failed", "failures": ["sqlite_missing"], "sqlite_path": rel(P82_SQLITE)}

    failures: list[str] = []
    with open_sqlite_readonly() as conn:
        risk_alerts = db_rows(
            conn,
            """
            SELECT alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,
                   prohibited_actions_json,suggested_actions_json,related_decision_id,related_audit_event_id,
                   resolved_at,resolution_reason,created_at,updated_at
            FROM risk_alerts
            WHERE alert_id LIKE 'risk_p75_sop_%'
            ORDER BY alert_id
            """,
        )
        lifecycle_audits = db_rows(
            conn,
            """
            SELECT audit_event_id,request_id,node_name,actor,action,node_action,output_ref,status,after_state,created_at
            FROM audit_events
            WHERE workflow_type='risk_alert_sop'
              AND node_action='update_risk_alert_lifecycle'
              AND output_ref LIKE 'risk_p75_sop_%'
            ORDER BY output_ref
            """,
        )
        decision_mark = db_one(
            conn,
            """
            SELECT decision_id,workflow_type,symbol,final_verdict_status,confirmation_status,
                   analyst_reports_json,context_snapshot_json,created_at
            FROM decision_records
            WHERE decision_id='decision_p75_mark_error'
            """,
        )
        error_case = db_one(
            conn,
            """
            SELECT error_case_id,decision_id,confirmation_id,actual_outcome,root_cause_tag,lesson_learned,created_at
            FROM error_cases
            WHERE decision_id='decision_p75_mark_error'
            """,
        )
        confirmation = db_one(
            conn,
            """
            SELECT confirmation_id,decision_id,confirmation_type,error_case_id,note,created_at
            FROM operation_confirmations
            WHERE decision_id='decision_p75_mark_error' AND confirmation_type='marked_error'
            """,
        )
        mark_error_audit = db_one(
            conn,
            """
            SELECT audit_event_id,request_id,decision_id,actor,action,status,before_state,after_state,
                   confirmation_id,error_case_id,created_at
            FROM audit_events
            WHERE decision_id='decision_p75_mark_error' AND action='mark_error'
            """,
        )
        proposals = db_rows(
            conn,
            """
            SELECT proposal_id,proposal_type,status,title,proposal_version,reason,sample_count,
                   final_confirmed_at,applied_rule_version,related_error_cases_json,created_at
            FROM rule_proposals
            WHERE proposal_id IN ('prop_p75_ui_send_gatekeeper','prop_p75_gatekeeper_denied','prop_p75_gatekeeper_review')
            ORDER BY proposal_id
            """,
        )
        gatekeeper_audits = db_rows(
            conn,
            """
            SELECT gatekeeper_audit_id,proposal_id,audit_result,audit_reason,required_changes,
                   violates_fundamental_rule,has_rule_conflict,backtest_metrics_json,allow_apply,
                   audited_rule_version,created_at
            FROM gatekeeper_audits
            WHERE proposal_id IN ('prop_p75_ui_send_gatekeeper','prop_p75_gatekeeper_denied','prop_p75_gatekeeper_review')
            ORDER BY proposal_id, gatekeeper_audit_id
            """,
        )
        node_audits = db_rows(
            conn,
            """
            SELECT audit_event_id,request_id,workflow_type,node_name,actor,action,node_action,
                   proposal_id,status,rule_version,input_ref_type,input_ref,output_ref_type,output_ref,created_at
            FROM audit_events
            WHERE proposal_id='prop_p75_ui_send_gatekeeper'
              AND action='audit_rule_change'
              AND node_name IN ('ProposalLoadNode','FundamentalRuleCheckNode','ConflictCheckNode','BacktestNode','AuditDecisionNode','AuditRecordNode')
            ORDER BY created_at,node_name
            """,
        )
        validation = db_one(
            conn,
            """
            SELECT validation_id,proposal_id,validation_status,sample_count,representativeness_status,
                   overfit_risk,replay_result,guardrail_decision,source_explanation_json,metrics_json,
                   risk_notes_json,related_error_cases_json,related_decision_ids_json,related_risk_alert_ids_json,
                   safety_note,created_at
            FROM rule_effect_validations
            WHERE validation_id='validation_p75_review'
            """,
        )
        send_validation = db_one(
            conn,
            """
            SELECT validation_id,proposal_id,validation_status,sample_count,representativeness_status,
                   overfit_risk,replay_result,guardrail_decision,safety_note,created_at
            FROM rule_effect_validations
            WHERE validation_id='validation_p75_ui_send_gatekeeper'
            """,
        )
        final_rule_audits = db_rows(
            conn,
            """
            SELECT audit_event_id,request_id,actor,action,proposal_id,status,before_state,after_state,
                   output_ref_type,output_ref,created_at
            FROM audit_events
            WHERE proposal_id='prop_p75_ui_send_gatekeeper'
              AND action='update_rule'
              AND after_state='applied'
              AND output_ref_type='rule_version'
            ORDER BY created_at
            """,
        )
        counts = {
            "risk_alerts": len(risk_alerts),
            "sop_lifecycle_audits": len(lifecycle_audits),
            "operation_confirmations_for_mark_error": db_one(conn, "SELECT COUNT(*) AS count FROM operation_confirmations WHERE decision_id='decision_p75_mark_error' AND confirmation_type='marked_error'")["count"],
            "error_cases_for_mark_error": db_one(conn, "SELECT COUNT(*) AS count FROM error_cases WHERE decision_id='decision_p75_mark_error'")["count"],
            "mark_error_audit_events": db_one(conn, "SELECT COUNT(*) AS count FROM audit_events WHERE decision_id='decision_p75_mark_error' AND action='mark_error'")["count"],
            "gatekeeper_audits": len(gatekeeper_audits),
            "gatekeeper_node_audits": len(node_audits),
            "position_transactions_total": db_one(conn, "SELECT COUNT(*) AS count FROM position_transactions")["count"],
            "rule_versions_total": db_one(conn, "SELECT COUNT(*) AS count FROM rule_versions")["count"],
            "forbidden_broker_order_push_tables": db_one(conn, "SELECT COUNT(*) AS count FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")["count"],
        }

    alert_by_id = {row["alert_id"]: row for row in risk_alerts}
    expected_alerts = {
        "risk_p75_sop_a": ("buy_thesis_broken", "observing"),
        "risk_p75_sop_b": ("valuation_high", "escalated"),
        "risk_p75_sop_c": ("sentiment_extreme", "observing"),
        "risk_p75_sop_d": ("sentiment_extreme", "escalated"),
        "risk_p75_sop_e": ("insufficient_evidence", "resolved"),
        "risk_p75_sop_f": ("data_degraded", "escalated"),
    }
    for alert_id, (risk_type, terminal_status) in expected_alerts.items():
        alert = alert_by_id.get(alert_id)
        if not alert:
            failures.append(f"risk_alert.{alert_id}.missing")
            continue
        if alert.get("risk_type") != risk_type:
            failures.append(f"risk_alert.{alert_id}.risk_type")
        if alert.get("sop_status") != terminal_status:
            failures.append(f"risk_alert.{alert_id}.sop_status")
        if not alert.get("trigger_context_json") or not alert.get("prohibited_actions_json") or not alert.get("suggested_actions_json"):
            failures.append(f"risk_alert.{alert_id}.action_context")
    if len(lifecycle_audits) < 6:
        failures.append("lifecycle_audits.count")

    if not decision_mark:
        failures.append("decision_mark_missing")
    elif not decision_mark.get("analyst_reports_json") or not decision_mark.get("context_snapshot_json"):
        failures.append("decision_mark.reasoning_or_context")
    if not error_case:
        failures.append("error_case_missing")
    else:
        for key in ["error_case_id", "decision_id", "confirmation_id", "actual_outcome", "root_cause_tag", "lesson_learned", "created_at"]:
            if not error_case.get(key):
                failures.append(f"error_case.{key}")
    if not confirmation or (error_case and confirmation.get("error_case_id") != error_case.get("error_case_id")):
        failures.append("confirmation_error_case_link")
    if not mark_error_audit or mark_error_audit.get("after_state") != "marked_error":
        failures.append("mark_error_audit")

    proposal_by_id = {row["proposal_id"]: row for row in proposals}
    send = proposal_by_id.get("prop_p75_ui_send_gatekeeper")
    if not send or send.get("status") != "applied":
        failures.append("final_rule_apply_status")
    if not send or not send.get("final_confirmed_at") or not send.get("applied_rule_version"):
        failures.append("final_rule_apply_readback")
    gate_by_proposal = {row["proposal_id"]: row for row in gatekeeper_audits}
    expected_gate = {
        "prop_p75_ui_send_gatekeeper": "approved",
        "prop_p75_gatekeeper_denied": "rejected",
        "prop_p75_gatekeeper_review": "needs_user_review",
    }
    for proposal_id, result in expected_gate.items():
        audit = gate_by_proposal.get(proposal_id)
        if not audit or audit.get("audit_result") != result or not audit.get("backtest_metrics_json"):
            failures.append(f"gatekeeper.{proposal_id}")
    if len(node_audits) != 6:
        failures.append("gatekeeper_node_audits.count")
    if not validation:
        failures.append("rule_effect_validation_missing")
    elif int(validation.get("sample_count") or 0) >= 3:
        failures.append("rule_effect_validation_sample_gate")
    if not send_validation:
        failures.append("send_rule_effect_validation_missing")
    elif send_validation.get("validation_status") != "passed" or send_validation.get("guardrail_decision") != "passed":
        failures.append("send_rule_effect_validation_not_passed")
    if len(final_rule_audits) != 1:
        failures.append("final_rule_audit_count")
    if counts["rule_versions_total"] < 2:
        failures.append("final_rule_version_missing")
    if counts["position_transactions_total"] != 0:
        failures.append("position_transactions_created")
    if counts["forbidden_broker_order_push_tables"] != 0:
        failures.append("forbidden_tables")

    registry_text = KNOWLEDGE_REGISTRY.read_text(encoding="utf-8")
    stable_ids = [
        "master.graham.margin_of_safety",
        "discipline.no_single_source_decision",
        "risk_sop.evidence_insufficient",
        "risk_sop.valuation_high",
        "symbol_profile.510300",
        "symbol_profile.159915",
    ]
    missing_ids = [item for item in stable_ids if item not in registry_text]
    if missing_ids:
        failures.append("knowledge_registry.stable_ids")

    return {
        "status": "passed" if not failures else "failed",
        "sqlite_path": rel(P82_SQLITE),
        "failures": failures,
        "counts": counts,
        "field_readback": {
            "risk_alerts": risk_alerts,
            "sop_lifecycle_audits": lifecycle_audits,
            "decision_mark_error": decision_mark or {},
            "error_case": error_case or {},
            "operation_confirmation": confirmation or {},
            "mark_error_audit": mark_error_audit or {},
            "rule_proposals": proposals,
            "gatekeeper_audits": gatekeeper_audits,
            "gatekeeper_node_audits": node_audits,
            "rule_effect_validation": validation or {},
            "send_rule_effect_validation": send_validation or {},
            "final_rule_audits": final_rule_audits,
            "stable_knowledge_ids": stable_ids,
            "missing_stable_knowledge_ids": missing_ids,
        },
    }


def row_evidence_map() -> dict[str, dict[str, str]]:
    out: dict[str, dict[str, str]] = {}
    for rid in sorted(P82_PLAN_IDS):
        base = SOP_MAP.get(rid) or GENERAL_MAP.get(rid)
        if not base:
            base = {"scenario": "p82-sop-action", "ui": "/risk-alerts", "sqlite": "risk_alerts,audit_events"}
        out[rid] = {
            "scenario": base["scenario"],
            "real_ui_paths": base["ui"],
            "sqlite_or_source_readback": base["sqlite"],
            "safety_negative_check": "no broker/order/external-push tables, no position_transactions in SOP journey, no automatic rule application, no trading confirmation",
            "p82_upgrade_decision": "deferred" if rid in P82_DEFER_IDS else "real_pass_candidate",
            "defer_reason": P82_DEFER_IDS.get(rid, ""),
        }
    return out


def summarize_evidence() -> dict[str, Any]:
    browser = read_json(P82_BROWSER)
    db_log = read_log_kv(P82_DB_CHECK)
    sqlite_evidence = extract_sqlite_evidence()
    failures: list[str] = []
    if browser.get("status") != "passed":
        failures.append("browser.status")
    if len(browser.get("sop_scenarios") or []) != 6:
        failures.append("browser.sop_scenarios")
    if len(browser.get("failure_states") or []) < 8:
        failures.append("browser.failure_states")
    if db_log.get("status") != "passed":
        failures.append("db_log.status")
    if sqlite_evidence.get("status") != "passed":
        failures.extend(f"sqlite.{item}" for item in sqlite_evidence.get("failures", []))
    return {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "status": "passed" if not failures else "failed",
        "failures": failures,
        "ui_command": P82_UI_COMMAND,
        "browser_results": rel(P82_BROWSER),
        "db_impact_log": rel(P82_DB_CHECK),
        "sqlite_summary": sqlite_evidence,
        "browser_summary": {
            "status": browser.get("status"),
            "generated_at": browser.get("generated_at"),
            "sop_scenarios": browser.get("sop_scenarios") or [],
            "failure_states": browser.get("failure_states") or [],
        },
        "db_log_summary": db_log,
        "row_evidence_map": row_evidence_map(),
    }


def p82_values(row: dict[str, str], evidence_ok: bool) -> dict[str, str]:
    rid = row["requirement_id"]
    if row["p81_status"] == "real_pass":
        return {
            "p82_status": "real_pass",
            "p82_closure_basis": "Carried forward from P81/P80/P79/P78/P77 real_pass; P82 does not rewrite prior evidence.",
            "p82_fresh_evidence_command": "N/A",
            "p82_fresh_evidence_artifact": row.get("p81_fresh_evidence_artifact", "N/A"),
            "p82_remaining_gap": "None for row already accepted before P82.",
            "p82_next_action": "Keep covered by future regression evidence.",
        }
    if rid in P82_UPGRADE_IDS and evidence_ok:
        return {
            "p82_status": "real_pass",
            "p82_closure_basis": "Fresh P82 real browser SOP/action journey plus field-level API/SQLite readback cover this exact UI operation, SOP status, audit, confirmation, review, gatekeeper, and safety-boundary behavior.",
            "p82_fresh_evidence_command": f"{P82_UI_COMMAND} && python3 scripts/p82_sop_action_ui_sqlite_closure.py --check",
            "p82_fresh_evidence_artifact": rel(P82_SUMMARY),
            "p82_remaining_gap": "None for this SOP/action UI-to-SQLite row; broker execution, external push, and automatic trading remain out of scope and absent.",
            "p82_next_action": "Keep in real UI SOP/action regression.",
        }
    if rid in P82_UPGRADE_IDS:
        return {
            "p82_status": row["p81_status"],
            "p82_closure_basis": "P82 upgrade candidate, but fresh real UI/API/SQLite evidence is not complete.",
            "p82_fresh_evidence_command": f"{P82_UI_COMMAND} && python3 scripts/p82_sop_action_ui_sqlite_closure.py --check",
            "p82_fresh_evidence_artifact": rel(P82_SUMMARY),
            "p82_remaining_gap": "Fresh P82 SOP/action evidence did not pass.",
            "p82_next_action": "Fix evidence failures and rerun P82.",
        }
    if rid in P82_DEFER_IDS:
        return {
            "p82_status": row["p81_status"],
            "p82_closure_basis": "P82 evaluated this row but did not upgrade it because the fresh SOP/action evidence would overclaim the original requirement.",
            "p82_fresh_evidence_command": P82_UI_COMMAND,
            "p82_fresh_evidence_artifact": rel(P82_SUMMARY),
            "p82_remaining_gap": P82_DEFER_IDS[rid],
            "p82_next_action": "Carry this row into the later portfolio, governance, dashboard, or core-goal closure batch with a direct scenario.",
        }
    return {
        "p82_status": row["p81_status"],
        "p82_closure_basis": "No P82 upgrade; row is owned by a later P83-P86 execution batch or remains previously scoped/reference.",
        "p82_fresh_evidence_command": "N/A",
        "p82_fresh_evidence_artifact": "N/A",
        "p82_remaining_gap": row.get("p81_remaining_gap", "Needs post-P82 evidence."),
        "p82_next_action": row.get("p81_next_action", "Create follow-up implementation or acceptance work."),
    }


def write_matrix(header: list[str], rows: list[dict[str, str]], evidence_ok: bool, write: bool) -> list[dict[str, str]]:
    columns = header + P82_COLUMNS
    output_rows: list[dict[str, str]] = []
    for row in rows:
        merged = dict(row)
        merged.update(p82_values(row, evidence_ok))
        output_rows.append(merged)

    counts = Counter(row["p82_status"] for row in output_rows)
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p82_status"] == "real_pass"]
    remaining = len(full_rows) - len(full_real)
    remaining_by_group = Counter(
        row.get("remediation_group", "unknown")
        for row in full_rows
        if row["p82_status"] != "real_pass"
    )
    lines = [
        "# P82 SOP Action UI SQLite Matrix",
        "",
        f"> Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
        f"> Source: `{rel(SOURCE_MATRIX)}`",
        "> Policy: P82 is a new evidence layer; it does not rewrite P75-P81 history.",
        "",
        "## Status Summary",
        "",
    ]
    for status, count in sorted(counts.items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend(
        [
            "",
            f"- full_release_required rows: {len(full_rows)}",
            f"- full_release_required `real_pass` rows: {len(full_real)}",
            f"- remaining full_release_required non-`real_pass` rows: {remaining}",
            f"- P82 evaluated rows: {len(P82_PLAN_IDS)}",
            f"- new P82 `real_pass` rows: {len([row for row in output_rows if row['requirement_id'] in P82_UPGRADE_IDS and row['p82_status'] == 'real_pass'])}",
            f"- P82 evaluated but deferred rows: {len(P82_DEFER_IDS)}",
            "",
            "## Remaining Non-Real-Pass By Remediation Group",
            "",
        ]
    )
    for group, count in sorted(remaining_by_group.items()):
        lines.append(f"- {group}: {count}")
    lines.extend(["", "## Matrix", "", "|" + "|".join(columns) + "|", "|" + "|".join(["---"] * len(columns)) + "|"])
    for row in output_rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in columns) + "|")
    if write:
        P82_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return output_rows


def write_acceptance(evidence: dict[str, Any], rows: list[dict[str, str]], write: bool) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p82_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in rows if row["requirement_id"] in P82_UPGRADE_IDS and row["p82_status"] == "real_pass"]
    lines = [
        "# P82 SOP Action UI SQLite Closure Acceptance",
        "",
        "> Date: 2026-06-22",
        "> Change: `p82-sop-action-ui-sqlite-closure`",
        "> Conclusion: `release_ready_scoped_with_p82_sop_action_progress`",
        "",
        "## Summary",
        "",
        f"- Source matrix: `{rel(SOURCE_MATRIX)}`",
        f"- P82 matrix: `{rel(P82_MATRIX)}`",
        f"- Summary JSON: `{rel(P82_SUMMARY)}`",
        f"- Full-release-required rows: {len(full_rows)}",
        f"- Full-release-required `real_pass` rows after P82: {len(full_real)}",
        f"- Remaining full-release-required non-`real_pass` rows: {len(full_rows) - len(full_real)}",
        f"- Evaluated by P82: {len(P82_PLAN_IDS)}",
        f"- Newly upgraded by P82: {len(upgraded)}",
        f"- Evaluated but deferred by P82: {len(P82_DEFER_IDS)}",
        "",
        "## P82 Upgrades",
        "",
    ]
    lines.extend(f"- `{rid}`" for rid in upgraded)
    lines.extend(["", "## P82 Evaluated But Deferred", ""])
    lines.extend(f"- `{rid}`: {reason}" for rid, reason in sorted(P82_DEFER_IDS.items()))
    lines.extend(
        [
            "",
            "## Fresh Real UI Evidence",
            "",
            f"- Artifact directory: `{rel(P82_ASSET_DIR)}`",
            f"- Browser result: `{rel(P82_BROWSER)}`",
            f"- DB impact log: `{rel(P82_DB_CHECK)}`",
            f"- Field-level summary: `{rel(P82_SUMMARY)}`",
            "",
            "Command:",
            "",
            "```bash",
            P82_UI_COMMAND,
            "python3 scripts/p82_sop_action_ui_sqlite_closure.py --check",
            "```",
            "",
            "## Field-Level Evidence Covered",
            "",
            "- SOP A-F are operated through the real browser `/risk-alerts` UI, then read back through page refresh and SQLite.",
            "- Risk alert lifecycle actions update `risk_alerts.sop_status` and write `audit_events` lifecycle rows.",
            "- Failure-state UI covers unsupported symbols, insufficient formal evidence, stale/degraded sources, model unavailability, validation errors, gatekeeper denial, and gatekeeper user-review states.",
            "- Mark-error UI creates exactly one local confirmation, exactly one error case, and a linked audit event with before/after state.",
            "- Rule proposal UI sends a proposal through the gatekeeper graph, records node-level audits, and applies a local rule version only after explicit user final confirmation.",
            "- Stable built-in knowledge IDs are verified from the registry source, including master, discipline, risk SOP, and symbol-profile IDs.",
            "- Forbidden broker/order/external-push tables are absent, no position transaction is created, and the journey does not confirm a trade.",
            "",
            "## Remaining Gaps",
            "",
            "P82 evaluated all 53 planned SOP/action rows, but only rows with direct fresh evidence were upgraded. Deferred rows remain owned by P83-P86 or a later direct scenario. P82 does not claim broker execution, external push, automatic trading, automatic rule application, or full original-requirement pass.",
            "",
            "## Evidence Status",
            "",
            f"- status: `{evidence['status']}`",
            f"- failures: `{', '.join(evidence['failures']) if evidence['failures'] else 'none'}`",
            "",
            "## Boundaries",
            "",
            "- P82 does not rewrite P75-P81 historical matrices.",
            "- P82 does not refresh distribution packages; a later package refresh is required before claiming package inclusion.",
            "- P82 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.",
        ]
    )
    if write:
        P82_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def verify_plan_rows(rows: list[dict[str, str]]) -> None:
    candidates = {
        row["requirement_id"]
        for row in rows
        if row.get("remediation_group") == "sop_action_data_impact"
        and row.get("p81_status") != "real_pass"
        and row.get("full_release_required") == "True"
    }
    if candidates != P82_PLAN_IDS:
        missing = sorted(candidates - P82_PLAN_IDS)
        extra = sorted(P82_PLAN_IDS - candidates)
        raise SystemExit(f"P82 row-set mismatch: missing={missing} extra={extra}")


def run(check_only: bool) -> dict[str, Any]:
    header, rows = read_source_rows()
    verify_plan_rows(rows)
    evidence = summarize_evidence()
    evidence_ok = evidence["status"] == "passed"
    output_rows = write_matrix(header, rows, evidence_ok, write=not check_only)
    write_acceptance(evidence, output_rows, write=not check_only)
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p82_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in output_rows if row["requirement_id"] in P82_UPGRADE_IDS and row["p82_status"] == "real_pass"]
    summary = {
        "generated_utc": evidence["generated_utc"],
        "source_matrix": rel(SOURCE_MATRIX),
        "matrix": rel(P82_MATRIX),
        "acceptance": rel(P82_ACCEPTANCE),
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "evaluated_rows": len(P82_PLAN_IDS),
        "newly_upgraded_rows": len(upgraded),
        "newly_upgraded_requirement_ids": upgraded,
        "deferred_rows": len(P82_DEFER_IDS),
        "deferred_requirement_ids": sorted(P82_DEFER_IDS),
        "conclusion": "release_ready_scoped_with_p82_sop_action_progress",
        "fresh_evidence": evidence,
    }
    if not check_only:
        P82_ASSET_DIR.mkdir(parents=True, exist_ok=True)
        P82_SUMMARY.write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    if check_only and not evidence_ok:
        raise SystemExit(f"P82 evidence failed: {evidence['failures']}")
    if check_only and len(upgraded) != len(P82_UPGRADE_IDS):
        raise SystemExit(f"P82 expected {len(P82_UPGRADE_IDS)} upgraded rows, got {len(upgraded)}")
    return summary


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="fail if evidence is incomplete")
    args = parser.parse_args()
    summary = run(check_only=args.check)
    print(
        "p82_sop_action_ui_sqlite:"
        f"status={summary['fresh_evidence']['status']}:"
        f"new_real={summary['newly_upgraded_rows']}:"
        f"remaining_full={summary['remaining_full_release_required_non_real_pass_rows']}"
    )


if __name__ == "__main__":
    main()
