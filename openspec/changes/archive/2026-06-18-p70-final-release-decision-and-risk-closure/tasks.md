# Tasks: P70 Final Release Decision And Risk Closure

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P68/P69 change 格式、P63-P69 release/acceptance/handoff/package 材料。
- [x] 1.3 创建 `p70-final-release-decision-and-risk-closure` OpenSpec change。
- [x] 1.4 写明 P70 只做最终发布决策和风险收口，不新增运行时能力。
- [x] 1.5 更新当前进度文档，标记 P70 active。
- [x] 1.6 运行 `openspec validate p70-final-release-decision-and-risk-closure --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 最终发布证据审计

- [x] 2.1 审计 P63-P69 release/acceptance 材料，确认产品/UI、package、repeat acceptance、current-data policy 和 scope exclusion 的事实链一致。
- [x] 2.2 审计 release candidate、handoff、release README、repeatability 是否仍存在 P68 之后的 stale wording。
- [x] 2.3 审计 P69 package wording，确认只声明覆盖 P68 source commit，不声称包含 P69/P70 文档。
- [x] 2.4 审计 current-data wording，确认不把 P67 scope exclusion 描述为 P66 pass 或 current data clean/healthy。
- [x] 2.5 审计 optional future stages，确认不把物理第二机器复验、P66 true pass、post-P69/P70 package refresh 或 VecLite hardening 描述为当前 limited release 必需项。

## 3. P70 决策记录与材料刷新

- [x] 3.1 新增 `docs/release/acceptance/2026-06-18-p70-final-release-decision.md`，记录最终 release decision、证据矩阵、剩余风险、Not Claimed 和 optional future work。
- [x] 3.2 更新 `docs/release/release-handoff-2026-06-18.md`：引用 P70，并把 Next Stage 改为“无必需下一阶段；仅保留可选后续”。
- [x] 3.3 更新 `docs/release/README.md` 和 `docs/release/acceptance-repeatability.md`：增加 P70 最终决策入口和复验边界。
- [x] 3.4 如 `docs/release/release-candidate-2026-06-18.md` 仍建议 P69，改为引用 P69/P70 结果。
- [x] 3.5 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。

## 4. 验证与安全扫描

- [x] 4.1 顺序运行 P67 resolution check，记录 `claim_state`、`policy`、`gate`、`resolution` 和 `clean_data_claim`，不得与 P66 strict gate 并行以避免 SQLite lock。
- [x] 4.2 顺序运行 P66 strict current-data gate，记录 expected non-zero 与 `policy=blocked` / `gate=block`，不得与 P67 resolution check 并行以避免 SQLite lock。
- [x] 4.3 运行 release 文案扫描，确认无 current data clean/healthy、P66 pass、scope exclusion as policy pass、P69 archive includes P69/P70 docs 的误导性声明。
- [x] 4.4 运行禁止能力扫描，确认无新增远程发布、Git tag、自动升级、自动迁移、自动恢复、自动修复、券商接口、交易、外推、自动确认、自动规则应用、收益承诺、登录源、付费源、授权源、Level2 或高频源承诺。
- [x] 4.5 运行 `openspec validate p70-final-release-decision-and-risk-closure --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 4.6 若 P70 意外修改运行时代码、scripts 或前端，追加 `go test ./...`、`npm --prefix web test`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh`。

## 5. 复审、归档与提交

- [x] 5.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 5.2 执行 OpenSpec archive，把 P70 delta 合并到 docs/OpenSpec specs。
- [x] 5.3 archive 后确认无活跃 change，并确认是否还有必需下一阶段。
- [x] 5.4 提交前子 agent 复审无 Critical / Important。
- [x] 5.5 提交 P70。
