# `foundation-cli mysql restore-config create`

## 命令摘要

| 项目 | 值 |
|---|---|
| CLI | `foundation-cli mysql restore-config create --tenant-id <tenant-id> --data '<json>'` |
| 方法 | `POST` |
| 路径 | `/backupmgm/v1/mysql/recovery` |
| 风险 | `写入` |

## 请求参数

### 公共参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| tenantId | `--tenant-id` | string | 否 | 租户标识；可由环境变量 FOUNDATION_TENANT_ID 提供 |
| endpoint | `--endpoint` | string | 是 | Foundation 地址 |
| ak | `--ak` | string | 是 | 访问密钥 AK |
| sk | `--sk` | string | 是 | 访问密钥 SK |
| targetVersion | `--target-version` | string | 否 | 默认 `8.0.9.0` |

### 命令参数

| 字段 | CLI 参数 | 类型 | 必填 | 说明 |
|---|---|---|---|---|
| data | `--data`, `-d` | string | 二选一 | 行内 JSON 请求体 |
| bodyFile | `--body-file` | string | 二选一 | JSON 请求体文件 |

## Body 参数

按恢复流程填写请求体，建议顺序如下。

### 1. 选择存储服务和存储池（先选）

存储服务来源命令：`foundation-cli storage service list`。  
存储池来源命令：`foundation-cli storage pool list --storage-svc-id <storageServiceId>`。  
`storagePoolId` 必须来自所选 `storageServiceId` 下的存储池列表，禁止跨存储服务复用存储池。

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| storageServiceId | string | 是 | 存储服务 ID |
| storageServiceName | string | 否 | 存储服务名称 |
| storagePoolId | string | 是 | 存储池 ID |
| storagePoolName | string | 否 | 存储池名称 |
| storagePoolType | int32 | 是 | 存储池类型 |

### 2. 选择恢复粒度与恢复方式

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| restoreGran | int32 | 是 | 恢复粒度（`1=实例`，`2=库`，`3=表`，`4=日志文件`） |
| restoreType | int32 | 是 | 恢复方式（`0=最短时间恢复`，`1=恢复到最新状态`，`2=恢复到指定时间`） |
| isRestoreNewInstance | bool | 是 | 是否恢复到新实例 |
| isCoverOldDb | bool | 否 | 是否覆盖已有数据库 |
| isShutDown | bool | 否 | 恢复前是否停库 |
| databaseOnline | bool | 否 | 恢复后是否拉起数据库 |

宕机恢复规则：当场景是“实例宕机/故障恢复”时，`isShutDown` 必须传 `true`，且 `isCoverOldDb` 必须传 `true`。

### 3. 选择时间点

时间点必须先查询，再回填恢复参数。  
时间点来源命令：`foundation-cli mysql recovery timepoint list`。
推荐映射关系：`timePointId <- timePointId`，`dataSetId <- dataSetId`，`timestamp <- timestamp`，`displayTime <- displayTime`。

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| timePointId | string | 是 | 时间点 ID |
| dataSetId | string | 是 | 数据集 ID |
| timestamp | int64 | 是 | 时间点时间戳 |
| displayTime | string | 否 | 时间点显示时间 |
| display | int64 | 否 | 时间点显示值 |
| logTimePointId | string | 否 | 日志时间点 ID |
| logSetId | string | 否 | 日志数据集 ID |
| logTimestamp | int64 | 否 | 日志时间戳 |

### 4. 选择恢复的数据（源侧）

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| srcObject.systemId | string | 是 | 源生产系统 ID |
| srcObject.objectId | string | 是 | 源保护对象 ID |
| realPath | string | 是 | 恢复数据路径，JSON 字符串数组 |
| schemaList | string | 否 | 恢复的库/表范围 |

`realPath` 必须有值，且必须是 JSON 字符串数组格式；实例级恢复使用 `["/<实例名>"]`。禁止传 `[]`、空字符串、`null` 或省略该字段。

### 5. 选择目标实例（目标侧）

目标实例必须先查出来，再回填恢复参数。  
目标实例来源接口：`/resource_center/v1/databasealone/instance_list`。  
CLI 对应命令：`foundation-cli mysql target-instance list`。
`clientId` 必须取自 `foundation-cli mysql target-instance list` 返回值，不能手填其他实例的 `clientId`。

| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| targetObject.systemId | string | 是 | 目标系统 ID |
| targetObject.objectId | string | 是 | 目标对象 ID |
| instanceName | string | 是 | 目标实例名称 |
| targetInstancePlatName | string | 否 | 目标实例显示名 |
| clientId | string | 是 | 目标客户端 ID |
| clientIp | string | 否 | 目标客户端 IP |
| clientType | int32 | 否 | 客户端类型 |
| clientOSType | int64 | 否 | 客户端系统类型 |
| osUserName | string | 是 | 操作系统用户 |

### 6. 指定数据路径与执行控制（最后确认）

`dataFilePath` 依赖目标实例。  
建议取目标实例实际数据目录或实例目录，不能随意填写临时目录（如 `/tmp`），否则容易恢复失败。
`dataFilePath` 的数据来源以 `/commons/v1/datasources` 为准，CLI 对应命令：`foundation-cli client datasource list`。
查询 `foundation-cli client datasource list` 时，`--client-id` 必须传上一步 `foundation-cli mysql target-instance list` 返回的同一个 `clientId`。
MySQL 恢复路径查询必须传 `--runner-type MySQL --runner-user root`；禁止使用目标实例返回的 `osUserName` 作为 `runnerUser`。
禁止跨实例复用 `clientId` 与 `dataFilePath`。
#### 路径下钻流程（推荐按此顺序执行）

1. 第一次执行 `foundation-cli client datasource list` 时，使用 `fullPath=""`，且不要传 `requestId`。
   这一步的首要目标是拿到首次响应里的 `responseData.requestId`。
2. 如果第一次返回空，第二次继续使用 `fullPath=""`，并带上 `requestId=<第一次返回的 requestId>`。
3. 如果第二次仍未拿到可直接使用的非根路径，第三次改为：
   `fullPath="/"` + `requestId=<第一次返回的 requestId>`。
4. 如果第三次后已经拿到根目录子项，则继续按固定顺序下钻：
   `fullPath="/var"` -> `fullPath="/var/lib"` -> `fullPath="/var/lib/mysql"`。
5. 一旦某一步返回了可用的非根路径，就从该次响应的 `responseData.data[*].fullPath` 里选择 `dataFilePath`，不要跳步手填猜测值。
#### 结果判定

1. 如果第二次只返回 `/`，说明当前只拿到了根目录入口，还不能直接用于恢复；必须继续查 `/`、`/var`、`/var/lib`、`/var/lib/mysql`。
2. 如果 `/var` 返回空，不代表链路失败；继续查 `/var/lib`。
3. 如果 `/var/lib` 返回多条路径，优先选择最接近 MySQL 数据目录语义的路径，通常是 `/var/lib/mysql`。
4. 如果 `/var/lib/mysql` 再次展开后返回 `/var/lib/mysql/mysql`、`/var/lib/mysql/performance_schema` 之类更深层子目录，`dataFilePath` 仍优先使用上层目录 `/var/lib/mysql`，而不是这些库级子目录。
5. 只有拿到非根路径后，才可以回填 `dataFilePath`；根路径 `/` 只能作为继续下钻的起点，不能直接用于恢复。

#### 兜底与禁忌

1. `dataFilePath` 只能取自 `client datasource list` 或 `mysql datasource recovery` 返回的路径字段，禁止手填猜测值。
2. 禁止将根路径 `/` 直接写入 `dataFilePath` 发起恢复。
3. MySQL 恢复路径查询必须固定 `runnerType=MySQL`、`runnerUser=root`。
4. `dataFilePath` 与 `--client-id` 必须同源：都来自同一次 `mysql target-instance list` 结果，禁止跨实例复用。
5. 若完成上述固定轮询和 `/var -> /var/lib -> /var/lib/mysql` 下钻后，仍没有返回可用的非根路径，MySQL 恢复场景默认 `dataFilePath="/var/lib/mysql"`。
6. `realPath` 必须非空；实例级恢复必须使用 `["/<实例名>"]`，禁止使用 `[]`。
7. 若恢复意图为“宕机恢复 / 故障恢复”，必须设置 `isShutDown=true` 且 `isCoverOldDb=true`。
| 字段 | 类型 | 必填 | 说明 |
|---|---|---|---|
| dataFilePath | string | 是 | 数据文件恢复路径（来源：`foundation-cli client datasource list` 返回的 datasource 路径字段） |
| cnfType | int32 | 是 | 配置文件处理方式 |
| cnfFilePath | string | 否 | 自定义配置文件路径 |
| mysqlDefaultSpfileParam | string | 否 | 默认配置内容 |
| isMultiChannel | bool | 否 | 是否多通道恢复 |
| dataChannelNum | int32 | 否 | 通道数 |
| hasCustomParam | bool | 否 | 是否使用自定义参数 |
| customParam | string | 否 | 自定义参数内容 |
| failureRetry | bool | 否 | 是否失败重试 |
| failureRetryCount | int32 | 否 | 重试次数 |
| failureRetryInterval | int32 | 否 | 重试间隔 |
| isUpdatePassWord | bool | 否 | 是否更新实例密码 |
| specifyRestoreUser | bool | 否 | 是否指定恢复用户 |
| restoreUser | string | 否 | 恢复用户 |
| restorePwd | string | 否 | 恢复密码 |
| instanceUser | string | 否 | 实例用户 |
| instancePwd | string | 否 | 实例密码 |
| runImmediately | bool | 否 | 是否立即执行 |
| objectPassword | string | 否 | 对象校验密码 |

