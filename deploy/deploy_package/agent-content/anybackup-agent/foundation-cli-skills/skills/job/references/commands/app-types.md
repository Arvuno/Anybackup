# `foundation-cli job app-types`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job app-types --tenant-id <tenant-id>` |
| 方法 | `GET` |
| 路径 | `/job_center/v1/app_types` |
| 风险 | `只读` |

## 请求参数

### 共享参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation console 基础 URL |
| ak | `--ak` | string | 是 | 访问密钥（AK） |
| sk | `--sk` | string | 是 | 访问密钥密文（SK） |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `8.0.9.0` |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| data | AppType[] | 应用类型列表 |

### `data[]` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| appType | string | 应用类型编码 |
| name | string | 应用类型名称 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "appType": "eso_backupengine_fileengine",
        "name": "文件"
      },
      {
        "appType": "eso_backupengine_hyperoracleengine",
        "name": "Oracle数据库"
      },
      {
        "appType": "eso_backupengine_hypermysqlengine",
        "name": "MySQL数据库"
      },
      {
        "appType": "eso_backupengine_hypersqlserverengine",
        "name": "SQL Server"
      },
      {
        "appType": "eso_backupengine_vmwareengine",
        "name": "VMware"
      },
      {
        "appType": "eso_backupengine_hypervengine",
        "name": "Hyper-V"
      }
    ]
  }
}
```

## 示例

```bash
foundation-cli job app-types \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```
