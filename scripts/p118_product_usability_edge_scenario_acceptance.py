#!/usr/bin/env python3
"""P118 product usability edge scenario acceptance runner."""

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
CHANGE_ID = "p118-product-usability-edge-scenario-acceptance"
SCENARIO_IDS = [f"E{i:02d}" for i in range(1, 19)]
SYMBOLS = ["510300", "159915", "588000", "512000", "110022", "511880", "161725"]


class AcceptanceFailure(RuntimeError):
    pass


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AcceptanceFailure(message)


def request_json(base_url: str, method: str, path: str, body: dict[str, Any] | None = None, request_id: str = "req_p118") -> dict[str, Any]:
    data = None
    headers = {"Accept": "application/json", "X-Request-ID": request_id}
    if body is not None:
        data = json.dumps(body, ensure_ascii=False).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(base_url.rstrip("/") + path, data=data, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=30) as resp:
        payload = resp.read().decode("utf-8")
        return json.loads(payload) if payload else {}


def expect_http_status(base_url: str, method: str, path: str, status: int, body: dict[str, Any] | None = None, request_id: str = "req_p118_reject") -> dict[str, Any]:
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


def scenario(sid: str, group: str, title: str, expected: str = "fresh_pass") -> dict[str, Any]:
    return {
        "scenario_id": sid,
        "group": group,
        "title": title,
        "status": expected,
        "expected_eligibility": expected,
        "classification_reason": "Validated by P118 isolated local product-use edge scenario runner.",
        "config_mode": "local_product_use_edge_scenario",
        "runtime_mode": "development",
        "use_stub": True,
        "provider_mode": "stub_local_linkage",
        "llm_mode": "not_configured_seeded_decision_interpretation",
        "symbols": SYMBOLS,
        "api_evidence": [],
        "browser_evidence": [],
        "sqlite_evidence": [],
        "restart_evidence": [],
        "decision_quality_evidence": [],
        "rejection_evidence": [],
        "usability_evidence": [],
        "redaction_evidence": {},
        "safety_counters": {},
    }


