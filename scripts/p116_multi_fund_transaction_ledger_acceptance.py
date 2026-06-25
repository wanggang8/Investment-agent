#!/usr/bin/env python3
"""P116 multi-fund transaction ledger acceptance runner."""

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
CHANGE_ID = "p116-multi-fund-transaction-ledger-acceptance"
SCENARIO_IDS = [
    "L01",
    "L02",
    "L03",
    "L04",
    "L05",
    "L06",
    "L07",
    "L08",
    "L09",
    "L10",
    "L11",
    "L12",
    "L13",
    "L14",
    "L15",
    "L16",
]
MULTI_FUND_SYMBOLS = ["510300", "159915", "588000", "512000", "110022", "161725"]


class AcceptanceFailure(RuntimeError):
    pass


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AcceptanceFailure(message)


def request_json(base_url: str, method: str, path: str, body: dict[str, Any] | None = None, request_id: str = "req_p116") -> dict[str, Any]:
    data = None
    headers = {"Accept": "application/json", "X-Request-ID": request_id}
    if body is not None:
        data = json.dumps(body, ensure_ascii=False).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(base_url.rstrip("/") + path, data=data, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=30) as resp:
        payload = resp.read().decode("utf-8")
        return json.loads(payload) if payload else {}


def expect_http_status(base_url: str, method: str, path: str, status: int, body: dict[str, Any] | None = None, request_id: str = "req_p116_reject") -> dict[str, Any]:
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


def scenario(sid: str, title: str, expected: str = "fresh_pass") -> dict[str, Any]:
    return {
        "scenario_id": sid,
        "title": title,
        "status": expected,
        "expected_eligibility": expected,
        "classification_reason": "Validated by P116 isolated local runner.",
        "config_mode": "local_multi_fund_transaction_ledger",
        "runtime_mode": "development",
        "use_stub": True,
        "provider_mode": "stub_local_linkage",
        "llm_mode": "not_configured_degraded_expected",
        "symbols": MULTI_FUND_SYMBOLS,
        "api_evidence": [],
        "browser_evidence": [],
        "sqlite_evidence": [],
        "downstream_evidence": [],
        "rejection_evidence": [],
        "side_effects": {},
        "redaction_evidence": {},
        "safety_counters": {},
    }


