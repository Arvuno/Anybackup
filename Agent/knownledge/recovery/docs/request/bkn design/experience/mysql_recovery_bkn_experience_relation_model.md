# MySQL 数据库恢复知识网络 - 经验知识 Relation 数据模型设计

**版本**：v5.0  
**文档类型**：关系模型设计文档  
**编写作者**：Andrew  
**目标读者**：BKN 建模工程师  
**最后更新**：2026-04-16

---

## 1. 设计目标

本文档定义 `recovery_experience` 的 Relation 模型。  
在新的边界下，经验层 Relation 不再表达“故障 -> 规则 -> 策略”的旧链路，而是表达恢复认知主数据之间的稳定语义连接。

---

## 2. 设计原则

### 2.1 关系用于连接恢复认知主数据

经验层 Relation 的目标是表达当前 data 已经稳定落库的恢复认知主数据连接：

- 某类故障由哪种恢复能力支撑
- 某类故障默认适用哪些恢复方案
- 某类恢复能力约束了哪些方案
- 某种方案有哪些风险规则

### 2.2 关系应能映射到数据库字段

每种 RelationType 都应能落到关系型数据库中的外键字段或稳定属性映射，不依赖会话态内存对象。

### 2.3 全部使用 direct 关系

当前 experience 层所有关系都采用 `direct`。  
因为它们本质上都是稳定的主数据映射关系，而不是运行态推导关系。

---

## 3. Relation 总览

| Relation Type | Source | Target | 语义 |
|---------------|--------|--------|------|
| `fault_pattern_supported_by_recovery_capability` | `fault_pattern` | `recovery_capability` | 故障模式受恢复能力支撑 |
| `fault_pattern_resolved_by_strategy_template` | `fault_pattern` | `recovery_strategy_template` | 故障模式默认由哪些策略模板解决 |
| `recovery_capability_constrains_strategy_template` | `recovery_capability` | `recovery_strategy_template` | 恢复能力约束可选策略模板 |
| `strategy_template_constrained_by_risk_rule` | `recovery_strategy_template` | `risk_rule` | 策略模板受风险规则约束 |
| `strategy_template_verified_by_checkpoint_template` | `recovery_strategy_template` | `availability_checkpoint_template` | 策略模板对应验证模板 |

---

## 4. FaultPatternSupportedByRecoveryCapability

### 4.1 语义定义

`fault_pattern_supported_by_recovery_capability` 表示识别故障后，需要先理解该类数据库在 Foundation 下具备哪些恢复能力。

补充约束：

- Agent 在沿该关系检索恢复能力之前，必须先根据 `fault_pattern.requiredClarification` 补齐必要澄清信息
- 本期 relation 仍按 `appType` 建立，不把 `vendor` 纳入映射规则

### 4.2 映射建议

- Source：`fault_pattern`
- Target：`recovery_capability`
- 映射字段：`appType -> appType`

---

## 5. FaultPatternResolvedByStrategyTemplate

### 5.1 语义定义

`fault_pattern_resolved_by_strategy_template` 表示某类故障默认适配哪些恢复方案。

补充约束：

- 只有当 `requiredClarification` 中的关键信息已经明确后，Agent 才能沿该关系进入正式策略选择
- 本期 `vendor` 仅作为 `recovery_strategy_template` 的对象属性，不参与 relation 映射
- 当前该关系仅表达“默认主模板映射”，不是穷尽映射；当多个故障模式复用同一策略模板时，可能只有其中一个故障模式通过 `faultPatternId` 建立直接关联
- 如果沿该关系未命中任何正式模板，Agent 可以基于 `fault_pattern` 与 `recovery_capability` 临时推理候选策略，但该策略不视为 relation 命中结果，必须按“未验证策略”走高风险确认流程

### 5.2 映射建议

- Source：`fault_pattern`
- Target：`recovery_strategy_template`
- 映射字段：`patternId -> faultPatternId`

---

## 6. RecoveryCapabilityConstrainsStrategyTemplate

### 6.1 语义定义

`recovery_capability_constrains_strategy_template` 表示某类数据库的能力边界约束了哪些策略模板可被采用。

补充约束：

- 本期仍按 `appType` 描述能力到策略的约束关系
- 多 vendor 约束问题留到后续阶段处理，本期不把 `vendor` 升级为 relation 级约束

### 6.2 映射建议

- Source：`recovery_capability`
- Target：`recovery_strategy_template`
- 映射字段：`appType -> appType`

---

## 7. StrategyTemplateConstrainedByRiskRule

### 7.1 语义定义

`strategy_template_constrained_by_risk_rule` 表示某个恢复方案在执行前需要遵循哪些风险控制和审批要求。

### 7.2 映射建议

- Source：`recovery_strategy_template`
- Target：`risk_rule`
- 映射字段：`strategyTemplateId -> strategyTemplateId`

---

## 8. StrategyTemplateVerifiedByCheckpointTemplate

### 8.1 语义定义

`strategy_template_verified_by_checkpoint_template` 表示某个恢复方案在执行完成后，应采用哪些验证模板来避免假恢复。

补充约束：

- 当前阶段虽然验证能力仍较基础，且初始化数据只保留一条通用 MySQL 基础 SQL 检索验证模板，但 relation 继续保留，用于保证 ontology 具备后续按恢复策略细化验证模板的扩展能力
- 后续随着 `run` 层验证能力增强，可继续按 `strategyTemplateId` 扩展更细的验证模板

### 8.2 映射建议

- Source：`recovery_strategy_template`
- Target：`availability_checkpoint_template`
- 映射字段：`strategyTemplateId -> strategyTemplateId`

---

## 9. 当前阶段未启用的验证关系

当前不再把验证模板按故障模式直接挂接，因此以下关系在当前 data 阶段不作为稳定 direct relation 使用：

- `fault_pattern_verified_by_checkpoint_template`

后续若重新沉淀出按故障场景细分的专项验证模板，再恢复该关系。

---

## 10. Relation 设计结论

M1 经验层 Relation 的目标，是把恢复认知主数据组织成一条清晰的业务主链：

```text
fault_pattern
  -> recovery_capability
  -> recovery_strategy_template
  -> risk_rule
  -> availability_checkpoint_template
```

当前阶段的 `availability_checkpoint_template` 以通用基础验证模板为主；`strategy_template -> availability_checkpoint_template` relation 继续保留，供后续按策略细化验证模板时使用；不再通过 `fault_pattern` 直接绑定。
