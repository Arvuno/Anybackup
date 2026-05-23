# Diagram Patterns

## Preferred Mermaid Types

- `flowchart TD` 或 `flowchart LR`
  - 用于总体架构、部署拓扑、运行时依赖。
- `sequenceDiagram`
  - 用于安装顺序、人工初始化、验证时序。

## Source-first Rule

- 先看代码和脚本里的真实组件、动作和依赖。
- 再看现有设计文档想表达的主线。
- 最后才组织图形布局。

## Authoring Rules

- 节点标签包含括号、冒号、斜杠或中文说明时,用 `A[\"label\"]` 这类写法。
- 同一个组件在同一张图里只保留一个主节点 ID。
- 同类边的方向保持一致。
- 关键人工步骤可以用备注节点或单独 actor,不要混进系统组件。

## Split Guidance

出现以下情况时拆图:

- 一张图同时表达部署时序和运行时拓扑。
- 节点超过 12 个且主线已经难以追踪。
- 一个文档同时需要“谁部署谁”和“谁依赖谁”。

## Final Checks

- 图名是否和章节名匹配。
- 图下文字是否说明了读图顺序和关键结论。
- 是否还有正文在描述旧图结构。