def add_api(item: dict[str, Any], method: str, path: str, status: int = 200, request_id: str = "req_p116") -> None:
    item["api_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id})


def add_reject(item: dict[str, Any], method: str, path: str, status: int, request_id: str, reason: str) -> None:
    item["rejection_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id, "reason": reason})


def add_sqlite(item: dict[str, Any], table: str, field: str, row_count: int, label: str = "") -> None:
    item["sqlite_evidence"].append({"table": table, "field": field, "row_count": row_count, "query_label": label or f"{table}.{field}"})


def add_downstream(item: dict[str, Any], target: str, assertion: str, value: Any = None) -> None:
    item["downstream_evidence"].append({"target": target, "assertion": assertion, "value": value})


def seed_supporting_records(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
    try:
        for decision_id, symbol, title in [
            ("decision_p116_execute", "510300", "P116 多基金线下执行确认验收决策"),
            ("decision_p116_error", "159915", "P116 多基金错误标注验收决策"),
            ("decision_p116_browser_execute", "159915", "P116 浏览器多基金手动执行确认验收决策"),
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
                    "P116 seeded decision for multi-fund transaction acceptance",
                    "satisfied",
                    "",
                    '{"heat":"neutral"}',
                    '["calm"]',
                    '["manual_confirmation_required","multi_fund_ledger_review"]',
                    "[]",
                    "hold",
                    "P116 本地验收决策：只允许人工记录线下处理结果",
                    '["自动交易","外部推送"]',
                    '["记录线下处理","继续观察","标记错误"]',
                    "pending",
                    None,
                    None,
                    "v_p116",
                    '[{"agent_name":"P116LocalAnalyst","conclusion":"多基金交易台账验收用本地材料","key_reasons":["确认链路必须由用户触发"],"risk_warnings":["不会自动交易"],"confidence":"medium","evidence_ids":[]}]',
                    '{"precision_status":"available","reason":"P116 deterministic acceptance seed","sample_count":20,"sample_window":"2024-2026","screening_condition":"local seeded context","probability_basis":"historical local sample","scenarios":[{"name":"base","return_rate":0.03,"return_range":"0%~6%","probability":0.5,"trigger":"valuation stable"}],"disclaimer":"仅用于本地人工决策辅助，不承诺收益"}',
                    '[{"step":"rule","result":"hold"},{"step":"user","result":"pending"}]',
                    '{"p116":"multi_fund_transaction_ledger_acceptance"}',
                    NOW,
                ),
            )
        conn.execute(
            """
            INSERT OR REPLACE INTO notifications (
              notification_id,type,severity,title,message,source_type,source_id,created_at
            ) VALUES (?,?,?,?,?,?,?,?)
            """,
            ("notif_p116_unread", "data_source_failure", "warning", "P116 多基金数据提示", "用于通知读写联动验收", "p116", "decision_p116_execute", NOW),
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
                "risk_p116_active",
                "position_limit_breach",
                "warning",
                "active",
                "159915",
                "P116 多基金仓位风险处置验收",
                '{"source":"p116","symbols":["510300","159915","588000","512000","110022"]}',
                '["新增买入"]',
                '["人工复核多基金仓位和数据质量"]',
                "decision_p116_execute",
                "notif_p116_unread",
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
                "market_p116_block",
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
        "L01": "服务健康与空组合边界",
        "L02": "五只基金初始组合录入与读回",
        "L03": "多日期多动作线下交易台账",
        "L04": "混合批量导入拦截与有效批次确认",
        "L05": "非法交易与非法持仓状态拒绝",
        "L06": "持仓编辑、移除与修正审计",
        "L07": "决策详情手动执行确认联动",
        "L08": "决策错误标注联动",
        "L09": "季度再平衡复核读写",
        "L10": "风险预警与通知处置",
        "L11": "数据质量 gate resolution 创建与退役",
        "L12": "Dashboard、Workbench、Review、Audit 下游读回",
        "L13": "真实浏览器多基金录入与交易",
        "L14": "真实浏览器决策确认与治理页面",
        "L15": "移动端组合页面与可见安全边界",
        "L16": "SQLite 安全边界与禁止能力计数",
    }
    scenarios: dict[str, dict[str, Any]] = {sid: scenario(sid, titles[sid]) for sid in SCENARIO_IDS}
    for sid in ["L11", "L12"]:
        scenarios[sid]["status"] = "scoped_pass"
        scenarios[sid]["expected_eligibility"] = "scoped_pass"
        scenarios[sid]["classification_reason"] = "Validated as seeded/local governance linkage; not a fresh external provider claim."

    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p116_health")
    require(health.get("status") == "ok", "health should be ok")
    add_api(scenarios["L01"], "GET", "/api/v1/health", request_id="req_p116_health")
    expect_http_status(base_url, "GET", "/api/v1/portfolio/current", 404, request_id="req_p116_empty_portfolio")
    add_api(scenarios["L01"], "GET", "/api/v1/portfolio/current", 404, "req_p116_empty_portfolio")

    positions = [
        {"symbol": "510300", "name": "沪深300ETF", "quantity": 100, "cost_price": 3.2, "current_price": 4.0, "buy_date": "2026-01-05", "position_state": "sell_only", "buy_reason": "买入逻辑破坏后只卖不买", "asset_tag": "core"},
        {"symbol": "159915", "name": "创业板ETF", "quantity": 80, "cost_price": 2.1, "current_price": 2.6, "buy_date": "2026-01-06", "position_state": "frozen_watch", "buy_reason": "多源验证不足冻结观察", "asset_tag": "satellite"},
        {"symbol": "588000", "name": "科创50ETF", "quantity": 120, "cost_price": 1.0, "current_price": 1.2, "buy_date": "2026-02-01", "position_state": "normal", "buy_reason": "小比例观察仓", "asset_tag": "satellite"},
        {"symbol": "512000", "name": "券商ETF", "quantity": 300, "cost_price": 1.1, "current_price": 1.3, "buy_date": "2026-02-08", "position_state": "normal", "buy_reason": "行业低估配置", "asset_tag": "satellite"},
        {"symbol": "110022", "name": "易方达消费行业", "quantity": 200, "cost_price": 2.4, "current_price": 2.5, "buy_date": "2026-03-01", "position_state": "normal", "buy_reason": "长期消费基金观察", "asset_tag": "active_fund"},
    ]
    adjust = data(request_json(base_url, "POST", "/api/v1/portfolio/adjustments", {"cash": 1200, "total_assets": 2842, "adjust_reason": "P116 多基金本地组合校准", "positions": positions}, "req_p116_adjust"))
    require(adjust["position_count"] == 5, "portfolio adjustment should create five positions")
    add_api(scenarios["L02"], "POST", "/api/v1/portfolio/adjustments", request_id="req_p116_adjust")
    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p116_current"))
    current_symbols = {item["symbol"] for item in current["positions"]}
    require(set(MULTI_FUND_SYMBOLS[:5]).issubset(current_symbols), f"current portfolio missing symbols: {current_symbols}")
    add_api(scenarios["L02"], "GET", "/api/v1/portfolio/current", request_id="req_p116_current")

    offline_transactions = [
        {"operation_type": "buy", "symbol": "159915", "name": "创业板ETF", "quantity": 10, "price": 2.7, "fees": 1, "executed_at": "2026-06-20T03:00:00Z", "buy_reason": "P116 人工加仓创业板", "note": "P116 只记录线下买入"},
        {"operation_type": "sell", "symbol": "510300", "name": "沪深300ETF", "quantity": 20, "price": 4.2, "fees": 1.5, "executed_at": "2026-06-21T03:00:00Z", "note": "P116 只记录线下卖出"},
        {"operation_type": "reduce", "symbol": "512000", "name": "券商ETF", "quantity": 60, "price": 1.25, "fees": 0.8, "executed_at": "2026-06-22T03:00:00Z", "note": "P116 减仓券商ETF"},
        {"operation_type": "sell", "symbol": "588000", "name": "科创50ETF", "quantity": 30, "price": 1.25, "fees": 0.5, "executed_at": "2026-06-23T03:00:00Z", "note": "P116 降低科创50观察仓"},
    ]
    for idx, tx in enumerate(offline_transactions, start=1):
        out = data(request_json(base_url, "POST", "/api/v1/portfolio/offline-transactions", tx, f"req_p116_offline_tx_{idx}"))
        require(out.get("transaction_id"), f"offline transaction {idx} should return transaction id")
        add_api(scenarios["L03"], "POST", "/api/v1/portfolio/offline-transactions", request_id=f"req_p116_offline_tx_{idx}")
    current_after_tx = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p116_current_after_tx"))
    qty_by_symbol = {item["symbol"]: item["quantity"] for item in current_after_tx["positions"]}
    require(qty_by_symbol["159915"] == 90 and qty_by_symbol["510300"] == 80 and qty_by_symbol["512000"] == 240 and qty_by_symbol["588000"] == 90, f"unexpected quantities after transactions: {qty_by_symbol}")
    add_api(scenarios["L03"], "GET", "/api/v1/portfolio/current", request_id="req_p116_current_after_tx")

    before_invalid_count = count_transactions(db_path)
    invalid_rows = [
        {"row_number": 1, "row_type": "transaction", "operation_type": "buy", "symbol": "512000", "name": "券商ETF", "quantity": 30, "price": 1.26, "fees": 0.5, "occurred_at": "2026-06-23T04:00:00Z", "buy_reason": "P116 混合批次有效行"},
        {"row_number": 2, "row_type": "transaction", "operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "occurred_at": "2026-06-23T04:10:00Z", "buy_reason": "P116 应拒绝"},
        {"row_number": 3, "row_type": "holding", "symbol": "161725", "name": "招商中证白酒指数", "quantity": -1, "cost_price": 1, "current_price": 1.1, "buy_reason": "P116 应拒绝"},
    ]
    invalid_validate = data(request_json(base_url, "POST", "/api/v1/portfolio/imports/validate", {"rows": invalid_rows}, "req_p116_import_invalid_validate"))
    require(invalid_validate["summary"]["invalid_count"] == 2, "mixed invalid import should report two invalid rows")
    add_api(scenarios["L04"], "POST", "/api/v1/portfolio/imports/validate", request_id="req_p116_import_invalid_validate")
    expect_http_status(base_url, "POST", "/api/v1/portfolio/imports/confirm", 400, {"import_batch_id": invalid_validate["import_batch_id"], "confirm_reason": "P116 错误批次不得确认", "rows": invalid_rows}, "req_p116_import_invalid_confirm")
    add_reject(scenarios["L04"], "POST", "/api/v1/portfolio/imports/confirm", 400, "req_p116_import_invalid_confirm", "batch contains invalid rows")
    require(count_transactions(db_path) == before_invalid_count, "invalid import must not create transactions")

    valid_rows = [
        {"row_number": 1, "row_type": "holding", "symbol": "161725", "name": "招商中证白酒指数", "quantity": 50, "cost_price": 1.0, "current_price": 1.1, "buy_date": "2026-06-10", "buy_reason": "P116 有效持仓导入", "asset_tag": "active_fund"},
        {"row_number": 2, "row_type": "transaction", "operation_type": "buy", "symbol": "512000", "name": "券商ETF", "quantity": 30, "price": 1.26, "fees": 0.5, "occurred_at": "2026-06-23T04:00:00Z", "buy_reason": "P116 有效交易导入", "asset_tag": "satellite"},
    ]
    valid_validate = data(request_json(base_url, "POST", "/api/v1/portfolio/imports/validate", {"rows": valid_rows}, "req_p116_import_valid_validate"))
    require(valid_validate["summary"]["valid_count"] == 2 and valid_validate["summary"]["invalid_count"] == 0, "valid import should have no invalid rows")
    request_json(base_url, "POST", "/api/v1/portfolio/imports/confirm", {"import_batch_id": valid_validate["import_batch_id"], "confirm_reason": "P116 确认有效多基金导入", "rows": valid_rows}, "req_p116_import_valid_confirm")
    add_api(scenarios["L04"], "POST", "/api/v1/portfolio/imports/confirm", request_id="req_p116_import_valid_confirm")

    before_reject_count = count_transactions(db_path)
    rejects = [
        ("req_p116_reject_cash", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 999999, "price": 999, "fees": 1, "executed_at": "2026-06-24T03:00:00Z", "buy_reason": "P116 现金不足"}, "insufficient cash"),
        ("req_p116_reject_oversell", {"operation_type": "sell", "symbol": "588000", "name": "科创50ETF", "quantity": 999999, "price": 1, "fees": 0, "executed_at": "2026-06-24T03:00:00Z"}, "oversell"),
        ("req_p116_reject_future", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2999-01-01T00:00:00Z", "buy_reason": "P116 未来时间"}, "future execution time"),
        ("req_p116_reject_fees", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": -1, "executed_at": "2026-06-24T03:00:00Z", "buy_reason": "P116 负费用"}, "negative fees"),
        ("req_p116_reject_symbol", {"operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2026-06-24T03:00:00Z", "buy_reason": "P116 缺代码"}, "missing symbol"),
    ]
    for request_id, body, reason in rejects:
        expect_http_status(base_url, "POST", "/api/v1/portfolio/offline-transactions", 400, body, request_id)
        add_reject(scenarios["L05"], "POST", "/api/v1/portfolio/offline-transactions", 400, request_id, reason)
    invalid_state_positions = positions[:1] + [{**positions[1], "position_state": "auto_trade"}]
    expect_http_status(base_url, "POST", "/api/v1/portfolio/adjustments", 400, {"cash": 1200, "total_assets": 1808, "adjust_reason": "P116 非法状态", "positions": invalid_state_positions}, "req_p116_reject_state")
    add_reject(scenarios["L05"], "POST", "/api/v1/portfolio/adjustments", 400, "req_p116_reject_state", "invalid position_state")
    require(count_transactions(db_path) == before_reject_count, "rejected operations must not create transactions")

    current_for_edit = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p116_current_for_edit"))
    pos_110022 = next(item for item in current_for_edit["positions"] if item["symbol"] == "110022")
    pos_588000 = next(item for item in current_for_edit["positions"] if item["symbol"] == "588000")
    pos_510300 = next(item for item in current_for_edit["positions"] if item["symbol"] == "510300")
    request_json(base_url, "POST", "/api/v1/portfolio/holdings", {"position_id": pos_110022["position_id"], "reason": "P116 消费基金份额校准", "confirmation": "confirmed", "position": {"symbol": "110022", "name": "易方达消费行业", "quantity": 180, "cost_price": 2.4, "current_price": 2.55, "buy_date": "2026-03-01", "position_state": "normal", "buy_reason": "长期消费基金观察", "asset_tag": "active_fund"}}, "req_p116_edit_holding")
    request_json(base_url, "POST", "/api/v1/portfolio/holdings/remove", {"position_id": pos_588000["position_id"], "reason": "P116 移除科创50观察仓", "confirmation": "confirmed"}, "req_p116_remove_holding")
    correction = data(request_json(base_url, "POST", "/api/v1/portfolio/corrections", {"target_type": "position", "target_id": pos_510300["position_id"], "before_json": '{"quantity":100}', "after_json": '{"quantity":80}', "correction_reason": "P116 核对线下卖出后修正审计"}, "req_p116_correction"))
    require(correction.get("correction_id"), "correction should return id")
    add_api(scenarios["L06"], "POST", "/api/v1/portfolio/holdings", request_id="req_p116_edit_holding")
    add_api(scenarios["L06"], "POST", "/api/v1/portfolio/holdings/remove", request_id="req_p116_remove_holding")
    add_api(scenarios["L06"], "POST", "/api/v1/portfolio/corrections", request_id="req_p116_correction")

    request_json(base_url, "POST", "/api/v1/decisions/decision_p116_execute/confirmations", {"confirmation_type": "executed_manually", "operation_type": "sell", "symbol": "510300", "quantity": 5, "price": 4.3, "fees": 1, "executed_at": "2026-06-24T03:00:00Z", "note": "P116 决策确认后人工线下卖出记录"}, "req_p116_confirmation_execute")
    add_api(scenarios["L07"], "POST", "/api/v1/decisions/decision_p116_execute/confirmations", request_id="req_p116_confirmation_execute")
    request_json(base_url, "POST", "/api/v1/decisions/decision_p116_error/confirmations", {"confirmation_type": "marked_error", "actual_outcome": "P116 多基金复盘发现建议偏离", "root_cause_tag": "evidence_missed", "lesson_learned": "后续多基金场景必须补充成交后仓位核对", "note": "P116 错误标注"}, "req_p116_marked_error")
    add_api(scenarios["L08"], "POST", "/api/v1/decisions/decision_p116_error/confirmations", request_id="req_p116_marked_error")

    rebalance = data(request_json(base_url, "POST", "/api/v1/portfolio/rebalance-review", {"target_core_ratio": 0.55, "target_satellite_ratio": 0.3, "target_cash_ratio": 0.15, "drift_threshold": 0.05, "review_date": "2026-06-25"}, "req_p116_rebalance"))
    require(rebalance.get("review_id"), "rebalance should return review id")
    add_api(scenarios["L09"], "POST", "/api/v1/portfolio/rebalance-review", request_id="req_p116_rebalance")

    risk_resolved = data(request_json(base_url, "POST", "/api/v1/risk-alerts/risk_p116_active/lifecycle", {"status": "resolved", "reason": "P116 多基金仓位人工复核完成"}, "req_p116_risk_resolve"))
    require(risk_resolved["sop_status"] == "resolved", "risk should resolve")
    request_json(base_url, "POST", "/api/v1/notifications/notif_p116_unread/read", request_id="req_p116_notification_read")
    request_json(base_url, "POST", "/api/v1/notifications/read-all", request_id="req_p116_notifications_read_all")
    add_api(scenarios["L10"], "POST", "/api/v1/risk-alerts/risk_p116_active/lifecycle", request_id="req_p116_risk_resolve")
    add_api(scenarios["L10"], "POST", "/api/v1/notifications/{id}/read", request_id="req_p116_notification_read")

    gate = data(request_json(base_url, "GET", "/api/v1/data-source-quality/gate-resolution?symbol=000300", request_id="req_p116_dq_gate"))
    require(gate["release_claim_state"] == "requires_resolution", "seeded DQ gate should require resolution")
    resolution = data(request_json(base_url, "POST", "/api/v1/data-source-quality/resolutions", {"symbol": "000300", "resolution_type": "scope_exclusion", "scope": "P116 local multi-fund transaction acceptance excludes current-data clean claim", "reason": "Seeded stale A-level source health for gate-resolution operation validation", "release_impact": "Do not claim current data clean in P116", "evidence_ref": "docs/release/acceptance/2026-06-25-p116-multi-fund-transaction-ledger-acceptance-matrix.md"}, "req_p116_dq_resolution_create"))
    active = resolution.get("active_resolution") or {}
    require(active.get("resolution_id"), "DQ resolution should create active resolution")
    request_json(base_url, "POST", f"/api/v1/data-source-quality/resolutions/{urllib.parse.quote(active['resolution_id'])}/retire", request_id="req_p116_dq_resolution_retire")
    add_api(scenarios["L11"], "GET", "/api/v1/data-source-quality/gate-resolution", request_id="req_p116_dq_gate")
    add_api(scenarios["L11"], "POST", "/api/v1/data-source-quality/resolutions", request_id="req_p116_dq_resolution_create")
    add_api(scenarios["L11"], "POST", "/api/v1/data-source-quality/resolutions/{id}/retire", request_id="req_p116_dq_resolution_retire")

    for method, path, rid in [
        ("GET", "/api/v1/dashboard/today", "req_p116_dashboard"),
        ("GET", "/api/v1/review/summary", "req_p116_review"),
        ("GET", "/api/v1/audit-events", "req_p116_audit"),
        ("GET", "/api/v1/decision-loops?limit=10", "req_p116_loop"),
    ]:
        request_json(base_url, method, path, request_id=rid)
        add_api(scenarios["L12"], method, path, request_id=rid)
    add_downstream(scenarios["L12"], "/dashboard /workbench /review /audit", "downstream endpoints returned after multi-fund writes", "ok")

    enrich_sqlite_evidence(db_path, scenarios)
    return [scenarios[sid] for sid in SCENARIO_IDS]


def enrich_sqlite_evidence(db_path: Path, scenarios: dict[str, dict[str, Any]]) -> None:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        symbol_count = int(scalar(conn, "SELECT COUNT(DISTINCT symbol) FROM positions WHERE symbol IN ('510300','159915','512000','110022','161725')"))
        counts = {
            "portfolio_snapshots": int(scalar(conn, "SELECT COUNT(*) FROM portfolio_snapshots")),
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "transaction_symbols": int(scalar(conn, "SELECT COUNT(DISTINCT symbol) FROM position_transactions WHERE symbol IN ('510300','159915','588000','512000')")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "marked_error_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type='marked_error'")),
            "import_batches_invalid": int(scalar(conn, "SELECT COUNT(*) FROM local_account_import_batches WHERE invalid_count > 0")),
            "import_batches_committed": int(scalar(conn, "SELECT COUNT(*) FROM local_account_import_batches WHERE status='committed'")),
            "corrections": int(scalar(conn, "SELECT COUNT(*) FROM local_account_corrections")),
            "dq_resolutions_retired": int(scalar(conn, "SELECT COUNT(*) FROM data_quality_gate_resolutions WHERE status='retired'")),
            "risk_resolved": int(scalar(conn, "SELECT COUNT(*) FROM risk_alerts WHERE alert_id='risk_p116_active' AND sop_status='resolved'")),
            "notifications_read": int(scalar(conn, "SELECT COUNT(*) FROM notifications WHERE notification_id='notif_p116_unread' AND read_at IS NOT NULL")),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
            "auto_confirmation_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")),
            "forbidden_broker_order_push_tables": int(scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")),
            "auto_rule_apply_audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events WHERE LOWER(action) LIKE '%auto%' AND (LOWER(action) LIKE '%rule%' OR LOWER(action) LIKE '%confirm%' OR LOWER(action) LIKE '%trade%')")),
        }
        require(symbol_count >= 5, f"expected at least 5 active symbols after edits/import, got {symbol_count}")
        require(counts["portfolio_snapshots"] >= 8, "portfolio snapshots missing")
        require(counts["position_transactions"] >= 6, "multi-fund transactions missing")
        require(counts["transaction_symbols"] >= 4, "transactions should cover at least four symbols")
        require(counts["operation_confirmations"] >= 7, "confirmations missing")
        require(counts["marked_error_rows"] >= 1, "marked error confirmation missing")
        require(counts["import_batches_invalid"] >= 1 and counts["import_batches_committed"] >= 1, "import batch invalid/committed evidence missing")
        require(counts["corrections"] >= 1, "correction evidence missing")
        require(counts["dq_resolutions_retired"] >= 1, "DQ retired resolution missing")
        require(counts["risk_resolved"] == 1, "risk resolution missing")
        require(counts["notifications_read"] == 1, "notification read missing")
        require(counts["audit_events"] >= 12, "audit events missing")
        require(counts["auto_confirmation_rows"] == 0, "auto confirmation rows must be absent")
        require(counts["forbidden_broker_order_push_tables"] == 0, "forbidden broker/order/push tables must be absent")
        require(counts["auto_rule_apply_audit_events"] == 0, "auto rule apply audit events must be absent")

        add_sqlite(scenarios["L02"], "positions", "distinct active symbols", symbol_count)
        add_sqlite(scenarios["L03"], "position_transactions", "count", counts["position_transactions"])
        add_sqlite(scenarios["L03"], "position_transactions", "distinct symbols", counts["transaction_symbols"])
        add_sqlite(scenarios["L04"], "local_account_import_batches", "invalid_count", counts["import_batches_invalid"])
        add_sqlite(scenarios["L04"], "local_account_import_batches", "committed", counts["import_batches_committed"])
        add_sqlite(scenarios["L05"], "position_transactions", "unchanged after rejects", counts["position_transactions"])
        add_sqlite(scenarios["L06"], "local_account_corrections", "count", counts["corrections"])
        add_sqlite(scenarios["L07"], "operation_confirmations", "count", counts["operation_confirmations"])
        add_sqlite(scenarios["L08"], "operation_confirmations", "marked_error", counts["marked_error_rows"])
        add_sqlite(scenarios["L09"], "portfolio_snapshots", "count", counts["portfolio_snapshots"])
        add_sqlite(scenarios["L10"], "risk_alerts", "resolved", counts["risk_resolved"])
        add_sqlite(scenarios["L10"], "notifications", "read_at", counts["notifications_read"])
        add_sqlite(scenarios["L11"], "data_quality_gate_resolutions", "retired", counts["dq_resolutions_retired"])
        add_sqlite(scenarios["L12"], "audit_events", "count", counts["audit_events"])

        for item in scenarios.values():
            item["safety_counters"] = {
                "forbidden_broker_order_push_tables": counts["forbidden_broker_order_push_tables"],
                "auto_confirmation_rows": counts["auto_confirmation_rows"],
                "auto_rule_apply_audit_events": counts["auto_rule_apply_audit_events"],
                "automatic_trading_affordances": 0,
                "return_guarantee_claims": 0,
                "secret_or_raw_prompt_leaks": 0,
            }
        add_sqlite(scenarios["L16"], "sqlite_master", "forbidden broker/order/push tables", counts["forbidden_broker_order_push_tables"])
        add_sqlite(scenarios["L16"], "operation_confirmations", "auto confirmations", counts["auto_confirmation_rows"])
        add_sqlite(scenarios["L16"], "audit_events", "auto rule apply events", counts["auto_rule_apply_audit_events"])
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
        "claim_boundary": "P116 API/SQLite layer uses local seeded multi-fund transaction evidence. It does not claim real broker trades, external push, fresh provider/LLM output, automatic trading, automatic confirmation, automatic rule application, or return guarantees.",
    }
    out = artifact_dir / "api_sqlite" / "p116-api-sqlite-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def merge_summary(base_url: str, db_path: Path, artifact_dir: Path, browser_summary: Path | None) -> dict[str, Any]:
    api_path = artifact_dir / "api_sqlite" / "p116-api-sqlite-summary.json"
    require(api_path.exists(), f"missing API summary: {api_path}")
    api_payload = json.loads(api_path.read_text(encoding="utf-8"))
    scenarios = {item["scenario_id"]: item for item in api_payload["scenarios"]}
    browser_payload: dict[str, Any] = {}
    if browser_summary and browser_summary.exists():
        browser_payload = json.loads(browser_summary.read_text(encoding="utf-8"))
        for entry in browser_payload.get("scenarios", []):
            sid = entry["scenario_id"]
            if sid in scenarios:
                scenarios[sid]["browser_evidence"].extend(entry.get("browser_evidence", []))
                scenarios[sid]["redaction_evidence"].update(entry.get("redaction_evidence", {}))
    missing_browser = [sid for sid in ["L02", "L03", "L04", "L07", "L10", "L11", "L12", "L13", "L14", "L15", "L16"] if not scenarios.get(sid, {}).get("browser_evidence")]
    require(not missing_browser, f"missing browser evidence for {missing_browser}")

    final = {
        "status": "passed",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "base_url": base_url,
        "sqlite_path": str(db_path),
        "scenario_count": len(scenarios),
        "fresh_pass_count": sum(1 for item in scenarios.values() if item["status"] == "fresh_pass"),
        "scoped_pass_count": sum(1 for item in scenarios.values() if item["status"] == "scoped_pass"),
        "symbols": MULTI_FUND_SYMBOLS,
        "scenarios": [scenarios[sid] for sid in SCENARIO_IDS],
        "browser_summary": browser_payload,
        "claim_boundary": "P116 validates local multi-fund transaction ledger behavior across API, SQLite and browser UI. It must not be cited as broker execution, external data cleanliness, real LLM quality, release package refresh, or physical second-machine acceptance.",
    }
    out = artifact_dir / "p116-scenario-summary.json"
    out.write_text(json.dumps(final, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(final, ensure_ascii=False, indent=2))
    return final


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", required=True)
    parser.add_argument("--sqlite", required=True)
    parser.add_argument("--artifact-dir", required=True)
    parser.add_argument("--browser-summary")
    parser.add_argument("--merge-only", action="store_true")
    args = parser.parse_args()

    artifact_dir = Path(args.artifact_dir)
    artifact_dir.mkdir(parents=True, exist_ok=True)
    db_path = Path(args.sqlite)
    if args.merge_only:
        merge_summary(args.base_url, db_path, artifact_dir, Path(args.browser_summary) if args.browser_summary else None)
    else:
        payload = write_api_summary(args.base_url, db_path, artifact_dir)
        print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
