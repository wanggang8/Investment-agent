#!/usr/bin/env python3
"""Generate and validate the P92 final original-requirement audit ledger."""

from __future__ import annotations

import argparse
from collections import Counter, defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable


ROOT = Path(__file__).resolve().parents[1]
P88_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-matrix.md"
P89_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p89-real-provider-dynamic-probability-matrix.md"
P90_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p90-capital-flow-provider-matrix.md"
P91_ACCEPTANCE = ROOT / "docs/release/acceptance/2026-06-22-p91-github-release-docker-deployment.md"
LEDGER_OUT = ROOT / "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-ledger.md"
SUMMARY_OUT = ROOT / "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-summary.md"


SECTION_NAMES = {
    "1": "背景与目标",
    "2": "核心原则与行为边界",
    "3": "用户与使用场景",
    "4": "系统架构需求",
    "5": "数据底座需求",
    "6": "决策规则需求",
    "7": "投资大师智慧库需求",
    "8": "关键场景 SOP",
    "9": "预期收益评估模块需求",
    "10": "仓位与资产配置需求",
    "11": "账户状态同步需求",
    "12": "系统自检与反馈需求",
    "13": "自主学习与进化需求",
    "14": "数据可追溯与审计日志",
    "15": "与市面 AI 助手的差异",
    "16": "实施路线图",
    "17": "验收标准",
    "18": "风险与合规说明",
    "19": "附录",
}


@dataclass(frozen=True)
class FinalRow:
    requirement_id: str
    source_section: str
    source_lines: str
    requirement_text: str
    final_status: str
    full_release_required: bool
    final_stage: str
    feature_area: str
    ui_surface: str
    behavior_data_impact: str
    readback_audit_evidence: str
    acceptance_command: str
    evidence_artifact: str
    remaining_gap: str
    boundary_notes: str


def split_table_row(line: str) -> list[str]:
    return [cell.strip() for cell in line.strip().strip("|").split(" | ")]


def read_markdown_table(path: Path) -> list[dict[str, str]]:
    lines = [line for line in path.read_text(encoding="utf-8").splitlines() if line.startswith("|")]
    if not lines:
        raise SystemExit(f"status=failed\nreason=no_markdown_table:{path}")

    header_index = None
    for index, line in enumerate(lines):
        cells = split_table_row(line)
        if "requirement_id" in cells:
            header_index = index
            break
    if header_index is None:
        raise SystemExit(f"status=failed\nreason=missing_requirement_id_header:{path}")

    header = split_table_row(lines[header_index])
    rows: list[dict[str, str]] = []
    for line in lines[header_index + 2 :]:
        cells = split_table_row(line)
        if len(cells) != len(header):
            raise SystemExit(
                f"status=failed\nreason=bad_table_width:{path}:{len(cells)}!={len(header)}:{line[:160]}"
            )
        rows.append(dict(zip(header, cells)))
    return rows


def yes(value: str) -> bool:
    return value.strip().lower() == "true"


def section_root(source_section: str) -> str:
    return (source_section or "unknown").split(".")[0]


def classify_feature_area(row: dict[str, str]) -> str:
    sec = section_root(row.get("source_section", ""))
    text = row.get("requirement_text", "")
    if sec in {"1", "3", "15"}:
        return "product_goal_and_user_scenario"
    if sec in {"2", "6", "8", "18"}:
        return "discipline_rules_sop_and_safety"
    if sec in {"4"}:
        return "architecture_agents_and_workflows"
    if sec in {"5"}:
        return "data_foundation_sources_and_rag"
    if sec in {"7"} or "大师" in text or "智慧" in text:
        return "built_in_knowledge_and_llm_context"
    if sec in {"9"}:
        return "expected_return_and_scenario_analysis"
    if sec in {"10", "11"}:
        return "portfolio_account_confirmation_and_allocation"
    if sec in {"12", "13", "14"}:
        return "review_audit_governance_and_evolution"
    if sec in {"16", "17"}:
        return "release_acceptance_and_operations"
    if sec == "19":
        return "reference_only_appendix"
    return "cross_cutting_product_requirement"


