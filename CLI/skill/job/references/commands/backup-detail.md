# `foundation-cli job backup-detail`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | `foundation-cli job backup-detail --tenant-id <tenant-id> --job-id <job-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/task/{jobId}/detail` |
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
| jobId | `--job-id` | string | 是 | 备份作业唯一标识 |

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| taskId | string | 任务唯一标识 |
| objectId | string | 保护对象唯一标识 |
| objectName | string | 保护对象名称（由 `/job_center/v1/job/{jobId}` 补充） |
| productionSystemId | string | 生产系统唯一标识 |
| appType | string | 应用类型 |
| operationType | int32 | 一级业务类型 |
| businessType | int32 | 二级业务类型 |
| status | int32 | 作业状态 |
| isAutoExecute | bool | 是否自动执行 |
| createUserId | string | 创建用户 ID |
| createUserName | string | 创建用户名称 |
| startTime | int64 | 开始时间（毫秒） |
| endTime | int64 | 结束时间（毫秒） |
| runTime | int64 | 运行时长（毫秒） |
| retryCount | int32 | 配置重试次数 |
| retriedCount | int32 | 已重试次数 |
| rawData | int64 | 原始数据量 |
| completedData | int64 | 已完成数据量 |
| sendSize | int64 | 发送数据量 |
| storedSize | int64 | 存储数据量 |
| speed | string | 执行速度 |
| deduplication | int32 | 去重开关状态 |
| resourceTotal | int64 | 资源总数 |
| resourceBackup | int64 | 备份资源数 |
| resourceSuccess | int64 | 成功资源数 |
| resourceFailed | int64 | 失败资源数 |
| scannedDirCount | int64 | 已扫描目录数 |
| scannedFileCount | int64 | 已扫描文件数 |
| appExtraInfo | object | 应用扩展信息 |
| isPermanentIncrement | bool | 是否永久增量时间点 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "taskId": "f134a7e65114116ab50400163e33df73",
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
    "productionSystemId": "a63ba04bfb530fefa3b09d1f89d47432",
    "appType": "eso_backupengine_hypermysqlengine",
    "operationType": 1,
    "businessType": 1,
    "status": 800,
    "isAutoExecute": false,
    "createUserId": "cdec38822f1811f19fd400163e33df73",
    "createUserName": "admin",
    "startTime": 1775803115000,
    "endTime": 1775803127000,
    "runTime": 12000,
    "retryCount": 0,
    "retriedCount": 0,
    "rawData": 19901940,
    "completedData": 19901940,
    "sendSize": 19901940,
    "storedSize": 19901940,
    "speed": "4.74 MiB/s",
    "deduplication": 0,
    "resourceTotal": 0,
    "resourceBackup": 0,
    "resourceSuccess": 0,
    "resourceFailed": 0,
    "scannedDirCount": 0,
    "scannedFileCount": 0,
    "appExtraInfo": {},
    "isPermanentIncrement": true,
    "objectName": "mysql3306_U0HYTDM3RENXS4EJ"
  }
}
```

## 示例

```bash
foundation-cli job backup-detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --job-id <job-id>
```
