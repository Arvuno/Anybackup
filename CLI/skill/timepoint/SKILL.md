---
name: foundation-cli-timepoint
description: 当用户需要通过 `foundation-cli timepoint` 执行对象时间点列表查询或发起时间点清理任务时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli timepoint --help"
---

# Foundation CLI Timepoint 技能

当任务属于 `foundation-cli timepoint` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否属于 `timepoint` 域。
2. 选择正确的 `timepoint` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要查询某个保护对象的时间点列表。
- 用户要为某个保护对象发起时间点清理任务。
- 请求针对所有应用对象时间点，而 MySQL 恢复场景下有专属的时间点接口。

## 不适用场景

- 用户要查询 MySQL 恢复场景下的专属时间点能力，此时应优先切换到对应应用域技能。
- 用户需要的是其他业务域对象能力，而不是通用对象时间点相关操作。
- 用户只是要走通用 API 透传，且没有标准 `timepoint` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“对象时间点列表”，使用 `timepoint list`。
- 如果用户要发起“时间点清理任务”，使用 `timepoint clean start`。
- `timepoint list` 必须带 `--object-id`。
- `timepoint clean start` 必须同时带 `--object-id` 和 `--data`。
- 如果用户说的是 MySQL 恢复场景下的专属时间点，不要优先留在 `timepoint` 域。

## 写入与只读

写入命令：

- `timepoint clean start`

只读命令：

- `timepoint list`

对于写入命令，优先使用 `--data` 传递结构化载荷。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把通用 `timepoint` 域和 MySQL 恢复场景下的专属时间点接口混淆。
- 不要在没有 `--object-id` 的情况下调用 `timepoint list`。
- 不要在没有 `--data` 的情况下调用 `timepoint clean start`。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询对象时间点列表 | `foundation-cli timepoint list --tenant-id <tenant-id> --object-id <object-id> [filters]` | [list.md](./references/commands/list.md) |
| 发起时间点清理任务 | `foundation-cli timepoint clean start --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | [clean-start.md](./references/commands/clean-start.md) |

## 参考资料

- [命令映射](./references/commands.md)
