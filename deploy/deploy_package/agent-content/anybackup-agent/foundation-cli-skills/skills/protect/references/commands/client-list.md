# `foundation-cli protect client list`

## 命令摘要

查询指定保护对象关联的客户端列表。

## 示例

```bash
foundation-cli protect client list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect client list --tenant-id <tenant-id> --object-id <object-id>` |
| Method | `GET` |
| Path | `/backupmgm/v2/protect_object/{objectId}/clients` |
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

### 业务参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | 保护对象 ID |

### Query 参数

无。

## 返回结果

来源：`params/protect/client.go` -> `GetObjectClientsResponse`  
嵌套对象定义复用：`params/protect/create_backup_config.go` -> `ClientParams`

### 顶层 envelope

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | any | 错误信息，成功通常为 `null` |
| responseData | object | 业务数据，对应 `GetObjectClientsResponse` |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clients | object[] | 关联客户端列表 |

### responseData.clients[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端唯一标识 |
| clientName | string | 客户端名称 |
| clientOsType | int | 客户端操作系统类型 |
| clientIp | string | 客户端 IP |
| clientStatus | int | 客户端状态 |

## 枚举类型列表

| 字段 | 值 | 含义 |
|---|---|---|
| `clientOsType` | `0`-`7` | 当前源码仅限制取值范围，未给出逐值中文定义 |

说明：`clientStatus` 在当前 `params/protect` 目录下未给出显式枚举定义。

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
    "error":  null,
    "status":  "success",
    "responseData":  {
                         "clients":  [
                                         {
                                             "clientId":  "sample-id",
                                             "clientName":  "sample-name",
                                             "clientOsType":  1
                                         }
                                     ]
                     }
}
```

## 示例命令
```bash
foundation-cli protect client list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```

## 当前环境验证

- 已使用 `object-id=428e4c05bf87deef20a86e670f3eb9a9` 实测成功。
- 当前环境返回 1 个客户端，`clientId=13caa58fd84bab1beb0b7c9d76b7efae`，`clientIp=172.31.0.2`。
