#!/usr/bin/env python3
"""Validate the P88 inventory against the P86 remaining blocker matrix."""

from __future__ import annotations

import json
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-matrix.md"
ARTIFACT_DIR = ROOT / "docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers"
ARTIFACT = ARTIFACT_DIR / "p88-inventory.json"

EXPECTED_P88_ROWS = {
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


def read_matrix(path: Path) -> list[dict[str, str]]:
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
    return rows


def main() -> None:
    rows = read_matrix(SOURCE_MATRIX)
    remaining = {
        row["requirement_id"]
        for row in rows
        if row.get("full_release_required") == "True" and row.get("p86_status") != "real_pass"
    }
    missing = sorted(EXPECTED_P88_ROWS - remaining)
    unexpected = sorted(remaining - EXPECTED_P88_ROWS)
    by_section: dict[str, list[dict[str, str]]] = defaultdict(list)
    for row in rows:
        if row["requirement_id"] in EXPECTED_P88_ROWS:
            by_section[row["requirement_id"].split("-")[1]].append(
                {
                    "requirement_id": row["requirement_id"],
                    "source_section": row["source_section"],
                    "source_start_line": row["source_start_line"],
                    "requirement_text": row["requirement_text"],
                    "p86_remaining_gap": row.get("p86_remaining_gap", ""),
                    "p86_next_action": row.get("p86_next_action", ""),
                }
            )

    passed = not missing and not unexpected and len(remaining) == 27
    payload = {
        "status": "passed" if passed else "failed",
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "source_matrix": str(SOURCE_MATRIX.relative_to(ROOT)),
        "remaining_full_release_required_non_real_pass_rows": len(remaining),
        "p88_owned_rows": len(EXPECTED_P88_ROWS),
        "missing_expected_rows": missing,
        "unexpected_remaining_rows": unexpected,
        "rows_by_section": {key: value for key, value in sorted(by_section.items())},
    }
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    ARTIFACT.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    if not passed:
        raise SystemExit(
            "p88_inventory:status=failed:"
            f"remaining={len(remaining)}:missing={missing}:unexpected={unexpected}:artifact={ARTIFACT}"
        )
    print(f"p88_inventory:status=passed:remaining=27:p88=27:artifact={ARTIFACT}")


if __name__ == "__main__":
    main()
