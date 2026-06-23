#!/usr/bin/env python3
"""Read-only SQLite checks for P89 dynamic probability acceptance."""

from __future__ import annotations

import json
import sqlite3
import sys
from pathlib import Path


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def decision(conn: sqlite3.Connection, decision_id: str) -> sqlite3.Row:
    row = conn.execute(
        "SELECT decision_id,symbol,final_verdict_status,prohibited_actions_json,expected_return_scenarios_json FROM decision_records WHERE decision_id=?",
        (decision_id,),
    ).fetchone()
    require(row is not None, f"missing_decision:{decision_id}")
    return row


def probs(expected: dict) -> list[float]:
    out: list[float] = []
    for item in expected.get("scenarios") or []:
        value = item.get("probability")
        if value is None:
            value = item.get("Probability")
        out.append(round(float(value), 4))
    return out


def main() -> None:
    if len(sys.argv) != 4:
        raise SystemExit("usage: p89_sqlite_readback_check.py <sqlite> <browser-results.json> <source-preverification.json>")
    db_path = Path(sys.argv[1])
    browser = json.loads(Path(sys.argv[2]).read_text(encoding="utf-8"))
    source_preverification = json.loads(Path(sys.argv[3]).read_text(encoding="utf-8"))

    baseline_id = browser.get("baseline", {}).get("decision_id")
    dynamic_id = browser.get("dynamic", {}).get("decision_id")
    extreme_id = browser.get("extreme_fear", {}).get("decision_id")
    provider_snapshot_id = browser.get("provider_readback", {}).get("market_snapshot_id")
    require(baseline_id and dynamic_id and extreme_id, "missing_browser_decision_ids")
    require(provider_snapshot_id and provider_snapshot_id != "market_p89_600000", "missing_runtime_provider_snapshot_id")

    conn = sqlite3.connect(f"file:{db_path}?mode=ro", uri=True)
    conn.row_factory = sqlite3.Row

    baseline = decision(conn, baseline_id)
    baseline_expected = json.loads(baseline["expected_return_scenarios_json"] or "{}")
    require(probs(baseline_expected) == [0.2, 0.6, 0.2], f"baseline_probabilities:{probs(baseline_expected)}")
    require(baseline_expected.get("probability_basis") == "historical_similar_sample_proportion", "baseline_probability_basis")

    dynamic = decision(conn, dynamic_id)
    dynamic_expected = json.loads(dynamic["expected_return_scenarios_json"] or "{}")
    dynamic_probs = probs(dynamic_expected)
    require(dynamic_probs[0] < 0.2 and dynamic_probs[1] < 0.6 and dynamic_probs[2] > 0.2, f"dynamic_probabilities:{dynamic_probs}")
    triggers = (dynamic_expected.get("sell_evaluation") or {}).get("triggers") or []
    for trigger in ("scenario_probability_downshift", "two_month_assumption_downshift", "one_month_pessimistic_path"):
        require(trigger in triggers, f"dynamic_trigger:{trigger}")
    checks = dynamic_expected.get("assumption_checks") or []
    require(any(item.get("months_below") == 2 for item in checks), "two_month_assumption_check")

    extreme = decision(conn, extreme_id)
    extreme_expected = json.loads(extreme["expected_return_scenarios_json"] or "{}")
    actions = json.loads(extreme["prohibited_actions_json"] or "[]")
    require("主动交易建议" in actions, "extreme_fear_prohibited_action")
    contexts = extreme_expected.get("historical_contexts") or []
    require(len(contexts) == 1 and contexts[0].get("label") == "极端恐惧样本", "extreme_historical_context")
    extreme_triggers = (extreme_expected.get("sell_evaluation") or {}).get("triggers") or []
    require("extreme_fear_historical_context" in extreme_triggers, "extreme_historical_trigger")

    source_categories = {item.get("category"): item for item in source_preverification.get("categories") or []}
    require(source_categories.get("margin_financing", {}).get("real_pass_eligible") is True, "margin_provider_eligible")
    require(source_categories.get("constituent_financial", {}).get("real_pass_eligible") is True, "financial_provider_eligible")
    require(source_categories.get("capital_flow", {}).get("real_pass_eligible") is False, "capital_flow_should_remain_blocked")

    market = conn.execute(
        "SELECT market_snapshot_id,symbol,margin_balance,margin_balance_change,market_metrics_json FROM market_snapshots WHERE market_snapshot_id=?",
        (provider_snapshot_id,),
    ).fetchone()
    require(market is not None, "missing_structured_market_snapshot")
    require(market["symbol"] == "600000", f"provider_snapshot_symbol:{market['symbol']}")
    market_metrics = json.loads(market["market_metrics_json"] or "{}")
    structured = ((market_metrics.get("metadata") or {}).get("p88_structured_fields") or {})
    margin_fields = source_categories["margin_financing"]["fields"]
    financial_fields = source_categories["constituent_financial"]["fields"]
    require(abs(float(market["margin_balance"]) - float(margin_fields["margin_balance"])) < 0.01, "sqlite_margin_balance")
    require(abs(float(market["margin_balance_change"]) - float(margin_fields["balance_change_rate"])) < 0.000001, "sqlite_margin_change_rate")
    require((structured.get("margin_financing") or {}).get("date") == margin_fields["date"], "sqlite_margin_date")
    require(float((structured.get("constituent_financial") or {}).get("revenue") or 0) == float(financial_fields["revenue"]), "sqlite_financial_revenue")
    require(float((structured.get("constituent_financial") or {}).get("net_profit") or 0) == float(financial_fields["net_profit"]), "sqlite_financial_net_profit")
    require((structured.get("constituent_financial") or {}).get("disclosure_date") == financial_fields["disclosure_date"], "sqlite_financial_disclosure_date")
    require((structured.get("capital_flow") or {}) == {}, "sqlite_capital_flow_should_not_be_synthesized")

    print("status=passed")
    print(f"baseline_decision_id={baseline_id}")
    print(f"dynamic_decision_id={dynamic_id}")
    print(f"extreme_fear_decision_id={extreme_id}")
    print(f"runtime_provider_snapshot_id={provider_snapshot_id}")
    print("baseline_probabilities=0.2,0.6,0.2")
    print("dynamic_probability_downshift=passed")
    print("two_month_assumption_check=passed")
    print("one_month_pessimistic_path=passed")
    print("extreme_fear_historical_context=passed")
    print("margin_provider_eligible=passed")
    print("financial_provider_eligible=passed")
    print("sqlite_margin_financing_readback=passed")
    print("sqlite_constituent_financial_readback=passed")
    print("api_ui_provider_readback=passed")
    print("capital_flow_provider_blocked=preserved")


if __name__ == "__main__":
    main()
