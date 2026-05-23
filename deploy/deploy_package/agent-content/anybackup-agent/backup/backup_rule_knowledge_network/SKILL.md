# 备份规则知识网络 - Agent 使用指南

> **网络ID**: backup_rule_knowledge_network
> **版本**:
> **标签**: 备份方案执行, 备份方案推荐

## 网络概览

面向备份方案推荐的新知识网络。该网络以 `backup_solution` 作为推荐结果的中心对象，并将多个 `backup_window`、多个 `backup_plan` 和一个 `retention_plan` 统一纳入方案内部。`basic_config_rule` 与 `app_config_rule` 采用“配置规则项”建模，不再作为单个模板对象直连到方案，而是在推荐阶段基于应用类型、数据量、RPO、窗口和数据级别等上下文匹配并汇总多条规则；同一候选方案在最终输出时按 `config_category` 去重，保证同一配置类别不会重复出现。`backup_plan`、`backup_window`、`retention_plan` 的字段表达与规格定义中的参数项保持一致。`recommendation_policy` 负责为候选备份方案提供行业级默认值、法规缩圈与评分权重，并通过显式关系表达其与候选方案的适配范围。`app_output_rule` 用于定义不同应用类型的最终输出方案结构，并通过显式关系表达不同应用方案的输出归属，避免在 action 中硬编码 mysql、oracle、file 等应用的展示层级。应用类型、数据类型、数据量、RPO、RTO、备份窗口和行业类型仍作为外部输入参数；规则对象统一作为约束来源，通过显式关系收敛候选备份方案。

### 核心对象

| 对象 | 文件路径 | 说明 |
|------|----------|------|
| 应用备份能力 | `object_types/app_backup_capability.bkn` | 用于描述某类应用支持的备份能力事实，确保推荐阶段不会输出应用本身不支持的备份计划条目或应用配置。 |
| 应用配置 | `object_types/app_config_rule.bkn` | 用于表达应用配置规则项。每条记录只表达一个应用配置类别，例如 MySQL 的备份方式、即时合成备份、实时日志备份、归档日志删除策略、通道数、自动转完备、页面追踪、自动清理追踪数据、自定义参数等。推荐阶段需按 `config_category` 汇总，并保证同一候选方案的最终应用配置结果中，同一配置类别只出现一次。 |
| 应用输出规则 | `object_types/app_output_rule.bkn` | 用于定义不同应用类型的推荐结果展示层级与节点命名规则，避免在 action 中为 mysql、oracle、file 等应用硬编码输出格式。每条记录描述一种应用输出方案，推荐阶段先匹配输出规则，再按规则组织最终返回结构。 |
| 应用级规则 | `object_types/app_rule.bkn` | 用于根据应用类型、数据类型与法规要求筛选可用备份方案的规则对象。该对象不再内嵌候选方案 ID，而是通过关系类型显式关联到 `backup_solution`。 |
| 备份计划 | `object_types/backup_plan.bkn` | 用于表达备份方案中的计划条目。该对象按计划类型建模，支持分钟级计划、小时级计划、天级计划、周级计划、月级计划、季度级计划、年级计划；每条记录只表达一个计划条目，并使用对应计划类型实际需要的参数。 |
| 备份方案 | `object_types/backup_solution.bkn` | 作为推荐结果的中心对象，用于统一承载多个备份窗口、多个备份计划和一个数据保留计划。基础配置与应用配置不再以单个模板对象直连，而是在推荐阶段基于方案上下文匹配并汇总多条配置规则项。 |
| 备份窗口 | `object_types/backup_window.bkn` | 用于表达备份方案中的窗口条目。一个备份方案下允许存在多个窗口条目，每条记录只表达一个窗口片段，并使用统一参数定义：窗口类型、开始星期、结束星期、开始时间、结束时间、最大速度。 |
| 基础配置 | `object_types/basic_config_rule.bkn` | 用于表达基础配置规则项。每条记录只表达一个基础配置类别，例如备份存储、数据压缩、备份失败重试、强制数据保留、重复数据删除、传输加密等。推荐阶段需按 `config_category` 汇总，并保证同一候选方案的最终基础配置结果中，同一配置类别只出现一次。 |
| 数据规则 | `object_types/data_rule.bkn` | 描述业务数据类型、数据分级属性和典型内容之间的映射关系，用于将外部输入的数据类型归入合适的保护层级，并作为数据侧规则对象关联到备份方案。 |
| 法规要求 | `object_types/legal_rule.bkn` | 描述影响备份方案选择的法规约束对象。法规规则通过关系类型显式关联到 `backup_solution`，用于提供合规约束和评分依据。 |
| 候选推荐策略 | `object_types/recommendation_policy.bkn` | 描述候选方案推荐时使用的默认值、综合评估权重以及最多输出数量。它不承载查询输入本身，而是为 action 提供稳定的评估与回退规则。 |
| 数据保留计划 | `object_types/retention_plan.bkn` | 定义备份方案中的唯一副本保留计划对象。该对象使用与规格定义一致的三类参数分组：数据保留设置、特定周期保留、日志保留设置。 |
| RPO等级规则 | `object_types/rpo_rule.bkn` | 定义恢复点目标所对应的时间窗口和适用条件，用于筛选不同等级的备份方案。 |

