# `foundation-cli job child list`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job child list --tenant-id <tenant-id> --job-id <job-id> [--index <n>] [--count <n>]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/task/{jobId}/child` |
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
| jobId | `--job-id` | string | 是 | 父作业 ID，同时映射到路径参数和 query 中的 `taskId` |
| index | `--index` | int | 否 | 分页起始位置，默认 `0` |
| count | `--count` | int | 否 | 分页数量，默认 `30` |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "data": [
      {
        "taskId": "f134a7e651141183b50400163e33df73",
        "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
        "objectName": "mysql3306_U0HYTDM3RENXS4EJ",
        "productionSystemId": "a63ba04bfb530fefa3b09d1f89d47432",
        "appType": "eso_backupengine_hypermysqlengine",
        "businessType": 1,
        "status": 800,
        "startTime": 1775803115000,
        "endTime": 1775803127000,
        "runTime": 12000,
        "speed": "4.74 MiB/s",
        "completedData": 19901940,
        "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
        "clientName": "iv-yej4xkshkwwh2yram2x1",
        "clientIp": "172.31.0.2"
      }
    ],
    "totalNum": 1
  }
}
```

## 示例

```bash
foundation-cli job child list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id f134a7e65114116ab50400163e33df73 \
  --index 0 \
  --count 30
```
