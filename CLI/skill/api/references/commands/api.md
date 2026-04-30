# `foundation-cli api`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli api --tenant-id <tenant-id> --method <method> --path </relative/path>` |
| 方法 | `GET` / `POST` / `PUT` / `PATCH` / `DELETE` |
| 路径 | 用户传入的 Foundation 相对 REST 路径 |
| 风险 | `按 method 判定` |

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
| method | `--method` | string | 是 | HTTP 方法 |
| path | `--path` | string | 是 | Foundation API 相对路径 |
| data | `--data` | string | 否 | 写入请求的 JSON 请求体 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [],
    "totalNum": 0
  }
}
```

## 示例

只读透传：

```bash
foundation-cli api \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --method GET \
  --path /job_center/v1/business_types
```

写入透传：

```bash
foundation-cli api \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --method POST \
  --path /deploy/v1/hostConfig/linux \
  --data '{"name":"demo","script":"echo ok"}'
```
