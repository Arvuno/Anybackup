# `foundation-cli client runner list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client runner list --tenant-id <tenant-id> [--index <n>] [--count <n>]` |
| Method | `GET` |
| Path | `/commons/clientRunner` |
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
| index | `--index` | int | 否 | 分页起始位置 |
| count | `--count` | int | 否 | 分页数量 |

## 返回结果

依据 `params/client/client_runner.go` 中 `APIResponse`、`RunnerResponseData`、`ClientRunner`。

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | string/null | 错误信息 |
| status | string | 请求状态 |
| responseData | object | 业务数据 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int | 总数 |
| data | ClientRunner[] | runner 列表 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| clientRunnerId | string | runner ID |
| name | string | runner 名称 |
| user | string | 用户 |
| status | int | 状态码 |
| appSupport | string[] | 支持应用列表 |
| version | string | runner 版本 |
| productVersion | string | 产品版本 |
| clientIp | string | 客户端 IP |
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| osType | int | 操作系统类型，见枚举 |
| updateStatus | int | 升级状态，见枚举 |
| lastestVersion | string | 最新版本 |
| jobId | string | 作业 ID |
| machineCode | string | 机器码 |
| runnerTypeName | string | runner 类型中文名 |
| runnerTypeNameEn | string | runner 类型英文名 |
| basicIsLatest | bool/null | 基础 runner 是否最新 |

## 枚举类型列表

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

### updateStatus（来自 `params/client/client_runner.go` 注释）

| 值 | 含义 |
|---|---|
| 0 | 全部/未知 |
| 1 | 升级中 |
| 2 | 可升级 |
| 3 | 已是最新，无需升级 |

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
        "clientRunnerId": "sample-runner-id",
        "name": "basic-runner-01",
        "user": "root",
        "status": 1,
        "appSupport": [
          "Oracle",
          "MySQL"
        ],
        "version": "8.0.9.0",
        "productVersion": "8.0.9.0",
        "clientIp": "10.10.10.10",
        "clientId": "sample-client-id",
        "clientName": "mysql-node-01",
        "osType": 2,
        "updateStatus": 3,
        "lastestVersion": "8.0.9.0",
        "jobId": "",
        "machineCode": "sample-machine-code",
        "runnerTypeName": "基础 Runner",
        "runnerTypeNameEn": "Basic Runner",
        "basicIsLatest": true
      }
    ]
  }
}
```

## 示例命令
```bash
foundation-cli client runner list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --index 0 \
  --count 30
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ runner\ list |
| 风险 | read-only |

## 示例

```bash
foundation-cli client runner list
```
