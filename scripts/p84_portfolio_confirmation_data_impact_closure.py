#!/usr/bin/env python3
"""Generate P84 portfolio and confirmation data-impact artifacts."""

from __future__ import annotations

import argparse
import json
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p83-governance-traceability-matrix.md"
P84_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p84-portfolio-confirmation-data-impact-matrix.md"
P84_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p84-portfolio-confirmation-data-impact-closure.md"
P84_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p84-portfolio-confirmation"
P84_SUMMARY = P84_ASSET_DIR / "portfolio-confirmation-summary.json"
P84_BROWSER = P84_ASSET_DIR / "browser-results.json"
P84_DB_CHECK = P84_ASSET_DIR / "db-readback-check.log"

P84_UI_COMMAND = (
    "P84_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation "
    "bash scripts/p84-portfolio-confirmation-acceptance.sh"
)

P84_COLUMNS = [
    "p84_status",
    "p84_closure_basis",
    "p84_fresh_evidence_command",
    "p84_fresh_evidence_artifact",
    "p84_remaining_gap",
    "p84_next_action",
]

P84_PLAN_IDS = {
    "REQ-01-001",
    "REQ-01-006",
    "REQ-02-006",
    "REQ-02-022",
    "REQ-02-024",
    "REQ-02-025",
    "REQ-02-031",
    "REQ-02-033",
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
    "REQ-11-002",
    "REQ-11-005",
    "REQ-11-019",
    "REQ-14-005",
    "REQ-14-007",
    "REQ-16-028",
    "REQ-16-033",
    "REQ-17-015",
    "REQ-17-024",
}

P84_UPGRADE_IDS = {
    "REQ-02-033",
    "REQ-11-002",
    "REQ-11-019",
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
    with path.open(encoding="utf-8") as fh:
        data = json.load(fh)
    return data if isinstance(data, dict) else {"status": "invalid"}


def read_log_kv(path: Path) -> dict[str, str]:
    if not path.exists():
        return {}
    out: dict[str, str] = {}
    for line in path.read_text(encoding="utf-8").splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key.strip()] = value.strip().replace(str(ROOT) + "/", "")
    return out


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


def evidence_passed(summary: dict[str, Any], browser: dict[str, Any], db: dict[str, str]) -> bool:
    handler = summary.get("go_tests", {}).get("handler", {})
    browser_routes = set(browser.get("downstream", {}).get("routes", []))
    required_routes = {"/decision-loop", "/review", "/audit", "/workbench"}
    required_db_minimums = {
        "position_count": 3,
        "snapshot_position_count": 3,
        "operation_confirmations_p84": 1,
        "position_transactions_p84": 1,
        "local_account_import_batches": 1,
        "local_account_corrections": 1,
        "portfolio_audit_events": 1,
    }
    db_passed = db.get("status") == "passed"
    for key, expected in required_db_minimums.items():
        try:
            db_passed = db_passed and int(db.get(key, "-1")) >= expected
        except ValueError:
            db_passed = False
    try:
        safety_passed = (
            int(db.get("forbidden_broker_order_push_tables", "-1")) == 0
            and int(db.get("auto_confirmation_rows", "-1")) == 0
        )
    except ValueError:
        safety_passed = False
    return (
        summary.get("status") == "passed"
        and browser.get("status") == "passed"
        and handler.get("status") == "passed"
        and required_routes.issubset(browser_routes)
        and db_passed
        and db.get("decision_p84_status") == "executed_manually"
        and safety_passed
    )


def p84_basis(requirement_id: str) -> str:
    if requirement_id == "REQ-02-033":
        return (
            "Fresh P84 real browser confirmation moved `decision_p84_pending` from pending to executed_manually; "
            "decision detail, decision-loop, review, audit, and SQLite `decision_records`/`operation_confirmations` readback agree."
        )
    if requirement_id == "REQ-11-002":
        return (
            "Fresh P84 real browser `/positions` journey manually entered cash, total assets, symbols, names, quantities, "
            "prices, buy reasons, and asset tags, then verified API and SQLite readback."
        )
    if requirement_id == "REQ-11-019":
        return (
            "Fresh P84 evidence shows pending advice did not alter account state until explicit user confirmation; after "
            "executed_manually confirmation, SQLite position transaction, latest portfolio snapshot, decision detail, and downstream UI readbacks changed."
        )
    return "Fresh P84 portfolio/confirmation evidence."


def p84_gap(requirement_id: str) -> str:
    if requirement_id in P84_UPGRADE_IDS:
        return "None for this P84 row; broader portfolio policy rows remain outside this upgrade."
    return (
        "P84 evaluated this row but did not upgrade it because the fresh evidence does not fully prove the complete row, "
        "such as core/satellite/cash target-policy enforcement, quarterly rebalance, sell-only/frozen-watch transitions, "
        "rule proposal confirmation, public-source collector readiness, or full product-goal effectiveness."
    )


