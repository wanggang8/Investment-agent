# Release Candidate: 2026-06-18

> Status: `release_ready_scoped_with_p86_final_integrated_progress`
> Acceptance run: `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`
> P68 decision: `docs/release/acceptance/2026-06-18-p68-release-readiness-governance.md`
> P70 final decision: `docs/release/acceptance/2026-06-18-p70-final-release-decision.md`
> P71 strict acceptance: `docs/release/acceptance/2026-06-18-p71-real-product-acceptance.md`
> P72 real user scenario: `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md`
> P73 product-effectiveness/UX validation: `docs/release/acceptance/2026-06-19-p73-product-effectiveness-ux-validation.md`
> P74 built-in knowledge/data readiness: `docs/release/acceptance/2026-06-19-p74-built-in-knowledge-and-data-readiness.md`
> P75 traceability closure: `docs/release/acceptance/2026-06-20-p75-real-use-closure.md`
> P83 governance traceability backfill: `docs/release/acceptance/2026-06-22-p83-governance-traceability-backfill.md`
> P84 portfolio confirmation data-impact closure: `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-closure.md`
> P85 expected return analysis-accuracy closure: `docs/release/acceptance/2026-06-22-p85-expected-return-analysis-accuracy-closure.md`
> P87 portfolio state allocation safety closure: `docs/release/acceptance/2026-06-22-p87-portfolio-state-allocation-safety-closure.md`
> P86 core goal knowledge safety final closure: `docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-closure.md`
> P76 package source commit: `8a317f25917b8ff18ec9b5049e6a6188206a22d3`
> P76 package: `tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz`
> Change: P63 `p63-full-ui-regression-release-refresh`; P68 governance refresh after P66/P67; P70 final risk closure after P69 package refresh; P71 strict real product acceptance true pass; P72 real user fund scenario data-impact acceptance; P75 scoped traceability closure; P76 post-P75 package refresh; P77 post-P75 real-pass upgrade gate; P78 real-pass batch closure; P79 real-use data-impact closure; P80 review/audit/governance closure; P81 dynamic source field coverage closure; P82 SOP/action UI-to-SQLite closure; P83 governance traceability backfill; P84 portfolio confirmation data-impact closure; P85 expected return analysis-accuracy closure; P87 portfolio state allocation safety closure; P86 core goal knowledge safety final closure

## Basis

- P52 acceptance matrix: `docs/project-acceptance-gate-matrix.md`
- P57 product polish roadmap: `docs/product-experience-polish-roadmap.md`
- P58-P62 product/UI hardening phases
- P63 acceptance execution: `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`
- P66 current-data policy gate: `docs/release/acceptance/2026-06-18-p66-current-data-policy.md`
- P67 current-data resolution: `docs/release/acceptance/2026-06-18-p67-current-data-resolution.md`
- P68 release readiness governance: `docs/release/acceptance/2026-06-18-p68-release-readiness-governance.md`
- P69 clean-tree package refresh: `docs/release/acceptance/2026-06-18-p69-clean-tree-package-refresh.md`
- P70 final release decision: `docs/release/acceptance/2026-06-18-p70-final-release-decision.md`
- P71 real product acceptance: `docs/release/acceptance/2026-06-18-p71-real-product-acceptance.md`
- P72 real user fund scenario acceptance: `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md`
- P83 governance traceability backfill: `docs/release/acceptance/2026-06-22-p83-governance-traceability-backfill.md`
- P84 portfolio confirmation data-impact closure: `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-closure.md`
- P85 expected return analysis-accuracy closure: `docs/release/acceptance/2026-06-22-p85-expected-return-analysis-accuracy-closure.md`
- P87 portfolio state allocation safety closure: `docs/release/acceptance/2026-06-22-p87-portfolio-state-allocation-safety-closure.md`

P63 refreshed the release candidate status after product experience polish. P68/P70 preserved that product/runtime readiness evidence but narrowed the top-level status because current local data health was excluded from clean release claims. P71 supersedes that limitation for the current accepted run: the current-data strict gate returns `policy=passed` / `gate=pass`, VecLite is healthy/fresh during real UI consultation, and the real LLM-backed UI journey passes without retrieval degradation.

