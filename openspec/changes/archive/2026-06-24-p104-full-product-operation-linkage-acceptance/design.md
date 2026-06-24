# Design

## Acceptance Shape

P104 is a quality gate rather than a product feature. It combines two artifacts:

1. A human-readable matrix that lists operation surfaces and expected evidence.
2. A machine runner that validates representative linked operations end-to-end.

The runner validates the following evidence types for each covered operation class:

- HTTP request succeeds or rejects correctly.
- API readback reflects the expected product state.
- SQLite tables contain the expected durable facts.
- Downstream endpoints can read the same state.
- `audit_events` and manual confirmation rows are present where appropriate.
- Forbidden automation tables/actions are absent.

## Runner Strategy

The runner uses an isolated temporary config and SQLite path, then starts the local Go backend on a caller-selectable localhost port. It does not use the user's configured production database. It does not call Docker or installer scripts.

The runner uses API operations for product behavior and SQLite checks for data impact. It may seed non-user-facing setup records directly in SQLite when an operation depends on historical facts, such as a pending decision or a risk alert.

## Boundary

P104 does not replace P92/P93 original-requirement and code-reality ledgers. It adds a fresh, repeatable product linkage gate over the local source tree. It also does not claim exhaustive branch coverage for every invalid input; focused Go/frontend tests continue to cover narrow branch behavior.
