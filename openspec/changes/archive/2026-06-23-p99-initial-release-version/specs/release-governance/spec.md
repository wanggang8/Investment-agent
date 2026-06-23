## ADDED Requirements

### Requirement: Initial Release Version Marker

The repository SHALL declare a single initial release version marker for release-facing handoff materials.

#### Scenario: Initial version is recorded

- **GIVEN** the project is prepared for a local release handoff
- **WHEN** a user or packaging operator inspects the repository version marker
- **THEN** the root `VERSION` file SHALL contain `v0.1.0`
- **AND** release materials SHALL describe `v0.1.0` as an initial local release version, not as proof of a Git tag, remote release publication, physical second-machine validation, or expanded runtime investment capability.
