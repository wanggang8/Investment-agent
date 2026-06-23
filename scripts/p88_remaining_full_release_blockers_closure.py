#!/usr/bin/env python3
"""Generate the P88 remaining-blocker closure matrix.

P88 only upgrades rows with direct fresh evidence. Parser contracts, accepted-local
fixtures, and seeded source-preverification are allowed to document gaps, but not
to claim real external-provider coverage.
"""

from __future__ import annotations

import argparse
import json
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-matrix.md"
OUT_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-matrix.md"
OUT_CLOSURE = ROOT / "docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-closure.md"
ARTIFACT_DIR = ROOT / "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers"
SUMMARY = ARTIFACT_DIR / "p88-acceptance-summary.json"
PREVERIFY = ARTIFACT_DIR / "p88-source-preverification.json"
INVENTORY = ARTIFACT_DIR / "p88-inventory.json"

P88_ROWS = {
    "REQ-02-022",
    "REQ-02-025",
    "REQ-04-016",
    "REQ-04-025",
    "REQ-05-003",
    "REQ-05-004",
    "REQ-05-005",
    "REQ-06-023",
    "REQ-06-024",
    "REQ-08-004",
    "REQ-08-023",
    "REQ-09-001",
    "REQ-09-003",
    "REQ-09-004",
    "REQ-09-006",
    "REQ-09-007",
    "REQ-09-008",
    "REQ-09-009",
    "REQ-09-010",
    "REQ-09-013",
    "REQ-09-023",
    "REQ-09-024",
    "REQ-09-025",
    "REQ-09-027",
    "REQ-10-004",
    "REQ-13-010",
    "REQ-17-015",
}

FULL_RUN_COMMAND = "bash scripts/p88-remaining-full-release-blockers-acceptance.sh && python3 scripts/p88_remaining_full_release_blockers_closure.py --check"
FOCUSED_GO_COMMAND = "go test ./internal/application/workflow ./internal/domain/rule -run 'TestP88|TestBuildExpectedReturn|TestExpectedReturnNode|TestExpectedReturnSampleCount|TestP88StructuredData' -count=1 && go test ./internal/application/handler ./internal/application/service -run 'P88|ExpectedReturn|Rebalance|RuleProposal|SOPAddendum' -count=1"
WEB_COMMAND = "npm --prefix web test -- --run DecisionTrace.test.tsx PortfolioPage.test.tsx RulesPage.test.tsx"


