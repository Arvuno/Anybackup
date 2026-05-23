---
name: foundation-cli-network
description: 当用户需要通过 `foundation-cli network` 执行子网、子网节点、节点业务 IP 查询，或设置、移除节点业务 IP 时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli network --help"
---

# Foundation CLI Network 技能

当任务属于 `foundation-cli network` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `network` 域。
2. 选择正确的 `network` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要获取子网列表。
- 用户要获取指定子网下的节点列表。
- 用户要获取指定节点的业务 IP 列表。
- 用户要为节点设置业务 IP。
- 用户要移除节点业务 IP。

## 不适用场景

- 用户要管理存储池、存储服务或建池节点，此时应优先切换到 `storage` 域技能。
- 用户要执行其他业务域对象操作，而不是子网、节点或节点业务 IP 相关操作。
- 用户只是要走通用 API 透传，且没有标准 `network` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“子网列表”，使用 `network subnet list`。
- 如果用户要查“某个子网下的节点”，使用 `network subnet node list`。
- 如果用户要查“某个节点的业务 IP”，使用 `network subnet node ip list`。
- 如果用户要“新增/设置”节点业务 IP，使用 `network subnet node ip set`。
- 如果用户要“移除”节点业务 IP，使用 `network subnet node ip remove`。
- 所有 `subnet node*` 命令都依赖 `--subnet-id`。
- 所有 `subnet node ip*` 命令都依赖 `--node-id`。
- `--plane-type` 是枚举类型：`1` 管理平面，`2` 存储平面，`3` 备份平面，`4` 复制平面。
- `set` 和 `remove` 都要求 `--ip` 是标准 IPv4 或 IPv6，不接受 CIDR。

## 写入与只读

写入命令：

- `network subnet node ip set`
- `network subnet node ip remove`

只读命令：

- `network subnet list`
- `network subnet node list`
- `network subnet node ip list`

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`
- `--storage-svc-id`

## 常见误判

- 不要把 `storage` 域和 `network` 域混淆。
- 不要在没有 `--subnet-id` 的情况下调用子网节点相关命令。
- 不要在没有 `--node-id` 的情况下调用节点业务 IP 相关命令。
- 不要把 `ip list` 和 `ip set/remove` 混淆。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 获取子网列表 | `foundation-cli network subnet list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> [--plane-type <1|2|3|4>]` | [subnet-list.md](./references/commands/subnet-list.md) |
| 获取子网节点列表 | `foundation-cli network subnet node list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id>` | [subnet-node-list.md](./references/commands/subnet-node-list.md) |
| 获取节点业务 IP 列表 | `foundation-cli network subnet node ip list --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id>` | [subnet-node-ip-list.md](./references/commands/subnet-node-ip-list.md) |
| 设置节点业务 IP | `foundation-cli network subnet node ip set --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id> --ip <ip>` | [subnet-node-ip-set.md](./references/commands/subnet-node-ip-set.md) |
| 移除节点业务 IP | `foundation-cli network subnet node ip remove --tenant-id <tenant-id> --storage-svc-id <storage-svc-id> --subnet-id <subnet-id> --node-id <node-id> --ip <ip>` | [subnet-node-ip-remove.md](./references/commands/subnet-node-ip-remove.md) |

## 参考资料

- [命令映射](./references/commands.md)
