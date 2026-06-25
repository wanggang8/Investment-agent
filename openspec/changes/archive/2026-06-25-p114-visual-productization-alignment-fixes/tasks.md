# Tasks: P114 视觉产品化与对齐残留修复

## 1. Governance

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 创建并校验 `p114-visual-productization-alignment-fixes`。
- [x] 1.3 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`，将 P114 标记为当前活跃 change。
- [x] 1.4 明确 P114 不扩大投资运行时能力，不处理发布/安装/物理第二机器。

## 2. Fresh Visual Audit

- [x] 2.1 启动本地 backend/frontend，采集桌面和 390px 移动截图。
- [x] 2.2 审查表单与按钮错位、动作区基线、移动端按钮堆叠。
- [x] 2.3 审查同层级卡片高度、标题区、内边距、底部动作位置和视觉重量。
- [x] 2.4 审查 raw/diagnostic/命令/路径/内部枚举是否进入首层产品 UI。
- [x] 2.5 建立 P114 finding ledger，逐项标记 `frontend-mapping` / `backend-summary-needed` / `intentional-technical-secondary`。

## 3. Fixes

- [x] 3.1 修复共享表单布局、button row、field/action alignment。
- [x] 3.2 修复同层级卡片 grid 等高、底部动作对齐和标题区一致性。
- [x] 3.3 修复诊断/命令/raw 内容产品化展示层级。
- [x] 3.4 修复页面级残留问题，覆盖表单和治理/运维页面。
- [x] 3.5 如确需后端 summary 字段，做最小只读 API/DTO 变更并补测试；如不需要，记录无需后端原因。

## 4. Rendered Re-review

- [x] 4.1 重新采集修复后桌面截图。
- [x] 4.2 重新采集修复后 390px 移动截图。
- [x] 4.3 逐页对比参考图和 P114 finding ledger。
- [x] 4.4 若仍有 P0/P1/P2 产品化视觉问题，继续修复并复验。

## 5. Validation

- [x] 5.1 `openspec validate p114-visual-productization-alignment-fixes --strict`。
- [x] 5.2 `npm --prefix web test -- --run`。
- [x] 5.3 `npm --prefix web run build`。
- [x] 5.4 forbidden affordance scan。
- [x] 5.5 sensitive/redaction scan。
- [x] 5.6 `openspec validate --all --strict`。
- [x] 5.7 `git diff --check`。

## 6. Documentation And Archive

- [x] 6.1 新增 P114 acceptance record，包含 findings、截图路径、后端是否需要修改的结论。
- [x] 6.2 更新 UI / frontend contract / roadmap / governance / progress 文档。
- [x] 6.3 OpenSpec archive。
- [ ] 6.4 最终验证后提交。
