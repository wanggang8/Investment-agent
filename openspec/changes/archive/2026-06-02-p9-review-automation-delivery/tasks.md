## 1. P9.1 `cmd/agent` 本地任务入口

- [x] 1.1 实现 `cmd/agent`，提供每日纪律、市场刷新、情报索引、复盘任务的本地触发入口。
- [x] 1.2 支持手动触发和本地调度配置，默认不启用任何自动交易能力。
- [x] 1.3 每次任务执行写入 `audit_events`，记录输入摘要、状态和错误码。
- [x] 1.4 任务失败时返回可读错误，并保留已有数据一致性。
- [x] 1.5 对调度、安全边界和审计写入添加必要中文注释。
- [x] 1.6 验收：执行 `go test ./...`。
- [x] 1.7 验收：执行 `go run ./cmd/agent --help`，确认本地任务入口可启动、任务行为可追踪、不会执行交易。

## 2. P9.2 月度/季度复盘与规则有效性评估

- [x] 2.1 完善月度复盘，汇总确认动作、错误案例、规则提案和审计事件。
- [x] 2.2 完善季度复盘，评估规则命中、误判、缺证据和降级情况。
- [x] 2.3 将规则有效性评估结果写入规则提案或复盘摘要，仍需守门人审计和用户最终确认。
- [x] 2.4 前端复盘页展示周期摘要、规则建议和追踪入口。
- [x] 2.5 保持规则提案不会自动应用。
- [x] 2.6 验收：执行 `go test ./internal/application/workflow/... ./internal/application/handler/...`。
- [x] 2.7 验收：执行 `cd web && npm run build`，确认月度/季度复盘可生成摘要，规则变更仍经审计和用户确认。

## 3. P9.3 本地交付与运维说明

- [x] 3.1 补充本地启动、初始化、数据备份、索引重建和恢复说明。
- [x] 3.2 补充常见故障处理：数据源不可用、VecLite 索引损坏、DeepSeek 缺配置、SQLite 写入失败。
- [x] 3.3 补充 P7-P9 后的完整验收命令。
- [x] 3.4 确认文档不包含真实密钥、账号、token 或个人敏感信息。
- [x] 3.5 验收：执行 `go test ./...`。
- [x] 3.6 验收：执行 `cd web && npm run build`，确认本地交付说明完整、关键故障有处理路径、安全边界清晰。

## 4. 对齐检查

- [x] 4.1 确认本 tasks.md 逐条覆盖 `docs/development-plan.md` 的 P9.1、P9.2、P9.3 全部任务。
- [x] 4.2 确认本 tasks.md 覆盖 P9.1、P9.2、P9.3 全部验收命令。
- [x] 4.3 确认 `specs/` 只包含 P9 delta，合并目标与 `openspec/project.md` 阶段表一致。
- [x] 4.4 确认没有加入 development-plan 之外的需求。

## 5. 验收记录

- [x] 5.1 `go test ./...`：通过（2026-06-02），输出 `Go test: 200 passed in 24 packages`。
- [x] 5.2 `go run ./cmd/agent --help`：通过（2026-06-02），本地任务入口展示每日纪律、市场刷新、情报索引、复盘任务，并声明不会执行交易。
- [x] 5.3 `go test ./internal/application/workflow/... ./internal/application/handler/...`：已由全量 `go test ./...` 覆盖；定向复验 `go test ./cmd/agent ./internal/application/handler` 通过（2026-06-02），输出 `Go test: 80 passed in 2 packages`。
- [x] 5.4 `cd web && npm run build && npm test`：通过（2026-06-02），构建 `✓ built in 97ms`，Vitest `18 passed (18)`、`52 passed (52)`。
