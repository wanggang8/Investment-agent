# P92 Tasks

## 1. Governance

- [x] Create P92 OpenSpec change and validate it.
- [x] Keep scope limited to final audit ledger generation and governance documentation.

## 2. Generator And Checker

- [x] Add `scripts/p92_final_requirement_audit.py`.
- [x] Parse the P88 full matrix and overlay P89/P90 final evidence.
- [x] Classify rows by original requirement section, feature area, UI/product surface, expected data impact, readback/audit evidence, and safety boundary.
- [x] Implement `--check` mode for completeness and staleness validation.

## 3. Generated Audit Artifacts

- [x] Generate final original-requirement audit ledger.
- [x] Generate final audit summary.
- [x] Confirm zero full-release-required rows remain non-`real_pass`.
- [x] Keep reference-only rows outside product pass claims.

## 4. Governance Updates

- [x] Update `AGENTS.md`, `docs/GOVERNANCE.md`, `openspec/PROGRESS.md`, and `openspec/project.md` with P92 status.
- [x] Archive P92 after validation passes.

## 5. Validation

- [x] Run `python3 scripts/p92_final_requirement_audit.py --check`.
- [x] Run `openspec validate --all --strict`.
- [x] Run `git diff --check`.
