# P31 每日自动运行闭环任务

## 1. 范围与基线

- [x] 1.1 确认当前 `cmd/agent`、`cmd/server`、workflow wiring、notification、audit 和本地调度文档的现状。
- [x] 1.2 明确每日自动运行的入口选择：server 内 scheduler、独立 agent scheduler，或两者之一的最小实现。
- [x] 1.3 明确默认关闭、安全边界、低频运行和本地-only 约束。
- [x] 1.4 明确每日运行 scope：账户/组合、当前持仓、关注标的或配置符号列表。

基线结论：

- `cmd/agent` 当前只提供显式手动任务与 `--schedule` 说明，不驻留后台；`cmd/server` 当前只启动 HTTP API，不启动 scheduler。
- P31 最小入口采用 server 内 scheduler：随本地服务生命周期运行，便于暴露状态 API 和前端展示；`cmd/agent` 保留为手动触发、配置校验和运维入口。
- 默认配置必须关闭自动运行；只有 `daily_auto_run.enabled: true` 后才允许本地低频调度。自动运行只写本地记录、应用内通知和审计，不接券商、不交易、不外推、不承诺收益、不预测确定涨跌。
- P31 scope 优先使用本地账户/组合的当前持仓；缺账户或缺持仓时记录 `missing_prerequisites` 并展示诊断，不生成正式交易建议。配置符号列表只作为后续显式补充 scope，不替代持仓事实。

## 2. 配置与状态模型

- [x] 2.1 新增 daily auto-run 配置项：enabled、run time、timezone/local time、scope、retry、timeout、max symbols。
- [x] 2.2 更新 `configs/config.example.yaml` 和 `docs/configuration.md`，记录默认关闭和安全边界。
- [x] 2.3 设计并实现每日运行状态持久化或可查询模型，记录 disabled/scheduled/running/success/degraded/failed。
- [x] 2.4 定义幂等 key：local date、portfolio/scope、symbol set、task version 或等价字段。

## 3. 后端自动运行编排

- [x] 3.1 实现 scheduler 或等价本地自动运行入口，默认不启动，显式 enable 后才运行。
- [x] 3.2 串联 market refresh、public evidence refresh、evidence/index preparation 和 daily discipline workflow。
- [x] 3.3 对缺账户、缺持仓、缺行情、缺证据、数据源失败、超时等情况做结构化状态和错误分类。
- [x] 3.4 自动运行成功、部分成功、降级和失败都写入 `audit_events`。
- [x] 3.5 自动运行结果写入应用内通知，避免重复通知刷屏。
- [x] 3.6 保证自动运行不会写入 `operation_confirmations` 的 executed 状态，不会创建交易流水，不会调用交易相关接口。

## 4. API 与前端展示

- [x] 4.1 增加或扩展 API，使前端可查询 auto-run enabled 状态、last run、next run、当前运行状态和失败原因。
- [x] 4.2 前端展示每日自动运行状态：关闭、已计划、运行中、成功、部分成功、降级、失败。
- [x] 4.3 前端提供跳转到最新每日决策、通知或审计详情的入口。
- [x] 4.4 前端在缺少账户/持仓/行情/证据等前提时展示可操作的缺失项说明。
- [x] 4.5 页面文案保持本地记录、人工复核和非自动交易边界。

## 5. 幂等、重试与诊断

- [x] 5.1 实现同一日期和 scope 的幂等保护，避免重复生成冲突的每日决策。
- [x] 5.2 支持有限重试或手动重跑，并在审计中区分 retry、reuse、rerun 或 degraded result。
- [x] 5.3 增加运行超时保护，超时后记录失败状态和可读原因。
- [x] 5.4 增加运行诊断日志或审计摘要，便于定位哪个步骤失败。

## 6. 测试与 E2E 验收

- [x] 6.1 为配置校验、scheduler 触发、幂等、重试和禁交易边界增加 Go 测试。
- [x] 6.2 为每日自动运行编排成功、部分失败、缺前提、数据不足增加 workflow/handler 测试。
- [x] 6.3 增加或扩展 Playwright E2E/smoke，验证本地 server + 前端可展示 auto-run 状态和最新结果。
- [x] 6.4 验证默认配置不会自动运行，不会因 server 启动而写入决策或通知。
- [x] 6.5 验证自动运行不会留下未跟踪临时产物。

## 7. 文档同步

- [x] 7.1 更新 `docs/development-plan.md`，将 P31 状态推进为进行中或完成后的状态。
- [x] 7.2 更新 `docs/ops-local-scheduler.md`，说明内置/本地自动运行与 cron/launchd 示例的关系。
- [x] 7.3 更新 `docs/testing-plan.md`，加入 P31 验收命令和结果记录。
- [x] 7.4 更新 `docs/README.md`、`AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/PROGRESS.md` 的当前阶段状态。

## 8. 验收

- [x] 8.1 运行 `go test ./...`。
- [x] 8.2 运行 `npm --prefix web test -- --run`。
- [x] 8.3 运行 `npm --prefix web run build`。
- [x] 8.4 运行 P31 新增或扩展的 E2E smoke 命令。
- [x] 8.5 运行 `openspec validate --all --strict`。
- [x] 8.6 执行 `git status --short`，确认没有 scheduler、E2E 或 smoke 临时产物遗留。
- [x] 8.7 归档前执行只读子 agent 复审，确认无 Critical / Important 问题。
