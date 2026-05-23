---
name: foundation-cli-api
description: 当用户需要通过 `foundation-cli api` 访问暂时没有标准业务命令覆盖的 Foundation REST API，或者用户直接给出原始路径、HTTP 方法、查询参数或请求体时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli api --help"
---

# Foundation CLI API 技能

当任务属于 `foundation-cli api` 域命令时，使用本技能。

本文件是入口，不是完整参考手册。
它的作用是帮助 agent：
1. 判断用户意图是否应落到 `api` 域。
2. 先排除已有标准业务命令覆盖的场景。
3. 选择正确的 `foundation-cli api` 透传方式。
4. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要访问尚未被标准业务命令覆盖的 Foundation REST API。
- 用户直接给出原始 REST 路径、HTTP 方法、查询参数或请求体。
- 用户明确要求通过签名透传方式调用某个接口。

## 不适用场景

- 用户要执行的其实是 `host`、`job`、`mysql`、`network`、`policy`、`protect`、`storage`、`timepoint`、`vmware` 等业务域已有标准命令覆盖的操作。
- 用户只是描述业务意图，但没有要求原始 API，且标准 CLI 已经可以表达该操作。
- 用户需要追踪异步任务结果时，不应继续停留在 `api` 域，而应切换到 `job` 域。

## 标准流程

1. 先使用“快速判断规则”确认是否真的需要 `api` 透传，而不是已有业务命令。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于单命令文档补齐 `--method`、`--path`、`--data` 和全局上下文。

## 快速判断规则

- 如果标准业务命令已经存在，优先使用标准业务命令，不要落到 `api`。
- 如果用户给的是原始相对路径和方法，通常应落到 `api`。
- 如果是只读透传，通常使用 `GET` 且不带 `--data`。
- 如果是写入透传，必须显式提供 `--data`，不要把请求体省略掉。
- `--path` 应使用相对路径，不要手写 curl 或自行拼签名。
- 如果接口返回异步任务标识，例如 `jobId`，下一步应切到 `job` 域跟踪。

## 写入与只读

写入命令：

- `api` 写入透传，请显式使用 `POST`、`PUT`、`PATCH` 或 `DELETE`，并补齐 `--data`

只读命令：

- `api` 只读透传，通常使用 `GET`

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把“已有标准命令的业务操作”错误地落到 `api`。
- 不要把原始 API 透传和业务域命令混用。
- 不要遗漏 `--method` 或把写入请求误当成只读请求。
- 写入透传不要省略 `--data`。
- 异步任务返回后不要继续停留在 `api`，应转到 `job` 域跟踪。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询一个尚无标准 CLI 覆盖的只读接口 | `foundation-cli api --tenant-id <tenant-id> --method GET --path </relative/path>` | [api.md](./references/commands/api.md) |
| 调用一个尚无标准 CLI 覆盖的写入接口 | `foundation-cli api --tenant-id <tenant-id> --method POST --path </relative/path> --data '<json>'` | [api.md](./references/commands/api.md) |

## 参考资料

- [命令映射](./references/commands.md)
