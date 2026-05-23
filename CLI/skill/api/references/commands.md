# API 命令映射

本文档只负责帮助 agent 选择 `foundation-cli api` 命令，并跳转到对应的单命令文档。

## 先区分只读和写入

| 类型 | 命令 |
|---|---|
| `写入` | `api` 写入透传，通常配合 `POST`、`PUT`、`PATCH`、`DELETE` 和 `--data` |
| `只读` | `api` 只读透传，通常配合 `GET` |

## 意图到命令

| 规范路径 | 用户常见意图 | 推荐命令 | 关键要求 | 单命令文档 |
|---|---|---|---|---|
| `api` | 调用一个尚无标准 CLI 封装的只读接口 | `foundation-cli api --tenant-id <tenant-id> --method GET --path </relative/path>` | 使用相对路径，不带 `--data` | [api.md](./commands/api.md) |
| `api` | 调用一个尚无标准 CLI 封装的写入接口 | `foundation-cli api --tenant-id <tenant-id> --method POST --path </relative/path> --data '<json>'` | 显式传入 `--data` | [api.md](./commands/api.md) |

## 最容易漏掉的参数

| 命令 | 易漏参数 | 说明 |
|---|---|---|
| `foundation-cli api` | `--method` | 必须明确请求方法 |
| `foundation-cli api` | `--path` | 传相对路径，而不是完整 curl 命令 |
| `foundation-cli api` | `--data` | 写入请求必须显式提供请求体 |