# `foundation-cli protect policy create-bind`

## 命令摘要

先创建备份策略，再将该策略绑定到指定保护对象。  
如果绑定失败，命令会自动调用删除接口回滚刚创建的策略。

## 请求体示例

```json
{
  "name": "demo-backup-policy",
  "validatePeriod": 1,
  "effectiveType": 1,
  "backupConfig": {
    "dayEnable": true,
    "dayConfiguration": [
      {
        "interval": 1,
        "triggerTime": {
          "hour": 2,
          "min": 0
        },
        "mode": 1
      }
    ]
  }
}
```

## 命令概览

| 项 | 值 |
|---|---|
| CLI | `foundation-cli protect policy create-bind --object-id <object-id> --data '<json>'` |
| Method | `POST` + `POST`（失败补偿：`DELETE`） |
| Path | 创建：`/api/sla/v1/group/backup_info`；绑定：`/backupmgm/v1/protect_object/{objectId}/slas`；回滚删除：`/api/sla/v1/groups` |
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
| objectId | `--object-id` | string | 是 | 保护对象 ID |
| createBody | `--data` | JSON string | 是 | 备份策略创建请求体 |

## Body 参数

对象：`AddSlaBackupRequest`

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 策略名称 |
| validatePeriod | int | 否 | 有效期类型；未传时 CLI 自动补 `1` |
| effectiveType | int | 否 | 生效类型；未传时 CLI 自动补 `1` |
| backupConfig | object | 是 | 备份计划配置 |

## 枚举列表

- 本命令的 `--data` 枚举含义与 `foundation-cli policy backup create` 保持一致。
- `bind` 阶段的请求体由 CLI 自动生成为 `{"slaIds":["<created-id>"]}`，无需额外枚举参数。

## 执行流程

1. 调用策略重名预检查：`GET /api/sla/v1/common/name/exists?name=<name>`。
2. 调用创建策略接口，解析返回 `id`（或 `groupId`）。
3. 构造绑定请求体 `{"objectId":"<object-id>","slaIds":["<created-id>"]}`，绑定到 `--object-id`。
4. 若绑定失败，自动调用删除接口 `DELETE /api/sla/v1/groups`，请求体 `{"groupIds":["<created-id>"]}`。

## 返回结果

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 执行状态 |
| error | object/null | 错误信息 |
| responseData.policyId | string | 本次创建出的策略 ID |
| responseData.create | object | 创建接口原始返回 |
| responseData.bind | object | 绑定接口原始返回 |
| responseData.rollback | object/null | 成功场景固定为 `null` |

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "policyId": "a15e75fbdab04effbc7bf076aac8c7a8",
    "create": {
      "status": "success",
      "error": null,
      "responseData": {
        "id": "a15e75fbdab04effbc7bf076aac8c7a8"
      }
    },
    "bind": {
      "status": "success",
      "error": null,
      "responseData": null
    },
    "rollback": null
  }
}
```

## 示例

```bash
foundation-cli protect policy create-bind \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id> \
  --data '{"name":"demo-backup-policy","validatePeriod":1,"effectiveType":1,"backupConfig":{"dayEnable":true,"dayConfiguration":[{"interval":1,"triggerTime":{"hour":2,"min":0},"mode":1}]}}'
```
