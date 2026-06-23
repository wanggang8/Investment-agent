## ADDED Requirements

### Requirement: Post-P67 release readiness decision

After P67, release materials SHALL include a P68 release readiness decision before any new final release handoff claim is made.

#### Scenario: P66 remains blocked and P67 scope exclusion is active

- **GIVEN** the P66 strict current-data gate reports `policy=blocked` and `gate=block`
- **AND** the P67 resolution check reports `claim_state=resolved_with_scope_exclusion`
- **WHEN** release readiness is described
- **THEN** the release status SHALL explicitly exclude current local data health from clean-data claims
- **AND** the materials SHALL NOT describe current local data as clean or healthy
- **AND** the materials SHALL NOT describe the P66 policy as passed

#### Scenario: Final distribution package evidence is stale after later commits

- **GIVEN** release package repeat evidence was generated before later P66, P67, or P68 commits
- **WHEN** final distribution readiness is described
- **THEN** the materials SHALL either require a package refresh stage or explicitly limit the package evidence to the earlier candidate archive
- **AND** the materials SHALL NOT imply that the earlier package artifact includes later commits.

#### Scenario: No active P67 resolution exists

- **GIVEN** the P66 strict current-data gate reports `gate=block` or `gate=waiver_required`
- **AND** the P67 resolution check reports `claim_state=requires_resolution`
- **WHEN** release readiness is described
- **THEN** the result SHALL be release-blocking for any claim depending on current local data health
- **AND** the materials SHALL direct the operator to record a valid waiver or scope exclusion before making a limited release-ready claim.
