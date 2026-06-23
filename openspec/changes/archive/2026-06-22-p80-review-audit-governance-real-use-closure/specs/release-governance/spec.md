## ADDED Requirements

### Requirement: P80 Review Audit Governance Closure

After P79, any P80 claim that review, audit, error-case, rule-proposal, or gatekeeper-governance rows have moved to `real_pass` SHALL be backed by fresh real UI execution and field-level SQLite/readback evidence.

#### Scenario: P80 upgrades require review and audit field proof

- **GIVEN** a P80 row is upgraded to `real_pass`
- **WHEN** the P80 checker evaluates the row
- **THEN** the row SHALL have fresh browser evidence from the P80 review/audit/governance journey
- **AND** data-bearing rows SHALL include SQLite readback for the expected tables and fields
- **AND** audit rows SHALL include `action`, `node_action`, `actor`, `status`, `before_state`, `after_state`, and `request_id` when those fields are part of the row claim
- **AND** governance rows SHALL include rule proposal, gatekeeper audit, and audit-event references when those fields are part of the row claim
- **AND** the row SHALL remain non-`real_pass` if the evidence is count-only, screenshot-only, route-smoke-only, fixture-only without UI operation, mock-only, waiver-only, or only partially covers the row text.

#### Scenario: P80 broad monthly and final-application rows remain bounded

- **GIVEN** a row requires monthly attribution, full quarterly review, final rule application time, or every SOP/data-impact branch
- **WHEN** P80 does not prove the exact required fields through fresh UI and readback
- **THEN** that row SHALL remain non-`real_pass`
- **AND** P80 SHALL record the missing field-level evidence as the remaining gap.

#### Scenario: P80 release conclusion remains scoped

- **WHEN** P80 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.