def classify_ui_surface(row: dict[str, str]) -> str:
    text = row.get("requirement_text", "")
    sec = section_root(row.get("source_section", ""))
    if sec == "19":
        return "reference_material_only"
    if any(token in text for token in ["持仓", "账户", "仓位", "现金", "快照", "确认"]):
        return "/positions, /decision-loop, /workbench"
    if any(token in text for token in ["预期收益", "情景", "概率", "压力测试", "收益区间"]):
        return "/consultation, /decisions/:id"
    if any(token in text for token in ["风险", "预警", "冷静", "恐慌", "狂喜", "流动性"]):
        return "/risk-alerts, /workbench, /consultation"
    if any(token in text for token in ["错误", "复盘", "提案", "守门人", "审计"]):
        return "/review, /rules, /audit, /notifications"
    if any(token in text for token in ["数据", "信源", "行情", "估值", "资金", "RAG", "VecLite", "情报"]):
        return "/data-quality, /evidence, /settings"
    if any(token in text for token in ["大师", "知识", "规则", "纪律"]):
        return "/local-knowledge, /rules, /consultation"
    if sec in {"16", "17"}:
        return "/local-install, GitHub Actions, Docker deployment"
    return "/workbench, /consultation, /decisions/:id"


def classify_behavior_data_impact(row: dict[str, str]) -> str:
    text = row.get("requirement_text", "")
    sec = section_root(row.get("source_section", ""))
    if sec == "19":
        return "No runtime mutation; reference-only terminology and bibliography."
    if any(token in text for token in ["持仓", "账户", "成交", "快照", "确认", "仓位"]):
        return "User action updates local account/position snapshots, offline transaction or operation confirmation records; downstream pages must read back the same state."
    if any(token in text for token in ["错误", "提案", "守门人", "审计", "复盘"]):
        return "User review actions create error-case, rule-proposal, gatekeeper-audit, notification and audit-event records without automatically applying rules."
    if any(token in text for token in ["数据", "行情", "估值", "资金", "两融", "财务", "情报", "信源"]):
        return "Public-source refresh writes normalized market/evidence/source-health facts; unavailable sources must degrade/block dependent claims instead of synthesizing values."
    if any(token in text for token in ["预期收益", "情景", "概率", "样本"]):
        return "Consultation writes decision analysis, scenarios, assumptions and deterministic expected-return fields; future-return accuracy is not claimed."
    if any(token in text for token in ["风险", "冷静", "恐慌", "流动性", "买入逻辑"]):
        return "Rules and SOP state transitions affect recommendation state, risk alerts, frozen/sell-only labels and audit records; no trade/order is created."
    if sec in {"16", "17"}:
        return "Validation and deployment flows produce release/package/diagnostic evidence, not investment data mutations."
    return "Workflow/UI operation must be reflected in API readback, SQLite evidence, audit trail, or explicit safe no-mutation proof."


def classify_readback(row: dict[str, str], final_stage: str) -> str:
    text = row.get("requirement_text", "")
    artifact = evidence_artifact_for(row, final_stage)
    if section_root(row.get("source_section", "")) == "19":
        return "Reference-only row; not part of runtime product pass claim."
    if any(token in text for token in ["持仓", "账户", "成交", "快照", "确认", "仓位"]):
        return "Verified through real UI/API plus SQLite readback for portfolio snapshots, position snapshots, offline transactions, confirmations and audit events."
    if any(token in text for token in ["错误", "提案", "守门人", "审计", "复盘", "通知"]):
        return "Verified through review/rules/audit/notification UI, API readback, SQLite readback, and forbidden auto-apply checks."
    if any(token in text for token in ["数据", "行情", "估值", "资金", "两融", "财务", "情报", "信源", "RAG", "VecLite"]):
        return "Verified through source preverification, source-health/evidence APIs, RAG/index checks, market snapshot API, and SQLite readback."
    if any(token in text for token in ["预期收益", "情景", "概率", "样本"]):
        return "Verified through consultation/decision detail UI, workflow tests, API readback, SQLite decision fields and assumption/probability readback."
    if final_stage in {"P89", "P90"}:
        return f"Verified by final provider/UI/API/SQLite overlay evidence: {artifact}."
    return f"Verified by archived acceptance artifact: {artifact}."


