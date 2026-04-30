---
name: foundation-cli-job
description: 当用户需要通过 `foundation-cli job` 查询作业信息，或执行作业终止/删除时使用。
metadata:
  requires:
    bins: ["foundation-cli"]
  cliHelp: "foundation-cli job --help"
---

# Foundation CLI Job 技能
当任务属于 `foundation-cli job` 域命令时，使用本技能。

## 适用场景
- 查询作业列表、作业详情、子作业列表、作业日志。
- 查询作业业务类型和应用类型字典。
- 终止作业。
- 删除作业。

## 命令分类
只读命令：
- `job list`
- `job backup-detail`
- `job sync-detail`
- `job child list`
- `job logs`
- `job business-types`
- `job app-types`

写入命令：
- `job stop`
- `job delete`

## 快速分流
- 查列表：`job list`
- 查父作业下的子作业：`job child list`
- 查备份作业详情：`job backup-detail`
- 查复制作业详情：`job sync-detail`
- 查执行输出：`job logs`
- 停止执行中的作业：`job stop`
- 删除历史作业：`job delete`
- `job stop` 仅建议用于状态 `200/300/310/500` 的作业。
- `job delete` 仅建议用于终态 `600/700/800/900/1000/1100` 的作业。

## 关键参数
- 全局连接参数：`--tenant-id`、`--endpoint`、`--ak`、`--sk`
- 详情/日志类查询参数：`--job-id`
- 写命令参数：
- 方式 A：`--data` 或 `--body-file`（传完整 JSON）
- 方式 B：重复传 `--job-id`（自动组装为 `{"jobIds":[...]}`）

## 参考
- [命令映射](./references/commands.md)

## 作业状态（status / --statuses）

| 状态码 | 枚举名 | 含义 |
|---|---|---|
| `0` | `ExecStatusProto_Not_Start` | 未启动 |
| `10` | `ExecStatusProto_Approving` | 审批中 |
| `100` | `ExecStatusProto_Queuing` | 排队中 |
| `200` | `ExecStatusProto_Ready` | 准备中 |
| `300` | `ExecStatusProto_Running` | 正在执行 |
| `310` | `ExecStatusProto_Retrying` | 重试中 |
| `400` | `ExecStatusProto_Stopping` | 正在停止 |
| `500` | `ExecStatusProto_Abnormal` | 异常 |
| `600` | `ExecStatusProto_Failed` | 失败 |
| `700` | `ExecStatusProto_Canceled` | 已取消 |
| `800` | `ExecStatusProto_Success` | 成功 |
| `900` | `ExecStatusProto_WarnSuss` | 成功但有警告 |
| `1000` | `ExecStatusProto_PartSuss` | 部分成功 |
| `1100` | `ExecStatusProto_Stopped` | 已停止 |

说明：
- `job stop` 推荐仅用于 `200/300/310/500`。
- `job delete` 推荐仅用于终态 `600/700/800/900/1000/1100`。
