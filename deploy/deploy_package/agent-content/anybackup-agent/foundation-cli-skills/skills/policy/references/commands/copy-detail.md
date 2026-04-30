# `foundation-cli policy copy detail`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy copy detail --tenant-id <tenant-id> --group-id <group-id>` |
| Method | `GET` |
| Path | `/api/sla/v1/group/{groupId}/copy_detail` |
| Risk | `read-only` |

## 公共引用

- 共享请求上下文：[`_shared-request-context.md`](./_shared-request-context.md)
- 通用响应包络：[`_shared-response-envelope.md`](./_shared-response-envelope.md)
- 通用枚举定义：[`_shared-enums.md`](./_shared-enums.md)
- 参数与枚举共性分析：[`_parameter-enum-common-analysis.md`](./_parameter-enum-common-analysis.md)

## 请求参数

### 共享参数

共享参数已统一收敛，详见：[`_shared-request-context.md`](./_shared-request-context.md)。

### 路径参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| groupId | `--group-id` | string | 是 | 计划组 ID |

## 返回结果

依据 `params/plan/copy_detail.go` 中 `CopyDetailResponse`，以及 `params/plan/backup_create.go`/`copy_create.go` 中嵌套类型。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息 |
| status | string | 请求状态 |
| responseData | object | 业务数据（`CopyDetailResponse`） |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | SLA ID |
| name | string | 策略名称 |
| remark | string | 备注 |
| validatePeriod | int | 有效期类型，见枚举 |
| periodType | int | 有效期周期类型，见枚举 |
| periodSpan | int | 周期跨度 |
| effectiveType | int | 生效类型，见枚举 |
| effectiveTime | int64 | 生效时间戳 |
| timeZone | string | 时区 |
| copyConfigs | object[] | 复制配置列表 |
| type | int | 业务类型，见枚举 |
| copyType | int | 复制类型，见枚举 |
| limitWindow | object | 限速窗口 |
| logRetention | object | 日志保留配置 |
| dataRetention | object | 数据保留配置 |
| bindResource | int | 绑定资源数量 |

### responseData.copyConfigs[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| copyId | string | 复制配置 ID |
| copyType | int | 复制类型，见枚举 |
| copyMode | int | 执行模式，见枚举 |
| copyConfig | object | 复制配置详情 |

### responseData.copyConfigs[].copyConfig 字段（核心）

`copyConfig` 对应 `CopyConfigResponse`，组合 `CopyPlan`、`CopyRule`、`CopyOption`、`CopySrcOption`、`CopyRetention` 等。

| 字段 | 类型 | 说明 |
|---|---|---|
| dayEnable/weekEnable/monthEnable/quarterEnable/yearEnable | bool | 周期计划开关 |
| dayConfiguration/weekConfiguration/monthConfiguration/quarterConfiguration/yearConfiguration | object | 周期配置 |
| copyOption | object | 目标端复制配置 |
| copySrcOption | object | 源端复制配置 |
| isRetentionOpen | bool | 是否开启保留 |
| retentionType | int | 保留类型，见枚举 |
| durationConfiguration | object | 按时长保留配置 |
| quantityNum | int | 按数量保留配置值 |
| limitWindow | object | 限速窗口 |
| dataType | int | 数据类型（0 数据，1 日志） |
| startOrigin | int64 | 起始时间点 |

### responseData.copyConfigs[].copyConfig 字段（对象展开）

以下对象结构来自 `params/plan/copy_detail.go` 与 `params/plan/copy_create.go`。

### copyConfig（`CopyConfigResponse`）