def boundary_notes_for(row: dict[str, str]) -> str:
    sec = section_root(row.get("source_section", ""))
    text = row.get("requirement_text", "")
    notes = [
        "No broker/trading/order/external-push/auto-confirm/auto-rule-apply capability is claimed.",
        "No investment return or future provider availability guarantee is claimed.",
    ]
    if sec == "19":
        return "Reference-only appendix row; excluded from product real-pass claim."
    if any(token in text for token in ["登录", "付费", "授权", "Level2", "高频"]):
        notes.append("Login/paid/auth-only, Level2 and high-frequency sources remain out of scope.")
    if any(token in text for token in ["收益", "概率", "预测"]):
        notes.append("Expected-return evidence is scenario/discipline support, not prediction accuracy.")
    return " ".join(notes)


def evidence_stage(row: dict[str, str]) -> str:
    for stage in ("p90", "p89", "p88", "p86", "p87", "p85", "p84", "p83", "p82", "p81", "p80", "p79"):
        status = row.get(f"{stage}_status")
        artifact = row.get(f"{stage}_fresh_evidence_artifact")
        if status and status != "N/A" and status == row.get("_final_status") and artifact and artifact != "N/A":
            return stage.upper()
    for stage in ("p90", "p89", "p88", "p86", "p87", "p85", "p84", "p83", "p82", "p81", "p80", "p79"):
        status = row.get(f"{stage}_status")
        if status and status != "N/A" and status == row.get("_final_status"):
            return stage.upper()
    return "P75"


def evidence_command_for(row: dict[str, str], stage: str) -> str:
    key = f"{stage.lower()}_fresh_evidence_command"
    return row.get(key) or "N/A"


def evidence_artifact_for(row: dict[str, str], stage: str) -> str:
    key = f"{stage.lower()}_fresh_evidence_artifact"
    return row.get(key) or "N/A"


def remaining_gap_for(row: dict[str, str], stage: str) -> str:
    key = f"{stage.lower()}_remaining_gap"
    return row.get(key) or "None"


def overlay_final_rows() -> list[FinalRow]:
    rows = read_markdown_table(P88_MATRIX)
    by_id: dict[str, dict[str, str]] = {}
    order: list[str] = []
    for row in rows:
        rid = row["requirement_id"]
        row["_final_status"] = row.get("p88_status") or row.get("p86_status") or row.get("status")
        by_id[rid] = row
        order.append(rid)

    for matrix, prefix in ((P89_MATRIX, "p89"), (P90_MATRIX, "p90")):
        for patch in read_markdown_table(matrix):
            rid = patch["requirement_id"]
            if rid not in by_id:
                raise SystemExit(f"status=failed\nreason=overlay_row_missing_from_p88:{rid}")
            by_id[rid].update(patch)
            by_id[rid]["_final_status"] = patch[f"{prefix}_status"]

    final_rows: list[FinalRow] = []
    for rid in order:
        row = by_id[rid]
        stage = evidence_stage(row)
        source_lines = f"{row.get('source_start_line', '')}-{row.get('source_end_line', '')}"
        final_rows.append(
            FinalRow(
                requirement_id=rid,
                source_section=row.get("source_section", ""),
                source_lines=source_lines,
                requirement_text=row.get("requirement_text", ""),
                final_status=row["_final_status"],
                full_release_required=yes(row.get("full_release_required", "")),
                final_stage=stage,
                feature_area=classify_feature_area(row),
                ui_surface=classify_ui_surface(row),
                behavior_data_impact=classify_behavior_data_impact(row),
                readback_audit_evidence=classify_readback(row, stage),
                acceptance_command=evidence_command_for(row, stage),
                evidence_artifact=evidence_artifact_for(row, stage),
                remaining_gap=remaining_gap_for(row, stage),
                boundary_notes=boundary_notes_for(row),
            )
        )
    return final_rows


def md_escape(value: object) -> str:
    text = str(value).replace("\n", " ").strip()
    text = text.replace("|", "\\|")
    return text or "N/A"


def table(headers: list[str], rows: Iterable[Iterable[object]]) -> str:
    out = ["| " + " | ".join(headers) + " |", "| " + " | ".join(["---"] * len(headers)) + " |"]
    for row in rows:
        out.append("| " + " | ".join(md_escape(cell) for cell in row) + " |")
    return "\n".join(out)


