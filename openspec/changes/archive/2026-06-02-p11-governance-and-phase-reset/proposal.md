## Why

P10 已完成后，项目需要进入后续 P11-P18 拆分执行阶段。此前 archive 后曾出现活跃 change 残留、进度文件与治理说明不一致的问题，需要先清理治理状态，避免后续 change 在错误阶段或错误真源上实施。

## What Changes

- 清理 `openspec/changes/` 下非预期活跃 change 残留，只保留 `archive/` 与当前正在执行的 change。
- 更新治理与进度说明，使 P11-P18 的阶段编号、执行门槛、archive 前复审要求一致。
- 明确每个后续 change 必须独立完成 propose、apply、verify、subagent review、archive。
- 明确 L1 契约仍只能通过 OpenSpec delta 修改，archive 时再合并到 `docs/`。

## Capabilities

### New Capabilities
- `governance-phase-reset`: Covers phase reset, active-change hygiene, and mandatory pre-archive review workflow for P11-P18.

### Modified Capabilities
- `product-completeness`: Records that P10 completion now hands off to a governed P11-P18 roadmap rather than more P10 work.

## Impact

- Documentation and governance files: `openspec/PROGRESS.md`, `docs/development-plan.md`, `docs/GOVERNANCE.md`, `AGENTS.md` if needed.
- OpenSpec state: active changes must be cleaned or archived before each next change.
- No runtime behavior change is expected for backend or frontend.
