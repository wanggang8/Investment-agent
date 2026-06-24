#!/usr/bin/env python3
"""P104 local product operation/linkage acceptance runner."""

from __future__ import annotations

import argparse
import json
import sqlite3
import time
import urllib.error
import urllib.parse
import urllib.request
from pathlib import Path
from typing import Any


NOW = "2026-06-24T03:00:00Z"


class AcceptanceFailure(RuntimeError):
    pass


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AcceptanceFailure(message)


def request_json(base_url: str, method: str, path: str, body: dict[str, Any] | None = None, request_id: str = "req_p104") -> dict[str, Any]:
    data = None
    headers = {
        "Accept": "application/json",
        "X-Request-ID": request_id,
    }
    if body is not None:
        data = json.dumps(body, ensure_ascii=False).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(base_url.rstrip("/") + path, data=data, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=20) as resp:
        payload = resp.read().decode("utf-8")
        return json.loads(payload) if payload else {}


def expect_http_status(base_url: str, method: str, path: str, status: int, body: dict[str, Any] | None = None) -> dict[str, Any]:
    try:
        out = request_json(base_url, method, path, body, "req_p104_expected_reject")
    except urllib.error.HTTPError as err:
        payload = err.read().decode("utf-8")
        if err.code != status:
            raise AcceptanceFailure(f"expected HTTP {status} for {path}, got {err.code}: {payload}") from err
        return json.loads(payload) if payload else {}
    raise AcceptanceFailure(f"expected HTTP {status} for {path}, got success: {out}")


def data(envelope: dict[str, Any]) -> Any:
    return envelope.get("data")


