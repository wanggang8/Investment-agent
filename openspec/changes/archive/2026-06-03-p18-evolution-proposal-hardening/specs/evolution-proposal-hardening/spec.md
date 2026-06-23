## ADDED Requirements

### Requirement: Review-generated proposal chain
The system SHALL convert eligible review-derived rule suggestions into durable rule proposal records without applying rule changes automatically.

#### Scenario: Review suggestion creates proposal record
- **WHEN** a review evaluation identifies an eligible rule improvement suggestion
- **THEN** the system SHALL create or expose a rule proposal record with source review metadata
- **THEN** the proposal SHALL enter an existing non-applied review status such as draft or pending user confirmation

#### Scenario: Proposal source is traceable
- **WHEN** a review-generated proposal is created
- **THEN** the system SHALL record input references to the review period, supporting decisions, audit events, or error cases where available
- **THEN** audit events SHALL allow the proposal to be traced back to its review source

#### Scenario: Review proposal never applies rules automatically
- **WHEN** a review-generated proposal is created
- **THEN** no new active rule version SHALL be written by that proposal creation step
- **THEN** gatekeeper audit and user final confirmation requirements SHALL remain required before any rule version can become active

### Requirement: Evolution proposal safety gates
The system SHALL handle insufficient review evidence safely when evaluating rule suggestions.

#### Scenario: Insufficient sample does not create applied proposal
- **WHEN** review effectiveness data has insufficient samples or missing source facts
- **THEN** the system SHALL NOT create an applied proposal or active rule version
- **THEN** the output SHALL remain a safe review result, blocked proposal, or draft requiring more evidence
