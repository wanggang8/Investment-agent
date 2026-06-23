#!/usr/bin/env python3
"""Generate P83 governance traceability backfill artifacts."""

from __future__ import annotations

import argparse
import json
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p82-sop-action-ui-sqlite-matrix.md"
P83_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p83-governance-traceability-matrix.md"
P83_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p83-governance-traceability-backfill.md"
P83_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p83-governance-traceability"
P83_SUMMARY = P83_ASSET_DIR / "governance-traceability-summary.json"
P83_BROWSER = P83_ASSET_DIR / "browser-results.json"
P83_DB_CHECK = P83_ASSET_DIR / "db-readback-check.log"

P83_UI_COMMAND = (
    "P83_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability "
    "bash scripts/p83-governance-traceability-acceptance.sh"
)

P83_COLUMNS = [
    "p83_status",
    "p83_closure_basis",
    "p83_fresh_evidence_command",
    "p83_fresh_evidence_artifact",
    "p83_remaining_gap",
    "p83_next_action",
]

P83_PLAN_IDS = {
    "REQ-12-002",
    "REQ-12-003",
    "REQ-13-011",
    "REQ-14-004",
    "REQ-15-001",
    "REQ-15-002",
    "REQ-15-003",
    "REQ-15-004",
    "REQ-15-005",
    "REQ-15-007",
    "REQ-15-009",
    "REQ-16-001",
    "REQ-16-002",
    "REQ-16-005",
    "REQ-16-006",
    "REQ-16-007",
    "REQ-16-008",
    "REQ-16-009",
    "REQ-16-010",
    "REQ-16-011",
    "REQ-16-013",
    "REQ-16-014",
    "REQ-16-015",
    "REQ-16-019",
    "REQ-16-021",
    "REQ-16-023",
    "REQ-16-025",
    "REQ-16-026",
    "REQ-16-030",
    "REQ-16-031",
    "REQ-16-032",
    "REQ-16-034",
    "REQ-17-003",
    "REQ-17-005",
    "REQ-17-007",
    "REQ-17-008",
    "REQ-17-011",
    "REQ-17-012",
    "REQ-17-014",
    "REQ-17-017",
    "REQ-17-020",
    "REQ-17-022",
    "REQ-17-023",
}

P83_UPGRADE_IDS = {
    "REQ-12-002",
    "REQ-12-003",
    "REQ-13-011",
    "REQ-15-009",
    "REQ-16-026",
    "REQ-16-032",
    "REQ-17-017",
    "REQ-17-020",
    "REQ-17-022",
    "REQ-17-023",
}

REVIEW_ROWS = {"REQ-12-002", "REQ-12-003", "REQ-13-011", "REQ-15-009", "REQ-16-026", "REQ-17-017"}
OPS_ROWS = {"REQ-17-020", "REQ-17-022", "REQ-17-023"}


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
    go_tests = summary.get("go_tests", {})
    go_passed = bool(go_tests) and all(isinstance(item, dict) and item.get("status") == "passed" for item in go_tests.values())
    required_counts = {
        "review_decisions": 2,
        "review_confirmations": 2,
        "review_error_cases": 1,
        "review_rule_proposals": 2,
        "master_weight_proposals": 1,
        "quarterly_effect_tracking": 1,
        "review_notifications": 1,
        "review_audit_events": 2,
        "forbidden_broker_order_push_tables": 0,
        "active_rule_versions_from_p83": 0,
    }
    db_passed = db.get("status") == "passed"
    for key, expected in required_counts.items():
        try:
            actual = int(db.get(key, "-1"))
        except ValueError:
            db_passed = False
            continue
        if key in {"forbidden_broker_order_push_tables", "active_rule_versions_from_p83"}:
            db_passed = db_passed and actual == expected
        else:
            db_passed = db_passed and actual >= expected
    return summary.get("status") == "passed" and browser.get("status") == "passed" and go_passed and db_passed


def p83_basis(requirement_id: str) -> str:
    if requirement_id in REVIEW_ROWS:
        return (
            "Fresh P83 real browser UI, review API period checks, SQLite readback, and Go workflow/handler tests prove "
            "monthly/quarterly review, rule-effect, proposal, audit, notification, and master-weight governance behavior."
        )
    if requirement_id in OPS_ROWS:
        return (
            "Fresh P83 local-install UI readback plus focused cmd/agent release/diagnostic tests prove the local release, "
            "upgrade, diagnostic, and redaction governance behavior."
        )
    return (
        "Fresh P83 governance traceability evidence links the delivered roadmap capability to current UI/API behavior, "
        "focused Go tests, acceptance artifacts, and explicit release-boundary checks."
    )


