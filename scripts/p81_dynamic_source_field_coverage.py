#!/usr/bin/env python3
"""Generate P81 dynamic source field coverage artifacts from the P80 matrix."""

from __future__ import annotations

import argparse
import json
import re
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
P80_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p80-review-audit-governance-matrix.md"
P81_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p81-dynamic-source-field-coverage-matrix.md"
P81_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p81-dynamic-source-field-coverage.md"
P81_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p81-dynamic-source-field-coverage"
P81_BROWSER = P81_ASSET_DIR / "browser-results.json"
P81_DB_SUMMARY = P81_ASSET_DIR / "non-510300-db-impact-summary.json"
P81_GO_TEST_LOG = P81_ASSET_DIR / "dynamic-source-go-test.log"
P81_SUMMARY = P81_ASSET_DIR / "dynamic-source-field-coverage-summary.json"

P81_UI_COMMAND = (
    "P75_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage "
    "bash scripts/p75-non-510300-real-ui-journey.sh"
)
P81_GO_TEST_COMMAND = (
    "go test -v ./cmd/agent -run TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness -count=1"
)
P81_CHECK_COMMAND = "python3 scripts/p81_dynamic_source_field_coverage.py --check"

P81_COLUMNS = [
    "p81_status",
    "p81_closure_basis",
    "p81_fresh_evidence_command",
    "p81_fresh_evidence_artifact",
    "p81_remaining_gap",
    "p81_next_action",
]

P81_UPGRADE_IDS = {
    "REQ-02-003",
    "REQ-02-009",
    "REQ-02-015",
    "REQ-02-016",
    "REQ-02-023",
    "REQ-02-027",
    "REQ-02-028",
    "REQ-02-030",
    "REQ-04-001",
    "REQ-04-002",
    "REQ-04-004",
    "REQ-04-006",
    "REQ-04-009",
    "REQ-04-010",
    "REQ-04-011",
    "REQ-04-012",
    "REQ-04-013",
    "REQ-04-014",
    "REQ-04-015",
    "REQ-04-017",
    "REQ-04-018",
    "REQ-04-021",
    "REQ-04-022",
    "REQ-04-023",
    "REQ-04-024",
    "REQ-04-026",
    "REQ-04-027",
    "REQ-05-001",
    "REQ-05-002",
    "REQ-05-006",
    "REQ-05-007",
    "REQ-05-008",
    "REQ-05-009",
    "REQ-05-011",
    "REQ-05-012",
    "REQ-05-013",
    "REQ-05-014",
    "REQ-05-015",
    "REQ-05-016",
    "REQ-05-017",
    "REQ-05-018",
    "REQ-05-019",
    "REQ-05-020",
    "REQ-06-001",
    "REQ-06-009",
    "REQ-07-013",
    "REQ-07-014",
    "REQ-14-001",
    "REQ-14-002",
    "REQ-14-003",
    "REQ-15-006",
    "REQ-15-008",
    "REQ-16-012",
    "REQ-16-018",
    "REQ-16-020",
    "REQ-17-006",
    "REQ-17-009",
    "REQ-17-013",
    "REQ-17-021",
}

