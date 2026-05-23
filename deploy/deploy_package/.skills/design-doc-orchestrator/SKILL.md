---
name: design-doc-orchestrator
description: 编排 `docs/00_设计文档/` 下的复杂设计文档更新、拆分、补图和回收。Use when requests mention 设计文档、架构图、部署图、时序图、组件边界、依赖矩阵、统一更新多份设计文档, or when a design change impacts more than one file under `docs/00_设计文档/`.
---

# Design Doc Orchestrator

把 `docs/00_设计文档/` 当成一组受控交付物来维护,先判断这次设计变更会影响哪些文档,再把具体写作、补图和 review 交给合适的 skill。

## First Reads

- 先读 [references/source-map.md](references/source-map.md),确认当前设计文档目录和代码侧事实来源。
- 如需新增文档,再读 [assets/design-doc-template.md](assets/design-doc-template.md)。
- 如果变更触及系统边界、组件归属或 KWeaver 集成口径,补读 [../anybackup-v9-architecture-charter/SKILL.md](../anybackup-v9-architecture-charter/SKILL.md)。
- 如果变更触及部署流程、chart、Ansible、验证脚本或交付前提,补读 [../v9-deployment-verifier/SKILL.md](../v9-deployment-verifier/SKILL.md)。

## Working Mode

### 1. Classify the change

先把任务归入一个或多个变更面:

- 总体拓扑变化: 组件增删、连接关系、命名空间归属、安装边界变化。
- 部署链路变化: 安装步骤、前置条件、控制机/目标机职责、验证路径变化。
- 运行时变化: 环境变量、中间件依赖、服务间调用关系变化。
- 文档结构变化: 需要把一份设计文档拆成多份,或把多份说明收敛成统一入口。

### 2. Define the output package

在动笔前先确定交付范围:

- 单文件更新: 直接更新现有设计文档。
- 多文件联动: 明确主文档、辅文档和图表落点。
- 新增文档: 只有当现有文档找不到合适落点,或单文档已经难以维持清晰主线时才新增。

默认保持 `docs/00_设计文档/` 的编号稳定。没有明确收益时,不要为了“更整齐”而重排整个目录。

### 3. Route the work

按最小必要原则调度:

- 文本改写、章节增删、表格同步: 用 [../design-doc-change-handler/SKILL.md](../design-doc-change-handler/SKILL.md)。
- Mermaid 架构图、部署图、时序图更新: 用 [../design-diagram-maintainer/SKILL.md](../design-diagram-maintainer/SKILL.md)。
- 设计文档 review、找事实偏差、找图文不一致: 用 [../design-doc-reviewer/SKILL.md](../design-doc-reviewer/SKILL.md)。
- 如果代码现状与目标边界冲突,先用 [../anybackup-v9-architecture-charter/SKILL.md](../anybackup-v9-architecture-charter/SKILL.md) 收口口径,再回写文档。

### 4. Apply the review gate

以下情况默认要走 review:

- 新增、删除、移动 `docs/00_设计文档/` 下的正式文档。
- 修改总体部署架构、运行时组件边界、命名空间归属或人工步骤职责。
- 更新 Mermaid 图,且图中节点、边或时序关系发生变化。

只有错别字、轻微措辞或不会改变事实的格式修正,才可以跳过专门 review。

### 5. Close the loop

完成前统一检查:

- 文档标题、编号、引用路径是否仍然稳定。
- 图文是否同步。
- 新增的设计结论是否能在代码、chart、Ansible 或脚本里找到事实支撑。
- 如果改动部署逻辑,是否需要补跑 `python scripts/validate_deploy.py --mode regression`。

## Guardrails

- 不要把临时现场约束写成长期产品边界。
- 不要只根据一处 YAML 或脚本就推导整套部署口径,至少交叉检查 `ansible/`、`helm-chart/` 和现有设计文档中的一个。
- 不要为了创建新文档而创建新文档; 先优先维护现有 `docs/00_设计文档/`。
- 不要发明代码里不存在的连接关系、环境变量或中间件职责。
- 不要跳过 review 后仍然宣称“设计文档已完全对齐”。

## Done Criteria

- 受影响的设计文档已明确落位。
- 如果需要,文本更新、图更新和 review 已串起来。
- 最终回复能说清:
  - 改了哪些设计文档。
  - 哪些事实来自代码或脚本。
  - 哪些地方仍需后续设计决策。