P72 adds a deeper real-user scenario acceptance layer on top of P71. It verifies `510300` portfolio setup and maintenance, local knowledge import, formal public evidence collection, healthy VecLite retrieval, real LLM consultation, manual offline confirmation, and deterministic SQLite data impact through browser UI and readback pages.

P73 adds a stricter product-effectiveness and UX validation layer, focused on whether real UI tasks support discipline, evidence sufficiency, traceability, review usefulness, and safe manual confirmation. It passes with real browser UI operations, screenshots, browser results, and deterministic SQLite effect replay.

## Acceptance Summary

| Area | Result | Release impact |
| --- | --- | --- |
| Governance/OpenSpec | pass | does_not_block |
| Go tests | pass | does_not_block |
| Frontend tests/build | pass | does_not_block |
| Browser E2E smoke | pass | does_not_block |
| Full UI route regression | pass | does_not_block |
| Local fixture/current smoke | degraded | does_not_block |
| Current data policy gate | pass | supports_clean_current_data_claim_for_p71 |
| Current data resolution | pass | supports_clean_current_data_claim_for_p71 |
| Real public source opt-in | pass | does_not_block |
| Real LLM opt-in | pass | does_not_block |
| P71 strict real UI acceptance | pass | does_not_block |
| P71 VecLite retrieval gate | pass | does_not_block |
| P72 real user fund scenario | pass | does_not_block |
| P72 SQLite data-impact verification | pass | does_not_block |
| P73 product-effectiveness/UX validation | pass | supports_product_effectiveness_ux_acceptance |
| Local install/release upgrade | pass | does_not_block |
| Safety/redaction | pass | does_not_block |

## Release Notes

- Full Go test suite passed.
- Focused integration packages passed.
- Frontend Vitest and production build passed.
- Browser E2E smoke passed.
- P63 full UI regression passed across 20 primary routes and 390px, 768px, and 1280px viewports.
- P63 committed 60 route screenshots and a redacted browser summary JSON under `docs/release/ui-audit-assets/2026-06-18-p63/`.
- Real UI consultation generated decision `decision_e6f6d404bb554d61`; the detail page opened successfully and displayed three parsed, quality-passed LLM analyst reports.
- The full UI regression recorded all `/api/v1/` `>=400` responses: 20 were classified expected client-state 404/409 responses and zero were unexpected failed API responses.
- Recovery smoke, retrieval quality smoke, fixture data-source regression, and current data-source regression executed.
- Real public evidence refresh executed against a temporary SQLite database and wrote evidence, RAG, verification, and audit rows.
- Real LLM smoke executed with model `gpt-5.4-mini`; parse and quality gates passed.
- P71 real local strict acceptance executed with `use_stub=false`, real LLM config, temporary SQLite, real backend, and Vite frontend.
- P71 current-data strict gate returned `policy=passed` / `gate=pass` for `000300`; this is not based on scope exclusion.
- P71 real UI consultation generated decision `decision_a3aed494f6b84ac4` with `workflow_status=completed`, three parsed and quality-passed LLM analyst reports, and retrieval quality `hit` / `veclite` / `healthy` / `fresh`.
- P71 post-P70 package refresh generated `tmp/p71-final-release/20260618T101504Z/investment-agent-p71-real-product-acceptance.tar.gz` from clean source commit `2c195a05cee3b6cdda031e86409d562bcc7ee379`; package verify and isolated repeat acceptance passed.
- P72 real user scenario added `csindex_index` and `eastmoney_fund` formal evidence collection for `510300`, fixed formal-evidence/background merge behavior, and persisted VecLite rebuild status back to SQLite.
- P72 full runner passed with `use_stub=false`, real public sources, real LLM smoke, real local backend/Vite UI, and SQLite impact verification; final decision `decision_5da93489f6f7f6c1` was `workflow_status=completed`, `final_verdict_status=hold`, and `confirmation_status=executed_manually`.
- P72 deterministic data-impact check passed with final cash `95630.5`, total assets `101265.0`, two `510300` position rows, total quantity `1390.0`, total market value `5634.5`, `rag_chunks_p72_indexed=1`, and no forbidden trading/external-push tables.
- Local install diagnostics and release upgrade checks passed.
- Safety and redaction review found no committed complete key, raw payload, full prompt, private log, raw DB, or new prohibited capability.

