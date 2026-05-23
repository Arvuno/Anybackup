---
name: ag-ui-design
description: Use when an Agent must design user-visible AG-UI Markdown constraints before ag-ui-response publishes MQ, including layout, components, semantic labels, thought summaries, tool-call presentation, results, errors, clarification, and hidden-content boundaries.
---

# AG-UI 设计技能

## 核心职责

本技能只负责生成 Markdown 设计约束，不返回 JSON，不生成 draft，不校验 draft，不发布 MQ，不调用业务工具。输出的 Markdown 作为后续 `ag-ui-response` 生成 draft JSON 的高优先级上下文约束。

`ag-ui-design` 完全负责用户可见内容定义：布局、组件选择、思考摘要、工具调用展示、主动询问、结果、错误、候选方案和最终报告的可见内容边界。`ag-ui-response` 只负责把这些约束协议化、校验和发布。

## 工作流

1. 在任何非极短 `thought` 的 AG-UI 输出前，先读取本技能。
2. 按当前用户请求、业务阶段和已知事实生成一段 Markdown 设计约束。
3. Markdown 必须包含固定章节，不得输出 JSON、代码块形式的 draft 或 MQ 消息。
4. 将 Markdown 设计约束作为高优先级上下文交给 `ag-ui-response`，再由 `ag-ui-response` 生成 draft JSON 字符串并发布。

## 固定输出章节

每次输出必须严格包含以下章节，章节名不可改：

```markdown
## 显示目标
## 内容约束
## 布局约束
## 思考链展示
## 工具调用展示
## 结果返回展示
## 禁止展示
```

## 渐进式披露

- 开始设计任何用户可见 AG-UI 输出时，先读取 `references/design-principles.md`，确定整体风格、信息层级、状态语义和动作分级。
- 将设计约束映射到 AG-UI 协议字段时，读取 `references/ag-ui-protocol-mapping.md`，明确 layout tree、activity、state、actions、sequence 和更新边界。
- 选择 `thought`、`tool_call`、`clarification`、`result`、`error` 的布局和组件时，读取 `references/layout-patterns.md`。
- 定义思考摘要、工具调用、结果、错误、候选方案和报告可见内容时，读取 `references/content-rules.md`。
- 判断哪些内部标识、工具名、参数、返回、知识网络信息或敏感内容不能展示时，读取 `references/visibility-policy.md`。
- 需要参考输出形态时，只读取相关 `references/examples/*.md` 示例；示例是 Markdown 约束，不是 draft JSON。资产选择、备份方案审查、运营概览、告警分析、恢复演练和方案对比优先参考对应业务场景示例。

## 设计原则

- 用户可见内容使用业务语义名称，不显示函数式工具名。
- 业务知识网络名称、ID、对象类 ID、关系类 ID、动作类 ID 不向用户展示。
- 工具输入和结果只展示脱敏摘要，不展示原始参数、原始返回、连接串、日志、堆栈或系统提示词。
- 思考链只展示阶段性判断摘要和下一步，不展示内部推理链。
- 业务结果必须有结构化布局建议，不能只要求自然语言段落。
- 空结果、失败和需要确认的情况也必须给出可渲染组件和用户可执行的下一步。
- 终态结果设计必须以 AI 已完成结果输出为结果源文本；发布前冻结该文本，并要求 AG-UI 卡片与该文本做一致性校验。
- 设计约束必须贴合 AG-UI 协议，明确卡片内容模型、layout tree 组件选择、`actions` 动作承载和状态更新方式。
- 用户问询、确认和审查类输出优先设计为轻量对话卡片 + 选项按钮组，正文短、问题明确、动作按钮贴近卡片底部。
- 方案类输出优先设计为方案审查卡：标题、对象摘要、关键参数、AI 生成思路、信心条、可展开详情和固定操作栏。
- 所有设计先抽象为“目标-结论-事实-风险-动作”的信息结构，再选择组件；不得只套用示例文本或图片形态。
