# P19-P24 历史审计追溯设计

## 设计目标

以“只读追溯+不补写历史”为原则，把 P19–P24 的审计信息固定为一个文本化证据入口，避免未来团队在审计时误把“未归档历史状态”当成可追溯事实。

## 做法

1. **新增追溯说明文档**
   - 在 `docs/p19-p24-historical-archive-traceability.md` 建立统一文档。
   - 按 P19–P24 顺序记录：
     - 功能交付状态摘要
     - 当前档案状态（无 archive / 已 archive）
     - 说明当前文档映射点（`development-plan`/`configuration`/`ops-local-scheduler` 等）
   - 明确注明：本文件为审计说明，不是历史 archive 正文。

2. **新增 change 级 spec**
   - 在 `openspec/changes/p19-p24-historical-archive-traceability/specs/traceability/spec.md` 增加治理级需求。
   - 需求要求包括：
     - 审计材料必须可核对。
     - 不得伪造完成证据。
     - 历史阶段边界保留为“交付/未归档/待补充”而非“已归档”。

3. **执行与复审**
   - 所有修改仅限文档。
   - 通过 `openspec validate` 和 `git diff --check` 验证变更一致性。

## 风险与边界

- 风险：信息不完整导致再次出现边界误读。缓解：直接指向事实源（`development-plan`、`docs/README.md`、`openspec/project.md`）与缺失说明。
- 边界：不改动功能实现，不回写历史阶段的运行代码与审计事实，不造假交付时间。

## 验收结果目标

- `openspec` 校验通过。
- 追溯文档中每个阶段都明确是否已有 archive 与缺口。
- 无敏感信息外露，未触及交易与自动化能力边界。
