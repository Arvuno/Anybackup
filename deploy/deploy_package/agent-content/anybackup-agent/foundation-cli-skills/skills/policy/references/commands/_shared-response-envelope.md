# Plan Commands Shared Response Envelope

适用范围：`skill/policy/references/commands` 下所有命令文档。

## Common Envelope

| Field | Type | Description |
|---|---|---|
| status | string | Request execution status |
| error | object/null | Error details when failed |
| responseData | object | Business payload returned by backend API |

## Notes

- 各命令文档中的“返回结果”都在此通用包络上扩展 `responseData` 的具体结构。
- 对于 `POST/PUT/DELETE` 命令，如果 params 未定义专用响应结构，默认按该包络解读。

