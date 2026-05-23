# Protect 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli protect` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `protect policy create-bind`、`protect policy bind`、`protect policy bind-batch`、`protect policy batch-unbind`、`protect backup start`、`protect backup batch-start` |
| `只读` | `protect policy bind-list`、`protect config-policy get`、`protect config-policy get-by-app-type`、`protect storage-pool auto-select`、`protect backup-config get`、`protect client list` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `protect/policy/bind` | 绑定单个对象的策略 | `foundation-cli protect policy bind --object-id <id> --data '<json>'` | 必填 `--object-id` 和 `--data` | [policy-bind.md](./commands/policy-bind.md) |
| `protect/policy/create-bind` | 创建备份策略并绑定到单个对象（失败自动回滚删除） | `foundation-cli protect policy create-bind --object-id <id> --data '<json>'` | 必填 `--object-id` 和 `--data`；`bind` 失败会尝试自动删除刚创建策略 | [policy-create-bind.md](./commands/policy-create-bind.md) |
| `protect/policy/bind-batch` | 批量绑定策略 | `foundation-cli protect policy bind-batch --data '<json>'` | 必填 `--data` | [policy-bind-batch.md](./commands/policy-bind-batch.md) |
| `protect/policy/bind-list` | 查询对象绑定策略列表 | `foundation-cli protect policy bind-list --object-id <id> [filters]` | 必填 `--object-id`；当前环境已验证成功 | [policy-bind-list.md](./commands/policy-bind-list.md) |
| `protect/policy/batch-unbind` | 批量解绑策略 | `foundation-cli protect policy batch-unbind --data '<json>'` | 必填 `--data` | [policy-batch-unbind.md](./commands/policy-batch-unbind.md) |
| `protect/backup/start` | 启动单对象备份 | `foundation-cli protect backup start --object-id <id> --data '<json>'` | 必填 `--object-id` 和 `--data` | [backup-start.md](./commands/backup-start.md) |
| `protect/backup/batch-start` | 启动批量备份 | `foundation-cli protect backup batch-start --data '<json>'` | 必填 `--data` | [backup-batch-start.md](./commands/backup-batch-start.md) |
| `protect/config-policy/get` | 查询对象级配置策略 | `foundation-cli protect config-policy get --object-id <id>` | 必填 `--object-id`；当前环境已验证成功 | [config-policy-get.md](./commands/config-policy-get.md) |
| `protect/config-policy/get-by-app-type` | 查询应用类型级配置策略 | `foundation-cli protect config-policy get-by-app-type --app-type <app-type>` | 必填 `--app-type`；当前环境返回 `Cli.UserMissingOrDisabled` / HTTP `403` | [config-policy-get-by-app-type.md](./commands/config-policy-get-by-app-type.md) |
| `protect/storage-pool/auto-select` | 自动选择存储池 | `foundation-cli protect storage-pool auto-select [--exclude-id <id> ...]` | 可传排除池 ID；当前环境已验证成功 | [storage-pool-auto-select.md](./commands/storage-pool-auto-select.md) |
| `protect/backup-config/get` | 查询备份配置 | `foundation-cli protect backup-config get --object-id <id>` | 必填 `--object-id`；当前环境已验证成功 | [backup-config-get.md](./commands/backup-config-get.md) |
| `protect/client/list` | 查询关联客户端 | `foundation-cli protect client list --object-id <id>` | 必填 `--object-id`；当前环境已验证成功 | [client-list.md](./commands/client-list.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `protect policy bind` | `--object-id`、`--data` | 单对象绑定依赖对象和请求体 |
| `protect policy create-bind` | `--object-id`、`--data` | 复合写命令：创建策略并绑定，绑定失败自动回滚删除 |
| `protect backup start` | `--object-id`、`--data` | 单对象备份依赖对象和请求体 |
| `protect config-policy get` | `--object-id` | 对象级配置策略依赖对象上下文 |
| `protect client list` | `--object-id` | 客户端关联查询依赖对象上下文 |
