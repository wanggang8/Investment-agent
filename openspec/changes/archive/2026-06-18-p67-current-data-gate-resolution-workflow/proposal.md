# P67 Current Data Gate Resolution Workflow

## Why

P66 made current data-source degradation machine-checkable: the current local database now returns `policy=blocked` and `gate=block`. That is useful, but incomplete for future releases. A release owner still needs a local, auditable way to record whether the blocked gate was handled by an explicit waiver or by excluding current local data health from the release scope.

Without P67, release materials can only say "blocked" in prose. P67 turns that into a repeatable local workflow: show the policy, capture an explicit human resolution record, keep the original P66 policy intact, and expose whether future release claims are unresolved, waiver-required policy resolved with waiver, blocked policy handled by scope exclusion, or cleanly passed.

## What Changes

- Add a local `data_quality_gate_resolutions` persistence model for manual current-data gate resolution records.
- Add service/API support to list, create, and retire current-data gate resolution records.
- Add a resolution check that combines the P66 current policy with the one active matching `policy_fingerprint` resolution and returns release-claim state.
- Extend `/data-quality` so users can see the current block, record an explicit waiver or scope exclusion, and see the resulting release statement boundary.
- Add CLI acceptance support for checking the current-data resolution state without changing the P66 strict policy gate.
- Add release acceptance material documenting the P67 run, the active resolution state, and what is still not claimed.

## In Scope

- Local SQLite schema, repository, service, DTOs, HTTP handlers, wiring, CLI task, frontend page/model/service/type updates, tests, E2E smoke, OpenSpec/docs/progress updates.
- Resolution types: `waiver` for `waiver_required` policy only, and `scope_exclusion` for `blocked` or `waiver_required` policy.
- Resolution statuses: `active` and `retired`.
- Release claim states: `pass`, `requires_resolution`, `resolved_with_waiver`, `resolved_with_scope_exclusion`. `resolved_with_waiver` is not available for `blocked` policy.
- Sanitized reason/scope/evidence text and audit-safe compact output.

## Out of Scope

- No new public data provider calls, source refresh, source repair, market snapshot mutation, index rebuild, LLM calls, or collector execution.
- No broker interface, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair promise, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No weakening of P66: an active resolution does not convert `policy=blocked` into `policy=passed`, and must not allow current local data quality to be described as clean.

## Impact

- Backend: migration, repository, service, handler, CLI and tests.
- Frontend: `/data-quality` state model, forms/actions for local resolution records, tests and E2E safety coverage.
- Docs/OpenSpec: data-source-quality-regression spec delta, release acceptance material, repeatability notes, progress/governance updates.
