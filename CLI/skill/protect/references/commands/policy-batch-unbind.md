# `foundation-cli protect policy batch-unbind`

## 命令摘要

批量移除保护对象已绑定的策略。

## 请求体示例

```json
{
  "objectIds": [
    "sample-object-id-1",
    "sample-object-id-2"
  ],
  "isAll": true
}
```

## 枚举列表

- 该命令当前请求体未定义显式枚举字段。

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect policy batch-unbind --data '<json>'` |
| Method | `DELETE` |
| Path | `/backupmgm/v1/protect_object/slas` |
| Risk | `write` |

## 请求参数

### 共享参数

| 字段 | CLI Flag | 必填 | 说明 |
|---|---|---|---|
| tenantId | `--tenant-id` | 否 | 租户标识 |
| endpoint | `--endpoint` | 是 | Foundation 服务地址 |
| ak | `--ak` | 是 | Access Key |
| sk | `--sk` | 是 | Secret Key |
| targetVersion | `--target-version` | 否 | 目标版本，默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI Flag | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| body | `--data` | JSON string | 是 | 批量解绑请求体 |

## Body 参数（`--data`）

对象：`BatchRemoveSLAsRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| objectIds | string[] | 是 | 保护对象 ID 列表，最多 100 条，每项长度 32 |
| isAll | bool | 否 | 是否解绑全部已绑定策略 |

## 枚举说明

该命令 body 在 `params/protect/sla_unbind_batch.go` 中未定义枚举类型。

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData | object | 解绑结果 |

## 返回案例

说明：
1. 当前仓库未附带该命令的真实接口成功响应。
2. 下例仅用于说明常见返回结构层级，字段名和取值以后端当前版本实际返回为准。
3. 建议后续补充 1 份真实成功响应和 1 份典型失败响应。

```json
{
  "status": "success",
  "error": null,
  "responseData": null
}
```

## 示例
```bash
foundation-cli protect policy batch-unbind \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '{"items":[{"objectId":"<object-id>","slaId":"<sla-id>"}]}'
```

