# `foundation-cli policy backup create`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy backup create --data '<json>'` |
| Method | `POST` |
| Path | `/api/sla/v1/group/backup_info` |
| Risk | `write` |

## 公共引用

- 共享请求上下文：[`_shared-request-context.md`](./_shared-request-context.md)
- 通用响应包络：[`_shared-response-envelope.md`](./_shared-response-envelope.md)
- 通用枚举定义：[`_shared-enums.md`](./_shared-enums.md)
- 参数与枚举共性分析：[`_parameter-enum-common-analysis.md`](./_parameter-enum-common-analysis.md)

## 请求参数

### 共享参数

共享参数已统一收敛，详见：[`_shared-request-context.md`](./_shared-request-context.md)。

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| body | `--data` | JSON string | 是 | 备份计划创建请求体 |

## Body 参数

顶层对象：`AddSlaBackupRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 计划名称，`3~256` |
| remark | string | 否 | 备注，`1~384` |
| validatePeriod | int | 是 | 有效期类型 |
| periodType | int | 否 | 周期类型（自定义有效期时使用） |
| periodSpan | int | 否 | 周期跨度，`1~9999` |
| effectiveType | int | 是 | 生效类型 |
| effectiveTime | int64 | 否 | 生效时间戳（毫秒） |
| timeZone | string | 否 | 时区 |
| backupConfig | object | 是 | 备份策略配置 |

`backupConfig`（`BackupConfigC`）关键字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| backupConfig.minEnable | bool | 否 | 分钟级计划开关 |
| backupConfig.minConfiguration | object[] | 否 | 分钟级计划配置 |
| backupConfig.hourEnable | bool | 否 | 小时级计划开关 |
| backupConfig.hourConfiguration | object[] | 否 | 小时级计划配置 |
| backupConfig.dayEnable | bool | 否 | 天级计划开关 |
| backupConfig.dayConfiguration | object[] | 否 | 天级计划配置 |
| backupConfig.weekEnable | bool | 否 | 周级计划开关 |
| backupConfig.weekConfiguration | object[] | 否 | 周级计划配置 |
| backupConfig.monthEnable | bool | 否 | 月级计划开关 |
| backupConfig.monthConfiguration | object[] | 否 | 月级计划配置 |
| backupConfig.quarterEnable | bool | 否 | 季度级计划开关 |
| backupConfig.quarterConfiguration | object[] | 否 | 季度级计划配置 |
| backupConfig.yearEnable | bool | 否 | 年级计划开关 |
| backupConfig.yearConfiguration | object[] | 否 | 年级计划配置 |
| backupConfig.disableWindow | object | 否 | 禁用窗口配置（按周一至周日） |
| backupConfig.limitWindow | object | 否 | 限速窗口配置（按周一至周日） |
| backupConfig.isRetentionOpen | bool | 否 | 是否开启数据保留 |
| backupConfig.retentionType | int | 否 | 保留类型 |
| backupConfig.durationConfiguration | object | 否 | 按时长保留配置 |
| backupConfig.quantityNum | int | 否 | 按副本数保留数量 |
| backupConfig.isGfs | bool | 否 | 是否开启按特定周期保留副本 |
| backupConfig.gfsConfiguration | object | 否 | GFS 配置 |
| backupConfig.logRetentionOpen | bool | 否 | 是否开启日志保留 |
| backupConfig.logRetentionConfig | object | 否 | 日志保留时长配置 |

`backupConfig.limitWindow`（`LimitWindowConfig`）展开字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| backupConfig.limitWindow.monday~sunday | object[] | 否 | 每天限速窗口列表，最多 24 段 |
| backupConfig.limitWindow.*[].start.hour | int | 是 | 开始小时，`0~24` |
| backupConfig.limitWindow.*[].start.minute | int | 是 | 开始分钟，`0/30` |
| backupConfig.limitWindow.*[].end.hour | int | 是 | 结束小时，`0~24` |
| backupConfig.limitWindow.*[].end.minute | int | 是 | 结束分钟，`0/30` |
| backupConfig.limitWindow.*[].speed | int | 否 | 限速值，`0~10240` |
| backupConfig.limitWindow.*[].speedUnit | int | 否 | 限速单位 |

`backupConfig.durationConfiguration` / `backupConfig.logRetentionConfig`（`DurationConfig`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| backupConfig.durationConfiguration.num | int | 否 | 保留时长数值，`1~9999` |
| backupConfig.durationConfiguration.type | int | 否 | 保留时长单位 |
| backupConfig.logRetentionConfig.num | int | 否 | 日志保留时长数值，`1~9999` |
| backupConfig.logRetentionConfig.type | int | 否 | 日志保留时长单位 |

计划配置对象字段（`minConfiguration[]` / `hourConfiguration[]` / `dayConfiguration[]` / `weekConfiguration[]` / `monthConfiguration[]` / `quarterConfiguration[]` / `yearConfiguration[]`）：

| 对象 | 字段 | 类型 | 说明 |
|---|---|---|---|
| MinConfig | interval | int | 分钟间隔，`1~59` |
| MinConfig | triggerTime | object | 触发时间 `{hour,min}` |
| MinConfig | mode | int | 备份方式 |
| HourConfig | interval/min/mode | int | 小时级间隔、分钟、备份方式 |
| DayConfig | interval/triggerTime/mode | mixed | 天级配置 |
| WeekConfig | interval/triggerTime/mode/triggerList | mixed | 周级配置（`triggerList` 为 `1~7` 的数字数组，分别表示周一到周日，不能传字符串数组） |
| MonthConfig | interval/triggerTime/mode/triggerType | mixed | 月级配置 |
| MonthConfig | triggerDay / triggerWeek | object | 月级日期/周模式配置 |
| QuarterConfig | mode/triggerTime/triggerList/monthDay | mixed | 季度级配置 |
| YearConfig | mode/triggerTime/dayOfMonth | mixed | 年级配置 |

### 周期计划硬约束

- 某个周期的 `*Enable=true` 时，应同时传对应的 `*Configuration`，且数组不能为空；`*Enable=false` 时，对应 `*Configuration` 建议传空数组 `[]`。
- 周、月、季度涉及“列表”语义的字段必须传数字数组，不能传字符串数组，也不能把整个数组再包成字符串。
- `backupConfig.weekConfiguration[].triggerList` 为 `1~7` 的数字数组，分别表示周一到周日。
- `backupConfig.monthConfiguration[]` 的 `triggerType=1` 时，应使用 `triggerDay`，不要同时传 `triggerWeek`；`triggerDay.day` 为“每月第几天”的数字数组，例如 `[3]`、`[4]`、`[3,4]`。
- `backupConfig.monthConfiguration[]` 的 `triggerType` 若走“按周”模式，应使用 `triggerWeek.week` 和 `triggerWeek.weekDays` 两个数字数组表达“第几周的周几”。
- `backupConfig.quarterConfiguration[].triggerList` 为季度列表，按数字数组传递；`monthDay.month` 表示季度内第几个月，按数字传 `1~3`；`monthDay.day` 表示该月第几天。
- `backupConfig.yearConfiguration[].dayOfMonth.month` 表示自然月，按数字传 `1~12`；`dayOfMonth.day` 表示该月第几天。

## 请求结构体说明

除非场景另有说明，以下示例默认都使用 `validatePeriod=1`、`effectiveType=1`、`effectiveTime=0`；未启用的周期开关传 `false`，对应配置传空数组 `[]`。

### 按天执行计划

该场景表示：

1. 策略名称为 `备份策略`。
2. 只开启天级计划，其他分钟、小时、周、月、季度、年计划全部关闭。
3. 每 `1` 天执行一次，触发时间为 `14:33`，备份方式为 `mode=1`。
4. 未开启数据保留、日志保留、按特定周期保留副本、限速窗口和禁用窗口。

关键字段：

| 字段路径 | 示例值 | 说明 |
|---|---|---|
| `backupConfig.dayEnable` | `true` | 开启天级计划 |
| `backupConfig.dayConfiguration[0].interval` | `1` | 每 1 天执行一次 |
| `backupConfig.dayConfiguration[0].triggerTime.hour` | `14` | 触发小时 |
| `backupConfig.dayConfiguration[0].triggerTime.min` | `33` | 触发分钟 |
| `backupConfig.dayConfiguration[0].mode` | `1` | 备份方式 |
| `backupConfig.minEnable/hourEnable/weekEnable/monthEnable/quarterEnable/yearEnable` | `false` | 其他周期计划均关闭 |
| `backupConfig.isRetentionOpen` | `false` | 不开启数据保留 |
| `backupConfig.logRetentionOpen` | `false` | 不开启日志保留 |
| `backupConfig.isGfs` | `false` | 不开启按特定周期保留副本 |

对应请求体示例：

```json
{
  "name": "备份策略",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
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
    "isRetentionOpen": false,
    "logRetentionOpen": false,
    "isGfs": false,
    "limitWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    },
    "disableWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    }
  },
  "effectiveTime": 0
}
```

### 按周执行计划

该场景表示：

1. 策略名称为 `按周执行的备份策略`。
2. 只开启周级计划，其他分钟、小时、天、月、季度、年计划全部关闭。
3. 每 `1` 周执行一次，在周日 `09:00` 触发，备份方式为 `mode=1`。
4. `triggerList` 必须传数字数组，取值范围为 `1~7`，分别表示周一到周日；例如周日应传 `[7]`，不能传成 `["7"]` 或 `"[7]"`。
5. 不要把周计划写成 `weekday: 7`、`weekDay: 7`、`dayOfWeek: 7` 这类单值字段；即使只选一天，也必须写成 `triggerList: [7]`。

关键字段：

| 字段路径 | 示例值 | 说明 |
|---|---|---|
| `backupConfig.weekEnable` | `true` | 开启周级计划 |
| `backupConfig.weekConfiguration[0].interval` | `1` | 每 1 周执行一次 |
| `backupConfig.weekConfiguration[0].triggerList` | `[7]` | 周触发列表；`1~7` 分别表示周一到周日，必须是数字数组 |
| `backupConfig.weekConfiguration[0].triggerTime.hour` | `9` | 触发小时 |
| `backupConfig.weekConfiguration[0].triggerTime.min` | `0` | 触发分钟 |
| `backupConfig.weekConfiguration[0].mode` | `1` | 备份方式 |
| `backupConfig.minEnable/hourEnable/dayEnable/monthEnable/quarterEnable/yearEnable` | `false` | 其他周期计划均关闭 |
| `backupConfig.isRetentionOpen` | `false` | 不开启数据保留 |
| `backupConfig.logRetentionOpen` | `false` | 不开启日志保留 |
| `backupConfig.isGfs` | `false` | 不开启按特定周期保留副本 |

常见错误：

- 错误：`"weekday": 7`
- 正确：`"triggerList": [7]`

对应请求体示例：

```json
{
  "name": "按周执行的备份策略",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
  "backupConfig": {
    "minEnable": false,
    "minConfiguration": [],
    "hourEnable": false,
    "hourConfiguration": [],
    "dayEnable": false,
    "dayConfiguration": [],
    "weekEnable": true,
    "weekConfiguration": [
      {
        "interval": 1,
        "triggerList": [
          7
        ],
        "triggerTime": {
          "hour": 9,
          "min": 0
        },
        "mode": 1
      }
    ],
    "monthEnable": false,
    "monthConfiguration": [],
    "quarterEnable": false,
    "quarterConfiguration": [],
    "yearEnable": false,
    "yearConfiguration": [],
    "isRetentionOpen": false,
    "logRetentionOpen": false,
    "isGfs": false,
    "limitWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    },
    "disableWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    }
  },
  "effectiveTime": 0
}
```

### 所有周期执行

该场景表示：

1. 策略名称为 `所有周期执行`。
2. 本示例同时开启周、月、季度、年计划，分钟、小时、天计划关闭。
3. 周计划使用 `triggerList` 数字数组表达多个周几；月计划使用 `triggerType=1 + triggerDay` 表达“每月第几天”；季度计划使用 `triggerList + monthDay` 表达“哪些季度的季度内第几个月第几天”；年计划使用 `dayOfMonth` 表达“每年几月几日”。
4. 这些字段都要按原生 JSON 数组/数字传递，不能写成字符串，否则容易出现参数格式错误。

关键字段：

| 字段路径 | 示例值 | 说明 |
|---|---|---|
| `backupConfig.weekEnable` | `true` | 开启周计划 |
| `backupConfig.weekConfiguration[0].triggerList` | `[2,3,6,7]` | 周触发列表；`1~7` 分别表示周一到周日，必须是数字数组 |
| `backupConfig.monthEnable` | `true` | 开启月计划 |
| `backupConfig.monthConfiguration[0].triggerType` | `1` | 按“每月第几天”触发 |
| `backupConfig.monthConfiguration[0].triggerDay.day` | `[3]` | 每月第 3 天；必须是数字数组 |
| `backupConfig.monthConfiguration[1].triggerDay.day` | `[4]` | 每月第 4 天；必须是数字数组 |
| `backupConfig.monthConfiguration[2].triggerDay.day` | `[3,4]` | 每月第 3、4 天；支持同一配置传多个日期 |
| `backupConfig.quarterEnable` | `true` | 开启季度计划 |
| `backupConfig.quarterConfiguration[0].triggerList` | `[1]` | 季度列表；必须是数字数组 |
| `backupConfig.quarterConfiguration[0].monthDay.month` | `1` | 季度内第 1 个月 |
| `backupConfig.quarterConfiguration[0].monthDay.day` | `2` | 该月第 2 天 |
| `backupConfig.yearEnable` | `true` | 开启年计划 |
| `backupConfig.yearConfiguration[0].dayOfMonth.month` | `4` | 每年 4 月 |
| `backupConfig.yearConfiguration[0].dayOfMonth.day` | `26` | 每年 26 日 |

对应请求体示例：

```json
{
  "name": "所有周期执行",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
  "backupConfig": {
    "minEnable": false,
    "minConfiguration": [],
    "hourEnable": false,
    "hourConfiguration": [],
    "dayEnable": false,
    "dayConfiguration": [],
    "weekEnable": true,
    "weekConfiguration": [
      {
        "interval": 1,
        "triggerList": [
          2,
          3,
          6,
          7
        ],
        "triggerTime": {
          "hour": 9,
          "min": 4
        },
        "mode": 1
      }
    ],
    "monthEnable": true,
    "monthConfiguration": [
      {
        "interval": 1,
        "triggerType": 1,
        "triggerDay": {
          "day": [
            3
          ],
          "isLastDayOfMonth": false
        },
        "triggerTime": {
          "hour": 9,
          "min": 3
        },
        "mode": 1
      },
      {
        "interval": 1,
        "triggerType": 1,
        "triggerDay": {
          "day": [
            4
          ],
          "isLastDayOfMonth": false
        },
        "triggerTime": {
          "hour": 9,
          "min": 4
        },
        "mode": 1
      },
      {
        "interval": 1,
        "triggerType": 1,
        "triggerDay": {
          "day": [
            3,
            4
          ],
          "isLastDayOfMonth": false
        },
        "triggerTime": {
          "hour": 9,
          "min": 4
        },
        "mode": 2
      }
    ],
    "quarterEnable": true,
    "quarterConfiguration": [
      {
        "triggerList": [
          1
        ],
        "monthDay": {
          "month": 1,
          "day": 2
        },
        "triggerTime": {
          "hour": 9,
          "min": 4
        },
        "mode": 1
      }
    ],
    "yearEnable": true,
    "yearConfiguration": [
      {
        "dayOfMonth": {
          "month": 4,
          "day": 26
        },
        "triggerTime": {
          "hour": 9,
          "min": 5
        },
        "mode": 3
      }
    ],
    "isRetentionOpen": false,
    "logRetentionOpen": false,
    "isGfs": false,
    "limitWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    },
    "disableWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    }
  },
  "effectiveTime": 0
}
```

### 禁止备份窗口

该场景表示：

1. 策略名称为 `备份禁用窗口`。
2. 所有分钟、小时、天、周、月、季度、年执行计划均关闭。
3. 不开启数据保留、日志保留和按特定周期保留副本，`isRetentionOpen=false`、`logRetentionOpen=false`、`isGfs=false`。
4. 不配置限速窗口，`limitWindow.monday~sunday` 全部为空数组。
5. 仅配置禁用窗口：周四、周五、周六各有一段 `09:00-13:00` 的禁止备份时间，其他日期为空。

关键字段：

| 字段路径 | 示例值 | 说明 |
|---|---|---|
| `backupConfig.disableWindow.thursday[0].start.hour` | `9` | 周四禁用开始小时 |
| `backupConfig.disableWindow.thursday[0].start.minute` | `0` | 周四禁用开始分钟 |
| `backupConfig.disableWindow.thursday[0].end.hour` | `13` | 周四禁用结束小时 |
| `backupConfig.disableWindow.thursday[0].end.minute` | `0` | 周四禁用结束分钟 |
| `backupConfig.disableWindow.friday[0]` | 同周四配置 | 周五复用同一禁用规则 |
| `backupConfig.disableWindow.saturday[0]` | 同周四配置 | 周六复用同一禁用规则 |
| `backupConfig.limitWindow.monday~sunday` | `[]` | 未配置限速窗口 |
| `backupConfig.minEnable/hourEnable/dayEnable/weekEnable/monthEnable/quarterEnable/yearEnable` | `false` | 所有执行计划均关闭 |
| `backupConfig.isRetentionOpen` | `false` | 不开启数据保留 |
| `backupConfig.logRetentionOpen` | `false` | 不开启日志保留 |
| `backupConfig.isGfs` | `false` | 不开启按特定周期保留副本 |

对应请求体示例：

```json
{
  "name": "备份禁用窗口",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
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
    "isRetentionOpen": false,
    "logRetentionOpen": false,
    "isGfs": false,
    "limitWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    },
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
    }
  },
  "effectiveTime": 0
}
```

### 开启按特定周期保留副本和日志保留

该场景表示：

1. 策略名称为 `备份策略特定周期保留和日志保留`。
2. 所有分钟、小时、天、周、月、季度、年执行计划均关闭。
3. 开启数据保留且按副本数量保留：`isRetentionOpen=true`、`retentionType=2`、`quantityNum=1`。
4. 开启按特定周期保留副本：`isGfs=true`，并开启日/周/月/季/年保留，各周期保留数量均为 `1`。
5. 开启日志保留：`logRetentionOpen=true`，日志保留配置 `logRetentionConfig={num:1,type:4}`。
6. 不配置限速窗口和禁用窗口：`limitWindow`、`disableWindow` 全部为空数组。

关键字段：

| 字段路径 | 示例值 | 说明 |
|---|---|---|
| `backupConfig.isRetentionOpen` | `true` | 开启数据保留 |
| `backupConfig.retentionType` | `2` | 按副本数量保留 |
| `backupConfig.quantityNum` | `1` | 保留副本数量 |
| `backupConfig.isGfs` | `true` | 开启按特定周期保留副本 |
| `backupConfig.gfsConfiguration.gfsDayEnable~gfsYearEnable` | `true` | 各周期保留开关全开 |
| `backupConfig.gfsConfiguration.gfsDayNum~gfsYearNum` | `1` | 各周期保留数量均为 1 |
| `backupConfig.logRetentionOpen` | `true` | 开启日志保留 |
| `backupConfig.logRetentionConfig.num` | `1` | 日志保留时长数值 |
| `backupConfig.logRetentionConfig.type` | `4` | 日志保留时长单位 |
| `backupConfig.limitWindow.monday~sunday` | `[]` | 未配置限速窗口 |
| `backupConfig.disableWindow.monday~sunday` | `[]` | 未配置禁用窗口 |

对应请求体示例：

```json
{
  "name": "备份策略特定周期保留和日志保留",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
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
    "isRetentionOpen": true,
    "logRetentionOpen": true,
    "retentionType": 2,
    "isGfs": true,
    "quantityNum": 1,
    "gfsConfiguration": {
      "gfsDayNum": 1,
      "gfsWeekNum": 1,
      "gfsMonthNum": 1,
      "gfsQuarterNum": 1,
      "gfsYearNum": 1,
      "gfsDayEnable": true,
      "gfsWeekEnable": true,
      "gfsMonthEnable": true,
      "gfsQuarterEnable": true,
      "gfsYearEnable": true
    },
    "logRetentionConfig": {
      "num": 1,
      "type": 4
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
    "disableWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    }
  },
  "effectiveTime": 0
}
```

## 枚举说明

本命令使用的枚举均为共享枚举，详见：[`_shared-enums.md`](./_shared-enums.md)。

字段与共享枚举映射：
- `validatePeriod` -> `ValidatePeriodType`
- `periodType` -> `ValidateType`
- `effectiveType` -> `EffectiveCType`
- `backupConfig.retentionType` -> `RetentionCType`
- `backupConfig.durationConfiguration.type` / `backupConfig.logRetentionConfig.type` -> `RetentionDurationType`
- `backupConfig.*Configuration[].mode` -> `BackupMode`
- `backupConfig.monthConfiguration[].triggerType` -> `MonthTriggerType`
- `backupConfig.limitWindow.*[].speedUnit` -> `SpeedUnit`

## 返回结果

CLI 标准返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息（失败时） |
| responseData | object | 接口业务返回 |

## 返回案例

### 成功示例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "id": "bb90fed8e2264a1789f3474cf6e960f6"
  }
}
```

### 失败示例（策略名称重复）

```json
{
  "status": "failure",
  "error": {
    "errorCode": "HyperSLAMgm.PolicyNameAlreadyExists",
    "errorArgs": {
      "name": "备份数据所有"
    }
  },
  "responseData": null
}
```

## 请求体示例

```json
{
  "name": "备份策略",
  "remark": "",
  "validatePeriod": 1,
  "effectiveType": 1,
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
          "hour": 9,
          "min": 0
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
    "isRetentionOpen": false,
    "logRetentionOpen": false,
    "isGfs": false,
    "limitWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    },
    "disableWindow": {
      "monday": [],
      "tuesday": [],
      "wednesday": [],
      "thursday": [],
      "friday": [],
      "saturday": [],
      "sunday": []
    }
  },
  "effectiveTime": 0
}
```

## 示例

```bash
foundation-cli policy backup create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '<json>'
```

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy backup create --data '<json>'` |
| 风险 | `write` |