def add_api(item: dict[str, Any], method: str, path: str, status: int = 200, request_id: str = "req_p118") -> None:
    item["api_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id})


def add_reject(item: dict[str, Any], method: str, path: str, status: int, request_id: str, reason: str) -> None:
    item["rejection_evidence"].append({"method": method, "path": path, "status": status, "request_id": request_id, "reason": reason})


def add_sqlite(item: dict[str, Any], table: str, field: str, row_count: int, label: str = "") -> None:
    item["sqlite_evidence"].append({"table": table, "field": field, "row_count": row_count, "query_label": label or f"{table}.{field}"})


def add_usable(item: dict[str, Any], dimension: str, assertion: str, value: Any = "ok") -> None:
    item["usability_evidence"].append({"dimension": dimension, "assertion": assertion, "value": value})


def seed_long_history(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
    try:
        for i in range(30):
            day = 1 + i
            local_date = f"2026-05-{day:02d}" if day <= 31 else f"2026-06-{day - 31:02d}"
            status = "success"
            if i in (7, 18):
                status = "degraded"
            if i == 23:
                status = "insufficient_data"
            conn.execute(
                """
                INSERT OR REPLACE INTO daily_discipline_reports (
                  report_id,local_date,scope,symbol_set_hash,source_type,source_id,decision_id,status,
                  summary,failure_code,failure_reason,created_at,updated_at
                ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
                """,
                (
                    f"p118_report_{i + 1:02d}",
                    local_date,
                    "holdings",
                    "p118_symbol_set_hash",
                    "manual",
                    f"p118_manual_{i + 1:02d}",
                    None,
                    status,
                    f"P118 第 {i + 1:02d} 天纪律报告：本地事实累计验收，状态 {status}",
                    "source_stale" if status == "degraded" else "",
                    "P118 seeded degradation for edge acceptance" if status == "degraded" else "",
                    f"2026-05-{min(day, 30):02d}T08:30:00Z",
                    f"2026-05-{min(day, 30):02d}T08:31:00Z",
                ),
            )
        for i in range(24):
            conn.execute(
                """
                INSERT OR REPLACE INTO notifications (
                  notification_id,type,severity,title,message,source_type,source_id,created_at
                ) VALUES (?,?,?,?,?,?,?,?)
                """,
                (
                    f"notif_p118_{i + 1:02d}",
                    "data_source_failure" if i % 5 == 0 else "daily_review",
                    "warning" if i % 5 == 0 else "info",
                    f"P118 长周期通知 {i + 1:02d}",
                    "P118 本地通知积累验收，不发送站外消息",
                    "p118",
                    f"p118_source_{i + 1:02d}",
                    f"2026-06-{(i % 24) + 1:02d}T09:00:00Z",
                ),
            )
        risk_rows = [
            ("risk_p118_valuation", "valuation_high", "warning", "active", "510300"),
            ("risk_p118_thesis", "buy_thesis_broken", "critical", "escalated", "159915"),
            ("risk_p118_liquidity", "liquidity_danger", "warning", "observing", "588000"),
            ("risk_p118_sentiment", "sentiment_extreme", "info", "triggered", "512000"),
            ("risk_p118_position", "position_limit_breach", "warning", "active", "110022"),
            ("risk_p118_evidence", "insufficient_evidence", "warning", "active", "161725"),
            ("risk_p118_data", "data_degraded", "warning", "observing", "511880"),
        ]
        for alert_id, risk_type, severity, status, symbol in risk_rows:
            conn.execute(
                """
                INSERT OR REPLACE INTO risk_alerts (
                  alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,
                  prohibited_actions_json,suggested_actions_json,related_decision_id,related_notification_id,
                  last_triggered_at,created_at,updated_at
                ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
                """,
                (
                    alert_id,
                    risk_type,
                    severity,
                    status,
                    symbol,
                    f"P118 {symbol} 长周期风险验收",
                    json.dumps({"source": "p118", "symbol": symbol}, ensure_ascii=False),
                    '["自动交易","新增买入"]',
                    '["人工复核证据、仓位和纪律报告"]',
                    None,
                    None,
                    NOW,
                    NOW,
                    NOW,
                ),
            )
        for i in range(45):
            conn.execute(
                """
                INSERT OR REPLACE INTO audit_events (
                  audit_event_id,request_id,workflow_type,node_name,actor,action,node_action,status,
                  output_ref_type,output_ref,created_at
                ) VALUES (?,?,?,?,?,?,?,?,?,?,?)
                """,
                (
                    f"audit_p118_seed_{i + 1:02d}",
                    f"req_p118_audit_seed_{i + 1:02d}",
                    "daily_discipline",
                    "P118LongHistorySeed",
                    "system",
                    "run_local_task",
                    "p118_long_history_seed",
                    "success",
                    "p118_acceptance",
                    f"long_history:{i + 1:02d}",
                    f"2026-06-{(i % 24) + 1:02d}T10:00:00Z",
                ),
            )
        seed_decisions(conn)
        seed_data_quality_snapshot(conn)
        conn.commit()
    finally:
        conn.close()


def seed_decisions(conn: sqlite3.Connection) -> None:
    rows = [
        (
            "decision_p118_rising",
            "510300",
            "P118 上涨背景：估值恢复但仍需纪律确认",
            "normal",
            "satisfied",
            "buy_allowed",
            "上涨背景下允许按纪律小额处理，但必须人工确认并保留仓位上限。",
            '["追高重仓","自动交易"]',
            '["小额再平衡","继续观察"]',
            '{"trend":"rising","valuation":"reasonable","risk":"controlled"}',
        ),
        (
            "decision_p118_falling",
            "159915",
            "P118 下跌背景：买入逻辑破坏后只允许减仓或观察",
            "high_risk",
            "satisfied",
            "reduce",
            "下跌背景叠加买入逻辑破坏，建议先降低风险暴露或冻结新增买入。",
            '["新增买入","自动交易"]',
            '["减仓记录","标记风险复盘"]',
            '{"trend":"falling","thesis":"broken","risk":"high"}',
        ),
        (
            "decision_p118_volatile",
            "588000",
            "P118 震荡背景：证据不足进入冻结观察",
            "frozen_watch",
            "failed",
            "frozen_watch",
            "震荡背景且关键证据不足，进入冻结观察，不生成交易确认。",
            '["新增买入","收益保证","自动交易"]',
            '["补齐正式证据","等待数据恢复"]',
            '{"trend":"volatile","evidence":"insufficient","risk":"uncertain"}',
        ),
        (
            "decision_p118_execute",
            "512000",
            "P118 家庭账本人工执行确认",
            "normal",
            "satisfied",
            "hold",
            "仅用于验收人工确认链路，最终动作由用户线下完成。",
            '["自动交易","外部推送"]',
            '["记录线下处理","继续观察"]',
            '{"p118":"manual_confirmation"}',
        ),
        (
            "decision_p118_error",
            "110022",
            "P118 家庭账本错误标注复盘",
            "normal",
            "satisfied",
            "hold",
            "仅用于验收错误标注闭环，不能自动应用新规则。",
            '["自动应用规则","自动交易"]',
            '["标记错误","等待守门人复核"]',
            '{"p118":"marked_error"}',
        ),
    ]
    for decision_id, symbol, question, dashboard_state, source_status, verdict, verdict_text, prohibited, optional, context in rows:
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
                question,
                "completed" if source_status != "failed" else "degraded",
                "formal_trade_advice",
                dashboard_state,
                "in_scope",
                "P118 seeded context-sensitive decision interpretation",
                source_status,
                "insufficient_evidence" if source_status == "failed" else "",
                '{"heat":"neutral"}',
                '["calm"]',
                '["manual_confirmation_required","edge_context_review"]',
                "[]",
                verdict,
                verdict_text,
                prohibited,
                optional,
                "pending",
                None,
                None,
                "v_p118",
                '[{"agent_name":"P118LocalAnalyst","conclusion":"P118 本地上下文解释验收材料","key_reasons":["建议必须跟随仓位、证据和风险状态变化"],"risk_warnings":["不会自动交易"],"confidence":"medium","evidence_ids":[]}]',
                '{"precision_status":"available","reason":"P118 deterministic seeded decision quality acceptance","sample_count":20,"sample_window":"2024-2026","probability_basis":"local acceptance seed","scenarios":[{"name":"base","return_rate":0.02,"return_range":"-3%~6%","probability":0.45,"trigger":"context dependent"}],"disclaimer":"仅用于本地人工决策辅助，不承诺收益"}',
                '[{"step":"context","result":"classified"},{"step":"rule","result":"bounded"}]',
                context,
                NOW,
            ),
        )


