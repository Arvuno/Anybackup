# `foundation-cli timepoint list`

## 命令概览

| 项 | 值 |
|---|---|
| 命令 | `foundation-cli timepoint list --object-id <object-id> [filters]` |
| 请求方法 | `GET` |
| 接口路径 | `/backupmgm/v1/protect_object/{objectId}/time_points` |
| 风险级别 | `read-only` |

## 请求参数

### 必填参数

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 服务地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 秘密密钥 |
| `--object-id` | 是 | 保护对象 ID |

### 常用筛选参数

| 参数 | 说明 |
|---|---|
| `--business` | 允许值：`1`、`4`、`5` |
| `--start-time` | 非负时间戳 |
| `--end-time` | 非负时间戳 |
| `--storage-pool-id` | 存储池 ID |
| `--storage-service-id` | 存储服务 ID |
| `--data-set-id` | 数据集 ID |
| `--businesses` | 可重复传参，允许值：`1`、`4`、`5`、`6` |
| `--usable` | 允许值：`1`、`2` |
| `--backup-types` | 可重复传参，允许值：`1~7` |
| `--include-storage-types` | 可重复传参 |
| `--exclude-storage-types` | 可重复传参 |
| `--time-point-type` | 允许值：`0`、`1`、`2`、`3` |

## 说明

1. 本命令为查询命令，不发送请求体。
2. 可重复筛选参数需要重复传入，例如：`--backup-types 1 --backup-types 4`。
3. 枚举类筛选参数和时间戳参数会进行严格校验。

## 返回列表说明

顶层返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

`responseData` 结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 总记录数 |
| data | object[] | 时间点列表，元素结构见下表（`TimePointInfo`） |

`responseData.data[]`（`TimePointInfo`）字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| objectId | string | 保护对象唯一标识 |
| timePointId | string | 时间点唯一标识 |
| timestamp | int64 | 时间戳 |
| display | int64 | 前端展示值 |
| displayTime | string | 前端展示时间 |
| dataSetId | string | 存储数据集唯一标识 |
| storagePoolId | string | 存储池唯一标识 |
| storagePoolType | int | 存储池类型，见 `PoolTypeProto` |
| storagePoolName | string | 存储池名称 |
| isClean | bool | 是否可清理 |
| backupType | int32 | 备份方式，见“备份方式约定值” |
| business | int32 | 时间点产生方式，见“业务来源约定值” |
| timePointLimit | int | 时间点限制，见 `TimePointLimit` |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| expireDeleteTime | string | 强制保留周期到期删除时间 |
| expirationTimeType | int | 数据保留到期时间类型，见 `ExpirationTimeType` |
| expirationTime | int64 | 数据保留到期时间 |
| isRetained | bool | 是否处于强制保留期或有效期内 |
| operationType | int | 业务类型，见 `OperationTypeProto` |
| fingerPrintId | string | 指纹 ID |
| fingerPrintName | string | 指纹名称 |
| usable | int32 | 时间点可用性，见“可用性约定值” |
| metaInfo | string | 扩展元信息 |
| timePointType | int | 时间点类型，见 `TimePointType` |
| tapeExpirationTime | int64 | 归档后备份时间点过期时间（毫秒） |

## 枚举说明

### PoolTypeProto（存储池类型）

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `ALL` | 所有类型 |
| `1` | `SAN` | SAN 存储池 |
| `2` | `LOCAL` | 本地磁盘存储池 |
| `3` | `DISTRIBUTED` | 分布式存储池 |
| `4` | `NAS` | NAS 存储池 |
| `5` | `TAPE_DATA` | 磁带 |
| `6` | `NATIVETAPE` | 磁带原生格式存储池 |
| `7` | `LOCAL_FILE` | 本地文件系统存储池 |
| `8` | `OBJECT` | 对象存储池 |

### TimePointLimit（时间点限制）

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `NO_LIMIT` | 无限制，可执行恢复、远程复制、归档、克隆副本等操作 |
| `1` | `RECOVERY_REMOTE` | 可执行数据恢复、远程复制 |
| `2` | `ONLY_RECOVERY` | 仅支持数据恢复 |
| `3` | `LIMIT_MOUNT_RECOVERY` | 限制 NAS 挂载恢复使用 |

### ExpirationTimeType（数据保留到期时间类型）

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `UNKNOWN_EXPIRATION_TIME` | `---` |
| `1` | `PERMANENT_RETENTION` | 永久保留 |
| `2` | `EXPIRATION_TIME` | 过期时间 |

### OperationTypeProto（业务类型）

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `UNKNOWN` | 未知 |
| `1` | `BACKUP` | 备份 |
| `2` | `RECOVERY` | 恢复 |
| `3` | `CLEAN` | 清理 |
| `4` | `REMOTE` | 远程复制 |
| `5` | `REVERSE` | 反向复制 |
| `6` | `TAPE` | 正常归档 |
| `7` | `TAPE_REVERSE` | 从磁带复制到本地存储池 |
| `8` | `BROWSE` | 浏览 |
| `9` | `INTRA` | 域内复制 |
| `10` | `INTER_DOMAIN` | 域间复制（预留，暂未使用） |
| `12` | `MACHINE_DRILL` | 机台演练 |

