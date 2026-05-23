# `foundation-cli mysql backup-config set`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql backup-config set --tenant-id <tenant-id> --data '<json>'` |
| 方法 | `PUT` |
| 路径 | `/backupmgm/v1/mysql/app_backup_config` |
| 风险 | `写入` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 地址 |
| ak | `--ak` | string | 是 | 访问密钥 AK |
| sk | `--sk` | string | 是 | 访问密钥 SK |
| targetVersion | `--target-version` | string | 否 | 默认 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| data | `--data`, `-d` | string | 二选一 | 行内 JSON 请求体 |
| bodyFile | `--body-file` | string | 二选一 | JSON 请求体文件 |

## Body 参数

### 顶层参数（完整字段）

| Field | Type | Required | Description |
|---|---|---|---|
| objectId | string | yes | 保护对象 ID |
| systemId | string | yes | 生产系统 ID |
| backupType | int | yes | 备份类型 |
| backupGran | int | yes | 备份粒度 |
| commonConfigParams | object | yes | 公共配置对象 |
| backupMedia | object/null | no | 备份介质配置 |
| isBackupMode | bool | no | 是否启用备份模式 |
| backupMode | int | no | 备份模式 |
| syntheticBackup | bool | no | 是否开启合成备份，映射 `isMergeBackup` |
| syntheticBackupValue | int | no | 合成备份附加值 |
| isMergeBackup | bool | no | 是否合成备份 |
| isNodeSequence | bool | no | 是否节点顺序备份 |
| isNotAutoToFull | bool | no | 是否关闭自动转全备 |
| usePageTrack | bool | no | 是否启用 PageTrack |
| cleanPageTrack | bool | no | 是否清理 PageTrack |
| isRealTimeLog | bool | no | 是否启用实时日志 |
| isTableLevelRecovery | bool | no | 兼容别名，映射 `isTableRestore` |
| isTableRestore | bool | no | 是否支持表级恢复 |
| isArchDelete | bool | no | 是否日志删除策略选项 |
| deleteArchPeriod | int | no | 日志删除策略数值 |
| deleteArchUnit | int | no | 日志删除策略单位 |
| dataChannelNum | int | no | 数据通道数 |
| isDataChannelNum | bool | no | 前端控制位；不是 `isMultiChannel` 的别名 |
| isMultiChannel | bool | no | 是否启用多通道 |
| lvmChunkSliceSize | int | no | LVM 切片大小 |
| nonAutomaticTransferComplete | bool | no | 前端控制位；不是 `isNotAutoToFull` 的别名 |
| hasCustomParam | bool | no | 是否启用自定义参数 |
| customParam | string | no | 自定义参数内容 |
| compress | int | no | 顶层压缩开关（兼容） |
| deduplication | int | no | 顶层重删开关（兼容） |
| failureRetry | bool | no | 顶层失败重试开关；与 `conventionalConfig.failureRetry` 分开传递 |
| failureRetryCount | int | no | 顶层重试次数；与 `conventionalConfig.failureRetryCount` 分开传递 |
| failureRetryInterval | int | no | 顶层重试间隔；与 `conventionalConfig.failureRetryInterval` 分开传递 |
| forcedRetentionSwitch | bool | no | 顶层强制保留开关；与 `conventionalConfig.forcedRetentionSwitch` 分开传递 |

### `commonConfigParams` 参数

| Field | Type | Required | Description |
|---|---|---|---|
| commonConfigParams.appType | string | yes | 应用类型 |
| commonConfigParams.overwriteConfig | bool | no | 是否覆盖原配置 |
| commonConfigParams.backupDestination | object | yes | 备份目标配置 |
| commonConfigParams.conventionalConfig | object | no | 常规配置 |

### `commonConfigParams.backupDestination` 参数

| Field | Type | Required | Description |
|---|---|---|---|
| ...backupDestination.multipleRegions | bool | no | 是否多区域 |
| ...backupDestination.regionParams | object[] | yes | 区域列表（至少 1 项） |
| ...regionParams[].storageAvailable | bool | no | 存储可用性开关 |
| ...regionParams[].storagePoolParams | object[] | yes | 存储池列表（至少 1 项） |

### `storagePoolParams[]` 参数（关键）

