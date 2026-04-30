# `foundation-cli vmware object info`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli vmware object info --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/virtual/vmware/{objectId}/info` |
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
| objectId | `--object-id` | string | Yes | VMware protect object ID |

## Response Parameters

### Top-level Envelope

| Field | Type | Description |
|---|---|---|
| status | string | 接口状态 |
| error | object/null | 错误信息 |
| responseData | object | VMware 保护对象详情 |

### `responseData`

| Field | Type | Description |
|---|---|---|
| protectObjectId | string | 保护对象 ID |
| vmId | string | 虚拟机 ID |
| name | string | 虚拟机名称 |
| type | string | 保护对象类型 |
| hostModePath | string | 主机模式路径 |
| templateModePath | string | 模板模式路径 |
| pool | string | 资源池 |
| guestOS | string | 客户机操作系统 |
| ip | string | IP 地址 |
| cpu | int | CPU 核数 |
| ram | int64 | 内存大小，单位字节 |
| powerState | string | 电源状态 |
| disk | object[] | 磁盘列表 |
| nic | object[] | 网卡列表 |

### `disk[]`

| Field | Type | Description |
|---|---|---|
| name | string | 磁盘名称 |
| size | int64 | 磁盘大小，单位字节 |
| storage | string | 所属存储 |

### `nic[]`

| Field | Type | Description |
|---|---|---|
| name | string | 网卡名称 |
| networkname | string | 网络名称 |
| macaddr | string | MAC 地址 |

## Enum Values

### `powerState`

| Value | Meaning |
|---|---|
| `poweredOff` | 关机 |
| `poweredOn` | 开机 |
| `suspended` | 挂起 |

## Verification

- 已按 [object_info.go](/e:/go/src/foundation-cli/internal/business/vmware/object_info.go) 校验返回结构。
- 当前命令参数 `--object-id` 与设计稿、实现代码一致。

## Examples

```bash
foundation-cli vmware object info \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ vmware\ object\ info |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli vmware object info
```
