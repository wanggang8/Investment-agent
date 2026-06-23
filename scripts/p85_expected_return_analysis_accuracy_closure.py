#!/usr/bin/env python3
"""Generate P85 expected-return analysis accuracy artifacts."""

from __future__ import annotations

import argparse
import json
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p84-portfolio-confirmation-data-impact-matrix.md"
P85_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p85-expected-return-analysis-accuracy-matrix.md"
P85_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p85-expected-return-analysis-accuracy-closure.md"
P85_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p85-expected-return-analysis"
P85_SUMMARY = P85_ASSET_DIR / "expected-return-summary.json"
P85_DB_CHECK = P85_ASSET_DIR / "db-readback-check.log"

P85_UI_COMMAND = (
    "P85_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p85-expected-return-analysis "
    "bash scripts/p85-expected-return-analysis-acceptance.sh"
)

P85_COLUMNS = [
    "p85_status",
    "p85_closure_basis",
    "p85_fresh_evidence_command",
    "p85_fresh_evidence_artifact",
    "p85_remaining_gap",
    "p85_next_action",
]

P85_PLAN_IDS = {
    "REQ-02-005",
    "REQ-02-014",
    "REQ-08-004",
    "REQ-08-023",
    "REQ-09-001",
    "REQ-09-003",
    "REQ-09-004",
    "REQ-09-005",
    "REQ-09-006",
    "REQ-09-007",
    "REQ-09-008",
    "REQ-09-009",
    "REQ-09-010",
    "REQ-09-011",
    "REQ-09-012",
    "REQ-09-013",
    "REQ-09-014",
    "REQ-09-015",
    "REQ-09-016",
    "REQ-09-018",
    "REQ-09-019",
    "REQ-09-020",
    "REQ-09-021",
    "REQ-09-022",
    "REQ-09-023",
    "REQ-09-024",
    "REQ-09-025",
    "REQ-09-027",
    "REQ-09-028",
    "REQ-13-010",
    "REQ-16-022",
}

P85_UPGRADE_IDS = {
    "REQ-02-005",
    "REQ-02-014",
    "REQ-09-005",
    "REQ-09-011",
    "REQ-09-012",
    "REQ-09-014",
    "REQ-09-015",
    "REQ-09-016",
    "REQ-09-018",
    "REQ-09-019",
    "REQ-09-020",
    "REQ-09-021",
    "REQ-09-022",
    "REQ-09-028",
    "REQ-16-022",
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
        and go_tests.get("workflow", {}).get("status") == "passed"
        and go_tests.get("handler", {}).get("status") == "passed"
        and db.get("available_precision_status") == "available"
        and db.get("available_sample_count") == "20"
        and db.get("unavailable_precision_status") == "unavailable"
        and db.get("operation_confirmations_p85") == "0"
        and db.get("forbidden_broker_order_push_tables") == "0"
        and db.get("auto_confirmation_rows") == "0"
    )


def p85_basis(requirement_id: str) -> str:
    bases = {
        "REQ-02-005": "Fresh P85 real UI/API/SQLite evidence displays scenario ranges, exact probabilities for sufficient samples, downside risk scenario, dynamic sell prompts, and non-absolute disclaimers.",
        "REQ-02-014": "Fresh P85 decision detail and SQLite readback show expected-return output is explicitly probabilistic reference material and contains no return promise or automatic action.",
        "REQ-09-005": "Fresh P85 complete-data UI path triggers sell evaluation when current return reaches the upside/base boundaries and target return.",
        "REQ-09-011": "Fresh P85 decision detail shows current date and current price/net value from the persisted market snapshot.",
        "REQ-09-012": "Fresh P85 decision detail shows current PE/PB percentile fields from the persisted market snapshot.",
        "REQ-09-014": "Fresh P85 covers available sample probabilities and unavailable sample degradation where no precise probability is shown.",
        "REQ-09-015": "Fresh P85 decision detail and SQLite readback show dynamic sell evaluation status, triggers, prompts, and manual-only actions.",
        "REQ-09-016": "Fresh P85 covers take-profit, stop-loss/recheck, reassessment, and user target-return triggers across complete and downside scenarios.",
        "REQ-09-018": "Fresh P85 complete-data scenario triggers `upside_lower_bound_reached` and prompts mobile take-profit review.",
        "REQ-09-019": "Fresh P85 complete-data scenario triggers `base_upper_bound_exceeded` and prompts staged take-profit review.",
        "REQ-09-020": "Fresh P85 downside scenario triggers `downside_lower_bound_breached` and prompts buy-thesis recheck.",
        "REQ-09-021": "Fresh P85 complete-data scenario triggers `base_midpoint_downshift` when previous base midpoint is more than 15 percentage points above the current base midpoint.",
        "REQ-09-022": "Fresh P85 UI input for target return rate triggers `target_return_reached` and only prompts a manual plan review.",
        "REQ-09-028": "Fresh P85 available-probability output is accompanied by sample count, sample window, and screening condition in UI/API/SQLite readback.",
        "REQ-16-022": "Fresh P85 real consultation workflow proves the expected-return module runs end-to-end through UI, persisted decision detail, SQLite readback, and focused Go tests.",
    }
    return bases[requirement_id]


