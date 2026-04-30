# `foundation-cli client runner-types`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client runner-types --tenant-id <tenant-id>` |
| Method | `GET` |
| Path | `/commons/client/runnerTypes` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `8.0.9.0` |

### Query 参数

无。

## 返回结果

依据 `params/client/runner_types.go`（样例 JSON + `TypeResponseData`/`Data` 结构）。

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
| data | Data[] | runner 类型列表 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| runnerType | string | runner 类型标识，例如 `Basic`、`Oracle` |
| text | string | 中文名称 |
| textEn | string | 英文名称 |

## 枚举类型列表

此接口在 `params/client/runner_types.go` 中未定义固定数值枚举，`runnerType` 由后端字典返回，建议按实际响应动态处理。

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
        "runnerType": "Basic",
        "text": "基础 Runner",
        "textEn": "Basic Runner"
      }
    ]
  }
}
```

## 示例命令
```bash
foundation-cli client runner-types \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ runner-types |
| 风险 | read-only |

## 示例

```bash
foundation-cli client runner-types
```
