# Tasks: P19–P24 历史审计追溯治理

## 1. 目标与范围确认

- [x] 1.1 确认本阶段只做历史审计追溯，不实现运行时功能。
- [x] 1.2 列出 P19–P24 的交付范围与 archive 状态（已交付/未归档）。
- [x] 1.3 明确该阶段不可伪造历史 archive 与交付时间。

## 2. 追溯说明文档

- [x] 2.1 在 `docs/` 新建 `p19-p24-historical-archive-traceability.md`。
- [x] 2.2 说明每个阶段（P19、P20、P21、P22、P23、P24）的当前状态：
  - [x] 已交付功能边界。
  - [x] 当前未归档项。
  - [x] 可核验文档入口（`docs/development-plan.md`、`docs/README.md` 等）。
- [x] 2.3 在说明中统一加入“仅审计说明、不得替代代码/历史 archive”的风险边界。

## 3. OpenSpec delta 与治理约束

- [x] 3.1 新增 `openspec/changes/p19-p24-historical-archive-traceability/specs/traceability/spec.md`。
- [x] 3.2 约束内容包括：
  - [x] 历史追溯必须可核对。
  - [x] 不新增或伪造历史 archive。
  - [x] 不改变 L1 契约和运行时边界。

## 4. 验证

- [x] 4.1 `openspec validate p19-p24-historical-archive-traceability --strict`
- [x] 4.2 `openspec validate --all --strict`
- [x] 4.3 `git diff --check`
- [x] 4.4 进行只读复审：确认文档无“伪造已归档/新增交易能力”类高风险表述。

## 5. 归档准备

- [x] 5.1 复审通过后执行 `/opsx:archive`。
- [x] 5.2 确认该 change 已移入 `openspec/changes/archive/`。
