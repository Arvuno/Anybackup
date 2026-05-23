# MySQL 数据库恢复知识网络 - Object 数据模型设计

**版本**：v1.0  
**文档类型**：数据模型设计文档  
**编写作者**：Andrew

---

## 1. 数据模型概述

本文档定义了 MySQL 数据库恢复知识网络（BKN）中 Object（对象）的数据模型设计，包括数据保护对象、存储服务、存储池、时间点副本等对象的结构、字段定义、约束条件等。

**对象总览**：
- **数据保护对象（ProtectionObject）**：Foundation 接入的数据库保护对象
- **存储服务对象（StorageService）**：备份数据存储的存储服务
- **存储池对象（StoragePool）**：存储服务下的存储池
- **时间点副本对象（TimePointCopy）**：数据库备份的时间点副本
- **恢复数据源对象（RecoveryDataSource）**：选定时间点副本下的可恢复数据源
- **客户端对象（Client）**：执行恢复操作的客户端主机
- **恢复任务对象（RecoveryTask）**：恢复动作参数契约对象
- **恢复作业对象（RecoveryJob）**：恢复任务执行产生的作业
- **可用性验证对象（AvailabilityVerification）**：验证执行契约 / 执行模板对象

---

## 2. 数据保护对象（ProtectionObject）

### 2.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `objectId` | `string` | 主键，非自增 | 对象唯一标识符（去掉横杠的 UUID） | `a1b2c3d4e5f67890a1b2c3d4e5f67890` |
| `name` | `string` | 非空 | 对象名称 | `MySQL-3306` |
| `appType` | `string` | 非空 | 应用类型（区分具体应用类型，如 MySQL、Oracle） | `backupengine_mysql_...` |
| `platName` | `string` | 可选 | 平台信息，格式通常为 `clientName(clientIp)`，用于识别备份源客户端 | `localhost.localdomain(192.168.151.13)` |
| `business_name` | `string` | 可选 | 对应的业务名称 | `生产订单数据库` |
| `online_status` | `string` | 可选 | 联机状态 | `running` |
| `authorization_status` | `string` | 可选 | 授权状态 | `authorized` |
| `system_username` | `string` | 可选 | 系统用户名 | `mysql` |
| `detail` | `string` | 可选 | 扩展字段（JSON 字符串格式） | `{"key": "value"}` |

### 2.2 状态枚举

**联机状态（online_status）**：
- `running`：运行中
- `stopped`：已停止
- `error`：错误
- `unknown`：未知

**授权状态（authorization_status）**：
- `authorized`：已授权
- `unauthorized`：未授权
- `expired`：已过期

### 2.3 数据源

| Type | ID | Name |
|------|-----|------|
| data_view | resource_center_db | resource_center_db |

**补充说明**：
- `platName` 是 Foundation 返回的原始字段
- 风险判断时通过解析 `platName` 中括号里的 IP 获取备份源客户端 IP
- `platName` 不作为 relation 映射主键使用

### 2.4 验证规则

1. **唯一性验证**：
   - `objectId` 字段必须唯一

2. **非空验证**：
   - `objectId`、`name`、`appType` 字段不能为空

3. **状态验证**：
   - 当 `online_status` 字段存在时，必须是预定义的枚举值
   - 当 `authorization_status` 字段存在时，必须是预定义的枚举值

4. **格式约束**：
   - `platName`（如果存在）应优先满足 `clientName(clientIp)` 格式
   - 当 `platName` 无法解析出 IP 时，风险判断应按高风险处理并要求人工确认

---

## 3. 存储服务对象（StorageService）

### 3.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `storageServiceId` | `string` | 主键，非自增 | 存储服务唯一标识符（去掉横杠的 UUID） | `s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890` |
| `storageServiceName` | `string` | 非空 | 存储服务名称 | `Storage-Service-01` |

### 3.2 数据源

| Type | ID | Name |
|------|-----|------|
| data_view | storageService | storageService |

### 3.3 验证规则

1. **唯一性验证**：
   - `storageServiceId` 字段必须唯一

