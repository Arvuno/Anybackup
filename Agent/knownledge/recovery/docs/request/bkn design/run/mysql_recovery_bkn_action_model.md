# MySQL 数据库恢复知识网络 - Action 数据模型设计

**版本**：v1.0  
**文档类型**：行动模型设计文档  
**编写作者**：Andrew

---

## 1. 行动模型概述

本文档定义了 MySQL 数据库恢复知识网络（BKN）中 Action（行动）的数据模型设计，包括对象操作工具的定义、参数配置、执行流程等。

**行动类型说明**：
- Action 定义了对 Object 对象执行的操作
- 包括查询、创建、修改、删除、执行等操作类型
- 每个行动关联工具、参数、前置条件、影响范围等

**关键说明**：
- 查询时间点副本列表需要提供两个参数：objectId、storagePoolId
- 数据保护对象在备份时会明确备份到某个存储服务的某个存储池中
- 通过存储池可以间接关联到存储服务，因此只需要提供 storagePoolId

---

## 2. 数据保护对象操作（ProtectionObject Actions）

### 2.1 查询数据库保护对象

#### 2.1.1 行动定义

**行动ID**：`query_protection_objects`  
**行动名称**：查询数据库保护对象  
**行动描述**：查询 Foundation 接入的所有数据库保护对象，支持按应用类型过滤。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 2.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| protection_object | query |

#### 2.1.3 触发条件

**触发场景**：
- 用户询问有哪些数据库保护对象
- 用户需要查看特定类型的数据库保护对象
- 恢复任务需要获取可用的数据保护对象列表

**触发条件说明**：
- 无特定触发条件，可随时调用
- 支持按 appType 参数过滤查询结果

#### 2.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| protection_object | data_source | exist | Foundation 数据源必须存在 |

#### 2.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| protection_object | 仅查询操作，不影响数据 |

#### 2.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_resource_center_db |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_resource_center_db（查询资源中心数据库）
- **工具功能**：查询 Foundation 接入的所有数据库保护对象
- **数据源**：resource_center_db（资源中心数据库）

#### 2.1.7 参数绑定

| Parameter | Source | Binding | Description |
|-----------|--------|---------|-------------|
| appType | input | - | 应用类型（mysql/oracle/sqlserver），为空则查询所有类型 |

**参数详细说明**：

**appType（应用类型）**：
- **类型**：string
- **来源**：input（用户输入）
- **是否必填**：否
- **可选值**：
  - `backupengine_mysql_...`：MySQL 数据库
  - `backupengine_oracle_...`：Oracle 数据库
  - `backupengine_sqlserver_...`：SQL Server 数据库
  - 空：查询所有类型的数据库保护对象
- **默认值**：空（查询所有）
- **示例**：
  - 查询 MySQL：`backupengine_mysql_8.0`
  - 查询 Oracle：`backupengine_oracle_19c`
  - 查询所有：不传参数或传空字符串

#### 2.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 appType 是否有效（如果提供）
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回数据库保护对象列表
   - 包含：objectId, name, appType, business_name, online_status, authorization_status 等
```

#### 2.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "objectId": "a1b2c3d4e5f67890a1b2c3d4e5f67890",
      "name": "MySQL-3306",
      "appType": "backupengine_mysql_8.0",
      "business_name": "生产订单数据库",
      "online_status": "running",
      "authorization_status": "authorized",
      "system_username": "mysql",
      "detail": "{}"
    },
    {
      "objectId": "b2c3d4e5f6789012b3c4d5e6f7890123",
      "name": "Oracle-1521",
      "appType": "backupengine_oracle_19c",
      "business_name": "财务系统数据库",
      "online_status": "running",
      "authorization_status": "authorized",
      "system_username": "oracle",
      "detail": "{}"
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 数据库保护对象列表 |
| total | integer | 总数量 |

---

## 3. 存储服务对象操作（StorageService Actions）

### 3.1 查询存储服务列表

#### 3.1.1 行动定义

**行动ID**：`query_storage_services`  
**行动名称**：查询存储服务列表  
**行动描述**：查询所有存储服务，支持通过 onlyStorage 参数过滤返回结果。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 3.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| storage_service | query |

#### 3.1.3 触发条件

**触发场景**：
- 用户询问有哪些存储服务
- 恢复任务需要选择存储服务
- 需要查看存储服务列表
- 需要查看管理控制台和存储服务

**触发条件说明**：
- 无特定触发条件，可随时调用
- 支持通过 onlyStorage 参数过滤返回结果
  - `true`（默认）：仅返回存储服务
  - `false`：返回管理控制台和存储服务

#### 3.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| storage_service | data_source | exist | 存储服务数据源必须存在 |

#### 3.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| storage_service | 仅查询操作，不影响数据 |

#### 3.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_storage_services |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_storage_services（查询存储服务）
- **工具功能**：查询所有存储服务，支持通过 onlyStorage 参数过滤
- **数据源**：storageService（存储服务表）

#### 3.1.7 参数绑定

| Parameter | Source | Binding | Description |
|-----------|--------|---------|-------------|
| onlyStorage | input | - | 是否仅返回存储服务（true=仅存储服务，false=管理控制台和存储服务） |

**参数详细说明**：

**onlyStorage（是否仅返回存储服务）**：
- **类型**：boolean
- **来源**：input（用户输入）
- **是否必填**：否
- **默认值**：true
- **可选值**：
  - `true`：仅返回存储服务
  - `false`：返回管理控制台和存储服务
- **说明**：控制查询结果的过滤条件，决定是否包含管理控制台
- **示例**：
  - 仅查询存储服务：`true` 或不传参数
  - 查询管理控制台和存储服务：`false`

#### 3.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 onlyStorage 参数是否有效（如果提供）
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回存储服务列表
   - 包含：storageServiceId, storageServiceName
```

