# `foundation-cli mysql object list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql object list --tenant-id <tenant-id> --app-type 202 [filters]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/database/object_list` |
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
| isAscSort | `--is-asc-sort` | bool | 否 | 升序排序时显式传 `true` 或 `false` |
| sortField | `--sort-field` | int | 否 | 允许值： `1/2/3/4/5/6` |
| index | `--index` | int | 否 | 默认值为 `0` |
| count | `--count` | int | 否 | 默认值为 `20`，取值范围 `1~100` |
| filter | `--filter` | string | 否 | 模糊筛选关键字，最大长度 `256` |
| includedPath | `--included-path` | bool | 否 | 是否启用路径匹配 |
| productionSystemIds | `--production-system-ids` | string[] | 否 | 重复 flag，生产系统 ID 列表 |
| objectIds | `--object-ids` | string[] | 否 | 重复 flag，保护对象 ID 列表 |
| objectTypes | `--object-types` | string[] | 否 | 重复 flag |
| appTypes | `--app-types` | string[] | 否 | 重复 flag |
| groupId | `--group-id` | string | 否 | 主机组 ID |
| hostId | `--host-id` | string | 否 | 主机 ID |
| emptyHostId | `--empty-host-id` | bool | 否 | 过滤未关联主机的对象 |
| hostTag | `--host-tag` | string | 否 | 主机标签 |
| isConfig | `--is-config` | int[] | 否 | 允许值： `1/2/3/4` |
| objectStatus | `--object-status` | int | 否 | 允许值： `1/2` |
| objectMode | `--object-mode` | int | 否 | 允许值： `1/2` |
| canBackup | `--can-backup` | int | 否 | 允许值： `1/2` |
| bindSlaStatus | `--bind-sla-status` | int[] | 否 | 允许值： `1~10` |
| includedBindSla | `--included-bind-sla` | int[] | 否 | 允许值： `1/2/3/4/5` |
| excludeBindSla | `--exclude-bind-sla` | int[] | 否 | 允许值： `1/2/3/4/5` |
| bindSlaIds | `--bind-sla-ids` | string[] | 否 | 重复 flag，最多 `10` 个 ID |
| protectStatus | `--protect-status` | int[] | 否 | 允许值： `1/2/3/4` |
| lastBackupStatus | `--last-backup-status` | int[] | 否 | 允许值： `600/700/800/900/1000/1100` |
| lastSnapshotStatus | `--last-snapshot-status` | int[] | 否 | 允许值： `600/700/800/900/1000/1100` |
| execStatus | `--exec-status` | int | 否 | 允许值： `1/2/3/4` |
| parentIds | `--parent-ids` | string[] | 否 | 重复 flag |
| allChild | `--all-child` | bool | 否 | 展开所有子节点 |
| nameFilter | `--name-filter` | string | 否 | 名称或位置筛选，最大长度 `1024` |
| appType | `--app-type` | int32 | 否 | 应用类型编码，MySQL 联调建议固定使用 `202` |
| isolation | `--isolation` | bool | 否 | 仅筛选隔离对象 |
| objectId | `--object-id` | string | 否 | 单个对象 ID |

## 响应参数

### 顶层响应字段

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态，成功通常为 `success` |
| error | object/null | 错误对象，成功时通常为 `null` |
| responseData | object | 业务响应体 |

### `responseData` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 对象总数 |
| data | object[] | MySQL 保护对象列表 |

### `data[]` 字段（完整）

| 字段 | 类型 | 说明 |
|---|---|---|
| systemId | string | 生产系统 ID |
| objectId | string | 保护对象 ID |
| name | string | 对象名称 |
| businessName | string | 业务名称 |
| platName | string | 平台名称 |
| clusterMode | int32 | 集群模式 |
| protectStatus | int32 | 保护状态 |
| execStatus | int32[] | 执行状态编码列表 |
| execStatuses | object[] | 当前执行状态列表 |
| bindStrategyStatus | int32 | SLA 绑定状态 |
| backupStrategyId | string | 备份策略 ID |
| backupStrategyName | string | 备份策略名称 |
| backupStrategyNum | int | 备份策略数量 |
| isGroupSLA | bool | 是否组 SLA |
| type | int32 | 实例类型编码 |
| latestTimePoint | string | 最近时间点 |
| lastBackupDataSize | int64 | 最近一次备份数据量 |
| lastRunTime | int64 | 最近一次备份执行时间（毫秒） |
| lastBackupResult | int32 | 最近备份结果 |
| lastBackupTaskId | string | 最近备份任务 ID |
| hasBackupConfig | bool | 是否存在备份配置 |
| hasAppBackupConfig | bool | 是否存在应用备份配置 |
| hasBackupData | bool | 是否存在备份数据 |
| objectStatus | int32 | 对象状态 |
| disableReason | string | 不可备份原因 |
| genMode | int32 | 资源来源模式 |
| isDeleted | bool | 是否已删除 |
| backupMode | int32 | 备份模式 |
| nextBackupTime | int64 | 下次备份时间 |
| backupConfig | object/null | 备份配置对象 |
| snapshotConfig | object/null | 快照配置对象 |
| fingerLibraryId | string | 指纹库 ID |
| fingerPrintLibrary | object/null | 指纹库信息 |
| backupSpeed | string | 备份速度 |
| recoverySpeed | string | 恢复速度 |
| lastSnapshotTime | int64 | 最近快照时间 |
| lastSnapshotSize | int64 | 最近快照数据量 |
| lastSnapshotStatus | int32 | 最近快照状态 |
| lastSnapshotTaskId | string | 最近快照任务 ID |
| lastSuccessSyncTime | int64 | 最近成功同步时间 |
| tenantId | string | 租户 ID |
| tenantName | string | 租户名称 |
| subObjectNumber | int32 | 子对象数量 |
| resourcestatus | int32 | 资源状态 |
| instanceInfo | string | 实例信息 |
| clientInfo | object[]/null | 客户端信息列表 |
| appExtend | object[]/null | 应用扩展字段列表 |
| creator | string | 创建者 |
| hasPassword | bool | 是否存在恢复校验密码 |
| backupDestination | object/null | 备份目的地配置 |
| conventionalConfig | object/null | 常规配置 |

