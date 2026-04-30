# `foundation-cli job sync-detail`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job sync-detail --tenant-id <tenant-id> --job-id <job-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/sync_task/{jobId}/detail` |
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

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| jobId | `--job-id` | string | 是 | 复制作业唯一标识 |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| jobId | string | 作业唯一标识 |
| status | int32 | 作业状态 |
| appType | string | 应用类型 |
| startTime | int64 | 开始时间（毫秒） |
| endTime | int64 | 结束时间（毫秒） |
| runTime | int64 | 运行时长 |
| storageSvcName | string | 存储服务名称 |
| storagePoolName | string | 存储池名称 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "jobId": "f138a59a121f1161b50400163e33df73",
    "status": 400,
    "appType": "MySQL",
    "startTime": 1775803115000,
    "endTime": 1775803145000,
    "runTime": 30000,
    "storageSvcName": "StorageService_192.168.1.1",
    "storagePoolName": "block"
  }
}
```

## 示例

```bash
foundation-cli job sync-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id>
```