## Known Degradations

| Item | Status | Impact |
| --- | --- | --- |
| Current data-source quality regression | degraded | The current local DB had one degraded case and zero failed cases. Fixture regression passed, so this does not block release but limits claims about the current local data snapshot. |
| P66/P67 historical current-data limitation | superseded by P71 for this run | P66/P67 remain accurate historical records for the P70 limited scope. P71 produced fresh strict evidence with `policy=passed` / `gate=pass` and `claim_state=pass`, so the P71 run may make a clean current-data gate claim. |
| P63 UI consultation workflow status | degraded | The real UI consultation returned HTTP 200, displayed LLM material, and opened decision detail. The degradation was `VECTOR_INDEX_UNAVAILABLE`, so claims about retrieval-enhanced context for this temporary run are limited. |
| P63 VecLite degradation | superseded by P71 for this run | P71 treats VecLite degradation as blocking and passed with `fallback_source=veclite`, `index_health=healthy`, and `index_freshness=fresh`. |
| P63 failed API response classification | classified | The browser regression recorded 11 expected 404 responses for `/api/v1/portfolio/current` and 9 expected 409 responses for `/api/v1/dashboard/today`; unexpected failed API responses were zero. |
| Initial parallel G5 current/retrieval command group | retried | The first grouped run hit SQLite `database is locked`; sequential gate reruns reached the expected pass/degraded outcomes. |
| P73 product-effectiveness/UX scope | accepted gap | P73 validates deterministic local product-effectiveness/UX behavior. It does not claim future returns, future market direction, future provider availability, a longitudinal real-world user study, or physical second-machine repeat. |
| P75 original-requirement traceability scope | accepted scoped gap | P75 records 341 atomic requirement rows and preserves `release_ready_scoped_with_traceability_gaps`; it is not a full original-requirement pass. |
| P76 package freshness | passed | P76 generated a clean package from source commit `8a317f25917b8ff18ec9b5049e6a6188206a22d3`; verify and isolated repeat acceptance passed, and committed P72-P75 evidence is included. |
| P77 real-pass upgrade gate | scoped progress | P77 generates a new upgrade layer from the P75 matrix: 17 rows are `real_pass`, 11 are `reference_only`, and 313 full-release-required rows remain non-real-pass; P77 does not refresh the P76 package. |
| P78 real-pass batch closure | scoped progress | P78 generates a new batch layer from the P77 matrix: 20 rows are now `real_pass`, 11 are `reference_only`, and 310 full-release-required rows remain non-real-pass; P78 does not refresh the P76 package. |
| P79 real-use data-impact closure | scoped progress | P79 generates a new evidence layer from the P78 matrix: 43 rows are now `real_pass`, 11 are `reference_only`, and 287 full-release-required rows remain non-real-pass; P79 does not refresh the P76 package. |
| P80 review/audit/governance closure | scoped progress | P80 generates a new evidence layer from the P79 matrix: 57 rows are now `real_pass`, 11 are `reference_only`, and 273 full-release-required rows remain non-real-pass; P80 does not refresh the P76 package. |
| P81 dynamic source field coverage | scoped progress | P81 generates a new evidence layer from the P80 matrix: 116 rows are now `real_pass`, 11 are `reference_only`, and 214 full-release-required rows remain non-real-pass; P81 does not refresh the P76 package. |
| P82 SOP/action UI-to-SQLite closure | scoped progress | P82 evaluates 53 SOP/action rows with fresh real browser UI, SQLite/readback, explicit final rule confirmation, and safety negative checks; 160 rows are now `real_pass`, 11 are `reference_only`, and 170 full-release-required rows remain non-real-pass; P82 does not refresh the P76 package. |
| P83 governance traceability backfill | scoped progress | P83 evaluates 43 candidate rows with fresh real browser UI, API readback, SQLite field checks, focused Go tests, and safety negative checks; 10 directly proven review/governance/release traceability rows upgrade, 33 broader rows defer to P86, 170 rows are now `real_pass`, 11 are `reference_only`, and 160 full-release-required rows remain non-real-pass; P83 does not refresh the P76 package. |
| P84 portfolio confirmation data-impact closure | scoped progress | P84 evaluates 35 candidate rows with fresh real browser UI, API readback, SQLite field checks, focused Go tests, downstream readbacks, and safety negative checks; 3 directly proven portfolio/confirmation rows upgrade, 32 broader rows defer to P85/P86, 173 rows are now `real_pass`, 11 are `reference_only`, and 157 full-release-required rows remain non-real-pass; P84 does not refresh the P76 package. |
| P85 expected return analysis-accuracy closure | scoped progress | P85 evaluates 31 candidate rows with fresh real browser UI consultation flows, target-return and previous-base-midpoint UI inputs, deterministic expected-return recomputation, SQLite readback, focused Go tests, degraded/unavailable sample handling, and safety negative checks; 15 directly proven rows upgrade, 16 broader historical-accuracy/backtest/probability rows defer, 188 rows are now `real_pass`, 11 are `reference_only`, and 142 full-release-required rows remain non-real-pass; P85 does not refresh the P76 package and does not claim fresh real LLM output in this environment. |
| P87 portfolio state allocation safety closure | scoped progress | P87 evaluates 32 candidate rows with fresh real `/positions` UI operation, API readback, SQLite field checks, focused handler/rule tests, sell-only/frozen-watch/information-insufficient decision readback, and safety negative checks; 5 directly proven rows upgrade, 27 broader source-transition/quarterly-rebalance/proposal/source-readiness/audit/release-safety rows defer, 193 rows are now `real_pass`, 11 are `reference_only`, and 137 full-release-required rows remain non-real-pass; P87 does not refresh the P76 package. |
| P86 core goal knowledge safety final closure | scoped progress | P86 replays P74/P81/P82/P83/P84/P85/P87 real UI/API/SQLite/Go evidence through the integrated runner, upgrades 110 additional rows, records 303 rows as `real_pass`, 11 as `reference_only`, and leaves 27 full-release-required rows non-real-pass; P86 does not refresh the P76 package and does not claim full original-requirement pass. |

