# Design: 已完成代码中文注释

## 范围

注释覆盖已实现模块：

- Go 后端入口、配置、日志、HTTP 工具。
- SQLite migration、迁移执行器、Repository 接口与实现。
- 领域模型和规则引擎。
- React 前端 API client、类型、页面与布局。

## 注释原则

- 注释解释业务含义、边界和约束，不复述显而易见的语法。
- Go 导出类型和函数优先补充中文注释，保持 gofmt 兼容。
- SQL 用 `--` 标注表用途和关键约束。
- TypeScript 用简短注释说明 DTO、API client 与页面职责。

## 验证

- `gofmt` 格式化 Go 文件。
- `go test ./...` 验证后端行为不变。
- `npm run build` 验证前端类型与构建不受影响。
