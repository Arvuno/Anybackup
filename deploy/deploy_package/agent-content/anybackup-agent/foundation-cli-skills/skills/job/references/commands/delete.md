# `foundation-cli job delete`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job delete (--job-id <job-id> ... \| --data <json> \| --body-file <file>)` |
| 方法 | `DELETE` |
| 路径 | `/job_center/v1/jobs` |
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
| jobIds | string[] | 是 | 需要删除的作业 ID 列表 |

## 枚举列表

- 无业务枚举字段。

## 删除状态约束

仅以下终态允许执行 `job delete`：

- `600`：失败（Failed）
- `700`：已取消（Canceled）
- `800`：成功（Success）
- `900`：成功但有警告（WarnSuccess）
- `1000`：部分成功（PartSuccess）
- `1100`：已停止（Stopped）

非上述终态不应删除。

## 请求体示例

```json
{
  "jobIds": [
    "f13a27003d391109b50400163e33df73"
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
foundation-cli job delete \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id f13a27003d391109b50400163e33df73
```

```bash
foundation-cli job delete \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"jobIds":["f13a27003d391109b50400163e33df73"]}'
```
