# `foundation-cli host list`

## 命令概要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli host list --tenant-id <tenant-id> [filters]` |
| Method | `GET` |
| Path | `/resource_center/v1/host/group/host_list` |
| Risk | `read-only` |

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
| runnerTypes | `--runner-types` | string[] | 否 | Runner 类型列表，重复 flag 传入 |
| groupId | `--group-id` | string | 否 | 主机组 ID |
| clientOsFilter | `--client-os-filter` | int32[] | 否 | 客户端操作系统过滤，允许 `0/1/2/3/4/5/6/7`，重复 flag 传入 |
| filter | `--filter` | string | 否 | 模糊过滤文本 |
| isChild | `--is-child` | bool | 否 | 是否查询子主机，需显式传 `true/false` |

## 返回结果

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 主机总数 |
| data | HostInfoV2[] | 主机列表 |

### `data[]` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | 主机唯一标识 |
| name | string | 主机名称 |
| type | int | 生产资源类型 |
| status | int32 | 生产资源状态 |
| access | string | 主机访问地址，通常为 IP |
| genMode | int | 资源生成方式 |
| resourceConstraints | int32 | 主机任务约束 |
| clientOs | int | 客户端操作系统类型 |
| clientOsUser | string | 客户端操作系统用户 |
| tags | string[] | 主机标签 |
| groupId | string | 所属组 ID |
| groupName | string | 所属组名称 |
| groupPath | string | 所属组路径 |
| runnerTypes | string[] | 已安装的 Runner 类型列表 |

## 示例

### 最小调用

```bash
foundation-cli host list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

### 按主机组和操作系统过滤

```bash
foundation-cli host list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --group-id group-1 \
  --client-os-filter 1 \
  --client-os-filter 7 \
  --filter oracle
```

### 按 Runner 类型和子主机范围过滤

```bash
foundation-cli host list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --runner-types VMBackup \
  --runner-types NasBackup \
  --is-child false
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ list |
| 风险 | read-only |

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
