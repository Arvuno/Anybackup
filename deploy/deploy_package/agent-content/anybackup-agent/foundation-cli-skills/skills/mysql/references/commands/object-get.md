# `foundation-cli mysql object get`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql object get --tenant-id <tenant-id> --object-id <object-id>` |
| 方法 | `GET` |
| 路径 | `/backupmgm/v1/mysql/object?objectId=<object-id>` |
| 风险 | `只读` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 基础地址 |
| ak | `--ak` | string | 是 | 访问密钥 |
| sk | `--sk` | string | 是 | 访问密钥密文 |
| targetVersion | `--target-version` | string | 否 | 默认值为 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| objectId | `--object-id` | string | 是 | MySQL 保护对象 ID |

## 响应参数

### 顶层结构

| 字段 | 类型 | 说明 |
|---|---|---|
| status | string | 接口执行状态 |
| error | object/null | 错误信息 |
| responseData | object | MySQL 保护对象详情 |

### `responseData` 核心字段

| 字段 | 类型 | 说明 |
|---|---|---|
| systemId | string | 生产系统 ID |
| objectId | string | 保护对象 ID |
| name | string | 对象名称 |
| businessName | string | 业务名称 |
| platName | string | 平台名称 |
| clusterMode | int32 | 集群模式 |
| protectStatus | int32 | 保护状态 |
| execStatus | int32[] | 执行状态码列表 |
| execStatuses | object[] | 扩展执行状态详情 |
| bindStrategyStatus | int32 | 策略绑定状态 |
| backupStrategyId | string | 备份策略 ID |
| backupStrategyName | string | 备份策略名称 |
| latestTimePoint | string | 最近时间点 |
| lastBackupResult | int32 | 最近备份结果 |
| lastBackupTaskId | string | 最近备份任务 ID |
| hasBackupConfig | bool | 是否存在备份配置 |
| hasAppBackupConfig | bool | 是否存在应用级备份配置 |
| hasBackupData | bool | 是否存在备份数据 |
| objectStatus | int32 | 对象状态 |
| disableReason | string | 禁止备份原因 |
| genMode | int32 | 资源来源模式 |
| nextBackupTime | int64 | 下次备份时间 |
| backupConfig | object | 备份配置摘要 |
| snapshotConfig | object | 快照配置摘要 |
| tenantId | string | 租户 ID |
| tenantName | string | 租户名称 |
| subObjectNumber | int32 | 子对象数量 |
| instanceInfo | string | 实例信息 |
| clientInfo | object[] | 关联客户端记录 |
| appExtend | object[] | 应用扩展字段 |
| creator | string | 创建人 |
| hasPassword | bool | 是否存在恢复校验密码 |

### `execStatuses[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| status | int32 | 执行状态码 |
| latestTaskId | string | 最近任务 ID |

### `clientInfo[]`

| 字段 | 类型 | 说明 |
|---|---|---|
| clientId | string | 客户端 ID |
| clientName | string | 客户端名称 |
| clientIp | string | 客户端 IP |
| clientMac | string | 客户端 MAC |
| clientOSType | int64 | 客户端操作系统类型 |
| clientType | int32 | 客户端类型 |
| clientStatus | int64 | 客户端状态 |

## 返回案例

说明：以下为已采集的成功响应示例。

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9",
    "name": "mysql3306_U0HYTDM3RENXS4EJ",
    "businessName": "mysql3306_U0HYTDM3RENXS4EJ",
    "platName": "iv-yej4xkshkwwh2yram2x1(172.31.0.2)",
    "clusterMode": 3,
    "protectStatus": 2,
    "execStatus": null,
    "execStatuses": null,
    "bindStrategyStatus": 2,
    "backupStrategyId": "",
    "backupStrategyName": "",
    "backupStrategyNum": 0,
    "isGroupSLA": false,
    "type": 202,
    "latestTimePoint": "2026-04-10T14:38:44.000000+08:00",
    "lastBackupDataSize": 19901940,
    "lastRunTime": 1775803115000,
    "lastBackupResult": 800,
    "hasBackupConfig": true,
    "hasAppBackupConfig": false,
    "hasBackupData": true,
    "objectStatus": 4,
    "disableReason": "There is no backup configuration for the current protected object.",
    "genMode": 0,
    "isDeleted": false,
    "backupMode": 0,
    "nextBackupTime": -1,
    "backupConfig": null,
    "snapshotConfig": null,
    "fingerLibraryId": "",
    "fingerPrintLibrary": null,
    "backupSpeed": "",
    "recoverySpeed": "",
    "lastSnapshotTime": 0,
    "lastSnapshotSize": 0,
    "lastSnapshotStatus": 0,
    "lastSnapshotTaskId": "",
    "lastSuccessSyncTime": 1775885977280,
    "tenantId": "",
    "tenantName": "",
    "subObjectNumber": 0,
    "resourcestatus": 0,
    "instanceInfo": "",
    "backupDestination": null,
    "conventionalConfig": null,
    "appExtend": null,
    "clientInfo": [
      {
        "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
        "clientName": "iv-yej4xkshkwwh2yram2x1",
        "clientIp": "172.31.0.2",
        "clientMac": "U0HYTDM3RENXS4EJ",
        "clientOSType": 2,
        "clientType": 1,
        "clientStatus": 0
      }
    ],
    "creator": "mysql"
  }
}
```

## 示例
```bash
foundation-cli mysql object get \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --object-id <object-id>
```



