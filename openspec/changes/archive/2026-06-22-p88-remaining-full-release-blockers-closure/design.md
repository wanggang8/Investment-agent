# P88 Design

## Evidence Strategy

P88 starts from `docs/release/acceptance/2026-06-22-p86-core-goal-knowledge-safety-final-matrix.md`. It owns only rows where `full_release_required=True` and `p86_status != real_pass`. The inventory must equal 27 rows; any drift blocks execution until the plan is updated.

Each P88 row can become `real_pass` only when the evidence directly proves the original row text. The minimum evidence package is:

- real browser UI operation for user-visible behavior;
- API/readback proof for the same state;
- read-only SQLite proof for persisted facts;
- workflow metadata or deterministic calculation proof where the row is workflow/model driven;
- forbidden-capability negative checks.

## Execution Tracks

### Track A: Source-Verified State Transitions

Rows: `REQ-02-022`, `REQ-02-025`, `REQ-06-023`, `REQ-06-024`, `REQ-17-015`.

Add or harden a deterministic source-verification transition path that reads formal evidence/source verification counts and writes local decision/portfolio states:

- buy logic broken with at least two independent A/S formal sources -> `sell_only`, prohibited buy/add actions, audit/readback;
- buy logic questioned, major positive information, or major negative information with fewer than two independent A/S formal sources -> `frozen_watch`, pause action, source-count provenance.

### Track B: Structured Data And Source Preverification

Rows: `REQ-04-016`, `REQ-04-025`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`.

Add a public-source preverification artifact for each P88 structured data category before treating it as production collector scope. `REQ-05-003`, `REQ-05-004`, and `REQ-05-005` may only upgrade to `real_pass` from a real runtime provider path with non-mock provenance, no-login/no-paid/no-Level2/no-high-frequency access, and SQLite readback. Accepted-local, fixture, stub, or manually seeded evidence may document a blocker, fallback, or parser contract, but cannot upgrade these structured-data rows to `real_pass`. The acceptance must prove field-level readback for:

- capital-flow `date`, `net_inflow`, `net_outflow`;
- margin-financing `date`, `margin_balance`, `balance_change_rate`;
- constituent financial `revenue`, `net_profit`, `growth`, `disclosure_date`.

### Track C: Expected Return Historical/Probability Closure

Rows: `REQ-08-004`, `REQ-08-023`, `REQ-09-001`, `REQ-09-003`, `REQ-09-004`, `REQ-09-006`, `REQ-09-007`, `REQ-09-008`, `REQ-09-009`, `REQ-09-010`, `REQ-09-013`, `REQ-09-023`, `REQ-09-024`, `REQ-09-025`, `REQ-09-027`.

Add a deterministic local historical-sample engine over persisted low-frequency market/valuation/fundamental inputs. It must expose sample count, sample window, screening condition, scenario ranges, probability basis, and missing supplement-data list. `REQ-09-001` requires a representative holding-class coverage matrix: broad ETF/index fund, sector/growth ETF or fund, and equity/security-like constituent-financial path. If P88 proves only one class/path, `REQ-09-001` must remain `partial` with exact uncovered classes. The engine must also support:

- extreme-fear active-trading lock with historical similar-scenario display;
- scenario rerun that visibly lowers affected probabilities;
- periodic assumption check;
- two-month below-expectation downshift warning;
- one-month pessimistic-path manual probability-adjustment suggestion;
- low-sample-below-5 degradation with no range and an explicit supplement-data list.

### Track D: Quarterly Rebalance

Rows: `REQ-10-004`.

Add quarterly +/-15% rebalance recommendation flow through UI/API/SQLite. The product may recommend offline manual buy/sell amounts but must not trade, confirm automatically, or write broker/order artifacts.

### Track E: SOP Addendum Proposal

Rows: `REQ-13-010`.

Use the existing rule proposal pipeline to generate a `sop` proposal when review/error-case evidence shows a high-frequency uncovered scenario. The proposal must create `rule_proposals`, `notifications`, and `audit_events`, remain pending for user review, and not modify active rules without the existing gatekeeper/final-confirm flow.

### Track F: P88 Matrix And Release Claims

Generate a P88 matrix and acceptance record. If all 27 rows become `real_pass`, P88 may state that the P88-owned full-release blockers reached 27/27 `real_pass` and may update the original-requirement matrix accordingly. Any reclassification requires explicit L1/OpenSpec rationale and must be reported separately from real-pass completion. If any row remains non-`real_pass`, P88 must list exact blockers and continue to avoid full-pass claims.

## UI Design Acceptance

P88 must verify the UI from a real user perspective:

- state transition labels explain why `sell_only` or `frozen_watch` was entered;
- expected-return report is readable and shows name/code, future-12-month label, samples, probabilities, missing data, and disclaimers;
- rebalance recommendations are clearly manual offline actions;
- SOP proposal pages distinguish pending proposal, review, audit, and final confirmation;
- mobile width has no overlapping text or unsafe action affordances.

## Safety Boundary

P88 must keep all forbidden capabilities absent. Any external data limitation must be surfaced as source unavailable or accepted-local evidence, not hidden as a pass.
