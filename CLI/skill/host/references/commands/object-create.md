# `foundation-cli host object create`

## 命令概要

| Item | Value |
|---|---|
| CLI | `foundation-cli host object create --tenant-id <tenant-id> --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/file/fileset` |
| Risk | `write` |

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
| data | `--data`, `-d` | string | 否 | 内联 JSON 请求体（与 `--body-file` 二选一；最终必须提供请求体） |
| bodyFile | `--body-file` | string | 否 | JSON 请求体文件路径（与 `--data/-d` 二选一；最终必须提供请求体） |

## Body 字段

### 顶层字段

- `hostInfos`
  - 类型：`HostInfo[]`
  - 必填：否
  - 说明：主机信息列表
- `name`
  - 类型：`string`
  - 必填：是
  - 说明：文件集（保护对象）名称
  - 约束：`nameform1`；长度 `3~256`
- `mode`
  - 类型：`int`
  - 必填：是
  - 说明：数据源模式
  - 约束：允许 `1/2`：`1` 自定义，`2` 模板
- `type`
  - 类型：`int`
  - 必填：否
  - 说明：自定义模式下的数据源类型
  - 约束：允许 `0/1/2/3`：`0` 未知，`1` 文件，`2` 卷，`3` 整机（注释提示部分类型暂不支持时，以后端为准）
- `pathType`
  - 类型：`int`
  - 必填：否
  - 说明：路径类型
  - 约束：允许 `0/1/2`：`0` 未知，`1` 静态文件路径，`2` 动态文件路径（注释提示 `2` 暂不支持时，以后端为准）
- `datasourcePaths`
  - 类型：`DataSourceInf[]`
  - 必填：否
  - 说明：数据源路径列表
- `ndmpService`
  - 类型：`int`
  - 必填：否
  - 说明：NAS 存储 NDMP 服务快照备份恢复方式
  - 约束：允许 `0/1/2`：`0` 关闭，`1` 粗粒度，`2` 细粒度
- `intelligentArchive`
  - 类型：`int`
  - 必填：否
  - 说明：是否智能目录归档
  - 约束：允许 `0/1`：`0` 非智能目录归档，`1` 智能目录归档
- `startDir`
  - 类型：`string`
  - 必填：否
  - 说明：起始目录（注释提示暂未启用时，以后端为准）
- `templateId`
  - 类型：`string`
  - 必填：否
  - 说明：模板 ID
  - 约束：长度固定 `32`
- `templateName`
  - 类型：`string`
  - 必填：否
  - 说明：模板名称
  - 约束：`nameform1`；长度 `3~256`
- `includeFilterRule`
  - 类型：`FilterRule[]`
  - 必填：否
  - 说明：智能目录的包含选项
- `excludeFilterRule`
  - 类型：`FilterRule[]`
  - 必填：否
  - 说明：智能目录的排除选项
- `openTimingFilterRule`
  - 类型：`bool`
  - 必填：否
  - 说明：定时备份过滤规则开关
- `timingFilterRule`
  - 类型：`TimingFilterRule`
  - 必填：否
  - 说明：普通定时过滤规则
- `incrementTraverseRange`
  - 类型：`IncrementTraverseRange`
  - 必填：否
  - 说明：增量遍历范围（数据保留期限）
- `protectGranularity`
  - 类型：`int32`
  - 必填：否
  - 说明：保护层级（注释提示默认 `1` 时，以后端默认值为准）

### `hostInfos[]` 字段

- `id`
  - 类型：`string`
  - 必填：是
  - 说明：主机 ID
  - 约束：长度固定 `32`
- `ip`
  - 类型：`string`
  - 必填：否
  - 说明：主机 IP
- `name`
  - 类型：`string`
  - 必填：否
  - 说明：主机名称
- `os`
  - 类型：`int`
  - 必填：是
  - 说明：主机操作系统大类
  - 约束：允许 `101/102/103`
- `clientOs`
  - 类型：`int`
  - 必填：是
  - 说明：主机操作系统
  - 约束：允许 `0/1/2/3/4/5/6/7`（`0` unknown，`1` windows，`2` linux，`3` aix，`4` solaris，`5` hp_ux，`6` neokylin，`7` yhkylin）

### `datasourcePaths[]` 字段

- `fullPath`
  - 类型：`string`
  - 必填：否
  - 说明：数据源路径
- `nodeType`
  - 类型：`int`
  - 必填：否
  - 说明：节点类型
  - 约束：允许 `1/2`：`1` 目录，`2` 文件
- `type`
  - 类型：`string`
  - 必填：否
  - 说明：数据源类型（例如 disk/volume；字段的 JSON tag 在代码中存在异常空格，但按字段名约定为 `type`）

### `includeFilterRule[]` 字段

