# `foundation-cli client datasource list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client datasource list --tenant-id <tenant-id> --client-id <id> [--full-path <path>] [--business-type <n>] [--runner-type <type>] [--runner-user <user>] [--request-id <id>] [--index <n>] [--count <n>]` |
| 方法 | `GET` |
| 路径 | `/commons/v1/datasources` |
| 风险 | `只读` |

## 何时使用

当用户要查看某个客户端下的数据源树、浏览目录，或在 Runner 上下文中继续下钻资源时，使用本命令。

## 必填上下文

- `--endpoint`
- `--ak`
- `--sk`
- `--client-id`

## 参数要点

| 参数 | 是否必填 | 说明 |
|---|---|---|
| `--client-id` | 是 | 客户端 ID |
| `--full-path` | 否 | 数据源路径；为空时按后端默认根路径查询 |
| `--business-type` | 否 | 业务类型过滤；枚举为 `1=恢复场景`、`2=备份场景`、`3=公共默认`；MySQL 恢复数据源常用 `1` |
| `--runner-type` | 否 | Runner 类型，如 `MySQL`；建议先用 [runner-types.md](./runner-types.md) 确认可用值 |
| `--runner-user` | 否 | Runner 用户；MySQL 恢复路径查询固定使用 `root` |
| `--request-id` | 否 | 分段查询时透传 |
| `--index` | 否 | 默认 `0` |
| `--count` | 否 | 默认 `100` |

## 最小可执行示例

```bash
foundation-cli client datasource list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --client-id <client-id> \
  --business-type 1 \
  --runner-type MySQL \
  --runner-user root \
  --index 0 \
  --count 100
```

MySQL 恢复场景必须使用 `--runner-type MySQL --runner-user root`，不要把目标实例返回的 `osUserName` 当作 `--runner-user`。
如果完整下钻后仍获取不到有效路径，MySQL 恢复场景默认恢复路径为 `/var/lib/mysql`。

## 常见错误

| 错误 | 原因 | 处理方式 |
|---|---|---|
| `missing --client-id` | 缺少客户端上下文 | 先用 [list.md](./list.md) 获取 `clientId` |
| 返回为空 | 路径、业务类型或 Runner 条件不匹配 | 放宽筛选条件，先从 `/` 开始；MySQL 恢复场景完整下钻仍为空时默认 `/var/lib/mysql` |

## 相关跳转

- 先获取客户端 ID： [list.md](./list.md)
- 字段来源：`params/client/datasource_list.go`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ datasource\ list |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回结果

### 顶层字段

| 字段 | 类型 | 说明 |
|---|---|---|
| error | any | 错误信息，成功通常为 `null` |
| status | string | 请求状态 |
| responseData | object | 数据源查询结果 |

### responseData 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| data | object[] | 当前层级节点列表 |
| partialSign | int | 分段标记 |
| more | bool | 是否还有更多数据 |
| requestId | string | 分段/下钻请求标识 |
| finished | bool | 本轮查询是否结束 |

### responseData.data[] 字段

| 字段 | 类型 | 说明 |
|---|---|---|
| fullPath | string | 节点完整路径 |
| disPath | string | 节点展示路径 |
| expandedFlag | bool | 是否可展开 |
| checkable | bool | 是否可选中 |

## 返回案例

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "data": [
      {
        "fullPath": "/",
        "disPath": "/",
        "expandedFlag": true,
        "checkable": true
      }
    ],
    "partialSign": 0,
    "more": false,
    "requestId": "39",
    "finished": true
  }
}
```

## 示例

```bash
foundation-cli client datasource list
```
