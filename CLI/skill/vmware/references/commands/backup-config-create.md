# `foundation-cli vmware backup-config create`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli vmware backup-config create --object-id <id> --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/virtual/vmware/{objectId}/backup_config` |
| 风险 | `write` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 `FOUNDATION_TENANT_ID` 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 服务地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | 保护对象 ID |
| body | `--data` | JSON string | 是 | 备份配置请求体 |

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| objectIds | string[] | 是 | 保护对象资源 ID 列表，建议与 `--object-id` 保持一致 |
| overwriteConfig | bool | 否 | 是否覆盖已有配置 |
| clientId | string | 是 | 客户端 ID |
| backupOnce | bool | 否 | 配置完成后是否立即备份一次 |
| storagePoolType | int | 是 | 存储池类型 |
| storagePoolId | string | 是 | 存储池 ID |
| storagePoolName | string | 否 | 存储池名称 |
| storageServiceId | string | 是 | 存储服务 ID |
| storageServiceName | string | 否 | 存储服务名称 |
| failureRetry | int | 否 | 失败重试开关 |
| failureRetryCount | int | 否 | 失败重试次数 |
| failureRetryInterval | int | 否 | 失败重试间隔 |
| encryptionTrans | int | 否 | 传输加密开关 |
| deduplication | int | 否 | 去重开关 |
| deduplicationLocation | int | 否 | 去重位置 |
| deduplicationThreadNum | int | 否 | 去重线程数（`0~32`） |
| fingerLibraryId | string | 否 | 指纹库 ID |
| fingerLibraryName | string | 否 | 指纹库名称 |
| fingerPrintLibrary | object | 否 | 指纹库对象 |
| forcedRetentionCycle | int | 否 | 强制保留周期（`0~9999`） |
| forcedRetentionSwitch | bool | 否 | 强制保留周期开关 |
| snapshotType | int | 否 | 快照类型 |
| changeBlockTracking | bool | 否 | 是否启用 CBT |
| transportMode | string | 是 | 传输模式 |
| usetoNBD | bool | 否 | 是否切换到 NBD |
| diskConcurrentThreads | int | 是 | 磁盘并发线程数（`1~8`） |
| concurrentThreadsPerDisk | int | 是 | 单磁盘并发线程数（`1~8`） |

## 枚举列表

- `snapshotType`：`0` 普通快照，`1` 静默快照
- `transportMode`：`nbd`、`nbdssl`、`hotadd`、`san`、`auto`
- `failureRetry`：`0` 关闭，`1` 开启
- `encryptionTrans`：`0` 关闭，`1` 开启
- `deduplication`：`0` 关闭，`1` 开启
- `deduplicationLocation`：`1` 源端去重，`2` 目标端去重

## 请求体示例

```json
{
  "objectIds": [
    "<object-id>"
  ],
  "overwriteConfig": true,
  "clientId": "<client-id>",
  "storagePoolType": 1,
  "storagePoolId": "<storage-pool-id>",
  "storageServiceId": "<storage-service-id>",
  "transportMode": "auto",
  "diskConcurrentThreads": 2,
  "concurrentThreadsPerDisk": 2
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
foundation-cli vmware backup-config create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --data '{"objectIds":["<object-id>"],"overwriteConfig":true,"clientId":"<client-id>","storagePoolType":1,"storagePoolId":"<storage-pool-id>","storageServiceId":"<storage-service-id>","transportMode":"auto","diskConcurrentThreads":2,"concurrentThreadsPerDisk":2}'
```
