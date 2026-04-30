# `foundation-cli vmware timepoint metadata`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli vmware timepoint metadata --tenant-id <tenant-id> --timestamp <ts> --data-set-id <id> [--business 1] [--request-id <id>] [--storage-service-id <id>]` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/time_point/get_metadata` |
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
| timestamp | `--timestamp` | string | Yes | Timepoint timestamp |
| dataSetId | `--data-set-id` | string | Yes | Dataset ID |
| business | `--business` | string | No | Only `1` is currently accepted by CLI validation |
| requestId | `--request-id` | string | No | Optional request correlation ID |
| storageServiceId | `--storage-service-id` | string | No | Storage service ID |

## Response Parameters

### Top-level Envelope

| Field | Type | Description |
|---|---|---|
| status | string | 接口状态 |
| error | object/null | 错误信息 |
| responseData | object | 时间点元数据响应 |

### `responseData`

| Field | Type | Description |
|---|---|---|
| tpInfo | object | 时间点信息 |
| finished | bool | 是否已完成 |
| requestId | string | 请求 ID |

### `tpInfo`

| Field | Type | Description |
|---|---|---|
| dataSetId | string | 数据集 ID |
| timestamp | int64 | 时间戳 |
| display | int64 | 前端展示值 |
| gnsPath | string | 路径 |
| displayTime | string | 展示时间 |
| crossPlatRecovery | object | 异构恢复能力 |
| metaData | object | 时间点元数据 |

### `crossPlatRecovery`

| Field | Type | Description |
|---|---|---|
| disableReason | string | 不支持异构恢复原因 |
| supportCrossPlatRecovery | bool | 是否支持异构恢复 |
| unsupportedAppTypes | string[] | 不支持的应用类型 |
| unsupportedReasons | string[] | 不支持原因列表 |

### `metaData`

| Field | Type | Description |
|---|---|---|
| version | int | 元数据版本 |
| platformType | int | 平台类型 |
| platform | object | 平台信息 |
| cpuSlots | int | CPU 插槽数 |
| cpuCores | int | 每 CPU 核心数 |
| memorySize | int64 | 内存大小，字节 |
| osSystem | string | 操作系统名称 |
| osBit | int | 系统位数 |
| bootType | string | 启动类型 |
| vmTools | bool | 是否安装 VMware Tools |
| arch | string | 架构 |
| vmId | string | 虚拟机 ID |
| vmName | string | 虚拟机名称 |
| vmFullName | string | 虚拟机完整路径 |
| totalSize | int64 | 总占用空间 |
| disk | object[] | 磁盘列表 |
| nic | object[] | 网卡列表 |
| computeResource | object | 计算资源 |
| location | object | 位置信息 |
| vmObjectId | string | 虚拟机对象 ID |

## Enum Values

### `platformType`

| Value | Meaning |
|---|---|
| 1 | XenServer |
| 2 | InCloud Sphere |
| 3 | VMware |
| 4 | FusionCompute |
| 5 | RHEV |
| 6 | Hyper-V |
| 7 | H3CLOUD |
| 8 | CAS |
| 11 | SANGFOR HCI |
| 12 | ZStack |
| 13 | OracleVM |
| 14 | IROS |
| 15 | InCloudSphereKvm |
| 16 | SmartX |
| 17 | Nutanix |
| 18 | QingCloud |
| 19 | OpenStack |
| 20 | HCS |
| 21 | AliyunPrivate |
| 22 | Bingo |
| 30 | FusionOne |
| 31 | Proxmox VE |

## Verification

- 已按 [timepoint_metadata.go](/e:/go/src/foundation-cli/internal/business/vmware/timepoint_metadata.go) 与 [restore_config_create.go](/e:/go/src/foundation-cli/internal/business/vmware/restore_config_create.go) 校验返回字段。
- 当前命令参数与设计稿一致，其中 `--business` 仍按实现约束仅允许 `1`。

## Examples

```bash
foundation-cli vmware timepoint metadata \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --timestamp <timestamp> \
  --data-set-id <data-set-id> \
  --business 1
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ timepoint\ metadata |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli vmware timepoint metadata
```