2. **非空验证**：
   - `storageServiceId`、`storageServiceName` 字段不能为空

---

## 4. 存储池对象（StoragePool）

### 4.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `storagePoolId` | `string` | 主键，非自增 | 存储池唯一标识符（去掉横杠的 UUID） | `p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890` |
| `storagePoolName` | `string` | 非空 | 存储池名称 | `Backup-Pool-01` |
| `storageServiceId` | `string` | 非空，外键 | 关联的存储服务 ID | `s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890` |

### 4.2 数据源

| Type | ID | Name |
|------|-----|------|
| data_view | storagePool | storagePool |

### 4.3 验证规则

1. **唯一性验证**：
   - `storagePoolId` 字段必须唯一

2. **非空验证**：
   - `storagePoolId`、`storagePoolName`、`storageServiceId` 字段不能为空

3. **外键验证**：
   - `storageServiceId` 必须关联到存在的存储服务对象

---

## 5. 时间点副本对象（TimePointCopy）

### 5.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `timePointId` | `string` | 主键，非自增 | 时间点副本唯一标识符（去掉横杠的 UUID） | `b2c3d4e5f6789012b3c4d5e6f7890123` |
| `timestamp` | `int64` | 可选 | 时间戳 | `1704067200000` |
| `displayTime` | `string` | 可选 | 显示时间（年-月-日 时:分:秒 时区） | `2024-01-01 00:00:00 UTC+8` |
| `backupType` | `int` | 可选 | 备份类型 | `1` |
| `storagePoolName` | `string` | 可选 | 存储池名称 | `backup_pool_01` |
| `isClean` | `bool` | 可选 | 是否干净备份 | `true` |
| `objectId` | `string` | 非空，外键 | 关联的数据保护对象 ID | `a1b2c3d4e5f67890a1b2c3d4e5f67890` |
| `storagePoolId` | `string` | 非空，外键 | 关联的存储池 ID | `p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890` |
| `dataSetId` | `string` | 可选 | 数据集 ID | `ds1a2b3c4d5e6f7890ds1a2b3c4d5e6f7890` |

**说明**：时间点副本通过 `storagePoolId` 关联到存储池，存储池再通过 `storageServiceId` 关联到存储服务，因此时间点副本不需要单独的 `storageServiceId` 字段。

### 5.2 备份类型枚举

**备份类型（backupType）**：
- `1`：全量备份
- `2`：增量备份
- `3`：日志备份

### 5.3 数据源

**说明**：时间点副本对象无数据源定义。当前 BKN 的数据源只支持 data_view，而时间点副本对象的数据需要通过调用工具即时获取，因此不适用 data_view 类型的数据源。

**查询说明**：查询时间点副本列表需要提供以下两个参数：
1. `objectId`：数据保护对象ID
2. `storagePoolId`：存储池ID

通过存储池可以间接关联到存储服务，因此不需要单独提供存储服务ID。

### 5.4 验证规则

1. **唯一性验证**：
   - `timePointId` 字段必须唯一

2. **非空验证**：
   - `timePointId`、`objectId`、`storagePoolId` 字段不能为空

3. **外键验证**：
   - `objectId` 必须关联到存在的数据保护对象
   - `storagePoolId` 必须关联到存在的存储池对象

4. **枚举值验证**：
   - `backupType`（如果存在）必须是预定义的枚举值（1/2/3）

---

## 6. 恢复数据源对象（RecoveryDataSource）

恢复数据源对象表示选定时间点副本下的可恢复数据源集合，是时间点副本资源的进一步细化结果。

### 6.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `recovery_datasourceId` | `string` | 主键，非自增 | 恢复数据源唯一标识符（去掉横杠的 UUID） | `c3d4e5f678901234c5d6e7f890123456` |
| `display` | `string` | 可选 | 展示名称 | `database_backup_001` |
| `fullPath` | `string` | 可选 | 全路径 | `/backup/mysql/2024/01/database_backup.sql` |
| `metadata` | `string` | 可选 | 对象元数据 | `{"format": "sql", "compression": "gzip"}` |
| `objectSize` | `int64` | 可选 | 对象大小（字节） | `104857600` |
| `timePointId` | `string` | 非空，外键 | 关联的时间点副本对象 ID | `b2c3d4e5f6789012b3c4d5e6f7890123` |

