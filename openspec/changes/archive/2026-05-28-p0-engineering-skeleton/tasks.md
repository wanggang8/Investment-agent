# Tasks: p0-engineering-skeleton

> 验收标准摘自 `docs/development-plan.md` P0。

## 1. Go 后端（P0.1）

- [x] 1.1 初始化 Go module（`go mod init`，模块名与目录一致）
- [x] 1.2 创建 `cmd/server/main.go`、`cmd/agent/main.go`（agent 可为占位）
- [x] 1.3 实现 `GET /api/v1/health` → `{"status":"ok"}`
- [x] 1.4 添加 `configs/config.example.yaml` 与配置加载（端口、SQLite/VecLite 路径、DeepSeek、日志级别）
- [x] 1.5 添加 `pkg/logger/`、`pkg/httputil/` 最小实现
- [x] 1.6 创建 `internal/` 目录骨架（domain/application/infrastructure 占位）
- [x] 1.7 编写 `docs/configuration.md` 启动说明
- [x] 1.8 验收：`go test ./...`、`go run ./cmd/server`、`curl localhost:8080/api/v1/health`

## 2. 前端（P0.2）

- [x] 2.1 在 `web/` 初始化 React + Vite + TypeScript
- [x] 2.2 配置基础路由（7 个页面占位）
- [x] 2.3 实现 API client（`request_id`、`data`、`error` 类型与封装）
- [x] 2.4 添加全局状态样式类（正常、信息不足、冻结观察、高危）
- [x] 2.5 验收：`cd web && npm install && npm run build && npm run dev`

## 3. 归档前

- [x] 3.1 将 `specs/api/spec.md` 中 ADDED 合并进 `docs/api.md`
- [x] 3.2 勾选 `docs/development-plan.md` P0 相关任务
- [x] 3.3 `/opsx:archive` 或 `openspec archive p0-engineering-skeleton`
