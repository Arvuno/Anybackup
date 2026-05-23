---
name: foundation-cli-mysql
description: 当用户需要通过 `foundation-cli mysql` 执行 MySQL 对象、时间点数据源、执行备份方案创建、备份配置、恢复配置或者发起恢复、恢复时间点、任务详情或授权相关命令时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli mysql --help"
---

# Foundation CLI MySQL 技能

当任务属于 `foundation-cli mysql` 域命令时，使用本技能。
本文档是入口，不是完整命令手册。它的作用是帮助 agent：
1. 判断用户意图是否属于 `mysql` 域。
2. 选择正确的 `mysql` 命令。
3. 跳转到 `references/commands/*.md` 中对应的单命令文档。

## 适用场景

- 用户要查询 MySQL 对象列表或对象详情。
- 用户要查询目标端 MySQL 实例列表。
- 用户要查看备份配置、恢复配置、备份任务或恢复任务详情。
- 用户要查询备份数据源、恢复时间点数据源、恢复范围或恢复时间点。
- 用户要设置备份配置、创建恢复配置或执行授权操作。

## 不适用场景

- 用户已经明确要查询通用异步作业详情或执行输出，此时应优先切换到 `job` 域技能。
- 用户需要的是 VMware、Host 等其它应用域能力，而不是 MySQL 域能力。
- 用户只是要走通用 API 透传，且没有标准 `mysql` 命令覆盖时，应回到共享技能重新判断是否改走 `api` 域。

## 标准流程

1. 先使用“快速判断规则”排除最常见的误选命令场景。
2. 再使用 [命令映射](./references/commands.md) 选择命令。
3. 立即打开 `./references/commands/` 下对应的单命令文档。
4. 基于该单命令文档构造最终 CLI。

## 快速判断规则

- 如果用户要查“对象列表”，使用 `mysql object list`。
- 如果用户要查“对象详情”，使用 `mysql object get`。
- 如果用户要查“目标端实例列表”，使用 `mysql target-instance list`。
- 如果用户要查“备份配置详情”，使用 `mysql backup-config detail`。
- 如果用户要“执行 MySQL 备份方案创建”（聚合策略创建、绑定与备份配置），按 3 步工作流执行：`policy backup create` -> `protect policy bind` -> `mysql backup-config set`。
- 如果用户要“设置备份配置”，使用 `mysql backup-config set`。
- 如果用户要查“备份数据源”，使用 `mysql datasource backup`。
- 如果用户要查“恢复时间点数据源”，使用 `mysql datasource recovery`。
- 如果用户要查“恢复时间点”，使用 `mysql recovery timepoint list`。
- 如果用户要“创建恢复配置”或者“发起恢复”，使用 `mysql restore-config create`。
- 如果用户明确是“宕机恢复/故障恢复”，`restore-config create` 请求体必须设置 `isShutDown=true`（恢复前停止数据库服务）。
- 如果用户明确是“宕机恢复/故障恢复”，`restore-config create` 请求体必须设置 `isCoverOldDb=true`（覆盖现有数据库）。
- `restore-config create` 请求体必须设置非空 `realPath`；实例级恢复使用 `["/<实例名>"]`，禁止传 `[]`、空字符串或省略。
- 进行 `mysql restore-config create` 前，必须先执行：`mysql target-instance list` -> 取 `clientId` -> `client datasource list` 取 `dataFilePath`。
- `client datasource list --client-id` 必须使用 `mysql target-instance list` 返回的同一个 `clientId`，禁止跨实例复用路径。
- MySQL 恢复场景调用 `client datasource list` 时必须传 `--runner-type MySQL --runner-user root`，禁止把目标实例的 `osUserName` 当作 `runnerUser`。
- `dataFilePath` 必须来自 `client datasource list` 返回的 `fullPath`；禁止直接填 `/`。
- 如果 `client datasource list` 返回空，不要立刻改填猜测路径；先按固定顺序继续查询。
- `client datasource list` 防空轮询必须按固定顺序执行，且 `requestId` 必须取第一次请求返回的 `responseData.requestId`。
- 第 1 次：`fullPath=""`，不传 `requestId` 或传空，读取返回的 `responseData.requestId`。
- 第 2 次：`fullPath=""` + `requestId=<第一次返回的 requestId>`。
- 第 3 次：`fullPath="/"` + `requestId=<第一次返回的 requestId>`。
- 若第 3 次后已拿到根目录子项，继续按固定顺序下钻：`fullPath="/var"` -> `fullPath="/var/lib"` -> `fullPath="/var/lib/mysql"`。
- 若 `/var/lib` 返回多条路径，优先选择最接近 MySQL 数据目录语义的路径，通常是 `/var/lib/mysql`。
- 若 `/var/lib/mysql` 再继续展开后返回 `/var/lib/mysql/mysql`、`/var/lib/mysql/performance_schema` 等更深层子目录，`dataFilePath` 仍优先使用上层目录 `/var/lib/mysql`。
- 若完成上述下钻仍获取不到有效路径，MySQL 恢复场景默认 `dataFilePath="/var/lib/mysql"`。
- `backup-config detail` 和 `recovery-config detail` 都依赖 `--system-id` 和 `--object-id`。
- `backup-detail` 和 `recovery-detail` 都依赖 `--task-id`。

