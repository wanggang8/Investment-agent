# Design: P77 Requirements Real Pass Upgrade Gate

## Decision

P77 will not mutate P75's historical acceptance matrix. It will create a new upgrade layer that references stable P75 requirement IDs and records the P77-specific decision for each row.

## Real-Pass Gate

A row may be upgraded to `real_pass` only when the evidence package proves every applicable dimension:

- implementation evidence exists and is traceable to code or runtime behavior;
- real UI evidence exists for user-visible behavior, unless the row is explicitly non-UI;
- data impact is verified for mutating flows, including changed tables, prohibited tables, audit events, and readback;
- workflow, rule, LLM, RAG, collector, or data-source behavior is verified when the requirement depends on it;
- scenario evidence is not limited to a single incompatible path unless the requirement itself is single-scope;
- safety evidence confirms no broker, trading, push, auto-confirm, auto-rule, auto-repair, provider-availability, or return-promise overclaim;
- evidence is fresh for P77 or explicitly accepted as deterministic current evidence.

Rows that do not satisfy every applicable dimension remain `partial`, `scoped_pass`, `deterministic_local_evidence`, `reference_only`, `not_implemented`, or `blocked`.

## First-Batch Strategy

The first P77 pass should be conservative. It may upgrade rows that have strong direct evidence from P75 reruns and deterministic safety checks, such as:

- explicit non-prediction, no proactive specific-target recommendation, no return promise, and user-final-decision boundaries;
- no automatic trading, broker, one-click, order delegation, external push, automatic confirmation, automatic rule application, automatic repair/migration/restore, or real database overwrite capabilities;
- rows proven by fresh real UI action plus SQLite and readback evidence;
- accepted-local non-`510300` evidence only where the claim is limited to accepted-local dynamic binding, not arbitrary live-source coverage.

Rows whose evidence remains only broad, historical, single-symbol, fixture-only, screenshot-only, route-smoke-only, or API-only for user-visible behavior must not be upgraded.

## Artifacts

P77 will produce:

- `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`
- `docs/release/acceptance/2026-06-21-p77-real-pass-upgrade-acceptance.md`
- `docs/release/ui-audit-assets/2026-06-21-p77/real-pass-upgrade-summary.json`

The summary JSON will include counts by original P75 status, P77 status, upgraded rows, remaining full-release-required gaps, and evidence commands.

## Release Claim Policy

P77 may produce `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`. Otherwise it must preserve a scoped conclusion such as `release_ready_scoped_with_p77_real_pass_progress` and list the exact remaining gaps.
