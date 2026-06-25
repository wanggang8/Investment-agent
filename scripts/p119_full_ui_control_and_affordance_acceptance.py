#!/usr/bin/env python3
"""P119 full UI control and affordance acceptance runner."""

from __future__ import annotations

import argparse
import json
import sqlite3
import time
import urllib.request
from pathlib import Path
from typing import Any


CHANGE_ID = "p119-full-ui-control-and-affordance-acceptance"
NOW = "2026-06-25T09:00:00Z"


class AcceptanceFailure(RuntimeError):
    pass


def require(condition: bool, message: str) -> None:
    if not condition:
        raise AcceptanceFailure(message)


def request_json(base_url: str, method: str, path: str, body: dict[str, Any] | None = None, request_id: str = "req_p119") -> dict[str, Any]:
    payload = None
    headers = {"Accept": "application/json", "X-Request-ID": request_id}
    if body is not None:
        payload = json.dumps(body, ensure_ascii=False).encode("utf-8")
        headers["Content-Type"] = "application/json"
    req = urllib.request.Request(base_url.rstrip("/") + path, data=payload, headers=headers, method=method)
    with urllib.request.urlopen(req, timeout=30) as resp:
        raw = resp.read().decode("utf-8")
        return json.loads(raw) if raw else {}


def data(envelope: dict[str, Any]) -> Any:
    return envelope.get("data")


def scalar(conn: sqlite3.Connection, sql: str, args: tuple[Any, ...] = ()) -> Any:
    row = conn.execute(sql, args).fetchone()
    if row is None:
        return 0
    return row[0] if row[0] is not None else 0


def seed_sqlite(db_path: Path) -> None:
    conn = sqlite3.connect(db_path)
    try:
        conn.execute(
            """
            INSERT OR REPLACE INTO daily_discipline_reports (
              report_id,local_date,scope,symbol_set_hash,source_type,source_id,decision_id,status,
              summary,failure_code,failure_reason,created_at,updated_at
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
            """,
            (
                "p119_report_01",
                "2026-06-25",
                "holdings",
                "p119_symbol_hash",
                "manual",
                "p119_seed",
                "decision_p119_confirm",
                "success",
                "P119 全页面控件验收纪律报告",
                "",
                "",
                NOW,
                NOW,
            ),
        )
        for notification_id, source_type, source_id, title in [
            ("notif_p119_unread", "risk_alert", "risk_p119_active", "P119 风险预警待读"),
            ("notif_p119_rule", "rule_proposal", "prop_seed_p119", "P119 规则提案待读"),
        ]:
            conn.execute(
                """
                INSERT OR REPLACE INTO notifications (
                  notification_id,type,severity,title,message,source_type,source_id,created_at,read_at
                ) VALUES (?,?,?,?,?,?,?,?,NULL)
                """,
                (notification_id, "p119_control_acceptance", "warning", title, "P119 本地通知控件验收，不发送外部消息", source_type, source_id, NOW),
            )
        conn.execute(
            """
            INSERT OR REPLACE INTO risk_alerts (
              alert_id,risk_type,severity,sop_status,symbol,trigger_summary,trigger_context_json,
              prohibited_actions_json,suggested_actions_json,related_decision_id,related_report_id,
              related_notification_id,related_audit_event_id,last_triggered_at,created_at,updated_at
            ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
            """,
            (
                "risk_p119_active",
                "valuation_high",
                "warning",
                "active",
                "510300",
                "P119 风险 SOP 按钮验收",
                '{"sop":"manual_review","data_prerequisites":["valuation","position"]}',
                '["新增买入","自动交易"]',
                '["记录继续观察","记录升级复核","记录本地解除预警"]',
                "decision_p119_confirm",
                "p119_report_01",
                "notif_p119_unread",
                "audit_p119_seed",
                NOW,
                NOW,
                NOW,
            ),
        )
        seed_decision(conn, "decision_p119_confirm", "510300", "P119 决策确认按钮验收", "pending", "hold", "持有并等待人工确认")
        seed_decision(conn, "decision_p119_error", "159915", "P119 错误标注按钮验收", "pending", "hold", "用于错误复盘标注")
        seed_market_and_evidence(conn)
        conn.commit()
    finally:
        conn.close()


