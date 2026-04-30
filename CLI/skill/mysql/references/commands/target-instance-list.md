# `foundation-cli mysql target-instance list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql target-instance list --tenant-id <tenant-id> [--index <n>] [--count <n>]` |
| 方法 | `GET` |
| 路径 | `/resource_center/v1/databasealone/instance_list` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| index | `--index` | int | 否 | 分页起始，默认 `0`，取值 `>=0` |
| count | `--count` | int | 否 | 分页条数，默认 `30`，取值 `1~100` |

## 固定查询参数

该命令会固定追加以下后端参数（不暴露为 CLI 参数）：

| 参数 | 固定值 | 说明 |
|---|---|---|
| type | `202` | 目标端 MySQL 类型编码 |
| customs | `clusterType:202` | 固定筛选条件 |

## 返回参数

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 请求状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

### `responseData` 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 结果总数 |
| data | object[] | 目标端 MySQL 实例列表 |

### `responseData.data[]` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| systemId | string | 系统 ID |
| resourceId | string | 资源 ID |
| name | string | 实例名称 |
| businessName | string | 业务名称 |
| resourceStatus | int32 | 资源状态 |
| authStatus | int32 | 授权状态 |
| osUserName | string | OS 用户名 |
| frozenState | string | 冻结状态 |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientId | string | 客户端 ID |
| clientType | int32 | 客户端类型 |
| clientOSType | int32 | 客户端操作系统类型 |
| clientStatus | int32 | 客户端状态 |
| assignStatus | int32 | 分配状态 |
| clusterType | int32 | 集群类型 |
| syncTime | int32 | 同步周期 |
| lastSuccessSyncTime | int64 | 最近成功同步时间 |
| subObjectGroupNumber | int32 | 子对象组数量 |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "totalNum": 1,
    "data": [
      {
        "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
        "resourceId": "428e4c05bf87deef20a86e670f3eb9a9",
        "name": "mysql3306_U0HYTDM3RENXS4EJ",
        "businessName": "mysql3306_U0HYTDM3RENXS4EJ",
        "resourceStatus": 1,
        "authStatus": 1,
        "osUserName": "mysql",
        "frozenState": "1",
        "clientName": "iv-yej4xkshkwwh2yram2x1",
        "clientIp": "172.31.0.2",
        "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
        "clientType": 1,
        "clientOSType": 2,
        "clientStatus": 1,
        "assignStatus": 1,
        "clusterType": 0,
        "syncTime": 300,
        "lastSuccessSyncTime": 1775885977280,
        "subObjectGroupNumber": 0
      }
    ]
  }
}
```

## 示例

```bash
foundation-cli mysql target-instance list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --index 0 \
  --count 30
```
