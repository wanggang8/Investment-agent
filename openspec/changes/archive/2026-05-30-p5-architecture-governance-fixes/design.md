## Context

The project already has strong product and contract boundaries: SQLite is the local fact base, VecLite is rebuildable auxiliary retrieval, DeepSeek does not make final verdicts, and every critical action writes audit facts. The current implementation satisfies P4 behavior but leaves some architectural concerns inside application-facing code:

- `handler.App` exposes `*sql.DB`, and several handlers own cross-table SQL transactions.
- `internal/application/workflow/dependencies.go` imports the SQLite implementation to build concrete repositories.
- Transaction handling exists both as `sqlite.withTx` and handwritten handler `BeginTx` blocks.
- Some contract enums and validation rules are local to handlers.
- ID and time creation are spread across request helpers, workflow builders, and handler logic.
- The frontend is still a skeleton and should adopt a feature organization before P5 pages expand.

The change keeps current product behavior and API semantics intact while improving internal maintainability and audit consistency.

## Goals / Non-Goals

**Goals:**

- Keep dependency direction aligned with `docs/architecture.md`: application depends on interfaces and domain rules, not concrete SQLite implementations.
- Make handlers thin: parse requests, call application use cases, return DTO envelopes and mapped errors.
- Provide one transaction coordination path for cross-repository writes.
- Make ID and Clock generation injectable and consistent across handlers, workflows, repositories, and tests.
- Move remaining contract enum validation into domain-level constants and validation methods.
- Define a P5 frontend feature layout that can grow without mixing page rendering, API calls, DTO mapping, and business components.

**Non-Goals:**

- No new external dependencies.
- No new public API groups.
- No changes to the no-auto-trading boundary.
- No rewrite of Eino graph semantics.
- No broad frontend implementation beyond directory and contract scaffolding needed for P5 organization.

## Decisions

### 1. Composition root owns concrete wiring

`cmd/server` or a small bootstrap package SHALL construct SQLite repositories and pass interfaces into application and workflow packages. `internal/application/workflow` SHALL no longer import `internal/infrastructure/persistence/sqlite`.

Alternative considered: keep `NewSQLiteWorkflowDependencies` in workflow for convenience. Rejected because it preserves the reverse dependency found in the architecture review.

### 2. Introduce a transaction coordinator

Add an application-facing transaction abstraction such as `Transactor` or `UnitOfWork` with a method like `WithinTx(ctx, func(ctx context.Context, repos Repositories) error) error`. The SQLite implementation may live under infrastructure and provide transactional repositories backed by `*sql.Tx`.

Alternative considered: keep handwritten handler transactions. Rejected because it duplicates rollback behavior and makes cross-table audit consistency harder to verify.

### 3. Use application services for write use cases

Move multi-table writes from handlers into use case services, for example confirmation, portfolio adjustment, and rule proposal confirmation services. The service validates DTO-derived command input, uses the transaction coordinator, writes facts and audit events, and returns response DTO data.

Alternative considered: keep SQL in helpers shared by handlers. Rejected because helpers would still couple HTTP and persistence responsibilities.

### 4. Centralize ID and Clock generation

Add small injectable interfaces in `internal/pkg` or application support packages. Production wiring uses UTC time and readable business IDs. Tests can inject fixed time and deterministic IDs.

Alternative considered: continue calling `time.Now()` and deriving IDs from request IDs. Rejected because it weakens deterministic testing and makes entity ID rules inconsistent with data model expectations.

### 5. Domain owns contract enums

Move remaining contract enum sets, including error-case root cause tags, into `internal/domain/model` with `Valid()` methods. Handler validation SHALL reuse these constants rather than local string switches.

Alternative considered: duplicate enums in DTO packages. Rejected because DTO duplication increases drift risk between API, data model, and frontend contract.

### 6. P5 frontend uses feature-based folders

Before P5 page implementation, introduce feature directories such as `web/src/features/dashboard`, `decision`, `evidence`, `rules`, `audit`, `settings`, and shared API/type utilities. Pages compose feature modules; generic layout remains under `web/src/app` and reusable primitives remain under shared/component folders.

Alternative considered: continue placing all pages and API functions in flat `pages` and `services`. Rejected because the P5 cockpit will add multiple pages, DTO mappers, state components, and error states.

## Risks / Trade-offs

- Application service extraction could touch many files → Keep behavior-preserving changes small and migrate one write use case at a time.
- Transaction coordinator can become over-abstracted → Define only the methods needed by current cross-table writes.
- Moving dependency construction may temporarily increase `cmd/server` size → Allow a small bootstrap package if `main.go` becomes noisy.
- Enum centralization may require test fixture updates → Update fixtures in the same task as enum migration.
- Frontend feature scaffolding may look ahead of P5 → Only create structure and shared conventions needed by upcoming pages, no speculative UI implementation.

## Migration Plan

1. Add ID, Clock, and transaction abstractions with tests.
2. Move SQLite dependency construction to `cmd/server` or bootstrap while preserving existing route registration.
3. Extract confirmation, portfolio, and rule write flows into application services using the transaction coordinator.
4. Centralize remaining enums and update handler tests.
5. Add frontend feature folder scaffold and move existing skeleton code without changing visible behavior.
6. Run focused Go tests, full Go tests, and frontend build.

Rollback strategy: each step is behavior-preserving and can be reverted independently. If transaction abstraction causes regressions, revert the service extraction while retaining enum and ID/Clock cleanup.

## Open Questions

- Whether the composition root should stay directly in `cmd/server/main.go` or use `internal/bootstrap` for readability.
- Whether transaction repositories should be exposed as a single aggregate interface or as named repository getters inside the transaction callback.
