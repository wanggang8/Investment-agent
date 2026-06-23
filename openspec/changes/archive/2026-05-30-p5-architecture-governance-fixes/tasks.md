## 1. Foundation

- [x] 1.1 Add shared ID generator and Clock interfaces with production implementations
- [x] 1.2 Add deterministic ID and Clock test helpers
- [x] 1.3 Add transaction coordinator interfaces for multi-repository write use cases
- [x] 1.4 Implement SQLite transaction coordinator with rollback and commit tests

## 2. Dependency Wiring

- [x] 2.1 Move SQLite repository construction out of `internal/application/workflow`
- [x] 2.2 Update workflow dependency construction to accept repository interfaces only
- [x] 2.3 Wire production dependencies in `cmd/server` or a small bootstrap package
- [x] 2.4 Add import-boundary tests or static checks covering workflow and domain packages

## 3. Application Write Services

- [x] 3.1 Extract confirmation write flow from handler into an application service
- [x] 3.2 Extract portfolio initialization and adjustment write flows into application services
- [x] 3.3 Extract rule proposal confirmation and final confirmation write flows into application services
- [x] 3.4 Ensure extracted services use the transaction coordinator for cross-table facts and audit events
- [x] 3.5 Keep HTTP handlers limited to request parsing, service calls, and response/error writing

## 4. Contract Enums and Validation

- [x] 4.1 Move error-case root cause tag enum into `internal/domain/model`
- [x] 4.2 Replace handler-local enum switches with domain `Valid()` methods
- [x] 4.3 Update DTO and handler tests to verify invalid enum values return stable `BAD_REQUEST` errors
- [x] 4.4 Verify API, data model, and frontend-contract enum values remain aligned

## 5. ID and Time Adoption

- [x] 5.1 Replace handler-local generated entity IDs with shared ID generator calls
- [x] 5.2 Replace direct `time.Now()` use in write paths with injected Clock
- [x] 5.3 Update workflow decision and evidence record builders to use shared ID and Clock helpers
- [x] 5.4 Add deterministic tests for generated audit IDs, decision IDs, and timestamps

## 6. Frontend Organization

- [x] 6.1 Create P5 feature folders for dashboard, decision, evidence, portfolio, rules, audit, settings, market, and review
- [x] 6.2 Move current dashboard skeleton into the dashboard feature without changing visible behavior
- [x] 6.3 Define shared frontend locations for API client, API types, DTO mappers, and reusable components
- [x] 6.4 Keep `web/src/pages` as route composition only for the migrated skeleton page

## 7. Verification

- [x] 7.1 Run focused backend tests for application handlers, workflow, domain model, infrastructure SQLite, and internal package helpers
- [x] 7.2 Run full backend test suite with `go test ./...`
- [x] 7.3 Run frontend build with `cd web && npm run build`
- [x] 7.4 Confirm OpenSpec tasks map to the `architecture-governance` spec scenarios
