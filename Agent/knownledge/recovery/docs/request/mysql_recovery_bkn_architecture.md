# MySQL 数据库恢复知识网络架构设计

**版本**：v5.0  
**文档类型**：架构设计文档  
**编写作者**：Andrew  
**目标读者**：产品经理 / BKN 建模工程师  
**最后更新**：2026-04-16

---

## 1. 设计目标

本文档定义 MySQL 恢复知识网络的总体架构。  
该架构采用上下两层知识网络：

- `recovery_experience`：经验知识网络，负责定义可复用的恢复认知主数据 schema
- `recovery_run`：运行知识网络，负责定义 Foundation 执行世界的 schema 与执行能力

两层网络都属于 BKN schema，真实数据存放在关系型数据库中，并通过 dataview 做映射检索。  
其中 `recovery_experience/database` 的建模方式面向 database 域统一扩展，但本期实现范围以 MySQL 为主，通过 `appType` 区分不同数据库类型。

---

## 2. 核心边界

### 2.1 `recovery_experience` 的职责

经验层只负责表达可复用的恢复认知知识，例如：

- 故障模式
- 恢复能力
- 恢复策略模板
- 风险规则
- 当前阶段的通用基础可恢复性验证模板

经验层不定义运行态 `action_type`，也不承担任何操作职责。  
其中 `recovery_capability` 需要同时包含 `vendor` 与 `appType`，用于区分不同备份软件在不同数据库类型上的能力边界。当前默认示例 vendor 为 `AISHU`，后续可扩展到 Veeam、NBU、Commvault 等。
本期 `vendor` 只进入 `recovery_strategy_template` 作为策略模板属性，不扩散到 `risk_rule` 或 `availability_checkpoint_template`。

### 2.2 `recovery_run` 的职责

运行层负责表达 Foundation 执行世界中的实体、作业数据与动作能力，例如：

- 保护对象
- 时间点副本
- 恢复数据源
- 客户端
- 恢复任务
- 恢复作业
- 可用性验证结果
- 查询与执行 Action

### 2.3 不进入 BKN 的内容

以下内容不进入 BKN：

- Agent 一次会话中的临时推理结果
- 临时执行计划
- 自然语言解释
- reasoning 过程

这些内容只保留在 conversation 或任务上下文中。

---

## 3. 总体工作方式

恢复 Agent 的工作方式应为：

1. 先从 `recovery_experience` 检索故障模式
2. 根据 `fault_pattern.requiredClarification` 先补齐必要澄清信息
3. 再检索恢复能力
4. 再检索恢复策略模板
5. 若未命中正式策略模板，则在会话上下文中基于故障与能力临时推理候选方案
6. 再检索风险规则与通用基础验证模板
7. 对临时推理方案按“未验证策略”处理，要求人工确认后才能继续
8. 在会话上下文中临时组织本次恢复方案与参数
9. 调用 `recovery_run` 的 Action 执行恢复
10. 基于经验层通用基础验证模板和运行层验证动作完成恢复验证

---

## 4. 经验层本体设计

### 4.1 核心 ObjectTypes

| Object Type | 作用 |
|-------------|------|
| `fault_pattern` | 可复用故障模式主数据，并声明必须先澄清的信息 |
| `recovery_capability` | 可复用恢复能力主数据 |
| `recovery_strategy_template` | 可复用恢复方案模板，本期包含 `vendor` 属性 |
| `risk_rule` | 风险与审批控制规则 |
| `availability_checkpoint_template` | 当前阶段的通用基础验证模板 |

### 4.2 核心 RelationTypes

| Relation Type | 作用 |
|---------------|------|
| `fault_pattern_supported_by_recovery_capability` | 故障模式受恢复能力支撑 |
| `fault_pattern_resolved_by_strategy_template` | 故障模式默认适配恢复方案 |
| `recovery_capability_constrains_strategy_template` | 恢复能力约束可选方案 |
| `strategy_template_constrained_by_risk_rule` | 方案受风险规则约束 |
| `strategy_template_verified_by_checkpoint_template` | 方案对应验证模板 |

### 4.3 检索方式

经验层通过 dataview 直接暴露主数据记录，不再额外定义 `action_type`。  
Agent 需要什么知识，就直接按对象和关系去检索什么知识。  
在进入恢复能力认知和正式策略选择之前，Agent 必须先根据 `requiredClarification` 补齐关键信息。

补充约束：

- 当前 `fault_pattern -> recovery_strategy_template` 更适合理解为“默认主模板映射”
- 当前 `availability_checkpoint_template` 以通用基础验证模板为主；`strategy_template -> availability_checkpoint_template` relation 继续保留，供后续按策略细化验证模板时使用；不再通过 `fault_pattern` 直接绑定

---

## 5. 运行层本体设计

运行层保持既有设计，不作为本次架构重写重点。  
其核心特点是：

- 直接映射 Foundation 执行世界
- Object/Relation/Action 面向运行和作业数据
- 为 Agent 提供真实可执行的恢复操作能力

---

## 6. 上下层联动关系

### 6.1 经验层提供恢复认知

经验层告诉 Agent：

- 这是什么故障
- 这次故障还缺哪些关键信息需要先澄清
- 当前 vendor 与数据库类型组合理论上具备什么恢复能力
- 可以采用什么恢复方案
- 该方案能否自动执行
- 恢复后应做哪些验证

如果未命中正式 `recovery_strategy_template`，Agent 可以在会话上下文中临时推理候选策略；但这类策略不属于经验层已验证知识，必须视为高风险并要求人工确认，不能自动执行。

### 6.2 运行层提供实时执行能力

运行层告诉 Agent：

- 当前有哪些保护对象
- 有哪些时间点副本可用
- 可以选择哪些客户端和恢复数据源
- 当前恢复任务与恢复作业状态如何
- 可以执行哪些查询和恢复动作

### 6.3 会话上下文负责临时组织方案

经验层和运行层之间不通过经验层 BKN 对象做正式桥接。  
桥接逻辑发生在 Agent 的会话上下文中：

- Agent 从经验层取回主数据
- Agent 从运行层取回实时数据
- Agent 在上下文中临时形成本次恢复计划
- Agent 再调用运行层 Action 执行

---

## 7. 关键架构结论

本次架构调整后，MySQL 恢复知识网络形成了清晰的三段式边界：

- `recovery_experience`：存恢复认知主数据
- `recovery_run`：存运行与作业数据 schema 及动作
- conversation context：存本次推理和临时计划

这种划分更符合 Ontology + dataview 的设计思想，也更适合后续将经验知识持续沉淀到关系型数据库中。
