# `foundation-cli policy delete`

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli policy delete --data '<json>'` |
| Method | `DELETE` |
| Path | `/api/sla/v1/groups` |
| Risk | `write` |

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
| body | `--data` | JSON string | 是 | 待删除计划组请求体 |

## Body 参数（`--data`）

`params/plan` 目录未定义 `plan delete` 的独立请求结构体。当前 CLI 按原样透传 JSON 到接口 `/api/sla/v1/groups`。  
建议按照现有命令实现传递删除条件对象，典型示例如下：

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| groupIds | string[] | 是（建议） | 待删除计划组 ID 列表 |

## 枚举说明

该命令 body 在 `params/plan` 中未定义枚举类型。

## 返回结果

CLI 标准返回结构：

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息（失败时） |
| responseData | object | 接口业务返回 |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "error": null,
  "status": "success",
  "responseData": {
    "groupIds": [
      "sample-group-id-1",
      "sample-group-id-2"
    ],
    "deleted": 2
  }
}
```

## 示例
```bash
foundation-cli policy delete \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"groupIds":["<group-id-1>","<group-id-2>"]}'
```

## 命令摘要

| 项 | 值 |
|---|---|
| CLI | foundation-cli policy delete |
| 风险 | write |

## 枚举列表

- 枚举值请以上文参数说明为准；未列出的取值按后端接口定义传入。

## 请求体示例

```json
{}
```
