# Client 命令映射

本文件只负责一件事：帮助 agent 根据用户意图选择正确的 `foundation-cli client` 命令，并跳转到对应的单命令文档。

如果已经选中命令，不要停留在本页，直接进入对应的 `./commands/*.md`。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `client deploy`、`client deploy-config create` |
| `只读` | `client deploy-history`、`client list`、`client detail`、`client datasource list`、`client runner list`、`client runner-types`、`client deploy-config list`、`client deploy-log list` |

写命令默认继续检查 `--data`；只读命令优先检查查询条件、分页参数和关联 ID。

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `client/deploy` | 下发客户端部署、发起部署任务 | `foundation-cli client deploy --tenant-id <tenant-id> --data '<json>'` | 必填 `--data` | [deploy.md](./commands/deploy.md) |
| `client/deploy-history` | 查询部署历史、拿 `jobId` | `foundation-cli client deploy-history --tenant-id <tenant-id>` | 无请求体 | [deploy-history.md](./commands/deploy-history.md) |
| `client/list` | 查询客户端列表 | `foundation-cli client list --tenant-id <tenant-id> [--index <n>] [--count <n>] [--type <type>] [--status <status>] [--client-type <type>]` | 默认按 `index=0&count=30&type=0&status=2&clientType=0` 查询在线客户端 | [list.md](./commands/list.md) |
| `client/detail` | 查询客户端详情 | `foundation-cli client detail --tenant-id <tenant-id> --client-id <id>` | 必填 `--client-id` | [detail.md](./commands/detail.md) |
| `client/datasource/list` | 查询某客户端的数据源树 | `foundation-cli client datasource list --tenant-id <tenant-id> --client-id <id> [--full-path <path>] [--business-type <n>] [--runner-type <type>] [--runner-user <user>] [--request-id <id>] [--index <n>] [--count <n>]` | 必填 `--client-id`；MySQL 场景常用 `--business-type 1 --runner-type MySQL --runner-user root` | [datasource-list.md](./commands/datasource-list.md) |
| `client/runner/list` | 查询 Runner 列表 | `foundation-cli client runner list --tenant-id <tenant-id> [--index <n>] [--count <n>]` | 可分页 | [runner-list.md](./commands/runner-list.md) |
| `client/runner-types` | 查询 Runner 类型字典 | `foundation-cli client runner-types --tenant-id <tenant-id>` | 无请求体 | [runner-types.md](./commands/runner-types.md) |
| `client/deploy-config/create` | 按操作系统创建部署主机配置 | `foundation-cli client deploy-config create --tenant-id <tenant-id> --os <os> --data '<json>'` | 必填 `--os` 和 `--data` | [deploy-config-create.md](./commands/deploy-config-create.md) |
| `client/deploy-config/list` | 查询部署主机配置列表 | `foundation-cli client deploy-config list --tenant-id <tenant-id>` | 无请求体 | [deploy-config-list.md](./commands/deploy-config-list.md) |
| `client/deploy-log/list` | 查询某次部署的日志 | `foundation-cli client deploy-log list --tenant-id <tenant-id> --job-id <job-id> [--index <n>] [--count <n>]` | 必填 `--job-id` | [deploy-log-list.md](./commands/deploy-log-list.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `client deploy` | `--data` | 写命令默认走结构化请求体 |
| `client datasource list` | `--client-id` | 没有客户端上下文无法查询 |
| `client detail` | `--client-id` | 没有客户端上下文无法查询详情 |
| `client datasource list` | `--business-type`、`--runner-type`、`--runner-user` | `business-type` 枚举为 `1=恢复场景`、`2=备份场景`、`3=公共默认`；MySQL 恢复数据源常用 `1 / MySQL / root`；`runner-type` 建议以 `client runner-types` 实际返回为准，不要在调用侧写死枚举 |
| `client deploy-config create` | `--os`、`--data` | 先确定 OS，再构造 body |
| `client deploy-log list` | `--job-id` | 日志查询依赖具体作业 |