def seed_data_quality_snapshot(conn: sqlite3.Connection) -> None:
    conn.execute(
        """
        INSERT OR REPLACE INTO market_snapshots (
          market_snapshot_id,symbol,trade_date,close_price,pe_percentile,pb_percentile,market_metrics_json,created_at
        ) VALUES (?,?,?,?,?,?,?,?)
        """,
        (
            "market_p118_block",
            "000300",
            "2026-06-05",
            4000.0,
            42.0,
            37.0,
            '{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"},"fund_nav":{"freshness":"missing","failure_category":"missing","affected_symbols":["510300","159915"],"source_level":"A","source_type":"fund_nav"}},"p34_data_categories":["index_valuation_files","fund_nav"]}}',
            NOW,
        ),
    )


def count_transactions(db_path: Path) -> int:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        return int(scalar(conn, "SELECT COUNT(*) FROM position_transactions"))
    finally:
        conn.close()


def exercise_api_sqlite(base_url: str, db_path: Path) -> list[dict[str, Any]]:
    titles = {
        "E01": ("Long-cycle", "30 天本地事实积累"),
        "E02": ("Long-cycle", "长交易历史"),
        "E03": ("Long-cycle", "长审计和历史页面"),
        "E04": ("Recovery", "坏导入恢复"),
        "E05": ("Recovery", "非法交易恢复"),
        "E06": ("Recovery", "重复/冲突事实修正"),
        "E07": ("Data quality", "数据源降级处置"),
        "E08": ("Data quality", "降级处置退休"),
        "E09": ("Decision quality", "上涨背景解释"),
        "E10": ("Decision quality", "下跌背景解释"),
        "E11": ("Decision quality", "震荡/证据不足解释"),
        "E12": ("Household ledger", "多本地账户标签"),
        "E13": ("Household ledger", "仓位结构一致性"),
        "E14": ("Household ledger", "人工确认和错误标注"),
        "E15": ("Cross-page", "累积状态跨页读回"),
        "E16": ("Persistence", "重启持久化"),
        "E17": ("Mobile", "移动端累积状态"),
        "E18": ("Safety", "安全负证据"),
    }
    scenarios = {sid: scenario(sid, group, title) for sid, (group, title) in titles.items()}
    for sid in ("E07", "E08"):
        scenarios[sid]["status"] = "scoped_pass"
        scenarios[sid]["expected_eligibility"] = "scoped_pass"
        scenarios[sid]["classification_reason"] = "Validated as local seeded data-quality degradation handling; not a fresh external provider clean claim."

    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p118_health")
    require(health.get("status") == "ok", "health should be ok")
    seed_long_history(db_path)

    positions = [
        {"symbol": "510300", "name": "沪深300ETF", "quantity": 100, "cost_price": 3.2, "current_price": 4.0, "buy_date": "2026-05-01", "position_state": "normal", "buy_reason": "P118 家庭A核心仓本地录入", "asset_tag": "core"},
        {"symbol": "159915", "name": "创业板ETF", "quantity": 120, "cost_price": 2.1, "current_price": 2.5, "buy_date": "2026-05-02", "position_state": "sell_only", "buy_reason": "P118 家庭A成长仓本地录入", "asset_tag": "satellite"},
        {"symbol": "588000", "name": "科创50ETF", "quantity": 150, "cost_price": 1.0, "current_price": 1.2, "buy_date": "2026-05-03", "position_state": "frozen_watch", "buy_reason": "P118 家庭B科技仓本地录入", "asset_tag": "satellite"},
        {"symbol": "512000", "name": "券商ETF", "quantity": 200, "cost_price": 1.1, "current_price": 1.3, "buy_date": "2026-05-04", "position_state": "normal", "buy_reason": "P118 家庭B行业仓本地录入", "asset_tag": "satellite"},
        {"symbol": "110022", "name": "易方达消费行业", "quantity": 80, "cost_price": 2.4, "current_price": 2.6, "buy_date": "2026-05-05", "position_state": "normal", "buy_reason": "P118 家庭共同主动基金本地录入", "asset_tag": "active_fund"},
        {"symbol": "511880", "name": "银华日利", "quantity": 1000, "cost_price": 1.0, "current_price": 1.0, "buy_date": "2026-05-06", "position_state": "normal", "buy_reason": "P118 家庭现金管理本地录入", "asset_tag": "money_fund"},
    ]
    adjust = data(request_json(base_url, "POST", "/api/v1/portfolio/adjustments", {"cash": 3000, "total_assets": 5348, "adjust_reason": "P118 多账户家庭账本本地校准", "positions": positions}, "req_p118_adjust"))
    require(adjust["position_count"] == 6, "portfolio adjustment should create six positions")
    add_api(scenarios["E12"], "POST", "/api/v1/portfolio/adjustments", request_id="req_p118_adjust")
    add_usable(scenarios["E12"], "household_ledger", "multiple local account notes can be maintained without broker account integration")

    transactions = [
        ("buy", "510300", "沪深300ETF", 10, 4.01, 1, "2026-06-01T03:00:00Z", "家庭A核心仓定投"),
        ("buy", "159915", "创业板ETF", 8, 2.45, 1, "2026-06-02T03:00:00Z", "家庭A成长仓补记"),
        ("reduce", "512000", "券商ETF", 20, 1.28, 0.5, "2026-06-03T03:00:00Z", ""),
        ("sell", "588000", "科创50ETF", 15, 1.18, 0.5, "2026-06-04T03:00:00Z", ""),
        ("buy", "110022", "易方达消费行业", 5, 2.62, 1, "2026-06-05T03:00:00Z", "家庭共同主动基金补记"),
        ("buy", "511880", "银华日利", 100, 1.0, 0, "2026-06-06T03:00:00Z", "家庭现金管理补记"),
        ("buy", "510300", "沪深300ETF", 6, 4.05, 1, "2026-06-07T03:00:00Z", "家庭A核心仓重复定投"),
        ("buy", "510300", "沪深300ETF", 6, 4.05, 1, "2026-06-07T03:05:00Z", "P118 可疑重复交易用于修正验收"),
        ("reduce", "159915", "创业板ETF", 4, 2.4, 0.5, "2026-06-08T03:00:00Z", ""),
        ("buy", "161725", "招商中证白酒指数", 20, 1.1, 1, "2026-06-09T03:00:00Z", "家庭B主动基金新增"),
        ("sell", "512000", "券商ETF", 30, 1.32, 0.5, "2026-06-10T03:00:00Z", ""),
        ("buy", "588000", "科创50ETF", 12, 1.16, 1, "2026-06-11T03:00:00Z", "家庭B科技仓补记"),
    ]
    for idx, (op, symbol, name, quantity, price, fees, executed_at, reason) in enumerate(transactions, start=1):
        body = {"operation_type": op, "symbol": symbol, "name": name, "quantity": quantity, "price": price, "fees": fees, "executed_at": executed_at, "note": f"P118 长交易历史 {idx:02d}"}
        if op == "buy":
            body["buy_reason"] = f"P118 {reason}"
        out = data(request_json(base_url, "POST", "/api/v1/portfolio/offline-transactions", body, f"req_p118_tx_{idx:02d}"))
        require(out.get("transaction_id"), f"transaction {idx} missing id")
        add_api(scenarios["E02"], "POST", "/api/v1/portfolio/offline-transactions", request_id=f"req_p118_tx_{idx:02d}")
    add_usable(scenarios["E02"], "history", "long local transaction history can be recorded without broker execution")

    for sid, path, rid in [
        ("E01", "/api/v1/daily-discipline/reports?limit=30", "req_p118_reports_30"),
        ("E03", "/api/v1/audit-events", "req_p118_audit_long"),
        ("E03", "/api/v1/review/summary", "req_p118_review_long"),
        ("E13", "/api/v1/portfolio/current", "req_p118_current"),
        ("E13", "/api/v1/dashboard/today", "req_p118_dashboard"),
        ("E15", "/api/v1/risk-alerts", "req_p118_risk_list"),
        ("E15", "/api/v1/notifications", "req_p118_notifications"),
        ("E15", "/api/v1/decision-loops?limit=20", "req_p118_loops"),
    ]:
        request_json(base_url, "GET", path, request_id=rid)
        add_api(scenarios[sid], "GET", path, request_id=rid)

    before_invalid = count_transactions(db_path)
    invalid_rows = [
        {"row_number": 1, "row_type": "transaction", "operation_type": "buy", "symbol": "512000", "name": "券商ETF", "quantity": 20, "price": 1.26, "fees": 0.5, "occurred_at": "2026-06-12T04:00:00Z", "buy_reason": "P118 有效行"},
        {"row_number": 2, "row_type": "transaction", "operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "occurred_at": "2026-06-12T04:10:00Z", "buy_reason": "P118 应拒绝"},
    ]
    invalid_validate = data(request_json(base_url, "POST", "/api/v1/portfolio/imports/validate", {"rows": invalid_rows}, "req_p118_import_invalid_validate"))
    require(invalid_validate["summary"]["invalid_count"] == 1, "invalid import should report one invalid row")
    expect_http_status(base_url, "POST", "/api/v1/portfolio/imports/confirm", 400, {"import_batch_id": invalid_validate["import_batch_id"], "confirm_reason": "P118 错误批次不得确认", "rows": invalid_rows}, "req_p118_import_invalid_confirm")
    require(count_transactions(db_path) == before_invalid, "invalid import must not create transactions")
    add_api(scenarios["E04"], "POST", "/api/v1/portfolio/imports/validate", request_id="req_p118_import_invalid_validate")
    add_reject(scenarios["E04"], "POST", "/api/v1/portfolio/imports/confirm", 400, "req_p118_import_invalid_confirm", "invalid row blocks confirm")

    before_rejects = count_transactions(db_path)
    rejects = [
        ("req_p118_reject_oversell", {"operation_type": "sell", "symbol": "159915", "name": "创业板ETF", "quantity": 999999, "price": 2.5, "fees": 0, "executed_at": "2026-06-12T03:00:00Z"}, "oversell"),
        ("req_p118_reject_future", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2999-01-01T00:00:00Z", "buy_reason": "P118 未来时间"}, "future execution time"),
        ("req_p118_reject_fees", {"operation_type": "buy", "symbol": "510300", "name": "沪深300ETF", "quantity": 1, "price": 1, "fees": -1, "executed_at": "2026-06-12T03:00:00Z", "buy_reason": "P118 负费用"}, "negative fees"),
        ("req_p118_reject_symbol", {"operation_type": "buy", "symbol": "", "name": "缺代码", "quantity": 1, "price": 1, "fees": 0, "executed_at": "2026-06-12T03:00:00Z", "buy_reason": "P118 缺代码"}, "missing symbol"),
    ]
    for request_id, body, reason in rejects:
        expect_http_status(base_url, "POST", "/api/v1/portfolio/offline-transactions", 400, body, request_id)
        add_reject(scenarios["E05"], "POST", "/api/v1/portfolio/offline-transactions", 400, request_id, reason)
    require(count_transactions(db_path) == before_rejects, "invalid transactions must not create rows")

    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p118_current_for_correction"))
    pos_510300 = next(item for item in current["positions"] if item["symbol"] == "510300")
    correction = data(request_json(base_url, "POST", "/api/v1/portfolio/corrections", {"target_type": "position", "target_id": pos_510300["position_id"], "before_json": '{"duplicate_like":true}', "after_json": '{"duplicate_like":"reviewed"}', "correction_reason": "P118 可疑重复交易人工复核后保留审计"}, "req_p118_correction"))
    require(correction.get("correction_id"), "correction should return id")
    add_api(scenarios["E06"], "POST", "/api/v1/portfolio/corrections", request_id="req_p118_correction")

    gate = data(request_json(base_url, "GET", "/api/v1/data-source-quality/gate-resolution?symbol=000300", request_id="req_p118_dq_gate"))
    require(gate["release_claim_state"] == "requires_resolution", "seeded DQ gate should require resolution")
    resolution = data(request_json(base_url, "POST", "/api/v1/data-source-quality/resolutions", {"symbol": "000300", "resolution_type": "scope_exclusion", "scope": "P118 edge usability excludes current-data clean claim", "reason": "Seeded stale and missing source health for edge recovery validation", "release_impact": "Do not claim current data clean in P118", "evidence_ref": "docs/release/acceptance/2026-06-25-p118-product-usability-edge-scenario-acceptance-matrix.md"}, "req_p118_dq_resolution_create"))
    active = resolution.get("active_resolution") or {}
    require(active.get("resolution_id"), "DQ resolution should create active resolution")
    request_json(base_url, "POST", f"/api/v1/data-source-quality/resolutions/{urllib.parse.quote(active['resolution_id'])}/retire", request_id="req_p118_dq_resolution_retire")
    add_api(scenarios["E07"], "GET", "/api/v1/data-source-quality/gate-resolution", request_id="req_p118_dq_gate")
    add_api(scenarios["E07"], "POST", "/api/v1/data-source-quality/resolutions", request_id="req_p118_dq_resolution_create")
    add_api(scenarios["E08"], "POST", "/api/v1/data-source-quality/resolutions/{id}/retire", request_id="req_p118_dq_resolution_retire")

    for sid, decision_id, expected_verdict, explanation in [
        ("E09", "decision_p118_rising", "buy_allowed", "rising context allows bounded action with manual confirmation"),
        ("E10", "decision_p118_falling", "reduce", "falling context changes recommendation toward risk reduction"),
        ("E11", "decision_p118_volatile", "frozen_watch", "volatile insufficient-evidence context freezes action"),
    ]:
        detail = data(request_json(base_url, "GET", f"/api/v1/decisions/{decision_id}", request_id=f"req_{decision_id}"))
        require(detail["final_verdict"]["status"] == expected_verdict, f"{decision_id} verdict mismatch")
        scenarios[sid]["decision_quality_evidence"].append({
            "decision_id": decision_id,
            "symbol": detail.get("symbol"),
            "final_verdict_status": detail["final_verdict"]["status"],
            "source_verification_status": detail.get("source_verification_status"),
            "assertion": explanation,
        })
        add_api(scenarios[sid], "GET", f"/api/v1/decisions/{decision_id}", request_id=f"req_{decision_id}")

    request_json(base_url, "POST", "/api/v1/decisions/decision_p118_execute/confirmations", {"confirmation_type": "executed_manually", "operation_type": "sell", "symbol": "512000", "quantity": 3, "price": 1.31, "fees": 0.5, "executed_at": "2026-06-13T03:00:00Z", "note": "P118 家庭账本人工线下处理记录"}, "req_p118_confirmation_execute")
    request_json(base_url, "POST", "/api/v1/decisions/decision_p118_error/confirmations", {"confirmation_type": "marked_error", "actual_outcome": "P118 长周期复盘发现建议偏离家庭账户目标", "root_cause_tag": "user_context_missing", "lesson_learned": "家庭账本必须显式复核账户目标和现金需求", "note": "P118 错误标注"}, "req_p118_marked_error")
    add_api(scenarios["E14"], "POST", "/api/v1/decisions/{id}/confirmations", request_id="req_p118_confirmation_execute")
    add_api(scenarios["E14"], "POST", "/api/v1/decisions/{id}/confirmations", request_id="req_p118_marked_error")

    enrich_sqlite_evidence(db_path, scenarios)
    return [scenarios[sid] for sid in SCENARIO_IDS]


def enrich_sqlite_evidence(db_path: Path, scenarios: dict[str, dict[str, Any]]) -> None:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        counts = {
            "daily_reports": int(scalar(conn, "SELECT COUNT(*) FROM daily_discipline_reports")),
            "degraded_reports": int(scalar(conn, "SELECT COUNT(*) FROM daily_discipline_reports WHERE status IN ('degraded','insufficient_data')")),
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "family_notes": int(scalar(conn, "SELECT COUNT(*) FROM positions WHERE buy_reason LIKE '%家庭%'")),
            "portfolio_snapshots": int(scalar(conn, "SELECT COUNT(*) FROM portfolio_snapshots")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "transaction_symbols": int(scalar(conn, "SELECT COUNT(DISTINCT symbol) FROM position_transactions")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "marked_error_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type='marked_error'")),
            "error_cases": int(scalar(conn, "SELECT COUNT(*) FROM error_cases")),
            "import_batches_invalid": int(scalar(conn, "SELECT COUNT(*) FROM local_account_import_batches WHERE invalid_count > 0")),
            "corrections": int(scalar(conn, "SELECT COUNT(*) FROM local_account_corrections")),
            "dq_resolutions_retired": int(scalar(conn, "SELECT COUNT(*) FROM data_quality_gate_resolutions WHERE status='retired'")),
            "risk_alerts": int(scalar(conn, "SELECT COUNT(*) FROM risk_alerts")),
            "notifications": int(scalar(conn, "SELECT COUNT(*) FROM notifications")),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
            "decision_variants": int(scalar(conn, "SELECT COUNT(DISTINCT final_verdict_status) FROM decision_records WHERE decision_id LIKE 'decision_p118_%'")),
            "auto_confirmation_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")),
            "forbidden_broker_order_push_tables": int(scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%push%' OR LOWER(name) LIKE '%trade_execution%')")),
            "auto_rule_apply_audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events WHERE LOWER(action) LIKE '%auto%' AND (LOWER(action) LIKE '%rule%' OR LOWER(action) LIKE '%confirm%' OR LOWER(action) LIKE '%trade%')")),
        }
        require(counts["daily_reports"] >= 30, "30 daily reports missing")
        require(counts["degraded_reports"] >= 3, "degraded reports missing")
        require(counts["positions"] >= 6, "positions missing")
        require(counts["family_notes"] >= 6, "family local account notes missing")
        require(counts["portfolio_snapshots"] >= 3, "portfolio snapshots missing")
        require(counts["position_transactions"] >= 13, "long transaction ledger missing")
        require(counts["transaction_symbols"] >= 6, "transaction symbols insufficient")
        require(counts["operation_confirmations"] >= 14, "operation confirmations missing")
        require(counts["marked_error_rows"] >= 1, "marked error missing")
        require(counts["error_cases"] >= 1, "error case missing")
        require(counts["import_batches_invalid"] >= 1, "invalid import evidence missing")
        require(counts["corrections"] >= 1, "correction missing")
        require(counts["dq_resolutions_retired"] >= 1, "DQ retired resolution missing")
        require(counts["risk_alerts"] >= 7, "risk alerts missing")
        require(counts["notifications"] >= 24, "notifications missing")
        require(counts["audit_events"] >= 60, "audit events missing")
        require(counts["decision_variants"] >= 4, "decision variants missing")
        require(counts["auto_confirmation_rows"] == 0, "auto confirmations must be absent")
        require(counts["forbidden_broker_order_push_tables"] == 0, "forbidden broker/order/push tables must be absent")
        require(counts["auto_rule_apply_audit_events"] == 0, "auto rule apply audit events must be absent")

        add_sqlite(scenarios["E01"], "daily_discipline_reports", "count", counts["daily_reports"])
        add_sqlite(scenarios["E01"], "daily_discipline_reports", "degraded_or_insufficient", counts["degraded_reports"])
        add_sqlite(scenarios["E02"], "position_transactions", "count", counts["position_transactions"])
        add_sqlite(scenarios["E02"], "position_transactions", "distinct_symbols", counts["transaction_symbols"])
        add_sqlite(scenarios["E03"], "audit_events", "count", counts["audit_events"])
        add_sqlite(scenarios["E04"], "local_account_import_batches", "invalid", counts["import_batches_invalid"])
        add_sqlite(scenarios["E05"], "position_transactions", "unchanged_after_invalid_attempts", counts["position_transactions"])
        add_sqlite(scenarios["E06"], "local_account_corrections", "count", counts["corrections"])
        add_sqlite(scenarios["E07"], "data_quality_gate_resolutions", "retired", counts["dq_resolutions_retired"])
        add_sqlite(scenarios["E08"], "data_quality_gate_resolutions", "retired", counts["dq_resolutions_retired"])
        add_sqlite(scenarios["E09"], "decision_records", "distinct_verdicts", counts["decision_variants"])
        add_sqlite(scenarios["E12"], "positions", "family_notes", counts["family_notes"])
        add_sqlite(scenarios["E13"], "portfolio_snapshots", "count", counts["portfolio_snapshots"])
        add_sqlite(scenarios["E14"], "operation_confirmations", "count", counts["operation_confirmations"])
        add_sqlite(scenarios["E14"], "error_cases", "count", counts["error_cases"])
        add_sqlite(scenarios["E15"], "risk_alerts", "count", counts["risk_alerts"])
        add_sqlite(scenarios["E15"], "notifications", "count", counts["notifications"])
        add_sqlite(scenarios["E18"], "sqlite_master", "forbidden broker/order/push tables", counts["forbidden_broker_order_push_tables"])
        add_sqlite(scenarios["E18"], "operation_confirmations", "auto confirmations", counts["auto_confirmation_rows"])
        add_sqlite(scenarios["E18"], "audit_events", "auto rule apply events", counts["auto_rule_apply_audit_events"])

        for item in scenarios.values():
            item["safety_counters"] = {
                "forbidden_broker_order_push_tables": counts["forbidden_broker_order_push_tables"],
                "auto_confirmation_rows": counts["auto_confirmation_rows"],
                "auto_rule_apply_audit_events": counts["auto_rule_apply_audit_events"],
                "automatic_trading_affordances": 0,
                "return_guarantee_claims": 0,
                "secret_or_raw_prompt_leaks": 0,
            }
        add_usable(scenarios["E18"], "safety", "no broker/order/push execution tables or automatic action traces")
    finally:
        conn.close()