| 字段 | 类型 | 说明 |
|---|---|---|
| dayEnable | bool | 天计划开关 |
| dayConfiguration | `DayConfig` | 天计划配置 |
| weekEnable | bool | 周计划开关 |
| weekConfiguration | `WeekConfig` | 周计划配置 |
| monthEnable | bool | 月计划开关 |
| monthConfiguration | `MonthConfig` | 月计划配置 |
| quarterEnable | bool | 季度计划开关 |
| quarterConfiguration | `QuarterConfig` | 季度计划配置 |
| yearEnable | bool | 年计划开关 |
| yearConfiguration | `YearConfig` | 年计划配置 |
| dayRule | `DayRuleConfig` | 天级规则 |
| weekRule | `DayRuleConfig` | 周级规则 |
| monthRule | `MonthRuleConfig` | 月级规则 |
| quarterRule | `DayRuleConfig` | 季度级规则 |
| yearRule | `DayRuleConfig` | 年级规则 |
| copyOption | `CopyOptionDetail` | 目标端复制配置 |
| copySrcOption | `CopySrcOptionDetail` | 源端复制配置 |
| isRetentionOpen | bool | 是否开启保留 |
| retentionType | int | 保留类型，见枚举 `RetentionCType` |
| durationConfiguration | `DurationConfig` | 按时长保留配置 |
| quantityNum | int | 按数量保留值 |
| limitWindow | `LimitWindowConfig` | 限速窗口 |
| dataType | int | 数据类型（0 数据，1 日志） |
| startOrigin | int64 | 起始时间点 |

### copyConfig 计划对象（`CopyPlan`）

| 字段 | 类型 | 说明 |
|---|---|---|
| dayEnable | bool | 天计划开关 |
| dayConfiguration | `DayConfig` | 天计划配置 |
| weekEnable | bool | 周计划开关 |
| weekConfiguration | `WeekConfig` | 周计划配置 |
| monthEnable | bool | 月计划开关 |
| monthConfiguration | `MonthConfig` | 月计划配置 |
| quarterEnable | bool | 季度计划开关 |
| quarterConfiguration | `QuarterConfig` | 季度计划配置 |
| yearEnable | bool | 年计划开关 |
| yearConfiguration | `YearConfig` | 年计划配置 |

注：`DayConfig/WeekConfig/MonthConfig/QuarterConfig/YearConfig` 结构与 `plan backup detail` 中的对应对象一致。

### copyConfig 规则对象（`CopyRule`）

| 字段 | 类型 | 说明 |
|---|---|---|
| dayRule | `DayRuleConfig` | 天级规则 |
| weekRule | `DayRuleConfig` | 周级规则 |
| monthRule | `MonthRuleConfig` | 月级规则 |
| quarterRule | `DayRuleConfig` | 季度级规则 |
| yearRule | `DayRuleConfig` | 年级规则 |

### DayRuleConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| ruleType | int | 规则类型 |
| ruleValue | int | 规则值 |

### MonthRuleConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| ruleType | int | 规则类型 |
| ruleValue | int | 规则值 |
| isLastDayOfMonth | bool | 是否月末最后一天 |

### copyConfig 保留对象（`CopyRetention`）

| 字段 | 类型 | 说明 |
|---|---|---|
| isRetentionOpen | bool | 是否开启保留 |
| retentionType | int | 保留类型，见枚举 `RetentionCType` |
| durationConfiguration | `DurationConfig` | 按时长保留配置 |
| quantityNum | int | 按数量保留值 |

### copyOption（`CopyOptionDetail`）

| 字段 | 类型 | 说明 |
|---|---|---|
| destId | string | 目标端 ID |
| destName | string | 目标端名称 |
| destUsers | `CopyOptionDestUser[]` | 目标端用户列表 |
| tenantId | string | 租户 ID |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int | 存储池类型 |
| destCopyIpList | `DestCopyIpConfig[]` | 目标端复制 IP 列表 |
| srcCopySubnetId | string | 源端子网 ID |
| destCopySubnetId | string | 目标端子网 ID |
| isRetry | bool | 是否重试 |
| copyRetryConfig | `RetryConfig` | 重试配置 |
| timePointSyncType | bool | 时间点同步类型 |
| isMultiChannel | bool | 是否多通道复制 |
| channelCount | int | 通道数 |
| isCopyComplete | bool | 是否复制完备时间点 |
| isCycFullArch | bool | 是否周期性全量复制 |
| fullArchTpCount | int | 周期性全量复制时间点数量 |
| isCleanBackup | bool | 复制完成后是否清理备份时间点 |
| cleanBackupExpire | int | 清理备份时间点过期时长（小时） |

