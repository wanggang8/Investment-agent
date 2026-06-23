## Context

Existing P10 work introduced configured data source entries and stable unavailable states. The next step is a small real-data slice: one readonly market provider and one readonly intelligence provider that can be enabled by configuration, write validated facts to SQLite, and degrade cleanly when unavailable. SQLite remains the source of truth; external providers only supply input facts.

## Goals / Non-Goals

**Goals:**
- Provide a minimal readonly market provider that can produce `market_snapshots`.
- Provide a minimal readonly intelligence provider that can produce `intelligence_items` / summaries through existing ingestion paths.
- Keep stub mode available and deterministic.
- Make source failures observable through stable errors and audit events.
- Document which fields still use stub, empty state, or unavailable state.

**Non-Goals:**
- No broker API or automatic trading capability.
- No complete financial statement integration.
- No complete sentiment data integration.
- No realtime availability or latency SLA.
- No active stock recommendation feature.
- No guaranteed return language.

## Decisions

1. Use provider interfaces already present in application/workflow boundaries where possible.
   - Rationale: The current code already wires market and intelligence sources into workflow dependencies.
   - Alternative considered: Introduce a new external data platform module. Rejected for P12 because it would exceed the minimum viable scope.

2. Prefer simple HTTP/CSV/JSON readonly providers with fixture-backed tests.
   - Rationale: Provider parsing and degradation can be tested without depending on public network availability.
   - Alternative considered: Require live public provider access in tests. Rejected because local acceptance must remain offline-capable.

3. Treat provider failures as degraded facts, not fatal system corruption.
   - Rationale: The product boundary requires information-insufficient or source-unavailable states, not unsafe advice.
   - Alternative considered: Retry indefinitely. Rejected because it can block local workflows and hide source status.

## Risks / Trade-offs

- Public source formats can change → Mitigation: parser tests use fixtures and failures return stable errors.
- “Real data” scope can expand unexpectedly → Mitigation: P12 explicitly excludes full finance, full sentiment, realtime SLA, broker APIs, and trading actions.
- External network tests can be flaky → Mitigation: tests use local fixture clients or injected HTTP clients.
