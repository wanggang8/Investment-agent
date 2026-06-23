# Design: P66 Current Data Zero-Degradation Policy

## Design Brief

P66 makes current data-source quality actionable. Today the regression service says `passed`, `degraded`, or `failed`; release materials then interpret those statuses manually. P66 adds a policy layer that says whether a current run is safe for release claims, requires a waiver, or blocks the release.

The policy must be conservative but not noisy. Core market/source-health categories should block when degraded. Optional or lower-priority sources can require a waiver if they are recognized and classified. Unknown or unrecognized states always block.

## Policy Shape

Add `policy` to `DataSourceQualityRegressionResponse`:

- `verdict`: `passed`, `waiver_required`, or `blocked`.
- `release_gate`: `pass`, `waiver_required`, or `block`.
- `degraded_count`, `failed_count`, `blocking_count`, `waiver_count`.
- `blocking_reasons`: sanitized strings explaining why release is blocked.
- `waiver_reasons`: sanitized strings explaining what must be documented.
- `next_actions`: manual local actions only.
- `safety_note`: repeats that the policy is read-only and does not repair or refresh data.

Fixture mode should return `passed` because fixture cases are deterministic classification assertions. Current mode should evaluate actual source-health cases.

## Classification Rules

Current policy rules:

1. If no source-health facts exist, verdict is `blocked`.
2. If any case has `status=failed`, verdict is `blocked`.
3. If any unrecognized freshness or failure category appears, verdict is `blocked`.
4. If a core category is degraded, verdict is `blocked`.
5. If only recognized optional categories are degraded, verdict is `waiver_required`.
6. If every case passes, verdict is `passed`.

Core categories are:

- `index_constituents`
- `index_weights`
- `index_valuation_files`
- any category with `source_level=A`
- any category with `source_type=index_basic`

Optional categories include lower-priority or proxy categories such as `sentiment_proxy` and other B/C-level non-core facts. P66 does not invent missing categories beyond the current local source-health metadata. Missing source-health metadata itself is blocking because the policy cannot verify current quality.

The policy must not treat the existing P48 `missing_categories` response field as a direct synonym for missing source-health metadata. Today that field also records degraded categories. P66 should block only when there are no source-health facts or when a source-health category is explicitly `missing`/failed/core-degraded. A recognized optional degraded category may still appear in `missing_categories` and must remain `waiver_required`, not automatically `blocked`.

Recognized failure categories should use the same bounded vocabulary as freshness classification: `no_data`, `source_unavailable`, `parse_error`, `stale`, `missing`, `unknown`, and empty. Any non-empty failure category outside that vocabulary blocks the policy even if freshness itself is recognized.

## API And CLI

The existing API response for `GET /api/v1/data-source-quality/regression` should include `policy`. This is not a new endpoint.

The existing CLI task should include policy in its compact audit output. Add a strict release gate option for local acceptance so a command can fail when the policy is blocked. The strict gate must not write raw diagnostics or attempt repairs.

## Frontend

The existing Data Quality page should display the policy verdict in the current data-source area. It should use the existing tone system:

- `passed` -> success.
- `waiver_required` -> warning.
- `blocked` -> danger.

The page should list short policy reasons and manual next actions. It must not add buttons that imply repair, refresh, auto-confirmation, rule application, or trading.

## Release Materials

Update release repeatability and handoff materials to state that future release-ready claims must include current-data policy evidence:

- `passed` can support a clean current-data claim.
- `waiver_required` must be explicitly documented as a waiver and cannot be called clean.
- `blocked` blocks release-ready claims until resolved or the release scope explicitly excludes current local data health.

## Safety Boundaries

P66 remains local and read-only:

- no public-provider calls;
- no LLM-provider calls;
- no new data sources;
- no automatic refresh, repair, migration, restore, rollback, or overwrite;
- no broker interface, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, or automatic rule application;
- no login-gated, paid, authorization-gated, Level2, or high-frequency data source;
- no future provider availability, current data health, or investment return promise unless the evidence supports that exact claim.
