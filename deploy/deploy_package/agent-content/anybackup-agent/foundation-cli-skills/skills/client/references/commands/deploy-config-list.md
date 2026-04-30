# `foundation-cli client deploy-config list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client deploy-config list --tenant-id <tenant-id>` |
| Method | `GET` |
| Path | `/deploy/v1/hostConfig` |
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

无。

## 返回结果

依据 `params/client/host_config.go` 中 `HostResponse`、`HostResponseData`、`Host`。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息 |
| status | string | 请求状态 |
| responseData | object | 业务数据 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int | 总数 |
| data | Host[] | 部署配置列表 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| hostId | string | 主机配置 ID |
| name | string | 配置名称 |
| user | string | 用户名 |
| port | int | 端口 |
| host | string | 主机地址 |
| hostType | int | 主机类型（整型编码） |
| administrator | int | 管理员模式标记（整型编码） |
| userType | int | 用户类型（整型编码） |
| status | int | 状态码 |
| errorMsg | string | 错误信息 |

## 枚举类型列表

`params/client/host_config.go` 未提供 `hostType`、`administrator`、`userType`、`status` 的显式枚举常量定义，当前按整数透传展示。

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
        "hostId": "sample-host-config-id",
        "name": "linux-config",
        "user": "root",
        "port": 22,
        "host": "10.10.10.10",
        "hostType": 1,
        "administrator": 1,
        "userType": 1,
        "status": 1,
        "errorMsg": ""
      }
    ]
  }
}
```

## 示例命令
```bash
foundation-cli client deploy-config list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ deploy-config\ list |
| 风险 | read-only |

## 示例

```bash
foundation-cli client deploy-config list
```
