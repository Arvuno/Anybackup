# `foundation-cli mysql backup-detail`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql backup-detail --tenant-id <tenant-id> --task-id <task-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/get_backup_task_detail?taskId=<task-id>` |
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
| taskId | `--task-id` | string | 是 | MySQL 备份任务 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | MySQL 备份任务详情 |

### `responseData` 字段（完整）

| 字段 | 类型 | 说明 |
|---|---|---|
| multipleBackupTaskDetail | object | 多区域备份任务详情 |
| strategyId | string | 策略 ID |
| strategyName | string | 策略名称 |
| isStrategyExecute | int32 | 是否由策略触发 |
| clientInfo | object[] | 关联客户端记录 |
| clusterMode | int32 | 集群模式 |
| realPath | string | 数据源路径 |
| backupType | int32 | 备份类型 |
| backupGran | int32 | 备份粒度 |
| isMergeBackup | bool | 是否启用合成备份 |
| isRealTimeLog | bool | 是否启用实时日志 |
| isBackupMode | bool | 是否开启备份模式 |
| backupMode | int32 | 备份模式 |
| isNodeSequence | bool | 是否节点顺序备份 |
| nodeSequence | string | 节点备份顺序 |
| isTableRestore | bool | 是否支持表级恢复 |
| isArchDelete | bool | 是否开启日志删除策略 |
| deleteArchUnit | int32 | 日志删除策略单位 |
| deleteArchPeriod | int32 | 日志删除策略数值 |
| isMultiChannel | bool | 是否启用多通道备份 |
| dataChannelNum | int32 | 数据通道数量 |
| isNotAutoToFull | bool | 是否关闭自动转全量 |
| lvmChunkSliceSize | int32 | LVM 切片大小 |
| isForwardSelect | bool | 是否启用正向选择 |
| isBackupCompress | bool | 是否启用高级压缩 |
| isRateControl | bool | 是否启用流量控制 |
| rateSize | int32 | 流量控制值 |
| gatewayId | string | SAN 网关 ID |
| gatewayName | string | SAN 网关名称 |
| gatewayIP | string | SAN 网关 IP |
| gatewayType | int | SAN 网关类型 |
| linkType | int32 | SAN 链路类型 |
| nodeWwn | string | SAN 节点 WWN/IQN |
| clientsWwn | string[]/null | SAN 客户端 WWN/IQN 列表 |
| lanFreeWwn | object[]/null | LAN-Free WWN 映射 |
| volumeInfo | string | 卷信息 |
| hasCustomParam | bool | 是否启用自定义参数 |
| customParam | string | 自定义参数 |
| usePageTrack | bool | 是否使用 PageTrack |
| cleanPageTrack | bool | 是否清理 PageTrack |

### `multipleBackupTaskDetail` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| taskId | string | 任务 ID |
| objectId | string | 保护对象 ID |
| productionSystemId | string | 生产系统 ID |
| appType | string | 应用类型 |
| appParam | map[string]string | 应用参数 |
| storagePoolConfigs | object[] | 存储池配置列表 |
| clientConfig | object[] | 客户端配置列表 |
| gateWayInfos | object[] | 网关配置列表 |
| failureRetry | bool | 是否开启重试 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔 |
| encryptionTrans | bool | 是否启用传输加密 |
| forcedRetentionSwitch | bool | 强制保留开关 |
| forcedRetentionCycle | int32 | 强制保留周期 |

### `storagePoolConfigs[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| regionId | string | 区域 ID |
| regionName | string | 区域名称 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int32 | 存储池类型 |
| encryption | bool | 是否启用存储加密 |
| encryptionLocation | int32 | 加密位置 |
| encryptionTreadNum | int32 | 加密线程数 |
| encryptionAlgorithm | int32 | 加密算法 |
| compress | bool | 是否启用压缩 |
| compressLocation | int32 | 压缩位置 |
| compressTreadNum | int32 | 压缩线程数 |
| compressAlgorithm | int32 | 压缩算法 |
| deduplication | bool | 是否启用重删 |
| deduplicationLocation | int32 | 重删位置 |
| deduplicationTreadNum | int32 | 重删线程数 |
| fingerLibraryId | string | 指纹库 ID |
| fingerLibraryConfig | object/null | 指纹库配置 |
| dataConsistencySwitch | bool | 数据一致性开关 |
| dataConsistencyLogic | int32 | 数据一致性算法 |

### `fingerLibraryConfig` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| name | string | 指纹库名称 |
| sliceMinValue | int32 | 切片最小值 |
| sliceMaxValue | int32 | 切片最大值 |
| dataSize | int32 | 预估数据量 |
| compressAlgorithm | int32 | 压缩算法 |
| creatorName | string | 创建者名称 |

### `clientConfig[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| regionId | string | 区域 ID |
| regionName | string | 区域名称 |
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientOsType | int | 客户端操作系统类型 |
| clientIp | string | 客户端 IP |
| clientStatus | int | 客户端状态 |

### `gateWayInfos[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| regionId | string | 区域 ID |
| regionName | string | 区域名称 |
| gatewayInfo | object/null | 网关配置对象 |

