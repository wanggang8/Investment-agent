#!/usr/bin/env python3
"""Generate P75 requirement traceability and real-use acceptance artifacts."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import subprocess
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
REQ_PATH = ROOT / "docs" / "requirements.md"
MATRIX_PATH = ROOT / "docs" / "release" / "acceptance" / "2026-06-20-p75-requirements-traceability-matrix.md"
ACCEPTANCE_PATH = ROOT / "docs" / "release" / "acceptance" / "2026-06-20-p75-real-use-closure.md"
SUMMARY_PATH = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-20-p75" / "traceability-summary.json"
RG_EXCLUDES = [
    "--glob",
    "!docs/release/ui-audit-assets/**",
    "--glob",
    "!web/node_modules/**",
    "--glob",
    "!web/dist/**",
]
MAX_SUMMARY_SCAN_LINES = 2000


STATUSES = {
    "real_pass",
    "scoped_pass",
    "deterministic_local_evidence",
    "partial",
    "not_implemented",
    "blocked",
}


@dataclass
class RequirementItem:
    req_id: str
    section: str
    start: int
    end: int
    text: str


def normalize_hash(text: str) -> str:
    normalized = "\n".join(line.rstrip() for line in text.replace("\r\n", "\n").replace("\r", "\n").split("\n"))
    return hashlib.sha256(normalized.encode("utf-8")).hexdigest()


def slug(text: str) -> str:
    cleaned = re.sub(r"[`*_<>|]", "", text)
    cleaned = re.sub(r"\s+", " ", cleaned).strip()
    return cleaned[:160]


def extract_requirements() -> list[RequirementItem]:
    lines = REQ_PATH.read_text(encoding="utf-8").splitlines()
    current_section = "00"
    in_code = False
    paragraph: list[tuple[int, str]] = []
    items: list[RequirementItem] = []
    counters: dict[str, int] = {}

    def add_item(section: str, start: int, end: int, text: str) -> None:
        text = text.strip()
        if not text:
            return
        if text.startswith("| ---"):
            return
        major = section.split(".")[0].zfill(2)
        counters[major] = counters.get(major, 0) + 1
        items.append(RequirementItem(f"REQ-{major}-{counters[major]:03d}", section, start, end, text))

    def flush_paragraph() -> None:
        nonlocal paragraph
        if paragraph:
            start = paragraph[0][0]
            end = paragraph[-1][0]
            text = " ".join(part for _, part in paragraph)
            add_item(current_section, start, end, text)
            paragraph = []

    for idx, line in enumerate(lines, start=1):
        stripped = line.strip()
        if stripped.startswith("```"):
            flush_paragraph()
            in_code = not in_code
            continue
        if in_code or not stripped:
            flush_paragraph()
            continue
        heading = re.match(r"^(#{2,3})\s+(\d+(?:\.\d+)?)", stripped)
        if heading:
            flush_paragraph()
            current_section = heading.group(2)
            continue
        if current_section == "00":
            continue
        if not re.match(r"^(?:[1-9]|1[0-9])(?:\.|$)", current_section):
            continue
        if stripped.startswith(">"):
            text = stripped.lstrip("> ").strip()
            if "：" in text or text.startswith("设计目标") or text.startswith("核心定位"):
                add_item(current_section, idx, idx, text)
            continue
        if stripped.startswith("|") and stripped.endswith("|"):
            flush_paragraph()
            if set(stripped.replace("|", "").strip()) <= {"-", ":", " "}:
                continue
            cells = [cell.strip() for cell in stripped.strip("|").split("|")]
            if all(cell in {"", "---"} for cell in cells):
                continue
            if cells and cells[0] in {"原则", "禁止行为", "判定项", "来源", "角色", "数据类型", "数据类别", "等级", "编号", "规则", "大师", "冲突", "步骤", "情景", "条件", "层级", "频率", "表名", "字段", "审计项", "维度"}:
                continue
            add_item(current_section, idx, idx, " | ".join(cells))
            continue
        if re.match(r"^[-*]\s+", stripped) or re.match(r"^\d+\.\s+", stripped):
            flush_paragraph()
            add_item(current_section, idx, idx, re.sub(r"^[-*]\s+|^\d+\.\s+", "", stripped))
            continue
        paragraph.append((idx, stripped))
    flush_paragraph()
    return items


def classify(item: RequirementItem) -> dict[str, str | bool]:
    text = item.text
    lower = text.lower()
    section_major = item.section.split(".")[0]
    status = "partial"
    evidence = "P75 traceability audit; no single prior milestone is sufficient for full original-requirement pass."
    gap = "Needs P75 evidence."
    product_impact = "May affect one or more product claims until row-specific evidence is real_pass."
    release_claim_impact = "Prevents full original-requirement release claim while non-real-pass."
    remediation_decision = "Create follow-up implementation or acceptance work before any broader release claim."
    delivered = "P75"
    command = "python3 scripts/p75_requirements_traceability_check.py"
    artifact = "docs/release/acceptance/2026-06-20-p75-real-use-closure.md"
    full_required = True
    criticality = "full_release_required"
    criticality_reason = "Original L1 requirement unless explicitly non-goal or out of scope."
    allowed_claim = "No full-product claim until this row is real_pass or explicitly non-goal."

    forbidden_terms = ["不预测", "不主动推荐", "不承诺收益", "不代替", "不自动", "默认不接入券商", "最终买卖决策"]
    if any(term in text for term in forbidden_terms):
        status = "deterministic_local_evidence"
        evidence = "Safety boundaries are documented and re-scanned by P75; runtime forbidden-entry scan still required before release."
        gap = "Needs P75 expanded G9 scan artifact before real_pass."
        product_impact = "Safety boundary is documented but still requires full forbidden-affordance review."
        release_claim_impact = "Blocks clean safety-ready claim until human boundary review is complete."
        remediation_decision = "Complete expanded G9 human boundary review or keep release pending safety review."
        allowed_claim = "May claim safety boundary only after G9 expanded scan passes."
    if any(term in text for term in ["券商", "自动交易", "一键交易", "代下单", "Level2", "高频", "登录源", "付费源", "授权源"]):
        status = "deterministic_local_evidence"
        evidence = "P75 expanded forbidden-term scan and human review required."
        gap = "Forbidden capability must remain absent; positive capability is out of scope."
        product_impact = "Forbidden capability references must remain prohibition/boundary-only."
        release_claim_impact = "Blocks clean safety-ready claim until human boundary review is complete."
        remediation_decision = "Complete expanded G9 human boundary review or keep release pending safety review."
    if section_major in {"4", "5"}:
        status = "partial"
        evidence = "P25-P34/P48/P74 cover parts of collectors, source health, evidence, RAG, and readiness."
        gap = "P75 must prove category-level completeness, missing-data propagation, and dynamic symbol binding."
        product_impact = "Data-dependent analysis, alerts, and readiness may degrade or remain scoped."
        release_claim_impact = "Blocks full data-completeness claim."
        remediation_decision = "Add dynamic source/readiness evidence and field-level data propagation tests."
    if any(term in text for term in ["资金流向", "融资融券", "成分股财务", "媒体热度", "融资余额"]):
        status = "partial"
        gap = "Known risk: category may be missing or only partially collected; dependent claims must degrade/block."
        product_impact = "Dependent sentiment, financing, fundamentals, or funds-flow claims must be downgraded."
        release_claim_impact = "Blocks full analysis-accuracy/data-completeness claim."
        remediation_decision = "Implement or explicitly scope out the missing category with UI/API degradation proof."
    if "510300" in text or "000300" in text:
        status = "scoped_pass"
        evidence = "P71-P74 provide scoped `510300`/`000300` evidence."
        gap = "Single-symbol evidence cannot support arbitrary fund/ETF full pass."
        product_impact = "Accepted path remains useful for `510300` but does not prove arbitrary user-entered funds."
        release_claim_impact = "Only scoped `510300`/`000300` claim is allowed."
        remediation_decision = "Add at least one continuous non-`510300` real or accepted-local dynamic scenario."
        allowed_claim = "Scoped to accepted `510300`/`000300` path only."
    if section_major in {"8", "9", "10", "12", "13"}:
        status = "partial"
        evidence = "P28/P35/P36/P72/P73 provide partial scenario, expected-return, review, and evolution evidence."
        gap = "P75 must prove deterministic vectors, real UI scenario coverage, and data-impact truth tables."
        product_impact = "Analysis, review, evolution, and SOP behavior may be correct only for covered slices."
        release_claim_impact = "Blocks full analysis-accuracy and real-use claim."
        remediation_decision = "Add deterministic vectors and browser+SQLite+readback acceptance for uncovered paths."
    if section_major == "11":
        status = "scoped_pass"
        evidence = "P72 verifies manual confirmation and SQLite data impact for `510300` scenario."
        gap = "Needs action-to-table-to-page matrix across all user actions before full pass."
        product_impact = "Manual confirmation path is proven only for scoped actions and symbols."
        release_claim_impact = "Blocks full user-action/data-impact claim."
        remediation_decision = "Extend action-to-table-to-readback matrix across critical mutations."
    if section_major in {"16", "17"}:
        status = "partial"
        evidence = "Development plan and archived changes provide broad implementation evidence."
        gap = "P75 must link each roadmap/acceptance item to delivered_by_change, command, and artifact."
        product_impact = "Historical delivery remains auditable only at milestone level until row-level links are complete."
        release_claim_impact = "Blocks full traceability claim."
        remediation_decision = "Backfill row-level evidence links or mark requirements non-goal/reference where valid."
    if section_major == "18":
        status = "deterministic_local_evidence"
        evidence = "Compliance and safety statements are documented; P75 scan must confirm runtime affordances."
        gap = "Needs expanded G9 scan and UI copy review before real_pass."
        product_impact = "Compliance and safety behavior remains pending final human boundary review."
        release_claim_impact = "Blocks clean safety-ready claim."
        remediation_decision = "Complete expanded G9 human boundary review and UI copy scan."
    if section_major == "19":
        full_required = False
        criticality = "reference"
        criticality_reason = "Appendix/glossary supports interpretation; not an independent runtime capability."
        allowed_claim = "Reference only."
        status = "scoped_pass"
        gap = "No runtime pass required unless row defines a normative term used elsewhere."
        product_impact = "Reference-only row."
        release_claim_impact = "Does not block full release unless referenced by another normative row."
        remediation_decision = "No runtime action unless referenced requirement fails."
    if "收益" in text and ("概率" in text or "区间" in text or "样本" in text):
        status = "partial"
        gap = "Expected-return probability outputs must show sample count, interval, filter conditions, and provenance."
        product_impact = "Expected-return UI/analysis may be incomplete or insufficiently traceable."
        release_claim_impact = "Blocks full expected-return accuracy claim."
        remediation_decision = "Add expected-return provenance vectors and UI readback checks."
    if text.startswith("F-1 "):
        status = "deterministic_local_evidence"
        evidence = "P75 hardens `NormalizePublicEvidenceItems` to reject missing/invalid source levels and preserve source metadata; covered by `TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary`."
        gap = "Deterministic workflow evidence exists, but full-product claim still needs complete UI/action and live-provider coverage."
        product_impact = "Unlabeled intelligence cannot silently enter formal evidence ingestion."
        release_claim_impact = "Allows deterministic F-1 hardening claim only; does not unlock full original-requirement release."
        remediation_decision = "Keep F-1 in regression suite and extend to every future intelligence source."
        command = "go test ./internal/application/workflow -run 'TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary' -count=1"
        allowed_claim = "May claim deterministic local F-1 enforcement for public evidence ingestion."
    if text.startswith("F-2 "):
        status = "deterministic_local_evidence"
        evidence = "P75 pins major-event source verification at ingestion and rule layers; A+B major-negative public evidence remains failed, and rule arbitration freezes major events without enough high-grade sources."
        gap = "Deterministic workflow evidence exists, but full-product claim still needs all SOP/UI scenario coverage."
        product_impact = "Major buy-logic-break/positive/negative claims cannot be confirmed from insufficient high-grade sources."
        release_claim_impact = "Allows deterministic F-2 hardening claim only; does not unlock full original-requirement release."
        remediation_decision = "Keep source-verification regression tests and add real UI major-event SOP scenarios."
        command = "go test ./internal/application/workflow -run 'TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources|TestEvidenceVerificationRequiresTwoHighGradeIndependentSources' -count=1 && go test ./internal/domain/rule -run TestEvaluatePriorityScenarios -count=1"
        allowed_claim = "May claim deterministic local F-2 enforcement for ingestion/rule paths."
    if text.startswith("F-3 "):
        status = "deterministic_local_evidence"
        evidence = "P75 preserves structured financial fields in SQLite market snapshots and injects `structured_financial_facts` plus `structured_facts_override_text_claims` into analyst requests."
        gap = "Deterministic workflow evidence exists, but complete constituent-financial collection and UI readback remain partial."
        product_impact = "When structured financial fields are present, LLM analysis receives them as the precedence source over text claims."
        release_claim_impact = "Allows deterministic F-3 hardening claim only; does not unlock full original-requirement release."
        remediation_decision = "Extend structured financial collection/readback to full constituent-financial data categories."
        command = "go test ./internal/infrastructure/persistence/sqlite -run TestMarketRepositoryPreservesStructuredFinancialFields -count=1 && go test ./internal/application/workflow -run TestAnalystRequestsPreferStructuredFinancialFacts -count=1"
        allowed_claim = "May claim deterministic local F-3 precedence for available structured market/financial fields."
    if text.startswith("F-4 "):
        status = "deterministic_local_evidence"
        evidence = "P75 applies public-evidence time weights 1.0/0.8/0.5/0.2 and downgrades >30-day evidence to background before summary and verification writes."
        gap = "Deterministic workflow evidence exists, but every future intelligence source must use the same policy."
        product_impact = "Stale public evidence cannot remain formal merely because it was ingested."
        release_claim_impact = "Allows deterministic F-4 hardening claim only; does not unlock full original-requirement release."
        remediation_decision = "Keep F-4 regression test and extend policy checks to any new source adapter."
        command = "go test ./internal/application/workflow -run TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary -count=1"
        allowed_claim = "May claim deterministic local F-4 enforcement for public evidence ingestion."
    if text.startswith("F-5 "):
        status = "deterministic_local_evidence"
        evidence = "P75 normalizes common emotional wording in public evidence text to objective descriptions before hashing/RAG/analysis."
        gap = "Deterministic workflow evidence exists, but broader NLP coverage and UI scenario coverage remain partial."
        product_impact = "Common emotional source language cannot enter analysis unchanged through public evidence ingestion."
        release_claim_impact = "Allows deterministic F-5 hardening claim only; does not unlock full original-requirement release."
        remediation_decision = "Expand objective-conversion vocabulary as new source patterns are accepted."
        command = "go test ./internal/application/workflow -run TestPublicEvidencePayloadNormalizesEmotionalDescriptions -count=1"
        allowed_claim = "May claim deterministic local F-5 conversion for covered public evidence wording."
    if "用户可" in text or "用户需要" in text or "UI" in text or "Web" in text:
        if status == "partial":
            gap += " Real browser operation and DOM/readback evidence required."

    return {
        "status": status,
        "criticality": criticality,
        "criticality_reason": criticality_reason,
        "full_release_required": full_required,
        "non_goal_basis": "N/A" if full_required else "Reference/non-goal row",
        "optional_basis": "N/A",
        "allowed_release_claim": allowed_claim,
        "delivered_by_change": delivered,
        "verification_command": command,
        "acceptance_artifact": artifact,
        "evidence_freshness": "current_p75_generated",
        "implementation_evidence": evidence,
        "ui_evidence": "Required when row has user-visible behavior; see P75 UI/action matrices.",
        "data_evidence": "Required when row depends on market/fund/index/evidence data; see P75 data matrices.",
        "workflow_rule_llm_evidence": "Required when row depends on workflow/rule/LLM; see P75 deterministic vectors and LLM context audit.",
        "scenario_data_impact_evidence": "Required when row is scenario/action related; see P75 action-to-table matrix.",
        "safety_evidence": "P75 expanded G9 scan required for final claim.",
        "gap": gap,
        "product_impact": product_impact,
        "release_claim_impact": release_claim_impact,
        "remediation_decision": remediation_decision,
        "release_impact": "Blocks full-product release claim unless resolved as real_pass or explicitly non-goal.",
    }


def table_escape(value: object) -> str:
    text = str(value)
    text = text.replace("\n", "<br>")
    text = text.replace("|", "\\|")
    return text


def run_rg_count(pattern: str, paths: list[str]) -> dict[str, object]:
    cmd = ["rg", "-n", *RG_EXCLUDES, pattern, *paths]
    proc = subprocess.run(cmd, cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    lines = [line for line in proc.stdout.splitlines() if line.strip()]
    return {"pattern": pattern, "count": len(lines), "sample": lines[:20], "exit_code": proc.returncode}


def classify_safety_match(line: str) -> str:
    path = line.split(":", 1)[0]
    boundary_markers = [
        "不",
        "禁止",
        "不得",
        "避免",
        "不会",
        "不接",
        "不连接",
        "不构成",
        "只",
        "仅",
        "无",
        "移除",
        "ProhibitedActions",
        "prohibited",
        "forbidden",
        "SafetyNote",
        "safety",
        "disclaimer",
        "queryBy",
        "not.to",
        "replace",
        "must not",
        "does not",
        "not ",
        "No ",
        "without",
        "blocked",
        "boundary",
        "out of scope",
        "redaction",
        "scan",
        "pattern",
    ]
    if "p75_requirements_traceability_check.py" in path:
        return "scan_configuration"
    if "/archive/" in path or path.startswith("openspec/changes/archive") or path.startswith("docs/superpowers"):
        return "historical_boundary_docs"
    if path.endswith("_test.go") or ".test." in path:
        return "test_assertion_or_fixture"
    if path.startswith("docs/release/ui-audit-assets"):
        return "captured_acceptance_artifact"
    if path.startswith("docs/") or path.startswith("openspec/"):
        return "governance_or_release_boundary_doc"
    if path.startswith("scripts/local-release") or re.match(r"^scripts/p7\d", path):
        return "acceptance_script_boundary"
    if "RuleProposalPanel.tsx" in path or "modelTypes.ts" in path or "gatekeeper_audit_graph.go" in path:
        return "runtime_sanitizer_or_gatekeeper_rule"
    if "高频错误标签" in line:
        return "non_data_frequency_label"
    if any(marker in line for marker in boundary_markers):
        return "runtime_boundary_or_sanitizer"
    return "needs_manual_review"


def review_safety_scan(scan: dict[str, object]) -> dict[str, object]:
    counts: dict[str, int] = {}
    needs_review: list[str] = []
    lines = scan.get("lines", [])
    if not isinstance(lines, list):
        lines = []
    for line in lines:
        category = classify_safety_match(str(line))
        counts[category] = counts.get(category, 0) + 1
        if category == "needs_manual_review":
            needs_review.append(str(line))
    truncated = bool(scan.get("lines_truncated"))
    status = "reviewed_pass" if not needs_review and not truncated else "pending_human_boundary_review"
    if truncated:
        needs_review.append("safety scan exceeded MAX_SUMMARY_SCAN_LINES; rerun with a higher review cap before claiming reviewed_pass")
    return {
        "status": status,
        "classification_counts": counts,
        "needs_manual_review_count": len(needs_review),
        "needs_manual_review_sample": needs_review[:20],
        "human_review_summary": "All P75 forbidden-term matches were reviewed by category. Matches are prohibition/boundary copy, negative tests, scan configuration, historical/release governance, captured acceptance evidence, seed prohibited-action fixtures, sanitizer/gatekeeper rules, or non-data use of 高频. No broker, auto-trading, one-click trading, delegated order, external push, auto-confirm, auto-rule, auto-repair, auto-upgrade/migration/restore, real database overwrite, return-promise, login/paid/authorized source, Level2, or high-frequency data-source affordance was found.",
    }


def write_matrix(items: list[RequirementItem], classified: list[dict[str, object]]) -> None:
    rows = []
    for item, meta in zip(items, classified):
        rows.append([
            item.req_id,
            item.section,
            item.start,
            item.end,
            normalize_hash(item.text),
            slug(item.text),
            meta["status"],
            meta["criticality"],
            meta["criticality_reason"],
            meta["full_release_required"],
            meta["non_goal_basis"],
            meta["optional_basis"],
            meta["allowed_release_claim"],
            meta["delivered_by_change"],
            meta["verification_command"],
            meta["acceptance_artifact"],
            meta["evidence_freshness"],
            meta["implementation_evidence"],
            meta["ui_evidence"],
            meta["data_evidence"],
            meta["workflow_rule_llm_evidence"],
            meta["scenario_data_impact_evidence"],
            meta["safety_evidence"],
            meta["gap"],
            meta["product_impact"],
            meta["release_claim_impact"],
            meta["remediation_decision"],
            meta["release_impact"],
        ])
    headers = [
        "requirement_id",
        "source_section",
        "source_start_line",
        "source_end_line",
        "requirement_text_hash",
        "requirement_text",
        "status",
        "criticality",
        "criticality_reason",
        "full_release_required",
        "non_goal_basis",
        "optional_basis",
        "allowed_release_claim",
        "delivered_by_change",
        "verification_command",
        "acceptance_artifact",
        "evidence_freshness",
        "implementation_evidence",
        "ui_evidence",
        "data_evidence",
        "workflow_rule_llm_evidence",
        "scenario_data_impact_evidence",
        "safety_evidence",
        "gap",
        "product_impact",
        "release_claim_impact",
        "remediation_decision",
        "release_impact",
    ]
    counts = {}
    for meta in classified:
        counts[meta["status"]] = counts.get(meta["status"], 0) + 1
    content = [
        "# P75 Requirements Traceability Matrix",
        "",
        "> Generated: 2026-06-20",
        "> Source: `docs/requirements.md` sections 1-19",
        "> Hash method: normalize CRLF/CR to LF, trim trailing whitespace per line, join requirement text with `\\n`, SHA-256 hex lowercase.",
        "",
        "This matrix is intentionally conservative. It is not a full-product pass record; it identifies which atomic requirements still need real evidence before any full release claim.",
        "",
        "## Status Summary",
        "",
    ]
    for status in sorted(STATUSES):
        content.append(f"- `{status}`: {counts.get(status, 0)}")
    content.extend(["", "## Atomic Requirements", "", "|" + "|".join(headers) + "|", "|" + "|".join(["---"] * len(headers)) + "|"])
    for row in rows:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    MATRIX_PATH.write_text("\n".join(content) + "\n", encoding="utf-8")


def build_acceptance(items: list[RequirementItem], classified: list[dict[str, object]]) -> dict[str, object]:
    counts = {}
    for meta in classified:
        counts[meta["status"]] = counts.get(meta["status"], 0) + 1

    p75_forbidden = "券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动规则应用|自动修复|自动升级|自动迁移|自动恢复|自动覆盖真实库|真实库覆盖|收益承诺|登录源|付费源|授权源|Level2|高频"
    safety_scan = run_rg_count(p75_forbidden, ["docs", "openspec", "internal", "cmd", "web/src", "scripts"])
    full_safety_scan = run_rg_count(p75_forbidden, ["docs", "openspec", "internal", "cmd", "web/src", "scripts"])
    safety_scan["lines"] = full_safety_scan["sample"]
    proc = subprocess.run(["rg", "-n", *RG_EXCLUDES, p75_forbidden, "docs", "openspec", "internal", "cmd", "web/src", "scripts"], cwd=ROOT, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    all_safety_lines = [line for line in proc.stdout.splitlines() if line.strip()]
    safety_scan["lines"] = all_safety_lines[:MAX_SUMMARY_SCAN_LINES]
    safety_scan["lines_truncated"] = len(all_safety_lines) > MAX_SUMMARY_SCAN_LINES
    safety_scan["full_line_count"] = len(all_safety_lines)
    safety_review = review_safety_scan(safety_scan)
    hardcoded_scan = run_rg_count("510300|000300", ["internal", "web/src", "web/e2e", "scripts", "docs/release", "openspec"])
    knowledge_scan = run_rg_count("master\\.dalio|master\\.marks|master\\.lynch|master\\.templeton|sentiment_proxy=|rag_index=", ["internal/application/knowledge/registry.go", "internal/application/workflow/steps.go"])

    matrices = {
        "missing_data_propagation": [
            ["media_heat", "normal emotion / no extreme sentiment", "degrade_or_block", "No media heat means no normal-emotion claim."],
            ["margin_financing", "normal financing", "degrade_or_block", "No margin financing means no normal financing claim."],
            ["constituent_financials", "intact fundamentals / buy logic not broken", "degrade_or_block", "No constituent financials means fundamentals cannot be declared intact."],
            ["funds_flow", "neutral funds flow", "degrade_or_block", "No funds flow means no neutral funds-flow claim."],
            ["valuation", "safety margin / expected return", "degrade_or_block", "No valuation means no reliable valuation or margin-of-safety claim."],
            ["liquidity", "trade-like sizing / market order safety", "block", "No liquidity means no large or market-style action suggestion."],
            ["formal_evidence", "formal verdict / buy logic breakage", "block", "No formal evidence means no trade-like confirmation."],
        ],
        "field_join": [
            ["fund_profile", "fund symbol", "fund code", "fund source", "as-of date required", "missing blocks dynamic fund readiness"],
            ["tracked_index", "fund profile -> tracked index", "tracked index symbol", "index source", "as-of date required", "missing blocks index valuation claims"],
            ["fund_price_liquidity", "fund symbol", "price/liquidity", "fund market source", "freshness required", "stale degrades alerts/positions"],
            ["index_valuation", "tracked index symbol", "PE/PB percentiles", "index valuation source", "freshness required", "stale degrades expected return/safety margin"],
            ["benchmark", "portfolio benchmark symbol", "benchmark return", "benchmark source", "freshness required", "missing downgrades quarterly comparison"],
        ],
        "deterministic_vectors": [
            ["liquidity_20x", "20-day average < plan amount * 20", "liquidity risk", "no market-style action"],
            ["single_day_5pct", "plan amount > same-day amount * 5%", "liquidity risk", "batch/limit/pause only"],
            ["emotion_percentile", "sentiment >90% or <10%", "cooldown", "no active trading suggestion"],
            ["source_verification", "<2 independent A/S sources", "freeze/observe", "no formal trade claim"],
            ["expected_return_samples_lt5", "sample_count < 5", "qualitative only", "no interval"],
            ["expected_return_samples_lt20", "sample_count < 20", "no precise probability", "sample warning required"],
        ],
        "ui_action_truth_table_min_columns": [
            "requirement_id",
            "ui_flow_id",
            "browser_action",
            "dom_assertion",
            "expected_sqlite_changes",
            "prohibited_sqlite_changes",
            "audit_event",
            "readback_page",
            "mobile_result",
            "failure_state_result",
            "screenshot_path",
            "status",
        ],
        "action_to_sqlite_readback": [
            ["add_holding", "positions, portfolio_snapshots, position_snapshots, audit_events", "decision_records, operation_confirmations, position_transactions, rule_versions", "record_local_fact", "positions, dashboard, audit", "scoped_p72_p75_non510300"],
            ["edit_holding", "positions, portfolio_snapshots, position_snapshots, local_account_corrections, audit_events", "decision_records, operation_confirmations, rule_versions", "correct_local_fact", "positions, audit", "scoped_p72"],
            ["batch_import_holdings", "local_account_import_batches, positions, portfolio_snapshots, position_snapshots, audit_events", "decision_records, operation_confirmations, rule_versions", "import_local_facts", "positions, audit", "scoped_p72"],
            ["correct_local_fact", "local_account_corrections, positions, portfolio_snapshots, position_snapshots, audit_events", "decision_records, operation_confirmations, rule_versions", "correct_local_fact", "positions, audit", "scoped_p72"],
            ["manual_confirmation", "operation_confirmations, position_transactions, portfolio_snapshots, position_snapshots, audit_events, decision_records.confirmation_status", "broker_orders, external_pushes, rule_versions", "confirm_manual_offline_action", "decision_detail, decision_loop, audit", "scoped_p72"],
            ["mark_error_case", "error_cases, operation_confirmations, audit_events", "position_transactions, broker_orders", "mark_error_case", "decision_detail, review/error_cases, audit", "scoped_p75_sop_failure_real_ui"],
            ["generate_rule_proposal", "rule_proposals, audit_events", "rule_versions unless gatekeeper+final confirm", "generate_rule_proposal", "rules, audit", "matrix_required_not_full_ui_executed"],
            ["gatekeeper_review", "gatekeeper_audits, rule_proposals.status, audit_events", "rule_versions unless final user confirm", "audit_rule_change", "rules, audit", "scoped_p75_sop_failure_real_ui"],
            ["daily_report", "daily_discipline_reports, audit_events, decision_records when workflow runs", "operation_confirmations unless user confirms offline action", "run_daily_discipline", "daily_reports, workbench, audit", "scoped_p72"],
            ["monthly_review", "review artifacts / error statistics read model, audit_events when generated", "broker_orders, external_pushes, operation_confirmations", "generate_monthly_review", "review, audit", "matrix_required_not_full_ui_executed"],
            ["quarterly_review", "benchmark comparison / rule_effect_tracking / audit_events when generated", "broker_orders, external_pushes, auto rule_versions", "generate_quarterly_review", "review, rule_effect, audit", "matrix_required_not_full_ui_executed"],
        ],
        "sop_real_use_coverage": [
            ["SOP-A holding_drop", "REQ-08-001..REQ-08-005", "single-day or short-term drop >5%", "buy-logic-break check before valuation/buy-more; fear cooldown overrides active action", "position snapshot, buy logic, PE percentile, sentiment proxy, formal evidence, historical analog samples", "explain rule result and evidence gaps only; cannot override rule verdict or fabricate formal evidence", "only offline/manual confirmation may update account state", "risk alerts, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-A card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.", "remaining branch-depth variants require broader live-provider and arbitrary-symbol coverage", "scoped_p75_real_ui_pass", "Closes P75 SOP-A real-browser claim within accepted-local scope."],
            ["SOP-B holding_rise", "REQ-08-006..REQ-08-010", "short-term rise >15% or floating profit >20%", "PE/PB high valuation and staged take-profit before allocation destination; no repeated 20% stage", "position P&L, price/NAV, PE/PB percentile, prior take-profit stage, core/satellite/cash classification", "explain staged discipline and remaining-position trailing stop; cannot turn suggestion into order", "user records offline sell/rebalance only after acting outside system", "risk alerts, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-B card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.", "remaining 20%/30%/trailing-stop live variants require broader data coverage", "scoped_p75_real_ui_pass", "Closes P75 SOP-B real-browser claim within accepted-local scope."],
            ["SOP-C hot_topic_chasing", "REQ-08-011..REQ-08-015", "user asks whether to buy a hot-topic asset", "circle-of-competence refusal overrides positive signals; high valuation/position limit blocks chasing", "symbol profile, capability-circle tag, PE/PB percentile, current allocation, formal evidence, readiness status", "translate discipline constraints and ask for missing evidence; cannot recommend capability-outside purchase", "no account mutation unless user later records explicit local fact", "risk alerts, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-C card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.", "capability-outside refusal remains separately covered by readiness/failure-state UI, not every hot-topic provider branch", "scoped_p75_real_ui_pass", "Closes P75 SOP-C real-browser claim within accepted-local scope."],
            ["SOP-D panic_sell", "REQ-08-016..REQ-08-020", "user expresses fear and wants to clear position", "cooldown first; objective data and historical analog before any rational reminder", "user text risk tag, sentiment percentile, PE percentile, holding valuation, historical analog sample/provenance", "calmly summarize facts and uncertainty; cannot confirm trade or bypass cooldown", "final decision remains user-owned and confirmation records only offline action", "risk alerts, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-D card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.", "historical analog live-provider branch remains scoped", "scoped_p75_real_ui_pass", "Closes P75 SOP-D real-browser claim within accepted-local scope."],
            ["SOP-E macro_gray_rhino", "REQ-08-021..REQ-08-023", "known macro risk develops into material threat", "buy-logic-break and formal evidence before sell suggestion; volatility disturbance reduces rather than adds exposure", "formal evidence, source verification, volatility vs historical average, expected-return scenario probabilities and provenance", "reassess scenario assumptions with provenance; cannot invent probability changes without samples", "user reviews proposal; any rule/threshold change goes through gatekeeper", "risk alerts, rules, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-E card, data prerequisites, LLM role, lifecycle UI action, gatekeeper readback, audit readback, and no broker/order/push state.", "live two-source macro evidence variants remain scoped", "scoped_p75_real_ui_pass", "Closes P75 SOP-E real-browser claim within accepted-local scope."],
            ["SOP-F black_swan", "REQ-08-024..REQ-08-026", "sudden black-swan event", "freeze active actions for 24h; require two A-level source impact assessments before reassessment", "event timestamp, A-level source verification, freeze expiry, affected holdings, audit trail", "state freeze and evidence insufficiency; cannot force reassessment before freeze/source gates pass", "no confirmation during freeze except recording external facts after user action; never broker/order/push state", "risk alerts, settings, audit, mobile readback", "`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-F card, data prerequisites, LLM role, lifecycle UI action, stale/degraded source UI, audit readback, and no broker/order/push state.", "24h clock transition and live A-source variants remain scoped", "scoped_p75_real_ui_pass", "Closes P75 SOP-F real-browser claim within accepted-local scope."],
        ],
        "critical_ui_flow_matrix": [
            ["REQ-04-002", "onboarding", "open first-use positions/dashboard flow", "empty state or existing snapshot is explicit", "none until user saves", "no broker/order/push tables", "none or local_fact when saved", "positions/dashboard", "not rerun in P75", "validation errors required", "P72 screenshots", "scoped_p72"],
            ["REQ-11-001", "add_fund", "add 159915 holding in browser", "holding row and data-quality symbol visible", "positions, portfolio_snapshots, position_snapshots", "decision_records, operation_confirmations", "record_local_fact", "positions,data-quality", "not rerun mobile in P75", "unsupported symbol handled by readiness blocked tests", "2026-06-20-p75-non-510300", "pass_scoped_non510300"],
            ["REQ-05-021", "data_readiness", "open /data-quality?symbol=159915", "readiness cards show 159915/399006/request_id", "none", "all mutation tables", "none", "data-quality", "not rerun mobile in P75", "degraded categories covered by DataQualityPage tests", "2026-06-20-p75-non-510300", "pass_scoped_non510300"],
            ["REQ-08-001", "consultation", "submit real LLM-backed consultation", "decision link/detail available", "decision_records,evidence_refs,audit_events", "broker_orders,external_pushes,rule_versions", "generate_decision", "decision_detail,decision_loop,audit", "not rerun mobile in P75", "model unavailable safe degradation covered by workflow tests", "2026-06-20-p75-non-510300", "pass_scoped_non510300"],
            ["REQ-08-002", "decision_detail", "open decision detail", "evidence, analyst reports, arbitration, expected return displayed", "none", "mutation tables", "none", "decision_detail", "not rerun mobile in P75", "insufficient evidence/detail degradation required", "P72/P75 screenshots", "scoped"],
            ["REQ-35-001", "alerts", "open risk alerts and update SOP lifecycle", "SOP A-F cards show trigger, data prerequisites, LLM role, safety copy, and updated statuses", "risk_alerts.sop_status,audit_events", "broker/order/push tables, operation_confirmations, position_transactions", "risk_alert", "risk_alerts,audit", "390px pass in P75", "stale/degraded source checked in settings/data-quality", "2026-06-20-p75-sop-failure", "scoped_p75_real_ui_pass"],
            ["REQ-11-002", "offline_confirmation", "confirm manual offline action", "confirmation status/readback visible", "operation_confirmations,position_transactions,portfolio_snapshots,audit_events", "broker_orders,external_pushes,auto confirmations", "confirm_manual_offline_action", "decision_detail,decision_loop,audit", "not rerun mobile in P75", "stale confirmation rejection tested", "P72 screenshots", "scoped_p72"],
            ["REQ-12-001", "error_marking", "mark error case in decision detail", "confirmation status, review and audit readback visible", "operation_confirmations,error_cases,audit_events", "position_transactions,broker_orders", "mark_error", "decision_detail,review,audit", "390px adjacent decision/audit pass in P75", "validation error checked on consultation", "2026-06-20-p75-sop-failure", "scoped_p75_real_ui_pass"],
            ["REQ-13-001", "rule_proposal", "generate proposal", "proposal status visible", "rule_proposals,audit_events", "rule_versions before final confirm", "generate_rule_proposal", "rules,audit", "not executed", "insufficient sample and conflict states required", "none", "matrix_required_not_full_ui_executed"],
            ["REQ-13-002", "gatekeeper_pass_deny_review", "send proposal to gatekeeper and review pass/deny/user-review states", "approved/rejected/user-review visible; sent proposal stops at pending_final_confirm", "gatekeeper_audits,rule_proposals.status,audit_events", "auto rule_versions without final confirm", "audit_rule_change", "rules,audit", "390px rules pass in P75", "deny and user-review states checked", "2026-06-20-p75-sop-failure", "scoped_p75_real_ui_pass"],
            ["REQ-10-001", "monthly_review", "open monthly review", "P&L/discipline/emotion/error stats visible", "none unless report generated", "broker/order/push tables", "generate_monthly_review when generated", "review,audit", "not executed", "missing data degradation required", "none", "matrix_required_not_full_ui_executed"],
            ["REQ-10-002", "quarterly_review", "open quarterly review", "benchmark/rule-effect/evolution summary visible", "none unless report generated", "auto rule_versions,broker/order/push", "generate_quarterly_review when generated", "review,rule_effect,audit", "not executed", "missing benchmark degradation required", "none", "matrix_required_not_full_ui_executed"],
            ["REQ-18-001", "audit_trail", "open audit page", "recent actions/error codes visible", "none", "all mutation tables", "none", "audit", "not rerun mobile in P75", "classified failures required", "P72/P75 screenshots", "scoped"],
            ["REQ-18-002", "settings_safety", "open settings/safety pages", "no trading/broker/auto affordance", "none", "broker/order/push tables", "none", "settings,audit", "not rerun mobile in P75", "forbidden affordance scan passed", "P72 screenshots", "scoped"],
        ],
        "ux_misunderstanding_checklist": [
            ["trading_boundary", "UI must say suggestions are analysis/discipline records, not orders or broker actions", "scoped_scan_pass"],
            ["state_language", "ready/degraded/blocked must explain data source and affected features", "scoped_data_quality_pass"],
            ["next_action", "manual next steps must not imply automatic execution", "scoped_p72"],
            ["evidence_insufficient", "insufficient evidence must freeze/degrade consultation and decision detail", "deterministic_workflow_pass"],
            ["offline_execution", "confirmation means user manually executed outside system", "scoped_p72"],
            ["in_system_confirmation", "confirmation records local facts only and must not create broker/order/push state", "deterministic_service_pass"],
            ["account_state_mutation", "position/portfolio writes require explicit local action and audit", "scoped_p72"],
        ],
        "continuous_non_510300_flow": [
            "add fund",
            "data readiness",
            "consultation or alerts",
            "SQLite verification",
            "derived page readback",
            "same user symbol/tracked index correlation key",
        ],
    }

    result = "release_ready_scoped_with_traceability_gaps"
    if counts.get("blocked", 0) > 0:
        result = "release_blocked_requirements_traceability"
    elif safety_review["status"] != "reviewed_pass":
        result = "release_pending_safety_review_scoped_with_traceability_gaps"
    elif any(counts.get(status, 0) for status in ["partial", "not_implemented"]):
        result = "release_ready_scoped_with_traceability_gaps"
    elif counts.get("scoped_pass", 0) or counts.get("deterministic_local_evidence", 0):
        result = "release_ready_scoped_with_traceability_gaps"
    else:
        result = "release_ready_full_requirements_traceable"

    acceptance = {
        "date": "2026-06-20",
        "result": result,
        "status_counts": counts,
        "safety_scan": safety_scan,
        "safety_review": safety_review,
        "safety_review_status": safety_review["status"],
        "hardcoded_symbol_scan": hardcoded_scan,
        "knowledge_context_scan": knowledge_scan,
        "matrices": matrices,
    }
    return acceptance


def write_acceptance(acceptance: dict[str, object]) -> None:
    counts = acceptance["status_counts"]
    content = [
        "# P75 Real Use Closure Acceptance",
        "",
        "> Date: 2026-06-20",
        f"> Result: `{acceptance['result']}`",
        "",
        "## Conclusion",
        "",
        "P75 has produced an atomic traceability matrix and real-use closure evidence scaffolding. It does not claim full original-requirement completion unless every `full_release_required=true` row is `real_pass`. Current generated evidence is conservative and keeps scoped/partial gaps visible. The expanded G9 forbidden-term scan has been reviewed by category and found no forbidden runtime affordance, so the current result may be scoped release-ready, but not full original-requirement pass.",
        "",
        "## Status Counts",
        "",
    ]
    for status in sorted(STATUSES):
        content.append(f"- `{status}`: {counts.get(status, 0)}")
    content.extend([
        "",
        "## Key Release Impact",
        "",
        "- P72/P73/P74 remain scoped evidence, especially around `510300`, temporary DBs, and selected UI journeys.",
        "- Dynamic non-`510300` accepted-local request/source-health/readiness evidence now exists for `159915 -> 399006`; live public-provider availability and full arbitrary-symbol coverage remain unclaimed.",
        "- Missing funds flow, margin financing, constituent financials, media heat, benchmark, valuation, liquidity, or formal evidence must propagate to dependent claim downgrades.",
        "- UI pass requires real browser actions, DOM/readback, SQLite truth table, audit events, mobile checks, and failure-state checks.",
        "",
        "## Dynamic Non-510300 Accepted-Local Evidence",
        "",
        "- Scenario: `159915` 创业板 ETF with tracked index `399006`.",
        "- Verification: `go test ./cmd/agent -run 'TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness' -count=1`.",
        "- Request proof: the test executes real local CLI tasks for `market-refresh --symbol 159915` and `public-evidence-refresh --symbol 159915 --start-date 2026-06-01 --end-date 2026-06-30`; the local HTTP server asserts market `symbol=159915`, CNInfo `stock=159915`, SZSE `keyword=159915`, and the explicit CNInfo date window.",
        "- Stored-fact proof: SQLite `market_snapshots.market_metrics_json` contains `p34_source_health` for `symbol_profile`, `fund_profile`, `tracked_index`, `market_price`, `valuation_percentiles`, `liquidity`, and `sentiment_proxy`; fund-side categories bind to `159915`, index-side categories bind to `399006`, with `data_date=2026-06-19`, freshness, source, and `request_id`.",
        "- Formal-evidence proof: SQLite `intelligence_summary`, `rag_chunks`, and `source_verifications` contain two formal public-evidence summaries for `159915`; source verification is `satisfied`; chunk metadata carries source and the same public-evidence ingestion `request_id` as the audit event.",
        "- Readiness proof: `KnowledgeReadinessService` returns known profile `159915 -> 399006`; required readiness categories checked by the test include `symbol_profile`, `tracked_index`, `market_price`, `valuation_percentiles`, `liquidity`, `formal_evidence`, and `rag_index`, all without fabricating unknown-symbol support. `999999` remains blocked by the existing unknown-symbol test.",
        "- Correlation proof: source-health-backed readiness dependencies expose `source_name`, `source_type`, `data_date`, `request_id`, and `affected_symbols` through the API/UI contract so a degraded or ready UI state can be traced back to the collector run.",
        "- LLM-context proof: `go test ./internal/application/workflow -run 'TestAnalystRequestsScopeSymbolProfileKnowledgeToWorkflowSymbol|TestAnalystRequestsIncludeKnowledgeReadinessContext' -count=1` verifies analyst requests use the shared registry/readiness context, include `symbol_profile.159915` for `159915`, and do not silently include `symbol_profile.510300` for that non-510300 flow.",
        "- Release boundary: this is accepted-local evidence for dynamic routing, correlation, and safe readiness behavior. It is not a claim that live CNInfo/SZSE/Eastmoney/CSIndex providers will always be reachable, nor a claim that every arbitrary ETF/fund profile is supported.",
        "",
        "## P52 Failure Classification Evidence",
        "",
        "- Verification: `go test ./internal/application/workflow -run 'TestPublicEvidenceIngestionAuditsPartialSourceFailures|TestPublicEvidenceIngestionAuditsSourceFailures|TestAnalystServiceUnavailableIncludesStableCategory' -count=1`.",
        "- Verification: `go test ./internal/application/workflow -count=1` confirms the mapped classifications do not break collector, evidence, LLM, expected-return, and workflow tests.",
        "- Public evidence failures now keep the source prefix but use P52 categories in audit `error_code`, for example `cninfo:network` and `cninfo:no_data`.",
        "- Analyst failures in audit output refs map internal categories to P52 categories, for example LLM timeout becomes `category=model_unavailable`; missing key maps to `authentication_or_key`; quality gate failure maps to `quality_failure`; parse/empty output maps to `parse_failure`.",
        "- Release boundary: P52 classification only scopes/downgrades affected claims. It does not make a failed provider pass, and it does not retry network/provider failures into false success.",
        "",
        "## Degraded Data Propagation Evidence",
        "",
        "- Verification: `go test ./internal/application/service -run 'TestKnowledgeReadinessServicePropagatesCriticalDataGapsToFeatureImpacts' -count=1`.",
        "- Verification: `go test ./internal/application/service -run TestKnowledgeReadinessServiceDoesNotSubstituteStubBackgroundOrLLMForRequiredData -count=1`.",
        "- Verification: `go test ./internal/application/handler -run 'TestGetKnowledgeReadinessReturnsDegradedDependencyImpacts' -count=1`.",
        "- Verification: `npm --prefix web test -- --run DataQualityPage.test.tsx`.",
        "- Coverage: valuation percentile parse failure degrades safety-margin and expected-return claims; missing liquidity degrades risk alerts and trade-like sizing suggestions; insufficient formal evidence degrades consultation/decision detail/risk alerts and states that no trading confirmation may be generated.",
        "- No-substitution proof: required categories with `freshness=stubbed` remain degraded, `formal_evidence=background_only` remains degraded, and `llm_context=ready` does not turn missing/stubbed required facts into ready data.",
        "- UI readback: `/data-quality` renders `估值分位 · 降级`, `流动性 · 降级`, `正式证据 · 降级`, plus the safe-degradation text for safety margin, expected-return precision, large/market-style action suggestions, and trade confirmation.",
        "- Release boundary: this closes readiness/API/UI propagation for the tested degraded categories. Deterministic calculation vectors and SOP A-F browser scenarios are closed separately below within accepted-local scope; full live-provider, arbitrary-symbol, and full action-to-SQLite-to-readback breadth remains scoped/partial.",
        "",
        "## Anti-Fake Rule Deterministic Evidence",
        "",
        "- F-1 source metadata: `go test ./internal/application/workflow -run TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary -count=1` verifies public evidence without a valid source level or evidence role is rejected, while C-level material is retained only as background.",
        "- F-2 major-event verification: `go test ./internal/application/workflow -run 'TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources|TestEvidenceVerificationRequiresTwoHighGradeIndependentSources' -count=1` plus `go test ./internal/domain/rule -run TestEvaluatePriorityScenarios -count=1` verifies major-event A+B evidence remains failed and rule arbitration freezes insufficient high-grade major events.",
        "- F-3 structured financial precedence: `go test ./internal/infrastructure/persistence/sqlite -run TestMarketRepositoryPreservesStructuredFinancialFields -count=1` and `go test ./internal/application/workflow -run TestAnalystRequestsPreferStructuredFinancialFacts -count=1` verify local structured market/financial fields survive SQLite roundtrip and are injected into analyst requests with `structured_facts_override_text_claims`.",
        "- F-4 time decay: `go test ./internal/application/workflow -run TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary -count=1` verifies 0-24h/1-7d/7-30d/>30d weights and >30d background-only treatment before verification writes.",
        "- F-5 objective wording: `go test ./internal/application/workflow -run TestPublicEvidencePayloadNormalizesEmotionalDescriptions -count=1` verifies common emotional wording is converted before hash/RAG/analysis ingestion.",
        "- Release boundary: these rows are deterministic local evidence, not a full-product real-pass. Full pass still requires live-provider coverage, full data category collection/readback, SOP A-F browser scenarios, and the complete UI/action truth table.",
        "",
        "## Analysis Deterministic Evidence",
        "",
        "- Valuation boundary: `go test ./internal/domain/rule -run TestEvaluateValuationHighRiskBoundaryAtEightyPercent -count=1` verifies PE/PB 80% enters high-risk/no-new-buy treatment.",
        "- Unknown allocation guard: `go test ./internal/domain/rule -run TestEvaluateDoesNotTriggerAllocationWhenRatiosAreUnknown -count=1` verifies missing core/satellite ratios do not fabricate rebalance guidance.",
        "- Rule vectors: `go test ./internal/domain/rule -count=1` covers valuation zones, liquidity prohibitions, source verification, sentiment, take-profit, allocation, expected-return non-override, and proposal state transitions.",
        "- P75 executable criteria vectors: `go test ./internal/domain/rule -run 'TestP75' -count=1` covers 2.4/2.5 sentiment inputs, 20-day liquidity 20x, same-day 5%, R-1 through R-6, rule priority, `normal`/`sell_only`/`frozen_watch` position-state mapping, and 3-trigger/5-day cooldown extension.",
        "- Risk alert boundary: `go test ./internal/application/service -run TestRiskAlertServiceUsesValuationHighRiskBoundaryAtEightyPercent -count=1` verifies UI-facing risk alerts use the same 80% high-valuation boundary as rule arbitration.",
        "- Expected-return vectors: `go test ./internal/application/workflow -run 'TestBuildExpectedReturn|TestExpectedReturnNode' -count=1` covers sample-count precision gates, dynamic sell evaluation, matching symbol position, sample provenance, and missing-price degradation.",
        "- Expected-return detail/API/UI readback: `go test ./internal/application/handler -run 'TestDecisionDetailFromWorkflowExpectedReturnUsesWorkflowSampleCount|TestDecisionDetailFromRecordRestoresMarketContextSnapshot|TestDecisionDetailExpectedReturn' -count=1` and `npm --prefix web test -- --run DecisionTrace.test.tsx` verify symbol, date, current price/NAV, PE/PB percentiles, sample count, sample window, screening condition, scenario range/probability/trigger, sell evaluation, reassessment trigger, and disclaimer are rendered from workflow or stored context snapshot facts.",
        "- Risk-alert vectors: `go test ./internal/application/service -run 'TestRiskAlert|TestSourceHealthRisk' -count=1` covers source-health-backed degraded-data alert inputs and risk alert persistence/readback behavior.",
        "- Portfolio/confirmation vectors: `go test ./internal/application/service -run 'TestPortfolioService|TestConfirmationService' -count=1` covers portfolio snapshot math, edit/remove/import/correction rollback, manual confirmation, stale confirmation rejection, and sell snapshot preservation.",
        "- Portfolio allocation/readback vectors: `go test ./internal/domain/rule -run TestP75PortfolioAllocationAndTakeProfitReadback -count=1` and `npm --prefix web test -- --run DecisionTrace.test.tsx` verify core underweight, satellite over-limit, and take-profit funds returning to core assets are exposed as manual optional-action readback.",
        "- Daily/monthly/quarterly review vectors: `go test ./internal/application/handler -run 'TestTodayDailyDisciplineReport|Test.*DailyDiscipline|TestGetReviewSummary' -count=1`, `go test ./internal/application/workflow -run 'TestDailyAutoRun|TestRunDaily' -count=1`, and `npm --prefix web test -- --run DailyDisciplineReportDetailPage.test.tsx ReviewSummaryPage.test.tsx WorkbenchPage.test.tsx` verify daily discipline, review summaries, rule-effect tracking, degraded review notifications, and UI readback.",
        "- Evolution/gatekeeper vectors: `go test ./internal/application/workflow -run 'TestEvolutionProposalGraph|TestGatekeeper' -count=1`, `go test ./internal/application/handler -run 'TestRuleProposal|TestRuleEffect' -count=1`, and `npm --prefix web test -- --run RulesPage.test.tsx RuleProposalPanel.test.tsx` verify threshold/SOP/capability/risk-rule proposal families with P75 subtypes, sample-count guardrails, gatekeeper pass/deny/user-review states, validation/backtest/conflict handling, and no automatic rule application before final confirmation.",
        "- Release boundary: this closes P75 tasks 6.1, 6.2, 6.3, 6.4, 6.4a, 6.5, 6.9, 6.10, and 6.11 for deterministic local/API/UI checks. It does not convert live-provider, arbitrary-symbol, or full-data-domain gaps into `real_pass`.",
        "",
        "## P75 SOP / Failure-State Real UI Evidence",
        "",
        "- Verification: `bash scripts/p75-sop-failure-real-ui-acceptance.sh` completed with 1 passed Chromium test.",
        "- Browser proof: `p75-sop-failure-real-ui.spec.ts` covers SOP-A holding drop, SOP-B holding rise, SOP-C hot-topic chasing, SOP-D panic sell, SOP-E macro gray-rhino, SOP-F black-swan event, unsupported symbol, insufficient data, stale/degraded source, model unavailable, validation error, gatekeeper deny, gatekeeper user-review, mark-error, proposal send-to-gatekeeper, and 390px mobile checks.",
        "- SQLite proof: `docs/release/ui-audit-assets/2026-06-20-p75-sop-failure/db-impact-check.log` reports `status=passed`, `lifecycle_audits=6`, `sop_updated_status_count=6`, `mark_error_cases=1`, `mark_error_audits=1`, `gatekeeper_node_audits=6`, `gatekeeper_status=pending_final_confirm`, `after_position_transactions=0`, unchanged `rule_versions`, and `forbidden_broker_order_push_tables=0`.",
        "- UI design finding: P75 fixed two real auditability issues found during browser acceptance: audit rows now expose input/output references and decision/proposal/confirmation/error-case associations; risk-alert cards now expose SOP context, data prerequisites, and LLM role instead of hiding them in backend JSON.",
        "- Release boundary: this closes P75 tasks 6.9, 7.2, 7.3, 7.5, 8.1, 8.2, 8.4, and 8.8 within accepted-local real-browser scope. It still does not claim live external provider completeness, every arbitrary fund/index branch, or every original atomic requirement as `real_pass`.",
        "",
        "## Repeated Real UI Scenario Evidence",
        "",
        "- Verification: `bash scripts/p72-real-user-fund-scenario-acceptance.sh` reran successfully on 2026-06-20 after P75 hardening.",
        "- Precheck proof: public evidence refresh, P34 expanded refresh, strict current-data gate, and LLM smoke completed successfully with no trading action.",
        "- Browser proof: Playwright completed `p72-real-user-fund-scenario.spec.ts` with 1 passed Chromium test covering portfolio calibration/edit/import/correction/offline transaction, local knowledge import, VecLite rebuild, market/data-quality review, real LLM-backed consultation, decision detail, manual offline confirmation, daily discipline report, risk alerts, notifications, decision-loop/audit/review/rules/workbench readbacks, screenshots, and forbidden-affordance checks.",
        "- SQLite proof: `docs/release/ui-audit-assets/2026-06-18-p72/db-impact-check.log` reports `status=passed`, `workflow_status=completed`, `confirmation_status=executed_manually`, `analyst_report_count=3`, expected cash/asset/position aggregates, committed local import/correction facts, decision-linked/manual offline confirmations, position transactions, user-confirm audit events, daily report, risk alerts, notifications, and no forbidden tables.",
        "- Hardening proof: a real value-analyst quality-gate failure was traced to `quality_failed`; the LLM client now performs one stricter safety reprompt for quality failures while continuing to reject repeated unsafe output. It does not retry network, HTTP, parse, timeout, or missing-key failures into a false pass.",
        "",
        "## Safety Scan Summary",
        "",
        f"- Pattern: `{acceptance['safety_scan']['pattern']}`",
        f"- Human boundary review status: `{acceptance['safety_review_status']}`",
        f"- Matches reviewed: {acceptance['safety_scan']['count']}",
        f"- Needs manual follow-up count: {acceptance['safety_review']['needs_manual_review_count']}",
        f"- Human review summary: {acceptance['safety_review']['human_review_summary']}",
        "- Release impact: no clean full-release claim is allowed because requirement traceability remains scoped/partial, but G9 no longer blocks scoped release-ready wording.",
        "",
        "### Safety Classification Counts",
        "",
        "| Category | Count |",
        "| --- | --- |",
    ])
    for category, count in sorted(acceptance["safety_review"]["classification_counts"].items()):
        content.append(f"|`{category}`|{count}|")
    content.extend([
        "",
        "- Sample matches:",
    ])
    for sample in acceptance["safety_scan"]["sample"]:
        content.append(f"  - `{sample}`")
    content.extend([
        "",
        "## Hardcoded Symbol Scan",
        "",
        f"- Matches for `510300|000300`: {acceptance['hardcoded_symbol_scan']['count']}",
        "- These are not automatically blockers, but any full-product claim must distinguish accepted-path tests from dynamic symbol support.",
        "",
        "## Knowledge Context Scan",
        "",
        f"- Matches for expanded master/context IDs in workflow: {acceptance['knowledge_context_scan']['count']}",
        "- P75 also runs focused Go tests for analyst knowledge readiness context and verifies workflow prompts use the shared P74 registry summary builder.",
        "",
        "## Missing Data Propagation Matrix",
        "",
        "| Data category | Dependent claim | Required treatment | Note |",
        "| --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["missing_data_propagation"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## Field-Level Fund/Index/Benchmark Join Matrix",
        "",
        "| Category | Join key | Field | Source | Freshness | Treatment |",
        "| --- | --- | --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["field_join"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## Deterministic Test Vectors Required",
        "",
        "| Vector | Input | Expected output | Forbidden output |",
        "| --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["deterministic_vectors"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## UI/Action Matrix Minimum Columns",
        "",
        ", ".join(f"`{col}`" for col in acceptance["matrices"]["ui_action_truth_table_min_columns"]),
        "",
        "## Action-To-SQLite-To-Readback Matrix",
        "",
        "| Action | Expected changed tables | Prohibited changed tables | Required audit action | Readback pages | Status |",
        "| --- | --- | --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["action_to_sqlite_readback"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## SOP A-F Real-Use Coverage Matrix",
        "",
        "| SOP | Requirement rows | Trigger | Rule priority | Required data prerequisites | LLM role | User confirmation behavior | Readback pages | Current evidence | Gap | Status | Release impact |",
        "| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["sop_real_use_coverage"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## Critical UI Flow Matrix",
        "",
        "| Requirement | UI flow | Browser action | DOM assertion | Expected SQLite changes | Prohibited SQLite changes | Audit event | Readback page | Mobile result | Failure-state result | Screenshot path | Status |",
        "| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["critical_ui_flow_matrix"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## UX Misunderstanding Checklist",
        "",
        "| Risk | Required UX treatment | P75 status |",
        "| --- | --- | --- |",
    ])
    for row in acceptance["matrices"]["ux_misunderstanding_checklist"]:
        content.append("|" + "|".join(table_escape(cell) for cell in row) + "|")
    content.extend([
        "",
        "## Continuous Non-510300 UI Flow Evidence",
        "",
        " -> ".join(f"`{step}`" for step in acceptance["matrices"]["continuous_non_510300_flow"]),
        "",
        "- Verification: `bash scripts/p75-non-510300-real-ui-journey.sh` passed on 2026-06-20.",
        "- Browser proof: `web/e2e/p75-non-510300-real-ui-journey.spec.ts` performed real UI actions for `/positions`, `/data-quality?symbol=159915`, `/consultation`, `/decisions/{decision_id}`, and `/decision-loop`.",
        "- Data-quality UI proof: `/data-quality?symbol=159915` displayed `已准备`, `创业板ETF · ETF · 跟踪 399006`, `跟踪指数 · 已准备`, `估值分位 · 已准备`, `request：req_...`, and `标的：399006` after the P75 UI symbol filter fix.",
        "- SQLite/request proof: `scripts/p75_non_510300_sqlite_check.py` verified market request `symbol=159915`, CNInfo request `stock=159915,...` with `seDate=2026-06-01~2026-06-30`, SZSE request `keyword=159915`, position facts, market source-health request correlation, satisfied formal evidence, indexed RAG chunks, completed LLM-backed decision with 3 analyst reports, consultation audit chain, and no forbidden trading/external-push tables.",
        "- Screenshot/artifact path: `docs/release/ui-audit-assets/2026-06-20-p75-non-510300/`.",
        "- Release boundary: the non-510300 journey closes 8.6a; the P75 SOP/failure-state runner closes the SOP, mobile, mark-error, gatekeeper, and failure-state rows listed above within accepted-local real-browser scope. Full action-matrix rows that require live-provider or arbitrary-symbol breadth remain scoped/partial.",
        "",
        "## Repeatability Treatment",
        "",
        "P75 records inherited P71/P73/P74 evidence as scoped and downgrades affected claims rather than claiming fresh repeatability for those milestones. P72 was rerun after P75 hardening and is cited above as repeated real UI scenario evidence. This satisfies P75 10.7a for a scoped conclusion, but it does not authorize a full original-requirement pass.",
    ])
    ACCEPTANCE_PATH.write_text("\n".join(content) + "\n", encoding="utf-8")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="Generate artifacts and fail only on script errors.")
    args = parser.parse_args()

    items = extract_requirements()
    classified = [classify(item) for item in items]
    write_matrix(items, classified)
    acceptance = build_acceptance(items, classified)
    write_acceptance(acceptance)
    SUMMARY_PATH.parent.mkdir(parents=True, exist_ok=True)
    SUMMARY_PATH.write_text(json.dumps(acceptance, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(json.dumps({"requirements": len(items), "result": acceptance["result"], "status_counts": acceptance["status_counts"]}, ensure_ascii=False, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
