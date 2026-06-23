#!/usr/bin/env python3
import json
import sqlite3
import sys
import time
import urllib.request
from pathlib import Path


def health_json(freshness):
    categories = {
        "symbol_profile": ("A", freshness.get("symbol_profile", "fresh"), ["510300"]),
        "fund_profile": ("B", freshness.get("fund_profile", "fresh"), ["510300"]),
        "tracked_index": ("A", freshness.get("tracked_index", "fresh"), ["000300"]),
        "market_price": ("B", freshness.get("market_price", "fresh"), ["510300"]),
        "valuation_percentiles": ("A", freshness.get("valuation_percentiles", "fresh"), ["000300"]),
        "liquidity": ("B", freshness.get("liquidity", "fresh"), ["510300"]),
        "sentiment_proxy": ("C", freshness.get("sentiment_proxy", "fresh"), ["510300"]),
    }
    source_health = {}
    for category, (level, state, symbols) in categories.items():
        source_health[category] = {
            "freshness": state,
            "data_date": "2026-06-19",
            "affected_symbols": symbols,
            "source_level": level,
        }
        if state not in ("fresh", "stubbed"):
            source_health[category]["failure_category"] = state
    return json.dumps({
        "source_name": "p74_acceptance_fixture",
        "source_level": "A",
        "source_type": "readiness_fixture",
        "captured_at": "2026-06-19T08:00:00Z",
        "metadata": {
            "p34_source_health": source_health,
            "p34_data_categories": list(categories.keys()),
        },
    }, ensure_ascii=False)


def ensure_active_rule(db):
    count = db.execute("SELECT COUNT(*) FROM rule_versions WHERE status='active'").fetchone()[0]
    if count:
        return
    db.execute(
        "INSERT INTO rule_versions (rule_version,status,rules_json,effective_at,created_at) VALUES (?,?,?,?,?)",
        ("p74-acceptance-active-rule", "active", '{"evidence":{"min_high_grade_sources":2},"safety":{"no_auto_trade":true}}', "2026-06-19T08:00:00Z", "2026-06-19T08:00:00Z"),
    )


def seed_market(db, scenario, freshness):
    db.execute(
        "INSERT INTO market_snapshots (market_snapshot_id,symbol,trade_date,close_price,turnover_rate,pe_percentile,pb_percentile,volume_percentile,volatility_percentile,liquidity_state,sentiment_state,market_metrics_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
        (
            f"market_p74_{scenario}_{int(time.time() * 1000)}",
            "510300",
            "2026-06-19",
            4.75,
            1.2,
            28,
            35,
            40,
            30,
            "normal",
            "neutral",
            health_json(freshness),
            now_for_scenario(scenario),
        ),
    )


def seed_verification(db, scenario, status, independent, high_grade, level="A", role="formal"):
    db.execute(
        "INSERT INTO source_verifications (verification_id,verification_group_id,event_id,symbol,event_type,evidence_role,verification_status,independent_source_count,high_grade_independent_source_count,highest_source_level,latest_published_at,evidence_ids_json,created_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)",
        (
            f"verify_p74_{scenario}_{int(time.time() * 1000)}",
            f"group_p74_{scenario}",
            f"event_p74_{scenario}",
            "510300",
            "normal",
            role,
            status,
            independent,
            high_grade,
            level,
            "2026-06-19T08:00:00Z",
            '["summary_p74_a","summary_p74_b"]',
            now_for_scenario(scenario),
        ),
    )


def now_for_scenario(name):
    order = {
        "ready": "2026-06-19T08:00:01Z",
        "missing_valuation": "2026-06-19T08:00:02Z",
        "background_only": "2026-06-19T08:00:03Z",
        "single_source": "2026-06-19T08:00:04Z",
        "multi_source": "2026-06-19T08:00:05Z",
        "ui_ready": "2026-06-19T08:00:06Z",
    }
    return order.get(name, "2026-06-19T08:00:00Z")


def readiness(base_url, symbol="510300"):
    with urllib.request.urlopen(f"{base_url}/api/v1/knowledge-readiness?symbol={symbol}", timeout=20) as response:
        payload = json.loads(response.read().decode("utf-8"))
    return payload["data"]


def dep(data, category):
    for item in data["data_dependencies"]:
        if item["category"] == category:
            return item
    raise AssertionError(f"missing dependency {category}: {data}")


def assert_status(data, expected):
    if data["status"] != expected:
        raise AssertionError(f"expected status {expected}, got {data['status']}: {json.dumps(data, ensure_ascii=False)}")


def run_scenario(db, base_url, scenario, market_freshness, verification, expected, checks):
    seed_market(db, scenario, market_freshness)
    seed_verification(db, scenario, *verification)
    db.commit()
    data = readiness(base_url)
    assert_status(data, expected)
    for check in checks:
        check(data)
    return {
        "id": scenario,
        "status": data["status"],
        "symbol": data["symbol"],
        "dependencies": {item["category"]: item["status"] for item in data["data_dependencies"]},
        "knowledge_reference_count": len(data["knowledge_references"]),
        "llm_context_summary_present": bool(data.get("llm_context_summary")),
    }


