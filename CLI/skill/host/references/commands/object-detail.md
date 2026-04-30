# `foundation-cli host object detail`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli host object detail --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v1/file/fileset/{objectId}` |
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
| objectId | `--object-id` | string | Yes | Host protect object ID |

## Response Parameters

### Top-level Response

该接口直接返回主机保护对象详情对象本身，不额外包裹 `status/responseData`。

### Root Fields

| Field | Type | Description |
|---|---|---|
| name | string | 文件集名称 |
| mode | int | 数据源模式 |
| type | int | 数据源类型 |
| pathType | int | 路径类型 |
| ndmpService | int | NDMP 服务方式 |
| intelligentArchive | int | 是否启用智能目录归档 |
| startDir | string | 起始目录 |
| datasourcePaths | object[] | 数据源路径列表 |
| includeFilterRule | object[] | 包含过滤规则 |
| excludeFilterRule | object[] | 排除过滤规则 |
| openTimingFilterRule | bool | 是否启用定时过滤规则 |
| timingFilterRule | object | 定时过滤规则 |
| isDeleted | bool | 保护对象是否已删除 |
| incrementTraverseRange | object | 增量遍历范围 |
| protectGranularity | int32 | 保护粒度 |

### `datasourcePaths[]`

| Field | Type | Description |
|---|---|---|
| fullPath | string | 数据源完整路径 |
| nodeType | int | 节点类型 |
| type | string | 数据源类别，如 `disk`、`volume` |

### `includeFilterRule[]` / `excludeFilterRule[]`

| Field | Type | Description |
|---|---|---|
| path | object | 路径过滤条件 |
| format | object | 格式过滤条件 |
| time | object | 时间过滤条件 |

### `timingFilterRule`

| Field | Type | Description |
|---|---|---|
| file | object | 文件路径过滤 |
| directory | object | 目录路径过滤 |
| format | object | 文件格式过滤 |
| time | object | 时间范围过滤 |

### `incrementTraverseRange`

| Field | Type | Description |
|---|---|---|
| open | bool | 是否启用增量遍历范围 |
| value | int | 范围值 |
| unit | int | 范围单位 |

## Enum Values

### `mode`

| Value | Meaning |
|---|---|
| 1 | 自定义模式 |
| 2 | 模板模式 |

### `type`

| Value | Meaning |
|---|---|
| 1 | 文件 |
| 2 | 卷 |
| 3 | 整机 |

### `pathType`

| Value | Meaning |
|---|---|
| 1 | 静态文件路径 |
| 2 | 动态文件路径 |

### `ndmpService`

| Value | Meaning |
|---|---|
| 0 | 关闭 |
| 1 | 粗粒度 |
| 2 | 细粒度 |

### `intelligentArchive`

| Value | Meaning |
|---|---|
| 0 | 普通备份 |
| 1 | 智能目录归档 |

### `datasourcePaths[].nodeType`

| Value | Meaning |
|---|---|
| 1 | 目录 |
| 2 | 文件 |

### `timingFilterRule.file.mode` / `timingFilterRule.directory.mode` / `timingFilterRule.format.mode` / `timingFilterRule.time.mode`

| Value | Meaning |
|---|---|
| 1 | 不包含 |
| 2 | 包含 |

### `timingFilterRule.format.fixedTypes[]`

| Value | Meaning |
|---|---|
| 1 | Office |
| 2 | Video |
| 3 | Picture |
| 4 | PDF |
| 5 | Web |
| 6 | Compress |

### `timingFilterRule.time.timeType`

| Value | Meaning |
|---|---|
| 2 | 修改时间 |

### `incrementTraverseRange.unit`

| Value | Meaning |
|---|---|
| 1 | 年 |
| 2 | 月 |
| 3 | 日 |

## Verification

- 已按 [object_detail.go](/e:/go/src/foundation-cli/internal/business/host/object_detail.go) 与 [object_create.go](/e:/go/src/foundation-cli/internal/business/host/object_create.go) 校验返回字段与嵌套结构。
- 当前命令参数 `--object-id` 与设计稿、实现代码一致。

## Examples

```bash
foundation-cli host object detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ object\ detail |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli host object detail
```
