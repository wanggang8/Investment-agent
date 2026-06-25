# release-governance Delta

## ADDED Requirements

### Requirement: P119 Full UI Control And Affordance Acceptance

The product SHALL provide repeatable evidence that every production route has been reviewed for visible controls, form alignment, productized UI state, backend consistency, and forbidden investment-action affordances before claiming the P119 UI/control acceptance layer.

#### Scenario: All visible route controls are inventoried and classified

- **WHEN** the P119 browser runner visits the production routes from `web/src/App.tsx`
- **THEN** it SHALL collect visible buttons, links, inputs, selects, textareas, details summaries, and navigation controls
- **AND** every collected interactive control SHALL have a non-empty accessible or visible name
- **AND** every control SHALL be classified as navigation, light interaction, read action, local fact write, governance confirmation, expected disabled, or boundary notice.

#### Scenario: Key visible write controls are backed by real local state

- **WHEN** P119 exercises safe-to-automate UI write operations
- **THEN** portfolio facts, offline transaction facts, decision confirmations, risk SOP lifecycle updates, notification read state, data-quality resolutions, rule proposals, local knowledge imports, evidence maintenance, or market refresh actions SHALL have API and SQLite readback evidence
- **AND** unsupported broker/order/auto paths SHALL remain absent.

#### Scenario: Productized visual layout and copy pass

- **WHEN** P119 captures desktop and mobile UI evidence
- **THEN** routes SHALL not be blank
- **AND** visible controls SHALL not overflow the viewport
- **AND** pages SHALL not expose raw debug, mock, placeholder, secret, stack trace, undefined, null, or NaN UI copy outside explicitly allowed product boundary language
- **AND** browser console errors, page errors, and API 5xx responses SHALL be zero.

#### Scenario: Upstream toggle interactions are exercised

- **WHEN** P119 classifies visible navigation toggles, details summaries, row expanders, filters, selects, and read-refresh controls as light or upstream interactions
- **THEN** the browser runner SHALL exercise a representative cross-route sweep of those controls with before/after assertions
- **AND** the final evidence SHALL include toggle interaction count and zero toggle issues.

#### Scenario: Safety boundary remains explicit

- **WHEN** P119 scans visible UI and SQLite schema/effects
- **THEN** the product SHALL not expose one-click trading, broker order placement, order delegation, external push execution, automatic confirmation, automatic rule application, or return-guarantee affordances
- **AND** any local-only operation SHALL continue to state or imply local fact recording rather than trade execution.
