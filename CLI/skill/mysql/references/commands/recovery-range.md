# `foundation-cli mysql recovery range`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql recovery range --tenant-id <tenant-id> [--object-id <id>] [--storage-service-id <id>] [--restore-object-type <type>] [--restore-gran <restore-gran>] [--timestamp <ts>] [--backup-task-type <type>] [--storage-pool-id <id>]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/get_recovery_range` |
| 风险 | `只读` |

`restore-gran` 允许值：`0/1/2/3/4/5/6/7`。

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `9.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 否 | MySQL 保护对象 ID |
| storageServiceId | `--storage-service-id` | string | 否 | 存储服务 ID |
| restoreObjectType | `--restore-object-type` | string | 否 | 恢复目标类型 |
| restoreGran | `--restore-gran` | string | 否 | 允许值： `0/1/2/3/4/5/6/7` |
| timestamp | `--timestamp` | string | 否 | 非负时间戳 |
| backupTaskType | `--backup-task-type` | string | 否 | 备份任务类型 |
| storagePoolId | `--storage-pool-id` | string | 否 | 存储池 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 恢复范围响应 |

### `responseData`

| 字段 | 类型 | 说明 |
|---|---|---|
| data | object[] | 恢复范围列表 |
| totalNum | int | 总记录数 |

### `data[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| incarnation | string | 恢复范围版本标记 |
| begin | int64 | 范围开始时间 |
| end | int64 | 范围结束时间 |
| ctrlend | int64 | 控制结束时间 |
| curtpbegin | int64 | 当前时间点开始时间 |
| absoluteBegin | int64 | 绝对开始时间 |
| timepointid | string | 时间点 ID |
| BackupType | int32 | 备份类型 |
| instanceName | string | 实例名称 |

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
        "incarnation": "1",
        "begin": 1712476800000,
        "end": 1712480400000,
        "ctrlend": 1712480400000,
        "curtpbegin": 1712476800000,
        "absoluteBegin": 1712476800000,
        "timepointid": "sample-timepoint-id",
        "BackupType": 1,
        "instanceName": "mysql3306"
      }
    ]
  }
}
```

## 示例
```bash
foundation-cli mysql recovery range \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --restore-gran 1 \
  --timestamp <timestamp>
```