def seed_supporting_records(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
    try:
        conn.execute(
            """
            INSERT OR REPLACE INTO decision_records (
              decision_id,request_id,workflow_type,symbol,question,workflow_status,record_type,dashboard_state,
              capability_status,capability_reason,source_verification_status,risk_reason_code,
              media_heat_summary_json,user_emotion_tags_json,triggered_rules_json,errors_json,
              final_verdict_status,final_verdict_text,prohibited_actions_json,optional_actions_json,
              confirmation_status,portfolio_snapshot_id,market_snapshot_id,rule_version,
              analyst_reports_json,expected_return_scenarios_json,arbitration_chain_json,context_snapshot_json,created_at
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
            """,
            (
                "decision_p104_execute",
                "req_p104_seed_decision",
                "consultation",
                "510300",
                "P104 operation linkage acceptance decision",
                "completed",
                "formal_trade_advice",
                "normal",
                "in_scope",
                "P104 seeded formal decision for manual confirmation linkage",
                "satisfied",
                "",
                '{"heat":"neutral"}',
                '["calm"]',
                '["manual_confirmation_required","portfolio_rebalance_review"]',
                "[]",
                "hold",
                "P104 本地验收决策：只允许人工记录线下处理结果",
                '["自动交易","外部推送"]',
                '["记录线下处理","继续观察"]',
                "pending",
                None,
                None,
                "v_p104",
                '[{"agent_name":"P104LocalAnalyst","conclusion":"用于产品联动验收的本地分析材料","key_reasons":["确认链路必须由用户触发","结果要能被闭环/复盘/审计读取"],"risk_warnings":["不会自动交易"],"confidence":"medium","evidence_ids":[]}]',
                '{"precision_status":"available","reason":"P104 deterministic acceptance seed","sample_count":20,"sample_window":"2024-2026","screening_condition":"local seeded context","probability_basis":"historical local sample","scenarios":[{"name":"base","return_rate":0.03,"return_range":"0%~6%","probability":0.5,"trigger":"valuation stable"}],"disclaimer":"仅用于本地人工决策辅助，不承诺收益"}',
                '[{"step":"rule","result":"hold"},{"step":"user","result":"pending"}]',
                '{"p104":"operation_linkage_acceptance"}',
                NOW,
            ),
        )
        conn.execute(
            """
            INSERT OR REPLACE INTO notifications (
              notification_id,type,severity,title,message,source_type,source_id,created_at
            ) VALUES (?,?,?,?,?,?,?,?)
            """,
            (
                "notif_p104_unread",
                "data_source_failure",
                "warning",
                "P104 数据源提示",
                "用于本地通知读写联动验收",
                "p104",
                "decision_p104_execute",
                NOW,
            ),
        )
        conn.execute(
            """
            INSERT OR REPLACE INTO risk_alerts (
              alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,
              prohibited_actions_json,suggested_actions_json,related_decision_id,related_notification_id,
              last_triggered_at,created_at,updated_at
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
            """,
            (
                "risk_p104_active",
                "data_degraded",
                "warning",
                "active",
                "510300",
                "P104 本地风险处置验收",
                '{"source":"p104"}',
                '["新增买入"]',
                '["人工复核数据质量"]',
                "decision_p104_execute",
                "notif_p104_unread",
                NOW,
                NOW,
                NOW,
            ),
        )
        conn.execute(
            """
            INSERT OR REPLACE INTO market_snapshots (
              market_snapshot_id,symbol,trade_date,close_price,pe_percentile,pb_percentile,market_metrics_json,created_at
            ) VALUES (?,?,?,?,?,?,?,?)
            """,
            (
                "market_p104_block",
                "000300",
                "2026-06-05",
                4000.0,
                42.0,
                37.0,
                '{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"}},"p34_data_categories":["index_valuation_files"]}}',
                NOW,
            ),
        )
        conn.commit()
    finally:
        conn.close()


def scalar(conn: sqlite3.Connection, sql: str, args: tuple[Any, ...] = ()) -> Any:
    row = conn.execute(sql, args).fetchone()
    if row is None:
        return 0
    return row[0] if row[0] is not None else 0


def exercise_operations(base_url: str, db_path: Path) -> dict[str, Any]:
    steps: list[dict[str, Any]] = []

    expect_http_status(
        base_url,
        "POST",
        "/api/v1/portfolio/adjustments",
        400,
        {"cash": 10, "total_assets": 20, "adjust_reason": "P104 rejects inconsistent total", "positions": []},
    )
    steps.append({"step": "portfolio_invalid_total_rejected", "status": "passed"})

    adjust = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/adjustments",
            {
                "cash": 1000,
                "total_assets": 2000,
                "adjust_reason": "P104 本地组合校准",
                "positions": [
                    {
                        "symbol": "510300",
                        "name": "沪深300ETF",
                        "quantity": 100,
                        "cost_price": 3.2,
                        "current_price": 4.0,
                        "buy_date": "2026-01-05",
                        "position_state": "sell_only",
                        "buy_reason": "买入逻辑破坏后只卖不买",
                        "asset_tag": "core",
                    },
                    {
                        "symbol": "159915",
                        "name": "创业板ETF",
                        "quantity": 100,
                        "cost_price": 2.0,
                        "current_price": 4.0,
                        "buy_date": "2026-01-06",
                        "position_state": "frozen_watch",
                        "buy_reason": "多源验证不足冻结观察",
                        "asset_tag": "satellite",
                    },
                    {
                        "symbol": "588000",
                        "name": "科创50ETF",
                        "quantity": 100,
                        "cost_price": 1.6,
                        "current_price": 2.0,
                        "buy_date": "2026-02-01",
                        "position_state": "normal",
                        "buy_reason": "小比例观察仓",
                        "asset_tag": "satellite",
                    },
                ],
            },
            "req_p104_adjust",
        )
    )
    require(adjust["position_count"] == 3, "portfolio adjustment should create three positions")
    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p104_current"))
    require(len(current["positions"]) == 3, "portfolio current should read adjusted positions")
    pos_510300 = next(item for item in current["positions"] if item["symbol"] == "510300")
    pos_588000 = next(item for item in current["positions"] if item["symbol"] == "588000")
    steps.append({"step": "portfolio_adjustment_readback", "snapshot_id": current["snapshot"]["snapshot_id"], "status": "passed"})

    edit = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/holdings",
            {
                "position_id": pos_510300["position_id"],
                "reason": "P104 持仓维护校准",
                "confirmation": "confirmed",
                "position": {
                    "symbol": "510300",
                    "name": "沪深300ETF",
                    "quantity": 90,
                    "cost_price": 3.2,
                    "current_price": 4.2,
                    "buy_date": "2026-01-05",
                    "position_state": "sell_only",
                    "buy_reason": "买入逻辑破坏后只卖不买",
                    "asset_tag": "core",
                },
            },
            "req_p104_edit_holding",
        )
    )
    require(edit.get("snapshot_id"), "holding edit should return a new snapshot")
    remove = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/holdings/remove",
            {"position_id": pos_588000["position_id"], "reason": "P104 移除观察仓", "confirmation": "confirmed"},
            "req_p104_remove_holding",
        )
    )
    require(remove.get("snapshot_id"), "holding remove should return a new snapshot")
    steps.append({"step": "holding_edit_remove", "status": "passed"})

    offline = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/offline-transactions",
            {
                "operation_type": "buy",
                "symbol": "159915",
                "name": "创业板ETF",
                "quantity": 5,
                "price": 4.1,
                "fees": 1,
                "executed_at": "2026-06-23T03:00:00Z",
                "buy_reason": "P104 线下人工买入记录",
                "note": "只记录线下动作，不连接券商",
            },
            "req_p104_offline_tx",
        )
    )
    require(offline.get("transaction_id"), "offline transaction should return transaction id")
    steps.append({"step": "offline_transaction", "status": "passed"})

    validate = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/imports/validate",
            {
                "rows": [
                    {
                        "row_number": 1,
                        "row_type": "transaction",
                        "operation_type": "buy",
                        "symbol": "512000",
                        "name": "券商ETF",
                        "quantity": 10,
                        "price": 1.2,
                        "fees": 0.5,
                        "occurred_at": "2026-06-22T03:00:00Z",
                        "buy_reason": "P104 批量导入交易",
                    }
                ]
            },
            "req_p104_import_validate",
        )
    )
    require(validate["summary"]["valid_count"] == 1 and validate["summary"]["invalid_count"] == 0, "import validate should accept one transaction")
    confirm_import = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/imports/confirm",
            {
                "import_batch_id": validate["import_batch_id"],
                "confirm_reason": "P104 确认批量导入",
                "rows": [
                    {
                        "row_number": 1,
                        "row_type": "transaction",
                        "operation_type": "buy",
                        "symbol": "512000",
                        "name": "券商ETF",
                        "quantity": 10,
                        "price": 1.2,
                        "fees": 0.5,
                        "occurred_at": "2026-06-22T03:00:00Z",
                        "buy_reason": "P104 批量导入交易",
                    }
                ],
            },
            "req_p104_import_confirm",
        )
    )
    require(confirm_import.get("transaction_id") or confirm_import.get("snapshot_id"), "import confirm should write facts")
    steps.append({"step": "batch_import_validate_confirm", "status": "passed"})

    current_after_import = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p104_current_after_import"))
    pos_for_correction = next(item for item in current_after_import["positions"] if item["symbol"] == "510300")
    correction = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/corrections",
            {
                "target_type": "position",
                "target_id": pos_for_correction["position_id"],
                "before_json": '{"quantity":100}',
                "after_json": '{"quantity":90}',
                "correction_reason": "P104 本地数量修正审计",
            },
            "req_p104_correction",
        )
    )
    require(correction.get("correction_id"), "correction should return correction id")
    rebalance = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/portfolio/rebalance-review",
            {
                "target_core_ratio": 0.6,
                "target_satellite_ratio": 0.3,
                "target_cash_ratio": 0.1,
                "drift_threshold": 0.05,
                "review_date": "2026-06-24",
            },
            "req_p104_rebalance",
        )
    )
    require(rebalance.get("review_id") and len(rebalance.get("items", [])) == 3, "rebalance should return three allocation buckets")
    steps.append({"step": "correction_and_rebalance", "status": "passed"})

    decision_before = data(request_json(base_url, "GET", "/api/v1/decisions/decision_p104_execute", request_id="req_p104_decision_before"))
    require(decision_before["user_confirmation"]["confirmation_status"] == "pending", "seeded decision should be pending")
    confirmation = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/decisions/decision_p104_execute/confirmations",
            {
                "confirmation_type": "executed_manually",
                "operation_type": "sell",
                "symbol": "510300",
                "quantity": 5,
                "price": 4.25,
                "fees": 1,
                "executed_at": "2026-06-23T04:00:00Z",
                "note": "P104 人工线下卖出记录",
            },
            "req_p104_confirmation_execute",
        )
    )
    require(confirmation["confirmation_status"] == "executed_manually", "confirmation should update status")
    decision_after = data(request_json(base_url, "GET", "/api/v1/decisions/decision_p104_execute", request_id="req_p104_decision_after"))
    require(decision_after["user_confirmation"]["confirmation_status"] == "executed_manually", "decision detail should read executed status")
    loop = data(request_json(base_url, "GET", "/api/v1/decision-loops/decision_p104_execute", request_id="req_p104_loop"))
    require(loop["decision_id"] == "decision_p104_execute", "decision loop detail should focus seeded decision")
    loop_list = data(request_json(base_url, "GET", "/api/v1/decision-loops?symbol=510300&limit=10", request_id="req_p104_loop_list"))
    loop_ids = [item["decision_id"] for item in loop_list.get("items", [])]
    require("decision_p104_execute" in loop_ids, "decision loop symbol-filtered list should include seeded decision")
    steps.append({"step": "decision_confirmation_loop", "status": "passed"})

    notification_list = data(request_json(base_url, "GET", "/api/v1/notifications", request_id="req_p104_notifications"))
    require(notification_list["unread_count"] >= 1, "notifications should include unread seed")
    request_json(base_url, "POST", "/api/v1/notifications/notif_p104_unread/read", request_id="req_p104_notification_read")
    request_json(base_url, "POST", "/api/v1/notifications/read-all", request_id="req_p104_notifications_read_all")
    steps.append({"step": "notifications_mark_read", "status": "passed"})

    risk = data(request_json(base_url, "GET", "/api/v1/risk-alerts/risk_p104_active", request_id="req_p104_risk_get"))
    require(risk["sop_status"] == "active", "risk alert should start active")
    risk_resolved = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/risk-alerts/risk_p104_active/lifecycle",
            {"status": "resolved", "reason": "P104 人工复核完成"},
            "req_p104_risk_resolve",
        )
    )
    require(risk_resolved["sop_status"] == "resolved", "risk alert should resolve")
    steps.append({"step": "risk_lifecycle", "status": "passed"})

    gate = data(request_json(base_url, "GET", "/api/v1/data-source-quality/gate-resolution?symbol=000300", request_id="req_p104_dq_gate"))
    require(gate["release_claim_state"] == "requires_resolution", "seeded current data policy should require resolution")
    resolution = data(
        request_json(
            base_url,
            "POST",
            "/api/v1/data-source-quality/resolutions",
            {
                "symbol": "000300",
                "resolution_type": "scope_exclusion",
                "scope": "P104 local-source operation linkage acceptance excludes current-data clean claim",
                "reason": "Seeded stale A-level source health for gate-resolution operation validation",
                "release_impact": "Do not claim current data clean in P104",
                "evidence_ref": "docs/release/acceptance/2026-06-24-p104-product-operation-linkage-matrix.md",
            },
            "req_p104_dq_resolution_create",
        )
    )
    active = resolution.get("active_resolution") or {}
    require(resolution["release_claim_state"] == "resolved_with_scope_exclusion" and active.get("resolution_id"), "data-quality scope exclusion should become active")
    retired = data(
        request_json(
            base_url,
            "POST",
            f"/api/v1/data-source-quality/resolutions/{urllib.parse.quote(active['resolution_id'])}/retire",
            request_id="req_p104_dq_resolution_retire",
        )
    )
    require(retired["release_claim_state"] == "requires_resolution", "retired data-quality resolution should restore requires_resolution")
    steps.append({"step": "data_quality_resolution_create_retire", "status": "passed"})

    dashboard = data(request_json(base_url, "GET", "/api/v1/dashboard/today", request_id="req_p104_dashboard"))
    require(dashboard["portfolio_summary"]["position_count"] >= 1, "dashboard should read latest portfolio summary")
    review = data(request_json(base_url, "GET", "/api/v1/review/summary", request_id="req_p104_review"))
    require(review is not None, "review summary should return data")
    audit = data(request_json(base_url, "GET", "/api/v1/audit-events", request_id="req_p104_audit"))
    require(audit.get("total", 0) >= 1, "audit list should include operation events")
    steps.append({"step": "downstream_dashboard_review_audit", "status": "passed"})

    return {"steps": steps}


