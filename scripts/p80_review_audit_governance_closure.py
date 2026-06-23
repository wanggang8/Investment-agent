#!/usr/bin/env python3
"""Generate P80 review/audit/governance real-use closure artifacts."""

from __future__ import annotations

import argparse
import json
import re
import sqlite3
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
P79_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md"
P80_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p80-review-audit-governance-matrix.md"
P80_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p80-review-audit-governance-closure.md"
P80_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p80-review-audit-governance"
P80_SUMMARY = P80_ASSET_DIR / "review-audit-governance-summary.json"
P80_BROWSER = P80_ASSET_DIR / "browser-results.json"
P80_DB_CHECK = P80_ASSET_DIR / "db-impact-check.log"
P80_SQLITE = ROOT / "tmp" / "p75-sop-failure-real-ui" / "investment-agent-p75-sop.db"

P80_UI_COMMAND = (
    "P75_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance "
    "bash scripts/p75-sop-failure-real-ui-acceptance.sh"
)

P80_COLUMNS = [
    "p80_status",
    "p80_closure_basis",
    "p80_fresh_evidence_command",
    "p80_fresh_evidence_artifact",
    "p80_remaining_gap",
    "p80_next_action",
]

P80_UPGRADE_IDS = {
    "REQ-04-007",
    "REQ-04-020",
    "REQ-11-018",
    "REQ-13-006",
    "REQ-13-013",
    "REQ-13-014",
    "REQ-13-018",
    "REQ-13-020",
    "REQ-13-021",
    "REQ-16-024",
    "REQ-16-027",
    "REQ-17-016",
    "REQ-17-018",
    "REQ-17-019",
}

P80_GUARD_IDS = {
    "REQ-01-006",
    "REQ-02-006",
    "REQ-03-010",
    "REQ-04-008",
    "REQ-04-014",
    "REQ-04-015",
    "REQ-12-002",
    "REQ-12-003",
    "REQ-13-009",
    "REQ-13-010",
    "REQ-13-011",
    "REQ-14-001",
    "REQ-14-005",
    "REQ-14-007",
    "REQ-15-009",
    "REQ-16-026",
    "REQ-16-028",
    "REQ-16-033",
    "REQ-17-010",
    "REQ-17-017",
}

GATEKEEPER_NODES = {
    "ProposalLoadNode": "load_proposal",
    "FundamentalRuleCheckNode": "check_fundamental_rules",
    "ConflictCheckNode": "check_rule_conflict",
    "BacktestNode": "check_backtest_samples",
    "AuditDecisionNode": "decide_gatekeeper_audit",
    "AuditRecordNode": "record_gatekeeper_audit",
}


def split_markdown_row(line: str) -> list[str]:
    stripped = line.rstrip("\n")
    if not stripped.startswith("|") or not stripped.endswith("|"):
        return []
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in stripped[1:-1]:
        if escaped:
            current.append(char)
            escaped = False
            continue
        if char == "\\":
            escaped = True
            continue
        if char == "|":
            cells.append("".join(current))
            current = []
            continue
        current.append(char)
    cells.append("".join(current))
    return [cell.strip() for cell in cells]


def escape_cell(value: object) -> str:
    text = str(value).replace("\n", " ").replace("\r", " ").strip()
    return text.replace("\\", "\\\\").replace("|", "\\|")


def rel(path: Path) -> str:
    return str(path.relative_to(ROOT))


def read_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        return {"status": "missing", "path": rel(path)}
    with path.open(encoding="utf-8") as fh:
        data = json.load(fh)
    return data if isinstance(data, dict) else {"status": "invalid", "path": rel(path)}


def read_log_kv(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}
    result: dict[str, str] = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        if "=" not in line:
            continue
        key, value = line.split("=", 1)
        result[key.strip()] = value.strip().replace(str(ROOT) + "/", "")
    return result