- `path`
  - 类型：`object`
  - 必填：否
  - 说明：路径过滤
  - 字段：
    - `pathContents`：`PathContent[]`
- `format`
  - 类型：`object`
  - 必填：否
  - 说明：格式过滤
  - 字段：
    - `fixedTypes`：`int[]`
    - `otherTypes`：`string[]`
- `time`
  - 类型：`object`
  - 必填：否
  - 说明：时间过滤（时间戳单位为秒/毫秒需以后端契约为准；参数类型为 `int64`）
  - 字段：
    - `fromTimestamp`：`int64`
    - `toTimestamp`：`int64`

### `excludeFilterRule[]` 字段

- `path`
  - 类型：`object`
  - 必填：否
  - 说明：路径过滤
  - 字段：
    - `pathContents`：`PathContent[]`
- `format`
  - 类型：`object`
  - 必填：否
  - 说明：格式过滤
  - 字段：
    - `fixedTypes`：`int[]`
    - `otherTypes`：`string[]`
- `time`
  - 类型：`object`
  - 必填：否
  - 说明：时间过滤（时间戳单位为秒/毫秒需以后端契约为准；参数类型为 `int64`）
  - 字段：
    - `fromTimestamp`：`int64`
    - `toTimestamp`：`int64`

### `timingFilterRule` 字段

- `file`
  - 类型：`TimingFileFilter`
  - 必填：否
  - 说明：文件过滤规则
- `directory`
  - 类型：`TimingDirectoryFilter`
  - 必填：否
  - 说明：目录过滤规则
- `format`
  - 类型：`TimingFormatFilter`
  - 必填：否
  - 说明：格式过滤规则
- `time`
  - 类型：`TimingTimeFilter`
  - 必填：否
  - 说明：时间过滤规则

### `timingFilterRule.file` 字段

- `mode`
  - 类型：`int`
  - 必填：否
  - 说明：过滤模式
  - 约束：允许 `1/2`：`1` 不包含，`2` 包含
- `paths`
  - 类型：`string[]`
  - 必填：否
  - 说明：过滤路径列表

### `timingFilterRule.directory` 字段

- `mode`
  - 类型：`int`
  - 必填：否
  - 说明：过滤模式
  - 约束：允许 `1/2`：`1` 不包含，`2` 包含
- `paths`
  - 类型：`string[]`
  - 必填：否
  - 说明：过滤路径列表

### `timingFilterRule.format` 字段

- `mode`
  - 类型：`int`
  - 必填：否
  - 说明：过滤模式
  - 约束：允许 `1/2`：`1` 不包含，`2` 包含
- `fixedTypes`
  - 类型：`int[]`
  - 必填：否
  - 说明：固定格式类型
  - 约束：允许 `1/2/3/4/5/6`：`1` office，`2` video，`3` picture，`4` pdf，`5` web，`6` compress
- `otherTypes`
  - 类型：`string[]`
  - 必填：否
  - 说明：其它格式类型

### `timingFilterRule.time` 字段

- `mode`
  - 类型：`int`
  - 必填：否
  - 说明：过滤模式
  - 约束：允许 `1/2`：`1` 不包含，`2` 包含
- `timeType`
  - 类型：`int`
  - 必填：否
  - 说明：时间类型
  - 约束：注释提示 `2` 表示修改时间
- `fromTimestamp`
  - 类型：`int64`
  - 必填：否
  - 说明：起始时间（时间戳）
- `toTimestamp`
  - 类型：`int64`
  - 必填：否
  - 说明：结束时间（时间戳）

### `incrementTraverseRange` 字段

- `open`
  - 类型：`bool`
  - 必填：否
  - 说明：是否开启增量遍历范围
- `value`
  - 类型：`int`
  - 必填：否
  - 说明：数值
- `unit`
  - 类型：`int`
  - 必填：否
  - 说明：单位
  - 约束：允许 `0/1/2/3`（注释提示 `1` 年，`2` 月，`3` 日）

## 返回结果

CLI 会原样透传后端返回的 JSON。

如果返回结果中包含异步任务标识（例如 `jobId`），可继续使用 `foundation-cli job backup-detail --job-id <job-id>` 或 `foundation-cli job sync-detail --job-id <job-id>` 查询任务详情。

## 示例

### 最小命令示例

```bash
foundation-cli host object create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"hostInfos":[{"id":"12345678901234567890123456789012","os":101,"clientOs":1}],"name":"fileset-a","mode":1}'
```

### 最小 JSON body 示例

```json
{
  "hostInfos": [
    {
      "id": "12345678901234567890123456789012",
      "os": 101,
      "clientOs": 1
    }
  ],
  "name": "fileset-a",
  "mode": 1
}
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ object\ create |
| 风险 | write |

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| data | object | 是 | 请求体对象，建议按“请求体示例”中的字段组织 |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
