## ADDED Requirements

### Requirement: Governed post-P10 phase progression
The project SHALL track P11-P18 as post-P10 changes with a single active change at a time, and each change SHALL complete propose, apply, verify, readonly subagent review, and archive before the next change starts implementation.

#### Scenario: Starting the next post-P10 change
- **WHEN** an agent starts a P11-P18 implementation step
- **THEN** there SHALL be exactly one intended active OpenSpec change for that step
- **THEN** the agent SHALL use that change's `tasks.md` as the implementation checklist

#### Scenario: Archiving a post-P10 change
- **WHEN** a P11-P18 change is ready for archive
- **THEN** verification commands SHALL have completed for the change's relevant scope
- **THEN** a readonly subagent review SHALL report no Critical or Important issues before archive

### Requirement: Active change hygiene
The project SHALL treat unexpected non-archive entries under `openspec/changes/` as governance drift that must be resolved before proposing or implementing later post-P10 changes.

#### Scenario: Unexpected active change exists
- **WHEN** `openspec list --json` reports a change that is not the current intended change
- **THEN** the agent SHALL stop later implementation work until the active change is archived, removed, or explicitly selected as the current work

#### Scenario: No active change exists after archive
- **WHEN** a change is archived
- **THEN** `openspec list --json` SHALL report no non-expected active changes before the next change is proposed

## MODIFIED Requirements

### Requirement: Product completeness handoff
The project SHALL treat P10 product completeness as complete, and any remaining product-grade work SHALL be executed through the P11-P18 roadmap rather than by reopening P10.

#### Scenario: Post-P10 work is identified
- **WHEN** remaining work is found after P10 archive
- **THEN** it SHALL be assigned to a new P11-P18 change or a later explicitly proposed change
- **THEN** it SHALL NOT be implemented by editing archived P10 artifacts
