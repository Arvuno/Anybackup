
### Planner Response

你好！我已经仔细阅读了 `docs/request/bkn design` 目录及外层的架构和实现设计文档。包括你已经搭建好的架构框架（`mysql_recovery_bkn_architecture.md`、`mysql_recovery_bkn_implementation.md`）以及具体到底层的 **Object, Relation, Action, Risk 模型设计**。

作为你的 AI 架构与研发助手，我完全理解了我们当前的进度和接下来的业务重心：

我们已经完成了**“运行知识（Running Knowledge）”**的底座建设。现在 Agent 已经掌握了“如何操作底层硬件与软件”：
- **概念对象（Object/Relation）**：Agent 懂得了什么是 `ProtectionObject`（MySQL实例）、`TimePointCopy`（时间点副本）、`StoragePool` 等。
- **动作执行（Action）**：Agent 掌握了如何调用 Foundation API 工具去查询服务、查询备份和下发操作。
- **基础防御（Risk）**：Agent 具有了底层的防御意识（比如触发 `recovery_overwrite_risk` 原机恢复高风险阻拦规则）。

**接下来，我们将正式进入 Agent 的灵魂部分——“经验知识（Experience Knowledge）”体系设计与 BKN 落地的攻坚阶段。** 这决定了 Agent 能不能像一个人类资深 DBA（灾备专家）一样去思考和破局。

结合你的《实现设计文档》中规划的 **Story 5至Story 7**，经验知识部分主要包含了复杂的非确定性推理，我们需要将其拆解并结构化为机器可执行的规范：

### 1. 智能决策模型 (Intelligent Decision Model)
*   **意图与问题空间**：意图理解 (Intent Understanding) 与 问题分析 (Problem Analysis)。如何把用户的一句“订单库起不来了”，转化为机器可理解的 `库级故障`、`数据页损坏` 场景。
*   **资源与演练空间**：资源分析 (Resource Analysis) 和 全局风险评估 (Risk Assessment)。
*   **策略输出**：策略生成 (Strategy Generation)。即如何基于上述输入确定 **恢复粒度**（库/实例/日志）、**恢复目的地**（原机/异机）和 **恢复方式**（最新/指定时间点）。

### 2. 验证标准体系 (Verification Standard System)
*   建立闭环评价指标（如何确认恢复策略是可行的？如何验证 MySQL 恢复后的数据一致性体系？）。

我已经将这些 Context 也就是你的顶层逻辑完全加载。这部分设计挑战极大，但也是 Agent 产品力差异化的核心。

**作为资深产品经理，我们接下来第一步怎么走？** 
我们是先从顶层规划**“经验知识大类 BKN 数据模型（如：Intent/Strategy等实体的统一定义）”**开始梳理？还是直接切入到具体的第一个核心模块**“故障场景与意图理解分析（Story 5）”**的设计讨论？请指示。

### User Input

按照当前的实现设计文档，经验知识应该拆分成几个object、几个action、几个relation？为什么？

### Planner Response

这是一个非常专业且切中核心的问题。

在“运行知识（Running Knowledge）”中，我们的 Object / Action / Relation 对应的是**物理和软件世界的实体**（如：MySQL实例对象、查询存储池动作、实例包含据库关系）。

而在“经验知识（Experience Knowledge）”中，我们抽象的是**资深 DBA 的认知域（Cognitive Domain）和决策网络**。这就要求我们把“思考过程”、“策略经验”和“验收标准”具象化为 BKN 的各类元素。

基于《实现设计文档》中规划的 **智能决策模型**（意图、问题、资源、风险、策略）和 **验证标准体系**，我建议将经验知识按照 **3 个 Object、3 个 Action、4 个 Relation** （并配合 Concept 和 Risk补充）来进行架构拆分。

以下是具体的拆分方案方案及深度产品逻辑剖析（为什么这么拆）：

---

### 一、 Object (对象) 拆分：抽象“思维实体”
经验知识的 Object 应该承载大模型的推理产物，将模糊的上下文转化为结构化的上下文状态。建议拆分为以下 **3 个核心思维对象**：

#### 1. `FaultScenario` (故障场景对象)
*   **承载域：** 意图理解、问题分析
*   **属性示例：** 故障根因分类（逻辑错误/硬件故障/勒索病毒）、影响范围（实例级/库级/表级）、RTO/RPO 容忍度评级。
*   **为什么拆分它：** 用户输入往往是非结构化的自然语言（如：“我们财务库误删了一张表”）。我们需要一个 Object 来接住 LLM “意图理解”的结果。当 Agent 识别出这是一个 `库级局部逻辑故障` 的 `FaultScenario` 时，它才能在知识图谱中根据关系去寻找对应的策略。

#### 2. `RecoveryStrategyBlueprint` (恢复策略蓝图/模板对象)
*   **承载域：** 策略生成
*   **属性示例：** 策略优先级、前置依赖矩阵、涉及的最佳恢复粒度、最佳恢复目的地组合。
*   **为什么拆分它：** 面对同一个故障，不同资深 DBA 会有不同解法。策略对象是连接“意图”和“底层执行任务（RecoveryTask）”的桥梁。将策略实例化为对象，Agent 就可以在对话中向用户呈现多个候选策略对象（如：A. 时间最优策略，B. 数据最安全策略）以供“审批（Approval）”。

#### 3. `VerificationStandard` (验证标准集对象)
*   **承载域：** 验证标准体系
*   **属性示例：** 验收指标清单（端口连通性/进程存活度/数据行数 Checksum比对/特定业务 SQL 探活）。
*   **为什么拆分它：** 没有验证的恢复等同于盲人摸象。将验证标准沉淀为独立的 Object，意味着它可以被复用。一个 `数据库联机验证标准集` 不仅可以用于故障恢复的最后一步，甚至可以被用在日常的演练系统中。

---

### 二、 Action (行动) 拆分：抽象“推演与验证工具”
对于经验知识，Action 不再是简单的“下发备份”“挂载存储”，而是 Agent 为了做出决策和兜底工作所需要执行的**分析型/校验型动作（Tools）**。建议拆分为以下 **3 个行动**：

#### 1. `diagnose_fault_environment` (环境侧辅助诊断)
*   **对应阶段：** 问题分析
*   **为什么拆分它：** 单凭用户的话往往无法准确定位问题。Agent 可能会为了完善 `FaultScenario` 对象，需要调用特定工具（例如执行一段探测 SQL，查看 MySQL 错误日志或特定表空间状态）来辅助 LLM 确认真正的故障影响面。

#### 2. `evaluate_strategy_resources` (策略资源连通性与容量验算)
*   **对应阶段：** 资源分析、策略生成
*   **为什么拆分它：** 所谓“经验”，就是不盲目执行。当 LLM 推理出一个 `异机全量恢复` 策略后，在正式下发任务前，Agent 需要一个经验知识类的 Action 去算账：目标异机剩余磁盘空间是否大于备份副本大小？端口是否放通？预评估不通过，则驱动 Agent 重新选择策略，形成自我纠正（Self-Healing / Reflection）。

#### 3. `execute_integrity_verification` (执行一致性与业务验证)
*   **对应阶段：** 观察反馈 & 验证标准
*   **为什么拆分它：** Foundation 底层的任务完成（状态变更为 successful）不代表业务恢复。经验知识层的特殊 Action 负责调用 `VerificationStandard` 中的配置（例如发起特定的 checksum 工具包），给出最终的“业务恢复成功”的裁判结论。

---

