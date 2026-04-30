# Policy Commands Shared Request Context

适用范围：`skill/policy/references/commands` 下所有命令文档。

## Shared Parameters

| Field | CLI Flag | Type | Required | Description |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | No | Tenant ID; can come from FOUNDATION_TENANT_ID env var |
| endpoint | `--endpoint` | string | Yes | Foundation Console endpoint |
| ak | `--ak` | string | Yes | Access Key |
| sk | `--sk` | string | Yes | Secret Key |
| targetVersion | `--target-version` | string | No | Target version, default `8.0.9.0` |

## Notes

- `tenantId/endpoint/ak/sk` 为所有 policy 命令的统一认证上下文参数。
- `targetVersion` 在未显式传入时按命令默认值处理。