def seed_decision(conn: sqlite3.Connection, decision_id: str, symbol: str, question: str, confirmation_status: str, verdict: str, text: str) -> None:
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
            "completed",
            "formal_trade_advice",
            "normal",
            "in_scope",
            "P119 seeded UI control acceptance",
            "satisfied",
            "",
            '{"heat":"neutral"}',
            '["calm"]',
            '[{"rule_id":"p119.manual_boundary","rule_name":"人工边界","severity":"warning","description":"必须由用户线下确认"}]',
            "[]",
            verdict,
            text,
            '["自动交易","券商接口","自动确认"]',
            '["记录计划","已手动执行","标记待观察","标记错误"]',
            confirmation_status,
            None,
            "market_p119_seed",
            "v_p119",
            '[{"agent_name":"P119LocalAnalyst","conclusion":"本地控件验收材料","key_reasons":["按钮必须落到本地事实"],"risk_warnings":["不会自动交易"],"confidence":"medium","evidence_ids":["sum_p119_seed"]}]',
            '{"precision_status":"available","reason":"P119 seeded scenario","sample_count":12,"sample_window":"2024-2026","probability_basis":"local acceptance seed","scenarios":[{"scenario":"base","return_range":"-2%~4%","probability":0.5,"trigger":"manual review"}],"disclaimer":"仅为本地情景分析，不承诺收益"}',
            '[{"step":"control_inventory","result":"seeded"}]',
            '{"p119":"control_acceptance"}',
            NOW,
        ),
    )


def seed_market_and_evidence(conn: sqlite3.Connection) -> None:
    conn.execute(
        """
        INSERT OR REPLACE INTO market_snapshots (
          market_snapshot_id,symbol,trade_date,close_price,pe_percentile,pb_percentile,market_metrics_json,created_at
        ) VALUES (?,?,?,?,?,?,?,?)
        """,
        (
            "market_p119_seed",
            "000300",
            "2026-06-05",
            4000,
            81,
            70,
            '{"source_name":"csindex","source_level":"A","source_type":"index_basic","metadata":{"p34_source_health":{"index_valuation_files":{"freshness":"stale","data_date":"2026-06-05","failure_category":"stale","affected_symbols":["000300"],"source_level":"A","source_type":"index_basic"}},"p34_data_categories":["index_valuation_files"]}}',
            NOW,
        ),
    )
    conn.execute(
        "INSERT OR REPLACE INTO intelligence_items (intelligence_id,source_name,source_level,original_url,published_at,captured_at,content_hash,raw_title,raw_text_ref,created_at) VALUES (?,?,?,?,?,?,?,?,?,?)",
        ("intel_p119_seed", "P119 本地验收材料", "A", "https://example.com/p119", NOW, NOW, "hash_p119_seed", "P119 控件验收证据", "p119_seed", NOW),
    )
    conn.execute(
        "INSERT OR REPLACE INTO intelligence_summary (summary_id,intelligence_id,symbol,entity,event_type,impact_direction,summary,source_level,evidence_role,time_weight,relevance_score,verification_group_id,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
        ("sum_p119_seed", "intel_p119_seed", "510300", "510300", "normal", "neutral", "P119 控件验收证据摘要", "A", "formal", 1, 1, "vg_p119_seed", NOW),
    )
    conn.execute(
        "INSERT OR REPLACE INTO rag_chunks (chunk_id,summary_id,chunk_text,chunk_hash,index_status,created_at) VALUES (?,?,?,?,?,?)",
        ("chunk_p119_seed", "sum_p119_seed", "P119 控件验收索引片段", "chunk_hash_p119_seed", "pending", NOW),
    )
    conn.execute(
        "INSERT OR REPLACE INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
        ("ver_p119_seed", "vg_p119_seed", "event_p119_seed", "510300", "normal", "formal", "satisfied", 2, 1, "A", NOW, '["sum_p119_seed"]', NOW),
    )


