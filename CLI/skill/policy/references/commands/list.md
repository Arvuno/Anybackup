# `foundation-cli policy list`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy list [--index <n>] [--count <n>] [--filter <text>] [--validate-status <1|2|3|4|5>] [--disable-mark <-1|0|1|2>] [--type <business-type>] [--types <business-type> ...] [--copy-mode <1|2>] [--backup-mode <1|2|3|4|5>]` |
| Method | `GET` |
| Path | `/api/sla/v1/templates` |
| Risk | `read-only` |

## 公共引用

- 共享请求上下文：[`_shared-request-context.md`](./_shared-request-context.md)
- 通用响应包络：[`_shared-response-envelope.md`](./_shared-response-envelope.md)
- 通用枚举定义：[`_shared-enums.md`](./_shared-enums.md)
- 参数与枚举共性分析：[`_parameter-enum-common-analysis.md`](./_parameter-enum-common-analysis.md)

## 请求参数

### 共享参数

共享参数已统一收敛，详见：[`_shared-request-context.md`](./_shared-request-context.md)。

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| index | `--index` | int | 否 | 分页起始，默认 `0` |
| count | `--count` | int | 否 | 分页条数，默认 `10` |
| filter | `--filter` | string | 否 | 模糊过滤 |
| validateStatus | `--validate-status` | string | 否 | 生效状态，`1/2/3/4/5`；当前环境下实测需要显式传入 |
| disableMark | `--disable-mark` | string | 否 | 禁用状态，`-1/0/1/2` |
| type | `--type` | string | 否 | 业务类型 |
| types | `--types` | string[] | 否 | 业务类型数组，最多 10 个 |
| copyMode | `--copy-mode` | string | 否 | 复制模式，`1/2` |
| backupMode | `--backup-mode` | string | 否 | 备份模式，`1/2/3/4/5` |

## 返回参数

CLI 通用外层结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 业务返回体 |

`responseData` 常见分页结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| totalNum | int64 | 总数 |
| data | object[] | 计划模板列表，元素结构见下表 |

`responseData.data[]`（`GetTemplateInfoResponse`）核心字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | SLA ID |
| name | string | SLA 名称 |
| validateStatus | int | 生效状态 |
| type | int | 业务类型 |
| status | int | 是否禁用 |
| bindResource | int | 绑定保护对象数量 |
| expiredTime | string | 到期时间展示值 |
| creatorName | string | 创建者 |
| originGroupId | string | 原计划组 ID |
| tenantName | string | 租户名称 |
| templateId | string | 模板 ID |
| isAddChild | bool | 是否可追加子策略 |
| isIncludeLog | bool | 是否包含日志 |
| execMode | int | 数据执行模式 |
| logExecMode | int | 日志执行模式 |
| dataCopyRule | object | 数据复制规则 |
| logCopyRule | object | 日志复制规则 |
| storageServiceId | string | 存储服务 ID |
| storagePoolId | string | 存储池 ID |
| storagePoolType | int | 存储池类型 |
| logRetention | object | 日志保留策略 |
| dataRetention | object | 数据保留策略 |
| snapshotRetention | object | 快照保留配置 |
| preRecoveryRetention | object | 预恢复保留配置 |
| children | object[] | 子策略列表（递归） |

`children[]`（`ChildTemplateInfoResponse`）字段：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | string | 子模板 ID |
| slaId | string | SLA ID |
| name | string | 模板名称 |
| level | int | 模板层级 |
| parentId | string | 父模板 ID |
| creatorName | string | 创建者 |
| type | int | 业务类型 |
| storageServiceId | string | 存储服务 ID |
| storagePoolId | string | 存储池 ID |
| storagePoolType | int | 存储池类型 |
| dataRetention | object | 数据保留策略 |
| children | object[] | 子策略列表 |

## 枚举说明

本命令使用的枚举均来自共享枚举，详见：[`_shared-enums.md`](./_shared-enums.md)。

字段与共享枚举映射：
- `validateStatus` -> `ValidateStatusType`
- `type` -> `BusinessType`
- `level` -> `TemplateLevel`
- `execMode` / `logExecMode` -> `ExecMode`
- `retentionType` -> `RetentionCType`

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
    "error":  null,
    "status":  "success",
    "responseData":  {
                         "totalNum":  1,
                         "data":  [
                                      {
                                          "id":  "sample-id",
                                          "name":  "sample-name",
                                          "validateStatus":  1
                                      }
                                  ]
                     }
}
```

## 示例
```bash
foundation-cli policy list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --validate-status 1
```

```bash
foundation-cli policy list \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --filter oracle \
  --validate-status 1 \
  --types mysql \
  --types vmware \
  --count 20
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli policy list |
| 风险 | read-only |
