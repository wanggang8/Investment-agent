#!/usr/bin/env python3
"""Generate P78 real-pass batch closure artifacts from the P77 matrix."""

from __future__ import annotations

import argparse
import json
import sqlite3
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
P77_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p77-requirements-real-pass-upgrade-matrix.md"
P78_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p78-requirements-real-pass-batch-matrix.md"
P78_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p78-requirements-real-pass-batch-closure.md"
P78_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p78"
P78_SUMMARY = P78_ASSET_DIR / "real-pass-batch-summary.json"
P78_EXPECTED_RETURN_LOG = P78_ASSET_DIR / "expected-return-go-tests.log"
P78_EXPECTED_RETURN_META = P78_ASSET_DIR / "expected-return-go-tests.json"
P78_EXPECTED_RETURN_READBACK = P78_ASSET_DIR / "expected-return-ui-readback.json"
P78_NON510300_ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p78-non-510300"
P78_NON510300_SUMMARY = P78_NON510300_ARTIFACT_DIR / "non-510300-db-impact-summary.json"

BASE_COLUMNS = [
    "requirement_id",
    "source_section",
    "source_start_line",
    "source_end_line",
    "requirement_text_hash",
    "requirement_text",
    "status",
    "criticality",
    "full_release_required",
    "release_impact",
    "original_status",
    "p77_status",
]
P78_COLUMNS = [
    "p78_status",
    "remediation_group",
    "batch",
    "closure_basis",
    "fresh_evidence_command",
    "fresh_evidence_artifact",
    "remaining_gap",
    "next_action",
]

EXPECTED_RETURN_UPGRADE_IDS = {
    "REQ-09-002",
    "REQ-09-017",
    "REQ-09-026",
}

EXPECTED_RETURN_OVERBROAD_IDS = {
    "REQ-09-014",
    "REQ-09-027",
    "REQ-09-028",
}

EXPECTED_RETURN_COMMAND = (
    "go test -v ./internal/application/workflow "
    "-run 'TestBuildExpectedReturnIncludesSampleContextForAllPrecisionStates|"
    "TestBuildExpectedReturnProducesAdvisorySellEvaluation|"
    "TestBuildExpectedReturnDoesNotTriggerTargetWithoutConfiguredTarget|"
    "TestBuildExpectedReturnUsesScenarioBoundsForSellTriggers|"
    "TestBuildExpectedReturnCoversAllSellEvaluationTriggers|"
    "TestExpectedReturnNodeUsesWorkflowPricesForSellEvaluation|"
    "TestExpectedReturnNodeUsesMatchingSymbolPosition|"
    "TestExpectedReturnNodeUsesWorkflowDynamicSellInputs|"
    "TestExpectedReturnNodeIncludesP34SupportingDataContext|"
    "TestExpectedReturnSampleCountFromWorkflowDataUsesMarketHistory|"
    "TestExpectedReturnSampleCountFromWorkflowDataDoesNotInventSamples|"
    "TestBuildExpectedReturnExplainsMissingPriceContext' -count=1 && "
    "go test -v ./internal/domain/rule -run TestExpectedReturnDoesNotOverrideVerdict -count=1"
)

NON510300_UI_COMMAND = (
    "P75_ARTIFACT_DIR=docs/release/ui-audit-assets/2026-06-21-p78-non-510300 "
    "bash scripts/p75-non-510300-real-ui-journey.sh"
)

EXPECTED_RETURN_TEST_NAMES = [
    "TestBuildExpectedReturnIncludesSampleContextForAllPrecisionStates",
    "TestBuildExpectedReturnProducesAdvisorySellEvaluation",
    "TestBuildExpectedReturnDoesNotTriggerTargetWithoutConfiguredTarget",
    "TestBuildExpectedReturnUsesScenarioBoundsForSellTriggers",
    "TestBuildExpectedReturnCoversAllSellEvaluationTriggers",
    "TestExpectedReturnNodeUsesWorkflowPricesForSellEvaluation",
    "TestExpectedReturnNodeUsesMatchingSymbolPosition",
    "TestExpectedReturnNodeUsesWorkflowDynamicSellInputs",
    "TestExpectedReturnNodeIncludesP34SupportingDataContext",
    "TestExpectedReturnSampleCountFromWorkflowDataUsesMarketHistory",
    "TestExpectedReturnSampleCountFromWorkflowDataDoesNotInventSamples",
    "TestBuildExpectedReturnExplainsMissingPriceContext",
    "TestExpectedReturnDoesNotOverrideVerdict",
]