#### 3.1.9 输出结果

**返回数据结构**：

**场景1：onlyStorage = true（默认）- 仅返回存储服务**

```json
{
  "success": true,
  "data": [
    {
      "storageServiceId": "s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890",
      "storageServiceName": "Storage-Service-01"
    },
    {
      "storageServiceId": "s2b3c4d5e6f78901s2b3c4d5e6f78901",
      "storageServiceName": "Storage-Service-02"
    }
  ],
  "total": 2
}
```

**场景2：onlyStorage = false - 返回管理控制台和存储服务**

```json
{
  "success": true,
  "data": [
    {
      "storageServiceId": "m1a2b3c4d5e6f7890m1a2b3c4d5e6f7890",
      "storageServiceName": "Management-Console-01"
    },
    {
      "storageServiceId": "s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890",
      "storageServiceName": "Storage-Service-01"
    },
    {
      "storageServiceId": "s2b3c4d5e6f78901s2b3c4d5e6f78901",
      "storageServiceName": "Storage-Service-02"
    }
  ],
  "total": 3
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 存储服务列表 |
| total | integer | 总数量 |
| storageServiceId | string | 存储服务唯一标识符 |
| storageServiceName | string | 存储服务名称 |

---

## 4. 存储池对象操作（StoragePool Actions）

### 4.1 查询存储池列表

#### 4.1.1 行动定义

**行动ID**：`query_storage_pools`  
**行动名称**：查询存储池列表  
**行动描述**：查询指定存储服务下的所有存储池。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 4.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| storage_pool | query |

#### 4.1.3 触发条件

**触发场景**：
- 用户询问某个存储服务下有哪些存储池
- 恢复任务需要选择存储池
- 需要查看存储池列表
- 需要确认存储池所属的存储服务

**触发条件说明**：
- 需要提供 storageServiceId 参数
- storageServiceId 必须是有效的存储服务ID

#### 4.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| storage_service | exist | true | 存储服务必须存在 |
| storage_service | property:storageServiceId | valid | storageServiceId 必须是有效的存储服务ID |
| storage_pool | data_source | exist | 存储池数据源必须存在 |

#### 4.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| storage_pool | 仅查询操作，不影响数据 |
| storage_service | 仅查询操作，不影响数据 |

#### 4.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_storage_pools |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_storage_pools（查询存储池）
- **工具功能**：查询指定存储服务下的所有存储池
- **数据源**：storagePool（存储池表）

#### 4.1.7 参数绑定

| Parameter | Source | Binding | Description |
|-----------|--------|---------|-------------|
| storageServiceId | input | - | 存储服务ID |

**参数详细说明**：

**storageServiceId（存储服务ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定要查询存储池的存储服务ID
- **示例**：`s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890`

#### 4.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 storageServiceId 是否存在
   - 检查 storageServiceId 是否有效
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回存储池列表
   - 包含：storagePoolId, storagePoolName, storageServiceId
```