### `gatewayInfo` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| gateWayId | string | 网关 ID |
| gatewayName | string | 网关名称 |
| gateWayIp | string | 网关 IP |
| shareName | string | 共享目录名称 |
| whiteList | string[] | 白名单 |
| linkType | int32 | 链路类型 |
| nodeWwn | string | 节点 WWN/IQN |
| clientsWwn | string[] | 客户端 WWN/IQN 列表 |
| lanFreeWwn | object[] | LAN-Free WWN 映射 |

### `clientInfo[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientMac | string | 客户端机器码 |
| clientOSType | int64 | 客户端操作系统类型 |
| clientType | int32 | 客户端类型 |
| clientStatus | int64 | 客户端状态 |

### `lanFreeWwn[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| nodeId | string | 节点 ID |
| serverWwn | string | 服务端 WWN |
| clientWwn | string | 客户端 WWN |
| clientId | string | 客户端 ID |

## 返回案例

```json
{
    "status": "success",
    "error": null,
    "responseData": {
        "realPath": "",
        "backupType": 1,
        "backupGran": 1,
        "isMergeBackup": false,
        "isRealTimeLog": false,
        "isBackupMode": false,
        "backupMode": 3,
        "isNodeSequence": false,
        "nodeSequence": "",
        "isTableRestore": false,
        "isArchDelete": false,
        "deleteArchUnit": 0,
        "deleteArchPeriod": 0,
        "isMultiChannel": false,
        "dataChannelNum": 0,
        "isNotAutoToFull": true,
        "lvmChunkSliceSize": 0,
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
        "multipleBackupTaskDetail": {
            "storagePoolConfigs": [
                {
                    "regionId": "5d2821cd9efc23d2269354c68395ad6d",
                    "regionName": "",
                    "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
                    "storageServiceName": "StorageService_172.31.12.91",
                    "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
                    "storagePoolName": "block",
                    "storagePoolType": 2,
                    "encryption": false,
                    "encryptionLocation": 0,
                    "encryptionTreadNum": 0,
                    "encryptionAlgorithm": 0,
                    "compress": false,
                    "compressLocation": 0,
                    "compressTreadNum": 0,
                    "compressAlgorithm": 0,
                    "deduplication": false,
                    "deduplicationLocation": 0,
                    "deduplicationTreadNum": 0,
                    "fingerLibraryId": "",
                    "fingerLibraryConfig": null,
                    "dataConsistencySwitch": false,
                    "dataConsistencyLogic": 0
                }
            ],
            "clientConfig": [
                {
                    "regionId": "5d2821cd9efc23d2269354c68395ad6d",
                    "regionName": "",
                    "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
                    "clientName": "iv-yej4xkshkwwh2yram2x1",
                    "clientOsType": 2,
                    "clientIp": "172.31.0.2",
                    "clientStatus": 1
                }
            ],
            "failureRetry": false,
            "failureRetryCount": 0,
            "failureRetryInterval": 0,
            "encryptionTrans": false,
            "forcedRetentionSwitch": false,
            "forcedRetentionCycle": 0,
            "taskId": "f134a7e65114116ab50400163e33df73",
            "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
            "productionSystemId": "a63ba04bfb530fefa3b09d1f89d47432",
            "appType": "eso_backupengine_hypermysqlengine",
            "appParam": {
                "EEE_CUSTOM_PARAMETERS_IS_OPEN": "0",
                "EEE_FORWARD_SELECT": "0",
                "EEE_IS_DELETE_ARCHIVE_LOG": "0",
                "EEE_MYSQL_AUTOLOG_SWITCH": "0",
                "EEE_MYSQL_BACKUP_COMPRESS": "0",
                "EEE_MYSQL_BACKUP_GRAN": "1",
                "EEE_MYSQL_BACKUP_MODE": "3",
                "EEE_MYSQL_BACKUP_MODE_IS_OPEN": "0",
                "EEE_MYSQL_BACKUP_TAG": "2026_04_10_14_38_34",
                "EEE_MYSQL_BACKUP_TYPE": "1",
                "EEE_MYSQL_CLEAN_PAGE_TRACK": "0",
                "EEE_MYSQL_DATABASE_RATE": "0",
                "EEE_MYSQL_MUTIL_CHANNEL": "0",
                "EEE_MYSQL_NODE_SEQUENCE": "",
                "EEE_MYSQL_NODE_SEQUENCE_IS_OPEN": "0",
                "EEE_MYSQL_NOT_TO_FULL": "1",
                "EEE_MYSQL_OS_USER": "mysql",
                "EEE_MYSQL_REALPATH": "",
                "EEE_MYSQL_USE_PAGE_TRACK": "0",
                "EEE_OBJECT_ID": "428e4c05bf87deef20a86e670f3eb9a9",
                "EEE_SYSTEM_ID": "a63ba04bfb530fefa3b09d1f89d47432"
            }
        },
        "isStrategyExecute": 0,
        "strategyId": "",
        "strategyName": "",
        "clientInfo": [
            {
                "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
                "clientName": "iv-yej4xkshkwwh2yram2x1",
                "clientIp": "172.31.0.2",
                "clientMac": "U0HYTDM3RENXS4EJ",
                "clientOSType": 2,
                "clientType": 1,
                "clientStatus": 0
            }
        ],
        "clusterMode": 3
    }
}
```

## 示例
```bash
foundation-cli mysql backup-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --task-id <task-id>
```



