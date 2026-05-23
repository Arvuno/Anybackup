# `foundation-cli job stop`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job stop (--job-id <job-id> ... \| --data <json> \| --body-file <file>)` |
| 方法 | `POST` |
| 路径 | `/job_center/v1/jobs/stop` |
| 风险 | `写入` |

## 请求参数

### 共享参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation console URL |
| ak | `--ak` | string | 是 | 访问密钥（AK） |
| sk | `--sk` | string | 是 | 访问密钥（SK） |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| jobId | `--job-id` | string[] | 条件必填 | 可重复；不传 `--data/--body-file` 时必填 |
| data | `--data` | string | 条件必填 | 内联 JSON；与 `--job-id` 互斥 |
| bodyFile | `--body-file` | string | 条件必填 | JSON 文件；与 `--job-id` 互斥 |

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| jobIds | string[] | 是 | 需要终止的作业 ID 列表 |

## 枚举列表

- 无业务枚举字段。

## 停止状态约束

仅以下状态允许执行 `job stop`：

- `200`：准备中（Ready）
- `300`：正在执行（Running）
- `310`：重试中（Retrying）
- `500`：异常（Abnormal）

以下状态不应再执行停止：

- `100`（排队中）、`400`（停止中）、`600/700/800/900/1000/1100`（终态）

## 请求体示例

```json
{
  "jobIds": [
    "f13c8b522b4111e8b50400163e33df73"
  ]
}
```

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": null
}
```

## 示例

```bash
foundation-cli job stop \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id f13c8b522b4111e8b50400163e33df73
```

```bash
foundation-cli job stop \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"jobIds":["f13c8b522b4111e8b50400163e33df73"]}'
```
