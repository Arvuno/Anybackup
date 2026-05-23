# 备份运行知识网络 - Agent 使用指南

> **网络ID**: backup_run_knowledge_network
> **版本**:

## 网络概览

### 核心对象

| 对象 | 文件路径 | 说明 |
|------|----------|------|
| 备份客户端配置 | `object_types/backup_client_config.bkn` | 用于描述保护对象公共配置下绑定的客户端配置项，包括关联客户端、所属区域和公共配置标识。 |
| 备份基础配置 | `object_types/backup_config.bkn` | 用于描述保护对象通用备份配置，覆盖存储、压缩、加密、重试和一致性校验等公共参数；当前使用关联存储池名称作为主要展示字段。 |
| 备份区域配置 | `object_types/backup_region_config.bkn` | 用于描述保护对象公共配置下的区域级配置，包括区域名称、故障切换开关和失败统计信息。 |
| 备份存储池配置 | `object_types/backup_storage_pool_config.bkn` | 用于描述保护对象公共配置下各存储池的细分配置项，包括压缩、加密、重删和所属区域等参数。 |
| 客户端 | `object_types/client.bkn` | 用于描述参与备份任务执行的客户端主数据，包括客户端名称、类型、状态和机器标识信息。 |
| 备份应用配置(mysql) | `object_types/mysql_backup_config.bkn` | 用于描述 MySQL 和 MariaDB 类型保护对象的专用备份配置；当前使用数据源路径作为主要展示字段。 |
| 生产资源 | `object_types/production_resource.bkn` | 用于描述被纳入备份管理的生产资源目录，记录资源名称、类型、归属生产系统及应用类型等信息。 |
| 生产系统 | `object_types/production_system.bkn` | 用于描述生产系统主数据，作为生产资源与保护对象归属关系中的上层系统实体。 |
| 保护对象 | `object_types/protected_object.bkn` | 用于描述已经纳入保护范围的对象实例，记录对象名称、所属系统、资源类型以及最近一次备份状态。 |
| 存储池 | `object_types/storage_pool.bkn` | 用于描述备份使用的存储池主数据，包括容量、状态以及所属存储服务。 |
| 存储服务 | `object_types/storage_service.bkn` | 用于描述备份相关的存储服务主数据，包括服务名称、类型、状态与唯一标识。 |

### 核心关系

| 关系 | 文件路径 | 说明 |
|------|----------|------|
| 备份配置使用客户端 | `relation_types/backup_config_to_client.bkn` | 描述通用备份配置与客户端之间的归属关系，来源于 `backup_config.client_id`。 |
| 备份配置依赖区域配置 | `relation_types/backup_config_to_region_cfg.bkn` | 描述备份配置与区域配置之间的依赖关系；当 `backup_config.config_version = 1` 时，备份配置会依赖对应的区域配置，关联键来源于 `backup_region_config.config_id`。 |
| 备份配置使用存储池 | `relation_types/backup_config_to_storage_pool.bkn` | 描述通用备份配置与存储池之间的关联关系，来源于 `backup_config.storage_pool_id`。 |
| 备份配置使用存储服务 | `relation_types/backup_config_to_storage_service.bkn` | 描述通用备份配置与存储服务之间的关联关系，来源于 `backup_config.storage_service_id`。 |
| 备份区域配置包含客户端配置 | `relation_types/backup_region_to_client_cfg.bkn` | 描述区域级备份配置与客户端配置之间的从属关系，来源于 `backup_client_config.region_id`；用于表达多区域场景下某一区域绑定了哪些备份客户端。 |
| 备份区域配置包含存储池配置 | `relation_types/backup_region_to_pool_cfg.bkn` | 描述区域级备份配置与存储池配置之间的从属关系，来源于 `backup_storage_pool_config.region_id`；用于表达多区域场景下某一区域下配置了哪些备份存储池。 |
| 生产资源生成保护对象 | `relation_types/prod_resource_gen_protected_object.bkn` | 描述生产资源在被纳入保护后生成保护对象的关系。 |
| 生产系统拥有生产资源 | `relation_types/prod_system_has_prod_resource.bkn` | 描述生产系统与其下生产资源之间的拥有关系，来源于 `production_resource.productionSystemId`。 |
| 生产系统属于父生产系统 | `relation_types/prod_system_to_parent_system.bkn` | 描述生产系统之间的层级关系，来源于 `resource_info.parent_id`。 |
| 保护对象生成MySQL配置 | `relation_types/protected_object_gen_mysql_config.bkn` | 描述生产资源完成备份配置并生成保护对象后，若该对象对应 MySQL 资源，则会进一步生成 MySQL 专属备份配置。 |
| 保护对象对应备份配置 | `relation_types/protected_object_has_backup_config.bkn` | 描述保护对象与其对应备份配置之间的关系，来源于 `protect_object.backup_config_id`。 |
| 保护对象属于生产系统 | `relation_types/protected_object_to_prod_system.bkn` | 描述保护对象与所属生产系统之间的归属关系，来源于 `protect_object.production_system_id`。 |
| 存储池属于存储服务 | `relation_types/storage_pool_to_storage_service.bkn` | 描述存储池与存储服务之间的从属关系，来源于 `pool_v8.storage_svc_id`。 |