### 三、 Relation (关系) 拆分：织就“因果网络”
关系是将散落的“场景”、“策略”、“风险”联系在一起的血管，指导 Agent 如何思考和检索知识。我建议增加以下 **4 条核心关系脉络**：

#### 1. `FaultScenario` --[resolved_by]--> `RecoveryStrategyBlueprint`
*   **描述：** 某个故障场景可以用哪些策略蓝图解决。
*   **为什么：** 思考起点的基石。例如：`单表数据页损坏` 场景指向了 `异机全库恢复再导表还原` 这个策略蓝图。这是基于经验沉淀下来的因果对应关系。

#### 2. `RecoveryStrategyBlueprint` --[requires_asset]--> `RunningKnowledge (Object)`
*   **描述：** 特定的经验策略，必须依赖特定的运行态资产（如要求必须存在 `Binlog备份源`，或者必须要求 `充足的异机存储实例`）。
*   **为什么：** 将高维度的经验知识与底层的物理运行知识强绑定。保证策略落地可行。

#### 3. `RecoveryStrategyBlueprint` --[mitigates/triggers_risk]--> `Risk` (风险实体)
*   **描述：** 选择某种恢复策略会触发或者规避什么风险。
*   **为什么：** 对应文档中原机恢复覆盖业务的 `Risk` 模型。Agent 顺着这条关系线，可以对比向用户宣读不同策略方案的风险优劣，增加 CoT（思维链）的话术信服力和安全门槛。

#### 4. `RecoveryStrategyBlueprint` --[verified_by]--> `VerificationStandard`
*   **描述：** 如果执行了这个策略，应该运用什么验证标准去验收。
*   **为什么：** 不同的策略产生的影响面不同：原机最新状态恢复，重点验证数据新旧；异机挂载恢复，重点验证表空间提取的完整性。此联系将策略和验收动作牢牢锁定。

---

### 总结：整体架构协同

这个拆分设计不仅严格遵循了你制定的 5 阶段模型（意图→信息→规划→执行→反馈），且将 Agent 的 **CoT 思考要素进行了结构化**。

另外，不得不提一下规范中原有的 **`Concept(概念)`** 和 **`Risk(风险)`** 节点，它们简直是为了经验知识量身定制的：
*   我们应把文档中定义的“恢复方式、恢复粒度”沉淀为 `Concept` 节点，被策略蓝图组合调用。
*   除了《风险模型设计》中已有的“数据覆盖风险”，我们将在这个经验大阶段注入更多的“业务风险”，完善 `Risk` 模型。

这份拆包蓝图是否符合您对现有架构演进的预期？如果方向确认OK，我们可以先着手针对 **“FaultScenario (故障场景)”** 的具体 Object 设计结构进行攻坚。

### User Input

按照当前的实现设计文档，经验知识应该拆分成几个object、几个action、几个relation？为什么？

### Planner Response

这个问题问得非常好，它直指经验知识设计的本质分歧：**经验知识到底是"新的实体"还是"对已有实体的推理逻辑"？**

让我回到文档本身，逐层拆解。

---

## 先看文档中经验知识的定义

《实现设计文档》第 2.2.2 节明确定义了经验知识的 8 个能力域：

**智能决策模型（5 个）**：意图理解、问题分析、资源分析、风险评估、策略生成

**验证标准体系（3 个）**：策略验证、结果验证、完整性验证

同时，文档第 3.4 节规划的目录结构把它们全部放在了 `experience/` 目录下，**并没有归入 `object/`、`act/`、`relation/` 目录**：

```
experience/
├── intent_understanding.md
├── problem_analysis.md
├── resource_analysis.md
├── risk_assessment.md
├── strategy_generation.md
└── verification_standard.md
```

这暴露了一个关键问题：**当前设计文档对经验知识用 Object/Action/Relation 体系建模这件事，实际上是留白的。**

---

## 核心判断：经验知识的本质是什么？

回顾运行知识的 Object/Action/Relation 体系：

| BKN 元素 | 运行知识中的含义 | 对应的是什么？ |
|----------|----------------|-------------|
| Object | 系统中的**持久化实体** | MySQL 实例、备份副本、客户端 —— 有 ID、可查询、有状态 |
| Action | 对实体的**工具操作** | 调 Foundation API 查询/下发恢复 —— 绑定 toolbox、有参数 |
| Relation | 实体间的**结构关系** | 实例包含数据库、副本属于存储池 —— 对象间 1:N / N:N |

经验知识本质上是：**Agent 基于已有的运行知识实体，进行诊断推理和策略决策的认知框架。**

核心问题是：
-  "故障场景"是一个新的 Object 吗？—— **不是**，它是对已有 ProtectionObject 状态的一种**分类判断**。
- "策略生成"是一个新的 Action 吗？—— **不完全是**，策略生成的核心逻辑是 LLM 推理，而非调用一个外部工具。
- "故障→策略的对应"是一个新的 Relation 吗？—— **这倒是**，它确实描述了两类知识之间的结构化映射。

---

## 我的拆分建议

### Object：**0 个新增**

**为什么不需要新 Object？**

经验知识推理的**所有实体**都已经存在于运行知识中：
- 用户说"OrderDB 坏了" → 推理对象是已有的 `ProtectionObject`
- Agent 评估备份可用性 → 评估对象是已有的 `TimePointCopy`
- Agent 选择恢复目标 → 选择对象是已有的 `Client`
- Agent 生成策略 → 最终填充的是已有的 `RecoveryTask` 的参数

如果把"故障场景"、"恢复策略"做成 Object，它们既没有独立的数据源（`data_view`），也不是系统中的持久化实体，生硬建模会导致 Object 概念膨胀、边界模糊。

**那这些"认知产物"放哪？**

| 认知域 | 正确的 BKN 归属 | 原因 |
|--------|----------------|------|
| 故障场景分类（数据页损坏/误删库/勒索…） | **Concept** | 是分类枚举，不是实体 |
| 恢复策略模板（异机全库恢复导表、原机最新恢复…） | **Concept** | 是策略模式的取值空间，不是持久化对象 |
| 风险场景扩展（版本不兼容、空间不足…） | **Risk** | 运行知识已有 Risk 模型，经验知识扩展风险项 |
| 验证标准（数据行数比对、端口连通性检查…） | 已有的 `AvailabilityVerification` Object + **Concept** 扩展 | 运行知识已定义了 AvailabilityVerification 对象 |

---

### Action：**2～3 个新增**

经验知识中有两类逻辑：
1. **LLM 推理逻辑**（意图理解、问题分析、策略生成）—— 这些**不是 BKN Action**，是 System Prompt / Dolphin 中的指令编排
2. **需要调工具获取信息的操作**（资源探测、实际验证）—— 这些**是 BKN Action**，因为它们需要绑定 toolbox 执行实际操作

| 新增 Action | 对应经验知识域 | 为什么需要 |
|------------|-------------|-----------|
| `check_target_environment` | 资源分析 | Agent 需要工具检查目标主机的磁盘空间、MySQL 端口状态、进程状况，才能判定策略的可行性。不调工具纯靠 LLM 猜，不可靠。 |
| `verify_recovery_integrity` | 恢复结果验证 + 完整性验证 | Foundation 报 success 只代表数据写入完毕，不代表业务可用。Agent 需要调工具（如执行探活 SQL、比对行数 checksum）来给出最终裁判。 |
| `diagnose_fault_indicators`（可选） | 问题分析 | 当用户描述模糊时，Agent 可能需要调工具查看 MySQL 错误日志或 `SHOW TABLE STATUS`，来辅助确认故障根因和影响范围。 |

