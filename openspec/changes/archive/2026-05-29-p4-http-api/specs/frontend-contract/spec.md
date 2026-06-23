## ADDED Requirements

### Requirement: P4 frontend DTO compatibility

P4 API DTOs SHALL use field names and stable error codes that are compatible with `docs/frontend-contract.md`.

#### Scenario: Frontend receives stable display states

- **WHEN** the API returns `DATA_REQUIRED`
- **THEN** frontend clients SHALL be able to map it to `first_use`
- **WHEN** the API returns `DATA_STALE`, `EVIDENCE_NOT_FOUND`, `VECTOR_INDEX_UNAVAILABLE`, `ANALYST_UNAVAILABLE`, or `DECISION_RECORD_FAILED`
- **THEN** frontend clients SHALL be able to map it to `insufficient_data`
- **WHEN** the API returns `SOURCE_VERIFICATION_FAILED`
- **THEN** frontend clients SHALL be able to map it to `frozen_watch`

#### Scenario: DTO fields match page contracts

- **WHEN** dashboard, decision, evidence, portfolio, rule, audit, settings, market, or review APIs return data
- **THEN** the JSON fields SHALL match the corresponding API fields referenced by `docs/frontend-contract.md`
