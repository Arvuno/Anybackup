# `foundation-cli job business-types`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job business-types --tenant-id <tenant-id>` |
| 方法 | `GET` |
| 路径 | `/job_center/v1/business_types` |
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

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| data | OperationBusinessType[] | 一级业务类型及二级业务类型列表 |

### `data[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| operationType | int32 | 一级业务类型 |
| name | string | 一级业务类型名称 |
| businessTypes | BusinessType[] | 二级业务类型列表 |

### `data[].businessTypes[]` 字段（`BusinessType[]`）

| 字段 | 类型 | 说明 |
|---|---|---|
| businessType | int32 | 二级业务类型编码 |
| name | string | 二级业务类型名称 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "operationType": 1,
        "name": "备份与恢复",
        "businessTypes": [
          {
            "businessType": 1,
            "name": "完全备份"
          },
          {
            "businessType": 2,
            "name": "增量备份"
          }
        ]
      }
    ]
  }
}
```

## 示例

```bash
foundation-cli job business-types \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```
