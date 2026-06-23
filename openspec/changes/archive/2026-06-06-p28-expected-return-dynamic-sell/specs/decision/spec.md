## ADDED Requirements

### Requirement: P28 expected return analysis SHALL expose sample context and precision limits

The system SHALL generate expected return analysis as explanatory material with explicit sample context, precision limits, and disclaimers.

#### Scenario: Expected return analysis has enough samples for probabilities

- **WHEN** expected return analysis is generated with at least 20 comparable samples
- **THEN** the response SHALL include upside, base, and downside scenarios
- **AND** it MAY include numeric probabilities
- **AND** it SHALL include `sample_count`, `sample_window`, `screening_condition`, `precision_status=available`, and a disclaimer that the analysis is not a return promise.

#### Scenario: Expected return analysis has insufficient samples

- **WHEN** expected return analysis is generated with 5 to 19 comparable samples
- **THEN** the response MAY include qualitative return ranges
- **AND** it SHALL NOT include precise numeric probabilities
- **AND** it SHALL include `precision_status=insufficient`, `sample_count`, `sample_window`, `screening_condition`, and a reason explaining the sample limitation.

#### Scenario: Expected return analysis has too few samples

- **WHEN** expected return analysis is generated with fewer than 5 comparable samples
- **THEN** the response SHALL use `precision_status=unavailable`
- **AND** it SHALL NOT include return ranges or precise probabilities
- **AND** it SHALL include `sample_count`, `sample_window`, `screening_condition`, a qualitative reason, and disclaimer
- **AND** unavailable sample window or screening condition SHALL be represented with an explicit missing reason rather than omitted silently.

### Requirement: P28 dynamic sell evaluation SHALL remain advisory and non-trading

The system SHALL translate expected return boundary conditions into advisory sell evaluation signals without changing account state or executing trades.

#### Scenario: Price reaches an expected return sell evaluation boundary

- **WHEN** current price or NAV enters the upside scenario lower bound, exceeds the base scenario upper bound, falls below the downside scenario lower bound, the base scenario midpoint shifts down by more than 15%, or an explicitly configured user target return is reached
- **THEN** the system SHALL produce a sell evaluation or reassessment prompt in expected return material or optional actions
- **AND** missing target return configuration SHALL be treated as not applicable rather than inferred
- **AND** prompts SHALL use advisory language such as evaluate, review, or record manual plan, not execute or place an order
- **AND** it SHALL NOT update positions, portfolio snapshots, confirmations, or transactions unless the user later records a manual action
- **AND** it SHALL NOT call broker, trading, order, cancellation, Webhook, Push, SMS, email, or WebSocket capabilities.

#### Scenario: Expected return conflicts with rule arbitration

- **WHEN** expected return analysis suggests a sell evaluation but rule arbitration returns a stricter status such as `insufficient_data`, `frozen_watch`, `sell_only`, or `rejected`
- **THEN** the final verdict SHALL continue to come from rule arbitration
- **AND** expected return analysis SHALL remain explanatory material only.
