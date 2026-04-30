# `foundation-cli protect backup-config get`

## 命令摘要

查询指定保护对象的备份配置。

## 示例

```bash
foundation-cli protect backup-config get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect backup-config get --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v2/protect_object/{objectId}/backup_config` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `8.0.9.0` |

### 业务参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | 保护对象 ID |

### Query 参数

无。

## 返回结果

来源：`params/protect/backup_config.go` -> `GetMultipleRegionConfigResponse`  
嵌套对象定义复用：`params/protect/create_backup_config.go`

### 顶层 envelope

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | any | 错误信息，成功通常为 `null` |
| responseData | object | 业务数据，对应 `GetMultipleRegionConfigResponse` |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| isConfig | bool | 是否已配置备份配置 |
| backupDestination | object | 备份目的地配置，对应 `BackupDestination` |
| conventionalConfig | object | 常规配置，对应 `ConventionalConfig` |

### `responseData.backupDestination` / `BackupDestination`

| 字段 | 类型 | 说明 |
|---|---|---|
| multipleRegions | bool | 是否启用多区域 |
| regionParams | object[] | 区域配置列表 |

### `responseData.backupDestination.regionParams[]` / `RegionParam`

| 字段 | 类型 | 说明 |
|---|---|---|
| regionName | string | 区域名称 |
| clientParams | object[] | 客户端配置列表 |
| storageAvailable | bool | 是否启用故障切换 |
| triggerCondition | object | 故障切换触发条件 |
| storagePoolParams | object[] | 存储池配置列表 |

### `clientParams[]` / `ClientParams`

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端唯一标识 |
| clientName | string | 客户端名称 |
| clientOsType | int | 客户端操作系统类型 |
| clientIp | string | 客户端 IP |
| clientStatus | int | 客户端状态 |

### `triggerCondition` / `TriggerCondition`

| 字段 | 类型 | 说明 |
|---|---|---|
| taskFailedCount | int32 | 作业失败触发次数 |

### `storagePoolParams[]` / `StoragePoolParam`

| 字段 | 类型 | 说明 |
|---|---|---|
| isMasterStorage | bool | 是否主存储池 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int | 存储池类型 |
| encryptionStorage | int32 | 是否开启存储加密 |
| encryptionLocation | int32 | 加密位置 |
| encryptionThreadNum | int32 | 加密线程数 |
| encryptAlgo | int32 | 加密算法 |
| compress | int32 | 是否开启压缩 |
| compressLocation | int32 | 压缩位置 |
| compressThreadNum | int32 | 压缩线程数 |
| compressAlgorithm | int32 | 压缩算法 |
| deduplication | int32 | 是否开启去重 |
| deduplicationLocation | int32 | 去重位置 |
| deduplicationThreadNum | int32 | 去重线程数 |
| fingerLibraryBindType | int32 | 指纹库绑定类型 |
| fingerLibraryId | string | 指纹库 ID |
| fingerLibraryConfig | object | 指纹库配置 |
| dataConsistencySwitch | bool | 数据一致性开关 |
| dataConsistencyLogic | int32 | 数据一致性算法 |

### `fingerLibraryConfig` / `FingerLibraryConfig`

| 字段 | 类型 | 说明 |
|---|---|---|
| name | string | 指纹库名称 |
| sliceMinValue | int32 | 切片最小值 |
| sliceMaxValue | int32 | 切片最大值 |
| dataSize | int32 | 预计数据量 |
| compressAlgorithm | int32 | 指纹库存储压缩算法 |
| creatorName | string | 创建人名称 |

### `conventionalConfig` / `ConventionalConfig`

| 字段 | 类型 | 说明 |
|---|---|---|
| permanentIncrement | int32 | 是否开启永久增量 |
| failureRetry | int32 | 是否开启失败重试 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔，单位分钟 |
| encryptionTrans | int32 | 是否开启传输加密 |
| forcedRetentionSwitch | bool | 是否开启强制数据保留周期 |
| forcedRetentionCycle | int32 | 强制保留天数 |

## 枚举类型列表

与 `params/protect/create_backup_config.go` 中的 `PoolTypeProto`、加密/压缩/去重、指纹库绑定和一致性字段定义一致。

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
    "isConfig": true,
    "backupDestination": {
      "multipleRegions": false,
      "regionParams": [
        {
          "regionName": "default",
          "clientParams": [
            {
              "clientId": "sample-client-id",
              "clientName": "mysql-node-01",
              "clientOsType": 2,
              "clientIp": "10.10.10.10",
              "clientStatus": 1
            }
          ],
          "storageAvailable": false,
          "triggerCondition": {
            "taskFailedCount": 3
          },
          "storagePoolParams": [
            {
              "isMasterStorage": true,
              "storageServiceId": "sample-storage-service-id",
              "storageServiceName": "StorageService_A",
              "storagePoolId": "sample-storage-pool-id",
              "storagePoolName": "pool-a",
              "storagePoolType": 1,
              "dataConsistencySwitch": true,
              "dataConsistencyLogic": 1
            }
          ]
        }
      ]
    },
    "conventionalConfig": {
      "permanentIncrement": 1,
      "failureRetry": 1,
      "failureRetryCount": 3,
      "failureRetryInterval": 5,
      "encryptionTrans": 1
    }
  }
}
```

## 示例命令
```bash
foundation-cli protect backup-config get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 当前环境验证

- 已使用 `object-id=428e4c05bf87deef20a86e670f3eb9a9` 实测成功。
- 当前环境返回 `isConfig=true`，主存储池为 `block`。
