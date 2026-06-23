#!/usr/bin/env python3
"""Generate/check P89 closure matrix and acceptance record."""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p89-real-provider-dynamic-probability"
SUMMARY = ARTIFACT_DIR / "p89-acceptance-summary.json"
FINAL_VALIDATION = ARTIFACT_DIR / "final-validation.log"
MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p89-real-provider-dynamic-probability-matrix.md"
CLOSURE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p89-real-provider-dynamic-probability-closure.md"


ROWS = {
    "REQ-04-016": {
        "text": "结构化数据中心广度，覆盖资金流、融资融券、成分财务等字段。",
        "status": "partial",
        "basis": "P89 verified live public provider and product runtime `/api/v1/market/refresh` UI/API/SQLite readback for margin_financing and constituent_financial, but capital_flow provider is blocked.",
        "evidence": "scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "capital_flow date/net_inflow/net_outflow still lacks a successful no-login/no-paid/no-auth/no-Level2/no-high-frequency runtime provider readback in this run.",
        "next": "Find or contract a stable public capital-flow provider; rerun provider preverification and SQLite/API/UI readback before upgrading.",
    },
    "REQ-05-003": {
        "text": "资金流向字段 date/net_inflow/net_outflow。",
        "status": "partial",
        "basis": "P89 attempted live public provider verification; Eastmoney push2 capital-flow endpoint returned curl exit 52 in this environment.",
        "evidence": "python3 scripts/p89_source_preverification.py",
        "artifact": str((ARTIFACT_DIR / "p89-source-preverification.json").relative_to(ROOT)),
        "gap": "No fresh successful runtime provider response and no SQLite readback for capital_flow fields.",
        "next": "Resolve provider availability or replace with a stable public source; do not synthesize values.",
    },
    "REQ-05-004": {
        "text": "融资融券字段 date/margin_balance/balance_change_rate。",
        "status": "real_pass",
        "basis": "P89 verified SSE public margin provider, triggered product market refresh from Settings UI, and read back margin fields through market snapshot API plus SQLite runtime snapshot.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for P89 margin-financing field proof.",
        "next": "Keep in provider/readback regression.",
    },
    "REQ-05-005": {
        "text": "成分财务 revenue/net_profit/growth/disclosure_date。",
        "status": "real_pass",
        "basis": "P89 verified Eastmoney public financial-report provider, triggered product market refresh from Settings UI, and read back revenue/net_profit/growth/disclosure_date through market snapshot API plus SQLite runtime snapshot.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for P89 constituent-financial field proof.",
        "next": "Keep in provider/readback regression.",
    },
    "REQ-08-004": {
        "text": "极端恐惧状态暂停主动交易建议并展示历史相似场景。",
        "status": "real_pass",
        "basis": "Fresh P89 browser UI/API/SQLite scenario for 600000 shows prohibited action 主动交易建议 and historical_contexts label/window/max_drawdown/recovery/source.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for this extreme-fear UI/API/SQLite path.",
        "next": "Keep in real browser regression.",
    },
    "REQ-08-023": {
        "text": "场景更新会随条件变化下修相关概率。",
        "status": "real_pass",
        "basis": "P89 baseline 510300 probabilities 0.2/0.6/0.2 and dynamic 159915 probabilities 0.15/0.54/0.31 prove downshift through UI/API/SQLite.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for tested dynamic downshift path.",
        "next": "Keep in expected-return regression.",
    },
    "REQ-09-004": {
        "text": "预期收益按估值、基本面、市场状态动态更新。",
        "status": "real_pass",
        "basis": "P89 metadata market_state=stress and fundamental_state=below_expectation produce scenario_probability_downshift and lower upside/base probabilities.",
        "evidence": "go test ./internal/application/workflow -run 'TestP89ExpectedReturn' -count=1; bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for tested valuation/fundamental/market-state path.",
        "next": "Keep in expected-return regression.",
    },
    "REQ-09-023": {
        "text": "周期性检查核心估值假设。",
        "status": "real_pass",
        "basis": "P89 UI/API/SQLite readback includes assumption_checks for 盈利增速 expected/actual/months_below.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for tested assumption-check readback.",
        "next": "Keep in expected-return regression.",
    },
    "REQ-09-024": {
        "text": "连续两个月低于预期触发情景下修预警。",
        "status": "real_pass",
        "basis": "P89 dynamic path triggers two_month_assumption_downshift in browser, API, and SQLite readback.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for tested two-month downshift warning.",
        "next": "Keep in expected-return regression.",
    },
    "REQ-09-025": {
        "text": "一个月偏向悲观实际路径提示用户手动调整概率。",
        "status": "real_pass",
        "basis": "P89 dynamic path triggers one_month_pessimistic_path and manual probability-adjustment action through UI/API/SQLite.",
        "evidence": "bash scripts/p89-real-provider-dynamic-probability-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for tested one-month pessimistic path.",
        "next": "Keep in expected-return regression.",
    },
}


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def load_summary() -> dict:
    require(SUMMARY.exists(), "missing_p89_acceptance_summary")
    payload = json.loads(SUMMARY.read_text(encoding="utf-8"))
    require(payload.get("status") == "passed", "acceptance_not_passed")
    source = payload.get("source_preverification") or {}
    categories = {item.get("category"): item for item in source.get("categories") or []}
    require(categories.get("capital_flow", {}).get("real_pass_eligible") is False, "capital_flow_should_be_blocked")
    require(categories.get("margin_financing", {}).get("real_pass_eligible") is True, "margin_provider_not_eligible")
    require(categories.get("constituent_financial", {}).get("real_pass_eligible") is True, "financial_provider_not_eligible")
    db = payload.get("db_readback") or {}
    for key in ("dynamic_probability_downshift", "two_month_assumption_check", "one_month_pessimistic_path", "extreme_fear_historical_context", "sqlite_margin_financing_readback", "sqlite_constituent_financial_readback", "api_ui_provider_readback"):
        require(db.get(key) == "passed", f"db_readback:{key}")
    require((payload.get("browser") or {}).get("provider_readback", {}).get("market_snapshot_id"), "missing_runtime_provider_browser_readback")
    require((payload.get("frontend_build") or {}).get("status") == "passed", "frontend_build_not_passed")
    require(FINAL_VALIDATION.exists() and "status=passed" in FINAL_VALIDATION.read_text(encoding="utf-8"), "missing_final_validation_log")
    return payload


