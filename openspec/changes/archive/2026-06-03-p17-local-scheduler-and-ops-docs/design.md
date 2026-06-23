# Design

## Goals
- Make local periodic operation understandable and safe by default.
- Keep scheduling examples as operator-owned configuration, not an enabled runtime feature.
- Preserve audit visibility for each local task execution.
- Keep all trading and rule-change safeguards unchanged.

## Approach

### Scheduler examples
Add versioned examples under a documentation or examples location for:
- macOS launchd plist using `cmd/agent` tasks.
- cron-compatible snippets for local machines.

Examples must be inert by default: placeholders for repo path, config path, and schedule are acceptable; no example should include secrets, broker endpoints, order placement commands, or automatic confirmation flags.

### `cmd/agent` help clarity
Review `cmd/agent` help output and tests. If current help is incomplete, extend it to show safe task names, required local configuration, audit behavior, and the no-automatic-trading boundary.

### Audit expectations
Scheduled and manual task execution should use existing audit behavior. If tests show a task path lacks an audit event, add the narrowest repository/service wiring needed.

### Operations documentation
Update non-L1 operational docs, such as configuration/testing/delivery docs, to cover:
- Local startup checklist.
- SQLite backup/restore expectations.
- VecLite index rebuild and degraded recovery.
- Data source/DeepSeek failure handling.
- Scheduler setup and disable/removal procedure.
- Safety boundary: no trading, no order placement, no automatic rule application.

## Risks and mitigations
- **Risk:** Scheduler examples may be mistaken as enabled automation.  
  **Mitigation:** Keep examples disabled-by-default and explicitly require user installation.
- **Risk:** Local scheduled tasks could appear to mutate portfolio state.  
  **Mitigation:** Tests and docs must state that portfolio mutations still require recorded manual execution through existing confirmation paths.
- **Risk:** Docs include local secrets or user-specific paths.  
  **Mitigation:** Use placeholders only.