**为什么意图理解和策略生成不是 Action？**

对照运行知识的 Action 定义模式（绑定 toolbox → 调 Foundation API → 返回结构化数据），意图理解和策略生成的核心是 **LLM 基于 Concept 和 Relation 做的推理**，不是一次工具调用。它们应该编码在 **System Prompt（角色指令 + Dolphin）** 中，由 Concept 和 Relation 提供推理的知识基底。

---

### Relation：**2～3 个新增**

| 新增 Relation | 描述 | 为什么需要 |
|--------------|------|-----------|
| `fault_recommends_strategy` | 故障场景 Concept →推荐→ 策略模板 Concept | **这是经验知识的核心脉络**：当 Agent 识别出"单表数据页损坏"这个故障类型后，它需要一条关系线来找到"异机全库恢复+导表还原"这个推荐策略。没有这条关系，Agent 只能靠 LLM 幻觉去猜策略。 |
| `strategy_triggers_risk` | 策略模板 Concept →触发→ Risk | Agent 选定一个策略后，需要沿关系查看这个策略会触发什么风险（原机恢复→数据覆盖风险 high；异机恢复→low），这为 CoT 中的风险提示提供结构化依据。 |
| `strategy_verified_by` | 策略模板 Concept →验收→ 验证标准 Concept | 不同策略的验收重点不同。原机最新恢复 → 验回放完整性；异机导表 → 验目标表行数与源库一致。 |

---

### 总结对比

| BKN 元素 | 运行知识（已完成） | 经验知识（待设计） | 差异原因 |
|----------|-----------------|-----------------|---------|
| **Object** | 9 个（ProtectionObject, Client 等） | **0 个** | 经验知识不产生新实体，而是推理已有实体 |
| **Action** | 8+ 个（query_*, execute_* 等） | **2～3 个** | 只有需要调工具的操作才是 Action；推理逻辑走 Prompt |
| **Relation** | 10 个（instance_has_db 等） | **2～3 个** | 经验知识关系连接的是 Concept 之间的因果映射 |
| **Concept** | 待设计 | **3～4 个**（故障分类、策略模板、验证标准集） | 经验知识的主力建模手段 |
| **Risk** | 1 个（recovery_overwrite_risk） | **扩展 2～3 个** | 经验域补充版本兼容、空间不足等风险规则 |

**一句话结论**：经验知识的 BKN 建模，**重 Concept、轻 Object、精 Action**。它的核心不在于引入新实体，而在于用 Concept 构建故障→策略→验证的决策知识图谱，用少量诊断/验证 Action 弥补 LLM 推理的信息盲区，最终由 System Prompt 串联整个推理流程。

---

需要我接着深入设计第一个模块——**故障场景 Concept 分类体系**（所有策略推理的起点）吗？

### User Input

1.数据完整性验证，这个需要修改，不是验证数据的完整性，是验证恢复后的应用能否对外提供业务服务，比如mysql恢复完能否use mysql，然后查询表数据能否查出来；2.经验业务知识网络，设计理念是参考palantir的ontology思想设计的。

### User Input

Continue

### User Input

Continue

### User Input

Continue

### Planner Response

收到，这两个信息非常关键，让我重新审视整个设计思路。

---

## 修正 1：验证标准的正确定义

你说得对。"数据完整性验证"这个提法容易让人理解成 checksum 比对，但实际需求是：

> **恢复后的应用能否对外提供业务服务？**

对 MySQL 来说就是：
1. 进程活了没？`mysqld` 是否正常监听端口
2. 能不能连上？`mysql -u root -p` 能否建立连接
3. 业务数据能不能用？`USE OrderDB` → `SELECT * FROM orders LIMIT 1` 能否返回数据

这本质上是**业务可用性验证**，而运行知识中已有的 `AvailabilityVerification` 对象（有 `appType`、`tool`、`params`、`description`）恰好就是承载这个能力的载体。经验知识要做的是：**针对不同的故障场景和恢复策略，配置不同的验证标准组合**，而不是新造一个验证实体。

---

## 修正 2：用 Palantir Ontology 思想重新审视

这个提示彻底改变了架构视角。

Palantir Ontology 的核心理念是：**把所有业务相关的概念都建模为一等公民 Object，通过 typed Link（关系）构建语义网络，通过 Action 定义对 Object 的操作。** 一切皆对象，对象皆可查询、可链接、可操作。

这意味着我之前"0 个新 Object"的判断是错的。在 Ontology 思想下：

- "误删库"不是一个模糊的 Concept 枚举值，它是一个**故障场景 Object 实例**，有自己的属性（典型症状、影响范围、紧急度）和到策略 Object 的链接。
- "异机全库恢复+导表还原"不是 Prompt 中的一段文字，它是一个**恢复策略 Object 实例**，有前置条件、执行步骤序列、和到验证标准的链接。

**Ontology 的价值在于：Agent 不是靠 LLM 幻觉去猜策略，而是沿着 Object → Link → Object 在知识图谱上做结构化检索和推理。**

---

## 重新拆分：基于 Palantir Ontology 的经验知识架构

### Object：3 个新增

| Object | 语义 | 关键属性 | 为什么在 Ontology 中它必须是 Object |
|--------|------|---------|----------------------------------|
| **FaultScenario** | 故障场景类型 | `faultCategory`（逻辑故障/硬件故障/安全攻击）、`affectedScope`（实例/库/表）、`typicalSymptoms`（错误日志特征）、`urgencyLevel`、`rtoRequirement` | 作为 Object，不同故障场景实例可以有不同的属性值，并通过 Link 指向不同的推荐策略。Agent 查询 `FaultScenario WHERE symptoms LIKE 'InnoDB 数据页损坏'` 就能找到对应对象，再沿关系找到策略。 |
| **RecoveryPlaybook** | 恢复策略剧本 | `strategyName`、`recommendedGranularity`、`recommendedDestination`、`recommendedMethod`、`executionSteps`（有序步骤描述）、`prerequisites` | 这是 Agent 的"经验手册"。每个 Playbook 都是一个可独立查询、可向用户展示的策略实例。Agent 可以在 CoT 中说："根据知识图谱，匹配到 Playbook-007：异机全库恢复→导表还原方案。" |
| **AvailabilityCheckpoint** | 业务可用性检查点 | `checkType`（进程检查/连接检查/业务查询检查）、`checkCommand`、`expectedResult`、`timeoutSeconds`、`appType` | 将验证步骤拆为独立的检查点 Object。一个策略可能关联多个检查点（先检查进程 → 再检查连接 → 再执行业务 SQL），每个检查点有明确的预期结果。这比把所有验证逻辑塞进一个 AvailabilityVerification.params JSON 中更清晰、更可扩展。 |

---

### Action：3 个新增

| Action | 绑定 Object | 类型 | 为什么需要 |
|--------|------------|------|-----------|
| **match_fault_scenario** | FaultScenario | query | Agent 将用户描述的症状（关键词/错误信息）作为输入，在 FaultScenario 知识库中匹配最相关的故障场景 Object，返回匹配的场景及其属性。这替代了纯 LLM 幻觉猜测。 |
| **evaluate_playbook_feasibility** | RecoveryPlaybook | query/check | 选定一个 Playbook 后，Agent 需要检查它的 `prerequisites` 是否满足（目标主机磁盘够不够、MySQL 版本兼容不兼容、备份副本是否可用）。这是"资源分析"能力的 Action 化。 |
| **execute_availability_check** | AvailabilityCheckpoint | execute | 恢复完成后，按照关联的 AvailabilityCheckpoint 逐项执行业务可用性验证（比如：`mysql -e "USE OrderDB; SELECT COUNT(*) FROM orders;"`），返回每个检查点的通过/失败结果。 |

