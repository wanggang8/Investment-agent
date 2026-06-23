# P98 Runtime Hardening And Code Reuse Cleanup

## Why

The final requirements and code-reality audits show the project is release-ready for the local/GitHub-Docker scope, but the latest review found small maintainability and operator-safety risks worth closing in one focused hardening pass:

- Development fallback paths are intentionally available, but the release path should make it harder to confuse stub/static fallback with real-data operation.
- Frontend redaction and safe diagnostic display logic is duplicated across multiple pages and components.
- A few large frontend/backend files remain acceptable but should not grow further without extracting shared utilities around the code touched by this change.

## What Changes

- Add explicit runtime-mode/profile validation so release-like configurations cannot accidentally run with local stub data.
- Keep development/test fallback support, but make release defaults and diagnostics explicit.
- Add a shared frontend redaction utility and migrate current duplicated redaction call sites to it.
- Add focused tests for runtime-mode validation and redaction behavior.
- Update code-reality/architecture documentation only where needed to describe the new guardrail.

## Out Of Scope

- Re-auditing or rewriting the P92 341-row requirement ledger.
- New investment strategy, new data provider breadth, new product workflows, or new HTTP business APIs.
- Broker integration, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
- Large-scale backend, frontend, CSS, Eino, SQLite, React, or Vite rewrite.

## Acceptance

P98 is acceptable only if:

- Release/runtime config validation rejects stub mode when the runtime mode is release.
- Docker release configuration remains non-stub by default.
- Development/test stub fallback remains available for existing tests and local examples.
- Frontend diagnostic/key/path/SQL/prompt/raw redaction is provided by a shared utility and current call sites use it.
- Existing backend, frontend, OpenSpec, P91/P92/P93, and API route checks remain passing.
