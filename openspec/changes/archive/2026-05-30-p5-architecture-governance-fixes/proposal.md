## Why

The architecture review found that the current P4 backend works functionally but still concentrates SQLite access, transaction handling, and concrete dependency wiring inside application-facing packages. This should be corrected before P5 frontend work grows the API surface and before P6 end-to-end hardening depends on stable audit and transaction boundaries.

## What Changes

- Introduce an architecture governance pass for application boundaries, transaction coordination, ID/Clock generation, and contract enum validation.
- Move SQLite-specific workflow dependency construction out of `internal/application/workflow` and into the composition root.
- Add a reusable transaction coordination pattern for multi-repository writes, so handlers do not manage SQL transactions directly.
- Centralize remaining contract enum validation, including review/error-case enums currently validated locally in handlers.
- Centralize request/entity ID and time generation behind injectable helpers for deterministic tests and consistent audit records.
- Define the P5 frontend feature-based organization before implementing cockpit pages.
- No new public API groups and no change to the no-auto-trading product boundary.

## Capabilities

### New Capabilities
- `architecture-governance`: Covers internal architecture requirements for dependency direction, transaction coordination, ID/Clock generation, enum validation, and frontend feature organization.

### Modified Capabilities

## Impact

- Affected backend packages:
  - `cmd/server`
  - `internal/application/handler`
  - `internal/application/workflow`
  - `internal/domain/model`
  - `internal/domain/repository`
  - `internal/infrastructure/persistence/sqlite`
  - `internal/pkg`
- Affected frontend packages:
  - `web/src/app`
  - `web/src/pages`
  - `web/src/services`
  - `web/src/types`
  - new `web/src/features/*` structure
- Affected docs at archive time:
  - `docs/architecture.md`
  - `docs/development-plan.md` if this change is completed before P5 implementation
