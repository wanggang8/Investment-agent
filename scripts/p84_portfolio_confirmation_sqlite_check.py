#!/usr/bin/env python3
"""Read-only SQLite checks for P84 portfolio and confirmation acceptance."""

from __future__ import annotations

import sqlite3
import sys
from pathlib import Path


def scalar(conn: sqlite3.Connection, sql: str, args: tuple = ()) -> float | int | str:
    row = conn.execute(sql, args).fetchone()
    if row is None:
        return 0
    return row[0] if row[0] is not None else 0


def main() -> None:
    if len(sys.argv) != 3:
        raise SystemExit("usage: p84_portfolio_confirmation_sqlite_check.py <sqlite> <artifact-dir>")
    db_path = Path(sys.argv[1])
    artifact_dir = sys.argv[2]
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    positions = conn.execute("SELECT symbol,quantity,current_price,market_value,asset_tag FROM positions ORDER BY symbol").fetchall()
    snapshot = conn.execute("SELECT cash,total_assets,cash_ratio,position_count FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1").fetchone()
    if snapshot is None:
        raise SystemExit("status=failed\nreason=missing_latest_snapshot")

    calculated_market_value = sum(float(row["quantity"]) * float(row["current_price"]) for row in positions)
    persisted_market_value = sum(float(row["market_value"]) for row in positions)
    expected_total_assets = float(snapshot["cash"]) + persisted_market_value
    if abs(calculated_market_value - persisted_market_value) > 0.01:
        raise SystemExit("status=failed\nreason=position_market_value_mismatch")
    if abs(expected_total_assets - float(snapshot["total_assets"])) > 0.01:
        raise SystemExit("status=failed\nreason=portfolio_total_assets_mismatch")

    confirm_count = int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE decision_id='decision_p84_pending' AND confirmation_type='executed_manually'"))
    tx_count = int(scalar(conn, "SELECT COUNT(*) FROM position_transactions WHERE confirmation_id IN (SELECT confirmation_id FROM operation_confirmations WHERE decision_id='decision_p84_pending')"))
    audit_count = int(scalar(conn, "SELECT COUNT(*) FROM audit_events WHERE decision_id='decision_p84_pending' OR action LIKE '%portfolio%' OR input_ref_type='portfolio'"))
    if confirm_count < 1 or tx_count < 1 or audit_count < 1:
        raise SystemExit("status=failed\nreason=missing_confirmation_transaction_or_audit")

    print("status=passed")
    print(f"position_count={len(positions)}")
    print(f"snapshot_position_count={int(snapshot['position_count'])}")
    print(f"snapshot_cash={float(snapshot['cash']):.2f}")
    print(f"snapshot_total_assets={float(snapshot['total_assets']):.2f}")
    print(f"calculated_market_value={calculated_market_value:.2f}")
    print(f"operation_confirmations_p84={confirm_count}")
    print(f"position_transactions_p84={tx_count}")
    print(f"local_account_import_batches={scalar(conn, 'SELECT COUNT(*) FROM local_account_import_batches')}")
    print(f"local_account_corrections={scalar(conn, 'SELECT COUNT(*) FROM local_account_corrections')}")
    print(f"portfolio_audit_events={audit_count}")
    print(f"review_confirmation_count={scalar(conn, 'SELECT COUNT(*) FROM operation_confirmations')}")
    decision_status = scalar(conn, "SELECT confirmation_status FROM decision_records WHERE decision_id='decision_p84_pending'")
    forbidden_table_count = scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")
    auto_confirmation_count = scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")
    print(f"decision_p84_status={decision_status}")
    print(f"forbidden_broker_order_push_tables={forbidden_table_count}")
    print(f"auto_confirmation_rows={auto_confirmation_count}")
    print(f"artifact_dir={artifact_dir}")


if __name__ == "__main__":
    main()
