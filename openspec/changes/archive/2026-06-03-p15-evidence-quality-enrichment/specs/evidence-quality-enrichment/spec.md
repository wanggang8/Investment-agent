## ADDED Requirements

### Requirement: Evidence quality metadata
The system SHALL preserve evidence quality metadata across normalization, persistence, retrieval, decision references, and DTOs.

#### Scenario: Quality metadata is preserved
- **WHEN** evidence is normalized or retrieved
- **THEN** source level, evidence role, published time, captured time, content hash, time weight, relevance score, independent source count, and high-grade independent source count SHALL remain available when known
- **THEN** the system SHALL NOT replace missing quality metadata with misleading placeholder values

#### Scenario: Decision evidence references retain quality fields
- **WHEN** a decision record persists evidence refs
- **THEN** evidence refs SHALL include source level, evidence role, published time, captured time, original URL, summary, content hash, time weight, relevance score, and high-grade independent source count

### Requirement: Formal/background boundary
The system SHALL enforce formal and background evidence boundaries before final rule arbitration.

#### Scenario: C-level source is background-only
- **WHEN** a source has level C
- **THEN** it SHALL be stored or returned as `background`
- **THEN** it SHALL NOT be used as formal decision evidence

#### Scenario: Unverified evidence is background-only
- **WHEN** source verification is not satisfied
- **THEN** evidence SHALL remain background-only or produce a non-satisfied verification state
- **THEN** final rule arbitration SHALL NOT treat it as satisfied formal evidence

### Requirement: High-grade independent source threshold
The system SHALL require at least two S/A independent sources before major-event evidence can be satisfied.

#### Scenario: Single high-grade source is insufficient
- **WHEN** a major positive, major negative, or buy-logic-break event has fewer than two S/A independent sources
- **THEN** verification SHALL NOT be `satisfied`
- **THEN** downstream rule arbitration SHALL not treat the event as satisfied formal evidence

#### Scenario: Two high-grade independent sources can satisfy
- **WHEN** a major event is supported by at least two S/A independent sources
- **THEN** verification MAY become `satisfied` if all other evidence checks pass

### Requirement: Structured facts and analyst materials are separated
The system SHALL distinguish structured facts used by rules from analyst materials used for explanation.

#### Scenario: Structured facts drive rule arbitration
- **WHEN** rule arbitration evaluates evidence
- **THEN** it SHALL use structured facts such as market snapshots, evidence sets, source verification status, and rule version context
- **THEN** it SHALL NOT let analyst text override the final verdict

#### Scenario: Analyst materials remain non-decisive
- **WHEN** analyst reports or expected return materials are present
- **THEN** they SHALL remain explanatory materials
- **THEN** final verdict generation SHALL remain rule-first