def write_matrix() -> None:
    lines = [
        "# P89 Real Provider And Dynamic Probability Matrix",
        "",
        "| requirement_id | p89_status | p89_closure_basis | p89_fresh_evidence_command | p89_fresh_evidence_artifact | p89_remaining_gap | p89_next_action |",
        "| --- | --- | --- | --- | --- | --- | --- |",
    ]
    for row_id, row in ROWS.items():
        lines.append(f"| {row_id} | {row['status']} | {row['basis']} | `{row['evidence']}` | `{row['artifact']}` | {row['gap']} | {row['next']} |")
    MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_closure(payload: dict) -> None:
    real_pass = sum(1 for row in ROWS.values() if row["status"] == "real_pass")
    partial = sum(1 for row in ROWS.values() if row["status"] != "real_pass")
    lines = [
        "# P89 Real Provider And Dynamic Probability Closure",
        "",
        "## Result",
        "",
        f"- P89 evaluated 10 remaining full-release-required rows from P88.",
        f"- P89 upgraded {real_pass} rows to `real_pass` and preserved {partial} rows as `partial`.",
        "- Remaining non-`real_pass` rows: `REQ-04-016`, `REQ-05-003`.",
        "- Release conclusion: `release_ready_scoped_with_p89_real_provider_dynamic_probability_progress`.",
        "- P89 does not claim full original-requirement pass because capital-flow provider verification is blocked.",
        "",
        "## Evidence",
        "",
        f"- Acceptance summary: `{SUMMARY.relative_to(ROOT)}`",
        f"- Matrix: `{MATRIX.relative_to(ROOT)}`",
        f"- Final validation: `{FINAL_VALIDATION.relative_to(ROOT)}`",
        "- Command: `bash scripts/p89-real-provider-dynamic-probability-acceptance.sh`",
        "- Inventory: `python3 scripts/p89_remaining_real_provider_dynamic_inventory_check.py`",
        "- Source preverification: `python3 scripts/p89_source_preverification.py`",
        "",
        "## Provider Boundary",
        "",
        "- `margin_financing`: verified public SSE provider; Settings UI market refresh, market snapshot API readback, and SQLite runtime snapshot readback passed.",
        "- `constituent_financial`: verified public Eastmoney financial-report provider; Settings UI market refresh, market snapshot API readback, and SQLite runtime snapshot readback passed.",
        "- `capital_flow`: provider verification failed with curl exit 52; no values were synthesized and no real_pass claim is made.",
        "",
        "## Safety Boundary",
        "",
        "- No broker interface, order table, one-click trading, automatic trading, external push, automatic confirmation, automatic rule application, automatic repair/migration/recovery, paid/login/auth source, Level2 source, high-frequency source, or return guarantee was added.",
    ]
    CLOSURE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> None:
    payload = load_summary()
    if "--check" not in sys.argv:
        MATRIX.parent.mkdir(parents=True, exist_ok=True)
        write_matrix()
        write_closure(payload)
    require(MATRIX.exists(), "missing_matrix")
    require(CLOSURE.exists(), "missing_closure")
    matrix_text = MATRIX.read_text(encoding="utf-8")
    closure_text = CLOSURE.read_text(encoding="utf-8")
    require("REQ-05-003 | partial" in matrix_text, "capital_flow_not_partial")
    require("REQ-09-025 | real_pass" in matrix_text, "pessimistic_path_not_real_pass")
    require("P89 upgraded 8 rows" in closure_text, "closure_count")
    print(f"p89_closure:status=passed:real_pass=8:partial=2:matrix={MATRIX}")


if __name__ == "__main__":
    main()
