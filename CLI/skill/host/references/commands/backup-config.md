# `foundation-cli host backup-config`

## 命令概要

| Item | Value |
|---|---|
| CLI | `foundation-cli host backup-config --tenant-id <tenant-id> --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/batch/backup_tasks/start` |
| Risk | `write` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation console base URL |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| data | `--data` / `-d` | string | 是 | inline JSON body |
| bodyFile | `--body-file` | string | 否 | JSON body file path（与 `--data` 二选一，由 CLI 读取后作为 body 发送） |

## Body 字段

### 顶层字段

- `filesetsInfo`
  - 类型：`array`
  - 必填：是
  - 说明：保护对象资源信息列表
  - 约束：数组长度 `1~100`
- `overwriteConfig`
  - 类型：`boolean`
  - 必填：否
  - 说明：是否覆盖配置
- `storageServiceId`
  - 类型：`string`
  - 必填：是
  - 说明：存储服务 ID
  - 约束：`min=1, max=48`
- `backupOnce`
  - 类型：`boolean`
  - 必填：否
  - 说明：配置完成后是否立即备份一次（`true/false`）

- `storagePoolType`
  - 类型：`int`
  - 必填：是
  - 说明：存储池类型
- `storagePoolId`
  - 类型：`string`
  - 必填：是
  - 说明：存储池 ID
- `storagePoolName`
  - 类型：`string`
  - 必填：否
  - 说明：存储池名称

- `encryptionTrans`
  - 类型：`int32`
  - 必填：否
  - 说明：是否开启数据传输加密
  - 约束：允许 `0/1`（默认 `0`）
- `encryptionStorage`
  - 类型：`int32`
  - 必填：否
  - 说明：是否开启数据存储加密
  - 约束：允许 `0/1`（默认 `0`）
- `encryptionLocation`
  - 类型：`int32`
  - 必填：否
  - 说明：加密位置/方式
  - 约束：允许 `1/2`（默认 `1`）
- `encryptionThreadNum`
  - 类型：`int32`
  - 必填：否
  - 说明：加密线程数
- `encryptAlgo`
  - 类型：`int32`
  - 必填：否
  - 说明：数据加密算法
  - 约束：允许 `0/1/2/3`（默认 `1`）

- `compress`
  - 类型：`int32`
  - 必填：否
  - 说明：是否开启数据压缩
  - 约束：允许 `0/1`
- `compressLocation`
  - 类型：`int32`
  - 必填：否
  - 说明：压缩位置/方式
  - 约束：允许 `1/2`
- `compressThreadNum`
  - 类型：`int32`
  - 必填：否
  - 说明：压缩线程数
  - 约束：允许 `-1/0/1/2/3/4/5/6/7/8`
- `compressAlgorithm`
  - 类型：`int32`
  - 必填：否
  - 说明：压缩算法
  - 约束：允许 `0/1/2/3`

- `deduplication`
  - 类型：`int32`
  - 必填：否
  - 说明：是否开启重删
  - 约束：允许 `0/1`
- `deduplicationLocation`
  - 类型：`int32`
  - 必填：否
  - 说明：重删类型
  - 约束：允许 `0/1/2`
- `deduplicationThreadNum`
  - 类型：`int32`
  - 必填：否
  - 说明：重删线程数
  - 约束：`min=-1, max=32`
- `fingerLibraryBindType`
  - 类型：`int32`
  - 必填：否
  - 说明：指纹库绑定类型
  - 约束：允许 `0/1/2/3`
- `fingerLibraryId`
  - 类型：`string`
  - 必填：否
  - 说明：指纹库唯一标识
  - 约束：`min=32, max=48`
- `fingerPrintLibrary`
  - 类型：`object`
  - 必填：否
  - 说明：指纹库配置对象（字段以服务端实际支持为准）

- `failureRetry`
  - 类型：`int32`
  - 必填：否
  - 说明：是否重试
  - 约束：允许 `0/1`（默认 `0`）
- `failureRetryCount`
  - 类型：`int32`
  - 必填：否
  - 说明：重试次数
- `failureRetryInterval`
  - 类型：`int32`
  - 必填：否
  - 说明：重试间隔

- `forcedRetentionCycle`
  - 类型：`int32`
  - 必填：否
  - 说明：强制数据保留周期（保留天数）
  - 约束：`-1~9999`
- `forcedRetentionSwitch`
  - 类型：`boolean`
  - 必填：否
  - 说明：强制数据保留开关（`true/false`）

- `dataConsistencyLogic`
  - 类型：`int32`
  - 必填：否
  - 说明：数据一致性算法参数
  - 约束：允许 `1/3/4/5`
