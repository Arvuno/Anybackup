# MySQL 数据库恢复知识网络 - 经验知识 Object 数据模型设计

**版本**：v5.0  
**文档类型**：数据模型设计文档  
**编写作者**：Andrew  
**目标读者**：BKN 建模工程师  
**最后更新**：2026-04-16

---

## 1. 设计目标

本文档定义 `recovery_experience` 的 Object 模型。  
在新的边界下，经验知识网络不是 Agent 会话态推理结果的存储层，而是数据库恢复领域的经验主数据 schema 层。真实数据存放在关系型数据库中，通过 dataview 做映射检索。  
虽然当前文档沿用 MySQL 命名，但本期实现范围以 MySQL 为主；experience 的建模方式仍保留 database 域扩展能力，并通过 `appType` 区分不同数据库类型。

其核心设计思想是：

- `recovery_experience` 描述“可复用恢复认知”
- `recovery_run` 描述“Foundation 执行世界”
- Agent 的临时推理结果只保留在会话上下文中，不进入 BKN

---

## 2. 设计原则

### 2.1 只建模可复用主数据

经验层只建模长期可复用的经验知识，不建模一次会话中的临时对象。

### 2.2 以恢复认知链组织对象

经验层核心对象对应 5 个问题：

1. 这是什么故障
2. Foundation 理论上具备什么恢复能力
3. 可以采用什么恢复方案
4. 该方案能否自动执行
5. 如何验证恢复不是假恢复

### 2.3 运行数据不在经验层重复建模

经验层不重复存储以下运行数据：

- `protection_object`
- `timepoint_copy`
- `recovery_datasource`
- `client`
- `recovery_task`
- `recovery_job`
- `availability_verification`

这些对象继续由 `recovery_run` 提供。

---

## 3. Object 总览

| Object Type | 语义角色 | 对应问题 |
|-------------|---------|---------|
| `fault_pattern` | 故障模式主数据 | 这是什么故障 |
| `recovery_capability` | 恢复能力主数据 | Foundation 理论上具备什么恢复能力 |
| `recovery_strategy_template` | 恢复方案主数据 | 可以采用什么恢复方案 |
| `risk_rule` | 风险控制主数据 | 该方案能否自动执行 |
| `availability_checkpoint_template` | 验证模板主数据 | 如何验证不是假恢复 |

---

## 4. FaultPattern

### 4.1 语义定义

`fault_pattern` 用于定义可复用的故障模式知识。  
像“业务宕机”“数据丢失”“勒索病毒”“误删数据”这类内容，应作为该对象的数据记录存放在关系型数据库中，而不是只在 Prompt 或会话中临时存在。

### 4.2 核心字段

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `patternId` | string | 是 | 故障模式唯一标识 |
| `name` | string | 是 | 故障模式名称 |
| `appType` | string | 是 | 应用类型，如 MySQL、Oracle、PostgreSQL、SQLServer |
| `faultCategory` | string | 是 | 一级故障分类 |
| `affectedGranularity` | string | 是 | 常见影响粒度 |
| `symptomKeywords` | text | 否 | 症状关键词集合 |
| `intentKeywords` | text | 否 | 用户表达恢复意图时的关键词集合 |
| `requiredClarification` | text | 否 | 必须先澄清的信息项，如发生时间、影响范围、提供业务的是实例还是库 |
| `disposalHint` | text | 否 | 处理提示 |
| `severityBaseline` | string | 否 | 默认严重度 |
| `enabled` | boolean | 是 | 是否启用 |

### 4.3 使用约束

- Agent 在进入恢复能力认知和恢复方案选择之前，必须先补齐 `requiredClarification`
- 本期不新增独立 `clarification_rule` 对象，必要追问信息先由 `fault_pattern.requiredClarification` 承载

---

## 5. RecoveryCapability

### 5.1 语义定义

`recovery_capability` 用于定义特定备份软件在特定数据库类型上的恢复能力边界。  
它回答的不是“当前有哪些实时资源”，而是“在当前 vendor 与数据库类型组合下，理论上能如何恢复”。

### 5.2 核心字段

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `capabilityId` | string | 是 | 恢复能力唯一标识 |
| `name` | string | 是 | 恢复能力名称 |
| `vendor` | string | 是 | 备份软件厂商或平台，如 AISHU、Veeam、NBU、Commvault |
| `appType` | string | 是 | 应用类型，如 MySQL、Oracle、PostgreSQL、SQLServer |
| `supportedGranularity` | text | 是 | 支持的恢复粒度集合，如 instance、database、log_file、table_indirect |
| `supportedTechnique` | text | 是 | 支持的恢复技术类型集合，如 restore、mount |
| `supportedMode` | text | 是 | 支持的恢复模式集合，如 shortest_time、latest_state、specified_time |
| `supportsOriginalRestore` | boolean | 是 | 是否支持原机恢复 |
| `supportsRemoteRestore` | boolean | 是 | 是否支持异机恢复 |
| `supportsPointInTimeRestore` | boolean | 是 | 是否支持定点恢复 |
| `supportsLogRestore` | boolean | 是 | 是否支持日志恢复 |
| `tableRecoveryMode` | string | 否 | 表级恢复模式 |
| `capabilitySummary` | text | 否 | vendor 与数据库类型组合下的能力边界摘要 |
| `enabled` | boolean | 是 | 是否启用 |

