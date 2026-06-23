# Design: P0 工程骨架

## 后端（P0.1）

- **布局**：按 `docs/architecture.md` 的 `cmd/`、`internal/`、`pkg/`、`configs/` 创建目录；首版 `internal/domain` 可为空包或占位。
- **HTTP**：`net/http` + `pkg/httputil`；统一 JSON 信封预留 `request_id`、`data`、`error`（与 `docs/frontend-contract.md` 对齐，health 可最小实现）。
- **配置**：`configs/config.example.yaml` + 环境变量覆盖；字段：服务端口、SQLite 路径、VecLite 路径、DeepSeek、日志级别（P0 可不连真实 DB/模型）。
- **健康检查**：`GET /api/v1/health` → `{"status":"ok"}`。

## 前端（P0.2）

- **工具链**：React 18+、Vite、TypeScript、根目录 `web/`。
- **路由**：今日纪律、持仓、决策咨询、情报与证据、规则与纪律、复盘与审计、设置（空页面 + 布局骨架）。
- **API client**：`web/src/services/`，统一处理 `request_id`、`data`、`error` 类型占位。
- **样式**：全局状态类名占位（正常、信息不足、冻结观察、高危），见 `docs/ui-design.md`。

## 决策

| 项 | 选择 | 理由 |
| --- | --- | --- |
| DI | `main.go` 手动注入 | 与 architecture 一致 |
| 前端状态 | 首版可用轻量本地 state | P5 再引入复杂状态库 |

## 风险

- 无；P0 不触碰交易与规则逻辑。
