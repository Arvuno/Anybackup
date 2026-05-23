# `foundation-cli job logs`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job logs --tenant-id <tenant-id> --job-id <job-id> [--index <n>] [--count <n>] [--start-time <ms>] [--end-time <ms>] [--level <1|2|3|4>]` |
| 方法 | `GET` |
| 路径 | `/job_center/v1/activity/{jobId}/logs` |
| 风险 | `只读` |

## 请求参数

### 共享参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation console 基础 URL |
| ak | `--ak` | string | 是 | 访问密钥（AK） |
| sk | `--sk` | string | 是 | 访问密钥密文（SK） |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| jobId | `--job-id` | string | 是 | 作业唯一标识 |
| index | `--index` | int | 否 | 分页起始位置，默认 `0` |
| count | `--count` | int | 否 | 分页数量，默认 `30` |
| startTime | `--start-time` | int64 | 否 | 起始时间（毫秒） |
| endTime | `--end-time` | int64 | 否 | 结束时间（毫秒） |
| level | `--level` | int | 否 | 日志级别，允许 `1/2/3/4` |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 日志总数 |
| data | Log[] | 日志列表 |

### `data[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| logId | string | 日志唯一标识 |
| level | int32 | 日志级别 |
| info | string | 日志内容 |
| code | string | 日志编码 |
| params | object | 日志参数 |
| createTime | int64 | 创建时间（毫秒） |

## 返回案例

```json
{
	"status": "success",
	"error": null,
	"responseData": {
		"totalNum": 10,
		"data": [{
			"logId": "f18b2b206eba11f399480050568943e2",
			"level": 1,
			"info": "开始执行作业。",
			"code": "BackupWorker.StartExecuteTask",
			"params": null,
			"createTime": 1785315276000
		}]
	}
}
```

## 示例

```bash
foundation-cli job logs \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id> \
  --index 0 \
  --count 30
```
