# `foundation-cli mysql backup-config detail`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql backup-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/app_backup_config?systemId=<system-id>&objectId=<object-id>` |
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

### 根字段

| 字段 | 类型 | 说明 |
|---|---|---|
| realPath | string | 数据源路径 |
| backupType | int32 | 备份类型 |
| backupGran | int32 | 备份粒度 |
| isMergeBackup | bool | 是否启用合成备份 |
| isRealTimeLog | bool | 是否启用实时日志 |
| isBackupMode | bool | 是否启用备份模式 |
| backupMode | int32 | 备份模式 |
| isNodeSequence | bool | 是否启用节点顺序执行 |
| nodeSequence | string | 节点顺序 |
| isTableRestore | bool | 是否支持表级恢复 |
| isArchDelete | bool | 是否启用归档删除 |
| deleteArchUnit | int32 | 归档删除单位 |
| deleteArchPeriod | int32 | 归档删除周期 |
| isMultiChannel | bool | 是否启用多通道备份 |
| dataChannelNum | int32 | 数据通道数量 |
| isNotAutoToFull | bool | 是否禁用自动转全量 |
| lvmChunkSliceSize | int32 | 分块大小 |
| isForwardSelect | bool | 是否启用正向选择 |
| isBackupCompress | bool | 是否启用高级压缩 |
| isRateControl | bool | 是否开启流量控制 |
| rateSize | int32 | 流量控制值 |
| volumeInfo | string | 卷信息 |
| hasCustomParam | bool | 是否启用自定义参数 |
| customParam | string | 自定义参数内容 |
| usePageTrack | bool | 是否启用 PageTrack |
| cleanPageTrack | bool | 是否启用 PageTrack 清理 |
| commonConfigParams | object | 共享配置块 |

### `commonConfigParams`

| 字段 | 类型 | 说明 |
|---|---|---|
| appType | string | 应用类型 |
| overwriteConfig | bool | 是否覆盖已有配置 |
| backupDestination | object | 备份目标配置 |
| conventionalConfig | object | 常规备份配置 |

### `commonConfigParams.backupDestination`（列表展开）

- `multipleRegions`（bool）：是否启用多区域。
- `regionParams`（object[]）：区域配置列表。
- `regionParams[].regionName`（string）：区域名称。
- `regionParams[].clientParams`（object[]）：客户端配置列表。
- `regionParams[].clientParams[].clientId`（string）：客户端唯一标识。
- `regionParams[].clientParams[].clientName`（string）：客户端名称。
- `regionParams[].clientParams[].clientOsType`（int）：客户端操作系统类型。
- `regionParams[].clientParams[].clientIp`（string）：客户端 IP。
- `regionParams[].clientParams[].clientStatus`（int）：客户端状态。
- `regionParams[].storageAvailable`（bool）：是否启用故障切换。
- `regionParams[].triggerCondition`（object）：故障切换触发条件。
- `regionParams[].triggerCondition.taskFailedCount`（int）：作业失败次数阈值。
- `regionParams[].storagePoolParams`（object[]）：存储池配置列表。
- `regionParams[].storagePoolParams[].isMasterStorage`（bool）：是否主存储池。
- `regionParams[].storagePoolParams[].storageServiceId`（string）：存储服务 ID。
- `regionParams[].storagePoolParams[].storageServiceName`（string）：存储服务名称。
- `regionParams[].storagePoolParams[].storagePoolId`（string）：存储池 ID。
- `regionParams[].storagePoolParams[].storagePoolName`（string）：存储池名称。
- `regionParams[].storagePoolParams[].storagePoolType`（int）：存储池类型。
- `regionParams[].storagePoolParams[].encryptionStorage`（int）：是否开启存储加密。
- `regionParams[].storagePoolParams[].encryptionLocation`（int）：加密位置。
- `regionParams[].storagePoolParams[].encryptionThreadNum`（int）：加密线程数。
- `regionParams[].storagePoolParams[].encryptAlgo`（int）：加密算法。
- `regionParams[].storagePoolParams[].compress`（int）：是否开启压缩。
- `regionParams[].storagePoolParams[].compressLocation`（int）：压缩位置。
- `regionParams[].storagePoolParams[].compressThreadNum`（int）：压缩线程数。
- `regionParams[].storagePoolParams[].compressAlgorithm`（int）：压缩算法。
- `regionParams[].storagePoolParams[].deduplication`（int）：是否开启重删。
- `regionParams[].storagePoolParams[].deduplicationLocation`（int）：重删位置。
- `regionParams[].storagePoolParams[].deduplicationThreadNum`（int）：重删线程数。
- `regionParams[].storagePoolParams[].fingerLibraryBindType`（int）：指纹库绑定类型。
- `regionParams[].storagePoolParams[].fingerLibraryId`（string）：指纹库 ID。
- `regionParams[].storagePoolParams[].fingerLibraryConfig`（object）：指纹库配置对象。
- `regionParams[].storagePoolParams[].dataConsistencySwitch`（bool）：数据一致性开关。
- `regionParams[].storagePoolParams[].dataConsistencyLogic`（int）：数据一致性算法。

