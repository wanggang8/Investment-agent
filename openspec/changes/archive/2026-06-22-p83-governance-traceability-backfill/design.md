# P83 Design

## Evidence Strategy

P83 is primarily a traceability and evidence-quality closure. It should avoid changing runtime behavior unless a row has a real product defect. The evidence layer must point to exact local artifacts and fresh command results, not narrative-only assertions.

Evidence categories:

- OpenSpec/archive/governance artifacts.
- Release acceptance records and package manifests.
- CLI/test command outputs.
- Safety scans and forbidden capability checks.
- UI/API evidence where the row concerns user-visible governance or ops behavior.

## Classification Rule

Rows that describe historical governance facts may become `reference_only` or scoped evidence if they are not full-release-required product behavior. Rows that are full-release-required can only become `real_pass` with concrete current evidence.