---

### Relation：3 个新增

| Relation | 源 Object → 目标 Object | 基数 | 语义 |
|----------|------------------------|------|------|
| **fault_resolved_by_playbook** | FaultScenario → RecoveryPlaybook | 1:N | 一种故障场景可以被多种策略剧本解决（Agent 呈现多选）。这是经验知识图谱中最核心的一条边——把"问题空间"和"解法空间"连起来。 |
| **playbook_triggers_risk** | RecoveryPlaybook → Risk | N:N | 一个策略可能触发多个风险（原机恢复→覆盖风险+业务中断风险）；同一个风险可被多个策略触发。Agent 沿此关系为 CoT 输出风险分析。 |
| **playbook_verified_by_checkpoint** | RecoveryPlaybook → AvailabilityCheckpoint | 1:N | 一个策略对应一组有序的验证检查点。异机全库恢复和原机恢复的验收重点不同，各自挂不同的 checkpoint 集合。 |

---

### 整体知识图谱骨架

```
用户描述 → [match_fault_scenario] → FaultScenario
                                          │
                                    fault_resolved_by_playbook (1:N)
                                          │
                                          ▼
                                    RecoveryPlaybook ──playbook_triggers_risk──→ Risk
                                          │
                                          │ [evaluate_playbook_feasibility]
                                          │   检查 prerequisites vs 运行知识
                                          │   (TimePointCopy/Client/StoragePool)
                                          │
                                    playbook_verified_by_checkpoint (1:N)
                                          │
                                          ▼
                                    AvailabilityCheckpoint
                                          │
                                          │ [execute_availability_check]
                                          │   USE OrderDB; SELECT COUNT(*)...
                                          ▼
                                    验证结果 → 反馈给用户
```

---

### 与运行知识的衔接

经验知识的 RecoveryPlaybook 最终产出的是运行知识中 `RecoveryTask` 的参数集：

```
RecoveryPlaybook.recommendedGranularity  → RecoveryTask.recovery_granularity
RecoveryPlaybook.recommendedDestination  → RecoveryTask.recovery_destination
RecoveryPlaybook.recommendedMethod       → RecoveryTask.recovery_method
```

**Ontology 的完整链路**：`FaultScenario` → `RecoveryPlaybook` → `RecoveryTask` → `RecoveryJob` → `AvailabilityCheckpoint` → 业务恢复确认。

---

### 汇总

| BKN 元素 | 经验知识新增数量 | 具体内容 |
|----------|---------------|---------|
| **Object** | **3** | FaultScenario, RecoveryPlaybook, AvailabilityCheckpoint |
| **Action** | **3** | match_fault_scenario, evaluate_playbook_feasibility, execute_availability_check |
| **Relation** | **3** | fault_resolved_by_playbook, playbook_triggers_risk, playbook_verified_by_checkpoint |
| **Risk** | **扩展 2~3** | 版本不兼容风险、目标空间不足风险等（扩展已有 Risk 模型） |

要继续深入设计 **FaultScenario Object 的属性模型**吗？这是整条链路的起点。

### User Input

1.mysql_recovery_bkn_implementation文档不进行修改吗？还是数据完整性完整。2.当前的设计是我分为两个业务知识网络 - 运行知识网络(已设计并开发完)和经验知识网络，经验知识网络是指导Agent如何理解用户的故障场景业务意图->获取当前的信息如故障发生的时间？损坏的业务或者数据范围？备份foundation有哪些备份数据、哪些恢复技术？等->恢复策略，如何进行恢复的方案->下发执行，执行生成的恢复策略->观察与反馈，监测当前执行情况，如遇问题及时进行修复；运行知识网络是指导Agent如何使用备份foundation提供的能力工具来落实恢复问题的解决。这两种业务知识网络是上下层的关系。

### Planner Response

你说得非常清晰，我完全理解了。让我对两个问题逐一回应。

---

## 关于第 1 点：文档需要修改

实现设计文档中"数据完整性验证（DataIntegrityVerification）"的描述应该修正，改为体现**业务可用性验证**的真实含义。让我来更新：

