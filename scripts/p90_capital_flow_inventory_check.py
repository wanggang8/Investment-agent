#!/usr/bin/env python3
"""Validate the P90 inventory against the P89 closure matrix."""

from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SOURCE_MATRIX = ROOT / "docs" / "release" / "acceptance" / "2026-06-22-p89-real-provider-dynamic-probability-matrix.md"
ARTIFACT_DIR = ROOT / "docs" / "release" / "ui-audit-assets" / "2026-06-22-p90-capital-flow-provider"
ARTIFACT = ARTIFACT_DIR / "p90-inventory.json"
EXPECTED_ROWS = {"REQ-04-016", "REQ-05-003"}


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
        if row.get("p89_status") != "real_pass"
    }
    missing = sorted(EXPECTED_ROWS - remaining)
    unexpected = sorted(remaining - EXPECTED_ROWS)
    owned_rows = [
        {
            "requirement_id": row["requirement_id"],
            "p89_status": row.get("p89_status", ""),
            "p89_remaining_gap": row.get("p89_remaining_gap", ""),
            "p89_next_action": row.get("p89_next_action", ""),
        }
        for row in rows
        if row.get("requirement_id") in EXPECTED_ROWS
    ]
    passed = not missing and not unexpected and len(remaining) == 2
    payload = {
        "status": "passed" if passed else "failed",
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "source_matrix": str(SOURCE_MATRIX.relative_to(ROOT)),
        "remaining_non_real_pass_rows": len(remaining),
        "p90_owned_rows": sorted(EXPECTED_ROWS),
        "missing_expected_rows": missing,
        "unexpected_remaining_rows": unexpected,
        "rows": owned_rows,
    }
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    ARTIFACT.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    if not passed:
        raise SystemExit(f"p90_inventory:status=failed:remaining={len(remaining)}:missing={missing}:unexpected={unexpected}:artifact={ARTIFACT}")
    print(f"p90_inventory:status=passed:remaining=2:p90=2:artifact={ARTIFACT}")


if __name__ == "__main__":
    main()
