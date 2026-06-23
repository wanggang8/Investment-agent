# P56 UI Design Review

> Date: 2026-06-17
> Method: Product Design get-context / audit / research applied to implementation and browser screenshots
> Evidence folder: `docs/release/ui-audit-assets/2026-06-17-p56/`

## Result

P56 resolves the P55 blocked issue and materially improves the product feel of the local investment discipline cockpit. It does not claim full WCAG conformance or final release readiness; it narrows the P55 blockers and records remaining product opportunities for later phases.

## Product Design Inputs Applied

| Input | Applied change |
| --- | --- |
| Product Design get-context | Treated the app as a local investment discipline workbench for daily manual review, not a broker terminal or marketing landing page. |
| Product Design audit | Prioritized P55-B1, mobile overflow, mobile navigation, form usability and decision detail readability. |
| Product Design research | Used dashboard, form, mobile reflow, data table and financial trust/explainability references in `openspec/changes/archive/2026-06-17-p56-ui-acceptance-blocker-fixes/design.md`. |
| NN/g dashboard guidance | App shell and dashboard/workbench now emphasize task scanning and grouped navigation rather than an ungrouped route list. |
| NN/g form guidance | Consultation and positions workflows use consistent field spacing, labels, action rows and input sizing. |
| WCAG reflow guidance | Mobile pages now reflow to 390px without page-level horizontal scrolling. |
| Material data table guidance | Holdings table keeps desktop row/column scanning and switches to labeled rows on mobile. |
| Financial trust/explainability research | Decision detail keeps final verdict, prohibited actions, optional actions, evidence and analyst material visible without implying automated action. |

## Implemented Design Changes

| Area | Before | After |
| --- | --- | --- |
| Navigation | 17 flat links in a permanent sidebar. | Grouped navigation: 今日、决策、组合、证据、治理、系统. Mobile uses a `导航` toggle instead of a permanent sidebar. |
| Visual system | Demo-like root defaults and inconsistent card/control styling. | Operational CSS tokens, smaller page typography, consistent card radius, button, field, table and link-row styling. |
| Decision detail | Direct array assumptions could crash page on real DTO shape. | Nullable list rendering uses safe fallback text; real LLM decision detail renders without page crash. |
| Positions | Inline labels and table layout caused mobile width 729px in P55. | Responsive `form-grid`, wrapped action rows and mobile labeled holdings rows. |
| Data quality | Long source/status tokens caused mobile width 528px in P55. | Quality list cards wrap source identifiers and diagnostics inside page width. |
| Consultation | Browser-default form presentation. | Structured field grid, explicit safety note and consistent action row. |

## Remaining Non-Blocking Opportunities

- Dashboard and workbench can still be further consolidated around a single daily decision summary, but P56 already improved the shell and first-screen scanning enough for current acceptance.
- Rules, audit and local install pages still contain dense JSON/preformatted sections. P56 made global wrapping safer, but deeper product treatment should be a future independent UI phase if needed.
- Full keyboard-only and screen-reader audit remains out of scope; P56 only validates visible semantics, accessible labels touched by tests, and mobile reflow.

## Safety Boundary

The productized UI preserves the core safety posture:

- LLM output remains analysis material; final verdict remains rule/discipline led.
- User confirmation remains a record of offline/manual action, not execution.
- The UI does not add broker APIs, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, automatic repair promises, real database overwrite or return promises.
