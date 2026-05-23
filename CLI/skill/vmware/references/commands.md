# VMware 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli vmware` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `vmware backup-config create`、`vmware restore-config create` |
| `只读` | `vmware object list`、`vmware object info`、`vmware datasource get`、`vmware backup-config detail`、`vmware backup-detail`、`vmware recovery-detail`、`vmware timepoint metadata` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `vmware/object/list` | 查询 VMware 对象列表 | `foundation-cli vmware object list --tenant-id <tenant-id> --production-system-id <id>` | 必填 `--production-system-id` | [object-list.md](./commands/object-list.md) |
| `vmware/object/info` | 查询 VMware 对象详情 | `foundation-cli vmware object info --tenant-id <tenant-id> --object-id <object-id>` | 必填 `--object-id` | [object-info.md](./commands/object-info.md) |
| `vmware/datasource/get` | 查询 VMware 数据源 | `foundation-cli vmware datasource get --tenant-id <tenant-id> --production-system-id <id>` | 必填 `--production-system-id` | [datasource-get.md](./commands/datasource-get.md) |
| `vmware/backup-config/detail` | 查询 VMware 备份配置详情 | `foundation-cli vmware backup-config detail --tenant-id <tenant-id> --object-id <object-id>` | 必填 `--object-id` | [backup-config-detail.md](./commands/backup-config-detail.md) |
| `vmware/backup-detail` | 查询 VMware 备份任务详情 | `foundation-cli vmware backup-detail --tenant-id <tenant-id> --task-id <task-id>` | 必填 `--task-id` | [backup-detail.md](./commands/backup-detail.md) |
| `vmware/recovery-detail` | 查询 VMware 恢复任务详情 | `foundation-cli vmware recovery-detail --tenant-id <tenant-id> --task-id <task-id>` | 必填 `--task-id` | [recovery-detail.md](./commands/recovery-detail.md) |
| `vmware/timepoint/metadata` | 查询 VMware 时间点元数据 | `foundation-cli vmware timepoint metadata --tenant-id <tenant-id> --timestamp <ts> --data-set-id <id>` | 必填 `--timestamp` 和 `--data-set-id` | [timepoint-metadata.md](./commands/timepoint-metadata.md) |
| `vmware/backup-config/create` | 创建 VMware 备份配置 | `foundation-cli vmware backup-config create --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | 必填 `--object-id` 和 `--data` | [backup-config-create.md](./commands/backup-config-create.md) |
| `vmware/restore-config/create` | 创建 VMware 恢复配置 | `foundation-cli vmware restore-config create --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | 必填 `--object-id` 和 `--data` | [restore-config-create.md](./commands/restore-config-create.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `vmware object list` | `--production-system-id` | 对象列表依赖生产系统上下文 |
| `vmware object info` | `--object-id` | 单对象详情依赖对象 ID |
| `vmware backup-detail` / `vmware recovery-detail` | `--task-id` | 任务详情依赖任务 ID |
| `vmware timepoint metadata` | `--timestamp`、`--data-set-id` | 元数据查询依赖时间点上下文 |
