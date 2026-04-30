---
name: design-diagram-maintainer
description: 维护 `docs/00_设计文档/` 中的 Mermaid 架构图、部署图、时序图和依赖关系图。Use when users ask to update 架构图、部署图、时序图、组件关系图, or when text changes make existing diagrams under `docs/00_设计文档/` stale.
---

# Design Diagram Maintainer

专门负责让 `docs/00_设计文档/` 里的 Mermaid 图和正文保持一致,避免“文字已经更新,图还停留在旧架构”。

## First Reads

- 先读 [../design-doc-orchestrator/references/source-map.md](../design-doc-orchestrator/references/source-map.md)。
- 再读 [references/diagram-patterns.md](references/diagram-patterns.md)。
- 然后读取目标文档里现有的 Mermaid 图和图前后的解释文字。

## Workflow

### 1. Choose the diagram type

根据问题选择最小图形:

- 总体组件和命名空间归属: `flowchart`
- 安装和初始化步骤: `sequenceDiagram`
- 运行时依赖矩阵或连接关系: `flowchart`

如果一张图同时想表达拓扑、时序和责任分工,优先拆成两张图。

### 2. Rebuild from source facts

只根据已确认事实更新节点和边:

- chart 中存在的组件和 Service。
- Ansible 中实际执行的步骤。
- 文档中已经写清楚并且能被代码验证的职责边界。

看不清楚的关系,宁可在正文里标成“待确认”,也不要在图里虚构箭头。

### 3. Keep the diagram readable

- 同一类节点用同一命名方式。
- 节点标签包含括号、中文标点或较长短语时,使用 Mermaid 引号标签。
- 让主路径从左到右或从上到下保持稳定,不要来回折线。
- 不要把解释性长句塞进节点里,长解释放到图下文字。

### 4. Sync the surrounding prose

每次改图后,同步更新:

- 图标题。
- 图下的“读图说明”。
- 引用了旧节点名或旧流程的正文段落。

## Guardrails

- 不要为了“好看”更换图的语义类型。
- 不要在没有事实支撑时引入新组件、新依赖或新时序。
- 不要让图里的命名比代码或正文更抽象,除非文中明确解释。
- 不要用一张超大图替代两张可读的小图。

## Done Criteria

- Mermaid 语法可渲染。
- 图中的节点、连线和顺序与正文一致。
- 图前图后没有遗留旧名称或旧解释。