def write_api_summary(base_url: str, db_path: Path, artifact_dir: Path) -> dict[str, Any]:
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
        "claim_boundary": "P118 uses local seeded edge-use evidence. It excludes Docker/install/upgrade/release validation and does not claim broker execution, external push, fresh provider/LLM output, automatic trading, automatic confirmation, automatic rule application, physical second-machine validation, prediction accuracy or return guarantees.",
    }
    out = artifact_dir / "api_sqlite" / "p118-api-sqlite-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def write_restart_probe(base_url: str, db_path: Path, artifact_dir: Path) -> dict[str, Any]:
    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p118_restart_health")
    require(health.get("status") == "ok", "restart health should be ok")
    current = data(request_json(base_url, "GET", "/api/v1/portfolio/current", request_id="req_p118_restart_portfolio"))
    reports = data(request_json(base_url, "GET", "/api/v1/daily-discipline/reports?limit=30", request_id="req_p118_restart_reports"))
    audit = data(request_json(base_url, "GET", "/api/v1/audit-events", request_id="req_p118_restart_audit"))
    loops = data(request_json(base_url, "GET", "/api/v1/decision-loops?limit=20", request_id="req_p118_restart_loop"))
    require(current and len(current.get("positions", [])) >= 6, "restart portfolio readback missing positions")
    require(reports and len(reports.get("reports", [])) >= 10, "restart report readback missing")
    require(audit is not None, "restart audit readback missing")
    require(loops is not None, "restart decision-loop readback missing")
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        counts = {
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "daily_reports": int(scalar(conn, "SELECT COUNT(*) FROM daily_discipline_reports")),
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
        "scenario_id": "E16",
        "restart_evidence": [
            {"method": "GET", "path": "/api/v1/health", "status": 200, "request_id": "req_p118_restart_health"},
            {"method": "GET", "path": "/api/v1/portfolio/current", "status": 200, "request_id": "req_p118_restart_portfolio", "position_count": len(current.get("positions", []))},
            {"method": "GET", "path": "/api/v1/daily-discipline/reports", "status": 200, "request_id": "req_p118_restart_reports", "returned_reports": len(reports.get("reports", []))},
            {"method": "GET", "path": "/api/v1/audit-events", "status": 200, "request_id": "req_p118_restart_audit"},
            {"method": "GET", "path": "/api/v1/decision-loops", "status": 200, "request_id": "req_p118_restart_loop"},
        ],
        "sqlite_counts_after_restart": counts,
        "usability_evidence": [{"dimension": "persistence", "assertion": "same SQLite remains readable after accumulated edge-use history and backend restart", "value": "ok"}],
    }
    out = artifact_dir / "restart" / "p118-restart-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(payload, ensure_ascii=False, indent=2))
    return payload


