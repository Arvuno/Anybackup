# `foundation-cli policy backup detail`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy backup detail --tenant-id <tenant-id> --group-id <group-id>` |
| Method | `GET` |
| Path | `/api/sla/v1/group/{groupId}/backup_detail` |
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

依据 `params/plan/backup_detail.go` 中 `BackupDetailResponse`，以及 `params/plan/backup_create.go` 中嵌套类型。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息 |
| status | string | 请求状态 |
| responseData | object | 业务数据（`BackupDetailResponse`） |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | SLA ID |
| backupConfig | object | 备份策略配置 |
| name | string | 策略名称 |
| remark | string | 备注 |
| validatePeriod | int | 有效期类型，见枚举 |
| periodType | int | 有效期周期类型，见枚举 |
| periodSpan | int | 周期跨度 |
| effectiveType | int | 生效类型，见枚举 |
| effectiveTime | int64 | 生效时间戳 |
| type | int | 业务类型，见枚举 |
| timeZone | string | 时区 |

### responseData.backupConfig 字段（核心）

`backupConfig` 对应 `BackupConfigC`，由 `PlanCommon`、`WindowsCommon`、`RetentionCommon` 组合。

| 字段 | 类型 | 说明 |
|---|---|---|
| minEnable/hourEnable/dayEnable/weekEnable/monthEnable/quarterEnable/yearEnable | bool | 各周期计划开关 |
| minConfiguration/hourConfiguration/dayConfiguration/weekConfiguration/monthConfiguration/quarterConfiguration/yearConfiguration | object[] | 各周期触发配置 |
| disableWindow | object | 禁用窗口配置 |
| limitWindow | object | 限速窗口配置 |
| isRetentionOpen | bool | 是否开启保留 |
| retentionType | int | 保留类型，见枚举 |
| durationConfiguration | object | 按时长保留配置 |
| quantityNum | int | 按数量保留配置值 |
| isGfs | bool | 是否开启 GFS |
| gfsConfiguration | object | GFS 配置 |
| logRetentionOpen | bool | 是否开启日志保留 |
| logRetentionConfig | object | 日志保留配置 |

### responseData.backupConfig 字段（对象展开）

以下对象结构来自 `params/plan/backup_create.go`，用于补全 `backupConfig` 的嵌套字段。

### backupConfig（`BackupConfigC`）

| 字段 | 类型 | 说明 |
|---|---|---|
| minEnable | bool | 分钟级计划开关 |
| minConfiguration | `MinConfig[]` | 分钟级计划配置 |
| hourEnable | bool | 小时级计划开关 |
| hourConfiguration | `HourConfig[]` | 小时级计划配置 |
| dayEnable | bool | 天级计划开关 |
| dayConfiguration | `DayConfig[]` | 天级计划配置 |
| weekEnable | bool | 周级计划开关 |
| weekConfiguration | `WeekConfig[]` | 周级计划配置 |
| monthEnable | bool | 月级计划开关 |
| monthConfiguration | `MonthConfig[]` | 月级计划配置 |
| quarterEnable | bool | 季度级计划开关 |
| quarterConfiguration | `QuarterConfig[]` | 季度级计划配置 |
| yearEnable | bool | 年级计划开关 |
| yearConfiguration | `YearConfig[]` | 年级计划配置 |
| disableWindow | `DisableWindowConfig` | 禁用窗口 |
| limitWindow | `LimitWindowConfig` | 限速窗口 |
| isRetentionOpen | bool | 是否开启保留 |
| retentionType | int | 保留类型，见枚举 `RetentionCType` |
| durationConfiguration | `DurationConfig` | 按时长保留配置 |
| quantityNum | int | 按数量保留值 |
| isGfs | bool | 是否开启 GFS |
| gfsConfiguration | `GfsConfig` | GFS 保留配置 |
| logRetentionOpen | bool | 是否开启日志保留 |
| logRetentionConfig | `DurationConfig` | 日志保留时长配置 |

### backupConfig.minConfiguration[]（`MinConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| interval | int | 间隔 |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| mode | int | 备份模式，见枚举 `BackupMode` |

### backupConfig.hourConfiguration[]（`HourConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| interval | int | 间隔 |
| min | int | 分钟 |
| mode | int | 备份模式，见枚举 `BackupMode` |

### backupConfig.dayConfiguration[]（`DayConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| interval | int | 间隔 |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| mode | int | 备份模式，见枚举 `BackupMode` |

### backupConfig.weekConfiguration[]（`WeekConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| interval | int | 间隔 |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| mode | int | 备份模式，见枚举 `BackupMode` |
| triggerList | `int[]` | 周触发列表 |

### backupConfig.monthConfiguration[]（`MonthConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| interval | int | 间隔 |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| mode | int | 备份模式，见枚举 `BackupMode` |
| triggerType | int | 月触发类型，见枚举 `MonthTriggerType` |
| triggerDay | `MonthConfigDay` | 按“每月第几天”触发配置 |
| triggerWeek | `MonthConfigWeek` | 按“每月第几周第几天”触发配置 |

