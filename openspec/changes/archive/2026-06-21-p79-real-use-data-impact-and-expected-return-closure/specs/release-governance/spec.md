## ADDED Requirements

### Requirement: P79 Real-Use Data-Impact Closure

After P78, any P79 claim that portfolio, confirmation, local-account, or expected-return rows have moved to `real_pass` SHALL be backed by fresh real UI execution and SQLite/readback evidence.

#### Scenario: P79 upgrades require action-to-data proof

- **GIVEN** a P79 row is upgraded to `real_pass`
- **WHEN** the P79 checker evaluates the row
- **THEN** the row SHALL have fresh P79 evidence from a real browser journey or direct runtime readback
- **AND** data-bearing rows SHALL include SQLite readback for expected changed tables
- **AND** local-account rows SHALL include field-level readback for the relevant position, confirmation, transaction, evidence, and audit fields rather than table counts alone
- **AND** data-bearing rows SHALL include negative evidence that prohibited broker/order/external-push/automatic-confirmation tables or claims were not created
- **AND** the row SHALL remain non-`real_pass` if the evidence is inherited-only, screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, or incompatible single-action-only.

#### Scenario: P79 expected-return rows remain bounded

- **GIVEN** a P79 expected-return row requires probabilities, scenario ranges, sell-evaluation triggers, valuation fields, sample counts, sample windows, screening conditions, source/provenance fields, or non-trading disclaimers
- **WHEN** fresh P79 evidence lacks any required field
- **THEN** that row SHALL remain non-`real_pass`
- **AND** P79 SHALL record the missing field-level evidence as the remaining gap.

#### Scenario: P79 release conclusion remains scoped

- **WHEN** P79 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.

#### Scenario: Expected-return quality failure uses safe local material

- **GIVEN** the expected-return LLM material is parseable but fails the analyst safety quality gate
- **WHEN** deterministic local expected-return scenarios have been generated
- **THEN** the failed LLM material SHALL be discarded
- **AND** ExpectedReturnNode SHALL emit safe deterministic local expected-return material with metadata showing `model=deterministic-local`, `parse_status=parsed`, `quality_status=passed`, and `fallback_reason=llm_quality_failure`
- **AND** ordinary analyst timeout, authentication, or model-unavailable errors SHALL continue to degrade the analyst node
- **AND** this fallback SHALL NOT be used to upgrade expected-return probability or scenario rows without separate field-level UI/readback evidence.