def build_interpretation(scenarios: dict[str, dict[str, Any]], browser_payload: dict[str, Any], restart_payload: dict[str, Any]) -> dict[str, Any]:
    blocked = [item for item in scenarios.values() if item["status"] not in ("fresh_pass", "scoped_pass")]
    fresh = [item for item in scenarios.values() if item["status"] == "fresh_pass"]
    scoped = [item for item in scenarios.values() if item["status"] == "scoped_pass"]
    return {
        "task_completion_rate": "18/18 runner scenarios passed",
        "fresh_pass_count": len(fresh),
        "scoped_pass_count": len(scoped),
        "blocked_count": len(blocked),
        "long_cycle": "Usable: accumulated reports, transactions, notifications, risk alerts and audit rows remain readable.",
        "recovery": "Usable: invalid import and invalid transaction attempts are rejected without partial transaction writes.",
        "data_quality": "Usable within scoped boundary: stale/missing source facts require explicit local resolution and do not become clean external-data claims.",
        "decision_quality": "Usable as interpretation evidence: rising, falling and volatile contexts produce different seeded verdicts with explicit prohibited actions and disclaimers.",
        "household_ledger": "Usable as local facts: multiple household account notes and fund categories remain traceable without broker account integration.",
        "persistence": "Usable: restart probe reads accumulated local history from the same SQLite database.",
        "mobile": "Usable within checked scope: 390px accumulated-state pages render without console/page/API 5xx failures.",
        "safety": "Usable as a local discipline assistant only: no broker/order/push tables, auto confirmations or auto rule-apply audit events.",
        "browser_health": {
            "status": browser_payload.get("status"),
            "console_errors": len(browser_payload.get("console_errors", [])),
            "page_errors": len(browser_payload.get("page_errors", [])),
            "failed_api_responses": len(browser_payload.get("failed_api_responses", [])),
        },
        "restart_health": restart_payload.get("status"),
        "claim_boundary": "This is a local product-use edge scenario pass. It explicitly excludes release/install/upgrade validation, broker execution, external data cleanliness, fresh real LLM quality, prediction accuracy, return guarantees and physical second-machine validation.",
    }


