#!/usr/bin/env python3
import json
import sqlite3
import sys
from pathlib import Path


def fail(message, details=None):
    payload = {"status": "failed", "message": message}
    if details is not None:
        payload["details"] = details
    print(json.dumps(payload, ensure_ascii=False, indent=2), file=sys.stderr)
    raise SystemExit(1)


def require(condition, message, details=None):
    if not condition:
        fail(message, details)


def scalar(db, sql, params=()):
    row = db.execute(sql, params).fetchone()
    return row[0] if row else None


def one(db, sql, params=()):
    row = db.execute(sql, params).fetchone()
    return dict(row) if row else {}


def rows(db, sql, params=()):
    return [dict(row) for row in db.execute(sql, params).fetchall()]


def main():
    if len(sys.argv) != 3:
        fail("usage: p73_effect_replay_check.py <sqlite_path> <artifact_dir>")

    sqlite_path = Path(sys.argv[1])
    artifact_dir = Path(sys.argv[2])
    require(sqlite_path.exists(), "sqlite database does not exist", str(sqlite_path))
    artifact_dir.mkdir(parents=True, exist_ok=True)

    db = sqlite3.connect(f"file:{sqlite_path}?mode=ro", uri=True)
    db.row_factory = sqlite3.Row

    background_decision = one(
        db,
        "SELECT decision_id,record_type,workflow_status,source_verification_status,final_verdict_status,confirmation_status "
        "FROM decision_records WHERE decision_id='decision_smoke_p73_background_only'",
    )
    require(background_decision, "missing P73 background-only decision")
    require(background_decision["record_type"] == "non_trade_record", "background-only decision must be non-trade", background_decision)
    require(background_decision["source_verification_status"] == "background_only", "background-only decision verification mismatch", background_decision)
    require(background_decision["final_verdict_status"] == "insufficient_data", "background-only decision must be insufficient_data", background_decision)
    require(background_decision["confirmation_status"] == "not_required", "background-only decision must not expose confirmation", background_decision)

    background_verification = one(
        db,
        "SELECT verification_id,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level "
        "FROM source_verifications WHERE verification_id='verification_smoke_p73_background'",
    )
    require(background_verification, "missing P73 background-only source verification")
    require(background_verification["evidence_role"] == "background", "background verification role mismatch", background_verification)
    require(background_verification["verification_status"] == "background_only", "C-level background material must not be satisfied", background_verification)
    require(background_verification["independent_source_count"] == 0, "background material must not count as independent formal source", background_verification)
    require(background_verification["high_grade_independent_source_count"] == 0, "background material must not count as high-grade source", background_verification)

    p73_summary = one(
        db,
        "SELECT summary_id,source_level,evidence_role,summary "
        "FROM intelligence_summary WHERE summary_id='summary_smoke_p73_background'",
    )
    require(p73_summary, "missing P73 background evidence summary")
    require(p73_summary["source_level"] == "C" and p73_summary["evidence_role"] == "background", "P73 evidence summary must remain C/background", p73_summary)

    planned_confirmation = one(
        db,
        "SELECT confirmation_id,decision_id,confirmation_type,symbol,note FROM operation_confirmations "
        "WHERE decision_id='decision_smoke_p30' AND confirmation_type='planned' ORDER BY created_at DESC LIMIT 1",
    )
    require(planned_confirmation, "missing planned confirmation from UX task")
    linked_transactions = scalar(
        db,
        "SELECT COUNT(*) FROM position_transactions tx JOIN operation_confirmations c ON c.confirmation_id=tx.confirmation_id "
        "WHERE c.decision_id='decision_smoke_p30' AND c.confirmation_type='planned'",
    )
    require(linked_transactions == 0, "planned confirmation must not mutate portfolio facts", {"linked_transactions": linked_transactions, "planned_confirmation": planned_confirmation})

    risk_alert = one(
        db,
        "SELECT alert_id,severity,sop_status,symbol,prohibited_actions_json,related_decision_id,related_report_id "
        "FROM risk_alerts WHERE alert_id='risk_smoke_p39'",
    )
    require(risk_alert, "missing risk alert fixture")
    require("自动交易" in (risk_alert["prohibited_actions_json"] or ""), "risk alert must expose prohibited actions", risk_alert)

    rule_effect = one(
        db,
        "SELECT validation_id,validation_status,sample_count,representativeness_status,overfit_risk,replay_result,guardrail_decision,safety_note "
        "FROM rule_effect_validations WHERE validation_id='val_smoke_p39'",
    )
    require(rule_effect, "missing rule effect validation fixture")
    require(rule_effect["validation_status"] == "passed", "rule effect validation status mismatch", rule_effect)
    require(rule_effect["sample_count"] >= 5, "rule effect validation must expose sample count", rule_effect)
    require(rule_effect["overfit_risk"] == "low", "rule effect validation must expose overfit risk", rule_effect)
    require("不自动应用规则" in rule_effect["safety_note"], "rule effect safety note missing no-auto-apply boundary", rule_effect)

    audit_count = scalar(db, "SELECT COUNT(*) FROM audit_events WHERE audit_event_id='audit_smoke_p73_effectiveness'")
    require(audit_count == 1, "missing P73 effectiveness audit event", {"audit_count": audit_count})

    forbidden_tables = rows(
        db,
        "SELECT name FROM sqlite_master WHERE type='table' AND ("
        "LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%trade_execution%' "
        "OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%webhook%')",
    )
    require(not forbidden_tables, "forbidden trading or external-push table exists", forbidden_tables)

    summary = {
        "status": "passed",
        "sqlite_path": str(sqlite_path),
        "background_decision": background_decision,
        "background_verification": background_verification,
        "background_summary": p73_summary,
        "planned_confirmation": planned_confirmation,
        "planned_confirmation_linked_transactions": linked_transactions,
        "risk_alert": risk_alert,
        "rule_effect": rule_effect,
        "p73_audit_event_count": audit_count,
        "forbidden_tables": forbidden_tables,
        "effectiveness_claim_boundary": "discipline behavior replay only; no future return or market-direction claim",
    }
    (artifact_dir / "effect-replay-summary.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