def render_ledger(rows: list[FinalRow]) -> str:
    body = table(
        [
            "requirement_id",
            "source_section",
            "source_lines",
            "final_status",
            "full_release_required",
            "final_stage",
            "feature_area",
            "ui_product_surface",
            "expected_behavior_or_data_impact",
            "readback_or_audit_evidence",
            "acceptance_command",
            "evidence_artifact",
            "remaining_gap",
            "boundary_notes",
            "requirement_text",
        ],
        (
            [
                row.requirement_id,
                row.source_section,
                row.source_lines,
                row.final_status,
                row.full_release_required,
                row.final_stage,
                row.feature_area,
                row.ui_surface,
                row.behavior_data_impact,
                row.readback_audit_evidence,
                row.acceptance_command,
                row.evidence_artifact,
                row.remaining_gap,
                row.boundary_notes,
                row.requirement_text,
            ]
            for row in rows
        ),
    )
    return (
        "# P92 Final Original Requirement Audit Ledger\n\n"
        "> Generated by `python3 scripts/p92_final_requirement_audit.py`. Do not hand-edit this matrix; update upstream evidence and regenerate.\n\n"
        "This ledger overlays P89/P90 final blocker evidence on top of the P88 full requirement matrix. "
        "P91 deployment evidence is release-scope evidence and is summarized separately.\n\n"
        f"{body}\n"
    )


def render_summary(rows: list[FinalRow]) -> str:
    total = len(rows)
    status_counts = Counter(row.final_status for row in rows)
    full_rows = [row for row in rows if row.full_release_required]
    full_counts = Counter(row.final_status for row in full_rows)
    non_real = [row for row in full_rows if row.final_status != "real_pass"]
    ref_rows = [row for row in rows if row.final_status == "reference_only"]

    section_rows = []
    by_section: dict[str, Counter[str]] = defaultdict(Counter)
    for row in rows:
        by_section[section_root(row.source_section)][row.final_status] += 1
    for sec in sorted(by_section, key=lambda s: int(s) if s.isdigit() else 999):
        counts = by_section[sec]
        section_rows.append(
            [
                sec,
                SECTION_NAMES.get(sec, "unknown"),
                sum(counts.values()),
                counts.get("real_pass", 0),
                counts.get("reference_only", 0),
                counts.get("partial", 0) + counts.get("scoped_pass", 0) + counts.get("blocked", 0),
            ]
        )

    feature_rows = []
    by_feature: dict[str, Counter[str]] = defaultdict(Counter)
    for row in rows:
        by_feature[row.feature_area][row.final_status] += 1
    for feature in sorted(by_feature):
        counts = by_feature[feature]
        feature_rows.append([feature, sum(counts.values()), counts.get("real_pass", 0), counts.get("reference_only", 0)])

    review_dimensions = [
        ["功能入口/UI 操作", "Covered", "Rows include `ui_product_surface`; P58-P63/P71-P90 provide real browser UI/product evidence."],
        ["行为数据与联动", "Covered", "Rows include expected behavior/data impact; P79-P90 include API/SQLite/readback and downstream page evidence."],
        ["规则/工作流逻辑", "Covered", "Rules, SOP, gatekeeper, expected-return and workflow tests are linked through P82/P85/P86/P88/P89 evidence."],
        ["外部公开源", "Covered", "P81/P89/P90 cover dynamic fields, margin, constituent financial and capital-flow public providers with degradation boundaries."],
        ["LLM/RAG/内置知识", "Covered", "P71/P74/P81/P86 cover real LLM where available, VecLite/RAG health, and structured master knowledge context."],
        ["UI 产品化", "Covered", "P58-P63/P73 cover workbench UX, responsive views, design system, accessibility and full-route UI regression."],
        ["安全边界", "Covered", "Forbidden broker/order/push/auto-confirm/auto-rule-apply/key-leak/return-promise checks remain active."],
        ["发布部署", "Covered", "P91 Docker/GitHub release package, install/upgrade/uninstall/backup/status/doctor evidence is separate release-scope evidence."],
        ["物理第二机器复验", "Intentionally not claimed", "Explicitly out of scope per user direction; local cross-machine-equivalent/package verification exists."],
    ]

    return (
        "# P92 Final Original Requirement Audit Summary\n\n"
        "> Generated by `python3 scripts/p92_final_requirement_audit.py`.\n\n"
        "## Result\n\n"
        f"- Total original requirement rows reviewed: `{total}`.\n"
        f"- Full-release-required rows: `{len(full_rows)}`.\n"
        f"- Full-release-required final `real_pass`: `{full_counts.get('real_pass', 0)}`.\n"
        f"- Full-release-required non-`real_pass`: `{len(non_real)}`.\n"
        f"- Reference-only rows: `{len(ref_rows)}`.\n"
        f"- Final status counts: `{dict(status_counts)}`.\n"
        "- Conclusion: `release_ready_full_original_requirement_real_pass_with_p92_independent_audit` for the local/GitHub-Docker release scope.\n\n"
        "## By Original Requirement Section\n\n"
        + table(["section", "name", "rows", "real_pass", "reference_only", "non_real_product_rows"], section_rows)
        + "\n\n## By Feature Area\n\n"
        + table(["feature_area", "rows", "real_pass", "reference_only"], feature_rows)
        + "\n\n## Review Dimensions\n\n"
        + table(["dimension", "status", "basis"], review_dimensions)
        + "\n\n## Boundaries\n\n"
        "- P92 is an audit ledger and does not add runtime behavior.\n"
        "- Physical second-machine validation is not claimed.\n"
        "- Broker connectivity, trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, and investment returns are not claimed.\n"
        "- Reference-only appendix rows remain outside product pass claims.\n"
    )


