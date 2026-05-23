# MySQL数据库恢复知识网络 - 风险模型设计

## 1. 风险模型概述

本文档定义MySQL数据库恢复知识网络的风险评估模型，用于在执行恢复操作前评估潜在风险，确保恢复操作的安全性。

### 1.1 风险评估目的

- 保护生产环境数据安全
- 防止误操作导致数据覆盖
- 提供风险提示和用户确认机制
- 支持自主恢复与人工确认的智能决策

### 1.2 风险等级定义

| 风险等级 | 说明 | 处理方式 |
|---------|------|---------|
| **high（高风险）** | 恢复操作可能覆盖原生产环境数据 | 必须用户确认后执行 |
| **low（低风险）** | 恢复操作不会影响原生产环境 | 可自主执行 |

---

## 2. 风险场景分析

### 2.1 原机恢复风险

**场景描述**：
- 恢复目标客户端IP与备份源客户端IP是同一个
- 恢复操作会覆盖原生产环境的数据
- 可能导致生产数据丢失或不可用

**风险评估**：
- **风险类型**：数据覆盖风险
- **风险等级**：high（高风险）
- **影响范围**：生产环境数据安全
- **处理方式**：必须用户确认后执行

**示例场景**：
```
备份源客户端：client-A（192.168.1.100）
恢复目标客户端：client-A（192.168.1.100）
时间点副本：timepoint-001（由client-A备份产生）

结论：原机恢复，高风险，需要用户确认
```

### 2.2 异机恢复风险

**场景描述**：
- 恢复目标客户端IP与备份源客户端IP不是同一个
- 恢复操作不会覆盖原生产环境的数据
- 在新的客户端上创建数据副本

**风险评估**：
- **风险类型**：无数据覆盖风险
- **风险等级**：low（低风险）
- **影响范围**：新客户端上的数据
- **处理方式**：可自主执行

**示例场景**：
```
备份源客户端：client-A（192.168.1.100）
恢复目标客户端：client-B（192.168.1.200）
时间点副本：timepoint-001（由client-A备份产生）

结论：异机恢复，低风险，可自主执行
```

---

## 3. 风险判断逻辑

### 3.1 判断流程

```
开始恢复任务
     ↓
获取恢复目标客户端ID（targetClientId）
     ↓
获取时间点副本ID（timePointId）
     ↓
根据 targetClientId 查询恢复目标客户端IP（targetClientIp）
     ↓
通过 timePointId 查询 protection_object.platName，并解析 sourceClientIp
     ↓
sourceClientIp 解析成功？
     ↓                    ↓
   否（失败）            是（成功）
     ↓                    ↓
高风险（high）       比较 targetClientIp 与 sourceClientIp
需要用户确认               ↓
                          ┌─────────────────────────────────────────┐
                          │ targetClientIp == sourceClientIp ?      │
                          └─────────────────────────────────────────┘
                               ↓                    ↓
                             是（相同）            否（不同）
                               ↓                    ↓
                           高风险（high）        低风险（low）
                           需要用户确认          可自主执行
                               ↓
                           用户确认通过？
                               ↓
                            是 → 执行恢复
                            否 → 取消恢复
```

### 3.2 判断规则

**规则ID**：`recovery_overwrite_risk`

**规则描述**：判断恢复操作是否会覆盖原生产环境数据

**判断条件**：
```python
targetClientIp = client.clientIp
sourceClientIp = parse_ip_from_plat_name(protection_object.platName)

if sourceClientIp is None:
    risk_level = "high"
    requires_confirmation = True
elif targetClientIp == sourceClientIp:
    risk_level = "high"
    requires_confirmation = True
else:
    risk_level = "low"
    requires_confirmation = False
```

**数据来源**：
- `targetClientId`：恢复任务参数中的 clientId
- `targetClientIp`：通过 `targetClientId -> client.clientIp` 获取
- `sourceClientIp`：通过 `timepoint_copy.objectId -> protection_object.platName` 解析得到

### 3.3 数据关系

```
recovery_task.clientId
    ↓ 关联
client.clientIp  →  targetClientIp（恢复目标客户端IP）

timepoint_copy
    ↓ objectId 关联
protection_object.platName
    ↓ parse_ip_from_plat_name
sourceClientIp（备份源客户端IP）
```

---

## 4. 风险评估配置

### 4.1 Action风险配置

**适用Action**：`execute_recovery_task`