#### 4.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "storagePoolId": "p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890",
      "storagePoolName": "Backup-Pool-01",
      "storageServiceId": "s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890"
    },
    {
      "storagePoolId": "p2b3c4d5e6f78901p2b3c4d5e6f78901",
      "storagePoolName": "Backup-Pool-02",
      "storageServiceId": "s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890"
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 存储池列表 |
| total | integer | 总数量 |
| storagePoolId | string | 存储池唯一标识符 |
| storagePoolName | string | 存储池名称 |
| storageServiceId | string | 所属存储服务ID |

---

## 5. 时间点副本对象操作（TimePointCopy Actions）

### 5.1 查询时间点副本列表

#### 5.1.1 行动定义

**行动ID**：`query_timepoint_copies`  
**行动名称**：查询时间点副本列表  
**行动描述**：查询指定数据库保护对象下的备份时间点副本列表，支持 MySQL、Oracle、SQL Server 等多种数据库类型。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 5.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| timepoint_copy | query |

#### 5.1.3 触发条件

**触发场景**：
- 用户询问某个数据库保护对象有哪些备份时间点
- 恢复任务需要选择恢复时间点
- 需要查看备份历史记录
- 需要确认某个数据库对象的备份情况

**触发条件说明**：
- 需要提供 objectId、storagePoolId 两个参数
- objectId 必须是有效的数据库保护对象ID
- storagePoolId 必须是有效的存储池ID

#### 5.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| protection_object | exist | true | 数据库保护对象必须存在 |
| protection_object | property:objectId | valid | objectId 必须是有效的数据库保护对象ID |
| storage_pool | exist | true | 存储池必须存在 |
| storage_pool | property:storagePoolId | valid | storagePoolId 必须是有效的存储池ID |

#### 5.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| timepoint_copy | 仅查询操作，不影响数据 |
| protection_object | 仅查询操作，不影响数据 |
| storage_pool | 仅查询操作，不影响数据 |

#### 5.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_timepoint_copies |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_timepoint_copies（查询时间点副本）
- **工具功能**：查询指定数据库保护对象下的备份时间点副本列表
- **数据源**：backup_catalog_db（备份目录数据库）

#### 5.1.7 参数绑定

| Parameter | Source | Binding | Description |
|-----------|--------|---------|-------------|
| objectId | input | - | 数据库保护对象ID |
| storagePoolId | input | - | 存储池ID |

**参数详细说明**：

**objectId（数据库保护对象ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定要查询时间点副本的数据库保护对象ID
- **示例**：
  - MySQL 对象：`a1b2c3d4e5f67890a1b2c3d4e5f67890`
  - Oracle 对象：`b2c3d4e5f6789012b3c4d5e6f7890123`
  - SQL Server 对象：`c3d4e5f678901234c5d6e7f8g9012345`

**storageServiceId（存储服务ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定存储服务ID，数据保护对象的备份数据存储在该存储服务下
- **示例**：`s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890`

**storagePoolId（存储池ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定存储池ID，数据保护对象的备份数据存储在该存储服务下的该存储池中
- **示例**：`p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890`

#### 5.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 objectId 是否存在
   - 检查 objectId 是否有效
   - 检查 storagePoolId 是否存在
   - 检查 storagePoolId 是否有效
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回时间点副本列表
   - 包含：timePointId, objectId, backup_time, backup_type, size, status 等
```

#### 5.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "timePointId": "tp_001_a1b2c3d4e5f67890",
      "objectId": "a1b2c3d4e5f67890a1b2c3d4e5f67890",
      "backup_time": "2024-01-15 02:00:00",
      "backup_type": "full",
      "size": "1024000000",
      "status": "completed",
      "retention_days": 30,
      "location": "/backup/mysql/3306/20240115_020000"
    },
    {
      "timePointId": "tp_002_a1b2c3d4e5f67890",
      "objectId": "a1b2c3d4e5f67890a1b2c3d4e5f67890",
      "backup_time": "2024-01-14 02:00:00",
      "backup_type": "incremental",
      "size": "51200000",
      "status": "completed",
      "retention_days": 30,
      "location": "/backup/mysql/3306/20240114_020000"
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 时间点副本列表 |
| total | integer | 总数量 |
| timePointId | string | 时间点副本ID |
| objectId | string | 数据库保护对象ID |
| backup_time | string | 备份时间 |
| backup_type | string | 备份类型（full/incremental/binlog） |
| size | string | 备份大小（字节） |
| status | string | 备份状态（completed/failed/running） |
| retention_days | integer | 保留天数 |
| location | string | 备份存储位置 |

---

## 6. 恢复数据源对象操作（RecoveryDataSource Actions）

### 6.1 查询恢复数据源列表

#### 6.1.1 行动定义

**行动ID**：`query_recovery_datasources`  
**行动名称**：查询恢复数据源列表  
**行动描述**：查询指定时间点副本下的恢复数据源列表。业务语义上，该动作面向“已选定时间点副本”的恢复数据源查询；当前底层工具实现需要提供存储服务ID、数据集ID和时间戳作为参数。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 6.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| recovery_datasource | query |

#### 6.1.3 触发条件

**触发场景**：
- 用户询问某个时间点副本有哪些恢复数据源
- 恢复任务需要选择恢复数据源
- 需要查看恢复数据源列表
- 需要确认恢复数据源的详细信息

**触发条件说明**：
- 业务上需要先明确选定的时间点副本
- 当前工具实现需要提供 storageServiceId、dataSetId、timestamp 三个参数
- storageServiceId 必须是有效的存储服务ID
- dataSetId 必须是有效的数据集ID
- timestamp 必须是有效的时间戳

#### 6.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| storage_service | exist | true | 存储服务必须存在 |
| storage_service | property:storageServiceId | valid | storageServiceId 必须是有效的存储服务ID |

#### 6.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| recovery_datasource | 仅查询操作，不影响数据 |
| storage_service | 仅查询操作，不影响数据 |

#### 6.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_recovery_datasources |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_recovery_datasources（查询恢复数据源）
- **工具功能**：查询指定时间点副本下的恢复数据源列表
- **数据源**：recovery_datasource（恢复数据源）

#### 6.1.7 参数绑定

| Parameter | Type | Source | Binding | Description |
|-----------|------|--------|---------|-------------|
| storageServiceId | string | input | - | 存储服务ID（必填） |
| dataSetId | string | input | - | 数据集ID（必填） |
| timestamp | integer | input | - | 时间戳（必填） |

**参数详细说明**：

**storageServiceId（存储服务ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：当前工具实现所需的存储服务ID，对应已选时间点副本所属的存储服务
- **示例**：`s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890`

**dataSetId（数据集ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：当前工具实现所需的数据集ID，对应已选时间点副本对象所属的数据集
- **示例**：`ds1a2b3c4d5e6f7890ds1a2b3c4d5e6f7890`

**timestamp（时间戳）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：当前工具实现所需的时间戳，对应已选时间点副本对象所对应的时间戳（毫秒）
- **示例**：`1704067200000`

#### 6.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 storageServiceId 是否存在、是否有效
   - 检查 dataSetId 是否存在、是否有效
   - 检查 timestamp 是否有效
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回恢复数据源列表
   - 包含：recovery_datasourceId, display, fullPath, metadata, objectSize, timePointId
```

#### 6.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "recovery_datasourceId": "c3d4e5f678901234c5d6e7f890123456",
      "display": "database_backup_001",
      "fullPath": "/backup/mysql/2024/01/database_backup.sql",
      "metadata": "{\"format\": \"sql\", \"compression\": \"gzip\"}",
      "objectSize": 104857600,
      "timePointId": "b2c3d4e5f6789012b3c4d5e6f7890123"
    },
    {
      "recovery_datasourceId": "d4e5f67890123456d7e8f901234567890",
      "display": "database_backup_002",
      "fullPath": "/backup/mysql/2024/01/database_backup_002.sql",
      "metadata": "{\"format\": \"sql\", \"compression\": \"gzip\"}",
      "objectSize": 52428800,
      "timePointId": "b2c3d4e5f6789012b3c4d5e6f7890123"
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 恢复数据源列表 |
| total | integer | 总数量 |
| recovery_datasourceId | string | 恢复数据源唯一标识符 |
| display | string | 展示名称 |
| fullPath | string | 全路径 |
| metadata | string | 对象元数据（JSON字符串） |
| objectSize | integer | 对象大小（字节） |
| timePointId | string | 关联的时间点副本ID |

---

## 7. 客户端对象操作（Client Actions）

### 7.1 查询客户端列表

#### 7.1.1 行动定义

**行动ID**：`query_clients`  
**行动名称**：查询客户端列表  
**行动描述**：查询 Foundation 接入的所有客户端列表，可根据客户端状态、操作系统类型和客户端类型进行筛选。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 7.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| client | query |

#### 7.1.3 触发条件

**触发场景**：
- 用户询问有哪些可用的客户端
- 恢复任务需要选择目标客户端
- 需要查看客户端列表
- 需要确认客户端的在线状态

**触发条件说明**：
- 可选提供 status、platform、clientType 参数进行筛选
- 如果不提供参数，则返回所有客户端

#### 7.1.4 前置条件

无特殊前置条件。

#### 7.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| client | 仅查询操作，不影响数据 |

#### 7.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_clients |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_clients（查询客户端列表）
- **工具功能**：查询 Foundation 接入的所有客户端列表
- **数据源**：client（客户端）

#### 7.1.7 参数绑定

| Parameter | Type | Source | Binding | Description |
|-----------|------|--------|---------|-------------|
| status | integer | input | - | 客户端状态（可选，默认为1在线） |
| type | integer | input | - | 操作系统类型（可选，1=Windows，2=Linux，3=Aix） |
| clientType | integer | input | - | 客户端类型（可选，1=物理客户端，2=虚拟客户端，3=其他类型） |

**参数详细说明**：

**status（客户端状态）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **默认值**：1（在线）
- **说明**：筛选客户端的在线或离线状态
- **枚举值**：
  - `0`：离线
  - `1`：在线
- **示例**：`1`

**type（操作系统类型）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **说明**：筛选客户端所属的操作系统类型
- **枚举值**：
  - `1`：Windows
  - `2`：Linux
  - `3`：Aix
- **示例**：`2`

**clientType（客户端类型）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **说明**：筛选客户端是虚拟客户端还是物理客户端或其他类型
- **枚举值**：
  - `1`：物理客户端
  - `2`：虚拟客户端
  - `3`：其他类型
- **示例**：`1`

#### 7.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 status 是否有效（如果提供）
   - 检查 type 是否有效（如果提供）
   - 检查 clientType 是否有效（如果提供）
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回客户端列表
   - 包含：clientId, clientName, clientType, clientIp, clientMac, type, platform, status, clientVersion
```

#### 7.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "clientId": "d4e5f67890123456d7e8f901234567890",
      "clientName": "MySQL-Client-01",
      "clientType": 1,
      "clientIp": "192.168.1.100",
      "clientMac": "00:1A:2B:3C:4D:5E",
      "type": 2,
      "platform": "Linux",
      "status": 1,
      "clientVersion": "8.0.8.0.318"
    },
    {
      "clientId": "e5f6789012345678e9f0a123456789012",
      "clientName": "MySQL-Client-02",
      "clientType": 2,
      "clientIp": "192.168.1.101",
      "clientMac": "00:1A:2B:3C:4D:5F",
      "type": 1,
      "platform": "Windows",
      "status": 1,
      "clientVersion": "8.0.8.0.318"
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 客户端列表 |
| total | integer | 总数量 |
| clientId | string | 客户端唯一标识符 |
| clientName | string | 客户端名称 |
| clientType | integer | 客户端类型（1=物理客户端，2=虚拟客户端，3=其他类型） |
| clientIp | string | 客户端IP地址 |
| clientMac | string | 客户端MAC地址 |
| type | integer | 操作系统类型（1=Windows，2=Linux，3=Aix） |
| platform | string | 平台类型 |
| status | integer | 客户端状态（0=离线，1=在线） |
| clientVersion | string | 客户端发布版本 |

---

## 8. 恢复任务对象操作（RecoveryTask Actions）

### 8.1 发起恢复任务

#### 8.1.1 行动定义

**行动ID**：`execute_recovery_task`  
**行动名称**：发起恢复任务  
**行动描述**：发起MySQL数据库恢复任务，执行恢复操作。恢复任务执行完成后（无论成功与否），会生成恢复作业记录。  
**风险等级**：high  
**需要审批**：true  
**启用状态**：true

#### 8.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| recovery_task | execute |

**绑定对象说明**：
- `recovery_task` 是恢复动作的参数契约对象
- 该对象不要求映射 Foundation 稳定实体表，而是用于承载恢复执行所需的输入参数

#### 8.1.3 触发条件

**触发场景**：
- 用户确认恢复策略后发起恢复任务
- 需要执行MySQL数据库恢复操作
- 需要将备份数据恢复到目标客户端

**触发条件说明**：
- 必须提供 clientId、clientIp、timePointId、timestamp、storageServiceId、storagePoolId 参数
- clientId 必须是有效的客户端ID
- clientIp 必须是有效的IP地址
- timePointId 必须是有效的时间点副本ID
- timestamp 必须是有效的时间戳
- storageServiceId 必须是有效的存储服务ID
- storagePoolId 必须是有效的存储池ID

#### 8.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| client | exist | true | 客户端必须存在 |
| client | property:status | 1 | 客户端必须在线 |
| timepoint_copy | exist | true | 时间点副本必须存在 |
| storage_service | exist | true | 存储服务必须存在 |
| storage_pool | exist | true | 存储池必须存在 |

#### 8.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| recovery_task | 创建恢复任务记录 |
| recovery_job | 恢复任务执行完成后生成恢复作业记录 |
| client | 在目标客户端上执行恢复操作，可能影响客户端上的数据库 |

#### 8.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | execute_recovery_task |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：execute_recovery_task（发起恢复任务）
- **工具功能**：发起MySQL数据库恢复任务，执行恢复操作
- **数据源**：无（恢复任务对象无数据源，通过工具即时创建）

#### 8.1.7 参数绑定

| Parameter | Type | Source | Binding | Description |
|-----------|------|--------|---------|-------------|
| clientId | string | input | - | 恢复目标客户端ID（必填） |
| clientIp | string | input | - | 恢复目标客户端IP（必填） |
| timePointId | string | input | - | 时间点副本ID（必填） |
| timestamp | integer | input | - | 选择的时间点的时间戳（毫秒，必填） |
| storageServiceId | string | input | - | 时间点所属的存储服务ID（必填） |
| storagePoolId | string | input | - | 时间点所属的存储池ID（必填） |

**参数详细说明**：

**clientId（客户端ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定恢复到哪一个ID的客户端
- **示例**：`d4e5f67890123456d7e8f901234567890`

**clientIp（客户端IP）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定恢复到哪一个IP的客户端
- **示例**：`192.168.1.100`

**timePointId（时间点副本ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定要恢复的时间点副本ID
- **示例**：`b2c3d4e5f6789012b3c4d5e6f7890123`

**timestamp（时间戳）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：选择的时间点的时间戳（毫秒）
- **示例**：`1704067200000`

**storageServiceId（存储服务ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：选择的时间点所属的存储服务ID
- **示例**：`s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890`

**storagePoolId（存储池ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：选择的时间点所属的存储服务下的存储池ID
- **示例**：`p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890`

#### 8.1.8 风险检查

**风险检查规则**：`recovery_overwrite_risk`

**检查目的**：判断恢复操作是否会覆盖原生产环境数据

**检查逻辑**：
```
1. 获取恢复目标客户端ID（targetClientId = 参数中的 clientId）
2. 通过 targetClientId 查询恢复目标客户端IP（targetClientIp）
3. 查询时间点副本关联的数据保护对象
4. 获取数据保护对象的 platName，并解析 sourceClientIp
5. 比较 targetClientIp 与 sourceClientIp
```

**风险评估结果**：

| 条件 | 风险等级 | 处理方式 | 说明 |
|------|---------|---------|------|
| sourceClientIp 无法解析 | high | 需要用户确认 | 无法确认备份源客户端，按高风险处理 |
| targetClientIp == sourceClientIp | high | 需要用户确认 | 原机恢复，会覆盖生产环境数据 |
| targetClientIp != sourceClientIp | low | 可自主执行 | 异机恢复，不影响生产环境 |

**高风险提示**：
```
⚠️ 原机恢复风险警告

检测到您正在执行原机恢复操作：
- 恢复目标客户端：{targetClientName}（{targetClientIp}）
- 备份源客户端：{sourceClientName}（{sourceClientIp}）
- 时间点副本：{timePointName}

此操作将覆盖原生产环境的数据，可能导致：
1. 当前生产数据被覆盖
2. 正在运行的业务中断
3. 数据不可恢复

请确认是否继续执行恢复操作？
```

**详细设计参考**：[mysql_recovery_bkn_risk_model.md](./mysql_recovery_bkn_risk_model.md)

#### 8.1.9 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 clientId 是否存在、是否有效
   - 检查 clientIp 是否有效
   - 检查 timePointId 是否存在、是否有效
   - 检查 timestamp 是否有效
   - 检查 storageServiceId 是否存在、是否有效
   - 检查 storagePoolId 是否存在、是否有效
   ↓
3. 检查前置条件
   - 检查客户端是否在线
   - 检查时间点副本是否存在
   - 检查存储服务和存储池是否存在
   ↓
4. 风险检查
   - 通过 clientId 查询恢复目标客户端IP（targetClientIp）
   - 查询时间点副本关联的数据保护对象
   - 获取 platName 并解析备份源客户端IP（sourceClientIp）
   - 比较恢复目标客户端IP（targetClientIp）与备份源客户端IP
   - 判断风险等级
   ↓
5. 风险处理
   - 高风险（原机恢复）：提示用户确认
     - 用户确认 → 继续执行
     - 用户取消 → 终止恢复
   - 低风险（异机恢复）：直接执行
   ↓
6. 调用 Foundation 工具执行恢复任务
   ↓
7. 返回结果
   - 返回恢复任务ID
   - 返回恢复任务执行状态
```

#### 8.1.10 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": {
    "recovery_taskId": "e5f6789012345678e9f0a123456789012",
    "clientId": "d4e5f67890123456d7e8f901234567890",
    "clientIp": "192.168.1.100",
    "timePointId": "b2c3d4e5f6789012b3c4d5e6f7890123",
    "timestamp": 1704067200000,
    "storageServiceId": "s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890",
    "storagePoolId": "p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890",
    "status": "executing",
    "message": "恢复任务已发起，正在执行中"
  }
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 操作是否成功 |
| data | object | 返回数据 |
| recovery_taskId | string | 恢复任务唯一标识符 |
| clientId | string | 客户端ID |
| clientIp | string | 客户端IP |
| timePointId | string | 时间点副本ID |
| timestamp | integer | 时间戳（毫秒） |
| storageServiceId | string | 存储服务ID |
| storagePoolId | string | 存储池ID |
| status | string | 恢复任务状态（executing/completed/failed） |
| message | string | 操作结果消息 |