def p83_gap(requirement_id: str) -> str:
    if requirement_id in OPS_ROWS:
        return "None for this local release/diagnostic governance row; P83 does not refresh the P76 package."
    return "None for this governance traceability row; broader original-requirement rows outside P83 remain in P84-P86."


def row_with_p83(row: dict[str, str], passed: bool) -> dict[str, str]:
    requirement_id = row["requirement_id"]
    out = dict(row)
    if requirement_id in P83_PLAN_IDS:
        if passed and requirement_id in P83_UPGRADE_IDS:
            out.update({
                "p83_status": "real_pass",
                "p83_closure_basis": p83_basis(requirement_id),
                "p83_fresh_evidence_command": P83_UI_COMMAND + " && python3 scripts/p83_governance_traceability_backfill.py --check",
                "p83_fresh_evidence_artifact": rel(P83_SUMMARY),
                "p83_remaining_gap": p83_gap(requirement_id),
                "p83_next_action": "Keep in P83 governance traceability regression and continue P84-P86 for non-P83 rows.",
            })
        elif passed:
            out.update({
                "p83_status": "partial",
                "p83_closure_basis": (
                    "P83 reviewed this row but did not upgrade it: fresh governance traceability evidence is not row-specific "
                    "enough to prove this broader implementation, analysis, knowledge/RAG, dashboard, or product-goal behavior."
                ),
                "p83_fresh_evidence_command": P83_UI_COMMAND + " && python3 scripts/p83_governance_traceability_backfill.py --check",
                "p83_fresh_evidence_artifact": rel(P83_SUMMARY),
                "p83_remaining_gap": (
                    "Needs dedicated P86 integrated real UI/API/SQLite/workflow evidence or another row-specific acceptance "
                    "before it can become real_pass."
                ),
                "p83_next_action": "Carry forward to P86 final integrated closure or a dedicated implementation-evidence batch.",
            })
        else:
            out.update({
                "p83_status": "partial",
                "p83_closure_basis": "P83 evidence did not pass; this row remains non-real-pass.",
                "p83_fresh_evidence_command": P83_UI_COMMAND,
                "p83_fresh_evidence_artifact": rel(P83_SUMMARY),
                "p83_remaining_gap": "Fresh P83 UI/API/SQLite/Go evidence must pass before upgrade.",
                "p83_next_action": "Rerun P83 acceptance after fixing the failing evidence.",
            })
        return out

    out.update({
        "p83_status": row.get("p82_status", row.get("p81_status", row.get("current_status", "partial"))),
        "p83_closure_basis": "No P83 upgrade; row is owned by P84-P86 or remains previously scoped/reference.",
        "p83_fresh_evidence_command": "N/A",
        "p83_fresh_evidence_artifact": "N/A",
        "p83_remaining_gap": row.get("p82_remaining_gap", row.get("p81_remaining_gap", "N/A")),
        "p83_next_action": row.get("p82_next_action", row.get("p81_next_action", "Continue assigned later batch.")),
    })
    return out


