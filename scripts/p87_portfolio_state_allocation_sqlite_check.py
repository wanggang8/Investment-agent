#!/usr/bin/env python3
"""Read-only SQLite checks for P87 portfolio state/allocation acceptance."""

from __future__ import annotations

import json
import sqlite3
import sys
from pathlib import Path


EXPECTED_POSITIONS = {
    "510300": {
        "market_value": 64000.0,
        "asset_tag": "core",
        "position_state": "normal",
        "buy_date": "2026-01-05",
    },
    "159915": {
        "market_value": 27000.0,
        "asset_tag": "satellite",
        "position_state": "sell_only",
        "buy_date": "2026-01-06",
    },
    "511880": {
        "market_value": 1000.0,
        "asset_tag": "cash",
        "position_state": "frozen_watch",
        "buy_date": "2026-01-07",
    },
}


def scalar(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> float | int | str:
    row = conn.execute(sql, args).fetchone()
    if row is None:
        return 0
    return row[0] if row[0] is not None else 0


def close_enough(value: float, expected: float, tolerance: float = 0.01) -> bool:
    return abs(value - expected) <= tolerance


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: p87_portfolio_state_allocation_sqlite_check.py <sqlite> <browser-results.json> <artifact-dir>")
    db_path = Path(sys.argv[1])
    browser_results = Path(sys.argv[2])
    artifact_dir = sys.argv[3]
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    browser = json.loads(browser_results.read_text(encoding="utf-8"))
    if browser.get("status") != "passed":
        raise SystemExit("status=failed\nreason=browser_results_not_passed")

    snapshot = conn.execute("SELECT cash,total_assets,cash_ratio,high_risk_ratio,position_count FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1").fetchone()
    if snapshot is None:
        raise SystemExit("status=failed\nreason=missing_latest_snapshot")
    if not close_enough(float(snapshot["cash"]), 8000.0) or not close_enough(float(snapshot["total_assets"]), 100000.0):
        raise SystemExit("status=failed\nreason=snapshot_cash_or_total_mismatch")
    if not close_enough(float(snapshot["cash_ratio"]), 0.08, 0.0001):
        raise SystemExit("status=failed\nreason=cash_ratio_mismatch")
    if int(snapshot["position_count"]) != 3:
        raise SystemExit("status=failed\nreason=position_count_mismatch")

    rows = conn.execute("SELECT symbol,market_value,asset_tag,position_state,buy_date FROM positions WHERE symbol IN ('510300','159915','511880') ORDER BY symbol").fetchall()
    if len(rows) != 3:
        raise SystemExit("status=failed\nreason=missing_expected_positions")

    total_assets = float(snapshot["total_assets"])
    bucket_values = {"core": 0.0, "satellite": 0.0, "cash": float(snapshot["cash"])}
    states: dict[str, str] = {}
    buy_dates: dict[str, str] = {}
    for row in rows:
        expected = EXPECTED_POSITIONS[row["symbol"]]
        if row["asset_tag"] != expected["asset_tag"]:
            raise SystemExit(f"status=failed\nreason=asset_tag_mismatch_{row['symbol']}")
        if row["position_state"] != expected["position_state"]:
            raise SystemExit(f"status=failed\nreason=position_state_mismatch_{row['symbol']}")
        if row["buy_date"] != expected["buy_date"]:
            raise SystemExit(f"status=failed\nreason=buy_date_mismatch_{row['symbol']}")
        if not close_enough(float(row["market_value"]), expected["market_value"]):
            raise SystemExit(f"status=failed\nreason=market_value_mismatch_{row['symbol']}")
        bucket_values[row["asset_tag"]] = bucket_values.get(row["asset_tag"], 0.0) + float(row["market_value"])
        states[row["symbol"]] = row["position_state"]
        buy_dates[row["symbol"]] = row["buy_date"]

    core_ratio = bucket_values["core"] / total_assets
    satellite_ratio = bucket_values["satellite"] / total_assets
    cash_ratio = float(snapshot["cash_ratio"])
    cash_bucket_ratio = bucket_values["cash"] / total_assets
    if not (0.60 <= core_ratio <= 0.70):
        raise SystemExit("status=failed\nreason=core_ratio_out_of_target")
    if not (0.20 <= satellite_ratio <= 0.30):
        raise SystemExit("status=failed\nreason=satellite_ratio_out_of_target")
    if not (0.05 <= cash_ratio <= 0.10):
        raise SystemExit("status=failed\nreason=cash_ratio_out_of_target")
    if not (0.05 <= cash_bucket_ratio <= 0.10):
        raise SystemExit("status=failed\nreason=cash_bucket_ratio_out_of_target")

    decision_expectations = {
        "decision_p87_sell_only": ("sell_only", "formal_trade_advice", "pending"),
        "decision_p87_frozen_watch": ("frozen_watch", "non_trade_record", "not_required"),
        "decision_p87_insufficient": ("insufficient_data", "non_trade_record", "not_required"),
    }
    for decision_id, expected in decision_expectations.items():
        row = conn.execute(
            "SELECT final_verdict_status,record_type,confirmation_status FROM decision_records WHERE decision_id=?",
            (decision_id,),
        ).fetchone()
        if row is None:
            raise SystemExit(f"status=failed\nreason=missing_{decision_id}")
        if (row["final_verdict_status"], row["record_type"], row["confirmation_status"]) != expected:
            raise SystemExit(f"status=failed\nreason=decision_state_mismatch_{decision_id}")

    frozen_confirmations = scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE decision_id IN ('decision_p87_frozen_watch','decision_p87_insufficient')")
    if int(frozen_confirmations) != 0:
        raise SystemExit("status=failed\nreason=frozen_or_insufficient_confirmation_created")

    print("status=passed")
    print(f"snapshot_cash={float(snapshot['cash']):.2f}")
    print(f"snapshot_total_assets={total_assets:.2f}")
    print(f"snapshot_cash_ratio={cash_ratio:.4f}")
    print(f"core_ratio={core_ratio:.4f}")
    print(f"satellite_ratio={satellite_ratio:.4f}")
    print(f"cash_bucket_ratio={cash_bucket_ratio:.4f}")
    print(f"position_count={len(rows)}")
    print(f"position_states_json={json.dumps(states, ensure_ascii=False, sort_keys=True)}")
    print(f"buy_dates_json={json.dumps(buy_dates, ensure_ascii=False, sort_keys=True)}")
    print(f"decision_p87_sell_only_status={decision_expectations['decision_p87_sell_only'][0]}")
    print(f"decision_p87_frozen_watch_status={decision_expectations['decision_p87_frozen_watch'][0]}")
    print(f"decision_p87_insufficient_status={decision_expectations['decision_p87_insufficient'][0]}")
    print(f"frozen_or_insufficient_confirmations={frozen_confirmations}")
    forbidden_table_count = scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")
    auto_confirmation_count = scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")
    print(f"forbidden_broker_order_push_tables={forbidden_table_count}")
    print(f"auto_confirmation_rows={auto_confirmation_count}")
    print(f"artifact_dir={artifact_dir}")


if __name__ == "__main__":
    main()
