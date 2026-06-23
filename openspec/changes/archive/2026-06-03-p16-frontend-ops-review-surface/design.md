## Context

Backend and local-task capabilities now expose data-source behavior, index health, review summaries, audit events, and evidence-quality metadata. P16 is a frontend-focused pass that makes these operational facts visible and testable without changing safety boundaries.

## Goals / Non-Goals

**Goals:**
- Display data source and index health states in a clear operator-facing panel.
- Display periodic review summaries and rule suggestions using API/service DTOs.
- Provide visible tracking entrypoints for audit events, rule proposals, error cases, and related decisions.
- Cover empty, failed, degraded, and successful states in frontend tests.

**Non-Goals:**
- No automatic trading controls.
- No automatic rule application controls.
- No direct access from frontend to SQLite, VecLite, or local files.
- No broad visual redesign unrelated to ops/review surfaces.

## Decisions

1. Use existing DTOs first.
   - Prefer mapping current review, evidence, audit, and status DTOs in frontend components.
   - Add DTO fields only when required facts are otherwise impossible to render.

2. Keep states explicit.
   - Distinguish empty, failed, degraded, and successful states with user-facing labels.
   - Unknown enum values continue to display safe fallback text.

3. Tracking links are navigational only.
   - Links may point to audit/review/rule pages or filtered views.
   - They cannot trigger rule application, portfolio mutation, or order placement.

## Risks / Trade-offs

- Frontend-only state labels can drift from backend enums; tests should cover known and unknown mappings.
- Review summary data may be empty in local development; empty-state copy must not imply failure.
- P17 may add scheduler details later; P16 should not prebuild scheduling controls.
