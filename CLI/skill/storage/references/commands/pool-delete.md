# `foundation-cli storage pool delete`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli storage pool delete --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --pool-id <pool-id>` |
| 方法 | `DELETE` |
| 路径 | `/storageresmgm/v1/{storageSvcId}/pool?id={poolId}` |
| 风险 | `写入` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--pool-id` | 是 | 存储池 ID |

## 参数映射

| CLI Flag | 位置 | 后端参数 |
|---|---|---|
| `--storage-svc-id` | 路径 | `{storageSvcId}` |
| `--pool-id` | Query | `id` |

## 示例

```bash
foundation-cli storage pool delete \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 52247d9c121745eab3a7880817457df5 \
  --pool-id e177c3fe311f11f1a1220050568943e2
```

## 使用建议

1. 先执行 `storage pool list`，确认 `pool-id`。
2. 该命令没有请求体，所有业务定位都靠路径和 query 参数完成。
3. 删除属于写操作，适合在 agent 工作流里明确增加确认步骤。
4. 该接口签名遵循 SDK demo 的 forward 规则：请求路径保留 `{storageSvcId}`，签名 path 使用 `/storageresmgm/v1/pool`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## Body 参数

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| data | object | 是 | 请求体对象，建议按“请求体示例”中的字段组织 |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