REAL_PASS = {
    "REQ-02-022": (
        "Fresh P88 real UI consultation for 159915 used A-level formal buy-logic-break evidence with A/S independent source count=2, produced sell_only, prohibited 新增买入/加仓, and SQLite readback passed.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/db-readback-check.log",
        "None for this row; keep P88 source-transition path in regression.",
    ),
    "REQ-02-025": (
        "Fresh P88 real UI consultation for 600000 used one A-level formal major-negative source, displayed current A/S independent source count=1, produced frozen_watch, prohibited 主动交易建议, and SQLite readback passed.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/db-readback-check.log",
        "None for this row; keep source-count provenance assertions in regression.",
    ),
    "REQ-04-025": (
        "P88 records source-preverification governance for capital_flow, margin_financing, and constituent_financial candidates before expanding production collector scope, including access limits, target SQLite paths, failure behavior, and blocked runtime-provider status.",
        "python3 scripts/p88_source_preverification.py --check",
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-preverification.json",
        "No provider-readiness gap for this governance row; the underlying data-category rows remain partial until real providers are verified.",
    ),
    "REQ-06-023": (
        "Fresh P88 UI/API/SQLite readback proves source-verified buy-logic-break transition to sell_only and prohibits buying or adding.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-06-024": (
        "Fresh P88 UI/API/SQLite readback proves multi-source-insufficient major-event transition to frozen_watch and pauses active trading advice while waiting for more evidence.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-001": (
        "Fresh P88 expected-return UI/API/SQLite evidence for 510300 shows historical-sample scenario ranges, current valuation context, three holding-class coverage, and triggered sell-evaluation prompts.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/db-readback-check.log",
        "None for this row; dynamic probability-update rows remain separate.",
    ),
    "REQ-09-003": (
        "Fresh P88 expected-return path derives probabilities from historical similar-sample proportions and displays sample window/screening/probability basis through real UI and SQLite readback.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-006": (
        "Fresh P88 UI/API/SQLite evidence shows optimistic scenario range 12.00%~18.00% with probability 20.0% from historical similar-sample proportion.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-007": (
        "Fresh P88 evidence shows the base scenario as the highest-frequency historical path with probability 60.0%.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-008": (
        "Fresh P88 UI/API/SQLite evidence always displays the pessimistic scenario, including -10.00%~-2.00% range and 20.0% probability for risk reminder.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-009": (
        "Fresh P88 report block contains target, horizon, scenarios, probability basis, sample metadata, support/missing/supplement fields, disclaimer, holding coverage, and sell-evaluation section in UI/API/readback.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for the report-field set proven by P88; periodic dynamic rows remain separate.",
    ),
    "REQ-09-010": (
        "Fresh P88 decision detail displays target name 沪深300ETF and code 510300 in the expected-return report.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-expected-return-detail.png; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-013": (
        "Fresh P88 expected-return UI displays the future-12-month label and all three scenario return ranges.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-expected-return-detail.png; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-09-027": (
        "Fresh P88 source-transition UI path also exercises sample_count=2, shows no return ranges, displays qualitative risk text, and lists supplement data to collect.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-transition-sell-only.png; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
    "REQ-10-004": (
        "Fresh P88 /positions UI triggers quarterly +/-15% rebalance review, API calculates manual buy/sell recommendation amounts, audit action run_local_task is written, and forbidden broker/order/auto-confirm checks pass.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-rebalance-review.png; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/db-readback-check.log",
        "None for this row.",
    ),
    "REQ-13-010": (
        "Fresh P88 /rules UI generates a high-frequency uncovered-scenario SOP addendum proposal with pending_user_confirm status, notification linkage, audit event, and no active rule auto-apply.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-sop-proposal.png; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/db-readback-check.log",
        "None for this row.",
    ),
    "REQ-17-015": (
        "Fresh P88 source-verified buy-logic-break UI/API/SQLite readback proves the workflow transition into sell_only.",
        FULL_RUN_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-acceptance-summary.json",
        "None for this row.",
    ),
}

PARTIAL = {
    "REQ-04-016": (
        "P88 adds structured parser/readback normalization for capital_flow, margin_financing, and constituent_financial fields, but no verified non-mock runtime provider proves the full structured data center breadth.",
        FOCUSED_GO_COMMAND + " && python3 scripts/p88_source_preverification.py --check",
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-preverification.json",
        "Needs real no-login/no-paid/no-Level2/no-high-frequency runtime provider collection and SQLite readback for the missing structured categories.",
    ),
    "REQ-05-003": (
        "P88 proves capital-flow field parser/readback contract only; source-preverification marks runtime provider status not_verified and real_pass_eligible=false.",
        "python3 scripts/p88_source_preverification.py --check && " + FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-preverification.json",
        "Needs verified real provider and persisted SQLite readback for date/net_inflow/net_outflow.",
    ),
    "REQ-05-004": (
        "P88 proves margin-financing field parser/readback contract only; source-preverification marks runtime provider status not_verified and real_pass_eligible=false.",
        "python3 scripts/p88_source_preverification.py --check && " + FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-preverification.json",
        "Needs verified real provider and persisted SQLite readback for date/margin_balance/balance_change_rate.",
    ),
    "REQ-05-005": (
        "P88 proves constituent-financial field parser/readback contract only; source-preverification marks runtime provider status not_verified and real_pass_eligible=false.",
        "python3 scripts/p88_source_preverification.py --check && " + FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-source-preverification.json",
        "Needs verified real provider and persisted SQLite readback for revenue/net_profit/growth/disclosure_date.",
    ),
    "REQ-08-004": (
        "P88 does not implement or run a fresh extreme-fear UI path with historical similar-scenario data display.",
        "N/A",
        "N/A",
        "Needs real UI/API/SQLite scenario for extreme fear locking active trading and showing historical similar-scenario data.",
    ),
    "REQ-08-023": (
        "P88 has focused deterministic tests for dynamic monitoring, but no fresh real UI/API/SQLite rerun proves scenario-probability downshift.",
        FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log",
        "Needs real UI rerun comparing before/after probabilities and persisted readback.",
    ),
    "REQ-09-004": (
        "P88 adds deterministic valuation/fundamental/market-state hooks and tests, but the full dynamic update path is not proven through real UI/API/SQLite.",
        FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log",
        "Needs real UI/API/SQLite proof that valuation, fundamentals, and market state update the report.",
    ),
    "REQ-09-023": (
        "P88 supports assumption-check fields and focused tests, but periodic assumption checking is not proven as a product workflow.",
        FOCUSED_GO_COMMAND + " && " + WEB_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log; docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/web-component-tests.log",
        "Needs scheduled or user-visible periodic assumption-check workflow with UI/API/SQLite readback.",
    ),
    "REQ-09-024": (
        "P88 has deterministic two-month below-expectation warning tests, but no fresh real UI/API/SQLite scenario exercises the warning end to end.",
        FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log",
        "Needs real UI scenario with two consecutive months below expectation and persisted downshift warning.",
    ),
    "REQ-09-025": (
        "P88 has deterministic one-month pessimistic-path monitoring tests, but no fresh real UI/API/SQLite scenario proves the user probability-adjustment suggestion.",
        FOCUSED_GO_COMMAND,
        "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/go-workflow-tests.log",
        "Needs real UI scenario for one-month pessimistic path and readback of manual probability-adjustment suggestion.",
    ),
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
    text = str(value)
    return text.replace("\\", "\\\\").replace("|", "\\|").replace("\n", "<br>")


def read_matrix(path: Path) -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in path.read_text(encoding="utf-8").splitlines():
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
            raise SystemExit(f"Invalid matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit(f"Matrix header not found: {path}")
    return header, rows


def p88_decision(row: dict[str, str]) -> tuple[str, str, str, str, str, str]:
    rid = row["requirement_id"]
    if rid in REAL_PASS:
        basis, command, artifact, gap = REAL_PASS[rid]
        return ("real_pass", basis, command, artifact, gap, "Keep in P88 regression; no row-specific follow-up.")
    if rid in PARTIAL:
        basis, command, artifact, gap = PARTIAL[rid]
        return ("partial", basis, command, artifact, gap, "Open a dedicated P89/P90 row-specific acceptance before full-release claim.")
    if row.get("full_release_required") == "True" and row.get("p86_status") == "real_pass":
        return ("real_pass", "Already resolved before P88; P88 does not rewrite prior evidence.", "N/A", row.get("p86_fresh_evidence_artifact", "N/A"), "None for this row.", "Keep covered by existing regression evidence.")
    return (row.get("p86_status") or row.get("status") or "partial", "No P88 change; row is outside the 27 P88 remaining blockers.", "N/A", "N/A", row.get("p86_remaining_gap") or "N/A", row.get("p86_next_action") or "N/A")


def generate() -> dict:
    header, rows = read_matrix(SOURCE_MATRIX)
    missing = sorted(P88_ROWS - {row["requirement_id"] for row in rows})
    if missing:
        raise SystemExit(f"P88 rows missing from source matrix: {missing}")
    if set(REAL_PASS) | set(PARTIAL) != P88_ROWS:
        raise SystemExit("P88 status maps do not cover exactly the P88 owned rows")

    p88_cols = ["p88_status", "p88_closure_basis", "p88_fresh_evidence_command", "p88_fresh_evidence_artifact", "p88_remaining_gap", "p88_next_action"]
    out_header = header + p88_cols
    out_rows: list[dict[str, str]] = []
    for row in rows:
        next_row = dict(row)
        status, basis, command, artifact, gap, next_action = p88_decision(row)
        next_row.update(
            {
                "p88_status": status,
                "p88_closure_basis": basis,
                "p88_fresh_evidence_command": command,
                "p88_fresh_evidence_artifact": artifact,
                "p88_remaining_gap": gap,
                "p88_next_action": next_action,
            }
        )
        out_rows.append(next_row)

    OUT_MATRIX.write_text(render_matrix(out_header, out_rows), encoding="utf-8")

    full_required = [row for row in out_rows if row.get("full_release_required") == "True"]
    remaining = [row for row in full_required if row.get("p88_status") != "real_pass"]
    p88_status_counts = Counter(row["p88_status"] for row in out_rows if row["requirement_id"] in P88_ROWS)
    payload = {
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "source_matrix": str(SOURCE_MATRIX.relative_to(ROOT)),
        "matrix": str(OUT_MATRIX.relative_to(ROOT)),
        "acceptance_summary": str(SUMMARY.relative_to(ROOT)),
        "source_preverification": str(PREVERIFY.relative_to(ROOT)),
        "p88_owned_rows": len(P88_ROWS),
        "p88_real_pass": p88_status_counts.get("real_pass", 0),
        "p88_partial": p88_status_counts.get("partial", 0),
        "full_release_required_rows": len(full_required),
        "full_release_required_non_real_pass_after_p88": len(remaining),
        "remaining_requirement_ids": [row["requirement_id"] for row in remaining],
    }
    OUT_CLOSURE.write_text(render_closure(payload), encoding="utf-8")
    return payload


def render_matrix(header: list[str], rows: list[dict[str, str]]) -> str:
    lines = [
        "| " + " | ".join(escape_cell(item) for item in header) + " |",
        "| " + " | ".join("---" for _ in header) + " |",
    ]
    for row in rows:
        lines.append("| " + " | ".join(escape_cell(row.get(item, "")) for item in header) + " |")
    return "\n".join(lines) + "\n"


def render_closure(payload: dict) -> str:
    remaining = payload["remaining_requirement_ids"]
    lines = [
        "# P88 Remaining Full-Release Blockers Closure",
        "",
        f"- Generated at: {payload['generated_at']}",
        f"- Source matrix: `{payload['source_matrix']}`",
        f"- P88 matrix: `{payload['matrix']}`",
        f"- Acceptance summary: `{payload['acceptance_summary']}`",
        f"- Source preverification: `{payload['source_preverification']}`",
        "",
        "## Result",
        "",
        f"- P88 owned rows: {payload['p88_owned_rows']}",
        f"- Upgraded to `real_pass`: {payload['p88_real_pass']}",
        f"- Kept `partial`: {payload['p88_partial']}",
        f"- Full-release-required rows still not `real_pass` after P88: {payload['full_release_required_non_real_pass_after_p88']}",
        "",
        "P88 does not claim full original-requirement pass. It closes the directly evidenced source-transition, expected-return, rebalance, SOP-proposal, and source-preverification governance rows. It keeps structured real-provider and dynamic-probability rows partial where P88 has tests or contracts but no real UI/API/SQLite end-to-end proof.",
        "",
        "## Remaining Full-Release Blockers",
        "",
    ]
    if remaining:
        for rid in remaining:
            lines.append(f"- `{rid}`")
    else:
        lines.append("- None")
    lines.extend(
        [
            "",
            "## Evidence Boundary",
            "",
            "- Accepted-local, fixture, stub, or manually seeded evidence does not upgrade capital-flow, margin-financing, or constituent-financial runtime-provider rows.",
            "- P88 verifies broker/order/external-push/auto-confirm/auto-rule-apply absence only on exercised P88 paths; it does not replace broader product G9 scans.",
            "- P88 package or full-release refresh is not implied by this closure.",
            "",
        ]
    )
    return "\n".join(lines)


def check() -> dict:
    if not OUT_MATRIX.exists() or not OUT_CLOSURE.exists():
        raise SystemExit("status=failed\nreason=missing_closure_outputs")
    for path in (SUMMARY, PREVERIFY, INVENTORY):
        if not path.exists():
            raise SystemExit(f"status=failed\nreason=missing_artifact:{path}")
    summary = json.loads(SUMMARY.read_text(encoding="utf-8"))
    preverify = json.loads(PREVERIFY.read_text(encoding="utf-8"))
    if summary.get("status") != "passed":
        raise SystemExit("status=failed\nreason=acceptance_summary_not_passed")
    if summary.get("source_preverification", {}).get("status") != "passed":
        raise SystemExit("status=failed\nreason=source_preverification_not_passed")
    if preverify.get("summary", {}).get("blocked_count") != 3:
        raise SystemExit("status=failed\nreason=preverification_blocked_count")

    _, rows = read_matrix(OUT_MATRIX)
    p88_rows = [row for row in rows if row["requirement_id"] in P88_ROWS]
    counts = Counter(row["p88_status"] for row in p88_rows)
    if counts.get("real_pass") != 17 or counts.get("partial") != 10:
        raise SystemExit(f"status=failed\nreason=p88_counts:{dict(counts)}")
    forbidden_real = {"REQ-04-016", "REQ-05-003", "REQ-05-004", "REQ-05-005", "REQ-08-004", "REQ-08-023", "REQ-09-004", "REQ-09-023", "REQ-09-024", "REQ-09-025"}
    mistakenly_real = sorted(row["requirement_id"] for row in p88_rows if row["requirement_id"] in forbidden_real and row["p88_status"] == "real_pass")
    if mistakenly_real:
        raise SystemExit(f"status=failed\nreason=overclaimed_rows:{mistakenly_real}")
    full_remaining = [row["requirement_id"] for row in rows if row.get("full_release_required") == "True" and row.get("p88_status") != "real_pass"]
    if len(full_remaining) != 10:
        raise SystemExit(f"status=failed\nreason=remaining_count:{len(full_remaining)}")
    return {"status": "passed", "p88_real_pass": counts.get("real_pass", 0), "p88_partial": counts.get("partial", 0), "remaining": full_remaining}


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    if args.check:
        result = check()
    else:
        result = generate()
    print(
        "p88_closure:"
        f"status=passed:"
        f"real_pass={result.get('p88_real_pass')}:"
        f"partial={result.get('p88_partial')}:"
        f"remaining={len(result.get('remaining', result.get('remaining_requirement_ids', [])))}:"
        f"matrix={OUT_MATRIX}"
    )


if __name__ == "__main__":
    main()