## Not Claimed

This release candidate does not claim:

- Future availability of public websites or model providers.
- Future clean current local data health beyond the P71 evidence window.
- Future public-source/model-provider availability beyond the P72 evidence window.
- P67 scope exclusion as a P66 policy pass. P71 pass is based on fresh strict-gate evidence, not scope exclusion.
- Do not claim P76 package-after-the-fact evidence is inside the P76 archive.
- Investment returns or deterministic market outcomes.
- Broker connectivity, automatic trading, one-click trading, order delegation, or external push.
- Automatic confirmation, automatic rule application, automatic repair, automatic migration, or automatic overwrite of real user databases.
- Login, paid, authorized, Level2, or high-frequency data source coverage.
- A longitudinal real-world user study or future investment outcome validation.
- Full original-requirement pass beyond P86 `release_ready_scoped_with_p86_final_integrated_progress`.

## Decision

`release_ready_scoped_with_p86_final_integrated_progress`.

Proceed only with the documented safety boundaries above. P71 resolves the current-data and VecLite acceptance limitations for the current accepted run; P72 verifies a real user fund scenario and its SQLite data impact end to end; P73 verifies product-effectiveness and UX behavior for the accepted local scope; P74 verifies built-in knowledge/data readiness; P75 preserves a scoped-with-gaps original-requirement conclusion; P76 refreshes package evidence through committed P72-P75 materials; P77 upgrades the first 17 atomic rows to `real_pass`; P78 upgrades 3 additional expected-return degradation/disclaimer rows; P79 upgrades 23 additional portfolio/local-account/confirmation data-impact rows; P80 upgrades 14 additional review/audit/gatekeeper rows; P81 upgrades 59 dynamic source/readiness rows; P82 upgrades 44 SOP/action UI-to-SQLite rows and defers 9 broader rows; P83 upgrades 10 directly proven review/governance/release traceability rows; P84 upgrades 3 directly proven portfolio/confirmation rows; P85 upgrades 15 directly proven expected-return/readback rows; P87 upgrades 5 directly proven portfolio-state/allocation/buy-date rows; P86 upgrades 110 additional integrated rows and leaves 27 full-release-required rows non-real-pass. A separate physical second-machine run remains unperformed.