### 6.2 数据源

**说明**：恢复数据源对象无数据源定义。当前 BKN 的数据源只支持 data_view，而恢复数据源对象的数据需要通过调用工具即时获取，因此不适用 data_view 类型的数据源。业务语义上，恢复数据源应按选定的时间点副本进行查询。

---

## 7. 客户端对象（Client）

### 7.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `clientId` | `string` | 主键，非自增 | 客户端唯一标识符（去掉横杠的 UUID） | `d4e5f67890123456d7e8f901234567890` |
| `clientName` | `string` | 可选 | 客户端名称 | `MySQL-Client-01` |
| `clientType` | `int` | 可选 | 客户端类型（虚拟客户端/物理客户端/其他类型） | `1` |
| `clientIp` | `string` | 可选 | 客户端 IP（IPv4 或 IPv6） | `192.168.1.100` |
| `clientMac` | `string` | 可选 | 客户端 MAC 地址 | `00:1A:2B:3C:4D:5E` |
| `type` | `int` | 可选 | 操作系统类型（Windows/Linux/Aix） | `1` |
| `platform` | `string` | 可选 | 平台类型 | `Linux` |
| `status` | `int` | 可选 | 客户端状态（在线/离线） | `1` |
| `clientVersion` | `string` | 可选 | 客户端发布版本 | `8.0.8.0.318` |

### 7.2 客户端类型枚举

**客户端类型（clientType）**：
- `1`：物理客户端
- `2`：虚拟客户端
- `3`：其他类型

**操作系统类型（type）**：
- `1`：Windows
- `2`：Linux
- `3`：Aix

**客户端状态（status）**：
- `0`：离线
- `1`：在线

### 7.3 数据源

| Type | ID | Name |
|------|-----|------|
| data_view | client | client |

---

## 8. 恢复任务对象（RecoveryTask）

恢复任务对象不是 Foundation 中稳定映射的实体表对象，而是 `execute_recovery_task` 动作使用的参数契约对象，用于承载恢复执行所需的输入参数。

### 8.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `recovery_taskId` | `string` | 主键，非自增 | 恢复任务唯一标识符（去掉横杠的 UUID） | `e5f6789012345678e9f0a123456789012` |
| `clientId` | `string` | 非空，外键 | 关联的客户端 ID | `d4e5f67890123456d7e8f901234567890` |
| `clientIp` | `string` | 非空 | 恢复目标客户端 IP | `192.168.1.100` |
| `timePointId` | `string` | 非空，外键 | 关联的时间点副本 ID | `b2c3d4e5f6789012b3c4d5e6f7890123` |
| `timestamp` | `int64` | 非空 | 选择的时间点的时间戳（毫秒） | `1704067200000` |
| `storageServiceId` | `string` | 非空，外键 | 时间点所属的存储服务 ID | `s1a2b3c4d5e6f7890s1a2b3c4d5e6f7890` |
| `storagePoolId` | `string` | 非空，外键 | 时间点所属的存储池 ID | `p1a2b3c4d5e6f7890p1a2b3c4d5e6f7890` |
| `recovery_datasourceId` | `string` | 可选，外键 | 关联的恢复数据源 ID（可选，不指定则使用整个时间点副本） | `c3d4e5f678901234c5d6e7f890123456` |
| `recovery_granularity` | `string` | 非空 | 恢复粒度 | `database` |
| `recovery_destination` | `string` | 非空 | 恢复目的地 | `original` |
| `recovery_method` | `string` | 非空 | 恢复方式 | `latest` |
| `recovery_type` | `string` | 非空 | 恢复类型 | `mount` |
| `parallel` | `int` | 可选 | 并发恢复线程数（默认 8，上限 64） | `8` |
| `recovery_user` | `string` | 可选 | 指定恢复用户 | `recovery_admin` |
| `overwrite` | `bool` | 可选 | 是否覆盖现有数据库 | `false` |
| `stop_service` | `bool` | 可选 | 恢复前是否停止数据库服务 | `false` |
| `update_credentials` | `bool` | 可选 | 恢复后是否更新用户名/密码 | `false` |
| `online` | `bool` | 可选 | 恢复后是否使数据库联机 | `true` |

