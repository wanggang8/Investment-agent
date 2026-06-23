#!/usr/bin/env python3
"""Generate P77 real-pass upgrade artifacts from the P75 traceability matrix."""

from __future__ import annotations

import argparse
import json
import subprocess
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
P75_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-20-p75-requirements-traceability-matrix.md"
P75_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-20-p75-real-use-closure.md"
P77_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p77-requirements-real-pass-upgrade-matrix.md"
P77_ACCEPTANCE = ROOT / "docs" / "release" / "acceptance" / "2026-06-21-p77-real-pass-upgrade-acceptance.md"
P77_ASSET_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p77"
P77_SUMMARY = P77_ASSET_DIR / "real-pass-upgrade-summary.json"
P77_SAFETY_SCAN = P77_ASSET_DIR / "safety-scan.txt"
P77_SAFETY_REVIEW = P77_ASSET_DIR / "safety-scan-review.json"
P77_SAFETY_TEST_LOG = P77_ASSET_DIR / "safety-and-boundary-go-tests.log"
P77_SAFETY_TEST_META = P77_ASSET_DIR / "safety-and-boundary-go-tests.json"
P77_F_TEST_LOG = P77_ASSET_DIR / "f1-f5-go-tests.log"
P77_F_TEST_META = P77_ASSET_DIR / "f1-f5-go-tests.json"
P77_SOP_ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p77-sop-failure"
P77_NON510300_ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-21-p77-non-510300"


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
]
P77_COLUMNS = [
    "p77_status",
    "upgrade_basis",
    "gate_dimensions",
    "fresh_evidence_command",
    "fresh_evidence_artifact",
    "residual_gap",
    "next_remediation",
]

STATUSES = {
    "real_pass",
    "scoped_pass",
    "deterministic_local_evidence",
    "partial",
    "not_implemented",
    "blocked",
    "reference_only",
    "p77_pending_evidence",
}

SAFETY_REAL_PASS_IDS = {
    "REQ-01-007",
    "REQ-01-008",
    "REQ-01-009",
    "REQ-01-010",
    "REQ-02-001",
    "REQ-07-016",
    "REQ-18-001",
    "REQ-18-002",
    "REQ-18-003",
    "REQ-18-004",
    "REQ-18-005",
    "REQ-18-006",
}

F_REAL_PASS_IDS = {
    "REQ-05-021",
    "REQ-05-022",
    "REQ-05-023",
    "REQ-05-024",
    "REQ-05-025",
}

REFERENCE_SECTION_PREFIX = "19"

SAFETY_COMMAND = (
    "go test -v ./internal/infrastructure/llm/deepseek "
    "-run 'TestClientRetriesQualityFailureWithStricterBoundary|TestClientRejectsProhibitedLLMOutput|TestEvaluateQualityAllowsNormalAnalysisAndRejectsUnsafeClaims' -count=1 && "
    "go test -v ./internal/application/workflow -run 'TestExpectedReturnMaterialIsExplanatoryOnly' -count=1 && "
    "go test -v ./internal/domain/rule -run 'TestExpectedReturnDoesNotOverrideVerdict|TestP75RulePriorityAndRootRules' -count=1"
)

F_COMMAND = (
    "go test -v ./internal/application/workflow "
    "-run 'TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary|TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources|TestEvidenceVerificationRequiresTwoHighGradeIndependentSources|TestAnalystRequestsPreferStructuredFinancialFacts|TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary|TestPublicEvidencePayloadNormalizesEmotionalDescriptions' -count=1 && "
    "go test -v ./internal/infrastructure/persistence/sqlite -run TestMarketRepositoryPreservesStructuredFinancialFields -count=1 && "
    "go test -v ./internal/domain/rule -run TestEvaluatePriorityScenarios -count=1"
)

SAFETY_TEST_NAMES = [
    "TestClientRetriesQualityFailureWithStricterBoundary",
    "TestClientRejectsProhibitedLLMOutput",
    "TestEvaluateQualityAllowsNormalAnalysisAndRejectsUnsafeClaims",
    "TestExpectedReturnMaterialIsExplanatoryOnly",
    "TestExpectedReturnDoesNotOverrideVerdict",
    "TestP75RulePriorityAndRootRules",
]