def main():
    if len(sys.argv) != 4:
        raise SystemExit("usage: p74_readiness_api_check.py <sqlite-path> <artifact-dir> <base-url>")
    sqlite_path, artifact_dir, base_url = sys.argv[1:4]
    Path(artifact_dir).mkdir(parents=True, exist_ok=True)
    db = sqlite3.connect(sqlite_path)
    ensure_active_rule(db)
    db.execute("DELETE FROM market_snapshots WHERE symbol IN ('510300', '999999')")
    db.execute("DELETE FROM source_verifications WHERE symbol IN ('510300', '999999')")
    db.commit()
    results = []

    results.append(run_scenario(
        db, base_url, "ready", {}, ("satisfied", 3, 2),
        "ready",
        [
            lambda data: dep(data, "valuation_percentiles")["status"] == "ready" or (_ for _ in ()).throw(AssertionError("valuation not ready")),
            lambda data: dep(data, "formal_evidence")["status"] == "ready" or (_ for _ in ()).throw(AssertionError("formal evidence not ready")),
            lambda data: dep(data, "active_rule")["status"] == "ready" or (_ for _ in ()).throw(AssertionError("active rule not ready")),
        ],
    ))
    results.append(run_scenario(
        db, base_url, "missing_valuation", {"valuation_percentiles": "parse_error"}, ("satisfied", 3, 2),
        "degraded",
        [lambda data: dep(data, "valuation_percentiles")["status"] == "degraded" or (_ for _ in ()).throw(AssertionError("valuation should degrade"))],
    ))
    results.append(run_scenario(
        db, base_url, "background_only", {}, ("background_only", 1, 0, "C", "background"),
        "degraded",
        [
            lambda data: dep(data, "formal_evidence")["freshness"] == "background_only" or (_ for _ in ()).throw(AssertionError("background-only formal evidence state missing")),
            lambda data: all(not item["formal_evidence_allowed"] for item in data["knowledge_references"]) or (_ for _ in ()).throw(AssertionError("built-in knowledge must not become formal evidence")),
        ],
    ))
    results.append(run_scenario(
        db, base_url, "single_source", {}, ("satisfied", 1, 1),
        "degraded",
        [lambda data: dep(data, "formal_evidence")["status"] == "degraded" or (_ for _ in ()).throw(AssertionError("single source should degrade"))],
    ))
    results.append(run_scenario(
        db, base_url, "multi_source", {}, ("satisfied", 3, 2),
        "ready",
        [lambda data: dep(data, "formal_evidence")["status"] == "ready" or (_ for _ in ()).throw(AssertionError("multi-source formal evidence should pass"))],
    ))

    blocked = readiness(base_url, "999999")
    assert_status(blocked, "blocked")
    if blocked["symbol_profile"]["known"]:
        raise AssertionError("unknown symbol profile must not be fabricated")
    results.append({
        "id": "out_of_scope_symbol_profile",
        "status": blocked["status"],
        "symbol": blocked["symbol"],
        "dependencies": {item["category"]: item["status"] for item in blocked["data_dependencies"]},
        "knowledge_reference_count": len(blocked["knowledge_references"]),
        "llm_context_summary_present": bool(blocked.get("llm_context_summary")),
    })

    seed_market(db, "ui_ready", {})
    seed_verification(db, "ui_ready", "satisfied", 3, 2)
    db.execute(
        "UPDATE decision_records SET analyst_reports_json=? WHERE decision_id='decision_smoke_p30'",
        (json.dumps([{
            "agent_name": "P74ReadinessAnalyst",
            "conclusion": "估值证据与多源正式证据满足本地准备度，保持人工复核。",
            "key_reasons": ["安全边际纪律已进入 LLM 摘要"],
            "risk_warnings": ["LLM 不能覆盖规则最终裁决"],
            "confidence": "medium",
            "evidence_ids": ["summary_p74_a"],
            "input_summary": "value 510300 principles=master.graham.margin_of_safety data_readiness=valuation_percentiles=ready,formal_evidence=ready boundary=背景知识不能满足正式证据",
            "prompt_version": "p74-knowledge-readiness-v1",
            "quality_status": "passed",
        }], ensure_ascii=False),),
    )
    db.commit()

    payload = {
        "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
        "status": "passed",
        "scenarios": results,
        "safety_note": "P74 API acceptance mutates only the temporary SQLite database used by this runner; it does not refresh public data, call LLM providers, trade, confirm actions, apply rules, or write release claims.",
    }
    Path(artifact_dir, "api-readiness-results.json").write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps(payload, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