def p85_gap(requirement_id: str) -> str:
    if requirement_id in P85_UPGRADE_IDS:
        return "None for this P85 row; future market-outcome accuracy remains explicitly out of scope."
    gaps = {
        "REQ-08-004": "P85 does not prove an extreme-fear active-trading lock plus historical similar-scenario display.",
        "REQ-08-023": "P85 does not prove automatic downward adjustment of scenario probabilities after a scenario update.",
        "REQ-09-001": "P85 proves deterministic scenario display and sell triggers, but not a full historical-law/current-valuation expected-return model for every holding class.",
        "REQ-09-003": "Current deterministic ranges are not a true historical backtest or similar-valuation frequency model.",
        "REQ-09-004": "Current probabilities/ranges are not dynamically recalibrated by valuation, fundamentals, and market state beyond sample/price trigger inputs.",
        "REQ-09-006": "P85 displays upside scenario, but does not prove probability is derived from historical similar-sample proportions.",
        "REQ-09-007": "P85 displays base scenario, but does not prove it is the highest-frequency path in historical samples.",
        "REQ-09-008": "P85 displays downside scenario and risk prompt, but does not prove the full pessimistic-business-performance model.",
        "REQ-09-009": "P85 proves many report fields, but not the complete report breadth including all child requirements such as 12-month label and historical-frequency provenance.",
        "REQ-09-010": "P85 detail shows code but not a complete fund/security display name in the expected-return report block.",
        "REQ-09-013": "P85 displays scenario ranges but not an explicit future-12-month horizon label in the report block.",
        "REQ-09-023": "P85 does not prove periodic checking of core valuation assumptions.",
        "REQ-09-024": "P85 does not prove a two-month below-expectation assumption tracker or scenario-downshift warning.",
        "REQ-09-025": "P85 does not prove one-month pessimistic-path tracking or user probability-adjustment suggestion.",
        "REQ-09-027": "P85 proves unavailable degradation with no ranges, but UI does not yet show a complete list of data to supplement.",
        "REQ-13-010": "P85 does not prove SOP addendum proposal generation for high-frequency uncovered scenarios.",
    }
    return gaps.get(requirement_id, "P85 evidence is adjacent but not complete for this row.")


def prior_status(row: dict[str, str]) -> str:
    return row.get("p84_status") or row.get("p83_status") or row.get("status") or "partial"


def row_with_p85(row: dict[str, str], passed: bool) -> dict[str, str]:
    requirement_id = row["requirement_id"]
    out = dict(row)
    if requirement_id in P85_PLAN_IDS:
        if passed and requirement_id in P85_UPGRADE_IDS:
            out.update({
                "p85_status": "real_pass",
                "p85_closure_basis": p85_basis(requirement_id),
                "p85_fresh_evidence_command": P85_UI_COMMAND + " && python3 scripts/p85_expected_return_analysis_accuracy_closure.py --check",
                "p85_fresh_evidence_artifact": rel(P85_SUMMARY),
                "p85_remaining_gap": p85_gap(requirement_id),
                "p85_next_action": "Keep in P85 expected-return regression and continue P87/P86 for remaining rows.",
            })
        elif passed:
            status = prior_status(row)
            if status not in {"real_pass", "scoped_pass", "reference_only"}:
                status = "partial"
            out.update({
                "p85_status": status,
                "p85_closure_basis": "P85 evaluated this expected-return or analysis-accuracy row but did not upgrade it because the fresh evidence does not prove the complete row text.",
                "p85_fresh_evidence_command": P85_UI_COMMAND + " && python3 scripts/p85_expected_return_analysis_accuracy_closure.py --check",
                "p85_fresh_evidence_artifact": rel(P85_SUMMARY),
                "p85_remaining_gap": p85_gap(requirement_id),
                "p85_next_action": "Carry forward to P86 final integrated closure or a dedicated row-specific acceptance.",
            })
        else:
            out.update({
                "p85_status": "partial",
                "p85_closure_basis": "P85 evidence did not pass; this row remains non-real-pass.",
                "p85_fresh_evidence_command": P85_UI_COMMAND,
                "p85_fresh_evidence_artifact": rel(P85_SUMMARY),
                "p85_remaining_gap": "Fresh P85 UI/API/SQLite/Go evidence must pass before upgrade.",
                "p85_next_action": "Rerun P85 acceptance after fixing the failing evidence.",
            })
        return out

    out.update({
        "p85_status": prior_status(row),
        "p85_closure_basis": "No P85 upgrade; row is owned by P87/P86 or remains previously scoped/reference.",
        "p85_fresh_evidence_command": "N/A",
        "p85_fresh_evidence_artifact": "N/A",
        "p85_remaining_gap": row.get("p84_remaining_gap", row.get("remaining_gap", "")),
        "p85_next_action": row.get("p84_next_action", row.get("next_action", "")),
    })
    return out


