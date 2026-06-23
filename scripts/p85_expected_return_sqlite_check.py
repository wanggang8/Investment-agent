#!/usr/bin/env python3
"""Read-only SQLite checks for P85 expected-return acceptance."""

from __future__ import annotations

import json
import sqlite3
import sys
from pathlib import Path


def scalar(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> int | str:
    row = conn.execute(sql, args).fetchone()
    if row is None or row[0] is None:
        return 0
    return row[0]


def read_expected(conn: sqlite3.Connection, decision_id: str) -> tuple[sqlite3.Row, dict]:
    row = conn.execute(
        "SELECT decision_id,symbol,workflow_status,final_verdict_status,confirmation_status,expected_return_scenarios_json,analyst_reports_json,context_snapshot_json FROM decision_records WHERE decision_id=?",
        (decision_id,),
    ).fetchone()
    if row is None:
        raise SystemExit(f"status=failed\nreason=missing_decision:{decision_id}")
    return row, json.loads(row["expected_return_scenarios_json"] or "{}")


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def scenario_value(item: dict, snake_key: str, workflow_key: str):
    return item.get(snake_key) if snake_key in item else item.get(workflow_key)


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: p85_expected_return_sqlite_check.py <sqlite> <browser-results.json> <artifact-dir>")
    db_path = Path(sys.argv[1])
    browser_results = json.loads(Path(sys.argv[2]).read_text(encoding="utf-8"))
    artifact_dir = sys.argv[3]
    decisions = browser_results.get("decisions", {})
    ids = {
        name: decisions.get(name, {}).get("decision_id")
        for name in ("available", "downside", "unavailable")
    }
    require(all(ids.values()), "missing_browser_decision_ids")

    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    available_row, available = read_expected(conn, ids["available"])
    downside_row, downside = read_expected(conn, ids["downside"])
    unavailable_row, unavailable = read_expected(conn, ids["unavailable"])

    scenarios = available.get("scenarios") or []
    require(available.get("precision_status") == "available", "available_precision_status")
    require(available.get("sample_count") == 20, "available_sample_count")
    require(len(scenarios) == 3, "available_scenario_count")
    require([scenario_value(item, "return_range", "ReturnRange") for item in scenarios] == ["8.00%~15.00%", "0.00%~8.00%", "-12.00%~0.00%"], "available_ranges")
    require([scenario_value(item, "probability", "Probability") for item in scenarios] == [0.25, 0.5, 0.25], "available_probabilities_storage_shape")
    triggers = available.get("sell_evaluation", {}).get("triggers") or []
    for trigger in ("upside_lower_bound_reached", "base_upper_bound_exceeded", "base_midpoint_downshift", "target_return_reached"):
        require(trigger in triggers, f"missing_available_trigger:{trigger}")
    require(available.get("reassessment_trigger", {}).get("boundary") == "base_midpoint_downshift", "reassessment_trigger")
    require("不构成收益承诺" in (available.get("sell_evaluation", {}).get("non_trading_disclaimer") or ""), "non_trading_disclaimer")

    downside_triggers = downside.get("sell_evaluation", {}).get("triggers") or []
    require(downside.get("precision_status") == "available", "downside_precision_status")
    require(downside.get("sample_count") == 20, "downside_sample_count")
    require("downside_lower_bound_breached" in downside_triggers, "downside_trigger")

    require(unavailable.get("precision_status") == "unavailable", "unavailable_precision_status")
    require(unavailable.get("sample_count") == 1, "unavailable_sample_count")
    require(len(unavailable.get("scenarios") or []) == 0, "unavailable_scenarios_empty")
    require(unavailable.get("sell_evaluation", {}).get("status") == "not_applicable", "unavailable_sell_eval")

    for row in (available_row, downside_row, unavailable_row):
        require(row["workflow_status"] in ("completed", "degraded"), f"workflow_status:{row['decision_id']}")
        require(row["confirmation_status"] in ("pending", "not_required"), f"confirmation_status:{row['decision_id']}")

    p85_confirmations = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM operation_confirmations WHERE decision_id IN (?,?,?)",
            (ids["available"], ids["downside"], ids["unavailable"]),
        )
    )
    forbidden_table_count = int(
        scalar(
            conn,
            "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')",
        )
    )
    auto_confirmation_count = int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'"))
    require(p85_confirmations == 0, "p85_must_not_create_confirmations")
    require(forbidden_table_count == 0, "forbidden_broker_order_push_tables")
    require(auto_confirmation_count == 0, "auto_confirmation_rows")

    print("status=passed")
    print(f"available_decision_id={ids['available']}")
    print(f"downside_decision_id={ids['downside']}")
    print(f"unavailable_decision_id={ids['unavailable']}")
    print("available_precision_status=available")
    print("available_sample_count=20")
    print("available_probability_sum=1.00")
    print(f"available_triggers={','.join(triggers)}")
    print(f"downside_triggers={','.join(downside_triggers)}")
    print("unavailable_precision_status=unavailable")
    print("unavailable_sample_count=1")
    print(f"operation_confirmations_p85={p85_confirmations}")
    print(f"forbidden_broker_order_push_tables={forbidden_table_count}")
    print(f"auto_confirmation_rows={auto_confirmation_count}")
    print(f"artifact_dir={artifact_dir}")


if __name__ == "__main__":
    main()