*Viewed [mysql_recovery_bkn_implementation.md](file:///e:/code/KWeaver/kweaver-sdk/docs/request/mysql_recovery_bkn_implementation.md) *

### Planner Response

现在让我看一下架构文档中对应的部分，一并修改：

*Viewed [mysql_recovery_bkn_architecture.md](file:///e:/code/KWeaver/kweaver-sdk/docs/request/mysql_recovery_bkn_architecture.md) *

### Planner Response

我还需要检查架构文档中其他引用到"数据完整性验证"的地方：

*Grep searched codebase*

### Planner Response

现在修改实现设计文档中的 3 处：

*Edited relevant file*

### Planner Response

同步修改架构文档：

*Edited relevant file*

### Planner Response

两份文档都已修改完成。

---

## 关于第 2 点：经验知识网络与运行知识网络的上下层关系

你对两层网络的定位描述得非常精准，让我确认我完全理解了这个架构分层：

```
┌─────────────────────────────────────────────────┐
│          经验知识网络（上层 · 认知层）              │
│                                                   │
│  指导 Agent "怎么想"：                             │
│  意图理解 → 信息获取 → 策略规划 → 下发执行 → 观察反馈 │
│                                                   │
│  关注的是：                                        │
│  · 用户说了什么？→ 故障场景是什么？                   │
│  · 故障发生的时间？损坏的业务/数据范围？               │
│  · Foundation 有哪些备份数据、哪些恢复技术可用？       │
│  · 如何制定恢复方案？                               │
│  · 执行中有没有问题？需不需要修正？                   │
└─────────────────────────────────────────────────┘
                    ↓  调用（上层编排下层）
┌─────────────────────────────────────────────────┐
│          运行知识网络（下层 · 执行层）              │
│                                                   │
│  指导 Agent "怎么做"：                             │
│  · 如何查询 Foundation 保护对象列表                  │
│  · 如何查询时间点副本                               │
│  · 如何创建恢复任务                                 │
│  · 如何监控恢复作业状态                              │
│  · 如何执行可用性验证                               │
│                                                   │
│  已完成：9 个 Object + 8+ 个 Action + 10 个 Relation │
└─────────────────────────────────────────────────┘
```

**核心关系**：经验知识网络编排 Agent 的 5 阶段认知流程，在每个阶段需要**调用**运行知识网络的能力来实际执行操作。例如：

| 经验知识阶段 | 调用运行知识的方式 |
|------------|-----------------|
| 信息获取（备份副本有哪些？） | 调用运行知识的 `query_timepoint_copies` Action |
| 信息获取（哪些客户端可做恢复目标？） | 调用运行知识的 `query_clients` Action |
| 下发执行（执行恢复策略） | 调用运行知识的 `execute_recovery_task` Action |
| 观察反馈（监控任务进度） | 调用运行知识的 `query_recovery_jobs` Action |

经验知识网络本身**不直接绑定 Foundation 工具**，它的 Object / Action / Relation 建模的是**决策逻辑和认知结构**，最终通过链接到运行知识网络的 Object 和 Action 来落地执行。

---

## 这个上下层关系对之前拆分的影响

回到之前讨论的经验知识 3 Object / 3 Action / 3 Relation，基于上下层关系视角做一个定位修正：

| 经验知识元素 | 定位修正 |
|------------|---------|
| **FaultScenario Object** | 上层认知对象。Agent 在"意图理解"阶段产出的结构化诊断结果。 |
| **RecoveryPlaybook Object** | 上层认知对象。Agent 在"策略规划"阶段，基于 FaultScenario 匹配到的恢复方案模板。Playbook 的 `executionSteps` 中引用的是下层运行知识的 Action（如"第1步：调用 query_timepoint_copies 获取备份列表"）。 |
| **AvailabilityCheckpoint Object** | 上层认知对象。定义"观察反馈"阶段的验收标准。实际执行验证时调用下层运行知识的工具。 |
| **match_fault_scenario Action** | 上层编排动作。在知识图谱中检索匹配的故障场景。 |
| **evaluate_playbook_feasibility Action** | 上层编排动作。验证策略可行性——它会**向下调用**运行知识的查询 Action（如检查备份是否存在、客户端是否在线）来获取判断依据。 |
| **execute_availability_check Action** | 上层编排动作。执行业务可用性验证——它会**向下调用**运行知识中的工具来实际执行 SQL 探活等操作。 |

这样，两层网络形成了清晰的**编排-执行**分离：经验知识说"做什么、为什么"，运行知识说"怎么做、用什么工具"。

接下来要进入具体的 Object 属性设计了吗？我建议从 **FaultScenario** 开始——它是整条经验知识链路的入口。

### User Input

按照最新的理解。按照当前的实现设计文档，经验知识应该拆分成几个object、几个action、几个relation？为什么？


### Planner Response

基于你确认的两个关键前提——**上下层架构**和 **Palantir Ontology 思想**——我沿着经验知识网络承载的 5 阶段认知流程逐阶段推导。

---

## 推导逻辑：从 5 阶段工作流反推需要哪些 BKN 元素

经验知识网络编排 Agent 的认知流程是：

```
意图理解 → 信息获取 → 策略规划 → 下发执行 → 观察与反馈
```

但这 5 个阶段与运行知识网络的关系**不对等**：

| 阶段 | 谁在主导？ | 说明 |
|------|----------|------|
| ① 意图理解 | **经验知识** | Agent 需要结构化认知来识别用户说的是什么故障 |
| ② 信息获取 | **经验知识编排，运行知识执行** | 经验知识告诉 Agent "你需要去查备份列表、查客户端状态"，运行知识的 Action（query_timepoint_copies 等）实际执行 |
| ③ 策略规划 | **经验知识** | Agent 需要结构化的策略模板来匹配故障 → 恢复方案 |
| ④ 下发执行 | **经验知识编排，运行知识执行** | 经验知识将策略参数映射为运行知识的 RecoveryTask，运行知识的 execute_recovery_task 实际执行 |
| ⑤ 观察与反馈 | **经验知识编排，运行知识执行** | 经验知识定义"验什么算恢复成功"，运行知识执行工具实际验证 |

**关键发现**：5 个阶段中，②④ 本质上是"调用下层运行知识"，不需要新建经验知识层的 Object/Action。经验知识真正需要建模的核心认知阶段是 **①③⑤**——这三个阶段各自产出一个认知实体。

---

## 结论：3 Object、3 Action、3 Relation

### Object：3 个

每个核心认知阶段需要一个 Ontology Object 来承载其结构化产出：

| Object | 对应阶段 | 它是什么 | 为什么必须是 Object |
|--------|---------|---------|------------------|
| **FaultScenario**（故障场景） | ① 意图理解 | 故障类型模板。属性：故障分类（逻辑故障/硬件故障/安全攻击）、影响范围（实例/库/表）、典型症状特征、紧急度、RTO 要求 | Agent 需要一个结构化的"故障知识库"来匹配用户描述。按 Ontology 思想，它是可查询、可链接的一等公民实体。用户说"InnoDB 数据页损坏"，Agent 在 FaultScenario 知识图谱中做语义匹配，找到对应的故障对象实例，再沿关系找到策略——而不是靠 LLM 幻觉猜测。 |
| **RecoveryPlaybook**（恢复策略剧本） | ③ 策略规划 | 策略方案模板。属性：推荐的恢复粒度、恢复目的地、恢复方式、前置条件列表、执行步骤序列说明 | 这是经验知识的核心——资深 DBA 的"解题套路"。一个 FaultScenario 可能关联多个 Playbook（比如"误删库"既可以原机恢复也可以异机导出），Agent 在 CoT 中向用户展示候选方案。Playbook 的 executionSteps 引用的是下层运行知识的 Action，实现上下层贯通。 |
| **AvailabilityCheckpoint**（业务可用性检查点） | ⑤ 观察与反馈 | 验证检查项。属性：检查类型（进程检查/连接检查/业务 SQL 检查）、检查命令、预期结果、超时时间、适用的 appType | 不同恢复策略的验收重点不同。将检查项拆成独立的 Object，一个 Playbook 可关联多个有序的 Checkpoint（先查进程 → 再查连接 → 最后执行 `USE OrderDB; SELECT COUNT(*) FROM orders`），每个 Checkpoint 有明确的预期结果判定。 |

**为什么信息获取（②）和下发执行（④）不需要新 Object？**

因为这两个阶段操作的实体（ProtectionObject、TimePointCopy、Client、RecoveryTask、RecoveryJob）全部已在运行知识网络中定义。经验知识只需要"知道去找谁"，不需要重复建模。

---

### Action：3 个

每个核心 Object 需要一个对应的操作动作：

| Action | 绑定 Object | 对应阶段 | 做什么 |
|--------|------------|---------|-------|
| **match_fault_scenario** | FaultScenario | ① 意图理解 | 将用户描述的症状关键词/错误信息输入，在故障场景知识库中匹配最相关的 FaultScenario 实例并返回。这替代了纯 LLM 猜测，让故障识别建立在结构化知识上。 |
| **generate_recovery_strategy** | RecoveryPlaybook | ③ 策略规划 | 基于匹配到的 FaultScenario，沿 `fault_resolved_by_playbook` 关系找到候选 Playbook，然后结合阶段②获取的实际资源状况（备份是否有效、目标主机是否可用），评估每个 Playbook 的前置条件是否满足，筛选出可行的策略方案。**资源分析和风险评估在此 Action 内完成**。 |
| **execute_availability_check** | AvailabilityCheckpoint | ⑤ 观察与反馈 | 恢复完毕后，按照 Playbook 关联的 AvailabilityCheckpoint 逐项执行业务可用性验证。向下调用运行知识的工具（如执行 `mysql -e "USE OrderDB; SELECT COUNT(*) FROM orders;"`），返回每个检查点的通过/失败结果。 |

**为什么"意图理解"中的问题分析不单独拆 Action？**

问题分析（故障根因判断、影响范围评估）和意图理解是**同一个认知动作的两个维度**——都是在回答"发生了什么"。拆成两个 Action 会导致 Agent 需要串行执行两次高度耦合的推理，不如在 `match_fault_scenario` 中一次完成（FaultScenario Object 的属性已包含根因分类和影响范围）。

**为什么"资源分析"和"风险评估"没有成为独立 Action？**

- **资源分析**的实际操作（查备份列表、查客户端状态）是调用下层运行知识的 Action，不需要经验知识层新建 Action
- **风险评估**通过 `playbook_triggers_risk` 关系实现——Agent 选定 Playbook 后沿关系查到关联的 Risk，这是图查询而非独立 Action
- 两者在 `generate_recovery_strategy` 中作为策略可行性判断的子步骤内联完成

---

### Relation：3 个

3 条关系构成经验知识网络的核心因果图谱：

| Relation | 源 → 目标 | 基数 | 为什么需要 |
|----------|----------|------|-----------|
| **fault_resolved_by_playbook** | FaultScenario → RecoveryPlaybook | 1:N | **经验知识最核心的一条边**。"单表数据页损坏"指向"异机全库恢复→导表还原"和"原机最新恢复"两个 Playbook。Agent 沿此关系做检索，而非凭空生成策略。N 个候选方案让用户在审批环节有选择权。 |
| **playbook_triggers_risk** | RecoveryPlaybook → Risk | N:N | Agent 选定 Playbook 后，沿此关系查出该策略触发的风险（原机恢复→数据覆盖风险 high；异机恢复→low）。直接复用运行知识已有的 Risk 模型，经验知识扩展更多风险场景。这条关系为 CoT 的风险分析话术提供结构化数据。 |
| **playbook_verified_by_checkpoint** | RecoveryPlaybook → AvailabilityCheckpoint | 1:N | 不同策略的验收重点不同：异机全库恢复重点验"数据导出是否完整"，原机最新恢复重点验"Binlog 回放是否到位"。有序的 Checkpoint 列表保证了验证的完整性和可追溯性。 |

---

### 上下层贯通的全链路视图

```
用户："OrderDB 数据页损坏，order_details 表查不了了"
│
│ ①意图理解
│  [match_fault_scenario]
│         ↓
│  FaultScenario: 库级故障·InnoDB数据页损坏
│         │
│         │ fault_resolved_by_playbook (1:N)
│         ↓
│ ②信息获取（经验知识编排 → 调用运行知识）
│  → query_timepoint_copies    ← 运行知识 Action
│  → query_clients             ← 运行知识 Action
│  → query_storage_pools       ← 运行知识 Action
│         ↓
│ ③策略规划
│  [generate_recovery_strategy]
│         ↓
│  RecoveryPlaybook-A: 异机全库恢复→导表还原（推荐）
│  RecoveryPlaybook-B: 原机最新恢复（备选）
│         │
│         │ playbook_triggers_risk (N:N)
│         │  → Playbook-A → Risk: low（异机，不覆盖生产）
│         │  → Playbook-B → Risk: high（原机，覆盖生产数据）
│         ↓
│  用户审批 → 选择 Playbook-A
│         ↓
│ ④下发执行（经验知识编排 → 调用运行知识）
│  → 经验知识将 Playbook-A 参数映射为 RecoveryTask
│  → execute_recovery_task     ← 运行知识 Action
│  → query_recovery_jobs       ← 运行知识 Action（监控进度）
│         ↓
│ ⑤观察与反馈
│         │ playbook_verified_by_checkpoint (1:N)
│         ↓
│  [execute_availability_check]
│  → Checkpoint-1: mysqld 进程是否存活      ✅
│  → Checkpoint-2: 能否建立连接              ✅
│  → Checkpoint-3: USE OrderDB; SELECT ... ✅
│         ↓
│  反馈报告：业务恢复成功
```

---

### 汇总

| BKN 元素 | 数量 | 具体内容 | 设计依据 |
|----------|------|---------|---------|
| **Object** | **3** | FaultScenario, RecoveryPlaybook, AvailabilityCheckpoint | 5 阶段中 ①③⑤ 三个核心认知阶段各产出一个 Ontology 实体；②④ 操作的是运行知识已有实体 |
| **Action** | **3** | match_fault_scenario, generate_recovery_strategy, execute_availability_check | 每个核心 Object 对应一个认知操作；资源分析/风险评估内联在策略生成中，不单独拆 |
| **Relation** | **3** | fault_resolved_by_playbook, playbook_triggers_risk, playbook_verified_by_checkpoint | 构成 故障→策略→验证 的因果链路，风险通过关系挂载而非独立推理 |

这个 **3-3-3** 结构的好处是：对齐 5 阶段工作流、上下层职责清晰、每个元素都有明确的 Ontology 存在理由。方向确认后，我们可以从 FaultScenario 的属性模型开始攻坚。

### User Input

信息获取不需要经验知识层吗？要分析故障发生的时间、粒度和where，以及备份Foundation的备份数据情况和恢复能力如是否支持细粒度恢复？是否支持挂载恢复？，这些获取的信息是后续规划策略的重要信息。

### User Input

Continue

### Planner Response

你说得对，我之前的分析把阶段②简化过度了。

## 信息获取为什么需要经验知识层

我之前把信息获取阶段定义为"只是调运行知识查数据"，这是错的。真正的信息获取需要经验知识来回答一个关键问题：**Agent 该去获取什么信息？**

一个初级工程师和一个资深 DBA 面对同一个故障，差距不在于"会不会调 API"，而在于：**资深 DBA 知道要从哪些维度去收集信息，而且知道收集来的信息意味着什么。**

具体来说你列出的信息维度：

| 信息维度 | 具体内容 | 为什么是经验知识而不只是调 API |
|---------|---------|---------------------------|
| **WHEN**（时间） | 故障发生的时间 | 决定了要找哪个时间点之前的备份，决定了 Binlog 回放的终点 |
| **WHAT**（粒度） | 影响的是实例级、库级还是表级 | 直接约束后续策略的恢复粒度选择 |
| **WHERE**（位置） | 具体是哪个实例、哪个库、哪张表 | 确定要查哪个 ProtectionObject 的备份 |
| **备份数据情况** | Foundation 有哪些备份时间点可用 | 不只是列表查询，还要判断"是否有故障时间点之前的干净备份" |
| **恢复能力** | Foundation 是否支持细粒度恢复？是否支持挂载恢复？ | 这直接决定了哪些 Playbook 是可行的——如果不支持挂载恢复，那"挂载恢复+导表"的策略就不可用 |

**这些维度的定义、获取顺序、以及如何解读结果，都是经验知识。** 运行知识只负责"查到数据"，经验知识负责"知道要查什么、怎么解读"。

---

## 修正后的拆分：4 Object、4 Action、4 Relation

### Object：4 个

新增 **RecoveryContext**（恢复上下文评估），对应阶段②信息获取的结构化产出：

| Object | 对应阶段 | 它是什么 |
|--------|---------|---------|
| **FaultScenario** | ① 意图理解 | 故障场景模板。属性：故障分类、影响范围、典型症状、紧急度 |
| **RecoveryContext** | ② 信息获取 | 恢复上下文评估。属性：故障时间、故障粒度（实例/库/表）、受影响目标、可用备份时间点列表、Foundation 可用恢复能力（是否支持挂载恢复、是否支持细粒度恢复、是否支持 Binlog 回放） |
| **RecoveryPlaybook** | ③ 策略规划 | 恢复策略剧本。属性：推荐恢复粒度、恢复目的地、恢复方式、前置条件、执行步骤序列 |
| **AvailabilityCheckpoint** | ⑤ 观察与反馈 | 业务可用性检查点。属性：检查类型、检查命令、预期结果、超时时间 |

**为什么 RecoveryContext 必须是独立 Object：**

RecoveryContext 是阶段②的**结构化产出物**，也是阶段③策略生成的**核心输入**。如果不把它建模为 Object：
- Agent 在策略规划阶段只能基于 FaultScenario（"什么类型的故障"）去选策略，但不知道"当前实际有哪些资源可用"
- RecoveryPlaybook 的前置条件（如"需要全量备份"、"需要支持挂载恢复"）无法与实际情况进行结构化比对
- 等于缺了**从问题空间到解法空间的约束条件**

---

### Action：4 个

| Action | 绑定 Object | 对应阶段 | 做什么 |
|--------|------------|---------|-------|
| **match_fault_scenario** | FaultScenario | ① 意图理解 | 匹配用户描述到结构化的故障场景 |
| **assess_recovery_context** | RecoveryContext | ② 信息获取 | 按维度系统收集恢复上下文。具体：1) 向用户确认故障时间/粒度/位置；2) 调用运行知识查询可用备份（query_timepoint_copies）；3) 评估 Foundation 对当前备份类型支持的恢复能力集（挂载？细粒度？Binlog？）。产出一个完整的 RecoveryContext 实例 |
| **generate_recovery_strategy** | RecoveryPlaybook | ③ 策略规划 | 基于 FaultScenario（沿关系找到候选 Playbook）+ RecoveryContext（用实际资源情况过滤可行的 Playbook），生成可执行的策略方案 |
| **execute_availability_check** | AvailabilityCheckpoint | ⑤ 观察与反馈 | 恢复完毕后逐项执行业务可用性检查 |