def seed_api(base_url: str, db_path: Path, artifact_dir: Path) -> dict[str, Any]:
    health = request_json(base_url, "GET", "/api/v1/health", request_id="req_p119_health")
    require(health.get("status") == "ok", "health should be ok")
    seed_sqlite(db_path)
    body = {
        "cash": 12000,
        "total_assets": 18020,
        "adjust_reason": "P119 UI 控件验收本地校准",
        "positions": [
            {"symbol": "510300", "name": "沪深300ETF", "quantity": 1000, "cost_price": 3.5, "current_price": 4.1, "buy_date": "2026-05-10", "position_state": "normal", "buy_reason": "P119 初始本地持仓", "asset_tag": "core"},
            {"symbol": "159915", "name": "创业板ETF", "quantity": 800, "cost_price": 2.1, "current_price": 2.4, "buy_date": "2026-05-12", "position_state": "sell_only", "buy_reason": "P119 卫星本地持仓", "asset_tag": "satellite"},
        ],
    }
    adjust = data(request_json(base_url, "POST", "/api/v1/portfolio/adjustments", body, request_id="req_p119_seed_portfolio"))
    require(adjust and adjust.get("position_count", 0) >= 2, "portfolio seed should create positions")
    summary = {
        "status": "seeded",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "api_evidence": [
            {"method": "GET", "path": "/api/v1/health", "status": 200, "request_id": "req_p119_health"},
            {"method": "POST", "path": "/api/v1/portfolio/adjustments", "status": 200, "request_id": "req_p119_seed_portfolio"},
        ],
    }
    out = artifact_dir / "api_sqlite" / "p119-api-sqlite-seed.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return summary


def final_counts(db_path: Path) -> dict[str, int]:
    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    try:
        return {
            "positions": int(scalar(conn, "SELECT COUNT(*) FROM positions")),
            "portfolio_snapshots": int(scalar(conn, "SELECT COUNT(*) FROM portfolio_snapshots")),
            "position_transactions": int(scalar(conn, "SELECT COUNT(*) FROM position_transactions")),
            "operation_confirmations": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations")),
            "error_cases": int(scalar(conn, "SELECT COUNT(*) FROM error_cases")),
            "risk_alert_resolved": int(scalar(conn, "SELECT COUNT(*) FROM risk_alerts WHERE alert_id='risk_p119_active' AND sop_status='resolved'")),
            "unread_p119_notifications": int(scalar(conn, "SELECT COUNT(*) FROM notifications WHERE notification_id LIKE 'notif_p119_%' AND read_at IS NULL")),
            "data_quality_resolutions": int(scalar(conn, "SELECT COUNT(*) FROM data_quality_gate_resolutions WHERE symbol='000300'")),
            "rule_proposals": int(scalar(conn, "SELECT COUNT(*) FROM rule_proposals WHERE title LIKE '%SOP 补充提案%'")),
            "intelligence_items": int(scalar(conn, "SELECT COUNT(*) FROM intelligence_items")),
            "rag_chunks": int(scalar(conn, "SELECT COUNT(*) FROM rag_chunks")),
            "audit_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events")),
            "forbidden_broker_order_push_tables": int(scalar(conn, "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND (LOWER(name) LIKE '%broker%' OR LOWER(name) LIKE '%order%' OR LOWER(name) LIKE '%external_push%' OR LOWER(name) LIKE '%trade_execution%')")),
            "auto_confirmation_rows": int(scalar(conn, "SELECT COUNT(*) FROM operation_confirmations WHERE confirmation_type LIKE 'auto%'")),
            "auto_rule_apply_events": int(scalar(conn, "SELECT COUNT(*) FROM audit_events WHERE LOWER(action) LIKE '%auto%' AND (LOWER(action) LIKE '%rule%' OR LOWER(action) LIKE '%confirm%' OR LOWER(action) LIKE '%trade%')")),
        }
    finally:
        conn.close()