def prior_status(row: dict[str, str]) -> str:
    return row.get("p83_status") or row.get("p82_status") or row.get("p81_status") or row.get("current_status") or "partial"


def row_with_p84(row: dict[str, str], passed: bool) -> dict[str, str]:
    requirement_id = row["requirement_id"]
    out = dict(row)
    if requirement_id in P84_PLAN_IDS:
        if passed and requirement_id in P84_UPGRADE_IDS:
            out.update({
                "p84_status": "real_pass",
                "p84_closure_basis": p84_basis(requirement_id),
                "p84_fresh_evidence_command": P84_UI_COMMAND + " && python3 scripts/p84_portfolio_confirmation_data_impact_closure.py --check",
                "p84_fresh_evidence_artifact": rel(P84_SUMMARY),
                "p84_remaining_gap": p84_gap(requirement_id),
                "p84_next_action": "Keep in P84 portfolio/confirmation regression and continue P85-P86 for broader non-P84 rows.",
            })
        elif passed:
            status = prior_status(row)
            if status not in {"real_pass", "scoped_pass", "reference_only"}:
                status = "partial"
            out.update({
                "p84_status": status,
                "p84_closure_basis": (
                    "P84 exercised adjacent real UI/API/SQLite portfolio behavior but did not produce complete row-specific proof "
                    "for this broader allocation, rule, source, master-knowledge, UX, or release-governance requirement."
                ),
                "p84_fresh_evidence_command": P84_UI_COMMAND + " && python3 scripts/p84_portfolio_confirmation_data_impact_closure.py --check",
                "p84_fresh_evidence_artifact": rel(P84_SUMMARY),
                "p84_remaining_gap": p84_gap(requirement_id),
                "p84_next_action": "Carry forward to P86 final integrated closure or a dedicated row-specific acceptance.",
            })
        else:
            out.update({
                "p84_status": "partial",
                "p84_closure_basis": "P84 evidence did not pass; this row remains non-real-pass.",
                "p84_fresh_evidence_command": P84_UI_COMMAND,
                "p84_fresh_evidence_artifact": rel(P84_SUMMARY),
                "p84_remaining_gap": "Fresh P84 UI/API/SQLite/Go evidence must pass before upgrade.",
                "p84_next_action": "Rerun P84 acceptance after fixing the failing evidence.",
            })
        return out

    out.update({
        "p84_status": prior_status(row),
        "p84_closure_basis": "No P84 upgrade; row is owned by P85-P86 or remains previously scoped/reference.",
        "p84_fresh_evidence_command": "N/A",
        "p84_fresh_evidence_artifact": "N/A",
        "p84_remaining_gap": row.get("p83_remaining_gap", row.get("p82_remaining_gap", "N/A")),
        "p84_next_action": row.get("p83_next_action", row.get("p82_next_action", "Continue assigned later batch.")),
    })
    return out