### `execStatuses[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| status | int32 | 执行状态（1:备份中，2:恢复中，3:清理中，4:远程复制中，5:归档中） |
| latestTaskId | string | 最新作业 ID |

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

### `appExtend[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| key | string | 扩展字段键 |
| value | string | 扩展字段值 |

### `backupDestination` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| multipleRegions | bool | 是否启用多区域 |
| regionParams | object[] | 区域配置列表 |

### `regionParams[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| regionName | string | 区域名称 |
| clientParams | object[]/null | 客户端参数列表 |
| storageAvailable | bool | 是否启用存储可用性能力 |
| triggerCondition | object/null | 触发条件 |
| storagePoolParams | object[] | 存储池参数列表 |

### `clientParams[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientOsType | int | 客户端操作系统类型 |
| clientIp | string | 客户端 IP |
| clientStatus | int | 客户端状态 |

### `triggerCondition` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| taskFailedCount | int32 | 任务失败次数阈值 |

### `storagePoolParams[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| isMasterStorage | bool | 是否主存储池 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int32 | 存储池类型 |
| encryptionStorage | int32 | 存储加密开关 |
| encryptionLocation | int32 | 加密位置 |
| encryptionThreadNum | int32 | 加密线程数 |
| encryptAlgo | int32 | 加密算法 |
| compress | int32 | 压缩开关 |
| compressLocation | int32 | 压缩位置 |
| compressThreadNum | int32 | 压缩线程数 |
| compressAlgorithm | int32 | 压缩算法 |
| deduplication | int32 | 重删开关 |
| deduplicationLocation | int32 | 重删位置 |
| deduplicationThreadNum | int32 | 重删线程数 |
| fingerLibraryBindType | int32 | 指纹库绑定类型 |
| fingerLibraryId | string | 指纹库 ID |
| fingerLibraryConfig | object/null | 指纹库配置 |
| dataConsistencySwitch | bool | 数据一致性开关 |
| dataConsistencyLogic | int32 | 数据一致性算法 |

### `fingerPrintLibrary` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| fingerLibraryId | string | 指纹库 ID |
| name | string | 指纹库名称 |
| sliceMinValue | int32 | 切片最小值 |
| sliceMaxValue | int32 | 切片最大值 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| dataSize | int32 | 预估数据量 |
| compressAlgorithm | int32 | 压缩算法 |
| objectUsedCount | int32 | 已绑定对象数量 |
| creatorName | string | 创建者名称 |

### `backupConfig` / `snapshotConfig` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | 配置 ID |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storagePoolType | int32 | 存储池类型 |
| clientId | string | 客户端 ID |
| status | int32 | 配置状态 |
| permanentIncrement | int32 | 永久增量开关 |
| failureRetry | int32 | 失败重试开关 |
| failureRetryCount | int32 | 失败重试次数 |
| failureRetryInterval | int32 | 失败重试间隔（分钟） |
| configType | int32 | 配置类型 |
| encryptionTrans | int32 | 传输加密开关 |
| encryptionStorage | int32 | 存储加密开关 |
| encryptionLocation | int32 | 加密位置 |
| encryptionTreadNum | int32 | 加密线程数 |
| encryptionAlgorithm | int32 | 加密算法 |
| compress | int32 | 压缩开关 |
| compressLocation | int32 | 压缩位置 |
| compressTreadNum | int32 | 压缩线程数 |
| compressAlgorithm | int32 | 压缩算法 |
| deduplication | int32 | 重删开关 |
| deduplicationLocation | int32 | 重删位置 |
| deduplicationTreadNum | int32 | 重删线程数 |
| fingerLibraryId | string | 指纹库 ID |
| forcedRetentionCycle | int32 | 强制保留周期（天） |
| forcedRetentionSwitch | bool | 强制保留开关 |
| dataConsistencyLogic | int32 | 数据一致性算法 |
| dataConsistencySwitch | bool | 数据一致性开关 |
| backupProtocolType | int32 | 备份协议类型 |
| gateWayInfo | object/null | 网关配置 |
| storageServiceId | string | 存储服务 ID |
| encryptionSecretKey | string | 加密密钥 |
| lanFree | bool | lan-free 开关 |
| version | int64 | 配置版本 |
| createTime | string | 创建时间 |
| updateTime | string | 更新时间 |
| configVersion | int32 | 配置版本号 |
| multipleRegions | bool | 是否多区域 |

