## ADDED Requirements

### Requirement: Rule proposals SHALL have effect validation facts
The system SHALL generate and persist local effect validation facts for rule proposals so the user can understand source evidence, sample representativeness, overfit risk, historical replay outcome, and safety notes before any rule can be applied.

#### Scenario: Proposal validation is generated
- **WHEN** a rule proposal is evaluated for effectiveness
- **THEN** the system SHALL persist proposal ID, candidate rule version, source explanation, sample count, sample window, impacted workflows, representativeness status, overfit risk, replay result, metrics snapshot, risk notes, and timestamps
- **AND** it SHALL write an audit event for the validation.

#### Scenario: Validation has insufficient samples
- **WHEN** the available sample count is below the configured minimum or source facts are missing
- **THEN** the validation status SHALL be insufficient or needs_more_samples
- **AND** the proposal SHALL NOT be allowed to progress to final application based on that validation.

#### Scenario: Validation remains local-only
- **WHEN** effect validation is generated or refreshed
- **THEN** the system SHALL NOT create orders, broker state, operation confirmations, position transactions, external notifications, or active rule versions.

### Requirement: Rule effect validation SHALL check overfit and replay risk
The system SHALL evaluate each candidate rule change for overfit risk and historical replay impact using local decisions, confirmations, error cases, review summaries, and risk alerts.

#### Scenario: Overfit risk is high
- **WHEN** validation detects narrow samples, single-case tuning, contradictory outcomes, or risk alerts that would worsen under the candidate rule
- **THEN** the validation SHALL mark overfit risk as high or medium with reasons
- **AND** the guardrail decision SHALL be rejected or needs_user_review.

#### Scenario: Historical replay is unfavorable
- **WHEN** replay metrics show worse misjudgment rate, missing evidence rate, degradation rate, or risk alert exposure than the baseline
- **THEN** the validation SHALL record the unfavorable metrics
- **AND** the proposal SHALL NOT be treated as ready for final application.

### Requirement: Applied rules SHALL have post-application tracking
The system SHALL track applied rule versions after user final confirmation so review surfaces can show whether rule effectiveness improves or degrades over time.

#### Scenario: Applied rule is tracked in review
- **WHEN** a review summary is generated after a proposal has been applied
- **THEN** the system SHALL summarize rule hit count, misjudgment count, missing evidence count, degradation count, related risk alert count, and trend direction for that applied rule version.

#### Scenario: Tracking finds deterioration
- **WHEN** post-application tracking detects worsening metrics or recurring risk alerts
- **THEN** the system SHALL expose a review warning or draft follow-up suggestion
- **AND** it SHALL NOT automatically change or revert the active rule version.

### Requirement: Rule effect validation API SHALL be local and traceable
The system SHALL expose local APIs for rule effect validation and tracking so frontend surfaces can display validation facts without changing rules automatically.

#### Scenario: Frontend requests proposal validation
- **WHEN** the frontend requests validation for a proposal
- **THEN** the API SHALL return validation status, source explanation, sample summary, overfit risk, replay metrics, guardrail decision, related risk alerts, audit links, safety note, and timestamps.

#### Scenario: Frontend requests applied tracking
- **WHEN** the frontend requests tracking for an applied rule version or review period
- **THEN** the API SHALL return trend metrics and related proposal, rule version, review, audit, and risk alert links.

#### Scenario: Validation API has no result
- **WHEN** no validation result exists for the requested proposal or rule version
- **THEN** the API SHALL return a typed empty state or not_found response without fabricating metrics.

### Requirement: Rule effect validation UI SHALL remain safe
The frontend SHALL display rule validation and applied tracking as rule governance information without presenting automatic rule application, trading, or return-prediction controls.

#### Scenario: Proposal detail shows validation
- **WHEN** a proposal has validation facts
- **THEN** the frontend SHALL show status, sample summary, overfit risk, replay result, guardrail decision, risk notes, and local trace links
- **AND** it SHALL state that final application still requires gatekeeper audit and user final confirmation.

#### Scenario: Review page shows tracking
- **WHEN** review output contains applied rule tracking
- **THEN** the frontend SHALL show trend summaries and links to proposals, audits, and risk alerts
- **AND** it SHALL NOT provide automatic apply, automatic rollback, broker, or external notification controls.

### Requirement: Rule effect validation contract SHALL be documented for archive
The P36 delta SHALL record API, data model, workflow, and frontend contract changes for merge into the L1 docs during archive.

#### Scenario: P36 API docs are archived
- **WHEN** P36 is archived
- **THEN** `docs/api.md` SHALL include rule effect validation and applied tracking API contracts, DTO fields, errors, and non-trading transaction boundaries.

#### Scenario: P36 data model docs are archived
- **WHEN** P36 is archived
- **THEN** `docs/data-model.md` SHALL include validation/tracking facts, statuses, indexes, audit relationships, and constraints that prevent automatic rule application.

#### Scenario: P36 workflow and frontend docs are archived
- **WHEN** P36 is archived
- **THEN** `docs/workflow.md` and `docs/frontend-contract.md` SHALL describe validation generation, review integration, guardrail display, and safe UI behavior.