# MySQL 命令映射

本文档只负责帮助 agent 选择正确的 `foundation-cli mysql` 命令，并跳转到对应的单命令文档。

## 先区分只读和写入

| 类型 | 命令 |
|---|---|
| `写入` | `mysql backup-plan create`、`mysql backup-config set`、`mysql restore-config create`、`mysql authorize` |
| `只读` | `mysql object list`、`mysql target-instance list`、`mysql object get`、`mysql backup-config detail`、`mysql datasource backup`、`mysql backup-detail`、`mysql recovery range`、`mysql recovery-config detail`、`mysql datasource recovery`、`mysql recovery timepoint list`、`mysql recovery-detail` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `mysql/object/list` | 查询 MySQL 对象列表 | `foundation-cli mysql object list --tenant-id <tenant-id> --app-type 202 [filters]` | MySQL 场景默认应使用 `--app-type 202` | [object-list.md](./commands/object-list.md) |
| `mysql/target-instance/list` | 查询目标端 MySQL 实例列表 | `foundation-cli mysql target-instance list --tenant-id <tenant-id> [--index <n>] [--count <n>]` | 默认分页 | [target-instance-list.md](./commands/target-instance-list.md) |
| `mysql/object/get` | 查询 MySQL 对象详情 | `foundation-cli mysql object get --tenant-id <tenant-id> --object-id <object-id>` | 必填 `--object-id` | [object-get.md](./commands/object-get.md) |
| `mysql/backup-config/detail` | 查询备份配置详情 | `foundation-cli mysql backup-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` | 必填 `--system-id` 和 `--object-id` | [backup-config-detail.md](./commands/backup-config-detail.md) |
| `mysql/backup-config/set` | 设置备份配置 | `foundation-cli mysql backup-config set --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` 或 `--body-file` | [backup-config-set.md](./commands/backup-config-set.md) |
| `mysql/backup-plan/create` | 创建 MySQL 备份方案（聚合） | 按工作流执行：`policy backup create` -> `protect policy bind` -> `mysql backup-config set` | Skill 编排意图（非二进制子命令）；`policyData/bindData/backupConfigData` 三段透传，`slaIds` 由第 1 步返回自动注入 | [backup-plan-create.md](./commands/backup-plan-create.md) |
| `mysql/datasource/backup` | 查询备份数据源 | `foundation-cli mysql datasource backup --tenant-id <tenant-id> --object-id <object-id> [--index <n>] [--count <n>]` | 必填 `--object-id` | [datasource-backup.md](./commands/datasource-backup.md) |
| `mysql/backup-detail` | 查询备份任务详情 | `foundation-cli mysql backup-detail --tenant-id <tenant-id> --task-id <task-id>` | 必填 `--task-id` | [backup-detail.md](./commands/backup-detail.md) |
| `mysql/recovery/range` | 查询恢复范围 | `foundation-cli mysql recovery range --tenant-id <tenant-id> [query flags]` | 依赖恢复查询参数 | [recovery-range.md](./commands/recovery-range.md) |
| `mysql/recovery-config/detail` | 查询恢复配置详情 | `foundation-cli mysql recovery-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` | 必填 `--system-id` 和 `--object-id` | [recovery-config-detail.md](./commands/recovery-config-detail.md) |
| `mysql/datasource/recovery` | 查询恢复时间点数据源 | `foundation-cli mysql datasource recovery --tenant-id <tenant-id> --data-set-id <data-set-id> --storage-service-id <id> --timestamp <ts>` | 必填 `--data-set-id`、`--storage-service-id`、`--timestamp` | [datasource-recovery.md](./commands/datasource-recovery.md) |
| `mysql/recovery/timepoint/list` | 查询恢复时间点 | `foundation-cli mysql recovery timepoint list --tenant-id <tenant-id> [query flags]` | 默认分页 | [recovery-timepoint-list.md](./commands/recovery-timepoint-list.md) |
| `mysql/recovery-detail` | 查询恢复任务详情 | `foundation-cli mysql recovery-detail --tenant-id <tenant-id> --task-id <task-id>` | 必填 `--task-id` | [recovery-detail.md](./commands/recovery-detail.md) |
| `mysql/restore-config/create` | 创建恢复配置 | `foundation-cli mysql restore-config create --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` 或 `--body-file` | [restore-config-create.md](./commands/restore-config-create.md) |
| `mysql/authorize` | 授权 MySQL 实例 | `foundation-cli mysql authorize --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` 或 `--body-file` | [authorize.md](./commands/authorize.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `mysql object list` | `--app-type` | MySQL 对象列表通常依赖 `--app-type 202` |
| `mysql object get` | `--object-id` | 单对象详情依赖对象 ID |
| `mysql backup-config detail` | `--system-id`、`--object-id` | 配置详情依赖系统与对象上下文 |
| `mysql datasource recovery` | `--data-set-id`、`--storage-service-id`、`--timestamp` | 恢复时间点数据源依赖数据集、存储服务和时间点上下文 |
| `mysql backup-detail` / `mysql recovery-detail` | `--task-id` | 任务详情依赖任务 ID |