---

## 6. RecoveryStrategyTemplate

### 6.1 语义定义

`recovery_strategy_template` 用于定义可复用的恢复方案模板。  
它表达在已识别故障和已知能力边界下可采用的恢复方案。

如果未命中正式 `recovery_strategy_template`，Agent 可以基于 `fault_pattern` 与 `recovery_capability` 临时推理候选方案。  
但这类策略不属于经验层已验证模板，默认视为“未验证策略”，不得自动执行。

### 6.2 核心字段

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `strategyTemplateId` | string | 是 | 策略模板唯一标识 |
| `name` | string | 是 | 模板名称 |
| `vendor` | string | 是 | 当前策略模板适用的备份软件平台。本期默认按 MySQL + AISHU 设计，并为未来扩展预留字段 |
| `appType` | string | 是 | 应用类型，如 MySQL、Oracle、PostgreSQL、SQLServer |
| `faultPatternId` | string | 是 | 关联故障模式 |
| `recoveryGranularity` | string | 是 | 推荐恢复粒度 |
| `destinationType` | string | 是 | 推荐恢复目标 |
| `recoveryMethod` | string | 是 | 推荐恢复方式 |
| `requiresRecovery` | boolean | 是 | 是否真的需要恢复 |
| `strategySummary` | text | 是 | 方案摘要 |
| `riskBaseline` | string | 否 | 默认风险等级 |
| `approvalRequired` | boolean | 是 | 是否默认需要审批 |
| `enabled` | boolean | 是 | 是否启用 |

### 6.3 约束说明

- `vendor` 本期只进入 `recovery_strategy_template`
- 本期不将 `vendor` 扩散到 `risk_rule`、`availability_checkpoint_template` 或 `fault_pattern`

---

## 7. RiskRule

### 7.1 语义定义

`risk_rule` 用于定义恢复方案的风险约束。  
它负责告诉 Agent：哪些方案需要审批、哪些场景风险更高、需要补充哪些缓解措施。

此外，`risk_rule` 还需要覆盖一种通用高风险场景：  
当 Agent 没有命中正式 `recovery_strategy_template`，只能基于 `fault_pattern` 与 `recovery_capability` 临时推理策略时，应将该策略视为“未验证策略”，默认 `high risk` 且必须人工确认。

### 7.2 核心字段

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `riskRuleId` | string | 是 | 风险规则唯一标识 |
| `name` | string | 是 | 风险规则名称 |
| `appType` | string | 是 | 应用类型，如 MySQL、Oracle、PostgreSQL、SQLServer |
| `strategyTemplateId` | string | 是 | 关联策略模板 |
| `triggerCondition` | text | 是 | 风险触发条件 |
| `riskLevel` | string | 是 | 风险等级 |
| `approvalRequired` | boolean | 是 | 是否必须审批 |
| `mitigationAdvice` | text | 否 | 缓解建议 |
| `enabled` | boolean | 是 | 是否启用 |

---

## 8. AvailabilityCheckpointTemplate

### 8.1 语义定义

`availability_checkpoint_template` 用于定义恢复完成后的验证模板。  
当前阶段它主要承载与 `run` 层实际工具能力一致的通用基础验证模板，而不是按故障场景细分的业务验证模板体系。

### 8.2 核心字段

| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `checkpointTemplateId` | string | 是 | 验证模板唯一标识 |
| `name` | string | 是 | 模板名称 |
| `appType` | string | 是 | 应用类型，如 MySQL、Oracle、PostgreSQL、SQLServer |
| `checkpointType` | string | 是 | 验证类型 |
| `targetScope` | string | 是 | 验证目标范围 |
| `verificationMethod` | text | 是 | 验证执行方式 |
| `successCriteria` | text | 是 | 验证通过标准 |
| `faultPatternId` | string | 否 | 可选，直接关联某类故障模式；当前阶段初始化数据可为空 |
| `strategyTemplateId` | string | 否 | 关联策略模板；当前阶段初始化数据可为空 |
| `enabled` | boolean | 是 | 是否启用 |

### 8.3 当前阶段约束

- 当前 MySQL 初始化数据仅保留一条通用基础 SQL 检索验证模板
- 当前初始化数据不按 `strategyTemplateId` 或 `faultPatternId` 建立 direct 映射
- `strategyTemplateId` 与 `faultPatternId` 保留为后续细化验证模板的扩展字段
- 后续随着 `run` 层验证工具能力增强，再逐步补齐场景化验证模板

---

## 9. Object 设计结论

M1 经验层 Object 设计的重点，是建立一条清晰的恢复认知链：

```text
FaultPattern
  -> RecoveryCapability
  -> RecoveryStrategyTemplate
  -> RiskRule
```

当前阶段 `AvailabilityCheckpointTemplate` 以单条通用基础验证模板存在；后续随着工具能力增强，再逐步按 `RecoveryStrategyTemplate` 或故障场景细化。