### 核心关系

| 关系 | 文件路径 | 说明 |
|------|----------|------|
| 应用能力约束计划 | `relation_types/app_capability_ct_backup_plan.bkn` | 通过 `app_backup_capability.app_object_type` 匹配 `backup_plan.app_object_type`，用于过滤应用不支持的计划条目。 |
| 应用配置规则约束方案 | `relation_types/app_config_rule_ct_backup_solution.bkn` | 通过 `app_config_rule.app_object_type` 与 `backup_solution.app_object_type` 建立应用配置规则项到候选备份方案的显式关系，用于先按应用类型收敛可匹配的应用配置规则集合，再由推荐 action 结合 `condition`、`config_category`、数据量、RPO 和数据分级等上下文完成二次筛选与去重。若规则项的 `app_object_type` 取值为 `all` 或逗号分隔的多应用类型列表，需由 action 继续解释其覆盖范围。 |
| 应用输出规则格式化方案 | `relation_types/app_output_rule_fmt_backup_solution.bkn` | 通过 `app_output_rule.app_object_type` 与 `backup_solution.app_object_type` 建立输出规则到候选备份方案的显式关系，用于按应用类型选择最终展示层级与章节命名规则。若未命中专属输出规则，推荐 action 应继续回退到 `app_object_type = all` 的通用输出规则。 |
| 应用规则约束方案 | `relation_types/app_rule_ct_backup_solution.bkn` | 通过 `app_rule.id` 匹配 `backup_solution.source_app_rule_id`，表示某条应用级规则可命中的备份方案集合。 |
| 方案使用备份计划 | `relation_types/backup_solution_use_backup_plan.bkn` | 通过 `backup_solution.id` 匹配 `backup_plan.backup_solution_id`，将备份方案关联到其下属的多个备份计划条目。 |
| 方案使用备份窗口 | `relation_types/backup_solution_use_backup_window.bkn` | 通过 `backup_solution.id` 匹配 `backup_window.backup_solution_id`，将备份方案关联到其下属的多个备份窗口条目。 |
| 方案使用保留计划 | `relation_types/backup_solution_use_retention_plan.bkn` | 通过 `backup_solution.id` 匹配 `retention_plan.backup_solution_id`，将备份方案关联到其唯一数据保留计划。 |
| 基础配置规则约束方案 | `relation_types/basic_config_rule_ct_backup_solution.bkn` | 通过 `basic_config_rule.app_object_type` 与 `backup_solution.app_object_type` 建立基础配置规则项到候选备份方案的显式关系，用于先按应用类型收敛可匹配的基础配置规则集合，再由推荐 action 结合 `condition`、`config_category`、数据量、RPO、备份窗口、行业类型和数据分级等上下文完成二次筛选与去重。若规则项的 `app_object_type` 取值为 `all` 或逗号分隔的多应用类型列表，需由 action 继续解释其覆盖范围。 |
| 数据规则约束方案 | `relation_types/data_rule_ct_backup_solution.bkn` | 通过 `data_rule.data_type` 匹配 `backup_solution.source_data_type`，将数据侧规则映射到对应的备份方案。 |
| 法规约束方案 | `relation_types/legal_rule_ct_backup_solution.bkn` | 通过 `legal_rule.id` 匹配 `backup_solution.source_legal_rule_id`，表示满足某条法规要求的备份方案集合。 |
| 候选推荐策略约束方案 | `relation_types/recommendation_policy_ct_backup_solution.bkn` | 通过 `recommendation_policy.industry_type` 与 `backup_solution.industry_type` 建立候选推荐策略到备份方案的显式关系，用于表达某个行业推荐策略主要评估哪些候选方案。该关系负责收敛行业专属方案集合；当未命中行业专属策略或需回退到通用策略时，仍应由推荐 action 继续选择 `industry_type = 未定义` 的默认推荐策略。 |
| RPO规则约束方案 | `relation_types/rpo_rule_ct_backup_solution.bkn` | 通过 `rpo_rule.id` 匹配 `backup_solution.source_rpo_rule_id`，表示满足某条 RPO 规则要求的备份方案集合。 |

## 目录结构

```
.
├── network.bkn
├── SKILL.md
├── CHECKSUM
├── object_types/
├── relation_types/
```

## 使用建议

### 查询场景

1. **获取所有对象定义**
   - 查看 `object_types/` 目录下的文件

2. **查找关系定义**
   - 查看 `relation_types/` 目录下的文件

## 索引表

### 按类型索引

- **对象定义**: `object_types/`
- **关系定义**: `relation_types/`

## 注意事项

1. 本网络由 BKN SDK 自动生成 SKILL.md
2. 所有定义遵循 BKN 规范
3. 使用 CHECKSUM 文件验证网络完整性