**assess_recovery_context 的关键价值**：它不是简单的"查一下 API"，而是一个**结构化的评估流程**——知道要查哪些维度、以什么顺序查、查到的结果对策略选择意味着什么。比如：
- 查到 Foundation 有挂载恢复能力 → RecoveryContext.capabilities 标记 `mount_recovery = true` → 后续策略规划中"挂载恢复+导表"的 Playbook 才被视为可行
- 查到最近的全量备份是 3 天前 → RecoveryContext 记录这个时间差 → 策略规划中需要评估 3 天的 Binlog 回放量是否可接受

---

### Relation：4 个

| Relation | 源 → 目标 | 基数 | 语义 |
|----------|----------|------|------|
| **scenario_requires_context** | FaultScenario → RecoveryContext | 1:1 | 每种故障场景定义了"需要收集哪些维度的信息"。比如：库级故障需要确认具体哪个库；实例级故障需要确认整个实例状态。Agent 匹配到 FaultScenario 后，沿此关系知道下一步该收集什么。 |
| **fault_resolved_by_playbook** | FaultScenario → RecoveryPlaybook | 1:N | 故障场景指向多个候选策略剧本。这是经验知识图谱的核心骨架。 |
| **playbook_triggers_risk** | RecoveryPlaybook → Risk | N:N | 策略触发的风险。原机恢复→覆盖风险 high；异机恢复→low。 |
| **playbook_verified_by_checkpoint** | RecoveryPlaybook → AvailabilityCheckpoint | 1:N | 策略关联的有序验证检查点列表。 |

