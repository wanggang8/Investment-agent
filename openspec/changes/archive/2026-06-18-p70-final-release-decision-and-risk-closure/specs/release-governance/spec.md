## ADDED Requirements

### Requirement: Final release decision and risk closure

After P69 clean-tree package refresh, release materials SHALL include a final release decision and risk-closure record before the project is described as having no mandatory next phase.

#### Scenario: Limited local release evidence is sufficient

- **GIVEN** P63 full UI regression and P65 repeat acceptance evidence remain passing
- **AND** P67 reports an active scope exclusion for the blocked current-data gate
- **AND** P69 clean-tree package verify and repeat acceptance have passed
- **WHEN** the release handoff describes final milestone status
- **THEN** the status SHALL be `release_ready_limited_current_data_scope`
- **AND** the handoff SHALL state that no mandatory next phase remains for that limited scope
- **AND** optional future stages SHALL be separated from release blockers.

#### Scenario: Current data remains blocked

- **GIVEN** the P66 strict current-data gate reports `policy=blocked` and `gate=block`
- **WHEN** the final release decision is written
- **THEN** it SHALL NOT claim current local data is clean or healthy
- **AND** it SHALL NOT describe P67 `resolved_with_scope_exclusion` as a P66 policy pass
- **AND** it SHALL preserve the current-data limitation in the release status.

#### Scenario: Package evidence does not include later documentation

- **GIVEN** the P69 package source commit predates P69 and P70 documentation
- **WHEN** package evidence is described in final handoff material
- **THEN** the material SHALL state the exact covered source commit or phase boundary
- **AND** it SHALL NOT imply that the P69 package archive includes P69 or P70 documents.