def read_p79_rows() -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in P79_MATRIX.read_text(encoding="utf-8").splitlines():
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if not cells:
            continue
        if header is None:
            if cells[0] == "requirement_id":
                header = cells
            continue
        if set("".join(cells)) <= {"-", ":"}:
            continue
        if len(cells) != len(header):
            raise SystemExit(f"Invalid P79 matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("P79 matrix header not found")
    return header, rows


def db_rows(conn: sqlite3.Connection, sql: str, params: tuple[Any, ...] = ()) -> list[dict[str, Any]]:
    conn.row_factory = sqlite3.Row
    cursor = conn.execute(sql, params)
    return [dict(row) for row in cursor.fetchall()]


def db_one(conn: sqlite3.Connection, sql: str, params: tuple[Any, ...] = ()) -> dict[str, Any] | None:
    rows = db_rows(conn, sql, params)
    return rows[0] if rows else None


def redact_row(row: dict[str, Any], keys: list[str]) -> dict[str, Any]:
    return {key: row.get(key) for key in keys}


def extract_sqlite_evidence() -> dict[str, Any]:
    if not P80_SQLITE.exists():
        return {"status": "failed", "failures": ["sqlite_missing"], "sqlite_path": rel(P80_SQLITE)}

    failures: list[str] = []
    with sqlite3.connect(P80_SQLITE) as conn:
        error_case = db_one(
            conn,
            """
            SELECT error_case_id, decision_id, confirmation_id, actual_outcome, root_cause_tag, lesson_learned, created_at
            FROM error_cases
            WHERE decision_id='decision_p75_mark_error'
            """,
        )
        confirmation = db_one(
            conn,
            """
            SELECT confirmation_id, decision_id, confirmation_type, error_case_id, note, created_at
            FROM operation_confirmations
            WHERE decision_id='decision_p75_mark_error' AND confirmation_type='marked_error'
            """,
        )
        mark_error_audit = db_one(
            conn,
            """
            SELECT audit_event_id, request_id, decision_id, actor, action, confirmation_id, error_case_id,
                   status, before_state, after_state, created_at
            FROM audit_events
            WHERE decision_id='decision_p75_mark_error' AND action='mark_error'
            """,
        )
        proposals = db_rows(
            conn,
            """
            SELECT proposal_id, proposal_type, status, title, proposal_version, reason,
                   sample_count, final_confirmed_at, applied_rule_version, created_at
            FROM rule_proposals
            WHERE proposal_id IN ('prop_p75_ui_send_gatekeeper','prop_p75_gatekeeper_denied','prop_p75_gatekeeper_review')
            ORDER BY proposal_id
            """,
        )
        gatekeeper = db_rows(
            conn,
            """
            SELECT gatekeeper_audit_id, proposal_id, audit_result, audit_reason, required_changes,
                   violates_fundamental_rule, has_rule_conflict, backtest_metrics_json,
                   allow_apply, audited_rule_version, created_at
            FROM gatekeeper_audits
            WHERE proposal_id IN ('prop_p75_ui_send_gatekeeper','prop_p75_gatekeeper_denied','prop_p75_gatekeeper_review')
            ORDER BY proposal_id, gatekeeper_audit_id
            """,
        )
        node_audits = db_rows(
            conn,
            """
            SELECT audit_event_id, request_id, workflow_type, node_name, actor, action, node_action,
                   proposal_id, status, rule_version, input_ref_type, input_ref, output_ref_type, output_ref, created_at
            FROM audit_events
            WHERE proposal_id='prop_p75_ui_send_gatekeeper'
              AND action='audit_rule_change'
              AND node_name IN ('ProposalLoadNode','FundamentalRuleCheckNode','ConflictCheckNode','BacktestNode','AuditDecisionNode','AuditRecordNode')
            ORDER BY created_at, node_name
            """,
        )
        user_gatekeeper_audit = db_one(
            conn,
            """
            SELECT audit_event_id, request_id, actor, action, proposal_id, status, before_state, after_state, created_at
            FROM audit_events
            WHERE proposal_id='prop_p75_ui_send_gatekeeper'
              AND action='audit_rule_change'
              AND actor='user'
            """,
        )
        risk_alerts = db_rows(
            conn,
            """
            SELECT alert_id, risk_type, severity, sop_status, symbol, trigger_summary,
                   prohibited_actions_json, suggested_actions_json, resolved_at, resolution_reason, updated_at
            FROM risk_alerts
            WHERE alert_id LIKE 'risk_p75_sop_%'
            ORDER BY alert_id
            """,
        )
        lifecycle_audits = db_rows(
            conn,
            """
            SELECT audit_event_id, request_id, node_name, actor, action, node_action, output_ref, status, after_state, created_at
            FROM audit_events
            WHERE workflow_type='risk_alert_sop'
              AND node_action='update_risk_alert_lifecycle'
              AND output_ref LIKE 'risk_p75_sop_%'
            ORDER BY output_ref
            """,
        )
        forbidden_tables = db_rows(
            conn,
            """
            SELECT name FROM sqlite_master
            WHERE type='table'
              AND (
                LOWER(name) LIKE '%broker%'
                OR LOWER(name) LIKE '%order%'
                OR LOWER(name) LIKE '%external_push%'
                OR LOWER(name) LIKE '%push%'
                OR LOWER(name) LIKE '%trade_execution%'
              )
            ORDER BY name
            """,
        )
        counts = {
            "operation_confirmations_for_mark_error": db_one(conn, "SELECT COUNT(*) AS count FROM operation_confirmations WHERE decision_id='decision_p75_mark_error' AND confirmation_type='marked_error'")["count"],
            "error_cases_for_mark_error": db_one(conn, "SELECT COUNT(*) AS count FROM error_cases WHERE decision_id='decision_p75_mark_error'")["count"],
            "mark_error_audit_events": db_one(conn, "SELECT COUNT(*) AS count FROM audit_events WHERE decision_id='decision_p75_mark_error' AND action='mark_error'")["count"],
            "gatekeeper_audits_for_p80": len(gatekeeper),
            "gatekeeper_node_audits": len(node_audits),
            "sop_risk_alerts": len(risk_alerts),
            "sop_lifecycle_audits": len(lifecycle_audits),
            "forbidden_broker_order_push_tables": len(forbidden_tables),
            "position_transactions_total": db_one(conn, "SELECT COUNT(*) AS count FROM position_transactions")["count"],
            "rule_versions_total": db_one(conn, "SELECT COUNT(*) AS count FROM rule_versions")["count"],
        }

    if not error_case:
        failures.append("error_case_missing")
    else:
        for key in ["decision_id", "confirmation_id", "actual_outcome", "root_cause_tag", "lesson_learned", "created_at"]:
            if not error_case.get(key):
                failures.append(f"error_case.{key}")
        if error_case.get("root_cause_tag") != "rule_threshold_issue":
            failures.append("error_case.root_cause_tag_expected")

    if not confirmation:
        failures.append("confirmation_missing")
    elif error_case and confirmation.get("error_case_id") != error_case.get("error_case_id"):
        failures.append("confirmation.error_case_link")

    if not mark_error_audit:
        failures.append("mark_error_audit_missing")
    else:
        expected = {"actor": "user", "action": "mark_error", "status": "success", "before_state": "pending", "after_state": "marked_error"}
        for key, value in expected.items():
            if mark_error_audit.get(key) != value:
                failures.append(f"mark_error_audit.{key}")
        if confirmation and mark_error_audit.get("confirmation_id") != confirmation.get("confirmation_id"):
            failures.append("mark_error_audit.confirmation_link")
        if error_case and mark_error_audit.get("error_case_id") != error_case.get("error_case_id"):
            failures.append("mark_error_audit.error_case_link")

    proposal_by_id = {row["proposal_id"]: row for row in proposals}
    send = proposal_by_id.get("prop_p75_ui_send_gatekeeper")
    if not send:
        failures.append("proposal.send_missing")
    else:
        if send.get("status") != "pending_final_confirm":
            failures.append("proposal.send_status")
        if send.get("final_confirmed_at") is not None or send.get("applied_rule_version") is not None:
            failures.append("proposal.final_application_boundary")

    gate_by_id = {row["proposal_id"]: row for row in gatekeeper}
    expected_gatekeeper = {
        "prop_p75_ui_send_gatekeeper": ("approved", 1, 0, 0),
        "prop_p75_gatekeeper_denied": ("rejected", 0, 1, 1),
        "prop_p75_gatekeeper_review": ("needs_user_review", 0, 0, 0),
    }
    for proposal_id, (result, allow_apply, fundamental, conflict) in expected_gatekeeper.items():
        audit = gate_by_id.get(proposal_id)
        if not audit:
            failures.append(f"gatekeeper.{proposal_id}.missing")
            continue
        if audit.get("audit_result") != result:
            failures.append(f"gatekeeper.{proposal_id}.result")
        if int(audit.get("allow_apply") or 0) != allow_apply:
            failures.append(f"gatekeeper.{proposal_id}.allow_apply")
        if int(audit.get("violates_fundamental_rule") or 0) != fundamental:
            failures.append(f"gatekeeper.{proposal_id}.fundamental_rule")
        if int(audit.get("has_rule_conflict") or 0) != conflict:
            failures.append(f"gatekeeper.{proposal_id}.conflict")
        if not audit.get("backtest_metrics_json"):
            failures.append(f"gatekeeper.{proposal_id}.backtest_metrics")

    node_by_name = {row["node_name"]: row for row in node_audits}
    for node_name, node_action in GATEKEEPER_NODES.items():
        audit = node_by_name.get(node_name)
        if not audit:
            failures.append(f"node_audit.{node_name}.missing")
            continue
        if audit.get("node_action") != node_action:
            failures.append(f"node_audit.{node_name}.node_action")
        if audit.get("actor") != "system" or audit.get("status") != "success":
            failures.append(f"node_audit.{node_name}.status")
        if not audit.get("request_id") or audit.get("proposal_id") != "prop_p75_ui_send_gatekeeper":
            failures.append(f"node_audit.{node_name}.trace_refs")
    if not user_gatekeeper_audit:
        failures.append("gatekeeper_user_audit_missing")
    elif user_gatekeeper_audit.get("before_state") != "pending_user_confirm" or user_gatekeeper_audit.get("after_state") != "under_gatekeeper_audit":
        failures.append("gatekeeper_user_audit.state_transition")

    terminal_statuses = {"observing", "escalated", "resolved"}
    if len(risk_alerts) != 6:
        failures.append("risk_alerts.count")
    if len(lifecycle_audits) < 6:
        failures.append("risk_alerts.lifecycle_audits")
    if any(row.get("sop_status") not in terminal_statuses for row in risk_alerts):
        failures.append("risk_alerts.terminal_status")
    if counts["forbidden_broker_order_push_tables"] != 0:
        failures.append("forbidden_tables")
    if counts["position_transactions_total"] != 0:
        failures.append("position_transactions_created")
    if counts["rule_versions_total"] != 1:
        failures.append("rule_versions_changed")

    field_readback = {
        "error_case": redact_row(error_case or {}, ["error_case_id", "decision_id", "confirmation_id", "actual_outcome", "root_cause_tag", "lesson_learned", "created_at"]),
        "operation_confirmation": redact_row(confirmation or {}, ["confirmation_id", "decision_id", "confirmation_type", "error_case_id", "note", "created_at"]),
        "mark_error_audit": mark_error_audit or {},
        "rule_proposals": proposals,
        "gatekeeper_audits": gatekeeper,
        "gatekeeper_node_audits": node_audits,
        "gatekeeper_user_audit": user_gatekeeper_audit or {},
        "risk_alerts": risk_alerts,
        "sop_lifecycle_audits": lifecycle_audits,
        "forbidden_tables": forbidden_tables,
    }
    return {
        "status": "passed" if not failures else "failed",
        "sqlite_path": rel(P80_SQLITE),
        "failures": failures,
        "counts": counts,
        "field_readback": field_readback,
    }


def summarize_evidence() -> dict[str, Any]:
    browser = read_json(P80_BROWSER)
    db_log = read_log_kv(P80_DB_CHECK)
    sqlite_evidence = extract_sqlite_evidence()
    failures: list[str] = []
    if browser.get("status") != "passed":
        failures.append("browser.status")
    if browser.get("failure_states") is None or len(browser.get("failure_states") or []) < 8:
        failures.append("browser.failure_states")
    if browser.get("sop_scenarios") is None or len(browser.get("sop_scenarios") or []) != 6:
        failures.append("browser.sop_scenarios")
    if db_log.get("status") != "passed":
        failures.append("db_log.status")
    if sqlite_evidence.get("status") != "passed":
        failures.extend(f"sqlite.{item}" for item in sqlite_evidence.get("failures", []))
    return {
        "status": "passed" if not failures else "failed",
        "failures": failures,
        "ui_command": P80_UI_COMMAND,
        "browser_results": rel(P80_BROWSER),
        "db_impact_log": rel(P80_DB_CHECK),
        "sqlite_summary": sqlite_evidence,
        "browser_summary": {
            "status": browser.get("status"),
            "generated_at": browser.get("generated_at"),
            "sop_scenarios": browser.get("sop_scenarios") or [],
            "failure_states": browser.get("failure_states") or [],
        },
        "db_log_summary": db_log,
    }


def p80_values(row: dict[str, str], evidence_ok: bool) -> dict[str, str]:
    rid = row["requirement_id"]
    if row["p79_status"] == "real_pass":
        return {
            "p80_status": "real_pass",
            "p80_closure_basis": "Carried forward from P79/P78/P77 real_pass; P80 does not rewrite prior evidence.",
            "p80_fresh_evidence_command": "N/A",
            "p80_fresh_evidence_artifact": row.get("p79_fresh_evidence_artifact", "N/A"),
            "p80_remaining_gap": "None for row already accepted before P80.",
            "p80_next_action": "Keep covered by future regression evidence.",
        }
    if rid in P80_UPGRADE_IDS and evidence_ok:
        return {
            "p80_status": "real_pass",
            "p80_closure_basis": "Fresh P80 real browser review/audit/rules journey plus field-level SQLite readback cover this exact review/audit/governance behavior.",
            "p80_fresh_evidence_command": P80_UI_COMMAND,
            "p80_fresh_evidence_artifact": rel(P80_SUMMARY),
            "p80_remaining_gap": "None for this row's review/audit/governance claim.",
            "p80_next_action": "Keep in real UI review/audit/governance regression.",
        }
    if rid in P80_UPGRADE_IDS and not evidence_ok:
        return {
            "p80_status": row["p79_status"],
            "p80_closure_basis": "P80 upgrade candidate, but fresh evidence is not complete.",
            "p80_fresh_evidence_command": P80_UI_COMMAND,
            "p80_fresh_evidence_artifact": rel(P80_SUMMARY),
            "p80_remaining_gap": "Fresh P80 real UI review/audit/governance evidence did not pass.",
            "p80_next_action": "Fix evidence failures and rerun P80 checker.",
        }
    if rid in P80_GUARD_IDS:
        return {
            "p80_status": row["p79_status"],
            "p80_closure_basis": "No P80 upgrade; row is broader than the directly proven review/audit/governance field evidence.",
            "p80_fresh_evidence_command": "N/A",
            "p80_fresh_evidence_artifact": "N/A",
            "p80_remaining_gap": "Needs direct proof for full rule-proposal generation, monthly/quarterly attribution, final rule application time, or complete flow coverage as applicable.",
            "p80_next_action": "Create a dedicated real UI/data-impact scenario before upgrading this broad row.",
        }
    return {
        "p80_status": row["p79_status"],
        "p80_closure_basis": "No P80 upgrade; classified for a later closure batch.",
        "p80_fresh_evidence_command": "N/A",
        "p80_fresh_evidence_artifact": "N/A",
        "p80_remaining_gap": row.get("p79_remaining_gap", "Needs post-P80 evidence."),
        "p80_next_action": row.get("p79_next_action", "Create follow-up implementation or acceptance work."),
    }


def write_matrix(header: list[str], rows: list[dict[str, str]], evidence_ok: bool) -> list[dict[str, str]]:
    columns = header + P80_COLUMNS
    output_rows: list[dict[str, str]] = []
    for row in rows:
        merged = dict(row)
        merged.update(p80_values(row, evidence_ok))
        output_rows.append(merged)

    counts = Counter(row["p80_status"] for row in output_rows)
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p80_status"] == "real_pass"]
    remaining = len(full_rows) - len(full_real)
    lines = [
        "# P80 Review Audit Governance Matrix",
        "",
        "> Generated: 2026-06-22",
        f"> Source: `{rel(P79_MATRIX)}`",
        "> Policy: P80 is a new evidence layer; it does not rewrite P75, P77, P78, or P79 history.",
        "",
        "## Status Summary",
        "",
    ]
    for status, count in sorted(counts.items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend([
        "",
        f"- full_release_required rows: {len(full_rows)}",
        f"- full_release_required real_pass rows: {len(full_real)}",
        f"- remaining full_release_required non-real-pass rows: {remaining}",
        "- conclusion: `release_ready_scoped_with_p80_review_audit_governance_progress`",
        "",
        "## Atomic Requirement Batch Rows",
        "",
        "|" + "|".join(columns) + "|",
        "|" + "|".join(["---"] * len(columns)) + "|",
    ])
    for row in output_rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in columns) + "|")
    P80_MATRIX.parent.mkdir(parents=True, exist_ok=True)
    P80_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return output_rows


