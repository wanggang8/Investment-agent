#!/usr/bin/env python3
"""Generate/read P88 structured public-source preverification evidence.

This artifact is intentionally conservative: parser/readback contracts are not
enough to upgrade real-provider rows. A category is real-pass eligible only
when this registry records a verified runtime provider and non-mock SQLite
readback.
"""

from __future__ import annotations

import json
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p88-remaining-full-release-blockers"
ARTIFACT = ARTIFACT_DIR / "p88-source-preverification.json"


CATEGORIES = [
    {
        "category": "capital_flow",
        "requirement_ids": ["REQ-05-003"],
        "required_fields": ["date", "net_inflow", "net_outflow"],
        "candidate_authority": "public market-data website/API candidate",
        "public_access_shape": "no-login/no-paid/no-Level2/no-high-frequency candidate must return stable dated flow fields",
        "stable_request_or_page_evidence": "not verified by P88 runtime provider",
        "update_frequency": "daily-or-lower only if provider terms and response shape are verified",
        "legal_access_limits": "no login, no paid entitlement, no authorization-only source, no Level2, no high-frequency source",
        "rate_limit_assumption": "manual/low-frequency acceptance only until source terms are verified",
        "failure_behavior": "source_unavailable or partial; do not synthesize values",
        "sqlite_target_path": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.capital_flow",
        "runtime_provider_status": "not_verified",
        "real_pass_eligible": False,
        "blocker": "P88 has parser/readback contract but no verified non-mock runtime provider and SQLite readback for capital-flow fields.",
    },
    {
        "category": "margin_financing",
        "requirement_ids": ["REQ-05-004"],
        "required_fields": ["date", "margin_balance", "balance_change_rate"],
        "candidate_authority": "public exchange or public finance data candidate",
        "public_access_shape": "no-login/no-paid/no-Level2/no-high-frequency candidate must return stable dated margin-financing fields",
        "stable_request_or_page_evidence": "not verified by P88 runtime provider",
        "update_frequency": "daily-or-lower only if provider terms and response shape are verified",
        "legal_access_limits": "no login, no paid entitlement, no authorization-only source, no Level2, no high-frequency source",
        "rate_limit_assumption": "manual/low-frequency acceptance only until source terms are verified",
        "failure_behavior": "source_unavailable or partial; do not synthesize values",
        "sqlite_target_path": "market_snapshots.margin_balance, market_snapshots.margin_balance_change, and market_snapshots.market_metrics_json.metadata.p88_structured_fields.margin_financing",
        "runtime_provider_status": "not_verified",
        "real_pass_eligible": False,
        "blocker": "P88 has parser/readback contract but no verified non-mock runtime provider and SQLite readback for margin-financing fields.",
    },
    {
        "category": "constituent_financial",
        "requirement_ids": ["REQ-05-005"],
        "required_fields": ["revenue", "net_profit", "growth", "disclosure_date"],
        "candidate_authority": "public exchange disclosure or public finance statement candidate",
        "public_access_shape": "no-login/no-paid/no-authorization candidate must return stable disclosed financial statement fields",
        "stable_request_or_page_evidence": "not verified by P88 runtime provider",
        "update_frequency": "quarterly/annual disclosure cadence only if provider terms and response shape are verified",
        "legal_access_limits": "no login, no paid entitlement, no authorization-only source, no Level2, no high-frequency source",
        "rate_limit_assumption": "manual/low-frequency acceptance only until source terms are verified",
        "failure_behavior": "source_unavailable or partial; do not synthesize values",
        "sqlite_target_path": "market_snapshots.market_metrics_json.metadata.p88_structured_fields.constituent_financial",
        "runtime_provider_status": "not_verified",
        "real_pass_eligible": False,
        "blocker": "P88 has parser/readback contract but no verified non-mock runtime provider and SQLite readback for constituent-financial fields.",
    },
]


def build() -> dict:
    return {
        "change": "p88-remaining-full-release-blockers-closure",
        "artifact": str(ARTIFACT.relative_to(ROOT)),
        "policy": "accepted-local, fixture, stub, or manually seeded evidence cannot upgrade structured-data collector rows to real_pass",
        "categories": CATEGORIES,
        "summary": {
            "category_count": len(CATEGORIES),
            "real_pass_eligible_count": sum(1 for item in CATEGORIES if item["real_pass_eligible"]),
            "blocked_count": sum(1 for item in CATEGORIES if not item["real_pass_eligible"]),
        },
    }


def validate(payload: dict) -> None:
    categories = payload.get("categories") or []
    if len(categories) != 3:
        raise SystemExit("status=failed\nreason=category_count")
    for item in categories:
        missing = [field for field in ("category", "required_fields", "sqlite_target_path", "runtime_provider_status", "real_pass_eligible", "blocker") if field not in item]
        if missing:
            raise SystemExit(f"status=failed\nreason=missing:{item.get('category')}:{','.join(missing)}")
        if item.get("real_pass_eligible"):
            raise SystemExit(f"status=failed\nreason=unexpected_real_pass_eligible:{item.get('category')}")


def main() -> None:
    check = "--check" in sys.argv
    if check:
        payload = json.loads(ARTIFACT.read_text(encoding="utf-8"))
    else:
        ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
        payload = build()
        ARTIFACT.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    validate(payload)
    print(f"p88_source_preverification:status=passed:blocked={payload['summary']['blocked_count']}:artifact={ARTIFACT}")


if __name__ == "__main__":
    main()
