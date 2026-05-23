---
name: foundation-cli-client
description: 当用户需要通过 `foundation-cli client` 执行客户端部署、部署历史、客户端清单、Runner、数据源、部署配置或部署日志相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli client --help"
---

# Foundation CLI Client 技能
当任务属于 `foundation-cli client` 域命令时，使用本技能。

本文档是入口，不是完整参考手册。它的作用是帮助 agent：
1. 判断用户意图是否属于 `client` 域。
2. 选择正确的 `client` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要下发客户端部署配置。
- 用户要查询部署历史。
- 用户要查询客户端列表。
- 用户要查询客户端详情。
- 用户要查询某个客户端下的数据源树。
- 用户要查询 Runner 列表或 Runner 类型信息。
- 用户要创建或查询部署主机配置。
- 用户要查看某个已知作业的部署日志。

## 不适用场景

- 用户已经明确要查询通用异步作业详情或作业输出，此时应优先切换到 `job` 域技能。
- 用户需要的是其他业务域对象能力，而不是客户端、Runner、数据源或部署相关操作。
- 用户只是要走通用 API 透传，且没有标准 `client` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

不要把本文档当成查询完整请求体、示例或排错细节的地方。

## 快速判断规则

选择命令前，先用下面这些区分规则：
- 如果用户要发起部署任务，使用 `client deploy`。
- 如果用户要查看历史部署记录，或需要拿到 `jobId`，使用 `client deploy-history`。
- 如果用户要查看某个已知作业的执行日志，使用 `client deploy-log list`。
- 如果用户要查询客户端清单，使用 `client list`。
- 如果用户要查询某个客户端详情，使用 `client detail`。
- 如果用户要查询某个客户端下的数据源信息，使用 `client datasource list`。
  - 这个命令必须带 `--client-id`。
- 如果用户要查询 Runner 实例列表，使用 `client runner list`。
- 如果用户要查询 Runner 类型字典，使用 `client runner-types`。
- 如果用户要创建部署主机配置，使用 `client deploy-config create`。
  - 这个命令必须同时带 `--os` 和 `--data`。
- 如果用户要查看已有部署主机配置，使用 `client deploy-config list`。

## 写入与只读

写入命令：
- `client deploy`
- `client deploy-config create`

只读命令：
- `client deploy-history`
- `client list`
- `client detail`
- `client datasource list`
- `client runner list`
- `client runner-types`
- `client deploy-config list`
- `client deploy-log list`

对于写入命令，优先使用 `--data` 传递结构化载荷。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：
- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `deploy-history` 和 `deploy-log list` 混淆。
  - 历史返回的是任务记录。
  - 部署日志返回的是日志明细，并且需要 `--job-id`。
- 不要把 `runner list` 和 `runner-types` 混淆。
  - `runner list` 返回 Runner 实例。
  - `runner-types` 返回类型字典。
- 不要在没有 `--client-id` 的情况下使用 `datasource list`。
- 不要在没有 `--client-id` 的情况下使用 `detail`。
- 不要把 `deploy-config create` 当成只读操作。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 下发客户端部署配置 | `foundation-cli client deploy --tenant-id <tenant-id> --data '<json>'` | [deploy.md](./references/commands/deploy.md) |
| 查询部署历史 | `foundation-cli client deploy-history --tenant-id <tenant-id>` | [deploy-history.md](./references/commands/deploy-history.md) |
| 查询客户端列表 | `foundation-cli client list --tenant-id <tenant-id> [--index <n>] [--count <n>] [--type <type>] [--status <status>] [--client-type <type>]` | [list.md](./references/commands/list.md) |
| 查询客户端详情 | `foundation-cli client detail --tenant-id <tenant-id> --client-id <id>` | [detail.md](./references/commands/detail.md) |
| 查询某客户端的数据源 | `foundation-cli client datasource list --tenant-id <tenant-id> --client-id <id> [--full-path <path>] [--business-type <n>] [--runner-type <type>] [--runner-user <user>] [--request-id <id>] [--index <n>] [--count <n>]` | [datasource-list.md](./references/commands/datasource-list.md) |
| 查询 Runner 列表 | `foundation-cli client runner list --tenant-id <tenant-id> [--index <n>] [--count <n>]` | [runner-list.md](./references/commands/runner-list.md) |
| 查询 Runner 类型 | `foundation-cli client runner-types --tenant-id <tenant-id>` | [runner-types.md](./references/commands/runner-types.md) |
| 按操作系统创建部署主机配置 | `foundation-cli client deploy-config create --tenant-id <tenant-id> --os <os> --data '<json>'` | [deploy-config-create.md](./references/commands/deploy-config-create.md) |
| 查询部署主机配置列表 | `foundation-cli client deploy-config list --tenant-id <tenant-id>` | [deploy-config-list.md](./references/commands/deploy-config-list.md) |
| 查询部署日志 | `foundation-cli client deploy-log list --tenant-id <tenant-id> --job-id <job-id> [--index <n>] [--count <n>]` | [deploy-log-list.md](./references/commands/deploy-log-list.md) |

## 参考资料

- [命令映射](./references/commands.md)

## `client datasource list` 关键说明

- `--business-type 1`：恢复场景。
- `--business-type 2`：备份场景。
- `--business-type 3`：公共默认。
- MySQL 恢复数据源常用组合：`--business-type 1 --runner-type MySQL --runner-user root`
- `--runner-type` 建议优先参考 `client runner-types` 的实际返回结果，不要在调用侧写死固定枚举。
