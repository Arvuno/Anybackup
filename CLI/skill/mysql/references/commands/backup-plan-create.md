# `foundation-cli mysql backup-plan create`

## 命令摘要

MySQL 备份方案聚合命令（Skill 编排命令）：按固定顺序执行 3 条已存在写命令，完成“创建策略 + 绑定策略 + 下发 MySQL 备份配置”。

说明：该命令是 Skill 层聚合能力，当前 `foundation-cli` 二进制未内置同名子命令，不能直接在终端执行 `foundation-cli mysql backup-plan create`。

## 命令概览

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql backup-plan create --policy-data '<json>' --bind-data '<json>' --backup-config-data '<json>'` |
| 方法 | `POST`（聚合编排） |
| 路径 | `/api/sla/v1/group/backup_info` -> `/backupmgm/v1/protect_object/{objectId}/slas` -> `/backupmgm/v1/mysql/app_backup_config` |
| 风险 | `写入` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 地址 |
| ak | `--ak` | string | 是 | 访问密钥 AK |
| sk | `--sk` | string | 是 | 访问密钥 SK |
| targetVersion | `--target-version` | string | 否 | 默认 `9.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| policyData | `--policy-data` | string | 是 | 透传到 `foundation-cli policy backup create --data` |
| bindData | `--bind-data` | string | 是 | 透传到 `foundation-cli protect policy bind --data`，其中 `slaIds` 由第 1 步返回值自动注入 |
| backupConfigData | `--backup-config-data` | string | 是 | 透传到 `foundation-cli mysql backup-config set --data` |

## Body 参数

聚合请求体由 3 段 JSON 组成，采用透传模式。

### `policyData`（策略创建）

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| name | string | 是 | 策略名称 |
| validatePeriod | int | 是 | 有效期类型 |
| effectiveType | int | 是 | 生效类型 |
| backupConfig | object | 是 | 备份计划配置 |

### `bindData`（策略绑定）

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| objectId | string | 是 | 保护对象 ID（用于路径参数） |
| slaIds | string[] | 否 | 绑定策略 ID 列表；由聚合命令自动注入第 1 步创建出的 `slaId` |

### `backupConfigData`（MySQL 备份配置）

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| objectId | string | 是 | 保护对象 ID |
| systemId | string | 是 | 生产系统 ID |
| commonConfigParams | object | 是 | 公共配置 |
| backupType | int | 是 | 备份类型 |
| backupGran | int | 是 | 备份粒度 |

## 枚举列表

- `backupConfigData.backupType`：`1=xtrabackup`，`2=LVM快照`
- `backupConfigData.backupGran`：`1=实例`，`2=库`
- `backupConfigData.backupMode`：`1=仅从备节点`，`2=仅从主节点`，`3=优先从节点`，`4=优先主节点`，`5=任意节点`
- `backupConfigData.deleteArchUnit`：`1=年`，`2=月`，`3=周`，`4=天`，`5=小时`

## 执行流程

1. 执行 `foundation-cli policy backup create`，从返回中提取 `responseData.id` 作为 `slaId`。
2. 将 `slaId` 写入 `bindData.slaIds=[slaId]`，执行 `foundation-cli protect policy bind`。
3. 执行 `foundation-cli mysql backup-config set` 完成 MySQL 备份配置下发。

## 回滚策略（全回滚）

1. 第 2 步失败：删除第 1 步创建的策略（按策略 ID 清理）。
2. 第 3 步失败：先解绑第 2 步策略，再删除第 1 步策略。
3. 若清理动作失败，返回中必须标记 `cleanupFailed=true` 并附带失败详情。

## 请求体示例

```json
{
  "policyData": {
    "name": "mysql-backup-plan-demo",
    "remark": "",
    "validatePeriod": 1,
    "effectiveType": 1,
    "effectiveTime": 0,
    "backupConfig": {
      "dayEnable": true,
      "dayConfiguration": [
        {
          "interval": 1,
          "triggerTime": {
            "hour": 1,
            "min": 0
          },
          "mode": 1
        }
      ],
      "weekEnable": false,
      "weekConfiguration": [],
      "monthEnable": false,
      "monthConfiguration": [],
      "isRetentionOpen": false,
      "logRetentionOpen": false,
      "isGfs": false
    }
  },
  "bindData": {
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
    "slaIds": []
  },
  "backupConfigData": {
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
    "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
    "commonConfigParams": {
      "appType": "eso_backupengine_hypermysqlengine",
      "backupDestination": {
        "multipleRegions": false,
        "regionParams": [
          {
            "storageAvailable": false,
            "storagePoolParams": [
              {
                "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
                "storagePoolName": "block",
                "storagePoolType": 2,
                "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
                "storageServiceName": "StorageService_172.31.12.91"
              }
            ]
          }
        ]
      }
    },
    "backupType": 1,
    "backupGran": 1
  }
}
```

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "policyCreate": {
      "status": "success",
      "slaId": "a15e75fbdab04effbc7bf076aac8c7a8"
    },
    "policyBind": {
      "status": "success"
    },
    "backupConfigSet": {
      "status": "success"
    },
    "cleanup": {
      "triggered": false,
      "cleanupFailed": false,
      "details": []
    }
  }
}
```

## 示例

以下为“等价可执行命令链”，请按顺序执行：

```bash
foundation-cli policy backup create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --data '<policyData-json>'
```

```bash
foundation-cli protect policy bind \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id-from-bindData> \
  --data '{"slaIds":["<slaId-from-step1>"]}'
```

```bash
foundation-cli mysql backup-config set \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --body-file <mysql-backup-config-json-file>
```
