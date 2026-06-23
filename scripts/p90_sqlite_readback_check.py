#!/usr/bin/env python3
"""Read-only SQLite checks for P90 capital-flow provider acceptance."""

from __future__ import annotations

import json
import math
import sqlite3
import sys
from pathlib import Path


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def close_enough(actual: object, expected: object) -> bool:
    try:
        return math.isclose(float(actual), float(expected), rel_tol=0, abs_tol=0.000001)
    except (TypeError, ValueError):
        return False


def market_metrics(conn: sqlite3.Connection, snapshot_id: str) -> dict:
    row = conn.execute(
        "SELECT market_snapshot_id,symbol,trade_date,market_metrics_json FROM market_snapshots WHERE market_snapshot_id=?",
        (snapshot_id,),
    ).fetchone()
    require(row is not None, f"missing_market_snapshot:{snapshot_id}")
    return {
        "market_snapshot_id": row["market_snapshot_id"],
        "symbol": row["symbol"],
        "trade_date": row["trade_date"],
        "market_metrics": json.loads(row["market_metrics_json"] or "{}"),
    }


def structured_capital_flow(snapshot: dict) -> dict:
    return (
        ((snapshot.get("market_metrics") or {}).get("metadata") or {})
        .get("p88_structured_fields", {})
        .get("capital_flow")
        or {}
    )


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: p90_sqlite_readback_check.py <sqlite> <browser-results.json> <source-preverification.json>")
    db_path = Path(sys.argv[1])
    browser = json.loads(Path(sys.argv[2]).read_text(encoding="utf-8"))
    source = json.loads(Path(sys.argv[3]).read_text(encoding="utf-8"))

    before_id = (browser.get("pre_refresh") or {}).get("market_snapshot_id")
    provider_id = (browser.get("provider_readback") or {}).get("market_snapshot_id")
    require(before_id == "market_p90_seed", f"pre_refresh_snapshot:{before_id}")
    require(provider_id and provider_id != before_id, "missing_runtime_provider_snapshot_id")

    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    before = market_metrics(conn, before_id)
    require(before["symbol"] == "600000", f"pre_refresh_symbol:{before['symbol']}")
    require(structured_capital_flow(before) == {}, "pre_refresh_capital_flow_should_be_absent")

    provider = market_metrics(conn, provider_id)
    require(provider["symbol"] == "600000", f"provider_symbol:{provider['symbol']}")
    stored = structured_capital_flow(provider)
    require(stored != {}, "sqlite_capital_flow_missing")

    api_flow = ((browser.get("provider_readback") or {}).get("capital_flow") or {})
    expected = ((source.get("category") or {}).get("fields") or {})
    for key in ("date", "net_inflow", "net_outflow", "raw_net_flow"):
        require(key in stored, f"sqlite_missing:{key}")
        require(key in api_flow, f"api_missing:{key}")
        require(key in expected, f"source_missing:{key}")
    require(stored["date"] == expected["date"], f"sqlite_date:{stored['date']}:{expected['date']}")
    require(api_flow["date"] == expected["date"], f"api_date:{api_flow['date']}:{expected['date']}")
    for key in ("net_inflow", "net_outflow", "raw_net_flow"):
        require(close_enough(stored[key], expected[key]), f"sqlite_{key}:{stored[key]}:{expected[key]}")
        require(close_enough(api_flow[key], expected[key]), f"api_{key}:{api_flow[key]}:{expected[key]}")
    require(close_enough(float(stored["net_inflow"]) + float(stored["net_outflow"]), abs(float(stored["raw_net_flow"]))), "directional_flow_mismatch")
    print("status=passed")
    print("pre_refresh_capital_flow_absent=passed")
    print(f"runtime_provider_snapshot_id={provider_id}")
    print("api_ui_capital_flow_readback=passed")
    print("sqlite_capital_flow_readback=passed")
    print("directional_net_flow_mapping=passed")


if __name__ == "__main__":
    main()
