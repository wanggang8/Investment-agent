## ADDED Requirements

### Requirement: P83 governance traceability backfill

After P80, governance and release traceability rows SHALL NOT be upgraded unless each row has exact artifact links, fresh validation where needed, and an honest status classification.

#### Scenario: P83 row inventory is complete before execution

- **GIVEN** P83 starts from the latest P82 evidence matrix
- **WHEN** execution begins
- **THEN** the P83 plan SHALL enumerate exactly 43 governance traceability rows
- **AND** each row SHALL have a target evidence or classification path.

#### Scenario: Evidence links are concrete

- **GIVEN** P83 marks a row as upgraded
- **WHEN** the evidence layer is reviewed
- **THEN** it SHALL include exact files, commands, tests, package manifests, acceptance records, UI/API evidence, or safety scans
- **AND** narrative-only assertions SHALL NOT be sufficient.

#### Scenario: Historical gaps remain honest

- **GIVEN** a historical archive or physical repeat acceptance was never performed
- **WHEN** P83 writes governance materials
- **THEN** it SHALL preserve that limitation
- **AND** it SHALL NOT fabricate historical archives, physical second-machine evidence, package refreshes, remote release, or Git tag evidence.
