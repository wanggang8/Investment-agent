# P93 Design

## Approach

P93 is a code-review and implementation-reality audit. It is intentionally stricter than a summary-only release note, but it does not replace P92's row-level ledger. P92 remains the 341-row original-requirement artifact; P93 cross-checks that every P92 row has required audit fields and resolves through its source section to current production code/evidence bundles. It also inspects route wiring, configuration defaults, deployment scripts, test gates, secret literals, dead code, and suspicious-token contexts.

The audit script produces a deterministic Markdown report. It combines static code inventory, P92 row-level ledger checks, known requirement-section-to-code mappings, route checks, config and secret checks, and suspicious token classification. The report is not a substitute for tests; it is a release audit trail that points to concrete implementation files and fails if active release-blocking findings remain.

## Outputs

- `scripts/p93_code_reality_audit.py`
- `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`

## Validation

`scripts/p93_code_reality_audit.py --check` regenerates the report in memory and fails if the checked-in report is stale, if the P92 row-level ledger counts/required fields drift, if any scanned non-test source/config file contains an unredacted `sk-...` key literal, or if another release-blocking pattern is found.
