# MySQL 数据库恢复知识网络 - Relation 数据模型设计

**版本**：v1.0  
**文档类型**：关系模型设计文档  
**编写作者**：Andrew

---

## 1. 关系模型概述

本文档定义了 MySQL 数据库恢复知识网络（BKN）中 Relation（关系）的数据模型设计，包括对象之间的关系定义、映射规则等。

**关系方向说明**：`Source → Target` 表示先有源对象，然后源对象通过属性指向目标对象进行操作。

**关系总览**：
- 共定义了10个关系类型
- 包括存储池与存储服务的关系
- 包括时间点副本与数据保护对象、存储池的关系
- 包括恢复任务与客户端、时间点副本、恢复数据源的关系
- 包括恢复作业与数据保护对象、客户端的关系
- 包括可用性验证与数据保护对象的关系

**说明**：时间点副本通过存储池间接关联到存储服务，因此不需要单独定义时间点副本与存储服务的关系。

---

## 2. 时间点副本与数据保护对象关系

### 2.1 关系定义

**关系ID**：`timepoint_belongs_object`  
**关系名称**：时间点副本属于数据保护对象  
**关系描述**：时间点副本对象通过 objectId 字段关联到数据保护对象，表示该时间点副本属于哪个数据保护对象。

### 2.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| timepoint_copy | protection_object | direct |

### 2.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| objectId | objectId |

---

## 2.5. 时间点副本与存储池关系

### 2.5.1 关系定义

**关系ID**：`timepoint_belongs_pool`  
**关系名称**：时间点副本属于存储池  
**关系描述**：时间点副本对象通过 storagePoolId 字段关联到存储池对象，表示该时间点副本存储在哪个存储池中。

### 2.5.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| timepoint_copy | storage_pool | direct |

### 2.5.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| storagePoolId | storagePoolId |

---

## 3. 恢复数据源与时间点副本关系

### 3.1 关系定义

**关系ID**：`datasource_belongs_timepoint`  
**关系名称**：恢复数据源属于时间点副本  
**关系描述**：恢复数据源对象通过 timePointId 字段关联到时间点副本对象，表示该恢复数据源属于哪个时间点副本。

### 3.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_datasource | timepoint_copy | direct |

### 3.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| timePointId | timePointId |

---

## 4. 恢复任务与客户端关系

### 4.1 关系定义

**关系ID**：`task_specifies_client`  
**关系名称**：恢复任务明确客户端  
**关系描述**：恢复任务对象通过 clientId 字段关联到客户端对象，明确恢复任务使用哪个客户端执行恢复操作。

### 4.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_task | client | direct |

### 4.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| clientId | clientId |

---

## 5. 恢复任务与时间点副本关系

### 5.1 关系定义

**关系ID**：`task_specifies_timepoint`  
**关系名称**：恢复任务明确时间点副本  
**关系描述**：恢复任务对象通过 timePointId 字段关联到时间点副本对象，明确恢复任务使用哪个时间点副本进行恢复。

### 5.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_task | timepoint_copy | direct |

### 5.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| timePointId | timePointId |

---

## 6. 恢复任务与恢复数据源关系

### 6.1 关系定义

**关系ID**：`task_specifies_datasource`  
**关系名称**：恢复任务明确恢复数据源  
**关系描述**：恢复任务对象通过 recovery_datasourceId 字段关联到恢复数据源对象，明确恢复任务使用哪个恢复数据源进行恢复（可选关系）。

### 6.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_task | recovery_datasource | direct |

### 6.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| recovery_datasourceId | recovery_datasourceId |

---

## 7. 恢复作业与数据保护对象关系

### 7.1 关系定义

**关系ID**：`job_belongs_object`  
**关系名称**：恢复作业属于数据保护对象  
**关系描述**：恢复作业对象通过 objectId 字段关联到数据保护对象，表示该恢复作业是针对哪个数据保护对象执行的恢复操作。

### 7.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_job | protection_object | direct |

