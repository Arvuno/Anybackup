# MySQL 数据库恢复知识网络实现设计

**版本**：v5.0  
**文档类型**：实现设计文档  
**编写作者**：Andrew  
**目标读者**：BKN 建模工程师  
**最后更新**：2026-04-16

---

## 1. 实现目标

本文档定义 `recovery_experience` 的落地实现方案。  
本次实现目标是构建一套可落库、可通过 dataview 检索的恢复认知主数据 schema，而不是构建会话态对象或经验层动作。  
experience 的建模方式面向 database 域统一扩展，但本期实现范围以 MySQL 为主，依靠 `appType` 区分不同数据库类型，并在 `recovery_capability` 中使用 `vendor` 区分不同备份软件平台。

---

## 2. 落地范围

### 2.1 In Scope

- 在 `examples/bkn/recovery_experience/database/` 下落经验知识网络 `.bkn`
- 定义经验层的 `object_types`、`relation_types`
- 为每个 ObjectType 补充 `Data Source`
- 明确经验层由关系型数据库 + dataview 承载

### 2.2 Out of Scope

- 不在经验层中建模临时会话对象
- 不在经验层中建模 Foundation 作业数据
- 不在经验层中存储临时执行计划
- 不在经验层中定义 `action_types`
- 不在本期实现经验数据库建表 SQL

---

## 3. 目录结构

经验层采用标准 BKN 模块化目录：

```text
examples/bkn/recovery_experience/database/
├── network.bkn
├── object_types/
├── relation_types/
├── data/
└── SKILL.md
```

---

## 4. 对象实现

### 4.1 ObjectTypes

| Object Type | 说明 | 建议 Data View |
|-------------|------|----------------|
| `fault_pattern` | 故障模式主数据 | `experience_fault_pattern` |
| `recovery_capability` | 恢复能力主数据 | `experience_recovery_capability` |
| `recovery_strategy_template` | 恢复方案模板主数据 | `experience_strategy_template` |
| `risk_rule` | 风险规则主数据 | `experience_risk_rule` |
| `availability_checkpoint_template` | 当前阶段的通用基础验证模板主数据 | `experience_availability_checkpoint_template` |

### 4.2 数据设计要求

- 每个 ObjectType 必须有稳定主键
- 每个 ObjectType 必须包含 `appType`
- 每个 ObjectType 必须能映射到一个 dataview
- `Mapped Field` 应优先对应数据库字段名
- 不使用会话态字段如 `sourceInput`、`planSummary`、`observedResult`

补充约束：

- `recovery_capability` 必须包含 `vendor`
- `recovery_strategy_template` 必须包含 `vendor`
- `fault_pattern` 必须包含 `requiredClarification`
- 当前原始数据默认 `vendor = AISHU`
- 后续可扩展到 `Veeam`、`NBU`、`Commvault`
- 本期不新增独立 `clarification_rule` 对象
- 本期不把 `vendor` 扩散到 `risk_rule` 或 `availability_checkpoint_template`

---

## 5. 关系实现

### 5.1 RelationTypes

| Relation Type | 实现方式 |
|---------------|---------|
| `fault_pattern_supported_by_recovery_capability` | `appType -> appType` |
| `fault_pattern_resolved_by_strategy_template` | `patternId -> faultPatternId` |
| `recovery_capability_constrains_strategy_template` | `appType -> appType` |
| `strategy_template_constrained_by_risk_rule` | `strategyTemplateId -> strategyTemplateId` |
| `strategy_template_verified_by_checkpoint_template` | `strategyTemplateId -> strategyTemplateId` |

### 5.2 设计要求

- 关系优先落在外键字段或稳定属性映射上
- 不通过自然语言描述隐式推理关系
- 不复制运行层实体关系
- Agent 在沿 `fault_pattern_supported_by_recovery_capability` 和 `fault_pattern_resolved_by_strategy_template` 检索前，必须先补齐 `requiredClarification`
- 本期 `vendor` 仅作为 `recovery_strategy_template` 的对象属性，不进入 relation 映射规则
- 当前 `fault_pattern_resolved_by_strategy_template` 只表达默认主模板映射，不保证穷尽所有可复用模板
- 当前 `availability_checkpoint_template` 初始化数据采用单条通用基础验证模板；`strategyTemplateId` relation 继续保留，供后续按策略细化验证模板时使用；不再通过 `fault_pattern` 建立稳定 relation

---

## 6. 检索实现

经验层不定义 `action_type`。  
Agent 对经验层的使用方式应为：

1. 直接通过 dataview 检索 `fault_pattern`
2. 根据 `requiredClarification` 补齐必要澄清信息
3. 再检索 `recovery_capability`
4. 再检索 `recovery_strategy_template`
5. 若未命中正式策略模板，则基于 `fault_pattern` 与 `recovery_capability` 在会话上下文中临时推理候选策略
6. 再检索 `risk_rule`
7. 再检索通用 `availability_checkpoint_template`
8. 若策略来源于临时推理，则按“未验证策略”处理并要求人工确认
9. 在会话上下文中形成临时恢复方案
10. 进入 `recovery_run` 查询和执行

---

## 7. 与运行层的集成方式

Agent 的实际恢复流程应为：

1. 从经验层主数据中识别故障模式
2. 根据 `requiredClarification` 先补齐故障发生时间、影响范围、实例或库归属等关键信息
3. 从经验层主数据中认知恢复能力
4. 优先从经验层主数据中匹配正式恢复方案；未命中时允许临时推理候选策略
5. 从经验层主数据中判断风险与验证要求
6. 若策略来自临时推理，则默认按高风险处理并要求人工确认
7. 在 `recovery_run` 中查询保护对象、时间点副本、客户端等实时数据
8. 在会话上下文中形成临时恢复方案
9. 调用运行层 Action 执行恢复
10. 读取通用基础验证模板并驱动运行层验证动作

因此：

- 经验层不定义 `execution_plan`
- `execution_plan` 只作为 Agent 上下文中的临时结构
- 运行层 Action 是唯一执行入口

---

## 8. CSV 原始数据要求

`data/` 目录下的 CSV 作为后续入库原始数据，需与 `.bkn` schema 保持一致：

- `fault_pattern.csv`
- `recovery_capability.csv`
- `recovery_strategy_template.csv`
- `risk_rule.csv`
- `availability_checkpoint_template.csv`

其中：

- `recovery_strategy_template.csv` 必须包含 `faultPatternId`
- `recovery_capability.csv` 必须与 `recovery_capability.bkn` 字段完全一致，并包含 `vendor`
- 若后续落地“未验证策略风险”主数据，则应在 `risk_rule.csv` 中补充对应规则记录

---

## 9. 实施顺序

1. 重写经验层设计文档
2. 重写 `network.bkn`
3. 重写 `object_types/*.bkn`
4. 重写 `relation_types/*.bkn`
5. 删除旧对象、旧关系与旧概念组
6. 校对 `data/*.csv`
7. 运行 `kweaver bkn validate examples/bkn/recovery_experience/database`

---

## 10. 实现结论

本次实现将 `recovery_experience` 收敛为恢复认知主数据 schema，避免把 BKN 误用为会话态推理存储层或操作层。  
这样可以保证：

- 经验知识能沉淀到数据库
- dataview 检索边界清晰
- 与 `recovery_run` 的职责不冲突
- Agent 的推理过程保持灵活，不被静态 schema 过度约束