def merge_summary(base_url: str, db_path: Path, artifact_dir: Path, browser_summary_path: Path) -> dict[str, Any]:
    browser = json.loads(browser_summary_path.read_text(encoding="utf-8"))
    counts = final_counts(db_path)
    require(browser.get("status") == "passed", "browser summary did not pass")
    require(browser.get("route_count") == 22, f"expected 22 route records, got {browser.get('route_count')}")
    require(browser.get("unclassified_controls") == 0, "unclassified controls found")
    require(browser.get("unnamed_controls") == 0, "unnamed controls found")
    require(browser.get("layout_issue_count") == 0, "layout issues found")
    require(browser.get("console_error_count") == 0, "console errors found")
    require(browser.get("page_error_count") == 0, "page errors found")
    require(browser.get("api_5xx_count") == 0, "API 5xx found")
    require(browser.get("toggle_interaction_count", 0) >= 13, "upstream toggle interaction coverage missing")
    require(browser.get("toggle_issue_count", 0) == 0, "upstream toggle interaction issues found")
    require(counts["positions"] >= 2, "positions readback missing")
    require(counts["portfolio_snapshots"] >= 2, "portfolio snapshot readback missing")
    require(counts["position_transactions"] >= 1, "UI transaction write missing")
    require(counts["operation_confirmations"] >= 2, "decision/portfolio confirmations missing")
    require(counts["error_cases"] >= 1, "marked error readback missing")
    require(counts["risk_alert_resolved"] == 1, "risk alert lifecycle write missing")
    require(counts["unread_p119_notifications"] == 0, "notification mark-read missing")
    require(counts["data_quality_resolutions"] >= 1, "data quality resolution missing")
    require(counts["rule_proposals"] >= 1, "rule proposal write missing")
    require(counts["intelligence_items"] >= 2, "local knowledge/evidence write missing")
    require(counts["rag_chunks"] >= 1, "evidence chunks missing")
    require(counts["forbidden_broker_order_push_tables"] == 0, "forbidden broker/order/push table present")
    require(counts["auto_confirmation_rows"] == 0, "automatic confirmations present")
    require(counts["auto_rule_apply_events"] == 0, "automatic rule application events present")

    payload = {
        "status": "passed",
        "change": CHANGE_ID,
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "base_url": base_url,
        "sqlite_path": str(db_path),
        "route_count": browser["route_count"],
        "control_count": browser["control_count"],
        "control_category_counts": browser["control_category_counts"],
        "mobile_route_count": browser["mobile_route_count"],
        "toggle_interaction_count": browser.get("toggle_interaction_count", 0),
        "toggle_issue_count": browser.get("toggle_issue_count", 0),
        "browser_summary": browser,
        "sqlite_counts": counts,
        "claim_boundary": "P119 validates local UI controls, layout, and backend consistency in an isolated local run. It does not claim install/upgrade/release validation, broker execution, external push, automatic trading, automatic confirmation, automatic rule application, fresh provider coverage, physical second-machine validation, prediction accuracy, or return guarantees.",
    }
    out = artifact_dir / "p119-ui-control-summary.json"
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", required=True)
    parser.add_argument("--sqlite", required=True)
    parser.add_argument("--artifact-dir", required=True)
    parser.add_argument("--browser-summary")
    parser.add_argument("--merge-only", action="store_true")
    args = parser.parse_args()

    db_path = Path(args.sqlite)
    artifact_dir = Path(args.artifact_dir)
    if args.merge_only:
        require(args.browser_summary, "--browser-summary required for merge")
        merge_summary(args.base_url, db_path, artifact_dir, Path(args.browser_summary))
    else:
        seed_api(args.base_url, db_path, artifact_dir)


if __name__ == "__main__":
    main()