| Field | Type | Required | Description |
|---|---|---|---|
| ...storagePoolParams[].storagePoolId | string | yes | 存储池 ID |
| ...storagePoolParams[].storagePoolName | string | yes | 存储池名称 |
| ...storagePoolParams[].storagePoolType | int | yes | 存储池类型 |
| ...storagePoolParams[].storageServiceId | string | yes | 存储服务 ID |
| ...storagePoolParams[].storageServiceName | string | no | 存储服务名称 |
| ...storagePoolParams[].isMasterStorage | bool | no | 是否主存储 |
| ...storagePoolParams[].encryptionStorage | int | no | 存储加密开关 |
| ...storagePoolParams[].compress | int | no | 压缩开关 |
| ...storagePoolParams[].deduplication | int | no | 重删开关 |
| ...storagePoolParams[].deduplicationLocation | int | no | 重删位置 |
| ...storagePoolParams[].deduplicationThreadNum | int | no | 重删线程数 |
| ...storagePoolParams[].fingerLibraryBindType | int | no | 指纹库绑定类型 |
| ...storagePoolParams[].fingerLibraryId | string | no | 指纹库 ID |
| ...storagePoolParams[].fingerLibraryConfig | object | no | 指纹库配置 |
| ...storagePoolParams[].dataConsistencyLogic | int | no | 数据一致性策略 |

### `fingerLibraryConfig` 参数

| Field | Type | Required | Description |
|---|---|---|---|
| ...fingerLibraryConfig.name | string | no | 指纹库名称 |
| ...fingerLibraryConfig.sliceMinValue | int | no | 最小切片值 |
| ...fingerLibraryConfig.sliceMaxValue | int | no | 最大切片值 |
| ...fingerLibraryConfig.dataSize | int | no | 数据量阈值 |
| ...fingerLibraryConfig.compressAlgorithm | int | no | 压缩算法 |
| ...fingerLibraryConfig.creatorName | string | no | 创建者 |

### 存储池自动选择约束

- 存储池信息必须来自 `foundation-cli protect storage-pool auto-select` 返回值。
- 不要手工猜测 `storageServiceId` / `storagePoolId`。
- `storagePoolId` 必须属于同一次 `auto-select` 返回的 `storageSvcId`。
- 映射关系：
- `storageSvcId -> storageServiceId`
- `storageSvcName -> storageServiceName`
- `storagePoolId -> storagePoolId`
- `storagePoolName -> storagePoolName`
- `storagePoolType -> storagePoolType`

### 自动补全行为（CLI）

- 当 `commonConfigParams.backupDestination.regionParams[0].storagePoolParams[0]` 中未选择存储池服务/存储池（例如缺少 `storageServiceId` 或 `storagePoolId`）时，CLI 会自动调用：
  `foundation-cli protect storage-pool auto-select`
- 自动补全时，CLI 仅填充缺失字段，不覆盖你已显式填写的值。
- 自动补全来源映射：
  - `responseData.storageSvcId -> storageServiceId`
  - `responseData.storageSvcName -> storageServiceName`
  - `responseData.storagePoolId -> storagePoolId`
  - `responseData.storagePoolName -> storagePoolName`
  - `responseData.storagePoolType -> storagePoolType`

### 兼容别名说明（CLI 自动归一）

- `syntheticBackup -> isMergeBackup`
- `isTableLevelRecovery -> isTableRestore`
- `isDataChannelNum` 和 `isMultiChannel` 不是别名；成功请求里允许 `isDataChannelNum=false` 且 `isMultiChannel=true`
- 当 `dataChannelNum > 0` 且缺少 `isMultiChannel` 时，CLI 会自动补 `isMultiChannel=true`
- `nonAutomaticTransferComplete` 和 `isNotAutoToFull` 不是别名；建议两个字段都显式传，常见成功组合为 `nonAutomaticTransferComplete=true`、`isNotAutoToFull=false`
- `conventionalConfig.failureRetry*`、`conventionalConfig.forcedRetentionSwitch` 不会覆盖同名顶层字段；若两层都传，CLI 保留用户显式值
- `conventionalConfig.encryptionTrans`、`conventionalConfig.failureRetry` 视为必填；未传时 CLI 自动补 `0`

### 成功参数约束（推荐直接遵守）

- 所有布尔字段必须传 `true/false`，不能传 `0/1`；例如 `isRealTimeLog` 必须写成 `false`，不能写成 `0`
- `backupMode` 优先使用 `3`（优先从节点），不要默认传 `5`
- `storagePoolParams[0]` 建议显式带上：`encryptionStorage=0`、`compress=0`、`deduplication=0`、`dataConsistencyLogic=1`、`isMasterStorage=true`
- 若 `dataChannelNum > 0`，请求体中应出现 `isMultiChannel=true`；`isDataChannelNum` 可按前端控制位单独传值，不要拿它替代 `isMultiChannel`
- 顶层 `failureRetry=false` / `failureRetryCount=1` / `failureRetryInterval=1` / `forcedRetentionSwitch=false` 与 `commonConfigParams.conventionalConfig` 中的同名字段是两层独立语义，建议都显式传
- 推荐显式带上这些稳定默认值：`backupMedia=null`、`isBackupMode=false`、`syntheticBackupValue=1`、`usePageTrack=false`、`cleanPageTrack=false`、`isArchDelete=false`、`deleteArchPeriod=1`、`deleteArchUnit=4`、`lvmChunkSliceSize=128`、`nonAutomaticTransferComplete=true`、`compress=0`、`deduplication=0`