### backupConfig.quarterConfiguration[]（`QuarterConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| mode | int | 备份模式，见枚举 `BackupMode` |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| triggerList | `int[]` | 季度触发列表 |
| monthDay | `MonthDay` | 季度中月/日触发配置 |

### backupConfig.yearConfiguration[]（`YearConfig`）

| 字段 | 类型 | 说明 |
|---|---|---|
| mode | int | 备份模式，见枚举 `BackupMode` |
| triggerTime | `TriggerTimeConfig` | 触发时间 |
| dayOfMonth | `DayOfMonthC` | 年触发日配置 |

### 时间触发对象

#### TriggerTimeConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| hour | int | 小时 |
| min | int | 分钟 |

#### WindowTriggerTimeConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| hour | int | 小时 |
| minute | int | 分钟（0/30） |

### 月/季度/年触发对象

#### MonthConfigDay

| 字段 | 类型 | 说明 |
|---|---|---|
| day | `int[]` | 月内日期列表 |
| isLastDayOfMonth | bool | 是否月末最后一天 |

#### MonthConfigWeek

| 字段 | 类型 | 说明 |
|---|---|---|
| week | `int[]` | 第几周列表 |
| weekDays | `int[]` | 周几列表 |
| isLastWeekOfMonth | bool | 是否月末最后一周 |

#### MonthDay

| 字段 | 类型 | 说明 |
|---|---|---|
| day | int | 日 |
| month | int | 季度内第几个月 |
| isLastDayOfMonth | bool | 是否月末最后一天 |

#### DayOfMonthC

| 字段 | 类型 | 说明 |
|---|---|---|
| day | int | 日 |
| month | int | 月 |

### 窗口对象

#### DisableWindowConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| monday | `StartEnd[]` | 周一禁用窗口 |
| tuesday | `StartEnd[]` | 周二禁用窗口 |
| wednesday | `StartEnd[]` | 周三禁用窗口 |
| thursday | `StartEnd[]` | 周四禁用窗口 |
| friday | `StartEnd[]` | 周五禁用窗口 |
| saturday | `StartEnd[]` | 周六禁用窗口 |
| sunday | `StartEnd[]` | 周日禁用窗口 |

#### StartEnd

| 字段 | 类型 | 说明 |
|---|---|---|
| start | `WindowTriggerTimeConfig` | 开始时间 |
| end | `WindowTriggerTimeConfig` | 结束时间 |

#### LimitWindowConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| monday | `StartEndSpeed[]` | 周一限速窗口 |
| tuesday | `StartEndSpeed[]` | 周二限速窗口 |
| wednesday | `StartEndSpeed[]` | 周三限速窗口 |
| thursday | `StartEndSpeed[]` | 周四限速窗口 |
| friday | `StartEndSpeed[]` | 周五限速窗口 |
| saturday | `StartEndSpeed[]` | 周六限速窗口 |
| sunday | `StartEndSpeed[]` | 周日限速窗口 |

#### StartEndSpeed

| 字段 | 类型 | 说明 |
|---|---|---|
| start | `WindowTriggerTimeConfig` | 开始时间 |
| end | `WindowTriggerTimeConfig` | 结束时间 |
| speed | int | 限速值 |
| speedUnit | int | 速率单位，见枚举 `SpeedUnit` |

### 保留对象

#### DurationConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| num | int | 时长数值 |
| type | int | 时长单位，见枚举 `RetentionDurationType` |

#### GfsConfig

| 字段 | 类型 | 说明 |
|---|---|---|
| gfsDayEnable | bool | 开启日级 GFS |
| gfsDayNum | int | 日级保留数量 |
| gfsWeekEnable | bool | 开启周级 GFS |
| gfsWeekNum | int | 周级保留数量 |
| gfsMonthEnable | bool | 开启月级 GFS |
| gfsMonthNum | int | 月级保留数量 |
| gfsQuarterEnable | bool | 开启季度级 GFS |
| gfsQuarterNum | int | 季度级保留数量 |
| gfsYearEnable | bool | 开启年级 GFS |
| gfsYearNum | int | 年级保留数量 |

## 枚举说明

本命令使用的枚举均为共享枚举，详见：[`_shared-enums.md`](./_shared-enums.md)。

字段与共享枚举映射：
- `validatePeriod` -> `ValidatePeriodType`
- `periodType` -> `ValidateType`
- `effectiveType` -> `EffectiveCType`
- `type` -> `BusinessType`
- `*.mode` -> `BackupMode`
- `retentionType` -> `RetentionCType`
- `durationConfiguration.type` / `logRetentionConfig.type` -> `RetentionDurationType`
- `monthConfiguration[].triggerType` -> `MonthTriggerType`
- `*.speedUnit` -> `SpeedUnit`

## 返回案例

说明：
1. 以下为 3 组已验证的成功响应样例，分别对应“按天执行”、“保留策略和限速”、“禁止备份窗口”。
2. 字段含义以本页“返回结果”章节为准；具体取值受后端版本和环境配置影响。