### 可用行动

| 行动 | 文件路径 | 说明 |
|------|----------|------|
| 查询生产资源的备份方案 | `action_types/query_prod_resource_backup_cfg.bkn` | 根据生产资源名称或生产资源ID查询当前已生效的备份方案。
1. 查询时先定位 `production_resource`
2. 再判断是否关联 `protected_object`；若未关联，则直接返回空结果。
3. 若已关联，则读取 `backup_config`
    3.1 backup_config.config_version=1
        3.1.1 `backup_config -> backup_region_config -> backup_client_config / backup_storage_pool_config` 的层级关系展开基础配置；
        3.1.2 按应用类型决定是否补充 `mysql_backup_config` 等应用配置。
    3.1 backup_config.config_version=0
        3.1 backup_config 作为基础配置
        3.2 按应用类型决定是否补充 `mysql_backup_config` 等应用配置。

最终输出应组织为两部分：
- `基础配置`：返回公共配置本身，以及其下的 `备份区域配置`；每个区域下继续返回 `备份客户端` 和 `备份存储池`
- `应用配置`：返回与当前保护对象匹配的应用专用配置

查询输入支持 `production_resource_name` 或 `production_resource_id`，至少提供一个；若两者同时提供，优先按 `production_resource_id` 精确定位。结果展示优先使用对象 `Display Name`，且不显示来源对象、字段英文名或底层表行细节。 |
| 查询保护对象备份配置 | `action_types/query_protected_object_backup_config.bkn` | 根据保护对象名称查询该保护对象的备份配置，并直接返回当前已生效的“备份方案内容”。输出口径应与 `backup_new_knowledge` 中的 `backup_solution`、`basic_config`、`app_config` 保持一致：优先输出可读的方案摘要、基础配置和应用配置，而不是仅返回底层配置行。

结果中的展示标签应优先使用各对象 `Data Properties` 中定义的 `Display Name`，而不直接使用字段名；例如 `client_name` 应显示为“客户端名称”。最终结果不应显示任何来源说明。 |
| 查询存储池 | `action_types/query_storage_pool.bkn` | 根据输入的存储服务 ID 查询当前运行态中的存储池，并直接返回面向用户展示的存储池结果。

当前 Action 面向存储池目录查询场景。输入 `storage_service_id` 后，应先定位该存储服务，再沿 `storage_pool_to_storage_service` 关系筛选属于该存储服务的存储池。结果展示标签应优先使用对象 `storage_pool` 的 `Display Name`，并将容量类字段按业务语义组织为“存储空间大小”“存储剩余空间”。

除返回全部存储池外，本 Action 还支持通过输入参数控制结果选择策略。例如可通过 `result_mode = max_remaining_space` 只返回“存储剩余空间最大”的存储池；若存在多条并列最大值，则按 `存储池ID` 升序返回第一条，保证结果稳定。 |
| 查询存储服务 | `action_types/query_storage_service.bkn` | 根据输入条件查询当前运行态中的存储服务，并直接返回面向用户展示的存储服务结果。

当前 Action 主要面向存储服务目录查询场景。若未提供筛选条件，则返回全部存储服务；若提供 `storage_service_id` 或 `storage_service_name`，则按输入做精确过滤。结果展示标签应优先使用对象 `storage_service` 的 `Display Name`，即“存储服务ID”“存储服务名称”“存储服务类型”，不直接暴露底层字段名。 |

## 目录结构

```
.
├── network.bkn
├── SKILL.md
├── CHECKSUM
├── object_types/
├── relation_types/
└── action_types/
```

## 使用建议

### 查询场景

1. **获取所有对象定义**
   - 查看 `object_types/` 目录下的文件

2. **查找关系定义**
   - 查看 `relation_types/` 目录下的文件

### 运维场景

1. **执行运维操作**
   - 查看 `action_types/` 目录下的行动定义
   - 了解触发条件和参数绑定

## 索引表

### 按类型索引

- **对象定义**: `object_types/`
- **关系定义**: `relation_types/`
- **行动定义**: `action_types/`

## 注意事项

1. 本网络由 BKN SDK 自动生成 SKILL.md
2. 所有定义遵循 BKN 规范
3. 使用 CHECKSUM 文件验证网络完整性