### 8.2 恢复粒度枚举

**恢复粒度（recovery_granularity）**：
- `instance`：实例级恢复（恢复整个 MySQL 实例）
- `database`：库级恢复（恢复指定数据库）
- `log`：日志文件恢复（恢复 Binlog 文件）

**说明**：表级恢复不是独立的恢复粒度，而是通过库级恢复后导出表数据再导入的间接方式。

### 8.3 恢复目的地枚举

**恢复目的地（recovery_destination）**：
- `original`：原机恢复（恢复到原生产服务器）
- `remote`：异机恢复（恢复到其他主机）

**说明**：无论是原机恢复还是异机恢复，都必须部署在主机上执行恢复操作。暂不支持往 Docker 容器或 RDS 中恢复。

### 8.4 恢复方式枚举

**恢复方式（recovery_method）**：
- `fastest`：最短时间恢复（选择最近备份时间点，不回放 Binlog）
- `latest`：恢复到最新状态（备份时间点 + Binlog 回放到最新）
- `specified_time`：恢复到指定时间（备份时间点 + Binlog 回放到指定时间）

### 8.5 恢复类型枚举

**恢复类型（recovery_type）**：
- `data`：数据恢复（直接恢复数据到目标位置）
- `mount`：挂载恢复（通过挂载方式访问备份数据）

### 8.6 恢复选项说明

**恢复选项适用性**：

| 恢复选项 | 实例级恢复 | 库级恢复 | 日志文件恢复 |
|---------|-----------|---------|-------------|
| `parallel`（通道数） | ✅ 适用 | ✅ 适用 | ✅ 适用 |
| `recovery_user`（指定恢复用户） | ✅ 适用 | ✅ 适用 | ✅ 适用 |
| `overwrite`（覆盖现有数据库） | ✅ 适用 | ✅ 适用 | ❌ 不适用 |
| `stop_service`（停止服务） | ✅ 适用 | ✅ 适用 | ❌ 不适用 |
| `update_credentials`（更新密码） | ✅ 适用 | ✅ 适用 | ❌ 不适用 |
| `online`（恢复后联机） | ✅ 适用 | ✅ 适用 | ❌ 不适用 |

### 8.7 数据源

**说明**：恢复任务对象无数据源定义。当前 BKN 的数据源只支持 data_view，而恢复任务对象不映射 Foundation 实体表，只作为恢复动作的参数契约，因此不适用 data_view 类型的数据源。

### 8.8 验证规则

1. **唯一性验证**：
   - `recovery_taskId` 字段必须唯一

2. **非空验证**：
   - `recovery_taskId`、`clientId`、`timePointId`、`recovery_granularity`、`recovery_destination`、`recovery_method`、`recovery_type` 字段不能为空

3. **外键验证**：
   - `clientId` 必须关联到存在的客户端对象
   - `timePointId` 必须关联到存在的时间点副本对象
   - `recovery_datasourceId`（如果指定）必须关联到存在的恢复数据源对象

4. **枚举值验证**：
   - `recovery_granularity` 必须是预定义的枚举值（instance/database/log）
   - `recovery_destination` 必须是预定义的枚举值（original/remote）
   - `recovery_method` 必须是预定义的枚举值（fastest/latest/specified_time）
   - `recovery_type` 必须是预定义的枚举值（data/mount）

5. **数值范围验证**：
   - `parallel` 取值范围：1-64

---

## 9. 恢复作业对象（RecoveryJob）