READY_CATEGORIES = {
    "symbol_profile": ["159915"],
    "fund_profile": ["159915"],
    "tracked_index": ["399006"],
    "market_price": ["159915"],
    "valuation_percentiles": ["399006"],
    "liquidity": ["159915"],
    "sentiment_proxy": ["159915"],
    "rag_index": ["159915"],
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


def rel(path: Path) -> str:
    return str(path.relative_to(ROOT))


def read_json(path: Path) -> dict[str, Any]:
    if not path.exists():
        return {"status": "missing", "path": rel(path)}
    with path.open(encoding="utf-8") as fh:
        data = json.load(fh)
    return data if isinstance(data, dict) else {"status": "invalid", "path": rel(path)}


def sanitize_decision_for_artifact(payload: dict[str, Any]) -> bool:
    decision = payload.get("decision")
    if not isinstance(decision, dict):
        return False
    reports = decision.get("analyst_reports")
    if not isinstance(reports, list):
        return False
    changed = False
    sanitized_reports: list[dict[str, Any]] = []
    for report in reports:
        if not isinstance(report, dict):
            sanitized_reports.append(report)
            continue
        conclusion = str(report.get("conclusion") or "").replace("\n", " ").strip()
        sanitized = {
            "agent_name": report.get("agent_name"),
            "model": report.get("model"),
            "prompt_version": report.get("prompt_version"),
            "parse_status": report.get("parse_status"),
            "quality_status": report.get("quality_status"),
            "confidence": report.get("confidence"),
            "evidence_ids": report.get("evidence_ids"),
            "input_summary": report.get("input_summary"),
            "output_summary": report.get("output_summary"),
            "conclusion_preview": (conclusion[:180] + "...") if len(conclusion) > 180 else conclusion,
        }
        if "conclusion" in report or sanitized != report:
            changed = True
        sanitized_reports.append(sanitized)
    if changed:
        decision["analyst_reports"] = sanitized_reports
    return changed


def artifact_has_raw_llm_conclusion(payload: dict[str, Any]) -> bool:
    decision = payload.get("decision")
    if not isinstance(decision, dict):
        return False
    reports = decision.get("analyst_reports")
    if not isinstance(reports, list):
        return False
    for report in reports:
        if isinstance(report, dict) and len(str(report.get("conclusion") or "")) > 180:
            return True
    return False


def sanitize_artifacts() -> None:
    if not P81_ASSET_DIR.exists():
        return
    for path in P81_ASSET_DIR.rglob("*"):
        if not path.is_file() or path.suffix.lower() not in {".json", ".log", ".txt", ".md"}:
            continue
        text = path.read_text(encoding="utf-8")
        sanitized = text.replace(str(ROOT) + "/", "")
        if path == P81_BROWSER:
            payload = json.loads(sanitized)
            if isinstance(payload, dict) and sanitize_decision_for_artifact(payload):
                sanitized = json.dumps(payload, ensure_ascii=False, indent=2) + "\n"
        if sanitized != text:
            path.write_text(sanitized, encoding="utf-8")


def read_p80_rows() -> tuple[list[str], list[dict[str, str]]]:
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in P80_MATRIX.read_text(encoding="utf-8").splitlines():
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
            raise SystemExit(f"Invalid P80 matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("P80 matrix header not found")
    return header, rows


def readiness_dependency(browser: dict[str, Any], category: str) -> dict[str, Any]:
    readiness = browser.get("readiness") if isinstance(browser.get("readiness"), dict) else {}
    for item in readiness.get("data_dependencies") or []:
        if isinstance(item, dict) and item.get("category") == category:
            return item
    return {}


def verify_evidence(*, sanitize: bool) -> dict[str, Any]:
    if sanitize:
        sanitize_artifacts()
    browser = read_json(P81_BROWSER)
    db = read_json(P81_DB_SUMMARY)
    go_log = P81_GO_TEST_LOG.read_text(encoding="utf-8") if P81_GO_TEST_LOG.exists() else ""
    failures: list[str] = []

    if browser.get("status") != "passed":
        failures.append("browser.status")
    if artifact_has_raw_llm_conclusion(browser):
        failures.append("browser.raw_llm_conclusion")
    if browser.get("symbol") != "159915" or browser.get("tracked_index_symbol") != "399006":
        failures.append("browser.symbol_binding")
    readiness = browser.get("readiness") if isinstance(browser.get("readiness"), dict) else {}
    if readiness.get("status") != "ready":
        failures.append("readiness.status")
    if ((readiness.get("symbol_profile") or {}).get("tracked_index_symbol")) != "399006":
        failures.append("readiness.tracked_index_symbol")
    for category, symbols in READY_CATEGORIES.items():
        dep = readiness_dependency(browser, category)
        if dep.get("status") != "ready":
            failures.append(f"readiness.{category}.status")
        if sorted(dep.get("affected_symbols") or []) != sorted(symbols):
            failures.append(f"readiness.{category}.affected_symbols")
        if not dep.get("request_id"):
            failures.append(f"readiness.{category}.request_id")
    for category in ["active_rule", "formal_evidence", "llm_context"]:
        dep = readiness_dependency(browser, category)
        if dep.get("status") != "ready":
            failures.append(f"readiness.{category}.status")

    decision = browser.get("decision") if isinstance(browser.get("decision"), dict) else {}
    if decision.get("workflow_status") != "completed" or (decision.get("analyst_reports") and len(decision.get("analyst_reports") or []) < 3):
        failures.append("browser.decision")

    if db.get("status") != "passed":
        failures.append("db.status")
    if db.get("symbol") != "159915" or db.get("tracked_index_symbol") != "399006":
        failures.append("db.symbol_binding")
    market = db.get("market") if isinstance(db.get("market"), dict) else {}
    if market.get("request_id") in {"", None}:
        failures.append("db.market.request_id")
    for key, expected_symbol in [("tracked_index_health", "399006"), ("valuation_health", "399006")]:
        item = market.get(key) if isinstance(market.get(key), dict) else {}
        if item.get("freshness") != "fresh":
            failures.append(f"db.market.{key}.freshness")
        if item.get("request_id") != market.get("request_id"):
            failures.append(f"db.market.{key}.request_id")
        if item.get("affected_symbols") != [expected_symbol]:
            failures.append(f"db.market.{key}.affected_symbols")
    source_verification = db.get("source_verification") if isinstance(db.get("source_verification"), dict) else {}
    if source_verification.get("verification_status") != "satisfied":
        failures.append("db.source_verification.status")
    if int(source_verification.get("high_grade_independent_source_count") or 0) < 2:
        failures.append("db.source_verification.high_grade")
    knowledge_counts = db.get("knowledge_counts") if isinstance(db.get("knowledge_counts"), dict) else {}
    for key in ["intelligence_items", "intelligence_summary", "rag_chunks_indexed"]:
        if int(knowledge_counts.get(key) or 0) < 2:
            failures.append(f"db.knowledge_counts.{key}")
    audit_counts = db.get("audit_counts") if isinstance(db.get("audit_counts"), dict) else {}
    for key in ["market_refresh", "public_evidence_command", "public_evidence_ingestion", "consultation_decision"]:
        if int(audit_counts.get(key) or 0) < 1:
            failures.append(f"db.audit_counts.{key}")
    if db.get("forbidden_tables"):
        failures.append("db.forbidden_tables")

    if "PASS: TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness" not in go_log or not re.search(r"\bPASS\b", go_log):
        failures.append("go_test.dynamic_source")

    return {
        "status": "passed" if not failures else "failed",
        "failures": failures,
        "browser_results": rel(P81_BROWSER),
        "db_summary": rel(P81_DB_SUMMARY),
        "go_test_log": rel(P81_GO_TEST_LOG),
        "symbol": "159915",
        "tracked_index_symbol": "399006",
        "ready_categories": sorted(READY_CATEGORIES),
        "source_verification": source_verification,
        "knowledge_counts": knowledge_counts,
        "audit_counts": audit_counts,
        "forbidden_tables": db.get("forbidden_tables") or [],
    }


def classify_row(row: dict[str, str], evidence: dict[str, Any]) -> dict[str, str]:
    rid = row["requirement_id"]
    p80_status = row.get("p80_status", "")
    if p80_status == "real_pass":
        return {
            "p81_status": "real_pass",
            "p81_closure_basis": "Carried forward from P80/P79/P78/P77 real_pass; P81 does not rewrite prior evidence.",
            "p81_fresh_evidence_command": "N/A",
            "p81_fresh_evidence_artifact": row.get("p80_fresh_evidence_artifact") or "N/A",
            "p81_remaining_gap": "None for row already accepted before P81.",
            "p81_next_action": "Keep covered by future regression evidence.",
        }
    if rid in P81_UPGRADE_IDS and evidence["status"] == "passed":
        artifact = "; ".join([evidence["browser_results"], evidence["db_summary"], evidence["go_test_log"], rel(P81_SUMMARY)])
        return {
            "p81_status": "real_pass",
            "p81_closure_basis": "Fresh P81 non-510300 real UI/API/SQLite/readback evidence proves dynamic source field binding, readiness, provenance, source health, RAG index, LLM context use, auditability, and forbidden-capability absence for this data-source row.",
            "p81_fresh_evidence_command": f"{P81_GO_TEST_COMMAND} && {P81_UI_COMMAND} && {P81_CHECK_COMMAND}",
            "p81_fresh_evidence_artifact": artifact,
            "p81_remaining_gap": "None for this dynamic source field coverage row.",
            "p81_next_action": "Keep in dynamic source/readiness regression; do not broaden to paid/login/Level2/high-frequency sources.",
        }
    return {
        "p81_status": p80_status,
        "p81_closure_basis": "No P81 upgrade; row is owned by a later P82-P86 execution batch or remains previously scoped/reference.",
        "p81_fresh_evidence_command": "N/A",
        "p81_fresh_evidence_artifact": row.get("p80_fresh_evidence_artifact") or "N/A",
        "p81_remaining_gap": row.get("p80_remaining_gap") or "Needs later batch evidence.",
        "p81_next_action": row.get("p80_next_action") or "Execute the assigned later P82-P86 closure batch.",
    }


def write_matrix(header: list[str], rows: list[dict[str, str]]) -> None:
    P81_MATRIX.parent.mkdir(parents=True, exist_ok=True)
    out_header = header + P81_COLUMNS
    lines = [
        "# P81 Dynamic Source Field Coverage Matrix",
        "",
        f"> Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
        f"> Source: `{rel(P80_MATRIX)}`",
        "> Policy: P81 is a new evidence layer; it does not rewrite P75-P80 history.",
        "",
        "## Status Summary",
        "",
    ]
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p81_status"] == "real_pass"]
    lines.extend(
        [
            f"- total rows: {len(rows)}",
            f"- full_release_required rows: {len(full_rows)}",
            f"- full_release_required real_pass rows: {len(full_real)}",
            f"- remaining full_release_required non-real-pass rows: {len(full_rows) - len(full_real)}",
            f"- new P81 real_pass rows: {len([row for row in rows if row['requirement_id'] in P81_UPGRADE_IDS and row['p81_status'] == 'real_pass'])}",
            "",
            "## Remaining Non-Real-Pass By Remediation Group",
            "",
        ]
    )
    for group, count in sorted(Counter(row.get("remediation_group") for row in full_rows if row["p81_status"] != "real_pass").items()):
        lines.append(f"- {group}: {count}")
    lines.extend(["", "## Matrix", "", "|" + "|".join(out_header) + "|", "|" + "|".join(["---"] * len(out_header)) + "|"])
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in out_header) + "|")
    P81_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]], evidence: dict[str, Any]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p81_status"] == "real_pass"]
    new_real = sorted(row["requirement_id"] for row in rows if row["requirement_id"] in P81_UPGRADE_IDS and row["p81_status"] == "real_pass")
    remaining = len(full_rows) - len(full_real)
    conclusion = "release_ready_scoped_with_p81_dynamic_source_progress"
    if remaining == 0:
        conclusion = "release_ready_full_requirements_traceable"
    P81_ACCEPTANCE.write_text(
        "\n".join(
            [
                "# P81 Dynamic Source Field Coverage Acceptance",
                "",
                "> Change: `p81-dynamic-source-field-coverage`",
                f"> Generated: {datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
                f"> Conclusion: `{conclusion}`",
                "",
                "## Summary",
                "",
                f"- Total rows: {len(rows)}",
                f"- Full-release-required rows: {len(full_rows)}",
                f"- Full-release-required `real_pass` rows after P81: {len(full_real)}",
                f"- Remaining full-release-required non-`real_pass` rows: {remaining}",
                f"- Newly upgraded P81 rows: {len(new_real)}",
                f"- Evidence status: `{evidence['status']}`",
                f"- Evidence symbol: `{evidence['symbol']}`",
                f"- Tracked index symbol: `{evidence['tracked_index_symbol']}`",
                "",
                "## Fresh Evidence",
                "",
                f"- `{P81_GO_TEST_COMMAND}`",
                f"- `{P81_UI_COMMAND}`",
                f"- `{P81_CHECK_COMMAND}`",
                f"- Browser results: `{evidence['browser_results']}`",
                f"- SQLite/readback summary: `{evidence['db_summary']}`",
                f"- Go test log: `{evidence['go_test_log']}`",
                f"- Summary JSON: `{rel(P81_SUMMARY)}`",
                "",
                "## Upgraded Rows",
                "",
                ", ".join(new_real),
                "",
                "## Remaining Boundaries",
                "",
                "- P81 only upgrades dynamic source field coverage rows directly proven by the fresh non-510300 user-symbol evidence.",
                "- P81 does not claim full original-requirement pass while any full-release-required row remains non-`real_pass`.",
                "- P81 does not refresh the P76 package.",
                "- P81 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, return promises, paid/login/authorized sources, Level2, or high-frequency sources.",
                "",
            ]
        ),
        encoding="utf-8",
    )


