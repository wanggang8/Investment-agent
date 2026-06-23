#!/usr/bin/env python3
"""Read-only SQLite checks for P88 acceptance."""

from __future__ import annotations

import json
import sqlite3
import sys
from pathlib import Path


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def scalar(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> int | str:
    row = conn.execute(sql, args).fetchone()
    if row is None or row[0] is None:
        return 0
    return row[0]


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: p88_sqlite_readback_check.py <sqlite> <browser-results.json> <artifact-dir>")
    db_path = Path(sys.argv[1])
    browser = json.loads(Path(sys.argv[2]).read_text(encoding="utf-8"))
    artifact_dir = sys.argv[3]
    decision_id = browser.get("decision", {}).get("decision_id")
    proposal_id = browser.get("proposal", {}).get("proposal_id")
    source_transitions = browser.get("source_transitions") or []
    require(decision_id, "missing_browser_decision_id")
    require(proposal_id, "missing_browser_proposal_id")
    require(len(source_transitions) >= 2, "missing_source_transition_results")

    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    decision = conn.execute(
        "SELECT decision_id,symbol,expected_return_scenarios_json FROM decision_records WHERE decision_id=?",
        (decision_id,),
    ).fetchone()
    require(decision is not None, "missing_decision")
    expected = json.loads(decision["expected_return_scenarios_json"] or "{}")
    require(expected.get("target_name") == "沪深300ETF", "expected_target_name")
    require(expected.get("target_code") == "510300", "expected_target_code")
    require(expected.get("horizon_label") == "未来 12 个月", "expected_horizon")
    require(expected.get("probability_basis") == "historical_similar_sample_proportion", "probability_basis")
    probs = [item.get("probability") if "probability" in item else item.get("Probability") for item in expected.get("scenarios") or []]
    require(probs == [0.2, 0.6, 0.2], f"probabilities:{probs}")
    coverage = [item.get("holding_class") for item in expected.get("holding_class_coverage") or []]
    require("equity_constituent_financial" in coverage, "holding_class_coverage")

    transitions_by_scenario = {item.get("scenario"): item for item in source_transitions}
    expected_transitions = {
        "source_verified_buy_logic_break_sell_only": ("sell_only", "新增买入", "加仓"),
        "single_high_grade_major_event_frozen_watch": ("frozen_watch", "主动交易建议", None),
    }
    sell_only_decision_id = ""
    for scenario, (want_status, required_action, second_action) in expected_transitions.items():
        item = transitions_by_scenario.get(scenario)
        require(item and item.get("decision_id"), f"missing_transition:{scenario}")
        if scenario == "source_verified_buy_logic_break_sell_only":
            sell_only_decision_id = item["decision_id"]
        row = conn.execute(
            "SELECT decision_id,symbol,final_verdict_status,prohibited_actions_json FROM decision_records WHERE decision_id=?",
            (item["decision_id"],),
        ).fetchone()
        require(row is not None, f"missing_transition_decision:{scenario}")
        require(row["final_verdict_status"] == want_status, f"transition_status:{scenario}:{row['final_verdict_status']}")
        actions = json.loads(row["prohibited_actions_json"] or "[]")
        require(required_action in actions, f"transition_action:{scenario}:{required_action}")
        if second_action:
            require(second_action in actions, f"transition_action:{scenario}:{second_action}")

    evidence = conn.execute(
        "SELECT COUNT(*) AS evidence_count, COUNT(DISTINCT source_name) AS source_count, MIN(high_grade_independent_source_count) AS min_high_grade FROM evidence_refs WHERE decision_id=? AND evidence_role='formal'",
        (sell_only_decision_id,),
    ).fetchone()
    require(evidence is not None, "sell_only_evidence_missing")
    require(int(evidence["evidence_count"]) >= 2, f"sell_only_evidence_count:{evidence['evidence_count']}")
    require(int(evidence["source_count"]) >= 2, f"sell_only_source_count:{evidence['source_count']}")
    require(int(evidence["min_high_grade"] or 0) >= 2, f"sell_only_high_grade:{evidence['min_high_grade']}")

    sell_only_expected_row = conn.execute(
        "SELECT expected_return_scenarios_json FROM decision_records WHERE decision_id=?",
        (sell_only_decision_id,),
    ).fetchone()
    require(sell_only_expected_row is not None, "sell_only_expected_return_missing")
    low_sample_expected = json.loads(sell_only_expected_row["expected_return_scenarios_json"] or "{}")
    require(low_sample_expected.get("sample_count") == 2, f"low_sample_count:{low_sample_expected.get('sample_count')}")
    require((low_sample_expected.get("scenarios") or []) == [], "low_sample_scenarios_should_be_empty")
    require("样本过少" in (low_sample_expected.get("reason") or ""), "low_sample_reason")
    supplement = low_sample_expected.get("supplement_data") or []
    for category in ("market_history", "valuation_percentiles", "fundamental_growth", "formal_evidence"):
        require(category in supplement, f"low_sample_supplement:{category}")

    rebalance_audits = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM audit_events WHERE action='run_local_task' AND input_ref_type='rebalance_review'",
        )
    )
    require(rebalance_audits >= 1, "rebalance_audit")

    proposal = conn.execute(
        "SELECT proposal_id,proposal_type,status,after_rule_json,risk_notes_json,applied_rule_version FROM rule_proposals WHERE proposal_id=?",
        (proposal_id,),
    ).fetchone()
    require(proposal is not None, "missing_proposal")
    require(proposal["proposal_type"] == "sop", "proposal_type")
    require(proposal["status"] == "pending_user_confirm", "proposal_status")
    require(not proposal["applied_rule_version"], "proposal_must_not_apply")
    require('"auto_apply":false' in (proposal["after_rule_json"] or ""), "proposal_auto_apply_boundary")
    require("不会自动应用规则" in (proposal["risk_notes_json"] or ""), "proposal_risk_notes")
    notification_count = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM notifications WHERE source_type='rule_proposal' AND source_id=?",
            (proposal_id,),
        )
    )
    audit_count = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM audit_events WHERE action='create_proposal' AND proposal_id=?",
            (proposal_id,),
        )
    )
    require(notification_count == 1, "proposal_notification")
    require(audit_count == 1, "proposal_audit")

    forbidden_table_count = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')",
        )
    )
    auto_confirmation_rows = int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'"))
    require(forbidden_table_count == 0, "forbidden_broker_order_push_tables")
    require(auto_confirmation_rows == 0, "auto_confirmation_rows")

    print("status=passed")
    print(f"decision_id={decision_id}")
    print(f"proposal_id={proposal_id}")
    print("probability_basis=historical_similar_sample_proportion")
    print("probabilities=0.2,0.6,0.2")
    print("holding_class_coverage=equity_constituent_financial")
    print("source_transition_sell_only=passed")
    print("source_transition_frozen_watch=passed")
    print(f"source_transition_sell_only_evidence_count={int(evidence['evidence_count'])}")
    print(f"source_transition_sell_only_source_count={int(evidence['source_count'])}")
    print("sample_below_5_sqlite_readback=passed")
    print(f"rebalance_audits={rebalance_audits}")
    print(f"proposal_notifications={notification_count}")
    print(f"proposal_audits={audit_count}")
    print(f"forbidden_broker_order_push_tables={forbidden_table_count}")
    print(f"auto_confirmation_rows={auto_confirmation_rows}")
    print(f"artifact_dir={artifact_dir}")


if __name__ == "__main__":
    main()