### TimePointType（时间点类型）

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `UnknownTimePointType` | 未知时间点类型 |
| `1` | `DataTimePointType` | 数据时间点 |
| `2` | `LogTimePointType` | 日志时间点 |
| `3` | `SnapshotTimePointType` | 快照时间点 |

### 备份方式约定值（backupType）

| 枚举值 | 说明 |
|---|---|
| `1` | 完全备份 |
| `2` | 增量备份 |
| `3` | 差异备份 |
| `4` | 日志备份 |
| `5` | 永久增量 |
| `6` | 快照 |
| `7` | 快照备份 |

### 业务来源约定值（business）

| 枚举值 | 说明 |
|---|---|
| `1` | 本地备份 |
| `4` | 远程复制 |
| `5` | 反向复制 |

### 可用性约定值（usable）

| 枚举值 | 说明 |
|---|---|
| `-1` | 不可用 |
| `0` | 可用 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "objectId": "755aa409c7d8e7e153960556ecfab7bd",
        "timePointId": "f18a24df1e63117c9ef40050568943e2",
        "timestamp": 1785202639163050,
        "display": 1785202652140328,
        "displayTime": "2026-07-28T09:37:32.140328+08:00",
        "dataSetId": "053e3c70d069ccfa9afbdd8b426dbd37",
        "storagePoolId": "a001ce7215e211f1a88b0050568952bd",
        "storagePoolType": 2,
        "storagePoolName": "50sdc",
        "isClean": false,
        "backupType": 5,
        "business": 1,
        "timePointLimit": 1,
        "storageServiceId": "52247d9c121745eab3a7880817457df5",
        "storageServiceName": "StorageService_10.4.105.50",
        "expireDeleteTime": "",
        "expirationTimeType": 0,
        "expirationTime": 0,
        "isRetained": false,
        "operationType": 1,
        "fingerPrintId": "f18a1f06580511c28d4a0050568943e2",
        "fingerPrintName": "patch_test_1785200117",
        "usable": 0,
        "metaInfo": "{\"snapInfo\":\"\",\"snapType\":0}",
        "timePointType": 1,
        "tapeExpirationTime": 0
      },
      {
        "objectId": "755aa409c7d8e7e153960556ecfab7bd",
        "timePointId": "f18a2425d96a11169ef40050568943e2",
        "timestamp": 1785202328255120,
        "display": 1785202342790624,
        "displayTime": "2026-07-28T09:32:22.790624+08:00",
        "dataSetId": "053e3c70d069ccfa9afbdd8b426dbd37",
        "storagePoolId": "a001ce7215e211f1a88b0050568952bd",
        "storagePoolType": 2,
        "storagePoolName": "50sdc",
        "isClean": true,
        "backupType": 5,
        "business": 1,
        "timePointLimit": 1,
        "storageServiceId": "52247d9c121745eab3a7880817457df5",
        "storageServiceName": "StorageService_10.4.105.50",
        "expireDeleteTime": "",
        "expirationTimeType": 0,
        "expirationTime": 0,
        "isRetained": false,
        "operationType": 1,
        "fingerPrintId": "f18a1f06580511c28d4a0050568943e2",
        "fingerPrintName": "patch_test_1785200117",
        "usable": 0,
        "metaInfo": "{\"snapInfo\":\"\",\"snapType\":0}",
        "timePointType": 1,
        "tapeExpirationTime": 0
      },
      {
        "objectId": "755aa409c7d8e7e153960556ecfab7bd",
        "timePointId": "f18a1f1b44a711329ef40050568943e2",
        "timestamp": 1785200163086405,
        "display": 1785200183155408,
        "displayTime": "2026-07-28T08:56:23.155408+08:00",
        "dataSetId": "053e3c70d069ccfa9afbdd8b426dbd37",
        "storagePoolId": "a001ce7215e211f1a88b0050568952bd",
        "storagePoolType": 2,
        "storagePoolName": "50sdc",
        "isClean": true,
        "backupType": 1,
        "business": 1,
        "timePointLimit": 1,
        "storageServiceId": "52247d9c121745eab3a7880817457df5",
        "storageServiceName": "StorageService_10.4.105.50",
        "expireDeleteTime": "",
        "expirationTimeType": 0,
        "expirationTime": 0,
        "isRetained": false,
        "operationType": 1,
        "fingerPrintId": "f18a1f06580511c28d4a0050568943e2",
        "fingerPrintName": "patch_test_1785200117",
        "usable": 0,
        "metaInfo": "{\"snapInfo\":\"\",\"snapType\":0}",
        "timePointType": 1,
        "tapeExpirationTime": 0
      }
    ],
    "totalNum": 3
  }
}
```

## 示例
```bash
foundation-cli timepoint list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --businesses 1 \
  --businesses 4 \
  --backup-types 1 \
  --backup-types 4
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli timepoint list |
| 风险 | read-only |
