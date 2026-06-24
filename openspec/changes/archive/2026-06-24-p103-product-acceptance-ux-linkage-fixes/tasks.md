# Tasks

## 1. Governance

- [x] Validate P103:
  - [x] `openspec validate p103-product-acceptance-ux-linkage-fixes --strict`

## 2. UX Fixes

- [x] Map portfolio `NOT_FOUND` empty state to first-use onboarding copy.
- [x] Collapse analyst analysis material by default in decision detail.
- [x] Support `decision_id` query focus/filter in decision loop.

## 3. Tests And Product Linkage Verification

- [x] Add/adjust focused frontend tests for the three fixes.
- [x] Run linked product workflow verification and record acceptance notes.
- [x] Document whether this is exhaustive feature validation or scoped linked-flow validation.

## 4. Gates

- [x] `npm --prefix web test -- --run`
- [x] `npm --prefix web run build`
- [x] `go test ./...`
- [x] `openspec validate --all --strict`
- [x] `python3 scripts/p92_final_requirement_audit.py --check`
- [x] `python3 scripts/p93_code_reality_audit.py --check`
- [x] `git diff --check`

## 5. Archive

- [x] Update governance/progress.
- [x] Archive P103.
