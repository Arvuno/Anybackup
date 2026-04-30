# `foundation-cli client list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client list --tenant-id <tenant-id> [--index <n>] [--count <n>] [--type <type>] [--status <status>] [--client-type <type>]` |
| Method | `GET` |
| Path | `/commons/all_clients` |
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

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| index | `--index` | int | 否 | 分页起始位置；默认 `0` |
| count | `--count` | int | 否 | 分页数量；默认 `30` |
| type | `--type` | string | 否 | 客户端应用类型过滤（透传到 `type`）；默认 `0` |
| status | `--status` | string | 否 | 客户端状态过滤（透传到 `status`）；默认 `2` |
| clientType | `--client-type` | string | 否 | 客户端分类过滤（透传到 `clientType`）；默认 `0` |

## 返回结果

依据 `params/client/all_clients.go` 中 `ClientResponse`、`ResponseData`、`Client`。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息，成功通常为 `null` |
| status | string | 请求状态 |
| responseData | object | 业务数据 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int | 总数 |
| filter | any | 后端过滤对象 |
| data | Client[] | 客户端列表 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientAppType | int | 客户端应用类型，见枚举 |
| clientType | int | 客户端分类，见枚举 |
| clientIp | string | 客户端 IP |
| clientMac | string | 客户端 MAC |
| clientStatus | int | 客户端状态 |
| clientRealStatus | int | 客户端真实状态 |
| clientVersion | string | 客户端版本 |
| clientArch | int | 架构位数 |
| clientOsUser | string | OS 用户 |
| clientCPUNum | int | CPU 核数 |
| clientOSType | int | 操作系统类型，见枚举 |
| clientOSName | string | 操作系统名称 |
| clientOSKind | int | OS 种类，见枚举 |
| clientMemory | int64 | 内存大小 |
| clientIsBuildin | bool | 是否内置 |
| clientBinds | object[] | 绑定信息列表 |
| platformIp | string | 平台 IP |
| port | int | 端口 |
| deleted | int | 删除标记 |
| isBuiltin | bool | 是否内置（冗余字段） |
| username | string | 用户名 |
| authUsers | string[] | 授权用户列表 |
| usersNum | int | 用户数量 |
| customer | string | 客户名 |
| clusterPort | string | 集群端口 |
| isEmail | bool | 是否邮件通知 |
| cloudHostInfo | string | 云主机信息 |
| updateRunnerNum | int | 待升级 runner 数量 |
| isUpdated | bool | 是否已升级 |
| isPortable | bool | 是否便携式 |
| basicVersion | string | 基础版本 |
| platform | string | 平台名称 |
| netAccessConfigId | string | 网络访问配置 ID |
| assignable | bool | 是否可分配 |
| creatorName | string | 创建人 |
| netAccessConfigName | string | 网络访问配置名称 |
| isBindNetAccessConfig | bool | 是否绑定网络访问配置 |

## 枚举类型列表

### ClientType (`params/client/all_clients.go`)

| 值 | 含义 |
|---|---|
| 0 | All |
| 1 | Builtin |
| 2 | Physics |
| 3 | Virtual |

### OsType (`params/client/all_clients.go`)

| 值 | 含义 |
|---|---|
| 1 | Windows |
| 2 | Linux |
| 3 | AIX |
| 4 | Solaris |
| 5 | HP-UX |
| 6 | NeoKylin |
| 7 | YHKylin |

### ClientOSKind (`params/client/all_clients.go`)

| 值 | 含义 |
|---|---|
| 1 | Normal Windows |
| 2 | Normal Linux |
| 4 | Windows PE |
| 8 | Linux PE |
| 16 | Drill Machine |

### ClientApplicationType (`params/client/all_clients.go`)

| 值范围 | 说明 |
|---|---|
| 0 | All |
| 2~52（离散） | 应用客户端类型（如 Oracle、VMware、MySQL、TiDB 等） |
| 101~107 | 按 OS 分类的客户端类型（Windows/Linux/AIX/Solaris/HP-UX/NeoKylin/YHKylin） |

注：该枚举常量较长，完整值集合以 `params/client/all_clients.go` 的 `ClientApplicationType` 常量定义为准。

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
                         "totalNum":  1,
                         "data":  [
                                      {
                                          "clientId":  "sample-id",
                                          "clientName":  "sample-name",
                                          "clientAppType":  1
                                      }
                                  ]
                     }
}
```

## 示例命令
```bash
foundation-cli client list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --index 0 \
  --count 30 \
  --status 2
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ list |
| 风险 | read-only |

## 示例

```bash
foundation-cli client list
```