### `commonConfigParams.conventionalConfig`（列表展开）

- `permanentIncrement`（int）：是否开启永久增量。
- `failureRetry`（int）：是否开启失败重试。
- `failureRetryCount`（int）：失败重试次数。
- `failureRetryInterval`（int）：失败重试间隔（分钟）。
- `encryptionTrans`（int）：是否开启数据传输加密。
- `forcedRetentionSwitch`（bool）：强制保留开关。
- `forcedRetentionCycle`（int）：强制保留周期（天）。

### 最小响应示例字段说明

- `backupType`：备份类型（决定使用 xtrabackup 或 LVM 快照）。
- `backupGran`：备份粒度（实例级或库级）。
- `isMergeBackup`：是否启用合成备份。
- `isMultiChannel`：是否启用多通道备份。
- `dataChannelNum`：数据通道数量。
- `commonConfigParams.appType`：应用类型标识。
- `commonConfigParams.backupDestination.multipleRegions`：是否启用多区域备份。
- `commonConfigParams.backupDestination.regionParams[].storagePoolParams[].storagePoolId`：目标存储池 ID。
- `commonConfigParams.backupDestination.regionParams[].storagePoolParams[].storageServiceId`：目标存储服务 ID。
- `commonConfigParams.conventionalConfig.failureRetry`：失败重试开关。
- `commonConfigParams.conventionalConfig.failureRetryCount`：失败重试次数。
- `commonConfigParams.conventionalConfig.encryptionTrans`：传输加密开关。

## 返回案例

```json
{
    "status": "success",
    "error": null,
    "responseData": {
        "realPath": "",
        "backupType": 2,
        "backupGran": 1,
        "isMergeBackup": true,
        "isRealTimeLog": false,
        "isBackupMode": false,
        "backupMode": 3,
        "isNodeSequence": false,
        "nodeSequence": "",
        "isTableRestore": false,
        "isArchDelete": true,
        "deleteArchUnit": 5,
        "deleteArchPeriod": 1,
        "isMultiChannel": true,
        "dataChannelNum": 8,
        "isNotAutoToFull": false,
        "lvmChunkSliceSize": 128,
        "isForwardSelect": false,
        "isBackupCompress": false,
        "isRateControl": false,
        "rateSize": 0,
        "gatewayId": "",
        "gatewayName": "",
        "gatewayIP": "",
        "gatewayType": 0,
        "linkType": 0,
        "nodeWwn": "",
        "clientsWwn": null,
        "lanFreeWwn": null,
        "volumeInfo": "",
        "hasCustomParam": false,
        "customParam": "",
        "usePageTrack": false,
        "cleanPageTrack": false,
        "commonConfigParams": {
            "appType": "",
            "overwriteConfig": true,
            "backupDestination": {
                "multipleRegions": false,
                "regionParams": [
                    {
                        "regionName": "",
                        "clientParams": null,
                        "storageAvailable": false,
                        "triggerCondition": {
                            "taskFailedCount": 0
                        },
                        "storagePoolParams": [
                            {
                                "isMasterStorage": true,
                                "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
                                "storageServiceName": "StorageService_172.31.12.91",
                                "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
                                "storagePoolName": "block",
                                "storagePoolType": 2,
                                "encryptionStorage": 0,
                                "encryptionLocation": 0,
                                "encryptionThreadNum": 0,
                                "encryptAlgo": 0,
                                "compress": 0,
                                "compressLocation": 0,
                                "compressThreadNum": 0,
                                "compressAlgorithm": 0,
                                "deduplication": 0,
                                "deduplicationLocation": 0,
                                "deduplicationThreadNum": 0,
                                "fingerLibraryBindType": 0,
                                "fingerLibraryId": "",
                                "fingerLibraryConfig": null,
                                "dataConsistencySwitch": false,
                                "dataConsistencyLogic": 1
                            }
                        ]
                    }
                ]
            },
            "conventionalConfig": {
                "permanentIncrement": 0,
                "failureRetry": 0,
                "failureRetryCount": 0,
                "failureRetryInterval": 0,
                "encryptionTrans": 0,
                "forcedRetentionSwitch": false,
                "forcedRetentionCycle": 0
            }
        }
    }
}
```

## 示例
```bash
foundation-cli mysql backup-config detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --system-id <system-id> \
  --object-id <object-id>
```