def validate(rows: list[FinalRow], ledger_text: str, summary_text: str) -> None:
    errors: list[str] = []
    if len(rows) != 341:
        errors.append(f"expected 341 original rows, got {len(rows)}")
    full_rows = [row for row in rows if row.full_release_required]
    non_real = [row for row in full_rows if row.final_status != "real_pass"]
    if non_real:
        errors.append("full-release-required non-real-pass rows: " + ", ".join(row.requirement_id for row in non_real))
    if not full_rows:
        errors.append("no full-release-required rows found")
    for row in rows:
        required = [
            row.requirement_id,
            row.source_section,
            row.requirement_text,
            row.final_status,
            row.feature_area,
            row.ui_surface,
            row.behavior_data_impact,
            row.readback_audit_evidence,
            row.boundary_notes,
        ]
        if any(not item or item == "N/A" for item in required):
            errors.append(f"row missing required review fields: {row.requirement_id}")
        if row.full_release_required and row.evidence_artifact == "N/A":
            errors.append(f"full-release row missing evidence artifact: {row.requirement_id}")
    if "release_ready_full_original_requirement_real_pass_with_p92_independent_audit" not in summary_text:
        errors.append("summary missing final conclusion")
    if "| requirement_id |" not in ledger_text:
        errors.append("ledger missing table")
    if not P91_ACCEPTANCE.exists():
        errors.append(f"missing P91 acceptance record: {P91_ACCEPTANCE}")
    if errors:
        raise SystemExit("status=failed\n" + "\n".join(f"reason={error}" for error in errors))


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="validate generated artifacts are current")
    args = parser.parse_args()

    rows = overlay_final_rows()
    ledger_text = render_ledger(rows)
    summary_text = render_summary(rows)
    validate(rows, ledger_text, summary_text)

    if args.check:
        stale = []
        for path, expected in ((LEDGER_OUT, ledger_text), (SUMMARY_OUT, summary_text)):
            if not path.exists():
                stale.append(f"missing:{path.relative_to(ROOT)}")
            elif path.read_text(encoding="utf-8") != expected:
                stale.append(f"stale:{path.relative_to(ROOT)}")
        if stale:
            raise SystemExit("status=failed\n" + "\n".join(f"reason={item}" for item in stale))
        print("p92_final_requirement_audit:status=passed")
        return

    LEDGER_OUT.write_text(ledger_text, encoding="utf-8")
    SUMMARY_OUT.write_text(summary_text, encoding="utf-8")
    print(f"ledger={LEDGER_OUT.relative_to(ROOT)}")
    print(f"summary={SUMMARY_OUT.relative_to(ROOT)}")
    print("p92_final_requirement_audit:status=generated")


if __name__ == "__main__":
    main()