**关于 RecoveryContext 与 RecoveryPlaybook 的关系**：不需要静态 Relation。RecoveryContext 是 `generate_recovery_strategy` Action 的**运行时输入**——Action 拿 FaultScenario（沿 fault_resolved_by_playbook 获取候选 Playbook 列表）+ RecoveryContext（实际资源状况），动态判断每个 Playbook 的 prerequisites 是否满足。这是 Action 内的逻辑，不是静态的图关系。

---

### 完整链路（修正版）

```
用户："OrderDB 数据页损坏，order_details 表查不了了"
│
│ ① 意图理解
│  [match_fault_scenario]
│         ↓
│  FaultScenario: 库级故障·InnoDB数据页损坏
│         │
│         │ scenario_requires_context
│         ↓
│ ② 信息获取
│  [assess_recovery_context]
│   → 向用户确认：故障发生时间约 10:30
│   → 确认粒度：库级（OrderDB）
│   → 调用运行知识 query_timepoint_copies → 最近全量备份 04-15 02:00
│   → 评估 Foundation 能力：支持挂载恢复 ✅、支持 Binlog 回放 ✅
│         ↓
│  RecoveryContext: {
│    faultTime: "2026-04-15 10:30",
│    granularity: "database",
│    target: "OrderDB",
│    latestFullBackup: "2026-04-15 02:00",
│    capabilities: { mount: true, binlog: true, fineGrained: false }
│  }
│         ↓
│ ③ 策略规划
│  [generate_recovery_strategy]
│   → 沿 fault_resolved_by_playbook 获取候选 Playbook
│   → 用 RecoveryContext 过滤可行方案
│         ↓
│  Playbook-A: 异机挂载恢复→导表还原→Binlog回放（推荐）
│     playbook_triggers_risk → Risk: low
│  Playbook-B: 原机全量恢复→Binlog回放
│     playbook_triggers_risk → Risk: high
│         ↓
│  用户审批 → 选择 Playbook-A
│         ↓
│ ④ 下发执行（调用运行知识 execute_recovery_task）
│         ↓
│ ⑤ 观察与反馈
│  playbook_verified_by_checkpoint → [1, 2, 3]
│  [execute_availability_check]
│   → Checkpoint-1: mysqld 进程存活     ✅
│   → Checkpoint-2: mysql 连接建立      ✅
│   → Checkpoint-3: USE OrderDB;
│       SELECT COUNT(*) FROM orders;   ✅ 返回 15832 行
│         ↓
│  反馈：业务恢复成功
```

---

### 汇总：4-4-4

| BKN 元素 | 数量 | 具体内容 | 对应5阶段 |
|----------|------|---------|----------|
| **Object** | **4** | FaultScenario, RecoveryContext, RecoveryPlaybook, AvailabilityCheckpoint | ①②③⑤ |
| **Action** | **4** | match_fault_scenario, assess_recovery_context, generate_recovery_strategy, execute_availability_check | ①②③⑤ |
| **Relation** | **4** | scenario_requires_context, fault_resolved_by_playbook, playbook_triggers_risk, playbook_verified_by_checkpoint | 串联全流程 |

阶段④（下发执行）完全由经验知识编排、运行知识执行，不产生新的经验知识层 BKN 元素。

方向OK的话，我们从哪个 Object 开始细化属性模型？

### User Input

下发执行，是Agent执行到这个workflow这一步的时候，查阅运行知识网络中执行恢复任务和执行可用性验证的知识，并通知Agent调用运行知识网络中对应的action工具来执行。这个不需要经验知识层吗？你仔细分析下。

### Planner Response

你说得对，我再仔细分析一下。

## 阶段④不只是"调用运行知识"

我之前犯的错误是：把"下发执行"等同于"直接调用运行知识的 execute_recovery_task"。但实际上，Agent 在这一步面临的问题是：

> **我有了一个 RecoveryPlaybook（策略），我有了一个 RecoveryContext（上下文），但我怎么把它们翻译成运行知识网络中具体的 Action 调用？**

这就是你说的关键点——Agent 需要：

1. **查阅运行知识网络** → "运行知识有哪些执行类 Action？各自需要什么参数？"
2. **建立映射** → "Playbook 说'做库级异机恢复'，对应到运行知识的哪个 Action？参数怎么填？"
3. **编排顺序** → "先创建恢复任务、再监控作业状态、中间如果失败怎么处理？"
4. **通知 Agent 执行** → "指导 Agent 按步骤调用运行知识的工具"

这些**映射和编排逻辑本身就是经验知识**。一个初级工程师看到"库级异机恢复"策略，不知道：
- RecoveryTask 的 `recovery_granularity` 该填 `database` 还是 `instance`
- 要不要先 `stop_service`
- `recovery_type` 该选 `data` 还是 `mount`
- 恢复完第一步后下一步该查什么

资深 DBA 知道这些——**这就是经验**。

---

## 所以阶段④也需要独立的 Object + Action