def write_matrix(header: list[str], rows: list[dict[str, str]]) -> None:
    out_header = header + P85_COLUMNS
    lines = [
        "| " + " | ".join(out_header) + " |",
        "| " + " | ".join("---" for _ in out_header) + " |",
    ]
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(col, "")) for col in out_header) + "|")
    P85_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]], summary: dict[str, Any], passed: bool) -> None:
    counts = Counter(row["p85_status"] for row in rows)
    full_rows = [row for row in rows if row.get("full_release_required") == "True"]
    remaining = [row for row in full_rows if row["p85_status"] != "real_pass"]
    upgraded = sorted(row["requirement_id"] for row in rows if row["requirement_id"] in P85_UPGRADE_IDS and row["p85_status"] == "real_pass")
    deferred = sorted(P85_PLAN_IDS - set(upgraded))
    browser = summary.get("browser", {})
    db = summary.get("db_readback", {})
    generated = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    lines = [
        "# P85 Expected Return Analysis Accuracy Closure",
        "",
        f"- Generated at: `{generated}`",
        f"- Status: `{'passed' if passed else 'failed'}`",
        f"- Source matrix: `{rel(SOURCE_MATRIX)}`",
        f"- Output matrix: `{rel(P85_MATRIX)}`",
        f"- Summary artifact: `{rel(P85_SUMMARY)}`",
        f"- Browser status: `{browser.get('status', 'missing')}`",
        f"- SQLite status: `{db.get('status', 'missing')}`",
        f"- LLM mode: `{summary.get('llm_mode', 'unknown')}`",
        "",
        "## Evidence",
        "",
        f"- Command: `{P85_UI_COMMAND}`",
        f"- Browser results: `{rel(P85_ASSET_DIR / 'browser-results.json')}`",
        f"- SQLite readback: `{rel(P85_DB_CHECK)}`",
        f"- Screenshots: `{rel(P85_ASSET_DIR)}/p85-*.png`",
        "- Scenarios: sufficient sample with target-return UI input, downside-boundary UI consult, unavailable-sample UI consult.",
        "",
        "## Row Outcome",
        "",
        f"- Total rows: `{len(rows)}`",
        f"- Counts: `{dict(counts)}`",
        f"- P85 planned rows: `{len(P85_PLAN_IDS)}`",
        f"- P85 upgraded rows: `{len(upgraded)}`",
        f"- Full-release-required rows still non-real-pass: `{len(remaining)}`",
        f"- Upgraded: `{', '.join(upgraded)}`",
        f"- Deferred: `{', '.join(deferred)}`",
        "",
        "## Boundary",
        "",
        "- P85 does not claim future return accuracy, future market-direction accuracy, a real historical backtest model, automatic probability downshift, longitudinal assumption tracking, broker connectivity, automatic trading, automatic confirmation, external push, or return promise.",
        "- Because `DEEPSEEK_API_KEY` was not present in this environment, P85 does not claim fresh real LLM output as acceptance evidence; deterministic workflow, UI/API/SQLite readback, and focused Go tests are the evidence basis.",
    ]
    P85_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def generate() -> tuple[bool, dict[str, int]]:
    header, source_rows = read_source_rows()
    summary = read_json(P85_SUMMARY)
    passed = evidence_passed(summary)
    rows = [row_with_p85(row, passed) for row in source_rows]
    write_matrix(header, rows)
    write_acceptance(rows, summary, passed)
    counts = Counter(row["p85_status"] for row in rows)
    full_remaining = sum(1 for row in rows if row.get("full_release_required") == "True" and row["p85_status"] != "real_pass")
    return passed, {"new_real": sum(1 for row in rows if row["requirement_id"] in P85_UPGRADE_IDS and row["p85_status"] == "real_pass"), "remaining_full": full_remaining, **counts}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    passed, counts = generate()
    ok = passed and counts["new_real"] == len(P85_UPGRADE_IDS) and counts["remaining_full"] == 142
    print(f"p85_expected_return:status={'passed' if ok else 'failed'}:new_real={counts['new_real']}:remaining_full={counts['remaining_full']}")
    if args.check and not ok:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
