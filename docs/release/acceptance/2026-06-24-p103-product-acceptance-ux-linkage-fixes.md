# P103 Product Acceptance UX Linkage Fixes

Date: 2026-06-24

Scope: fix P102 non-blocking product UX findings and add explicit linked-flow validation notes. This change used the local source runtime and did not validate Docker, installation, GitHub Release, distribution packages, or physical second-machine execution.

## Verdict

`p102_ux_findings_fixed_and_linked_flow_reverified`

P103 fixes the three P102 product UX findings:

1. Empty/missing portfolio snapshots no longer show a generic read failure on `/positions`; the page falls through to first-use onboarding and local calibration guidance.
2. Decision detail now keeps analyst analysis material collapsed by default, while still leaving full LLM analysis available through explicit expansion.
3. `/decision-loop?decision_id=<id>` now calls the single decision-loop read path and focuses that decision instead of showing the full recent list.

## What Was Fixed

### Portfolio Empty State

Root cause: the frontend treated `NOT_FOUND` from `GET /api/v1/portfolio/current` as a generic failure.

Fix: `PortfolioPage` now treats portfolio `NOT_FOUND` as a first-use onboarding state. Other API errors still use the safe `StatusNotice` path.

Validation:

- `PortfolioPage.test.tsx` covers missing snapshot as first-use onboarding and asserts no `读取失败`/generic failure text.
- Rendered `/positions` check on the current local DB confirmed no generic failure and a usable portfolio state.

### Decision Detail Density

Root cause: real LLM analyst text was expanded by default, and fallback story reasons could include full analyst conclusions.

Fix:

- `DecisionTrace` now collapses analyst material by default.
- `decisionExplanationModel` compacts fallback analyst conclusions in the story summary.
- Full analysis remains available through the explicit expand control.

Validation:

- `DecisionDetailPage.test.tsx` verifies long analyst material is hidden until the user expands it.
- `decisionExplanationModel.test.ts` verifies fallback analyst conclusions are compact.
- Rendered `/decisions/decision_6f9fa7db5afe919a` showed final verdict and safety boundary first, with `展开 3 份分析材料`.

### Decision Loop Deep Link

Root cause: `DecisionLoopPage` ignored `decision_id` query parameters and always loaded the recent list.

Fix: when `decision_id` is present, the page uses `getDecisionLoop(decisionId)` and renders a focused single-record view.

Validation:

- `DecisionLoopPage.test.tsx` verifies `getDecisionLoop('decision_focus')` is used, recent-list loading is skipped, and the focused record retains read-only trace links.
- Rendered `/decision-loop?decision_id=decision_6f9fa7db5afe919a` showed `闭环条数：1` and `当前聚焦：decision_6f9fa7db5afe919a · 510300`.

## Linked-Flow Validation

This P103 validation is scoped linked-flow validation, not a claim that every UI control and every error branch in the product was exhaustively interacted with in this turn.

Current linked-flow evidence covers:

- Portfolio facts: current local portfolio state renders without generic failure; empty snapshot handling is covered by focused test.
- Decision detail: existing real-model decision `decision_6f9fa7db5afe919a` renders final verdict, safety boundary, compact reasons, and collapsed analyst material.
- Decision loop: focused deep link renders the matching decision, manual planned confirmation state, audit links, and missing risk clue.
- Prior P102 evidence remains valid for the real workflow write/readback: portfolio adjustment, real-model consult, planned confirmation, decision-loop API, audit/review/workbench/notifications DOM readback, and SQLite readback.

Historical full-requirement coverage remains governed by P92/P93:

- P92 confirms the original requirement ledger status.
- P93 confirms production-code reality and release-claim safety.

P103 does not expand those claims. It only fixes and revalidates the P102 UX/linkage findings.

## Boundaries

P103 does not claim:

- Docker install/upgrade/uninstall validation.
- GitHub Release validation.
- Distribution package refresh.
- Physical second-machine verification.
- Broker connectivity.
- Automatic trading, one-click trading, delegated orders, external push, automatic confirmation, or automatic rule application.
- Paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
