# Tasks: P7 审查问题修复

## 1. P7.3 预期收益分析服务

- [x] 1.1 为 `expectedReturnStep` 增加 `AnalystService` 调用。
- [x] 1.2 将预期收益分析材料写入 `analyst_reports[expected_return]` 或等价结构。
- [x] 1.3 LLM 不可用时返回降级状态，最终裁决仍由规则引擎负责。
- [x] 1.4 补充对应单元测试与中文注释。

## 2. P7.2 检索降级与审计

- [x] 2.1 增加检索服务接口，表达命中证据、降级原因和检索状态。
- [x] 2.2 证据读取节点优先使用检索服务；VecLite 不可用时降级到 SQLite 摘要。
- [x] 2.3 摘要不足时返回信息不足或既有错误状态。
- [x] 2.4 审计或工作流上下文记录检索输入、命中证据和降级原因。
- [x] 2.5 修正独立信源数量按实际来源去重计算。
- [x] 2.6 补充对应单元测试与中文注释。

## 3. 配置到生产依赖串联

- [x] 3.1 增加基于配置构建 `WorkflowDependencies` 的组装函数。
- [x] 3.2 DeepSeek key 存在时使用 DeepSeek client；缺失时使用降级或 stub 服务。
- [x] 3.3 `data_sources.use_stub` 生效，默认本地 stub 不依赖公网。
- [x] 3.4 补充配置组装测试，不写真实密钥。

## 4. 验收

- [x] 4.1 执行 `go test ./internal/infrastructure/... ./internal/application/...`。
- [x] 4.2 执行 `go test ./internal/infrastructure/... ./internal/application/workflow/...`。
- [x] 4.3 执行 `go test ./internal/application/workflow/... ./internal/infrastructure/...`。
- [x] 4.4 确认未实现 P8/P9 范围，未新增自动交易、一键交易或真实密钥。
- [x] 4.5 执行 `openspec instructions apply --change p7-review-fixes --json`，确认 artifacts 完整且 tasks 全部完成。
