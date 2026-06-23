#!/usr/bin/env python3
"""Generate/check P90 capital-flow provider closure matrix and acceptance record."""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p90-capital-flow-provider"
SUMMARY = ARTIFACT_DIR / "p90-acceptance-summary.json"
MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p90-capital-flow-provider-matrix.md"
CLOSURE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p90-capital-flow-provider-closure.md"


ROWS = {
    "REQ-04-016": {
        "text": "结构化数据中心使用 SQLite 存储行情、估值、资金流向、财务等数据。",
        "status": "real_pass",
        "basis": "P89 had already proved margin-financing and constituent-financial provider/UI/API/SQLite readback; P90 adds fresh Eastmoney H5 capital-flow provider verification plus Settings UI refresh, market snapshot API readback, and SQLite runtime readback for date/net_inflow/net_outflow/raw_net_flow.",
        "evidence": "bash scripts/p90-capital-flow-provider-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for the P89-known structured-data breadth blocker.",
        "next": "Keep P90 source preverification, browser acceptance, and SQLite readback in regression before future release refreshes.",
    },
    "REQ-05-003": {
        "text": "资金流向 日期、资金净流入、资金净流出。",
        "status": "real_pass",
        "basis": "P90 verifies the public Eastmoney H5 capital-flow source, maps the daily raw net-flow value into directional net_inflow/net_outflow, triggers product refresh through the Settings UI, and reads back the same fields through API and SQLite.",
        "evidence": "bash scripts/p90-capital-flow-provider-acceptance.sh",
        "artifact": str(SUMMARY.relative_to(ROOT)),
        "gap": "None for date/net_inflow/net_outflow/raw_net_flow in the tested public-provider path.",
        "next": "Keep the H5 provider availability and directional mapping in real provider regression; do not synthesize values if the provider becomes unavailable.",
    },
}


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def load_summary() -> dict:
    require(SUMMARY.exists(), "missing_p90_acceptance_summary")
    payload = json.loads(SUMMARY.read_text(encoding="utf-8"))
    require(payload.get("status") == "passed", "acceptance_not_passed")
    source = payload.get("source_preverification") or {}
    category = source.get("category") or {}
    require(category.get("category") == "capital_flow", "source_category")
    require(category.get("real_pass_eligible") is True, "capital_flow_not_real_pass_eligible")
    fields = category.get("fields") or {}
    for key in ("date", "net_inflow", "net_outflow", "raw_net_flow"):
        require(key in fields, f"source_missing:{key}")
    browser = payload.get("browser") or {}
    require(browser.get("status") == "passed", "browser_not_passed")
    require((browser.get("pre_refresh") or {}).get("capital_flow_absent") is True, "pre_refresh_capital_flow_not_absent")
    readback = (browser.get("provider_readback") or {}).get("capital_flow") or {}
    for key in ("date", "net_inflow", "net_outflow", "raw_net_flow"):
        require(key in readback, f"browser_readback_missing:{key}")
    db = payload.get("db_readback") or {}
    for key in ("pre_refresh_capital_flow_absent", "api_ui_capital_flow_readback", "sqlite_capital_flow_readback", "directional_net_flow_mapping"):
        require(db.get(key) == "passed", f"db_readback:{key}")
    require((payload.get("go_tests") or {}).get("workflow", {}).get("status") == "passed", "go_workflow_not_passed")
    require((payload.get("web_tests") or {}).get("settings_page", {}).get("status") == "passed", "web_settings_not_passed")
    require((payload.get("frontend_build") or {}).get("status") == "passed", "frontend_build_not_passed")
    return payload


def write_matrix() -> None:
    lines = [
        "# P90 Capital-Flow Provider Closure Matrix",
        "",
        "| requirement_id | p90_status | p90_closure_basis | p90_fresh_evidence_command | p90_fresh_evidence_artifact | p90_remaining_gap | p90_next_action |",
        "| --- | --- | --- | --- | --- | --- | --- |",
    ]
    for row_id, row in ROWS.items():
        lines.append(f"| {row_id} | {row['status']} | {row['basis']} | `{row['evidence']}` | `{row['artifact']}` | {row['gap']} | {row['next']} |")
    MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_closure(payload: dict) -> None:
    fields = ((payload.get("source_preverification") or {}).get("category") or {}).get("fields") or {}
    snapshot_id = (payload.get("browser") or {}).get("provider_readback", {}).get("market_snapshot_id")
    lines = [
        "# P90 Capital-Flow Provider Closure",
        "",
        "## Result",
        "",
        "- P90 evaluated the 2 full-release-required rows that remained after P89: `REQ-04-016` and `REQ-05-003`.",
        "- P90 upgraded both rows to `real_pass` using fresh public-provider, real Settings UI, market snapshot API, and SQLite readback evidence.",
        "- P90 conclusion: `release_ready_full_original_requirement_real_pass_candidate_with_p90_capital_flow_closure`.",
        "- No P89-chain full-release-required row is known to remain non-`real_pass` after P90.",
        "- P90 does not refresh the P76 package and does not claim physical second-machine validation, remote release, Git tag, broker integration, trading, external push, automatic confirmation, automatic rule application, paid/login/auth source, Level2 source, high-frequency source, or return guarantee.",
        "",
        "## Evidence",
        "",
        f"- Acceptance summary: `{SUMMARY.relative_to(ROOT)}`",
        f"- Matrix: `{MATRIX.relative_to(ROOT)}`",
        "- Command: `bash scripts/p90-capital-flow-provider-acceptance.sh`",
        "- Source preverification: `python3 scripts/p90_source_preverification.py --check`",
        "- SQLite readback: `python3 scripts/p90_sqlite_readback_check.py <sqlite> <browser-results.json> <p90-source-preverification.json>`",
        "",
        "## Verified Capital-Flow Fields",
        "",
        f"- Runtime snapshot: `{snapshot_id}`",
        f"- `date`: `{fields.get('date')}`",
        f"- `net_inflow`: `{fields.get('net_inflow')}`",
        f"- `net_outflow`: `{fields.get('net_outflow')}`",
        f"- `raw_net_flow`: `{fields.get('raw_net_flow')}`",
        "- Directional mapping: positive raw value maps to `net_inflow`; negative raw value maps to `net_outflow`; raw value is preserved.",
        "",
        "## Safety Boundary",
        "",
        "- P90 only performs low-frequency read-only public data collection and local fact persistence.",
        "- If the H5 provider becomes unavailable, the product must degrade/block dependent claims and must not synthesize capital-flow values.",
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
    require("REQ-04-016 | real_pass" in matrix_text, "req_04_016_not_real_pass")
    require("REQ-05-003 | real_pass" in matrix_text, "req_05_003_not_real_pass")
    require("No P89-chain full-release-required row is known to remain non-`real_pass` after P90." in closure_text, "closure_remaining_statement")
    print(f"p90_closure:status=passed:real_pass=2:partial=0:matrix={MATRIX}")


if __name__ == "__main__":
    main()
