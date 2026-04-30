# `foundation-cli mysql recovery-detail`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql recovery-detail --tenant-id <tenant-id> --task-id <task-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/get_recovery_task_detail?taskId=<task-id>` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `9.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| taskId | `--task-id` | string | 是 | MySQL 恢复任务 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | MySQL 恢复任务详情 |

### `responseData` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int32 | 存储池类型 |
| storageServiceId | string | 存储服务 ID |
| timePointId | string | 数据时间点 ID |
| dataSetId | string | 数据集 ID |
| timestamp | int64 | 数据时间戳 |
| displayTime | string | 展示时间 |
| logTimePointId | string | 日志时间点 ID |
| logSetId | string | 日志数据集 ID |
| logTimestamp | int64 | 日志时间戳 |
| restoreType | int32 | 恢复类型 |
| restoreToTime | int64 | 恢复目标时间 |
| restoreGran | int32 | 恢复粒度 |
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| realPath | string | 数据源路径 |
| isRateControl | bool | 是否开启流量控制 |
| rateControlSize | int32 | 流量控制值 |
| isCoverOldDb | bool | 是否覆盖已有数据库 |
| isShutDown | bool | 恢复前是否执行停库 |
| databaseOnline | bool | 恢复后是否自动拉起数据库 |
| hasCustomParam | bool | 是否启用自定义参数 |
| customParam | string | 自定义参数内容 |
| failureRetry | bool | 是否开启重试 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔 |
| runImmediately | bool | 是否立即执行 |

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
    "timePointId": "",
    "dataSetId": "bdc5695e9b30f8aa822cabd43f82f791",
    "timestamp": 1775803115554800,
    "displayTime": "2026-04-10T14:38:44.000000+08:00",
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
    "mysqlDefaultSpfileParam": "",
    "failureRetry": true,
    "failureRetryCount": 1,
    "failureRetryInterval": 1,
    "runImmediately": false
  }
}
```

## 示例
```bash
foundation-cli mysql recovery-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --task-id <task-id>
```



