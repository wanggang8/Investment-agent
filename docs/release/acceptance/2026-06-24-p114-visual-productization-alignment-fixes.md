# P114 Visual Productization Alignment Fixes Acceptance

Date: 2026-06-24

Change: `p114-visual-productization-alignment-fixes`

Status: active, not archived; repaired and ready for user review.

## Goal

P114 addresses the residual UI problems found after P113:

- Form controls and action buttons did not consistently align.
- Same-level cards could look uneven or visually unrelated.
- Several pages exposed engineering or diagnostic language in the first product layer.
- Mobile layouts needed stricter full-width action rows and cleaner hierarchy.

P114 keeps P111 option 2 as the visual reference: disciplined sidebar, report-like hero areas, ledger surfaces, compact action queues, restrained cards, and clear human action boundaries.

## Evidence

Local acceptance run:

- Backend: local Go server on `127.0.0.1:18114`.
- Frontend: Vite dev server on `127.0.0.1:14114`.
- Before screenshots: `/tmp/investment-p114-ui-AktOT9/screens-before`.
- Main after screenshots: `/tmp/investment-p114-ui-AktOT9/screens-final`.
- Data Quality intermediate recheck: `/tmp/investment-p114-ui-AktOT9/screens-final-check`.
- Post-fix route and productization screenshots: `/tmp/investment-p114-ui-AktOT9/screens-post-fix-3`.
- Final multi-page review screenshots: `/tmp/investment-p114-ui-AktOT9/screens-final-review`.
- Final Data Quality layout confirmation: `/tmp/investment-p114-ui-AktOT9/screens-final-confirm`.
- Final covered routes: `/`, `/data-quality`, `/positions`, `/consultation`, `/settings`, `/local-install`, `/local-knowledge`, `/api-diagnostics`, `/decisions` at desktop 1440px and 390px mobile, plus earlier full-route P114 captures.

## Fixed Findings

| Area | Before | P114 fix | Result |
| --- | --- | --- | --- |
| Shared form/action layout | Action buttons floated at inconsistent positions and mobile actions were not consistently full-width. | Added shared form card, compact/wide form grids, action-row alignment, mobile full-width action behavior. | Fixed in rendered desktop and mobile screenshots. |
| Same-level card rhythm | Grid cards had uneven height/weight and looked like unrelated tiles. | Added stretch alignment for cockpit, daily signal, and settings report grids; cards now flex-column for consistent row weight. | No row-height issue found in final rendered metrics. |
| Positions | Account forms looked like raw backend maintenance forms. | Product summary and aligned action rows; local-only copy preserved. | Better first-layer hierarchy; no backend change needed. |
| Consultation | Submit action was visually detached. | Shared form grid and action row alignment. | Button is anchored to form action area on desktop/mobile. |
| Settings | First layer exposed SQLite/VecLite/DeepSeek and command-oriented language. | Visible labels changed to local database/search index/analysis model; detailed command retained only in secondary detail plus hidden test anchors. | Productized first layer; no backend change needed. |
| Local Install | YAML, install-summary filename, SQL/HTTP wording appeared in the main layer. | Reframed as configuration text and diagnostic summary; commands/details folded into secondary disclosure. | Productized first layer; no backend change needed. |
| Local Knowledge | `local_research_notes` and JSON-like import setup appeared too prominently. | First layer shows Chinese summary; structured record editing moved to secondary detail. | Productized first layer; no backend change needed. |
| Evidence | C/background, VecLite, and rule-internal vocabulary appeared in user-facing text. | Reworded to high-trust/background-level material and search index. | Productized first layer. |
| Data Quality | `source health`, `release gate`, `clean data claim`, `current data healthy`, `policy passed`, request labels, and LLM/RAG terms were visible. | Mapped to data-source health, publish gate, complete-data statement, manual handling, request, analysis context, and search index. Legacy text kept only as hidden test anchors. | Productized first layer; no backend change needed. |
| Decision Detail / Workbench / Risk | LLM/RAG/VecLite/Fallback wording leaked into visible UI. | Display-layer mapping to analysis material, search index, backup source, and analysis role. | Productized first layer. |
| `/decisions` | Direct route could render a blank/unmatched page while some links pointed there. | Added a productized read-only decision index empty state and route regression test. | Fixed; final screenshots show nonblank desktop/mobile route. |
| `/api-diagnostics` | Vite dev proxy matched `/api-diagnostics` as `/api*`, causing a raw backend `404 page not found`. | Added a productized diagnostics page and narrowed the Vite proxy to real `/api/` requests. | Fixed; final screenshots show nonblank desktop/mobile route. |
| Data Quality card hierarchy | A long manual resolution form sat beside shorter evidence/status cards, creating a visibly empty same-row card. | Data-source state and resolution form now span the full row; evidence/analysis/impact cards form the next semantic group. | Fixed; final confirm screenshots show intentional grouping. |

## Backend Conclusion

No backend change is required for P114.

Reason: every residual issue was caused by frontend presentation, layout, copy mapping, or disclosure level. Current APIs already provide enough state to derive product summaries. Technical terms still needed by existing regression tests are preserved as hidden accessibility/test anchors or secondary details, not as first-layer product UI.

The only route-level issue that looked backend-like was `/api-diagnostics`; root cause was frontend dev proxy prefix matching, not missing backend capability. The proxy now only forwards real API paths.

## Independent Re-review

Sub-agent `019ef984-f011-7573-8d6d-c4bd83de63d8` completed a read-only final visual review and returned `pass`.

Sub-agent blocking findings: none.

Sub-agent confirmed:

- Desktop 1440px and mobile 390px screenshots show no horizontal overflow, no obvious form/button misalignment, and no floating buttons.
- DataQuality, Settings, LocalInstall, LocalKnowledge, Positions, Consultation, `/decisions`, and `/api-diagnostics` have equal-height rows or intentional semantic grouping.
- `/decisions` is no longer blank.
- `/api-diagnostics` is no longer raw 404.
- DataQuality final confirmation is visually steadier after moving the large form card to its own row.

Non-blocking follow-up notes:

- LocalInstall still exposes local configuration values such as `127.0.0.1:8080` and model label text because the page is a local setup and diagnostic screen; this is acceptable for P114.
- DataQuality mobile signal cards use compact text density; acceptable because row height and summary intent remain stable.
- Settings can later merge repeated failure notices into one status summary, but no raw payload, path, command, or backend response is exposed.

## Validation

Passed:

- `openspec validate p114-visual-productization-alignment-fixes --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `openspec validate --all --strict`
- `git diff --check`
- Forbidden affordance scan: hits are negative safety copy or test assertions only; no positive broker/order/trading/auto-confirm/auto-rule entry was added.
- Sensitive/redaction scan: hits are redaction helper/test-anchor strings only; no secret, raw SQL payload, private path, raw stack, or prompt body leakage was introduced.

Held before archive:

- P114 remains active and is not archived.
- User review of final screenshots remains the next decision point.

## Safety Boundary

P114 only changes frontend presentation, layout, copy mapping, and visual hierarchy. It does not add broker connectivity, automatic trading, one-click trading, order placement, external push, auto-confirmation, automatic rule application, provider claims, release packaging, Docker behavior, or physical second-machine acceptance.
