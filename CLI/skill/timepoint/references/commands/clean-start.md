# `foundation-cli timepoint clean start`

## 命令概览

| 项 | 值 |
|---|---|
| 命令 | `foundation-cli timepoint clean start --object-id <object-id> --data '<json>'` |
| 请求方法 | `POST` |
| 接口路径 | `/backupmgm/v1/protect_object/{objectId}/clean_task/start` |
| 风险级别 | `write` |

## 请求参数

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 服务地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 秘密密钥 |
| `--object-id` | 是 | 保护对象 ID |
| `--data` 或 `--body-file` | 是 | 清理任务请求体 |

## Body 参数（`--data`）

顶层对象：`StartCleanTaskRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| cleanAll | bool | 否 | 是否清理全部时间点。该字段已废弃，仅兼容旧版本，后续建议使用 `cleanType` |
| isCleanSubObject | bool | 否 | 是否清理子对象时间点 |
| cleanType | int | 否 | 清理类型，取值范围见“枚举说明（CleanTypeProto）” |
| timePointIds | string[] | 否 | 时间点 ID 列表，最多 `100` 个，元素长度必须为 `32` |
| dataSetIds | string[] | 否 | 数据集 ID 列表，最多 `100` 个，元素长度必须为 `32` |
| startTime | int64 | 否 | 起始时间戳（毫秒），`>=0` |
| endTime | int64 | 否 | 结束时间戳（毫秒），`>=0` |
| storagePoolId | string | 否 | 存储池 ID，长度 `32~48` |
| storagePoolType | int32 | 否 | 存储池类型 |
| storageServiceId | string | 否 | 存储服务 ID，长度 `1~48` |

参数组合建议（按清理类型）：

| 场景 | 推荐字段组合 |
|---|---|
| 清理指定时间点 | `cleanType=2` + `timePointIds` |
| 清理时间区间数据 | `cleanType=7` + `startTime` + `endTime` |
| 清理时间点和数据集信息 | `cleanType=8` + `timePointIds` / `dataSetIds` |
| 清理指定存储池或存储服务数据 | `cleanType=9` + `storageServiceId`（可选配 `storagePoolId`、`storagePoolType`） |

## 说明

1. `timepoint clean start` 保持后端请求结构不变，仅调整 CLI 命令路径。
2. 可使用 `--data` 传入内联 JSON，也可使用 `--body-file` 传入 JSON 文件。
3. 缺少请求体时会在本地参数校验阶段直接失败，不会发起后端请求。
4. `startTime` 和 `endTime` 的时间戳单位为毫秒。

## 请求结构体说明

### 清理全部时间点

适用场景：按存储服务和存储池维度发起“全量时间点清理”。

```json
{
  "cleanAll": true,
  "storagePoolId": "e177c3fe311f11f1a1220050568943e2",
  "storageServiceId": "9fc189d638944a57ac042874a21d94b7",
  "storagePoolType": 2
}
```

### 清理单个时间点

适用场景：仅清理指定的一个或多个时间点 ID。

```json
{
  "cleanAll": false,
  "timePointIds": [
    "f18a2425d96a11169ef40050568943e2"
  ]
}
```

关键字段说明：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| cleanAll | bool | 是 | `true` 表示清理全部时间点；`false` 表示按 `timePointIds` 精确清理 |
| storagePoolId | string | 条件必填 | `cleanAll=true` 时建议传入，表示目标存储池 ID |
| storageServiceId | string | 条件必填 | `cleanAll=true` 时建议传入，表示目标存储服务 ID |
| storagePoolType | int | 条件必填 | `cleanAll=true` 时建议传入，表示存储池类型 |
| timePointIds | string[] | 条件必填 | `cleanAll=false` 时必填，表示待清理时间点 ID 列表 |

## 枚举说明

### CleanTypeProto

| 枚举值 | 枚举名 | 说明 |
|---|---|---|
| `0` | `UNKNOWN_CLEAN` | 未知清理类型 |
| `1` | `CLEAN_ALL_DATA` | 清理全部数据 |
| `2` | `CLEAN_SPECIFIED_DATA` | 清理指定数据 |
| `3` | `CLEAN_BACKUP_DATA` | 清理全部备份数据 |
| `4` | `CLEAN_SNAPSHOT_DATA` | 清理全部快照数据 |
| `5` | `CLEAN_APP_RESOURCE` | 清理应用资源 |
| `6` | `CLEAN_DATASET_DATA` | 清理文件归档数据集数据 |
| `7` | `CLEAN_TIME_INTERVAL_DATA` | 清理时间区间数据 |
| `8` | `CLEAN_TIME_AND_DATASET` | 清理时间点和数据集融合数据 |
| `9` | `CLEAN_SPECIFIED_STORAGE_POOL` | 清理指定存储池或存储服务数据 |
| `10` | `CLEAN_ALL_LOG_DATA` | 清理全部日志数据 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "taskId": "f19ad7647e6c116299480050568943e2"
  }
}
```

## 示例
```bash
foundation-cli timepoint clean start \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --data '{"cleanAll":false,"timePointIds":["<timepoint-id>"]}'
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli timepoint clean start |
| 风险 | read-only |
