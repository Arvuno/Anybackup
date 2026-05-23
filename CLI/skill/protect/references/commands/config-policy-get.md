# `foundation-cli protect config-policy get`

## 命令摘要

查询指定保护对象的配置策略可见性。

## 示例

```bash
foundation-cli protect config-policy get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect config-policy get --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v2/{objectId}/config_police` |
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

### 业务参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | 保护对象 ID |

### Query 参数

无。

## 返回结果

来源：`params/protect/config_police.go` -> `PoliceConfig`

### 顶层 envelope

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | any | 错误信息，成功通常为 `null` |
| responseData | object | 业务数据，对应 `PoliceConfig` |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| appType | string | 应用类型 |
| backupDestination | object | 备份目的地策略可见性 |
| basicConfig | object | 基础策略可见性 |

### responseData.backupDestination 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| multipleRegions.visible | bool | 是否展示多区域配置 |
| disturbanceSwitching.visible | bool | 是否展示故障切换配置 |
| storagePool.visible | bool | 是否展示存储池配置 |
| deduplication.visible | bool | 是否展示去重配置 |
| dataEncryption.visible | bool | 是否展示数据加密配置 |
| dataCompression.visible | bool | 是否展示数据压缩配置 |
| dataConsistency.visible | bool | 是否展示数据一致性配置 |

### responseData.basicConfig 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| transEncryption.visible | bool | 是否展示传输加密配置 |
| failedRetry.visible | bool | 是否展示失败重试配置 |
| forcedRetention.visible | bool | 是否展示强制保留配置 |

## 枚举类型列表

该接口返回结构中未定义显式枚举。

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
    "appType": "vmware",
    "backupDestination": {
      "multipleRegions": {
        "visible": true
      },
      "disturbanceSwitching": {
        "visible": true
      },
      "storagePool": {
        "visible": true
      },
      "deduplication": {
        "visible": true
      },
      "dataEncryption": {
        "visible": true
      },
      "dataCompression": {
        "visible": true
      },
      "dataConsistency": {
        "visible": false
      }
    },
    "basicConfig": {
      "transEncryption": {
        "visible": true
      },
      "failedRetry": {
        "visible": true
      },
      "forcedRetention": {
        "visible": false
      }
    }
  }
}
```

## 示例命令
```bash
foundation-cli protect config-policy get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 当前环境验证

- 已使用 `object-id=428e4c05bf87deef20a86e670f3eb9a9` 实测成功。
- 当前环境返回 `appType=eso_backupengine_hypermysqlengine`。