### `gateWayInfo` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| gateWayId | string | 网关 ID |
| gatewayName | string | 网关名称 |
| gateWayIp | string | 网关 IP |
| shareName | string | 共享目录名称 |
| whiteList | string[] | 白名单 |
| linkType | int32 | SAN 链路类型 |
| nodeWwn | string | 节点 WWN/IQN |
| clientsWwn | string[] | 客户端 WWN/IQN 列表 |
| lanFreeWwn | object[] | lan-free WWN 配置 |

### `lanFreeWwn[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| nodeId | string | 节点 ID |
| serverWwn | string | 服务端 WWN |
| clientWwn | string | 客户端 WWN |
| clientId | string | 客户端 ID |

### `conventionalConfig` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| permanentIncrement | int32 | 是否开启永久增量 |
| failureRetry | int32 | 是否开启失败重试 |
| failureRetryCount | int32 | 失败重试次数 |
| failureRetryInterval | int32 | 失败重试间隔（分钟） |
| encryptionTrans | int32 | 传输加密开关 |
| forcedRetentionSwitch | bool | 强制保留开关 |
| forcedRetentionCycle | int32 | 强制保留周期（天） |

## 枚举值

### `clusterMode`

| 值 | 含义 |
|---|---|
| 1 | 非分布式 |
| 2 | 分布式 |

### `protectStatus`

| 值 | 含义 |
|---|---|
| 1 | 未保护 |
| 2 | 全保护 |
| 3 | 部分保护 |
| 4 | 已删除 |

### `bindStrategyStatus`

| 值 | 含义 |
|---|---|
| 1 | 已绑定 |
| 2 | 未绑定 |

### `objectStatus`

| 值 | 含义 |
|---|---|
| 1 | 不具备备份能力 |
| 2 | 无备份配置 |
| 3 | 尚未备份 |
| 4 | 已完成备份 |

### `genMode`

| 值 | 含义 |
|---|---|
| 0 | 本地注册 |
| 1 | 数据同步 |
| 2 | 反向复制 |

## 返回案例

```json
{
    "status": "success",
    "error": null,
    "responseData": {
        "data": [
            {
                "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
                "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
                "name": "mysql3306_U0HYTDM3RENXS4EJ",
                "businessName": "mysql3306_U0HYTDM3RENXS4EJ",
                "platName": "iv-yej4xkshkwwh2yram2x1(172.31.0.2)",
                "clusterMode": 3,
                "protectStatus": 2,
                "execStatus": null,
                "execStatuses": null,
                "bindStrategyStatus": 0,
                "backupStrategyId": "",
                "backupStrategyName": "",
                "backupStrategyNum": 0,
                "isGroupSLA": false,
                "type": 202,
                "latestTimePoint": "2026-04-10T14:38:44.000000+08:00",
                "lastBackupDataSize": 19901940,
                "lastRunTime": 1775803115000,
                "lastBackupResult": 800,
                "hasBackupConfig": true,
                "hasAppBackupConfig": true,
                "hasBackupData": true,
                "objectStatus": 4,
                "disableReason": "",
                "genMode": 0,
                "isDeleted": false,
                "backupMode": 0,
                "nextBackupTime": -1,
                "backupConfig": null,
                "snapshotConfig": null,
                "fingerLibraryId": "",
                "fingerPrintLibrary": null,
                "backupSpeed": "",
                "recoverySpeed": "",
                "lastSnapshotTime": 0,
                "lastSnapshotSize": 0,
                "lastSnapshotStatus": 0,
                "lastSnapshotTaskId": "",
                "lastSuccessSyncTime": 1775885977280,
                "tenantId": "",
                "tenantName": "",
                "subObjectNumber": 0,
                "resourcestatus": 1,
                "instanceInfo": "",
                "clientInfo": null,
                "appExtend": null,
                "creator": "mysql",
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
        ],
        "totalNum": 1
    }
}
```

## 示例
### 最小调用

```bash
foundation-cli mysql object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --app-type 202
```

### 按主机、分组和关键字过滤

```bash
foundation-cli mysql object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --app-type 202 \
  --group-id default \
  --host-id 12345678901234567890123456789012 \
  --filter prod \
  --included-path true \
  --count 30
```