### copyOption.destUsers[]（`CopyOptionDestUser`）

| 字段 | 类型 | 说明 |
|---|---|---|
| destUserAkId | string | 目标端用户 AKID |
| destUserAk | string | 目标端用户 AK |
| destUserName | string | 目标端用户名 |
| destTenantName | string | 目标端租户名 |

### copyOption.destCopyIpList[]（`DestCopyIpConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| destCopyId | string | 目标复制节点 ID |
| destCopyIp | string | 目标复制节点 IP |

### copyOption.copyRetryConfig（`RetryConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| inMinute | int | 重试总分钟数 |
| minuteInterval | int | 重试间隔分钟数 |

### copySrcOption（`CopySrcOptionDetail`）

| 字段 | 类型 | 说明 |
|---|---|---|
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int | 存储池类型 |

### copyConfig.limitWindow（`LimitWindowConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| monday | `StartEndSpeed[]` | 周一限速窗口 |
| tuesday | `StartEndSpeed[]` | 周二限速窗口 |
| wednesday | `StartEndSpeed[]` | 周三限速窗口 |
| thursday | `StartEndSpeed[]` | 周四限速窗口 |
| friday | `StartEndSpeed[]` | 周五限速窗口 |
| saturday | `StartEndSpeed[]` | 周六限速窗口 |
| sunday | `StartEndSpeed[]` | 周日限速窗口 |

### StartEndSpeed

| 字段 | 类型 | 说明 |
|---|---|---|
| start | `WindowTriggerTimeConfig` | 开始时间 |
| end | `WindowTriggerTimeConfig` | 结束时间 |
| speed | int | 限速值 |
| speedUnit | int | 速率单位，见枚举 `SpeedUnit` |

### WindowTriggerTimeConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| hour | int | 小时 |
| minute | int | 分钟（0/30） |

## 枚举说明

本命令使用的枚举均为共享枚举，详见：[`_shared-enums.md`](./_shared-enums.md)。

字段与共享枚举映射：
- `validatePeriod` -> `ValidatePeriodType`
- `periodType` -> `ValidateType`
- `effectiveType` -> `EffectiveCType`
- `type` -> `BusinessType`
- `copyType` -> `CopyType`
- `copyMode` -> `ExecMode`
- `retentionType` -> `RetentionCType`
- `durationConfiguration.type` -> `RetentionDurationType`
- `*.mode` -> `BackupMode`
- `monthConfiguration.triggerType` -> `MonthTriggerType`
- `*.speedUnit` -> `SpeedUnit`

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "id": "sample-copy-policy-id",
    "name": "demo-copy-plan",
    "remark": "cross-region copy",
    "copyType": 1,
    "bindResource": 12,
    "copyConfigs": [
      {
        "copyId": "sample-copy-config-id",
        "copyType": 1,
        "copyMode": 1,
        "copyConfig": {
          "dayEnable": true,
          "dayConfiguration": {
            "interval": 1,
            "triggerTime": {
              "hour": 1,
              "min": 0
            },
            "mode": 1
          },
          "copyOption": {
            "storageServiceId": "sample-storage-service-id",
            "storagePoolId": "sample-storage-pool-id",
            "channelCount": 4
          }
        }
      }
    ]
  }
}
```

## 示例命令
```bash
foundation-cli policy copy detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --group-id <group-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli policy copy detail |
| 风险 | read-only |

## 示例

```bash
foundation-cli policy copy detail
```
