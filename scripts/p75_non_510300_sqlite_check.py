#!/usr/bin/env python3
import json
import sqlite3
import sys
from pathlib import Path


SYMBOL = "159915"
TRACKED_INDEX = "399006"


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


def row(db, sql, params=()):
    result = db.execute(sql, params).fetchone()
    return dict(result) if result else {}


def load_json_file(path):
    require(path.exists(), f"missing JSON file: {path}")
    return json.loads(path.read_text(encoding="utf-8"))


def request_matches(requests, path, *, query=None, form=None):
    query = query or {}
    form = form or {}
    for item in requests:
        if item.get("path") != path:
            continue
        item_query = item.get("query") or {}
        item_form = item.get("form") or {}
        if all(item_query.get(key) == value for key, value in query.items()) and all(item_form.get(key) == value for key, value in form.items()):
            return True
    return False


def cninfo_request_matches(requests):
    for item in requests:
        if item.get("path") != "/new/hisAnnouncement/query":
            continue
        form = item.get("form") or {}
        stock = form.get("stock", "")
        if (stock == SYMBOL or stock.startswith(f"{SYMBOL},")) and form.get("seDate") == "2026-06-01~2026-06-30":
            return True
    return False


def main():
    if len(sys.argv) != 4:
        fail("usage: p75_non_510300_sqlite_check.py <sqlite_path> <request_log_json> <artifact_dir>")

    sqlite_path = Path(sys.argv[1])
    request_log_path = Path(sys.argv[2])
    artifact_dir = Path(sys.argv[3])
    require(sqlite_path.exists(), "sqlite database does not exist", str(sqlite_path))
    artifact_dir.mkdir(parents=True, exist_ok=True)

    requests = load_json_file(request_log_path)
    require(request_matches(requests, "/market", query={"symbol": SYMBOL}), "market collector did not request 159915", requests)
    require(cninfo_request_matches(requests), "CNInfo collector did not request 159915 with the explicit date window", requests)
    require(
        request_matches(requests, "/api/disc/announcement/searchQuery", query={"keyword": SYMBOL}),
        "SZSE collector did not request 159915",
        requests,
    )

    db = sqlite3.connect(f"file:{sqlite_path}?mode=ro", uri=True)
    db.row_factory = sqlite3.Row

    position = row(
        db,
        "SELECT symbol,name,quantity,current_price,market_value,buy_reason,asset_tag FROM positions WHERE symbol=? ORDER BY updated_at DESC LIMIT 1",
        (SYMBOL,),
    )
    require(position, "missing UI-created 159915 position")
    require(position["name"] == "创业板ETF" and position["asset_tag"] == "satellite", "159915 position content mismatch", position)

    market = row(
        db,
        "SELECT market_snapshot_id,symbol,trade_date,close_price,market_metrics_json FROM market_snapshots WHERE symbol=? ORDER BY created_at DESC LIMIT 1",
        (SYMBOL,),
    )
    require(market, "missing 159915 market snapshot")
    metrics = json.loads(market["market_metrics_json"] or "{}")
    metric_text = json.dumps(metrics, ensure_ascii=False)
    for expected in ["p34_source_health", SYMBOL, TRACKED_INDEX, "valuation_percentiles", "request_id"]:
        require(expected in metric_text, f"market metrics missing {expected}", metrics)
    request_id = metrics.get("request_id")
    source_health = ((metrics.get("metadata") or {}).get("p34_source_health") or {})
    tracked_health = source_health.get("tracked_index") or {}
    valuation_health = source_health.get("valuation_percentiles") or {}
    require(tracked_health.get("request_id") == request_id, "tracked index health lost request correlation", tracked_health)
    require(valuation_health.get("affected_symbols") == [TRACKED_INDEX], "valuation health is not bound to tracked index", valuation_health)

    source_verification = row(
        db,
        "SELECT verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level "
        "FROM source_verifications WHERE symbol=? ORDER BY created_at DESC LIMIT 1",
        (SYMBOL,),
    )
    require(source_verification.get("verification_status") == "satisfied", "159915 formal evidence is not satisfied", source_verification)
    require(source_verification.get("high_grade_independent_source_count", 0) >= 2, "159915 formal evidence lacks two high-grade sources", source_verification)

    knowledge_counts = {
        "intelligence_items": scalar(
            db,
            "SELECT COUNT(DISTINCT i.intelligence_id) FROM intelligence_items i JOIN intelligence_summary s ON s.intelligence_id=i.intelligence_id WHERE s.symbol=?",
            (SYMBOL,),
        ),
        "intelligence_summary": scalar(db, "SELECT COUNT(*) FROM intelligence_summary WHERE symbol=?", (SYMBOL,)),
        "rag_chunks_indexed": scalar(
            db,
            "SELECT COUNT(*) FROM rag_chunks r JOIN intelligence_summary s ON s.summary_id=r.summary_id "
            "WHERE s.symbol=? AND r.index_status='indexed'",
            (SYMBOL,),
        ),
    }
    require(knowledge_counts["intelligence_items"] >= 2, "missing 159915 public evidence items", knowledge_counts)
    require(knowledge_counts["intelligence_summary"] >= 2, "missing 159915 public evidence summaries", knowledge_counts)
    require(knowledge_counts["rag_chunks_indexed"] >= 2, "159915 evidence was not indexed into VecLite", knowledge_counts)

    decision = row(
        db,
        "SELECT decision_id,request_id,symbol,workflow_status,final_verdict_status,confirmation_status,analyst_reports_json,context_snapshot_json "
        "FROM decision_records WHERE symbol=? AND question LIKE '%P75 非510300连续真实UI场景%' ORDER BY created_at DESC LIMIT 1",
        (SYMBOL,),
    )
    require(decision, "missing P75 159915 consultation decision")
    require(decision["workflow_status"] == "completed", "P75 159915 decision did not complete", decision)
    reports = json.loads(decision["analyst_reports_json"] or "[]")
    require(len(reports) >= 3, "P75 159915 decision lacks analyst reports", reports)
    require(all(item.get("parse_status") == "parsed" for item in reports), "analyst report parse_status mismatch", reports)
    require(all(item.get("quality_status") == "passed" for item in reports), "analyst report quality_status mismatch", reports)

    audit_counts = {
        "market_refresh": scalar(db, "SELECT COUNT(*) FROM audit_events WHERE action='refresh_market_data' AND input_ref='market-refresh'", ()),
        "public_evidence_command": scalar(
            db,
            "SELECT COUNT(*) FROM audit_events WHERE input_ref='public-evidence-refresh:symbol=159915:start=2026-06-01:end=2026-06-30'",
            (),
        ),
        "public_evidence_ingestion": scalar(
            db,
            "SELECT COUNT(*) FROM audit_events WHERE input_ref=? AND output_ref='source=public_evidence count=2'",
            (SYMBOL,),
        ),
        "consultation_decision": scalar(db, "SELECT COUNT(*) FROM audit_events WHERE request_id=? AND decision_id=?", (decision["request_id"], decision["decision_id"])),
    }
    require(audit_counts["market_refresh"] >= 1, "missing market refresh audit", audit_counts)
    require(audit_counts["public_evidence_command"] >= 1, "missing public evidence command audit", audit_counts)
    require(audit_counts["public_evidence_ingestion"] >= 1, "missing public evidence ingestion audit", audit_counts)
    require(audit_counts["consultation_decision"] >= 1, "missing consultation decision audit chain", audit_counts)

    browser_results = load_json_file(artifact_dir / "browser-results.json")
    require(browser_results.get("status") == "passed", "browser journey did not report passed", browser_results)
    require(browser_results.get("symbol") == SYMBOL and browser_results.get("tracked_index_symbol") == TRACKED_INDEX, "browser result symbol binding mismatch", browser_results)

    forbidden_tables = [
        dict(item)
        for item in db.execute(
            "SELECT name FROM sqlite_master WHERE type='table' AND ("
            "LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%trade_execution%' "
            "OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%webhook%')"
        ).fetchall()
    ]
    require(not forbidden_tables, "forbidden trading or external-push table exists", forbidden_tables)

    summary = {
        "status": "passed",
        "sqlite_path": str(sqlite_path),
        "request_log_path": str(request_log_path),
        "symbol": SYMBOL,
        "tracked_index_symbol": TRACKED_INDEX,
        "position": position,
        "market": {
            "market_snapshot_id": market["market_snapshot_id"],
            "trade_date": market["trade_date"],
            "close_price": market["close_price"],
            "request_id": request_id,
            "tracked_index_health": tracked_health,
            "valuation_health": valuation_health,
        },
        "source_verification": source_verification,
        "knowledge_counts": knowledge_counts,
        "decision": {
            "decision_id": decision["decision_id"],
            "request_id": decision["request_id"],
            "workflow_status": decision["workflow_status"],
            "final_verdict_status": decision["final_verdict_status"],
            "confirmation_status": decision["confirmation_status"],
            "analyst_report_count": len(reports),
        },
        "audit_counts": audit_counts,
        "forbidden_tables": forbidden_tables,
    }
    (artifact_dir / "non-510300-db-impact-summary.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(summary, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
