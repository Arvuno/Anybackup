# `foundation-cli protect storage-pool auto-select`

## 命令摘要

根据当前保护配置自动选择可用存储池。

## 示例

```bash
foundation-cli protect storage-pool auto-select \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect storage-pool auto-select --tenant-id <tenant-id> [--exclude-id <id> ...]` |
| Method | `GET` |
| Path | `/backupmgm/v2/auto_select/storage_pool` |
| Risk | `read-only` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation Console 地址 |
| ak | `--ak` | string | 是 | Access Key |
| sk | `--sk` | string | 是 | Secret Key |
| targetVersion | `--target-version` | string | 否 | 目标版本，默认 `9.0.9.0` |

### Query 参数

来源：`params/protect/auto_select_storage_pool.go` -> `AutoSelectStoragePoolRequest`

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| excludeIds | `--exclude-id` | string[] | 否 | 需要排除的存储服务 ID，可重复传递 |

## 返回结果

来源：`params/protect/storage_pool.go` -> `StoragePoolInfo`

### 顶层 envelope

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | any | 错误信息，成功通常为 `null` |
| responseData | object | 业务数据，对应 `StoragePoolInfo` |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| storageSvcId | string | 存储服务唯一标识 |
| storageSvcType | int | 存储服务类型 |
| storageSvcName | string | 存储服务名称 |
| storageSvcStatus | int | 存储服务状态 |
| storagePoolId | string | 存储池唯一标识 |
| storagePoolType | int | 存储池类型 |
| storagePoolName | string | 存储池名称 |
| storagePoolStatus | int | 存储池状态 |

## 枚举类型列表

### `storagePoolType` / `PoolTypeProto`

来源：`params/protect/create_backup_config.go`

| 值 | 含义 |
|---|---|
| `0` | 所有类型 |
| `1` | SAN 存储池 |
| `2` | 本地磁盘存储池 |
| `3` | 分布式存储池 |
| `4` | NAS 存储池 |
| `5` | 磁带 |
| `6` | 磁带原生格式存储池 |
| `7` | 本地文件系统存储池 |
| `8` | 对象存储池 |

### 未显式定义的整型字段

`storageSvcType`、`storageSvcStatus`、`storagePoolStatus` 在当前 `params/protect` 目录下未给出枚举常量，文档仅按原始整型返回展示。

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "storageSvcId": "4235757746e447d18d4d7a8b3f5edc54",
    "storageSvcType": 2,
    "storageSvcName": "StorageService_172.31.12.91",
    "storageSvcStatus": 1,
    "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
    "storagePoolType": 2,
    "storagePoolName": "block",
    "storagePoolStatus": 1
  }
}
```

## 示例命令
```bash
foundation-cli protect storage-pool auto-select \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

```bash
foundation-cli protect storage-pool auto-select \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --exclude-id <storage-service-id-1> \
  --exclude-id <storage-service-id-2>
```

## 当前环境验证

- 已使用当前测试环境实测成功。
- 当前环境返回 `storageSvcId=4235757746e447d18d4d7a8b3f5edc54`、`storagePoolId=0c6022822f1d11f19b0600163e33df73`、`storagePoolName=block`。