def split_markdown_row(line: str) -> list[str]:
    stripped = line.rstrip("\n")
    if not stripped.startswith("|") or not stripped.endswith("|"):
        return []
    body = stripped[1:-1]
    cells: list[str] = []
    current: list[str] = []
    escaped = False
    for char in body:
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


def redact_path(path: object) -> str:
    raw = str(path or "")
    if not raw:
        return raw
    candidate = Path(raw)
    if candidate.is_absolute():
        try:
            return str(candidate.relative_to(ROOT))
        except ValueError:
            return "[redacted-absolute-path]"
    return raw


def resolve_artifact_path(path: object) -> Path:
    raw = str(path or "")
    if not raw:
        return Path("")
    candidate = Path(raw)
    return candidate if candidate.is_absolute() else ROOT / candidate


def read_p77_rows() -> tuple[list[str], list[dict[str, str]]]:
    lines = P77_MATRIX.read_text(encoding="utf-8").splitlines()
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in lines:
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
            raise SystemExit(f"Invalid P77 matrix row column count: expected={len(header)} got={len(cells)}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("P77 matrix header not found")
    missing = [column for column in BASE_COLUMNS if column not in header]
    if missing:
        raise SystemExit(f"P77 matrix missing required columns: {missing}")
    return header, rows


def remediation_group(row: dict[str, str]) -> str:
    section = row["source_section"].split(".")[0]
    text = row["requirement_text"]
    gap = f"{row.get('residual_gap', '')} {row.get('next_remediation', '')}"
    if section == "9" or "预期收益" in text or "概率" in text or "情景" in text:
        return "expected_return"
    if section == "11" or "持仓" in text or "SQLite" in text or "确认" in text:
        return "portfolio_confirmation_data"
    if "SOP" in text or "预警" in text or "action" in gap or "readback" in gap:
        return "sop_action_data_impact"
    if "数据" in text or "source" in gap or "collector" in gap or "源" in text:
        return "data_source_dynamic"
    if "大师" in text or "LLM" in text or "RAG" in text or "知识" in text:
        return "knowledge_llm_rag"
    if section in {"13", "14", "15", "16", "17"}:
        return "governance_traceability"
    if "不" in text and ("交易" in text or "承诺" in text or "预测" in text):
        return "safety_boundary"
    if section in {"1", "2", "3", "4"}:
        return "core_product_goal"
    return "unclassified_followup"


def batch_for(row: dict[str, str], group: str) -> str:
    if row["requirement_id"] in EXPECTED_RETURN_UPGRADE_IDS:
        return "P78-A-expected-return-degradation"
    if group == "expected_return":
        return "P79-expected-return-full-scenario-ui"
    if group == "portfolio_confirmation_data":
        return "P79-portfolio-confirmation-data-impact"
    if group == "sop_action_data_impact":
        return "P80-sop-action-matrix"
    if group == "data_source_dynamic":
        return "P81-dynamic-source-field-coverage"
    if group == "knowledge_llm_rag":
        return "P82-knowledge-llm-rag-coverage"
    if group == "governance_traceability":
        return "P83-governance-traceability-backfill"
    return "P84-core-goal-decomposition"


def read_expected_return_log() -> dict[str, object]:
    if not P78_EXPECTED_RETURN_LOG.exists():
        return {"status": "missing", "missing_artifact": str(P78_EXPECTED_RETURN_LOG)}
    text = P78_EXPECTED_RETURN_LOG.read_text(encoding="utf-8", errors="replace")
    missing = [name for name in EXPECTED_RETURN_TEST_NAMES if f"=== RUN   {name}" not in text or f"--- PASS: {name}" not in text]
    package_pass_lines = [line for line in text.splitlines() if line.strip() == "PASS"]
    status = "passed" if not missing and len(package_pass_lines) >= 2 else "failed"
    meta = {
        "status": status,
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "command": EXPECTED_RETURN_COMMAND,
        "log": str(P78_EXPECTED_RETURN_LOG.relative_to(ROOT)),
        "required_tests": EXPECTED_RETURN_TEST_NAMES,
        "missing_tests": missing,
        "package_pass_line_count": len(package_pass_lines),
    }
    P78_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P78_EXPECTED_RETURN_META.write_text(json.dumps(meta, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return meta


def no_precise_probability(payload: dict[str, object]) -> bool:
    scenarios = payload.get("scenarios") or []
    if not isinstance(scenarios, list):
        return True
    return all(item.get("probability") is None for item in scenarios if isinstance(item, dict))


def read_expected_return_ui() -> dict[str, object]:
    if not P78_NON510300_SUMMARY.exists():
        return {"status": "missing", "missing_artifact": str(P78_NON510300_SUMMARY)}
    summary = json.loads(P78_NON510300_SUMMARY.read_text(encoding="utf-8"))
    if summary.get("status") != "passed":
        return {"status": "failed", "reason": "non-510300 summary did not pass", "summary": summary}
    sqlite_path = resolve_artifact_path(summary.get("sqlite_path", ""))
    if not sqlite_path.exists():
        return {"status": "failed", "reason": "sqlite path missing", "sqlite_path": str(sqlite_path)}
    summary["sqlite_path"] = redact_path(sqlite_path)
    if "request_log_path" in summary:
        summary["request_log_path"] = redact_path(summary.get("request_log_path"))
    P78_NON510300_SUMMARY.write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")

    db = sqlite3.connect(f"file:{sqlite_path}?mode=ro", uri=True)
    db.row_factory = sqlite3.Row
    decision_id = ((summary.get("decision") or {}).get("decision_id") or "")
    row = db.execute(
        "SELECT decision_id,symbol,workflow_status,expected_return_scenarios_json "
        "FROM decision_records WHERE decision_id=?",
        (decision_id,),
    ).fetchone()
    if row is None:
        return {"status": "failed", "reason": "decision not found", "decision_id": decision_id}
    payload = json.loads(row["expected_return_scenarios_json"] or "{}")
    sell_evaluation = payload.get("sell_evaluation") or {}
    source_health = payload.get("source_health") or []
    checks = {
        "decision_completed": row["workflow_status"] == "completed",
        "symbol_bound": row["symbol"] == "159915",
        "precision_unavailable_under_low_samples": payload.get("precision_status") == "unavailable" and int(payload.get("sample_count") or 0) < 5,
        "no_precise_probability": no_precise_probability(payload),
        "sample_window_present": bool(payload.get("sample_window")),
        "screening_condition_present": bool(payload.get("screening_condition")),
        "source_health_present": isinstance(source_health, list) and len(source_health) >= 3,
        "non_trading_disclaimer_present": "不会自动交易" in str(sell_evaluation.get("non_trading_disclaimer", "")),
        "no_return_promise_disclaimer_present": "不构成收益承诺" in str(sell_evaluation.get("non_trading_disclaimer", "")),
    }
    status = "passed" if all(checks.values()) else "failed"
    readback = {
        "status": status,
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "source_summary": str(P78_NON510300_SUMMARY.relative_to(ROOT)),
        "sqlite_path": redact_path(sqlite_path),
        "decision_id": decision_id,
        "symbol": row["symbol"],
        "checks": checks,
        "expected_return": {
            "precision_status": payload.get("precision_status"),
            "reason": payload.get("reason"),
            "sample_count": payload.get("sample_count"),
            "sample_window": payload.get("sample_window"),
            "screening_condition": payload.get("screening_condition"),
            "scenario_count": len(payload.get("scenarios") or []),
            "sell_evaluation_status": sell_evaluation.get("status"),
            "sell_evaluation_prompts": sell_evaluation.get("prompts"),
            "non_trading_disclaimer": sell_evaluation.get("non_trading_disclaimer"),
            "source_health_count": len(source_health) if isinstance(source_health, list) else 0,
        },
    }
    P78_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P78_EXPECTED_RETURN_READBACK.write_text(json.dumps(readback, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    return readback


def expected_return_evidence_ready() -> tuple[bool, dict[str, object], dict[str, object]]:
    log_meta = read_expected_return_log()
    ui_readback = read_expected_return_ui()
    sanitize_p78_text_artifacts()
    ready = log_meta.get("status") == "passed" and ui_readback.get("status") == "passed"
    return ready, log_meta, ui_readback


def sanitize_p78_text_artifacts() -> None:
    prefixes = [str(ROOT) + "/", str(ROOT.resolve()) + "/"]
    for base in [P78_ASSET_DIR, P78_NON510300_ARTIFACT_DIR]:
        if not base.exists():
            continue
        for path in base.rglob("*"):
            if path.suffix.lower() not in {".json", ".log", ".txt", ".md"}:
                continue
            try:
                text = path.read_text(encoding="utf-8")
            except UnicodeDecodeError:
                continue
            updated = text
            for prefix in prefixes:
                updated = updated.replace(prefix, "")
            if updated != text:
                path.write_text(updated, encoding="utf-8")


def enrich_rows(rows: list[dict[str, str]], evidence_ready: bool) -> list[dict[str, str]]:
    out: list[dict[str, str]] = []
    for row in rows:
        group = remediation_group(row)
        batch = batch_for(row, group)
        p77_status = row["p77_status"]
        rid = row["requirement_id"]

        if p77_status == "real_pass":
            p78_status = "real_pass"
            basis = "Carried forward from P77 real_pass; P78 does not rewrite P77 evidence."
            command = "N/A"
            artifact = "docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md"
            gap = "None for the row already accepted by P77."
            next_action = "Keep covered by future regression evidence."
        elif rid in EXPECTED_RETURN_UPGRADE_IDS and evidence_ready:
            p78_status = "real_pass"
            basis = "P78 batch A fresh expected-return Go tests plus accepted-local non-510300 real UI SQLite readback prove this expected-return degradation/disclaimer row."
            command = f"{EXPECTED_RETURN_COMMAND} && {NON510300_UI_COMMAND} && python3 scripts/p78_requirements_real_pass_batch_closure.py --check"
            artifact = "; ".join(
                [
                    "docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.log",
                    "docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.json",
                    "docs/release/ui-audit-assets/2026-06-21-p78/expected-return-ui-readback.json",
                    "docs/release/ui-audit-assets/2026-06-21-p78-non-510300",
                ]
            )
            gap = "None for this narrow expected-return degradation/disclaimer row; broader expected-return scenario, provenance, and sell-trigger rows remain separate."
            next_action = "Keep in P78/P79 regression set and extend real UI coverage for the remaining REQ-09 rows."
        elif rid in EXPECTED_RETURN_UPGRADE_IDS:
            p78_status = "p78_pending_evidence"
            basis = "P78 batch A candidate, but fresh evidence artifacts are not complete yet."
            command = f"{EXPECTED_RETURN_COMMAND} && {NON510300_UI_COMMAND}"
            artifact = "docs/release/ui-audit-assets/2026-06-21-p78; docs/release/ui-audit-assets/2026-06-21-p78-non-510300"
            gap = "Missing fresh expected-return Go log or real UI SQLite readback."
            next_action = "Run P78 batch A evidence commands, then rerun the P78 checker."
        else:
            p78_status = p77_status
            basis = "No P78 upgrade; classified for a later closure batch."
            command = "N/A"
            artifact = "N/A"
            gap = row.get("residual_gap") or "Still non-real-pass after P77."
            next_action = next_action_for(group, row)

        updated = dict(row)
        updated.update(
            {
                "p78_status": p78_status,
                "remediation_group": group,
                "batch": batch,
                "closure_basis": basis,
                "fresh_evidence_command": command,
                "fresh_evidence_artifact": artifact,
                "remaining_gap": gap,
                "next_action": next_action,
            }
        )
        out.append(updated)
    return out


def next_action_for(group: str, row: dict[str, str]) -> str:
    if group == "expected_return":
        return "Add real UI scenario coverage for available precision, sell-evaluation triggers, scenario ranges, valuation fields, and readback."
    if group == "portfolio_confirmation_data":
        return "Run a multi-action real UI action-to-SQLite-to-readback matrix covering add/edit/import/confirm/offline transactions."
    if group == "sop_action_data_impact":
        return "Split each SOP action into exact UI action, changed tables, audit event, prohibited tables, and readback assertions."
    if group == "data_source_dynamic":
        return "Add dynamic source/readiness evidence and field-level propagation tests for each required data category."
    if group == "knowledge_llm_rag":
        return "Add row-level LLM context, built-in knowledge, local knowledge, and RAG retrieval evidence with source references."
    if group == "governance_traceability":
        return "Backfill delivered_by_change, command, artifact, and release-boundary evidence for the row."
    if group == "core_product_goal":
        return "Decompose broad goal into scenario rows and link each scenario to UI, data, workflow, and audit evidence."
    return row.get("next_remediation") or "Assign to a later closure batch."


def write_matrix(rows: list[dict[str, str]]) -> None:
    header = BASE_COLUMNS + P78_COLUMNS
    lines = [
        "# P78 Requirements Real-Pass Batch Matrix",
        "",
        "> Generated: 2026-06-21",
        "> Source: `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`",
        "> Policy: P78 is a new batch closure layer; it does not rewrite P75 or P77 history.",
        "",
        "## Status Summary",
        "",
    ]
    counts = Counter(row["p78_status"] for row in rows)
    for status, count in sorted(counts.items()):
        lines.append(f"- `{status}`: {count}")
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p78_status"] == "real_pass"]
    lines.extend(
        [
            "",
            f"- full_release_required rows: {len(full_rows)}",
            f"- full_release_required real_pass rows: {len(full_real)}",
            f"- remaining full_release_required non-real-pass rows: {len(full_rows) - len(full_real)}",
            f"- conclusion: `{conclusion(rows)}`",
            "",
            "## Remediation Groups",
            "",
        ]
    )
    for group, count in sorted(Counter(row["remediation_group"] for row in full_rows if row["p78_status"] != "real_pass").items()):
        lines.append(f"- `{group}`: {count}")
    lines.extend(["", "## Atomic Requirement Batch Rows", "", "|" + "|".join(header) + "|", "|" + "|".join(["---"] * len(header)) + "|"])
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in header) + "|")
    P78_MATRIX.parent.mkdir(parents=True, exist_ok=True)
    P78_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def conclusion(rows: list[dict[str, str]]) -> str:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    if all(row["p78_status"] == "real_pass" for row in full_rows):
        return "release_ready_full_requirements_traceable"
    if any(row["p78_status"] == "p78_pending_evidence" for row in full_rows):
        return "release_pending_p78_batch_evidence"
    return "release_ready_scoped_with_p78_real_pass_batch_progress"


def write_summary(rows: list[dict[str, str]], log_meta: dict[str, object], ui_readback: dict[str, object]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p78_status"] == "real_pass"]
    newly_upgraded = [row["requirement_id"] for row in rows if row["requirement_id"] in EXPECTED_RETURN_UPGRADE_IDS and row["p78_status"] == "real_pass"]
    summary = {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "source_matrix": str(P77_MATRIX.relative_to(ROOT)),
        "matrix": str(P78_MATRIX.relative_to(ROOT)),
        "acceptance": str(P78_ACCEPTANCE.relative_to(ROOT)),
        "status_counts": dict(sorted(Counter(row["p78_status"] for row in rows).items())),
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "new_p78_real_pass_rows": newly_upgraded,
        "new_p78_real_pass_count": len(newly_upgraded),
        "remediation_group_counts": dict(sorted(Counter(row["remediation_group"] for row in full_rows if row["p78_status"] != "real_pass").items())),
        "conclusion": conclusion(rows),
        "expected_return_go_tests": log_meta,
        "expected_return_ui_readback": ui_readback,
        "package_boundary": "P78 does not refresh the P76 package.",
        "not_claimed": [
            "full original-requirement pass",
            "P78 evidence inside any existing P76 distribution archive",
            "physical second-machine repeat",
            "broker connectivity",
            "automatic trading",
            "one-click trading",
            "order delegation",
            "external push",
            "automatic confirmation",
            "automatic rule application",
            "automatic repair/migration/restore",
            "future provider availability",
            "investment returns",
        ],
    }
    P78_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P78_SUMMARY.write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


def write_acceptance(rows: list[dict[str, str]]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p78_status"] == "real_pass"]
    newly_upgraded = [row for row in rows if row["requirement_id"] in EXPECTED_RETURN_UPGRADE_IDS and row["p78_status"] == "real_pass"]
    lines = [
        "# P78 Requirements Real-Pass Batch Closure Acceptance",
        "",
        "> Date: 2026-06-21",
        "> Change: `p78-requirements-real-pass-batch-closure`",
        f"> Conclusion: `{conclusion(rows)}`",
        "",
        "## Summary",
        "",
        f"- Source matrix: `{P77_MATRIX.relative_to(ROOT)}`",
        f"- P78 matrix: `{P78_MATRIX.relative_to(ROOT)}`",
        f"- Summary JSON: `{P78_SUMMARY.relative_to(ROOT)}`",
        f"- Full-release-required rows: {len(full_rows)}",
        f"- Full-release-required `real_pass` rows after P78: {len(full_real)}",
        f"- Remaining full-release-required non-`real_pass` rows: {len(full_rows) - len(full_real)}",
        f"- Newly upgraded by P78: {len(newly_upgraded)}",
        "",
        "## P78 Batch A Upgrades",
        "",
    ]
    for row in newly_upgraded:
        lines.append(f"- `{row['requirement_id']}` `{row['source_section']}`: {row['requirement_text']}")
    if not newly_upgraded:
        lines.append("- None yet; P78 evidence is pending.")
    lines.extend(
        [
            "",
            "## Fresh Evidence",
            "",
            f"- Expected-return Go tests: `{P78_EXPECTED_RETURN_LOG.relative_to(ROOT)}` and `{P78_EXPECTED_RETURN_META.relative_to(ROOT)}`",
            f"- Accepted-local non-`510300` real UI journey: `{P78_NON510300_ARTIFACT_DIR.relative_to(ROOT)}`",
            f"- Expected-return SQLite readback: `{P78_EXPECTED_RETURN_READBACK.relative_to(ROOT)}`",
            "",
            "Commands:",
            "",
            "```bash",
            EXPECTED_RETURN_COMMAND,
            NON510300_UI_COMMAND,
            "python3 scripts/p78_requirements_real_pass_batch_closure.py --check",
            "```",
            "",
            "## Boundaries",
            "",
            "- P78 does not rewrite P75 or P77 historical matrices.",
            "- P78 does not refresh the P76 package; a separate package refresh is required before claiming distribution archives include P78 materials.",
            "- P78 does not claim full original-requirement pass while any full-release-required row remains non-`real_pass`.",
            "- P78 does not add broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.",
        ]
    )
    P78_ACCEPTANCE.parent.mkdir(parents=True, exist_ok=True)
    P78_ACCEPTANCE.write_text("\n".join(lines) + "\n", encoding="utf-8")


def validate_outputs(rows: list[dict[str, str]], log_meta: dict[str, object], ui_readback: dict[str, object]) -> None:
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    if conclusion(rows) == "release_ready_full_requirements_traceable" and not all(row["p78_status"] == "real_pass" for row in full_rows):
        raise SystemExit("P78 attempted full traceability conclusion while gaps remain")
    if log_meta.get("status") != "passed":
        raise SystemExit(f"Expected-return Go evidence not passed: {log_meta}")
    if ui_readback.get("status") != "passed":
        raise SystemExit(f"Expected-return UI readback not passed: {ui_readback}")
    missing = [path for path in [P78_MATRIX, P78_ACCEPTANCE, P78_SUMMARY, P78_EXPECTED_RETURN_META, P78_EXPECTED_RETURN_READBACK, P78_NON510300_SUMMARY] if not path.exists()]
    if missing:
        raise SystemExit(f"Missing P78 artifacts: {missing}")
    for row in rows:
        if row["requirement_id"] in EXPECTED_RETURN_UPGRADE_IDS and row["p78_status"] != "real_pass":
            raise SystemExit(f"P78 expected-return candidate did not upgrade: {row['requirement_id']}")
        if row["requirement_id"] in EXPECTED_RETURN_OVERBROAD_IDS and row["p78_status"] == "real_pass":
            raise SystemExit(f"P78 overbroad expected-return row must not be upgraded by low-sample batch A evidence: {row['requirement_id']}")
        if row["p78_status"] == "real_pass" and row["fresh_evidence_artifact"] == "N/A":
            raise SystemExit(f"real_pass row missing evidence artifact: {row['requirement_id']}")
    full_non_real = [row for row in full_rows if row["p78_status"] != "real_pass"]
    if len(full_non_real) == 0:
        return
    if conclusion(rows) != "release_ready_scoped_with_p78_real_pass_batch_progress":
        raise SystemExit(f"Unexpected scoped conclusion for remaining gaps: {conclusion(rows)}")
    validate_release_claims()


def validate_release_claims() -> None:
    scan_files = [
        P78_MATRIX,
        P78_ACCEPTANCE,
        P78_SUMMARY,
        P78_EXPECTED_RETURN_META,
        P78_EXPECTED_RETURN_READBACK,
        P78_NON510300_SUMMARY,
        ROOT / "docs" / "release" / "README.md",
        ROOT / "docs" / "release" / "acceptance-repeatability.md",
        ROOT / "docs" / "release" / "release-candidate-2026-06-18.md",
        ROOT / "docs" / "release" / "release-handoff-2026-06-18.md",
        ROOT / "docs" / "README.md",
        ROOT / "docs" / "GOVERNANCE.md",
        ROOT / "AGENTS.md",
        ROOT / "openspec" / "PROGRESS.md",
        ROOT / "openspec" / "project.md",
        ROOT / "docs" / "development-plan.md",
    ]
    for base in [P78_ASSET_DIR, P78_NON510300_ARTIFACT_DIR]:
        if base.exists():
            scan_files.extend(path for path in base.rglob("*") if path.suffix.lower() in {".json", ".log", ".txt", ".md"})
    private_path_hits: list[str] = []
    p76_overclaims: list[str] = []
    full_pass_overclaims: list[str] = []
    for path in scan_files:
        if not path.exists():
            continue
        text = path.read_text(encoding="utf-8", errors="replace")
        if "/Users/" in text:
            private_path_hits.append(str(path.relative_to(ROOT)))
        for line in text.splitlines():
            if not any(pattern in line for pattern in ["P76 package includes P78", "archive includes P78", "包内包含 P78"]):
                continue
            lowered = line.lower()
            is_boundary = (
                "not_claimed" in lowered
                or "does not" in lowered
                or "not " in lowered
                or "required before claiming" in lowered
                or "separate package" in lowered
                or "不" in line
                or "不得" in line
            )
            if not is_boundary:
                p76_overclaims.append(str(path.relative_to(ROOT)))
        if "release_ready_full_requirements_traceable" in text and path not in {P78_MATRIX}:
            allowed = "not `release_ready_full_requirements_traceable`" in text or "only if every" in text or "不得" in text
            if not allowed:
                full_pass_overclaims.append(str(path.relative_to(ROOT)))
    if private_path_hits:
        raise SystemExit(f"P78 release artifacts contain private absolute paths: {private_path_hits}")
    if p76_overclaims:
        raise SystemExit(f"P78 materials overclaim P76 package freshness: {p76_overclaims}")
    if full_pass_overclaims:
        raise SystemExit(f"P78 materials overclaim full requirements pass: {full_pass_overclaims}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="validate generated P78 artifacts")
    args = parser.parse_args()

    _, p77_rows = read_p77_rows()
    evidence_ready, log_meta, ui_readback = expected_return_evidence_ready()
    rows = enrich_rows(p77_rows, evidence_ready)
    write_matrix(rows)
    write_summary(rows, log_meta, ui_readback)
    write_acceptance(rows)
    if args.check:
        validate_outputs(rows, log_meta, ui_readback)
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p78_status"] == "real_pass"]
    print(
        "p78_real_pass_batch:"
        f"rows={len(rows)}:"
        f"real_pass={len([row for row in rows if row['p78_status'] == 'real_pass'])}:"
        f"new={len([row for row in rows if row['requirement_id'] in EXPECTED_RETURN_UPGRADE_IDS and row['p78_status'] == 'real_pass'])}:"
        f"remaining_full={len(full_rows) - len(full_real)}:"
        f"conclusion={conclusion(rows)}"
    )


if __name__ == "__main__":
    main()
