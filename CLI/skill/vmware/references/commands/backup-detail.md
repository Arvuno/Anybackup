# `foundation-cli vmware backup-detail`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli vmware backup-detail --tenant-id <tenant-id> --task-id <task-id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/backup_task/{taskId}/detail` |
| Risk | `read-only` |

## Request Parameters

### Shared Parameters

| Field | CLI Flag | Type | Required | Notes |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | No | Tenant identifier; can come from FOUNDATION_TENANT_ID env var |
| endpoint | `--endpoint` | string | Yes | Foundation base URL |
| ak | `--ak` | string | Yes | Access key |
| sk | `--sk` | string | Yes | Secret key |
| targetVersion | `--target-version` | string | No | Defaults to `9.0.9.0` |

### Command Parameters

| Field | CLI Flag | Type | Required | Notes |
|---|---|---|---|---|
| taskId | `--task-id` | string | Yes | VMware backup task ID |

## Response Parameters

### Root Fields

| Field | Type | Description |
|---|---|---|
| backupConfig | object | 备份任务配置详情 |

### `backupConfig`

| Field | Type | Description |
|---|---|---|
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientId | string | 客户端 ID |
| storagePoolType | int32 | 存储池类型 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| failureRetry | int32 | 失败重试开关 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔 |
| snapshotType | int | 快照类型 |
| changeBlockTracking | bool | CBT 开关 |
| transportMode | string | 传输模式 |
| usetoNBD | bool | 是否转 NBD |
| diskConcurrentThreads | int | 磁盘并发线程数 |
| concurrentThreadsPerDisk | int | 单磁盘并发线程数 |
| encryptionTrans | int32 | 传输加密开关 |
| deduplication | int32 | 重删开关 |
| deduplicationLocation | int32 | 重删位置 |
| deduplicationThreadNum | int32 | 重删线程数 |
| fingerLibraryId | string | 指纹库 ID |
| forcedRetentionCycle | int32 | 强制保留周期 |
| forcedRetentionSwitch | bool | 强制保留开关 |

## Enum Values

### `failureRetry` / `encryptionTrans` / `deduplication`

| Value | Meaning |
|---|---|
| 0 | 关闭 |
| 1 | 开启 |

### `deduplicationLocation`

| Value | Meaning |
|---|---|
| 1 | 源端重删 |
| 2 | 目标端重删 |

### `snapshotType`

| Value | Meaning |
|---|---|
| 0 | 普通快照 |
| 1 | 静默快照 |

### `transportMode`

| Value | Meaning |
|---|---|
| `nbd` | NBD |
| `nbdssl` | NBDSSL |
| `hotadd` | HotAdd |
| `san` | SAN |
| `auto` | 自动 |

## Verification

- 已按 [backup_detail.go](/e:/go/src/foundation-cli/internal/business/vmware/backup_detail.go)、[backup_config_detail.go](/e:/go/src/foundation-cli/internal/business/vmware/backup_config_detail.go) 与 [object_list.go](/e:/go/src/foundation-cli/internal/business/vmware/object_list.go) 校验返回字段。
- 当前命令参数 `--task-id` 与设计稿、实现代码一致。

## Examples

```bash
foundation-cli vmware backup-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --task-id <task-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ backup-detail |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli vmware backup-detail
```