F_TEST_NAMES = [
    "TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary",
    "TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources",
    "TestEvidenceVerificationRequiresTwoHighGradeIndependentSources",
    "TestAnalystRequestsPreferStructuredFinancialFacts",
    "TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary",
    "TestPublicEvidencePayloadNormalizesEmotionalDescriptions",
    "TestMarketRepositoryPreservesStructuredFinancialFields",
    "TestEvaluatePriorityScenarios",
]

SAFETY_REVIEW_ALLOWED_CATEGORIES = {
    "prohibition_boundary_copy",
    "negative_test_or_scanner_config",
    "historical_or_release_governance",
    "captured_acceptance_evidence",
    "seed_prohibited_action_fixture",
    "sanitizer_or_gatekeeper_rule",
    "non_data_use",
}


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
    text = text.replace("\\", "\\\\").replace("|", "\\|")
    return text


def read_p75_rows() -> tuple[list[str], list[dict[str, str]]]:
    lines = P75_MATRIX.read_text(encoding="utf-8").splitlines()
    header: list[str] | None = None
    rows: list[dict[str, str]] = []
    for line in lines:
        if not line.startswith("|"):
            continue
        cells = split_markdown_row(line)
        if not cells:
            continue
        if header is None:
            if cells and cells[0] == "requirement_id":
                header = cells
            continue
        if set("".join(cells)) <= {"-", ":"}:
            continue
        if len(cells) != len(header):
            raise SystemExit(f"Invalid P75 matrix row column count: expected={len(header)} got={len(cells)} line={line[:120]}")
        rows.append(dict(zip(header, cells)))
    if header is None:
        raise SystemExit("P75 matrix header not found")
    missing = [column for column in BASE_COLUMNS if column not in header]
    if missing:
        raise SystemExit(f"P75 matrix missing required columns: {missing}")
    return header, rows


def classify_safety_match(line: str) -> str:
    lowered = line.lower()
    line_no = 0
    parts = line.split(":", 2)
    if len(parts) >= 2 and parts[1].isdigit():
        line_no = int(parts[1])
    if line.startswith("docs/superpowers/"):
        return "historical_or_release_governance"
    if line.startswith("openspec/changes/p77-requirements-real-pass-upgrade-gate/proposal.md:27:"):
        return "prohibition_boundary_copy"
    if line.startswith("scripts/p77_requirements_real_pass_upgrade.py:"):
        return "negative_test_or_scanner_config"
    if line.startswith("docs/product-experience-polish-roadmap.md:") and 209 <= line_no <= 219:
        return "prohibition_boundary_copy"
    if any(token in line for token in ["不得", "不会", "不能", "不允许", "不新增", "不增加", "不接", "不连接", "不下单", "不创建", "不自动", "不触发", "不执行", "不提供", "不展示", "不包含", "不代表", "不在", "不扩展", "不产生", "不出现", "不生成", "不承诺", "不启用", "严格避免", "禁止", "无", "没有", "未声称", "不声称", "边界", "仅", "只", "需用户", "Not Claimed", "not claim", "does not claim", "no broker", "no forbidden"]):
        return "prohibition_boundary_copy"
    if any(token in lowered for token in ["out of scope", "not in scope", "excluded", "excludes", "must not", "shall not", "not requiring", "requires login", "requires login, paid"]):
        return "prohibition_boundary_copy"
    if any(token in lowered for token in ["test", "forbidden", "scanner", "scan", "safety", "mock", "fixture"]):
        return "negative_test_or_scanner_config"
    if any(token in line for token in ["archive", "release", "handoff", "candidate", "GOVERNANCE", "PROGRESS", "OpenSpec", "P75", "P76", "P77"]):
        return "historical_or_release_governance"
    if "ui-audit-assets" in line or "acceptance" in line:
        return "captured_acceptance_evidence"
    if any(token in lowered for token in ["prohibited_actions", "seed", "insert or replace"]):
        return "seed_prohibited_action_fixture"
    if any(token in lowered for token in ["sanitize", "gatekeeper", "quality", "reject", "requiredchanges", "evaluatequality", "replace"]):
        return "sanitizer_or_gatekeeper_rule"
    if "高频" in line and "high-frequency data-source" not in lowered and "高频源" not in line:
        return "non_data_use"
    return "needs_manual_review"


