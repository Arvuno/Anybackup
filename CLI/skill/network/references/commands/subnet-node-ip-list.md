# `foundation-cli network subnet node ip list`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli network subnet node ip list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id>` |
| 方法 | `GET` |
| 路径 | `/clusters/v1/{storageSvcId}/subnet/nodes/node_ip_addresses` |
| 风险 | `只读` |

## 必填输入

| 参数 | 必填 | 说明 |
|---|---|---|
| `--tenant-id` | 否 | 租户上下文 |
| `--endpoint` | 是 | Foundation 接口地址 |
| `--ak` | 是 | 访问密钥 |
| `--sk` | 是 | 访问密钥密文 |
| `--storage-svc-id` | 是 | 存储服务 ID |
| `--subnet-id` | 是 | 子网 ID |
| `--node-id` | 是 | 节点 ID |

## 可选输入

| 参数 | 默认值 | 说明 |
|---|---|---|
| `--index` | `0` | 分页起始位置 |
| `--count` | `30` | 分页数量 |

## 查询参数映射

| CLI Flag | 后端参数 |
|---|---|
| `--index` | `index` |
| `--count` | `count` |
| `--node-id` | `nodeId` |
| `--subnet-id` | `subnetId` |

## 返回结果关注点

常见字段：

| 字段 | 说明 |
|---|---|
| `responseData.data[].ip` | 节点业务 IP |
| `responseData.data[].cardName` | 网卡名称 |
| `responseData.data[].usable` | 该 IP 是否可用于后续设置 |

## 示例

```bash
foundation-cli network subnet node ip list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --storage-svc-id 1234567890abcdef1234567890abcdef \
  --subnet-id fab4ec3015df11f1966b0050568952bd \
  --node-id 839b4b217e214c348152221288257773 \
  --index 0 \
  --count 30
```

## 使用建议

1. 在执行 `network subnet node ip set` 之前，先用该命令确认目标 IP 是否存在且 `usable=true`。
2. 如果一个节点返回多条 IP，优先结合 `cardName` 和 `usable` 做筛选。
3. 该命令通常依赖上一步 `network subnet node list` 返回的 `nodeId`。

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```
