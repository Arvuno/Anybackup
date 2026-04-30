# Job 命令映射

## 先区分只读/写入
| 类型 | 命令 |
|---|---|
| `只读` | `job list`、`job backup-detail`、`job sync-detail`、`job child list`、`job logs`、`job business-types`、`job app-types` |
| `写入` | `job stop`、`job delete` |

## 意图到命令
| 规范路径 | 常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `job/list` | 查询作业列表 | `foundation-cli job list --tenant-id <tenant-id> [filters]` | 支持分页和筛选 | [list.md](./commands/list.md) |
| `job/backup-detail` | 查询备份作业详情 | `foundation-cli job backup-detail --job-id <job-id>` | 必填 `--job-id` | [backup-detail.md](./commands/backup-detail.md) |
| `job/sync-detail` | 查询复制作业详情 | `foundation-cli job sync-detail --job-id <job-id>` | 必填 `--job-id` | [sync-detail.md](./commands/sync-detail.md) |
| `job/child/list` | 查询子作业列表 | `foundation-cli job child list --job-id <job-id> [--index <n>] [--count <n>]` | 必填 `--job-id`；默认 `--index 0 --count 30` | [child-list.md](./commands/child-list.md) |
| `job/logs` | 查询作业执行输出 | `foundation-cli job logs --job-id <job-id> [filters]` | 必填 `--job-id` | [logs.md](./commands/logs.md) |
| `job/business-types` | 查询业务类型字典 | `foundation-cli job business-types` | 无请求体 | [business-types.md](./commands/business-types.md) |
| `job/app-types` | 查询应用类型字典 | `foundation-cli job app-types` | 无请求体 | [app-types.md](./commands/app-types.md) |
| `job/stop` | 终止作业 | `foundation-cli job stop (--job-id <id> ... \| --data <json> \| --body-file <file>)` | `--job-id` 与 `--data/--body-file` 互斥 | [stop.md](./commands/stop.md) |
| `job/delete` | 删除作业 | `foundation-cli job delete (--job-id <id> ... \| --data <json> \| --body-file <file>)` | `--job-id` 与 `--data/--body-file` 互斥 | [delete.md](./commands/delete.md) |
