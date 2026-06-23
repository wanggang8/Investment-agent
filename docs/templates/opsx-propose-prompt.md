# `/opsx:propose` 提示词模板

## 自动下一阶（推荐，不用每次写 P几）

复制下面整段到 Cursor（**无需改 P 编号**）：

```text
/opsx:propose

读取 @openspec/PROGRESS.md 与 @docs/development-plan.md：
1. 若存在未 archive 的活跃 change → 提示继续 /opsx:apply 或 archive，不要新建
2. 否则取 next_change_id，为下一阶段创建 change
3. tasks.md 必须逐条对齐 development-plan 中该阶段全部小节与验收命令
4. specs/ 只写本阶段 delta，合并目标见 openspec/project.md 阶段表
5. 不要发明 development-plan 以外的需求

生成后更新 proposal 的 in/out scope，并说明与 plan 条目是否一一对应。
```

进度真源：`openspec/PROGRESS.md`（archive 时更新）。

---

## 手动指定阶段（可选）

每次新阶段：复制下面「完整提示词」到 Cursor 聊天，**只改【】里的 5 处**，并 @ 对应文档。

---

## 每次必改（5 处）

| 序号 | 占位符 | 示例 |
| --- | --- | --- |
| ① | `【change-id】` | `p1-data-foundation` |
| ② | `【阶段标题】` | `P1：数据底座` |
| ③ | `【development-plan 小节】` | `P1.1、P1.2、P1.3`（照抄 plan 里的小节号） |
| ④ | `【要 @ 的契约文档】` | `data-model.md`（本阶段相关的 L1 文件） |
| ⑤ | `【delta 合并目标】` | `docs/data-model.md` |
| ⑥ | `【out of scope】` | 下一阶段及之后（见下方阶段表） |

---

## 阶段速查表（⑥ out of scope 照抄本行「排除」列）

| 阶段 | ① change-id | ② 标题 | ③ plan 小节 | ④ @ 文档 | ⑤ delta 合并到 | 排除（out of scope） |
| --- | --- | --- | --- | --- | --- | --- |
| P0 | `p0-engineering-skeleton` | P0 工程骨架 | P0.1、P0.2 | architecture, frontend-contract, ui-design | `docs/api.md`（health） | P1–P6（**已有 change 勿重复 propose**） |
| P1 | `p1-data-foundation` | P1 数据底座 | P1 下全部小节 | data-model, architecture | `docs/data-model.md` | P2–P6 |
| P2 | `p2-domain-rules` | P2 领域规则 | P2 下全部小节 | data-model, workflow, requirements | `docs/data-model.md`、`docs/workflow.md` | P3–P6 |
| P3 | `p3-eino-workflows` | P3 工作流 | P3 下全部小节 | workflow, architecture | `docs/workflow.md` | P4–P6 |
| P4 | `p4-http-api` | P4 HTTP API | P4 下全部小节 | api, data-model, frontend-contract | `docs/api.md` | P5–P6 |
| P5 | `p5-frontend-cockpit` | P5 前端驾驶舱 | P5 下全部小节 | frontend-contract, ui-design, ui-flow, api | `docs/frontend-contract.md` | P6 |
| P6 | `p6-e2e-hardening` | P6 验收加固 | P6 下全部小节 | functional-spec, development-plan | 按 change 实际改动的 docs | 无（收尾阶段） |

---

## 完整提示词（复制到 Cursor）

```text
/opsx:propose 【change-id】

创建 OpenSpec change，严格依据现有文档，不要发明新需求。

1. 阶段：docs/development-plan.md 的「【阶段标题】」，小节：【development-plan 小节】
2. 必读：@docs/development-plan.md @docs/【要 @ 的契约文档】 @docs/GOVERNANCE.md @openspec/project.md
3. tasks.md：
   - 逐条对应 development-plan 的任务列表与验收 bash
   - 每条用 - [ ] checkbox，保留原文验收命令
4. specs/：只写本阶段 delta（ADDED/MODIFIED/REMOVED），合并目标：【delta 合并目标】
5. proposal.md：写清 in scope / out of scope

out of scope：【out of scope】

生成后列出 tasks 条数与 development-plan 是否一一对应。
```

---

## 生成后检查（30 秒）

- [ ] `openspec/changes/【change-id】/tasks.md` 与 plan 任务一致
- [ ] `proposal.md` 的 out of scope 正确
- [ ] delta 未改无关契约

通过 → `/opsx:apply`  
做完验收 → 合并 delta 到 docs → `/opsx:archive`

---

## P1 填好示例（对照格式）

```text
/opsx:propose p1-data-foundation

创建 OpenSpec change，严格依据现有文档，不要发明新需求。

1. 阶段：docs/development-plan.md 的「P1：数据底座」，小节：P1 下全部小节
2. 必读：@docs/development-plan.md @docs/data-model.md @docs/architecture.md @docs/GOVERNANCE.md @openspec/project.md
3. tasks.md：
   - 逐条对应 development-plan 的任务列表与验收 bash
   - 每条用 - [ ] checkbox，保留原文验收命令
4. specs/：只写本阶段 delta（ADDED/MODIFIED/REMOVED），合并目标：docs/data-model.md
5. proposal.md：写清 in scope / out of scope

out of scope：P2–P6（领域规则、Eino、业务 API、前端页面、E2E）

生成后列出 tasks 条数与 development-plan 是否一一对应。
```