常见错误：

- 错误：`"isRealTimeLog": 0`
- 正确：`"isRealTimeLog": false`

### `HyperBackupMgm.ForcedRetentionSwitchError` 防错规则（必须遵守）

- 现网场景下，`forcedRetentionSwitch` 不能随意传 `false`，否则容易触发 `HyperBackupMgm.ForcedRetentionSwitchError`。
- 为保证每次都正确，MySQL 备份配置统一使用以下安全模板：
- `commonConfigParams.conventionalConfig.forcedRetentionSwitch = true`
- `commonConfigParams.conventionalConfig.forcedRetentionCycle = 1`
- 若同时传顶层字段，则保持一致：`forcedRetentionSwitch = true`
- 推荐先执行 `foundation-cli mysql backup-config detail` 读取当前值，再显式带上以上字段，避免后端默认值差异导致失败。

## 枚举列表

- `backupType`：`1=xtrabackup`，`2=LVM快照`
- `backupGran`：`1=实例`，`2=库`
- `backupMode`：`1=仅从备节点`，`2=仅从主节点`，`3=优先从节点`，`4=优先主节点`，`5=任意节点`
- `deleteArchUnit`：`1=年`，`2=月`，`3=周`，`4=天`，`5=小时`
- `storagePoolType`：`1=SAN`，`2=本地磁盘`，`3=分布式`，`4=NAS`，`5=磁带`，`7=本地文件系统`，`8=对象存储`

## 请求体示例

```json
{
  "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
  "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
  "commonConfigParams": {
    "appType": "eso_backupengine_hypermysqlengine",
    "backupDestination": {
      "multipleRegions": false,
      "regionParams": [
        {
          "storageAvailable": false,
          "storagePoolParams": [
            {
              "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
              "storagePoolName": "block",
              "storagePoolType": 2,
              "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
              "storageServiceName": "StorageService_172.31.12.91",
              "encryptionStorage": 0,
              "compress": 0,
              "deduplication": 1,
              "deduplicationLocation": 1,
              "deduplicationThreadNum": 4,
              "fingerLibraryBindType": 1,
              "fingerLibraryId": "f13c864382cb11de8d2600163e33df73",
              "fingerLibraryConfig": {
                "compressAlgorithm": 2,
                "creatorName": "chp",
                "dataSize": 30,
                "name": "mysql3306_U0HYTDM3RENXS4EJ_1776668217",
                "sliceMaxValue": 256,
                "sliceMinValue": 128
              },
              "dataConsistencyLogic": 1,
              "isMasterStorage": true
            }
          ]
        }
      ]
    },
    "conventionalConfig": {
      "encryptionTrans": 1,
      "failureRetry": 1,
      "failureRetryCount": 2,
      "failureRetryInterval": 5,
      "forcedRetentionSwitch": true,
      "forcedRetentionCycle": 1
    }
  },
  "backupMedia": null,
  "isBackupMode": false,
  "backupMode": 3,
  "backupType": 1,
  "syntheticBackup": false,
  "syntheticBackupValue": 1,
  "isMergeBackup": false,
  "isNodeSequence": false,
  "isNotAutoToFull": false,
  "usePageTrack": false,
  "cleanPageTrack": false,
  "isRealTimeLog": false,
  "isTableLevelRecovery": false,
  "isArchDelete": false,
  "deleteArchPeriod": 1,
  "deleteArchUnit": 4,
  "dataChannelNum": 2,
  "isDataChannelNum": false,
  "isMultiChannel": true,
  "lvmChunkSliceSize": 128,
  "nonAutomaticTransferComplete": true,
  "hasCustomParam": false,
  "backupGran": 1,
  "compress": 0,
  "deduplication": 0,
  "failureRetry": false,
  "failureRetryCount": 1,
  "failureRetryInterval": 1,
  "forcedRetentionSwitch": false
}
```

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": null
}
```

## 示例

```bash
foundation-cli mysql backup-config set \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --body-file backup-config.json
```
