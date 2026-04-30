# `foundation-cli host backup-config detail`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli host backup-config detail --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/file/backup_config/{objectId}` |
| Risk | `read-only` |

## Request Parameters

### Shared Parameters

| Field | CLI Flag | Type | Required | Notes |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | No | Tenant identifier; can come from FOUNDATION_TENANT_ID env var |
| endpoint | `--endpoint` | string | Yes | Foundation console base URL |
| ak | `--ak` | string | Yes | Access key |
| sk | `--sk` | string | Yes | Secret key |
| targetVersion | `--target-version` | string | No | Defaults to `9.0.9.0` |

### Command Parameters

| Field | CLI Flag | Type | Required | Notes |
|---|---|---|---|---|
| objectId | `--object-id` | string | Yes | Host protect object ID |

## Response Parameters

### Root Fields

| Field | Type | Description |
|---|---|---|
| storageServiceId | string | Storage service ID |
| storageServiceName | string | Storage service name |
| storagePoolType | int | Storage pool type |
| storagePoolId | string | Storage pool ID |
| storagePoolName | string | Storage pool name |
| encryptionTrans | int32 | Transport encryption switch |
| encryptionStorage | int32 | Storage encryption switch |
| encryptionLocation | int32 | Encryption location |
| encryptionThreadNum | int32 | Encryption thread count |
| encryptAlgo | int32 | Encryption algorithm |
| compress | int32 | Compression switch |
| compressLocation | int32 | Compression location |
| compressThreadNum | int32 | Compression thread count |
| compressAlgorithm | int32 | Compression algorithm |
| deduplication | int32 | Deduplication switch |
| deduplicationLocation | int32 | Deduplication location |
| deduplicationThreadNum | int32 | Deduplication thread count |
| fingerLibraryBindType | int32 | Fingerprint library bind type |
| fingerLibraryId | string | Fingerprint library ID |
| fingerPrintLibrary | object | Fingerprint library configuration |
| failureRetry | int32 | Failure retry switch |
| failureRetryCount | int32 | Retry count |
| failureRetryInterval | int32 | Retry interval |
| forcedRetentionCycle | int32 | Forced retention cycle |
| forcedRetentionSwitch | bool | Forced retention switch |
| lanFree | int32 | LAN-Free switch |
| dataConsistencySwitch | bool | Data consistency switch |
| dataConsistencyLogic | int32 | Data consistency mode |
| customScript | int | Custom script switch |
| preScriptPath | string | Pre-backup script path |
| failureScriptPath | string | Failure script path |
| successScriptPath | string | Success script path |
| backupChannel | int | Multi-channel backup switch |
| backupChannelReadThread | int | Read thread count |
| backupChannelTraversalThread | int | Traversal thread count |
| backupSize | int | Incremental backup granularity |
| resumableDataTransfer | int | Resume transfer switch |
| gateWayId | string | Gateway ID |
| gateWayName | string | Gateway name |
| gateWayIp | string | Gateway IP |
| linkType | int32 | Link type |
| lanFreeWwns | object[] | LAN-Free WWN list |

### `fingerPrintLibrary`

| Field | Type | Description |
|---|---|---|
| name | string | Fingerprint library name |
| sliceMinValue | int32 | Minimum slice size |
| sliceMaxValue | int32 | Maximum slice size |
| storagePoolId | string | Fingerprint library storage pool ID |
| storagePoolName | string | Fingerprint library storage pool name |
| dataSize | int32 | Estimated data size |
| compressAlgorithm | int32 | Compression algorithm |
| creatorName | string | Creator name |

## Enum Values

### `encryptionTrans` / `encryptionStorage` / `compress` / `deduplication` / `failureRetry`

| Value | Meaning |
|---|---|
| 0 | Disabled |
| 1 | Enabled |

### `encryptionLocation` / `compressLocation`

| Value | Meaning |
|---|---|
| 1 | Source side |
| 2 | Target side |

### `encryptAlgo`

| Value | Meaning |
|---|---|
| 0 | No encryption |
| 1 | AES256GCM |
| 2 | SM4 |
| 3 | BlowFish |

### `compressAlgorithm`

| Value | Meaning |
|---|---|
| 0 | No compression |
| 1 | Fast |
| 2 | Standard |
| 3 | Strong |

### `fingerLibraryBindType`

| Value | Meaning |
|---|---|
| 0 | Unspecified |
| 1 | Existing |
| 2 | Auto-configured |
| 3 | Newly created |

### `lanFree`

| Value | Meaning |
|---|---|
| 0 | Disabled |
| 1 | Enabled |

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ backup-config\ detail |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli host backup-config detail
```
