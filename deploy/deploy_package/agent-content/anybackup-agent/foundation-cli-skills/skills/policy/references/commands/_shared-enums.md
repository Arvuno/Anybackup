# Plan Commands Shared Enums

适用范围：`skill/policy/references/commands` 下所有命令文档。

## ValidatePeriodType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 永久有效 |
| 2 | 自定义有效期 |

## ValidateType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 分钟 |
| 2 | 小时 |
| 3 | 天 |
| 4 | 周 |
| 5 | 月 |
| 6 | 季度 |
| 7 | 年 |

## EffectiveCType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 立即生效 |
| 2 | 自定义生效时间 |

## BusinessType

| 值 | 枚举名 | 含义（中文） |
|---|---|---|
| 0 | SlaGroup | SLA 类型（聚合） |
| 1 | SlaBackup | SLA 下的备份策略 |
| 2 | SlaRemoteSync | SLA 下远程复制策略 |
| 3 | SlaArchive | SLA 下归档策略 |
| 4 | StrategyK8sBackup | K8s 备份策略 |
| 5 | StrategyK8sSnapShot | K8s 快照策略 |
| 6 | SlaPreRecovery | SLA 下预恢复策略 |
| 7 | StrategyGather | 采集策略 |
| 8 | StrategyGatherCopy | 采集复制策略 |

## TemplateLevel

| 值 | 含义 |
|---|---|
| 0 | 一级策略 |
| 1 | 二级策略 |
| 2 | 三级策略 |

## ExecMode

| 值 | 含义 |
|---|---|
| 1 | 立即执行 |
| 2 | 按频率执行 |
| 3 | 按频率执行（指定副本） |

## RetentionCType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 按时长保留 |
| 2 | 按副本数量保留 |

## RetentionDurationType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 天 |
| 2 | 周 |
| 3 | 月 |
| 4 | 年 |

## BackupMode

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 完全备份 |
| 2 | 增量备份 |
| 3 | 差异备份 |
| 4 | 日志备份 |

## MonthTriggerType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 按日期触发 |
| 2 | 按第 N 周第 N 天触发 |

## SpeedUnit

| 值 | 含义 |
|---|---|
| 0 | MiB |
| 1 | KiB |

## CopyType

| 值 | 含义 |
|---|---|
| 1 | 域间复制/远程复制 |
| 2 | 域内复制 |

## ValidateStatusType

| 值 | 含义 |
|---|---|
| 0 | 未知 |
| 1 | 已生效 |
| 2 | 未生效 |
| 3 | 已到期 |
| 4 | 全部状态 |
| 5 | 审批中 |

## 缺失项补充说明

- `BusinessType`：在 `skill/policy/references/commands` 现有命令文档中，主要是英文枚举名（如 `SlaBackup`），没有完整中文释义；本文件已基于 `params/plan/copy_create.go` 注释补充中文含义。

