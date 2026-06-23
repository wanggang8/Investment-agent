#!/usr/bin/env python3
"""Verify that P86 owns the full P87-after remainder."""

from __future__ import annotations

import json
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path


MATRIX = Path("docs/release/acceptance/2026-06-22-p87-portfolio-state-allocation-safety-matrix.md")
ARTIFACT = Path("docs/release/ui-audit-assets/2026-06-22-p86-core-goal-knowledge-safety-final/p86-inventory.json")

P86 = set(
    """
    REQ-01-001 REQ-01-002 REQ-01-003 REQ-01-004 REQ-01-005 REQ-01-006
    REQ-01-011 REQ-02-002 REQ-02-006 REQ-02-007 REQ-02-008 REQ-02-010
    REQ-02-011 REQ-02-012 REQ-02-013 REQ-02-017 REQ-02-020 REQ-02-021
    REQ-02-022 REQ-02-024 REQ-02-025 REQ-02-026 REQ-02-029 REQ-02-032
    REQ-03-001 REQ-03-002 REQ-03-003 REQ-03-004 REQ-03-005 REQ-03-006
    REQ-03-007 REQ-03-008 REQ-03-009 REQ-03-010 REQ-04-003 REQ-04-008
    REQ-04-016 REQ-04-025 REQ-05-003 REQ-05-004 REQ-05-005 REQ-05-010
    REQ-06-002 REQ-06-003 REQ-06-004 REQ-06-005 REQ-06-006 REQ-06-007
    REQ-06-008 REQ-06-010 REQ-06-011 REQ-06-012 REQ-06-013 REQ-06-014
    REQ-06-015 REQ-06-016 REQ-06-017 REQ-06-018 REQ-06-019 REQ-06-020
    REQ-06-021 REQ-06-022 REQ-06-023 REQ-06-024 REQ-06-025 REQ-07-001
    REQ-07-002 REQ-07-003 REQ-07-004 REQ-07-005 REQ-07-006 REQ-07-007
    REQ-07-008 REQ-07-009 REQ-07-010 REQ-07-012 REQ-07-015 REQ-08-004
    REQ-08-018 REQ-08-020 REQ-08-023 REQ-09-001 REQ-09-003 REQ-09-004
    REQ-09-006 REQ-09-007 REQ-09-008 REQ-09-009 REQ-09-010 REQ-09-013
    REQ-09-023 REQ-09-024 REQ-09-025 REQ-09-027 REQ-10-004 REQ-13-010
    REQ-14-004 REQ-14-005 REQ-14-007 REQ-15-001 REQ-15-002 REQ-15-003
    REQ-15-004 REQ-15-005 REQ-15-007 REQ-16-001 REQ-16-002 REQ-16-005
    REQ-16-006 REQ-16-007 REQ-16-008 REQ-16-009 REQ-16-010 REQ-16-011
    REQ-16-013 REQ-16-014 REQ-16-015 REQ-16-019 REQ-16-021 REQ-16-023
    REQ-16-025 REQ-16-028 REQ-16-029 REQ-16-030 REQ-16-031 REQ-16-033
    REQ-16-034 REQ-17-003 REQ-17-004 REQ-17-005 REQ-17-007 REQ-17-008
    REQ-17-011 REQ-17-012 REQ-17-014 REQ-17-015 REQ-17-024
    """.split()
)


def parse_markdown_table(path: Path) -> list[dict[str, str]]:
    lines = [line for line in path.read_text().splitlines() if line.startswith("|")]
    header = [cell.strip() for cell in lines[0].strip("|").split("|")]
    rows: list[dict[str, str]] = []
    for line in lines[2:]:
        cells = [cell.strip() for cell in line.strip("|").split("|")]
        if len(cells) > len(header):
            extra = len(cells) - len(header)
            merged_text = "|".join(cells[5 : 6 + extra])
            cells = cells[:5] + [merged_text] + cells[6 + extra :]
        if len(cells) != len(header):
            raise SystemExit(f"cannot parse row with {len(cells)} cells: {line[:160]}")
        rows.append(dict(zip(header, cells)))
    return rows


def main() -> int:
    rows = parse_markdown_table(MATRIX)
    status_col = "p87_status"
    remaining = {
        row["requirement_id"]
        for row in rows
        if row["full_release_required"] == "True" and row[status_col] != "real_pass"
    }
    missing = remaining - P86
    extra = P86 - remaining

    if missing or extra or len(remaining) != 137 or len(P86) != 137:
        print("p86_inventory:status=failed")
        print(f"remaining={len(remaining)} p86={len(P86)}")
        print(f"missing={','.join(sorted(missing)) or '-'}")
        print(f"extra={','.join(sorted(extra)) or '-'}")
        return 1

    by_section: dict[str, list[str]] = defaultdict(list)
    for requirement_id in sorted(P86):
        by_section[requirement_id.split("-")[1]].append(requirement_id)
    ARTIFACT.parent.mkdir(parents=True, exist_ok=True)
    ARTIFACT.write_text(
        json.dumps(
            {
                "generated_at": datetime.now(timezone.utc).isoformat(),
                "source_matrix": str(MATRIX),
                "status": "passed",
                "remaining_full_release_required_non_real_pass_rows": len(remaining),
                "p86_owned_rows": len(P86),
                "rows_by_section": dict(sorted(by_section.items())),
            },
            ensure_ascii=False,
            indent=2,
        )
        + "\n"
    )
    print(f"p86_inventory:status=passed:remaining=137:p86=137:artifact={ARTIFACT}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