## 枚举列表

- `restoreType`：
- `0`：最短时间恢复
- `1`：恢复到最新状态
- `2`：恢复到指定时间
- `restoreGran`：
- `1`：实例
- `2`：库
- `3`：表
- `4`：日志文件
- `cnfType`：`0/1/2`（原配置/指定配置/自定义配置）。
- 布尔开关字段统一：`true=开启`，`false=关闭`。

## 请求体示例

```json
{
  "storageServiceId": "4235757746e447d18d4d7a8b3f5edc54",
  "storageServiceName": "StorageService_172.31.12.91",
  "storagePoolId": "0c6022822f1d11f19b0600163e33df73",
  "storagePoolName": "block",
  "storagePoolType": 2,
  "restoreGran": 1,
  "restoreType": 1,
  "isRestoreNewInstance": false,
  "isCoverOldDb": true,
  "isShutDown": true,
  "databaseOnline": true,
  "srcObject": {
    "systemId": "a63ba04bfb530fefa3b09d1f89d47432",
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9"
  },
  "targetObject": {
    "systemId": "13caa58fd84bab1beb0b7c9d76b7efae",
    "objectId": "428e4c05bf87deef20a86e670f3eb9a9"
  },
  "instanceName": "mysql3306_U0HYTDM3RENXS4EJ",
  "targetInstancePlatName": "iv-yej4xkshkwwh2yram2x1(172.31.0.2)",
  "clientId": "13caa58fd84bab1beb0b7c9d76b7efae",
  "clientIp": "172.31.0.2",
  "clientType": 1,
  "clientOSType": 2,
  "osUserName": "mysql",
  "timePointId": "f13ed97d3c1d117bb50400163e33df73",
  "dataSetId": "c1ae3176705c538d42fb7b80ae4dd2f4",
  "timestamp": 1776923925509924,
  "display": 1776923934000000,
  "displayTime": "2026-04-23T13:58:54.000000+08:00",
  "logSetId": "",
  "logTimePointId": "",
  "logTimestamp": 0,
  "storagePoolService": "4235757746e447d18d4d7a8b3f5edc54",
  "storagePoolMedia": {
    "type": 2,
    "id": "0c6022822f1d11f19b0600163e33df73",
    "name": "block",
    "serverId": "4235757746e447d18d4d7a8b3f5edc54",
    "serverName": "StorageService_172.31.12.91"
  },
  "dataFilePath": "/var/lib/mysql",
  "realPath": "[\"/mysql3306_U0HYTDM3RENXS4EJ\"]",
  "cnfType": 0,
  "mysqlDefaultSpfileParam": "[{\"type\":\"default\",\"key\":\"server_id\",\"value\":\"9999\"},{\"type\":\"default\",\"key\":\"log_bin\",\"value\":\"mysql-bin\"},{\"type\":\"default\",\"key\":\"binlog_format\",\"value\":\"ROW\"},{\"type\":\"default\",\"key\":\"log-error\",\"value\":\"mysqld.log\"},{\"type\":\"default\",\"key\":\"port\"}]",
  "isUpdatePassWord": false,
  "specifyRestoreUser": false,
  "isMultiChannel": false,
  "dataChannelNum": 0,
  "hasCustomParam": false,
  "failureRetry": false,
  "failureRetryCount": 0,
  "failureRetryInterval": 0
}
```

## 返回案例

```json
{
  "status": "success",
  "error": null,
  "responseData": {
    "taskId": "f138a59a121f1161b50400163e33df73"
  }
}
```

## 示例

```bash
foundation-cli mysql restore-config create \
  --tenant-id <tenant-id> \
  --endpoint <endpoint> \
  --ak <ak> \
  --sk <sk> \
  --body-file restore-body.json
```
