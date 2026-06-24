# Proposal: P113 布局高保真精修

## Summary

P113 承接 P111/P112 的参考图视觉方向，针对用户继续指出的“好多布局问题、错位、不够精致”执行一次更严格的布局精修与逐页复审。P113 不改变产品信息架构的大方向，不新增后端能力，而是把已发现的移动端横向溢出、卡片压缩、触控目标过小、二级页面旧式块状信息、决策详情层级不足和局部文本错位问题修到可验收。

参考图仍为 P111 锁定的第二方案：

`/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

## Why

P112 虽已完成全路由截图与复审，但用户复核后仍能看到布局问题，说明上一轮 pass 标准过于乐观。P113 将已知问题作为阻断项处理：

- 移动端 `/data-quality`、`/settings`、`/local-install`、`/local-knowledge` 等页面的 metric/card 区域可能横向溢出、裁切或出现滚动轨。
- 桌面 `/data-quality` 顶部 hero 将多个状态卡和侧栏动作挤在同一行，信息密度不稳。
- `/decisions/:id` 决策详情仍有普通后台卡片感，与参考图的 report/status/ledger 层级差距明显。
- 多个移动页面的 action link / 文本按钮高度不足，不符合可点击目标和精致感要求。
- `/rules` 等治理页面过早暴露 raw JSON/工程化内容，破坏用户可读的 report surface。
- 运维/设置类页面的第二屏仍像表单或读数堆叠，缺少参考图式的分区、边界和密度。

## What Changes

- 收紧共享 responsive CSS，取消移动端关键 metric 区的横向 rail，改为稳定的 compact grid/list。
- 调整 report hero 在桌面和 390px 移动视口下的栅格、卡片尺寸、文字换行和边界，避免挤压、裁切、错位。
- 提升移动端链接和操作控件的最小高度、触控区域与视觉 affordance。
- 改善决策详情、数据质量、规则、设置、本地安装、本地知识等二级页面的 report composition。
- 隐藏或折叠工程化 raw 内容，让用户先看到摘要、状态和下一步动作。
- 重新采集全 18 个路由的桌面与 390px 移动截图，逐页检查横向溢出、错位、重叠、压缩、触控尺寸和参考图一致性。

## In Scope

- 前端 React/Vite/TypeScript 页面与共享组件。
- CSS tokens、响应式布局、卡片密度、可点击目标、页面首屏 composition。
- P113 screenshot evidence、mismatch ledger、逐页 layout QA 和验收记录。

## Out Of Scope

- 不新增后端 API、SQLite schema、Eino workflow、LLM 能力、RAG/VecLite 能力、真实数据源或投资规则。
- 不改变最终裁决逻辑、规则提案逻辑、确认流程、审计语义、release/package/version 行为。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、真实库覆盖或收益承诺。
- 不新增登录源、付费源、授权源、Level2、高频源或外部商业数据依赖。

## Validation

- `openspec validate p113-layout-fidelity-polish --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- 启动真实本地后端与 Vite 前端。
- 采集全 18 个桌面路由和全 18 个 390px 移动路由截图。
- 对每个页面检查 no-overflow、文本裁切、按钮高度、首屏信息层级和参考图一致性。
- 使用 `view_image` 对参考图、代表性桌面页面、代表性移动页面做人工对比。
- 生成 P113 mismatch ledger；P0/P1/P2 布局问题必须修复后重新截图。
- Forbidden affordance scan：不得新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺。
- Sensitive/redaction scan：不得暴露完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。
- `openspec validate --all --strict`
- `git diff --check`
