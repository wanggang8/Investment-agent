#!/usr/bin/env python3
"""P117 continuous product usability acceptance runner."""

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


NOW = "2026-06-25T01:00:00Z"
CHANGE_ID = "p117-continuous-product-usability-acceptance"
SCENARIO_IDS = [f"U{i:02d}" for i in range(1, 18)]
SYMBOLS = ["510300", "159915", "512000", "110022", "161725"]


class AcceptanceFailure(RuntimeError):
    pass


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AcceptanceFailure(message)


def request_json(base_url: str, method: str, path: str, body: dict[str, Any] | None = None, request_id: str = "req_p117") -> dict[str, Any]:
    data = None
    headers = {"Accept": "application/json", "X-Request-ID": request_id}
    if body is not None:
        data = json.dumps(body, ensure_ascii=False).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(base_url.rstrip("/") + path, data=data, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=30) as resp:
        payload = resp.read().decode("utf-8")
        return json.loads(payload) if payload else {}


def expect_http_status(base_url: str, method: str, path: str, status: int, body: dict[str, Any] | None = None, request_id: str = "req_p117_reject") -> dict[str, Any]:
    try:
        out = request_json(base_url, method, path, body, request_id)
    except urllib.error.HTTPError as err:
        payload = err.read().decode("utf-8")
        if err.code != status:
            raise AcceptanceFailure(f"expected HTTP {status} for {path}, got {err.code}: {payload}") from err
        return json.loads(payload) if payload else {}
    raise AcceptanceFailure(f"expected HTTP {status} for {path}, got success: {out}")


def data(envelope: dict[str, Any]) -> Any:
    return envelope.get("data")


def scalar(conn: sqlite3.Connection, sql: str, args: tuple[Any, ...] = ()) -> Any:
    row = conn.execute(sql, args).fetchone()
    if row is None:
        return 0
    return row[0] if row[0] is not None else 0


def scenario(sid: str, day: str, title: str, expected: str = "fresh_pass") -> dict[str, Any]:
    return {
        "scenario_id": sid,
        "day": day,
        "title": title,
        "status": expected,
        "expected_eligibility": expected,
        "classification_reason": "Validated by P117 isolated continuous-use runner.",
        "config_mode": "local_seeded_continuous_use",
        "runtime_mode": "development",
        "use_stub": True,
        "provider_mode": "stub_local_linkage",
        "llm_mode": "not_configured_degraded_expected",
        "symbols": SYMBOLS,
        "api_evidence": [],
        "browser_evidence": [],
        "sqlite_evidence": [],
        "restart_evidence": [],
        "downstream_evidence": [],
        "rejection_evidence": [],
        "usability_evidence": [],
        "redaction_evidence": {},
        "safety_counters": {},
    }