| 配置项 | 值 | 说明 |
|-------|-----|------|
| risk_check_enabled | true | 启用风险检查 |
| risk_rule | recovery_overwrite_risk | 风险判断规则 |
| risk_on_high | confirm | 高风险时需要确认 |
| risk_on_low | auto | 低风险时自动执行 |

### 4.2 风险检查配置

```yaml
risk_check:
  enabled: true
  rules:
    - id: recovery_overwrite_risk
      name: 原机恢复风险检查
      description: 检查恢复操作是否会覆盖原生产环境数据
      condition: "sourceClientIp is None or targetClientIp == sourceClientIp"
      on_match:
        risk_level: high
        action: confirm
        message: "检测到原机恢复或无法确认备份源客户端的恢复操作，恢复可能覆盖原生产环境数据，是否确认执行？"
      on_not_match:
        risk_level: low
        action: auto
        message: "异机恢复操作，可安全执行"
```

---

## 5. 风险提示信息

### 5.1 高风险提示

**提示标题**：⚠️ 原机恢复风险警告

**提示内容**：
```
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

### 5.2 低风险提示

**提示标题**：✅ 异机恢复

**提示内容**：
```
检测到您正在执行异机恢复操作：

- 恢复目标客户端：{targetClientName}（{targetClientIp}）
- 备份源客户端：{sourceClientName}（{sourceClientIp}）
- 时间点副本：{timePointName}

此操作将在新客户端上创建数据副本，不会影响原生产环境。
```

---

## 6. 风险评估执行流程

### 6.1 执行步骤

```
1. 接收恢复任务参数
   - targetClientId
   - timePointId
   ↓
2. 查询时间点副本信息
   - 获取 timepoint_copy 记录
   ↓
3. 查询恢复目标客户端信息
   - 通过 recovery_task.clientId 获取 client.clientIp（targetClientIp）
   ↓
4. 查询数据保护对象信息
   - 通过 timepoint_copy.objectId 获取 protection_object
   - 获取 protection_object.platName
   - 解析出 sourceClientIp
   ↓
5. 确定风险等级
   - platName 无法解析出 IP → high
   - targetClientIp == sourceClientIp → high
   - targetClientIp != sourceClientIp → low
   ↓
6. 执行相应操作
   - high → 提示用户确认
   - low → 自主执行
```

### 6.2 数据查询路径

```sql
-- 查询备份源客户端 platName
SELECT po.platName
FROM timepoint_copy tc
JOIN protection_object po ON tc.objectId = po.objectId
WHERE tc.timePointId = {timePointId}
```

```python
def parse_ip_from_plat_name(plat_name: str) -> str | None:
    if not plat_name:
        return None
    match = re.search(r"\(([^()]+)\)\s*$", plat_name)
    if not match:
        return None
    candidate = match.group(1).strip()
    return candidate if is_valid_ip(candidate) else None
```

---

## 7. 后续扩展

### 7.1 可扩展的风险场景

当前只考虑原机恢复风险，后续可扩展：

| 风险场景 | 风险等级 | 说明 |
|---------|---------|------|
| 数据库版本不兼容 | medium | 目标数据库版本与备份版本不兼容 |
| 存储空间不足 | medium | 目标客户端存储空间不足以存放恢复数据 |
| 数据库正在运行 | low | 目标数据库服务正在运行 |
| 网络不稳定 | low | 客户端网络连接不稳定 |

### 7.2 风险组合评估

后续可支持多风险组合评估：

```python
def evaluate_risks(context):
    risks = []
    
    # 原机恢复风险
    if check_overwrite_risk(context):
        risks.append(Risk("overwrite", "high", "原机恢复将覆盖生产数据"))
    
    # 版本兼容性风险
    if check_version_compatibility(context):
        risks.append(Risk("version", "medium", "数据库版本不兼容"))
    
    # 存储空间风险
    if check_storage_space(context):
        risks.append(Risk("storage", "medium", "存储空间不足"))
    
    return aggregate_risks(risks)
```

---

## 8. 设计总结

### 8.1 核心设计原则

1. **安全优先**：默认采用最严格的风险评估标准
2. **用户友好**：提供清晰的风险提示和确认机制
3. **可扩展**：支持后续添加更多风险场景
4. **自动化**：低风险场景支持自主执行，提高效率

### 8.2 当前实现范围

- ✅ 原机恢复风险检查
- ✅ 高风险需要用户确认
- ✅ 低风险可自主执行
- ✅ 清晰的风险提示信息

### 8.3 待实现功能

- [ ] 多风险组合评估
- [ ] 风险历史记录
- [ ] 风险评估报告
- [ ] 自定义风险规则
