## ADDED Requirements

### Requirement: P85 expected return analysis accuracy closure

After P84, expected-return and analysis-accuracy rows SHALL NOT be marked `real_pass` unless fresh real local execution proves deterministic calculation correctness, provenance, degradation safety, LLM boundary safety, and user-visible readback.

#### Scenario: P85 row inventory is complete before execution

- **GIVEN** P85 starts from the P84 evidence matrix
- **WHEN** execution begins
- **THEN** the P85 plan SHALL enumerate exactly 31 expected-return and analysis-accuracy rows
- **AND** each row SHALL map to a concrete acceptance mode and evidence target.

#### Scenario: Expected-return calculations are deterministic checks

- **GIVEN** P85 evaluates expected-return or scenario fields
- **WHEN** a row is marked `real_pass`
- **THEN** the product output SHALL be compared with independently computed deterministic expectations where the value is deterministic
- **AND** future return, future market direction, or investment performance accuracy SHALL NOT be claimed.

#### Scenario: LLM remains analysis-only

- **GIVEN** a P85 scenario includes real LLM output, unavailable LLM, or LLM quality failure
- **WHEN** the decision workflow completes or degrades
- **THEN** LLM material SHALL remain analysis-only
- **AND** it SHALL NOT override final rule verdict, create confirmations, trigger trades, or suppress required data-quality blockers.
