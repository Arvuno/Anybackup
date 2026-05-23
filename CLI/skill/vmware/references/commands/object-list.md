# `foundation-cli vmware object list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli vmware object list --production-system-id <id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/{productionSystemId}/get_vms` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 必填 | 说明 |
|---|---|---|---|
| tenantId | `--tenant-id` | 否 | 租户标识 |
| endpoint | `--endpoint` | 是 | Foundation 服务地址 |
| ak | `--ak` | 是 | Access Key |
| sk | `--sk` | 是 | Secret Key |
| targetVersion | `--target-version` | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| productionSystemId | `--production-system-id` | string | 是 | 生产系统 ID |

## 返回参数

CLI 外层结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

`responseData`（`VMwareVmListResponse`）：

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 总数 |
| data | object[] | VMware 保护对象列表 |

`responseData.data[]`（`ProtectObject`）核心字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| productionSystemId | string | 生产系统 ID |
| objectId | string | 保护对象资源 ID |
| vmId | string | VM ID |
| name | string | VM 名称 |
| hasPassword | bool | 是否存在密码 |
| execStatus | object[] | 执行状态列表 |
| bindStrategyStatus | int | 策略绑定状态 |
| protectStatus | int | 保护状态 |
| backupSpeed | string | 备份速度 |
| recoverySpeed | string | 恢复速度 |
| lastRunTime | int64 | 上次备份时间 |
| lastBackupDataSize | int64 | 上次备份数据量 |
| lastBackupResult | int | 上次备份结果 |
| lastBackupTaskId | string | 上次备份作业 ID |
| latestTimePoint | string | 最新时间点 |
| hasBackupConfig | bool | 是否存在备份配置 |
| hasBackupData | bool | 是否已有备份数据 |
| disableReason | string | 不支持备份原因 |
| backupStatus | int | 备份状态 |
| backupStrategyId | string | 备份策略 ID |
| backupStrategyName | string | 备份策略名称 |
| backupStrategyNum | int | 额外策略数量 |
| isGroupSLA | bool | 是否组合策略 |
| nextBackupTime | int64 | 下次备份时间 |
| backupConfig | object | 备份参数配置 |

`execStatus[]` 字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| latestJobId | string | 最新作业 ID |
| status | int | 执行状态 |

## 枚举说明

### `bindStrategyStatus`

| 值 | 含义 |
|---|---|
| 1 | 已绑定 |
| 2 | 未绑定 |

### `protectStatus`

| 值 | 含义 |
|---|---|
| 1 | 支持备份 |
| 2 | 不支持备份 |

### `execStatus[].status`

| 值 | 含义 |
|---|---|
| 1 | 备份中 |
| 2 | 恢复中 |
| 3 | 清理中 |
| 4 | 复制中 |
| 5 | 浏览中 |

### `lastBackupResult`

| 值 | 含义 |
|---|---|
| 600 | 失败 |
| 700 | 已取消 |
| 800 | 成功 |
| 900 | 成功有告警 |
| 1100 | 已终止 |

### `backupStatus`

| 值 | 含义 |
|---|---|
| 1 | 不支持备份 |
| 2 | 未配置保护策略 |
| 3 | 未产生备份 |
| 4 | 已有备份 |
| 5 | 已配置备份策略 |

### `backupConfig` 相关枚举

| 字段 | 值 | 含义 |
|---|---|---|
| encryptionTrans | 0/1 | 关闭/开启传输加密 |
| encryptionStorage | 0/1 | 关闭/开启存储加密 |
| encryptionLocation | 1/2 | 源端加密/目标端加密 |
| encryptionAlgorithm | 0/1/2/3 | 不加密/AES256GCM/SM4/BlowFish |
| compress | 0/1 | 关闭/开启压缩 |
| compressLocation | 1/2 | 源端压缩/目标端压缩 |
| compressAlgorithm | 0/1/2/3 | 不压缩/快速/标准/强力 |
| deduplication | 0/1 | 关闭/开启重删 |
| deduplicationLocation | 1/2 | 源端重删/目标端重删 |
| failureRetry | 0/1 | 关闭/开启失败重试 |

## 示例

```bash
foundation-cli vmware object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --production-system-id <production-system-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ object\ list |
| 风险 | read-only |

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
