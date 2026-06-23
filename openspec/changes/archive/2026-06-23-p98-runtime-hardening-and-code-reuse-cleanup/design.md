# P98 Design

## Scope

P98 is a narrow hardening change. It does not change investment behavior or the P92/P93 release claim. It closes two engineering risks: release-mode stub confusion and duplicated frontend redaction.

## Runtime Guardrail

Add a small runtime-mode concept to config, defaulting to `development` for example/local fallback compatibility. Docker release config and `.env.example` will set release mode. Config validation will reject `runtime.mode=release` combined with `data_sources.use_stub=true`.

This keeps tests and local examples simple while making the release path explicit. The existing `use_stub=false` Docker default remains the primary release behavior; runtime mode adds a second, clearer guardrail.

## Frontend Redaction Utility

Create `web/src/shared/utils/redaction.ts` with a shared `redactSensitiveText` function for key-shaped tokens, SQL fragments, prompt fragments, raw diagnostic payloads, stack traces, and local paths. Current call sites in `ErrorState`, `LocalInstallPage`, and `DataQualityPage` will delegate to this utility. Page-specific replacement wording can be handled through options instead of maintaining separate regex sets.

## Testing

Use TDD:

- Add config tests proving release+stub is invalid and release+non-stub remains valid.
- Add frontend utility tests proving the shared redactor covers existing sensitive diagnostic patterns.
- Update existing page/component tests to keep passing through the shared utility.

## Documentation

Update configuration/deployment docs only where needed to explain `runtime.mode`. Do not create new L1 capability claims.
