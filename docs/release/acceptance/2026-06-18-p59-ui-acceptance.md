# P59 UI Acceptance Run

> Date: 2026-06-18
> Change: `p59-decision-explainability-experience`
> Scope: `/consultation`, `/decisions/:decisionId`, `/evidence`, `/decision-loop`.

## Summary

P59 frontend experience acceptance passed for the scoped decision explanation surfaces.

The live app was started with `configs/config.local.yaml` against the local database. The temporary LLM key stayed in the untracked local config and was not written to tracked files or screenshots.

## Commands

| Gate | Command | Result |
| --- | --- | --- |
| Target frontend tests | `npm test -- --run src/features/decision/decisionExplanationModel.test.ts src/pages/DecisionDetailPage.test.tsx src/components/decision/DecisionTrace.test.tsx src/pages/EvidencePage.test.tsx src/pages/DecisionLoopPage.test.tsx` | Pass: 5 files, 22 tests |
| Full frontend tests | `npm test` | Pass: 34 files, 119 tests |
| Frontend build | `npm run build` | Pass |
| Backend tests | `go test ./...` | Pass |
| Real smoke / E2E | `bash scripts/e2e-smoke.sh` | Pass: 2 Playwright tests |
| OpenSpec | `openspec validate p59-decision-explainability-experience --strict && openspec validate --all --strict` | Pass: 33 items |
| Diff check | `git diff --check` | Pass |
| UI CTA scan | P59 browser metrics button/link scan | Pass: no forbidden execution CTA |
| Sensitive scan | Changed and untracked text files | Pass: no secret matches |

## Real LLM Result

Real LLM calls were attempted through the frontend consultation flow and the CLI smoke task. The Analyst nodes executed and recorded model metadata for `gpt-5.4-mini`, but the remote `/v1/chat/completions` endpoint returned HTTP 503. The UI correctly rendered the safe degraded state with `LLM 材料 0 份；解析/质量通过 0 份` and preserved the rule-based final verdict, safety boundary, evidence, audit, and decision-loop links.

This is recorded as an external dependency failure, not as a P59 UI implementation failure.

## Browser Acceptance

| Route | Viewport | Screenshot | Result |
| --- | --- | --- | --- |
| `/consultation` after real submit | 1280x900 | `docs/release/ui-audit-assets/2026-06-18-p59/desktop-consultation-result.png` | Pass |
| `/consultation` after real submit | 390x844 | `docs/release/ui-audit-assets/2026-06-18-p59/mobile-consultation-result.png` | Pass |
| `/decisions/decision_77f47ddbd489d66a` | 1280x900 | `docs/release/ui-audit-assets/2026-06-18-p59/desktop-decision-detail.png` | Pass |
| `/decisions/decision_77f47ddbd489d66a` | 390x844 | `docs/release/ui-audit-assets/2026-06-18-p59/mobile-decision-detail.png` | Pass |
| `/evidence` | 1280x900 | `docs/release/ui-audit-assets/2026-06-18-p59/desktop-evidence.png` | Pass |
| `/evidence` | 390x844 | `docs/release/ui-audit-assets/2026-06-18-p59/mobile-evidence.png` | Pass |
| `/decision-loop` | 1280x900 | `docs/release/ui-audit-assets/2026-06-18-p59/desktop-decision-loop.png` | Pass |
| `/decision-loop` | 390x844 | `docs/release/ui-audit-assets/2026-06-18-p59/mobile-decision-loop.png` | Pass |

Recorded metrics:

- `docs/release/ui-audit-assets/2026-06-18-p59/browser-results.json`
- Decision used for aligned screenshots: `decision_77f47ddbd489d66a`
- All P59 scoped desktop and mobile checks reported no page-level horizontal overflow.

## Findings

- Consultation result shows the generated explanation path, decision story, safety boundary, trust summary, and generated decision detail link after real UI submission.
- Decision detail now presents story-first explanation before the technical trace: verdict, context, safety boundary, key reasons, trust summary, then layered details.
- Evidence page exposes source trust and decision explanation navigation before raw evidence detail.
- Decision loop reads as a read-only lifecycle and links back to the latest decision detail without write-action buttons.
- The degraded LLM state remains safe: no automatic confirmation, trading, rule application, external push, auto-fix, or return promise entrypoint was exposed.

## Notes

- In-app browser interaction was used for live UI operation. Screenshot capture through the in-app browser timed out, so screenshots were captured with the project Playwright runtime against the same live Vite URL.
- A local capability setting for `510300` and a local test holding were created through local UI/API preconditions so consultation could proceed through the real workflow.
