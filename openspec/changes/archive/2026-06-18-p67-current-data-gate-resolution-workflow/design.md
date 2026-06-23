# Design: P67 Current Data Gate Resolution Workflow

## Context

P66 introduced current-data policy verdicts and a strict gate. The current local DB is blocked by a core source-health degradation. P67 must not hide that fact. It should add a narrow local workflow for release owners to record how they handled the block and to prove that future release materials are not claiming a clean current-data state.

Product design brief: extend the existing `/data-quality` operational cockpit. Use current `daily-hero`, `daily-signal-grid`, `cockpit-card`, `quality-list`, `StatusNotice`, and link-row patterns. Interactivity is full local app interactivity: load current state, create a resolution, retire a resolution, and show updated state without a page refresh if practical. No new visual direction is required.

## Data Model

Add `data_quality_gate_resolutions`:

- `resolution_id` TEXT primary key.
- `symbol` TEXT.
- `policy_fingerprint` TEXT canonical fingerprint derived from symbol, policy verdict, release gate, blocking reasons, waiver reasons, degraded/failed/blocking/waiver counts, and normalized current regression case categories.
- `policy_verdict` TEXT and `release_gate` TEXT captured from the current P66 check.
- `policy_summary` TEXT compact sanitized summary.
- `resolution_type` TEXT: `waiver` or `scope_exclusion`.
- `status` TEXT: `active` or `retired`.
- `scope` TEXT sanitized human-readable release scope.
- `reason` TEXT sanitized human-readable reason.
- `release_impact` TEXT sanitized statement of what may and may not be claimed.
- `evidence_ref` TEXT optional sanitized local reference.
- `blocking_reasons_json` TEXT and `waiver_reasons_json` TEXT copied from P66 policy.
- `created_by` TEXT fixed local actor label, default `local_user`.
- `retired_by` TEXT nullable local actor label.
- `created_at` TEXT RFC3339 UTC.
- `retired_at` TEXT nullable RFC3339 UTC.
- `safety_note` TEXT fixed boundary copy.

No raw provider payload, private path, complete key, raw SQL, full prompt, raw HTTP exchange, stack trace, or local DB path may be stored.

## Service Rules

The service reads P66 current regression first and computes a canonical `policy_fingerprint`. The fingerprint, not the display summary, controls whether a resolution matches the current policy. `policy_summary` remains display/audit copy only.

Validation:

- `resolution_type=waiver` is allowed only when the current policy is `waiver_required`.
- `resolution_type=scope_exclusion` is allowed for `blocked` or `waiver_required`.
- `scope`, `reason`, and `release_impact` are required after trimming.
- Only one active resolution per `symbol + policy_fingerprint` may exist. Duplicate creates with the same type reuse the active record; attempts to create a different active type for the same fingerprint are rejected until the existing record is retired.
- Retire changes only the local resolution status and writes an audit event; it must not change market data or P66 policy.

Resolution check output:

- If P66 policy is `passed`: `release_claim_state=pass`, no resolution required.
- If P66 policy is `waiver_required` and no active resolution: `requires_resolution`.
- If P66 policy is `blocked` and no active resolution: `requires_resolution`.
- If active `waiver` and policy is `waiver_required`: `resolved_with_waiver`.
- If active `scope_exclusion`: `resolved_with_scope_exclusion`.
- `clean_data_claim_allowed` is true only when P66 policy is `passed`.

Claim labels are fixed:

- Allowed for `pass`: `可以声明当前本地数据门禁通过`.
- Allowed for `resolved_with_waiver`: `可以声明已记录当前数据质量豁免`.
- Allowed for `resolved_with_scope_exclusion`: `可以声明当前本地数据健康已排除在 clean claim 外`.
- Prohibited unless `pass`: `不得声明当前本地数据 clean`、`不得声明 current data healthy`、`不得把 resolution 描述为 policy passed`.

## API

Add local APIs under `/api/v1/data-source-quality`:

- `GET /gate-resolution?symbol=000300`: returns P66 current policy, active resolution if any, release claim state, allowed claims, prohibited claims, and safety note.
- `GET /resolutions?symbol=000300`: lists local resolution records newest first.
- `POST /resolutions`: creates/reuses a local manual resolution record from `symbol`, `resolution_type`, `scope`, `reason`, `release_impact`, optional `evidence_ref`.
- `POST /resolutions/{resolution_id}/retire`: retires a local resolution record.

POST APIs write local audit events with sanitized input/output refs. GET APIs are read-only and write no audit records.

## CLI

Add a local acceptance task:

```bash
go run ./cmd/agent --task data-source-quality-resolution-check --source current --symbol 000300
```

It returns exit 0 for `pass`, `resolved_with_waiver`, or `resolved_with_scope_exclusion`; exit 1 for `requires_resolution`. Because `resolved_with_waiver` is only valid for `waiver_required`, a `blocked` policy can exit 0 only through an active `scope_exclusion`. Output must include compact sanitized `policy`, `gate`, `fingerprint`, `resolution`, and `claim_state` fields. It must not change P66 strict gate behavior.

## Frontend

Extend `/data-quality`:

- Show a "当前数据门禁处置" section near the current data policy.
- Display release claim state, clean-data claim allowed/not allowed, and active resolution details.
- Provide a compact local form when policy requires resolution: resolution type, scope, reason, release impact, evidence ref. The form must only show `scope_exclusion` when policy is `blocked`; it may show `waiver` and `scope_exclusion` when policy is `waiver_required`.
- Provide a local retire action for active resolution.
- Keep all controls as local record actions, not data refresh or source repair.
- Sanitize displayed user text and keep forbidden capability words out of user-facing body and action labels.

## Testing

- Repository tests for create/reuse/list/retire.
- Service tests for pass, requires resolution, resolved with waiver, resolved with scope exclusion, duplicate reuse, retired record ignored, sanitization.
- Handler tests for GET/POST/retire and audit boundaries.
- CLI tests for unresolved exit 1 and resolved exit 0.
- Frontend model/page tests for blocked unresolved, create form, active resolution, retire action, safety copy and nullable arrays.
- E2E smoke must cover `/data-quality` release state and forbidden copy scan.

## Risks

- Risk: a resolution is misread as data quality becoming clean. Mitigation: keep P66 policy visible and set `clean_data_claim_allowed=false` unless policy passes.
- Risk: storing sensitive free text. Mitigation: sanitize before persistence and before display; tests include secret/path/SQL/prompt samples.
- Risk: UI suggests automation. Mitigation: labels use "记录" / "撤销记录" and avoid refresh/repair/confirm/apply/trade language.
