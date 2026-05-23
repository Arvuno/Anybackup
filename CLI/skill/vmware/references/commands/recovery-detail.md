# `foundation-cli vmware recovery-detail`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli vmware recovery-detail --tenant-id <tenant-id> --task-id <task-id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/recovery_task/{taskId}/detail` |
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
| taskId | `--task-id` | string | Yes | VMware recovery task ID |

## Response Parameters

### Root Fields

| Field | Type | Description |
|---|---|---|
| recoveryObject | object | 恢复源对象信息 |
| recoveryDestination | object | 恢复目标信息 |
| vmConfig | object | 恢复后的虚拟机配置 |
| recoveryOptions | object | 恢复执行选项 |

### `recoveryObject`

| Field | Type | Description |
|---|---|---|
| waitRecoveryVmName | string | 待恢复虚拟机名称 |
| waitRecoveryVmId | string | 待恢复虚拟机 ID |
| storagePoolType | int32 | 存储池类型 |
| storagePoolId | string | 存储池 ID |
| storagePoolName | string | 存储池名称 |
| storageServiceId | string | 存储服务 ID |
| storageServiceName | string | 存储服务名称 |
| recoveryTimeStamp | int64 | 恢复时间戳 |
| recoveryDisplayTime | string | 恢复时间展示值 |
| recoveryTpType | int32 | 恢复时间点类型 |

### `recoveryDestination`

| Field | Type | Description |
|---|---|---|
| platformType | int32 | 目标平台类型 |
| platformName | string | 目标平台名称 |
| platformIp | string | 目标平台 IP |
| vSphereType | string | vSphere 平台类型 |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientId | string | 客户端 ID |
| recoveryMode | int | 恢复方式 |
| configMode | int | 配置方式 |
| overlayMode | int | 覆盖方式 |

### `vmConfig`

| Field | Type | Description |
|---|---|---|
| vmName | string | 原虚拟机名称 |
| vmId | string | 原虚拟机 ID |
| recoveryVMName | string | 恢复后的虚拟机名称 |
| location | string | 恢复位置 |
| resource | string | 资源池 |
| disksConfig | object[] | 磁盘配置 |
| vifsConfig | object[] | 网卡配置 |
| autoStartVm | int32 | 恢复后是否自动开机 |
| useOriginalMac | int32 | 是否使用原 MAC 地址 |
| nameSuffix | string | 新名称后缀 |
| enableCustomCpuAndMem | bool | 是否启用自定义 CPU/内存 |
| memorySize | int | 内存大小，GiB |
| numCoresPerSocket | int | 每插槽核心数 |
| cpuSlots | int | CPU 插槽数 |

### `recoveryOptions`

| Field | Type | Description |
|---|---|---|
| transportMode | string | 传输模式 |
| enableAutoRepair | int32 | 异构恢复自动修复开关 |
| failureRetry | int32 | 失败重试开关 |
| failureRetryCount | int32 | 重试次数 |
| failureRetryInterval | int32 | 重试间隔 |
| encryptionTrans | int32 | 传输加密开关 |

## Enum Values

### `recoveryObject.recoveryTpType`

| Value | Meaning |
|---|---|
| 1 | 完全备份 |
| 2 | 增量备份 |
| 3 | 差异备份 |
| 4 | 日志备份 |
| 5 | 永久增量 |

### `recoveryDestination.recoveryMode`

| Value | Meaning |
|---|---|
| 1 | 新建恢复 |
| 2 | 覆盖恢复 |
| 3 | 异构恢复 |

### `recoveryDestination.configMode`

| Value | Meaning |
|---|---|
| 1 | 原配置 |
| 2 | 指定配置 |

### `recoveryDestination.overlayMode`

| Value | Meaning |
|---|---|
| 1 | 全量恢复 |

### `vmConfig.autoStartVm` / `vmConfig.useOriginalMac` / `recoveryOptions.enableAutoRepair`

| Value | Meaning |
|---|---|
| 0 | 否 |
| 1 | 是 |

### `recoveryOptions.transportMode`

| Value | Meaning |
|---|---|
| `nbd` | NBD |
| `nbdssl` | NBDSSL |
| `hotadd` | HotAdd |
| `san` | SAN |
| `auto` | 自动 |

## Verification

- 已按 [recovery_detail.go](/e:/go/src/foundation-cli/internal/business/vmware/recovery_detail.go) 与 [restore_config_create.go](/e:/go/src/foundation-cli/internal/business/vmware/restore_config_create.go) 校验返回字段和枚举。
- 当前命令参数 `--task-id` 与设计稿、实现代码一致。

## Examples

```bash
foundation-cli vmware recovery-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --task-id <task-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ recovery-detail |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli vmware recovery-detail
```
