# Investment Agent v0.1.3 Release Notes

Date: 2026-06-25  
Tag: `v0.1.3`  
Change: `p121-final-review-and-v0-1-3-tag-release`  
Status: `ready_after_p121_final_review`

## Summary

`v0.1.3` is a post-redesign source release that packages the P114-P120 product/UI acceptance layer after a fresh P121 release review. It focuses on the productized local investment research workstation experience, scenario acceptance evidence, and release-governance clarity.

This release does not add new investment runtime capability. It publishes the reviewed source state after the high-fidelity UI/productization pass and the expanded real-use acceptance runs.

## Highlights

- P114 repaired remaining productization and alignment issues after the high-fidelity redesign, including form/button alignment, card hierarchy, action baselines, `/decisions`, and `/api-diagnostics`.
- P115 executed 34 real user scenario checks and fixed the `DecisionTrace` duplicate React key issue found during acceptance.
- P116 added complex multi-fund transaction ledger acceptance across multiple funds/ETFs, buy/sell/reduce flows, invalid trade rejection, downstream readback, and mobile portfolio checks.
- P117 covered a 7-day continuous product usability story, including restart persistence and cross-page consistency.
- P118 covered edge scenarios such as long-history accumulation, malformed import recovery, stale/missing data handling, household/local account notes, mobile layout, and safety negative evidence.
- P119 covered all production UI controls and upstream/light interaction toggles, including 603 desktop-visible controls and 24 toggle interactions with zero reported layout/toggle issues.
- P120 archived the P114-P119 closure summary and kept the P93 stale boundary explicit.

## Verification

The final P121 gate results are recorded in `docs/release/acceptance/2026-06-25-p121-final-review-and-v0.1.3-tag-release.md`.

Current expected gates before tag publication:

| Gate | Required result |
| --- | --- |
| `openspec validate p121-final-review-and-v0-1-3-tag-release --strict` | passed |
| `openspec validate --all --strict` | passed |
| `go test ./...` | passed |
| `go vet ./...` | passed |
| `npm --prefix web test -- --run` | passed |
| `npm --prefix web run build` | passed |
| `python3 scripts/p92_final_requirement_audit.py --check` | passed |
| `python3 scripts/p121_final_release_review.py --check` | passed |
| `git diff --check` | passed |
| local release package smoke and verify for `v0.1.3` | passed |

## Scope Boundaries

P93 remains a historical final code-reality/design audit from before P114-P120. Because P114-P120 changed source and evidence, `v0.1.3` does not claim a fresh P93 pass. P121 is the fresh release-governance review for the current tree.

`v0.1.3` does not claim Docker installation validation, upgrade validation, uninstall validation, physical second-machine validation, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, fresh real LLM/provider quality, future provider availability, prediction accuracy, investment returns, login-only sources, paid sources, authorized sources, Level2 data, or high-frequency data unless separately validated.

## Tag Content

The annotated tag message should summarize:

- P114-P120 UI/product/real-use acceptance closure.
- P121 fresh release review and `v0.1.3` version synchronization.
- The exact gate list and package artifact identity.
- The explicit P93 stale boundary and safety exclusions.