def write_summary(rows: list[dict[str, str]], evidence: dict[str, Any]) -> dict[str, Any]:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p80_status"] == "real_pass"]
    upgraded = [
        row["requirement_id"]
        for row in rows
        if row["p80_status"] == "real_pass" and row.get("p79_status") != "real_pass"
    ]
    payload = {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "source_matrix": rel(P79_MATRIX),
        "matrix": rel(P80_MATRIX),
        "acceptance": rel(P80_ACCEPTANCE),
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "newly_upgraded_rows": len(upgraded),
        "newly_upgraded_requirement_ids": upgraded,
        "conclusion": "release_ready_scoped_with_p80_review_audit_governance_progress",
        "fresh_evidence": evidence,
        "not_claimed": [
            "full original-requirement pass",
            "P80 evidence inside any existing P76 distribution archive",
            "monthly or quarterly P&L attribution completeness",
            "final rule application time or applied rule version",
            "automatic rule application",
            "broker connectivity or automatic trading",
            "rule proposal generation from every error case",
            "future return or provider availability",
        ],
    }
    P80_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P80_SUMMARY.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def write_acceptance(summary: dict[str, Any]) -> None:
    upgraded = summary["newly_upgraded_requirement_ids"]
    lines = [
        "# P80 Review Audit Governance Real-Use Closure Acceptance",
        "",
        "> Date: 2026-06-22",
        "> Change: `p80-review-audit-governance-real-use-closure`",
        "> Conclusion: `release_ready_scoped_with_p80_review_audit_governance_progress`",
        "",
        "## Summary",
        "",
        f"- Source matrix: `{rel(P79_MATRIX)}`",
        f"- P80 matrix: `{rel(P80_MATRIX)}`",
        f"- Summary JSON: `{rel(P80_SUMMARY)}`",
        f"- Full-release-required rows: {summary['full_release_required_rows']}",
        f"- Full-release-required `real_pass` rows after P80: {summary['full_release_required_real_pass_rows']}",
        f"- Remaining full-release-required non-`real_pass` rows: {summary['remaining_full_release_required_non_real_pass_rows']}",
        f"- Newly upgraded by P80: {summary['newly_upgraded_rows']}",
        "",
        "## P80 Upgrades",
        "",
    ]
    lines.extend(f"- `{rid}`" for rid in upgraded)
    lines.extend([
        "",
        "## Fresh Real UI Evidence",
        "",
        f"- Artifact directory: `{rel(P80_ASSET_DIR)}`",
        f"- Browser result: `{rel(P80_BROWSER)}`",
        f"- DB impact log: `{rel(P80_DB_CHECK)}`",
        f"- Field-level summary: `{rel(P80_SUMMARY)}`",
        "",
        "Command:",
        "",
        "```bash",
        P80_UI_COMMAND,
        "python3 scripts/p80_review_audit_governance_closure.py --check",
        "```",
        "",
        "## Field-Level Evidence Covered",
        "",
        "- Mark-error UI creates exactly one `operation_confirmations` row and exactly one `error_cases` row.",
        "- `error_cases` readback includes `decision_id`, `confirmation_id`, `actual_outcome`, `root_cause_tag`, `lesson_learned`, and `created_at`.",
        "- Mark-error audit event readback includes user actor, action, status, before/after state, request id, confirmation id, and error-case id.",
        "- Rule proposal UI confirmation sends the proposal to gatekeeper audit and stops at `pending_final_confirm` without applying a rule version.",
        "- Gatekeeper readback covers `approved`, `rejected`, and `needs_user_review` states, with fundamental-rule, conflict, backtest, decision, and audit-record node events.",
        "- SOP A-F UI actions update risk alert statuses and create lifecycle audit events.",
        "- Forbidden broker/order/external-push tables are absent, and no position transaction is created by this review/audit/governance journey.",
        "",
        "## Remaining Gaps",
        "",
        "P80 deliberately does not upgrade broad rows requiring full monthly or quarterly attribution, final rule application time, automatic proposal generation from every error case, or complete original-requirement pass. Those rows remain scoped/partial until a dedicated real UI/data-impact scenario proves the exact behavior.",
        "",
        "## Boundaries",
        "",
        "- P80 does not rewrite P75, P77, P78, or P79 historical matrices.",
        "- P80 does not refresh the P76 package; a separate package refresh is required before claiming distribution archives include P80 materials.",
        "- P80 does not claim full original-requirement pass while any full-release-required row remains non-`real_pass`.",
        "- P80 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.",
    ])
    P80_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def sanitize_artifacts() -> None:
    if not P80_ASSET_DIR.exists():
        return
    for path in P80_ASSET_DIR.rglob("*"):
        if not path.is_file() or path.suffix.lower() not in {".json", ".log", ".txt", ".md"}:
            continue
        text = path.read_text(encoding="utf-8")
        sanitized = text.replace(str(ROOT) + "/", "")
        if sanitized != text:
            path.write_text(sanitized, encoding="utf-8")


