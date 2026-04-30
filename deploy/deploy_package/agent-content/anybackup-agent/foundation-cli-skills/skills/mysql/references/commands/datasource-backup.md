# `foundation-cli mysql datasource backup`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql datasource backup --tenant-id <tenant-id> --object-id <object-id> [--index <n>] [--count <n>] [--system-id <id>] [--full-path <path>] [--request-id <id>]` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/backup_datasource` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| index | `--index` | int | 否 | 默认值为 `0` |
| count | `--count` | int | 否 | 默认值为 `20`，取值范围 `1~100` |
| systemId | `--system-id` | string | 否 | MySQL 生产系统 ID |
| objectId | `--object-id` | string | 是 | MySQL 保护对象 ID |
| fullPath | `--full-path` | string | 否 | 数据源路径筛选 |
| requestId | `--request-id` | string | 否 | 可选请求关联 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 备份数据源响应体 |

### `responseData`

| 字段 | 类型 | 说明 |
|---|---|---|
| data | object[] | 数据源列表 |
| finished | bool | 是否已枚举完成 |
| more | bool | 是否还有更多记录 |
| partialSign | int32 | 部分结果标记 |
| requestId | string | 请求 ID |
| totalNum | int32 | 总记录数 |

### `data[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| authType | string | 授权类型 |
| checkAble | bool | 节点是否可选 |
| checked | bool | 节点是否已选中 |
| code | string | 节点编码 |
| custom | string | 自定义扩展字段 |
| dbName | string | 数据库名称 |
| dispPath | string | 展示路径 |
| fullPath | string | 完整路径 |
| expandedFlag | bool | 节点是否可展开 |
| isAuthorized | bool | 节点是否已授权 |
| nodeType | int32 | 节点类型 |
| userName | string | 用户名 |
| uuid | string | UUID |

## 枚举值

### `nodeType`

| 值 | 含义 |
|---|---|
| `params` 中未固定 | 按集成返回值解释 |

## 校验说明

- 已对照 [datasource_backup.go](/e:/go/src/foundation-cli/internal/business/mysql/datasource_backup.go).
- 已确认 CLI 参数正确映射到 `index`, `count`, `systemId`, `objectId`, `fullPath`, 和 `requestId`。

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
                                          "finished":  true,
                                          "more":  true,
                                          "partialSign":  true
                                      }
                                  ]
                     }
}
```

## 示例
```bash
foundation-cli mysql datasource backup \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --count 30
```



