# `foundation-cli protect policy bind-list`

## 命令摘要

查询指定保护对象已绑定的策略列表。

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect policy bind-list --object-id <object-id> [--index <n>] [--count <n>] [--filter <text>] [--validate-status <1|2|3|4|5>] [--disable-mark <-1|0|1|2>] [--type <business-type>]` |
| Method | `GET` |
| Path | `/api/sla/v1/group/object/{objectId}/templates` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 必填 | 说明 |
|---|---|---|---|
| tenantId | `--tenant-id` | 是 | 租户标识 |
| endpoint | `--endpoint` | 是 | Foundation 服务地址 |
| ak | `--ak` | 是 | Access Key |
| sk | `--sk` | 是 | Secret Key |
| targetVersion | `--target-version` | 否 | 目标版本，默认 `8.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | 保护对象 ID |
| index | `--index` | int | 否 | 分页起始，默认 `0` |
| count | `--count` | int | 否 | 分页条数，默认 `10` |
| filter | `--filter` | string | 否 | 模糊过滤 |
| validateStatus | `--validate-status` | string | 否 | 生效状态，`1/2/3/4/5` |
| disableMark | `--disable-mark` | string | 否 | 禁用状态，`-1/0/1/2` |
| type | `--type` | string | 否 | 业务类型 |

## 返回参数

CLI 外层结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

`responseData.data[]`（`GetTemplateListByObjectIdResponse`）核心字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | SLA ID |
| name | string | SLA 名称 |
| type | int | 业务类型 |
| validateStatus | int | 生效状态 |
| status | int | 是否禁用 |
| expiredTime | string | 到期时间展示值 |
| execMode | int | 数据执行模式 |
| logExecMode | int | 日志执行模式 |
| isIncludeLog | bool | 是否包含日志 |
| dataCopyRule | object | 数据复制规则 |
| logCopyRule | object | 日志复制规则 |
| logPlanCommon | object | 日志计划配置 |
| logRetention | object | 日志保留策略 |
| dataRetention | object | 数据保留策略 |
| snapshotRetention | object | 快照保留策略 |
| preRecoveryRetention | object | 预恢复保留策略 |
| copyType | int | 复制类型 |
| copyDestInfos | object[] | 复制目标信息 |
| children | object[] | 子策略列表（递归） |

`children[]`（`ChildGetTemplateListByObjectIdResponse`）字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | 子模板 ID |
| name | string | 子模板名称 |
| level | int | 层级 |
| parentId | string | 父模板 ID |
| creatorName | string | 创建者 |
| type | int | 业务类型 |
| children | object[] | 子节点 |

`copyDestInfos[]` 字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| dataType | int | 数据类型 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int | 存储池类型 |

## 枚举说明

### `validateStatus`

| 值 | 含义 |
|---|---|
| 0 | Unknown |
| 1 | 已生效 |
| 2 | 未生效 |
| 3 | 已到期 |
| 4 | 全部状态 |
| 5 | 审批中 |

### `type`（`BusinessType`）

| 值 | 含义 |
|---|---|
| 0 | SlaGroup |
| 1 | SlaBackup |
| 2 | SlaRemoteSync |
| 3 | SlaArchive |
| 4 | StrategyK8sBackup |
| 5 | StrategyK8sSnapShot |
| 6 | SlaPreRecovery |
| 7 | StrategyGather |
| 8 | StrategyGatherCopy |

### `level`（`TemplateLevel`）

| 值 | 含义 |
|---|---|
| 0 | 一级策略 |
| 1 | 二级策略 |
| 2 | 三级策略 |

### `execMode` / `logExecMode`（`ExecMode`）

| 值 | 含义 |
|---|---|
| 1 | 立即执行 |
| 2 | 按频率执行 |
| 3 | 按频率执行（指定副本） |

### `copyType`

| 值 | 含义 |
|---|---|
| 1 | 域间复制 / 远程复制 |
| 2 | 域内复制 |

### `copyDestInfos.dataType`

| 值 | 含义 |
|---|---|
| 0 | 数据 |
| 1 | 日志 |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "id": "a15e75fbdab04effbc7bf076aac8c7a8",
        "name": "备份策略",
        "type": 1,
        "validateStatus": 1,
        "status": 0,
        "expiredTime": "-",
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
        "snapshotMinEnable": false,
        "snapshotMinConfiguration": null,
        "snapshotDayEnable": false,
        "snapshotDayConfiguration": null,
        "dataCopyRule": {
          "dayRule": null,
          "weekRule": null,
          "monthRule": null
        },
        "logCopyRule": {
          "dayRule": null,
          "weekRule": null,
          "monthRule": null
        },
        "children": null,
        "logPlanCommon": {
          "minEnable": false,
          "minConfiguration": null,
          "hourEnable": false,
          "hourConfiguration": null,
          "dayEnable": false,
          "dayConfiguration": null,
          "weekEnable": false,
          "weekConfiguration": null,
          "monthEnable": false,
          "monthConfiguration": null,
          "quarterEnable": false,
          "quarterConfiguration": null,
          "yearEnable": false,
          "yearConfiguration": null
        },
        "logExecMode": 0,
        "isIncludeLog": false,
        "execMode": 2
      }
    ],
    "totalNum": 0
  }
}
```

## 示例
```bash
foundation-cli protect policy bind-list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --count 20
```

## 当前环境验证

- 已使用 `object-id=428e4c05bf87deef20a86e670f3eb9a9` 实测成功。
- 当前环境返回 `totalNum=0`。