### 9.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `jobId` | `string` | 主键，非自增 | 作业唯一标识符（去掉横杠的 UUID） | `f6789012345678901a234567890123456` |
| `objectId` | `string` | 非空，外键 | 关联的数据保护对象 ID | `a1b2c3d4e5f67890a1b2c3d4e5f67890` |
| `clientId` | `string` | 非空，外键 | 关联的客户端 ID | `d4e5f67890123456d7e8f901234567890` |
| `speed` | `string` | 可选 | 恢复速度 | `100 MB/s` |
| `startTime` | `int64` | 可选 | 恢复开始执行的时间（时间戳） | `1704067200000` |
| `endTime` | `int64` | 可选 | 恢复完成执行的时间（时间戳） | `1704070800000` |
| `completedData` | `int64` | 可选 | 恢复写入的数据量（字节） | `52428800000` |
| `sendSize` | `int64` | 可选 | 恢复作业执行过程中从服务器存储中读出来的数据量（字节） | `52428800000` |
| `businessType` | `int` | 可选 | 业务类型（默认 10，区分备份/恢复类型） | `10` |

### 9.2 业务类型枚举

**业务类型（businessType）**：
- `10`：数据恢复作业
- `11`：挂载恢复作业
- `12`：细粒度恢复作业

**说明**：businessType 用于区分作业类型，10 表示数据恢复类型，11 表示挂载恢复类型，12 表示细粒度恢复类型。

### 9.3 数据源

| Type | ID | Name |
|------|-----|------|
| data_view | job | job |

**说明**：恢复作业对象的数据源为 Foundation 的 job 表，业务知识网络从这个表中获取数据。

### 9.4 验证规则

1. **唯一性验证**：
   - `jobId` 字段必须唯一

2. **非空验证**：
   - `jobId`、`objectId`、`clientId` 字段不能为空

3. **外键验证**：
   - `objectId` 必须关联到存在的数据保护对象
   - `clientId` 必须关联到存在的客户端对象

4. **业务类型验证**：
   - `businessType` 默认值为 10（数据恢复作业）

5. **时间逻辑验证**：
   - `endTime` 必须大于或等于 `startTime`（如果两个字段都存在）

6. **数据量验证**：
   - `completedData` 和 `sendSize` 应为非负整数

---

## 10. 可用性验证对象（AvailabilityVerification）

可用性验证对象不是验证结果对象，而是 `execute_availability_verification` 动作使用的执行契约 / 执行模板对象。验证结果由动作输出返回，再由 LLM 在上下文中理解和总结。

### 10.1 数据属性

| 字段名 | 数据类型 | 约束 | 描述 | 示例值 |
|--------|----------|------|------|--------|
| `appType` | `string` | 非空 | 应用类型 | `backupengine_mysql_...` |
| `description` | `string` | 可选 | 验证说明（如何进行验证，验证结果是什么才能证明可用还是不可用） | `通过执行SQL查询验证数据库是否可正常读写` |
| `tool` | `string` | 可选 | 验证工具 | `mysql_verify_tool` |
| `params` | `string` | 可选 | 验证工具的参数（JSON字符串格式，可json.loads()解析） | `{"query": "SELECT 1", "timeout": 30}` |

### 10.2 业务逻辑说明

**应用类型关联**：
- 可用性验证对象通过 `appType` 字段按应用类型适配数据保护对象
- 根据数据保护对象的 `appType`，可以确定需要采用的验证方法和工具
- 不同的应用类型（如 MySQL、Oracle）对应不同的验证工具和验证参数

**验证流程**：
1. 根据数据保护对象查询其 `appType`
2. 根据 `appType` 确定对应的可用性验证对象
3. 获取验证工具（`tool`）和验证参数（`params`）
4. 执行可用性验证，判断恢复后的应用是否能正常提供业务
5. 将动作输出交由 LLM 结合上下文理解，不单独定义验证结果对象

### 10.3 数据源

**说明**：可用性验证对象无数据源定义。当前 BKN 的数据源只支持 data_view，而可用性验证对象不定义结果表映射，只作为验证执行契约 / 执行模板，因此不适用 data_view 类型的数据源。

### 10.4 验证规则

1. **非空验证**：
   - `appType` 字段不能为空

2. **JSON格式验证**：
   - `params` 字段如果存在，必须是有效的 JSON 字符串格式

3. **应用类型验证**：
   - `appType` 应与数据保护对象的 `appType` 保持一致

