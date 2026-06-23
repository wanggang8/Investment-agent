# P93 Tasks

## 1. Governance

- [x] Create P93 OpenSpec change and validate it.
- [x] Keep scope limited to code/design audit unless a blocker is found.

## 2. Code Reality Audit

- [x] Add `scripts/p93_code_reality_audit.py`.
- [x] Map original requirement sections to production implementation files.
- [x] Cross-check the P92 341-row requirement ledger and resolve every row to a P93 code/evidence bundle.
- [x] Scan suspicious demo/mock/stub/hardcoded/placeholder/dead-code tokens by context.
- [x] Scan current non-test source/config files for `sk-...` API key literals and fail on unredacted secrets.
- [x] Verify production route wiring does not use placeholder/demo pages.
- [x] Verify deployment config uses real-data defaults and embeds no secrets.

## 3. Audit Report

- [x] Generate `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- [x] Record findings with severity and release impact.
- [x] Record P92 row-level ledger cross-check counts and clarify that P92 remains the 341-row artifact.
- [x] Record requirement-section implementation evidence and design reasonableness.

## 4. Validation

- [x] Run `python3 scripts/p93_code_reality_audit.py --check`.
- [x] Run `python3 scripts/p92_final_requirement_audit.py --check`.
- [x] Run `openspec validate --all --strict`.
- [x] Run `go test ./...`.
- [x] Run `npm --prefix web test`.
- [x] Run `npm --prefix web run build`.
- [x] Run `git diff --check`.
- [x] Archive P93 after validation passes.
