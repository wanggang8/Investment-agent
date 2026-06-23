# P75 Requirements Traceability And Real Use Closure Design

## Design Summary

P75 treats `docs/requirements.md` as the full product truth and previous P71-P74 results as scoped evidence, not as automatic full-product proof. The phase has two outputs:

1. a machine-readable and human-readable traceability matrix;
2. a real-use acceptance runner that verifies representative user journeys through UI/API/data/audit evidence.

The core rule is simple: a requirement can be marked `real_pass` only when there is current evidence for implementation behavior, UI operation where relevant, real non-stub product data where applicable, data-impact/readback where relevant, and safety/release-claim boundaries. Fixture, mock, stub, temporary-only, waiver, scope-exclusion, or single-symbol evidence can support deterministic checks, but it cannot by itself become `real_pass`. Anything narrower becomes `scoped_pass`, `partial`, `not_implemented`, or `blocked`.

## Status Taxonomy

| Status | Meaning |
| --- | --- |
| `real_pass` | The requirement has current end-to-end evidence against real local product behavior, including real non-stub data where applicable and UI/data impact when applicable. |
| `scoped_pass` | The requirement passed only for an explicit accepted scope, such as `510300`/`000300`, a temporary database, or a single user journey. |
| `deterministic_local_evidence` | Fixture or temporary local facts prove calculation or state-transition correctness, but do not prove real data/source readiness by themselves. |
| `partial` | Some implementation exists, but one or more required evidence dimensions are missing. |
| `not_implemented` | No meaningful product implementation exists beyond docs or placeholders. |
| `blocked` | The requirement cannot safely pass because data, source, provider, design, or product constraints prevent it. |

## Traceability Matrix

P75 will add a traceability artifact under `docs/release/acceptance/` with one row per atomic requirement. Atomic means every normative paragraph, bullet, table row, SOP step, acceptance bullet, and safety/compliance statement in `docs/requirements.md` sections 1-19 gets a stable ID and line range. Each row must include:

- requirement ID, source section, source start/end lines, and requirement text hash;
- product expectation;
- criticality, criticality reason, full-release requirement flag, non-goal or optional basis when applicable, and allowed release claim;
- implementation evidence;
- UI evidence;
- external/built-in data evidence;
- LLM/workflow/rule evidence;
- scenario and data-impact evidence;
- security/safety boundary evidence;
- delivered-by change, verification command, acceptance artifact, and evidence freshness;
- status;
- gap;
- remediation task or explicit no-go reason;
- release-claim impact.

The first pass is deliberately allowed to find failures. A P75 pass requires either closing every `full_release_required=true` gap or downgrading the release conclusion honestly.

## Real-Use Acceptance Dimensions

P75 acceptance must cover the user's specific concerns:

- **Dynamic fund input:** user enters a fund/ETF symbol; the product resolves what data to query instead of relying on hardcoded `510300`.
- **External data usage:** fund-side and tracked-index-side data dependencies are evaluated from the entered symbol; at least one non-`510300` real or accepted-local read-only collector/request-construction path must prove request parameters, stored facts, source health, freshness, audit events, readiness, and correlation keys are bound to that symbol and tracked index, or the claim is downgraded.
- **Built-in knowledge usage:** master wisdom and SOPs are structured, visible, and actually referenced by workflow/LLM context without becoming formal evidence.
- **Data completeness:** required categories from the original requirements are checked: market price, valuation, liquidity, sentiment proxy, formal evidence, RAG/index health, fund profile, tracked index, and where applicable funds flow, margin financing, and constituent financials.
- **Missing-data propagation:** missing funds flow, margin financing, constituent financials, media heat, sentiment proxy, valuation, liquidity, or formal evidence must block or downgrade the exact claims that depend on them.
- **Accuracy checks:** expected-return, risk alert, liquidity, valuation-zone, source-verification, and manual-confirmation calculations are verified against deterministic expected values.
- **Field-level joins:** fund profile, NAV/price/liquidity, tracked index, benchmark symbol/profile, index valuation, constituent/financial data, and formal evidence must expose source, join key, as-of date, freshness, and conflict handling.
- **Cross-feature impact:** a UI action must be traced through SQLite state, derived pages, decision detail, daily discipline, risk alerts, notifications, audit, and release/readiness readback where applicable, using an action-to-table-to-page truth table.
- **SOP scenarios:** A-F scenarios must be mapped to accepted evidence or explicit gap status. User-visible SOP behavior requires real browser UI operation; API evidence can only supplement rule/data/database assertions.
- **Self-check and evolution:** daily, monthly, and quarterly review outputs plus error-case, proposal, and gatekeeper flows must be verified against source data, UI readback, and audit.
- **UI design:** task flows must be reviewed from a real-user perspective: discoverability, wording, mobile layout, states, next action clarity, and absence of misleading trading affordances.
- **Repeatability gates:** final release conclusions must inherit P52 G0-G9, P66/P67 current-data policy/resolution, P71-P74 repeatability rules, and P52 G6/G7 failure classifications by rerunning concrete repeat commands or explicitly downgrading inherited historical/scoped evidence.
- **Release claim:** no scoped pass may be used as a full-product pass.

## Gap Closure Strategy

P75 will not blindly implement every missing requirement. It will first classify gaps:

- **Must fix before full product pass:** hardcoded symbol paths, missing data blockers that affect accepted claims, unsafe UI wording, LLM context inconsistency, false release language.
- **Can remain scoped with honest claim:** sources that are unavailable without paid/login/authorized access, physical second-machine repeat, future provider availability.
- **Must remain out of scope:** broker connectivity, automatic trading, return promises, Level2/high-frequency data, automatic external push, automatic confirmation, automatic rule application.

## Acceptance Bar

P75 can end in one of three outcomes. Any use of "product-critical", "optional", "non-goal", or "future optional work" must be backed by a row-level criticality decision that cites `docs/requirements.md` text or a P75 L1 delta.

- `release_ready_full_requirements_traceable`: every `full_release_required=true` original requirement is `real_pass`, and remaining gaps are explicitly non-goals or future optional work.
- `release_ready_scoped_with_traceability_gaps`: the product remains usable and safety-reviewed, but some original requirements are partial/scoped and cannot be marketed as full completion.
- `release_pending_safety_review_scoped_with_traceability_gaps`: the product has scoped/partial traceability gaps and the expanded G9 forbidden-term scan still needs human boundary review, so clean release-ready wording is not allowed.
- `release_blocked_requirements_traceability`: critical gaps make the current release claim unsafe or misleading.