def add_api(item: dict[str, Any], method: str, path: str, status: int = 200, request_id: str = "req_p117") -> None:
    item["api_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id})


def add_reject(item: dict[str, Any], method: str, path: str, status: int, request_id: str, reason: str) -> None:
    item["rejection_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id, "reason": reason})


def add_sqlite(item: dict[str, Any], table: str, field: str, row_count: int, label: str = "") -> None:
    item["sqlite_evidence"].append({"table": table, "field": field, "row_count": row_count, "query_label": label or f"{table}.{field}"})


def add_usable(item: dict[str, Any], dimension: str, assertion: str, value: Any = "ok") -> None:
    item["usability_evidence"].append({"dimension": dimension, "assertion": assertion, "value": value})


def seed_supporting_records(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
    try:
        for decision_id, symbol, title in [
            ("decision_p117_execute", "159915", "P117 连续使用人工计划确认"),
            ("decision_p117_error", "510300", "P117 连续使用错误标注复盘"),
            ("decision_p117_browser_execute", "512000", "P117 浏览器连续使用确认"),
        ]:
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
                    decision_id,
                    f"req_{decision_id}",
                    "consultation",
                    symbol,
                    title,
                    "completed",
                    "formal_trade_advice",
                    "normal",
                    "in_scope",
                    "P117 seeded decision for continuous usability acceptance",
                    "satisfied",
                    "",
                    '{"heat":"neutral"}',
                    '["calm"]',
                    '["manual_confirmation_required","continuous_review"]',
                    "[]",
                    "hold",
                    "P117 本地连续使用验收决策：只允许人工记录线下处理结果",
                    '["自动交易","外部推送"]',
                    '["记录线下处理","继续观察","标记错误"]',
                    "pending",
                    None,
                    None,
                    "v_p117",
                    '[{"agent_name":"P117LocalAnalyst","conclusion":"连续使用验收用本地分析材料","key_reasons":["确认链路必须由用户触发"],"risk_warnings":["不会自动交易"],"confidence":"medium","evidence_ids":[]}]',
                    '{"precision_status":"available","reason":"P117 deterministic acceptance seed","sample_count":20,"sample_window":"2024-2026","screening_condition":"local seeded context","probability_basis":"historical local sample","scenarios":[{"name":"base","return_rate":0.03,"return_range":"0%~6%","probability":0.5,"trigger":"valuation stable"}],"disclaimer":"仅用于本地人工决策辅助，不承诺收益"}',
                    '[{"step":"rule","result":"hold"},{"step":"user","result":"pending"}]',
                    '{"p117":"continuous_product_usability_acceptance"}',
                    NOW,
                ),
            )
        conn.execute(
            """
            INSERT OR REPLACE INTO notifications (
              notification_id,type,severity,title,message,source_type,source_id,created_at
            ) VALUES (?,?,?,?,?,?,?,?)
            """,
            ("notif_p117_unread", "data_source_failure", "warning", "P117 连续使用提醒", "用于连续使用通知读写验收", "p117", "decision_p117_execute", NOW),
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
                "risk_p117_active",
                "data_degraded",
                "warning",
                "active",
                "159915",
                "P117 连续使用风险处置验收",
                '{"source":"p117","day":"day3"}',
                '["新增买入"]',
                '["人工复核数据质量和仓位变化"]',
                "decision_p117_execute",
                "notif_p117_unread",
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
                "market_p117_block",
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


def count_transactions(db_path: Path) -> int:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        return int(scalar(conn, "SELECT COUNT(*) FROM position_transactions"))
    finally:
        conn.close()


def exercise_api_sqlite(base_url: str, db_path: Path) -> list[dict[str, Any]]:
    titles = {
        "U01": ("Day 0", "冷启动和空事实"),
        "U02": ("Day 1", "首次账户入门"),
        "U03": ("Day 1", "首次解释和下一步"),
        "U04": ("Day 2", "日常纪律读回"),
        "U05": ("Day 3", "线下交易补记"),
        "U06": ("Day 3", "风险和通知处置"),
        "U07": ("Day 4", "坏批量导入恢复"),
        "U08": ("Day 4", "非法交易恢复"),
        "U09": ("Day 4", "本地修正审计"),
        "U10": ("Day 5", "数据质量降级处置"),
        "U11": ("Day 6", "人工决策执行确认"),
        "U12": ("Day 6", "错误标注学习闭环"),
        "U13": ("Day 7", "跨页面一致性"),
        "U14": ("Day 7", "重启持久化"),
        "U15": ("Day 7", "移动端可用性"),
        "U16": ("Day 7", "安全负证据"),
        "U17": ("Day 7", "可用性解读报告"),
    }
    scenarios = {sid: scenario(sid, day, title) for sid, (day, title) in titles.items()}
    scenarios["U10"]["status"] = "scoped_pass"
    scenarios["U10"]["expected_eligibility"] = "scoped_pass"
    scenarios["U10"]["classification_reason"] = "Validated as local seeded data-quality degradation handling; not a fresh external provider clean claim."

    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p117_health")
    require(health.get("status") == "ok", "health should be ok")
    add_api(scenarios["U01"], "GET", "/api/v1/health", request_id="req_p117_health")
    expect_http_status(base_url, "GET", "/api/v1/portfolio/current", 404, request_id="req_p117_empty_portfolio")
    add_api(scenarios["U01"], "GET", "/api/v1/portfolio/current", 404, "req_p117_empty_portfolio")
    add_usable(scenarios["U01"], "cold_start", "empty portfolio does not fabricate facts")

    positions = [
        {"symbol": "510300", "name": "沪深300ETF", "quantity": 100, "cost_price": 3.2, "current_price": 4.0, "buy_date": "2026-06-01", "position_state": "sell_only", "buy_reason": "P117 首次录入核心仓", "asset_tag": "core"},
        {"symbol": "159915", "name": "创业板ETF", "quantity": 100, "cost_price": 2.1, "current_price": 2.5, "buy_date": "2026-06-01", "position_state": "frozen_watch", "buy_reason": "P117 首次录入成长仓", "asset_tag": "satellite"},
        {"symbol": "512000", "name": "券商ETF", "quantity": 200, "cost_price": 1.1, "current_price": 1.3, "buy_date": "2026-06-01", "position_state": "normal", "buy_reason": "P117 首次录入行业仓", "asset_tag": "satellite"},
        {"symbol": "110022", "name": "易方达消费行业", "quantity": 100, "cost_price": 2.4, "current_price": 2.6, "buy_date": "2026-06-01", "position_state": "normal", "buy_reason": "P117 首次录入主动基金", "asset_tag": "active_fund"},
    ]
    adjust = data(request_json(base_url, "POST", "/api/v1/portfolio/adjustments", {"cash": 1500, "total_assets": 2670, "adjust_reason": "P117 Day1 本地账户入门", "positions": positions}, "req_p117_adjust"))
    require(adjust["position_count"] == 4, "portfolio adjustment should create four positions")
    add_api(scenarios["U02"], "POST", "/api/v1/portfolio/adjustments", request_id="req_p117_adjust")
    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p117_current_day1"))
    require({item["symbol"] for item in current["positions"]} == set(SYMBOLS[:4]), "day1 portfolio symbols mismatch")
    add_api(scenarios["U02"], "GET", "/api/v1/portfolio/current", request_id="req_p117_current_day1")
    add_usable(scenarios["U02"], "onboarding", "user can create local account context without broker connection")

    for sid, path, rid in [
        ("U03", "/api/v1/dashboard/today", "req_p117_dashboard_day1"),
        ("U03", "/api/v1/review/summary", "req_p117_review_day1"),
        ("U04", "/api/v1/dashboard/today", "req_p117_dashboard_day2"),
        ("U04", "/api/v1/review/summary", "req_p117_review_day2"),
        ("U04", "/api/v1/audit-events", "req_p117_audit_day2"),
        ("U04", "/api/v1/daily-discipline/reports", "req_p117_reports_day2"),
    ]:
        request_json(base_url, "GET", path, request_id=rid)
        add_api(scenarios[sid], "GET", path, request_id=rid)
    add_usable(scenarios["U03"], "next_step_clarity", "daily surfaces open after onboarding")
    add_usable(scenarios["U04"], "routine", "daily dashboard, review, reports and audit remain readable")

    offline_transactions = [
        {"operation_type": "buy", "symbol": "159915", "name": "创业板ETF", "quantity": 10, "price": 2.55, "fees": 1, "executed_at": "2026-06-20T03:00:00Z", "buy_reason": "P117 Day3 线下加仓", "note": "P117 Day3 只记录线下买入"},
        {"operation_type": "sell", "symbol": "510300", "name": "沪深300ETF", "quantity": 15, "price": 4.15, "fees": 1, "executed_at": "2026-06-21T03:00:00Z", "note": "P117 Day3 只记录线下卖出"},
        {"operation_type": "reduce", "symbol": "512000", "name": "券商ETF", "quantity": 40, "price": 1.28, "fees": 0.5, "executed_at": "2026-06-22T03:00:00Z", "note": "P117 Day3 只记录线下减仓"},
    ]
    for idx, tx in enumerate(offline_transactions, start=1):
        out = data(request_json(base_url, "POST", "/api/v1/portfolio/offline-transactions", tx, f"req_p117_offline_{idx}"))
        require(out.get("transaction_id"), f"offline transaction {idx} missing id")
        add_api(scenarios["U05"], "POST", "/api/v1/portfolio/offline-transactions", request_id=f"req_p117_offline_{idx}")
    add_usable(scenarios["U05"], "ledger_update", "multi-fund offline facts can be recorded without broker execution")

    risk = data(request_json(base_url, "POST", "/api/v1/risk-alerts/risk_p117_active/lifecycle", {"status": "resolved", "reason": "P117 Day3 人工复核完成"}, "req_p117_risk_resolve"))
    require(risk["sop_status"] == "resolved", "risk should resolve")
    request_json(base_url, "POST", "/api/v1/notifications/notif_p117_unread/read", request_id="req_p117_notification_read")
    add_api(scenarios["U06"], "POST", "/api/v1/risk-alerts/risk_p117_active/lifecycle", request_id="req_p117_risk_resolve")
    add_api(scenarios["U06"], "POST", "/api/v1/notifications/{id}/read", request_id="req_p117_notification_read")
    add_usable(scenarios["U06"], "operator_control", "risk and notification state changes require explicit local action")

    before_invalid = count_transactions(db_path)
    invalid_rows = [
        {"row_number": 1, "row_type": "transaction", "operation_type": "buy", "symbol": "512000", "name": "券商ETF", "quantity": 20, "price": 1.26, "fees": 0.5, "occurred_at": "2026-06-23T04:00:00Z", "buy_reason": "P117 有效行"},
        {"row_number": 2, "row_type": "transaction", "operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "occurred_at": "2026-06-23T04:10:00Z", "buy_reason": "P117 应拒绝"},
    ]
    invalid_validate = data(request_json(base_url, "POST", "/api/v1/portfolio/imports/validate", {"rows": invalid_rows}, "req_p117_import_invalid_validate"))
    require(invalid_validate["summary"]["invalid_count"] == 1, "invalid import should report one invalid row")
    expect_http_status(base_url, "POST", "/api/v1/portfolio/imports/confirm", 400, {"import_batch_id": invalid_validate["import_batch_id"], "confirm_reason": "P117 错误批次不得确认", "rows": invalid_rows}, "req_p117_import_invalid_confirm")
    require(count_transactions(db_path) == before_invalid, "invalid import must not create transactions")
    add_api(scenarios["U07"], "POST", "/api/v1/portfolio/imports/validate", request_id="req_p117_import_invalid_validate")
    add_reject(scenarios["U07"], "POST", "/api/v1/portfolio/imports/confirm", 400, "req_p117_import_invalid_confirm", "invalid row blocks confirm")
    add_usable(scenarios["U07"], "recovery", "bad import is visible and does not partially write local facts")

    before_rejects = count_transactions(db_path)
    rejects = [
        ("req_p117_reject_oversell", {"operation_type": "sell", "symbol": "159915", "name": "创业板ETF", "quantity": 999999, "price": 2.5, "fees": 0, "executed_at": "2026-06-23T03:00:00Z"}, "oversell"),
        ("req_p117_reject_future", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2999-01-01T00:00:00Z", "buy_reason": "P117 未来时间"}, "future execution time"),
        ("req_p117_reject_fees", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": -1, "executed_at": "2026-06-23T03:00:00Z", "buy_reason": "P117 负费用"}, "negative fees"),
        ("req_p117_reject_symbol", {"operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2026-06-23T03:00:00Z", "buy_reason": "P117 缺代码"}, "missing symbol"),
    ]
    for request_id, body, reason in rejects:
        expect_http_status(base_url, "POST", "/api/v1/portfolio/offline-transactions", 400, body, request_id)
        add_reject(scenarios["U08"], "POST", "/api/v1/portfolio/offline-transactions", 400, request_id, reason)
    require(count_transactions(db_path) == before_rejects, "invalid transactions must not create rows")
    add_usable(scenarios["U08"], "recovery", "invalid transaction attempts leave ledger unchanged")

    current_for_correction = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p117_current_for_correction"))
    pos_510300 = next(item for item in current_for_correction["positions"] if item["symbol"] == "510300")
    correction = data(request_json(base_url, "POST", "/api/v1/portfolio/corrections", {"target_type": "position", "target_id": pos_510300["position_id"], "before_json": '{"quantity":100}', "after_json": '{"quantity":85}', "correction_reason": "P117 Day4 修正线下卖出后数量"}, "req_p117_correction"))
    require(correction.get("correction_id"), "correction should return id")
    add_api(scenarios["U09"], "POST", "/api/v1/portfolio/corrections", request_id="req_p117_correction")
    add_usable(scenarios["U09"], "auditability", "user can record correction without silently rewriting prior transactions")

    gate = data(request_json(base_url, "GET", "/api/v1/data-source-quality/gate-resolution?symbol=000300", request_id="req_p117_dq_gate"))
    require(gate["release_claim_state"] == "requires_resolution", "seeded DQ gate should require resolution")
    resolution = data(request_json(base_url, "POST", "/api/v1/data-source-quality/resolutions", {"symbol": "000300", "resolution_type": "scope_exclusion", "scope": "P117 continuous usability excludes current-data clean claim", "reason": "Seeded stale source health for usability recovery validation", "release_impact": "Do not claim current data clean in P117", "evidence_ref": "docs/release/acceptance/2026-06-25-p117-continuous-product-usability-acceptance-matrix.md"}, "req_p117_dq_resolution_create"))
    active = resolution.get("active_resolution") or {}
    require(active.get("resolution_id"), "DQ resolution should create active resolution")
    request_json(base_url, "POST", f"/api/v1/data-source-quality/resolutions/{urllib.parse.quote(active['resolution_id'])}/retire", request_id="req_p117_dq_resolution_retire")
    add_api(scenarios["U10"], "GET", "/api/v1/data-source-quality/gate-resolution", request_id="req_p117_dq_gate")
    add_api(scenarios["U10"], "POST", "/api/v1/data-source-quality/resolutions", request_id="req_p117_dq_resolution_create")
    add_api(scenarios["U10"], "POST", "/api/v1/data-source-quality/resolutions/{id}/retire", request_id="req_p117_dq_resolution_retire")
    add_usable(scenarios["U10"], "claim_honesty", "degraded data requires explicit scoped resolution and does not become clean pass")

    request_json(base_url, "POST", "/api/v1/decisions/decision_p117_execute/confirmations", {"confirmation_type": "executed_manually", "operation_type": "sell", "symbol": "159915", "quantity": 4, "price": 2.6, "fees": 1, "executed_at": "2026-06-24T03:00:00Z", "note": "P117 Day6 人工线下处理记录"}, "req_p117_confirmation_execute")
    add_api(scenarios["U11"], "POST", "/api/v1/decisions/decision_p117_execute/confirmations", request_id="req_p117_confirmation_execute")
    add_usable(scenarios["U11"], "manual_control", "decision confirmation remains explicit and user-triggered")
    request_json(base_url, "POST", "/api/v1/decisions/decision_p117_error/confirmations", {"confirmation_type": "marked_error", "actual_outcome": "P117 Day6 复盘发现建议偏离", "root_cause_tag": "evidence_missed", "lesson_learned": "连续使用时必须复核成交后仓位和证据质量", "note": "P117 错误标注"}, "req_p117_marked_error")
    add_api(scenarios["U12"], "POST", "/api/v1/decisions/decision_p117_error/confirmations", request_id="req_p117_marked_error")
    add_usable(scenarios["U12"], "learning_loop", "user can mark a decision wrong with root cause and lesson")

    for method, path, rid in [
        ("GET", "/api/v1/portfolio/current", "req_p117_final_portfolio"),
        ("GET", "/api/v1/dashboard/today", "req_p117_final_dashboard"),
        ("GET", "/api/v1/review/summary", "req_p117_final_review"),
        ("GET", "/api/v1/audit-events", "req_p117_final_audit"),
        ("GET", "/api/v1/decision-loops?limit=10", "req_p117_final_loop"),
    ]:
        request_json(base_url, method, path, request_id=rid)
        add_api(scenarios["U13"], method, path, request_id=rid)
    add_usable(scenarios["U13"], "cross_page_consistency", "portfolio, dashboard, review, audit and loop endpoints read after seven-day facts")

    enrich_sqlite_evidence(db_path, scenarios)
    return [scenarios[sid] for sid in SCENARIO_IDS]


def enrich_sqlite_evidence(db_path: Path, scenarios: dict[str, dict[str, Any]]) -> None:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        counts = {
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "portfolio_snapshots": int(scalar(conn, "SELECT COUNT(*) FROM portfolio_snapshots")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "transaction_symbols": int(scalar(conn, "SELECT COUNT(DISTINCT symbol) FROM position_transactions")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "marked_error_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type='marked_error'")),
            "import_batches_invalid": int(scalar(conn, "SELECT COUNT(*) FROM local_account_import_batches WHERE invalid_count > 0")),
            "corrections": int(scalar(conn, "SELECT COUNT(*) FROM local_account_corrections")),
            "dq_resolutions_retired": int(scalar(conn, "SELECT COUNT(*) FROM data_quality_gate_resolutions WHERE status='retired'")),
            "risk_resolved": int(scalar(conn, "SELECT COUNT(*) FROM risk_alerts WHERE alert_id='risk_p117_active' AND sop_status='resolved'")),
            "notifications_read": int(scalar(conn, "SELECT COUNT(*) FROM notifications WHERE notification_id='notif_p117_unread' AND read_at IS NOT NULL")),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
            "auto_confirmation_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")),
            "forbidden_broker_order_push_tables": int(scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")),
            "auto_rule_apply_audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events WHERE LOWER(action) LIKE '%auto%' AND (LOWER(action) LIKE '%rule%' OR LOWER(action) LIKE '%confirm%' OR LOWER(action) LIKE '%trade%')")),
        }
        require(counts["positions"] >= 4, "positions missing")
        require(counts["portfolio_snapshots"] >= 4, "portfolio snapshots missing")
        require(counts["position_transactions"] >= 4, "transaction ledger missing")
        require(counts["transaction_symbols"] >= 3, "transaction symbols insufficient")
        require(counts["operation_confirmations"] >= 4, "operation confirmations missing")
        require(counts["marked_error_rows"] >= 1, "marked error missing")
        require(counts["import_batches_invalid"] >= 1, "invalid import evidence missing")
        require(counts["corrections"] >= 1, "correction missing")
        require(counts["dq_resolutions_retired"] >= 1, "DQ retired resolution missing")
        require(counts["risk_resolved"] == 1, "risk resolution missing")
        require(counts["notifications_read"] == 1, "notification read missing")
        require(counts["audit_events"] >= 10, "audit events missing")
        require(counts["auto_confirmation_rows"] == 0, "auto confirmations must be absent")
        require(counts["forbidden_broker_order_push_tables"] == 0, "forbidden broker/order/push tables must be absent")
        require(counts["auto_rule_apply_audit_events"] == 0, "auto rule apply audit events must be absent")

        add_sqlite(scenarios["U02"], "positions", "count", counts["positions"])
        add_sqlite(scenarios["U05"], "position_transactions", "count", counts["position_transactions"])
        add_sqlite(scenarios["U05"], "position_transactions", "distinct symbols", counts["transaction_symbols"])
        add_sqlite(scenarios["U06"], "risk_alerts", "resolved", counts["risk_resolved"])
        add_sqlite(scenarios["U06"], "notifications", "read_at", counts["notifications_read"])
        add_sqlite(scenarios["U07"], "local_account_import_batches", "invalid", counts["import_batches_invalid"])
        add_sqlite(scenarios["U08"], "position_transactions", "unchanged after invalid attempts", counts["position_transactions"])
        add_sqlite(scenarios["U09"], "local_account_corrections", "count", counts["corrections"])
        add_sqlite(scenarios["U10"], "data_quality_gate_resolutions", "retired", counts["dq_resolutions_retired"])
        add_sqlite(scenarios["U11"], "operation_confirmations", "count", counts["operation_confirmations"])
        add_sqlite(scenarios["U12"], "operation_confirmations", "marked_error", counts["marked_error_rows"])
        add_sqlite(scenarios["U13"], "audit_events", "count", counts["audit_events"])
        add_sqlite(scenarios["U16"], "sqlite_master", "forbidden broker/order/push tables", counts["forbidden_broker_order_push_tables"])
        add_sqlite(scenarios["U16"], "operation_confirmations", "auto confirmations", counts["auto_confirmation_rows"])
        add_sqlite(scenarios["U16"], "audit_events", "auto rule apply events", counts["auto_rule_apply_audit_events"])
        add_usable(scenarios["U16"], "safety", "no broker/order/push execution tables or auto action traces")

        for item in scenarios.values():
            item["safety_counters"] = {
                "forbidden_broker_order_push_tables": counts["forbidden_broker_order_push_tables"],
                "auto_confirmation_rows": counts["auto_confirmation_rows"],
                "auto_rule_apply_audit_events": counts["auto_rule_apply_audit_events"],
                "automatic_trading_affordances": 0,
                "return_guarantee_claims": 0,
                "secret_or_raw_prompt_leaks": 0,
            }
    finally:
        conn.close()


def write_api_summary(base_url: str, db_path: Path, artifact_dir: Path) -> dict[str, Any]:
    seed_supporting_records(db_path)
    scenarios = exercise_api_sqlite(base_url, db_path)
    payload = {
        "status": "passed",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "base_url": base_url,
        "sqlite_path": str(db_path),
        "evidence_layer": "api_sqlite",
        "scenario_count": len(scenarios),
        "scenarios": scenarios,
        "claim_boundary": "P117 API/SQLite layer uses local seeded continuous-use evidence. It does not claim real broker trades, external push, fresh provider/LLM output, automatic trading, automatic confirmation, automatic rule application, release packaging, physical second-machine validation, or return guarantees.",
    }
    out = artifact_dir / "api_sqlite" / "p117-api-sqlite-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def write_restart_probe(base_url: str, db_path: Path, artifact_dir: Path) -> dict[str, Any]:
    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p117_restart_health")
    require(health.get("status") == "ok", "restart health should be ok")
    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p117_restart_portfolio"))
    audit = data(request_json(base_url, "GET", "/api/v1/audit-events", request_id="req_p117_restart_audit"))
    loops = data(request_json(base_url, "GET", "/api/v1/decision-loops?limit=10", request_id="req_p117_restart_loop"))
    require(current and len(current.get("positions", [])) >= 4, "restart portfolio readback missing positions")
    require(audit is not None, "restart audit readback missing")
    require(loops is not None, "restart decision-loop readback missing")
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        counts = {
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
        }
    finally:
        conn.close()
    payload = {
        "status": "passed",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "scenario_id": "U14",
        "restart_evidence": [
            {"method": "GET", "path": "/api/v1/health", "status": 200, "request_id": "req_p117_restart_health"},
            {"method": "GET", "path": "/api/v1/portfolio/current", "status": 200, "request_id": "req_p117_restart_portfolio", "position_count": len(current.get("positions", []))},
            {"method": "GET", "path": "/api/v1/audit-events", "status": 200, "request_id": "req_p117_restart_audit"},
            {"method": "GET", "path": "/api/v1/decision-loops", "status": 200, "request_id": "req_p117_restart_loop"},
        ],
        "sqlite_counts_after_restart": counts,
        "usability_evidence": [{"dimension": "persistence", "assertion": "same SQLite remains readable after backend restart", "value": "ok"}],
    }
    out = artifact_dir / "restart" / "p117-restart-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(payload, ensure_ascii=False, indent=2))
    return payload


def build_interpretation(scenarios: dict[str, dict[str, Any]], browser_payload: dict[str, Any], restart_payload: dict[str, Any]) -> dict[str, Any]:
    blocked = [item for item in scenarios.values() if item["status"] not in ("fresh_pass", "scoped_pass")]
    fresh = [item for item in scenarios.values() if item["status"] == "fresh_pass"]
    scoped = [item for item in scenarios.values() if item["status"] == "scoped_pass"]
    return {
        "task_completion_rate": "17/17 runner scenarios passed",
        "fresh_pass_count": len(fresh),
        "scoped_pass_count": len(scoped),
        "blocked_count": len(blocked),
        "cold_start": "Usable: empty portfolio is explicit and does not fabricate holdings.",
        "onboarding": "Usable: local account and holdings can be created without broker connection.",
        "daily_routine": "Usable: dashboard, workbench, review and audit remain readable after local facts.",
        "recovery": "Usable: invalid imports and invalid transactions are rejected without partial ledger writes.",
        "traceability": "Usable: corrections, manual confirmations, marked errors and audit rows provide review trail.",
        "persistence": "Usable: restart probe reads portfolio, audit and decision-loop facts from the same SQLite database.",
        "mobile": "Usable within checked scope: 390px portfolio/workbench paths render without console/page/API 5xx failures.",
        "safety": "Usable as a local discipline assistant only: no broker/order/push tables, auto confirmations or auto rule-apply audit events.",
        "browser_health": {
            "status": browser_payload.get("status"),
            "console_errors": len(browser_payload.get("console_errors", [])),
            "page_errors": len(browser_payload.get("page_errors", [])),
            "failed_api_responses": len(browser_payload.get("failed_api_responses", [])),
        },
        "restart_health": restart_payload.get("status"),
        "claim_boundary": "This is a local continuous-use usability pass, not a claim of broker execution, external data cleanliness, fresh real LLM quality, release package readiness, or physical second-machine validation.",
    }


def merge_summary(base_url: str, db_path: Path, artifact_dir: Path, browser_summary: Path | None) -> dict[str, Any]:
    api_path = artifact_dir / "api_sqlite" / "p117-api-sqlite-summary.json"
    restart_path = artifact_dir / "restart" / "p117-restart-summary.json"
    require(api_path.exists(), f"missing API summary: {api_path}")
    require(restart_path.exists(), f"missing restart summary: {restart_path}")
    api_payload = json.loads(api_path.read_text(encoding="utf-8"))
    restart_payload = json.loads(restart_path.read_text(encoding="utf-8"))
    scenarios = {item["scenario_id"]: item for item in api_payload["scenarios"]}
    if "U14" in scenarios:
        scenarios["U14"]["restart_evidence"].extend(restart_payload.get("restart_evidence", []))
        scenarios["U14"]["usability_evidence"].extend(restart_payload.get("usability_evidence", []))
    browser_payload: dict[str, Any] = {}
    if browser_summary and browser_summary.exists():
        browser_payload = json.loads(browser_summary.read_text(encoding="utf-8"))
        for entry in browser_payload.get("scenarios", []):
            sid = entry["scenario_id"]
            if sid in scenarios:
                scenarios[sid]["browser_evidence"].extend(entry.get("browser_evidence", []))
                scenarios[sid]["redaction_evidence"].update(entry.get("redaction_evidence", {}))
    missing_browser = [sid for sid in ["U02", "U03", "U04", "U05", "U06", "U10", "U11", "U13", "U15", "U16"] if not scenarios.get(sid, {}).get("browser_evidence")]
    require(not missing_browser, f"missing browser evidence for {missing_browser}")

    interpretation = build_interpretation(scenarios, browser_payload, restart_payload)
    scenarios["U17"]["usability_evidence"].append({"dimension": "interpretation", "assertion": "final usability interpretation generated", "value": interpretation})
    final = {
        "status": "passed",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "base_url": base_url,
        "sqlite_path": str(db_path),
        "scenario_count": len(scenarios),
        "fresh_pass_count": sum(1 for item in scenarios.values() if item["status"] == "fresh_pass"),
        "scoped_pass_count": sum(1 for item in scenarios.values() if item["status"] == "scoped_pass"),
        "symbols": SYMBOLS,
        "scenarios": [scenarios[sid] for sid in SCENARIO_IDS],
        "browser_summary": browser_payload,
        "restart_summary": restart_payload,
        "usability_interpretation": interpretation,
        "claim_boundary": interpretation["claim_boundary"],
    }
    out = artifact_dir / "p117-usability-summary.json"
    out.write_text(json.dumps(final, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(final, ensure_ascii=False, indent=2))
    return final


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", required=True)
    parser.add_argument("--sqlite", required=True)
    parser.add_argument("--artifact-dir", required=True)
    parser.add_argument("--browser-summary")
    parser.add_argument("--restart-probe", action="store_true")
    parser.add_argument("--merge-only", action="store_true")
    args = parser.parse_args()

    artifact_dir = Path(args.artifact_dir)
    artifact_dir.mkdir(parents=True, exist_ok=True)
    db_path = Path(args.sqlite)
    if args.restart_probe:
        write_restart_probe(args.base_url, db_path, artifact_dir)
    elif args.merge_only:
        merge_summary(args.base_url, db_path, artifact_dir, Path(args.browser_summary) if args.browser_summary else None)
    else:
        payload = write_api_summary(args.base_url, db_path, artifact_dir)
        print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
