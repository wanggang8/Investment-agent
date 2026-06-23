# P30 真实环境 E2E / Playwright smoke 验收任务

## 1. 范围与基线

- [x] 1.1 确认现有前端测试、构建、后端服务启动方式和本地配置入口。
- [x] 1.2 确认 Playwright 或浏览器 smoke 需要的临时目录、端口和清理策略。
- [x] 1.3 明确 smoke 使用 stub / 临时 SQLite / 本地配置，不依赖真实密钥或真实交易能力。

## 2. E2E smoke 实现

- [x] 2.1 新增或扩展本地 E2E smoke 脚本，能够启动后端 server 与前端页面，验证健康检查和关键 UI 可达。
- [x] 2.2 覆盖决策详情或主动咨询相关展示，至少验证 expected return 空数组、情景、动态卖出提示或降级态不会导致前端崩溃。
- [x] 2.3 覆盖证据/审计相关只读路径，验证公开证据刷新或降级审计能够在 UI/API 中可见。
- [x] 2.4 确保测试使用临时 SQLite、临时配置和可控数据，执行后不留下未跟踪日志目录或数据库文件。

## 3. 临时产物治理

- [x] 3.1 将 `.playwright-mcp/`、Playwright 输出、trace、截图、临时 SQLite 等本地 smoke 产物纳入忽略或清理策略。
- [x] 3.2 确认 smoke 失败时仍尽量保留有用诊断，但不会默认污染 git 工作树。

## 4. 文档同步

- [x] 4.1 更新 `docs/development-plan.md`，新增 P30 阶段、目标和验收命令。
- [x] 4.2 更新 `docs/testing-plan.md` 或相关运维文档，记录如何运行本地 E2E smoke。
- [x] 4.3 更新 `docs/README.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 和 `openspec/PROGRESS.md` 的当前阶段状态。

## 5. 验收

- [x] 5.1 运行新增 E2E smoke 命令并确认通过。
- [x] 5.2 运行 `go test ./...`。
- [x] 5.3 运行 `npm --prefix web test -- --run`。
- [x] 5.4 运行 `npm --prefix web run build`。
- [x] 5.5 运行 `openspec validate --all --strict`。
- [x] 5.6 执行 `git status --short`，确认没有 `.playwright-mcp/` 或其他 smoke 临时产物遗留。
