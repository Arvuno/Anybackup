# Policy Commands Parameter & Enum Commonality Analysis

分析范围：`skill/policy/references/commands`

## Commands

- `foundation-cli policy list`
- `foundation-cli policy backup create`
- `foundation-cli policy copy create`
- `foundation-cli policy backup detail`
- `foundation-cli policy copy detail`
- `foundation-cli policy delete`

## Common Parameter Blocks

| Common Block | Included Commands | Notes |
|---|---|---|
| Shared request context (`tenant-id/endpoint/ak/sk/target-version`) | all 6 | 强公共块，已抽取到 `_shared-request-context.md` |
| Standard response envelope (`status/error/responseData`) | all 6 | 强公共块，已抽取到 `_shared-response-envelope.md` |
| Pagination (`index/count`) | list | 弱公共块，仅 list 使用 |
| Path param (`group-id`) | backup detail, copy detail | 同形态公共块，可在详情类命令复用 |
| Body transport (`--data`) | backup create, copy create, delete | 写命令公共输入形态 |

## Common Enum Blocks

| Enum | Included Commands |
|---|---|
| `ValidatePeriodType` | backup create, copy create, backup detail, copy detail |
| `ValidateType` | backup create, copy create, backup detail, copy detail |
| `EffectiveCType` | backup create, copy create, backup detail, copy detail |
| `BusinessType` | list, copy create, backup detail, copy detail |
| `ExecMode` | list, copy create, copy detail |
| `RetentionCType` | list, backup create, copy create, backup detail, copy detail |
| `RetentionDurationType` | backup create, copy create, backup detail, copy detail |
| `BackupMode` | backup create, copy create, backup detail, copy detail |
| `MonthTriggerType` | backup create, backup detail, copy detail |
| `SpeedUnit` | backup create, copy create, backup detail, copy detail |
| `CopyType` | copy create, copy detail |
| `TemplateLevel` | list |
| `ValidateStatusType` | list |

## Split Strategy

- 高复用（3 个命令以上）统一收敛到 `_shared-enums.md`。
- 所有命令统一引用 `_shared-request-context.md` 与 `_shared-response-envelope.md`。
- 命令文档只保留“命令特有参数/特有枚举”的详细说明，公共项通过引用跳转。

