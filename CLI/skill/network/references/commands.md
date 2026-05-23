# Network 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli network` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `network subnet node ip set`、`network subnet node ip remove` |
| `只读` | `network subnet list`、`network subnet node list`、`network subnet node ip list` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `network/subnet/list` | 获取子网列表 | `foundation-cli network subnet list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> [--plane-type <1|2|3|4>]` | 必填 `--storage-svc-id`；`--plane-type` 为平面枚举 | [subnet-list.md](./commands/subnet-list.md) |
| `network/subnet/node/list` | 获取子网节点列表 | `foundation-cli network subnet node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id>` | 必填 `--subnet-id` | [subnet-node-list.md](./commands/subnet-node-list.md) |
| `network/subnet/node/ip/list` | 获取节点业务 IP 列表 | `foundation-cli network subnet node ip list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id>` | 必填 `--node-id` | [subnet-node-ip-list.md](./commands/subnet-node-ip-list.md) |
| `network/subnet/node/ip/set` | 设置节点业务 IP | `foundation-cli network subnet node ip set --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id> --ip <ip>` | 必填 `--ip` | [subnet-node-ip-set.md](./commands/subnet-node-ip-set.md) |
| `network/subnet/node/ip/remove` | 移除节点业务 IP | `foundation-cli network subnet node ip remove --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id> --ip <ip>` | 必填 `--ip` | [subnet-node-ip-remove.md](./commands/subnet-node-ip-remove.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `network subnet list` | `--storage-svc-id`、`--plane-type` | 子网查询依赖存储服务上下文；`--plane-type` 仅允许 `1/2/3/4` |
| `network subnet node list` | `--subnet-id` | 节点查询依赖子网上下文 |
| `network subnet node ip list` | `--node-id` | 业务 IP 查询依赖节点上下文 |
| `network subnet node ip set/remove` | `--ip` | 只接受标准 IPv4 或 IPv6 |
