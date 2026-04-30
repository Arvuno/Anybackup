> 命名约定：标准命令统一使用 `foundation-cli {{domain}} <resource> <action>`；兼容别名仅用于迁移或兜底，不作为首选命令。

## 核心规则

1. 有标准业务命令时，优先使用标准业务命令，不优先使用 `foundation-cli curl`。
2. 写操作前先确认租户上下文：`tenant-id`、`endpoint`、`targetVersion` 是否正确。
3. 请求体统一通过 `--data` 传递，并保持 JSON 脱敏且尽量精简。
4. 数组型 query 参数统一使用重复 flag，不使用逗号拼接作为默认约定。
5. 长任务启动命令成功时，只表示“控制台已受理”，不表示后端任务已完成；应返回并保留 `jobId`。
6. 若控制台返回 `401/402/403`，必须按“AK/SK 失败 / token 失效 / 权限不足”三类分别解释。
7. `stdout` 只保留结果包；审计、版本命中、告警和排障信息统一进入 `stderr`。

## 意图 -> 命令索引

| 意图 | 推荐命令 | 备注 |
|---|---|---|
{{intent_rows}}

建议至少覆盖：

- 查看列表
- 查看详情
- 创建/绑定/配置
- 启动任务
- 查询任务结果
- 异常场景下的兜底命令

## 操作注意事项

1. 查询类命令应说明分页参数、过滤参数和数组参数传递方式。
2. 写入类命令应说明请求体结构、幂等责任边界和失败后的排障入口。
3. 长任务类命令应说明如何使用 `foundation-cli job backup-detail --job-id <job-id>` 或 `foundation-cli job sync-detail --job-id <job-id>` 继续查询。
4. 若使用 `foundation-cli curl` 作为兜底方案，必须写清允许访问的 endpoint、签名要求和输出模式。
5. 每个命令最好配一条最小可执行示例，并标出关键参数含义。

## 推荐示例骨架

```bash
# 查询类
foundation-cli {{domain}} {{resource}} list --tenant-id <tenant-id> [flags]

# 写入类
foundation-cli {{domain}} {{resource}} {{write_action}} --tenant-id <tenant-id> --data '<json>'

# 长任务查询
foundation-cli job backup-detail --job-id <job-id>

# 特殊兜底
foundation-cli curl --tenant-id <tenant-id> -X <METHOD> <URL> [--data "<json-body>"]
```

## 常见错误说明

| 场景 | 控制台 HTTP 状态码 | CLI 退出码 | 建议说明 |
|---|---|---|---|
| `AK/SK` 验证失败 | `401` | `201` | 检查 `AK`、`SK`、签名和租户凭据文件 |
| `token` 过期/失效 | `402` | `201` | 重新获取 `token` 或重新登录 |
| 无权限访问 | `403` | `201` | 检查租户、角色和目标资源授权 |
| 参数非法 | `400` 或业务失败 | `101` | 检查 flag、body 和字段格式 |
| 兼容不支持 | 业务侧返回或本地判定 | `551` | 检查 `targetVersion`、支持矩阵和本地映射 |

## 审计与回包要求

1. 统一返回包格式：`status/error/responseData`。
2. 长任务命令返回成功时，`responseData` 中应优先包含 `jobId`。
3. `stderr` 中建议输出 `traceId`、版本命中来源、目标 endpoint 和关键告警。
4. 日志和示例中不得出现真实 `AK/SK`、生产租户 ID 或生产主机地址。

