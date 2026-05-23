# `foundation-cli protect backup batch-start`

## 命令摘要

批量触发多个保护对象的备份任务。

## 请求体示例

```json
{
  "objectIds": [
    "sample-object-id-1",
    "sample-object-id-2"
  ],
  "backupType": 1
}
```

## 枚举列表

- `engineType`：常见值包括 `1001`、`1005`。
- `backupType`：常见值包括 `1`、`2`、`3`、`4`、`6`、`7`。

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect backup batch-start --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/batch/backup_tasks/start` |
| Risk | `write` |

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
| body | `--data` | JSON string | 是 | 批量备份触发请求体 |

## Body 参数（`--data`）

对象：`BatchStartBackupTaskRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| objectIds | string[] | 是 | 保护对象 ID 列表，最多 100 条，每项长度 32 |
| engineType | int32 | 否 | 作业类型 |
| backupType | int32 | 是 | 备份方式 |

## 枚举说明

### `engineType`

| 值 | 含义 |
|---|---|
| 1001 | 普通备份 |
| 1005 | 快照作业 |

### `backupType`

| 值 | 含义 |
|---|---|
| 1 | 完全备份 |
| 2 | 增量备份 |
| 3 | 差异备份 |
| 4 | 日志备份 |
| 6 | 快照 |
| 7 | 快照备份 |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 批量触发结果 |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{"status":"success","error":null,"responseData":null}
```

## 示例
```bash
foundation-cli protect backup batch-start \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"objectIds":["<object-id-1>","<object-id-2>"],"backupType":1}'
```
