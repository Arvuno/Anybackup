# `foundation-cli mysql recovery timepoint list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql recovery timepoint list --tenant-id <tenant-id> [--index <n>] [--count <n>] [--object-id <id>] [--storage-service-id <id>] [--storage-pool-id <id>] [--start-time <ts>] [--end-time <ts>] [--backup-task-type <type>] [--restore-gran <restore-gran>]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/time_points` |
| 风险 | `只读` |

`restore-gran` 允许值：`0/1/2/3`。

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
| index | `--index` | int | 否 | 默认值为 `0` |
| count | `--count` | int | 否 | 默认值为 `20`，取值范围 `1~100` |
| objectId | `--object-id` | string | 否 | MySQL 保护对象 ID |
| storageServiceId | `--storage-service-id` | string | 否 | 存储服务 ID |
| storagePoolId | `--storage-pool-id` | string | 否 | 存储池 ID |
| startTime | `--start-time` | string | 否 | 非负起始时间戳 |
| endTime | `--end-time` | string | 否 | 非负结束时间戳 |
| backupTaskType | `--backup-task-type` | string | 否 | 备份任务类型 |
| restoreGran | `--restore-gran` | string | 否 | 允许值： `0/1/2/3` |

## 响应参数

### 根字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int | 总记录数 |
| data | object[] | 时间点列表 |

### `data[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| objectId | string | 保护对象 ID |
| timePointId | string | 时间点 ID |
| timestamp | int64 | 时间戳 |
| display | int64 | 展示值 |
| displayTime | string | 可读展示时间 |
| dataSetId | string | 数据集 ID |
| storagePoolId | string | 存储池 ID |
| storagePoolType | int32 | 存储池类型 |
| storagePoolName | string | 存储池名称 |
| isClean | bool | 是否允许清理 |
| backupType | int32 | 备份类型 |
| business | int32 | 时间点产生方式 |
| timePointLimit | int32 | 时间点限制标记 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| fingerPrintId | string | 指纹库 ID |
| fingerPrintName | string | 指纹库名称 |

## 返回案例

说明：以下为已采集的成功响应示例。

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
        "timePointId": "f134a7e67138119fb50400163e33df73",
        "timestamp": 1775803115554800,
        "display": 1775803124000000,
        "displayTime": "2026-04-10T14:38:44.000000+08:00",
        "dataSetId": "bdc5695e9b30f8aa822cabd43f82f791",
        "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
        "storagePoolType": 2,
        "storagePoolName": "block",
        "isClean": false,
        "backupType": 1,
        "business": 1,
        "timePointLimit": 0,
        "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
        "storageServiceName": "StorageService_172.31.12.91",
        "expireDeleteTime": "",
        "expirationTimeType": 0,
        "expirationTime": 0,
        "isRetained": false,
        "operationType": 1,
        "fingerPrintId": "",
        "fingerPrintName": "",
        "usable": 0,
        "metaInfo": "{\"backupObjectId\":\"428e4c05bf87deef20a86e670f3eb9a9\",\"dataBackupEndTime\":\"1775803124000000\",\"dataCurBackupTp\":\"1775803122187868\",\"dataFullBackupEndTime\":\"1775803124000000\",\"dataFullBackupTp\":\"1775803122187868\",\"instanceName\":\"mysql3306_U0HYTDM3RENXS4EJ\",\"mergeBackup\":\"0\"}",
        "timePointType": 1
      }
    ],
    "totalNum": 1
  }
}
```

## 示例
```bash
foundation-cli mysql recovery timepoint list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --count 30
```




