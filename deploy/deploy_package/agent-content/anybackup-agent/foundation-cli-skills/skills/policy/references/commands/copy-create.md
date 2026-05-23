# `foundation-cli policy copy create`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy copy create --data '<json>'` |
| Method | `POST` |
| Path | `/api/sla/v1/group/copy_info` |
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
| body | `--data` | JSON string | 是 | 复制计划创建请求体 |



| Field | Type | Required | Description |
|---|---|---|---|
| data | object | Yes | Request body JSON |

顶层对象：`AddSlaCopyRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 计划名称，`3~256` |
| id | string | 否 | SLA ID |
| templateId | string | 否 | 模板 ID |
| remark | string | 否 | 备注，`1~384` |
| type | int | 否 | 策略类型 |
| validatePeriod | int | 是 | 有效期类型 |
| periodType | int | 否 | 周期类型 |
| periodSpan | int | 否 | 周期跨度，`1~9999` |
| effectiveType | int | 是 | 生效类型 |
| effectiveTime | int64 | 否 | 生效时间戳（毫秒） |
| timeZone | string | 否 | 时区 |
| copyType | int | 是 | 复制类型 |
| copyConfigs | object[] | 是 | 复制策略配置集合 |
| limitWindow | object | 否 | 全局限速窗口 |
| logRetention | object | 否 | 日志保留配置 |
| DataRetention | object | 否 | 数据保留配置 |

`copyConfigs[]`（`CopyConfigCInfo`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| copyConfigs[].copyMode | int | 是 | 复制执行模式 |
| copyConfigs[].copyConfig | object | 是 | 单条复制配置 |

`copyConfigs[].copyConfig`（`CopyConfigC`）关键字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| copyConfigs[].copyConfig.dayEnable/weekEnable/monthEnable/quarterEnable/yearEnable | bool | 否 | 各计划层级开关 |
| copyConfigs[].copyConfig.dayConfiguration/weekConfiguration/monthConfiguration/quarterConfiguration/yearConfiguration | object | 否 | 各计划层级配置 |
| copyConfigs[].copyConfig.dayRule/weekRule/monthRule/quarterRule/yearRule | object | 否 | 复制规则 |
| copyConfigs[].copyConfig.copyOption | object | 否 | 复制目标配置 |
| copyConfigs[].copyConfig.copySrcOption | object | 否 | 复制源配置 |
| copyConfigs[].copyConfig.isRetentionOpen | bool | 否 | 保留策略开关 |
| copyConfigs[].copyConfig.retentionType | int | 否 | 保留类型 |
| copyConfigs[].copyConfig.durationConfiguration | object | 否 | 时长保留配置 |
| copyConfigs[].copyConfig.quantityNum | int | 否 | 副本数保留数量 |
| copyConfigs[].copyConfig.limitWindow | object | 否 | 限速窗口 |
| copyConfigs[].copyConfig.dataType | int | 否 | 数据类型（由后端定义） |
| copyConfigs[].copyConfig.startOrigin | int64 | 否 | 起始时间点 |

`copyConfigs[].copyConfig.copyOption`（`CopyOptionUpdateRequest`）展开字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| copyConfigs[].copyConfig.copyOption.destId | string | 否 | 目标域 ID |
| copyConfigs[].copyConfig.copyOption.destName | string | 否 | 目标域名称 |
| copyConfigs[].copyConfig.copyOption.tenantId | string | 否 | 目标租户 ID |
| copyConfigs[].copyConfig.copyOption.destUsers | object[] | 否 | 目标用户信息，最多 3 条 |
| copyConfigs[].copyConfig.copyOption.destUsers[].destUserAkId | string | 是（传 `destUsers` 时） | 目标用户 AK ID，长度 32 |
| copyConfigs[].copyConfig.copyOption.destUsers[].destUserAk | string | 是（传 `destUsers` 时） | 目标用户 AK |
| copyConfigs[].copyConfig.copyOption.storageServiceId | string | 是 | 目标存储服务 ID |
| copyConfigs[].copyConfig.copyOption.storageServiceName | string | 是 | 目标存储服务名称 |
| copyConfigs[].copyConfig.copyOption.storagePoolId | string | 是 | 目标存储池 ID |
| copyConfigs[].copyConfig.copyOption.storagePoolName | string | 是 | 目标存储池名称 |
| copyConfigs[].copyConfig.copyOption.storagePoolType | int | 是 | 目标存储池类型 |
| copyConfigs[].copyConfig.copyOption.srcCopySubnetId | string | 否 | 源端复制子网 ID |
| copyConfigs[].copyConfig.copyOption.destCopySubnetId | string | 否 | 目的端复制子网 ID |
| copyConfigs[].copyConfig.copyOption.destCopyIpList | object[] | 否 | 复制 IP 配置列表 |
| copyConfigs[].copyConfig.copyOption.destCopyIpList[].destCopyId | string | 否 | 复制网卡/通道 ID |
| copyConfigs[].copyConfig.copyOption.destCopyIpList[].destCopyIp | string | 否 | 复制 IP |
| copyConfigs[].copyConfig.copyOption.isRetry | bool | 否 | 是否启用复制失败重试 |
| copyConfigs[].copyConfig.copyOption.copyRetryConfig | object | 否 | 重试配置 |
| copyConfigs[].copyConfig.copyOption.copyRetryConfig.inMinute | int | 否 | 重试窗口分钟数，`1~60` |
| copyConfigs[].copyConfig.copyOption.copyRetryConfig.minuteInterval | int | 否 | 重试间隔分钟，`1~30` |
| copyConfigs[].copyConfig.copyOption.timePointSyncType | bool | 否 | 时间点同步类型 |
| copyConfigs[].copyConfig.copyOption.isMultiChannel | bool | 否 | 是否启用多通道 |
| copyConfigs[].copyConfig.copyOption.channelCount | int | 否 | 通道数量，`1~256` |
| copyConfigs[].copyConfig.copyOption.isCopyComplete | bool | 否 | 是否复制完整数据 |
| copyConfigs[].copyConfig.copyOption.isCycFullArch | bool | 否 | 是否开启周期性收敛 |
| copyConfigs[].copyConfig.copyOption.fullArchTpCount | int | 否 | 收敛时间点数量，`0~100` |
| copyConfigs[].copyConfig.copyOption.isCleanBackup | bool | 否 | 复制完成后是否清理备份点 |
| copyConfigs[].copyConfig.copyOption.cleanBackupExpire | int | 否 | 清理过期时长（小时） |

`copyConfigs[].copyConfig.copySrcOption`（`CopySrcOption`）展开字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| copyConfigs[].copyConfig.copySrcOption.storageServiceId | string | 否 | 源端存储服务 ID |
| copyConfigs[].copyConfig.copySrcOption.storageServiceName | string | 否 | 源端存储服务名称 |
| copyConfigs[].copyConfig.copySrcOption.storagePoolId | string | 否 | 源端存储池 ID |
| copyConfigs[].copyConfig.copySrcOption.storagePoolName | string | 否 | 源端存储池名称 |
| copyConfigs[].copyConfig.copySrcOption.storagePoolType | int | 否 | 源端存储池类型 |

`copyConfigs[].copyConfig.durationConfiguration`（`DurationConfig`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| copyConfigs[].copyConfig.durationConfiguration.num | int | 否 | 保留时长数值，`1~9999` |
| copyConfigs[].copyConfig.durationConfiguration.type | int | 否 | 保留时长单位 |

`limitWindow` / `copyConfigs[].copyConfig.limitWindow`（`LimitWindowConfig`）展开字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| limitWindow.monday~sunday | object[] | 否 | 顶层限速窗口列表，最多 24 段/天 |
| copyConfigs[].copyConfig.limitWindow.monday~sunday | object[] | 否 | 子复制策略限速窗口列表 |
| *.start.hour | int | 是（传窗口对象时） | 开始小时，`0~24` |
| *.start.minute | int | 是（传窗口对象时） | 开始分钟，`0/30` |
| *.end.hour | int | 是（传窗口对象时） | 结束小时，`0~24` |
| *.end.minute | int | 是（传窗口对象时） | 结束分钟，`0/30` |
| *.speed | int | 否 | 限速值，`0~10240` |
| *.speedUnit | int | 否 | 限速单位 |

`logRetention` / `DataRetention`（`CopyRetention`）字段：

| 字段路径 | 类型 | 必填 | 说明 |
|---|---|---|---|
| logRetention.isRetentionOpen | bool | 否 | 日志保留开关 |
| logRetention.retentionType | int | 否 | 日志保留类型 |
| logRetention.durationConfiguration | object | 否 | 日志按时长保留配置 |
| logRetention.durationConfiguration.num | int | 否 | 日志保留时长数值，`1~9999` |
| logRetention.durationConfiguration.type | int | 否 | 日志保留时长单位 |
| logRetention.quantityNum | int | 否 | 日志按副本数保留数量 |
| DataRetention.isRetentionOpen | bool | 否 | 数据保留开关 |
| DataRetention.retentionType | int | 否 | 数据保留类型 |
| DataRetention.durationConfiguration | object | 否 | 数据按时长保留配置 |
| DataRetention.durationConfiguration.num | int | 否 | 数据保留时长数值，`1~9999` |
| DataRetention.durationConfiguration.type | int | 否 | 数据保留时长单位 |
| DataRetention.quantityNum | int | 否 | 数据按副本数保留数量 |

## 枚举说明

共享枚举定义详见：[`_shared-enums.md`](./_shared-enums.md)。

字段与共享枚举映射：
- `type` -> `BusinessType`
- `validatePeriod` -> `ValidatePeriodType`
- `periodType` -> `ValidateType`
- `effectiveType` -> `EffectiveCType`
- `copyType` -> `CopyType`
- `copyConfigs[].copyMode` -> `ExecMode`（`1/2/3`）
- `copyConfigs[].copyConfig.retentionType` / `logRetention.retentionType` / `DataRetention.retentionType` -> `RetentionCType`
- `copyConfigs[].copyConfig.durationConfiguration.type` / `logRetention.durationConfiguration.type` / `DataRetention.durationConfiguration.type` -> `RetentionDurationType`
- `limitWindow.*[].speedUnit` -> `SpeedUnit`

命令特有补充（非共享枚举正文）：

### `copyConfigs[].copyMode` 扩展值

| 值 | 含义 |
|---|---|
| 4 | 后端扩展值（当前 `params` 未命名） |

## 返回结果

CLI 标准返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息（失败时） |
| responseData | object | 接口业务返回 |

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
    "groupId": "sample-group-id",
    "name": "demo-copy-plan"
  }
}
```

## 示例
```bash
foundation-cli policy copy create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"name":"demo-copy-plan","validatePeriod":1,"effectiveType":1,"copyType":1,"copyConfigs":[{"copyMode":1,"copyConfig":{"dayEnable":true,"dayConfiguration":{"interval":1,"triggerTime":{"hour":1,"min":0},"mode":1},"copyOption":{"storageServiceId":"<storage-service-id>","storageServiceName":"<storage-service-name>","storagePoolId":"<storage-pool-id>","storagePoolName":"<storage-pool-name>","storagePoolType":1}}}]}'
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli policy copy create |
| 风险 | write |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```

## Body 参数

| Field | Type | Required | Description |
|---|---|---|---|
| data | object | Yes | Request body JSON |