- `dataConsistencySwitch`
  - 类型：`boolean`
  - 必填：否
  - 说明：数据一致性开关（`true/false`）

- `customScript`
  - 类型：`int`
  - 必填：否
  - 说明：自定义脚本开关
  - 约束：允许 `0/1`（默认 `0`）
- `preScriptPath`
  - 类型：`string`
  - 必填：否
  - 说明：前置脚本路径
  - 约束：`min=1, max=4096`
- `failureScriptPath`
  - 类型：`string`
  - 必填：否
  - 说明：失败脚本路径
  - 约束：`min=1, max=4096`
- `successScriptPath`
  - 类型：`string`
  - 必填：否
  - 说明：成功脚本路径
  - 约束：`min=1, max=4096`
- `backupChannel`
  - 类型：`int`
  - 必填：否
  - 说明：是否多通道备份
  - 约束：允许 `0/1`（默认 `0`）
- `backupChannelReadThread`
  - 类型：`int`
  - 必填：否
  - 说明：多通道备份读取并发数
  - 约束：`1~50`（默认 `1`）
- `backupChannelTraversalThread`
  - 类型：`int`
  - 必填：否
  - 说明：多通道备份遍历并发数
  - 约束：`1~50`（默认 `1`）
- `backupSize`
  - 类型：`int`
  - 必填：否
  - 说明：增量备份粒度
  - 约束：`1~9999`（默认 `1`）
- `resumableDataTransfer`
  - 类型：`int`
  - 必填：否
  - 说明：断点续备开关
  - 约束：允许 `0/1`（默认 `0`）

- `lanFree`
  - 类型：`int32`
  - 必填：否
  - 说明：lanFree 开关
  - 约束：允许 `0/1`
- `gateWayId`
  - 类型：`string`
  - 必填：否
  - 说明：网关 ID（NAS 或 SAN 备份时需要）
  - 约束：`min=3, max=128`
- `gateWayName`
  - 类型：`string`
  - 必填：否
  - 说明：网关名称（NAS 或 SAN 备份时需要）
  - 约束：`min=3, max=128`
- `gateWayIp`
  - 类型：`string`
  - 必填：否
  - 说明：网关 IP（NAS 备份时需要）
  - 约束：`ip, max=254`
- `linkType`
  - 类型：`int32`
  - 必填：否
  - 说明：链路类型（SAN 备份时需要）
  - 约束：允许 `0/1/2/3`
- `lanFreeWwns`
  - 类型：`array`
  - 必填：否
  - 说明：lanFree 关联参数

### `filesetsInfo[]` 字段

- `filesetId`
  - 类型：`string`
  - 必填：是
  - 说明：保护对象资源 ID
  - 约束：长度 `len=32`
- `hostId`
  - 类型：`string`
  - 必填：是
  - 说明：生产系统 ID（主机 ID）
  - 约束：长度 `len=32`

### `lanFreeWwns[]` 字段

- `lanFreeWwn`
  - 类型：`array`
  - 必填：否
  - 说明：lanFree WWN 列表（开启 lanFree 时使用）
- `hostInfo`
  - 类型：`object`
  - 必填：否
  - 说明：主机信息对象（字段以服务端实际支持为准）

### `lanFreeWwns[].lanFreeWwn[]` 字段

- `nodeId`
  - 类型：`string`
  - 必填：否
  - 说明：节点 ID
- `serverWwn`
  - 类型：`string`
  - 必填：否
  - 说明：服务端 WWN
- `clientWwn`
  - 类型：`string`
  - 必填：否
  - 说明：客户端 WWN
- `clientId`
  - 类型：`string`
  - 必填：否
  - 说明：客户端 ID

## 返回结果

该命令直接透传后端返回的 JSON（不做结构化改造）。如果返回结果中包含 `jobId`，可以使用 `foundation-cli job backup-detail --tenant-id <tenant-id> --job-id <jobId>` 或 `foundation-cli job sync-detail --tenant-id <tenant-id> --job-id <jobId>` 查询任务详情。

## 示例

最小命令示例：

```bash
foundation-cli host backup-config --tenant-id <tenant-id> --data '<json>'
```

最小 JSON body 示例：

```json
{
  "filesetsInfo": [
    {
      "filesetId": "<fileset-id>",
      "hostId": "<host-id>"
    }
  ],
  "storageServiceId": "1",
  "storagePoolType": 1,
  "storagePoolId": "<storage-pool-id>"
}
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ backup-config |
| 风险 | write |

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| data | object | 是 | 请求体对象，建议按“请求体示例”中的字段组织 |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
