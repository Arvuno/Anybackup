# MySQL 恢复 Agent 的 Experience / Run 联动说明

**版本**：v4.0  
**文档类型**：联动说明文档  
**编写作者**：Andrew  
**目标读者**：Agent 编排工程师  
**最后更新**：2026-04-16

---

## 1. 文档目的

本文档说明恢复 Agent 应如何联合使用 `recovery_experience` 与 `recovery_run`。  
在新的边界下：

- `recovery_experience` 提供恢复认知主数据，本期以 MySQL 为主，并通过 `appType` 区分不同数据库类型
- `recovery_run` 提供运行与作业数据 schema 及动作能力
- 会话上下文负责保存临时推理和临时执行计划

---

## 2. 三层职责划分

### 2.1 Experience 层

Experience 层回答：

- 这更像哪种故障模式
- 还缺哪些关键信息需要先追问和澄清
- 当前备份软件与数据库类型组合理论上具备什么恢复能力
- 哪个恢复方案更合适
- 哪些风险规则需要触发
- 恢复后该做哪些验证

Experience 层不提供 Action。

### 2.2 Run 层

Run 层回答：

- 当前可恢复对象有哪些
- 哪些时间点副本可用
- 哪个客户端和恢复数据源可用
- 当前恢复任务与作业状态如何
- 可以执行哪些查询与恢复动作

### 2.3 Conversation Context

Conversation 上下文负责保存：

- 当前轮次命中的故障模式
- 当前轮次认定的恢复能力边界
- 当前轮次选中的策略模板
- 当前轮次拼装出来的参数
- 当前轮次临时执行计划

这些内容不应回写到经验知识网络。

---

## 3. Agent 标准流程

```text
用户问题
  -> 检索 fault_pattern
  -> 根据 requiredClarification 补齐关键信息
  -> 检索 recovery_capability
  -> 检索 recovery_strategy_template
  -> 若未命中模板，则临时推理候选策略
  -> 检索 risk_rule
  -> 检索 availability_checkpoint_template
  -> 若策略为 inferred，则必须人工确认
  -> query_protection_objects / query_timepoint_copies / query_clients / query_recovery_datasources
  -> 会话上下文中临时生成执行方案
  -> execute_recovery_task / query_recovery_jobs
  -> execute_availability_verification
```

---

## 4. Experience 层的使用方式

### 4.1 故障识别

Agent 不需要在经验层中创建 `fault_scenario` 对象，也不需要调用经验层 Action。  
应改为：

1. 从用户输入中抽取症状关键词和恢复意图关键词
2. 直接检索 `fault_pattern`
3. 在返回的主数据记录中选择最匹配项
4. 根据 `requiredClarification` 继续追问，直到发生时间、影响范围、实例还是库等关键信息明确

### 4.2 恢复能力认知

Agent 不需要把能力判断交给运行层。  
应改为：

1. 仅在 `requiredClarification` 已补齐后，基于已识别故障和 `appType` 检索 `recovery_capability`
2. 明确当前数据库类型理论上支持哪些恢复粒度和恢复技术
3. 将能力边界带入后续策略决策

### 4.3 恢复方案选择

Agent 不需要在经验层中创建 `recovery_playbook` 对象。  
应改为：

1. 基于 `faultPatternId`、`appType` 和本期策略模板中的 `vendor` 检索 `recovery_strategy_template`
2. 结合 `recovery_capability` 过滤不可用方案
3. 若命中正式模板，则在上下文中临时形成本次推荐策略
4. 若未命中正式模板，则允许基于 `fault_pattern` 与 `recovery_capability` 临时推理候选策略，但必须标记为 `inferred` / `unverified`

### 4.4 风险与审批

Agent 不需要在经验层中创建 `approval` 或 `risk assessment` 结果对象。  
应改为：

1. 检索 `risk_rule`
2. 判断是否命中高风险或必须审批条件
3. 若策略来源于正式模板，按正常风险规则判断
4. 若策略来源于临时推理，则默认视为高风险未验证策略
5. 命中高风险、必须审批或未验证策略时，暂停自动执行并转入人审

### 4.5 验证模板获取

Agent 不需要在经验层中创建 `availability_checkpoint` 对象。  
应改为：

1. 检索 `availability_checkpoint_template`
2. 在上下文中组织恢复后验证步骤
3. 再调用运行层验证动作执行

---

## 5. Run 层的使用方式

当经验层已经给出故障模式、能力边界和策略模板后，Agent 应转向运行层查询实时数据：

- `query_protection_objects`
- `query_timepoint_copies`
- `query_recovery_datasources`
- `query_clients`
- `query_recovery_jobs`

然后在上下文中组合：

- 恢复目标对象
- 恢复时间点
- 恢复数据源
- 执行客户端
- 参数绑定

---

## 6. 执行计划的正确位置

`execution_plan` 不再是经验层 BKN 对象。  
它应只存在于 Agent 当前任务的上下文中，作为临时结构被组织和消费。

这样做的原因是：

- 它强依赖当前会话输入
- 它强依赖当前运行层实时数据
- 它不具备长期复用价值

---

## 7. 必须人审的场景

Agent 在以下情况下必须停止自动执行并请求人工确认：

- 命中 `risk_rule.approvalRequired = true`
- 风险等级为 `high`
- 未命中正式 `recovery_strategy_template`，只能依赖临时推理策略
- 目标是原机恢复且会覆盖生产数据
- 运行层返回的候选资源不完整或歧义较大

---

## 8. 联动结论

新的联动原则可以概括为：

- Experience 层提供“这是什么故障、能怎么恢复、应如何验证”
- Run 层提供“当前能怎么做”
- Conversation Context 保存“这次准备怎么做”

这样 Agent 才能在不污染 BKN 本体的前提下完成恢复决策闭环。