def write_summary(rows: list[dict[str, str]], evidence: dict[str, Any]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p81_status"] == "real_pass"]
    payload = {
        "generated_at": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "change": "p81-dynamic-source-field-coverage",
        "status": evidence["status"],
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "new_p81_real_pass_count": len([row for row in rows if row["requirement_id"] in P81_UPGRADE_IDS and row["p81_status"] == "real_pass"]),
        "new_p81_real_pass_rows": sorted(row["requirement_id"] for row in rows if row["requirement_id"] in P81_UPGRADE_IDS and row["p81_status"] == "real_pass"),
        "remaining_by_group": dict(sorted(Counter(row.get("remediation_group") for row in full_rows if row["p81_status"] != "real_pass").items())),
        "evidence": evidence,
    }
    P81_SUMMARY.parent.mkdir(parents=True, exist_ok=True)
    P81_SUMMARY.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def validate(rows: list[dict[str, str]], evidence: dict[str, Any]) -> None:
    if evidence["status"] != "passed":
        raise SystemExit(f"P81 evidence failed: {evidence['failures']}")
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    p80_remaining = {row["requirement_id"] for row in full_rows if row.get("p80_status") != "real_pass"}
    missing = sorted(P81_UPGRADE_IDS - p80_remaining)
    if missing:
        raise SystemExit(f"P81 upgrade ids not present in P80 remaining full rows: {missing}")
    for row in rows:
        rid = row["requirement_id"]
        if rid in P81_UPGRADE_IDS and row["p81_status"] != "real_pass":
            raise SystemExit(f"P81 row did not upgrade: {rid}")
        if rid not in P81_UPGRADE_IDS and row.get("p80_status") != "real_pass" and row["p81_status"] == "real_pass":
            raise SystemExit(f"P81 over-upgraded non-P81 row: {rid}")
        if row["p81_status"] == "real_pass" and row["p81_fresh_evidence_artifact"] == "N/A":
            raise SystemExit(f"P81 real_pass row missing evidence artifact: {rid}")
    remaining = [row for row in full_rows if row["p81_status"] != "real_pass"]
    if len(remaining) != 214:
        raise SystemExit(f"Expected 214 remaining full-release-required non-real-pass rows after P81, got {len(remaining)}")


def run(check_only: bool) -> None:
    header, p80_rows = read_p80_rows()
    evidence = verify_evidence(sanitize=not check_only)
    rows: list[dict[str, str]] = []
    for source in p80_rows:
        row = dict(source)
        row.update(classify_row(row, evidence))
        rows.append(row)
    validate(rows, evidence)
    if check_only:
        return
    write_summary(rows, evidence)
    write_matrix(header, rows)
    write_acceptance(rows, evidence)
    print(
        "p81_dynamic_source_field_coverage:"
        f"status={evidence['status']}:"
        f"new_real={len([row for row in rows if row['requirement_id'] in P81_UPGRADE_IDS and row['p81_status'] == 'real_pass'])}:"
        "remaining_full=214"
    )


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="Validate existing evidence without rewriting P81 artifacts.")
    args = parser.parse_args()
    run(check_only=args.check)


if __name__ == "__main__":
    main()
