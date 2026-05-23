# `foundation-cli client deploy-history`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli client deploy-history --tenant-id <tenant-id> [--index <n>] [--count <n>]` |
| 方法 | `GET` |
| 路径 | 由标准 `client deploy-history` 命令封装 |
| 风险 | `只读` |

## 何时使用

当用户要查看历史部署任务、获取某次部署的 `jobId`，或判断部署任务状态时，使用本命令。

## 必填上下文

- `--tenant-id`
- `--endpoint`
- `--ak`
- `--sk`

## 最小可执行示例

```bash
foundation-cli client deploy-history \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk>
```

## 返回结果

依据 `params/client/history.go`，常用字段包括：

- `responseData.totalNum`
- `responseData.data[].jobId`
- `responseData.data[].jobName`
- `responseData.data[].jobType`
- `responseData.data[].status`
- `responseData.data[].startTime`
- `responseData.data[].endTime`

## 常见错误

| 错误 | 原因 | 处理方式 |
|---|---|---|
| 误把历史当日志 | 该命令返回任务记录，不返回日志明细 | 查看日志请用 [deploy-log-list.md](./deploy-log-list.md) |
| 需要继续跟踪过程 | 只有历史还不够 | 取出 `jobId` 后继续查日志或 `job` 域详情 |

## 相关跳转

- 查部署日志： [deploy-log-list.md](./deploy-log-list.md)
- 字段来源：`params/client/history.go`

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli\ client\ deploy-history |
| 风险 | read-only |

## 请求参数

请结合上文参数表传参；常见查询命令可按需传 `--index` 与 `--count` 控制分页。

## 返回案例

```json
{"status":"success","error":null,"responseData":null}
```

## 示例

```bash
foundation-cli client deploy-history
```