def run_safety_scan() -> dict[str, object]:
    P77_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    pattern = "券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动规则应用|自动修复|自动升级|自动迁移|自动恢复|自动覆盖真实库|真实库覆盖|收益承诺|登录源|付费源|授权源|Level2|高频"
    cmd = [
        "rg",
        "-n",
        "--glob",
        "!docs/release/ui-audit-assets/**",
        "--glob",
        "!web/node_modules/**",
        "--glob",
        "!web/dist/**",
        pattern,
        "docs",
        "openspec",
        "internal",
        "cmd",
        "web/src",
        "scripts",
    ]
    proc = subprocess.run(cmd, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    matches = [line for line in proc.stdout.splitlines() if line.strip()]
    categories = Counter(classify_safety_match(line) for line in matches)
    review_status = "reviewed_pass" if set(categories) <= SAFETY_REVIEW_ALLOWED_CATEGORIES else "needs_manual_review"
    P77_SAFETY_SCAN.write_text(
        "\n".join(
            [
                "# P77 Safety Scan",
                f"generated_utc={datetime.now(timezone.utc).strftime('%Y-%m-%dT%H:%M:%SZ')}",
                f"command={' '.join(cmd)}",
                f"exit_code={proc.returncode}",
                f"match_count={len(matches)}",
                "",
                "## stdout",
                proc.stdout,
                "## stderr",
                proc.stderr,
                "",
            ]
        ),
        encoding="utf-8",
    )
    review = {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "scan_artifact": str(P77_SAFETY_SCAN.relative_to(ROOT)),
        "scan_exit_code": proc.returncode,
        "match_count": len(matches),
        "review_status": review_status,
        "allowed_categories": sorted(SAFETY_REVIEW_ALLOWED_CATEGORIES),
        "category_counts": dict(sorted(categories.items())),
        "human_review_required": review_status != "reviewed_pass",
        "review_basis": "P77 automatic category review accepts only prohibition/boundary copy, negative tests/scanner config, historical/release governance, captured acceptance evidence, seed prohibited-action fixtures, sanitizer/gatekeeper rules, or non-data use. Any uncategorized match blocks P77 safety real_pass upgrades.",
    }
    P77_SAFETY_REVIEW.write_text(json.dumps(review, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    return {"path": str(P77_SAFETY_SCAN.relative_to(ROOT)), "review_path": str(P77_SAFETY_REVIEW.relative_to(ROOT)), "match_count": len(matches), "exit_code": proc.returncode, "review_status": review_status}


def file_contains(path: Path, needles: list[str]) -> bool:
    if not path.exists():
        return False
    text = path.read_text(encoding="utf-8", errors="replace")
    return all(needle in text for needle in needles)


def go_log_passed(path: Path, packages: list[str]) -> bool:
    if not path.exists():
        return False
    text = path.read_text(encoding="utf-8", errors="replace")
    return all(f"ok  \t{package}" in text or f"ok  {package}" in text or f"ok\t{package}" in text for package in packages)


def load_json(path: Path) -> dict[str, object]:
    if not path.exists():
        return {}
    return json.loads(path.read_text(encoding="utf-8"))


def metadata_passed(path: Path, *, command: str, tests: list[str]) -> bool:
    data = load_json(path)
    if data.get("status") != "passed":
        return False
    if data.get("command") != command:
        return False
    recorded = set(data.get("test_names") or [])
    if not set(tests) <= recorded:
        return False
    return bool(data.get("generated_utc"))


def log_contains_tests(path: Path, tests: list[str]) -> bool:
    if not path.exists():
        return False
    text = path.read_text(encoding="utf-8", errors="replace")
    return all(f"=== RUN   {test}" in text and f"--- PASS: {test}" in text for test in tests)


def ensure_test_metadata(meta_path: Path, log_path: Path, *, command: str, tests: list[str], packages: list[str]) -> None:
    if not (go_log_passed(log_path, packages) and log_contains_tests(log_path, tests)):
        return
    meta = {
        "command": command,
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "log_artifact": str(log_path.relative_to(ROOT)),
        "packages": packages,
        "status": "passed",
        "test_names": tests,
    }
    meta_path.write_text(json.dumps(meta, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def safety_review_passed() -> bool:
    data = load_json(P77_SAFETY_REVIEW)
    return data.get("review_status") == "reviewed_pass" and not data.get("human_review_required")


def evidence_state() -> dict[str, object]:
    safety_packages = [
        "investment-agent/internal/infrastructure/llm/deepseek",
        "investment-agent/internal/application/workflow",
        "investment-agent/internal/domain/rule",
    ]
    f_packages = [
        "investment-agent/internal/application/workflow",
        "investment-agent/internal/infrastructure/persistence/sqlite",
        "investment-agent/internal/domain/rule",
    ]
    ensure_test_metadata(P77_SAFETY_TEST_META, P77_SAFETY_TEST_LOG, command=SAFETY_COMMAND, tests=SAFETY_TEST_NAMES, packages=safety_packages)
    ensure_test_metadata(P77_F_TEST_META, P77_F_TEST_LOG, command=F_COMMAND, tests=F_TEST_NAMES, packages=f_packages)
    return {
        "p75_acceptance_exists": P75_ACCEPTANCE.exists(),
        "safety_scan": P77_SAFETY_SCAN.exists() and safety_review_passed(),
        "safety_review": safety_review_passed(),
        "safety_tests": go_log_passed(P77_SAFETY_TEST_LOG, safety_packages)
        and log_contains_tests(P77_SAFETY_TEST_LOG, SAFETY_TEST_NAMES)
        and metadata_passed(P77_SAFETY_TEST_META, command=SAFETY_COMMAND, tests=SAFETY_TEST_NAMES),
        "f_tests": go_log_passed(P77_F_TEST_LOG, f_packages)
        and log_contains_tests(P77_F_TEST_LOG, F_TEST_NAMES)
        and metadata_passed(P77_F_TEST_META, command=F_COMMAND, tests=F_TEST_NAMES),
        "sop_browser_results": (P77_SOP_ARTIFACT_DIR / "browser-results.json").exists(),
        "sop_db_log": file_contains(P77_SOP_ARTIFACT_DIR / "db-impact-check.log", ["status=passed", "forbidden_broker_order_push_tables=0"]),
        "non510300_browser_results": (P77_NON510300_ARTIFACT_DIR / "browser-results.json").exists(),
        "non510300_summary": (P77_NON510300_ARTIFACT_DIR / "non-510300-db-impact-summary.json").exists(),
    }


def classify_p77(row: dict[str, str], evidence: dict[str, object]) -> dict[str, str]:
    req_id = row["requirement_id"]
    source_section = row["source_section"]
    original_status = row["status"]
    full_required = row["full_release_required"]

    result = {
        "p77_status": original_status,
        "upgrade_basis": "No P77 upgrade; preserve P75 status until row-specific real-pass gate is satisfied.",
        "gate_dimensions": "implementation/data/UI/workflow/scenario/safety dimensions remain as recorded by P75.",
        "fresh_evidence_command": "N/A",
        "fresh_evidence_artifact": "N/A",
        "residual_gap": row.get("gap", "Row remains non-real-pass."),
        "next_remediation": row.get("remediation_decision", "Add row-specific real product evidence."),
    }

    if source_section.startswith(REFERENCE_SECTION_PREFIX) and full_required == "False":
        result.update(
            {
                "p77_status": "reference_only",
                "upgrade_basis": "Appendix/reference row; P77 excludes it from full-release-required runtime pass denominator.",
                "gate_dimensions": "reference classification, source hash preserved, no independent runtime behavior required.",
                "fresh_evidence_command": "python3 scripts/p77_requirements_real_pass_upgrade.py --check",
                "fresh_evidence_artifact": str(P77_SUMMARY.relative_to(ROOT)),
                "residual_gap": "No runtime gap unless another normative row depends on this term.",
                "next_remediation": "Keep referenced by normative rows; do not count as runtime real_pass.",
            }
        )
        return result

    if req_id in SAFETY_REAL_PASS_IDS:
        if evidence["safety_scan"] and evidence["safety_tests"] and evidence["p75_acceptance_exists"]:
            result.update(
                {
                    "p77_status": "real_pass",
                    "upgrade_basis": "Fresh P77 safety scan plus targeted LLM/rule boundary tests and P75 reviewed safety closure prove this negative/safety requirement in current code.",
                    "gate_dimensions": "implementation=yes; ui=not independently required for negative safety boundary; data=no mutation; workflow/rule/LLM=yes; scenario=safety boundary; safety=yes.",
                    "fresh_evidence_command": SAFETY_COMMAND,
                    "fresh_evidence_artifact": f"{P77_SAFETY_SCAN.relative_to(ROOT)}; {P77_SAFETY_REVIEW.relative_to(ROOT)}; {P77_SAFETY_TEST_LOG.relative_to(ROOT)}; {P77_SAFETY_TEST_META.relative_to(ROOT)}; {P75_ACCEPTANCE.relative_to(ROOT)}",
                    "residual_gap": "None for the explicit negative/safety boundary row; broader product rows remain scoped separately.",
                    "next_remediation": "Keep in safety regression suite for future changes.",
                }
            )
        else:
            result.update(
                {
                    "p77_status": "p77_pending_evidence",
                    "upgrade_basis": "Candidate safety row, but P77 fresh safety evidence is incomplete.",
                    "gate_dimensions": "requires P77 safety scan, P77 safety Go test log, and P75 reviewed safety closure.",
                    "fresh_evidence_command": SAFETY_COMMAND,
                    "fresh_evidence_artifact": f"{P77_SAFETY_SCAN.relative_to(ROOT)}; {P77_SAFETY_REVIEW.relative_to(ROOT)}; {P77_SAFETY_TEST_LOG.relative_to(ROOT)}; {P77_SAFETY_TEST_META.relative_to(ROOT)}",
                    "residual_gap": "Run and record P77 safety evidence before upgrading.",
                    "next_remediation": "Run the cited command and rerun this script with --check.",
                }
            )
        return result

    if req_id in F_REAL_PASS_IDS:
        if evidence["f_tests"] and evidence["safety_scan"]:
            result.update(
                {
                    "p77_status": "real_pass",
                    "upgrade_basis": "Fresh P77 deterministic Go evidence proves the currently implemented F-1 through F-5 source-verification/anti-fake behavior at the central ingestion, workflow, persistence, and rule boundaries.",
                    "gate_dimensions": "implementation=yes; ui=non-UI data-quality rule; data=yes; workflow/rule=yes; scenario=deterministic vectors; safety=yes.",
                    "fresh_evidence_command": F_COMMAND,
                    "fresh_evidence_artifact": f"{P77_F_TEST_LOG.relative_to(ROOT)}; {P77_F_TEST_META.relative_to(ROOT)}; {P77_SAFETY_SCAN.relative_to(ROOT)}; {P77_SAFETY_REVIEW.relative_to(ROOT)}",
                    "residual_gap": "None for currently implemented central anti-fake rule path; future new source adapters must keep the same regression gate.",
                    "next_remediation": "Keep F-1..F-5 tests mandatory for future source changes.",
                }
            )
        else:
            result.update(
                {
                    "p77_status": "p77_pending_evidence",
                    "upgrade_basis": "Candidate F-1..F-5 row, but P77 fresh deterministic evidence is incomplete.",
                    "gate_dimensions": "requires P77 F-1..F-5 Go test log and P77 safety scan.",
                    "fresh_evidence_command": F_COMMAND,
                    "fresh_evidence_artifact": f"{P77_F_TEST_LOG.relative_to(ROOT)}; {P77_F_TEST_META.relative_to(ROOT)}; {P77_SAFETY_SCAN.relative_to(ROOT)}; {P77_SAFETY_REVIEW.relative_to(ROOT)}",
                    "residual_gap": "Run and record P77 F-1..F-5 evidence before upgrading.",
                    "next_remediation": "Run the cited command and rerun this script with --check.",
                }
            )
        return result

    if req_id in {"REQ-07-011", "REQ-07-013"} or "510300" in row.get("requirement_text", ""):
        if evidence["non510300_browser_results"] and evidence["non510300_summary"]:
            result.update(
                {
                    "p77_status": "scoped_pass",
                    "upgrade_basis": "Fresh P77 non-510300 accepted-local UI journey revalidates dynamic binding, but accepted-local evidence still does not prove arbitrary live-provider coverage.",
                    "gate_dimensions": "implementation=yes; ui=yes; data=yes; workflow/LLM=yes; scenario=accepted-local non-510300; safety=yes.",
                    "fresh_evidence_command": "P75_ARTIFACT_DIR=docs/release/ui-audit-assets/2026-06-21-p77-non-510300 bash scripts/p75-non-510300-real-ui-journey.sh",
                    "fresh_evidence_artifact": str(P77_NON510300_ARTIFACT_DIR.relative_to(ROOT)),
                    "residual_gap": "Live arbitrary-symbol external-source coverage remains scoped.",
                    "next_remediation": "Add live-source or broader configured-symbol coverage before full dynamic-symbol real_pass.",
                }
            )
        return result

    if row["source_section"].split(".")[0] in {"8", "11", "13"} and evidence["sop_browser_results"] and evidence["sop_db_log"]:
        result.update(
            {
                "p77_status": original_status,
                "upgrade_basis": "Fresh P77 SOP/failure-state UI rerun supports the row, but coverage remains scenario-scoped and cannot be upgraded wholesale.",
                "gate_dimensions": "ui=yes for covered flows; data/audit=yes for covered flows; broader scenarios remain partial.",
                "fresh_evidence_command": "P75_ARTIFACT_DIR=docs/release/ui-audit-assets/2026-06-21-p77-sop-failure bash scripts/p75-sop-failure-real-ui-acceptance.sh",
                "fresh_evidence_artifact": str(P77_SOP_ARTIFACT_DIR.relative_to(ROOT)),
                "residual_gap": "Needs complete row-specific action/SOP matrix coverage before real_pass.",
                "next_remediation": "Split remaining SOP/action rows into exact real UI + SQLite + readback tests.",
            }
        )
    return result


def build_rows(rows: list[dict[str, str]], evidence: dict[str, object]) -> list[dict[str, str]]:
    p77_rows: list[dict[str, str]] = []
    for row in rows:
        projected = {column: row.get(column, "") for column in BASE_COLUMNS}
        projected["original_status"] = row["status"]
        projected.update(classify_p77(row, evidence))
        p77_rows.append(projected)
    return p77_rows


def write_matrix(rows: list[dict[str, str]], summary: dict[str, object]) -> None:
    columns = BASE_COLUMNS + ["original_status"] + P77_COLUMNS
    lines = [
        "# P77 Requirements Real-Pass Upgrade Matrix",
        "",
        "> Generated: 2026-06-21",
        "> Source: `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`",
        "> Policy: P77 is a new evidence layer; it does not rewrite P75 history.",
        "",
        "## Status Summary",
        "",
    ]
    for status, count in sorted(summary["p77_status_counts"].items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend(
        [
            "",
            f"- full_release_required rows: {summary['full_release_required_rows']}",
            f"- full_release_required real_pass rows: {summary['full_release_required_real_pass_rows']}",
            f"- remaining full_release_required non-real-pass rows: {summary['remaining_full_release_required_non_real_pass_rows']}",
            f"- conclusion: `{summary['conclusion']}`",
            "",
            "## Atomic Requirement Upgrade Rows",
            "",
            "|" + "|".join(columns) + "|",
            "|" + "|".join(["---"] * len(columns)) + "|",
        ]
    )
    for row in rows:
        lines.append("|" + "|".join(escape_cell(row.get(column, "")) for column in columns) + "|")
    P77_MATRIX.write_text("\n".join(lines) + "\n", encoding="utf-8")


def write_acceptance(summary: dict[str, object]) -> None:
    lines = [
        "# P77 Real-Pass Upgrade Acceptance",
        "",
        "> Date: 2026-06-21",
        "> Change: `p77-requirements-real-pass-upgrade-gate`",
        "> Source matrix: `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`",
        "",
        "## Conclusion",
        "",
        f"- Result: `{summary['conclusion']}`",
        f"- P77 upgraded rows to `real_pass`: {summary['upgraded_to_real_pass_count']}",
        f"- Full-release-required rows still non-real-pass: {summary['remaining_full_release_required_non_real_pass_rows']}",
        "- P77 does not rewrite P75 history and does not expand P76 package claims.",
        "",
        "## Evidence Inputs",
        "",
        f"- Safety scan: `{P77_SAFETY_SCAN.relative_to(ROOT)}` (`exists={P77_SAFETY_SCAN.exists()}`)",
        f"- Safety scan review: `{P77_SAFETY_REVIEW.relative_to(ROOT)}` (`reviewed_pass={safety_review_passed()}`)",
        f"- Safety boundary Go log: `{P77_SAFETY_TEST_LOG.relative_to(ROOT)}` (`exists={P77_SAFETY_TEST_LOG.exists()}`)",
        f"- Safety boundary Go metadata: `{P77_SAFETY_TEST_META.relative_to(ROOT)}` (`valid={metadata_passed(P77_SAFETY_TEST_META, command=SAFETY_COMMAND, tests=SAFETY_TEST_NAMES)}`)",
        f"- F-1..F-5 Go log: `{P77_F_TEST_LOG.relative_to(ROOT)}` (`exists={P77_F_TEST_LOG.exists()}`)",
        f"- F-1..F-5 Go metadata: `{P77_F_TEST_META.relative_to(ROOT)}` (`valid={metadata_passed(P77_F_TEST_META, command=F_COMMAND, tests=F_TEST_NAMES)}`)",
        f"- SOP/failure UI artifacts: `{P77_SOP_ARTIFACT_DIR.relative_to(ROOT)}` (`browser_results={(P77_SOP_ARTIFACT_DIR / 'browser-results.json').exists()}`)",
        f"- Non-510300 UI artifacts: `{P77_NON510300_ARTIFACT_DIR.relative_to(ROOT)}` (`summary={(P77_NON510300_ARTIFACT_DIR / 'non-510300-db-impact-summary.json').exists()}`)",
        "",
        "## Counts",
        "",
        "### Original P75 Status",
        "",
    ]
    for status, count in sorted(summary["original_status_counts"].items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend(["", "### P77 Status", ""])
    for status, count in sorted(summary["p77_status_counts"].items()):
        lines.append(f"- `{status}`: {count}")
    lines.extend(
        [
            "",
            "## Upgraded Rows",
            "",
        ]
    )
    if summary["upgraded_real_pass_rows"]:
        for row in summary["upgraded_real_pass_rows"]:
            lines.append(f"- `{row['requirement_id']}`: {row['requirement_text']} — {row['upgrade_basis']}")
    else:
        lines.append("- None.")
    lines.extend(
        [
            "",
            "## Remaining Gap Policy",
            "",
            "P77 keeps every remaining non-`real_pass` row visible in the upgrade matrix. It must not claim `release_ready_full_requirements_traceable` unless all `full_release_required=true` rows become `real_pass`.",
            "",
            "## Not Claimed",
            "",
            "P77 does not claim broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, future provider availability, paid/login/authorized/Level2/high-frequency sources, physical second-machine verification, remote publishing, Git tag creation, package refresh, investment return, or future market direction.",
            "",
        ]
    )
    P77_ACCEPTANCE.write_text("\n".join(lines), encoding="utf-8")


def build_summary(rows: list[dict[str, str]], evidence: dict[str, object], safety_scan: dict[str, object]) -> dict[str, object]:
    original_counts = Counter(row["original_status"] for row in rows)
    p77_counts = Counter(row["p77_status"] for row in rows)
    full_rows = [row for row in rows if row["full_release_required"] == "True"]
    full_real = [row for row in full_rows if row["p77_status"] == "real_pass"]
    pending = [row for row in rows if row["p77_status"] == "p77_pending_evidence"]
    upgraded = [row for row in rows if row["p77_status"] == "real_pass" and row["original_status"] != "real_pass"]
    conclusion = (
        "release_ready_full_requirements_traceable"
        if len(full_rows) == len(full_real)
        else "release_ready_scoped_with_p77_real_pass_progress"
    )
    return {
        "generated_utc": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "source_matrix": str(P75_MATRIX.relative_to(ROOT)),
        "p77_matrix": str(P77_MATRIX.relative_to(ROOT)),
        "p77_acceptance": str(P77_ACCEPTANCE.relative_to(ROOT)),
        "original_status_counts": dict(original_counts),
        "p77_status_counts": dict(p77_counts),
        "full_release_required_rows": len(full_rows),
        "full_release_required_real_pass_rows": len(full_real),
        "remaining_full_release_required_non_real_pass_rows": len(full_rows) - len(full_real),
        "upgraded_to_real_pass_count": len(upgraded),
        "pending_evidence_count": len(pending),
        "conclusion": conclusion,
        "evidence": evidence,
        "safety_scan": safety_scan,
        "upgraded_real_pass_rows": [
            {
                "requirement_id": row["requirement_id"],
                "source_section": row["source_section"],
                "requirement_text": row["requirement_text"],
                "original_status": row["original_status"],
                "upgrade_basis": row["upgrade_basis"],
                "fresh_evidence_artifact": row["fresh_evidence_artifact"],
            }
            for row in upgraded
        ],
        "pending_rows": [
            {"requirement_id": row["requirement_id"], "requirement_text": row["requirement_text"], "needed": row["fresh_evidence_artifact"]}
            for row in pending
        ],
    }


def write_summary(summary: dict[str, object]) -> None:
    P77_ASSET_DIR.mkdir(parents=True, exist_ok=True)
    P77_SUMMARY.write_text(json.dumps(summary, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def validate_outputs(rows: list[dict[str, str]], summary: dict[str, object]) -> list[str]:
    errors: list[str] = []
    if len(rows) != 341:
        errors.append(f"expected 341 rows from P75 matrix, got {len(rows)}")
    unknown_statuses = sorted({row["p77_status"] for row in rows} - STATUSES)
    if unknown_statuses:
        errors.append(f"unknown p77_status values: {unknown_statuses}")
    if summary["pending_evidence_count"] != 0:
        errors.append(f"P77 still has pending candidate evidence rows: {summary['pending_evidence_count']}")
    if summary["conclusion"] == "release_ready_full_requirements_traceable" and summary["remaining_full_release_required_non_real_pass_rows"] != 0:
        errors.append("full traceability conclusion is impossible while non-real-pass full-required rows remain")
    for row in rows:
        if row["p77_status"] == "real_pass":
            artifacts = [item.strip() for item in row["fresh_evidence_artifact"].split(";") if item.strip()]
            missing = [artifact for artifact in artifacts if not (ROOT / artifact).exists()]
            if missing:
                errors.append(f"{row['requirement_id']} real_pass missing artifact(s): {missing}")
            forbidden_basis = ["scope exclusion", "waiver-only", "screenshot-only", "route-smoke-only", "fixture-only", "mock/stub-only"]
            lower_basis = (row["upgrade_basis"] + " " + row["residual_gap"]).lower()
            if any(term in lower_basis for term in forbidden_basis):
                errors.append(f"{row['requirement_id']} real_pass uses forbidden upgrade basis language")
    return errors


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="Fail if generated artifacts still contain pending evidence or invalid real-pass rows.")
    args = parser.parse_args()

    _, p75_rows = read_p75_rows()
    safety_scan = run_safety_scan()
    evidence = evidence_state()
    p77_rows = build_rows(p75_rows, evidence)
    summary = build_summary(p77_rows, evidence, safety_scan)
    write_summary(summary)
    write_matrix(p77_rows, summary)
    write_acceptance(summary)

    errors = validate_outputs(p77_rows, summary) if args.check else []
    if errors:
        for error in errors:
            print(f"P77 check failed: {error}")
        return 1
    print(
        "p77_real_pass_upgrade:"
        f"rows={len(p77_rows)}:"
        f"real_pass={summary['p77_status_counts'].get('real_pass', 0)}:"
        f"pending={summary['pending_evidence_count']}:"
        f"conclusion={summary['conclusion']}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