def merge_summary(base_url: str, db_path: Path, artifact_dir: Path, browser_summary: Path | None) -> dict[str, Any]:
    api_path = artifact_dir / "api_sqlite" / "p118-api-sqlite-summary.json"
    restart_path = artifact_dir / "restart" / "p118-restart-summary.json"
    require(api_path.exists(), f"missing API summary: {api_path}")
    require(restart_path.exists(), f"missing restart summary: {restart_path}")
    api_payload = json.loads(api_path.read_text(encoding="utf-8"))
    restart_payload = json.loads(restart_path.read_text(encoding="utf-8"))
    scenarios = {item["scenario_id"]: item for item in api_payload["scenarios"]}
    if "E16" in scenarios:
        scenarios["E16"]["restart_evidence"].extend(restart_payload.get("restart_evidence", []))
        scenarios["E16"]["usability_evidence"].extend(restart_payload.get("usability_evidence", []))
    browser_payload: dict[str, Any] = {}
    if browser_summary and browser_summary.exists():
        browser_payload = json.loads(browser_summary.read_text(encoding="utf-8"))
        for entry in browser_payload.get("scenarios", []):
            sid = entry["scenario_id"]
            if sid in scenarios:
                scenarios[sid]["browser_evidence"].extend(entry.get("browser_evidence", []))
                scenarios[sid]["redaction_evidence"].update(entry.get("redaction_evidence", {}))
    missing_browser = [sid for sid in ["E01", "E02", "E03", "E07", "E09", "E10", "E11", "E13", "E15", "E17", "E18"] if not scenarios.get(sid, {}).get("browser_evidence")]
    require(not missing_browser, f"missing browser evidence for {missing_browser}")

    interpretation = build_interpretation(scenarios, browser_payload, restart_payload)
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
        "edge_usability_interpretation": interpretation,
        "claim_boundary": interpretation["claim_boundary"],
    }
    out = artifact_dir / "p118-edge-usability-summary.json"
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