### 按天执行备份策略详情

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "id": "b287446028c4420d8e4738531aa7b18a",
    "backupConfig": {
      "minEnable": false,
      "minConfiguration": [],
      "hourEnable": false,
      "hourConfiguration": [],
      "dayEnable": true,
      "dayConfiguration": [
        {
          "interval": 1,
          "triggerTime": {
            "hour": 14,
            "min": 33
          },
          "mode": 1
        }
      ],
      "weekEnable": false,
      "weekConfiguration": [],
      "monthEnable": false,
      "monthConfiguration": [],
      "quarterEnable": false,
      "quarterConfiguration": [],
      "yearEnable": false,
      "yearConfiguration": [],
      "disableWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [],
        "friday": [],
        "saturday": [],
        "sunday": []
      },
      "limitWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [],
        "friday": [],
        "saturday": [],
        "sunday": []
      },
      "isRetentionOpen": false,
      "retentionType": 0,
      "durationConfiguration": null,
      "quantityNum": 0,
      "isGfs": false,
      "gfsConfiguration": null,
      "logRetentionOpen": false,
      "logRetentionConfig": null
    },
    "name": "备份策略",
    "remark": "",
    "validatePeriod": 1,
    "periodType": 0,
    "periodSpan": 0,
    "effectiveType": 1,
    "effectiveTime": 1787036458095,
    "type": 1,
    "timeZone": ""
  }
}
```

### 备份策略保留策略和限速详情

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "id": "9e2144508ff54a70a94c1824d87012cc",
    "backupConfig": {
      "minEnable": false,
      "minConfiguration": [],
      "hourEnable": false,
      "hourConfiguration": [],
      "dayEnable": false,
      "dayConfiguration": [],
      "weekEnable": false,
      "weekConfiguration": [],
      "monthEnable": false,
      "monthConfiguration": [],
      "quarterEnable": false,
      "quarterConfiguration": [],
      "yearEnable": false,
      "yearConfiguration": [],
      "disableWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [],
        "friday": [],
        "saturday": [],
        "sunday": []
      },
      "limitWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [
          {
            "start": {
              "hour": 5,
              "minute": 30
            },
            "end": {
              "hour": 10,
              "minute": 0
            },
            "speed": 10,
            "speedUnit": 0
          }
        ],
        "friday": [
          {
            "start": {
              "hour": 5,
              "minute": 30
            },
            "end": {
              "hour": 10,
              "minute": 0
            },
            "speed": 10,
            "speedUnit": 0
          }
        ],
        "saturday": [
          {
            "start": {
              "hour": 5,
              "minute": 30
            },
            "end": {
              "hour": 10,
              "minute": 0
            },
            "speed": 10,
            "speedUnit": 0
          }
        ],
        "sunday": []
      },
      "isRetentionOpen": true,
      "retentionType": 1,
      "durationConfiguration": {
        "num": 1,
        "type": 4
      },
      "quantityNum": 0,
      "isGfs": false,
      "gfsConfiguration": null,
      "logRetentionOpen": false,
      "logRetentionConfig": null
    },
    "name": "备份策略保留策略和限速",
    "remark": "",
    "validatePeriod": 1,
    "periodType": 0,
    "periodSpan": 0,
    "effectiveType": 1,
    "effectiveTime": 1787037129759,
    "type": 1,
    "timeZone": ""
  }
}
```

### 禁止备份策略详情

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "id": "a230e2b778a142a1a53702b2e002b01f",
    "backupConfig": {
      "minEnable": false,
      "minConfiguration": [],
      "hourEnable": false,
      "hourConfiguration": [],
      "dayEnable": false,
      "dayConfiguration": [],
      "weekEnable": false,
      "weekConfiguration": [],
      "monthEnable": false,
      "monthConfiguration": [],
      "quarterEnable": false,
      "quarterConfiguration": [],
      "yearEnable": false,
      "yearConfiguration": [],
      "disableWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [
          {
            "start": {
              "hour": 9,
              "minute": 0
            },
            "end": {
              "hour": 13,
              "minute": 0
            }
          }
        ],
        "friday": [
          {
            "start": {
              "hour": 9,
              "minute": 0
            },
            "end": {
              "hour": 13,
              "minute": 0
            }
          }
        ],
        "saturday": [
          {
            "start": {
              "hour": 9,
              "minute": 0
            },
            "end": {
              "hour": 13,
              "minute": 0
            }
          }
        ],
        "sunday": []
      },
      "limitWindow": {
        "monday": [],
        "tuesday": [],
        "wednesday": [],
        "thursday": [],
        "friday": [],
        "saturday": [],
        "sunday": []
      },
      "isRetentionOpen": false,
      "retentionType": 0,
      "durationConfiguration": null,
      "quantityNum": 0,
      "isGfs": false,
      "gfsConfiguration": null,
      "logRetentionOpen": false,
      "logRetentionConfig": null
    },
    "name": "备份禁用窗口",
    "remark": "",
    "validatePeriod": 1,
    "periodType": 0,
    "periodSpan": 0,
    "effectiveType": 1,
    "effectiveTime": 1787037429819,
    "type": 1,
    "timeZone": ""
  }
}
```

## 示例命令
```bash
foundation-cli policy backup detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --group-id <group-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli policy backup detail |
| 风险 | read-only |

## 示例

```bash
foundation-cli policy backup detail
```