| 维度 | 分析 |
|------|------|
| **知道该调什么** | Playbook.executionSteps 只描述高层步骤（如"Step1: 执行异机恢复"），经验知识需要把它映射为具体的运行知识 Action（execute_recovery_task） |
| **知道参数怎么填** | Playbook 的 recommendedGranularity="database" + RecoveryContext 的 target="OrderDB" → 映射为 RecoveryTask 的具体参数集 |
| **知道如何编排** | 恢复完成后要主动查询作业状态（query_recovery_jobs），而不是傻等；失败了要判断原因并决定重试还是切换策略 |
| **知道异常怎么处理** | 执行中遇到"目标路径空间不足"，经验知识指导 Agent 重新评估或切换 Playbook |

这催生了一个新 Object：**ExecutionPlan**（执行计划）——它是 RecoveryPlaybook + RecoveryContext 的**具象化产物**，把高层策略翻译成具体的、可逐步执行的运行知识 Action 调用序列。

---

## 修正后的拆分：5 Object、5 Action、5 Relation

每个阶段恰好对应一个 Object + 一个 Action：

### Object：5 个

| Object | 对应阶段 | 承载什么认知 |
|--------|---------|------------|
| **FaultScenario** | ① 意图理解 | 故障模式识别（"这是什么类型的问题"） |
| **RecoveryContext** | ② 信息获取 | 恢复上下文评估（"当前情况如何"） |
| **RecoveryPlaybook** | ③ 策略规划 | 恢复策略模板（"应该怎么恢复"） |
| **ExecutionPlan** | ④ 下发执行 | 具体执行计划（"具体怎么调用运行知识的工具来落地"） |
| **AvailabilityCheckpoint** | ⑤ 观察与反馈 | 业务可用性检查点（"怎么判定恢复成功"） |

**ExecutionPlan 的属性**：

| 属性 | 说明 | 示例 |
|------|------|------|
| `steps` | 有序执行步骤列表 | Step1 → Step2 → Step3 |
| `steps[].runningKnowledgeAction` | 映射的运行知识 Action ID | `execute_recovery_task` |
| `steps[].parameterMapping` | 策略参数 → Action 参数的映射规则 | `playbook.granularity → task.recovery_granularity` |
| `steps[].successCriteria` | 该步骤什么算成功 | `job.status == 'completed'` |
| `steps[].errorHandling` | 失败后怎么处理 | `retry_3_times` / `switch_playbook` / `escalate_to_user` |
| `monitoringStrategy` | 执行过程中的监控策略 | `poll_job_status_every_30s` |

---

### Action：5 个

| Action | 绑定 Object | 对应阶段 | 做什么 |
|--------|------------|---------|-------|
| **match_fault_scenario** | FaultScenario | ① | 匹配故障场景 |
| **assess_recovery_context** | RecoveryContext | ② | 系统评估恢复上下文 |
| **generate_recovery_strategy** | RecoveryPlaybook | ③ | 生成可行的恢复策略 |
| **dispatch_execution** | ExecutionPlan | ④ | 将 Playbook + Context 翻译为具体的执行计划，逐步调用运行知识的 Action 来执行，监控执行过程，处理异常 |
| **execute_availability_check** | AvailabilityCheckpoint | ⑤ | 执行业务可用性验证 |

**dispatch_execution 的关键逻辑**：
```
输入: RecoveryPlaybook + RecoveryContext
  ↓
1. 将 Playbook 的高层步骤映射为运行知识的具体 Action 调用
   例: Playbook.step "执行库级异机恢复"
       → 运行知识 Action: execute_recovery_task
       → 参数: { granularity: "database",
                 destination: "remote",
                 clientId: context.targetClient,
                 timePointId: context.selectedTimepoint,
                 recovery_type: "mount" }
  ↓
2. 按顺序调用运行知识 Action
   → execute_recovery_task → 得到 taskId
   → query_recovery_jobs  → 轮询监控进度
  ↓
3. 异常处理
   → 如果作业失败 → 分析错误原因
   → 可重试错误 → 自动重试
   → 不可恢复错误 → 上报用户决策
```

---

### Relation：5 个

| Relation | 源 → 目标 | 对应链路 | 语义 |
|----------|----------|---------|------|
| **scenario_requires_context** | FaultScenario → RecoveryContext | ①→② | 故障场景定义需要收集哪些维度的信息 |
| **fault_resolved_by_playbook** | FaultScenario → RecoveryPlaybook | ①→③ | 故障场景指向候选策略剧本（1:N） |
| **playbook_triggers_risk** | RecoveryPlaybook → Risk | ③→风险 | 策略触发的风险评估 |
| **playbook_produces_execution** | RecoveryPlaybook → ExecutionPlan | ③→④ | 策略剧本关联其对应的执行计划模板（包含映射规则和编排逻辑） |
| **playbook_verified_by_checkpoint** | RecoveryPlaybook → AvailabilityCheckpoint | ③→⑤ | 策略关联有序的验证检查点 |

---

### 完整 5-5-5 架构

```
5个阶段 × (1 Object + 1 Action) + 5 Relation

         经验知识网络（上层 · 认知层）
   ┌───────────────────────────────────────┐
   │                                       │
   │  ① FaultScenario ──scenario_requires_context──→ ② RecoveryContext
   │       │                                              │
   │       │ fault_resolved_by_playbook                   │ (输入)
   │       ↓                                              ↓
   │  ③ RecoveryPlaybook ←── generate_recovery_strategy uses both
   │       │         │
   │       │         │ playbook_triggers_risk──→ Risk (扩展)
   │       │         │
   │       │         │ playbook_produces_execution
   │       │         ↓
   │       │    ④ ExecutionPlan
   │       │         │
   │       │         │ dispatch_execution
   │       │         ↓
   │       │    调用运行知识网络 Action ──────────────────────┐
   │       │                                               │
   │       │ playbook_verified_by_checkpoint                │
   │       ↓                                               │
   │  ⑤ AvailabilityCheckpoint                             │
   │       │                                               │
   │       │ execute_availability_check                    │
   │       ↓                                               │
   └───────┼───────────────────────────────────────────────┘
           │                                               │
           ↓                                               ↓
   ┌───────────────────────────────────────────────────────┐
   │          运行知识网络（下层 · 执行层）                   │
   │  execute_recovery_task / query_recovery_jobs / ...    │
   └───────────────────────────────────────────────────────┘
```

### 汇总

| BKN 元素 | 数量 | 内容 | 对应1-1关系 |
|----------|------|------|-----------|
| **Object** | **5** | FaultScenario, RecoveryContext, RecoveryPlaybook, ExecutionPlan, AvailabilityCheckpoint | 每个阶段产出一个认知实体 |
| **Action** | **5** | match_fault_scenario, assess_recovery_context, generate_recovery_strategy, dispatch_execution, execute_availability_check | 每个阶段执行一个认知操作 |
| **Relation** | **5** | scenario_requires_context, fault_resolved_by_playbook, playbook_triggers_risk, playbook_produces_execution, playbook_verified_by_checkpoint | 串联全认知流程 |

**5-5-5 结构与 5 阶段工作流完全对齐**。每个阶段都在经验知识层有自己的 Object（认知产出）、Action（认知操作），通过 Relation 串联形成一条完整的从"理解问题"到"确认恢复"的认知链路。

方向确认后，从哪个 Object 开始细化属性模型？

### User Input

输出bkn design/experience的object、action、relation的设计文档。要求：1.你参照bkn design/run中object、action、relation的设计文档的布局；2.先写好设计文档，没问题之后再继续编写bkn文件。

### User Input

Continue
