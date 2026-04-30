# `foundation-cli host restore start`

## Command Summary

| Item | Value |
|---|---|
| CLI | `foundation-cli host restore start --tenant-id <tenant-id> --data '<json>'` |
| Method | `POST` |
| Path | `/backupmgm/v1/file/recovery` |
| Risk | `write` |

## Request Parameters

### Shared Parameters

| Field | CLI Flag | Type | Required | Description |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | No | Tenant identifier; can come from FOUNDATION_TENANT_ID env var |
| endpoint | `--endpoint` | string | yes | Foundation base URL |
| ak | `--ak` | string | yes | Access Key |
| sk | `--sk` | string | yes | Secret Key |
| targetVersion | `--target-version` | string | no | Target backend version, default `9.0.9.0` |

### Command Parameters

| Field | CLI Flag | Type | Required | Description |
|---|---|---|---|---|
| data | `--data` / `-d` | string | yes | Inline JSON request body |
| bodyFile | `--body-file` | string | no | JSON body file path (CLI reads it and sends as request body) |

## Request Body

`host restore start` directly forwards the JSON body to backend. Body schema follows backend contract for `/backupmgm/v1/file/recovery`.

## Response

Response is passthrough JSON from backend. If response includes `jobId`, follow with:

```bash
foundation-cli job backup-detail --job-id <job-id>
```

## Examples

```bash
foundation-cli host restore start --tenant-id <tenant-id> --data '<json>'
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ host\ restore\ start |
| 风险 | write |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

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

## 示例

```bash
foundation-cli host restore start
```