---

## 9. 恢复作业对象操作（RecoveryJob Actions）

### 9.1 查询恢复作业列表

#### 9.1.1 行动定义

**行动ID**：`query_recovery_jobs`  
**行动名称**：查询恢复作业列表  
**行动描述**：查询某一数据保护对象的恢复作业列表，可根据业务类型进行筛选。  
**风险等级**：low  
**需要审批**：false  
**启用状态**：true

#### 9.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| recovery_job | query |

#### 9.1.3 触发条件

**触发场景**：
- 用户查看某一数据保护对象的恢复作业历史
- 需要了解恢复作业的执行情况
- 需要筛选特定类型的恢复作业

**触发条件说明**：
- 必须提供 objectId 参数
- 可选提供 businessType 参数进行筛选

#### 9.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| protection_object | exist | true | 数据保护对象必须存在 |

#### 9.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| recovery_job | 仅查询操作，不影响数据 |

#### 9.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | query_recovery_jobs |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：query_recovery_jobs（查询恢复作业列表）
- **工具功能**：查询某一数据保护对象的恢复作业列表
- **数据源**：job（恢复作业）

#### 9.1.7 参数绑定

| Parameter | Type | Source | Binding | Description |
|-----------|------|--------|---------|-------------|
| objectId | string | input | - | 数据保护对象ID（必填） |
| businessType | integer | input | - | 业务类型（可选，10=数据恢复，11=挂载恢复，12=细粒度恢复） |

