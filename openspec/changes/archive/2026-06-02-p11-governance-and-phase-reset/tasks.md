## 1. Governance state cleanup

- [x] 1.1 Verify `openspec list --json` reports only the intended active change.
- [x] 1.2 Inspect `openspec/changes/` for non-archive residual directories and resolve unexpected entries.
- [x] 1.3 Confirm archived P10 artifacts remain in `openspec/changes/archive/` and are not edited.

## 2. Phase tracking alignment

- [x] 2.1 Update `openspec/PROGRESS.md` to add P11-P18 rows and mark P11 in progress.
- [x] 2.2 Update `docs/development-plan.md` with the P11-P18 roadmap summary and execution gates.
- [x] 2.3 Update `docs/GOVERNANCE.md` active change table and workflow notes for pre-archive readonly subagent review.
- [x] 2.4 Update `AGENTS.md` current change guidance if it is stale.

## 3. Validation

- [x] 3.1 Run `openspec validate p11-governance-and-phase-reset --strict`.
- [x] 3.2 Run `openspec validate --all --strict`.
- [x] 3.3 Run `go test ./...` to confirm governance-only changes did not affect backend behavior.
- [x] 3.4 Run `cd web && npm run test && npm run build` to confirm frontend remains unchanged and valid.
