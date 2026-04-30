---
name: foundation-cli-vmware
description: 当用户需要通过 `foundation-cli vmware` 执行 VMware 对象、数据源、备份配置、恢复配置、任务详情或时间点元数据相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli vmware --help"
---

# Foundation CLI VMware 技能

当任务属于 `foundation-cli vmware` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `vmware` 域。
2. 选择正确的 `vmware` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要查询 VMware 对象列表或对象详情。
- 用户要查询 VMware 数据源。
- 用户要查看备份配置详情、备份任务详情或恢复任务详情。
- 用户要查询时间点元数据。
- 用户要创建 VMware 备份配置或恢复配置。

## 不适用场景

- 用户已经明确要查询通用异步作业详情或执行输出，此时应优先切换到 `job` 域技能。
- 用户需要的是 Host、MySQL 等其他应用域能力，而不是 VMware 域能力。
- 用户只是要走通用 API 透传，且没有标准 `vmware` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“对象列表”，使用 `vmware object list`。
- 如果用户要查“对象详情”，使用 `vmware object info`。
- 如果用户要查“数据源”，使用 `vmware datasource get`。
- 如果用户要查“备份配置详情”，使用 `vmware backup-config detail`。
- 如果用户要“创建备份配置”，使用 `vmware backup-config create`。
- 如果用户要查“备份任务详情”，使用 `vmware backup-detail`。
- 如果用户要查“恢复任务详情”，使用 `vmware recovery-detail`。
- 如果用户要“创建恢复配置”，使用 `vmware restore-config create`。
- 如果用户要查“时间点元数据”，使用 `vmware timepoint metadata`。
- `object info` 与 `backup-config detail` 都依赖 `--object-id`。
- `backup-detail` 与 `recovery-detail` 都依赖 `--task-id`。

## 写入与只读

写入命令：

- `vmware backup-config create`
- `vmware restore-config create`

只读命令：

- `vmware object list`
- `vmware object info`
- `vmware datasource get`
- `vmware backup-config detail`
- `vmware backup-detail`
- `vmware recovery-detail`
- `vmware timepoint metadata`

对于写入命令，优先使用 `--data` 传递结构化载荷。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `object list` 和 `object info` 混淆。
- 不要把 `backup-config detail` 和 `backup-config create` 混淆。
- 不要把 `backup-detail` 和 `recovery-detail` 混淆。
- 不要在缺少 `--object-id`、`--task-id`、`--production-system-id` 的情况下直接拼命令。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询 VMware 对象列表 | `foundation-cli vmware object list --tenant-id <tenant-id> --production-system-id <id>` | [object-list.md](./references/commands/object-list.md) |
| 查询 VMware 对象详情 | `foundation-cli vmware object info --tenant-id <tenant-id> --object-id <object-id>` | [object-info.md](./references/commands/object-info.md) |
| 查询 VMware 数据源 | `foundation-cli vmware datasource get --tenant-id <tenant-id> --production-system-id <id>` | [datasource-get.md](./references/commands/datasource-get.md) |
| 查询 VMware 备份配置详情 | `foundation-cli vmware backup-config detail --tenant-id <tenant-id> --object-id <object-id>` | [backup-config-detail.md](./references/commands/backup-config-detail.md) |
| 查询 VMware 备份任务详情 | `foundation-cli vmware backup-detail --tenant-id <tenant-id> --task-id <task-id>` | [backup-detail.md](./references/commands/backup-detail.md) |
| 查询 VMware 恢复任务详情 | `foundation-cli vmware recovery-detail --tenant-id <tenant-id> --task-id <task-id>` | [recovery-detail.md](./references/commands/recovery-detail.md) |
| 查询 VMware 时间点元数据 | `foundation-cli vmware timepoint metadata --tenant-id <tenant-id> --timestamp <ts> --data-set-id <id> [--business 1]` | [timepoint-metadata.md](./references/commands/timepoint-metadata.md) |
| 创建 VMware 备份配置 | `foundation-cli vmware backup-config create --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [backup-config-create.md](./references/commands/backup-config-create.md) |
| 创建 VMware 恢复配置 | `foundation-cli vmware restore-config create --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [restore-config-create.md](./references/commands/restore-config-create.md) |

## 参考资料

- [命令映射](./references/commands.md)
