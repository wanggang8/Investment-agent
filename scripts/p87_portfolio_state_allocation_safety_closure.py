#!/usr/bin/env python3
"""Generate P87 portfolio state/allocation safety artifacts."""

from __future__ import annotations

import argparse
import json
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p85-expected-return-analysis-accuracy-matrix.md"
P87_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p87-portfolio-state-allocation-safety-matrix.md"
P87_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p87-portfolio-state-allocation-safety-closure.md"
P87_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p87-portfolio-state-allocation-safety"
P87_SUMMARY = P87_ASSET_DIR / "portfolio-state-allocation-summary.json"
P87_DB_CHECK = P87_ASSET_DIR / "db-readback-check.log"

P87_UI_COMMAND = (
    "P87_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety "
    "bash scripts/p87-portfolio-state-allocation-acceptance.sh"
)

P87_COLUMNS = [
    "p87_status",
    "p87_closure_basis",
    "p87_fresh_evidence_command",
    "p87_fresh_evidence_artifact",
    "p87_remaining_gap",
    "p87_next_action",
]

P87_PLAN_IDS = {
    "REQ-01-001",
    "REQ-01-006",
    "REQ-02-006",
    "REQ-02-022",
    "REQ-02-024",
    "REQ-02-025",
    "REQ-02-031",
    "REQ-03-004",
    "REQ-03-005",
    "REQ-03-006",
    "REQ-04-003",
    "REQ-04-008",
    "REQ-04-016",
    "REQ-04-025",
    "REQ-05-010",
    "REQ-06-023",
    "REQ-06-024",
    "REQ-07-006",
    "REQ-07-015",
    "REQ-08-018",
    "REQ-08-020",
    "REQ-10-001",
    "REQ-10-002",
    "REQ-10-003",
    "REQ-10-004",
    "REQ-11-005",
    "REQ-14-005",
    "REQ-14-007",
    "REQ-16-028",
    "REQ-16-033",
    "REQ-17-015",
    "REQ-17-024",
}

P87_UPGRADE_IDS = {
    "REQ-02-031",
    "REQ-10-001",
    "REQ-10-002",
    "REQ-10-003",
    "REQ-11-005",
}


def split_markdown_row(line: str) -> list[str]:
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in line.rstrip("\n")[1:-1]:
        if escaped:
            current.append(char)
            escaped = False
        elif char == "\\":
            escaped = True
        elif char == "|":
            cells.append("".join(current).strip())
            current = []
        else:
            current.append(char)
    cells.append("".join(current).strip())
    return cells


def escape_cell(value: object) -> str:
    return str(value).replace("\n", " ").replace("\r", " ").strip().replace("\\", "\\\\").replace("|", "\\|")


def rel(path: Path) -> str:
    return str(path.relative_to(ROOT))


def read_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        return {"status": "missing", "path": rel(path)}
    data = json.loads(path.read_text(encoding="utf-8"))
    return data if isinstance(data, dict) else {"status": "invalid"}


