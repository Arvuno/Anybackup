# `foundation-cli client detail`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client detail --tenant-id <tenant-id> --client-id <id>` |
| Method | `GET` |
| Path | `/commons/client/{clientId}` |
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

### 路径参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| clientId | `--client-id` | string | 是 | 客户端 ID |

## 返回结果

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息，成功通常为 `null` |
| status | string | 请求状态 |
| responseData | object | 客户端详情数据 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientAppType | int | 客户端应用类型 |
| clientType | int | 客户端分类 |
| clientIp | string | 客户端 IP |
| clientMac | string | 客户端 MAC |
| clientStatus | int | 客户端状态 |
| clientRealStatus | int | 客户端真实状态 |
| clientVersion | string | 客户端版本 |
| clientArch | int | 架构位数 |
| clientOsUser | string | OS 用户 |
| clientCPUNum | int | CPU 核数 |
| clientOSType | string | 操作系统类型 |
| clientOSName | string | 操作系统名称 |
| clientOSKind | int | 操作系统分类 |
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
| updateRunnerNum | int | 待升级 Runner 数量 |
| isUpdated | bool | 是否已升级 |
| isPortable | bool | 是否便携式 |
| basicVersion | string | 基础版本 |
| platform | string | 平台标识 |
| netAccessConfigId | string | 网络访问配置 ID |
| assignable | bool | 是否可分配 |
| creatorName | string | 创建人 |

## 返回案例

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
    "clientName": "iv-yej4xkshkwwh2yram2x1",
    "clientAppType": 102,
    "clientType": 1,
    "clientIp": "172.31.0.2",
    "clientMac": "U0HYTDM3RENXS4EJ",
    "clientStatus": 1,
    "clientRealStatus": 1,
    "clientVersion": "9.0.9.0.378",
    "clientArch": 64,
    "clientOsUser": "",
    "clientCPUNum": 16,
    "clientOSType": "linux",
    "clientOSName": "Linux",
    "clientOSKind": 2,
    "clientMemory": 66496913408,
    "clientIsBuildin": false,
    "clientBinds": [],
    "platformIp": "",
    "port": 0,
    "deleted": 0,
    "isBuiltin": false,
    "username": "",
    "authUsers": [],
    "usersNum": 0,
    "customer": "",
    "clusterPort": "",
    "isEmail": false,
    "cloudHostInfo": "",
    "updateRunnerNum": 0,
    "isUpdated": false,
    "isPortable": false,
    "basicVersion": "008000009000",
    "platform": "Linux_el7_x64",
    "netAccessConfigId": "",
    "assignable": false,
    "creatorName": ""
  }
}
```

## 示例命令

```bash
foundation-cli client detail \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --client-id <client-id>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ detail |
| 风险 | read-only |

## 示例

```bash
foundation-cli client detail
```