def sqlite_readback(db_path: Path) -> dict[str, Any]:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row
    try:
        latest = conn.execute(
            "SELECT cash,total_assets,cash_ratio,high_risk_ratio,position_count FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1"
        ).fetchone()
        require(latest is not None, "missing latest portfolio snapshot")
        position_market_value = float(scalar(conn, "SELECT COALESCE(SUM(market_value),0) FROM positions"))
        expected_total = float(latest["cash"]) + position_market_value
        require(abs(expected_total - float(latest["total_assets"])) < 0.05, "latest snapshot total must equal cash plus positions")

        checks = {
            "portfolio_snapshots": int(scalar(conn, "SELECT COUNT(*) FROM portfolio_snapshots")),
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "p104_executed_confirmations": int(
                scalar(
                    conn,
                    "SELECT COUNT(*) FROM operation_confirmations WHERE decision_id='decision_p104_execute' AND confirmation_type='executed_manually'",
                )
            ),
            "import_batches_committed": int(scalar(conn, "SELECT COUNT(*) FROM local_account_import_batches WHERE status='committed'")),
            "corrections": int(scalar(conn, "SELECT COUNT(*) FROM local_account_corrections")),
            "risk_resolved": int(scalar(conn, "SELECT COUNT(*) FROM risk_alerts WHERE alert_id='risk_p104_active' AND sop_status='resolved'")),
            "p104_notification_read": int(scalar(conn, "SELECT COUNT(*) FROM notifications WHERE notification_id='notif_p104_unread' AND read_at IS NOT NULL")),
            "notifications_unread_total": int(scalar(conn, "SELECT COUNT(*) FROM notifications WHERE read_at IS NULL")),
            "dq_resolutions_retired": int(
                scalar(
                    conn,
                    "SELECT COUNT(*) FROM data_quality_gate_resolutions WHERE symbol='000300' AND status='retired'",
                )
            ),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
            "auto_confirmation_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")),
            "forbidden_broker_order_push_tables": int(
                scalar(
                    conn,
                    "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')",
                )
            ),
            "auto_rule_apply_audit_events": int(
                scalar(
                    conn,
                    "SELECT COUNT(*) FROM audit_events WHERE LOWER(action) LIKE '%auto%' AND (LOWER(action) LIKE '%rule%' OR LOWER(action) LIKE '%confirm%' OR LOWER(action) LIKE '%trade%')",
                )
            ),
        }
        require(checks["p104_executed_confirmations"] == 1, "P104 confirmation should be exactly one executed row")
        require(checks["position_transactions"] >= 3, "portfolio operations should create transaction rows")
        require(checks["import_batches_committed"] >= 1, "batch import should be committed")
        require(checks["corrections"] >= 1, "correction audit should be present")
        require(checks["risk_resolved"] == 1, "risk alert should be resolved")
        require(checks["p104_notification_read"] == 1, "P104 notification read operation should set read_at")
        require(checks["dq_resolutions_retired"] >= 1, "data-quality resolution should be retired")
        require(checks["audit_events"] >= 10, "audit trace should include broad operation evidence")
        require(checks["auto_confirmation_rows"] == 0, "automatic confirmation rows must be absent")
        require(checks["forbidden_broker_order_push_tables"] == 0, "forbidden broker/order/push tables must be absent")
        require(checks["auto_rule_apply_audit_events"] == 0, "automatic rule/trade/confirmation audit events must be absent")
        checks["latest_total_assets"] = round(float(latest["total_assets"]), 2)
        checks["latest_position_market_value"] = round(position_market_value, 2)
        checks["latest_cash"] = round(float(latest["cash"]), 2)
        return checks
    finally:
        conn.close()


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", required=True)
    parser.add_argument("--sqlite", required=True)
    parser.add_argument("--artifact-dir", required=True)
    args = parser.parse_args()

    db_path = Path(args.sqlite)
    artifact_dir = Path(args.artifact_dir)
    artifact_dir.mkdir(parents=True, exist_ok=True)

    seed_supporting_records(db_path)
    operations = exercise_operations(args.base_url, db_path)
    db = sqlite_readback(db_path)

    payload = {
        "status": "passed",
        "change": "p104-full-product-operation-linkage-acceptance",
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "base_url": args.base_url,
        "sqlite_path": str(db_path),
        "operation_steps": operations["steps"],
        "sqlite_readback": db,
        "safety": {
            "forbidden_broker_order_push_tables": db["forbidden_broker_order_push_tables"],
            "auto_confirmation_rows": db["auto_confirmation_rows"],
            "auto_rule_apply_audit_events": db["auto_rule_apply_audit_events"],
            "claim_boundary": "P104 validates local-source product operation linkage only. It does not claim Docker/install/package/physical second-machine validation, broker integration, automatic trading, automatic confirmation, automatic rule application, or return guarantees.",
        },
    }
    out = artifact_dir / "p104-operation-linkage-summary.json"
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
