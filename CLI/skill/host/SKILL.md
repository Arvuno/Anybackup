---
name: foundation-cli-host
description: 当用户需要通过 `foundation-cli host` 执行主机列表、主机保护对象、备份配置或主机恢复相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli host --help"
---

# Foundation CLI Host 技能

当任务属于 `foundation-cli host` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `host` 域。
2. 选择正确的 `host` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要查询主机列表。
- 用户要创建或查询主机保护对象。
- 用户要创建或查看主机备份配置。
- 用户要发起主机恢复任务。

## 不适用场景

- 用户已经明确要查询通用异步作业详情或执行输出，此时应优先切换到 `job` 域技能。
- 用户需要的是 MySQL、VMware 等应用专属对象能力，而不是主机对象能力。
- 用户只是要走通用 API 透传，且没有标准 `host` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“主机列表”，使用 `host list`。
- 如果用户要查“主机保护对象列表”，使用 `host object list`。
- 如果用户要查“单个主机保护对象详情”，使用 `host object detail`。
- 如果用户要“创建主机保护对象”，使用 `host object create`。
- 如果用户要查“备份配置详情”，使用 `host backup-config detail`。
- 如果用户要“设置/提交备份配置”，使用 `host backup-config`。
- 如果用户要“发起恢复”，使用 `host restore start`。
- `host object detail` 与 `host backup-config detail` 都依赖 `--object-id`。

## 写入与只读

写入命令：

- `host object create`
- `host backup-config`
- `host restore start`

只读命令：

- `host list`
- `host object list`
- `host object detail`
- `host backup-config detail`

对于写入命令，优先使用 `--data` 传递结构化载荷。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `host list` 和 `host object list` 混淆。
- 不要把 `host object detail` 和 `host backup-config detail` 混淆。
- 不要在没有 `--object-id` 的情况下调用 detail 命令。
- 如果写命令返回 `jobId`，后续应切到 `job` 域继续跟踪。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询主机列表 | `foundation-cli host list --tenant-id <tenant-id> [filters]` | [list.md](./references/commands/list.md) |
| 创建主机保护对象 | `foundation-cli host object create --tenant-id <tenant-id> --data '<json>'` | [object-create.md](./references/commands/object-create.md) |
| 查询主机保护对象列表 | `foundation-cli host object list --tenant-id <tenant-id> [filters]` | [object-list.md](./references/commands/object-list.md) |
| 查询主机保护对象详情 | `foundation-cli host object detail --tenant-id <tenant-id> --object-id <object-id>` | [object-detail.md](./references/commands/object-detail.md) |
| 创建主机备份配置 | `foundation-cli host backup-config --tenant-id <tenant-id> --data '<json>'` | [backup-config.md](./references/commands/backup-config.md) |
| 查询主机备份配置详情 | `foundation-cli host backup-config detail --tenant-id <tenant-id> --object-id <object-id>` | [backup-config-detail.md](./references/commands/backup-config-detail.md) |
| 发起主机恢复 | `foundation-cli host restore start --tenant-id <tenant-id> --data '<json>'` | [restore-start.md](./references/commands/restore-start.md) |

## 参考资料

- [命令映射](./references/commands.md)