**参数详细说明**：

**objectId（数据保护对象ID）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：指定查询哪一个数据保护对象的恢复作业
- **示例**：`a1b2c3d4e5f67890a1b2c3d4e5f67890`

**businessType（业务类型）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **说明**：用于过滤恢复作业的类型
- **枚举值**：
  - `10`：数据恢复作业
  - `11`：挂载恢复作业
  - `12`：细粒度恢复作业
- **示例**：`10`

#### 9.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查 objectId 是否存在、是否有效
   - 检查 businessType 是否有效（如果提供）
   ↓
3. 调用 Foundation 工具
   ↓
4. 返回结果
   - 返回恢复作业列表
   - 包含：jobId, objectId, clientId, speed, startTime, endTime, completedData, sendSize, businessType
```

#### 9.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": [
    {
      "jobId": "f6789012345678901a234567890123456",
      "objectId": "a1b2c3d4e5f67890a1b2c3d4e5f67890",
      "clientId": "d4e5f67890123456d7e8f901234567890",
      "speed": "100 MB/s",
      "startTime": 1704067200000,
      "endTime": 1704070800000,
      "completedData": 52428800000,
      "sendSize": 52428800000,
      "businessType": 10
    },
    {
      "jobId": "g7890123456789012b345678901234567",
      "objectId": "a1b2c3d4e5f67890a1b2c3d4e5f67890",
      "clientId": "d4e5f67890123456d7e8f901234567890",
      "speed": "80 MB/s",
      "startTime": 1704153600000,
      "endTime": 1704157200000,
      "completedData": 41943040000,
      "sendSize": 41943040000,
      "businessType": 11
    }
  ],
  "total": 2
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 查询是否成功 |
| data | array | 恢复作业列表 |
| total | integer | 总数量 |
| jobId | string | 作业唯一标识符 |
| objectId | string | 数据保护对象ID |
| clientId | string | 客户端ID |
| speed | string | 恢复速度 |
| startTime | integer | 恢复开始时间（时间戳） |
| endTime | integer | 恢复完成时间（时间戳） |
| completedData | integer | 恢复写入的数据量（字节） |
| sendSize | integer | 从服务器存储中读出的数据量（字节） |
| businessType | integer | 业务类型（10=数据恢复，11=挂载恢复，12=细粒度恢复） |

---

## 10. 可用性验证对象操作（AvailabilityVerification Actions）

### 10.1 执行可用性验证

#### 10.1.1 行动定义

**行动ID**：`execute_availability_verification`  
**行动名称**：执行可用性验证  
**行动描述**：执行对某数据库对象进行可用性验证，验证其是否能对外提供业务。通过连接数据库并执行验证操作，判断恢复后的数据库是否可以正常工作。  
**风险等级**：medium  
**需要审批**：true  
**启用状态**：true

#### 10.1.2 绑定对象

| Bound Object | Action Type |
|--------------|-------------|
| availability_verification | execute |

**绑定对象说明**：
- `availability_verification` 是验证执行契约 / 执行模板对象
- 它用于承载验证动作所需的执行信息，不单独定义验证结果对象
- 验证结果由动作输出返回，再由 LLM 在上下文中进行理解与总结

#### 10.1.3 触发条件

**触发场景**：
- 恢复作业执行完成后需要验证数据库可用性
- 需要确认恢复后的数据库是否能正常提供业务
- 需要对数据库进行健康检查

**触发条件说明**：
- 必须提供所有数据库连接参数（dbUser, dbPwd, dbHost, dbPort）
- 必须提供数据库类型（dbType）
- 可选提供配置文件路径和安装路径

#### 10.1.4 前置条件

| Object | Check | Condition | Message |
|--------|-------|-----------|---------|
| client | exist | true | 客户端必须存在 |
| client | property:status | 1 | 客户端必须在线 |

#### 10.1.5 影响范围

| Object | Impact Description |
|--------|---------------------|
| availability_verification | 执行可用性验证契约，不单独创建结果对象 |
| client | 在目标客户端上执行验证操作，可能影响客户端上的数据库 |

#### 10.1.6 工具配置

| Type | Toolbox ID | Tool ID |
|------|------------|---------|
| tool | foundation_toolbox | execute_availability_verification |

**工具说明**：
- **工具箱**：foundation_toolbox（Foundation 工具箱）
- **工具ID**：execute_availability_verification（执行可用性验证）
- **工具功能**：执行数据库可用性验证，验证数据库是否能正常提供业务
- **数据源**：无（通过action参数获取验证信息）

#### 10.1.7 参数绑定

| Parameter | Type | Source | Binding | Description |
|-----------|------|--------|---------|-------------|
| dbUser | string | input | - | 连接要验证的数据库的所需用户名（必填） |
| dbPwd | string | input | - | 连接要验证的数据库的所需密码（必填） |
| dbHost | string | input | - | 连接要验证的数据库安装的主机的IP或主机名（必填） |
| dbPort | integer | input | - | 连接要验证的数据库实例或库对外提供服务的端口（必填） |
| confPath | string | input | - | 数据库的配置文件路径（可选） |
| installPath | string | input | - | 数据库的安装路径（可选） |
| dbType | string | input | - | 数据库类型（必填，如mysql/oracle/sqlserver） |

**参数详细说明**：

**dbUser（数据库用户名）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：连接要验证的数据库所需的用户名
- **示例**：`root`

**dbPwd（数据库密码）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：连接要验证的数据库所需的密码
- **示例**：`password123`

**dbHost（数据库主机）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：连接要验证的数据库安装的主机的IP或主机名
- **示例**：`192.168.1.100` 或 `db-server.example.com`

**dbPort（数据库端口）**：
- **类型**：integer
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：连接要验证的数据库实例或库对外提供服务的端口
- **示例**：`3306`

**confPath（配置文件路径）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **说明**：数据库的配置文件路径
- **示例**：`/etc/mysql/my.cnf`

**installPath（安装路径）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：否
- **说明**：数据库的安装路径
- **示例**：`/usr/local/mysql`

**dbType（数据库类型）**：
- **类型**：string
- **来源**：input（用户输入或上下文获取）
- **是否必填**：是
- **说明**：数据库类型，用于选择对应的验证工具和方法
- **枚举值**：
  - `mysql`：MySQL数据库
  - `oracle`：Oracle数据库
  - `sqlserver`：SQL Server数据库
- **示例**：`mysql`

#### 10.1.8 执行流程

```
1. 接收参数
   ↓
