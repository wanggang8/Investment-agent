# P92 Final Original Requirement Audit Ledger

## Why

After P90 and P91, the product has functional acceptance and deployment readiness evidence, but the user asked for a final independent review that walks the original requirements row by row and ties each requirement to UI behavior, data effects, readback evidence, audit/safety proof, and final release status. This stage makes that review explicit and repeatable.

## What Changes

- Add a generator/checker for a final original-requirement audit ledger.
- Generate a final matrix derived from `docs/requirements.md` and P75-P90 evidence.
- Generate a concise summary grouped by original requirement section and review dimension.
- Verify that all full-release-required rows end as `real_pass` after applying P89/P90 final evidence.
- Record P91 deployment as release/distribution readiness, separate from product behavior.

## Out Of Scope

- No product runtime behavior changes.
- No new UI routes, APIs, SQLite schema, workflow nodes, LLM prompts, providers, or investment logic.
- No physical second-machine validation.
- No broker integration, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/auth-only sources, Level2 data, or high-frequency data.

## Acceptance

P92 is acceptable only if:

- The generated ledger contains every original P75 requirement row.
- Every full-release-required row has a final status and there are zero non-`real_pass` full-release-required rows after P90 overlay.
- Reference-only rows remain separated from product pass claims.
- Each ledger row includes at least: requirement id, source section, requirement text, final status, feature area, UI/product surface, expected behavior/data impact, readback/audit evidence, acceptance command/artifact, and boundary notes.
- The summary clearly states remaining non-product gaps, especially physical second-machine validation and forbidden capabilities that remain intentionally absent.
- Validation includes the P92 checker, OpenSpec strict validation, and `git diff --check`.