## 写入与只读

写入命令：
- `mysql backup-plan create`
- `mysql backup-config set`
- `mysql restore-config create`
- `mysql authorize`

只读命令：
- `mysql object list`
- `mysql target-instance list`
- `mysql object get`
- `mysql backup-config detail`
- `mysql datasource backup`
- `mysql backup-detail`
- `mysql recovery range`
- `mysql recovery-config detail`
- `mysql datasource recovery`
- `mysql recovery timepoint list`
- `mysql recovery-detail`

对于写入命令，优先使用 `--data`；请求体较长时可改用 `--body-file`。

## 必要的全局上下文

除非调用环境已经隐式提供，否则先确认这些参数：
- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 常见误判

- 不要把 `object list` 和 `object get` 混淆。
- 不要把 `backup-config detail` 和 `backup-config set` 混淆。
- 不要把 `datasource backup` 和 `datasource recovery` 混淆。
- 不要把 `backup-detail` 和 `recovery-detail` 混淆。
- 不要在缺少 `--system-id`、`--object-id`、`--task-id`、`--storage-service-id` 的情况下直接拼命令。

## 意图到命令

| 意图 | 推荐命令 | 下一步打开的文件 |
|---|---|---|
| 查询 MySQL 对象列表 | `foundation-cli mysql object list --tenant-id <tenant-id> --app-type 202 [filters]` | [object-list.md](./references/commands/object-list.md) |
| 查询目标端 MySQL 实例列表 | `foundation-cli mysql target-instance list --tenant-id <tenant-id> [--index <n>] [--count <n>]` | [target-instance-list.md](./references/commands/target-instance-list.md) |
| 查询 MySQL 对象详情 | `foundation-cli mysql object get --tenant-id <tenant-id> --object-id <object-id>` | [object-get.md](./references/commands/object-get.md) |
| 查询备份配置详情 | `foundation-cli mysql backup-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` | [backup-config-detail.md](./references/commands/backup-config-detail.md) |
| 执行 MySQL 备份方案创建（聚合） | 按工作流执行：`foundation-cli policy backup create --data <policyData>` -> `foundation-cli protect policy bind --data <bindData>` -> `foundation-cli mysql backup-config set --data <backupConfigData>` | [backup-plan-create.md](./references/commands/backup-plan-create.md) |
| 设置备份配置 | `foundation-cli mysql backup-config set --tenant-id <tenant-id> --data '<json>'` | [backup-config-set.md](./references/commands/backup-config-set.md) |
| 查询备份数据源 | `foundation-cli mysql datasource backup --tenant-id <tenant-id> --object-id <object-id> [query flags]` | [datasource-backup.md](./references/commands/datasource-backup.md) |
| 查询备份任务详情 | `foundation-cli mysql backup-detail --tenant-id <tenant-id> --task-id <task-id>` | [backup-detail.md](./references/commands/backup-detail.md) |
| 查询恢复范围 | `foundation-cli mysql recovery range --tenant-id <tenant-id> [query flags]` | [recovery-range.md](./references/commands/recovery-range.md) |
| 查询恢复配置详情 | `foundation-cli mysql recovery-config detail --tenant-id <tenant-id> --system-id <system-id> --object-id <object-id>` | [recovery-config-detail.md](./references/commands/recovery-config-detail.md) |
| 查询恢复时间点数据源 | `foundation-cli mysql datasource recovery --tenant-id <tenant-id> --data-set-id <data-set-id> --storage-service-id <id> --timestamp <ts> [query flags]` | [datasource-recovery.md](./references/commands/datasource-recovery.md) |
| 查询恢复时间点 | `foundation-cli mysql recovery timepoint list --tenant-id <tenant-id> [query flags]` | [recovery-timepoint-list.md](./references/commands/recovery-timepoint-list.md) |
| 查询恢复任务详情 | `foundation-cli mysql recovery-detail --tenant-id <tenant-id> --task-id <task-id>` | [recovery-detail.md](./references/commands/recovery-detail.md) |
| 创建恢复配置或发起恢复 | `foundation-cli mysql restore-config create --tenant-id <tenant-id> --data '<json>'` | [restore-config-create.md](./references/commands/restore-config-create.md) |
| 授权 MySQL 实例 | `foundation-cli mysql authorize --tenant-id <tenant-id> --data '<json>'` | [authorize.md](./references/commands/authorize.md) |

## 参考资料

- [命令映射](./references/commands.md)

## 执行 MySQL 备份方案创建（聚合）参数映射

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| policyData | `--policy-data` | string | 是 | 透传到 `foundation-cli policy backup create --data` |
| bindData | `--bind-data` | string | 是 | 透传到 `foundation-cli protect policy bind --data`，其中 `slaIds` 由第 1 步返回值自动注入 |
| backupConfigData | `--backup-config-data` | string | 是 | 透传到 `foundation-cli mysql backup-config set --data` |
