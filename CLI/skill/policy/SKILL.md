---
name: foundation-cli-policy
description: 当用户需要通过 `foundation-cli policy` 执行策略列表、备份策略、复制策略或策略删除相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli policy --help"
---

# Foundation CLI Policy 技能

当任务属于 `foundation-cli policy` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `policy` 域。
2. 选择正确的 `policy` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要查询策略列表。
- 用户要查看备份策略详情。
- 用户要创建备份策略。
- 用户要查看复制策略详情。
- 用户要创建复制策略。
- 用户要删除策略。

## 不适用场景

- 用户要执行通用保护动作绑定或启动备份，此时应优先切换到 `protect` 域技能。
- 用户需要的是其他业务域对象能力，而不是策略模板本身。
- 用户只是要走通用 API 透传，且没有标准 `policy` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“策略列表”，使用 `policy list`，并优先补上 `--validate-status`。
- 如果用户要查“备份策略详情”，使用 `policy backup detail`。
- 如果用户要“创建备份策略”，使用 `policy backup create`。
- 如果用户要查“复制策略详情”，使用 `policy copy detail`。
- 如果用户要“创建复制策略”，使用 `policy copy create`。
- 如果用户要“删除策略”，使用 `policy delete`。
- `backup detail` 与 `copy detail` 都依赖 `--group-id`。

## 写入与只读

写入命令：

- `policy backup create`
- `policy copy create`
- `policy delete`

只读命令：

- `policy list`
- `policy backup detail`
- `policy copy detail`

对于写入命令，优先使用 `--data`；请求体较长时可改用 `--body-file`。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `backup detail` 和 `copy detail` 混淆。
- 不要把策略模板管理和 `protect` 域中的策略绑定动作混淆。
- 不要在没有 `--group-id` 的情况下调用 detail 命令。
- `delete` 是写操作，并且可能是批量删除。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询策略列表 | `foundation-cli policy list --tenant-id <tenant-id> --validate-status <1|2|3|4|5> [filters]` | [list.md](./references/commands/list.md) |
| 查看备份策略详情 | `foundation-cli policy backup detail --tenant-id <tenant-id> --group-id <group-id>` | [backup-detail.md](./references/commands/backup-detail.md) |
| 创建备份策略 | `foundation-cli policy backup create --tenant-id <tenant-id> --data '<json>'` | [backup-create.md](./references/commands/backup-create.md) |
| 查看复制策略详情 | `foundation-cli policy copy detail --tenant-id <tenant-id> --group-id <group-id>` | [copy-detail.md](./references/commands/copy-detail.md) |
| 创建复制策略 | `foundation-cli policy copy create --tenant-id <tenant-id> --data '<json>'` | [copy-create.md](./references/commands/copy-create.md) |
| 删除策略 | `foundation-cli policy delete --tenant-id <tenant-id> --data '<json>'` | [delete.md](./references/commands/delete.md) |

## 参考资料

- [命令映射](./references/commands.md)
