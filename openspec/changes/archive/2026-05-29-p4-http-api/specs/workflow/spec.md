## ADDED Requirements

### Requirement: P4 workflow-backed API behavior

P4 APIs SHALL expose workflow-backed behavior without changing P3 workflow semantics.

#### Scenario: Consultation API runs the consultation workflow

- **WHEN** `POST /api/v1/decisions/consult` is called
- **THEN** the API SHALL synchronously execute the consultation workflow
- **AND** return a renderable decision detail response
- **AND** preserve the rule-first constraint that DeepSeek analysis cannot write the final verdict

#### Scenario: Evidence refresh keeps SQLite facts when vector indexing fails

- **WHEN** `POST /api/v1/evidence/refresh` writes SQLite facts successfully but vector index update fails
- **THEN** the SQLite facts SHALL NOT be rolled back
- **AND** the response SHALL communicate the index failure through stable API fields or error state

#### Scenario: Gatekeeper approval requires final confirmation

- **WHEN** a rule proposal passes gatekeeper audit
- **THEN** the API SHALL set or expose `pending_final_confirm`
- **AND** it SHALL NOT create a new active `rule_versions` entry until final confirmation succeeds