def write_matrix(header: list[str], rows: list[dict[str, str]]) -> None:
    output_header = header + P83_COLUMNS
    lines = [
        "# P83 Governance Traceability Matrix",
        "",
        f"> Generated: {datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace('+00:00', 'Z')}",
        f"> Source: `{rel(SOURCE_MATRIX)}`",
        "> Policy: P83 is a new evidence layer; it does not rewrite P75-P82 history.",
        "",
        "|" + "|".join(output_header) + "|",
        "|" + "|".join("---" for _ in output_header) + "|",
    ]
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(col, "")) for col in output_header) + "|")
    P83_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]], summary: dict[str, Any], db: dict[str, str]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p83_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in rows if row["requirement_id"] in P83_PLAN_IDS and row["p83_status"] == "real_pass"]
    counts = Counter(row["p83_status"] for row in rows)
    lines = [
        "# P83 Governance Traceability Backfill Acceptance",
        "",
        "> Date: 2026-06-22",
        "> Change: `p83-governance-traceability-backfill`",
        "> Conclusion: `release_ready_scoped_with_p83_governance_traceability_progress`",
        "",
        "## Evidence Commands",
        "",
        "```bash",
        P83_UI_COMMAND,
        "python3 scripts/p83_governance_traceability_backfill.py --check",
        "```",
        "",
        "## Result",
        "",
        f"- Total rows: {len(rows)}",
        f"- Full-release-required rows: {len(full_rows)}",
        f"- Full-release-required `real_pass` rows after P83: {len(full_real)}",
        f"- Remaining full-release-required non-`real_pass` rows: {len(full_rows) - len(full_real)}",
        f"- P83 evaluated rows: {len(P83_PLAN_IDS)}",
        f"- Newly upgraded by P83: {len(upgraded)}",
        f"- P83 evaluated but deferred rows: {len(P83_PLAN_IDS) - len(upgraded)}",
        f"- Matrix counts: {dict(sorted(counts.items()))}",
        "",
        "## Fresh Evidence",
        "",
        f"- Runtime summary: `{rel(P83_SUMMARY)}`",
        f"- Browser results: `{rel(P83_BROWSER)}`",
        f"- SQLite/readback log: `{rel(P83_DB_CHECK)}`",
        f"- Handler tests: `{rel(P83_ASSET_DIR / 'go-handler-tests.log')}`",
        f"- Workflow tests: `{rel(P83_ASSET_DIR / 'go-workflow-tests.log')}`",
        f"- Agent CLI tests: `{rel(P83_ASSET_DIR / 'go-agent-tests.log')}`",
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
        "## Boundary",
        "",
        "P83 evaluates 43 governance/review traceability candidate rows and upgrades only the rows directly backed by fresh UI/API/SQLite/Go evidence. Broader implementation, analysis, knowledge/RAG, dashboard, and product-goal rows remain partial for P86 or another row-specific acceptance. P83 does not refresh the P76 package, fabricate historical archives, perform physical second-machine acceptance, connect a broker, trade, create external push, automatically confirm user actions, or automatically apply rules.",
        "",
        "P84-P86 remain required for the remaining full-release-required non-`real_pass` rows.",
        "",
        "## Machine Summary",
        "",
        "```json",
        json.dumps({
            "summary_status": summary.get("status"),
            "p83_rows": len(P83_PLAN_IDS),
            "new_real": len(upgraded),
            "deferred_rows": len(P83_PLAN_IDS) - len(upgraded),
            "remaining_full_release_required_non_real_pass": len(full_rows) - len(full_real),
        }, ensure_ascii=False, indent=2),
        "```",
    ])
    P83_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="validate existing P83 evidence without writing artifacts")
    args = parser.parse_args()

    header, source_rows = read_source_rows()
    ids = {row["requirement_id"] for row in source_rows}
    missing = sorted(P83_PLAN_IDS - ids)
    if missing:
        raise SystemExit(f"P83 planned rows missing from source matrix: {missing}")
    if len(P83_PLAN_IDS) != 43:
        raise SystemExit(f"P83 planned row count mismatch: {len(P83_PLAN_IDS)}")

    summary = read_json(P83_SUMMARY)
    browser = read_json(P83_BROWSER)
    db = read_log_kv(P83_DB_CHECK)
    passed = evidence_passed(summary, browser, db)
    output_rows = [row_with_p83(row, passed) for row in source_rows]
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p83_status"] == "real_pass"]
    upgraded = [row["requirement_id"] for row in output_rows if row["requirement_id"] in P83_PLAN_IDS and row["p83_status"] == "real_pass"]
    remaining = len(full_rows) - len(full_real)

    if not args.check:
        write_matrix(header, output_rows)
        write_acceptance(output_rows, summary, db)
        enriched = dict(summary)
        enriched["p83_matrix"] = {
            "source_matrix": rel(SOURCE_MATRIX),
            "matrix": rel(P83_MATRIX),
            "acceptance": rel(P83_ACCEPTANCE),
            "evaluated_rows": len(P83_PLAN_IDS),
            "newly_upgraded_rows": len(upgraded),
            "newly_upgraded_requirement_ids": upgraded,
            "deferred_rows": len(P83_PLAN_IDS) - len(upgraded),
            "deferred_requirement_ids": sorted(P83_PLAN_IDS - set(upgraded)),
            "full_release_required_rows": len(full_rows),
            "full_release_required_real_pass_rows": len(full_real),
            "remaining_full_release_required_non_real_pass_rows": remaining,
            "conclusion": "release_ready_scoped_with_p83_governance_traceability_progress",
        }
        P83_SUMMARY.write_text(json.dumps(enriched, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    if not passed:
        raise SystemExit("P83 governance traceability evidence has not passed")
    if len(upgraded) != len(P83_UPGRADE_IDS):
        raise SystemExit(f"P83 upgraded row count mismatch: expected={len(P83_UPGRADE_IDS)} actual={len(upgraded)}")
    if remaining != 160:
        raise SystemExit(f"P83 remaining full-release-required mismatch: expected=160 actual={remaining}")

    print(f"p83_governance_traceability:status=passed:new_real={len(upgraded)}:remaining_full={remaining}")


if __name__ == "__main__":
    main()
