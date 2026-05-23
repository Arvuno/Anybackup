# Policy 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli policy` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `policy backup create`、`policy copy create`、`policy delete` |
| `只读` | `policy list`、`policy backup detail`、`policy copy detail` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `policy/list` | 查询策略列表 | `foundation-cli policy list --tenant-id <tenant-id> --validate-status <1|2|3|4|5> [filters]` | 当前环境下需显式传 `--validate-status` | [list.md](./commands/list.md) |
| `policy/backup/detail` | 查看备份策略详情 | `foundation-cli policy backup detail --tenant-id <tenant-id> --group-id <group-id>` | 必填 `--group-id` | [backup-detail.md](./commands/backup-detail.md) |
| `policy/backup/create` | 创建备份策略 | `foundation-cli policy backup create --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [backup-create.md](./commands/backup-create.md) |
| `policy/copy/detail` | 查看复制策略详情 | `foundation-cli policy copy detail --tenant-id <tenant-id> --group-id <group-id>` | 必填 `--group-id` | [copy-detail.md](./commands/copy-detail.md) |
| `policy/copy/create` | 创建复制策略 | `foundation-cli policy copy create --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [copy-create.md](./commands/copy-create.md) |
| `policy/delete` | 删除策略 | `foundation-cli policy delete --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [delete.md](./commands/delete.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `policy backup detail` | `--group-id` | 详情查询依赖策略组 ID |
| `policy copy detail` | `--group-id` | 详情查询依赖策略组 ID |
| `policy backup create` | `--data` | 写命令默认走结构化请求体 |
| `policy delete` | `--data` | 删除通常是批量动作 |
