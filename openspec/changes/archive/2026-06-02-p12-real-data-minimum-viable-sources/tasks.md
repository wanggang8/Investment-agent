## 1. Provider contracts and tests

- [x] 1.1 Add tests for readonly market provider success, timeout/unavailable, stale data, and parse failure.
- [x] 1.2 Add tests for readonly intelligence provider success, timeout/unavailable, and parse failure.
- [x] 1.3 Add tests proving stub mode remains deterministic and does not require public network access.

## 2. Provider implementation

- [x] 2.1 Implement one configurable readonly market data provider with injected client/fixture support.
- [x] 2.2 Implement one configurable readonly intelligence provider with injected client/fixture support.
- [x] 2.3 Preserve source metadata: provider name, URL when available, published time, captured time, and source level.
- [x] 2.4 Return stable degraded errors for timeout, unavailable source, stale data, and parse failure.

## 3. Wiring and audit

- [x] 3.1 Wire provider selection from configuration while preserving `data_sources.use_stub`.
- [x] 3.2 Ensure successful provider refresh writes SQLite facts through existing repositories/workflows.
- [x] 3.3 Ensure provider failures write audit events and do not write invalid facts.
- [x] 3.4 Document which fields remain stub, empty, or unavailable in P12.

## 4. Safety boundaries

- [x] 4.1 Verify P12 introduces no broker API, automatic trading, active recommendation, or return guarantee.
- [x] 4.2 Verify account state is not mutated by provider refresh except existing user-recorded offline confirmation paths.

## 5. Validation

- [x] 5.1 Run `openspec validate p12-real-data-minimum-viable-sources --strict`.
- [x] 5.2 Run `openspec validate --all --strict`.
- [x] 5.3 Run `go test ./...`.
- [x] 5.4 Run `cd web && npm run test && npm run build`.
