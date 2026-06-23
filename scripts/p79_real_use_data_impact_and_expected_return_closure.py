#!/usr/bin/env python3
"""Generate P79 real-use data-impact and expected-return closure artifacts."""

from __future__ import annotations

import argparse
import json
import re
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
P78_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p78-requirements-real-pass-batch-matrix.md"
P79_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md"
P79_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p79-real-use-data-impact-and-expected-return-closure.md"
P79_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p79"
P79_SUMMARY = P79_ASSET_DIR / "real-use-data-impact-summary.json"
P79_REAL_USER_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p79-real-user-fund"
P79_REAL_USER_DB_SUMMARY = P79_REAL_USER_DIR / "db-impact-summary.json"
P79_REAL_USER_BROWSER = P79_REAL_USER_DIR / "browser-results.json"
P79_NON510300_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p79-non-510300"
P79_NON510300_DB_SUMMARY = P79_NON510300_DIR / "non-510300-db-impact-summary.json"
P79_NON510300_BROWSER = P79_NON510300_DIR / "browser-results.json"

P72_UI_COMMAND = (
    "P72_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-21-p79-real-user-fund "
    "bash scripts/p72-real-user-fund-scenario-acceptance.sh"
)
P75_NON510300_COMMAND = (
    "P75_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-21-p79-non-510300 "
    "bash scripts/p75-non-510300-real-ui-journey.sh"
)

P79_COLUMNS = [
    "p79_status",
    "p79_closure_basis",
    "p79_fresh_evidence_command",
    "p79_fresh_evidence_artifact",
    "p79_remaining_gap",
    "p79_next_action",
]

PORTFOLIO_UPGRADE_IDS = {
    "REQ-04-019",
    "REQ-11-001",
    "REQ-11-003",
    "REQ-11-004",
    "REQ-11-006",
    "REQ-11-007",
    "REQ-11-008",
    "REQ-11-009",
    "REQ-11-010",
    "REQ-11-011",
    "REQ-11-012",
    "REQ-11-013",
    "REQ-11-014",
    "REQ-11-015",
    "REQ-11-016",
    "REQ-11-017",
    "REQ-11-020",
    "REQ-14-006",
    "REQ-16-003",
    "REQ-16-004",
    "REQ-16-017",
    "REQ-17-001",
    "REQ-17-002",
}

