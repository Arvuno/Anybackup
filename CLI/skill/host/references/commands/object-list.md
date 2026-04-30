# `foundation-cli host object list`

## 命令概要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli host object list --tenant-id <tenant-id> [filters]` |
| Method | `GET` |
| Path | `/backupmgm/v1/file/object_list` |
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
| productionSystemId | `--production-system-id` | string | 否 | 主机 ID，长度固定 32 |
| index | `--index` | int | 否 | 分页起始位置，默认 `0`，范围 `0~20000` |
| count | `--count` | int | 否 | 分页数量，默认 `10`，范围 `1~100` |
| name | `--name` | string | 否 | 保护对象名称 |
| datasourceType | `--datasource-type` | int32 | 否 | 数据源类型，允许 `0/1/2/3/4/5` |
| intelligentArchive | `--intelligent-archive` | int32 | 否 | 智能归档过滤，允许 `0/1/2` |
| execStatus | `--exec-status` | int32 | 否 | 执行状态过滤，允许 `0/1/2/3/4` |
| isIncludeTenantId | `--is-include-tenant-id` | bool | 否 | 是否在结果中携带租户 ID，需显式传 `true/false` |
| tenantId | `--query-tenant-id` | string | 否 | 查询条件中的租户 ID，避免与全局 `--tenant-id` 冲突 |
| groupId | `--group-id` | string | 否 | 所属组 ID |
| hasTp | `--has-tp` | int | 否 | 是否存在时间点，允许 `0/1/2` |
| isBackupConfig | `--is-backup-config` | int | 否 | 备份配置状态，允许 `1/2` |
| objectId | `--object-id` | string | 否 | 保护对象 ID，长度固定 32 |

## 返回结果

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int32 | 保护对象总数 |
| data | ProtectObject[] | 主机保护对象列表 |

### `data[]` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| hostId | string | 主机 ID |
| hostName | string | 主机名称 |
| hostIp | string | 主机 IP |
| objectId | string | 保护对象 ID |
| name | string | 保护对象名称 |
| execStatus | ExecStatus[] | 当前执行状态列表 |
| bindStrategyStatus | int | 策略绑定状态 |
| intelligentArchive | int | 智能归档开关状态 |
| latestTimePoint | string | 最新时间点 |
| hasBackupConfig | bool | 是否存在备份配置 |
| backupStrategyId | string | 备份策略 ID |
| backupStrategyName | string | 备份策略名称 |
| datasourceType | int32 | 数据源类型 |
| lastBackupStatus | int32 | 最近一次备份状态 |
| storageServiceName | string | 存储服务名称 |
| hasTp | int | 是否存在时间点 |
| isDeleted | bool | 保护对象是否已删除 |
| clientIsDelete | bool | 客户端是否已删除 |
| backupStrategyNum | int | 绑定策略数量 |
| isGroupSLA | bool | 是否绑定组策略 |

## 示例

### 最小调用

```bash
foundation-cli host object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

### 按主机和对象名称过滤

```bash
foundation-cli host object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --production-system-id 12345678901234567890123456789012 \
  --name fileset-a \
  --index 0 \
  --count 20
```

### 按备份配置和时间点状态过滤

```bash
foundation-cli host object list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --has-tp 1 \
  --is-backup-config 2 \
  --exec-status 4 \
  --group-id group-1
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ object\ list |
| 风险 | read-only |

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
