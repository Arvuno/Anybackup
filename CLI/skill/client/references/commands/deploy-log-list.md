# `foundation-cli client deploy-log list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client deploy-log list --tenant-id <tenant-id> --job-id <job-id> [--index <n>] [--count <n>]` |
| Method | `GET` |
| Path | `/deploy/v1/job/jobLog` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### Query 参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| jobId | `--job-id` | string | 是 | 部署作业 ID |
| index | `--index` | int | 否 | 分页起始位置 |
| count | `--count` | int | 否 | 分页数量 |

## 返回结果

依据 `params/client/job_log.go` 中 `JobLogResponse`、`JobLogResponseData`、`LogEntry`。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息 |
| status | string | 请求状态 |
| responseData | object | 业务数据 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int | 总数 |
| data | LogEntry[] | 日志列表 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| time | int64 | 时间戳（毫秒） |
| level | int | 日志级别，见枚举 |
| logInfo | string | 日志内容 |
| errorCode | string | 错误码 |
| errorArgs | map[string]string | 错误参数 |

## 枚举类型列表

### ExecLogLevel (`params/client/job_log.go`)

| 值 | 含义 |
|---|---|
| 1 | INFO |
| 2 | WARNING |
| 3 | ERROR |
| 4 | FATAL |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "totalNum": 1,
    "data": [
      {
        "time": 1712476800000,
        "level": 1,
        "logInfo": "deploy runner package to host 10.10.10.10",
        "errorCode": "",
        "errorArgs": {}
      }
    ]
  }
}
```

## 示例命令
```bash
foundation-cli client deploy-log list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ deploy-log\ list |
| 风险 | read-only |

## 示例

```bash
foundation-cli client deploy-log list
```