def read_source_rows() -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in SOURCE_MATRIX.read_text(encoding="utf-8").splitlines():
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if header is None:
            if cells and cells[0] == "requirement_id":
                header = cells
            continue
        if set("".join(cells)) <= {"-", ":"}:
            continue
        if len(cells) != len(header):
            raise SystemExit(f"Invalid source matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("Source matrix header not found")
    return header, rows


def evidence_passed(summary: dict[str, Any]) -> bool:
    browser = summary.get("browser", {})
    db = summary.get("db_readback", {})
    go_tests = summary.get("go_tests", {})
    return (
        summary.get("status") == "passed"
        and browser.get("status") == "passed"
        and db.get("status") == "passed"
        and go_tests.get("handler", {}).get("status") == "passed"
        and go_tests.get("rule", {}).get("status") == "passed"
        and db.get("position_count") == "3"
        and db.get("snapshot_cash_ratio") == "0.0800"
        and db.get("core_ratio") == "0.6400"
        and db.get("satellite_ratio") == "0.2700"
        and db.get("cash_bucket_ratio") == "0.0900"
        and db.get("decision_p87_sell_only_status") == "sell_only"
        and db.get("decision_p87_frozen_watch_status") == "frozen_watch"
        and db.get("decision_p87_insufficient_status") == "insufficient_data"
        and db.get("frozen_or_insufficient_confirmations") == "0"
        and db.get("forbidden_broker_order_push_tables") == "0"
        and db.get("auto_confirmation_rows") == "0"
    )


def p87_basis(requirement_id: str) -> str:
    bases = {
        "REQ-02-031": "Fresh P87 `/positions` UI writes and reads back `normal`, `sell_only`, and `frozen_watch` discipline states through API and SQLite.",
        "REQ-10-001": "Fresh P87 real UI portfolio calibration persists core asset tag with 64.00% of total assets, inside the 60%-70% target.",
        "REQ-10-002": "Fresh P87 real UI portfolio import persists satellite asset tag with 27.00% of total assets, inside the 20%-30% target.",
        "REQ-10-003": "Fresh P87 real UI portfolio calibration persists cash of 8.00% of total assets and cash-plus-money-fund bucket of 9.00%, both inside the 5%-10% target.",
        "REQ-11-005": "Fresh P87 `/positions` UI exposes buy-date entry and the API/SQLite/table readback proves buy dates for all three holdings.",
    }
    return bases[requirement_id]


def p87_gap(requirement_id: str) -> str:
    if requirement_id in P87_UPGRADE_IDS:
        return "None for this P87 row; broader product-wide closure remains governed by P86 and later matrices."
    gaps = {
        "REQ-01-001": "P87 proves portfolio-state/allocation slices, but not the complete product goal across all discipline, evidence, feedback, and effect-validation workflows.",
        "REQ-01-006": "P87 does not complete the full controlled-evolution loop of error cases, user confirmation, proposal acceptance, safety audit, and rule application.",
        "REQ-02-006": "P87 does not execute a full rule-change confirmation and safety-audit proposal lifecycle.",
        "REQ-02-022": "P87 proves sell-only position-state UI/API/SQLite readback and seeded decision safety display, but not a workflow-generated buy-logic break backed by at least 2 A/S independent sources.",
        "REQ-02-024": "P87 proves seeded information-insufficient decision readback with no confirmation actions, but not a workflow-generated data-insufficient state across all listed dependency failures.",
        "REQ-02-025": "P87 proves seeded frozen-watch decision readback, but not a workflow-generated multi-source verification failure with source-count evidence.",
        "REQ-03-004": "P87 verifies portfolio snapshot fields, but not the complete daily view of holdings, valuation, allocation, and all risk alerts.",
        "REQ-03-005": "P87 does not prove a full large-gain staged take-profit user scenario.",
        "REQ-03-006": "P87 does not prove a full large-drop hold-vs-buy-thesis-recheck scenario.",
        "REQ-04-003": "P87 does not prove the entire decision cockpit breadth including all account state, risk redline, rules, evidence chain, agent views, verdict, and confirmation areas.",
        "REQ-04-008": "P87 does not execute proposal view/confirm/reject UI acceptance.",
        "REQ-04-016": "P87 proves portfolio SQLite facts, but not market, valuation, capital-flow, and financial structured data breadth.",
        "REQ-04-025": "P87 does not complete public data-source preverification before production collector scope.",
        "REQ-05-010": "P87 does not prove structured intelligence summaries are written to SQLite for this row.",
        "REQ-06-023": "P87 proves sell-only UI/API/SQLite readback, but not the complete buy-logic-break transition and source-verification evidence required by the row.",
        "REQ-06-024": "P87 proves frozen-watch UI/API/SQLite readback, but not the complete multi-source-insufficient transition evidence required by the row.",
        "REQ-07-006": "P87 shows classification-related asset tags, but does not prove the full Peter Lynch classification and industry-logic pause requirement.",
        "REQ-07-015": "P87 proves local seeded safe-degradation readback, but not every required dependency-missing/degraded path listed by this safety row.",
        "REQ-08-018": "P87 does not prove the full objective-data display including PE percentile, sentiment index, and holding valuation together.",
        "REQ-08-020": "P87 does not prove the complete rational reminder flow based on buy thesis and valuation data ending in user confirmation.",
        "REQ-10-004": "P87 proves target ratios, but not a complete quarterly ±15% rebalance action flow.",
        "REQ-14-005": "P87 proves account snapshot facts, but not monthly attribution and full discipline audit.",
        "REQ-14-007": "P87 does not prove complete action/node action/user confirmation/error/proposal/application-time audit history breadth.",
        "REQ-16-028": "P87 does not complete feedback, proposal, confirmation, audit, and write workflow breadth.",
        "REQ-16-033": "P87 does not connect daily report, instant query, error marking, and proposal confirmation end to end.",
        "REQ-17-024": "P87 checks no forbidden broker/order/push/auto-confirmation surfaces in this scenario, but does not run the full release/upgrade preflight prohibition suite.",
        "REQ-17-015": "P87 proves sell-only state can be entered as a local fact, but not a complete buy-logic-broken workflow transition into sell-only.",
    }
    return gaps.get(requirement_id, "P87 evidence is adjacent but not complete for this row.")


def prior_status(row: dict[str, str]) -> str:
    return row.get("p85_status") or row.get("p84_status") or row.get("status") or "partial"


def row_with_p87(row: dict[str, str], passed: bool) -> dict[str, str]:
    requirement_id = row["requirement_id"]
    out = dict(row)
    if requirement_id in P87_PLAN_IDS:
        if passed and requirement_id in P87_UPGRADE_IDS:
            out.update({
                "p87_status": "real_pass",
                "p87_closure_basis": p87_basis(requirement_id),
                "p87_fresh_evidence_command": P87_UI_COMMAND + " && python3 scripts/p87_portfolio_state_allocation_safety_closure.py --check",
                "p87_fresh_evidence_artifact": rel(P87_SUMMARY),
                "p87_remaining_gap": p87_gap(requirement_id),
                "p87_next_action": "Keep in P87 portfolio-state regression and continue P86 for remaining rows.",
            })
        elif passed:
            status = prior_status(row)
            if status not in {"real_pass", "scoped_pass", "reference_only"}:
                status = "partial"
            out.update({
                "p87_status": status,
                "p87_closure_basis": "P87 evaluated this row but did not upgrade it because the fresh evidence does not prove the complete row text.",
                "p87_fresh_evidence_command": P87_UI_COMMAND + " && python3 scripts/p87_portfolio_state_allocation_safety_closure.py --check",
                "p87_fresh_evidence_artifact": rel(P87_SUMMARY),
                "p87_remaining_gap": p87_gap(requirement_id),
                "p87_next_action": "Carry forward to P86 final integrated closure or a dedicated row-specific acceptance.",
            })
        else:
            out.update({
                "p87_status": "partial",
                "p87_closure_basis": "P87 evidence did not pass; this row remains non-real-pass.",
                "p87_fresh_evidence_command": P87_UI_COMMAND,
                "p87_fresh_evidence_artifact": rel(P87_SUMMARY),
                "p87_remaining_gap": "Fresh P87 UI/API/SQLite/Go evidence must pass before upgrade.",
                "p87_next_action": "Rerun P87 acceptance after fixing the failing evidence.",
            })
        return out

    out.update({
        "p87_status": prior_status(row),
        "p87_closure_basis": "No P87 upgrade; row is owned by P86, already passed, or remains previously scoped/reference.",
        "p87_fresh_evidence_command": "N/A",
        "p87_fresh_evidence_artifact": "N/A",
        "p87_remaining_gap": row.get("p85_remaining_gap", row.get("remaining_gap", "")),
        "p87_next_action": row.get("p85_next_action", row.get("next_action", "")),
    })
    return out


def write_matrix(header: list[str], rows: list[dict[str, str]]) -> None:
    out_header = header + P87_COLUMNS
    lines = [
        "| " + " | ".join(out_header) + " |",
        "| " + " | ".join("---" for _ in out_header) + " |",
    ]
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(col, "")) for col in out_header) + "|")
    P87_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]], summary: dict[str, Any], passed: bool) -> None:
    counts = Counter(row["p87_status"] for row in rows)
    full_rows = [row for row in rows if row.get("full_release_required") == "True"]
    remaining = [row for row in full_rows if row["p87_status"] != "real_pass"]
    upgraded = sorted(row["requirement_id"] for row in rows if row["requirement_id"] in P87_UPGRADE_IDS and row["p87_status"] == "real_pass")
    deferred = sorted(P87_PLAN_IDS - set(upgraded))
    browser = summary.get("browser", {})
    db = summary.get("db_readback", {})
    generated = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    lines = [
        "# P87 Portfolio State Allocation Safety Closure",
        "",
        f"- Generated at: `{generated}`",
        f"- Status: `{'passed' if passed else 'failed'}`",
        f"- Source matrix: `{rel(SOURCE_MATRIX)}`",
        f"- Output matrix: `{rel(P87_MATRIX)}`",
        f"- Summary artifact: `{rel(P87_SUMMARY)}`",
        f"- Browser status: `{browser.get('status', 'missing')}`",
        f"- SQLite status: `{db.get('status', 'missing')}`",
        "",
        "## Evidence",
        "",
        f"- Command: `{P87_UI_COMMAND}`",
        f"- Browser results: `{rel(P87_ASSET_DIR / 'browser-results.json')}`",
        f"- SQLite readback: `{rel(P87_DB_CHECK)}`",
        f"- Screenshots: `{rel(P87_ASSET_DIR)}/p87-*.png`",
        "- Scenarios: core/satellite/cash portfolio UI write/readback, buy-date/state persistence, sell-only decision, frozen-watch decision, information-insufficient decision.",
        "",
        "## Row Outcome",
        "",
        f"- Total rows: `{len(rows)}`",
        f"- Counts: `{dict(counts)}`",
        f"- P87 planned rows: `{len(P87_PLAN_IDS)}`",
        f"- P87 upgraded rows: `{len(upgraded)}`",
        f"- Full-release-required rows still non-real-pass: `{len(remaining)}`",
        f"- Upgraded: `{', '.join(upgraded)}`",
        f"- Deferred: `{', '.join(deferred)}`",
        "",
        "## Boundary",
        "",
        "- P87 does not claim complete quarterly rebalance automation, monthly attribution, full audit-history breadth, proposal accept/reject closure, public collector production readiness, broker connectivity, automatic trading, automatic confirmation, external push, release/upgrade preflight full closure, or full original-requirement pass.",
        "- P87 treats user operation as local fact recording only. Sell-only, frozen-watch, and information-insufficient evidence does not authorize automatic trade execution.",
    ]
    P87_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def generate() -> tuple[bool, dict[str, int]]:
    header, source_rows = read_source_rows()
    summary = read_json(P87_SUMMARY)
    passed = evidence_passed(summary)
    rows = [row_with_p87(row, passed) for row in source_rows]
    write_matrix(header, rows)
    write_acceptance(rows, summary, passed)
    counts = Counter(row["p87_status"] for row in rows)
    full_remaining = sum(1 for row in rows if row.get("full_release_required") == "True" and row["p87_status"] != "real_pass")
    return passed, {"new_real": sum(1 for row in rows if row["requirement_id"] in P87_UPGRADE_IDS and row["p87_status"] == "real_pass"), "remaining_full": full_remaining, **counts}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    passed, counts = generate()
    ok = passed and counts["new_real"] == len(P87_UPGRADE_IDS) and counts["remaining_full"] == 137
    print(f"p87_portfolio_state_allocation:status={'passed' if ok else 'failed'}:new_real={counts['new_real']}:remaining_full={counts['remaining_full']}")
    if args.check and not ok:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
