# Tasks: P68 Post-P67 Release Readiness Governance

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P63 release candidate/handoff、P64 packaging、P65 repeat acceptance、P66 policy、P67 resolution 和 repeatability 材料。
- [x] 1.3 创建 `p68-post-p67-release-readiness-governance` OpenSpec change。
- [x] 1.4 写明 P68 只做发布状态治理、发布材料边界和下一阶段判断，不新增运行时能力。
- [x] 1.5 更新当前进度文档，标记 P68 active。
- [x] 1.6 运行 `openspec validate p68-post-p67-release-readiness-governance --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 发布材料一致性审计

- [x] 2.1 审计 release candidate、handoff、release README、repeatability 是否仍把 current degraded/block 描述为普通 non-blocking。
- [x] 2.2 审计 P64/P65 包材料是否需要在 P65-P67/P68 commit 后重新生成最终分发包，或只需保留 P69 建议。
- [x] 2.3 审计所有 release-ready 文案是否明确排除 current local data clean/healthy claim。
- [x] 2.4 审计安全边界文案，确认未新增交易、外部推送、自动修复、自动迁移、自动升级、真实 provider 保证或收益承诺。

## 3. P68 决策记录与材料刷新

- [x] 3.1 新增 `docs/release/acceptance/2026-06-18-p68-release-readiness-governance.md`，记录 P66/P67 当前命令证据、文档审计结论、发布状态结论和下一阶段建议。
- [x] 3.2 更新 `docs/release/release-candidate-2026-06-18.md`：把顶层状态改为 P68 后的限域状态，引用 P66/P67/P68 证据，禁止 current data clean claim。
- [x] 3.3 更新 `docs/release/release-handoff-2026-06-18.md`：明确当前可交付状态、P67 scope exclusion 边界、是否需要 P69 package refresh。
- [x] 3.4 更新 `docs/release/README.md` 和 `docs/release/acceptance-repeatability.md`：让后续操作者能重复 P66/P67/P68 判定。
- [x] 3.5 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。

## 4. 验证与安全扫描

- [x] 4.1 运行 P67 resolution check，记录 `claim_state`、`policy`、`gate`、`resolution` 和 `clean_data_claim`。
- [x] 4.2 运行 P66 strict gate，记录 expected non-zero 与 `policy=blocked` / `gate=block`。
- [x] 4.3 运行 release 文案扫描，确认无 `current data clean`、`current data healthy`、`policy passed` 等误导性声明。
- [x] 4.4 运行禁止能力扫描，确认无新增券商、交易、一键下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺、登录源、付费源、授权源、Level2 或高频源承诺。
- [x] 4.5 运行 `openspec validate p68-post-p67-release-readiness-governance --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 4.6 若 P68 意外修改运行时代码，追加 `go test ./...`、`npm --prefix web test`、`npm --prefix web run build`、`bash scripts/e2e-smoke.sh`。

## 5. 复审、归档与提交

- [x] 5.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 5.2 执行 OpenSpec archive，把 P68 delta 合并到 `docs/` 真源。
- [x] 5.3 archive 后确认无活跃 change，并规划下一阶段。
- [x] 5.4 提交前子 agent 复审无 Critical / Important。
- [x] 5.5 提交 P68。
