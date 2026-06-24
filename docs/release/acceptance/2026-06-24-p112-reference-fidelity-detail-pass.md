# P112 Reference Fidelity Detail Pass Acceptance

Date: 2026-06-24

Change: `p112-reference-fidelity-detail-pass`

Reference image: `/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

Screenshot evidence directory: `docs/release/ui-audit-assets/2026-06-24-p112-reference-fidelity-detail-pass/`

## Scope

P112 is a visual fidelity closure pass after P111. It tightens the product UI against the selected high-fidelity reference direction across the full product, not only the home page.

P112 changes only frontend visual structure, layout density, navigation grouping, responsive rules, and screenshot evidence capture. It does not add backend APIs, SQLite schema, Eino workflow behavior, LLM behavior, data providers, broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, automatic repair, or return promises.

## Implemented Changes

- Reduced sidebar fragmentation from many small groups to three product groups: core workspace, system/evidence, and governance/local.
- Tightened reference surfaces: hero padding, card radius, border contrast, shadow strength, type scale, action density, metric card density, and checklist rhythm.
- Removed leftover pre-hero page headers from secondary pages so the first content surface is the report/status panel.
- Reworked mobile reference hero layout so status and prohibited actions remain visible in the first screen without stacking into an oversized hero.
- Converted mobile daily-hero metric summaries to a horizontal compact rail and shortened secondary action copy in the hero area.
- Changed the mobile topbar from sticky to static to avoid visual duplication during long-page capture and simplify real scrolling behavior.
- Replaced Browser-plugin full-page screenshot capture with local Playwright screenshot capture after the plugin produced repeated topbar stitching artifacts.

## Page Matrix

All routes were captured at desktop `1492x1058` and mobile `390x844` with local Playwright.

| Page | Route | Desktop | Mobile | Notes |
| --- | --- | --- | --- | --- |
| Dashboard | `/` | pass_with_minor_notes | pass_with_minor_notes | Seed data creates safe empty/readback states below the first screen. |
| Workbench | `/workbench` | pass_with_minor_notes | pass_with_minor_notes | Shares Dashboard rhythm and compact reference hero. |
| Positions | `/positions` | pass | pass | First screen shows report state before detailed maintenance content. |
| Data Quality | `/data-quality` | pass_with_minor_notes | pass_with_minor_notes | Highest mobile density; no overflow or blocking overlap. |
| Risk Alerts | `/risk-alerts` | pass | pass | Hero is compact and no old header remains. |
| Consultation | `/consultation` | pass | pass_with_minor_notes | Mobile hero is tall but first downstream content remains visible. |
| Decision Detail | `/decisions/decision_smoke_p30` | pass | pass_with_minor_notes | Explanation pages keep taller report hero by design. |
| Decision Loop | `/decision-loop` | pass | pass_with_minor_notes | No overflow; follows decision evidence rhythm. |
| Evidence | `/evidence` | pass | pass_with_minor_notes | No overflow; follows ledger rhythm. |
| Rules | `/rules` | pass | pass | Governance page no longer opens with generic page header. |
| Review | `/review` | pass | pass_with_minor_notes | No overflow; follows review report rhythm. |
| Audit | `/audit` | pass_with_minor_notes | pass | Tool-like detail remains below first-screen summary. |
| Notifications | `/notifications` | pass | pass | Mark-all action moved into compact report/action surface. |
| Daily Reports | `/daily-discipline/reports` | pass | pass | Report history opens with compact status surface. |
| Daily Auto Run | `/daily-auto-run` | pass | pass | Local automation boundary remains explicit. |
| Local Install | `/local-install` | pass_with_minor_notes | pass | Diagnostic detail remains below compact status surface. |
| Local Knowledge | `/local-knowledge` | pass | pass | Local import boundary remains explicit. |
| Settings | `/settings` | pass_with_minor_notes | pass_with_minor_notes | Capability labels can be further polished later; no blocker. |

## Metrics

Desktop:

- 18/18 routes: `overflow=false`.
- 18/18 routes: `scroll_w=win_w=1492`.
- Hero top y is stable at `92`.
- First downstream section appears between y `240` and y `322`.

Mobile:

- 18/18 routes: `overflow=false`.
- 18/18 routes: `scroll_w=win_w=390`.
- Mobile topbar is `position=static`.
- First downstream section appears between y `343` and y `489`.
- Long screenshots no longer contain repeated topbar/first-screen stitching.

## Subagent Review

Desktop subagent result: `pass_with_minor_notes`.

- No release-blocking desktop visual issues.
- No horizontal overflow.
- Hero is consistently the first content surface.
- Minor notes are seed-data empty states and a few tool-like dense areas below the first screen.
- Backend adjustment required: no.

Mobile subagent first result: `fail`.

- Blocking issue was screenshot evidence quality: Browser-plugin full-page screenshots repeated topbar/first-screen content.
- The issue was treated as blocking because it prevented honest mobile visual evidence.

Mobile corrective action:

- Mobile topbar changed to static.
- All screenshots recaptured using local Playwright full-page and first-viewport screenshots.

Mobile subagent second result: `pass_with_minor_notes`.

- Repeated topbar stitching is gone.
- No release-blocking mobile visual issues across 18 routes.
- No backend adjustment required.
- Only non-blocking note: `/data-quality` remains the densest mobile page.

## Backend Impact

No backend change is needed for P112. The remaining differences from the reference image are visual density, seed data, and page composition concerns, not API/data-model gaps.

## Validation Commands

Passed:

- `openspec validate p112-reference-fidelity-detail-pass --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- `go vet ./...`
- `openspec validate --all --strict`
- Forbidden runtime affordance scan over changed runtime UI files: no matching button/link affordances.
- Sensitive pattern scan over changed runtime UI files, P112 docs, and P112 OpenSpec package: no matches.
- `git diff --check`

Note: `python3 scripts/p93_code_reality_audit.py --check` was also tried but returned `status=failed` / `reason=stale:docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md` after P112 added new documentation. P112 therefore does not claim P93 historical checker pass; P112 scope validation uses the commands above plus fresh desktop/mobile screenshot and subagent review evidence.

## Conclusion

P112 resolves the P111 residual fidelity gap to a `pass_with_minor_notes` state across desktop and mobile. There are no remaining P0/P1/P2 or release-blocking visual findings from subagent review.
