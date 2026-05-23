---
name: foundation-cli-protect
description: 当用户需要通过 `foundation-cli protect` 执行跨应用保护动作、策略绑定、备份启动、配置策略查询、客户端关联查询或存储池自动选择时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli protect --help"
---

# Foundation CLI Protect 技能

当任务属于 `foundation-cli protect` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `protect` 域。
2. 选择正确的 `protect` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要绑定、批量绑定、查询绑定或批量解绑策略。
- 用户要对一个或多个对象发起备份。
- 用户要查询对象级或应用类型级配置策略。
- 用户要查询备份配置。
- 用户要自动选择存储池。
- 用户要查询某个保护对象关联的客户端列表。

## 不适用场景

- 用户已经明确要管理策略模板本身，此时应优先切换到 `policy` 域技能。
- 用户需要的是 MySQL、VMware、Host 等应用域专属对象能力，而不是跨应用保护动作。
- 用户只是要走通用 API 透传，且没有标准 `protect` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要“绑定一个对象的策略”，使用 `protect policy bind`。
- 如果用户要“基于单个保护对象创建备份策略（命令会自动绑定并在失败时回滚删除）”，使用 `protect policy create-bind`。
- 如果用户要“批量绑定”，使用 `protect policy bind-batch`。
- 如果用户要“查询绑定列表”，使用 `protect policy bind-list`。
- 如果用户要“批量解绑”，使用 `protect policy batch-unbind`。
- 如果用户要“启动单对象备份”，使用 `protect backup start`。
- 如果用户要“启动批量备份”，使用 `protect backup batch-start`。
- 如果用户要“查询对象级配置策略”，使用 `protect config-policy get`。
- 如果用户要“查询应用类型级配置策略”，使用 `protect config-policy get-by-app-type`。
- 如果用户要“查询关联客户端”，使用 `protect client list`。
- 多个对象级命令依赖 `--object-id`。

## 写入与只读

写入命令：

- `protect policy bind`
- `protect policy create-bind`
- `protect policy bind-batch`
- `protect policy batch-unbind`
- `protect backup start`
- `protect backup batch-start`

只读命令：

- `protect policy bind-list`
- `protect config-policy get`
- `protect config-policy get-by-app-type`
- `protect storage-pool auto-select`
- `protect backup-config get`
- `protect client list`

对于写入命令，优先使用 `--data` 传递结构化载荷。

当前环境实测：
- `protect policy bind-list` 成功，返回 `totalNum=0`。
- `protect config-policy get` 成功，返回 `appType=eso_backupengine_hypermysqlengine`。
- `protect storage-pool auto-select` 成功，返回存储池 `block`。
- `protect backup-config get` 成功，返回 `isConfig=true`。
- `protect client list` 成功，返回 1 个客户端。
- `protect config-policy get-by-app-type` 当前返回 `Cli.UserMissingOrDisabled` / HTTP `403`；同一路径通过 `foundation-cli api` 透传也返回一致结果。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `policy` 域里的策略模板管理和 `protect` 域里的策略绑定动作混淆。
- 不要把单对象操作和批量操作混淆。
- 不要在没有 `--object-id` 的情况下调用对象级命令。
- 不要在没有 `--data` 的情况下调用写命令。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 绑定单个对象的策略 | `foundation-cli protect policy bind --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [policy-bind.md](./references/commands/policy-bind.md) |
| 基于单个保护对象创建备份策略（命令会自动绑定并在失败时回滚删除） | `foundation-cli protect policy create-bind --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [policy-create-bind.md](./references/commands/policy-create-bind.md) |
| 批量绑定策略 | `foundation-cli protect policy bind-batch --tenant-id <tenant-id> --data '<json>'` | [policy-bind-batch.md](./references/commands/policy-bind-batch.md) |
| 查询对象绑定策略列表 | `foundation-cli protect policy bind-list --tenant-id <tenant-id> --object-id <object-id>` | [policy-bind-list.md](./references/commands/policy-bind-list.md) |
| 批量解绑策略 | `foundation-cli protect policy batch-unbind --tenant-id <tenant-id> --data '<json>'` | [policy-batch-unbind.md](./references/commands/policy-batch-unbind.md) |
| 启动单对象备份 | `foundation-cli protect backup start --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [backup-start.md](./references/commands/backup-start.md) |
| 启动批量备份 | `foundation-cli protect backup batch-start --tenant-id <tenant-id> --data '<json>'` | [backup-batch-start.md](./references/commands/backup-batch-start.md) |
| 查询对象级配置策略 | `foundation-cli protect config-policy get --tenant-id <tenant-id> --object-id <object-id>` | [config-policy-get.md](./references/commands/config-policy-get.md) |
| 查询应用类型级配置策略 | `foundation-cli protect config-policy get-by-app-type --tenant-id <tenant-id> --app-type <app-type>` | [config-policy-get-by-app-type.md](./references/commands/config-policy-get-by-app-type.md) |
| 自动选择存储池 | `foundation-cli protect storage-pool auto-select --tenant-id <tenant-id> [--exclude-id <id> ...]` | [storage-pool-auto-select.md](./references/commands/storage-pool-auto-select.md) |
| 查询备份配置 | `foundation-cli protect backup-config get --tenant-id <tenant-id> --object-id <object-id>` | [backup-config-get.md](./references/commands/backup-config-get.md) |
| 查询关联客户端 | `foundation-cli protect client list --tenant-id <tenant-id> --object-id <object-id>` | [client-list.md](./references/commands/client-list.md) |

## 参考资料

- [命令映射](./references/commands.md)
