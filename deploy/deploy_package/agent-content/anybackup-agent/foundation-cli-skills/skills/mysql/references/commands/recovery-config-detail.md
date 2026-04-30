# `foundation-cli mysql recovery-config detail`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql recovery-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/recovery_config?systemId=<system-id>&objectId=<object-id>` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| systemId | `--system-id` | string | 是 | MySQL 生产系统 ID |
| objectId | `--object-id` | string | 是 | MySQL 保护对象 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | MySQL 恢复配置 |

### `responseData` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int32 | 存储池类型 |
| storageServiceId | string | 存储服务 ID |
| timePointId | string | 数据时间点 ID |
| dataSetId | string | 数据集 ID |
| timestamp | int64 | 恢复时间戳 |
| displayTime | string | 展示时间 |
| logTimePointId | string | 日志时间点 ID |
| logSetId | string | 日志数据集 ID |
| logTimestamp | int64 | 日志时间戳 |
| restoreType | int32 | 恢复类型 |
| restoreToDay | int64 | 恢复目标日期 |
| restoreToTimeRange | string | 恢复时间范围 |
| restoreToTime | int64 | 恢复目标时间 |
| isRestoreNewInstance | bool | 是否恢复到新实例 |
| instanceName | string | 实例名称 |
| businessName | string | 业务名称 |
| restoreGran | int32 | 恢复粒度 |
| schemaList | string | 数据库或 Schema 列表 |
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientMac | string | 客户端 MAC |
| clientOSType | int64 | 客户端操作系统类型 |
| clientType | int32 | 客户端类型 |
| realPath | string | 数据源路径 |
| isRateControl | bool | 是否开启流量控制 |
| rateControlSize | int32 | 流量控制值 |
| isCoverOldDb | bool | 是否覆盖已有数据库 |
| isShutDown | bool | 恢复前是否执行停库 |
| databaseOnline | bool | 恢复后是否自动拉起数据库 |
| isMultiChannel | bool | 是否启用多通道恢复 |
| dataChannelNum | int32 | 数据通道数量 |
| hasCustomParam | bool | 是否启用自定义参数 |
| customParam | string | 自定义参数内容 |
| logRestoreParam | object | 日志恢复参数 |
| osUserName | string | 操作系统用户 |
| databasePort | int32 | 数据库端口 |
| isUpdatePassWord | bool | 是否更新密码 |
| instanceUser | string | 实例用户 |
| instancePwd | string | 实例密码 |
| specifyRestoreUser | bool | 是否显式指定恢复用户 |
| restoreUser | string | 恢复用户 |
| restorePwd | string | 恢复密码 |
| dataFilePath | string | 数据文件恢复路径 |
| cnfType | int32 | 配置文件类型 |
| cnfFilePath | string | 配置文件路径 |
| mysqlDefaultSpfileParam | string | 默认配置文件内容 |
| failureRetry | bool | 是否开启重试 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔 |
| srcObject | object | 恢复源对象 |
| targetObject | object | 恢复目标对象 |
| runImmediately | bool | 是否立即执行 |
| objectPassword | string | 恢复校验密码 |

### `logRestoreParam`

| 字段 | 类型 | 说明 |
|---|---|---|
| logTimestampRange | string | 日志恢复范围 |
| beginTimestamp | int64 | 开始时间 |
| endTimestamp | int64 | 结束时间 |
| archiveLogPath | string | 归档日志路径 |

### `srcObject`（列表展开）

- `systemId`（string）：恢复源系统 ID。
- `objectId`（string）：恢复源对象 ID。

### `targetObject`（列表展开）

- `systemId`（string）：恢复目标系统 ID。
- `objectId`（string）：恢复目标对象 ID。
- `name`（string）：恢复目标对象名称。
- `businessName`（string）：恢复目标对象业务名称。

## 返回案例

说明：以下为已采集的成功响应示例。

```json
{
    "status": "success",
    "error": null,
    "responseData": {
        "storagePoolType": 2,
        "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
        "storagePoolName": "block",
        "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
        "storageServiceName": "StorageService_172.31.12.91",
        "timePointId": "f134a7e67138119fb50400163e33df73",
        "dataSetId": "bdc5695e9b30f8aa822cabd43f82f791",
        "timestamp": 1775803115554800,
        "displayTime": "2026-04-10T14:38:44.000000+08:00",
        "logTimePointId": "",
        "logSetId": "",
        "logTimestamp": 0,
        "restoreType": 0,
        "restoreToDay": 0,
        "restoreToTimeRange": "",
        "restoreToTime": 0,
        "isRestoreNewInstance": false,
        "instanceName": "mysql3306_U0HYTDM3RENXS4EJ",
        "businessName": "mysql3306_U0HYTDM3RENXS4EJ",
        "restoreGran": 2,
        "schemaList": "",
        "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
        "realPath": "[\"/mysql3306_U0HYTDM3RENXS4EJ/test\"]",
        "isRateControl": false,
        "rateControlSize": 0,
        "isCoverOldDb": true,
        "isShutDown": true,
        "databaseOnline": true,
        "isMultiChannel": true,
        "dataChannelNum": 8,
        "hasCustomParam": false,
        "customParam": "",
        "logRestoreParam": {
            "logTimestampRange": "",
            "beginTimestamp": 0,
            "endTimestamp": 0,
            "archiveLogPath": ""
        },
        "osUserName": "mysql",
        "databasePort": 0,
        "isUpdatePassWord": false,
        "instanceUser": "",
        "instancePwd": "",
        "specifyRestoreUser": false,
        "restoreUser": "",
        "restorePwd": "",
        "dataFilePath": "/tmp",
        "cnfType": 0,
        "cnfFilePath": "",
        "mysqlDefaultSpfileParam": "[{\"type\":\"default\",\"key\":\"server_id\",\"value\":\"9999\"},{\"type\":\"default\",\"key\":\"log_bin\",\"value\":\"mysql-bin\"},{\"type\":\"default\",\"key\":\"binlog_format\",\"value\":\"ROW\"},{\"type\":\"default\",\"key\":\"log-error\",\"value\":\"mysqld.log\"},{\"type\":\"default\",\"key\":\"port\"}]",
        "failureRetry": true,
        "failureRetryCount": 1,
        "failureRetryInterval": 1,
        "srcObject": {
            "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
            "objectId": "428e4c05bf87deef20a86e670f3eb9a9"
        },
        "targetObject": {
            "systemId": "13caa58fd84bab1beb0b7c9d76b7efae",
            "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
            "name": "",
            "businessName": ""
        },
        "runImmediately": false
    }
}
```

## 示例
```bash
foundation-cli mysql recovery-config detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --system-id <system-id> \
  --object-id <object-id>
```



