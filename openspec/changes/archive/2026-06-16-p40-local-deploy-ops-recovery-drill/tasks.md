## 1. OpenSpec 与范围

- [x] 1.1 确认 P40 只覆盖本地部署、运维、自检、恢复演练和诊断治理。
- [x] 1.2 确认 P40 不接券商 API、不自动交易、不外部推送、不自动应用规则、不承诺收益。
- [x] 1.3 对齐 P33-P39 已有账户/持仓、公开数据、风险预警、规则治理、LLM、retrieval quality 和浏览器 E2E 能力。

## 2. 本地预检与启动诊断

- [x] 2.1 增加本地预检命令或脚本，检查 Go、Node、npm、Playwright browser、SQLite path、VecLite path、配置文件和目录权限。
- [x] 2.2 预检输出 pass / warning / failed / skipped 状态和可修复提示。
- [x] 2.3 启动前诊断复用配置校验、迁移检查和本地依赖状态，不输出密钥原文。
- [x] 2.4 诊断失败写入安全日志、审计或诊断文件，不触发交易、外部推送或规则应用。

## 3. 备份恢复演练

- [x] 3.1 增加备份恢复 smoke，默认使用临时 SQLite / VecLite / 配置路径。
- [x] 3.2 恢复目标已有 DB 时必须显式确认或拒绝覆盖。
- [x] 3.3 恢复后通过 API 或浏览器读取历史决策、审计、持仓或报告事实。
- [x] 3.4 记录恢复诊断和失败分类，避免污染真实本地数据。

## 4. 数据源健康与前端运维面板

- [x] 4.1 展示数据源最近成功时间、最近失败时间、失败分类、新鲜度和影响范围。
- [x] 4.2 健康面板覆盖 fresh / stale / failed / missing / unknown 状态。
- [x] 4.3 面板只提供查看、过滤或本地刷新入口，不提供交易、外部推送或自动恢复承诺。
- [x] 4.4 前端状态、空态和错误态具备安全文案。

## 5. 临时文件、日志与 gitignore 治理

- [x] 5.1 明确 tmp、Playwright output、诊断文件、恢复演练输出和本地日志目录。
- [x] 5.2 确认生成物被 `.gitignore` 覆盖或脚本退出清理。
- [x] 5.3 文档说明不要使用真实私有数据库复现 smoke。

## 6. 文档与验收

- [x] 6.1 在 P40 delta 中记录待归档合并到 `docs/configuration.md`、`docs/ops-local-scheduler.md`、`docs/frontend-contract.md` 的运行/恢复/健康面板契约。
- [x] 6.2 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P40 active 状态。
- [x] 6.3 运行 `go test ./...`。
- [x] 6.4 运行 `npm --prefix web test -- --run`。
- [x] 6.5 运行 `npm --prefix web run build`。
- [x] 6.6 运行本地部署/恢复 smoke。
- [x] 6.7 运行 `bash scripts/e2e-smoke.sh`。
- [x] 6.8 运行 archive 前只读子 agent 复审，且无 Critical / Important 问题。
- [x] 6.9 运行 `openspec validate p40-local-deploy-ops-recovery-drill --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
