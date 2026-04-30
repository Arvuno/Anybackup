# Host 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli host` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `host object create`、`host backup-config`、`host restore start` |
| `只读` | `host list`、`host object list`、`host object detail`、`host backup-config detail` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `host/list` | 查询主机列表 | `foundation-cli host list --tenant-id <tenant-id> [filters]` | 支持数组和布尔筛选 | [list.md](./commands/list.md) |
| `host/object/create` | 创建主机保护对象 | `foundation-cli host object create --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [object-create.md](./commands/object-create.md) |
| `host/object/list` | 查询主机保护对象列表 | `foundation-cli host object list --tenant-id <tenant-id> [filters]` | 默认分页 `index=0,count=10` | [object-list.md](./commands/object-list.md) |
| `host/object/detail` | 查询主机保护对象详情 | `foundation-cli host object detail --tenant-id <tenant-id> --object-id <object-id>` | 必填 `--object-id` | [object-detail.md](./commands/object-detail.md) |
| `host/backup-config` | 创建主机备份配置 | `foundation-cli host backup-config --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [backup-config.md](./commands/backup-config.md) |
| `host/backup-config/detail` | 查询主机备份配置详情 | `foundation-cli host backup-config detail --tenant-id <tenant-id> --object-id <object-id>` | 必填 `--object-id` | [backup-config-detail.md](./commands/backup-config-detail.md) |
| `host/restore/start` | 发起主机恢复 | `foundation-cli host restore start --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [restore-start.md](./commands/restore-start.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `host object detail` | `--object-id` | 单对象详情依赖对象 ID |
| `host backup-config detail` | `--object-id` | 配置详情依赖对象 ID |
| `host backup-config` | `--data` | 写命令默认走结构化请求体 |
| `host restore start` | `--data` | 恢复任务需要请求体 |