2. 验证参数
   - 检查必填参数是否存在：dbUser, dbPwd, dbHost, dbPort, dbType
   - 检查参数格式是否正确
   ↓
3. 检查前置条件
   - 检查客户端是否在线
   ↓
4. 根据dbType选择验证工具
   - mysql：使用MySQL验证工具
   - oracle：使用Oracle验证工具
   - sqlserver：使用SQL Server验证工具
   ↓
5. 执行可用性验证
   - 连接数据库
   - 执行验证SQL或操作
   - 获取验证结果
   ↓
6. 返回结果
   - 返回验证状态（成功/失败）
   - 返回验证详情
   - 返回错误信息（如果失败）
```

#### 10.1.9 输出结果

**返回数据结构**：

```json
{
  "success": true,
  "data": {
    "verificationId": "v1234567890abcdef1234567890abcdef",
    "dbType": "mysql",
    "dbHost": "192.168.1.100",
    "dbPort": 3306,
    "status": "passed",
    "message": "数据库可用性验证通过，可以正常提供业务",
    "details": {
      "connectionTest": "passed",
      "queryTest": "passed",
      "responseTime": "50ms"
    },
    "timestamp": 1704067200000
  }
}
```

**字段说明**：

| 字段名 | 类型 | 说明 |
|--------|------|------|
| success | boolean | 验证是否成功执行 |
| data | object | 验证结果数据 |
| verificationId | string | 验证记录ID |
| dbType | string | 数据库类型 |
| dbHost | string | 数据库主机 |
| dbPort | integer | 数据库端口 |
| status | string | 验证状态（passed/failed） |
| message | string | 验证结果消息 |
| details | object | 验证详情 |
| connectionTest | string | 连接测试结果 |
| queryTest | string | 查询测试结果 |
| responseTime | string | 响应时间 |
| timestamp | integer | 验证时间戳 |

---

## 11. 工具箱定义

### 11.1 Foundation 工具箱

**工具箱ID**：`foundation_toolbox`  
**工具箱名称**：Foundation 工具箱  
**工具箱描述**：Foundation 备份软件提供的工具集合

**工具列表**：

| Tool ID | Tool Name | Description |
|---------|-----------|-------------|
| query_resource_center_db | 查询资源中心数据库 | 查询 Foundation 接入的所有数据库保护对象 |
| query_storage_services | 查询存储服务 | 查询所有存储服务 |
| query_storage_pools | 查询存储池 | 查询指定存储服务下的所有存储池 |
| query_timepoint_copies | 查询时间点副本 | 查询指定数据库保护对象下的备份时间点副本列表 |
| query_recovery_datasources | 查询恢复数据源 | 查询指定时间点副本下的恢复数据源列表 |
| query_clients | 查询客户端列表 | 查询 Foundation 接入的所有客户端列表 |
| execute_recovery_task | 发起恢复任务 | 发起MySQL数据库恢复任务，执行恢复操作 |
| query_recovery_jobs | 查询恢复作业列表 | 查询某一数据保护对象的恢复作业列表 |
| execute_availability_verification | 执行可用性验证 | 执行数据库可用性验证，验证数据库是否能正常提供业务 |

---

## 12. 设计说明

### 12.1 行动类型分类

**Query（查询）**：
- 风险等级：low
- 需要审批：false
- 影响：只读操作，不影响数据

**Create（创建）**：
- 风险等级：high
- 需要审批：true
- 影响：创建新资源，可能影响系统状态

**Execute（执行）**：
- 风险等级：medium
- 需要审批：true
- 影响：执行操作，可能影响系统状态

### 12.2 参数来源分类

**input（输入）**：
- 来源：用户输入或上下文获取
- 说明：需要用户提供或从上下文中获取

**property（属性）**：
- 来源：对象属性
- 说明：从绑定对象的属性中获取

**const（常量）**：
- 来源：固定值
- 说明：预定义的常量值

### 12.3 风险等级说明

**low（低风险）**：
- 只读操作
- 不影响系统状态
- 可随时执行

**medium（中风险）**：
- 可能影响系统状态
- 需要审批后执行
- 建议在测试环境验证

**high（高风险）**：
- 可能影响业务运行
- 必须审批后执行
- 建议在测试环境充分验证

---

## 13. 后续工作

### 13.1 待设计的行动类型

#### 数据保护对象操作
- [x] 查询数据库保护对象（query_protection_objects）

#### 存储服务操作
- [x] 查询存储服务列表（query_storage_services）

#### 存储池操作
- [x] 查询存储池列表（query_storage_pools）

#### 时间点副本操作
- [x] 查询时间点副本列表（query_timepoint_copies）

#### 恢复数据源操作
- [x] 查询恢复数据源列表

#### 客户端操作
- [x] 查询客户端列表

#### 恢复任务操作
- [x] 发起恢复任务（execute_recovery_task）

#### 恢复作业操作
- [x] 查询恢复作业列表（query_recovery_jobs）

#### 可用性验证操作
- [x] 执行可用性验证（execute_availability_verification）

### 13.2 待编写的 BKN 文件

根据设计文档，后续需要编写以下 action_type BKN 文件：
- query_protection_objects.bkn
- query_storage_services.bkn
- query_storage_pools.bkn
- query_timepoint_copies.bkn
- query_recovery_datasources.bkn
- query_clients.bkn
- execute_recovery_task.bkn
- query_recovery_jobs.bkn
- execute_availability_verification.bkn
- 其他行动类型 BKN 文件（待设计完成后编写）