EXPECTED_RETURN_GUARD_IDS = {
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


def split_markdown_row(line: str) -> list[str]:
    stripped = line.rstrip("\n")
    if not stripped.startswith("|") or not stripped.endswith("|"):
        return []
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in stripped[1:-1]:
        if escaped:
            current.append(char)
            escaped = False
            continue
        if char == "\\":
            escaped = True
            continue
        if char == "|":
            cells.append("".join(current))
            current = []
            continue
        current.append(char)
    cells.append("".join(current))
    return [cell.strip() for cell in cells]


def escape_cell(value: object) -> str:
    text = str(value).replace("\n", " ").replace("\r", " ").strip()
    return text.replace("\\", "\\\\").replace("|", "\\|")


def redact_path(value: object) -> str:
    raw = str(value or "")
    if not raw:
        return raw
    path = Path(raw)
    if path.is_absolute():
        try:
            return str(path.relative_to(ROOT))
        except ValueError:
            return "[redacted-absolute-path]"
    return raw


def read_json(path: Path) -> dict:
    if not path.exists():
        return {"status": "missing", "path": str(path.relative_to(ROOT))}
    with path.open(encoding="utf-8") as fh:
        data = json.load(fh)
    return data if isinstance(data, dict) else {"status": "invalid", "path": str(path.relative_to(ROOT))}


def read_p78_rows() -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in P78_MATRIX.read_text(encoding="utf-8").splitlines():
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if not cells:
            continue
        if header is None:
            if cells[0] == "requirement_id":
                header = cells
            continue
        if set("".join(cells)) <= {"-", ":"}:
            continue
        if len(cells) != len(header):
            raise SystemExit(f"Invalid P78 matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("P78 matrix header not found")
    return header, rows


def sanitize_artifacts() -> None:
    roots = [P79_ASSET_DIR, P79_REAL_USER_DIR, P79_NON510300_DIR]
    for directory in roots:
        if not directory.exists():
            continue
        for path in directory.rglob("*"):
            if path.suffix.lower() not in {".json", ".log", ".txt", ".md"} or not path.is_file():
                continue
            text = path.read_text(encoding="utf-8")
            sanitized = text.replace(str(ROOT) + "/", "")
            if sanitized != text:
                path.write_text(sanitized, encoding="utf-8")


def missing_checks(payload: dict, section: str, required: list[str]) -> list[str]:
    checks = ((payload.get(section) or {}).get("checks") if isinstance(payload.get(section), dict) else None) or {}
    return [f"{section}.{key}" for key in required if checks.get(key) is not True]


def summarize_real_user_evidence() -> dict:
    db = read_json(P79_REAL_USER_DB_SUMMARY)
    browser = read_json(P79_REAL_USER_BROWSER)
    if db.get("sqlite_path"):
        db["sqlite_path"] = redact_path(db["sqlite_path"])
    required_counts = {
        "portfolio_snapshots": ["portfolio_counts", "portfolio_snapshots", 5],
        "position_snapshots": ["portfolio_counts", "position_snapshots", 5],
        "local_account_import_batches_committed": ["portfolio_counts", "local_account_import_batches_committed", 1],
        "local_account_corrections": ["portfolio_counts", "local_account_corrections", 1],
        "operation_confirmations_total": ["portfolio_counts", "operation_confirmations_total", 2],
        "position_transactions_total": ["portfolio_counts", "position_transactions_total", 2],
        "manual_daily_reports": ["daily_counts", "manual_daily_reports", 1],
        "risk_alerts": ["daily_counts", "risk_alerts", 1],
        "notifications": ["daily_counts", "notifications", 1],
    }
    failures: list[str] = []
    if browser.get("status") in {"missing", "invalid"} or not P79_REAL_USER_BROWSER.exists():
        failures.append("browser_results")
    for label, (section, key, minimum) in required_counts.items():
        actual = ((db.get(section) or {}).get(key) if isinstance(db.get(section), dict) else None)
        if actual is None or int(actual) < minimum:
            failures.append(f"{label}<{minimum}")
    decision = db.get("decision") or {}
    if decision.get("workflow_status") != "completed":
        failures.append("decision.workflow_status")
    if decision.get("confirmation_status") != "executed_manually":
        failures.append("decision.confirmation_status")
    if (decision.get("analyst_report_count") or 0) < 1:
        failures.append("decision.analyst_report_count")
    if db.get("forbidden_tables"):
        failures.append("forbidden_tables")
    if browser.get("unexpected_failed_api_responses"):
        failures.append("unexpected_failed_api_responses")
    if browser.get("page_errors"):
        failures.append("page_errors")
    if browser.get("console_errors"):
        failures.append("console_errors")
    field_failures: list[str] = []
    field_failures.extend(missing_checks(db, "position_field_readback", ["symbol", "name", "quantity", "cost_price", "buy_reason", "asset_tag", "position_state"]))
    field_failures.extend(missing_checks(db, "operation_confirmation_readback", ["decision_linked_manual_execution", "offline_local_transaction", "quantity_and_price"]))
    field_failures.extend(missing_checks(db, "transaction_readback", ["transaction_rows", "symbol", "quantity_and_price", "before_after_state"]))
    table_readback = db.get("table_readback") or {}
    for key in ["portfolio_snapshots_have_source", "position_snapshots_have_required_fields", "decision_record_has_snapshot_refs", "evidence_refs_for_decision", "audit_events_include_user_confirm"]:
        if table_readback.get(key) is not True:
            field_failures.append(f"table_readback.{key}")
    if field_failures:
        failures.extend(field_failures)
    return {
        "status": "passed" if db.get("status") == "passed" and not failures else "failed",
        "db_summary": str(P79_REAL_USER_DB_SUMMARY.relative_to(ROOT)),
        "browser_results": str(P79_REAL_USER_BROWSER.relative_to(ROOT)),
        "failures": failures,
        "decision": decision,
        "portfolio_counts": db.get("portfolio_counts") or {},
        "daily_counts": db.get("daily_counts") or {},
        "field_level_evidence_failures": field_failures,
        "forbidden_tables": db.get("forbidden_tables") or [],
    }


def summarize_non510300_evidence() -> dict:
    db = read_json(P79_NON510300_DB_SUMMARY)
    browser = read_json(P79_NON510300_BROWSER)
    if db.get("sqlite_path"):
        db["sqlite_path"] = redact_path(db["sqlite_path"])
    decision = db.get("decision") or {}
    failures: list[str] = []
    if db.get("status") != "passed":
        failures.append("db_status")
    if decision.get("workflow_status") != "completed":
        failures.append("decision.workflow_status")
    if db.get("symbol") != "159915" or browser.get("symbol") != "159915":
        failures.append("symbol")
    source_verification = db.get("source_verification") or {}
    if source_verification.get("verification_status") != "satisfied":
        failures.append("source_verification.status")
    if (source_verification.get("high_grade_independent_source_count") or 0) < 2:
        failures.append("source_verification.high_grade_independent_source_count")
    if browser.get("status") != "passed":
        failures.append("browser_status")
    field_failures: list[str] = []
    position = db.get("position") or {}
    if position.get("symbol") != "159915":
        field_failures.append("position.symbol")
    if position.get("name") != "创业板ETF":
        field_failures.append("position.name")
    if (position.get("quantity") or 0) <= 0:
        field_failures.append("position.quantity")
    if (position.get("current_price") or 0) <= 0 or (position.get("market_value") or 0) <= 0:
        field_failures.append("position.market_value")
    if "创业板ETF" not in str(position.get("buy_reason") or ""):
        field_failures.append("position.buy_reason")
    if position.get("asset_tag") != "satellite":
        field_failures.append("position.asset_tag")
    if field_failures:
        failures.extend(field_failures)
    return {
        "status": "passed" if not failures else "failed",
        "db_summary": str(P79_NON510300_DB_SUMMARY.relative_to(ROOT)),
        "browser_results": str(P79_NON510300_BROWSER.relative_to(ROOT)),
        "failures": failures,
        "decision": decision,
        "symbol": db.get("symbol"),
        "tracked_index_symbol": db.get("tracked_index_symbol"),
        "source_verification": source_verification,
        "field_level_evidence_failures": field_failures,
    }


def p79_values(row: dict[str, str], evidence_ok: bool) -> dict[str, str]:
    rid = row["requirement_id"]
    if row["p78_status"] == "real_pass":
        return {
            "p79_status": "real_pass",
            "p79_closure_basis": "Carried forward from P78/P77 real_pass; P79 does not rewrite prior evidence.",
            "p79_fresh_evidence_command": "N/A",
            "p79_fresh_evidence_artifact": row.get("fresh_evidence_artifact", "N/A"),
            "p79_remaining_gap": "None for row already accepted before P79.",
            "p79_next_action": "Keep covered by future regression evidence.",
        }
    if rid in PORTFOLIO_UPGRADE_IDS and evidence_ok:
        return {
            "p79_status": "real_pass",
            "p79_closure_basis": "Fresh P79 real UI portfolio/confirmation journeys plus SQLite readback cover the exact local-account data-impact behavior.",
            "p79_fresh_evidence_command": f"{P72_UI_COMMAND}; {P75_NON510300_COMMAND}",
            "p79_fresh_evidence_artifact": f"{P79_SUMMARY.relative_to(ROOT)}",
            "p79_remaining_gap": "None for this row's local-account data-impact claim.",
            "p79_next_action": "Keep in real UI data-impact regression.",
        }
    if rid in PORTFOLIO_UPGRADE_IDS and not evidence_ok:
        return {
            "p79_status": row["p78_status"],
            "p79_closure_basis": "P79 upgrade candidate, but fresh evidence is not complete.",
            "p79_fresh_evidence_command": f"{P72_UI_COMMAND}; {P75_NON510300_COMMAND}",
            "p79_fresh_evidence_artifact": f"{P79_SUMMARY.relative_to(ROOT)}",
            "p79_remaining_gap": "Fresh P79 real UI data-impact evidence did not pass.",
            "p79_next_action": "Fix evidence failures and rerun P79 checker.",
        }
    if rid in EXPECTED_RETURN_GUARD_IDS:
        return {
            "p79_status": row["p78_status"],
            "p79_closure_basis": "No P79 upgrade; expected-return field-level UI/readback evidence is still incomplete.",
            "p79_fresh_evidence_command": "N/A",
            "p79_fresh_evidence_artifact": "N/A",
            "p79_remaining_gap": "Needs direct UI/readback proof for probabilities, scenario ranges, sell triggers, valuation fields, sample count/window/screening, provenance, and disclaimer as applicable.",
            "p79_next_action": "Run a dedicated expected-return available/degraded scenario UI matrix before upgrading.",
        }
    return {
        "p79_status": row["p78_status"],
        "p79_closure_basis": "No P79 upgrade; classified for a later closure batch.",
        "p79_fresh_evidence_command": "N/A",
        "p79_fresh_evidence_artifact": "N/A",
        "p79_remaining_gap": row.get("remaining_gap", "Needs post-P79 evidence."),
        "p79_next_action": row.get("next_action", "Create follow-up implementation or acceptance work."),
    }


def write_matrix(header: list[str], rows: list[dict[str, str]], evidence_ok: bool) -> list[dict[str, str]]:
    output_rows: list[dict[str, str]] = []
    columns = header + P79_COLUMNS
    for row in rows:
        merged = dict(row)
        merged.update(p79_values(row, evidence_ok))
        output_rows.append(merged)
    counts = Counter(row["p79_status"] for row in output_rows)
    full_rows = [row for row in output_rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p79_status"] == "real_pass"]
    remaining = len(full_rows) - len(full_real)
    lines = [
        "# P79 Real-Use Data-Impact And Expected-Return Matrix",
        "",
        "> Generated: 2026-06-21",
        f"> Source: `{P78_MATRIX.relative_to(ROOT)}`",
        "> Policy: P79 is a new evidence layer; it does not rewrite P75, P77, or P78 history.",
        "",
        "## Status Summary",
        "",
    ]
    for status, count in sorted(counts.items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend([
        "",
        f"- full_release_required rows: {len(full_rows)}",
        f"- full_release_required real_pass rows: {len(full_real)}",
        f"- remaining full_release_required non-real-pass rows: {remaining}",
        "- conclusion: `release_ready_scoped_with_p79_real_use_data_impact_progress`",
        "",
        "## Atomic Requirement Batch Rows",
        "",
        "|" + "|".join(columns) + "|",
        "|" + "|".join(["---"] * len(columns)) + "|",
    ])
    for row in output_rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in columns) + "|")
    P79_MATRIX.parent.mkdir(parents=True, exist_ok=True)
    P79_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return output_rows


def write_summary(rows: list[dict[str, str]], real_user: dict, non510300: dict) -> dict:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p79_status"] == "real_pass"]
    upgraded = [
        row["requirement_id"]
        for row in rows
        if row["p79_status"] == "real_pass" and row.get("p78_status") != "real_pass"
    ]
    payload = {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "source_matrix": str(P78_MATRIX.relative_to(ROOT)),
        "matrix": str(P79_MATRIX.relative_to(ROOT)),
        "acceptance": str(P79_ACCEPTANCE.relative_to(ROOT)),
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "newly_upgraded_rows": len(upgraded),
        "newly_upgraded_requirement_ids": upgraded,
        "conclusion": "release_ready_scoped_with_p79_real_use_data_impact_progress",
        "p72_real_user_fund": real_user,
        "p75_non_510300": non510300,
        "not_claimed": [
            "full original-requirement pass",
            "P79 evidence inside any existing P76 distribution archive",
            "broker connectivity or automatic trading",
            "expected-return field-level real_pass beyond evidence",
            "monthly attribution real_pass from P79",
            "future return or provider availability",
        ],
    }
    P79_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P79_SUMMARY.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return payload


def write_acceptance(summary: dict) -> None:
    upgraded = summary["newly_upgraded_requirement_ids"]
    lines = [
        "# P79 Real-Use Data-Impact And Expected-Return Closure Acceptance",
        "",
        "> Date: 2026-06-21",
        "> Change: `p79-real-use-data-impact-and-expected-return-closure`",
        "> Conclusion: `release_ready_scoped_with_p79_real_use_data_impact_progress`",
        "",
        "## Summary",
        "",
        f"- Source matrix: `{P78_MATRIX.relative_to(ROOT)}`",
        f"- P79 matrix: `{P79_MATRIX.relative_to(ROOT)}`",
        f"- Summary JSON: `{P79_SUMMARY.relative_to(ROOT)}`",
        f"- Full-release-required rows: {summary['full_release_required_rows']}",
        f"- Full-release-required `real_pass` rows after P79: {summary['full_release_required_real_pass_rows']}",
        f"- Remaining full-release-required non-`real_pass` rows: {summary['remaining_full_release_required_non_real_pass_rows']}",
        f"- Newly upgraded by P79: {summary['newly_upgraded_rows']}",
        "",
        "## P79 Upgrades",
        "",
    ]
    lines.extend(f"- `{rid}`" for rid in upgraded)
    lines.extend([
        "",
        "## Fresh Evidence",
        "",
        f"- P72 real-user fund scenario rerun: `{P79_REAL_USER_DIR.relative_to(ROOT)}`",
        f"- P75 accepted-local non-`510300` rerun: `{P79_NON510300_DIR.relative_to(ROOT)}`",
        f"- P79 summary/readback: `{P79_SUMMARY.relative_to(ROOT)}`",
        "",
        "Commands:",
        "",
        "```bash",
        P72_UI_COMMAND,
        P75_NON510300_COMMAND,
        "python3 scripts/p79_real_use_data_impact_and_expected_return_closure.py --check",
        "```",
        "",
        "## Expected-Return Remaining Gap",
        "",
        "P79 does not upgrade broad expected-return probability/scenario rows. Those rows still require direct UI/readback proof for available and degraded precision states, scenario ranges, sell-evaluation triggers, valuation fields, sample count/window/screening, source/provenance fields, and non-trading disclaimers.",
        "",
        "P79 also does not upgrade broad monthly attribution rows such as `REQ-14-005`; P79 proves daily/local account snapshot readback and confirmation data impact, not monthly attribution completeness.",
        "",
        "## Boundaries",
        "",
        "- P79 does not rewrite P75, P77, or P78 historical matrices.",
        "- P79 does not refresh the P76 package; a separate package refresh is required before claiming distribution archives include P79 materials.",
        "- P79 does not claim full original-requirement pass while any full-release-required row remains non-`real_pass`.",
        "- P79 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.",
    ])
    P79_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def validate_rows(rows: list[dict[str, str]], evidence_ok: bool) -> None:
    upgraded = {row["requirement_id"] for row in rows if row["p79_status"] == "real_pass" and row.get("p78_status") != "real_pass"}
    invalid = upgraded - PORTFOLIO_UPGRADE_IDS
    if invalid:
        raise SystemExit(f"P79 invalid upgraded rows: {sorted(invalid)}")
    if PORTFOLIO_UPGRADE_IDS & upgraded and not evidence_ok:
        raise SystemExit("P79 portfolio rows upgraded without passing evidence")
    overbroad = EXPECTED_RETURN_GUARD_IDS & upgraded
    if overbroad:
        raise SystemExit(f"P79 expected-return rows upgraded without field-level proof: {sorted(overbroad)}")
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    if any(row["p79_status"] != "real_pass" for row in full_rows):
        return
    raise SystemExit("P79 unexpectedly claims all full-release-required rows are real_pass; review full-pass claim gate")


def validate_claim_text() -> None:
    scan_files = [
        P79_MATRIX,
        P79_ACCEPTANCE,
        P79_SUMMARY,
        ROOT / "docs" / "release" / "README.md",
        ROOT / "docs" / "release" / "release-candidate-2026-06-18.md",
        ROOT / "docs" / "release" / "release-handoff-2026-06-18.md",
    ]
    for directory in [P79_ASSET_DIR, P79_REAL_USER_DIR, P79_NON510300_DIR]:
        if directory.exists():
            scan_files.extend(path for path in directory.rglob("*") if path.suffix.lower() in {".json", ".log", ".txt", ".md"})
    forbidden = [
        "release_ready_full_requirements_traceable",
        "status is full original-requirement pass",
        "conclusion is full original-requirement pass",
        "P76 package includes P79",
        "P76 archive includes P79",
        "monthly attribution is real_pass",
        "REQ-14-005 is real_pass",
        "automatic trading enabled",
        "one-click trading enabled",
        "guaranteed return",
        "DEEPSEEK_API_KEY",
        "OPENAI_API_KEY",
        "\"prompt_messages\"",
        "\"raw_request\"",
        "\"raw_response\"",
        "\"raw_payload\"",
    ]
    for path in scan_files:
        if not path.exists() or not path.is_file():
            continue
        text = path.read_text(encoding="utf-8", errors="ignore")
        if "/Users/" in text:
            raise SystemExit(f"P79 private absolute path leaked in {path.relative_to(ROOT)}")
        if re.search(r"sk-[A-Za-z0-9_-]{12,}", text):
            raise SystemExit(f"P79 potential API key leaked in {path.relative_to(ROOT)}")
        for phrase in forbidden:
            if phrase in text:
                raise SystemExit(f"P79 overbroad claim '{phrase}' found in {path.relative_to(ROOT)}")


def run(check: bool) -> None:
    sanitize_artifacts()
    header, source_rows = read_p78_rows()
    real_user = summarize_real_user_evidence()
    non510300 = summarize_non510300_evidence()
    evidence_ok = real_user["status"] == "passed" and non510300["status"] == "passed"
    rows = write_matrix(header, source_rows, evidence_ok)
    summary = write_summary(rows, real_user, non510300)
    write_acceptance(summary)
    sanitize_artifacts()
    if check:
        if real_user["status"] != "passed":
            raise SystemExit(f"P79 real-user evidence failed: {real_user['failures']}")
        if non510300["status"] != "passed":
            raise SystemExit(f"P79 non-510300 evidence failed: {non510300['failures']}")
        validate_rows(rows, evidence_ok)
        validate_claim_text()
    print(
        "p79_real_use_data_impact:"
        f"rows={len(rows)}:"
        f"real_pass={summary['full_release_required_real_pass_rows']}:"
        f"new={summary['newly_upgraded_rows']}:"
        f"remaining_full={summary['remaining_full_release_required_non_real_pass_rows']}:"
        f"conclusion={summary['conclusion']}"
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true")
    args = parser.parse_args()
    run(check=args.check)


if __name__ == "__main__":
    main()
