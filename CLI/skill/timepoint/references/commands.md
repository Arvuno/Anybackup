# Timepoint 命令映射

本文件只负责帮助 agent 选择正确的 `foundation-cli timepoint` 命令，并跳转到对应的单命令文档。

## 先区分读写

| 类型 | 命令 |
|---|---|
| `写入` | `timepoint clean start` |
| `只读` | `timepoint list` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `timepoint/list` | 查询对象时间点列表 | `foundation-cli timepoint list --tenant-id <tenant-id> --object-id <object-id> [filters]` | 必填 `--object-id` | [list.md](./commands/list.md) |
| `timepoint/clean/start` | 发起时间点清理任务 | `foundation-cli timepoint clean start --tenant-id <tenant-id> --object-id <object-id> --data '<json>'` | 必填 `--object-id` 和 `--data` | [clean-start.md](./commands/clean-start.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `timepoint list` | `--object-id` | 时间点查询依赖对象上下文 |
| `timepoint clean start` | `--object-id`、`--data` | 清理任务必须同时给对象和请求体 |
