# Tasks: P73 Product Effectiveness And UX Validation

## 1. Change Setup

- [x] 1.1 Confirm no active change is present before opening P73.
- [x] 1.2 Create OpenSpec change `p73-product-effectiveness-ux-validation`.
- [x] 1.3 Update governance/progress docs to mark P73 active.
- [x] 1.4 Run `openspec validate p73-product-effectiveness-ux-validation --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Effectiveness And UX Acceptance Design

- [x] 2.1 Define product-effectiveness metrics tied to discipline adherence, evidence sufficiency, traceability, review usefulness, and UX comprehension.
- [x] 2.2 Define a real UX task matrix covering first-use, daily discipline, portfolio, data quality/evidence, consultation, confirmation, risk/notification/audit/rules/review, and unsafe input.
- [x] 2.3 Define effect replay scenarios that validate discipline behavior rather than investment returns.
- [x] 2.4 Define pass/block/gap wording so safe degradation cannot be counted as product-effectiveness pass.

## 3. Browser UX Task Runner

- [x] 3.1 Add a P73 browser spec that operates representative real UX tasks against the local app.
- [x] 3.2 Capture screenshots and sanitized task results under `docs/release/ui-audit-assets/2026-06-19-p73/`.
- [x] 3.3 Treat page errors, unexpected API failures, console errors, forbidden affordances, and critical UX confusion as blockers.

## 4. Effect Replay Checks

- [x] 4.1 Add deterministic replay/check script for background-only evidence, insufficient evidence, rule-effect gate state, risk/readback links, and manual-confirmation-only portfolio mutation.
- [x] 4.2 Write sanitized JSON summary for replay results.
- [x] 4.3 Verify replay checks do not use future outcomes or investment return claims as pass criteria.

## 5. UX Audit Record

- [x] 5.1 Review screenshots and task results for information hierarchy, next-action clarity, state labels, navigation, mobile/reflow, and copy safety.
- [x] 5.2 Record critical/major/minor UX findings with evidence and disposition.
- [x] 5.3 If critical UX findings exist, mark P73 blocked until fixed and rerun affected tasks.

## 6. Release Materials

- [x] 6.1 Add `docs/release/acceptance/2026-06-19-p73-product-effectiveness-ux-validation.md`.
- [x] 6.2 Update release README, release candidate/handoff, repeatability, docs README, development plan, governance, AGENTS, OpenSpec project, and progress materials.
- [x] 6.3 Record gaps explicitly: no future return guarantee, no full real-world user study unless separately performed, no physical second-machine repeat, no post-P72/P73 package refresh unless separately performed.

## 7. Verification And Archive

- [x] 7.1 Run targeted tests for new checks.
- [x] 7.2 Run `npm --prefix web test`.
- [x] 7.3 Run `npm --prefix web run build`.
- [x] 7.4 Run `go test ./...`.
- [x] 7.5 Run the P73 browser/effectiveness runner.
- [x] 7.6 Run safety/redaction scans and classify expected forbidden-boundary strings.
- [x] 7.7 Run `openspec validate p73-product-effectiveness-ux-validation --strict`, `openspec validate --all --strict`, and `git diff --check`.
- [x] 7.8 Archive P73 only after product-effectiveness evidence, UX audit, and gap review are recorded.
