# P55 UI Design Audit

> Date: 2026-06-17  
> Method: Product Design audit from browser-operated UI screenshots  
> Evidence folder: `docs/release/ui-audit-assets/2026-06-17-p55/`

## Audit Result

Design quality is serviceable for desktop operational review, but mobile and error-resilience need work before calling the frontend fully accepted. This audit is screenshot and browser-behavior based; it does not claim complete WCAG compliance.

## Strengths

- The product boundary is consistently visible: pages repeatedly state that the system records local facts and supports manual review only.
- Desktop information architecture is understandable. The left navigation exposes all major workflows, and repeated cards make cross-page scanning predictable.
- High-risk areas such as risk alerts, daily discipline, local knowledge import, settings, and decision details keep no-auto-trading copy visible.
- Core workflows can be operated from the UI: portfolio calibration, consultation, manual confirmation, notification read state, risk SOP update, audit expansion, local install draft editing, and local knowledge import.

## Blocked Issue

| ID | Priority | Area | Finding | Recommendation |
| --- | --- | --- | --- | --- |
| P55-B1 | blocked | Real LLM decision detail | A real LLM consultation can produce `final_verdict.optional_actions: null`, causing the decision detail UI to crash at `DecisionTrace.tsx:44`. | In P56, normalize nullable arrays at the API or frontend adapter boundary and add a regression test using a real-like DTO with `optional_actions: null`. |

## Needs Optimization

| ID | Priority | Area | Finding | Evidence | Recommendation |
| --- | --- | --- | --- | --- |
| P55-D1 | important | Mobile layout | `/positions` overflows horizontally on 390px mobile, with measured width 729px. | `32-mobile-positions.png` | Stack calibration labels/inputs vertically, constrain table/form width, and add responsive wrapping for action rows. |
| P55-D2 | important | Mobile layout | `/data-quality` overflows horizontally on 390px mobile, with measured width 528px. | `30-mobile-data-quality.png` | Wrap long source-health tokens, avoid fixed-width cards, and allow dense data rows to become vertical definition lists on mobile. |
| P55-D3 | important | Mobile navigation | The fixed sidebar remains visible on mobile and consumes a large portion of the viewport, leaving the main content too narrow even when there is no horizontal overflow. | `29-mobile-workbench.png` | Convert sidebar to a compact top nav, drawer, or collapsible rail below a mobile breakpoint. |
| P55-D4 | important | Form usability | Several forms use browser-default inputs with tight inline labels, making finance workflows harder to scan and operate. | `24-local-knowledge-validated.png`, `32-mobile-positions.png` | Use consistent field groups, stacked labels, stable input widths, and clearer primary/secondary button styling. |

## Minor Findings

| ID | Priority | Area | Finding | Recommendation |
| --- | --- | --- | --- | --- |
| P55-M1 | minor | Accessibility semantics | Some section labels such as "最终裁决" are visually useful but not always robust as exact accessible text targets. | Prefer semantic headings or `aria-labelledby` for key card sections. |
| P55-M2 | minor | Long-page capture | Browser full-page screenshots timed out or produced repeated content on long decision pages; fixed clip screenshots were stable. | Long decision pages should have clearer anchored sections and less excessive vertical repetition. |
| P55-M3 | minor | Visual density | Desktop dashboard and decision pages are readable but card-heavy and vertically long. | Keep operational density, but reduce repeated card padding and use clearer compact section hierarchy for recurring facts. |

## Accessibility Notes

- Keyboard-targetable controls were present for the checked forms and buttons, but this run did not perform a complete keyboard-only audit.
- The fixed mobile sidebar creates a practical accessibility issue because it reduces reading width and increases scroll effort.
- No complete WCAG conformance statement is made from this evidence.

## Recommended Next Phase

Create P56 as a UI acceptance fix phase:

- Fix nullable decision DTO rendering for real LLM outputs.
- Add frontend regression tests for nullable `final_verdict.prohibited_actions` and `final_verdict.optional_actions`.
- Fix `/positions` and `/data-quality` mobile overflow.
- Improve mobile navigation behavior.
- Re-run the P55 route matrix after fixes and update release status only if the real LLM UI path renders successfully.