def write_matrix(header: list[str], rows: list[dict[str, str]]) -> None:
    output_header = header + P84_COLUMNS
    lines = [
        "# P84 Portfolio Confirmation Data Impact Matrix",
        "",
        f"> Generated: {datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace('+00:00', 'Z')}",
        f"> Source: `{rel(SOURCE_MATRIX)}`",
        "> Policy: P84 is a new evidence layer; it does not rewrite P75-P83 history.",
        "",
        "|" + "|".join(output_header) + "|",
        "|" + "|".join("---" for _ in output_header) + "|",
    ]
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(col, "")) for col in output_header) + "|")
    P84_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]], summary: dict[str, Any], db: dict[str, str]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p84_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in rows if row["requirement_id"] in P84_PLAN_IDS and row["p84_status"] == "real_pass"]
    deferred = sorted(P84_PLAN_IDS - set(upgraded))
    counts = Counter(row["p84_status"] for row in rows)
    lines = [
        "# P84 Portfolio Confirmation Data Impact Closure",
        "",
        "> Date: 2026-06-22",
        "> Change: `p84-portfolio-confirmation-data-impact-closure`",
        "> Conclusion: `release_ready_scoped_with_p84_portfolio_confirmation_progress`",
        "",
        "## Evidence Commands",
        "",
        "```bash",
        P84_UI_COMMAND,
        "python3 scripts/p84_portfolio_confirmation_data_impact_closure.py --check",
        "```",
        "",
        "## Result",
        "",
        f"- Total rows: {len(rows)}",
        f"- Full-release-required rows: {len(full_rows)}",
        f"- Full-release-required `real_pass` rows after P84: {len(full_real)}",
        f"- Remaining full-release-required non-`real_pass` rows: {len(full_rows) - len(full_real)}",
        f"- P84 evaluated rows: {len(P84_PLAN_IDS)}",
        f"- Newly upgraded by P84: {len(upgraded)}",
        f"- P84 evaluated but deferred rows: {len(deferred)}",
        f"- Matrix counts: {dict(sorted(counts.items()))}",
        "",
        "## Fresh Evidence",
        "",
        f"- Runtime summary: `{rel(P84_SUMMARY)}`",
        f"- Browser results: `{rel(P84_BROWSER)}`",
        f"- SQLite/readback log: `{rel(P84_DB_CHECK)}`",
        f"- Handler tests: `{rel(P84_ASSET_DIR / 'go-handler-tests.log')}`",
        f"- Screenshots: `{rel(P84_ASSET_DIR)}/p84-portfolio-after-actions.png`, `{rel(P84_ASSET_DIR)}/p84-decision-confirmed.png`, `{rel(P84_ASSET_DIR)}/p84-audit-readback.png`",
        "",
        "## SQLite Readback",
        "",
    ]
    for key in sorted(db):
        lines.append(f"- `{key}` = `{db[key]}`")
    lines.extend([
        "",
        "## Upgraded Rows",
        "",
    ])
    lines.extend(f"- `{rid}`" for rid in upgraded)
    lines.extend([
        "",
        "## Deferred Rows",
        "",
    ])
    lines.extend(f"- `{rid}`" for rid in deferred)
    lines.extend([
        "",
        "## Boundary",
        "",
        "P84 proves a real local portfolio and manual-confirmation data-impact path: browser UI operations, API readbacks, SQLite field checks, audit UI readbacks, downstream decision-loop/review/workbench readbacks, and no broker/order/external-push/auto-confirm persistence. P84 does not claim complete core/satellite/cash target allocation enforcement, quarterly rebalance execution, sell-only/frozen-watch transitions, public collector readiness, rule proposal application, broker sync, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, or return promises.",
        "",
        "P85-P86 remain required for the remaining full-release-required non-`real_pass` rows.",
        "",
        "## Machine Summary",
        "",
        "```json",
        json.dumps({
            "summary_status": summary.get("status"),
            "p84_rows": len(P84_PLAN_IDS),
            "new_real": len(upgraded),
            "deferred_rows": len(deferred),
            "remaining_full_release_required_non_real_pass": len(full_rows) - len(full_real),
        }, ensure_ascii=False, indent=2),
        "```",
    ])
    P84_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="validate existing P84 evidence without writing artifacts")
    args = parser.parse_args()

    header, source_rows = read_source_rows()
    ids = {row["requirement_id"] for row in source_rows}
    missing = sorted(P84_PLAN_IDS - ids)
    if missing:
        raise SystemExit(f"P84 planned rows missing from source matrix: {missing}")
    if len(P84_PLAN_IDS) != 35:
        raise SystemExit(f"P84 planned row count mismatch: {len(P84_PLAN_IDS)}")

    summary = read_json(P84_SUMMARY)
    browser = read_json(P84_BROWSER)
    db = read_log_kv(P84_DB_CHECK)
    passed = evidence_passed(summary, browser, db)
    output_rows = [row_with_p84(row, passed) for row in source_rows]
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p84_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in output_rows if row["requirement_id"] in P84_PLAN_IDS and row["p84_status"] == "real_pass"]
    remaining = len(full_rows) - len(full_real)

    if not args.check:
        write_matrix(header, output_rows)
        write_acceptance(output_rows, summary, db)
        enriched = dict(summary)
        enriched["p84_matrix"] = {
            "source_matrix": rel(SOURCE_MATRIX),
            "matrix": rel(P84_MATRIX),
            "acceptance": rel(P84_ACCEPTANCE),
            "evaluated_rows": len(P84_PLAN_IDS),
            "newly_upgraded_rows": len(upgraded),
            "newly_upgraded_requirement_ids": upgraded,
            "deferred_rows": len(P84_PLAN_IDS) - len(upgraded),
            "deferred_requirement_ids": sorted(P84_PLAN_IDS - set(upgraded)),
            "full_release_required_rows": len(full_rows),
            "full_release_required_real_pass_rows": len(full_real),
            "remaining_full_release_required_non_real_pass_rows": remaining,
            "conclusion": "release_ready_scoped_with_p84_portfolio_confirmation_progress",
        }
        P84_SUMMARY.write_text(json.dumps(enriched, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    if not passed:
        raise SystemExit("P84 portfolio confirmation evidence has not passed")
    if len(upgraded) != len(P84_UPGRADE_IDS):
        raise SystemExit(f"P84 upgraded row count mismatch: expected={len(P84_UPGRADE_IDS)} actual={len(upgraded)}")
    if remaining != 157:
        raise SystemExit(f"P84 remaining full-release-required mismatch: expected=157 actual={remaining}")

    print(f"p84_portfolio_confirmation:status=passed:new_real={len(upgraded)}:remaining_full={remaining}")


if __name__ == "__main__":
    main()