### 7.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| objectId | objectId |

---

## 8. 恢复作业与客户端关系

### 8.1 关系定义

**关系ID**：`job_specifies_client`  
**关系名称**：恢复作业明确客户端  
**关系描述**：恢复作业对象通过 clientId 字段关联到客户端对象，明确恢复作业在哪个客户端上接收和写入恢复数据。

### 8.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| recovery_job | client | direct |

### 8.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| clientId | clientId |

---

## 9. 存储池与存储服务关系

### 9.1 关系定义

**关系ID**：`pool_belongs_service`  
**关系名称**：存储池属于存储服务  
**关系描述**：存储池对象通过 storageServiceId 字段关联到存储服务对象，表示该存储池属于哪个存储服务。

### 9.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| storage_pool | storage_service | direct |

### 9.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| storageServiceId | storageServiceId |

---

## 10. 可用性验证按应用类型适配关系

### 10.1 关系定义

**关系ID**：`verification_applies_to_app_type`  
**关系名称**：可用性验证按应用类型适配  
**关系描述**：可用性验证对象通过 appType 字段按应用类型适配数据保护对象，根据数据保护对象的应用类型确定需要采用的可用性验证方法和工具。

### 10.2 Endpoint

| Source | Target | Type |
|--------|--------|------|
| availability_verification | protection_object | direct |

### 10.3 Mapping Rules

| Source Property | Target Property |
|-----------------|-----------------|
| appType | appType |

---

## 11. 关系总览

| 关系ID | 关系名称 | 源对象 | 目标对象 | 映射字段 |
|--------|----------|--------|----------|----------|
| `timepoint_belongs_object` | 时间点副本属于数据保护对象 | timepoint_copy | protection_object | objectId → objectId |
| `timepoint_belongs_pool` | 时间点副本属于存储池 | timepoint_copy | storage_pool | storagePoolId → storagePoolId |
| `datasource_belongs_timepoint` | 恢复数据源属于时间点副本 | recovery_datasource | timepoint_copy | timePointId → timePointId |
| `task_specifies_client` | 恢复任务明确客户端 | recovery_task | client | clientId → clientId |
| `task_specifies_timepoint` | 恢复任务明确时间点副本 | recovery_task | timepoint_copy | timePointId → timePointId |
| `task_specifies_datasource` | 恢复任务明确恢复数据源 | recovery_task | recovery_datasource | recovery_datasourceId → recovery_datasourceId |
| `job_belongs_object` | 恢复作业属于数据保护对象 | recovery_job | protection_object | objectId → objectId |
| `job_specifies_client` | 恢复作业明确客户端 | recovery_job | client | clientId → clientId |
| `pool_belongs_service` | 存储池属于存储服务 | storage_pool | storage_service | storageServiceId → storageServiceId |
| `verification_applies_to_app_type` | 可用性验证按应用类型适配 | availability_verification | protection_object | appType → appType |

---

## 12. 关系图示

```
StorageService (存储服务)
    ↑ storageServiceId
    |
    +-- StoragePool (存储池)
            ↑ storagePoolId
            |
            | (备份数据存储位置)
            |
ProtectionObject (数据保护对象)
    ↑ objectId
    |
    +-- TimePointCopy (时间点副本)
    |       | objectId (关联数据保护对象)
    |       | storagePoolId (关联存储池，通过存储池间接关联存储服务)
    |       ↑ timePointId
    |       |
    |       +-- RecoveryDataSource (恢复数据源)
    |               ↑ recovery_datasourceId
    |               |
    |               +-- RecoveryTask (恢复任务)
    |                       | clientId
    |                       | timePointId
    |                       | recovery_datasourceId
    |                       ↓
    |                   Client (客户端)
    |
    +-- RecoveryJob (恢复作业)
    |       | objectId
    |       | clientId
    |       ↓
    |   Client (客户端)
    |
    +-- AvailabilityVerification (可用性验证)
            | appType
            ↓
        ProtectionObject (数据保护对象)
```