def validate_rows(rows: list[dict[str, str]], evidence_ok: bool) -> None:
    upgraded = {row["requirement_id"] for row in rows if row["p80_status"] == "real_pass" and row.get("p79_status") != "real_pass"}
    invalid = upgraded - P80_UPGRADE_IDS
    if invalid:
        raise SystemExit(f"P80 invalid upgraded rows: {sorted(invalid)}")
    if P80_UPGRADE_IDS & upgraded and not evidence_ok:
        raise SystemExit("P80 review/audit/governance rows upgraded without passing evidence")
    guarded = P80_GUARD_IDS & upgraded
    if guarded:
        raise SystemExit(f"P80 guarded broad rows upgraded without exact proof: {sorted(guarded)}")
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    if all(row["p80_status"] == "real_pass" for row in full_rows):
        raise SystemExit("P80 unexpectedly claims all full-release-required rows are real_pass; review full-pass claim gate")


def validate_claim_text() -> None:
    scan_files = [P80_MATRIX, P80_ACCEPTANCE, P80_SUMMARY]
    for directory in [P80_ASSET_DIR]:
        if directory.exists():
            scan_files.extend(path for path in directory.rglob("*") if path.suffix.lower() in {".json", ".log", ".txt", ".md"})
    forbidden = [
        "release_ready_full_requirements_traceable",
        "status is full original-requirement pass",
        "conclusion is full original-requirement pass",
        "P76 package includes P80",
        "automatic trading is supported",
        "broker connectivity is supported",
    ]
    private_path = re.compile(r"/Users/[^\\s`\"']+")
    secret_like = re.compile(r"(sk-[A-Za-z0-9_-]{12,}|api[_-]?key['\"]?\\s*[:=]\\s*['\"][^'\"]{8,})", re.IGNORECASE)
    for path in scan_files:
        if not path.exists() or not path.is_file():
            continue
        text = path.read_text(encoding="utf-8")
        for term in forbidden:
            if term in text:
                raise SystemExit(f"P80 forbidden overbroad claim in {rel(path)}: {term}")
        if private_path.search(text):
            raise SystemExit(f"P80 private absolute path leak in {rel(path)}")
        if secret_like.search(text):
            raise SystemExit(f"P80 secret-like value leak in {rel(path)}")


def run(check: bool) -> None:
    sanitize_artifacts()
    header, source_rows = read_p79_rows()
    evidence = summarize_evidence()
    evidence_ok = evidence["status"] == "passed"
    rows = write_matrix(header, source_rows, evidence_ok)
    summary = write_summary(rows, evidence)
    write_acceptance(summary)
    sanitize_artifacts()
    validate_rows(rows, evidence_ok)
    validate_claim_text()
    if check and not evidence_ok:
        raise SystemExit(f"P80 evidence failed: {evidence['failures']}")
    print(
        "P80 review/audit/governance closure generated: "
        f"rows={len(rows)} real_pass={summary['full_release_required_real_pass_rows']} "
        f"new={summary['newly_upgraded_rows']} "
        f"remaining_full={summary['remaining_full_release_required_non_real_pass_rows']} "
        f"conclusion={summary['conclusion']}"
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="fail if fresh P80 evidence is incomplete")
    args = parser.parse_args()
    run(check=args.check)


if __name__ == "__main__":
    main()
