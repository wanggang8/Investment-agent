# local-knowledge-import Specification

## Purpose
TBD - created by archiving change p46-local-knowledge-import-governance. Update Purpose after archive.
## Requirements
### Requirement: Local knowledge imports SHALL be validated before facts are written

The system SHALL provide a local-only knowledge import validation flow for user-provided notes, research excerpts, and CSV/table-like rows before any local evidence or RAG facts are written.

#### Scenario: User validates local knowledge rows

- **WHEN** the user submits local knowledge rows for validation
- **THEN** the system SHALL return row-level validation status, redacted previews, content hashes, risk flags, blocking count, warning count, and an index plan with `rag_chunk_count` and `index_status`
- **AND** it SHALL generate `import_batch_id` from normalized `source_label`, `default_symbol`, row order, and row content hashes
- **AND** validation SHALL NOT write `intelligence_items`, `intelligence_summary`, `rag_chunks`, `source_verifications`, `audit_events`, portfolio facts, rule versions, or decision records
- **AND** the response SHALL NOT expose complete secrets, private local paths, raw SQL, full prompts, or raw HTTP responses

#### Scenario: Validation detects unsafe local knowledge content

- **WHEN** a row contains suspected API keys, private key material, raw SQL, private filesystem paths, raw HTTP responses, or full prompts
- **THEN** the row SHALL be marked as blocking or warning according to severity
- **AND** the preview SHALL redact the sensitive portion
- **AND** raw HTTP responses or full prompts SHALL be replaced in preview rather than partially displayed
- **AND** the row SHALL NOT be eligible for confirmation while blocking issues remain

### Requirement: Local knowledge import confirmation SHALL require explicit user confirmation

The system SHALL write local knowledge facts only after explicit user confirmation and only when server-side validation passes.

#### Scenario: User confirms a valid local knowledge import

- **GIVEN** the submitted rows pass validation with no blocking issues
- **WHEN** the user confirms the import with an `import_batch_id`, `confirm_reason`, `source_label`, `default_symbol`, and rows
- **THEN** the system SHALL re-run validation on the submitted rows
- **AND** it SHALL recompute `import_batch_id` and reject the request when the recomputed value does not match the submitted `import_batch_id`
- **AND** it SHALL write local facts in a single transaction
- **AND** it SHALL write `intelligence_items`, `intelligence_summary`, `rag_chunks`, `source_verifications`, and `audit_events`
- **AND** imported local knowledge SHALL default to `source_level=C`, `evidence_role=background`, and `rag_chunks.index_status=pending`
- **AND** P46 imported `intelligence_items.original_url` SHALL remain empty even if the request included `source_url`
- **AND** the response SHALL include write counts, audit ids, and `index_status`

#### Scenario: User confirms rows with blocking validation issues

- **GIVEN** the submitted rows contain blocking validation issues
- **WHEN** the user attempts to confirm the import
- **THEN** the system SHALL reject the confirmation
- **AND** it SHALL NOT write local evidence facts, RAG chunks, source verifications, portfolio facts, rule versions, decision records, or audit success events

### Requirement: Local knowledge imports SHALL preserve investment safety boundaries

The system SHALL keep local knowledge import materials as user-provided background material and SHALL NOT expand trading, external delivery, or rule-application capabilities.

#### Scenario: Local knowledge import is used by downstream retrieval

- **WHEN** imported local knowledge is later visible in evidence or retrieval surfaces
- **THEN** it SHALL remain distinguishable as local user-provided background material
- **AND** it SHALL NOT count as A/S high-grade independent formal evidence
- **AND** it SHALL NOT by itself satisfy major-event multi-source verification requirements

#### Scenario: Local knowledge import controls are rendered in the frontend

- **WHEN** the user opens the local knowledge import page
- **THEN** the page SHALL show validation, redacted preview, index plan, and explicit confirmation controls
- **AND** it SHALL NOT provide broker API, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair promise, return promise, login-only source, paid source, authorization-gated source, Level2 source, or high-frequency source capabilities

