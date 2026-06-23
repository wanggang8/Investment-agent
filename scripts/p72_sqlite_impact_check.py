#!/usr/bin/env python3
import json
import math
import sqlite3
import sys
from pathlib import Path


EXPECTED = {
    "symbol": "510300",
    "cash": 95630.5,
    "total_assets": 101265.0,
    "position_count": 2,
    "total_quantity": 1390.0,
    "total_market_value": 5634.5,
}


def approx(actual, expected, tolerance=0.02):
    return math.isclose(float(actual), float(expected), abs_tol=tolerance)


def scalar(db, sql, params=()):
    row = db.execute(sql, params).fetchone()
    return row[0] if row else None


def rows(db, sql, params=()):
    return [dict(row) for row in db.execute(sql, params).fetchall()]


def fail(message, details=None):
    payload = {"status": "failed", "message": message}
    if details is not None:
      payload["details"] = details
    print(json.dumps(payload, ensure_ascii=False, indent=2), file=sys.stderr)
    raise SystemExit(1)


def require(condition, message, details=None):
    if not condition:
        fail(message, details)


def main():
    if len(sys.argv) != 3:
        fail("usage: p72_sqlite_impact_check.py <sqlite_path> <artifact_dir>")

    sqlite_path = Path(sys.argv[1])
    artifact_dir = Path(sys.argv[2])
    require(sqlite_path.exists(), "sqlite database does not exist", str(sqlite_path))
    artifact_dir.mkdir(parents=True, exist_ok=True)

    db = sqlite3.connect(f"file:{sqlite_path}?mode=ro", uri=True)
    db.row_factory = sqlite3.Row

    latest_snapshot = dict(db.execute(
        "SELECT snapshot_id,cash,total_assets,cash_ratio,high_risk_ratio,position_count,source,created_at "
        "FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1"
    ).fetchone() or {})
    require(latest_snapshot, "missing latest portfolio snapshot")
    require(approx(latest_snapshot["cash"], EXPECTED["cash"]), "latest cash mismatch", latest_snapshot)
    require(approx(latest_snapshot["total_assets"], EXPECTED["total_assets"]), "latest total_assets mismatch", latest_snapshot)
    require(latest_snapshot["position_count"] == EXPECTED["position_count"], "latest position_count mismatch", latest_snapshot)

    aggregate = dict(db.execute(
        "SELECT COUNT(*) AS row_count, SUM(quantity) AS total_quantity, SUM(market_value) AS total_market_value "
        "FROM positions WHERE symbol=?",
        (EXPECTED["symbol"],),
    ).fetchone())
    require(aggregate["row_count"] == EXPECTED["position_count"], "position row count mismatch", aggregate)
    require(approx(aggregate["total_quantity"], EXPECTED["total_quantity"]), "position quantity mismatch", aggregate)
    require(approx(aggregate["total_market_value"], EXPECTED["total_market_value"]), "position market value mismatch", aggregate)
    position_rows = rows(
        db,
        "SELECT symbol,name,quantity,cost_price,current_price,market_value,buy_reason,asset_tag,position_state "
        "FROM positions WHERE symbol=? ORDER BY quantity DESC, asset_tag",
        (EXPECTED["symbol"],),
    )
    position_field_readback = {
        "rows": position_rows,
        "checks": {
            "symbol": len(position_rows) == 2 and all(row["symbol"] == EXPECTED["symbol"] for row in position_rows),
            "name": all(row["name"] == "沪深300ETF" for row in position_rows),
            "quantity": sorted(float(row["quantity"]) for row in position_rows) == [100.0, 1290.0] or sorted(float(row["quantity"]) for row in position_rows) == [100.0, 1300.0],
            "cost_price": any(approx(row["cost_price"], 3.2, 0.08) for row in position_rows) and any(approx(row["cost_price"], 3.9, 0.02) for row in position_rows),
            "buy_reason": any("人工复核后更新成本与现价" in (row["buy_reason"] or "") for row in position_rows) and any("批量导入追加线下确认持仓" in (row["buy_reason"] or "") for row in position_rows),
            "asset_tag": {"core", "satellite"}.issubset({row["asset_tag"] for row in position_rows}),
            "position_state": all(row["position_state"] == "normal" for row in position_rows),
        },
    }
    require(all(position_field_readback["checks"].values()), "position field readback mismatch", position_field_readback)

    portfolio_counts = {
        "portfolio_snapshots": scalar(db, "SELECT COUNT(*) FROM portfolio_snapshots"),
        "position_snapshots": scalar(db, "SELECT COUNT(*) FROM position_snapshots WHERE symbol=?", (EXPECTED["symbol"],)),
        "local_account_import_batches_committed": scalar(db, "SELECT COUNT(*) FROM local_account_import_batches WHERE status='committed'"),
        "local_account_corrections": scalar(db, "SELECT COUNT(*) FROM local_account_corrections"),
        "operation_confirmations_total": scalar(db, "SELECT COUNT(*) FROM operation_confirmations WHERE symbol=?", (EXPECTED["symbol"],)),
        "operation_confirmations_decision_linked": scalar(db, "SELECT COUNT(*) FROM operation_confirmations WHERE symbol=? AND decision_id<>''", (EXPECTED["symbol"],)),
        "operation_confirmations_offline": scalar(db, "SELECT COUNT(*) FROM operation_confirmations WHERE symbol=? AND decision_id=''", (EXPECTED["symbol"],)),
        "position_transactions_total": scalar(db, "SELECT COUNT(*) FROM position_transactions WHERE symbol=?", (EXPECTED["symbol"],)),
        "position_transactions_decision_linked": scalar(
            db,
            "SELECT COUNT(*) FROM position_transactions tx JOIN operation_confirmations c ON c.confirmation_id=tx.confirmation_id "
            "WHERE tx.symbol=? AND c.decision_id<>''",
            (EXPECTED["symbol"],),
        ),
        "audit_events_user_confirm": scalar(db, "SELECT COUNT(*) FROM audit_events WHERE action='confirm_operation'"),
    }
    require(portfolio_counts["portfolio_snapshots"] >= 5, "expected multiple portfolio snapshots", portfolio_counts)
    require(portfolio_counts["position_snapshots"] >= 5, "expected position snapshots for scenario symbol", portfolio_counts)
    require(portfolio_counts["local_account_import_batches_committed"] >= 1, "missing committed import batch", portfolio_counts)
    require(portfolio_counts["local_account_corrections"] >= 1, "missing correction audit record", portfolio_counts)
    require(portfolio_counts["operation_confirmations_total"] >= 2, "missing operation confirmations", portfolio_counts)
    require(portfolio_counts["operation_confirmations_decision_linked"] >= 1, "missing decision-linked confirmation", portfolio_counts)
    require(portfolio_counts["operation_confirmations_offline"] >= 1, "missing offline local transaction confirmation", portfolio_counts)
    require(portfolio_counts["position_transactions_total"] >= 2, "missing position transactions", portfolio_counts)
    require(portfolio_counts["position_transactions_decision_linked"] >= 1, "missing decision-linked transaction", portfolio_counts)
    operation_confirmation_readback = {
        "rows": rows(
            db,
            "SELECT confirmation_id,decision_id,confirmation_type,operation_type,symbol,quantity,price,executed_at,note "
            "FROM operation_confirmations WHERE symbol=? ORDER BY created_at",
            (EXPECTED["symbol"],),
        ),
    }
    operation_confirmation_readback["checks"] = {
        "decision_linked_manual_execution": any(row["decision_id"] and row["confirmation_type"] == "executed_manually" for row in operation_confirmation_readback["rows"]),
        "offline_local_transaction": any(not row["decision_id"] and row["operation_type"] in {"buy", "sell", "reduce"} for row in operation_confirmation_readback["rows"]),
        "quantity_and_price": all(float(row["quantity"] or 0) > 0 and float(row["price"] or 0) > 0 for row in operation_confirmation_readback["rows"]),
    }
    require(all(operation_confirmation_readback["checks"].values()), "operation confirmation field readback mismatch", operation_confirmation_readback)
    transaction_readback = {
        "rows": rows(
            db,
            "SELECT transaction_id,confirmation_id,symbol,operation_type,quantity,price,occurred_at,before_position_json,after_position_json "
            "FROM position_transactions WHERE symbol=? ORDER BY created_at",
            (EXPECTED["symbol"],),
        ),
    }
    transaction_readback["checks"] = {
        "transaction_rows": len(transaction_readback["rows"]) >= 2,
        "symbol": all(row["symbol"] == EXPECTED["symbol"] for row in transaction_readback["rows"]),
        "quantity_and_price": all(float(row["quantity"]) > 0 and float(row["price"]) > 0 for row in transaction_readback["rows"]),
        "before_after_state": any(row["before_position_json"] and row["after_position_json"] for row in transaction_readback["rows"]),
    }
    require(all(transaction_readback["checks"].values()), "position transaction field readback mismatch", transaction_readback)

    decision = dict(db.execute(
        "SELECT decision_id,workflow_status,final_verdict_status,confirmation_status,analyst_reports_json,context_snapshot_json "
        "FROM decision_records WHERE question LIKE '%P72 真实用户场景%' ORDER BY created_at DESC LIMIT 1"
    ).fetchone() or {})
    require(decision, "missing P72 consultation decision")
    require(decision["workflow_status"] == "completed", "P72 decision did not complete", decision)
    require(decision["confirmation_status"] == "executed_manually", "P72 decision confirmation not recorded", decision)
    reports = json.loads(decision["analyst_reports_json"] or "[]")
    require(len(reports) > 0, "missing analyst reports", decision)
    require(all(item.get("parse_status") == "parsed" for item in reports), "analyst report parse_status mismatch", reports)
    require(all(item.get("quality_status") == "passed" for item in reports), "analyst report quality_status mismatch", reports)
    table_readback = {
        "portfolio_snapshots_have_source": scalar(db, "SELECT COUNT(*) FROM portfolio_snapshots WHERE source IN ('manual','system')") >= portfolio_counts["portfolio_snapshots"],
        "position_snapshots_have_required_fields": scalar(db, "SELECT COUNT(*) FROM position_snapshots WHERE symbol=? AND quantity>0 AND cost_price>0 AND market_value>0", (EXPECTED["symbol"],)) >= 5,
        "decision_record_has_snapshot_refs": bool(decision.get("decision_id")),
        "evidence_refs_for_decision": scalar(db, "SELECT COUNT(*) FROM evidence_refs WHERE decision_id=?", (decision["decision_id"],)) >= 1,
        "audit_events_include_user_confirm": portfolio_counts["audit_events_user_confirm"] >= 2,
    }
    require(all(table_readback.values()), "table readback mismatch", table_readback)

    knowledge_counts = {
        "intelligence_items_p72": scalar(db, "SELECT COUNT(*) FROM intelligence_items WHERE raw_title LIKE '%P72 510300%'"),
        "intelligence_summary_p72": scalar(db, "SELECT COUNT(*) FROM intelligence_summary WHERE symbol=? AND summary LIKE '%P72 真实用户验收%'", (EXPECTED["symbol"],)),
        "rag_chunks_p72_indexed": scalar(
            db,
            "SELECT COUNT(*) FROM rag_chunks r JOIN intelligence_summary s ON s.summary_id=r.summary_id "
            "WHERE s.symbol=? AND r.index_status='indexed' AND r.chunk_text LIKE '%P72 真实用户验收%'",
            (EXPECTED["symbol"],),
        ),
        "source_verifications": scalar(db, "SELECT COUNT(*) FROM source_verifications"),
    }
    require(knowledge_counts["intelligence_items_p72"] >= 1, "missing P72 intelligence item", knowledge_counts)
    require(knowledge_counts["intelligence_summary_p72"] >= 1, "missing P72 intelligence summary", knowledge_counts)
    require(knowledge_counts["rag_chunks_p72_indexed"] >= 1, "missing indexed P72 RAG chunk", knowledge_counts)
    require(knowledge_counts["source_verifications"] >= 1, "missing source verification facts", knowledge_counts)

    daily_counts = {
        "manual_daily_reports": scalar(db, "SELECT COUNT(*) FROM daily_discipline_reports WHERE source_type='manual'"),
        "manual_daily_reports_with_decision": scalar(db, "SELECT COUNT(*) FROM daily_discipline_reports WHERE source_type='manual' AND COALESCE(decision_id,'')<>''"),
        "risk_alerts": scalar(db, "SELECT COUNT(*) FROM risk_alerts"),
        "risk_alerts_510300": scalar(db, "SELECT COUNT(*) FROM risk_alerts WHERE symbol=?", (EXPECTED["symbol"],)),
        "notifications": scalar(db, "SELECT COUNT(*) FROM notifications"),
    }
    require(daily_counts["manual_daily_reports"] >= 1, "missing manual daily discipline report", daily_counts)
    require(daily_counts["manual_daily_reports_with_decision"] >= 1, "missing daily report decision link", daily_counts)
    require(daily_counts["risk_alerts"] >= 1, "missing risk alert readback facts", daily_counts)
    require(daily_counts["notifications"] >= 1, "missing notification facts", daily_counts)

    forbidden_tables = rows(
        db,
        "SELECT name FROM sqlite_master WHERE type='table' AND ("
        "LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%trade_execution%' "
        "OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%webhook%')"
    )
    require(not forbidden_tables, "forbidden trading or external-push table exists", forbidden_tables)

    summary = {
        "status": "passed",
        "sqlite_path": str(sqlite_path),
        "expected": EXPECTED,
        "latest_snapshot": latest_snapshot,
        "position_aggregate": aggregate,
        "position_field_readback": position_field_readback,
        "portfolio_counts": portfolio_counts,
        "operation_confirmation_readback": operation_confirmation_readback,
        "transaction_readback": transaction_readback,
        "decision": {
            "decision_id": decision["decision_id"],
            "workflow_status": decision["workflow_status"],
            "final_verdict_status": decision["final_verdict_status"],
            "confirmation_status": decision["confirmation_status"],
            "analyst_report_count": len(reports),
        },
        "table_readback": table_readback,
        "knowledge_counts": knowledge_counts,
        "daily_counts": daily_counts,
        "forbidden_tables": forbidden_tables,
    }
    (artifact_dir / "db-impact-summary.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
