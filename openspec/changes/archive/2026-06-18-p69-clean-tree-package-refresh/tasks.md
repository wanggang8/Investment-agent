# Tasks: P69 Clean Tree Package Refresh

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P64 packaging、P65 repeat acceptance、P68 release readiness governance、package/repeat scripts。
- [x] 1.3 创建 `p69-clean-tree-package-refresh` OpenSpec change。
- [x] 1.4 写明 P69 只刷新 clean-tree package evidence 和 release materials，不新增运行时能力。
- [x] 1.5 明确 clean tree 包从已提交 P68 HEAD 生成；P69 文档是包后验收记录，不声称包内包含 P69 文档。
- [x] 1.6 更新当前进度文档，标记 P69 active。
- [x] 1.7 运行 `openspec validate p69-clean-tree-package-refresh --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.8 子 agent 方案复审无 Critical / Important 后执行。

## 2. Clean tree package generation

- [x] 2.1 记录 P68 source commit：`git rev-parse HEAD`。
- [x] 2.2 确认主工作树在创建 P69 change 前为 clean，并记录 P69 change 创建后主工作树不再用于 source_status evidence。
- [x] 2.3 创建 detached temporary worktree：`git worktree add --detach tmp/p69-clean-tree-source <p68_commit>`。
- [x] 2.4 在 temporary worktree 中确认 `git status --short` 为空。
- [x] 2.5 在 temporary worktree 中运行 `npm --prefix web ci`，为 package script 的 frontend build 准备依赖。
- [x] 2.6 重新确认 temporary worktree 的 `git status --short` 为空，确保 ignored `web/node_modules/` 不污染 source_status。
- [x] 2.7 在 temporary worktree 中运行 `bash scripts/local-release-package.sh --release-label p69-clean-tree --output-dir tmp/p69-release`。
- [x] 2.8 解析生成的 archive、manifest、SHA-256、source commit、source_status、archive size、entry count。
- [x] 2.9 在 temporary worktree 中运行 package verify：`bash scripts/local-release-package.sh --verify <archive> --output-dir tmp/p69-release`。

## 3. Package repeat acceptance

- [x] 3.1 在 temporary worktree 中运行 `bash scripts/local-release-repeat-acceptance.sh --archive <archive> --output-dir tmp/p69-repeat`。
- [x] 3.2 记录 repeat summary：commands、durations、status、source commit、source_status、skip flags。
- [x] 3.3 确认 repeat command matrix 至少覆盖 OpenSpec validation、Go tests、npm ci、frontend tests、frontend build、E2E smoke。
- [x] 3.4 确认 repeat artifacts 只位于 `tmp/`，不提交 archive、manifest、logs、node_modules、dist、SQLite DB、trace 或 private config。
- [x] 3.5 移除 temporary worktree：`git worktree remove tmp/p69-clean-tree-source`，或记录不能移除的原因。

## 4. 发布材料刷新

- [x] 4.1 新增 `docs/release/acceptance/2026-06-18-p69-clean-tree-package-refresh.md`，记录 package identity、verify result、repeat acceptance、source cleanliness 和 Not Claimed。
- [x] 4.2 更新 `docs/release/release-packaging-2026-06-18.md`：增加 P69 clean-tree package refresh 证据，保留 P64/P65 为历史候选证据。
- [x] 4.3 更新 `docs/release/release-handoff-2026-06-18.md`：将 package freshness 从 recommended 转为 P69 clean-tree package repeat passed。
- [x] 4.4 更新 `docs/release/README.md` 和 `docs/release/acceptance-repeatability.md`：增加 P69 repeat command 和状态。
- [x] 4.5 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 4.6 明确 release status 仍受 P68/P66/P67 限制：`release_ready_limited_current_data_scope`，不得声明 current data clean。

## 5. 验证与安全扫描

- [x] 5.1 运行 `openspec validate p69-clean-tree-package-refresh --strict`。
- [x] 5.2 运行 `openspec validate --all --strict`。
- [x] 5.3 运行 `git diff --check`。
- [x] 5.4 运行 release 文案扫描，确认无 current data clean/healthy、P66 pass、scope exclusion as policy pass、archive includes P69 的误导性声明。
- [x] 5.5 运行禁止能力扫描，确认无新增远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、券商接口、交易、外推、自动确认、自动规则应用、收益承诺、登录源、付费源、授权源、Level2 或高频源承诺。
- [x] 5.6 若 P69 意外修改运行时代码或 scripts，追加 `go test ./...`、`npm --prefix web test`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh`。

## 6. 复审、归档与提交

- [x] 6.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 6.2 执行 OpenSpec archive，把 P69 delta 合并到 docs/OpenSpec specs。
- [x] 6.3 archive 后确认无活跃 change，并规划下一阶段。
- [x] 6.4 提交前子 agent 复审无 Critical / Important。
- [x] 6.5 提交 P69。
